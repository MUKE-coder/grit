#import "../template.typ": *

= Services: where the rules live, and why not the handler

The thin-handler rule has come up three times now. This chapter is the one that
justifies it properly and shows what a generated service actually contains.

== The rule, and the reason

*Every database read and write for a resource lives in that resource's
service. The handler reads the request, calls one method, and writes the
answer.*

The reason is not tidiness. It is this:

```go
func (s *InvoiceService) Create(ctx context.Context,
    in CreateInvoice) (*models.Invoice, error)
```

That signature takes a `context.Context`, not a `*gin.Context`. The service has
no idea it is being called over HTTP. So the same method serves:

- the HTTP route
- a background job
- a CLI command in `cmd/`
- the CSV importer
- a test, with no server running

A rule enforced here is enforced in all five. A rule written in the handler
holds until the first time something other than a browser does the same thing,
and then it silently does not. That is the whole argument, and it is why "put
the validation in the handler, it is right there" is the most expensive
shortcut in a web codebase.

#bridge[
  In an Express app the equivalent discipline is a service module that never
  imports `req` or `res`. Most teams intend this and drift, because the
  handler is where the data already is. Go makes the drift visible: the moment
  a service needs `*gin.Context` it must import Gin, and that import is a
  review comment waiting to happen.
]

== What a generated service contains

Four things, in this order.

*The struct, with one dependency.*

```go
type ContactService struct {
    DB *gorm.DB
}
```

*The list configuration,* which is a security boundary:

```go
// What a client may search, sort and filter by. Whitelisted, because each
// name ends up in SQL.
var contactListConfig = paginate.Config{
    Searchable: []string{"name", "phone"},
    Sortable: map[string]bool{
        "id": true, "created_at": true, "name": true, "phone": true,
    },
    Filterable: map[string]bool{
        "id": true, "name": true, "phone": true, "group_id": true,
    },
}
```

*The writable set,* which is the other one:

```go
// Every column Patch and Bulk may write. id, the timestamps and the
// version are the framework's, and are dropped.
var writableContact = map[string]bool{
    "name": true, "phone": true, "group_id": true,
}
```

*The methods:* `List`, `GetByID`, `Create`, `Update`, `Patch`, `Delete`,
`Bulk`, `BulkCreate`, `BulkEdit`, `Export`.

== The two whitelists are the chapter

Chapter twenty two said a value can be a placeholder and a column name cannot.
These maps are the consequence.

A client sends `?sort=name&filter[group_id]=3`. Those strings end up in SQL as
*identifiers*, which no placeholder can protect. So each is checked against a
map the generator wrote:

```go
// Roughly what paginate does with it.
if !cfg.Sortable[requested] {
    requested = "created_at"      // fall back, never interpolate
}
```

`writableContact` is the same defence for writes, and it closes a different
hole. A `PATCH` that merged the client's JSON straight into the model would let
anybody set `id`, `created_at`, `version`, or on a user record, `role`. That is
mass assignment, and it is how a normal user becomes an administrator.

#gotcha[
  When you add a field by hand rather than with `grit generate field`, *you
  must add it to these maps yourself*. The symptom of forgetting is not an
  error: sorting by the new column silently falls back, and patching it
  silently does nothing.

  This is the main reason to use `grit generate field` rather than editing the
  struct: the command updates the model, the migration, the two maps, the Zod
  schema, the type and the admin form together.
]

== Ownership travels in the context

```go
func (s *ContactService) db(ctx context.Context) *gorm.DB {
    return s.DB.WithContext(ctx)
}
```

That looks trivial and is not. The context carries the organisation the
multitenant plugin resolved and the actor that per-record ownership is scoped
by, and the scoping is applied through the context rather than by each method
remembering to add a `WHERE`.

So a service method that uses `s.db(ctx)` is automatically scoped, and one that
reaches for `s.DB` directly is not.

#gotcha[
  This is the strongest argument for threading the request's context all the
  way down rather than creating a fresh one in the middle. A
  `context.Background()` halfway through a request does not merely lose
  cancellation, it loses the caller's identity, and the query comes back
  unscoped: every organisation's rows, not just this one's.

  When a background job genuinely must act as nobody, there is an explicit way
  to say so. The explicitness is the point: an unscoped query should be a
  sentence someone wrote on purpose.
]

== Errors, by design

The service returns the sentinels and typed errors from chapter eleven:

```go
func (s *ContactService) GetByID(ctx context.Context,
    id uint) (*models.Contact, error) {
    var c models.Contact
    err := s.db(ctx).First(&c, id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, fmt.Errorf("contact %d: %w", id, respond.ErrNotFound)
    }
    if err != nil {
        return nil, fmt.Errorf("getting contact %d: %w", id, err)
    }
    return &c, nil
}
```

Note what it does *not* do: it does not know about status codes, and it does
not log. It returns, and the boundary decides. That is chapter eleven's
"handle an error once" as a structural rule rather than advice.

== Adding your own method

This is the part people worry about. The generated service is a starting point,
and adding to it is ordinary Go:

```go
// InvoiceService methods you add live alongside the generated ones.
func (s *InvoiceService) MarkPaid(ctx context.Context,
    id uint, paidAt time.Time) (*models.Invoice, error) {

    inv, err := s.GetByID(ctx, id)      // reuse the generated method
    if err != nil {
        return nil, err
    }
    if inv.PaidAt != nil {
        return nil, fmt.Errorf("invoice %d: %w", id, ErrAlreadyPaid)
    }

    inv.PaidAt = &paidAt
    if err := s.db(ctx).Save(inv).Error; err != nil {
        return nil, fmt.Errorf("marking invoice %d paid: %w", id, err)
    }

    events.Emit(ctx, events.Event{
        Name: "invoices.paid", Resource: "invoices",
        ID: inv.ID, After: inv,
    })
    return inv, nil
}
```

Three things worth copying from that. It reuses `GetByID` rather than writing
another query. It returns a sentinel for the business rule, so the handler can
map it to 409. And it emits an event rather than calling the audit log, the
webhook sender and the realtime pusher by hand.

#gotcha[
  Put it in the *same file* as the generated methods, not a separate one.

  The manifest from chapter twenty nine hashes generated files, so an edited
  file is reported on a re-run rather than silently overwritten. That
  protection only applies to the file it knows about. A method in
  `invoice_custom.go` is invisible to that machinery, and the next person to
  regenerate will not be warned about anything.
]

== Where a rule does not belong in a service

Two exceptions worth naming.

*A rule about the shape of a request* belongs in the validation layer: "email
must look like an email" is a Zod rule and a binding tag, not a service method.

*A rule that spans resources* may belong in a saga rather than in either
service, because neither owns the sequence. Marking an invoice paid, adjusting
stock and sending a receipt is three services and a failure mode in the middle,
which is what `internal/saga` exists for.

#project(title: "Project 29: add a business rule properly")[
  Take the bookshop from chapter twenty seven and add a real rule: a book
  cannot be discounted below its cost price.

  ```bash
  grit generate field Book cost:money
  grit generate field Book discount_percent:percent
  grit migrate
  ```

  Then add the method to `services/book.go`, below the generated ones:

  ```go
  var ErrBelowCost = errors.New("discount would price below cost")

  func (s *BookService) ApplyDiscount(ctx context.Context,
      id uint, percent int) (*models.Book, error) {

      if percent < 0 || percent > 90 {
          return nil, fmt.Errorf("discount %d: %w", percent,
              respond.ErrValidation)
      }

      b, err := s.GetByID(ctx, id)
      if err != nil {
          return nil, err
      }

      final := b.Price.MinusPercent(percent)
      if final.Cents() < b.Cost.Cents() {
          return nil, fmt.Errorf("book %d: %w", id, ErrBelowCost)
      }

      b.DiscountPercent = percent
      if err := s.db(ctx).Save(b).Error; err != nil {
          return nil, fmt.Errorf("discounting book %d: %w", id, err)
      }
      return b, nil
  }
  ```

  Now prove the rule holds in more than one place, which is the point of the
  chapter. Write a test that calls it directly, with no HTTP:

  ```go
  func TestApplyDiscountRefusesBelowCost(t *testing.T) {
      s := newTestBookService(t)
      b, err := s.Create(context.Background(), CreateBook{
          Title: "Go", Price: money.FromCents(2000),
          Cost: money.FromCents(1500),
      })
      if err != nil {
          t.Fatal(err)
      }

      _, err = s.ApplyDiscount(context.Background(), b.ID, 50)
      if err == nil {
          t.Fatal("want an error for a discount below cost")
      } else if !errors.Is(err, ErrBelowCost) {
          t.Errorf("got %v, want ErrBelowCost", err)
      }

      _, err = s.ApplyDiscount(context.Background(), b.ID, 10)
      if err != nil {
          t.Errorf("a 10 percent discount should be allowed: %v", err)
      }
  }
  ```

  That test runs with `context.Background()` and no server, which only works
  because the rule lives in the service. Move the check into a handler and try
  to write the same test: you now need a request, a recorder and a router, and
  the CSV importer bypasses your rule entirely.

  Finally, write a four-line `cmd/discount/main.go` that calls the same method
  from the command line, and watch the rule apply there too, for free.
]

#tryit[
  Add the discount field to the admin by hand instead of with
  `grit generate field`: edit only the model. Then try to sort the Books table
  by it and patch it through the API. Both fail silently, because the two maps
  in the service do not mention it. Fix them, and you will not forget what
  those whitelists are for again.
]

#assignment(name: "Services - Assignment One")[
  Put a rule in the right place and prove it holds from two directions.

  *One.* Add a method to a generated service that enforces something real:
  stock cannot go below zero, an invoice cannot be paid twice, a booking cannot
  overlap. Return a typed error, not a string.

  *Two.* Call it from a handler. Then call the same method from a second
  caller: a CLI command, a seeder or a job. Both must be refused the same way,
  and that is the whole argument of the chapter.

  *Three.* Add the new sortable and filterable columns to both whitelists, then
  deliberately leave one out and watch the admin column quietly do nothing.

  *Four.* A test for the rule that does not go through HTTP at all.]
