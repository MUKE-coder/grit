#import "../template.typ": *

= Variables, types, and the compiler that argues with you

This is the chapter where Go feels most foreign, and it is worth knowing why in
advance. In JavaScript, a variable is a box that holds anything and can hold
something else tomorrow. In Go, a variable has one type forever, decided when
it is created, and the compiler refuses to build your program if you break that
promise.

That refusal is the feature. Work with it for a week and you will notice a
whole category of bug has stopped happening to you.

== Declaring things

Go has two ways to make a variable, and you will use one of them ninety percent
of the time.

```go
var name string = "Ada"   // the long form, type stated
var age = 36              // type inferred, still a var
// short form: declare and infer, inside functions
name := "Ada"
```

The `:=` operator is the one you want. It declares a new variable and infers
its type from the value on the right. It only works inside a function, which is
why package level variables use `var`.

#bridge[
  `:=` is closest to `let` in that it creates a new binding. There is no
  equivalent of `const` for variables created with `:=`. Go's `const` exists,
  but it is only for values known at compile time: numbers, strings, booleans.
  You cannot declare a constant struct or slice. If you want "do not reassign
  this", Go's answer is to write a small enough function that it is obvious.
]

An unused variable is a *compile error*, not a warning.

```go
func main() {
    count := 10
    fmt.Println("hello")
}
// ./main.go:4:5: declared and not used: count
```

This feels hostile on day one and reasonable by day ten. Dead variables are
usually the residue of a half-finished edit, and Go would rather stop you than
let the residue accumulate. The same rule applies to unused imports.

#gotcha[
  The escape hatch is `_`, the blank identifier, which means "I know there is a
  value here and I am deliberately ignoring it". You will use it constantly in
  loops and multiple returns. Do not reach for it to silence a real warning:
  `_ = count` compiles and tells the next reader nothing.
]

== The types you will actually use

Go has a lot of numeric types. You will use about four of them.

#table(
  columns: (auto, auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 6pt,
  [*Type*], [*JavaScript*], [*Notes*],
  [`string`], [`string`], [Immutable, UTF-8 bytes. Indexing gives a byte, not a character],
  [`int`], [`number`], [Whole numbers. 64 bit on any machine you own],
  [`float64`], [`number`], [Decimals. The default for a fractional literal],
  [`bool`], [`boolean`], [`true` or `false`. No truthiness, ever],
  [`byte`], [-], [Alias for `uint8`. A byte of data],
  [`rune`], [-], [Alias for `int32`. One Unicode code point],
)

The one to notice is that JavaScript's single `number` splits into `int` and
`float64`, and *Go will not convert between them for you*.

```go
count := 10          // int
price := 2.5         // float64
total := count * price
// invalid operation: count * price (mismatched types int and float64)
```

You must say what you mean:

```go
total := float64(count) * price   // 25
```

#bridge[
  JavaScript coerces almost anything: `"3" * 2` is `6`, `[] + {}` is a string,
  and `0.1 + 0.2` is famously not `0.3`. Go does none of this. There is no
  implicit conversion between any two types, not even `int` and `int64`. Every
  conversion is written out, which is more typing and removes an entire family
  of surprise.
]

== Zero values: there is no undefined

Every type in Go has a *zero value*, and a declared variable always has it.
There is no `undefined`, and no uninitialised memory.

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 6pt,
  [`string`], [`""`],
  [`int`, `float64`], [`0`],
  [`bool`], [`false`],
  [pointers, slices, maps, interfaces, functions], [`nil`],
)

```go
var name string
var count int
var ready bool
fmt.Println(name == "", count == 0, ready == false) // true true true
```

This is genuinely useful. A struct you create is immediately valid, with every
field at a sensible default, and you never check whether a field "exists".

#gotcha[
  It is also the one place zero values hurt: you cannot tell "the caller set
  this to zero" apart from "the caller did not set it". If that distinction
  matters, for example a form field where empty and absent mean different
  things, you need a pointer (`*int`, where `nil` means absent) or a dedicated
  type. Chapter eight covers pointers properly. This problem is real, and it is
  the usual reason you see `*string` in a struct.
]

== No truthiness

An `if` in Go takes a `bool`. Not a value that is bool-ish. A `bool`.

#compare(
  "if (name) {
  greet(name);
}
if (items.length) {
  render(items);
}",
  "if name != \"\" {
    greet(name)
}
if len(items) > 0 {
    render(items)
}",
)

More characters, and no memorising which values are falsy. There is no `==` and
`===` distinction either, because there is only one `==` and it does not
coerce.

== Strings are bytes, and that matters

A Go string is a read-only slice of bytes holding UTF-8. Indexing it gives you
a byte.

```go
s := "héllo"
fmt.Println(len(s))     // 6, not 5: é takes two bytes
fmt.Println(s[1])       // 195, a byte value, not "é"
```

To work with characters, range over the string, which decodes runes:

```go
for i, r := range "héllo" {
    fmt.Printf("%d: %c\n", i, r)
}
// 0: h
// 1: é      <- index jumps to 3 next, because é was two bytes
// 3: l
// 4: l
// 5: o
```

#bridge[
  JavaScript strings are UTF-16 and have the same class of problem: `"👍".length`
  is `2`. Both languages make you think about encoding the moment you leave
  ASCII. Go is at least explicit that you are holding bytes.
]

Strings are immutable, so building one in a loop with `+=` allocates a new
string every time. For anything beyond a few concatenations use a builder:

```go
var b strings.Builder
for _, word := range words {
    b.WriteString(word)
    b.WriteString(" ")
}
result := b.String()
```

== Control flow, briefly

`if` and `for` are almost familiar, with two differences worth stating.

There are no parentheses around the condition, and braces are mandatory even
for one line.

```go
if count > 10 {
    fmt.Println("plenty")
} else if count > 0 {
    fmt.Println("some")
} else {
    fmt.Println("none")
}
```

`for` is the only loop keyword. It covers every case JavaScript spreads across
`for`, `while` and `for...of`.

```go
for i := 0; i < 10; i++ { }      // classic
for count < 10 { }               // while
for { }                          // infinite
for i, v := range items { }      // for...of, with the index
```

An `if` can carry a statement before the condition, scoped to the `if`. You
will see this constantly in real code:

```go
if err := doSomething(); err != nil {
    return err
}
// err does not exist here, which is the point
```

#bridge[
  `switch` in Go needs no `break`: cases do not fall through. You can switch on
  no value at all, which makes it a tidy replacement for a long `else if`
  chain:

  ```go
  switch {
  case score > 90:
      grade = "A"
  case score > 80:
      grade = "B"
  default:
      grade = "C"
  }
  ```
]

#project(title: "Project 2: a temperature table")[
  Write a program that prints a table converting Celsius to Fahrenheit from
  minus 40 to 100 in steps of 20. The formula is `f = c*9/5 + 32`.

  ```go
  package main

  import "fmt"

  func main() {
      fmt.Printf("%8s %8s\n", "C", "F")
      for c := -40; c <= 100; c += 20 {
          f := float64(c)*9/5 + 32
          fmt.Printf("%8d %8.1f\n", c, f)
      }
  }
  ```

  Two things to notice. `float64(c)` is required: without it you get integer
  division and wrong answers for anything not divisible by five, which is the
  bug this chapter exists to prevent. And `%8.1f` means "right align in eight
  columns, one decimal place", which is how Go does string formatting.

  Now break it on purpose: remove the `float64(...)` conversion and read the
  error. Then change it to `c*9/5` with `c` still an `int` and see that it
  compiles but produces wrong numbers. The compiler catches the type error. It
  cannot catch the maths error, and knowing which is which is the skill.
]

#tryit[
  Add a third column showing Kelvin (`k = c + 273.15`), and right align the
  header over it. Then make the step size a variable at the top of `main` and
  confirm the table still lines up when you change it to 10.
]
