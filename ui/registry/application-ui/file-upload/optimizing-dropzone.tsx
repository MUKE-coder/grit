'use client'

import { useCallback, useEffect, useId, useRef, useState } from 'react'
import { AlertCircle, CheckCircle2, FileIcon, ImageIcon, UploadCloud, X } from 'lucide-react'

/*
 * A dropzone that shrinks images before it uploads them.
 *
 * The work happens on the client, which is not a shortcut. Uploads go from the
 * browser straight to storage through a presigned URL and never pass through
 * an API, so bytes not removed here are bytes paid for on the user's
 * connection. A 5 MB phone photo becomes about 35 KB before anything is sent.
 * On mobile data that is most of the wait.
 *
 * The optimisation is @gritframework/upload. It is a peer rather than
 * something reimplemented inline because the parts that are easy to get wrong
 * are not the drag handlers: EXIF orientation has to be applied or portrait
 * photos arrive sideways, canvas.toBlob silently returns a PNG when asked for
 * an unsupported type so a "WebP" can quietly be larger than the JPEG it
 * replaced, and a transparent source must never be flattened onto a black
 * JPEG background.
 *
 * What this file is responsible for is the part a user sees.
 *
 * The drop target is a <label> wrapping a real file input, not a div with a
 * click handler. That gets keyboard activation, focus, the space and enter
 * keys and the mobile file picker for free, all of which a div has to
 * reimplement and usually reimplements incompletely.
 *
 * Progress is per file rather than one bar for the batch. A single bar for six
 * photos tells you nothing about which one is stuck, and the common failure
 * here is one file failing while the rest succeed.
 *
 * The saving is shown per file: "5.2 MB -> 41 KB WebP". Without it the only
 * visible effect of all this work is that the progress bar finishes sooner,
 * and the number is the most reassuring thing on the screen when somebody has
 * just dropped a photo they know is huge.
 *
 * Completion is announced through one polite live region rather than a toast
 * per file. Six toasts for six files is a queue you have to wait out.
 *
 * A file that fails does not remove itself. It stays in the list with what
 * went wrong, because a row that vanishes reads as success.
 */

/* --------------------------------------------------------------------------
 * The shapes this needs from @gritframework/upload. Declared structurally so
 * the block type-checks in a project that has not installed it yet, and so the
 * preview below can supply a stand-in.
 * ----------------------------------------------------------------------- */

export interface UploadedRef {
  url: string
  key: string
  name: string
  mime: string
  size: number
  thumbnail_url?: string
  format?: string
  optimised?: boolean
}

export interface DropzoneUploader {
  upload(
    file: Blob,
    filename: string,
    options?: {
      profile?: string
      accepts?: string[]
      onProgress?: (fraction: number) => void
      onOptimized?: (result: { primary: { blob: Blob; mime: string } }) => void
    },
  ): Promise<UploadedRef>
}

type ItemState = 'optimizing' | 'uploading' | 'done' | 'error'

interface Item {
  id: string
  name: string
  state: ItemState
  progress: number
  originalSize: number
  optimizedSize?: number
  format?: string
  previewUrl?: string
  error?: string
  ref?: UploadedRef
}

export interface OptimizingDropzoneProps {
  /** From createUploader() in @gritframework/upload. */
  uploader: DropzoneUploader
  /** Optimisation profile name, e.g. "product-image". */
  profile?: string
  /** MIME aliases the field accepts, e.g. ["image"]. */
  accepts?: string[]
  accept?: string
  maxFiles?: number
  disabled?: boolean
  label?: string
  hint?: string
  onChange?: (refs: UploadedRef[]) => void
}

function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(0)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

export function OptimizingDropzone({
  uploader,
  profile,
  accepts,
  accept = 'image/*',
  maxFiles = 8,
  disabled = false,
  label = 'Drop images here',
  hint = 'or click to browse. Large photos are shrunk before they upload.',
  onChange,
}: OptimizingDropzoneProps) {
  const inputId = useId()
  const [items, setItems] = useState<Item[]>([])
  const [dragging, setDragging] = useState(false)
  const [announcement, setAnnouncement] = useState('')
  const counter = useRef(0)
  /* Nested dragenter/dragleave fire for every child element, so a boolean flag
     flickers as the pointer crosses the icon or the text. Counting entries and
     leaves is the only version that stays steady. */
  const dragDepth = useRef(0)

  /* Object URLs are revoked on unmount. Without this every preview holds its
     full-size original in memory for the life of the page, which on a form
     with eight photos is most of a phone's budget. */
  useEffect(() => {
    return () => {
      for (const it of items) if (it.previewUrl) URL.revokeObjectURL(it.previewUrl)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const patch = useCallback((id: string, next: Partial<Item>) => {
    setItems((prev) => prev.map((it) => (it.id === id ? { ...it, ...next } : it)))
  }, [])

  const handleFiles = useCallback(
    async (files: File[]) => {
      if (disabled || files.length === 0) return

      const room = maxFiles - items.length
      const accepted = files.slice(0, Math.max(0, room))
      if (accepted.length === 0) return

      const started: Item[] = accepted.map((f) => ({
        id: `f${++counter.current}`,
        name: f.name,
        state: 'optimizing',
        progress: 0,
        originalSize: f.size,
        previewUrl: f.type.startsWith('image/') ? URL.createObjectURL(f) : undefined,
      }))
      setItems((prev) => [...prev, ...started])
      setAnnouncement(`Processing ${accepted.length} file${accepted.length === 1 ? '' : 's'}`)

      const done: UploadedRef[] = []
      for (let i = 0; i < accepted.length; i++) {
        const item = started[i]
        try {
          const ref = await uploader.upload(accepted[i], accepted[i].name, {
            profile,
            accepts,
            onOptimized: (res) =>
              patch(item.id, {
                state: 'uploading',
                optimizedSize: res.primary.blob.size,
                format: res.primary.mime.replace('image/', ''),
              }),
            onProgress: (f) => patch(item.id, { progress: f }),
          })
          patch(item.id, { state: 'done', progress: 1, ref })
          done.push(ref)
        } catch (err) {
          /* One bad file does not abandon the rest of the selection. */
          patch(item.id, {
            state: 'error',
            error: err instanceof Error ? err.message : String(err),
          })
        }
      }

      setAnnouncement(
        done.length === accepted.length
          ? `${done.length} file${done.length === 1 ? '' : 's'} uploaded`
          : `${done.length} of ${accepted.length} files uploaded`,
      )
      if (done.length) onChange?.(done)
    },
    [accepts, disabled, items.length, maxFiles, onChange, patch, profile, uploader],
  )

  const remove = useCallback((id: string) => {
    setItems((prev) => {
      const gone = prev.find((it) => it.id === id)
      if (gone?.previewUrl) URL.revokeObjectURL(gone.previewUrl)
      return prev.filter((it) => it.id !== id)
    })
  }, [])

  const full = items.length >= maxFiles

  return (
    <div className="w-full">
      <label
        htmlFor={inputId}
        onDragEnter={(e) => {
          e.preventDefault()
          dragDepth.current += 1
          setDragging(true)
        }}
        onDragOver={(e) => e.preventDefault()}
        onDragLeave={(e) => {
          e.preventDefault()
          dragDepth.current -= 1
          if (dragDepth.current <= 0) setDragging(false)
        }}
        onDrop={(e) => {
          e.preventDefault()
          dragDepth.current = 0
          setDragging(false)
          void handleFiles(Array.from(e.dataTransfer.files))
        }}
        className={[
          'relative flex cursor-pointer flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed px-6 py-10 text-center transition-colors',
          'focus-within:outline-none focus-within:ring-2 focus-within:ring-ring focus-within:ring-offset-2',
          dragging
            ? 'border-primary bg-primary/5'
            : 'border-border bg-muted/20 hover:border-primary/40 hover:bg-muted/30',
          disabled || full ? 'pointer-events-none opacity-60' : '',
        ].join(' ')}
      >
        <span
          className={[
            'flex h-11 w-11 items-center justify-center rounded-full transition-colors',
            dragging ? 'bg-primary/15 text-primary' : 'bg-muted text-muted-foreground',
          ].join(' ')}
        >
          <UploadCloud className="h-5 w-5" aria-hidden="true" />
        </span>
        <span className="space-y-1">
          <span className="block text-sm font-medium text-foreground">
            {full ? `Maximum of ${maxFiles} files reached` : label}
          </span>
          <span className="block text-xs text-muted-foreground">{full ? 'Remove one to add another.' : hint}</span>
        </span>

        {/* A real input, visually hidden rather than display:none, so it stays
            focusable and the label activates it with the keyboard. */}
        <input
          id={inputId}
          type="file"
          accept={accept}
          multiple={maxFiles > 1}
          disabled={disabled || full}
          className="sr-only"
          onChange={(e) => {
            void handleFiles(Array.from(e.target.files ?? []))
            /* Cleared so picking the same file twice in a row still fires a
               change event. */
            e.target.value = ''
          }}
        />
      </label>

      <p aria-live="polite" className="sr-only">
        {announcement}
      </p>

      {items.length > 0 && (
        <ul className="mt-3 space-y-2">
          {items.map((item) => (
            <li
              key={item.id}
              className="flex items-center gap-3 rounded-lg border border-border bg-card px-3 py-2.5"
            >
              <span className="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-md bg-muted">
                {item.previewUrl ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img src={item.previewUrl} alt="" className="h-full w-full object-cover" />
                ) : (
                  <FileIcon className="h-4 w-4 text-muted-foreground" aria-hidden="true" />
                )}
              </span>

              <span className="min-w-0 flex-1">
                <span className="flex items-center gap-2">
                  <span className="truncate text-sm font-medium text-foreground">{item.name}</span>
                  {item.state === 'done' && (
                    <CheckCircle2 className="h-3.5 w-3.5 shrink-0 text-emerald-500" aria-hidden="true" />
                  )}
                  {item.state === 'error' && (
                    <AlertCircle className="h-3.5 w-3.5 shrink-0 text-destructive" aria-hidden="true" />
                  )}
                </span>

                <span className="mt-0.5 block text-xs text-muted-foreground">
                  {item.state === 'error' ? (
                    <span className="text-destructive">{item.error}</span>
                  ) : item.optimizedSize && item.optimizedSize < item.originalSize ? (
                    <>
                      <span className="line-through opacity-60">{formatBytes(item.originalSize)}</span>{' '}
                      <span className="text-emerald-600 dark:text-emerald-400">
                        {formatBytes(item.optimizedSize)}
                      </span>
                      {item.format && <span className="uppercase"> · {item.format}</span>}
                    </>
                  ) : item.state === 'optimizing' ? (
                    'Optimising…'
                  ) : (
                    formatBytes(item.originalSize)
                  )}
                </span>

                {(item.state === 'uploading' || item.state === 'optimizing') && (
                  <span
                    role="progressbar"
                    aria-valuenow={Math.round(item.progress * 100)}
                    aria-valuemin={0}
                    aria-valuemax={100}
                    aria-label={`Uploading ${item.name}`}
                    className="mt-1.5 block h-1 w-full overflow-hidden rounded-full bg-muted"
                  >
                    <span
                      className="block h-full rounded-full bg-primary transition-[width] duration-200"
                      style={{ width: `${Math.max(4, item.progress * 100)}%` }}
                    />
                  </span>
                )}
              </span>

              <button
                type="button"
                onClick={() => remove(item.id)}
                aria-label={`Remove ${item.name}`}
                className="shrink-0 rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              >
                <X className="h-4 w-4" aria-hidden="true" />
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

/* --------------------------------------------------------------------------
 * Preview
 *
 * Wired to a stand-in uploader so this page can run without a backend. In a
 * real project the uploader is the one from @gritframework/upload:
 *
 *   import { createUploader, createFetchTransport } from '@gritframework/upload'
 *   import { optimizeImage } from '@gritframework/upload/web'
 *
 *   const uploader = createUploader({
 *     transport: createFetchTransport('/api/v1', () => getToken()),
 *     optimize: optimizeImage,
 *   })
 *
 * The stand-in reports a plausible saving so the layout is honest about what
 * the real one shows. It does not invent a ratio for non-images.
 * ----------------------------------------------------------------------- */

const demoUploader: DropzoneUploader = {
  async upload(file, filename, options) {
    const isImage = file.type.startsWith('image/')
    const optimized = isImage ? Math.max(18_000, Math.round(file.size * 0.03)) : file.size

    await new Promise((r) => setTimeout(r, 350))
    options?.onOptimized?.({
      primary: { blob: new Blob([new Uint8Array(optimized)]), mime: isImage ? 'image/webp' : file.type },
    })

    for (let p = 0; p <= 1.0001; p += 0.2) {
      options?.onProgress?.(Math.min(1, p))
      await new Promise((r) => setTimeout(r, 120))
    }

    return {
      url: `https://cdn.example.com/uploads/${filename}`,
      key: `uploads/${filename}`,
      name: filename,
      mime: isImage ? 'image/webp' : file.type,
      size: optimized,
      format: isImage ? 'webp' : undefined,
      optimised: isImage,
    }
  },
}

export default function Preview() {
  return (
    <div className="mx-auto w-full max-w-xl p-6">
      <p className="mb-2 text-sm font-medium text-foreground">Product images</p>
      <OptimizingDropzone uploader={demoUploader} profile="product-image" accepts={['image']} maxFiles={4} />
    </div>
  )
}
