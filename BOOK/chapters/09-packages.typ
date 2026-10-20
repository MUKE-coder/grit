#import "../template.typ": *

= Packages, modules and imports, against npm and ESM

You have been writing single-file programs. Real projects are many files, and
Go's rules for organising them are stricter and simpler than what you are used
to. This chapter is the one that makes a Go repository stop looking strange.

== A package is a directory

Every `.go` file belongs to a package, declared on its first line. *Every file
in a directory must declare the same package name*, and a directory is exactly
one package. That is the whole rule.

```
myapp/
├── go.mod              module github.com/you/myapp
├── main.go             package main
├── store/
│   ├── store.go        package store
│   └── user.go         package store
└── money/
    └── money.go        package money
```

Files in `store/` can see each other's identifiers with no import at all, even
the lowercase ones. There is no equivalent of importing between files in the
same folder, because they are already the same package.

#bridge[
  This is the biggest structural difference from JavaScript. In a Next.js
  project every file is its own module, and you import between them constantly.
  In Go, `store/store.go` and `store/user.go` share a namespace: a function
  declared in one is callable from the other with no ceremony. You import
  across *directories*, never within one.
]

== Importing

```go
package main

import (
    "fmt"               // standard library
    "net/http"

    "github.com/google/uuid"          // third party

    "github.com/you/myapp/store"      // your own
)

func main() {
    s := store.New()
    fmt.Println(uuid.NewString())
}
```

You refer to things by the *package name*, which is the last element of the
path by convention but is really whatever the files declare. `store.New`, not
`New`.

The grouping above is conventional: standard library first, then third party,
then your own module, separated by blank lines. `gofmt` sorts within groups and
`goimports` adds and removes imports as you type.

#gotcha[
  An unused import is a compile error, exactly like an unused variable. This is
  why `goimports` in your editor is not optional: without it you will spend
  real time adding and removing import lines by hand.

  The escape hatch, for a package imported only for its side effects, is the
  blank identifier:

  ```go
  import _ "github.com/lib/pq"   // registers a database driver
  ```
]

== The import path is the module path plus the directory

There is no `baseUrl`, no `paths` mapping, no `@/components` alias, and no
relative imports. `import "./store"` is not valid Go.

If your module is `github.com/you/myapp`, then the directory `store/` is
imported as `github.com/you/myapp/store`. Always the full path.

#bridge[
  Verbose compared to `@/store`, and it buys one thing: an import path is
  globally unambiguous. You can paste it into a browser and often land on the
  source. There is no build-time configuration that has to agree with your
  editor, your bundler and your test runner, because there is no mapping at
  all.
]

== internal/ is a real access rule

A directory named `internal` is enforced by the compiler: only code rooted at
its parent may import it.

```
myapp/
├── internal/
│   └── db/            importable by myapp and anything under it
└── pkg/
    └── client/        importable by anybody
```

If someone imports `github.com/you/myapp/internal/db` from a different module,
the build fails. This is how a library publishes a small surface while keeping
the rest private, and it is why nearly every serious Go project puts most of
its code under `internal/`.

#bridge[
  The nearest npm equivalent is the `exports` field in `package.json`, which
  controls what consumers can reach. Go's version is a directory name, needs no
  configuration and cannot be worked around.
]

== Naming, and the stutter rule

Package names are short, lowercase, single words. No underscores, no camelCase,
no plurals.

The identifier is read *with* the package name, which leads to the most quoted
style rule in Go:

```go
// Bad: the package name is already there.
store.NewStore()
http.HTTPServer
money.MoneyAmount

// Good:
store.New()
http.Server
money.Amount
```

#gotcha[
  Avoid `util`, `helpers`, `common` and `shared`. They are not packages, they
  are places things go to hide, and they grow into a dependency that everything
  imports and nobody understands. If a function does not belong anywhere, that
  usually means the thing it operates on wants a package of its own.
]

== Modules, versions and go.mod

A module is a tree of packages with one `go.mod` at the root. You saw it in
chapter two; here is what the rest of it means.

```
module github.com/you/myapp

go 1.25

require (
    github.com/gin-gonic/gin v1.10.0
    gorm.io/gorm v1.25.12
)

require (
    github.com/bytedance/sonic v1.12.0 // indirect
)
```

The second `require` block is indirect dependencies: things your dependencies
need. `go mod tidy` manages both blocks, adding what is imported and removing
what is not. Run it before every commit.

#bridge[
  `go.sum` is the lockfile and must be committed. Unlike `package-lock.json` it
  is not a full dependency tree: it is a list of cryptographic hashes. Go uses
  *minimal version selection*, meaning it picks the lowest version that
  satisfies everyone, so builds are reproducible from `go.mod` alone and the
  sum file exists to verify integrity rather than to pin.

  There is also no `devDependencies`. A tool you use at development time is
  either run with `go run` or tracked in a `tools.go` file.
]

Upgrading is explicit:

```bash
go get -u ./...          # upgrade everything (minor and patch)
go get github.com/gin-gonic/gin@v1.10.1   # one package, one version
go mod tidy
```

== Major versions are in the path

This surprises everyone. From v2 onward, the major version is part of the
import path:

```go
import "github.com/you/myapp/v2/store"
```

It looks strange and it solves a real problem: two major versions of the same
library can coexist in one build, because as far as the compiler is concerned
they are different packages. You will see this in Grit itself, which imports as
`github.com/MUKE-coder/grit/v3`.

== init, and why to avoid it

A package can declare `func init()`, which runs automatically before `main`.

```go
func init() {
    // runs once, before main, in no order you control across packages
}
```

Use it rarely. It makes startup order implicit, it cannot return an error, and
it runs even when a test only wanted one function from the package. Explicit
setup called from `main` is almost always better.

#project(title: "Project 8: split the inventory into packages")[
  Take the inventory from chapter seven and give it a real project layout. The
  point is to feel the import rules rather than to add features.

  ```
  inventory/
  ├── go.mod                  module github.com/you/inventory
  ├── main.go                 package main
  └── internal/
      └── inventory/
          ├── item.go         package inventory
          └── store.go        package inventory
  ```

  `internal/inventory/item.go`:

  ```go
  package inventory

  import "errors"

  var ErrInsufficientStock = errors.New("insufficient stock")

  type Item struct {
      SKU   string
      Name  string
      Stock int
  }

  // low is unexported, so only this package can call it. item.go and
  // store.go
  // are the same package, so store.go can use it with no import.
  func (i *Item) low(threshold int) bool {
      return i.Stock < threshold
  }
  ```

  `internal/inventory/store.go`:

  ```go
  package inventory

  import "fmt"

  type Store struct {
      items map[string]*Item
  }

  func New() *Store {
      return &Store{items: map[string]*Item{}}
  }

  func (s *Store) Add(sku, name string, stock int) {
      s.items[sku] = &Item{SKU: sku, Name: name, Stock: stock}
  }

  func (s *Store) Take(sku string, n int) error {
      it, ok := s.items[sku]
      if !ok {
          return fmt.Errorf("taking %s: not found", sku)
      }
      if it.Stock < n {
          return fmt.Errorf("taking %d of %s: %w",
              n, sku, ErrInsufficientStock)
      }
      it.Stock -= n
      return nil
  }

  // Low uses the unexported method from the other file. No import needed.
  func (s *Store) Low(threshold int) []*Item {
      var out []*Item
      for _, it := range s.items {
          if it.low(threshold) {
              out = append(out, it)
          }
      }
      return out
  }
  ```

  `main.go`:

  ```go
  package main

  import (
      "errors"
      "fmt"

      "github.com/you/inventory/internal/inventory"
  )

  func main() {
      inv := inventory.New()
      inv.Add("A1", "Blue widget", 10)

      if err := inv.Take("A1", 99); err != nil {
          // true
          fmt.Println(errors.Is(err, inventory.ErrInsufficientStock))
      }

      for _, it := range inv.Low(20) {
          fmt.Println("low:", it.Name, it.Stock)
      }
  }
  ```

  Now run three experiments. Call `inv.low(5)` from `main.go` and read the
  error: lowercase means package-private, and the compiler enforces it. Move
  the whole `internal` directory to `pkg` and confirm it still builds, then put
  it back. And try `import "./internal/inventory"` to see that relative imports
  genuinely do not exist.
]

#tryit[
  Add a second package, `internal/report`, with a function that takes a
  `*inventory.Store` and prints a summary. You will hit a design question
  immediately: `report` needs to read every item, but `Store.items` is
  unexported. Solve it by adding an exported method to `inventory` rather than
  by exporting the field, and notice that the package boundary just forced you
  into a better API than you would have written in one file.
]
