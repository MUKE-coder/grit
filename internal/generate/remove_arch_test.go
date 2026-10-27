package generate

import (
	"path/filepath"
	"testing"
)

// `grit remove resource` spelled two of its roots as apps/api and apps/web.
//
// A single keeps its Go code in api/ (or at the project root) and IS the
// frontend, so on one of those the removal deleted only the two files whose
// roots were resolved properly, left the model, the handler, the service, the
// seeder and the hook exactly where they were, and printed "removed
// successfully". Ten files on a single project, two of them deleted.
//
// Found while fixing grit#111, which is the same fault from the other side:
// there an upgrade put back what removal had taken out, here removal never
// took it out at all.
func TestRemoveResolvesRootsPerArchitecture(t *testing.T) {
	root := filepath.Join("some", "project")

	for _, c := range []struct {
		arch string
		web  string
	}{
		// A single is its own frontend: its pages are at the project root.
		{"single", root},
		// A double has one frontend with the admin as a route group inside it,
		// and a triple has a separate admin. Both keep the web app in the same
		// place.
		{"double", filepath.Join(root, "apps", "web")},
		{"triple", filepath.Join(root, "apps", "web")},
	} {
		if got := removeWebRoot(root, c.arch); got != c.web {
			t.Errorf("%s: web root is %s, want %s", c.arch, got, c.web)
		}
	}
}

// And the API root comes from the generator, which already knows, rather than
// being spelled out a second time.
func TestRemoveUsesTheGeneratorsAPIRoot(t *testing.T) {
	root := t.TempDir()

	single := (&Generator{Root: root, Architecture: "single"}).APIRoot()
	// No api/go.mod in the temp dir, so a single's Go code is at its root.
	if single != root {
		t.Errorf("single API root is %s, want the project root %s", single, root)
	}

	triple := (&Generator{Root: root, Architecture: "triple"}).APIRoot()
	if want := filepath.Join(root, "apps", "api"); triple != want {
		t.Errorf("triple API root is %s, want %s", triple, want)
	}
}
