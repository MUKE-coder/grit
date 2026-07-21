import Link from 'next/link'
import { ArrowLeft, ArrowRight, Puzzle, RotateCcw, PackageCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Plugins — Grit',
  description:
    'What Grit plugins are: reversible code generation recorded in a lockfile. The plugin model, the CLI, the first-party catalog, and how removal is derived rather than written.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/overview' },
}

export default function PluginsOverviewPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <Puzzle className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Plugins
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
            Plugins
          </h1>
          <p className="mb-10 text-lg leading-relaxed text-muted-foreground">
            A Grit plugin generates real code into your repo — models, routes, pages, migrations —
            not a runtime dependency that hides in <code>node_modules</code>. And because every
            install is recorded, plugins do the one thing most plugin systems can&apos;t: cleanly
            uninstall.
          </p>

          {/* What a plugin is */}
          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">What a plugin is</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Most frameworks extend through packages: you install one, it lives in your dependency
            folder, and at runtime it reaches in and registers itself. That&apos;s fine until you
            want to <em>see</em> what it did or <em>remove</em> it — the code isn&apos;t yours, it&apos;s
            behind a version number.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Grit plugins work the other way around. <code>grit plugin add multitenant</code> writes
            Go models, handlers, routes, admin pages and migrations directly into{' '}
            <code>apps/api</code>, <code>apps/admin</code> and the rest of your project. After it
            runs, the code is <strong>yours</strong>: read it, edit it, step through it in a
            debugger. There&apos;s no hidden runtime library deciding your app&apos;s behavior.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            This is the same idea as <Link href="/docs/concepts/code-generation" className="text-primary hover:underline">code generation</Link>:{' '}
            <code>grit generate resource</code> gives you a model, a service, a handler and an admin
            page you own. A plugin is that, bigger and reusable — a packaged code generation.
          </p>

          {/* The lockfile */}
          <div className="my-8 rounded-xl border border-border bg-card/40 p-6">
            <div className="mb-2 flex items-center gap-2">
              <RotateCcw className="h-4 w-4 text-primary" />
              <span className="font-semibold">The part most systems can&apos;t do: uninstall</span>
            </div>
            <p className="mb-3 text-sm leading-relaxed text-muted-foreground">
              When a plugin runs, Grit records <strong>everything</strong> it did — every file it
              created and every injection it made — into a lockfile:
            </p>
            <CodeBlock language="text" code={`.grit/plugins.lock.json`} />
            <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
              So removal isn&apos;t a second pile of code the author maintains — it&apos;s derived.{' '}
              <code>grit plugin remove</code> reads the lockfile and replays the install backwards:
              it deletes the files it added and reverts the injections, in reverse order. Commit the
              lockfile, and your plugins are as reversible as a git branch. A framework whose plugins
              publish migrations into your app can&apos;t cleanly reverse them — publishing is a
              one-way door. Grit&apos;s can, because the install is a recorded, replayable
              transaction.
            </p>
          </div>

          {/* CLI */}
          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">The commands</h2>
          <CodeBlock language="bash" code={`grit plugin list                 # what's available
grit plugin info impersonate     # what a plugin does before you run it
grit plugin add impersonate      # install (refuses to run on a dirty git tree)
grit plugin remove impersonate   # replay the install backwards`} />
          <p className="mt-4 leading-relaxed text-muted-foreground">
            <code>grit plugin add</code> refuses to run with uncommitted changes, so the diff it
            produces is always clean and reviewable. Read the diff, commit it, and it&apos;s part of
            your app like anything else.
          </p>

          {/* Guarantees */}
          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            What the installer guarantees
          </h2>
          <ul className="mb-4 space-y-3 leading-relaxed text-muted-foreground">
            <li className="flex gap-3">
              <PackageCheck className="mt-0.5 h-4 w-4 shrink-0 text-primary" />
              <span><strong>No overwrite.</strong> A plugin can&apos;t clobber a file that already exists — it creates, or it fails loudly.</span>
            </li>
            <li className="flex gap-3">
              <PackageCheck className="mt-0.5 h-4 w-4 shrink-0 text-primary" />
              <span><strong>A missing marker is an error.</strong> If an injection&apos;s anchor isn&apos;t found, the install stops rather than silently doing nothing.</span>
            </li>
            <li className="flex gap-3">
              <PackageCheck className="mt-0.5 h-4 w-4 shrink-0 text-primary" />
              <span><strong>No double install.</strong> Installing an already-installed plugin is refused.</span>
            </li>
            <li className="flex gap-3">
              <PackageCheck className="mt-0.5 h-4 w-4 shrink-0 text-primary" />
              <span><strong>Requirements both ways.</strong> A plugin&apos;s dependencies must be present to install it, and can&apos;t be removed while it&apos;s still installed.</span>
            </li>
          </ul>

          {/* The catalog */}
          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            First-party plugins
          </h2>
          <p className="mb-6 leading-relaxed text-muted-foreground">
            The built-in catalog covers the big, opinionated pieces that don&apos;t belong in every
            app, and the small high-value features that the best admin ecosystems treat as add-ons.
          </p>
          <div className="grid gap-3 sm:grid-cols-2">
            {[
              { href: '/docs/plugins/multitenant', name: 'multitenant', desc: 'Organizations, per-org roles, and automatic query scoping that fails closed.' },
              { href: '/docs/plugins/impersonate', name: 'impersonate', desc: 'An admin signs in as another user, with an audit trail and one-click return.' },
              { href: '/docs/plugins/command-palette', name: 'command-palette', desc: '⌘K navigation across the admin. Frontend-only — touches no Go.' },
              { href: '/docs/plugins/saved-views', name: 'saved-views', desc: "Save a table's filters and sort as named, per-user views." },
            ].map((p) => (
              <Link
                key={p.href}
                href={p.href}
                className="group rounded-xl border border-border bg-card/40 p-4 transition-colors hover:border-primary/40"
              >
                <div className="mb-1 font-mono text-sm font-semibold text-foreground group-hover:text-primary">
                  {p.name}
                </div>
                <p className="text-sm leading-relaxed text-muted-foreground">{p.desc}</p>
              </Link>
            ))}
          </div>

          <p className="mt-8 leading-relaxed text-muted-foreground">
            Want to build one? See{' '}
            <Link href="/docs/plugins/authoring" className="text-primary hover:underline">
              Writing a plugin
            </Link>
            . Plugins are how Grit stays small at the core and still grows: the batteries every app
            needs are built in, and the bigger, opinionated pieces are plugins — real code, in your
            repo, that you can always take back out.
          </p>

          {/* Nav */}
          <div className="mt-14 flex items-center justify-between border-t border-border pt-6">
            <Button variant="ghost" asChild>
              <Link href="/docs/batteries/modules">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Turning modules off
              </Link>
            </Button>
            <Button variant="ghost" asChild>
              <Link href="/docs/plugins/multitenant">
                Multi-tenancy plugin
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
