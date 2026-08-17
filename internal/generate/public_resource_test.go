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

	if !strings.Contains(src, "archived_at IS NULL") {
		t.Error("a public list must exclude archived rows")
	}
	// archived_at is scoped in Go and must never become a query parameter, or
	// ?archived_at= hands back rows somebody took down on purpose.
	if strings.Contains(src, `"archived_at": true`) {
		t.Error("archived_at must not be filterable from the query string")
	}
	// It is an allowlist struct, never the model, so a column added next month
	// is private until somebody says otherwise.
	if strings.Contains(src, "res.Data\n") && !strings.Contains(src, "toPublicProduct(m)") {
		t.Error("the list must map through the allowlist, not return models")
	}
}

// The filter lists are built from the published fields, which is the property
// that keeps them honest: a storefront can filter on anything it can see, and
// on nothing it cannot. A held-back column reachable through ?cost_price= would
// leak by comparison what the allowlist refused to leak directly.
func TestPublicFiltersOnlyPublishedColumns(t *testing.T) {
	g := &Generator{
		Module: "shopfront/apps/api",
		Definition: &ResourceDefinition{
			Name:   "Product",
			Public: true,
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "price", Type: "float"},
				{Name: "cost_price", Type: "float"},
				{Name: "stock", Type: "int"},
				{Name: "active", Type: "bool"},
				{Name: "category", Type: "belongs_to", RelatedModel: "Category"},
			},
		},
	}
	included, _ := PublicFields(g.Definition.Fields)
	src := g.publicHandlerSource(g.Names(), included)

	filters := src[strings.Index(src, "Filterable:"):]
	filters = filters[:strings.Index(filters, "},\n\t\t},")+1]

	for _, want := range []string{
		`"name": true`,
		`"price": true`,
		// The foreign key, which the response holds back but a category page
		// cannot do without.
		`"category_id": true`,
	} {
		if !strings.Contains(filters, want) {
			t.Errorf("expected %s in the filter lists, got:\n%s", want, filters)
		}
	}

	// Held back from the response, so held back from the query string.
	for _, forbidden := range []string{"cost_price", "stock", "active"} {
		if strings.Contains(filters, `"`+forbidden+`": true`) {
			t.Errorf("%q is not published, so it must not be filterable:\n%s", forbidden, filters)
		}
	}

	// A window on the numbers, because equality on a price is never the
	// question a storefront is asking.
	if !strings.Contains(src, "RangeFilterable: map[string]bool{\"price\": true}") {
		t.Errorf("price should accept a min/max window:\n%s", filters)
	}
}

// The similar-items strip. Keyed on a relation the generator picks, so it
// cannot be turned into a filter on something unpublished.
func TestPublicRelatedEndpoint(t *testing.T) {
	withParent := &Generator{
		Module: "shopfront/apps/api",
		Definition: &ResourceDefinition{
			Name:   "Product",
			Public: true,
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "slug", Type: "slug"},
				{Name: "category", Type: "belongs_to", RelatedModel: "Category"},
			},
		},
	}
	included, _ := PublicFields(withParent.Definition.Fields)
	src := withParent.publicHandlerSource(withParent.Names(), included)

	for _, want := range []string{
		"func (h *ProductHandler) RelatedPublic(",
		// The model spells it CategoryID, not CategoryId.
		"item.CategoryID",
		`Where("category_id = ?"`,
		// Excludes itself, or the strip shows the page you are already on.
		`Where("id <> ? AND archived_at IS NULL"`,
		// Uncapped limits on a public route are a free way to make the
		// database work.
		"if limit > 24 {",
		`"strconv"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("related endpoint missing %q", want)
		}
	}

	// No parent, no endpoint: similarity would be arbitrary. And then strconv
	// must not be imported either, or the file does not build.
	noParent := &Generator{
		Module: "shopfront/apps/api",
		Definition: &ResourceDefinition{
			Name:   "Page",
			Public: true,
			Fields: []Field{{Name: "title", Type: "string"}, {Name: "slug", Type: "slug"}},
		},
	}
	noParentIncluded, _ := PublicFields(noParent.Definition.Fields)
	src = noParent.publicHandlerSource(noParent.Names(), noParentIncluded)
	if strings.Contains(src, "RelatedPublic") {
		t.Error("a resource with no parent must not get a related endpoint")
	}
	if strings.Contains(src, `"strconv"`) {
		t.Error("strconv is only needed by the related endpoint, and an unused import will not build")
	}
}

// A file field pulls in an extra import, and it has to land in the right group.
// gofmt sorts inside a group but never moves a line between them, so a local
// import emitted beside "net/http" sorts above it and stays there.
func TestPublicHandlerGroupsItsImports(t *testing.T) {
	g := &Generator{
		Module:     "shopfront/apps/api",
		Definition: &ResourceDefinition{Name: "Product", Public: true},
	}
	src := g.publicHandlerSource(g.Names(), []Field{
		{Name: "name", Type: "string"},
		{Name: "images", Type: "files"},
	})

	files := strings.Index(src, `"shopfront/apps/api/internal/files"`)
	paginate := strings.Index(src, `"shopfront/apps/api/internal/paginate"`)
	http := strings.Index(src, `"net/http"`)
	if files < 0 {
		t.Fatal("a files field must import the files package")
	}
	if files < http || files < paginate {
		t.Errorf("the local import must sit in the local group, after net/http (%d) and paginate (%d), got %d",
			http, paginate, files)
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
