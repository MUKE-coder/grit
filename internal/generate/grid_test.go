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

// Bulk edit must not validate the whole struct.
//
// It sends only the cells that changed, like Patch, and a patch of one cell
// cannot be made to satisfy every required column on the model: validating
// there would refuse an edit of a stock level because some other column the
// operator never touched is empty.
func TestBulkEditDoesNotValidateTheWholeRow(t *testing.T) {
	src := renderService(t, "shop/apps/api", "Gadget", "name:string,category:belongs_to:Category", nil)

	bulkEdit := methodBody(t, src, "func (s *GadgetService) BulkEdit(")
	if strings.Contains(bulkEdit, "ValidateStruct") {
		t.Error("BulkEdit validates a whole row, so a one-cell edit needs every required column")
	}
	if got := strings.Count(src, "respond.ValidateStruct("); got != 1 {
		t.Errorf("ValidateStruct appears %d times; only BulkCreate should call it", got)
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
