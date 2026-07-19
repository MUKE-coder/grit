// Package plugin provides Grit's extension mechanism.
//
// A Grit plugin is not a runtime-loaded module. Grit is a code generator, so a
// plugin *generates code into your project*: it writes files, injects snippets
// at marker comments, and declares dependencies. What you get is ordinary code
// in your repo that you can read, edit and delete — no vendored magic, and no
// framework indirection at runtime.
//
// The design constraint that shapes everything here: **installation records
// exactly what it did**, and removal replays that record backwards.
//
// Grit already had a feature that generated code one way and removed it via a
// separately hand-maintained list. The two drifted, and `grit remove resource`
// silently left behind handlers, services and switch arms that no longer
// compiled. A plugin system multiplies that risk by every plugin anyone writes,
// so removal here is driven by the lockfile written at install time — never by
// a second list a plugin author has to remember to update.
package plugin

import (
	"fmt"
	"sort"
	"strings"
)

// Context describes the project a plugin is being applied to. Plugins use it to
// emit code that matches the project's module path and shape.
type Context struct {
	// Root is the project root (the directory holding grit.json).
	Root string
	// Module is the Go module path of the API, e.g. "myapp/apps/api".
	Module string
	// APIRoot is where Go code lives — the project root for --single, or
	// apps/api for a monorepo.
	APIRoot string
	// Architecture is single | double | triple | api | mobile.
	Architecture string
	// Frontend is next | tanstack.
	Frontend string
}

// Injection inserts a snippet at a marker comment in an existing file.
//
// Code is inserted on the line BEFORE Marker, matching how the rest of Grit's
// codegen works. The exact text is recorded in the lockfile so removal can take
// out precisely what was put in.
type Injection struct {
	// File is relative to the project root.
	File string
	// Marker is the anchor comment, e.g. "// grit:models".
	Marker string
	// Code is the snippet to insert.
	Code string
	// Optional marks an injection that may legitimately not apply — e.g. a
	// frontend file in an --api project. A missing marker is then a warning
	// rather than a failed install.
	Optional bool
}

// Dependency is a package the plugin needs.
type Dependency struct {
	// Name is the import path (Go) or package name (npm).
	Name string
	// Version is a Go module version or npm semver range.
	Version string
	// Workspace targets an npm dependency at a specific app, e.g. "apps/admin".
	// Empty means the Go API.
	Workspace string
}

// Plugin describes an installable extension.
//
// Files and Injections are functions rather than data so a plugin can adapt to
// the project — a triple app needs admin pages an --api project does not.
type Plugin struct {
	Name    string
	Version string
	// Summary is one line, shown by `grit plugin list`.
	Summary string
	// Description is the longer text for `grit plugin info`.
	Description string

	// Requires names plugins that must be installed first.
	Requires []string

	// Files returns path (relative to project root) -> file content.
	Files func(ctx Context) map[string]string

	// Injections returns edits to existing files.
	Injections func(ctx Context) []Injection

	// GoDeps and NodeDeps are added to go.mod / package.json.
	GoDeps   []Dependency
	NodeDeps []Dependency

	// NextSteps are printed after a successful install — migrations to run,
	// env vars to set, and so on.
	NextSteps []string
}

// registry holds the plugins built into the CLI.
//
// Built-in for now. The Plugin shape is deliberately data-driven so a future
// release can load a manifest from a git repo or registry without changing how
// install and removal work.
var registry = map[string]Plugin{}

// Register adds a plugin to the built-in registry. Called from init() in each
// plugin's file.
func Register(p Plugin) {
	if p.Name == "" {
		panic("plugin: Name is required")
	}
	if _, dup := registry[p.Name]; dup {
		panic("plugin: duplicate registration: " + p.Name)
	}
	registry[p.Name] = p
}

// Get returns a plugin by name.
func Get(name string) (Plugin, error) {
	p, ok := registry[strings.ToLower(strings.TrimSpace(name))]
	if !ok {
		return Plugin{}, fmt.Errorf("unknown plugin %q — run `grit plugin list` to see what's available", name)
	}
	return p, nil
}

// All returns every registered plugin, sorted by name.
func All() []Plugin {
	out := make([]Plugin, 0, len(registry))
	for _, p := range registry {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
