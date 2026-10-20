#import "../template.typ": *

= defer, panic and recover

Three keywords that look like `finally`, `throw` and `catch` and are not. The
resemblance is close enough to mislead, which is why they get a chapter.

== defer: cleanup that cannot be forgotten

`defer` schedules a function call to run when the *surrounding function*
returns, by any path.

```go
func readConfig(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()     // runs on every return below, and on a panic

    // twenty lines with six return statements
    return nil
}
```

Write the `defer` immediately after the thing that needs cleaning up. That
adjacency is the point: open and close are one line apart, so a reviewer can
see the pair without scrolling.

#bridge[
  `try { ... } finally { f.close() }` does the same job by wrapping, which
  indents the body and separates the close from the open. `defer` keeps them
  together and costs no indentation. It also composes: ten resources means ten
  `defer` lines, not ten nested blocks.
]

== The three rules of defer

*One: arguments are evaluated immediately, the call happens later.*

```go
func main() {
    msg := "first"
    defer fmt.Println(msg)   // prints "first", not "second"
    msg = "second"
}
```

`msg` was read when the `defer` statement ran. To defer the *evaluation*, wrap
it in a closure:

```go
defer func() { fmt.Println(msg) }()   // prints "second"
```

*Two: deferred calls run last in, first out.*

```go
for i := 0; i < 3; i++ {
    defer fmt.Println(i)
}
// prints 2, then 1, then 0
```

This is right for cleanup: resources are released in the reverse of the order
they were acquired, like unwinding a stack.

*Three: defer is function-scoped, not block-scoped.*

#gotcha[
  This is the one that causes real bugs:

  ```go
  func processAll(paths []string) error {
      for _, p := range paths {
          f, err := os.Open(p)
          if err != nil {
              return err
          }
          defer f.Close()     // does NOT run at the end of this iteration
          // ... use f
      }
      return nil
  }
  ```

  All the files stay open until `processAll` returns. With a thousand paths you
  run out of file descriptors.

  The fix is to give each iteration its own function, so there is a return for
  the `defer` to fire on:

  ```go
  for _, p := range paths {
      if err := processOne(p); err != nil {
          return err
      }
  }

  func processOne(p string) error {
      f, err := os.Open(p)
      if err != nil {
          return err
      }
      defer f.Close()     // runs when processOne returns
      // ... use f
      return nil
  }
  ```

  Extracting the loop body into a function is the standard answer, and it
  usually improves the code for other reasons too.
]

== defer can change the return value

A deferred closure can modify a *named* return value after the `return`
statement has chosen it. This is the one legitimate use of named returns.

```go
func doWork() (err error) {
    defer func() {
        if err != nil {
            err = fmt.Errorf("doWork: %w", err)
        }
    }()

    return doSomethingThatFails()    // gets wrapped on the way out
}
```

Used sparingly this is tidy. Used often it makes a function's return value hard
to trace, so prefer explicit wrapping at each return unless you have many of
them.

== defer and errors from Close

A deferred `Close` discards its error, which is usually fine for reading and
sometimes wrong for writing: a failed `Close` on a file you wrote may mean the
data never reached the disk.

```go
func writeReport(path string, data []byte) (err error) {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer func() {
        if cerr := f.Close(); cerr != nil && err == nil {
            err = fmt.Errorf("closing %s: %w", path, cerr)
        }
    }()

    _, err = f.Write(data)
    return err
}
```

That is verbose and correct. For a read-only file, `defer f.Close()` is
accepted practice and nobody will object.

== panic: stop, this should not have happened

`panic` unwinds the stack, running deferred functions as it goes, and if it
reaches the top of a goroutine the program crashes with a stack trace.

```go
func main() {
    defer fmt.Println("this still runs")
    panic("something impossible")
}
```

```
this still runs
panic: something impossible

goroutine 1 [running]:
main.main()
    /tmp/main.go:7 +0x65
exit status 2
```

Go panics on its own for a handful of runtime faults: nil pointer dereference,
index out of range, division by zero, sending on a closed channel, and writing
to a nil map.

#bridge[
  `throw` in JavaScript is for *expected* failures: a 404, a validation
  problem, a parse error. `panic` is not. Go's expected failures are returned
  values, and `panic` is reserved for "this program's assumptions are broken".

  If you find yourself panicking in a request handler because a user sent bad
  input, that is the JavaScript habit talking. Return an error.
]

== recover: only for a boundary

`recover` stops a panic, and it only works *inside a deferred function*.

```go
func safe() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered: %v", r)
        }
    }()

    panic("boom")
}

err := safe()
fmt.Println(err)    // recovered: boom
```

There are exactly two good reasons to use it.

*One: a server that must not die because one request was bad.* A panic in one
handler should return 500 for that request, not take down the process serving
everybody else.

```go
func Recover(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                log.Printf("panic on %s %s: %v\n%s",
                    r.Method, r.URL.Path, rec, debug.Stack())
                http.Error(w, "internal server error",
                    http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

Log the stack. A recovered panic with no stack trace is a bug you will never
find. Every Go web framework ships this middleware, Gin included, and Grit
wires it in for you.

*Two: a library turning an internal panic into an error at its edge.* Some
parsers panic internally to unwind out of deep recursion, then recover at the
public function and return an error. The caller never sees a panic. This is
legitimate and rare.

#gotcha[
  `recover` only catches panics in *its own goroutine*. A panic in a goroutine
  you started is not caught by a recover in the function that started it:

  ```go
  func main() {
      defer func() { recover() }()     // does not help
      go func() {
          panic("in a goroutine")      // kills the whole program
      }()
      time.Sleep(time.Second)
  }
  ```

  Every goroutine that could panic needs its own recover. This is the single
  most common way a Go server dies in production: a background worker with no
  recover, panicking on a nil map six hours after deployment.
]

#project(title: "Project 11: a worker that survives bad input")[
  Build a small job runner where one bad job must not stop the rest, and where
  cleanup is guaranteed. This is the panic boundary in miniature.

  ```go
  package main

  import (
      "errors"
      "fmt"
      "strings"
  )

  type Job struct {
      Name string
      Run  func() error
  }

  // runOne is the boundary. Whatever the job does, this returns an error
  // rather than letting a panic escape.
  func runOne(j Job) (err error) {
      defer func() {
          if r := recover(); r != nil {
              err = fmt.Errorf("job %q panicked: %v", j.Name, r)
          }
      }()
      defer fmt.Printf("  [%s] finished\n", j.Name)   // runs even on panic

      if j.Run == nil {
          return errors.New("job has no work")
      }
      return j.Run()
  }

  func runAll(jobs []Job) (failed int) {
      for _, j := range jobs {
          fmt.Printf("running %s\n", j.Name)
          if err := runOne(j); err != nil {
              fmt.Println("  error:", err)
              failed++
          }
      }
      return failed
  }

  func main() {
      jobs := []Job{
          {Name: "ok", Run: func() error { return nil }},
          {Name: "fails", Run: func() error {
              return errors.New("disk full")
          }},
          {Name: "panics", Run: func() error {
              var m map[string]int
              m["boom"] = 1              // panic: assignment to nil map
              return nil
          }},
          {Name: "index", Run: func() error {
              s := []int{1, 2, 3}
              fmt.Println(s[10])         // panic: index out of range
              return nil
          }},
          {Name: "no work"},
          {Name: "still runs", Run: func() error {
              fmt.Println("  the runner survived four failures")
              return nil
          }},
      }

      failed := runAll(jobs)
      fmt.Printf("\n%d of %d jobs failed\n", failed, len(jobs))
      _ = strings.TrimSpace("")
  }
  ```

  Three things to notice when you run it.

  The `finished` line prints for every job including the ones that panicked,
  because deferred calls run while the stack unwinds. That is what makes
  `defer` safe for cleanup in a way that a line at the end of the function is
  not.

  The two deferred functions run in reverse order, so `finished` prints before
  the recover's error is assembled.

  And the program exits normally with a count, rather than dying on job three.
  Delete the `recover` block and run it again: you get a stack trace and jobs
  four, five and six never run. That difference is the reason every server has
  this middleware.
]

#tryit[
  Move the `recover` from `runOne` up into `runAll`, wrapping the whole loop,
  and run it again. You will catch the first panic and then the loop stops,
  because the panic unwound past the loop. That experiment shows exactly what
  "a boundary" means: it has to be inside the unit of work you want to isolate,
  not outside it.
]
