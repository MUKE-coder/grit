"use client";

import type { FieldDefinition } from "@/lib/resource";
import { AvatarDropzone, type UploadedFile } from "@/components/ui/dropzone";

interface ImageFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export function ImageField({ field, value, onChange, error }: ImageFieldProps) {
  const existingFiles: UploadedFile[] = value
    ? [{ url: value, name: "Current image", size: 0, type: "image/jpeg" }]
    : [];

  return (
    <AvatarDropzone
      maxSize={field.maxSize ?? 5 * 1024 * 1024}
      accept={{ "image/*": [".jpeg", ".jpg", ".png", ".gif", ".webp"] }}
      value={existingFiles}
      onFilesChange={(files) => {
        onChange(files[0]?.url || "");
      }}
      label={field.label}
      description={field.description}
      error={error}
    />
  );
}
