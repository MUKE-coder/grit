#import "../template.typ": *

= Build your own test runner, and what go test really does

`go test` feels like magic in a way the rest of Go does not. There is no
registration, no decorator, no configuration file. You write a function whose
name starts with `Test`, and it runs.

This chapter is how. By the end you will have a runner that discovers tests,
builds a binary, runs them in parallel with timeouts and isolation, reports
failures usefully, and emits machine-readable output. And you will know
precisely which of Go's conveniences are compiler features and which are
ordinary code you could have written.

== The trick: it writes a program

`go test` does not interpret your test files. It *generates a `main` package*
that imports your package, lists your test functions in a slice, and calls
`testing.Main` with them. Then it compiles that program and runs it.

```go
// Roughly what go test generates, simplified
func main() {
    tests := []testing.InternalTest{
        {Name: "TestAdd", F: pkg.TestAdd},
        {Name: "TestSub", F: pkg.TestSub},
    }
    testing.Main(matchString, tests, benchmarks, examples)
}
```

#bridge[
  This is the answer to a question most Go developers never ask: how is there
  no registration step?

  In JavaScript, `describe` and `it` are function calls that push into a
  registry at import time. In Go there is no import-time registry, because the
  *tooling* walks the source, finds the functions and writes the list. The
  magic is in the build step, not the language.

  Which means you can build it yourself, and that is exactly what this chapter
  does.
]

== Step one: discovery

Parse the package with `go/parser` and find the functions. This is the same
`go/ast` work `grit sync` does to read models, so the shape should be familiar:

```go
func discover(dir string) ([]Test, error) {
    fset := token.NewFileSet()
    pkgs, err := parser.ParseDir(fset, dir, isTestFile,
        parser.ParseComments)
    if err != nil {
        return nil, err
    }

    var out []Test
    for _, pkg := range pkgs {
        for path, file := range pkg.Files {
            for _, decl := range file.Decls {
                fn, ok := decl.(*ast.FuncDecl)
                if !ok || !isTestFunc(fn) {
                    continue
                }
                out = append(out, Test{
                    Name: fn.Name.Name,
                    File: path,
                    Line: fset.Position(fn.Pos()).Line,
                })
            }
        }
    }
    sort.Slice(out, func(i, j int) bool {
        return out[i].Name < out[j].Name
    })
    return out, nil
}
```

`isTestFunc` is the rule, written out:

```go
func isTestFunc(fn *ast.FuncDecl) bool {
    if fn.Recv != nil {                     // not a method
        return false
    }
    if !strings.HasPrefix(fn.Name.Name, "Test") {
        return false
    }
    // TestXxx: the character after "Test" must not be lowercase,
    // so Testing() is not a test and TestAdd is.
    rest := fn.Name.Name[len("Test"):]
    if rest != "" {
        r, _ := utf8.DecodeRuneInString(rest)
        if unicode.IsLower(r) {
            return false
        }
    }
    return hasSignature(fn, "*testing.T")
}
```

That lowercase rule is real and catches people: a function called `Testing` is
not a test, and `Testadd` is not either.

== Step two: generate and build

```go
const mainTmpl = `package main

import (
    "os"
    "testing"

    target "{{.ImportPath}}"
)

func main() {
    tests := []testing.InternalTest{
{{range .Tests}}        {Name: "{{.Name}}", F: target.{{.Name}}},
{{end}}    }
    m := testing.MainStart(deps{}, tests, nil, nil, nil)
    os.Exit(m.Run())
}
`
```

Write it to a temporary directory, then:

```bash
go build -o runner.test ./tmpdir
```

Notice what this means: *your test runner is a compiler front end*. It produces
Go source and hands it to the Go toolchain, which is how `go test` gets types,
build tags and vet checks for free.

#gotcha[
  The generated file imports your package, so a test in package `foo` can only
  see exported identifiers. That is why Go has *two* kinds of test file:
  `package foo` tests, compiled into the package and able to see everything,
  and `package foo_test` tests, which see only the public API.

  `go test` handles both by generating two sets of entries. It is worth
  supporting in your version, because the distinction is one of the better
  habits Go encourages: a `_test` package forces you to use your own API the
  way a caller does.
]

== Step three: running, in parallel, in isolation

```go
func (r *Runner) Run(ctx context.Context, tests []Test) Report {
    sem := make(chan struct{}, r.Parallel)
    results := make([]Result, len(tests))
    var wg sync.WaitGroup

    for i, t := range tests {
        wg.Add(1)
        go func(i int, t Test) {
            defer wg.Done()
            sem <- struct{}{}
            defer func() { <-sem }()

            results[i] = r.runOne(ctx, t)
        }(i, t)
    }
    wg.Wait()
    return report(results)
}
```

Each test runs as a subprocess invocation of the built binary with
`-test.run=^TestName$`:

```go
func (r *Runner) runOne(ctx context.Context, t Test) Result {
    ctx, cancel := context.WithTimeout(ctx, r.Timeout)
    defer cancel()

    dir, err := os.MkdirTemp("", "testrun-")
    if err != nil {
        return Result{Test: t, Err: err}
    }
    defer os.RemoveAll(dir)

    cmd := exec.CommandContext(ctx, r.Binary,
        "-test.run=^"+t.Name+"$", "-test.v")
    cmd.Dir = dir
    var out bytes.Buffer
    cmd.Stdout, cmd.Stderr = &out, &out

    start := time.Now()
    err = cmd.Run()

    return Result{
        Test:     t,
        Passed:   err == nil,
        Output:   out.String(),
        Duration: time.Since(start),
        TimedOut: ctx.Err() == context.DeadlineExceeded,
    }
}
```

A process per test is heavier than `go test`'s goroutine per test, and it buys
two things worth having: a panic kills one test rather than the run, and a test
that leaks goroutines or global state cannot affect its neighbours.

#gotcha[
  Each test gets its own temporary working directory, and that is not
  politeness.

  A test that writes `output.json` in the current directory passes alone and
  fails when run next to another test that does the same. Those failures are
  order dependent, so they appear when somebody adds an unrelated test, and
  they are blamed on the new test.

  Isolate the filesystem and the environment, and *randomise the order*
  occasionally. A suite that passes in a random order is a suite whose tests
  are genuinely independent.
]

== Step four: output people can use

Two audiences, as in chapter forty nine.

For a human, failures only, with the file and line so an editor can jump there:

```
FAIL  TestInvoiceTotal  (0.03s)
  invoice_test.go:42: Total() = 1999, want 2000

ok    38 passed, 1 failed, 2 skipped  in 1.4s
```

For a machine, one JSON object per event, which is what `go test -json`
produces and what every CI system consumes:

```json
{"Action":"run","Test":"TestInvoiceTotal"}
{"Action":"output","Test":"TestInvoiceTotal","Output":"..."}
{"Action":"fail","Test":"TestInvoiceTotal","Elapsed":0.03}
```

Streaming the events rather than printing a report at the end matters more than
it sounds: a twenty minute suite that prints nothing until it finishes is
indistinguishable from a hung one.

#gotcha[
  Print the *failures* last, or print them twice.

  Scrolling back through four thousand lines of passing output to find the one
  failure is a daily annoyance in a hundred CI systems. A summary at the end,
  repeating each failure with its output, costs nothing and is the single most
  appreciated feature of a runner.

  This is also why you never pipe a test log through `tail`: the summary
  survives and the context that explains it does not.
]

== The features that make a runner good

Having built the core, these are each an afternoon and each one you will miss
if you leave it out:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Filtering], [`-run` with a regular expression. The most used flag there
    is],
  [Timeouts], [Per test and overall, with the stack dumped on timeout],
  [Retry flaky], [`-retry 2`, and *report* the retry rather than hiding it],
  [Shuffle], [Random order with a printed seed, so a failure is reproducible],
  [Sharding], [`-shard 2/5`, for splitting a suite across CI machines],
  [Slowest], [Print the ten slowest tests. The list is always surprising],
)

#bridge[
  The retry one deserves a warning. A runner that silently retries turns a
  flaky test into a passing test, and the flakiness becomes invisible until it
  is twenty tests and nobody trusts the suite.

  Retry if you must, but print `FLAKY: TestX passed on attempt 2` and count it
  in the summary. Make the flakiness visible, or you are building a machine for
  ignoring problems.
]

== What `go test` does that yours will not

Worth knowing, because the list is a map of the remaining depth:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Caching], [Results keyed by the inputs: source, environment, files read.
    This is why an unchanged test prints `(cached)`],
  [Coverage], [Rewriting the source to add counters before compiling],
  [Benchmarks], [Adaptive iteration counts until the timing is stable],
  [Fuzzing], [Corpus management, mutation, minimisation of a failing input],
  [Subtests], [`t.Run`, with its own parallelism and naming rules],
  [Vet], [A subset of checks run automatically before the tests],
)

The caching one is the most interesting piece of engineering in the list: it
hashes everything that could affect the outcome, including the files the test
read, and that is why it is safe to trust.

#project(title: "Project 54: a runner you would use")[
  In `part6-from-scratch/58-runner`, build it and then run it on real code.

  *One: discovery.* Parse a package, list its tests with file and line.
  Compare your list with `go test -list '.*'` on the same package: they must
  match exactly, including the `Testing` and `Testadd` edge cases.

  *Two: generate and build* the runner binary.

  *Three: run* with configurable parallelism, a per-test timeout, a temporary
  working directory each, and a summary that repeats the failures.

  *Four: both outputs.* Human and JSON-lines. Pipe the JSON into `jq` to print
  the slowest five tests.

  *Five: the ordering bug.* Write two tests that both write `out.txt` in the
  current directory. Show them failing without isolation and passing with it.

  *Six: shuffle and shard.* Random order with a printed seed, and
  `-shard i/n`. Prove the shards are disjoint and cover everything.

  *Seven: run it on something real.* Point it at the test suite of a project
  from Part V and compare the output with `go test`. Every difference is
  either a feature you have not built or a bug.

  *Eight: time it.* Your runner against `go test` on the same package, both
  cold. Report both numbers and explain the difference, which will mostly be
  process startup.
]

#tryit[
  Run `go test -x ./...` on any package and read the commands it prints.

  You will see the generated `_testmain.go`, the build, and the invocation.
  Everything in this chapter is visible in that output, and after building your
  own, it reads as a description rather than as noise.
]

#assignment(name: "Test Runner - Assignment One")[
  Build the runner, and use it on real code rather than on examples.

  *One.* Discovery, generation, parallel isolated execution, timeouts, both
  output formats, shuffle with a seed, and sharding.

  *Two.* `COMPARE.md`: your test list against `go test -list`, your pass and
  fail counts against `go test`, and the timing of both on a real package.

  *Three.* The isolation demonstration: two tests that collide, failing and
  then passing.

  *Four.* One flaky test, run with retries, showing that your output *reports*
  the retry rather than hiding it.

  *Five.* A paragraph on caching: what would you have to hash to make a cached
  result trustworthy? List everything you can think of, then look up what the
  Go tool actually uses and note what you missed.

  This is the last build in the book. When it runs your own earlier chapters'
  tests, you have gone from "I have never written Go" to "I wrote the thing
  that runs my tests", which is a reasonable place to stop.
]
