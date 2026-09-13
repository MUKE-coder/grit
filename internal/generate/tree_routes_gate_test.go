package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func treeRoutesFor(t *testing.T, routesGo string, tenantOwned bool) string {
	t.Helper()
	const module = "shop/apps/api"
	def, err := ParseInlineFields("Category", "name:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	def.Tree = true
	def.TenantOwned = tenantOwned
	g := newTestGenerator(setupMinimalProject(t, module), module, def)
	path := filepath.Join(g.APIRoot(), "internal", "routes", "routes.go")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(routesGo), 0o644); err != nil {
		t.Fatal(err)
	}
	g.ensureTreeRoutes(g.Names())
	return readTestFile(t, path)
}

const staffRoutesGo = "package routes\n\nfunc Setup() {\n\tstaff.Use(middleware.RequireStaff())\n\t// grit:routes:custom\n}\n"

// --tree mounted move, reorder and rebuild-tree on the protected group, so any
// signed-in account could rearrange a shared tree or rewrite every path in it.
func TestTreeRoutesAskForAPermission(t *testing.T) {
	src := treeRoutesFor(t, staffRoutesGo, false)
	for _, want := range []string{
		`staff.GET("/categories/tree", middleware.RequireRole("ADMIN", "perm:categories.view"), categoryTreeHandler.GetTree)`,
		`staff.GET("/categories/:id/breadcrumbs", middleware.RequireRole("ADMIN", "perm:categories.view"), categoryTreeHandler.GetBreadcrumbs)`,
		`staff.PATCH("/categories/:id/move", middleware.RequireRole("ADMIN", "perm:categories.edit"), categoryTreeHandler.Move)`,
		`staff.POST("/categories/reorder", middleware.RequireRole("ADMIN", "perm:categories.edit"), categoryTreeHandler.Reorder)`,
		`staff.POST("/categories/rebuild-tree", middleware.RequireRole("ADMIN", "perm:categories.edit"), categoryTreeHandler.RebuildPaths)`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("routes.go is missing %s\n%s", want, src)
		}
	}
	if strings.Contains(src, "protected.") {
		t.Errorf("a tree route is open to any signed-in account:\n%s", src)
	}
}

func TestTreeRoutesWithoutAStaffGroupAreAdmins(t *testing.T) {
	src := treeRoutesFor(t, "package routes\n\nfunc Setup() {\n\t// grit:routes:custom\n}\n", false)
	if !strings.Contains(src, `admin.PATCH("/categories/:id/move", categoryTreeHandler.Move)`) {
		t.Errorf("tree routes did not go to the admin group:\n%s", src)
	}
}

func TestTenantOwnedTreeRoutesStayProtected(t *testing.T) {
	src := treeRoutesFor(t, staffRoutesGo, true)
	if !strings.Contains(src, `protected.PATCH("/categories/:id/move", categoryTreeHandler.Move)`) {
		t.Errorf("a tenant-owned tree was moved off the protected group:\n%s", src)
	}
}
