import Link from 'next/link'
import { BookOpen, Download } from 'lucide-react'
import { cn } from '@/lib/utils'

/**
 * The download card for the book, on every page where a reader is likely to
 * want it.
 *
 * Status lives here and nowhere else. The book is written in the open and
 * ships as it grows, so the page count on the card is a fact that changes
 * every few chapters: keeping it in one constant means updating it is one edit
 * rather than a hunt through five pages that have quietly disagreed for weeks.
 *
 * Update BOOK_PAGES and BOOK_CHAPTERS when BOOK/build.py reports new numbers.
 */
export const BOOK_PAGES = 101
export const BOOK_CHAPTERS = 13
export const BOOK_TARGET_PAGES = 500
export const BOOK_HREF = '/book/go-for-javascript-developers.pdf'
export const BOOK_TITLE = 'Go for JavaScript Developers'

export function BookDownload({
  className,
  compact = false,
}: {
  className?: string
  compact?: boolean
}) {
  const complete = BOOK_PAGES >= BOOK_TARGET_PAGES

  if (compact) {
    return (
      <Link
        href={BOOK_HREF}
        className={cn(
          'group inline-flex items-center gap-2 text-sm text-primary hover:underline',
          className,
        )}
      >
        <BookOpen className="h-4 w-4" />
        Free book: {BOOK_TITLE} (PDF)
      </Link>
    )
  }

  return (
    <div
      className={cn(
        'rounded-xl border border-primary/20 bg-primary/5 p-5 sm:p-6',
        className,
      )}
    >
      <div className="flex flex-col gap-4 sm:flex-row sm:items-start">
        <span className="inline-flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
          <BookOpen className="h-5 w-5" />
        </span>
        <div className="min-w-0 flex-1">
          <h3 className="mb-1 text-base font-semibold text-foreground">
            Free book: {BOOK_TITLE}
          </h3>
          <p className="mb-4 text-sm leading-relaxed text-muted-foreground">
            A textbook that teaches Go from scratch to someone who already ships React and
            Next.js, then builds real applications with Grit. Every concept starts from the
            JavaScript you already know, and every chapter ends with something you build.
          </p>
          <div className="flex flex-wrap items-center gap-3">
            <Link
              href={BOOK_HREF}
              className="group inline-flex h-9 items-center gap-2 rounded-full bg-primary px-5 text-sm font-semibold text-primary-foreground transition-colors hover:bg-primary/90"
            >
              <Download className="h-3.5 w-3.5" />
              Download the PDF
            </Link>
            <span className="font-mono text-xs text-muted-foreground">
              {complete
                ? `${BOOK_PAGES} pages`
                : `${BOOK_PAGES} pages, ${BOOK_CHAPTERS} chapters so far`}
            </span>
          </div>
          {!complete && (
            // Saying so is not modesty, it is the difference between an early
            // edition and a download that misrepresents itself.
            <p className="mt-3 text-xs leading-relaxed text-muted-foreground/80">
              Early edition, written in the open and growing toward{' '}
              {BOOK_TARGET_PAGES} pages. Download it again any time for the latest chapters.
            </p>
          )}
        </div>
      </div>
    </div>
  )
}
