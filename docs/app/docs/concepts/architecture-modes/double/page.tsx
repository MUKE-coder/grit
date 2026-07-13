import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { FileTree } from '@/components/diagram'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/architecture-modes/double')

export default function DoubleArchitecturePage() {
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
              Double Architecture: Web + API
            </h1>
            <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl">
              A streamlined Turborepo monorepo with two applications: a public-facing web app
              and a Go API. Admin features live inside the web app as role-protected routes
              instead of a separate application.
            </p>
          </div>

          {/* ── Section 1: Overview ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              Overview
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Double architecture strips away the dedicated admin panel and gives you a two-app monorepo.
              The Go API handles all backend logic (authentication, CRUD, file uploads, jobs), and a
              single frontend serves both regular users and administrators. Admin functionality is
              implemented through role-protected routes within the web app — users with the ADMIN role
              see additional navigation items and have access to management pages.
            </p>
            <CodeBlock language="bash" filename="scaffold command" code={`# Double with TanStack Router (recommended)
grit new myapp --double --vite

# Double with Next.js
grit new myapp --double --next`} />

            <div className="mt-6 rounded-lg border border-border/40 bg-accent/20 p-6">
              <h3 className="text-sm font-semibold text-foreground mb-3">Best suited for</h3>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li className="flex items-start gap-2"><span className="text-primary mt-0.5">-</span> Projects where admins and users share the same interface with extra privileges</li>
                <li className="flex items-start gap-2"><span className="text-primary mt-0.5">-</span> Blogs, portfolios, and content sites where the author manages content inline</li>
                <li className="flex items-start gap-2"><span className="text-primary mt-0.5">-</span> Simpler SaaS applications where admin features are minimal (user management, settings)</li>
                <li className="flex items-start gap-2"><span className="text-primary mt-0.5">-</span> MERN/MEAN stack developers transitioning to Go who want a familiar two-app structure</li>
                <li className="flex items-start gap-2"><span className="text-primary mt-0.5">-</span> Teams that want fewer apps to deploy and maintain</li>
              </ul>
            </div>

            <div className="mt-4 rounded-lg border border-primary/20 bg-primary/5 p-6">
              <h3 className="text-sm font-semibold text-foreground mb-3">Key differences from Triple</h3>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li className="flex items-start gap-2"><span className="text-primary">-</span> <strong>No apps/admin/ directory</strong> — the admin panel does not exist as a separate application</li>
                <li className="flex items-start gap-2"><span className="text-primary">-</span> <strong>No admin resource definitions</strong> — no <code className="text-xs font-mono bg-accent/30 px-1 rounded">defineResource()</code>, no DataTable/FormBuilder auto-generation</li>
                <li className="flex items-start gap-2"><span className="text-primary">-</span> <strong>Admin features in web app</strong> — role-protected routes (e.g., <code className="text-xs font-mono bg-accent/30 px-1 rounded">/admin/users</code>) guarded by <code className="text-xs font-mono bg-accent/30 px-1 rounded">RequireRole(&quot;ADMIN&quot;)</code></li>
                <li className="flex items-start gap-2"><span className="text-primary">-</span> <strong>Fewer generated files</strong> — <code className="text-xs font-mono bg-accent/30 px-1 rounded">grit generate</code> creates Go files + shared types + web hooks, but no admin page or resource definition</li>
                <li className="flex items-start gap-2"><span className="text-primary">-</span> <strong>Simpler deployment</strong> — 2 apps instead of 3, fewer Docker images, less infrastructure</li>
              </ul>
            </div>
          </section>

          {/* ── Section 2: Complete Folder Structure ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              Complete Folder Structure
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              This is the full tree generated by <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new myapp --double --vite</code>.
              Compare it with the triple structure — the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/admin/</code> directory
              is absent, and the web app gains role-protected admin routes.
            </p>
            <FileTree
              title="myapp/"
              nodes={[
                { name: '.env', type: 'file', depth: 0, comment: 'Shared environment variables' },
                { name: '.env.example', type: 'file', depth: 0, comment: 'Template for other developers' },
                { name: '.gitignore', type: 'file', depth: 0 },
                { name: '.prettierrc', type: 'file', depth: 0, comment: 'Code formatting' },
                { name: '.prettierignore', type: 'file', depth: 0 },
                { name: 'docker-compose.yml', type: 'file', depth: 0, comment: 'Dev: PostgreSQL, Redis, MinIO, Mailhog' },
                { name: 'docker-compose.prod.yml', type: 'file', depth: 0, comment: 'Production deployment' },
                { name: 'grit.json', type: 'file', depth: 0, comment: 'Project manifest (architecture: "double")' },
                { name: 'turbo.json', type: 'file', depth: 0, comment: 'Turborepo pipeline config' },
                { name: 'pnpm-workspace.yaml', type: 'file', depth: 0, comment: 'pnpm workspace definition' },
                { name: 'package.json', type: 'file', depth: 0, comment: 'Root scripts (dev, build, lint)' },
                { name: '.claude/', type: 'folder', depth: 0 },
                { name: 'skills/grit/', type: 'folder', depth: 1 },
                { name: 'SKILL.md', type: 'file', depth: 2, comment: 'AI assistant guide (tailored to double)' },
                { name: 'reference.md', type: 'file', depth: 2, comment: 'Detailed API conventions' },
                { name: 'packages/', type: 'folder', depth: 0 },
                { name: 'shared/', type: 'folder', depth: 1, comment: 'Shared between API and web' },
                { name: 'schemas/', type: 'folder', depth: 2, comment: 'Zod validation schemas' },
                { name: 'index.ts', type: 'file', depth: 3, comment: '// grit:schemas marker' },
                { name: 'user.ts', type: 'file', depth: 3, comment: 'User schema (register, login, update)' },
                { name: 'types/', type: 'folder', depth: 2, comment: 'TypeScript interfaces' },
                { name: 'index.ts', type: 'file', depth: 3, comment: '// grit:types marker' },
                { name: 'user.ts', type: 'file', depth: 3, comment: 'User type' },
                { name: 'constants/', type: 'folder', depth: 2, comment: 'API_ROUTES, ROLES // grit:api-routes' },
                { name: 'apps/', type: 'folder', depth: 0 },
                { name: 'api/', type: 'folder', depth: 1, comment: 'Go backend (identical to triple)' },
                { name: 'Dockerfile', type: 'file', depth: 2, comment: 'Multi-stage build (golang → alpine)' },
                { name: 'go.mod', type: 'file', depth: 2, comment: 'Go module: myapp/apps/api' },
                { name: 'go.sum', type: 'file', depth: 2 },
                { name: 'cmd/', type: 'folder', depth: 2 },
                { name: 'server/main.go', type: 'file', depth: 3, comment: 'Entry point: config, db, services, router' },
                { name: 'migrate/main.go', type: 'file', depth: 3, comment: 'Migration runner' },
                { name: 'seed/main.go', type: 'file', depth: 3, comment: 'Database seeder' },
                { name: 'internal/', type: 'folder', depth: 2 },
                { name: 'config/config.go', type: 'file', depth: 3, comment: 'Loads .env, all env vars' },
                { name: 'database/db.go', type: 'file', depth: 3, comment: 'GORM connection + AutoMigrate' },
                { name: 'models/', type: 'folder', depth: 3, comment: 'GORM models' },
                { name: 'user.go', type: 'file', depth: 4, comment: 'User model // grit:models marker' },
                { name: 'upload.go', type: 'file', depth: 4, comment: 'Upload model' },
                { name: 'handlers/', type: 'folder', depth: 3, comment: 'HTTP handlers (thin — call services)' },
                { name: 'auth.go', type: 'file', depth: 4, comment: 'Register, login, refresh, me' },
                { name: 'user.go', type: 'file', depth: 4, comment: 'User CRUD' },
                { name: 'upload.go', type: 'file', depth: 4, comment: 'File upload (presigned URLs)' },
                { name: 'ai.go', type: 'file', depth: 4, comment: 'AI completions + streaming' },
                { name: 'services/', type: 'folder', depth: 3, comment: 'Business logic' },
                { name: 'auth_service.go', type: 'file', depth: 4, comment: 'JWT, bcrypt, token generation' },
                { name: 'user_service.go', type: 'file', depth: 4, comment: 'User queries' },
                { name: 'middleware/', type: 'folder', depth: 3, comment: 'Gin middleware' },
                { name: 'auth.go', type: 'file', depth: 4, comment: 'JWT verification' },
                { name: 'cors.go', type: 'file', depth: 4, comment: 'CORS configuration' },
                { name: 'logger.go', type: 'file', depth: 4, comment: 'Structured logging' },
                { name: 'cache.go', type: 'file', depth: 4, comment: 'Redis response caching' },
                { name: 'maintenance.go', type: 'file', depth: 4, comment: 'grit down/up support' },
                { name: 'rate_limit.go', type: 'file', depth: 4, comment: 'Sentinel rate limiting' },
                { name: 'routes/', type: 'folder', depth: 3 },
                { name: 'routes.go', type: 'file', depth: 4, comment: 'Route registration // grit:handlers, grit:routes:*' },
                { name: 'mail/', type: 'folder', depth: 3, comment: 'Email service (Resend)' },
                { name: 'mailer.go', type: 'file', depth: 4, comment: 'Send function' },
                { name: 'templates/', type: 'folder', depth: 4, comment: 'HTML email templates' },
                { name: 'storage/', type: 'folder', depth: 3, comment: 'S3-compatible file storage' },
                { name: 'jobs/', type: 'folder', depth: 3, comment: 'Background jobs (asynq)' },
                { name: 'cron/', type: 'folder', depth: 3, comment: 'Scheduled tasks' },
                { name: 'cache/', type: 'folder', depth: 3, comment: 'Redis cache service' },
                { name: 'ai/', type: 'folder', depth: 3, comment: 'AI service (Vercel AI Gateway)' },
                { name: 'auth/', type: 'folder', depth: 3, comment: 'TOTP 2FA service' },
                { name: 'totp.go', type: 'file', depth: 4, comment: 'Setup, verify, backup codes, trusted devices' },
                { name: 'web/', type: 'folder', depth: 1, comment: 'Frontend (public + admin routes)' },
                { name: 'Dockerfile', type: 'file', depth: 2, comment: 'Next.js standalone or Vite static build' },
                { name: 'package.json', type: 'file', depth: 2, comment: 'Dependencies + scripts' },
                { name: 'vite.config.ts', type: 'file', depth: 2, comment: '(or next.config.js for --next)' },
                { name: 'tailwind.config.ts', type: 'file', depth: 2 },
                { name: 'src/', type: 'folder', depth: 2, comment: 'TanStack Router (or app/ for Next.js)' },
                { name: 'routes/', type: 'folder', depth: 3 },
                { name: '__root.tsx', type: 'file', depth: 4, comment: 'Root layout' },
                { name: 'index.tsx', type: 'file', depth: 4, comment: 'Landing page' },
                { name: '(auth)/', type: 'folder', depth: 4, comment: 'Auth pages (login, register)' },
                { name: '(app)/', type: 'folder', depth: 4, comment: 'Protected user pages' },
                { name: '(admin)/', type: 'folder', depth: 4, comment: 'Role-protected admin pages' },
                { name: '_layout.tsx', type: 'file', depth: 5, comment: 'Admin layout (checks ADMIN role)' },
                { name: 'users/', type: 'folder', depth: 5, comment: 'User management' },
                { name: 'settings/', type: 'folder', depth: 5, comment: 'System settings' },
                { name: 'hooks/', type: 'folder', depth: 3, comment: 'React Query hooks' },
                { name: 'components/', type: 'folder', depth: 3, comment: 'Shared components' },
                { name: 'lib/', type: 'folder', depth: 3, comment: 'Utilities, API client' },
              ]}
            />
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
                The root configuration is nearly identical to triple. The key difference is
                in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit.json</code> which
                stores <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">{`"architecture": "double"`}</code>,
                telling the CLI to skip admin-specific file generation.
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
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">grit.json</td><td className="p-3">Project manifest. Architecture is set to <code className="text-xs font-mono bg-accent/30 px-1 rounded">double</code> — the CLI uses this to determine which files to generate</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">turbo.json</td><td className="p-3">Turborepo pipelines: <code className="text-xs font-mono bg-accent/30 px-1 rounded">dev</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">build</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">lint</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">test</code>. Simpler dependency graph with only 2 apps</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">docker-compose.yml</td><td className="p-3">Development infrastructure: PostgreSQL 16, Redis 7, MinIO, Mailhog</td></tr>
                    <tr><td className="p-3 font-mono text-xs">docker-compose.prod.yml</td><td className="p-3">Production deployment with 2 app containers (API + web) instead of 3</td></tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* packages/shared */}
            <div className="mb-8">
              <h3 className="text-lg font-semibold text-foreground mb-3">packages/shared/</h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Identical to the triple architecture. Contains Zod schemas, TypeScript types, and
                shared constants. The only difference is that there is one consumer (web) instead
                of two (web + admin). The marker comments and injection system work the same way.
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
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">schemas/</td><td className="p-3">Zod validation schemas per resource. Re-exported from <code className="text-xs font-mono bg-accent/30 px-1 rounded">index.ts</code> via the <code className="text-xs font-mono bg-accent/30 px-1 rounded">// grit:schemas</code> marker.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">types/</td><td className="p-3">TypeScript interfaces per resource. Re-exported via <code className="text-xs font-mono bg-accent/30 px-1 rounded">// grit:types</code> marker.</td></tr>
                    <tr><td className="p-3 font-mono text-xs">constants/</td><td className="p-3">API route constants and role enums. Injected via <code className="text-xs font-mono bg-accent/30 px-1 rounded">// grit:api-routes</code> marker.</td></tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* apps/api */}
            <div className="mb-8">
              <h3 className="text-lg font-semibold text-foreground mb-3">apps/api/ (Go Backend)</h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The Go API is <strong>identical</strong> to the triple architecture. It uses the same
                handler-service-model pattern, the same middleware stack, the same batteries (cache, storage,
                email, jobs, cron, AI). The only difference is that the CORS middleware is configured to allow
                a single frontend origin instead of two.
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
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">cmd/server/</td><td className="p-3">Entry point: loads config, connects to DB, starts Gin on :8080</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">cmd/migrate/</td><td className="p-3">Standalone migration runner via GORM AutoMigrate</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">cmd/seed/</td><td className="p-3">Database seeder with admin user and sample data</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/config/</td><td className="p-3">Typed config struct loaded from <code className="text-xs font-mono bg-accent/30 px-1 rounded">.env</code></td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/database/</td><td className="p-3">GORM connection: PostgreSQL (production) or SQLite (dev/test)</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/models/</td><td className="p-3">GORM model definitions with struct tags. <code className="text-xs font-mono bg-accent/30 px-1 rounded">// grit:models</code> marker for injection.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/handlers/</td><td className="p-3">Thin Gin handlers: parse input, call service, return JSON</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/services/</td><td className="p-3">Business logic: CRUD with pagination, filtering, error handling</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/middleware/</td><td className="p-3">JWT auth, CORS, logger, cache, rate limiter, maintenance mode</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/routes/</td><td className="p-3">Centralized route registration with marker comments</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/mail/</td><td className="p-3">Resend email service with HTML templates</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/storage/</td><td className="p-3">S3-compatible file storage with presigned URLs</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/jobs/</td><td className="p-3">Background job processing with asynq + Redis</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/cron/</td><td className="p-3">Scheduled tasks with cron expressions</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/cache/</td><td className="p-3">Redis cache service with TTL</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">internal/ai/</td><td className="p-3">Claude + OpenAI with streaming</td></tr>
                    <tr><td className="p-3 font-mono text-xs">internal/auth/</td><td className="p-3">TOTP 2FA: setup, verify, backup codes, trusted devices</td></tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* apps/web */}
            <div className="mb-8">
              <h3 className="text-lg font-semibold text-foreground mb-3">apps/web/ (Public + Admin Frontend)</h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                In double architecture, the web app serves double duty. Public pages are accessible to
                everyone, while admin routes are protected by role checks. The navigation dynamically
                shows admin links (e.g., &quot;Users&quot;, &quot;Settings&quot;) only when the logged-in user has
                the ADMIN role.
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
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">src/routes/__root.tsx</td><td className="p-3">Root layout: providers, fonts, global styles (TanStack Router). For Next.js: <code className="text-xs font-mono bg-accent/30 px-1 rounded">app/layout.tsx</code></td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">src/routes/index.tsx</td><td className="p-3">Landing page visible to all visitors</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">src/routes/(auth)/</td><td className="p-3">Auth route group: login, register, forgot-password</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">src/routes/(app)/</td><td className="p-3">Protected user pages: dashboard, profile, resource pages</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">src/routes/(admin)/</td><td className="p-3">Role-protected admin routes. The layout component checks the user role and redirects non-admins. Contains user management, system settings, and other admin-only pages.</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">src/hooks/</td><td className="p-3">React Query hooks for data fetching (generated by <code className="text-xs font-mono bg-accent/30 px-1 rounded">grit generate</code>)</td></tr>
                    <tr><td className="p-3 font-mono text-xs">src/lib/</td><td className="p-3">Utilities: API client (Axios/fetch wrapper), auth helpers, formatters</td></tr>
                  </tbody>
                </table>
              </div>

              <div className="mt-4 rounded-lg border border-border/40 bg-accent/20 p-5">
                <h4 className="text-sm font-semibold text-foreground mb-3">How role-protected admin routes work</h4>
                <CodeBlock language="tsx" filename="src/routes/(admin)/_layout.tsx" code={`// The admin layout checks the user's role before rendering children
import { useAuth } from '@/hooks/use-auth'
import { Navigate, Outlet } from '@tanstack/react-router'

export default function AdminLayout() {
  const { user } = useAuth()

  // Redirect non-admin users
  if (user?.role !== 'ADMIN') {
    return <Navigate to="/" />
  }

  return (
    <div className="flex">
      {/* Admin sidebar with management links */}
      <AdminSidebar />
      <main className="flex-1 p-6">
        <Outlet />
      </main>
    </div>
  )
}`} />
                <p className="text-xs text-muted-foreground/60 mt-3">
                  The backend also enforces role checks via the <code className="text-xs font-mono bg-accent/30 px-1 rounded">RequireRole(&quot;ADMIN&quot;)</code> middleware
                  on admin API routes, so the frontend check is for UX only — not a security boundary.
                </p>
              </div>
            </div>
          </section>

          {/* ── Section 4: Data Flow ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              Data Flow
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              The double architecture has a simpler data flow with a single frontend talking to the API.
              Admin and user requests follow the same path — the only difference is the role attached to
              the JWT token.
            </p>
            <LaneFlow
              id="double"
              lanes={['Web App :3000', 'Go API — :8080', 'Data & Services']}
              groups={[{ lane: 1, rows: [0, 2], label: 'Request pipeline', tone: 'primary' }]}
              nodes={[
                { id: 'public', lane: 0, row: 0, title: 'Public Pages', sub: '/ · /login', tone: 'blue' },
                { id: 'admin', lane: 0, row: 1, title: 'Admin Pages', sub: 'role-gated /admin/*', tone: 'cyan' },
                { id: 'router', lane: 1, row: 0, title: 'Gin Router', sub: 'REST + JWT', tone: 'primary' },
                { id: 'mw', lane: 1, row: 1, title: 'Middleware', sub: 'Auth · RequireRole · Cache', tone: 'primary' },
                { id: 'svc', lane: 1, row: 2, title: 'Handler → Service', sub: 'GORM', tone: 'primary' },
                { id: 'pg', lane: 2, row: 0, title: 'PostgreSQL', sub: ':5434 · data', tone: 'green' },
                { id: 'redis', lane: 2, row: 1, title: 'Redis', sub: ':6380 · cache + jobs', tone: 'rose' },
                { id: 'minio', lane: 2, row: 2, title: 'MinIO', sub: ':9002 · files', tone: 'amber' },
                { id: 'resend', lane: 2, row: 3, title: 'Resend', sub: 'email', tone: 'violet' },
              ]}
              edges={[
                { from: 'public', to: 'router', label: 'REST + JWT', tone: 'blue' },
                { from: 'admin', to: 'router', label: 'ADMIN role', tone: 'cyan' },
                { from: 'router', to: 'mw', tone: 'primary' },
                { from: 'mw', to: 'svc', tone: 'primary' },
                { from: 'svc', to: 'pg', label: 'query', tone: 'green' },
                { from: 'svc', to: 'redis', label: 'cache', tone: 'rose' },
                { from: 'svc', to: 'minio', label: 'files', tone: 'amber' },
                { from: 'svc', to: 'resend', label: 'mail', tone: 'violet' },
              ]}
              legend={[
                { tone: 'blue', label: 'Public UI' },
                { tone: 'cyan', label: 'Admin UI (role-gated)' },
                { tone: 'primary', label: 'Go API' },
                { tone: 'green', label: 'Data & services' },
              ]}
            />
            <div className="mt-3 rounded-md border border-border/40 bg-card/20 px-3 py-2 text-[12px] text-muted-foreground/80">
              <span className="font-semibold text-foreground/70">Background</span> — asynq workers · cron · GORM Studio (<code>/studio</code>)
            </div>

            <div className="mt-6 rounded-lg border border-border/40 bg-accent/20 p-4">
              <h4 className="text-sm font-semibold text-foreground mb-2">Step-by-step: Admin managing users</h4>
              <ol className="space-y-2 text-sm text-muted-foreground list-decimal list-inside">
                <li>Admin navigates to <code className="text-xs font-mono bg-accent/30 px-1 rounded">/admin/users</code> in the web app</li>
                <li>The admin layout component verifies the user&apos;s role is ADMIN (frontend guard)</li>
                <li>React Query hook calls <code className="text-xs font-mono bg-accent/30 px-1 rounded">GET /api/v3/admin/users</code> with the JWT token</li>
                <li>Gin router matches the admin route group and applies <code className="text-xs font-mono bg-accent/30 px-1 rounded">RequireRole(&quot;ADMIN&quot;)</code> middleware (backend guard)</li>
                <li>Handler calls <code className="text-xs font-mono bg-accent/30 px-1 rounded">userService.GetAll()</code> and returns the paginated user list</li>
                <li>The web app renders the user table — built with your own components (not the auto-generated DataTable from triple)</li>
              </ol>
            </div>
          </section>

          {/* ── Section 5: Code Generation ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              How Code Generation Works
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              In double architecture, <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate</code> creates
              fewer files compared to triple because there is no admin panel to scaffold into. The Go backend
              files and shared types are identical — the difference is the absence of admin resource definitions
              and admin pages.
            </p>
            <CodeBlock language="bash" filename="example command" code={`grit generate resource Product --fields "name:string, price:float, description:text, category_id:belongs_to, is_active:bool"`} />

            <div className="mt-6 space-y-4">
              <h3 className="text-lg font-semibold text-foreground">Files Created</h3>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-5">
                <h4 className="text-sm font-mono text-primary mb-3">Go Backend (apps/api/)</h4>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> <code className="text-xs font-mono bg-accent/30 px-1 rounded">internal/models/product.go</code> — GORM model with struct tags, belongs_to relationship</li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> <code className="text-xs font-mono bg-accent/30 px-1 rounded">internal/services/product_service.go</code> — CRUD operations with pagination</li>
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
                <h4 className="text-sm font-mono text-primary mb-3">Web App (apps/web/)</h4>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> <code className="text-xs font-mono bg-accent/30 px-1 rounded">hooks/use-products.ts</code> — React Query hooks: <code className="text-xs font-mono bg-accent/30 px-1 rounded">useProducts()</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">useProduct(id)</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">useCreateProduct()</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">useUpdateProduct()</code>, <code className="text-xs font-mono bg-accent/30 px-1 rounded">useDeleteProduct()</code></li>
                </ul>
              </div>

              <div className="rounded-lg border border-primary/20 bg-primary/5 p-5">
                <h4 className="text-sm font-mono text-primary mb-3">NOT generated (compared to triple)</h4>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li className="flex items-start gap-2"><span className="text-muted-foreground">-</span> <s>resources/products.ts</s> — No admin resource definition</li>
                  <li className="flex items-start gap-2"><span className="text-muted-foreground">-</span> <s>app/(dashboard)/products/page.tsx</s> — No auto-generated admin page with DataTable and FormBuilder</li>
                </ul>
                <p className="text-xs text-muted-foreground/60 mt-3">
                  To manage products as an admin in double architecture, you build the management UI yourself
                  inside the <code className="text-xs font-mono bg-accent/30 px-1 rounded">src/routes/(admin)/</code> route group,
                  using the generated React Query hooks and shared types.
                </p>
              </div>
            </div>

            <div className="mt-6">
              <h3 className="text-lg font-semibold text-foreground mb-3">Marker Injections</h3>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Double uses 5 marker injections instead of 6 (no admin resource marker).
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
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">routes/routes.go</td><td className="p-3 font-mono text-xs">// grit:routes:*</td><td className="p-3">Route group registration</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">schemas/index.ts</td><td className="p-3 font-mono text-xs">// grit:schemas</td><td className="p-3">Schema re-export</td></tr>
                    <tr className="border-b border-border/20"><td className="p-3 font-mono text-xs">types/index.ts</td><td className="p-3 font-mono text-xs">// grit:types</td><td className="p-3">Type re-export</td></tr>
                    <tr><td className="p-3 font-mono text-xs">constants/index.ts</td><td className="p-3 font-mono text-xs">// grit:api-routes</td><td className="p-3">API route constant</td></tr>
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
              With one fewer frontend, the port layout is simpler. No port 3001 is occupied.
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
                  <tr className="border-b border-border/20"><td className="p-3 font-semibold text-foreground">Web App</td><td className="p-3 font-mono text-xs">3000</td><td className="p-3 font-mono text-xs">http://localhost:3000</td><td className="p-3">Next.js or Vite (:5173) — serves both public and admin UI</td></tr>
                  <tr className="border-b border-border/20"><td className="p-3 font-semibold text-foreground">PostgreSQL</td><td className="p-3 font-mono text-xs">5434</td><td className="p-3 font-mono text-xs">localhost:5434</td><td className="p-3">Docker container (host port 5434 → 5432)</td></tr>
                  <tr className="border-b border-border/20"><td className="p-3 font-semibold text-foreground">Redis</td><td className="p-3 font-mono text-xs">6380</td><td className="p-3 font-mono text-xs">localhost:6380</td><td className="p-3">Cache + job queue (host port 6380 → 6379)</td></tr>
                  <tr className="border-b border-border/20"><td className="p-3 font-semibold text-foreground">MinIO API</td><td className="p-3 font-mono text-xs">9002</td><td className="p-3 font-mono text-xs">http://localhost:9002</td><td className="p-3">S3-compatible API</td></tr>
                  <tr className="border-b border-border/20"><td className="p-3 font-semibold text-foreground">MinIO Console</td><td className="p-3 font-mono text-xs">9003</td><td className="p-3 font-mono text-xs">http://localhost:9003</td><td className="p-3">Web UI for managing buckets</td></tr>
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
              Deploying double is simpler than triple — two containers instead of three, one frontend
              domain instead of two.
            </p>

            <div className="space-y-4">
              <div className="rounded-lg border border-border/40 bg-accent/20 p-6">
                <h3 className="text-base font-semibold text-foreground mb-2">Docker Compose (Self-hosted)</h3>
                <p className="text-sm text-muted-foreground mb-3">
                  The production Docker Compose file builds the Go API and the web app as optimized images.
                  With only 2 app containers, resource usage is lower and the configuration is simpler.
                </p>
                <CodeBlock language="bash" code={`# Build and deploy
docker compose -f docker-compose.prod.yml up -d --build`} />
              </div>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-6">
                <h3 className="text-base font-semibold text-foreground mb-2">grit deploy</h3>
                <p className="text-sm text-muted-foreground">
                  Works the same as triple — <code className="text-xs font-mono bg-accent/30 px-1 rounded">grit deploy</code> reads
                  the <code className="text-xs font-mono bg-accent/30 px-1 rounded">grit.json</code> manifest to determine which
                  apps to build and deploy. With double architecture, it skips the admin container entirely.
                </p>
              </div>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-6">
                <h3 className="text-base font-semibold text-foreground mb-2">Hybrid: Vercel + VPS</h3>
                <p className="text-sm text-muted-foreground">
                  Deploy the web app (Next.js) to Vercel and the Go API to a VPS. Simpler than triple
                  because you only have one Vercel project to configure instead of two. Set
                  <code className="text-xs font-mono bg-accent/30 px-1 rounded">NEXT_PUBLIC_API_URL</code> in Vercel
                  to point to your API server.
                </p>
              </div>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-6">
                <h3 className="text-base font-semibold text-foreground mb-2">Static Deploy (Vite only)</h3>
                <p className="text-sm text-muted-foreground">
                  If you chose <code className="text-xs font-mono bg-accent/30 px-1 rounded">--vite</code>, the web app builds
                  to static files that can be served from any CDN (Cloudflare Pages, Netlify, S3 + CloudFront).
                  The Go API is the only server-side component. This is the most cost-effective deployment
                  option.
                </p>
              </div>
            </div>
          </section>

          {/* ── Section 8: When to Choose Double ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              When to Choose Double
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Double is the right choice when you want the power of a monorepo with shared types
              but don&apos;t need a separate admin application.
            </p>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="rounded-lg border border-primary/20 bg-primary/5 p-6">
                <h3 className="text-sm font-semibold text-foreground mb-3">Choose Double when...</h3>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> Admins and users share the same UI with role-based visibility</li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> Your admin features are simple (user list, settings page)</li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> You are building a blog, portfolio, or personal SaaS</li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> You want fewer apps to deploy and maintain</li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> You prefer building admin UI manually with your own components</li>
                  <li className="flex items-start gap-2"><span className="text-primary">-</span> You are familiar with MERN/MEAN stack structure (API + SPA)</li>
                </ul>
              </div>
              <div className="rounded-lg border border-border/40 bg-accent/20 p-6">
                <h3 className="text-sm font-semibold text-foreground mb-3">Consider Triple instead when...</h3>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li className="flex items-start gap-2"><span className="text-muted-foreground">-</span> You need DataTable + FormBuilder auto-generated admin pages</li>
                  <li className="flex items-start gap-2"><span className="text-muted-foreground">-</span> Non-technical admins need a polished, dedicated dashboard</li>
                  <li className="flex items-start gap-2"><span className="text-muted-foreground">-</span> You want <code className="text-xs font-mono bg-accent/30 px-1 rounded">grit generate</code> to create admin pages automatically</li>
                  <li className="flex items-start gap-2"><span className="text-muted-foreground">-</span> Your admin operations are complex (multi-step forms, widgets, analytics)</li>
                  <li className="flex items-start gap-2"><span className="text-muted-foreground">-</span> Separate teams work on the public site vs admin panel</li>
                </ul>
              </div>
            </div>

            <div className="mt-6 rounded-lg border border-border/40 bg-accent/20 p-5">
              <h3 className="text-sm font-semibold text-foreground mb-2">Upgrading from Double to Triple</h3>
              <p className="text-sm text-muted-foreground">
                If your project outgrows double architecture, you can upgrade to triple later. Run
                <code className="text-xs font-mono bg-accent/30 px-1 rounded">grit upgrade --triple</code> to scaffold
                the admin app into your existing project. Your existing Go API, shared types, and web app
                remain untouched — the command only adds the <code className="text-xs font-mono bg-accent/30 px-1 rounded">apps/admin/</code> directory
                and updates <code className="text-xs font-mono bg-accent/30 px-1 rounded">grit.json</code>.
              </p>
            </div>
          </section>

          {/* ── Section 9: Example Project ── */}
          <section className="mb-14">
            <h2 className="text-2xl font-semibold tracking-tight text-foreground mb-4">
              Example Project
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              The same Job Portal application built with double architecture. The web app includes
              role-protected admin routes for managing job listings, companies, and applications.
              Compare it with the triple example to see the structural differences.
            </p>
            <div className="rounded-lg border border-border/40 bg-accent/20 p-6">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-base font-semibold text-foreground mb-1">Job Portal — Double + TanStack Router</h3>
                  <p className="text-sm text-muted-foreground">
                    Full source code with README, .env template, step-by-step guide, and production Docker Compose.
                  </p>
                </div>
                <a
                  href="https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-double-vite"
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
              <Link href="/docs/concepts/architecture-modes/triple" className="gap-1.5">
                <ArrowLeft className="h-3.5 w-3.5" />
                Triple Architecture
              </Link>
            </Button>
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/concepts/architecture-modes/single" className="gap-1.5">
                Single Architecture
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
