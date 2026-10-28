import Link from 'next/link'
import { ArrowRight, Blocks, Boxes, Layers, Package, Puzzle } from 'lucide-react'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Extending Grit: resources, plugins, UI and blueprints — Grit',
  description:
    'Four ways to add something to a Grit project, what each one is for, and how to tell which one a given job needs. Resources generate a noun. Plugins package a feature. Grit UI gives a component. A blueprint is a whole application.',
  alternates: { canonical: 'https://gritframework.dev/docs/extending' },
}

const MECHANISMS = [
  {
    name: 'Resource',
    command: 'grit generate resource Invoice',
    icon: Layers,
    unit: 'A noun in your domain',
    gives: 'Model, migration, service, handler, routes, Zod schema, TypeScript types, React hooks, an admin screen',
    yours: 'Everything it writes',
    remove: 'grit remove resource Invoice',
    href: '/docs/cli/generate',
    when: 'The thing you are adding is a kind of record: an invoice, a patient, a shipment.',
  },
  {
    name: 'Plugin',
    command: 'grit plugin add kanban',
    icon: Puzzle,
    unit: 'A feature, across every layer it touches',
    gives: 'Whatever the plugin writes: models, services, routes, pages, components, dependencies',
    yours: 'Everything it writes',
    remove: 'grit plugin remove kanban',
    href: '/docs/plugins/overview',
    when: 'The thing you are adding is a capability rather than a record, and more than one project would want it.',
  },
  {
    name: 'Grit UI block',
    command: 'grit ui add pricing-table',
    icon: Blocks,
    unit: 'A component',
    gives: 'One React component file, in your frontend',
    yours: 'The file',
    remove: 'Delete the file',
    href: '/docs/frontend/ui-components',
    when: 'You need something to look right, and the data behind it already exists.',
  },
  {
    name: 'Blueprint',
    command: 'grit blueprint list',
    icon: Boxes,
    unit: 'A whole application',
    gives: 'A repository you clone and work from',
    yours: 'The repository',
    remove: 'Not applicable: it is the project',
    href: '/docs/blueprints',
    when: 'You are starting something, and something close to it already exists.',
  },
]

export default function ExtendingPage() {
  return (
    <DocsShell>
          <div className="mb-3 flex items-center gap-2">
            <Package className="h-5 w-5 text-accent" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Extending
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
            Four ways to add something
          </h1>
          <p className="mb-10 text-lg leading-relaxed text-muted-foreground">
            Grit has four mechanisms for adding to a project, and they are not alternatives to each
            other. Each one has a unit: a noun, a feature, a component, an application. Picking the
            right one is mostly a matter of naming the unit.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">The short answer</h2>

          <div className="mb-8 overflow-hidden rounded-xl border border-border">
            <table className="w-full text-sm">
              <thead className="bg-card/60">
                <tr className="border-b border-border text-left">
                  <th className="px-4 py-3 font-medium">Mechanism</th>
                  <th className="px-4 py-3 font-medium">The unit</th>
                  <th className="px-4 py-3 font-medium">Reach for it when</th>
                </tr>
              </thead>
              <tbody>
                {MECHANISMS.map((m) => (
                  <tr key={m.name} className="border-b border-border/60 last:border-0 align-top">
                    <td className="px-4 py-3">
                      <Link href={m.href} className="font-medium text-accent hover:underline">
                        {m.name}
                      </Link>
                    </td>
                    <td className="px-4 py-3 text-muted-foreground">{m.unit}</td>
                    <td className="px-4 py-3 text-muted-foreground">{m.when}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            What they have in common
          </h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            All four write ordinary code into your repository. None of them installs a runtime
            library that reaches into your app and registers itself. After any of them runs, what
            you have is files you can read, edit, step through in a debugger and delete.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            That is the property that makes the choice between them low-stakes. Picking the wrong
            one costs you a refactor, not a migration off a dependency.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            Resource: a noun in your domain
          </h2>
          <CodeBlock
            language="bash"
            code={`grit generate resource Invoice --fields "number:string,total:money,due:date,client:belongs_to:Client"`}
          />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            This is the one you use most, and the one the framework is built around. One definition
            produces the Go model and migration, the service and handler, the routes, the shared Zod
            schema, the TypeScript types, the React Query hooks and the admin screen. They agree
            with each other because they came from the same description.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            The test is whether the thing is a kind of record. If you would put it in a table and
            list it, it is a resource.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            Plugin: a feature, packaged
          </h2>
          <CodeBlock language="bash" code={`grit plugin add kanban`} />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            Some things are not nouns. Multi-tenancy is not a record. Neither is a command palette,
            or signing in as another user, or a board you drag cards around on. They cut across
            models, routes, pages and dependencies at once, and the pieces only make sense together.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            A plugin is that bundle. It writes files, patches existing ones at named markers,
            declares the dependencies it needs, and records every one of those actions in a
            lockfile. Removal replays the lockfile backwards, which is why{' '}
            <code>grit plugin remove</code> leaves a project that still compiles rather than a
            project with half a feature in it.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            The <Link href="/docs/plugins/kanban" className="text-accent hover:underline">kanban plugin</Link>{' '}
            is the worked example: it is the smallest feature that needs every part of the system at
            once, so the{' '}
            <Link href="/docs/plugins/authoring" className="text-accent hover:underline">
              authoring guide
            </Link>{' '}
            is written around it.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            Grit UI: a component, nothing behind it
          </h2>
          <CodeBlock language="bash" code={`grit ui list\ngrit ui add pricing-table`} />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            Grit UI is a component library in the shadcn format: the command downloads the source
            and writes it into your frontend, into the right app for your architecture. It is a file
            from that point on.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            The thing to be clear about is what it does not do. A block is presentation. It has no
            model, no endpoint and no migration. If you add a board component this way you will have
            a board that displays nothing and saves nothing, which is a worse place to be than
            having no board, because it looks finished.
          </p>
          <p className="mb-6 leading-relaxed text-muted-foreground">
            So: a UI block when the data already exists and you need it to look right. A plugin when
            the data does not exist yet.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            Blueprint: a whole application
          </h2>
          <CodeBlock language="bash" code={`grit blueprint list\ngrit blueprint info whatsapp`} />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            A blueprint is not installed. It is a complete application in its own repository, built
            with Grit, that you clone and work from. A plugin adds a feature to your project; a
            blueprint is a project.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            They exist for two reasons. One is a starting point: if you are building a storefront,
            starting from a storefront beats starting from an empty project. The other is that each
            blueprint is a load test for the framework. Every gap one finds gets fixed in Grit
            first, released, and then the blueprint upgrades into it, so its code stays Grit&apos;s
            own rather than drifting into a fork.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            How to tell which one you need
          </h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Two questions settle almost every case.
          </p>
          <div className="mb-6 rounded-xl border border-border bg-card/40 p-5">
            <p className="mb-3 leading-relaxed">
              <strong>Is it a kind of record?</strong> If yes, it is a resource. Stop here.
            </p>
            <p className="mb-3 leading-relaxed">
              <strong>Would a second project want it?</strong> If yes, it is a plugin. If no, it is
              just code you write in your project, and that is fine. Not everything needs a
              mechanism.
            </p>
            <p className="leading-relaxed text-muted-foreground">
              A UI block is the answer to a narrower question: the feature exists and the screen
              looks wrong.
            </p>
          </div>

          <h3 className="mb-3 mt-8 text-xl font-semibold tracking-tight">Worked examples</h3>
          <div className="mb-8 overflow-hidden rounded-xl border border-border">
            <table className="w-full text-sm">
              <thead className="bg-card/60">
                <tr className="border-b border-border text-left">
                  <th className="px-4 py-3 font-medium">You want</th>
                  <th className="px-4 py-3 font-medium">Use</th>
                  <th className="px-4 py-3 font-medium">Because</th>
                </tr>
              </thead>
              <tbody className="align-top">
                <tr className="border-b border-border/60">
                  <td className="px-4 py-3">Invoices</td>
                  <td className="px-4 py-3 font-medium">Resource</td>
                  <td className="px-4 py-3 text-muted-foreground">An invoice is a record.</td>
                </tr>
                <tr className="border-b border-border/60">
                  <td className="px-4 py-3">A Kanban board</td>
                  <td className="px-4 py-3 font-medium">Plugin</td>
                  <td className="px-4 py-3 text-muted-foreground">
                    Three models plus a reordering endpoint plus a drag-and-drop screen, useless
                    apart.
                  </td>
                </tr>
                <tr className="border-b border-border/60">
                  <td className="px-4 py-3">A chat feature</td>
                  <td className="px-4 py-3 font-medium">Plugin, and read the WhatsApp blueprint</td>
                  <td className="px-4 py-3 text-muted-foreground">
                    Same shape as a board, with realtime delivery and read state on top.
                  </td>
                </tr>
                <tr className="border-b border-border/60">
                  <td className="px-4 py-3">A nicer pricing page</td>
                  <td className="px-4 py-3 font-medium">Grit UI block</td>
                  <td className="px-4 py-3 text-muted-foreground">
                    The plans already exist; this is presentation.
                  </td>
                </tr>
                <tr className="border-b border-border/60">
                  <td className="px-4 py-3">Per-organisation data isolation</td>
                  <td className="px-4 py-3 font-medium">Plugin</td>
                  <td className="px-4 py-3 text-muted-foreground">
                    It changes how every query is scoped, which is not a record.
                  </td>
                </tr>
                <tr>
                  <td className="px-4 py-3">A shop to start from</td>
                  <td className="px-4 py-3 font-medium">Blueprint</td>
                  <td className="px-4 py-3 text-muted-foreground">
                    You are starting a project, not adding to one.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            When none of them fits
          </h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Sometimes the answer is that the thing belongs in the framework rather than on top of
            it. The test Grit applies to itself: is this machinery a team would otherwise rebuild
            for every client target, and can one definition drive it everywhere? Ordered lists that
            survive two people reordering at once is machinery. A Kanban board is not; it is one
            screen that happens to need that machinery.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            If you hit something in that first category, that is worth{' '}
            <Link
              href="https://github.com/MUKE-coder/grit/issues/new/choose"
              className="text-accent hover:underline"
            >
              saying so
            </Link>
            . Most of what is in Grit got there because building something real ran into its
            absence.
          </p>

          <div className="mt-12 grid gap-3 sm:grid-cols-2">
            {MECHANISMS.map((m) => {
              const Icon = m.icon
              return (
                <Link
                  key={m.name}
                  href={m.href}
                  className="group flex items-center gap-3 rounded-xl border border-border p-4 transition-colors hover:border-accent"
                >
                  <Icon className="h-5 w-5 shrink-0 text-accent" />
                  <span className="flex-1 font-medium">{m.name}</span>
                  <ArrowRight className="h-4 w-4 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
                </Link>
              )
            })}
          </div>
    </DocsShell>
  )
}
