package plugin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// pkgJSON writes a package.json with the given dependency lines.
func pkgJSON(t *testing.T, deps ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "package.json")
	body := "{\n  \"name\": \"admin\",\n  \"dependencies\": {\n"
	for i, d := range deps {
		body += "    " + d
		if i < len(deps)-1 {
			body += ","
		}
		body += "\n"
	}
	body += "  }\n}\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func mustParse(t *testing.T, path string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("%s no longer parses as JSON:\n%s\n%v", path, raw, err)
	}
	return out
}

// Removing a dependency leaves JSON that parses, from the middle of the list.
func TestRemoveNodeDependencyFromTheMiddle(t *testing.T) {
	path := pkgJSON(t,
		`"react": "19.1.0"`,
		`"calendarkit-pro": "2.0.0"`,
		`"zod": "3.24.1"`,
	)

	removed, err := removeNodeDependency(path, "calendarkit-pro", "2.0.0")
	if err != nil || !removed {
		t.Fatalf("removed = %v, err = %v", removed, err)
	}
	deps := mustParse(t, path)["dependencies"].(map[string]any)
	if _, still := deps["calendarkit-pro"]; still {
		t.Error("the dependency is still there")
	}
	if len(deps) != 2 || deps["react"] != "19.1.0" || deps["zod"] != "3.24.1" {
		t.Errorf("the neighbours changed: %v", deps)
	}
}

// And from the end, where the comma that joined it dangles.
//
// This is the case that produces a package.json pnpm refuses to read, so it is
// asserted by parsing rather than by matching text.
func TestRemoveNodeDependencyFromTheEnd(t *testing.T) {
	path := pkgJSON(t,
		`"react": "19.1.0"`,
		`"calendarkit-pro": "2.0.0"`,
	)

	removed, err := removeNodeDependency(path, "calendarkit-pro", "2.0.0")
	if err != nil || !removed {
		t.Fatalf("removed = %v, err = %v", removed, err)
	}
	deps := mustParse(t, path)["dependencies"].(map[string]any)
	if len(deps) != 1 || deps["react"] != "19.1.0" {
		t.Errorf("deps = %v", deps)
	}
}

// The only dependency, which empties the object.
func TestRemoveTheOnlyNodeDependency(t *testing.T) {
	path := pkgJSON(t, `"calendarkit-pro": "2.0.0"`)

	removed, err := removeNodeDependency(path, "calendarkit-pro", "2.0.0")
	if err != nil || !removed {
		t.Fatalf("removed = %v, err = %v", removed, err)
	}
	deps := mustParse(t, path)["dependencies"].(map[string]any)
	if len(deps) != 0 {
		t.Errorf("deps = %v", deps)
	}
}

// A version somebody has changed is left alone.
//
// The same rule an edited injection gets: the project has made a decision, and
// a removal that overrides it is a removal that deletes somebody's work. The
// caller turns the false into a warning naming the package.
func TestADifferentVersionIsNotRemoved(t *testing.T) {
	path := pkgJSON(t, `"calendarkit-pro": "2.1.0"`)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	removed, err := removeNodeDependency(path, "calendarkit-pro", "2.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if removed {
		t.Error("a version the plugin did not install was removed")
	}
	after, _ := os.ReadFile(path)
	if string(before) != string(after) {
		t.Error("the file was rewritten anyway")
	}
}

// A package that is not there is not an error.
func TestRemovingAnAbsentDependencyIsQuiet(t *testing.T) {
	path := pkgJSON(t, `"react": "19.1.0"`)
	removed, err := removeNodeDependency(path, "calendarkit-pro", "2.0.0")
	if err != nil || removed {
		t.Errorf("removed = %v, err = %v", removed, err)
	}
}

// add then remove returns the file to its original bytes.
//
// The pair is what matters: either half on its own can look right while the
// round trip leaves a comma, a blank line or a changed indent behind.
func TestAddThenRemoveIsARoundTrip(t *testing.T) {
	for _, start := range [][]string{
		{`"react": "19.1.0"`},
		{`"react": "19.1.0"`, `"zod": "3.24.1"`},
	} {
		path := pkgJSON(t, start...)
		before, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		added, err := addNodeDependency(path, "calendarkit-pro", "2.0.0")
		if err != nil || !added {
			t.Fatalf("added = %v, err = %v", added, err)
		}
		mustParse(t, path) // still JSON with the entry in

		removed, err := removeNodeDependency(path, "calendarkit-pro", "2.0.0")
		if err != nil || !removed {
			t.Fatalf("removed = %v, err = %v", removed, err)
		}

		after, _ := os.ReadFile(path)
		if string(before) != string(after) {
			t.Errorf("the round trip changed the file:\nbefore:\n%s\nafter:\n%s", before, after)
		}
	}
}

// Remove reverts the dependencies it recorded.
//
// Guarded at the source, because the behaviour needs a whole project to
// exercise and the mistake is a missing loop rather than a wrong one: every
// file and patch came back and the npm dependency stayed, which is how the
// first community plugin left calendarkit-pro behind in package.json.
func TestRemoveRevertsTheRecordedDependencies(t *testing.T) {
	body := read(t, "install.go")
	at := strings.Index(body, "func Remove(")
	if at < 0 {
		t.Fatal("no Remove")
	}
	removeBody := body[at:]
	if end := strings.Index(removeBody, "\n// injectBefore"); end > 0 {
		removeBody = removeBody[:end]
	}

	if !strings.Contains(removeBody, "record.NodeDeps") {
		t.Error("Remove does not walk the recorded node dependencies, so a " +
			"removed plugin leaves its npm packages in package.json and " +
			"pnpm install keeps fetching them")
	}
	if !strings.Contains(removeBody, "removeNodeDependency(") {
		t.Error("Remove records the dependencies and does not take them out")
	}
	if !strings.Contains(removeBody, "record.GoDeps") {
		t.Error("Remove says nothing about the Go modules it added, so an " +
			"unused require stays in go.mod with nothing pointing at it")
	}
}
