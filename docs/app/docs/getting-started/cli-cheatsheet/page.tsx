import Link from "next/link";
import { ArrowLeft, ArrowRight, Download, Terminal } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from "@/components/code-block";
import { getDocMetadata } from "@/config/docs-metadata";

export const metadata = getDocMetadata("/docs/getting-started/cli-cheatsheet");

/* ------------------------------------------------------------------ */
/*  Reusable terminal card (matches Docker Cheat Sheet style)         */
/* ------------------------------------------------------------------ */
function TerminalCard({
  cmd,
  desc,
}: {
  cmd: string;
  desc: string;
}) {
  return (
    <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
      <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
        <div className="flex items-center gap-1.5">
          <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
          <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
          <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
        </div>
        <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">
          terminal
        </span>
      </div>
      <div className="px-5 py-3 font-mono text-sm">
        <div>
          <span className="text-primary/50 select-none">$ </span>
          <span className="text-foreground/80">{cmd}</span>
        </div>
        <p className="text-xs text-muted-foreground/50 mt-1">{desc}</p>
      </div>
    </div>
  );
}

/* ------------------------------------------------------------------ */
/*  Page                                                              */
/* ------------------------------------------------------------------ */
export default function CLICheatsheetPage() {
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
                CLI Cheatsheet
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed mb-6">
                Every command the Grit CLI offers in one place. Bookmark this
                page or print it out &mdash; it&apos;s your pocket reference for
                scaffolding projects, generating resources, running migrations,
                and everything in between.
              </p>
              <a
                href="https://14j7oh8kso.ufs.sh/f/HLxTbDBCDLwfeHHJl34ZKSqNhOvVj6p9rg3Icmo05TAEwQ4a"
                target="_blank"
                rel="noopener noreferrer"
              >
                <Button size="sm" className="gap-1.5">
                  <Download className="h-3.5 w-3.5" />
                  Download Grit Handbook PDF
                </Button>
              </a>
            </div>

            {/* Quick-reference table */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Command Overview
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                All top-level commands at a glance. Scroll down for detailed
                usage, flags, and examples.
              </p>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                        Command
                      </th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                        Alias
                      </th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                        Description
                      </th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    {[
                      {
                        cmd: "grit init",
                        alias: "",
                        desc: "Write CLAUDE.md / AGENTS.md convention docs",
                      },
                      {
                        cmd: "grit new",
                        alias: "",
                        desc: "Scaffold a new project",
                      },
                      {
                        cmd: "grit new-desktop",
                        alias: "",
                        desc: "Scaffold a standalone Wails desktop app",
                      },
                      {
                        cmd: "grit generate",
                        alias: "g",
                        desc: "Generate resources, seeders, sequences",
                      },
                      {
                        cmd: "grit remove",
                        alias: "rm",
                        desc: "Remove a generated resource",
                      },
                      {
                        cmd: "grit add",
                        alias: "",
                        desc: "Add roles / web-auth helpers",
                      },
                      {
                        cmd: "grit expose",
                        alias: "",
                        desc: "Scaffold a web page for a resource (form/table)",
                      },
                      {
                        cmd: "grit start",
                        alias: "",
                        desc: "Start dev servers",
                      },
                      {
                        cmd: "grit compile",
                        alias: "",
                        desc: "Build desktop app executable (Wails)",
                      },
                      {
                        cmd: "grit package",
                        alias: "",
                        desc: "Build a distributable desktop installer",
                      },
                      {
                        cmd: "grit studio",
                        alias: "",
                        desc: "Open the GORM Studio database browser",
                      },
                      {
                        cmd: "grit sync",
                        alias: "",
                        desc: "Sync Go types → TypeScript",
                      },
                      {
                        cmd: "grit migrate",
                        alias: "",
                        desc: "Run database migrations",
                      },
                      {
                        cmd: "grit seed",
                        alias: "",
                        desc: "Seed the database",
                      },
                      {
                        cmd: "grit backup",
                        alias: "",
                        desc: "Back up the entire database",
                      },
                      {
                        cmd: "grit restore",
                        alias: "",
                        desc: "Restore the database from a backup archive",
                      },
                      {
                        cmd: "grit routes",
                        alias: "",
                        desc: "List all registered API routes",
                      },
                      {
                        cmd: "grit down",
                        alias: "",
                        desc: "Enter maintenance mode (503)",
                      },
                      {
                        cmd: "grit up",
                        alias: "",
                        desc: "Exit maintenance mode",
                      },
                      {
                        cmd: "grit deploy",
                        alias: "",
                        desc: "Deploy to a remote server",
                      },
                      {
                        cmd: "grit upgrade",
                        alias: "",
                        desc: "Upgrade project templates",
                      },
                      {
                        cmd: "grit update",
                        alias: "self-update",
                        desc: "Update the CLI binary",
                      },
                      {
                        cmd: "grit version",
                        alias: "",
                        desc: "Print CLI version",
                      },
                    ].map((row) => (
                      <tr key={row.cmd} className="border-b border-border/20 last:border-0">
                        <td className="px-4 py-2.5 font-mono text-xs text-foreground/90 font-medium">
                          {row.cmd}
                        </td>
                        <td className="px-4 py-2.5 font-mono text-xs text-primary/60">
                          {row.alias || "—"}
                        </td>
                        <td className="px-4 py-2.5 text-muted-foreground/70">
                          {row.desc}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>

            <div className="prose-grit">
              {/* -------------------------------------------------------- */}
              {/*  grit init                                               */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit init
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Write the framework&apos;s hard-rules convention docs to the
                  project root as{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    CLAUDE.md
                  </code>{" "}
                  and{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    AGENTS.md
                  </code>{" "}
                  (same content &mdash; different AI tools look for different
                  filenames). Skips files that already exist unless{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    --force
                  </code>{" "}
                  is passed.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit init"
                    desc="Write CLAUDE.md and AGENTS.md (skips existing files)"
                  />
                  <TerminalCard
                    cmd="grit init --force"
                    desc="Overwrite the convention docs (e.g. after a major upgrade)"
                  />
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit new                                                */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit new
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Scaffold a brand-new Grit monorepo with the Go API, Next.js
                  web app, admin panel, shared packages, Docker configs, and all
                  the batteries.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit new myapp"
                    desc="Create a full-stack project (API + web + admin + shared)"
                  />
                  <TerminalCard
                    cmd="grit new myapp --api"
                    desc="Scaffold only the Go API (no frontend apps)"
                  />
                  <TerminalCard
                    cmd="grit new myapp --expo"
                    desc="Full stack + Expo mobile app"
                  />
                  <TerminalCard
                    cmd="grit new myapp --mobile"
                    desc="API + Expo mobile app only (no web/admin)"
                  />
                  <TerminalCard
                    cmd="grit new myapp --mobile --desktop"
                    desc="API + Expo + Wails desktop — the multi-client pattern (new in v3.9)"
                  />
                  <TerminalCard
                    cmd="grit new myapp --triple --next --desktop"
                    desc="API + web + admin + desktop, one monorepo, shared types"
                  />
                  <TerminalCard
                    cmd="grit new myapp --full"
                    desc="Scaffold everything including docs site"
                  />
                  <TerminalCard
                    cmd='grit new myapp --style modern'
                    desc='Set admin panel style variant (default, modern, minimal, glass, centered)'
                  />
                </div>

                {/* Flags table */}
                <div className="mt-6 rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <div className="px-4 py-2.5 border-b border-border/30 bg-accent/20">
                    <span className="text-xs font-mono font-medium text-foreground/70">
                      FLAGS
                    </span>
                  </div>
                  <table className="w-full text-sm">
                    <tbody className="text-muted-foreground">
                      {[
                        {
                          flag: "--api",
                          type: "bool",
                          desc: "Scaffold only the Go API",
                        },
                        {
                          flag: "--expo",
                          type: "bool",
                          desc: "Include Expo mobile app",
                        },
                        {
                          flag: "--mobile",
                          type: "bool",
                          desc: "API + Expo mobile only",
                        },
                        {
                          flag: "--desktop",
                          type: "bool",
                          desc: "Add a Wails desktop client (combinable with --triple/--double/--mobile/--api; not --single)",
                        },
                        {
                          flag: "--full",
                          type: "bool",
                          desc: "Include docs site",
                        },
                        {
                          flag: "--style",
                          type: "string",
                          desc: "Admin style: default, modern, minimal, glass, centered",
                        },
                      ].map((f) => (
                        <tr key={f.flag} className="border-b border-border/15 last:border-0">
                          <td className="px-4 py-2 font-mono text-xs text-primary/70 w-28">
                            {f.flag}
                          </td>
                          <td className="px-4 py-2 font-mono text-xs text-muted-foreground/50 w-16">
                            {f.type}
                          </td>
                          <td className="px-4 py-2 text-muted-foreground/70 text-sm">
                            {f.desc}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                <div className="p-4 rounded-lg border border-primary/15 bg-primary/5 mt-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    The project name must be lowercase, alphanumeric, and hyphens
                    only (e.g.{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      my-saas-app
                    </code>
                    ). It must start with a letter and cannot end with a hyphen.
                    Only one architecture flag (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--api</code>,{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--single</code>,{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--double</code>,{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--triple</code>,{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--mobile</code>,{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--full</code>)
                    can be used at a time.{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--desktop</code>{" "}
                    and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--expo</code>{" "}
                    are additive — combine them with any architecture (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--desktop</code> is not supported with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--single</code>).
                  </p>
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit generate resource                                  */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit generate resource
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-1">
                  <span className="text-xs font-mono text-primary/60 bg-primary/10 rounded-full px-2 py-0.5 mr-2">
                    alias: grit g resource
                  </span>
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4 mt-3">
                  Generate a complete full-stack CRUD resource: Go model,
                  handler, service, Zod schemas, TypeScript types, React Query
                  hooks, and an admin page &mdash; all wired together
                  automatically.
                </p>

                <div className="space-y-3">
                  <TerminalCard
                    cmd='grit generate resource Post --fields "title:string,content:text,published:bool"'
                    desc="Generate a resource with inline field definitions"
                  />
                  <TerminalCard
                    cmd='grit g resource Post --fields "title:string,slug:string:unique,views:int"'
                    desc="Use the short alias and add a unique constraint"
                  />
                  <TerminalCard
                    cmd="grit generate resource Post --from post.yaml"
                    desc="Generate from a YAML field definition file"
                  />
                  <TerminalCard
                    cmd="grit generate resource Post -i"
                    desc="Interactive mode — define fields with prompts"
                  />
                  <TerminalCard
                    cmd='grit g resource Post --fields "title:string,content:text" --roles "ADMIN,EDITOR"'
                    desc="Restrict generated routes to specific roles"
                  />
                  <TerminalCard
                    cmd='grit g resource Post --fields "title:string,views:int" --faker --count 50'
                    desc="Also generate a seeder that inserts 50 fake rows"
                  />
                </div>

                {/* Flags table */}
                <div className="mt-6 rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <div className="px-4 py-2.5 border-b border-border/30 bg-accent/20">
                    <span className="text-xs font-mono font-medium text-foreground/70">
                      FLAGS
                    </span>
                  </div>
                  <table className="w-full text-sm">
                    <tbody className="text-muted-foreground">
                      {[
                        {
                          flag: "--fields",
                          type: "string",
                          desc: 'Inline fields (e.g. "title:string,published:bool")',
                        },
                        {
                          flag: "--from",
                          type: "string",
                          desc: "Path to YAML file defining the resource fields",
                        },
                        {
                          flag: "-i, --interactive",
                          type: "bool",
                          desc: "Interactively define fields via prompts",
                        },
                        {
                          flag: "--roles",
                          type: "string",
                          desc: 'Restrict routes to roles (e.g. "ADMIN,EDITOR")',
                        },
                        {
                          flag: "--seed",
                          type: "bool",
                          desc: "Also generate a seeder with one example record",
                        },
                        {
                          flag: "--faker",
                          type: "bool",
                          desc: "Also generate a gofakeit seeder (implies --seed)",
                        },
                        {
                          flag: "--count",
                          type: "int",
                          desc: "Number of rows for the faker seeder (default 10)",
                        },
                      ].map((f) => (
                        <tr key={f.flag} className="border-b border-border/15 last:border-0">
                          <td className="px-4 py-2 font-mono text-xs text-primary/70 w-40">
                            {f.flag}
                          </td>
                          <td className="px-4 py-2 font-mono text-xs text-muted-foreground/50 w-16">
                            {f.type}
                          </td>
                          <td className="px-4 py-2 text-muted-foreground/70 text-sm">
                            {f.desc}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                {/* Generated files */}
                <div className="mt-6">
                  <h3 className="text-base font-semibold text-foreground/90 mb-3">
                    Generated Files
                  </h3>
                  <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                    <table className="w-full text-sm">
                      <tbody className="text-muted-foreground">
                        {[
                          { file: "apps/api/internal/models/<name>.go", what: "GORM model with struct tags" },
                          { file: "apps/api/internal/handlers/<name>.go", what: "Full CRUD handler with pagination" },
                          { file: "apps/api/internal/services/<name>.go", what: "Business logic layer" },
                          { file: "packages/shared/schemas/<name>.ts", what: "Zod validation schemas" },
                          { file: "packages/shared/types/<name>.ts", what: "TypeScript types" },
                          { file: "apps/admin/hooks/use-<names>.ts", what: "React Query hooks" },
                          { file: "apps/admin/app/resources/<names>/page.tsx", what: "Admin page with data table" },
                        ].map((r) => (
                          <tr key={r.file} className="border-b border-border/15 last:border-0">
                            <td className="px-4 py-2 font-mono text-xs text-foreground/70">
                              {r.file}
                            </td>
                            <td className="px-4 py-2 text-muted-foreground/60 text-sm">
                              {r.what}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                  <p className="text-sm text-muted-foreground/50 mt-2">
                    Routes are auto-registered in{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      routes.go
                    </code>
                    , the model is added to auto-migrations, and the resource is
                    injected into the admin sidebar.
                  </p>
                </div>

                {/* Field types */}
                <div className="mt-6">
                  <h3 className="text-base font-semibold text-foreground/90 mb-3">
                    Supported Field Types
                  </h3>
                  <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                    <table className="w-full text-sm">
                      <thead>
                        <tr className="border-b border-border/30 bg-accent/20">
                          <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                            Type
                          </th>
                          <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                            Go Type
                          </th>
                          <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                            TS Type
                          </th>
                          <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                            Form Field
                          </th>
                        </tr>
                      </thead>
                      <tbody className="text-muted-foreground">
                        {[
                          { type: "string", go: "string", ts: "string", form: "text" },
                          { type: "text", go: "string", ts: "string", form: "textarea" },
                          { type: "richtext", go: "string", ts: "string", form: "richtext" },
                          { type: "int", go: "int", ts: "number", form: "number" },
                          { type: "uint", go: "uint", ts: "number", form: "number" },
                          { type: "float", go: "float64", ts: "number", form: "number" },
                          { type: "bool", go: "bool", ts: "boolean", form: "toggle" },
                          { type: "datetime", go: "*time.Time", ts: "string | null", form: "datetime" },
                          { type: "date", go: "*time.Time", ts: "string | null", form: "date" },
                          { type: "slug", go: "string", ts: "string", form: "auto" },
                          { type: "belongs_to", go: "uint", ts: "number", form: "select" },
                          { type: "many_to_many", go: "[]uint", ts: "number[]", form: "multi-select" },
                          { type: "string_array", go: "JSONSlice[string]", ts: "string[]", form: "images" },
                        ].map((f) => (
                          <tr key={f.type} className="border-b border-border/15 last:border-0">
                            <td className="px-4 py-2 font-mono text-xs text-primary/70 font-medium">
                              {f.type}
                            </td>
                            <td className="px-4 py-2 font-mono text-xs text-foreground/70">
                              {f.go}
                            </td>
                            <td className="px-4 py-2 font-mono text-xs text-foreground/70">
                              {f.ts}
                            </td>
                            <td className="px-4 py-2 text-sm text-muted-foreground/60">
                              {f.form}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>

                {/* Inline field syntax */}
                <div className="mt-6">
                  <h3 className="text-base font-semibold text-foreground/90 mb-3">
                    Inline Field Syntax
                  </h3>
                  <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                    Fields are comma-separated. Each field follows the pattern{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      name:type
                    </code>{" "}
                    or{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      name:type:modifier
                    </code>
                    . Available modifiers:
                  </p>
                  <div className="grid gap-2 sm:grid-cols-2">
                    {[
                      { mod: "unique", desc: "Adds a unique database index" },
                      { mod: "required", desc: "Marks the field as required" },
                    ].map((m) => (
                      <div
                        key={m.mod}
                        className="rounded-lg border border-border/30 bg-card/30 px-4 py-2.5"
                      >
                        <code className="text-xs font-mono text-primary/70 font-medium">
                          :{m.mod}
                        </code>
                        <span className="text-sm text-muted-foreground/60 ml-2">
                          {m.desc}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit generate seeder                                    */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit generate seeder
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-1">
                  <span className="text-xs font-mono text-primary/60 bg-primary/10 rounded-full px-2 py-0.5 mr-2">
                    alias: grit g seeder
                  </span>
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4 mt-3">
                  Generate a database seeder for one or more already-generated
                  resources. Writes{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    internal/database/&lt;name&gt;_seeder.go
                  </code>{" "}
                  with one editable example record, registers it in{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    seed.go
                  </code>
                  , and runs with{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    grit seed
                  </code>
                  . Pass <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--faker</code> to fill many rows instead.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit generate seeder Customer"
                    desc="Seed one example Customer record"
                  />
                  <TerminalCard
                    cmd="grit generate seeder Customer Order Product"
                    desc="Generate seeders for several resources at once"
                  />
                  <TerminalCard
                    cmd="grit g seeder Customer --faker --count 50"
                    desc="Fill 50 rows with gofakeit data"
                  />
                </div>
                <div className="mt-6 rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <div className="px-4 py-2.5 border-b border-border/30 bg-accent/20">
                    <span className="text-xs font-mono font-medium text-foreground/70">
                      FLAGS
                    </span>
                  </div>
                  <table className="w-full text-sm">
                    <tbody className="text-muted-foreground">
                      {[
                        { flag: "--faker", type: "bool", desc: "Fill many rows with gofakeit instead of one example" },
                        { flag: "--count", type: "int", desc: "Number of rows for the faker seeder (default 10)" },
                      ].map((f) => (
                        <tr key={f.flag} className="border-b border-border/15 last:border-0">
                          <td className="px-4 py-2 font-mono text-xs text-primary/70 w-28">
                            {f.flag}
                          </td>
                          <td className="px-4 py-2 font-mono text-xs text-muted-foreground/50 w-16">
                            {f.type}
                          </td>
                          <td className="px-4 py-2 text-muted-foreground/70 text-sm">
                            {f.desc}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit generate sequence                                  */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit generate sequence
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-1">
                  <span className="text-xs font-mono text-primary/60 bg-primary/10 rounded-full px-2 py-0.5 mr-2">
                    alias: grit g sequence
                  </span>
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4 mt-3">
                  Generate atomic sequential numbers for a resource (e.g.{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    INV-202605-0001
                  </code>
                  ). Creates a database-backed counter package plus a typed helper
                  so handlers call{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    services.Next&lt;Name&gt;Number(db, t)
                  </code>
                  .
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit generate sequence Invoice"
                    desc="Create an Invoice number sequence with defaults"
                  />
                  <TerminalCard
                    cmd="grit g sequence Invoice --prefix INV --reset monthly --width 4"
                    desc="Custom prefix, monthly reset, 4-digit width"
                  />
                  <TerminalCard
                    cmd="grit g sequence Receipt --reset never"
                    desc="A counter that never resets"
                  />
                </div>
                <div className="mt-6 rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <div className="px-4 py-2.5 border-b border-border/30 bg-accent/20">
                    <span className="text-xs font-mono font-medium text-foreground/70">
                      FLAGS
                    </span>
                  </div>
                  <table className="w-full text-sm">
                    <tbody className="text-muted-foreground">
                      {[
                        { flag: "--prefix", type: "string", desc: "Alphabetic prefix (default: first 3 chars of name, uppercased)" },
                        { flag: "--reset", type: "string", desc: "When the counter resets: monthly, yearly, never (default monthly)" },
                        { flag: "--width", type: "int", desc: "Zero-padded width of the numeric portion (default 4)" },
                      ].map((f) => (
                        <tr key={f.flag} className="border-b border-border/15 last:border-0">
                          <td className="px-4 py-2 font-mono text-xs text-primary/70 w-28">
                            {f.flag}
                          </td>
                          <td className="px-4 py-2 font-mono text-xs text-muted-foreground/50 w-16">
                            {f.type}
                          </td>
                          <td className="px-4 py-2 text-muted-foreground/70 text-sm">
                            {f.desc}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit remove resource                                    */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit remove resource
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-1">
                  <span className="text-xs font-mono text-primary/60 bg-primary/10 rounded-full px-2 py-0.5 mr-2">
                    alias: grit rm resource
                  </span>
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4 mt-3">
                  Delete all generated files for a resource and reverse all
                  marker-based injections (routes, migrations, sidebar entries).
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit remove resource Post"
                    desc="Remove the Post resource (prompts for confirmation)"
                  />
                  <TerminalCard
                    cmd="grit rm resource Post --force"
                    desc="Skip the confirmation prompt"
                  />
                </div>
                <div className="mt-6 rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <div className="px-4 py-2.5 border-b border-border/30 bg-accent/20">
                    <span className="text-xs font-mono font-medium text-foreground/70">
                      FLAGS
                    </span>
                  </div>
                  <table className="w-full text-sm">
                    <tbody className="text-muted-foreground">
                      <tr>
                        <td className="px-4 py-2 font-mono text-xs text-primary/70 w-28">
                          --force
                        </td>
                        <td className="px-4 py-2 font-mono text-xs text-muted-foreground/50 w-16">
                          bool
                        </td>
                        <td className="px-4 py-2 text-muted-foreground/70 text-sm">
                          Skip the confirmation prompt
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit add role                                           */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit add role
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Add a new role constant across the entire stack: Go models,
                  TypeScript types, Zod schemas, constants, and admin resource
                  definitions.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit add role EDITOR"
                    desc="Add EDITOR role to Go, TypeScript, and admin files"
                  />
                  <TerminalCard
                    cmd="grit add role MODERATOR"
                    desc="Add a custom MODERATOR role across the stack"
                  />
                </div>
                <div className="p-4 rounded-lg border border-primary/15 bg-primary/5 mt-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    Role names should be UPPERCASE. The command updates all the
                    right files so the role is instantly available in both Go
                    middleware and the admin panel&apos;s role dropdowns.
                  </p>
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit add web-auth                                       */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit add web-auth
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Add page-protection helpers to{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    apps/web/
                  </code>
                  : a{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    middleware.ts
                  </code>{" "}
                  SSR cookie check that redirects to{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    /login
                  </code>
                  , and a{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    ProtectedWebRoute
                  </code>{" "}
                  client wrapper. Existing files are left alone unless{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    --force
                  </code>
                  .
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit add web-auth"
                    desc="Scaffold middleware.ts + ProtectedWebRoute.tsx"
                  />
                  <TerminalCard
                    cmd="grit add web-auth --force"
                    desc="Overwrite the helpers if they already exist"
                  />
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit expose                                             */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit expose
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Scaffold a public-facing Next.js page in{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    apps/web/
                  </code>{" "}
                  that consumes an already-generated resource &mdash; reusing its
                  shared Zod schema and React Query hook instead of
                  re-implementing the form or table.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit expose form Contact --to apps/web/app/contact-us/page.tsx"
                    desc="A create-styled form page (authenticated submit)"
                  />
                  <TerminalCard
                    cmd="grit expose form Contact --to apps/web/app/contact-us/page.tsx --public-share --token 9CkLh7..."
                    desc="A no-auth public form backed by a FormShare token"
                  />
                  <TerminalCard
                    cmd="grit expose table Contact --to apps/web/app/contacts/page.tsx"
                    desc="A paginated, searchable list page"
                  />
                </div>
                <div className="mt-6 rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <div className="px-4 py-2.5 border-b border-border/30 bg-accent/20">
                    <span className="text-xs font-mono font-medium text-foreground/70">
                      FLAGS
                    </span>
                  </div>
                  <table className="w-full text-sm">
                    <tbody className="text-muted-foreground">
                      {[
                        { flag: "--to", type: "string", desc: "Destination page path (required)" },
                        { flag: "--force", type: "bool", desc: "Overwrite the destination if it exists" },
                        { flag: "--public-share", type: "bool", desc: "form only — submit via the public FormShare endpoint (no auth)" },
                        { flag: "--token", type: "string", desc: "form only — FormShare token (falls back to NEXT_PUBLIC_FORM_TOKEN)" },
                      ].map((f) => (
                        <tr key={f.flag} className="border-b border-border/15 last:border-0">
                          <td className="px-4 py-2 font-mono text-xs text-primary/70 w-40">
                            {f.flag}
                          </td>
                          <td className="px-4 py-2 font-mono text-xs text-muted-foreground/50 w-16">
                            {f.type}
                          </td>
                          <td className="px-4 py-2 text-muted-foreground/70 text-sm">
                            {f.desc}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit start                                              */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit start
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Start development servers. With no argument,{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    grit start
                  </code>{" "}
                  runs every app in the project in parallel (API + frontends, plus
                  the Wails desktop app when present) and stops them all on Ctrl+C.
                  Pass an app name to run just one.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit start"
                    desc="Start every app in the project in parallel"
                  />
                  <TerminalCard
                    cmd="grit start server"
                    desc="Start the Go API only (hot-reload via air)"
                  />
                  <TerminalCard
                    cmd="grit start client"
                    desc="Start all frontend apps via Turborepo (runs pnpm dev)"
                  />
                  <TerminalCard
                    cmd="grit start web"
                    desc="Start a single app: web, admin, expo, or desktop"
                  />
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit sync                                               */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit sync
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Parse Go model files and regenerate TypeScript types and Zod
                  schemas in{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    packages/shared
                  </code>
                  . Run this whenever you manually edit a Go model to keep
                  frontend types in sync.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit sync"
                    desc="Sync Go types → TypeScript types and Zod schemas"
                  />
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit migrate                                            */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit migrate
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Connect to the database and run GORM AutoMigrate for all
                  registered models. Use{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    --fresh
                  </code>{" "}
                  to drop all tables first for a clean slate.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit migrate"
                    desc="Run database migrations for all models"
                  />
                  <TerminalCard
                    cmd="grit migrate --fresh"
                    desc="Drop all tables, then re-run migrations from scratch"
                  />
                </div>
                <div className="mt-6 rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <div className="px-4 py-2.5 border-b border-border/30 bg-accent/20">
                    <span className="text-xs font-mono font-medium text-foreground/70">
                      FLAGS
                    </span>
                  </div>
                  <table className="w-full text-sm">
                    <tbody className="text-muted-foreground">
                      <tr>
                        <td className="px-4 py-2 font-mono text-xs text-primary/70 w-28">
                          --fresh
                        </td>
                        <td className="px-4 py-2 font-mono text-xs text-muted-foreground/50 w-16">
                          bool
                        </td>
                        <td className="px-4 py-2 text-muted-foreground/70 text-sm">
                          Drop all tables before migrating (destructive)
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div className="p-4 rounded-lg border border-red-500/20 bg-red-500/5 mt-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-red-500/90">Danger:</strong>{" "}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                      --fresh
                    </code>{" "}
                    permanently deletes all data in every table. Only use this in
                    development.
                  </p>
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit seed                                               */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit seed
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Populate the database with initial data including an admin
                  user and demo records. Perfect for bootstrapping a fresh
                  development environment.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit seed"
                    desc="Run all database seeders"
                  />
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit backup                                             */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit backup
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Dump every registered model to a ZIP archive: one CSV per table,
                  a{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    dump.sql
                  </code>{" "}
                  of INSERTs in parent-to-child order, and a{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    metadata.json
                  </code>{" "}
                  manifest. By default the archive is uploaded to object storage;
                  pass{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    --output
                  </code>{" "}
                  to write a local file instead.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit backup"
                    desc="Back up to object storage (R2 / S3 / MinIO)"
                  />
                  <TerminalCard
                    cmd="grit backup --output ./backup.zip"
                    desc="Write a local archive (no storage credentials needed)"
                  />
                </div>
                <div className="mt-6 rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <div className="px-4 py-2.5 border-b border-border/30 bg-accent/20">
                    <span className="text-xs font-mono font-medium text-foreground/70">
                      FLAGS
                    </span>
                  </div>
                  <table className="w-full text-sm">
                    <tbody className="text-muted-foreground">
                      <tr>
                        <td className="px-4 py-2 font-mono text-xs text-primary/70 w-28">
                          -o, --output
                        </td>
                        <td className="px-4 py-2 font-mono text-xs text-muted-foreground/50 w-16">
                          string
                        </td>
                        <td className="px-4 py-2 text-muted-foreground/70 text-sm">
                          Write the archive to a local file instead of uploading it
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit restore                                            */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit restore
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Run migrations, then replay a backup archive&apos;s{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    dump.sql
                  </code>{" "}
                  inside a single transaction &mdash; every row lands or none does.
                  Point it at an empty database; the archive carries data, not
                  schema.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit restore backup.zip"
                    desc="Migrate, then replay the archive in one transaction"
                  />
                  <TerminalCard
                    cmd="grit restore backup.zip --no-migrate"
                    desc="Skip migrations (schema already exists)"
                  />
                </div>
                <div className="mt-6 rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <div className="px-4 py-2.5 border-b border-border/30 bg-accent/20">
                    <span className="text-xs font-mono font-medium text-foreground/70">
                      FLAGS
                    </span>
                  </div>
                  <table className="w-full text-sm">
                    <tbody className="text-muted-foreground">
                      <tr>
                        <td className="px-4 py-2 font-mono text-xs text-primary/70 w-28">
                          --no-migrate
                        </td>
                        <td className="px-4 py-2 font-mono text-xs text-muted-foreground/50 w-16">
                          bool
                        </td>
                        <td className="px-4 py-2 text-muted-foreground/70 text-sm">
                          Skip running migrations before restoring
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit studio                                             */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit studio
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Open the GORM Studio database browser at{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    http://localhost:8080/studio
                  </code>
                  . For web projects, make sure your API server is running first.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit studio"
                    desc="Open GORM Studio in your browser"
                  />
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit routes                                             */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit routes
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Parse{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    routes.go
                  </code>{" "}
                  and print a table of every registered HTTP route with its
                  method, path, handler, and middleware group.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit routes"
                    desc="List all registered API routes"
                  />
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit upgrade                                            */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit upgrade
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Upgrade an existing project to the latest scaffold templates.
                  This regenerates framework components (admin panel, web app,
                  configs) while preserving your resource definitions and API
                  code.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit upgrade"
                    desc="Upgrade project templates (prompts before overwriting)"
                  />
                  <TerminalCard
                    cmd="grit upgrade --force"
                    desc="Overwrite all files without prompting"
                  />
                </div>
                <div className="mt-6 rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <div className="px-4 py-2.5 border-b border-border/30 bg-accent/20">
                    <span className="text-xs font-mono font-medium text-foreground/70">
                      FLAGS
                    </span>
                  </div>
                  <table className="w-full text-sm">
                    <tbody className="text-muted-foreground">
                      <tr>
                        <td className="px-4 py-2 font-mono text-xs text-primary/70 w-28">
                          -f, --force
                        </td>
                        <td className="px-4 py-2 font-mono text-xs text-muted-foreground/50 w-16">
                          bool
                        </td>
                        <td className="px-4 py-2 text-muted-foreground/70 text-sm">
                          Overwrite all files without prompting
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit update                                             */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit update
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Update the Grit CLI binary to the latest version. It checks
                  GitHub first and exits if you are already up to date. Otherwise,
                  if the Go toolchain is on your PATH it runs{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    go install ...@&lt;latest&gt;
                  </code>
                  ; if Go is not installed it downloads the matching prebuilt
                  binary from the GitHub release and atomically swaps it in (on
                  Windows the old{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    .exe
                  </code>{" "}
                  is renamed to{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    .old
                  </code>{" "}
                  first). Aliased as{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    grit self-update
                  </code>
                  .
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit update"
                    desc="Update the Grit CLI to the latest release"
                  />
                  <TerminalCard
                    cmd="grit update --from-release"
                    desc="Skip go install and pull the prebuilt GitHub binary"
                  />
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit version                                            */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit version
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Print the installed Grit CLI version.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit version"
                    desc="Print the current CLI version number"
                  />
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit down / grit up                                     */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit down / grit up
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Toggle maintenance mode.{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    grit down
                  </code>{" "}
                  writes a{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    .maintenance
                  </code>{" "}
                  file that makes the scaffolded middleware return 503 for every
                  request;{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    grit up
                  </code>{" "}
                  removes it and resumes normal handling.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit down"
                    desc="Enter maintenance mode (all requests get 503)"
                  />
                  <TerminalCard
                    cmd="grit up"
                    desc="Exit maintenance mode"
                  />
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  grit deploy                                             */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  grit deploy
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Build the app, upload it over SSH, configure a systemd service,
                  and optionally set up a Caddy reverse proxy with auto-TLS. Flags
                  fall back to{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    DEPLOY_*
                  </code>{" "}
                  environment variables.
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit deploy --host user@server.com --domain myapp.com"
                    desc="Deploy and serve behind Caddy with auto-TLS"
                  />
                  <TerminalCard
                    cmd="grit deploy"
                    desc="Use DEPLOY_HOST / DEPLOY_DOMAIN / DEPLOY_KEY_FILE from .env"
                  />
                </div>
                <div className="mt-6 rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <div className="px-4 py-2.5 border-b border-border/30 bg-accent/20">
                    <span className="text-xs font-mono font-medium text-foreground/70">
                      FLAGS
                    </span>
                  </div>
                  <table className="w-full text-sm">
                    <tbody className="text-muted-foreground">
                      {[
                        { flag: "--host", type: "string", desc: "SSH host (user@server.com) or DEPLOY_HOST" },
                        { flag: "--port", type: "string", desc: "SSH port (default 22)" },
                        { flag: "--key", type: "string", desc: "Path to SSH private key or DEPLOY_KEY_FILE" },
                        { flag: "--domain", type: "string", desc: "Domain for the Caddy reverse proxy or DEPLOY_DOMAIN" },
                        { flag: "--app-port", type: "string", desc: "Port the app listens on (default 8080)" },
                      ].map((f) => (
                        <tr key={f.flag} className="border-b border-border/15 last:border-0">
                          <td className="px-4 py-2 font-mono text-xs text-primary/70 w-28">
                            {f.flag}
                          </td>
                          <td className="px-4 py-2 font-mono text-xs text-muted-foreground/50 w-16">
                            {f.type}
                          </td>
                          <td className="px-4 py-2 text-muted-foreground/70 text-sm">
                            {f.desc}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  Desktop apps                                            */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                  Desktop apps
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Scaffold and ship standalone Wails desktop applications (Go +
                  React + SQLite).
                </p>
                <div className="space-y-3">
                  <TerminalCard
                    cmd="grit new-desktop myapp"
                    desc="Scaffold a standalone Wails desktop app"
                  />
                  <TerminalCard
                    cmd="grit compile"
                    desc="Build the desktop app into a binary (wails build)"
                  />
                  <TerminalCard
                    cmd="grit package"
                    desc="Build a distributable installer (.exe / .app / binary)"
                  />
                  <TerminalCard
                    cmd="grit package --platform windows/amd64 --clean"
                    desc="Target a platform and clean the build dir first"
                  />
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  Common Workflows                                        */}
              {/* -------------------------------------------------------- */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Common Workflows
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Copy-paste recipes for everyday development tasks.
                </p>

                {/* Workflow: New project from scratch */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/90 mb-3">
                    New project from scratch
                  </h3>
                  <CodeBlock
                    terminal
                    code={`# Create the project
grit new myapp

# Start infrastructure
cd myapp && docker compose up -d

# Start the API
grit start server

# In another terminal — start frontend
grit start client`}
                    className="mb-0"
                  />
                </div>

                {/* Workflow: Add a new resource */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/90 mb-3">
                    Add a new resource
                  </h3>
                  <CodeBlock
                    terminal
                    code={`# Generate the resource
grit g resource Product --fields "name:string,price:float,description:text,published:bool"

# Sync types to keep everything in sync
grit sync

# Run migrations to create the table
grit migrate`}
                    className="mb-0"
                  />
                </div>

                {/* Workflow: Add a resource with relationships */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/90 mb-3">
                    Resource with relationships
                  </h3>
                  <CodeBlock
                    terminal
                    code={`# Generate a Category resource first
grit g resource Category --fields "name:string,description:text"

# Generate Product that belongs to Category
grit g resource Product --fields "name:string,price:float,category:belongs_to"

# Apply database changes
grit migrate`}
                    className="mb-0"
                  />
                </div>

                {/* Workflow: Fresh database reset */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/90 mb-3">
                    Fresh database reset
                  </h3>
                  <CodeBlock
                    terminal
                    code={`# Drop all tables and re-migrate
grit migrate --fresh

# Re-seed with initial data
grit seed`}
                    className="mb-0"
                  />
                </div>

                {/* Workflow: Remove and regenerate */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/90 mb-3">
                    Remove and regenerate a resource
                  </h3>
                  <CodeBlock
                    terminal
                    code={`# Remove the old resource
grit rm resource Post --force

# Regenerate with updated fields
grit g resource Post --fields "title:string,slug:slug,content:richtext,published:bool,views:int"`}
                    className="mb-0"
                  />
                </div>

                {/* Workflow: Upgrade existing project */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/90 mb-3">
                    Upgrade an existing project
                  </h3>
                  <CodeBlock
                    terminal
                    code={`# Update the CLI first
grit update

# Then upgrade project templates
grit upgrade`}
                    className="mb-0"
                  />
                </div>
              </div>

              {/* -------------------------------------------------------- */}
              {/*  Command Tree                                            */}
              {/* -------------------------------------------------------- */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Full Command Tree
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The complete hierarchy of every command and subcommand.
                </p>
                <CodeBlock
                  code={`grit
├── init                      # Write CLAUDE.md / AGENTS.md convention docs
│   └── --force               # Overwrite existing files
│
├── new <project-name|.>      # Scaffold a new project
│   ├── --single              # Go API + embedded SPA (one binary)
│   ├── --double              # API + web (no admin)
│   ├── --triple              # API + web + admin (default)
│   ├── --api                 # Go API only
│   ├── --mobile              # API + Expo mobile only
│   ├── --expo                # Add Expo mobile (additive)
│   ├── --desktop             # Add Wails desktop client (combinable, not --single)
│   ├── --full                # Everything + docs site
│   ├── --next / --vite       # Frontend: Next.js or TanStack (Vite)
│   ├── --style <variant>     # Admin style: default, modern, minimal, glass, centered
│   ├── --theme <name>        # Theme: atlas, aurora, pulse
│   ├── --here                # Scaffold into the current directory
│   └── --force               # Allow a non-empty directory
│
├── new-desktop <name>        # Scaffold a standalone Wails desktop app
│
├── generate (g)              # Code generation
│   ├── resource <Name>       # Generate full-stack CRUD resource
│   │   ├── --fields "..."    # Inline field definitions
│   │   ├── --from file.yaml  # YAML field definitions
│   │   ├── -i, --interactive # Interactive field prompts
│   │   ├── --roles "..."     # Restrict routes to roles
│   │   ├── --seed            # Also generate a seeder (one example row)
│   │   ├── --faker           # Also generate a gofakeit seeder (implies --seed)
│   │   └── --count <n>       # Rows for the faker seeder (default 10)
│   ├── seeder <Resource>...  # Seeder for existing resources (--faker, --count)
│   └── sequence <Name>       # Sequential numbering helper (--prefix, --reset, --width)
│
├── remove (rm)               # Remove components
│   └── resource <Name>       # Remove a generated resource
│       └── --force           # Skip confirmation
│
├── add                       # Add components
│   ├── role <ROLE_NAME>      # Add a role across the stack
│   └── web-auth              # Add page-protection helpers to apps/web/ (--force)
│
├── expose                    # Scaffold a web page for a resource
│   ├── form <Resource>       # Public form page (--to, --public-share, --token, --force)
│   └── table <Resource>      # Paginated list page (--to, --force)
│
├── start                     # Development servers (all apps if no arg)
│   ├── server                # Start the Go API (air hot-reload)
│   ├── client                # Start frontend apps (Turborepo)
│   ├── web / admin           # Start a single Next.js app
│   ├── expo                  # Start the Expo mobile app
│   └── desktop               # Start the Wails desktop app
│
├── compile                   # Build desktop app executable (Wails)
├── package                   # Build a distributable desktop installer
│   ├── --platform <os/arch>  # Target platform (default: host)
│   ├── --no-installer        # Raw binary only (skip NSIS on Windows)
│   └── --clean               # Clean the build dir first
│
├── studio                    # Open GORM Studio database browser
├── sync                      # Sync Go types → TypeScript
├── migrate                   # Run database migrations
│   └── --fresh               # Drop all tables first
├── seed                      # Run database seeders
├── backup                    # Back up the entire database
│   └── -o, --output          # Write a local archive instead of uploading
├── restore <backup.zip>      # Restore the database from an archive
│   └── --no-migrate          # Skip migrations before restoring
├── routes                    # List all registered API routes
├── down                      # Enter maintenance mode (503)
├── up                        # Exit maintenance mode
├── deploy                    # Deploy to a remote server
│   ├── --host / --port       # SSH host and port
│   ├── --key                 # SSH private key
│   ├── --domain              # Domain for Caddy + auto-TLS
│   └── --app-port            # Port the app listens on
├── upgrade                   # Upgrade project templates
│   └── -f, --force           # Overwrite without prompting
├── update (self-update)      # Update the CLI binary
│   └── --from-release        # Pull the prebuilt GitHub binary
└── version                   # Print CLI version`}
                  className="mb-0"
                />
              </div>

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button
                  variant="ghost"
                  size="sm"
                  asChild
                  className="text-muted-foreground/60 hover:text-foreground"
                >
                  <Link
                    href="/docs/getting-started/troubleshooting"
                    className="gap-1.5"
                  >
                    <ArrowLeft className="h-3.5 w-3.5" />
                    Troubleshooting
                  </Link>
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  asChild
                  className="text-muted-foreground/60 hover:text-foreground"
                >
                  <Link
                    href="/docs/concepts/cli"
                    className="gap-1.5"
                  >
                    CLI Commands (Deep Dive)
                    <ArrowRight className="h-3.5 w-3.5" />
                  </Link>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
