package scaffold

import (
	"strings"
	"testing"
)

// A --tree resource's routes went into routes.go on the protected group, where
// any signed-in account could move nodes or rebuild every path.
func TestRepairTreeRoutes(t *testing.T) {
	src := "package routes\n\nfunc Setup() {\n" +
		"\tstaff.Use(middleware.RequireStaff())\n" +
		"\t// Custom role-restricted routes\n" +
		"\t\tcategoryTreeHandler := handlers.NewCategoryTreeHandler(db)\n" +
		"\t\tprotected.GET(\"/categories/tree\", categoryTreeHandler.GetTree)\n" +
		"\t\tprotected.PATCH(\"/categories/:id/move\", categoryTreeHandler.Move)\n" +
		"\t\tprotected.POST(\"/categories/rebuild-tree\", categoryTreeHandler.RebuildPaths)\n" +
		"\t\tprotected.GET(\"/regions/tree\", regionTreeHandler.GetTree)\n" +
		"\t// grit:routes:custom\n}\n"
	out, fixed, warn := repairTreeRoutesSource(src, map[string]bool{"Region": true})
	if len(warn) > 0 {
		t.Errorf("warnings: %v", warn)
	}
	mustFormat(t, "routes.go", out)
	for _, want := range []string{
		`staff.GET("/categories/tree", middleware.RequireRole("ADMIN", "perm:categories.view"), categoryTreeHandler.GetTree)`,
		`staff.PATCH("/categories/:id/move", middleware.RequireRole("ADMIN", "perm:categories.edit"), categoryTreeHandler.Move)`,
		`staff.POST("/categories/rebuild-tree", middleware.RequireRole("ADMIN", "perm:categories.edit"), categoryTreeHandler.RebuildPaths)`,
		// Tenant-owned, so scoped already.
		`protected.GET("/regions/tree", regionTreeHandler.GetTree)`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repaired routes.go is missing %s", want)
		}
	}
	if len(fixed) != 1 || !strings.Contains(fixed[0], "3 tree routes") {
		t.Errorf("reported %v, want the 3 routes it moved", fixed)
	}
	if again, fixed, _ := repairTreeRoutesSource(out, map[string]bool{"Region": true}); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed routes.go again")
	}
}

func TestRepairTreeRoutesWithoutAStaffGroup(t *testing.T) {
	src := "package routes\n\nfunc Setup() {\n\t\tprotected.PATCH(\"/categories/:id/move\", categoryTreeHandler.Move)\n}\n"
	out, _, _ := repairTreeRoutesSource(src, map[string]bool{})
	if !strings.Contains(out, `admin.PATCH("/categories/:id/move", categoryTreeHandler.Move)`) {
		t.Errorf("the tree route did not go to the admin group:\n%s", out)
	}
}
