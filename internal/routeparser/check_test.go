package routeparser

import (
	"strings"
	"testing"
)

// declSource is a routes file using the three declarations, with the groups
// the generated routes.go declares.
const declSource = `package routes

func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	protected := v1.Group("")
	protected.Use(middleware.Auth(cfg))
	staff := protected.Group("")
	staff.Use(middleware.RequireStaff())

	// Declared with a permission.
	staff.POST("/payslips/:id/email", perm("payroll.edit"), h.EmailOne)
	staff.GET("/payslips", perm("payroll.view", "payroll.admin"), h.List)

	// Declared open, with the reason.
	protected.GET("/payslips/mine", openToAnySignedIn("the handler scopes to the caller"), h.Mine)

	// A role rather than a permission.
	staff.POST("/payroll/run", role("FINANCE"), h.Run)

	// Forgotten.
	protected.POST("/payslips/:id/void", h.Void)
}
`

func routeNamed(t *testing.T, routes []Route, method, path string) Route {
	t.Helper()
	for _, r := range routes {
		if r.Method == method && r.Path == path {
			return r
		}
	}
	t.Fatalf("no %s %s in the parse", method, path)
	return Route{}
}

// perm("x") has to mean exactly what the generator's
// RequireRole("ADMIN", "perm:x") means, or there are two spellings of one
// requirement and they drift (grit#179).
func TestPermMeansWhatTheGeneratorWrites(t *testing.T) {
	routes := ParseSource(declSource)

	one := routeNamed(t, routes, "POST", "/api/v1/payslips/:id/email")
	if got := one.Access.Summary(); got != "ADMIN or payroll.edit" {
		t.Errorf("perm(\"payroll.edit\") reads as %q", got)
	}
	if one.Access.NeedsPermission() {
		t.Error("a route with a permission reads as needing one")
	}

	// Variadic: any of them, or an ADMIN.
	many := routeNamed(t, routes, "GET", "/api/v1/payslips")
	if len(many.Access.Permissions) != 2 {
		t.Errorf("perm with two arguments recorded %v", many.Access.Permissions)
	}
	if !strings.Contains(many.Access.Summary(), "payroll.admin") {
		t.Errorf("the second permission is missing from %q", many.Access.Summary())
	}
}

// A route declared open is not the same as a route that forgot, and the whole
// check rests on telling them apart.
func TestADeclaredOpenRouteIsNotAForgottenOne(t *testing.T) {
	routes := ParseSource(declSource)

	declared := routeNamed(t, routes, "GET", "/api/v1/payslips/mine")
	if !declared.Access.DeclaredOpen() {
		t.Error("openToAnySignedIn did not record a declaration")
	}
	if declared.Access.OpenReason != "the handler scopes to the caller" {
		t.Errorf("the reason reads %q", declared.Access.OpenReason)
	}
	if declared.Access.NeedsPermission() {
		t.Error("a declared-open route is reported as needing a permission")
	}
	// And the table says so, because the table is where somebody scans for a
	// route nobody guarded.
	if got := declared.Access.Summary(); got != "signed in (declared)" {
		t.Errorf("the table shows %q, which does not distinguish it from a forgotten route", got)
	}

	forgotten := routeNamed(t, routes, "POST", "/api/v1/payslips/:id/void")
	if forgotten.Access.DeclaredOpen() {
		t.Error("a route with no declaration reads as declared")
	}
	if !forgotten.Access.NeedsPermission() {
		t.Error("the forgotten route is not reported, which is the whole point")
	}
	if got := forgotten.Access.Summary(); got != "signed in" {
		t.Errorf("the forgotten route shows %q", got)
	}
}

// role() is RequireRole without the ADMIN, for the rare route that is about
// who somebody is.
func TestRoleDeclaresARoleAndNoPermission(t *testing.T) {
	routes := ParseSource(declSource)
	r := routeNamed(t, routes, "POST", "/api/v1/payroll/run")

	if len(r.Access.Permissions) != 0 {
		t.Errorf("role() recorded permissions: %v", r.Access.Permissions)
	}
	if len(r.Access.Roles) != 1 || r.Access.Roles[0] != "FINANCE" {
		t.Errorf("roles are %v, want [FINANCE]", r.Access.Roles)
	}
	if r.Access.NeedsPermission() {
		t.Error("a route demanding a role reads as demanding nothing")
	}
}

// A public route is a decision made elsewhere and is not this check's
// business. Reporting them would make the check mostly wrong, and a check
// that is mostly wrong is one people learn to skip.
func TestAPublicRouteIsNotAskedForAPermission(t *testing.T) {
	src := `package routes

func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	public := v1.Group("")
	public.GET("/health", h.Health)
}
`
	routes := ParseSource(src)
	r := routeNamed(t, routes, "GET", "/api/v1/health")
	if r.Access.Authenticated {
		t.Fatal("the fixture's public route reads as authenticated")
	}
	if r.Access.NeedsPermission() {
		t.Error("a public route is reported as needing a permission")
	}
	if got := r.Access.Summary(); got != "public" {
		t.Errorf("summary is %q", got)
	}
}

// A route behind RequireStaff with no permission of its own is exactly the
// case in question: staff is a gate that admits anybody holding any
// permission at all, so the route behind it has to name the one it needs.
func TestAStaffRouteWithNoPermissionIsReported(t *testing.T) {
	src := `package routes

func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	protected := v1.Group("")
	protected.Use(middleware.Auth(cfg))
	staff := protected.Group("")
	staff.Use(middleware.RequireStaff())
	staff.POST("/exports", h.Export)
}
`
	routes := ParseSource(src)
	r := routeNamed(t, routes, "POST", "/api/v1/exports")
	if !r.Access.Staff {
		t.Fatal("the fixture's route is not behind RequireStaff")
	}
	if !r.Access.NeedsPermission() {
		t.Error("a staff route naming no permission was not reported")
	}
}
