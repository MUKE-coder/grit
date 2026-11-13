package doctor

import (
	"os"
	"path/filepath"
	"testing"
)

// doctor could not be run on a single project at all.
//
// apiRootOf knew apps/api and the project root. A --single project keeps its
// Go module in api/ with the frontend owning the root, so loadProject looked
// for <root>/internal, did not find it, and every run ended with
//
//	no Go API at <root>: run grit doctor inside a Grit project
//
// The command whose whole job is finding problems was the one that could not
// run. The seventh place this same resolver was written out by hand: v3.516.0
// fixed six others in the generator.
func TestAPIRootFindsEveryLayout(t *testing.T) {
	withGoMod := func(t *testing.T, rel ...string) string {
		t.Helper()
		root := t.TempDir()
		dir := root
		if len(rel) > 0 {
			dir = filepath.Join(append([]string{root}, rel...)...)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module app\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		return root
	}

	mono := withGoMod(t, "apps", "api")
	if got, want := apiRootOf(mono), filepath.Join(mono, "apps", "api"); got != want {
		t.Errorf("monorepo API root is %s, want %s", got, want)
	}

	single := withGoMod(t, "api")
	if got, want := apiRootOf(single), filepath.Join(single, "api"); got != want {
		t.Errorf("single API root is %s, want %s", got, want)
	}

	// A single scaffolded before the module moved into api/ keeps everything
	// at the root, and still has to be readable.
	legacy := withGoMod(t)
	if got := apiRootOf(legacy); got != legacy {
		t.Errorf("legacy single API root is %s, want the project root %s", got, legacy)
	}
}

// And the error that said doctor could not run has to stop happening for the
// layout it could not read. loadProject is the thing that printed it.
func TestDoctorLoadsASingleProject(t *testing.T) {
	root := t.TempDir()
	internal := filepath.Join(root, "api", "internal", "models")
	if err := os.MkdirAll(internal, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "api", "go.mod"), []byte("module app\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	p, err := loadProject(root)
	if err != nil {
		t.Fatalf("doctor refused a single project: %v", err)
	}
	if want := filepath.Join(root, "api"); p.apiRoot != want {
		t.Errorf("read the API from %s, want %s", p.apiRoot, want)
	}
}
