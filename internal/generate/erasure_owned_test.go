package generate

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
)

// Erasing a user left every row they owned in a generated resource in place,
// while the deletion journal recorded the erasure as complete. A resource
// generated with --owned-by now registers itself for erasure.
func TestOwnedResourceRegistersForErasure(t *testing.T) {
	const module = "clinic/apps/api"
	root := setupMinimalProject(t, module)
	def := mustFields(t, "Note", "body:text,user:belongs_to:User")
	def.OwnedBy = "user"
	g := newTestGenerator(root, module, def)
	if err := g.writeGoModel(g.Names()); err != nil {
		t.Fatalf("model: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "note.go"))

	if !strings.Contains(src, `erasure.Register(&Note{}, "user_id")`) {
		t.Error("an owned resource does not register for erasure, so its rows outlive their owner")
	}
	if !strings.Contains(src, `"clinic/apps/api/internal/erasure"`) {
		t.Error("the model registers with erasure without importing it")
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "note.go", src, parser.AllErrors); err != nil {
		t.Fatalf("the model does not parse: %v", err)
	}
}

// Append-only rows cannot be deleted, so they are not registered: an erasure
// that tried would fail the whole transaction, and the conflict between
// retention and erasure is the project's call, not the generator's.
func TestAppendOnlyOwnedResourceIsNotErased(t *testing.T) {
	const module = "clinic/apps/api"
	root := setupMinimalProject(t, module)
	def := mustFields(t, "Consent", "text:text,user:belongs_to:User")
	def.OwnedBy = "user"
	def.AppendOnly = true
	g := newTestGenerator(root, module, def)
	if err := g.writeGoModel(g.Names()); err != nil {
		t.Fatalf("model: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "consent.go"))
	if strings.Contains(src, "erasure.Register") {
		t.Error("an append-only resource was registered for erasure")
	}
}

// Unowned resources are not the user's, and do not register either.
func TestUnownedResourceIsNotErased(t *testing.T) {
	const module = "clinic/apps/api"
	root := setupMinimalProject(t, module)
	g := newTestGenerator(root, module, mustFields(t, "Clinic", "name:string"))
	if err := g.writeGoModel(g.Names()); err != nil {
		t.Fatalf("model: %v", err)
	}
	if strings.Contains(readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "clinic.go")), "erasure") {
		t.Error("an unowned resource mentions erasure")
	}
}
