package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// ownedGenerator builds a resource scoped to its owner, the way the --owned-by
// flag does: the flag adds the belongs_to-User field when it is absent.
func ownedGenerator(t *testing.T, root, module string) (*Generator, Names) {
	t.Helper()
	def, err := ParseInlineFields("Invoice", "number:string,paid:bool")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	def.OwnedBy = "user"
	def.Fields = append(def.Fields, Field{
		Name: "user", Type: "belongs_to", RelatedModel: "User",
	})
	g := newTestGenerator(root, module, def)
	return g, g.Names()
}

// ownedFiles writes and reads back an owned resource's handler and service.
func ownedFiles(t *testing.T) (*Generator, string, string, string) {
	t.Helper()
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	g, names := ownedGenerator(t, root, module)
	if err := g.writeGoHandler(names); err != nil {
		t.Fatalf("handler: %v", err)
	}
	if err := g.writeGoService(names); err != nil {
		t.Fatalf("service: %v", err)
	}
	h := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "invoice.go"))
	s := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "services", "invoice.go"))
	mustParse(t, "handlers/invoice.go", h)
	mustParse(t, "services/invoice.go", s)
	return g, root, h, s
}

// Without scoping, "protected" means any account you have ever issued.
//
// Generated routes sit on the authenticated group, so before --owned-by every
// signed-in user could list, read, edit and delete every row of every generated
// resource. Reproduced on a stock scaffold: a second account read the first
// account's record straight off GET /api/v1/conversations/:id.
func TestOwnedResourceScopesEveryAccessPath(t *testing.T) {
	g, root, h, s := ownedFiles(t)
	if err := g.writeGoModel(g.Names()); err != nil {
		t.Fatalf("model: %v", err)
	}
	model := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "invoice.go"))
	if !strings.Contains(model, "func (m *Invoice) GetOwnerID() string") {
		t.Error("the model does not implement authz.Ownable, so the ownership " +
			"helpers have nothing to compare against")
	}

	// The list is the path people forget, and the export is the list without
	// pages. Guarding only the by-id routes still hands the whole table over.
	if n := strings.Count(s, `authz.ScopeOwned(ctx, query, "user_id")`); n != 2 {
		t.Errorf("List and Export are scoped %d time(s), want both", n)
	}
	if !strings.Contains(method(t, s, "func (s *InvoiceService) Bulk("), `authz.ScopeOwned(ctx, scope, "user_id")`) {
		t.Error("Bulk acts on anyone's rows")
	}

	// Every row fetched by id comes through a guarded read: GetByID for GetByID
	// and PDF, load for Update, Patch and Delete. PDF and Patch were once
	// generated without the check.
	for _, sig := range []string{"func (s *InvoiceService) GetByID(", "func (s *InvoiceService) load("} {
		if !strings.Contains(method(t, s, sig), "if !authz.Owns(ctx, &item) {") {
			t.Errorf("%s returns somebody else's row", sig)
		}
	}
	for _, fn := range []string{"Update", "Patch", "Delete"} {
		if !strings.Contains(method(t, s, "func (s *InvoiceService) "+fn+"("), "s.load(ctx, id)") {
			t.Errorf("%s does not load its row through the guarded read", fn)
		}
	}
	for _, fn := range []string{"GetByID", "PDF"} {
		if !strings.Contains(method(t, h, "func (h *InvoiceHandler) "+fn+"("), "h.service().GetByID(") {
			t.Errorf("the %s handler does not read through the guarded GetByID", fn)
		}
	}

	if !strings.Contains(method(t, s, "func (s *InvoiceService) Create("), "item.UserID = authz.UserIDFrom(ctx)") {
		t.Error("Create does not stamp the owner from the caller")
	}
	if !strings.Contains(s, `"shop/apps/api/internal/authz"`) {
		t.Error("the authz import is missing, so the service will not compile")
	}
}

// The owner comes from the session, never from the request body.
//
// Generated as an ordinary belongs_to, the field arrived with
// binding:"required", so create failed with a validation error demanding a
// user_id, and a caller who supplied one would have been naming the owner of
// a row they were about to create. Caught by running it, not by reading it.
func TestOwnerFieldIsNotAcceptedFromTheClient(t *testing.T) {
	_, _, h, s := ownedFiles(t)

	// Both request structs, and the writable columns PATCH and bulk use.
	for name, src := range map[string]string{"handler": h, "service": s} {
		for _, forbidden := range []string{
			"UserID string `json:\"user_id\" binding:\"required\"`",
			"UserID *string `json:\"user_id\"`",
			"\"user_id\": true,",
		} {
			if strings.Contains(src, forbidden) {
				t.Errorf("the %s still carries %q: the owner is server-assigned, and "+
					"a client that can set it can hand rows to other accounts", name, forbidden)
			}
		}
	}

	// It is still preloaded, so reads return the owner.
	if !strings.Contains(s, "Preload(\"User\")") {
		t.Error("the owner association is no longer preloaded, so responses lost the user object")
	}
}

// A resource generated without the flag has no ownership logic at all.
func TestSharedResourcesAreUnchanged(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Product", "name:string,price:money")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	names := g.Names()
	if err := g.writeGoHandler(names); err != nil {
		t.Fatalf("handler: %v", err)
	}
	if err := g.writeGoService(names); err != nil {
		t.Fatalf("service: %v", err)
	}
	h := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "product.go"))
	s := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "services", "product.go"))

	for _, owned := range []string{`"shop/apps/api/internal/authz"`, "authz.ScopeOwned(", "authz.Owns(", "authz.UserIDFrom("} {
		if strings.Contains(s, owned) {
			t.Errorf("a resource generated without --owned-by grew ownership checks (%s); "+
				"a shared catalogue would become invisible to everyone but its creator", owned)
		}
	}
	if strings.Contains(h, "ScopeOwned") || strings.Contains(h, "authz.Owns") {
		t.Error("the handler of a shared resource scopes by owner")
	}
}
