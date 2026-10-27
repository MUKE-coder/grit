package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// A refinement is the one validator in the schema that reads a sibling, so it
// is the one that breaks when a field is taken out.
//
// server_set took `status` out of the create schema, and the refinement
// generated for `lost_reason` went on saying `v.status === "lost"`. That is a
// TypeScript error, which stops @repo/shared compiling, which fails the build
// of every app that imports it (grit#115). A feature shipped one release
// earlier broke the build of any project that combined it with a condition.
func TestARefinementIsDroppedWhenItsFieldIsNotInTheSchema(t *testing.T) {
	rules := []conditionalRequired{
		{Field: "lost_reason", On: "status", Value: "lost"},
		{Field: "cancel_note", On: "state", Value: "cancelled"},
	}
	keys := map[string]bool{"lost_reason": true, "cancel_note": true, "state": true}

	kept, dropped := refinementsFor(rules, keys)

	if len(kept) != 1 || kept[0].On != "state" {
		t.Errorf("kept %v, want only the one whose field is present", kept)
	}
	if len(dropped) != 1 || dropped[0].On != "status" {
		t.Errorf("dropped %v, want the one reading the absent status", dropped)
	}
}

// And the ordinary case is untouched: a condition on a field the schema has is
// exactly why refinements exist.
func TestARefinementSurvivesWhenItsFieldIsPresent(t *testing.T) {
	rules := []conditionalRequired{{Field: "lost_reason", On: "status", Value: "lost"}}
	kept, dropped := refinementsFor(rules, map[string]bool{"status": true, "lost_reason": true})

	if len(kept) != 1 {
		t.Fatalf("the refinement was dropped although status is in the schema: %v", dropped)
	}
	if !strings.Contains(zodRefinements(kept), `v.status === "lost"`) {
		t.Error("the rendered refinement does not read the field")
	}
	if len(dropped) != 0 {
		t.Errorf("dropped %v", dropped)
	}
}

// End to end through the generator, because the unit above would pass even if
// nothing called it.
func TestServerSetFieldDoesNotLeaveADanglingRefinement(t *testing.T) {
	root := t.TempDir()
	def := &ResourceDefinition{
		Name: "Deal",
		Fields: []Field{
			{Name: "title", Type: "string", Required: true},
			{Name: "status", Type: "select", Required: true, Default: "open", ServerSet: true,
				Options: []FieldOption{{Value: "open", Label: "Open"}, {Value: "lost", Label: "Lost"}}},
			{Name: "lost_reason", Type: "text", Required: true,
				ShowWhen: &FieldCondition{Field: "status", Value: "lost"}},
		},
	}
	g := newTestGenerator(root, "demo/apps/api", def)
	names := BuildNames(def)
	if err := g.writeZodSchema(names); err != nil {
		t.Fatal(err)
	}
	schema := readFile(t, filepath.Join(g.SharedRoot(), "schemas", "deal.ts"))

	if strings.Contains(schema, "v.status") {
		t.Error("the schema reads v.status, which server_set removed: this does not " +
			"compile, and it fails the build of every app importing @repo/shared")
	}
	if strings.Contains(schema, "status:") {
		t.Error("a server_set field is in the create schema")
	}
	// The field it guards is still there, and still optional.
	if !strings.Contains(schema, "lost_reason: z.string().optional()") {
		t.Error("lost_reason should remain, optional, enforced by the server")
	}
}

// With the same definition minus server_set, the refinement is emitted. This
// is the pair: a fix that drops refinements everywhere would pass the test
// above and be useless.
func TestTheRefinementIsStillEmittedNormally(t *testing.T) {
	root := t.TempDir()
	def := &ResourceDefinition{
		Name: "Deal",
		Fields: []Field{
			{Name: "title", Type: "string", Required: true},
			{Name: "status", Type: "select", Required: true, Default: "open",
				Options: []FieldOption{{Value: "open", Label: "Open"}, {Value: "lost", Label: "Lost"}}},
			{Name: "lost_reason", Type: "text", Required: true,
				ShowWhen: &FieldCondition{Field: "status", Value: "lost"}},
		},
	}
	g := newTestGenerator(root, "demo/apps/api", def)
	names := BuildNames(def)
	if err := g.writeZodSchema(names); err != nil {
		t.Fatal(err)
	}
	schema := readFile(t, filepath.Join(g.SharedRoot(), "schemas", "deal.ts"))

	if !strings.Contains(schema, `v.status === "lost"`) {
		t.Error("the refinement was dropped although status is in the schema")
	}
	if !strings.Contains(schema, "status:") {
		t.Error("status should be in the create schema when it is not server_set")
	}
}
