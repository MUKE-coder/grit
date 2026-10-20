#import "../template.typ": *

= sync: mutexes, and when a channel is the wrong tool

The previous chapter sold channels. This one is the correction: *channels are
not always the answer*, and reaching for one where a mutex belongs produces
code that is slower and harder to read.

The Go team's own guidance, from the wiki: use a channel to pass ownership of
data or to distribute work; use a mutex to protect shared state. The rest of
this chapter is what that means in practice.

== Mutex: one goroutine at a time

```go
type Counter struct {
    mu sync.Mutex
    n  int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.n++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.n
}
```

`Lock` blocks until no other goroutine holds the lock. `defer c.mu.Unlock()` on
the line after is the pattern: a path that returns early without unlocking
deadlocks the whole program.

#gotcha[
  Four rules, each the source of a real bug:

  *The zero value works.* `var mu sync.Mutex` is ready. Never `mu := sync.Mutex{}`
  then copy it around.

  *Never copy a struct containing a mutex.* Copying duplicates the lock, so two
  copies protect nothing. This is why `Inc` has a pointer receiver. `go vet`
  catches most of these.

  *Lock the reads too.* A read racing with a write is still a data race, and
  the race detector will say so. Reading an `int` is not atomic in any sense
  you can rely on.

  *Keep the critical section small.* Never do I/O, call an unknown function, or
  send on a channel while holding a lock. That is how a lock becomes a queue
  and then a deadlock.
]

== Embedding the mutex

Common style, and it reads well:

```go
type SafeMap struct {
    sync.Mutex                   // embedded: m.Lock() works directly
    items map[string]int
}

func (m *SafeMap) Set(k string, v int) {
    m.Lock()
    defer m.Unlock()
    m.items[k] = v
}
```

The cost is that `Lock` and `Unlock` become part of the type's public API, so
callers can lock your struct. Keep the mutex a named field when the type is
exported from a package.

== RWMutex: many readers, one writer

When reads vastly outnumber writes, `sync.RWMutex` lets readers run
concurrently.

```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]string
}

func (c *Cache) Get(k string) (string, bool) {
    c.mu.RLock()                 // read lock: others may also read
    defer c.mu.RUnlock()
    v, ok := c.items[k]
    return v, ok
}

func (c *Cache) Set(k, v string) {
    c.mu.Lock()                  // write lock: exclusive
    defer c.mu.Unlock()
    c.items[k] = v
}
```

#gotcha[
  `RWMutex` is not automatically faster. It has more bookkeeping than a plain
  `Mutex`, so for short critical sections with moderate contention a `Mutex`
  often wins. Use `RWMutex` when reads are frequent *and* the critical section
  is long enough to matter, and confirm with a benchmark rather than a feeling.

  Also: never upgrade an `RLock` to a `Lock` while holding it. That deadlocks.
  Release, acquire the write lock, and re-check the condition, because things
  may have changed in between.
]

== sync.Once: exactly once, no matter who asks

```go
var (
    once   sync.Once
    client *http.Client
)

func Client() *http.Client {
    once.Do(func() {
        client = &http.Client{Timeout: 10 * time.Second}
    })
    return client
}
```

`Do` runs the function once, ever. Concurrent callers block until the first one
finishes, so nobody sees a half-built value.

#bridge[
  This is the lazy singleton you would write in JavaScript as a module-level
  `let instance` with an `if (!instance)` check. That check is a race in Go,
  because two goroutines can both see nil. `sync.Once` is the correct version
  and is shorter anyway.
]

== atomic: for single values

For a counter or a flag, a lock is heavier than necessary.

```go
var requests atomic.Int64

func handle() {
    requests.Add(1)
}

func report() {
    fmt.Println(requests.Load())
}
```

The typed atomics (`atomic.Int64`, `atomic.Bool`, `atomic.Pointer[T]`) arrived
in Go 1.19 and are much nicer than the old function forms. They cannot be
copied, which the compiler enforces.

Use atomics only for a *single* value. The moment two values must change
together, you need a mutex, because two atomic operations are not one atomic
operation:

```go
// WRONG: another goroutine can see total updated and count not.
total.Add(n)
count.Add(1)
```

== sync.Map, and why you probably do not want it

There is a concurrent map in the standard library. It is specialised, and the
documentation says so plainly: it is optimised for two cases, where a key is
written once and read many times, or where goroutines work on disjoint key
sets.

For everything else, a plain map with a mutex is faster and far more readable,
because `sync.Map` has no type safety: everything is `any` and you assert on
the way out.

```go
// Usually better:
type Store struct {
    mu sync.RWMutex
    m  map[string]User
}
```

Reach for `sync.Map` after measuring, not before.

== When to use which

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 6pt,
  [*Channel*], [Passing ownership of a value, distributing work, signalling an event],
  [*Mutex*], [Protecting state that several goroutines read and write in place],
  [*RWMutex*], [The same, when reads dominate and the section is long],
  [*atomic*], [One number or flag, nothing else],
  [*Once*], [Lazy initialisation],
  [*WaitGroup*], [Waiting for a known set of goroutines to finish],
)

The counter from chapter thirteen makes the point. With a channel:

```go
// Works, and is the wrong tool.
func counter(inc <-chan int, get chan<- int) {
    n := 0
    for {
        select {
        case v := <-inc:
            n += v
        case get <- n:
        }
    }
}
```

That is a goroutine, two channels, a `select` loop that never exits, and a
lifetime you now have to manage. Against:

```go
var n atomic.Int64
n.Add(1)
```

The channel version is the one that looks clever and the atomic is the one you
want. *Sharing a number is not communication.*

#project(title: "Project 14: an in-memory cache with expiry")[
  Build something with real shared state: a cache with per-entry expiry and a
  background cleaner. It needs a mutex, a `Once`, an atomic and a goroutine,
  and the combination is typical of real service code.

  ```go
  package main

  import (
      "fmt"
      "sync"
      "sync/atomic"
      "time"
  )

  type entry struct {
      value   string
      expires time.Time
  }

  type Cache struct {
      mu      sync.RWMutex
      items   map[string]entry
      hits    atomic.Int64      // counters: atomic, not under the lock
      misses  atomic.Int64
      stopper sync.Once
      done    chan struct{}
  }

  func NewCache(cleanEvery time.Duration) *Cache {
      c := &Cache{
          items: map[string]entry{},
          done:  make(chan struct{}),
      }
      go c.clean(cleanEvery)
      return c
  }

  func (c *Cache) Set(k, v string, ttl time.Duration) {
      c.mu.Lock()
      defer c.mu.Unlock()
      c.items[k] = entry{value: v, expires: time.Now().Add(ttl)}
  }

  func (c *Cache) Get(k string) (string, bool) {
      c.mu.RLock()
      e, ok := c.items[k]
      c.mu.RUnlock()            // released before the branch below

      if !ok || time.Now().After(e.expires) {
          c.misses.Add(1)
          return "", false
      }
      c.hits.Add(1)
      return e.value, true
  }

  func (c *Cache) clean(every time.Duration) {
      t := time.NewTicker(every)
      defer t.Stop()
      for {
          select {
          case <-t.C:
              now := time.Now()
              c.mu.Lock()
              for k, e := range c.items {
                  if now.After(e.expires) {
                      delete(c.items, k)     // deleting while ranging is legal
                  }
              }
              c.mu.Unlock()
          case <-c.done:
              return
          }
      }
  }

  // Close is safe to call more than once, which a Once gives us free.
  func (c *Cache) Close() {
      c.stopper.Do(func() { close(c.done) })
  }

  func (c *Cache) Stats() (hits, misses int64) {
      return c.hits.Load(), c.misses.Load()
  }

  func main() {
      c := NewCache(50 * time.Millisecond)
      defer c.Close()

      c.Set("a", "alpha", 200*time.Millisecond)
      c.Set("b", "bravo", 40*time.Millisecond)

      // Hammer it from several goroutines to prove the locking holds.
      var wg sync.WaitGroup
      for i := 0; i < 50; i++ {
          wg.Go(func() {
              c.Get("a")
              c.Set(fmt.Sprintf("k%d", i), "v", time.Second)
          })
      }
      wg.Wait()

      time.Sleep(100 * time.Millisecond)   // let b expire and be cleaned

      if _, ok := c.Get("b"); !ok {
          fmt.Println("b expired, as expected")
      }
      h, m := c.Stats()
      fmt.Printf("hits=%d misses=%d\n", h, m)

      c.Close()
      c.Close()                            // safe, thanks to Once
      fmt.Println("closed twice without panicking")
  }
  ```

  Run it with `-race`. It must be clean.

  Then break it three ways, one at a time, and run with `-race` after each.
  Change `hits` from an `atomic.Int64` to a plain `int64` and watch the
  detector fire. Change `RLock` to no lock at all in `Get`, same result.
  Finally remove the `sync.Once` from `Close` and call it twice: the second
  `close(c.done)` panics, which is the "closing twice" rule from the last
  chapter meeting a real program.

  One design note worth absorbing: `Get` releases the read lock *before*
  checking expiry and updating counters. The counters are atomic so they need
  no lock, and holding a lock while doing work you do not have to is the habit
  that turns a cache into a bottleneck.
]

#tryit[
  Add a `Len() int` method and decide which lock it takes. Then add a maximum
  size, evicting the entry closest to expiry when the cache is full. That
  second one is harder than it looks: finding the minimum requires scanning
  under the write lock, which makes `Set` O(n). Think about what you would
  need to make it fast, and notice you are now designing an LRU, which is why
  real cache libraries exist.
]
