import Link from 'next/link'
import {
  ArrowUpRight,
  Bell,
  CreditCard,
  Download,
  Film,
  Globe,
  KeyRound,
  Languages,
  MessageSquare,
  QrCode,
  Search,
  ShieldAlert,
  SquareKanban,
  Table2,
  Terminal,
  UserCog,
  Users,
  Webhook,
  Zap,
} from 'lucide-react'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import {
  CommunityEmpty,
  ModuleRow,
  PluginCard,
  Stat,
  type PluginEntry,
} from '@/components/plugin-card'
import communityIndex from '@/config/community-plugins.json'

export const metadata = {
  title: 'Plugin marketplace — Grit',
  description:
    'Every Grit plugin: the ten built into the CLI, the nine that ship as Go modules, and the community index. What each one does, how to install it, and how to list your own.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/marketplace' },
}

// The official plugins, as `grit plugin list` reports them.
//
// Hand-maintained because a Next page cannot import a Go registry, so
// TestMarketplacePageListsEveryPlugin in internal/plugin compares the two in
// both directions and fails when they drift.
const OFFICIAL: PluginEntry[] = [
  {
    name: 'kanban',
    title: 'Kanban boards',
    summary: 'Boards, columns and cards, dragged with a pointer, a finger or the keyboard.',
    href: '/docs/plugins/kanban',
    icon: <SquareKanban />,
    tags: ['admin', 'models', 'dnd'],
    kind: 'official',
    install: 'grit plugin add kanban',
  },
  {
    name: 'multitenant',
    title: 'Multi-tenancy',
    summary: 'Organizations, per-org roles, and query scoping a forgotten WHERE cannot leak past.',
    href: '/docs/plugins/multitenant',
    icon: <Users />,
    tags: ['models', 'auth', 'scoping'],
    kind: 'official',
    install: 'grit plugin add multitenant',
  },
  {
    name: 'impersonate',
    title: 'Impersonate',
    summary: 'Sign in as another user to reproduce their bug, with an audit trail of every swap.',
    href: '/docs/plugins/impersonate',
    icon: <UserCog />,
    tags: ['auth', 'admin', 'audit'],
    kind: 'official',
    install: 'grit plugin add impersonate',
  },
  {
    name: 'command-palette',
    title: 'Command palette',
    summary: 'Ctrl+K to anywhere in the admin, built from the resource registry.',
    href: '/docs/plugins/command-palette',
    icon: <Terminal />,
    tags: ['admin', 'frontend only'],
    kind: 'official',
    install: 'grit plugin add command-palette',
  },
  {
    name: 'saved-views',
    title: 'Saved views',
    summary: 'A table’s filters, sort and search saved as a named view, per user.',
    href: '/docs/plugins/saved-views',
    icon: <Table2 />,
    tags: ['admin', 'tables'],
    kind: 'official',
    install: 'grit plugin add saved-views',
  },
  {
    name: 'device-pairing',
    title: 'Device pairing',
    summary: 'Sign in a browser by scanning a code from a phone that is already signed in.',
    href: '/docs/plugins/device-pairing',
    icon: <QrCode />,
    tags: ['auth', 'web', 'admin'],
    kind: 'official',
    install: 'grit plugin add device-pairing',
  },
  {
    name: 'push',
    title: 'Push notifications',
    summary: 'Notifications to the Expo app through Expo, so no certificates reach the API.',
    href: '/docs/plugins/push',
    icon: <Bell />,
    tags: ['mobile', 'expo'],
    kind: 'official',
    install: 'grit plugin add push',
  },
  {
    name: 'video',
    title: 'Video',
    summary: 'Clips converted by ffmpeg to an MP4 and a poster every client can play.',
    href: '/docs/plugins/video',
    icon: <Film />,
    tags: ['media', 'jobs'],
    kind: 'official',
    install: 'grit plugin add video',
  },
  {
    name: 'stripe',
    title: 'Stripe',
    summary: 'Checkout and billing, with webhooks that survive a redelivery.',
    href: '/docs/plugins/stripe',
    icon: <CreditCard />,
    tags: ['payments', 'webhooks'],
    kind: 'official',
    install: 'grit plugin add stripe',
  },
  {
    name: 'webhooks',
    title: 'Webhooks',
    summary: 'Outbound webhooks to the Standard Webhooks spec, signed and retried.',
    href: '/docs/backend/webhooks',
    icon: <Webhook />,
    tags: ['api', 'integrations'],
    kind: 'official',
    install: 'grit plugin add webhooks',
  },
]

const MODULES = [
  { name: 'websockets', summary: 'Realtime hub with rooms, presence and typing indicators', icon: <Zap /> },
  { name: 'notifications', summary: 'In-app and push notifications with a delivery log', icon: <Bell /> },
  { name: 'search', summary: 'Full-text search across resources', icon: <Search /> },
  { name: 'stripe', summary: 'Subscriptions and recurring billing', icon: <CreditCard /> },
  { name: 'oauth', summary: 'Sign in with Google, GitHub and friends', icon: <KeyRound /> },
  { name: 'i18n', summary: 'Translated API messages and locale negotiation', icon: <Languages /> },
  { name: 'video', summary: 'Video upload, transcode and playback', icon: <Film /> },
  { name: 'export', summary: 'Scheduled exports to CSV and XLSX', icon: <Download /> },
  { name: 'conference', summary: 'Audio and video rooms', icon: <MessageSquare /> },
]

type CommunityRecord = {
  name: string
  author: string
  repo: string
  summary: string
  added: string
}

const COMMUNITY = ((communityIndex.plugins ?? []) as CommunityRecord[]).map(
  (p): PluginEntry => ({
    name: p.name,
    title: p.name,
    summary: p.summary,
    href: p.repo.startsWith('http') ? p.repo : `https://${p.repo}`,
    icon: <Globe />,
    tags: [],
    kind: 'community',
    author: p.author,
    install: `grit plugin add ${p.repo}`,
  })
)

export default function MarketplacePage() {
  return (
    <DocsShell>
        {/* ---------------------------------------------------------- hero */}
        <section className="relative overflow-hidden border-b border-border">
          <div
            aria-hidden="true"
            className="pointer-events-none absolute inset-0 opacity-[0.05]"
            style={{
              backgroundImage: 'radial-gradient(currentColor 1px, transparent 1px)',
              backgroundSize: '32px 32px',
            }}
          />
          <div className="relative mx-auto max-w-5xl px-6 py-16">
            <p className="font-mono text-[11px] uppercase tracking-[0.12em] text-muted-foreground">
              Plugins
            </p>
            <h1 className="mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
              Marketplace
            </h1>
            <p className="mt-4 max-w-2xl text-lg leading-relaxed text-muted-foreground">
              Every plugin, in three groups. All of them write ordinary code into your project,
              and all of them come back out with one command.
            </p>

            <div className="mt-8 flex flex-wrap gap-3">
              <Link
                href="/docs/plugins/authoring"
                className="inline-flex h-11 items-center gap-2 rounded-full bg-primary px-6 text-sm font-medium text-primary-foreground transition-opacity duration-150 hover:opacity-90"
              >
                Write a plugin
                <ArrowUpRight className="h-4 w-4" aria-hidden="true" />
              </Link>
              <Link
                href="/docs/extending"
                className="inline-flex h-11 items-center gap-2 rounded-full border border-border px-6 text-sm font-medium text-foreground transition-colors duration-150 hover:border-accent/50"
              >
                Is a plugin the right tool?
              </Link>
            </div>

            {/* Stats. Hairline-divided rather than boxed: three boxes in a row
                would be three more objects on a page that is mostly cards. */}
            <div className="mt-10 inline-grid grid-cols-3 divide-x divide-border rounded-lg border border-border">
              <Stat value={String(OFFICIAL.length)} label="Official" />
              <Stat value={String(MODULES.length)} label="Go modules" />
              <Stat value={String(COMMUNITY.length)} label="Community" />
            </div>
          </div>
        </section>

        <div className="mx-auto max-w-5xl px-6">
          {/* ------------------------------------------------------ official */}
          <section className="py-16">
            <div className="mb-6 flex flex-wrap items-end justify-between gap-3">
              <div>
                <h2 className="text-2xl font-semibold tracking-tight">Official</h2>
                <p className="mt-1.5 text-sm text-muted-foreground">
                  Built into the CLI, maintained with the framework, tested on every release.
                </p>
              </div>
              <code className="rounded-sm border border-border px-2.5 py-1.5 font-mono text-xs text-muted-foreground">
                grit plugin list
              </code>
            </div>

            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {OFFICIAL.map((p) => (
                <PluginCard key={p.name} plugin={p} />
              ))}
            </div>
          </section>

          {/* ------------------------------------------------------- modules */}
          <section className="border-t border-border py-16">
            <h2 className="text-2xl font-semibold tracking-tight">Go module plugins</h2>
            <p className="mt-1.5 max-w-2xl text-sm leading-relaxed text-muted-foreground">
              Larger pieces that ship as Go modules from{' '}
              <Link
                href="https://github.com/MUKE-coder/grit-plugins"
                className="text-accent hover:underline"
              >
                grit-plugins
              </Link>
              . Installed with <code className="font-mono text-xs">go get</code> rather than{' '}
              <code className="font-mono text-xs">grit plugin add</code>, because they are
              libraries your code imports rather than code generated into it.
            </p>

            <div className="mt-6 rounded-lg border border-border px-5">
              {MODULES.map((m) => (
                <ModuleRow key={m.name} name={m.name} summary={m.summary} />
              ))}
            </div>

            <div className="mt-4">
              <CodeBlock
                language="bash"
                code={`go get github.com/MUKE-coder/grit-plugins/grit-websockets`}
              />
            </div>
          </section>

          {/* ----------------------------------------------------- community */}
          <section className="border-t border-border py-16">
            <h2 className="text-2xl font-semibold tracking-tight">Community</h2>
            <p className="mt-1.5 max-w-2xl text-sm text-muted-foreground">
              Written by other people and listed here on request.
            </p>

            <div className="mt-6 flex gap-3 rounded-lg border border-border bg-card p-5">
              <ShieldAlert className="h-5 w-5 shrink-0 text-warning" aria-hidden="true" />
              <div className="text-sm leading-relaxed text-muted-foreground">
                <p className="font-medium text-foreground">Listing is not endorsement</p>
                <p className="mt-1.5">
                  The review is deliberately shallow: that the repository exists, that it
                  installs and removes cleanly on a fresh project, and that it is not obviously
                  hostile. Nobody audits the code for you. A plugin writes files into your
                  project and patches existing ones, so read what it did with{' '}
                  <code className="font-mono text-xs">git diff</code>, the same as any
                  dependency.
                </p>
              </div>
            </div>

            <div className="mt-6">
              {COMMUNITY.length === 0 ? (
                <CommunityEmpty />
              ) : (
                <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                  {COMMUNITY.map((p) => (
                    <PluginCard key={p.name} plugin={p} />
                  ))}
                </div>
              )}
            </div>

            <p className="mt-6 text-sm text-muted-foreground">
              Installing one takes the repository directly. Grit names the source and confirms
              before cloning, and nothing from the repository is executed: it reads the
              manifest and writes the files that manifest names.
            </p>
            <div className="mt-3">
              <CodeBlock language="bash" code={`grit plugin add github.com/someone/grit-calendar`} />
            </div>
          </section>

          {/* -------------------------------------------------------- submit */}
          <section className="border-t border-border py-16">
            <h2 className="text-2xl font-semibold tracking-tight">Listing your own</h2>
            <p className="mt-1.5 max-w-2xl text-sm text-muted-foreground">
              <code className="font-mono text-xs">grit plugin submit</code> prints the checklist
              and the submission link. Three commands on a freshly scaffolded project, and the
              third is the one people skip:
            </p>

            <ol className="mt-5 space-y-2.5">
              {[
                ['grit plugin add <your repo>', 'installs without errors'],
                ['go build ./...', 'compiles afterwards'],
                ['grit plugin remove <name>', 'leaves the project compiling again'],
              ].map(([cmd, what], i) => (
                <li key={cmd} className="flex items-center gap-4">
                  <span
                    aria-hidden="true"
                    className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-accent/10 font-mono text-xs font-semibold text-accent"
                  >
                    {i + 1}
                  </span>
                  <code className="font-mono text-sm text-foreground">{cmd}</code>
                  <span className="text-sm text-muted-foreground">{what}</span>
                </li>
              ))}
            </ol>

            <p className="mt-5 max-w-2xl text-sm leading-relaxed text-muted-foreground">
              Removal replays the lockfile backwards, so it only works if every change the
              plugin made is one the lockfile recorded. A plugin that edits a file outside the
              injection mechanism installs fine and leaves wreckage when it is taken out.
            </p>
          </section>
        </div>
    </DocsShell>
  )
}
