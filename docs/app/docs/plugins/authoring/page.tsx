import Link from 'next/link'
import { ArrowLeft, ArrowRight, Boxes, AlertTriangle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Plugins — Grit',
  description:
    'How Grit plugins work: installing, removing, and writing your own. Plugins generate code into your project rather than hiding behind a runtime dependency.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/authoring' },
}

export default function PluginsPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <Boxes className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Extending Grit
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">Plugins</h1>
          <p className="mb-10 text-lg text-muted-foreground">
            A Grit plugin <strong>generates code into your project</strong>. It isn&apos;t a
            runtime dependency you import &mdash; it writes models, routes and pages into
            your repo, which you then own and can edit like anything else.
          </p>

          <div className="prose-grit">
            <h2>Using plugins</h2>
            <CodeBlock language="bash" code={`grit plugin list              # what's available
grit plugin info multitenant  # what it does
grit plugin add multitenant   # install
grit plugin remove multitenant`} />

            <h2>Why generated code, not a library</h2>
            <p>
              Grit is a code generator. <code>grit new</code> emits a Go + React app that has
              no dependency on Grit afterwards &mdash; there is no runtime framework object
              for a plugin to hook into. So plugins write code instead.
            </p>
            <p>That trade is deliberate, and it cuts both ways:</p>
            <ul>
              <li>
                <strong>You own the code.</strong> Read it, edit it, delete it. No digging
                through a vendor directory to work out what a plugin does.
              </li>
              <li>
                <strong>No automatic upgrades.</strong> A plugin that ships a fix
                won&apos;t update code already in your repo. Remove and re-add to take a new
                version, and review the diff.
              </li>
            </ul>

            <h2>Removal is exact</h2>
            <p>
              Installation records every file written and every snippet injected in{' '}
              <code>.grit/plugins.lock.json</code>. Removal replays that record backwards.
              Commit the lockfile &mdash; it&apos;s how your teammates know what&apos;s
              installed.
            </p>
            <p>
              This is the part worth understanding if you write a plugin: there is{' '}
              <strong>no uninstall code to write</strong>. You describe what to install;
              removal is derived. A separate hand-maintained removal list is exactly how
              this kind of tooling drifts and starts leaving projects that don&apos;t
              compile.
            </p>

            <div className="my-8 rounded-xl border border-amber-500/30 bg-amber-500/5 p-5">
              <div className="mb-2 flex items-center gap-2">
                <AlertTriangle className="h-4 w-4 text-amber-500" />
                <strong className="text-foreground">Edited code is never overwritten</strong>
              </div>
              <p className="mb-0 text-sm text-muted-foreground">
                If you change an injected block, removal reports it and leaves it alone
                rather than guessing. You&apos;ll be told which file to tidy by hand.
              </p>
            </div>

            <h2>Writing a plugin</h2>
            <p>A plugin is a value describing files, injections and dependencies:</p>
            <CodeBlock language="go" code={`plugin.Plugin{
    Name:    "audit-trail",
    Version: "1.0.0",
    Summary: "Records who changed what",

    // Requires other plugins be installed first.
    Requires: []string{},

    Files: func(ctx plugin.Context) map[string]string {
        return map[string]string{
            "apps/api/internal/audit/audit.go": auditSource(ctx),
        }
    },

    Injections: func(ctx plugin.Context) []plugin.Injection {
        return []plugin.Injection{{
            File:   "apps/api/internal/routes/routes.go",
            Marker: "// grit:routes:protected",
            Code:   "\\t\\tprotected.GET(\\"/audit\\", auditHandler.List)",
        }}
    },

    NextSteps: []string{"Run migrations: go run cmd/migrate/main.go"},
}`} />

            <p>
              <code>Files</code> and <code>Injections</code> are functions, not data, so a
              plugin can adapt to the project. <code>ctx</code> carries the module path,
              architecture and frontend &mdash; a <code>--triple</code> app needs admin
              pages an <code>--api</code> project does not.
            </p>

            <h3>Injections</h3>
            <p>
              Code is inserted on the line <em>before</em> a marker comment. The markers a
              scaffolded project provides:
            </p>
            <CodeBlock language="text" code={`// grit:models             AutoMigrate registry
// grit:handlers           handler construction
// grit:routes:protected   authenticated routes
// grit:routes:admin       admin-only routes
// grit:routes:custom      public routes
// grit:seeders            seed registration
// grit:cron-tasks         scheduled tasks`} />

            <p>
              Mark an injection <code>Optional</code> when the target legitimately may not
              exist &mdash; a frontend file in an <code>--api</code> project. A missing
              marker is then a warning instead of a failed install.
            </p>

            <h3>Rules the installer enforces</h3>
            <ul>
              <li>
                <strong>Never overwrites an existing file.</strong> Removal would otherwise
                delete something your plugin didn&apos;t create.
              </li>
              <li>
                <strong>A missing marker fails the install</strong> unless the injection is
                optional. A silent partial install is worse than a clear error.
              </li>
              <li>
                <strong>Installing twice is refused</strong> &mdash; remove first.
              </li>
              <li>
                <strong>Requirements are enforced both ways.</strong> You can&apos;t install
                without dependencies, or remove something still depended on.
              </li>
            </ul>

            <h2>A worked example</h2>
            <p>
              The <code>multitenant</code> plugin is the reference implementation: models, a
              GORM callback, middleware, an API, and tests &mdash; five files and three
              injections. Read it in <code>internal/plugin/multitenant.go</code>.
            </p>
          </div>

          <div className="mt-12 flex items-center justify-between border-t border-border/40 pt-6">
            <Button asChild variant="ghost">
              <Link href="/docs">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Docs
              </Link>
            </Button>
            <Button asChild variant="ghost">
              <Link href="/docs/plugins/multitenant">
                Multi-tenancy
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
