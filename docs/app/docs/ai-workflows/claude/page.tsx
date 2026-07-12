import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/ai-workflows/claude')

export default function UsingGritWithClaudePage() {
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
                Using Grit with Claude
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Build entire full-stack applications using AI spec-driven development.
                Plan with Claude Web, build with Claude Code &mdash; ship in hours, not weeks.
              </p>
            </div>

            <div className="prose-grit">
              {/* ============================================================ */}
              {/* THE WORKFLOW */}
              {/* ============================================================ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The Workflow
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-6">
                  The Grit + Claude workflow is a 3-step process that turns a project idea into a
                  fully scaffolded, working application. You use <strong className="text-foreground/90">Claude Web</strong> (claude.ai)
                  to plan and generate spec files, then hand those specs to <strong className="text-foreground/90">Claude Code</strong> (the
                  CLI agent) to execute every command and build the entire project automatically.
                </p>

                {/* Visual flow diagram */}
                <div className="flex flex-col md:flex-row items-stretch gap-3 mb-6">
                  <div className="flex-1 rounded-xl border border-border/40 bg-card/50 p-5 text-center">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 border border-primary/15 text-sm font-mono font-bold text-primary mx-auto mb-3">
                      1
                    </div>
                    <h4 className="text-sm font-semibold mb-1.5">Plan</h4>
                    <p className="text-xs text-muted-foreground/70 leading-relaxed">
                      Open <strong>claude.ai</strong> and describe your project. Ask Claude to generate 3 spec files.
                    </p>
                  </div>
                  <div className="hidden md:flex items-center text-muted-foreground/30 text-2xl font-light">&rarr;</div>
                  <div className="flex-1 rounded-xl border border-border/40 bg-card/50 p-5 text-center">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 border border-primary/15 text-sm font-mono font-bold text-primary mx-auto mb-3">
                      2
                    </div>
                    <h4 className="text-sm font-semibold mb-1.5">Spec</h4>
                    <p className="text-xs text-muted-foreground/70 leading-relaxed">
                      Claude outputs 3 files: <strong>project-description.md</strong>, <strong>claude-code-prompt.md</strong>, and <strong>project-phases.md</strong>.
                    </p>
                  </div>
                  <div className="hidden md:flex items-center text-muted-foreground/30 text-2xl font-light">&rarr;</div>
                  <div className="flex-1 rounded-xl border border-border/40 bg-card/50 p-5 text-center">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 border border-primary/15 text-sm font-mono font-bold text-primary mx-auto mb-3">
                      3
                    </div>
                    <h4 className="text-sm font-semibold mb-1.5">Build</h4>
                    <p className="text-xs text-muted-foreground/70 leading-relaxed">
                      Feed the specs to <strong>Claude Code</strong>. It runs every grit command, generates resources, and builds the app.
                    </p>
                  </div>
                </div>

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mb-6">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-primary/90">Why two Claudes?</strong> Claude Web excels at
                    creative planning and writing detailed specifications. Claude Code excels at
                    executing commands, editing files, and building working software. Together, they
                    cover the full spectrum from idea to deployed application.
                  </p>
                </div>
              </div>

              {/* ============================================================ */}
              {/* STEP 1: PLANNING WITH CLAUDE WEB */}
              {/* ============================================================ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Step 1: Planning with Claude Web
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Open{' '}
                  <Link href="https://claude.ai" target="_blank" className="text-primary hover:underline">
                    claude.ai
                  </Link>{' '}
                  and describe your project idea in detail. The more specific you are about your
                  models, fields, relationships, user roles, and features, the better the output
                  will be. Tell Claude about Grit so it generates specs that use the correct
                  CLI commands and field types.
                </p>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Here is the exact prompt template you can paste into Claude Web. Replace the
                  bracketed sections with your own project details:
                </p>

                <CodeBlock filename="Prompt template for Claude Web" code={`I want to build a [PROJECT NAME] using the Grit framework.

Grit is a full-stack meta-framework that uses:
- Go backend (Gin web framework + GORM ORM)
- Next.js frontend (React + App Router + TypeScript)
- Admin panel (resource-based, Filament-like)
- Monorepo with shared types and validation (Zod)
- Docker for infrastructure (PostgreSQL, Redis, MinIO)

Here's what I want to build:
[DESCRIBE YOUR PROJECT IN DETAIL]

Please create these 3 files for me:

1. project-description.md — Full project spec with all models, their
   fields (use Grit field types: string, text, richtext, int, uint,
   float, bool, datetime, date, slug, string_array, belongs_to,
   many_to_many, file, files), relationships, user roles, and features.

2. claude-code-prompt.md — Step-by-step instructions that I'll give
   to Claude Code (AI coding agent). Include exact grit CLI commands
   to run in order. Use \`grit generate resource <Name> --fields "..."\`
   syntax.

3. project-phases.md — A phased build plan with checkboxes. Phase 1:
   scaffold + core models. Phase 2: relationships + complex features.
   Phase 3: frontend customization + polish.`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-8">
                  Example: Job Board Application
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Here is what a filled-in prompt looks like for a Job Board app:
                </p>

                <CodeBlock filename="Example prompt for a Job Board" code={`I want to build a Job Board using the Grit framework.

Grit is a full-stack meta-framework that uses:
- Go backend (Gin web framework + GORM ORM)
- Next.js frontend (React + App Router + TypeScript)
- Admin panel (resource-based, Filament-like)
- Monorepo with shared types and validation (Zod)
- Docker for infrastructure (PostgreSQL, Redis, MinIO)

Here's what I want to build:
A job board where companies can post job listings and
candidates can browse and apply. Companies have profiles
with a name, logo, website, and description. Jobs have a
title, description, location, salary range (min and max),
job type (full-time, part-time, contract, remote), and
a category. Applications track which user applied to which
job, with a cover letter and resume URL. There should be
an admin panel for managing all resources, and a public
frontend for browsing published jobs.

User roles: ADMIN (manages everything), EMPLOYER (posts
jobs for their company), USER (browses and applies).

Please create these 3 files for me:

1. project-description.md
2. claude-code-prompt.md
3. project-phases.md`} />

                <p className="text-muted-foreground leading-relaxed">
                  Claude Web will respond with all 3 files, fully detailed and ready to use.
                  Copy each file into your project directory before moving to Step 2.
                </p>
              </div>

              {/* ============================================================ */}
              {/* THE 3 SPEC FILES */}
              {/* ============================================================ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The 3 Spec Files
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-6">
                  Claude Web generates three Markdown files that together form a complete
                  blueprint for your application. Each file serves a different purpose in the
                  build pipeline.
                </p>

                {/* project-description.md */}
                <div className="mb-8">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      1
                    </div>
                    <h3 className="text-lg font-semibold tracking-tight">
                      project-description.md
                    </h3>
                  </div>
                  <p className="text-muted-foreground leading-relaxed mb-4 pl-10">
                    The full project specification. Contains the project name, description,
                    every model with its fields and relationships, user roles and permissions,
                    feature list, and API endpoints. This is the single source of truth for
                    what the application does.
                  </p>

                  <CodeBlock filename="project-description.md (Job Board example)" code={`# Job Board — Project Description

## Overview
A full-stack job board where companies post listings
and candidates browse and apply.

## Models

### Company
- name: string (required)
- slug: slug (unique)
- logo: string
- website: string
- description: text
- user: belongs_to (the employer who owns it)

### Job
- title: string (required)
- slug: slug (unique)
- description: text (required)
- location: string
- salary_min: int
- salary_max: int
- job_type: string (full-time | part-time | contract | remote)
- published: bool (default: false)
- company: belongs_to
- category: belongs_to

### Category
- name: string (required, unique)
- slug: slug (unique)
- description: text

### Application
- cover_letter: text
- resume_url: string
- status: string (pending | reviewed | accepted | rejected)
- job: belongs_to
- user: belongs_to (the applicant)

## User Roles
- ADMIN: full access to all resources
- EMPLOYER: can manage their own company and jobs
- USER: can browse jobs and submit applications

## Features
- Public job listing page with search and filters
- Company profile pages
- Job detail pages with "Apply" button
- Application tracking for users
- Admin dashboard with stats and management`} />
                </div>

                {/* claude-code-prompt.md */}
                <div className="mb-8">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      2
                    </div>
                    <h3 className="text-lg font-semibold tracking-tight">
                      claude-code-prompt.md
                    </h3>
                  </div>
                  <p className="text-muted-foreground leading-relaxed mb-4 pl-10">
                    Step-by-step instructions that you feed directly to Claude Code. Contains
                    the exact <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit</code> CLI
                    commands to run, in the correct order. Claude Code reads this file and
                    executes each step automatically.
                  </p>

                  <CodeBlock filename="claude-code-prompt.md (Job Board example)" code={`# Claude Code — Build Instructions

Read GRIT_SKILL.md first to understand the framework.
Then read project-description.md for the full spec.
Execute these steps in order:

## Phase 1: Scaffold and Core Models

1. Scaffold the project:
   grit new jobboard
   cd jobboard

2. Start infrastructure:
   docker compose up -d

3. Add the EMPLOYER role:
   grit add role EMPLOYER

4. Generate core resources:
   grit generate resource Category --fields "name:string,slug:slug,description:text"
   grit generate resource Company --fields "name:string,slug:slug,logo:string,website:string,description:text,user_id:belongs_to"
   grit generate resource Job --fields "title:string,slug:slug,description:text,location:string,salary_min:int,salary_max:int,job_type:string,published:bool,company_id:belongs_to,category_id:belongs_to"
   grit generate resource Application --fields "cover_letter:text,resume_url:string,status:string,job_id:belongs_to,user_id:belongs_to"

5. Verify the build:
   cd apps/api && go build ./... && cd ../..

## Phase 2: Relationships and Business Logic

6. Update Go models to add Preload in services
7. Add a public endpoint for published jobs
8. Add role-based route guards for EMPLOYER routes
9. Run grit sync to update TypeScript types

## Phase 3: Frontend Customization

10. Customize admin resource definitions
11. Build the public job listing page
12. Build the company profile page
13. Add search and filter functionality`} />
                </div>

                {/* project-phases.md */}
                <div className="mb-8">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      3
                    </div>
                    <h3 className="text-lg font-semibold tracking-tight">
                      project-phases.md
                    </h3>
                  </div>
                  <p className="text-muted-foreground leading-relaxed mb-4 pl-10">
                    A phased checklist that breaks the build into manageable stages.
                    Phase 1 covers scaffolding and core models. Phase 2 handles
                    relationships and advanced features. Phase 3 is frontend customization
                    and polish. This prevents Claude Code from trying to do everything at once.
                  </p>

                  <CodeBlock filename="project-phases.md (Job Board example)" code={`# Job Board — Build Phases

## Phase 1: Scaffold + Core Models
- [ ] Run grit new jobboard
- [ ] Start Docker services
- [ ] Add EMPLOYER role with grit add role
- [ ] Generate Category resource
- [ ] Generate Company resource
- [ ] Generate Job resource
- [ ] Generate Application resource
- [ ] Verify go build ./... passes

## Phase 2: Relationships + Business Logic
- [ ] Add Preload("Company") and Preload("Category") to Job service
- [ ] Add Preload("Job") to Application service
- [ ] Create public /api/jobs/published endpoint
- [ ] Add EMPLOYER route guards for company/job management
- [ ] Add status workflow for applications (pending -> reviewed -> accepted/rejected)
- [ ] Run grit sync to update TypeScript types
- [ ] Verify go build ./... passes

## Phase 3: Frontend + Polish
- [ ] Customize Job admin resource (badges, filters, relation columns)
- [ ] Customize Application admin resource (status badges, job relation)
- [ ] Build public job listing page with search and category filter
- [ ] Build job detail page with Apply button
- [ ] Build company profile page
- [ ] Add application tracking page for logged-in users
- [ ] Test full flow: post job -> browse -> apply -> review`} />
                </div>
              </div>

              {/* ============================================================ */}
              {/* STEP 2: BUILDING WITH CLAUDE CODE */}
              {/* ============================================================ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Step 2: Building with Claude Code
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-6">
                  Once you have the 3 spec files from Claude Web, it&apos;s time to hand them
                  to Claude Code and let it build the application. Claude Code is an AI coding
                  agent that runs in your terminal &mdash; it can execute shell commands, read
                  and edit files, and build entire projects autonomously.
                </p>

                <div className="space-y-6">
                  {/* Install */}
                  <div className="relative pl-10">
                    <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      1
                    </div>
                    <div>
                      <h3 className="text-lg font-semibold tracking-tight mb-2">
                        Install Claude Code
                      </h3>
                      <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                        Install Claude Code globally using npm. You need Node.js 18+ installed.
                      </p>
                      <CodeBlock language="bash" filename="Terminal" code={`npm install -g @anthropic-ai/claude-code`} />
                    </div>
                  </div>

                  {/* Open terminal */}
                  <div className="relative pl-10">
                    <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      2
                    </div>
                    <div>
                      <h3 className="text-lg font-semibold tracking-tight mb-2">
                        Open your terminal in an empty directory
                      </h3>
                      <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                        Create a new folder for your project and navigate into it. Place the 3 spec
                        files from Claude Web into this directory.
                      </p>
                      <CodeBlock language="bash" filename="Terminal" code={`mkdir jobboard-project && cd jobboard-project

# Place your 3 spec files here:
# - project-description.md
# - claude-code-prompt.md
# - project-phases.md`} />
                    </div>
                  </div>

                  {/* Start Claude Code */}
                  <div className="relative pl-10">
                    <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      3
                    </div>
                    <div>
                      <h3 className="text-lg font-semibold tracking-tight mb-2">
                        Start a Claude Code session
                      </h3>
                      <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                        Launch Claude Code in the project directory. It will detect the spec files
                        and be ready to execute instructions.
                      </p>
                      <CodeBlock language="bash" filename="Terminal" code={`claude`} />
                    </div>
                  </div>

                  {/* Feed the prompt */}
                  <div className="relative pl-10">
                    <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      4
                    </div>
                    <div>
                      <h3 className="text-lg font-semibold tracking-tight mb-2">
                        Feed it the claude-code-prompt.md
                      </h3>
                      <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                        Tell Claude Code to read the prompt file and execute the instructions.
                        It will scaffold the project, generate every resource, and build the app step by step.
                      </p>
                      <CodeBlock filename="Message to Claude Code" code={`Read claude-code-prompt.md and project-description.md.
Execute all the steps in Phase 1 of the build instructions.
After each grit command, verify the output looks correct
before moving to the next step.`} />
                    </div>
                  </div>

                  {/* Review and iterate */}
                  <div className="relative pl-10">
                    <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      5
                    </div>
                    <div>
                      <h3 className="text-lg font-semibold tracking-tight mb-2">
                        Review and iterate
                      </h3>
                      <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                        After each phase, review what Claude Code built. Check that the Go API
                        compiles, models look correct, and the admin panel renders properly.
                        Then tell it to proceed with the next phase, or ask for adjustments.
                      </p>
                      <CodeBlock filename="Follow-up messages" code={`# After Phase 1 completes:
Phase 1 looks good. Now execute Phase 2 from the build instructions.

# If something needs fixing:
The Job model needs a "deadline" field (datetime). Add it and re-run grit sync.

# When ready for frontend work:
Execute Phase 3. Customize the admin resources and build the public pages.`} />
                    </div>
                  </div>
                </div>
              </div>

              {/* ============================================================ */}
              {/* THE GRIT SKILL FILE */}
              {/* ============================================================ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The Grit Skill File
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every Grit project generated with{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new</code>{' '}
                  includes a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">GRIT_SKILL.md</code> file
                  in the project root. This file teaches Claude Code (and any other AI assistant)
                  how your Grit project works &mdash; the folder structure, CLI commands, code
                  patterns, naming conventions, and injection markers.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When Claude Code opens your project directory, it automatically reads{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">GRIT_SKILL.md</code> and
                  gains full context about how to work with Grit. If you are starting from scratch
                  (before running <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new</code>),
                  you can create this file manually so Claude Code already understands Grit before
                  scaffolding begins.
                </p>

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mb-6">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-primary/90">Tip:</strong> Copy the skill file below and
                    save it as <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">GRIT_SKILL.md</code> in
                    your project directory <em>before</em> starting Claude Code. This ensures the AI
                    understands Grit conventions from the very first message.
                  </p>
                </div>

                <CodeBlock filename="GRIT_SKILL.md" code={`# Grit Framework — AI Skill Context

> Give this file to any AI coding assistant to teach it how to work with your Grit project.

## What is Grit?

Grit is a full-stack framework: Go (Gin + GORM) backend + Next.js (React + TypeScript) frontend in a monorepo.

## CLI Commands

grit new <name>                    # Scaffold full project
grit new <name> --api              # API only
grit generate resource <Name> --fields "field:type,..."  # Full-stack CRUD
grit generate resource <Name> --fields "..." --seed --faker --count 50  # + seeder
grit generate seeder <Resource>    # Add a seeder to an existing resource
grit remove resource <Name>        # Remove a resource
grit sync                          # Sync Go types → TypeScript
grit migrate                       # Run GORM AutoMigrate
grit seed                          # Run all seeders (idempotent)
grit start server                  # Run the Go API
grit start client                  # Run web + admin
grit add role <ROLE>               # Add a new role

## Field Types (exactly 15 — anything else fails validation)

string, text, richtext, int, uint, float, bool, datetime, date, slug
string_array (e.g., sizes:string_array)
belongs_to (e.g., category:belongs_to or author:belongs_to:User)
many_to_many (e.g., tags:many_to_many:Tag)
file  (e.g., image:file:image, doc:file:all)
files (e.g., gallery:files:image)

## Models

UUID string primary keys (NOT gorm.Model, NOT uint):
  ID string \`gorm:"type:varchar(36);primaryKey"\`  (set in BeforeCreate)
belongs_to FK columns are string (varchar(36)), e.g. CategoryID string.
User has FirstName + LastName; roles are UPPERCASE: ADMIN, EDITOR, USER (default USER).

## Project Structure

project/
├── apps/api/           # Go backend (Gin + GORM)
│   ├── cmd/server/     # Entry point
│   └── internal/       # models, handlers, services, middleware, routes
├── apps/web/           # Next.js frontend
├── apps/admin/         # Admin panel (resource-based)
├── packages/shared/    # Zod schemas, TS types, constants
└── docker-compose.yml  # PostgreSQL, Redis, MinIO, Mailhog

## Conventions

- Go models: PascalCase structs with gorm/json/binding tags
- API responses: { "data": ..., "meta": { "total", "page", "page_size", "pages" } }
- Errors: { "error": { "code": "...", "message": "..." } }
- TS files: kebab-case (use-users.ts, api-client.ts)
- Routes: plural, lowercase (/api/users, /api/posts)
- Validation: Zod schemas in packages/shared/schemas/
- Data fetching: React Query hooks (useQuery, useMutation)
- Styling: Tailwind CSS + shadcn/ui, dark theme default
- Admin: defineResource() API with DataTable, FormBuilder, widgets

## Code Generation Markers (DO NOT REMOVE — lowercase, no END markers)

// grit:models, // grit:routes:protected, // grit:routes:admin,
// grit:routes:custom, // grit:schemas, // grit:resources, // grit:seeders

## How to Build with Grit

1. grit new myapp && cd myapp && docker compose up -d
2. grit generate resource <Name> --fields "..." for each model
   (seed parents first, then children; add --seed --faker to generate data)
3. grit migrate && grit seed
4. Terminal 1: grit start server   |   Terminal 2: grit start client
5. Customize handlers, add business logic, build frontend pages`} />
              </div>

              {/* ============================================================ */}
              {/* TIPS & BEST PRACTICES */}
              {/* ============================================================ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Tips &amp; Best Practices
                </h2>
                <div className="space-y-4">
                  {[
                    {
                      title: 'Let Claude Code run grit commands',
                      desc: 'Do not manually type grit commands in a separate terminal. Let Claude Code execute them so it can read the output, detect errors, and fix issues automatically. The whole point of the workflow is hands-free building.',
                    },
                    {
                      title: 'Review generated models before adding relationships',
                      desc: 'After generating core resources, check that the Go models have the correct field types and tags before adding belongs_to or many_to_many relationships. It is easier to regenerate a resource early than to fix relationships later.',
                    },
                    {
                      title: 'Use phases — do not try to build everything at once',
                      desc: 'The phased approach prevents Claude Code from getting overwhelmed. Complete Phase 1 (scaffold + core models), verify it compiles, then move to Phase 2 (relationships), and so on. Each phase builds on a working foundation.',
                    },
                    {
                      title: 'Always check go build ./... after each phase',
                      desc: 'This catches compilation errors early. If Claude Code introduces a bug in Phase 1, you want to find it before Phase 2 adds more complexity on top. Include this step in your claude-code-prompt.md.',
                    },
                    {
                      title: 'Use grit remove resource if you need to redo something',
                      desc: 'If a resource was generated with the wrong fields, remove it completely and regenerate. This is cleaner than manually editing multiple files across the monorepo. The CLI handles cleanup in all layers.',
                    },
                    {
                      title: 'Keep GRIT_SKILL.md in your project root',
                      desc: 'Claude Code automatically reads files in the project root. Having the skill file there means every new Claude Code session starts with full Grit context, no manual pasting required.',
                    },
                  ].map((item) => (
                    <div key={item.title} className="p-4 rounded-lg border border-border/30 bg-card/30">
                      <h3 className="text-sm font-semibold mb-1.5">{item.title}</h3>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.desc}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* ============================================================ */}
              {/* EXAMPLE PROMPT FOR CLAUDE WEB */}
              {/* ============================================================ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Complete Prompt Template for Claude Web
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Here is the full, ready-to-copy prompt that you can paste directly into{' '}
                  <Link href="https://claude.ai" target="_blank" className="text-primary hover:underline">
                    claude.ai
                  </Link>.
                  Replace the bracketed sections with your own project details and Claude will
                  generate all 3 spec files in a single response.
                </p>

                <CodeBlock filename="Copy this entire prompt into claude.ai" code={`I want to build a [PROJECT NAME] using the Grit framework.

Grit is a full-stack meta-framework that uses:
- Go backend (Gin web framework + GORM ORM)
- Next.js frontend (React + App Router + TypeScript)
- Admin panel (resource-based, Filament-like)
- Monorepo with shared types and validation (Zod)
- Docker for infrastructure (PostgreSQL, Redis, MinIO)

Here's what I want to build:
[DESCRIBE YOUR PROJECT IN DETAIL]

Please create these 3 files for me:

1. project-description.md — Full project spec with all models, their
   fields (use Grit field types: string, text, richtext, int, uint,
   float, bool, datetime, date, slug, string_array, belongs_to,
   many_to_many, file, files), relationships, user roles, and features.

2. claude-code-prompt.md — Step-by-step instructions that I'll give
   to Claude Code (AI coding agent). Include exact grit CLI commands
   to run in order. Use \`grit generate resource <Name> --fields "..."\`
   syntax.

3. project-phases.md — A phased build plan with checkboxes. Phase 1:
   scaffold + core models. Phase 2: relationships + complex features.
   Phase 3: frontend customization + polish.`} />

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mb-6">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-primary/90">Tip:</strong> The more detail you provide in
                    the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">[DESCRIBE YOUR PROJECT]</code> section,
                    the better your spec files will be. Include specific field names and types,
                    relationships between models, user roles and permissions, and any custom
                    business logic or API endpoints you need.
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
                  <Link href="/docs/ai-workflows/antigravity" className="gap-1.5">
                    Using Grit with Antigravity
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
