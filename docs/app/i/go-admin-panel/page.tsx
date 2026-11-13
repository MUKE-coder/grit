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
  title: 'Admin panels for Go, compared',
  description:
    'Go has no Filament. GoAdmin, go-admin, QOR, PocketBase, hosted tools and generated admins compared, and the question that decides between them: where your schema lives.',
  alternates: { canonical: 'https://gritframework.dev/i/go-admin-panel' },
  openGraph: {
    title: 'Admin panels for Go: the four approaches compared',
    description: 'Why Go has no Filament, and the four ways teams solve it.',
    url: 'https://gritframework.dev/i/go-admin-panel',
    type: 'article',
  },
}

export default function Page() {
  return (
    <InfoPage
      slug="/i/go-admin-panel"
      eyebrow="Guide"
      title="Admin panels for Go"
      checked="13 November 2026"
      lead={
        <>
          <p className="mb-4">
            Every application needs a screen where somebody internal can find a record, correct it,
            and get on with their day. In PHP that screen is Filament and it takes an afternoon. In
            Go there is no equivalent with that standing, which is why this search exists.
          </p>
          <p>
            There are four ways teams solve it, and they differ on one question that matters more
            than any feature list: <strong>where does the admin&apos;s idea of your data live?</strong>{' '}
            In a config you write, in reflection over your structs, in a separate service, or in
            code generated into your repository. That answer decides what happens when the schema
            changes.
          </p>
        </>
      }
    >
      <Section id="approaches" title="The four approaches">
        <H3>1. A Go admin library</H3>
        <p>
          You register tables and fields with a library that renders server-side HTML. Fast to
          start. The UI is the library&apos;s, and customising a screen means learning its
          extension model.
        </p>

        <H3>2. A general-purpose admin tool pointed at your database</H3>
        <p>
          Retool, Forest Admin, Appsmith and similar connect to Postgres directly. No Go code at
          all, excellent for internal operations, and your business rules are now in two places:
          your API and their query builder.
        </p>

        <H3>3. Build it yourself in React</H3>
        <p>
          What most teams end up doing. Total control, and the honest estimate is a month per
          application for tables with search, filters, sorting, pagination, bulk actions, CSV export
          and a form builder that handles relations.
        </p>

        <H3>4. Generate it from the model</H3>
        <p>
          A command writes the admin screen next to your Go model, into your repository. Regenerate
          when the model changes. You own the output and can edit it; the cost is a generator that
          has to be re-run and a convention to follow.
        </p>
      </Section>

      <Section id="glance" title="The options at a glance">
        <CompareTable
          caption="Go admin panel options compared"
          columns={['Option', 'UI technology', 'Schema lives in', 'Stars', 'Last activity']}
          rows={[
            ['GoAdmin', 'Server HTML, AdminLTE', 'Registered tables', '9.0k', 'June 2025'],
            ['go-admin (team)', 'Vue, separate app', 'Database and config', '12.8k', 'October 2026'],
            ['QOR', 'Server HTML', 'Struct reflection', '5.3k', 'November 2026'],
            ['PocketBase', 'Bundled dashboard', 'The database itself', '61.6k', 'October 2026'],
            ['Retool and similar', 'Hosted, drag and drop', 'Their query builder', 'Commercial', 'Active'],
            ['Grit', 'React, generated', 'Your Go model', 'See repo', 'November 2026'],
          ]}
        />
      </Section>

      <Section id="options" title="The options, in detail">
        <Option
          n={1}
          name="GoAdmin"
          href="https://github.com/GoAdminGroup/go-admin"
          meta="9,016 stars · last push June 2025 · Apache 2.0"
          what={
            <p>
              The best-known Go admin library, explicitly inspired by laravel-admin. You register a
              table, describe its columns and form fields in Go, and it renders a themed admin with
              AdminLTE. It supports several databases and has a plugin system.
            </p>
          }
          pros={[
            'A working admin from a short Go file per table',
            'Several themes and multiple database drivers',
            'No JavaScript build step anywhere',
          ]}
          cons={[
            'No commits since June 2025',
            'The visual style is of its era and not easy to escape',
            'Customising beyond its model means fighting it',
          ]}
          bestFor="An internal tool where the look does not matter and the tables are simple."
        />

        <Option
          n={2}
          name="go-admin (go-admin-team)"
          href="https://github.com/go-admin-team/go-admin"
          meta="12,832 stars · active October 2026 · MIT"
          what={
            <p>
              A different project with a confusingly similar name: a complete back-office
              application rather than a library. Gin and GORM on the server, a Vue and Element UI
              front end, with users, roles, menus and permissions already built. Much of the
              documentation is in Chinese.
            </p>
          }
          pros={[
            'A whole working back office, not a library to assemble',
            'RBAC, menu management and audit logging included',
            'Actively maintained',
          ]}
          cons={[
            'You adopt its project structure wholesale',
            'Documentation is largely Chinese',
            'The Vue front end is a second application to maintain',
          ]}
          bestFor="Teams happy to start from its codebase rather than add an admin to their own."
        />

        <Option
          n={3}
          name="QOR"
          href="https://github.com/qor/qor"
          meta="5,349 stars · active November 2026 · MIT"
          what={
            <p>
              The oldest approach here. QOR reflects over your GORM structs and produces admin
              screens, with a set of adjacent modules for media, workflows and localisation. Being
              reflection-driven, it follows your model automatically.
            </p>
          }
          pros={[
            'Follows your GORM models with little configuration',
            'A family of modules: media, l10n, publishing, workflow',
            'Long history in production',
          ]}
          cons={[
            'A dated interface and a dated extension model',
            'Reflection means surprises are discovered at runtime',
            'Documentation is thin for its surface area',
          ]}
          bestFor="An existing GORM codebase that needs admin screens with minimal new code."
        />

        <Option
          n={4}
          name="PocketBase"
          href="https://pocketbase.io"
          meta="61,597 stars on its GitHub mirror · active, developed on its own git host"
          what={
            <p>
              Here the admin <em>is</em> the product. One Go binary contains SQLite, the API, auth,
              file storage and a dashboard where you define collections, and it alters the real
              tables underneath. Nothing else on this page is as fast to a working admin.
            </p>
          }
          pros={[
            'Download one file and you have an admin, an API and auth',
            'The dashboard is genuinely well made',
            'Realtime subscriptions and file handling included',
          ]}
          cons={[
            'SQLite only',
            'Your schema lives in the database rather than your repository',
            'Adopting it is adopting the whole backend, not just an admin',
          ]}
          bestFor="New internal tools where one binary on one server is the right deployment."
        />

        <Option
          n={5}
          name="Retool, Forest Admin and similar"
          meta="Commercial · hosted or self-hosted · per-seat pricing"
          what={
            <p>
              Point a hosted tool at your Postgres and assemble screens by dragging. No Go involved.
              For operations teams this is frequently the correct answer, and it is the fastest path
              to a usable internal screen that exists.
            </p>
          }
          pros={[
            'Usable screens in an afternoon with no backend work',
            'Non-engineers can maintain them',
            'Connects to the database and your API alike',
          ]}
          cons={[
            'Business rules end up duplicated outside your codebase',
            'Per-seat cost grows with the operations team',
            'Not reviewable in a pull request, and not in your repository',
          ]}
          bestFor="Operations tooling that changes weekly and is not part of the product."
        />

        <Option
          n={6}
          name="Grit"
          href="https://github.com/MUKE-coder/grit"
          meta="Active November 2026 · MIT · this site"
          what={
            <p>
              Ours. <code>grit generate resource Order</code> writes the Go model and the admin
              screen together: a React table with search, filters, sorting, pagination, bulk
              actions, CSV and Excel export, and a form built from the same field definitions. The
              output is TypeScript in your repository, reviewed in a pull request and edited like
              any other file.
            </p>
          }
          pros={[
            'The admin is generated from the Go model, so the two cannot drift',
            'The output is yours: readable, editable, in version control',
            'Postgres, MySQL or SQLite rather than one database',
          ]}
          cons={[
            'React and Next.js or Vite, not server-rendered Go',
            'Regeneration is a command somebody has to run',
            'You adopt a project layout along with the admin',
          ]}
          bestFor="Product teams who need an admin that keeps up with a schema that is still moving."
        />
      </Section>

      <Section id="choose" title="How to choose">
        <H3>Ask who changes the admin, and how often</H3>
        <p>
          If an operations team changes screens weekly and engineers are not involved, a hosted tool
          wins and the duplication is worth it. If the admin is part of the product and moves with
          the schema, you want it in the repository where a pull request can catch a mistake.
        </p>

        <H3>Ask where the schema is allowed to live</H3>
        <p>
          This is the question that bites later. If the admin keeps its own description of your
          tables, you now maintain two, and they drift silently. Reflection and generation both
          avoid that, in opposite ways: reflection reads your structs at runtime, generation writes
          code from them at build time. Generation is the one you can read, review and step through
          in a debugger.
        </p>

        <H3>Be honest about the month</H3>
        <p>
          A table with server-side search, filters that combine, sorting, pagination that survives a
          deep page, row selection, bulk actions, an export that matches the filters, and a form
          that handles relations and file uploads is not a week. Every team that has built one knows
          this and every estimate forgets it.
        </p>
      </Section>

      <WhereGritFits>
        <p>
          Grit exists largely because of this gap. The admin panel is generated from the same Go
          model that defines the table, so adding a column and re-running the generator puts it in
          the list, the filters and the form.
        </p>
        <p>
          What you get per resource: a data table with search, combinable filters, sorting,
          pagination, row selection and bulk actions; create, edit and detail screens built from the
          field types; CSV and Excel import and export; role and permission checks on every route;
          and an audit trail of who changed what.
        </p>
        <p>
          It is the wrong choice if you want server-rendered Go with no JavaScript build, or if a
          hosted tool pointed at Postgres would genuinely do.
        </p>
      </WhereGritFits>

      <Faq
        slug="/i/go-admin-panel"
        items={[
          {
            q: 'Is there a Filament for Go?',
            a: (
              <p>
                Nothing with Filament&apos;s standing. The nearest in spirit are GoAdmin for
                server-rendered screens, PocketBase if its dashboard can be your backend, and Grit
                if you want Filament&apos;s resource-definition model producing a React admin in
                your own repository.
              </p>
            ),
            text:
              'Nothing in Go has the standing Filament has in PHP. The closest options are GoAdmin for server-rendered screens, PocketBase if its bundled dashboard can serve as the backend, and Grit for a Filament-style resource definition generating a React admin into your repository.',
          },
          {
            q: 'Can I use Retool with a Go backend?',
            a: (
              <p>
                Yes. Retool connects to Postgres or to your REST API and does not care what wrote
                it. The trade is that query logic and permissions then live in Retool rather than in
                your codebase, which is fine for operations tooling and awkward for product
                behaviour.
              </p>
            ),
            text:
              'Yes. Retool connects to Postgres or to a REST API regardless of the language behind it. The trade-off is that query logic and permissions then live in Retool rather than your codebase, which suits operations tooling more than product behaviour.',
          },
          {
            q: 'How long does it take to build an admin panel from scratch?',
            a: (
              <p>
                For one resource with search, filters, sorting, pagination, bulk actions, export and
                a form that handles relations and uploads, plan on weeks rather than days, and
                multiply by the number of resources unless you build a generic table first, which is
                itself the month.
              </p>
            ),
            text:
              'For a single resource with search, filters, sorting, pagination, bulk actions, export and a form handling relations and uploads, plan on weeks rather than days, multiplied by the number of resources unless you first build a generic table component.',
          },
          {
            q: 'Should the admin be a separate application?',
            a: (
              <p>
                It depends on who signs in. A separate deployment keeps staff screens off the
                customer-facing surface and lets you lock it to an internal network. One application
                with a route group is simpler to run. Grit supports both, and{' '}
                <Link href="/docs/concepts/architecture-modes">the shapes are described here</Link>.
              </p>
            ),
            text:
              'It depends on who signs in. A separate deployment keeps staff screens off the customer-facing surface and can be restricted to an internal network, while a single application with an admin route group is simpler to operate. Both are valid shapes.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/i/go-full-stack-framework',
            label: 'Full-stack frameworks for Go',
            blurb: 'Where an admin fits in the wider framework question.',
          },
          {
            href: '/i/laravel-alternative-go',
            label: 'Coming to Go from Laravel',
            blurb: 'What transfers, what does not, and what to reach for instead.',
          },
          {
            href: '/docs/admin/resources',
            label: 'How a Grit admin resource is defined',
            blurb: 'The resource definition, the table and the form builder.',
          },
          {
            href: '/compare/laravel',
            label: 'Grit vs Laravel',
            blurb: 'Filament against a generated admin, compared honestly.',
          },
        ]}
      />
    </InfoPage>
  )
}
