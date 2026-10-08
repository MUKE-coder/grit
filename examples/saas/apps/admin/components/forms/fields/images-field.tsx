"use client";

import type { FieldDefinition } from "@/lib/resource";
import { BoxDropzone, type UploadedFile } from "@/components/ui/dropzone";

interface ImagesFieldProps {
  field: FieldDefinition;
  // Loose typing for the same reason FilesField does: guard against
  // non-array values arriving from a stale default or API response.
  value: string[] | unknown;
  onChange: (value: string[]) => void;
  error?: string;
}

export function ImagesField({ field, value, onChange, error }: ImagesFieldProps) {
  const urls = Array.isArray(value) ? (value as string[]) : [];
  const existingFiles: UploadedFile[] = urls.map((url, i) => ({
    url,
    name: `Image ${i + 1}`,
    size: 0,
    type: "image/jpeg",
  }));

  return (
    <BoxDropzone
      maxFiles={field.max ?? 10}
      maxSize={field.maxSize ?? 5 * 1024 * 1024}
      accept={{ "image/*": [".jpeg", ".jpg", ".png", ".gif", ".webp"] }}
      value={existingFiles}
      onFilesChange={(files) => {
        onChange(files.map((f) => f.url));
      }}
      label={field.label}
      description={field.description ?? "Upload up to " + String(field.max ?? 10) + " images"}
      error={error}
    />
  );
}
