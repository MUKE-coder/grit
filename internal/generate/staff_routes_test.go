package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeRegistry(t *testing.T, g *Generator, withStaff bool) {
	t.Helper()
	body := "package routes\n\ntype Mount struct {\n\tAdmin     *gin.RouterGroup\n}\n"
	if withStaff {
		body = "package routes\n\ntype Mount struct {\n\tAdmin     *gin.RouterGroup\n\tStaff     *gin.RouterGroup\n}\n"
	}
	path := filepath.Join(g.APIRoot(), "internal", "routes", "resources.go")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// A role granted products.delete was refused by the delete route: it sat on
// the ADMIN-only group with no permission of its own.
func TestDeleteAndBulkAskForTheResourcePermission(t *testing.T) {
	const module = "shop/apps/api"
	def, err := ParseInlineFields("Product", "name:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(setupMinimalProject(t, module), module, def)
	writeRegistry(t, g, true)
	src, err := g.resourceRoutesSource(g.Names())
	if err != nil {
		t.Fatalf("routes: %v", err)
	}
	for _, want := range []string{
		`m.Staff.DELETE("/products/:id", middleware.RequireRole("ADMIN", "perm:products.delete"), h.Delete)`,
		`m.Staff.POST("/products/bulk", middleware.RequireRole("ADMIN", "perm:products.delete"), h.Bulk)`,
		`"shop/apps/api/internal/middleware"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the route file is missing %s", want)
		}
	}
	mustParse(t, "product_routes.go", src)
}

// A project whose registry predates the staff group keeps ADMIN-only deletes,
// rather than generating code that does not compile.
func TestOlderProjectsKeepTheAdminGroup(t *testing.T) {
	const module = "shop/apps/api"
	def, err := ParseInlineFields("Product", "name:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(setupMinimalProject(t, module), module, def)
	writeRegistry(t, g, false)
	src, err := g.resourceRoutesSource(g.Names())
	if err != nil {
		t.Fatalf("routes: %v", err)
	}
	if strings.Contains(src, "m.Staff") || strings.Contains(src, "internal/middleware") {
		t.Error("an older project got the staff group, which its Mount does not have")
	}
	if !strings.Contains(src, `m.Admin.DELETE("/products/:id", h.Delete)`) {
		t.Error("the admin delete route is gone")
	}
}

// An owned append-only resource has no delete or bulk route and protects the
// rest, so importing middleware would be an unused import.
func TestAppendOnlyGetsNoStrayImport(t *testing.T) {
	const module = "shop/apps/api"
	def, err := ParseInlineFields("Entry", "memo:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	def.AppendOnly = true
	def.OwnedBy = "user"
	def.Fields = append(def.Fields, Field{Name: "user", Type: "belongs_to", RelatedModel: "User"})
	g := newTestGenerator(setupMinimalProject(t, module), module, def)
	writeRegistry(t, g, true)
	src, err := g.resourceRoutesSource(g.Names())
	if err != nil {
		t.Fatalf("routes: %v", err)
	}
	if strings.Contains(src, "internal/middleware") {
		t.Error("an append-only resource imports middleware it does not use")
	}
	mustParse(t, "entry_routes.go", src)
}

// C3 in the contact-app review: a generated resource's list, get, export,
// import, create and update were open to any signed-in account, which on an app
// with open registration is anybody. The review exported a whole address book
// that way. Each verb now asks for the permission the roles UI grants for it.
func TestSharedResourceAsksForAPermissionPerVerb(t *testing.T) {
	const module = "shop/apps/api"
	def, err := ParseInlineFields("Contact", "name:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(setupMinimalProject(t, module), module, def)
	writeRegistry(t, g, true)
	src, err := g.resourceRoutesSource(g.Names())
	if err != nil {
		t.Fatalf("routes: %v", err)
	}
	for _, want := range []string{
		`m.Staff.GET("/contacts", middleware.RequireRole("ADMIN", "perm:contacts.view"), h.List)`,
		`m.Staff.GET("/contacts/export", middleware.RequireRole("ADMIN", "perm:contacts.view"), h.Export)`,
		`m.Staff.POST("/contacts/import", middleware.RequireRole("ADMIN", "perm:contacts.create"), h.Import)`,
		`m.Staff.GET("/contacts/:id", middleware.RequireRole("ADMIN", "perm:contacts.view"), h.GetByID)`,
		`m.Staff.POST("/contacts", middleware.RequireRole("ADMIN", "perm:contacts.create"), h.Create)`,
		`m.Staff.PUT("/contacts/:id", middleware.RequireRole("ADMIN", "perm:contacts.edit"), h.Update)`,
		`m.Staff.PATCH("/contacts/:id", middleware.RequireRole("ADMIN", "perm:contacts.edit"), h.Patch)`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the route file is missing %s", want)
		}
	}
	if strings.Contains(src, "m.Protected.") {
		t.Errorf("a shared resource still has a route any signed-in account reaches:\n%s", src)
	}
	mustParse(t, "contact_routes.go", src)
}

// An owned resource's service scopes every query to the caller, so its routes
// stay open to signed-in accounts: that is what lets a customer read their own
// invoices. Delete keeps its permission.
func TestOwnedResourceStaysProtected(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	g, names := ownedGenerator(t, root, module)
	writeRegistry(t, g, true)
	src, err := g.resourceRoutesSource(names)
	if err != nil {
		t.Fatalf("routes: %v", err)
	}
	for _, want := range []string{
		`m.Protected.GET("/invoices", h.List)`,
		`m.Protected.POST("/invoices", h.Create)`,
		`m.Protected.PUT("/invoices/:id", h.Update)`,
		`m.Staff.DELETE("/invoices/:id", middleware.RequireRole("ADMIN", "perm:invoices.delete"), h.Delete)`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the route file is missing %s", want)
		}
	}
	mustParse(t, "invoice_routes.go", src)
}

// A shared append-only resource asks for view and create like any other.
func TestAppendOnlySharedResourceAsksForAPermission(t *testing.T) {
	const module = "shop/apps/api"
	def, err := ParseInlineFields("Entry", "memo:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	def.AppendOnly = true
	g := newTestGenerator(setupMinimalProject(t, module), module, def)
	writeRegistry(t, g, true)
	src, err := g.resourceRoutesSource(g.Names())
	if err != nil {
		t.Fatalf("routes: %v", err)
	}
	plural := g.Names().Plural
	for _, want := range []string{
		`m.Staff.GET("/` + plural + `", middleware.RequireRole("ADMIN", "perm:` + plural + `.view"), h.List)`,
		`m.Staff.POST("/` + plural + `", middleware.RequireRole("ADMIN", "perm:` + plural + `.create"), h.Create)`,
		`"shop/apps/api/internal/middleware"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the route file is missing %s", want)
		}
	}
	mustParse(t, "entry_routes.go", src)
}

// A tenant-owned resource is narrowed to the caller's organization by the tenant
// middleware, the way an owned one is narrowed to its owner, so its routes stay
// protected: members work with their organization's rows.
func TestTenantOwnedResourceStaysProtected(t *testing.T) {
	const module = "shop/apps/api"
	def, err := ParseInlineFields("Deal", "title:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	def.TenantOwned = true
	g := newTestGenerator(setupMinimalProject(t, module), module, def)
	writeRegistry(t, g, true)
	src, err := g.resourceRoutesSource(g.Names())
	if err != nil {
		t.Fatalf("routes: %v", err)
	}
	if !strings.Contains(src, `m.Protected.GET("/deals", h.List)`) || strings.Contains(src, "perm:deals.view") {
		t.Errorf("a tenant-owned resource was moved off the protected group:\n%s", src)
	}
	mustParse(t, "deal_routes.go", src)
}
