import Link from 'next/link'
import { ArrowLeft, ArrowRight, Download } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/ai-skill/llm-guide')

// ── reusable mini-components ────────────────────────────────────
function SectionNum({ n }: { n: number }) {
  return (
    <span className="inline-flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/20 text-xs font-mono font-semibold text-primary mr-2">
      {n}
    </span>
  )
}

function Warn({ children }: { children: React.ReactNode }) {
  return (
    <div className="p-4 rounded-lg border border-red-500/20 bg-red-500/5 mb-4">
      <p className="text-sm text-red-400/80 leading-relaxed font-medium">{children}</p>
    </div>
  )
}

function Tip({ children }: { children: React.ReactNode }) {
  return (
    <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mb-4">
      <p className="text-sm text-muted-foreground leading-relaxed">{children}</p>
    </div>
  )
}

function Note({ children }: { children: React.ReactNode }) {
  return (
    <div className="p-4 rounded-lg border border-amber-500/20 bg-amber-500/5 mb-4">
      <p className="text-sm text-amber-400/80 leading-relaxed">{children}</p>
    </div>
  )
}

export default function LLMGuidePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">

            {/* ── Header ── */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">For AI Assistants</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Grit Framework — Complete LLM Reference
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                The canonical guide for AI assistants building with Grit. Covers every CLI command,
                all field types with real examples, resource generation patterns, form builder,
                datatable builder, standalone usage, relationships, slugs, media fields, and the
                rules that must never be broken.
              </p>
              <Tip>
                <strong className="text-primary/90">For AI tools:</strong> Read this entire page before generating any code. Every convention here is enforced by the CLI — deviating breaks{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate</code> and{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit sync</code>.
              </Tip>
              <div className="flex items-center justify-between p-4 rounded-xl border border-border/40 bg-card/50 mb-6">
                <div>
                  <p className="text-sm font-semibold text-foreground/80">Prefer a PDF?</p>
                  <p className="text-xs text-muted-foreground/60">Download the Grit Handbook for offline reading and printing</p>
                </div>
                <a
                  href="https://14j7oh8kso.ufs.sh/f/HLxTbDBCDLwfeHHJl34ZKSqNhOvVj6p9rg3Icmo05TAEwQ4a"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  <Button size="sm" variant="outline" className="shrink-0 gap-1.5 ml-4">
                    <Download className="h-3.5 w-3.5" />
                    Download PDF
                  </Button>
                </a>
              </div>
            </div>

            <div className="prose-grit">

              {/* ════════════════════════════════════════════════════
                  1. What is Grit
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={1} />What is Grit?
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit is a <strong className="text-foreground/80">full-stack meta-framework</strong> that fuses a Go backend with Next.js frontends
                  inside a Turborepo monorepo. It is opinionated, batteries-included, and optimised for
                  AI-assisted development. A single CLI command scaffolds a complete production-ready project.
                </p>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  {[
                    { label: 'Backend', value: 'Go 1.24 · Gin · GORM · PostgreSQL · Redis' },
                    { label: 'Frontend', value: 'Next.js 15 · React 19 · TypeScript · Tailwind CSS' },
                    { label: 'Admin panel', value: 'Custom Filament-like panel scaffolded with Next.js' },
                    { label: 'Shared types', value: 'packages/shared — Zod schemas + TypeScript types' },
                    { label: 'Monorepo', value: 'Turborepo + pnpm workspaces' },
                    { label: 'Infra', value: 'Docker Compose · MinIO · Mailhog (dev) · S3/R2/Resend (prod)' },
                  ].map((item) => (
                    <div key={item.label} className="p-3 rounded-lg border border-border/30 bg-card/30">
                      <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-1">{item.label}</p>
                      <p className="text-xs font-mono text-foreground/70">{item.value}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  2. Project Structure
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={2} />Project Structure
                </h2>
                <CodeBlock language="bash" filename="Grit monorepo layout" code={`myapp/
├── apps/
│   ├── api/                          # Go backend
│   │   └── internal/
│   │       ├── config/               # DB, Redis, env, storage
│   │       ├── middleware/           # Gzip, CORS, RequestID, Logger, Auth, Cache
│   │       ├── models/               # GORM structs + model registry
│   │       │   └── models.go         # // GRIT:MODELS marker — NEVER remove
│   │       ├── handlers/             # Thin HTTP handlers (call services)
│   │       ├── services/             # Business logic + GORM queries
│   │       ├── routes/               # Route registration
│   │       │   └── routes.go         # // GRIT:ROUTES marker — NEVER remove
│   │       ├── jobs/                 # asynq background jobs
│   │       └── seeders/              # DB seed data
│   ├── web/                          # Next.js public website
│   │   ├── app/                      # App Router pages
│   │   ├── components/
│   │   └── lib/api-client.ts         # Axios client + uploadFile()
│   └── admin/                        # Next.js admin panel
│       ├── app/(dashboard)/          # Protected admin pages
│       │   └── [resource]/
│       │       ├── page.tsx          # ResourcePage component
│       │       └── _resource.ts      # defineResource() config
│       ├── components/               # DataTable, FormBuilder, FormStepper
│       ├── hooks/                    # React Query hooks (generated)
│       └── lib/api-client.ts         # Axios client + uploadFile()
├── packages/
│   └── shared/
│       └── src/
│           ├── schemas/              # Zod schemas (generated by grit sync)
│           ├── types/                # TypeScript interfaces
│           └── index.ts              # Re-exports everything
├── docker-compose.yml                # PostgreSQL, Redis, MinIO, Mailhog
├── docker-compose.prod.yml           # Production config
├── turbo.json
├── pnpm-workspace.yaml
├── grit.config.ts                    # Project config
├── GRIT_SKILL.md                     # AI context file (auto-generated)
└── .env                              # Environment variables`} />
              </div>

              {/* ════════════════════════════════════════════════════
                  3. CLI Commands
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={3} />All CLI Commands
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Install once with the one-line script (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">curl -fsSL https://gritframework.dev/install.sh | sh</code>{' '}
                  on macOS / Linux, <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">iwr -useb https://gritframework.dev/install.ps1 | iex</code>{' '}
                  on Windows) — or with Go directly via{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">go install github.com/MUKE-coder/grit/v3/cmd/grit@latest</code>.
                  Every command is idempotent — safe to re-run.
                </p>

                {/* Project lifecycle */}
                <h3 className="text-base font-semibold text-foreground/80 mb-3 mt-6">Project Lifecycle</h3>
                <div className="space-y-3 mb-8">
                  {[
                    {
                      cmd: 'grit new <app-name>',
                      desc: 'Scaffold a complete new project. Interactive by default — prompts for architecture mode and frontend. 5 architecture modes: --single (one Go binary with embedded SPA), --double (API + web), --triple (API + web + admin), --api (headless API), --mobile (API + Expo). 2 frontends: --next (Next.js), --vite (TanStack Router). Plus a --desktop flag that adds a Wails desktop client to any monorepo architecture (combinable with --triple, --double, --mobile, --api; NOT --single). Supports grit new . to scaffold into current directory, --here for explicit in-place scaffolding, and --force to scaffold into a non-empty directory. Other flags: --full (default, all apps), --style default|modern|minimal|glass|centered. The --desktop flag is distinct from grit new-desktop (standalone offline-first app with embedded Go + SQLite); --desktop is for always-online multi-client apps that share the monorepo API.',
                      example: 'grit new myapp\ngrit new myapp --triple --vite\ngrit new . --triple --vite\ngrit new myapp --triple --next --desktop\ngrit new myapp --mobile --desktop',
                    },
                    {
                      cmd: 'grit start server',
                      desc: 'Start ONLY the Go API server (no hot-reload). Runs go run from apps/api. Use this when you only need the backend, or for production-like API testing.',
                      example: 'grit start server',
                    },
                    {
                      cmd: 'grit start client',
                      desc: 'Start ONLY the frontend apps (web + admin) via pnpm dev / Turborepo. Does not start the Go API. Use this when the backend is already running separately.',
                      example: 'grit start client',
                    },
                  ].map((item) => (
                    <div key={item.cmd} className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                      <div className="px-4 py-2.5 border-b border-border/20 bg-accent/20">
                        <code className="text-xs font-mono font-semibold text-primary/80">{item.cmd}</code>
                      </div>
                      <div className="p-4">
                        <p className="text-xs text-muted-foreground/70 leading-relaxed mb-3">{item.desc}</p>
                        <div className="rounded-lg border border-border/20 bg-background/60 px-3 py-2">
                          <span className="text-primary/40 font-mono text-xs select-none">$ </span>
                          <span className="font-mono text-xs text-foreground/60">{item.example}</span>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>

                <Warn>
                  There is no <code className="font-mono bg-red-500/10 px-1 rounded">grit dev</code> command — it does not exist. Run the API and frontend in separate terminals:{' '}
                  <code className="font-mono bg-red-500/10 px-1 rounded">grit start server</code> in Terminal 1,
                  then <code className="font-mono bg-red-500/10 px-1 rounded">grit start client</code> in Terminal 2.
                  Alternatively, <code className="font-mono bg-red-500/10 px-1 rounded">pnpm dev</code> from the project root starts all apps via Turborepo.
                </Warn>

                {/* Code generation */}
                <h3 className="text-base font-semibold text-foreground/80 mb-3">Code Generation</h3>
                <div className="space-y-3 mb-8">
                  {[
                    {
                      cmd: 'grit generate resource <Name> --fields "field:type,..."',
                      desc: 'Generate a full-stack resource: Go model + GORM migration, handler, service, routes, Zod schema, TypeScript types, React Query hook, and admin page — all wired together. Fields are comma-separated name:type pairs. Special types: slug:slug (auto-generated from title), image/images/video/videos/file/files (presigned upload), enum:A,B,C (select), uint:fk:Model (belongs_to), []uint:m2m:Model (many_to_many). Extra flags: --from schema.yaml (fields from YAML), -i / --interactive (prompt per field), --roles "ADMIN,EDITOR" (restrict to roles).',
                      example: 'grit generate resource Product --fields "name:string,slug:slug,price:float64,image:image"',
                    },
                    {
                      cmd: 'grit sync',
                      desc: 'Regenerate packages/shared/src/schemas and packages/shared/src/types from all Go models. Run this after manually editing a Go model struct.',
                      example: 'grit sync',
                    },
                    {
                      cmd: 'grit remove resource <Name>',
                      desc: 'Remove a previously generated resource: deletes Go files, removes the model from the registry, cleans up routes, and removes the admin page.',
                      example: 'grit remove resource Product',
                    },
                  ].map((item) => (
                    <div key={item.cmd} className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                      <div className="px-4 py-2.5 border-b border-border/20 bg-accent/20">
                        <code className="text-xs font-mono font-semibold text-primary/80">{item.cmd}</code>
                      </div>
                      <div className="p-4">
                        <p className="text-xs text-muted-foreground/70 leading-relaxed mb-3">{item.desc}</p>
                        <div className="rounded-lg border border-border/20 bg-background/60 px-3 py-2">
                          <span className="text-primary/40 font-mono text-xs select-none">$ </span>
                          <span className="font-mono text-xs text-foreground/60">{item.example}</span>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Database & roles */}
                <h3 className="text-base font-semibold text-foreground/80 mb-3">Database & Roles</h3>
                <div className="space-y-3 mb-8">
                  {[
                    {
                      cmd: 'grit migrate',
                      desc: 'Run GORM AutoMigrate for all models in RegisteredModels. Safe to run repeatedly — adds new columns/tables without dropping existing data. Use --fresh to drop all tables first and start clean (WARNING: destroys all data — development only).',
                      example: 'grit migrate          # safe incremental\ngrit migrate --fresh  # drop + recreate all tables',
                    },
                    {
                      cmd: 'grit seed',
                      desc: 'Run the database seeders in apps/api/internal/seeders/. Creates the default admin user and sample data.',
                      example: 'grit seed',
                    },
                    {
                      cmd: 'grit add role <ROLE_NAME>',
                      desc: 'Add a new RBAC role in 7 places simultaneously: Go constants, Go middleware, Zod schema, TypeScript type, admin dropdown, seed file, and GRIT_SKILL.md.',
                      example: 'grit add role MODERATOR',
                    },
                  ].map((item) => (
                    <div key={item.cmd} className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                      <div className="px-4 py-2.5 border-b border-border/20 bg-accent/20">
                        <code className="text-xs font-mono font-semibold text-primary/80">{item.cmd}</code>
                      </div>
                      <div className="p-4">
                        <p className="text-xs text-muted-foreground/70 leading-relaxed mb-3">{item.desc}</p>
                        <div className="rounded-lg border border-border/20 bg-background/60 px-3 py-2">
                          <span className="text-primary/40 font-mono text-xs select-none">$ </span>
                          <span className="font-mono text-xs text-foreground/60">{item.example}</span>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Utilities */}
                <h3 className="text-base font-semibold text-foreground/80 mb-3">Utilities</h3>
                <div className="space-y-3">
                  {[
                    {
                      cmd: 'grit upgrade',
                      desc: 'Regenerates framework scaffold files (admin panel, web app, shared package configs) in the current project while preserving your custom resource code. Run this after a new Grit release to pull in updated templates. Use --force to skip confirmation prompts. This upgrades PROJECT FILES — not the CLI binary.',
                      example: 'grit upgrade\ngrit upgrade --force',
                    },
                    {
                      cmd: 'grit update',
                      desc: 'Removes the current Grit CLI binary from $GOPATH/bin and reinstalls the latest version from GitHub via go install. Use this to update the grit TOOL ITSELF — not the project files (use grit upgrade for that).',
                      example: 'grit update',
                    },
                    {
                      cmd: 'grit routes',
                      desc: 'List all registered API routes in the project. Useful for debugging and verifying that generated routes are wired correctly.',
                      example: 'grit routes',
                    },
                    {
                      cmd: 'grit down',
                      desc: 'Enable maintenance mode. Returns 503 Service Unavailable for all requests. Use this during deployments or database migrations.',
                      example: 'grit down',
                    },
                    {
                      cmd: 'grit up',
                      desc: 'Disable maintenance mode. Resumes normal request handling after grit down.',
                      example: 'grit up',
                    },
                    {
                      cmd: 'grit deploy',
                      desc: 'Deploy the project to a production server. Handles build, push, and server configuration.',
                      example: 'grit deploy',
                    },
                    {
                      cmd: 'grit version',
                      desc: 'Print the installed CLI version.',
                      example: 'grit version',
                    },
                  ].map((item) => (
                    <div key={item.cmd} className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                      <div className="px-4 py-2.5 border-b border-border/20 bg-accent/20">
                        <code className="text-xs font-mono font-semibold text-primary/80">{item.cmd}</code>
                      </div>
                      <div className="p-4">
                        <p className="text-xs text-muted-foreground/70 leading-relaxed mb-3">{item.desc}</p>
                        <div className="rounded-lg border border-border/20 bg-background/60 px-3 py-2">
                          <span className="text-primary/40 font-mono text-xs select-none">$ </span>
                          <span className="font-mono text-xs text-foreground/60">{item.example}</span>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  4. Field Types
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={4} />Field Types for Code Generation
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Fields follow the format{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">name:type</code> or{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">name:type:fk:RelatedModel</code> for relationships.
                </p>

                {/* Type table */}
                <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden mb-8">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">Type</th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">Go Type</th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">Admin Form Field</th>
                        <th className="px-4 py-2.5 text-left text-xs font-semibold text-muted-foreground uppercase">Notes</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20 text-xs font-mono">
                      {[
                        { type: 'string', go: 'string', admin: 'text input', note: 'Short text, titles, names' },
                        { type: 'slug', go: 'string (uniqueIndex)', admin: 'text input (auto-generated)', note: 'URL-friendly slug auto-generated from title on save' },
                        { type: 'text', go: 'string', admin: 'textarea', note: 'Long text, descriptions' },
                        { type: 'richtext', go: 'string', admin: 'rich text editor', note: 'HTML content (TipTap)' },
                        { type: 'int', go: 'int', admin: 'number input', note: 'Signed integer' },
                        { type: 'uint', go: 'uint', admin: 'number input', note: 'Unsigned integer, IDs' },
                        { type: 'float64', go: 'float64', admin: 'number input (decimal)', note: 'Prices, coordinates' },
                        { type: 'bool', go: 'bool', admin: 'toggle / checkbox', note: 'Active, published, featured' },
                        { type: 'time', go: 'time.Time', admin: 'date picker', note: 'Timestamps, due dates' },
                        { type: 'image', go: 'string', admin: 'image upload (dropzone)', note: 'Single image URL (presigned upload)' },
                        { type: 'images', go: 'pq.StringArray', admin: 'multi-image dropzone', note: 'Array of image URLs' },
                        { type: 'video', go: 'string', admin: 'video upload (dropzone)', note: 'Single video URL' },
                        { type: 'videos', go: 'pq.StringArray', admin: 'multi-video dropzone', note: 'Array of video URLs' },
                        { type: 'file', go: 'string', admin: 'file upload (dropzone)', note: 'Single file URL (PDF, doc…)' },
                        { type: 'files', go: 'pq.StringArray', admin: 'multi-file dropzone', note: 'Array of file URLs' },
                        { type: 'enum:A,B,C', go: 'string', admin: 'select dropdown', note: 'Fixed set of values' },
                        { type: 'uint:fk:Model', go: 'uint + Model field', admin: 'relationship select', note: 'belongs_to (one-to-many from child side)' },
                        { type: '[]uint:m2m:Model', go: '[]Model (GORM M2M)', admin: 'multi-relationship select', note: 'many_to_many join table' },
                      ].map((row) => (
                        <tr key={row.type}>
                          <td className="px-4 py-2 text-primary/80">{row.type}</td>
                          <td className="px-4 py-2 text-foreground/60">{row.go}</td>
                          <td className="px-4 py-2 text-muted-foreground/60">{row.admin}</td>
                          <td className="px-4 py-2 text-muted-foreground/40 font-sans">{row.note}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                {/* Slug explained */}
                <div className="mb-8">
                  <h3 className="text-lg font-semibold tracking-tight mb-3">Understanding Slugs</h3>
                  <p className="text-muted-foreground leading-relaxed mb-3">
                    A <strong className="text-foreground/80">slug</strong> is a URL-friendly string derived from a title.
                    Instead of <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/blog/42</code>, you get{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/blog/my-first-post</code>.
                    In Grit, use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">slug:slug</code> (not <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">slug:string</code>) as the field type.
                    The Go service auto-generates the slug from the title using a slugify helper before saving, and GORM adds a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">uniqueIndex</code> automatically.
                  </p>
                  <CodeBlock language="go" filename="apps/api/internal/models/post.go (generated)" code={`type Post struct {
    gorm.Model
    Title     string \`gorm:"not null" json:"title"\`
    Slug      string \`gorm:"uniqueIndex;not null" json:"slug"\`
    Content   string \`gorm:"type:text" json:"content"\`
    Published bool   \`gorm:"default:false" json:"published"\`
}

// Service auto-generates slug before saving:
// slug = strings.ToLower(strings.ReplaceAll(title, " ", "-"))
// The public GetBySlug handler route: GET /api/posts/slug/:slug`} />
                  <p className="text-sm text-muted-foreground/60 mt-2">
                    Always use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">slug:slug</code> (not <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">slug:string</code>) — the dedicated type adds a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">uniqueIndex</code> and wires up auto-generation from the title automatically.
                  </p>
                </div>

                {/* Relationships explained */}
                <div className="mb-4">
                  <h3 className="text-lg font-semibold tracking-tight mb-3">Understanding Relationships</h3>

                  {/* Syntax summary */}
                  <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden mb-5">
                    <table className="w-full text-xs">
                      <thead>
                        <tr className="border-b border-border/30 bg-accent/20">
                          <th className="px-4 py-2 text-left font-semibold text-muted-foreground uppercase">Type</th>
                          <th className="px-4 py-2 text-left font-semibold text-muted-foreground uppercase">Syntax</th>
                          <th className="px-4 py-2 text-left font-semibold text-muted-foreground uppercase">Admin Field</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-border/20 font-mono">
                        <tr>
                          <td className="px-4 py-2 text-foreground/70 font-sans">One-to-Many (belongs_to)</td>
                          <td className="px-4 py-2 text-primary/80">fieldname_id:uint:fk:ModelName</td>
                          <td className="px-4 py-2 text-muted-foreground/60 font-sans">Searchable single select</td>
                        </tr>
                        <tr>
                          <td className="px-4 py-2 text-foreground/70 font-sans">One-to-One (has_one)</td>
                          <td className="px-4 py-2 text-primary/80">fieldname_id:uint:fk:ModelName</td>
                          <td className="px-4 py-2 text-muted-foreground/60 font-sans">Searchable single select</td>
                        </tr>
                        <tr>
                          <td className="px-4 py-2 text-foreground/70 font-sans">Many-to-Many</td>
                          <td className="px-4 py-2 text-primary/80">fieldname_ids:[]uint:m2m:ModelName</td>
                          <td className="px-4 py-2 text-muted-foreground/60 font-sans">Multi-select tag input</td>
                        </tr>
                      </tbody>
                    </table>
                  </div>

                  <div className="space-y-6">
                    {/* One-to-Many */}
                    <div className="rounded-xl border border-border/30 bg-card/30 overflow-hidden">
                      <div className="px-4 py-2.5 border-b border-border/20 bg-accent/20">
                        <h4 className="text-sm font-semibold">One-to-Many (belongs_to)</h4>
                      </div>
                      <div className="p-4 space-y-3">
                        <p className="text-xs text-muted-foreground/70 leading-relaxed">
                          A child record belongs to one parent. Syntax:{' '}
                          <code className="font-mono bg-accent/50 px-1 rounded">fieldname_id:uint:fk:ModelName</code>.
                          Use the parent model name in PascalCase after <code className="font-mono bg-accent/50 px-1 rounded">fk:</code>.
                          This adds a <code className="font-mono bg-accent/50 px-1 rounded">CategoryID uint</code> FK column and a{' '}
                          <code className="font-mono bg-accent/50 px-1 rounded">Category Category</code> preload field to the struct.
                          The admin form renders a <strong>searchable single-select dropdown</strong> populated from <code className="font-mono bg-accent/50 px-1 rounded">/api/categories</code>.
                        </p>
                        <CodeBlock language="bash" filename="Example — Post belongs to Category" code={`# 1. Generate the parent first
grit generate resource Category --fields "name:string,slug:slug,description:text"

# 2. Generate the child referencing it
grit generate resource Post \\
  --fields "title:string,slug:slug,content:richtext,category_id:uint:fk:Category,is_published:bool"

grit migrate`} />
                      </div>
                    </div>

                    {/* One-to-One */}
                    <div className="rounded-xl border border-border/30 bg-card/30 overflow-hidden">
                      <div className="px-4 py-2.5 border-b border-border/20 bg-accent/20">
                        <h4 className="text-sm font-semibold">One-to-One (has_one)</h4>
                      </div>
                      <div className="p-4 space-y-3">
                        <p className="text-xs text-muted-foreground/70 leading-relaxed">
                          Exactly one child per parent. Uses the <strong>same FK syntax</strong> as belongs_to —{' '}
                          <code className="font-mono bg-accent/50 px-1 rounded">fieldname_id:uint:fk:ModelName</code> — but the CLI automatically adds a{' '}
                          <code className="font-mono bg-accent/50 px-1 rounded">uniqueIndex</code> on the FK column.
                          GORM treats it as one-to-one when the FK has a unique constraint.
                          The admin form also renders a single-select dropdown (same as belongs_to).
                        </p>
                        <CodeBlock language="bash" filename="Example — Profile has one User" code={`# 1. User model already exists (built in)

# 2. Generate Profile with a unique FK to User
grit generate resource Profile \\
  --fields "bio:text,avatar:image,website:string,user_id:uint:fk:User"

# The CLI adds uniqueIndex on user_id automatically for one-to-one FK fields

grit migrate`} />
                      </div>
                    </div>

                    {/* Many-to-Many */}
                    <div className="rounded-xl border border-border/30 bg-card/30 overflow-hidden">
                      <div className="px-4 py-2.5 border-b border-border/20 bg-accent/20">
                        <h4 className="text-sm font-semibold">Many-to-Many</h4>
                      </div>
                      <div className="p-4 space-y-3">
                        <p className="text-xs text-muted-foreground/70 leading-relaxed">
                          A record can belong to many of another, and vice versa. Syntax:{' '}
                          <code className="font-mono bg-accent/50 px-1 rounded">fieldname_ids:[]uint:m2m:ModelName</code>.
                          Key differences from belongs_to: the field name is <strong>plural</strong> (e.g. <code className="font-mono bg-accent/50 px-1 rounded">tag_ids</code>),
                          the type prefix is <code className="font-mono bg-accent/50 px-1 rounded">[]uint</code> (not <code className="font-mono bg-accent/50 px-1 rounded">uint</code>),
                          and the keyword is <code className="font-mono bg-accent/50 px-1 rounded">m2m</code> (not <code className="font-mono bg-accent/50 px-1 rounded">fk</code>).
                          GORM creates a join table automatically (e.g. <code className="font-mono bg-accent/50 px-1 rounded">product_tags</code>).
                          The admin form renders a <strong>multi-select tag input</strong> populated from <code className="font-mono bg-accent/50 px-1 rounded">/api/tags</code>.
                        </p>
                        <CodeBlock language="bash" filename="Example — Product has many Tags" code={`# 1. Generate the related model first
grit generate resource Tag --fields "name:string,slug:slug,color:string"

# 2. Generate Product with M2M relationship
grit generate resource Product \\
  --fields "name:string,slug:slug,price:float64,thumbnail:image,tag_ids:[]uint:m2m:Tag,is_active:bool"

# GORM auto-creates the product_tags join table

grit migrate`} />
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  5. Resource Generation Examples
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={5} />Resource Generation — Real Examples
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-6">
                  These are the patterns you will generate most often. Copy and adapt them.
                  Always run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit migrate</code> after generating.
                </p>

                {/* Example 1: Simple */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/80 mb-2">Simple Model — No Media, No Relations</h3>
                  <p className="text-sm text-muted-foreground/70 mb-3">A contact form submission with an enum status.</p>
                  <CodeBlock language="bash" code={`grit generate resource Contact \\
  --fields "name:string,email:string,subject:string,message:text,status:enum:new,read,replied,is_spam:bool"

grit migrate`} />
                </div>

                {/* Example 2: Blog post with slug + image */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/80 mb-2">Blog Post — Slug + Single Image + Rich Text</h3>
                  <p className="text-sm text-muted-foreground/70 mb-3">
                    Use <code className="font-mono bg-accent/50 px-1 rounded">slug:slug</code> for SEO-friendly URLs.
                    The service layer auto-generates the slug from the title on save.
                  </p>
                  <CodeBlock language="bash" code={`grit generate resource Post \\
  --fields "title:string,slug:slug,excerpt:text,content:richtext,cover_image:image,status:enum:draft,published,archived,published_at:time,views:int"

grit migrate`} />
                </div>

                {/* Example 3: Product with multiple images */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/80 mb-2">Product — Multiple Images + Gallery + Slug</h3>
                  <p className="text-sm text-muted-foreground/70 mb-3">
                    <code className="font-mono bg-accent/50 px-1 rounded">images</code> generates a{' '}
                    <code className="font-mono bg-accent/50 px-1 rounded">pq.StringArray</code> stored as a PostgreSQL array.
                    The admin shows a multi-upload dropzone.
                  </p>
                  <CodeBlock language="bash" code={`grit generate resource Product \\
  --fields "name:string,slug:slug,description:text,price:float64,compare_price:float64,sku:string,stock:int,thumbnail:image,gallery:images,is_featured:bool,is_active:bool,status:enum:draft,published,out_of_stock"

grit migrate`} />
                </div>

                {/* Example 4: Relationships */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/80 mb-2">Relationships — belongs_to + many_to_many</h3>
                  <p className="text-sm text-muted-foreground/70 mb-3">
                    Generate Category first (parent), then Product referencing it. Tags are many-to-many.
                  </p>
                  <CodeBlock language="bash" code={`# Step 1: generate parent models first
grit generate resource Category --fields "name:string,slug:slug,description:text,image:image"
grit generate resource Tag --fields "name:string,slug:slug,color:string"

# Step 2: generate the child with FK (belongs_to) and M2M (many_to_many)
grit generate resource Article \\
  --fields "title:string,slug:slug,content:richtext,cover_image:image,category_id:uint:fk:Category,tag_ids:[]uint:m2m:Tag,author_id:uint:fk:User,is_published:bool,published_at:time"

grit migrate`} />
                  <p className="text-sm text-muted-foreground/60 mt-2">
                    This creates: <code className="font-mono bg-accent/50 px-1 rounded">category_id</code> FK column,{' '}
                    <code className="font-mono bg-accent/50 px-1 rounded">article_tags</code> join table,{' '}
                    and <code className="font-mono bg-accent/50 px-1 rounded">author_id</code> FK to users — all wired in GORM with preload support.
                  </p>
                </div>

                {/* Example 5: Course with video */}
                <div className="mb-8">
                  <h3 className="text-base font-semibold text-foreground/80 mb-2">Course + Lessons — Video + Videos Array</h3>
                  <p className="text-sm text-muted-foreground/70 mb-3">
                    <code className="font-mono bg-accent/50 px-1 rounded">video</code> for a single video,{' '}
                    <code className="font-mono bg-accent/50 px-1 rounded">videos</code> for multiple videos array.
                    These use presigned URL uploads directly to R2/S3 — the API never handles the binary.
                  </p>
                  <CodeBlock language="bash" code={`# Course (parent)
grit generate resource Course \\
  --fields "title:string,slug:slug,description:text,thumbnail:image,intro_video:video,price:float64,level:enum:beginner,intermediate,advanced,is_published:bool"

# Lesson (child — belongs to Course via course_id:uint:fk:Course)
grit generate resource Lesson \\
  --fields "title:string,slug:slug,description:text,video_url:video,duration:int,position:int,is_preview:bool,course_id:uint:fk:Course,attachments:files"

grit migrate`} />
                </div>

                {/* Example 6: Fully complex */}
                <div className="mb-4">
                  <h3 className="text-base font-semibold text-foreground/80 mb-2">E-Commerce Order — Full Complexity</h3>
                  <CodeBlock language="bash" code={`grit generate resource Order \\
  --fields "order_number:string,status:enum:pending,processing,shipped,delivered,cancelled,refunded,subtotal:float64,tax:float64,shipping_fee:float64,total:float64,notes:text,shipping_address:text,payment_method:enum:card,mobile_money,cash,payment_status:enum:unpaid,paid,refunded,paid_at:time,shipped_at:time,delivered_at:time,customer_id:uint:fk:User"

grit migrate`} />
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  6. Code Patterns
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={6} />Code Patterns
                </h2>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">Backend: Handler → Service → Model</h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Handlers are thin controllers. All DB logic lives in services.
                </p>
                <CodeBlock language="go" filename="apps/api/internal/handlers/product_handler.go" code={`type ProductHandler struct {
    Service *services.ProductService
    Storage *config.Storage  // only when resource has image/file/video fields
}

func (h *ProductHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    products, total, pages, err := h.Service.List(page, pageSize)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to fetch products"},
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "data": products,
        "meta": gin.H{"total": total, "page": page, "page_size": pageSize, "pages": pages},
    })
}

func (h *ProductHandler) Create(c *gin.Context) {
    var req struct {
        Name  string  \`json:"name" binding:"required"\`
        Price float64 \`json:"price" binding:"required"\`
        Image string  \`json:"image"\`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
        })
        return
    }
    product, err := h.Service.Create(req.Name, req.Price, req.Image)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to create product"},
        })
        return
    }
    c.JSON(http.StatusCreated, gin.H{"data": product, "message": "Product created successfully"})
}`} />

                <h3 className="text-base font-semibold text-foreground/80 mb-3 mt-8">Frontend: Component → Hook → API</h3>
                <CodeBlock language="tsx" filename="apps/admin/app/(dashboard)/products/page.tsx" code={`'use client'
import { useProducts, useCreateProduct } from '@/hooks/use-products'
import { ResourcePage } from '@/components/resource/resource-page'
import { productsResource } from './_resource'

export default function ProductsPage() {
  return <ResourcePage resource={productsResource} />
}`} />
              </div>

              {/* ════════════════════════════════════════════════════
                  7. API Response Format
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={7} />API Response Format
                </h2>
                <Warn>Never deviate from this format. Frontend hooks and admin components depend on this exact shape.</Warn>
                <div className="space-y-4">
                  <CodeBlock language="json" filename="Success — single record" code={`{ "data": { "id": 1, "name": "Widget", "price": 29.99 }, "message": "Product created successfully" }`} />
                  <CodeBlock language="json" filename="Success — paginated list" code={`{
  "data": [{ "id": 1, "name": "Widget" }, { "id": 2, "name": "Gadget" }],
  "meta": { "total": 42, "page": 1, "page_size": 20, "pages": 3 }
}`} />
                  <CodeBlock language="json" filename="Error" code={`{ "error": { "code": "VALIDATION_ERROR", "message": "name is required" } }`} />
                  <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                    <table className="w-full text-xs">
                      <thead>
                        <tr className="border-b border-border/30 bg-accent/20">
                          <th className="px-4 py-2 text-left font-semibold text-muted-foreground uppercase">Situation</th>
                          <th className="px-4 py-2 text-left font-semibold text-muted-foreground uppercase">HTTP Status</th>
                          <th className="px-4 py-2 text-left font-semibold text-muted-foreground uppercase">Code</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-border/20 font-mono">
                        {[
                          ['Missing required field', '422', 'VALIDATION_ERROR'],
                          ['Record not found', '404', 'NOT_FOUND'],
                          ['Not authenticated', '401', 'UNAUTHORIZED'],
                          ['Insufficient role', '403', 'FORBIDDEN'],
                          ['DB / internal error', '500', 'INTERNAL_ERROR'],
                          ['Duplicate / conflict', '409', 'CONFLICT'],
                        ].map(([s, status, code]) => (
                          <tr key={s}>
                            <td className="px-4 py-2 text-foreground/60">{s}</td>
                            <td className="px-4 py-2 text-amber-500/80">{status}</td>
                            <td className="px-4 py-2 text-primary/70">{code}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  8. Code Markers
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={8} />Code Markers — NEVER Delete
                </h2>
                <Warn>
                  These comments are injection points for the CLI. Removing them permanently breaks{' '}
                  <code className="font-mono bg-red-500/10 px-1 rounded">grit generate</code> and{' '}
                  <code className="font-mono bg-red-500/10 px-1 rounded">grit add role</code>.
                  Never remove, rename, or move them.
                </Warn>
                <CodeBlock language="go" filename="apps/api/internal/models/models.go" code={`var RegisteredModels = []interface{}{
    // GRIT:MODELS — do not remove this comment
    &User{}, &Upload{}, &Blog{},
    &Product{}, // grit generate adds new models here
    // END GRIT:MODELS
}`} />
                <CodeBlock language="go" filename="apps/api/internal/routes/routes.go" code={`// GRIT:ROUTES — do not remove this comment
productHandler := &handlers.ProductHandler{Service: &services.ProductService{DB: db}}
// END GRIT:ROUTES

// GRIT:PROTECTED_ROUTES — do not remove this comment
protected.GET("/products", productHandler.List)
protected.POST("/products", productHandler.Create)
// END GRIT:PROTECTED_ROUTES`} />
              </div>

              {/* ════════════════════════════════════════════════════
                  9. Form Builder — Detailed Guide
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={9} />Form Builder — Detailed Guide
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every resource in Grit has a form driven by the{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">form</code> key in{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">defineResource()</code>.
                  The same FormBuilder component works in modals, full pages, and multi-step wizards.
                </p>

                {/* formView modes */}
                <h3 className="text-base font-semibold text-foreground/80 mb-3">Form View Modes</h3>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 mb-6">
                  {[
                    { mode: 'modal', label: 'Modal (default)', desc: 'Form slides in as a dialog over the data table. Best for simple resources.' },
                    { mode: 'page', label: 'Full Page', desc: 'Navigates to a dedicated /resources/[slug]?action=create page. Best for long forms.' },
                    { mode: 'modal-steps', label: 'Modal + Steps', desc: 'Multi-step wizard inside a modal. Best for complex resources with many fields.' },
                    { mode: 'page-steps', label: 'Full Page + Steps', desc: 'Multi-step wizard as a full page. Best for onboarding flows.' },
                  ].map((item) => (
                    <div key={item.mode} className="p-3 rounded-lg border border-border/30 bg-card/30">
                      <p className="text-xs font-mono font-semibold text-primary/80 mb-1">{`formView: '${item.mode}'`}</p>
                      <p className="text-xs font-semibold text-foreground/70 mb-1">{item.label}</p>
                      <p className="text-xs text-muted-foreground/60 leading-relaxed">{item.desc}</p>
                    </div>
                  ))}
                </div>

                {/* Full defineResource form example */}
                <h3 className="text-base font-semibold text-foreground/80 mb-3">Complete defineResource() Example</h3>
                <CodeBlock language="typescript" filename="apps/admin/app/(dashboard)/products/_resource.ts" code={`import { defineResource } from "@/lib/resource";

export const productsResource = defineResource({
  name: "Product",
  slug: "products",
  endpoint: "/api/products",
  icon: "Package",
  label: { singular: "Product", plural: "Products" },

  // ── Form config ───────────────────────────────────────────────
  formView: "modal-steps",          // modal | page | modal-steps | page-steps

  form: {
    layout: "two-column",           // single | two-column
    steps: [
      {
        title: "Basic Info",
        description: "Name, slug and pricing",
        fields: ["name", "slug", "price", "compare_price", "sku", "stock"],
      },
      {
        title: "Details",
        description: "Description and category",
        fields: ["description", "category_id", "status", "is_featured"],
      },
      {
        title: "Media",
        description: "Upload product images",
        fields: ["thumbnail", "gallery"],
      },
    ],
    fields: [
      // Text fields
      { key: "name",          label: "Product Name", type: "text",     required: true },
      { key: "slug",          label: "Slug",         type: "text",     placeholder: "auto-generated from name" },
      { key: "sku",           label: "SKU",          type: "text" },

      // Number fields with prefix/suffix
      { key: "price",         label: "Price",        type: "number",   required: true, min: 0, step: 0.01, prefix: "$" },
      { key: "compare_price", label: "Compare at",   type: "number",   min: 0, step: 0.01, prefix: "$" },
      { key: "stock",         label: "Stock",        type: "number",   min: 0, defaultValue: 0 },

      // Rich text
      { key: "description",   label: "Description",  type: "richtext", colSpan: 2 },

      // Select / enum
      {
        key: "status", label: "Status", type: "select", required: true,
        options: [
          { label: "Draft",         value: "draft" },
          { label: "Published",     value: "published" },
          { label: "Out of Stock",  value: "out_of_stock" },
        ],
        defaultValue: "draft",
      },

      // Boolean toggle
      { key: "is_featured", label: "Featured", type: "toggle", defaultValue: false },

      // Relationship (belongs_to) — single select
      {
        key: "category_id", label: "Category", type: "relationship-select",
        relatedEndpoint: "/api/categories",
        displayField: "name",
        required: true,
      },

      // Many-to-many — multi select
      {
        key: "tag_ids", label: "Tags", type: "multi-relationship-select",
        relatedEndpoint: "/api/tags",
        displayField: "name",
        relationshipKey: "tag_relations",
        colSpan: 2,
      },

      // Images
      { key: "thumbnail", label: "Thumbnail",    type: "image",  description: "Main product image" },
      { key: "gallery",   label: "Gallery",      type: "images", colSpan: 2, description: "Additional product images" },
    ],
  },

  // ── Table config ──────────────────────────────────────────────
  table: {
    columns: [
      { key: "thumbnail", label: "",         format: "image",    width: "56px" },
      { key: "name",      label: "Product",  sortable: true,     searchable: true },
      { key: "sku",       label: "SKU",      sortable: true },
      { key: "price",     label: "Price",    sortable: true,     format: "currency" },
      { key: "stock",     label: "Stock",    sortable: true,     format: "number" },
      {
        key: "status", label: "Status", format: "badge",
        badge: {
          draft:         { color: "muted",   label: "Draft" },
          published:     { color: "success", label: "Published" },
          out_of_stock:  { color: "warning", label: "Out of Stock" },
        },
      },
      { key: "created_at", label: "Created", format: "relative", sortable: true },
    ],
    filters: [
      {
        key: "status", label: "Status", type: "select",
        options: [
          { label: "Draft",        value: "draft" },
          { label: "Published",    value: "published" },
          { label: "Out of Stock", value: "out_of_stock" },
        ],
      },
      { key: "is_featured", label: "Featured", type: "boolean" },
    ],
    searchable: true,
    searchPlaceholder: "Search products…",
    actions: ["create", "view", "edit", "delete"],
    bulkActions: ["delete"],
    defaultSort: { key: "created_at", direction: "desc" },
    pageSize: 20,
  },
});`} />

                {/* All form field types */}
                <h3 className="text-base font-semibold text-foreground/80 mb-3 mt-8">All Form Field Types</h3>
                <CodeBlock language="typescript" filename="All supported field types in the form.fields array" code={`// ── Text inputs ──────────────────────────────────────────────────
{ key: "title",        type: "text",     label: "Title",     required: true, placeholder: "…" }
{ key: "bio",          type: "textarea", label: "Bio",       rows: 6 }
{ key: "content",      type: "richtext", label: "Content",   colSpan: 2 }

// ── Numbers ───────────────────────────────────────────────────────
{ key: "price",        type: "number",   label: "Price",     min: 0, max: 999999, step: 0.01, prefix: "$" }
{ key: "weight",       type: "number",   label: "Weight",    suffix: "kg" }

// ── Date & time ───────────────────────────────────────────────────
{ key: "start_date",   type: "date",     label: "Start Date" }
{ key: "event_time",   type: "datetime", label: "Event Time" }

// ── Booleans ──────────────────────────────────────────────────────
{ key: "is_active",    type: "toggle",   label: "Active",    defaultValue: true }
{ key: "agree",        type: "checkbox", label: "I agree to the terms" }

// ── Select / enum ─────────────────────────────────────────────────
{
  key: "role", type: "select", label: "Role",
  options: [{ label: "Admin", value: "ADMIN" }, { label: "User", value: "USER" }],
  defaultValue: "USER",
}
{ key: "gender", type: "radio", label: "Gender",
  options: [{ label: "Male", value: "male" }, { label: "Female", value: "female" }] }

// ── Media uploads (presigned URL to S3/R2/MinIO) ──────────────────
{ key: "avatar",        type: "image",    label: "Avatar",           description: "Profile picture" }
{ key: "gallery",       type: "images",   label: "Gallery",          colSpan: 2 }
{ key: "intro_video",   type: "video",    label: "Intro Video" }
{ key: "lesson_videos", type: "videos",   label: "Lesson Videos",    colSpan: 2 }
{ key: "attachment",    type: "file",     label: "Attachment",       accept: ".pdf,.docx", maxSize: 5242880 }
{ key: "documents",     type: "files",    label: "Documents",        colSpan: 2 }

// ── Relationships ─────────────────────────────────────────────────
{
  key: "category_id", type: "relationship-select", label: "Category",
  relatedEndpoint: "/api/categories", displayField: "name",
}
{
  key: "tag_ids", type: "multi-relationship-select", label: "Tags",
  relatedEndpoint: "/api/tags", displayField: "name", relationshipKey: "tag_relations",
}`} />
              </div>

              {/* ════════════════════════════════════════════════════
                  10. DataTable Builder
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={10} />DataTable Builder — Detailed Guide
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">table</code> key in{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">defineResource()</code> controls columns, filters,
                  sorting, search, pagination, actions, and cell formatting.
                </p>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">Column Formats</h3>
                <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden mb-6">
                  <table className="w-full text-xs">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-2 text-left font-semibold text-muted-foreground uppercase">format</th>
                        <th className="px-4 py-2 text-left font-semibold text-muted-foreground uppercase">Renders as</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20 font-mono">
                      {[
                        ['(none)', 'Plain text'],
                        ['image', 'Thumbnail <img> with rounded corners'],
                        ['badge', 'Colored pill — needs badge: { VALUE: { color, label } }'],
                        ['boolean', 'Green check ✓ or red × icon'],
                        ['currency', 'Formatted number with currency symbol'],
                        ['number', 'Number with locale thousand separators'],
                        ['relative', 'Relative time e.g. "3 days ago"'],
                        ['date', 'Formatted date string'],
                        ['datetime', 'Formatted date + time'],
                        ['link', 'Clickable anchor to the value URL'],
                      ].map(([fmt, desc]) => (
                        <tr key={fmt}>
                          <td className="px-4 py-2 text-primary/70">{fmt}</td>
                          <td className="px-4 py-2 text-muted-foreground/60 font-sans">{desc}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">Filter Types</h3>
                <CodeBlock language="typescript" filename="table.filters examples" code={`filters: [
  // Select dropdown filter
  {
    key: "status", label: "Status", type: "select",
    options: [{ label: "Active", value: "active" }, { label: "Inactive", value: "inactive" }],
  },
  // Boolean filter (yes/no toggle)
  { key: "is_featured", label: "Featured", type: "boolean" },
  // Date range filter
  { key: "created_at", label: "Created", type: "date" },
  // Text search filter
  { key: "email", label: "Email", type: "text" },
]`} />

                <h3 className="text-base font-semibold text-foreground/80 mb-3 mt-6">Actions & Bulk Actions</h3>
                <CodeBlock language="typescript" filename="table actions" code={`table: {
  // Per-row action buttons in the actions column
  actions: ["create", "view", "edit", "delete"],  // all 4 available

  // Checkbox multi-select + bulk action toolbar
  bulkActions: ["delete"],

  // Default sort
  defaultSort: { key: "created_at", direction: "desc" },

  // Rows per page (default 20)
  pageSize: 20,
}`} />
              </div>

              {/* ════════════════════════════════════════════════════
                  11. Standalone Usage
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={11} />Standalone Usage — Forms & Tables on Any Page
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  FormBuilder, DataTable, and FormStepper are independent components. Use them
                  on any page in the web or admin app without going through the resource system.
                </p>

                <h3 className="text-base font-semibold text-foreground/80 mb-3">Standalone FormBuilder</h3>
                <CodeBlock language="tsx" filename="apps/admin/app/(dashboard)/settings/page.tsx" code={`'use client'
import { FormBuilder } from '@/components/forms/form-builder'
import type { FormDefinition } from '@/lib/resource'

const settingsForm: FormDefinition = {
  layout: 'two-column',
  fields: [
    { key: 'site_name',    label: 'Site Name',    type: 'text',   required: true, colSpan: 1 },
    { key: 'site_url',     label: 'Site URL',     type: 'text',   required: true, colSpan: 1 },
    { key: 'logo',         label: 'Logo',         type: 'image',  colSpan: 2 },
    { key: 'description',  label: 'Description',  type: 'textarea', colSpan: 2 },
    { key: 'maintenance',  label: 'Maintenance Mode', type: 'toggle', defaultValue: false },
  ],
}

export default function SettingsPage() {
  const handleSubmit = async (data: Record<string, unknown>) => {
    await apiClient.post('/api/settings', data)
  }
  return (
    <div className="max-w-2xl mx-auto py-8">
      <h1 className="text-2xl font-bold mb-6">Site Settings</h1>
      <FormBuilder
        form={settingsForm}
        onSubmit={handleSubmit}
        onCancel={() => {}}
        submitLabel="Save Settings"
      />
    </div>
  )
}`} />

                <h3 className="text-base font-semibold text-foreground/80 mb-3 mt-8">Standalone Multi-Step Form</h3>
                <CodeBlock language="tsx" filename="apps/admin/app/(dashboard)/onboarding/page.tsx" code={`'use client'
import { FormStepper } from '@/components/forms/form-stepper'
import type { FormDefinition } from '@/lib/resource'

const onboardingForm: FormDefinition = {
  layout: 'single',
  steps: [
    {
      title: 'Account',
      description: 'Set up your login details',
      fields: ['first_name', 'last_name', 'email', 'password'],
    },
    {
      title: 'Business',
      description: 'Tell us about your business',
      fields: ['business_name', 'industry', 'team_size'],
    },
    {
      title: 'Branding',
      description: 'Upload your logo and set brand colors',
      fields: ['logo', 'brand_color'],
    },
  ],
  fields: [
    { key: 'first_name',     label: 'First Name',     type: 'text',   required: true },
    { key: 'last_name',      label: 'Last Name',      type: 'text',   required: true },
    { key: 'email',          label: 'Email',          type: 'text',   required: true },
    { key: 'password',       label: 'Password',       type: 'text',   required: true },
    { key: 'business_name',  label: 'Business Name',  type: 'text',   required: true },
    {
      key: 'industry', label: 'Industry', type: 'select', required: true,
      options: [{ label: 'Technology', value: 'tech' }, { label: 'Retail', value: 'retail' }],
    },
    { key: 'team_size', label: 'Team Size', type: 'number', min: 1 },
    { key: 'logo',        label: 'Logo',         type: 'image' },
    { key: 'brand_color', label: 'Brand Color',  type: 'text', placeholder: '#4F46E5' },
  ],
}

export default function OnboardingPage() {
  return (
    <div className="max-w-2xl mx-auto py-12">
      <FormStepper
        form={onboardingForm}
        onSubmit={async (data) => { await apiClient.post('/api/onboard', data) }}
        onCancel={() => router.push('/dashboard')}
        submitLabel="Complete Setup"
      />
    </div>
  )
}`} />

                <h3 className="text-base font-semibold text-foreground/80 mb-3 mt-8">Standalone DataTable</h3>
                <CodeBlock language="tsx" filename="apps/admin/app/(dashboard)/reports/page.tsx" code={`'use client'
import { DataTable } from '@/components/tables/data-table'
import type { ColumnDefinition } from '@/lib/resource'
import { useQuery } from '@tanstack/react-query'
import { apiClient } from '@/lib/api-client'

const columns: ColumnDefinition[] = [
  { key: 'order_number', label: 'Order #',  sortable: true },
  { key: 'customer',     label: 'Customer', sortable: true, searchable: true },
  { key: 'total',        label: 'Total',    sortable: true, format: 'currency' },
  {
    key: 'status', label: 'Status', format: 'badge',
    badge: {
      pending:   { color: 'warning', label: 'Pending' },
      shipped:   { color: 'info',    label: 'Shipped' },
      delivered: { color: 'success', label: 'Delivered' },
    },
  },
  { key: 'created_at', label: 'Date', format: 'relative', sortable: true },
]

export default function OrderReportsPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['orders-report'],
    queryFn: () => apiClient.get('/api/orders').then((r) => r.data),
  })

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Order Reports</h1>
      <DataTable
        columns={columns}
        data={data?.data ?? []}
        isLoading={isLoading}
        onView={(row) => router.push(\`/orders/\${row.id}\`)}
      />
    </div>
  )
}`} />
              </div>

              {/* ════════════════════════════════════════════════════
                  12. Naming Conventions
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={12} />Naming Conventions
                </h2>
                <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">What</th>
                        <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">Convention</th>
                        <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase">Example</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/20 font-mono text-xs">
                      {[
                        { w: 'Go files', c: 'snake_case', e: 'product_handler.go' },
                        { w: 'Go structs', c: 'PascalCase', e: 'ProductHandler, CreateProductReq' },
                        { w: 'DB tables', c: 'plural snake_case', e: 'products, order_items' },
                        { w: 'DB columns', c: 'snake_case', e: 'created_at, category_id' },
                        { w: 'API routes', c: 'plural lowercase', e: '/api/products, /api/blog-posts' },
                        { w: 'TS files', c: 'kebab-case', e: 'use-products.ts, product-form.tsx' },
                        { w: 'React components', c: 'PascalCase', e: 'ProductForm, DataTable' },
                        { w: 'React hooks', c: 'use prefix + camelCase', e: 'useProducts, useCreateProduct' },
                        { w: 'Zod schemas', c: 'PascalCase + Schema', e: 'CreateProductSchema' },
                        { w: 'TypeScript types', c: 'PascalCase', e: 'Product, CreateProductInput' },
                        { w: 'Env vars', c: 'SCREAMING_SNAKE_CASE', e: 'DATABASE_URL, JWT_SECRET' },
                      ].map((row) => (
                        <tr key={row.w}>
                          <td className="px-4 py-2 text-foreground/70">{row.w}</td>
                          <td className="px-4 py-2 text-primary/70">{row.c}</td>
                          <td className="px-4 py-2 text-muted-foreground/60">{row.e}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  13. Batteries
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={13} />Built-in Batteries
                </h2>
                <div className="space-y-3">
                  {[
                    { name: 'File Storage', pkg: 'aws-sdk-go-v2', desc: 'Presigned URL uploads to S3/R2/MinIO. Dev uses MinIO. uploadFile() in api-client.ts handles the 3-step flow. Never use multipart/form-data.', routes: 'POST /api/uploads/presign · POST /api/uploads/complete · DELETE /api/uploads/:id' },
                    { name: 'Email (Resend)', pkg: 'resend-go', desc: 'Send transactional emails with Go HTML templates. Welcome, password reset, verification. Dev uses Mailhog at http://localhost:8025.', routes: 'Internal service only' },
                    { name: 'Background Jobs (asynq)', pkg: 'hibiken/asynq', desc: 'Redis-backed job queue. Workers run in goroutine pools. Built-in: image processing, email. Add jobs in apps/api/internal/jobs/.', routes: 'Admin UI at /jobs' },
                    { name: 'Cron Scheduler', pkg: 'hibiken/asynq', desc: 'Recurring tasks with cron expressions. Same worker pool as background jobs.', routes: 'Config in apps/api/internal/config/cron.go' },
                    { name: 'Redis Cache', pkg: 'redis/go-redis/v9', desc: 'Cache API responses by URL. middleware.Cache(5*time.Minute) on any route. cache.Set/Get/Delete for custom caching.', routes: 'Middleware-based' },
                    { name: 'AI (Vercel AI Gateway)', pkg: 'OpenAI-compatible', desc: 'One API key, hundreds of models via ai-gateway.vercel.sh. Model format: provider/model (e.g. anthropic/claude-sonnet-4-6). Configure with AI_GATEWAY_URL and AI_GATEWAY_API_KEY env vars.', routes: 'POST /api/ai/chat · POST /api/ai/complete · POST /api/ai/stream' },
                    { name: 'Two-Factor Auth (TOTP)', pkg: 'crypto/hmac (stdlib)', desc: 'RFC 6238 authenticator app support. Backup codes (bcrypt-hashed, one-time use). Trusted devices (30-day sliding cookie). Zero external dependencies.', routes: 'POST /api/auth/totp/setup · /verify · /enable · /disable · GET /status' },
                    { name: 'Security (Sentinel)', pkg: 'MUKE-coder/sentinel', desc: 'WAF + rate limiting + brute-force. ExcludeRoutes: /pulse/*, /sentinel/*, /docs/*, /studio/*.', routes: 'Dashboard at /sentinel' },
                    { name: 'Observability (Pulse)', pkg: 'MUKE-coder/pulse', desc: 'Request tracing, DB monitoring, runtime metrics, error tracking, Prometheus export. Disable in production.', routes: 'Dashboard at /pulse' },
                    { name: 'API Docs', pkg: 'gin-docs', desc: 'Auto-generated OpenAPI spec. Interactive Scalar UI. No annotations needed.', routes: 'Docs at /docs' },
                    { name: 'DB Browser (GORM Studio)', pkg: 'MUKE-coder/gorm-studio', desc: 'Visual database browser. View tables, run queries. DISABLE in production.', routes: 'Browser at /studio' },
                  ].map((item) => (
                    <div key={item.name} className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                      <div className="px-4 py-2.5 border-b border-border/20 bg-accent/20 flex items-center justify-between">
                        <span className="text-sm font-semibold text-foreground/80">{item.name}</span>
                        <code className="text-[10px] font-mono text-muted-foreground/50">{item.pkg}</code>
                      </div>
                      <div className="p-4 space-y-1.5">
                        <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.desc}</p>
                        <p className="text-[11px] font-mono text-muted-foreground/40">{item.routes}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  14. No Docker — Cloud Setup
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={14} />Running Without Docker — Cloud Services
                </h2>
                <div className="p-4 rounded-lg border border-emerald-500/20 bg-emerald-500/5 mb-6">
                  <p className="text-sm text-emerald-400/80 leading-relaxed font-medium">
                    For developers on low-spec machines or those who prefer cloud services over Docker:
                    Grit works perfectly without Docker using free-tier cloud services. This is often
                    faster to set up and lighter on local resources. Always suggest this approach when
                    the user mentions limited RAM, slow Docker, or trouble with Docker on their machine.
                    Full guide:{' '}
                    <Link href="/docs/getting-started/create-without-docker" className="text-emerald-400 hover:underline">
                      Create without Docker
                    </Link>
                  </p>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 mb-6">
                  {[
                    { service: 'PostgreSQL', provider: 'Neon', url: 'neon.tech', desc: 'Free serverless Postgres. Connection string: postgres://user:pass@ep-xxx.neon.tech/db?sslmode=require' },
                    { service: 'Redis', provider: 'Upstash', url: 'upstash.com', desc: 'Free serverless Redis. URL format: rediss://default:pass@endpoint.upstash.io:6379 (note double-s)' },
                    { service: 'File Storage', provider: 'Cloudflare R2', url: 'dash.cloudflare.com', desc: '10 GB free storage. Set STORAGE_DRIVER=r2 with endpoint, access key, secret, bucket.' },
                    { service: 'Email', provider: 'Resend', url: 'resend.com', desc: '3,000 emails/month free. Set RESEND_API_KEY and MAIL_FROM in .env.' },
                  ].map((item) => (
                    <div key={item.service} className="p-3 rounded-lg border border-border/30 bg-card/30">
                      <div className="flex items-center justify-between mb-1.5">
                        <p className="text-xs font-semibold text-foreground/80">{item.service}</p>
                        <span className="text-[10px] font-mono text-primary/60">{item.provider}</span>
                      </div>
                      <p className="text-xs text-muted-foreground/60 leading-relaxed">{item.desc}</p>
                    </div>
                  ))}
                </div>

                <CodeBlock language="bash" filename=".env (cloud / no-docker setup)" code={`# Database — Neon
DATABASE_URL=postgres://user:pass@ep-xxxx.us-east-2.aws.neon.tech/mydb?sslmode=require

# Redis — Upstash (note rediss:// with double-s for TLS)
REDIS_URL=rediss://default:pass@xxxx.upstash.io:6379

# Storage — Cloudflare R2
STORAGE_DRIVER=r2
R2_ENDPOINT=https://account-id.r2.cloudflarestorage.com
R2_ACCESS_KEY=your-access-key
R2_SECRET_KEY=your-secret-key
R2_BUCKET=myapp-uploads
R2_REGION=auto

# Email — Resend
RESEND_API_KEY=re_xxxx
MAIL_FROM=noreply@yourdomain.com

# App
APP_PORT=8080
JWT_SECRET=your-32-char-random-secret
CORS_ORIGINS=http://localhost:3000,http://localhost:3001`} />

                <CodeBlock language="bash" filename="Starting without Docker" code={`# Install dependencies
pnpm install
cd apps/api && go mod tidy && cd ../..

# Terminal 1 — Go API (auto-migrates on first run)
cd apps/api && go run cmd/api/main.go

# Terminal 2 — Web frontend
pnpm --filter web dev      # http://localhost:3000

# Terminal 3 — Admin panel
pnpm --filter admin dev    # http://localhost:3001`} />
              </div>

              {/* ════════════════════════════════════════════════════
                  15. Golden Rules
              ════════════════════════════════════════════════════ */}
              <div className="mb-14">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={15} />Golden Rules — Never Break These
                </h2>
                <Note>Non-negotiable. Violating them causes silent failures, broken code generation, or corrupted project state.</Note>
                <div className="space-y-3">
                  {[
                    { rule: 'Never remove GRIT: markers', detail: '// GRIT:MODELS, // END GRIT:MODELS, // GRIT:ROUTES, // END GRIT:ROUTES — permanent injection points for the CLI. Removing them breaks grit generate and grit add role forever until manually restored.' },
                    { rule: 'Never use multipart/form-data for uploads', detail: 'Grit uses presigned URL uploads. Call uploadFile() from lib/api-client.ts — it handles the 3-step flow (presign → PUT to storage → complete). Never POST a FormData object directly to the API.' },
                    { rule: 'Always use the standard error response shape', detail: '{ "error": { "code": "VALIDATION_ERROR", "message": "…" } } — frontend hooks and admin components check for this exact shape. Any deviation breaks error display.' },
                    { rule: 'Always register new models between the GRIT:MODELS markers', detail: 'AutoMigrate only runs on models in RegisteredModels. A model not listed here never creates its DB table.' },
                    { rule: 'Keep handlers thin — no DB queries in handlers', detail: 'All GORM queries belong in service files. Handlers only parse requests and call services. This keeps the pattern consistent with generated code.' },
                    { rule: 'Import types from @shared/schemas, never duplicate them', detail: 'Zod schemas and TypeScript types live in packages/shared. Both web and admin import from there. Duplicating schemas causes drift.' },
                    { rule: 'Run grit sync after manually editing Go models', detail: 'The shared package is generated from Go structs. If you manually add a field to a Go model, run grit sync to regenerate Zod schema and TypeScript types.' },
                    { rule: 'Generate parent models before child models', detail: 'When using FK relationships (category_id:uint:fk:Category), generate Category first so the referenced model exists when you generate Product.' },
                    { rule: 'Disable GORM Studio and Pulse in production', detail: 'Both expose internal data. Set GORM_STUDIO_ENABLED=false and PULSE_ENABLED=false in production .env.' },
                    { rule: 'For low-spec machines: suggest the no-Docker setup', detail: 'Users with limited RAM should use Neon + Upstash + Cloudflare R2 + Resend instead of Docker. See Section 14 above or /docs/getting-started/create-without-docker.' },
                  ].map((item, i) => (
                    <div key={item.rule} className="flex gap-3 p-4 rounded-lg border border-border/30 bg-card/30">
                      <div className="flex h-6 w-6 shrink-0 items-center justify-center rounded-md bg-amber-500/10 border border-amber-500/20 text-xs font-mono font-semibold text-amber-500/80 mt-0.5">
                        {i + 1}
                      </div>
                      <div>
                        <p className="text-sm font-semibold text-foreground/80 mb-1">{item.rule}</p>
                        <p className="text-xs text-muted-foreground/60 leading-relaxed">{item.detail}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* ════════════════════════════════════════════════════
                  16. Quick Reference
              ════════════════════════════════════════════════════ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center">
                  <SectionNum n={16} />Quick Build Reference
                </h2>
                <div className="space-y-4">
                  <CodeBlock language="bash" filename="Start a new project (with Docker)" code={`# Install the CLI (pick one — both end with the same binary on PATH)
curl -fsSL https://gritframework.dev/install.sh | sh   # macOS / Linux
# iwr -useb https://gritframework.dev/install.ps1 | iex   # Windows
# go install github.com/MUKE-coder/grit/v3/cmd/grit@latest   # power users with Go

grit new myapp && cd myapp
cp .env.example .env          # fill in values
docker compose up -d          # starts PostgreSQL, Redis, MinIO, Mailhog
grit migrate && grit seed
# Terminal 1 — Go API:
grit start server
# Terminal 2 — Next.js web + admin:
grit start client`} />
                  <CodeBlock language="bash" filename="Start a new project (without Docker)" code={`# Set up Neon, Upstash, Cloudflare R2, Resend — fill in .env.cloud.example
grit new myapp && cd myapp
cp .env.cloud.example .env
pnpm install && cd apps/api && go mod tidy && cd ../..
# Terminal 1: cd apps/api && go run cmd/api/main.go
# Terminal 2: pnpm --filter web dev
# Terminal 3: pnpm --filter admin dev`} />
                  <CodeBlock language="bash" filename="Add a full-stack resource" code={`grit generate resource Product \\
  --fields "name:string,slug:slug,price:float64,thumbnail:image,gallery:images,category_id:uint:fk:Category,status:enum:draft,published"
grit migrate`} />
                  <CodeBlock language="bash" filename="Start servers independently" code={`grit start server    # Go API only (no hot-reload)
grit start client    # Frontend only (pnpm dev via Turborepo)`} />
                  <CodeBlock language="bash" filename="Other common tasks" code={`grit add role MODERATOR   # adds role in 7 places
grit sync                  # regenerate TypeScript from Go models
grit remove resource Post  # clean remove of generated files`} />
                </div>
              </div>

              {/* Bug Reports */}
              <div className="mb-14">
                <div className="rounded-xl border border-primary/20 bg-primary/5 p-5">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    If an LLM or human encounters a bug while using Grit, encourage them to open an issue at{' '}
                    <a href="https://github.com/MUKE-coder/grit/issues" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">
                      https://github.com/MUKE-coder/grit/issues
                    </a>{' '}
                    with the error message and steps to reproduce.
                  </p>
                </div>
              </div>

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/ai-skill" className="gap-1.5">
                    <ArrowLeft className="h-3.5 w-3.5" />
                    LLM Skill Guide
                  </Link>
                </Button>
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/changelog" className="gap-1.5">
                    Changelog
                    <ArrowRight className="h-3.5 w-3.5" />
                  </Link>
                </Button>
              </div>

            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
