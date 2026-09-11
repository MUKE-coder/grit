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

// Without scoping, "protected" means any account you have ever issued.
//
// Generated routes sit on the authenticated group, so before --owned-by every
// signed-in user could list, read, edit and delete every row of every generated
// resource. Reproduced on a stock scaffold: a second account read the first
// account's record straight off GET /api/v1/conversations/:id.
func TestOwnedResourceScopesEveryAccessPath(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	g, names := ownedGenerator(t, root, module)

	if err := g.writeGoModel(names); err != nil {
		t.Fatalf("model: %v", err)
	}
	if err := g.writeGoHandler(names); err != nil {
		t.Fatalf("handler: %v", err)
	}

	model := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "invoice.go"))
	if !strings.Contains(model, "func (m *Invoice) GetOwnerID() string") {
		t.Error("the model does not implement authz.Ownable, so the ownership " +
			"helpers have nothing to compare against")
	}

	h := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "invoice.go"))

	// The list is the path people forget. Guarding only the by-id routes still
	// hands the whole table over on GET /invoices.
	if !strings.Contains(h, `authz.ScopeToOwner(c, query, "user_id")`) {
		t.Error("List is not scoped to the owner, so every row is returned to " +
			"every authenticated caller")
	}
	// GetByID, PDF, Update, Patch and Delete each load a row by id and each
	// need the check. PDF and Patch were generated without it.
	if n := strings.Count(h, "authz.OwnsOr404(c, &item)"); n != 5 {
		t.Errorf("ownership is checked on %d of the 5 by-id handlers "+
			"(GetByID, PDF, Update, Patch, Delete)", n)
	}
	if !strings.Contains(h, "item.UserID = authz.CurrentUserID(c)") {
		t.Error("Create does not stamp the owner from the session")
	}
	if !strings.Contains(h, `"shop/apps/api/internal/authz"`) {
		t.Error("the authz import is missing, so the handler will not compile")
	}
}

// The owner comes from the session, never from the request body.
//
// Generated as an ordinary belongs_to, the field arrived with
// binding:"required", so create failed with a validation error demanding a
// user_id, and a caller who supplied one would have been naming the owner of
// a row they were about to create. Caught by running it, not by reading it.
func TestOwnerFieldIsNotAcceptedFromTheClient(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	g, names := ownedGenerator(t, root, module)

	if err := g.writeGoHandler(names); err != nil {
		t.Fatalf("handler: %v", err)
	}
	h := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "invoice.go"))

	// Both request structs, and the PATCH whitelist, must be silent about it.
	for _, forbidden := range []string{
		"UserID string `json:\"user_id\" binding:\"required\"`",
		"UserID *string `json:\"user_id\"`",
		"\"user_id\": true,",
	} {
		if strings.Contains(h, forbidden) {
			t.Errorf("the request surface still carries %q: the owner is "+
				"server-assigned and a client that can set it can hand rows to "+
				"other accounts", forbidden)
		}
	}

	// It is still preloaded, so reads return the owner.
	if !strings.Contains(h, "Preload(\"User\")") {
		t.Error("the owner association is no longer preloaded, so responses " +
			"lost the user object")
	}
}

// A resource generated without the flag keeps the old behaviour exactly.
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
	h := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "product.go"))

	if strings.Contains(h, "authz.") {
		t.Error("a resource generated without --owned-by grew ownership checks; " +
			"a shared catalogue would become invisible to everyone but its creator")
	}
}
