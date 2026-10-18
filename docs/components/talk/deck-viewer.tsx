'use client'

import { useCallback, useEffect, useLayoutEffect, useRef, useState } from 'react'
import {
  ChevronLeft,
  ChevronRight,
  Grid2x2,
  Maximize,
  Minimize,
  NotebookText,
  X,
} from 'lucide-react'
import { cn } from '@/lib/utils'
import type { TalkDeck, TalkSlide } from '@/lib/talk-deck'

// The slides are authored on a fixed canvas and every length in them is in px,
// so the viewer never reflows a slide: it scales the whole 1920x1080 box to
// whatever width it has. That is the only way the deck on the site is the deck
// that gets presented.
const CANVAS_W = 1920
const CANVAS_H = 1080

/** One slide, scaled to fill its container. */
function Stage({
  slide,
  className,
  ariaHidden,
}: {
  slide: TalkSlide
  className?: string
  ariaHidden?: boolean
}) {
  const box = useRef<HTMLDivElement>(null)
  const [scale, setScale] = useState(0)

  useLayoutEffect(() => {
    const el = box.current
    if (!el) return

    const measure = () => setScale(el.clientWidth / CANVAS_W)
    measure()

    // ResizeObserver rather than a window listener: the stage also changes
    // width when the notes panel opens, which resizes nothing globally.
    if (typeof ResizeObserver === 'undefined') {
      window.addEventListener('resize', measure)
      return () => window.removeEventListener('resize', measure)
    }
    const ro = new ResizeObserver(measure)
    ro.observe(el)
    return () => ro.disconnect()
  }, [])

  return (
    <div
      ref={box}
      className={cn('relative w-full overflow-hidden', className)}
      style={{ aspectRatio: '16 / 9' }}
      aria-hidden={ariaHidden}
    >
      {/* Until the first measurement the slide would render at full 1920px and
          overflow, so hold it back one frame. */}
      {scale > 0 && (
        <div
          style={{
            width: CANVAS_W,
            height: CANVAS_H,
            transform: `scale(${scale})`,
            transformOrigin: 'top left',
          }}
          dangerouslySetInnerHTML={{ __html: slide.html }}
        />
      )}
    </div>
  )
}

export function DeckViewer({ deck }: { deck: TalkDeck }) {
  const { slides, sections } = deck
  const [current, setCurrent] = useState(0)
  const [notesOpen, setNotesOpen] = useState(false)
  const [overviewOpen, setOverviewOpen] = useState(false)
  const [isFullscreen, setIsFullscreen] = useState(false)
  const frame = useRef<HTMLDivElement>(null)

  const slide = slides[current]

  const go = useCallback(
    (next: number) => {
      setCurrent((prev) => {
        const clamped = Math.max(0, Math.min(slides.length - 1, next))
        if (clamped !== prev) {
          // replaceState, not a hash assignment: assigning to location.hash
          // makes the browser scroll to the element and adds a history entry
          // for every arrow press.
          window.history.replaceState(null, '', `#${slides[clamped].id}`)
        }
        return clamped
      })
    },
    [slides],
  )

  // Open on the slide in the URL, so #d7 and a shared link land in one place.
  //
  // The listener is the half that is easy to leave out: the outline below the
  // deck links to #sec-biz, and clicking it from this page is a fragment
  // change, not a load. Without this the URL moves and the deck does not.
  useEffect(() => {
    const fromHash = (scroll: boolean) => {
      const id = window.location.hash.replace('#', '')
      if (!id) return
      const i = slides.findIndex((s) => s.id === id)
      if (i < 0) return
      setCurrent(i)
      if (scroll) frame.current?.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }
    fromHash(false)

    const onHashChange = () => fromHash(true)
    window.addEventListener('hashchange', onHashChange)
    return () => window.removeEventListener('hashchange', onHashChange)
  }, [slides])

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      // Never swallow a key someone is typing into the search box.
      const target = e.target as HTMLElement | null
      if (target && /^(INPUT|TEXTAREA|SELECT)$/.test(target.tagName)) return
      if (target?.isContentEditable) return

      switch (e.key) {
        case 'ArrowRight':
        case 'PageDown':
        case ' ':
          e.preventDefault()
          go(current + 1)
          break
        case 'ArrowLeft':
        case 'PageUp':
          e.preventDefault()
          go(current - 1)
          break
        case 'Home':
          e.preventDefault()
          go(0)
          break
        case 'End':
          e.preventDefault()
          go(slides.length - 1)
          break
        case 'n':
        case 'N':
          setNotesOpen((v) => !v)
          break
        case 'g':
        case 'G':
          setOverviewOpen((v) => !v)
          break
        case 'Escape':
          setOverviewOpen(false)
          break
        default:
          break
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [current, go, slides.length])

  // The browser can leave fullscreen without us (Escape, the window button),
  // so the button's icon follows the document rather than our own state.
  useEffect(() => {
    const sync = () => setIsFullscreen(Boolean(document.fullscreenElement))
    document.addEventListener('fullscreenchange', sync)
    return () => document.removeEventListener('fullscreenchange', sync)
  }, [])

  const toggleFullscreen = () => {
    if (document.fullscreenElement) {
      void document.exitFullscreen()
    } else {
      void frame.current?.requestFullscreen?.()
    }
  }

  const sectionNow = [...sections].reverse().find((s) => s.number <= slide.number)

  return (
    <div className="flex flex-col gap-4">
      <div
        ref={frame}
        className="flex flex-col gap-0 overflow-hidden rounded-xl border border-border bg-black"
      >
        <div className="flex flex-col lg:flex-row">
          <Stage slide={slide} className={notesOpen ? 'lg:w-[68%]' : 'w-full'} />

          {notesOpen && (
            <div className="flex max-h-[420px] flex-1 flex-col gap-3 overflow-y-auto border-t border-border bg-background p-5 lg:max-h-none lg:border-l lg:border-t-0">
              <p className="font-mono text-[11px] uppercase tracking-[0.2em] text-muted-foreground">
                Speaker notes
              </p>
              {slide.notes.length === 0 ? (
                <p className="text-sm text-muted-foreground">No notes on this slide.</p>
              ) : (
                slide.notes.map((note, i) => (
                  <p key={i} className="text-sm leading-relaxed text-foreground/90">
                    {note}
                  </p>
                ))
              )}
            </div>
          )}
        </div>

        {/* Progress. Two divs rather than a <progress>, because a styled
            progress element is more work than it is worth here. */}
        <div className="h-[3px] w-full bg-border">
          <div
            className="h-full bg-primary transition-[width] duration-200"
            style={{ width: `${((current + 1) / slides.length) * 100}%` }}
          />
        </div>

        <div className="flex flex-wrap items-center gap-2 border-t border-border bg-background px-3 py-2">
          <button
            type="button"
            onClick={() => go(current - 1)}
            disabled={current === 0}
            aria-label="Previous slide"
            className="inline-flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent/50 hover:text-foreground disabled:pointer-events-none disabled:opacity-40"
          >
            <ChevronLeft className="h-4 w-4" />
          </button>
          <button
            type="button"
            onClick={() => go(current + 1)}
            disabled={current === slides.length - 1}
            aria-label="Next slide"
            className="inline-flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent/50 hover:text-foreground disabled:pointer-events-none disabled:opacity-40"
          >
            <ChevronRight className="h-4 w-4" />
          </button>

          <p
            className="font-mono text-xs text-muted-foreground"
            aria-live="polite"
            aria-atomic="true"
          >
            {slide.number} / {slides.length}
            <span className="hidden sm:inline"> &middot; {slide.id}</span>
          </p>

          {sectionNow && (
            <p className="hidden min-w-0 flex-1 truncate text-xs text-muted-foreground md:block">
              {sectionNow.description}
            </p>
          )}

          <div className="ml-auto flex items-center gap-1">
            <button
              type="button"
              onClick={() => setOverviewOpen(true)}
              className="inline-flex h-8 items-center gap-1.5 rounded-md px-2.5 text-xs text-muted-foreground transition-colors hover:bg-accent/50 hover:text-foreground"
            >
              <Grid2x2 className="h-3.5 w-3.5" />
              <span className="hidden sm:inline">All slides</span>
            </button>
            <button
              type="button"
              onClick={() => setNotesOpen((v) => !v)}
              aria-pressed={notesOpen}
              className={cn(
                'inline-flex h-8 items-center gap-1.5 rounded-md px-2.5 text-xs transition-colors hover:bg-accent/50',
                notesOpen ? 'text-primary' : 'text-muted-foreground hover:text-foreground',
              )}
            >
              <NotebookText className="h-3.5 w-3.5" />
              <span className="hidden sm:inline">Notes</span>
            </button>
            <button
              type="button"
              onClick={toggleFullscreen}
              aria-label={isFullscreen ? 'Leave fullscreen' : 'Present fullscreen'}
              className="inline-flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent/50 hover:text-foreground"
            >
              {isFullscreen ? <Minimize className="h-4 w-4" /> : <Maximize className="h-4 w-4" />}
            </button>
          </div>
        </div>
      </div>

      <p className="text-center text-xs text-muted-foreground">
        Arrow keys to page through.{' '}
        <kbd className="rounded border border-border px-1 py-0.5 font-mono text-[10px]">N</kbd> for
        notes,{' '}
        <kbd className="rounded border border-border px-1 py-0.5 font-mono text-[10px]">G</kbd> for
        every slide at once.
      </p>

      {overviewOpen && (
        <div
          role="dialog"
          aria-modal="true"
          aria-label="All slides"
          className="fixed inset-0 z-50 overflow-y-auto bg-background/95 backdrop-blur-sm"
        >
          <div className="sticky top-0 z-10 flex items-center justify-between border-b border-border bg-background/90 px-5 py-3">
            <p className="font-mono text-[11px] uppercase tracking-[0.2em] text-muted-foreground">
              {slides.length} slides
            </p>
            <button
              type="button"
              onClick={() => setOverviewOpen(false)}
              aria-label="Close"
              className="inline-flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent/50 hover:text-foreground"
            >
              <X className="h-4 w-4" />
            </button>
          </div>

          <div className="grid grid-cols-1 gap-4 p-5 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {slides.map((s, i) => (
              <button
                key={s.id}
                type="button"
                onClick={() => {
                  go(i)
                  setOverviewOpen(false)
                }}
                className={cn(
                  'group overflow-hidden rounded-lg border text-left transition-colors',
                  i === current
                    ? 'border-primary'
                    : 'border-border hover:border-muted-foreground/50',
                )}
              >
                {/* The thumbnail is the slide itself, scaled. A separate
                    thumbnail would be a second thing to keep in step. */}
                <Stage slide={s} ariaHidden />
                <p className="flex items-baseline gap-2 border-t border-border bg-background px-2.5 py-1.5 font-mono text-[11px] text-muted-foreground">
                  <span className="text-foreground">{s.number}</span>
                  <span className="truncate">{s.id}</span>
                </p>
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
