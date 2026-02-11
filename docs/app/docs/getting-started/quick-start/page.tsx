import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'

export default function QuickStartPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Getting Started</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Quick Start
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Get a full-stack Grit project running in 5 minutes. This guide assumes you have
                Go, Node.js, pnpm, and Docker already installed. If not, see the{' '}
                <Link href="/docs/getting-started/installation" className="text-primary hover:underline">
                  Installation guide
                </Link>{' '}
                first.
              </p>
            </div>

            {/* Prerequisites */}
            <div className="prose-grit mb-10">
              <h2>Prerequisites</h2>
              <p>
                Make sure the following tools are installed on your machine before proceeding:
              </p>
            </div>

            <div className="grid gap-3 sm:grid-cols-2 mb-10">
              {[
                { name: 'Go', version: '1.21+', check: 'go version' },
                { name: 'Node.js', version: '18+', check: 'node --version' },
                { name: 'pnpm', version: '8+', check: 'pnpm --version' },
                { name: 'Docker', version: 'Latest', check: 'docker --version' },
              ].map((tool) => (
                <div key={tool.name} className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-sm font-semibold">{tool.name}</span>
                    <span className="text-xs font-mono text-primary/60">{tool.version}</span>
                  </div>
                  <code className="text-xs font-mono text-muted-foreground/50">{tool.check}</code>
                </div>
              ))}
            </div>

            {/* Step 1 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  1
                </div>
                <h2 className="text-xl font-semibold tracking-tight">Install the Grit CLI</h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Install the Grit CLI globally using <code>go install</code>. This gives you the{' '}
                  <code>grit</code> command available anywhere on your system.
                </p>
              </div>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <span className="text-primary/50 select-none">$ </span>
                  <span className="text-foreground/80">go install github.com/MUKE-coder/grit@latest</span>
                </div>
              </div>
              <div className="prose-grit mt-4">
                <p>
                  Verify the installation by running <code>grit --help</code>. You should see
                  the Grit ASCII art logo and a list of available commands.
                </p>
              </div>
            </div>

            {/* Step 2 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  2
                </div>
                <h2 className="text-xl font-semibold tracking-tight">Create a New Project</h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Scaffold a complete full-stack project with one command. This creates the entire
                  monorepo: Go API, Next.js web app, admin panel, shared types, Docker setup, and
                  all the batteries.
                </p>
              </div>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm space-y-2">
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">grit new myapp</span>
                  </div>
                  <div className="text-muted-foreground/40 text-xs space-y-1 pl-4 mt-2">
                    <div className="flex items-center gap-2">
                      <span className="text-primary/60">+</span>
                      <span>Creating directory structure...</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-primary/60">+</span>
                      <span>Scaffolding Go API...</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-primary/60">+</span>
                      <span>Adding batteries (cache, storage, mail, jobs, cron, AI)...</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-primary/60">+</span>
                      <span>Setting up Next.js web app...</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-primary/60">+</span>
                      <span>Creating admin panel...</span>
                    </div>
                    <div className="flex items-center gap-2 text-primary/80">
                      <span>&#10003;</span>
                      <span>Project &quot;myapp&quot; created successfully!</span>
                    </div>
                  </div>
                </div>
              </div>
              <div className="prose-grit mt-4">
                <p>
                  The project name must be lowercase, alphanumeric, and hyphens only (e.g.,{' '}
                  <code>my-saas-app</code>). It must start with a letter and cannot end with
                  a hyphen.
                </p>
              </div>
            </div>

            {/* Step 3 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  3
                </div>
                <h2 className="text-xl font-semibold tracking-tight">Start Infrastructure Services</h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Navigate into the project and start the Docker services. This launches PostgreSQL,
                  Redis, MinIO (S3-compatible storage), and Mailhog (email testing).
                </p>
              </div>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm space-y-2">
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">cd myapp</span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">docker compose up -d</span>
                  </div>
                </div>
              </div>
              <div className="prose-grit mt-4">
                <p>
                  This starts the following services in the background:
                </p>
                <ul>
                  <li><strong>PostgreSQL 16</strong> on port 5432 -- your primary database</li>
                  <li><strong>Redis 7</strong> on port 6379 -- caching and job queues</li>
                  <li><strong>MinIO</strong> on port 9000 (console: 9001) -- local S3-compatible file storage</li>
                  <li><strong>Mailhog</strong> on port 8025 -- catch-all email testing UI</li>
                </ul>
                <blockquote>
                  Do not have Docker? See the{' '}
                  <Link href="/docs/getting-started/installation" className="text-primary hover:underline">
                    Installation guide
                  </Link>{' '}
                  for a cloud-only setup using Neon (Postgres) and Upstash (Redis) instead.
                </blockquote>
              </div>
            </div>

            {/* Step 4 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  4
                </div>
                <h2 className="text-xl font-semibold tracking-tight">Start the Go API</h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Navigate into the Go API directory, install dependencies with <code>go mod tidy</code>,
                  then start the server. This runs the Go backend on port 8080 with auto-migration
                  enabled.
                </p>
              </div>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm space-y-2">
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">cd apps/api</span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">go mod tidy</span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">go run cmd/server/main.go</span>
                  </div>
                </div>
              </div>
              <div className="prose-grit mt-4">
                <p>
                  You should see the Gin router start up and log all registered routes. The API
                  is now running at <code>http://localhost:8080</code> and GORM Studio is
                  available at <code>http://localhost:8080/studio</code>.
                </p>
                <blockquote>
                  The first run of <code>go mod tidy</code> may take a minute as Go downloads
                  all dependencies. Subsequent runs will be instant.
                </blockquote>
              </div>
            </div>

            {/* Step 5 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  5
                </div>
                <h2 className="text-xl font-semibold tracking-tight">Start the Frontend</h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Open a <strong>new terminal</strong> (keep the API running), navigate back to the
                  project root, install Node.js dependencies, then start the Next.js web app.
                </p>
              </div>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal (new tab)</span>
                </div>
                <div className="p-5 font-mono text-sm space-y-2">
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">cd myapp</span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">pnpm install</span>
                  </div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">cd apps/web &amp;&amp; pnpm dev</span>
                  </div>
                </div>
              </div>
              <div className="prose-grit mt-4">
                <p>
                  To also run the admin panel, open <strong>another terminal</strong> and run:
                </p>
              </div>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm mt-4">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal (another tab)</span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <span className="text-primary/50 select-none">$ </span>
                  <span className="text-foreground/80">cd myapp/apps/admin &amp;&amp; pnpm dev</span>
                </div>
              </div>
              <div className="prose-grit mt-4">
                <p>
                  Alternatively, you can run everything at once with Turborepo from the project root:
                </p>
              </div>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm mt-4">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm space-y-2">
                  <div className="text-muted-foreground/40 text-xs">{`# From the project root (myapp/)`}</div>
                  <div>
                    <span className="text-primary/50 select-none">$ </span>
                    <span className="text-foreground/80">turbo dev</span>
                  </div>
                </div>
              </div>
              <div className="prose-grit mt-4">
                <p>
                  Once started, you can access:
                </p>
              </div>
              <div className="grid gap-3 sm:grid-cols-2 mt-4">
                {[
                  { name: 'Go API', url: 'http://localhost:8080', desc: 'Backend + GORM Studio at /studio' },
                  { name: 'Web App', url: 'http://localhost:3000', desc: 'Next.js frontend with auth pages' },
                  { name: 'Admin Panel', url: 'http://localhost:3001', desc: 'Resource-based admin dashboard' },
                  { name: 'Mailhog', url: 'http://localhost:8025', desc: 'Email testing inbox' },
                ].map((item) => (
                  <div key={item.name} className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-sm font-semibold">{item.name}</span>
                    </div>
                    <code className="text-xs font-mono text-primary/60 block mb-1">{item.url}</code>
                    <span className="text-xs text-muted-foreground/50">{item.desc}</span>
                  </div>
                ))}
              </div>
              <div className="prose-grit mt-4">
                <p>
                  Try registering a user at <code>http://localhost:3000/register</code>, then log
                  in and explore the dashboard. Open <code>http://localhost:3001</code> to see the
                  admin panel. Visit <code>http://localhost:8080/studio</code> to browse your
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
                <h2 className="text-xl font-semibold tracking-tight">Generate a Resource</h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Now for the magic. Generate a complete full-stack resource with a single command.
                  This creates the Go model, CRUD handler, service layer, React Query hooks,
                  Zod schemas, TypeScript types, and an admin page -- all wired together.
                </p>
              </div>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <div className="flex items-center gap-1.5">
                    <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <span className="text-primary/50 select-none">$ </span>
                  <span className="text-foreground/80">grit generate resource Post --fields &quot;title:string,content:text,published:bool&quot;</span>
                </div>
              </div>
              <div className="prose-grit mt-4">
                <p>
                  This generates the following files:
                </p>
                <ul>
                  <li><code>apps/api/internal/models/post.go</code> -- GORM model with struct tags</li>
                  <li><code>apps/api/internal/handlers/post.go</code> -- Full CRUD handler with pagination</li>
                  <li><code>apps/api/internal/services/post.go</code> -- Business logic layer</li>
                  <li><code>packages/shared/schemas/post.ts</code> -- Zod validation schemas</li>
                  <li><code>packages/shared/types/post.ts</code> -- TypeScript types</li>
                  <li><code>apps/admin/hooks/use-posts.ts</code> -- React Query hooks</li>
                  <li><code>apps/admin/app/resources/posts/page.tsx</code> -- Admin page with data table</li>
                </ul>
                <p>
                  It also automatically registers the routes in <code>routes.go</code>, adds the
                  model to auto-migrations, and injects the resource into the admin sidebar.
                  Restart <code>turbo dev</code> and visit the admin panel to see your new
                  Posts resource with a fully functional data table and create form.
                </p>
              </div>
            </div>

            {/* What's Next */}
            <div className="mb-8">
              <h2 className="text-2xl font-semibold tracking-tight mb-6">
                What&apos;s Next?
              </h2>
              <div className="grid gap-3 sm:grid-cols-2">
                {[
                  { title: 'Project Structure', desc: 'Understand the monorepo layout and where things live', href: '/docs/getting-started/project-structure' },
                  { title: 'Configuration', desc: 'All .env variables explained', href: '/docs/getting-started/configuration' },
                  { title: 'CLI Commands', desc: 'Every command the Grit CLI offers', href: '/docs/concepts/cli' },
                  { title: 'Code Generation', desc: 'Deep dive into resource generation', href: '/docs/concepts/code-generation' },
                  { title: 'Troubleshooting', desc: 'Common errors and how to fix them', href: '/docs/getting-started/troubleshooting' },
                ].map((item) => (
                  <Link key={item.href} href={item.href}>
                    <div className="group rounded-lg border border-border/40 bg-card/50 p-4 hover:border-primary/20 hover:bg-card/80 transition-all duration-200">
                      <h3 className="text-sm font-semibold mb-1 group-hover:text-primary transition-colors flex items-center gap-1.5">
                        {item.title}
                        <ArrowRight className="h-3 w-3 opacity-0 -translate-x-1 group-hover:opacity-100 group-hover:translate-x-0 transition-all" />
                      </h3>
                      <p className="text-xs text-muted-foreground/60">{item.desc}</p>
                    </div>
                  </Link>
                ))}
              </div>
            </div>

            {/* Nav */}
            <div className="flex flex-wrap gap-3 mt-12 pt-6 border-t border-border/30">
              <Button variant="outline" asChild className="border-border/60 bg-transparent hover:bg-accent/50">
                <Link href="/docs/getting-started/philosophy">
                  Philosophy
                </Link>
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
  )
}
