package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writePlugin lays out a plugin directory and returns its path.
func writePlugin(t *testing.T, manifest string, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for rel, body := range files {
		full := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "plugin.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func testCtx(root string) Context {
	return Context{
		Root: root, Module: "shop/apps/api",
		APIRoot:      filepath.Join(root, "apps", "api"),
		Architecture: "triple", Frontend: "next",
	}
}

// A plugin somebody wrote themselves installs like a built-in.
//
// Before this, the authoring page described a plugin.Plugin value and ended by
// pointing at internal/plugin/multitenant.go, which lives in the CLI. A reader
// who followed it had nowhere to put what they had written: writing a plugin
// meant forking Grit and rebuilding the binary, and the page never said so.
func TestLoadDirBuildsAnInstallablePlugin(t *testing.T) {
	dir := writePlugin(t, `{
	  "name": "product-reviews",
	  "version": "1.0.0",
	  "summary": "Customer reviews",
	  "files": [{"from": "files/review.go", "to": "{{API_ROOT}}/internal/models/review.go"}],
	  "injections": [{
	    "file": "{{API_ROOT}}/internal/models/user.go",
	    "marker": "// grit:models",
	    "code": "\t\t&Review{},"
	  }],
	  "nextSteps": ["Run grit migrate"]
	}`, map[string]string{
		"files/review.go": "package models\n\nimport \"{{MODULE}}/internal/ids\"\n",
	})

	p, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}
	if p.Name != "product-reviews" || p.Version != "1.0.0" {
		t.Errorf("got %s v%s", p.Name, p.Version)
	}

	ctx := testCtx(t.TempDir())
	files := p.Files(ctx)
	body, ok := files["apps/api/internal/models/review.go"]
	if !ok {
		t.Fatalf("the file landed nowhere useful: %v", keysOf(files))
	}
	// The plugin cannot know the module path of the project installing it.
	if !strings.Contains(body, `"shop/apps/api/internal/ids"`) {
		t.Errorf("{{MODULE}} was not substituted in the file body:\n%s", body)
	}

	inj := p.Injections(ctx)
	if len(inj) != 1 || inj[0].File != "apps/api/internal/models/user.go" {
		t.Fatalf("injection target is wrong: %+v", inj)
	}
}

// A "to" that climbs out of the project is refused.
//
// A plugin is somebody else's code, and the installer writes wherever it is
// told. "../../.ssh/authorized_keys" is a file write, not an install.
func TestLoadDirRefusesPathsThatEscapeTheProject(t *testing.T) {
	for _, bad := range []string{
		"../../.ssh/authorized_keys",
		"/etc/passwd",
	} {
		dir := writePlugin(t, `{
		  "name": "evil", "version": "1.0.0", "summary": "x",
		  "files": [{"from": "files/x.go", "to": "`+bad+`"}]
		}`, map[string]string{"files/x.go": "package x\n"})

		if _, err := LoadDir(dir); err == nil {
			t.Errorf("a plugin writing to %q was accepted", bad)
		}
	}
}

// A file the manifest names but does not ship is caught before anything is
// written, because a half-applied install is the state removal cannot undo.
func TestLoadDirCatchesAMissingFileUpFront(t *testing.T) {
	dir := writePlugin(t, `{
	  "name": "x", "version": "1.0.0", "summary": "x",
	  "files": [{"from": "files/absent.go", "to": "apps/api/x.go"}]
	}`, nil)

	_, err := LoadDir(dir)
	if err == nil {
		t.Fatal("a manifest naming a file it does not ship was accepted")
	}
	if !strings.Contains(err.Error(), "absent.go") {
		t.Errorf("the error does not name the missing file: %v", err)
	}
}

// when limits a file to the projects that have somewhere to put it.
func TestLoadDirHonoursTheWhenClause(t *testing.T) {
	dir := writePlugin(t, `{
	  "name": "x", "version": "1.0.0", "summary": "x",
	  "files": [
	    {"from": "files/a.go", "to": "apps/api/a.go"},
	    {"from": "files/b.tsx", "to": "apps/admin/b.tsx", "when": {"architecture": ["triple"]}}
	  ]
	}`, map[string]string{"files/a.go": "a", "files/b.tsx": "b"})

	p, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("LoadDir: %v", err)
	}

	triple := testCtx(t.TempDir())
	if len(p.Files(triple)) != 2 {
		t.Errorf("a triple project should get both files, got %v", keysOf(p.Files(triple)))
	}

	apiOnly := testCtx(t.TempDir())
	apiOnly.Architecture = "api"
	got := p.Files(apiOnly)
	if len(got) != 1 {
		t.Errorf("an --api project has no admin to write into, got %v", keysOf(got))
	}
}

// A manifest with nothing in it is a mistake worth naming.
func TestLoadDirRejectsAnEmptyPlugin(t *testing.T) {
	dir := writePlugin(t, `{"name":"x","version":"1.0.0","summary":"x"}`, nil)
	if _, err := LoadDir(dir); err == nil {
		t.Error("a plugin with no files and no injections was accepted")
	}
}

// Only an explicit path is treated as a directory, so a folder cannot shadow
// a built-in that happens to share its name.
func TestOnlyExplicitPathsAreTreatedAsDirectories(t *testing.T) {
	for _, arg := range []string{"./plugins/x", "../x", "/abs/x"} {
		if !IsDirRef(arg) {
			t.Errorf("%q should be read as a path", arg)
		}
	}
	for _, arg := range []string{"multitenant", "impersonate", "product-reviews"} {
		if IsDirRef(arg) {
			t.Errorf("%q should be read as a built-in name", arg)
		}
	}
}

func keysOf(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
