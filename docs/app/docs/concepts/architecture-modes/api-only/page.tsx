import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/architecture-modes/api-only')

export default function ApiOnlyArchitecturePage() {
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
              API Only Architecture: Headless Go Backend
            </h1>
            <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl">
              The leanest architecture Grit offers. Pure Go API with no frontend,
              no React, no Node.js. Your API docs page is the primary interface.
            </p>
          </div>

          <div className="prose-grit">
            {/* Overview */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Overview
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The API-only architecture strips away everything except the Go backend. There is no
                React, no TypeScript, no Node.js, no pnpm, no Turborepo. You get a clean Go API
                that serves JSON endpoints and auto-generated API documentation at{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/docs</code>.
                Test your endpoints with curl, Postman, or the built-in Scalar/Swagger UI.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Despite being the smallest architecture, it still includes every backend battery:
                authentication (JWT + OAuth), file storage (S3), email (Resend), background jobs (asynq),
                cron scheduling, Redis caching, AI integration, and TOTP two-factor auth. You get the
                full Go API -- you just bring your own frontend.
              </p>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-5 mb-6">
                <h4 className="text-sm font-semibold text-foreground mb-3">Scaffold command</h4>
                <CodeBlock language="bash" code={`grit new myapp --api`} />
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
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">API Only Architecture</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Frontend</td>
                      <td className="px-4 py-2.5">None -- pure Go, no React, no Node.js</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Monorepo</td>
                      <td className="px-4 py-2.5">None -- no turbo.json, no pnpm-workspace</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Interface</td>
                      <td className="px-4 py-2.5">API docs at /docs (Scalar/Swagger UI)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Testing</td>
                      <td className="px-4 py-2.5">curl, Postman, or built-in API docs</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Batteries</td>
                      <td className="px-4 py-2.5">All included (auth, storage, email, jobs, AI, TOTP)</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">Code generation</td>
                      <td className="px-4 py-2.5">Go files only (model, service, handler + route injection)</td>
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
                The project contains just the Go API inside{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/api/</code>{' '}
                plus root-level config files. No frontend directories, no packages directory,
                no TypeScript configuration.
              </p>
              <CodeBlock language="bash" filename="myapp/" code={`myapp/
\u251c\u2500\u2500 apps/
\u2502   \u2514\u2500\u2500 api/
\u2502       \u251c\u2500\u2500 Dockerfile
\u2502       \u251c\u2500\u2500 go.mod                    # Module: myapp/apps/api
\u2502       \u251c\u2500\u2500 go.sum
\u2502       \u251c\u2500\u2500 cmd/
\u2502       \u2502   \u251c\u2500\u2500 server/main.go        # API entry point
\u2502       \u2502   \u251c\u2500\u2500 migrate/main.go       # Run migrations
\u2502       \u2502   \u2514\u2500\u2500 seed/main.go          # Seed database
\u2502       \u2514\u2500\u2500 internal/                 # All Go backend code
\u2502           \u251c\u2500\u2500 config/config.go
\u2502           \u251c\u2500\u2500 database/db.go
\u2502           \u251c\u2500\u2500 models/               # // grit:models
\u2502           \u2502   \u251c\u2500\u2500 user.go
\u2502           \u2502   \u2514\u2500\u2500 upload.go
\u2502           \u251c\u2500\u2500 handlers/
\u2502           \u2502   \u251c\u2500\u2500 auth_handler.go
\u2502           \u2502   \u251c\u2500\u2500 upload_handler.go
\u2502           \u2502   \u2514\u2500\u2500 user_handler.go
\u2502           \u251c\u2500\u2500 services/
\u2502           \u2502   \u251c\u2500\u2500 auth_service.go
\u2502           \u2502   \u251c\u2500\u2500 upload_service.go
\u2502           \u2502   \u2514\u2500\u2500 user_service.go
\u2502           \u251c\u2500\u2500 middleware/
\u2502           \u2502   \u251c\u2500\u2500 auth.go
\u2502           \u2502   \u251c\u2500\u2500 cors.go
\u2502           \u2502   \u251c\u2500\u2500 logger.go
\u2502           \u2502   \u2514\u2500\u2500 cache.go
\u2502           \u251c\u2500\u2500 routes/routes.go      # // grit:handlers, grit:routes:*
\u2502           \u251c\u2500\u2500 mail/
\u2502           \u251c\u2500\u2500 storage/
\u2502           \u251c\u2500\u2500 jobs/
\u2502           \u251c\u2500\u2500 cache/
\u2502           \u251c\u2500\u2500 ai/
\u2502           \u2514\u2500\u2500 auth/totp.go
\u251c\u2500\u2500 .env
\u251c\u2500\u2500 .env.example
\u251c\u2500\u2500 .gitignore
\u251c\u2500\u2500 docker-compose.yml                # PostgreSQL, Redis, MinIO, Mailhog
\u251c\u2500\u2500 docker-compose.prod.yml
\u251c\u2500\u2500 grit.json                         # architecture: "api"
\u2514\u2500\u2500 .claude/skills/grit/              # Tailored \u2014 no frontend rules
    \u251c\u2500\u2500 SKILL.md
    \u2514\u2500\u2500 reference.md`} />
            </div>

            {/* Directory Explanation */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Directory Breakdown
              </h2>

              <div className="space-y-6">
                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    <code className="text-sm font-mono text-primary">cmd/</code> -- Entry Points
                  </h3>
                  <p className="text-muted-foreground leading-relaxed">
                    Three separate entry points for different operations.{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">cmd/server/main.go</code>{' '}
                    starts the API server.{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">cmd/migrate/main.go</code>{' '}
                    runs GORM AutoMigrate to apply schema changes.{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">cmd/seed/main.go</code>{' '}
                    populates the database with initial data (admin user, sample records).
                  </p>
                </div>

                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    <code className="text-sm font-mono text-primary">internal/</code> -- Backend Code
                  </h3>
                  <p className="text-muted-foreground leading-relaxed">
                    Identical to every other Grit architecture. The handler-service-model pattern,
                    middleware stack, mail templates, storage abstraction, job queue, cache layer,
                    and AI integration are all present. The only difference is that no frontend code
                    exists alongside it.
                  </p>
                </div>

                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    <code className="text-sm font-mono text-primary">.claude/skills/grit/</code> -- AI Skill
                  </h3>
                  <p className="text-muted-foreground leading-relaxed">
                    The AI skill file is tailored for the API-only architecture. It contains no
                    frontend rules, no React patterns, and no TypeScript conventions. This makes
                    AI assistants more effective because they only see relevant Go patterns.
                  </p>
                </div>

                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    <code className="text-sm font-mono text-primary">grit.json</code> -- Project Config
                  </h3>
                  <CodeBlock language="json" filename="grit.json" code={`{
  "architecture": "api",
  "go_module": "myapp/apps/api"
}`} />
                  <p className="text-sm text-muted-foreground/60 mt-3">
                    No <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">frontend</code> key.
                    The CLI knows to skip all frontend file generation when it reads this config.
                  </p>
                </div>
              </div>
            </div>

            {/* What Gets Generated */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Code Generation
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Running{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>{' '}
                in an API-only project creates only Go files. No Zod schemas, no TypeScript types,
                no React hooks, no admin pages.
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
                      <td className="px-4 py-2.5 font-mono text-xs">apps/api/internal/models/post.go</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Go service</td>
                      <td className="px-4 py-2.5 font-mono text-xs">apps/api/internal/services/post_service.go</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Go handler</td>
                      <td className="px-4 py-2.5 font-mono text-xs">apps/api/internal/handlers/post_handler.go</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5">Route injection</td>
                      <td className="px-4 py-2.5 font-mono text-xs">apps/api/internal/routes/routes.go</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mt-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">NOT Generated</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Reason</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Zod schemas</td>
                      <td className="px-4 py-2.5">No frontend to validate</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">TypeScript types</td>
                      <td className="px-4 py-2.5">No TypeScript in project</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">React Query hooks</td>
                      <td className="px-4 py-2.5">No React in project</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5">Admin page</td>
                      <td className="px-4 py-2.5">No admin panel</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <p className="text-sm text-muted-foreground/60 mt-3">
                Example:{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                  grit generate resource Post title:string body:text published:bool
                </code>{' '}
                creates 3 Go files and injects routes -- that&apos;s it.
              </p>
            </div>

            {/* API Endpoints */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Built-in API Endpoints
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Every API-only project ships with these endpoints out of the box, before you
                generate any resources.
              </p>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Method</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Endpoint</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Purpose</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">GET</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/api/health</td>
                      <td className="px-4 py-2.5">Health check</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/api/auth/register</td>
                      <td className="px-4 py-2.5">Create account</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/api/auth/login</td>
                      <td className="px-4 py-2.5">Get JWT tokens</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/api/auth/refresh</td>
                      <td className="px-4 py-2.5">Refresh access token</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">GET</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/api/auth/me</td>
                      <td className="px-4 py-2.5">Current user profile</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/api/uploads</td>
                      <td className="px-4 py-2.5">File upload (presigned URL)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/api/ai/chat</td>
                      <td className="px-4 py-2.5">AI chat (streaming)</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">GET</td>
                      <td className="px-4 py-2.5 font-mono text-xs">/docs</td>
                      <td className="px-4 py-2.5">Interactive API documentation</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* Data Flow */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Data Flow
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                With no frontend, the data flow is straightforward. Clients send HTTP requests
                directly to the Go API.
              </p>
              <CodeBlock language="bash" filename="data flow" code={`Client (curl / Postman / mobile app / external frontend)
    \u2502
    \u251c\u2500\u2500 POST /api/auth/login      \u2192 Auth handler \u2192 JWT service
    \u251c\u2500\u2500 GET  /api/posts            \u2192 Post handler \u2192 Post service \u2192 PostgreSQL
    \u251c\u2500\u2500 POST /api/uploads          \u2192 Upload handler \u2192 S3/MinIO
    \u251c\u2500\u2500 POST /api/ai/chat          \u2192 AI handler \u2192 Claude/OpenAI
    \u2514\u2500\u2500 GET  /docs                 \u2192 Scalar/Swagger UI`} />
            </div>

            {/* Development */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Development Workflow
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Development is simpler than any other architecture -- just one terminal for the Go server.
              </p>
              <CodeBlock language="bash" filename="getting started" code={`# Start infrastructure
docker compose up -d

# Run the API with hot reload
cd apps/api && air

# Or without air:
cd apps/api && go run cmd/server/main.go

# Test with curl
curl http://localhost:8080/api/health
curl -X POST http://localhost:8080/api/auth/register \\
  -H "Content-Type: application/json" \\
  -d '{"name":"admin","email":"admin@example.com","password":"password123"}'

# Or open the API docs in your browser
open http://localhost:8080/docs`} />
            </div>

            {/* Deployment */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Deployment
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Deploy the Go binary with Docker or directly. No frontend build step needed.
              </p>

              <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                Docker
              </h3>
              <CodeBlock language="dockerfile" filename="apps/api/Dockerfile" code={`FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server cmd/server/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]`} />

              <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                Docker Compose (production)
              </h3>
              <CodeBlock language="yaml" filename="docker-compose.prod.yml" code={`services:
  api:
    build:
      context: ./apps/api
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    env_file: .env
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:16-alpine
    volumes:
      - pgdata:/var/lib/postgresql/data
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: \${DB_PASSWORD}

  redis:
    image: redis:7-alpine

volumes:
  pgdata:`} />
            </div>

            {/* When to Choose */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                When to Choose API Only
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="rounded-lg border border-primary/20 bg-primary/5 p-5">
                  <h4 className="text-sm font-semibold text-foreground mb-3">Good fit</h4>
                  <ul className="space-y-2 text-sm text-muted-foreground">
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Mobile app backends (iOS/Android consume your API)
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Microservices in a larger system
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Headless CMS or content API
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Third-party integrations and webhooks
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Frontend built by a different team or in a different repo
                    </li>
                  </ul>
                </div>
                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h4 className="text-sm font-semibold text-foreground mb-3">Not ideal for</h4>
                  <ul className="space-y-2 text-sm text-muted-foreground">
                    <li className="flex items-start gap-2">
                      <span className="text-muted-foreground/50 mt-0.5">-</span>
                      Projects that need an admin panel (use triple or double)
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-muted-foreground/50 mt-0.5">-</span>
                      Full-stack apps where you want React scaffolded too
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-muted-foreground/50 mt-0.5">-</span>
                      Teams that want shared TypeScript types from the same repo
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
                  The Job Portal example built with the API-only architecture. Includes all endpoints,
                  Docker setup, seed data, and API documentation.
                </p>
                <a
                  href="https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-api-only"
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:text-primary/80 transition-colors"
                >
                  View Job Portal (API Only) on GitHub &rarr;
                </a>
              </div>
            </div>
          </div>

          {/* Nav */}
          <div className="flex items-center justify-between pt-8 border-t border-border/40">
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/concepts/architecture-modes/single" className="gap-1.5">
                <ArrowLeft className="h-3.5 w-3.5" />
                Single Architecture
              </Link>
            </Button>
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/concepts/architecture-modes/mobile" className="gap-1.5">
                Mobile Architecture
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
