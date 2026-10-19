#import "../template.typ": *

= Goroutines, against the event loop

You already write concurrent code. Every `await` you have typed was
concurrency. What you have not had is *parallelism*: two pieces of your own
code running at the same instant on two cores. This chapter is about the
difference, because it is the whole reason Go's model exists.

== The one-keyword version

Put `go` in front of a function call and it runs independently.

```go
func main() {
    go fmt.Println("from a goroutine")
    fmt.Println("from main")
    time.Sleep(100 * time.Millisecond)
}
```

That is the entire syntax. No `async` on the function, no `await` at the call,
no promise returned. Any function can be a goroutine; nothing about its
declaration changes.

#gotcha[
  The `time.Sleep` is there because *when `main` returns, the program exits
  immediately*, killing every goroutine still running. Remove it and you will
  probably see only "from main".

  Sleeping is never the real answer. The real answers are `sync.WaitGroup` and
  channels, both below. But almost every Go tutorial starts with a sleep, and
  it is worth knowing that it is a placeholder for synchronisation rather than
  a technique.
]

== What is actually different

In Node, your JavaScript runs on *one* thread. The event loop interleaves work:
while one request waits on the database, another gets a turn. That is excellent
for I/O and useless for computation, because a function that spends two hundred
milliseconds doing arithmetic blocks everything.

Go runs goroutines on a pool of operating system threads, usually one per CPU
core. Two goroutines can execute at literally the same moment.

#table(
  columns: (auto, 1fr, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 6pt,
  [], [*Node*], [*Go*],
  [Unit], [Promise / task], [Goroutine],
  [Threads], [One for your code], [One per core, by default],
  [Waiting on I/O], [Excellent], [Excellent],
  [CPU work], [Blocks everything], [Uses every core],
  [Cost each], [Small object], [About 2 KB of stack],
  [Practical limit], [Many thousands], [Hundreds of thousands],
)

A goroutine starts with a two kilobyte stack that grows as needed. An operating
system thread starts with one or two *megabytes*. That thousand-fold difference
is why "start a goroutine per request" is a normal design in Go and "start a
thread per request" was never viable in Java or C\#.

#bridge[
  The mental shift: in Node you *structure* code so it never blocks, because
  blocking is catastrophic. In Go you write straight-line blocking code, and
  the runtime moves the goroutine off the thread while it waits. There is no
  colour to functions, no `async` infecting every caller up the stack. A
  function that reads a file looks like a function that reads a file.
]

== Waiting properly: sync.WaitGroup

A `WaitGroup` counts outstanding work.

```go
func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1)                   // one more to wait for
        go func() {
            defer wg.Done()         // one fewer, even if this panics
            fmt.Println("worker", i)
        }()
    }

    wg.Wait()                       // blocks until the count is zero
    fmt.Println("all done")
}
```

Three rules, and the bugs that come from breaking them:

*`wg.Add` before the `go`, never inside it.* Inside, `Wait` may run before the
goroutine starts and see a count of zero.

*`defer wg.Done()` as the first line.* A panic or an early return otherwise
leaves the count above zero and `Wait` blocks forever.

*Pass the `WaitGroup` by pointer* if it crosses a function boundary. Copying
one gives the callee its own counter, and the original never reaches zero.

Since Go 1.25 there is a shorthand that removes the first two:

```go
var wg sync.WaitGroup
for i := 1; i <= 3; i++ {
    wg.Go(func() {              // Add and Done handled for you
        fmt.Println("worker", i)
    })
}
wg.Wait()
```

== Collecting results

Goroutines return nothing. A `go` call is a statement, not an expression, so
there is no value to take.

Write into a pre-sized slice, one index per goroutine. No lock is needed
because no two goroutines touch the same element.

```go
type result struct {
    url string
    ms  int64
    err error
}

func fetchAll(urls []string) []result {
    results := make([]result, len(urls))      // sized up front
    var wg sync.WaitGroup

    for i, url := range urls {
        wg.Go(func() {
            start := time.Now()
            resp, err := http.Get(url)
            if err != nil {
                results[i] = result{url: url, err: err}
                return
            }
            defer resp.Body.Close()
            io.Copy(io.Discard, resp.Body)    // read it, or the connection leaks
            results[i] = result{url: url, ms: time.Since(start).Milliseconds()}
        })
    }

    wg.Wait()
    return results
}
```

#bridge[
  This is `Promise.all`, with the ordering preserved the same way: by index
  rather than by completion. The difference is that an error in one goroutine
  does not reject the whole batch. Here every result carries its own error,
  which is usually what you wanted from `Promise.allSettled` anyway.

  Appending to a shared slice instead of indexing would be a data race, because
  `append` reads and writes the slice header. Index assignment to distinct
  elements is safe.
]

== Races, and the detector that finds them

Two goroutines touching the same variable, at least one of them writing, with
no synchronisation, is a *data race*. The result is not "the wrong number", it
is undefined.

```go
func main() {
    count := 0
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Go(func() { count++ })    // RACE
    }
    wg.Wait()
    fmt.Println(count)               // 973, 981, 1000, who knows
}
```

`count++` is three operations: read, add, write. Two goroutines interleaving
them lose updates.

Go ships a race detector, and it is the most valuable tool in the toolchain:

```bash
go run -race main.go
go test -race ./...
```

```
WARNING: DATA RACE
Write at 0x00c0000140a8 by goroutine 8:
  main.main.func1()
      /tmp/main.go:12 +0x3c
Previous write at 0x00c0000140a8 by goroutine 7:
  ...
```

#gotcha[
  Run your tests with `-race` in CI. Always. A data race can run correctly ten
  million times and corrupt data on the ten-million-and-first, on a different
  machine, under load. The detector finds it the first time the racing access
  happens, which turns an unreproducible production mystery into a failing test.

  It costs roughly ten times the CPU and a few times the memory, so it is a CI
  and development tool rather than something you ship.
]

The fix for the counter is a mutex:

```go
var (
    mu    sync.Mutex
    count int
)

wg.Go(func() {
    mu.Lock()
    defer mu.Unlock()
    count++
})
```

Or, for exactly this shape, an atomic:

```go
var count atomic.Int64
wg.Go(func() { count.Add(1) })
fmt.Println(count.Load())
```

== Bounding how many run at once

Unbounded goroutines are a real hazard. Ten thousand URLs means ten thousand
simultaneous connections, which will exhaust file descriptors or get you
blocked.

A buffered channel used as a semaphore caps the concurrency:

```go
func fetchAll(urls []string, limit int) []result {
    results := make([]result, len(urls))
    sem := make(chan struct{}, limit)    // room for `limit` tokens
    var wg sync.WaitGroup

    for i, url := range urls {
        wg.Go(func() {
            sem <- struct{}{}            // take a token, blocking if full
            defer func() { <-sem }()     // give it back

            // ... fetch, as before
            _ = i
            _ = url
        })
    }

    wg.Wait()
    return results
}
```

Channels are the next chapter. For now, read `sem <- struct{}{}` as "wait until
there is room, then claim a slot".

#bridge[
  The JavaScript equivalent is a library like `p-limit`, because the language
  gives you nothing for it. Here it is two lines and a channel, which is a fair
  summary of the trade Go makes across the board: fewer batteries, better
  primitives.
]

== When not to use a goroutine

Starting one is cheap, which makes it tempting. Three cases where it is wrong:

*When the work is faster than the coordination.* Spawning a goroutine to add
two numbers costs more than adding them.

*When you do not wait for it.* A goroutine whose result nobody collects, with
no `WaitGroup` and no channel, is a leak. If `main` exits it is killed
silently; if it blocks forever it holds its stack forever.

*When it needs ordering.* Goroutines have no execution order. If step two must
follow step one, write them in sequence.

#project(title: "Project 12: a concurrent link checker")[
  Check a list of URLs in parallel, with a concurrency limit and a timeout, and
  report which are broken. This is a genuinely useful tool and it exercises
  everything in this chapter.

  ```go
  package main

  import (
      "fmt"
      "io"
      "net/http"
      "os"
      "sort"
      "sync"
      "time"
  )

  type check struct {
      url    string
      status int
      ms     int64
      err    error
  }

  func checkAll(urls []string, limit int) []check {
      results := make([]check, len(urls))
      sem := make(chan struct{}, limit)
      var wg sync.WaitGroup

      client := &http.Client{Timeout: 10 * time.Second}

      for i, url := range urls {
          wg.Go(func() {
              sem <- struct{}{}
              defer func() { <-sem }()

              start := time.Now()
              resp, err := client.Get(url)
              if err != nil {
                  results[i] = check{url: url, err: err}
                  return
              }
              defer resp.Body.Close()
              io.Copy(io.Discard, resp.Body)

              results[i] = check{
                  url:    url,
                  status: resp.StatusCode,
                  ms:     time.Since(start).Milliseconds(),
              }
          })
      }

      wg.Wait()
      return results
  }

  func main() {
      urls := os.Args[1:]
      if len(urls) == 0 {
          fmt.Fprintln(os.Stderr, "usage: check <url>...")
          os.Exit(1)
      }

      start := time.Now()
      results := checkAll(urls, 5)
      total := time.Since(start)

      // Slowest first, so the problems are at the top.
      sort.Slice(results, func(a, b int) bool {
          return results[a].ms > results[b].ms
      })

      broken := 0
      for _, r := range results {
          switch {
          case r.err != nil:
              broken++
              fmt.Printf("  ERR  %-45s %v\n", r.url, r.err)
          case r.status >= 400:
              broken++
              fmt.Printf("  %d  %-45s %dms\n", r.status, r.url, r.ms)
          default:
              fmt.Printf("  %d  %-45s %dms\n", r.status, r.url, r.ms)
          }
      }

      fmt.Printf("\nchecked %d urls in %v, %d broken\n",
          len(results), total.Round(time.Millisecond), broken)
      if broken > 0 {
          os.Exit(1)
      }
  }
  ```

  Run it against a handful of sites and notice that the total time is roughly
  the *slowest* request, not the sum of all of them. That is the payoff.

  Three details that matter in real code. The `http.Client` has a timeout:
  without one, a hung server holds a goroutine forever and `wg.Wait` never
  returns. The response body is read and closed, because an unread body leaks
  the connection. And the exit code is non-zero when something is broken, which
  is what makes this usable in CI.

  Now run it with `-race`. It will be clean, because each goroutine writes only
  to `results[i]`. Then introduce a race on purpose: add a shared
  `var count int` and `count++` inside the goroutine, and run with `-race`
  again to see the detector's report.
]

#tryit[
  Add a `-n` flag for the concurrency limit and try the same URL list at 1, 5
  and 50. At 1 the total is the sum of all requests. At 5 it drops sharply. At
  50 it probably stops improving, because you are now limited by the network or
  by the servers themselves. Finding the point where more concurrency stops
  helping is the actual skill, and it is measured rather than guessed.
]
