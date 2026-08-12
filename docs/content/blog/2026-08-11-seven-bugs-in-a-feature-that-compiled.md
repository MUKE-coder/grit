---
title: "Seven bugs in a feature that compiled"
subtitle: "Grit shipped a way to replace any admin table, form or page with your own. Tests passed, types checked, docs were accurate. Then we built a real app with it and found seven bugs — including one where the component was correct, the DOM was correct, and the screen was blank. A tour of the customisation system, told through everything that was wrong with it."
series: "The Daily Grit"
edition: 14
date: 2026-08-11
readingTime: "11 min"
author: "Muke JohnBaptist"
tags: [grit, admin, customisation, testing, typescript, tailwind, go, react]
canonical: "https://gritframework.dev/blog/seven-bugs-in-a-feature-that-compiled"
---

Three releases ago Grit got a customisation system. The idea: you bought a
dashboard template, you like its table, and you want it inside your admin without
giving up the sorting, paging, filters, selection, bulk delete and toasts that
the generated page already has. So a resource grew a sibling file:

```
apps/admin/resources/
  products.ts          # generated — rewritten on every grit generate
  products.custom.tsx  # yours — created once, never touched again
```

Put a component in the second file and it replaces the corresponding piece of the
first. Four slots — `Table`, `Form`, `EmptyState`, `Page` — plus per-column and
per-field patches. Everything else keeps working.

It compiled. The Go tests passed. `tsc --noEmit` was clean. The docs described it
accurately, because I wrote them from the source.

Then I built an app with it, and found seven bugs.

This post is a tour of the feature by way of its failures, because that turns out
to be the better tour. Every bug is a place where the design met reality and one
of them was wrong.

## The app

`portkit` — a small ops console, scaffolded fresh, three resources, each one
exercising a different part of the system:

```bash
grit new portkit --triple --next --theme pulse
grit g resource Product  --fields "name:string,sku:string,price:float,stock:int,status:select:active|draft|archived" --faker --count 24
grit g resource Deal     --fields "title:string,company:string,value:float,stage:select:lead|qualified|won|lost,owner:string" --faker --count 18
grit g resource Enquiry  --fields "subject:string,requester:string,priority:select:low|normal|high,resolved:bool,body:text"
```

Products keeps the stock table and patches three cells. Deals throws the table
away entirely and becomes a kanban board. Enquiries replaces the list, the form
and the empty state at once. Between them, every slot gets used the way somebody
porting a template would use it.

## Bug 1: the generator deleted the support desk

The first resource I tried to generate was called `Ticket`.

```bash
grit g resource Ticket --fields "subject:string,requester:string,..."
  ✅ Resource Ticket generated successfully!
```

Then:

```
internal/models/user.go:129:4: undefined: TicketReply
```

Grit ships a support desk. It has a `Ticket` model and a `TicketReply` model, in
`internal/models/ticket.go`. The generator wrote its own `ticket.go` over the top
of it. `TicketReply` went with it, and the build broke in `user.go` — a file I had
never touched, naming a symbol I had never heard of.

Note the shape of this failure. The command reported success. The error surfaced
somewhere else entirely, several minutes later, with nothing connecting the two.
And when I ran `grit remove resource Ticket` to undo it, that finished the job:
it deleted the model, the handler and the routes that the *scaffold* had written,
because as far as it knew they were mine.

Thirty-odd built-in model names are reserved now:

```
"Ticket" is a built-in model — it belongs to the support desk

Generating over it would overwrite apps/api/internal/models/ticket.go and
break the build, and `grit remove resource Ticket` would then delete the
original. Pick another name:

  grit generate resource SupportTicket --fields "..."

If you really mean to replace the built-in, pass --force.
```

The reserved list is hand-written, which means the interesting failure mode is
drift: someone adds a built-in model next year, nobody adds it to the list, and
`Ticket` happens again under a different name. So there is a test that scaffolds
a real project, reads every model that `AutoMigrate` is given, and fails if the
list has fallen behind. I deleted one entry to make sure it actually fails. It
does.

## Bug 2: `--faker` filled a dropdown with dictionary words

The seeded data looked like this:

```
   status   | count
------------+-------
 huh        |     1
 sufficient |     1
 moreover   |     1
 ouch       |     1
```

`status` is a `select` with three declared options. The faker seeder had no case
for choice fields, so it fell through to the string branch and called
`gofakeit.Word()`.

This is worse than untidy. Those values contradict the form's own dropdown, the
API's validation would reject them on write, and the TypeScript union that
`grit sync` generates from the Go struct says they cannot exist:

```ts
status: "active" | "draft" | "archived"
```

Choice fields now seed from their own options. One line in the generator, and
suddenly a freshly seeded app is internally consistent.

## Bug 3: the row type is a promise about the API, not a guarantee

Bug 2 had a passenger. My status cell looked like this:

```tsx
const STATUS = {
  active:   { label: "Active",   className: "bg-emerald-700 text-white" },
  draft:    { label: "Draft",    className: "bg-gray-600 text-white" },
  archived: { label: "Archived", className: "bg-amber-700 text-white" },
};

columns: {
  status: {
    cell: (row) => {
      const s = STATUS[row.status];
      return <span className={s.className}>{s.label}</span>;
    },
  },
}
```

`row` is a `Product`. `row.status` is that three-value union. TypeScript is
completely satisfied, and the page crashed on load:

```
TypeError: Cannot read properties of undefined (reading 'className')
```

Fixing the seeder made the crash go away, and that is exactly why it is worth
writing down. The type came from the Go struct. The *value* came from the
database — and a database picks up values from imports, migrations and
hand-written `UPDATE`s that no type ever saw. A cell renderer runs against
whatever actually arrived.

```tsx
const s = STATUS[row.status] ?? UNKNOWN;
```

One `??`, and an unexpected value renders a grey dash instead of taking down the
table in front of whoever opened the page.

## Bug 4: you could not wrap the thing you were replacing

The docs said a slot receives the stock component's own props, so you could put
the original inside yours:

```tsx
components: {
  Table: (props) => <TemplateCard><DataTable {...props} /></TemplateCard>,
}
```

I wrote a small file to check that claim rather than trust it:

```
error TS2322: Type '{ columns: ColumnDefinition<Product>[]; data: Product[]; ... }'
is not assignable to type 'DataTableProps'.
```

The customisation surface had been made generic over the row — that was the whole
of the previous release. The components it hands rows to had not. `DataTable`
took `Record<string, unknown>`, a typed overlay gives it `Product[]`, and a
`Product` has no index signature. Same story for `FormSheet`, which meant a custom
page could not pass `controller.form.item` to the stock form either.

So the two most-recommended patterns in the documentation — wrap the default,
reuse the default dialogs — did not compile. Both are now generic over the row,
with the erasure done once at each component's boundary instead of scattered
through its render.

The lesson I keep relearning: **a documented claim is a testable claim.** It took
fifteen lines to find out this one was false.

## Bug 5: `grit upgrade` had been updating files nothing imports

Having fixed the components in the scaffold, I ran `grit upgrade` on the test
project to pick them up.

```
✓ Admin panel updated (65 files)
```

Same type error. The components had not changed.

The admin's components were renamed to kebab-case a long time ago —
`data-table.tsx`, `form-sheet.tsx`. The upgrade command's path list still had the
old PascalCase names. So every upgrade for however long had been writing
`components/tables/DataTable.tsx` next to the real `data-table.tsx` and leaving it
there. Thirty-one files, none of them imported by anything.

The symptom was the opposite of an error. The command reported dozens of files
updated and exited zero. It just updated the wrong ones — which means **no
component fix shipped in an upgrade had reached anybody since the rename.** That
is the most consequential bug in this post, and I only found it because a fix I
had just written failed to appear.

The paths are correct now, the strays are cleaned up on the next upgrade, and
five components that were never in the list at all — including
`use-resource-controller`, the hook this entire feature is built on — are
refreshed too.

### A footnote from Windows

The cleanup deletes a stray only when the real file is present, so a half-finished
rename cannot take the last copy. One pair differs only in case:
`Providers.tsx` and `providers.tsx`. On Windows and macOS,
`os.Stat("providers.tsx")` cheerfully returns the entry for `Providers.tsx` — so
the guard passed, and the delete took the only copy. I watched the admin lose its
provider.

```go
// os.Stat is not good enough here. Reading the directory and comparing names
// byte for byte is the only answer that means the same thing on every platform.
func existsExact(path string) bool { ... }
```

## Bug 6: the admin never type-checked clean, and it was our fault

While counting errors I noticed three that had nothing to do with me:

```
components/language-switcher.tsx(5,27): Cannot find module 'next-intl'
components/language-switcher.tsx(16,8): Cannot find module '@/components/ui/dropdown-menu'
i18n/request.ts(1,34): Cannot find module 'next-intl/server'
```

The scaffold was writing six i18n files without the `next-intl` dependency that
compiles them. Nothing imported them. And the kicker: `grit add i18n` writes those
same six files properly — with the dependency, the provider and the plugin — but
it skips files that already exist. The broken copies were blocking the command
that would have fixed them.

Gone from the scaffold, pruned on upgrade when `next-intl` is absent. A fresh
Grit admin now reports **zero** errors from `tsc --noEmit`, which it never has
before.

## Bug 7: the one nothing but a browser could find

Everything compiled. Zero type errors. I opened the page.

![The status column, empty](/blog/portkit-blank-status.png)

The `STATUS` column is blank. The stock bars have no fill. Every other column is
perfect — the price is formatted, the numbers are there.

The DOM was correct:

```html
<span class="inline-flex rounded-full px-2 py-0.5 text-xs font-medium
             bg-emerald-700 text-white">Active</span>
```

The right element, the right classes, the right text. And nothing on screen,
because `bg-emerald-700` did not exist in the stylesheet. `text-white` did —
inherited from elsewhere in the app — so I was looking at white text on a
background that was never painted.

```ts
content: [
  "./app/**/*.{ts,tsx}",
  "./components/**/*.{ts,tsx}",
  "./lib/**/*.{ts,tsx}",
],
```

Tailwind was not scanning `resources/`. The one directory the entire feature
invites you to write markup in, and the compiler that turns those class names into
CSS had never been told it existed.

Add the glob, and the same page:

![The same page with the glob added](/blog/portkit-status-fixed.png)

I want to sit on this one, because it is the reason the whole exercise was worth
doing. There was no error. No warning. No failing test. The types were right, the
component was right, the markup was right, and the feature was unusable. A
compiler tells you your code is well-formed. It has never once told you that
anyone can see it.

## What the feature actually does, now that it works

With all seven fixed, here is the thing itself.

**Patch a cell and keep the table.** The cheapest customisation there is:

```tsx
const custom: ResourceCustomisation<Product> = {
  columns: {
    price:  { cell: (row) => <span className="font-mono">{money.format(row.price)}</span> },
    stock:  { cell: (row) => <StockBar value={row.stock} /> },
    status: { cell: (row) => <StatusPill value={row.status} /> },
  },
  fields: {
    sku: { placeholder: "PK-0000", description: "Printed on the picking slip." },
  },
};
```

`row` is a `Product`, so `row.price` is a number and `row.prise` is a compile
error. Columns and fields are patched **by key**, so `grit sync` can go on adding
new columns from the Go model without touching your renderers.

**Replace the page and keep the machinery.** Deals is not a table at all:

```tsx
function PipelineBoard({ resource }: ResourcePageSlotProps) {
  const c = useResourceController<Deal>(resource, { initialPageSize: 100 });
  const { mutate: patch } = usePatchResource(resource.endpoint, "Deal");
  // ...four columns of cards, drag-free stage buttons, totals per stage
}

const custom: ResourceCustomisation<Deal> = {
  components: { Page: PipelineBoard },
};
```

Rows, loading, search, the create/edit/delete actions and both confirm dialogs
still come from the controller. There is no second copy of the fetching logic —
the stock page is built on the same hook, which is what makes it safe to claim
your page can do anything the default one can.

One thing a `Page` slot does inherit is responsibility for its dialogs. The state
is still in the controller; you just have to render them:

```tsx
{c.form.open && <FormSheet resource={resource} item={c.form.item} onClose={c.form.close} />}
```

That line is the one bug 4 was blocking.

**Replace three slots at once.** Enquiries is an inbox, so it gets a custom list,
a custom composer and a custom empty state:

```tsx
components: {
  Table: InboxList,       // props are DataTable's, unchanged
  Form: Composer,         // saves via useCreateResource / useUpdateResource
  EmptyState: EmptyInbox, // shown when the query finishes with nothing
}
```

Everything around them is still the generated page: the search box, the date
filter, the exporter, the column picker, the pagination. Create a row in the
custom composer and the empty state gives way to the custom list, the stat cards
tick to 1, and the pager reads "Showing 1–1 of 1" — none of which is code anybody
wrote twice.

## And it survives regeneration

The point of the whole two-file split. After building all of the above, I
re-ran the generator over a resource that already had a customisation:

```bash
grit g resource Product --fields "name:string,sku:string,price:float,..."
  ✅ Resource Product generated successfully!
```

`products.ts` was rewritten from scratch — every column back to its generated
form. `products.custom.tsx` was not opened. The cells still render, because they
were never in the file that got replaced.

Deleting a resource is the one case where the overlay moves, and it took a bug to
notice: left in place, it imports a type the shared package no longer exports and
the admin stops type-checking. So `grit remove resource` deletes an untouched
stub, and renames one you have written in:

```
→ apps/admin/resources/gadgets.custom.tsx.bak (your customisations, kept out of the build)
```

## The through-line

Seven bugs. One found by a compiler, two by a test I wrote to check a sentence in
my own documentation, three by running a command and reading what it did, and one
— the one that made the feature useless — by looking at a screen.

The feature was not badly built. It was *unexercised*. Every one of these bugs
lived in the gap between "the code is correct" and "a person can use this," and
nothing closes that gap except being the person.

If you are shipping something this week: build the smallest real thing with it
before you announce it. Not a test. A thing you would be embarrassed to hand
someone.

📖 **Docs:** [Custom pages and tables](/docs/admin/custom-pages) ·
[Code generation](/docs/concepts/code-generation) ·
[Changelog v3.141.0](/docs/changelog)
