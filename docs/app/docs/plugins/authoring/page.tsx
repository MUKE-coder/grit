import Link from 'next/link'
import { AlertTriangle, ArrowRight, Boxes } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Writing a plugin — Grit',
  description:
    'The anatomy of a Grit plugin, walked through the kanban plugin: files, injections, dependencies, preflight, and the lockfile that makes removal exact. The rules that matter and the failures each one prevents.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/authoring' },
}

const RULES = [
  {
    rule: 'Never write a file you did not create',
    why: 'Removal deletes every file the lockfile records. If a plugin overwrites a file the project already had, removal deletes the user’s file. The installer refuses to overwrite for exactly this reason: change an existing file with an injection instead.',
  },
  {
    rule: 'Every injection needs a marker that will still be there',
    why: 'An injection anchors on a comment like // grit:models. If a plugin anchors on ordinary code, a reformat or a rename silently breaks the install. Anchor on markers, and if the marker you need does not exist, add it to the scaffold rather than anchoring on something else.',
  },
  {
    rule: 'Declare dependencies, do not assume them',
    why: 'GoDeps and NodeDeps are fetched before any file is written, so a network failure aborts before the project changes. A plugin that imports a package it did not declare installs cleanly and then fails to compile.',
  },
  {
    rule: 'Refuse early rather than install half',
    why: 'Preflight runs before the first write. A plugin whose UI needs an admin app should say so on a project that has none, not write the Go half and report success while the feature is invisible.',
  },
  {
    rule: 'Adapt to the architecture',
    why: 'A --single project has no apps/api, and a TanStack admin keeps its code under src/ with a different router. Path helpers exist for both. Plugins that wrote the Next shape unconditionally landed their files where the alias does not point, which is not a build error: the plugin reports success and the feature is simply not there.',
  },
  {
    rule: 'Pin a column name GORM would not guess',
    why: 'GORM derives column names by looking for known initialisms anywhere in a field name. WIPLimit contains IP, so the column is w_ip_limit. Reading the struct works, AutoMigrate is happy, and only a write by column name fails, in a handler, in production.',
  },
  {
    rule: 'Test that it removes, not just that it installs',
    why: 'Installing is the easy half, and the half everybody tests. Removal replays the lockfile backwards, so it only works if every change the plugin made is one the lockfile recorded. A plugin that edits a file outside the injection mechanism installs fine and leaves wreckage behind.',
  },
]

export default function PluginAuthoringPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <Boxes className="h-5 w-5 text-accent" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Plugins
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">Writing a plugin</h1>
          <p className="mb-10 text-lg leading-relaxed text-muted-foreground">
            This guide is written around the{' '}
            <Link href="/docs/plugins/kanban" className="text-accent hover:underline">
              kanban plugin
            </Link>
            , because a board is the smallest feature that needs every part of the system at once:
            models and migrations, a service with real logic, routes, React Query hooks, two admin
            pages, a nav entry, an npm dependency, and injections into six existing files.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            What a plugin actually is
          </h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            A Grit plugin is not loaded at runtime. It is a code generator with a receipt. Install
            writes files, patches existing files at named markers, and records every one of those
            actions in <code>.grit/plugins.lock.json</code>. Removal reads that record and replays
            it backwards.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            That one design decision shapes everything else here. Grit previously had a feature
            that generated code one way and removed it via a separately maintained list. The two
            drifted, and removal left behind handlers and switch arms that no longer compiled. A
            plugin system multiplies that risk by every plugin anyone writes, so removal is derived
            from what install did, never from a second list an author has to remember to update.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">The anatomy</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            A plugin is one value. Here is kanban&apos;s, trimmed:
          </p>
          <CodeBlock
            language="go"
            code={`func kanbanPlugin() Plugin {
	return Plugin{
		Name:    "kanban",
		Version: "1.0.0",
		Summary: "Boards, columns and cards with drag-and-drop ordering",

		// Refuses before the first write if the project cannot host it.
		Preflight: kanbanPreflight,

		// Fetched before any file is written, so a network failure aborts
		// while the project is still untouched.
		NodeDeps: []Dependency{
			{Name: "@dnd-kit/core", Version: "6.3.1", Workspace: "apps/admin"},
			{Name: "@dnd-kit/sortable", Version: "10.0.0", Workspace: "apps/admin"},
		},

		// Functions, not data: a triple project needs admin pages an --api
		// project does not.
		Files:      kanbanFiles,
		Injections: kanbanInjections,

		NextSteps: []string{
			"Install the new admin dependency:  pnpm install",
			"Run migrations:                    cd apps/api && go run cmd/migrate/main.go",
		},
	}
}`}
          />

          <h3 className="mb-3 mt-8 text-xl font-semibold tracking-tight">Files</h3>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            <code>Files</code> returns a map of project-relative path to content. It is a function
            because the right paths depend on the project: <code>apiPath</code> resolves{' '}
            <code>apps/api</code> or the root, <code>adminFile</code> resolves the Next admin or the
            TanStack one, and <code>adminSource</code> converts a component written in the Next
            dialect for whichever admin the project has.
          </p>
          <CodeBlock
            language="go"
            code={`func kanbanFiles(ctx Context) map[string]string {
	files := map[string]string{
		apiPath(ctx, "internal/kanban/rank.go"):         kanbanRankSource,
		apiPath(ctx, "internal/models/board.go"):        kanbanModelsGo(ctx),
		apiPath(ctx, "internal/services/board.go"):      kanbanServiceGo(ctx),
		apiPath(ctx, "internal/handlers/board.go"):      kanbanHandlerGo(ctx),
		apiPath(ctx, "internal/handlers/board_test.go"): kanbanHandlerTestGo(ctx),
	}

	if !hasAdmin(ctx) {
		return files
	}

	files[adminFile(ctx, "hooks/use-boards.ts")] = adminSource(ctx, kanbanHook())
	files[adminFile(ctx, "components/kanban/board-view.tsx")] = adminSource(ctx, kanbanBoardView())

	// A Next page is its own route; TanStack needs a separate shim.
	list, listRoute := adminPageFile(ctx, "boards", "boards")
	files[list] = adminSource(ctx, kanbanBoardsPage())
	if listRoute != "" {
		files[listRoute] = tanStackRouteShim("boards", "@/pages/boards")
	}

	return files
}`}
          />

          <div className="my-6 rounded-xl border border-accent/30 bg-accent/5 p-5">
            <p className="mb-2 font-medium">Embed code you want tested</p>
            <p className="leading-relaxed text-muted-foreground">
              <code>kanbanRankSource</code> is not a Go string. It is{' '}
              <code>//go:embed kanban/rank.go</code>, so the ordering algorithm the project receives
              is byte for byte the code Grit&apos;s own test suite runs on every CI run. A plugin
              that keeps its algorithm in a string literal has no tests for what it ships: the
              string compiles nowhere until somebody installs it.
            </p>
          </div>

          <h3 className="mb-3 mt-8 text-xl font-semibold tracking-tight">Injections</h3>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            An injection inserts a snippet on the line before a marker comment. The exact text is
            recorded in the lockfile, so removal takes out precisely what was put in.
          </p>
          <CodeBlock
            language="go"
            code={`{
	File:   apiPath(ctx, "internal/models/user.go"),
	Marker: "// grit:models",
	Code:   "\\t\\t&Board{},\\n\\t\\t&BoardColumn{},\\n\\t\\t&BoardCard{},",
},
{
	File:   apiPath(ctx, "internal/routes/routes.go"),
	Marker: "// grit:handlers",
	Code:   "\\tboardHandler := handlers.NewBoardHandler(db)",
},`}
          />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            The markers available in a scaffolded project:
          </p>
          <div className="mb-6 overflow-hidden rounded-xl border border-border">
            <table className="w-full text-sm">
              <thead className="bg-card/60">
                <tr className="border-b border-border text-left">
                  <th className="px-4 py-3 font-medium">Marker</th>
                  <th className="px-4 py-3 font-medium">In</th>
                  <th className="px-4 py-3 font-medium">For</th>
                </tr>
              </thead>
              <tbody className="align-top">
                <tr className="border-b border-border/60">
                  <td className="px-4 py-3 font-mono text-xs">// grit:models</td>
                  <td className="px-4 py-3 text-muted-foreground">models/user.go</td>
                  <td className="px-4 py-3 text-muted-foreground">Register a model for AutoMigrate</td>
                </tr>
                <tr className="border-b border-border/60">
                  <td className="px-4 py-3 font-mono text-xs">// grit:handlers</td>
                  <td className="px-4 py-3 text-muted-foreground">routes/routes.go</td>
                  <td className="px-4 py-3 text-muted-foreground">Construct a handler</td>
                </tr>
                <tr className="border-b border-border/60">
                  <td className="px-4 py-3 font-mono text-xs">// grit:routes:protected</td>
                  <td className="px-4 py-3 text-muted-foreground">routes/routes.go</td>
                  <td className="px-4 py-3 text-muted-foreground">Routes behind auth</td>
                </tr>
                <tr className="border-b border-border/60">
                  <td className="px-4 py-3 font-mono text-xs">// grit:routes:custom</td>
                  <td className="px-4 py-3 text-muted-foreground">routes/routes.go</td>
                  <td className="px-4 py-3 text-muted-foreground">
                    Routes with no auth middleware, for genuinely anonymous endpoints
                  </td>
                </tr>
                <tr className="border-b border-border/60">
                  <td className="px-4 py-3 font-mono text-xs">// grit:icons:import</td>
                  <td className="px-4 py-3 text-muted-foreground">admin lib/icons.ts</td>
                  <td className="px-4 py-3 text-muted-foreground">Import a lucide icon</td>
                </tr>
                <tr className="border-b border-border/60">
                  <td className="px-4 py-3 font-mono text-xs">// grit:icons:map</td>
                  <td className="px-4 py-3 text-muted-foreground">admin lib/icons.ts</td>
                  <td className="px-4 py-3 text-muted-foreground">
                    Make it reachable by name. Both are needed: getIcon falls back to a document
                    icon and nothing reports the miss
                  </td>
                </tr>
                <tr className="border-b border-border/60">
                  <td className="px-4 py-3 font-mono text-xs">// grit:nav:content</td>
                  <td className="px-4 py-3 text-muted-foreground">admin CollapsibleSidebar.tsx</td>
                  <td className="px-4 py-3 text-muted-foreground">
                    A working surface, in the rail beside the resources
                  </td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-mono text-xs">// grit:nav:system</td>
                  <td className="px-4 py-3 text-muted-foreground">admin CollapsibleSidebar.tsx</td>
                  <td className="px-4 py-3 text-muted-foreground">
                    An operational surface, in the System Hub&apos;s Extensions tab
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Those last two are worth getting right, and kanban got it wrong first. A board filed
            under System, Extensions is three clicks from the dashboard and nobody looks for their
            board there. <code>grit:nav:content</code> exists because building this plugin found
            that there was no anchor for a surface people open while working rather than while
            configuring.
          </p>

          <h3 className="mb-3 mt-8 text-xl font-semibold tracking-tight">Preflight</h3>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Preflight states a requirement the plugin cannot inject its way out of. It runs before
            anything is written, so a project that is not ready is told what to do rather than left
            half-installed.
          </p>
          <CodeBlock
            language="go"
            code={`func kanbanPreflight(ctx Context) error {
	if !hasAdmin(ctx) {
		return fmt.Errorf("kanban needs an admin app, and this project is --%s\\n\\n"+
			"The board UI lives in the admin. Scaffold with --triple, or install the\\n"+
			"API half by hand from the files this plugin would have written.", ctx.Architecture)
	}
	return nil
}`}
          />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            Say what to do instead. &quot;Incompatible&quot; is not a fix; &quot;run grit upgrade
            first&quot; is.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">The rules</h2>
          <p className="mb-6 leading-relaxed text-muted-foreground">
            Each of these exists because something went wrong without it.
          </p>
          <ol className="mb-8 space-y-5">
            {RULES.map((r, i) => (
              <li key={r.rule} className="rounded-xl border border-border bg-card/40 p-5">
                <p className="mb-2 font-medium">
                  <span className="mr-2 font-mono text-sm text-accent">{i + 1}</span>
                  {r.rule}
                </p>
                <p className="leading-relaxed text-muted-foreground">{r.why}</p>
              </li>
            ))}
          </ol>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            Writing one outside the Grit repo
          </h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            The plugins above are Go values compiled into the CLI. You do not need to fork Grit to
            write one: the same shape is expressible as a directory with a JSON manifest, and
            install and removal behave identically.
          </p>
          <CodeBlock
            language="text"
            code={`my-plugin/
  plugin.json          the manifest
  files/               file bodies, referenced by "from"`}
          />
          <CodeBlock
            language="json"
            code={`{
  "name": "calendar",
  "version": "1.0.0",
  "summary": "A month view for any resource with a date field",
  "nodeDeps": [
    { "name": "date-fns", "version": "4.1.0", "workspace": "apps/admin" }
  ],
  "files": [
    { "from": "files/calendar.tsx", "to": "apps/admin/components/calendar.tsx" }
  ],
  "injections": [
    {
      "file": "apps/admin/components/chrome/CollapsibleSidebar.tsx",
      "marker": "// grit:nav:content",
      "code": "  { href: \\"/calendar\\", label: \\"Calendar\\", iconKey: \\"Calendar\\", adminOnly: false },"
    }
  ],
  "nextSteps": ["pnpm install"]
}`}
          />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            File bodies get <code>{'{{MODULE}}'}</code> and <code>{'{{PROJECT}}'}</code> substituted,
            so a plugin does not need to know the module path of the project installing it. Install
            it from a path or straight from a repository:
          </p>
          <CodeBlock
            language="bash"
            code={`grit plugin add ./my-plugin
grit plugin add github.com/you/grit-calendar`}
          />

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Testing it</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            On a freshly scaffolded project, in this order:
          </p>
          <CodeBlock
            language="bash"
            code={`grit new testbed --triple --next
cd testbed

grit plugin add ./my-plugin       # installs without errors
cd apps/api && go build ./...     # still compiles
go test ./internal/...            # the tests it shipped pass
cd ../.. && pnpm install && cd apps/admin && npx tsc --noEmit

cd ../.. && grit plugin remove calendar --force
cd apps/api && go build ./...     # STILL compiles`}
          />
          <div className="my-6 flex gap-3 rounded-xl border border-amber-500/30 bg-amber-500/5 p-5">
            <AlertTriangle className="h-5 w-5 shrink-0 text-amber-500" />
            <p className="leading-relaxed text-muted-foreground">
              That last build is the one people skip, and it is the one that matters. Installing is
              the easy half. If removal leaves a project that does not compile, the plugin changed
              something the lockfile did not record.
            </p>
          </div>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Publishing it</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Put it in a public repository with the manifest at the root or in <code>plugin/</code>,
            then run <code>grit plugin submit</code> for the checklist and the submission link.
            Accepted plugins appear in <code>grit plugin list</code> and on the{' '}
            <Link href="/docs/plugins/marketplace" className="text-accent hover:underline">
              marketplace
            </Link>
            .
          </p>

          <div className="mt-12 flex flex-col gap-3 sm:flex-row">
            <Link
              href="/docs/plugins/marketplace"
              className="group flex flex-1 items-center gap-3 rounded-xl border border-border p-4 transition-colors hover:border-accent"
            >
              <span className="flex-1">
                <span className="block font-medium">Marketplace</span>
                <span className="text-sm text-muted-foreground">Official and community plugins</span>
              </span>
              <ArrowRight className="h-4 w-4 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
            </Link>
            <Link
              href="/docs/plugins/kanban"
              className="group flex flex-1 items-center gap-3 rounded-xl border border-border p-4 transition-colors hover:border-accent"
            >
              <span className="flex-1">
                <span className="block font-medium">The kanban plugin</span>
                <span className="text-sm text-muted-foreground">The example in full</span>
              </span>
              <ArrowRight className="h-4 w-4 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
            </Link>
          </div>
        </div>
      </main>
    </div>
  )
}
