import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { FileTree, type TreeNode } from '@/components/diagram'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/architecture-modes')

interface ArchMode {
  name: string
  flag: string
  tagline: string
  color: string
  ideal: string
  example: string
  tree: TreeNode[]
  features: string[]
}

const architectures: ArchMode[] = [
  {
    name: 'Single',
    flag: '--single',
    tagline: 'Go API + embedded React SPA — one binary',
    color: 'sky',
    ideal: 'Laravel/Rails developers, solo devs, simple deploys',
    example: 'https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-single-vite',
    tree: [
      { name: 'cmd/server/main.go', type: 'file', depth: 0, comment: 'go:embed frontend/dist/*' },
      { name: 'internal/', type: 'folder', depth: 0, comment: 'Go backend' },
      { name: 'frontend/', type: 'folder', depth: 0, comment: 'React + Vite + TanStack Router' },
      { name: 'go.mod', type: 'file', depth: 0 },
      { name: 'Makefile', type: 'file', depth: 0, comment: 'make dev, make build' },
    ],
    features: [
      'Single binary deployment via go:embed',
      'Dev: Go on :8080, Vite on :5173 with proxy',
      'make build produces one executable',
      'No Node.js needed in production',
    ],
  },
  {
    name: 'Double',
    flag: '--double',
    tagline: 'Web + API Turborepo monorepo',
    color: 'violet',
    ideal: 'MERN stack developers, API + SPA projects',
    example: 'https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-double-vite',
    tree: [
      { name: 'apps/', type: 'folder', depth: 0 },
      { name: 'api/', type: 'folder', depth: 1, comment: 'Go backend (Gin + GORM)' },
      { name: 'web/', type: 'folder', depth: 1, comment: 'React frontend (Next.js or TanStack)' },
      { name: 'packages/shared/', type: 'folder', depth: 0, comment: 'Types, schemas, constants' },
      { name: 'turbo.json', type: 'file', depth: 0 },
      { name: 'pnpm-workspace.yaml', type: 'file', depth: 0 },
    ],
    features: [
      'Turborepo for parallel builds',
      'Shared TypeScript types + Zod schemas',
      'Independent deployment of API and web',
      'pnpm workspaces for dependency management',
    ],
  },
  {
    name: 'Triple',
    flag: '--triple (default)',
    tagline: 'Web + Admin + API Turborepo monorepo',
    color: 'emerald',
    ideal: 'Full-stack teams, SaaS products, content platforms',
    example: 'https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-triple-next',
    tree: [
      { name: 'apps/', type: 'folder', depth: 0 },
      { name: 'api/', type: 'folder', depth: 1, comment: 'Go backend (Gin + GORM)' },
      { name: 'web/', type: 'folder', depth: 1, comment: 'Public-facing frontend' },
      { name: 'admin/', type: 'folder', depth: 1, comment: 'Admin panel (DataTable, FormBuilder)' },
      { name: 'packages/shared/', type: 'folder', depth: 0, comment: 'Types, schemas, constants' },
      { name: 'turbo.json', type: 'file', depth: 0 },
      { name: 'pnpm-workspace.yaml', type: 'file', depth: 0 },
    ],
    features: [
      'Full admin panel with DataTable + FormBuilder',
      'Resource definitions for zero-code CRUD',
      'Dashboard widgets, system pages',
      'grit generate creates Go + admin page',
    ],
  },
  {
    name: 'API Only',
    flag: '--api',
    tagline: 'Go API with no frontend',
    color: 'amber',
    ideal: 'Microservices, backend teams, mobile-first apps',
    example: 'https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-api-only',
    tree: [
      { name: 'apps/', type: 'folder', depth: 0 },
      { name: 'api/', type: 'folder', depth: 1, comment: 'Go backend (Gin + GORM)' },
      { name: 'cmd/server/', type: 'folder', depth: 2 },
      { name: 'internal/', type: 'folder', depth: 2 },
      { name: 'docker-compose.yml', type: 'file', depth: 0 },
    ],
    features: [
      'Minimal footprint — Go only',
      'All batteries included (auth, storage, jobs)',
      'No Node.js, no frontend build step',
      'Perfect for REST/gRPC APIs',
    ],
  },
  {
    name: 'Mobile',
    flag: '--mobile',
    tagline: 'API + Expo React Native',
    color: 'rose',
    ideal: 'Mobile-first products, cross-platform apps',
    example: 'https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-mobile-expo',
    tree: [
      { name: 'apps/', type: 'folder', depth: 0 },
      { name: 'api/', type: 'folder', depth: 1, comment: 'Go backend' },
      { name: 'expo/', type: 'folder', depth: 1, comment: 'React Native (Expo)' },
      { name: 'packages/shared/', type: 'folder', depth: 0, comment: 'Shared types' },
      { name: 'turbo.json', type: 'file', depth: 0 },
    ],
    features: [
      'Expo managed workflow',
      'Shared types between API and mobile',
      'Same auth system (JWT)',
      'File upload with presigned URLs',
    ],
  },
  {
    name: 'Multi-Client (combine flags)',
    flag: '--desktop (combinable)',
    tagline: 'Add a Wails desktop app to any architecture',
    color: 'cyan',
    ideal: 'SaaS products needing native + mobile + web clients — the Linear / Notion / Slack pattern',
    example: '/docs/concepts/architecture-modes/multi-client',
    tree: [
      { name: 'apps/', type: 'folder', depth: 0 },
      { name: 'api/', type: 'folder', depth: 1, comment: 'Shared Go backend' },
      { name: 'web/', type: 'folder', depth: 1, comment: 'Next.js or TanStack (optional)' },
      { name: 'admin/', type: 'folder', depth: 1, comment: 'Admin panel (optional)' },
      { name: 'expo/', type: 'folder', depth: 1, comment: 'Mobile (optional)' },
      { name: 'desktop/', type: 'folder', depth: 1, comment: 'Wails desktop (always)' },
      { name: 'packages/shared/', type: 'folder', depth: 0, comment: 'Types shared across ALL clients' },
      { name: 'turbo.json', type: 'file', depth: 0 },
    ],
    features: [
      'Combine --desktop with --triple, --double, --mobile, or --api',
      'Frameless Wails window (platform-aware window controls)',
      'Command palette (⌘K) built in',
      'OS keychain token storage (not localStorage)',
      'Not compatible with --single (single apps bundle their own SPA)',
      'Distinct from grit new-desktop (standalone offline-first)',
    ],
  },
]

export default function ArchitectureModesPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="max-w-4xl mx-auto py-12 px-6 lg:px-8">
          <div className="mb-14">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">
              Core Concepts
            </p>
            <h1 className="text-3xl lg:text-4xl font-bold tracking-tight text-foreground mb-4">
              Architecture Modes
            </h1>
            <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl">
              Grit supports 5 architecture modes plus a <code className="text-base font-mono bg-accent/50 px-1.5 py-0.5 rounded">--desktop</code> flag
              that combines with any of them (except <code className="text-base font-mono bg-accent/50 px-1.5 py-0.5 rounded">--single</code>).
              Choose the one that fits your team, your deployment target, and the frameworks you already know.
            </p>
            <LaneFlow
              id="modes-overview"
              lanes={['Clients — vary by mode', 'The constant: one Go API']}
              nodes={[
                { id: 'web', lane: 0, row: 0, title: 'Web app', sub: 'Next.js / Vite', tone: 'blue' },
                { id: 'admin', lane: 0, row: 1, title: 'Admin panel', sub: 'role-gated', tone: 'cyan' },
                { id: 'mobile', lane: 0, row: 2, title: 'Mobile', sub: 'Expo', tone: 'green' },
                { id: 'desktop', lane: 0, row: 3, title: 'Desktop', sub: 'Wails', tone: 'violet' },
                { id: 'api', lane: 1, row: 1, title: 'Go API', sub: 'Gin · GORM · batteries', tone: 'primary' },
                { id: 'data', lane: 1, row: 2, title: 'Data & services', sub: 'Postgres · Redis · S3', tone: 'amber' },
              ]}
              edges={[
                { from: 'web', to: 'api', label: 'REST + JWT', tone: 'blue' },
                { from: 'admin', to: 'api', tone: 'cyan' },
                { from: 'mobile', to: 'api', tone: 'green' },
                { from: 'desktop', to: 'api', tone: 'violet' },
                { from: 'api', to: 'data', tone: 'amber' },
              ]}
              legend={[
                { tone: 'primary', label: 'Same in every mode' },
                { tone: 'blue', label: 'Clients the mode includes' },
              ]}
              caption="Every mode ships the same Go API core — the mode just decides which clients come with it"
            />
          </div>

          <div className="space-y-10">
            {/* Architecture cards */}
            {architectures.map((arch) => (
              <div key={arch.name} className="rounded-xl border border-border/40 bg-accent/20 overflow-hidden">
                <div className="p-6 border-b border-border/30">
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className="text-xl font-semibold text-foreground">{arch.name}</h3>
                    <code className="text-[11px] font-mono text-muted-foreground/70 bg-accent/30 px-2 py-0.5 rounded">
                      {arch.flag}
                    </code>
                  </div>
                  <p className="text-muted-foreground">{arch.tagline}</p>
                  <p className="text-sm text-muted-foreground/70 mt-1">Ideal for: {arch.ideal}</p>
                  {arch.example && (
                    arch.example.startsWith('/') ? (
                      <Link href={arch.example} className="inline-flex items-center gap-1.5 mt-3 text-xs font-medium text-primary hover:text-primary/80 transition-colors">
                        Read the multi-client guide &rarr;
                      </Link>
                    ) : (
                      <a href={arch.example} target="_blank" rel="noreferrer" className="inline-flex items-center gap-1.5 mt-3 text-xs font-medium text-primary hover:text-primary/80 transition-colors">
                        View example project (Job Portal) &rarr;
                      </a>
                    )
                  )}
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-2 divide-y lg:divide-y-0 lg:divide-x divide-border/30">
                  <div className="p-6">
                    <h4 className="text-xs font-mono text-muted-foreground/70 uppercase tracking-wider mb-1">Structure</h4>
                    <FileTree title="my-app/" nodes={arch.tree} />
                  </div>
                  <div className="p-6">
                    <h4 className="text-xs font-mono text-muted-foreground/70 uppercase tracking-wider mb-3">Features</h4>
                    <ul className="space-y-2">
                      {arch.features.map((f, i) => (
                        <li key={i} className="flex items-start gap-2 text-sm text-muted-foreground">
                          <span className="text-primary mt-0.5">-</span>
                          {f}
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* Frontend choice */}
          <div className="mt-14 mb-14">
            <h2 className="text-xl font-semibold text-foreground mb-6">Frontend Framework Choice</h2>
            <p className="text-muted-foreground mb-6">
              For any architecture that includes a frontend (single, double, triple), you can choose between:
            </p>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="rounded-lg border border-border/40 bg-accent/20 p-6">
                <div className="flex items-center gap-2 mb-3">
                  <h3 className="font-semibold text-foreground">Next.js</h3>
                  <code className="text-[10px] font-mono text-muted-foreground/70 bg-accent/30 px-1.5 py-0.5 rounded">--next</code>
                </div>
                <ul className="space-y-1.5 text-sm text-muted-foreground">
                  <li>- Server-side rendering (SSR)</li>
                  <li>- SEO-friendly by default</li>
                  <li>- App Router with layouts</li>
                  <li>- Larger bundle, Node.js runtime</li>
                </ul>
              </div>
              <div className="rounded-lg border border-primary/20 bg-primary/5 p-6">
                <div className="flex items-center gap-2 mb-3">
                  <h3 className="font-semibold text-foreground">TanStack Router</h3>
                  <code className="text-[10px] font-mono text-primary/70 bg-primary/10 px-1.5 py-0.5 rounded">--vite</code>
                </div>
                <ul className="space-y-1.5 text-sm text-muted-foreground">
                  <li>- Vite — instant HMR, fast builds</li>
                  <li>- Small bundle size (SPA)</li>
                  <li>- File-based routing via plugin</li>
                  <li>- No Node.js server needed</li>
                </ul>
              </div>
            </div>
          </div>

          {/* Example Projects */}
          <div className="mt-14 mb-14">
            <h2 className="text-xl font-semibold text-foreground mb-4">Example Projects</h2>
            <p className="text-muted-foreground mb-6">
              We built the same <strong>Job Portal</strong> app with every architecture so you can see exactly how each one works.
              Each example includes a README, step-by-step guide, .env template, and production Docker Compose.
            </p>
            <div className="rounded-lg border border-border/40 bg-accent/20 p-6">
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                {[
                  { label: 'Triple + Next.js', href: 'https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-triple-next' },
                  { label: 'Triple + TanStack', href: 'https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-triple-vite' },
                  { label: 'Double + TanStack', href: 'https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-double-vite' },
                  { label: 'Single (one binary)', href: 'https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-single-vite' },
                  { label: 'API Only', href: 'https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-api-only' },
                  { label: 'Mobile + Expo', href: 'https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-mobile-expo' },
                ].map((ex) => (
                  <a key={ex.label} href={ex.href} target="_blank" rel="noreferrer" className="flex items-center justify-between rounded-md border border-border/30 bg-background/50 px-4 py-3 text-sm text-muted-foreground hover:text-foreground hover:border-primary/30 transition-colors">
                    {ex.label}
                    <span className="text-primary text-xs">&rarr;</span>
                  </a>
                ))}
              </div>
              <p className="text-xs text-muted-foreground/60 mt-4">
                All examples: <a href="https://github.com/MUKE-coder/grit/tree/main/examples" target="_blank" rel="noreferrer" className="text-primary hover:underline">github.com/MUKE-coder/grit/tree/main/examples</a>
              </p>
            </div>
          </div>

          {/* Nav */}
          <div className="flex items-center justify-between pt-8 border-t border-border/40">
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/concepts/architecture" className="gap-1.5">
                <ArrowLeft className="h-3.5 w-3.5" />
                Architecture Overview
              </Link>
            </Button>
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/concepts/cli" className="gap-1.5">
                CLI Commands
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
