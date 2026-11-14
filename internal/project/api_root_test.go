package project

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Every shape Grit scaffolds has to be resolvable.
//
// Not a theoretical list: --single and --api --admin keep the module in api/,
// and the two resolvers that did not know that each failed in the vocabulary of
// their own command. grit new --i18n said "no Go API found" about a module it
// had written a second earlier. grit versions said every framework library was
// "absent, not in go.mod", which is also what a library you have not mounted
// looks like, so the wrong answer was indistinguishable from a right one.
func TestAPIRootFindsTheModuleInEveryShape(t *testing.T) {
	for name, dir := range map[string]string{
		"triple":        filepath.Join("apps", "api"),
		"double":        filepath.Join("apps", "api"),
		"mobile":        filepath.Join("apps", "api"),
		"single":        "api",
		"--api --admin": "api",
		"api-only":      "",
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			want := filepath.Join(root, dir)
			writeModule(t, want)

			if got := APIRoot(root); got != want {
				t.Errorf("a %s project resolved to %q, want %q", name, got, want)
			}
			if !HasAPIModule(root) {
				t.Errorf("a %s project reports no Go module", name)
			}
		})
	}
}

// Most specific first: a monorepo that also holds an api/ directory is still a
// monorepo, and resolving to the stray one would point every check at the
// wrong module while finding a real go.mod, so nothing would complain.
func TestTheMonorepoWinsOverAStrayAPIDirectory(t *testing.T) {
	root := t.TempDir()
	mono := filepath.Join(root, "apps", "api")
	writeModule(t, mono)
	writeModule(t, filepath.Join(root, "api"))

	if got := APIRoot(root); got != mono {
		t.Errorf("resolved to %q, want %q", got, mono)
	}
}

// A directory called api with nothing in it is not an API.
func TestAnEmptyAPIDirectoryIsNotTheAPI(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "api"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeModule(t, root)

	if got := APIRoot(root); got != root {
		t.Errorf("resolved to %q, want the project root %q", got, root)
	}
}

// And a directory with no module anywhere is reported as such rather than
// silently resolving to the root, so a caller that needs certainty has it.
func TestHasAPIModuleIsHonestAboutFindingNothing(t *testing.T) {
	root := t.TempDir()

	if HasAPIModule(root) {
		t.Error("a directory with no go.mod anywhere reports a Go module")
	}
	if got := APIRoot(root); got != root {
		t.Errorf("the fallback is %q, want the project root %q", got, root)
	}
}

func writeModule(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module example.com/app\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

// There is one resolver, and the other two call it.
//
// This is the guard on the actual bug. The fix was not "teach detectLayout
// about api/": it was "stop having three answers to one question", because the
// day api/ was added only one of the three learned it and nothing noticed for
// two releases. If a copy is re-inlined, this fails.
//
// Comments are stripped first. Both delegating functions explain in prose what
// their old candidate list used to be, which is the string a check for a
// candidate list looks for: a guard that matches its own explanation passes
// forever.
func TestNothingElseResolvesTheAPIRootForItself(t *testing.T) {
	callers := map[string]string{
		filepath.Join("..", "doctor", "doctor.go"):                  "gritproject.APIRoot(root)",
		filepath.Join("..", "..", "cmd", "grit", "versions_cmd.go"): "project.APIRoot(root)",
		filepath.Join("..", "scaffold", "add_i18n.go"):              "project.APIRoot(root)",
	}

	for path, call := range callers {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("cannot read %s: %v", path, err)
		}
		body := stripGoComments(string(raw))

		if !strings.Contains(body, call) {
			t.Errorf("%s no longer calls %s, so it has its own answer again", path, call)
		}
		// The shape of a hand-rolled copy: the monorepo path joined inside the
		// file that is supposed to be delegating.
		if strings.Contains(body, `filepath.Join(root, "apps", "api")`) {
			t.Errorf("%s joins the monorepo API path itself; call project.APIRoot instead", path)
		}
	}
}

// stripGoComments removes /* ... */ and // ... so a check cannot match the
// prose that explains it.
func stripGoComments(src string) string {
	block := regexp.MustCompile(`(?s)/\*.*?\*/`)
	line := regexp.MustCompile(`(?m)//.*$`)
	return line.ReplaceAllString(block.ReplaceAllString(src, ""), "")
}
