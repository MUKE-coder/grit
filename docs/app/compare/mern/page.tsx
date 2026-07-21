import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, Check, X, Minus } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

export const metadata: Metadata = {
  title: 'Grit vs MERN — one generated stack vs Mongo + Express + React + Node by hand',
  description:
    'MERN is a popular DIY stack you assemble yourself in JavaScript; Grit is one generated, typed, batteries-included full-stack framework in Go + React. When to pick each, honestly compared.',
  alternates: { canonical: 'https://gritframework.dev/compare/mern' },
}

const rows: [string, string, 'yes' | 'no' | 'partial', string][] = [
  ['What it is', 'One generated framework', 'yes', 'A stack you assemble yourself'],
  ['Backend language', 'Go (Gin + GORM)', 'yes', 'JavaScript / Node (Express)'],
  ['Database', 'PostgreSQL (relational)', 'yes', 'MongoDB (document / NoSQL)'],
  ['CRUD + resources', 'Generated in one command', 'yes', 'Hand-write each route + model'],
  ['Typed API → client', 'Generated Go → TS + Zod', 'yes', 'Manual, or add tRPC / codegen'],
  ['Admin panel', 'Generated, Filament-like', 'yes', 'Build it yourself'],
  ['Auth + RBAC', 'Built in (roles + permissions)', 'yes', 'Wire up JWT / Passport yourself'],
  ['Background jobs / cron', 'Built in (asynq)', 'yes', 'Add a queue (BullMQ, etc.)'],
  ['Storage, email, cache', 'Built in (S3/R2, Resend, Redis)', 'yes', 'Wire up each yourself'],
  ['Mobile + desktop clients', 'Generated (Expo, Wails)', 'yes', 'Assemble separately'],
  ['One language end to end', 'No — Go + TypeScript', 'partial', 'Yes — JavaScript everywhere'],
]

export default function CompareMernPage() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <span className="font-mono text-xs uppercase tracking-wider text-primary">Compare</span>
        <h1 className="mb-4 mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Grit vs the MERN stack
        </h1>
        <p className="mb-10 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          These aren&apos;t the same kind of thing. MERN — MongoDB, Express, React, Node — is a
          popular <em>stack</em> you assemble yourself, all in JavaScript. Grit is one generated,
          batteries-included framework: a Go backend, a generated admin, and typed clients from a
          single command. MERN is one of the friendliest ways to learn and ship full-stack JS; the
          question is whether you want to wire the pieces or have them generated.
        </p>

        {/* Short version */}
        <div className="mb-12 grid gap-4 sm:grid-cols-2">
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">The short version</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Reach for <strong>MERN</strong> when you want to stay in one language, keep full
              control over every layer, and lean on a huge community and a document database. Reach
              for <strong>Grit</strong> when you&apos;d rather skip the assembly — a typed backend,
              auth, roles, jobs, an admin, and multiple clients generated for you — and you&apos;re
              happy with a relational database and a second language for the backend.
            </p>
          </div>
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">DIY vs generated</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              With MERN, nothing is generated: you build CRUD, auth, validation, and an admin by
              hand, exactly how you like them. With Grit, <code>grit generate resource</code>
              &nbsp;emits the model, service, handler, Zod schema, TS types, React Query hooks, and
              an admin page at once — the wiring MERN leaves to you.
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
                <th className="px-4 py-3 font-semibold text-foreground/70">MERN (assembled)</th>
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
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You&apos;d rather generate CRUD, an admin, and clients than build them by hand.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> Your data is relational and Postgres fits it better than documents.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want types flowing end to end — Go → TS + Zod — with no drift.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want Go&apos;s performance and batteries (auth, jobs, storage) included.</li>
            </ul>
          </div>
          <div>
            <h2 className="mb-3 text-xl font-semibold tracking-tight">Choose MERN when</h2>
            <ul className="space-y-2 text-sm leading-relaxed text-muted-foreground">
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You want one language — JavaScript — across the entire stack.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> Your data is unstructured or document-shaped and a NoSQL DB fits.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You want full control and no opinions imposed on your structure.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You value the huge community and the gentlest learning ramp.</li>
            </ul>
          </div>
        </div>

        {/* Honest note */}
        <div className="mb-12 rounded-2xl border border-border bg-card/40 p-6">
          <h2 className="mb-2 font-semibold text-foreground">The honest trade-off</h2>
          <p className="text-sm leading-relaxed text-muted-foreground">
            Grit is opinionated, adds a second language (Go), and rides a younger ecosystem than the
            JS stacks — MERN has years of tutorials, jobs, and Stack Overflow answers behind it, and
            keeping everything in JavaScript is genuinely simpler to learn and hire for. In exchange,
            Grit hands you a typed backend, an admin, and clients you didn&apos;t assemble. If you
            want to stay all-JS with a document database and control every layer, MERN is a great
            choice. If you&apos;d rather skip the wiring and start from a generated, typed,
            batteries-included stack, that&apos;s Grit.
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
