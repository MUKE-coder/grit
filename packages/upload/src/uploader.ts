import type {
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
      // The stored name carries the format actually produced, so a WebP is not
      // filed as photo.jpg. The Content-Type is right either way and browsers
      // go by that, but a key whose extension contradicts its bytes confuses
      // every tool that reads the key instead: CDN rules, lifecycle policies,
      // and whoever is looking through the bucket six months from now.
      //
      // ext is empty when nothing was optimised, and then the original name is
      // exactly right.
      const name =
        part.name === "primary"
          ? part.ext
            ? stem(filename) + part.ext
            : filename
          : stem(filename) + "-" + part.name + part.ext;
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
