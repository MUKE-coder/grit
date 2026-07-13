import Link from 'next/link'
import {
  ArrowRight,
  KeyRound,
  ShieldCheck,
  HardDrive,
  Mail,
  ListChecks,
  Clock,
  DatabaseZap,
  Sparkles,
  ShieldAlert,
  Activity,
  Database,
  Flag,
  Webhook,
  RadioTower,
  type LucideIcon,
} from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { PageHelp } from '@/components/page-help'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/batteries')

interface Battery {
  title: string
  blurb: string
  href: string
  icon: LucideIcon
  file: string
  code: string
}

const BATTERIES: Battery[] = [
  {
    title: 'Authentication',
    blurb: 'JWT register/login, bcrypt hashing and protected routes — /api/auth/* is live on day one.',
    href: '/docs/backend/authentication',
    icon: KeyRound,
    file: 'auth.http',
    code: 'POST /api/auth/register\nPOST /api/auth/login\n// -> { "token": "eyJ...", "user": {...} }',
  },
  {
    title: 'RBAC & Roles',
    blurb: 'Uppercase roles and a RequireRole guard — lock any route to ADMIN, EDITOR or USER.',
    href: '/docs/backend/rbac',
    icon: ShieldCheck,
    file: 'routes.go',
    code: 'admin := r.Group("/api/admin")\nadmin.Use(middleware.RequireRole("ADMIN"))\nadmin.GET("/users", h.ListUsers)',
  },
  {
    title: 'File Storage',
    blurb: 'S3-compatible storage across MinIO, R2 and AWS S3. Hand the client a presigned upload URL.',
    href: '/docs/batteries/storage',
    icon: HardDrive,
    file: 'upload.go',
    code: 'url, err := storage.PresignPut(ctx, "avatars/u1.png")\nif err != nil {\n  return err\n}\n// client PUTs the file straight to the bucket',
  },
  {
    title: 'Email (Resend)',
    blurb: 'Transactional email through Resend with typed HTML templates for welcome, reset and verify.',
    href: '/docs/batteries/email',
    icon: Mail,
    file: 'mailer.go',
    code: 'mailer.Send(ctx, mail.Message{\n  To:       user.Email,\n  Template: mail.Welcome,\n  Data:     map[string]any{"Name": user.Name},\n})',
  },
  {
    title: 'Background Jobs',
    blurb: 'Redis-backed async queue (asynq). Enqueue work off the request path and process it in a worker.',
    href: '/docs/batteries/jobs',
    icon: ListChecks,
    file: 'handler.go',
    code: 'jobs.Enqueue(ctx, jobs.SendEmail{\n  UserID: user.ID, // UUID string\n  Kind:   "welcome",\n})',
  },
  {
    title: 'Cron Scheduler',
    blurb: 'Register recurring tasks with plain cron expressions — monitored from the admin dashboard.',
    href: '/docs/batteries/cron',
    icon: Clock,
    file: 'scheduler.go',
    code: 'scheduler.Register("0 * * * *", func(ctx context.Context) error {\n  return report.GenerateHourly(ctx)\n})',
  },
  {
    title: 'Redis Caching',
    blurb: 'A tiny cache service plus response middleware. Cache hot reads with a TTL, invalidate on write.',
    href: '/docs/batteries/caching',
    icon: DatabaseZap,
    file: 'service.go',
    code: 'if v, ok := cache.Get(ctx, key); ok {\n  return v, nil\n}\ncache.Set(ctx, key, stats, 5*time.Minute)',
  },
  {
    title: 'AI',
    blurb: 'Stream completions through the Vercel AI Gateway — swap models with one string, no SDK churn.',
    href: '/docs/batteries/ai',
    icon: Sparkles,
    file: 'ai.go',
    code: 'stream, err := ai.Stream(ctx, ai.Request{\n  Model:  "anthropic/claude-sonnet-4-6",\n  Prompt: "Summarise this ticket",\n})',
  },
  {
    title: 'Sentinel',
    blurb: 'WAF, rate limiting, brute-force and anomaly detection — one mount and every route is shielded.',
    href: '/docs/batteries/security',
    icon: ShieldAlert,
    file: 'main.go',
    code: 'if err := sentinel.MountE(r, db, sentinel.Config{\n  RateLimit: 100, // req/min per IP\n}); err != nil {\n  log.Fatal(err)\n}',
  },
  {
    title: 'Pulse',
    blurb: 'Self-hosted observability: request tracing, DB monitoring, runtime metrics and a live dashboard.',
    href: '/docs/backend/pulse',
    icon: Activity,
    file: 'main.go',
    code: 'pulse.Mount(ctx, r, db, pulse.Config{\n  Path: "/pulse",\n})\n// dashboard live at /pulse',
  },
  {
    title: 'GORM Studio',
    blurb: 'A visual database browser embedded in the API. Inspect and edit rows without leaving the app.',
    href: '/docs/infrastructure/database',
    icon: Database,
    file: 'terminal',
    code: '$ grit studio\n  ✓ GORM Studio running at http://localhost:8080/studio',
  },
  {
    title: 'Feature Flags',
    blurb: 'Ship behind flags with percentage rollouts, allow/block lists and sticky per-user bucketing.',
    href: '/docs/backend/feature-flags',
    icon: Flag,
    file: 'handler.go',
    code: 'if flags.IsEnabled(c, "new-checkout") {\n  return h.newCheckout(c)\n}\nreturn h.legacyCheckout(c)',
  },
  {
    title: 'Webhooks',
    blurb: 'One inbound endpoint with built-in Stripe, GitHub and HMAC verifiers plus idempotent dedupe.',
    href: '/docs/backend/webhooks',
    icon: Webhook,
    file: 'webhooks.go',
    code: 'hooks.Register("stripe", webhook.StripeVerifier(secret))\n// POST /webhooks/:provider  ->  verified & deduped',
  },
  {
    title: 'Realtime',
    blurb: 'A WebSocket hub with JWT handshake auth. Push events to one user or broadcast to everyone.',
    href: '/docs/backend/realtime',
    icon: RadioTower,
    file: 'hub.go',
    code: 'hub.SendToUser(user.ID, realtime.Event{\n  Type: "job.finished",\n  Data: map[string]any{"id": job.ID},\n})',
  },
]

const FAQS = [
  {
    q: 'Do I have to use all of these?',
    a: (
      <>
        No. Every battery is opt-in at the call site &mdash; if you never call{' '}
        <code>mailer.Send</code> or mount Pulse, it simply does nothing. Nothing here forces a
        pattern on the rest of your app; unused batteries add no runtime cost and can be deleted
        outright.
      </>
    ),
  },
  {
    q: 'Can I swap the storage / email / AI provider?',
    a: (
      <>
        Yes &mdash; that&apos;s the point of the wrappers. Storage speaks the S3 API, so MinIO in
        dev and Cloudflare R2 or AWS S3 in prod are a config change. Email is Resend by default but
        the mailer is an interface. AI routes through the Vercel AI Gateway, so switching from{' '}
        <code>anthropic/claude-sonnet-4-6</code> to another model is a one-string edit.
      </>
    ),
  },
  {
    q: 'Are these external services or built into Grit?',
    a: (
      <>
        Both, depending on the battery. Auth, RBAC, caching, jobs, cron, Sentinel, Pulse, GORM
        Studio, feature flags, webhooks and realtime are pure Go that ships inside your binary. File
        storage, email and AI are thin, swappable clients in front of services you bring (an S3
        bucket, a Resend key, an AI Gateway key). No proprietary Grit cloud is required.
      </>
    ),
  },
  {
    q: 'What do the batteries cost me in bundle size or dependencies?',
    a: (
      <>
        Nothing on the frontend &mdash; these are backend features, so your Next.js bundle is
        untouched. On the Go side they compile into one static binary, and unused packages are
        eliminated by the compiler. The heaviest shared dependency is Redis, and only jobs, cron and
        caching need it.
      </>
    ),
  },
  {
    q: 'Do the batteries work in every architecture mode?',
    a: (
      <>
        Yes. Because they live in the Go API, they&apos;re available whether you scaffold the full
        triple stack, a single embedded binary, an API-only backend, or a mobile/desktop client
        talking to a Grit API. See{' '}
        <Link href="/docs/concepts/architecture-modes">architecture modes</Link> for the details.
      </>
    ),
  },
]

export default function BatteriesPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-5xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Batteries Included</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Everything, out of the box</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit is <em>batteries included</em> &mdash; auth, storage, email, jobs, AI, security
                and observability are wired into every scaffolded project. No plugins to install, no
                glue to write. Here is the full set, each with the one line you actually call.
              </p>
            </div>

            <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-2">
              {BATTERIES.map((b) => {
                const Icon = b.icon
                return (
                  <Link
                    key={b.href}
                    href={b.href}
                    className="group flex flex-col rounded-xl border border-border/50 bg-card/40 overflow-hidden hover:border-primary/30 transition-colors"
                  >
                    {/* Code mockup */}
                    <div className="relative border-b border-border/40 bg-[#0d1117] px-4 py-3.5 overflow-hidden">
                      <div className="flex items-center gap-1.5 mb-2">
                        <span className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                        <span className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                        <span className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                        <span className="ml-2 text-[10px] font-mono tracking-wider text-primary/50">
                          {b.file}
                        </span>
                      </div>
                      <pre className="text-[11px] leading-5 font-mono text-slate-300/80 whitespace-pre-wrap">
                        {b.code}
                      </pre>
                    </div>
                    {/* Body */}
                    <div className="flex flex-1 flex-col p-5">
                      <h2 className="text-base font-semibold mb-1.5 flex items-center gap-2 group-hover:text-primary transition-colors">
                        <Icon className="h-4 w-4 text-primary/80 shrink-0" />
                        {b.title}
                      </h2>
                      <p className="text-sm text-muted-foreground/70 leading-relaxed">{b.blurb}</p>
                      <span className="mt-4 inline-flex items-center gap-1.5 text-sm font-medium text-primary/90">
                        Read the docs
                        <ArrowRight className="h-3.5 w-3.5 -translate-x-0.5 group-hover:translate-x-0 transition-transform" />
                      </span>
                    </div>
                  </Link>
                )
              })}
            </div>

            <PageHelp faqs={FAQS} />
          </div>
        </div>
      </main>
    </div>
  )
}
