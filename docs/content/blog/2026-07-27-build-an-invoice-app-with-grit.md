---
title: "Build an invoice app with Grit in ten minutes"
subtitle: "A hands-on build that uses every one of Grit's new pieces in context — a Cloudflare-style theme, belongs_to and line-item relations, dropdown and toggle fields that carry their own options, atomic invoice numbering, the column you forgot added in place, and a Print button that was already there. One resource at a time, one command each, with a link to the docs for every step."
series: "The Daily Grit"
edition: 12
date: 2026-07-27
readingTime: "12 min"
author: "Muke JohnBaptist"
tags: [grit, tutorial, invoices, codegen, relations, forms, printing, go, react]
canonical: "https://gritframework.dev/blog/build-an-invoice-app-with-grit"
---

Feature lists are easy to skim and hard to remember. So instead of listing what's
new in Grit, let's *build* something with it — a small but real invoice app — and
pick up every new piece along the way. By the end you'll have customers, invoices
with line items, auto-generated invoice numbers, a status dropdown, and a
printable invoice, and you'll have touched relations, the new field types,
`grit g field`, `grit generate sequence`, and the print view without any of them
feeling like a detour.

Ten minutes. One command per step. Each step links to the docs that explain it in
full.

## Step 1 — Scaffold, with a theme that fits

An invoicing tool is a business app, so let's give it a business look. Grit ships
three themes; `pulse` is the Cloudflare-inspired one — confident blue CTAs, cool
grey-blue canvas, white elevated cards. Pick it right at scaffold time:

```bash
grit new invoicer --triple --next --theme pulse
cd invoicer
```

That's a Turborepo with a Go API, a Next.js web app, and an admin panel, all
sharing one set of types and Zod schemas. (`aurora` is the Apple-style
black/white/grey theme if you'd rather; `atlas` is the default.)

📖 **Docs:** [Creating a project](/docs/getting-started/create-a-project) ·
[Architecture modes](/docs/concepts/architecture-modes)

## Step 2 — Customers, so invoices have someone to belong to

Every invoice needs a customer. Generate the resource first — it's the parent side
of our first relation:

```bash
grit g resource Customer --fields "name:string,email:string,company:string"
```

One command, and you have the Go model, service, and handler; the REST routes; the
Zod schema and TypeScript types in the shared package; the React Query hooks; and
a fully working admin page with a searchable, sortable table and a form. Nothing
to wire.

📖 **Docs:** [Code generation](/docs/concepts/code-generation) — what the eight
generated files are and how to customize them.

## Step 3 — The invoice: two relations and two choice fields

This is the heart of the app, so let's slow down. One `generate resource` call
models an invoice that **belongs to** a customer, **has many** line items, and
uses the new option-backed field types for its status and a flag:

```bash
grit g resource Invoice --fields \
  "number:string:optional,\
   status:select:draft=Draft|sent=Sent|paid=Paid,\
   sent:toggle,\
   customer:belongs_to:Customer" \
  --items "InvoiceItem:description:string,qty:int,unit_rate:float"
```

> **Why `number:string:optional`?** We're going to auto-generate the invoice
> number in [Step 4](#step-4--auto-number-the-invoices), so the user should never
> have to type it. String fields are *required* by default — the `:optional`
> modifier tells Grit not to demand it in the form, which is what lets the server
> fill it in. More on the timing in Step 4.

That's a lot on one line, so here's each piece, on its own.

### The two relations

**`customer:belongs_to:Customer`** is the *to-one* side. The invoice gets a
`customer_id` foreign key, and — the part you'd otherwise hand-build — the admin
form renders a **searchable customer picker** instead of a raw ID box, backed by a
live query against `/api/customers`.

**`--items "InvoiceItem:description:string,qty:int,unit_rate:float"`** is the
*to-many* side. It does three things in one flag: generates the whole
`InvoiceItem` resource, gives it a `belongs_to` back to the invoice, and renders
it as an inline, add/remove **line-items table right inside the invoice form**.
When you save the invoice, its rows are written in the *same database
transaction* — so an invoice and its items are always consistent.

📖 **Docs:** [Relationships in the admin](/docs/admin/relationships) ·
[Invoices &amp; line items](/docs/backend/invoices) (the full `--items` breakdown,
including how to do it as two separate commands).

### The two choice fields — and how they render

`status` and `sent` use Grit's new **option-backed field types**. There are four,
and the whole point is that you declare the *choices* in the field spec and Grit
generates the control, the validation, and the types to match. Here's each one:

| You write | Renders in the form as | Stored as | Types generated |
|-----------|------------------------|-----------|-----------------|
| `status:select:draft=Draft\|paid=Paid` | a **dropdown** (`<select>`) | Go `string` | Zod `z.enum([...])`, TS `"draft" \| "paid"` |
| `rating:radio:low=Low\|high=High` | a **radio-button group** (single choice) | Go `string` | Zod `z.enum([...])`, TS `"low" \| "high"` |
| `channels:check:email=Email\|sms=SMS` | a **checkbox group** (multi-select) | Go `datatypes.JSONSlice[string]` (JSON array) | Zod `z.array(z.enum([...]))`, TS `("email" \| "sms")[]` |
| `sent:toggle` | a **switch** | Go `bool` | Zod `z.boolean()`, TS `boolean` |

So in our invoice, **`status:select:…`** becomes a dropdown with Draft / Sent /
Paid, stored as a string and validated against exactly those three values
everywhere — the Go binding, the Zod schema, and the TypeScript union all agree
because they're generated from the same token. And **`sent:toggle`** is a simple
on/off switch backed by a boolean.

**`select` and `radio` are both single-choice** from a fixed list and generate the
exact same Go/Zod/TS — the only difference is the control. `select` is a compact
dropdown; `radio` lays every option out as a button. Reach for `radio` when there
are just a few choices you want visible at once (a plan tier, a priority), and
`select` when the list is longer or space is tight.

Want *multiple* choices instead of one? That's `check` — it renders a checkbox
group and stores the ticked values as a JSON array. For example, if invoices could
be delivered several ways: `channels:check:email=Email|sms=SMS|push=Push`.

### Labels are optional

The syntax is `type:value=Label|value=Label`, but **the `=Label` part is
optional**. Give just the values and Grit generates the labels for you by
capitalizing each one:

```bash
# These two are equivalent:
status:select:draft=Draft|sent=Sent|paid=Paid
status:select:draft|sent|paid          # labels auto-generated: Draft, Sent, Paid
```

Multi-word values are humanized, not just capitalized — `in_progress` becomes
**In Progress**, `awaiting_payment` becomes **Awaiting Payment**. You only reach
for `value=Label` when the label needs to differ from the stored value (say
`net_30=Net 30`, or a value that's an abbreviation you want spelled out). Mix and
match freely: `status:select:draft|sent=Sent to client|paid` works.

📖 **Docs:** [Field types](/docs/concepts/field-types) — the full table mapping
every type to its Go, GORM, TypeScript, Zod, and form representation.

### The same shape, four other apps

An invoice is just one instance of *a parent with a lifecycle status and a table
of child rows*. Once you see the shape, it's everywhere. Here are four other apps,
each one command, each leaning on the new field types the same way:

**E-commerce — orders and line items.** A `select` for the order lifecycle, a
`radio` for the two shipping options, line items as order lines.

```bash
grit g resource Order --fields \
  "number:string,\
   status:select:pending|paid|shipped|delivered,\
   fulfillment:radio:standard|express,\
   customer:belongs_to:Customer" \
  --items "OrderLine:product:string,qty:int,unit_price:float"
```

**Loan management — loans and repayments.** A `select` for the loan state, a
`radio` for the risk band, and each repayment carries a `date` and a `toggle`.

```bash
grit g resource Loan --fields \
  "reference:string,\
   status:select:pending|active|closed|defaulted,\
   risk:radio:low|medium|high,\
   borrower:belongs_to:Customer" \
  --items "Repayment:due_date:date,amount:float,paid:toggle"
```

**Subscriptions — plans and line charges.** Two `select`s (the plan and the
billing state) and an `auto_renew` `toggle`; line charges are the children.

```bash
grit g resource Subscription --fields \
  "plan:select:starter|pro|enterprise,\
   status:select:trialing|active|past_due|canceled,\
   auto_renew:toggle,\
   customer:belongs_to:Customer" \
  --items "LineCharge:description:string,qty:int,unit_price:float"
```

**Clinic — prescriptions and medications.** A `select` for the prescription
status, a `radio` for priority, and the prescribed drugs as line items.

```bash
grit g resource Prescription --fields \
  "code:string,\
   status:select:draft|issued|dispensed,\
   priority:radio:routine|urgent,\
   patient:belongs_to:Customer" \
  --items "Medication:name:string,dose:string,quantity:int"
```

Notice none of these spell out a single label — the bare `select`/`radio` values
get capitalized automatically (`past_due` → **Past Due**). Swap the nouns, keep
the moves.

## Step 4 — Auto-number the invoices

Nobody should type `INV-0001` by hand, and two people creating invoices at the
same moment must never collide. This is three small moves.

**4a. Generate the counter.**

```bash
grit generate sequence Invoice --prefix INV --reset monthly --width 4
```

**One thing to be clear about up front:** this command **does not create or modify
any column** on the Invoice model. It doesn't know about your `number:string`
field, and it doesn't need to. It generates a **standalone counter** — its own
tiny storage plus a function — and the word `Invoice` here is just a *name* for
that counter and its helper, not a link to the model. (You could name it
`grit generate sequence Receipt` and get `NextReceiptNumber`.) That's exactly why
step 4c is *you* deciding which field receives the value.

**4b. See what it wrote.** Two things:

- `internal/sequence/` — a generic counter package. The count lives in a database
  row that's locked and incremented in a transaction, which is what makes it
  atomic and gap-free.
- `internal/services/invoice_sequence.go` — a typed helper exposing one function,
  `NextInvoiceNumber(db, t)`, that returns the next formatted string.

The flags shape the output: `--prefix INV` is the leading text, `--width 4` is the
zero-padding, and `--reset` decides when the counter rolls back to 1. **Don't want
a monthly reset?** Use `--reset yearly` (rolls over each year) or `--reset never`
for one continuous, ever-incrementing series — `INV-0001`, `INV-0002`, … with no
date segment at all.

**4c. Wire it into the model.** Call the helper from the invoice's `BeforeCreate`
hook, so a number is assigned automatically — and only when one isn't already set,
so an imported invoice keeps its original:

```go
// apps/api/internal/models/invoice.go
func (m *Invoice) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.Number == "" {
		number, err := services.NextInvoiceNumber(tx, time.Now())
		if err != nil {
			return err
		}
		m.Number = number
	}
	return nil
}
```

Now every new invoice is numbered `INV-202607-0001`, `INV-202607-0002`, … The
counter is incremented in the *same transaction* as the insert, so it's safe under
concurrent load — no duplicates, no gaps.

**So what does the user see?** This is the part worth being precise about, because
it surprises people. The number is filled **on the server, at create time, only
when it's blank** — *not* in the browser as you type. So:

- Because we declared `number:string:**optional**` back in Step 3, the create form
  doesn't demand a number — you can submit with it empty.
- On submit, `BeforeCreate` runs and assigns the next number.
- The number then shows up on the invoice's **detail page and in the list**,
  immediately after you save. It was never typed by hand.

If you'd rather the empty Number box not appear in the create form at all, delete
the `number` field from the invoice's `form.fields` in
`apps/admin/resources/invoices.ts` — it stays in the table and on the detail page,
it just leaves the create form. (And if you want a shape like `2026/Q3/0001`, the
helper is plain Go you can edit; the sequence package just hands you the next
integer.)

📖 **Docs:** [Invoices &amp; line items → Auto-numbering](/docs/backend/invoices)

## Step 5 — The column you forgot

You build for ten minutes and realize invoices need a due date. In most
generators that means regenerating (and clobbering the `BeforeCreate` hook you
just wrote) or hand-editing five files. In Grit it's one command that adds the
column *in place*:

```bash
grit g field Invoice due_date:date
```

It injects the field into the Go model, the create **and** update Zod schemas, the
TypeScript type, and the admin form and table — at structural anchors, so your
`BeforeCreate` edit is untouched. There's no migration file to manage; the model
is the source of truth, so the database column appears on the next migrate:

```text
~ *models.Invoice — added 1 column(s): due_date
Migration done — 0 table(s) created, 1 altered (+1 column(s))
```

The same command adds a dropdown just as easily, labels-optional and all:
`grit g field Invoice terms:select:net_15=Net 15|net_30=Net 30`.

📖 **Docs:** [grit g field](/docs/concepts/cli) ·
[Migrations](/docs/backend/migrations)

## Step 6 — Run it, and print an invoice

Bring the schema up and start everything:

```bash
grit migrate
grit dev
```

Open the admin panel, add a customer, then create an invoice: pick the customer
from the searchable dropdown, choose a **status** from the select, flip the
**sent** switch, add a couple of line items in the inline table, and save. You
never touch the number — it fills itself in on save (Step 4).

### Where does the line-item total come from?

As you type quantities and rates, the inline table shows a per-row **Total** and a
grand total at the bottom — and you didn't ask for either. That's a small
convenience, so it's worth knowing exactly how it works, because it's driven by
**column names**, not magic:

- The line-items table looks for one column whose name matches a **quantity**
  pattern (`qty` or `quantity`) and one that matches a **money** pattern
  (`unit_rate`, `unit_price`, `rate`, `price`, or `amount`).
- If it finds both, it shows `Total = quantity × money` per row, and sums them.
  Our `qty:int` and `unit_rate:float` match, so you get it for free.

Two things follow from that. First, the total is **display-only** — it's computed
in the browser to help data entry; only your declared columns (`description`,
`qty`, `unit_rate`) are actually submitted and stored. If you need a *persisted*
total, add an `amount` column and compute it in the item's `BeforeSave` hook — the
same pattern as the invoice number. Second, because the trigger is the column
*name*, if you rename `qty` to `count` or `unit_rate` to `cost`, the auto-total
quietly stops appearing — rename them back to a name in those patterns (or keep
`qty`/`unit_rate`) and it returns. Nothing breaks either way; you just gain or lose
the convenience column.

Now open that invoice's detail page and hit **Print**. Every generated resource
detail page ships with the button, and a print stylesheet does the rest: the
detail content lives in a `#print-area`, and the sidebar, navbar, Edit/Delete
controls, and unrelated tables are all marked `no-print`. What reaches the paper
is exactly the invoice and its line items — customer, status, dates, the itemized
table — with no per-resource template to write.

📖 **Docs:** [Migrations](/docs/backend/migrations) ·
[Invoices &amp; line items → Printing](/docs/backend/invoices)

## What you actually built

Ten minutes, six commands, and you have:

- **Customers and invoices** with a `belongs_to` relation and a searchable picker.
- **Line items** as an inline, atomically-saved table (`--items`).
- A **status dropdown** and a **switch**, with their options — labels optional —
  generated across Go, Zod, and TypeScript so they can't disagree.
- **Auto-numbered** invoices that are safe under concurrent load.
- A **due date** you added after the fact without regenerating anything.
- A **printable** invoice, for free.

And "Invoice" is just the example — the exact same moves build orders and
order-items, subscriptions and line charges, or purchase-orders and receipts.
Swap the nouns; the commands don't change.

```bash
grit update   # get the latest, then build your own
```
