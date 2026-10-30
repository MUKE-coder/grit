package generate

import (
	"strings"
	"testing"
)

// The grid: bulk create and bulk edit, the two endpoints behind the admin's
// spreadsheet-shaped data entry.
//
// Bulk create had a hole worth a test of its own. It decodes each row with
// encoding/json, which knows nothing about binding: tags, so the model's own
// required never fired on that path: the grid inserted a gadget with no
// category into a table whose model requires one, and the edit form then
// refused to save that row, because a form does enforce it. A row you can
// create and cannot edit is the worst of both, and it looks to whoever hits it
// like the edit form is broken.

func TestBulkCreateValidatesEveryRow(t *testing.T) {
	src := renderService(t, "shop/apps/api", "Gadget", "name:string,category:belongs_to:Category", nil)

	if !strings.Contains(src, "respond.ValidateStruct(&item)") {
		t.Error("BulkCreate decodes rows without running the model's binding tags")
	}
	// Against the row, not the request: the grid puts the message back on the
	// line the operator is looking at.
	if !strings.Contains(src, "GadgetRowError{Index: i, Message: err.Error()}") {
		t.Error("a row that fails validation is not reported against its row")
	}
}

// Bulk edit validates against the stored row, never against a bare patch.
//
// The distinction is the whole thing. A grid sends only the cells that
// changed, so validating the patch on its own would refuse an edit of a stock
// level because some other column the operator never touched is empty. That
// was the reason this path did no validation at all until v3.451.0, and the
// reason it could then write anything the column whitelist allowed:
// {"count":"x"} reached the database as a type error and came back a 500
// (enterprise review 2.4).
//
// Validating a copy of the loaded row with the patch applied answers both.
// The untouched column is not empty, it holds what is stored, so a one-cell
// edit passes; and the cell that did change is checked. Measured on a running
// project: a one-cell grid edit with a required column untouched answers 200,
// and the same edit with a wrong type answers 422 naming the column.
func TestBulkEditValidatesAgainstTheStoredRow(t *testing.T) {
	src := renderService(t, "shop/apps/api", "Gadget", "name:string,category:belongs_to:Category", nil)

	bulkEdit := methodBody(t, src, "func (s *GadgetService) BulkEdit(")

	// Against the row the database holds, not against the patch.
	if !strings.Contains(bulkEdit, "s.checked(row, w.updates)") {
		t.Error("BulkEdit writes a map nothing has validated: binding:\"required\" " +
			"never fires and a wrong type is a 500 rather than a 422 (review 2.4)")
	}
	// Never the patch alone. ValidateStruct on a freshly decoded patch is the
	// mistake this test originally existed to prevent.
	if strings.Contains(bulkEdit, "respond.ValidateStruct(") {
		t.Error("BulkEdit validates a struct directly. It must go through " +
			"s.checked, which applies the patch onto the stored row first, or " +
			"a one-cell edit needs every required column")
	}

	// And s.checked is the thing that starts from the row.
	checked := methodBody(t, src, "func (s *GadgetService) checked(")
	if !strings.Contains(checked, "patched := *current") {
		t.Error("checked does not start from the current row, so every " +
			"map-shaped write now needs every required column")
	}
	if !strings.Contains(checked, "respond.ValidateStruct(&patched)") {
		t.Error("checked does not run the model's rules over the patched row")
	}
}

// methodBody is src from the start of a method to the closing brace in column
// one, which for gofmt-ed Go is the end of that function.
func methodBody(t *testing.T, src, signature string) string {
	t.Helper()
	start := strings.Index(src, signature)
	if start < 0 {
		t.Fatalf("no %s in the generated service", signature)
	}
	rest := src[start:]
	end := strings.Index(rest, "\n}\n")
	if end < 0 {
		t.Fatalf("%s has no end", signature)
	}
	return rest[:end]
}
