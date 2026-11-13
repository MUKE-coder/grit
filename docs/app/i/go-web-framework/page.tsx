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
  title: 'Go web frameworks: Gin, Echo, Fiber and Chi compared',
  description:
    'Four of the five are routers, not frameworks, and that decides your next six months. What each gives you, what none of them give you, and how to choose.',
  alternates: { canonical: 'https://gritframework.dev/i/go-web-framework' },
  openGraph: {
    title: 'Go web frameworks compared: net/http, Chi, Gin, Echo, Fiber',
    description: 'Where the router ends and the framework begins, and how to choose.',
    url: 'https://gritframework.dev/i/go-web-framework',
    type: 'article',
  },
}

export default function Page() {
  return (
    <InfoPage
      slug="/i/go-web-framework"
      eyebrow="Guide"
      title="Go web frameworks: net/http, Chi, Gin, Echo and Fiber"
      checked="13 November 2026"
      lead={
        <>
          <p className="mb-4">
            Every comparison of these ends in a benchmark, and the benchmark is almost never what
            decides the project. On any workload that touches a database, the difference between
            these five is smaller than one missing index.
          </p>
          <p>
            The decision that matters is a different one. Four of the five are{' '}
            <strong>routers</strong>: they match a URL to a function and run middleware. One of them
            is the standard library doing the same thing. None is a framework in the sense Rails or
            Laravel are, and knowing exactly where that line falls is worth more than any
            requests-per-second figure.
          </p>
        </>
      }
    >
      <Section id="glance" title="At a glance">
        <CompareTable
          caption="Go routers and web frameworks compared"
          columns={['Project', 'Built on net/http', 'Stars', 'Last activity', 'Shape']}
          rows={[
            ['net/http', 'yes', 'Standard library', 'Every Go release', 'Router and server'],
            ['Chi', 'yes', '23.1k', 'November 2026', 'Router, 100% compatible'],
            ['Gin', 'yes', '89.5k', 'October 2026', 'Router plus helpers'],
            ['Echo', 'yes', '32.8k', 'November 2026', 'Router plus more helpers'],
            ['Fiber', 'no', '40.3k', 'November 2026', 'Router on fasthttp'],
          ]}
        />
        <p className="text-sm">
          The <em>Built on net/http</em> column is the most consequential one in the table. It
          decides whether the ecosystem&apos;s middleware, test helpers and tracing libraries work
          with your choice without an adapter.
        </p>
      </Section>

      <Section id="options" title="The five">
        <Option
          n={1}
          name="net/http (standard library)"
          href="https://pkg.go.dev/net/http"
          meta="Standard library · ships with Go · BSD"
          what={
            <p>
              Since Go 1.22 the standard multiplexer matches methods and path wildcards:{' '}
              <code>mux.HandleFunc(&quot;GET /orders/{'{'}id{'}'}&quot;, h)</code>. That removed the
              main reason most projects reached for a router at all. It has no route groups and no
              request binding, and for a service with a dozen endpoints you will not miss them.
            </p>
          }
          pros={[
            'No dependency to audit, upgrade or outlive',
            'Every middleware and test helper in Go works with it',
            'Method and wildcard routing since Go 1.22',
          ]}
          cons={[
            'No route groups, so shared middleware is manual',
            'No binding or validation helpers',
            'Error handling is yours to design',
          ]}
          bestFor="Services with a small, stable set of endpoints, and anyone minimising dependencies."
        />

        <Option
          n={2}
          name="Chi"
          href="https://github.com/go-chi/chi"
          meta="23,112 stars · active November 2026 · MIT"
          what={
            <p>
              A router and nothing else, deliberately. Chi&apos;s handlers are plain{' '}
              <code>http.HandlerFunc</code>, so anything written for the standard library works
              unchanged. It adds the two things the standard mux lacks: nestable route groups and a
              good middleware stack.
            </p>
          }
          pros={[
            'Completely compatible with net/http, no adapters anywhere',
            'Nested routers and middleware that compose cleanly',
            'Small, stable and quiet in the best way',
          ]}
          cons={[
            'No binding, validation or rendering helpers by design',
            'You write more boilerplate per handler than with Gin',
          ]}
          bestFor="Teams who want the standard library plus route groups, and nothing else."
        />

        <Option
          n={3}
          name="Gin"
          href="https://github.com/gin-gonic/gin"
          meta="89,541 stars · active October 2026 · MIT"
          what={
            <p>
              The default. Gin adds a context object with binding, validation and JSON rendering on
              top of <code>net/http</code>. <code>c.ShouldBindJSON(&amp;req)</code> with{' '}
              <code>binding:&quot;required&quot;</code> struct tags removes a genuinely tedious
              amount of code, and the ecosystem of Gin middleware is the largest in Go.
            </p>
          }
          pros={[
            'The largest community and middleware ecosystem in Go',
            'Binding and validation from struct tags',
            'Still net/http underneath, so standard middleware adapts easily',
          ]}
          cons={[
            'Its own Context type, so handlers are Gin-shaped',
            'Error handling conventions are looser than Echo’s',
          ]}
          bestFor="Most API work. If you have no strong opinion, this is the safe default."
        />

        <Option
          n={4}
          name="Echo"
          href="https://github.com/labstack/echo"
          meta="32,821 stars · active November 2026 · MIT"
          what={
            <p>
              Gin&apos;s closest peer, with a tidier error model. Handlers return an{' '}
              <code>error</code>, which centralises error handling instead of scattering{' '}
              <code>c.JSON(500, ...)</code> through the codebase. It ships more in the box:
              binding, validation, JWT, rate limiting and more.
            </p>
          }
          pros={[
            'Handlers return errors, so one place decides the response',
            'More built in than Gin, less to go and find',
            'Good HTTP/2 and automatic TLS support',
          ]}
          cons={[
            'A smaller middleware ecosystem than Gin',
            'Some breaking changes between major versions',
          ]}
          bestFor="Teams who want Gin’s conveniences with a stricter error-handling convention."
        />

        <Option
          n={5}
          name="Fiber"
          href="https://github.com/gofiber/fiber"
          meta="40,299 stars · active November 2026 · MIT"
          what={
            <p>
              Fiber is built on fasthttp rather than <code>net/http</code>, and that one fact is the
              whole decision. It is genuinely fast, the API is familiar to anyone from Express, and
              in exchange nothing in the standard ecosystem works without an adapter: middleware,
              <code> httptest</code>, many tracing libraries, and some HTTP/2 features.
            </p>
          }
          pros={[
            'Very fast on benchmarks, with low allocation counts',
            'An Express-shaped API that Node developers read immediately',
            'A large collection of first-party middleware',
          ]}
          cons={[
            'Not net/http, so the standard ecosystem needs adapters',
            'httptest does not work directly, which changes how you test',
            'The performance gain is usually invisible behind a database',
          ]}
          bestFor="Throughput-bound services that do little I/O, where the benchmark is the requirement."
        />
      </Section>

      <Section id="gap" title="What none of them give you">
        <p>
          This is the part a router comparison leaves out, and it is where the next six months go.
          Having chosen any of the five above, you still need:
        </p>
        <ul className="ml-5 list-disc space-y-2">
          <li>
            <strong>Authentication.</strong> Registration, login, password hashing and reset,
            sessions or JWTs, refresh and revocation.
          </li>
          <li>
            <strong>Authorization.</strong> Roles, permissions, and the per-record checks that
            decide which rows a user may see.
          </li>
          <li>
            <strong>Migrations.</strong> A tool, a convention, and a way to run them on deploy.
          </li>
          <li>
            <strong>An admin interface.</strong> Tables with search, filters, sorting, pagination,
            bulk actions and export. Always estimated at a week and never finished in one.
          </li>
          <li>
            <strong>Background jobs and scheduling.</strong> A queue, retries, dead letters and
            somewhere to see what failed.
          </li>
          <li>
            <strong>File storage.</strong> Uploads, validation, an object store and signed URLs.
          </li>
          <li>
            <strong>The frontend, and its types.</strong> Whatever renders the UI, and some way of
            keeping its idea of your data in step with Go&apos;s.
          </li>
        </ul>
        <p>
          None of that is hard. All of it is weeks, and all of it is work every team redoes. That is
          the gap a full-stack framework exists to fill, and{' '}
          <Link href="/i/go-full-stack-framework">the options for filling it are here</Link>.
        </p>
      </Section>

      <Section id="choose" title="How to choose">
        <H3>Do not start with the benchmark</H3>
        <p>
          A framework&apos;s routing overhead is measured in microseconds. One unindexed query is
          measured in milliseconds, and one unnecessary round trip to a database across a network is
          worse than both. If throughput is genuinely your constraint, measure your own workload;
          the published numbers disagree with each other and with themselves between versions.
        </p>

        <H3>Decide whether you are leaving net/http</H3>
        <p>
          Fiber on one side, everything else on the other. Leaving the standard interface is a
          defensible trade for a narrow, throughput-bound service and an expensive one for a
          long-lived application, because every library you adopt afterwards has to be checked for
          compatibility.
        </p>

        <H3>Then pick on ergonomics, and stop</H3>
        <p>
          Between Chi, Gin and Echo there is no wrong answer. Chi if you want the standard library
          with groups. Gin if you want binding and the biggest ecosystem. Echo if returning errors
          from handlers appeals. Spend the time you saved on the list in the previous section.
        </p>
      </Section>

      <WhereGritFits>
        <p>
          Grit is not an alternative to Gin: it is built <em>on</em> Gin, and on GORM. The question
          it answers is the one after the router.
        </p>
        <p>
          <code>grit new my-app</code> gives you a Gin API with authentication, roles and
          permissions, migrations, background jobs, file storage and an admin panel already wired
          together, plus a React frontend whose types come from the Go models.{' '}
          <code>grit generate resource Order</code> then writes the model, service, handler, routes,
          validation schema, TypeScript types, data-fetching hooks and admin screen as one typed
          unit.
        </p>
        <p>
          If your project is a single service with no user interface, keep Gin and four libraries.
          That is a good answer and this is not aimed at it.
        </p>
      </WhereGritFits>

      <Faq
        slug="/i/go-web-framework"
        items={[
          {
            q: 'Which Go web framework is fastest?',
            a: (
              <p>
                Fiber usually leads synthetic benchmarks because fasthttp allocates less per
                request. Behind a database the difference is almost always invisible. Published
                figures disagree with each other and change between versions, so measure your own
                workload before treating speed as the deciding factor.
              </p>
            ),
            text:
              'Fiber usually leads synthetic benchmarks because fasthttp allocates less per request. Behind a database the difference is almost always invisible, and published figures disagree with each other, so measure your own workload before treating speed as the deciding factor.',
          },
          {
            q: 'Do I still need a router in Go 1.22 and later?',
            a: (
              <p>
                Often not. The standard multiplexer now matches methods and path wildcards, which
                covers most services. Reach for Chi when you want nested route groups, or Gin and
                Echo when you want request binding and validation helpers.
              </p>
            ),
            text:
              'Often not. Since Go 1.22 the standard multiplexer matches methods and path wildcards, which covers most services. Chi adds nested route groups; Gin and Echo add request binding and validation helpers.',
          },
          {
            q: 'What is the difference between Gin and Echo?',
            a: (
              <p>
                Mostly error handling and how much ships in the box. Echo handlers return an{' '}
                <code>error</code>, so one handler decides every error response; Gin handlers write
                the response themselves. Gin has the larger ecosystem, Echo has more built in.
              </p>
            ),
            text:
              'Mostly error handling and how much ships in the box. Echo handlers return an error so one place decides every error response, while Gin handlers write the response themselves. Gin has the larger middleware ecosystem and Echo includes more out of the box.',
          },
          {
            q: 'Is Fiber production ready?',
            a: (
              <p>
                Yes, and it is widely used. The real question is not stability but compatibility:
                because it is not <code>net/http</code>, standard middleware, <code>httptest</code>{' '}
                and some observability tooling need adapters or alternatives.
              </p>
            ),
            text:
              'Yes, and it is widely used. The question is compatibility rather than stability: because Fiber is built on fasthttp rather than net/http, standard middleware, httptest and some observability tooling need adapters or alternatives.',
          },
          {
            q: 'Is Gin a full-stack framework?',
            a: (
              <p>
                No. Gin routes requests and helps you bind them. Authentication, permissions,
                migrations, jobs, an admin interface and a frontend are all yours.{' '}
                <Link href="/i/go-full-stack-framework">
                  The projects that cover those are compared here
                </Link>
                .
              </p>
            ),
            text:
              'No. Gin routes requests and helps bind them. Authentication, permissions, migrations, background jobs, an admin interface and a frontend are not included and remain your responsibility.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/i/go-full-stack-framework',
            label: 'Full-stack frameworks for Go',
            blurb: 'The seven projects that cover what a router does not.',
          },
          {
            href: '/i/go-react-stack',
            label: 'A Go backend with a React frontend',
            blurb: 'Types, auth, CORS and deployment across the boundary.',
          },
          {
            href: '/i/go-crud-generator',
            label: 'CRUD and scaffolding generators for Go',
            blurb: 'sqlc, ent, GORM gen, goa, and what each one actually writes.',
          },
          {
            href: '/docs/backend/handlers',
            label: 'How Grit writes a handler',
            blurb: 'The generated Gin handler, service and routes, annotated.',
          },
        ]}
      />
    </InfoPage>
  )
}
