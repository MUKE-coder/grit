---
title: "Build an invoice app with Grit in ten minutes"
subtitle: "A hands-on build that uses every one of Grit's new pieces in context — a Cloudflare-style theme, belongs_to and line-item relations, dropdown and toggle fields that carry their own options, atomic invoice numbering, the column you forgot added in place, and a Print button that was already there. One resource at a time, one command each."
series: "The Daily Grit"
edition: 12
date: 2026-07-27
readingTime: "11 min"
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

Ten minutes. One command per step.

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

## Step 3 — Invoices: relations two ways, and fields that carry their options

Here's where the new pieces earn their keep. One `generate resource` call models
an invoice that **belongs to** a customer, **has many** line items, and uses the
new option-backed field types for its status and a flag:

```bash
grit g resource Invoice --fields \
  "number:string,\
   status:select:draft=Draft|sent=Sent|paid=Paid,\
   sent:toggle,\
   customer:belongs_to:Customer" \
  --items "InvoiceItem:description:string,qty:int,unit_rate:float"
```

Read that as four field decisions and one relation:

- **`customer:belongs_to:Customer`** — the *to-one* relation. The invoice gets a
  `customer_id`, and the admin form gets a searchable **customer picker** instead
  of a raw ID box.
- **`--items "InvoiceItem:…"`** — the *to-many* relation. It generates the whole
  `InvoiceItem` resource, gives it a `belongs_to` back to the invoice, and renders
  it as an inline, add/remove **line-items table** right inside the invoice form.
  Rows are saved atomically with the invoice, in one transaction.
- **`status:select:draft=Draft|sent=Sent|paid=Paid`** — a **dropdown**. From this
  one token Grit generates a Go `string`, a Zod `z.enum(["draft","sent","paid"])`,
  a TypeScript union `"draft" | "sent" | "paid"`, and the `<select>` with those
  labels. The enum, the validation, the type, and the control are generated
  together, so they can't drift apart.
- **`sent:toggle`** — a boolean rendered as a **switch**.

The field grammar for choices is `type:value=Label|value=Label`. Need a
multi-select instead of a single one? Use `check` — `channels:check:email=Email|sms=SMS`
stores an array and renders a checkbox group.

## Step 4 — Auto-number the invoices

Nobody should type `INV-0001` by hand, and two people creating invoices at the
same moment must never collide. `grit generate sequence` gives you an atomic,
gap-free counter backed by a database row:

```bash
grit generate sequence Invoice --prefix INV --reset monthly --width 4
```

That writes a generic counter package plus a typed helper,
`NextInvoiceNumber(db, t)`. Wire it into the invoice's `BeforeCreate` hook so the
number is assigned automatically — and only when blank, so an imported invoice
keeps its original:

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

Now every new invoice is `INV-202607-0001`, `INV-202607-0002`, … The counter is
incremented in the same transaction as the insert, so it's safe under concurrent
load — no duplicates, no gaps. `--reset monthly` rolls it over each month;
`--width 4` sets the padding.

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

While we're here, the same command adds a dropdown just as easily —
`grit g field Invoice terms:select:net_15=Net 15|net_30=Net 30`.

## Step 6 — Run it, and print an invoice

Bring the schema up and start everything:

```bash
grit migrate
grit dev
```

Open the admin panel, add a customer, then create an invoice: pick the customer
from the searchable dropdown, choose a **status** from the select, add a couple of
line items in the inline table, and save. The number fills itself in.

Now open that invoice's detail page and hit **Print**. Every generated resource
detail page ships with the button, and a print stylesheet does the rest: the
detail content lives in a `#print-area`, and the sidebar, navbar, Edit/Delete
controls, and unrelated tables are all marked `no-print`. What reaches the paper
is exactly the invoice and its line items — customer, status, dates, the itemized
table — with no per-resource template to write.

## What you actually built

Ten minutes, six commands, and you have:

- **Customers and invoices** with a `belongs_to` relation and a searchable picker.
- **Line items** as an inline, atomically-saved table (`--items`).
- A **status dropdown** and a **switch**, with their options generated across Go,
  Zod, and TypeScript so they can't disagree.
- **Auto-numbered** invoices that are safe under concurrent load.
- A **due date** you added after the fact without regenerating anything.
- A **printable** invoice, for free.

And "Invoice" is just the example — the exact same moves build orders and
order-items, subscriptions and line charges, or purchase-orders and receipts.
Swap the nouns; the commands don't change.

```bash
grit update   # get the latest, then build your own
```

Full walkthrough, including the combined-vs-separate `--items` breakdown and the
numbering internals, is in the [Invoices &amp; Line Items](/docs/backend/invoices)
guide.
