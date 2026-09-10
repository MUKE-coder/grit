package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// An append-only resource is created and read, never changed or deleted.
//
// Found building a double-entry ledger. The generator gave a JournalEntry PUT,
// PATCH, DELETE, bulk and import routes and an admin with Edit and Delete, and
// the hand-rolled guard that followed was beaten by one UPDATE typed into the
// GORM Studio SQL editor. These pin each piece of the flag separately, so a
// later change to the ordinary CRUD templates cannot quietly hand the verbs
// back.

func appendOnlyGenerator(t *testing.T, fields string) (*Generator, string) {
	t.Helper()
	const module = "ledger/apps/api"
	root := setupMinimalProject(t, module)
	def := mustFields(t, "JournalEntry", fields)
	def.AppendOnly = true
	return newTestGenerator(root, module, def), root
}

func TestAppendOnlyRoutesAreReadAndCreate(t *testing.T) {
	g, _ := appendOnlyGenerator(t, "reference:string")
	src, err := g.resourceRoutesSource(g.Names())
	if err != nil {
		t.Fatalf("routes: %v", err)
	}

	for _, want := range []string{
		`m.Protected.GET("/journal_entries", h.List)`,
		`m.Protected.GET("/journal_entries/:id", h.GetByID)`,
		`m.Protected.POST("/journal_entries", h.Create)`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("missing %s", want)
		}
	}
	for _, verb := range []string{"h.Update", "h.Patch", "h.Delete", "h.Bulk", "h.Import"} {
		if strings.Contains(src, verb) {
			t.Errorf("an append-only resource still routes %s", verb)
		}
	}
}

func TestAppendOnlyRoutesKeepTheRoleGuard(t *testing.T) {
	g, _ := appendOnlyGenerator(t, "reference:string")
	g.Roles = []string{"ADMIN"}
	src, err := g.resourceRoutesSource(g.Names())
	if err != nil {
		t.Fatalf("routes: %v", err)
	}
	if !strings.Contains(src, `middleware.RequireRole("ADMIN")`) {
		t.Error("--roles was dropped for an append-only resource")
	}
	if strings.Contains(src, "h.Update") || strings.Contains(src, "h.Delete") {
		t.Error("the role-restricted variant still routes update or delete")
	}
}

func TestAppendOnlyModelRegistersItself(t *testing.T) {
	g, root := appendOnlyGenerator(t, "reference:string")
	if err := g.writeGoModel(g.Names()); err != nil {
		t.Fatalf("model: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "journal_entry.go"))

	if !strings.Contains(src, "appendonly.Register(&JournalEntry{})") {
		t.Error("the model never registers with the guard, so nothing is protected")
	}
	if !strings.Contains(src, `"ledger/apps/api/internal/appendonly"`) {
		t.Error("the model calls appendonly without importing it")
	}
}

// An ordinary resource must not change at all.
func TestOrdinaryResourceIsUnaffected(t *testing.T) {
	const module = "ledger/apps/api"
	root := setupMinimalProject(t, module)
	g := newTestGenerator(root, module, mustFields(t, "Invoice", "number:string"))

	if err := g.writeGoModel(g.Names()); err != nil {
		t.Fatalf("model: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "invoice.go"))
	if strings.Contains(src, "appendonly") {
		t.Error("an ordinary model mentions appendonly")
	}

	routes, err := g.resourceRoutesSource(g.Names())
	if err != nil {
		t.Fatalf("routes: %v", err)
	}
	if !strings.Contains(routes, "h.Update") || !strings.Contains(routes, "h.Delete") {
		t.Error("an ordinary resource lost its update or delete route")
	}
}

func TestAppendOnlyAdminOffersCreateAndView(t *testing.T) {
	g, _ := appendOnlyGenerator(t, "reference:string")
	src := g.adminDefinitionContent(g.Names())

	if !strings.Contains(src, `actions: ["create", "view"]`) {
		t.Error("the admin table still offers edit and delete")
	}
	if !strings.Contains(src, `bulkActions: ["export"]`) {
		t.Error("bulk edit, archive or delete survived")
	}
	// If the ordinary template's bulkActions line ever changes, the Replace in
	// adminDefinitionContent becomes a silent no-op. This is what notices.
	if strings.Contains(src, `"archive"`) {
		t.Error("the editable bulk actions were not replaced: has the template line changed?")
	}
}

func TestAppendOnlyDocsOmitUpdateAndDelete(t *testing.T) {
	g, _ := appendOnlyGenerator(t, "reference:string")
	docs := "\tdocs.Route(\"GET /x\").\n\t\tSummary(\"List\")\n" +
		"\tdocs.Route(\"PUT /x/:id\").\n\t\tSummary(\"Update\")\n" +
		"\tdocs.Route(\"DELETE /x/:id\").\n\t\tSummary(\"Delete\")"
	got := g.appendOnlyDocs(docs)
	if strings.Contains(got, "PUT") || strings.Contains(got, "DELETE") {
		t.Errorf("the API reference documents routes that do not exist:\n%s", got)
	}
	if !strings.Contains(got, "GET /x") {
		t.Error("the read route was dropped from the docs")
	}
}

func TestAppendOnlyRefusesTree(t *testing.T) {
	g, _ := appendOnlyGenerator(t, "name:string")
	g.Definition.Tree = true
	if err := g.prepareAppendOnly(); err == nil || !strings.Contains(err.Error(), "--tree") {
		t.Fatalf("a tree moves nodes, which is an update: want a refusal, got %v", err)
	}
}

// A project too old to have per-resource route files is told what to do rather
// than handed half a feature.
func TestAppendOnlyNeedsTheRouteLayout(t *testing.T) {
	g, root := appendOnlyGenerator(t, "name:string")
	_ = os.Remove(filepath.Join(root, "apps", "api", "internal", "routes", "resources.go"))
	err := g.prepareAppendOnly()
	if err == nil || !strings.Contains(err.Error(), "grit upgrade") {
		t.Fatalf("want a pointer to grit upgrade, got %v", err)
	}
}
