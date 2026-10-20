#import "../template.typ": *

= Project: an invoicing application

An invoice is not a record. It is a document somebody can be held to, and that
changes what the software has to guarantee: the number must be unique and
gap-free, the document must not change after it is sent, the arithmetic must be
exact, and every change must leave a trace.

Chapter forty did the arithmetic. This chapter is the rest, and most of it is
about refusing to do things.

== What makes this harder than CRUD

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Numbering], [INV-202610-0003 must be unique, sequential and gap-free, under
    concurrent creates],
  [Immutability], [A sent invoice is evidence. You correct it with another
    document, not with an edit],
  [State], [Draft, sent, part paid, paid, void, and the transitions that are
    not allowed],
  [Money], [Totals, tax, partial payments, in a currency that is not always
    yours],
  [Trace], [Who changed what, when, kept in something nobody can quietly
    rewrite],
)

== Numbering, which is a concurrency problem in disguise

The naive version is `SELECT MAX(number) + 1`, and it is wrong the first time
two people create an invoice in the same second. Both read the same maximum,
both write the same number, and one of them wins depending on your unique
index.

```bash
grit generate sequence Invoice --prefix INV \
  --reset monthly --width 4
```

That creates a counter package backed by a table, and a typed wrapper:

```go
number, err := services.NextInvoiceNumber(db, time.Now())
// INV-202610-0001
```

The table is one row per sequence and period:

```go
type Counter struct {
    Name      string `gorm:"primaryKey;size:50"`
    Bucket    string `gorm:"primaryKey;size:20"`
    NextValue uint64 `gorm:"not null;default:0"`
}
```

`Bucket` is the period key: `202610` for a monthly reset, `2026` for yearly,
empty for never. The increment runs inside a transaction holding a row-level
lock:

```go
tx.Clauses(clause.Locking{Strength: "UPDATE"}).
    Where("name = ? AND bucket = ?", cfg.Name, bucket).
    First(&c)
```

so concurrent callers serialise on that row rather than racing.

#gotcha[
  Call the sequence from the model's `BeforeCreate`, and call
  `sequence.Next` *directly* rather than the generated `services.Next...`
  wrapper.

  Services already import models. A model importing services is an import
  cycle, and Go refuses to compile it. The error names two packages and not the
  line you are thinking about, so it is worth knowing before you meet it:

  ```
  import cycle not allowed
  ```

  The wrapper exists for handlers and jobs. The hook, which lives in the model
  package, uses the generic function.
]

#bridge[
  The honest limit, stated in the package's own comment: one counter row
  serialises its writers, so beyond roughly a hundred invoices a second on a
  single counter you need a different design.

  That number is far above what an invoicing application does and far below
  what an event stream does. Knowing which side of it you are on is the useful
  part. Most business documents are rare events; telemetry is not.
]

== Immutability, and why an edit is the wrong verb

A draft is yours to change. A sent invoice is not, because a copy exists
outside your database, in somebody's inbox and somebody else's accounts.

So the model has two regimes:

```go
func (s *InvoiceService) Update(ctx context.Context,
    id uint, in UpdateInvoice) (*models.Invoice, error) {

    inv, err := s.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if inv.Status != "draft" {
        return nil, ErrNotEditable
    }
    ...
}
```

and a correction is a second document:

```go
func (s *InvoiceService) CreditNote(ctx context.Context,
    id uint, reason string) (*models.CreditNote, error)
```

#gotcha[
  The pressure to allow "just a small edit" on a sent invoice will come from
  the client, and it will sound reasonable, because the mistake is small.

  It is still wrong, and the argument that works is not technical: the number
  on the document in their customer's hand no longer matches the number in the
  system, and the first person to notice is an auditor. Offer a credit note and
  a reissue, which is what accountants already expect.
]

The ledger of what actually happened goes in a resource nobody can rewrite:

```bash
grit generate resource LedgerEntry --append-only \
  --fields "invoice:belongs_to:Invoice, kind:string, \
amount:money, note:text"
```

`--append-only` means no update and no delete routes, a GORM guard, and a
database trigger installed by `grit migrate`. Three layers, because the one
that matters is the one that still holds when somebody writes a quick script
against the database directly.

== The state machine, written as methods

Statuses in a dropdown are a suggestion. Statuses as methods are a rule:

```go
func (s *InvoiceService) Send(ctx context.Context, id uint) (
    *models.Invoice, error) {

    inv, err := s.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    switch inv.Status {
    case "sent", "paid":
        return inv, nil          // idempotent
    case "void":
        return nil, ErrVoided
    }
    if len(inv.Items) == 0 {
        return nil, ErrNoLines
    }

    now := time.Now().UTC()
    inv.Status, inv.SentAt = "sent", &now

    if err := s.db(ctx).Save(inv).Error; err != nil {
        return nil, fmt.Errorf("sending %d: %w", id, err)
    }
    events.Emit(ctx, events.InvoiceSent{ID: inv.ID})
    return inv, nil
}
```

The transitions worth refusing, in any document workflow:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Send an empty invoice], [A document with no lines is a mistake, not a zero
    bill],
  [Void a paid invoice], [Money moved. Refund it, do not pretend],
  [Pay a draft], [The customer has not been asked yet],
  [Edit after sending], [Covered above],
)

== Payments, and the status you should not store twice

A payment is its own record, because partial payments exist and because a
payment has a date and a method that the invoice does not:

```bash
grit generate resource Payment \
  --fields "invoice:belongs_to:Invoice, amount:money, \
paid_on:date, method:select:bank=Bank|card=Card|cash=Cash"
```

The balance is arithmetic over payments, not a column somebody updates:

```go
func (inv *Invoice) Balance() (money.Money, error) {
    total, err := inv.Total()
    if err != nil {
        return money.Money{}, err
    }
    for _, p := range inv.Payments {
        total, err = total.Sub(p.Amount)
        if err != nil {
            return money.Money{}, err
        }
    }
    return total, nil
}
```

#gotcha[
  Storing `paid_amount` on the invoice as well as storing payments is two
  sources of truth, and they will disagree. They always disagree after the
  first deleted payment, the first failed transaction, or the first manual
  correction in the database.

  Derive it. If the derivation is too slow at scale, that is a *caching*
  decision made later with measurements, and a cached value is clearly labelled
  as one. It is not the same as two columns that both claim to be the truth.
]

== The document itself

The admin gives you two ways out of the screen:

*Print and PDF.* The detail page has a print action, and PDF rendering uses
`@react-pdf/renderer` in the browser, lazy loaded so it is not in the first
load of every list page.

*Excel.* The list exports with the current filters applied, through `xlsx`,
also lazy loaded.

#bridge[
  Client-side PDF is right until it is not. It is free, it uses the layout you
  can see, and it needs no fonts installed on a server.

  It stops being right when the PDF must exist without a browser: attached to
  an automated email, archived nightly, generated by a webhook. At that point
  the job moves server side and becomes a real project with fonts, page breaks
  and a rendering budget.

  Start with the browser. Move it when a background job needs the document,
  not before.
]

== Reminders, in one line for now

Overdue invoices should chase themselves:

```bash
grit generate job SendInvoiceReminders --cron "0 9 * * *"
```

That writes a job that runs at nine every morning. Chapter forty three is about
what makes a job safe to run unattended: retries, idempotency, and what happens
when it fails at three in the morning. For now, note the shape: the job calls
the same service method a human would, so the rules cannot be bypassed by
automation.

#project(title: "Project 38: build the invoicing application")[
  ```bash
  grit new billing && cd billing
  grit generate resource Client \
    --fields "name:string, email:email, address:text"
  grit generate resource Invoice \
    --fields "number:string, client:belongs_to:Client, \
issued:date, due:date, \
status:select:draft=Draft|sent=Sent|paid=Paid|void=Void" \
    --items "InvoiceItem:description:string,qty:int,unit:money"
  grit generate sequence Invoice --prefix INV --reset monthly
  grit generate resource Payment \
    --fields "invoice:belongs_to:Invoice, amount:money, paid_on:date"
  grit generate resource LedgerEntry --append-only \
    --fields "invoice:belongs_to:Invoice, kind:string, amount:money"
  grit migrate
  ```

  *One: numbering under pressure.* Wire the sequence into
  `BeforeCreate`. Then create fifty invoices concurrently:

  ```go
  var wg sync.WaitGroup
  for i := 0; i < 50; i++ {
      wg.Add(1)
      go func() { defer wg.Done(); createInvoice(t) }()
  }
  wg.Wait()
  ```

  Assert fifty distinct numbers with no gaps. Run it with `-race`.

  *Two: the state machine.* `Send`, `Pay`, `Void` and `CreditNote` as service
  methods, each refusing the transitions in the table above, each idempotent
  where that makes sense. Tests for every refusal, not just the happy paths.

  *Three: immutability.* Prove that updating a sent invoice fails through the
  API and through the admin. Then try it in `/studio` and notice that you can,
  which is why the ledger exists.

  *Four: the balance.* Partial payments, with the balance derived. Add a
  payment in a different currency and confirm the error surfaces as a 422
  rather than a wrong number.

  *Five: the document.* Print an invoice to PDF from the detail page and read
  it as a customer would. Then export the list to Excel with a filter applied
  and confirm the filter was respected.

  *Six: the trace.* Write a ledger entry on every transition. Then try to
  delete one through the API, through the admin, and directly with SQL.
  Record which of the three refused you, and why that is the right number.
]

#tryit[
  Create an invoice in October and one in November and look at the numbers. The
  monthly reset means the counter starts again, which is what most accounting
  systems expect and which surprises people who assumed numbers are global.

  Then decide, for your own domain, whether a reset is right at all. Gap-free
  within a period and restarting is one convention; a single unbroken sequence
  forever is another. Both are defensible, and picking one after the first
  thousand documents is not.
]

#assignment(name: "Invoicing - Assignment One")[
  In `part5-building/41-invoicing`, build the application with its guarantees
  demonstrated rather than claimed.

  *One.* The resources, the sequence in `BeforeCreate`, and the concurrency
  test proving distinct gap-free numbers under fifty parallel creates.

  *Two.* The four state transitions as service methods, with a test for every
  refusal in the table.

  *Three.* `GUARANTEES.md` with one short section each: how numbering survives
  concurrency, what happens to a sent invoice somebody wants to change, how the
  balance is derived, and what the append-only ledger does and does not
  protect you from.

  *Four.* One deliberate failure. Remove the row lock from the sequence, run
  the concurrency test again, and paste the duplicate numbers into the notes.
  Then put it back.

  *Five.* A paragraph for a non-technical reader, the one you would send a
  client, explaining why you cannot simply edit a sent invoice. If you can
  write that paragraph, you understand the chapter.
]
