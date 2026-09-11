package generate

import (
	"strings"
	"testing"
)

// The --public endpoints query through the resource's service.
//
// The public handler queried h.DB directly, so its reads carried no request
// context, and the multitenant plugin refuses a tenant-owned query that has no
// organization on it. What a client may filter by stays in the public handler,
// where the developer edits it; what the queries are is the service's, like
// every other query the resource has.
func TestPublicQueriesAreTheServices(t *testing.T) {
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
	svc := withParent.publicServiceMethods(withParent.Names())
	for _, want := range []string{
		"func (s *ProductService) ListPublic(ctx context.Context, params paginate.Params, cfg paginate.Config)",
		"func (s *ProductService) GetPublic(ctx context.Context, key string)",
		"func (s *ProductService) RelatedPublic(ctx context.Context, item *models.Product, limit int)",
		// The model spells it CategoryID, not CategoryId.
		`if item.CategoryID != "" {`,
		`query = query.Where("category_id = ?", item.CategoryID)`,
		// Excludes itself, or the strip shows the page you are already on.
		`Where("id <> ? AND archived_at IS NULL", item.ID)`,
	} {
		if !strings.Contains(svc, want) {
			t.Errorf("the service is missing %q", want)
		}
	}
	if strings.Contains(svc, "s.DB.") {
		t.Error("a public query bypasses the request context")
	}

	// No parent, no related query; no tree, no tree queries.
	flat := &Generator{
		Module: "shopfront/apps/api",
		Definition: &ResourceDefinition{
			Name: "Page", Public: true,
			Fields: []Field{{Name: "title", Type: "string"}, {Name: "slug", Type: "slug"}},
		},
	}
	svc = flat.publicServiceMethods(flat.Names())
	for _, absent := range []string{"RelatedPublic", "PublicSubtreeIDs", "PublicTreeRows"} {
		if strings.Contains(svc, absent) {
			t.Errorf("a flat resource with no parent grew %s", absent)
		}
	}

	// A tree publishes its shape, and its parent is a nullable pointer.
	tree := &Generator{
		Module: "shopfront/apps/api",
		Definition: &ResourceDefinition{
			Name: "Category", Public: true, Tree: true,
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "slug", Type: "slug"},
				{Name: "parent", Type: "belongs_to", RelatedModel: "Category"},
			},
		},
	}
	svc = tree.publicServiceMethods(tree.Names())
	for _, want := range []string{
		"func (s *CategoryService) PublicSubtreeIDs(ctx context.Context, path string)",
		"func (s *CategoryService) PublicTreeRows(ctx context.Context)",
		`if item.ParentID != nil && *item.ParentID != "" {`,
	} {
		if !strings.Contains(svc, want) {
			t.Errorf("the tree's public service is missing %q", want)
		}
	}

	// Without --public, the service has none of it.
	withParent.Definition.Public = false
	if withParent.publicServiceMethods(withParent.Names()) != "" {
		t.Error("a resource without --public grew public queries")
	}
}
