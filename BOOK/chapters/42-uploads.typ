#import "../template.typ": *

= File uploads and the image pipeline

D9 was "a user uploads a ten megabyte photo from their phone, and it is served
to every visitor at ten megabytes". It is the most predictable performance
problem in web software, it is never the thing anybody plans for, and it has a
second half nobody mentions: that photo arrives with the GPS coordinates of the
house it was taken in.

This chapter is what a file field does, what happens to an image between the
form and the disk, and the two attacks that live in this feature.

== A file field, and what it stores

```bash
grit generate resource Product \
  --fields "name:string, photo:file, gallery:files"
```

`file` is one, `files` is many. The column does not hold a path. It holds a
structured reference:

```json
{
  "url": "https://cdn.example.com/u/2026/10/x7f3.webp",
  "key": "u/2026/10/x7f3.webp",
  "name": "kitchen-final.jpg",
  "mime": "image/webp",
  "size": 147204,
  "width": 1600,
  "height": 1067,
  "thumbnail_url": ".../x7f3_thumb.webp",
  "format": "webp",
  "optimised": true
}
```

#bridge[
  A bare URL string is what everybody stores first, and it is missing exactly
  the fields a page needs.

  Without `width` and `height` the browser cannot reserve space, so the layout
  jumps as images load, which is most of a bad Cumulative Layout Shift score.
  Without `name` you cannot show the user what they uploaded. Without `mime`
  you are guessing from the extension. Without `thumbnail_url` a list page
  downloads full-size images to draw at 60 pixels.

  `optimised: false` is the one people miss: it says the pipeline failed and
  the original was stored untouched, which is a thing you want to be able to
  find later.
]

== Where the bytes go

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [S3, R2, Spaces], [Production. Any S3-compatible object store],
  [MinIO], [The same API, in Docker, for local development],
  [Local disk], [`STORAGE_DRIVER=local`. No object store at all],
)

The local driver is the one to use while learning. It needs nothing running,
and uploads land in a directory the API serves.

#gotcha[
  With `STORAGE_DRIVER=local`, every stored file's URL is built from `APP_URL`.
  Change `APP_PORT` without changing `APP_URL` and every upload gets a URL
  pointing at a port nothing is listening on. The upload succeeds, the file is
  on disk, and the image is broken, which sends people hunting through the
  storage code for a bug that is in an environment variable.

  On a phone or another machine, the same applies to `localhost`: an image URL
  that says `localhost` resolves to the phone itself. Mobile clients rewrite
  the host for this reason, and MinIO has to be published on the LAN address
  rather than `127.0.0.1`.
]

== What happens to an image on the way in

Every image field runs through a profile. The default one is a set of decisions
somebody would otherwise have to research:

```go
func DefaultProfile() Profile {
    return Profile{
        Max:        Fit(1600, 1600),
        Quality:    0.82,
        Format:     Auto,
        Renditions: map[string]Size{"thumb": Fill(400, 400)},
        OnError:    StoreOriginal,
        MaxPixels:  50_000_000,
    }
}
```

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [1600], [Covers a 2x retina display at a typical content width. Storing more
    is paying to keep pixels no browser will draw],
  [0.82], [The point at which JPEG is visually indistinguishable from source.
    Below about 0.75, gradients and skin show artifacts],
  [400], [The thumbnail an admin row or a card grid needs at 2x],
  [Auto], [Encode to the best format the client can take, usually WebP],
)

Measured on a 5.3 MB phone photo: *147 KB out, including the thumbnail*. That
is D9 closed, and it happened without anybody remembering to do it.

Three more things happen in the same pass, and two of them are not about size.

*EXIF orientation is applied.* A portrait photo from a phone is usually stored
landscape with a rotation flag, and a naive decode produces a sideways
thumbnail.

*Metadata is stripped.* Including GPS coordinates. A photograph taken at home
and uploaded to a public listing carries the house's location in a field nobody
looks at, and stripping it is not a feature request anybody files: it is
something a framework should do by default.

*Renditions are derived.* The thumbnail comes from the decoded original rather
than from the resized copy, so it is not a resize of a resize.

== Your own profile, when the default is wrong

```go
media.Define("product-image", media.Profile{
    Max:        media.Fit(2400, 2400),
    Quality:    0.9,
    Renditions: map[string]media.Size{
        "thumb": media.Fill(400, 400),
        "card":  media.Fit(800, 0),
    },
})
```

A clothing retailer wants more pixels and higher quality than a blog does. A
document scan wants no lossy compression at all. The profile is the one place
that decision lives, rather than being spread across the upload handler, the
list page and whatever the mobile app does.

The chainable API is there when you need something bespoke:

```go
img, err := media.Open(file)
thumb, err := img.Cover(150, 150).ToWebP().Encode()
large, err := img.Fit(1600, 0).Quality(80).ToJPEG().Encode()
```

An operation that fails does not panic and does not return its own error. The
error travels down the chain and comes out of `Encode` or `Store`, which is
what makes the chain readable.

== The decompression bomb, which the size limit cannot see

```go
MaxPixels: 50_000_000
```

This is the line that stops a specific attack, and it is worth understanding
because the obvious defence does not work.

A solid-colour PNG compresses to almost nothing whatever its dimensions. *165
kilobytes on the wire can be 144 megapixels decoded.* Measured, one such upload
allocated 224 MB, and ten concurrent uploads would have been 2.2 GB, which is
the process gone.

Your file size limit cannot see it, because on disk the file is small. The
header is read first, the pixel count is computed before any pixels are
allocated, and an image claiming more than the profile allows is refused
without being decoded.

#gotcha[
  Fifty megapixels is chosen to pass a 48 megapixel professional camera frame
  and refuse the bomb. If you raise it, raise it with a number in mind rather
  than to make one complaint go away.

  The same reasoning applies to the upload cap of 50 MB. Limits that exist for
  a reason should carry the reason in a comment, or the next person removes
  them during a demo.
]

== When the pipeline fails

```go
OnError: StoreOriginal
```

A transform that fails stores the original untouched and marks the ref
`optimised: false`. The alternative, `Reject`, refuses the upload.

Storing a large image is a worse outcome than storing a small one, and both are
much better than losing the user's file because your encoder did not like
something. The failures are findable afterwards, and a reprocess pass can retry
them.

#gotcha[
  `DiscardOriginal` is deliberately phrased in the negative, and the story is a
  general lesson about defaults in Go.

  It began as `KeepOriginal`, which read better. But profiles are merged with
  the defaults field by field, and a `bool` has no way to say "unset": every
  profile that did not mention `KeepOriginal` had it as `false`, so the
  original was silently discarded while the documentation promised the
  opposite.

  The rule: *in a struct that gets merged with defaults, the zero value must be
  the recommended behaviour*. Where it cannot be, use a pointer and treat
  `nil` as unset.
]

== Uploading from the frontend

The admin has a dropzone component, and the generated form uses it for any
`file` or `files` field. Behind it is an ordinary multipart POST:

```ts
const body = new FormData()
body.append("file", file)

const res = await fetch(`${API}/api/v1/uploads`, {
  method: "POST",
  body,                       // no Content-Type header
  credentials: "include",
})
```

#gotcha[
  Do not set `Content-Type` yourself on a multipart request. The browser adds
  it *with the boundary parameter*, and a hand-written
  `multipart/form-data` header has no boundary, so the server cannot parse the
  body.

  This is the single most common upload bug in every stack, and on React Native
  it has a sibling: some upload helpers send an empty body, and the fix is the
  same `fetch` plus `FormData` with no content type.
]

For very large files, the direct-to-storage path skips the API: the client asks
for a presigned URL, uploads straight to the bucket, and tells the API the key.
The image helpers run the same transform afterwards, so a presigned upload is
not a way around the pipeline.

== Serving them

Three things make the difference between a fast page and a slow one, and none
of them is the upload:

*Use the thumbnail in lists.* A table of fifty rows showing full-size images is
fifty megabytes for fifty thumbnails.

*Give the browser the dimensions.* They are in the ref, so pass them to the
image component and the layout stops jumping.

*Let the CDN cache it.* Objects are content-addressed by key, so a long cache
lifetime is safe: a changed image is a new key.

#project(title: "Project 39: measure the pipeline")[
  Numbers, not impressions.

  ```bash
  grit new gallery && cd gallery
  grit generate resource Photo \
    --fields "caption:string, image:file, album:belongs_to:Album"
  grit generate resource Album --fields "name:string"
  grit migrate && grit start
  ```

  *One: the ratio.* Upload a real photo straight off a phone. Record the file
  size before, the stored size after, and the thumbnail size. Compute the
  ratio and put all three in `NOTES.md`.

  *Two: the metadata.* Before uploading, read the EXIF of your photo:

  ```bash
  exiftool kitchen.jpg | grep -i "gps\|orientation"
  ```

  Then read the stored file. The GPS fields should be gone and the image should
  be the right way up.

  *Three: the bomb.* Generate a large solid-colour PNG and try to upload it:

  ```bash
  python -c "from PIL import Image; \
Image.new('RGB',(12000,12000),'blue').save('bomb.png')"
  ```

  Record the file size on disk and what the API says. Then raise `MaxPixels`,
  retry, and watch the memory of the API process while it decodes.

  *Four: a profile of your own.* Define a profile for `Photo.image` with a
  second rendition at 800 pixels wide, re-upload, and confirm both renditions
  exist in storage.

  *Five: the failure path.* Upload a file with a `.jpg` extension that is not
  an image. Confirm the ref says `optimised: false` and that the row still
  exists, rather than the request failing with a 500.

  *Six: the page.* Render the album with thumbnails and dimensions set, then
  again without, and compare the layout shift in Lighthouse.
]

#tryit[
  Switch `STORAGE_DRIVER` to `local`, upload a file, and find it on disk. Then
  change `APP_PORT` without changing `APP_URL` and reload the page.

  Seeing the broken image once, with the cause known, is worth an hour of
  debugging it later with the cause unknown.
]

#assignment(name: "Uploads - Assignment One")[
  In `part5-building/42-uploads`, build the gallery and report the pipeline
  honestly.

  *One.* The resources, the upload working through the admin and through a
  `fetch` you wrote yourself.

  *Two.* `PIPELINE.md` with a table of five real photos: original size, stored
  size, thumbnail size, dimensions in and out, and the format chosen.

  *Three.* The EXIF evidence, before and after, with the GPS fields shown
  present and then absent.

  *Four.* The decompression bomb, with the file size on disk next to the
  decoded megapixels, and what the API did. One sentence on why a file size
  limit alone would not have saved you.

  *Five.* One custom profile with a justification: name the screen that needs
  the extra pixels, and the trade you are making in storage to get them.
]
