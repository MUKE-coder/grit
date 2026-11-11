'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { ArrowRight, BookOpen, Home, RefreshCw, RotateCcw } from 'lucide-react'
import { GridFrame } from '@/components/grid-frame'

/*
 * The route-level error boundary.
 *
 * This site had none, so an exception anywhere under the root layout fell
 * through to Next's own fallback: one line of unstyled black-on-white reading
 * "Application error: a client-side exception has occurred (see the browser
 * console for more information)". Every project Grit generates ships an
 * error.tsx and a global-error.tsx; the site that generates them did not.
 *
 * The specific failure that is worth handling, rather than only decorating, is
 * the one below.
 */

/*
 * A stale chunk, which is what a deploy does to an open tab.
 *
 * Next splits the app into hashed chunks. A deploy replaces them, so a tab
 * that has been open across one asks for a filename that no longer exists,
 * gets a 404 and throws. Nothing is wrong with the app, and reset() cannot
 * help: it re-renders the same tree, which asks for the same missing file.
 * The fix is a hard reload, which fetches the new HTML and the new names.
 *
 * Matched on the message rather than the type, because the error has crossed
 * a minifier and the class name has not survived: webpack says "Loading chunk
 * N failed", turbopack and native ESM say "Failed to fetch dynamically
 * imported module", and Safari says "Importing a module script failed".
 */
function isStaleChunk(error: Error): boolean {
  const text = `${error.name} ${error.message}`.toLowerCase()
  return (
    text.includes('chunkloaderror') ||
    text.includes('loading chunk') ||
    text.includes('loading css chunk') ||
    text.includes('failed to fetch dynamically imported module') ||
    text.includes('importing a module script failed')
  )
}

// One reload, ever, per tab. An unguarded auto-reload turns any error this
// misidentifies into an infinite refresh, which is worse than the error.
const RELOAD_GUARD = 'grit:reloaded-after-stale-chunk'

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  const [stale, setStale] = useState(false)

  useEffect(() => {
    // Keep this: the console is the only place the real stack exists once the
    // boundary has swallowed it, and the fallback text tells people to look.
    console.error(error)

    if (!isStaleChunk(error)) return
    setStale(true)

    let reloaded = false
    try {
      reloaded = window.sessionStorage.getItem(RELOAD_GUARD) === '1'
      window.sessionStorage.setItem(RELOAD_GUARD, '1')
    } catch {
      // Private mode, or storage disabled. Without a guard an automatic reload
      // could loop, so do not start one: the button below still works.
      return
    }
    if (!reloaded) window.location.reload()
  }, [error])

  return (
    <div className="relative min-h-screen bg-background isolate">
      <GridFrame />

      <main className="max-w-2xl mx-auto px-6 py-24 md:py-32">
        <div className="text-center">
          <span className="tag-mono text-accent mb-5 block">
            {stale ? 'New version' : 'Error'}
          </span>

          <h1 className="font-display text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-4">
            {stale ? 'This page was built by an older version' : 'Something went wrong on this page'}
          </h1>

          <p className="text-lg text-muted-foreground leading-relaxed">
            {stale ? (
              <>
                The site was deployed while this tab was open, so part of it is
                no longer where this page expects. Reloading fixes it, and
                should have happened already.
              </>
            ) : (
              <>
                The page failed while rendering in your browser. It is not
                something you did, and nothing you were reading is lost.
              </>
            )}
          </p>

          {error.digest && (
            <p className="mt-6 font-mono text-xs text-muted-foreground/60">
              Reference: {error.digest}
            </p>
          )}

          <div className="mt-10 flex flex-wrap items-center justify-center gap-3">
            <button
              type="button"
              onClick={() => window.location.reload()}
              className="group inline-flex h-11 items-center gap-2 rounded-full bg-primary px-6 text-sm font-semibold text-primary-foreground transition-colors hover:bg-fg/85 glow-primary-sm"
            >
              <RefreshCw className="h-3.5 w-3.5" />
              Reload the page
            </button>

            {/* Only offered when it can work. After a deploy, re-rendering the
                same tree asks for the same missing file, and a button that
                does nothing twice is worse than one that is not there. */}
            {!stale && (
              <button
                type="button"
                onClick={reset}
                className="inline-flex h-11 items-center gap-2 rounded-full border border-border bg-card/60 px-6 text-sm font-medium text-foreground transition-colors hover:bg-surface-2"
              >
                <RotateCcw className="h-3.5 w-3.5" />
                Try again
              </button>
            )}

            <Link
              href="/"
              className="inline-flex h-11 items-center gap-2 rounded-full border border-border bg-card/60 px-6 text-sm font-medium text-foreground transition-colors hover:bg-surface-2"
            >
              <Home className="h-3.5 w-3.5" />
              Home
            </Link>
          </div>

          <div className="mt-12 flex flex-wrap items-center justify-center gap-x-6 gap-y-2 text-sm">
            <Link
              href="/docs"
              className="inline-flex items-center gap-1.5 text-muted-foreground transition-colors hover:text-accent"
            >
              <BookOpen className="h-3.5 w-3.5" />
              Documentation
            </Link>
            <a
              href="https://github.com/MUKE-coder/grit/issues/new"
              target="_blank"
              rel="noreferrer"
              className="inline-flex items-center gap-1.5 text-muted-foreground transition-colors hover:text-accent"
            >
              Report this
              <ArrowRight className="h-3.5 w-3.5" />
            </a>
          </div>
        </div>
      </main>
    </div>
  )
}
