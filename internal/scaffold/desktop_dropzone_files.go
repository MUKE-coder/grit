package scaffold

// desktopClientFileDropzone is the offline form's file/image field. It uploads
// through the API (POST /uploads) when the app is online and stores the
// resulting FileRef(s) in the form payload. Uploads need a connection — when
// offline it disables itself with a clear hint, since a blob can't be pushed
// through the JSON sync outbox. No extra deps: a hidden <input> + drag events.
func desktopClientFileDropzone() string {
	return `import { useRef, useState } from "react";
import { Upload, X, FileIcon, Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { uploadFile, type FileRef } from "@/lib/api-client";
import { useSyncStatus } from "@/hooks/use-sync-status";

interface FileDropzoneProps {
  value: FileRef | FileRef[] | null;
  onChange: (value: FileRef | FileRef[] | null) => void;
  multiple?: boolean;
  accept?: string;
  label?: string;
}

function isImage(f: FileRef): boolean {
  return !!f.thumbnail_url || (f.mime ?? "").startsWith("image/");
}

// A ref whose url is an inline data: URL hasn't been uploaded yet — it was
// picked offline and is waiting for the reconnect reconciler.
function isPending(f: FileRef): boolean {
  return typeof f.url === "string" && f.url.startsWith("data:");
}

function readAsDataUrl(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const r = new FileReader();
    r.onload = () => resolve(String(r.result));
    r.onerror = reject;
    r.readAsDataURL(file);
  });
}

export function FileDropzone({ value, onChange, multiple, accept, label }: FileDropzoneProps) {
  const { status } = useSyncStatus();
  const online = status.reachable && !status.force_offline;
  const inputRef = useRef<HTMLInputElement>(null);
  const [dragging, setDragging] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const items: FileRef[] = value ? (Array.isArray(value) ? value : [value]) : [];

  const emit = (next: FileRef[]) => {
    if (multiple) onChange(next);
    else onChange(next[0] ?? null);
  };

  const handleFiles = async (fileList: FileList | null) => {
    if (!fileList || fileList.length === 0) return;
    setError(null);
    setBusy(true);
    try {
      const picked = multiple ? Array.from(fileList) : [fileList[0]];
      const refs: FileRef[] = [];
      for (const file of picked) {
        if (online) {
          refs.push(await uploadFile(file));
        } else {
          // Offline: keep the file inline as a data: URL so it saves and
          // renders now. usePendingUploads() uploads it and swaps in the real
          // FileRef once the connection is back.
          const dataUrl = await readAsDataUrl(file);
          refs.push({ url: dataUrl, name: file.name, mime: file.type, size: file.size });
        }
      }
      emit(multiple ? [...items, ...refs] : refs);
    } catch {
      setError(online ? "Upload failed. Check your connection and try again." : "Could not read that file.");
    } finally {
      setBusy(false);
      if (inputRef.current) inputRef.current.value = "";
    }
  };

  const removeAt = (i: number) => emit(items.filter((_, j) => j !== i));

  return (
    <div>
      {label && <label className="mb-1.5 block text-[13px] font-medium text-foreground">{label}</label>}

      <div
        onClick={() => inputRef.current?.click()}
        onDragOver={(e) => { e.preventDefault(); setDragging(true); }}
        onDragLeave={() => setDragging(false)}
        onDrop={(e) => { e.preventDefault(); setDragging(false); handleFiles(e.dataTransfer.files); }}
        className={cn(
          "flex cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed px-4 py-6 text-center transition-colors",
          dragging ? "border-accent bg-accent/5" : "border-border hover:border-accent/50 hover:bg-surface-hover/30",
        )}
      >
        <input
          ref={inputRef}
          type="file"
          accept={accept}
          multiple={multiple}
          className="hidden"
          onChange={(e) => handleFiles(e.target.files)}
        />
        {busy ? (
          <Loader2 className="h-5 w-5 animate-spin text-accent" />
        ) : (
          <span className="inline-flex h-9 w-9 items-center justify-center rounded-lg bg-surface-2 text-foreground-secondary">
            <Upload className="h-4 w-4" />
          </span>
        )}
        <p className="mt-2 text-[12px] text-foreground-secondary">
          {online ? "Click or drag to upload" : "Click or drag — uploads when you reconnect"}
        </p>
      </div>

      {error && <p className="mt-1.5 text-[12px] text-danger">{error}</p>}

      {items.length > 0 && (
        <ul className="mt-3 space-y-2">
          {items.map((f, i) => (
            <li key={i} className="flex items-center gap-3 rounded-lg border border-border bg-surface-2 px-3 py-2">
              {isImage(f) ? (
                <img src={f.thumbnail_url || f.url} alt="" className="h-9 w-9 shrink-0 rounded object-cover" />
              ) : (
                <span className="inline-flex h-9 w-9 shrink-0 items-center justify-center rounded bg-surface-3 text-foreground-secondary">
                  <FileIcon className="h-4 w-4" />
                </span>
              )}
              <span className="min-w-0 flex-1 truncate text-[13px] text-foreground">{f.name}</span>
              {isPending(f) && (
                <span className="shrink-0 rounded-full bg-warning/10 px-2 py-0.5 text-[11px] font-medium text-warning">Pending</span>
              )}
              <button
                type="button"
                onClick={() => removeAt(i)}
                className="rounded p-1 text-foreground-muted hover:bg-danger/10 hover:text-danger"
                title="Remove"
              >
                <X className="h-3.5 w-3.5" />
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
`
}
