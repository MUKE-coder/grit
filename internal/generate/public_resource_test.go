package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func names(fields []Field) []string {
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		out = append(out, f.Name)
	}
	return out
}

// The default has to be safe, because the audience is the internet and the
// person running the command is usually not thinking about it.
func TestPublicFieldsHoldsBackWhatShouldNotBePublished(t *testing.T) {
	fields := []Field{
		{Name: "name", Type: "string"},
		{Name: "slug", Type: "slug"},
		{Name: "price", Type: "float"},
		{Name: "description", Type: "richtext"},
		{Name: "images", Type: "files"},
		// Money nobody publishes.
		{Name: "cost_price", Type: "float"},
		{Name: "margin", Type: "float"},
		{Name: "supplier_code", Type: "string"},
		// Notes nobody publishes.
		{Name: "internal_note", Type: "text"},
		// A count competitors would enjoy.
		{Name: "stock", Type: "int"},
		// Always the same value on an endpoint that filters on it.
		{Name: "active", Type: "bool"},
		// A relation would publish a whole record nobody vetted.
		{Name: "category", Type: "belongs_to", RelatedModel: "Category"},
	}

	included, excluded := PublicFields(fields)

	wantIn := map[string]bool{
		"name": true, "slug": true, "price": true,
		"description": true, "images": true,
	}
	for _, f := range included {
		if !wantIn[f.Name] {
			t.Errorf("%s should not be published by default", f.Name)
		}
		delete(wantIn, f.Name)
	}
	for missing := range wantIn {
		t.Errorf("%s should be published by default and was not", missing)
	}

	for _, want := range []string{
		"cost_price", "margin", "supplier_code", "internal_note",
		"stock", "active", "category",
	} {
		found := false
		for _, f := range excluded {
			if f.Name == want {
				found = true
			}
		}
		if !found {
			t.Errorf("%s was published by default and must not be (got %v)", want, names(excluded))
		}
	}
}

// A name containing a held-back word is held back whatever its type, because
// the word is the signal and somebody will store a cost as a string.
func TestPublicFieldsIsDrivenByNameNotJustType(t *testing.T) {
	for _, name := range []string{
		"cost_price", "unit_cost", "profit_margin", "internal_ref",
		"admin_note", "api_secret", "access_token", "wholesale_price",
	} {
		included, _ := PublicFields([]Field{{Name: name, Type: "string"}})
		if len(included) != 0 {
			t.Errorf("%q was published and should not have been", name)
		}
	}
}

// featured is genuinely wanted on a homepage, so it is not swept up with the
// visibility flags.
func TestPublicFieldsKeepsFeatured(t *testing.T) {
	included, _ := PublicFields([]Field{{Name: "featured", Type: "bool"}})
	if len(included) != 1 {
		t.Error("featured should be publishable: a homepage needs it")
	}
}

func TestPublicHandlerSource(t *testing.T) {
	g := &Generator{
		Module:     "shopfront/apps/api",
		Definition: &ResourceDefinition{Name: "Product", Public: true},
	}
	included := []Field{
		{Name: "name", Type: "string"},
		{Name: "slug", Type: "slug"},
		{Name: "price", Type: "float"},
	}
	src := g.publicHandlerSource(g.Names(), included)

	for _, want := range []string{
		"type publicProduct struct",
		"func toPublicProduct(",
		"func (h *ProductHandler) ListPublic(",
		"func (h *ProductHandler) GetPublic(",
		// Looked up by slug, so a public URL reads as something a person could
		// type rather than a UUID.
		`Where("slug = ? AND archived_at IS NULL"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated handler missing %q", want)
		}
	}

	// The two properties that make this endpoint safe rather than just public.
	// "Filterable:" with the colon, because the generated file explains in a
	// comment why it has no Filterable, and matching the bare word finds that.
	if strings.Contains(src, "Filterable:") {
		t.Error("a public endpoint must not accept arbitrary column filters")
	}
	if !strings.Contains(src, "archived_at IS NULL") {
		t.Error("a public list must exclude archived rows")
	}
	// It is an allowlist struct, never the model, so a column added next month
	// is private until somebody says otherwise.
	if strings.Contains(src, "res.Data\n") && !strings.Contains(src, "toPublicProduct(m)") {
		t.Error("the list must map through the allowlist, not return models")
	}
}

// A resource whose every field is held back would publish an endpoint that
// returns nothing but ids, which is worse than refusing.
func TestPublicRefusesAResourceWithNothingToPublish(t *testing.T) {
	g := &Generator{
		Module: "shopfront/apps/api",
		Definition: &ResourceDefinition{
			Name:   "Ledger",
			Public: true,
			Fields: []Field{
				{Name: "cost_price", Type: "float"},
				{Name: "internal_note", Type: "text"},
			},
		},
	}
	err := g.writePublicResource(g.Names())
	if err == nil {
		t.Fatal("expected a refusal rather than an endpoint returning only ids")
	}
	if !strings.Contains(err.Error(), "publishable") {
		t.Errorf("the error should explain why: %v", err)
	}
}

func TestPublicIsOptIn(t *testing.T) {
	g := &Generator{
		Module: "shopfront/apps/api",
		Definition: &ResourceDefinition{
			Name:   "Product",
			Fields: []Field{{Name: "name", Type: "string"}},
			// Public not set.
		},
	}
	if err := g.writePublicResource(g.Names()); err != nil {
		t.Fatalf("a resource without --public should be a no-op, got %v", err)
	}
}

// The generated file tells the reader that regenerating will not overwrite it,
// because the allowlist in it is theirs to edit. That promise has to be kept in
// code, or somebody loses the two columns they deliberately published.
func TestPublicHandlerIsNotOverwritten(t *testing.T) {
	root := t.TempDir()
	g := &Generator{
		Root:         root,
		Module:       "shopfront/apps/api",
		Architecture: "single",
		Definition: &ResourceDefinition{
			Name:   "Product",
			Public: true,
			Fields: []Field{{Name: "name", Type: "string"}},
		},
	}
	names := g.Names()

	if err := g.writePublicResource(names); err != nil {
		t.Fatalf("first write: %v", err)
	}
	path := filepath.Join(root, "internal", "handlers", "product_public.go")

	// The developer publishes a column the default held back.
	edited := "package handlers\n\n// I added stock on purpose.\n"
	if err := os.WriteFile(path, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := g.writePublicResource(names); err != nil {
		t.Fatalf("second write: %v", err)
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != edited {
		t.Error("regenerating overwrote the public handler and took the developer's allowlist with it")
	}
}
