import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, Check, X, Minus } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

export const metadata: Metadata = {
  title: 'Grit vs Django — Go + React vs Python’s batteries-included classic',
  description:
    'Django gives you batteries and a famous auto-admin on Python with server-rendered templates; Grit gives you batteries and a generated modern React admin on Go, typed end-to-end, with clients for web, mobile, and desktop. When to pick each, honestly compared.',
  alternates: { canonical: 'https://gritframework.dev/compare/django' },
}

const rows: [string, string, 'yes' | 'no' | 'partial', string][] = [
  ['Backend language', 'Go (Gin + GORM)', 'yes', 'Python (Django ORM)'],
  ['Admin panel', 'Generated React, yours to style', 'yes', 'Auto, server-rendered, ORM-coupled'],
  ['Typed API → client', 'Generated Go → TS + Zod', 'yes', 'Manual, or DRF serializers'],
  ['Auth + RBAC', 'Built in (roles + permissions)', 'yes', 'Built in (auth + permissions)'],
  ['Background jobs / cron', 'Built in (asynq)', 'yes', 'Celery / add-on'],
  ['File storage, email, cache', 'Built in (S3/R2, Resend, Redis)', 'yes', 'Built in / configurable'],
  ['Mobile + desktop clients', 'Generated (Expo, Wails)', 'yes', 'Not in scope'],
  ['Database backups', 'Built in, scheduled + restore', 'yes', 'Roll your own'],
  ['Ecosystem + maturity', 'Younger, smaller', 'partial', 'Vast, battle-tested'],
  ['Data / ML gravity', 'Not the focus', 'partial', 'Python-native, first-class'],
]

export default function CompareDjangoPage() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <span className="font-mono text-xs uppercase tracking-wider text-primary">Compare</span>
        <h1 className="mb-4 mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Grit vs Django
        </h1>
        <p className="mb-10 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          These two share a philosophy: batteries included. Django was the landmark that made it
          famous — an ORM, migrations, and the auto-generated admin that&apos;s launched countless
          internal tools. Grit takes the same &ldquo;everything in the box&rdquo; idea to Go and a
          modern React frontend. Django gives you batteries and an auto-admin on Python; Grit gives
          you batteries and a generated React admin on Go, typed end-to-end.
        </p>

        {/* Short version */}
        <div className="mb-12 grid gap-4 sm:grid-cols-2">
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">The short version</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Reach for <strong>Django</strong> when Python is home — your team, your data and ML
              stack, or the vast ecosystem you want to lean on — and the famous auto-admin covers
              your internal tooling out of the box. Reach for <strong>Grit</strong> when you want
              Go&apos;s speed and end-to-end types, a modern React admin you own and style, and one
              API feeding web, mobile, and desktop clients.
            </p>
          </div>
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">Same idea, different stack</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Django&apos;s admin is <em>automatic</em> but server-rendered and coupled to your
              models — powerful and zero-config, if a little dated to look at. Grit&apos;s admin is
              <em> generated</em> React code you own: run <code>grit generate resource</code> and you
              get a model, service, handler, Zod schema, TS types, React Query hooks, and an admin
              page in one command — then style it however you like.
            </p>
          </div>
        </div>

        {/* Table */}
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">Side by side</h2>
        <div className="mb-12 overflow-x-auto rounded-2xl border border-border">
          <table className="w-full min-w-[640px] text-sm">
            <thead>
              <tr className="border-b border-border bg-card/40 text-left">
                <th className="px-4 py-3 font-medium text-foreground/70"> </th>
                <th className="px-4 py-3 font-semibold text-foreground">Grit</th>
                <th className="px-4 py-3 font-semibold text-foreground/70">Django</th>
              </tr>
            </thead>
            <tbody>
              {rows.map(([label, grit, mark, other]) => (
                <tr key={label} className="border-b border-border/50 last:border-b-0">
                  <td className="px-4 py-3 font-medium text-foreground">{label}</td>
                  <td className="px-4 py-3 text-muted-foreground">
                    <span className="inline-flex items-center gap-1.5">
                      {mark === 'yes' && <Check className="h-4 w-4 shrink-0 text-emerald-500" />}
                      {mark === 'partial' && <Minus className="h-4 w-4 shrink-0 text-amber-500" />}
                      {grit}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-muted-foreground/80">{other}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* When each */}
        <div className="mb-12 grid gap-6 md:grid-cols-2">
          <div>
            <h2 className="mb-3 text-xl font-semibold tracking-tight">Choose Grit when</h2>
            <ul className="space-y-2 text-sm leading-relaxed text-muted-foreground">
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want a modern React admin you own and style, not a server-rendered one.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want end-to-end types: Go models generate TypeScript and Zod.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> More than one client consumes the API (web + mobile + desktop + admin).</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want Go&apos;s performance and a single self-hosted binary.</li>
            </ul>
          </div>
          <div>
            <h2 className="mb-3 text-xl font-semibold tracking-tight">Choose Django when</h2>
            <ul className="space-y-2 text-sm leading-relaxed text-muted-foreground">
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> Python is your team&apos;s language and you want to stay in one.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You lean on Python&apos;s data and ML gravity, or a data-adjacent team.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You want the deepest, most mature ecosystem and documentation.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> The zero-config auto-admin is exactly enough for your internal tools.</li>
            </ul>
          </div>
        </div>

        {/* Honest note */}
        <div className="mb-12 rounded-2xl border border-border bg-card/40 p-6">
          <h2 className="mb-2 font-semibold text-foreground">The honest trade-off</h2>
          <p className="text-sm leading-relaxed text-muted-foreground">
            Django is more mature, Python-native, and one language — its ecosystem and community are
            hard to match, and its auto-admin is genuinely a landmark. Grit is younger with a smaller
            ecosystem, opinionated about its stack, and asks you to work in two languages (Go + TS).
            In exchange you get Go&apos;s speed, types that flow all the way to the client, a React
            admin you own, and web/mobile/desktop from one API. If Python and its ecosystem are your
            gravity, Django is the safer, deeper choice. If you want a typed Go backend and a modern
            admin across many clients, that&apos;s where Grit earns its keep.
          </p>
        </div>

        {/* CTA */}
        <div className="flex flex-wrap items-center gap-4">
          <Link
            href="/docs/getting-started/quick-start"
            className="inline-flex items-center gap-2 rounded-xl bg-primary px-5 py-3 font-medium text-primary-foreground transition-opacity hover:opacity-90"
          >
            Try Grit in 5 minutes
            <ArrowRight className="h-4 w-4" />
          </Link>
          <Link href="/compare" className="text-sm text-muted-foreground hover:text-foreground">
            See all comparisons →
          </Link>
        </div>
      </main>
    </div>
  )
}
