"use client";

// v3.31.30 — FileField. Single-file variant. Value is a FileRef (the
// JSON shape returned by POST /api/uploads), not a bare URL string.
// Storing the full ref means previews can render without re-fetching
// metadata and the parent record carries enough info for the storage
// admin page to compute usage totals.

import type { FieldDefinition } from "@/lib/resource";
import type { FileRef } from "@repo/shared/schemas";
import { Dropzone, type UploadedFile } from "@/components/ui/dropzone";
import { acceptsToReactDropzoneFormat, buildUploadEndpoint } from "@/lib/file-accepts";

interface FileFieldProps {
  field: FieldDefinition;
  value: FileRef | null;
  onChange: (value: FileRef | null) => void;
  error?: string;
}

function refToUploaded(ref: FileRef | null): UploadedFile[] {
  if (!ref) return [];
  return [{ url: ref.url, key: ref.key, name: ref.name, size: ref.size, type: ref.mime, thumbnail_url: ref.thumbnail_url }];
}

function uploadedToRef(u: UploadedFile): FileRef {
  return {
    url: u.url,
    key: u.key || extractKeyFromUrl(u.url),
    name: u.name,
    mime: u.type,
    size: u.size,
    thumbnail_url: u.thumbnail_url,
  };
}

// Fallback for the rare case where the Dropzone never round-tripped
// the file through /api/uploads (e.g. value loaded from server-side
// state without a key column). Pathname is good enough for the
// storage admin page to deduplicate but won't survive a CDN rewrite.
function extractKeyFromUrl(url: string): string {
  try {
    return new URL(url).pathname.replace(/^\//, "");
  } catch {
    return url;
  }
}

export function FileField({ field, value, onChange, error }: FileFieldProps) {
  const maxBytes = (field.maxSizeMB ?? 5) * 1024 * 1024;
  return (
    <Dropzone
      variant={field.dropzone ?? "default"}
      progress={field.progress ?? "bar"}
      maxFiles={1}
      maxSize={maxBytes}
      accept={acceptsToReactDropzoneFormat(field.accepts ?? ["all"])}
      uploadEndpoint={buildUploadEndpoint(field.accepts, maxBytes)}
      value={refToUploaded(value)}
      onFilesChange={(files) => {
        onChange(files[0] ? uploadedToRef(files[0]) : null);
      }}
      label={field.label}
      description={field.description}
      error={error}
    />
  );
}
