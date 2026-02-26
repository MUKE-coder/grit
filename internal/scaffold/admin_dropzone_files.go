package scaffold

// adminDropzone returns the reusable Dropzone component with multiple variants.
func adminDropzone() string {
	return `"use client";

import { useCallback, useState } from "react";
import { useDropzone, type Accept } from "react-dropzone";
import { Upload, X, File, Image as ImageIcon, Loader2 } from "@/lib/icons";
import { uploadFile } from "@/lib/api-client";

// ── Types ────────────────────────────────────────────────────────

export type DropzoneVariant = "default" | "compact" | "minimal" | "avatar" | "inline";
export type DropzoneProvider = "cloudflare" | "aws" | "minio" | "local";

export interface UploadedFile {
  id?: number;
  url: string;
  name: string;
  size: number;
  type: string;
  thumbnail_url?: string;
}

export interface DropzoneProps {
  /** Storage provider hint (for display only) */
  provider?: DropzoneProvider;
  /** Visual variant */
  variant?: DropzoneVariant;
  /** Maximum number of files */
  maxFiles?: number;
  /** Maximum file size in bytes (default 10MB) */
  maxSize?: number;
  /** Accepted MIME types */
  accept?: Accept;
  /** Callback when files change */
  onFilesChange?: (files: UploadedFile[]) => void;
  /** Callback for raw File objects before upload */
  onDrop?: (files: File[]) => void;
  /** Whether to auto-upload to /api/uploads */
  autoUpload?: boolean;
  /** Custom upload endpoint */
  uploadEndpoint?: string;
  /** Whether the dropzone is disabled */
  disabled?: boolean;
  /** Existing files to display */
  value?: UploadedFile[];
  /** Label text */
  label?: string;
  /** Helper text */
  description?: string;
  /** Error message */
  error?: string;
  /** CSS class overrides */
  className?: string;
}

// ── Main Component ───────────────────────────────────────────────

export function Dropzone({
  variant = "default",
  maxFiles = 1,
  maxSize = 10 * 1024 * 1024,
  accept,
  onFilesChange,
  onDrop: onDropProp,
  autoUpload = true,
  uploadEndpoint = "/api/uploads",
  disabled = false,
  value = [],
  label,
  description,
  error,
  className = "",
}: DropzoneProps) {
  const [files, setFiles] = useState<UploadedFile[]>(value);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [uploadError, setUploadError] = useState<string | null>(null);

  const onDrop = useCallback(
    async (acceptedFiles: File[]) => {
      if (onDropProp) {
        onDropProp(acceptedFiles);
      }

      if (!autoUpload) {
        const newFiles: UploadedFile[] = acceptedFiles.map((f) => ({
          url: URL.createObjectURL(f),
          name: f.name,
          size: f.size,
          type: f.type,
        }));
        const updated = maxFiles === 1 ? newFiles : [...files, ...newFiles].slice(0, maxFiles);
        setFiles(updated);
        onFilesChange?.(updated);
        return;
      }

      setUploading(true);
      setUploadError(null);
      setUploadProgress(0);

      const uploaded: UploadedFile[] = [];

      for (let i = 0; i < acceptedFiles.length; i++) {
        const file = acceptedFiles[i];

        try {
          const result = await uploadFile(file, uploadEndpoint, (percent) => {
            const fileProgress = ((i + percent / 100) / acceptedFiles.length) * 100;
            setUploadProgress(Math.round(fileProgress));
          });
          const d = result.data as Record<string, unknown>;
          uploaded.push({
            id: d.id as number,
            url: (d.url as string) || "",
            name: (d.original_name as string) || file.name,
            size: (d.size as number) || file.size,
            type: (d.mime_type as string) || file.type,
            thumbnail_url: d.thumbnail_url as string,
          });
        } catch (err: unknown) {
          const axiosErr = err as { response?: { data?: { error?: { message?: string } } } };
          const msg = axiosErr?.response?.data?.error?.message || (err as Error)?.message || ` + "`" + `Failed to upload ${file.name}` + "`" + `;
          setUploadError(msg);
        }
      }

      const updated = maxFiles === 1 ? uploaded : [...files, ...uploaded].slice(0, maxFiles);
      setFiles(updated);
      onFilesChange?.(updated);
      setUploading(false);
      setUploadProgress(0);
    },
    [files, maxFiles, autoUpload, uploadEndpoint, onFilesChange, onDropProp]
  );

  const removeFile = (index: number) => {
    const updated = files.filter((_, i) => i !== index);
    setFiles(updated);
    onFilesChange?.(updated);
  };

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept,
    maxFiles,
    maxSize,
    disabled: disabled || uploading,
    onDropRejected: (rejections) => {
      const msg = rejections[0]?.errors[0]?.message || "File rejected";
      setUploadError(msg);
    },
  });

  // ── Render by variant ─────────────────────────────────────────

  const VariantComponent = {
    default: DefaultVariant,
    compact: CompactVariant,
    minimal: MinimalVariant,
    avatar: AvatarVariant,
    inline: InlineVariant,
  }[variant];

  return (
    <div className={` + "`" + `space-y-1.5 ${className}` + "`" + `}>
      {label && (
        <label className="block text-sm font-medium text-foreground">{label}</label>
      )}

      <VariantComponent
        getRootProps={getRootProps}
        getInputProps={getInputProps}
        isDragActive={isDragActive}
        uploading={uploading}
        uploadProgress={uploadProgress}
        files={files}
        removeFile={removeFile}
        maxFiles={maxFiles}
        maxSize={maxSize}
        disabled={disabled}
      />

      {description && !error && !uploadError && (
        <p className="text-xs text-text-muted">{description}</p>
      )}
      {(error || uploadError) && (
        <p className="text-xs text-danger">{error || uploadError}</p>
      )}
    </div>
  );
}

// ── Shared Types ─────────────────────────────────────────────────

interface VariantProps {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  getRootProps: () => any;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  getInputProps: () => any;
  isDragActive: boolean;
  uploading: boolean;
  uploadProgress: number;
  files: UploadedFile[];
  removeFile: (index: number) => void;
  maxFiles: number;
  maxSize: number;
  disabled: boolean;
}

function formatSize(bytes: number) {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + " " + sizes[i];
}

function isImage(type: string) {
  return type.startsWith("image/");
}

// ── Variant: Default ─────────────────────────────────────────────

function DefaultVariant({
  getRootProps,
  getInputProps,
  isDragActive,
  uploading,
  uploadProgress,
  files,
  removeFile,
  maxSize,
}: VariantProps) {
  return (
    <div className="space-y-3">
      <div
        {...getRootProps()}
        className={` + "`" + `flex flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed p-8 cursor-pointer transition-all ${
          isDragActive
            ? "border-accent bg-accent/5 scale-[1.01]"
            : "border-border hover:border-accent/50 hover:bg-bg-hover/30"
        } ${uploading ? "opacity-60 cursor-not-allowed" : ""}` + "`" + `}
      >
        <input {...getInputProps()} />
        {uploading ? (
          <div className="flex flex-col items-center gap-2">
            <Loader2 className="h-8 w-8 animate-spin text-accent" />
            <p className="text-sm text-text-secondary">Uploading... {uploadProgress}%</p>
            <div className="w-48 h-1.5 bg-bg-tertiary rounded-full overflow-hidden">
              <div
                className="h-full bg-accent rounded-full transition-all duration-300"
                style={{ width: ` + "`" + `${uploadProgress}%` + "`" + ` }}
              />
            </div>
          </div>
        ) : (
          <>
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-bg-tertiary">
              {isDragActive ? (
                <Upload className="h-6 w-6 text-accent" />
              ) : (
                <Upload className="h-6 w-6 text-text-muted" />
              )}
            </div>
            <div className="text-center">
              <p className="text-sm font-medium text-foreground">
                {isDragActive ? "Drop files here" : "Click to upload or drag and drop"}
              </p>
              <p className="text-xs text-text-muted mt-1">
                Max size: {formatSize(maxSize)}
              </p>
            </div>
          </>
        )}
      </div>

      {files.length > 0 && (
        <div className="space-y-2">
          {files.map((file, i) => (
            <FilePreview key={i} file={file} onRemove={() => removeFile(i)} />
          ))}
        </div>
      )}
    </div>
  );
}

// ── Variant: Compact ─────────────────────────────────────────────

function CompactVariant({
  getRootProps,
  getInputProps,
  isDragActive,
  uploading,
  uploadProgress,
  files,
  removeFile,
  maxSize,
}: VariantProps) {
  return (
    <div className="space-y-2">
      <div
        {...getRootProps()}
        className={` + "`" + `flex items-center gap-3 rounded-lg border-2 border-dashed px-4 py-3 cursor-pointer transition-all ${
          isDragActive
            ? "border-accent bg-accent/5"
            : "border-border hover:border-accent/50 hover:bg-bg-hover/30"
        } ${uploading ? "opacity-60 cursor-not-allowed" : ""}` + "`" + `}
      >
        <input {...getInputProps()} />
        {uploading ? (
          <Loader2 className="h-5 w-5 animate-spin text-accent shrink-0" />
        ) : (
          <Upload className="h-5 w-5 text-text-muted shrink-0" />
        )}
        <div className="flex-1 min-w-0">
          {uploading ? (
            <div className="flex items-center gap-2">
              <div className="flex-1 h-1.5 bg-bg-tertiary rounded-full overflow-hidden">
                <div
                  className="h-full bg-accent rounded-full transition-all"
                  style={{ width: ` + "`" + `${uploadProgress}%` + "`" + ` }}
                />
              </div>
              <span className="text-xs text-text-muted">{uploadProgress}%</span>
            </div>
          ) : (
            <p className="text-sm text-text-secondary">
              {isDragActive ? "Drop here..." : ` + "`" + `Drop files or click to browse (max ${formatSize(maxSize)})` + "`" + `}
            </p>
          )}
        </div>
      </div>

      {files.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {files.map((file, i) => (
            <div key={i} className="flex items-center gap-1.5 rounded-md bg-bg-tertiary px-2.5 py-1.5 text-xs">
              <File className="h-3 w-3 text-text-muted" />
              <span className="text-foreground truncate max-w-[150px]">{file.name}</span>
              <button
                type="button"
                onClick={(e) => { e.stopPropagation(); removeFile(i); }}
                className="ml-1 rounded p-0.5 text-text-muted hover:text-danger hover:bg-danger/10 transition-colors"
              >
                <X className="h-3 w-3" />
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

// ── Variant: Minimal ─────────────────────────────────────────────

function MinimalVariant({
  getRootProps,
  getInputProps,
  isDragActive,
  uploading,
  files,
  removeFile,
}: VariantProps) {
  return (
    <div className="space-y-2">
      <div
        {...getRootProps()}
        className={` + "`" + `inline-flex items-center gap-2 rounded-lg border px-3 py-2 cursor-pointer transition-all ${
          isDragActive
            ? "border-accent bg-accent/5 text-accent"
            : "border-border hover:border-accent/50 text-text-secondary hover:text-foreground"
        } ${uploading ? "opacity-60 cursor-not-allowed" : ""}` + "`" + `}
      >
        <input {...getInputProps()} />
        {uploading ? (
          <Loader2 className="h-4 w-4 animate-spin" />
        ) : (
          <Upload className="h-4 w-4" />
        )}
        <span className="text-sm font-medium">
          {uploading ? "Uploading..." : isDragActive ? "Drop here" : "Upload file"}
        </span>
      </div>

      {files.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {files.map((file, i) => (
            <div key={i} className="flex items-center gap-1.5 text-xs text-text-secondary">
              <File className="h-3 w-3" />
              <span className="truncate max-w-[150px]">{file.name}</span>
              <button
                type="button"
                onClick={() => removeFile(i)}
                className="rounded p-0.5 hover:text-danger transition-colors"
              >
                <X className="h-3 w-3" />
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

// ── Variant: Avatar ──────────────────────────────────────────────

function AvatarVariant({
  getRootProps,
  getInputProps,
  isDragActive,
  uploading,
  files,
  removeFile,
}: VariantProps) {
  const preview = files[0];

  return (
    <div className="inline-block">
      {preview ? (
        <div className="relative group">
          {isImage(preview.type) ? (
            <img
              src={preview.url}
              alt={preview.name}
              className="h-24 w-24 rounded-full object-cover border-2 border-border"
            />
          ) : (
            <div className="flex h-24 w-24 items-center justify-center rounded-full bg-bg-tertiary border-2 border-border">
              <File className="h-8 w-8 text-text-muted" />
            </div>
          )}
          <button
            type="button"
            onClick={() => removeFile(0)}
            className="absolute -top-1 -right-1 rounded-full bg-danger p-1 text-white opacity-0 group-hover:opacity-100 transition-opacity shadow-lg"
          >
            <X className="h-3 w-3" />
          </button>
          {uploading && (
            <div className="absolute inset-0 flex items-center justify-center rounded-full bg-black/50">
              <Loader2 className="h-6 w-6 animate-spin text-white" />
            </div>
          )}
          <div
            {...getRootProps()}
            className="absolute inset-0 rounded-full cursor-pointer"
          >
            <input {...getInputProps()} />
          </div>
        </div>
      ) : (
        <div
          {...getRootProps()}
          className={` + "`" + `flex h-24 w-24 flex-col items-center justify-center gap-1 rounded-full border-2 border-dashed cursor-pointer transition-all ${
            isDragActive
              ? "border-accent bg-accent/5"
              : "border-border hover:border-accent/50 hover:bg-bg-hover/30"
          } ${uploading ? "opacity-60 cursor-not-allowed" : ""}` + "`" + `}
        >
          <input {...getInputProps()} />
          {uploading ? (
            <Loader2 className="h-6 w-6 animate-spin text-accent" />
          ) : (
            <>
              <ImageIcon className="h-6 w-6 text-text-muted" />
              <span className="text-[10px] text-text-muted">Upload</span>
            </>
          )}
        </div>
      )}
    </div>
  );
}

// ── Variant: Inline ──────────────────────────────────────────────

function InlineVariant({
  getRootProps,
  getInputProps,
  isDragActive,
  uploading,
  files,
  removeFile,
  maxSize,
}: VariantProps) {
  return (
    <div className="space-y-3">
      <div
        {...getRootProps()}
        className={` + "`" + `flex items-center justify-between rounded-lg border px-4 py-3 cursor-pointer transition-all ${
          isDragActive
            ? "border-accent bg-accent/5"
            : "border-border hover:border-accent/50"
        } ${uploading ? "opacity-60 cursor-not-allowed" : ""}` + "`" + `}
      >
        <input {...getInputProps()} />
        <div className="flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-bg-tertiary">
            {uploading ? (
              <Loader2 className="h-4 w-4 animate-spin text-accent" />
            ) : (
              <Upload className="h-4 w-4 text-text-muted" />
            )}
          </div>
          <div>
            <p className="text-sm font-medium text-foreground">
              {uploading ? "Uploading..." : isDragActive ? "Drop files here" : "Choose files"}
            </p>
            <p className="text-xs text-text-muted">
              Max {formatSize(maxSize)} per file
            </p>
          </div>
        </div>
        <span className="rounded-lg bg-accent px-3 py-1.5 text-xs font-medium text-white">
          Browse
        </span>
      </div>

      {files.length > 0 && (
        <div className="space-y-2">
          {files.map((file, i) => (
            <FilePreview key={i} file={file} onRemove={() => removeFile(i)} />
          ))}
        </div>
      )}
    </div>
  );
}

// ── File Preview Card ────────────────────────────────────────────

function FilePreview({ file, onRemove }: { file: UploadedFile; onRemove: () => void }) {
  return (
    <div className="flex items-center gap-3 rounded-lg border border-border bg-bg-secondary px-3 py-2.5">
      {isImage(file.type) ? (
        <img
          src={file.thumbnail_url || file.url}
          alt={file.name}
          className="h-10 w-10 rounded-md object-cover"
        />
      ) : (
        <div className="flex h-10 w-10 items-center justify-center rounded-md bg-bg-tertiary">
          <File className="h-5 w-5 text-text-muted" />
        </div>
      )}
      <div className="flex-1 min-w-0">
        <p className="text-sm font-medium text-foreground truncate">{file.name}</p>
        <p className="text-xs text-text-muted">{formatSize(file.size)}</p>
      </div>
      <button
        type="button"
        onClick={onRemove}
        className="rounded-lg p-1.5 text-text-muted hover:text-danger hover:bg-danger/10 transition-colors"
      >
        <X className="h-4 w-4" />
      </button>
    </div>
  );
}
`
}
