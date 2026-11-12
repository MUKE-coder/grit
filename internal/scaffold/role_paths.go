package scaffold

import (
	"os"
	"path/filepath"
	"strings"
)

// Where a project's role declarations live, in any architecture.
//
// `grit add role` hardcoded apps/api, packages/shared and apps/admin, which
// are the triple layout and only the triple layout. On a --single project the
// Go module is api/, the admin is a section inside the one SPA, and the
// command found none of its files and refused with "no files were updated".
// So the command was simply unavailable there, and nothing said why.
//
// Six places are resolved rather than assumed, by the same walk every other
// architecture-aware part of the scaffolder uses: try each layout, take the
// first that is there.

// rolePaths are the files a role has to be declared in.
type rolePaths struct {
	Model     string // internal/models/user.go
	Schema    string // packages/shared/schemas/user.ts
	Types     string // packages/shared/types/user.ts
	Constants string // packages/shared/constants/index.ts
	Admin     string // the admin's users resource definition
}

// firstExistingDir returns the first candidate that is a directory, or "".
func firstExistingDir(candidates ...string) string {
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return ""
}

// firstExistingFile returns the first candidate that is a file, or "".
func firstExistingFile(candidates ...string) string {
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && !info.IsDir() {
			return c
		}
	}
	return ""
}

// roleFilesFor locates a project's role declarations.
//
// Every field may be empty: a project without an admin has no users resource
// to put a badge in, and that is a fact about the project rather than a
// failure. The caller reports what it could not find rather than guessing.
func roleFilesFor(root string) rolePaths {
	// The Go module: apps/api on a monorepo, api/ on --single, the root
	// itself on an --api project.
	apiRoot := firstExistingDir(
		filepath.Join(root, "apps", "api"),
		filepath.Join(root, "api"),
		root,
	)

	// The shared package is packages/shared in every shape that has one,
	// including --single, which keeps a packages/ directory.
	sharedRoot := firstExistingDir(
		filepath.Join(root, "packages", "shared"),
	)

	p := rolePaths{
		Model: firstExistingFile(filepath.Join(apiRoot, "internal", "models", "user.go")),
	}
	if sharedRoot != "" {
		p.Schema = firstExistingFile(filepath.Join(sharedRoot, "schemas", "user.ts"))
		p.Types = firstExistingFile(filepath.Join(sharedRoot, "types", "user.ts"))
		p.Constants = firstExistingFile(filepath.Join(sharedRoot, "constants", "index.ts"))
	}

	// The admin's users resource, wherever the panel is: its own app, a route
	// group inside the web app, or a section of the single SPA.
	for _, resources := range []string{
		filepath.Join(root, "apps", "admin", "resources"),
		filepath.Join(root, "apps", "web", "admin-panel", "resources"),
		filepath.Join(root, "apps", "web", "src", "admin-panel", "resources"),
		filepath.Join(root, "src", "admin", "resources"),
		filepath.Join(root, "src", "routes", "admin", "resources"),
		filepath.Join(root, "frontend", "src", "admin-panel", "resources"),
	} {
		if def := FindResourceDef(resources, "users"); def != "" {
			p.Admin = def
			break
		}
	}
	return p
}

// missing names the role files a project does not have, for a message that
// says what was not updated rather than a count that is quietly short.
func (p rolePaths) missing(root string) []string {
	var out []string
	for _, f := range []struct{ path, what string }{
		{p.Model, "internal/models/user.go"},
		{p.Schema, "packages/shared/schemas/user.ts"},
		{p.Types, "packages/shared/types/user.ts"},
		{p.Constants, "packages/shared/constants/index.ts"},
		{p.Admin, "the admin's users resource"},
	} {
		if f.path == "" {
			out = append(out, f.what)
		}
	}
	return out
}

// shown renders a path relative to the project root, for the terminal.
func shown(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return strings.ReplaceAll(rel, string(filepath.Separator), "/")
}
