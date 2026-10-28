import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CommunityCTA } from '@/components/community-cta'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/architecture-style')

/**
 * "Which architecture is Grit?" and "can I refactor to another one?"
 *
 * Deliberately separate from /docs/concepts/architecture, which answers a
 * different question with the same word. In this project "architecture" names
 * the single/double/triple/api/mobile choice, so the style question needed a
 * slug of its own rather than a section inside that page.
 */

/** The five styles a Grit project is already running, and where each one sits. */
const STYLES_IN_USE = [
  {
    style: 'Modular monolith',
    where: 'The whole API: one binary, one deployment, around fifty internal packages.',
    why: 'One process removes every distributed-systems failure mode from a product that does not yet have the traffic to need them.',
  },
  {
    style: 'Layered (n-tier)',
    where: 'routes to handlers to services to models.',
    why: 'The service owns the rules and every database call, so a job, a CLI command and a test get the same rules as the HTTP route.',
  },
  {
    style: 'Event-driven, in process',
    where: 'internal/events. A write emits once; audit, webhooks and realtime subscribe.',
    why: 'The handler does not have to remember to call three systems, so adding a fourth does not mean editing every handler.',
  },
  {
    style: 'Transactional outbox',
    where: 'internal/outbox, behind the bus Durable delivery mode.',
    why: 'The event is queued in the same transaction as the write, so a crash between the commit and the HTTP call loses nothing.',
  },
  {
    style: 'Orchestration (saga)',
    where: 'internal/saga, for work that spans modules and can half fail.',
    why: 'The sequence and its compensation are readable in one place, which is the thing choreography gives up.',
  },
  {
    style: 'Microkernel / plug-in',
    where: 'The framework itself: grit plugin add.',
    why: 'Capabilities that only some projects want are not carried by the ones that do not.',
  },
]

/** What Grit is not, and the honest reason in each case. */
const NOT_THIS = [
  {
    style: 'Microservices',
    why: 'A split is worth paying for when teams need to release independently. Until then it buys network failures, eventual consistency and a deployment pipeline, in exchange for nothing.',
  },
  {
    style: 'CQRS and event sourcing',
    why: 'Both earn their complexity in a hot spot, not across a codebase. The read path is already separate enough to point at another store when one query needs it.',
  },
  {
    style: 'Hexagonal / ports and adapters',
    why: 'The generated service talks to GORM directly rather than through a port. That is a deliberate trade: less indirection to read, at the cost of swapping the ORM being a real job.',
  },
  {
    style: 'Serverless',
    why: 'Background jobs, a scheduler, realtime connections and a warm connection pool all assume a process that stays up.',
  },
]

/** Refactors people actually ask about, cheapest first. */
const REFACTORS = [
  {
    move: 'Add CQRS to one endpoint',
    cost: 'Hours',
    verdict: 'Yes',
    tone: 'good' as const,
    detail:
      'The read path is already a separate whitelist of searchable, sortable and filterable columns. Point one list query at a denormalised table or a read replica and leave every other endpoint alone. This is the refactor that is almost always the right one, because it is local and reversible.',
  },
  {
    move: 'Extract one module into its own service',
    cost: 'Weeks',
    verdict: 'Yes, with care',
    tone: 'good' as const,
    detail:
      'The service layer is the seam: it is already the only thing that touches the database for its resource, and it already takes a context rather than a request. Replace the in-process call with an HTTP client and give the new service its own schema. What makes it expensive is not the code, it is the foreign keys you have to cut.',
  },
  {
    move: 'Go fully hexagonal inside the monolith',
    cost: 'Weeks, and it keeps costing',
    verdict: 'Possible, rarely worth it',
    tone: 'warn' as const,
    detail:
      'Putting a port in front of every service is mechanical. The ongoing cost is that the generator keeps writing the un-ported version, so every generate is followed by a hand edit and every upgrade reports a conflict. If you want this shape, the honest answer is that Grit is fighting you.',
  },
  {
    move: 'Split into microservices across the board',
    cost: 'Months',
    verdict: 'Only with a reason',
    tone: 'warn' as const,
    detail:
      'The reason should be organisational: separate teams with separate release cadences, or one module whose scaling profile is nothing like the rest. "It will scale better" is not a reason on its own, and a modular monolith on one large machine outruns most products for years.',
  },
  {
    move: 'Move to event sourcing',
    cost: 'A rewrite of the write path',
    verdict: 'No, not as a refactor',
    tone: 'bad' as const,
    detail:
      'Current state stops being a row and becomes a fold over events, which changes every read, every migration and every admin screen. The append-only audit log is not a head start on this: it records what happened for a reviewer, it is not the source the state is rebuilt from.',
  },
  {
    move: 'Collapse existing microservices into Grit',
    cost: 'Weeks',
    verdict: 'Yes, and often the right direction',
    tone: 'good' as const,
    detail:
      'Easier than the reverse and more common than people admit. Each service becomes a package with its own services and models, cross-service HTTP calls become function calls, and the events between them become bus subscriptions.',
  },
]

/**
 * Every style on the ByteByteGo "Software Architecture Styles" wheel, which is
 * the chart people have in mind when they ask this question, with Grit's
 * position on each. Written so that whichever segment someone points at, the
 * page has an answer rather than a gap.
 *
 * `use` is deliberately three-valued. "Partly" is the honest answer more often
 * than either of the other two, and collapsing it into yes or no is how a
 * reference page starts overclaiming.
 */
const WHEEL: {
  family: string
  styles: { name: string; use: 'yes' | 'partly' | 'no'; note: string }[]
}[] = [
  {
    family: 'Layered',
    styles: [
      {
        name: 'Layered (n-tier)',
        use: 'yes',
        note: 'The API itself: routes, handlers, services, models.',
      },
      {
        name: 'Clean / Onion',
        use: 'no',
        note: 'The service imports GORM rather than depending inward on a port.',
      },
    ],
  },
  {
    family: 'Component-Based',
    styles: [
      {
        name: 'Object-Oriented',
        use: 'partly',
        note: 'A paradigm rather than a deployment shape. Go structs with methods, no inheritance.',
      },
      {
        name: 'Microkernel',
        use: 'yes',
        note: 'The framework: a core project plus capabilities added with grit plugin add.',
      },
      {
        name: 'Plug-in',
        use: 'yes',
        note: 'The same thing seen from the plugin side. Eleven modules also switch off with MODULE_<NAME>=false.',
      },
    ],
  },
  {
    family: 'Service-Oriented',
    styles: [
      { name: 'SOA', use: 'no', note: 'No enterprise bus, no coarse-grained shared contracts.' },
      {
        name: 'Broker',
        use: 'partly',
        note: 'Redis carries background jobs through asynq. That is a queue, not a request bus: nothing routes calls between services.',
      },
      {
        name: 'Microservices',
        use: 'no',
        note: 'One binary by design. The seams for a split are listed below.',
      },
      {
        name: 'Serverless (FaaS)',
        use: 'no',
        note: 'Jobs, a scheduler, realtime connections and a warm pool all assume a process that stays up.',
      },
    ],
  },
  {
    family: 'Distributed System',
    styles: [
      { name: 'Space-Based', use: 'no', note: 'The database is the system of record, not an in-memory grid.' },
      { name: 'Peer-to-Peer', use: 'no', note: 'Clients talk to one API.' },
    ],
  },
  {
    family: 'Domain-Driven',
    styles: [
      {
        name: 'DDD',
        use: 'partly',
        note: 'It borrows the vocabulary, domain events and a service per resource, not the method: no bounded contexts, aggregates or ubiquitous-language modelling.',
      },
      {
        name: 'Hexagonal (ports & adapters)',
        use: 'no',
        note: 'Deliberate: less indirection to read, at the cost of swapping the ORM being real work.',
      },
    ],
  },
  {
    family: 'Event-Driven',
    styles: [
      {
        name: 'Event-Driven',
        use: 'yes',
        note: 'internal/events. A write emits once, with Sync, Async and Durable delivery.',
      },
      {
        name: 'Publish-Subscribe',
        use: 'yes',
        note: 'That bus is pub-sub: subscribers match a name, a resource prefix or everything, and the emitter never knows who listens.',
      },
    ],
  },
  {
    family: 'Data-Centric',
    styles: [
      {
        name: 'CQRS',
        use: 'no',
        note: 'One model for reads and writes, but the read path is already a separate whitelist, which is what makes adding it to one endpoint cheap.',
      },
      {
        name: 'Event Sourcing',
        use: 'no',
        note: 'The audit chain is append-only evidence for a reviewer. State is rows, not a fold over events.',
      },
      { name: 'Kappa', use: 'no', note: 'No streaming pipeline.' },
      { name: 'Lambda (batch + speed)', use: 'no', note: 'Not a big-data system.' },
    ],
  },
  {
    family: 'Concurrency',
    styles: [
      {
        name: 'Orchestration',
        use: 'yes',
        note: 'internal/saga, where a sequence spans modules and can half fail, with its compensation in one place.',
      },
      {
        name: 'Choreography',
        use: 'yes',
        note: 'The other half of the same story: audit, webhooks and realtime react to an event rather than being called.',
      },
      {
        name: 'Primary-Secondary',
        use: 'partly',
        note: 'Reads go to DATABASE_REPLICA_URLS when you set them, and the outbox relay claims messages with a lease so another replica takes over if one dies.',
      },
      {
        name: 'Pipeline / Pipe-Filter',
        use: 'partly',
        note: 'Two places, both internal: the seeder (build, bound, batch, write) and the image pipeline (decode, orient, resize, encode).',
      },
    ],
  },
  {
    family: 'Interpreter',
    styles: [
      {
        name: 'Interpreter',
        use: 'no',
        note: 'The field definition is a small language, but the generator reads it once at build time. Nothing interprets a DSL at runtime.',
      },
    ],
  },
  {
    family: 'Separation of Concern',
    styles: [
      {
        name: 'MVC',
        use: 'no',
        note: 'A handler returns JSON; it does not pick a view. The frontend is React.',
      },
      { name: 'MVP', use: 'no', note: 'A desktop and Android pattern.' },
      {
        name: 'MVVM',
        use: 'no',
        note: 'React Query holds server state declaratively, which rhymes with it, but calling it MVVM would be a stretch.',
      },
    ],
  },
]

const USE_LABEL = {
  yes: { text: 'Yes', className: 'text-emerald-400' },
  partly: { text: 'Partly', className: 'text-amber-400' },
  no: { text: 'No', className: 'text-muted-foreground' },
}

/**
 * The two ways to lay the same four resources out. Kept as literal trees
 * because the whole argument is about where files sit, and a prose summary of
 * a directory tree persuades nobody.
 */
const LAYOUTS = [
  {
    title: 'By layer',
    subtitle: 'What Grit generates',
    tree: `internal/
├── models/
│   ├── invoice.go
│   ├── customer.go
│   └── ... 31 more, one package
├── services/
│   ├── invoice.go
│   └── customer.go
└── handlers/
    ├── invoice.go
    └── customer.go`,
  },
  {
    title: 'By feature',
    subtitle: 'A folder per noun',
    tree: `internal/
├── invoice/
│   ├── model.go
│   ├── invoice.go
│   └── stores/invoicesql/
├── customer/
│   ├── model.go
│   └── stores/customersql/
└── web/
    └── handlers/`,
  },
]

const TONE = {
  good: 'text-emerald-400',
  warn: 'text-amber-400',
  bad: 'text-rose-400',
}

export default function ArchitectureStylePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-accent/80 mb-3 block">Core Concepts</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Architecture style</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Which of the named architecture styles Grit actually is, which ones it
                deliberately is not, and what it costs to move a Grit project from one to
                another. Written because both questions get asked in the same breath.
              </p>
            </div>

            <div className="rounded-lg border border-accent/20 bg-accent/5 p-4 mb-10">
              <p className="text-sm text-muted-foreground leading-relaxed">
                <strong className="text-foreground">A word with two jobs.</strong> In Grit,{' '}
                <em>architecture</em> normally means the{' '}
                <Link
                  href="/docs/concepts/architecture-modes"
                  className="text-accent hover:underline"
                >
                  project shape
                </Link>{' '}
                you pick with <code>--architecture</code>: single, double, triple, api or
                mobile. This page is about something else: the architecture <em>style</em> the
                code inside is written in.
              </p>
            </div>

            <div className="prose-grit">
              <h2>The short answer</h2>
              <p>
                <strong>A modular monolith with a layered API and an in-process event bus.</strong>{' '}
                One binary, one database, around fifty internal packages, and a domain event bus
                that lets the audit log, outbound webhooks and realtime clients subscribe to a
                write rather than be called by it.
              </p>
              <p>
                That is the honest headline, and it is also slightly misleading, because no real
                system is one style. Styles compose. A Grit project is running six of them at once,
                each in the place where it earns its keep.
              </p>
            </div>

            <div className="mt-8 overflow-hidden rounded-xl border border-border/40">
              <table className="w-full text-sm">
                <thead className="bg-card/60">
                  <tr className="border-b border-border/40 text-left">
                    <th className="px-4 py-3 font-semibold text-foreground">Style</th>
                    <th className="px-4 py-3 font-semibold text-foreground">Where it sits</th>
                  </tr>
                </thead>
                <tbody>
                  {STYLES_IN_USE.map((s) => (
                    <tr key={s.style} className="border-b border-border/30 last:border-0">
                      <td className="px-4 py-4 align-top font-medium text-foreground whitespace-nowrap">
                        {s.style}
                      </td>
                      <td className="px-4 py-4 align-top text-muted-foreground leading-relaxed">
                        {s.where}
                        <span className="mt-1.5 block text-foreground/70">{s.why}</span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="prose-grit mt-10">
              <h2>The one rule the layering exists for</h2>
              <p>
                Handlers are thin on purpose. A handler reads the request, calls one service
                method and writes the answer. Every database read and write for a resource lives
                in that resource&apos;s service, and every service method takes a{' '}
                <code>context.Context</code> rather than a <code>*gin.Context</code>.
              </p>
              <p>
                The payoff is not tidiness. It is that a rule enforced in the service is enforced
                in the background job, the CLI command, the import, the export and the test, not
                only on the HTTP route. A rule written in a handler is a rule that holds until the
                first time something other than a browser does the same thing.
              </p>

              <h2>What Grit deliberately is not</h2>
            </div>

            <div className="mt-6 grid gap-4 sm:grid-cols-2">
              {NOT_THIS.map((n) => (
                <div key={n.style} className="rounded-lg border border-border/40 bg-card/40 p-5">
                  <div className="mb-2 text-sm font-semibold text-foreground">{n.style}</div>
                  <p className="text-sm leading-relaxed text-muted-foreground">{n.why}</p>
                </div>
              ))}
            </div>

            <div className="prose-grit mt-12">
              <h2>The whole wheel, segment by segment</h2>
              <p>
                The question usually arrives with a chart attached: the one with twenty-seven
                styles arranged in a circle. Here is Grit&apos;s position on every segment of it,
                so whichever one you are pointing at has an answer.
              </p>
              <p>
                Three of the answers are <strong>partly</strong>, and that is the honest one more
                often than either of the others. Most styles are not a thing a system is or is
                not: they are a thing a system does in one place and not another.
              </p>
            </div>

            <div className="mt-6 overflow-hidden rounded-xl border border-border/40">
              {WHEEL.map((group) => (
                <div key={group.family}>
                  <div className="border-b border-border/40 bg-card/60 px-4 py-2">
                    <span className="font-mono text-[11px] uppercase tracking-[0.18em] text-muted-foreground">
                      {group.family}
                    </span>
                  </div>
                  {group.styles.map((s) => (
                    <div
                      key={s.name}
                      className="flex flex-col gap-1 border-b border-border/25 px-4 py-3 last:border-b-0 sm:flex-row sm:gap-4"
                    >
                      <div className="flex shrink-0 items-baseline gap-2.5 sm:w-64">
                        <span className={`text-xs font-medium ${USE_LABEL[s.use].className}`}>
                          {USE_LABEL[s.use].text}
                        </span>
                        <span className="text-sm font-medium text-foreground">{s.name}</span>
                      </div>
                      <p className="text-sm leading-relaxed text-muted-foreground">{s.note}</p>
                    </div>
                  ))}
                </div>
              ))}
            </div>

            <div className="prose-grit mt-12">
              <h2>Why Grit is organised by layer and not by feature</h2>
              <p>
                This is the question that follows the one above, and it is a better question.
                Package-by-feature, a folder per noun with that noun&apos;s model, rules and
                storage inside it, is the layout most Go teams land on, and it is taught as the
                end state of the tour for good reasons. Grit does not do it. Here is the trade,
                stated honestly, because it is a real one.
              </p>
            </div>

            <div className="mt-6 grid gap-4 sm:grid-cols-2">
              {LAYOUTS.map((l) => (
                <div
                  key={l.title}
                  className="overflow-hidden rounded-xl border border-border/40 bg-card/40"
                >
                  <div className="border-b border-border/40 px-4 py-2.5">
                    <div className="text-sm font-semibold text-foreground">{l.title}</div>
                    <div className="font-mono text-[11px] uppercase tracking-[0.15em] text-muted-foreground">
                      {l.subtitle}
                    </div>
                  </div>
                  <pre className="overflow-x-auto px-4 py-3 font-mono text-[12px] leading-relaxed text-muted-foreground">
                    {l.tree}
                  </pre>
                </div>
              ))}
            </div>

            <div className="prose-grit mt-8">
              <h3>What feature folders buy, and Grit gives up</h3>
              <p>
                A feature is a folder. Everything about an invoice is in one place and nothing
                else is, you can delete a feature by deleting a directory, and the root stops
                growing a package every time a new technology arrives. Grit pays for not having
                this: at forty resources, <code>services/</code> is forty files and the directory
                no longer tells you which three of them form a subsystem. That cost is real and
                it arrives exactly when the codebase is big enough to need the help.
              </p>

              <h3>What Grit buys back: the entity graph</h3>
              <p>
                The examples that make feature folders look obvious have two or three domains with
                a clean one-way reference. A business application has thirty entities with a dense
                graph: an invoice points at a customer, which points at an organisation, which
                points back at users who raised invoices.
              </p>
              <p>
                <strong>Go forbids import cycles.</strong> So the moment <code>customer</code>{' '}
                needs to refer back to <code>invoice</code>, feature folders force you into an
                event announcement or an interface dance to express what is really just a foreign
                key. That is the correct answer when two domains must stay independent, and it is
                a lot of machinery for <code>invoice.customer_id</code>.
              </p>
              <p>
                Grit&apos;s thirty-odd models are one flat <code>package models</code>. Any model
                can reference any other with no ceremony, which is the thing that lets{' '}
                <code>client:belongs_to</code> be a single word in a generate command rather than
                a design conversation. The layout is not laziness; it is what buys the generator
                its simplicity.
              </p>
              <p>
                The rule still bites in the other direction, and it is worth knowing before you go
                looking for it: <code>models</code> cannot import <code>services</code>, so a hook
                on a model calls the lower-level package directly rather than the service that
                wraps it.
              </p>

              <h3>The honest summary</h3>
              <p>
                Feature folders are the better default for an application a team writes by hand,
                with a handful of genuinely independent domains. Layer folders are the better
                default for a <em>generated</em> application with a dense entity graph, because
                one command has to know where every file goes and the entities all point at each
                other. Those are different problems, and the layouts are not competing for the
                same job.
              </p>
              <p>
                If your application really is a few independent domains, package-by-feature is a
                good choice and Grit is the wrong tool, which is a better thing to find out now
                than in month three.
              </p>

              <h3>One thing worth taking from it either way</h3>
              <p>
                The strongest idea in the feature-folder literature is not the folders. It is that
                a dependency should be declared by the code that <em>uses</em> it and named for
                what it needs, rather than taken whole: a two-method interface instead of somebody
                else&apos;s four-method one. That is right in every layout, Grit included.
              </p>
              <p>
                It is worth being precise about where it pays, because the obvious application of
                it here does not. A generated handler calls ten of its service&apos;s ten exported
                methods, so a &ldquo;narrow&rdquo; interface in front of that service would be a
                verbatim copy of the service: the same thing written twice, which is the cost the
                idea exists to avoid. Grit does not generate one, on purpose.
              </p>
              <p>
                Where it does pay is the weaker form of the same rule:{' '}
                <strong>do not reach around a thing to rebuild what it already knows.</strong> The
                blog list handler used to skip its service and query the database directly,
                re-typing the search clause the service already had, so the list and the stat
                cards above it held two definitions of the same search in two different files. Add
                a searchable column to one and they disagree. The service now states that search
                once and exposes it as <code>CountQuery</code>, which the handler asks for by
                name. Shipped in v3.403.0.
              </p>
              <p>
                The lesson generalises better than the interface does. Count how many methods the
                caller actually uses before you narrow something, and spend the effort instead on
                the places where a consumer has quietly grown its own copy of a rule.
              </p>

              <h3>A note on the name</h3>
              <p>
                Package-by-feature is often called domain-driven, and the folders are named after
                domain nouns, so the label is not wrong. It is worth keeping separate from{' '}
                <strong>Domain-Driven Design</strong> as a method, though, which is bounded
                contexts, aggregates and a language agreed with domain experts. That is a
                modelling discipline, not a directory layout. You can do feature folders with no
                DDD in them, and real DDD inside any layout on this page.
              </p>
            </div>

            <div className="prose-grit mt-12">
              <h2>Can you refactor from one style to another?</h2>
              <p>
                Yes, and the answer is more interesting than yes. The thing that decides whether a
                refactor is a weekend or a quarter is almost never the code. It is the{' '}
                <strong>data</strong>.
              </p>
              <p>
                Moving a package is easy. Moving a table out from under three foreign keys, two
                joins and an admin screen that lists both sides of a relationship is the actual
                work, and no amount of clean code makes it cheaper. So when you are weighing a
                refactor, look at the schema first and the package layout second.
              </p>
              <p>Here are the moves people ask about, cheapest first.</p>
            </div>

            <div className="mt-6 flex flex-col gap-4">
              {REFACTORS.map((r) => (
                <div key={r.move} className="rounded-xl border border-border/40 bg-card/40 p-5">
                  <div className="mb-2 flex flex-wrap items-baseline gap-x-3 gap-y-1">
                    <span className="text-[15px] font-semibold text-foreground">{r.move}</span>
                    <span className={`text-xs font-medium ${TONE[r.tone]}`}>{r.verdict}</span>
                    <span className="font-mono text-xs text-muted-foreground">{r.cost}</span>
                  </div>
                  <p className="text-sm leading-relaxed text-muted-foreground">{r.detail}</p>
                </div>
              ))}
            </div>

            <div className="prose-grit mt-12">
              <h2>The seams, if you are planning to split</h2>
              <p>
                A Grit project has four places where a boundary can be cut without inventing one
                first. If you expect to split later, these are the ones to keep clean.
              </p>
              <ol>
                <li>
                  <strong>The service.</strong> One resource, one service, all of its database
                  access. A service that reaches into another resource&apos;s tables is the single
                  thing that will make a future split expensive, and it is invisible until then.
                </li>
                <li>
                  <strong>The event bus.</strong> A <code>Durable</code> subscriber already
                  receives work at least once, through the outbox, retried until it succeeds. That
                  is the same contract a queue between two services gives you, so turning a
                  subscriber into a remote call changes the transport and not the semantics.
                </li>
                <li>
                  <strong>The outbox.</strong> It is the thing that makes the first boundary
                  correct rather than hopeful: the event is written in the same transaction as the
                  data, so the two cannot disagree.
                </li>
                <li>
                  <strong>Ownership scoping.</strong> Projects generated with per-record ownership
                  or multitenancy already carry the owner through the context into every query. A
                  split that keeps that scoping is a split that keeps its access rules.
                </li>
              </ol>

              <h2>What a heavy refactor costs you with the generator</h2>
              <p>
                This is the part that is specific to a framework that writes code. The further your
                structure moves from what the generator emits, the more often{' '}
                <code>grit generate</code> writes something you then have to adjust, and the more
                files <code>grit upgrade</code> reports as conflicts.
              </p>
              <p>
                Nothing is overwritten: every generated file is hashed into a manifest when it is
                written, and a file whose hash has moved is reported rather than replaced. But a
                conflict is still a decision someone has to make, and twenty of them per upgrade is
                a tax. The practical guidance is to keep the refactor local. One module behind a
                port, or one endpoint on a read model, costs almost nothing. A different shape
                across the whole codebase means you have taken over the framework half, which is a
                legitimate choice as long as it is made on purpose.
              </p>

              <h2>If you are choosing a style for a new product</h2>
              <p>
                The modular monolith is not a compromise you are stuck with until you are big
                enough for the real thing. For most products it is the end state. It keeps one
                deployment, one transaction boundary and one place to look when something is wrong,
                and it leaves the seams where a split would go if the organisation ever needs one.
              </p>
              <p>
                Reach past it when the reason is concrete: separate teams that cannot coordinate a
                release, a module whose traffic is nothing like the rest, or a compliance boundary
                that has to be physical. Those are reasons. Anticipated scale is not one.
              </p>
            </div>

            <div className="mt-10 rounded-xl border border-border/40 bg-card/40 p-6">
              <h2 className="mb-2 text-lg font-semibold tracking-tight text-foreground">
                The same argument, in a talk
              </h2>
              <p className="mb-4 text-sm leading-relaxed text-muted-foreground">
                Why the unit of generation is the resource rather than the layer, and how the
                framework half keeps improving after the code is yours, are the middle of the
                GopherCon Africa talk.
              </p>
              <Link
                href="/talk#sec-concepts"
                className="group inline-flex items-center gap-2 text-sm font-medium text-accent hover:underline"
              >
                Watch that part of the deck
                <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
              </Link>
            </div>

            <CommunityCTA className="mt-10" />

            <div className="mt-12 flex flex-wrap gap-3 border-t border-border/30 pt-6">
              <Button
                variant="outline"
                asChild
                className="border-border/60 bg-transparent hover:bg-surface-2"
              >
                <Link href="/docs/concepts/architecture">Architecture Overview</Link>
              </Button>
              <Button asChild className="glow-purple-sm ml-auto">
                <Link href="/docs/concepts/code-generation">
                  Code Generation
                  <ArrowRight className="ml-1.5 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
