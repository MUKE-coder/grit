import Link from 'next/link'
import type { ReactNode } from 'react'
import { Check, Minus } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { StructuredData } from '@/components/structured-data'

/**
 * The shell for /vs/<tool>: a short, high-intent head-to-head.
 *
 * Deliberately not the /i shell. Those pages survey a field and carry a
 * "checked against each project's source on <date>" line, which is a claim
 * about star counts and activity that these pages do not make. A page that
 * borrows furniture making a promise it does not keep is worse than a plain
 * one.
 *
 * The rules these pages hold to, because a comparison page written by one
 * side is read sceptically and should be:
 *
 *   - Only architectural claims. No versions, no prices, no star counts, no
 *     benchmark numbers. Those go stale and make the page wrong rather than
 *     merely dated, and nobody comes back to fix a competitor's figures.
 *   - Every page says plainly when the other tool is the better choice, in
 *     its own section, not as a hedge in a closing sentence.
 *   - No feature table scored by us. A table where one column happens to win
 *     every row is an advert with gridlines.
 */

const BASE = 'https://gritframework.dev'

export interface VersusProps {
  slug: string
  /** The other tool, as it spells itself. */
  other: string
  title: string
  lead: ReactNode
  children: ReactNode
}

export function VersusPage({ slug, other, title, lead, children }: VersusProps) {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <StructuredData
        data={{
          '@context': 'https://schema.org',
          '@type': 'BreadcrumbList',
          itemListElement: [
            { '@type': 'ListItem', position: 1, name: 'Home', item: BASE },
            { '@type': 'ListItem', position: 2, name: 'Comparisons', item: `${BASE}/vs` },
            { '@type': 'ListItem', position: 3, name: `Grit and ${other}`, item: `${BASE}${slug}` },
          ],
        }}
      />

      <main className="mx-auto max-w-3xl px-6 py-16">
        <nav aria-label="Breadcrumb" className="mb-6 text-xs text-muted-foreground">
          <Link href="/" className="hover:text-foreground">
            Home
          </Link>
          <span className="mx-2 text-muted-foreground/50">/</span>
          <Link href="/vs" className="hover:text-foreground">
            Comparisons
          </Link>
        </nav>

        <h1 className="mb-5 font-display text-4xl font-bold tracking-tight md:text-5xl">
          {title}
        </h1>
        <div className="mb-12 text-lg leading-relaxed text-muted-foreground">{lead}</div>

        {children}

        <section className="mt-14 border-t border-border/60 pt-8">
          <p className="text-sm leading-relaxed text-muted-foreground">
            This page is written by the people who make Grit, so read it the way you would read
            any comparison written by one side. It sticks to how the two are built rather than to
            numbers, because numbers go stale and nobody comes back to fix a competitor&apos;s.
            Both projects are open and take an afternoon to try, which is worth more than this
            page.
          </p>
        </section>
      </main>
    </div>
  )
}

export function VsSection({
  id,
  title,
  children,
}: {
  id: string
  title: string
  children: ReactNode
}) {
  return (
    <section id={id} className="mb-12 scroll-mt-24">
      <h2 className="mb-4 font-display text-2xl font-semibold tracking-tight">{title}</h2>
      <div className="space-y-4 leading-relaxed text-muted-foreground">{children}</div>
    </section>
  )
}

/** Two columns of plain statements. No scoring, no ticks against the other side. */
export function SideBySide({
  left,
  right,
  rows,
}: {
  left: string
  right: string
  rows: { aspect: string; left: ReactNode; right: ReactNode }[]
}) {
  return (
    <div className="my-6 overflow-hidden rounded-xl border border-border/60">
      <div className="grid grid-cols-[minmax(0,7rem)_1fr_1fr] gap-px bg-border/40 text-sm">
        <div className="bg-card/60 px-4 py-3" />
        <div className="bg-card/60 px-4 py-3 font-display font-semibold text-foreground">
          {left}
        </div>
        <div className="bg-card/60 px-4 py-3 font-display font-semibold text-foreground">
          {right}
        </div>
        {rows.map((r) => (
          <div key={r.aspect} className="contents">
            <div className="bg-background px-4 py-3 text-xs font-medium text-muted-foreground">
              {r.aspect}
            </div>
            <div className="bg-background px-4 py-3 leading-relaxed text-muted-foreground">
              {r.left}
            </div>
            <div className="bg-background px-4 py-3 leading-relaxed text-muted-foreground">
              {r.right}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

/** When the other one is the right answer. Every page has one of these. */
export function ChooseOther({ name, items }: { name: string; items: ReactNode[] }) {
  return (
    <div className="my-6 rounded-xl border border-border bg-card/40 p-6">
      <p className="mb-3 font-display text-base font-semibold text-foreground">
        Choose {name} when
      </p>
      <ul className="space-y-2.5">
        {items.map((item, i) => (
          <li key={i} className="flex gap-2.5 text-sm leading-relaxed text-muted-foreground">
            <Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-400" aria-hidden="true" />
            <span>{item}</span>
          </li>
        ))}
      </ul>
    </div>
  )
}

export function ChooseGrit({ items }: { items: ReactNode[] }) {
  return (
    <div className="my-6 rounded-xl border border-accent/30 bg-accent/[0.06] p-6">
      <p className="mb-3 font-display text-base font-semibold text-foreground">Choose Grit when</p>
      <ul className="space-y-2.5">
        {items.map((item, i) => (
          <li key={i} className="flex gap-2.5 text-sm leading-relaxed text-muted-foreground">
            <Check className="mt-0.5 h-4 w-4 shrink-0 text-accent" aria-hidden="true" />
            <span>{item}</span>
          </li>
        ))}
      </ul>
    </div>
  )
}

/** Things neither one does, so the page is not only about winning. */
export function NeitherDoes({ items }: { items: ReactNode[] }) {
  return (
    <ul className="my-6 space-y-2.5">
      {items.map((item, i) => (
        <li key={i} className="flex gap-2.5 text-sm leading-relaxed text-muted-foreground">
          <Minus className="mt-0.5 h-4 w-4 shrink-0 text-muted-foreground/50" aria-hidden="true" />
          <span>{item}</span>
        </li>
      ))}
    </ul>
  )
}
