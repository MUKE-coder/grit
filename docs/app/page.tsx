import Link from 'next/link'
import type { Metadata } from 'next'
import { Activity, ArrowRight, Bot, Building2, Check, ChevronDown, Database, FileCheck, Flag, Github, HardDrive, Heart, Layers, LayoutDashboard, Lock, Mail, Monitor, Radio, Rocket, Server, Shield, Smartphone, Terminal, TestTube2, TrendingUp, UploadCloud, UserCheck, Webhook, Zap } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock } from '@/components/code-block'
import { HeroCodeTabs, InstallTabs } from '@/components/hero-code-tabs'
import { HomepageBenchmarks } from '@/components/homepage-benchmarks'
import { OneApiClients } from '@/components/one-api-clients'
import { ThemeShowcase } from '@/components/theme-showcase'
import { PlatformShowcase } from '@/components/platform-showcase'
import { SystemShowcase } from '@/components/system-showcase'
import { FormShowcase } from '@/components/form-showcase'
import { AuthShowcase } from '@/components/auth-showcase'
import { ResourceDefinitionShowcase } from '@/components/resource-definition-showcase'
import { InfraShowcase } from '@/components/infra-showcase'
import { ComplianceShowcase } from '@/components/compliance-showcase'
import { DeployShowcase } from '@/components/deploy-showcase'
import { SoftwareApplicationSchema, FAQPageSchema } from '@/components/structured-data'
import { FeatureTabs } from '@/components/feature-tabs'
import { CpuArchitecture } from '@/components/ui/cpu-architecture'
import { MagneticButton, GSAPSection, FadeIn, GlowOrb } from '@/components/motion-primitives'
import { ReactLogo, TanStackLogo } from '@/components/framework-logos'
import { CommunityCTA, WhatsAppIcon, WHATSAPP_COMMUNITY_URL } from '@/components/community-cta'
import { HubAndSpoke } from '@/components/hub-and-spoke'
import { WhatIsGrit } from '@/components/what-is-grit'
import { GRIT_VERSION } from '@/config/site'
import { GridFrame } from '@/components/grid-frame'
import { HomeSponsors } from '@/components/sponsors'

export const metadata: Metadata = {
  title: 'Grit — Go + React Full-Stack Framework',
  description: 'Build production-ready full-stack applications with Go and React. One CLI, 5 architectures, batteries included.',
  alternates: { canonical: 'https://gritframework.dev' },
}

// AsciiCube — Inertia-style ASCII/dotted 3D cube wireframe for hero corners.
// Pure SVG so it scales infinitely + costs nothing to render. The `flip`
// prop mirrors it for the opposing corner so the two read as a matched pair.
function AsciiCube({ className, flip = false }: { className?: string; flip?: boolean }) {
  return (
    <svg
      viewBox="0 0 200 200"
      className={className}
      style={{ transform: flip ? 'scaleX(-1)' : undefined }}
      aria-hidden
      fill="none"
      stroke="currentColor"
      strokeWidth="0.6"
    >
      {/* Bottom face (a parallelogram) */}
      <g strokeDasharray="2 3">
        <path d="M 30 130 L 100 170 L 170 130 L 100 90 Z" />
        {/* Front face */}
        <path d="M 30 130 L 30 60 L 100 20 L 100 90 Z" />
        {/* Right face */}
        <path d="M 100 90 L 100 20 L 170 60 L 170 130 Z" />
        {/* Inner grid — dotted lines suggesting cube subdivisions */}
        <line x1="30" y1="90" x2="100" y2="50" />
        <line x1="60" y1="115" x2="60" y2="40" />
        <line x1="80" y1="130" x2="80" y2="55" />
        <line x1="170" y1="90" x2="100" y2="50" />
        <line x1="130" y1="115" x2="130" y2="40" />
        <line x1="150" y1="105" x2="150" y2="45" />
        <line x1="65" y1="150" x2="135" y2="150" />
        <line x1="50" y1="140" x2="150" y2="140" />
      </g>
      {/* Dots at vertices */}
      {[[30, 130], [30, 60], [100, 20], [100, 90], [170, 130], [170, 60], [100, 170]].map(([cx, cy], i) => (
        <circle key={i} cx={cx} cy={cy} r="1.4" fill="currentColor" stroke="none" />
      ))}
    </svg>
  )
}

export default function HomePage() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SoftwareApplicationSchema />
      <FAQPageSchema />
      <SiteHeader />
      <GridFrame />

      {/* ═══ HERO — consistent dark base, blueprint grid (via GridFrame), glass pills, GitHub editor ═══ */}
      <section className="relative overflow-hidden">
        {/* Soft primary glow — the hero's only colour moment; the dark/light base
            and the grid + rails come from the shared GridFrame so the hero reads
            consistently with the rest of the page in both themes. */}
        <div
          className="absolute inset-0 -z-10"
          style={{
            background:
              'radial-gradient(ellipse 60% 55% at 50% -8%, hsl(var(--primary) / 0.18), transparent 60%)',
          }}
        />

        {/* ASCII-cube corner decorations */}
        <AsciiCube className="absolute -top-8 -left-8 w-[280px] h-[280px] hidden md:block opacity-[0.18] text-foreground" />
        <AsciiCube className="absolute -top-12 -right-8 w-[320px] h-[320px] hidden md:block opacity-[0.18] text-foreground" flip />
        <AsciiCube className="absolute -bottom-16 -left-16 w-[260px] h-[260px] hidden lg:block opacity-[0.15] text-foreground" />
        <AsciiCube className="absolute -bottom-8 -right-12 w-[280px] h-[280px] hidden lg:block opacity-[0.15] text-foreground" flip />

        {/* Firecrawl-style crosshair accents pinned to the grid */}
        <span className="crosshair absolute top-28 left-[11%] text-foreground/30 hidden md:block" style={{ width: 16, height: 16 }} />
        <span className="crosshair absolute top-44 right-[13%] text-primary/40 hidden md:block" style={{ width: 16, height: 16 }} />
        <span className="crosshair absolute bottom-36 left-[18%] text-foreground/25 hidden lg:block" style={{ width: 14, height: 14 }} />
        <span className="crosshair absolute bottom-24 right-[22%] text-foreground/30 hidden lg:block" style={{ width: 14, height: 14 }} />

        {/* Floating glow orbs */}
        <GlowOrb className="-top-32 left-1/4 h-[400px] w-[400px] bg-sky-400/30" duration={18} />
        <GlowOrb className="top-40 right-1/4 h-[300px] w-[300px] bg-cyan-300/20" delay={2} duration={16} />

        <div className="relative max-w-7xl mx-auto pt-14 pb-16 md:pt-20 md:pb-24 px-6">
          {/* Two-column, left-aligned — the shape Nuxt, Wasp and Encore all use.
              A centred hero sends the eye back to the middle for every line; an
              asymmetric split lets the copy read as a paragraph and gives the
              code somewhere permanent to sit. */}
          {/* minmax(0, …) on every track, at every breakpoint. A bare `grid`
              creates one auto-sized column that grows to its max-content width,
              and a grid item defaults to min-width:auto so it will not shrink
              below that. Together those made this hero 542px wide inside a
              370px phone and pushed the whole page into horizontal scroll —
              the headline read "The Batteries-Include" with the rest off-screen. */}
          <div className="grid grid-cols-[minmax(0,1fr)] lg:grid-cols-[minmax(0,1.15fr)_minmax(0,1fr)] gap-12 lg:gap-14 items-center">

            {/* ── Left: the pitch ─────────────────────────────────────── */}
            <div className="text-left">
              <FadeIn>
                <Link
                  href="/docs/changelog"
                  className="inline-flex items-center gap-2 rounded-full border border-primary/30 bg-primary/10 px-3 py-1 text-xs font-mono font-medium text-primary hover:bg-primary/15 transition-colors mb-7"
                >
                  <span className="h-1.5 w-1.5 rounded-full bg-primary" />
                  Grit v{GRIT_VERSION} is out
                  <ArrowRight className="h-3 w-3" />
                </Link>
              </FadeIn>

              {/* Says what it IS, in the words the audience already uses.
                  "Opinionated" and "batteries-included" are what people search
                  for once they are tired of assembling a stack themselves. */}
              <FadeIn delay={0.08}>
                <h1 className="font-display text-4xl md:text-5xl lg:text-[3.4rem] xl:text-[3.75rem] font-bold tracking-tight text-foreground mb-6 leading-[1.08]">
                  <span className="lg:whitespace-nowrap">The Batteries-Included</span>
                  <br />
                  <span className="lg:whitespace-nowrap">Full-Stack Framework</span>
                  <br />
                  <span className="bg-gradient-to-r from-primary via-sky-500 to-primary bg-clip-text text-transparent">
                    for Go &amp; React
                  </span>
                </h1>
              </FadeIn>

              <FadeIn delay={0.14}>
                <p className="text-base md:text-lg text-muted-foreground mb-8 leading-relaxed max-w-xl">
                  Opinionated by design. Describe a resource and Grit writes the Go model,
                  API, migrations, TypeScript types, React hooks and admin screen. Auth,
                  RBAC, jobs, storage, realtime, observability and deploy &mdash; all
                  configured out of the box.
                </p>
              </FadeIn>

              <FadeIn delay={0.2}>
                <div className="flex flex-wrap items-center gap-3 mb-8">
                  <MagneticButton>
                    <Link
                      href="/docs/start"
                      className="group relative inline-flex items-center justify-center gap-2 h-12 px-7 rounded-full bg-primary text-primary-foreground font-semibold text-sm glow-primary-sm hover:bg-primary/90 transition-all"
                    >
                      Get started
                      <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
                    </Link>
                  </MagneticButton>
                  <MagneticButton>
                    <Link
                      href="/docs/getting-started/philosophy"
                      className="inline-flex items-center justify-center h-12 px-7 rounded-full border border-border bg-card/60 backdrop-blur-xl text-foreground font-medium text-sm hover:bg-accent/40 transition-all"
                    >
                      Our philosophy
                    </Link>
                  </MagneticButton>
                </div>
              </FadeIn>

              {/* Platform-tabbed, Windows first — see InstallTabs. */}
              <FadeIn delay={0.26}>
                <InstallTabs />
                <p className="text-xs text-muted-foreground mt-2.5">
                  Detects an existing install and runs{' '}
                  <code className="text-foreground/80">grit update</code>, otherwise pulls the
                  right binary for your OS.{' '}
                  <Link
                    href="/docs/getting-started/installation"
                    className="text-foreground/80 hover:text-primary underline underline-offset-2"
                  >
                    All install options
                  </Link>
                </p>
              </FadeIn>

              {/* The stack, stated rather than implied. Small on purpose — it is
                  reassurance, not the headline. */}
              <FadeIn delay={0.3}>
                <div className="flex items-center gap-4 mt-9 pt-7 border-t border-border/40">
                  <span className="text-[11px] font-mono uppercase tracking-wider text-muted-foreground/70 shrink-0">
                    Built on
                  </span>
                  <div className="flex items-center gap-2.5">
                    {[
                      { src: '/images/icons/go.svg', alt: 'Go' },
                      { src: '/images/icons/postgressql.png', alt: 'Postgres' },
                      { src: '/images/icons/redis-logo-svgrepo-com.svg', alt: 'Redis' },
                      { src: '/images/icons/docker-svgrepo-com.svg', alt: 'Docker' },
                    ].map((logo) => (
                      <div
                        key={logo.alt}
                        title={logo.alt}
                        className="h-7 w-7 rounded-full bg-white shadow-[0_2px_8px_rgba(0,0,0,0.12)] flex items-center justify-center"
                      >
                        {/* eslint-disable-next-line @next/next/no-img-element */}
                        <img src={logo.src} alt={logo.alt} className="h-4 w-4 object-contain" />
                      </div>
                    ))}
                    <span aria-hidden className="mx-0.5 text-muted-foreground/40">
                      +
                    </span>
                    <div className="h-7 w-7 rounded-full bg-white shadow-[0_2px_8px_rgba(0,0,0,0.12)] flex items-center justify-center" title="React">
                      <ReactLogo className="h-4 w-4" />
                    </div>
                    <div className="h-7 w-7 rounded-full bg-white shadow-[0_2px_8px_rgba(0,0,0,0.12)] flex items-center justify-center" title="Next.js">
                      {/* eslint-disable-next-line @next/next/no-img-element */}
                      <img src="/images/icons/Next.js.svg" alt="Next.js" className="h-4 w-4 object-contain" />
                    </div>
                    <div className="h-7 w-7 rounded-full bg-white shadow-[0_2px_8px_rgba(0,0,0,0.12)] flex items-center justify-center" title="TanStack Router">
                      <TanStackLogo className="h-4 w-4" />
                    </div>
                    <div className="h-7 w-7 rounded-full bg-white shadow-[0_2px_8px_rgba(0,0,0,0.12)] flex items-center justify-center" title="Expo (React Native)">
                      <Smartphone className="h-4 w-4 text-slate-800" strokeWidth={2} />
                    </div>
                  </div>
                </div>
              </FadeIn>
            </div>

            {/* ── Right: four tabs, because "batteries included" is a claim
                   until someone sees the batteries ──────────────────── */}
            <FadeIn delay={0.34}>
              <HeroCodeTabs />
            </FadeIn>

          </div>
        </div>
      </section>

      {/* ═══ BENCHMARKS ═══

          Directly under the hero because it is the first question anyone asks
          about a Go framework, and because the answer is the strongest single
          argument this project has. Numbers live in config/benchmarks.ts so the
          chart cannot drift from the methodology pages behind it. */}
      <HomepageBenchmarks />

      {/* ═══ WHAT IS ALREADY IN THE BOX ═══

          Added after four independent reviews of this site each recommended
          building things that shipped months ago — an MCP server, generated
          tests, an audit log, webhooks, multi-tenancy, backups. All documented,
          all a thousand lines further down this page, all invisible to someone
          deciding in ninety seconds whether to keep reading.

          This is not new content. It is the same capabilities, scannable, near
          the top, each one linking to where it is explained properly. */}
      <section className="py-16 px-6 border-b border-border/40">
        <div className="max-w-5xl mx-auto">
          <div className="text-center mb-10">
            <h2 className="text-2xl md:text-3xl font-bold tracking-tight text-foreground mb-3">
              You are not going to build these again
            </h2>
            <p className="text-sm md:text-base text-muted-foreground max-w-2xl mx-auto leading-relaxed">
              Every one of these ships in a generated project. Switch on what you need,
              leave the rest off.
            </p>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-2.5">
            {[
              { icon: Lock, label: 'Auth, OAuth & 2FA', href: '/docs/backend/authentication' },
              { icon: UserCheck, label: 'Roles & permissions', href: '/docs/backend/rbac' },
              { icon: LayoutDashboard, label: 'Admin panel', href: '/docs/admin/overview' },
              { icon: Zap, label: 'Background jobs', href: '/docs/batteries' },
              { icon: UploadCloud, label: 'File storage & uploads', href: '/docs/batteries' },
              { icon: Radio, label: 'Realtime WebSockets', href: '/docs/backend/realtime' },
              { icon: Mail, label: 'Email & notifications', href: '/docs/batteries' },
              { icon: Bot, label: 'AI gateway', href: '/docs/ai-integration' },
              { icon: Activity, label: 'Observability (Pulse)', href: '/docs/backend/pulse' },
              { icon: Shield, label: 'WAF & rate limiting', href: '/docs/security' },
              { icon: FileCheck, label: 'Tamper-evident audit log', href: '/docs/security' },
              { icon: Building2, label: 'Multi-tenancy', href: '/docs/plugins/multitenant' },
              { icon: Webhook, label: 'Webhooks', href: '/docs/backend/webhooks' },
              { icon: Flag, label: 'Feature flags', href: '/docs/backend/feature-flags' },
              { icon: TestTube2, label: 'Generated tests', href: '/docs/testing' },
              { icon: HardDrive, label: 'Backups & restore', href: '/docs/deployment/checklist' },
              { icon: Monitor, label: 'Offline-first desktop', href: '/docs/desktop' },
              { icon: Smartphone, label: 'Mobile (Expo)', href: '/docs/mobile/getting-started' },
              { icon: Bot, label: 'MCP server for AI agents', href: '/docs/ai-integration' },
              { icon: Rocket, label: 'One-command deploy', href: '/docs/deployment' },
            ].map((cap) => (
              <Link
                key={cap.label}
                href={cap.href}
                className="group flex items-center gap-2.5 rounded-lg border border-border/50 bg-card/40 px-3 py-2.5 text-[13px] transition-colors hover:border-border hover:bg-card/70"
              >
                <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary transition-colors group-hover:bg-primary/15">
                  <cap.icon className="h-3.5 w-3.5" />
                </span>
                <span className="text-muted-foreground group-hover:text-foreground transition-colors leading-snug">
                  {cap.label}
                </span>
              </Link>
            ))}
          </div>

          <div className="mt-8 text-center">
            <Link
              href="/docs/start"
              className="inline-flex items-center gap-1.5 text-sm font-semibold text-primary hover:underline"
            >
              Start here &mdash; nothing to deployed in seven steps
              <ArrowRight className="h-3.5 w-3.5" />
            </Link>
          </div>
        </div>
      </section>

      {/* ═══ NEW: WHAT IS GRIT (60-second tour) ═══ */}
      <WhatIsGrit />

      {/* ═══ ONE FRAMEWORK, EVERY PLATFORM ═══

          Screenshots of real generated projects — the admin, the Wails window,
          the Expo app on an emulator, and the Scalar reference the API serves.
          Each tab pairs the shot with the exact commands that produce it, so
          the claim and the proof sit in the same frame. */}
      <section className="py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-6xl mx-auto">
          <div className="max-w-2xl mb-10">
            <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
              <span className="h-1.5 w-1.5 rounded-full bg-primary" /> Every platform
            </span>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-4 leading-tight">
              One framework.<br />Every platform.
            </h2>
            <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
              Web, desktop, mobile and a documented API &mdash; from one Go backend and one
              set of generated types. Every screenshot below is a real generated project,
              next to the commands that produce it.
            </p>
          </div>

          <PlatformShowcase />
        </div>
      </section>

      {/* ═══ WHAT ONE COMMAND GENERATES ═══

          Screenshots of the forms Grit produced from the field definitions
          shown beside them. "It generates a form for you" is a cheap claim;
          the uploads, the searchable FK with inline-create, the line-items
          table with live totals and the wizard are not. */}
      <section className="py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-6xl mx-auto">
          <div className="max-w-2xl mb-10">
            <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
              <span className="h-1.5 w-1.5 rounded-full bg-primary" /> Code generation
            </span>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-4 leading-tight">
              One command.<br />The whole form.
            </h2>
            <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
              Not a scaffold you finish by hand. Uploads that know what they accept, foreign
              keys you can fill without leaving the page, line-item tables that total
              themselves, wizards that save a step at a time &mdash; and a desktop app that
              keeps taking input with the network off.
            </p>
          </div>

          <FormShowcase />
        </div>
      </section>

      {/* ═══ THEMES ═══

          Real screenshots of a generated admin, not mockups. A mockup proves a
          designer can draw; a screenshot proves the framework produces it. This
          is also the one thing on the page no competitor can copy without
          building the admin first. */}
      <section className="py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-6xl mx-auto">
          <div className="max-w-2xl mb-10">
            <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
              <span className="h-1.5 w-1.5 rounded-full bg-primary" /> Themes
            </span>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-4 leading-tight">
              Four themes.<br />None of them look generated.
            </h2>
            <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
              Every scaffolded admin ships with four complete themes &mdash; not palettes, but
              coordinated typography, auth pages, dashboard and brand colour. Pick one at{' '}
              <code className="text-foreground/80 text-sm">grit new</code>, or change your mind
              later with one line in <code className="text-foreground/80 text-sm">.env</code>.
            </p>
          </div>

          <ThemeShowcase />
        </div>
      </section>

      {/* ═══ THE SYSTEM HUB ═══

          The batteries, one screenshot at a time. Several of these needed real
          infrastructure to capture — the jobs page shows work that five actual
          uploads enqueued through Redis. */}
      <section className="py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-6xl mx-auto">
          <div className="max-w-2xl mb-10">
            <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
              <span className="h-1.5 w-1.5 rounded-full bg-primary" /> Batteries included
            </span>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-4 leading-tight">
              The parts you always end up<br />building anyway.
            </h2>
            <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
              Permissions, brute-force protection, golden-signal metrics, an audit timeline,
              scheduled backups, S3 uploads, a job queue and a cron scheduler &mdash; each with
              a real screen, not a config file and a README. This is what{' '}
              <code className="text-foreground/80 text-sm">grit new</code> gives you on day one.
            </p>
          </div>

          <SystemShowcase />
        </div>
      </section>

      {/* ═══ AUTHENTICATION ═══

          Screenshots plus the env/route detail, because "auth included" means
          nothing until you can see which parts. The closing paragraph names
          what is absent on purpose. */}
      <section className="py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-6xl mx-auto">
          <div className="max-w-2xl mb-10">
            <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
              <span className="h-1.5 w-1.5 rounded-full bg-primary" /> Authentication
            </span>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-4 leading-tight">
              Auth you would otherwise<br />spend a month on.
            </h2>
            <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
              Sign-in pages, a JWT pair, revocable server-side sessions, Google and GitHub,
              TOTP with backup codes, database-backed roles, and enterprise SSO over OIDC or
              SAML. Working on the first run, not a tutorial to follow.
            </p>
          </div>

          <AuthShowcase />
        </div>
      </section>

      {/* ═══ COMPLIANCE & ENTERPRISE ═══

          The section a buyer's security reviewer reads. Every claim is one the
          code backs — "tamper-evident" is a real hash chain with a verify
          endpoint, not a turn of phrase. */}
      <section className="py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-6xl mx-auto">
          <div className="max-w-2xl mb-10">
            <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
              <span className="h-1.5 w-1.5 rounded-full bg-primary" /> Compliance
            </span>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-4 leading-tight">
              The questions enterprise<br />buyers ask first.
            </h2>
            <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
              GDPR export and erasure with a tamper-evident journal, SSO per customer over OIDC
              or SAML, access-review campaigns for SOC 2, and a hash-chained audit trail that
              exports to your SIEM. Built in, not bought later.
            </p>
          </div>

          <ComplianceShowcase />
        </div>
      </section>

      {/* ═══ RESOURCE DEFINITION ═══ */}
      <section className="py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-6xl mx-auto">
          <div className="max-w-2xl mb-10">
            <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
              <span className="h-1.5 w-1.5 rounded-full bg-primary" /> Customisation
            </span>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-4 leading-tight">
              One file describes<br />the whole screen.
            </h2>
            <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
              The generator writes a resource definition; after that it is ordinary TypeScript
              you own. Columns, filters, row and bulk actions, twenty field types, wizards,
              dropzone variants &mdash; changed by editing a file, not by fighting a generator
              that wants to overwrite it.
            </p>
          </div>

          <ResourceDefinitionShowcase />
        </div>
      </section>

      {/* ═══ INFRASTRUCTURE ═══ */}
      <section className="py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-6xl mx-auto">
          <div className="max-w-2xl mb-10">
            <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
              <span className="h-1.5 w-1.5 rounded-full bg-primary" /> Infrastructure
            </span>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-4 leading-tight">
              Your database.<br />Your bucket. Your Redis.
            </h2>
            <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
              Grit is opinionated about structure, not about who you rent from. Each of these
              is one environment variable, with a local default that works before you have any
              cloud account at all.
            </p>
          </div>

          <InfraShowcase />
        </div>
      </section>

      {/* ═══ DEPLOYMENT ═══

          Reads DEPLOYMENT_PROVIDERS, the same config the docs pages use, so
          the two cannot drift. */}
      <section className="py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-6xl mx-auto">
          <div className="max-w-2xl mb-10">
            <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
              <span className="h-1.5 w-1.5 rounded-full bg-primary" /> Deployment
            </span>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-4 leading-tight">
              It is a Go binary<br />and some containers.
            </h2>
            <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
              Which means it runs anywhere &mdash; a $5 VPS, a managed platform, or your own
              Docker host. Pick a target for the actual steps, what it costs, and the thing that
              catches people out.
            </p>
          </div>

          <DeployShowcase />
        </div>
      </section>

      {/* ═══ ONE API, EVERY CLIENT ═══

          The argument nothing else in this space can make, and it is invisible
          unless the four are side by side and clickable. A framework that
          scaffolds a backend is common; one where the same generated types and
          hooks drive Next.js, TanStack, Expo and an offline desktop binary is
          not. The left pane deliberately never changes. */}
      <section className="py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-6xl mx-auto">
          <div className="max-w-2xl mb-12">
            <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
              <span className="h-1.5 w-1.5 rounded-full bg-primary" /> One backend
            </span>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-4 leading-tight">
              Write the API once.<br />Ship it to every client.
            </h2>
            <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
              The generator emits the Go handler <em>and</em> the typed client — schemas,
              TypeScript types and React Query hooks. The same hook then works in a Next.js
              app, a TanStack SPA, an Expo phone app and an offline-capable desktop binary.
              Rename a field in the Go struct and every one of them stops compiling until
              you fix it.
            </p>
          </div>

          <OneApiClients />
        </div>
      </section>

      {/* ═══ CPU ARCHITECTURE — Grit as the central chip ═══ */}
      <section className="relative py-24 px-6 overflow-hidden border-t border-border/40">
        {/* Layered backdrop — subtle radial + dotted grid */}
        <div className="absolute inset-0 -z-30 bg-gradient-to-b from-background via-card/20 to-background" />
        <div
          className="absolute inset-0 -z-20 opacity-[0.4]"
          style={{
            backgroundImage:
              'radial-gradient(circle at 1px 1px, hsl(var(--foreground) / 0.08) 1px, transparent 0)',
            backgroundSize: '24px 24px',
          }}
        />
        <GlowOrb className="-top-40 left-1/2 -translate-x-1/2 h-[700px] w-[700px] bg-primary/[0.06]" duration={20} />

        <div className="max-w-6xl mx-auto">

          <GSAPSection>
            <div className="text-center mb-14" data-gsap-reveal>
              <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-5">
                <span className="h-1.5 w-1.5 rounded-full bg-primary" /> The Grit Core
              </span>
              <h2 className="text-3xl md:text-5xl font-bold tracking-tight text-foreground mb-5 leading-[1.1]">
                One CLI &mdash; eight production<br />primitives wired together
              </h2>
              <p className="text-base md:text-lg text-muted-foreground max-w-2xl mx-auto leading-relaxed">
                Grit is the chip on the board. Auth, jobs, storage, AI, observability,
                webhooks, realtime, and cache all light up the moment you scaffold,
                so you spend your time on product not plumbing.
              </p>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-12 gap-10 items-center">

              {/* CPU canvas — center 6 cols on desktop, full width below */}
              <div className="lg:col-span-7 order-2 lg:order-1" data-gsap-reveal>
                <div className="relative rounded-2xl border-2 border-border/60 bg-card/60 backdrop-blur p-8 shadow-[0_24px_64px_-16px_rgba(0,0,0,0.4),inset_0_1px_0_rgba(255,255,255,0.04),0_0_0_4px_hsl(var(--primary)/0.04)]">
                  {/* Corner bracket decorations */}
                  <div className="absolute top-3 left-3 h-4 w-4 border-t-2 border-l-2 border-primary/40" />
                  <div className="absolute top-3 right-3 h-4 w-4 border-t-2 border-r-2 border-primary/40" />
                  <div className="absolute bottom-3 left-3 h-4 w-4 border-b-2 border-l-2 border-primary/40" />
                  <div className="absolute bottom-3 right-3 h-4 w-4 border-b-2 border-r-2 border-primary/40" />
                  <CpuArchitecture
                    className="text-foreground/30 dark:text-foreground/25"
                    text="GRIT"
                  />
                  {/* Stamp footer */}
                  <div className="mt-4 flex items-center justify-between px-2">
                    <div className="flex items-center gap-2">
                      <div className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse" />
                      <span className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground">{`v${GRIT_VERSION} · production-ready`}</span>
                    </div>
                    <span className="text-[10px] font-mono text-muted-foreground/50">GRIT-FW-A1</span>
                  </div>
                </div>
              </div>

              {/* Feature spokes — 2-col on desktop, full grid on mobile */}
              <div className="lg:col-span-5 order-1 lg:order-2 grid grid-cols-2 gap-3" data-gsap-reveal>
                {[
                  { label: 'Auth',          sub: 'JWT · OAuth · 2FA',     dot: 'bg-sky-400' },
                  { label: 'AI Gateway',    sub: '100+ models · stream',  dot: 'bg-violet-400' },
                  { label: 'File Storage',  sub: 'S3 · R2 · MinIO',       dot: 'bg-amber-400' },
                  { label: 'Background Jobs', sub: 'asynq · retries',     dot: 'bg-orange-400' },
                  { label: 'Webhooks',      sub: 'Stripe · HMAC · replay', dot: 'bg-emerald-400' },
                  { label: 'Realtime Hub',  sub: 'WebSockets · channels', dot: 'bg-rose-400' },
                  { label: 'Redis Cache',   sub: 'middleware · TTL',      dot: 'bg-red-400' },
                  { label: 'Transactional Mail', sub: 'Resend · templates', dot: 'bg-cyan-400' },
                ].map((f) => (
                  <div
                    key={f.label}
                    className="group relative rounded-xl border border-border/50 bg-card/50 p-4 hover:border-primary/40 hover:bg-card/80 transition-all shadow-sm hover:shadow-[0_8px_24px_-8px_rgba(0,0,0,0.3)]"
                  >
                    <div className="flex items-center gap-2 mb-1">
                      <span className={`h-1.5 w-1.5 rounded-full ${f.dot} shadow-[0_0_8px_currentColor]`} />
                      <span className="font-semibold text-foreground text-sm">{f.label}</span>
                    </div>
                    <p className="text-[11px] text-muted-foreground font-mono">{f.sub}</p>
                  </div>
                ))}
              </div>

            </div>
          </GSAPSection>

        </div>
      </section>

      {/* ═══ HUB & SPOKE — Hubfly-style with animated flow lines ═══ */}
      <HubAndSpoke />

      {/* ═══ FRAMEWORK FOR DEVELOPERS & AGENTS — tabbed code section ═══ */}
      <section className="relative py-24 px-6 overflow-hidden">
        {/* Soft warm gradient backdrop */}
        <div className="absolute inset-0 -z-10 bg-gradient-to-b from-card/30 via-background to-background" />
        <div className="absolute top-20 right-10 w-[600px] h-[600px] -z-10 rounded-full bg-primary/[0.04] blur-[120px]" />

        <div className="max-w-6xl mx-auto">
          <div className="max-w-3xl mx-auto text-center mb-10">
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-4 leading-tight">
              A framework for developers and agents
            </h2>
            <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
              Grit has opinions on everything: routing, queues, auth, storage, AI. That is
              thousands of decisions an AI agent does not have to make &mdash; and the code it
              writes lands in the same shape a person would have written.
            </p>
          </div>

          <div className="flex flex-wrap justify-center gap-x-6 gap-y-2.5 mb-10">
            {[
              'Generates Go + React from one CLI command',
              'Ships a SKILL.md so agents know the patterns',
              'AI Gateway: 100+ models via one API key',
              'OWASP 2025 hardened — secure by default',
            ].map((line) => (
              <span key={line} className="flex items-center gap-2 text-sm text-foreground/80">
                <Check className="h-4 w-4 text-primary shrink-0" strokeWidth={2.5} />
                {line}
              </span>
            ))}
          </div>

          <FeatureTabs />

          <div className="flex justify-center mt-10">
            <Button variant="outline" className="border-border/60 text-foreground hover:bg-accent/30 rounded-full" asChild>
              <Link href="/docs">
                Explore the framework <ArrowRight className="ml-2 h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>
        </div>
      </section>

      {/* ═══ PULSE DASHBOARD + FRONTEND-AGNOSTIC CARDS ═══ */}
      <section className="py-20 px-6 border-t border-border/40">
        <div className="max-w-6xl mx-auto">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-10">

            {/* LEFT: Monitor with Pulse */}
            <div>
              <h2 className="text-2xl md:text-3xl font-bold tracking-tight text-foreground mb-3">
                Monitor and fix issues with Pulse
              </h2>
              <p className="text-base text-muted-foreground leading-relaxed mb-5">
                Pulse gives full observability — find errors and performance issues
                before your team does. Mounted at <code className="text-primary text-sm bg-primary/5 px-1.5 py-0.5 rounded">/pulse/ui</code> on every Grit project.
              </p>
              <ul className="space-y-2.5 mb-6">
                {[
                  'Fix errors and performance with recommended solutions',
                  'Trace requests, jobs, DB queries, cache hits, errors',
                  'Wire k6 test runs into the live latency timeline',
                ].map((line) => (
                  <li key={line} className="flex items-start gap-2.5 text-sm text-foreground/80">
                    <Check className="h-4 w-4 text-primary mt-0.5 shrink-0" strokeWidth={2.5} />
                    {line}
                  </li>
                ))}
              </ul>
              <Button variant="outline" className="border-border/60 text-foreground hover:bg-accent/30 rounded-full mb-8" asChild>
                <Link href="/docs/backend/pulse">
                  Explore Pulse <ArrowRight className="ml-2 h-3.5 w-3.5" />
                </Link>
              </Button>

              {/* Pulse dashboard mockup */}
              <div className="relative rounded-xl border-2 border-border/60 bg-card/80 overflow-hidden shadow-[0_24px_48px_-16px_rgba(0,0,0,0.35),0_2px_8px_-2px_rgba(0,0,0,0.15),0_0_0_4px_hsl(var(--primary)/0.04)]">
                <div className="flex">
                  {/* Sidebar */}
                  <div className="hidden sm:block w-32 border-r border-border/40 bg-background/60 px-3 py-3">
                    <div className="flex items-center gap-1.5 mb-3">
                      <div className="h-5 w-5 rounded-md bg-sky-500/15 flex items-center justify-center">
                        <span className="text-sky-400 font-mono font-bold text-[9px]">P</span>
                      </div>
                      <div className="text-[10px] font-semibold text-foreground/80">Pulse</div>
                    </div>
                    <div className="text-[8px] text-muted-foreground/60 font-mono uppercase tracking-wider mb-1.5">Production</div>
                    {[
                      { l: 'Dashboard', sel: false },
                      { l: 'Requests', sel: true },
                      { l: 'Jobs', sel: false },
                      { l: 'DB Queries', sel: false },
                      { l: 'Errors', sel: false, badge: '12' },
                      { l: 'Slow Queries', sel: false },
                    ].map((row) => (
                      <div key={row.l} className={`flex items-center justify-between rounded-md px-1.5 py-1 mb-0.5 text-[10px] ${row.sel ? 'bg-primary/10 text-primary' : 'text-muted-foreground'}`}>
                        <span>{row.l}</span>
                        {row.badge && <span className="text-[8px] font-mono bg-rose-500/15 text-rose-400 px-1 rounded">{row.badge}</span>}
                      </div>
                    ))}
                  </div>

                  {/* Body */}
                  <div className="flex-1 p-4 min-w-0">
                    <div className="flex items-baseline justify-between mb-1">
                      <div className="text-sm font-semibold text-foreground">Requests</div>
                      <div className="flex items-center gap-1 text-[10px] text-emerald-400">
                        <TrendingUp className="h-2.5 w-2.5" />
                        +14% vs yesterday
                      </div>
                    </div>
                    <div className="flex items-center gap-4 text-[10px] text-muted-foreground/70 mb-3 font-mono">
                      <div><span className="text-foreground/80 font-semibold text-base mr-1">124.2K</span>requests</div>
                      <div className="flex items-center gap-1.5">
                        <span className="inline-block h-1.5 w-1.5 rounded-full bg-emerald-400" />2xx 122.5K
                      </div>
                      <div className="flex items-center gap-1.5">
                        <span className="inline-block h-1.5 w-1.5 rounded-full bg-amber-400" />4xx 1,151
                      </div>
                      <div className="flex items-center gap-1.5">
                        <span className="inline-block h-1.5 w-1.5 rounded-full bg-rose-400" />5xx 324
                      </div>
                    </div>

                    {/* Bar chart */}
                    <div className="flex items-end gap-[3px] h-28 mb-2">
                      {[35, 48, 42, 58, 52, 65, 38, 72, 55, 48, 62, 90, 68, 55, 78, 65, 82, 70, 58, 45, 52, 62, 75, 88, 70, 65].map((h, i) => {
                        const isAlert = i === 11 || i === 17 || i === 23
                        const tone = isAlert
                          ? 'from-rose-400 to-amber-400'
                          : i % 5 === 0
                            ? 'from-amber-400 to-amber-400/60'
                            : 'from-emerald-400/80 to-emerald-400/30'
                        return (
                          <div
                            key={i}
                            className={`flex-1 rounded-sm bg-gradient-to-t ${tone}`}
                            style={{ height: `${h}%` }}
                          />
                        )
                      })}
                    </div>
                    <div className="flex justify-between text-[8px] font-mono text-muted-foreground/50">
                      <span>02 Nov 18:00 UTC</span>
                      <span>03 Nov 18:00 UTC</span>
                    </div>

                    {/* Duration mini-strip */}
                    <div className="mt-3 pt-3 border-t border-border/30 flex items-center justify-between">
                      <div>
                        <div className="text-[10px] text-muted-foreground/70 font-mono uppercase tracking-wider">Duration</div>
                        <div className="text-sm font-semibold text-foreground">125ms — 2.2s</div>
                      </div>
                      <svg className="h-8 w-32" viewBox="0 0 120 30">
                        <polyline
                          fill="none"
                          strokeWidth="1.5"
                          stroke="hsl(var(--primary))"
                          points="0,18 10,15 20,20 30,12 40,16 50,10 60,14 70,8 80,12 90,6 100,10 110,5 120,8"
                        />
                      </svg>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* RIGHT: Frontend-agnostic cascading file cards */}
            <div>
              <h2 className="text-2xl md:text-3xl font-bold tracking-tight text-foreground mb-3">
                The best partner to any front-end
              </h2>
              <p className="text-base text-muted-foreground leading-relaxed mb-5">
                Grit scaffolds Next.js and TanStack Router front-ends, and generates the
                typed client for both. The API itself is plain REST with an OpenAPI
                document, so anything that speaks HTTP can call it &mdash; but these are the
                ones we generate and test.
              </p>
              <Button variant="outline" className="border-border/60 text-foreground hover:bg-accent/30 rounded-full mb-12" asChild>
                <Link href="/docs/frontend">
                  Explore front-ends <ArrowRight className="ml-2 h-3.5 w-3.5" />
                </Link>
              </Button>

              {/* Cascading frontend cards */}
              <div className="relative h-[280px]">
                {[
                  {
                    name: 'users.expo.tsx',
                    color: 'bg-indigo-500/15 text-indigo-300 border-indigo-500/30',
                    icon: <svg viewBox="0 0 24 24" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth="1.6"><rect x="7" y="2" width="10" height="20" rx="2.5" /><path d="M11 18.5h2" /></svg>,
                    style: { top: '0%', right: '0%', width: '78%' },
                  },
                  {
                    name: 'users.tsx',
                    color: 'bg-sky-500/15 text-sky-400 border-sky-500/30',
                    icon: <svg viewBox="0 0 24 24" className="h-4 w-4" fill="currentColor"><circle cx="12" cy="12" r="2.1" /><g fill="none" stroke="currentColor" strokeWidth="1"><ellipse cx="12" cy="12" rx="10" ry="4.5" /><ellipse cx="12" cy="12" rx="10" ry="4.5" transform="rotate(60 12 12)" /><ellipse cx="12" cy="12" rx="10" ry="4.5" transform="rotate(-60 12 12)" /></g></svg>,
                    style: { top: '24%', right: '4%', width: '74%' },
                  },
                  {
                    name: 'users.desktop.tsx',
                    color: 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30',
                    icon: <svg viewBox="0 0 24 24" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth="1.6"><rect x="2.5" y="4" width="19" height="13" rx="2" /><path d="M9 20h6" /></svg>,
                    style: { top: '48%', right: '8%', width: '70%' },
                  },
                  {
                    name: 'users.next.tsx',
                    color: 'bg-violet-500/15 text-violet-300 border-violet-500/30',
                    icon: <svg viewBox="0 0 24 24" className="h-4 w-4" fill="currentColor"><circle cx="12" cy="12" r="10" /><path d="M9 7v10M15 7l-6 10" stroke="white" strokeWidth="1.2" fill="none" /></svg>,
                    style: { top: '72%', right: '12%', width: '66%' },
                  },
                ].map((card) => (
                  <div
                    key={card.name}
                    className={`absolute flex items-center gap-2.5 rounded-xl border ${card.color} bg-card/95 backdrop-blur shadow-lg px-4 py-3`}
                    style={card.style}
                  >
                    {card.icon}
                    <span className="font-mono text-sm font-medium text-foreground/90">{card.name}</span>
                    <div className="flex items-center gap-1 ml-auto">
                      <span className="inline-block h-1 w-8 rounded-full bg-border/60" />
                      <span className="inline-block h-1 w-5 rounded-full bg-border/40" />
                    </div>
                  </div>
                ))}
              </div>
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
              <div key={arch.name} className="rounded-xl border border-border/40 bg-card/50 card-gradient p-5 text-center card-grit hover:border-primary/40">
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
                  { feature: 'OWASP 2025 Hardened', grit: true, next: false, laravel: 'partial', goravel: false },
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

      {/* ═══ TESTIMONIALS ═══

          This section previously carried quotes attributed to real, named people
          and institutions who never gave them — placeholder copy from the design
          mock that shipped by accident.

          It is deliberately empty until there are real ones. Inventing an
          endorsement is both legally actionable and, for a framework asking to be
          trusted with auth and audit logs, self-defeating: nothing else on this
          page survives being caught doing it.

          Add real quotes here as they arrive, with a link to the source. Never
          add a name without written permission. */}
      <section className="relative py-24 px-6 border-t border-border/40 overflow-hidden">
        <div className="absolute inset-0 -z-10 bg-gradient-to-b from-card/40 via-background to-background" />
        <div className="absolute top-40 left-1/4 w-[500px] h-[500px] -z-10 rounded-full bg-primary/[0.04] blur-[120px]" />

        <div className="max-w-3xl mx-auto text-center">
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight text-foreground leading-tight">
            Built something with Grit?
          </h2>
          <p className="mt-5 text-lg text-muted-foreground leading-relaxed">
            There are no testimonials here yet, and there will not be any invented ones.
            If Grit is running something of yours in production, tell us about it — the
            good and the parts that hurt — and it goes on this page with your name and a
            link back to you.
          </p>
          <div className="mt-9 flex flex-col sm:flex-row items-center justify-center gap-3">
            <Link
              href="https://github.com/MUKE-coder/grit/issues/new?template=testimonial.yml"
              target="_blank"
              rel="noreferrer"
              className="inline-flex h-11 items-center gap-2 rounded-lg bg-primary px-6 text-sm font-semibold text-primary-foreground transition-colors hover:bg-primary/90"
            >
              <Github className="h-4 w-4" />
              Share your experience
            </Link>
            <Link
              href="/showcase"
              className="inline-flex h-11 items-center gap-2 rounded-lg border border-border/60 px-6 text-sm font-semibold text-foreground transition-colors hover:bg-accent/50"
            >
              See what people are building
              <ArrowRight className="h-4 w-4" />
            </Link>
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
              <Link key={project.name} href={project.url} target={project.url.startsWith('http') ? '_blank' : undefined} className="group rounded-xl border border-border/40 bg-card/50 card-gradient p-6 card-grit hover:border-primary/40 block">
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

      {/* ═══ SPONSORS ═══ */}
      <section className="relative py-24 px-6 border-t border-border/40 overflow-hidden">
        <div className="pointer-events-none absolute inset-0 -z-10">
          <div className="absolute left-1/2 top-1/2 h-[320px] w-[640px] -translate-x-1/2 -translate-y-1/2 rounded-full bg-primary/[0.05] blur-[120px]" />
        </div>
        <div className="max-w-5xl mx-auto">
          <div className="text-center mb-12">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">Sponsors</p>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Backed by developers like you</h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Grit is free and MIT licensed. Sponsors fund the features, docs and
              releases &mdash; and get their name in front of everyone who builds with it.
            </p>
          </div>

          <HomeSponsors />

          <div className="mt-10 flex flex-col sm:flex-row items-center justify-center gap-3">
            <Button asChild size="lg" className="rounded-full">
              <Link href="/sponsor">
                <Heart className="mr-2 h-4 w-4" />
                Become a sponsor
              </Link>
            </Button>
            <Button asChild variant="ghost" size="lg" className="rounded-full">
              <Link href="/sponsors">
                See all sponsors
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
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

      {/* ═══ COMMUNITY ═══ */}
      <section className="relative py-20 px-6 border-t border-border/40">
        <div className="max-w-5xl mx-auto">
          <CommunityCTA />
        </div>
      </section>

      {/* ═══ CTA — "Start using Grit today" ═══ */}
      <section className="relative py-32 px-6 border-t border-border/40 overflow-hidden">
        {/* faint decorative grid */}
        <div className="absolute inset-0 -z-10 opacity-[0.04]" style={{ backgroundImage: 'linear-gradient(to right, hsl(var(--foreground)) 1px, transparent 1px), linear-gradient(to bottom, hsl(var(--foreground)) 1px, transparent 1px)', backgroundSize: '40px 40px' }} />
        <div className="max-w-3xl mx-auto text-center">
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight text-foreground mb-4">
            Start using Grit today
          </h2>
          <p className="text-muted-foreground mb-8 text-lg">
            Install the CLI and scaffold your first project. Or dive into the docs
            to plan your architecture first.
          </p>
          <CodeBlock language="bash" className="mb-8 text-left" code={`go install github.com/MUKE-coder/grit/v3/cmd/grit@latest
grit new my-app`} />
          <div className="flex flex-col sm:flex-row items-center justify-center gap-3">
            <Button size="lg" className="bg-primary hover:bg-primary/90 text-primary-foreground px-7 h-11 text-sm rounded-full" asChild>
              <Link href="/docs/getting-started/quick-start">Read the docs <ArrowRight className="ml-2 h-3.5 w-3.5" /></Link>
            </Button>
            <Button variant="outline" size="lg" className="border-border/60 text-foreground hover:bg-accent/30 px-7 h-11 text-sm rounded-full" asChild>
              <Link href="/docs/stack-selector">Explore architectures</Link>
            </Button>
          </div>
        </div>
      </section>

      {/* ═══ RICH FOOTER ═══ */}
      <footer className="border-t border-border/40 px-6 pt-16 pb-8 overflow-hidden">
        <div className="max-w-6xl mx-auto">
          <div className="grid grid-cols-2 md:grid-cols-5 gap-8 mb-12">
            {/* Brand column */}
            <div className="col-span-2 md:col-span-2">
              <div className="flex items-center gap-2 mb-4">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img src="/grit_logo.png" alt="Grit" className="h-7 w-7 rounded-md" />
                <span className="font-semibold text-foreground">Grit Framework</span>
              </div>
              <p className="text-sm text-muted-foreground leading-relaxed mb-5 max-w-xs">
                The fastest way to build, deploy, and operate full-stack apps in Go.
              </p>
              <div className="flex items-center gap-3">
                <Link href="https://github.com/MUKE-coder/grit" target="_blank" className="h-8 w-8 rounded-md border border-border/40 flex items-center justify-center text-muted-foreground hover:text-foreground hover:border-border/60 transition-colors">
                  <Github className="h-3.5 w-3.5" />
                </Link>
                <Link href="https://www.youtube.com/@GritFramework" target="_blank" className="h-8 w-8 rounded-md border border-border/40 flex items-center justify-center text-muted-foreground hover:text-foreground hover:border-border/60 transition-colors">
                  <svg className="h-3.5 w-3.5" viewBox="0 0 24 24" fill="currentColor"><path d="M21.582 7.146a2.78 2.78 0 0 0-1.957-1.967C17.882 4.7 12 4.7 12 4.7s-5.882 0-7.625.479A2.78 2.78 0 0 0 2.418 7.146C1.94 8.892 1.94 12 1.94 12s0 3.108.478 4.854a2.78 2.78 0 0 0 1.957 1.967C6.118 19.3 12 19.3 12 19.3s5.882 0 7.625-.479a2.78 2.78 0 0 0 1.957-1.967C22.06 15.108 22.06 12 22.06 12s0-3.108-.478-4.854zM9.94 15.3V8.7l5.715 3.3-5.715 3.3z" /></svg>
                </Link>
                <Link href="https://www.linkedin.com/company/grit-framework" target="_blank" className="h-8 w-8 rounded-md border border-border/40 flex items-center justify-center text-muted-foreground hover:text-foreground hover:border-border/60 transition-colors">
                  <svg className="h-3.5 w-3.5" viewBox="0 0 24 24" fill="currentColor"><path d="M19 3a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h14zM8.339 18.337v-8.59H5.667v8.59h2.672zM7.003 8.574a1.548 1.548 0 1 0 0-3.096 1.548 1.548 0 0 0 0 3.096zm11.335 9.763V13.64c0-2.465-1.338-3.612-3.123-3.612-1.44 0-2.087.791-2.448 1.348v-1.157h-2.671c.034.751 0 8.59 0 8.59h2.671v-4.79c0-.243.018-.487.09-.66.196-.485.642-.989 1.39-.989.982 0 1.376.752 1.376 1.852v4.587h2.715z" /></svg>
                </Link>
              </div>
            </div>

            {/* Products */}
            <div>
              <h4 className="text-xs font-mono font-medium text-foreground/70 tracking-wider mb-4">PRODUCTS</h4>
              <ul className="space-y-2.5 text-sm">
                <li><Link href="/docs" className="text-muted-foreground hover:text-foreground transition-colors">Grit CLI</Link></li>
                <li><Link href="https://ui.gritframework.dev" target="_blank" className="text-muted-foreground hover:text-foreground transition-colors">Grit UI</Link></li>
                <li><Link href="https://gritcms.com" target="_blank" className="text-muted-foreground hover:text-foreground transition-colors">GritCMS</Link></li>
                <li><Link href="https://github.com/MUKE-coder/gorm-studio" target="_blank" className="text-muted-foreground hover:text-foreground transition-colors">GORM Studio</Link></li>
                <li><Link href="/showcase" className="text-muted-foreground hover:text-foreground transition-colors">Showcase</Link></li>
              </ul>
            </div>

            {/* Packages */}
            <div>
              <h4 className="text-xs font-mono font-medium text-foreground/70 tracking-wider mb-4">PACKAGES</h4>
              <ul className="space-y-2.5 text-sm">
                <li><Link href="https://github.com/MUKE-coder/sentinel" target="_blank" className="text-muted-foreground hover:text-foreground transition-colors">Sentinel</Link></li>
                <li><Link href="https://github.com/MUKE-coder/pulse" target="_blank" className="text-muted-foreground hover:text-foreground transition-colors">Pulse</Link></li>
                <li><Link href="https://github.com/MUKE-coder/gin-docs" target="_blank" className="text-muted-foreground hover:text-foreground transition-colors">gin-docs</Link></li>
                <li><Link href="/docs/plugins" className="text-muted-foreground hover:text-foreground transition-colors">Plugins</Link></li>
              </ul>
            </div>

            {/* Resources */}
            <div>
              <h4 className="text-xs font-mono font-medium text-foreground/70 tracking-wider mb-4">RESOURCES</h4>
              <ul className="space-y-2.5 text-sm">
                <li><Link href="/docs" className="text-muted-foreground hover:text-foreground transition-colors">Documentation</Link></li>
                <li><Link href="/docs/getting-started/quick-start" className="text-muted-foreground hover:text-foreground transition-colors">Quick Start</Link></li>
                <li><Link href="/docs/security" className="text-muted-foreground hover:text-foreground transition-colors">Security Guide</Link></li>
                <li><Link href="/docs/testing" className="text-muted-foreground hover:text-foreground transition-colors">Testing</Link></li>
                <li><Link href="/docs/changelog" className="text-muted-foreground hover:text-foreground transition-colors">Changelog</Link></li>
                <li><Link href="/docs/tech-kits" className="text-muted-foreground hover:text-foreground transition-colors">Tech Kits</Link></li>
                <li><Link href="/courses" className="text-muted-foreground hover:text-foreground transition-colors">Courses</Link></li>
                <li><Link href="/sponsors" className="text-muted-foreground hover:text-foreground transition-colors">Sponsors</Link></li>
                <li><Link href="/pitch" className="text-muted-foreground hover:text-foreground transition-colors">Pitch</Link></li>
                <li>
                  <Link href={WHATSAPP_COMMUNITY_URL} target="_blank" rel="noreferrer" className="inline-flex items-center gap-1.5 text-emerald-500/90 hover:text-emerald-500 transition-colors">
                    <WhatsAppIcon className="h-3 w-3" />
                    WhatsApp Community
                  </Link>
                </li>
                <li><Link href="/hire" className="text-muted-foreground hover:text-foreground transition-colors">Hire Us</Link></li>
              </ul>
            </div>
          </div>

          <div className="flex flex-col sm:flex-row items-center justify-between gap-3 pt-8 border-t border-border/40">
            <p className="text-xs text-muted-foreground">© 2026 Grit Framework. MIT licensed.</p>
            <div className="flex items-center gap-5">
              <Link href="/docs" className="text-xs text-muted-foreground hover:text-foreground transition-colors">Docs</Link>
              <Link href="https://github.com/MUKE-coder/grit" target="_blank" className="text-xs text-muted-foreground hover:text-foreground transition-colors">GitHub</Link>
              <Link href="/showcase" className="text-xs text-muted-foreground hover:text-foreground transition-colors">Showcase</Link>
              <span className="text-xs text-muted-foreground/60">{`Built with Grit v${GRIT_VERSION}`}</span>
            </div>
          </div>

          {/* Giant decorative GRIT wordmark */}
          <div className="relative mt-12 -mb-4 select-none pointer-events-none">
            <div className="text-center font-bold tracking-tighter leading-none text-primary/[0.08]" style={{ fontSize: 'clamp(80px, 22vw, 320px)' }}>
              GRIT
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
