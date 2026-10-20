#import "../template.typ": *

= Designing with errors: what a package should return

Chapter four covered the mechanics: `if err != nil`, wrapping with `%w`,
`errors.Is` and `errors.As`. This chapter is about the decision those mechanics
leave you with, which is harder and matters more: *when someone calls your
function and it fails, what should they be able to do about it?*

Get this right and the caller writes three lines. Get it wrong and they end up
matching on error message strings, which is the clearest sign of an API that
did not think about failure.

== Three kinds of error, and when to use each

*Opaque.* The caller can log it and give up. Nothing else.

```go
return fmt.Errorf("parsing config: %w", err)
```

This is the default and it is correct far more often than people expect. If
there is nothing useful a caller could do differently, do not invent a way for
them to try.

*Sentinel.* A specific, exported value the caller compares against.

```go
var ErrNotFound = errors.New("not found")

func (s *Store) Get(id int) (*User, error) {
    ...
    return nil, ErrNotFound
}
```

The caller:

```go
u, err := store.Get(id)
if errors.Is(err, ErrNotFound) {
    http.Error(w, "no such user", http.StatusNotFound)
    return
}
```

Use a sentinel when the condition is a *fact with no detail*: not found, already
exists, timed out. The name is the entire payload.

*Typed.* A struct implementing `error`, carrying fields.

```go
type ValidationError struct {
    Field  string
    Reason string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Reason)
}
```

The caller:

```go
var ve *ValidationError
if errors.As(err, &ve) {
    respond(w, 422, map[string]string{ve.Field: ve.Reason})
}
```

Use a typed error when the caller needs *data* from the failure, not just its
identity. A validation failure needs the field name. A rate limit needs the
retry-after. A conflict needs the version that is current.

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 6pt,
  [*Opaque*], [Caller logs and gives up. The default],
  [*Sentinel*], [Caller branches on identity. `errors.Is`],
  [*Typed*], [Caller needs fields out of it. `errors.As`],
)

#gotcha[
  A sentinel or an exported error type is *part of your public API*. Renaming
  `ErrNotFound`, or changing which conditions return it, breaks callers exactly
  like changing a function signature. Declare them deliberately and few.

  The corollary: do not export an error type just because you made one. If
  nothing outside the package uses `errors.As` on it, keep it unexported and
  return it as a plain `error`.
]

== Wrap, but not at every level

`%w` preserves the chain, and the temptation is to wrap at every function
boundary. Do not. The result is a message like:

```
handling request: getting user: loading user: querying users:
executing query:
scanning row: sql: no rows in result set
```

Every layer added a word and none added information.

The rule worth following: *wrap when you add context the caller does not
already have*. The database layer knows the query. The service layer knows the
user ID. The handler knows the request. Each of those is worth one wrap. "I am
a function and I called another function" is not context.

```go
// Good: adds the thing the caller cannot see.
return fmt.Errorf("loading user %d: %w", id, err)

// Noise: the caller already knows they called LoadUser.
return fmt.Errorf("LoadUser: %w", err)
```

#bridge[
  This is the same discipline as a JavaScript stack trace, with one difference:
  Go errors have no stack trace by default, so the chain of wraps *is* your
  trace. That is why each one should earn its place. Some teams add stack
  traces with a library; the standard library deliberately does not, on the
  grounds that a well-wrapped error is more readable than sixty frames.
]

== Error strings are a contract, lowercase and unpunctuated

Convention, followed throughout the standard library:

```go
errors.New("connection refused")     // good
errors.New("Connection refused.")    // wrong: capital and full stop
```

The reason is composition. Errors get wrapped into sentences, and a capital in
the middle reads badly:

```
loading config: Connection refused.: retrying
```

Also: do not end with a newline, and do not include the word "error". The
reader knows.

== Do not return both a value and an error

Pick one. A function that returns a partially filled struct *and* an error
invites the caller to use the broken value.

```go
// Bad: what is the caller supposed to do with a half-built User?
func Load(id int) (User, error)

// Good: nil on failure, impossible to misuse.
func Load(id int) (*User, error)
```

The exceptions are real but narrow: `io.Reader` returns a byte count with an
error because a partial read is meaningful, and `strconv.Atoi` returns zero
because it has nothing better to offer. If you are not implementing something
in that shape, return nil.

== Handle an error once

The most common error-handling bug in any language is doing two things with one
failure:

```go
if err != nil {
    log.Printf("saving user: %v", err)   // handled it here
    return err                           // and again, up the stack
}
```

Now it is logged at every level, and the operator reading the logs sees the
same failure five times and has to work out whether it happened once or five
times.

Choose: *either* handle it, *or* return it. The usual split is that libraries
and services return, and the top layer, a handler or `main`, logs once and
decides what the user sees.

== Mapping errors at the boundary

This is where the three kinds pay off. One place in your program turns errors
into HTTP responses:

```go
func writeError(w http.ResponseWriter, err error) {
    var ve *ValidationError
    switch {
    case errors.Is(err, ErrNotFound):
        respond(w, 404, "not found")
    case errors.Is(err, ErrConflict):
        respond(w, 409, "already exists")
    case errors.As(err, &ve):
        respond(w, 422, map[string]string{ve.Field: ve.Reason})
    default:
        log.Printf("unhandled: %v", err)   // the only place that logs
        respond(w, 500, "internal error")
    }
}
```

Everything below this function returns errors and logs nothing. Everything
above it gets a correct status code. That is the whole architecture of error
handling in a Go service, and you will meet it again in part four, because it
is exactly what Grit generates.

#bridge[
  In an Express app this is the error middleware, `app.use((err, req, res,
  next) => ...)`, and the design problem is identical: the middleware has to
  work out what kind of failure it received. JavaScript usually does it with
  `err instanceof ValidationError` or a `code` property. Go's `errors.As` is
  the same idea with the chain-walking built in, so it still works after the
  error has been wrapped three times.
]

== panic is not an error

Go has `panic`, and it is not for expected failures. A panic unwinds the stack
and, uncaught, kills the program.

Use `panic` for *programmer error that cannot be recovered from*: an impossible
state, a nil dependency at startup, a configuration that cannot be parsed
before the server starts serving.

```go
func MustCompile(pattern string) *Regexp {
    re, err := Compile(pattern)
    if err != nil {
        panic("regexp: Compile(" + pattern + "): " + err.Error())
    }
    return re
}
```

The `Must` prefix is the convention: it means "this panics instead of returning
an error, and you should only call it where failure means the program is
broken", typically at package initialisation with a literal you wrote yourself.

Never panic for a bad request, a missing row, or a failed network call. Those
are expected, and expected failures are values.

#project(title: "Project 10: a user store with a designed error surface")[
  Build a small store that uses all three error kinds deliberately, plus the
  boundary mapping. This is a miniature of the structure every Go service ends
  up with.

  ```go
  package main

  import (
      "errors"
      "fmt"
      "strings"
  )

  // Sentinels: facts with no detail. Part of the public API.
  var (
      ErrNotFound = errors.New("user not found")
      ErrConflict = errors.New("email already registered")
  )

  // Typed: the caller needs the field name out of it.
  type ValidationError struct {
      Field  string
      Reason string
  }

  func (e *ValidationError) Error() string {
      return fmt.Sprintf("%s: %s", e.Field, e.Reason)
  }

  type User struct {
      ID    int
      Name  string
      Email string
  }

  type Store struct {
      users  map[int]*User
      nextID int
  }

  func NewStore() *Store {
      return &Store{users: map[int]*User{}, nextID: 1}
  }

  func (s *Store) Create(name, email string) (*User, error) {
      if strings.TrimSpace(name) == "" {
          return nil, &ValidationError{
              Field: "name", Reason: "must not be empty"}
      }
      if !strings.Contains(email, "@") {
          return nil, &ValidationError{
              Field: "email", Reason: "must be an email"}
      }
      for _, u := range s.users {
          if u.Email == email {
              // Wrapped, so the caller gets the identity and the detail.
              return nil, fmt.Errorf("creating user %s: %w",
                  email, ErrConflict)
          }
      }
      u := &User{ID: s.nextID, Name: name, Email: email}
      s.users[u.ID] = u
      s.nextID++
      return u, nil
  }

  func (s *Store) Get(id int) (*User, error) {
      u, ok := s.users[id]
      if !ok {
          return nil, fmt.Errorf("getting user %d: %w", id, ErrNotFound)
      }
      return u, nil
  }

  // The boundary. The only place that decides what a failure means, and
  // the
  // only place that logs.
  func describe(err error) string {
      var ve *ValidationError
      switch {
      case err == nil:
          return "200 ok"
      case errors.Is(err, ErrNotFound):
          return "404 not found"
      case errors.Is(err, ErrConflict):
          return "409 conflict"
      case errors.As(err, &ve):
          return fmt.Sprintf("422 invalid %s (%s)", ve.Field, ve.Reason)
      default:
          return "500 internal: " + err.Error()
      }
  }

  func main() {
      s := NewStore()

      _, err := s.Create("Ada", "ada@example.com")
      fmt.Println(describe(err))                  // 200 ok

      _, err = s.Create("", "x@example.com")
      fmt.Println(describe(err))                  // 422 invalid name

      _, err = s.Create("Grace", "not-an-email")
      fmt.Println(describe(err))                  // 422 invalid email

      _, err = s.Create("Someone", "ada@example.com")
      fmt.Println(describe(err))                  // 409 conflict

      _, err = s.Get(99)
      fmt.Println(describe(err))                  // 404 not found
  }
  ```

  Two things to prove to yourself. First, `errors.Is(err, ErrConflict)` returns
  true even though the error was wrapped with extra text, which is the entire
  reason to use `%w` rather than `%v`. Change that one verb to `%v` and watch
  the conflict case fall through to 500.

  Second, note that `Create` never logs. It returns. `describe` is the only
  function that decides what a failure means to the outside world, so moving
  this store behind an HTTP API means changing one function.
]

#tryit[
  Add an `Update(id int, name string) error` that can fail in three ways: the
  user is missing, the name is invalid, or the name is unchanged. Decide for
  each whether it is opaque, sentinel or typed, and be able to defend the
  choice. "Unchanged" is the interesting one: is that an error at all, or
  should `Update` quietly succeed? There is a defensible answer either way, and
  the act of choosing is the skill this chapter is teaching.
]

#assignment(name: "Errors - Assignment One")[
  Design the errors of a package before writing its happy path.

  In `part2-ideas/11-errors`, write a `store` package with `Get`, `Create` and
  `Delete` over an in-memory map.

  *One.* A sentinel `ErrNotFound`, a typed `ValidationError` carrying the field
  name, and an unexpected error wrapped with context.

  *Two.* A boundary function that maps any of them to an HTTP status and a
  code, using `errors.Is` and `errors.As` rather than string matching.

  *Three.* Check every error string against the chapter's rule: lowercase, no
  full stop, no "failed to". Fix the ones that fail.

  *Four.* Write a test for the boundary that asserts the status for each kind.
  That test is the contract your handlers will rely on in Part IV.]
