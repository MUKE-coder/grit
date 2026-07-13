import Link from 'next/link'
import { ArrowLeft, ArrowRight, Download } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/desktop/getting-started')

export default function DesktopGettingStartedPage() {
  return (
    <div className="min-h-screen bg-background isolate">
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
                Getting Started with Desktop
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Scaffold a native desktop application with one command. This guide walks you
                through prerequisites, project creation, development workflow, and running
                GORM Studio.
              </p>
              <LaneFlow
                id="desk-gs"
                lanes={['1 · Scaffold', '2 · Dev', '3 · Native window']}
                nodes={[
                  { id: 'new', lane: 0, row: 0, title: 'grit new-desktop', sub: 'Wails + Go', tone: 'primary' },
                  { id: 'dev', lane: 1, row: 0, title: 'grit start desktop', sub: 'hot reload', tone: 'cyan' },
                  { id: 'win', lane: 2, row: 0, title: 'Native window', sub: 'embedded API', tone: 'green' },
                  { id: 'studio', lane: 2, row: 1, title: 'GORM Studio', sub: 'browse data', tone: 'amber' },
                ]}
                edges={[
                  { from: 'new', to: 'dev', tone: 'cyan' },
                  { from: 'dev', to: 'win', label: 'launches', tone: 'green' },
                  { from: 'dev', to: 'studio', tone: 'amber' },
                ]}
                legend={[{ tone: 'primary', label: 'One command' }, { tone: 'green', label: 'Native app' }]}
                caption="One command scaffolds a native Wails app with an embedded Go API and GORM Studio"
              />
            </div>

            {/* Prerequisites */}
            <div className="prose-grit mb-10">
              <h2>Prerequisites</h2>
              <p>
                Make sure the following tools are installed before creating a desktop project:
              </p>
            </div>

            <div className="grid gap-3 sm:grid-cols-2 mb-10">
              {[
                { name: 'Go', version: '1.21+', check: 'go version' },
                { name: 'Node.js', version: '18+', check: 'node --version' },
                { name: 'Grit CLI', version: 'Latest', check: 'grit --help' },
                { name: 'Wails CLI', version: 'v2', check: 'wails version' },
              ].map((tool) => (
                <div
                  key={tool.name}
                  className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                >
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-[15px] font-semibold">{tool.name}</span>
                    <span className="text-sm font-mono text-primary/60">
                      {tool.version}
                    </span>
                  </div>
                  <code className="text-sm font-mono text-muted-foreground/50">
                    {tool.check}
                  </code>
                </div>
              ))}
            </div>

            <div className="prose-grit mb-4">
              <p>
                Install the Wails CLI if you don&apos;t have it yet:
              </p>
            </div>
            <CodeBlock terminal code="go install github.com/wailsapp/wails/v2/cmd/wails@latest" className="mb-4" />
            <div className="prose-grit mb-10">
              <p>
                Run <code>wails doctor</code> to verify your environment. It checks for
                Go, Node.js, npm/pnpm, and platform-specific build tools (GCC on Linux,
                Xcode on macOS, or WebView2 on Windows).
              </p>
            </div>

            {/* Step 1 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  1
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Scaffold the Project
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Create a new desktop project with the Grit CLI. This generates a complete
                  Wails application with Go backend, React frontend, authentication, CRUD
                  resources, and all batteries included.
                </p>
              </div>
              <CodeBlock terminal code="grit new-desktop myapp" className="mb-0 glow-purple-sm" />
            </div>

            {/* Step 2: What you get */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  2
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Project Structure
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  The scaffolded project has the following layout. Go code lives at the root
                  and in <code>internal/</code>, while the React frontend lives in <code>frontend/</code>.
                  Alongside the Wails bindings, an embedded Gin REST API runs on{' '}
                  <code>127.0.0.1:34999</code> (see <code>internal/api/</code>) to serve native file
                  uploads, a health check, and any external clients.
                </p>
              </div>
              <CodeBlock language="bash" filename="myapp/" code={`myapp/
├── main.go                  # Wails entry point
├── app.go                   # App struct with bound methods
├── wails.json               # Wails project configuration
├── go.mod
├── go.sum
├── internal/
│   ├── config/
│   │   └── config.go        # App configuration
│   ├── db/
│   │   └── db.go            # GORM database setup (SQLite)
│   ├── models/
│   │   ├── user.go          # User model + AutoMigrate
│   │   ├── blog.go          # Blog post model
│   │   └── contact.go       # Contact model
│   ├── service/
│   │   ├── auth.go          # Authentication service
│   │   ├── blog.go          # Blog CRUD service
│   │   ├── contact.go       # Contact CRUD service
│   │   └── export.go        # PDF / Excel export service
│   ├── api/
│   │   ├── router.go        # Embedded Gin REST API (127.0.0.1:34999)
│   │   └── uploads.go       # POST /api/uploads, GET /api/health
│   ├── storage/
│   │   └── storage.go       # Writes uploads into the OS app-data dir
│   ├── files/
│   │   └── file_ref.go      # FileRef — the stored upload's JSON shape
│   └── types/
│       └── types.go         # Shared request/response types
├── frontend/
│   ├── src/
│   │   ├── main.tsx          # React entry point (TanStack Router)
│   │   ├── routes/           # File-based routes (TanStack Router)
│   │   │   ├── __root.tsx    # Root route
│   │   │   ├── _layout.tsx   # Auth guard + sidebar layout
│   │   │   └── _layout/      # Protected page routes
│   │   ├── components/       # Reusable UI components
│   │   ├── hooks/            # TanStack Query hooks
│   │   └── lib/              # Utilities
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   └── tailwind.config.ts
└── cmd/
    └── studio/
        └── main.go           # GORM Studio standalone server`} className="mb-0" />
            </div>

            {/* How File-Based Routing Works */}
            <div className="mb-10">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                How File-Based Routing Works
              </h2>
              <div className="prose-grit mb-4">
                <p>
                  Grit Desktop uses <a href="https://tanstack.com/router" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">TanStack Router</a> with
                  file-based routing. Instead of a centralized route configuration, each page
                  is a route file in <code>routes/</code>. The TanStack Router Vite plugin
                  auto-discovers these files and generates a type-safe route tree.
                </p>
              </div>

              <div className="overflow-x-auto mb-6">
                <table className="w-full text-sm border border-border rounded-lg">
                  <thead>
                    <tr className="border-b border-border bg-muted/50">
                      <th className="text-left p-3 font-medium">File</th>
                      <th className="text-left p-3 font-medium">URL</th>
                      <th className="text-left p-3 font-medium">Purpose</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border text-muted-foreground">
                    {[
                      ['routes/__root.tsx', '(wrapper)', 'Root layout — wraps everything, renders Toaster'],
                      ['routes/login.tsx', '/login', 'Login page (public)'],
                      ['routes/register.tsx', '/register', 'Register page (public)'],
                      ['routes/_layout.tsx', '(wrapper)', 'Auth guard + sidebar — no URL segment added'],
                      ['routes/_layout/index.tsx', '/', 'Dashboard (protected)'],
                      ['routes/_layout/blogs.index.tsx', '/blogs', 'Blog list page'],
                      ['routes/_layout/blogs.new.tsx', '/blogs/new', 'Blog create form'],
                      ['routes/_layout/blogs.$id.edit.tsx', '/blogs/:id/edit', 'Blog edit form'],
                    ].map(([file, url, purpose]) => (
                      <tr key={file}>
                        <td className="p-3 font-mono text-xs">{file}</td>
                        <td className="p-3 font-mono text-xs text-primary">{url}</td>
                        <td className="p-3">{purpose}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <div className="prose-grit mb-4">
                <p>
                  Key conventions:
                </p>
              </div>

              <div className="space-y-2 mb-4">
                {[
                  'The __root.tsx file is the topmost layout. It renders an Outlet where child routes appear.',
                  'Files prefixed with _ (like _layout.tsx) create pathless layouts — they wrap child routes without adding a URL segment.',
                  'Dot notation in filenames creates nested paths: blogs.index.tsx becomes /blogs, blogs.$id.edit.tsx becomes /blogs/:id/edit.',
                  'The $ prefix denotes dynamic parameters: $id in the filename becomes a typed route parameter accessible via Route.useParams().',
                  'Each route file exports a Route constant created with createFileRoute() — this registers the route and its component.',
                ].map((item, i) => (
                  <div
                    key={i}
                    className="flex items-start gap-3 rounded-lg border border-border/30 bg-card/30 px-4 py-2.5"
                  >
                    <span className="text-sm font-mono text-primary/60 shrink-0 mt-0.5">
                      {i + 1}.
                    </span>
                    <span className="text-sm text-muted-foreground">
                      {item}
                    </span>
                  </div>
                ))}
              </div>

              <div className="prose-grit mb-0">
                <blockquote>
                  The route tree (<code>routeTree.gen.ts</code>) is auto-generated by the
                  TanStack Router Vite plugin and gitignored. You never edit it manually.
                  Just create or delete route files — the plugin handles the rest.
                </blockquote>
              </div>
            </div>

            {/* Step 3: Development */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  3
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Start Development
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Navigate into the project and start Wails in development mode. This launches
                  the desktop window with hot-reload for both Go and React code.
                </p>
              </div>
              <CodeBlock terminal code={`cd myapp
wails dev`} className="mb-0 glow-purple-sm" />
              <div className="prose-grit mt-4">
                <p>
                  Alternatively, if you prefer using the Grit CLI:
                </p>
              </div>
              <CodeBlock terminal code="grit start" className="mb-0" />
              <div className="prose-grit mt-4">
                <p>
                  The app window opens automatically. Changes to Go files trigger a rebuild,
                  and changes to React files trigger a Vite HMR refresh. The frontend dev
                  server runs on <code>http://localhost:34115</code> during development.
                </p>
                <blockquote>
                  On first run, <code>wails dev</code> installs frontend dependencies
                  automatically. Subsequent starts are much faster.
                </blockquote>
              </div>
            </div>

            {/* Step 4: Studio */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  4
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Open GORM Studio
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  GORM Studio provides a visual database browser. For desktop projects, it
                  starts a standalone Gin server with the{' '}
                  <code>gorm-studio</code> package mounted at <code>/studio</code> on port 8080
                  and automatically opens your browser.
                </p>
              </div>
              <CodeBlock terminal code="grit studio" className="mb-0 glow-purple-sm" />
              <div className="prose-grit mt-4">
                <p>
                  Your browser opens automatically at{' '}
                  <code>http://localhost:8080/studio</code> where you can browse tables,
                  inspect rows, and explore your SQLite database visually. This uses the same{' '}
                  <code>gorm-studio</code> package as the web scaffold.
                </p>
              </div>
            </div>

            {/* Default URLs */}
            <div className="mb-10">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Development URLs
              </h2>
              <div className="grid gap-3 sm:grid-cols-2">
                {[
                  {
                    name: 'Desktop App',
                    url: 'Native window',
                    desc: 'Opens automatically with wails dev',
                  },
                  {
                    name: 'Frontend Dev',
                    url: 'http://localhost:34115',
                    desc: 'Vite dev server (also viewable in browser)',
                  },
                  {
                    name: 'GORM Studio',
                    url: 'http://localhost:8080/studio',
                    desc: 'Visual database browser (grit studio)',
                  },
                ].map((item) => (
                  <div
                    key={item.name}
                    className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                  >
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-[15px] font-semibold">{item.name}</span>
                    </div>
                    <code className="text-sm font-mono text-primary/60 block mb-1">
                      {item.url}
                    </code>
                    <span className="text-sm text-muted-foreground/50">
                      {item.desc}
                    </span>
                  </div>
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
                <Link href="/docs/desktop" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Desktop Overview
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/desktop/resource-generation" className="gap-1.5">
                  Resource Generation
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
