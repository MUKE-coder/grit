import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
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
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
