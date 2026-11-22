import Link from "next/link";
import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { DocsShell } from "@/components/docs/shell";
import { DocsPageHeader } from "@/components/docs/page-header";
import { CodeBlock } from "@/components/code-block";
import { Steps, Step } from "@/components/steps";
import { CommunityCTA } from "@/components/community-cta";
import { getDocMetadata } from '@/config/docs-metadata'
import { GO_MIN } from '@/config/site'
import { BookDownload } from '@/components/book-download'
import { DocsShot } from '@/components/docs-shot'
import { GRIT_VERSION } from '@/config/site'

export const metadata = getDocMetadata('/docs/getting-started/quick-start')

export default function QuickStartPage() {
  return (
    <DocsShell
      header={
        <DocsPageHeader
          crumbs={[
            { label: 'Docs', href: '/docs' },
            { label: 'Getting started', href: '/docs' },
            { label: 'Quick Start' },
          ]}
          title="Quick Start"
          readingTime={5}
          sourcePath="app/docs/getting-started/quick-start/page.tsx"
          lead={
            <>
              Get a full-stack Grit project running in five minutes. This guide
              assumes Go, Node.js, pnpm and Docker are installed. If not, start
              with{' '}
              <Link
                href="/docs/getting-started/installation"
                className="g-link"
              >
                Installation
              </Link>
              .
            </>
          }
        />
      }
    >
      <BookDownload className="mb-8" />

            {/* Prerequisites */}
            <div className="prose-grit mb-10">
              <h2>Prerequisites</h2>
              <p>
                Make sure the following tools are installed on your machine
                before proceeding:
              </p>
            </div>

            <div className="grid gap-3 sm:grid-cols-2 mb-10">
              {[
                { name: "Go", version: `${GO_MIN}+`, check: "go version" },
                { name: "Node.js", version: "22+", check: "node --version" },
                { name: "pnpm", version: "9+", check: "pnpm --version" },
                {
                  name: "Docker",
                  version: "Latest",
                  check: "docker --version",
                },
              ].map((tool) => (
                <div
                  key={tool.name}
                  className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                >
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-[15px] font-semibold">{tool.name}</span>
                    <span className="text-sm font-mono text-accent/60">
                      {tool.version}
                    </span>
                  </div>
                  <code className="text-sm font-mono text-muted-foreground/50">
                    {tool.check}
                  </code>
                </div>
              ))}
            </div>

            <Steps>
            {/* Step 1 */}
            <Step title="Install the Grit CLI">
              <div className="prose-grit mb-4">
                <p>
                  One line, works on every platform. The script detects an
                  existing install and runs <code>grit update</code>; otherwise it
                  downloads the matching release binary for your OS / arch.
                </p>
              </div>
              <CodeBlock
                terminal
                code={`# macOS / Linux
curl -fsSL https://gritframework.dev/install.sh | sh

# Windows (PowerShell)
iwr -useb https://gritframework.dev/install.ps1 | iex`}
                className="mb-0 glow-purple-sm"
              />
              <div className="prose-grit mt-4">
                <p>
                  Verify the installation by running <code>grit --help</code>.
                  You should see the Grit ASCII art logo and a list of available
                  commands.
                </p>
              </div>

              <div className="mt-6 space-y-4">
                <h3 className="text-base font-semibold text-foreground/90">
                  Update to the latest version
                </h3>
                <p className="text-[15px] text-muted-foreground">
                  Already have Grit installed? Update the CLI to the latest version
                  with a single command. It checks for a newer release and, if there is
                  one, swaps in the new binary:
                </p>
                <CodeBlock terminal code="grit update" className="mb-0" />
                <p className="text-[15px] text-muted-foreground">
                  Run <code className="text-[15px] font-mono bg-surface-2 px-1.5 py-0.5 rounded-xs text-accent">grit version</code> afterwards
                  to confirm you&apos;re on the latest release.
                </p>
              </div>
            </Step>

            {/* Step 2 */}
            <Step title="Create a New Project">
              <div className="prose-grit mb-4">
                <p>
                  Scaffold a complete full-stack project with one command. This
                  creates the entire monorepo: Go API, Next.js web app, admin
                  panel, shared types, Docker setup, and all the batteries.
                </p>
              </div>
              <CodeBlock terminal code="grit new myapp" className="mb-0 glow-purple-sm" />

              {/* Captured from a real run. A reader who cannot tell a slow step
                  from a failed one has no way to know whether to keep waiting. */}
              <details className="mt-4 rounded-xl border border-border/60 bg-card/30 px-5 py-4">
                <summary className="cursor-pointer text-[14px] font-medium">
                  What it prints
                </summary>
                <div className="mt-3">
                  <CodeBlock
                    output="grit new myapp"
                    language="text"
                    code={`  ✓ Directory structure
  ✓ Configuration files
  ✓ Go API                                        391ms
  ✓ Migration and seed tools
  ✓ Batteries                                     721ms
  ✓ Blog example
  ✓ Go dependencies                               2.3s
  ✓ Docker setup
  ✓ Shared package
  ✓ Web app                                       106ms
  ✓ Admin panel                                   748ms
  ✓ Frontend tests

  Created myapp  Triple · Atlas · postgres · 4.4s

  Next steps

    cd myapp
    docker compose up -d    Postgres, Redis, MinIO, Mailhog
    pnpm install            frontend deps, one time
    grit migrate            create database tables
    grit seed               sample data, optional
    grit start              run everything`}
                    className="mb-0"
                  />
                  <p className="text-[13px] text-muted-foreground mt-3 leading-relaxed">
                    Then a table of every address the project will serve on. The timings are
                    from one machine and yours will differ; what matters is that every line
                    is a tick and the last one says <strong>Created</strong>.
                  </p>
                </div>
              </details>

              <div className="prose-grit mt-4">
                <p>
                  The <strong>interactive CLI</strong> walks you through selecting an architecture mode
                  and frontend option. Grit supports <strong>6 architecture modes</strong>:{' '}
                  <code>full</code> (everything), <code>triple</code> (API + web + admin),{' '}
                  <code>double</code> (API + web), <code>single</code> (one binary: Go API with
                  the React app embedded), <code>api</code> (headless, no frontend), and{' '}
                  <code>mobile</code> (API + Expo). It also supports <strong>2 frontend options</strong>:{' '}
                  <code>Next.js</code> and <code>TanStack Router (Vite)</code>.
                </p>
                <p>
                  The project name must be lowercase, alphanumeric, and hyphens
                  only (e.g., <code>my-saas-app</code>). It must start with a
                  letter and cannot end with a hyphen.
                </p>
              </div>

              {/* Five questions is four more decisions than a newcomer has any
                  basis for making, so say what Enter does. The CLI prints the
                  same line above its first question (grit#240). */}
              <div className="mt-4 rounded-xl border border-accent/25 bg-accent/[0.04] p-5">
                <p className="text-[15px] leading-relaxed">
                  <strong>New to Grit? Press Enter on every question.</strong> Every picker opens
                  on the recommended answer, so pressing Enter through the prompt gives you the
                  same project as{' '}
                  <code className="font-mono text-[14px] text-accent">grit new myapp --triple --next</code>{' '}
                  with Postgres and MinIO: the combination every guide here assumes. The CLI
                  prints this line above the first question too.
                </p>
              </div>

              <div className="mt-6 space-y-4">
                <h3 className="text-base font-semibold text-foreground/90">
                  Flag shortcuts (skip prompts)
                </h3>
                <p className="text-[15px] text-muted-foreground">
                  Use flags to skip the interactive prompts and scaffold instantly:
                </p>
                <CodeBlock terminal code="grit new myapp --triple --vite" className="mb-0" />
                <p className="text-[15px] text-muted-foreground mt-3">
                  Available flags: <code className="text-[15px] font-mono bg-surface-2 px-1.5 py-0.5 rounded-xs text-accent">--single</code>,{' '}
                  <code className="text-[15px] font-mono bg-surface-2 px-1.5 py-0.5 rounded-xs text-accent">--double</code>,{' '}
                  <code className="text-[15px] font-mono bg-surface-2 px-1.5 py-0.5 rounded-xs text-accent">--triple</code>,{' '}
                  <code className="text-[15px] font-mono bg-surface-2 px-1.5 py-0.5 rounded-xs text-accent">--api</code>,{' '}
                  <code className="text-[15px] font-mono bg-surface-2 px-1.5 py-0.5 rounded-xs text-accent">--mobile</code> for architecture,{' '}
                  <code className="text-[15px] font-mono bg-surface-2 px-1.5 py-0.5 rounded-xs text-accent">--vite</code>,{' '}
                  <code className="text-[15px] font-mono bg-surface-2 px-1.5 py-0.5 rounded-xs text-accent">--next</code> for frontend,
                  and <code className="text-[15px] font-mono bg-surface-2 px-1.5 py-0.5 rounded-xs text-accent">--desktop</code> to add a
                  Wails desktop client (<a href="/docs/concepts/architecture-modes/multi-client" className="text-accent hover:underline">multi-client pattern</a>).
                </p>
              </div>

              <div className="mt-6 space-y-4">
                <h3 className="text-base font-semibold text-foreground/90">
                  In-place scaffolding
                </h3>
                <p className="text-[15px] text-muted-foreground">
                  Scaffold directly into the current directory instead of creating a new folder:
                </p>
                <CodeBlock terminal code="grit new . --triple --vite" className="mb-0" />
                <p className="text-sm text-muted-foreground/60 mt-2">
                  Use <code className="text-[15px] font-mono bg-surface-2 px-1.5 py-0.5 rounded-xs text-accent">--force</code> to
                  scaffold into a non-empty directory.
                </p>
              </div>
            </Step>

            {/* Step 3 */}
            <Step title="Start Infrastructure Services">
              <div className="prose-grit mb-4">
                <p>
                  Navigate into the project and start the Docker services. This
                  launches PostgreSQL, Redis, MinIO (S3-compatible storage), and
                  Mailhog (email testing).
                </p>
              </div>
              <CodeBlock terminal code={`cd myapp
docker compose up -d`} className="mb-0 glow-purple-sm" />
              <div className="prose-grit mt-4">
                <p>This starts the following services in the background:</p>
                <ul>
                  <li>
                    <strong>PostgreSQL 16</strong> on port 5434 -- your primary
                    database
                  </li>
                  <li>
                    <strong>Redis 7</strong> on port 6380 -- caching and job
                    queues
                  </li>
                  <li>
                    <strong>MinIO</strong> on port 9002 (console: 9003) -- local
                    S3-compatible file storage
                  </li>
                  <li>
                    <strong>Mailhog</strong> on port 8025 -- catch-all email
                    testing UI
                  </li>
                </ul>
                <blockquote>
                  Do not have Docker? See the{" "}
                  <Link
                    href="/docs/getting-started/installation"
                    className="text-accent hover:underline"
                  >
                    Installation guide
                  </Link>{" "}
                  for a cloud-only setup using Neon (Postgres) and Upstash
                  (Redis) instead.
                </blockquote>
              </div>
            </Step>

            {/* Step 4 */}
            <Step title="Install Deps & Set Up the Database">
              <div className="prose-grit mb-4">
                <p>
                  Install the frontend dependencies once, then create your database
                  tables and load some sample data &mdash; all from the project root.
                  You never <code>cd</code> into a sub-folder or touch <code>go</code> by hand.
                </p>
              </div>
              <CodeBlock terminal code={`pnpm install      # install frontend deps (one-time)
grit migrate      # create database tables from your models
grit seed         # (optional) add a demo admin + sample rows`} className="mb-0 glow-purple-sm" />
              <div className="prose-grit mt-4">
                <p>
                  <code>grit migrate</code> runs GORM AutoMigrate for every registered
                  model and prints a created / altered / unchanged summary.{" "}
                  <code>grit seed</code> fills the tables with starter data, including a
                  demo admin you can log in with.
                </p>
                <blockquote>
                  Migrations are an <strong>explicit command</strong>, not something that
                  runs on server start &mdash; so you always control when your schema
                  changes. Run <code>grit migrate</code> again any time you add or change a
                  model.
                </blockquote>
              </div>

              <details className="mt-4 rounded-xl border border-border/60 bg-card/30 px-5 py-4">
                <summary className="cursor-pointer text-[14px] font-medium">
                  What <code className="font-mono text-[13px]">grit migrate</code> prints the first time
                </summary>
                <div className="mt-3">
                  <CodeBlock
                    output="grit migrate"
                    language="text"
                    code={`  Running database migrations...
Database pool: up to 25 connections from this process (DB_MAX_OPEN_CONNS)
Database connected successfully
================================================================
DATABASE MIGRATION: 44 model(s) registered
================================================================
  + created *models.User
  + created *models.Session
  + created *models.PasswordResetToken
  ... 41 more
----------------------------------------------------------------
Migration done: 44 table(s) created, 0 altered (+0 column(s)), 0 unchanged.
================================================================

  + table blogs
  + table posts
  + table sessions
  + table users
  ... 40 more

This is the baseline run, the one that built the schema. It cannot be
rolled back: use grit migrate --fresh to start over.`}
                    className="mb-0"
                  />
                  <p className="text-[13px] text-muted-foreground mt-3 leading-relaxed">
                    Run it again and the summary reads{' '}
                    <code className="font-mono text-[12.5px]">0 table(s) created, 0 altered, 44 unchanged</code>.
                    That is the normal result rather than a failure: it means your models and
                    your schema already agree.
                  </p>
                </div>
              </details>

              {/* The seeder prints these and /docs/start lists them, and the
                  Quick Start left a reader at a login screen with no password
                  (grit#238). */}
              <div className="mt-4 rounded-xl border border-accent/25 bg-accent/[0.04] p-5">
                <p className="text-[15px] font-semibold mb-2">Your login</p>
                <p className="text-[15px] leading-relaxed">
                  <code className="font-mono text-[14px] text-accent">admin@example.com</code>
                  {' / '}
                  <code className="font-mono text-[14px] text-accent">admin123</code>
                </p>
                <p className="text-[13px] text-muted-foreground mt-3 leading-relaxed">
                  That is the dev default, and the seeder says so as it creates the account:
                  change it before production, or set{' '}
                  <code className="font-mono text-[12.5px]">SEED_ADMIN_PASSWORD</code> before
                  seeding and it uses that instead. Four demo users (jane, robert, emily and
                  michael, all at <code className="font-mono text-[12.5px]">example.com</code>)
                  share the same password.
                </p>
                <p className="text-[13px] text-muted-foreground mt-3 leading-relaxed">
                  <code className="font-mono text-[12.5px]">grit seed</code> also prints a
                  publishable and a secret API key, and writes them into{' '}
                  <code className="font-mono text-[12.5px]">apps/web/.env.local</code> and{' '}
                  <code className="font-mono text-[12.5px]">apps/admin/.env.local</code> for
                  you. The secret one is shown that once and never again.
                </p>
              </div>
            </Step>

            {/* Step 5 */}
            <Step title="Run Your App">
              <div className="prose-grit mb-4">
                <p>
                  One command starts the whole stack &mdash; the Go API, the web app,
                  and the admin panel &mdash; from the project root. No juggling
                  terminals, no <code>cd</code>-ing into folders. Press <code>Ctrl+C</code>{" "}
                  to stop everything.
                </p>
              </div>
              <CodeBlock terminal code="grit start" className="mb-0 glow-purple-sm" />

              <details className="mt-4 rounded-xl border border-border/60 bg-card/30 px-5 py-4">
                <summary className="cursor-pointer text-[14px] font-medium">
                  What it prints
                </summary>
                <div className="mt-3">
                  <CodeBlock
                    output="grit start"
                    language="text"
                    code={`  Starting API + client apps in parallel...
  API hot-reload via air (bundled, first run downloads it, then cached).
  Press Ctrl+C to stop everything.

  Recording output to .grit/logs/dev.log
[web] apps/admin dev$ rm -rf .next && next dev --webpack --port 3001
[web] apps/web dev$ rm -rf .next && next dev --webpack --port 3000
[web] apps/web dev: ▲ Next.js 16.4.6 (webpack)
[web] apps/web dev: - Local:         http://localhost:3000
[web] apps/web dev: ✓ Ready in 564ms
[web] apps/admin dev: - Local:         http://localhost:3001
[web] apps/admin dev: ✓ Ready in 566ms
[api] watching .
[api] building...
[api] running...
[api] [GIN-debug] GET  /api/v1/posts  --> handlers.(*PostHandler).List-fm
  ... one line per route
[api] Server starting on port 8080
[api] GORM Studio available at http://localhost:8080/studio
[api] API Documentation at http://localhost:8080/docs
[api] Pulse dashboard at http://localhost:8080/pulse/ui/
[api] Sentinel dashboard at http://localhost:8080/sentinel/ui`}
                    className="mb-0"
                  />
                  <p className="text-[13px] text-muted-foreground mt-3 leading-relaxed">
                    Two prefixes, interleaved: <code>[api]</code> is the Go server under air,
                    and <code>[web]</code> carries both Next apps. The frontends are ready
                    when each prints its own <strong>Ready in</strong>; the API is ready at{' '}
                    <strong>Server starting on port 8080</strong>, which arrives last because
                    it compiles first. The hundreds of <code>[GIN-debug]</code> lines between
                    them are one per route and are not a problem. All of it is also written
                    to <code>.grit/logs/dev.log</code>, so you can read what scrolled past.
                    If a port is already taken, that app fails and{' '}
                    <code>grit start</code> stops the rest rather than leaving half a stack
                    running.
                  </p>
                </div>
              </details>

              <DocsShot
                src="/images/quickstart/admin-login.png"
                alt="The admin panel sign-in screen, with email and password fields"
                label="localhost:3001/login"
                caption="The admin at :3001 opens here. Sign in with the account grit seed created: admin@example.com / admin123."
              />

              <DocsShot
                src="/images/quickstart/admin-dashboard.png"
                alt="The admin dashboard, showing stat cards, recent activity and a severity chart"
                label="localhost:3001"
                caption="What you land on after signing in. The numbers come from the rows grit seed wrote, so a dashboard of zeroes means the seed step was skipped rather than that something failed."
              />

              <div className="prose-grit mt-4">
                <p>
                  Need just one app? Use a subcommand: <code>grit start server</code>{" "}
                  (Go API only), <code>grit start web</code>, <code>grit start admin</code>,{" "}
                  <code>grit start expo</code>, or <code>grit start desktop</code>.
                </p>
                <p>Once started, you can access:</p>
              </div>
              <div className="grid gap-3 sm:grid-cols-2 mt-4">
                {[
                  {
                    name: "Go API",
                    url: "http://localhost:8080",
                    desc: "Backend + GORM Studio at /studio",
                  },
                  {
                    name: "Web App",
                    url: "http://localhost:3000",
                    desc: "Next.js frontend with auth pages",
                  },
                  {
                    name: "Admin Panel",
                    url: "http://localhost:3001",
                    desc: "Resource-based admin dashboard",
                  },
                  {
                    name: "API Docs",
                    url: "http://localhost:8080/docs",
                    desc: "Interactive Scalar API reference",
                  },
                  {
                    name: "Mailhog",
                    url: "http://localhost:8025",
                    desc: "Email testing inbox",
                  },
                ].map((item) => (
                  <div
                    key={item.name}
                    className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                  >
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-[15px] font-semibold">{item.name}</span>
                    </div>
                    <code className="text-sm font-mono text-accent/60 block mb-1">
                      {item.url}
                    </code>
                    <span className="text-sm text-muted-foreground/50">
                      {item.desc}
                    </span>
                  </div>
                ))}
              </div>
              <div className="prose-grit mt-4">
                <p>
                  Try registering a user at{" "}
                  <code>http://localhost:3000/register</code>, then log in and
                  explore the dashboard. Open <code>http://localhost:3001</code>{" "}
                  to see the admin panel. Visit{" "}
                  <code>http://localhost:8080/studio</code> to browse your
                  database visually with GORM Studio.
                </p>
              </div>
            </Step>

            {/* Step 6 */}
            <Step title="Generate a Resource">
              <div className="prose-grit mb-4">
                <p>
                  Now for the magic. Generate a complete full-stack resource
                  with a single command. This creates the Go model, CRUD
                  handler, service layer, React Query hooks, Zod schemas,
                  TypeScript types, and an admin page -- all wired together.
                </p>
              </div>
              <CodeBlock terminal code={`grit generate resource Post --fields "title:string,content:text,published:bool"`} className="mb-0 glow-purple-sm" />
              <div className="prose-grit mt-4">
                <p>This generates the following files:</p>
                <ul>
                  <li>
                    <code>apps/api/internal/models/post.go</code> -- GORM model
                    with struct tags
                  </li>
                  <li>
                    <code>apps/api/internal/handlers/post.go</code> -- Full CRUD
                    handler with pagination
                  </li>
                  <li>
                    <code>apps/api/internal/services/post.go</code> -- Business
                    logic layer
                  </li>
                  <li>
                    <code>packages/shared/schemas/post.ts</code> -- Zod
                    validation schemas
                  </li>
                  <li>
                    <code>packages/shared/types/post.ts</code> -- TypeScript
                    types
                  </li>
                  <li>
                    <code>apps/api/internal/handlers/post_import.go</code> -- CSV/Excel
                    bulk-import handler
                  </li>
                  <li>
                    <code>apps/web/hooks/use-posts.ts</code> -- React Query
                    hooks
                  </li>
                  <li>
                    <code>apps/admin/app/(dashboard)/resources/posts/page.tsx</code> --
                    Admin page with data table
                  </li>
                </ul>
                <p>
                  It also registers the routes in <code>routes.go</code>, adds the model
                  to the migration registry, and injects the resource into the admin
                  sidebar. To finish, create the new table and refresh:
                </p>
                <CodeBlock terminal code="grit migrate" className="mb-0" />
                <p className="mt-4">
                  Then reload the admin panel &mdash; your <strong>Posts</strong> resource
                  is there with a working data table and create form. (Full breakdown of
                  every generated file:{" "}
                  <Link href="/docs/concepts/generated-files" className="text-accent hover:underline">Generated File Map</Link>.)
                </p>
              </div>

              <details className="mt-4 rounded-xl border border-border/60 bg-card/30 px-5 py-4">
                <summary className="cursor-pointer text-[14px] font-medium">
                  What it prints
                </summary>
                <div className="mt-3">
                  <CodeBlock
                    output="grit generate resource Post"
                    language="text"
                    code={`  Generating resource: Post

  ✓ apps/api/internal/models/post.go
  ✓ apps/api/internal/services/post.go
  ✓ apps/api/internal/handlers/post.go
  ✓ apps/api/internal/policies/post.go
  ✓ packages/shared/schemas/post.ts
  ✓ packages/shared/types/post.ts
  ✓ apps/web/hooks/use-posts.ts
  ✓ apps/admin/resources/posts/posts.ts
  ✓ apps/admin/app/(dashboard)/resources/posts/page.tsx
  ... 6 more

  Injecting into existing files...
  ✓ Injected model into AutoMigrate
  ✓ Routes in internal/routes/post_routes.go
  ✓ Injected permissions into the authz catalog
  ✓ Injected resource into registry list
  ... 12 more

  Resource Post generated successfully!

  Next steps:
    1. cd apps/api && go build ./...
    2. grit migrate   (creates the posts table)
    3. Restart the API server
    4. The admin panel will show Posts in the sidebar

  ✓ apps/api/internal/access/registry.go now covers 209 routes`}
                    className="mb-0"
                  />
                </div>
              </details>

              {/* Step 3 of the generator's own next steps reads as a chore and
                  usually is not one. Answer it here rather than leaving a
                  reader to guess (grit#238). */}
              <div className="mt-4 rounded-xl border border-border/60 bg-card/30 p-5">
                <p className="text-[15px] font-semibold mb-2">
                  Do I need to restart <code className="font-mono text-[14px]">grit start</code>?
                </p>
                <p className="text-[15px] leading-relaxed text-muted-foreground">
                  No. <code className="font-mono text-[14px]">grit start</code> runs the Go
                  API under <strong>air</strong>, so the new handler and model are rebuilt
                  and reloaded on their own, and both Next apps pick up the new admin page
                  the same way. The step that is not automatic is{' '}
                  <code className="font-mono text-[14px]">grit migrate</code>: the table does
                  not exist until you run it, and until then the new screen answers with a
                  database error. The generator&apos;s own &quot;Restart the API server&quot;
                  line is for running the server by hand with{' '}
                  <code className="font-mono text-[14px]">go run</code>.
                </p>
              </div>

              <DocsShot
                src="/images/quickstart/posts-screen.png"
                alt="The generated Posts screen in the admin, with a data table and a create button"
                label="localhost:3001/resources/posts"
                caption="The Posts screen one command later: a sidebar entry, a searchable table with column controls, and a create form built from the same field list you typed."
              />
            </Step>
            </Steps>

            {/* Upgrade */}
            <div className="mb-10 rounded-xl border border-accent/20 bg-accent/5 p-5">
              <h3 className="text-lg font-semibold mb-2 flex items-center gap-2">
                <span className="text-accent">Upgrading?</span>
                <span className="text-xs font-mono text-accent/60 bg-accent/10 rounded-full px-2 py-0.5">v{GRIT_VERSION}</span>
              </h3>
              <p className="text-[15px] text-muted-foreground mb-3">
                If you have an existing Grit project and want to update the framework
                components (admin panel, configs, web app) to the latest version, run:
              </p>
              <CodeBlock terminal code="grit upgrade" className="mb-0" />
              <p className="text-sm text-muted-foreground/50 mt-2">
                This preserves your resource definitions and API code while updating
                all framework-generated files. Use <code className="text-accent/60">grit upgrade --force</code> to
                overwrite without prompting.
              </p>
            </div>

            {/* What's Next */}
            <div className="mb-8">
              <h2 className="text-2xl font-semibold tracking-tight mb-6">
                What&apos;s Next?
              </h2>
              <div className="grid gap-3 sm:grid-cols-2">
                {[
                  {
                    title: "Show It on the Web App",
                    desc: "Products in the admin, read publicly on the home and detail pages",
                    href: "/docs/frontend/fetching-data",
                  },
                  {
                    title: "Sign-In and Private Pages",
                    desc: "Auth hooks, the edge gate, and protecting a page properly",
                    href: "/docs/frontend/auth",
                  },
                  {
                    title: "Add Authentication",
                    desc: "Register, login, JWT, OAuth and 2FA, built in",
                    href: "/docs/backend/authentication",
                  },
                  {
                    title: "Seed Sample Data",
                    desc: "grit seed, faker, and per-resource seeders",
                    href: "/docs/backend/seeders",
                  },
                  {
                    title: "Build a Real App",
                    desc: "Step-by-step: your first Grit app end to end",
                    href: "/docs/tutorials/contact-app",
                  },
                  {
                    title: "CLI Commands",
                    desc: "Every command the Grit CLI offers",
                    href: "/docs/concepts/cli",
                  },
                  {
                    title: "Code Generation",
                    desc: "Deep dive into resource generation",
                    href: "/docs/concepts/code-generation",
                  },
                  {
                    title: "Project Structure",
                    desc: "The monorepo layout and where things live",
                    href: "/docs/getting-started/project-structure",
                  },
                ].map((item) => (
                  <Link key={item.href} href={item.href}>
                    <div className="group rounded-lg border border-border/40 bg-card/50 p-4 hover:border-accent/20 hover:bg-card/80 transition-all duration-200">
                      <h3 className="text-[15px] font-semibold mb-1 group-hover:text-accent transition-colors flex items-center gap-1.5">
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

            {/* Bug Reports */}
            <div className="mb-8 rounded-xl border border-accent/20 bg-accent/5 p-5">
              <p className="text-[15px] text-muted-foreground leading-relaxed">
                Found a bug or something doesn&apos;t work? Please open an issue at{' '}
                <a href="https://github.com/MUKE-coder/grit/issues" target="_blank" rel="noopener noreferrer" className="text-accent hover:underline">
                  https://github.com/MUKE-coder/grit/issues
                </a>
                . Your feedback helps us improve Grit for everyone.
              </p>
            </div>

            {/* Example Projects */}
            <div className="mt-10 rounded-lg border border-accent/20 bg-accent/5 p-5">
              <h3 className="text-sm font-semibold text-foreground mb-2">Example Projects</h3>
              <p className="text-sm text-muted-foreground mb-3">
                See the same Job Portal app built with every architecture: full source code, setup guide, and deployment config.
              </p>
              <a href="https://github.com/MUKE-coder/grit/tree/main/examples" target="_blank" rel="noreferrer" className="inline-flex items-center gap-1.5 text-sm font-medium text-accent hover:text-accent/80 transition-colors">
                Browse all 6 examples on GitHub &rarr;
              </a>
            </div>

            {/* Community, placed at the end of the first thing anyone reads,
                where "it didn't work on my machine" is most likely. */}
            <CommunityCTA className="mt-10" />

            {/* Nav */}
            <div className="flex flex-wrap gap-3 mt-12 pt-6 border-t border-border/30">
              <Button
                variant="outline"
                asChild
                className="border-border/60 bg-transparent hover:bg-surface-2"
              >
                <Link href="/docs/getting-started/philosophy">Philosophy</Link>
              </Button>
              <Button asChild className="glow-purple-sm ml-auto">
                <Link href="/docs/getting-started/installation">
                  Installation
                  <ArrowRight className="ml-1.5 h-4 w-4" />
                </Link>
              </Button>
            </div>
    </DocsShell>
  );
}
