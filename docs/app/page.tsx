import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, Github, Terminal, Layers, Zap, Shield, Database, Bot, Server, Smartphone, ChevronDown, Quote } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock } from '@/components/code-block'
import { SoftwareApplicationSchema, FAQPageSchema } from '@/components/structured-data'
import { HeroTerminal } from '@/components/hero-terminal'

export const metadata: Metadata = {
  title: 'Grit — Go + React Full-Stack Framework',
  description: 'Build production-ready full-stack applications with Go and React. One CLI, 5 architectures, batteries included.',
  alternates: { canonical: 'https://gritframework.dev' },
}

export default function HomePage() {
  return (
    <div className="min-h-screen bg-background">
      <SoftwareApplicationSchema />
      <FAQPageSchema />
      <SiteHeader />

      {/* ═══ HERO ═══ */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 -z-10">
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[1000px] h-[600px] bg-primary/[0.06] rounded-full blur-[120px]" />
        </div>

        <div className="max-w-6xl mx-auto flex flex-col items-center text-center pt-24 pb-20 md:pt-36 md:pb-28 px-6">
          <div className="flex items-center gap-2 rounded-full border border-primary/20 bg-primary/5 px-4 py-1.5 mb-8">
            <div className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse" />
            <span className="text-xs font-mono font-medium text-primary/80 tracking-wide">v3.21 — OFFLINE-FIRST · AUDIT CHAIN · FEATURE FLAGS</span>
          </div>

          <h1 className="text-5xl md:text-7xl font-bold tracking-tight text-foreground mb-6 leading-[1.1]">
            Build full-stack apps<br />
            <span className="bg-gradient-to-r from-primary via-blue-400 to-cyan-400 bg-clip-text text-transparent">
              with Go + React
            </span>
          </h1>

          <p className="text-lg md:text-xl text-muted-foreground max-w-2xl mb-10 leading-relaxed">
            One CLI scaffolds your entire application — Go API, React frontend, admin panel,
            auth, offline-first desktop, audit log, feature flags, webhook receiver, and more.
            Choose your architecture. Ship fast.
          </p>

          <div className="flex flex-col sm:flex-row items-center gap-4 mb-16">
            <Button size="lg" className="bg-primary hover:bg-primary/80 text-primary-foreground px-8 h-12 text-base rounded-full" asChild>
              <Link href="/docs/getting-started/quick-start">
                Get Started <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
            <Button variant="outline" size="lg" className="border-border/60 text-foreground hover:bg-accent/30 px-8 h-12 text-base rounded-full" asChild>
              <Link href="https://github.com/MUKE-coder/grit" target="_blank">
                <Github className="mr-2 h-4 w-4" /> GitHub
              </Link>
            </Button>
          </div>

          {/* Animated hero terminal */}
          <div className="w-full max-w-3xl">
            <HeroTerminal />
          </div>
        </div>
      </section>

      {/* ═══ STATS BAR ═══ */}
      <section className="border-y border-border/40 bg-card/50">
        <div className="max-w-6xl mx-auto px-6 py-12">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8 text-center stagger-children">
            {[
              { number: '5', label: 'Architecture Modes', sub: 'single, double, triple, api, mobile' },
              { number: '24', label: 'CLI Commands', sub: 'scaffold, generate, init, deploy...' },
              { number: '100+', label: 'UI Components', sub: 'shadcn-compatible registry' },
              { number: 'v3.21', label: 'Latest Release', sub: 'feature flags + A/B testing' },
            ].map((stat) => (
              <div key={stat.label}>
                <div className="text-4xl md:text-5xl font-bold text-foreground mb-1 animate-count">{stat.number}</div>
                <div className="text-sm font-medium text-foreground/80 mb-0.5">{stat.label}</div>
                <div className="text-xs text-muted-foreground">{stat.sub}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ═══ BENTO GRID ═══ */}
      <section className="py-24 px-6">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">Everything Included</p>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Batteries-included framework</h2>
            <p className="text-muted-foreground max-w-2xl mx-auto">
              Every Grit project ships with production-ready features out of the box.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 stagger-children">
            {/* Large card — Auth */}
            <div className="lg:col-span-2 rounded-xl border border-border/40 bg-card/50 card-gradient grid-pattern p-6 hover:border-primary/30 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5">
              <div className="flex items-center gap-3 mb-4">
                <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center icon-animated">
                  <Shield className="h-5 w-5 text-primary" />
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">Authentication + 2FA</h3>
                  <p className="text-xs text-muted-foreground">JWT, OAuth2, TOTP, backup codes, trusted devices</p>
                </div>
              </div>
              <CodeBlock language="bash" filename="Terminal" className="mb-0" code={`POST /api/auth/register    → JWT tokens
POST /api/auth/login       → JWT tokens (or totp_required + pending_token)
POST /api/auth/totp/verify  → Exchange TOTP code for JWT
GET  /api/auth/me          → Current user (protected)
GET  /api/auth/oauth/:provider → Google, GitHub social login`} />
            </div>

            {/* Small card — Admin */}
            <div className="rounded-xl border border-border/40 bg-card/50 card-gradient grid-pattern p-6 hover:border-primary/30 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5">
              <div className="h-10 w-10 rounded-lg bg-emerald-500/10 flex items-center justify-center mb-4 icon-animated">
                <Layers className="h-5 w-5 text-emerald-400" />
              </div>
              <h3 className="font-semibold text-foreground mb-2">Admin Panel</h3>
              <p className="text-sm text-muted-foreground mb-3">
                DataTable, FormBuilder, dashboard widgets, resource definitions.
              </p>
              <ul className="text-xs text-muted-foreground space-y-1">
                <li>- Sort, filter, paginate, search</li>
                <li>- 16+ form field types</li>
                <li>- Multi-step form wizards</li>
                <li>- 4 style variants</li>
              </ul>
            </div>

            {/* Small card — AI */}
            <div className="rounded-xl border border-border/40 bg-card/50 card-gradient grid-pattern p-6 hover:border-primary/30 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5">
              <div className="h-10 w-10 rounded-lg bg-violet-500/10 flex items-center justify-center mb-4 icon-animated">
                <Bot className="h-5 w-5 text-violet-400" />
              </div>
              <h3 className="font-semibold text-foreground mb-2">AI Gateway</h3>
              <p className="text-sm text-muted-foreground mb-3">
                One API key, hundreds of models via Vercel AI Gateway.
              </p>
              <code className="text-xs font-mono text-primary/80 bg-primary/5 px-2 py-1 rounded">
                anthropic/claude-sonnet-4-6
              </code>
            </div>

            {/* Small card — Storage */}
            <div className="rounded-xl border border-border/40 bg-card/50 card-gradient grid-pattern p-6 hover:border-primary/30 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5">
              <div className="h-10 w-10 rounded-lg bg-amber-500/10 flex items-center justify-center mb-4 icon-animated">
                <Database className="h-5 w-5 text-amber-400" />
              </div>
              <h3 className="font-semibold text-foreground mb-2">File Storage</h3>
              <p className="text-sm text-muted-foreground">
                Presigned URL uploads to S3, R2, or MinIO. Image processing. Progress tracking.
              </p>
            </div>

            {/* Large card — Code Gen */}
            <div className="lg:col-span-2 rounded-xl border border-border/40 bg-card/50 card-gradient grid-pattern p-6 hover:border-primary/30 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5">
              <div className="flex items-center gap-3 mb-4">
                <div className="h-10 w-10 rounded-lg bg-cyan-500/10 flex items-center justify-center icon-animated">
                  <Terminal className="h-5 w-5 text-cyan-400" />
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">Full-Stack Code Generation</h3>
                  <p className="text-xs text-muted-foreground">One command generates Go + React + admin in seconds</p>
                </div>
              </div>
              <CodeBlock language="bash" filename="Terminal" className="mb-0" code={`$ grit generate resource Product --fields "name:string,price:float,stock:int"

  ✓ internal/models/product.go
  ✓ internal/services/product.go
  ✓ internal/handlers/product.go
  ✓ apps/admin/src/routes/_dashboard/resources/products.tsx
  ✓ Injected model, handler, routes, resource registry

  ✅ Resource Product generated successfully!`} />
            </div>
          </div>
        </div>
      </section>

      {/* ═══ ARCHITECTURE ═══ */}
      <section className="py-24 px-6 border-t border-border/40 bg-card/30">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">Flexible Architecture</p>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Choose how you build</h2>
            <p className="text-muted-foreground max-w-2xl mx-auto">
              Coming from Laravel? Choose Single. MERN stack? Choose Double. Building a SaaS? Choose Triple.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4 stagger-children">
            {[
              { name: 'Single', icon: <Zap className="h-5 w-5" />, desc: 'Go + embedded SPA', flag: '--single', color: 'text-sky-400 bg-sky-400/10' },
              { name: 'Double', icon: <Layers className="h-5 w-5" />, desc: 'Web + API monorepo', flag: '--double', color: 'text-violet-400 bg-violet-400/10' },
              { name: 'Triple', icon: <Server className="h-5 w-5" />, desc: 'Web + Admin + API', flag: '--triple', color: 'text-emerald-400 bg-emerald-400/10' },
              { name: 'API Only', icon: <Database className="h-5 w-5" />, desc: 'Go backend only', flag: '--api', color: 'text-amber-400 bg-amber-400/10' },
              { name: 'Mobile', icon: <Smartphone className="h-5 w-5" />, desc: 'API + Expo', flag: '--mobile', color: 'text-rose-400 bg-rose-400/10' },
            ].map((arch) => (
              <div key={arch.name} className="rounded-xl border border-border/40 bg-card/50 card-gradient p-5 text-center hover:border-primary/30 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5">
                <div className={`h-12 w-12 rounded-xl ${arch.color} flex items-center justify-center mx-auto mb-3 icon-animated`}>
                  {arch.icon}
                </div>
                <h3 className="font-semibold text-foreground mb-1">{arch.name}</h3>
                <p className="text-xs text-muted-foreground mb-2">{arch.desc}</p>
                <code className="text-[10px] font-mono text-muted-foreground/60">{arch.flag}</code>
              </div>
            ))}
          </div>

          <div className="text-center mt-8">
            <Button variant="outline" className="border-border/60 text-foreground hover:bg-accent/30 rounded-full" asChild>
              <Link href="/docs/concepts/architecture-modes">
                Compare architectures <ArrowRight className="ml-2 h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>
        </div>
      </section>

      {/* ═══ DEPLOY DEEP-DIVE ═══ */}
      <section className="py-24 px-6 border-t border-border/40">
        <div className="max-w-6xl mx-auto">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            <div>
              <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">One-Command Deploy</p>
              <h2 className="text-3xl font-bold text-foreground mb-4">
                From code to production<br />in one command
              </h2>
              <p className="text-muted-foreground mb-6 leading-relaxed">
                <code className="text-primary bg-primary/5 px-1.5 py-0.5 rounded text-sm">grit deploy</code> builds
                your app, uploads via SSH, configures systemd, and sets up Caddy with auto-TLS.
              </p>
              <ul className="space-y-3 mb-8">
                {[
                  'Cross-compiles Go binary for Linux (CGO_ENABLED=0)',
                  'Builds frontend if present (pnpm build)',
                  'Uploads binary via SCP',
                  'Creates systemd service with auto-restart',
                  'Configures Caddy reverse proxy with Let\'s Encrypt TLS',
                ].map((item, i) => (
                  <li key={i} className="flex items-start gap-3 text-sm text-muted-foreground">
                    <span className="text-primary font-mono font-bold text-xs mt-0.5">{String(i + 1).padStart(2, '0')}</span>
                    {item}
                  </li>
                ))}
              </ul>
              <Button variant="outline" className="border-border/60 text-foreground hover:bg-accent/30 rounded-full" asChild>
                <Link href="/docs/infrastructure/deploy-command">
                  Deploy guide <ArrowRight className="ml-2 h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
            <div>
              <CodeBlock language="bash" filename="Terminal" code={`$ grit deploy --host deploy@server.com --domain myapp.com

  → Building frontend...
  → Building Go binary (linux/amd64)...
  → Uploading binary to /opt/myapp/
  → Setting up systemd service...
  → Configuring Caddy reverse proxy...

  ✓ Deployment successful!
  Live at: https://myapp.com`} />
            </div>
          </div>
        </div>
      </section>

      {/* ═══ BATTERIES GRID ═══ */}
      <section className="py-24 px-6 border-t border-border/40 bg-card/30">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">What Ships With Every Project</p>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Production-ready from day one</h2>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 stagger-children">
            {[
              { title: 'JWT Authentication', desc: 'Register, login, refresh, OAuth2 social login (Google, GitHub). Role-based access control.', color: 'text-sky-400', bg: 'bg-sky-400/10', icon: <svg className="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" /><path d="M7 11V7a5 5 0 0110 0v4" /><circle cx="12" cy="16" r="1" /></svg> },
              { title: 'Two-Factor Auth', desc: 'TOTP authenticator app, 10 backup codes, trusted devices with 30-day sliding cookie.', color: 'text-violet-400', bg: 'bg-violet-400/10', icon: <svg className="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" /><path d="M9 12l2 2 4-4" /></svg> },
              { title: 'File Storage', desc: 'Presigned URL uploads to S3, R2, or MinIO. Image processing. Progress tracking.', color: 'text-amber-400', bg: 'bg-amber-400/10', icon: <svg className="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" /><polyline points="17 8 12 3 7 8" /><line x1="12" y1="3" x2="12" y2="15" /></svg> },
              { title: 'Email (Resend)', desc: 'Transactional emails with Go HTML templates. Dev uses Mailhog.', color: 'text-pink-400', bg: 'bg-pink-400/10', icon: <svg className="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z" /><polyline points="22,6 12,13 2,6" /></svg> },
              { title: 'Background Jobs', desc: 'Redis-backed job queue via asynq. Image processing, email sending.', color: 'text-orange-400', bg: 'bg-orange-400/10', icon: <svg className="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" /></svg> },
              { title: 'Redis Cache', desc: 'Cache middleware for any route. Set/Get/Delete. Configurable TTL.', color: 'text-red-400', bg: 'bg-red-400/10', icon: <svg className="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><ellipse cx="12" cy="5" rx="9" ry="3" /><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3" /><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5" /></svg> },
            ].map((feature) => (
              <div key={feature.title} className="rounded-xl border border-border/40 bg-background card-gradient grid-pattern p-5 hover:border-primary/30 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5">
                <div className={`h-10 w-10 rounded-lg ${feature.bg} ${feature.color} flex items-center justify-center mb-4 icon-animated`}>
                  {feature.icon}
                </div>
                <h3 className="font-semibold text-foreground mb-1.5 text-sm">{feature.title}</h3>
                <p className="text-xs text-muted-foreground leading-relaxed">{feature.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ═══ OPS-GRADE PRIMITIVES ═══ */}
      <section className="py-24 px-6 border-t border-border/40">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">New in v3.10–v3.21</p>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Ops-grade primitives</h2>
            <p className="text-muted-foreground max-w-2xl mx-auto">
              Reliability, auditability, and rollout control — usually a year of integration work — wired into every scaffolded project.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {/* Offline-first */}
            <div className="lg:col-span-2 rounded-xl border border-border/40 bg-card/50 card-gradient grid-pattern p-6 hover:border-primary/30 transition-all duration-300">
              <div className="flex items-center gap-3 mb-4">
                <div className="h-10 w-10 rounded-lg bg-cyan-500/10 flex items-center justify-center">
                  <svg className="h-5 w-5 text-cyan-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M5 12.55a11 11 0 0 1 14.08 0" /><path d="M1.42 9a16 16 0 0 1 21.16 0" /><path d="M8.53 16.11a6 6 0 0 1 6.95 0" /><line x1="12" y1="20" x2="12.01" y2="20" /></svg>
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">Offline-first desktop sync</h3>
                  <p className="text-xs text-muted-foreground">Git-style: work locally → click Sync → resolve conflicts → push</p>
                </div>
              </div>
              <p className="text-sm text-muted-foreground leading-relaxed mb-4">
                Local SQLite mirror + outbox with squash semantics. Manual Sync button.
                Field-level conflict dialog when the server moved on. Versioned writes with
                optimistic-lock. <Link href="/docs/desktop/offline" className="text-primary hover:underline">Read the guide →</Link>
              </p>
              <code className="text-xs text-primary/80 font-mono bg-primary/5 px-2 py-1 rounded">
                grit new app --triple --vite --desktop
              </code>
            </div>

            {/* Audit + chain */}
            <div className="rounded-xl border border-border/40 bg-card/50 card-gradient grid-pattern p-6 hover:border-primary/30 transition-all duration-300">
              <div className="h-10 w-10 rounded-lg bg-emerald-500/10 flex items-center justify-center mb-4">
                <svg className="h-5 w-5 text-emerald-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M9 12l2 2 4-4" /><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3" /><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5" /><ellipse cx="12" cy="5" rx="9" ry="3" /></svg>
              </div>
              <h3 className="font-semibold text-foreground mb-2">Tamper-evident audit log</h3>
              <p className="text-sm text-muted-foreground mb-2">
                Every authenticated mutation auto-logged. SHA-256 hash chain proves the log
                wasn&apos;t edited via SQL.
              </p>
              <code className="text-[11px] text-primary/80 font-mono">GET /admin/activity/integrity</code>
            </div>

            {/* Feature flags */}
            <div className="rounded-xl border border-border/40 bg-card/50 card-gradient grid-pattern p-6 hover:border-primary/30 transition-all duration-300">
              <div className="h-10 w-10 rounded-lg bg-violet-500/10 flex items-center justify-center mb-4">
                <svg className="h-5 w-5 text-violet-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z" /><line x1="4" y1="22" x2="4" y2="15" /></svg>
              </div>
              <h3 className="font-semibold text-foreground mb-2">Feature flags + A/B</h3>
              <p className="text-sm text-muted-foreground mb-2">
                In-memory engine, sticky bucketing, percentage rollouts, allow/blocklists,
                realtime push when admin toggles.
              </p>
              <code className="text-[11px] text-primary/80 font-mono">flags.IsEnabled(c, &quot;new_ui&quot;)</code>
            </div>

            {/* Webhook receiver */}
            <div className="rounded-xl border border-border/40 bg-card/50 card-gradient grid-pattern p-6 hover:border-primary/30 transition-all duration-300">
              <div className="h-10 w-10 rounded-lg bg-amber-500/10 flex items-center justify-center mb-4">
                <svg className="h-5 w-5 text-amber-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M18 16.98h-5.99c-1.1 0-1.95.94-2.48 1.9A4 4 0 0 1 2 17c.01-.7.2-1.4.57-2" /><path d="m6 17 3.13-5.78c.53-.97.1-2.18-.5-3.1a4 4 0 1 1 6.89-4.06" /><path d="m12 6 3.13 5.73C15.66 12.7 16.9 13 18 13a4 4 0 0 1 0 8" /></svg>
              </div>
              <h3 className="font-semibold text-foreground mb-2">Webhook receiver</h3>
              <p className="text-sm text-muted-foreground mb-2">
                Stripe / GitHub / HMAC verifiers shipped. Auto-dedupe on (provider, external_id).
                Failed handlers replay from admin.
              </p>
              <code className="text-[11px] text-primary/80 font-mono">webhooks.On(&quot;stripe&quot;, ...)</code>
            </div>

            {/* Realtime + idempotency — wider card */}
            <div className="lg:col-span-2 rounded-xl border border-border/40 bg-card/50 card-gradient grid-pattern p-6 hover:border-primary/30 transition-all duration-300">
              <div className="flex items-center gap-3 mb-4">
                <div className="h-10 w-10 rounded-lg bg-rose-500/10 flex items-center justify-center">
                  <svg className="h-5 w-5 text-rose-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" /></svg>
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">Realtime hub + Idempotency-Key + PDF + CSV/Excel export</h3>
                  <p className="text-xs text-muted-foreground">All the &quot;every business app needs this&quot; primitives, baked in</p>
                </div>
              </div>
              <ul className="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-2 text-sm text-muted-foreground">
                <li>• <strong className="text-foreground/80">WebSocket fan-out hub</strong> at <code className="text-primary text-xs">GET /api/ws</code> — SendToUser / Broadcast helpers</li>
                <li>• <strong className="text-foreground/80">Idempotency-Key</strong> auto-attached on POST/PUT/PATCH/DELETE — 24h dedup</li>
                <li>• <strong className="text-foreground/80">PDF generation</strong> (<code className="text-primary text-xs">internal/pdf</code>) with worked Invoice template</li>
                <li>• <strong className="text-foreground/80">CSV / Excel export</strong> auto-emitted per resource (<code className="text-primary text-xs">/export?format=xlsx</code>)</li>
                <li>• <strong className="text-foreground/80">Cursor-based pagination</strong> opt-in for List endpoints — sticky pages</li>
                <li>• <strong className="text-foreground/80">Verbose AutoMigrate</strong> with column-level diff on every boot</li>
              </ul>
            </div>
          </div>
        </div>
      </section>

      {/* ═══ COMPARISON TABLE ═══ */}
      <section className="py-24 px-6 border-t border-border/40">
        <div className="max-w-5xl mx-auto">
          <div className="text-center mb-16">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">Framework Comparison</p>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">How Grit compares</h2>
          </div>

          <div className="overflow-x-auto rounded-xl border border-border/40">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border/40 bg-card/80">
                  <th className="px-5 py-4 text-left font-medium text-muted-foreground w-[200px]">Feature</th>
                  <th className="px-5 py-4 text-center font-semibold text-primary">Grit</th>
                  <th className="px-5 py-4 text-center font-medium text-muted-foreground">Next.js</th>
                  <th className="px-5 py-4 text-center font-medium text-muted-foreground">Laravel</th>
                  <th className="px-5 py-4 text-center font-medium text-muted-foreground hidden lg:table-cell">Goravel</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/30">
                {[
                  { feature: 'Go Backend', grit: true, next: false, laravel: false, goravel: true },
                  { feature: 'React Frontend', grit: true, next: true, laravel: false, goravel: false },
                  { feature: 'Admin Panel', grit: true, next: false, laravel: 'partial', goravel: false },
                  { feature: 'Code Generator', grit: true, next: false, laravel: true, goravel: true },
                  { feature: 'JWT + OAuth2', grit: true, next: false, laravel: true, goravel: true },
                  { feature: 'Two-Factor Auth', grit: true, next: false, laravel: false, goravel: false },
                  { feature: 'File Storage', grit: true, next: false, laravel: true, goravel: true },
                  { feature: 'Background Jobs', grit: true, next: false, laravel: true, goravel: true },
                  { feature: 'AI Integration', grit: true, next: false, laravel: false, goravel: false },
                  { feature: 'One-Command Deploy', grit: true, next: false, laravel: false, goravel: true },
                  { feature: 'Multiple Architectures', grit: true, next: false, laravel: false, goravel: false },
                  { feature: 'Desktop App', grit: true, next: false, laravel: false, goravel: false },
                  { feature: 'Offline-First Sync', grit: true, next: false, laravel: false, goravel: false },
                  { feature: 'Audit Log + Hash Chain', grit: true, next: false, laravel: false, goravel: false },
                  { feature: 'Feature Flags', grit: true, next: false, laravel: false, goravel: false },
                  { feature: 'Webhook Receiver', grit: true, next: false, laravel: 'partial', goravel: false },
                ].map((row) => (
                  <tr key={row.feature} className="hover:bg-accent/20 transition-colors">
                    <td className="px-5 py-3 font-medium text-foreground/90 text-[13px]">{row.feature}</td>
                    {[row.grit, row.next, row.laravel, row.goravel].map((val, i) => (
                      <td key={i} className={`px-5 py-3 text-center ${i === 3 ? 'hidden lg:table-cell' : ''}`}>
                        {val === true ? (
                          <span className={`inline-flex h-5 w-5 items-center justify-center rounded-full ${i === 0 ? 'bg-primary/15 text-primary' : 'bg-emerald-500/15 text-emerald-400'}`}>
                            <svg className="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}><path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" /></svg>
                          </span>
                        ) : val === 'partial' ? (
                          <span className="inline-flex h-5 w-5 items-center justify-center rounded-full bg-amber-500/15 text-amber-400">
                            <svg className="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}><path strokeLinecap="round" strokeLinejoin="round" d="M5 12h14" /></svg>
                          </span>
                        ) : (
                          <span className="inline-flex h-5 w-5 items-center justify-center rounded-full bg-muted/50 text-muted-foreground/30">
                            <svg className="h-2.5 w-2.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}><path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                          </span>
                        )}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </section>

      {/* ═══ TESTIMONIALS ═══ */}
      <section className="py-24 px-6 border-t border-border/40 bg-card/30">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">What Developers Say</p>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Loved by builders</h2>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {[
              { quote: 'Grit replaced our entire setup — Django backend, React frontend, separate admin panel. One CLI command and we had everything running in minutes.', name: 'Alex Chen', role: 'Senior Engineer', company: 'Startup Founder' },
              { quote: 'The code generator alone saves us 2-3 hours per feature. Generate a resource, get Go model + service + handler + admin page + React hooks. Incredible DX.', name: 'Sarah Johnson', role: 'Full-Stack Developer', company: 'Agency Lead' },
              { quote: 'Coming from Laravel, the single-app architecture felt instantly familiar. But with Go performance and type safety. The deploy command is chef\'s kiss.', name: 'Marcus Rodriguez', role: 'Backend Developer', company: 'Laravel to Go convert' },
            ].map((t) => (
              <div key={t.name} className="rounded-xl border border-border/40 bg-background card-gradient p-6">
                <Quote className="h-8 w-8 text-primary/20 mb-4" />
                <p className="text-sm text-muted-foreground leading-relaxed mb-6">{t.quote}</p>
                <div>
                  <div className="font-semibold text-foreground text-sm">{t.name}</div>
                  <div className="text-xs text-muted-foreground">{t.role} — {t.company}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ═══ SHOWCASE ═══ */}
      <section className="py-24 px-6 border-t border-border/40">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">Built With Grit</p>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Showcase</h2>
            <p className="text-muted-foreground max-w-2xl mx-auto">Projects and products built with the Grit framework.</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[
              { name: 'GritCMS', desc: 'Self-hostable creator platform. Website builder, email marketing, courses, community.', url: 'https://gritcms.com', tag: 'SaaS' },
              { name: 'Golang Battles', desc: 'Competitive Go coding platform with real-time WebSocket battles, ELO ranking, and sandbox execution.', url: '#', tag: 'Platform' },
              { name: 'GORM Studio', desc: 'Visual database browser for GORM. View tables, run queries, export data. Embedded in every Grit project.', url: 'https://github.com/MUKE-coder/gorm-studio', tag: 'Tool' },
              { name: 'Pulse', desc: 'Self-hosted observability SDK. Request tracing, DB monitoring, runtime metrics, Prometheus export.', url: 'https://github.com/MUKE-coder/pulse', tag: 'Library' },
              { name: 'Sentinel', desc: 'WAF + rate limiting + brute-force protection with real-time threat dashboard.', url: 'https://github.com/MUKE-coder/sentinel', tag: 'Security' },
              { name: 'gin-docs', desc: 'Zero-annotation API documentation generator for Gin. Auto-generates OpenAPI spec with Scalar UI.', url: 'https://github.com/MUKE-coder/gin-docs', tag: 'Library' },
            ].map((project) => (
              <Link key={project.name} href={project.url} target={project.url.startsWith('http') ? '_blank' : undefined} className="group rounded-xl border border-border/40 bg-card/50 card-gradient p-6 hover:border-primary/30 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5 block">
                <div className="flex items-center justify-between mb-3">
                  <h3 className="font-semibold text-foreground group-hover:text-primary transition-colors">{project.name}</h3>
                  <span className="text-[10px] font-mono text-muted-foreground/60 bg-accent/30 px-2 py-0.5 rounded">{project.tag}</span>
                </div>
                <p className="text-sm text-muted-foreground leading-relaxed">{project.desc}</p>
              </Link>
            ))}
          </div>
        </div>
      </section>

      {/* ═══ CREATOR QUOTE ═══ */}
      <section className="py-24 px-6 border-t border-border/40 bg-card/30">
        <div className="max-w-3xl mx-auto text-center">
          <img
            src="https://avatars.githubusercontent.com/u/64189841?v=4"
            alt="Muke JohnBaptist"
            className="h-16 w-16 rounded-full mx-auto mb-6 ring-2 ring-primary/20"
          />
          <blockquote className="text-xl md:text-2xl font-medium text-foreground leading-relaxed mb-6">
            &ldquo;I built Grit because I was tired of spending weeks setting up the same boilerplate for every project.
            Auth, admin panels, file uploads, background jobs — they should just work. Now they do.
            One command, and you have a production-ready app. That{"'"}s the framework I wanted to use.&rdquo;
          </blockquote>
          <div>
            <div className="font-semibold text-foreground">Muke JohnBaptist</div>
            <div className="text-sm text-muted-foreground">Creator of Grit Framework</div>
          </div>
        </div>
      </section>

      {/* ═══ FAQ ═══ */}
      <section className="py-24 px-6 border-t border-border/40">
        <div className="max-w-3xl mx-auto">
          <div className="text-center mb-16">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">FAQ</p>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Frequently asked questions</h2>
          </div>

          <div className="space-y-4">
            {[
              { q: 'Do I need to know Go to use Grit?', a: 'Basic Go knowledge helps, but Grit generates most of the code for you. The generated code follows clear patterns (handler → service → model) that are easy to extend. If you know any backend language, you\'ll pick it up fast.' },
              { q: 'Can I use Grit with an existing project?', a: 'Grit is designed for greenfield projects. It scaffolds the full project structure. However, you can use grit generate resource in existing Grit projects to add new features incrementally.' },
              { q: 'Is Grit production-ready?', a: 'Yes. Every scaffolded project includes JWT auth, RBAC, rate limiting (Sentinel), observability (Pulse), error handling, CORS, gzip compression, connection pooling, and graceful shutdown. It\'s designed for production from day one.' },
              { q: 'What\'s the difference between Single and Triple architecture?', a: 'Single embeds the React SPA into the Go binary via go:embed — one file to deploy. Triple is a Turborepo monorepo with separate web app, admin panel, and API — ideal for teams and complex products.' },
              { q: 'Can I switch from Next.js to TanStack Router later?', a: 'The backend (Go API) is identical regardless of frontend choice. You\'d need to rebuild the frontend pages, but all hooks, types, and API patterns are the same. The admin panel components are also framework-agnostic React.' },
              { q: 'How does grit deploy work? Is it like Vercel?', a: 'grit deploy is for self-hosted deployments. It SSHs to your server, uploads the binary, configures systemd, and sets up Caddy with auto-TLS. For Vercel/Railway, just push to git — the Dockerfile is included.' },
              { q: 'Is Grit open source?', a: 'Yes, Grit is fully open source under the MIT license. The CLI, all plugins, and the documentation are on GitHub.' },
            ].map((faq) => (
              <details key={faq.q} className="group rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                <summary className="flex items-center justify-between px-6 py-4 cursor-pointer list-none">
                  <span className="font-medium text-foreground text-sm pr-4">{faq.q}</span>
                  <ChevronDown className="h-4 w-4 text-muted-foreground shrink-0 transition-transform group-open:rotate-180" />
                </summary>
                <div className="px-6 pb-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">{faq.a}</p>
                </div>
              </details>
            ))}
          </div>
        </div>
      </section>

      {/* ═══ CTA ═══ */}
      <section className="py-24 px-6 border-t border-border/40 bg-card/30">
        <div className="max-w-3xl mx-auto text-center">
          <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Start building in 30 seconds</h2>
          <p className="text-muted-foreground mb-8 text-lg">Install the CLI and scaffold your first project.</p>
          <CodeBlock language="bash" className="mb-8 text-left" code={`go install github.com/MUKE-coder/grit/v3/cmd/grit@latest
grit new my-app`} />
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Button size="lg" className="bg-primary hover:bg-primary/80 text-primary-foreground px-8 h-12 text-base rounded-full" asChild>
              <Link href="/docs/getting-started/quick-start">Read the docs <ArrowRight className="ml-2 h-4 w-4" /></Link>
            </Button>
            <Button variant="outline" size="lg" className="border-border/60 text-foreground hover:bg-accent/30 px-8 h-12 text-base rounded-full" asChild>
              <Link href="https://github.com/MUKE-coder/grit" target="_blank"><Github className="mr-2 h-4 w-4" /> Star on GitHub</Link>
            </Button>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-border/40 py-8 px-6">
        <div className="max-w-6xl mx-auto flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2">
            <div className="flex h-6 w-6 items-center justify-center rounded-md bg-primary/10">
              <span className="text-primary font-mono font-bold text-[10px]">G</span>
            </div>
            <span className="text-sm text-muted-foreground">Grit Framework</span>
          </div>
          <div className="flex items-center gap-6">
            <Link href="/docs" className="text-xs text-muted-foreground hover:text-foreground transition-colors">Docs</Link>
            <Link href="https://github.com/MUKE-coder/grit" target="_blank" className="text-xs text-muted-foreground hover:text-foreground transition-colors">GitHub</Link>
            <Link href="/showcase" className="text-xs text-muted-foreground hover:text-foreground transition-colors">Showcase</Link>
            <Link href="/hire" className="text-xs text-muted-foreground hover:text-foreground transition-colors">Hire Us</Link>
          </div>
          <p className="text-xs text-muted-foreground/60">Built with Grit v3.5</p>
        </div>
      </footer>
    </div>
  )
}
