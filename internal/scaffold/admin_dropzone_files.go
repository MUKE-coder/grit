package scaffold

// adminDropzone returns the admin's upload component: one upload engine
// (Dropzone.Root and its parts) and five named looks built from it, plus
// the variant-driven Dropzone the file fields and older code use.
func adminDropzone() string {
	return `"use client";

import { createContext, useCallback, useContext, useEffect, useRef, useState, type ReactNode } from "react";
import { useDropzone, type Accept, type DropzoneState as DropTargetState } from "react-dropzone";
import {
  Upload,
  X,
  File,
  Image as ImageIcon,
  Loader2,
  ChevronUp,
  ChevronDown,
  Play,
  FileText,
  FileSpreadsheet,
  Music,
} from "@/lib/icons";
import { uploadFile } from "@/lib/api-client";

// One upload engine and five looks. Pick the look by name:
//
//   <AvatarDropzone />    a round picture, for a profile photo
//   <InlineDropzone />    a row with a Browse button, and file cards below
//   <CompactDropzone />   a one-line strip, and file chips below
//   <MinimalDropzone />   an "Upload file" button
//   <BoxDropzone />       the large dashed box, and file cards below
//
// or build your own from the parts every one of them is made of:
//
//   <Dropzone.Root maxFiles={3} onFilesChange={save}>
//     <Dropzone.Target className="...">Drop receipts here</Dropzone.Target>
//     <Dropzone.Progress />
//     <Dropzone.FileList reorderable />
//   </Dropzone.Root>
//
// <Dropzone variant="..."> still works, for a look chosen at runtime (the
// file fields read it from the resource definition) and for code written
// before the named components existed.

// ── Types ────────────────────────────────────────────────────────

export type DropzoneVariant = "default" | "compact" | "minimal" | "avatar" | "inline";
export type ProgressVariant = "bar" | "circular" | "pulse";
export type DropzoneProvider = "cloudflare" | "aws" | "minio" | "local";

export interface UploadedFile {
  id?: number;
  url: string;
  /** S3 key, needed by FileField/FilesField to round-trip a FileRef
   * without re-deriving the key from the URL. v3.31.31. */
  key?: string;
  name: string;
  size: number;
  /** MIME type. Kept named "type" to avoid churning every existing
   * call site; the FileField bridge maps it to FileRef.mime. */
  type: string;
  thumbnail_url?: string;
  /** The size of the file the user picked, before optimisation. */
  originalSize?: number;
  /** What the stored file actually is: "jpeg", "webp". */
  format?: string;
  /** Whether the pipeline transformed it, as opposed to storing it as-is. */
  optimised?: boolean;
}

/** What files are allowed, and what happens to them. Shared by every look. */
export interface DropzoneOptions {
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
}

/** The frame around every look: a label above, a hint or an error below. */
export interface DropzoneFrameProps {
  /** Label text */
  label?: string;
  /** Helper text */
  description?: string;
  /** Error message */
  error?: string;
  /** CSS class overrides */
  className?: string;
}

export type DropzoneBaseProps = DropzoneOptions & DropzoneFrameProps;

export interface BoxDropzoneProps extends DropzoneBaseProps {
  /** Progress indicator shown while uploading. Default "bar". */
  progress?: ProgressVariant;
  /** Up/down buttons on each file card (multi-file only). */
  reorderable?: boolean;
}

export interface CompactDropzoneProps extends DropzoneBaseProps {
  /** Progress indicator shown while uploading. Default "bar". */
  progress?: ProgressVariant;
}

/** An avatar holds one picture, so it takes no maxFiles. */
export type AvatarDropzoneProps = Omit<DropzoneBaseProps, "maxFiles">;

export interface DropzoneProps extends DropzoneBaseProps {
  /** Storage provider hint. Nothing reads it; kept so existing call sites compile. */
  provider?: DropzoneProvider;
  /** Visual variant. Prefer the named component when the look is fixed. */
  variant?: DropzoneVariant;
  /** Progress indicator variant (v3.31.31). Default "bar". Box and compact only. */
  progress?: ProgressVariant;
  /** Allow up/down reordering of files in the preview list (box, multi only). */
  reorderable?: boolean;
}

// ── State ────────────────────────────────────────────────────────

export interface DropzoneState
  extends Pick<DropTargetState, "getRootProps" | "getInputProps" | "isDragActive"> {
  files: UploadedFile[];
  uploading: boolean;
  /** 0 to 100, across every file in the current drop. */
  uploadProgress: number;
  uploadError: string | null;
  maxFiles: number;
  maxSize: number;
  removeFile: (index: number) => void;
  moveFile: (index: number, direction: -1 | 1) => void;
}

function useDropzoneState({
  maxFiles = 1,
  maxSize = 10 * 1024 * 1024,
  accept,
  onFilesChange,
  onDrop: onDropProp,
  autoUpload = true,
  uploadEndpoint = "/api/uploads",
  disabled = false,
  value = [],
}: DropzoneOptions): DropzoneState {
  const [files, setFiles] = useState<UploadedFile[]>(value);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [uploadError, setUploadError] = useState<string | null>(null);

  // Blob URLs this dropzone made for local previews. Each one pins its File in
  // memory until it is revoked, so they are released when the file leaves the
  // list (removed, or pushed out by a newer pick) and when the dropzone
  // unmounts. URLs that came from the server are never in here.
  const previewUrls = useRef<Set<string>>(new Set());
  const releasePreviews = useCallback((kept: UploadedFile[], previous: UploadedFile[]) => {
    const still = new Set(kept.map((f) => f.url));
    for (const f of previous) {
      if (!still.has(f.url) && previewUrls.current.delete(f.url)) {
        URL.revokeObjectURL(f.url);
      }
    }
  }, []);
  useEffect(() => {
    const urls = previewUrls.current;
    return () => {
      for (const url of urls) URL.revokeObjectURL(url);
      urls.clear();
    };
  }, []);

  const onDrop = useCallback(
    async (acceptedFiles: File[]) => {
      if (onDropProp) {
        onDropProp(acceptedFiles);
      }

      if (!autoUpload) {
        const newFiles: UploadedFile[] = acceptedFiles.map((f) => {
          const url = URL.createObjectURL(f);
          previewUrls.current.add(url);
          return { url, name: f.name, size: f.size, type: f.type };
        });
        const updated = maxFiles === 1 ? newFiles : [...files, ...newFiles].slice(0, maxFiles);
        releasePreviews(updated, [...files, ...newFiles]);
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
          // v3.31.30: /api/uploads now returns a FileRef shape with
          // name + mime (not original_name + mime_type). Fall back to
          // the legacy keys if an older API is in front of us so the
          // dropzone stays compatible across upgrade boundaries.
          const d = result.data as Record<string, unknown>;
          uploaded.push({
            id: d.id as number,
            url: (d.url as string) || "",
            key: (d.key as string) || (d.path as string) || "",
            name: (d.name as string) || (d.original_name as string) || file.name,
            size: (d.size as number) || file.size,
            type: (d.mime as string) || (d.mime_type as string) || file.type,
            thumbnail_url: (d.thumbnail_url as string) || undefined,
            // What the optimisation pipeline did. originalSize is the file the
            // user picked; size is what was actually stored. Showing both is
            // the difference between "uploaded" and "uploaded 5.3 MB as
            // 147 KB", and it is the only place that work is visible.
            originalSize: (d.original_size as number) || file.size,
            format: (d.format as string) || undefined,
            optimised: d.optimised === true,
          });
        } catch (err: unknown) {
          const axiosErr = err as { response?: { data?: { error?: { message?: string } } } };
          const msg = axiosErr?.response?.data?.error?.message || (err as Error)?.message || ` + "`" + `Failed to upload ${file.name}` + "`" + `;
          setUploadError(msg);
        }
      }

      const updated = maxFiles === 1 ? uploaded : [...files, ...uploaded].slice(0, maxFiles);
      releasePreviews(updated, files);
      setFiles(updated);
      onFilesChange?.(updated);
      setUploading(false);
      setUploadProgress(0);
    },
    [files, maxFiles, autoUpload, uploadEndpoint, onFilesChange, onDropProp, releasePreviews]
  );

  const removeFile = (index: number) => {
    const updated = files.filter((_, i) => i !== index);
    releasePreviews(updated, files);
    setFiles(updated);
    onFilesChange?.(updated);
  };

  // v3.31.31: simple up/down reorder. When reorderable && multi, the
  // FilePreview rows show small arrow buttons that swap adjacent files.
  // Drag-reorder needs dnd-kit and lands later; this covers 80% of the
  // ergonomics with zero new deps.
  const moveFile = (index: number, direction: -1 | 1) => {
    const target = index + direction;
    if (target < 0 || target >= files.length) return;
    const updated = [...files];
    [updated[index], updated[target]] = [updated[target], updated[index]];
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

  return {
    getRootProps,
    getInputProps,
    isDragActive,
    files,
    uploading,
    uploadProgress,
    uploadError,
    maxFiles,
    maxSize,
    removeFile,
    moveFile,
  };
}

const DropzoneContext = createContext<DropzoneState | null>(null);

/** The state of the enclosing Dropzone.Root, for a part you write yourself. */
export function useDropzoneContext(): DropzoneState {
  const state = useContext(DropzoneContext);
  if (!state) throw new Error("Dropzone parts must be rendered inside <Dropzone.Root>");
  return state;
}

// ── Parts ────────────────────────────────────────────────────────

/** Holds the files and the upload, and draws the label, hint and error. */
function DropzoneRoot({
  label,
  description,
  error,
  className = "",
  children,
  ...options
}: DropzoneBaseProps & { children: ReactNode }) {
  const state = useDropzoneState(options);
  return (
    <DropzoneContext.Provider value={state}>
      <div className={` + "`" + `space-y-1.5 ${className}` + "`" + `}>
        {label && (
          <label className="block text-sm font-medium text-foreground">{label}</label>
        )}

        {children}

        {description && !error && !state.uploadError && (
          <p className="text-xs text-text-muted">{description}</p>
        )}
        {(error || state.uploadError) && (
          <p className="text-xs text-danger">{error || state.uploadError}</p>
        )}
      </div>
    </DropzoneContext.Provider>
  );
}

/** The element that opens the file picker and takes a drop. */
function DropzoneTarget({ className, children }: { className?: string; children?: ReactNode }) {
  const { getRootProps, getInputProps } = useDropzoneContext();
  return (
    <div {...getRootProps()} className={className}>
      <input {...getInputProps()} />
      {children}
    </div>
  );
}

/** The files so far, one card each. Renders nothing while the list is empty. */
function DropzoneFileList({ reorderable = false }: { reorderable?: boolean }) {
  const { files, removeFile, moveFile, maxFiles } = useDropzoneContext();
  if (files.length === 0) return null;
  return (
    <div className="space-y-2">
      {files.map((file, i) => (
        <FilePreview
          key={i}
          file={file}
          onRemove={() => removeFile(i)}
          index={i}
          total={files.length}
          reorderable={reorderable && maxFiles > 1}
          onMove={moveFile}
        />
      ))}
    </div>
  );
}

/** How far the upload has got. Renders nothing while nothing is uploading. */
function DropzoneProgress({ variant = "bar" }: { variant?: ProgressVariant }) {
  const { uploading, uploadProgress } = useDropzoneContext();
  if (!uploading) return null;
  if (variant === "circular") return <UploadProgressCircular percent={uploadProgress} />;
  if (variant === "pulse") return <UploadProgressPulse percent={uploadProgress} />;
  return <UploadProgressBar percent={uploadProgress} />;
}

// ── The five looks ───────────────────────────────────────────────

function targetState(isDragActive: boolean, uploading: boolean, idle: string, active: string) {
  return ` + "`" + `${isDragActive ? active : idle} ${uploading ? "opacity-60 cursor-not-allowed" : ""}` + "`" + `;
}

function BoxBody({ progress, reorderable }: { progress: ProgressVariant; reorderable: boolean }) {
  const { isDragActive, uploading, maxSize } = useDropzoneContext();
  return (
    <div className="space-y-3">
      <DropzoneTarget
        className={` + "`" + `flex flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed p-8 cursor-pointer transition-all ${targetState(
          isDragActive,
          uploading,
          "border-border hover:border-accent/50 hover:bg-bg-hover/30",
          "border-accent bg-accent/5 scale-[1.01]"
        )}` + "`" + `}
      >
        {uploading ? (
          <DropzoneProgress variant={progress} />
        ) : (
          <>
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-bg-tertiary">
              <Upload className={` + "`" + `h-6 w-6 ${isDragActive ? "text-accent" : "text-text-muted"}` + "`" + `} />
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
      </DropzoneTarget>

      <DropzoneFileList reorderable={reorderable} />
    </div>
  );
}

function CompactBody({ progress }: { progress: ProgressVariant }) {
  const { isDragActive, uploading, maxSize, files, removeFile } = useDropzoneContext();
  return (
    <div className="space-y-2">
      <DropzoneTarget
        className={` + "`" + `flex items-center gap-3 rounded-lg border-2 border-dashed px-4 py-3 cursor-pointer transition-all ${targetState(
          isDragActive,
          uploading,
          "border-border hover:border-accent/50 hover:bg-bg-hover/30",
          "border-accent bg-accent/5"
        )}` + "`" + `}
      >
        {!uploading && <Upload className="h-5 w-5 text-text-muted shrink-0" />}
        <div className="flex-1 min-w-0">
          {uploading ? (
            <DropzoneProgress variant={progress} />
          ) : (
            <p className="text-sm text-text-secondary">
              {isDragActive ? "Drop here..." : ` + "`" + `Drop files or click to browse (max ${formatSize(maxSize)})` + "`" + `}
            </p>
          )}
        </div>
      </DropzoneTarget>

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

function MinimalBody() {
  const { isDragActive, uploading, files, removeFile } = useDropzoneContext();
  return (
    <div className="space-y-2">
      <DropzoneTarget
        className={` + "`" + `inline-flex items-center gap-2 rounded-lg border px-3 py-2 cursor-pointer transition-all ${targetState(
          isDragActive,
          uploading,
          "border-border hover:border-accent/50 text-text-secondary hover:text-foreground",
          "border-accent bg-accent/5 text-accent"
        )}` + "`" + `}
      >
        {uploading ? (
          <Loader2 className="h-4 w-4 animate-spin" />
        ) : (
          <Upload className="h-4 w-4" />
        )}
        <span className="text-sm font-medium">
          {uploading ? "Uploading..." : isDragActive ? "Drop here" : "Upload file"}
        </span>
      </DropzoneTarget>

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

function AvatarBody() {
  const { isDragActive, uploading, files, removeFile } = useDropzoneContext();
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
          <DropzoneTarget className="absolute inset-0 rounded-full cursor-pointer" />
        </div>
      ) : (
        <DropzoneTarget
          className={` + "`" + `flex h-24 w-24 flex-col items-center justify-center gap-1 rounded-full border-2 border-dashed cursor-pointer transition-all ${targetState(
            isDragActive,
            uploading,
            "border-border hover:border-accent/50 hover:bg-bg-hover/30",
            "border-accent bg-accent/5"
          )}` + "`" + `}
        >
          {uploading ? (
            <Loader2 className="h-6 w-6 animate-spin text-accent" />
          ) : (
            <>
              <ImageIcon className="h-6 w-6 text-text-muted" />
              <span className="text-[10px] text-text-muted">Upload</span>
            </>
          )}
        </DropzoneTarget>
      )}
    </div>
  );
}

function InlineBody() {
  const { isDragActive, uploading, maxSize } = useDropzoneContext();
  return (
    <div className="space-y-3">
      <DropzoneTarget
        className={` + "`" + `flex items-center justify-between rounded-lg border px-4 py-3 cursor-pointer transition-all ${targetState(
          isDragActive,
          uploading,
          "border-border hover:border-accent/50",
          "border-accent bg-accent/5"
        )}` + "`" + `}
      >
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
      </DropzoneTarget>

      <DropzoneFileList />
    </div>
  );
}

// ── Named components ─────────────────────────────────────────────

/** The large dashed box, with a card per file below it. */
export function BoxDropzone({ progress = "bar", reorderable = false, ...props }: BoxDropzoneProps) {
  return (
    <DropzoneRoot {...props}>
      <BoxBody progress={progress} reorderable={reorderable} />
    </DropzoneRoot>
  );
}

/** A one-line strip, with a chip per file below it. */
export function CompactDropzone({ progress = "bar", ...props }: CompactDropzoneProps) {
  return (
    <DropzoneRoot {...props}>
      <CompactBody progress={progress} />
    </DropzoneRoot>
  );
}

/** An "Upload file" button, with the file names below it. */
export function MinimalDropzone(props: DropzoneBaseProps) {
  return (
    <DropzoneRoot {...props}>
      <MinimalBody />
    </DropzoneRoot>
  );
}

/** A round picture that is its own drop target. One file. */
export function AvatarDropzone(props: AvatarDropzoneProps) {
  return (
    <DropzoneRoot {...props} maxFiles={1}>
      <AvatarBody />
    </DropzoneRoot>
  );
}

/** A row with a Browse button, with a card per file below it. */
export function InlineDropzone(props: DropzoneBaseProps) {
  return (
    <DropzoneRoot {...props}>
      <InlineBody />
    </DropzoneRoot>
  );
}

// ── Dropzone: the look chosen by a prop ──────────────────────────

// props still carries provider, a display hint nothing ever displayed. The
// root takes what it knows from props and ignores the rest.
function VariantDropzone({ variant = "default", progress = "bar", reorderable = false, ...props }: DropzoneProps) {
  return (
    <DropzoneRoot {...props}>
      {variant === "compact" ? (
        <CompactBody progress={progress} />
      ) : variant === "minimal" ? (
        <MinimalBody />
      ) : variant === "avatar" ? (
        <AvatarBody />
      ) : variant === "inline" ? (
        <InlineBody />
      ) : (
        <BoxBody progress={progress} reorderable={reorderable} />
      )}
    </DropzoneRoot>
  );
}

export const Dropzone = Object.assign(VariantDropzone, {
  Root: DropzoneRoot,
  Target: DropzoneTarget,
  FileList: DropzoneFileList,
  Progress: DropzoneProgress,
});

// ── File preview ─────────────────────────────────────────────────

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

// v3.31.31: type-aware preview thumbnail. Images render the actual
// image, video shows a play badge over a dark thumb, audio gets a
// music icon, PDF / doc / excel get format-specific glyphs, anything
// else falls back to a generic file icon.
function PreviewThumb({ file }: { file: UploadedFile }) {
  const mime = file.type || "";
  if (mime.startsWith("image/")) {
    return (
      <img
        src={file.thumbnail_url || file.url}
        alt={file.name}
        className="h-10 w-10 rounded-md object-cover"
      />
    );
  }
  if (mime.startsWith("video/")) {
    return (
      <div className="relative h-10 w-10 rounded-md overflow-hidden bg-bg-tertiary">
        <div className="absolute inset-0 flex items-center justify-center bg-black/30">
          <Play className="h-3.5 w-3.5 text-white fill-white" />
        </div>
      </div>
    );
  }
  let Icon = File;
  let tint = "text-text-muted";
  if (mime.startsWith("audio/")) { Icon = Music; tint = "text-info"; }
  else if (mime === "application/pdf") { Icon = FileText; tint = "text-danger"; }
  else if (mime.includes("wordprocessing") || mime === "application/msword") { Icon = FileText; tint = "text-info"; }
  else if (mime.includes("spreadsheet") || mime === "application/vnd.ms-excel" || mime === "text/csv") { Icon = FileSpreadsheet; tint = "text-success"; }
  return (
    <div className="flex h-10 w-10 items-center justify-center rounded-md bg-bg-tertiary">
      <Icon className={"h-5 w-5 " + tint} />
    </div>
  );
}

interface FilePreviewProps {
  file: UploadedFile;
  onRemove: () => void;
  /** Multi-file list position (0-based). v3.31.31. */
  index: number;
  /** Total count in the list. */
  total: number;
  /** Whether to show up/down arrow buttons. */
  reorderable: boolean;
  /** Reorder callback (delta = -1 for up, 1 for down). */
  onMove: (index: number, direction: -1 | 1) => void;
}

function FilePreview({ file, onRemove, index, total, reorderable, onMove }: FilePreviewProps) {
  const canReorder = reorderable && total > 1;
  return (
    <div className="flex items-center gap-3 rounded-lg border border-border bg-bg-secondary px-3 py-2.5">
      <PreviewThumb file={file} />
      <div className="flex-1 min-w-0">
        <p className="text-sm font-medium text-foreground truncate">{file.name}</p>
        <p className="text-xs text-text-muted">
          {file.optimised && file.originalSize && file.originalSize > file.size ? (
            <>
              <span className="line-through opacity-60">{formatSize(file.originalSize)}</span>
              {" -> "}
              <span className="text-success">{formatSize(file.size)}</span>
              {file.format ? <span className="uppercase"> {file.format}</span> : null}
            </>
          ) : (
            formatSize(file.size)
          )}
        </p>
      </div>
      {canReorder && (
        <div className="flex flex-col gap-0.5">
          <button
            type="button"
            disabled={index === 0}
            onClick={(e) => { e.stopPropagation(); onMove(index, -1); }}
            className="rounded p-0.5 text-text-muted hover:text-foreground hover:bg-bg-hover transition-colors disabled:opacity-30 disabled:hover:bg-transparent"
            title="Move up"
          >
            <ChevronUp className="h-3.5 w-3.5" />
          </button>
          <button
            type="button"
            disabled={index === total - 1}
            onClick={(e) => { e.stopPropagation(); onMove(index, 1); }}
            className="rounded p-0.5 text-text-muted hover:text-foreground hover:bg-bg-hover transition-colors disabled:opacity-30 disabled:hover:bg-transparent"
            title="Move down"
          >
            <ChevronDown className="h-3.5 w-3.5" />
          </button>
        </div>
      )}
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

// v3.31.31: three progress looks, all driven by a 0..100 percent:
// "bar" (linear horizontal), "circular" (donut), "pulse" (dots and a number).

function UploadProgressBar({ percent }: { percent: number }) {
  return (
    <div className="flex flex-col items-center gap-2">
      <Loader2 className="h-8 w-8 animate-spin text-accent" />
      <p className="text-sm text-text-secondary">Uploading... {percent}%</p>
      <div className="w-48 h-1.5 bg-bg-tertiary rounded-full overflow-hidden">
        <div
          className="h-full bg-accent rounded-full transition-all duration-300"
          style={{ width: percent + "%" }}
        />
      </div>
    </div>
  );
}

function UploadProgressCircular({ percent }: { percent: number }) {
  // SVG donut. 36x36 viewbox, stroke 4, radius 14.
  // Circumference 2 * PI * 14 = 87.96; dashoffset = circumference * (1 - p).
  const radius = 14;
  const circumference = 2 * Math.PI * radius;
  const dashoffset = circumference * (1 - Math.min(1, Math.max(0, percent / 100)));
  return (
    <div className="flex flex-col items-center gap-2">
      <div className="relative h-12 w-12">
        <svg className="h-12 w-12 -rotate-90" viewBox="0 0 36 36">
          <circle
            cx="18"
            cy="18"
            r={radius}
            fill="none"
            className="stroke-bg-tertiary"
            strokeWidth="4"
          />
          <circle
            cx="18"
            cy="18"
            r={radius}
            fill="none"
            className="stroke-accent transition-[stroke-dashoffset] duration-300"
            strokeWidth="4"
            strokeLinecap="round"
            strokeDasharray={circumference}
            strokeDashoffset={dashoffset}
          />
        </svg>
        <div className="absolute inset-0 flex items-center justify-center">
          <span className="text-[10px] font-semibold tabular-nums text-foreground">{percent}%</span>
        </div>
      </div>
      <p className="text-sm text-text-secondary">Uploading...</p>
    </div>
  );
}

function UploadProgressPulse({ percent }: { percent: number }) {
  // Minimal: three pulsing dots and a percentage. No bar -- intended
  // for embedded contexts (compact dropzones inside dense tables)
  // where every pixel of chrome counts.
  return (
    <div className="flex items-center gap-2">
      <div className="flex items-center gap-1">
        <span className="h-1.5 w-1.5 rounded-full bg-accent animate-pulse" style={{ animationDelay: "0ms" }} />
        <span className="h-1.5 w-1.5 rounded-full bg-accent animate-pulse" style={{ animationDelay: "150ms" }} />
        <span className="h-1.5 w-1.5 rounded-full bg-accent animate-pulse" style={{ animationDelay: "300ms" }} />
      </div>
      <span className="text-xs font-medium text-text-secondary tabular-nums">Uploading {percent}%</span>
    </div>
  );
}
`
}

// adminDropzoneTest is the Vitest file for the dropzone: one test per named
// look, the variant prop, and a look composed from the parts.
func adminDropzoneTest() string {
	return `import { describe, it, expect, vi, beforeAll, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

// The upload goes to the API; here it answers the way /api/uploads does.
vi.mock("@/lib/api-client", () => ({
  uploadFile: vi.fn(async (file: File) => ({
    data: { url: "https://cdn.test/" + file.name, key: "uploads/" + file.name, name: file.name, size: file.size, mime: file.type },
  })),
}));

import {
  AvatarDropzone,
  BoxDropzone,
  CompactDropzone,
  Dropzone,
  InlineDropzone,
  MinimalDropzone,
  type UploadedFile,
} from "@/components/ui/dropzone";

beforeAll(() => {
  // jsdom has no object URLs, and a dropzone that does not upload makes one
  // per local preview.
  URL.createObjectURL = vi.fn(() => "blob:preview");
  URL.revokeObjectURL = vi.fn();
});

let onFilesChange: ReturnType<typeof vi.fn>;
beforeEach(() => {
  onFilesChange = vi.fn();
});

function picker(container: HTMLElement): HTMLInputElement {
  const input = container.querySelector('input[type="file"]');
  if (!(input instanceof HTMLInputElement)) throw new Error("no file input");
  return input;
}

const png = (name: string) => new File(["x"], name, { type: "image/png" });
const stored = (name: string, type = "image/png"): UploadedFile => ({
  url: "https://cdn.test/" + name,
  name,
  size: 2048,
  type,
});

describe("AvatarDropzone", () => {
  it("uploads one picture and shows it in the circle", async () => {
    const user = userEvent.setup();
    const { container } = render(
      <AvatarDropzone label="Photo" accept={{ "image/*": [".png"] }} onFilesChange={onFilesChange} />
    );
    expect(screen.getByText("Photo")).toBeInTheDocument();
    expect(screen.getByText("Upload")).toBeInTheDocument();

    await user.upload(picker(container), png("me.png"));

    await waitFor(() =>
      expect(onFilesChange).toHaveBeenCalledWith([
        expect.objectContaining({ url: "https://cdn.test/me.png", key: "uploads/me.png" }),
      ])
    );
    expect(screen.getByAltText("me.png")).toHaveAttribute("src", "https://cdn.test/me.png");
  });

  it("removes the current picture", async () => {
    const user = userEvent.setup();
    render(<AvatarDropzone value={[stored("current.png")]} onFilesChange={onFilesChange} />);
    expect(screen.getByAltText("current.png")).toBeInTheDocument();

    await user.click(screen.getByRole("button"));

    expect(onFilesChange).toHaveBeenCalledWith([]);
    expect(screen.queryByAltText("current.png")).not.toBeInTheDocument();
  });
});

describe("InlineDropzone", () => {
  it("lists the files it holds and keeps a local pick without uploading", async () => {
    const user = userEvent.setup();
    const { container } = render(
      <InlineDropzone autoUpload={false} maxFiles={3} value={[stored("a.pdf", "application/pdf")]} onFilesChange={onFilesChange} />
    );
    expect(screen.getByText("Choose files")).toBeInTheDocument();
    expect(screen.getByText("a.pdf")).toBeInTheDocument();
    expect(screen.getByText("2 KB")).toBeInTheDocument();

    await user.upload(picker(container), png("b.png"));

    await waitFor(() =>
      expect(onFilesChange).toHaveBeenCalledWith([
        expect.objectContaining({ name: "a.pdf" }),
        expect.objectContaining({ name: "b.png", url: "blob:preview" }),
      ])
    );
    // Inline cards have no reorder buttons.
    expect(screen.queryByTitle("Move up")).not.toBeInTheDocument();
  });
});

describe("BoxDropzone", () => {
  it("reorders its files when reorderable", async () => {
    const user = userEvent.setup();
    render(
      <BoxDropzone maxFiles={5} reorderable value={[stored("one.png"), stored("two.png")]} onFilesChange={onFilesChange} />
    );
    expect(screen.getByText("Click to upload or drag and drop")).toBeInTheDocument();

    await user.click(screen.getAllByTitle("Move down")[0]);

    expect(onFilesChange).toHaveBeenCalledWith([
      expect.objectContaining({ name: "two.png" }),
      expect.objectContaining({ name: "one.png" }),
    ]);
  });

  it("shows the error in place of the description", () => {
    render(<BoxDropzone description="PNG only" error="An image is required" />);
    expect(screen.getByText("An image is required")).toBeInTheDocument();
    expect(screen.queryByText("PNG only")).not.toBeInTheDocument();
  });
});

describe("CompactDropzone", () => {
  it("shows a chip per file and removes one", async () => {
    const user = userEvent.setup();
    render(<CompactDropzone maxFiles={2} value={[stored("clip.mp4", "video/mp4")]} onFilesChange={onFilesChange} />);
    expect(screen.getByText(/Drop files or click to browse/)).toBeInTheDocument();
    expect(screen.getByText("clip.mp4")).toBeInTheDocument();

    await user.click(screen.getByRole("button"));

    expect(onFilesChange).toHaveBeenCalledWith([]);
  });
});

describe("MinimalDropzone", () => {
  it("is an upload button that lists file names", () => {
    render(<MinimalDropzone value={[stored("notes.txt", "text/plain")]} />);
    expect(screen.getByText("Upload file")).toBeInTheDocument();
    expect(screen.getByText("notes.txt")).toBeInTheDocument();
  });
});

describe("Dropzone", () => {
  it("still picks the look from the variant prop", () => {
    const { unmount } = render(<Dropzone variant="avatar" value={[stored("face.png")]} />);
    expect(screen.getByAltText("face.png")).toBeInTheDocument();
    unmount();

    render(<Dropzone variant="inline" provider="minio" />);
    expect(screen.getByText("Choose files")).toBeInTheDocument();
  });

  it("builds a custom look from its parts", () => {
    render(
      <Dropzone.Root label="Receipts" maxFiles={4} value={[stored("receipt.png")]}>
        <Dropzone.Target className="receipts">Drop receipts here</Dropzone.Target>
        <Dropzone.Progress />
        <Dropzone.FileList />
      </Dropzone.Root>
    );
    expect(screen.getByText("Receipts")).toBeInTheDocument();
    expect(screen.getByText("Drop receipts here")).toHaveClass("receipts");
    expect(screen.getByText("receipt.png")).toBeInTheDocument();
    expect(screen.queryByText(/Uploading/)).not.toBeInTheDocument();
  });
});
`
}
