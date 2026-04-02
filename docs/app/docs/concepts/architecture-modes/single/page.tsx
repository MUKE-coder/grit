import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/architecture-modes/single')

export default function SingleArchitecturePage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="max-w-4xl mx-auto py-12 px-6 lg:px-8">
          {/* Header */}
          <div className="mb-14">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">
              Architecture Modes
            </p>
            <h1 className="text-3xl lg:text-4xl font-bold tracking-tight text-foreground mb-4">
              Single Architecture: One Binary
            </h1>
            <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl">
              One Go binary serves both the API and the frontend. No monorepo tooling,
              no separate deployments. Build it, ship one file, done.
            </p>
          </div>

          <div className="prose-grit">
            {/* Overview */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Overview
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The single architecture is the most different from the other Grit modes. Instead of
                a Turborepo monorepo with separate apps, you get a flat project where a single Go
                binary serves both the REST API and the compiled React frontend using{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">go:embed</code>.
                Think of it like Laravel or an embedded Next.js app -- one deployment unit that handles everything.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                There is no Turborepo, no{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">pnpm-workspace.yaml</code>,
                no{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">turbo.json</code>.
                The Go module path is just{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">myapp</code>{' '}
                (not <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">myapp/apps/api</code>).
                Schemas and types live directly in{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">frontend/src/</code>{' '}
                -- no shared package needed.
              </p>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-5 mb-6">
                <h4 className="text-sm font-semibold text-foreground mb-3">Scaffold command</h4>
                <CodeBlock language="bash" code={`grit new myapp --single --vite`} />
              </div>
            </div>

            {/* Key Characteristics */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Key Characteristics
              </h2>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Property</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Single Architecture</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Binary</td>
                      <td className="px-4 py-2.5">One Go binary serves API + frontend</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Monorepo</td>
                      <td className="px-4 py-2.5">None -- flat project, no Turborepo, no pnpm workspaces</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Module path</td>
                      <td className="px-4 py-2.5"><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">myapp</code></td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Frontend</td>
                      <td className="px-4 py-2.5">React + Vite + TanStack Router (embedded via go:embed)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Shared types</td>
                      <td className="px-4 py-2.5">Inline in frontend/src/ (no packages/shared)</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">Deployment</td>
                      <td className="px-4 py-2.5">Upload one binary + .env file</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* Full Folder Structure */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Full Folder Structure
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The entire project is flat. Go code lives in{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">internal/</code>,
                React code lives in{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">frontend/</code>,
                and{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">main.go</code>{' '}
                sits at the root.
              </p>
              <CodeBlock language="bash" filename="myapp/" code={`myapp/
\u251c\u2500\u2500 main.go                           # Entry point with go:embed frontend/dist/*
\u251c\u2500\u2500 go.mod                            # Module: myapp (not myapp/apps/api)
\u251c\u2500\u2500 go.sum
\u251c\u2500\u2500 .env
\u251c\u2500\u2500 .env.example
\u251c\u2500\u2500 .gitignore
\u251c\u2500\u2500 .prettierrc
\u251c\u2500\u2500 .prettierignore
\u251c\u2500\u2500 docker-compose.yml                # PostgreSQL, Redis, MinIO, Mailhog
\u251c\u2500\u2500 docker-compose.prod.yml
\u251c\u2500\u2500 grit.json                         # architecture: "single", frontend: "tanstack"
\u251c\u2500\u2500 Makefile                          # make dev, make build, make deploy
\u251c\u2500\u2500 .claude/skills/grit/
\u2502   \u251c\u2500\u2500 SKILL.md                      # Tailored to single architecture
\u2502   \u2514\u2500\u2500 reference.md
\u251c\u2500\u2500 internal/                         # ALL Go backend code
\u2502   \u251c\u2500\u2500 config/config.go
\u2502   \u251c\u2500\u2500 database/db.go
\u2502   \u251c\u2500\u2500 models/                       # // grit:models
\u2502   \u251c\u2500\u2500 handlers/
\u2502   \u251c\u2500\u2500 services/
\u2502   \u251c\u2500\u2500 middleware/
\u2502   \u251c\u2500\u2500 routes/routes.go              # // grit:handlers, grit:routes:*
\u2502   \u251c\u2500\u2500 mail/
\u2502   \u251c\u2500\u2500 storage/
\u2502   \u251c\u2500\u2500 jobs/
\u2502   \u251c\u2500\u2500 cache/
\u2502   \u251c\u2500\u2500 ai/
\u2502   \u2514\u2500\u2500 auth/totp.go
\u2514\u2500\u2500 frontend/                         # React + Vite + TanStack Router
    \u251c\u2500\u2500 package.json
    \u251c\u2500\u2500 vite.config.ts                # Proxy /api \u2192 localhost:8080
    \u251c\u2500\u2500 tailwind.config.ts
    \u251c\u2500\u2500 tsconfig.json
    \u251c\u2500\u2500 index.html
    \u2514\u2500\u2500 src/
        \u251c\u2500\u2500 main.tsx
        \u251c\u2500\u2500 routes/                   # TanStack Router file-based routes
        \u2502   \u251c\u2500\u2500 __root.tsx
        \u2502   \u251c\u2500\u2500 index.tsx
        \u2502   \u2514\u2500\u2500 ...
        \u251c\u2500\u2500 components/
        \u251c\u2500\u2500 hooks/
        \u251c\u2500\u2500 lib/
        \u251c\u2500\u2500 schemas/                  # Zod schemas (inline, not shared package)
        \u2514\u2500\u2500 types/                    # TypeScript types (inline)`} />
            </div>

            {/* Directory Explanation */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Directory Breakdown
              </h2>

              <div className="space-y-6">
                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    <code className="text-sm font-mono text-primary">main.go</code> -- Entry Point
                  </h3>
                  <p className="text-muted-foreground leading-relaxed mb-3">
                    The root <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">main.go</code> file
                    is where everything starts. It uses Go&apos;s{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">embed</code> package
                    to bundle the compiled frontend into the Go binary. In production, this means
                    you get a single executable that can serve both the API and the static frontend files.
                  </p>
                  <CodeBlock language="go" filename="main.go" code={`package main

import (
    "embed"
    "io/fs"
    "net/http"

    "myapp/internal/config"
    "myapp/internal/database"
    "myapp/internal/routes"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

func main() {
    cfg := config.Load()
    db := database.Connect(cfg)

    // Serve API routes
    r := routes.Setup(db, cfg)

    // Serve embedded frontend
    dist, _ := fs.Sub(frontendFS, "frontend/dist")
    r.NoRoute(gin.WrapH(http.FileServer(http.FS(dist))))

    r.Run(":" + cfg.Port)
}`} />
                  <p className="text-sm text-muted-foreground/60 mt-3">
                    The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">{`//go:embed frontend/dist/*`}</code> directive
                    tells the Go compiler to include the entire{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">frontend/dist/</code> directory
                    inside the compiled binary. The{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">NoRoute</code> handler
                    serves the frontend for any path that doesn&apos;t match an API route.
                  </p>
                </div>

                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    <code className="text-sm font-mono text-primary">internal/</code> -- Go Backend
                  </h3>
                  <p className="text-muted-foreground leading-relaxed">
                    Identical structure to the other architectures. Contains config, database connection,
                    models, handlers, services, middleware, mail, storage, jobs, cache, and AI integration.
                    The only difference is import paths: you import{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">{`"myapp/internal/models"`}</code>{' '}
                    instead of{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">{`"myapp/apps/api/internal/models"`}</code>.
                    The{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">routes/routes.go</code> file
                    contains the same code markers{' '}
                    (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">{`// grit:handlers`}</code>,{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">{`// grit:routes:*`}</code>)
                    for code generation.
                  </p>
                </div>

                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    <code className="text-sm font-mono text-primary">frontend/</code> -- React SPA
                  </h3>
                  <p className="text-muted-foreground leading-relaxed">
                    A standard React + Vite + TanStack Router application. File-based routing
                    lives in{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">src/routes/</code>.
                    Unlike the triple or double architecture, there is no{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">packages/shared/</code> directory.
                    Zod schemas live in{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">src/schemas/</code>{' '}
                    and TypeScript types in{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">src/types/</code> --
                    everything is self-contained within the frontend directory.
                  </p>
                </div>

                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    <code className="text-sm font-mono text-primary">grit.json</code> -- Project Config
                  </h3>
                  <p className="text-muted-foreground leading-relaxed">
                    Tells the Grit CLI which architecture this project uses. Code generation reads this
                    file to determine where to place generated files and which templates to use.
                  </p>
                  <CodeBlock language="json" filename="grit.json" code={`{
  "architecture": "single",
  "frontend": "tanstack",
  "go_module": "myapp"
}`} />
                </div>
              </div>
            </div>

            {/* Data Flow */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Data Flow
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The data flow differs between development and production. In development, you run
                two processes: Go on port 8080 and Vite on port 5173. Vite proxies API requests
                to Go. In production, there is only one process: the Go binary serves everything.
              </p>

              <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                Production
              </h3>
              <CodeBlock language="bash" filename="production data flow" code={`Browser \u2192 Go binary (:8080)
               \u251c\u2500\u2500 /api/*        \u2192 Gin router \u2192 handlers \u2192 services \u2192 database
               \u2514\u2500\u2500 everything else \u2192 serves frontend/dist/ (HTML, JS, CSS)`} />

              <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                Development
              </h3>
              <CodeBlock language="bash" filename="development data flow" code={`Browser \u2192 Vite dev server (:5173)
               \u251c\u2500\u2500 /api/*  \u2192 proxied to Go (:8080) \u2192 handlers \u2192 services \u2192 database
               \u2514\u2500\u2500 other   \u2192 Vite HMR (instant refresh)`} />
              <p className="text-sm text-muted-foreground/60 mt-3">
                The Vite proxy is configured in{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">vite.config.ts</code>{' '}
                to forward any request starting with{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/api</code>{' '}
                to the Go server running on port 8080.
              </p>
            </div>

            {/* Dev Workflow */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Development Workflow
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                During development, you run two processes in separate terminals. The Makefile
                provides a convenience command that runs both at once.
              </p>

              <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                Terminal 1: Go API
              </h3>
              <CodeBlock language="bash" code={`# Start the Go API with hot reload
air
# or without air:
go run main.go`} />

              <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                Terminal 2: Vite Frontend
              </h3>
              <CodeBlock language="bash" code={`cd frontend && pnpm dev`} />

              <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                Or use the Makefile
              </h3>
              <CodeBlock language="bash" code={`# Runs both Go and Vite concurrently
make dev`} />

              <div className="rounded-lg border border-border/40 bg-accent/20 p-5 mt-6">
                <h4 className="text-sm font-semibold text-foreground mb-2">Vite Proxy Config</h4>
                <p className="text-sm text-muted-foreground mb-3">
                  The Vite config proxies API requests so the frontend can call{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/api/users</code>{' '}
                  without worrying about CORS or port differences during development.
                </p>
                <CodeBlock language="typescript" filename="frontend/vite.config.ts" code={`import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { TanStackRouterVite } from '@tanstack/router-plugin/vite'

export default defineConfig({
  plugins: [TanStackRouterVite(), react()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})`} />
              </div>
            </div>

            {/* Production Build */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Production Build
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Building for production is a two-step process: build the frontend, then build the
                Go binary. The Go compiler embeds the compiled frontend assets into the binary
                at compile time.
              </p>
              <CodeBlock language="bash" filename="build steps" code={`# Step 1: Build the frontend
cd frontend && pnpm build
# This creates frontend/dist/ with optimized HTML, JS, and CSS

# Step 2: Build the Go binary
cd .. && go build -o myapp main.go
# The binary now contains the embedded frontend

# Step 3: Deploy
scp myapp server:/opt/myapp/
scp .env server:/opt/myapp/
# That's it. One binary + one .env file.`} />
              <p className="text-sm text-muted-foreground/60 mt-3">
                No Node.js runtime needed on the production server. No{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">npm install</code>,
                no{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">pnpm</code>,
                no build tools. Just the binary and your environment variables.
              </p>
            </div>

            {/* Code Generation */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Code Generation
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                When you run{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>{' '}
                in a single-architecture project, the CLI creates files in different locations compared
                to the monorepo architectures.
              </p>

              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Generated File</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Location</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Go model</td>
                      <td className="px-4 py-2.5 font-mono text-xs">internal/models/post.go</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Go service</td>
                      <td className="px-4 py-2.5 font-mono text-xs">internal/services/post_service.go</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Go handler</td>
                      <td className="px-4 py-2.5 font-mono text-xs">internal/handlers/post_handler.go</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Route injection</td>
                      <td className="px-4 py-2.5 font-mono text-xs">internal/routes/routes.go</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Zod schema</td>
                      <td className="px-4 py-2.5 font-mono text-xs">frontend/src/schemas/post.ts</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">TypeScript types</td>
                      <td className="px-4 py-2.5 font-mono text-xs">frontend/src/types/post.ts</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5">React Query hooks</td>
                      <td className="px-4 py-2.5 font-mono text-xs">frontend/src/hooks/use-posts.ts</td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <p className="text-sm text-muted-foreground/60 mt-3">
                Notice that schemas and types go into{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">frontend/src/</code>{' '}
                directly, not into a{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">packages/shared/</code>{' '}
                directory. There is no shared package in the single architecture.
              </p>
            </div>

            {/* Deployment */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Deployment
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Deploying a single-architecture app is as simple as it gets. You only need to transfer
                two files to your server: the compiled binary and your{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.env</code> file.
              </p>

              <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                Docker
              </h3>
              <CodeBlock language="dockerfile" filename="Dockerfile" code={`FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm build

FROM golang:1.21-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 go build -o server main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=backend /app/server .
EXPOSE 8080
CMD ["./server"]`} />

              <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                Direct Deploy (no Docker)
              </h3>
              <CodeBlock language="bash" code={`# Build locally for your server's OS/arch
GOOS=linux GOARCH=amd64 go build -o myapp main.go

# Transfer and run
scp myapp .env yourserver:/opt/myapp/
ssh yourserver 'cd /opt/myapp && ./myapp'`} />
            </div>

            {/* When to Choose */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                When to Choose Single
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="rounded-lg border border-primary/20 bg-primary/5 p-5">
                  <h4 className="text-sm font-semibold text-foreground mb-3">Good fit</h4>
                  <ul className="space-y-2 text-sm text-muted-foreground">
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Solo developers who want the simplest deploy story
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Laravel or Rails developers who like one-app-does-everything
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Internal tools and admin dashboards
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Microservices where each service has its own UI
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Apps where SSR and SEO are not critical (SPA is fine)
                    </li>
                  </ul>
                </div>
                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h4 className="text-sm font-semibold text-foreground mb-3">Not ideal for</h4>
                  <ul className="space-y-2 text-sm text-muted-foreground">
                    <li className="flex items-start gap-2">
                      <span className="text-muted-foreground/50 mt-0.5">-</span>
                      SEO-heavy public sites that need server-side rendering
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-muted-foreground/50 mt-0.5">-</span>
                      Large teams that need separate frontend and backend deploys
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-muted-foreground/50 mt-0.5">-</span>
                      Projects that need a separate admin panel (use triple instead)
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-muted-foreground/50 mt-0.5">-</span>
                      Apps that share types with multiple frontend clients
                    </li>
                  </ul>
                </div>
              </div>
            </div>

            {/* Example Project */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Example Project
              </h2>
              <div className="rounded-lg border border-border/40 bg-accent/20 p-5">
                <p className="text-muted-foreground mb-3">
                  The Job Portal example built with the single architecture demonstrates all the patterns
                  described above: go:embed, Vite proxy, TanStack Router, and production Dockerfile.
                </p>
                <a
                  href="https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-single-vite"
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:text-primary/80 transition-colors"
                >
                  View Job Portal (Single + Vite) on GitHub &rarr;
                </a>
              </div>
            </div>
          </div>

          {/* Nav */}
          <div className="flex items-center justify-between pt-8 border-t border-border/40">
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/concepts/architecture-modes/double" className="gap-1.5">
                <ArrowLeft className="h-3.5 w-3.5" />
                Double Architecture
              </Link>
            </Button>
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/concepts/architecture-modes/api-only" className="gap-1.5">
                API Only Architecture
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
