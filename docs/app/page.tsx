import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, Github, Terminal, Layers, Zap, Shield, Database, Bot, Server, Monitor, Smartphone } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock } from '@/components/code-block'
import { SoftwareApplicationSchema, FAQPageSchema } from '@/components/structured-data'

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

      {/* ═══════════════════════════════════════════════════════
          HERO
      ═══════════════════════════════════════════════════════ */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 -z-10">
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[1000px] h-[600px] bg-primary/[0.06] rounded-full blur-[120px]" />
        </div>

        <div className="max-w-6xl mx-auto flex flex-col items-center text-center pt-24 pb-20 md:pt-36 md:pb-28 px-6">
          <div className="flex items-center gap-2 rounded-full border border-primary/20 bg-primary/5 px-4 py-1.5 mb-8">
            <div className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse" />
            <span className="text-xs font-mono font-medium text-primary/80 tracking-wide">v3.5 — 5 ARCHITECTURES, 2 FRONTENDS</span>
          </div>

          <h1 className="text-5xl md:text-7xl font-bold tracking-tight text-foreground mb-6 leading-[1.1]">
            Build full-stack apps<br />
            <span className="bg-gradient-to-r from-primary via-blue-400 to-cyan-400 bg-clip-text text-transparent">
              with Go + React
            </span>
          </h1>

          <p className="text-lg md:text-xl text-muted-foreground max-w-2xl mb-10 leading-relaxed">
            One CLI scaffolds your entire application — Go API, React frontend, admin panel,
            auth, file uploads, background jobs, and more. Choose your architecture. Ship fast.
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

          {/* Hero code snippet */}
          <div className="w-full max-w-3xl">
            <CodeBlock language="bash" filename="Terminal" code={`$ grit new my-saas --triple --vite

  ██████╗ ██████╗ ██╗████████╗
 ██╔════╝ ██╔══██╗██║╚══██╔══╝
 ██║  ███╗██████╔╝██║   ██║
 ██║   ██║██╔══██╗██║   ██║
 ╚██████╔╝██║  ██║██║   ██║
  ╚═════╝ ╚═╝  ╚═╝╚═╝   ╚═╝

  Creating new Grit project: my-saas
  Architecture: triple | Frontend: tanstack

  → Creating directory structure...
  → Scaffolding Go API...
  → Adding batteries (cache, storage, mail, jobs, AI, TOTP)...
  → Scaffolding TanStack Router web app (Vite)...
  → Scaffolding TanStack Router admin panel (Vite)...

  ✓ Project created successfully!`} />
          </div>
        </div>
      </section>

      {/* ═══════════════════════════════════════════════════════
          STATS BAR
      ═══════════════════════════════════════════════════════ */}
      <section className="border-y border-border/40 bg-card/50">
        <div className="max-w-6xl mx-auto px-6 py-12">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
            {[
              { number: '5', label: 'Architecture Modes', sub: 'single, double, triple, api, mobile' },
              { number: '21', label: 'CLI Commands', sub: 'scaffold, generate, deploy, routes...' },
              { number: '100+', label: 'UI Components', sub: 'shadcn-compatible registry' },
              { number: '10', label: 'Official Plugins', sub: 'websockets, stripe, oauth...' },
            ].map((stat) => (
              <div key={stat.label}>
                <div className="text-4xl md:text-5xl font-bold text-foreground mb-1">{stat.number}</div>
                <div className="text-sm font-medium text-foreground/80 mb-0.5">{stat.label}</div>
                <div className="text-xs text-muted-foreground">{stat.sub}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ═══════════════════════════════════════════════════════
          BENTO GRID — FEATURES
      ═══════════════════════════════════════════════════════ */}
      <section className="py-24 px-6">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">Everything Included</p>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Batteries-included framework</h2>
            <p className="text-muted-foreground max-w-2xl mx-auto">
              Every Grit project ships with production-ready features out of the box.
              No boilerplate, no third-party setup. Just build.
            </p>
          </div>

          {/* Bento Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {/* Large card — Auth */}
            <div className="lg:col-span-2 rounded-xl border border-border/40 bg-card/50 p-6 hover:border-primary/20 transition-colors">
              <div className="flex items-center gap-3 mb-4">
                <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center">
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
            <div className="rounded-xl border border-border/40 bg-card/50 p-6 hover:border-primary/20 transition-colors">
              <div className="h-10 w-10 rounded-lg bg-emerald-500/10 flex items-center justify-center mb-4">
                <Layers className="h-5 w-5 text-emerald-400" />
              </div>
              <h3 className="font-semibold text-foreground mb-2">Admin Panel</h3>
              <p className="text-sm text-muted-foreground mb-3">
                DataTable, FormBuilder, dashboard widgets, resource definitions.
                Zero-code CRUD from a single config.
              </p>
              <ul className="text-xs text-muted-foreground space-y-1">
                <li>- Sort, filter, paginate, search</li>
                <li>- 16+ form field types</li>
                <li>- Multi-step form wizards</li>
                <li>- 4 style variants</li>
              </ul>
            </div>

            {/* Small card — AI */}
            <div className="rounded-xl border border-border/40 bg-card/50 p-6 hover:border-primary/20 transition-colors">
              <div className="h-10 w-10 rounded-lg bg-violet-500/10 flex items-center justify-center mb-4">
                <Bot className="h-5 w-5 text-violet-400" />
              </div>
              <h3 className="font-semibold text-foreground mb-2">AI Gateway</h3>
              <p className="text-sm text-muted-foreground mb-3">
                One API key, hundreds of models via Vercel AI Gateway.
                Streaming responses, chat, completions.
              </p>
              <code className="text-xs font-mono text-primary/80 bg-primary/5 px-2 py-1 rounded">
                anthropic/claude-sonnet-4-6
              </code>
            </div>

            {/* Small card — Storage */}
            <div className="rounded-xl border border-border/40 bg-card/50 p-6 hover:border-primary/20 transition-colors">
              <div className="h-10 w-10 rounded-lg bg-amber-500/10 flex items-center justify-center mb-4">
                <Database className="h-5 w-5 text-amber-400" />
              </div>
              <h3 className="font-semibold text-foreground mb-2">File Storage</h3>
              <p className="text-sm text-muted-foreground">
                Presigned URL uploads to S3, R2, or MinIO.
                Image processing, progress tracking. No multipart.
              </p>
            </div>

            {/* Large card — Code Gen */}
            <div className="lg:col-span-2 rounded-xl border border-border/40 bg-card/50 p-6 hover:border-primary/20 transition-colors">
              <div className="flex items-center gap-3 mb-4">
                <div className="h-10 w-10 rounded-lg bg-cyan-500/10 flex items-center justify-center">
                  <Terminal className="h-5 w-5 text-cyan-400" />
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">Full-Stack Code Generation</h3>
                  <p className="text-xs text-muted-foreground">One command generates Go + React + admin in seconds</p>
                </div>
              </div>
              <CodeBlock language="bash" filename="Terminal" className="mb-0" code={`$ grit generate resource Product --fields "name:string,price:float,stock:int,active:bool"

  ✓ internal/models/product.go
  ✓ internal/services/product.go
  ✓ internal/handlers/product.go
  ✓ apps/admin/src/routes/_dashboard/resources/products.tsx
  ✓ apps/web/src/hooks/use-products.ts
  ✓ Injected model, handler, routes, resource registry

  ✅ Resource Product generated successfully!`} />
            </div>
          </div>
        </div>
      </section>

      {/* ═══════════════════════════════════════════════════════
          ARCHITECTURE SECTION
      ═══════════════════════════════════════════════════════ */}
      <section className="py-24 px-6 border-t border-border/40 bg-card/30">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">Flexible Architecture</p>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Choose how you build</h2>
            <p className="text-muted-foreground max-w-2xl mx-auto">
              Coming from Laravel? Choose Single. MERN stack? Choose Double.
              Building a SaaS? Choose Triple. Grit adapts to your workflow.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4">
            {[
              { name: 'Single', icon: <Zap className="h-5 w-5" />, desc: 'Go + embedded SPA', flag: '--single', color: 'text-sky-400 bg-sky-400/10' },
              { name: 'Double', icon: <Layers className="h-5 w-5" />, desc: 'Web + API monorepo', flag: '--double', color: 'text-violet-400 bg-violet-400/10' },
              { name: 'Triple', icon: <Server className="h-5 w-5" />, desc: 'Web + Admin + API', flag: '--triple', color: 'text-emerald-400 bg-emerald-400/10' },
              { name: 'API Only', icon: <Database className="h-5 w-5" />, desc: 'Go backend only', flag: '--api', color: 'text-amber-400 bg-amber-400/10' },
              { name: 'Mobile', icon: <Smartphone className="h-5 w-5" />, desc: 'API + Expo', flag: '--mobile', color: 'text-rose-400 bg-rose-400/10' },
            ].map((arch) => (
              <div key={arch.name} className="rounded-xl border border-border/40 bg-card/50 p-5 text-center hover:border-primary/20 transition-colors">
                <div className={`h-12 w-12 rounded-xl ${arch.color} flex items-center justify-center mx-auto mb-3`}>
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

      {/* ═══════════════════════════════════════════════════════
          FEATURE DEEP-DIVE — Deploy
      ═══════════════════════════════════════════════════════ */}
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
                No CI/CD pipeline needed for getting started.
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

      {/* ═══════════════════════════════════════════════════════
          FEATURE DEEP-DIVE — Batteries
      ═══════════════════════════════════════════════════════ */}
      <section className="py-24 px-6 border-t border-border/40 bg-card/30">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">What Ships With Every Project</p>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Production-ready from day one</h2>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {[
              { title: 'JWT Authentication', desc: 'Register, login, refresh, OAuth2 social login (Google, GitHub). Role-based access control with ADMIN, EDITOR, USER.', icon: '🔐' },
              { title: 'Two-Factor Auth (TOTP)', desc: 'Authenticator app support, 10 backup codes (bcrypt-hashed), trusted devices with 30-day sliding cookie.', icon: '🛡️' },
              { title: 'File Storage (S3)', desc: 'Presigned URL uploads to AWS S3, Cloudflare R2, or MinIO. Image processing. Progress tracking.', icon: '📁' },
              { title: 'Email (Resend)', desc: 'Transactional emails with Go HTML templates. Welcome, password reset, verification. Dev uses Mailhog.', icon: '✉️' },
              { title: 'Background Jobs', desc: 'Redis-backed job queue via asynq. Image processing, email sending, cleanup workers.', icon: '⚡' },
              { title: 'Cron Scheduler', desc: 'Recurring tasks with cron expressions. Same worker pool as background jobs.', icon: '🕐' },
              { title: 'Redis Cache', desc: 'Cache middleware for any route. Set/Get/Delete for custom caching. Configurable TTL.', icon: '💾' },
              { title: 'AI Gateway', desc: 'Vercel AI Gateway — one key, hundreds of models. Streaming chat, completions. Zero markup on tokens.', icon: '🤖' },
              { title: 'Security (Sentinel)', desc: 'WAF, rate limiting, brute-force protection, anomaly detection. Real-time threat dashboard.', icon: '🔒' },
              { title: 'API Docs (gin-docs)', desc: 'Auto-generated OpenAPI spec from routes + GORM models. Interactive Scalar UI. No annotations.', icon: '📖' },
              { title: 'DB Browser (GORM Studio)', desc: 'Visual database browser. View tables, run queries, export data. Built into every project.', icon: '🗄️' },
              { title: 'Observability (Pulse)', desc: 'Request tracing, DB monitoring, runtime metrics, error tracking, Prometheus export.', icon: '📊' },
            ].map((feature) => (
              <div key={feature.title} className="rounded-xl border border-border/40 bg-background p-5 hover:border-primary/20 transition-colors group">
                <span className="text-2xl mb-3 block">{feature.icon}</span>
                <h3 className="font-semibold text-foreground mb-1.5 text-sm">{feature.title}</h3>
                <p className="text-xs text-muted-foreground leading-relaxed">{feature.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ═══════════════════════════════════════════════════════
          COMPARISON TABLE
      ═══════════════════════════════════════════════════════ */}
      <section className="py-24 px-6 border-t border-border/40">
        <div className="max-w-5xl mx-auto">
          <div className="text-center mb-16">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">Framework Comparison</p>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">How Grit compares</h2>
            <p className="text-muted-foreground max-w-2xl mx-auto">
              See what you get out of the box versus setting it up yourself with other frameworks.
            </p>
          </div>

          <div className="overflow-x-auto rounded-xl border border-border/40">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border/40 bg-card/80">
                  <th className="px-5 py-4 text-left font-medium text-muted-foreground w-[200px]">Feature</th>
                  <th className="px-5 py-4 text-center font-semibold text-primary">Grit</th>
                  <th className="px-5 py-4 text-center font-medium text-muted-foreground">Next.js</th>
                  <th className="px-5 py-4 text-center font-medium text-muted-foreground">Laravel</th>
                  <th className="px-5 py-4 text-center font-medium text-muted-foreground">T3 Stack</th>
                  <th className="px-5 py-4 text-center font-medium text-muted-foreground hidden lg:table-cell">Goravel</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/30">
                {[
                  { feature: 'Go Backend', grit: true, next: false, laravel: false, t3: false, goravel: true },
                  { feature: 'React Frontend', grit: true, next: true, laravel: false, t3: true, goravel: false },
                  { feature: 'Admin Panel', grit: true, next: false, laravel: 'partial', t3: false, goravel: false },
                  { feature: 'Code Generator (CLI)', grit: true, next: false, laravel: true, t3: false, goravel: true },
                  { feature: 'JWT + OAuth2 Auth', grit: true, next: false, laravel: true, t3: 'partial', goravel: true },
                  { feature: 'Two-Factor Auth (TOTP)', grit: true, next: false, laravel: false, t3: false, goravel: false },
                  { feature: 'File Storage (S3)', grit: true, next: false, laravel: true, t3: false, goravel: true },
                  { feature: 'Background Jobs', grit: true, next: false, laravel: true, t3: false, goravel: true },
                  { feature: 'Email System', grit: true, next: false, laravel: true, t3: false, goravel: true },
                  { feature: 'AI Integration', grit: true, next: false, laravel: false, t3: false, goravel: false },
                  { feature: 'API Documentation', grit: true, next: false, laravel: false, t3: false, goravel: false },
                  { feature: 'Database Browser', grit: true, next: false, laravel: false, t3: false, goravel: false },
                  { feature: 'Security (WAF)', grit: true, next: false, laravel: false, t3: false, goravel: false },
                  { feature: 'Observability', grit: true, next: false, laravel: false, t3: false, goravel: false },
                  { feature: 'One-Command Deploy', grit: true, next: false, laravel: false, t3: false, goravel: true },
                  { feature: 'Maintenance Mode', grit: true, next: false, laravel: true, t3: false, goravel: true },
                  { feature: 'Desktop App (Wails)', grit: true, next: false, laravel: false, t3: false, goravel: false },
                  { feature: 'Multiple Architectures', grit: true, next: false, laravel: false, t3: false, goravel: false },
                  { feature: 'TanStack Router Option', grit: true, next: false, laravel: false, t3: false, goravel: false },
                  { feature: 'Shared Types (Zod)', grit: true, next: false, laravel: false, t3: true, goravel: false },
                ].map((row) => (
                  <tr key={row.feature} className="hover:bg-accent/20 transition-colors">
                    <td className="px-5 py-3 font-medium text-foreground/90 text-[13px]">{row.feature}</td>
                    {[row.grit, row.next, row.laravel, row.t3, row.goravel].map((val, i) => (
                      <td key={i} className={`px-5 py-3 text-center ${i === 4 ? 'hidden lg:table-cell' : ''}`}>
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

          <p className="text-center text-xs text-muted-foreground/60 mt-4">
            Comparison based on out-of-the-box features without additional packages or configuration.
          </p>
        </div>
      </section>

      {/* ═══════════════════════════════════════════════════════
          TECH STACK
      ═══════════════════════════════════════════════════════ */}
      <section className="py-24 px-6 border-t border-border/40 bg-card/30">
        <div className="max-w-4xl mx-auto text-center">
          <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">Tech Stack</p>
          <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-12">Built on proven technologies</h2>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            {[
              { name: 'Go', sub: 'Gin + GORM', color: 'text-cyan-400' },
              { name: 'React', sub: 'Next.js / TanStack', color: 'text-blue-400' },
              { name: 'TypeScript', sub: 'Zod + TanStack Query', color: 'text-blue-300' },
              { name: 'Tailwind CSS', sub: 'shadcn/ui', color: 'text-sky-400' },
              { name: 'PostgreSQL', sub: 'GORM ORM', color: 'text-blue-400' },
              { name: 'Redis', sub: 'Cache + Jobs', color: 'text-red-400' },
              { name: 'Docker', sub: 'Compose', color: 'text-blue-300' },
              { name: 'Turborepo', sub: 'Monorepo', color: 'text-violet-400' },
            ].map((tech) => (
              <div key={tech.name} className="p-4 rounded-xl border border-border/40 bg-card/50 hover:border-primary/20 transition-colors">
                <div className={`font-semibold ${tech.color} mb-0.5`}>{tech.name}</div>
                <div className="text-xs text-muted-foreground">{tech.sub}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ═══════════════════════════════════════════════════════
          CTA
      ═══════════════════════════════════════════════════════ */}
      <section className="py-24 px-6 border-t border-border/40 bg-card/30">
        <div className="max-w-3xl mx-auto text-center">
          <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Start building in 30 seconds
          </h2>
          <p className="text-muted-foreground mb-8 text-lg">
            Install the CLI and scaffold your first project. No config, no boilerplate.
          </p>
          <CodeBlock language="bash" className="mb-8 text-left" code={`go install github.com/MUKE-coder/grit/v2/cmd/grit@latest
grit new my-app`} />
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Button size="lg" className="bg-primary hover:bg-primary/80 text-primary-foreground px-8 h-12 text-base rounded-full" asChild>
              <Link href="/docs/getting-started/quick-start">
                Read the docs <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
            <Button variant="outline" size="lg" className="border-border/60 text-foreground hover:bg-accent/30 px-8 h-12 text-base rounded-full" asChild>
              <Link href="https://github.com/MUKE-coder/grit" target="_blank">
                <Github className="mr-2 h-4 w-4" /> Star on GitHub
              </Link>
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
