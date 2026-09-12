package plugin

import (
	"strings"
	"testing"
)

// The multitenant plugin has to resolve the active organization on every
// authenticated route group.
//
// It used to patch the protected group alone. The staff group carries every DELETE
// and every bulk route and the admin group every admin-only endpoint, so on those
// the request never resolved an organization: the scoping callback failed closed
// and deleting a tenant-owned row answered 500. Found by building the reviewer's
// Project 1 and trying to delete a deal.
func TestMultitenantMountsTheResolverOnEveryAuthenticatedGroup(t *testing.T) {
	ctx := Context{Module: "example.com/app", APIRoot: "apps/api"}
	injections := multitenantPlugin().Injections(ctx)

	want := map[string]string{
		"// grit:middleware:protected": "protected.Use(middleware.Tenant(db))",
		"// grit:middleware:staff":     "staff.Use(middleware.Tenant(db))",
		"// grit:middleware:admin":     "admin.Use(middleware.Tenant(db))",
	}
	found := map[string]bool{}
	for _, injection := range injections {
		if code, ok := want[injection.Marker]; ok && strings.Contains(injection.Code, code) {
			found[injection.Marker] = true
		}
	}
	for marker, code := range want {
		if !found[marker] {
			t.Errorf("nothing patches %s with %s, so that route group never resolves an organization",
				marker, code)
		}
	}
}

// The error it raises says what it is on the wire.
//
// "No active organization" used to fall through respond.WriteError to an opaque
// 500 reading "Failed to fetch deals", for a request whose only problem was not
// naming an organization.
func TestNoOrganizationErrorCarriesItsCode(t *testing.T) {
	src := mtTenantPackage(Context{Module: "example.com/app", APIRoot: "apps/api"})
	for _, want := range []string{
		"func (noOrganizationError) ErrorCode() respond.Code { return respond.CodeNoOrganization }",
		"var ErrNoOrganization error = noOrganizationError{}",
		`"example.com/app/internal/respond"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("internal/tenant/tenant.go is missing %q", want)
		}
	}
}

// And a role held through a membership grants inside that organization.
//
// The middleware put org_role_id on the request context and nothing read it, so
// somebody made an administrator of one organization held the same permissions in
// every organization they belonged to.
func TestTenantMiddlewareAppliesTheOrganizationRole(t *testing.T) {
	src := mtMiddleware(Context{Module: "example.com/app", APIRoot: "apps/api"})
	for _, want := range []string{
		"addOrgGrants(c, db, roleID)",
		"func addOrgGrants(c *gin.Context, db *gorm.DB, roleID string)",
		"authz.GrantsForRole(db, roleID)",
		`"example.com/app/internal/authz"`,
		// Merged, not replaced: an organization cannot take away a platform grant.
		`c.Set("user_grants", merged)`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("internal/middleware/tenant.go is missing %q", want)
		}
	}
}
