#import "../template.typ": *

= Installing Go and your first program

This chapter gets Go on your machine, gets a program running, and introduces
the toolchain. If you have ever set up a Node project you will find this
shorter than you expect, because most of the steps you are used to do not
exist.

== Installing

Go to #link("https://go.dev/dl/")[go.dev/dl] and download the installer for
your platform. There is no version manager to install first, no equivalent of
`nvm`, and for most people there never needs to be. Go's compatibility promise
means a program written for Go 1.0 in 2012 still compiles today, so "which
version" is rarely a question that costs you an afternoon.

Check it worked:

```bash
go version
```

You should see something like `go version go1.25.0 windows/amd64`. If the
command is not found, the installer did not put Go on your `PATH`; restart your
terminal first, because that fixes it more often than anything else.

#bridge[
  There is no `nvm`, no `.nvmrc`, and no per-project runtime. One Go
  installation builds every project on your machine. When a project needs a
  specific version, it says so in one line of `go.mod` and the toolchain
  downloads that compiler itself.
]

== Your first program

Make a directory and put one file in it.

```go
// main.go
package main

import "fmt"

func main() {
    fmt.Println("Hello from Go")
}
```

Run it:

```bash
go run main.go
```

Four things in that tiny file deserve explanation, because each is a rule you
will meet constantly.

`package main` declares which package this file belongs to. Every Go file
starts with a package declaration. The name `main` is special: it means this
package builds into a runnable program rather than a library.

`import "fmt"` pulls in the formatting package from the standard library. Note
there is no `npm install`. `fmt` ships with Go, as does an HTTP server, a JSON
encoder, a test runner, a crypto suite and a templating engine.

`func main()` is the entry point. A program starts here. There is exactly one
`main` function in a program.

`fmt.Println` prints a line. Capital `P` matters enormously, and it is the next
section.

== Capital letters are the access modifier

This is the first genuinely surprising rule in Go, and it catches everybody.

*An identifier starting with a capital letter is exported. One starting with a
lowercase letter is private to its package.* That is the entire visibility
system. There is no `public`, no `private`, no `export` keyword.

```go
package money

// Convert is callable from other packages, because it is capitalised.
func Convert(cents int) string { ... }

// round is not. Only code inside package money can call it.
func round(f float64) int { ... }
```

#bridge[
  In JavaScript you choose what leaves a module with `export`. In Go the name
  itself decides. `fmt.Println` works because `Println` is capitalised inside
  package `fmt`; if it were `println` you could not reach it. This applies to
  functions, types, struct fields, constants and variables alike.
]

#gotcha[
  This bites hardest with struct fields and JSON. A lowercase field is
  invisible to the JSON encoder, because the encoder lives in another package
  and genuinely cannot see it. A struct whose fields are all lowercase encodes
  to `{}`, with no error and no warning. Chapter sixteen covers this properly,
  but remember the symptom now: empty JSON usually means lowercase fields.
]

== Modules: go.mod is your package.json

A single file is fine for a demo. A real project is a *module*. Initialise one:

```bash
mkdir hello && cd hello
go mod init github.com/yourname/hello
```

That writes `go.mod`:

```
module github.com/yourname/hello

go 1.25
```

That is the whole file. Compare it to a `package.json` for an equivalent
project and notice what is missing: no `scripts`, no `devDependencies`, no
`main`, no `type: module`. The module path doubles as the import path, which is
why it looks like a URL: if you push this to GitHub, anybody can import it at
exactly that address with no registry publish step.

Add a dependency and `go.mod` grows a line:

```bash
go get github.com/google/uuid
```

```
module github.com/yourname/hello

go 1.25

require github.com/google/uuid v1.6.0
```

#bridge[
  `go.sum` is your lockfile, equivalent to `package-lock.json`. Commit it.
  There is no `node_modules`: downloaded packages go to a shared module cache
  outside your project, and the compiler links what it needs directly into the
  binary. Deleting a project directory does not take 300MB of dependencies with
  it, because they were never there.
]

== The commands you will actually use

Six commands cover almost everything.

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 6pt,
  [`go run .`], [Compile and run the current package, without keeping the binary],
  [`go build`], [Compile to an executable in the current directory],
  [`go test ./...`], [Run every test in this module],
  [`go mod tidy`], [Add what is imported, remove what is not. Run it often],
  [`gofmt -w .`], [Format every file. Your editor probably does this on save],
  [`go vet ./...`], [Report suspicious code the compiler allows],
)

The `./...` pattern means "this directory and everything under it", and you
will type it constantly.

== Compiling a real binary

This is the part that sells Go to people who have deployed Node.

```bash
go build -o hello .
./hello
```

You now have one executable file. It contains your code, every dependency, and
the Go runtime. The target machine needs nothing installed. You can also build
for a platform you are not on:

```bash
GOOS=linux GOARCH=amd64 go build -o hello-linux .
```

That produces a Linux binary from a Mac or Windows laptop, with no container,
no cross-compilation toolchain and no configuration. Copy it to a Linux server
and it runs.

#project(title: "Project 1: a greeting with arguments")[
  Write a program that takes a name as a command line argument and greets it,
  falling back to "world" when none is given. You need `os.Args`, which is a
  slice of strings where `os.Args[0]` is the program name itself.

  ```go
  package main

  import (
      "fmt"
      "os"
  )

  func main() {
      name := "world"
      if len(os.Args) > 1 {
          name = os.Args[1]
      }
      fmt.Printf("Hello, %s\n", name)
  }
  ```

  Build it, run `./hello Ada`, then run it with no argument. Then break it
  deliberately: change `len(os.Args) > 1` to `len(os.Args) > 0` and run it with
  no argument to see what a real Go panic looks like. Reading panics early is
  worth more than avoiding them.
]

#tryit[
  Extend the greeting to accept several names and greet each on its own line.
  You will need a `for` loop over a slice, which looks like
  `for _, n := range os.Args[1:] { ... }`. The next chapter explains every
  piece of that line, so guessing now is fine and the compiler will tell you
  when you are wrong.
]

#assignment(name: "Hello Go - Assignment One")[
  A module, a program, and a binary you can hand to somebody.

  *One.* In `part1-language/02-hello`, run `go mod init`, write a program that
  prints a greeting, and run it with `go run .`.

  *Two.* Add a second file in the same package with one exported and one
  unexported function. Call both from `main`. Then try to call the unexported
  one from a different package and keep the compiler error in a comment.

  *Three.* Build a real binary and commit its size, not the binary itself:

  ```bash
  go build -o hello .
  ls -lh hello > BUILD.txt
  ```

  *Four.* Cross compile the same program for another operating system and add
  that line to `BUILD.txt`. You now have the one-file deployment this chapter
  claimed, and you did it in two commands.

  The website's *Your First Program* and *Print/Println and Raw Strings*
  exercises cover the same ground in the browser if you want the warm-up.]
