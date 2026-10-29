#import "../template.typ": *

= Pointers, without the fear

Pointers have a reputation they do not deserve, earned in languages where you
can do arithmetic on them and walk off the end of memory. Go's pointers cannot
do that. There is no pointer arithmetic, and the garbage collector means you
never free anything. What is left is a small, useful idea: *a pointer is the
address of a value rather than a copy of it*.

You have already been using them. This chapter names what was happening.


#bridge(title: "Things you have not met yet")[
  A couple of examples borrow from later chapters: errors.Is is chapter
  11, goroutines are chapter 13, and the json tags are chapter 19. They
  are there because pointers are what those chapters trip over.
]

== Everything in Go is passed by value

This is the single sentence that explains pointers, and it is worth reading
twice because it is the opposite of your JavaScript instinct.

When you pass anything to a function, Go copies it. An int, a struct, a slice
header, all copied. So a function cannot change what you passed:

```go
func rename(u User) {
    u.Name = "changed"      // changes the copy
}

u := User{Name: "Ada"}
rename(u)
fmt.Println(u.Name)         // "Ada"
```

#bridge[
  In JavaScript, objects are passed by *reference to the same object*, so a
  function that sets `user.name` changes the caller's object. Primitives are
  copied. Go removes the distinction: everything is copied, every time, and
  when you want the other behaviour you ask for it explicitly with a pointer.
  Less to remember, more to type.
]

== Two operators

```go
u := User{Name: "Ada"}

p := &u        // & takes the address: p is a *User
fmt.Println(*p)     // * dereferences: the User itself
fmt.Println(p.Name) // but you rarely need *, see below
```

`&x` means "the address of x". `*T` in a type means "pointer to T". `*p` as an
expression means "the value at p".

Now the function that works:

```go
func rename(u *User, name string) {
    u.Name = name
}

rename(&u, "Grace")
fmt.Println(u.Name)     // "Grace"
```

#gotcha[
  You almost never write `(*u).Name`. Go dereferences automatically for field
  access and method calls, so `u.Name` works whether `u` is a `User` or a
  `*User`. This is a deliberate convenience and it means pointer-ness mostly
  disappears from the code that uses the value. The place it reappears is `nil`.
]

== nil, and the panic you will meet

A pointer's zero value is `nil`. Using a field on a nil pointer crashes:

```go
var u *User          // nil
fmt.Println(u.Name)  // panic: runtime error: invalid memory address
                     // or nil pointer dereference
```

This is Go's `Cannot read properties of undefined`, with the same cause and a
louder failure: the program stops rather than continuing with a bad value.

The guard is an ordinary `if`:

```go
func display(u *User) string {
    if u == nil {
        return "(none)"
    }
    return u.Name
}
```

#bridge[
  TypeScript's `strictNullChecks` would make this a compile error. Go has no
  equivalent: a `*User` might be nil and the compiler will not stop you. This
  is the one area where TypeScript's type system is genuinely ahead, and it is
  worth saying rather than glossing over. Your defences in Go are: return
  values rather than pointers where you can, check at the boundary, and let the
  panic happen in development rather than coding defensively everywhere.
]

== When to use a pointer

Three reasons, and nothing else.

*One: you need to modify the caller's value.* This is the common case, and it
is why methods that change state take pointer receivers.

*Two: the struct is large and copied often.* A struct with twenty fields copied
in a tight loop is real work. For a three-field struct it does not matter, and
chasing this before measuring is how code gets uglier for no gain.

*Three: nil is a meaningful value.* A `*string` field distinguishes "the user
did not supply this" from "the user supplied an empty string". You met this in
chapter three as the limit of zero values. This is the fix.

```go
type UpdateUser struct {
    Name  *string `json:"name"`   // nil means "leave it alone"
    Email *string `json:"email"`
}

func apply(u *User, patch UpdateUser) {
    if patch.Name != nil {
        u.Name = *patch.Name       // explicitly set, even to ""
    }
    if patch.Email != nil {
        u.Email = *patch.Email
    }
}
```

That is how a PATCH endpoint tells "clear this field" apart from "do not touch
this field", and you will write it for real in part four.

== When not to use a pointer

Slices, maps and channels already hold an internal pointer. Passing
`*[]int` or `*map[string]int` is almost always wrong and marks code as written
by someone still finding their feet.

```go
func addItem(items []string) []string {   // correct
    return append(items, "new")
}

func addItem(items *[]string) {           // almost always wrong
    *items = append(*items, "new")
}
```

The slice version returns the new slice, which is the convention `append`
itself follows. The pointer version works and reads badly.

Small structs you only read are also better by value. A `Point{X, Y int}` is
sixteen bytes; copying it is cheaper than following a pointer.

== The loop variable, and a trap that Go fixed

For years this was the most common Go bug in existence:

```go
// Go 1.21 and earlier
for _, u := range users {
    go func() {
        fmt.Println(u.Name)   // every goroutine saw the SAME u
    }()
}
```

The loop variable was reused across iterations, so all the goroutines shared
one variable and usually printed the last element several times.

*Go 1.22 changed this.* Loop variables are now per-iteration, and the code
above does what it looks like it does.

#gotcha[
  You will still find the old workaround everywhere in blog posts and older
  codebases:

  ```go
  for _, u := range users {
      u := u              // "shadow" it, making a per-iteration copy
      go func() { ... }()
  }
  ```

  It is harmless on modern Go and no longer needed. But check your `go.mod`:
  the new behaviour only applies if the `go` directive says 1.22 or higher. A
  module that still says `go 1.21` gets the old semantics *even on a new
  toolchain*, which is a subtle and genuinely nasty way to inherit this bug.
]

== Pointers to struct fields

You can take the address of a field, and of a slice element, which is how the
contact book in chapter five worked:

```go
for i := range book.contacts {
    if book.contacts[i].Email == email {
        return &book.contacts[i]    // address of the real element
    }
}
```

What you cannot take the address of is a map entry, because the map may move
its contents when it grows:

```go
m := map[string]User{"a": {Name: "Ada"}}
p := &m["a"]        // compile error: cannot take address of m["a"]
```

This is the same restriction that stopped `m["a"].Name++` in chapter seven, and
it has the same two fixes: read-modify-write, or store `*User` in the map.

#project(title: "Project 7: an inventory with adjustable stock")[
  Build something where the pointer decisions actually matter, and where
  getting one wrong is visible.

  ```go
  package main

  import (
      "errors"
      "fmt"
  )

  type Item struct {
      SKU   string
      Name  string
      Stock int
  }

  type Inventory struct {
      items map[string]*Item    // pointers: we mutate these in place
  }

  func NewInventory() *Inventory {
      return &Inventory{items: map[string]*Item{}}
  }

  func (inv *Inventory) Add(sku, name string, stock int) {
      inv.items[sku] = &Item{SKU: sku, Name: name, Stock: stock}
  }

  // Get returns nil when absent, so callers must check. The bool makes
  // that
  // impossible to forget by accident.
  func (inv *Inventory) Get(sku string) (*Item, bool) {
      it, ok := inv.items[sku]
      return it, ok
  }

  var ErrInsufficientStock = errors.New("insufficient stock")

  func (inv *Inventory) Take(sku string, n int) error {
      it, ok := inv.items[sku]
      if !ok {
          return fmt.Errorf("taking %s: not found", sku)
      }
      if it.Stock < n {
          return fmt.Errorf("taking %d of %s: %w",
              n, sku, ErrInsufficientStock)
      }
      it.Stock -= n       // through the pointer, so the map sees it
      return nil
  }

  func main() {
      inv := NewInventory()
      inv.Add("A1", "Blue widget", 10)

      if err := inv.Take("A1", 3); err != nil {
          fmt.Println(err)
      }
      if err := inv.Take("A1", 99); err != nil {
          fmt.Println("expected:", err)
          fmt.Println("is stock error:",
              errors.Is(err, ErrInsufficientStock))
      }

      if it, ok := inv.Get("A1"); ok {
          fmt.Printf("%s: %d left\n", it.Name, it.Stock)   // 7 left
      }
  }
  ```

  Now do the experiment that teaches the chapter. Change the map to
  `map[string]Item`, holding values rather than pointers, and try to compile.
  The line `it.Stock -= n` will still compile, because `it` is now a copy, and
  the stock will never change. Run it and watch `10 left` print after you took
  three. Then make it work by reading, modifying and writing back.

  That bug, with values where you meant pointers, is the single most common
  mistake newcomers make with Go collections, and having caused it once on
  purpose you will recognise it instantly.
]

#tryit[
  Add a `Restock(sku string, n int) error` method, then add a `Low(threshold
  int) []*Item` that returns every item below a stock level. Decide whether
  `Low` should return `[]*Item` or `[]Item` and be able to say why. One of them
  lets a caller change your inventory by accident, which might be exactly what
  you want or exactly what you do not.
]

#assignment(name: "Pointers - Assignment One")[
  Three programs, each one line different.

  In `part1-language/08-pointers`:

  *One.* A `Counter` with `Inc()` on a value receiver. Call it ten times and
  print zero. Then change one character and print ten. Commit both, the broken
  one as `value_receiver.go.txt` so it does not compile into the build.

  *Two.* A function taking a struct by value and another taking a pointer.
  Print the addresses with `%p` to show what was copied.

  *Three.* Cause a nil pointer panic deliberately, read the stack trace, and
  write in `NOTES.md` which line of it actually told you where the bug was.

  *Four.* Write a function returning `*User` that returns `nil` for not found,
  and a caller that checks before dereferencing. That shape is everywhere in
  Part III.]
