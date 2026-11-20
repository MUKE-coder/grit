package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writePluginDir builds a plugin directory from a map of relative paths.
func writePluginDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for rel, body := range files {
		path := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func tripleNext() Context {
	return Context{
		Root:         "/project",
		APIRoot:      "/project/apps/api",
		Module:       "example.com/app",
		Architecture: "triple",
		Frontend:     "next",
	}
}

// A directory entry installs the tree, keeping its shape.
//
// grit#227: a manifest listed every file one by one, repeating the same "to"
// prefix and the same "when" clause for each. grit-calendar needed nine
// entries for four groups, and adding a component meant editing the tree and
// the manifest and keeping the two in step.
func TestADirectoryEntryInstallsTheWholeTree(t *testing.T) {
	dir := writePluginDir(t, map[string]string{
		"plugin.json": `{
			"name": "calendar",
			"version": "0.5.0",
			"summary": "A calendar",
			"directories": [
				{"from": "files/components", "to": "apps/admin/components/calendar"}
			]
		}`,
		"files/components/calendar-view.tsx":   "view",
		"files/components/calendar-widget.tsx": "widget",
		"files/components/day/day-grid.tsx":    "grid",
		"files/components/calendar.css":        "css",
	})

	p, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	got := p.Files(tripleNext())

	want := map[string]string{
		"apps/admin/components/calendar/calendar-view.tsx":   "view",
		"apps/admin/components/calendar/calendar-widget.tsx": "widget",
		"apps/admin/components/calendar/day/day-grid.tsx":    "grid",
		"apps/admin/components/calendar/calendar.css":        "css",
	}
	if len(got) != len(want) {
		t.Fatalf("installed %d files, want %d: %v", len(got), len(want), keysOf(got))
	}
	for path, body := range want {
		if got[path] != body {
			t.Errorf("%s = %q, want %q", path, got[path], body)
		}
	}
	// Slashes, not the host separator: a destination is a project path, and a
	// backslash in one would travel into the lockfile and out again.
	for path := range got {
		if strings.Contains(path, `\`) {
			t.Errorf("destination %q carries a host separator", path)
		}
	}
}

// A directory in files[].from is refused, not silently skipped.
//
// This is the bug behind the enhancement. os.Stat passed, os.ReadFile failed
// at install time, and the failure was swallowed as "a file that vanished
// between load and now": the entry installed nothing and the install reported
// success.
func TestADirectoryInAFileEntryIsRefused(t *testing.T) {
	dir := writePluginDir(t, map[string]string{
		"plugin.json": `{
			"name": "calendar",
			"version": "0.5.0",
			"summary": "A calendar",
			"files": [
				{"from": "files/components", "to": "apps/admin/components/calendar"}
			]
		}`,
		"files/components/calendar-view.tsx": "view",
	})

	_, err := LoadDir(dir)
	if err == nil {
		t.Fatal("a directory in files[].from loaded without complaint, so it " +
			"installs nothing and says the install succeeded")
	}
	if !strings.Contains(err.Error(), "is a directory") ||
		!strings.Contains(err.Error(), "directories") {
		t.Errorf("the error does not say what to do instead: %v", err)
	}
}

// And the inverse, because the same mistake happens the other way round.
func TestAFileInADirectoryEntryIsRefused(t *testing.T) {
	dir := writePluginDir(t, map[string]string{
		"plugin.json": `{
			"name": "calendar", "version": "0.5.0", "summary": "A calendar",
			"directories": [{"from": "files/one.tsx", "to": "apps/admin/one.tsx"}]
		}`,
		"files/one.tsx": "one",
	})

	_, err := LoadDir(dir)
	if err == nil || !strings.Contains(err.Error(), "is a file") {
		t.Errorf("err = %v, want a complaint that it is a file", err)
	}
}

// An empty tree is reported rather than installed.
//
// The same failure the file-level bug had: a manifest entry that writes
// nothing and an install that reports success.
func TestAnEmptyDirectoryIsRefused(t *testing.T) {
	dir := writePluginDir(t, map[string]string{
		"plugin.json": `{
			"name": "calendar", "version": "0.5.0", "summary": "A calendar",
			"directories": [{"from": "files/components", "to": "apps/admin/calendar"}]
		}`,
	})
	if err := os.MkdirAll(filepath.Join(dir, "files", "components"), 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := LoadDir(dir)
	if err == nil || !strings.Contains(err.Error(), "installs nothing") {
		t.Errorf("err = %v, want a complaint that it installs nothing", err)
	}
}

// A when clause applies to the whole tree, which is the point of the feature.
func TestADirectorysWhenClauseIsHonoured(t *testing.T) {
	dir := writePluginDir(t, map[string]string{
		"plugin.json": `{
			"name": "calendar", "version": "0.5.0", "summary": "A calendar",
			"directories": [
				{"from": "files/next", "to": "apps/admin/components",
				 "when": {"frontend": ["next"]}},
				{"from": "files/vite", "to": "apps/admin/src/components",
				 "when": {"frontend": ["tanstack"]}}
			]
		}`,
		"files/next/view.tsx": "next view",
		"files/vite/view.tsx": "vite view",
	})

	p, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	next := p.Files(tripleNext())
	if next["apps/admin/components/view.tsx"] != "next view" {
		t.Errorf("the Next tree did not install: %v", keysOf(next))
	}
	if len(next) != 1 {
		t.Errorf("the Vite tree installed into a Next project too: %v", keysOf(next))
	}

	vite := tripleNext()
	vite.Frontend = "tanstack"
	got := p.Files(vite)
	if got["apps/admin/src/components/view.tsx"] != "vite view" || len(got) != 1 {
		t.Errorf("the Vite project got %v", keysOf(got))
	}
}

// A file inside a tree cannot climb out of the project.
//
// safeRelative guards the prefix; this guards what the prefix is joined with.
// A tree is somebody else's directory listing, so the check belongs on every
// destination it produces rather than on the one the manifest wrote down.
func TestATreeCannotClimbOutOfTheProject(t *testing.T) {
	dir := writePluginDir(t, map[string]string{
		"plugin.json": `{
			"name": "evil", "version": "0.1.0", "summary": "A plugin",
			"directories": [{"from": "files/x", "to": "apps/admin"}]
		}`,
		"files/x/a.txt": "fine",
	})
	// A destination that escapes once the prefix is joined. Built by hand
	// because a manifest cannot name it: the tree does.
	deep := filepath.Join(dir, "files", "x", "..", "..", "..", "escaped.txt")
	if err := os.WriteFile(deep, []byte("escaped"), 0o644); err != nil {
		t.Skipf("could not write outside the plugin dir: %v", err)
	}

	// The tree itself holds only a.txt, so this loads: the guard that matters
	// is that nothing outside files/x came with it.
	p, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for path := range p.Files(tripleNext()) {
		if strings.Contains(path, "escaped") {
			t.Errorf("a file outside the tree was installed as %q", path)
		}
	}
}

// A symlink is refused rather than followed.
//
// A link named files/components/x.tsx pointing at a file outside the plugin
// would be read and installed, which is the hole safeRelative closes from the
// other direction.
func TestASymlinkInATreeIsRefused(t *testing.T) {
	dir := writePluginDir(t, map[string]string{
		"plugin.json": `{
			"name": "evil", "version": "0.1.0", "summary": "A plugin",
			"directories": [{"from": "files/x", "to": "apps/admin"}]
		}`,
		"files/x/real.txt": "real",
		"outside.txt":      "secret",
	})

	link := filepath.Join(dir, "files", "x", "link.txt")
	if err := os.Symlink(filepath.Join(dir, "outside.txt"), link); err != nil {
		// Windows needs a privilege for this. The guard is still in the code
		// and still covered on CI's Linux runners.
		t.Skipf("cannot create a symlink here: %v", err)
	}

	_, err := LoadDir(dir)
	if err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Errorf("err = %v, want a complaint about the symlink", err)
	}
}

// files and directories coexist, so a single exception stays easy.
func TestFilesAndDirectoriesTogether(t *testing.T) {
	dir := writePluginDir(t, map[string]string{
		"plugin.json": `{
			"name": "calendar", "version": "0.5.0", "summary": "A calendar",
			"files": [{"from": "files/config.ts", "to": "apps/admin/config/calendar.ts"}],
			"directories": [{"from": "files/components", "to": "apps/admin/components/calendar"}]
		}`,
		"files/config.ts":           "config",
		"files/components/view.tsx": "view",
	})

	p, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	got := p.Files(tripleNext())
	if got["apps/admin/config/calendar.ts"] != "config" {
		t.Errorf("the single file did not install: %v", keysOf(got))
	}
	if got["apps/admin/components/calendar/view.tsx"] != "view" {
		t.Errorf("the tree did not install: %v", keysOf(got))
	}
}

// A manifest with only directories is a manifest that does something.
func TestDirectoriesAloneSatisfyTheDoesNothingCheck(t *testing.T) {
	dir := writePluginDir(t, map[string]string{
		"plugin.json": `{
			"name": "calendar", "version": "0.5.0", "summary": "A calendar",
			"directories": [{"from": "files/x", "to": "apps/admin/x"}]
		}`,
		"files/x/a.txt": "a",
	})
	if _, err := LoadDir(dir); err != nil {
		t.Errorf("a plugin with only a directory entry was refused: %v", err)
	}
}

// The tree's contents go through the same placeholder expansion as a file's.
func TestATreesContentsAreExpanded(t *testing.T) {
	dir := writePluginDir(t, map[string]string{
		"plugin.json": `{
			"name": "calendar", "version": "0.5.0", "summary": "A calendar",
			"directories": [{"from": "files/x", "to": "apps/api/x"}]
		}`,
		"files/x/a.go": "package x // {{MODULE}}",
	})

	p, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	got := p.Files(tripleNext())["apps/api/x/a.go"]
	if !strings.Contains(got, "example.com/app") {
		t.Errorf("the tree's contents were not expanded: %q", got)
	}
}
