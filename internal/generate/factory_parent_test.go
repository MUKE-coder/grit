package generate

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// A generated factory has to produce a row that saves, and a generated test has
// to compile. Both were false, and together they are why the job that runs a
// generated project's tests was red from v3.464.0 to v3.487.0.

// TestFactoryFillsAFrameworkParent: a belongs_to to one of the framework's own
// models gets its foreign key filled, because the column is a string whose zero
// value is "" rather than NULL, and no parent row has that id.
func TestFactoryFillsAFrameworkParent(t *testing.T) {
	for _, parent := range scaffold.FrameworkFactoryModels {
		t.Run(parent, func(t *testing.T) {
			const module = "myapp/apps/api"
			root := setupMinimalProject(t, module)

			def, err := ParseInlineFields("Note",
				fmt.Sprintf("title:string,owner:belongs_to:%s", parent))
			if err != nil {
				t.Fatalf("ParseInlineFields: %v", err)
			}
			g := newTestGenerator(root, module, def)
			if err := g.Run(); err != nil {
				t.Fatalf("Generator.Run(): %v", err)
			}

			src := readTestFile(t, filepath.Join(root, "apps", "api",
				"internal", "factory", "note.go"))

			want := fmt.Sprintf("\tif record.OwnerID == \"\" {\n\t\trecord.OwnerID = Create%s(tb, db).ID\n\t}", parent)
			if !strings.Contains(src, want) {
				t.Errorf("CreateNote does not fill OwnerID from Create%s, so the insert fails "+
					"with a foreign key violation whenever nobody passes one:\n%s", parent, src)
			}
		})
	}
}

// TestFactoryLeavesASelfReferenceAlone: a resource whose parent is its own kind
// must not create one, or every fixture recurses.
func TestFactoryLeavesASelfReferenceAlone(t *testing.T) {
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

	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "factory", "folder.go"))
	if strings.Contains(src, "CreateFolder(tb, db).ID") {
		t.Error("CreateFolder creates its own parent, so every fixture recurses until the stack ends")
	}
}

// TestFactoryDoesNotCallAFactoryThatIsNotThere: a parent with no factory is
// left to the caller. Emitting the call would trade a failing test for a
// package that will not compile, which is strictly worse.
func TestFactoryDoesNotCallAFactoryThatIsNotThere(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	// Session is a framework model with no factory.
	def, err := ParseInlineFields("Note", "title:string,origin:belongs_to:Session")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "factory", "note.go"))
	if strings.Contains(src, "CreateSession(tb, db)") {
		t.Error("the factory calls CreateSession, which nothing writes, so the package does not compile")
	}
}

// fixtureDecl finds every `name := factory.CreateX(t, db)` in a generated test.
var fixtureDecl = regexp.MustCompile(`(?m)^\t(\w+) := factory\.Create\w+\(t, db`)

// TestGeneratedTestDeclaresEachFixtureOnce.
//
// The bulk-edit stale-row test interpolated the parent fixtures once per row,
// so a resource with a belongs_to produced two `file := factory.CreateUpload`
// lines in one function and the package failed with "no new variables on left
// side of :=". v3.483.0 shipped that, and the unit test for it asserted the
// emitted text contained the fixture line, which it did, twice.
func TestGeneratedTestDeclaresEachFixtureOnce(t *testing.T) {
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

	// Per function, because the same variable in two functions is fine.
	for _, fn := range strings.Split(src, "\nfunc ") {
		seen := map[string]int{}
		for _, m := range fixtureDecl.FindAllStringSubmatch(fn, -1) {
			seen[m[1]]++
		}
		for name, count := range seen {
			if count > 1 {
				head := fn
				if i := strings.Index(head, "("); i > 0 {
					head = head[:i]
				}
				t.Errorf("%s declares %q %d times with :=, so the test package does not compile",
					head, name, count)
			}
		}
	}
}
