package scaffold

import (
	"strings"
	"testing"
)

const routesWithProtectedVariants = `package routes

func Setup() {
	protected := v1.Group("")
	staff := v1.Group("")

	// Custom role-restricted routes
		productVariantHandler := handlers.NewProductVariantHandler(db)
		protected.GET("/options", productVariantHandler.ListOptions)
		protected.POST("/options", productVariantHandler.CreateOption)
		protected.GET("/products/:id/variants", productVariantHandler.List)
		protected.PATCH("/product-variants/:id", productVariantHandler.Update)
	// grit:routes:custom
	_ = protected
	_ = staff
}
`

// A customer must not be able to reprice a variant: every variant route asks
// for the resource's permissions, reads for view and writes for edit.
func TestVariantRoutesMoveToTheStaffGroup(t *testing.T) {
	out, changes, warnings := VariantRoutesOnStaff(routesWithProtectedVariants)
	if len(changes) != 1 || len(warnings) != 0 {
		t.Fatalf("changes %v, warnings %v", changes, warnings)
	}
	for _, want := range []string{
		`	staff.GET("/options", middleware.RequireRole("ADMIN", "perm:products.view"), productVariantHandler.ListOptions)`,
		`	staff.POST("/options", middleware.RequireRole("ADMIN", "perm:products.edit"), productVariantHandler.CreateOption)`,
		`	staff.PATCH("/product-variants/:id", middleware.RequireRole("ADMIN", "perm:products.edit"), productVariantHandler.Update)`,
		"\n\tproductVariantHandler := handlers.NewProductVariantHandler(db)\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in\n%s", want, out)
		}
	}
	if strings.Contains(out, "protected.GET(\"/options\"") {
		t.Error("a variant route is still on the protected group")
	}
	if again, changes, _ := VariantRoutesOnStaff(out); again != out || len(changes) != 0 {
		t.Error("a second run changed something")
	}
}

// Without a staff group ahead of them there is nowhere safe to move them, and
// saying so beats writing code that does not compile.
func TestVariantRoutesWithoutAStaffGroupWarn(t *testing.T) {
	src := strings.Replace(routesWithProtectedVariants, "\tstaff := v1.Group(\"\")\n", "", 1)
	out, changes, warnings := VariantRoutesOnStaff(src)
	if out != src || len(changes) != 0 || len(warnings) != 1 {
		t.Errorf("changes %v, warnings %v", changes, warnings)
	}
}
