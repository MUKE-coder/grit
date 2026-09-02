package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// withRouteRegistry adds internal/routes/resources.go, which is what tells the
// generator this project takes the per-resource split.
func withRouteRegistry(t *testing.T, root string) {
	t.Helper()
	writeTestFile(t, filepath.Join(root, "apps", "api", "internal", "routes", "resources.go"),
		"package routes\n\ntype Mount struct{}\n\nfunc RegisterRoutes(fn func(*Mount)) {}\n")
}

// The point of the split: routes.go stops growing.
//
// It used to take four edits per resource, hundreds of lines apart, in a file
// that passed a thousand lines with a handful of resources.
func TestResourceRoutesGoToTheirOwnFile(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	withRouteRegistry(t, root)

	routesPath := filepath.Join(root, "apps", "api", "internal", "routes", "routes.go")
	before := readTestFile(t, routesPath)

	def, err := ParseInlineFields("Product", "name:string,price:money")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	if err := newTestGenerator(root, module, def).Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "routes", "product_routes.go"))
	assertContains(t, "product_routes.go", src,
		"package routes",
		"func init() {",
		"RegisterRoutes(func(m *Mount) {",
		"h := &handlers.ProductHandler{",
		`m.Protected.GET("/products", h.List)`,
		`m.Protected.GET("/products/:id", h.GetByID)`,
		// Bulk sits with DELETE because it can delete.
		`m.Admin.DELETE("/products/:id", h.Delete)`,
		`m.Admin.POST("/products/bulk", h.Bulk)`,
		module+"/internal/handlers",
	)
	// No roles, so middleware is never named, and an unused import does not
	// compile in Go.
	if strings.Contains(src, "internal/middleware") {
		t.Error("imported middleware without using it")
	}

	after := readTestFile(t, routesPath)
	for _, leaked := range []string{"productHandler", "/products"} {
		if strings.Contains(after, leaked) {
			t.Errorf("routes.go gained %q; the whole point is that it does not", leaked)
		}
	}
	if strings.Count(after, "\n") > strings.Count(before, "\n")+2 {
		t.Errorf("routes.go grew by more than the studio and sync one-liners: %d -> %d lines",
			strings.Count(before, "\n"), strings.Count(after, "\n"))
	}
}

// A role-restricted resource puts its guard in the same file as the routes it
// guards, so the two cannot drift apart.
func TestRoleRestrictedResourceRoutes(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	withRouteRegistry(t, root)

	def, err := ParseInlineFields("Report", "title:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	g.Roles = []string{"ADMIN", "EDITOR"}
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "routes", "report_routes.go"))
	assertContains(t, "report_routes.go", src,
		`g := m.Protected.Group("/reports")`,
		`g.Use(middleware.RequireRole("ADMIN", "EDITOR"))`,
		`g.DELETE("/:id", h.Delete)`,
		module+"/internal/middleware",
	)
	// Every route on the guarded group. A resource where DELETE escaped onto
	// m.Admin would be guarded by role in one place and not the other.
	if strings.Contains(src, "m.Admin.") {
		t.Error("a role-restricted resource put routes on m.Admin, outside its own guard")
	}
}

// Deleting the file is the whole removal, and it must not disturb routes.go.
func TestRemovingASplitResourceDeletesItsFile(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	withRouteRegistry(t, root)

	def, err := ParseInlineFields("Product", "name:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	if err := newTestGenerator(root, module, def).Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	apiRoot := filepath.Join(root, "apps", "api")
	path := filepath.Join(apiRoot, "internal", "routes", "product_routes.go")
	if !fileExists(path) {
		t.Fatal("nothing to remove: the routes file was never written")
	}

	removed, err := removeResourceRoutes(apiRoot, "product")
	if err != nil {
		t.Fatalf("removeResourceRoutes: %v", err)
	}
	if !removed {
		t.Error("removeResourceRoutes reported nothing removed")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("product_routes.go is still there")
	}

	// Removing it twice is not an error: a pre-split project has no such file
	// and `grit remove resource` runs this unconditionally.
	removed, err = removeResourceRoutes(apiRoot, "product")
	if err != nil {
		t.Fatalf("second removeResourceRoutes: %v", err)
	}
	if removed {
		t.Error("reported a removal on the second pass")
	}
}

// A project without the registry keeps the old behaviour exactly, because
// upgrading is the user's decision and a generate must not depend on it.
func TestProjectWithoutRegistryStillInjectsIntoRoutes(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Product", "name:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	if err := newTestGenerator(root, module, def).Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	if fileExists(filepath.Join(root, "apps", "api", "internal", "routes", "product_routes.go")) {
		t.Error("wrote a split routes file into a project with no registry to load it")
	}
	routes := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "routes", "routes.go"))
	assertContains(t, "routes.go", routes, "productHandler", "/products")
}
