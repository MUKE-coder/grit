#import "../template.typ": *

= Testing: the runner is already installed

Go ships a test runner, an assertion-free testing philosophy, a benchmark
harness, a coverage tool, a fuzzer and a race detector. There is nothing to
install and nothing to configure. After a JavaScript project's testing setup,
this takes some adjusting to, mostly because you keep looking for the config
file.

== The rules

A test file ends in `_test.go`, sits beside the code it tests, and contains
functions named `TestXxx` taking `*testing.T`.

```go
// math.go
package math

func Add(a, b int) int { return a + b }
```

```go
// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    got := Add(2, 3)
    if got != 5 {
        t.Errorf("Add(2, 3) = %d, want 5", got)
    }
}
```

```bash
go test ./...
```

```
ok      github.com/you/myapp/math    0.002s
```

#bridge[
  No `describe`, no `it`, no `expect`. A test is a function, and it fails by
  calling `t.Errorf`. There is no assertion library in the standard library,
  and this is deliberate: the Go team's position is that `if got != want` is
  clearer than learning which of `toBe`, `toEqual` and `toStrictEqual` you
  wanted.

  You will still meet `testify` in real codebases, which adds
  `assert.Equal(t, want, got)`. It is popular and it is fine. Start without it.
]

== Error messages are the whole skill

A failing test must tell you what happened without opening the source. The
convention is `got, want`:

```go
t.Errorf("Add(2, 3) = %d, want 5", got)
```

```
    math_test.go:8: Add(2, 3) = 6, want 5
```

That line names the input, the result and the expectation. Compare it to
`t.Error("failed")`, which tells you nothing you did not already know.

== Errorf against Fatalf

`t.Errorf` records a failure and *keeps going*. `t.Fatalf` records it and stops
the test function immediately.

```go
u, err := Load(1)
if err != nil {
    t.Fatalf("Load(1): unexpected error: %v", err)   // stop: u is unusable
}
if u.Name != "Ada" {
    t.Errorf("Name = %q, want Ada", u.Name)          // keep checking
}
```

Use `Fatalf` when continuing would panic or produce noise. Use `Errorf` when
you want to see every failure in one run.

== Table-driven tests

This is the dominant Go testing idiom. One test function, a slice of cases.

```go
func TestSlugify(t *testing.T) {
    tests := []struct {
        name string
        in   string
        want string
    }{
        {"simple", "Hello World", "hello-world"},
        {"punctuation", "Hello, World!", "hello-world"},
        {"already a slug", "hello-world", "hello-world"},
        {"leading space", "  padded  ", "padded"},
        {"empty", "", ""},
        {"unicode", "Café Noir", "cafe-noir"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Slugify(tt.in); got != tt.want {
                t.Errorf("Slugify(%q) = %q, want %q", tt.in, got, tt.want)
            }
        })
    }
}
```

`t.Run` creates a subtest with its own name, so output identifies the case:

```
--- FAIL: TestSlugify (0.00s)
    --- FAIL: TestSlugify/unicode (0.00s)
        slug_test.go:21: Slugify("Café Noir") =
            "caf-noir", want "cafe-noir"
```

You can run one case directly:

```bash
go test -run 'TestSlugify/unicode' ./...
```

#bridge[
  This replaces `describe.each` and `test.each`. The advantage over both is
  that the cases are an ordinary slice of an ordinary struct, so adding a field
  to the table is a normal code change and the compiler checks every row.
]

== Comparing things that are not numbers

`==` works on structs of comparable fields. For slices, maps and nested
structs, use `reflect.DeepEqual` or, much better, `go-cmp`:

```go
import "github.com/google/go-cmp/cmp"

if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf("Load() mismatch (-want +got):\n%s", diff)
}
```

```
    user_test.go:31: Load() mismatch (-want +got):
          user.User{
                ID:    1,
        -       Name:  "Ada",
        +       Name:  "ada",
                Email: "ada@example.com",
          }
```

A readable diff instead of `false`. This is the one testing dependency worth
adding on day one.

== Helpers and t.Helper

A helper that reports failures should call `t.Helper()` so the line number
points at the *caller*:

```go
func mustLoad(t *testing.T, id int) *User {
    t.Helper()                       // failures blame the caller's line
    u, err := Load(id)
    if err != nil {
        t.Fatalf("Load(%d): %v", id, err)
    }
    return u
}
```

Without it, every failure points at the line inside `mustLoad`, which is
useless when ten tests call it.

== Setup and cleanup

`t.Cleanup` registers work to run when the test finishes, including on failure.

```go
func newTestDB(t *testing.T) *sql.DB {
    t.Helper()
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatal(err)
    }
    t.Cleanup(func() { db.Close() })    // runs even if the test fails
    return db
}
```

For package-wide setup, `TestMain`:

```go
func TestMain(m *testing.M) {
    // before all tests
    code := m.Run()
    // after all tests
    os.Exit(code)
}
```

#bridge[
  `t.Cleanup` is `afterEach`, scoped to one test and registered where the
  resource is created rather than in a separate block. `TestMain` is
  `globalSetup` and `globalTeardown`. There is no `beforeEach`: you call a
  helper at the top of each test, which is more typing and means a reader can
  see what a test depends on without scrolling.
]

== Parallel tests

```go
func TestSomething(t *testing.T) {
    t.Parallel()
    // ...
}
```

Tests in different packages already run in parallel. `t.Parallel()` marks tests
*within* a package as safe to interleave.

#gotcha[
  A parallel test must not share mutable state with another test. Writing to
  the same database table, the same temp file or the same package-level
  variable will fail intermittently, which is worse than failing reliably.

  Combine `t.Parallel()` with `-race` so that a shared-state mistake is
  reported rather than merely flaky.
]

== Testing errors

Test the *kind*, not the message. Message text is not a contract.

```go
// Brittle: breaks when someone improves the wording.
if err.Error() != "user not found" { ... }

// Correct.
if !errors.Is(err, ErrNotFound) {
    t.Errorf("got %v, want ErrNotFound", err)
}
```

And always test that an error *is* returned when you expect one:

```go
{"missing email", "", "", true},   // a wantErr field in the table
...
_, err := Create(tt.name, tt.email)
if (err != nil) != tt.wantErr {
    t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
}
```

== Coverage

```bash
go test -cover ./...
go test -coverprofile=cover.out ./... && go tool cover -html=cover.out
```

The second opens a browser showing every line in green or red. It is the
fastest way to find the branch nobody tested.

#gotcha[
  Coverage measures lines executed, not behaviour verified. A test that calls
  every function and asserts nothing reports high coverage and proves nothing.
  Use it to find what you *missed*, never as a target to hit.
]

== Benchmarks

```go
func BenchmarkSlugify(b *testing.B) {
    for b.Loop() {
        Slugify("Hello, World! This is a title")
    }
}
```

```bash
go test -bench=. -benchmem ./...
```

```
BenchmarkSlugify-8    2419371    489.2 ns/op    144 B/op    3 allocs/op
```

Operations per second, nanoseconds each, bytes allocated and allocation count.
`b.Loop()` replaced the older `for i := 0; i < b.N; i++` in Go 1.24 and is
harder to get wrong, because it stops the compiler optimising your work away.

#project(title: "Project 17: test the user store properly")[
  Take the store from chapter eleven and give it a real test suite: table
  driven, testing error kinds rather than messages, with a helper and cleanup.

  ```go
  package store

  import (
      "errors"
      "testing"
  )

  func newStore(t *testing.T) *Store {
      t.Helper()
      s := NewStore()
      t.Cleanup(func() {
          // Nothing to close here, but the hook is where a real test would
          // drop a temp database. Wire it now so it is not retrofitted
          // later.
      })
      return s
  }

  func TestCreate(t *testing.T) {
      tests := []struct {
          name    string
          user    string
          email   string
          wantErr error       // nil, a sentinel, or errSomeValidation
          field   string      // for validation errors
      }{
          {"valid", "Ada", "ada@example.com", nil, ""},
          {"empty name", "", "x@example.com", errValidation, "name"},
          {"bad email", "Ada", "not-an-email", errValidation, "email"},
      }

      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              t.Parallel()
              s := newStore(t)

              u, err := s.Create(tt.user, tt.email)

              if tt.wantErr == nil {
                  if err != nil {
                      t.Fatalf("Create(%q, %q): unexpected error: %v",
                          tt.user, tt.email, err)
                  }
                  if u.Name != tt.user {
                      t.Errorf("Name = %q, want %q", u.Name, tt.user)
                  }
                  if u.ID == 0 {
                      t.Error("ID was not assigned")
                  }
                  return
              }

              if err == nil {
                  t.Fatalf("Create(%q, %q): want error, got nil",
                      tt.user, tt.email)
              }
              var ve *ValidationError
              if !errors.As(err, &ve) {
                  t.Fatalf("got %v, want a *ValidationError", err)
              }
              if ve.Field != tt.field {
                  t.Errorf("Field = %q, want %q", ve.Field, tt.field)
              }
          })
      }
  }

  func TestCreateDuplicateEmail(t *testing.T) {
      t.Parallel()
      s := newStore(t)

      if _, err := s.Create("Ada", "ada@example.com"); err != nil {
          t.Fatalf("first create: %v", err)
      }
      _, err := s.Create("Grace", "ada@example.com")

      // The identity of the failure is the contract. The wording is not.
      if !errors.Is(err, ErrConflict) {
          t.Errorf("got %v, want ErrConflict", err)
      }
  }

  func TestGetMissing(t *testing.T) {
      t.Parallel()
      s := newStore(t)

      if _, err := s.Get(999); !errors.Is(err, ErrNotFound) {
          t.Errorf("got %v, want ErrNotFound", err)
      }
  }
  ```

  Run it four ways and notice what each tells you.

  `go test ./...` gives you pass or fail. `go test -v ./...` names every
  subtest, which is how you confirm the table ran every row rather than
  silently skipping. `go test -race ./...` would catch the shared-state bug if
  `newStore` returned a package-level store instead of a fresh one, which is
  the mistake `t.Parallel()` makes dangerous. And `go test -cover ./...` tells
  you which branch of `Create` you forgot.

  The important design point is in `TestCreateDuplicateEmail`: it asserts
  `errors.Is(err, ErrConflict)` and never looks at the message. Change the
  wording of that error and the test still passes, which is correct, because
  the wording was never the promise.
]

#tryit[
  Add a benchmark for `Create` and run it with `-benchmem`. Then change the
  duplicate-email check from a linear scan to a second map keyed by email, and
  run the benchmark again with a thousand existing users. The allocation count
  and the nanoseconds per operation will both move, and seeing a design change
  show up as a number is the habit this chapter is really teaching.
]
