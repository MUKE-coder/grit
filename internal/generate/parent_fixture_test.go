package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// grit#137: the generated handler test inserted a child row with an empty
// parent id.
//
// The factory leaves a belongs_to column empty on purpose and its doc comment
// says to pass the parent in. The test did not, so the row went in with
// author_id = "", which is an empty string rather than NULL because the
// column is a string. Postgres refused it against the foreign key; SQLite
// accepted it, because the test database did not enforce foreign keys.
//
// Green on the database nobody ships, red on the one everybody does.
func TestGeneratedTestCreatesItsBelongsToParents(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Book", "title:string,author:belongs_to:Author")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "book_test.go"))

	assertContains(t, "book_test.go", src,
		"author := factory.CreateAuthor(t, db)",
		"m.AuthorID = author.ID",
	)

	// The bare call is the bug: it is what inserted the empty foreign key.
	if strings.Contains(src, "factory.CreateBook(t, db)") {
		t.Error("the test still creates the record with no parent, so it inserts " +
			"an empty foreign key and fails on Postgres (grit#137)")
	}

	// The parent has to exist before the child is made.
	parent := strings.Index(src, "factory.CreateAuthor(t, db)")
	child := strings.Index(src, "factory.CreateBook(t, db,")
	if parent < 0 || child < 0 || parent > child {
		t.Error("the parent is not created before the record that points at it")
	}
}

// A resource whose parent is its own kind can be made with none: creating one
// first would need the same override and recurse.
func TestGeneratedTestSkipsASelfReference(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Folder", "name:string,parent:belongs_to:Folder")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "folder_test.go"))

	if strings.Contains(src, "factory.CreateFolder(t, db)\n\n\trecord := factory.CreateFolder") {
		t.Error("a self-referencing resource creates itself as its own parent")
	}
	if !strings.Contains(src, "record := factory.CreateFolder(t, db)") {
		t.Error("a self-referencing resource should be created with no parent; " +
			"its column is nullable")
	}
}

// Two parents, two fixtures. A resource with several belongs_to fields had
// every one of them empty.
func TestGeneratedTestCreatesEveryParent(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Task",
		"title:string,project:belongs_to:Project,assignee:belongs_to:User")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "task_test.go"))

	assertContains(t, "task_test.go", src,
		"project := factory.CreateProject(t, db)",
		"assignee := factory.CreateUser(t, db)",
		"m.ProjectID = project.ID",
		"m.AssigneeID = assignee.ID",
	)
}
