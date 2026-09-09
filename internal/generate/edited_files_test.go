package generate

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Re-running a generate does not silently rewrite files you have edited.
//
// Adding a field by re-running the same command with one more entry in
// --fields is the obvious thing to do, and it rewrote the model, service and
// handler whole. Anything hand-written in them was gone, with a green tick
// printed for each file as it went. The docs call the service layer "your
// custom business logic goes here" and promise regeneration never touches it.
func TestRegenerateRefusesToOverwriteEditedFiles(t *testing.T) {
	const module = "notes/apps/api"
	root := setupMinimalProject(t, module)

	svc := filepath.Join(root, "apps", "api", "internal", "services", "note.go")
	recordAs(t, root, "resource:Note", svc, "package services\n\n// as generated\n")

	// The developer adds a method, the way the documentation invites them to.
	writeTestFile(t, svc, "package services\n\n// as generated\n\nfunc (s *NoteService) PinnedCount() {}\n")

	g := newTestGenerator(root, module, mustFields(t, "Note", "body:text"))
	err := g.checkEditedFiles(g.Names())
	if err == nil {
		t.Fatal("regenerate would have overwritten an edited service and said nothing")
	}
	msg := err.Error()
	if !strings.Contains(msg, "apps/api/internal/services/note.go") {
		t.Errorf("the message does not name the file at risk:\n%s", msg)
	}
	if !strings.Contains(msg, "--force") {
		t.Errorf("the message offers no way through:\n%s", msg)
	}
}

// An untouched resource regenerates exactly as before. This is the common case
// and the guard must stay out of its way.
func TestRegenerateAllowedWhenNothingWasEdited(t *testing.T) {
	const module = "notes/apps/api"
	root := setupMinimalProject(t, module)

	svc := filepath.Join(root, "apps", "api", "internal", "services", "note.go")
	recordAs(t, root, "resource:Note", svc, "package services\n\n// as generated\n")

	g := newTestGenerator(root, module, mustFields(t, "Note", "body:text"))
	if err := g.checkEditedFiles(g.Names()); err != nil {
		t.Fatalf("refused a regenerate over untouched files: %v", err)
	}
}

// The .custom.tsx overlay exists to be edited and is never rewritten, so an
// edit there is not a reason to refuse. Flagging it would block the command
// for doing the one thing the overlay was added for.
func TestEditedCustomOverlayDoesNotBlockRegenerate(t *testing.T) {
	const module = "notes/apps/api"
	root := setupMinimalProject(t, module)

	overlay := filepath.Join(root, "apps", "admin", "resources", "notes", "notes.custom.tsx")
	recordAs(t, root, "resource:Note", overlay, "const custom = {};\n")
	writeTestFile(t, overlay, "const custom = { columns: {} };\n")

	g := newTestGenerator(root, module, mustFields(t, "Note", "body:text"))
	if err := g.checkEditedFiles(g.Names()); err != nil {
		t.Fatalf("an edited overlay blocked a regenerate: %v", err)
	}
}

// Another resource's edited files are not this resource's problem.
func TestAnotherResourcesEditsAreIgnored(t *testing.T) {
	const module = "notes/apps/api"
	root := setupMinimalProject(t, module)

	svc := filepath.Join(root, "apps", "api", "internal", "services", "post.go")
	recordAs(t, root, "resource:Post", svc, "package services\n")
	writeTestFile(t, svc, "package services\n\nfunc (s *PostService) Mine() {}\n")

	g := newTestGenerator(root, module, mustFields(t, "Note", "body:text"))
	if err := g.checkEditedFiles(g.Names()); err != nil {
		t.Fatalf("Post's edits blocked Note: %v", err)
	}
}

// recordAs writes content and records it in the manifest under generator, the
// way a real generation run would.
func recordAs(t *testing.T, root, generator, abs, content string) {
	t.Helper()
	writeTestFile(t, abs, content)

	m, err := manifest.Load(root)
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}
	rel, ok := manifest.Rel(root, abs)
	if !ok {
		t.Fatalf("%s is not inside %s", abs, root)
	}
	m.Record(rel, generator, "test", content)
	if err := m.Save(root, "test"); err != nil {
		t.Fatalf("save manifest: %v", err)
	}
}

func mustFields(t *testing.T, name, fields string) *ResourceDefinition {
	t.Helper()
	def, err := ParseInlineFields(name, fields)
	if err != nil {
		t.Fatalf("fields: %v", err)
	}
	return def
}
