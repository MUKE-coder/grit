import Link from 'next/link'
import type { Metadata } from 'next'
import {
  GuidePage,
  Section,
  H3,
  Outcome,
  Note,
  Warn,
  Checklist,
  Faq,
  Related,
} from '@/components/guide'
import { CodeBlock } from '@/components/code-block'
import { DocsShot } from '@/components/docs-shot'

export const metadata: Metadata = {
  title: 'Build an admin dashboard without writing it from scratch',
  description:
    'What an internal admin panel actually needs, what it costs to hand-build, and how to describe a screen instead of writing one.',
  alternates: { canonical: 'https://gritframework.dev/guides/build-admin-dashboard-go-react' },
  openGraph: {
    title: 'Build an admin dashboard without writing it from scratch',
    description: 'Tables, forms, filters and widgets, and what to declare rather than build.',
    url: 'https://gritframework.dev/guides/build-admin-dashboard-go-react',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="build-admin-dashboard-go-react"
      lead={
        <>
          <p className="mb-4">
            Every application grows an admin panel, because somebody eventually has to look up a
            customer, refund an order or fix a typo in production. It is also the part of the
            product nobody budgets for, so it gets built in a hurry, three times, slightly
            differently.
          </p>
          <p>
            This is what an admin screen genuinely needs, what each part costs to build by hand,
            and the alternative: describing the screen in a file and letting the table, the form
            and the filters follow from that description.
          </p>
        </>
      }
    >
      <Outcome>
        <p>A realistic estimate of what a hand-built admin costs, and which parts are worth it.</p>
        <p>
          A working understanding of a resource definition: one file that produces the list, the
          form, the filters and the detail view at once.
        </p>
        <p>The four decisions that separate a usable admin from a data dump.</p>
      </Outcome>

      <Section id="what" title="What an admin panel actually needs">
        <p>
          Not &quot;CRUD&quot;. CRUD is four verbs, and an admin panel that only does those four
          verbs is unusable past about fifty rows. The real list:
        </p>
        <Checklist
          items={[
            <>
              <strong>A table that is a scanning surface.</strong> Sorting by any column, filtering
              per column, text search, pagination that works at a million rows, and a column set
              chosen for scanning rather than every field the model has.
            </>,
            <>
              <strong>Forms that validate the way the API does.</strong> Same rules, same messages,
              so the form never accepts what the server then rejects with a 422 nobody can explain.
            </>,
            <>
              <strong>Relationships that are usable.</strong> A dropdown showing a customer&apos;s
              name and not their id, searchable, because a select with four thousand options is not
              a control.
            </>,
            <>
              <strong>Bulk actions.</strong> The reason anyone opens an admin panel is usually to
              do one thing to forty rows.
            </>,
            <>
              <strong>Dates and money formatted correctly,</strong> in the viewer&apos;s timezone
              and the record&apos;s currency, with thousands separators.{' '}
              <code>50000000</code> is unreadable and <code>50,000,000</code> is not.
            </>,
            <>
              <strong>Every state designed.</strong> Loading, empty, error, and far too much data.
              The empty state should say what to do, not &quot;No data&quot;.
            </>,
            <>
              <strong>An audit trail.</strong> The first serious question an admin panel is asked
              is &quot;who changed this&quot;.
            </>,
          ]}
        />
        <p>
          Build that once, carefully, and it is a few weeks. Build it per resource, and it is a few
          weeks per resource, which is why admin panels end up inconsistent: the fifth table is
          always worse than the first.
        </p>
      </Section>

      <Section id="options" title="The four ways people do this">
        <H3>Hand-written React, per screen</H3>
        <p>
          Total control, and the cost is linear in the number of resources with no ceiling. It is
          the right answer for the two or three screens that are genuinely bespoke: an operations
          view stitching four queries together is a real page. It is the wrong answer for a list of
          invoices.
        </p>
        <H3>A hosted admin tool</H3>
        <p>
          Retool, Forest and similar: fast to a first screen, and your internal tooling now lives
          in somebody else&apos;s product, usually priced per seat, often requiring a connection
          into your database. Fine for a back office, awkward when the admin is part of the thing
          you sell.
        </p>
        <H3>A database browser</H3>
        <p>
          pgAdmin or a GUI client. Useful for you, dangerous for anybody else: there are no
          business rules, no validation and no audit, and &quot;just edit the row&quot; is how
          invariants get broken in ways the application cannot recover from.
        </p>
        <H3>Generated screens from a definition</H3>
        <p>
          Describe the resource, and the table, form, filters and detail view are derived from that
          description. This is the model Django admin and Laravel Filament popularised, and the one
          worth understanding, because it changes what a change costs: adding a column becomes a
          line rather than an edit in four files.
        </p>
      </Section>

      <Section id="definition" title="Describing a screen instead of writing one">
        <p>
          The whole idea in one file. This is a Grit resource definition, but the shape is the same
          in any system of this kind:
        </p>
        <CodeBlock
          language="ts"
          filename="apps/admin/resources/invoice.ts"
          code={`export const invoiceResource = defineResource({
  name: 'invoice',
  label: 'Invoices',
  icon: 'receipt',

  // The table. Chosen for scanning: six columns, not the twenty the model has.
  columns: [
    { name: 'number',    label: 'Invoice' },
    { name: 'customer',  label: 'Customer', relation: 'customer.name' },
    { name: 'total',     label: 'Total',    type: 'money', align: 'right' },
    { name: 'status',    label: 'Status',   type: 'badge', filter: true },
    { name: 'due_at',    label: 'Due',      type: 'date',  filter: 'date-range' },
  ],

  // The form. Validation comes from the Go model's binding tags, so there is
  // no second copy of the rules to drift.
  form: [
    { name: 'customer_id', label: 'Customer', type: 'relation', searchable: true },
    { name: 'number',      label: 'Invoice number' },
    { name: 'total',       label: 'Total', type: 'money' },
    { name: 'due_at',      label: 'Due date', type: 'date' },
    { name: 'notes',       label: 'Notes', type: 'textarea',
      // A conditional belongs here rather than in a useEffect that hides a
      // field: this declaration is read by the form AND by the validation,
      // so the API cannot accept what the form hid.
      showIf: { field: 'status', equals: 'disputed' } },
  ],

  search: ['number', 'customer.name'],
  defaultSort: { field: 'due_at', direction: 'desc' },

  emptyState: {
    title: 'No invoices yet',
    body: 'Invoices appear here once an order is marked fulfilled.',
  },
})`}
        />
        <p>
          Add <code>{'{ name: \'paid_at\', label: \'Paid\', type: \'date\' }'}</code> to{' '}
          <code>columns</code> and the column appears in the table, sortable, with the date
          formatted in the viewer&apos;s timezone. No component, no query change, no migration of
          four files.
        </p>
        <DocsShot
          src="/images/admin/table.png"
          alt="An admin list screen with sortable column headers, a search field, per-column filters, row selection checkboxes and pagination"
          label="localhost:3001/contacts"
          caption="A generated list screen in a real project. Sorting, filtering, selection and pagination come from the definition rather than from a page somebody wrote."
        />
        <Note title="The declaration is the thing to change">
          <p>
            Once screens are generated, the strongest move is almost always to change the
            declaration, and the second strongest is to change the model. Writing a custom page is
            a last resort, because you have just taken on the maintenance of something that used to
            update itself.
          </p>
        </Note>
      </Section>

      <Section id="decisions" title="The four decisions that matter">
        <p>
          Generation removes the typing. It does not remove the design, and these are where an
          admin panel is actually won or lost.
        </p>
        <H3>Which columns</H3>
        <p>
          The default is every field, and the default is almost always wrong. A table is something
          people scan: the columns are the ones they scan <em>for</em>. Everything else belongs on
          the detail view. A table with fourteen columns is a table nobody reads, and it will be
          horizontally scrolling on every laptop.
        </p>
        <H3>What the label column is</H3>
        <p>
          A relationship renders the related record&apos;s label. If a model has no human-readable
          label, the dropdown shows ids and the table shows nothing useful. The fix is on the
          related model, not on this screen, and getting it right once fixes every screen that
          references it.
        </p>
        <H3>What the empty state says</H3>
        <p>
          A new user sees the empty state first, and &quot;No data&quot; tells them nothing. Say
          what would cause a row to appear here. It is one sentence and it is the highest-value
          sentence on the screen.
        </p>
        <H3>What the dashboard asks</H3>
        <p>
          A widget earns its place by answering a question somebody actually asks at the start of
          the day. &quot;Total contacts&quot; is a number. &quot;Contacts with no activity in
          thirty days&quot; is a question, and somebody will act on it. Most generated dashboards
          are full of the first kind.
        </p>
        <Warn title="Charts are where a generated dashboard looks generated">
          <p>
            Six arbitrary colours in a chart is the clearest single tell. Use the product&apos;s
            accent and its semantic colours, success, warning, danger, and nothing else. A chart
            that uses the same palette as the rest of the screen reads as designed even when it was
            not.
          </p>
        </Warn>
      </Section>

      <Section id="security" title="An admin panel is the most dangerous surface you have">
        <p>
          It reads every table and writes most of them, which means the ordinary rules apply twice
          over.
        </p>
        <Checklist
          items={[
            <>
              <strong>Permissions per action, not per panel.</strong> &quot;Can see the admin&quot;
              is not a permission model. Delete, refund and export are each worth their own.
            </>,
            <>
              <strong>Every way into a table needs the same scope.</strong> The list, the bulk
              edit, the CSV export and the import are four separate paths to the same rows. A
              scoped list next to an unscoped export is an unscoped system.
            </>,
            <>
              <strong>Audit the writes, and the sensitive reads.</strong> Who looked at this
              customer&apos;s record is a question that gets asked, usually at the worst time.
            </>,
            <>
              <strong>Confirm the destructive actions, with the name typed out.</strong> A modal
              that only needs a click is a modal people click.
            </>,
            <>
              <strong>Soft delete, with a trash screen.</strong> Someone will delete the wrong row.
              Plan for the undo rather than for the apology.
            </>,
          ]}
        />
      </Section>

      <Section id="with-grit" title="Doing it with Grit">
        <p>
          Generating a resource writes the Go model, service and handler, the TypeScript types and
          Zod schema, a React Query hook, and the admin definition and page, all from one command:
        </p>
        <CodeBlock
          terminal
          code={`grit generate resource Invoice \\
  --fields "number:string:required,customer:belongs_to:Customer,total:money,status:enum(draft|sent|paid),due_at:date"`}
        />
        <p>
          The admin that comes with a new project already has the parts nobody enjoys building:
          sortable and filterable tables, a form builder covering thirty field types, searchable
          relationship pickers, bulk create and bulk edit grids, soft delete with a trash screen,
          an audit log, saved views, and system pages for jobs, files, cron and mail.
        </p>
        <DocsShot
          src="/images/forms/relationship.png"
          alt="A form field for selecting a related record, showing a searchable dropdown listing records by name rather than by id"
          label="localhost:3001/invoices/new"
          caption="A relationship field. Searchable, and showing the related record's label, which is the difference between a usable control and a list of numbers."
        />
        <DocsShot
          src="/images/system/audit.png"
          alt="An audit log screen listing changes with the actor, the action, the record and the time"
          label="localhost:3001/system/audit"
          caption="Who changed what, which is the first question an admin panel is asked once other people use it."
        />
        <p>
          It is generated into your repository as ordinary React, so a screen that needs to be
          different can be, and the ones that do not stay generated.{' '}
          <Link href="/docs/admin/resources" className="text-accent hover:underline">
            The admin reference
          </Link>{' '}
          covers the definition format in full.
        </p>
      </Section>

      <Faq
        slug="build-admin-dashboard-go-react"
        items={[
          {
            q: 'Should the admin panel be a separate app or part of the main one?',
            a: (
              <>
                Separate if the audiences differ, which they usually do: a customer-facing app and
                an operator console want different navigation, different density and different
                deploy cadence. A route group inside the main app is reasonable when the
                &quot;admins&quot; are just your team and there are a handful of screens. The one
                thing that is not reasonable is sharing a layout and hoping a role check hides the
                dangerous buttons.
              </>
            ),
            text:
              'Separate if the audiences differ, which they usually do: a customer app and an operator console want different navigation, density and deploy cadence. A route group inside the main app is fine when the admins are your own team and there are a few screens. Sharing one layout and relying on a role check to hide dangerous buttons is not.',
          },
          {
            q: 'Is a generated admin flexible enough for a real product?',
            a: (
              <>
                For the list and form screens, yes, and those are most of them. The honest limit is
                that a screen which is not really a resource, an operations view, a reconciliation
                tool, a wizard, is a hand-written page and should be. The value is that generating
                the eighty per cent leaves you time to design the twenty that matters.
              </>
            ),
            text:
              'For list and form screens, yes, and those are most of an admin panel. The limit is that a screen which is not really a resource, such as an operations view or a wizard, should be hand-written. Generating the routine eighty per cent leaves time to design the twenty that matters.',
          },
          {
            q: 'Can I use the admin components in my customer-facing app?',
            a: (
              <>
                You can copy them, and it is usually better than rebuilding a table from scratch.
                Two caveats. The admin is a dense operator console, so the copy is a starting point
                rather than the finished customer screen. And the two apps may use different design
                tokens, in which case a component moved across renders with no styling at all and
                reports no error, which is a confusing half hour.
              </>
            ),
            text:
              'You can copy them, which usually beats rebuilding a table. Two caveats: the admin is a dense operator console so the copy is a starting point rather than a finished customer screen, and if the two apps use different design tokens a moved component renders unstyled while reporting no error.',
          },
          {
            q: 'How do I stop the admin becoming the place business logic lives?',
            a: (
              <>
                Keep the rules in the API. If &quot;an invoice cannot be marked paid before it is
                sent&quot; is enforced in a React form, then the API, the mobile app and every
                script can break it. The admin should be one caller of rules that live in the
                service layer, not the place the rules are written.
              </>
            ),
            text:
              'Keep the rules in the API. A constraint enforced in a React form can be broken by the API, a mobile client or a script. The admin should be one caller of rules that live in the service layer rather than the place the rules are written.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/rbac-explained-golang',
            label: 'Roles and permissions explained',
            blurb: 'Who may open the panel, and who may press the dangerous buttons.',
          },
          {
            href: '/guides/first-go-react-app',
            label: 'Build your first full-stack app',
            blurb: 'The API the admin sits on top of, written out by hand.',
          },
          {
            href: '/docs/admin/resources',
            label: 'Resource definitions',
            blurb: 'Every option in the definition, and what each produces.',
          },
          {
            href: '/docs/admin/widgets',
            label: 'Dashboard widgets',
            blurb: 'Stats, charts and activity feeds, and choosing what to put on them.',
          },
        ]}
      />
    </GuidePage>
  )
}
