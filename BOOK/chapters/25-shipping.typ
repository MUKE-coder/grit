#import "../template.typ": *

= Building, cross compiling and shipping one binary

This is the chapter that explains why people who deploy Node sound envious when
Go comes up. Deployment is `scp` and `./app`, and the rest of this chapter is
the detail around that sentence.

== The build

```bash
go build -o app .
./app
```

One file. It contains your code, every dependency, and the Go runtime. The
target machine needs nothing: no Go, no package manager, no shared libraries.

#bridge[
  Deploying a Next.js app means Node of the right major version, a
  `node_modules` tree, a build step, and a process manager. Deploying this
  means copying one file. That is the whole difference, and it is the reason Go
  dominates CLI tools and containers.
]

== Smaller binaries

The default build includes debug symbols and a symbol table.

```bash
go build -ldflags="-s -w" -o app .
```

`-s` drops the symbol table, `-w` the DWARF debug information. Expect roughly
25 to 30 percent off. The cost is that a stack trace from a panic loses some
detail and the binary cannot be inspected with a debugger, so many teams ship
unstripped and accept the size.

For genuinely small, `upx` compresses further, at the cost of slower startup
and antivirus software treating your binary with suspicion. Rarely worth it.

== Embedding the version

A binary that cannot tell you which commit it is costs real time during an
incident.

```go
package main

var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

func main() {
    if len(os.Args) > 1 && os.Args[1] == "version" {
        fmt.Printf("%s (%s) built %s\n", version, commit, date)
        return
    }
    ...
}
```

```bash
go build -ldflags="\
  -X main.version=1.4.2 \
  -X main.commit=$(git rev-parse --short HEAD) \
  -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o app .
```

`-X` sets a string variable at link time. It only works on `string` variables
in the `main` package or a named package path, and only ones that are not
constants.

#bridge[
  Grit does exactly this, which is why `grit version` reports a real version
  rather than "dev". Its release process bumps three places that must agree:
  the Go variable, the scaffolder's default, and the documentation site's
  constant. Chapter 47 covers why a version that disagrees with itself is
  worse than no version.
]

== Cross compiling

Two environment variables, no toolchain to install.

```bash
GOOS=linux   GOARCH=amd64 go build -o app-linux-amd64 .
GOOS=linux   GOARCH=arm64 go build -o app-linux-arm64 .
GOOS=darwin  GOARCH=arm64 go build -o app-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o app-windows-amd64.exe .
```

From any machine, to any target. This is not an exaggeration: a Windows laptop
builds the Linux ARM binary that runs on your Raspberry Pi, with no container
and no configuration.

```bash
go tool dist list      # every supported pair
```

#gotcha[
  Cross compilation is free *until you use cgo*. Any package binding to a C
  library (the classic `mattn/go-sqlite3`, some image codecs, some crypto)
  needs a C cross-compiler for the target, and the simplicity evaporates.

  The defence is to prefer pure Go packages. This is why the database chapter
  used `modernc.org/sqlite` rather than `mattn/go-sqlite3`: the former is Go
  all the way down, so `GOOS=linux go build` works from anywhere. Check for
  cgo with:

  ```bash
  CGO_ENABLED=0 go build ./...
  ```

  If that fails, something in your tree needs C.
]

== Containers, in two stages

A multi-stage Dockerfile builds in one image and ships in a tiny one.

```dockerfile
# Build
FROM golang:1.25-alpine AS build
WORKDIR /src

# Copy manifests first: this layer caches unless dependencies change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/app .

# Ship
FROM gcr.io/distroless/static-debian12
COPY --from=build /bin/app /app
USER nonroot:nonroot
ENTRYPOINT ["/app"]
```

The final image contains your binary, CA certificates and nothing else: no
shell, no package manager, no libc. Typical size is your binary plus about two
megabytes.

`CGO_ENABLED=0` is required for the `static` base, because there is no libc in
it to dynamically link against.

#gotcha[
  Copy `go.mod` and `go.sum` and run `go mod download` *before* copying your
  source. Docker caches layers, and if the source is copied first then every
  code change re-downloads every dependency. This one ordering decision is the
  difference between a ten second rebuild and a two minute one.
]

== Embedding files into the binary

`go:embed` puts files inside the executable, so your single-file deployment
can still serve a frontend.

```go
import "embed"

//go:embed static/*
var staticFiles embed.FS

//go:embed schema.sql
var schema string

func main() {
    http.Handle("/static/", http.FileServerFS(staticFiles))
    db.Exec(schema)
}
```

The comment is the directive and must sit immediately above the variable with
no blank line. `embed.FS` for a tree, `string` or `[]byte` for one file.

This is how a Go program ships a built React app inside itself: compile the
frontend, embed `dist/`, and deploy one file that serves both the API and the
UI. Grit's single-binary architecture mode does precisely this.

== Production flags worth knowing

```bash
go build -trimpath -ldflags="-s -w" -o app .
```

`-trimpath` removes absolute file paths from the binary, so a panic trace says
`main.go:42` rather than `/home/you/projects/secret-name/main.go:42`. It also
makes builds reproducible, which matters if you sign releases.

```bash
GOAMD64=v3 go build -o app .
```

Targets a newer x86 instruction set for a few percent. Only do this if you
control the hardware, because the binary will not run on older CPUs.

== Automating releases

A tag-triggered workflow that cross compiles and attaches binaries:

```yaml
name: Release
on:
  push:
    tags: ['v*']

jobs:
  release:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
        with: { fetch-depth: 0 }
      - uses: actions/setup-go@v5
        with: { go-version: '1.25' }
      - run: go test -race ./...
      - uses: goreleaser/goreleaser-action@v6
        with: { args: release --clean }
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

`goreleaser` reads a config file, builds every platform, writes checksums,
generates release notes from commits and uploads everything. It is the standard
tool, and it is what turns "I tagged v1.2.0" into binaries for six platforms
without further thought.

#project(title: "Project 23: ship the shortener")[
  Take the URL shortener and make it genuinely deployable: versioned,
  containerised and cross compiled.

  First, the version:

  ```go
  // version.go
  package main

  import (
      "fmt"
      "runtime"
  )

  var (
      version = "dev"
      commit  = "none"
      date    = "unknown"
  )

  func versionString() string {
      return fmt.Sprintf("shortener %s (%s) built %s with %s",
          version, commit, date, runtime.Version())
  }
  ```

  Wire it into `main` before anything else happens:

  ```go
  if len(os.Args) > 1 && os.Args[1] == "version" {
      fmt.Println(versionString())
      return
  }
  ```

  A `Makefile`, because the build command is now too long to retype:

  ```makefile
  VERSION ?= $(shell git describe --tags --always --dirty)
  COMMIT  := $(shell git rev-parse --short HEAD)
  DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
  LDFLAGS := -s -w \
      -X main.version=$(VERSION) \
      -X main.commit=$(COMMIT) \
      -X main.date=$(DATE)

  .PHONY: build test release clean

  build:
  	go build -trimpath -ldflags="$(LDFLAGS)" -o bin/shortener .

  test:
  	go test -race -cover ./...

  release:
  	GOOS=linux   GOARCH=amd64 go build -trimpath \
  	    -ldflags="$(LDFLAGS)" -o bin/shortener-linux-amd64 .
  	GOOS=linux   GOARCH=arm64 go build -trimpath \
  	    -ldflags="$(LDFLAGS)" -o bin/shortener-linux-arm64 .
  	GOOS=darwin  GOARCH=arm64 go build -trimpath \
  	    -ldflags="$(LDFLAGS)" -o bin/shortener-darwin-arm64 .
  	GOOS=windows GOARCH=amd64 go build -trimpath \
  	    -ldflags="$(LDFLAGS)" -o bin/shortener-windows-amd64.exe .

  clean:
  	rm -rf bin
  ```

  The `Dockerfile`:

  ```dockerfile
  FROM golang:1.25-alpine AS build
  WORKDIR /src
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .
  ARG VERSION=docker
  RUN CGO_ENABLED=0 go build -trimpath \
      -ldflags="-s -w -X main.version=$VERSION" -o /bin/app .

  FROM gcr.io/distroless/static-debian12
  COPY --from=build /bin/app /app
  USER nonroot:nonroot
  EXPOSE 8080
  ENV DB_PATH=/data/links.db
  ENTRYPOINT ["/app"]
  ```

  Build and check the results:

  ```bash
  make build && ./bin/shortener version
  make release && ls -lh bin/

  docker build --build-arg VERSION=1.0.0 -t shortener .
  docker images shortener          # a few megabytes, not a few hundred
  docker run -p 8080:8080 -v $(pwd)/data:/data shortener
  ```

  Three things to notice. `make release` produces four binaries for four
  platforms from one machine, in seconds, with no cross-compilation toolchain,
  because nothing in this project uses cgo. The container has no shell at all,
  so `docker exec -it ... sh` fails: that is the point of distroless, and it is
  worth experiencing once so the tradeoff is a decision rather than a surprise.
  And the whole image is small enough that pulling it is not a deployment
  consideration.

  Compare the image size to a Node service doing the same job. That number is
  the argument of this chapter.
]

#tryit[
  Add `go:embed` for a small HTML page at `/` that posts to `/shorten`, so the
  single binary serves a usable interface as well as an API. Then add a
  `/healthz` that reports the version, and confirm with `docker run` that the
  version baked in by the build argument is what the endpoint reports. That
  round trip, from build flag to running container to HTTP response, is the
  whole deployment story in one check.
]

#assignment(name: "Shipping - Assignment One")[
  Ship the shortener as if somebody else had to run it.

  *One.* Build with the version embedded, and have `-version` print it:

  ```bash
  go build -ldflags "-X main.version=$(git describe --always)" .
  ```

  *Two.* Cross compile for Linux, macOS and Windows, and record the three
  sizes in `RELEASE.md`.

  *Three.* A two-stage `Dockerfile` ending in a minimal image. Record the image
  size next to the binary sizes.

  *Four.* Embed the static files with `embed`, so one file is genuinely the
  whole program.

  *Five.* A GitHub Actions workflow that does all of it on a tag and attaches
  the binaries to the release. Push a tag and let it run. That workflow file,
  in your repository, is the assignment.]
