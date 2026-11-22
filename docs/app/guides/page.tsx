import Link from 'next/link'
import { ArrowRight, Clock } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { StructuredData } from '@/components/structured-data'
import { LIVE_GUIDES, PLANNED_GUIDES, type Guide } from '@/config/guides'

export const metadata = {
  // The root layout appends the site name, so the title is the page's own.
  title: 'Guides',
  description:
    'Practical guides to building with Go and React: REST APIs, JWT auth, GORM relationships, Docker, deploys, multi-tenancy, caching and migrations. Each one answers the problem in plain Go first.',
  alternates: { canonical: 'https://gritframework.dev/guides' },
}

function Card({ g }: { g: Guide }) {
  return (
    <li>
      <Link
        href={`/guides/${g.slug}`}
        className="group flex h-full flex-col rounded-2xl border border-border/60 bg-card/30 p-6 transition-colors hover:border-accent/50"
      >
        <h3 className="font-display text-lg font-semibold leading-snug text-foreground">
          {g.title}
        </h3>
        <p className="mt-2 flex-1 text-sm leading-relaxed text-muted-foreground">{g.blurb}</p>
        <p className="mt-4 flex items-center gap-4 text-xs text-muted-foreground/80">
          <span className="inline-flex items-center gap-1.5">
            <Clock className="h-3 w-3" aria-hidden="true" />
            {g.reading}
          </span>
          <span className="inline-flex items-center gap-1 text-accent opacity-0 transition-opacity group-hover:opacity-100">
            Read
            <ArrowRight className="h-3 w-3" aria-hidden="true" />
          </span>
        </p>
      </Link>
    </li>
  )
}

export default function Page() {
  const beginner = LIVE_GUIDES.filter((g) => g.audience === 'beginner')
  const senior = LIVE_GUIDES.filter((g) => g.audience === 'senior')

  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <StructuredData
        data={{
          '@context': 'https://schema.org',
          '@type': 'CollectionPage',
          '@id': 'https://gritframework.dev/guides',
          name: 'Grit guides',
          description: metadata.description,
          hasPart: LIVE_GUIDES.map((g) => ({
            '@type': 'TechArticle',
            name: g.title,
            url: `https://gritframework.dev/guides/${g.slug}`,
            abstract: g.blurb,
          })),
        }}
      />

      <main className="mx-auto max-w-6xl px-6 py-16">
        <h1 className="mb-5 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Guides to building with Go and React
        </h1>
        <p className="mb-16 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          Each of these answers one question properly, in plain Go and React, with code you can run
          whether or not you use Grit. Where Grit does the same thing in one command, that comes at
          the end, as the payoff rather than the pitch.
        </p>

        {beginner.length > 0 && (
          <section className="mb-20">
            <h2 className="mb-2 font-display text-2xl font-semibold tracking-tight">
              Starting out
            </h2>
            <p className="mb-7 max-w-2xl text-sm text-muted-foreground">
              Nothing assumed beyond knowing JavaScript. Every step has something you can check
              before moving on.
            </p>
            <ul className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {beginner.map((g) => (
                <Card key={g.slug} g={g} />
              ))}
            </ul>
          </section>
        )}

        {senior.length > 0 && (
          <section className="mb-20">
            <h2 className="mb-2 font-display text-2xl font-semibold tracking-tight">
              Production concerns
            </h2>
            <p className="mb-7 max-w-2xl text-sm text-muted-foreground">
              The decisions that are expensive to get wrong, with the trade-offs stated rather than
              one answer asserted.
            </p>
            <ul className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {senior.map((g) => (
                <Card key={g.slug} g={g} />
              ))}
            </ul>
          </section>
        )}

        {PLANNED_GUIDES.length > 0 && (
          <section className="mb-12 rounded-2xl border border-border/60 bg-card/20 p-7">
            <h2 className="mb-2 font-display text-lg font-semibold tracking-tight">
              Being written
            </h2>
            <p className="mb-5 max-w-2xl text-sm text-muted-foreground">
              Listed because the order is worth knowing, and because the set is published as it is
              written rather than all at once. Nothing here is a link yet.
            </p>
            <ul className="grid gap-x-8 gap-y-2 text-sm text-muted-foreground sm:grid-cols-2">
              {PLANNED_GUIDES.map((g) => (
                <li key={g.slug} className="flex gap-2">
                  <span className="mt-[9px] h-1 w-1 shrink-0 rounded-full bg-muted-foreground/40" />
                  <span>{g.title}</span>
                </li>
              ))}
            </ul>
          </section>
        )}

        <section className="border-t border-border/60 pt-10">
          <h2 className="mb-4 font-display text-xl font-semibold tracking-tight">
            Looking for the reference instead?
          </h2>
          <p className="max-w-2xl text-sm leading-relaxed text-muted-foreground">
            These pages teach a topic. For how a specific part of Grit behaves, with every flag and
            field type, read <Link href="/docs" className="text-accent hover:underline">the
            documentation</Link>. For the full applications built with it, start with{' '}
            <Link href="/docs/tutorials" className="text-accent hover:underline">the tutorials</Link>
            .
          </p>
        </section>
      </main>
    </div>
  )
}
