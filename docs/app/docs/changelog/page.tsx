import Link from 'next/link'
import { ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/changelog')

export default function ChangelogPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Release History</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Changelog
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                All notable changes to Grit are documented here. Each release includes new features,
                bug fixes, and any breaking changes you need to be aware of.
              </p>
            </div>

            {/* v2.3.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v2.3.0
                </span>
                <span className="text-sm text-muted-foreground">March 5, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Fixes</h3>
                <ul>
                  <li>
                    <strong>Desktop build script</strong> &mdash; Removed <code>tsc</code> from the frontend
                    build script. TanStack Router{`'`}s Vite plugin generates <code>routeTree.gen.ts</code> during
                    the Vite build, so running <code>tsc</code> before Vite caused{' '}
                    <code>Cannot find module {`'`}./routeTree.gen{`'`}</code> errors. Vite handles type-stripping
                    via esbuild.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.2.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.2.0
                </span>
                <span className="text-sm text-muted-foreground">March 4, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Fixes</h3>
                <ul>
                  <li>
                    <strong>Desktop auth hook file extension</strong> &mdash; Renamed the scaffolded{' '}
                    <code>use-auth.ts</code> to <code>use-auth.tsx</code> so TypeScript correctly handles
                    the JSX in <code>&lt;AuthContext.Provider&gt;</code>. Previously, <code>grit new-desktop</code>{' '}
                    projects would fail to compile with <code>TS1005: {`'>'`} expected</code> errors.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    Added Desktop Handbook PDF download links to all 8 desktop documentation pages.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.1.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.1.0
                </span>
                <span className="text-sm text-muted-foreground">March 4, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>TanStack Router for desktop</strong> &mdash; Migrated the desktop frontend from
                    React Router to{' '}
                    <a href="https://tanstack.com/router" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">TanStack Router</a>{' '}
                    with file-based routing. Routes are auto-discovered by the Vite plugin &mdash; no centralized
                    route registry. Uses <code>createHashHistory()</code> for Wails compatibility and{' '}
                    <code>Route.useParams()</code> for type-safe params. Resource generation now creates 5 files
                    (list, new, edit routes + model + service) and performs 10 injections (down from 12).
                  </li>
                  <li>
                    <strong>Mobile navigation</strong> &mdash; Added a hamburger menu to the docs site header,
                    visible below the <code>lg</code> breakpoint. Opens a Sheet sidebar with all navigation links.
                    Auto-closes on link click.
                  </li>
                  <li>
                    <strong>CGO-free SQLite</strong> &mdash; Replaced <code>gorm.io/driver/sqlite</code> (requires
                    CGO) with <code>github.com/glebarez/sqlite</code> (pure Go) in all scaffold templates. Desktop
                    apps now build and run without CGO or a C compiler.
                  </li>
                  <li>
                    <strong>20 Desktop Project Ideas</strong> &mdash; New{' '}
                    <Link href="/docs/desktop/project-ideas" className="text-primary hover:underline">project ideas page</Link>{' '}
                    with 20 ready-to-build desktop app ideas across business, education, healthcare, logistics,
                    and more. Each includes resources, field definitions, and <code>grit generate</code> commands.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    Added TanStack Router explanations to all desktop doc pages: overview, getting started,
                    first app, resource generation, and POS app.
                  </li>
                  <li>
                    Updated{' '}
                    <Link href="/docs/desktop/llm-reference" className="text-primary hover:underline">LLM Reference</Link>,{' '}
                    <Link href="/docs/ai-skill" className="text-primary hover:underline">GRIT_SKILL.md</Link>, and
                    database docs to reflect TanStack Router and CGO-free SQLite changes.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.0.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.0.0
                </span>
                <span className="text-sm text-muted-foreground">March 4, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Native desktop apps (Wails)</strong> &mdash; New <code>grit new-desktop</code> command
                    scaffolds a complete desktop application with Go backend, React frontend (Vite + TanStack Router +
                    TanStack Query), SQLite database, JWT authentication, blog and contact CRUD, PDF/Excel export,
                    custom title bar, dark theme, and GORM Studio. Compiles to a single native executable for
                    Windows, macOS, and Linux. See{' '}
                    <Link href="/docs/desktop" className="text-primary hover:underline">Desktop docs</Link>.
                  </li>
                  <li>
                    <strong>Desktop resource generation</strong> &mdash; <code>grit generate resource</code> now
                    works inside desktop projects. Generates Go model, service, and TanStack Router route files
                    (list, new, edit), then injects code into 10 locations (db.go, main.go, app.go, types.go,
                    sidebar.tsx, studio/main.go) using <code>grit:</code> markers. See{' '}
                    <Link href="/docs/desktop/resource-generation" className="text-primary hover:underline">Desktop Resource Generation</Link>.
                  </li>
                  <li>
                    <strong>Project type auto-detection</strong> &mdash; All CLI commands now auto-detect whether
                    you are inside a web (Turborepo) or desktop (Wails) project. No flags needed.
                  </li>
                  <li>
                    <strong><code>grit start</code> for desktop</strong> &mdash; Running <code>grit start</code>{' '}
                    inside a desktop project launches <code>wails dev</code> with hot-reload for both Go and React.
                  </li>
                  <li>
                    <strong><code>grit compile</code></strong> &mdash; New command that runs <code>wails build</code>{' '}
                    to produce a distributable native binary.
                  </li>
                  <li>
                    <strong><code>grit studio</code></strong> &mdash; New command that launches GORM Studio. For
                    desktop projects it starts a standalone server on port 4000. For web projects it opens the
                    browser to the embedded Studio route.
                  </li>
                  <li>
                    <strong><code>grit remove resource</code> for desktop</strong> &mdash; Removes a previously
                    generated desktop resource, deleting files and reversing all 10 marker injections.
                  </li>
                  <li>
                    <strong>Grit UI component registry (91 components)</strong> &mdash; Every scaffolded web
                    project now includes a shadcn-compatible component registry with 91 pre-built components
                    across 5 categories: marketing (14), auth (10), SaaS (30), ecommerce (20), and layout (18).
                    Install via <code>npx shadcn@latest add</code> from <code>/r</code> endpoints.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/desktop" className="text-primary hover:underline">Desktop (Wails)</Link>{' '}
                    section &mdash; 8 pages covering overview, getting started, first app tutorial, POS app
                    tutorial, resource generation, building/distribution, project ideas, and LLM reference.
                  </li>
                  <li>
                    Updated{' '}
                    <Link href="/docs/ai-skill/llm-guide" className="text-primary hover:underline">LLM Reference</Link>{' '}
                    with complete desktop section: project structure, CLI commands, markers, and architecture comparison.
                  </li>
                </ul>
              </div>
            </div>

            {/* v1.4.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v1.4.0
                </span>
                <span className="text-sm text-muted-foreground">March 2, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Gzip response compression</strong> &mdash; All API responses are now
                    compressed automatically via a custom <code>Gzip()</code> middleware using the Go
                    standard library <code>compress/gzip</code> at <code>BestSpeed</code>.
                    JSON payloads shrink by 60–80%, reducing bandwidth on paginated list endpoints
                    with zero external dependencies.
                  </li>
                  <li>
                    <strong>Request ID tracing</strong> &mdash; A <code>RequestID()</code> middleware
                    injects a unique <code>X-Request-ID</code> header on every request (echoes the
                    upstream header or generates a nanosecond-based ID). The ID is stored in Gin
                    context and included in every structured log line for end-to-end request tracing.
                  </li>
                  <li>
                    <strong>Database connection pool tuning</strong> &mdash; The scaffold now sets
                    four GORM pool parameters: <code>MaxIdleConns(10)</code>,{' '}
                    <code>MaxOpenConns(100)</code>, <code>ConnMaxLifetime(30m)</code>, and{' '}
                    <code>ConnMaxIdleTime(10m)</code>. This prevents stale connections after network
                    interruptions and avoids connection exhaustion under load.
                  </li>
                  <li>
                    <strong>Cache-Control headers on public blog endpoints</strong> &mdash; The{' '}
                    <code>ListPublished</code> handler now returns{' '}
                    <code>Cache-Control: public, max-age=300</code> (5 minutes) and{' '}
                    <code>GetBySlug</code> returns <code>Cache-Control: public, max-age=3600</code>{' '}
                    (1 hour). CDNs and edge caches can now serve public blog content without hitting
                    the Go API.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/concepts/performance" className="text-primary hover:underline">Performance</Link>{' '}
                    page &mdash; comprehensive guide to all backend (Go/API) and frontend (Next.js)
                    performance optimisations that ship with every Grit project out of the box.
                    Covers Gzip, Request ID, connection pool, Cache-Control, presigned uploads,
                    background jobs, Redis caching, Server Components, ISR, React Query, next/image,
                    Turborepo, and code splitting.
                  </li>
                  <li>
                    New{' '}
                    <Link href="/docs/ai-skill/llm-guide" className="text-primary hover:underline">Complete LLM Reference</Link>{' '}
                    page &mdash; a dedicated machine-readable guide that teaches AI assistants
                    everything about Grit: project structure, all CLI commands, every field type,
                    code patterns, API response format, code markers, naming conventions, all
                    batteries, performance features, and the golden rules that must never be broken.
                  </li>
                </ul>
              </div>
            </div>

            {/* v1.3.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v1.3.0
                </span>
                <span className="text-sm text-muted-foreground">February 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Presigned URL uploads</strong> &mdash; File uploads now bypass the API server entirely.
                    The browser gets a presigned PUT URL, uploads directly to S3/R2/MinIO, then records the upload
                    in the database. This fixes file uploads breaking behind reverse proxies (Dokploy/Traefik/Nginx)
                    due to request body size limits and timeouts. Includes progress tracking via XHR.
                  </li>
                  <li>
                    <strong>Error pages for scaffolded apps</strong> &mdash; New <code>grit new</code> projects now include{' '}
                    <code>error.tsx</code>, <code>not-found.tsx</code>, and <code>global-error.tsx</code> for both
                    admin and web apps. Errors are displayed with styled UI instead of the default Next.js error page.
                  </li>
                  <li>
                    <strong>Production-ready Docker config</strong> &mdash; <code>docker-compose.prod.yml</code> now uses{' '}
                    <code>expose</code> instead of <code>ports</code>, <code>env_file</code> for secrets, MinIO service,
                    named bridge network, build args for <code>NEXT_PUBLIC_API_URL</code>, and Go 1.24.
                  </li>
                  <li>
                    <strong>Sentinel ExcludePaths</strong> &mdash; Pulse, GORM Studio, Sentinel, and API docs paths are
                    now excluded from rate limiting by default, fixing Pulse health checks triggering rate limits.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/getting-started/create-without-docker" className="text-primary hover:underline">Create without Docker</Link>{' '}
                    guide &mdash; set up a Grit project using Neon, Upstash, Cloudflare R2, and Resend instead of Docker.
                  </li>
                </ul>

                <h3>Infrastructure</h3>
                <ul>
                  <li>
                    Scaffold Dockerfile updated from Go 1.23 to Go 1.24
                  </li>
                  <li>
                    Next.js Dockerfile now accepts <code>NEXT_PUBLIC_API_URL</code> as a build argument
                  </li>
                  <li>
                    <code>.env</code> template includes Docker Compose production variables (<code>POSTGRES_USER</code>,{' '}
                    <code>POSTGRES_PASSWORD</code>, <code>POSTGRES_DB</code>, <code>API_URL</code>)
                  </li>
                </ul>
              </div>
            </div>

            {/* v1.1.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v1.1.0
                </span>
                <span className="text-sm text-muted-foreground">February 25, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Default font changed to Onest</strong> &mdash; New projects scaffolded with{' '}
                    <code>grit new</code> now use the{' '}
                    <a href="https://fonts.google.com/specimen/Onest" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">Onest</a>{' '}
                    Google Font for all UI text instead of DM Sans. JetBrains Mono remains the code font.
                    The font is loaded via <code>next/font/google</code> with weights 400, 500, 600, and 700.
                  </li>
                  <li>
                    <strong>Hire Us page</strong> &mdash; New{' '}
                    <Link href="/hire" className="text-primary hover:underline">/hire</Link>{' '}
                    page for professional Grit development services. Includes service offerings,
                    tech stack overview, and contact CTA.
                  </li>
                  <li>
                    <strong>Monetization banners</strong> &mdash; Docs sidebar now shows promotional cards for{' '}
                    <a href="https://gritcms.com" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">GritCMS</a>,
                    developer hiring services, and{' '}
                    <Link href="/donate" className="text-primary hover:underline">donations</Link>{' '}
                    &mdash; visible on every documentation page.
                  </li>
                  <li>
                    <strong>Grit Fullstack Course page</strong> &mdash; New{' '}
                    <Link href="/course" className="text-primary hover:underline">/course</Link>{' '}
                    page with a 10-module curriculum covering Go, React, Next.js, and the full Grit stack.
                  </li>
                </ul>

                <h3>Improvements</h3>
                <ul>
                  <li>
                    Top navigation now includes GritCMS, Hire Us, and a Sponsor heart icon for quick access
                    to all revenue channels.
                  </li>
                  <li>
                    <code>richtext</code> added to the FieldType union for better type safety in the code generator.
                  </li>
                </ul>

                <h3>Bug Fixes</h3>
                <ul>
                  <li>
                    <strong>OAuth callback fix</strong> &mdash; Fixed <code>TokenPair</code> struct field access
                    in the social login callback handler (was using map indexing instead of struct fields).
                  </li>
                  <li>
                    <strong>Course waitlist fix</strong> &mdash; Fixed Google Sheets submission to use
                    form-encoded data instead of JSON.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/getting-started/cli-cheatsheet" className="text-primary hover:underline">CLI Cheatsheet</Link>{' '}
                    page &mdash; complete reference for all Grit CLI commands with flags, field types,
                    generated files, common workflows, and full command tree.
                  </li>
                  <li>
                    New{' '}
                    <Link href="/docs/backend/oauth" className="text-primary hover:underline">Social Login (OAuth2)</Link>{' '}
                    setup guide for Google and GitHub authentication.
                  </li>
                  <li>
                    Updated Docker Cheat Sheet with force remove commands for containers and volumes.
                  </li>
                  <li>
                    Updated AI skill guide with social login (OAuth2) section.
                  </li>
                </ul>
              </div>
            </div>

            {/* v1.0.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v1.0.0
                </span>
                <span className="text-sm text-muted-foreground">February 24, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Social Login (Google + GitHub)</strong> &mdash; Every <code>grit new</code> project now
                    includes OAuth2 social authentication via{' '}
                    <a href="https://github.com/markbates/goth" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">Gothic</a>.
                    Users can sign in with Google or GitHub on all auth pages (login, register, admin).
                    Accounts are linked by email &mdash; existing users who sign in with a social provider are automatically connected.
                    Configurable via <code>GOOGLE_CLIENT_ID</code>, <code>GITHUB_CLIENT_ID</code> environment variables.
                  </li>
                  <li>
                    <strong>GORM Studio v1.0.1</strong> &mdash; Updated to the first stable tagged release of GORM Studio.
                  </li>
                </ul>

                <h3>Improvements</h3>
                <ul>
                  <li>
                    User model now includes <code>Provider</code>, <code>GoogleID</code>, and <code>GithubID</code> fields
                    for social account linking. Password field is now nullable to support OAuth-only accounts.
                  </li>
                  <li>
                    Admin users table shows Provider column with badges (Email, Google, GitHub) and new filter option.
                  </li>
                  <li>
                    Social login buttons (Google + GitHub) appear on all 4 admin style variants (default, modern, minimal, glass).
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.19.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.19.0
                </span>
                <span className="text-sm text-muted-foreground">February 24, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Fixes</h3>
                <ul>
                  <li>
                    <strong>gin-docs AuthConfig</strong> &mdash; Updated scaffold template to use the new{' '}
                    <code>gindocs.AuthConfig</code> struct instead of the deprecated <code>gindocs.AuthBearer</code> constant,
                    fixing compilation errors in newly scaffolded projects.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/tutorials/contact-app" className="text-primary hover:underline">Your First App</Link>{' '}
                    tutorial &mdash; step-by-step Contact Manager guide covering project setup, resource generation, and CRUD
                  </li>
                  <li>
                    New{' '}
                    <Link href="/docs/deployment/dokploy" className="text-primary hover:underline">Dokploy Deployment</Link>{' '}
                    guide with Dockerfile examples
                  </li>
                  <li>
                    Improved terminal blocks across all tutorials with copy buttons and horizontal scroll
                  </li>
                  <li>
                    Updated{' '}
                    <Link href="/docs/backend/api-docs" className="text-primary hover:underline">API Documentation</Link>{' '}
                    page to reflect the new <code>AuthConfig</code> struct format
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.18.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.18.0
                </span>
                <span className="text-sm text-muted-foreground">February 22, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Pulse (Observability)</strong> &mdash; Every <code>grit new</code> project now includes{' '}
                    <a href="https://github.com/MUKE-coder/pulse" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">Pulse</a>,
                    a self-hosted observability SDK. Provides request tracing, database monitoring, runtime metrics,
                    error tracking, health checks, alerting, Prometheus export, and an embedded React dashboard
                    at <code>/pulse</code>. Enabled by default, configurable via <code>PULSE_ENABLED</code>.
                    See <Link href="/docs/backend/pulse" className="text-primary hover:underline">Pulse docs</Link>.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/backend/pulse" className="text-primary hover:underline">Pulse (Observability)</Link> page
                    covering configuration, endpoints, health checks, alerting, Prometheus metrics, and data storage
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.17.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.17.0
                </span>
                <span className="text-sm text-muted-foreground">February 22, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>API Documentation (gin-docs)</strong> &mdash; Replaced hand-written Scalar/OpenAPI
                    spec with{' '}
                    <a href="https://github.com/MUKE-coder/gin-docs" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">gin-docs</a>,
                    a zero-annotation API documentation generator. Routes and GORM models are introspected
                    automatically to produce an OpenAPI 3.1 spec with interactive Scalar or Swagger UI,
                    plus Postman and Insomnia export.
                  </li>
                  <li>
                    <strong>Dark/Light mode for Go Playground</strong> &mdash; The playground now follows the
                    site-wide theme toggle, switching between VS Code dark and light CodeMirror themes.
                  </li>
                  <li>
                    <strong>Umami Analytics</strong> &mdash; Optional visitor analytics via self-hosted Umami,
                    configured with <code>NEXT_PUBLIC_UMAMI_WEBSITE_ID</code> environment variable.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/backend/api-docs" className="text-primary hover:underline">API Documentation</Link> page
                    covering gin-docs configuration, GORM model schemas, route customization, UI switching, and spec export
                  </li>
                  <li>Full SEO + AEO implementation: sitemap, robots.txt, JSON-LD structured data, per-page metadata</li>
                </ul>

                <h3>Infrastructure</h3>
                <ul>
                  <li>Added Dockerfile for docs site deployment (Next.js standalone output)</li>
                  <li>Google Search Console verification</li>
                </ul>
              </div>
            </div>

            {/* v0.16.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.16.0
                </span>
                <span className="text-sm text-muted-foreground">February 21, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Go Playground</strong> &mdash; Interactive code editor at{' '}
                    <Link href="/playground" className="text-primary hover:underline">/playground</Link> with
                    Go syntax highlighting, code execution via the official Go Playground API, example snippets,
                    share links, and keyboard shortcuts (Ctrl+Enter to run).
                  </li>
                  <li>
                    <strong>GORM Studio updated</strong> &mdash; Updated to latest version with raw SQL editor,
                    schema export (SQL/JSON/YAML/DBML/ERD), data import/export (JSON/CSV/SQL/XLSX),
                    and Go model generation from database schema.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    <Link href="/docs/prerequisites/golang" className="text-primary hover:underline">Go for Grit Developers</Link> &mdash;
                    comprehensive rewrite with 22 sections covering methods, Gin routing, middleware, CORS,
                    handler/service architecture, GORM CRUD, migrations, seeding, JWT auth flow, and RBAC
                  </li>
                  <li>Fixed right-side table of contents for the Go prerequisites page</li>
                  <li>New Middleware and CORS sections added to Go guide</li>
                </ul>
              </div>
            </div>

            {/* v0.15.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.15.0
                </span>
                <span className="text-sm text-muted-foreground">February 20, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Security (Sentinel)</strong> &mdash; Every <code>grit new</code> project now ships with
                    a production-grade security suite powered by{' '}
                    <Link href="https://github.com/MUKE-coder/sentinel" className="text-primary hover:underline">Sentinel</Link>.
                    Includes WAF, rate limiting, brute-force protection, anomaly detection, IP geolocation,
                    security headers, and a real-time threat dashboard at <code>/sentinel/ui</code>.
                    See <Link href="/docs/batteries/security" className="text-primary hover:underline">Security docs</Link>.
                  </li>
                  <li>
                    <strong>Admin security page</strong> &mdash; New System &rarr; Security page in the admin panel
                    embeds the Sentinel dashboard for monitoring threats without leaving the admin UI.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>New: <Link href="/docs/batteries/security" className="text-primary hover:underline">Security (Sentinel)</Link> documentation page</li>
                  <li>Migrated getting-started pages (Installation, Quick Start, Troubleshooting) to use CodeBlock component</li>
                  <li>Added prerequisite learning pages for Go, Next.js, and Docker</li>
                </ul>
              </div>
            </div>

            {/* v0.14.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.14.0
                </span>
                <span className="text-sm text-muted-foreground">February 18, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Multi-step forms</strong> &mdash; New <code>formView: &quot;modal-steps&quot;</code> and{' '}
                    <code>&quot;page-steps&quot;</code> variants with horizontal/vertical step indicators,
                    per-step validation, progress bar, and clickable step navigation.
                    See <Link href="/docs/admin/multi-step-forms" className="text-primary hover:underline">Multi-Step Forms</Link>.
                  </li>
                  <li>
                    <strong>Standalone component usage</strong> &mdash; FormBuilder, FormStepper, and DataTable
                    can now be used on any page in both web and admin apps without the resource system.
                    See <Link href="/docs/admin/standalone-usage" className="text-primary hover:underline">Standalone Usage</Link>.
                  </li>
                  <li>
                    <strong>Richtext field type</strong> &mdash; New <code>richtext</code> field with Tiptap WYSIWYG
                    editor (bold, italic, headings, lists, code blocks, links, undo/redo).
                  </li>
                  <li>
                    <strong><code>string_array</code> field type</strong> &mdash; Store arrays of strings
                    using <code>datatypes.JSONSlice[string]</code>. Works with PostgreSQL and SQLite.
                    Maps to <code>string[]</code> in TypeScript and <code>z.array(z.string())</code> in Zod.
                  </li>
                  <li>
                    <strong>Built-in blog example</strong> &mdash; <code>grit new</code> now scaffolds a complete
                    blog with model, service, handler, seed data, public web pages, and admin resource definition.
                  </li>
                  <li>
                    <strong>Sidebar user avatar</strong> &mdash; Admin sidebar shows the current user&apos;s avatar
                    with a dropdown menu for profile and logout.
                  </li>
                  <li>
                    <strong>Profile avatar upload</strong> &mdash; Profile page now supports avatar image upload.
                  </li>
                  <li>
                    <strong><code>react-hook-form</code> in web app</strong> &mdash; Web app scaffold now includes{' '}
                    <code>react-hook-form</code> as a dependency, enabling standalone FormBuilder usage out of the box.
                  </li>
                </ul>

                <h3>Bug Fixes</h3>
                <ul>
                  <li>
                    <strong>Scalar API docs crash</strong> &mdash; Fixed <code>c.String</code> treating HTML as
                    a format string. Now uses <code>c.Data</code> to avoid panics when Scalar HTML
                    contains <code>%</code> characters in CSS/JS.
                  </li>
                  <li>
                    <strong>Blog route conflict</strong> &mdash; Admin blog CRUD routes moved
                    from <code>/api/blogs</code> to <code>/api/admin/blogs</code> to avoid conflict
                    with public blog routes.
                  </li>
                  <li>
                    <strong>Select dropdown styling</strong> &mdash; Fixed relationship select dropdown
                    rendering behind modals using portal-based positioning.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>New: <Link href="/docs/tutorials/product-catalog" className="text-primary hover:underline">Build a Product Catalog</Link> tutorial &mdash; resource generation, multi-step forms, standalone DataTable &amp; FormBuilder</li>
                  <li>New: <Link href="/docs/admin/multi-step-forms" className="text-primary hover:underline">Multi-Step Forms</Link> guide</li>
                  <li>New: <Link href="/docs/admin/standalone-usage" className="text-primary hover:underline">Standalone Usage</Link> guide</li>
                  <li>New: Changelog page</li>
                  <li>Updated CLI Commands, Code Generation, Quick Start, Resources, Shared Package, Web App, Seeders, and Forms pages</li>
                </ul>
              </div>
            </div>

            {/* v0.12.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.12.0
                </span>
                <span className="text-sm text-muted-foreground">February 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Relationship support</strong> &mdash; New <code>belongs_to</code> and{' '}
                    <code>many_to_many</code> field types for the code generator. Automatically creates
                    foreign keys, junction tables, and relationship-aware form fields.
                  </li>
                  <li>
                    <strong>Relationship select fields</strong> &mdash; New <code>relationship-select</code> and{' '}
                    <code>multi-relationship-select</code> form field components with search, portal-based
                    dropdowns, and tag-based multi-select.
                  </li>
                  <li>
                    <strong>Beginner tutorial</strong> &mdash; &quot;Learn Grit Step by Step&quot; tutorial
                    walking through building a full-stack app from scratch.
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.11.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.11.0
                </span>
                <span className="text-sm text-muted-foreground">February 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Full-page form view</strong> &mdash; New <code>formView: &quot;page&quot;</code> option
                    renders forms as dedicated pages instead of modals.
                  </li>
                  <li>
                    <strong><code>slug</code> field type</strong> &mdash; Auto-generates URL-friendly slugs with
                    unique suffixes. Excluded from create/update forms and Zod schemas.
                  </li>
                  <li>
                    <strong>DataTable column customization</strong> &mdash; Hide/show columns, column visibility
                    toggle in table toolbar.
                  </li>
                  <li>
                    <strong><code>grit start</code> commands</strong> &mdash; <code>grit start client</code> and{' '}
                    <code>grit start server</code> for running frontend and API separately.
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.10.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.10.0
                </span>
                <span className="text-sm text-muted-foreground">January 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Style variants</strong> &mdash; <code>--style</code> flag for <code>grit new</code> with
                    4 admin panel styles: default, modern, minimal, and glass.
                  </li>
                  <li>
                    <strong>Air hot reloading</strong> &mdash; Go API development with automatic rebuild on file
                    changes using Air.
                  </li>
                  <li>
                    <strong><code>grit remove resource</code></strong> &mdash; Remove a generated resource and
                    clean up all injected code (model, handler, routes, schemas, types, hooks, admin pages).
                  </li>
                  <li>
                    <strong>AI workflow docs</strong> &mdash; Guides for using Grit with Claude and Antigravity AI assistants.
                  </li>
                </ul>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Introduction
                </Link>
              </Button>
              <div />
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
