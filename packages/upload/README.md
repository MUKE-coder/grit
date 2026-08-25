# @gritframework/upload

Client-side image optimisation and direct-to-storage uploads, for React, Next.js and Expo.

Your user picks a 6 MB photo off their phone. This shrinks it to about 35 KB **before it leaves the device**, then uploads it straight to S3 through a presigned URL. Your server never sees the bytes: no upload bandwidth, no CPU, no request timeouts.

```
5.02 MB  4000x3000 JPEG        what the user picked
  31.6 KB  1600x1200 WebP      what gets stored
   3.4 KB   400x400  WebP      the thumbnail, alongside it
-------------------------------------------------------
   147x smaller, measured in Chromium and Mobile Chrome
```

That is not a compromise against doing it on a server. The browser has a lossy WebP encoder built in, so it matches [libvips](https://www.libvips.org/) (33.9 KB on the same image) without any native dependency at all.

## Install

```bash
npm install @gritframework/upload
```

Expo also needs the manipulator, because React Native has no canvas:

```bash
npx expo install expo-image-manipulator
```

## Use

The optimiser is **injected rather than imported**, so the same uploader works on Next.js, a Vite SPA and Expo, and no bundler has to resolve a platform at build time.

```ts
import { createUploader, createFetchTransport } from "@gritframework/upload";
import { optimizeImage } from "@gritframework/upload/web";

export const uploader = createUploader({
  transport: createFetchTransport("/api/v1", () => localStorage.getItem("token")),
  optimize: optimizeImage,
});

const ref = await uploader.upload(file, file.name, { profile: "product-image" });
// ref.url, ref.thumbnail_url, ref.renditions.thumb, ref.format, ref.optimised
```

### React

```tsx
import { useUpload, describeSaving } from "@gritframework/upload/react";

function Upload() {
  const { upload, items } = useUpload(uploader);

  return (
    <>
      <input
        type="file"
        multiple
        onChange={(e) =>
          upload([...(e.target.files ?? [])].map((f) => ({ blob: f, name: f.name })))
        }
      />
      {items.map((it) => (
        <div key={it.id}>
          {it.name} — {it.state} {Math.round(it.progress * 100)}%
          {describeSaving(it) && <em> {describeSaving(it)}</em>}
        </div>
      ))}
    </>
  );
}
```

`useUpload` is deliberately not a TanStack Query mutation, though it composes with one. An upload has per-file progress and partial success; `useMutation` models one request with one result. Use this for the upload and TanStack Query for the record the refs end up attached to.

### Expo

```ts
import { optimizeNativeImage } from "@gritframework/upload/expo";

const picked = await ImagePicker.launchImageLibraryAsync();
const optimize = (_blob: Blob, profile) => optimizeNativeImage(picked.assets[0], profile);
```

## What your API needs to provide

Three endpoints. If you use [Grit](https://gritframework.dev) they already exist.

| Endpoint | Purpose |
| --- | --- |
| `GET /media/profiles` | The optimisation profiles, so the client uses your numbers rather than its own copy |
| `POST /uploads/presign` | Takes `{filename, content_type, file_size}`, returns `{presigned_url, key, public_url}` |
| `POST /uploads/complete` | Records the row once the bytes are in storage |

Two things your presign endpoint should do, because moving optimisation to the client changes the server's job from transforming to constraining:

**Sign the exact byte count into the URL.** The client optimises first, so it knows its size before it asks. Without this the URL is an unbounded write capability: a client can ask to upload two megabytes and then send five gigabytes.

**Re-read the object on completion** rather than trusting the reported size. The bytes never came through your API, so every number in that request is a claim.

## Profiles

Fetched from `GET /media/profiles`, with a sensible built-in fallback if that request fails.

```ts
{
  name: "default",
  max_width: 1600, max_height: 1600, crop: false,
  quality: 0.82,
  format: "auto",
  max_pixels: 50_000_000,
  renditions: { thumb: { width: 400, height: 400, crop: true } },
}
```

`format: "auto"` picks per image rather than making you decide: anything with real transparency stays lossless, everything else goes lossy. A transparent logo can never be flattened onto a black JPEG background.

`max_pixels` refuses a decompression bomb. A solid-colour PNG compresses to almost nothing whatever its dimensions, so a 165 KB file can be 144 megapixels decoded, which will hang or kill a mobile browser tab.

## Details that are easy to get wrong

These are handled, and are the reason this is a package rather than forty lines of canvas code:

- **EXIF orientation** is applied on decode, or portrait photos come out sideways. It is then stripped by re-encoding, which matters because phone photos carry GPS coordinates.
- **`canvas.toBlob` silently returns a PNG** when asked for an unsupported type instead of failing. Asking for WebP and getting PNG back would make photos *larger* than the JPEG they replaced, and nothing would look broken. The result is checked.
- **Transparency detection** uses a threshold just under opaque, because resampling leaves alpha a hair below 255 and treating that as transparency sends every photo down the lossless path at many times the size.
- **`Fit` never upscales.** Enlarging a small upload spends bytes on pixels the source never had.
- **GIF is skipped**, because drawing one to a canvas keeps the first frame only.
- **An undecodable file is uploaded untouched** rather than lost. Losing somebody's file because a codec choked is the worse failure.

## Licence

MIT
