import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/desktop/getting-started')

export default function DesktopGettingStartedPage() {
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
                Getting Started with Desktop
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Scaffold a native desktop application with one command. This guide walks you
                through prerequisites, project creation, development workflow, and running
                GORM Studio.
              </p>
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
│   ├── services/
│   │   ├── auth.go          # Authentication service
│   │   ├── blog.go          # Blog CRUD service
│   │   └── contact.go       # Contact CRUD service
│   └── types/
│       └── types.go         # Shared request/response types
├── frontend/
│   ├── src/
│   │   ├── App.tsx           # React Router setup
│   │   ├── main.tsx          # React entry point
│   │   ├── pages/            # Page components
│   │   ├── components/       # Reusable UI components
│   │   ├── hooks/            # TanStack Query hooks
│   │   └── lib/              # Utilities
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   └── tailwind.config.js
└── cmd/
    └── studio/
        └── main.go           # GORM Studio standalone server`} className="mb-0" />
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
                  runs as a separate process on port 4000 since there is no HTTP server
                  embedded in the Wails app.
                </p>
              </div>
              <CodeBlock terminal code="grit studio" className="mb-0 glow-purple-sm" />
              <div className="prose-grit mt-4">
                <p>
                  Open <code>http://localhost:4000</code> in your browser to browse tables,
                  run queries, and inspect your SQLite database visually.
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
                    url: 'http://localhost:4000',
                    desc: 'Visual database browser',
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
