import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, Check, X, Minus } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

export const metadata: Metadata = {
  title: 'Grit vs Laravel — Go + React vs PHP’s productivity king',
  description:
    'Laravel + Filament give you batteries and a great admin on PHP; Grit gives you batteries and a generated admin on Go + React. Laravel&apos;s admin is a runtime package; Grit&apos;s is code you own. Honestly compared.',
  alternates: { canonical: 'https://gritframework.dev/compare/laravel' },
}

const rows: [string, string, 'yes' | 'no' | 'partial', string][] = [
  ['Backend language', 'Go (Gin + GORM)', 'yes', 'PHP (Eloquent)'],
  ['CLI scaffolding', 'grit generate resource', 'yes', 'Artisan make:*'],
  ['Admin panel', 'Generated code you own', 'yes', 'Filament (runtime package)'],
  ['Typed clients from the API', 'Go → TS + Zod, web/mobile/desktop', 'yes', 'Web via Blade/Livewire/Inertia'],
  ['Auth + RBAC', 'Built in (roles + permissions)', 'yes', 'Sanctum + spatie/permission'],
  ['Background jobs / cron', 'Built in (asynq)', 'yes', 'Queues + Horizon'],
  ['File storage, email, cache', 'Built in (S3/R2, Resend, Redis)', 'yes', 'First-class, well documented'],
  ['Ecosystem & community', 'Younger, smaller', 'partial', 'Enormous, mature'],
  ['First-party paid tooling', 'Self-hosted, open', 'partial', 'Forge, Vapor, Nova, Cashier'],
  ['Number of languages', 'Two (Go + TS)', 'partial', 'One (PHP)'],
]

export default function CompareLaravelPage() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <span className="font-mono text-xs uppercase tracking-wider text-primary">Compare</span>
        <h1 className="mb-4 mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Grit vs Laravel
        </h1>
        <p className="mb-10 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          Laravel is the framework a lot of us measure everything else against — batteries included,
          superb docs, a mature ecosystem, and in Filament, one of the best admin panels anywhere.
          Grit owes it a real debt: Filament directly inspired Grit&apos;s generated admin. The
          honest difference is the stack underneath — Laravel gives you all of this on PHP; Grit
          gives you a similar shape on Go + React, with typed clients for more than the web.
        </p>

        {/* Short version */}
        <div className="mb-12 grid gap-4 sm:grid-cols-2">
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">The short version</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Reach for <strong>Laravel</strong> when you want the most mature, best-documented
              batteries-included framework on earth, a single language, and an ecosystem with an
              answer for everything — plus Filament for a polished admin. Reach for <strong>Grit</strong>{' '}
              when you want that same all-in-one feel on a Go backend with a typed React frontend,
              and you want one API to feed web, admin, mobile, and desktop.
            </p>
          </div>
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">Credit where it&apos;s due</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Grit&apos;s admin exists because Filament showed how good a resource-driven admin can
              be. The key difference is ownership: Filament is a runtime package your admin depends
              on, while Grit&apos;s admin is generated code that lands in your repo — you read it,
              edit it, and it&apos;s yours. Same idea, different delivery.
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
                <th className="px-4 py-3 font-semibold text-foreground/70">Laravel</th>
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
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> One API should feed web, admin, mobile (Expo), and desktop (Wails).</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want end-to-end types: Go models generate TS + Zod for the client.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You&apos;d rather own the admin as generated code than depend on a package.</li>
            </ul>
          </div>
          <div>
            <h2 className="mb-3 text-xl font-semibold tracking-tight">Choose Laravel when</h2>
            <ul className="space-y-2 text-sm leading-relaxed text-muted-foreground">
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You want the deepest ecosystem and the best docs in the business.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> One language (PHP) across the whole team is a priority.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> First-party tooling — Forge, Vapor, Nova, Cashier, Horizon — is worth it.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> Filament plus Livewire/Inertia already give you exactly the shape you need.</li>
            </ul>
          </div>
        </div>

        {/* Honest note */}
        <div className="mb-12 rounded-2xl border border-border bg-card/40 p-6">
          <h2 className="mb-2 font-semibold text-foreground">The honest trade-off</h2>
          <p className="text-sm leading-relaxed text-muted-foreground">
            Laravel is more mature, has a far bigger community, is easier to hire for, and its
            ecosystem is genuinely hard to beat — Grit is younger, smaller, and opinionated, and it
            asks you to work in two languages (Go and TypeScript) instead of one. What you get back
            is a Go-fast, typed-end-to-end backend, an admin you own as code rather than pull in as a
            package, and generated clients for web, mobile, and desktop from the same API. If the
            ecosystem and single-language simplicity matter most, Laravel is the safe, excellent
            choice. If you want Go and typed multi-client output with the same batteries-included
            feel, that&apos;s where Grit fits.
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
