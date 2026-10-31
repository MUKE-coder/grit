package plugin

import "strings"

// The skeleton's prose, kept apart from its code because a Markdown fence is
// three backticks and a Go raw string cannot hold one. Spliced rather than
// escaped: an interpreted string of this length is unreadable.
const fence = "```"

// skeletonReadme is the README, and it is most of the value of the skeleton.
//
// The issue asked for the full loop, and the loop is the part that cannot be
// inferred from the manifest: install, build, type-check, remove, build again.
// That last build is what catches a plugin which edited a file outside an
// injection, because removal replays the lockfile and cannot take back an edit
// it has no record of.
func skeletonReadme(name, pkg string) string {
	return `# ` + name + `

A Grit plugin. Installing it writes a Go package, an admin component and two
injections into a project; removing it takes all of that back out.

` + fence + `bash
grit plugin add ./` + name + `
grit plugin remove ` + name + `
` + fence + `

## What is in here

` + fence + `text
plugin.json                      the manifest
files/api/` + pkg + `.go` + padTo(pkg, 24) + `  the endpoint, installed always
files/api/` + pkg + `_test.go` + padTo(pkg, 19) + `  its test, which runs in the project
files/admin/` + name + `-card.tsx` + padTo(name, 15) + `  installed only in a triple
.github/workflows/plugin.yml     the install and remove loop, in CI
` + fence + `

## The loop

Run this before you publish anything, and let CI run it on every push. The
last two commands are the ones that matter: removal replays the lockfile
backwards, so a plugin that edited a file outside an injection installs
cleanly and leaves the project broken when it comes out.

` + fence + `bash
grit new testbed --triple
cd testbed
grit plugin add ../` + name + `

cd apps/api && go build ./... && go test ./...
cd ../.. && pnpm install && cd apps/admin && npx tsc --noEmit

cd ../.. && grit plugin remove ` + name + ` --force
cd apps/api && go build ./...
` + fence + `

## The manifest

Required: ` + "`name`" + `, ` + "`version`" + `, ` + "`summary`" + `, and at least one file or
injection. A plugin with neither does nothing, and ` + "`grit plugin add`" + ` says so
rather than reporting success.

### Placeholders

File bodies and the ` + "`to`" + ` and ` + "`file`" + ` paths are substituted, so a plugin
does not need to know anything about the project installing it:

| Placeholder | Becomes |
|---|---|
| ` + "`{{MODULE}}`" + ` | the Go module path, e.g. ` + "`myapp/apps/api`" + ` |
| ` + "`{{PROJECT}}`" + ` | the project directory's name |
| ` + "`{{API_ROOT}}`" + ` | ` + "`apps/api`" + `, or ` + "`.`" + ` in a ` + "`--single`" + ` project |

Use ` + "`{{API_ROOT}}`" + ` for every Go path. Writing ` + "`apps/api/...`" + ` by hand
works until somebody installs the plugin into a single project, where the Go
module is at the root.

### when

A file or an injection can be limited to certain projects. It is data, not an
expression, and leaving it out means always:

` + fence + `json
"when": { "architecture": ["triple", "double"], "frontend": ["next"] }
` + fence + `

` + "`architecture`" + ` is ` + "`single`" + `, ` + "`double`" + `, ` + "`triple`" + `, ` + "`api`" + ` or
` + "`mobile`" + `. ` + "`frontend`" + ` is ` + "`next`" + ` or ` + "`tanstack`" + `.

The admin file here is gated to ` + "`triple`" + ` on purpose. A standalone
` + "`apps/admin`" + ` imports its own code as ` + "`@/components/...`" + `; a panel
embedded in a web app keeps the same components under ` + "`admin-panel/`" + ` and
imports them as ` + "`@admin/components/...`" + `. One file cannot satisfy both, so
supporting a double means a second entry with its own ` + "`when`" + `:

` + fence + `json
{
  "from": "files/admin/` + name + `-card.embedded.tsx",
  "to": "apps/web/admin-panel/components/` + name + `/` + name + `-card.tsx",
  "when": { "architecture": ["double"], "frontend": ["next"] }
}
` + fence + `

### Injections

An injection inserts a snippet on the line **before** a marker comment that is
already in the file. The two here are a pair: the route needs the import, and
an import with no route does not compile either.

` + fence + `text
// grit:imports              inside the import block of routes.go
// grit:routes:protected     a route requiring a signed-in user
// grit:routes:public        a route requiring nothing
// grit:routes:admin         a route requiring the admin role
// grit:routes:custom        function level, with the route groups in scope
// grit:models               a model to migrate
// grit:handlers             a handler to construct
` + fence + `

A marker that is not in the file is not an error at install time: the
injection simply does not happen, and the first sign is a route that answers
404. So check the marker against a real project rather than against this list.

Snippets can also come from a file, which is easier to read for anything
longer than a line:

` + fence + `json
{ "file": "{{API_ROOT}}/internal/routes/routes.go",
  "marker": "// grit:handlers",
  "codeFrom": "files/api/wiring.go.txt" }
` + fence + `

Set ` + "`\"optional\": true`" + ` for an injection into a file that may not exist, so
a project without that feature installs the rest of the plugin instead of
failing.

### Dependencies

` + fence + `json
"goDeps":   [{ "name": "github.com/x/y", "version": "v1.2.3" }],
"nodeDeps": [{ "name": "recharts", "version": "^2.12.0", "workspace": "apps/admin" }]
` + fence + `

Leave ` + "`workspace`" + ` out to target the Go API. Pin versions: an install that
resolves differently next week is an install nobody can reproduce.

## What not to do

- **Do not overwrite a file the project already has.** Install refuses, and it
  is right to. If your plugin needs to change an existing file, inject at a
  marker.
- **Do not anchor an injection on ordinary code.** Markers exist so the file
  can be reformatted and edited without breaking you. Anchoring on
  ` + "`func Register(`" + ` works until somebody adds a comment above it.
- **Do not forget ` + "`when`" + `.** An admin file written into an ` + "`--api`" + `
  project lands in a directory nothing compiles, which nobody notices until a
  build somewhere else.
- **Do not climb out of the project.** ` + "`\"to\": \"../../.ssh/authorized_keys\"`" + `
  is refused at load time, by design.

## Publishing

Bump ` + "`version`" + ` in ` + "`plugin.json`" + `, add a ` + "`CHANGELOG.md`" + ` entry, tag
the repository, then:

` + fence + `bash
grit plugin submit
` + fence + `

which opens the submission issue against the community index with the fields
already filled in.
`
}

// padTo lines up the tree listing without a table, since the names vary in
// length with the plugin's own name.
func padTo(s string, width int) string {
	if n := width - len(s); n > 0 {
		return strings.Repeat(" ", n)
	}
	return ""
}

func skeletonChangelog() string {
	return `# Changelog

## 0.1.0

First release.
`
}

// skeletonWorkflow runs the loop from the README on both of the two shapes
// worth covering: the monorepo the plugin's admin file targets, and the
// API-only project where that file must not be written at all.
func skeletonWorkflow(name string) string {
	return `name: plugin

on:
  push:
  pull_request:

jobs:
  install-and-remove:
    runs-on: ubuntu-latest
    strategy:
      fail-fast: false
      matrix:
        # Both shapes, because the when clause behaves differently in each:
        # the admin file is written for triple and skipped for api. An api
        # project is also the one where a missing when clause would put a
        # .tsx file somewhere nothing compiles.
        include:
          - arch: "--triple"
            frontend: true
          - arch: "--api"
            frontend: false
    steps:
      - uses: actions/checkout@v4
        with:
          path: plugin

      - uses: actions/setup-go@v5
        with:
          go-version: "1.26"

      - name: Install the Grit CLI
        run: go install github.com/MUKE-coder/grit/v3/cmd/grit@latest

      - name: Scaffold a project
        run: grit new testbed ${{ matrix.arch }}

      - name: Install the plugin
        working-directory: testbed
        run: grit plugin add ../plugin

      - name: The Go side builds and its tests pass
        working-directory: testbed/apps/api
        run: |
          go mod tidy
          go build ./...
          go test ./...

      - uses: pnpm/action-setup@v4
        if: matrix.frontend
        with:
          version: 10.33.4
      - uses: actions/setup-node@v4
        if: matrix.frontend
        with:
          node-version: "22"

      - name: The admin type-checks
        if: matrix.frontend
        working-directory: testbed
        run: |
          pnpm install
          cd apps/admin && npx tsc --noEmit

      # The step that catches a plugin which edited a file outside an
      # injection: removal replays the lockfile, so it cannot take back an
      # edit it has no record of, and the build is where that shows up.
      - name: Remove it, and the project still builds
        working-directory: testbed
        run: |
          grit plugin remove ` + name + ` --force
          cd apps/api && go build ./...
`
}
