package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// --tenant-owned embeds tenant.Owned and imports the package.
//
// Without the flag the multitenant plugin's instruction ("embed tenant.Owned")
// had to be carried out by hand on every generated model, and again after any
// --force regeneration. Nothing warned when it was missed: scoping is per
// model, so a table without OrgID is simply never filtered. The endpoint works,
// returns rows, and returns everyone's. Found by building a CRM whose four
// generated resources were four shared tables.
func TestTenantOwnedEmbedsTheScopingMixin(t *testing.T) {
	const module = "crm/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Contact", "name:string,email:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	def.TenantOwned = true
	g := newTestGenerator(root, module, def)
	if err := g.writeGoModel(g.Names()); err != nil {
		t.Fatalf("model: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "contact.go"))

	if !strings.Contains(src, `"crm/apps/api/internal/tenant"`) {
		t.Error("the tenant package is not imported, so the model will not compile")
	}
	if !strings.Contains(src, "\ttenant.Owned\n") {
		t.Errorf("tenant.Owned is not embedded, so this model is never scoped:\n%s", src)
	}
	// Directly under ID, where someone asking "which tenant owns this row"
	// looks first.
	idIdx := strings.Index(src, `json:"id"`)
	ownedIdx := strings.Index(src, "tenant.Owned")
	nameIdx := strings.Index(src, "Name ")
	if !(idIdx < ownedIdx && ownedIdx < nameIdx) {
		t.Error("tenant.Owned should sit between ID and the resource's own fields")
	}
}

// Without the flag, nothing changes. A shared reference table scoped by
// accident becomes invisible to every query, with no error to explain it.
func TestResourcesAreSharedUnlessAsked(t *testing.T) {
	const module = "crm/apps/api"
	root := setupMinimalProject(t, module)

	def, _ := ParseInlineFields("Country", "name:string,code:string")
	g := newTestGenerator(root, module, def)
	if err := g.writeGoModel(g.Names()); err != nil {
		t.Fatalf("model: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "country.go"))

	if strings.Contains(src, "tenant.Owned") || strings.Contains(src, "internal/tenant") {
		t.Error("a resource generated without --tenant-owned grew organization " +
			"scoping; a shared catalogue would become invisible to every query")
	}
}

// The flag composes with the other model-shaping flags.
func TestTenantOwnedComposesWithTreeAndOwnedBy(t *testing.T) {
	const module = "crm/apps/api"
	root := setupMinimalProject(t, module)

	def, _ := ParseInlineFields("Deal", "title:string,user:belongs_to:User")
	def.TenantOwned = true
	def.OwnedBy = "user"
	g := newTestGenerator(root, module, def)
	if err := g.writeGoModel(g.Names()); err != nil {
		t.Fatalf("model: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "deal.go"))

	// Organization scoping and per-user ownership are different questions and
	// a deal in a CRM needs both: the org decides which company's data it is,
	// the owner decides which rep may touch it.
	if !strings.Contains(src, "tenant.Owned") {
		t.Error("tenant scoping was dropped when combined with --owned-by")
	}
	if !strings.Contains(src, "func (m *Deal) GetOwnerID() string") {
		t.Error("ownership was dropped when combined with --tenant-owned")
	}
}
