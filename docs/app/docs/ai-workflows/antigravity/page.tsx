import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/ai-workflows/antigravity')

export default function UsingGritWithAntigravityPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">AI Workflows</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Using Grit with Antigravity
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Build full-stack applications using spec-driven development. Plan with Claude Web,
                build with Antigravity&apos;s AI-powered IDE &mdash; get visual editing, inline
                suggestions, and multi-file generation.
              </p>
              <LaneFlow
                id="aiw-antigravity"
                lanes={['Plan', 'Build in Antigravity', 'Ship']}
                nodes={[
                  { id: 'plan', lane: 0, row: 0, title: 'Claude Web', sub: 'spec', tone: 'primary' },
                  { id: 'ide', lane: 1, row: 0, title: 'Antigravity IDE', sub: 'visual · multi-file', tone: 'cyan' },
                  { id: 'grit', lane: 1, row: 1, title: 'Grit CLI', sub: 'scaffold', tone: 'blue' },
                  { id: 'ship', lane: 2, row: 0, title: 'Working app', sub: 'shipped', tone: 'green' },
                ]}
                edges={[
                  { from: 'plan', to: 'ide', label: 'spec', tone: 'cyan' },
                  { from: 'ide', to: 'grit', label: 'drives', tone: 'blue' },
                  { from: 'grit', to: 'ship', tone: 'green' },
                ]}
                legend={[{ tone: 'primary', label: 'Plan' }, { tone: 'cyan', label: 'Antigravity IDE' }]}
                caption="Plan with Claude Web, build in Antigravity's AI IDE driving Grit — visual, multi-file, fast"
              />
            </div>

            <div className="prose-grit">
              {/* The Workflow */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The Workflow
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The spec-driven workflow is a 3-step process that separates planning from building.
                  You use <strong className="text-foreground/90">Claude Web</strong> (claude.ai) for
                  high-level planning, then hand off a structured spec to{' '}
                  <strong className="text-foreground/90">Antigravity</strong> for implementation.
                </p>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 mb-6">
                  {[
                    { step: '1', title: 'Plan', desc: 'Describe your project to Claude Web and get 3 spec files back.' },
                    { step: '2', title: 'Spec', desc: 'project-description.md, antigravity-prompt.md, project-phases.md' },
                    { step: '3', title: 'Build', desc: 'Feed the spec to Antigravity and let its AI composer build your app.' },
                  ].map((item) => (
                    <div key={item.step} className="p-4 rounded-xl border border-border/40 bg-card/50 text-center">
                      <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 border border-primary/15 text-sm font-mono font-bold text-primary mx-auto mb-2">
                        {item.step}
                      </div>
                      <h3 className="text-sm font-semibold mb-1">{item.title}</h3>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.desc}</p>
                    </div>
                  ))}
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Antigravity is a Cursor-like AI-powered IDE that provides inline code editing,
                  a multi-file composer, and an integrated AI chat. It can run terminal commands,
                  apply diffs across multiple files, and let you review every change visually
                  before accepting it &mdash; making it an ideal companion for Grit&apos;s structured
                  code generation workflow.
                </p>
                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mb-6">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-primary/90">Why spec-driven?</strong> AI assistants
                    produce dramatically better code when given a complete plan upfront. Instead of
                    going back and forth, you front-load the thinking in Claude Web and let
                    Antigravity execute the plan methodically.
                  </p>
                </div>
              </div>

              {/* Step 1: Planning with Claude Web */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Step 1: Planning with Claude Web
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Open{' '}
                  <Link href="https://claude.ai" target="_blank" className="text-primary hover:underline">
                    claude.ai
                  </Link>{' '}
                  and describe the application you want to build. Be specific about models,
                  relationships, user roles, and features. Claude will generate 3 structured
                  files that Antigravity&apos;s AI can follow step by step.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Here&apos;s an example using a <strong className="text-foreground/90">Job Board</strong> application.
                  Copy the prompt template below, replace the bracketed sections with your project details,
                  and paste it into Claude Web:
                </p>
                <CodeBlock filename="Prompt for Claude Web" code={`I want to build a [PROJECT NAME] using the Grit framework.

Grit is a full-stack meta-framework that uses:
- Go backend (Gin web framework + GORM ORM)
- Next.js frontend (React + App Router + TypeScript)
- Admin panel (resource-based, Filament-like)
- Monorepo with shared types and validation (Zod)
- Docker for infrastructure (PostgreSQL, Redis, MinIO)

Here's what I want to build:
[DESCRIBE YOUR PROJECT IN DETAIL]

Please create these 3 files for me:

1. project-description.md — Full project spec with all models,
   their fields (use Grit field types: string, text, richtext,
   int, uint, float, bool, datetime, date, slug, string_array,
   belongs_to, many_to_many, file, files), relationships,
   user roles, and features.

2. antigravity-prompt.md — Step-by-step instructions that I'll
   give to Antigravity AI IDE. Include exact grit CLI commands
   to run in order. Use \`grit generate resource <Name>
   --fields "..."\` syntax.

3. project-phases.md — A phased build plan with checkboxes.
   Phase 1: scaffold + core models.
   Phase 2: relationships + complex features.
   Phase 3: frontend customization + polish.`} />
                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mb-6">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-primary/90">Tip:</strong> The more detail you provide about your
                    project, the better the spec files will be. Include specific field names, relationships
                    between models, user roles, and any business logic requirements.
                  </p>
                </div>
              </div>

              {/* The 3 Spec Files */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The 3 Spec Files
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Claude Web will generate these 3 files. Save each one to your project folder
                  before opening Antigravity.
                </p>
                <div className="space-y-4 mb-6">
                  {[
                    {
                      file: 'project-description.md',
                      title: 'Project Description',
                      desc: 'The complete project specification. Contains every model with its fields and types, all relationships between models, user roles and permissions, and a description of every feature. This is the source of truth for your application.',
                    },
                    {
                      file: 'antigravity-prompt.md',
                      title: 'Antigravity Prompt',
                      desc: 'Step-by-step instructions formatted for Antigravity\'s AI composer. Includes exact grit CLI commands in the correct order, file modifications to make after generation, and customization instructions. Feed this directly into Antigravity\'s AI chat.',
                    },
                    {
                      file: 'project-phases.md',
                      title: 'Project Phases',
                      desc: 'A phased build plan with checkboxes for tracking progress. Phase 1 covers scaffolding and core models. Phase 2 adds relationships and complex features. Phase 3 handles frontend customization and polish.',
                    },
                  ].map((item) => (
                    <div key={item.file} className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                      <div className="p-4">
                        <div className="flex items-center gap-2 mb-2">
                          <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">{item.file}</code>
                          <span className="text-sm font-semibold">{item.title}</span>
                        </div>
                        <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.desc}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Step 2: Setting Up Antigravity */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Step 2: Setting Up Antigravity
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Before you start building, configure Antigravity so its AI understands
                  Grit&apos;s conventions and can generate correct code from the start.
                </p>

                <div className="mb-8">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      1
                    </div>
                    <h3 className="text-lg font-semibold tracking-tight">
                      Open Antigravity IDE
                    </h3>
                  </div>
                  <p className="text-muted-foreground leading-relaxed mb-4 pl-10">
                    Launch Antigravity and either create a new project folder or open an
                    existing one where you want to scaffold your Grit application.
                  </p>
                </div>

                <div className="mb-8">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      2
                    </div>
                    <h3 className="text-lg font-semibold tracking-tight">
                      Configure AI context
                    </h3>
                  </div>
                  <p className="text-muted-foreground leading-relaxed mb-4 pl-10">
                    Add the{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">GRIT_SKILL.md</code>{' '}
                    file to Antigravity&apos;s AI context. You can do this by:
                  </p>
                  <ul className="space-y-2 pl-10 mb-4">
                    <li className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-1">&#10003;</span>
                      Adding it as a <strong className="text-foreground/90">rules file</strong> in Antigravity&apos;s AI settings (preferred &mdash; persists across sessions)
                    </li>
                    <li className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-1">&#10003;</span>
                      Pinning it to the <strong className="text-foreground/90">context window</strong> so the AI always has it available
                    </li>
                    <li className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-1">&#10003;</span>
                      Placing it in the project root &mdash; Antigravity will index it as part of the workspace
                    </li>
                  </ul>
                  <div className="pl-10 p-4 rounded-lg border border-primary/20 bg-primary/5 mb-4">
                    <p className="text-sm text-muted-foreground leading-relaxed">
                      <strong className="text-primary/90">Tip:</strong> The GRIT_SKILL.md file teaches Antigravity&apos;s
                      AI everything about Grit&apos;s conventions, CLI commands, project structure, and code
                      patterns. Without it, the AI may generate code that doesn&apos;t follow Grit&apos;s standards.
                    </p>
                  </div>
                </div>

                <div className="mb-8">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      3
                    </div>
                    <h3 className="text-lg font-semibold tracking-tight">
                      Open the integrated terminal
                    </h3>
                  </div>
                  <p className="text-muted-foreground leading-relaxed pl-10">
                    You&apos;ll need the terminal to run Grit CLI commands like{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new</code>,{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>, and{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker compose up -d</code>.
                    Open Antigravity&apos;s integrated terminal so you can run these commands without
                    leaving the IDE.
                  </p>
                </div>
              </div>

              {/* Step 3: Building with Antigravity */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Step 3: Building with Antigravity
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  With your spec files ready and Antigravity configured, it&apos;s time to build.
                  Here&apos;s the process:
                </p>

                <div className="mb-8">
                  <h3 className="text-lg font-semibold tracking-tight mb-3">
                    Scaffold the project
                  </h3>
                  <p className="text-muted-foreground leading-relaxed mb-4">
                    Use the integrated terminal to scaffold a new Grit project and start the
                    Docker infrastructure:
                  </p>
                  <CodeBlock language="bash" filename="Terminal" code={`grit new jobboard
cd jobboard
docker compose up -d`} />
                </div>

                <div className="mb-8">
                  <h3 className="text-lg font-semibold tracking-tight mb-3">
                    Open the AI Composer
                  </h3>
                  <p className="text-muted-foreground leading-relaxed mb-4">
                    Open Antigravity&apos;s AI Composer (the multi-file chat interface) and paste in
                    the contents of your{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">antigravity-prompt.md</code>{' '}
                    file. The AI will read through the step-by-step instructions and begin
                    executing them.
                  </p>
                </div>

                <div className="mb-8">
                  <h3 className="text-lg font-semibold tracking-tight mb-3">
                    What Antigravity does
                  </h3>
                  <p className="text-muted-foreground leading-relaxed mb-4">
                    With the prompt loaded, Antigravity&apos;s AI will:
                  </p>
                  <ul className="space-y-2.5 mb-4">
                    {[
                      'Run grit CLI commands in the integrated terminal to generate resources',
                      'Edit generated files inline with visual diffs you can review',
                      'Make multi-file changes across the Go backend, TypeScript types, and React frontend simultaneously',
                      'Add custom business logic to handlers and services',
                      'Customize admin panel pages, forms, and dashboards',
                    ].map((item) => (
                      <li key={item} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                        <span className="text-primary mt-1">&#10003;</span>
                        {item}
                      </li>
                    ))}
                  </ul>
                </div>

                <div className="mb-8">
                  <h3 className="text-lg font-semibold tracking-tight mb-3">
                    Review and iterate
                  </h3>
                  <p className="text-muted-foreground leading-relaxed mb-4">
                    Use Antigravity&apos;s inline diff view to review every change before accepting it.
                    If something doesn&apos;t look right, ask the AI to adjust it. You can also use
                    the inline chat to make quick edits to individual files without opening the
                    full composer.
                  </p>
                  <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mb-6">
                    <p className="text-sm text-muted-foreground leading-relaxed">
                      <strong className="text-primary/90">Tip:</strong> Work through the phases in your{' '}
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">project-phases.md</code>{' '}
                      file one at a time. Complete Phase 1 (scaffold + core models) before moving to Phase 2
                      (relationships + complex features). This keeps the AI focused and reduces errors.
                    </p>
                  </div>
                </div>
              </div>

              {/* The Grit Skill File */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The Grit Skill File
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">GRIT_SKILL.md</code> file
                  is a structured context document that teaches Antigravity&apos;s AI how Grit projects work.
                  Every project scaffolded with{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new</code>{' '}
                  includes this file automatically. Add it to Antigravity&apos;s AI context so the
                  AI understands your project&apos;s conventions from the start.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Here is the full content of the skill file:
                </p>
                <CodeBlock filename="GRIT_SKILL.md" code={`# Grit Framework — AI Skill Context

> Give this file to any AI coding assistant to teach it
> how to work with your Grit project.

## What is Grit?

Grit is a full-stack framework: Go (Gin + GORM) backend +
Next.js (React + TypeScript) frontend in a monorepo.

## CLI Commands

grit new <name>                    # Scaffold full project
grit new <name> --api              # API only
grit generate resource <Name> --fields "field:type,..."
                                   # Full-stack CRUD
                                   # add --seed --faker --count 50 for data
grit generate seeder <Resource>    # Seeder for an existing resource
grit remove resource <Name>        # Remove a resource
grit sync                          # Sync Go types → TypeScript
grit migrate                       # Run GORM AutoMigrate
grit seed                          # Run all seeders (idempotent)
grit start server                  # Run the Go API
grit start client                  # Run web + admin
grit add role <ROLE>               # Add a new role

## Field Types (exactly 15 — anything else fails validation)

string, text, richtext, int, uint, float, bool, datetime,
date, slug, string_array
belongs_to (e.g., category:belongs_to or author:belongs_to:User)
many_to_many (e.g., tags:many_to_many:Tag)
file  (e.g., image:file:image, doc:file:all)
files (e.g., gallery:files:image)

## Models

UUID string primary keys (NOT gorm.Model, NOT uint):
  ID string \`gorm:"type:varchar(36);primaryKey"\`
belongs_to FK columns are string (varchar(36)).
User has FirstName + LastName; roles are UPPERCASE:
ADMIN, EDITOR, USER (default USER).

## Project Structure

project/
├── apps/api/           # Go backend (Gin + GORM)
│   ├── cmd/server/     # Entry point
│   └── internal/       # models, handlers, services,
│                       # middleware, routes
├── apps/web/           # Next.js frontend
├── apps/admin/         # Admin panel (resource-based)
├── packages/shared/    # Zod schemas, TS types, constants
└── docker-compose.yml  # PostgreSQL, Redis, MinIO, Mailhog

## Conventions

- Go models: PascalCase structs with gorm/json/binding tags
- API responses:
  { "data": ..., "meta": { "total", "page",
    "page_size", "pages" } }
- Errors:
  { "error": { "code": "...", "message": "..." } }
- TS files: kebab-case (use-users.ts, api-client.ts)
- Routes: plural, lowercase (/api/users, /api/posts)
- Validation: Zod schemas in packages/shared/schemas/
- Data fetching: React Query hooks (useQuery, useMutation)
- Styling: Tailwind CSS + shadcn/ui, dark theme default
- Admin: defineResource() API with DataTable, FormBuilder,
  widgets

## Code Generation Markers (DO NOT REMOVE — lowercase, no END markers)

// grit:models, // grit:routes:protected,
// grit:routes:admin, // grit:routes:custom,
// grit:schemas, // grit:resources, // grit:seeders

## How to Build with Grit

1. grit new myapp && cd myapp && docker compose up -d
2. grit generate resource <Name> --fields "..."
   for each model (seed parents before children;
   add --seed --faker to generate sample data)
3. grit migrate && grit seed
4. Terminal 1: grit start server
   Terminal 2: grit start client
5. Customize handlers, add business logic,
   build frontend pages`} />
              </div>

              {/* Antigravity-Specific Tips */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Antigravity-Specific Tips
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Get the most out of Antigravity&apos;s AI features when building with Grit:
                </p>
                <div className="space-y-4">
                  {[
                    {
                      title: 'Pin GRIT_SKILL.md as context',
                      desc: 'Add GRIT_SKILL.md as a pinned context file in Antigravity\'s AI settings. This ensures the AI always has Grit\'s conventions loaded without you needing to paste the file into every conversation.',
                    },
                    {
                      title: 'Use the Composer for multi-file changes',
                      desc: 'When adding a new feature that spans the model, handler, and frontend (e.g., adding a new field or creating a custom endpoint), use the Composer. It can apply changes across multiple files simultaneously and show you diffs for each one.',
                    },
                    {
                      title: 'Use inline chat for quick edits',
                      desc: 'For small changes to a single file — like tweaking a handler\'s validation logic or adjusting a column in the admin DataTable — use the inline chat. Select the code, ask the AI to modify it, and accept the diff.',
                    },
                    {
                      title: 'Run grit CLI commands in the integrated terminal',
                      desc: 'Always use Antigravity\'s integrated terminal for CLI commands like grit generate resource, grit migrate, and docker compose up. This keeps everything in one window and lets the AI see the terminal output if needed.',
                    },
                    {
                      title: 'Review diffs before accepting',
                      desc: 'Antigravity shows a visual diff for every AI-generated change. Always review these before accepting — especially for generated Go code where incorrect struct tags or missing error handling can cause runtime issues.',
                    },
                    {
                      title: 'Install the Go extension',
                      desc: 'When working with Go files, make sure the Go language extension is installed in Antigravity for proper type checking, auto-completion, and error detection. This helps catch issues the AI might miss before you even run the code.',
                    },
                  ].map((item) => (
                    <div key={item.title} className="p-4 rounded-lg border border-border/30 bg-card/30">
                      <h3 className="text-sm font-semibold mb-1.5">{item.title}</h3>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.desc}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Example Prompt for Claude Web */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Example Prompt for Claude Web
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Copy this prompt template and paste it into{' '}
                  <Link href="https://claude.ai" target="_blank" className="text-primary hover:underline">
                    claude.ai
                  </Link>{' '}
                  to generate your 3 spec files. Replace the bracketed sections with your
                  project details.
                </p>
                <CodeBlock filename="prompt-template.md" code={`I want to build a [PROJECT NAME] using the Grit framework.

Grit is a full-stack meta-framework that uses:
- Go backend (Gin web framework + GORM ORM)
- Next.js frontend (React + App Router + TypeScript)
- Admin panel (resource-based, Filament-like)
- Monorepo with shared types and validation (Zod)
- Docker for infrastructure (PostgreSQL, Redis, MinIO)

Here's what I want to build:
[DESCRIBE YOUR PROJECT IN DETAIL]

Please create these 3 files for me:

1. project-description.md — Full project spec with all
   models, their fields (use Grit field types: string,
   text, richtext, int, uint, float, bool, datetime, date,
   slug, string_array, belongs_to, many_to_many, file,
   files), relationships, user roles, and features.

2. antigravity-prompt.md — Step-by-step instructions that
   I'll give to Antigravity AI IDE. Include exact grit CLI
   commands to run in order. Use \`grit generate resource
   <Name> --fields "..."\` syntax.

3. project-phases.md — A phased build plan with checkboxes.
   Phase 1: scaffold + core models.
   Phase 2: relationships + complex features.
   Phase 3: frontend customization + polish.`} />
                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mb-6">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-primary/90">Tip:</strong> For a Job Board, you might describe it as:
                    &quot;A job board where companies can post job listings and candidates can apply.
                    Models include Company, Job, Application, and Category. Companies have a name,
                    logo, website, and description. Jobs belong to a company and a category, with
                    fields for title, description, salary range, location, and type (full-time,
                    part-time, contract). Applications belong to a job and track candidate name,
                    email, resume URL, cover letter, and status.&quot;
                  </p>
                </div>
              </div>

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/ai-workflows/claude" className="gap-1.5">
                    <ArrowLeft className="h-3.5 w-3.5" />
                    Using Grit with Claude
                  </Link>
                </Button>
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/ai-skill" className="gap-1.5">
                    LLM Skill Guide
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
