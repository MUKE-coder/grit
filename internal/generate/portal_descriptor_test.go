package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// A column added to the model turns up in the portal descriptor.
//
// Every descriptor says so at the top of itself, and for one release it was
// not true: the file was written by generate resource and never touched
// again, so a customer screen built against it went on rendering yesterday's
// fields. A generated file whose own comment is wrong about it is worse than
// no file, because the comment is why somebody trusted it.
func TestAddingAFieldReachesThePortalDescriptor(t *testing.T) {
	const module = "keep/apps/api"
	root := saasProject(t, module, true)
	g := ownedNoteGenerator(t, root, module)
	names := g.Names()

	if err := g.writePortalFiles(names); err != nil {
		t.Fatalf("writePortalFiles: %v", err)
	}
	path := filepath.Join(root, "apps", "web", "lib", "portal", "notes.ts")
	before := readTestFile(t, path)

	added := mustFields(t, "Note", "remind_at:datetime").Fields[0]
	if err := g.refreshPortalDescriptor(names, added); err != nil {
		t.Fatalf("refreshPortalDescriptor: %v", err)
	}

	after := readTestFile(t, path)
	if !strings.Contains(after, `{ name: "remind_at", label: "Remind at", kind: "datetime" },`) {
		t.Errorf("the new column is not in the descriptor:\n%s", after)
	}

	// Everything that was there before is still there, unchanged. The first
	// version of this rewrote the file from the model instead of injecting,
	// and the model is a lossy round trip: a text column came back a string,
	// a required column came back optional, and the scope flipped from
	// organization to user, which is a lie about who the rows belong to.
	for _, line := range strings.Split(before, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if !strings.Contains(after, line) {
			t.Errorf("the refresh changed a line nobody touched:\n  %s", line)
		}
	}

	// Twice is once.
	if err := g.refreshPortalDescriptor(names, added); err != nil {
		t.Fatalf("second refresh: %v", err)
	}
	if n := strings.Count(readTestFile(t, path), `name: "remind_at"`); n != 1 {
		t.Errorf("the column appears %d times after two runs, want 1", n)
	}
}

// A tenant-owned resource keeps saying so. This is the specific regression
// the rewrite caused, and it is the one that matters most: scope is what a
// screen reads to know whether these rows belong to a person or to their
// organisation.
func TestTheRefreshKeepsTheScope(t *testing.T) {
	const module = "keep/apps/api"
	root := saasProject(t, module, true)
	def := mustFields(t, "Doc", "title:string,body:text")
	def.TenantOwned = true
	g := newTestGenerator(root, module, def)
	g.Architecture, g.Frontend = "triple", "next"
	names := g.Names()

	if err := g.writePortalFiles(names); err != nil {
		t.Fatalf("writePortalFiles: %v", err)
	}
	path := filepath.Join(root, "apps", "web", "lib", "portal", "docs.ts")
	if !strings.Contains(readTestFile(t, path), `scope: "organization"`) {
		t.Fatal("the descriptor did not say organization to begin with, so this test checks nothing")
	}

	added := mustFields(t, "Doc", "pinned:bool").Fields[0]
	if err := g.refreshPortalDescriptor(names, added); err != nil {
		t.Fatalf("refreshPortalDescriptor: %v", err)
	}

	after := readTestFile(t, path)
	if strings.Contains(after, `scope: "user"`) {
		t.Error("adding a column changed the resource's scope to user, which is false")
	}
	if !strings.Contains(after, `scope: "organization"`) {
		t.Error("the scope is gone")
	}
	// The textarea is the other thing the rewrite lost: a text column came
	// back from the model as a plain string.
	if !strings.Contains(after, `{ name: "body", label: "Body", kind: "textarea" },`) {
		t.Error("the text column came back as something else")
	}
}

// A resource with no descriptor does not acquire one by having a field added.
func TestARefreshWritesNothingWhereThereIsNoDescriptor(t *testing.T) {
	const module = "shop/apps/api"
	root := saasProject(t, module, false) // not a SaaS, so no portal at all
	g := ownedNoteGenerator(t, root, module)
	names := g.Names()

	added := mustFields(t, "Note", "pinned:bool").Fields[0]
	if err := g.refreshPortalDescriptor(names, added); err != nil {
		t.Fatalf("refreshPortalDescriptor: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "apps", "web", "lib", "portal", "notes.ts")); err == nil {
		t.Error("a descriptor was created for a project that never asked for one")
	}
}

// The row the injection writes and the row the full write writes are built by
// one function, so a field type added to one cannot be missing from the other.
func TestOneRowSpelling(t *testing.T) {
	f := portalFormField{
		JSON:     "colour",
		Label:    "Colour",
		Kind:     "select",
		Required: true,
		Options:  []string{"yellow", "blue"},
	}
	row := portalFieldRow(f)
	for _, want := range []string{`name: "colour"`, `label: "Colour"`, `kind: "select"`, "required: true", `options: ["yellow", "blue"]`} {
		if !strings.Contains(row, want) {
			t.Errorf("the row is missing %s: %s", want, row)
		}
	}
	if !strings.HasSuffix(strings.TrimSpace(row), "},") {
		t.Errorf("the row does not close: %s", row)
	}
	_ = scaffold.WebAccountLayout // the fixture's layout comes from the template
}
