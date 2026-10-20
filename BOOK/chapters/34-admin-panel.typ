#import "../template.typ": *

= The generated admin panel

D4 was "the same UI for every table, twenty tables, twenty slightly different
list pages". D11 was "the frontend is somebody else's job, and there is nobody
else". This chapter is the answer to both, and the part of the framework that
is easiest to undersell and easiest to oversell.

== What you get without asking

Generate a resource and the admin gains a working screen: a table with
pagination, sorting, search and filters, a form with the right input for every
field type, a detail view, bulk actions, Excel export and CSV import.

```bash
grit generate resource Invoice \
  --fields "number:string, total:money, due:date, client:belongs_to"
```

Reload the admin and Invoices is in the sidebar, listing rows.

The form is the part worth looking at. `total` is a money input with a currency
prefix and two decimals. `due` is a date picker rather than a text box.
`client` is a searchable dropdown showing client *names*. None of that was
configured; it follows from the field types, which is chapter thirty's argument
arriving where you can see it.

== How it fits together

```
apps/admin/
├── app/
│   └── _dashboard/
│       └── resources/
│           └── invoice/
│               └── index.tsx     written once, never regenerated
└── components/                   the shared table, form and inputs
```

The generator writes `index.tsx` *once*. After that it is yours, and a
re-generate will not touch it. Everything underneath, the table component, the
form builder, the inputs, is framework code shared by every resource.

#gotcha[
  That single file is the seam, and it is worth understanding before you fight
  it.

  Customising a column, adding a computed cell, changing how a row renders:
  all of it happens in `index.tsx`, and it survives regeneration. Changing how
  *every* table behaves happens in the shared components, and that is
  framework code which `grit upgrade` will want to update.

  The rule: per-resource changes in the resource file, cross-cutting changes in
  the shared components and expect to merge them on upgrade.
]

== Customising a resource

The file is a definition rather than a page, which is why it is short:

```tsx
export default defineResource({
  name: "invoice",
  label: "Invoices",
  columns: [
    { key: "number", label: "Number", sortable: true },
    { key: "client.name", label: "Client" },
    { key: "total", label: "Total", align: "right" },
    { key: "due", label: "Due", sortable: true },
  ],
})
```

To change a cell, give the column a renderer:

```tsx
{
  key: "due",
  label: "Due",
  sortable: true,
  render: (row) => {
    const overdue = new Date(row.due) < new Date() && !row.paid_at
    return (
      <span className={overdue ? "text-red-500 font-medium" : ""}>
        {formatDate(row.due)}
      </span>
    )
  },
}
```

This is ordinary React. Everything you know applies, which is the point: the
generated screen is a starting point written in the same language you would
have used.

#bridge[
  If you have used Django admin or Laravel Filament, the shape is familiar: a
  resource definition rather than a page. The difference from Django is that
  this is React you can read and change, not a template system with its own
  idioms, and the difference from a hand-built panel is that twenty resources
  cost twenty definitions rather than twenty pages.
]

== The list page, honestly

What the table does for you:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Pagination], [Server side, using the `meta` from the last chapter],
  [Sorting], [Only on columns the service whitelisted],
  [Search], [Across the `Searchable` columns, server side],
  [Filters], [Per column, from the `Filterable` map],
  [Selection], [With bulk delete, edit and export],
  [Import], [CSV, with per-row errors rather than all-or-nothing],
  [Export], [Excel, honouring the current filters],
)

Notice that four of those are driven by the service's configuration from
chapter thirty two. The frontend does not decide what is sortable; it asks, and
the server is the authority. That is why a column you added by hand without
updating the maps appears in the table and refuses to sort.

== Where it stops

This is the part worth being plain about, because an admin panel demo always
looks like it does more than it does.

*It is an operator interface, not your product.* It is excellent for the
screens staff use: managing records, fixing data, running an import. It is not
a customer-facing experience, and trying to make it one fights the grain.

*A genuinely custom screen is a page, not a resource.* A dashboard that joins
four models with a chart and a date range is not a table of one resource. Write
it as a normal Next.js page using the generated hooks, and put it in the
sidebar by hand.

*The shared components are opinionated.* They are dark-themed, dense, built for
operators. Reskinning them entirely is possible and is a real project, and you
will be merging your changes against upgrades afterwards.

#gotcha[
  The failure mode to avoid is bending the resource definition into something
  it is not. If you find yourself writing a `render` for every column, hiding
  the table and nesting a custom layout inside it, you have written a page
  inside a resource file.

  Write the page. The generated hooks and types are available to it, and you
  lose nothing except the table you were not using.
]

== The hooks underneath

Each resource generates typed React Query hooks:

```tsx
const { data, isLoading } = useInvoices({ page: 1, search: "ACME" })
const { mutate: create } = useCreateInvoice()
const { mutate: update } = useUpdateInvoice()
```

These are the same hooks the generated table uses, and they are available to
any page you write yourself. The types come from `packages/shared`, which came
from the same field definitions as the Go model, so the autocomplete in your
editor is describing the actual API.

That chain is D3 closed end to end: change a field, run one command, and the
Go struct, the validation, the TypeScript type and this hook all move together.

#project(title: "Project 31: customise a resource screen")[
  Take the bookshop and make the Books screen genuinely useful, which means
  touching exactly one file.

  ```bash
  grit generate field Book stock:int
  grit generate field Book cover:file
  grit migrate
  grit seed Book --count 300
  grit start
  ```

  Open `apps/admin/app/_dashboard/resources/book/index.tsx`.

  *One: a computed cell.* Add a column that shows stock as a badge: red at
  zero, amber under ten, plain otherwise.

  ```tsx
  {
    key: "stock",
    label: "Stock",
    sortable: true,
    render: (row) => {
      const tone = row.stock === 0 ? "bg-red-500/15 text-red-400"
        : row.stock < 10 ? "bg-amber-500/15 text-amber-400"
        : "text-muted-foreground"
      return <span className={`rounded px-2 py-0.5 ${tone}`}>
        {row.stock}</span>
    },
  }
  ```

  *Two: a thumbnail.* Render `cover` as an image rather than a filename.

  *Three: prove the whitelist is the authority.* Make `stock` sortable in the
  definition, then try to sort by it. If you added the field with
  `grit generate field` it works. Now add a column for a field you added to
  the model by hand, mark it sortable, and watch it do nothing, because
  `services/book.go` never heard of it.

  *Four: prove the seam.* Re-run
  `grit generate resource Book --fields "..."` with the original fields. Your
  `index.tsx` edits survive, because that file is written once. Check the model
  and service: those are the files the manifest reports on.

  *Five: write a real page.* Add `apps/admin/app/_dashboard/reports/page.tsx`
  that uses `useBooks()` to show total stock value by category, with a chart.
  Notice you are writing ordinary React against typed hooks, and that this was
  never going to be a resource definition.
]

#tryit[
  Count the lines in your customised `index.tsx`. Then estimate what the same
  screen would have cost hand-written: a table with server-side pagination,
  sorting, search, filters, selection, a modal form with eleven input types,
  validation wired to the API's error details, CSV import with per-row errors
  and Excel export.

  That difference is D4 and D11 priced. It is also the fairest way to judge
  whether this part of the framework is worth it for your project, because if
  your application has three tables the answer might genuinely be no.
]

#assignment(name: "Admin Panel - Assignment One")[
  Customise one screen, then write the page the screen could not be.

  *One.* In one resource's `index.tsx`, add a computed cell, reorder the
  columns, and render one field as something other than text.

  *Two.* Re-run the generator for that resource and confirm your edits
  survived. Note which files the manifest reported instead.

  *Three.* Write a real page under the dashboard that uses the generated hooks
  to show something no table could: a chart, a comparison, a daily summary.

  *Four.* In `NOTES.md`, price it honestly. Count the lines you wrote, list
  what the table does for free, and say whether this part of the framework
  earns its place in your application. If your app has three tables, "no" is a
  respectable answer.]
