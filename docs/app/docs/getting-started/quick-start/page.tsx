import Link from "next/link";
import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from "@/components/code-block";
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/getting-started/quick-start')

export default function QuickStartPage() {
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
                Getting Started
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Quick Start
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Get a full-stack Grit project running in 5 minutes. This guide
                assumes you have Go, Node.js, pnpm, and Docker already
                installed. If not, see the{" "}
                <Link
                  href="/docs/getting-started/installation"
                  className="text-primary hover:underline"
                >
                  Installation guide
                </Link>{" "}
                first.
              </p>
            </div>

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
                { name: "Go", version: "1.21+", check: "go version" },
                { name: "Node.js", version: "18+", check: "node --version" },
                { name: "pnpm", version: "8+", check: "pnpm --version" },
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

            {/* Step 1 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  1
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Install the Grit CLI
                </h2>
              </div>
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
                  with a single command. This removes the old binary and installs the
                  newest release:
                </p>
                <CodeBlock terminal code="grit update" className="mb-0" />
                <p className="text-[15px] text-muted-foreground">
                  Run <code className="text-[15px] font-mono bg-accent/80 px-1.5 py-0.5 rounded text-primary">grit version</code> afterwards
                  to confirm you&apos;re on the latest release.
                </p>
              </div>
            </div>

            {/* Step 2 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  2
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Create a New Project
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Scaffold a complete full-stack project with one command. This
                  creates the entire monorepo: Go API, Next.js web app, admin
                  panel, shared types, Docker setup, and all the batteries.
                </p>
              </div>
              <CodeBlock terminal code="grit new myapp" className="mb-0 glow-purple-sm" />
              <div className="prose-grit mt-4">
                <p>
                  The <strong>interactive CLI</strong> walks you through selecting an architecture mode
                  and frontend option. Grit supports <strong>5 architecture modes</strong>:{' '}
                  <code>single</code> (API only), <code>double</code> (API + web),{' '}
                  <code>triple</code> (API + web + admin), <code>api</code> (headless API),
                  and <code>mobile</code> (API + Expo). It also supports <strong>2 frontend options</strong>:{' '}
                  <code>Next.js</code> and <code>TanStack Router (Vite)</code>.
                </p>
                <p>
                  The project name must be lowercase, alphanumeric, and hyphens
                  only (e.g., <code>my-saas-app</code>). It must start with a
                  letter and cannot end with a hyphen.
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
                  Available flags: <code className="text-[15px] font-mono bg-accent/80 px-1.5 py-0.5 rounded text-primary">--single</code>,{' '}
                  <code className="text-[15px] font-mono bg-accent/80 px-1.5 py-0.5 rounded text-primary">--double</code>,{' '}
                  <code className="text-[15px] font-mono bg-accent/80 px-1.5 py-0.5 rounded text-primary">--triple</code>,{' '}
                  <code className="text-[15px] font-mono bg-accent/80 px-1.5 py-0.5 rounded text-primary">--api</code>,{' '}
                  <code className="text-[15px] font-mono bg-accent/80 px-1.5 py-0.5 rounded text-primary">--mobile</code> for architecture,{' '}
                  <code className="text-[15px] font-mono bg-accent/80 px-1.5 py-0.5 rounded text-primary">--vite</code>,{' '}
                  <code className="text-[15px] font-mono bg-accent/80 px-1.5 py-0.5 rounded text-primary">--next</code> for frontend,
                  and <code className="text-[15px] font-mono bg-accent/80 px-1.5 py-0.5 rounded text-primary">--desktop</code> to add a
                  Wails desktop client (<a href="/docs/concepts/architecture-modes/multi-client" className="text-primary hover:underline">multi-client pattern</a>).
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
                  Use <code className="text-[15px] font-mono bg-accent/80 px-1.5 py-0.5 rounded text-primary">--force</code> to
                  scaffold into a non-empty directory.
                </p>
              </div>
            </div>

            {/* Step 3 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  3
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Start Infrastructure Services
                </h2>
              </div>
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
                    className="text-primary hover:underline"
                  >
                    Installation guide
                  </Link>{" "}
                  for a cloud-only setup using Neon (Postgres) and Upstash
                  (Redis) instead.
                </blockquote>
              </div>
            </div>

            {/* Step 4 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  4
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Start the Go API
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Navigate into the Go API directory, install dependencies with{" "}
                  <code>go mod tidy</code>, then start the server. This runs the
                  Go backend on port 8080 with auto-migration enabled.
                </p>
              </div>
              <CodeBlock terminal code={`cd apps/api
go mod tidy
go run cmd/server/main.go`} className="mb-0 glow-purple-sm" />
              <div className="prose-grit mt-4">
                <p>
                  You should see the Gin router start up and log all registered
                  routes. The API is now running at{" "}
                  <code>http://localhost:8080</code>. Interactive API docs (Scalar) are at{" "}
                  <code>http://localhost:8080/docs</code> and GORM Studio is
                  available at <code>http://localhost:8080/studio</code>.
                </p>
                <blockquote>
                  The first run of <code>go mod tidy</code> may take a minute as
                  Go downloads all dependencies. Subsequent runs will be
                  instant.
                </blockquote>
              </div>
            </div>

            {/* Step 5 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  5
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Start the Frontend
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Open a <strong>new terminal</strong> (keep the API running),
                  navigate back to the project root, install Node.js
                  dependencies, then start the Next.js web app.
                </p>
              </div>
              <CodeBlock terminal filename="terminal (new tab)" code={`cd myapp
pnpm install
cd apps/web && pnpm dev`} className="mb-0 glow-purple-sm" />
              <div className="prose-grit mt-4">
                <p>
                  To also run the admin panel, open{" "}
                  <strong>another terminal</strong> and run:
                </p>
              </div>
              <CodeBlock terminal filename="terminal (another tab)" code="cd myapp/apps/admin && pnpm dev" className="mt-4 mb-0 glow-purple-sm" />
              <div className="prose-grit mt-4">
                <p>
                  Alternatively, you can run everything at once from the project
                  root:
                </p>
              </div>
              <CodeBlock terminal code={`# From the project root (myapp/)
pnpm dev`} className="mt-4 mb-0 glow-purple-sm" />
              <div className="prose-grit mt-4">
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
                    <code className="text-sm font-mono text-primary/60 block mb-1">
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
            </div>

            {/* Step 6 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  6
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Generate a Resource
                </h2>
              </div>
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
                    <code>apps/admin/hooks/use-posts.ts</code> -- React Query
                    hooks
                  </li>
                  <li>
                    <code>apps/admin/app/resources/posts/page.tsx</code> --
                    Admin page with data table
                  </li>
                </ul>
                <p>
                  It also automatically registers the routes in{" "}
                  <code>routes.go</code>, adds the model to auto-migrations, and
                  injects the resource into the admin sidebar. Restart{" "}
                  <code>pnpm dev</code> and visit the admin panel to see your
                  new Posts resource with a fully functional data table and
                  create form.
                </p>
              </div>
            </div>

            {/* Upgrade */}
            <div className="mb-10 rounded-xl border border-primary/20 bg-primary/5 p-5">
              <h3 className="text-lg font-semibold mb-2 flex items-center gap-2">
                <span className="text-primary">Upgrading?</span>
                <span className="text-xs font-mono text-primary/60 bg-primary/10 rounded-full px-2 py-0.5">v3.55.0</span>
              </h3>
              <p className="text-[15px] text-muted-foreground mb-3">
                If you have an existing Grit project and want to update the framework
                components (admin panel, configs, web app) to the latest version, run:
              </p>
              <CodeBlock terminal code="grit upgrade" className="mb-0" />
              <p className="text-sm text-muted-foreground/50 mt-2">
                This preserves your resource definitions and API code while updating
                all framework-generated files. Use <code className="text-primary/60">grit upgrade --force</code> to
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
                    title: "Project Structure",
                    desc: "Understand the monorepo layout and where things live",
                    href: "/docs/getting-started/project-structure",
                  },
                  {
                    title: "Configuration",
                    desc: "All .env variables explained",
                    href: "/docs/getting-started/configuration",
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
                    title: "Troubleshooting",
                    desc: "Common errors and how to fix them",
                    href: "/docs/getting-started/troubleshooting",
                  },
                ].map((item) => (
                  <Link key={item.href} href={item.href}>
                    <div className="group rounded-lg border border-border/40 bg-card/50 p-4 hover:border-primary/20 hover:bg-card/80 transition-all duration-200">
                      <h3 className="text-[15px] font-semibold mb-1 group-hover:text-primary transition-colors flex items-center gap-1.5">
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
            <div className="mb-8 rounded-xl border border-primary/20 bg-primary/5 p-5">
              <p className="text-[15px] text-muted-foreground leading-relaxed">
                Found a bug or something doesn&apos;t work? Please open an issue at{' '}
                <a href="https://github.com/MUKE-coder/grit/issues" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">
                  https://github.com/MUKE-coder/grit/issues
                </a>{' '}
                — your feedback helps us improve Grit for everyone.
              </p>
            </div>

            {/* Example Projects */}
            <div className="mt-10 rounded-lg border border-primary/20 bg-primary/5 p-5">
              <h3 className="text-sm font-semibold text-foreground mb-2">Example Projects</h3>
              <p className="text-sm text-muted-foreground mb-3">
                See the same Job Portal app built with every architecture — full source code, setup guide, and deployment config.
              </p>
              <a href="https://github.com/MUKE-coder/grit/tree/main/examples" target="_blank" rel="noreferrer" className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:text-primary/80 transition-colors">
                Browse all 6 examples on GitHub &rarr;
              </a>
            </div>

            {/* Nav */}
            <div className="flex flex-wrap gap-3 mt-12 pt-6 border-t border-border/30">
              <Button
                variant="outline"
                asChild
                className="border-border/60 bg-transparent hover:bg-accent/50"
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
          </div>
        </div>
      </main>
    </div>
  );
}
