"use client";

// v3.31.30 — FilesField. Multi-file variant. Value is FileRef[].

import type { FieldDefinition } from "@/lib/resource";
import type { FileRef } from "@repo/shared/schemas";
import { Dropzone, type UploadedFile } from "@/components/ui/dropzone";
import { acceptsToReactDropzoneFormat, buildUploadEndpoint } from "@/lib/file-accepts";

interface FilesFieldProps {
  field: FieldDefinition;
  // Loosely typed so refsToUploaded can guard against non-array values
  // arriving from react-hook-form defaults or a stale API response.
  value: FileRef[] | unknown;
  onChange: (value: FileRef[]) => void;
  error?: string;
}

function refsToUploaded(refs: FileRef[] | unknown): UploadedFile[] {
  // Defensive: when react-hook-form's initial value falls back to ""
  // (the legacy buildDefaults behaviour pre-fix) or the API returns a
  // non-array, calling .map() throws. Bail to [] so the dropzone
  // mounts cleanly and the user can still upload files.
  if (!Array.isArray(refs)) return [];
  return refs.map((r) => ({
    url: r.url,
    key: r.key,
    name: r.name,
    size: r.size,
    type: r.mime,
    thumbnail_url: r.thumbnail_url,
  }));
}

function uploadedToRef(u: UploadedFile): FileRef {
  return {
    url: u.url,
    key: extractKeyFromUrl(u.url),
    name: u.name,
    mime: u.type,
    size: u.size,
    thumbnail_url: u.thumbnail_url,
  };
}

function extractKeyFromUrl(url: string): string {
  try {
    return new URL(url).pathname.replace(/^\//, "");
  } catch {
    return url;
  }
}

export function FilesField({ field, value, onChange, error }: FilesFieldProps) {
  const maxFiles = field.max ?? 5;
  const maxBytes = (field.maxSizeMB ?? 5) * 1024 * 1024;
  return (
    <Dropzone
      variant={field.dropzone ?? "default"}
      progress={field.progress ?? "bar"}
      reorderable={field.reorderable ?? true}
      maxFiles={maxFiles}
      maxSize={maxBytes}
      accept={acceptsToReactDropzoneFormat(field.accepts ?? ["all"])}
      uploadEndpoint={buildUploadEndpoint(field.accepts, maxBytes)}
      value={refsToUploaded(value)}
      onFilesChange={(files) => {
        onChange(files.map(uploadedToRef));
      }}
      label={field.label}
      description={field.description}
      error={error}
    />
  );
}
