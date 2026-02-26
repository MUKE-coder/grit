import Link from 'next/link'
import Image from 'next/image'
import type { Metadata } from 'next'
import { ArrowRight, Shield, Layers, Zap, Code2, Rocket, Github, Globe, Server, Gauge, CheckCircle2, X, Minus, Lock, HardDrive, LayoutDashboard, FileText, Database, Mail, Bot, Activity, Wand2, Table2, ClipboardList, ShieldCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { SiteHeader } from '@/components/site-header'
import { AnimatedTerminal } from '@/components/animated-terminal'
import { SoftwareApplicationSchema, FAQPageSchema } from '@/components/structured-data'

export const metadata: Metadata = {
  title: 'Grit — Go + React Full-Stack Framework',
  description:
    'Grit is a full-stack meta-framework that combines Go (Gin + GORM) with React (Next.js) and a Filament-like admin panel. One CLI command scaffolds an entire monorepo with API, web app, admin dashboard, shared types, and Docker setup.',
  alternates: { canonical: 'https://gritframework.dev' },
}

export default function HomePage() {
  return (
    <div className="min-h-screen bg-background">
      <SoftwareApplicationSchema />
      <FAQPageSchema />
      <SiteHeader />

      {/* Hero Section */}
      <section className="relative overflow-hidden">
        {/* Background effects */}
        <div className="absolute inset-0 -z-10">
          <div className="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[600px] bg-primary/[0.04] rounded-full blur-3xl" />
          <div className="absolute top-1/3 left-1/4 w-[400px] h-[400px] bg-purple-500/[0.03] rounded-full blur-3xl" />
          <div className="absolute top-1/3 right-1/4 w-[400px] h-[400px] bg-indigo-500/[0.03] rounded-full blur-3xl" />
        </div>

        <div className="container max-w-screen-2xl flex flex-col items-center justify-center gap-8 pt-24 pb-20 md:pt-40 md:pb-32 px-6">
          {/* Top badge */}
          <div className="flex items-center gap-2 rounded-full border border-border/60 bg-accent/40 px-4 py-1.5 animate-fade-in">
            <div className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse-glow" />
            <span className="tag-mono text-muted-foreground">The CRM Framework for Builders Who Ship</span>
          </div>

          {/* Main heading */}
          <h1 className="text-balance text-center text-5xl font-bold tracking-tight sm:text-6xl md:text-7xl lg:text-[5.25rem] leading-[1.05] animate-fade-in">
            <span className="gradient-text">Go + React.</span>
            <br />
            Built with Grit.
          </h1>

          {/* Subtitle */}
          <p className="text-balance text-center text-lg text-muted-foreground max-w-2xl leading-relaxed animate-fade-in">
            The batteries-included framework for CRMs, admin dashboards, SaaS products, and internal tools.
            Go performance on the backend. React on the frontend. A Filament-like admin panel.
            Full-stack code generation. Self-host everything. Ship in hours, not weeks.
          </p>

          {/* Action buttons */}
          <div className="flex flex-wrap items-center justify-center gap-3 mt-2 animate-fade-in">
            <Button size="lg" asChild className="h-11 px-6 text-sm font-medium rounded-lg gap-2 bg-primary hover:bg-primary/90 glow-purple-sm">
              <Link href="/docs/getting-started/quick-start">
                <Rocket className="h-4 w-4" />
                Get Started
              </Link>
            </Button>
            <Button size="lg" variant="outline" asChild className="h-11 px-6 text-sm font-medium rounded-lg border-border/60 hover:bg-accent/60 bg-transparent gap-2">
              <Link href="https://github.com/MUKE-coder/grit" target="_blank" rel="noreferrer">
                <Github className="h-4 w-4" />
                View on GitHub
              </Link>
            </Button>
          </div>

          {/* Platform tags */}
          <div className="flex items-center gap-4 mt-4 animate-fade-in flex-wrap justify-center">
            {['Go', 'React', 'TypeScript', 'PostgreSQL', 'Redis', 'Docker'].map((platform) => (
              <span key={platform} className="tag-mono text-muted-foreground/60 flex items-center gap-1.5">
                <span className="h-1 w-1 rounded-full bg-primary/40" />
                {platform}
              </span>
            ))}
          </div>

          {/* Terminal Demo */}
          <div className="mt-12 w-full max-w-2xl animate-fade-in">
            <AnimatedTerminal />
          </div>
        </div>
      </section>

      {/* Divider */}
      <div className="w-full h-px bg-gradient-to-r from-transparent via-border/60 to-transparent" />

      {/* Built for Builders */}
      <section className="container max-w-screen-2xl py-24 md:py-32 px-6">
        <div className="text-center mb-16">
          <span className="tag-mono text-primary/80 mb-4 block">Built for Builders</span>
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl mb-4">
            Ship fast. Self-host. Own everything.
          </h2>
          <p className="text-muted-foreground max-w-2xl mx-auto">
            Grit is for developers who want to build real products without fighting infrastructure.
            No vendor lock-in. No serverless surprises. Just a solid framework that gets out of your way.
          </p>
        </div>

        <div className="grid gap-5 md:grid-cols-2 max-w-5xl mx-auto">
          <div className="group rounded-xl border border-border/40 bg-card/50 p-6 hover:border-primary/20 transition-all duration-300">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 border border-primary/10 group-hover:bg-primary/15 transition-colors mb-4">
              <Rocket className="h-5 w-5 text-primary" />
            </div>
            <h3 className="text-lg font-semibold mb-2">Ship in Hours, Not Weeks</h3>
            <p className="text-sm text-muted-foreground/80 leading-relaxed">
              One command scaffolds your entire stack. Another generates full-stack CRUD resources &mdash;
              Go model, API handler, React Query hooks, Zod schema, and admin page &mdash; all wired
              together. Stop gluing boilerplate and start building what your users actually need.
            </p>
          </div>

          <div className="group rounded-xl border border-border/40 bg-card/50 p-6 hover:border-emerald-500/20 transition-all duration-300">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-emerald-500/10 border border-emerald-500/10 group-hover:bg-emerald-500/15 transition-colors mb-4">
              <Server className="h-5 w-5 text-emerald-400" />
            </div>
            <h3 className="text-lg font-semibold mb-2">Self-Host Everything</h3>
            <p className="text-sm text-muted-foreground/80 leading-relaxed">
              No vendor lock-in. No serverless cold starts. No per-request pricing surprises. Grit compiles
              to a single Go binary and static Next.js bundles. Deploy on a $5 VPS, your own
              servers, or any cloud. You own every byte of your infrastructure and data.
            </p>
          </div>

          <div className="group rounded-xl border border-border/40 bg-card/50 p-6 hover:border-sky-500/20 transition-all duration-300">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-sky-500/10 border border-sky-500/10 group-hover:bg-sky-500/15 transition-colors mb-4">
              <Gauge className="h-5 w-5 text-sky-400" />
            </div>
            <h3 className="text-lg font-semibold mb-2">Go Performance</h3>
            <p className="text-sm text-muted-foreground/80 leading-relaxed">
              Your API compiles to native machine code. Handle thousands of concurrent requests with
              minimal memory. No interpreter overhead, no JIT warm-up, no garbage collection pauses
              that matter. Just raw speed backed by goroutines and a compiled runtime.
            </p>
          </div>

          <div className="group rounded-xl border border-border/40 bg-card/50 p-6 hover:border-amber-500/20 transition-all duration-300">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-amber-500/10 border border-amber-500/10 group-hover:bg-amber-500/15 transition-colors mb-4">
              <Shield className="h-5 w-5 text-amber-400" />
            </div>
            <h3 className="text-lg font-semibold mb-2">Rock Solid</h3>
            <p className="text-sm text-muted-foreground/80 leading-relaxed">
              Opinionated by design. One folder structure, one auth system, one state management approach.
              Every pattern is predictable and consistent. Any developer &mdash; or AI assistant &mdash; can jump
              into any Grit project and immediately understand it. Convention over configuration.
            </p>
          </div>
        </div>
      </section>

      {/* Divider */}
      <div className="w-full h-px bg-gradient-to-r from-transparent via-border/60 to-transparent" />

      {/* Features Section */}
      <section className="container max-w-screen-2xl py-24 md:py-32 px-6">
        <div className="text-center mb-16">
          <span className="tag-mono text-primary/80 mb-4 block">Why Grit</span>
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl mb-4">
            Everything you need, nothing you don&apos;t
          </h2>
          <p className="text-muted-foreground max-w-xl mx-auto">
            A full-stack CRM framework that eliminates boilerplate and lets you focus on building your product.
          </p>
        </div>

        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {[
            {
              icon: Code2,
              title: 'Full-Stack Code Gen',
              desc: 'One command generates Go models, API handlers, React Query hooks, Zod schemas, and admin pages. The full stack, wired together.',
            },
            {
              icon: Shield,
              title: 'CRM-Grade Admin Panel',
              desc: 'Resource-based admin with data tables, form builders, multi-step wizards, RBAC, widgets, and a premium dark theme.',
            },
            {
              icon: Layers,
              title: 'End-to-End Type Safety',
              desc: 'Go struct tags generate TypeScript types and Zod schemas automatically. Change the backend, the frontend stays in sync.',
            },
            {
              icon: Zap,
              title: 'Batteries Included',
              desc: 'Auth, file storage, email, background jobs, cron, Redis caching, AI integration, and GORM Studio all ship out of the box.',
            },
          ].map((feature) => (
            <Card
              key={feature.title}
              className="group border-border/40 bg-card/50 hover:bg-card/80 hover:border-primary/20 transition-all duration-300"
            >
              <CardHeader className="space-y-3 pb-4">
                <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 border border-primary/10 group-hover:bg-primary/15 transition-colors">
                  <feature.icon className="h-4 w-4 text-primary" />
                </div>
                <CardTitle className="text-base font-semibold">{feature.title}</CardTitle>
                <CardDescription className="text-sm leading-relaxed text-muted-foreground/80">
                  {feature.desc}
                </CardDescription>
              </CardHeader>
            </Card>
          ))}
        </div>
      </section>

      {/* Divider */}
      <div className="w-full h-px bg-gradient-to-r from-transparent via-border/60 to-transparent" />

      {/* Creator Quote */}
      <section className="container max-w-screen-2xl py-24 md:py-32 px-6">
        <div className="max-w-3xl mx-auto text-center">
          <div className="text-6xl font-serif text-primary/20 leading-none mb-6">&ldquo;</div>
          <blockquote className="text-xl md:text-2xl font-medium text-foreground/90 leading-relaxed mb-8 text-balance">
            I kept rebuilding the same Go + React stack for every project &mdash;
            authentication, admin panels, file uploads, background jobs. The same
            plumbing every time. Grit packages all of that into one CLI command
            so you can focus on what makes your app unique.
          </blockquote>
          <div className="flex flex-col items-center gap-3">
            <div className="relative h-14 w-14 rounded-full overflow-hidden border-2 border-primary/20">
              <Image
                src="https://14j7oh8kso.ufs.sh/f/HLxTbDBCDLwfAUUBxSZezIN7vwylkF1PXSCqAuseUG0gx8mh"
                alt="JB Web Developer"
                fill
                className="object-cover"
              />
            </div>
            <div>
              <div className="text-sm font-semibold text-foreground/90">JB Web Developer</div>
              <div className="text-xs text-muted-foreground/60">Creator of Grit</div>
            </div>
          </div>
        </div>
      </section>

      {/* Divider */}
      <div className="w-full h-px bg-gradient-to-r from-transparent via-border/60 to-transparent" />

      {/* Grit Has Everything */}
      <section className="container max-w-screen-2xl py-24 md:py-32 px-6">
        <div className="mb-16">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl md:text-[2.75rem] mb-4">
            Grit has everything.
          </h2>
          <p className="text-muted-foreground max-w-2xl leading-relaxed">
            Grit&apos;s mission is to give you every tool you need to build production-ready
            full-stack apps with Go and React. If it can be included, we&apos;ve built it in.
          </p>
        </div>

        <div className="grid gap-5 md:grid-cols-2 max-w-5xl">
          {/* Authentication */}
          <div className="rounded-xl border border-border/40 bg-card/50 p-6">
            <div className="flex items-center gap-2.5 mb-2">
              <Lock className="h-4.5 w-4.5 text-primary/70" />
              <h3 className="text-lg font-semibold">Authentication</h3>
            </div>
            <p className="text-sm text-muted-foreground/70 leading-relaxed mb-5">
              JWT-based auth with login, register, password reset, email verification, and role-based access control.
              Three roles out of the box (Admin, Editor, User) and a CLI command to add more.
            </p>
            <div className="rounded-lg border border-border/30 bg-accent/20 p-4 space-y-2.5">
              <div className="flex items-center gap-2.5">
                <div className="h-8 w-8 rounded-md bg-primary/10 flex items-center justify-center">
                  <Lock className="h-3.5 w-3.5 text-primary/50" />
                </div>
                <div className="flex-1">
                  <div className="h-2 w-16 rounded-full bg-foreground/10" />
                  <div className="h-1.5 w-24 rounded-full bg-foreground/5 mt-1.5" />
                </div>
              </div>
              <div className="space-y-2">
                <div className="h-8 rounded-md bg-foreground/[0.03] border border-border/20 px-3 flex items-center">
                  <div className="h-1.5 w-28 rounded-full bg-foreground/8" />
                </div>
                <div className="h-8 rounded-md bg-foreground/[0.03] border border-border/20 px-3 flex items-center">
                  <div className="h-1.5 w-20 rounded-full bg-foreground/8" />
                </div>
              </div>
              <div className="h-8 rounded-md bg-primary/15 flex items-center justify-center">
                <div className="h-1.5 w-12 rounded-full bg-primary/40" />
              </div>
            </div>
          </div>

          {/* Admin Dashboard */}
          <div className="rounded-xl border border-border/40 bg-card/50 p-6">
            <div className="flex items-center gap-2.5 mb-2">
              <LayoutDashboard className="h-4.5 w-4.5 text-primary/70" />
              <h3 className="text-lg font-semibold">Admin Dashboard</h3>
            </div>
            <p className="text-sm text-muted-foreground/70 leading-relaxed mb-5">
              A Filament-like admin panel with resource definitions, data tables with sorting/filtering/search,
              form builders, multi-step wizards, dashboard widgets, and dark/light themes.
            </p>
            <div className="rounded-lg border border-border/30 bg-accent/20 p-4">
              {/* Table mockup */}
              <div className="flex items-center gap-3 mb-3">
                <div className="h-2 w-14 rounded-full bg-foreground/8" />
                <div className="h-2 w-14 rounded-full bg-foreground/8" />
                <div className="h-2 w-14 rounded-full bg-foreground/8" />
                <div className="h-2 w-14 rounded-full bg-foreground/8" />
              </div>
              <div className="space-y-2">
                {[1, 2, 3, 4].map((row) => (
                  <div key={row} className="flex items-center gap-3">
                    <div className="h-2 w-14 rounded-full bg-foreground/5" />
                    <div className="h-2 w-14 rounded-full bg-foreground/5" />
                    <div className="h-2 w-14 rounded-full bg-foreground/5" />
                    <div className="h-2 w-14 rounded-full bg-foreground/5" />
                  </div>
                ))}
              </div>
              <div className="flex items-center justify-end gap-2 mt-3">
                <div className="h-5 w-5 rounded bg-foreground/5" />
                <div className="h-5 w-5 rounded bg-foreground/5" />
              </div>
            </div>
          </div>

          {/* File Storage */}
          <div className="rounded-xl border border-border/40 bg-card/50 p-6">
            <div className="flex items-center gap-2.5 mb-2">
              <HardDrive className="h-4.5 w-4.5 text-primary/70" />
              <h3 className="text-lg font-semibold">File Storage</h3>
            </div>
            <p className="text-sm text-muted-foreground/70 leading-relaxed mb-5">
              S3-compatible file uploads that work with MinIO in development and AWS S3 or
              Cloudflare R2 in production. Built-in image processing, resizing, and a file browser in the admin panel.
            </p>
            <div className="rounded-lg border border-border/30 bg-accent/20 p-4">
              <div className="border-2 border-dashed border-border/30 rounded-lg p-6 flex flex-col items-center justify-center gap-2">
                <HardDrive className="h-6 w-6 text-muted-foreground/20" />
                <div className="h-1.5 w-28 rounded-full bg-foreground/8" />
              </div>
              <div className="flex items-center gap-2 mt-3">
                <div className="h-7 w-7 rounded bg-foreground/5 flex items-center justify-center">
                  <div className="h-3 w-3 rounded-sm bg-foreground/10" />
                </div>
                <div className="flex-1">
                  <div className="h-1.5 w-full rounded-full bg-primary/20" />
                </div>
                <X className="h-3 w-3 text-muted-foreground/20" />
              </div>
            </div>
          </div>

          {/* API Documentation */}
          <div className="rounded-xl border border-border/40 bg-card/50 p-6">
            <div className="flex items-center gap-2.5 mb-2">
              <FileText className="h-4.5 w-4.5 text-primary/70" />
              <h3 className="text-lg font-semibold">API Documentation</h3>
            </div>
            <p className="text-sm text-muted-foreground/70 leading-relaxed mb-5">
              Zero-configuration interactive API docs generated automatically from your routes.
              Browse endpoints, test requests, and export to Postman or Insomnia.
            </p>
            <div className="rounded-lg border border-border/30 bg-accent/20 p-4 space-y-2">
              {[
                { method: 'GET', path: '/api/users', color: 'bg-emerald-500/20 text-emerald-400/70' },
                { method: 'POST', path: '/api/users', color: 'bg-blue-500/20 text-blue-400/70' },
                { method: 'GET', path: '/api/users/:id', color: 'bg-emerald-500/20 text-emerald-400/70' },
                { method: 'PUT', path: '/api/users/:id', color: 'bg-amber-500/20 text-amber-400/70' },
                { method: 'DEL', path: '/api/users/:id', color: 'bg-red-500/20 text-red-400/70' },
              ].map((route) => (
                <div key={`${route.method}-${route.path}`} className="flex items-center gap-2.5 text-xs font-mono">
                  <span className={`px-1.5 py-0.5 rounded text-[10px] font-semibold ${route.color}`}>
                    {route.method}
                  </span>
                  <span className="text-foreground/30">{route.path}</span>
                </div>
              ))}
            </div>
          </div>

          {/* Database Browser */}
          <div className="rounded-xl border border-border/40 bg-card/50 p-6">
            <div className="flex items-center gap-2.5 mb-2">
              <Database className="h-4.5 w-4.5 text-primary/70" />
              <h3 className="text-lg font-semibold">Database Browser</h3>
            </div>
            <p className="text-sm text-muted-foreground/70 leading-relaxed mb-5">
              GORM Studio ships with every project &mdash; a visual database browser to inspect tables,
              view relationships, browse records, and run queries without leaving your browser.
            </p>
            <div className="rounded-lg border border-border/30 bg-accent/20 p-4">
              <div className="flex items-center gap-2 mb-3 border-b border-border/20 pb-2">
                {['users', 'posts', 'groups'].map((table) => (
                  <span key={table} className={`text-[10px] font-mono px-2 py-0.5 rounded ${table === 'users' ? 'bg-primary/10 text-primary/60' : 'text-foreground/20'}`}>
                    {table}
                  </span>
                ))}
              </div>
              <div className="space-y-1.5">
                <div className="flex items-center gap-2 text-[10px] font-mono text-foreground/15">
                  <span className="w-6 text-right text-foreground/10">1</span>
                  <span className="w-16 text-foreground/25">admin@...</span>
                  <span className="w-12">John</span>
                  <span className="px-1 rounded bg-emerald-500/10 text-emerald-400/30">ADMIN</span>
                </div>
                <div className="flex items-center gap-2 text-[10px] font-mono text-foreground/15">
                  <span className="w-6 text-right text-foreground/10">2</span>
                  <span className="w-16 text-foreground/25">jane@...</span>
                  <span className="w-12">Jane</span>
                  <span className="px-1 rounded bg-blue-500/10 text-blue-400/30">EDITOR</span>
                </div>
                <div className="flex items-center gap-2 text-[10px] font-mono text-foreground/15">
                  <span className="w-6 text-right text-foreground/10">3</span>
                  <span className="w-16 text-foreground/25">mike@...</span>
                  <span className="w-12">Mike</span>
                  <span className="px-1 rounded bg-foreground/5 text-foreground/20">USER</span>
                </div>
              </div>
            </div>
          </div>

          {/* Background Jobs & Cron */}
          <div className="rounded-xl border border-border/40 bg-card/50 p-6">
            <div className="flex items-center gap-2.5 mb-2">
              <Activity className="h-4.5 w-4.5 text-primary/70" />
              <h3 className="text-lg font-semibold">Background Jobs &amp; Cron</h3>
            </div>
            <p className="text-sm text-muted-foreground/70 leading-relaxed mb-5">
              Asynq-powered job queues backed by Redis. Send emails, process images, and run
              cleanup tasks in the background. Schedule recurring cron jobs. Monitor everything from the admin panel.
            </p>
            <div className="rounded-lg border border-border/30 bg-accent/20 p-4 space-y-2">
              {[
                { name: 'email:welcome', status: 'completed', color: 'bg-emerald-500/20 text-emerald-400/60' },
                { name: 'image:resize', status: 'running', color: 'bg-blue-500/20 text-blue-400/60' },
                { name: 'cleanup:temp', status: 'scheduled', color: 'bg-amber-500/20 text-amber-400/60' },
                { name: 'email:digest', status: 'completed', color: 'bg-emerald-500/20 text-emerald-400/60' },
              ].map((job) => (
                <div key={job.name} className="flex items-center justify-between text-xs font-mono">
                  <span className="text-foreground/30">{job.name}</span>
                  <span className={`px-1.5 py-0.5 rounded text-[10px] ${job.color}`}>
                    {job.status}
                  </span>
                </div>
              ))}
            </div>
          </div>

          {/* Email */}
          <div className="rounded-xl border border-border/40 bg-card/50 p-6">
            <div className="flex items-center gap-2.5 mb-2">
              <Mail className="h-4.5 w-4.5 text-primary/70" />
              <h3 className="text-lg font-semibold">Email</h3>
            </div>
            <p className="text-sm text-muted-foreground/70 leading-relaxed mb-5">
              Send transactional emails with Resend. Four HTML templates included: welcome, password reset,
              email verification, and notifications. Mailhog catches everything in development.
            </p>
            <div className="rounded-lg border border-border/30 bg-accent/20 p-4 space-y-2">
              <div className="flex items-center gap-2.5">
                <div className="h-6 w-6 rounded bg-primary/10 flex items-center justify-center">
                  <Mail className="h-3 w-3 text-primary/40" />
                </div>
                <div>
                  <div className="h-1.5 w-24 rounded-full bg-foreground/10" />
                  <div className="h-1 w-16 rounded-full bg-foreground/5 mt-1" />
                </div>
              </div>
              <div className="border-t border-border/20 pt-2 space-y-1.5">
                <div className="h-1.5 w-full rounded-full bg-foreground/5" />
                <div className="h-1.5 w-4/5 rounded-full bg-foreground/5" />
                <div className="h-1.5 w-3/5 rounded-full bg-foreground/5" />
              </div>
              <div className="h-7 rounded bg-primary/10 flex items-center justify-center">
                <div className="h-1.5 w-16 rounded-full bg-primary/30" />
              </div>
            </div>
          </div>

          {/* Code Generator */}
          <div className="rounded-xl border border-border/40 bg-card/50 p-6">
            <div className="flex items-center gap-2.5 mb-2">
              <Wand2 className="h-4.5 w-4.5 text-primary/70" />
              <h3 className="text-lg font-semibold">Code Generator</h3>
            </div>
            <p className="text-sm text-muted-foreground/70 leading-relaxed mb-5">
              One command generates the full stack for a resource: Go model, CRUD handler, service layer,
              Zod schema, TypeScript types, React Query hooks, and admin page &mdash; all wired together.
            </p>
            <div className="rounded-lg border border-border/30 bg-accent/20 p-4">
              <div className="font-mono text-xs space-y-1.5">
                <div className="text-foreground/30">
                  <span className="text-primary/50 select-none">$ </span>
                  grit generate resource Product
                </div>
                <div className="text-muted-foreground/25 pl-3 space-y-0.5">
                  <div><span className="text-emerald-400/40">+</span> models/product.go</div>
                  <div><span className="text-emerald-400/40">+</span> handlers/product.go</div>
                  <div><span className="text-emerald-400/40">+</span> services/product.go</div>
                  <div><span className="text-emerald-400/40">+</span> schemas/product.ts</div>
                  <div><span className="text-emerald-400/40">+</span> types/product.ts</div>
                  <div><span className="text-emerald-400/40">+</span> hooks/use-products.ts</div>
                  <div><span className="text-emerald-400/40">+</span> resources/products.ts</div>
                </div>
              </div>
            </div>
          </div>

          {/* DataTable */}
          <div className="rounded-xl border border-border/40 bg-card/50 p-6">
            <div className="flex items-center gap-2.5 mb-2">
              <Table2 className="h-4.5 w-4.5 text-primary/70" />
              <h3 className="text-lg font-semibold">DataTable</h3>
            </div>
            <p className="text-sm text-muted-foreground/70 leading-relaxed mb-5">
              A powerful data table with column sorting, full-text search, faceted filters,
              row selection, pagination, column visibility toggles, bulk actions, and custom cell renderers.
            </p>
            <div className="rounded-lg border border-border/30 bg-accent/20 p-4">
              <div className="flex items-center gap-2 mb-3">
                <div className="h-6 flex-1 rounded bg-foreground/[0.03] border border-border/20 px-2 flex items-center">
                  <div className="h-1 w-12 rounded-full bg-foreground/8" />
                </div>
                <div className="h-6 w-14 rounded bg-primary/10 flex items-center justify-center">
                  <div className="h-1 w-8 rounded-full bg-primary/30" />
                </div>
              </div>
              <div className="space-y-1.5">
                {[
                  { cols: ['John Doe', 'john@...', 'Admin'], active: true },
                  { cols: ['Jane Smith', 'jane@...', 'Editor'], active: false },
                  { cols: ['Bob Wilson', 'bob@...', 'User'], active: false },
                ].map((row) => (
                  <div key={row.cols[0]} className="flex items-center gap-2 text-[10px] font-mono">
                    <div className={`h-3 w-3 rounded border ${row.active ? 'bg-primary/20 border-primary/40' : 'border-border/30'}`} />
                    <span className="w-20 text-foreground/25 truncate">{row.cols[0]}</span>
                    <span className="w-16 text-foreground/15 truncate">{row.cols[1]}</span>
                    <span className="px-1 rounded bg-foreground/5 text-foreground/20 text-[9px]">{row.cols[2]}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Form Builder */}
          <div className="rounded-xl border border-border/40 bg-card/50 p-6">
            <div className="flex items-center gap-2.5 mb-2">
              <ClipboardList className="h-4.5 w-4.5 text-primary/70" />
              <h3 className="text-lg font-semibold">Form Builder</h3>
            </div>
            <p className="text-sm text-muted-foreground/70 leading-relaxed mb-5">
              Define forms declaratively with 10+ field types: text, textarea, number, select, date, toggle,
              checkbox, radio, richtext, and relationship fields. Multi-step wizards with per-step validation included.
            </p>
            <div className="rounded-lg border border-border/30 bg-accent/20 p-4 space-y-2.5">
              <div className="space-y-1">
                <div className="h-1.5 w-8 rounded-full bg-foreground/10" />
                <div className="h-7 rounded-md bg-foreground/[0.03] border border-border/20 px-2 flex items-center">
                  <div className="h-1.5 w-24 rounded-full bg-foreground/8" />
                </div>
              </div>
              <div className="space-y-1">
                <div className="h-1.5 w-10 rounded-full bg-foreground/10" />
                <div className="h-7 rounded-md bg-foreground/[0.03] border border-border/20 px-2 flex items-center justify-between">
                  <div className="h-1.5 w-16 rounded-full bg-foreground/8" />
                  <div className="h-2.5 w-2.5 border-l border-b border-foreground/15 rotate-[-45deg] -mt-1" />
                </div>
              </div>
              <div className="flex items-center gap-2">
                <div className="h-4 w-7 rounded-full bg-primary/20 relative">
                  <div className="absolute right-0.5 top-0.5 h-3 w-3 rounded-full bg-primary/50" />
                </div>
                <div className="h-1.5 w-14 rounded-full bg-foreground/10" />
              </div>
            </div>
          </div>

          {/* Security */}
          <div className="rounded-xl border border-border/40 bg-card/50 p-6">
            <div className="flex items-center gap-2.5 mb-2">
              <ShieldCheck className="h-4.5 w-4.5 text-primary/70" />
              <h3 className="text-lg font-semibold">Security</h3>
            </div>
            <p className="text-sm text-muted-foreground/70 leading-relaxed mb-5">
              Sentinel ships with every project &mdash; a built-in security layer with WAF,
              rate limiting, brute-force protection, anomaly detection, IP geolocation, and a real-time threat dashboard.
            </p>
            <div className="rounded-lg border border-border/30 bg-accent/20 p-4 space-y-2">
              {[
                { label: 'WAF Protection', status: 'active', color: 'bg-emerald-500/20 text-emerald-400/60' },
                { label: 'Rate Limiting', status: 'active', color: 'bg-emerald-500/20 text-emerald-400/60' },
                { label: 'Brute-Force Guard', status: 'active', color: 'bg-emerald-500/20 text-emerald-400/60' },
                { label: 'Threats Blocked', status: '1,247', color: 'bg-red-500/15 text-red-400/60' },
              ].map((item) => (
                <div key={item.label} className="flex items-center justify-between text-xs">
                  <span className="text-foreground/30 font-mono">{item.label}</span>
                  <span className={`px-1.5 py-0.5 rounded text-[10px] font-mono ${item.color}`}>
                    {item.status}
                  </span>
                </div>
              ))}
            </div>
          </div>

          {/* Performance Monitoring */}
          <div className="rounded-xl border border-border/40 bg-card/50 p-6">
            <div className="flex items-center gap-2.5 mb-2">
              <Gauge className="h-4.5 w-4.5 text-primary/70" />
              <h3 className="text-lg font-semibold">Performance Monitoring</h3>
            </div>
            <p className="text-sm text-muted-foreground/70 leading-relaxed mb-5">
              Pulse gives you self-hosted observability: request tracing, database query monitoring,
              runtime metrics, error tracking, health checks, alerting, and Prometheus export.
            </p>
            <div className="rounded-lg border border-border/30 bg-accent/20 p-4 space-y-2">
              <div className="flex items-center justify-between text-xs font-mono">
                <span className="text-foreground/30">Avg Response</span>
                <span className="text-emerald-400/60">12ms</span>
              </div>
              <div className="flex items-center justify-between text-xs font-mono">
                <span className="text-foreground/30">Requests/sec</span>
                <span className="text-blue-400/60">2,841</span>
              </div>
              <div className="flex items-center justify-between text-xs font-mono">
                <span className="text-foreground/30">Error Rate</span>
                <span className="text-emerald-400/60">0.02%</span>
              </div>
              <div className="flex items-center justify-between text-xs font-mono">
                <span className="text-foreground/30">Memory</span>
                <span className="text-foreground/25">48 MB</span>
              </div>
              {/* Mini chart */}
              <div className="flex items-end gap-0.5 h-6 pt-1">
                {[40, 55, 45, 60, 50, 70, 65, 80, 60, 75, 85, 70, 90, 75, 80].map((h, i) => (
                  <div key={i} className="flex-1 bg-primary/15 rounded-t-sm" style={{ height: `${h}%` }} />
                ))}
              </div>
            </div>
          </div>

          {/* And much more */}
          <div className="rounded-xl border border-border/40 bg-card/50 p-6 flex flex-col items-center justify-center text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 border border-primary/15 mb-4">
              <Bot className="h-6 w-6 text-primary/60" />
            </div>
            <h3 className="text-lg font-semibold mb-2">And much more</h3>
            <p className="text-sm text-muted-foreground/70 leading-relaxed max-w-sm">
              Redis caching with middleware, AI integration (Claude + OpenAI with streaming),
              multi-step form wizards, GORM Studio, Docker setup, dark/light themes,
              and shared Zod schemas across frontend and backend &mdash; all built in.
            </p>
          </div>
        </div>
      </section>

      {/* Divider */}
      <div className="w-full h-px bg-gradient-to-r from-transparent via-border/60 to-transparent" />

      {/* How It Works Section */}
      <section className="container max-w-screen-2xl py-24 md:py-32 px-6">
        <div className="text-center mb-16">
          <span className="tag-mono text-primary/80 mb-4 block">How It Works</span>
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl mb-4">
            Three steps to production
          </h2>
          <p className="text-muted-foreground max-w-xl mx-auto">
            From zero to a full-stack application with authentication, admin panel, and all the infrastructure you need.
          </p>
        </div>

        <div className="grid gap-8 md:grid-cols-3 max-w-5xl mx-auto">
          {[
            {
              step: '01',
              title: 'Scaffold',
              desc: 'Run grit new myapp and get a complete monorepo: Go API with auth, Next.js frontend, admin panel, shared types, Docker setup, and GORM Studio.',
              code: 'grit new myapp',
            },
            {
              step: '02',
              title: 'Generate',
              desc: 'Use grit generate resource to create full-stack resources. One command produces a Go model, CRUD handler, React hooks, Zod schema, and admin page.',
              code: 'grit generate resource Post --fields "title:string,content:text,published:bool"',
            },
            {
              step: '03',
              title: 'Ship',
              desc: 'Your app is production-ready from day one. Self-host on a VPS, your own servers, or any cloud. No vendor lock-in, no recurring platform fees.',
              code: 'docker compose up -d && turbo dev',
            },
          ].map((item) => (
            <div key={item.step} className="relative group">
              <div className="rounded-xl border border-border/40 bg-card/50 p-6 hover:border-primary/20 transition-all duration-300 h-full flex flex-col">
                <div className="text-4xl font-bold text-primary/15 mb-4 font-mono">{item.step}</div>
                <h3 className="text-lg font-semibold mb-3">{item.title}</h3>
                <p className="text-sm text-muted-foreground/70 leading-relaxed mb-5 flex-1">{item.desc}</p>
                <div className="rounded-lg border border-border/30 bg-accent/30 px-4 py-2.5 font-mono text-xs text-muted-foreground/60 overflow-x-auto">
                  <span className="text-primary/50 select-none">$ </span>
                  <span className="text-foreground/60">{item.code}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Divider */}
      <div className="w-full h-px bg-gradient-to-r from-transparent via-border/60 to-transparent" />

      {/* Comparison Table */}
      <section className="container max-w-screen-2xl py-24 md:py-32 px-6">
        <div className="text-center mb-16">
          <span className="tag-mono text-primary/80 mb-4 block">Comparison</span>
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl mb-4">
            How Grit stacks up
          </h2>
          <p className="text-muted-foreground max-w-2xl mx-auto">
            Grit combines the developer experience of Laravel, the type safety of the T3 stack,
            and the performance of Go into one cohesive framework.
          </p>
        </div>

        <div className="max-w-5xl mx-auto rounded-xl border border-border/40 overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-sm border-collapse">
              <thead>
                <tr>
                  <th className="text-left px-5 py-4 text-xs font-semibold text-muted-foreground/70 min-w-[180px] bg-accent/20 border-b border-border/40" />
                  <th className="text-center px-5 py-4 min-w-[90px] bg-primary/[0.08] border-b-2 border-primary/40">
                    <div className="flex flex-col items-center gap-1">
                      <span className="inline-flex items-center rounded-full bg-primary/15 border border-primary/25 px-2.5 py-0.5 text-[10px] font-bold text-primary tracking-wider uppercase">Grit</span>
                    </div>
                  </th>
                  <th className="text-center px-4 py-4 text-xs font-semibold text-muted-foreground/60 min-w-[80px] bg-accent/20 border-b border-border/40">Laravel</th>
                  <th className="text-center px-4 py-4 text-xs font-semibold text-muted-foreground/60 min-w-[80px] bg-accent/20 border-b border-border/40">Django</th>
                  <th className="text-center px-4 py-4 text-xs font-semibold text-muted-foreground/60 min-w-[80px] bg-accent/20 border-b border-border/40">Rails</th>
                  <th className="text-center px-4 py-4 text-xs font-semibold text-muted-foreground/60 min-w-[80px] bg-accent/20 border-b border-border/40">T3 Stack</th>
                  <th className="text-center px-4 py-4 text-xs font-semibold text-muted-foreground/60 min-w-[80px] bg-accent/20 border-b border-border/40">Refine</th>
                </tr>
              </thead>
              <tbody>
                {([
                  { feature: 'Language', grit: 'Go + TS', laravel: 'PHP', django: 'Python', rails: 'Ruby', t3: 'TypeScript', refine: 'TypeScript', type: 'text' as const },
                  { feature: 'Compiled / Native', grit: true, laravel: false, django: false, rails: false, t3: false, refine: false },
                  { feature: 'Built-in Admin Panel', grit: true, laravel: 'partial', django: true, rails: false, t3: false, refine: true },
                  { feature: 'Full-Stack Code Gen', grit: true, laravel: 'partial', django: false, rails: 'partial', t3: false, refine: false },
                  { feature: 'End-to-End Type Safety', grit: true, laravel: false, django: false, rails: false, t3: true, refine: true },
                  { feature: 'React Frontend', grit: true, laravel: false, django: false, rails: false, t3: true, refine: true },
                  { feature: 'Self-Hostable', grit: true, laravel: true, django: true, rails: true, t3: 'partial', refine: true },
                  { feature: 'Background Jobs', grit: true, laravel: true, django: 'partial', rails: true, t3: false, refine: false },
                  { feature: 'File Storage (S3)', grit: true, laravel: true, django: 'partial', rails: true, t3: false, refine: false },
                  { feature: 'Email Service', grit: true, laravel: true, django: true, rails: true, t3: false, refine: false },
                  { feature: 'AI Integration', grit: true, laravel: false, django: false, rails: false, t3: false, refine: false },
                  { feature: 'Database Browser', grit: true, laravel: false, django: true, rails: false, t3: false, refine: false },
                  { feature: 'Multi-Step Forms', grit: true, laravel: false, django: false, rails: false, t3: false, refine: 'partial' },
                  { feature: 'Docker Setup', grit: true, laravel: 'partial', django: false, rails: false, t3: false, refine: false },
                ] as const).map((row, idx) => (
                  <tr key={row.feature} className={idx % 2 === 0 ? 'bg-card/30' : 'bg-card/50'}>
                    <td className="px-5 py-3 text-xs text-foreground/80 font-medium border-r border-border/20">{row.feature}</td>
                    <td className="text-center px-5 py-3 bg-primary/[0.04] border-x border-primary/10">
                      {row.type === 'text' ? (
                        <span className="text-xs text-primary font-semibold">{row.grit as string}</span>
                      ) : row.grit === true ? (
                        <CheckCircle2 className="h-4.5 w-4.5 mx-auto text-primary" />
                      ) : row.grit === false ? (
                        <X className="h-4 w-4 mx-auto text-muted-foreground/25" />
                      ) : (
                        <Minus className="h-4 w-4 mx-auto text-amber-500/50" />
                      )}
                    </td>
                    {(['laravel', 'django', 'rails', 't3', 'refine'] as const).map((fw) => {
                      const val = row[fw]
                      return (
                        <td key={fw} className="text-center px-4 py-3">
                          {row.type === 'text' ? (
                            <span className="text-xs text-muted-foreground/50">{val as string}</span>
                          ) : val === true ? (
                            <CheckCircle2 className="h-4 w-4 mx-auto text-emerald-500/60" />
                          ) : val === false ? (
                            <X className="h-4 w-4 mx-auto text-muted-foreground/20" />
                          ) : (
                            <Minus className="h-4 w-4 mx-auto text-amber-500/40" />
                          )}
                        </td>
                      )
                    })}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="px-5 py-3.5 border-t border-border/30 bg-accent/10 flex flex-wrap items-center gap-x-5 gap-y-1">
            <span className="text-[11px] text-muted-foreground/40 flex items-center gap-1">
              <CheckCircle2 className="h-3 w-3 text-emerald-500/70" /> Built-in
            </span>
            <span className="text-[11px] text-muted-foreground/40 flex items-center gap-1">
              <Minus className="h-3 w-3 text-amber-500/50" /> Partial / Plugin
            </span>
            <span className="text-[11px] text-muted-foreground/40 flex items-center gap-1">
              <X className="h-3 w-3 text-muted-foreground/25" /> Not included
            </span>
            <span className="text-[11px] text-muted-foreground/30 ml-auto hidden sm:inline">
              Laravel includes Filament &middot; T3 = Next.js + tRPC + Prisma + NextAuth
            </span>
          </div>
        </div>
      </section>

      {/* Divider */}
      <div className="w-full h-px bg-gradient-to-r from-transparent via-border/60 to-transparent" />

      {/* Tech Stack Section */}
      <section className="container max-w-screen-2xl py-24 md:py-32 px-6">
        <div className="text-center mb-16">
          <span className="tag-mono text-primary/80 mb-4 block">Tech Stack</span>
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl mb-4">
            Opinionated. On purpose.
          </h2>
          <p className="text-muted-foreground max-w-xl mx-auto">
            Every tool is chosen for a reason. One way to do things means any developer can jump into any Grit project and instantly understand it.
          </p>
        </div>

        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4 max-w-5xl mx-auto">
          {[
            { label: 'Backend', tech: 'Go + Gin + GORM' },
            { label: 'Frontend', tech: 'Next.js 14+ (App Router)' },
            { label: 'Admin', tech: 'Next.js + shadcn/ui' },
            { label: 'Styling', tech: 'Tailwind CSS + shadcn/ui' },
            { label: 'Validation', tech: 'Zod' },
            { label: 'Data Fetching', tech: 'React Query' },
            { label: 'Database', tech: 'PostgreSQL + GORM' },
            { label: 'Cache & Queue', tech: 'Redis + Asynq' },
            { label: 'File Storage', tech: 'S3 / R2 / MinIO' },
            { label: 'Email', tech: 'Resend' },
            { label: 'Monorepo', tech: 'Turborepo + pnpm' },
            { label: 'DB Browser', tech: 'GORM Studio' },
          ].map((item) => (
            <div
              key={item.label}
              className="flex items-center gap-3 rounded-lg border border-border/30 bg-card/30 px-4 py-3 hover:border-primary/20 transition-colors"
            >
              <span className="text-xs font-mono text-primary/60 w-20 shrink-0">{item.label}</span>
              <span className="text-sm text-foreground/70">{item.tech}</span>
            </div>
          ))}
        </div>
      </section>

      {/* Divider */}
      <div className="w-full h-px bg-gradient-to-r from-transparent via-border/60 to-transparent" />

      {/* CTA Section */}
      <section className="container max-w-screen-2xl py-24 md:py-32 px-6">
        <div className="relative rounded-2xl border border-border/40 bg-card/30 overflow-hidden">
          {/* Background glow */}
          <div className="absolute inset-0 -z-10">
            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[300px] bg-primary/[0.04] rounded-full blur-3xl" />
          </div>

          <div className="text-center space-y-6 py-20 md:py-28 px-6">
            <span className="tag-mono text-primary/80">Ready to ship?</span>
            <h2 className="text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight">
              Start building with{' '}
              <span className="gradient-text">Grit</span>
            </h2>
            <p className="text-muted-foreground max-w-lg mx-auto">
              Go performance. React ecosystem. Self-hosted. One framework.
              Five minutes to a full-stack app with auth, admin panel, and all the batteries you need.
            </p>
            <div className="flex flex-wrap items-center justify-center gap-3 pt-4">
              <Button size="lg" asChild className="h-11 px-6 text-sm font-medium rounded-lg glow-purple-sm">
                <Link href="/docs/getting-started/quick-start">
                  Get Started
                  <ArrowRight className="ml-1.5 h-4 w-4" />
                </Link>
              </Button>
              <Button size="lg" variant="outline" asChild className="h-11 px-6 text-sm font-medium rounded-lg border-border/60 hover:bg-accent/60 bg-transparent">
                <Link href="/docs">
                  Read the Docs
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-border/30 py-8 px-6">
        <div className="container max-w-screen-2xl flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2">
            <div className="flex h-5 w-5 items-center justify-center rounded bg-primary/15 border border-primary/20">
              <span className="text-primary font-mono font-bold text-[9px]">G</span>
            </div>
            <span className="text-xs text-muted-foreground/50">Grit Framework</span>
          </div>
          <div className="flex items-center gap-6">
            {[
              { label: 'Docs', href: '/docs' },
              { label: 'Showcase', href: '/showcase' },
              { label: 'Tutorials', href: '/docs/tutorials/blog' },
              { label: 'Course', href: '/course' },
              { label: 'Playground', href: '/playground' },
              { label: 'GitHub', href: 'https://github.com/MUKE-coder/grit' },
            ].map((link) => (
              <Link
                key={link.label}
                href={link.href}
                className="text-xs text-muted-foreground/40 hover:text-muted-foreground transition-colors"
              >
                {link.label}
              </Link>
            ))}
          </div>
          <div className="flex items-center gap-3">
            <Link href="https://www.youtube.com/@JBWEBDEVELOPER" target="_blank" rel="noreferrer" className="text-muted-foreground/40 hover:text-muted-foreground transition-colors">
              <svg className="h-4 w-4" viewBox="0 0 24 24" fill="currentColor"><path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/></svg>
              <span className="sr-only">YouTube</span>
            </Link>
            <Link href="https://www.linkedin.com/in/muke-johnbaptist-95bb82198/" target="_blank" rel="noreferrer" className="text-muted-foreground/40 hover:text-muted-foreground transition-colors">
              <svg className="h-4 w-4" viewBox="0 0 24 24" fill="currentColor"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>
              <span className="sr-only">LinkedIn</span>
            </Link>
            <Link href="https://jb.desishub.com" target="_blank" rel="noreferrer" className="text-muted-foreground/40 hover:text-muted-foreground transition-colors">
              <Globe className="h-4 w-4" />
              <span className="sr-only">Website</span>
            </Link>
            <Link href="https://github.com/MUKE-coder/grit" target="_blank" rel="noreferrer" className="text-muted-foreground/40 hover:text-muted-foreground transition-colors">
              <Github className="h-4 w-4" />
              <span className="sr-only">GitHub</span>
            </Link>
          </div>
        </div>
      </footer>
    </div>
  )
}
