package generate

import (
	"strings"
	"testing"
)

func catalogueGenerator() *Generator {
	return &Generator{
		Module: "shopfront/apps/api",
		Root:   "shopfront",
		Definition: &ResourceDefinition{
			Name:   "Product",
			Public: true,
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "slug", Type: "slug"},
				{Name: "price", Type: "int"},
				{Name: "stock", Type: "int"},
				{Name: "image", Type: "file"},
				{Name: "category", Type: "belongs_to", RelatedModel: "Category"},
				{Name: "active", Type: "bool"},
			},
		},
	}
}

// grit#91 item 2: only archived_at hid a row, so an admin's "active: off" was
// still on sale.
func TestPublicVisibilityColumn(t *testing.T) {
	cases := []struct {
		name   string
		fields []Field
		want   string
	}{
		{"active", []Field{{Name: "name", Type: "string"}, {Name: "active", Type: "bool"}}, "active"},
		{"published", []Field{{Name: "published", Type: "bool"}}, "published"},
		{"is_visible", []Field{{Name: "is_visible", Type: "toggle"}}, "is_visible"},
		// featured is a bool about the row, not about who may see it.
		{"featured only", []Field{{Name: "featured", Type: "bool"}}, ""},
		// The name has to be on a boolean: a string column called "status" is
		// not a switch.
		{"string named active", []Field{{Name: "active", Type: "string"}}, ""},
		{"none", []Field{{Name: "name", Type: "string"}}, ""},
		// Declaration order decides, so regenerating gives the same answer.
		{"first wins", []Field{{Name: "published", Type: "bool"}, {Name: "active", Type: "bool"}}, "published"},
	}
	for _, tc := range cases {
		if got := publicVisibilityColumn(tc.fields); got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestPublicScopeGatesOnVisibility(t *testing.T) {
	g := catalogueGenerator()
	svc := g.publicServiceMethods(g.Names())
	if !strings.Contains(svc, `q = q.Where("active = ?", true)`) {
		t.Errorf("the public scope does not gate on active:\n%s", svc)
	}
	// Bound, not spliced: "true" is not the same word in raw SQL on every
	// engine Grit supports.
	if strings.Contains(svc, "active = true") {
		t.Error("the visibility value is a SQL literal rather than a bound parameter")
	}

	// A resource with no such column keeps the archived_at scope and grows no
	// condition it cannot satisfy.
	flat := &Generator{
		Module:     "shopfront/apps/api",
		Definition: &ResourceDefinition{Name: "Page", Public: true, Fields: []Field{{Name: "title", Type: "string"}}},
	}
	svc = flat.publicServiceMethods(flat.Names())
	if !strings.Contains(svc, `q = q.Where("archived_at IS NULL")`) || strings.Contains(svc, "= ?\", true)") {
		t.Errorf("a resource with no visibility column got one:\n%s", svc)
	}
}

// grit#91 item 3: the response omitted the foreign key along with the relation,
// so a storefront could not link a product to its category.
func TestPublicResponsePublishesForeignKeys(t *testing.T) {
	g := catalogueGenerator()
	src := g.publicHandlerSource(g.Names(), firstReturn(PublicFields(g.Definition.Fields)))
	for _, want := range []string{
		"CategoryID string `json:\"category_id\"`",
		"CategoryID: m.CategoryID,",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the public response is missing %q", want)
		}
	}
	// The relation itself stays unpublished: publishing it publishes a whole
	// record nobody vetted.
	if strings.Contains(src, "Category *models.Category") || strings.Contains(src, "Category: m.Category,") {
		t.Error("the relation was published along with its id")
	}
	// A count competitors enjoy, and a visibility switch that is the same on
	// every row returned, both stay out.
	if strings.Contains(src, `json:"stock"`) || strings.Contains(src, `json:"active"`) {
		t.Error("stock or active leaked into the public response")
	}
}

// A tree's parent is the self-referencing key and the tree block already
// publishes it, so it must not be published twice.
func TestPublicResponseDoesNotRepeatTheTreeParent(t *testing.T) {
	g := &Generator{
		Module: "shopfront/apps/api",
		Definition: &ResourceDefinition{
			Name: "Category", Public: true, Tree: true,
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "parent", Type: "belongs_to", RelatedModel: "Category"},
			},
		},
	}
	src := g.publicHandlerSource(g.Names(), firstReturn(PublicFields(g.Definition.Fields)))
	if n := strings.Count(src, `json:"parent_id"`); n != 1 {
		t.Errorf("parent_id appears %d times, want 1", n)
	}
}

// grit#91 item 1: the generated hooks called the authenticated routes and
// nothing sent the publishable key, so a storefront got 401s.
func TestPublicWebReadsCallThePublicRoutesWithTheKey(t *testing.T) {
	g := catalogueGenerator()
	src := g.publicWebReadsSource(g.Names(), true)
	for _, want := range []string{
		`const API_KEY = process.env.NEXT_PUBLIC_API_KEY ?? "";`,
		`headers: API_KEY ? { "X-API-Key": API_KEY } : {}`,
		`"/api/v1/public/products"`,
		"export const getPublicProducts",
		"export const getPublicProduct ",
		// The resource has a belongs_to, so it has a related endpoint.
		"export const getRelatedProducts",
		// Typed against the allowlist, not the model.
		"export interface PublicProduct {",
		"  category_id: string;",
		"  price: number;",
		// A FileRef object, which is what the Go view struct carries. As a
		// string this compiled and rendered [object Object].
		"  image: FileRef | null;",
		`import type { FileRef, PaginatedResponse } from "@repo/shared/types";`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the public web reads are missing %q", want)
		}
	}
	// Held back on the Go side, so absent here too, or the type lies.
	for _, absent := range []string{"  stock:", "  active:"} {
		if strings.Contains(src, absent) {
			t.Errorf("%q is in the public type and not in the public response", absent)
		}
	}
	// Next reads these on the server, where the container reaches the API by
	// its compose service name.
	if !strings.Contains(src, "process.env.API_INTERNAL_URL") || !strings.Contains(src, "REVALIDATE_SECONDS") {
		t.Error("the Next.js variant has no server-side URL or revalidation")
	}
}

func TestPublicWebReadsForVite(t *testing.T) {
	g := catalogueGenerator()
	src := g.publicWebReadsSource(g.Names(), false)
	if !strings.Contains(src, "import.meta.env.VITE_API_KEY") {
		t.Error("the Vite variant does not read VITE_API_KEY")
	}
	for _, absent := range []string{"REVALIDATE_SECONDS", `from "react"`, "process.env"} {
		if strings.Contains(src, absent) {
			t.Errorf("the Vite variant carries %q, which only Next has", absent)
		}
	}
}

func TestPublicWebReadsSkipTheEndpointsThatDoNotExist(t *testing.T) {
	g := &Generator{
		Module: "shopfront/apps/api",
		Definition: &ResourceDefinition{
			Name: "Page", Public: true,
			Fields: []Field{{Name: "title", Type: "string"}, {Name: "slug", Type: "slug"}},
		},
	}
	src := g.publicWebReadsSource(g.Names(), true)
	for _, absent := range []string{"getRelatedPages", "getPublicPagesTree"} {
		if strings.Contains(src, absent) {
			t.Errorf("a flat resource grew %s, which the API does not serve", absent)
		}
	}
}

func firstReturn(included, _ []Field) []Field { return included }
