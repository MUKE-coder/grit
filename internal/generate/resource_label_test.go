package generate

import (
	"strings"
	"testing"
)

// A multi-word resource was labelled with its Go type name, so the sidebar,
// the page title and every form heading read "InventoryItems".
func TestResourceLabelsAreSpaced(t *testing.T) {
	const module = "erp/apps/api"
	for _, tc := range []struct {
		name, singular, plural, slug string
	}{
		{"InventoryItem", "Inventory Item", "Inventory Items", "inventory-items"},
		{"Budget", "Budget", "Budgets", "budgets"},
	} {
		def, err := ParseInlineFields(tc.name, "sku:string")
		if err != nil {
			t.Fatalf("ParseInlineFields: %v", err)
		}
		g := newTestGenerator(setupMinimalProject(t, module), module, def)
		src := g.resourceDefinitionFileContent(g.Names())
		want := `label: { singular: "` + tc.singular + `", plural: "` + tc.plural + `" }`
		if !strings.Contains(src, want) {
			t.Errorf("%s: want %s in the definition", tc.name, want)
		}
		// The slug is an identifier, not a label, and routes and catalogue
		// keys depend on it: it must not change.
		if !strings.Contains(src, `slug: "`+tc.slug+`"`) {
			t.Errorf("%s: the slug changed", tc.name)
		}
	}
}
