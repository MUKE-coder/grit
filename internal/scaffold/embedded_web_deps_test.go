package scaffold

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// An existing double project whose apps/web/package.json has been edited is left
// alone by the whole-file rewrite, so it never received the two libraries the
// phone and country inputs import, and the build failed on an unresolved import.
// The upgrade now adds only the missing entries, by name.
func TestEmbeddedWebDepsReachAnEditedPackageJSON(t *testing.T) {
	opts := Options{ProjectName: "shop", Architecture: ArchDouble}
	dir := t.TempDir()
	pkg := filepath.Join(dir, "package.json")
	// An edited file: its own xlsx pin, no libphonenumber-js, no Base UI.
	edited := `{
  "name": "@shop/web",
  "dependencies": {
    "next": "^16.1.6",
    "xlsx": "https://cdn.sheetjs.com/xlsx-0.20.3/xlsx-0.20.3.tgz",
    "recharts": "^2.13.0"
  }
}
`
	if err := os.WriteFile(pkg, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}

	ok, err := insertMissingAfter(pkg, `"dependencies": {`, embeddedWebDependencyLines(opts))
	if err != nil || !ok {
		t.Fatalf("insertMissingAfter: ok=%v err=%v", ok, err)
	}
	raw, _ := os.ReadFile(pkg)
	var parsed struct {
		Dependencies map[string]string `json:"dependencies"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("the result is not valid JSON: %v\n%s", err, raw)
	}
	for _, want := range []string{"libphonenumber-js", "@base-ui/react", "react-dropzone"} {
		if parsed.Dependencies[want] == "" {
			t.Errorf("%s was not added", want)
		}
	}
	if got := parsed.Dependencies["recharts"]; got != "^2.13.0" {
		t.Errorf("recharts is %q, want the developer's ^2.13.0 kept", got)
	}
	if n := strings.Count(string(raw), `"recharts"`); n != 1 {
		t.Errorf("recharts appears %d times, want 1", n)
	}

	again, err := insertMissingAfter(pkg, `"dependencies": {`, embeddedWebDependencyLines(opts))
	if err != nil || again {
		t.Errorf("a second run changed the file again (ok=%v err=%v)", again, err)
	}
}

// Only a project with the panel inside the web app has anything to add.
func TestEmbeddedWebDepsAreEmptyWithoutAnEmbeddedPanel(t *testing.T) {
	if lines := embeddedWebDependencyLines(Options{ProjectName: "shop", Architecture: ArchTriple}); len(lines) != 0 {
		t.Errorf("a triple project got %d dependency lines, want none", len(lines))
	}
}
