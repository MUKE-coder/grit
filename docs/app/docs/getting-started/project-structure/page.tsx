import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/getting-started/project-structure')

export default function ProjectStructurePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Getting Started</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Project Structure
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                A complete guide to the Grit monorepo layout. Every Grit project follows this
                exact structure, so any developer (or AI assistant) can jump in and know exactly
                where everything lives.
              </p>
            </div>

            <div className="prose-grit">
              <h2>Overview</h2>
              <p>
                Grit uses a Turborepo-powered monorepo with three applications and one shared
                package. The Go API, Next.js web app, Next.js admin panel, and shared TypeScript
                package all live in a single repository with shared configuration.
              </p>
            </div>

            {/* Full Tree */}
            <div className="mb-10">
              <CodeBlock filename="monorepo root" code={`myapp/
├── .env                        # Environment variables
├── .env.example                # Template with documentation
├── .env.cloud.example          # Cloud-only setup (no Docker)
├── .gitignore                  # Git ignore rules
├── docker-compose.yml          # Dev: PostgreSQL, Redis, MinIO, Mailhog
├── docker-compose.prod.yml     # Production: multi-stage builds
├── grit.config.ts              # Grit framework configuration
├── package.json                # Root package.json (workspace scripts)
├── pnpm-workspace.yaml         # pnpm workspace definition
├── turbo.json                  # Turborepo task configuration
├── README.md                   # Project documentation
│
├── apps/
│   ├── api/                    # Go backend (Gin + GORM)
│   ├── web/                    # Next.js main frontend
│   └── admin/                  # Next.js admin panel
│
└── packages/
    └── shared/                 # Shared Zod schemas, TS types, constants`} />
            </div>

            <div className="prose-grit">
              <h2>Root Files</h2>
              <p>
                The root of the monorepo contains configuration files that apply to the entire project:
              </p>
            </div>

            <div className="mb-10 space-y-3">
              {[
                { file: '.env', desc: 'All environment variables for the project -- database URL, JWT secret, Redis, storage, email, AI config. This file is gitignored.' },
                { file: '.env.example', desc: 'Documented template of all environment variables with sensible defaults. Committed to git so new developers know what variables are needed.' },
                { file: 'docker-compose.yml', desc: 'Development services: PostgreSQL 16, Redis 7, MinIO (S3-compatible storage), and Mailhog (email testing). Run with docker compose up -d.' },
                { file: 'docker-compose.prod.yml', desc: 'Production setup with multi-stage Docker builds for the Go API and Next.js apps.' },
                { file: 'grit.config.ts', desc: 'Grit framework configuration -- project name, API URL, and other framework-level settings.' },
                { file: 'turbo.json', desc: 'Turborepo configuration defining build, dev, and lint tasks with dependency relationships and caching.' },
                { file: 'pnpm-workspace.yaml', desc: 'Defines the pnpm workspace: apps/* and packages/* directories are included.' },
                { file: 'package.json', desc: 'Root package.json with workspace-level scripts like dev, build, and lint.' },
              ].map((item) => (
                <div key={item.file} className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                  <code className="text-sm font-mono text-primary/70 font-medium">{item.file}</code>
                  <p className="text-xs text-muted-foreground/60 mt-1 leading-relaxed">{item.desc}</p>
                </div>
              ))}
            </div>

            {/* Go API */}
            <div className="prose-grit">
              <h2>Go API (apps/api/)</h2>
              <p>
                The Go backend is a Gin web server with GORM ORM. It follows Go conventions
                with an <code>internal/</code> directory for private packages and <code>cmd/</code>{' '}
                for the entry point.
              </p>
            </div>

            <div className="mb-4">
              <CodeBlock filename="apps/api/" code={`apps/api/
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
├── Dockerfile                  # Multi-stage production build
├── .air.toml                   # Hot reload configuration
│
├── cmd/
│   └── server/
│       └── main.go             # Entry point: init config, DB, router, start server
│
└── internal/
    ├── config/
    │   └── config.go           # Load .env, parse config struct
    ├── database/
    │   └── database.go         # GORM connection, auto-migration
    ├── models/
    │   ├── user.go             # User model (built-in)
    │   └── post.go             # Generated models go here
    ├── handlers/
    │   ├── auth.go             # Auth endpoints (login, register, etc.)
    │   ├── user.go             # User CRUD endpoints
    │   └── post.go             # Generated handlers go here
    ├── services/
    │   ├── auth.go             # JWT generation, token validation
    │   ├── user.go             # User business logic
    │   └── post.go             # Generated services go here
    ├── middleware/
    │   ├── auth.go             # JWT validation, role-based access
    │   ├── cors.go             # CORS configuration
    │   └── logger.go           # Structured JSON logging
    ├── routes/
    │   └── routes.go           # Route registration (all endpoints)
    ├── cache/
    │   └── cache.go            # Redis caching service
    ├── storage/
    │   └── storage.go          # S3/R2/MinIO file storage
    ├── mail/
    │   ├── mailer.go           # Resend email service
    │   └── templates/          # HTML email templates
    ├── jobs/
    │   └── jobs.go             # Asynq background job queue
    ├── cron/
    │   └── cron.go             # Asynq cron scheduler
    └── ai/
        └── ai.go               # AI integration (Claude + OpenAI)`} />
            </div>

            <div className="prose-grit mb-10">
              <h3>Key Conventions</h3>
              <ul>
                <li><strong>Models</strong> define GORM structs with json, gorm, and binding tags. One file per model.</li>
                <li><strong>Handlers</strong> are thin HTTP controllers. They parse requests, call services, and return responses. No business logic in handlers.</li>
                <li><strong>Services</strong> contain business logic. They interact with the database through GORM and are called by handlers.</li>
                <li><strong>Middleware</strong> runs before handlers. Auth, CORS, and logging are pre-configured.</li>
                <li><strong>Routes</strong> are registered in a single file. Each resource group is clearly separated.</li>
              </ul>
            </div>

            {/* Next.js Web App */}
            <div className="prose-grit">
              <h2>Web App (apps/web/)</h2>
              <p>
                The main Next.js frontend application. Uses the App Router with route groups
                for auth and dashboard sections.
              </p>
            </div>

            <div className="mb-4">
              <CodeBlock filename="apps/web/" code={`apps/web/
├── package.json
├── next.config.ts
├── tailwind.config.ts
├── Dockerfile
│
├── app/
│   ├── layout.tsx              # Root layout with providers
│   ├── page.tsx                # Landing page
│   ├── (auth)/                 # Auth route group (no layout)
│   │   ├── login/page.tsx
│   │   ├── register/page.tsx
│   │   └── forgot-password/page.tsx
│   └── (dashboard)/            # Protected routes
│       ├── layout.tsx          # Dashboard layout (sidebar + navbar)
│       └── dashboard/page.tsx  # Main dashboard
│
├── components/
│   ├── ui/                     # shadcn/ui components
│   └── shared/                 # App-specific components
│
├── hooks/
│   ├── use-auth.ts             # Auth hooks (login, register, logout, me)
│   └── use-posts.ts            # Generated resource hooks
│
└── lib/
    ├── api-client.ts           # Axios instance with JWT interceptor
    ├── auth.ts                 # Auth utilities (token storage)
    └── utils.ts                # Utility functions`} />
            </div>

            <div className="prose-grit mb-10">
              <h3>Key Conventions</h3>
              <ul>
                <li><strong>Route groups</strong> <code>(auth)</code> and <code>(dashboard)</code> organize pages without affecting the URL structure.</li>
                <li><strong>Hooks</strong> wrap React Query mutations and queries. All data fetching goes through hooks, never raw fetch calls in components.</li>
                <li><strong>API client</strong> is a pre-configured Axios instance that automatically injects JWT tokens and handles token refresh.</li>
                <li><strong>UI components</strong> come from shadcn/ui. Add more with <code>pnpm dlx shadcn@latest add button</code>.</li>
              </ul>
            </div>

            {/* Admin Panel */}
            <div className="prose-grit">
              <h2>Admin Panel (apps/admin/)</h2>
              <p>
                The Filament-like admin panel. A separate Next.js app that provides resource
                management with data tables, forms, and dashboard widgets.
              </p>
            </div>

            <div className="mb-4">
              <CodeBlock filename="apps/admin/" code={`apps/admin/
├── package.json
├── next.config.ts
├── tailwind.config.ts
│
├── app/
│   ├── layout.tsx              # Admin layout with sidebar
│   ├── page.tsx                # Dashboard with widgets
│   └── resources/
│       ├── users/page.tsx      # User management page
│       └── posts/page.tsx      # Generated resource pages
│
├── components/
│   ├── layout/
│   │   ├── admin-layout.tsx    # Admin shell
│   │   ├── sidebar.tsx         # Collapsible sidebar
│   │   └── navbar.tsx          # Top navigation bar
│   ├── tables/
│   │   ├── data-table.tsx      # Server-side paginated table
│   │   ├── columns.tsx         # Column definitions
│   │   └── filters.tsx         # Table filters
│   ├── forms/
│   │   ├── form-builder.tsx    # Dynamic form renderer
│   │   ├── fields/             # Field type components
│   │   └── form-modal.tsx      # Modal form wrapper
│   └── widgets/
│       ├── stats-card.tsx      # Stat number + trend
│       ├── chart-widget.tsx    # Recharts wrapper
│       └── recent-activity.tsx # Activity feed
│
├── hooks/
│   ├── use-auth.ts             # Admin auth hooks
│   └── use-posts.ts            # Generated resource hooks
│
└── resources/                  # Resource definitions (THE MAGIC)
    ├── index.ts                # Resource registry
    ├── users.ts                # User resource config
    └── posts.ts                # Generated resource configs`} />
            </div>

            <div className="prose-grit mb-10">
              <h3>Key Conventions</h3>
              <ul>
                <li><strong>Resource definitions</strong> in <code>resources/</code> define the table columns, form fields, filters, and actions for each resource. This is how the admin panel generates its UI.</li>
                <li><strong>Data tables</strong> are server-side paginated. They communicate directly with the Go API for sorting, filtering, and searching.</li>
                <li><strong>Form builder</strong> renders forms dynamically from resource definitions. Field types include text, number, select, date, toggle, and file upload.</li>
                <li><strong>Widgets</strong> are dashboard components that fetch data from the API and display stats, charts, and activity feeds.</li>
              </ul>
            </div>

            {/* Shared Package */}
            <div className="prose-grit">
              <h2>Shared Package (packages/shared/)</h2>
              <p>
                The shared package contains TypeScript types, Zod validation schemas, and constants
                used by both the web app and admin panel. This is the glue that keeps the frontend
                in sync with the Go backend.
              </p>
            </div>

            <div className="mb-4">
              <CodeBlock filename="packages/shared/" code={`packages/shared/
├── package.json
│
├── schemas/                    # Zod validation schemas
│   ├── user.ts                 # User create/update schemas
│   ├── post.ts                 # Generated schemas
│   └── index.ts                # Re-exports all schemas
│
├── types/                      # TypeScript types
│   ├── user.ts                 # User type + API response types
│   ├── post.ts                 # Generated types
│   ├── api.ts                  # Pagination, error, response types
│   └── index.ts                # Re-exports all types
│
└── constants/
    └── index.ts                # Roles, API routes, config constants`} />
            </div>

            <div className="prose-grit mb-10">
              <h3>Key Conventions</h3>
              <ul>
                <li><strong>Schemas</strong> are Zod validation schemas that match the Go model struct tags. They are the source of truth for frontend validation.</li>
                <li><strong>Types</strong> are TypeScript interfaces that mirror Go structs. Generated by <code>grit sync</code> from the Go model definitions.</li>
                <li><strong>Constants</strong> include role strings, API route paths, and configuration values shared between all frontend apps.</li>
                <li>Both <code>apps/web</code> and <code>apps/admin</code> import from <code>@shared/schemas</code>, <code>@shared/types</code>, and <code>@shared/constants</code>.</li>
              </ul>
            </div>

            {/* Where Things Go */}
            <div className="prose-grit">
              <h2>Where Things Go</h2>
              <p>
                A quick reference for where to put different types of code:
              </p>
            </div>

            <div className="mb-10">
              <table className="w-full border-collapse">
                <thead>
                  <tr className="border-b border-border">
                    <th className="text-left text-sm font-semibold text-foreground pb-3 pr-4">What</th>
                    <th className="text-left text-sm font-semibold text-foreground pb-3">Where</th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    { what: 'New database model', where: 'apps/api/internal/models/<name>.go' },
                    { what: 'API endpoint handler', where: 'apps/api/internal/handlers/<name>.go' },
                    { what: 'Business logic', where: 'apps/api/internal/services/<name>.go' },
                    { what: 'Auth middleware', where: 'apps/api/internal/middleware/auth.go' },
                    { what: 'API route registration', where: 'apps/api/internal/routes/routes.go' },
                    { what: 'Environment config', where: 'apps/api/internal/config/config.go' },
                    { what: 'Background job', where: 'apps/api/internal/jobs/jobs.go' },
                    { what: 'Email template', where: 'apps/api/internal/mail/templates/' },
                    { what: 'Zod validation schema', where: 'packages/shared/schemas/<name>.ts' },
                    { what: 'TypeScript type', where: 'packages/shared/types/<name>.ts' },
                    { what: 'React Query hook', where: 'apps/web/hooks/use-<names>.ts' },
                    { what: 'Admin resource page', where: 'apps/admin/app/resources/<names>/page.tsx' },
                    { what: 'Admin resource definition', where: 'apps/admin/resources/<names>.ts' },
                    { what: 'Reusable UI component', where: 'apps/web/components/shared/' },
                    { what: 'shadcn/ui component', where: 'apps/web/components/ui/' },
                    { what: 'Dashboard widget', where: 'apps/admin/components/widgets/' },
                  ].map((row) => (
                    <tr key={row.what} className="border-b border-border/50">
                      <td className="text-sm text-foreground py-2.5 pr-4">{row.what}</td>
                      <td className="text-sm text-muted-foreground py-2.5 font-mono text-xs text-primary/60">{row.where}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="prose-grit">
              <h2>Generated vs. Hand-Written Code</h2>
              <p>
                When you run <code>grit generate resource</code>, the CLI creates files in all the
                locations listed above. These generated files are <strong>yours to modify</strong>.
                The CLI uses marker comments like <code>{`// grit:inject-routes`}</code> and{' '}
                <code>{`// grit:inject-models`}</code> to know where to inject new code into
                existing files (like <code>routes.go</code> and <code>database.go</code>).
              </p>
              <p>
                Do not remove these marker comments. They are how the CLI knows where to add new
                routes and model registrations when you generate additional resources.
              </p>
            </div>

            {/* Nav */}
            <div className="flex flex-wrap gap-3 mt-12 pt-6 border-t border-border/30">
              <Button variant="outline" asChild className="border-border/60 bg-transparent hover:bg-accent/50">
                <Link href="/docs/getting-started/installation">
                  Installation
                </Link>
              </Button>
              <Button asChild className="glow-purple-sm ml-auto">
                <Link href="/docs/getting-started/configuration">
                  Configuration
                  <ArrowRight className="ml-1.5 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
