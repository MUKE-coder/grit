/**
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
