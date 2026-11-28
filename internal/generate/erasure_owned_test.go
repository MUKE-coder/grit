package generate

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
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

// An owned resource with a file column registers a cleanup, because its rows
// name objects in storage and a DELETE does not reach them.
//
// What was lost without it: a file somebody else uploaded onto one of these
// rows. Its uploads row belongs to the uploader, so erasing the owner left it
// alone, and the record that claimed it was hard-deleted, so the orphan sweep
// could not see it either (grit#311).
func TestOwnedResourceWithAFileColumnErasesItsFiles(t *testing.T) {
	const module = "clinic/apps/api"
	root := setupMinimalProject(t, module)
	// internal/files comes from the scaffold template rather than a stub. The
	// generator writes the cleanup only when the project's own copy declares
	// what the cleanup calls, so a hand-written stub would prove this works
	// on a project shape that does not exist.
	writeTestFile(t, filepath.Join(root, "apps", "api", "internal", "files", "lifecycle.go"),
		strings.ReplaceAll(scaffold.FilesLifecycleGo(), "{{MODULE}}", module))

	def := mustFields(t, "Scan", "title:string,report:file,user:belongs_to:User")
	def.OwnedBy = "user"
	g := newTestGenerator(root, module, def)
	if err := g.writeGoModel(g.Names()); err != nil {
		t.Fatalf("model: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "scan.go"))

	for _, want := range []string{
		`erasure.RegisterWithCleanup(&Scan{}, "user_id", eraseScanFiles)`,
		"func eraseScanFiles(tx *gorm.DB, ownerColumn, userID string) error {",
		"files.EraseOwnedFiles(context.Background(), tx, disk, &Scan{}, ownerColumn, userID)",
		`"context"`,
		`"clinic/apps/api/internal/storage"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the model is missing %q, so erasing a user leaves the files on its rows", want)
		}
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "scan.go", src, parser.AllErrors); err != nil {
		t.Fatalf("the model does not parse: %v", err)
	}
}

// A project whose internal/files predates files.EraseOwnedFiles gets the
// registration it always got, rather than a model that calls a function it
// does not have. prepareErasure refreshes that package where it can and says
// so where it cannot, and a generated resource that does not build would be
// worse than the leak.
func TestOwnedResourceWithAFileColumnDegradesOnAnOlderProject(t *testing.T) {
	const module = "clinic/apps/api"
	root := setupMinimalProject(t, module) // no internal/files at all
	def := mustFields(t, "Scan", "title:string,report:file,user:belongs_to:User")
	def.OwnedBy = "user"
	g := newTestGenerator(root, module, def)
	if err := g.writeGoModel(g.Names()); err != nil {
		t.Fatalf("model: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "scan.go"))

	if !strings.Contains(src, `erasure.Register(&Scan{}, "user_id")`) {
		t.Error("the resource does not register for erasure at all, so its rows outlive their owner")
	}
	if strings.Contains(src, "EraseOwnedFiles") {
		t.Error("the model calls files.EraseOwnedFiles, which this project's internal/files does not declare")
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "scan.go", src, parser.AllErrors); err != nil {
		t.Fatalf("the model does not parse: %v", err)
	}
}
