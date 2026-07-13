import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { PageHelp } from '@/components/page-help'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/getting-started/coming-from')

interface CmdRow {
  task: string
  laravel: string
  django: string
  grit: string
}

const COMMANDS: CmdRow[] = [
  { task: 'Create a project', laravel: 'laravel new blog', django: 'django-admin startproject', grit: 'grit new blog' },
  { task: 'Model + CRUD + API', laravel: 'make:model -mcr', django: 'startapp + models.py', grit: 'grit generate resource' },
  { task: 'Run migrations', laravel: 'artisan migrate', django: 'manage.py migrate', grit: 'grit migrate' },
  { task: 'Seed the database', laravel: 'artisan db:seed', django: 'loaddata (fixtures)', grit: 'grit seed' },
  { task: 'Run the dev server', laravel: 'artisan serve', django: 'manage.py runserver', grit: 'grit start' },
  { task: 'Browse the DB', laravel: 'artisan tinker', django: 'manage.py shell', grit: '/studio (GORM Studio)' },
  { task: 'Admin panel', laravel: 'Filament', django: 'django.contrib.admin', grit: 'generated per resource' },
]

const CONCEPTS: { grit: string; is: string }[] = [
  { grit: 'Resource', is: 'One declaration → model, migration, CRUD API, validation, admin page, and typed hooks. Laravel’s make:model -mcr + a Filament resource + your API + your TS client, in a single command.' },
  { grit: 'GORM', is: 'The ORM — the Eloquent / Django-ORM of the Go world. Models are Go structs with tags instead of PHP/Python classes.' },
  { grit: 'grit migrate', is: 'Explicit, like artisan/manage.py migrate. Uses GORM AutoMigrate (adds tables & columns); there are no Laravel-style down-migrations — reset in dev with grit migrate --fresh.' },
  { grit: 'Shared types', is: 'Go struct → TypeScript types + Zod, kept in sync by grit sync. There is no Laravel/Django equivalent — your frontend is typed against your backend for free, so you never hand-write an API client.' },
  { grit: 'The admin panel', is: 'Filament-style and resource-driven, but it is generated code you own — not a black box. Style it, override it, or ignore it.' },
]

const NEXT_DOCS: { title: string; desc: string; href: string }[] = [
  { title: 'Models & Database', desc: 'GORM structs, tags, relations — your Eloquent/Django models', href: '/docs/backend/models' },
  { title: 'Handlers', desc: 'Thin Gin handlers — the controllers of a Grit API', href: '/docs/backend/handlers' },
  { title: 'Services', desc: 'Where business logic lives, called by handlers', href: '/docs/backend/services' },
  { title: 'Migrations', desc: 'grit migrate, AutoMigrate, and --fresh in dev', href: '/docs/backend/migrations' },
  { title: 'Architecture Overview', desc: 'How the handler → service → model split fits together', href: '/docs/concepts/architecture' },
  { title: 'Database & GORM Studio', desc: 'Connections, DSNs, and the visual row browser', href: '/docs/infrastructure/database' },
]

export default function ComingFromPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Getting Started</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Coming from Laravel, Django, or Next.js
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                You already think in models, migrations, seeders, an admin panel, and a CLI that
                scaffolds everything. Grit works the same way &mdash; a Go backend and a React
                frontend, driven by one <code>grit</code> command. Here&apos;s the translation so
                you feel at home in minutes.
              </p>
            </div>

            <div className="prose-grit">
              {/* Command mapping */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Command cheat sheet</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Your muscle memory maps almost one-to-one. The big difference: Grit&apos;s{' '}
                  <code>generate resource</code> does in one command what Laravel splits across{' '}
                  <code>make:model</code>, <code>make:controller</code>, <code>make:migration</code>,
                  and a Filament resource.
                </p>
                <div className="overflow-x-auto">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-3 py-2 font-medium">Task</th>
                        <th className="px-3 py-2 font-medium">Laravel</th>
                        <th className="px-3 py-2 font-medium">Django</th>
                        <th className="px-3 py-2 font-medium text-primary">Grit</th>
                      </tr>
                    </thead>
                    <tbody>
                      {COMMANDS.map((r) => (
                        <tr key={r.task} className="border-b border-border/50 align-top">
                          <td className="px-3 py-2 text-muted-foreground">{r.task}</td>
                          <td className="px-3 py-2 font-mono text-[12px] text-muted-foreground/70">{r.laravel}</td>
                          <td className="px-3 py-2 font-mono text-[12px] text-muted-foreground/70">{r.django}</td>
                          <td className="px-3 py-2 font-mono text-[12px] text-primary">{r.grit}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* The whole app in grit commands */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">The whole workflow, in grit</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Like <code>artisan</code> or <code>manage.py</code>, the Grit CLI is the only
                  interface you need &mdash; you never <code>cd</code> into a sub-folder or run raw{' '}
                  <code>go</code>/<code>pnpm</code> to work on your app.
                </p>
                <CodeBlock
                  terminal
                  code={`grit new blog                 # scaffold the project
cd blog
docker compose up -d          # start Postgres, Redis, MinIO, Mailhog
pnpm install                  # frontend deps (one-time)

grit generate resource Post --fields "title:string,body:text,published:bool"
grit migrate                  # create the tables
grit seed                     # (optional) sample data
grit start                    # run the API + web + admin together`}
                />
              </div>

              {/* Concept mapping */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Concepts that translate</h2>
                <div className="space-y-4">
                  {CONCEPTS.map((c) => (
                    <div key={c.grit} className="rounded-lg border border-border/40 bg-card/30 p-4">
                      <div className="font-mono text-[13px] text-primary mb-1">{c.grit}</div>
                      <p className="text-sm text-muted-foreground leading-relaxed">{c.is}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Supported databases */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Supported databases</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit talks to your database through <strong>GORM</strong>, so the supported engines
                  are <strong>PostgreSQL</strong> (production) and <strong>SQLite</strong> (dev,
                  desktop, and tests). The Go API builds its connection string from the{' '}
                  <code>POSTGRES_*</code> environment variables, or you can hand it a single{' '}
                  <code>DATABASE_URL</code> and it will use that instead. Desktop apps default to a
                  local SQLite file so they run with zero infrastructure.
                </p>
                <CodeBlock
                  language="bash"
                  code={`# Option A — the API assembles the DSN from these
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=grit
POSTGRES_PASSWORD=grit
POSTGRES_DB=blog

# Option B — one URL wins over the POSTGRES_* vars
DATABASE_URL=postgres://grit:grit@localhost:5432/blog
DATABASE_URL=sqlite:./app.db        # local file (desktop default)
DATABASE_URL=sqlite::memory:        # ephemeral, used by tests`}
                />
              </div>

              {/* Docker */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Docker (what it runs, and why)</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  <code>docker compose up -d</code> spins up the <em>infrastructure</em> your app
                  talks to &mdash; not your app. Think Laravel Sail: one command gives you a local
                  Postgres, Redis, MinIO and Mailhog. You still run the app itself with{' '}
                  <code>grit start</code>.
                </p>
                <div className="overflow-x-auto mb-4">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-3 py-2 font-medium">Service</th>
                        <th className="px-3 py-2 font-medium">What it is</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 font-mono text-[12px] text-primary">postgres</td>
                        <td className="px-3 py-2 text-muted-foreground">Your primary database in dev.</td>
                      </tr>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 font-mono text-[12px] text-primary">redis</td>
                        <td className="px-3 py-2 text-muted-foreground">Cache and background-job queue (asynq).</td>
                      </tr>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 font-mono text-[12px] text-primary">minio</td>
                        <td className="px-3 py-2 text-muted-foreground">S3-compatible object storage for uploads.</td>
                      </tr>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 font-mono text-[12px] text-primary">mailhog</td>
                        <td className="px-3 py-2 text-muted-foreground">Catches outgoing email so you can preview it locally.</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <p className="text-muted-foreground leading-relaxed">
                  Docker is optional. Point the same env vars at hosted services &mdash; Neon or
                  Supabase for Postgres, Upstash for Redis, Cloudflare R2 for storage &mdash; and you
                  can develop and deploy completely Docker-free.
                </p>
              </div>

              {/* Architecture: MVC vs Grit */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Architecture: MVC vs Grit</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Laravel and Django are <strong>MVC</strong> &mdash; Model, View, Controller in one
                  app. Grit keeps the same responsibilities but splits them cleanly: on the backend a
                  thin <strong>Gin handler</strong> receives the request and calls a{' '}
                  <strong>service</strong> that owns the business logic, which in turn works with{' '}
                  <strong>GORM models</strong>. The &ldquo;view&rdquo; is a separate React frontend,
                  fed by typed hooks generated from those same models.
                </p>
                <div className="overflow-x-auto">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-3 py-2 font-medium">Laravel / Django</th>
                        <th className="px-3 py-2 font-medium text-primary">Grit equivalent</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 text-muted-foreground">Model</td>
                        <td className="px-3 py-2 text-primary">GORM model (Go struct with tags)</td>
                      </tr>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 text-muted-foreground">Controller</td>
                        <td className="px-3 py-2 text-primary">Handler (thin Gin handler)</td>
                      </tr>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 text-muted-foreground">Fat controllers / business logic</td>
                        <td className="px-3 py-2 text-primary">Service layer</td>
                      </tr>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 text-muted-foreground">Blade / template / View</td>
                        <td className="px-3 py-2 text-primary">React frontend + shared generated types</td>
                      </tr>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 text-muted-foreground">Route</td>
                        <td className="px-3 py-2 text-primary">Gin route + <code>grit:routes</code> markers</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Monorepos */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Monorepos</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The default Grit modes (<code>triple</code>, <code>double</code>, and{' '}
                  <code>mobile</code>) scaffold a <strong>Turborepo monorepo</strong> &mdash; one repo
                  holding every part of your product, with the Go&rarr;TypeScript types shared through
                  a single package.
                </p>
                <div className="overflow-x-auto mb-4">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-3 py-2 font-medium">Path</th>
                        <th className="px-3 py-2 font-medium">What it holds</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 font-mono text-[12px] text-primary">apps/api</td>
                        <td className="px-3 py-2 text-muted-foreground">The Go backend (Gin + GORM).</td>
                      </tr>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 font-mono text-[12px] text-primary">apps/web</td>
                        <td className="px-3 py-2 text-muted-foreground">The Next.js user-facing frontend.</td>
                      </tr>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 font-mono text-[12px] text-primary">apps/admin</td>
                        <td className="px-3 py-2 text-muted-foreground">The generated admin panel.</td>
                      </tr>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 font-mono text-[12px] text-primary">packages/shared</td>
                        <td className="px-3 py-2 text-muted-foreground">The Go&rarr;TS types &amp; Zod schemas that keep everything in sync.</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <p className="text-muted-foreground leading-relaxed">
                  The <code>single</code> and <code>api-only</code> modes are flat instead of a
                  monorepo. Either way it&apos;s one repo, shared types, and a single{' '}
                  <code>grit start</code>.
                </p>
              </div>

              {/* Prisma vs GORM */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Prisma vs GORM</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Coming from a Node/Prisma world, the mental model carries over &mdash; only the
                  tooling changes. Grit&apos;s ORM is <strong>GORM</strong> (Go).
                </p>
                <div className="overflow-x-auto">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-3 py-2 font-medium">Prisma</th>
                        <th className="px-3 py-2 font-medium text-primary">Grit / GORM</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 text-muted-foreground"><code>schema.prisma</code> DSL</td>
                        <td className="px-3 py-2 text-primary">Go structs with <code>gorm</code>/<code>json</code> tags</td>
                      </tr>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 text-muted-foreground"><code>prisma migrate</code></td>
                        <td className="px-3 py-2 text-primary"><code>grit migrate</code> (GORM AutoMigrate)</td>
                      </tr>
                      <tr className="border-b border-border/50 align-top">
                        <td className="px-3 py-2 text-muted-foreground">Generated Prisma Client types</td>
                        <td className="px-3 py-2 text-primary"><code>grit sync</code> &rarr; TypeScript from your Go structs</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <p className="text-muted-foreground leading-relaxed mt-4">
                  So the end-to-end type-safety you loved from Prisma Client is still there &mdash; it
                  just flows from your Go models out to the frontend instead of from a schema file.
                </p>
              </div>

              {/* GORM Studio */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">GORM Studio</h2>
                <p className="text-muted-foreground leading-relaxed">
                  Grit ships a visual database browser &mdash; the equivalent of Prisma Studio,
                  Adminer, or Django admin&apos;s data view. It&apos;s served by the Go API at{' '}
                  <code>http://localhost:8080/studio</code>, and lets you browse and edit rows
                  directly while developing, no SQL required.
                </p>
              </div>

              {/* Where to go next */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Where to go next</h2>
                <div className="grid gap-3 sm:grid-cols-2">
                  {NEXT_DOCS.map((item) => (
                    <Link key={item.href} href={item.href}>
                      <div className="group rounded-lg border border-border/40 bg-card/50 p-4 hover:border-primary/20 hover:bg-card/80 transition-all duration-200">
                        <h3 className="text-[15px] font-semibold mb-1 group-hover:text-primary transition-colors flex items-center gap-1.5">
                          {item.title}
                          <ArrowRight className="h-3.5 w-3.5 opacity-0 -translate-x-1 group-hover:opacity-100 group-hover:translate-x-0 transition-all" />
                        </h3>
                        <p className="text-sm text-muted-foreground/60">{item.desc}</p>
                      </div>
                    </Link>
                  ))}
                </div>
              </div>

              <Callout type="tip" title="The one real mindset shift">
                In a Laravel/Django + Next.js setup you maintain two codebases and hand-write the
                API client between them. In Grit the Go backend and React frontend are one project
                with <strong>shared, generated types</strong> &mdash; change a Go struct, run{' '}
                <code>grit sync</code>, and the frontend knows. That single flow is the thing worth
                unlearning your old habits for.
              </Callout>

              <Callout type="note" title="What Grit does differently">
                It&apos;s Go, not PHP/Python &mdash; a compiled, single-binary backend. Migrations use
                GORM AutoMigrate (no down-migrations; use <code>grit migrate --fresh</code> in dev).
                And the frontend is React (Next.js or Vite), not Blade/Livewire or Django templates.
                If you know Go and can read React/TypeScript, everything else is familiar.
              </Callout>

              <div className="flex items-center justify-between border-t border-border pt-8 mt-12">
                <Button variant="ghost" asChild>
                  <Link href="/docs/getting-started/quick-start" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Quick Start
                  </Link>
                </Button>
                <Button variant="ghost" asChild>
                  <Link href="/docs/tutorials/contact-app" className="gap-2">
                    Build your first app
                    <ArrowRight className="h-4 w-4" />
                  </Link>
                </Button>
              </div>

              <PageHelp
                faqs={[
                  {
                    q: 'Do I have to learn Go?',
                    a: (
                      <>
                        Enough to read and edit structs, handlers, and services &mdash; but{' '}
                        <code>grit generate resource</code> writes the boilerplate for you. If you
                        know one C-family language you&apos;ll be productive fast. See the{' '}
                        <Link href="/docs/prerequisites/golang">Go for Grit Developers</Link> primer.
                      </>
                    ),
                  },
                  {
                    q: 'Can I use MySQL?',
                    a: (
                      <>
                        Grit officially supports <strong>PostgreSQL</strong> (production) and{' '}
                        <strong>SQLite</strong> (dev, desktop, tests). GORM can talk to MySQL, but the
                        scaffolds, DSN builder, and defaults are tuned for Postgres and SQLite &mdash;
                        those are the supported paths.
                      </>
                    ),
                  },
                  {
                    q: 'Is there an Eloquent / Prisma-style query builder?',
                    a: (
                      <>
                        Yes &mdash; that&apos;s GORM. You get a chainable API like{' '}
                        <code>db.Where(&quot;published = ?&quot;, true).Find(&amp;posts)</code>, plus
                        associations, preloading, and scopes. It&apos;s the Go equivalent of Eloquent
                        or the Prisma Client.
                      </>
                    ),
                  },
                  {
                    q: 'Where do migrations live?',
                    a: (
                      <>
                        Grit uses GORM AutoMigrate driven by your models, run with{' '}
                        <code>grit migrate</code>. There are no hand-written up/down migration files
                        like Laravel or Django &mdash; in dev, reset with{' '}
                        <code>grit migrate --fresh</code>. See{' '}
                        <Link href="/docs/backend/migrations">Migrations</Link>.
                      </>
                    ),
                  },
                  {
                    q: 'Can I keep my existing Postgres database?',
                    a: (
                      <>
                        Yes. Point <code>DATABASE_URL</code> (or the <code>POSTGRES_*</code> vars) at
                        your existing database and skip the Docker Postgres. AutoMigrate adds new
                        tables and columns without dropping your data.
                      </>
                    ),
                  },
                ]}
              />
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
