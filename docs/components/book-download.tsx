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
export const BOOK_PAGES = 503
export const BOOK_CHAPTERS = 61
export const BOOK_TARGET_PAGES = 500
export const BOOK_HREF = '/book/go-for-javascript-developers.pdf'
export const BOOK_TITLE = 'Go for JavaScript Developers'

/**
 * The leaderboard the book's assignments are scored on. Chapter 0 sets it up,
 * and the join code is the same one printed there: if it ever changes, both
 * have to change together, so the code lives next to the link rather than
 * inside a sentence.
 */
export const LEAGUE_CODE = '5TCLJ27Z'
export const LEAGUE_JOIN_HREF = `https://committed-api.gritframework.dev/join/${LEAGUE_CODE}`

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
          'group inline-flex items-center gap-2 text-sm text-accent hover:underline',
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
        'rounded-xl border border-accent/20 bg-accent/5 p-5 sm:p-6',
        className,
      )}
    >
      <div className="flex flex-col gap-4 sm:flex-row sm:items-start">
        <span className="inline-flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-accent/10 text-accent">
          <BookOpen className="h-5 w-5" />
        </span>
        <div className="min-w-0 flex-1">
          {/* A <p>, not an <h3>.
              This card sits above the page's first real section, so a level-3
              heading here made every page carrying it skip h1 -> h3, which is
              what a screen reader user navigating by heading level hears as a
              missing section. It is a promotional card, not a section of the
              document. */}
          <p className="mb-1 text-base font-semibold text-foreground">
            Free book: {BOOK_TITLE}
          </p>
          <p className="mb-4 text-sm leading-relaxed text-muted-foreground">
            A textbook that teaches Go from scratch to someone who already ships React and
            Next.js, then builds real applications with Grit. Every concept starts from the
            JavaScript you already know, and every chapter ends with something you build.
          </p>
          <p className="mb-4 text-sm leading-relaxed text-muted-foreground">
            Every chapter ends with an assignment you push to your own GitHub repository.
            Chapter 0 sets that up and gets you onto{' '}
            <a
              href={LEAGUE_JOIN_HREF}
              target="_blank"
              rel="noreferrer"
              className="font-medium text-accent hover:underline"
            >
              Gopher League
            </a>
            , a leaderboard for developers learning Go together, where every commit,
            review and merged pull request scores.
          </p>
          <div className="flex flex-wrap items-center gap-3">
            <Link
              href={BOOK_HREF}
              className="group inline-flex h-9 items-center gap-2 rounded-full bg-primary px-5 text-sm font-semibold text-primary-foreground transition-colors hover:bg-fg/85"
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
