import { useCallback, useRef, useState } from "react";
import type { FileRef } from "./types";
import type { UploadOptions, createUploader } from "./uploader";

/**
 * React binding for the uploader.
 *
 * Deliberately not a TanStack Query mutation, though it composes with one. An
 * upload has per-file progress and partial success, and useMutation models a
 * single request with a single result. Use this for the upload itself and
 * TanStack Query for the record the refs end up attached to:
 *
 *   const { upload, items } = useUpload(uploader);
 *   const save = useMutation({ mutationFn: (refs: FileRef[]) => api.patch(...) });
 */

export type UploadState = "optimizing" | "uploading" | "done" | "error";

export interface UploadItem {
  id: string;
  name: string;
  state: UploadState;
  progress: number;
  /** The size of the file the user picked. */
  originalSize: number;
  /** What is actually being stored, once optimisation has run. */
  optimizedSize?: number;
  format?: string;
  error?: string;
  ref?: FileRef;
}

type Uploader = ReturnType<typeof createUploader>;

export function useUpload(uploader: Uploader) {
  const [items, setItems] = useState<UploadItem[]>([]);
  const counter = useRef(0);

  const patch = useCallback((id: string, next: Partial<UploadItem>) => {
    setItems((prev) => prev.map((it) => (it.id === id ? { ...it, ...next } : it)));
  }, []);

  const upload = useCallback(
    async (files: Array<{ blob: Blob; name: string }>, options: UploadOptions = {}) => {
      const started: UploadItem[] = files.map((f) => ({
        id: "u" + ++counter.current,
        name: f.name,
        state: "optimizing" as const,
        progress: 0,
        originalSize: f.blob.size,
      }));
      setItems((prev) => [...prev, ...started]);

      const refs: FileRef[] = [];
      for (let i = 0; i < files.length; i++) {
        const item = started[i];
        try {
          const ref = await uploader.upload(files[i].blob, files[i].name, {
            ...options,
            onOptimized: (res) => {
              patch(item.id, {
                state: "uploading",
                optimizedSize: res.primary.blob.size,
                format: res.primary.mime.replace("image/", ""),
              });
              options.onOptimized?.(res);
            },
            onProgress: (f) => {
              patch(item.id, { progress: f });
              options.onProgress?.(f);
            },
          });
          patch(item.id, { state: "done", progress: 1, ref });
          refs.push(ref);
        } catch (err) {
          // One bad file does not abandon the rest of the selection.
          patch(item.id, {
            state: "error",
            error: err instanceof Error ? err.message : String(err),
          });
        }
      }
      return refs;
    },
    [uploader, patch],
  );

  const reset = useCallback(() => setItems([]), []);
  const remove = useCallback(
    (id: string) => setItems((prev) => prev.filter((it) => it.id !== id)),
    [],
  );

  return { upload, items, reset, remove };
}

/**
 * "6.1 MB -> 41 KB", or null when there is nothing to boast about.
 *
 * Worth showing. The work happens before the upload starts, so without this
 * the only visible effect is that the progress bar finishes sooner.
 */
export function describeSaving(item: UploadItem): string | null {
  if (!item.optimizedSize || item.optimizedSize >= item.originalSize) return null;
  return formatBytes(item.originalSize) + " -> " + formatBytes(item.optimizedSize);
}

function formatBytes(bytes: number): string {
  if (bytes < 1024) return bytes + " B";
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(0) + " KB";
  return (bytes / (1024 * 1024)).toFixed(1) + " MB";
}
