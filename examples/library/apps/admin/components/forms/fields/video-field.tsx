"use client";

import type { FieldDefinition } from "@/lib/resource";
import { CompactDropzone, type UploadedFile } from "@/components/ui/dropzone";

interface VideoFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export function VideoField({ field, value, onChange, error }: VideoFieldProps) {
  const existingFiles: UploadedFile[] = value
    ? [{ url: value, name: "Current video", size: 0, type: "video/mp4" }]
    : [];

  return (
    <CompactDropzone
      maxFiles={1}
      maxSize={field.maxSize ?? 100 * 1024 * 1024}
      accept={{ "video/*": [".mp4", ".webm", ".mov"] }}
      value={existingFiles}
      onFilesChange={(files) => {
        onChange(files[0]?.url || "");
      }}
      label={field.label}
      description={field.description ?? "MP4, WebM, or MOV up to 100MB"}
      error={error}
    />
  );
}
