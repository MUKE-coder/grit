#import "../template.typ": *

= Interfaces, and why Go's point the other way

If you take one idea from this book, take this one. Go's interfaces look like
TypeScript's and work backwards from them, and once that clicks a great deal of
Go code stops looking odd and starts looking inevitable.


#bridge(title: "Things you have not met yet")[
  Two examples look ahead: context.Context is chapter 16 and the test is
  chapter 18. Interfaces are what make both of them work, which is why
  they are here rather than there.
]

== The mechanical part

An interface is a set of method signatures.

```go
type Shape interface {
    Area() float64
    Perimeter() float64
}
```

Any type with those two methods *is* a `Shape`. You do not say so anywhere.

```go
type Rect struct{ W, H float64 }

func (r Rect) Area() float64      { return r.W * r.H }
func (r Rect) Perimeter() float64 { return 2 * (r.W + r.H) }

// Rect is now a Shape. Nothing was declared. Nothing imports Shape.
var s Shape = Rect{W: 3, H: 4}
```

This is *structural typing*, and it is called implicit interface satisfaction.

#bridge[
  TypeScript interfaces are structural too, so this much is familiar: an object
  with the right shape satisfies the interface. The difference is where the
  interface *lives* and who declares it, which is the next section and is the
  part that actually matters.
]

== The direction: interfaces belong to the consumer

In Java or C\#, the implementer declares what it implements: `class Rect
implements Shape`. The interface is defined near the thing that provides it,
and the provider must import it.

In Go, the convention is the reverse. *The interface is declared by the package
that consumes it, and lists only what that consumer needs.*

Consider a function that writes a report somewhere.

```go
// In package report. It does not import os, or a logging library, or a
// storage SDK. It says what it needs and nothing more.
type Writer interface {
    Write(p []byte) (int, error)
}

func Generate(w Writer, data []Row) error {
    for _, r := range data {
        _, err := fmt.Fprintf(w, "%s,%d\n", r.Name, r.Count)
        if err != nil {
            return err
        }
    }
    return nil
}
```

Now `Generate` accepts a file, a network connection, an HTTP response, an
in-memory buffer, a gzip stream, or anything else with a `Write` method. None
of those types know `report` exists. The dependency points *from* the consumer
*to* nothing at all.

#bridge[
  In TypeScript you would often import a type from the module that provides it,
  or define a shared `types.ts` that both sides import. Go's version removes
  the shared file: the consumer owns a small interface, the provider owns a
  concrete type, and neither imports the other. That is the decoupling, and it
  is why Go codebases have surprisingly few imports between their own packages.
]

== Keep them small

That `Writer` above is the real one. `io.Writer` in the standard library is
exactly those four words, and it is arguably the most successful interface ever
designed:

```go
type Writer interface {
    Write(p []byte) (int, error)
}
```

One method. Hundreds of types satisfy it. Everything composes.

The proverb is *the bigger the interface, the weaker the abstraction*. A
one-method interface is satisfied by nearly anything and is easy to fake in a
test. A ten-method interface is satisfied by one type and is a class
declaration wearing a disguise.

#gotcha[
  The beginner mistake, carried in from other languages, is to define an
  interface for every struct up front:

  ```go
  // Do not do this.
  type UserService interface {
      Create(...) error
      Update(...) error
      Delete(...) error
      List(...) ([]User, error)
      // ... nine more
  }

  type userService struct{}   // the only implementation, ever
  ```

  One implementation and one interface that mirrors it exactly is pure
  ceremony: you have doubled the code and decoupled nothing. The Go advice is
  *accept interfaces, return structs*: return the concrete type from your
  constructor, and let whichever caller needs an abstraction declare the two
  methods it actually uses.
]

== The empty interface, and `any`

An interface with no methods is satisfied by everything.

```go
var x any = 42          // any is an alias for interface{}
x = "hello"
x = User{}
```

`any` is Go's `unknown`, not its `any`. You cannot do anything with the value
until you establish what it is. The tool is a *type assertion*:

```go
s, ok := x.(string)
if !ok {
    // x was not a string
}
```

Or a type switch, which is the readable form when several types are possible:

```go
switch v := x.(type) {
case string:
    fmt.Println("string of length", len(v))
case int:
    fmt.Println("int doubled:", v*2)
case User:
    fmt.Println("user named", v.Name)
default:
    fmt.Println("something else")
}
```

#gotcha[
  The single-value assertion `s := x.(string)` *panics* if the type is wrong.
  The two-value form `s, ok := x.(string)` never panics. Use the two-value form
  unless you genuinely want to crash, which in a server you almost never do.
]

Reach for `any` rarely. It throws away everything the type system knows, and
most uses of it in a new codebase are better served by an interface with one
method or, since Go 1.18, by generics.

== The nil interface trap

This one catches experienced people and is worth the page.

An interface value has two parts internally: a *type* and a *value*. It is
`nil` only when both are nil.

```go
type MyError struct{}
func (e *MyError) Error() string { return "boom" }

func doWork() error {
    var e *MyError = nil    // a nil pointer
    return e                // wrapped in an error interface
}

err := doWork()
fmt.Println(err == nil)     // false!
```

The returned interface holds type `*MyError` and value `nil`. The type is not
nil, so the interface is not nil, so `err != nil` is true and your error
handling runs for an error that is not there.

#bridge[
  There is no JavaScript equivalent because there is no such two-part value.
  The practical rule is the one from chapter four, and now you know why it
  exists: *declare the return type as `error` and return a literal `nil`*,
  never a typed nil variable.

  ```go
  func doWork() error {
      if somethingWrong {
          return &MyError{}
      }
      return nil          // literal nil, not a typed nil
  }
  ```
]

== Interfaces you will meet constantly

Four from the standard library carry an enormous amount of Go code.

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 6pt,
  [`io.Reader`], [`Read(p []byte) (int, error)`. Anything you can read bytes from],
  [`io.Writer`], [`Write(p []byte) (int, error)`. Anything you can write bytes to],
  [`error`], [`Error() string`. You met it in chapter four],
  [`fmt.Stringer`], [`String() string`. How a type prints itself],
)

`Stringer` is a quick win. Implement it and `fmt.Println` uses it:

```go
func (u User) String() string {
    return fmt.Sprintf("%s <%s>", u.Name, u.Email)
}

fmt.Println(u)      // Ada <ada@example.com>
```

#bridge[
  That is `toString()`, with one important difference: nothing is called
  implicitly during string concatenation, because Go has no implicit
  conversion. `"x" + u` does not compile. Only the `fmt` family looks for
  `String()`.
]

== Why this matters for testing

The payoff arrives when you test. A function that takes an interface can be
handed a fake with no mocking library and no dependency injection framework.

```go
// In the package that needs it, listing only what it uses.
type Mailer interface {
    Send(to, subject, body string) error
}

func Register(u User, m Mailer) error {
    // ...
    return m.Send(u.Email, "Welcome", "Hello "+u.Name)
}
```

The test:

```go
type fakeMailer struct {
    sent []string
}

func (f *fakeMailer) Send(to, subject, body string) error {
    f.sent = append(f.sent, to)
    return nil
}

func TestRegisterSendsWelcome(t *testing.T) {
    m := &fakeMailer{}
    if err := Register(User{Email: "ada@example.com"}, m); err != nil {
        t.Fatal(err)
    }
    if len(m.sent) != 1 || m.sent[0] != "ada@example.com" {
        t.Errorf("got %v, want one mail to ada@example.com", m.sent)
    }
}
```

Six lines of fake, no library, no configuration. That is the whole argument for
consumer-declared interfaces, and it only works because `Mailer` has one method
rather than the fifteen a real mail client has.

#project(title: "Project 9: a notifier with three backends")[
  Build something that uses an interface the way Go intends, then prove the
  decoupling by adding a backend without touching the code that uses it.

  ```go
  package main

  import (
      "fmt"
      "io"
      "os"
      "strings"
  )

  // Notifier is declared here, by the consumer, and has one method.
  type Notifier interface {
      Notify(message string) error
  }

  // Three unrelated implementations. None of them mention Notifier.

  type ConsoleNotifier struct{}

  func (ConsoleNotifier) Notify(msg string) error {
      _, err := fmt.Println("[console]", msg)
      return err
  }

  type WriterNotifier struct{ W io.Writer }

  func (w WriterNotifier) Notify(msg string) error {
      _, err := fmt.Fprintf(w.W, "[file] %s\n", msg)
      return err
  }

  type CollectingNotifier struct{ Messages []string }

  func (c *CollectingNotifier) Notify(msg string) error {
      c.Messages = append(c.Messages, msg)
      return nil
  }

  // The consumer. It never learns which backend it has.
  func alertAll(ns []Notifier, msg string) error {
      for _, n := range ns {
          if err := n.Notify(msg); err != nil {
              return fmt.Errorf("notifying: %w", err)
          }
      }
      return nil
  }

  func main() {
      var buf strings.Builder
      collector := &CollectingNotifier{}

      notifiers := []Notifier{
          ConsoleNotifier{},
          WriterNotifier{W: &buf},
          WriterNotifier{W: os.Stdout},
          collector,
      }

      if err := alertAll(notifiers, "disk almost full"); err != nil {
          fmt.Fprintln(os.Stderr, err)
          os.Exit(1)
      }

      fmt.Print("buffer contains: ", buf.String())
      fmt.Println("collector saw:", collector.Messages)
  }
  ```

  Notice three things. `WriterNotifier` is constructed twice with completely
  different destinations, a string builder and standard output, because both
  satisfy `io.Writer`. `CollectingNotifier` uses a pointer receiver because it
  mutates, so the slice holds `collector` rather than `CollectingNotifier{}`,
  and putting the value in instead is a compile error worth seeing. And
  `alertAll` would work unchanged against an email sender, a Slack client or a
  test fake.

  Now add a fourth: a `PrefixNotifier` that wraps another `Notifier` and adds a
  prefix before delegating. That is the decorator pattern, in Go, in eight
  lines, with no base class:

  ```go
  type PrefixNotifier struct {
      Prefix string
      Inner  Notifier
  }

  func (p PrefixNotifier) Notify(msg string) error {
      return p.Inner.Notify(p.Prefix + msg)
  }
  ```
]

#tryit[
  Give `CollectingNotifier` a `String() string` method so it satisfies
  `fmt.Stringer`, then print the collector directly with `fmt.Println(collector)`
  and watch the output change. Then write a test for `alertAll` that proves it
  stops and returns a wrapped error when one notifier fails, using a fake that
  always errors. That test is four lines and needs no library, which is the
  point of the chapter.
]

#assignment(name: "Interfaces - Assignment One")[
  Declare an interface where it is used, and feel the difference.

  In `part2-ideas/10-interfaces`, write a `notify` package that sends a
  message.

  *One.* The package defines what it needs, not what anybody provides:

  ```go
  type Sender interface {
      Send(ctx context.Context, to, body string) error
  }
  ```

  *Two.* Two implementations in other packages: one printing to stdout, one
  collecting into a slice. Neither imports `notify` to say it implements
  anything.

  *Three.* A test that uses the collecting one and asserts on what was sent.
  No mocking library, no interface file, no registration.

  *Four.* Add a third method to the interface and watch both implementations
  break at compile time. Then decide whether the method belonged there, and say
  so in `NOTES.md`.

  The website's *Interfaces* and *Swapping Implementations* cover the mechanics
  if this is new.]
