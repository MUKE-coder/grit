#import "../template.typ": *

= Slices: the array that is not an array

Go has arrays and it has slices. You will use slices for essentially
everything, and the distinction matters exactly twice: when you declare one,
and when a surprising bug appears. This chapter covers both.


#bridge(title: "Things you have not met yet")[
  One example counts words with a map. Maps are chapter 7; here it is
  just somewhere to put a count.
]

== Arrays have a fixed length, in the type

```go
var a [3]int        // an array of exactly three ints
a[0] = 1
fmt.Println(len(a)) // 3
```

The length is *part of the type*. `[3]int` and `[4]int` are different types,
and a function taking `[3]int` will not accept `[4]int`. Arrays are also copied
whole on assignment.

This is almost never what you want, which is why arrays are rare in Go code.

== Slices are a view onto an array

```go
nums := []int{1, 2, 3}      // no length in the brackets: a slice
nums = append(nums, 4)
fmt.Println(len(nums))      // 4
```

A slice is a small three-field value: a pointer to an underlying array, a
length, and a capacity. That is the whole model, and understanding it explains
every slice surprise you will ever hit.

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 6pt,
  [`len(s)`], [How many elements you can index right now],
  [`cap(s)`], [How many fit before the underlying array must be reallocated],
)

```go
s := make([]int, 2, 5)   // len 2, cap 5
fmt.Println(len(s), cap(s))  // 2 5
s = append(s, 10)            // fits in the existing array
fmt.Println(len(s), cap(s))  // 3 5
```

`make([]T, len, cap)` creates one with room reserved. When you know roughly how
many elements are coming, say so: it avoids repeated reallocation.

#bridge[
  A JavaScript array handles all of this invisibly. `push` grows it and you
  never think about capacity. Go exposes the machinery, which costs you some
  attention and buys you control over allocation, which matters a great deal in
  a hot loop and not at all in a handler that runs once per request.
]

== append, and the thing it does not promise

`append` returns a slice. *You must use the return value.*

```go
nums := []int{1, 2, 3}
append(nums, 4)          // compile error: result not used
nums = append(nums, 4)   // correct
```

The reason is that `append` may or may not allocate a new array. If there is
spare capacity it writes in place and returns a slice with a longer length. If
there is not, it allocates a bigger array, copies everything over, and returns
a slice pointing at the *new* array. The caller's original slice still points
at the old one.

This is where the famous slice bug lives.

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3]              // [2 3], sharing a's array
b[0] = 99
fmt.Println(a)           // [1 99 3 4 5]  <- a changed
```

Slicing does not copy. `b` is a window onto the same memory. Writing through
one is visible through the other.

#gotcha[
  It gets worse with `append`, because whether you corrupt the original depends
  on capacity, which depends on history:

  ```go
  a := []int{1, 2, 3, 4, 5}
  b := a[1:3]              // len 2, cap 4 (to the end of a)
  b = append(b, 99)        // fits in cap, writes into a's array
  fmt.Println(a)           // [1 2 3 99 5]  <- clobbered a[3]
  ```

  The fix, when you want a genuinely independent slice, is a three-index slice
  that caps capacity, or an explicit copy:

  ```go
  b := a[1:3:3]            // cap is now 2, so append must allocate
  // or
  b := slices.Clone(a[1:3])
  ```

  This bug is rare in everyday handler code and vicious when it appears. The
  tell is "a value changed and nothing in my code assigned to it".
]

== Iterating

```go
for i, v := range nums {
    fmt.Println(i, v)
}

for _, v := range nums {    // index discarded
    fmt.Println(v)
}

for i := range nums {       // value discarded
    nums[i] *= 2
}
```

`range` gives you a *copy* of each element. For ints that is irrelevant. For
structs it is the bug from the previous chapter: modifying `v` changes nothing.
Use the index when you need to write.

#bridge[
  Go has no `map`, `filter`, `reduce`, `find` or `some` for slices in the
  language. Since Go 1.21 the `slices` package covers some of it:

  ```go
  slices.Contains(nums, 3)        // bool
  slices.Index(nums, 3)           // position or -1
  slices.Sort(nums)               // in place
  slices.Reverse(nums)            // in place
  slices.Max(nums)                // and Min
  ```

  For anything else, write the loop. A filter is four lines:

  ```go
  var active []User
  for _, u := range users {
      if u.Active {
          active = append(active, u)
      }
  }
  ```
]

== Passing slices to functions

A slice is cheap to pass: you are copying three words, not the data. That means
a function *can* modify the caller's elements.

```go
func zero(nums []int) {
    for i := range nums {
        nums[i] = 0
    }
}

nums := []int{1, 2, 3}
zero(nums)
fmt.Println(nums)   // [0 0 0]
```

But it cannot change the caller's *length*, because the slice header itself was
copied:

```go
func addOne(nums []int) {
    nums = append(nums, 1)   // changes the local copy's header only
}

nums := []int{1, 2, 3}
addOne(nums)
fmt.Println(len(nums))   // still 3
```

The convention, and the reason `append` works the way it does: a function that
grows a slice returns it.

```go
func addOne(nums []int) []int {
    return append(nums, 1)
}
nums = addOne(nums)
```

== Sorting

```go
people := []User{{Name: "Zoe"}, {Name: "Ada"}}

slices.SortFunc(people, func(a, b User) int {
    return strings.Compare(a.Name, b.Name)
})
```

The comparison function returns a negative number, zero, or a positive number,
the same contract as JavaScript's `Array.prototype.sort` comparator. For stable
sorting use `slices.SortStableFunc`.

#project(title: "Project 5: word frequency, top ten")[
  Read a file, count how often each word appears, and print the ten most
  common. This uses slices, sorting and a map together, and it is a real tool
  you can point at a book.

  ```go
  package main

  import (
      "fmt"
      "os"
      "slices"
      "strings"
  )

  type pair struct {
      word  string
      count int
  }

  func main() {
      if len(os.Args) < 2 {
          fmt.Fprintln(os.Stderr, "usage: freq <file>")
          os.Exit(1)
      }

      data, err := os.ReadFile(os.Args[1])
      if err != nil {
          fmt.Fprintln(os.Stderr, "error:", err)
          os.Exit(1)
      }

      counts := map[string]int{}
      for _, w := range strings.Fields(string(data)) {
          w = strings.ToLower(strings.Trim(w, ".,!?;:\"'()"))
          if w != "" {
              counts[w]++
          }
      }

      pairs := make([]pair, 0, len(counts))
      for w, c := range counts {
          pairs = append(pairs, pair{w, c})
      }

      slices.SortFunc(pairs, func(a, b pair) int {
          if b.count != a.count {
              return b.count - a.count       // count descending
          }
          return strings.Compare(a.word, b.word) // then alphabetical
      })

      for _, p := range pairs[:min(10, len(pairs))] {
          fmt.Printf("%6d  %s\n", p.count, p.word)
      }
  }
  ```

  Three things worth noticing. `make([]pair, 0, len(counts))` reserves capacity
  up front because we know exactly how many there will be. The tie-break on
  word makes the output deterministic, which matters because map iteration
  order is random and without it two runs would disagree. And `min` is a
  builtin since Go 1.21, so `pairs[:min(10, len(pairs))]` does not panic on a
  file with fewer than ten distinct words.
]

#tryit[
  Add a flag so the user can ask for a different number of results, for example
  `freq book.txt 25`. Then add a stop-word list so "the", "a" and "and" are
  excluded, and decide whether a slice or a map is the better structure for
  holding those stop words. The next chapter will tell you whether you chose
  right.
]

#assignment(name: "Slices - Assignment One")[
  Write the three helpers you miss from JavaScript, then prove the trap.

  In `part1-language/06-slices`:

  *One.* `Filter`, `Chunk` and `Dedupe` over `[]string`, each with a
  table-driven test. You will write these generically in chapter seventeen, so
  keep them.

  *Two.* Reproduce the aliasing bug on purpose. Take a slice of a slice, append
  to it, and show the original changing. Then fix it with a copy, and write the
  one-line rule in a comment.

  *Three.* Sort a slice of structs by two keys, falling back to the second when
  the first ties.

  ```go
  sort.Slice(xs, func(i, j int) bool { ... })
  ```

  *Four.* In `NOTES.md`, state in one sentence what `append` does not promise.]
