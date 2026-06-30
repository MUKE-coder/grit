import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/architecture-modes/triple')

export default function TripleArchitecturePage() {
  return (
    <div className="min-h-screen bg-background isolate">
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
              Triple Architecture: Web + Admin + API
            </h1>
            <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl">
              The default and most feature-rich architecture mode. A Turborepo monorepo
              with three applications (web, admin, api) sharing types and validation
              through a shared package.
            </p>
          </div>

          {/* ── Section 1: Overview ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              Overview
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Triple is the default architecture when you run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new</code>.
              It scaffolds a Turborepo monorepo containing three applications that communicate through a shared
              TypeScript package and a REST API. The Go backend serves as the single source of truth for data
              and authentication, while the two frontends serve distinct audiences: the <strong>web</strong> app
              faces your end users, and the <strong>admin</strong> panel provides a Filament-like dashboard
              for managing resources, viewing analytics, and configuring the system.
            </p>
            <CodeBlock language="bash" filename="scaffold command" code={`# Default — triple is implied
grit new myapp --next

# Explicit triple flag
grit new myapp --triple --next

# Triple with TanStack Router instead of Next.js
grit new myapp --triple --vite`} />

            <div className="mt-6 rounded-lg border border-border/40 bg-accent/20 p-6">
              <h3 className="text-sm font-semibold text-foreground mb-3">Best suited for</h3>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li className="flex items-start gap-2"><span className="text-primary mt-0.5">-</span> SaaS products with distinct admin operations (user management, billing, analytics)</li>
                <li className="flex items-start gap-2"><span className="text-primary mt-0.5">-</span> Marketplaces where sellers/admins need a separate interface from buyers</li>
                <li className="flex items-start gap-2"><span className="text-primary mt-0.5">-</span> Content platforms (CMS, LMS, blog networks) where editors manage content separately</li>
                <li className="flex items-start gap-2"><span className="text-primary mt-0.5">-</span> Multi-tenant applications requiring tenant administration</li>
                <li className="flex items-start gap-2"><span className="text-primary mt-0.5">-</span> Any project that benefits from a dedicated admin panel with DataTable, FormBuilder, and dashboard widgets</li>
              </ul>
            </div>
          </section>

          {/* ── Section 2: Complete Folder Structure ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              Complete Folder Structure
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              This is the full, unabbreviated tree that <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new myapp --triple --next</code> generates.
              Every file and directory is listed with its purpose.
            </p>
            <CodeBlock language="bash" filename="myapp/" code={`myapp/
├── .env                              # Shared environment variables
├── .env.example                      # Template for other developers
├── .gitignore
├── .prettierrc                       # Code formatting
├── .prettierignore
├── docker-compose.yml                # Dev: PostgreSQL, Redis, MinIO, Mailhog
├── docker-compose.prod.yml           # Production deployment
├── grit.json                         # Project manifest (architecture, frontend)
├── turbo.json                        # Turborepo pipeline config
├── pnpm-workspace.yaml               # pnpm workspace definition
├── package.json                      # Root scripts (dev, build, lint)
├── .claude/
│   └── skills/grit/
│       ├── SKILL.md                  # AI assistant guide (tailored to triple)
│       └── reference.md              # Detailed API conventions
├── packages/
│   └── shared/                       # Shared between all frontend apps
│       ├── package.json
│       ├── schemas/                  # Zod validation schemas
│       │   ├── index.ts              # // grit:schemas marker
│       │   └── user.ts              # User schema (register, login, update)
│       ├── types/                    # TypeScript interfaces
│       │   ├── index.ts              # // grit:types marker
│       │   └── user.ts              # User type
│       └── constants/                # Shared constants
│           └── index.ts              # API_ROUTES, ROLES // grit:api-routes
├── apps/
│   ├── api/                          # Go backend
│   │   ├── Dockerfile                # Multi-stage build (golang → alpine)
│   │   ├── go.mod                    # Go module: myapp/apps/api
│   │   ├── go.sum
│   │   ├── cmd/
│   │   │   ├── server/main.go        # Entry point: config, db, services, router
│   │   │   ├── migrate/main.go       # Migration runner
│   │   │   └── seed/main.go          # Database seeder
│   │   └── internal/
│   │       ├── config/config.go      # Loads .env, all env vars
│   │       ├── database/db.go        # GORM connection + AutoMigrate
│   │       ├── models/               # GORM models
│   │       │   ├── user.go           # User model // grit:models marker
│   │       │   └── upload.go         # Upload model
│   │       ├── handlers/             # HTTP handlers (thin — call services)
│   │       │   ├── auth.go           # Register, login, refresh, me
│   │       │   ├── user.go           # User CRUD
│   │       │   ├── upload.go         # File upload (presigned URLs)
│   │       │   └── ai.go             # AI completions + streaming
│   │       ├── services/             # Business logic
│   │       │   ├── auth_service.go   # JWT, bcrypt, token generation
│   │       │   └── user_service.go   # User queries
│   │       ├── middleware/           # Gin middleware
│   │       │   ├── auth.go           # JWT verification
│   │       │   ├── cors.go           # CORS configuration
│   │       │   ├── logger.go         # Structured logging
│   │       │   ├── cache.go          # Redis response caching
│   │       │   ├── maintenance.go    # grit down/up support
│   │       │   └── rate_limit.go     # Sentinel rate limiting
│   │       ├── routes/
│   │       │   └── routes.go         # Route registration // grit:handlers, grit:routes:*
│   │       ├── mail/                 # Email service (Resend)
│   │       │   ├── mailer.go         # Send function
│   │       │   └── templates/        # HTML email templates
│   │       ├── storage/              # S3-compatible file storage
│   │       ├── jobs/                 # Background jobs (asynq)
│   │       ├── cron/                 # Scheduled tasks
│   │       ├── cache/                # Redis cache service
│   │       ├── ai/                   # AI service (Vercel AI Gateway)
│   │       └── auth/                 # TOTP 2FA service
│   │           └── totp.go           # Setup, verify, backup codes, trusted devices
│   ├── web/                          # Public-facing frontend
│   │   ├── Dockerfile                # Next.js standalone build
│   │   ├── package.json              # Dependencies + scripts
│   │   ├── next.config.js            # (or vite.config.ts for --vite)
│   │   ├── tailwind.config.ts
│   │   └── app/                      # Next.js App Router (or src/routes/ for --vite)
│   │       ├── layout.tsx            # Root layout
│   │       ├── page.tsx              # Landing page
│   │       ├── (auth)/               # Auth pages (login, register)
│   │       └── (app)/                # Protected app pages
│   └── admin/                        # Admin panel
│       ├── Dockerfile
│       ├── package.json
│       └── app/                      # (or src/routes/)
│           ├── layout.tsx            # Admin layout (sidebar, navbar)
│           ├── page.tsx              # Dashboard
│           ├── resources/            # Resource definitions // grit:resources
│           └── (dashboard)/          # Admin pages (users, system)
└── packages/
    └── grit-ui/                      # 100 shadcn-compatible components
        ├── registry.json
        └── registry/                 # Per-component JSON + TSX`} />
          </section>

          {/* ── Section 3: Directory Explanations ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              Directory-by-Directory Breakdown
            </h2>

            {/* Root */}
            <div className="mb-8">
              <h3 className="text-lg font-semibold text-foreground mb-3">Root Files</h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The monorepo root contains configuration files that apply to all applications. The
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit.json</code> manifest
                tells the CLI which architecture mode and frontend framework the project uses, so commands like
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate</code> know
                exactly which files to create.
              </p>
              <div className="rounded-lg border border-border/40 bg-accent/20 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30">
                      <th className="text-left p-3 text-muted-foreground font-medium">File</th>
                      <th className="text-left p-3 text-muted-foreground font-medium">Purpose</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">.env</td><td className="p-3">Shared environment variables: DB connection, JWT secret, Redis URL, S3 credentials, Resend API key</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">grit.json</td><td className="p-3">Project manifest read by the CLI. Contains architecture mode (<code className="text-xs font-mono bg-accent/30 px-1 rounded">triple</code>), frontend framework, and project name</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">turbo.json</td><td className="p-3">Defines Turborepo pipelines: <code className="text-xs font-mono bg-accent/30 px-1 rounded">dev</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">build</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">lint</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">test</code>. Handles dependency ordering (shared builds before web/admin)</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">pnpm-workspace.yaml</td><td className="p-3">Declares <code className="text-xs font-mono bg-accent/30 px-1 rounded">apps/*</code> and <code className="text-xs font-mono bg-accent/30 px-1 rounded">packages/*</code> as workspace members</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">docker-compose.yml</td><td className="p-3">Development infrastructure: PostgreSQL 16, Redis 7, MinIO (S3-compatible storage), Mailhog (email testing)</td></tr>
                    <tr><td className="p-3 font-mono text-xs">docker-compose.prod.yml</td><td className="p-3">Production deployment: multi-stage Docker builds for API, web, and admin with Nginx reverse proxy</td></tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* packages/shared */}
            <div className="mb-8">
              <h3 className="text-lg font-semibold text-foreground mb-3">packages/shared/</h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The shared package is the bridge between frontend applications. Both <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">web</code> and
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">admin</code> import from it to ensure consistent validation and types.
                When you run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate</code>, new schemas and types
                are injected here automatically via marker comments.
              </p>
              <div className="rounded-lg border border-border/40 bg-accent/20 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30">
                      <th className="text-left p-3 text-muted-foreground font-medium">Directory</th>
                      <th className="text-left p-3 text-muted-foreground font-medium">Contents</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">schemas/</td><td className="p-3">Zod validation schemas. Each resource gets a file (e.g., <code className="text-xs font-mono bg-accent/30 px-1 rounded">product.ts</code>) with <code className="text-xs font-mono bg-accent/30 px-1 rounded">CreateProductSchema</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">UpdateProductSchema</code>. Re-exported from <code className="text-xs font-mono bg-accent/30 px-1 rounded">index.ts</code> via the <code className="text-xs font-mono bg-accent/30 px-1 rounded">// grit:schemas</code> marker.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">types/</td><td className="p-3">TypeScript interfaces inferred from Zod schemas. Each resource file exports the interface (e.g., <code className="text-xs font-mono bg-accent/30 px-1 rounded">Product</code>). Re-exported via the <code className="text-xs font-mono bg-accent/30 px-1 rounded">// grit:types</code> marker.</td></tr>
                    <tr><td className="p-3 font-mono text-xs">constants/</td><td className="p-3">Shared constants including <code className="text-xs font-mono bg-accent/30 px-1 rounded">API_ROUTES</code> (injected via <code className="text-xs font-mono bg-accent/30 px-1 rounded">// grit:api-routes</code>) and <code className="text-xs font-mono bg-accent/30 px-1 rounded">ROLES</code> enum.</td></tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* apps/api */}
            <div className="mb-8">
              <h3 className="text-lg font-semibold text-foreground mb-3">apps/api/ (Go Backend)</h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The Go API follows a strict handler-service-model pattern. Handlers are thin HTTP layers
                that parse requests and call services. Services contain business logic and GORM queries.
                Models define the database schema. This separation makes the codebase testable and maintainable.
              </p>
              <div className="rounded-lg border border-border/40 bg-accent/20 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30">
                      <th className="text-left p-3 text-muted-foreground font-medium">Directory</th>
                      <th className="text-left p-3 text-muted-foreground font-medium">Role</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">cmd/server/</td><td className="p-3">Application entry point. Loads config, connects to database, initializes all services, registers routes, and starts the Gin server on port 8080.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">cmd/migrate/</td><td className="p-3">Standalone migration runner. Calls GORM AutoMigrate on all registered models.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">cmd/seed/</td><td className="p-3">Database seeder. Creates admin user, sample data, and UI component registry seed.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/config/</td><td className="p-3">Loads all environment variables from <code className="text-xs font-mono bg-accent/30 px-1 rounded">.env</code> into a typed <code className="text-xs font-mono bg-accent/30 px-1 rounded">Config</code> struct. Accessed throughout the application.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/database/</td><td className="p-3">GORM connection setup. Opens PostgreSQL (production) or SQLite (development/testing). Runs AutoMigrate for all models.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/models/</td><td className="p-3">GORM model definitions. Each model is a Go struct with <code className="text-xs font-mono bg-accent/30 px-1 rounded">gorm</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">json</code>, and <code className="text-xs font-mono bg-accent/30 px-1 rounded">binding</code> struct tags. New models are injected at the <code className="text-xs font-mono bg-accent/30 px-1 rounded">// grit:models</code> marker.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/handlers/</td><td className="p-3">Gin HTTP handlers. Each resource gets a handler file with Create, GetAll, GetByID, Update, Delete functions. Handlers parse input, call the corresponding service, and return JSON responses.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/services/</td><td className="p-3">Business logic layer. Each service receives a <code className="text-xs font-mono bg-accent/30 px-1 rounded">*gorm.DB</code> and implements CRUD operations with pagination, filtering, and error handling.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/middleware/</td><td className="p-3">Gin middleware stack: JWT auth verification, CORS policy, structured request logging, Redis response caching, rate limiting (Sentinel), and maintenance mode toggle.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/routes/</td><td className="p-3">Centralized route registration. Contains marker comments (<code className="text-xs font-mono bg-accent/30 px-1 rounded">// grit:handlers</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">// grit:routes:public</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">// grit:routes:protected</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">// grit:routes:admin</code>) where generated routes are injected.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/mail/</td><td className="p-3">Email sending via Resend API. Includes HTML templates for welcome, password reset, verification, and notification emails.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/storage/</td><td className="p-3">S3-compatible file storage. Generates presigned upload URLs for direct browser-to-S3 uploads. Works with MinIO (dev), Cloudflare R2, or AWS S3 (production).</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/jobs/</td><td className="p-3">Background job processing with asynq + Redis. Includes workers for email delivery, image processing, and cleanup tasks.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/cron/</td><td className="p-3">Scheduled tasks with asynq cron. Define recurring jobs with standard cron expressions.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/cache/</td><td className="p-3">Redis cache service. Set/get/delete with TTL. Used by the cache middleware to cache entire API responses.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/ai/</td><td className="p-3">AI service supporting Claude and OpenAI. Includes streaming responses for real-time completions.</td></tr>
                    <tr><td className="p-3 font-mono text-xs">internal/auth/</td><td className="p-3">TOTP two-factor authentication. Setup flow, verification, backup code generation, and trusted device management.</td></tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* apps/web */}
            <div className="mb-8">
              <h3 className="text-lg font-semibold text-foreground mb-3">apps/web/ (Public Frontend)</h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The web application is what your end users interact with. It consumes the Go API via React Query
                hooks and validates form inputs using the Zod schemas from the shared package. With Next.js,
                pages can be server-rendered for SEO. With TanStack Router (Vite), the app is a client-side SPA
                served as static files.
              </p>
              <div className="rounded-lg border border-border/40 bg-accent/20 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30">
                      <th className="text-left p-3 text-muted-foreground font-medium">Path</th>
                      <th className="text-left p-3 text-muted-foreground font-medium">Description</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">app/layout.tsx</td><td className="p-3">Root layout: global providers (React Query, theme), font loading, metadata</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">app/page.tsx</td><td className="p-3">Landing page — the first page visitors see</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">app/(auth)/</td><td className="p-3">Auth route group: login, register, forgot-password pages with shared auth layout</td></tr>
                    <tr><td className="p-3 font-mono text-xs">app/(app)/</td><td className="p-3">Protected route group: dashboard, settings, and resource pages. Requires JWT token.</td></tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* apps/admin */}
            <div className="mb-8">
              <h3 className="text-lg font-semibold text-foreground mb-3">apps/admin/ (Admin Panel)</h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The admin panel is a Filament-inspired dashboard with runtime resource definitions. It provides
                DataTable (sorting, filtering, pagination, row selection), FormBuilder (10+ field types),
                multi-step forms, dashboard widgets, and system administration pages. Resources are defined
                declaratively using <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">defineResource()</code> and
                the admin automatically generates CRUD pages from those definitions.
              </p>
              <div className="rounded-lg border border-border/40 bg-accent/20 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30">
                      <th className="text-left p-3 text-muted-foreground font-medium">Path</th>
                      <th className="text-left p-3 text-muted-foreground font-medium">Description</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">app/layout.tsx</td><td className="p-3">Admin layout: collapsible sidebar with Lucide icons, top navbar, dark/light theme toggle</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">app/page.tsx</td><td className="p-3">Dashboard: StatsCards, ChartWidgets (Recharts), ActivityWidget, WidgetGrid</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">app/resources/</td><td className="p-3">Resource definition files. Each file exports a <code className="text-xs font-mono bg-accent/30 px-1 rounded">defineResource()</code> call that configures columns, filters, search, form fields, and permissions. The <code className="text-xs font-mono bg-accent/30 px-1 rounded">// grit:resources</code> marker is used for injection.</td></tr>
                    <tr><td className="p-3 font-mono text-xs">app/(dashboard)/</td><td className="p-3">Admin pages: users management, system pages (Jobs, Files, Cron, Mail Preview), and generated resource pages</td></tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* packages/grit-ui */}
            <div className="mb-8">
              <h3 className="text-lg font-semibold text-foreground mb-3">packages/grit-ui/ (Component Registry)</h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                A registry of 100 shadcn-compatible components organized into five categories:
                marketing (21), auth (10), saas (30), ecommerce (20), and layout (20). Each component
                has a JSON metadata file and TSX source file. The API serves them
                at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">GET /r.json</code> (full registry)
                and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">GET /r/:name.json</code> (individual component),
                making them installable via the shadcn CLI pattern.
              </p>
            </div>

            {/* .claude */}
            <div className="mb-8">
              <h3 className="text-lg font-semibold text-foreground mb-3">.claude/ (AI Skill)</h3>
              <p className="text-muted-foreground leading-relaxed">
                Every scaffolded project includes an AI skill file tailored to its architecture mode.
                The skill teaches Claude Code (or any LLM) the project conventions, folder structure,
                naming rules, available CLI commands, and code markers. This means AI assistants can
                generate correct Grit code without additional context.
              </p>
            </div>
          </section>

          {/* ── Section 4: Data Flow ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              Data Flow
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              The following diagram shows how a typical request flows through the triple architecture,
              from the user&apos;s browser all the way to the database and back.
            </p>
            <CodeBlock language="bash" filename="request flow" code={`┌─────────────────────────────────────────────────────────┐
│                       BROWSER                           │
│  ┌──────────────┐   ┌──────────────┐                    │
│  │  Web App     │   │ Admin Panel  │                    │
│  │  :3000       │   │ :3001        │                    │
│  └──────┬───────┘   └──────┬───────┘                    │
└─────────┼──────────────────┼────────────────────────────┘
          │                  │
     REST + JWT         REST + JWT
          │                  │
          ▼                  ▼
┌─────────────────────────────────────────────────────────┐
│                    GO API (:8080)                        │
│                                                         │
│  Request → Gin Router                                   │
│         → Middleware Stack                               │
│           ├── CORS (allow web:3000, admin:3001)         │
│           ├── Logger (structured request logging)       │
│           ├── Rate Limiter (Sentinel)                   │
│           ├── Auth (JWT verification, extract user)     │
│           └── Cache (Redis response cache)              │
│         → Handler (parse input, validate)               │
│         → Service (business logic, GORM queries)        │
│         → GORM                                          │
│         → PostgreSQL                                    │
│                                                         │
│  Background:                                            │
│  ├── asynq workers (email, image processing, cleanup)   │
│  ├── cron scheduler (recurring tasks)                   │
│  └── GORM Studio (/studio — visual DB browser)          │
└─────────────────────────────────────────────────────────┘
          │         │         │         │
          ▼         ▼         ▼         ▼
     PostgreSQL   Redis     MinIO    Resend
      (data)     (cache    (files)   (email)
                  + jobs)`} />

            <div className="mt-6 space-y-4">
              <div className="rounded-lg border border-border/40 bg-accent/20 p-4">
                <h4 className="text-sm font-semibold text-foreground mb-2">Step-by-step for a typical API call</h4>
                <ol className="space-y-2 text-sm text-muted-foreground list-decimal list-inside">
                  <li>User clicks &quot;Save&quot; in the admin panel form (or the web app)</li>
                  <li>React Query mutation fires a POST request to <code className="text-xs font-mono bg-accent/30 px-1 rounded">/api/v3/products</code> with the JWT in the Authorization header</li>
                  <li>Gin router matches the route and passes the request through the middleware stack</li>
                  <li>CORS middleware validates the origin. Rate limiter checks the IP. Auth middleware extracts and verifies the JWT, injecting the user into the Gin context.</li>
                  <li>The <code className="text-xs font-mono bg-accent/30 px-1 rounded">CreateProduct</code> handler parses the JSON body using Gin&apos;s binding tags for validation</li>
                  <li>Handler calls <code className="text-xs font-mono bg-accent/30 px-1 rounded">productService.Create()</code> which runs the GORM insert query</li>
                  <li>Service returns the created product. Handler responds with <code className="text-xs font-mono bg-accent/30 px-1 rounded">201 Created</code> and the standard JSON envelope</li>
                  <li>React Query invalidates the products list query, triggering a background refetch of the table data</li>
                </ol>
              </div>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-4">
                <h4 className="text-sm font-semibold text-foreground mb-2">Shared types keep everything in sync</h4>
                <p className="text-sm text-muted-foreground">
                  The Go model struct tags define the API contract. When you run <code className="text-xs font-mono bg-accent/30 px-1 rounded">grit generate</code>,
                  matching Zod schemas and TypeScript types are created in <code className="text-xs font-mono bg-accent/30 px-1 rounded">packages/shared/</code>.
                  Both the web and admin apps import from the same package, so a type change in Go propagates to both frontends
                  via <code className="text-xs font-mono bg-accent/30 px-1 rounded">grit sync</code>.
                </p>
              </div>
            </div>
          </section>

          {/* ── Section 5: Code Generation ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              How Code Generation Works
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              When you run a generate command, Grit creates files across all three applications and the shared
              package, then injects import statements and route registrations into existing files using marker comments.
            </p>
            <CodeBlock language="bash" filename="example command" code={`grit generate resource Product --fields "name:string, price:float, description:text, category_id:belongs_to, is_active:bool"`} />

            <div className="mt-6 space-y-4">
              <h3 className="text-lg font-semibold text-foreground">Files Created</h3>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-5">
                <h4 className="text-sm font-mono text-primary mb-3">Go Backend (apps/api/)</h4>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> <code className="text-xs font-mono bg-accent/30 px-1 rounded">internal/models/product.go</code> — GORM model with struct tags, belongs_to relationship to Category</li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> <code className="text-xs font-mono bg-accent/30 px-1 rounded">internal/services/product_service.go</code> — CRUD operations with pagination, preloaded relationships</li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> <code className="text-xs font-mono bg-accent/30 px-1 rounded">internal/handlers/product.go</code> — Gin handlers: Create, GetAll, GetByID, Update, Delete</li>
                </ul>
              </div>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-5">
                <h4 className="text-sm font-mono text-primary mb-3">Shared Package (packages/shared/)</h4>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> <code className="text-xs font-mono bg-accent/30 px-1 rounded">schemas/product.ts</code> — Zod schemas: <code className="text-xs font-mono bg-accent/30 px-1 rounded">CreateProductSchema</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">UpdateProductSchema</code></li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> <code className="text-xs font-mono bg-accent/30 px-1 rounded">types/product.ts</code> — TypeScript interface: <code className="text-xs font-mono bg-accent/30 px-1 rounded">Product</code></li>
                </ul>
              </div>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-5">
                <h4 className="text-sm font-mono text-primary mb-3">Admin Panel (apps/admin/)</h4>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> <code className="text-xs font-mono bg-accent/30 px-1 rounded">resources/products.ts</code> — Resource definition with columns, filters, search config, form fields</li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> <code className="text-xs font-mono bg-accent/30 px-1 rounded">app/(dashboard)/products/page.tsx</code> — Admin page with DataTable and FormBuilder, ready to use</li>
                </ul>
              </div>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-5">
                <h4 className="text-sm font-mono text-primary mb-3">Web App (apps/web/)</h4>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> <code className="text-xs font-mono bg-accent/30 px-1 rounded">hooks/use-products.ts</code> — React Query hooks: <code className="text-xs font-mono bg-accent/30 px-1 rounded">useProducts()</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">useProduct(id)</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">useCreateProduct()</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">useUpdateProduct()</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">useDeleteProduct()</code></li>
                </ul>
              </div>
            </div>

            <div className="mt-6">
              <h3 className="text-lg font-semibold text-foreground mb-3">Marker Injections</h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                In addition to creating new files, the generator injects code into 6 existing files
                using marker comments. This is how routes get registered, imports get added, and types
                get re-exported — without manual wiring.
              </p>
              <div className="rounded-lg border border-border/40 bg-accent/20 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30">
                      <th className="text-left p-3 text-muted-foreground font-medium">File</th>
                      <th className="text-left p-3 text-muted-foreground font-medium">Marker</th>
                      <th className="text-left p-3 text-muted-foreground font-medium">What Gets Injected</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">routes/routes.go</td><td className="p-3 font-mono text-xs">// grit:handlers</td><td className="p-3">Handler initialization</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">routes/routes.go</td><td className="p-3 font-mono text-xs">// grit:routes:*</td><td className="p-3">Route group registration (public, protected, or admin)</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">schemas/index.ts</td><td className="p-3 font-mono text-xs">// grit:schemas</td><td className="p-3">Schema re-export line</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">types/index.ts</td><td className="p-3 font-mono text-xs">// grit:types</td><td className="p-3">Type re-export line</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">constants/index.ts</td><td className="p-3 font-mono text-xs">// grit:api-routes</td><td className="p-3">API route constant</td></tr>
                    <tr><td className="p-3 font-mono text-xs">resources/ (admin)</td><td className="p-3 font-mono text-xs">// grit:resources</td><td className="p-3">Resource registration in the admin registry</td></tr>
                  </tbody>
                </table>
              </div>
            </div>
          </section>

          {/* ── Section 6: Ports & URLs ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              Ports & URLs
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              During development, all services run on different ports. The CORS middleware in the Go API
              is configured to allow requests from all frontend origins.
            </p>
            <div className="rounded-lg border border-border/40 bg-accent/20 overflow-hidden">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/30">
                    <th className="text-left p-3 text-muted-foreground font-medium">Service</th>
                    <th className="text-left p-3 text-muted-foreground font-medium">Port</th>
                    <th className="text-left p-3 text-muted-foreground font-medium">URL</th>
                    <th className="text-left p-3 text-muted-foreground font-medium">Notes</th>
                  </tr>
                </thead>
                <tbody className="text-muted-foreground">
                  <tr className="border-b border-border/20"><td className="p-3 font-semibold text-foreground">Go API</td><td className="p-3 font-mono text-xs">8080</td><td className="p-3 font-mono text-xs">http://localhost:8080</td><td className="p-3">REST endpoints + GORM Studio at /studio</td></tr>
                  <tr className="border-b border-border/20"><td className="p-3 font-semibold text-foreground">Web App</td><td className="p-3 font-mono text-xs">3000</td><td className="p-3 font-mono text-xs">http://localhost:3000</td><td className="p-3">Next.js dev server (or Vite :5173)</td></tr>
                  <tr className="border-b border-border/20"><td className="p-3 font-semibold text-foreground">Admin Panel</td><td className="p-3 font-mono text-xs">3001</td><td className="p-3 font-mono text-xs">http://localhost:3001</td><td className="p-3">Next.js dev server (or Vite :5174)</td></tr>
                  <tr className="border-b border-border/20"><td className="p-3 font-semibold text-foreground">Docs</td><td className="p-3 font-mono text-xs">3002</td><td className="p-3 font-mono text-xs">http://localhost:3002</td><td className="p-3">Documentation site (if running locally)</td></tr>
                  <tr className="border-b border-border/20"><td className="p-3 font-semibold text-foreground">PostgreSQL</td><td className="p-3 font-mono text-xs">5432</td><td className="p-3 font-mono text-xs">localhost:5432</td><td className="p-3">Docker container</td></tr>
                  <tr className="border-b border-border/20"><td className="p-3 font-semibold text-foreground">Redis</td><td className="p-3 font-mono text-xs">6379</td><td className="p-3 font-mono text-xs">localhost:6379</td><td className="p-3">Cache + job queue</td></tr>
                  <tr className="border-b border-border/20"><td className="p-3 font-semibold text-foreground">MinIO API</td><td className="p-3 font-mono text-xs">9000</td><td className="p-3 font-mono text-xs">http://localhost:9000</td><td className="p-3">S3-compatible API</td></tr>
                  <tr className="border-b border-border/20"><td className="p-3 font-semibold text-foreground">MinIO Console</td><td className="p-3 font-mono text-xs">9001</td><td className="p-3 font-mono text-xs">http://localhost:9001</td><td className="p-3">Web UI for managing buckets</td></tr>
                  <tr className="border-b border-border/20"><td className="p-3 font-semibold text-foreground">Mailhog SMTP</td><td className="p-3 font-mono text-xs">1025</td><td className="p-3 font-mono text-xs">localhost:1025</td><td className="p-3">Catches outgoing email</td></tr>
                  <tr><td className="p-3 font-semibold text-foreground">Mailhog UI</td><td className="p-3 font-mono text-xs">8025</td><td className="p-3 font-mono text-xs">http://localhost:8025</td><td className="p-3">View caught emails in browser</td></tr>
                </tbody>
              </table>
            </div>
          </section>

          {/* ── Section 7: Deployment ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              Deployment Options
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Triple architecture supports multiple deployment strategies depending on your
              infrastructure preferences and scale requirements.
            </p>

            <div className="space-y-4">
              <div className="rounded-lg border border-border/40 bg-accent/20 p-6">
                <h3 className="text-base font-semibold text-foreground mb-2">Docker Compose (Self-hosted)</h3>
                <p className="text-sm text-muted-foreground mb-3">
                  The scaffolded <code className="text-xs font-mono bg-accent/30 px-1 rounded">docker-compose.prod.yml</code> builds
                  all three apps as optimized Docker images and runs them behind an Nginx reverse proxy with
                  automatic SSL. This is the simplest path to production on any VPS.
                </p>
                <CodeBlock language="bash" code={`# Build and deploy everything
docker compose -f docker-compose.prod.yml up -d --build`} />
              </div>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-6">
                <h3 className="text-base font-semibold text-foreground mb-2">grit deploy</h3>
                <p className="text-sm text-muted-foreground mb-3">
                  The <code className="text-xs font-mono bg-accent/30 px-1 rounded">grit deploy</code> command
                  automates the deployment process. It builds Docker images, pushes them to your registry,
                  and deploys to your configured target (VPS, Dokploy, or cloud).
                </p>
                <CodeBlock language="bash" code={`# Deploy to configured target
grit deploy`} />
              </div>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-6">
                <h3 className="text-base font-semibold text-foreground mb-2">Hybrid: Vercel + VPS</h3>
                <p className="text-sm text-muted-foreground">
                  Deploy the Next.js frontends (web + admin) to Vercel for zero-config hosting with edge
                  caching, and deploy the Go API to a VPS or cloud provider (Railway, Fly.io, Render).
                  This gives you Vercel&apos;s CDN for the frontends and full control over the API infrastructure.
                  Set the <code className="text-xs font-mono bg-accent/30 px-1 rounded">NEXT_PUBLIC_API_URL</code> environment
                  variable in each Vercel project to point to your API server.
                </p>
              </div>
            </div>
          </section>

          {/* ── Section 8: When to Choose Triple ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              When to Choose Triple
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Triple is the right choice when your project has distinct user roles with different interfaces.
              Here are the decision criteria:
            </p>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="rounded-lg border border-primary/20 bg-primary/5 p-6">
                <h3 className="text-sm font-semibold text-foreground mb-3">Choose Triple when...</h3>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> You need a dedicated admin panel separate from the user-facing app</li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> Your admin needs DataTable, FormBuilder, dashboard widgets</li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> You have non-technical admins who need a polished management UI</li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> You want <code className="text-xs font-mono bg-accent/30 px-1 rounded">grit generate</code> to create admin pages automatically</li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> You are building a SaaS, marketplace, CMS, or LMS</li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> Your team has separate frontend developers for admin and public UI</li>
                </ul>
              </div>
              <div className="rounded-lg border border-border/40 bg-accent/20 p-6">
                <h3 className="text-sm font-semibold text-foreground mb-3">Consider a simpler mode when...</h3>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li className="flex items-start gap-2"><span className="text-muted-foreground">-</span> Admins and users share the same interface (use Double)</li>
                  <li className="flex items-start gap-2"><span className="text-muted-foreground">-</span> You need a single binary deployment (use Single)</li>
                  <li className="flex items-start gap-2"><span className="text-muted-foreground">-</span> You only need an API with no frontend (use API Only)</li>
                  <li className="flex items-start gap-2"><span className="text-muted-foreground">-</span> You are building a mobile app (use Mobile)</li>
                  <li className="flex items-start gap-2"><span className="text-muted-foreground">-</span> Your project is a simple blog or portfolio</li>
                </ul>
              </div>
            </div>
          </section>

          {/* ── Section 9: Example Project ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              Example Project
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              We built a complete Job Portal using the triple architecture to demonstrate every feature.
              The example includes a public job board (web), an employer/admin dashboard (admin), and a Go API
              with authentication, file uploads, and background job processing.
            </p>
            <div className="rounded-lg border border-border/40 bg-accent/20 p-6">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-base font-semibold text-foreground mb-1">Job Portal — Triple + Next.js</h3>
                  <p className="text-sm text-muted-foreground">
                    Full source code with README, .env template, step-by-step guide, and production Docker Compose.
                  </p>
                </div>
                <a
                  href="https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-triple-next"
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:text-primary/80 transition-colors whitespace-nowrap"
                >
                  View on GitHub &rarr;
                </a>
              </div>
            </div>
          </section>

          {/* ── Navigation ── */}
          <div className="flex items-center justify-between pt-8 border-t border-border/40">
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/concepts/architecture-modes" className="gap-1.5">
                <ArrowLeft className="h-3.5 w-3.5" />
                Architecture Modes
              </Link>
            </Button>
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/concepts/architecture-modes/double" className="gap-1.5">
                Double Architecture
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
