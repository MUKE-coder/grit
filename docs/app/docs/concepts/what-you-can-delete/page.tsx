import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/what-you-can-delete')

/**
 * What you can take out of a generated project (grit#248).
 *
 * Every number on this page was measured on a `grit new myapp --triple --next`
 * at v3.568.0, not estimated. The import analysis in particular: "nothing
 * imports it" is a fact you can check with grep, and it is the only honest
 * basis for telling somebody a directory is safe to delete.
 *
 * The thing this page exists to correct is the reading of the file count as
 * "files I now maintain". Most of them are not that, and saying which are
 * which is more useful than making the number smaller.
 */

/**
 * Measured on `grit new measured --triple --next` at v3.569.0, with nothing
 * generated into it. Not estimated, and not taken from an older project: the
 * first draft of this page counted one that already had a resource in it and
 * an empty internal/docs that this release stopped writing.
 */
const COUNTS = {
  total: 819,
  go: 441,
  goTests: 124,
  adminTs: 242,
  webTs: 34,
}

type Row = { name: string; lines: number; what: string }

/**
 * Go packages nothing imports, in a project that has generated a resource.
 *
 * That condition is the whole subtlety. On a project with no resources of its
 * own, ten packages look unimported: pdf, rowrules, export, factory and
 * testkit join this list, and every one of them is used the moment you run
 * `grit generate resource`. Listing those as deletable would be the worst
 * kind of wrong, so the measurement is taken after generating one.
 *
 * Checked by grepping every import in apps/api against each package, not by
 * reading names, and then by deleting all six and running `go build ./...`
 * and `go test ./...`, which passed.
 */
const WAITING: Row[] = [
  { name: 'money', lines: 480, what: 'exact money: integer minor units and a currency. Used by a money field.' },
  { name: 'jsontime', lines: 295, what: 'the date and datetime types that survive a browser date picker. Used by a date or datetime field.' },
  { name: 'docx', lines: 1077, what: 'Word documents. Used by the document generation a resource opts into.' },
  { name: 'stock', lines: 448, what: 'stock levels and movements, for a resource that tracks them.' },
  { name: 'workflow', lines: 387, what: 'status state machines, for a resource with a workflow.' },
  { name: 'notify', lines: 392, what: 'the notification fan-out, for code that sends one.' },
]

/** Imported from everywhere. Deleting one is deleting the application. */
const LOAD_BEARING = [
  'handlers', 'services', 'models', 'routes', 'respond', 'middleware',
  'authz', 'config', 'database', 'migrate', 'paginate', 'ids',
]

export default function WhatYouCanDeletePage() {
  const waitingLines = WAITING.reduce((total, row) => total + row.lines, 0)

  return (
    <DocsShell>
      <div className="mb-10">
        <span className="tag-mono text-accent/80 mb-3 block">Concepts</span>
        <h1 className="text-4xl font-bold tracking-tight mb-4">What you can delete</h1>
        <p className="text-lg text-muted-foreground leading-relaxed">
A new triple project is {COUNTS.total} files, and the common reaction to that
          number is &quot;{COUNTS.total} files I now maintain&quot;. Most of them are not
          that: {COUNTS.goTests} are tests that run themselves, and the largest group is
          code you will never open. Here is what each group is, which ones you can take
          out, and which ones are the application.
        </p>
      </div>

      <div className="prose-grit">
        <p>
          Every number here was measured on a{' '}
          <code>grit new measured --triple --next</code> at v3.569.0, not estimated. The
          &quot;nothing imports it&quot; claims were checked by grepping every import in{' '}
          <code>apps/api</code>, and then by deleting the six directories and running{' '}
          <code>go build ./...</code> and <code>go test ./...</code>, which passed.
        </p>

        <h2 id="what-the-number-is">What the number actually is</h2>
      </div>

      <div className="grid gap-3 sm:grid-cols-2 my-6">
        {[
          { n: COUNTS.go, label: 'Go files in apps/api', note: `including ${COUNTS.goTests} test files that ship with your project and run themselves` },
          { n: COUNTS.adminTs, label: 'TypeScript in apps/admin', note: 'the admin panel: components, screens, hooks' },
          { n: COUNTS.webTs, label: 'TypeScript in apps/web', note: 'the public site, which is close to empty on purpose' },
          { n: COUNTS.total - COUNTS.go - COUNTS.adminTs - COUNTS.webTs, label: 'everything else', note: 'config, Docker, CI, the shared package, assets' },
        ].map((card) => (
          <div key={card.label} className="rounded-xl border border-border/40 bg-card/30 p-5">
            <p className="text-2xl font-semibold">{card.n}</p>
            <p className="text-sm font-medium mt-0.5">{card.label}</p>
            <p className="text-[13px] text-muted-foreground mt-1 leading-relaxed">{card.note}</p>
          </div>
        ))}
      </div>

      <div className="prose-grit">
        <p>
          The useful split is not by language. It is: code you will edit, code you will
          read once, and code you will never open. The third group is most of it, and it
          is the group this page is about.
        </p>

        <h2 id="sample-resources">The sample resources, which have a command</h2>
        <p>
          A new project ships a <strong>Blog</strong> and a <strong>Post</strong> so
          there is something in the admin on the first run. They are ordinary generated
          resources, and the generator can take them back out:
        </p>
      </div>

      <CodeBlock
        language="bash"
        code={`grit remove resource Blog
grit remove resource Post    # --force to skip the confirmation`}
      />

      <div className="prose-grit">
        <p>
          This is the one removal you should not do by hand. A resource with two plain
          fields is fourteen files plus thirty edits to files it does not own: the
          migration registry, the route table, the authz catalogue, the admin registry,
          the shared type exports, the route-coverage table.{' '}
          <code>grit remove</code> reverses every one of them and drops the manifest
          entries, and it prints what it did:
        </p>
      </div>

      <CodeBlock
        output="grit remove resource Invoice --force"
        language="text"
        code={`  ✗ Removed permissions from the authz catalog
  ✓ apps/api/internal/access/registry.go now covers 209 routes
  ✗ Removed 16 manifest entries

  ✅ Resource Invoice removed successfully!`}
      />

      <div className="prose-grit">
        <p>
          Deleting the files instead leaves a project that does not compile, and if it
          does compile, leaves permissions in the catalogue for a resource that is gone.
        </p>

        <h2 id="waiting">The libraries waiting for a field type</h2>
        <p>
          Six Go packages, {waitingLines.toLocaleString()} lines between them, that{' '}
          <strong>nothing imports</strong>. They are not dead code and they are not
          loaded at runtime: they are there so that the day you type{' '}
          <code>price:money</code>, the generator has somewhere to point.
        </p>
        <p>
          Measured on a project that has already generated a resource, and that
          condition is the whole subtlety. Before you generate anything,{' '}
          <code>pdf</code>, <code>rowrules</code>, <code>export</code>,{' '}
          <code>factory</code> and <code>testkit</code> look unimported too, and every
          one of them is used by the first resource you create. Do not read a bare
          project&apos;s import graph as a deletion list.
        </p>
      </div>

      <div className="my-6 overflow-hidden rounded-xl border border-border/40">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40 bg-card/50">
              <th className="px-4 py-2.5 text-left font-medium">Package</th>
              <th className="px-4 py-2.5 text-right font-medium">Lines</th>
              <th className="px-4 py-2.5 text-left font-medium">What needs it</th>
            </tr>
          </thead>
          <tbody>
            {WAITING.map((row) => (
              <tr key={row.name} className="border-b border-border/20 last:border-0">
                <td className="px-4 py-2.5 font-mono text-[13px] text-accent align-top">
                  internal/{row.name}
                </td>
                <td className="px-4 py-2.5 text-right text-muted-foreground align-top">
                  {row.lines}
                </td>
                <td className="px-4 py-2.5 text-[13px] text-muted-foreground leading-relaxed">
                  {row.what}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="prose-grit">
        <p>
          You can delete all six today and the project builds and its tests pass. What
          happens next is the part worth knowing, because it is not a subtle failure:
        </p>
      </div>

      <CodeBlock
        output="grit generate resource Invoice --fields total:money,due:date"
        language="text"
        code={`$ grit generate resource Invoice --fields "total:money,due:date"
$ cd apps/api && go build ./...

internal/models/invoice.go:9:2: package myapp/apps/api/internal/jsontime is not in std
internal/models/invoice.go:10:2: package myapp/apps/api/internal/money is not in std`}
      />

      <div className="prose-grit">
        <p>
          So the rule is simple: delete one of these when you are confident you will
          never use the field type it serves, and put it back from a fresh{' '}
          <code>grit new</code> if you were wrong. Nothing else in the project refers to
          them, so there is no second edit to make in either direction.
        </p>
        <p>
          <code>grit upgrade</code> does not bring back a directory you deleted. It
          updates the files it wrote and still finds; a file that is not there is not a
          file it is maintaining.
        </p>

        <h2 id="whole-apps">Whole apps</h2>
        <p>
          <code>apps/web</code>, <code>apps/admin</code>, <code>apps/expo</code> and{' '}
          <code>apps/desktop</code> are each a directory and a line in the pnpm
          workspace. Deleting one is deleting both, and nothing in the Go API knows they
          exist.
        </p>
        <p>
          Better: do not scaffold it. The{' '}
          <Link href="/docs/concepts/architecture-modes">architecture you pick</Link> at{' '}
          <code>grit new</code> decides which apps exist, and{' '}
          <code>--double</code> is the triple shape without the admin,{' '}
          <code>--api</code> is the Go server alone. A project you never generated is
          cheaper than one you deleted.
        </p>

        <h2 id="load-bearing">What is the application</h2>
        <p>
          These are imported from everywhere, and the import graph is the evidence rather
          than the names:
        </p>
        <p>
          {LOAD_BEARING.map((name) => (
            <code key={name} className="mr-2">internal/{name}</code>
          ))}
        </p>
        <p>
          The same goes for anything mounted in <code>routes.go</code>. If you want a
          feature off rather than gone, most of them are a line in{' '}
          <code>.env</code> or a module switch in the admin, which is reversible and does
          not need a diff.
        </p>

        <h2 id="lean">Is there a --lean preset?</h2>
        <p>
          No, and after measuring this page I do not think there should be. Of the{' '}
          {COUNTS.total} files, the ones a preset could leave out are the six directories
          above, {waitingLines.toLocaleString()} lines that nothing imports and nothing
          loads. The cost of them is disk, and the cost of leaving them out is a
          generator that cannot emit half its field types until you notice why.
        </p>
        <p>
          What actually makes a project smaller is choosing the architecture: an{' '}
          <code>--api</code> project is the Go server with no frontend at all. What makes
          it feel smaller is knowing which files you will never open, which is what this
          page is for.
        </p>
      </div>

      <div className="mt-12 flex flex-wrap gap-3">
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
          <Link href="/docs/concepts/generated-files" className="gap-1.5">
            <ArrowLeft className="h-3.5 w-3.5" />
            What each generated file is
          </Link>
        </Button>
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
          <Link href="/docs/concepts/architecture-modes" className="gap-1.5">
            Pick a smaller architecture
            <ArrowRight className="h-3.5 w-3.5" />
          </Link>
        </Button>
      </div>
    </DocsShell>
  )
}
