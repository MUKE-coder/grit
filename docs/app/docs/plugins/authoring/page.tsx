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
            <CodeBlock language="bash" code={`grit plugin list                       # what's available
grit plugin info multitenant           # what it does
grit plugin add multitenant            # install a built-in
grit plugin add ./plugins/my-plugin    # install one of your own
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
            <CodeBlock language="text" code={`// grit:models                AutoMigrate registry
// grit:handlers              handler construction
// grit:imports               imports in routes.go
// grit:middleware:protected  request middleware, after auth
// grit:routes:public         API-key routes
// grit:routes:protected      authenticated routes
// grit:routes:admin          admin-only routes
// grit:routes:custom         unauthenticated routes
// grit:seeders               seed registration
// grit:cron-tasks            scheduled tasks
// grit:nav:system            a System entry in the admin sidebar
// grit:icons:import          an icon import in the admin
// grit:icons:map             that icon in the admin's iconMap`} />

              <p className="mt-4 leading-relaxed text-muted-foreground">
                An icon needs both of its markers. <code>getIcon</code> falls back to
                a document icon for a key it does not know rather than failing, so a
                sidebar entry with only the import renders the wrong icon and nothing
                reports it.
              </p>

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

            <h2>Writing one you can install</h2>
            <p className="leading-relaxed text-muted-foreground">
              The value above is how the built-in plugins are written, and they are
              compiled into the CLI. Yours does not have to be: the same shape,
              described in JSON, installs from a directory.
            </p>

            <CodeBlock language="text" code={`my-plugin/
  plugin.json          the manifest
  files/               the file bodies it installs`} />

            <CodeBlock language="json" code={`{
  "name": "product-reviews",
  "version": "1.0.0",
  "summary": "Customer reviews with a moderation queue",
  "files": [
    { "from": "files/review_model.go", "to": "{{API_ROOT}}/internal/models/review.go" }
  ],
  "injections": [
    {
      "file": "{{API_ROOT}}/internal/models/user.go",
      "marker": "// grit:models",
      "code": "\t\t&Review{},"
    },
    {
      "file": "{{API_ROOT}}/internal/routes/routes.go",
      "marker": "// grit:routes:public",
      "code": "\t\tpublicAPI.GET(\"/products/:key/reviews\", reviewHandler.ForProduct)",
      "when": { "architecture": ["triple", "double", "full"] }
    }
  ],
  "nextSteps": ["Run the migration:  grit migrate"]
}`} />

              <p className="leading-relaxed text-muted-foreground">
                <code>from</code> is relative to the plugin directory,{' '}
                <code>to</code> to the project root. File bodies and paths get{' '}
                <code>{'{{MODULE}}'}</code>, <code>{'{{PROJECT}}'}</code> and{' '}
                <code>{'{{API_ROOT}}'}</code> substituted, so a plugin does not need to
                know the module path of the project installing it, or whether that
                project keeps its Go code in <code>apps/api</code> or at the root.
              </p>
              <p className="leading-relaxed text-muted-foreground">
                <code>when</code> limits a file or injection to certain projects, by{' '}
                <code>architecture</code> and <code>frontend</code>. It is data rather
                than an expression language on purpose: a plugin that needs real logic
                is better written in Go and contributed upstream, and a
                half-implemented expression parser is a support burden nobody asked
                for.
              </p>

              <div className="not-prose my-5 rounded-xl border border-border p-4 text-sm leading-relaxed">
                <p className="text-foreground"><strong>Nothing else changes.</strong></p>
                <p className="mt-2 text-muted-foreground">
                  A directory plugin is recorded in{' '}
                  <code>.grit/plugins.lock.json</code> exactly like a built-in, and{' '}
                  <code>grit plugin remove</code> replays that record backwards the same
                  way. The installer refuses the same things: an existing file, a
                  missing marker that is not optional, a second install, a dirty git
                  tree. Paths that climb out of the project are refused at load time,
                  because a plugin is somebody else&apos;s code and{' '}
                  <code>&quot;to&quot;: &quot;../../.ssh/authorized_keys&quot;</code>{' '}
                  is a file write rather than an install.
                </p>
              </div>

              <p className="leading-relaxed text-muted-foreground">
                Only an explicit path is read as a directory:{' '}
                <code>./plugins/x</code>, <code>../x</code> or an absolute path. A bare
                name is always a built-in, so a folder cannot shadow one by sharing its
                name.
              </p>

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
