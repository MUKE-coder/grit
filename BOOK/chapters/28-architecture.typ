#import "../template.typ": *

= The architecture: a modular monolith, layer by layer

You have the project. This chapter is the map: how the code is arranged, which
layer is allowed to call which, and why those rules exist. Get this and
everything in part five is predictable.

== The short answer

*A modular monolith with a layered API and an in-process event bus.* One
binary, one database, about fifty internal packages.

That is the headline and it is slightly misleading, because no real system is
one style. A Grit project runs six at once, each where it earns its keep:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [*Modular monolith*], [The whole API: one deployment, one transaction boundary],
  [*Layered (n-tier)*], [routes, handlers, services, models],
  [*Event-driven*], [`internal/events`: a write emits once, subscribers react],
  [*Transactional outbox*], [`internal/outbox`: durable delivery of those events],
  [*Orchestration*], [`internal/saga`: multi-step work that can half fail],
  [*Microkernel*], [The framework itself: `grit plugin add`],
)

#bridge[
  Coming from Next.js, the mental shift is that there is one process rather
  than a function per route. State can live in memory, a connection pool is
  warm, a background worker runs, and nothing is cold-started. The cost is that
  you deploy and scale one thing, which for most products is a relief rather
  than a limitation.
]

== The request path

Follow one request all the way down. This is the shape every generated
resource has.

```
HTTP request
  -> routes       which URL maps to which handler, and what guards it
  -> middleware   auth, CORS, rate limit, recover, request id
  -> handler      read the request, call one service method,
                  write the answer
  -> service      the rules, and every database call
  -> model        the struct, its GORM tags, its hooks
  -> database
```

Each layer talks only to the one below it. A handler never touches the
database; a model never calls a service.

== The rule the layering exists for

Handlers are thin *on purpose*, and the reason is not tidiness.

```go
// A generated handler, in essence.
func (h *ContactHandler) List(c *gin.Context) {
    opts := paginate.FromQuery(c, contactListConfig)
    result, err := h.service().List(h.ctx(c), opts)
    if err != nil {
        respond.Error(c, err)
        return
    }
    respond.Paginated(c, result)
}
```

Read the request, call one method, write the answer. Every rule lives one layer
down, in the service.

The payoff: a rule enforced in the service is enforced *everywhere*. The
background job uses it. The CSV import uses it. The CLI command uses it. The
test uses it. A rule written in a handler holds until the first time something
other than a browser does the same thing, and then it silently does not.

```go
// Every service method takes a context, not a *gin.Context.
func (s *ContactService) List(ctx context.Context,
    opts paginate.Options) (*paginate.Result[models.Contact], error)
```

That signature is the whole architecture in one line. `context.Context` rather
than `*gin.Context` means the service has no idea it is being called over HTTP,
so a job can call it with `context.Background()`.

#gotcha[
  The context carries more than cancellation here. It holds the organisation
  the multitenant plugin resolved and the actor that per-record ownership is
  scoped by. That is why passing the request's context down matters rather than
  creating a fresh one: a `context.Background()` in the middle of a request
  loses the caller's identity, and the query comes back unscoped.

  When a background job genuinely needs to act as nobody, there is an explicit
  way to say so rather than an accident.
]

== Where the rules live, concretely

Open any generated service and you find three things above the methods:

```go
// What a client may search, sort and filter by. Whitelisted, because each
// name ends up in SQL.
var contactListConfig = paginate.Config{
    Searchable: []string{"name", "phone"},
    Sortable:   map[string]bool{"id": true, "created_at": true,
        "name": true, "phone": true},
    Filterable: map[string]bool{"id": true, "name": true,
        "phone": true, "group_id": true},
}

// Every column Patch and Bulk may write. id, the timestamps and
// the version are the framework's, and are dropped.
var writableContact = map[string]bool{
    "name": true, "phone": true, "group_id": true,
}
```

These are the answer to chapter twenty two's warning that a column name cannot
be a placeholder. The client sends `?sort=name`, and `name` is checked against
a map the generator wrote, rather than being interpolated into SQL.

#bridge[
  This is the difference between an ORM protecting you and a convention
  protecting you. Prisma builds the query so the column is validated by its
  schema. Here the query is yours, so the whitelist is explicit and sitting at
  the top of the file where a reviewer sees it.
]

== The event bus

A write emits once, and the systems that care subscribe.

```go
events.Emit(c, events.Event{
    Name: "invoices.paid", Resource: "invoices", ID: inv.ID, After: inv,
})
```

The audit log, outbound webhooks and realtime clients are *subscribers*, not
calls the handler has to remember. Adding a fourth consumer does not mean
editing every handler.

Three delivery modes, and choosing correctly matters:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`Sync`], [Runs inline, before the response. For the audit record, which must be written before the caller is told the write succeeded. A slow Sync subscriber is a slow endpoint],
  [`Async`], [Runs on a worker pool after the handler returns. For anything with a network call: a webhook that takes four seconds must not make the API take four seconds],
  [`Durable`], [Through the transactional outbox: at least once, retried, surviving a restart. For work another module depends on],
)

`Durable` is the one worth understanding, because it is the difference between
an event system and a reliable one. The event is written *in the same database
transaction as the data*, so a crash between the commit and the HTTP call loses
nothing. A relay then delivers it, claiming messages with a lease so another
replica takes over if one dies.

#bridge[
  The naive version is a `POST` at the end of the handler. It works in the demo
  and loses events in production: the commit succeeded, the HTTP call did not,
  and nothing records that anything was owed. The outbox is the standard fix
  and it is already wired, which is B13 from chapter twenty six.
]

== What it deliberately is not

Four, each with a reason rather than a dismissal:

*Not microservices.* One binary by design. A split is worth paying for when
teams need to release independently, and until then it buys network failures
and a deployment pipeline in exchange for nothing.

*Not CQRS or event sourcing.* Both earn their complexity in a hot spot, not
across a codebase. The read path is already separate enough to point at another
store when one query needs it.

*Not hexagonal.* The generated service talks to GORM directly rather than
through a port. Less indirection to read, at the cost of swapping the ORM being
real work. That is a trade, stated as one.

*Not serverless.* Jobs, a scheduler, realtime connections and a warm pool all
assume a process that stays up.

== Organised by layer, not by feature

A fair question, since package-by-feature is what many Go teams land on:
why is it `models/`, `services/`, `handlers/` rather than `invoice/`,
`customer/`?

What feature folders buy is real: a feature is a folder, you can delete one by
deleting a directory, and the root stops growing. At forty resources Grit's
`services/` is forty files and the directory stops telling you which three form
a subsystem.

What the layered layout buys back is the *entity graph*. A business application
has thirty entities pointing at each other: invoice to customer to organisation
and back. *Go forbids import cycles*, so with feature folders the moment
`customer` refers to `invoice` you need an event or an interface dance to
express what is really a foreign key.

Grit's thirty-odd models are one flat `package models`. Any model can reference
any other with no ceremony, which is what lets `client:belongs_to` be a single
word in a generate command rather than a design conversation.

#gotcha[
  The import-cycle rule still bites, in the other direction. `models` cannot
  import `services`, so a hook on a model calls the lower-level package
  directly rather than the service that wraps it. If you add a `BeforeCreate`
  that needs a sequence number, call the sequence package, not the service.

  This is the kind of constraint that is obvious once stated and costs an hour
  when discovered by trial.
]

== Can you refactor to something else?

Yes, and the useful part of the answer is this: *what decides whether a
refactor is a weekend or a quarter is almost never the code. It is the data.*

Moving a package is easy. Moving a table out from under three foreign keys, two
joins and an admin screen that lists both sides of a relationship is the real
work, and no amount of clean code makes it cheaper.

#table(
  columns: (auto, auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [*CQRS on one endpoint*], [hours], [The read path is already a separate whitelist. Point one query at a read replica],
  [*Extract one module*], [weeks], [The service is the seam. The cost is the foreign keys you must cut],
  [*Fully hexagonal*], [ongoing], [Mechanical, and the generator keeps writing the un-ported version],
  [*Microservices*], [months], [Only with an organisational reason. "It will scale" is not one],
  [*Event sourcing*], [a rewrite], [State stops being a row. Not a refactor],
  [*Collapse services in*], [weeks], [Easier than the reverse, and often the right direction],
)

#project(title: "Project 25: trace a request end to end")[
  No new code. Follow one request through every layer of the project you made
  last chapter, with the profiler open, and the architecture stops being a
  diagram.

  ```bash
  cd bookshop
  grit start
  ```

  *One.* Open `apps/api/internal/routes/` and find where the user routes are
  registered. Note which middleware guards them and in what order. The order is
  the composition from chapter twenty one: outermost runs first.

  *Two.* Open the handler that route points at. Count its lines, then count how
  many are about HTTP rather than about users. Follow the one service call.

  *Three.* In the service, find `paginate.Config` and the writable map. Make a
  request with `?sort=nonsense` and see what happens. Then make one with a
  column that *is* in the whitelist and watch it work. You have just tested the
  defence from chapter twenty two.

  *Four.* Open `/pulse` and make the same request again. Note the query count
  for one request. Then create a resource with a `belongs_to` relationship,
  list it, and watch the query count: if it jumps by the number of rows, you
  are looking at an N+1, and the "Issued from" column names the file and line
  that caused it.

  *Five.* Make a write, then open `/studio` and look at the audit table. The
  handler never called the audit log; it emitted an event, and audit was a
  subscriber. That indirection is the event bus doing its job.
]

#tryit[
  Find a service method that takes `ctx context.Context` and call it from a
  `cmd/` command with `context.Background()`, bypassing HTTP entirely. It
  should work, because the service never knew it was being called over HTTP.
  That is the thin-handler rule paying off, and proving it once is worth more
  than reading the argument for it.
]
