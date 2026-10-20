#import "../template.typ": *

#backmatter.update(true)

= Appendix: the JavaScript to Go phrasebook

The fastest way to write Go in your first month is to know what the thing you
already do is called. This is that list: the JavaScript you reach for, and the
Go you want instead.

It is a lookup table, not a tutorial. Each row names the chapter that explains
why.

== Values and types

#table(
  columns: (1fr, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [*JavaScript*], [*Go*],
  [`let x = 1`], [`x := 1`],
  [`const x = 1`], [`const x = 1`, or `x := 1` you do not reassign],
  [`undefined` / `null`], [The zero value, or a pointer that is `nil`. Ch 3],
  [`typeof x`], [`fmt.Sprintf("%T", x)`, or a type switch. Ch 10],
  [`Number(s)`], [`strconv.Atoi(s)`, which returns an error too],
  [`String(n)`], [`strconv.Itoa(n)` or `fmt.Sprint(n)`],
  [`` `hi ${name}` ``], [`fmt.Sprintf("hi %s", name)`],
  [`x ?? y`], [`if x == zero { x = y }`. No coalescing operator],
  [`x?.y`], [`if x != nil { ... }`. No optional chaining. Ch 8],
)

== Collections

#table(
  columns: (1fr, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [*JavaScript*], [*Go*],
  [`[]`], [`[]T{}` or `make([]T, 0, n)`. Ch 6],
  [`arr.push(v)`], [`arr = append(arr, v)`, and the assignment matters],
  [`arr.length`], [`len(arr)`],
  [`arr.map(f)`], [A `for` loop, or a generic helper you write once. Ch 17],
  [`arr.filter(f)`], [A `for` loop with an `if`],
  [`arr.find(f)`], [A `for` loop with an early `return`],
  [`arr.includes(v)`], [`slices.Contains(arr, v)`],
  [`arr.sort(cmp)`], [`sort.Slice(arr, less)` or `slices.SortFunc`],
  [`[...a, ...b]`], [`append(append([]T{}, a...), b...)`],
  [`{}`], [`map[string]T{}` or a struct. Ch 5, Ch 7],
  [`obj.key`], [`m["key"]`, or `s.Field` on a struct],
  [`key in obj`], [`v, ok := m[key]`. Ch 7],
  [`delete obj.key`], [`delete(m, key)`],
  [`Object.keys(o)`], [A loop over the map, then sort. Order is random],
  [`new Set()`], [`map[T]struct{}{}`],
)

== Functions and errors

#table(
  columns: (1fr, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [*JavaScript*], [*Go*],
  [`function f(a) {}`], [`func f(a int) {}`],
  [`const f = (a) => a * 2`], [`f := func(a int) int { return a * 2 }`],
  [`try { } catch (e) { }`], [`v, err := f()` then `if err != nil`. Ch 4],
  [`throw new Error("x")`], [`return errors.New("x")`],
  [`e instanceof MyError`], [`errors.As(err, &myErr)`. Ch 11],
  [`e.cause`], [`errors.Unwrap(err)`, from `%w`],
  [`finally { }`], [`defer`. Ch 12],
  [`process.exit(1)`], [`os.Exit(1)`, or return an error to `main`],
)

== Concurrency

#table(
  columns: (1fr, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [*JavaScript*], [*Go*],
  [`async function f()`], [Nothing. Every function can block. Ch 13],
  [`await f()`], [`f()`],
  [`f().then(...)`], [`go func() { ... }()`],
  [`Promise.all(ps)`], [`sync.WaitGroup`, or `errgroup`. Ch 13],
  [`Promise.race(ps)`], [`select` over channels. Ch 14],
  [`AbortController`], [`context.WithCancel`. Ch 16],
  [`setTimeout(f, ms)`], [`time.AfterFunc(d, f)`],
  [`setInterval(f, ms)`], [`time.NewTicker(d)` and a loop],
  [A shared mutable object], [`sync.Mutex`, or do not share it. Ch 15],
)

== Modules and tooling

#table(
  columns: (1fr, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [*JavaScript*], [*Go*],
  [`package.json`], [`go.mod`. Ch 9],
  [`node_modules/`], [A module cache outside your project],
  [`npm install x`], [`go get x`],
  [`npm install`], [`go mod download`, usually implicit],
  [`import { x } from "y"`], [`import "y"`, then `y.X`],
  [`export const x`], [A capital letter. Ch 2],
  [`export default`], [Nothing. Name the thing],
  [`npm test`], [`go test ./...`. Ch 18],
  [`npm run build`], [`go build ./...`. Ch 25],
  [`nodemon`], [`air`, or just rebuild. It is fast],
  [`eslint`], [`go vet`, plus `gofmt` which is not optional],
  [`prettier`], [`gofmt`. There are no options and that is the feature],
)

== The standard library, by what you would npm install

#table(
  columns: (1fr, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [*npm package*], [*Go*],
  [`express`], [`net/http`. Ch 21],
  [`axios` / `node-fetch`], [`net/http` client],
  [`dotenv`], [`os.Getenv`, plus twenty lines. Ch 24],
  [`uuid`], [`crypto/rand`, or one small dependency],
  [`lodash`], [`slices`, `maps`, and a loop],
  [`jest` / `vitest`], [`testing`. Ch 18],
  [`zod`], [Struct tags plus a validator. Ch 19, Ch 36],
  [`date-fns`], [`time`. Ch 40],
  [`winston` / `pino`], [`log/slog`],
  [`commander` / `yargs`], [`flag`, or Ch 49],
)

== Five habits to unlearn

*Checking for truthiness.* `if x` is not valid Go, and that is deliberate. An
empty string, a zero and a `nil` are different facts. Ch 3.

*Reaching for a library first.* The standard library covers more than you
expect, and a Go project with four dependencies is normal rather than
impressive. Ch 9.

*Returning different shapes from one function.* Go has one return type and an
error. A function that sometimes returns a string and sometimes an object is
not expressible, and the discipline is the point. Ch 4.

*Mutating what you were handed.* Everything is passed by value, but a slice or
a map holds a reference to shared memory. The function that sorts its argument
has changed the caller's data. Ch 6, Ch 8.

*Ignoring the second return value.* `_` is a decision to proceed on something
you did not look at, and the compiler will let you. The cases where that is
genuinely right are rare and deserve a comment. Ch 4, Ch 40.
