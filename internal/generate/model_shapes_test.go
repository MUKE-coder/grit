package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// The readers are run against what the generator actually writes.
//
// This is the test the pair was missing. `grit doctor` reads a model's source
// to decide whether its suggested `retention.Register` line needs
// `WithLines` or `AsTree`, and it used its own copy of both patterns. A test
// on the reader's side passes on a fixture somebody typed; a test on the
// writer's side passes on output nobody reads back. Only one test that runs
// the generator and then the reader fails when they come apart, which is the
// habit this repository arrived at the hard way.
//
// Fixtures from a real generator run, not hand-written: a model written by
// hand here tests a shape no project has.
func TestTheModelReadersMatchWhatTheGeneratorWrites(t *testing.T) {
	t.Run("a record with lines", func(t *testing.T) {
		src := generatedModel(t, "orders/apps/api", "PurchaseOrder", "reference:string",
			func(d *ResourceDefinition) { d.Items = mustFields(t, "PurchaseOrderLine", "description:string") })

		got := ItemsChildModel(src)
		if got != "PurchaseOrderLine" {
			t.Errorf("ItemsChildModel = %q, want PurchaseOrderLine.\n"+
				"The has-many --items writes has changed shape, so grit doctor would suggest a\n"+
				"retention line with no WithLines and the sweep would fail on a foreign key.\n"+
				"The model's slice fields are:\n%s", got, sliceFields(src))
		}
	})

	t.Run("a tree", func(t *testing.T) {
		src := generatedModel(t, "catalogue/apps/api", "Category", "name:string",
			func(d *ResourceDefinition) { d.Tree = true })

		if !IsTreeModel(src) {
			t.Errorf("IsTreeModel is false for a generated --tree model, so grit doctor would\n"+
				"suggest a retention line with no AsTree() and a parent would be refused\n"+
				"while its children remain. Looked for %q.", treeBanner)
		}
	})

	t.Run("an ordinary resource is neither", func(t *testing.T) {
		src := generatedModel(t, "plain/apps/api", "Note", "body:text", nil)

		if got := ItemsChildModel(src); got != "" {
			t.Errorf("ItemsChildModel = %q on a resource with no --items", got)
		}
		if IsTreeModel(src) {
			t.Error("IsTreeModel is true on a resource with no --tree")
		}
	})
}

// generatedModel writes one resource's model with the real generator and
// returns its source.
func generatedModel(t *testing.T, module, name, fields string, shape func(*ResourceDefinition)) string {
	t.Helper()

	root := setupMinimalProject(t, module)
	def := mustFields(t, name, fields)
	if shape != nil {
		shape(def)
	}
	g := newTestGenerator(root, module, def)
	names := g.Names()
	if err := g.writeGoModel(names); err != nil {
		t.Fatalf("writing the %s model: %v", name, err)
	}
	return readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", names.Snake+".go"))
}

// sliceFields is the model's has-many lines, so a failure shows what the
// generator writes now rather than only that the pattern missed.
func sliceFields(src string) string {
	var out []string
	for _, line := range strings.Split(src, "\n") {
		if strings.Contains(line, "[]") && strings.Contains(line, "gorm:") {
			out = append(out, "  "+strings.TrimSpace(line))
		}
	}
	if len(out) == 0 {
		return "  (none)"
	}
	return strings.Join(out, "\n")
}
