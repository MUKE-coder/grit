import Link from 'next/link'
import { ArrowLeft, ArrowRight, Command, Search, Keyboard } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Command palette plugin — Grit',
  description:
    'A keyboard-first ⌘K command palette for the admin, built from the resource registry — pure client code, no Go.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/command-palette' },
}

export default function CommandPalettePage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <Command className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Plugin
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
            Command palette
          </h1>
          <p className="mb-10 text-lg text-muted-foreground">
            Press ⌘K, type a few letters, hit Enter &mdash; jump to any resource or
            system page without touching the mouse.
          </p>

          <div className="prose-grit">
            <CodeBlock language="bash" code={`grit plugin add command-palette`} />

            <p>
              This one is unusual: it adds a full feature to the admin and never touches
              Go. No models, no migrations, no server changes &mdash; it is pure client
              code, mounted through the admin layout&apos;s injection markers. It works on
              the <code>triple</code> and <code>full</code> architectures, the two that ship
              an admin app.
            </p>

            <h2>The problem it solves</h2>
            <p>
              Navigating a growing admin by mouse &mdash; sidebar, then section, then page
              &mdash; is slow for people who live in it all day. A keyboard palette lets a
              power user jump anywhere by typing a few letters, the same muscle memory they
              already have from Linear, VS Code, or Raycast.
            </p>
            <p>
              The admin already ships a floating <strong>QuickAccess</strong> button &mdash;
              a grid of cards you click. The command palette is the complementary keyboard
              path, built from the same registry so the two never drift apart.
            </p>

            <h2>How it works</h2>
            <p>
              A global <code>keydown</code> listener is registered once when the admin
              loads. Press <code>⌘K</code> on macOS or <code>Ctrl+K</code> anywhere to open
              a centered overlay with a search box; <code>Escape</code> closes it. Until you
              open it, the component renders nothing.
            </p>
            <p>
              The list of destinations is assembled from two sources. The first is the
              resource registry &mdash; <code>resources</code> imported from{' '}
              <code>@/resources</code>. Every registered resource contributes two commands:
            </p>
            <CodeBlock language="text" code={`Go to <Plural>    →  /resources/<slug>
New <Singular>    →  /resources/<slug>?action=create`} />
            <p>
              The second source is the fixed set of system pages: Dashboard, System Hub,
              Roles &amp; permissions, User Activity, System Health, Data &amp; Backup, and
              Profile.
            </p>
            <p>
              Because the palette reads the registry directly, a resource you just created
              with <code>grit generate resource</code> shows up automatically &mdash; there
              is nothing to wire up.
            </p>

            <div className="my-8 rounded-xl border border-primary/30 bg-primary/5 p-5">
              <div className="mb-2 flex items-center gap-2">
                <Search className="h-4 w-4 text-primary" />
                <strong className="text-foreground">Fuzzy search and keyboard control</strong>
              </div>
              <p className="mb-0 text-sm text-muted-foreground">
                Type to fuzzy-filter against each command&apos;s label and hint. Arrow Up and
                Arrow Down move the highlight, Enter opens the highlighted command, and
                clicking a row opens it too. Navigation goes through Next.js{' '}
                <code>router.push</code>, so every jump is client-side and instant.
              </p>
            </div>

            <h2>Install</h2>
            <p>
              Add the plugin, then start the admin &mdash; there is no migration step because
              there is no database work:
            </p>
            <CodeBlock language="bash" code={`grit plugin add command-palette
cd apps/admin && pnpm dev`} />
            <p>
              The palette is mounted in the dashboard layout via injection markers, so it is
              available on every admin page the moment the app boots.
            </p>

            <h2>Using it</h2>
            <p>
              Press <code>Ctrl+K</code> (or <code>⌘K</code>), start typing a resource or page
              name, and press <code>Enter</code> to go there:
            </p>
            <CodeBlock language="text" code={`⌘K            open the palette
roles         filter to "Roles & permissions"
Enter         jump to it

⌘K
new           see every "New <Resource>" create action
↑ / ↓         move the highlight
Enter         open the highlighted command`} />

            <h2>Use cases</h2>
            <ul>
              <li>
                <strong>Fast keyboard navigation</strong> for people in the admin every day
                &mdash; no reaching for the sidebar.
              </li>
              <li>
                <strong>Jump straight to create</strong> &mdash; type <code>new</code> and a
                resource name to land on the create form without clicking through.
              </li>
              <li>
                <strong>Discoverability</strong> &mdash; a single searchable list of every
                place the admin can take you.
              </li>
              <li>
                <strong>Muscle memory</strong> &mdash; the ⌘K reflex you already have from
                Linear, VS Code, and Raycast, now in your own admin.
              </li>
            </ul>

            <div className="my-8 rounded-xl border border-primary/30 bg-primary/5 p-5">
              <div className="mb-2 flex items-center gap-2">
                <Keyboard className="h-4 w-4 text-primary" />
                <strong className="text-foreground">A plugin can be pure frontend</strong>
              </div>
              <p className="mb-0 text-sm text-muted-foreground">
                Most plugins touch the Go API. This one proves the injection-marker system
                is general enough for client-only features: the command palette is entirely
                admin code and still installs, mounts, and stays in sync with your resources
                like any other plugin.
              </p>
            </div>
          </div>

          <div className="mt-12 flex items-center justify-between border-t border-border/40 pt-6">
            <Button asChild variant="ghost">
              <Link href="/docs/plugins/impersonate">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Impersonate
              </Link>
            </Button>
            <Button asChild variant="ghost">
              <Link href="/docs/plugins/saved-views">
                Saved views
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
