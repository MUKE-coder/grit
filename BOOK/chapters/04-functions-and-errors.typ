#import "../template.typ": *

= Functions, multiple returns, and the end of try/catch

This is the chapter that decides whether you get on with Go. Error handling is
where the language is most unlike JavaScript, it is the thing newcomers
complain about loudest, and the complaint is half right. Let us do the
mechanics first and the argument second.

== Functions

```go
func add(a int, b int) int {
    return a + b
}

func add(a, b int) int {   // same thing, shared type
    return a + b
}
```

The type comes *after* the name, which reads backwards for about a day and
then stops mattering. The reason is that it makes complicated declarations read
left to right instead of spiralling outwards, which anyone who has read a C
function pointer declaration will appreciate.

Functions are values, so the things you do in JavaScript still work:

```go
double := func(n int) int { return n * 2 }

func apply(nums []int, f func(int) int) []int {
    out := make([]int, len(nums))
    for i, n := range nums {
        out[i] = f(n)
    }
    return out
}

apply([]int{1, 2, 3}, double)  // [2 4 6]
```

#bridge[
  That `apply` is `Array.prototype.map`, written by hand. Go had no generics
  until 1.18 and still has no `map`, `filter` or `reduce` in the standard
  library for slices. You write the loop. Go's argument is that a three-line
  loop needs no explanation and a chain of five callbacks sometimes does. You
  are allowed to find this annoying.
]

== Two return values, and why

A Go function can return several values, and this is not a curiosity: it is the
foundation of the error system.

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("divide by zero")
    }
    return a / b, nil
}
```

The caller gets both:

```go
result, err := divide(10, 2)
if err != nil {
    fmt.Println("failed:", err)
    return
}
fmt.Println(result)
```

The convention is absolute, across the standard library and every package you
will use: *the error is the last return value, and you check it immediately*.

#bridge[
  JavaScript has one return value, so a function that can fail either throws or
  returns something like `{ data, error }`. Go makes the second shape the
  language default. If you have written Go-style error handling in TypeScript
  with a `Result` type, or used Rust, you already have the mental model.
]

== `error` is just an interface

There is nothing magic about errors. `error` is an interface with one method:

```go
type error interface {
    Error() string
}
```

Anything with an `Error() string` method is an error. That means you can make
your own, carrying whatever fields you want:

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return e.Field + ": " + e.Message
}

func validate(age int) error {
    if age < 0 {
        return ValidationError{Field: "age", Message: "must not be negative"}
    }
    return nil
}
```

== Wrapping: adding context without losing the original

A bare error from deep in a call stack is useless on its own. `fmt.Errorf` with
the `%w` verb wraps one error inside another, keeping both.

```go
func loadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("loading config from %s: %w", path, err)
    }
    ...
}
```

The message becomes a trail that reads from the outside in:

```
loading config from /etc/app.json: open /etc/app.json: no such file or directory
```

Use `%w`, not `%v`. `%v` formats the error as text and throws the original
away, and then both of the checks below stop working.

Two functions read a wrapped chain. `errors.Is` asks "is this error, anywhere
in the chain, that specific one":

```go
if errors.Is(err, os.ErrNotExist) {
    // the file was missing, however deeply that was wrapped
}
```

`errors.As` asks "is there an error of this type in the chain", and gives it to
you with its fields intact:

```go
var ve ValidationError
if errors.As(err, &ve) {
    fmt.Println("bad field:", ve.Field)
}
```

#gotcha[
  `err != nil` is the only check that matters, and comparing an error to `nil`
  has one notorious trap. A function returning a *concrete* error type rather
  than the `error` interface can produce an error value that is not `nil` even
  when it holds nothing:

  ```go
  func bad() error {
      var e *ValidationError   // nil pointer
      return e                 // not nil! it is a non-nil interface holding a nil pointer
  }
  ```

  The fix is a rule: *always declare the return type as `error`, and return a
  literal `nil`, never a typed nil variable.* This catches experienced people
  and it is worth reading twice.
]

== The argument about verbosity

Here is the complaint, stated fairly. This:

#compare(
  "async function setup() {
  try {
    const cfg = await loadConfig();
    const db = await connect(cfg.dsn);
    await migrate(db);
    return db;
  } catch (e) {
    throw new Error(`setup failed: ${e.message}`);
  }
}",
  "func setup() (*DB, error) {
    cfg, err := loadConfig()
    if err != nil {
        return nil, fmt.Errorf(\"setup: %w\", err)
    }
    db, err := connect(cfg.DSN)
    if err != nil {
        return nil, fmt.Errorf(\"setup: %w\", err)
    }
    if err := migrate(db); err != nil {
        return nil, fmt.Errorf(\"setup: %w\", err)
    }
    return db, nil
}",
)

The Go is twice as long. That is not a matter of opinion.

What you get for it is that every failure point is visible at the point it can
fail. In the JavaScript, which of the three awaited calls threw? You cannot
tell from the code, and the catch block treats all three identically whether or
not that is what you wanted. In the Go, each failure is a separate decision,
and you could handle the migration failure differently from the connection
failure by changing one block.

The honest summary: `try`/`catch` is better when you want to treat a whole
region of code as one unit that either works or does not. Explicit returns are
better when failures deserve individual handling. Most real backend code is the
second kind more often than people expect, and Go bet the language on it.

#bridge(title: "What you lose, and should know about")[
  In JavaScript an uncaught exception unwinds to the nearest handler, so a
  forgotten `catch` still surfaces somewhere. In Go a forgotten error check
  silently continues with a zero value. Nothing warns you.

  This is a real hazard and the tooling answers it: `errcheck`, bundled into
  `golangci-lint`, fails the build on an unchecked error. Install it early. It
  is the one linter that is not optional.
]

== Named returns and `defer`

A function can name its return values, which is occasionally useful and
frequently overused:

```go
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return          // a "naked" return, sending back x and y
}
```

Use naked returns only in short functions. In a long one, a bare `return`
forces the reader to scroll up to learn what is being returned.

Named returns do earn their place with `defer`, which schedules a call to run
when the surrounding function exits, no matter how it exits:

```go
func process(path string) (err error) {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()     // runs on every return path below

    // ... many returns, all of which close the file
    return nil
}
```

`defer f.Close()` immediately after a successful open is one of the most common
patterns in Go, and it is how the language avoids the cleanup bugs that
`try`/`finally` exists to solve.

#gotcha[
  Deferred calls run in last-in, first-out order, and their *arguments* are
  evaluated immediately, not at exit:

  ```go
  for i := 0; i < 3; i++ {
      defer fmt.Println(i)   // prints 2, 1, 0
  }
  ```

  And a `defer` inside a loop does not run at the end of each iteration, it
  runs when the *function* returns. Opening a thousand files in a loop with
  `defer f.Close()` keeps all thousand open until the function ends. Put the
  body in its own function if you need per-iteration cleanup.
]

#project(title: "Project 3: a word counter that handles failure")[
  Write a program that takes a filename, counts the words in it, and reports
  errors properly rather than crashing.

  ```go
  package main

  import (
      "fmt"
      "os"
      "strings"
  )

  func countWords(path string) (int, error) {
      data, err := os.ReadFile(path)
      if err != nil {
          return 0, fmt.Errorf("counting words in %s: %w", path, err)
      }
      return len(strings.Fields(string(data))), nil
  }

  func main() {
      if len(os.Args) < 2 {
          fmt.Fprintln(os.Stderr, "usage: wc <file>")
          os.Exit(1)
      }
      n, err := countWords(os.Args[1])
      if err != nil {
          fmt.Fprintln(os.Stderr, "error:", err)
          os.Exit(1)
      }
      fmt.Printf("%d words\n", n)
  }
  ```

  Run it on a real file, then on one that does not exist, and read the wrapped
  message. Notice `os.Exit(1)` and writing to `os.Stderr`: a command line tool
  that reports failure on stdout with exit code zero is a tool that breaks
  every script that uses it.

  `strings.Fields` splits on any run of whitespace, which is the right
  behaviour here and is worth remembering: it is not the same as
  `strings.Split(s, " ")`, which gives empty strings between consecutive
  spaces.
]

#tryit[
  Add a check that reports a clearer message when the path is a directory
  rather than a file. `os.ReadFile` returns an error you can inspect; try
  `errors.Is(err, syscall.EISDIR)` on Unix, or simply call `os.Stat` first and
  check `info.IsDir()`, which is portable and is the version you should
  probably ship.
]
