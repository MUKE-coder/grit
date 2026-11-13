import Link from 'next/link'
import type { ReactNode } from 'react'
import { Check, Minus, X } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { StructuredData } from '@/components/structured-data'

// The shell and the small parts the /i/ pages are built from.
//
// These pages exist to be found by somebody searching for a problem rather
// than for Grit, so they are written as the answer to that question first.
// Grit gets one section, near the end, where it is genuinely one of the
// options. A page that mentions the product in every paragraph reads as an
// advert and is treated as one, by readers and by search engines.
//
// Everything is real prose in each page's own file rather than a config the
// pages render. Two different indexers walk this site's pages for their text,
// and content that only exists as data is invisible to both.

const BASE = 'https://gritframework.dev'

export interface InfoPageProps {
  /** The URL path, for the canonical and the breadcrumb. */
  slug: string
  eyebrow: string
  title: string
  /** The opening paragraph, which should answer the query before the fold. */
  lead: ReactNode
  /** When the facts on this page were last checked against their sources. */
  checked: string
  children: ReactNode
}

export function InfoPage({ slug, eyebrow, title, lead, checked, children }: InfoPageProps) {
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
            { '@type': 'ListItem', position: 2, name: 'Guides', item: `${BASE}/i` },
            { '@type': 'ListItem', position: 3, name: title, item: `${BASE}${slug}` },
          ],
        }}
      />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <nav aria-label="Breadcrumb" className="mb-6 text-xs text-muted-foreground">
          <Link href="/" className="hover:text-foreground">
            Home
          </Link>
          <span className="mx-2 text-muted-foreground/50">/</span>
          <Link href="/i" className="hover:text-foreground">
            Guides
          </Link>
        </nav>

        <span className="font-mono text-xs uppercase tracking-wider text-accent">{eyebrow}</span>
        <h1 className="mb-5 mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
          {title}
        </h1>
        <div className="mb-4 max-w-3xl text-lg leading-relaxed text-muted-foreground">{lead}</div>
        <p className="mb-12 text-xs text-muted-foreground/70">
          Star counts and activity dates checked against each project&apos;s source on {checked}.
          Numbers move; the reasoning does not.
        </p>

        {children}
      </main>
    </div>
  )
}

/** A top-level section with its own heading. */
export function Section({
  id,
  title,
  children,
}: {
  id: string
  title: string
  children: ReactNode
}) {
  return (
    <section id={id} className="mb-14 scroll-mt-24">
      <h2 className="mb-5 font-display text-2xl font-semibold tracking-tight">{title}</h2>
      <div className="space-y-4 leading-relaxed text-muted-foreground">{children}</div>
    </section>
  )
}

export function H3({ children }: { children: ReactNode }) {
  return <h3 className="mb-3 mt-8 font-display text-lg font-semibold text-foreground">{children}</h3>
}

/** One option in a comparison: what it is, what it is good at, what it is not. */
export function Option({
  n,
  name,
  meta,
  what,
  pros,
  cons,
  bestFor,
  href,
}: {
  n: number
  name: string
  /** Stars, last activity, licence: the facts a reader checks before trusting a list. */
  meta: string
  what: ReactNode
  pros: string[]
  cons: string[]
  bestFor: string
  href?: string
}) {
  return (
    <article className="mb-8 rounded-xl border border-border/60 bg-card/30 p-6">
      <h3 className="font-display text-xl font-semibold text-foreground">
        {n}. {href ? (
          <a href={href} target="_blank" rel="noreferrer" className="hover:text-accent">
            {name}
          </a>
        ) : (
          name
        )}
      </h3>
      <p className="mb-4 mt-1 font-mono text-xs text-muted-foreground/80">{meta}</p>
      <div className="mb-5 space-y-3 leading-relaxed text-muted-foreground">{what}</div>

      <div className="grid gap-5 sm:grid-cols-2">
        <div>
          <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground/70">
            Strengths
          </p>
          <ul className="space-y-1.5 text-sm text-muted-foreground">
            {pros.map((p) => (
              <li key={p} className="flex gap-2">
                <Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" aria-hidden />
                <span>{p}</span>
              </li>
            ))}
          </ul>
        </div>
        <div>
          <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground/70">
            Trade-offs
          </p>
          <ul className="space-y-1.5 text-sm text-muted-foreground">
            {cons.map((c) => (
              <li key={c} className="flex gap-2">
                <X className="mt-0.5 h-4 w-4 shrink-0 text-rose-500" aria-hidden />
                <span>{c}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>

      <p className="mt-5 border-t border-border/50 pt-4 text-sm text-foreground">
        <span className="font-semibold">Best for: </span>
        <span className="text-muted-foreground">{bestFor}</span>
      </p>
    </article>
  )
}

export type Cell = string | 'yes' | 'no' | 'partial'

/** The at-a-glance table every one of these pages opens its comparison with. */
export function CompareTable({
  columns,
  rows,
  caption,
}: {
  columns: string[]
  rows: Cell[][]
  caption?: string
}) {
  return (
    <div className="my-6 overflow-x-auto rounded-xl border border-border/60">
      <table className="w-full min-w-[640px] text-left text-sm">
        {caption && <caption className="sr-only">{caption}</caption>}
        <thead className="bg-card/50">
          <tr>
            {columns.map((c) => (
              <th key={c} scope="col" className="px-4 py-3 font-semibold text-foreground">
                {c}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={i} className="border-t border-border/50">
              {row.map((cell, j) => (
                <td key={j} className="px-4 py-3 align-top text-muted-foreground">
                  {j === 0 ? <span className="font-medium text-foreground">{cell}</span> : mark(cell)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function mark(cell: Cell) {
  if (cell === 'yes') return <Check className="h-4 w-4 text-emerald-500" aria-label="Yes" />
  if (cell === 'no') return <X className="h-4 w-4 text-rose-500" aria-label="No" />
  if (cell === 'partial') return <Minus className="h-4 w-4 text-amber-500" aria-label="Partly" />
  return cell
}

/** The section where Grit is one of the answers, said once and plainly. */
export function WhereGritFits({ children }: { children: ReactNode }) {
  return (
    <section id="grit" className="mb-14 scroll-mt-24 rounded-xl border border-accent/30 bg-accent/5 p-6">
      <h2 className="mb-4 font-display text-2xl font-semibold tracking-tight">Where Grit fits</h2>
      <div className="space-y-4 leading-relaxed text-muted-foreground">{children}</div>
      <div className="mt-6 flex flex-wrap gap-3">
        <Link
          href="/docs/getting-started/quick-start"
          className="rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white hover:opacity-90"
        >
          Start a project
        </Link>
        <Link
          href="/compare"
          className="rounded-lg border border-border px-4 py-2 text-sm font-semibold text-foreground hover:bg-card"
        >
          Read the honest comparisons
        </Link>
      </div>
    </section>
  )
}

export interface QA {
  q: string
  a: ReactNode
  /** The same answer as plain text, for the structured data. */
  text: string
}

/**
 * The questions, visible on the page and described once in JSON-LD.
 *
 * The schema describes what a reader can see, which is the only form of it
 * worth having: a FAQPage block over answers that are not on the page is the
 * kind of markup that gets a site a manual action.
 */
export function Faq({ items, slug }: { items: QA[]; slug: string }) {
  return (
    <section id="faq" className="mb-14 scroll-mt-24">
      <StructuredData
        data={{
          '@context': 'https://schema.org',
          '@type': 'FAQPage',
          '@id': `${BASE}${slug}#faq`,
          mainEntity: items.map((i) => ({
            '@type': 'Question',
            name: i.q,
            acceptedAnswer: { '@type': 'Answer', text: i.text },
          })),
        }}
      />
      <h2 className="mb-5 font-display text-2xl font-semibold tracking-tight">
        Frequently asked questions
      </h2>
      <dl className="space-y-6">
        {items.map((i) => (
          <div key={i.q}>
            <dt className="mb-2 font-display text-base font-semibold text-foreground">{i.q}</dt>
            <dd className="leading-relaxed text-muted-foreground">{i.a}</dd>
          </div>
        ))}
      </dl>
    </section>
  )
}

export function Related({ links }: { links: { href: string; label: string; blurb: string }[] }) {
  return (
    <section className="mb-8">
      <h2 className="mb-5 font-display text-2xl font-semibold tracking-tight">Keep reading</h2>
      <ul className="grid gap-3 sm:grid-cols-2">
        {links.map((l) => (
          <li key={l.href}>
            <Link
              href={l.href}
              className="block h-full rounded-xl border border-border/60 bg-card/30 p-4 transition-colors hover:border-accent/50"
            >
              <span className="block font-medium text-foreground">{l.label}</span>
              <span className="mt-1 block text-sm text-muted-foreground">{l.blurb}</span>
            </Link>
          </li>
        ))}
      </ul>
    </section>
  )
}
