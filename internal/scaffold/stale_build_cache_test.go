package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

// An upgrade rewrites route files, so the compiled output of those files is
// stale the moment it finishes.
//
// Next keeps its build in .next, keyed to the route tree it compiled. A dev
// server running through an upgrade holds a client tree describing the routes
// as they were, and the App Router does not always recover: its layout router
// reads the render tree for the level it is drawing and finds undefined, which
// surfaces as "Cannot read properties of undefined (reading 'slots')" from a
// stack with no sign of which file caused it.
func TestUpgradeClearsStaleBuildCaches(t *testing.T) {
	root := t.TempDir()
	caches := []string{
		filepath.Join(root, "apps", "web", ".next"),
		filepath.Join(root, "apps", "admin", ".next"),
		filepath.Join(root, "apps", "docs", ".next"),
		filepath.Join(root, "node_modules", ".vite"),
	}
	for _, dir := range caches {
		if err := os.MkdirAll(filepath.Join(dir, "static"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	// Something that is not a build cache, to be sure the sweep is narrow.
	keep := filepath.Join(root, "apps", "web", "app")
	if err := os.MkdirAll(keep, 0o755); err != nil {
		t.Fatal(err)
	}

	cleared := clearFrontendBuildCaches(root)
	if len(cleared) != len(caches) {
		t.Errorf("cleared %d caches, want %d: %v", len(cleared), len(caches), cleared)
	}
	for _, dir := range caches {
		if dirExists(dir) {
			t.Errorf("%s survived", dir)
		}
	}
	if !dirExists(keep) {
		t.Error("the sweep removed app/, which is the source and not a cache")
	}

	// A project with no caches at all is not an error, and says nothing.
	if got := clearFrontendBuildCaches(t.TempDir()); len(got) != 0 {
		t.Errorf("a project with no build caches reported %v", got)
	}
}
