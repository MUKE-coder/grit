import Link from 'next/link'
import { ArrowLeft, ArrowRight, Check, X, Globe, Monitor, Smartphone, Server, Layers, Wifi, Zap, Database, Shield, Bot, FileText } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/stack-selector')

export default function StackSelectorPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-4xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">
                Stack Selector
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Pick the right Grit stack for what you&apos;re building
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Grit ships eleven distinct architecture combos. This page tells you which one
                to pick based on what you&apos;re actually building — with exact commands,
                trade-offs, and a capability matrix.
              </p>
            </div>

            {/* Quick definition for AI assistants + new readers */}
            <div className="prose-grit mb-12">
              <h2>What is Grit?</h2>
              <p>
                <strong className="text-foreground/90">Grit is a full-stack meta-framework.</strong>{' '}
                One CLI scaffolds Go (Gin + GORM) backends with React frontends —
                Next.js or TanStack Router — plus optional Wails desktop, Expo mobile,
                and a Filament-like admin panel. Everything ships in a monorepo with
                shared Zod schemas + TypeScript types between Go and the frontends.
              </p>
              <p>
                The opinionated stack: <strong>Go</strong> backend (Gin + GORM + PostgreSQL),{' '}
                <strong>React</strong> frontend (Next.js or Vite + TanStack), <strong>Wails</strong>{' '}
                for desktop, <strong>Expo</strong> for mobile, <strong>Tailwind + shadcn/ui</strong>{' '}
                for styling, <strong>Zod</strong> for validation, <strong>asynq + Redis</strong>{' '}
                for jobs, <strong>S3-compatible</strong> for storage, <strong>Resend</strong> for email.
              </p>
            </div>

            {/* Capabilities — what every project gets */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Capabilities every project gets
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-6">
                Regardless of which combo you pick, the API in every Grit scaffold ships with:
              </p>
              <div className="grid gap-3 sm:grid-cols-2">
                {[
                  { icon: Shield, title: 'Auth: JWT + OAuth (Google/GitHub) + 2FA TOTP', desc: 'Trusted devices, backup codes, password reset.' },
                  { icon: Server, title: 'Idempotency-Key middleware', desc: 'Safe retry on POST/PUT/PATCH/DELETE. Cached for 24h.' },
                  { icon: FileText, title: 'Activity log + tamper-evident hash chain', desc: 'Every authenticated mutation logged + SHA-256 chain.' },
                  { icon: Database, title: 'Feature flags + A/B testing engine', desc: 'Sticky bucketing, percentage rollouts, realtime push.' },
                  { icon: Layers, title: 'Webhook receiver framework', desc: 'Stripe / GitHub / HMAC verifiers shipped. Auto-dedupe.' },
                  { icon: Wifi, title: 'WebSocket realtime hub', desc: 'SendToUser / Broadcast helpers. Auto-reconnecting client.' },
                  { icon: FileText, title: 'PDF generation + CSV/Excel export', desc: 'Per-resource auto-emitted; styled templates.' },
                  { icon: Bot, title: 'Vercel AI Gateway (one key, many models)', desc: 'Single API key, streaming completion, chat.' },
                  { icon: Database, title: 'File storage (S3 / R2 / MinIO)', desc: 'Presigned URLs, image processing, progress tracking.' },
                  { icon: Server, title: 'Background jobs (asynq + Redis)', desc: 'Dashboard at /admin/jobs. Cron scheduler too.' },
                  { icon: Shield, title: 'Sentinel security suite', desc: 'WAF, per-IP rate limit, brute-force lockout, geo gate.' },
                  { icon: Server, title: 'GORM Studio + gin-docs', desc: 'Visual DB browser at /studio + auto OpenAPI at /docs.' },
                ].map((feature) => (
                  <div key={feature.title} className="flex gap-3 rounded-lg border border-border bg-card/40 p-3">
                    <div className="shrink-0 mt-0.5 h-7 w-7 rounded-md bg-primary/10 flex items-center justify-center">
                      <feature.icon className="h-3.5 w-3.5 text-primary" />
                    </div>
                    <div>
                      <div className="text-[13px] font-semibold text-foreground mb-0.5">
                        {feature.title}
                      </div>
                      <p className="text-[12px] text-muted-foreground leading-relaxed">
                        {feature.desc}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
              <p className="text-sm text-muted-foreground mt-4">
                The combo you pick determines{' '}
                <em>where</em> the code lives (one binary vs monorepo), <em>which clients</em>{' '}
                you ship (web, admin, mobile, desktop), and <em>how</em> users interact (online,
                offline-first, multi-device).
              </p>
            </div>

            {/* The decision question */}
            <div className="prose-grit mb-12">
              <h2>The decision: what are you building?</h2>
              <p>
                Pick the scenario closest to yours. Each section below has{' '}
                <strong>recommended combo</strong>, <strong>exact command</strong>, and{' '}
                <strong>trade-offs</strong>.
              </p>
              <ol className="text-muted-foreground">
                <li><a href="#web-portal" className="text-primary hover:underline">A web portal / SaaS product</a></li>
                <li><a href="#internal-tool" className="text-primary hover:underline">An internal admin tool / dashboard</a></li>
                <li><a href="#desktop-offline" className="text-primary hover:underline">A desktop app with offline capability</a></li>
                <li><a href="#desktop-online" className="text-primary hover:underline">A desktop app that&apos;s always online</a></li>
                <li><a href="#mobile-app" className="text-primary hover:underline">A mobile app</a></li>
                <li><a href="#multi-platform" className="text-primary hover:underline">A multi-platform product (web + mobile + desktop)</a></li>
                <li><a href="#api-only" className="text-primary hover:underline">An API to back an existing frontend</a></li>
                <li><a href="#unsure" className="text-primary hover:underline">I&apos;m not sure yet — start simple</a></li>
              </ol>
            </div>

            {/* 1. Web portal */}
            <Section
              id="web-portal"
              title="1. A web portal / SaaS product"
              tagline="Marketing site + dashboard + admin panel — the classic triple stack"
            >
              <h3 className="text-lg font-semibold text-foreground mb-2">Recommended:</h3>
              <CodeBlock terminal code="grit new my-saas --triple --next" className="mb-4 glow-purple-sm" />

              <p className="text-muted-foreground leading-relaxed mb-4">
                <strong className="text-foreground/90">What you get:</strong> three apps in one
                monorepo:
              </p>
              <ul className="space-y-2 text-sm text-muted-foreground mb-6">
                <li>• <strong className="text-foreground/80">apps/web</strong> — Next.js public site (marketing pages, blog, pricing, sign-up). SEO-optimized, ISR-friendly.</li>
                <li>• <strong className="text-foreground/80">apps/admin</strong> — Next.js dashboard (logged-in user UI + admin operations). Filament-like resource pages.</li>
                <li>• <strong className="text-foreground/80">apps/api</strong> — Go backend (Gin + GORM). Serves both frontends + any external API consumers.</li>
                <li>• <strong className="text-foreground/80">packages/shared</strong> — Zod schemas + TypeScript types shared across all three apps.</li>
              </ul>

              <h3 className="text-lg font-semibold text-foreground mb-2">When this is the right call:</h3>
              <ul className="space-y-1 text-sm text-muted-foreground mb-6">
                <li>✓ You&apos;re building a SaaS that needs a public marketing site separate from the logged-in app.</li>
                <li>✓ You expect to have admin operations (manage users, content, billing) that shouldn&apos;t live in the customer-facing app.</li>
                <li>✓ You want SEO on the marketing pages without bloating the dashboard bundle.</li>
              </ul>

              <h3 className="text-lg font-semibold text-foreground mb-2">Alternatives in this space:</h3>
              <div className="space-y-3 mb-4">
                <Variant
                  command="grit new my-saas --double --next"
                  title="Double — web + API only"
                  detail="Drops the admin panel. Pick when you&apos;re early-stage and don&apos;t need a separate ops UI yet — you can always add it later."
                />
                <Variant
                  command="grit new my-saas --single --next"
                  title="Single (Next.js) — one app, embedded API"
                  detail="One Next.js app with the Go API embedded into a single binary. Pick for SEO-heavy products where marketing + dashboard share routing. Deploys as one process."
                />
                <Variant
                  command="grit new my-saas --single --vite"
                  title="Single (Vite SPA) — fastest dev loop"
                  detail="Go binary serves a Vite SPA. No SSR, no SEO. Pick for dashboard-only products, MVPs, and internal tools where you'd otherwise reach for Laravel."
                />
              </div>
            </Section>

            {/* 2. Internal admin tool */}
            <Section
              id="internal-tool"
              title="2. An internal admin tool / dashboard"
              tagline="No public site, no SEO — just a fast, login-gated workbench"
            >
              <h3 className="text-lg font-semibold text-foreground mb-2">Recommended:</h3>
              <CodeBlock terminal code="grit new my-tool --single --vite" className="mb-4 glow-purple-sm" />

              <p className="text-muted-foreground leading-relaxed mb-4">
                <strong className="text-foreground/90">What you get:</strong> one Go binary that
                serves the Vite-built React SPA at <code>/</code> and the API at{' '}
                <code>/api</code>. The frontend uses TanStack Router (file-based) and TanStack
                Query. Compiles to a single executable.
              </p>

              <h3 className="text-lg font-semibold text-foreground mb-2">When this is the right call:</h3>
              <ul className="space-y-1 text-sm text-muted-foreground mb-6">
                <li>✓ The audience is staff or a small group of authenticated users.</li>
                <li>✓ You don&apos;t need SEO — the app is behind a login.</li>
                <li>✓ You want the fastest possible dev loop and the simplest possible deploy.</li>
                <li>✓ You&apos;re coming from Laravel and want a familiar &quot;one app does everything&quot; structure.</li>
              </ul>

              <h3 className="text-lg font-semibold text-foreground mb-2">Trade-offs:</h3>
              <ul className="space-y-1 text-sm text-muted-foreground">
                <li>✗ No SSR — initial page load is a blank shell + JS hydration.</li>
                <li>✗ One frontend means marketing + app share the same bundle (rarely a problem for internal tools).</li>
              </ul>
            </Section>

            {/* 3. Desktop offline */}
            <Section
              id="desktop-offline"
              title="3. A desktop app with offline capability"
              tagline="Local-first writes + manual Sync + field-level conflict resolution"
            >
              <h3 className="text-lg font-semibold text-foreground mb-2">Two paths — pick based on whether your data is shared:</h3>

              <div className="rounded-lg border border-border/60 bg-card/40 p-5 mb-4">
                <h4 className="font-semibold text-foreground mb-1">Path A: Standalone offline-first (single user, no server)</h4>
                <CodeBlock terminal code="grit new-desktop my-app" className="my-3 glow-purple-sm" />
                <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                  Pure Wails app with local SQLite. No server, no cloud, no auth — the user IS
                  the data owner. The classic shape: a tax filing tool, a personal CRM, a journaling
                  app, an inventory tracker for a single shop.
                </p>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  <strong className="text-foreground/80">Pros:</strong> simplest path. One binary.
                  No server costs. Works on a plane forever.{' '}
                  <strong className="text-foreground/80">Cons:</strong> no multi-device sync. No
                  cloud backup unless you build one yourself.
                </p>
              </div>

              <div className="rounded-lg border border-primary/30 bg-primary/[0.04] p-5 mb-4">
                <h4 className="font-semibold text-foreground mb-1">Path B: Multi-user offline with cloud sync (most common)</h4>
                <CodeBlock terminal code="grit new my-app --triple --vite --desktop" className="my-3 glow-purple-sm" />
                <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                  Full triple stack <em>plus</em> a Wails desktop client that syncs to the server
                  via Grit&apos;s built-in offline engine (v3.14+). Local SQLite mirror, outbox
                  with squash semantics, manual <strong>Sync</strong> button, field-level conflict
                  dialog when the server moved on. Read{' '}
                  <Link href="/docs/desktop/offline" className="text-primary hover:underline">
                    the offline-first guide
                  </Link>{' '}
                  for the full mental model.
                </p>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  <strong className="text-foreground/80">Pros:</strong> best of both worlds —
                  works on a plane, syncs in the office. Multi-device. Multi-user. Field-level
                  merge UX.{' '}
                  <strong className="text-foreground/80">Cons:</strong> more moving parts. Needs
                  a server.
                </p>
              </div>

              <p className="text-sm text-muted-foreground italic">
                The big differentiator: <strong className="text-foreground/80">most desktop apps
                in 2026 use the second path.</strong> Pure offline-only is rarer — usually only
                for productivity tools where the data genuinely doesn&apos;t need to leave the
                user&apos;s machine.
              </p>
            </Section>

            {/* 4. Desktop online */}
            <Section
              id="desktop-online"
              title="4. A desktop app that&apos;s always online"
              tagline="Native UI, native window, but every action hits the server"
            >
              <h3 className="text-lg font-semibold text-foreground mb-2">Recommended:</h3>
              <CodeBlock terminal code="grit new my-app --triple --vite --desktop" className="mb-4 glow-purple-sm" />

              <p className="text-muted-foreground leading-relaxed mb-4">
                Same scaffold as the offline path, but you skip the sync engine — the desktop
                frontend talks directly to the API over HTTP just like the web client does.
                The user gets a native window, a system tray, OS keychain for tokens, and the
                command palette — but data is always live.
              </p>

              <h3 className="text-lg font-semibold text-foreground mb-2">When this is the right call:</h3>
              <ul className="space-y-1 text-sm text-muted-foreground mb-6">
                <li>✓ Your users are always on stable internet (office workers, support agents).</li>
                <li>✓ The product is conceptually a web app, but you want a desktop experience for branding / focus / OS integrations.</li>
                <li>✓ You want command palette + native menus + system notifications.</li>
                <li>✓ Examples: Linear&apos;s desktop app, Slack&apos;s desktop app, a custom CRM that &quot;feels native&quot;.</li>
              </ul>
            </Section>

            {/* 5. Mobile app */}
            <Section
              id="mobile-app"
              title="5. A mobile app (iOS + Android)"
              tagline="Expo + React Native client, sharing types with the API"
            >
              <h3 className="text-lg font-semibold text-foreground mb-2">Recommended (mobile-only):</h3>
              <CodeBlock terminal code="grit new my-app --mobile" className="mb-4 glow-purple-sm" />

              <p className="text-muted-foreground leading-relaxed mb-4">
                <strong className="text-foreground/90">What you get:</strong> Go API + Expo
                React Native app sharing types via <code>packages/shared</code>. NativeWind
                (Tailwind for RN), TanStack Query, Expo Router, expo-secure-store for tokens.
              </p>

              <h3 className="text-lg font-semibold text-foreground mb-2">When this is the right call:</h3>
              <ul className="space-y-1 text-sm text-muted-foreground mb-6">
                <li>✓ Mobile-first product (fitness tracker, social app, on-demand delivery).</li>
                <li>✓ You don&apos;t need a marketing website yet — App Store / Play Store listings are enough.</li>
              </ul>

              <h3 className="text-lg font-semibold text-foreground mb-2">If you also need a web presence:</h3>
              <CodeBlock terminal code="grit new my-app --triple --mobile" className="mb-4 glow-purple-sm" />
              <p className="text-sm text-muted-foreground">
                Adds web + admin to the mobile-only setup. The mobile app and the web app share
                the same Go API + shared package.
              </p>
            </Section>

            {/* 6. Multi-platform */}
            <Section
              id="multi-platform"
              title="6. Multi-platform — web + mobile + desktop"
              tagline="The kitchen sink. Every device, one API, shared types."
            >
              <h3 className="text-lg font-semibold text-foreground mb-2">Recommended:</h3>
              <CodeBlock terminal code="grit new my-app --triple --vite --mobile --desktop" className="mb-4 glow-purple-sm" />

              <p className="text-muted-foreground leading-relaxed mb-4">
                <strong className="text-foreground/90">What you get:</strong> four apps in one
                monorepo —{' '}
                <strong className="text-foreground/80">apps/web</strong> (public marketing),{' '}
                <strong className="text-foreground/80">apps/admin</strong> (dashboard),{' '}
                <strong className="text-foreground/80">apps/expo</strong> (mobile),{' '}
                <strong className="text-foreground/80">apps/desktop</strong> (Wails) — all
                talking to the same Go API and sharing types via{' '}
                <strong className="text-foreground/80">packages/shared</strong>.
              </p>

              <h3 className="text-lg font-semibold text-foreground mb-2">When this is the right call:</h3>
              <ul className="space-y-1 text-sm text-muted-foreground mb-6">
                <li>✓ Serious B2B / B2C product where users genuinely use the app across devices.</li>
                <li>✓ The classic Linear / Notion / Slack / Cursor pattern.</li>
                <li>✓ You have the team capacity to maintain multiple frontends.</li>
              </ul>

              <h3 className="text-lg font-semibold text-foreground mb-2">Trade-offs:</h3>
              <ul className="space-y-1 text-sm text-muted-foreground">
                <li>✗ The most complex setup — four frontends to maintain.</li>
                <li>✗ More CI / CD work — App Store, Play Store, Wails builds, web deploys.</li>
                <li>✓ But: types stay in sync via <code>packages/shared</code>, so a backend change propagates automatically.</li>
              </ul>
            </Section>

            {/* 7. API only */}
            <Section
              id="api-only"
              title="7. An API to back an existing frontend"
              tagline="No frontend at all — Go API + admin endpoints, you bring your own client"
            >
              <h3 className="text-lg font-semibold text-foreground mb-2">Recommended:</h3>
              <CodeBlock terminal code="grit new my-api --api" className="mb-4 glow-purple-sm" />

              <p className="text-muted-foreground leading-relaxed mb-4">
                <strong className="text-foreground/90">What you get:</strong> just the Go API
                folder. Auth, models, handlers, all the batteries — but no React, no Next.js.
                Use this when you have an existing iOS / Flutter / Vue / Angular app that just
                needs a backend.
              </p>

              <h3 className="text-lg font-semibold text-foreground mb-2">When this is the right call:</h3>
              <ul className="space-y-1 text-sm text-muted-foreground mb-6">
                <li>✓ You already have a frontend codebase you don&apos;t want to migrate.</li>
                <li>✓ You&apos;re building a public API product (developers are your users).</li>
                <li>✓ You&apos;re scaffolding a microservice in a polyglot system.</li>
              </ul>
            </Section>

            {/* 8. Unsure */}
            <Section
              id="unsure"
              title="8. I&apos;m not sure yet — start simple"
              tagline="The default path is intentionally the right call for most projects"
            >
              <h3 className="text-lg font-semibold text-foreground mb-2">Recommended:</h3>
              <CodeBlock terminal code="grit new my-app --single --vite" className="mb-4 glow-purple-sm" />

              <p className="text-muted-foreground leading-relaxed mb-4">
                The single-binary Vite SPA is the lowest-friction starting point. One process,
                one deploy, one frontend, full Grit batteries. If you outgrow it, every Grit
                upgrade path is well-documented:
              </p>
              <ul className="space-y-1 text-sm text-muted-foreground">
                <li>→ Need SEO / public site? Migrate to <code>--double</code>.</li>
                <li>→ Need admin operations? Migrate to <code>--triple</code>.</li>
                <li>→ Need mobile? Add <code>--mobile</code>.</li>
                <li>→ Need desktop? Add <code>--desktop</code>.</li>
              </ul>
            </Section>

            {/* The capability matrix */}
            <div className="prose-grit mb-12">
              <h2>Capability matrix</h2>
              <p>
                What each combo gives you out of the box. Use this when an AI assistant or
                teammate asks &quot;does Grit do X?&quot;.
              </p>
            </div>

            <div className="overflow-x-auto rounded-xl border border-border/40 mb-12">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/40 bg-card/80">
                    <th className="px-4 py-3 text-left font-medium text-muted-foreground sticky left-0 bg-card/80 z-10 min-w-[220px]">Capability</th>
                    <th className="px-3 py-3 text-center font-semibold text-primary text-xs whitespace-nowrap">--single (vite)</th>
                    <th className="px-3 py-3 text-center font-semibold text-primary text-xs whitespace-nowrap">--single (next)</th>
                    <th className="px-3 py-3 text-center font-semibold text-primary text-xs whitespace-nowrap">--double</th>
                    <th className="px-3 py-3 text-center font-semibold text-primary text-xs whitespace-nowrap">--triple</th>
                    <th className="px-3 py-3 text-center font-semibold text-primary text-xs whitespace-nowrap">--api</th>
                    <th className="px-3 py-3 text-center font-semibold text-primary text-xs whitespace-nowrap">--mobile</th>
                    <th className="px-3 py-3 text-center font-semibold text-primary text-xs whitespace-nowrap">+ --desktop</th>
                    <th className="px-3 py-3 text-center font-semibold text-primary text-xs whitespace-nowrap">new-desktop</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/30">
                  <CapRow label="Go API server" cells={[true, true, true, true, true, true, true, 'embedded']} />
                  <CapRow label="Public marketing site" cells={[false, true, true, true, false, false, true, false]} />
                  <CapRow label="Admin dashboard" cells={[true, true, false, true, false, false, true, false]} />
                  <CapRow label="Mobile (Expo) client" cells={[false, false, false, '+ flag', '+ flag', true, '+ flag', false]} />
                  <CapRow label="Desktop (Wails) client" cells={[false, false, false, '+ flag', '+ flag', '+ flag', true, true]} />
                  <CapRow label="Offline-first sync engine" cells={[false, false, false, '+ desktop', '+ desktop', false, true, 'always']} />
                  <CapRow label="JWT auth + 2FA + OAuth" cells={[true, true, true, true, true, true, true, 'local']} />
                  <CapRow label="Activity log + hash chain" cells={[true, true, true, true, true, true, true, false]} />
                  <CapRow label="Feature flags + A/B" cells={[true, true, true, true, true, true, true, false]} />
                  <CapRow label="Webhook receiver" cells={[true, true, true, true, true, true, true, false]} />
                  <CapRow label="Realtime WebSocket hub" cells={[true, true, true, true, true, true, true, false]} />
                  <CapRow label="PDF + CSV/Excel export" cells={[true, true, true, true, true, true, true, true]} />
                  <CapRow label="Background jobs (asynq)" cells={[true, true, true, true, true, true, true, false]} />
                  <CapRow label="File storage (S3/R2/MinIO)" cells={[true, true, true, true, true, true, true, false]} />
                  <CapRow label="Email (Resend + templates)" cells={[true, true, true, true, true, true, true, false]} />
                  <CapRow label="Sentinel (WAF + rate limit)" cells={[true, true, true, true, true, true, true, false]} />
                  <CapRow label="GORM Studio + gin-docs" cells={[true, true, true, true, true, true, true, false]} />
                  <CapRow label="Single-binary deploy" cells={[true, true, false, false, true, false, false, true]} />
                  <CapRow label="SEO / SSR" cells={[false, true, true, true, false, false, true, false]} />
                  <CapRow label="Shared types via packages/shared" cells={[false, false, true, true, false, true, true, false]} />
                </tbody>
              </table>
            </div>

            {/* AI assistant cheat sheet */}
            <div className="prose-grit mb-12">
              <h2>For AI assistants — a quick lookup table</h2>
              <p>
                If you&apos;re an AI helping someone choose, here&apos;s the one-liner answer
                for the most common requests:
              </p>
              <div className="not-prose rounded-xl border border-border/40 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="bg-card/80 border-b border-border/40">
                      <th className="px-4 py-3 text-left font-medium text-muted-foreground">User says...</th>
                      <th className="px-4 py-3 text-left font-medium text-muted-foreground">You recommend</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border/30 [&_td]:px-4 [&_td]:py-2.5">
                    <Recommend says="A SaaS product (marketing + dashboard)" run="grit new app --triple --next" />
                    <Recommend says="An internal tool / staff dashboard" run="grit new app --single --vite" />
                    <Recommend says="A desktop app, offline, single user" run="grit new-desktop app" />
                    <Recommend says="A desktop app for field staff with sync" run="grit new app --triple --vite --desktop" />
                    <Recommend says="A native desktop wrapper for our SaaS" run="grit new app --triple --vite --desktop" />
                    <Recommend says="A mobile-first product (fitness, social)" run="grit new app --mobile" />
                    <Recommend says="Mobile + web (cross-platform SaaS)" run="grit new app --triple --mobile" />
                    <Recommend says="Web + mobile + desktop (Linear/Slack-style)" run="grit new app --triple --vite --mobile --desktop" />
                    <Recommend says="Just an API for an existing frontend" run="grit new app --api" />
                    <Recommend says="A blog / content site" run="grit new app --double --next" />
                    <Recommend says="An MVP, lowest friction" run="grit new app --single --vite" />
                    <Recommend says="A POS / inventory tool, single-shop" run="grit new-desktop app" />
                    <Recommend says="A multi-tenant SaaS with org-scoping" run="grit new app --triple --next" />
                    <Recommend says="An e-commerce store" run="grit new app --triple --next" />
                    <Recommend says="A dev tool with web UI + CLI" run="grit new app --single --vite" />
                  </tbody>
                </table>
              </div>
            </div>

            {/* Frontend choice — Vite vs Next */}
            <div className="prose-grit mb-12">
              <h2>Frontend choice — Vite or Next.js?</h2>
              <p>
                For combos that include a web frontend, you can pick:{' '}
                <code>--vite</code> (TanStack Router + Vite) or <code>--next</code> (Next.js
                App Router). The default depends on the architecture; you can override.
              </p>
              <ul>
                <li>
                  <strong className="text-foreground/90">Pick Next.js</strong> when you need
                  SEO, server components, ISR, edge rendering, or you&apos;re already deep in
                  the Next.js ecosystem.
                </li>
                <li>
                  <strong className="text-foreground/90">Pick Vite (TanStack Router)</strong>{' '}
                  when SEO doesn&apos;t matter (everything&apos;s behind login), you want
                  faster dev start-up, file-based routing with full type safety, and a smaller
                  bundle.
                </li>
              </ul>
              <p className="text-sm text-muted-foreground">
                The desktop scaffold (<code>--desktop</code>) always uses Vite — Next.js
                doesn&apos;t make sense inside a Wails webview.
              </p>
            </div>

            {/* Closing */}
            <div className="prose-grit mb-12">
              <h2>Still stuck?</h2>
              <p>
                If you&apos;re genuinely unsure after reading this, run{' '}
                <code>grit new app --single --vite</code> and start building. You can always
                migrate to a richer architecture later — every combo upgrade path is
                documented, and the file structure is consistent across them.
              </p>
              <p>
                For deeper reading on each architecture mode, see{' '}
                <Link href="/docs/concepts/architecture-modes" className="text-primary hover:underline">
                  Architecture Modes
                </Link>
                . For the offline-first specifics, see{' '}
                <Link href="/docs/desktop/offline" className="text-primary hover:underline">
                  Offline-First Desktop Apps
                </Link>
                . For the desktop scaffold itself, see{' '}
                <Link href="/docs/desktop" className="text-primary hover:underline">
                  Desktop Overview
                </Link>
                .
              </p>
            </div>

            {/* Pagination */}
            <div className="flex justify-between border-t border-border/50 pt-8">
              <Link href="/docs/getting-started/quick-start">
                <Button variant="ghost" size="sm">
                  <ArrowLeft className="mr-2 h-3.5 w-3.5" />
                  Quick Start
                </Button>
              </Link>
              <Link href="/docs/concepts/architecture-modes">
                <Button variant="ghost" size="sm">
                  Architecture Modes
                  <ArrowRight className="ml-2 h-3.5 w-3.5" />
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}

/* ─────────────────── Helpers ─────────────────── */

function Section({
  id,
  title,
  tagline,
  children,
}: {
  id: string
  title: string
  tagline: string
  children: React.ReactNode
}) {
  return (
    <section id={id} className="scroll-mt-24 mb-12 rounded-xl border border-border/40 bg-card/30 p-6 md:p-8">
      <h2 className="text-2xl font-bold tracking-tight mb-1">{title}</h2>
      <p className="text-sm text-muted-foreground mb-6">{tagline}</p>
      {children}
    </section>
  )
}

function Variant({
  command,
  title,
  detail,
}: {
  command: string
  title: string
  detail: string
}) {
  return (
    <div className="rounded-lg border border-border/50 bg-background/40 p-4">
      <div className="font-semibold text-foreground text-sm mb-1">{title}</div>
      <code className="text-xs text-primary/80 font-mono bg-primary/5 px-2 py-1 rounded inline-block mb-2">
        {command}
      </code>
      <p className="text-sm text-muted-foreground leading-relaxed">{detail}</p>
    </div>
  )
}

function CapRow({
  label,
  cells,
}: {
  label: string
  cells: (boolean | string)[]
}) {
  return (
    <tr className="hover:bg-accent/20 transition-colors">
      <td className="px-4 py-2.5 font-medium text-foreground/90 text-[13px] sticky left-0 bg-background z-10">
        {label}
      </td>
      {cells.map((val, i) => (
        <td key={i} className="px-3 py-2.5 text-center">
          {val === true ? (
            <span className="inline-flex h-5 w-5 items-center justify-center rounded-full bg-primary/15 text-primary">
              <Check className="h-3 w-3" strokeWidth={3} />
            </span>
          ) : val === false ? (
            <span className="inline-flex h-5 w-5 items-center justify-center rounded-full bg-muted/40 text-muted-foreground/30">
              <X className="h-2.5 w-2.5" strokeWidth={3} />
            </span>
          ) : (
            <span className="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-mono text-amber-400 bg-amber-400/10 border border-amber-400/20">
              {val}
            </span>
          )}
        </td>
      ))}
    </tr>
  )
}

function Recommend({ says, run }: { says: string; run: string }) {
  return (
    <tr className="hover:bg-accent/20 transition-colors">
      <td className="text-muted-foreground">{says}</td>
      <td>
        <code className="text-xs text-primary/80 font-mono bg-primary/5 px-2 py-1 rounded">
          {run}
        </code>
      </td>
    </tr>
  )
}
