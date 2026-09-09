package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// An inline line item is typed the same way the child's own form types it.
//
// The line-items builder carried a second switch covering five field types and
// falling through to "text" for the rest, so money, toggles, selects,
// textareas and rich text were all plain text inputs. Found on an invoicing
// app: unit_rate was `type: "money"` in the standalone InvoiceItem resource and
// `type: "text"` inline, in the same generation run. A currency amount in a
// plain text box, on the line the whole invoice multiplies.
func TestLineItemFieldsAreTypedLikeAnyOtherForm(t *testing.T) {
	const module = "billing/apps/api"
	root := setupMinimalProject(t, module)

	child, err := ParseInlineFields("InvoiceItem",
		"description:string,qty:int,unit_rate:money,taxable:bool,kind:select:service=Service|goods=Goods")
	if err != nil {
		t.Fatalf("child fields: %v", err)
	}
	parent, err := ParseInlineFields("Invoice", "number:string,discount:money")
	if err != nil {
		t.Fatalf("parent fields: %v", err)
	}
	parent.Items = child

	g := newTestGenerator(root, module, parent)
	if err := g.writeResourceDefinition(g.Names()); err != nil {
		t.Fatalf("resource definition: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "admin", "resources", "invoices", "invoices.ts"))

	items := src[strings.Index(src, "itemFields: ["):]
	if end := strings.Index(items, "] }"); end > 0 {
		items = items[:end]
	}

	for _, want := range []struct{ key, typ string }{
		{"description", "text"},
		{"qty", "number"},
		{"unit_rate", "money"},  // was "text"
		{"taxable", "toggle"},   // was "text"
		{"kind", "select"},      // was "text"
	} {
		frag := `{ key: "` + want.key + `", label: `
		i := strings.Index(items, frag)
		if i < 0 {
			t.Errorf("no line-item field for %q", want.key)
			continue
		}
		line := items[i:]
		if j := strings.Index(line, "\n"); j > 0 {
			line = line[:j]
		}
		if !strings.Contains(line, `type: "`+want.typ+`"`) {
			t.Errorf("line item %q renders as the wrong control:\n  %s\n  want type %q",
				want.key, strings.TrimSpace(line), want.typ)
		}
	}

	// A select with no options renders an empty dropdown, which reads as a
	// loading bug rather than a missing argument.
	if !strings.Contains(items, `{ value: "service", label: "Service" }`) {
		t.Errorf("the line-item select has no options:\n%s", items)
	}
}
