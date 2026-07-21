import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, Check, X, Minus } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

export const metadata: Metadata = {
  title: 'Grit vs Next.js — full-stack Go backend or a React framework?',
  description:
    'Next.js is a React framework with API routes bolted on; Grit is a Go backend, a React admin, and typed clients generated from one command. When to pick each, honestly compared.',
  alternates: { canonical: 'https://gritframework.dev/compare/nextjs' },
}

const rows: [string, string, 'yes' | 'no' | 'partial', string][] = [
  ['Backend language', 'Go (Gin + GORM)', 'yes', 'JavaScript / TypeScript'],
  ['Typed API → client', 'Generated Go → TS + Zod', 'yes', 'Manual, or tRPC add-on'],
  ['Admin panel', 'Generated, Filament-like', 'yes', 'Build it yourself'],
  ['Auth + RBAC', 'Built in (roles + permissions)', 'yes', 'Add a library (NextAuth, etc.)'],
  ['Background jobs / cron', 'Built in (asynq)', 'yes', 'External queue / service'],
  ['File storage, email, cache', 'Built in (S3/R2, Resend, Redis)', 'yes', 'Wire up each yourself'],
  ['Mobile + desktop clients', 'Generated (Expo, Wails)', 'yes', 'Not in scope'],
  ['Database backups', 'Built in, scheduled + restore', 'yes', 'Roll your own'],
  ['Rendering / SSR / SEO', 'Next.js on the frontend', 'partial', 'Best-in-class'],
  ['Edge / serverless deploy', 'Container (self-host)', 'partial', 'First-class (Vercel)'],
]

export default function CompareNextjsPage() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <span className="font-mono text-xs uppercase tracking-wider text-primary">Compare</span>
        <h1 className="mb-4 mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Grit vs Next.js
        </h1>
        <p className="mb-10 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          They&apos;re not really the same tool. Next.js is a React framework that can also serve API
          routes; Grit is a Go backend, a generated React admin, and typed web/mobile/desktop
          clients from one command. In fact, Grit apps often <em>use</em> Next.js for the web
          frontend — the question is what runs behind it.
        </p>

        {/* Short version */}
        <div className="mb-12 grid gap-4 sm:grid-cols-2">
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">The short version</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Reach for <strong>Next.js alone</strong> when the app is mostly a frontend — a
              marketing site, a content site, a dashboard whose data comes from a service you
              already run. Reach for <strong>Grit</strong> when you need a real backend — auth,
              roles, a database, jobs, an admin, an API several clients consume — and you don&apos;t
              want to assemble ten libraries to get there.
            </p>
          </div>
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">They compose</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              This isn&apos;t either/or. A Grit project&apos;s web app is Next.js (App Router).
              You get Next.js&apos;s rendering and SEO on the frontend and a Go API, admin, and
              typed clients behind it — instead of pushing database and auth logic into API routes
              that get hard to test and reuse.
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
                <th className="px-4 py-3 font-semibold text-foreground/70">Next.js (alone)</th>
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
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You need a real backend: database, auth, roles, jobs, an admin.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> More than one client will consume the API (web + mobile + admin).</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want Go&apos;s performance and a single self-hosted binary.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You&apos;d rather generate CRUD than hand-write it for every model.</li>
            </ul>
          </div>
          <div>
            <h2 className="mb-3 text-xl font-semibold tracking-tight">Choose Next.js alone when</h2>
            <ul className="space-y-2 text-sm leading-relaxed text-muted-foreground">
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> The app is frontend-first: marketing, content, docs, a storefront.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> Your data already lives in a CMS or a backend you run elsewhere.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You want edge/serverless deploys on Vercel with zero infra.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> Your team is all-in on TypeScript and doesn&apos;t want a second language.</li>
            </ul>
          </div>
        </div>

        {/* Honest note */}
        <div className="mb-12 rounded-2xl border border-border bg-card/40 p-6">
          <h2 className="mb-2 font-semibold text-foreground">The honest trade-off</h2>
          <p className="text-sm leading-relaxed text-muted-foreground">
            Grit adds a second language (Go) and expects you to run a container rather than deploy to
            the edge. In exchange you get a typed, testable backend and an admin you didn&apos;t
            build. If your product is essentially a frontend, that&apos;s overhead you don&apos;t
            need — use Next.js. If it&apos;s a product with a backend, pushing all of it into Next.js
            API routes is the overhead you&apos;ll feel later.
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
