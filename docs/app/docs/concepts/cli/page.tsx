import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/cli')

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
                  Update the Grit CLI itself to the latest version. This removes the current binary
                  and installs the newest release from GitHub using <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">go install</code>.
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
                      <div>  → Removing old binary: /home/user/go/bin/grit</div>
                      <div>  → Installing latest version...</div>
                      <div className="text-primary/60">  ✓ Grit CLI updated successfully!</div>
                    </div>
                  </div>
                </div>
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
                    <div className="mt-1"><span className="text-muted-foreground/60">grit version 0.13.0</span></div>
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
                        <td className="px-4 py-2.5 font-mono text-xs">grit remove resource &lt;Name&gt;</td>
                        <td className="px-4 py-2.5">Remove a generated resource and clean up markers</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit add role &lt;ROLE&gt;</td>
                        <td className="px-4 py-2.5">Add a new role across all project files</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit start client</td>
                        <td className="px-4 py-2.5">Start frontend apps via pnpm dev</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit start server</td>
                        <td className="px-4 py-2.5">Start Go API server</td>
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
                        <td className="px-4 py-2.5 font-mono text-xs">grit upgrade</td>
                        <td className="px-4 py-2.5">Update project scaffold files to latest</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">grit update</td>
                        <td className="px-4 py-2.5">Remove old CLI and install latest version</td>
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
