import type { MediaProfile, OptimizedImage, Rendition, RenditionSize } from "./types";

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
