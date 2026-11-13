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
  title: 'Full-stack frameworks for Go, compared',
  description:
    'Go has no Rails, and "full-stack" means three different things here. Buffalo, Beego, Bud, templ, Encore, PocketBase and Grit compared, with what each leaves you to build.',
  alternates: { canonical: 'https://gritframework.dev/i/go-full-stack-framework' },
  openGraph: {
    title: 'Full-stack frameworks for Go, compared',
    description:
      'What "full-stack" means in Go, the seven real options, and how to choose between them.',
    url: 'https://gritframework.dev/i/go-full-stack-framework',
    type: 'article',
  },
}

export default function Page() {
  return (
    <InfoPage
      slug="/i/go-full-stack-framework"
      eyebrow="Guide"
      title="Full-stack frameworks for Go, compared"
      checked="13 November 2026"
      lead={
        <>
          <p className="mb-4">
            Most comparisons of this kind are a list of routers with a benchmark attached. That is
            not the question people are asking when they search for a full-stack framework for Go.
            The question is: <em>what do I not have to build?</em>
          </p>
          <p>
            Go has no Rails. It has an excellent standard library, a handful of very good routers,
            and a culture that prefers assembling libraries to adopting a framework. That is a real
            strength and a real cost, and the cost lands on whoever has to produce auth, an admin
            screen, migrations, background jobs and a frontend before anybody can use what they
            built.
          </p>
        </>
      }
    >
      <Section id="what" title="What &quot;full-stack&quot; means in Go">
        <p>
          The phrase covers three different architectures, and most arguments about it are two
          people meaning different ones. Decide which you want before you read any comparison,
          including this one.
        </p>

        <H3>1. Go renders the HTML</H3>
        <p>
          The Rails shape. Templates on the server, forms that post back, a sprinkle of JavaScript.
          In Go this is <code>html/template</code>, or{' '}
          <a href="https://github.com/a-h/templ" target="_blank" rel="noreferrer">
            templ
          </a>{' '}
          for typed components, usually with HTMX for interactivity. Fastest to a working page,
          hardest to reuse when a mobile app turns up.
        </p>

        <H3>2. Go serves an API, JavaScript renders the UI</H3>
        <p>
          The majority of new work. A Go API plus React, Vue or Svelte, talking over JSON. The
          framework question becomes a <em>two</em> framework question, and the interesting problems
          move to the boundary: keeping types in step, auth across two origins, and shipping two
          build artefacts that have to agree.{' '}
          <Link href="/i/go-react-stack">Covered in detail here</Link>.
        </p>

        <H3>3. A batteries-included application framework</H3>
        <p>
          Something that gives you auth, an admin, migrations, jobs and a project layout on the
          first day, whichever of the two above you picked. This is the part Go is thinnest on, and
          the reason this page exists.
        </p>
      </Section>

      <Section id="glance" title="The options at a glance">
        <p>
          Eight projects that genuinely attempt more than routing. Dates are the last commit on the branch you get when you clone, not
          GitHub’s "last pushed", which counts a push to any of a project’s
          branches and reads as much more recent. Stars and dates are a proxy for
          whether you will be alone when something breaks, not for quality.
        </p>
        <CompareTable
          caption="Go full-stack frameworks compared"
          columns={['Project', 'Frontend model', 'Admin UI', 'Stars', 'Last activity']}
          rows={[
            ['Buffalo', 'Go templates', 'no', '8.4k', 'March 2026'],
            ['Beego', 'Go templates', 'partial', '32.4k', 'March 2026'],
            ['Pagoda', 'Go components, HTMX', 'yes', '3.0k', 'October 2026'],
            ['Bud', 'Svelte, server rendered', 'no', '5.6k', 'April 2023'],
            ['templ + HTMX', 'Go templates, typed', 'no', '10.6k', 'November 2026'],
            ['Encore', 'Bring your own', 'no', '12.7k', 'November 2026'],
            ['PocketBase', 'Bring your own', 'yes', '61.6k', 'October 2026'],
            ['Grit', 'React, generated', 'yes', 'See repo', 'November 2026'],
          ]}
        />
        <p className="text-sm">
          &quot;Admin UI&quot; means a working screen for your own tables on day one, not a
          dashboard you assemble. It is the single biggest difference between these projects in
          practice, and only three of the eight have one.
        </p>
      </Section>

      <Section id="options" title="The eight, in detail">
        <Option
          n={1}
          name="Buffalo"
          href="https://github.com/gobuffalo/buffalo"
          meta="8,424 stars · last push March 2026 · MIT"
          what={
            <p>
              The closest thing Go has to Rails, and it says so. Buffalo gives you a generator, a
              project layout, Pop for the database, migrations, asset pipelines and sessions. If you
              want the Rails shape in Go, this is the most complete attempt at it.
            </p>
          }
          pros={[
            'A real convention-over-configuration project layout',
            'Generators for resources, models and migrations',
            'Pop and Fizz give you migrations without choosing a library',
          ]}
          cons={[
            'Development has slowed markedly; last push March 2026',
            'Server-rendered templates first, so a React frontend is your own problem',
            'The asset pipeline predates the current JavaScript tooling world',
          ]}
          bestFor="A server-rendered Go app by a team that likes Rails conventions and accepts a quiet project."
        />

        <Option
          n={2}
          name="Beego"
          href="https://github.com/beego/beego"
          meta="32,435 stars · master last updated March 2026 · Apache 2.0"
          what={
            <p>
              The oldest and largest of these. Beego is a full MVC framework with its own ORM,
              router, cache, task scheduler and a command-line tool. It is genuinely batteries
              included, and it is also an older idea of what a Go framework looks like: a lot of the
              API predates modules, generics and the standard library&apos;s routing improvements.
            </p>
          }
          pros={[
            'The most complete feature list of any Go framework here',
            'Still actively maintained after more than a decade',
            'Very large installed base, especially in China',
          ]}
          cons={[
            'Opinionated in ways that feel dated against modern Go',
            'Its ORM is a world of its own rather than GORM or sqlc',
            'Documentation in English lags the Chinese original',
          ]}
          bestFor="A team that wants one vendor for everything and is comfortable with its conventions."
        />

        <Option
          n={3}
          name="Bud"
          href="https://github.com/livebud/bud"
          meta="5,586 stars · main last updated April 2023 · MIT"
          what={
            <p>
              The most interesting idea on this list. Bud generated the glue between a Go backend
              and a Svelte frontend as you added files, the way Vite generates a dev server. It
              deserved to win. Its main branch has not moved since April 2023, so it belongs here as prior
              art rather than as a recommendation.
            </p>
          }
          pros={[
            'Genuinely novel: generated glue rather than hand-wired plumbing',
            'Excellent developer experience when it worked',
          ]}
          cons={[
            'Its main branch has not moved since April 2023',
            'Svelte only, with the frontend choice baked in',
            'Do not start new work on it',
          ]}
          bestFor="Reading, if you are designing something in this space."
        />

        <Option
          n={4}
          name="Pagoda"
          href="https://github.com/mikestefanello/pagoda"
          meta="2,964 stars · main last updated October 2026 · MIT"
          what={
            <p>
              A starter kit rather than a library: you clone it and it is your application. Echo,
              Ent, SQLite, HTMX and gomponents, with authentication, sessions, email verification,
              a task queue and an admin panel already wired. The admin is the clever part: a
              generated extension reads Ent&apos;s schema graph and produces admin screens for
              every entity, so adding a field to the schema adds it to the admin.
            </p>
          }
          pros={[
            'A working application on the first clone, not a library to assemble',
            'An admin generated from the Ent schema, so the two cannot drift',
            'Type-safe HTML in pure Go with gomponents, and no JavaScript build',
          ]}
          cons={[
            'A starter kit, so you inherit the code rather than upgrade a dependency',
            'SQLite and Ent are baked in, not choices',
            'Its admin is explicitly beta: no filters or sorting yet',
          ]}
          bestFor="A server-rendered Go application where you want the whole thing scaffolded and are happy owning it."
        />

        <Option
          n={5}
          name="templ with HTMX"
          href="https://github.com/a-h/templ"
          meta="10,608 stars · active November 2026 · MIT"
          what={
            <p>
              Not a framework, and that is the point. templ compiles typed components to Go, HTMX
              gives you interactivity without a build step, and you assemble the rest. It has become
              the default answer in the Go community for &quot;I want to ship HTML&quot;, and the
              typed templates are a real improvement on <code>html/template</code>.
            </p>
          }
          pros={[
            'Typed templates, checked at compile time',
            'No JavaScript build, no second deployment artefact',
            'Composes with any router you already use',
          ]}
          cons={[
            'You still build auth, admin, migrations and jobs yourself',
            'A mobile app later means writing the API you skipped',
            'HTMX hits a ceiling on genuinely interactive screens',
          ]}
          bestFor="Content-shaped applications where the UI is forms and tables and will stay that way."
        />

        <Option
          n={6}
          name="Encore"
          href="https://github.com/encoredev/encore"
          meta="12,712 stars · active November 2026 · MPL 2.0"
          what={
            <p>
              Encore is a backend framework built on a strong idea: declare your databases, queues
              and cron jobs in Go, and it provisions and wires them, locally and in the cloud. It
              also gives you distributed tracing out of the box, which nothing else here does.
              Frontend is explicitly not in scope.
            </p>
          }
          pros={[
            'Infrastructure from code, with local and cloud environments',
            'Built-in distributed tracing and a service catalogue',
            'Type-safe service-to-service calls',
          ]}
          cons={[
            'No admin UI and no frontend story',
            'Its model is opinionated enough to be hard to leave',
            'Best value is tied to its own cloud',
          ]}
          bestFor="Backend teams building several services who want infrastructure and observability handled."
        />

        <Option
          n={7}
          name="PocketBase"
          href="https://pocketbase.io"
          meta="61,597 stars on its GitHub mirror · active, now developed on its own git host"
          what={
            <p>
              One Go binary containing SQLite, a REST API, auth, file storage, realtime
              subscriptions and an admin dashboard. You define collections in the dashboard and it
              alters the real tables underneath. For an internal tool or a prototype it is the
              fastest path on this page by a wide margin.
            </p>
          }
          pros={[
            'Download one file, run it, you have a backend with an admin',
            'Realtime subscriptions and 32 OAuth providers included',
            'Extensible in Go or JavaScript',
          ]}
          cons={[
            'SQLite only, which rules out some deployments',
            'Your schema lives in the database, not in your repository',
            'Still pre-1.0 and says so',
          ]}
          bestFor="Internal tools, prototypes, and anything where one file on one server is the right answer."
        />

        <Option
          n={8}
          name="Grit"
          href="https://github.com/MUKE-coder/grit"
          meta="Active November 2026 · MIT · this site"
          what={
            <p>
              Ours, so read it with that in mind. Grit is a Go API (Gin and GORM) plus a React
              frontend and a generated admin panel, from one command.{' '}
              <code>grit generate resource Order</code> writes the Go model, service, handler,
              routes, the Zod schema, the TypeScript types, the React Query hooks and the admin
              screen, all typed against each other.
            </p>
          }
          pros={[
            'An admin panel for your own tables on day one',
            'Types generated from the Go model through to the React form',
            'Auth, RBAC, jobs, storage, migrations and Docker included',
          ]}
          cons={[
            'A larger project than a router: more to learn on day one',
            'Opinionated about Gin, GORM and React',
            'Younger than Beego or Buffalo',
          ]}
          bestFor="Teams that want a Go backend and a real UI for it without assembling the stack first."
        />
      </Section>

      <Section id="choose" title="How to choose">
        <H3>Start with who renders the HTML</H3>
        <p>
          This decides more than anything else on the page. If the answer is Go, you are choosing
          between templ, Buffalo and Beego, and the admin question is yours to solve. If the answer
          is React or Vue, half this list stops applying and the real work moves to the boundary
          between the two.
        </p>

        <H3>Then ask what you are not willing to build</H3>
        <p>
          Be specific about it. Authentication with password reset and sessions is a fortnight.
          Role-based permissions is another. An admin table with search, filters, sorting,
          pagination, bulk actions and CSV export is a month, and it is the part nobody budgets for
          because it is &quot;just CRUD&quot;. A framework is worth adopting exactly when it
          removes work you would otherwise do badly.
        </p>

        <H3>Then check the project is still alive</H3>
        <p>
          Two of the eight above are effectively dormant. For a dependency you will carry for years,
          the date of the last commit tells you more than the star count. Open the repository and
          look before you commit to anything on any comparison page, including this one.
        </p>

        <H3>Finally, consider not choosing one</H3>
        <p>
          Gin or Chi plus a handful of libraries is a legitimate answer, and for a single service
          with no UI it is usually the right one. The framework question only becomes pressing when
          there are screens, users and roles.{' '}
          <Link href="/i/go-web-framework">That comparison is here</Link>.
        </p>
      </Section>

      <WhereGritFits>
        <p>
          Grit is the third shape: a batteries-included application framework, built for the case
          where a Go backend needs a real user interface and an admin.
        </p>
        <p>
          It is the right answer when you would otherwise spend the first month on auth, roles, an
          admin table and the types between Go and TypeScript. It is the wrong answer for a single
          stateless service with no UI, where Gin and four libraries will serve you better, and for
          an internal tool that one SQLite file would solve, where PocketBase is faster.
        </p>
        <p>
          We publish side-by-side comparisons against Rails, Laravel, Django, Next.js, NestJS,
          Encore and the MERN stack, including where each of them wins.
        </p>
      </WhereGritFits>

      <Faq
        slug="/i/go-full-stack-framework"
        items={[
          {
            q: 'Does Go have a Rails equivalent?',
            a: (
              <p>
                Not one with Rails&apos; standing. Buffalo is the closest in design and its
                development has slowed; Beego is the largest and is a different, older idea of a
                framework. The Go community has largely preferred composing libraries, which is why
                the gap persists.
              </p>
            ),
            text:
              'Not one with the standing Rails has. Buffalo is the closest in design but its development has slowed, and Beego is the largest but reflects an older idea of a framework. The Go community generally prefers composing libraries, which is why the gap persists.',
          },
          {
            q: 'Is Gin a full-stack framework?',
            a: (
              <p>
                No. Gin is an HTTP router with middleware, and an excellent one. It does not give
                you auth, an admin, migrations, background jobs or a frontend. Several projects on
                this page, Grit included, are built on top of it.
              </p>
            ),
            text:
              'No. Gin is an HTTP router with middleware. It does not provide authentication, an admin panel, migrations, background jobs or a frontend. Several full-stack projects are built on top of Gin.',
          },
          {
            q: 'Can I use Go with React?',
            a: (
              <p>
                Yes, and it is the most common shape for new Go web work. The problems to solve are
                keeping types in step across the boundary, authentication across two origins, and
                whether you ship one artefact or two.{' '}
                <Link href="/i/go-react-stack">We go through each of them here</Link>.
              </p>
            ),
            text:
              'Yes, and it is the most common shape for new Go web work. The problems to solve are keeping types in step across the boundary, authentication across two origins, and whether you deploy one artefact or two.',
          },
          {
            q: 'What is the fastest way to get an admin panel for a Go app?',
            a: (
              <p>
                PocketBase if a single SQLite file suits you, since the admin is the product. Grit
                if you need Postgres and the admin generated into your own repository.{' '}
                <Link href="/i/go-admin-panel">The full comparison is here</Link>.
              </p>
            ),
            text:
              'PocketBase if a single SQLite file suits the deployment, because its admin dashboard is the product. Grit if you need Postgres and want the admin generated into your own repository. Older options include GoAdmin and QOR.',
          },
          {
            q: 'Should I just use the standard library?',
            a: (
              <p>
                For a service with no user interface, very often yes. Go 1.22 added method and
                wildcard routing to <code>net/http</code>, which removed the main reason to reach
                for a router at all. The calculation changes when you have screens, users and roles.
              </p>
            ),
            text:
              'For a service with no user interface, often yes. Go 1.22 added method matching and wildcard routing to net/http, which removed the main reason to reach for a third-party router. The calculation changes once you have screens, users and roles.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/i/go-web-framework',
            label: 'Go web frameworks compared',
            blurb: 'net/http, Chi, Gin, Echo and Fiber, and where the router ends.',
          },
          {
            href: '/i/go-admin-panel',
            label: 'Admin panels for Go',
            blurb: 'Why Go has no Filament, and the four ways teams solve it.',
          },
          {
            href: '/i/go-react-stack',
            label: 'A Go backend with a React frontend',
            blurb: 'Types, auth, CORS and deployment across the boundary.',
          },
          {
            href: '/compare/encore',
            label: 'Grit vs Encore',
            blurb: 'The closest comparison we have: two Go frameworks, different problems.',
          },
        ]}
      />
    </InfoPage>
  )
}
