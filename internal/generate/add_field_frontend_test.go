package generate

import (
	"os"
	"path/filepath"
	"testing"
)

// An --api project has no packages/shared and no apps/admin. grit generate field
// wrote to both without looking, so on an API-only project it failed on the first
// one: the model had already gained the field, the command reported an error, and
// what it had done was left half applied.
func TestAddFieldSkipsAFrontendThatIsNotThere(t *testing.T) {
	root := t.TempDir()
	def, err := ParseInlineFields("Widget", "name:string,colour:text")
	if err != nil {
		t.Fatal(err)
	}
	g := newTestGenerator(root, "app/apps/api", def)
	names := BuildNames(def)
	field := def.Fields[len(def.Fields)-1]

	if err := g.injectFrontendField(names, field); err != nil {
		t.Fatalf("an API-only project should need no frontend injection: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "packages")); !os.IsNotExist(err) {
		t.Error("skipping the frontend should not create packages/")
	}
	if _, err := os.Stat(filepath.Join(root, "apps", "admin")); !os.IsNotExist(err) {
		t.Error("skipping the frontend should not create apps/admin/")
	}
}

// And the guard is a check, not a blanket skip: with the directory there, the
// injection is attempted and its failure is reported rather than swallowed.
func TestAddFieldStillReportsAFrontendFailure(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "packages", "shared", "schemas"), 0o755); err != nil {
		t.Fatal(err)
	}

	def, err := ParseInlineFields("Widget", "name:string,colour:text")
	if err != nil {
		t.Fatal(err)
	}
	g := newTestGenerator(root, "app/apps/api", def)

	// The directory exists but widget.ts does not, which is a real problem in a
	// monorepo: the resource's schema file should be there.
	if err := g.injectFrontendField(BuildNames(def), def.Fields[len(def.Fields)-1]); err == nil {
		t.Error("a missing schema file in a monorepo should be reported, not skipped")
	}
}
