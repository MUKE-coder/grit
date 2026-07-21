import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, Check, X, Minus } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

export const metadata: Metadata = {
  title: 'Grit vs NestJS — batteries-included Go+React vs a structured Node backend',
  description:
    'NestJS is a structured, opinionated Node.js backend framework; Grit is a batteries-included Go backend with a generated React admin and typed clients. When to pick each, honestly compared.',
  alternates: { canonical: 'https://gritframework.dev/compare/nestjs' },
}

const rows: [string, string, 'yes' | 'no' | 'partial', string][] = [
  ['Backend language', 'Go (Gin + GORM)', 'yes', 'Node.js / TypeScript'],
  ['Typed API → client', 'Generated Go → TS + Zod', 'yes', 'Manual, or a codegen add-on'],
  ['Admin panel', 'Generated, Filament-like', 'yes', 'Build it yourself'],
  ['Full CRUD generation', 'One command (model → admin)', 'yes', 'Scaffold pieces, wire the rest'],
  ['Auth + RBAC', 'Built in (roles + permissions)', 'yes', 'Add Passport / a library'],
  ['Background jobs / cron', 'Built in (asynq)', 'yes', 'Add BullMQ / a queue'],
  ['File storage, email, cache', 'Built in (S3/R2, Resend, Redis)', 'yes', 'Wire up each yourself'],
  ['Mobile + desktop clients', 'Generated (Expo, Wails)', 'yes', 'Not in scope'],
  ['ORM flexibility', 'GORM (opinionated)', 'partial', 'TypeORM / Prisma / Mongoose'],
  ['Microservices / WebSockets', 'Possible, less turnkey', 'partial', 'First-class'],
]

export default function CompareNestjsPage() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <span className="font-mono text-xs uppercase tracking-wider text-primary">Compare</span>
        <h1 className="mb-4 mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Grit vs NestJS
        </h1>
        <p className="mb-10 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          Both are opinionated and both give your backend a real shape — but they draw the box
          differently. NestJS is a Node.js backend framework: modules, controllers, providers,
          dependency injection, and you assemble auth, jobs, storage, and an admin on top. Grit is
          a batteries-included full-stack meta-framework that generates most of that — plus a React
          admin and typed clients — in Go.
        </p>

        {/* Short version */}
        <div className="mb-12 grid gap-4 sm:grid-cols-2">
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">The short version</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Reach for <strong>NestJS</strong> when you want a well-structured backend and the
              freedom to pick your own ORM, auth, and pieces — all in TypeScript, sharing types with
              a TS frontend without codegen. Reach for <strong>Grit</strong> when you&apos;d rather
              generate the batteries — auth, RBAC, jobs, storage, an admin, typed web/mobile/desktop
              clients — from one command, and you&apos;re happy running Go behind your frontend.
            </p>
          </div>
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">Different scope</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              NestJS is a <em>backend</em> framework — it doesn&apos;t ship an admin panel, generate
              CRUD, or produce typed clients; those stay your job. Grit spans the stack: one command
              gives you a model, service, handler, Zod schema, TS types, React Query hooks, and an
              admin page — with the trade-off that Grit adds Go as a second language where Nest is
              all-TypeScript.
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
                <th className="px-4 py-3 font-semibold text-foreground/70">NestJS</th>
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
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want batteries generated for you: auth, RBAC, jobs, storage, an admin.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> More than one client will consume the API (web + mobile + desktop + admin).</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want Go&apos;s performance and a single self-hosted binary.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You&apos;d rather generate CRUD and typed clients than hand-write them per model.</li>
            </ul>
          </div>
          <div>
            <h2 className="mb-3 text-xl font-semibold tracking-tight">Choose NestJS when</h2>
            <ul className="space-y-2 text-sm leading-relaxed text-muted-foreground">
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You want one language end to end and shared types with a TS frontend, no codegen.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You want to choose your own ORM, auth, and pieces rather than accept defaults.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You&apos;re building microservices, GraphQL, or WebSocket-heavy services.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You want the depth of the npm ecosystem and your team is all-in on TypeScript.</li>
            </ul>
          </div>
        </div>

        {/* Honest note */}
        <div className="mb-12 rounded-2xl border border-border bg-card/40 p-6">
          <h2 className="mb-2 font-semibold text-foreground">The honest trade-off</h2>
          <p className="text-sm leading-relaxed text-muted-foreground">
            Grit trades flexibility for batteries: you get a generated admin, typed clients, and
            wired-in auth, jobs, and storage — but in Go, a second language alongside your frontend,
            with GORM as the ORM and a younger ecosystem than Node&apos;s. NestJS keeps everything in
            TypeScript and lets you pick every piece — but that admin, that CRUD, that RBAC, and those
            clients are yours to build and maintain. If you value structure plus freedom of choice,
            NestJS fits. If you value generated batteries, multi-client output, and Go performance,
            Grit fits.
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
