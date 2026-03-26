import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/architecture-modes')

export default function ArchitectureModesPage() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="max-w-4xl mx-auto py-12 px-6 lg:px-8">
          <div className="mb-14">
            <p className="text-sm font-mono font-medium text-sky-400 mb-3 tracking-wide uppercase">
              Core Concepts
            </p>
            <h1 className="text-3xl lg:text-4xl font-bold tracking-tight text-white mb-4">
              Architecture Modes
            </h1>
            <p className="text-lg text-slate-400 leading-relaxed max-w-2xl">
              Grit supports 5 architecture modes. Choose the one that fits your team,
              your deployment target, and the frameworks you already know.
            </p>
          </div>

          <div className="space-y-10">
            {/* Architecture cards */}
            {[
              {
                name: 'Single',
                flag: '--single',
                tagline: 'Go API + embedded React SPA — one binary',
                color: 'sky',
                ideal: 'Laravel/Rails developers, solo devs, simple deploys',
                structure: `my-app/
├── cmd/server/main.go   # go:embed frontend/dist/*
├── internal/            # Go backend
├── frontend/            # React + Vite + TanStack Router
├── go.mod
└── Makefile             # make dev, make build`,
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
                structure: `my-app/
├── apps/
│   ├── api/             # Go backend (Gin + GORM)
│   └── web/             # React frontend (Next.js or TanStack)
├── packages/shared/     # Types, schemas, constants
├── turbo.json
└── pnpm-workspace.yaml`,
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
                structure: `my-app/
├── apps/
│   ├── api/             # Go backend (Gin + GORM)
│   ├── web/             # Public-facing frontend
│   └── admin/           # Admin panel (DataTable, FormBuilder)
├── packages/shared/     # Types, schemas, constants
├── turbo.json
└── pnpm-workspace.yaml`,
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
                structure: `my-app/
├── apps/
│   └── api/             # Go backend (Gin + GORM)
│       ├── cmd/server/
│       └── internal/
└── docker-compose.yml`,
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
                structure: `my-app/
├── apps/
│   ├── api/             # Go backend
│   └── expo/            # React Native (Expo)
├── packages/shared/     # Shared types
└── turbo.json`,
                features: [
                  'Expo managed workflow',
                  'Shared types between API and mobile',
                  'Same auth system (JWT)',
                  'File upload with presigned URLs',
                ],
              },
            ].map((arch) => (
              <div key={arch.name} className="rounded-xl border border-white/[0.06] bg-white/[0.02] overflow-hidden">
                <div className="p-6 border-b border-white/[0.04]">
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className="text-xl font-semibold text-white">{arch.name}</h3>
                    <code className="text-[11px] font-mono text-slate-500 bg-white/[0.04] px-2 py-0.5 rounded">
                      {arch.flag}
                    </code>
                  </div>
                  <p className="text-slate-400">{arch.tagline}</p>
                  <p className="text-sm text-slate-500 mt-1">Ideal for: {arch.ideal}</p>
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-2 divide-y lg:divide-y-0 lg:divide-x divide-white/[0.04]">
                  <div className="p-6">
                    <h4 className="text-xs font-mono text-slate-500 uppercase tracking-wider mb-3">Structure</h4>
                    <pre className="text-[12px] font-mono text-slate-400 leading-5 whitespace-pre">{arch.structure}</pre>
                  </div>
                  <div className="p-6">
                    <h4 className="text-xs font-mono text-slate-500 uppercase tracking-wider mb-3">Features</h4>
                    <ul className="space-y-2">
                      {arch.features.map((f, i) => (
                        <li key={i} className="flex items-start gap-2 text-sm text-slate-400">
                          <span className="text-sky-400 mt-0.5">-</span>
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
            <h2 className="text-xl font-semibold text-white mb-6">Frontend Framework Choice</h2>
            <p className="text-slate-400 mb-6">
              For any architecture that includes a frontend (single, double, triple), you can choose between:
            </p>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="rounded-lg border border-white/[0.06] bg-white/[0.02] p-6">
                <div className="flex items-center gap-2 mb-3">
                  <h3 className="font-semibold text-white">Next.js</h3>
                  <code className="text-[10px] font-mono text-slate-500 bg-white/[0.04] px-1.5 py-0.5 rounded">--next</code>
                </div>
                <ul className="space-y-1.5 text-sm text-slate-400">
                  <li>- Server-side rendering (SSR)</li>
                  <li>- SEO-friendly by default</li>
                  <li>- App Router with layouts</li>
                  <li>- Larger bundle, Node.js runtime</li>
                </ul>
              </div>
              <div className="rounded-lg border border-sky-400/20 bg-sky-400/5 p-6">
                <div className="flex items-center gap-2 mb-3">
                  <h3 className="font-semibold text-white">TanStack Router</h3>
                  <code className="text-[10px] font-mono text-sky-400/70 bg-sky-400/10 px-1.5 py-0.5 rounded">--vite</code>
                </div>
                <ul className="space-y-1.5 text-sm text-slate-400">
                  <li>- Vite — instant HMR, fast builds</li>
                  <li>- Small bundle size (SPA)</li>
                  <li>- File-based routing via plugin</li>
                  <li>- No Node.js server needed</li>
                </ul>
              </div>
            </div>
          </div>

          {/* Nav */}
          <div className="flex items-center justify-between pt-8 border-t border-white/[0.06]">
            <Button variant="ghost" size="sm" asChild className="text-slate-500 hover:text-white">
              <Link href="/docs/concepts/architecture" className="gap-1.5">
                <ArrowLeft className="h-3.5 w-3.5" />
                Architecture Overview
              </Link>
            </Button>
            <Button variant="ghost" size="sm" asChild className="text-slate-500 hover:text-white">
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
