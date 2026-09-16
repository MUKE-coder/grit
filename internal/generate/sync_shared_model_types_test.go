package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// A new project ships packages/shared/types files for the scaffold's own
// models, so the admin can import them before grit sync has ever run
// (contact-app review M41). They must be what sync would write: otherwise the
// first sync rewrites them and a page that compiled stops compiling.
func TestScaffoldSharedModelTypesMatchSync(t *testing.T) {
	dir := t.TempDir()
	for _, typ := range scaffold.SharedModelTypes {
		model := strings.ReplaceAll(typ.Model(), "{{MODULE}}", "example.com/app")
		path := filepath.Join(dir, typ.Kebab+".go")
		if err := os.WriteFile(path, []byte(model), 0o644); err != nil {
			t.Fatal(err)
		}
		structs, err := parseGoStructs(path)
		if err != nil {
			t.Fatalf("%s: %v", typ.Name, err)
		}
		found := false
		for _, s := range structs {
			if s.Name != typ.Name {
				continue
			}
			found = true
			if got := buildTSType(s); got != typ.Source {
				t.Errorf("packages/shared/types/%s.ts is not what grit sync writes.\nsync:\n%s\nscaffold:\n%s", typ.Kebab, got, typ.Source)
			}
		}
		if !found {
			t.Errorf("no %s struct in its model file", typ.Name)
		}
	}
}
