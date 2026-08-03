import Link from 'next/link'
import { ArrowLeft, ArrowRight, Server, Smartphone, Monitor, Layers } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { Tabs } from '@/components/tabs'
import { Steps, Step } from '@/components/steps'
import { PageHelp } from '@/components/page-help'
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
            Every target at once &mdash; Go API, Next.js web app, admin panel, Wails desktop
            app, Expo mobile app, and a docs site. More than most projects need, but the
            fastest way to see what Grit does.
          </p>
          <CodeBlock
            terminal
            code={`grit new myapp --full      # api + web + admin + desktop + expo + docs
cd myapp
docker compose up -d       # Postgres, Redis, MinIO, Mailhog
pnpm install               # frontend deps (one-time)
grit migrate               # create database tables
grit seed                  # sample data + a demo admin login
grit start                 # runs every app, Ctrl+C stops them all`}
          />
          <WhatYouGet
            rows={[
              ['Web app', 'http://localhost:3000'],
              ['Admin panel', 'http://localhost:3001'],
              ['Go API + docs', 'http://localhost:8080/docs'],
              ['GORM Studio', 'http://localhost:8080/studio'],
              ['Desktop + Expo', 'grit start desktop / grit start expo'],
            ]}
          />
          <p className="text-sm text-muted-foreground/70 mt-3">
            Six apps is a lot to run at once. If you only want the web stack, use the{' '}
            <strong>Web + Admin + API</strong> tab &mdash; you can add desktop or mobile later.
          </p>
        </div>
      ),
    },
    {
      id: 'triple',
      label: 'Web + Admin + API',
      icon: <Layers className="h-4 w-4" />,
      content: (
        <div>
          <p className="text-muted-foreground leading-relaxed mb-4">
            The recommended starting point. Go API, a Next.js web app and the Filament-style
            admin panel &mdash; nothing you will not use on day one.
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
      id: 'triple-mobile',
      label: '+ Mobile',
      icon: <Smartphone className="h-4 w-4" />,
      content: (
        <div>
          <p className="text-muted-foreground leading-relaxed mb-4">
            The web stack plus an Expo app, all sharing one API and one set of generated
            types. Generate a resource and it appears on the phone too.
          </p>
          <CodeBlock
            terminal
            code={`grit new myapp --triple --expo   # api + web + admin + expo
cd myapp
docker compose up -d
pnpm install
grit migrate
grit seed
grit start                       # api + web + admin
grit start expo                  # second terminal — Expo dev server`}
          />
          <WhatYouGet
            rows={[
              ['Web app', 'http://localhost:3000'],
              ['Admin panel', 'http://localhost:3001'],
              ['Go API + docs', 'http://localhost:8080/docs'],
              ['Expo (Expo Go / emulator)', 'scan the QR code'],
            ]}
          />
          <p className="text-sm text-muted-foreground/70 mt-3">
            Expo runs in its own terminal on purpose &mdash; its dev server wants the
            foreground for the QR code and the keyboard shortcuts. On a physical device, point
            the app at your machine&apos;s LAN IP: see{' '}
            <Link href="/docs/mobile/getting-started" className="text-primary hover:underline">
              Mobile &middot; Getting Started
            </Link>
            .
          </p>
        </div>
      ),
    },
    {
      id: 'triple-desktop',
      label: '+ Desktop',
      icon: <Monitor className="h-4 w-4" />,
      content: (
        <div>
          <p className="text-muted-foreground leading-relaxed mb-4">
            The web stack plus a Wails desktop app that shares the same API and models, with a
            local SQLite mirror so it keeps working with no network.
          </p>
          <CodeBlock
            terminal
            code={`grit new myapp --triple --desktop   # api + web + admin + desktop
cd myapp
docker compose up -d
pnpm install
grit migrate
grit seed
grit start                          # api + web + admin + desktop`}
          />
          <WhatYouGet
            rows={[
              ['Web app', 'http://localhost:3000'],
              ['Admin panel', 'http://localhost:3001'],
              ['Go API + docs', 'http://localhost:8080/docs'],
              ['Desktop app', 'opens as a native window'],
            ]}
          />
          <p className="text-sm text-muted-foreground/70 mt-3">
            Wails needs a C toolchain and a WebView runtime on your machine &mdash; see{' '}
            <Link href="/docs/desktop" className="text-primary hover:underline">
              Desktop
            </Link>{' '}
            for the per-OS prerequisites. Build an installer later with{' '}
            <code className="text-foreground/80">grit package</code>.
          </p>
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
              <Steps>
                {/* Install */}
                <Step title="Install the Grit CLI">
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
                </Step>

                {/* Scaffold — tabs */}
                <Step title="Scaffold your app">
                  <p className="text-muted-foreground leading-relaxed mb-2">
                    Pick what you&apos;re building. Each tab is a complete, copy-pasteable sequence
                    &mdash; scaffold, set up the database, and run, all with <code>grit</code>.
                  </p>
                  <Tabs items={tabs} defaultId="all" />
                  <Callout type="tip" title="Prefer to be asked?">
                    Run <code>grit new myapp</code> with no flags and the CLI walks you through
                    architecture and frontend choices interactively. Add <code>--vite</code> to use
                    TanStack Router (Vite) instead of Next.js.
                  </Callout>
                </Step>

                {/* Next */}
                <Step title="You're running — now build">
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
                  <Callout type="escape" title="Escape hatch">
                    Everything generated is <strong>your code</strong> &mdash; edit the model, tune the
                    admin resource, restyle the screens. Grit generates opinions, not a cage. See the{' '}
                    <Link href="/docs/concepts/generated-files">Generated File Map</Link> for exactly
                    what each command writes.
                  </Callout>
                </Step>
              </Steps>

              <PageHelp
                faqs={[
                  {
                    q: 'Do I need Docker?',
                    a: (
                      <>
                        For the web/API/mobile kits, Docker runs Postgres, Redis, MinIO and Mailhog
                        locally &mdash; or use cloud services (Neon, Upstash, R2) instead. The{' '}
                        <strong>Desktop</strong> kit needs no Docker at all: it uses local SQLite.
                      </>
                    ),
                  },
                  {
                    q: 'Which architecture should I pick?',
                    a: (
                      <>
                        Start with <strong>Web + Admin + API</strong> unless you know you don&apos;t
                        need a frontend. Pick <strong>API only</strong> for a headless backend,{' '}
                        <strong>Mobile</strong> for Expo, <strong>Desktop</strong> for a native app.
                        You can always add more later. See{' '}
                        <Link href="/docs/concepts/architecture-modes">Architecture Modes</Link>.
                      </>
                    ),
                  },
                  {
                    q: 'Next.js or Vite?',
                    a: (
                      <>
                        Next.js is the default. Add <code>--vite</code> to scaffold a TanStack Router
                        (Vite) frontend instead, e.g. <code>grit new myapp --triple --vite</code>.
                      </>
                    ),
                  },
                  {
                    q: 'How do I add a database table?',
                    a: (
                      <>
                        <code>grit generate resource &lt;Name&gt; --fields &quot;…&quot;</code> then{' '}
                        <code>grit migrate</code>. It creates the model, API, admin page, types and
                        hooks in one go &mdash; see{' '}
                        <Link href="/docs/concepts/generated-files">the Generated File Map</Link>.
                      </>
                    ),
                  },
                ]}
              />

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
