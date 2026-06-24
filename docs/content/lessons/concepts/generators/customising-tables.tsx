import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The table is the first thing operators see — the column choices,
        the formatting, the empty state. The generator gets it 80%
        right; this lesson covers the last 20%. By the end you can pick
        which columns show, format values as money/relative
        dates/badges, pack two fields into one column, hide noisy
        columns, and add server-side filters that drive a real CRM-style
        list.
      </p>

      <h2>Where the table lives</h2>
      <p>
        Same file as the form — <code>apps/admin/resources/&lt;plural&gt;.ts</code>{' '}
        — under the <code>table</code> key:
      </p>

      <CodeBlock
        language="ts"
        filename="apps/admin/resources/contacts.ts (excerpt)"
        code={`table: {
  columns: [
    { key: "name",       label: "Name",    sortable: true, searchable: true },
    { key: "email",      label: "Email",   sortable: true, searchable: true },
    { key: "phone",      label: "Phone" },
    { key: "group.name", label: "Group" },
    { key: "created_at", label: "Created", sortable: true, format: "relative" },
  ],
  filters: [],
  defaultSort: { key: "created_at", direction: "desc" },
  searchable: true,
  pageSize: 20,
},`}
      />

      <p>Every column accepts these keys:</p>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Key</th>
              <th className="text-left px-3 py-2 font-medium">What it does</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr><td className="px-3 py-2 font-mono text-[12px]">key</td><td className="text-[12px]">Field name. Supports dotted paths (<code>group.name</code>) for preloaded relations.</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">label</td><td className="text-[12px]">Column header.</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">sortable</td><td className="text-[12px]">Adds clickable sort arrows (server-side sort).</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">searchable</td><td className="text-[12px]">Includes this column in the global search box.</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">hidden</td><td className="text-[12px]">Defined but not shown — handy for columns you want filterable but invisible.</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">width</td><td className="text-[12px]">CSS width — keeps narrow columns narrow (<code>&quot;80px&quot;</code>, <code>&quot;15%&quot;</code>).</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">format</td><td className="text-[12px]">Pre-built renderer (see table below).</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">badge</td><td className="text-[12px]">For status-style columns — pairs with <code>format: &quot;badge&quot;</code>.</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">currencyPrefix</td><td className="text-[12px]">For <code>format: &quot;currency&quot;</code> — e.g. <code>&quot;$&quot;</code>, <code>&quot;€&quot;</code>.</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">className</td><td className="text-[12px]">Extra Tailwind classes on the cell.</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">cell</td><td className="text-[12px]"><strong>v3.31.15+</strong> — custom render function. Receives the full row. Overrides format/badge.</td></tr>
          </tbody>
        </table>
      </div>

      <h2>Built-in column formats</h2>

      <CodeBlock
        language="ts"
        code={`// 14 ready-made cell renderers
{ key: "name",         format: "text"     }   // default
{ key: "is_active",    format: "boolean"  }   // green check / grey dash
{ key: "status",       format: "badge",  badge: { active: { color: "success", label: "Active" }, … } }
{ key: "price",        format: "currency", currencyPrefix: "$" }
{ key: "starts_at",    format: "date"     }   // 2026-06-21
{ key: "created_at",   format: "relative" }   // "5 minutes ago"
{ key: "avatar",       format: "image"    }   // small thumbnail
{ key: "video_url",    format: "video"    }   // play icon link
{ key: "spec_sheet",   format: "file"     }   // FileRef preview (v3.31.30)
{ key: "photos",       format: "files"    }   // FileRef[] gallery stack
{ key: "website",      format: "link"     }   // external link with icon
{ key: "email",        format: "email"    }   // mailto: link
{ key: "color",        format: "color"    }   // colored swatch
{ key: "body",         format: "richtext" }   // stripped of HTML for preview
{ key: "user",         format: "user"     }   // packed avatar + name + email`}
      />

      <TipBox tone="info">
        The <code>ColumnFormat</code> TS union currently lists 13 of
        these — <code>file</code> and <code>files</code> render fine
        at runtime but you may need <code>as never</code> to satisfy
        the compiler if you set them by hand. They&apos;re emitted
        automatically by the generator for <code>:file:</code> and{' '}
        <code>:files:</code> columns.
      </TipBox>

      <h2>Recipe 1 — hide the noisy timestamp from the default list</h2>
      <p>
        The generator adds <code>created_at</code> automatically. For a
        contacts list you might prefer just last-seen or to drop it
        entirely:
      </p>
      <CodeBlock
        language="ts"
        code={`columns: [
  { key: "name",       label: "Name",  sortable: true, searchable: true },
  { key: "email",      label: "Email", sortable: true, searchable: true },
  { key: "phone",      label: "Phone" },
  { key: "group.name", label: "Group" },
  // dropped created_at
],`}
      />

      <h2>Recipe 2 — money + percent + relative date</h2>
      <CodeBlock
        language="ts"
        code={`columns: [
  { key: "title",       label: "Project", sortable: true, searchable: true },
  { key: "budget",      label: "Budget",  format: "currency", currencyPrefix: "$" },
  { key: "progress",    label: "Progress", format: "text",
    cell: (row) => <span>{Math.round(Number(row.progress) * 100)}%</span> },
  { key: "due_date",    label: "Due",     format: "relative", sortable: true },
],`}
      />

      <h2>Recipe 3 — status badge with colour coding</h2>
      <CodeBlock
        language="ts"
        code={`{
  key: "status",
  label: "Status",
  format: "badge",
  badge: {
    active:   { color: "success",  label: "Active"   },
    inactive: { color: "muted",    label: "Inactive" },
    pending:  { color: "warning",  label: "Pending"  },
    archived: { color: "danger",   label: "Archived" },
  },
}`}
      />

      <h2>Auto-packed columns (v3.31.19+)</h2>
      <p>
        Starting in v3.31.19, the generator detects two common pack
        patterns and emits a stacked column automatically. Generate a
        resource with both <code>name</code> and <code>email</code> and
        you&apos;ll get a single <strong>Contact</strong> column showing
        name on top and email below — no manual <code>cell:</code>{' '}
        callback needed.
      </p>

      <CodeBlock
        terminal
        code={`grit generate resource Customer \\
  --fields "name:string,email:string:unique,phone:string,company:string"`}
      />

      <p>The generated <code>resources/customers.ts</code> contains:</p>

      <CodeBlock
        language="ts"
        filename="apps/admin/resources/customers.ts (excerpt)"
        code={`import { defineResource } from "@/lib/resource";
import { StackedCell } from "@/components/tables/stacked-cell";

columns: [
  // grit:cols:auto-start
  // Packed automatically: name + email → single Contact column.
  { key: "name", label: "Contact", sortable: true, searchable: true,
    cell: (row) => StackedCell({ top: String(row.name ?? ""), bottom: String(row.email ?? "") }) },
  { key: "phone", label: "Phone", sortable: true, searchable: true },
  { key: "company", label: "Company", sortable: true, searchable: true },
  // grit:cols:auto-end
],`}
      />

      <h3>Patterns the generator recognises</h3>
      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Fields detected</th>
              <th className="text-left px-3 py-2 font-medium">Packed column</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">name + email</td>
              <td className="px-3 py-2 text-[12px]">"Contact" — name top, email muted below</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">first_name + last_name</td>
              <td className="px-3 py-2 text-[12px]">"Name" — concatenated full name</td>
            </tr>
          </tbody>
        </table>
      </div>

      <TipBox tone="info">
        The <code>StackedCell</code> helper is just a function — calling
        <code> StackedCell({`{ top, bottom }`})</code> returns the
        rendered JSX. That&apos;s why this works from a <code>.ts</code>{' '}
        file without forcing <code>.tsx</code>. Use the helper for any
        custom pack you write by hand too.
      </TipBox>

      <p>
        Resources scaffolded before v3.31.19 don&apos;t auto-pack —
        either add the pack manually following Recipe 4 below, or wait
        for the <code>grit pack table &lt;Resource&gt;</code> CLI
        command (roadmap).
      </p>

      <h2>Recipe 4 — pack multiple fields into one column (manual)</h2>
      <p>
        The big one. Sometimes &quot;Name&quot; really means
        &quot;Name + email stacked&quot;. Sometimes price + currency
        should be one cell, not two. As of <strong>v3.31.15</strong>,
        the <code>cell</code> property accepts a render function that
        receives the entire row — pack as many fields as you like:
      </p>

      <CodeBlock
        language="ts"
        filename="apps/admin/resources/contacts.ts"
        code={`columns: [
  {
    key: "name",
    label: "Contact",
    sortable: true,
    searchable: true,
    cell: (row) => (
      <div className="flex flex-col">
        <span className="font-medium text-foreground">{String(row.name)}</span>
        <span className="text-xs text-text-muted">{String(row.email)}</span>
      </div>
    ),
  },
  {
    key: "phone",
    label: "Contact info",
    cell: (row) => (
      <div className="flex flex-col">
        <span className="text-sm">{String(row.phone ?? "—")}</span>
        <span className="text-xs text-text-muted">{String((row.group as { name?: string })?.name ?? "no group")}</span>
      </div>
    ),
  },
  { key: "created_at", label: "Created", format: "relative", sortable: true },
],`}
      />
      <p>
        Five fields collapse into two columns. The <code>cell</code>{' '}
        function gets the full row (typed as{' '}
        <code>Record&lt;string, unknown&gt;</code>) so dotted keys
        aren&apos;t needed — read whatever properties make sense.
      </p>

      <TipBox tone="info">
        When <code>cell</code> is set, it <em>replaces</em> the entire
        cell rendering — <code>format</code>, <code>badge</code>, and{' '}
        <code>currencyPrefix</code> are ignored. The{' '}
        <code>key</code> still drives sorting + search, so put the most
        meaningful sort field there (usually <code>name</code> for the
        leftmost packed column).
      </TipBox>

      <h2>Recipe 5 — three-column dashboard look in a flat table</h2>
      <p>
        A pattern from CRMs: avatar + name + email in column 1, company
        + role in column 2, last activity in column 3. Three packed
        columns instead of seven thin ones:
      </p>
      <CodeBlock
        language="ts"
        code={`columns: [
  {
    key: "name",
    label: "Person",
    cell: (row) => (
      <div className="flex items-center gap-3">
        <img src={String(row.avatar ?? "/default-avatar.png")} className="h-9 w-9 rounded-full" />
        <div className="flex flex-col">
          <span className="font-medium">{String(row.name)}</span>
          <span className="text-xs text-text-muted">{String(row.email)}</span>
        </div>
      </div>
    ),
  },
  {
    key: "company",
    label: "Company",
    cell: (row) => (
      <div className="flex flex-col">
        <span>{String(row.company)}</span>
        <span className="text-xs text-text-muted">{String(row.role)}</span>
      </div>
    ),
  },
  {
    key: "last_activity_at",
    label: "Last activity",
    format: "relative",
    sortable: true,
  },
],`}
      />

      <h2>Filters — server-side query knobs</h2>
      <p>
        Filters add UI controls that translate to query-string params on{' '}
        <code>GET /api/&lt;plural&gt;</code>. The API&apos;s generated
        List handler already understands{' '}
        <code>?group_id=…</code>,{' '}
        <code>?created_at_from=…</code>, etc.
      </p>

      <CodeBlock
        language="ts"
        code={`filters: [
  {
    key: "status",
    label: "Status",
    type: "select",
    options: [
      { label: "Active",   value: "active"   },
      { label: "Inactive", value: "inactive" },
    ],
  },
  {
    key: "group_id",
    label: "Group",
    type: "select",
    options: [/* populate at runtime from a fetch */],
  },
  {
    key: "created_at",
    label: "Created date",
    type: "date-range",
  },
  {
    key: "is_active",
    label: "Active only",
    type: "boolean",
  },
],`}
      />

      <h2>Default sort, page size, search</h2>
      <CodeBlock
        language="ts"
        code={`table: {
  columns: [/* … */],
  defaultSort: { key: "created_at", direction: "desc" },
  pageSize: 20,           // 20 by default. 10/50/100 are nice round options.
  searchable: true,       // shows the global search box. Searches all searchable: true columns.
},`}
      />

      <h2>Hidden / always-defined-but-not-shown columns</h2>
      <p>
        Useful for fields you want filterable or queryable without
        cluttering the visible table:
      </p>
      <CodeBlock
        language="ts"
        code={`{ key: "internal_notes", label: "Internal notes", hidden: true, searchable: true }`}
      />

      <KnowledgeCheck
        question="You want a Contacts table that shows Name + email in column 1 (stacked), the group name in column 2, and the relative created-at in column 3 — just three columns. Best approach?"
        choices={[
          {
            label: 'Edit page.tsx by hand and replace ResourcePage with a custom table',
            feedback:
              "Works but heavyweight. The declarative resource definition is meant for this case — use the cell render function.",
          },
          {
            label: 'Use cell: (row) => <div>...</div> on the Name column and drop the email/phone columns from the array',
            correct: true,
            feedback:
              "Exactly. cell wins over format, you can stack JSX however you like, and the rest of the table machinery (sorting on the key, search) still works.",
          },
          {
            label: 'Set width: "0px" on email and phone',
            feedback:
              "Visually hides them but you're shipping invisible columns + their content with every page render. Drop them or use hidden: true.",
          },
          {
            label: "Sort by name and let the user understand which column is which",
            feedback:
              "Not the question — the user explicitly wants 3 columns showing 5 fields.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Customise your Contact table to look like a CRM:
            </p>
            <ol>
              <li>
                Drop the <code>created_at</code> column from the visible
                list.
              </li>
              <li>
                Replace the Name column with a stacked cell showing
                Name on top and email below (smaller, muted).
              </li>
              <li>
                Replace the Phone column with a stacked cell showing
                phone on top and group name below.
              </li>
              <li>
                Add a filter on <code>group_id</code> with options
                fetched from <code>/api/groups</code>.
              </li>
            </ol>
            <p>Reload the page and confirm the layout is tighter and the filter works.</p>
          </>
        }
        hint={
          <>
            For the dynamic filter options, you&apos;ll either need to
            pre-populate them (an effect at the top of the file) or
            use a <code>relationship-filter</code> wrapper. Easiest
            v1: hardcode 3–5 options to confirm the filter wiring,
            then make them dynamic.
          </>
        }
        solution={
          <>
            <CodeBlock
              language="ts"
              filename="apps/admin/resources/contacts.ts (table block)"
              code={`table: {
  columns: [
    {
      key: "name",
      label: "Contact",
      sortable: true,
      searchable: true,
      cell: (row) => (
        <div className="flex flex-col">
          <span className="font-medium text-foreground">{String(row.name)}</span>
          <span className="text-xs text-text-muted">{String(row.email)}</span>
        </div>
      ),
    },
    {
      key: "phone",
      label: "Phone / Group",
      cell: (row) => (
        <div className="flex flex-col">
          <span>{String(row.phone ?? "—")}</span>
          <span className="text-xs text-text-muted">
            {String((row.group as { name?: string })?.name ?? "no group")}
          </span>
        </div>
      ),
    },
  ],
  filters: [
    {
      key: "group_id",
      label: "Group",
      type: "select",
      options: [
        { label: "Clients",  value: "01HX…" },
        { label: "Leads",    value: "01HX…" },
        { label: "Vendors",  value: "01HX…" },
      ],
    },
  ],
  defaultSort: { key: "created_at", direction: "desc" },
  searchable: true,
  pageSize: 20,
},`}
            />
            <p>
              Two columns, five fields, one filter. The created_at sort
              still works because the <code>defaultSort.key</code>{' '}
              doesn&apos;t require the column to be visible.
            </p>
          </>
        }
      />

      <h2>Date-window filter (v3.31.34+)</h2>
      <p>
        Every list page ships with a toolbar Date filter — Today,
        Last 7 days, Last 30 days, This month, or a custom range. It
        rewrites <code>?created_since=</code> / <code>?created_from=</code> /{' '}
        <code>?created_to=</code> into the list query and feeds the
        same params to every stat card on the page, so the numbers
        above the table always match the rows below.
      </p>

      <CodeBlock
        language="ts"
        code={`table: {
  // Defaults to enabled on created_at.
  dateFilter: { enabled: true, field: "created_at", label: "Created" },
  // Hide entirely:
  // dateFilter: { enabled: false },
  // Filter on a different timestamp instead — Booking by scheduled_for,
  // BlogPost by published_at, Order by placed_at:
  // dateFilter: { field: "scheduled_for", label: "Scheduled" },
},`}
      />
      <p>
        The selection is URL-persisted via{' '}
        <code>?date=preset</code> / <code>?date_from</code> /{' '}
        <code>?date_to</code>, so a refresh or shared link rehydrates
        the same view. No code changes needed on the API — the{' '}
        <code>paginate</code> package already parses all four params.
      </p>

      <h2>Export + Import (v3.31.35+)</h2>
      <p>
        The toolbar also ships a split <strong>Export</strong> button
        (CSV / Excel / JSON) and an <strong>Import</strong> button.
        Both run client-side via SheetJS — no new API routes, no
        async cutoff.
      </p>

      <CodeBlock
        language="ts"
        code={`table: {
  // All three formats on by default; default click triggers Excel.
  // Hide a format from the menu:
  export: { csv: true, excel: true, json: false },
  // Disable export entirely:
  // export: false,

  // Excel import on by default. Restrict to a subset of form fields:
  import: { fields: ["title", "price", "stock"] },
  // Disable import entirely:
  // import: false,
},`}
      />
      <p>
        Export loops the paginated endpoint until every row is fetched
        (using the current filter + sort + date range), then writes
        the file. Import opens a modal that parses a dropped file,
        previews per-row validation errors, and POSTs each valid row
        at concurrency 4. File-field columns (<code>:file:</code> /{' '}
        <code>:files:</code>) are auto-excluded from imports —
        spreadsheets can&apos;t carry binary blobs.
      </p>

      <h2>What&apos;s next</h2>
      <p>
        Forms and tables are the operator side. The final lesson in this
        chapter takes the same generated resource and shows you how to
        consume it from the customer-facing web app — using the
        auto-generated React Query hook from{' '}
        <code>apps/web/hooks/use-&lt;plural&gt;.ts</code>.
      </p>
    </>
  )
}
