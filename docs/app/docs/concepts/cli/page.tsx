import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/cli')

export default function CLICommandsPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Core Concepts</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                CLI Commands
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                The Grit CLI is a single binary that scaffolds projects, generates full-stack
                resources, and syncs types between Go and TypeScript. Running{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new</code> is
                interactive by default -- it walks you through architecture mode and frontend
                selection. Install it once and use it across all your Grit projects.
              </p>
            </div>

            <div className="prose-grit">
              {/* Lifecycle */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The CLI across a project&apos;s life
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Four commands carry a Grit project from empty folder to production. Everything
                  else &mdash; <code>migrate</code>, <code>seed</code>, <code>sync</code>,{' '}
                  <code>studio</code> &mdash; slots in between these stages.
                </p>
                <LaneFlow
                  id="cli-life"
                  lanes={['1 · Scaffold', '2 · Generate', '3 · Run', '4 · Ship']}
                  nodes={[
                    { id: 'new', lane: 0, row: 0, title: 'grit new', sub: 'scaffold project', tone: 'primary' },
                    { id: 'gen', lane: 1, row: 0, title: 'grit generate', sub: 'full-stack resource', tone: 'cyan' },
                    { id: 'run', lane: 2, row: 0, title: 'grit start', sub: 'dev servers', tone: 'blue' },
                    { id: 'ship', lane: 3, row: 0, title: 'grit deploy', sub: 'to production', tone: 'green' },
                  ]}
                  edges={[
                    { from: 'new', to: 'gen', label: 'add features', tone: 'cyan' },
                    { from: 'gen', to: 'run', label: 'develop', tone: 'blue' },
                    { from: 'run', to: 'ship', label: 'release', tone: 'green' },
                  ]}
                  legend={[
                    { tone: 'primary', label: 'Scaffold' },
                    { tone: 'cyan', label: 'Generate' },
                    { tone: 'blue', label: 'Run' },
                    { tone: 'green', label: 'Ship' },
                  ]}
                  caption="migrate · seed · sync · studio slot in between these four stages"
                />
              </div>

              {/* Installation */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Installing the CLI
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Install the Grit CLI globally using Go:
                </p>
                <CodeBlock language="bash" code={`$ grit generate resource Invoice -i

  Defining fields for Invoice
  Enter fields as name:type (e.g., title:string)
  Valid types: string, text, int, uint, float, bool, datetime, date, slug, belongs_to, many_to_many
  Press Enter with no input when done.

  > number:string
  \u2713 Added number (string)
  > amount:float
  \u2713 Added amount (float)
  > status:string
  \u2713 Added status (string)
  > due_date:date
  \u2713 Added due_date (date)
  > paid:bool
  \u2713 Added paid (bool)
  >`} />

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  More Examples
                </h3>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-4">
                    <div>
                      <span className="text-muted-foreground/40"># Blog post with title, content, and published flag</span>
                      <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit g resource Post --fields &quot;title:string,content:text,published:bool&quot;</span></div>
                    </div>
                    <div>
                      <span className="text-muted-foreground/40"># Product with name, price, and stock</span>
                      <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit g resource Product --fields &quot;name:string,description:text,price:float,stock:uint&quot;</span></div>
                    </div>
                    <div>
                      <span className="text-muted-foreground/40"># Event with title, date, and description</span>
                      <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit g resource Event --fields &quot;title:string,description:text,start_date:datetime,end_date:datetime&quot;</span></div>
                    </div>
                    <div>
                      <span className="text-muted-foreground/40"># Category with just a name</span>
                      <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit g resource Category --fields &quot;name:string,description:text&quot;</span></div>
                    </div>
                    <div>
                      <span className="text-muted-foreground/40"># Article with an auto-generated slug</span>
                      <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit g resource Article --fields &quot;title:string,slug:slug,body:text,published:bool&quot;</span></div>
                    </div>
                  </div>
                </div>
              </div>

              {/* grit remove resource */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  grit remove resource
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Remove a previously generated resource. This deletes the Go model, service, handler,
                  Zod schemas, TypeScript types, React hooks, resource definition, and admin page. It also
                  cleans up all injection markers that were added when the resource was generated.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit remove resource Post</span></div>
                  </div>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Syntax
                </h3>
                <CodeBlock filename="usage" code={`grit remove resource <Name>

# Shorthand alias
grit rm resource <Name>`} />

                <p className="text-sm text-muted-foreground/60 mt-3">
                  The resource name should be the singular PascalCase name (e.g.,{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Post</code>,{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Product</code>,{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">BlogCategory</code>)
                  — the same name you used with{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>.
                </p>
              </div>

              {/* grit add role */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  grit add role
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Add a new role to your project. This command updates all relevant files across the stack
                  in one step — Go model constants, TypeScript types, Zod schemas, shared constants, and
                  admin panel resource definitions (badge, filter, and form options).
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit add role MODERATOR</span></div>
                  </div>
                </div>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  This single command updates <strong className="text-foreground/90">7 locations</strong> across your project:
                </p>

                <ul className="space-y-2 mb-6">
                  {[
                    'Go model constants (RoleModerator = "MODERATOR")',
                    'Zod schema enum validation',
                    'TypeScript union type',
                    'ROLES constants object',
                    'Admin badge configuration',
                    'Admin table filter options',
                    'Admin form select options',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="mt-1.5 h-1.5 w-1.5 rounded-full bg-primary/50 shrink-0" />
                      <span>{item}</span>
                    </li>
                  ))}
                </ul>

                <p className="text-sm text-muted-foreground/60">
                  The role name is automatically uppercased. Multi-word roles use underscores:
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded ml-1">grit add role CONTENT_MANAGER</code>
                </p>
              </div>

              {/* grit start */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  grit start
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Start development servers for your Grit project. Use subcommands to launch
                  the frontend client apps or the Go API server individually.
                </p>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  grit start client
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Runs <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">pnpm dev</code> from
                  the project root, which starts all frontend apps (web, admin, expo, docs) via Turborepo.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm mb-8">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit start client</span></div>
                  </div>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  grit start server
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Runs <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">go run cmd/server/main.go</code> from
                  the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/api</code> directory
                  to start the Go API server.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm mb-4">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit start server</span></div>
                  </div>
                </div>

                <p className="text-sm text-muted-foreground/60">
                  Both commands auto-detect the project root by looking for <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker-compose.yml</code> or <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">turbo.json</code>,
                  so you can run them from any subdirectory within your project.
                </p>
              </div>

              {/* grit sync */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  grit sync
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Parse all Go model files and regenerate the corresponding TypeScript types and
                  Zod schemas in the shared package. Use this command whenever you manually modify
                  a Go model and want the frontend types to stay in sync.
                </p>

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
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit sync</span></div>
                  </div>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  How It Works
                </h3>
                <ol className="space-y-2.5 mb-4 list-decimal list-inside">
                  {[
                    'Finds the project root by walking up directories looking for docker-compose.yml or turbo.json',
                    'Scans all .go files in apps/api/internal/models/',
                    'Parses each file using Go\'s AST (Abstract Syntax Tree) parser to extract struct definitions',
                    'For each struct, reads field names, Go types, JSON tags, and GORM tags',
                    'Maps Go types to TypeScript types and Zod validators',
                    'Writes TypeScript interface files to packages/shared/types/',
                    'Writes Zod schema files to packages/shared/schemas/',
                    'Skips the User model (which has custom hand-written schemas)',
                  ].map((item, i) => (
                    <li key={i} className="text-[14px] text-muted-foreground pl-1">
                      {item}
                    </li>
                  ))}
                </ol>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  When to Use It
                </h3>
                <ul className="space-y-2 mb-4">
                  {[
                    'After manually adding or removing fields from a Go model',
                    'After changing a field\'s type in a Go struct',
                    'After modifying GORM tags (e.g., adding type:text)',
                    'After adding a completely new model file manually (without using grit generate)',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-1">&#10003;</span>
                      {item}
                    </li>
                  ))}
                </ul>
                <p className="text-sm text-muted-foreground/60">
                  Note: You do not need to run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit sync</code> after
                  using <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>.
                  The generator already creates the TypeScript types and Zod schemas for the new resource.
                </p>
              </div>

              {/* grit update */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  grit update
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Update the Grit CLI itself to the latest version. It first checks GitHub for the
                  newest release and exits immediately if you are already up to date. Otherwise it
                  picks a strategy: if the Go toolchain is on your PATH it runs{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">go install ...@&lt;latest&gt;</code>;
                  if Go is not installed it downloads the matching prebuilt binary from the GitHub
                  release and atomically swaps it in. On Windows the running{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.exe</code> is locked, so
                  the old binary is renamed to <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.old</code>{" "}
                  before the new one is written (and restored if the update fails).
                </p>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit update</span></div>
                    <div className="mt-1 text-muted-foreground/40 text-xs space-y-0.5">
                      <div>  Grit self-update — current: v3.55.0</div>
                      <div>  → Checking GitHub for the latest release...</div>
                      <div>  → New version available: v3.55.0 → v3.56.0</div>
                      <div>  → Running: go install github.com/MUKE-coder/grit/v3/cmd/grit@v3.56.0</div>
                      <div className="text-primary/60">  ✓ Updated to v3.56.0</div>
                    </div>
                  </div>
                </div>
                <p className="text-sm text-muted-foreground/60 mt-3">
                  Aliased as <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit self-update</code>.
                  Pass <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--from-release</code> to skip
                  the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">go install</code> path and
                  always pull the prebuilt binary from the GitHub release.
                </p>
                <div className="mt-4 rounded-lg border border-border/20 bg-accent/20 p-3">
                  <p className="text-sm text-muted-foreground/70">
                    <strong className="text-foreground/80">Note:</strong>{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit update</code> updates
                    the <strong className="text-foreground/80">CLI tool</strong> itself.
                    To update your <strong className="text-foreground/80">project&apos;s scaffold files</strong> (admin
                    panel, configs, web app), use{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit upgrade</code> instead.
                  </p>
                </div>
              </div>

              {/* grit version */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  grit version
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Print the current version of the Grit CLI.
                </p>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit version</span></div>
                    <div className="mt-1"><span className="text-muted-foreground/60">grit version 3.55.0</span></div>
                  </div>
                </div>
              </div>

              {/* Operational Commands */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Operational Commands
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-6">
                  Grit ships a set of operational commands inspired by Laravel/Goravel for
                  day-to-day workflows: route inspection, maintenance mode, and one-command
                  deployment.
                </p>

                {/* grit routes */}
                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  grit routes
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  List all registered API routes in a formatted table. Parses your{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">routes.go</code> file
                  and shows the HTTP method, path, handler function, and middleware group.
                </p>
                <CodeBlock language="bash" filename="Terminal" code={`$ grit routes

  METHOD  PATH                              HANDLER                GROUP
  ──────  ────                              ───────                ─────
  GET     /api/health                       func1                  public
  POST    /api/auth/register                authHandler.Register   public
  POST    /api/auth/login                   authHandler.Login      public
  POST    /api/auth/refresh                 authHandler.Refresh    public
  GET     /api/auth/me                      authHandler.Me         protected
  POST    /api/auth/logout                  authHandler.Logout     protected
  POST    /api/auth/totp/setup              totpHandler.Setup      protected
  GET     /api/auth/totp/status             totpHandler.Status     protected
  GET     /api/users/:id                    userHandler.GetByID    protected
  POST    /api/uploads                      uploadHandler.Create   protected
  POST    /api/ai/chat                      aiHandler.Chat         protected
  DELETE  /api/admin/users/:id              userHandler.Delete     admin

  16 routes total`} />
                <p className="text-sm text-muted-foreground/60 mt-3">
                  Works for both monorepo (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/api/internal/routes/</code>) and
                  single app (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">internal/routes/</code>) projects.
                </p>

                {/* grit down / grit up */}
                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  grit down / grit up
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Toggle maintenance mode. When enabled, all API requests receive a{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">503 Service Unavailable</code> response.
                </p>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                  <CodeBlock language="bash" filename="Enable maintenance" code={`$ grit down

  Application is now in maintenance mode.
  All requests will receive 503.
  Run 'grit up' to bring it back online.`} />
                  <CodeBlock language="bash" filename="Disable maintenance" code={`$ grit up

  Application is back online!
  Normal request handling has resumed.`} />
                </div>

                <div className="rounded-lg border border-border/40 bg-accent/20 p-5 mb-6">
                  <h4 className="text-sm font-semibold text-foreground mb-2">How it works</h4>
                  <ul className="space-y-2 text-sm text-muted-foreground">
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">1.</span>
                      <span><code className="text-foreground/80">grit down</code> creates a <code className="text-foreground/80">.maintenance</code> file in the project root</span>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">2.</span>
                      <span>The scaffolded <code className="text-foreground/80">Maintenance()</code> middleware checks for this file on every request</span>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">3.</span>
                      <span><code className="text-foreground/80">grit up</code> removes the file, resuming normal operation</span>
                    </li>
                  </ul>
                </div>

                <CodeBlock language="go" filename="middleware/maintenance.go" code={`func Maintenance() gin.HandlerFunc {
    return func(c *gin.Context) {
        if _, err := os.Stat(".maintenance"); err == nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{
                "error": gin.H{
                    "code":    "MAINTENANCE",
                    "message": "Application is in maintenance mode.",
                },
            })
            c.Abort()
            return
        }
        c.Next()
    }
}`} highlightLines={[3, 4, 5, 6, 7, 8, 9]} />

                {/* grit deploy */}
                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  grit deploy
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  One-command production deployment. See the dedicated{' '}
                  <Link href="/docs/infrastructure/deploy-command" className="text-primary hover:underline">
                    Deploy Command
                  </Link>{' '}
                  guide for full details.
                </p>
                <CodeBlock language="bash" filename="Terminal" code={`# Deploy with flags
grit deploy --host user@server.com --domain myapp.com

# Or set env vars in .env
DEPLOY_HOST=user@server.com
DEPLOY_DOMAIN=myapp.com
DEPLOY_KEY_FILE=~/.ssh/id_rsa

grit deploy`} />
              </div>

              {/* Command Reference Table */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Quick Reference
                </h2>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Command</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit init</td>
                        <td className="px-4 py-2.5">Write CLAUDE.md / AGENTS.md convention docs to the project root</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit new &lt;name&gt;</td>
                        <td className="px-4 py-2.5">Scaffold a new project (interactive by default)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit new-desktop &lt;name&gt;</td>
                        <td className="px-4 py-2.5">Scaffold a standalone Wails desktop app</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit new .</td>
                        <td className="px-4 py-2.5">Scaffold into the current directory</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit new . --force</td>
                        <td className="px-4 py-2.5">Scaffold into a non-empty directory</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit new myapp --here</td>
                        <td className="px-4 py-2.5">Explicit in-place scaffolding</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit new &lt;name&gt; --api</td>
                        <td className="px-4 py-2.5">Scaffold Go API only</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit new &lt;name&gt; --full</td>
                        <td className="px-4 py-2.5">Scaffold everything including docs</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit new &lt;name&gt; --single</td>
                        <td className="px-4 py-2.5">Single architecture (API only)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit new &lt;name&gt; --double</td>
                        <td className="px-4 py-2.5">Double architecture (API + web)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit new &lt;name&gt; --triple</td>
                        <td className="px-4 py-2.5">Triple architecture (API + web + admin)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit new &lt;name&gt; --vite</td>
                        <td className="px-4 py-2.5">Use TanStack Router (Vite) frontend</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit new &lt;name&gt; --next</td>
                        <td className="px-4 py-2.5">Use Next.js frontend</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit new &lt;name&gt; --desktop</td>
                        <td className="px-4 py-2.5">Add Wails desktop client (combinable with --triple/--double/--mobile/--api)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit new &lt;name&gt; --mobile --desktop</td>
                        <td className="px-4 py-2.5">Multi-client: API + Expo mobile + Wails desktop (shared types)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">--arch, --frontend</td>
                        <td className="px-4 py-2.5">Long-form flags for architecture and frontend</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit generate resource &lt;Name&gt;</td>
                        <td className="px-4 py-2.5">Generate full-stack CRUD resource</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit g resource &lt;Name&gt;</td>
                        <td className="px-4 py-2.5">Shorthand for generate resource</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit generate seeder &lt;Resource&gt;</td>
                        <td className="px-4 py-2.5">Generate a database seeder for existing resources</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit generate sequence &lt;Name&gt;</td>
                        <td className="px-4 py-2.5">Generate a sequential numbering helper (e.g. INV-202605-0001)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit remove resource &lt;Name&gt;</td>
                        <td className="px-4 py-2.5">Remove a generated resource and clean up markers</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit add role &lt;ROLE&gt;</td>
                        <td className="px-4 py-2.5">Add a new role across all project files</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit add web-auth</td>
                        <td className="px-4 py-2.5">Add page protection helpers to apps/web/</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit expose form &lt;Resource&gt;</td>
                        <td className="px-4 py-2.5">Scaffold a public form page for a resource</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit expose table &lt;Resource&gt;</td>
                        <td className="px-4 py-2.5">Scaffold a paginated list page for a resource</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit start</td>
                        <td className="px-4 py-2.5">Start every app in the project in parallel</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit start client</td>
                        <td className="px-4 py-2.5">Start frontend apps via pnpm dev</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit start server</td>
                        <td className="px-4 py-2.5">Start Go API server (hot-reload via air)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit compile</td>
                        <td className="px-4 py-2.5">Build the desktop app executable (Wails)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit package</td>
                        <td className="px-4 py-2.5">Build a distributable desktop installer (.exe / .app / binary)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit studio</td>
                        <td className="px-4 py-2.5">Open the GORM Studio database browser</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit sync</td>
                        <td className="px-4 py-2.5">Sync Go models to TypeScript + Zod</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit migrate</td>
                        <td className="px-4 py-2.5">Run GORM AutoMigrate for all models</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit migrate --fresh</td>
                        <td className="px-4 py-2.5">Drop all tables then re-migrate</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit seed</td>
                        <td className="px-4 py-2.5">Populate database with initial data</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit backup</td>
                        <td className="px-4 py-2.5">Back up the entire database to a ZIP archive</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit restore &lt;backup.zip&gt;</td>
                        <td className="px-4 py-2.5">Restore the database from a backup archive</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit upgrade</td>
                        <td className="px-4 py-2.5">Update project scaffold files to latest</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit update</td>
                        <td className="px-4 py-2.5">Remove old CLI and install latest version</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit routes</td>
                        <td className="px-4 py-2.5">List all registered API routes</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit down</td>
                        <td className="px-4 py-2.5">Enable maintenance mode (503)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit up</td>
                        <td className="px-4 py-2.5">Disable maintenance mode</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit deploy</td>
                        <td className="px-4 py-2.5">Deploy to production server</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">grit version</td>
                        <td className="px-4 py-2.5">Print CLI version</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>

            {/* Bug Reports */}
            <div className="mb-10 rounded-xl border border-primary/20 bg-primary/5 p-5">
              <p className="text-[15px] text-muted-foreground leading-relaxed">
                Found a bug? Open an issue at{' '}
                <a href="https://github.com/MUKE-coder/grit/issues" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">
                  https://github.com/MUKE-coder/grit/issues
                </a>{' '}
                — we fix bugs fast and appreciate every report.
              </p>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/concepts/architecture" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Architecture
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/concepts/code-generation" className="gap-1.5">
                  Code Generation
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
