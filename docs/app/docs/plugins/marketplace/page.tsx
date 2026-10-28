import Link from 'next/link'
import { ArrowRight, Github, ShieldAlert, Store } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import communityIndex from '@/config/community-plugins.json'

export const metadata = {
  title: 'Plugin marketplace — Grit',
  description:
    'Every Grit plugin: the official ones built into the CLI, the Go module plugins, and the community index. How to install one and how to list your own.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/marketplace' },
}

// The official plugins, as `grit plugin list` reports them.
//
// This list is hand-maintained and the CLI's registry is authoritative, so
// TestMarketplacePageListsEveryPlugin in internal/plugin compares the two and
// fails when a plugin is added without reaching this page. A marketplace that
// quietly omits a plugin is worse than no marketplace.
const OFFICIAL = [
  { name: 'kanban', summary: 'Boards, columns and cards with drag-and-drop ordering', href: '/docs/plugins/kanban' },
  { name: 'multitenant', summary: 'Organizations, per-org roles, and automatic query scoping', href: '/docs/plugins/multitenant' },
  { name: 'impersonate', summary: 'Admin signs in as another user, with an audit trail', href: '/docs/plugins/impersonate' },
  { name: 'command-palette', summary: 'Command palette to jump anywhere in the admin', href: '/docs/plugins/command-palette' },
  { name: 'saved-views', summary: 'Save a table’s filters, sort and columns as named views', href: '/docs/plugins/saved-views' },
  { name: 'device-pairing', summary: 'Sign in a browser by scanning a QR code from a signed-in phone', href: '/docs/plugins/device-pairing' },
  { name: 'push', summary: 'Push notifications to the Expo app, through Expo’s push service', href: '/docs/plugins/push' },
  { name: 'video', summary: 'Video upload, transcode and playback', href: '/docs/plugins/video' },
  { name: 'stripe', summary: 'Checkout and billing', href: '/docs/plugins/stripe' },
  { name: 'webhooks', summary: 'Outbound webhooks, to the Standard Webhooks spec', href: '/docs/backend/webhooks' },
]

const MODULES = [
  { name: 'websockets', summary: 'Realtime hub with rooms, presence and typing indicators' },
  { name: 'notifications', summary: 'In-app and push notifications with a delivery log' },
  { name: 'search', summary: 'Full-text search across resources' },
  { name: 'stripe', summary: 'Subscriptions and billing' },
  { name: 'oauth', summary: 'Sign in with Google, GitHub and friends' },
  { name: 'i18n', summary: 'Translated API messages and locale negotiation' },
  { name: 'video', summary: 'Video upload, transcode and playback' },
  { name: 'export', summary: 'Scheduled exports to CSV and XLSX' },
  { name: 'conference', summary: 'Audio and video rooms' },
]

type CommunityPlugin = {
  name: string
  author: string
  repo: string
  summary: string
  added: string
}

const COMMUNITY = (communityIndex.plugins ?? []) as CommunityPlugin[]

export default function MarketplacePage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <Store className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Plugins
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">Marketplace</h1>
          <p className="mb-8 text-lg leading-relaxed text-muted-foreground">
            Every plugin, in three groups: the ones built into the CLI, the ones that ship as Go
            modules, and the ones other people wrote. All of them write ordinary code into your
            project, and all of them come back out with{' '}
            <code>grit plugin remove</code>.
          </p>

          <CodeBlock language="bash" code={`grit plugin list`} />
          <p className="mb-10 mt-3 text-sm text-muted-foreground">
            The same list, for the version of Grit you have installed. That one is authoritative;
            this page is the browsable copy.
          </p>

          <h2 className="mb-2 mt-12 text-2xl font-semibold tracking-tight">Official</h2>
          <p className="mb-5 leading-relaxed text-muted-foreground">
            Built into the CLI, maintained with the framework, and tested on every release.
          </p>
          <div className="mb-8 space-y-2">
            {OFFICIAL.map((p) => (
              <Link
                key={p.name}
                href={p.href}
                className="group flex items-start gap-4 rounded-xl border border-border p-4 transition-colors hover:border-primary"
              >
                <code className="mt-0.5 shrink-0 font-mono text-sm text-primary">{p.name}</code>
                <span className="flex-1 text-sm text-muted-foreground">{p.summary}</span>
                <ArrowRight className="mt-0.5 h-4 w-4 shrink-0 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
              </Link>
            ))}
          </div>
          <CodeBlock language="bash" code={`grit plugin add kanban`} />

          <h2 className="mb-2 mt-12 text-2xl font-semibold tracking-tight">Go module plugins</h2>
          <p className="mb-5 leading-relaxed text-muted-foreground">
            Larger pieces that ship as Go modules from{' '}
            <Link
              href="https://github.com/MUKE-coder/grit-plugins"
              className="text-primary hover:underline"
            >
              grit-plugins
            </Link>
            . Installed with <code>go get</code> rather than <code>grit plugin add</code>, because
            they are libraries your code imports rather than code generated into it.
          </p>
          <div className="mb-6 grid gap-2 sm:grid-cols-2">
            {MODULES.map((p) => (
              <div key={p.name} className="rounded-xl border border-border p-4">
                <code className="font-mono text-sm text-primary">{p.name}</code>
                <p className="mt-1 text-sm text-muted-foreground">{p.summary}</p>
              </div>
            ))}
          </div>
          <CodeBlock
            language="bash"
            code={`go get github.com/MUKE-coder/grit-plugins/grit-websockets`}
          />

          <h2 className="mb-2 mt-12 text-2xl font-semibold tracking-tight">Community</h2>
          <p className="mb-5 leading-relaxed text-muted-foreground">
            Written by other people and listed here on request.
          </p>

          <div className="mb-6 flex gap-3 rounded-xl border border-amber-500/30 bg-amber-500/5 p-5">
            <ShieldAlert className="h-5 w-5 shrink-0 text-amber-500" />
            <div className="leading-relaxed text-muted-foreground">
              <p className="mb-2 font-medium text-foreground">Listing is not endorsement</p>
              <p>
                The review is deliberately shallow: that the repository exists, that it installs and
                removes cleanly on a fresh project, and that it is not obviously hostile. Nobody
                audits the code for you. A plugin writes files into your project and patches
                existing ones, so read what it did with <code>git diff</code>, the same as you would
                any dependency.
              </p>
            </div>
          </div>

          {COMMUNITY.length === 0 ? (
            <div className="mb-6 rounded-xl border border-dashed border-border p-8 text-center">
              <p className="mb-2 font-medium">Nothing listed yet</p>
              <p className="text-sm text-muted-foreground">
                If you have written one, yours can be the first.
              </p>
            </div>
          ) : (
            <div className="mb-6 space-y-2">
              {COMMUNITY.map((p) => (
                <div key={p.name} className="rounded-xl border border-border p-4">
                  <div className="flex items-baseline gap-3">
                    <code className="font-mono text-sm text-primary">{p.name}</code>
                    <span className="text-xs text-muted-foreground">by {p.author}</span>
                    <span className="ml-auto text-xs text-muted-foreground">{p.added}</span>
                  </div>
                  <p className="mt-1 text-sm text-muted-foreground">{p.summary}</p>
                  <Link
                    href={p.repo.startsWith('http') ? p.repo : `https://${p.repo}`}
                    className="mt-2 inline-flex items-center gap-1.5 text-xs text-primary hover:underline"
                  >
                    <Github className="h-3 w-3" />
                    {p.repo}
                  </Link>
                </div>
              ))}
            </div>
          )}

          <p className="mb-3 leading-relaxed text-muted-foreground">
            Installing one takes the repository directly:
          </p>
          <CodeBlock language="bash" code={`grit plugin add github.com/someone/grit-calendar`} />
          <p className="mb-4 mt-3 leading-relaxed text-muted-foreground">
            Grit confirms the source before cloning, and nothing from the repository is executed: it
            reads the plugin&apos;s manifest and writes the files that manifest names.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Listing your own</h2>
          <CodeBlock language="bash" code={`grit plugin submit`} />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            That prints the checklist and the submission link. The checklist is three commands on a
            freshly scaffolded project, and the third is the one people skip:
          </p>
          <ol className="mb-6 ml-5 list-decimal space-y-1.5 leading-relaxed text-muted-foreground">
            <li>
              <code>grit plugin add &lt;your repo&gt;</code> installs without errors
            </li>
            <li>
              <code>go build ./...</code> compiles afterwards
            </li>
            <li>
              <code>grit plugin remove &lt;name&gt;</code> leaves the project compiling again
            </li>
          </ol>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Removal replays the lockfile backwards, so it only works if every change the plugin made
            is one the lockfile recorded. A plugin that edits a file outside the injection mechanism
            installs fine and leaves wreckage when it is taken out.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Accepted submissions are added to <code>internal/plugin/community.json</code>, which is
            what both <code>grit plugin list</code> and this page read.
          </p>

          <div className="mt-12 flex flex-col gap-3 sm:flex-row">
            <Link
              href="/docs/plugins/authoring"
              className="group flex flex-1 items-center gap-3 rounded-xl border border-border p-4 transition-colors hover:border-primary"
            >
              <span className="flex-1">
                <span className="block font-medium">Writing a plugin</span>
                <span className="text-sm text-muted-foreground">Anatomy, rules, testing</span>
              </span>
              <ArrowRight className="h-4 w-4 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
            </Link>
            <Link
              href="/docs/extending"
              className="group flex flex-1 items-center gap-3 rounded-xl border border-border p-4 transition-colors hover:border-primary"
            >
              <span className="flex-1">
                <span className="block font-medium">Is a plugin the right thing?</span>
                <span className="text-sm text-muted-foreground">Four mechanisms compared</span>
              </span>
              <ArrowRight className="h-4 w-4 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
            </Link>
          </div>
        </div>
      </main>
    </div>
  )
}
