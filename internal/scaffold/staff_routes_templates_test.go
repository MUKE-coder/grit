package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// routeBlock returns the lines registering routes on one group in the routes
// template, e.g. every "staff.GET(" line.
func routeBlock(src, group string) []string {
	re := regexp.MustCompile(`(?m)^\t\t` + group + `\.(GET|POST|PUT|PATCH|DELETE)\(.*$`)
	return re.FindAllString(src, -1)
}

// A custom role granted users.view was refused by every admin endpoint: the
// whole admin group required the ADMIN role, so every per-route permission
// guard behind it was dead code. The admin UI already showed pages to an
// EDITOR that the API then refused.
func TestStaffRoutesNameTheirPermission(t *testing.T) {
	src := apiRoutesGo()
	staff := routeBlock(src, "staff")
	if len(staff) < 40 {
		t.Fatalf("only %d staff routes; the permission-mapped admin routes did not move", len(staff))
	}
	// The staff gate admits anyone with any grant, so a staff route that names
	// no permission would open to all of them.
	for _, line := range staff {
		if !strings.Contains(line, "middleware.Require") {
			t.Errorf("a staff route names no permission: %s", strings.TrimSpace(line))
		}
	}
	if !strings.Contains(src, "staff.Use(middleware.RequireStaff())") {
		t.Error("the staff group is not gated")
	}
	if !strings.Contains(src, "admin.Use(middleware.RequireRole(\"ADMIN\"))") {
		t.Error("the admin group no longer requires ADMIN, so unguarded routes and plugins opened up")
	}
}

// gin panics at start-up when one method and path is registered twice.
func TestNoRouteIsOnBothGroups(t *testing.T) {
	src := apiRoutesGo()
	key := regexp.MustCompile(`\.(GET|POST|PUT|PATCH|DELETE)\("([^"]+)"`)
	seen := map[string]bool{}
	for _, line := range routeBlock(src, "staff") {
		m := key.FindStringSubmatch(line)
		seen[m[1]+" "+m[2]] = true
	}
	for _, line := range routeBlock(src, "admin") {
		m := key.FindStringSubmatch(line)
		if seen[m[1]+" "+m[2]] {
			t.Errorf("%s %s is registered on both groups", m[1], m[2])
		}
	}
}

func TestStaffMiddlewareAndMount(t *testing.T) {
	mw := apiAuthMiddlewareGo()
	for _, want := range []string{"func RequireStaff() gin.HandlerFunc", "func RequirePermissionFor(param, action string) gin.HandlerFunc"} {
		if !strings.Contains(mw, want) {
			t.Errorf("the auth middleware is missing %s", want)
		}
	}
	reg := apiRoutesRegistryGo()
	if !regexp.MustCompile(`Staff\s+\*gin\.RouterGroup`).MatchString(reg) {
		t.Error("the Mount has no Staff group for generated routes")
	}
	// An older routes.go builds the Mount without Staff; generated routes
	// then fall back to the ADMIN-only group instead of a nil one.
	if !strings.Contains(reg, "m.Staff = m.Admin") {
		t.Error("a Mount without Staff is not given the admin group")
	}
	if !strings.Contains(apiRoutesGo(), "Staff:     staff,") {
		t.Error("routes.go does not hand the staff group to generated routes")
	}
}
