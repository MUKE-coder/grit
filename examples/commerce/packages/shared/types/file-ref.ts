// FileRef — re-export of the Zod-inferred type for code that only
// needs the type, not the schema. The schema lives in
// schemas/file-ref.ts; importing the type from here avoids pulling in
// Zod just to get a type definition.


/** One derived size of a file, produced by the optimisation pipeline. */
export type Rendition = {
  url: string;
  key: string;
  width?: number;
  height?: number;
  size?: number;
  mime?: string;
};

export type FileRef = {
  url: string;
  key: string;
  name: string;
  mime: string;
  size: number;
  width?: number;
  height?: number;
  duration?: number;
  thumbnail_url?: string;

  /** What the stored file actually is, after optimisation. */
  format?: string;
  /** false when the transform failed and the upload was stored as it arrived. */
  optimised?: boolean;
  /** The untouched upload, kept privately for reprocessing. Not for serving. */
  original_key?: string;
  original_size?: number;
  /** Extra sizes the profile asked for, keyed by name ("thumb"). */
  renditions?: Record<string, Rendition>;
  /** The profile that produced this ref. */
  profile?: string;
};
