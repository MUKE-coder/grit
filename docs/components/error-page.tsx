import Link from 'next/link'
import type { ReactNode } from 'react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

/**
 * The shell for /e/<slug>: one problem, explained.
 *
 * These are the pages the CLI prints the address of when something fails, so
 * they are read by somebody who is stuck, in a hurry, and possibly annoyed.
 * That shapes them:
 *
 *   - The fix is above the explanation, because most readers want it and
 *     leave. The reasoning is underneath for the ones who stay.
 *   - Short. A page that needs scrolling to reach the command has failed.
 *   - Every slug here matches an entry in internal/clierr, and a test pairs
 *     the two sets. A renamed slug otherwise leaves the CLI printing a 404.
 */

export interface ErrorPageProps {
  /** The slug, which must match the one in internal/clierr. */
  slug: string
  /** What happened, as the thing that happened. */
  title: string
  /** The one-line cause, shown under the title. */
  cause: string
  /** The fix, first thing on the page. */
  fix: ReactNode
  /** Why it happens, and anything else worth knowing. */
  children: ReactNode
  /** How to tell it worked. */
  confirm?: ReactNode
}

export function ErrorPage({ slug, title, cause, fix, children, confirm }: ErrorPageProps) {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-2xl px-6 py-16">
        <nav aria-label="Breadcrumb" className="mb-6 text-xs text-muted-foreground">
          <Link href="/" className="hover:text-foreground">
            Home
          </Link>
          <span className="mx-2 text-muted-foreground/50">/</span>
          <Link href="/e" className="hover:text-foreground">
            Errors
          </Link>
          <span className="mx-2 text-muted-foreground/50">/</span>
          <span className="font-mono">{slug}</span>
        </nav>

        <h1 className="mb-3 font-display text-3xl font-bold tracking-tight md:text-4xl">
          {title}
        </h1>
        <p className="mb-10 text-lg leading-relaxed text-muted-foreground">{cause}</p>

        {/* The fix, first. Most people came for this and will leave with it. */}
        <section className="mb-12 rounded-2xl border border-accent/30 bg-accent/[0.06] p-6">
          <h2 className="mb-4 font-display text-base font-semibold text-foreground">The fix</h2>
          <div className="space-y-4 text-sm leading-relaxed text-muted-foreground">{fix}</div>
        </section>

        <div className="space-y-4 leading-relaxed text-muted-foreground">{children}</div>

        {confirm && (
          <section className="mt-12 border-t border-border/60 pt-8">
            <h2 className="mb-3 font-display text-lg font-semibold tracking-tight">
              How to tell it worked
            </h2>
            <div className="space-y-3 text-sm leading-relaxed text-muted-foreground">
              {confirm}
            </div>
          </section>
        )}

        <section className="mt-12 border-t border-border/60 pt-8">
          <p className="text-sm leading-relaxed text-muted-foreground">
            Still stuck?{' '}
            <Link href="/docs/getting-started/troubleshooting" className="text-accent hover:underline">
              The troubleshooting page
            </Link>{' '}
            covers more, and{' '}
            <Link href="/guides/common-go-docker-errors" className="text-accent hover:underline">
              the errors guide
            </Link>{' '}
            lists every message with its cause. Running{' '}
            <code className="font-mono text-[13px]">grit doctor --env</code> checks this machine
            for all of them at once.
          </p>
        </section>
      </main>
    </div>
  )
}

export function Why({ children }: { children: ReactNode }) {
  return (
    <>
      <h2 className="mb-3 font-display text-lg font-semibold tracking-tight text-foreground">
        Why it happens
      </h2>
      {children}
    </>
  )
}
