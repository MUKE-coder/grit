"use client";

import type { FieldDefinition } from "@/lib/resource";
import { BoxDropzone, type UploadedFile } from "@/components/ui/dropzone";

interface VideosFieldProps {
  field: FieldDefinition;
  value: string[] | unknown;
  onChange: (value: string[]) => void;
  error?: string;
}

export function VideosField({ field, value, onChange, error }: VideosFieldProps) {
  const urls = Array.isArray(value) ? (value as string[]) : [];
  const existingFiles: UploadedFile[] = urls.map((url, i) => ({
    url,
    name: `Video ${i + 1}`,
    size: 0,
    type: "video/mp4",
  }));

  return (
    <BoxDropzone
      maxFiles={field.max ?? 5}
      maxSize={field.maxSize ?? 100 * 1024 * 1024}
      accept={{ "video/*": [".mp4", ".webm", ".mov"] }}
      value={existingFiles}
      onFilesChange={(files) => {
        onChange(files.map((f) => f.url));
      }}
      label={field.label}
      description={field.description ?? "Upload up to " + String(field.max ?? 5) + " videos"}
      error={error}
    />
  );
}
