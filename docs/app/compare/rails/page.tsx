import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, Check, X, Minus } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

export const metadata: Metadata = {
  title: 'Grit vs Rails — Go + React vs Ruby’s batteries-included original',
  description:
    'Rails brought batteries-included, convention-over-configuration to Ruby with server-rendered views; Grit brings the same philosophy to Go with a typed React admin and generated web/mobile/desktop clients. When to pick each, honestly compared.',
  alternates: { canonical: 'https://gritframework.dev/compare/rails' },
}

const rows: [string, string, 'yes' | 'no' | 'partial', string][] = [
  ['Language', 'Go + TypeScript', 'yes', 'Ruby (one language)'],
  ['Philosophy', 'Batteries-included, convention', 'yes', 'Batteries-included, convention'],
  ['ORM', 'GORM', 'yes', 'ActiveRecord (mature)'],
  ['Admin panel', 'Generated as code you own', 'yes', 'DIY or a gem (ActiveAdmin/Avo)'],
  ['Frontend model', 'Typed React (Next.js / Vite)', 'yes', 'Server-rendered + Hotwire'],
  ['End-to-end types', 'Go → TS + Zod', 'yes', 'Dynamic (no static types)'],
  ['Mobile + desktop clients', 'Generated (Expo, Wails)', 'yes', 'DIY / API + separate app'],
  ['Auth, jobs, cache, mail', 'Built in', 'yes', 'Built in / mature gems'],
  ['Ecosystem maturity', 'Young (2024+)', 'partial', '20 years, a gem for everything'],
  ['Community + hiring pool', 'Small and new', 'partial', 'Huge, well-established'],
]

export default function CompareRailsPage() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <span className="font-mono text-xs uppercase tracking-wider text-primary">Compare</span>
        <h1 className="mb-4 mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Grit vs Ruby on Rails
        </h1>
        <p className="mb-10 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          Rails is the framework that made &quot;batteries-included, convention-over-configuration&quot;
          the standard — and it directly inspired Grit&apos;s philosophy. The difference is the
          stack underneath: Rails gives you batteries and convention on Ruby with server-rendered
          views; Grit gives you the same on Go with a typed React admin and generated typed clients
          for web, mobile and desktop.
        </p>

        {/* Short version */}
        <div className="mb-12 grid gap-4 sm:grid-cols-2">
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">The short version</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Reach for <strong>Rails</strong> when you want a proven, one-language stack with two
              decades of gems, a huge hiring pool, and the ergonomics that make server-rendered apps
              a joy to build. Reach for <strong>Grit</strong> when you want that same
              batteries-included feel on a Go backend — typed end-to-end, with a generated admin and
              ready-made web, mobile and desktop clients from one API.
            </p>
          </div>
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">Same spirit, different stack</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Both bet that convention beats configuration and that a framework should ship auth,
              jobs, mail and an ORM in the box. Rails renders HTML on the server and layers
              reactivity with Hotwire; Grit serves a Go API that a typed React SPA and admin consume,
              plus generated Expo and Wails clients. One favors Ruby&apos;s simplicity; the other
              favors Go&apos;s speed and static types.
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
                <th className="px-4 py-3 font-semibold text-foreground/70">Ruby on Rails</th>
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
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want Go&apos;s performance and a single self-hosted binary.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want end-to-end type safety from Go structs to TS + Zod.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You need a typed React admin generated as code you own.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> One API must feed web, mobile and desktop clients out of the box.</li>
            </ul>
          </div>
          <div>
            <h2 className="mb-3 text-xl font-semibold tracking-tight">Choose Rails when</h2>
            <ul className="space-y-2 text-sm leading-relaxed text-muted-foreground">
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You value a mature, 20-year ecosystem with a gem for everything.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You want one language and Ruby&apos;s famous developer happiness.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> Server-rendered views with Hotwire cover your interactivity needs.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> A large community, deep docs, and a wide hiring pool matter to you.</li>
            </ul>
          </div>
        </div>

        {/* Honest note */}
        <div className="mb-12 rounded-2xl border border-border bg-card/40 p-6">
          <h2 className="mb-2 font-semibold text-foreground">The honest trade-off</h2>
          <p className="text-sm leading-relaxed text-muted-foreground">
            Rails is more mature by two decades — a bigger community, more gems, more answers, and one
            language to learn. Grit is young and opinionated: it commits to Gin, GORM and Next.js
            rather than letting you swap them, and it asks you to work in two languages (Go and
            TypeScript). In exchange you get Go&apos;s speed, static types all the way to the client,
            a generated admin, and typed web/mobile/desktop clients from one API. If Ruby&apos;s
            ecosystem and server-rendered simplicity fit your team, Rails is a fantastic choice. If
            you want the batteries-included feel on a fast, typed, multi-client Go stack, that&apos;s
            what Grit is for.
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
