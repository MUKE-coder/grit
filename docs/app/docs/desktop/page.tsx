import Link from 'next/link'
import { ArrowRight, Download, Terminal } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/desktop')

export default function DesktopOverviewPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">
                Desktop (Wails)
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Desktop App Development
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Grit supports building native desktop applications with Wails. One command
                scaffolds a complete desktop project with Go backend, React frontend, authentication,
                CRUD, data export, and GORM Studio -- all compiled into a single executable.
              </p>
            </div>

            <div className="prose-grit">
              {/* Intro */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  What is Grit Desktop?
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit Desktop extends the framework beyond web applications. Using{' '}
                  <a href="https://wails.io" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Wails</a>,
                  it combines a Go backend with a React frontend (Vite + TanStack Router + TanStack Query)
                  to produce native desktop apps for Windows, macOS, and Linux.
                </p>
                <CodeBlock terminal code="grit new-desktop myapp" className="mb-0 glow-purple-sm" />
                <p className="text-muted-foreground leading-relaxed mt-4">
                  This single command generates a fully working desktop application with SQLite or
                  PostgreSQL, JWT authentication, blog and contact CRUD, PDF and Excel export, a
                  custom title bar, dark theme, and GORM Studio for database browsing.
                </p>
              </div>

              {/* Install */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Install Grit
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  You need <strong className="text-foreground/80">Go 1.21+</strong>,{' '}
                  <strong className="text-foreground/80">Node.js 18+</strong>, and{' '}
                  <a href="https://wails.io/docs/gettingstarted/installation" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Wails v2</a>{' '}
                  installed. Then install the Grit CLI:
                </p>
                <CodeBlock
                  terminal
                  code={`# macOS / Linux
curl -fsSL https://gritframework.dev/install.sh | sh

# Windows (PowerShell)
iwr -useb https://gritframework.dev/install.ps1 | iex`}
                  className="mb-4 glow-purple-sm"
                />
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Scaffold a new desktop app and start developing:
                </p>
                <CodeBlock terminal code={`grit new-desktop my-app
cd my-app
wails dev`} className="mb-0 glow-purple-sm" />
              </div>

              {/* Architecture */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Architecture
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Instead of HTTP requests between client and server, Wails uses direct Go bindings.
                  The React frontend calls Go functions directly, with Wails handling the bridge
                  between the two runtimes.
                </p>
                <div className="space-y-3 mb-6">
                  {[
                    { step: '1', title: 'Go Backend (Wails Bindings)', desc: 'Business logic lives in Go structs. Methods on the App struct are automatically exposed to the frontend via Wails bindings. No HTTP server needed.' },
                    { step: '2', title: 'React Frontend (Vite)', desc: 'A Vite-powered React app with TanStack Router (file-based routing) and TanStack Query for state management. Calls Go functions through the generated Wails bindings.' },
                    { step: '3', title: 'SQLite / PostgreSQL', desc: 'GORM handles the database layer. SQLite is the default for portable desktop apps. PostgreSQL is supported for networked setups.' },
                    { step: '4', title: 'Single Binary Output', desc: 'The entire application -- Go backend, React frontend, and all assets -- compiles into a single executable that you distribute.' },
                  ].map((item) => (
                    <div key={item.step} className="flex gap-3 p-3 rounded-lg border border-border/30 bg-card/30">
                      <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                        {item.step}
                      </div>
                      <div>
                        <p className="text-sm font-medium">{item.title}</p>
                        <p className="text-xs text-muted-foreground/70 mt-0.5">{item.desc}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* TanStack Router */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Routing with TanStack Router
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit Desktop uses{' '}
                  <a href="https://tanstack.com/router" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">TanStack Router</a>{' '}
                  with <strong className="text-foreground/80">file-based routing</strong>. Instead of declaring
                  routes in a centralized file, each page is a self-contained route file in the{' '}
                  <code>routes/</code> directory. The TanStack Router Vite plugin automatically discovers
                  route files and generates a type-safe route tree at build time.
                </p>
                <div className="space-y-3 mb-6">
                  {[
                    { step: '1', title: 'File-Based Routes', desc: 'Each page is a file in routes/_layout/. The filename determines the URL: blogs.index.tsx maps to /blogs, blogs.$id.edit.tsx maps to /blogs/:id/edit. No route registry to maintain.' },
                    { step: '2', title: 'Type-Safe Navigation', desc: 'Route params, search params, and navigation are fully typed. Route.useParams() returns typed params scoped to the current route. navigate({ to: "/blogs/$id/edit", params: { id } }) is validated at compile time.' },
                    { step: '3', title: 'Pathless Layouts', desc: 'The _layout.tsx file creates a layout route (auth guard + sidebar) without adding a URL segment. All pages inside _layout/ inherit this wrapper automatically.' },
                    { step: '4', title: 'Zero-Config for Generation', desc: 'When grit generate resource creates a new resource, it simply creates route files. No imports to add, no route registry to update. Deleting a resource just deletes the files.' },
                  ].map((item) => (
                    <div key={item.step} className="flex gap-3 p-3 rounded-lg border border-border/30 bg-card/30">
                      <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                        {item.step}
                      </div>
                      <div>
                        <p className="text-sm font-medium">{item.title}</p>
                        <p className="text-xs text-muted-foreground/70 mt-0.5">{item.desc}</p>
                      </div>
                    </div>
                  ))}
                </div>
                <p className="text-muted-foreground leading-relaxed">
                  TanStack Router uses{' '}
                  <code>createHashHistory()</code> for desktop apps, so routes work correctly when
                  the app is loaded from disk (no web server needed). The Vite plugin generates{' '}
                  <code>routeTree.gen.ts</code> automatically — this file is gitignored and
                  regenerated on every build.
                </p>
              </div>

              {/* Features */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Built-in Features
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every scaffolded desktop project includes the following out of the box:
                </p>
                <div className="grid gap-3 sm:grid-cols-2 mb-6">
                  {[
                    { name: 'Authentication', desc: 'Register and login with JWT tokens and bcrypt password hashing' },
                    { name: 'Blog CRUD', desc: 'Full create, read, update, delete for blog posts with rich text' },
                    { name: 'Contact CRUD', desc: 'Contact management with form validation and search' },
                    { name: 'PDF Export', desc: 'Export data tables to styled PDF documents' },
                    { name: 'Excel Export', desc: 'Export data to .xlsx spreadsheets' },
                    { name: 'GORM Studio', desc: 'Built-in database browser at localhost:8080/studio' },
                    { name: 'Custom Title Bar', desc: 'Frameless window with custom drag region and window controls' },
                    { name: 'Dark Theme', desc: 'Grit dark theme with Tailwind CSS and shadcn/ui components' },
                    { name: 'Resource Generation', desc: 'Generate new resources with grit generate resource' },
                    { name: 'Toast Notifications', desc: 'User feedback with sonner toast notifications' },
                  ].map((feature) => (
                    <div
                      key={feature.name}
                      className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                    >
                      <span className="text-[15px] font-semibold block mb-1">{feature.name}</span>
                      <span className="text-sm text-muted-foreground/60">{feature.desc}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Comparison Table */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Web vs Desktop
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Both modes share Grit conventions (GORM models, code generation, Studio), but
                  the runtime and distribution model are different.
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Aspect</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Web (grit new)</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Desktop (grit new-desktop)</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs font-medium">Backend</td>
                        <td className="px-4 py-2.5 text-xs">Gin HTTP server</td>
                        <td className="px-4 py-2.5 text-xs">Wails bindings (direct Go calls)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs font-medium">Frontend</td>
                        <td className="px-4 py-2.5 text-xs">Next.js (App Router)</td>
                        <td className="px-4 py-2.5 text-xs">Vite + React + TanStack Router</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs font-medium">Database</td>
                        <td className="px-4 py-2.5 text-xs">PostgreSQL</td>
                        <td className="px-4 py-2.5 text-xs">SQLite (default) or PostgreSQL</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs font-medium">Project Structure</td>
                        <td className="px-4 py-2.5 text-xs">Turborepo monorepo</td>
                        <td className="px-4 py-2.5 text-xs">Single project directory</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs font-medium">State Management</td>
                        <td className="px-4 py-2.5 text-xs">React Query + fetch</td>
                        <td className="px-4 py-2.5 text-xs">React Query + Wails bindings</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 text-xs font-medium">Distribution</td>
                        <td className="px-4 py-2.5 text-xs">Deploy to cloud / VPS</td>
                        <td className="px-4 py-2.5 text-xs">Distribute .exe / .app / binary</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>

            {/* What's Next */}
            <div className="mb-8">
              <h2 className="text-2xl font-semibold tracking-tight mb-6">
                What&apos;s Next?
              </h2>
              <div className="grid gap-3 sm:grid-cols-2">
                {[
                  {
                    title: 'Getting Started',
                    desc: 'Install prerequisites and scaffold your first desktop app',
                    href: '/docs/desktop/getting-started',
                  },
                  {
                    title: 'Resource Generation',
                    desc: 'Generate full-stack CRUD resources for desktop',
                    href: '/docs/desktop/resource-generation',
                  },
                  {
                    title: 'Building & Distribution',
                    desc: 'Compile and distribute your desktop app',
                    href: '/docs/desktop/building',
                  },
                ].map((item) => (
                  <Link key={item.href} href={item.href}>
                    <div className="group rounded-lg border border-border/40 bg-card/50 p-4 hover:border-primary/20 hover:bg-card/80 transition-all duration-200">
                      <h3 className="text-[15px] font-semibold mb-1 group-hover:text-primary transition-colors flex items-center gap-1.5">
                        {item.title}
                        <ArrowRight className="h-3.5 w-3.5 opacity-0 -translate-x-1 group-hover:opacity-100 group-hover:translate-x-0 transition-all" />
                      </h3>
                      <p className="text-sm text-muted-foreground/60">
                        {item.desc}
                      </p>
                    </div>
                  </Link>
                ))}
              </div>
            </div>

            {/* Desktop Handbook */}
            <div className="mb-8">
              <a
                href="https://14j7oh8kso.ufs.sh/f/HLxTbDBCDLwfpiJDPD3QgNG9hYzVFo5iLR0yrDPTJedWnBH7"
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-3 rounded-lg border border-primary/25 bg-primary/5 px-5 py-4 hover:bg-primary/10 transition-colors group"
              >
                <Download className="h-5 w-5 text-primary/70 group-hover:text-primary shrink-0" />
                <div>
                  <span className="text-sm font-semibold text-foreground block">Download Desktop Handbook (PDF)</span>
                  <span className="text-xs text-muted-foreground/60">Complete offline reference for Grit Desktop development</span>
                </div>
              </a>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/getting-started/quick-start" className="gap-1.5">
                  Web Quick Start
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/desktop/getting-started" className="gap-1.5">
                  Getting Started
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
