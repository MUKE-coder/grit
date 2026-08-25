import { useCallback, useEffect, useId, useRef, useState } from "react";
import { useUpload, describeSaving, type UploadItem } from "./react";
import type { FileRef } from "./types";
import type { createUploader } from "./uploader";

/**
 * A ready-made dropzone.
 *
 * The package is otherwise logic only, and that is usually the right shape for
 * a library: shipping opinionated CSS into somebody else's design system is a
 * fight nobody wins. This exists because the alternative is that every project
 * rewrites the same forty lines of drag handling, progress rows and object-URL
 * cleanup, and most of them get the last one wrong.
 *
 * The compromise is markup with stable class names and no styling of its own.
 * Every element takes a class you can target, and the classNames prop replaces
 * any of them outright. If you want it to look like something without doing
 * that work:
 *
 *   import "@gritframework/upload/styles.css";
 *
 * That sheet is optional, plain CSS, themed through custom properties, and
 * scoped to the grit-dz namespace so it cannot leak into the rest of a page.
 *
 * If you use Tailwind and want a version to own and edit rather than
 * configure, take the Grit UI block instead:
 *   npx shadcn@latest add https://ui.gritframework.dev/r/application-ui-file-upload-optimizing-dropzone.json
 */

type Uploader = ReturnType<typeof createUploader>;

export interface DropzoneClassNames {
  root?: string;
  zone?: string;
  zoneDragging?: string;
  icon?: string;
  label?: string;
  hint?: string;
  list?: string;
  item?: string;
  thumb?: string;
  body?: string;
  name?: string;
  meta?: string;
  saving?: string;
  error?: string;
  progress?: string;
  progressBar?: string;
  remove?: string;
}

export interface DropzoneProps {
  /** From createUploader(). */
  uploader: Uploader;
  /** Optimisation profile name, e.g. "product-image". */
  profile?: string;
  /** MIME aliases the field accepts, e.g. ["image"]. */
  accepts?: string[];
  /** The input's accept attribute. */
  accept?: string;
  maxFiles?: number;
  disabled?: boolean;
  label?: string;
  hint?: string;
  /** Called with the refs of everything that uploaded successfully. */
  onChange?: (refs: FileRef[]) => void;
  classNames?: DropzoneClassNames;
  /** Replace a row entirely, keeping the drop target and the upload logic. */
  renderItem?: (item: UploadItem, remove: () => void) => React.ReactNode;
}

function cx(...parts: Array<string | undefined | false>): string {
  return parts.filter(Boolean).join(" ");
}

function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(0)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

export function Dropzone({
  uploader,
  profile,
  accepts,
  accept = "image/*",
  maxFiles = 8,
  disabled = false,
  label = "Drop files here",
  hint = "or click to browse",
  onChange,
  classNames: cn = {},
  renderItem,
}: DropzoneProps) {
  const inputId = useId();
  const { upload, items, remove } = useUpload(uploader);
  const [dragging, setDragging] = useState(false);
  const [previews, setPreviews] = useState<Record<string, string>>({});

  /**
   * Nested dragenter and dragleave fire for every child element, so a boolean
   * flickers as the pointer crosses the icon or the text. Counting entries and
   * leaves is the only version that stays steady.
   */
  const depth = useRef(0);

  /**
   * Object URLs are revoked on unmount. Without this every preview pins its
   * full-size original in memory for the life of the page, which on a form
   * with eight photos is most of a phone's budget.
   */
  const previewsRef = useRef<Record<string, string>>({});
  previewsRef.current = previews;
  useEffect(() => {
    return () => {
      for (const url of Object.values(previewsRef.current)) URL.revokeObjectURL(url);
    };
  }, []);

  const handle = useCallback(
    async (files: File[]) => {
      if (disabled || files.length === 0) return;
      const room = maxFiles - items.length;
      const accepted = files.slice(0, Math.max(0, room));
      if (accepted.length === 0) return;

      // Keyed by name and size rather than by the item id, which useUpload
      // assigns internally and does not hand back before the upload starts.
      setPreviews((prev) => {
        const next = { ...prev };
        for (const f of accepted) {
          if (f.type.startsWith("image/")) next[f.name + ":" + f.size] = URL.createObjectURL(f);
        }
        return next;
      });

      const refs = await upload(
        accepted.map((f) => ({ blob: f, name: f.name })),
        { profile, accepts },
      );
      if (refs.length) onChange?.(refs);
    },
    [accepts, disabled, items.length, maxFiles, onChange, profile, upload],
  );

  const full = items.length >= maxFiles;

  return (
    <div className={cx("grit-dz", cn.root)}>
      <label
        htmlFor={inputId}
        onDragEnter={(e) => {
          e.preventDefault();
          depth.current += 1;
          setDragging(true);
        }}
        onDragOver={(e) => e.preventDefault()}
        onDragLeave={(e) => {
          e.preventDefault();
          depth.current -= 1;
          if (depth.current <= 0) setDragging(false);
        }}
        onDrop={(e) => {
          e.preventDefault();
          depth.current = 0;
          setDragging(false);
          void handle(Array.from(e.dataTransfer.files));
        }}
        className={cx(
          "grit-dz__zone",
          dragging && "grit-dz__zone--dragging",
          (disabled || full) && "grit-dz__zone--disabled",
          cn.zone,
          dragging && cn.zoneDragging,
        )}
        data-dragging={dragging || undefined}
        data-disabled={disabled || full || undefined}
      >
        <span className={cx("grit-dz__label", cn.label)}>
          {full ? `Maximum of ${maxFiles} files reached` : label}
        </span>
        <span className={cx("grit-dz__hint", cn.hint)}>
          {full ? "Remove one to add another." : hint}
        </span>

        {/*
          A real input behind a label, not a div with a click handler. That is
          what gives keyboard activation, focus, space and enter, and the
          native mobile picker, none of which a div gets without reimplementing
          them. Visually hidden rather than display:none so it stays focusable.
        */}
        <input
          id={inputId}
          type="file"
          accept={accept}
          multiple={maxFiles > 1}
          disabled={disabled || full}
          className="grit-dz__input"
          onChange={(e) => {
            void handle(Array.from(e.target.files ?? []));
            // Cleared so picking the same file twice in a row still fires.
            e.target.value = "";
          }}
        />
      </label>

      {items.length > 0 && (
        <ul className={cx("grit-dz__list", cn.list)}>
          {items.map((item) => {
            if (renderItem) {
              return <li key={item.id}>{renderItem(item, () => remove(item.id))}</li>;
            }
            const preview = previews[item.name + ":" + item.originalSize];
            const saving = describeSaving(item);
            return (
              <li key={item.id} className={cx("grit-dz__item", cn.item)} data-state={item.state}>
                {preview ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img src={preview} alt="" className={cx("grit-dz__thumb", cn.thumb)} />
                ) : (
                  <span className={cx("grit-dz__thumb", cn.thumb)} aria-hidden="true" />
                )}

                <span className={cx("grit-dz__body", cn.body)}>
                  <span className={cx("grit-dz__name", cn.name)}>{item.name}</span>
                  <span className={cx("grit-dz__meta", cn.meta)}>
                    {item.state === "error" ? (
                      <span className={cx("grit-dz__error", cn.error)}>{item.error}</span>
                    ) : saving ? (
                      // The whole point of optimising: without this the only
                      // visible effect is that the bar finishes sooner.
                      <span className={cx("grit-dz__saving", cn.saving)}>
                        {saving}
                        {item.format ? ` · ${item.format.toUpperCase()}` : ""}
                      </span>
                    ) : item.state === "optimizing" ? (
                      "Optimising…"
                    ) : (
                      formatBytes(item.originalSize)
                    )}
                  </span>

                  {(item.state === "uploading" || item.state === "optimizing") && (
                    <span
                      role="progressbar"
                      aria-valuenow={Math.round(item.progress * 100)}
                      aria-valuemin={0}
                      aria-valuemax={100}
                      aria-label={`Uploading ${item.name}`}
                      className={cx("grit-dz__progress", cn.progress)}
                    >
                      <span
                        className={cx("grit-dz__progress-bar", cn.progressBar)}
                        style={{ width: `${Math.max(4, item.progress * 100)}%` }}
                      />
                    </span>
                  )}
                </span>

                <button
                  type="button"
                  onClick={() => remove(item.id)}
                  aria-label={`Remove ${item.name}`}
                  className={cx("grit-dz__remove", cn.remove)}
                >
                  ×
                </button>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}
