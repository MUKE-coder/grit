import Link from 'next/link'
import type { Metadata } from 'next'
import { Heart, Star, ArrowRight } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { SponsorCard, FoundingSponsorCTA } from '@/components/sponsors'
import {
  SPONSORS_LAST_SYNC,
  SPONSOR_TIERS,
  activeSponsors,
  specialSponsors,
  sponsorStats,
} from '@/config/sponsors'

export const metadata: Metadata = {
  title: 'Sponsors — Grit',
  description:
    'The companies and developers funding Grit — the Go + React full-stack framework.',
  alternates: { canonical: 'https://gritframework.dev/sponsors' },
}

function Stat({
  label,
  value,
  hint,
}: {
  label: string
  value: string
  hint: string
}) {
  return (
    <div className="rounded-xl border border-border/40 bg-card/50 p-5">
      <p className="font-mono text-[10px] uppercase tracking-wider text-muted-foreground/60">
        {label}
      </p>
      <p className="mt-3 font-display text-3xl font-bold text-foreground">{value}</p>
      <p className="mt-1 text-xs text-muted-foreground/60">{hint}</p>
    </div>
  )
}

export default function SponsorsPage() {
  const stats = sponsorStats()
  const special = specialSponsors()
  const everyone = activeSponsors()
  const rest = everyone.filter((s) => !special.includes(s))

  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />

      {/* Header */}
      <section className="relative overflow-hidden border-b border-border/30">
        <div className="pointer-events-none absolute inset-0 -z-10">
          <div className="absolute left-1/2 top-0 h-[360px] w-[680px] -translate-x-1/2 rounded-full bg-primary/[0.06] blur-[120px]" />
          <div className="absolute inset-0 bg-grit-grid-sm mask-fade-y opacity-[0.35]" />
        </div>
        <div className="container max-w-screen-xl px-6 py-16 md:py-20">
          <div className="flex flex-col gap-6 md:flex-row md:items-end md:justify-between">
            <div>
              <div className="mb-4 flex items-center gap-2">
                <Heart className="h-5 w-5 text-primary" />
                <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
                  Sponsors
                </span>
              </div>
              <h1 className="font-display text-4xl font-bold tracking-tight md:text-5xl">
                <span className="gradient-text">Who funds Grit</span>
              </h1>
              <p className="mt-4 max-w-xl text-lg text-muted-foreground">
                The companies and developers paying for the features, docs and
                releases everyone else gets for free.
              </p>
            </div>
            <div className="flex flex-col items-start gap-3 md:items-end">
              <span className="font-mono text-[10px] uppercase tracking-wider text-muted-foreground/50">
                Last sync: {SPONSORS_LAST_SYNC}
              </span>
              <Link
                href="/sponsor"
                className="inline-flex items-center gap-2 rounded-lg border border-primary/40 bg-primary/10 px-4 py-2 text-sm font-semibold text-primary transition-colors hover:bg-primary/20"
              >
                <Heart className="h-4 w-4" />
                Become a sponsor
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Stats — only meaningful once someone has sponsored. */}
      {stats.hasAny && (
        <section className="border-b border-border/30">
          <div className="container max-w-screen-xl px-6 py-10">
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
              <Stat
                label="Lifetime funding"
                value={`$${stats.lifetimeFunding.toLocaleString()}`}
                hint="All-time processed"
              />
              <Stat
                label="Monthly recurring"
                value={`$${stats.monthlyRecurring.toLocaleString()}`}
                hint="Per month right now"
              />
              <Stat
                label="Active sponsors"
                value={String(stats.activeCount)}
                hint={`${stats.specialCount} special`}
              />
              <Stat
                label="All-time sponsors"
                value={String(stats.allTimeCount)}
                hint={`Including ${stats.allTimeCount - stats.activeCount} past`}
              />
            </div>
          </div>
        </section>
      )}

      {/* Roster */}
      <section>
        <div className="container max-w-screen-xl px-6 py-16">
          {!stats.hasAny ? (
            <div className="mx-auto max-w-3xl">
              <FoundingSponsorCTA />

              {/* With an empty roster, show what sponsorship buys instead of zeros. */}
              <div className="mt-12">
                <h2 className="mb-6 text-center font-mono text-xs uppercase tracking-wider text-muted-foreground/60">
                  Sponsorship tiers
                </h2>
                <div className="grid gap-3 sm:grid-cols-2">
                  {SPONSOR_TIERS.map((tier) => (
                    <Link
                      key={tier.id}
                      href="/sponsor"
                      className="card-grit flex items-center justify-between rounded-xl border border-border/40 bg-card/50 px-5 py-4"
                    >
                      <div>
                        <p className="font-display text-sm font-bold text-foreground">
                          {tier.name}
                        </p>
                        <p className="mt-0.5 text-xs text-muted-foreground">
                          {tier.blurb}
                        </p>
                      </div>
                      <span className="shrink-0 font-mono text-sm font-semibold text-primary">
                        ${tier.price}/mo
                      </span>
                    </Link>
                  ))}
                </div>
              </div>
            </div>
          ) : (
            <div className="space-y-14">
              {special.length > 0 && (
                <div>
                  <div className="mb-6 flex items-center gap-3">
                    <Star className="h-4 w-4 shrink-0 fill-current text-amber-400" />
                    <h2 className="font-mono text-xs uppercase tracking-wider text-foreground">
                      Special sponsors
                    </h2>
                    <div className="h-px flex-1 bg-border/40" />
                    <span className="font-mono text-[10px] uppercase tracking-wider text-muted-foreground/50">
                      [{special.length} {special.length === 1 ? 'record' : 'records'}]
                    </span>
                  </div>
                  <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                    {special.map((sponsor) => (
                      <SponsorCard key={sponsor.name} sponsor={sponsor} />
                    ))}
                  </div>
                </div>
              )}

              {rest.length > 0 && (
                <div>
                  <div className="mb-6 flex items-center gap-3">
                    <h2 className="font-mono text-xs uppercase tracking-wider text-foreground">
                      Sponsors
                    </h2>
                    <div className="h-px flex-1 bg-border/40" />
                    <span className="font-mono text-[10px] uppercase tracking-wider text-muted-foreground/50">
                      [{rest.length} {rest.length === 1 ? 'record' : 'records'}]
                    </span>
                  </div>
                  <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                    {rest.map((sponsor) => (
                      <SponsorCard key={sponsor.name} sponsor={sponsor} />
                    ))}
                  </div>
                </div>
              )}

              <div className="rounded-xl border border-dashed border-border/50 px-6 py-10 text-center">
                <h3 className="font-display text-lg font-bold text-foreground">
                  Your logo could be here
                </h3>
                <p className="mx-auto mt-2 max-w-md text-sm text-muted-foreground">
                  Sponsors get their name in front of every developer who builds with
                  Grit — on this site, in the README, and inside the CLI.
                </p>
                <Link
                  href="/sponsor"
                  className="mt-5 inline-flex items-center gap-2 text-sm font-semibold text-primary hover:underline"
                >
                  See the tiers
                  <ArrowRight className="h-4 w-4" />
                </Link>
              </div>
            </div>
          )}
        </div>
      </section>

      {/* Badge kit — the "give something back" half of sponsorship. */}
      <section className="border-t border-border/30 bg-accent/20">
        <div className="container max-w-screen-xl px-6 py-16">
          <div className="mx-auto max-w-3xl">
            <div className="mb-8 text-center">
              <h2 className="font-display text-2xl font-bold tracking-tight md:text-3xl">
                Badges
              </h2>
              <p className="mx-auto mt-3 max-w-xl text-sm text-muted-foreground">
                Shipping with Grit? Add a badge to your README. Sponsors are welcome to
                use it anywhere.
              </p>
            </div>
            <div className="grid gap-4 sm:grid-cols-2">
              {[
                {
                  label: 'Built with Grit',
                  file: 'built-with-grit.svg',
                  alt: 'Built with Grit',
                },
                {
                  label: 'Powered by Grit',
                  file: 'powered-by-grit.svg',
                  alt: 'Powered by Grit',
                },
              ].map((badge) => (
                <div
                  key={badge.file}
                  className="rounded-xl border border-border/40 bg-card/50 p-5"
                >
                  <div className="mb-4 flex items-center justify-center rounded-lg border border-border/30 bg-background/50 py-6">
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img
                      src={`/badge/${badge.file}`}
                      alt={badge.alt}
                      width={168}
                      height={28}
                    />
                  </div>
                  <p className="mb-2 font-mono text-[10px] uppercase tracking-wider text-muted-foreground/60">
                    Markdown
                  </p>
                  <pre className="overflow-x-auto rounded-lg border border-border/30 bg-background/70 p-3">
                    <code className="font-mono text-[11px] text-foreground/80">
                      {`[![${badge.alt}](https://gritframework.dev/badge/${badge.file})](https://gritframework.dev)`}
                    </code>
                  </pre>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      <footer className="border-t border-border/30 px-6 py-8">
        <div className="container flex max-w-screen-xl flex-col items-center justify-between gap-4 text-sm text-muted-foreground/50 sm:flex-row">
          <span>Grit Framework — Go + React. Built with Grit.</span>
          <div className="flex items-center gap-4">
            <Link href="/docs" className="transition-colors hover:text-foreground">
              Docs
            </Link>
            <Link href="/sponsor" className="transition-colors hover:text-foreground">
              Become a sponsor
            </Link>
            <Link
              href="https://github.com/MUKE-coder/grit"
              target="_blank"
              rel="noreferrer"
              className="transition-colors hover:text-foreground"
            >
              GitHub
            </Link>
          </div>
        </div>
      </footer>
    </div>
  )
}
