package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

// Every shape Grit scaffolds has to be findable by the thing that installs
// i18n into it.
//
// grit new --i18n runs AddI18n straight after the scaffold, so a shape
// detectLayout cannot resolve does not degrade: it scaffolds the entire
// project, then fails with "no Go API found", pointing at the module it wrote
// itself a second earlier. --single and --api --admin both keep the module in
// api/, and neither was in the candidate list.
func TestEveryProjectShapeHasAFindableAPI(t *testing.T) {
	shapes := map[string]string{
		"triple / double / mobile": filepath.Join("apps", "api"),
		"single":                   "api",
		"--api --admin":            "api",
		"api-only":                 "",
	}

	for name, apiDir := range shapes {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			want := filepath.Join(root, apiDir)
			if err := os.MkdirAll(want, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(want, "go.mod"),
				[]byte("module example.com/app\n"), 0o644); err != nil {
				t.Fatal(err)
			}

			layout, err := detectLayout(root)
			if err != nil {
				t.Fatalf("a %s project cannot be resolved: %v", name, err)
			}
			if layout.APIRoot != want {
				t.Errorf("a %s project resolved to %q, want %q", name, layout.APIRoot, want)
			}
		})
	}
}

// A monorepo that also happens to hold an api/ directory still resolves to
// apps/api, because the candidates are ordered most specific first.
func TestTheMonorepoAPIWinsOverAStrayAPIDirectory(t *testing.T) {
	root := t.TempDir()

	mono := filepath.Join(root, "apps", "api")
	stray := filepath.Join(root, "api")
	for _, dir := range []string{mono, stray} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "go.mod"),
			[]byte("module example.com/app\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	layout, err := detectLayout(root)
	if err != nil {
		t.Fatal(err)
	}
	if layout.APIRoot != mono {
		t.Errorf("resolved to %q, want the monorepo API at %q", layout.APIRoot, mono)
	}
}

// And a directory called api with no module in it is not an API.
func TestAnEmptyAPIDirectoryIsNotMistakenForTheAPI(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "api"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"),
		[]byte("module example.com/app\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	layout, err := detectLayout(root)
	if err != nil {
		t.Fatal(err)
	}
	if layout.APIRoot != root {
		t.Errorf("resolved to %q, want the project root %q", layout.APIRoot, root)
	}
}
