import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, Check, X, Minus } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

export const metadata: Metadata = {
  title: 'Grit vs Encore — two Go frameworks: batteries + React admin vs infrastructure-from-code',
  description:
    'Encore is a Go backend framework built around infrastructure-from-code and observability; Grit is a full-stack Go framework with a generated React admin, batteries, and typed multi-platform clients. When to pick each, honestly compared.',
  alternates: { canonical: 'https://gritframework.dev/compare/encore' },
}

const rows: [string, string, 'yes' | 'no' | 'partial', string][] = [
  ['Backend language', 'Go (Gin + GORM)', 'yes', 'Go or TypeScript'],
  ['Core idea', 'Full-stack app + generated admin', 'yes', 'Infrastructure from code'],
  ['Admin panel', 'Generated, Filament-like', 'yes', 'Not in scope'],
  ['Frontend / mobile / desktop clients', 'Generated (Next.js/Vite, Expo, Wails)', 'yes', 'Not in scope'],
  ['Typed API → client', 'Generated Go → TS + Zod', 'yes', 'Type-safe service-to-service'],
  ['Auth + RBAC', 'Built in (roles + permissions)', 'yes', 'Bring your own'],
  ['Infra provisioning', 'You wire it (Docker)', 'partial', 'Automatic, local + cloud'],
  ['Observability / tracing', 'Audit log + Sentinel', 'partial', 'Built-in distributed tracing'],
  ['Background jobs / cron', 'Built in (asynq)', 'yes', 'Built in (cron, pub/sub)'],
  ['Deploy model', 'Self-host (container / single binary)', 'yes', 'Own cloud or Encore Cloud'],
]

export default function CompareEncorePage() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <span className="font-mono text-xs uppercase tracking-wider text-primary">Compare</span>
        <h1 className="mb-4 mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Grit vs Encore
        </h1>
        <p className="mb-10 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          This is the most apples-to-apples backend comparison we have — both are Go frameworks, and
          Encore is a strong, well-engineered one. But they aim at different problems. Encore is a
          backend framework built around <em>infrastructure from code</em>: you declare databases,
          pub/sub, caches, and cron in your Go, and Encore provisions and wires them. Grit is a
          full-stack meta-framework whose headline is a generated React admin, batteries, and typed
          web/mobile/desktop clients from one command.
        </p>

        {/* Short version */}
        <div className="mb-12 grid gap-4 sm:grid-cols-2">
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">The short version</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Reach for <strong>Encore</strong> when you&apos;re building a distributed Go backend —
              microservices, service-to-service calls, and infrastructure you&apos;d rather declare
              than hand-provision, with tracing and API docs out of the box. Reach for{' '}
              <strong>Grit</strong> when you&apos;re building a full-stack product — a Go API, a
              generated admin panel, auth and roles, and web/mobile/desktop clients — and you&apos;d
              rather generate the whole app than assemble the frontend and admin yourself.
            </p>
          </div>
          <div className="rounded-2xl border border-border bg-card/40 p-6">
            <h2 className="mb-2 font-semibold text-foreground">Different centers of gravity</h2>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Encore&apos;s gravity is the backend and its infrastructure — provisioning, wiring, and
              observability, with a strong local dev dashboard and an optional managed cloud. Grit&apos;s
              gravity is the full stack — one <code>grit generate resource</code> gives you a model,
              service, handler, Zod schema, TS types, React Query hooks, and an admin page. Neither
              tries to be the other.
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
                <th className="px-4 py-3 font-semibold text-foreground/70">Encore</th>
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
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want a generated, Filament-like admin panel out of the box.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You need typed web, mobile, and desktop clients from one API.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You want batteries — auth, RBAC, storage, email, jobs, backups — included.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" /> You&apos;d rather self-host plain Go + Docker with no platform to adopt.</li>
            </ul>
          </div>
          <div>
            <h2 className="mb-3 text-xl font-semibold tracking-tight">Choose Encore when</h2>
            <ul className="space-y-2 text-sm leading-relaxed text-muted-foreground">
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You&apos;re building a distributed backend of microservices, not a single app.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> You want infrastructure declared in code and provisioned automatically.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> Built-in distributed tracing, API docs, and a dev dashboard matter to you.</li>
              <li className="flex gap-2"><Check className="mt-0.5 h-4 w-4 shrink-0 text-sky-500" /> Type-safe service-to-service calls and cloud provisioning are the point.</li>
            </ul>
          </div>
        </div>

        {/* Honest note */}
        <div className="mb-12 rounded-2xl border border-border bg-card/40 p-6">
          <h2 className="mb-2 font-semibold text-foreground">The honest trade-off</h2>
          <p className="text-sm leading-relaxed text-muted-foreground">
            Encore does things Grit doesn&apos;t: it automates infrastructure provisioning, gives you
            distributed tracing and API docs for free, and is genuinely built for microservices — if
            that&apos;s your shape, it&apos;s excellent, and its managed cloud is a real convenience.
            Grit doesn&apos;t provision infrastructure or trace distributed calls; it assumes a more
            conventional single-service app you deploy as a container. What Grit gives back is the
            whole front of the app — a generated admin and typed web/mobile/desktop clients — plus
            batteries and self-hosted plain Go with no platform lock-in. If your product is a
            distributed backend, use Encore. If it&apos;s a full-stack app that needs an admin and
            clients, that&apos;s the part Encore leaves to you — and the part Grit generates.
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
