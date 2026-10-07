package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

// The project's name comes from go.mod, because the module path is what every
// generated import is built from.
//
// It came from the root package.json, with the scope stripped by taking the
// last path segment. A single's root package.json is its FRONTEND's, named
// "@<project>/web", whose last segment is "web". So every --single --next
// project came out of its first `grit upgrade` with every internal import
// rewritten to web/internal/..., a module that does not exist, and an API that
// would not compile. The error named the Go standard library, which is about
// as far from the cause as an error can point.
func TestProjectNameComesFromTheModule(t *testing.T) {
	write := func(path, body string) {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("single, whose root package.json belongs to the frontend", func(t *testing.T) {
		root := t.TempDir()
		write(filepath.Join(root, "package.json"), `{"name": "@contacts/web"}`)
		write(filepath.Join(root, "api", "go.mod"), "module contacts\n\ngo 1.26\n")

		got, err := readProjectName(root)
		if err != nil {
			t.Fatal(err)
		}
		if got != "contacts" {
			t.Errorf("got %q, want contacts: a module called %q rewrites every import to one that does not exist", got, got)
		}
	})

	t.Run("monorepo, whose module carries the apps/api suffix", func(t *testing.T) {
		root := t.TempDir()
		write(filepath.Join(root, "package.json"), `{"name": "contacts"}`)
		write(filepath.Join(root, "apps", "api", "go.mod"), "module contacts/apps/api\n")

		got, _ := readProjectName(root)
		if got != "contacts" {
			t.Errorf("got %q, want contacts", got)
		}
	})

	t.Run("a module published under a URL", func(t *testing.T) {
		root := t.TempDir()
		write(filepath.Join(root, "apps", "api", "go.mod"), "module github.com/acme/shop/apps/api\n")

		got, _ := readProjectName(root)
		if got != "shop" {
			t.Errorf("got %q, want shop", got)
		}
	})

	t.Run("no go.mod at all falls back to package.json", func(t *testing.T) {
		root := t.TempDir()
		write(filepath.Join(root, "package.json"), `{"name": "frontend-only"}`)

		got, _ := readProjectName(root)
		if got != "frontend-only" {
			t.Errorf("got %q, want frontend-only", got)
		}
	})
}
