#import "../template.typ": *

= Generics, briefly and honestly

Go got generics in 1.18, after a decade of arguing about whether it should. The
result is deliberately less powerful than TypeScript's, and the community
consensus is that you should use them rarely. This chapter covers enough to
read other people's code and to recognise the handful of cases where writing
your own is right.

== The syntax

Type parameters go in square brackets before the arguments.

```go
func First[T any](items []T) (T, bool) {
    var zero T
    if len(items) == 0 {
        return zero, false
    }
    return items[0], true
}

n, ok := First([]int{1, 2, 3})        // T inferred as int
s, ok := First([]string{"a", "b"})    // T inferred as string
```

`T` is a type parameter and `any` is its *constraint*, meaning "any type at
all". Inference usually works, so you rarely write `First[int](...)`.

Note `var zero T`. You cannot return `nil` for an arbitrary `T`, because `T`
might be `int`. The zero value is the only thing you can produce without
knowing the type, and this little declaration appears in almost every generic
function.

#bridge[
  TypeScript's `function first<T>(items: T[]): T | undefined` is the same idea,
  with a much more capable type system behind it. TypeScript has conditional
  types, mapped types, template literal types and inference that borders on
  computation. Go has type parameters and constraints, and that is the lot. If
  you arrive expecting the TypeScript toolkit you will be disappointed, which
  is by design: the proposal was accepted on the promise that it would not make
  Go code hard to read.
]

== Constraints are interfaces

A constraint says what you may *do* with a value of that type.

```go
// Will not compile: you cannot add two values of an arbitrary type.
func Sum[T any](items []T) T {
    var total T
    for _, v := range items {
        // invalid operation: operator + not defined on total
        total += v
    }
    return total
}
```

You must constrain `T` to types that support `+`. Constraints are written as
interfaces containing *type sets*:

```go
type Number interface {
    ~int | ~int64 | ~float64
}

func Sum[T Number](items []T) T {
    var total T
    for _, v := range items {
        total += v
    }
    return total
}

Sum([]int{1, 2, 3})          // 6
Sum([]float64{1.5, 2.5})     // 4.0
```

The `~` means "any type whose *underlying* type is this", so a named type still
qualifies:

```go
type Celsius float64
Sum([]Celsius{20.5, 21.0})   // works because of ~float64
```

Without the tilde, `Celsius` would be rejected, which is almost never what you
want. Write `~` unless you have a reason not to.

== The constraints you will actually use

Two come from the standard library and cover most cases.

```go
import "cmp"

func Max[T cmp.Ordered](a, b T) T {    // anything with < and >
    if a > b { return a }
    return b
}
```

`cmp.Ordered` is every type supporting `<`, which is the numeric types and
`string`. `comparable` is built in and means "supports `==`", which is what map
keys require:

```go
func Keys[K comparable, V any](m map[K]V) []K {
    out := make([]K, 0, len(m))
    for k := range m {
        out = append(out, k)
    }
    return out
}
```

== Generic types

Types can take parameters too.

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(v T) {
    s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    v := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return v, true
}

s := &Stack[string]{}
s.Push("a")
v, ok := s.Pop()
```

#gotcha[
  *Methods cannot introduce new type parameters.* This compiles nowhere:

  ```go
  func (s *Stack[T]) Map[U any](f func(T) U) *Stack[U] {   // illegal
  ```

  If you want a `Map` that changes the element type, it has to be a package
  function, not a method:

  ```go
  func MapStack[T, U any](s *Stack[T], f func(T) U) *Stack[U] { ... }
  ```

  This single restriction is why Go will never have fluent chained collection
  methods like `items.map(...).filter(...)`. It was a deliberate choice to keep
  the type checker simple, and it is the main reason generics changed Go code
  far less than people expected.
]

== What the standard library gives you

Since 1.21 there are generic helpers, and you should use these rather than
write your own.

```go
import ("slices"; "maps")

slices.Contains(s, v)        // bool
slices.Index(s, v)           // int, -1 if absent
slices.Sort(s)               // for cmp.Ordered
slices.SortFunc(s, cmp)      // for anything
slices.Reverse(s)
slices.Max(s); slices.Min(s)
slices.Clone(s)
slices.Equal(a, b)

maps.Keys(m)                 // an iterator, not a slice
maps.Values(m)
maps.Clone(m)
```

#gotcha[
  `maps.Keys` returns an *iterator* (`iter.Seq[K]`), not a slice. To get a
  slice, collect it:

  ```go
  keys := slices.Collect(maps.Keys(m))
  ```

  This trips up anybody reading older blog posts, where an earlier experimental
  version did return a slice.
]

== When to write your own, and when not to

The honest guidance, which matches what the Go team says: *write the concrete
version first*. Write the generic version when you have written the same
function for the third type and the bodies are identical.

Good reasons:

- A container: a stack, a set, a typed cache.
- A function over any slice or map where the body does not care about the type.
- Removing an `any` and the type assertions that came with it.

Bad reasons:

- "It might be reused later." It usually is not.
- To avoid writing a four-line loop. The loop is clearer.
- To model a hierarchy. Go has interfaces for that, and they are better at it.

#bridge[
  The decisive difference from TypeScript: TypeScript's generics are erased at
  runtime and cost nothing, so using them liberally is free. Go's are compiled,
  and a generic function is often *slower* than the concrete one because of how
  the compiler shares instantiations for pointer-shaped types. If you are
  reaching for generics for performance, measure. The usual answer is that the
  concrete function wins.
]

== Interfaces are still usually the answer

The mistake to avoid is reaching for a type parameter where an interface does
the job better.

```go
// Generic: the caller must have a concrete type, and the body can only do
// what the constraint allows.
func PrintAll[T fmt.Stringer](items []T) {
    for _, v := range items {
        fmt.Println(v.String())
    }
}

// Interface: accepts a mixed slice of different types. Usually better.
func PrintAll(items []fmt.Stringer) {
    for _, v := range items {
        fmt.Println(v.String())
    }
}
```

The rule of thumb: if the function's behaviour depends only on *methods*, use
an interface. Use a type parameter when you need the *actual type* back out,
which is what interfaces cannot do.

#project(title: "Project 16: a typed result cache")[
  A small generic container that is a genuinely good use of type parameters:
  the body does not care what it stores, and the caller gets its type back with
  no assertion.

  ```go
  package main

  import (
      "fmt"
      "sync"
      "time"
  )

  type cached[V any] struct {
      value   V
      expires time.Time
  }

  // Cache is generic over both key and value. K must be comparable because
  // it is a map key; V can be anything at all.
  type Cache[K comparable, V any] struct {
      mu    sync.RWMutex
      items map[K]cached[V]
      ttl   time.Duration
  }

  func NewCache[K comparable, V any](ttl time.Duration) *Cache[K, V] {
      return &Cache[K, V]{items: map[K]cached[V]{}, ttl: ttl}
  }

  func (c *Cache[K, V]) Set(k K, v V) {
      c.mu.Lock()
      defer c.mu.Unlock()
      c.items[k] = cached[V]{value: v, expires: time.Now().Add(c.ttl)}
  }

  func (c *Cache[K, V]) Get(k K) (V, bool) {
      c.mu.RLock()
      e, ok := c.items[k]
      c.mu.RUnlock()

      var zero V
      if !ok || time.Now().After(e.expires) {
          return zero, false
      }
      return e.value, true
  }

  // GetOrLoad is where the design earns itself: the loader returns a V and
  // the caller receives a V, with no any and no type assertion anywhere.
  func (c *Cache[K, V]) GetOrLoad(k K,
      load func(K) (V, error)) (V, error) {
      if v, ok := c.Get(k); ok {
          return v, nil
      }
      v, err := load(k)
      if err != nil {
          var zero V
          return zero, err
      }
      c.Set(k, v)
      return v, nil
  }

  type User struct {
      ID   int
      Name string
  }

  func main() {
      users := NewCache[int, User](time.Minute)
      counts := NewCache[string, int](time.Minute)

      calls := 0
      loadUser := func(id int) (User, error) {
          calls++
          return User{ID: id, Name: fmt.Sprintf("user-%d", id)}, nil
      }

      u, _ := users.GetOrLoad(1, loadUser)
      fmt.Println(u.Name, "loader calls:", calls)   // user-1  1

      u, _ = users.GetOrLoad(1, loadUser)
      fmt.Println(u.Name, "loader calls:", calls)   // user-1  1 (cached)

      // u is a User, not an any. This compiles, and would not if the cache
      // stored any.
      fmt.Println("id plus one:", u.ID+1)

      counts.Set("hits", 42)
      n, _ := counts.Get("hits")
      fmt.Println("count doubled:", n*2)            // 84, n is an int
  }
  ```

  Now prove the value. Rewrite `Cache` with `map[any]any` and no type
  parameters, then try `u.ID+1`: you need `u.(User).ID+1`, an assertion that
  can panic, and the compiler no longer checks that you put a `User` in and got
  a `User` out. That assertion, repeated at every call site, is exactly what
  generics removed.

  Then try to add a `Map` method that converts a `Cache[K, V]` into a
  `Cache[K, W]` and watch the compiler refuse, because methods cannot introduce
  type parameters. Make it a package function instead and notice how that
  constraint pushes Go code away from fluent chaining.
]

#tryit[
  Add a `Keys() []K` method, then write a package-level `Filter[K comparable, V
  any](c *Cache[K, V], keep func(V) bool) []V`. Then stop and ask whether
  either of them earns its place, or whether the caller ranging over a plain
  map would have been clearer. That question, asked honestly every time, is the
  whole discipline of generics in Go.
]

#assignment(name: "Generics - Assignment One")[
  Write the generic versions of what you already wrote, then judge them.

  In `part2-ideas/17-generics`:

  *One.* Rewrite chapter six's `Filter`, `Chunk` and `Dedupe` to work for any
  element type. `Dedupe` needs `comparable`, which is the constraint doing real
  work.

  *Two.* A generic `Set[T comparable]` with `Add`, `Has` and `Items`.

  *Three.* A `Map[T, U any]` that changes the element type, and one use of it
  that would have needed `any` and a cast before.

  *Four.* Now the honest part. In `NOTES.md`, name one place in your own code
  where an interface is the better answer than a type parameter, and say why.
  A chapter on generics that does not leave you more careful has failed.]
