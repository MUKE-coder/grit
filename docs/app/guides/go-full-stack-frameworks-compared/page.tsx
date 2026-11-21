import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Choosing a Go full-stack approach in 2026',
  description:
    'Runtime backend, infrastructure from code, code generator or library assembly: four different bargains, what each costs when you outgrow it, and when not to pick Grit.',
  alternates: { canonical: 'https://gritframework.dev/guides/go-full-stack-frameworks-compared' },
  openGraph: {
    title: 'Choosing a Go full-stack approach in 2026',
    description: 'Four kinds of answer, and the exit cost of each.',
    url: 'https://gritframework.dev/guides/go-full-stack-frameworks-compared',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="go-full-stack-frameworks-compared"
      lead={
        <>
          <p className="mb-4">
            &quot;Which Go framework&quot; is the wrong first question, because the options are
            not competitors. PocketBase, Encore, Goravel and a code generator are four different
            kinds of thing, and most comparison tables line them up as though choosing between
            them were a matter of features.
          </p>
          <p>
            The useful question is which bargain you are making, and what it costs to leave when
            the product outgrows it. That cost is the thing nobody puts in the table, and it is
            the thing you will care about in year two.
          </p>
        </>
      }
    >
      <Outcome>
        <p>Four categories, and which questions put you in each.</p>
        <p>The exit cost of each, stated plainly.</p>
        <p>A section on when not to pick Grit, which is a real list.</p>
      </Outcome>

      <Note title="For a per-tool comparison">
        <p>
          This page is about the decision. For eight specific tools, Buffalo, Beego, Bud, templ,
          Encore, PocketBase and the rest, with what each one actually includes,{' '}
          <Link href="/i/go-full-stack-framework" className="text-accent hover:underline">
            the options are compared individually here
          </Link>
          .
        </p>
      </Note>

      <Section id="kinds" title="Four kinds of answer">
        <H3>1. A runtime backend</H3>
        <p>
          PocketBase, Supabase, Appwrite. You get a running server with a database, auth, file
          storage, realtime and an admin interface, and you write against its API. The product is
          the server.
        </p>
        <p>
          <strong>The bargain:</strong> almost nothing to build, and your application lives inside
          somebody else&apos;s process model. Custom logic goes in hooks or an extension
          mechanism rather than in your own request handlers. PocketBase is a single binary with
          embedded SQLite, which is genuinely excellent for a tool, an internal application or a
          product that stays small.
        </p>
        <p>
          <strong>The exit cost:</strong> high, and back-loaded. The data is portable; the
          application is not, because your logic is written in terms of the platform&apos;s
          extension points. Leaving means rewriting the backend.
        </p>
        <H3>2. Infrastructure from code</H3>
        <p>
          Encore. You declare a database, a queue, a cron job in Go, and the framework provisions
          and wires them, locally and in the cloud, with tracing and dashboards.
        </p>
        <p>
          <strong>The bargain:</strong> an enormous amount of infrastructure work disappears, and
          your code is written against the framework&apos;s declarations. It is a strong idea
          executed well, and it is most valuable when you have several services and nobody who
          enjoys Terraform.
        </p>
        <p>
          <strong>The exit cost:</strong> moderate to high. The business logic is ordinary Go; the
          wiring is not, and reproducing what the framework was doing for you is the work.
        </p>
        <H3>3. A code generator</H3>
        <p>
          Grit, and in other ecosystems Rails scaffolding or Laravel with Filament. The tool
          writes ordinary code into your repository: models, handlers, routes, frontend, admin
          screens. Then it gets out of the way.
        </p>
        <p>
          <strong>The bargain:</strong> you own everything it wrote, so there is no ceiling on
          what you can change, and there is also no magic doing anything for you at runtime. What
          you give up is that generated code is code you did not write and must still review,
          understand and maintain.
        </p>
        <p>
          <strong>The exit cost:</strong> low by construction. Stop running the generator and you
          have a Gin and GORM application, which is what you would have written. &quot;Ejecting&quot;
          is not a feature because there is nothing to eject from.
        </p>
        <H3>4. Libraries, assembled</H3>
        <p>
          Gin or Echo, GORM or sqlc, your own structure. The default, and for many services the
          right one.
        </p>
        <p>
          <strong>The bargain:</strong> complete control, and you build auth, the admin,
          migrations, file uploads, background jobs and the frontend plumbing. Every one of those
          is a known quantity and a week or two, and they add up to a quarter you did not spend
          on the product.
        </p>
        <p>
          <strong>The exit cost:</strong> none, and that is the point. The risk is not lock-in, it
          is the five identical plumbing problems each solved slightly differently by whoever got
          there first.
        </p>
      </Section>

      <Section id="decide" title="The questions that actually decide it">
        <Checklist
          items={[
            <>
              <strong>Will you need custom business logic in the request path?</strong> If yes, a
              runtime backend will be fighting you by month three. If no, it may be the fastest
              route to a working product by a wide margin.
            </>,
            <>
              <strong>Does an operator need to administer the data?</strong> Almost every product
              eventually says yes, and the admin panel is where unbudgeted weeks go. A database
              browser is not an admin panel: it has no business rules, no validation and no audit
              trail.
            </>,
            <>
              <strong>Who maintains this in three years?</strong> Generated plain Go can be
              maintained by any Go developer. A framework-specific runtime requires somebody to
              learn the framework, and a popular framework is a much safer bet than an elegant
              one with four contributors.
            </>,
            <>
              <strong>How much infrastructure are you willing to operate?</strong> One SQLite
              file, one Postgres and one Redis, or a cluster. This narrows the options faster
              than any feature does.
            </>,
            <>
              <strong>Is the frontend part of the same decision?</strong> If you have a React team
              and a React application already, anything that wants you to write Go templates is
              not in the running whatever else it offers.
            </>,
          ]}
        />
      </Section>

      <Section id="not-grit" title="When not to pick Grit">
        <p>
          A list like this is worthless if it concludes that the tool you are reading about is
          always right. These are the cases where it is not.
        </p>
        <Checklist
          items={[
            <>
              <strong>You are building a library or a single-purpose service.</strong> A metrics
              collector, a webhook relay, a CLI. There is no frontend, no admin and no CRUD. Use
              the standard library and one router.
            </>,
            <>
              <strong>Your frontend is not React.</strong> The generated frontend is Next.js or a
              Vite React application. With Vue, Svelte or server-rendered Go templates, you are
              using the backend half and writing the rest, which is a smaller benefit than the
              page suggests.
            </>,
            <>
              <strong>You want no generated code in your repository.</strong> Some teams hold
              that position and it is coherent. A generator is the wrong shape of tool for it.
            </>,
            <>
              <strong>You already have a large, working Go codebase.</strong> Adopting a
              generator mid-life means reconciling its conventions with yours. Take the
              individual pieces if they help, and do not restructure a working system around a
              tool.
            </>,
            <>
              <strong>Your data model is not relational.</strong> Event-sourced, document-shaped
              or graph-shaped domains do not fit a resource generator built on GORM, and forcing
              them to is worse than not using it.
            </>,
            <>
              <strong>You need a very small deployment footprint.</strong> A single binary with
              embedded SQLite and an admin interface is PocketBase, and nothing here beats it on
              that specific axis.
            </>,
          ]}
        />
        <Warn title="And the honest risk">
          <p>
            Grit is young and maintained by a small team. That is a real consideration, and the
            mitigation is structural rather than promissory: it generates ordinary Gin and GORM
            code into your repository, so a project survives the tool. If the generator stopped
            being maintained tomorrow, you would have a Go application that builds, runs and can
            be changed by anybody who knows Go. That is the main argument for the generator shape
            over the runtime shape, and it applies to the risk as much as to the ceiling.
          </p>
        </Warn>
      </Section>

      <Section id="lock-in" title="Measuring lock-in properly">
        <p>
          &quot;Lock-in&quot; is used loosely. A useful definition: how much of what you have
          written would have to be rewritten if the tool disappeared.
        </p>
        <CodeBlock
          language="text"
          code={`Libraries assembled     nothing. You wrote it.
Code generator          nothing. It wrote ordinary code, in your repo.
Infrastructure-as-code  the wiring, the deployment, the observability.
Runtime backend         the backend. The data exports; the logic does not.`}
        />
        <p>
          That ordering is not a ranking. A high exit cost can be entirely worth paying: if a
          runtime backend takes your product to market three months earlier, and most products
          never need to leave, the bargain was good. The mistake is making it without noticing.
        </p>
        <p>
          The honest test for any tool: <strong>what does the code look like if I stop using
          it?</strong> Ask it before adopting, not during the migration.
        </p>
      </Section>

      <Section id="benchmarks" title="Why performance is not in this page">
        <p>
          All of these are Go, and Go is fast enough that the framework is not what makes your
          application slow. Your database is, and specifically the queries you wrote against it.
        </p>
        <p>
          Router benchmarks measure a hello-world handler over a loopback interface. The
          difference between the fastest and the slowest is tens of microseconds, against a
          database query of tens of milliseconds, which makes it about a thousandth of the
          response. Choosing on that number is choosing on noise.
        </p>
        <p>
          What does differ, and is worth asking about: how many queries the default list endpoint
          makes, whether pagination is on by default, whether a relation is preloaded or produces
          an N+1, and whether there is a profiler to answer those questions when something is
          slow.
        </p>
      </Section>

      <Section id="try" title="Spend an afternoon rather than a week reading">
        <p>
          Each of these can be stood up in under an hour. Build the same small thing in two of
          them: one resource with a relationship, a list with a filter, an admin screen, and a
          deploy. You will learn more from that than from any comparison, including this one.
        </p>
        <CodeBlock
          terminal
          code={`# Grit
grit new demo && cd demo && grit start

# PocketBase
./pocketbase serve

# Encore
encore app create demo && encore run`}
        />
        <p>Three things to notice while doing it, because they are what you will live with:</p>
        <Checklist
          items={[
            <>
              <strong>What happens when you need something the tool did not anticipate.</strong>{' '}
              Is there a door, and where does it lead?
            </>,
            <>
              <strong>What the error messages are like when you get it wrong.</strong> You will
              spend more time with those than with the documentation.
            </>,
            <>
              <strong>Whether you can read the code that runs your application.</strong> Not
              whether it is pretty: whether it exists somewhere you can open it.
            </>,
          ]}
        />
      </Section>

      <Faq
        slug="go-full-stack-frameworks-compared"
        items={[
          {
            q: 'Is there a Rails or Laravel for Go?',
            a: (
              <>
                Not in the sense of one framework the community agrees on. Go&apos;s culture
                pushes against large frameworks: the standard library is capable, the ecosystem
                favours small composable packages, and the result is several partial answers
                rather than one dominant one. Goravel is the closest attempt at a deliberate
                Laravel port; Buffalo was the closest in spirit and has slowed.
              </>
            ),
            text:
              'Not in the sense of one framework the community agrees on. Go culture favours a capable standard library and small composable packages, so there are several partial answers rather than one dominant one. Goravel is the closest deliberate Laravel port, and Buffalo was closest in spirit before slowing.',
          },
          {
            q: 'Is PocketBase good enough for production?',
            a: (
              <>
                Yes, for the shape it fits: one binary, SQLite, an admin interface included,
                vertical scaling. It runs real products. The limits are the ones that follow from
                that shape, a single writer and a single machine, and the fact that your custom
                logic lives in its extension model rather than in your own handlers. Both are
                acceptable until they are not, and the point at which they stop being acceptable
                is predictable.
              </>
            ),
            text:
              'Yes, for the shape it fits: one binary, SQLite, an included admin interface and vertical scaling. The limits follow from that shape, a single writer and a single machine, and from custom logic living in its extension model rather than your own handlers.',
          },
          {
            q: 'Why is generated code better than a framework doing it at runtime?',
            a: (
              <>
                It is not better in every way; it is a different trade. Generated code is
                readable, debuggable with ordinary tools and changeable without permission, and
                it is also code you did not write that you now maintain, and it does not improve
                when the tool does unless you regenerate. Runtime magic is less to maintain and
                harder to see into. Pick based on which failure you would rather have at three in
                the morning.
              </>
            ),
            text:
              'It is a different trade rather than a better one. Generated code is readable, debuggable with ordinary tools and changeable without permission, but it is code you now maintain and it does not improve when the tool does unless regenerated. Runtime behaviour is less to maintain and harder to see into.',
          },
          {
            q: 'What if I pick wrong?',
            a: (
              <>
                For most of these the answer is a migration rather than a catastrophe, and the
                cost is the exit cost above. The decision that is genuinely hard to reverse is
                not the framework, it is the data model, and in a multi-tenant product
                specifically it is whether a tenant is a user or an organisation. Spend your
                anxiety there.
              </>
            ),
            text:
              'For most of these it is a migration rather than a catastrophe, and the cost is the exit cost described above. The genuinely hard decision to reverse is the data model, and in a multi-tenant product specifically whether a tenant is a user or an organisation.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/i/go-full-stack-framework',
            label: 'The eight options, individually',
            blurb: 'Buffalo, Beego, Bud, templ, Encore, PocketBase and the rest, in detail.',
          },
          {
            href: '/guides/go-vs-nodejs-backend',
            label: 'Go or Node.js for your backend',
            blurb: 'The decision before this one, if it is not settled.',
          },
          {
            href: '/guides/go-project-structure',
            label: 'Go project structure',
            blurb: 'What you are choosing between if you assemble it yourself.',
          },
          {
            href: '/docs/getting-started/philosophy',
            label: 'What Grit is and is not',
            blurb: 'The reasoning behind the generator shape, from the inside.',
          },
        ]}
      />
    </GuidePage>
  )
}
