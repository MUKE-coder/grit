import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { Files, Folder, File } from '@/components/files'
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
              <Files title="myapp/">
                <File name=".env" comment="Environment variables" />
                <File name=".env.example" comment="Template with documentation" />
                <File name=".env.cloud.example" comment="Cloud-only setup (no Docker)" />
                <File name=".gitignore" comment="Git ignore rules" />
                <File name="docker-compose.yml" comment="Dev: PostgreSQL, Redis, MinIO, Mailhog" />
                <File name="docker-compose.prod.yml" comment="Production: multi-stage builds" />
                <File name="grit.config.ts" comment="Grit framework configuration" />
                <File name="package.json" comment="Root package.json (workspace scripts)" />
                <File name="pnpm-workspace.yaml" comment="pnpm workspace definition" />
                <File name="turbo.json" comment="Turborepo task configuration" />
                <File name="README.md" comment="Project documentation" />
                <Folder name="apps" defaultOpen>
                  <Folder name="api" comment="Go backend (Gin + GORM)" />
                  <Folder name="web" comment="Next.js main frontend" />
                  <Folder name="admin" comment="Next.js admin panel" />
                </Folder>
                <Folder name="packages" defaultOpen>
                  <Folder name="shared" comment="Shared Zod schemas, TS types, constants" />
                </Folder>
              </Files>
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
              <Files title="apps/api/">
                <File name="go.mod" comment="Go module definition" />
                <File name="go.sum" comment="Dependency checksums" />
                <File name="Dockerfile" comment="Multi-stage production build" />
                <File name=".air.toml" comment="Hot reload configuration" />
                <Folder name="cmd" defaultOpen>
                  <Folder name="server" defaultOpen>
                    <File name="main.go" comment="Entry point: init config, DB, router, start server" />
                  </Folder>
                </Folder>
                <Folder name="internal" defaultOpen>
                  <Folder name="config" defaultOpen>
                    <File name="config.go" comment="Load .env, parse config struct" />
                  </Folder>
                  <Folder name="database" defaultOpen>
                    <File name="database.go" comment="GORM connection, auto-migration" />
                  </Folder>
                  <Folder name="models" defaultOpen>
                    <File name="user.go" comment="User model (built-in)" />
                    <File name="post.go" comment="Generated models go here" highlight />
                  </Folder>
                  <Folder name="handlers" defaultOpen>
                    <File name="auth.go" comment="Auth endpoints (login, register, etc.)" />
                    <File name="user.go" comment="User CRUD endpoints" />
                    <File name="post.go" comment="Generated handlers go here" highlight />
                  </Folder>
                  <Folder name="services" defaultOpen>
                    <File name="auth.go" comment="JWT generation, token validation" />
                    <File name="user.go" comment="User business logic" />
                    <File name="post.go" comment="Generated services go here" highlight />
                  </Folder>
                  <Folder name="middleware" defaultOpen>
                    <File name="auth.go" comment="JWT validation, role-based access" />
                    <File name="cors.go" comment="CORS configuration" />
                    <File name="logger.go" comment="Structured JSON logging" />
                  </Folder>
                  <Folder name="routes" defaultOpen>
                    <File name="routes.go" comment="Route registration (all endpoints)" />
                  </Folder>
                  <Folder name="cache" defaultOpen>
                    <File name="cache.go" comment="Redis caching service" />
                  </Folder>
                  <Folder name="storage" defaultOpen>
                    <File name="storage.go" comment="S3/R2/MinIO file storage" />
                  </Folder>
                  <Folder name="mail" defaultOpen>
                    <File name="mailer.go" comment="Resend email service" />
                    <Folder name="templates" comment="HTML email templates" />
                  </Folder>
                  <Folder name="jobs" defaultOpen>
                    <File name="jobs.go" comment="Asynq background job queue" />
                  </Folder>
                  <Folder name="cron" defaultOpen>
                    <File name="cron.go" comment="Asynq cron scheduler" />
                  </Folder>
                  <Folder name="ai" defaultOpen>
                    <File name="ai.go" comment="AI integration (Vercel AI Gateway — one key, many models)" />
                  </Folder>
                </Folder>
              </Files>
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
              <Files title="apps/web/">
                <File name="package.json" />
                <File name="next.config.ts" />
                <File name="tailwind.config.ts" />
                <File name="Dockerfile" />
                <Folder name="app" defaultOpen>
                  <File name="layout.tsx" comment="Root layout with providers" />
                  <File name="page.tsx" comment="Landing page" />
                  <Folder name="(auth)" comment="Auth route group (no layout)" defaultOpen>
                    <File name="login/page.tsx" />
                    <File name="register/page.tsx" />
                    <File name="forgot-password/page.tsx" />
                  </Folder>
                  <Folder name="(dashboard)" comment="Protected routes" defaultOpen>
                    <File name="layout.tsx" comment="Dashboard layout (sidebar + navbar)" />
                    <File name="dashboard/page.tsx" comment="Main dashboard" />
                  </Folder>
                </Folder>
                <Folder name="components" defaultOpen>
                  <Folder name="ui" comment="shadcn/ui components" />
                  <Folder name="shared" comment="App-specific components" />
                </Folder>
                <Folder name="hooks" defaultOpen>
                  <File name="use-auth.ts" comment="Auth hooks (login, register, logout, me)" />
                  <File name="use-posts.ts" comment="Generated resource hooks" highlight />
                </Folder>
                <Folder name="lib" defaultOpen>
                  <File name="api-client.ts" comment="Axios instance with JWT interceptor" />
                  <File name="auth.ts" comment="Auth utilities (token storage)" />
                  <File name="utils.ts" comment="Utility functions" />
                </Folder>
              </Files>
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
              <Files title="apps/admin/">
                <File name="package.json" />
                <File name="next.config.ts" />
                <File name="tailwind.config.ts" />
                <Folder name="app" defaultOpen>
                  <File name="layout.tsx" comment="Admin layout with sidebar" />
                  <File name="page.tsx" comment="Dashboard with widgets" />
                  <Folder name="resources" defaultOpen>
                    <File name="users/page.tsx" comment="User management page" />
                    <File name="posts/page.tsx" comment="Generated resource pages" highlight />
                  </Folder>
                </Folder>
                <Folder name="components" defaultOpen>
                  <Folder name="layout" defaultOpen>
                    <File name="admin-layout.tsx" comment="Admin shell" />
                    <File name="sidebar.tsx" comment="Collapsible sidebar" />
                    <File name="navbar.tsx" comment="Top navigation bar" />
                  </Folder>
                  <Folder name="tables" defaultOpen>
                    <File name="data-table.tsx" comment="Server-side paginated table" />
                    <File name="columns.tsx" comment="Column definitions" />
                    <File name="filters.tsx" comment="Table filters" />
                  </Folder>
                  <Folder name="forms" defaultOpen>
                    <File name="form-builder.tsx" comment="Dynamic form renderer" />
                    <Folder name="fields" comment="Field type components" />
                    <File name="form-modal.tsx" comment="Modal form wrapper" />
                  </Folder>
                  <Folder name="widgets" defaultOpen>
                    <File name="stats-card.tsx" comment="Stat number + trend" />
                    <File name="chart-widget.tsx" comment="Recharts wrapper" />
                    <File name="recent-activity.tsx" comment="Activity feed" />
                  </Folder>
                </Folder>
                <Folder name="hooks" defaultOpen>
                  <File name="use-auth.ts" comment="Admin auth hooks" />
                  <File name="use-posts.ts" comment="Generated resource hooks" highlight />
                </Folder>
                <Folder name="resources" comment="Resource definitions (THE MAGIC)" defaultOpen>
                  <File name="index.ts" comment="Resource registry" />
                  <File name="users.ts" comment="User resource config" />
                  <File name="posts.ts" comment="Generated resource configs" highlight />
                </Folder>
              </Files>
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
              <Files title="packages/shared/">
                <File name="package.json" />
                <Folder name="schemas" comment="Zod validation schemas" defaultOpen>
                  <File name="user.ts" comment="User create/update schemas" />
                  <File name="post.ts" comment="Generated schemas" highlight />
                  <File name="index.ts" comment="Re-exports all schemas" />
                </Folder>
                <Folder name="types" comment="TypeScript types" defaultOpen>
                  <File name="user.ts" comment="User type + API response types" />
                  <File name="post.ts" comment="Generated types" highlight />
                  <File name="api.ts" comment="Pagination, error, response types" />
                  <File name="index.ts" comment="Re-exports all types" />
                </Folder>
                <Folder name="constants" defaultOpen>
                  <File name="index.ts" comment="Roles, API routes, config constants" />
                </Folder>
              </Files>
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
