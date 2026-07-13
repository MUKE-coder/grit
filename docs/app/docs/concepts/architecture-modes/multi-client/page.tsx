import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { FileTree } from '@/components/diagram'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/architecture-modes/multi-client')

export default function MultiClientArchitecturePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="max-w-4xl mx-auto py-12 px-6 lg:px-8">
          {/* Header */}
          <div className="mb-14">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">
              Architecture Modes
            </p>
            <h1 className="text-3xl lg:text-4xl font-bold tracking-tight text-foreground mb-4">
              Multi-Client: API + Mobile + Desktop
            </h1>
            <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl">
              One Go API serving an Expo mobile app AND a Wails desktop app -- all
              in a single monorepo with shared types. The pattern real SaaS products
              use (Linear, Notion, Slack).
            </p>
            <p className="text-sm text-muted-foreground/70 mt-3">
              Added in v3.9.0 via the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--desktop</code> flag.
            </p>
          </div>

          <div className="prose-grit">
            {/* Overview */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Overview
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The multi-client pattern combines a Go API with two or more native clients
                (mobile via Expo, desktop via Wails, optionally web via Next.js or TanStack).
                All clients share the same{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">packages/shared/</code>{' '}
                types and schemas, and all talk to the same API over HTTP. No code
                duplication across clients.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                This isn&apos;t a separate architecture mode -- it&apos;s a{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--desktop</code> flag
                you combine with an existing architecture (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--mobile</code>,{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--triple</code>,{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--double</code>, or{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--api</code>).
              </p>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-5 mb-6">
                <h4 className="text-sm font-semibold text-foreground mb-3">The command that started it all</h4>
                <CodeBlock language="bash" code={`grit new myapp --mobile --desktop --next`} />
                <p className="text-xs text-muted-foreground mt-3 leading-relaxed">
                  Scaffolds: Go API + Next.js web + admin panel + Expo mobile + Wails desktop,
                  all in one monorepo, all sharing types.
                </p>
              </div>
            </div>

            {/* When to use */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                When to Use This Pattern
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Pick the multi-client pattern when you&apos;re building a product that needs to
                live in multiple places:
              </p>
              <ul className="list-disc pl-6 space-y-2 text-muted-foreground mb-4">
                <li>
                  <strong className="text-foreground">Native apps AND a web presence</strong> -- e.g. a Slack-style team
                  tool with desktop, mobile, and a browser fallback
                </li>
                <li>
                  <strong className="text-foreground">A mobile-first product that also needs a power-user desktop app</strong>
                  {' '}-- e.g. a task manager like Things or Linear
                </li>
                <li>
                  <strong className="text-foreground">Any product where customers expect to pick up where they left off</strong>
                  {' '}across phone and laptop -- e.g. chat apps, note-takers, calendars
                </li>
                <li>
                  <strong className="text-foreground">Enterprise tools</strong> where clients specifically request a native
                  desktop installer instead of a browser tab
                </li>
              </ul>
              <p className="text-muted-foreground leading-relaxed">
                <strong className="text-foreground">Don&apos;t</strong> use this for simple SaaS apps where a browser-only
                experience is enough -- stick with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--triple</code>.
                Don&apos;t use it for <strong className="text-foreground">offline-first</strong> desktop apps either --
                use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new-desktop</code> instead
                (standalone Wails + embedded Go + local SQLite).
              </p>
            </div>

            {/* --desktop vs grit new-desktop */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                <code className="text-xl font-mono bg-accent/50 px-1.5 py-0.5 rounded">--desktop</code> vs{' '}
                <code className="text-xl font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new-desktop</code>
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                These are two different desktop patterns. Pick the one that matches your product:
              </p>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80"></th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                        <code className="text-xs font-mono">grit new-desktop</code>
                      </th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                        <code className="text-xs font-mono">grit new --desktop</code>
                      </th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Network</td>
                      <td className="px-4 py-2.5">Offline-first, no remote API</td>
                      <td className="px-4 py-2.5">Always online, calls remote API</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Backend</td>
                      <td className="px-4 py-2.5">Embedded in the desktop binary</td>
                      <td className="px-4 py-2.5">Shared monorepo API (Gin + GORM)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Database</td>
                      <td className="px-4 py-2.5">Local SQLite</td>
                      <td className="px-4 py-2.5">Remote PostgreSQL via API</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Wails bindings</td>
                      <td className="px-4 py-2.5">Full CRUD, auth, exports (19+ methods)</td>
                      <td className="px-4 py-2.5">Only native OS (window, file dialogs, keychain)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Token storage</td>
                      <td className="px-4 py-2.5">Local-only session</td>
                      <td className="px-4 py-2.5">OS keychain via 99designs/keyring</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Works with mobile?</td>
                      <td className="px-4 py-2.5">No -- it&apos;s standalone</td>
                      <td className="px-4 py-2.5">Yes -- shares packages/shared with Expo</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">Use case</td>
                      <td className="px-4 py-2.5">Invoicer, notes, local utilities</td>
                      <td className="px-4 py-2.5">Multi-client SaaS, team tools</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* Scaffold combinations */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Scaffold Combinations
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--desktop</code> flag
                combines with every architecture mode except <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--single</code>
                (single apps already bundle their own frontend). Here are the common setups:
              </p>

              <div className="space-y-5">
                <div>
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    1. Full multi-client SaaS (the Linear setup)
                  </h3>
                  <CodeBlock language="bash" code={`grit new myapp --triple --next --desktop`} />
                  <p className="text-sm text-muted-foreground leading-relaxed mt-2">
                    Scaffolds: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/api</code>,{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/web</code>,{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/admin</code>,{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/desktop</code>, plus{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">packages/shared</code>.
                    You get a web marketing site, an admin dashboard for ops, and a native desktop app -- all sharing one API.
                  </p>
                </div>

                <div>
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    2. Mobile + desktop with a marketing site
                  </h3>
                  <CodeBlock language="bash" code={`grit new myapp --triple --next --mobile --desktop`} />
                  <p className="text-sm text-muted-foreground leading-relaxed mt-2">
                    Scaffolds everything above PLUS <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/expo</code>.
                    This is the full Slack / Notion setup: every client surface.
                  </p>
                </div>

                <div>
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    3. Desktop + mobile, no web app
                  </h3>
                  <CodeBlock language="bash" code={`grit new myapp --mobile --desktop`} />
                  <p className="text-sm text-muted-foreground leading-relaxed mt-2">
                    Scaffolds <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/api</code>,{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/expo</code>,{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/desktop</code>. Good for
                    products that are native-only (no marketing page yet).
                  </p>
                </div>

                <div>
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    4. Minimal: API + desktop only
                  </h3>
                  <CodeBlock language="bash" code={`grit new myapp --api --desktop`} />
                  <p className="text-sm text-muted-foreground leading-relaxed mt-2">
                    Scaffolds <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/api</code> and{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/desktop</code>. The leanest
                    multi-client setup.
                  </p>
                </div>

                <div>
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    5. Not allowed: <code className="text-lg font-mono bg-accent/50 px-1.5 py-0.5 rounded">--single --desktop</code>
                  </h3>
                  <CodeBlock language="bash" code={`# This errors out:
grit new myapp --single --desktop
# → --desktop is not supported with --single architecture`} />
                  <p className="text-sm text-muted-foreground leading-relaxed mt-2">
                    Single-app architectures already bundle their own SPA into one Go binary. Adding a
                    separate desktop client would be redundant. Use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--triple</code>{' '}
                    or <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--api</code> if you
                    want a desktop app.
                  </p>
                </div>
              </div>
            </div>

            {/* Project structure */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Project Structure
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The shape of <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new myapp --triple --next --desktop --mobile</code>:
              </p>
              <FileTree
                title="myapp/"
                nodes={[
                  { name: 'apps/', type: 'folder', depth: 0 },
                  { name: 'api/', type: 'folder', depth: 1, comment: 'Go backend (Gin + GORM + PostgreSQL)' },
                  { name: 'cmd/server/main.go', type: 'file', depth: 2 },
                  { name: 'internal/', type: 'folder', depth: 2 },
                  { name: 'handlers/', type: 'folder', depth: 3, comment: 'HTTP handlers' },
                  { name: 'services/', type: 'folder', depth: 3, comment: 'Business logic' },
                  { name: 'models/', type: 'folder', depth: 3, comment: 'GORM models' },
                  { name: 'routes/', type: 'folder', depth: 3, comment: 'Route registration' },
                  { name: 'web/', type: 'folder', depth: 1, comment: 'Next.js marketing/landing' },
                  { name: 'admin/', type: 'folder', depth: 1, comment: 'Next.js admin dashboard' },
                  { name: 'expo/', type: 'folder', depth: 1, comment: 'Expo mobile (iOS + Android)' },
                  { name: 'app/', type: 'folder', depth: 2, comment: 'Expo Router screens' },
                  { name: 'lib/', type: 'folder', depth: 2, comment: 'API client, auth, SecureStore' },
                  { name: 'desktop/', type: 'folder', depth: 1, comment: 'Wails desktop (macOS/Windows/Linux)' },
                  { name: 'main.go', type: 'file', depth: 2, comment: 'Wails bootstrap, frameless window' },
                  { name: 'app.go', type: 'file', depth: 2, comment: 'Native OS bindings only' },
                  { name: 'internal/', type: 'folder', depth: 2 },
                  { name: 'keychain.go', type: 'file', depth: 3, comment: 'OS keychain wrapper' },
                  { name: 'frontend/', type: 'folder', depth: 2 },
                  { name: 'src/', type: 'folder', depth: 3 },
                  { name: 'routes/', type: 'folder', depth: 4, comment: 'TanStack Router file-based routes' },
                  { name: 'components/', type: 'folder', depth: 4, comment: 'Shared UI (layout, ui/, command-palette)' },
                  { name: 'lib/', type: 'folder', depth: 4, comment: 'API client, Wails bridge, shortcuts' },
                  { name: 'hooks/', type: 'folder', depth: 4, comment: 'useAuth, useMe, useLogout' },
                  { name: 'vite.config.ts', type: 'file', depth: 3 },
                  { name: 'packages/', type: 'folder', depth: 0 },
                  { name: 'shared/', type: 'folder', depth: 1, comment: '🔑 Shared across ALL clients' },
                  { name: 'schemas/', type: 'folder', depth: 2, comment: 'Zod schemas (login, register, user...)' },
                  { name: 'types/', type: 'folder', depth: 2, comment: 'TypeScript types' },
                  { name: 'constants/', type: 'folder', depth: 2, comment: 'ROUTES, ROLES' },
                  { name: 'docker-compose.yml', type: 'file', depth: 0 },
                  { name: 'turbo.json', type: 'file', depth: 0 },
                  { name: 'pnpm-workspace.yaml', type: 'file', depth: 0 },
                ]}
              />
              <p className="text-muted-foreground leading-relaxed mt-4">
                The key insight: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">packages/shared</code>{' '}
                is the glue. Change a Zod schema there and both the Expo app and the desktop app pick it up
                immediately. Change a Go model, run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit sync</code>,
                and the TypeScript types regenerate for every client.
              </p>
            </div>

            {/* Development workflow */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Development Workflow
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                With all apps running, you get live hot-reload across everything. One
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit start</code> runs
                them all, or start each app on its own:
              </p>

              <CodeBlock
                language="bash"
                code={`# Terminal 1: Start infrastructure (Postgres, Redis, MinIO, Mailhog)
docker compose up -d

# Terminal 2: Go API (hot reload)
grit start server

# Terminal 3: Mobile app (Expo)
grit start expo

# Terminal 4: Desktop app (Wails)
grit start desktop

# All of these can be run together with:
grit start`}
              />

              <p className="text-muted-foreground leading-relaxed mt-4 mb-4">
                Or start a single app from anywhere in the project:
              </p>

              <CodeBlock
                language="bash"
                code={`grit start            # All apps
grit start server     # Just the Go API
grit start web        # Just the web frontend
grit start admin      # Just the admin panel
grit start expo       # Just the Expo app
grit start desktop    # Just the desktop app`}
              />
            </div>

            {/* How the clients share code */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                How Clients Share Code
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The three client patterns (web, mobile, desktop) use the same primitives but with
                platform-appropriate adapters:
              </p>

              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-4">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Concept</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Web</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Mobile</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Desktop</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Types/schemas</td>
                      <td className="px-4 py-2.5" colSpan={3}>
                        <code className="text-xs">packages/shared</code> -- identical across all three
                      </td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">HTTP client</td>
                      <td className="px-4 py-2.5">axios</td>
                      <td className="px-4 py-2.5">axios</td>
                      <td className="px-4 py-2.5">axios</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Data fetching</td>
                      <td className="px-4 py-2.5">TanStack Query</td>
                      <td className="px-4 py-2.5">TanStack Query</td>
                      <td className="px-4 py-2.5">TanStack Query</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Forms</td>
                      <td className="px-4 py-2.5">react-hook-form + zod</td>
                      <td className="px-4 py-2.5">react-hook-form + zod</td>
                      <td className="px-4 py-2.5">react-hook-form + zod</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Router</td>
                      <td className="px-4 py-2.5">Next.js App Router</td>
                      <td className="px-4 py-2.5">Expo Router</td>
                      <td className="px-4 py-2.5">TanStack Router</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Token storage</td>
                      <td className="px-4 py-2.5">HttpOnly cookies set by API (never localStorage)</td>
                      <td className="px-4 py-2.5">expo-secure-store</td>
                      <td className="px-4 py-2.5">OS keychain (via Wails bridge)</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">Styling</td>
                      <td className="px-4 py-2.5">Tailwind + shadcn/ui</td>
                      <td className="px-4 py-2.5">NativeWind</td>
                      <td className="px-4 py-2.5">Tailwind + custom primitives</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <p className="text-muted-foreground leading-relaxed">
                The useAuth hook signature is identical across all three:{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useMe()</code>,{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useLogin()</code>,{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useRegister()</code>,{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useLogout()</code>. Only the{' '}
                <em>adapter</em> (where tokens get stored) differs. Copy a feature from web to desktop and the
                structure stays the same.
              </p>
            </div>

            {/* Desktop app details */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                What Ships in the Desktop App
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The Wails desktop scaffold is intentionally minimal on the Go side and rich on the React side.
                Go handles only what the browser can&apos;t: window chrome, file dialogs, and OS keychain.
                Everything else lives in the React frontend and calls the API.
              </p>

              <h3 className="text-lg font-semibold text-foreground mb-2 mt-6">Go side (minimal)</h3>
              <ul className="list-disc pl-6 space-y-1 text-muted-foreground mb-4">
                <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">main.go</code> -- Wails bootstrap, frameless 1280x800 window (min 1000x700)</li>
                <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">app.go</code> -- Native bindings: window controls, file dialogs, keychain, platform detection</li>
                <li><code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">internal/keychain.go</code> -- Cross-platform keychain via <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">99designs/keyring</code></li>
              </ul>

              <h3 className="text-lg font-semibold text-foreground mb-2">React side (TanStack Router)</h3>
              <ul className="list-disc pl-6 space-y-1 text-muted-foreground mb-4">
                <li>
                  <strong className="text-foreground">Custom title bar</strong> with platform-aware window controls
                  (macOS traffic lights on the left, Windows/Linux controls on the right)
                </li>
                <li>
                  <strong className="text-foreground">Command palette</strong> (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Cmd+K</code> / <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Ctrl+K</code>)
                  with keyboard-first navigation
                </li>
                <li>
                  <strong className="text-foreground">Fixed 240px sidebar</strong> (not collapsible -- desktop windows are wide enough)
                </li>
                <li>
                  <strong className="text-foreground">Topbar</strong> with search, notifications, theme toggle, user menu
                </li>
                <li>
                  <strong className="text-foreground">Pre-built auth flow</strong> (login / register / forgot-password)
                  using react-hook-form + zod
                </li>
                <li>
                  <strong className="text-foreground">Dashboard, Profile, Settings</strong> pages out of the box
                </li>
                <li>
                  <strong className="text-foreground">Global keyboard shortcuts</strong> via{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useShortcuts()</code>
                </li>
              </ul>
            </div>

            {/* Keyboard shortcuts */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Keyboard Shortcuts (Default)
              </h2>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Shortcut</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Action</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">⌘K / Ctrl+K</td>
                      <td className="px-4 py-2.5">Open command palette</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">⌘, / Ctrl+,</td>
                      <td className="px-4 py-2.5">Open Settings</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">⌘L / Ctrl+L</td>
                      <td className="px-4 py-2.5">Log out</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Esc</td>
                      <td className="px-4 py-2.5">Close modal / palette</td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <p className="text-muted-foreground leading-relaxed mt-4">
                Add your own with the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useShortcuts()</code> hook:
              </p>
              <CodeBlock
                language="tsx"
                code={`import { useShortcuts } from "@/lib/use-shortcuts";

useShortcuts({
  "mod+n": () => createNewItem(),        // ⌘N / Ctrl+N
  "mod+shift+p": () => openProfile(),    // ⌘⇧P / Ctrl+Shift+P
  "mod+slash": () => openHelp(),         // ⌘/ / Ctrl+/
});`}
              />
            </div>

            {/* Deployment */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Deploying a Multi-Client App
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Each client ships differently. The shared API is deployed once, and the clients
                point to it:
              </p>

              <h3 className="text-lg font-semibold text-foreground mb-2 mt-5">
                1. API (one-time deploy)
              </h3>
              <CodeBlock
                language="bash"
                code={`grit deploy --host user@server.com --domain api.myapp.com`}
              />
              <p className="text-sm text-muted-foreground leading-relaxed mt-2">
                Cross-compiles, uploads via SCP, configures systemd + Caddy with auto-TLS.
              </p>

              <h3 className="text-lg font-semibold text-foreground mb-2 mt-5">
                2. Web / Admin (any time)
              </h3>
              <CodeBlock
                language="bash"
                code={`# Deploy to Vercel (pointing NEXT_PUBLIC_API_URL=https://api.myapp.com)
cd apps/web && vercel --prod
cd apps/admin && vercel --prod`}
              />

              <h3 className="text-lg font-semibold text-foreground mb-2 mt-5">
                3. Mobile (EAS Build)
              </h3>
              <CodeBlock
                language="bash"
                code={`cd apps/expo
npx eas-cli build --platform ios
npx eas-cli build --platform android
npx eas-cli submit --platform ios
npx eas-cli submit --platform android`}
              />

              <h3 className="text-lg font-semibold text-foreground mb-2 mt-5">
                4. Desktop (Wails build)
              </h3>
              <CodeBlock
                language="bash"
                code={`cd apps/desktop
wails build -platform darwin/amd64   # macOS Intel
wails build -platform darwin/arm64   # macOS Apple Silicon
wails build -platform windows/amd64
wails build -platform linux/amd64
# Binaries land in build/bin/`}
              />
              <p className="text-sm text-muted-foreground leading-relaxed mt-2">
                Code-sign and notarize as needed (macOS) or sign with a certificate (Windows).
                Distribute as standalone binaries, via a download page, or through app stores.
              </p>
            </div>

            {/* Tips */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Tips &amp; Gotchas
              </h2>
              <ul className="list-disc pl-6 space-y-3 text-muted-foreground">
                <li>
                  <strong className="text-foreground">Set <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">VITE_API_URL</code></strong>{' '}
                  in the desktop frontend&apos;s env when deploying. In dev, Vite proxies{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/api</code> to{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">localhost:8080</code>, but once built,
                  the frontend needs to know the full production URL.
                </li>
                <li>
                  <strong className="text-foreground">CORS matters again.</strong> Native clients don&apos;t enforce CORS the
                  way browsers do, but your API still needs to allow the desktop origin if you proxy
                  through one. Usually fine with Grit&apos;s default CORS middleware; check if you customize.
                </li>
                <li>
                  <strong className="text-foreground">Keep Wails bindings minimal.</strong> Resist the urge to expose CRUD
                  methods as Wails bindings -- it creates a second API surface to maintain. Do data work
                  over HTTP, use bindings only for window/file/keychain.
                </li>
                <li>
                  <strong className="text-foreground">Dev without Wails.</strong> You can run the desktop frontend in a
                  regular browser with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">cd apps/desktop/frontend &amp;&amp; pnpm dev</code>{' '}
                  (it falls back to HttpOnly cookies set by the dev API). Useful for quick UI iteration without
                  launching the full Wails window. Tokens are never written to <code>localStorage</code> in either mode.
                </li>
                <li>
                  <strong className="text-foreground">Match platform expectations.</strong> macOS users expect traffic lights
                  on the left; Windows/Linux users expect close/min/max on the right. The scaffold handles
                  this automatically via <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">getPlatform()</code>.
                </li>
                <li>
                  <strong className="text-foreground">Read the style guide.</strong>{' '}
                  <Link href="https://github.com/MUKE-coder/grit/blob/main/GRIT_STYLE_GUIDE.md" className="text-primary hover:underline">
                    GRIT_STYLE_GUIDE.md §14.5
                  </Link>{' '}
                  covers desktop-specific patterns: no breadcrumbs, no collapsible sidebar, 32px content
                  padding (vs web&apos;s 24px), tighter focus rings, OS keychain only -- never localStorage.
                </li>
              </ul>
            </div>

            {/* Related */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Related Pages
              </h2>
              <ul className="list-disc pl-6 space-y-2 text-muted-foreground">
                <li>
                  <Link href="/docs/concepts/architecture-modes/triple" className="text-primary hover:underline">
                    Triple architecture
                  </Link>
                  {' '}-- the base you&apos;ll usually combine <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">--desktop</code> with
                </li>
                <li>
                  <Link href="/docs/concepts/architecture-modes/mobile" className="text-primary hover:underline">
                    Mobile architecture
                  </Link>
                  {' '}-- details on the Expo app
                </li>
                <li>
                  <Link href="/docs/concepts/architecture-modes/api-only" className="text-primary hover:underline">
                    API-only architecture
                  </Link>
                  {' '}-- for minimal API + desktop setups
                </li>
                <li>
                  <Link href="/docs/concepts/cli" className="text-primary hover:underline">
                    CLI reference
                  </Link>
                  {' '}-- full flag list
                </li>
              </ul>
            </div>
          </div>

          {/* Footer nav */}
          <div className="flex items-center justify-between pt-8 border-t border-border/40">
            <Button variant="ghost" asChild>
              <Link href="/docs/concepts/architecture-modes/mobile" className="gap-2">
                <ArrowLeft className="h-4 w-4" />
                Mobile Architecture
              </Link>
            </Button>
            <Button variant="ghost" asChild>
              <Link href="/docs/concepts/cli" className="gap-2">
                CLI Reference
                <ArrowRight className="h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
