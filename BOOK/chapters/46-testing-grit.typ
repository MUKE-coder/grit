#import "../template.typ": *

= Testing a Grit application

Chapter eighteen taught the mechanics. This chapter is the harder question: in
an application where most of the code was generated, *what is worth testing*,
and which tests pay for the time they cost.

The short answer: do not test the generator. Test the things you wrote, the
boundaries you configured, and the one property whose failure is silent.

== What arrives already written

A new project comes with tests, and they are worth reading before you add any:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`auth_test.go`], [Register and login, success and every failure],
  [`user_test.go`], [The auth guard, an admin list, a 404, a success],
  [`bench_test.go`], [Health check, login, register],
  [`perf_test.go`], [Allocation budgets for the per-request hot path],
  [`__tests__/`], [Component tests with Vitest and Testing Library],
  [`e2e/`], [Playwright: the auth flow and the admin],
)

```bash
grit test                 # Go and frontend
grit test --go --race     # Go only, with the race detector
grit test --cover         # Go coverage
grit test --e2e           # include Playwright, app must be running
```

One report, and suites that do not apply to your architecture are reported as
*skipped with the reason* rather than quietly dropped. A test suite that
silently does not run is worse than no suite, because it reports green.

== The in-memory database trap

This one is worth the whole chapter, because it produces a test that fails six
runs in ten and looks like a broken login.

```go
db, err := gorm.Open(sqlite.Open(":memory:"), ...)
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(1)     // this line is the whole fix
```

*Every connection to `:memory:` gets its own empty database.* The pool opens a
second connection as soon as two queries overlap, and in this project they do,
because registering a user writes an activity row on a goroutine. The login
that follows lands on a fresh, empty database, finds no user, and returns 401.

#gotcha[
  The failure reads as a broken login rather than a test harness problem, which
  is why it costs an afternoon. The tell is the flakiness: *a test that passes
  about six runs in ten is almost never a race in your code and almost always a
  race in your fixture*.

  One connection, and only one. If you need concurrency in a test, use a file
  database in a temporary directory instead.
]

The same harness makes a second decision worth copying. The email verification
token table is *deliberately not migrated*, because that write happens on a
goroutine so a slow SMTP call cannot hold the response open. Creating the table
would turn a harmless early return into a background writer racing whatever the
test does next.

If you are testing verification itself, migrate that table in that test and
call the service directly rather than through the register handler.

== What is worth testing, and what is not

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Do not test], [That a generated create returns 201. The generator's own
    tests cover it, and your copy adds nothing],
  [Test], [Every rule you wrote in a service, especially the refusals],
  [Test], [The whitelists, by asserting their contents],
  [Test], [Ownership isolation between two users],
  [Test], [The error boundary: each error kind maps to the status you expect],
  [Test], [Anything that broke once. A bug without a test comes back],
)

The whitelist one is not obvious, and chapter forty's date-sorting story is
why: the pagination layer *drops* a sort it does not recognise rather than
rejecting it, so a missing column fails silently and the endpoint still returns
200.

```go
func TestInvoiceSortableColumns(t *testing.T) {
    want := []string{"number", "issued", "due", "total_amount"}
    for _, col := range want {
        if !invoiceListConfig.Sortable[col] {
            t.Errorf("%s is not sortable", col)
        }
    }
}
```

Four lines, and they catch a class of bug that no endpoint test catches.

== Service tests, without HTTP

Most of your tests belong here. The rules live in services, so testing them
through HTTP adds a router, a body, a status code and a JSON decode to every
assertion, for nothing.

```go
func TestCloseRefusesOpenSubtasks(t *testing.T) {
    svc, ctx := newTaskService(t)

    parent := mustCreate(t, svc, ctx, "parent")
    mustCreateChild(t, svc, ctx, parent.ID, "child")

    _, err := svc.Close(ctx, parent.ID)
    if !errors.Is(err, ErrHasOpenSubtasks) {
        t.Fatalf("err = %v, want ErrHasOpenSubtasks", err)
    }
}
```

Note the assertion is on the *sentinel*, not on a message. The message is
allowed to change; the identity of the error is the contract, which is chapter
eleven's whole point arriving in a test file.

== Handler tests, for the contract

Test at the HTTP level only what is about HTTP: the status code, the envelope,
the headers.

```go
func TestStaleVersionReturns409WithNumericVersion(t *testing.T) {
    api := newTestAPI(t)
    rec := patchInvoice(t, api, 1, "If-Match: 1", `{"title":"x"}`)

    if rec.Code != 409 {
        t.Fatalf("status = %d, want 409", rec.Code)
    }
    var body struct {
        Error struct {
            CurrentVersion int `json:"current_version"`
        } `json:"error"`
    }
    if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
        t.Fatalf("decode: %v", err)
    }
}
```

The decode into an `int` is the assertion. If the server ever regresses to
sending `"9"` as a string, this fails, which is the bug chapter thirty three
said shipped once.

== The highest-value test in the whole project

```go
func TestNotesAreIsolatedBetweenUsers(t *testing.T) {
    svc, _ := newNoteService(t)
    ada := ctxFor(t, "ada")
    grace := ctxFor(t, "grace")

    note := mustCreate(t, svc, ada, "ada's note")

    _, err := svc.GetByID(grace, note.ID)
    if !errors.Is(err, ErrNotFound) {
        t.Fatalf("grace could read ada's note: %v", err)
    }
}
```

In a multi-user application this is the test to write first and run in CI
forever. The bug it catches is silent, is not a crash, produces no error in any
log, and is the one you cannot apologise your way out of.

#gotcha[
  Write it once per *ownership-scoped resource*, not once per project. The
  scoping is applied per service, and the one that was added last month by
  somebody in a hurry is the one missing `s.db(ctx)`.

  A table-driven version over every scoped service is about thirty lines and
  catches the regression on the day it is introduced.
]

== Allocation budgets, and a trap in the toolchain

`perf_test.go` asserts how many allocations the hot path makes:

```
token verification: 56 allocations under go1.26.9
                    41 allocations under go1.27.2
```

A benchmark tells you how long something took on the machine that ran it, which
on a shared CI runner moves on its own. An allocation count does not move with
the machine or the load, which is what makes these assertions rather than
measurements, and why they can fail a build when a benchmark cannot.

#gotcha[
  Allocation counts *do* move with the Go toolchain, by more than you would
  guess: fifteen fewer allocations for the same work between two releases.

  So a budget is set from the *oldest* toolchain the project supports, because
  that is the one that allocates most. Set it from whatever is installed on
  your laptop and it passes locally and fails in CI, which is how that file
  first failed.

  When a budget fails, it is telling you a change allocates more than it used
  to. That is not automatically wrong: read the change, decide whether the cost
  is worth it, and if it is, raise the budget *in the same commit*, so the next
  person can see the new number was a decision.
]

== The frontend, in two layers

*Component tests* with Vitest and Testing Library, for logic that lives in the
browser: a form that disables its button, a list that shows an empty state, a
hook that retries. Fast, and they run on every save.

*End-to-end tests* with Playwright, for the three or four flows that must never
break: sign in, create the central record, and whatever your application exists
to do. They are slow and they are flaky in proportion to how many you have, so
keep the set small and keep it meaningful.

#bridge[
  The temptation is to write end-to-end tests for everything, because they feel
  like proof. The cost arrives later: a suite that takes twenty minutes and
  fails for unrelated reasons twice a week gets marked "allowed to fail", and
  then it is decoration.

  Rule of thumb: an end-to-end test for every flow you would stop a release
  for, and not one more.
]

== What CI should run

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Every push], [`go build`, `go vet`, Go tests, component tests, type check],
  [Every push], [The race detector on the packages with concurrency],
  [Every push], [The isolation tests. They are fast and they are the ones that
    matter],
  [Nightly], [End-to-end, the full race suite, the dependency drift check],
)

One practical note from this project's own CI: *never pipe `go test ./...`
through `tail`*. Keep the whole log and grep it when something fails. A
truncated log turns a five minute diagnosis into a re-run with more logging.

#project(title: "Project 43: make the suite worth having")[
  Take the invoicing application from chapter forty one and bring its tests to
  a state you would defend.

  *One: the rules.* A test for every refusal in the state machine, asserting
  sentinels rather than messages.

  *Two: the whitelists.* A test asserting the sortable and filterable column
  sets. Then remove a column from the service, watch the test fail, and
  confirm that the endpoint test still passes. That gap is the point.

  *Three: isolation.* Two users, every scoped resource, table driven.

  *Four: the flake.* Remove `SetMaxOpenConns(1)` from the test harness and run
  the suite ten times. Record how many passed. Put it back.

  *Five: the contract.* A handler test for the 409 that decodes
  `current_version` into an `int`.

  *Six: end to end.* Exactly three Playwright tests: sign in, create an
  invoice, send it. Time the suite and write the number down.

  *Seven: CI.* A workflow running all of it on push, with the end-to-end suite
  on a schedule instead. Push a deliberate break and confirm it fails for the
  right reason.
]

#tryit[
  Open your largest test file and count how many tests would still pass if the
  business rules were deleted.

  Every one of those is testing the framework rather than your application.
  Deleting them makes the suite faster and no weaker, and the exercise tends to
  be uncomfortable in a useful way.
]

#assignment(name: "Testing Grit - Assignment One")[
  In `part5-building/46-testing`, build the suite and report what it is for.

  *One.* Service tests for every rule you wrote, handler tests only for the
  contract, and the isolation tests table driven across every scoped resource.

  *Two.* The whitelist test, with evidence that removing a column fails it and
  does not fail an endpoint test.

  *Three.* `TESTS.md` naming, for each suite, the bug class it exists to catch
  and roughly how long it takes to run.

  *Four.* The flake experiment: ten runs with and without the single
  connection, with the counts.

  *Five.* One bug you have actually hit in your own work, written as a
  regression test. If you cannot remember one, go through your commit history
  for the word "fix" and take the first one.
]
