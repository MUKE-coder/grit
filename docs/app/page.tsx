import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, Github, Terminal, Layers, Zap, Shield, Database, Bot, Server, Smartphone, ChevronDown, Check, AlertCircle, Upload, TrendingUp } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock } from '@/components/code-block'
import { SoftwareApplicationSchema, FAQPageSchema } from '@/components/structured-data'
import { FeatureTabs } from '@/components/feature-tabs'
import { CpuArchitecture } from '@/components/ui/cpu-architecture'
import { MagneticButton, GSAPSection, FadeIn, GlowOrb } from '@/components/motion-primitives'
import { GoLogo, ReactLogo, VueLogo, SvelteLogo, NextLogo, TanStackLogo, TypeScriptLogo, TailwindLogo, PostgresLogo, RedisLogo, DockerLogo } from '@/components/framework-logos'
import { HubAndSpoke } from '@/components/hub-and-spoke'

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
    <div className="min-h-screen bg-background">
      <SoftwareApplicationSchema />
      <FAQPageSchema />
      <SiteHeader />

      {/* ═══ HERO — saturated block, glass pills, ASCII grid corners, GitHub editor ═══ */}
      <section className="relative overflow-hidden isolate">
        {/* Layer 1 — saturated gradient backdrop */}
        <div
          className="absolute inset-0 -z-30"
          style={{
            background:
              'radial-gradient(ellipse 80% 60% at 50% 0%, #0284c7 0%, #0369a1 35%, #082f49 100%)',
          }}
        />
        {/* Layer 2 — fine grid */}
        <div
          className="absolute inset-0 -z-20 opacity-[0.10]"
          style={{
            backgroundImage:
              'linear-gradient(to right, white 1px, transparent 1px), linear-gradient(to bottom, white 1px, transparent 1px)',
            backgroundSize: '48px 48px',
          }}
        />
        {/* Layer 3 — radial vignette so the grid fades at edges */}
        <div
          className="absolute inset-0 -z-10"
          style={{
            background:
              'radial-gradient(ellipse 70% 50% at 50% 30%, transparent 0%, rgba(8,47,73,0.55) 80%)',
          }}
        />

        {/* ASCII-cube corner decorations */}
        <AsciiCube className="absolute -top-8 -left-8 w-[280px] h-[280px] hidden md:block opacity-[0.18] text-white" />
        <AsciiCube className="absolute -top-12 -right-8 w-[320px] h-[320px] hidden md:block opacity-[0.18] text-white" flip />
        <AsciiCube className="absolute -bottom-16 -left-16 w-[260px] h-[260px] hidden lg:block opacity-[0.15] text-white" />
        <AsciiCube className="absolute -bottom-8 -right-12 w-[280px] h-[280px] hidden lg:block opacity-[0.15] text-white" flip />

        {/* Floating glow orbs */}
        <GlowOrb className="-top-32 left-1/4 h-[400px] w-[400px] bg-sky-400/30" duration={18} />
        <GlowOrb className="top-40 right-1/4 h-[300px] w-[300px] bg-cyan-300/20" delay={2} duration={16} />

        <div className="relative max-w-6xl mx-auto pt-20 pb-16 md:pt-28 md:pb-20 px-6">

          {/* GLASS PILL ROW — backend stack → arrow → frontend stack */}
          <FadeIn delay={0.05}>
            <div className="flex items-center justify-center gap-2 sm:gap-3 mb-12">
              {/* Backend pill */}
              <div className="pill-pulse flex items-center gap-1.5 rounded-full border border-white/25 bg-white/[0.08] backdrop-blur-xl px-2 py-1.5 shadow-[0_8px_32px_-8px_rgba(0,0,0,0.4),inset_0_1px_0_rgba(255,255,255,0.15)]">
                {[
                  { src: '/images/icons/go.svg', alt: 'Go' },
                  { src: '/images/icons/postgressql.png', alt: 'Postgres' },
                  { src: '/images/icons/redis-logo-svgrepo-com.svg', alt: 'Redis' },
                  { src: '/images/icons/docker-svgrepo-com.svg', alt: 'Docker' },
                ].map((logo) => (
                  <div
                    key={logo.alt}
                    title={logo.alt}
                    className="h-8 w-8 rounded-full bg-white shadow-[0_2px_8px_rgba(0,0,0,0.15),inset_0_-1px_0_rgba(0,0,0,0.05)] flex items-center justify-center"
                  >
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img src={logo.src} alt={logo.alt} className="h-5 w-5 object-contain" />
                  </div>
                ))}
              </div>

              {/* Dotted connector with arrowhead */}
              <svg className="h-3 w-12 hidden sm:block" viewBox="0 0 48 12" fill="none">
                <line x1="0" y1="6" x2="40" y2="6" stroke="white" strokeOpacity="0.55" strokeWidth="1.5" strokeDasharray="3 3" className="dash-animate" />
                <path d="M40 2 L46 6 L40 10" stroke="white" strokeOpacity="0.8" strokeWidth="1.5" fill="none" strokeLinecap="round" strokeLinejoin="round" />
              </svg>

              {/* Frontend pill */}
              <div className="pill-pulse flex items-center gap-1.5 rounded-full border border-white/25 bg-white/[0.08] backdrop-blur-xl px-2 py-1.5 shadow-[0_8px_32px_-8px_rgba(0,0,0,0.4),inset_0_1px_0_rgba(255,255,255,0.15)]">
                {[
                  { node: <ReactLogo className="h-5 w-5" />, alt: 'React' },
                  { img: '/images/icons/Next.js.svg', alt: 'Next.js' },
                  { node: <VueLogo className="h-5 w-5" />, alt: 'Vue' },
                  { node: <SvelteLogo className="h-5 w-5" />, alt: 'Svelte' },
                  { img: '/images/icons/tanstack-seeklogo.svg', alt: 'TanStack' },
                ].map((logo, i) => (
                  <div
                    key={i}
                    title={logo.alt}
                    className="h-8 w-8 rounded-full bg-white shadow-[0_2px_8px_rgba(0,0,0,0.15),inset_0_-1px_0_rgba(0,0,0,0.05)] flex items-center justify-center"
                  >
                    {logo.node ?? (
                      // eslint-disable-next-line @next/next/no-img-element
                      <img src={logo.img} alt={logo.alt} className="h-5 w-5 object-contain" />
                    )}
                  </div>
                ))}
              </div>
            </div>
          </FadeIn>

          <FadeIn delay={0.12}>
            <h1 className="text-4xl md:text-6xl lg:text-[5.5rem] font-bold tracking-tight text-white text-center mb-6 leading-[1.05]"
                style={{ textShadow: '0 4px 20px rgba(0,0,0,0.15)' }}>
              Build full-stack apps<br />
              <span className="bg-gradient-to-b from-white via-white to-sky-100 bg-clip-text text-transparent">
                with the backend you trust
              </span>
            </h1>
          </FadeIn>

          <FadeIn delay={0.18}>
            <p className="text-base md:text-lg text-white/85 max-w-2xl mx-auto text-center mb-10 leading-relaxed">
              Scaffold a Go API + React frontend + admin panel in one command.
              Auth, OAuth, file storage, jobs, AI, observability, OWASP-2025 hardened
              — meticulously optimized for production. No boilerplate required.
            </p>
          </FadeIn>

          {/* Premium CTAs — magnetic + glass */}
          <FadeIn delay={0.24}>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-3 mb-16">
              <MagneticButton>
                <Link
                  href="/docs/getting-started/quick-start"
                  className="group relative inline-flex items-center justify-center gap-2 h-12 px-7 rounded-full bg-white text-sky-700 font-semibold text-sm shadow-[0_8px_24px_-4px_rgba(0,0,0,0.3),0_2px_4px_rgba(0,0,0,0.1),inset_0_-2px_0_rgba(2,132,199,0.15)] hover:shadow-[0_12px_32px_-4px_rgba(0,0,0,0.35),0_2px_4px_rgba(0,0,0,0.1),inset_0_-2px_0_rgba(2,132,199,0.2)] transition-all"
                >
                  Get started
                  <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
                </Link>
              </MagneticButton>
              <MagneticButton>
                <Link
                  href="/docs"
                  className="inline-flex items-center justify-center h-12 px-7 rounded-full border border-white/25 bg-white/[0.08] backdrop-blur-xl text-white font-medium text-sm hover:bg-white/[0.14] hover:border-white/35 shadow-[0_8px_24px_-8px_rgba(0,0,0,0.3),inset_0_1px_0_rgba(255,255,255,0.1)] transition-all"
                >
                  Read docs
                </Link>
              </MagneticButton>
            </div>
          </FadeIn>

          {/* INSTALL pill — one-line copyable command, lives right above the editor */}
          <FadeIn delay={0.28}>
            <div className="max-w-2xl mx-auto mb-10 mt-2">
              <div className="rounded-xl border border-white/15 bg-white/[0.05] backdrop-blur-xl shadow-[0_8px_32px_-8px_rgba(0,0,0,0.4),inset_0_1px_0_rgba(255,255,255,0.08)]">
                <CodeBlock
                  terminal
                  filename="Quick install"
                  code={`go install github.com/MUKE-coder/grit/v3/cmd/grit@latest && grit new my-app`}
                  className="!border-0 !rounded-xl !bg-transparent dark:!bg-transparent !m-0"
                />
              </div>
            </div>
          </FadeIn>

          {/* SIDE-BY-SIDE GITHUB EDITOR — bold border, layered shadow, file tabs row at top */}
          <FadeIn delay={0.32}>
            <div className="relative rounded-2xl overflow-hidden bg-[#ffffff] dark:bg-[#0d1117] border-2 border-white/30 shadow-[0_24px_64px_-16px_rgba(2,6,23,0.6),0_8px_24px_-8px_rgba(0,0,0,0.4),inset_0_1px_0_rgba(255,255,255,0.15)]">
              {/* Editor chrome — file tab strip */}
              <div className="flex items-center gap-0 bg-[#f6f8fa] dark:bg-[#161b22] border-b border-[#d0d7de] dark:border-white/[0.08]">
                {/* Window dots */}
                <div className="flex items-center gap-1.5 px-4 py-2.5 border-r border-[#d0d7de]/60 dark:border-white/[0.06]">
                  <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                  <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                  <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                </div>
                {/* Active file tabs */}
                <div className="flex items-center gap-2 px-3 py-2.5 border-r border-[#d0d7de]/60 dark:border-white/[0.06] bg-white dark:bg-[#0d1117] -mb-px">
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img src="/images/icons/go.svg" alt="" className="h-3.5 w-3.5" />
                  <span className="text-[12px] font-mono text-[#24292f] dark:text-slate-300">internal/handlers/product.go</span>
                </div>
                <div className="flex items-center gap-2 px-3 py-2.5">
                  <ReactLogo className="h-3.5 w-3.5" />
                  <span className="text-[12px] font-mono text-[#57606a] dark:text-slate-500">frontend/src/hooks/use-products.ts</span>
                  <span className="ml-1 text-[9px] font-mono text-[#57606a]/70 dark:text-slate-500/70 uppercase tracking-wider">tsx</span>
                </div>
                <div className="ml-auto text-[10px] font-mono text-[#57606a] dark:text-slate-500 px-4 hidden md:block">
                  generated · 2 files
                </div>
              </div>

              {/* Two-pane code body — GitHub theme, light/dark aware */}
              <div className="grid md:grid-cols-2 divide-y md:divide-y-0 md:divide-x divide-[#d0d7de] dark:divide-white/[0.06] bg-white dark:bg-[#0d1117]">
                <div className="text-left">
                  <CodeBlock
                    code={`package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "myapp/internal/authz"
)

func (h *ProductHandler) List(c *gin.Context) {
    var products []models.Product
    h.DB.
        Where("user_id = ?", c.GetString("user_id")).
        Find(&products)

    c.JSON(http.StatusOK, gin.H{
        "data": products,
    })
}`}
                    language="go"
                    className="!border-0 !rounded-none !shadow-none !bg-transparent dark:!bg-transparent !m-0"
                  />
                </div>
                <div className="text-left">
                  <CodeBlock
                    code={`import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api'

export function useProducts() {
  return useQuery({
    queryKey: ['products'],
    queryFn: async () => {
      const res = await api
        .get('/api/products')
      return res.data.data
    },
  })
}`}
                    language="tsx"
                    className="!border-0 !rounded-none !shadow-none !bg-transparent dark:!bg-transparent !m-0"
                  />
                </div>
              </div>

              {/* Footer attribution strip */}
              <div className="flex items-center justify-center gap-2 px-4 py-2.5 bg-[#f6f8fa] dark:bg-[#161b22] border-t border-[#d0d7de] dark:border-white/[0.08]">
                <Terminal className="h-3 w-3 text-[#57606a] dark:text-slate-500" />
                <span className="text-[11px] font-mono text-[#57606a] dark:text-slate-400">
                  Both files generated by{' '}
                  <span className="text-sky-700 dark:text-sky-400 font-semibold">grit generate resource Product</span>
                </span>
              </div>
            </div>
          </FadeIn>

        </div>
      </section>

      {/* ═══ FEATURES INTRO + STATS BY NUMBERS ═══ */}
      <section className="py-20 px-6 border-b border-border/40">
        <div className="max-w-5xl mx-auto text-center">
          <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
            <span className="h-1.5 w-1.5 rounded-full bg-primary" /> Features
          </span>
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight text-foreground mb-5">
            Elevate your stack<br />
            with the modern Go monolith
          </h2>
          <p className="text-base md:text-lg text-muted-foreground max-w-2xl mx-auto leading-relaxed">
            Grit ships with production-ready features that accelerate development
            and make building full-stack apps in Go a breeze.
          </p>
        </div>
      </section>

      {/* "BY NUMBERS" stats row — Inertia-style minimal */}
      <section className="border-b border-border/40">
        <div className="max-w-6xl mx-auto px-6 py-12">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8 items-center">
            <div className="flex flex-col">
              <div className="flex items-center gap-2 mb-1">
                <div className="h-7 w-7 rounded-md bg-primary/10 flex items-center justify-center">
                  <span className="text-primary font-mono font-bold text-xs">G</span>
                </div>
                <span className="text-xs font-mono font-medium text-foreground/70 tracking-wider">GRIT<br className="hidden lg:block" /> BY NUMBERS</span>
              </div>
            </div>
            {[
              { number: '5', label: 'ARCHITECTURE MODES' },
              { number: '100+', label: 'UI COMPONENTS' },
              { number: 'v3.23', label: 'LATEST RELEASE' },
            ].map((stat) => (
              <div key={stat.label} className="text-center md:text-left">
                <div className="text-4xl md:text-5xl font-bold text-foreground mb-1 tracking-tight">{stat.number}</div>
                <div className="text-[10px] font-mono font-medium text-muted-foreground tracking-widest">{stat.label}</div>
              </div>
            ))}
          </div>
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
                      <span className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground">v3.23 · production-ready</span>
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
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-10 mb-12 items-start">
            <div className="lg:col-span-4">
              <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-4 leading-tight">
                A framework for<br />developers and agents
              </h2>
              <p className="text-base text-muted-foreground leading-relaxed mb-6">
                Grit has opinions on everything: routing, queues, auth, storage, AI.
                That&apos;s thousands of decisions an AI agent doesn&apos;t have to make.
                The result? Clean Go code that anyone — human or assistant — can extend.
              </p>
              <ul className="space-y-3 mb-7">
                {[
                  'Generates Go + React from one CLI command',
                  'Ships an SKILL.md so agents know the patterns',
                  'AI Gateway: 100+ models via one API key',
                  'OWASP 2025 hardened — secure by default',
                ].map((line) => (
                  <li key={line} className="flex items-start gap-2.5 text-sm text-foreground/80">
                    <Check className="h-4 w-4 text-primary mt-0.5 shrink-0" strokeWidth={2.5} />
                    {line}
                  </li>
                ))}
              </ul>
              <Button variant="outline" className="border-border/60 text-foreground hover:bg-accent/30 rounded-full" asChild>
                <Link href="/docs">
                  Explore the framework <ArrowRight className="ml-2 h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>

            <div className="lg:col-span-8">
              <FeatureTabs />
            </div>
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
                Easily craft frontend experiences with React, TanStack Router, Vue, or
                Svelte alongside the Grit API. Or accelerate development with a generated
                Next.js admin panel.
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
                    name: 'users.svelte',
                    color: 'bg-orange-500/15 text-orange-400 border-orange-500/30',
                    icon: <svg viewBox="0 0 24 24" className="h-4 w-4" fill="currentColor"><path d="M19.6 5.2C17.6 2.6 14.1 2.1 11.4 4l-4.5 2.9c-1.3.8-2.2 2.1-2.5 3.6-.3 1.2-.1 2.4.4 3.5-.4.5-.6 1.2-.8 1.8-.3 1.4-.1 2.8.6 4 1.9 2.6 5.4 3.2 8.1 1.3l4.5-2.9c1.3-.8 2.2-2.1 2.5-3.6.3-1.2.1-2.4-.4-3.5.4-.5.6-1.2.8-1.8.3-1.4.1-2.8-.6-4Z" /></svg>,
                    style: { top: '0%', right: '0%', width: '78%' },
                  },
                  {
                    name: 'users.tsx',
                    color: 'bg-sky-500/15 text-sky-400 border-sky-500/30',
                    icon: <svg viewBox="0 0 24 24" className="h-4 w-4" fill="currentColor"><circle cx="12" cy="12" r="2.1" /><g fill="none" stroke="currentColor" strokeWidth="1"><ellipse cx="12" cy="12" rx="10" ry="4.5" /><ellipse cx="12" cy="12" rx="10" ry="4.5" transform="rotate(60 12 12)" /><ellipse cx="12" cy="12" rx="10" ry="4.5" transform="rotate(-60 12 12)" /></g></svg>,
                    style: { top: '24%', right: '4%', width: '74%' },
                  },
                  {
                    name: 'users.vue',
                    color: 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30',
                    icon: <svg viewBox="0 0 24 24" className="h-4 w-4" fill="currentColor"><path d="M12 2 2 22h6l4-7 4 7h6L12 2Zm0 4 6 10h-3l-3-5-3 5H6L12 6Z" /></svg>,
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

      {/* ═══ 3-COLUMN FEATURE GRID — illustrated cards ═══ */}
      <section className="py-20 px-6">
        <div className="max-w-6xl mx-auto">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">

            {/* Multi-architecture */}
            <div>
              <h3 className="font-semibold text-foreground mb-2">5 architectures, one CLI</h3>
              <p className="text-sm text-muted-foreground mb-5 leading-relaxed">
                Pick the scaffold that fits — embed your SPA in the binary, or split web /
                admin / API into a monorepo. Same generators across all of them.
              </p>
              <div className="rounded-xl border border-border/40 bg-card/50 p-6">
                <div className="flex items-center justify-center h-32 relative">
                  <div className="absolute inset-0 flex items-center justify-center">
                    <div className="grid grid-cols-3 gap-2.5">
                      {['single', 'double', 'triple', 'api', 'mobile', 'desktop'].map((name, i) => (
                        <div key={name} className={`h-10 w-10 rounded-xl border ${i === 2 ? 'border-primary/40 bg-primary/15 ring-2 ring-primary/20' : 'border-border/40 bg-background'} flex items-center justify-center`}>
                          {i === 0 && <Zap className="h-4 w-4 text-sky-400" />}
                          {i === 1 && <Layers className="h-4 w-4 text-violet-400" />}
                          {i === 2 && <Server className="h-4 w-4 text-emerald-400" />}
                          {i === 3 && <Database className="h-4 w-4 text-amber-400" />}
                          {i === 4 && <Smartphone className="h-4 w-4 text-rose-400" />}
                          {i === 5 && <Terminal className="h-4 w-4 text-cyan-400" />}
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* More than just Go */}
            <div>
              <h3 className="font-semibold text-foreground mb-2">More than just Go</h3>
              <p className="text-sm text-muted-foreground mb-5 leading-relaxed">
                Grit works with Gin, GORM, asynq, Resend, Sentinel, Pulse, and a curated
                stack designed to play together — so you always have what your team adores.
              </p>
              <div className="rounded-xl border border-border/40 bg-card/50 p-6">
                <div className="grid grid-cols-3 gap-2.5">
                  {[
                    { c: 'bg-sky-500/15 text-sky-400', l: 'Go' },
                    { c: 'bg-cyan-500/15 text-cyan-400', l: 'Gin' },
                    { c: 'bg-rose-500/15 text-rose-400', l: 'GORM' },
                    { c: 'bg-violet-500/15 text-violet-400', l: 'PG' },
                    { c: 'bg-red-500/15 text-red-400', l: 'Redis' },
                    { c: 'bg-amber-500/15 text-amber-400', l: 'R2' },
                    { c: 'bg-emerald-500/15 text-emerald-400', l: 'JWT' },
                    { c: 'bg-pink-500/15 text-pink-400', l: 'Resend' },
                    { c: 'bg-indigo-500/15 text-indigo-400', l: 'AI' },
                  ].map((t) => (
                    <div key={t.l} className={`h-12 rounded-xl ${t.c} flex items-center justify-center font-mono font-bold text-xs`}>
                      {t.l}
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Flexible frontend */}
            <div>
              <h3 className="font-semibold text-foreground mb-2">Flexible frontend</h3>
              <p className="text-sm text-muted-foreground mb-5 leading-relaxed">
                Ship a Next.js SaaS or a Vite + TanStack Router SPA. Generated hooks +
                types match either, so swapping the frontend never rewrites your API.
              </p>
              <div className="rounded-xl border border-border/40 bg-card/50 p-6">
                <div className="space-y-2.5">
                  {[
                    { name: 'users.tsx', sel: false, ico: '⚛' },
                    { name: 'use-users.ts', sel: false, ico: '⚛' },
                    { name: 'products.tsx', sel: true, ico: '⚛' },
                  ].map((f, i) => (
                    <div key={i} className={`flex items-center gap-2 rounded-lg px-3 py-2 text-xs font-mono ${f.sel ? 'bg-primary/10 border border-primary/30 text-primary' : 'bg-background border border-border/40 text-muted-foreground'}`}>
                      <span className={f.sel ? 'text-primary' : 'text-cyan-400/60'}>{f.ico}</span>
                      {f.name}
                    </div>
                  ))}
                </div>
              </div>
            </div>

          </div>
        </div>
      </section>

      {/* ═══ FORMS / SECURITY / PERFORMANCE — Inertia-style mockup cards ═══ */}
      <section className="py-12 px-6">
        <div className="max-w-6xl mx-auto">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">

            {/* Forms */}
            <div>
              <h3 className="font-semibold text-foreground mb-2">Form Builder</h3>
              <p className="text-sm text-muted-foreground mb-5 leading-relaxed">
                Grit streamlines form management with simple submissions, intuitive slots
                and props, and fully typed event handlers for a seamless experience.
              </p>
              <div className="rounded-xl border border-border/40 bg-card/50 p-6">
                <div className="space-y-2">
                  {[
                    { icon: <Check className="h-3.5 w-3.5 text-emerald-400" />, label: 'Validation' },
                    { icon: <Upload className="h-3.5 w-3.5 text-sky-400" />, label: 'File Uploads' },
                    { icon: <AlertCircle className="h-3.5 w-3.5 text-rose-400" />, label: 'Error Handling' },
                    { icon: <Check className="h-3.5 w-3.5 text-emerald-400" />, label: '16+ Field Types' },
                  ].map((row) => (
                    <div key={row.label} className="flex items-center gap-3 rounded-lg border border-border/40 bg-background px-4 py-2.5">
                      <span className="h-5 w-5 rounded-full bg-card flex items-center justify-center">{row.icon}</span>
                      <span className="text-sm text-foreground/80">{row.label}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Security / Performance / Observability — three score circles */}
            <div>
              <h3 className="font-semibold text-foreground mb-2">Production-grade by default</h3>
              <p className="text-sm text-muted-foreground mb-5 leading-relaxed">
                Enjoy out-of-the-box security headers, percentile latency tracking, and a
                Lighthouse-ready frontend the moment you run <code>grit new</code>.
              </p>
              <div className="rounded-xl border border-border/40 bg-card/50 p-6">
                <div className="grid grid-cols-3 gap-3">
                  {[
                    { label: 'SECURITY', score: '100', tone: 'text-emerald-400', ring: 'stroke-emerald-400', sub: 'OWASP 2025' },
                    { label: 'PERFORMANCE', score: '98', tone: 'text-amber-400', ring: 'stroke-amber-400', sub: 'Pulse SLO' },
                    { label: 'OBSERVABILITY', score: '100', tone: 'text-sky-400', ring: 'stroke-sky-400', sub: 'p50/p95/p99' },
                  ].map((s) => (
                    <div key={s.label} className="flex flex-col items-center">
                      <div className="relative h-16 w-16 mb-2">
                        <svg className="h-16 w-16 -rotate-90" viewBox="0 0 36 36">
                          <circle cx="18" cy="18" r="15.9" fill="none" className="stroke-border/40" strokeWidth="2.5" />
                          <circle cx="18" cy="18" r="15.9" fill="none" className={s.ring} strokeWidth="2.5" strokeDasharray={`${s.score}, 100`} strokeLinecap="round" />
                        </svg>
                        <div className={`absolute inset-0 flex items-center justify-center font-bold text-base ${s.tone}`}>{s.score}</div>
                      </div>
                      <div className="text-[9px] font-mono font-medium text-foreground/80 tracking-wider">{s.label}</div>
                      <div className="text-[9px] text-muted-foreground mt-0.5">{s.sub}</div>
                    </div>
                  ))}
                </div>
              </div>
            </div>

          </div>
        </div>
      </section>

      {/* ═══ 6-CELL TEXT FEATURE LIST — Inertia 3×2 ═══ */}
      <section className="py-12 px-6">
        <div className="max-w-6xl mx-auto">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-x-8 gap-y-8">
            {[
              { title: 'Background jobs', desc: 'Redis-backed asynq queue. Image processing, email sending, scheduled cleanup — wired in.' },
              { title: 'Realtime hub', desc: 'WebSocket fan-out at /api/ws. SendToUser and Broadcast helpers; useRealtimeEvent hook on the client.' },
              { title: 'Idempotency-Key', desc: 'Auto-attached on unsafe methods. Replay the original 2xx response for 24h — safe retries everywhere.' },
              { title: 'CSV / Excel export', desc: 'Every generated resource ships /export?format=xlsx. Streaming CSV, chunked XLSX, constant memory.' },
              { title: 'PDF generation', desc: 'internal/pdf with a worked Invoice template — typeset, branded, ready to email or download.' },
              { title: 'Feature flags', desc: 'In-memory engine, sticky bucketing, percentage rollouts, allow/blocklists, realtime admin push.' },
            ].map((item) => (
              <div key={item.title}>
                <h4 className="font-semibold text-foreground text-sm mb-2">{item.title}</h4>
                <p className="text-sm text-muted-foreground leading-relaxed">{item.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ═══ 3-COLUMN VISUAL FEATURE CARDS — Inertia bottom row ═══ */}
      <section className="py-12 px-6 mb-8">
        <div className="max-w-6xl mx-auto">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">

            {/* Offline-first */}
            <div>
              <h3 className="font-semibold text-foreground mb-2">Offline-first sync</h3>
              <p className="text-sm text-muted-foreground mb-5 leading-relaxed">
                Local SQLite mirror + outbox with squash semantics. Click Sync, resolve
                field-level conflicts, push.
              </p>
              <div className="rounded-xl border border-border/40 bg-card/50 p-6">
                <Button className="w-full rounded-full bg-primary/90 hover:bg-primary text-primary-foreground" size="sm">
                  Sync now
                </Button>
                <div className="mt-3 flex items-center justify-between text-[11px] text-muted-foreground">
                  <span>Last synced</span>
                  <span className="font-mono">12s ago</span>
                </div>
              </div>
            </div>

            {/* Audit chain */}
            <div>
              <h3 className="font-semibold text-foreground mb-2">Tamper-evident audit log</h3>
              <p className="text-sm text-muted-foreground mb-5 leading-relaxed">
                Every mutation is appended to a SHA-256 hash chain. Verify integrity from
                the admin in one click.
              </p>
              <div className="rounded-xl border border-border/40 bg-card/50 p-6 space-y-2">
                {[
                  { hash: 'a3f1...c92', user: 'maya@acme.com', action: 'POST /api/invoices' },
                  { hash: 'b9d2...44e', user: 'jb@grit.dev', action: 'PUT /api/users/42' },
                  { hash: 'c1e7...8af', user: 'maya@acme.com', action: 'DELETE /api/blogs/9' },
                ].map((row) => (
                  <div key={row.hash} className="flex items-center justify-between text-[11px]">
                    <code className="text-emerald-400/80 font-mono">{row.hash}</code>
                    <span className="text-muted-foreground/70 truncate">{row.action}</span>
                  </div>
                ))}
                <div className="mt-2 pt-2 border-t border-border/30 flex items-center gap-1.5 text-[10px] text-emerald-400 font-mono">
                  <Check className="h-3 w-3" /> integrity verified
                </div>
              </div>
            </div>

            {/* Infinite scrolling / pagination */}
            <div>
              <h3 className="font-semibold text-foreground mb-2">Cursor pagination</h3>
              <p className="text-sm text-muted-foreground mb-5 leading-relaxed">
                Sticky-page pagination on every generated list endpoint. No skipped rows
                when data shifts mid-scroll.
              </p>
              <div className="rounded-xl border border-border/40 bg-card/50 p-6 space-y-2">
                {[
                  { id: '#1284', name: 'Invoice — Acme Inc', tone: 'text-emerald-400 bg-emerald-400/10', tag: 'Paid' },
                  { id: '#1283', name: 'Invoice — Globex', tone: 'text-amber-400 bg-amber-400/10', tag: 'Pending' },
                  { id: '#1282', name: 'Invoice — Initech', tone: 'text-rose-400 bg-rose-400/10', tag: 'Overdue' },
                ].map((row) => (
                  <div key={row.id} className="flex items-center justify-between text-[11px]">
                    <div className="flex items-center gap-2">
                      <span className="font-mono text-muted-foreground/70">{row.id}</span>
                      <span className="text-foreground/80 truncate">{row.name}</span>
                    </div>
                    <span className={`px-1.5 py-0.5 rounded text-[10px] font-medium ${row.tone}`}>{row.tag}</span>
                  </div>
                ))}
                <div className="mt-2 pt-2 border-t border-border/30 text-[10px] text-muted-foreground font-mono">
                  has_more · cursor: eyJ0ID...
                </div>
              </div>
            </div>

          </div>
        </div>
      </section>

      {/* ═══ BENTO GRID — original auth/admin/AI/storage/codegen ═══ */}
      <section className="py-24 px-6 border-t border-border/40">
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
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">New in v3.23</p>
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">OWASP 2025 hardened</h2>
            <p className="text-muted-foreground max-w-2xl mx-auto">
              Reliability, auditability, and security — usually a year of integration work — wired into every scaffolded project.
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

            {/* SSRF defence */}
            <div className="rounded-xl border border-border/40 bg-card/50 card-gradient grid-pattern p-6 hover:border-primary/30 transition-all duration-300">
              <div className="h-10 w-10 rounded-lg bg-rose-500/10 flex items-center justify-center mb-4">
                <Shield className="h-5 w-5 text-rose-400" />
              </div>
              <h3 className="font-semibold text-foreground mb-2">SSRF + IDOR closed</h3>
              <p className="text-sm text-muted-foreground mb-2">
                <code className="text-primary text-xs">safefetch.Get</code> blocks private/IMDS IPs with DNS-rebind defence.{' '}
                <code className="text-primary text-xs">authz.MustOwn</code> returns 404 to prevent enumeration.
              </p>
              <code className="text-[11px] text-primary/80 font-mono">OWASP A01 closed</code>
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

            {/* Realtime + idempotency — wider card */}
            <div className="lg:col-span-2 rounded-xl border border-border/40 bg-card/50 card-gradient grid-pattern p-6 hover:border-primary/30 transition-all duration-300">
              <div className="flex items-center gap-3 mb-4">
                <div className="h-10 w-10 rounded-lg bg-amber-500/10 flex items-center justify-center">
                  <svg className="h-5 w-5 text-amber-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" /></svg>
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">Webhook receiver + k6 load test suite + security CI</h3>
                  <p className="text-xs text-muted-foreground">The &quot;every business app needs this&quot; primitives, baked in</p>
                </div>
              </div>
              <ul className="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-2 text-sm text-muted-foreground">
                <li>• <strong className="text-foreground/80">Stripe / GitHub / HMAC verifiers</strong> with auto-dedup on (provider, external_id)</li>
                <li>• <strong className="text-foreground/80">k6 test suite</strong> — smoke / load / stress / spike / soak / breakpoint</li>
                <li>• <strong className="text-foreground/80">Dependabot + govulncheck + CodeQL</strong> CI on every PR</li>
                <li>• <strong className="text-foreground/80">CSRF middleware</strong> double-submit cookie for OAuth flows</li>
                <li>• <strong className="text-foreground/80">Security event log</strong> with typed events + hash chain</li>
                <li>• <strong className="text-foreground/80">CSP / COOP / CORP</strong> headers strict by default</li>
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

      {/* ═══ TESTIMONIALS — asymmetric grid with dark accent cards ═══ */}
      <section className="relative py-24 px-6 border-t border-border/40 overflow-hidden">
        <div className="absolute inset-0 -z-10 bg-gradient-to-b from-card/40 via-background to-background" />
        <div className="absolute top-40 left-1/4 w-[500px] h-[500px] -z-10 rounded-full bg-primary/[0.04] blur-[120px]" />

        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-5xl font-bold tracking-tight text-foreground leading-tight">
              Trusted by builders<br />all over the world
            </h2>
          </div>

          {/* Asymmetric 3-column grid */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 md:gap-5">

            {/* Tall dark card #1 — Sentinel "logo" treatment */}
            <div className="md:row-span-2 rounded-2xl bg-slate-950 border border-slate-800/80 p-7 flex flex-col">
              <div className="flex items-center gap-2 mb-7">
                <div className="h-7 w-7 rounded-md bg-rose-500/20 flex items-center justify-center">
                  <Shield className="h-3.5 w-3.5 text-rose-400" />
                </div>
                <span className="font-semibold text-white text-sm tracking-tight">sentinel</span>
              </div>
              <p className="text-2xl font-medium text-white leading-tight tracking-tight mb-auto">
                &ldquo;Grit&apos;s code generator and Sentinel
                integration meant we shipped a secure
                WAF + audit dashboard the same week
                we started the project.&rdquo;
              </p>
              <div className="flex items-center gap-3 mt-7">
                <div className="h-10 w-10 rounded-full bg-gradient-to-br from-rose-400 to-orange-500 flex items-center justify-center font-bold text-white">AC</div>
                <div>
                  <div className="font-semibold text-white text-sm">Alex Chen</div>
                  <div className="text-xs text-white/50">Founder, Skywatcher</div>
                </div>
              </div>
            </div>

            {/* Light cards */}
            <div className="rounded-2xl border border-border/40 bg-card/60 p-6 flex flex-col">
              <p className="text-base text-foreground/85 leading-relaxed flex-1 mb-5">
                &ldquo;Grit is our sourdough starter and multitool for full-stack projects
                large and small. The single-app mode in particular is fresh and useful.&rdquo;
              </p>
              <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-full bg-gradient-to-br from-emerald-400 to-cyan-500 flex items-center justify-center font-bold text-white">IC</div>
                <div>
                  <div className="font-semibold text-foreground text-sm">Ian Callahan</div>
                  <div className="text-xs text-muted-foreground">Harvard Art Museums</div>
                </div>
              </div>
            </div>

            <div className="rounded-2xl border border-border/40 bg-card/60 p-6 flex flex-col">
              <p className="text-base text-foreground/85 leading-relaxed flex-1 mb-5">
                &ldquo;Grit takes the pain out of building modern, scalable Go web apps. The
                tabbed code generator is a love-letter to senior devs.&rdquo;
              </p>
              <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-full bg-gradient-to-br from-amber-400 to-rose-500 flex items-center justify-center font-bold text-white">AF</div>
                <div>
                  <div className="font-semibold text-foreground text-sm">Aaron Francis</div>
                  <div className="text-xs text-muted-foreground">Co-founder, Try Hard Studios</div>
                </div>
              </div>
            </div>

            <div className="rounded-2xl border border-border/40 bg-card/60 p-6 flex flex-col">
              <p className="text-base text-foreground/85 leading-relaxed flex-1 mb-5">
                &ldquo;Grit&apos;s elegance, performance, and developer experience are unmatched
                for Go. The generated code is clean enough to teach from.&rdquo;
              </p>
              <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-full bg-gradient-to-br from-violet-400 to-indigo-500 flex items-center justify-center font-bold text-white">CP</div>
                <div>
                  <div className="font-semibold text-foreground text-sm">Chandresh Patel</div>
                  <div className="text-xs text-muted-foreground">CEO, Bacancy</div>
                </div>
              </div>
            </div>

            {/* Tall dark card #2 — Pulse "logo" treatment */}
            <div className="md:row-span-2 rounded-2xl bg-slate-950 border border-slate-800/80 p-7 flex flex-col">
              <div className="flex items-center gap-2 mb-7">
                <div className="h-7 w-7 rounded-md bg-sky-500/20 flex items-center justify-center">
                  <TrendingUp className="h-3.5 w-3.5 text-sky-400" />
                </div>
                <span className="font-semibold text-white text-sm tracking-tight">pulse</span>
              </div>
              <p className="text-2xl font-medium text-white leading-tight tracking-tight mb-auto">
                &ldquo;The Grit ecosystem has been
                integral to the success of our
                product. The framework lets us
                move fast and ship regularly
                without dropping a single SLO.&rdquo;
              </p>
              <div className="flex items-center gap-3 mt-7">
                <div className="h-10 w-10 rounded-full bg-gradient-to-br from-sky-400 to-violet-500 flex items-center justify-center font-bold text-white">JE</div>
                <div>
                  <div className="font-semibold text-white text-sm">Jack Ellis</div>
                  <div className="text-xs text-white/50">Founder, Fathom Analytics</div>
                </div>
              </div>
            </div>

            <div className="rounded-2xl border border-border/40 bg-card/60 p-6 flex flex-col">
              <p className="text-base text-foreground/85 leading-relaxed flex-1 mb-5">
                &ldquo;Grit is a breath of fresh air in the Go ecosystem, with a brilliant
                community around it. Generators that actually feel like Rails.&rdquo;
              </p>
              <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-full bg-gradient-to-br from-pink-400 to-rose-500 flex items-center justify-center font-bold text-white">EH</div>
                <div>
                  <div className="font-semibold text-foreground text-sm">Erika Heidi</div>
                  <div className="text-xs text-muted-foreground">Creator, Minicli</div>
                </div>
              </div>
            </div>

            <div className="rounded-2xl border border-border/40 bg-card/60 p-6 flex flex-col">
              <p className="text-base text-foreground/85 leading-relaxed flex-1 mb-5">
                &ldquo;The framework, the ecosystem, and the community — it&apos;s the perfect
                package for shipping production Go apps.&rdquo;
              </p>
              <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-full bg-gradient-to-br from-emerald-400 to-teal-500 flex items-center justify-center font-bold text-white">ZK</div>
                <div>
                  <div className="font-semibold text-foreground text-sm">Zuzana Kunckova</div>
                  <div className="text-xs text-muted-foreground">Founder, GoBuilders</div>
                </div>
              </div>
            </div>

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
                <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10">
                  <span className="text-primary font-mono font-bold text-xs">G</span>
                </div>
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
                <li><Link href="/courses" className="text-muted-foreground hover:text-foreground transition-colors">Courses</Link></li>
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
              <span className="text-xs text-muted-foreground/60">Built with Grit v3.23</span>
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
