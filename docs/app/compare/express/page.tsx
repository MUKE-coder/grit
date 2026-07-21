import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, Check, X, Minus } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

export const metadata: Metadata = {
  title: 'Grit vs Express — an assembled full stack vs a minimal Node router',
  description:
    'Express is a minimal, unopinionated Node router you assemble a stack around; Grit is the assembled, opinionated full stack — Go backend, generated admin, and typed clients from one command. When to pick each, honestly compared.',
  alternates: { canonical: 'https://gritframework.dev/compare/express' },
}

const rows: [string, string, 'yes' | 'no' | 'partial', string][] = [
  ['Backend language', 'Go (Gin + GORM)', 'yes', 'JavaScript / TypeScript'],
  ['What it is', 'Assembled full stack', 'yes', 'Routing + middleware only'],
  ['ORM / database layer', 'Built in (GORM)', 'yes', 'Bring your own (Prisma, etc.)'],
  ['Typed API → client', 'Generated Go → TS + Zod', 'yes', 'Wire it up yourself'],
  ['Admin panel', 'Generated, Filament-like', 'yes', 'Build it yourself'],
  ['Auth + RBAC', 'Built in (roles + permissions)', 'yes', 'Add Passport / custom'],
  ['Background jobs / cron', 'Built in (asynq)', 'yes', 'Add a queue yourself'],
  ['File storage, email, cache', 'Built in (S3/R2, Resend, Redis)', 'yes', 'Wire up each yourself'],
  ['Mobile + desktop clients', 'Generated (Expo, Wails)', 'yes', 'Not in scope'],
  ['Flexibility to swap parts', 'Opinionated (Gin + GORM)', 'partial', 'Total — you choose everything'],
  ['Ecosystem + community', 'Younger', 'partial', 'Enormous, the Node default'],
]

export default function CompareExpressPage() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <span className="font-mono text-xs uppercase tracking-wider text-primary">Compare</span>
        <h1 className="mb-4 mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Grit vs Express
        </h1>
        <p className="mb-10 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          They sit at opposite ends of the same spectrum. Express is a minimal, unopinionated Node
          router — just routing and middleware — that you assemble a stack around. Grit is the
          assembled stack: a Go backend, a generated React admin, and typed clients from one command.
          Express is the glue; Grit is the thing already glued together.
        </p>

        {/* Short version */}
        <div className="mb-12 grid gap-4 sm:grid-cols-2">
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">The short version</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Reach for <strong>Express</strong> when you want minimalism and full control — a small
              or simple API, an unusual architecture, a stack you want to assemble yourself piece by
              piece. Reach for <strong>Grit</strong> when you don&apos;t want to spend the first
              weeks wiring auth, RBAC, an admin, jobs, and storage together — because it ships with
              all of them already connected.
            </p>
          </div>
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">Glue vs assembled</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Express gives you routing and middleware, then you choose and wire up everything
              else — ORM, auth, validation, jobs, structure. Grit makes those choices for you and
              generates them: <code>grit generate resource</code> emits a model, service, handler,
              Zod schema, TS types, React Query hooks, and an admin page in one command.
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
                <th className="px-4 py-3 font-semibold text-foreground/70">Express</th>
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
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You don&apos;t want to hand-assemble auth, RBAC, admin, and jobs first.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want the batteries included: storage, email, cache, cron, backups.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> More than one client will consume the API (web + mobile + admin).</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You&apos;d rather generate CRUD and typed clients than write them each time.</li>
            </ul>
          </div>
          <div>
            <h2 className="mb-3 text-xl font-semibold tracking-tight">Choose Express when</h2>
            <ul className="space-y-2 text-sm leading-relaxed text-muted-foreground">
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You want minimalism and full control over every layer of the stack.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> The API is small or simple and doesn&apos;t need the batteries.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You want to assemble exactly your own stack with no baked-in opinions.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> One language (JS/TS), a huge ecosystem, and no lock-in matter most.</li>
            </ul>
          </div>
        </div>

        {/* Honest note */}
        <div className="mb-12 rounded-2xl border border-border bg-card/40 p-6">
          <h2 className="mb-2 font-semibold text-foreground">The honest trade-off</h2>
          <p className="text-sm leading-relaxed text-muted-foreground">
            Express&apos;s minimalism is its strength: it makes no choices for you, runs everywhere
            in Node, and has an ecosystem and community nothing here can match. That&apos;s also the
            trade-off — you assemble and maintain the stack yourself. Grit makes the opposite bet: it
            picks Gin and GORM, adds a second language (Go), and hands you the batteries and admin
            pre-wired, at the cost of that flexibility and a younger ecosystem. If you want to build
            your stack exactly, use Express. If you want the stack already built, use Grit.
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
