import Link from 'next/link'
import { ArrowLeft, ArrowRight, ExternalLink, Globe, Github, Sparkles, ShieldCheck, Activity, Bike, Receipt, HandCoins, Users } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/demo')

const FEATURE_TOURS: { icon: any; title: string; body: string; href: string }[] = [
  {
    icon: Receipt,
    title: 'Full POS + sales',
    body: "Cart, discount + tax math, multi-payment-method (cash / mobile money / credit), printable receipt, returns. Real business logic — not a todo list.",
    href: '/pos',
  },
  {
    icon: HandCoins,
    title: 'Loans + repayments',
    body: "Loan products, borrower profiles, payment schedules, overdue flagging cron, mobile-money collection via DGateway. The kind of domain logic Grit is built for.",
    href: '/loans/new',
  },
  {
    icon: Bike,
    title: 'Motorcycle inventory + Daily Boda',
    body: "Track motorcycles by chassis number, assign to riders, capture daily earnings. Shows off image uploads, presigned URLs, list filtering.",
    href: '/motorcycles',
  },
  {
    icon: Users,
    title: 'Multi-tenant + RBAC',
    body: 'Workspaces (Loans · Spares), per-business roles, invitation flow, audit log on every mutation. The auth + tenancy primitives every SaaS rebuilds.',
    href: '/staff',
  },
  {
    icon: ShieldCheck,
    title: 'In-app Security dashboard',
    body: "Sentinel summary inside the admin: security score, recent threats with CVSS, AuthShield state, CSP violations. Source visible — copy the pattern into your own app.",
    href: '/system/security',
  },
  {
    icon: Activity,
    title: 'In-app Observability dashboard',
    body: "Pulse summary inside the admin: p95/p99 latency, SLO bars, USE method grid, top N+1 by impact, error feed. Same pattern as Security.",
    href: '/system/observability',
  },
]

export default function DemoPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Demo</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Demo Application
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                A live, working full-stack Grit app — kitchen-sink edition. Built with
                {' '}<code className="text-primary text-sm bg-primary/5 px-1.5 py-0.5 rounded">grit new --single --vite</code>{' '}
                then extended into a real motorcycle dealership management system with
                POS, loans, inventory, multi-tenant RBAC, and live Security + Observability
                dashboards. Source code lives next to the framework so you can copy the
                patterns straight into your own project.
              </p>
            </div>

            {/* Quick links — live + source + login pill */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 mb-10">
              <Button
                size="lg"
                className="bg-primary hover:bg-primary/90 text-primary-foreground h-12 text-base rounded-xl"
                asChild
              >
                <Link href="https://demo.gritframework.dev" target="_blank" rel="noopener noreferrer">
                  <Globe className="mr-2 h-4 w-4" /> Open live demo
                </Link>
              </Button>
              <Button
                size="lg"
                variant="outline"
                className="border-border/60 h-12 text-base rounded-xl"
                asChild
              >
                <Link href="https://github.com/MUKE-coder/grit/tree/main/demo" target="_blank" rel="noopener noreferrer">
                  <Github className="mr-2 h-4 w-4" /> Source on GitHub
                </Link>
              </Button>
              <Button
                size="lg"
                variant="outline"
                className="border-border/60 h-12 text-base rounded-xl"
                asChild
              >
                <Link href="#self-host">
                  <Sparkles className="mr-2 h-4 w-4" /> Self-host
                </Link>
              </Button>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 mb-12 flex gap-3 items-start">
              <ShieldCheck className="h-5 w-5 text-primary mt-0.5 shrink-0" />
              <div className="text-sm text-foreground/85">
                <p className="mb-1"><strong>The login form is pre-filled.</strong> Credentials:</p>
                <p className="font-mono text-xs">
                  <span className="text-foreground">admin@grit.demo</span>
                  {' / '}
                  <span className="text-foreground">password123</span>
                </p>
                <p className="text-xs text-muted-foreground mt-2">
                  The database is reset every midnight (local time). Please be respectful when editing data.
                  Outbound email is disabled while in <code className="text-xs px-1 py-0.5 rounded bg-background">DEMO_MODE=true</code>.
                </p>
              </div>
            </div>

            <div className="prose-grit">

              <h2 id="what-it-shows">What the demo shows</h2>
              <p>
                Inertia&apos;s demo is a CRM kitchen sink. Ours is a motorcycle dealership
                + loan-financing app — a richer domain that exercises more Grit primitives
                in their natural habitat: server-side pagination, optimistic mutations,
                file uploads, role gating, audit logging, async cron, webhook receivers,
                and the new in-app Security &amp; Observability summary pages.
              </p>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 my-6 not-prose">
                {FEATURE_TOURS.map(({ icon: Icon, title, body, href }) => (
                  <div key={title} className="rounded-xl border border-border/50 bg-card/50 p-4">
                    <div className="flex items-center gap-2.5 mb-2">
                      <Icon className="h-4 w-4 text-primary" />
                      <h3 className="font-semibold text-foreground text-sm">{title}</h3>
                    </div>
                    <p className="text-sm text-muted-foreground leading-relaxed mb-3">{body}</p>
                    <a
                      href={`https://demo.gritframework.dev${href}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-1 text-xs font-mono text-primary hover:underline"
                    >
                      {href} <ExternalLink className="h-3 w-3" />
                    </a>
                  </div>
                ))}
              </div>

              <h2 id="grit-features-in-use">Grit features wired in the demo</h2>
              <p>
                Each of these maps to a real file in <code>demo/</code> so you can read the
                source after seeing the behaviour in the live app:
              </p>
              <ul>
                <li><strong>Auth + JWT + RBAC</strong> — <code>internal/handlers/auth.go</code>, role gates in <code>internal/middleware/auth.go</code></li>
                <li><strong>Multi-tenant business switcher</strong> — <code>internal/middleware/business_scope.go</code> + the workspace switcher in <code>frontend/src/components/layout/AppSidebar.tsx</code></li>
                <li><strong>File uploads</strong> — presigned S3/R2 in <code>internal/handlers/uploads.go</code>, used by motorcycle photos</li>
                <li><strong>Background cron</strong> — overdue-installment flagger + (in demo mode) nightly DB reset in <code>internal/cron/cron.go</code></li>
                <li><strong>Audit log</strong> — every mutation streams into <code>activities</code> via <code>internal/services/audit.go</code></li>
                <li><strong>Webhook receiver</strong> — DGateway mobile-money callbacks (HMAC verify + dedup) in <code>internal/handlers/dgateway.go</code></li>
                <li><strong>Sentinel v2.0 + Pulse v1.0</strong> — mounted at <code>/sentinel</code> and <code>/pulse</code> in <code>internal/routes/routes.go</code></li>
                <li><strong>In-app Security + Observability</strong> — <code>internal/handlers/security.go</code> + <code>observability.go</code>, the <code>SecObsBridge</code> for loopback API calls, and <code>SecObsPoller</code> turning findings into bell notifications</li>
              </ul>

              <h2 id="self-host">Run it locally / self-host</h2>
              <p>
                The demo is just a Grit project — clone, set env, run.
              </p>
              <CodeBlock language="bash" filename="Terminal" code={`# Clone
git clone https://github.com/MUKE-coder/grit.git
cd grit/demo

# Configure (DEMO_MODE=true ships in .env.example)
cp .env.example .env
# Edit DATABASE_URL — point at your Postgres (Neon, Supabase, local).
# Set JWT_SECRET to anything long + random.

# Frontend bundle
cd frontend && pnpm install && pnpm build && cd ..

# Run
go run .`} />
              <p>
                The seeder runs on every boot — first run creates{' '}
                <code>admin@grit.demo / password123</code> plus a <code>Grit Motors</code>
                business and a default branch. On subsequent boots it&apos;s a no-op unless
                something is missing.
              </p>

              <h2 id="demo-mode-behaviour">What <code>DEMO_MODE=true</code> changes</h2>
              <ul>
                <li><strong>Resend mailer is not constructed.</strong> Invitation + low-stock alert emails are silently skipped (the call sites are already <code>nil</code>-guarded).</li>
                <li><strong>Nightly reset cron registers.</strong> At 00:00 local time, <code>WipeMutableData</code> truncates sales / loans / payments / activities / notifications, then <code>Seed()</code> re-runs to restore the canonical admin + business + branch.</li>
                <li><strong>Auth credentials are pre-filled.</strong> The login form ships with <code>admin@grit.demo</code> + <code>password123</code> in <code>frontend/src/routes/_auth/login.tsx</code> so visitors can hit Enter.</li>
                <li><strong>A demo banner</strong> sits at the top of the login form explaining the reset cadence.</li>
              </ul>

              <h2 id="security-observability-pages">Security &amp; Observability pages</h2>
              <p>
                Two of the most interesting things to study in this demo are{' '}
                <Link href="/docs/security">/system/security</Link> and{' '}
                <Link href="/docs/testing">/system/observability</Link> — the same pattern
                every Grit project gets when scaffolded. They use a small
                <code> SecObsBridge</code> that logs in to the locally-mounted Sentinel /
                Pulse APIs over loopback, caches the JWT, and proxies the dashboard
                summary endpoints into one envelope. The React pages render KPI scorecards,
                live threat / SLO tables, and Brendan-Gregg USE-method grids — and a deep
                link out to the full <code>/sentinel/ui</code> or <code>/pulse/ui</code>
                for when you want to dig further. A <code>SecObsPoller</code> runs once a
                minute and writes high-severity findings into the in-app notification feed,
                so the bell icon at the top of the admin lights up when something firing.
              </p>

              <h2 id="what-its-not">What this demo isn&apos;t</h2>
              <ul>
                <li><strong>Not the only way to use Grit.</strong> It&apos;s a single-app (<code>--single --vite</code>) deployment; Grit also scaffolds triple-app monorepos, mobile (<code>--mobile</code>), and offline-first desktop (<code>--desktop</code>).</li>
                <li><strong>Not a starter template.</strong> If you&apos;re building a motorcycle dealership, copy from here. If you&apos;re building anything else, run <code>grit new</code> and pull selectively.</li>
                <li><strong>Not production-grade billing.</strong> The DGateway integration is real but the demo&apos;s sandbox keys make it safe to click through.</li>
              </ul>

            </div>

            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/testing" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Testing
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/changelog" className="gap-1.5">
                  Changelog
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
