#import "../template.typ": *

= Context: cancellation, deadlines, and request scope

You will see `ctx context.Context` as the first parameter of almost every
function in a Go codebase, and it looks like noise until you understand what it
buys. It is the done channel from chapter fourteen, standardised so that every
library in the ecosystem agrees on how to be cancelled.

== The problem it solves

A user opens a page, your handler starts a database query, and the user closes
the tab. Without cancellation the query runs to completion, the connection
stays busy, and the result is thrown away. Multiply by a thousand users on a
slow page and the database falls over serving nobody.

Context is a signal that travels *down* the call stack so that work can stop.

== Creating one

You rarely create a context from nothing. Two roots exist:

```go
ctx := context.Background()   // the root, in main or a top-level handler
ctx := context.TODO()         // a placeholder: "a context belongs here"
```

`TODO` is for when you are threading context through an existing codebase and
have not reached the caller yet. It behaves identically; it exists so that
`grep TODO` finds the unfinished work.

From a root you *derive*:

```go
ctx, cancel := context.WithCancel(parent)
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
ctx, cancel := context.WithDeadline(parent, someTime)
```

Each returns a child and a `cancel` function.

#gotcha[
  *Always call cancel, even on the timeout versions.* It is not optional:

  ```go
  ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
  defer cancel()
  ```

  Not calling it leaks the timer and the goroutine watching it until the
  deadline passes. `go vet` catches the obvious cases and not all of them.
  Writing `defer cancel()` on the very next line is the habit to build.
]

== Using one

A context has a `Done()` channel that closes when the context is cancelled, and
an `Err()` saying why.

```go
func work(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()           // Canceled or DeadlineExceeded
        default:
        }
        doChunk(i)
    }
    return nil
}
```

That is the shape: check between units of work, return promptly when told.

For anything that blocks, put `ctx.Done()` in the `select` directly:

```go
select {
case result := <-ch:
    return result, nil
case <-ctx.Done():
    return nil, ctx.Err()
}
```

#bridge[
  This is `AbortController` and `AbortSignal`. You create a controller, pass
  `signal` to `fetch`, and call `abort()`. Go's version has two differences:
  the convention of passing it as the *first parameter of every function* is
  near-universal, so the whole ecosystem supports it, and contexts form a tree
  where cancelling a parent cancels every descendant.
]

== Cancellation is cooperative

Nothing is killed. A context cancellation sets a flag and closes a channel;
code that does not check it keeps running.

```go
func ignoresCancellation(ctx context.Context) {
    time.Sleep(10 * time.Second)    // runs the full ten seconds regardless
}
```

There are no interrupts and no thread killing in Go. If a library does not take
a context, it cannot be cancelled, and that is a reason to prefer libraries
that do.

== The standard library already honours it

This is where it pays off, because you mostly pass a context rather than check
one:

```go
req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := http.DefaultClient.Do(req)   // aborts if ctx is cancelled

rows, err := db.QueryContext(ctx, "SELECT ...")   // cancels the query
```

An incoming HTTP request already carries one, cancelled automatically when the
client disconnects:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // cancelled if the client goes away
    ctx := r.Context()
    rows, err := db.QueryContext(ctx, "SELECT ...")
    ...
}
```

So passing `r.Context()` into your service layer, which passes it to the
database, gives you end-to-end cancellation for free. That is the single most
valuable thing context does, and it requires no `select` anywhere in your code.

== Timeouts compose

A derived context cannot outlive its parent, which gives you sane behaviour
without arithmetic.

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Overall budget for this request.
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    // One slow dependency gets a tighter budget.
    userCtx, cancel2 := context.WithTimeout(ctx, 500*time.Millisecond)
    defer cancel2()
    u, err := users.Get(userCtx, id)
    ...
}
```

If the request's three seconds elapse first, the inner context is cancelled
too, even though its own 500ms had not expired.

== Values, and why to be sparing

A context can carry request-scoped values.

```go
type ctxKey string
const userKey ctxKey = "user"

ctx = context.WithValue(ctx, userKey, currentUser)

u, ok := ctx.Value(userKey).(*User)
```

Use an unexported custom key type, never a plain string, or two packages using
`"user"` collide silently.

#gotcha[
  Context values are the most abused part of the package. They are untyped,
  invisible in a function signature, and found at runtime or not at all.

  The rule from the standard library documentation: use them for *request-scoped
  data that crosses API boundaries*, such as a request ID, an authenticated
  user, or a trace span. Never for optional parameters. If a function needs a
  value to do its job, that value is a parameter:

  ```go
  // No.
  func Send(ctx context.Context) error    // reads the recipient from ctx

  // Yes.
  func Send(ctx context.Context, to string) error
  ```
]

== The conventions

Four, and they are followed almost universally.

*Context is the first parameter, named `ctx`.* `func Get(ctx context.Context, id int)`.

*Never store it in a struct.* A context describes one call's lifetime; a struct
outlives calls. The exception is `http.Request`, which has its own reasons.

*Never pass nil.* Use `context.TODO()` if you have nothing better.

*Do not pass a context to something that outlives the request.* Starting a
background job with the request's context means the job is cancelled when the
response is sent. Use `context.WithoutCancel(ctx)` to keep the values and drop
the cancellation.

#project(title: "Project 15: a search that races two sources with a budget")[
  Query two backends at once, take whichever answers first, give up on both
  after a deadline, and never leak a goroutine. This is context, channels and
  select doing the job they exist for.

  ```go
  package main

  import (
      "context"
      "errors"
      "fmt"
      "math/rand"
      "time"
  )

  type Result struct {
      Source string
      Hits   int
  }

  // search simulates a backend that takes an unpredictable amount of time
  // and
  // respects cancellation.
  func search(ctx context.Context, source string,
      slow time.Duration) (Result, error) {
      select {
      case <-time.After(slow):
          return Result{Source: source, Hits: rand.Intn(100)}, nil
      case <-ctx.Done():
          // Stop promptly and say why.
          return Result{}, fmt.Errorf("%s: %w", source, ctx.Err())
      }
  }

  // first returns whichever source answers first, cancelling the other.
  func first(ctx context.Context, budget time.Duration) (Result, error) {
      ctx, cancel := context.WithTimeout(ctx, budget)
      defer cancel()              // cancels the loser, and the timer

      type outcome struct {
          res Result
          err error
      }
      // Buffered with room for both, so the loser's send never blocks and
      // its goroutine always exits. With an unbuffered channel the loser
      // would block forever on the send: a goroutine leak.
      ch := make(chan outcome, 2)

      go func() {
          r, err := search(ctx, "index-a", 300*time.Millisecond)
          ch <- outcome{r, err}
      }()
      go func() {
          r, err := search(ctx, "index-b", 120*time.Millisecond)
          ch <- outcome{r, err}
      }()

      var lastErr error
      for i := 0; i < 2; i++ {
          select {
          case o := <-ch:
              if o.err == nil {
                  return o.res, nil      // first success wins
              }
              lastErr = o.err
          case <-ctx.Done():
              return Result{}, ctx.Err()
          }
      }
      return Result{}, lastErr
  }

  func main() {
      // Generous budget: the fast index wins.
      r, err := first(context.Background(), time.Second)
      fmt.Printf("got %+v err=%v\n", r, err)

      // Tight budget: both are cancelled.
      r, err = first(context.Background(), 50*time.Millisecond)
      fmt.Printf("got %+v err=%v\n", r, err)
      fmt.Println("deadline:",
          errors.Is(err, context.DeadlineExceeded))
  }
  ```

  Three things this program teaches.

  The buffered channel is a correctness requirement, not an optimisation. When
  `first` returns after the winner, the loser still sends. With no buffer that
  send blocks forever and the goroutine never exits. Change the `2` to a `0`,
  run with `-race` and then inspect goroutine counts, and you have reproduced
  one of the most common leaks in Go services.

  The `defer cancel()` is what makes the loser stop early rather than running
  its full three hundred milliseconds into a void.

  And `errors.Is(err, context.DeadlineExceeded)` works through the `%w` wrap in
  `search`, which is chapter eleven's design paying off: the caller can tell a
  timeout from any other failure.
]

#tryit[
  Add a third source and make `first` return as soon as *two* of the three
  agree, rather than on the first answer. Then add a `context.WithoutCancel`
  branch that logs the losing results to an audit sink after the winner has
  been returned, and confirm that work still completes after `first` has
  returned. That is the pattern for "respond now, finish the bookkeeping
  afterwards", and getting the context wrong there is why audit records go
  missing under load.
]
