import Link from 'next/link'
import { ArrowLeft, ArrowRight, Download } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/ai-skill')

export default function AISkillPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">For AI Assistants</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                LLM Skill Guide
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed mb-6">
                Teach AI assistants to understand your Grit project. Drop a single file
                into your project root and your AI can generate resources, write handlers,
                and follow Grit conventions perfectly.
              </p>
              <a
                href="https://14j7oh8kso.ufs.sh/f/HLxTbDBCDLwfeHHJl34ZKSqNhOvVj6p9rg3Icmo05TAEwQ4a"
                target="_blank"
                rel="noopener noreferrer"
              >
                <Button size="sm" variant="outline" className="gap-1.5 border-primary/20 text-primary/80 hover:bg-primary/10 hover:text-primary">
                  <Download className="h-3.5 w-3.5" />
                  Download Grit Handbook PDF
                </Button>
              </a>
            </div>

            <div className="prose-grit">
              {/* What is GRIT_SKILL.md */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  What is GRIT_SKILL.md?
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">GRIT_SKILL.md</code> is
                  a structured context file that teaches AI coding assistants how your Grit project works.
                  It contains everything an AI needs to write code that follows Grit&apos;s conventions:
                  project structure, CLI commands, code patterns, API format, naming rules, and more.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Think of it as a &quot;skill file&quot; &mdash; a document that gives any AI assistant
                  instant expertise in Grit development. Instead of explaining your project setup in every
                  conversation, the AI reads the skill file and immediately understands the codebase.
                </p>
                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-primary/90">Why this matters:</strong> Grit&apos;s strict conventions
                    and predictable structure make it uniquely suited for AI-assisted development. An AI
                    that understands Grit&apos;s patterns can generate complete, working code &mdash; not just
                    snippets that need manual wiring.
                  </p>
                </div>
              </div>

              {/* How to use it */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  How to Use It
                </h2>

                <div className="mb-8">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      1
                    </div>
                    <h3 className="text-lg font-semibold tracking-tight">
                      Add the file to your project root
                    </h3>
                  </div>
                  <p className="text-muted-foreground leading-relaxed mb-4 pl-10">
                    Every Grit project generated with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new</code> includes
                    a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">GRIT_SKILL.md</code> file
                    in the project root automatically. If you don&apos;t have one, create it manually:
                  </p>
                  <CodeBlock filename="Project structure" code={`myapp/
├── GRIT_SKILL.md          <-- AI reads this
├── grit.config.ts
├── docker-compose.yml
├── packages/shared/
├── apps/
│   ├── api/
│   ├── web/
│   └── admin/
└── ...`} />
                </div>

                <div className="mb-8">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      2
                    </div>
                    <h3 className="text-lg font-semibold tracking-tight">
                      Open your project with an AI assistant
                    </h3>
                  </div>
                  <p className="text-muted-foreground leading-relaxed mb-4 pl-10">
                    The following AI tools will automatically detect and use the skill file:
                  </p>
                  <div className="ml-10 space-y-3">
                    {[
                      {
                        name: 'Claude Code (CLI)',
                        desc: 'Reads GRIT_SKILL.md as project context automatically when you open the project directory.',
                      },
                      {
                        name: 'Cursor',
                        desc: 'Add GRIT_SKILL.md to your project. Cursor indexes it and uses the context when generating code.',
                      },
                      {
                        name: 'Windsurf',
                        desc: 'The file is indexed as part of your workspace context and informs all code generation.',
                      },
                      {
                        name: 'GitHub Copilot',
                        desc: 'When working in VS Code, reference @workspace to include GRIT_SKILL.md context in chat.',
                      },
                      {
                        name: 'Any AI Chat (ChatGPT, Claude)',
                        desc: 'Copy and paste the GRIT_SKILL.md content into the conversation as system context.',
                      },
                    ].map((tool) => (
                      <div key={tool.name} className="p-3 rounded-lg border border-border/30 bg-card/30">
                        <h4 className="text-sm font-semibold mb-1">{tool.name}</h4>
                        <p className="text-xs text-muted-foreground/70 leading-relaxed">{tool.desc}</p>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="mb-8">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                      3
                    </div>
                    <h3 className="text-lg font-semibold tracking-tight">
                      Start building with AI
                    </h3>
                  </div>
                  <p className="text-muted-foreground leading-relaxed pl-10">
                    The AI now understands your project structure, conventions, and patterns.
                    Ask it to generate resources, add fields, create endpoints, or build features &mdash;
                    and it will follow Grit conventions automatically.
                  </p>
                </div>
              </div>

              {/* What it teaches */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  What It Teaches AI Assistants
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The GRIT_SKILL.md file covers everything an AI needs to write
                  idiomatic Grit code:
                </p>
                <div className="space-y-4">
                  {[
                    {
                      title: 'Project Structure',
                      desc: 'The complete folder layout: where models, handlers, services, and middleware live in the Go API; where pages, components, and hooks live in Next.js; how the shared package connects them.',
                      example: 'apps/api/internal/models/ — Go models\napps/api/internal/handlers/ — HTTP handlers\napps/web/app/ — Next.js pages\npackages/shared/src/ — Zod schemas + types',
                    },
                    {
                      title: 'CLI Commands',
                      desc: 'All available Grit CLI commands: grit new, grit generate resource, grit generate seeder, grit start (server / client), grit migrate, grit seed, grit sync, grit studio, grit add role, grit upgrade. There is no grit dev command. The AI knows which commands to suggest.',
                      example: 'grit generate resource Invoice --fields "..." --seed\n# Creates model, handler, service, hook, schema, admin page + seeder\n\ngrit generate seeder Invoice --faker --count 50\n# Fills 50 realistic rows (belongs_to linked to real parents)',
                    },
                    {
                      title: 'Code Patterns',
                      desc: 'Standard patterns for handlers (thin controllers that call services), services (business logic), models (GORM structs with tags), and hooks (React Query wrappers).',
                      example: 'Handler -> Service -> Database (GORM)\nComponent -> React Query Hook -> API Client',
                    },
                    {
                      title: 'API Conventions',
                      desc: 'The standard JSON response format: { data, message } for success, { error: { code, message, details } } for errors. HTTP status codes. Pagination format with meta.',
                      example: '{ "data": [...], "meta": { "total": 100, "page": 1, "page_size": 20 } }',
                    },
                    {
                      title: 'Markers (Do Not Delete!)',
                      desc: 'Lowercase code markers like // grit:models, // grit:routes:protected, and // grit:seeders that the CLI uses to inject generated code. There are no uppercase or END markers. The AI knows to preserve these and insert new code just before them.',
                      example: '&models.User{},\n&models.Invoice{}, // new models go above\n// grit:models',
                    },
                    {
                      title: 'Batteries & Dashboards',
                      desc: 'All built-in services: file storage (S3/R2), email (Resend), background jobs (asynq), cron, Redis caching, AI (Vercel AI Gateway — one AI_GATEWAY_API_KEY, many models), security (Sentinel), observability (Pulse), and API docs (gin-docs).',
                      example: '/studio   → GORM Studio (DB browser)\n/docs     → gin-docs (API docs)\n/sentinel → Security dashboard\n/pulse    → Observability dashboard',
                    },
                    {
                      title: 'Naming Conventions',
                      desc: 'Go files: snake_case. Go structs: PascalCase. TypeScript files: kebab-case. React components: PascalCase. API routes: plural lowercase. Database tables: plural snake_case.',
                      example: 'user_handler.go, CreatePostSchema, use-users.ts, /api/users',
                    },
                  ].map((item) => (
                    <div key={item.title} className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                      <div className="p-4">
                        <h3 className="text-sm font-semibold mb-1.5">{item.title}</h3>
                        <p className="text-xs text-muted-foreground/70 leading-relaxed mb-3">{item.desc}</p>
                        <div className="rounded-lg border border-border/30 bg-accent/20 overflow-hidden">
                          <pre className="px-4 py-2.5 font-mono text-[11px] text-foreground/60 overflow-x-auto">{item.example}</pre>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Vibe Coding with Grit */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  &quot;Vibe Coding&quot; with Grit
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit&apos;s strict conventions make it ideal for AI-assisted development &mdash; what some
                  call &quot;vibe coding.&quot; Because every Grit project follows the same patterns,
                  AI assistants can generate correct, production-ready code with minimal guidance.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  What You Can Ask Your AI
                </h3>
                <div className="space-y-3 mb-6">
                  {[
                    {
                      prompt: '"Add a Product resource with name, price, description, and category fields"',
                      result: 'The AI creates the Go model, handler, service, Zod schema, TypeScript types, React Query hook, and admin page — all following Grit conventions.',
                    },
                    {
                      prompt: '"Add an image upload field to the Product model"',
                      result: 'The AI adds the field to the Go struct, updates the Zod schema, adds the upload UI to the admin form, and wires the storage service.',
                    },
                    {
                      prompt: '"Create an API endpoint that returns monthly revenue stats"',
                      result: 'The AI adds the handler, service function with the GORM query, route registration, and React Query hook.',
                    },
                    {
                      prompt: '"Add a sidebar link for Products in the admin panel"',
                      result: 'The AI updates the sidebar component with the correct path, icon, and positioning.',
                    },
                    {
                      prompt: '"Add email verification to the registration flow"',
                      result: 'The AI adds the verification token to the User model, creates the email template, updates the auth handler, and adds the verification page.',
                    },
                  ].map((item) => (
                    <div key={item.prompt} className="p-4 rounded-lg border border-border/30 bg-card/30">
                      <p className="text-sm font-medium text-foreground/90 mb-1.5">{item.prompt}</p>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.result}</p>
                    </div>
                  ))}
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Why Grit Works Well with AI
                </h3>
                <ul className="space-y-2.5">
                  {[
                    'Predictable structure — every file has exactly one place it belongs',
                    'Consistent patterns — handlers, services, and hooks all follow the same shape',
                    'Code markers — the AI knows exactly where to inject new code',
                    'Shared types — changes to Go structs can be propagated to TypeScript automatically',
                    'Convention over configuration — fewer decisions means fewer mistakes',
                    'Single framework — the AI does not need to guess which library you are using',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-1">&#10003;</span>
                      {item}
                    </li>
                  ))}
                </ul>
              </div>

              {/* Tips for best results */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Tips for Best Results
                </h2>
                <div className="space-y-4">
                  {[
                    {
                      title: 'Be specific about field types',
                      desc: 'Instead of "add a products table," say "add a Product resource with name (string), price (float), description (text), category (belongs_to Category), and is_active (bool)." The AI generates better code with precise Grit field types.',
                    },
                    {
                      title: 'Reference existing patterns',
                      desc: 'Say "follow the same pattern as the User handler" or "use the same form layout as the Users admin page." The AI understands Grit patterns and will replicate them consistently.',
                    },
                    {
                      title: 'Ask for full-stack changes',
                      desc: 'Grit\'s monorepo structure means the AI can update the Go model, TypeScript types, and React UI in a single conversation. Ask for the complete change rather than one layer at a time.',
                    },
                    {
                      title: 'Never delete markers',
                      desc: 'Comments like // grit:models and // grit:routes:protected are used by the CLI code generator. If the AI tries to remove them, ask it to keep them. The skill file already instructs AI assistants to preserve markers.',
                    },
                    {
                      title: 'Keep the skill file updated',
                      desc: 'When you add custom patterns or change conventions, update GRIT_SKILL.md so the AI stays in sync. The file is your single source of truth for AI context.',
                    },
                  ].map((item) => (
                    <div key={item.title} className="p-4 rounded-lg border border-border/30 bg-card/30">
                      <h3 className="text-sm font-semibold mb-1.5">{item.title}</h3>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.desc}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Get the Skill File */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Get the Skill File
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">GRIT_SKILL.md</code> file
                  is generated automatically when you create a new Grit project:
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
                  <div className="p-5 font-mono text-sm space-y-2">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">grit new myapp</span>
                    </div>
                    <div className="text-muted-foreground/50 text-xs mt-2">
                      # GRIT_SKILL.md is created in the project root
                    </div>
                  </div>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  You can also find the latest version on the{' '}
                  <Link href="https://github.com/MUKE-coder/grit" target="_blank" className="text-primary hover:underline">
                    Grit GitHub repository
                  </Link>.
                  Copy it into any existing Grit project that does not have one.
                </p>
                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-primary/90">Remember:</strong> The skill file works best when it
                    reflects the current state of your project. After generating resources with{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>,
                    the CLI automatically updates the skill file with the new models and routes.
                  </p>
                </div>
              </div>

              {/* Nav */}
              <div className="flex items-center justify-between pt-6 border-t border-border/30">
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/tutorials/ecommerce" className="gap-1.5">
                    <ArrowLeft className="h-3.5 w-3.5" />
                    Tutorials
                  </Link>
                </Button>
                <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                  <Link href="/docs/ai-workflows/claude" className="gap-1.5">
                    Using Grit with Claude
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
