import Link from 'next/link'
import type { Metadata } from 'next'
import {
  InfoPage,
  Section,
  H3,
  Option,
  CompareTable,
  WhereGritFits,
  Faq,
  Related,
} from '@/components/informational'

export const metadata: Metadata = {
  title: 'CRUD generators for Go, compared',
  description:
    'sqlc, ent, GORM gen, goa and full-stack scaffolders compared: what each one writes, what it leaves to you, and why removal matters as much as generation.',
  alternates: { canonical: 'https://gritframework.dev/i/go-crud-generator' },
  openGraph: {
    title: 'CRUD and scaffolding generators for Go, compared',
    description: 'What each generator writes, and what it leaves you to do.',
    url: 'https://gritframework.dev/i/go-crud-generator',
    type: 'article',
  },
}

export default function Page() {
  return (
    <InfoPage
      slug="/i/go-crud-generator"
      eyebrow="Guide"
      title="CRUD and scaffolding generators for Go"
      checked="13 November 2026"
      lead={
        <>
          <p className="mb-4">
            Go has no runtime reflection culture and no decorators, so the ecosystem solved
            repetition the other way: by generating code you can read. <code>go generate</code> has
            been in the toolchain since 2014, and the best Go libraries of the last decade are
            generators.
          </p>
          <p>
            They do not all generate the same thing, though, and the word &quot;CRUD
            generator&quot; covers at least four different jobs. Knowing which one you are shopping
            for is most of the decision.
          </p>
        </>
      }
    >
      <Section id="kinds" title="Four different jobs">
        <H3>SQL to Go</H3>
        <p>
          You write queries, the generator writes typed Go functions for them. No ORM, no surprises
          about what SQL runs. sqlc is the standard here.
        </p>

        <H3>Schema to everything</H3>
        <p>
          You describe entities in code or a DSL, the generator writes the model, the query builder
          and the migrations. ent and GORM gen live here.
        </p>

        <H3>Contract to both sides</H3>
        <p>
          You describe the API, the generator writes the server stubs and the clients. goa and
          protobuf with gRPC.
        </p>

        <H3>Model to the whole feature</H3>
        <p>
          One command writes the model, the service, the HTTP handlers, the routes, the validation,
          the frontend types and the admin screen. This is scaffolding in the Rails sense, and it is
          the thinnest category in Go.
        </p>
      </Section>

      <Section id="glance" title="At a glance">
        <CompareTable
          caption="Go code generators compared"
          columns={['Tool', 'You write', 'It writes', 'Stars', 'Last activity']}
          rows={[
            ['sqlc', 'SQL queries', 'Typed Go functions', '18.5k', 'November 2026'],
            ['ent', 'A Go schema', 'Models, queries, migrations, GraphQL', '17.2k', 'October 2026'],
            ['GORM gen', 'GORM models', 'Type-safe query code', '2.6k', 'October 2026'],
            ['goa', 'A design DSL', 'Server, client, OpenAPI', '6.1k', 'November 2026'],
            ['protobuf and gRPC', 'A .proto file', 'Server, client, types', 'Google', 'Active'],
            ['Grit', 'One command', 'Model to admin screen', 'See repo', 'November 2026'],
          ]}
        />
      </Section>

      <Section id="options" title="The options">
        <Option
          n={1}
          name="sqlc"
          href="https://github.com/sqlc-dev/sqlc"
          meta="18,467 stars · active November 2026 · MIT"
          what={
            <p>
              You write a query in a <code>.sql</code> file with a name and a return shape; sqlc
              reads your schema and generates a typed Go function. No ORM between you and the
              database, no surprise N+1, and a compile error when a column disappears. It is the
              best-liked tool in this category by some distance.
            </p>
          }
          pros={[
            'The SQL that runs is the SQL you wrote',
            'A schema change breaks the build rather than production',
            'No runtime dependency at all',
          ]}
          cons={[
            'Dynamic filters are awkward: conditional WHERE clauses fight the model',
            'You write every query, including the dull ones',
            'Nothing above the database layer',
          ]}
          bestFor="Teams who like SQL and want the compiler to check it."
        />

        <Option
          n={2}
          name="ent"
          href="https://github.com/ent/ent"
          meta="17,230 stars · active October 2026 · Apache 2.0 · from Meta"
          what={
            <p>
              Describe entities and edges in Go, and ent generates a fully typed client: queries,
              mutations, eager loading, migrations, and optionally a GraphQL layer. The graph model
              makes traversals genuinely pleasant, and the generated API is large enough to feel
              like learning a framework.
            </p>
          }
          pros={[
            'Completely type-safe queries, including traversals',
            'Migrations and a GraphQL layer come from the same schema',
            'Strong on modelling relationships',
          ]}
          cons={[
            'A large generated API to learn',
            'Its conventions reach into how you structure the application',
            'Harder to drop into an existing database than sqlc',
          ]}
          bestFor="Graph-shaped domains where relationships are the complexity."
        />

        <Option
          n={3}
          name="GORM gen"
          href="https://github.com/go-gorm/gen"
          meta="2,560 stars · active October 2026 · MIT"
          what={
            <p>
              If you already use GORM, gen produces type-safe query code from your models or
              straight from the database, so <code>Where(&quot;naem = ?&quot;)</code> becomes a
              compile error instead of an empty result. The smallest change of the three.
            </p>
          }
          pros={[
            'Adds type safety without leaving GORM',
            'Can generate models from an existing database',
            'Small adoption cost',
          ]}
          cons={[
            'You inherit GORM’s behaviour, good and bad',
            'Much smaller community than sqlc or ent',
          ]}
          bestFor="Existing GORM codebases that want fewer stringly-typed queries."
        />

        <Option
          n={4}
          name="goa"
          href="https://github.com/goadesign/goa"
          meta="6,123 stars · active November 2026 · MIT"
          what={
            <p>
              Design-first. Write the API in goa&apos;s DSL and it generates the server scaffolding,
              the client, the OpenAPI document and the transport code for HTTP and gRPC. The
              contract is the source of truth and the code follows it.
            </p>
          }
          pros={[
            'One description produces server, client and documentation',
            'The spec cannot drift from the implementation',
            'HTTP and gRPC from the same design',
          ]}
          cons={[
            'A DSL to learn, and the team has to edit it first every time',
            'Heavy for a small service',
            'Nothing about the database or the frontend',
          ]}
          bestFor="APIs with several consumers where the contract has to be enforced."
        />

        <Option
          n={5}
          name="Full-stack scaffolders"
          meta="Grit, Buffalo, PocketBase · see the full-stack comparison"
          what={
            <p>
              The Rails sense of the word: one command writes the model, the service, the handlers,
              the routes, the validation and a screen. Buffalo does this for server-rendered Go,
              PocketBase effectively does it through its dashboard, and Grit does it through to a
              React admin.
            </p>
          }
          pros={[
            'The widest reach: database to user interface in one command',
            'Conventions that make a team’s code look alike',
            'Removes the work nobody budgets for',
          ]}
          cons={[
            'You adopt a project layout, not just a library',
            'Regeneration is a workflow to learn',
            'Hardest to retrofit into an existing codebase',
          ]}
          bestFor="New applications where the user interface is part of the deliverable."
        />
      </Section>

      <Section id="choose" title="How to choose">
        <H3>Decide what the source of truth is</H3>
        <p>
          SQL, a Go schema, a design document, or a model. Every tool here is a bet on one of those,
          and the bet is hard to change later because everything downstream is generated from it.
        </p>

        <H3>Check the generated code is readable</H3>
        <p>
          You will debug it. sqlc and ent both produce code you can step through, which is the
          strongest argument for generation over reflection: when something is wrong you can open
          the file and look, rather than reasoning about what a library is doing at runtime.
        </p>

        <H3>Check that removal is supported</H3>
        <p>
          Generators are judged on the second day, when a field is renamed or a resource is deleted.
          A tool that only adds leaves orphans: a registry entry pointing at a deleted file, a route
          to a handler that no longer exists. Ask what the undo looks like before you adopt one.
        </p>

        <H3>Do not combine too many</H3>
        <p>
          sqlc plus goa plus a frontend generator is three code-generation steps in one build, three
          sets of conventions, and three things to re-run in the right order. One is usually enough.
        </p>
      </Section>

      <WhereGritFits>
        <p>
          Grit is in the fourth category, the widest one.{' '}
          <code>grit generate resource Order --fields &quot;total:money,status:select[new,paid]&quot;</code>{' '}
          writes the Go model, the service, the handler, the routes, the migration, the Zod
          validation schema, the TypeScript types, the React Query hooks and the admin screen,
          typed against each other.
        </p>
        <p>
          Two details that matter for a generator you will live with. Everything it writes goes into
          your repository and is reviewed in a pull request. And{' '}
          <code>grit remove resource Order</code> takes all of it back out, including the
          injections: the route registration, the sync registry, the permission catalogue, the
          seeder call. A generator without a removal leaves a codebase nobody can tidy.
        </p>
      </WhereGritFits>

      <Faq
        slug="/i/go-crud-generator"
        items={[
          {
            q: 'What is the best CRUD generator for Go?',
            a: (
              <p>
                It depends which job you mean. For SQL to typed Go, sqlc. For a schema producing
                models, queries and migrations, ent. For a contract producing both sides, goa. For a
                model producing an entire feature including the user interface, a full-stack
                scaffolder.
              </p>
            ),
            text:
              'It depends on the job. sqlc for turning SQL into typed Go, ent for a schema that produces models, queries and migrations, goa for a contract that produces server and client, and a full-stack scaffolder for a model that produces an entire feature including the user interface.',
          },
          {
            q: 'Is sqlc better than GORM?',
            a: (
              <p>
                They answer different questions. sqlc gives you exactly the SQL you wrote with
                compile-time checking, and makes dynamic filters awkward. GORM gives you
                associations, hooks and dynamic query building, and makes the SQL less predictable.
                Plenty of codebases use both.
              </p>
            ),
            text:
              'They answer different questions. sqlc gives exactly the SQL you wrote with compile-time checking but makes dynamic filters awkward. GORM gives associations, hooks and dynamic query building at the cost of predictable SQL. Many codebases use both.',
          },
          {
            q: 'Does Go have scaffolding like rails generate?',
            a: (
              <p>
                Not in the standard toolchain. Buffalo has generators for server-rendered Go,
                PocketBase creates collections through its dashboard, and Grit generates from the Go
                model through to a React admin screen.{' '}
                <Link href="/i/go-full-stack-framework">The full comparison is here</Link>.
              </p>
            ),
            text:
              'Not in the standard toolchain. Buffalo provides generators for server-rendered Go, PocketBase creates collections through its dashboard, and Grit generates from the Go model through to a React admin screen.',
          },
          {
            q: 'Should generated code be committed?',
            a: (
              <p>
                Usually yes. It keeps the build reproducible without the generator installed, it
                makes the diff of a schema change visible in review, and it means you can read the
                code you are debugging. The cost is noisier pull requests, which is a smaller
                problem than the one it solves.
              </p>
            ),
            text:
              'Usually yes. Committing generated code keeps builds reproducible without the generator installed, makes schema changes visible in review, and lets you read the code you are debugging. The cost is noisier pull requests.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/i/go-full-stack-framework',
            label: 'Full-stack frameworks for Go',
            blurb: 'The projects that generate more than the database layer.',
          },
          {
            href: '/i/go-react-stack',
            label: 'A Go backend with a React frontend',
            blurb: 'Generating types across the language boundary.',
          },
          {
            href: '/docs/concepts/code-generation',
            label: 'What grit generate writes',
            blurb: 'Every file produced for one resource, annotated.',
          },
          {
            href: '/i/go-admin-panel',
            label: 'Admin panels for Go',
            blurb: 'The screen a full-stack generator produces.',
          },
        ]}
      />
    </InfoPage>
  )
}
