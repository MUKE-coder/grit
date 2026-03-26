import Link from 'next/link'
import { ArrowRight, Rocket, Terminal, Layers, Zap, Shield, Code2, Database, Server, Gauge, CheckCircle2, X, Minus, Download, BookOpen } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs')

export default function DocsIntroductionPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Introduction</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                The CRM Framework for Builders Who Ship
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit is a full-stack meta-framework purpose-built for CRMs, admin dashboards, SaaS
                products, and internal tools. It fuses a{' '}
                <strong className="text-foreground">Go backend</strong> (Gin + GORM) with your choice of{' '}
                <strong className="text-foreground">Next.js or TanStack Router (Vite)</strong> frontend and a{' '}
                <strong className="text-foreground">Filament-like admin panel</strong> &mdash; with 5 architecture
                modes, an interactive CLI, authentication, file storage, email, background jobs, AI
                integration, and a visual database browser all wired together out of the box.
              </p>
            </div>

            {/* Value proposition pills */}
            <div className="flex flex-wrap gap-2 mb-10">
              {[
                { label: 'Ship in hours, not weeks', color: 'bg-primary/10 text-primary border-primary/20' },
                { label: 'Self-host everything', color: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' },
                { label: 'Go performance', color: 'bg-sky-500/10 text-sky-400 border-sky-500/20' },
                { label: 'CRM-ready', color: 'bg-amber-500/10 text-amber-400 border-amber-500/20' },
              ].map((pill) => (
                <span
                  key={pill.label}
                  className={`inline-flex items-center rounded-full border px-3 py-1 text-xs font-medium ${pill.color}`}
                >
                  {pill.label}
                </span>
              ))}
            </div>

            {/* YouTube intro video */}
            <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-10 shadow-lg">
              <div className="px-4 py-2.5 border-b border-border/30 bg-accent/20 flex items-center gap-2">
                <div className="flex items-center gap-1.5">
                  <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                  <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                  <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                </div>
                <span className="ml-2 text-[11px] font-mono text-muted-foreground/50">What is Grit?</span>
              </div>
              <div className="aspect-video">
                <iframe
                  src="https://www.youtube.com/embed/F-6ZEJLy8Yc"
                  title="What is Grit Framework?"
                  allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
                  allowFullScreen
                  className="w-full h-full"
                />
              </div>
            </div>

            {/* Handbook download banner */}
            <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4 p-5 rounded-xl border border-primary/20 bg-primary/5 mb-10">
              <div className="flex items-center gap-3 flex-1">
                <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 border border-primary/20">
                  <BookOpen className="h-4 w-4 text-primary/70" />
                </div>
                <div>
                  <p className="text-sm font-semibold text-foreground/80">Grit Developer Handbook</p>
                  <p className="text-xs text-muted-foreground/60">The complete Grit reference — printable PDF for offline reading</p>
                </div>
              </div>
              <a
                href="https://14j7oh8kso.ufs.sh/f/HLxTbDBCDLwfeHHJl34ZKSqNhOvVj6p9rg3Icmo05TAEwQ4a"
                target="_blank"
                rel="noopener noreferrer"
              >
                <Button size="sm" className="shrink-0 gap-1.5">
                  <Download className="h-3.5 w-3.5" />
                  Download PDF
                </Button>
              </a>
            </div>

            {/* Built for Builders */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-6">
                Built for Builders
              </h2>
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="rounded-xl border border-border/40 bg-card/50 p-5 hover:border-primary/20 transition-colors">
                  <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 border border-primary/10 mb-3">
                    <Rocket className="h-4.5 w-4.5 text-primary" />
                  </div>
                  <h3 className="text-sm font-semibold mb-1.5">Ship Fast</h3>
                  <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                    One command scaffolds your entire stack. Another generates full-stack CRUD resources.
                    Stop gluing boilerplate together and start building what matters. Go from zero to a
                    production-ready app with auth, admin panel, and API in minutes.
                  </p>
                </div>

                <div className="rounded-xl border border-border/40 bg-card/50 p-5 hover:border-primary/20 transition-colors">
                  <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-emerald-500/10 border border-emerald-500/10 mb-3">
                    <Server className="h-4.5 w-4.5 text-emerald-400" />
                  </div>
                  <h3 className="text-sm font-semibold mb-1.5">Self-Host Everything</h3>
                  <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                    No vendor lock-in. No serverless cold starts. No per-request pricing. Grit compiles
                    to a single Go binary and static frontend bundles (Next.js or Vite). Deploy on a $5 VPS, your own
                    servers, or any cloud. You own every byte of your infrastructure.
                  </p>
                </div>

                <div className="rounded-xl border border-border/40 bg-card/50 p-5 hover:border-primary/20 transition-colors">
                  <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-sky-500/10 border border-sky-500/10 mb-3">
                    <Gauge className="h-4.5 w-4.5 text-sky-400" />
                  </div>
                  <h3 className="text-sm font-semibold mb-1.5">Go Performance</h3>
                  <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                    Go compiles to native machine code. Your API handles thousands of concurrent requests
                    with minimal memory. No garbage collection pauses that matter, no interpreter overhead,
                    no JIT warm-up. Just raw speed backed by goroutines and a compiled runtime.
                  </p>
                </div>

                <div className="rounded-xl border border-border/40 bg-card/50 p-5 hover:border-primary/20 transition-colors">
                  <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-amber-500/10 border border-amber-500/10 mb-3">
                    <Shield className="h-4.5 w-4.5 text-amber-400" />
                  </div>
                  <h3 className="text-sm font-semibold mb-1.5">Rock Solid</h3>
                  <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                    Opinionated by design. One folder structure, one auth system, one state management
                    approach. Every pattern is predictable. Any developer (or AI assistant) can jump
                    into any Grit project and immediately understand it. Convention over configuration.
                  </p>
                </div>
              </div>
            </div>

            {/* What is Grit */}
            <div className="prose-grit mb-12">
              <h2>Why Grit?</h2>
              <p>
                Think of Grit as <strong>Laravel&apos;s developer experience, but with Go&apos;s performance
                and React&apos;s frontend ecosystem.</strong> Instead of stitching together 15+ tools to build
                a modern full-stack app, you run one command and get everything: a Go API server, a React
                frontend (Next.js or TanStack Router), an admin dashboard, shared TypeScript types, Docker infrastructure, and a CLI
                that generates full-stack resources for you.
              </p>
              <p>
                Grit is designed for the kind of software that runs businesses: <strong>CRMs, admin panels,
                internal tools, SaaS dashboards, and data-heavy applications</strong>. The admin panel is
                not an afterthought &mdash; it&apos;s a first-class citizen with resource definitions,
                advanced data tables, form builders, multi-step wizards, charts, and widgets. Define a
                resource once, and Grit generates the entire stack from Go model to admin UI.
              </p>
              <p>
                Every piece of code Grit generates is yours. No lock-in, no black boxes,
                no runtime magic. The CLI produces clean, readable, editable files that you own and control.
                Self-host on your own servers with zero ongoing costs beyond hosting.
              </p>
            </div>

            {/* Comparison table */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-2">
                How Grit Compares
              </h2>
              <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-6">
                Grit combines the best parts of existing frameworks into one cohesive tool.
                Here&apos;s how it stacks up:
              </p>

              <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/40 bg-accent/20">
                        <th className="text-left px-4 py-3 text-xs font-semibold text-muted-foreground/70 w-[180px]">Feature</th>
                        <th className="text-center px-3 py-3 text-xs font-semibold text-primary">Grit</th>
                        <th className="text-center px-3 py-3 text-xs font-semibold text-muted-foreground/70">Laravel</th>
                        <th className="text-center px-3 py-3 text-xs font-semibold text-muted-foreground/70">Django</th>
                        <th className="text-center px-3 py-3 text-xs font-semibold text-muted-foreground/70">Rails</th>
                        <th className="text-center px-3 py-3 text-xs font-semibold text-muted-foreground/70">T3 Stack</th>
                        <th className="text-center px-3 py-3 text-xs font-semibold text-muted-foreground/70">Refine</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/30">
                      {[
                        {
                          feature: 'Language',
                          grit: 'Go + TypeScript',
                          laravel: 'PHP',
                          django: 'Python',
                          rails: 'Ruby',
                          t3: 'TypeScript',
                          refine: 'TypeScript',
                          type: 'text',
                        },
                        {
                          feature: 'Compiled / Native',
                          grit: true,
                          laravel: false,
                          django: false,
                          rails: false,
                          t3: false,
                          refine: false,
                        },
                        {
                          feature: 'Built-in Admin Panel',
                          grit: true,
                          laravel: 'partial',
                          django: true,
                          rails: false,
                          t3: false,
                          refine: true,
                        },
                        {
                          feature: 'Full-Stack Code Gen',
                          grit: true,
                          laravel: 'partial',
                          django: false,
                          rails: 'partial',
                          t3: false,
                          refine: false,
                        },
                        {
                          feature: 'End-to-End Type Safety',
                          grit: true,
                          laravel: false,
                          django: false,
                          rails: false,
                          t3: true,
                          refine: true,
                        },
                        {
                          feature: 'React Frontend',
                          grit: true,
                          laravel: false,
                          django: false,
                          rails: false,
                          t3: true,
                          refine: true,
                        },
                        {
                          feature: 'Self-Hostable',
                          grit: true,
                          laravel: true,
                          django: true,
                          rails: true,
                          t3: 'partial',
                          refine: true,
                        },
                        {
                          feature: 'Background Jobs',
                          grit: true,
                          laravel: true,
                          django: 'partial',
                          rails: true,
                          t3: false,
                          refine: false,
                        },
                        {
                          feature: 'File Storage (S3)',
                          grit: true,
                          laravel: true,
                          django: 'partial',
                          rails: true,
                          t3: false,
                          refine: false,
                        },
                        {
                          feature: 'Email Service',
                          grit: true,
                          laravel: true,
                          django: true,
                          rails: true,
                          t3: false,
                          refine: false,
                        },
                        {
                          feature: 'AI Integration',
                          grit: true,
                          laravel: false,
                          django: false,
                          rails: false,
                          t3: false,
                          refine: false,
                        },
                        {
                          feature: 'Database Browser',
                          grit: true,
                          laravel: false,
                          django: true,
                          rails: false,
                          t3: false,
                          refine: false,
                        },
                        {
                          feature: 'Multi-Step Forms',
                          grit: true,
                          laravel: false,
                          django: false,
                          rails: false,
                          t3: false,
                          refine: 'partial',
                        },
                        {
                          feature: 'Docker Setup',
                          grit: true,
                          laravel: 'partial',
                          django: false,
                          rails: false,
                          t3: false,
                          refine: false,
                        },
                      ].map((row) => (
                        <tr key={row.feature} className="hover:bg-accent/10 transition-colors">
                          <td className="px-4 py-2.5 text-xs text-foreground/80 font-medium">
                            {row.feature}
                          </td>
                          {['grit', 'laravel', 'django', 'rails', 't3', 'refine'].map((fw) => {
                            const val = row[fw as keyof typeof row]
                            return (
                              <td key={fw} className="text-center px-3 py-2.5">
                                {row.type === 'text' ? (
                                  <span className={`text-xs ${fw === 'grit' ? 'text-primary font-medium' : 'text-muted-foreground/60'}`}>
                                    {val as string}
                                  </span>
                                ) : val === true ? (
                                  <CheckCircle2 className={`h-4 w-4 mx-auto ${fw === 'grit' ? 'text-primary' : 'text-emerald-500/70'}`} />
                                ) : val === false ? (
                                  <X className="h-4 w-4 mx-auto text-muted-foreground/25" />
                                ) : (
                                  <Minus className="h-4 w-4 mx-auto text-amber-500/50" />
                                )}
                              </td>
                            )
                          })}
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
                <div className="px-4 py-2.5 border-t border-border/30 bg-accent/10">
                  <p className="text-[11px] text-muted-foreground/40">
                    <CheckCircle2 className="h-3 w-3 inline-block mr-1 text-emerald-500/70 -mt-px" /> Built-in
                    <span className="mx-2">&middot;</span>
                    <Minus className="h-3 w-3 inline-block mr-1 text-amber-500/50 -mt-px" /> Partial / via plugin
                    <span className="mx-2">&middot;</span>
                    <X className="h-3 w-3 inline-block mr-1 text-muted-foreground/25 -mt-px" /> Not included
                    <span className="mx-2">&middot;</span>
                    Laravel comparison includes Filament. T3 Stack = Next.js + tRPC + Prisma + NextAuth.
                  </p>
                </div>
              </div>
            </div>

            {/* Key features cards */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-6">
                Key Features
              </h2>
              <div className="grid gap-4 sm:grid-cols-2">
                {[
                  {
                    icon: Terminal,
                    title: 'Interactive CLI Scaffolder',
                    desc: 'grit new myapp launches an interactive CLI to pick your architecture (single/double/triple) and frontend (Next.js or TanStack Router). Five modes from API-only to full monorepo.',
                  },
                  {
                    icon: Code2,
                    title: 'Full-Stack Code Gen',
                    desc: 'grit generate resource Post creates a Go model, CRUD handler, React Query hooks, Zod schema, and admin page in one shot.',
                  },
                  {
                    icon: Shield,
                    title: 'CRM-Grade Admin Panel',
                    desc: 'Resource-based admin with data tables, form builders, multi-step wizards, charts, widgets, RBAC, and a polished dark theme.',
                  },
                  {
                    icon: Layers,
                    title: 'End-to-End Type Safety',
                    desc: 'Go struct tags auto-generate TypeScript types and Zod schemas. Change the backend, run grit sync, and the frontend stays in sync.',
                  },
                  {
                    icon: Zap,
                    title: 'Batteries Included',
                    desc: 'JWT auth, file storage (S3/R2/MinIO), email (Resend), background jobs, cron, Redis caching, and AI integration ship with every project.',
                  },
                  {
                    icon: Database,
                    title: 'GORM Studio',
                    desc: 'A visual database browser embedded at /studio. Browse tables, inspect data, and manage your database without leaving the browser.',
                  },
                ].map((feature) => (
                  <Card key={feature.title} className="border-border/40 bg-card/50 hover:border-primary/20 transition-colors">
                    <CardHeader className="space-y-2 pb-4">
                      <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 border border-primary/10">
                        <feature.icon className="h-4 w-4 text-primary" />
                      </div>
                      <CardTitle className="text-sm font-semibold">{feature.title}</CardTitle>
                      <CardDescription className="text-[13px] leading-relaxed">{feature.desc}</CardDescription>
                    </CardHeader>
                  </Card>
                ))}
              </div>
            </div>

            {/* How it works */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-6">
                How It Works
              </h2>
              <div className="space-y-6">
                {[
                  {
                    step: '1',
                    title: 'Scaffold your project',
                    description: 'Run grit new myapp to launch the interactive CLI. Choose your architecture mode and frontend framework (Next.js or TanStack Router), and Grit generates the full project -- Go API with auth, frontend app, admin panel, shared types, Docker Compose, and GORM Studio.',
                  },
                  {
                    step: '2',
                    title: 'Start the dev environment',
                    description: 'Run docker compose up -d to start PostgreSQL, Redis, MinIO, and Mailhog. Then run turbo dev to start the Go API and Next.js apps with hot reload.',
                  },
                  {
                    step: '3',
                    title: 'Generate resources',
                    description: 'Use grit generate resource Post --fields "title:string,content:text" to create a full-stack resource. The Go model, handler, React hooks, Zod schema, and admin page are all created and wired together.',
                  },
                  {
                    step: '4',
                    title: 'Build and ship',
                    description: 'Your app is production-ready from day one. Multi-stage Docker builds for Go and Next.js. Self-host on a VPS, your own servers, or any cloud provider. No vendor lock-in, no recurring platform fees.',
                  },
                ].map((item) => (
                  <div key={item.step} className="relative pl-10">
                    <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      {item.step}
                    </div>
                    <div>
                      <h3 className="text-[15px] font-semibold mb-1">{item.title}</h3>
                      <p className="text-[13px] text-muted-foreground/70 leading-relaxed">{item.description}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Quick example */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Quick Example
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Create a project and generate your first resource in under a minute:
              </p>

              {/* Terminal block: scaffold */}
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm mb-4">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm space-y-2">
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">grit new myapp</span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">cd myapp</span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">docker compose up -d</span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">grit generate resource Post --fields &quot;title:string,content:text,published:bool&quot;</span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">turbo dev</span>
                  </div>
                </div>
              </div>
              <p className="text-sm text-muted-foreground/60">
                That&apos;s it. You now have a Go API serving CRUD endpoints for Posts, a React frontend
                with data fetching hooks, Zod validation, and an admin panel with a data table and
                form -- all connected and ready to use.
              </p>
            </div>

            {/* What gets generated */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                What Gets Generated
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                When you run <code className="text-sm font-mono bg-accent/80 px-1.5 py-0.5 rounded text-primary">grit generate resource Post</code>, the CLI creates:
              </p>
              <div className="space-y-3">
                {[
                  { file: 'apps/api/internal/models/post.go', desc: 'GORM model with struct tags, timestamps, and soft deletes' },
                  { file: 'apps/api/internal/handlers/post.go', desc: 'CRUD handler with pagination, sorting, filtering, and search' },
                  { file: 'apps/api/internal/services/post.go', desc: 'Business logic layer with reusable query scopes' },
                  { file: 'packages/shared/schemas/post.ts', desc: 'Zod validation schemas (create + update)' },
                  { file: 'packages/shared/types/post.ts', desc: 'TypeScript types inferred from the Go model' },
                  { file: 'apps/admin/hooks/use-posts.ts', desc: 'React Query hooks for all CRUD operations' },
                  { file: 'apps/admin/resources/posts.ts', desc: 'Resource definition with table, form, and multi-step wizard config' },
                  { file: 'apps/admin/app/resources/posts/page.tsx', desc: 'Admin page with data table, search, filters, and form actions' },
                ].map((item) => (
                  <div key={item.file} className="flex items-start gap-3 rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                    <code className="text-xs font-mono text-primary/70 shrink-0 mt-0.5">{item.file}</code>
                    <span className="text-xs text-muted-foreground/60">{item.desc}</span>
                  </div>
                ))}
              </div>
            </div>

            {/* Who is Grit for */}
            <div className="prose-grit mb-12">
              <h2>Who is Grit For?</h2>
              <ul>
                <li><strong>Freelancers and agencies</strong> who need to ship client projects fast with a complete stack that works out of the box.</li>
                <li><strong>Solo SaaS developers</strong> who want Go&apos;s performance for their backend but React&apos;s ecosystem for their frontend.</li>
                <li><strong>Small to mid-size teams</strong> building internal tools, CRMs, dashboards, or SaaS products.</li>
                <li><strong>Laravel developers</strong> who want to move to Go but miss the developer experience and Filament&apos;s admin panel.</li>
                <li><strong>Next.js developers</strong> who are frustrated with serverless limitations and want a real, self-hosted backend.</li>
                <li><strong>Anyone tired of vendor lock-in</strong> who wants to own their infrastructure and deploy anywhere.</li>
              </ul>
            </div>

            {/* Next steps */}
            <div className="mb-8">
              <h2 className="text-2xl font-semibold tracking-tight mb-6">
                Next Steps
              </h2>
              <div className="grid gap-3 sm:grid-cols-2">
                {[
                  { title: 'Quick Start', desc: 'Get a project running in 5 minutes', href: '/docs/getting-started/quick-start' },
                  { title: 'Philosophy', desc: 'Why Grit exists and the design principles behind it', href: '/docs/getting-started/philosophy' },
                  { title: 'Installation', desc: 'System requirements and CLI setup', href: '/docs/getting-started/installation' },
                  { title: 'Tutorials', desc: 'Step-by-step guides to build real apps', href: '/docs/tutorials/learn' },
                ].map((item) => (
                  <Link key={item.href} href={item.href}>
                    <div className="group rounded-lg border border-border/40 bg-card/50 p-4 hover:border-primary/20 hover:bg-card/80 transition-all duration-200">
                      <h3 className="text-sm font-semibold mb-1 group-hover:text-primary transition-colors flex items-center gap-1.5">
                        {item.title}
                        <ArrowRight className="h-3 w-3 opacity-0 -translate-x-1 group-hover:opacity-100 group-hover:translate-x-0 transition-all" />
                      </h3>
                      <p className="text-xs text-muted-foreground/60">{item.desc}</p>
                    </div>
                  </Link>
                ))}
              </div>
            </div>

            {/* Bottom nav */}
            <div className="flex flex-wrap gap-3 pt-4 border-t border-border/30">
              <Button asChild className="glow-purple-sm">
                <Link href="/docs/getting-started/quick-start">
                  Quick Start
                  <ArrowRight className="ml-1.5 h-4 w-4" />
                </Link>
              </Button>
              <Button variant="outline" asChild className="border-border/60 bg-transparent hover:bg-accent/50">
                <Link href="/docs/getting-started/philosophy">
                  Read the Philosophy
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
