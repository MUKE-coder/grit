import { z } from "zod";

// FileRef — canonical shape of a stored file. The API's POST /api/uploads
// returns this; resource forms embed it in their submit body.
//
// Fields width / height / duration / thumbnail_url are populated by the
// server when the source format makes them cheap to compute (images get
// dimensions, audio gets duration). They're optional because not every
// upload has them — a PDF has no width, a CSV has no thumbnail.
export const RenditionSchema = z.object({
  url: z.string().url(),
  key: z.string().min(1),
  width: z.number().int().positive().optional(),
  height: z.number().int().positive().optional(),
  size: z.number().int().nonnegative().optional(),
  mime: z.string().optional(),
});

export type Rendition = z.infer<typeof RenditionSchema>;

export const FileRefSchema = z.object({
  url: z.string().url(),
  key: z.string().min(1),
  name: z.string().min(1),
  mime: z.string().min(1),
  size: z.number().int().nonnegative(),
  width: z.number().int().positive().optional(),
  height: z.number().int().positive().optional(),
  duration: z.number().int().nonnegative().optional(),
  thumbnail_url: z.string().url().optional(),

  // Written by the image pipeline. format is what the file actually is after
  // optimisation, which is not always what was uploaded: a PNG photograph
  // comes back as a JPEG, and a PNG with transparency comes back as WebP.
  format: z.string().optional(),
  // false when the transform failed and the upload was stored as it arrived.
  optimised: z.boolean().optional(),
  // The untouched upload, kept privately so a profile change can be replayed.
  // A key rather than a URL, because the original is not for serving.
  original_key: z.string().optional(),
  original_size: z.number().int().nonnegative().optional(),
  renditions: z.record(z.string(), RenditionSchema).optional(),
  profile: z.string().optional(),
});

export type FileRef = z.infer<typeof FileRefSchema>;
