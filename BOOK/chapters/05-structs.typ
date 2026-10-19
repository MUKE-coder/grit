#import "../template.typ": *

= Structs and methods, against objects and classes

A JavaScript object is a bag of keys you can add to at any time. A Go struct is
a fixed set of named fields with fixed types, decided when you declare the type
and never changing. Everything else in this chapter follows from that one
difference.

== Declaring and creating

```go
type User struct {
    ID    int
    Name  string
    Email string
    Admin bool
}
```

There are four ways to make one, and you will see all four in real code.

```go
// 1. Field names. Use this one. It survives a field being added.
u := User{ID: 1, Name: "Ada", Email: "ada@example.com"}

// 2. Positional. Fragile: adding a field breaks every call site.
u := User{1, "Ada", "ada@example.com", false}

// 3. Zero value, then assign.
var u User
u.Name = "Ada"

// 4. new(), giving a pointer to a zero value.
u := new(User)   // *User
```

Notice what is missing from the first example: `Admin` was not set, so it is
`false`. There is no `undefined`. Every field always has a value, and a struct
you create is immediately complete.

#bridge[
  The closest JavaScript equivalent is a class with a constructor that sets
  every property. The difference is that you cannot add a field later. `u.Age =
  36` on a `User` with no `Age` field is a compile error, not a new property.
  This is the whole point: the shape of your data is checked, everywhere, for
  free.
]

== Methods: functions with a receiver

Go has no classes. You attach behaviour to a type by declaring a function with
a *receiver*, the parameter in parentheses before the name.

```go
func (u User) DisplayName() string {
    if u.Name == "" {
        return u.Email
    }
    return u.Name
}

fmt.Println(u.DisplayName())
```

`u User` is the receiver. It behaves like a parameter, because that is exactly
what it is: `u.DisplayName()` is sugar for passing `u` as the first argument.

#bridge[
  There is no `this`, and that removes an entire class of JavaScript bug. The
  receiver is a normal named parameter, bound at the call site, and it cannot
  be rebound by how the function was called. No `.bind(this)`, no arrow
  function workaround, no `const self = this`.
]

== Value receivers and pointer receivers

This is the first thing in Go that genuinely trips people up, and it is worth
getting right now rather than debugging later.

A *value receiver* gets a copy. Changes to it do not affect the original.

```go
func (u User) Rename(name string) {
    u.Name = name       // modifies the copy, then throws it away
}

u := User{Name: "Ada"}
u.Rename("Grace")
fmt.Println(u.Name)     // still "Ada"
```

A *pointer receiver* gets the address, so changes stick.

```go
func (u *User) Rename(name string) {
    u.Name = name       // modifies the original
}

u := User{Name: "Ada"}
u.Rename("Grace")
fmt.Println(u.Name)     // "Grace"
```

You do not write `(&u).Rename(...)`. Go inserts the address-of automatically
when the method needs a pointer and you have an addressable value.

#gotcha[
  The rule in practice: *if a method modifies the receiver, it must be a
  pointer receiver.* A value receiver that assigns to a field compiles happily
  and silently does nothing, which is one of the few truly quiet bugs in Go.

  The second rule: *be consistent within a type.* If any method on `User` takes
  a pointer, give them all pointers. Mixing the two is legal and makes the
  rules about which calls are addressable surface at the worst moment.
]

A practical reason to prefer pointer receivers even when not mutating: a value
receiver copies the whole struct on every call. For a struct with twenty fields
called in a loop, that is real work for no benefit.

== Composition instead of inheritance

Go has no `extends`. It has *embedding*, which looks similar and behaves
differently in one important way.

```go
type Timestamps struct {
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Post struct {
    Timestamps          // embedded: no field name
    Title string
    Body  string
}

p := Post{Title: "Hello"}
p.CreatedAt = time.Now()        // promoted, as if declared on Post
fmt.Println(p.Timestamps.CreatedAt)  // the long way, same field
```

The embedded type's fields and methods are *promoted* to the outer type. If
`Timestamps` had a method `Age()`, then `p.Age()` would work.

#bridge[
  This resembles `class Post extends Timestamps`, and the difference matters:
  there is no polymorphism here. A `Post` is *not* a `Timestamps` and cannot be
  passed where one is expected. Embedding is a convenience for forwarding
  fields and methods, not a type relationship. In Go, the thing that lets one
  type stand in for another is an interface, which is chapter ten.
]

== Struct tags: metadata for encoders

A field can carry a string of metadata, used by libraries through reflection.
You will meet this constantly, because it is how Go does JSON.

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Admin bool   `json:"-"`                  // never encoded
    Bio   string `json:"bio,omitempty"`      // dropped when empty
}
```

Without tags, `json.Marshal` uses the Go field names, so you get `{"ID":1,
"Name":"Ada"}`. With them you get the lowercase JSON your frontend expects.

#gotcha[
  Those backticks are a raw string literal, and the format inside is strict:
  `json:"name"` with no space after the colon. `json: "name"` silently does
  nothing, the tag is ignored, and the field encodes as `Name`. There is no
  error. If a field is coming out with the wrong key, check for that space
  first.
]

== Constructors are just functions

Go has no `constructor` keyword. The convention is a package-level function
named `New` or `NewThing` that returns the type.

```go
func NewUser(name, email string) (*User, error) {
    if email == "" {
        return nil, errors.New("email is required")
    }
    return &User{
        Name:  name,
        Email: email,
    }, nil
}
```

Returning `*User` rather than `User` is the usual choice, because the caller
almost always wants to pass it around without copying, and because it lets you
return `nil` on failure.

#bridge[
  A JavaScript constructor cannot fail cleanly: it either throws or returns a
  half-built object. A Go constructor is an ordinary function, so it returns an
  error like everything else, and you cannot get a `*User` back unless it
  succeeded.
]

== Comparing structs

Structs are comparable with `==` when all their fields are, and the comparison
is by value, field by field.

```go
a := User{ID: 1, Name: "Ada"}
b := User{ID: 1, Name: "Ada"}
fmt.Println(a == b)     // true
```

That is deep equality for free, and it is a genuine surprise if you are used to
`a === b` on objects being an identity check.

#gotcha[
  It stops working the moment a struct contains a slice, a map or a function,
  because those are not comparable. The compiler tells you. For those, use
  `reflect.DeepEqual`, or in tests `google/go-cmp`, which prints a readable
  diff instead of just `false`.
]

#project(title: "Project 4: a contact book in memory")[
  Build a small type with behaviour, and exercise value and pointer receivers.

  ```go
  package main

  import (
      "errors"
      "fmt"
      "strings"
  )

  type Contact struct {
      Name  string
      Email string
      Tags  []string
  }

  // Value receiver: reads only, so a copy is fine.
  func (c Contact) Display() string {
      if len(c.Tags) == 0 {
          return c.Name
      }
      return fmt.Sprintf("%s [%s]", c.Name, strings.Join(c.Tags, ", "))
  }

  // Pointer receiver: it modifies, so it must be a pointer.
  func (c *Contact) Tag(tag string) {
      for _, t := range c.Tags {
          if t == tag {
              return // already there, do nothing
          }
      }
      c.Tags = append(c.Tags, tag)
  }

  type Book struct {
      contacts []Contact
  }

  func (b *Book) Add(c Contact) error {
      if c.Email == "" {
          return errors.New("contact needs an email")
      }
      b.contacts = append(b.contacts, c)
      return nil
  }

  func (b *Book) Find(email string) (*Contact, bool) {
      for i := range b.contacts {
          if b.contacts[i].Email == email {
              return &b.contacts[i], true   // address of the real one
          }
      }
      return nil, false
  }

  func main() {
      var book Book
      if err := book.Add(Contact{Name: "Ada", Email: "ada@example.com"}); err != nil {
          fmt.Println(err)
          return
      }

      if c, ok := book.Find("ada@example.com"); ok {
          c.Tag("maths")
          c.Tag("maths") // ignored, already present
      }

      for _, c := range book.contacts {
          fmt.Println(c.Display())  // Ada [maths]
      }
  }
  ```

  The subtle line is `return &b.contacts[i], true`. Ranging with
  `for _, c := range b.contacts` would give you a *copy* of each contact, and
  tagging that copy would change nothing. Indexing with `i` and taking the
  address reaches the real element. Change it to the copying version and watch
  the tag vanish: that is the lesson of this chapter in one experiment.
]

#tryit[
  Add a `Remove(email string) bool` method to `Book`. You will need to build a
  new slice or use the `slices.Delete` helper from the standard library. Then
  add a `Count()` method and decide, with a reason you can say out loud,
  whether it takes a value or a pointer receiver.
]
