import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, Check, X, Minus } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

export const metadata: Metadata = {
  title: 'Grit vs MEAN — Go + React generated vs Mongo + Express + Angular + Node',
  description:
    'MEAN is a JavaScript stack you assemble yourself — MongoDB, Express, Angular, Node. Grit is one generated, typed, batteries-included stack: a Go backend, a React admin, and clients from one command. When to pick each, honestly compared.',
  alternates: { canonical: 'https://gritframework.dev/compare/mean' },
}

const rows: [string, string, 'yes' | 'no' | 'partial', string][] = [
  ['Backend language', 'Go (Gin + GORM)', 'yes', 'JavaScript / TypeScript (Node + Express)'],
  ['Frontend framework', 'React (Next.js or Vite/TanStack)', 'yes', 'Angular (opinionated, TS)'],
  ['Database', 'PostgreSQL (relational)', 'yes', 'MongoDB (document / NoSQL)'],
  ['Resource generation', 'model + service + handler + hooks + admin', 'yes', 'Hand-written, nothing generated'],
  ['Typed API → client', 'Generated Go → TS + Zod', 'yes', 'Manual, or share TS interfaces'],
  ['Admin panel', 'Generated, Filament-like', 'yes', 'Build it yourself'],
  ['Auth + RBAC', 'Built in (roles + permissions, 2FA)', 'yes', 'Add libraries and wire it up'],
  ['Background jobs / cron', 'Built in (asynq)', 'yes', 'External queue / service'],
  ['File storage, email, cache', 'Built in (S3/R2, Resend, Redis)', 'yes', 'Wire up each yourself'],
  ['Mobile + desktop clients', 'Generated (Expo, Wails)', 'yes', 'Assemble separately'],
  ['Database backups', 'Built in, scheduled + restore', 'yes', 'Roll your own'],
  ['One language everywhere', 'No (Go + TypeScript)', 'no', 'Yes (all JS/TS)'],
]

export default function CompareMeanPage() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <span className="font-mono text-xs uppercase tracking-wider text-primary">Compare</span>
        <h1 className="mb-4 mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Grit vs the MEAN stack
        </h1>
        <p className="mb-10 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          MEAN — MongoDB, Express, Angular, Node — is a popular full-stack JavaScript stack you
          assemble yourself, in one language, choosing and wiring each piece. Grit is one generated,
          typed, batteries-included stack: a Go backend, a React admin, and web/mobile/desktop
          clients from a single command. Same goal — ship a full-stack app — very different bet on
          how you get there.
        </p>

        {/* Short version */}
        <div className="mb-12 grid gap-4 sm:grid-cols-2">
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">The short version</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Reach for <strong>MEAN</strong> when you want to stay all-JavaScript, lean on
              Angular&apos;s built-in structure for a large frontend team, and keep the freedom to
              assemble each layer — including a document database — exactly how you like. Reach for{' '}
              <strong>Grit</strong> when you&apos;d rather have the backend, admin, auth, jobs, and
              typed clients generated and wired for you out of the box.
            </p>
          </div>
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">Assembly vs generation</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              MEAN is a set of parts you connect: Express routes, an Angular app, a Mongo schema,
              and every battery (auth, admin, jobs) added by hand. Grit is one opinionated stack
              where <code>grit generate resource</code> emits the model, service, handler, Zod
              schema, TS types, React Query hooks, and an admin page together — typed end to end.
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
                <th className="px-4 py-3 font-semibold text-foreground/70">MEAN (assembled)</th>
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
                      {mark === 'no' && <X className="h-4 w-4 shrink-0 text-muted-foreground/60" />}
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
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want the backend, admin, auth, and jobs generated — not assembled by hand.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> A relational database (Postgres) fits your data better than documents.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want Go&apos;s performance and a single self-hosted binary.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You prefer React, and want typed web, mobile, and desktop clients from one API.</li>
            </ul>
          </div>
          <div>
            <h2 className="mb-3 text-xl font-semibold tracking-tight">Choose MEAN when</h2>
            <ul className="space-y-2 text-sm leading-relaxed text-muted-foreground">
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You want to stay in one language — JavaScript/TypeScript top to bottom.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> A large team benefits from Angular&apos;s built-in structure (modules, DI, RxJS).</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> Your data is a natural fit for a document (NoSQL) database like MongoDB.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You value flexibility and a huge, mature JS community over an opinionated stack.</li>
            </ul>
          </div>
        </div>

        {/* Honest note */}
        <div className="mb-12 rounded-2xl border border-border bg-card/40 p-6">
          <h2 className="mb-2 font-semibold text-foreground">The honest trade-off</h2>
          <p className="text-sm leading-relaxed text-muted-foreground">
            Grit asks you to work in two languages (Go and TypeScript) and lives in a younger
            ecosystem than the battle-tested JS stacks. In exchange you skip the assembly: the admin,
            auth, jobs, storage, and typed clients are generated and wired from day one. MEAN keeps
            everything in one language with maximum freedom, but that freedom is also the work — you
            build the admin, choose and connect each battery, and keep the pieces in sync yourself.
            If staying all-JS with Angular and a document DB matters most, MEAN is the honest pick;
            if you want a generated, typed Go backend with React and the batteries already in place,
            that&apos;s Grit.
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
