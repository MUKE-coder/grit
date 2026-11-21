import Link from 'next/link'
import type { ReactNode } from 'react'
import { AlertTriangle, Check, Lightbulb } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { StructuredData } from '@/components/structured-data'
import { guide } from '@/config/guides'

// The shell for /guides: problem-first articles for somebody searching for
// their problem rather than for Grit.
//
// A sibling of components/informational.tsx rather than a reuse of it. That
// shell's breadcrumb points at /i and its dateline reads "star counts and
// activity dates checked against each project's source", which is true of a
// comparison page and false of a tutorial. Two shells with honest furniture
// beat one shell with a prop that makes it lie.
//
// The order of a guide is fixed, and it is the order that earns the ranking:
// answer the question in plain Go and React first, then show the Grit
// one-liner as the payoff. A page that opens with the product is an advert,
// and is read as one.
//
// Everything is prose in each page's own file rather than a config the pages
// render. Two indexers walk this site for text, and content that exists only
// as data is invisible to both.

const BASE = 'https://gritframework.dev'

export interface GuidePageProps {
  /** The slug, without the /guides prefix. Its entry carries the rest. */
  slug: string
  /** The answer, before the fold. A reader who stops here should have it. */
  lead: ReactNode
  children: ReactNode
}

export function GuidePage({ slug, lead, children }: GuidePageProps) {
  const entry = guide(slug)
  const path = `/guides/${slug}`
  const title = entry.title
  const eyebrow = entry.audience === 'senior' ? 'For senior developers' : 'For beginners'
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
            { '@type': 'ListItem', position: 2, name: 'Guides', item: `${BASE}/guides` },
            { '@type': 'ListItem', position: 3, name: title, item: `${BASE}${path}` },
          ],
        }}
      />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <nav aria-label="Breadcrumb" className="mb-6 text-xs text-muted-foreground">
          <Link href="/" className="hover:text-foreground">
            Home
          </Link>
          <span className="mx-2 text-muted-foreground/50">/</span>
          <Link href="/guides" className="hover:text-foreground">
            Guides
          </Link>
        </nav>

        <span className="font-mono text-xs uppercase tracking-wider text-accent">{eyebrow}</span>
        <h1 className="mb-5 mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
          {title}
        </h1>
        <div className="mb-4 max-w-3xl text-lg leading-relaxed text-muted-foreground">{lead}</div>
        <p className="mb-12 text-xs text-muted-foreground/70">
          {entry.reading} read. Commands and output last run against Grit {entry.updated}.
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

/**
 * What the reader will have when they finish, at the top.
 *
 * Beginners abandon a tutorial that does not say where it ends, and it is the
 * single most requested thing in the review this section set came from.
 */
export function Outcome({ children }: { children: ReactNode }) {
  return (
    <div className="mb-12 rounded-2xl border border-accent/25 bg-accent/[0.06] p-6">
      <p className="mb-3 flex items-center gap-2 font-display text-sm font-semibold text-foreground">
        <Check className="h-4 w-4 text-accent" aria-hidden="true" />
        What you will have at the end
      </p>
      <div className="space-y-2 text-sm leading-relaxed text-muted-foreground">{children}</div>
    </div>
  )
}

/** An aside worth reading but not worth interrupting the flow for. */
export function Note({ title, children }: { title: string; children: ReactNode }) {
  return (
    <div className="my-6 flex gap-3 rounded-xl border border-border bg-card/30 p-5">
      <Lightbulb className="mt-0.5 h-4 w-4 shrink-0 text-accent" aria-hidden="true" />
      <div className="space-y-2 text-sm leading-relaxed text-muted-foreground">
        <p className="font-medium text-foreground">{title}</p>
        {children}
      </div>
    </div>
  )
}

/** A thing that will go wrong, with the fix. */
export function Warn({ title, children }: { title: string; children: ReactNode }) {
  return (
    <div className="my-6 flex gap-3 rounded-xl border border-amber-500/30 bg-amber-500/[0.06] p-5">
      <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0 text-amber-500" aria-hidden="true" />
      <div className="space-y-2 text-sm leading-relaxed text-muted-foreground">
        <p className="font-medium text-foreground">{title}</p>
        {children}
      </div>
    </div>
  )
}

/**
 * One error, with the exact text as the heading.
 *
 * The heading is the error somebody pastes into a search box, so it is the
 * thing that has to match. Wrapping it in prose ("you may see an error
 * about...") is how a troubleshooting page becomes unfindable.
 */
export function Mistake({ error, children }: { error: string; children: ReactNode }) {
  return (
    <div className="mb-8">
      <h3 className="mb-2 mt-8 font-mono text-[13px] font-semibold leading-snug text-foreground">
        {error}
      </h3>
      <div className="space-y-3 leading-relaxed text-muted-foreground">{children}</div>
    </div>
  )
}

/** The numbered body of a tutorial. */
export function Step({
  n,
  title,
  children,
}: {
  n: number
  title: string
  children: ReactNode
}) {
  return (
    <section className="mb-12 scroll-mt-24" id={`step-${n}`}>
      <h2 className="mb-4 flex items-baseline gap-3 font-display text-2xl font-semibold tracking-tight">
        <span className="font-mono text-base text-accent">{String(n).padStart(2, '0')}</span>
        {title}
      </h2>
      <div className="space-y-4 leading-relaxed text-muted-foreground">{children}</div>
    </section>
  )
}

/**
 * The checklist a senior reads before putting it in production.
 *
 * Deliberately not a summary of the page. Every line is something that is not
 * true by default and that somebody has to decide.
 */
export function Checklist({ items }: { items: ReactNode[] }) {
  return (
    <ul className="my-6 space-y-3">
      {items.map((item, i) => (
        <li key={i} className="flex gap-3 text-sm leading-relaxed text-muted-foreground">
          <span
            className="mt-[7px] h-1.5 w-1.5 shrink-0 rounded-full bg-accent"
            aria-hidden="true"
          />
          <span>{item}</span>
        </li>
      ))}
    </ul>
  )
}

export interface QA {
  q: string
  a: ReactNode
  /** The answer as plain text, for the schema. It must match what is rendered. */
  text: string
}

/**
 * FAQ, with schema over answers a reader can actually see.
 *
 * A FAQPage block describing answers that are not on the page is the kind of
 * markup that earns a manual action, so the text field is the same answer in
 * plain prose rather than a keyword-stuffed twin.
 */
export function Faq({ items, slug }: { items: QA[]; slug: string }) {
  return (
    <section id="faq" className="mb-14 scroll-mt-24">
      <StructuredData
        data={{
          '@context': 'https://schema.org',
          '@type': 'FAQPage',
          '@id': `${BASE}/guides/${slug}#faq`,
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
      <h2 className="mb-5 font-display text-2xl font-semibold tracking-tight">Next steps</h2>
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
