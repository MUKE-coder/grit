package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// Two follow-ups to the parent auto-fill added in v3.488.0, both found by
// generating resources it had not met.

// grit#193: an :optional belongs_to gives the model a *string foreign key, and
// the auto-fill emitted the required form for both cardinalities:
//
//	internal/factory/ai_run.go:56:25: invalid operation: record.ProjectID == ""
//	  (mismatched types *string and untyped string)
//
// A pointer-aware version would compile and still be wrong: it would make an
// optional relation mandatory in every fixture, so one row would write three,
// and the case the column exists to express, a row with no parent, would never
// be the default. So nothing is emitted, in the factory and in the test.
func TestOptionalParentIsLeftNil(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	// The parents first. Their factories have to exist, or the auto-fill is
	// skipped because it has nothing to call and this test proves nothing: it
	// passed with the bug reintroduced until they were generated.
	generateInto(t, root, module, "Project", "name:string")
	generateInto(t, root, module, "Document", "title:string")

	def, err := ParseInlineFields("AIRun",
		"project:belongs_to[Project]:optional,document:belongs_to[Document]:optional,model:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}
	if !g.parentFactoryExists("Project") {
		t.Fatal("Project has no factory, so the auto-fill would be skipped for the wrong " +
			"reason and this test would pass with the bug present")
	}

	factory := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "factory", "ai_run.go"))
	for _, forbidden := range []string{
		`record.ProjectID == ""`,
		`record.DocumentID == ""`,
		"record.ProjectID = CreateProject(tb, db).ID",
	} {
		if strings.Contains(factory, forbidden) {
			t.Errorf("the factory emits %q for an optional foreign key, which is a *string "+
				"and does not compile against it:\n%s", forbidden, factory)
		}
	}
	// And it says so, because a fixture that quietly has no parent is worth a
	// line of explanation.
	if !strings.Contains(factory, "stay nil") {
		t.Errorf("the factory does not say the optional parents are left nil:\n%s", factory)
	}

	// The handlers test had the same problem from the same helper: it assigned
	// parent.ID, a string, to the same *string column.
	test := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "ai_run_test.go"))
	for _, forbidden := range []string{
		"m.ProjectID = project.ID",
		"m.DocumentID = document.ID",
	} {
		if strings.Contains(test, forbidden) {
			t.Errorf("the generated test emits %q, assigning a string to a *string column:\n%s",
				forbidden, test)
		}
	}
}

// A required parent still gets its fill: the fix for the optional case must not
// take the thing it was built for.
func TestRequiredParentStillFilled(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	// User is a framework model, so its factory ships with the project and
	// needs no generating.
	def, err := ParseInlineFields("Note", "title:string,author:belongs_to:User")
	if err != nil {
		t.Fatal(err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	factory := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "factory", "note.go"))
	if !strings.Contains(factory, `record.AuthorID == ""`) {
		t.Errorf("a required parent is no longer filled in, so CreateNote cannot save:\n%s", factory)
	}
}

// grit#195: the bulk-edit stale-row test gave both rows the same parent, which
// v3.488.0 introduced to fix a redeclared variable. A one_to_one is a
// belongs_to with a UNIQUE foreign key, so the second insert is refused:
//
//	UNIQUE constraint failed: methodology_profiles.project_id
//
// The second row omits the assignment and the factory gives it a parent of its
// own, because the column is empty. Sharing is kept for a plain belongs_to,
// where two of every parent would be rows nothing asked for.
func TestGridTestDoesNotShareAUniqueParent(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)
	generateInto(t, root, module, "Project", "name:string")

	def, err := ParseInlineFields("MethodologyProfile",
		"project:one_to_one[Project],notes:text")
	if err != nil {
		t.Fatal(err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	src := readTestFile(t, filepath.Join(root, "apps", "api",
		"internal", "handlers", "methodology_profile_test.go"))

	grid := src[strings.Index(src, "BulkEditRefusesAStaleRow"):]
	if i := strings.Index(grid, "client = client.As"); i != -1 {
		grid = grid[:i]
	}
	if strings.Count(grid, "m.ProjectID = project.ID") != 1 {
		t.Errorf("the two rows of the bulk-edit test do not share exactly one assignment of a "+
			"unique foreign key; the second row must make its own parent:\n%s", grid)
	}
	if !strings.Contains(grid, "second := factory.CreateMethodologyProfile(t, db)") {
		t.Errorf("the second row does not let the factory create its own parent:\n%s", grid)
	}
}

// And a plain belongs_to still shares, so the fix is narrow.
func TestGridTestSharesAPlainParent(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)
	generateInto(t, root, module, "Project", "name:string")

	def, err := ParseInlineFields("Task", "title:string,project:belongs_to[Project]")
	if err != nil {
		t.Fatal(err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "task_test.go"))
	grid := src[strings.Index(src, "BulkEditRefusesAStaleRow"):]
	if i := strings.Index(grid, "client = client.As"); i != -1 {
		grid = grid[:i]
	}
	if strings.Count(grid, "m.ProjectID = project.ID") != 2 {
		t.Errorf("both rows of the bulk-edit test should share a plain belongs_to parent, "+
			"because two of them would be rows nothing asked for:\n%s", grid)
	}
}

// generateInto generates a resource into an existing test project, so a later
// resource can point at it and find its factory.
func generateInto(t *testing.T, root, module, name, fields string) {
	t.Helper()
	def, err := ParseInlineFields(name, fields)
	if err != nil {
		t.Fatalf("ParseInlineFields(%s): %v", name, err)
	}
	if err := newTestGenerator(root, module, def).Run(); err != nil {
		t.Fatalf("generating %s: %v", name, err)
	}
}
