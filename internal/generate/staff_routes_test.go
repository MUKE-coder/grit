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

// An append-only resource has no delete or bulk route, so importing
// middleware for them would be an unused import.
func TestAppendOnlyGetsNoStrayImport(t *testing.T) {
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
	if strings.Contains(src, "internal/middleware") {
		t.Error("an append-only resource imports middleware it does not use")
	}
	mustParse(t, "entry_routes.go", src)
}
