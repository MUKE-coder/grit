package scaffold

import (
	"strings"
	"testing"
)

// A belongs_to column is generated with a dotted key ("board.name"), because
// the API returns the related object rather than a flattened field. Every
// surface that renders a column has to resolve that the same way.
//
// The list table did and the record page did not, so on a real CRM a deal's
// page showed a dash for Pipeline, Stage, Contact, Company and Owner while
// the list beside it showed all five (grit#118). Nothing failed: a missing
// nested property is undefined, and undefined renders as an em dash.
func TestEveryColumnSurfaceResolvesDottedKeys(t *testing.T) {
	renderers := adminCellRenderers()

	if !strings.Contains(renderers, "export function getNestedValue(") {
		t.Fatal("cell-renderers no longer exports getNestedValue; the surfaces that " +
			"import it will each need their own copy, which is how they drifted before")
	}

	table := adminDataTable()
	detail := adminResourceDetailPage()

	for name, src := range map[string]string{
		"data-table.tsx":           table,
		"resource-detail-page.tsx": detail,
	} {
		if strings.Count(src, "function getNestedValue(") != 0 {
			t.Errorf("%s declares its own getNestedValue; it should import the shared one", name)
		}
		if !strings.Contains(src, "getNestedValue") {
			t.Errorf("%s does not use getNestedValue, so a dotted column key reads as a flat property", name)
		}
	}

	// The specific line the bug was in.
	if strings.Contains(detail, "renderCell(col, record[col.key]") {
		t.Error("the record page still reads record[col.key], which is always undefined " +
			"for a belongs_to column and renders as a dash")
	}
	if !strings.Contains(detail, "renderCell(col, getNestedValue(record, col.key), record)") {
		t.Error("the record page does not resolve its column key through getNestedValue")
	}
}
