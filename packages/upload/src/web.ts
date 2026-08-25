import type { MediaProfile, OptimizedImage, Rendition, RenditionSize } from "./types";

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
