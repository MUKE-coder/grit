import Link from 'next/link'
import { ArrowLeft, ArrowRight, Server, Smartphone, Monitor, Layers } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { Tabs } from '@/components/tabs'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/getting-started/create-a-project')

function WhatYouGet({ rows }: { rows: [string, string][] }) {
  return (
    <div className="mt-4 grid gap-2 sm:grid-cols-2">
      {rows.map(([name, url]) => (
        <div key={name} className="rounded-lg border border-border/40 bg-card/40 px-4 py-2.5">
          <div className="text-[13px] font-semibold">{name}</div>
          <code className="text-xs font-mono text-primary/70">{url}</code>
        </div>
      ))}
    </div>
  )
}

export default function CreateAProjectPage() {
  const tabs = [
    {
      id: 'all',
      label: 'Everything',
      icon: <Layers className="h-4 w-4" />,
      content: (
        <div>
          <p className="text-muted-foreground leading-relaxed mb-4">
            The full stack &mdash; Go API, a Next.js web app, and the Filament-style admin panel.
            The best place to start.
          </p>
          <CodeBlock
            terminal
            code={`grit new myapp --triple    # Go API + web + admin
cd myapp
docker compose up -d       # Postgres, Redis, MinIO, Mailhog
pnpm install               # frontend deps (one-time)
grit migrate               # create database tables
grit seed                  # sample data + a demo admin login
grit start                 # run all three, Ctrl+C stops them`}
          />
          <WhatYouGet
            rows={[
              ['Web app', 'http://localhost:3000'],
              ['Admin panel', 'http://localhost:3001'],
              ['Go API + docs', 'http://localhost:8080/docs'],
              ['GORM Studio', 'http://localhost:8080/studio'],
            ]}
          />
        </div>
      ),
    },
    {
      id: 'api',
      label: 'API only',
      icon: <Server className="h-4 w-4" />,
      content: (
        <div>
          <p className="text-muted-foreground leading-relaxed mb-4">
            A headless Go API &mdash; no frontend. Perfect for a mobile/SPA backend or a
            microservice. Ships with auth, storage, jobs, and interactive API docs.
          </p>
          <CodeBlock
            terminal
            code={`grit new myapp --api       # headless Go API
cd myapp
docker compose up -d       # Postgres, Redis, MinIO, Mailhog
grit migrate               # create database tables
grit seed                  # sample data + a demo admin login
grit start server          # run the API`}
          />
          <WhatYouGet
            rows={[
              ['API', 'http://localhost:8080'],
              ['Interactive docs', 'http://localhost:8080/docs'],
              ['GORM Studio', 'http://localhost:8080/studio'],
            ]}
          />
        </div>
      ),
    },
    {
      id: 'mobile',
      label: 'Mobile',
      icon: <Smartphone className="h-4 w-4" />,
      content: (
        <div>
          <p className="text-muted-foreground leading-relaxed mb-4">
            A Go API plus an Expo (React Native) app that share types. Generate a resource and
            you get typed screens on your phone.
          </p>
          <CodeBlock
            terminal
            code={`grit new myapp --mobile    # Go API + Expo app
cd myapp
docker compose up -d
pnpm install
grit migrate
grit seed
grit start server          # terminal 1 — the API
grit start expo            # terminal 2 — the Expo dev server (scan the QR)`}
          />
          <WhatYouGet
            rows={[
              ['Go API', 'http://localhost:8080'],
              ['Expo (Expo Go / emulator)', 'exp://…'],
            ]}
          />
          <p className="text-sm text-muted-foreground/70 mt-3">
            On a physical device, point the app at your machine&apos;s LAN IP &mdash; see{' '}
            <Link href="/docs/mobile/getting-started" className="text-primary hover:underline">Mobile · Getting Started</Link>.
          </p>
        </div>
      ),
    },
    {
      id: 'desktop',
      label: 'Desktop',
      icon: <Monitor className="h-4 w-4" />,
      content: (
        <div>
          <p className="text-muted-foreground leading-relaxed mb-4">
            A native desktop app (Wails) with an embedded API on local SQLite &mdash;
            offline-first, no Docker required. Uses its own command:
          </p>
          <CodeBlock
            terminal
            code={`grit new-desktop myapp     # native Wails desktop app (SQLite)
cd myapp
grit start                 # launches the desktop window`}
          />
          <p className="text-sm text-muted-foreground/70 mt-3">
            The local database is created automatically on first run. To build a distributable
            installer, see{' '}
            <Link href="/docs/desktop/building" className="text-primary hover:underline">Desktop · Building &amp; Distribution</Link>.
          </p>
        </div>
      ),
    },
  ]

  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Getting Started</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Create a project</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Everything you need to scaffold and run your first Grit app. Install the CLI once,
                pick what you&apos;re building, and copy the block &mdash; you&apos;ll be running a
                full-stack app in a few minutes.
              </p>
            </div>

            <div className="prose-grit">
              {/* Install */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">1. Install the Grit CLI</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  One line, every platform. The script installs the latest release (or updates an
                  existing install):
                </p>
                <CodeBlock
                  terminal
                  code={`# macOS / Linux
curl -fsSL https://gritframework.dev/install.sh | sh

# Windows (PowerShell)
iwr -useb https://gritframework.dev/install.ps1 | iex`}
                />
                <p className="text-muted-foreground leading-relaxed mt-4">
                  Verify with <code>grit --help</code>. You&apos;ll also need a few tools installed
                  &mdash; <strong>Go 1.24+</strong>, <strong>Node 22+</strong>,{' '}
                  <strong>pnpm 9+</strong>, and <strong>Docker</strong> (skippable for desktop).
                  New to any of them? The{' '}
                  <Link href="/docs/getting-started/prerequisites">Prerequisites</Link> page has a
                  short primer for each.
                </p>
              </div>

              {/* Scaffold — tabs */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">2. Scaffold your app</h2>
                <p className="text-muted-foreground leading-relaxed mb-2">
                  Pick what you&apos;re building. Each tab is a complete, copy-pasteable sequence
                  &mdash; scaffold, set up the database, and run, all with <code>grit</code>.
                </p>
                <Tabs items={tabs} defaultId="all" />
              </div>

              <Callout type="tip" title="Prefer to be asked?">
                Run <code>grit new myapp</code> with no flags and the CLI walks you through
                architecture and frontend choices interactively. Add <code>--vite</code> to use
                TanStack Router (Vite) instead of Next.js.
              </Callout>

              {/* Next */}
              <div className="mb-10 mt-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">3. You&apos;re running &mdash; now build</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Register a user, log into the admin panel, and browse your database in GORM Studio.
                  Then generate your first full-stack resource &mdash; model, API, admin page, types
                  and hooks &mdash; in one command:
                </p>
                <CodeBlock
                  terminal
                  code={`grit generate resource Post --fields "title:string,body:text,published:bool"
grit migrate    # create the new table`}
                />
                <p className="text-muted-foreground leading-relaxed mt-4">
                  Refresh the admin panel and your <strong>Posts</strong> resource is there with a
                  working table and form.
                </p>
              </div>

              <Callout type="escape" title="Escape hatch">
                Everything generated is <strong>your code</strong> &mdash; edit the model, tune the
                admin resource, restyle the screens. Grit generates opinions, not a cage. See the{' '}
                <Link href="/docs/concepts/generated-files">Generated File Map</Link> for exactly
                what each command writes.
              </Callout>

              <div className="flex items-center justify-between border-t border-border pt-8 mt-12">
                <Button variant="ghost" asChild>
                  <Link href="/docs" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Home
                  </Link>
                </Button>
                <Button variant="ghost" asChild>
                  <Link href="/docs/getting-started/coming-from" className="gap-2">
                    Coming from Laravel / Django
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
