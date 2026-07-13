import Link from 'next/link'
import {
  ArrowRight,
  ArrowLeft,
  Zap,
  Server,
  Globe,
  Shield,
  HardDrive,
  Mail,
  Layers,
  Clock,
  Database,
  Bot,
  ShieldCheck,
  Activity,
  Table2,
  Flag,
  Webhook,
  Radio,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { Callout } from '@/components/callout'
import { PageHelp } from '@/components/page-help'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/getting-started/performance')

/* ─── Batteries data ─────────────────────────────────────────────── */
const BATTERIES: {
  icon: typeof Shield
  name: string
  desc: string
  href: string
}[] = [
  {
    icon: Shield,
    name: 'Authentication',
    desc: 'JWT login/register/refresh, bcrypt hashing, OAuth (Google/GitHub), TOTP 2FA, RBAC.',
    href: '/docs/backend/authentication',
  },
  {
    icon: ShieldCheck,
    name: 'RBAC',
    desc: 'ADMIN / EDITOR / USER roles, RequireRole middleware, role-restricted routes.',
    href: '/docs/backend/rbac',
  },
  {
    icon: HardDrive,
    name: 'File storage',
    desc: 'S3-compatible (AWS S3 / Cloudflare R2 / MinIO) with presigned uploads + image processing.',
    href: '/docs/batteries/storage',
  },
  {
    icon: Mail,
    name: 'Email',
    desc: 'Resend transactional sender with HTML templates — welcome, reset, verify, notify.',
    href: '/docs/batteries/email',
  },
  {
    icon: Layers,
    name: 'Background jobs',
    desc: 'Redis-backed asynq queue with email / image / cleanup workers + admin dashboard.',
    href: '/docs/batteries/jobs',
  },
  {
    icon: Clock,
    name: 'Cron scheduler',
    desc: 'Recurring tasks via asynq scheduler, cron expressions, admin monitoring.',
    href: '/docs/batteries/cron',
  },
  {
    icon: Database,
    name: 'Redis caching',
    desc: 'Cache service + response-cache middleware; hot public reads served in under a millisecond.',
    href: '/docs/batteries/caching',
  },
  {
    icon: Bot,
    name: 'AI Gateway',
    desc: 'Vercel AI Gateway — one key, many models — with streaming completion and chat.',
    href: '/docs/batteries/ai',
  },
  {
    icon: Shield,
    name: 'Security (Sentinel)',
    desc: 'Built-in WAF, per-IP rate limiting, brute-force lockout, anomaly + geo gating.',
    href: '/docs/batteries/security',
  },
  {
    icon: Activity,
    name: 'Pulse observability',
    desc: 'Self-hosted request tracing, DB + runtime metrics, error tracking, Prometheus export.',
    href: '/docs/backend/pulse',
  },
  {
    icon: Table2,
    name: 'GORM Studio',
    desc: 'Visual database browser embedded at /studio — browse and edit rows, no external tool.',
    href: '/docs/infrastructure/database',
  },
  {
    icon: Flag,
    name: 'Feature flags',
    desc: 'FeatureFlag engine: percentage rollouts, allow/block lists, sticky bucketing, A/B variants.',
    href: '/docs/backend/feature-flags',
  },
  {
    icon: Webhook,
    name: 'Webhooks',
    desc: 'Universal receiver with Stripe / GitHub / HMAC verifiers, idempotent dedupe, admin replay.',
    href: '/docs/backend/webhooks',
  },
  {
    icon: Radio,
    name: 'Realtime',
    desc: 'WebSocket hub with JWT handshake, SendToUser vs Broadcast, multi-device fan-out.',
    href: '/docs/backend/realtime',
  },
]

/* ─── Comparison matrix data ─────────────────────────────────────── */
type Cell = 'yes' | 'partial' | 'no'
const FRAMEWORKS = ['Grit', 'Laravel', 'Django', 'Next.js', 'Rails', 'T3'] as const
const MATRIX: { capability: string; cells: Cell[] }[] = [
  // Grit, Laravel, Django, Next.js, Rails, T3
  { capability: 'Compiled single-binary deploy', cells: ['yes', 'no', 'no', 'no', 'no', 'no'] },
  { capability: 'Generated admin panel', cells: ['yes', 'partial', 'yes', 'no', 'partial', 'no'] },
  { capability: 'End-to-end types (backend ↔ frontend)', cells: ['yes', 'no', 'no', 'yes', 'no', 'yes'] },
  { capability: 'Full-stack code generator', cells: ['yes', 'yes', 'partial', 'no', 'yes', 'no'] },
  { capability: 'Built-in auth (JWT / session)', cells: ['yes', 'yes', 'yes', 'partial', 'partial', 'yes'] },
  { capability: '2FA + OAuth out of the box', cells: ['yes', 'partial', 'partial', 'partial', 'partial', 'partial'] },
  { capability: 'File storage (S3 / R2)', cells: ['yes', 'yes', 'yes', 'no', 'yes', 'no'] },
  { capability: 'Background jobs + scheduler', cells: ['yes', 'yes', 'partial', 'no', 'yes', 'no'] },
  { capability: 'Type-safe API hooks generated', cells: ['yes', 'no', 'no', 'partial', 'no', 'yes'] },
  { capability: 'Built-in WAF / rate-limit / observability', cells: ['yes', 'partial', 'partial', 'no', 'partial', 'no'] },
  { capability: 'Single-language stack', cells: ['no', 'no', 'no', 'yes', 'no', 'yes'] },
]

function CellMark({ v }: { v: Cell }) {
  if (v === 'yes')
    return (
      <span className="inline-flex h-5 w-5 items-center justify-center rounded-full bg-primary/15 text-primary text-xs font-bold">
        ✓
      </span>
    )
  if (v === 'partial')
    return (
      <span className="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-mono text-amber-500 bg-amber-500/10 border border-amber-500/20">
        ~
      </span>
    )
  return (
    <span className="inline-flex h-5 w-5 items-center justify-center rounded-full bg-muted/40 text-muted-foreground/30 text-xs">
      —
    </span>
  )
}

export default function PerformanceBenchmarksPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Getting Started</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Performance &amp; Benchmarks
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit is a compiled Go backend with a React frontend, scaffolded with
                production-grade performance defaults and a full set of batteries already wired.
                This page covers <strong>why it&apos;s fast</strong>, <strong>what ships out of the
                box</strong>, and <strong>how it compares</strong> to Laravel, Django, and Next.js.
              </p>
              <LaneFlow
                id="gs-perf"
                lanes={['Fast by default', 'Wired from the first scaffold']}
                nodes={[
                  { id: 'core', lane: 0, row: 2, title: 'Grit defaults', sub: 'zero config', tone: 'primary' },
                  { id: 'go', lane: 1, row: 0, title: 'Compiled Go', sub: 'no runtime', tone: 'cyan' },
                  { id: 'gzip', lane: 1, row: 1, title: 'Gzip', sub: '60–80% smaller', tone: 'blue' },
                  { id: 'pool', lane: 1, row: 2, title: 'Conn pool', sub: 'tuned', tone: 'green' },
                  { id: 'cache', lane: 1, row: 3, title: 'Redis cache', sub: 'sub-ms reads', tone: 'rose' },
                  { id: 'cdn', lane: 1, row: 4, title: 'ISR / CDN', sub: 'static pages', tone: 'amber' },
                ]}
                edges={[
                  { from: 'core', to: 'go', tone: 'cyan' },
                  { from: 'core', to: 'gzip', tone: 'blue' },
                  { from: 'core', to: 'pool', label: 'ships with', tone: 'green' },
                  { from: 'core', to: 'cache', tone: 'rose' },
                  { from: 'core', to: 'cdn', tone: 'amber' },
                ]}
                legend={[
                  { tone: 'primary', label: 'Defaults' },
                  { tone: 'cyan', label: 'Backend & frontend wins' },
                ]}
                caption="Production-grade performance defaults are wired in from the very first scaffold"
              />
            </div>

            <div className="prose-grit">

              {/* ── PART 1: Why Grit is fast ─────────────────────── */}
              <div className="mb-16">
                <div className="flex items-center gap-3 mb-6">
                  <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 border border-primary/20">
                    <Zap className="h-5 w-5 text-primary" />
                  </div>
                  <div>
                    <h2 className="text-2xl font-bold tracking-tight m-0">Why Grit is fast</h2>
                    <p className="text-sm text-muted-foreground m-0">Compiled backend, production defaults, zero config</p>
                  </div>
                </div>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  The API is a single statically-linked <strong>Go binary</strong> — no interpreter,
                  no per-request VM warm-up, low memory footprint. On top of that raw speed, every
                  Grit project is scaffolded with performance optimisations baked into both the Go
                  API and the Next.js frontends. Nothing to configure.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3 flex items-center gap-2">
                  <Server className="h-4 w-4 text-primary/70" /> Backend (Go — Gin + GORM)
                </h3>
                <div className="overflow-x-auto mb-6">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-4 py-2 font-medium">Optimisation</th>
                        <th className="px-4 py-2 font-medium">What it does</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-medium text-foreground/90">Gzip compression</td><td className="px-4 py-2 text-muted-foreground">Every response compressed at BestSpeed — JSON payloads shrink 60–80%.</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-medium text-foreground/90">Connection pool tuning</td><td className="px-4 py-2 text-muted-foreground">GORM&apos;s sql pool capped and recycled (MaxOpen 100, MaxIdle 10, 30m lifetime) — no exhaustion or stale connections under load.</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-medium text-foreground/90">Request-ID tracing</td><td className="px-4 py-2 text-muted-foreground">Every request gets an X-Request-ID echoed into logs and Pulse for trivial correlation.</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-medium text-foreground/90">Cache-Control headers</td><td className="px-4 py-2 text-muted-foreground">Public read endpoints are CDN/edge-cacheable so requests never hit Go at all.</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-medium text-foreground/90">Redis response cache</td><td className="px-4 py-2 text-muted-foreground">Hot public endpoints served from memory in under a millisecond; auth routes skipped automatically.</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-medium text-foreground/90">Presigned uploads</td><td className="px-4 py-2 text-muted-foreground">Binaries go straight to S3/R2/MinIO — the API never buffers large files.</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-medium text-foreground/90">Async background jobs</td><td className="px-4 py-2 text-muted-foreground">Slow work (thumbnails, emails, webhooks) pushed to an asynq queue so handlers return immediately.</td></tr>
                      <tr><td className="px-4 py-2 font-medium text-foreground/90">Sentinel rate limiting</td><td className="px-4 py-2 text-muted-foreground">WAF + per-IP limiter sheds abusive traffic before it degrades latency for real users.</td></tr>
                    </tbody>
                  </table>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3 flex items-center gap-2">
                  <Globe className="h-4 w-4 text-emerald-500/70" /> Frontend (Next.js — App Router)
                </h3>
                <div className="overflow-x-auto mb-6">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-4 py-2 font-medium">Optimisation</th>
                        <th className="px-4 py-2 font-medium">What it does</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-medium text-foreground/90">Server Components</td><td className="px-4 py-2 text-muted-foreground">Pages fetch on the server and stream HTML — zero JS for data fetching, no client waterfalls.</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-medium text-foreground/90">ISR / revalidate</td><td className="px-4 py-2 text-muted-foreground">Public pages rendered once and served from the CDN edge; revalidated in the background.</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-medium text-foreground/90">React Query caching</td><td className="px-4 py-2 text-muted-foreground">Admin data cached in memory; back-navigation is instant and mutations auto-invalidate.</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-medium text-foreground/90">next/image</td><td className="px-4 py-2 text-muted-foreground">Automatic WebP/AVIF, correct sizing, lazy-load, CDN caching.</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-medium text-foreground/90">Turborepo cache</td><td className="px-4 py-2 text-muted-foreground">Hashed build outputs replayed in milliseconds — cuts a typical CI build from 4+ min to under 30s on the second run.</td></tr>
                      <tr><td className="px-4 py-2 font-medium text-foreground/90">Automatic code splitting</td><td className="px-4 py-2 text-muted-foreground">JS bundle split per route; heavy admin components dynamically imported.</td></tr>
                    </tbody>
                  </table>
                </div>

                <Callout type="note" title="Measured on a real run">
                  The <Link href="/docs/learnings/stateless-service-load-test">stateless-service load test</Link>{' '}
                  walks through scaffolding an <code>--api</code> project and load-testing the tiny{' '}
                  <code>/api/health</code> endpoint with k6, capturing p50 / p95 / p99. In the walkthrough&apos;s
                  example run (Gin in release mode, SQLite, k6 co-located on one box) the endpoint returns at
                  roughly <strong>p50 ~2.8 ms, p95 ~7.8 ms, p99 ~18 ms</strong> at ~165 req/s — illustrative
                  numbers for a local machine, not a hardware-normalised benchmark. See that page for the full
                  methodology and how to reproduce it against your own deployment.
                </Callout>

                <p className="text-muted-foreground leading-relaxed mt-4">
                  For the code behind each of these — middleware snippets, the connection-pool config,
                  the presigned-upload flow — see the deep dive in{' '}
                  <Link href="/docs/concepts/performance">Core Concepts → Performance</Link>.
                </p>
              </div>

              {/* ── PART 2: Batteries included ───────────────────── */}
              <div className="mb-16">
                <div className="flex items-center gap-3 mb-6">
                  <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-500/10 border border-emerald-500/20">
                    <Layers className="h-5 w-5 text-emerald-500" />
                  </div>
                  <div>
                    <h2 className="text-2xl font-bold tracking-tight m-0">Batteries included</h2>
                    <p className="text-sm text-muted-foreground m-0">Everything below is wired the moment you scaffold</p>
                  </div>
                </div>

                <p className="text-muted-foreground leading-relaxed mb-6">
                  The list of things you&apos;d normally spend the first month wiring — auth, an admin
                  panel, uploads, jobs, email, security — is treated as already solved. Every one of
                  these ships in the API of a fresh Grit project. Click through for the full docs on each.
                </p>

                <div className="grid gap-3 sm:grid-cols-2 not-prose mb-6">
                  {BATTERIES.map((b) => (
                    <Link
                      key={b.name}
                      href={b.href}
                      className="group flex gap-3 rounded-lg border border-border bg-card/40 p-3.5 transition-colors hover:border-primary/40 hover:bg-card/70"
                    >
                      <div className="shrink-0 mt-0.5 h-8 w-8 rounded-md bg-primary/10 flex items-center justify-center">
                        <b.icon className="h-4 w-4 text-primary" />
                      </div>
                      <div className="min-w-0">
                        <div className="text-[13px] font-semibold text-foreground mb-0.5 flex items-center gap-1">
                          {b.name}
                          <ArrowRight className="h-3 w-3 text-muted-foreground/0 -translate-x-1 transition-all group-hover:text-primary group-hover:translate-x-0" />
                        </div>
                        <p className="text-[12px] text-muted-foreground leading-relaxed m-0">{b.desc}</p>
                      </div>
                    </Link>
                  ))}
                </div>

                <Callout type="tip" title="One code generator drives all of it">
                  Every one of these batteries is available to any resource you scaffold.{' '}
                  <code>grit generate resource Product</code> emits the Go model, service, handler,
                  routes, Zod schema, TypeScript types, React Query hooks, and an admin page — all
                  type-safe and consistent, with the batteries above already available to it.
                </Callout>
              </div>

              {/* ── PART 3: How Grit compares ────────────────────── */}
              <div className="mb-12">
                <div className="flex items-center gap-3 mb-6">
                  <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 border border-primary/20">
                    <Server className="h-5 w-5 text-primary" />
                  </div>
                  <div>
                    <h2 className="text-2xl font-bold tracking-tight m-0">How Grit compares</h2>
                    <p className="text-sm text-muted-foreground m-0">Grit vs Laravel, Django, Next.js, Rails, T3</p>
                  </div>
                </div>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every framework below is excellent. The difference is what you have to add versus
                  what&apos;s already there. Grit&apos;s bet: a compiled Go backend, a generated admin
                  panel, and types shared across the language boundary — batteries in the core, not
                  in a marketplace of add-ons.
                </p>

                <div className="flex flex-wrap items-center gap-4 text-xs text-muted-foreground mb-4 not-prose">
                  <span className="inline-flex items-center gap-1.5"><CellMark v="yes" /> built-in</span>
                  <span className="inline-flex items-center gap-1.5"><CellMark v="partial" /> add-on / partial</span>
                  <span className="inline-flex items-center gap-1.5"><CellMark v="no" /> not applicable</span>
                </div>

                <div className="overflow-x-auto rounded-xl border border-border/40 mb-6 not-prose">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/40 bg-card/80">
                        <th className="px-4 py-3 text-left font-medium text-muted-foreground sticky left-0 bg-card/80 z-10 min-w-[240px]">Capability</th>
                        {FRAMEWORKS.map((f) => (
                          <th
                            key={f}
                            className={`px-3 py-3 text-center font-semibold text-xs whitespace-nowrap ${f === 'Grit' ? 'text-primary' : 'text-muted-foreground'}`}
                          >
                            {f}
                          </th>
                        ))}
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/30">
                      {MATRIX.map((row) => (
                        <tr key={row.capability} className="hover:bg-accent/20 transition-colors">
                          <td className="px-4 py-2.5 font-medium text-foreground/90 text-[13px] sticky left-0 bg-background z-10">
                            {row.capability}
                          </td>
                          {row.cells.map((c, i) => (
                            <td key={i} className="px-3 py-2.5 text-center">
                              <CellMark v={c} />
                            </td>
                          ))}
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                <p className="text-xs text-muted-foreground/70 mb-8 not-prose">
                  Ratings reflect what ships in a default install: e.g. Django&apos;s admin is core, but
                  Laravel&apos;s admin (Filament/Nova) and Rails&apos; (ActiveAdmin) are add-ons; Next.js and
                  T3 are single-language TypeScript stacks, so end-to-end types are native but there&apos;s no
                  compiled server binary. Add-ons can close most gaps — the column just shows what you get
                  before installing anything.
                </p>

                <Callout type="warning" title="When NOT to reach for Grit">
                  Grit is a two-language stack (Go + TypeScript). If your team is all-in on a single
                  language, a native option is a better fit — a pure-TypeScript app (Next.js / T3) or a
                  Python/PHP/Ruby team staying in Django / Laravel / Rails. Grit also assumes you want a
                  generated admin panel and the batteries above; if you need a bespoke architecture or a
                  runtime other than Go, the opinionated defaults will fight you more than they help. It
                  shines when you want one CLI to scaffold a hardened, type-safe Go + React product fast —
                  not when you want a blank slate.
                </Callout>
              </div>

              {/* Prev / Next */}
              <PageHelp
                faqs={[
                  {
                    q: 'Are these benchmark numbers hardware-normalised?',
                    a: 'No — the latency figures are illustrative local-machine results from the load-test walkthrough, not an independent benchmark. Run the k6 suite on your own hardware for real numbers.',
                  },
                  {
                    q: 'Do I have to use every battery?',
                    a: (
                      <>
                        No. They ship configured but are opt-in &mdash; use what you need, ignore the
                        rest. See <Link href="/docs/batteries">Batteries Included</Link>.
                      </>
                    ),
                  },
                  {
                    q: 'Is the comparison table fair?',
                    a: 'It compares built-in capabilities, not raw performance. Every framework here can be extended to match — Grit just ships more of it by default. The "when not to use Grit" note above is deliberate.',
                  },
                ]}
              />

              <div className="flex items-center justify-between border-t border-border pt-8 mt-12">
                <Button variant="ghost" asChild>
                  <Link href="/docs/getting-started/philosophy" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Philosophy
                  </Link>
                </Button>
                <Button variant="ghost" asChild>
                  <Link href="/docs/concepts/performance" className="gap-2">
                    Performance (deep dive)
                    <ArrowRight className="h-4 w-4" />
                  </Link>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
