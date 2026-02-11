import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'

export default function CLICommandsPage() {
  return (
    <div className="min-h-screen bg-background">
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
                resources, and syncs types between Go and TypeScript. Install it once and use it
                across all your Grit projects.
              </p>
            </div>

            <div className="prose-grit">
              {/* Installation */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Installing the CLI
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Install the Grit CLI globally using Go:
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
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">go install github.com/MUKE-coder/grit/cmd/grit@latest</span></div>
                  </div>
                </div>
                <p className="text-sm text-muted-foreground/60 mt-3">
                  This installs the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit</code> binary
                  to your <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">$GOPATH/bin</code>.
                  Make sure that directory is in your system PATH.
                </p>
              </div>

              {/* grit new */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  grit new
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Scaffold a new Grit project. This creates the full monorepo directory structure,
                  initializes all configuration files, and sets up the Go API, Next.js frontends,
                  shared package, and Docker infrastructure.
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
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit new my-saas-app</span></div>
                  </div>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Syntax
                </h3>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">usage</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`grit new <project-name> [flags]`}</pre>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Flags
                </h3>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Flag</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">(no flag)</td>
                        <td className="px-4 py-2.5">Default. Creates API + web app + admin panel + shared package + Docker</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">--api</td>
                        <td className="px-4 py-2.5">Go API only. No Next.js frontends, no shared package. Ideal for pure backend projects</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">--expo</td>
                        <td className="px-4 py-2.5">Include an Expo (React Native) mobile app alongside the web and admin apps</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">--mobile</td>
                        <td className="px-4 py-2.5">API + Expo mobile app only. No web or admin frontends</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">--full</td>
                        <td className="px-4 py-2.5">Everything including a documentation site app</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <p className="text-sm text-muted-foreground/60 mt-3">
                  Only one mode flag can be used at a time. Using multiple flags together (e.g.,
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--api --expo</code>) will
                  produce an error.
                </p>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Project Name Validation
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The project name must follow these rules:
                </p>
                <ul className="space-y-2 mb-4">
                  {[
                    'Lowercase letters, numbers, and hyphens only',
                    'Must start with a letter',
                    'Cannot end with a hyphen',
                    'No consecutive hyphens',
                    'Minimum 2 characters, maximum 64 characters',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-1">&#10003;</span>
                      {item}
                    </li>
                  ))}
                </ul>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  What Gets Created
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  With the default flags (no <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--api</code>),
                  the following files and directories are scaffolded:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Category</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Files</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">Root</td>
                        <td className="px-4 py-2.5 text-xs">.env, .env.example, .gitignore, turbo.json, pnpm-workspace.yaml, README.md, docker-compose.yml, docker-compose.prod.yml</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">Go API</td>
                        <td className="px-4 py-2.5 text-xs">go.mod, main.go, config, database, models (User), handlers (auth, user), services (auth), middleware (auth, cors, logger), routes, Dockerfile</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">Web App</td>
                        <td className="px-4 py-2.5 text-xs">package.json, next.config, tailwind.config, layout, auth pages (login, register, forgot-password), dashboard, sidebar, hooks, api-client, Dockerfile</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">Admin Panel</td>
                        <td className="px-4 py-2.5 text-xs">package.json, layout, dashboard with stats, users management page, sidebar, navbar, data table, hooks, resource definitions</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">Shared</td>
                        <td className="px-4 py-2.5 text-xs">package.json, schemas (user.ts), types (user.ts, api.ts), constants (index.ts)</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Examples
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
                      <span className="text-muted-foreground/40"># Full-stack monorepo (default)</span>
                      <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit new my-crm</span></div>
                    </div>
                    <div>
                      <span className="text-muted-foreground/40"># Go API only (no frontend)</span>
                      <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit new billing-service --api</span></div>
                    </div>
                    <div>
                      <span className="text-muted-foreground/40"># Full stack with Expo mobile app</span>
                      <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit new my-app --expo</span></div>
                    </div>
                    <div>
                      <span className="text-muted-foreground/40"># Everything including docs site</span>
                      <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit new my-framework --full</span></div>
                    </div>
                  </div>
                </div>
              </div>

              {/* grit generate resource */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  grit generate resource
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Generate a complete full-stack CRUD resource. This is the most powerful command
                  in the Grit CLI. It creates 8 new files and modifies up to 10 existing files to
                  wire everything together automatically.
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
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit generate resource Post --fields &quot;title:string,content:text,published:bool&quot;</span></div>
                  </div>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Syntax
                </h3>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">usage</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`grit generate resource <Name> [flags]

# Shorthand alias
grit g resource <Name> [flags]`}</pre>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Flags
                </h3>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Flag</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">--fields</td>
                        <td className="px-4 py-2.5">Inline field definitions as comma-separated name:type pairs</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">--from &lt;file&gt;</td>
                        <td className="px-4 py-2.5">Load field definitions from a YAML file</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">-i, --interactive</td>
                        <td className="px-4 py-2.5">Define fields interactively in the terminal, one at a time</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <p className="text-sm text-muted-foreground/60 mt-3">
                  You must provide exactly one of <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--fields</code>,
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--from</code>, or
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">-i</code>. Running the
                  command without any of these flags prints a usage example.
                </p>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Field Syntax
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Fields are defined as <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">name:type</code> pairs.
                  The name should be a simple identifier (e.g., <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">title</code>,
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">due_date</code>,
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">is_active</code>).
                  The generator automatically converts names to the appropriate case for each language.
                </p>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Field Type Reference
                </h3>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Type</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Go Type</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">TS Type</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Zod Type</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Form</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">string</td>
                        <td className="px-4 py-2.5 font-mono text-xs">string</td>
                        <td className="px-4 py-2.5 font-mono text-xs">string</td>
                        <td className="px-4 py-2.5 font-mono text-xs">{`z.string().min(1)`}</td>
                        <td className="px-4 py-2.5 font-mono text-xs">text input</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">text</td>
                        <td className="px-4 py-2.5 font-mono text-xs">string</td>
                        <td className="px-4 py-2.5 font-mono text-xs">string</td>
                        <td className="px-4 py-2.5 font-mono text-xs">z.string()</td>
                        <td className="px-4 py-2.5 font-mono text-xs">textarea</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">int</td>
                        <td className="px-4 py-2.5 font-mono text-xs">int</td>
                        <td className="px-4 py-2.5 font-mono text-xs">number</td>
                        <td className="px-4 py-2.5 font-mono text-xs">z.number().int()</td>
                        <td className="px-4 py-2.5 font-mono text-xs">number input</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">uint</td>
                        <td className="px-4 py-2.5 font-mono text-xs">uint</td>
                        <td className="px-4 py-2.5 font-mono text-xs">number</td>
                        <td className="px-4 py-2.5 font-mono text-xs">{`z.number().int().nonnegative()`}</td>
                        <td className="px-4 py-2.5 font-mono text-xs">number input</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">float</td>
                        <td className="px-4 py-2.5 font-mono text-xs">float64</td>
                        <td className="px-4 py-2.5 font-mono text-xs">number</td>
                        <td className="px-4 py-2.5 font-mono text-xs">z.number()</td>
                        <td className="px-4 py-2.5 font-mono text-xs">number input</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">bool</td>
                        <td className="px-4 py-2.5 font-mono text-xs">bool</td>
                        <td className="px-4 py-2.5 font-mono text-xs">boolean</td>
                        <td className="px-4 py-2.5 font-mono text-xs">z.boolean()</td>
                        <td className="px-4 py-2.5 font-mono text-xs">toggle</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">datetime</td>
                        <td className="px-4 py-2.5 font-mono text-xs">*time.Time</td>
                        <td className="px-4 py-2.5 font-mono text-xs">string | null</td>
                        <td className="px-4 py-2.5 font-mono text-xs">z.string().nullable()</td>
                        <td className="px-4 py-2.5 font-mono text-xs">datetime picker</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">date</td>
                        <td className="px-4 py-2.5 font-mono text-xs">*time.Time</td>
                        <td className="px-4 py-2.5 font-mono text-xs">string | null</td>
                        <td className="px-4 py-2.5 font-mono text-xs">z.string().nullable()</td>
                        <td className="px-4 py-2.5 font-mono text-xs">date picker</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <p className="text-sm text-muted-foreground/60 mt-3">
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">string</code> type
                  maps to a GORM column with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">size:255</code> and
                  is required by default. The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">text</code> type
                  maps to a GORM <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">type:text</code> column
                  and is optional by default. Both <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">datetime</code> and
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">date</code> use Go
                  pointer types (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">*time.Time</code>) to
                  allow null values.
                </p>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Using --from (YAML Definition)
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  For resources with many fields or complex configurations, use a YAML definition file:
                </p>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">post.yaml</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`name: Post
fields:
  - name: title
    type: string
    required: true
    unique: true
  - name: content
    type: text
  - name: excerpt
    type: string
  - name: published
    type: bool
    default: "false"
  - name: views
    type: int
  - name: published_at
    type: datetime`}</pre>
                </div>
                <div className="mt-4 rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm">
                    <div><span className="text-primary/50 select-none">$ </span><span className="text-foreground/80">grit generate resource Post --from post.yaml</span></div>
                  </div>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Using -i (Interactive Mode)
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Interactive mode prompts you for fields one at a time. Enter each field
                  as <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">name:type</code> and
                  press Enter. Press Enter on an empty line when done.
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
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`$ grit generate resource Invoice -i

  Defining fields for Invoice
  Enter fields as name:type (e.g., title:string)
  Valid types: string, text, int, uint, float, bool, datetime, date
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
  >`}</pre>
                </div>

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
                  </div>
                </div>
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
                    <div className="mt-1"><span className="text-muted-foreground/60">grit version 0.1.0</span></div>
                  </div>
                </div>
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
                        <td className="px-4 py-2.5 font-mono text-xs">grit new &lt;name&gt;</td>
                        <td className="px-4 py-2.5">Scaffold a new full-stack project</td>
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
                        <td className="px-4 py-2.5 font-mono text-xs">grit generate resource &lt;Name&gt;</td>
                        <td className="px-4 py-2.5">Generate full-stack CRUD resource</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit g resource &lt;Name&gt;</td>
                        <td className="px-4 py-2.5">Shorthand for generate resource</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit sync</td>
                        <td className="px-4 py-2.5">Sync Go models to TypeScript + Zod</td>
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
