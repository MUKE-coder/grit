#import "../template.typ": *

= Channels, and the select statement

Goroutines run independently. Channels are how they talk. The slogan Go is
famous for is *do not communicate by sharing memory; share memory by
communicating*, which sounds like a fortune cookie until you have written both
versions of the same program.

== A channel is a typed pipe

```go
ch := make(chan int)    // unbuffered
ch <- 42                // send
v := <-ch               // receive
```

The arrow always points in the direction the value travels. `ch <- v` puts a
value in, `<-ch` takes one out.

```go
func main() {
    ch := make(chan string)

    go func() {
        ch <- "hello from the goroutine"
    }()

    msg := <-ch          // blocks until something arrives
    fmt.Println(msg)
}
```

No `time.Sleep` this time. The receive blocks until the send happens, which is
synchronisation and message passing in one operation.

#bridge[
  The closest JavaScript analogy is an async generator or a stream, with one
  crucial difference: a channel receive *blocks the goroutine*. There is no
  callback, no `await`, and no promise. The code reads top to bottom, and the
  runtime parks the goroutine while it waits.

  If you have used `postMessage` between workers, the mental model is closer:
  separate lines of execution passing values, except channels are typed and
  synchronous by default.
]

== Unbuffered channels are a handshake

An unbuffered channel has no room. A send blocks until a receiver is ready, and
a receive blocks until a sender is ready. The two goroutines meet.

```go
ch := make(chan int)
ch <- 1                 // deadlock: nobody is receiving
```

```
fatal error: all goroutines are asleep - deadlock!
```

Go detects when every goroutine is blocked and stops the program rather than
hanging forever. This is a feature, and it is how you will find most channel
bugs.

== Buffered channels have room

```go
ch := make(chan int, 3)  // capacity 3
ch <- 1                  // does not block
ch <- 2
ch <- 3
ch <- 4                  // now it blocks: the buffer is full
```

A send blocks only when the buffer is full; a receive blocks only when it is
empty.

#gotcha[
  Buffering is not a performance trick to sprinkle on. The capacity is a design
  decision: it is how many items may be in flight before the producer must slow
  down. That back pressure is usually what you want, because an unbounded queue
  in front of a slow consumer is just a memory leak with extra steps.

  Choose zero when you want a handshake, and a specific number when you know
  how much slack the system should tolerate. "I made it 100 so it would not
  block" is how a queue silently grows to 100 and then blocks anyway.
]

== Closing, and ranging

The sender closes a channel to say "no more values".

```go
func produce(ch chan<- int) {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)           // only the sender closes
}

func main() {
    ch := make(chan int)
    go produce(ch)

    for v := range ch {  // ends when the channel is closed
        fmt.Println(v)
    }
}
```

Receiving from a closed channel returns the zero value immediately. The
two-value form tells you which happened:

```go
v, ok := <-ch
if !ok {
    // the channel is closed and drained
}
```

#gotcha[
  Three rules that cover every close-related bug:

  *Only the sender closes.* A receiver closing a channel causes the sender to
  panic on its next send.

  *Closing twice panics.* So does sending on a closed channel.

  *You do not have to close.* A channel that is garbage collected is fine.
  Close only to signal completion to a `range` or a `select`.
]

Note `chan<- int` in that signature: a *directional* channel type, meaning
send-only. `<-chan int` is receive-only. Using them in function parameters
makes the data flow obvious and lets the compiler stop a receiver from closing.

== select: waiting on several channels

`select` blocks until one of its cases can proceed.

```go
select {
case msg := <-ch1:
    fmt.Println("from ch1:", msg)
case msg := <-ch2:
    fmt.Println("from ch2:", msg)
case ch3 <- "value":
    fmt.Println("sent to ch3")
}
```

If several are ready it picks one at random, which prevents one busy channel
starving another.

A `default` case makes it non-blocking:

```go
select {
case msg := <-ch:
    fmt.Println(msg)
default:
    fmt.Println("nothing waiting")   // does not block
}
```

#bridge[
  `select` is `Promise.race`, with two differences worth noting. It can wait on
  sends as well as receives. And the losing cases are not "cancelled", they
  simply were not chosen: the other channels still hold their values for the
  next receive. `Promise.race` leaves the losers running and their results are
  discarded, which is a frequent source of leaked work in JavaScript.
]

== Timeouts, the most common select

```go
select {
case result := <-work:
    fmt.Println("finished:", result)
case <-time.After(2 * time.Second):
    fmt.Println("timed out")
}
```

`time.After` returns a channel that delivers a value after the duration. This
is the idiomatic timeout, and it is three lines.

#gotcha[
  `time.After` leaks until it fires. In a tight loop, creating one per
  iteration piles up timers for the full duration. For a loop, create a
  `time.NewTimer` once and reset it, or use a context deadline, which is the
  next chapter.
]

== Signalling with a done channel

A channel carrying no data, used only to signal, is a common shape:

```go
func worker(done <-chan struct{}) {
    for {
        select {
        case <-done:
            fmt.Println("stopping")
            return
        default:
            // do a unit of work
        }
    }
}

func main() {
    done := make(chan struct{})
    go worker(done)

    time.Sleep(time.Second)
    close(done)          // closing broadcasts to every receiver
    time.Sleep(10 * time.Millisecond)
}
```

Closing is the broadcast mechanism: *every* goroutine blocked on `<-done`
wakes up at once, because a closed channel always yields immediately. Sending a
value would wake only one.

This pattern is so common that the standard library wrapped it, and that
wrapper is `context`, which is chapter sixteen. Understand the done channel
first and `context` becomes obvious rather than magical.

== The nil channel, and why it is useful

Sending or receiving on a nil channel blocks forever. That sounds like a pure
bug and it has one real use: disabling a `select` case.

```go
for {
    select {
    case v, ok := <-in:
        if !ok {
            in = nil        // disable this case, keep the loop alive
            continue
        }
        process(v)
    case <-done:
        return
    }
}
```

Once `in` is nil its case can never be chosen, so the loop keeps serving `done`
without spinning on a closed channel returning zeros forever.

== Worker pools

The canonical channel program: a fixed number of workers pulling from a jobs
channel.

```go
func worker(id int, jobs <-chan int, results chan<- string, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs {           // ends when jobs is closed
        results <- fmt.Sprintf("worker %d did job %d", id, j)
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan string, 100)
    var wg sync.WaitGroup

    for w := 1; w <= 3; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs)                     // tells the workers no more are coming

    wg.Wait()                       // wait for them to finish
    close(results)                  // now safe to close

    for r := range results {
        fmt.Println(r)
    }
}
```

The ordering is the whole lesson. Close `jobs` so the workers' `range` ends.
`Wait` for the workers. *Then* close `results`, because closing it earlier
would panic a worker still writing.

#bridge[
  Compare this to the semaphore from chapter thirteen. Both bound concurrency.
  The semaphore starts a goroutine per item and limits how many run; the worker
  pool starts a fixed number and feeds them. The pool is better when goroutine
  startup matters or when each worker holds an expensive resource, like a
  database connection. The semaphore is simpler and usually enough.
]

#project(title: "Project 13: a pipeline with cancellation")[
  Build a three-stage pipeline where each stage is a goroutine connected by
  channels, and where the whole thing can be stopped early. This is the shape
  of a great deal of real Go.

  ```go
  package main

  import (
      "fmt"
      "strings"
      "sync"
      "time"
  )

  // Stage 1: produce lines. Returns a channel rather than taking one, which
  // is the convention: the producer owns and closes what it creates.
  func generate(done <-chan struct{}, lines []string) <-chan string {
      out := make(chan string)
      go func() {
          defer close(out)
          for _, l := range lines {
              select {
              case out <- l:
              case <-done:
                  return          // stop early without leaking
              }
          }
      }()
      return out
  }

  // Stage 2: normalise.
  func normalise(done <-chan struct{}, in <-chan string) <-chan string {
      out := make(chan string)
      go func() {
          defer close(out)
          for l := range in {
              select {
              case out <- strings.ToLower(strings.TrimSpace(l)):
              case <-done:
                  return
              }
          }
      }()
      return out
  }

  // Stage 3: count words, fanned out across several goroutines.
  func count(done <-chan struct{}, in <-chan string, workers int) <-chan int {
      out := make(chan int)
      var wg sync.WaitGroup

      for i := 0; i < workers; i++ {
          wg.Go(func() {
              for l := range in {
                  time.Sleep(10 * time.Millisecond)   // pretend it is slow
                  select {
                  case out <- len(strings.Fields(l)):
                  case <-done:
                      return
                  }
              }
          })
      }

      // One goroutine whose only job is to close out once every worker is
      // finished. Without this, nothing ever closes out and main hangs.
      go func() {
          wg.Wait()
          close(out)
      }()

      return out
  }

  func main() {
      lines := []string{
          "  The quick brown fox  ", "jumps over", "THE LAZY DOG",
          "and then", "runs away quickly", "into the woods",
      }

      done := make(chan struct{})
      defer close(done)            // every stage stops when main returns

      start := time.Now()
      total := 0
      for n := range count(done, normalise(done, generate(done, lines)), 3) {
          total += n
      }

      fmt.Printf("%d words in %v\n", total, time.Since(start).Round(time.Millisecond))
  }
  ```

  Run it and note the time: six lines at ten milliseconds each would be sixty
  milliseconds serially, and with three workers it is closer to twenty.

  Then do the experiment that proves the cancellation works. Change the loop in
  `main` to stop after the first two results:

  ```go
  n := 0
  for range count(done, normalise(done, generate(done, lines)), 3) {
      n++
      if n == 2 {
          break
      }
  }
  ```

  The `defer close(done)` fires, every stage's `select` takes the `<-done`
  case, and all the goroutines exit. Remove the `done` channel entirely and the
  producer blocks forever trying to send into a pipeline nobody is reading:
  a goroutine leak, which in a long-running server is a slow memory leak you
  will eventually have to find.

  Run the whole thing with `-race` as well. It is clean, because every value
  moves through a channel and nothing is shared.
]

#tryit[
  Add a fourth stage between normalise and count that drops any line shorter
  than three words, and confirm the total changes. Then count how many
  goroutines the program runs at its peak. If your answer is not "one per
  stage, plus the workers, plus the closer", re-read `count`: that little
  goroutine whose only job is `wg.Wait()` then `close(out)` is the piece people
  forget, and forgetting it is a deadlock rather than a wrong answer.
]
