package generate

import (
	"os"
	"strings"
	"testing"
)

// Enterprise review 2.4: Patch, BulkEdit and Bulk wrote maps nothing had
// checked.
//
// Create and Update bind into the model, so gin runs the binding tags, and
// they load their row through s.load, which consults a policy. The three
// map-shaped paths did neither: a column whitelist answers "may this column
// be written at all" and nothing else, so binding:"required" never fired and
// {"count":"x"} reached the database as a type error and came back a 500.
//
// Read from the template source. The emitted code is compiled by CI against a
// real scaffold; what these pin is that the checks are still called at all,
// which is the part a refactor silently drops.

// Every map-shaped write runs the model's own rules.
func TestMapWritesAreValidated(t *testing.T) {
	src := readTemplateSource(t)

	// Patch, before the money expansion: that rewrites a field into two
	// columns whose names are not json tags, so a decode after it drops them.
	patch := between(src, "func (s *{{Pascal}}Service) Patch(", "// Delete soft-deletes one")
	if !strings.Contains(patch, "s.checked(item, updates)") {
		t.Error("Patch writes a map nothing has validated: binding:\"required\" " +
			"never fires and a wrong type is a 500 rather than a 422 (review 2.4)")
	}
	if strings.Index(patch, "s.checked(") > strings.Index(patch, "expand{{Pascal}}Embedded(updates)") {
		t.Error("Patch validates after expanding the money columns, whose names " +
			"are not json tags, so those values are silently not checked")
	}

	edit := between(src, "func (s *{{Pascal}}Service) BulkEdit(", "// BulkCreate inserts")
	if !strings.Contains(edit, "s.checked(row, w.updates)") {
		t.Error("BulkEdit writes rows nothing has validated (review 2.4)")
	}

	bulk := between(src, "func (s *{{Pascal}}Service) Bulk(", "// BulkDelete")
	if !strings.Contains(bulk, "s.checked(&items[i], result.Updates)") {
		t.Error("the bulk patch action writes a map nothing has validated. One " +
			"patch over many rows can be invalid for some of them: a " +
			"conditional field's requirement depends on the row (review 2.4)")
	}
}

// And every bulk path consults the policy, which the single-row paths get
// from s.load. Without it, a rule in internal/policies is evadable by
// selecting the rows and using a bulk action.
func TestBulkWritesConsultThePolicy(t *testing.T) {
	src := readTemplateSource(t)

	edit := between(src, "func (s *{{Pascal}}Service) BulkEdit(", "// BulkCreate inserts")
	if !strings.Contains(edit, `s.permitted(ctx, allowed, "update")`) {
		t.Error("BulkEdit does not run the policy per row, so any rule in " +
			"internal/policies is evadable through the grid (review 2.4)")
	}

	bulk := between(src, "func (s *{{Pascal}}Service) Bulk(", "// BulkDelete")
	if !strings.Contains(bulk, "s.permitted(ctx, items, ability)") {
		t.Error("the bulk action does not run the policy per row (review 2.4)")
	}
	// The ability has to follow the action, or a policy that allows archiving
	// and refuses deleting is asked the wrong question.
	if !strings.Contains(bulk, `ability = "delete"`) {
		t.Error("every bulk action asks the policy about \"update\", so a rule " +
			"that refuses a delete and allows an edit is not honoured")
	}
	// A denial is 403, the same answer the single-row path gives. A client
	// branching on status should not see one code for one row and another
	// for twelve.
	if !strings.Contains(bulk, `authz.Denied("{{plural}}."+ability`) {
		t.Error("a bulk denial is not reported as a denial, so it answers 422 " +
			"where the single-row path answers 403")
	}
}

// readTemplateSource reads crud_templates.go, which is where the service
// template lives as a Go string. Reading the file rather than the constant
// because the template is assembled from several of them.
func readTemplateSource(t *testing.T) string {
	t.Helper()
	b, err := os.ReadFile("crud_templates.go")
	if err != nil {
		t.Fatalf("reading the service template: %v", err)
	}
	return string(b)
}

// between returns the slice of src from the first marker to the second, so a
// test about Patch cannot accidentally pass on something in Update.
func between(src, from, to string) string {
	a := strings.Index(src, from)
	if a < 0 {
		return ""
	}
	b := strings.Index(src[a:], to)
	if b < 0 {
		return src[a:]
	}
	return src[a : a+b]
}
