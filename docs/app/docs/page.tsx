import Link from 'next/link'
import { ArrowRight, Rocket, Box, Terminal, Layers, Zap, Shield, Code2, Database } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'

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
                Welcome to Grit
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit is a full-stack meta-framework that fuses a{' '}
                <strong className="text-foreground">Go backend</strong> (Gin + GORM) with a{' '}
                <strong className="text-foreground">Next.js frontend</strong> (React + App Router) and a{' '}
                <strong className="text-foreground">Filament-like admin panel</strong> in a single monorepo.
                It ships with authentication, file storage, email, background jobs, AI integration,
                and a visual database browser -- all wired together out of the box.
              </p>
            </div>

            {/* What is Grit */}
            <div className="prose-grit mb-12">
              <h2>What is Grit?</h2>
              <p>
                Think of Grit as <strong>Laravel&apos;s developer experience, but with Go&apos;s performance
                and React&apos;s frontend ecosystem.</strong> Instead of stitching together 15+ tools to build
                a modern full-stack app, you run one command and get everything: a Go API server, a React
                frontend, an admin dashboard, shared TypeScript types, Docker infrastructure, and a CLI
                that generates full-stack resources for you.
              </p>
              <p>
                Grit is opinionated by design. There is one way to structure your project, one auth system,
                one state management approach, one folder structure. This means any developer can jump into
                any Grit project and immediately understand it. It also means AI assistants can generate
                Grit code confidently because the patterns are predictable and consistent.
              </p>
              <p>
                Every piece of code Grit generates is yours. No lock-in, no black boxes,
                no runtime magic. The CLI produces clean, readable, editable files that you own and control.
              </p>
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
                    title: 'CLI Scaffolder',
                    desc: 'grit new myapp creates a complete monorepo with Go API, Next.js frontend, admin panel, shared types, and Docker setup.',
                  },
                  {
                    icon: Code2,
                    title: 'Full-Stack Code Gen',
                    desc: 'grit generate resource Post creates a Go model, CRUD handler, React Query hooks, Zod schema, and admin page in one shot.',
                  },
                  {
                    icon: Shield,
                    title: 'Admin Panel',
                    desc: 'A Filament-like resource-based admin dashboard with data tables, form builders, charts, widgets, and a polished dark theme.',
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
                    description: 'Run grit new myapp to generate the full monorepo -- Go API with auth, Next.js web app, admin panel, shared types, Docker Compose, and GORM Studio.',
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
                    description: 'Your app is production-ready from day one. Multi-stage Docker builds for Go and Next.js. Deploy anywhere -- VPS, Kubernetes, or your favorite cloud provider.',
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
                  { file: 'apps/admin/app/resources/posts/page.tsx', desc: 'Admin page with data table, search, and delete actions' },
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
                <li><strong>Laravel developers</strong> who want to move to Go but miss the developer experience.</li>
                <li><strong>Next.js developers</strong> who are frustrated with serverless limitations and want a real backend.</li>
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
                  { title: 'Project Structure', desc: 'Understand the monorepo layout', href: '/docs/getting-started/project-structure' },
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
