package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeUploadPackageFiles writes packages/upload: client-side image
// optimisation and direct-to-storage uploads.
//
// Uploads go from the browser straight to S3 via a presigned URL and never
// pass through the API, so the optimisation has to happen before the bytes
// leave the device. That is not only where it has to be, it is where it is
// best: the browser has a lossy WebP encoder built in, which is exactly what
// the pure-Go server backend cannot do without cgo.
func writeUploadPackageFiles(root string, opts Options) error {
	pkgRoot := filepath.Join(root, "packages", "upload")

	files := map[string]string{
		filepath.Join(pkgRoot, "package.json"):        uploadPackageJSON(),
		filepath.Join(pkgRoot, "tsconfig.json"):       uploadTSConfig(),
		filepath.Join(pkgRoot, "src", "types.ts"):     uploadTypesTS(),
		filepath.Join(pkgRoot, "src", "index.ts"):     uploadIndexTS(),
		filepath.Join(pkgRoot, "src", "web.ts"):       uploadWebTS(),
		filepath.Join(pkgRoot, "src", "expo.ts"):      uploadExpoTS(),
		filepath.Join(pkgRoot, "src", "uploader.ts"):  uploadUploaderTS(),
		filepath.Join(pkgRoot, "src", "react.ts"):     uploadReactTS(),
		filepath.Join(pkgRoot, "src", "transport.ts"): uploadTransportTS(),
		filepath.Join(pkgRoot, "src", "web.test.ts"):  uploadWebTestTS(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", opts.Module())
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func uploadPackageJSON() string {
	return `{
  "name": "@repo/upload",
  "version": "0.0.0",
  "private": true,
  "type": "module",
  "main": "./src/index.ts",
  "types": "./src/index.ts",
  "exports": {
    ".": "./src/index.ts",
    "./web": "./src/web.ts",
    "./transport": "./src/transport.ts",
    "./expo": "./src/expo.ts",
    "./react": "./src/react.ts"
  },
  "scripts": {
    "test": "vitest run",
    "typecheck": "tsc --noEmit"
  },
  "peerDependencies": {
    "react": ">=18"
  },
  "peerDependenciesMeta": {
    "react": { "optional": true }
  },
  "devDependencies": {
    "typescript": "^5.6.3",
    "vitest": "^2.1.4"
  }
}
`
}

func uploadTSConfig() string {
	return `{
  "compilerOptions": {
    "target": "ES2022",
    "lib": ["ES2022", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "strict": true,
    "noEmit": true,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "jsx": "react-jsx"
  },
  "include": ["src/**/*.ts"]
}
`
}

func uploadTypesTS() string {
	return `/**
 * Shapes shared by every platform backend.
 *
 * MediaProfile mirrors what GET /api/v1/media/profiles returns. Fetch it
 * rather than hardcoding numbers: the server has the same values, and two
 * copies drift the first time one of them changes.
 */

export type MediaFormat = "auto" | "jpeg" | "png" | "webp" | "avif";

export interface RenditionSize {
  width: number;
  height: number;
  /** Crop to exactly these dimensions instead of fitting inside them. */
  crop: boolean;
}

export interface MediaProfile {
  name: string;
  max_width: number;
  max_height: number;
  crop: boolean;
  /** 0 to 1. Ignored by lossless formats. */
  quality: number;
  format: MediaFormat;
  /**
   * Refuse anything larger, read from the decoded dimensions.
   *
   * A decompression bomb on the client is the user's own device rather than
   * your server, so this is a crash guard rather than a security control: a
   * 144 megapixel image will hang or kill a mobile browser tab.
   */
  max_pixels: number;
  renditions?: Record<string, RenditionSize>;
}

/** One encoded output. */
export interface Rendition {
  name: string;
  blob: Blob;
  width: number;
  height: number;
  mime: string;
  ext: string;
}

export interface OptimizedImage {
  primary: Rendition;
  extra: Rendition[];
  /** The source dimensions, before resizing. */
  originalWidth: number;
  originalHeight: number;
  /** The size of the file the user picked. */
  originalSize: number;
}

/**
 * Turns a picked file into what should be stored.
 *
 * Implemented once per platform: canvas on the web, expo-image-manipulator on
 * React Native. Injected into the uploader rather than imported by it, so no
 * bundler has to resolve a platform at build time.
 */
export type Optimizer = (
  file: Blob,
  profile: MediaProfile,
) => Promise<OptimizedImage>;

/** What the API records once the bytes are in storage. */
export interface FileRef {
  url: string;
  key: string;
  name: string;
  mime: string;
  size: number;
  width?: number;
  height?: number;
  thumbnail_url?: string;
  format?: string;
  optimised?: boolean;
  renditions?: Record<string, { url: string; key: string; width?: number; height?: number; size?: number; mime?: string }>;
  profile?: string;
}

/** The minimum the uploader needs from your API client. */
export interface UploadTransport {
  /** GET, returning parsed JSON. */
  get<T>(path: string): Promise<T>;
  /** POST JSON, returning parsed JSON. */
  post<T>(path: string, body: unknown): Promise<T>;
  /**
   * PUT raw bytes to an absolute URL, with no auth headers.
   *
   * Deliberately separate from post: a presigned URL carries its own
   * authorisation in the query string, and attaching an Authorization header
   * to it makes S3 reject the request.
   */
  put(url: string, body: Blob, contentType: string, onProgress?: (fraction: number) => void): Promise<void>;
}
`
}

func uploadIndexTS() string {
	return `export * from "./types";
export * from "./uploader";
export * from "./transport";
`
}

func uploadWebTS() string {
	return `import type { MediaProfile, OptimizedImage, Rendition, RenditionSize } from "./types";

/**
 * Browser image optimisation, via canvas.
 *
 * This runs before the file leaves the device, which is the point: a presigned
 * upload goes straight to storage, so bytes not removed here are bytes paid
 * for on the user's connection. On a phone that is the difference between a
 * 6 MB upload and a 40 KB one.
 *
 * It also gets something the Go server cannot do without cgo: the browser has
 * a lossy WebP encoder built in, so the default output here is lossy WebP,
 * roughly four times smaller than the JPEG the pure-Go backend would produce.
 */

const MIME: Record<string, string> = {
  jpeg: "image/jpeg",
  png: "image/png",
  webp: "image/webp",
  avif: "image/avif",
};

const EXT: Record<string, string> = {
  jpeg: ".jpg",
  png: ".png",
  webp: ".webp",
  avif: ".avif",
};

let webpSupport: boolean | null = null;

/**
 * Whether this browser can actually encode WebP from a canvas.
 *
 * Worth testing rather than assuming, because toBlob does not fail on an
 * unsupported type: it silently returns a PNG. Asking for WebP and getting a
 * PNG back would make photographs larger than the JPEG they replaced, and
 * nothing would look broken.
 */
export async function canEncodeWebP(): Promise<boolean> {
  if (webpSupport !== null) return webpSupport;
  try {
    const canvas = document.createElement("canvas");
    canvas.width = 1;
    canvas.height = 1;
    const blob = await new Promise<Blob | null>((resolve) =>
      canvas.toBlob(resolve, "image/webp", 0.8),
    );
    webpSupport = blob?.type === "image/webp";
  } catch {
    webpSupport = false;
  }
  return webpSupport;
}

/** Decode, applying EXIF orientation, or a portrait photo comes out sideways. */
async function decode(file: Blob): Promise<ImageBitmap | HTMLImageElement> {
  if (typeof createImageBitmap === "function") {
    try {
      return await createImageBitmap(file, { imageOrientation: "from-image" });
    } catch {
      // Older Safari rejects the options object. Fall through.
    }
  }
  const url = URL.createObjectURL(file);
  try {
    const img = new Image();
    img.decoding = "sync";
    await new Promise<void>((resolve, reject) => {
      img.onload = () => resolve();
      img.onerror = () => reject(new Error("could not decode image"));
      img.src = url;
    });
    return img;
  } finally {
    // Revoked after load: the bitmap has already been decoded into memory.
    URL.revokeObjectURL(url);
  }
}

function sizeOf(src: ImageBitmap | HTMLImageElement) {
  const w = "naturalWidth" in src ? src.naturalWidth : src.width;
  const h = "naturalHeight" in src ? src.naturalHeight : src.height;
  return { w, h };
}

/** Target dimensions for a fit or a fill. Fit never scales up. */
function target(srcW: number, srcH: number, box: RenditionSize) {
  if (box.crop) return { w: box.width, h: box.height };
  if (srcW <= box.width && srcH <= box.height) return { w: srcW, h: srcH };
  const scale = Math.min(box.width / srcW, box.height / srcH);
  return { w: Math.round(srcW * scale), h: Math.round(srcH * scale) };
}

/**
 * Whether any pixel is meaningfully transparent.
 *
 * The threshold is not 255, because resampling leaves alpha a hair under
 * opaque; treating that as transparency would push every resized photograph
 * down the lossless path, where it is many times larger. A JPEG source skips
 * the scan entirely, since the format has no alpha channel.
 */
function hasAlpha(ctx: CanvasRenderingContext2D, w: number, h: number, sourceType: string): boolean {
  if (sourceType === "image/jpeg") return false;
  const { data } = ctx.getImageData(0, 0, w, h);
  for (let i = 3; i < data.length; i += 4) {
    if (data[i] < 250) return true;
  }
  return false;
}

async function encode(
  canvas: HTMLCanvasElement,
  format: string,
  quality: number,
): Promise<{ blob: Blob; format: string }> {
  const want = MIME[format] ?? "image/jpeg";
  const blob = await new Promise<Blob | null>((resolve) =>
    canvas.toBlob(resolve, want, quality),
  );
  if (!blob) throw new Error("canvas produced no image");
  // toBlob falls back to PNG rather than failing on an unsupported type, so
  // trust what came back rather than what was asked for.
  const got = blob.type || want;
  if (got !== want && got === "image/png" && want !== "image/png") {
    const jpeg = await new Promise<Blob | null>((resolve) =>
      canvas.toBlob(resolve, "image/jpeg", quality),
    );
    if (jpeg) return { blob: jpeg, format: "jpeg" };
  }
  return { blob, format: got.replace("image/", "") };
}

async function render(
  src: ImageBitmap | HTMLImageElement,
  box: RenditionSize,
  profile: MediaProfile,
  sourceType: string,
  name: string,
): Promise<Rendition> {
  const { w: srcW, h: srcH } = sizeOf(src);
  const { w, h } = target(srcW, srcH, box);

  const canvas = document.createElement("canvas");
  canvas.width = w;
  canvas.height = h;
  const ctx = canvas.getContext("2d", { alpha: true });
  if (!ctx) throw new Error("no 2d canvas context");
  ctx.imageSmoothingEnabled = true;
  ctx.imageSmoothingQuality = "high";

  if (box.crop) {
    // Cover, centred: scale so the box is filled, then take the middle. The
    // alternative letterboxes a portrait photo inside a square avatar.
    const scale = Math.max(w / srcW, h / srcH);
    const dw = srcW * scale;
    const dh = srcH * scale;
    ctx.drawImage(src as CanvasImageSource, (w - dw) / 2, (h - dh) / 2, dw, dh);
  } else {
    ctx.drawImage(src as CanvasImageSource, 0, 0, w, h);
  }

  let format = profile.format as string;
  if (format === "auto" || format === "avif") {
    // avif collapses to the same choice: no browser encodes it from a canvas.
    format = (await canEncodeWebP()) ? "webp" : hasAlpha(ctx, w, h, sourceType) ? "png" : "jpeg";
  }
  if (format === "jpeg" && hasAlpha(ctx, w, h, sourceType)) {
    // Never flatten transparency onto black, which is what JPEG would do.
    format = (await canEncodeWebP()) ? "webp" : "png";
  }

  const { blob, format: actual } = await encode(canvas, format, profile.quality);
  return {
    name,
    blob,
    width: w,
    height: h,
    mime: MIME[actual] ?? blob.type,
    ext: EXT[actual] ?? "",
  };
}

/** Optimise an image in the browser. */
export async function optimizeImage(
  file: Blob,
  profile: MediaProfile,
): Promise<OptimizedImage> {
  const src = await decode(file);
  const { w: srcW, h: srcH } = sizeOf(src);

  if (profile.max_pixels > 0 && srcW * srcH > profile.max_pixels) {
    throw new Error(
      "That image is " + srcW + "x" + srcH + ", larger than this app accepts (" +
        Math.round(profile.max_pixels / 1000000) + " megapixels).",
    );
  }

  const primary = await render(
    src,
    { width: profile.max_width, height: profile.max_height, crop: profile.crop },
    profile,
    file.type,
    "primary",
  );

  const extra: Rendition[] = [];
  for (const [name, box] of Object.entries(profile.renditions ?? {})) {
    // Rendered from the source, not from the primary, so a crop is taken at
    // full detail rather than from an already-downscaled copy.
    extra.push(await render(src, box, profile, file.type, name));
  }

  if ("close" in src && typeof src.close === "function") src.close();

  return {
    primary,
    extra,
    originalWidth: srcW,
    originalHeight: srcH,
    originalSize: file.size,
  };
}
`
}

func uploadExpoTS() string {
	return `import type { MediaProfile, OptimizedImage, Rendition, RenditionSize } from "./types";

/**
 * Expo / React Native image optimisation.
 *
 * There is no canvas here, so the work goes through expo-image-manipulator,
 * which resizes and re-encodes natively. Install it in the app that uses this:
 *
 *   npx expo install expo-image-manipulator
 *
 * Two differences from the web backend are worth knowing rather than
 * discovering. The manipulator only writes JPEG and PNG, so there is no WebP
 * here: a photograph becomes a JPEG. And it takes a URI rather than a Blob,
 * because on React Native a picked image is a file path, and reading it into
 * memory first is exactly what you are trying to avoid on a phone.
 */

type ManipulateResult = { uri: string; width: number; height: number };
type SaveFormat = "jpeg" | "png";

interface Manipulator {
  manipulateAsync(
    uri: string,
    actions: Array<Record<string, unknown>>,
    options: { compress: number; format: unknown; base64?: boolean },
  ): Promise<ManipulateResult>;
  SaveFormat: { JPEG: unknown; PNG: unknown };
}

/**
 * An image on this platform: the URI the picker returned, plus what is known
 * about it. There is no Blob until the upload actually reads the file.
 */
export interface NativeImage {
  uri: string;
  width: number;
  height: number;
  mimeType?: string;
  fileSize?: number;
}

async function loadManipulator(): Promise<Manipulator> {
  // The specifier is a variable on purpose. expo-image-manipulator is an
  // optional peer that only exists inside an Expo app, and a literal import
  // would make this file fail to type-check in the web app, which never uses
  // it. A non-literal specifier is resolved at runtime instead.
  const id = "expo-image-manipulator";
  try {
    return (await import(/* @vite-ignore */ id)) as unknown as Manipulator;
  } catch {
    throw new Error(
      "expo-image-manipulator is not installed. Run: npx expo install expo-image-manipulator",
    );
  }
}

function target(srcW: number, srcH: number, box: RenditionSize) {
  if (box.crop) return { w: box.width, h: box.height };
  if (srcW <= box.width && srcH <= box.height) return { w: srcW, h: srcH };
  const scale = Math.min(box.width / srcW, box.height / srcH);
  return { w: Math.round(srcW * scale), h: Math.round(srcH * scale) };
}

async function render(
  image: NativeImage,
  box: RenditionSize,
  profile: MediaProfile,
  name: string,
): Promise<Rendition> {
  const m = await loadManipulator();
  const { w, h } = target(image.width, image.height, box);

  const actions: Array<Record<string, unknown>> = [{ resize: { width: w, height: h } }];

  // PNG only when the source is a PNG, since this backend cannot write WebP
  // and re-encoding a transparent PNG as JPEG would flatten it onto black.
  const keepPNG = (image.mimeType ?? "").includes("png") || profile.format === "png";
  const format: SaveFormat = keepPNG ? "png" : "jpeg";

  const result = await m.manipulateAsync(image.uri, actions, {
    compress: profile.quality,
    format: format === "png" ? m.SaveFormat.PNG : m.SaveFormat.JPEG,
  });

  const response = await fetch(result.uri);
  const blob = await response.blob();

  return {
    name,
    blob,
    width: result.width,
    height: result.height,
    mime: format === "png" ? "image/png" : "image/jpeg",
    ext: format === "png" ? ".png" : ".jpg",
  };
}

/**
 * Optimise a picked image on Expo.
 *
 * Takes a NativeImage rather than a Blob, which is why it is not the same
 * signature as the web Optimizer. Wrap it where you use it:
 *
 *   const optimize = (_: Blob, p: MediaProfile) => optimizeNativeImage(picked, p);
 */
export async function optimizeNativeImage(
  image: NativeImage,
  profile: MediaProfile,
): Promise<OptimizedImage> {
  if (profile.max_pixels > 0 && image.width * image.height > profile.max_pixels) {
    throw new Error(
      "That image is " + image.width + "x" + image.height + ", larger than this app accepts.",
    );
  }

  const primary = await render(
    image,
    { width: profile.max_width, height: profile.max_height, crop: profile.crop },
    profile,
    "primary",
  );

  const extra: Rendition[] = [];
  for (const [name, box] of Object.entries(profile.renditions ?? {})) {
    extra.push(await render(image, box, profile, name));
  }

  return {
    primary,
    extra,
    originalWidth: image.width,
    originalHeight: image.height,
    originalSize: image.fileSize ?? 0,
  };
}
`
}

func uploadUploaderTS() string {
	return `import type {
  FileRef,
  MediaProfile,
  OptimizedImage,
  Optimizer,
  Rendition,
  UploadTransport,
} from "./types";

/**
 * The upload flow, with no platform in it.
 *
 * Optimise, presign, PUT straight to storage, then tell the API what landed.
 * The bytes never touch your server, so it costs no bandwidth and no CPU there
 * and does not care how large the file was to begin with.
 *
 * The optimizer is injected rather than imported, so this file works
 * unchanged on Next.js, a Vite SPA and Expo, and no bundler has to resolve a
 * platform at build time.
 */

export interface UploaderOptions {
  transport: UploadTransport;
  optimize: Optimizer;
  /** Overrides the profiles fetched from the API. Mostly for tests. */
  profiles?: MediaProfile[];
}

export interface UploadOptions {
  /** Profile name. Falls back to "default". */
  profile?: string;
  /** MIME aliases the field accepts, e.g. ["image"]. */
  accepts?: string[];
  /** 0 to 1, across every rendition. */
  onProgress?: (fraction: number) => void;
  /** Called once the sizes are known, before anything is uploaded. */
  onOptimized?: (result: OptimizedImage) => void;
}

const DEFAULT_PROFILE: MediaProfile = {
  name: "default",
  max_width: 1600,
  max_height: 1600,
  crop: false,
  quality: 0.82,
  format: "auto",
  max_pixels: 50000000,
  renditions: { thumb: { width: 400, height: 400, crop: true } },
};

export function createUploader({ transport, optimize, profiles }: UploaderOptions) {
  let cache: MediaProfile[] | null = profiles ?? null;

  /**
   * The profiles the server defines.
   *
   * Fetched once and cached. If the request fails the built-in default is
   * used, because an upload failing because a config endpoint was briefly
   * unreachable is a worse outcome than optimising to standard settings.
   */
  async function getProfiles(): Promise<MediaProfile[]> {
    if (cache) return cache;
    try {
      const res = await transport.get<{ data: { profiles: MediaProfile[] } }>(
        "/media/profiles",
      );
      cache = res.data.profiles?.length ? res.data.profiles : [DEFAULT_PROFILE];
    } catch {
      cache = [DEFAULT_PROFILE];
    }
    return cache;
  }

  async function profileFor(name?: string): Promise<MediaProfile> {
    const all = await getProfiles();
    return (
      all.find((p) => p.name === (name ?? "default")) ??
      all.find((p) => p.name === "default") ??
      DEFAULT_PROFILE
    );
  }

  /** Presign, PUT, and return the stored key. */
  async function putOne(
    rendition: Rendition,
    filename: string,
    accepts: string[] | undefined,
    onProgress?: (f: number) => void,
  ): Promise<{ key: string; url: string }> {
    const presign = await transport.post<{
      data: { presigned_url: string; key: string; public_url: string };
    }>("/uploads/presign", {
      filename,
      content_type: rendition.mime,
      // The exact byte count, which the server signs into the URL. Possible
      // only because the optimisation already happened: the size is known
      // before the URL is asked for.
      file_size: rendition.blob.size,
      accepts,
    });

    await transport.put(
      presign.data.presigned_url,
      rendition.blob,
      rendition.mime,
      onProgress,
    );

    return { key: presign.data.key, url: presign.data.public_url };
  }

  /**
   * Optimise a file and upload it straight to storage.
   *
   * Non-images skip the optimiser and are uploaded as they are, since there is
   * nothing useful to do to a PDF in a browser.
   */
  async function upload(
    file: Blob,
    filename: string,
    options: UploadOptions = {},
  ): Promise<FileRef> {
    const profile = await profileFor(options.profile);
    const isImage = file.type.startsWith("image/") && !file.type.includes("svg");

    let optimized: OptimizedImage | null = null;
    if (isImage && file.type !== "image/gif") {
      // GIF is skipped on purpose: drawing one to a canvas keeps the first
      // frame only, so "optimising" an animation would throw it away.
      try {
        optimized = await optimize(file, profile);
        options.onOptimized?.(optimized);
      } catch (err) {
        // A file this browser cannot decode is still a file the user chose.
        // Upload it untouched rather than losing it, and record that nothing
        // was done to it.
        optimized = null;
        if (String(err).includes("megapixels")) throw err;
      }
    }

    const parts: Rendition[] = optimized
      ? [optimized.primary, ...optimized.extra]
      : [
          {
            name: "primary",
            blob: file,
            width: 0,
            height: 0,
            mime: file.type || "application/octet-stream",
            ext: "",
          },
        ];

    const total = parts.reduce((sum, p) => sum + p.blob.size, 0);
    let done = 0;

    const stored: Record<string, { key: string; url: string; rendition: Rendition }> = {};
    for (const part of parts) {
      const name = part.name === "primary" ? filename : stem(filename) + "-" + part.name + part.ext;
      const { key, url } = await putOne(part, name, options.accepts, (f) => {
        options.onProgress?.((done + f * part.blob.size) / total);
      });
      done += part.blob.size;
      options.onProgress?.(done / total);
      stored[part.name] = { key, url, rendition: part };
    }

    const primary = stored.primary;
    const renditions: FileRef["renditions"] = {};
    for (const [name, s] of Object.entries(stored)) {
      if (name === "primary") continue;
      renditions[name] = {
        url: s.url,
        key: s.key,
        width: s.rendition.width,
        height: s.rendition.height,
        size: s.rendition.blob.size,
        mime: s.rendition.mime,
      };
    }

    // Records the row, and re-reads the object from storage to confirm what
    // actually arrived rather than trusting these numbers.
    const completed = await transport.post<{ data: FileRef }>("/uploads/complete", {
      key: primary.key,
      filename,
      content_type: primary.rendition.mime,
      size: primary.rendition.blob.size,
      accepts: options.accepts,
    });

    return {
      ...completed.data,
      width: primary.rendition.width || undefined,
      height: primary.rendition.height || undefined,
      format: primary.rendition.mime.replace("image/", ""),
      optimised: optimized !== null,
      profile: profile.name,
      thumbnail_url: renditions.thumb?.url,
      renditions: Object.keys(renditions).length ? renditions : undefined,
    };
  }

  return { upload, getProfiles, profileFor };
}

function stem(filename: string): string {
  const i = filename.lastIndexOf(".");
  return i > 0 ? filename.slice(0, i) : filename;
}
`
}
