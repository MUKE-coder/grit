#import "../template.typ": *

= Maps, and why iteration order is random on purpose

A Go map is a hash table with typed keys and typed values. It is the closest
thing Go has to a JavaScript object used as a dictionary, and it differs in
four ways that each cause a specific bug the first time you meet them.


#bridge(title: "Things you have not met yet")[
  The concurrency-safe example at the end uses sync.Mutex (chapter 15) and
  defer (chapter 12). Skip it on a first read if you like: the map itself is
  the chapter.
]

== Declaring and using

```go
ages := map[string]int{
    "ada":   36,
    "grace": 45,
}

ages["alan"] = 41
fmt.Println(ages["ada"])    // 36
fmt.Println(len(ages))      // 3
delete(ages, "alan")
```

The type is `map[KeyType]ValueType`. Keys must be comparable, so strings, ints,
booleans and structs of those are fine; slices, maps and functions are not.

#gotcha[
  A `nil` map is readable but *not writable*, and this is the first map bug
  everyone hits:

  ```go
  var m map[string]int      // nil
  fmt.Println(m["a"])       // 0, fine
  fmt.Println(len(m))       // 0, fine
  m["a"] = 1                // panic: assignment to entry in nil map
  ```

  Always create one before writing: `m := map[string]int{}` or
  `m := make(map[string]int)`. The same trap hides inside structs, where a map
  field is nil until you assign one, so a constructor that forgets it produces
  a struct that panics on first write.
]

== Missing keys return the zero value

Reading a key that is not there is not an error. You get the zero value of the
value type.

```go
counts := map[string]int{}
fmt.Println(counts["nothing"])   // 0, not undefined
counts["word"]++                 // works: reads 0, writes 1
```

That last line is why counting in Go is so tidy: no "if the key exists, add
one, otherwise set it to one".

To tell "absent" from "present but zero", use the two-value form:

```go
v, ok := ages["nobody"]
if !ok {
    fmt.Println("not in the map")
}
```

`ok` is a bool. This is the *comma ok* idiom, and you will see the same shape
for type assertions and channel receives.

#bridge[
  In JavaScript, `obj.missing` gives `undefined`, and you distinguish absent
  from set-to-undefined with `in` or `hasOwnProperty`. Go's `v, ok :=` is the
  same idea with better ergonomics, because the check and the read are one
  operation and you cannot forget which you did.
]

== Iteration order is deliberately random

```go
for k, v := range ages {
    fmt.Println(k, v)
}
```

Run that twice and the order differs. This is not an implementation accident
that might settle down: the runtime randomises the starting point on purpose,
so that code which accidentally depends on order fails immediately and loudly
rather than in production two years later when the table grows.

To iterate in a stable order, sort the keys:

```go
keys := make([]string, 0, len(ages))
for k := range ages {
    keys = append(keys, k)
}
slices.Sort(keys)

for _, k := range keys {
    fmt.Println(k, ages[k])
}
```

#bridge[
  JavaScript objects do have an order: integer-like keys ascending, then string
  keys in insertion order. Plenty of code quietly relies on it. Go removes the
  temptation. If your output needs an order, you must say what that order is,
  which is usually an improvement to the code anyway.
]

== Maps are reference types

Passing a map to a function does not copy it. The function can modify the
caller's map.

```go
func add(m map[string]int) {
    m["new"] = 1        // visible to the caller
}
```

This is unlike a struct and unlike an array, and like a slice's elements. There
is no `const` map and no way to pass one read-only. If a function must not
modify a map, the only guards are a copy (`maps.Clone`) or a comment.

== Structs as values, and the thing you cannot do

A map of structs has one sharp edge:

```go
type Stat struct{ Count int }
m := map[string]Stat{"a": {Count: 1}}

m["a"].Count++      // compile error: cannot assign to struct field
```

The map returns a *copy* of the struct, so incrementing it would modify
something about to be discarded, and Go refuses rather than letting you write
code that does nothing. Two fixes:

```go
// 1. Read, modify, write back.
s := m["a"]
s.Count++
m["a"] = s

// 2. Store pointers instead.
m := map[string]*Stat{"a": {Count: 1}}
m["a"].Count++      // fine
```

The second is usually what you want when the values are large or frequently
updated. The first is fine for small structs and avoids a layer of indirection.

== Sets

Go has no set type. The idiom is a map with an empty struct value:

```go
seen := map[string]struct{}{}
seen["ada"] = struct{}{}

if _, ok := seen["ada"]; ok {
    fmt.Println("already seen")
}
```

`struct{}` is a struct with no fields and it occupies zero bytes, so the map
stores only keys. Using `map[string]bool` is also common and slightly more
readable at the cost of one byte per entry; both are accepted Go.

#bridge[
  JavaScript has a real `Set` with `.has()` and `.add()`. Go's version is
  wordier, and that `struct{}{}` reads as line noise until you have seen it a
  dozen times. Read it as: the type `struct{}`, then an empty literal of it.
]

== Maps are not safe for concurrent use

Two goroutines writing to the same map at the same time will crash the program
with a runtime error, deliberately. Go detects it and stops rather than
corrupting the table silently.

The answer is usually a mutex, which is chapter fifteen:

```go
type Counter struct {
    mu     sync.Mutex
    counts map[string]int
}

func (c *Counter) Inc(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.counts[key]++
}
```

#bridge[
  This has no JavaScript equivalent, because the event loop means two pieces of
  your code never run at literally the same instant. In Go they do. Any value
  shared between goroutines needs protection, and a map is the one where
  forgetting is immediately fatal rather than subtly wrong, which is mercifully
  the better failure.
]

#project(title: "Project 6: a word index with line numbers")[
  Extend the frequency counter into something genuinely useful: for each word,
  record every line it appears on. This is a `map[string][]int`, a map whose
  values are slices, and it is a shape you will use often.

  ```go
  package main

  import (
      "bufio"
      "fmt"
      "os"
      "slices"
      "strings"
  )

  func main() {
      if len(os.Args) < 2 {
          fmt.Fprintln(os.Stderr, "usage: index <file>")
          os.Exit(1)
      }

      f, err := os.Open(os.Args[1])
      if err != nil {
          fmt.Fprintln(os.Stderr, "error:", err)
          os.Exit(1)
      }
      defer f.Close()

      index := map[string][]int{}

      scanner := bufio.NewScanner(f)
      line := 0
      for scanner.Scan() {
          line++
          for _, w := range strings.Fields(scanner.Text()) {
              w = strings.ToLower(strings.Trim(w, ".,!?;:\"'()"))
              if len(w) < 4 {
                  continue    // skip short words
              }
              // Appending to a nil slice works, so no "if absent" branch.
              if !slices.Contains(index[w], line) {
                  index[w] = append(index[w], line)
              }
          }
      }
      if err := scanner.Err(); err != nil {
          fmt.Fprintln(os.Stderr, "reading:", err)
          os.Exit(1)
      }

      words := make([]string, 0, len(index))
      for w := range index {
          words = append(words, w)
      }
      slices.Sort(words)

      for _, w := range words {
          if len(index[w]) > 2 {   // only words on three or more lines
              fmt.Printf("%-16s %v\n", w, index[w])
          }
      }
  }
  ```

  The line that rewards attention is `index[w] = append(index[w], line)`. If
  `w` is not in the map, `index[w]` is a nil slice, and appending to a nil
  slice is legal and returns a one-element slice. So the "create the entry if
  missing" branch that you would write in most languages simply does not exist
  here. Two zero-value rules, slices and maps, combining to delete a line of
  code.

  Also note `bufio.Scanner`: it reads line by line without loading the file
  into memory, so this works on a file larger than your RAM, unlike the
  `os.ReadFile` in the last project. And `scanner.Err()` must be checked,
  because the loop ends identically on end-of-file and on a read error.
]

#tryit[
  Change the value type to a small struct holding both the line numbers and a
  total count, and watch the compiler stop you from writing `index[w].Count++`.
  Fix it both ways described above, then decide which you prefer and why. That
  decision, made consciously once, is the difference between copying Go idioms
  and understanding them.
]

#assignment(name: "Maps - Assignment One")[
  A word counter, which is the smallest program that forces every map lesson.

  In `part1-language/07-maps`, read a text file and report the ten most common
  words.

  *One.* Count with a `map[string]int`, using the fact that a missing key reads
  as zero rather than checking first.

  *Two.* Sort the results before printing, and say in a comment why printing
  straight from the map would make the output differ between runs.

  *Three.* Group with a `map[string][]string`: words by first letter. Note what
  happens when you append to a key that does not exist yet.

  *Four.* Try to assign to a field of a struct stored in a map, keep the
  compiler error in a comment, and fix it the way the chapter says.

  The website's *Grouping With a Map of Slices* is the same shape, smaller.]
