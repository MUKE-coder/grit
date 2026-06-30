import Link from "next/link"
import { ArrowLeft, ArrowRight } from "lucide-react"
import { Button } from "@/components/ui/button"
import { SiteHeader } from "@/components/site-header"
import { DocsSidebar } from "@/components/docs-sidebar"
import { CodeBlock, StepWithCode } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/getting-started/installation')

export default function InstallationPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="max-w-5xl mx-auto py-12 px-6 lg:px-8">
          {/* Header */}
          <div className="mb-14">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">
              Getting Started
            </p>
            <h1 className="text-3xl lg:text-4xl font-bold tracking-tight text-foreground mb-4">
              Installation
            </h1>
            <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl">
              Get up and running with Grit in minutes. Install the CLI, choose your architecture,
              and scaffold a production-ready full-stack application.
            </p>
          </div>

          {/* Prerequisites */}
          <div className="mb-14">
            <h2 className="text-xl font-semibold text-foreground mb-6">System Requirements</h2>
            <div className="overflow-x-auto rounded-lg border border-border/40">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/40 bg-accent/20">
                    <th className="px-4 py-3 text-left font-medium text-foreground/80">Tool</th>
                    <th className="px-4 py-3 text-left font-medium text-foreground/80">Min Version</th>
                    <th className="px-4 py-3 text-left font-medium text-foreground/80 hidden sm:table-cell">Required For</th>
                    <th className="px-4 py-3 text-left font-medium text-foreground/80 hidden md:table-cell">Check</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/30">
                  {[
                    { tool: 'Go', version: '1.21+', for: 'API server', check: 'go version' },
                    { tool: 'Node.js', version: '18+', for: 'Frontend apps', check: 'node --version' },
                    { tool: 'pnpm', version: '8+', for: 'Package management', check: 'pnpm --version' },
                    { tool: 'Docker', version: 'Latest', for: 'Infrastructure', check: 'docker --version' },
                  ].map((row) => (
                    <tr key={row.tool} className="hover:bg-accent/20 transition-colors">
                      <td className="px-4 py-3 font-medium text-foreground">{row.tool}</td>
                      <td className="px-4 py-3 font-mono text-sm text-primary">{row.version}</td>
                      <td className="px-4 py-3 text-muted-foreground hidden sm:table-cell">{row.for}</td>
                      <td className="px-4 py-3 hidden md:table-cell">
                        <code className="text-xs font-mono text-muted-foreground/70 bg-accent/30 px-1.5 py-0.5 rounded">{row.check}</code>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Installation Steps */}
          <div className="mb-14">
            <h2 className="text-xl font-semibold text-foreground mb-8">Install the Grit CLI</h2>

            <StepWithCode
              number="01"
              title="Install Grit"
              description={
                <>
                  <p>
                    Use the one-line install script. It detects an existing install and
                    runs <code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">grit update</code> (idempotent —
                    no-ops if already on latest); otherwise it downloads the binary from the
                    matching GitHub release for your OS / arch and puts it on PATH.
                  </p>
                  <p className="mt-3 text-sm text-muted-foreground">
                    Prefer Go? <code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[12px]">go install github.com/MUKE-coder/grit/v3/cmd/grit@latest</code> does
                    the same thing if you have the Go toolchain.
                  </p>
                </>
              }
              code={`# macOS / Linux
curl -fsSL https://gritframework.dev/install.sh | sh

# Windows (PowerShell)
iwr -useb https://gritframework.dev/install.ps1 | iex`}
              filename="Terminal"
              language="bash"
            />

            <StepWithCode
              number="02"
              title="Create your project"
              description={
                <p>
                  Run <code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">grit new</code> and
                  follow the interactive prompts to select your architecture and frontend framework.
                </p>
              }
              code={`grit new my-app

? Select architecture:
  > Triple — Web + Admin + API (Turborepo)
    Double — Web + API (Turborepo)
    Single — Go API + embedded React SPA
    API Only — Go API (no frontend)
    Mobile — API + Expo (React Native)

? Select frontend framework:
  > Next.js — SSR, SEO, App Router
    TanStack Router — Vite, fast builds, small bundle`}
              filename="Terminal"
              language="bash"
            />

            <StepWithCode
              number="03"
              title="Start infrastructure"
              description={
                <p>
                  Start PostgreSQL, Redis, MinIO, and Mailhog with Docker Compose.
                  These services are pre-configured in the generated{' '}
                  <code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">docker-compose.yml</code>.
                </p>
              }
              code={`cd my-app
docker compose up -d`}
              filename="Terminal"
              language="bash"
            />

            <StepWithCode
              number="04"
              title="Install dependencies & start"
              description={
                <p>
                  Install frontend dependencies and start all development servers.
                  The Go API runs on <code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">:8080</code>,
                  web app on <code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">:3000</code>,
                  and admin on <code className="text-primary bg-accent/30 px-1.5 py-0.5 rounded text-[13px]">:3001</code>.
                </p>
              }
              code={`pnpm install
pnpm dev`}
              filename="Terminal"
              language="bash"
            />
          </div>

          {/* Quick command reference */}
          <div className="mb-14">
            <h2 className="text-xl font-semibold text-foreground mb-4">Architecture shortcuts</h2>
            <p className="text-muted-foreground mb-6">
              Skip the interactive prompts with flags:
            </p>
            <CodeBlock language="bash" filename="Terminal" code={`# Triple (default) with Next.js
grit new my-app --triple --next

# Single app with TanStack Router (Vite)
grit new my-app --single --vite

# Double (web + api) with TanStack Router
grit new my-app --double --vite

# API only (no frontend)
grit new my-app --api

# Desktop app (Wails)
grit new-desktop my-app`} />
          </div>

          {/* Services table */}
          <div className="mb-14">
            <h2 className="text-xl font-semibold text-foreground mb-4">Default services</h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {[
                { name: 'Go API', url: 'http://localhost:8080', desc: 'Backend server' },
                { name: 'Web App', url: 'http://localhost:3000', desc: 'Next.js or TanStack' },
                { name: 'Admin Panel', url: 'http://localhost:3001', desc: 'Resource management' },
                { name: 'API Docs', url: 'http://localhost:8080/docs', desc: 'Auto-generated' },
                { name: 'GORM Studio', url: 'http://localhost:8080/studio', desc: 'Database browser' },
                { name: 'Mailhog', url: 'http://localhost:8025', desc: 'Email testing' },
              ].map((svc) => (
                <div key={svc.name} className="rounded-lg border border-border/40 bg-accent/20 p-4 hover:bg-accent/30 transition-colors">
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-sm font-medium text-foreground">{svc.name}</span>
                    <code className="text-[11px] font-mono text-primary/70">{svc.url.replace('http://', '')}</code>
                  </div>
                  <p className="text-xs text-muted-foreground/70">{svc.desc}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Nav */}
          <div className="flex items-center justify-between pt-8 border-t border-border/40">
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/getting-started/quick-start" className="gap-1.5">
                <ArrowLeft className="h-3.5 w-3.5" />
                Quick Start
              </Link>
            </Button>
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/getting-started/project-structure" className="gap-1.5">
                Project Structure
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
