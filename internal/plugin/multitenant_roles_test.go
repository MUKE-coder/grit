package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The per-organization role API (grit#130), and the comparator that should
// have existed before any of it was written.

// Every permission key the plugin declares must be one an organization role
// can actually hold. Three separate things decide that, in three files, and
// nothing compared them:
//
//   - authz.Keys() builds "<feature>.<action>" from the catalogue, two
//     segments, so a three-segment key is an unknown permission to
//     validateGrants and cannot be authored through any roles UI.
//   - CheckOrgRoleGrants in services/role.go refuses a list of prefixes that
//     belong to the platform.
//   - the plugin invents the key.
//
// organizations.members.manage failed the first two for four releases. It was
// added in v3.438.0 so a customer's administrator could hold it; the
// "organizations." prefix was added in v3.441.0 to keep platform
// administration out of a customer's role; and the role the feature existed
// for could not be saved. This test is what notices.
func TestOrgPermissionKeysCanBeHeldByAnOrgRole(t *testing.T) {
	src := memberHandlerSource(t)

	keys := orgPermissionKeys(t, src)
	if len(keys) == 0 {
		t.Fatal("no permission keys found in the organization handler; this test checks nothing")
	}

	for _, key := range keys {
		if n := strings.Count(key, ".") + 1; n != 2 {
			t.Errorf("%q has %d segments: a permission key is \"<feature>.<action>\", "+
				"and authz.Keys() cannot produce anything else, so no role can be "+
				"authored holding this (grit#130)", key, n)
		}
		for _, prefix := range forbiddenOrgRolePrefixes(t) {
			if strings.HasPrefix(key, prefix) {
				t.Errorf("%q starts with %q, which CheckOrgRoleGrants refuses in an "+
					"organization's own role. The key exists so a customer's "+
					"administrator can hold it, and it cannot be saved (grit#130)", key, prefix)
			}
		}
	}
}

// orgPermissionKeys reads the Permission* constants out of the emitted handler.
func orgPermissionKeys(t *testing.T, src string) []string {
	t.Helper()
	var keys []string
	for _, line := range strings.Split(src, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "Permission") && !strings.HasPrefix(line, "const Permission") {
			continue
		}
		i := strings.Index(line, "\"")
		j := strings.LastIndex(line, "\"")
		if i < 0 || j <= i {
			continue
		}
		keys = append(keys, line[i+1:j])
	}
	return keys
}

// forbiddenOrgRolePrefixes reads the list out of the role service template, so
// a prefix added there is compared against the plugin's keys the same day.
func forbiddenOrgRolePrefixes(t *testing.T) []string {
	t.Helper()
	path := filepath.Join("..", "scaffold", "templates", "api", "services", "role.go.tmpl")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the role service template is where the forbidden prefixes live: %v", err)
	}
	src := string(body)
	start := strings.Index(src, "var OrgRoleForbiddenPrefixes = []string{")
	if start < 0 {
		t.Fatal("OrgRoleForbiddenPrefixes is gone from the role service; this test checks nothing")
	}
	end := strings.Index(src[start:], "}")
	if end < 0 {
		t.Fatal("OrgRoleForbiddenPrefixes does not close")
	}
	var prefixes []string
	for _, line := range strings.Split(src[start:start+end], "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "\"") {
			continue
		}
		prefixes = append(prefixes, strings.Trim(strings.TrimSuffix(line, ","), "\""))
	}
	if len(prefixes) == 0 {
		t.Fatal("no prefixes parsed out of OrgRoleForbiddenPrefixes")
	}
	return prefixes
}

// Every write to an organization's roles is guarded, and every read of one is
// scoped to the organization in the route. The second is what stops another
// organization's role being edited by passing its id.
func TestOrgRoleEndpointsAreGuardedAndScoped(t *testing.T) {
	src := memberHandlerSource(t)

	for _, fn := range []string{"CreateRole", "UpdateRole", "DeleteRole"} {
		body := handlerBody(t, src, fn)
		if !strings.Contains(body, "h.canManageRoles(c, orgID)") {
			t.Errorf("%s does not check who may manage this organization's roles: "+
				"authoring a role is how permissions enter an organization (grit#130)", fn)
		}
	}

	if body := handlerBody(t, src, "Roles"); !strings.Contains(body, "h.isMember(c, orgID)") {
		t.Error("Roles does not check membership, so one organization can list another's roles")
	}
	if body := handlerBody(t, src, "Roles"); !strings.Contains(body, "ForOrg(") {
		t.Error("Roles does not narrow to the organization, so it lists the platform's roles too")
	}

	for _, fn := range []string{"UpdateRole", "DeleteRole"} {
		body := handlerBody(t, src, fn)
		if !strings.Contains(body, "OwnedBy(c.Request.Context(), orgID,") {
			t.Errorf("%s reads the role by id alone, so another organization's role "+
				"can be changed by passing its id (grit#130)", fn)
		}
	}

	// A role nobody holds is safe to delete; one somebody holds is not.
	if body := handlerBody(t, src, "DeleteRole"); !strings.Contains(body, "held > 0") {
		t.Error("DeleteRole does not check who holds the role, so deleting one " +
			"silently changes what its members can do")
	}
}

// The owner holds no role by default, so a ceiling of "only what you hold"
// would have refused them every role they tried to write, including the one
// that delegates their own job. They are taken to hold the two organization
// keys and nothing else.
func TestOnlyTheTwoOrgKeysAreDelegableByTheOwner(t *testing.T) {
	src := memberHandlerSource(t)
	body := handlerBody(t, src, "grantsBeyondAuthor")

	if !strings.Contains(body, "authz.CannotGrant(c, grants)") {
		t.Error("the ceiling is not the author's own grants, so a member holding " +
			"orgroles.manage can write themselves any permission (grit#130)")
	}
	if !strings.Contains(body, "h.whyNotDelegable(g)") {
		t.Error("the owner's exception is not explained per permission, so a refusal " +
			"names the rule instead of the reason")
	}
	if !strings.Contains(src, "var orgAdminKeys = []string{PermissionMembersManage, PermissionRolesManage}") {
		t.Error("orgAdminKeys is not the two organization keys, so the owner's " +
			"exception covers something nobody chose")
	}
}

func handlerBody(t *testing.T, src, fn string) string {
	t.Helper()
	start := strings.Index(src, "func (h *OrganizationHandler) "+fn+"(")
	if start < 0 {
		t.Fatalf("%s is not declared in the organization handler; this test checks nothing", fn)
	}
	body := src[start:]
	if end := strings.Index(body[1:], "\nfunc "); end > 0 {
		body = body[:end+1]
	}
	return body
}

// The endpoints are mounted. Handler methods nothing calls are not a feature,
// and the routes live in a different file from the handler.
func TestOrgRoleRoutesAreMounted(t *testing.T) {
	var routes string
	for _, inj := range multitenantPlugin().Injections(Context{Module: "example/apps/api"}) {
		if strings.Contains(inj.Code, "orgHandler.") {
			routes += inj.Code
		}
	}
	if routes == "" {
		t.Fatal("the plugin injects no organization routes at all")
	}
	for _, want := range []string{
		`protected.GET("/organizations/:id/roles", orgHandler.Roles)`,
		`protected.POST("/organizations/:id/roles", orgHandler.CreateRole)`,
		`protected.PUT("/organizations/:id/roles/:roleId", orgHandler.UpdateRole)`,
		`protected.DELETE("/organizations/:id/roles/:roleId", orgHandler.DeleteRole)`,
	} {
		if !strings.Contains(routes, want) {
			t.Errorf("not mounted: %s", want)
		}
	}
}

// The upgrade path carries a second copy of this API, because a plugin's files
// are written at install time and `grit upgrade` cannot reinstall a plugin
// over a project's own edits. Two copies of security-shaped code is the thing
// that goes wrong quietly, so they are compared here.
//
// Not byte for byte: the repair holds the same source inside a raw string,
// where a backtick has to be concatenated in, so the struct tags differ by
// construction. Every declaration and every guard is compared instead, which
// is what would actually go missing.
func TestOrgRoleRepairMatchesPlugin(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "scaffold", "org_members_repair.go"))
	if err != nil {
		t.Fatalf("the upgrade repair is the second copy of this API: %v", err)
	}
	repair := string(body)
	if !strings.Contains(repair, "const orgRoleSource = `") {
		t.Fatal("org_members_repair.go no longer carries orgRoleSource, so an " +
			"upgraded project gets the routes and not the handlers (grit#130)")
	}

	src := memberHandlerSource(t)
	for _, decl := range []string{
		"func (h *OrganizationHandler) canManageRoles(",
		"func (h *OrganizationHandler) isOwner(",
		"func (h *OrganizationHandler) Roles(",
		"func (h *OrganizationHandler) CreateRole(",
		"func (h *OrganizationHandler) UpdateRole(",
		"func (h *OrganizationHandler) DeleteRole(",
		"func (h *OrganizationHandler) grantsAcceptable(",
		"func (h *OrganizationHandler) grantsBeyondAuthor(",
		"func (h *OrganizationHandler) whyNotDelegable(",
		"func unknownOrgGrants(",
		"func isOrgAdminKey(",
	} {
		if !strings.Contains(src, decl) {
			t.Errorf("the plugin no longer declares %s; this test checks less than it should", decl)
			continue
		}
		if !strings.Contains(repair, decl) {
			t.Errorf("%s is in the plugin and not in the upgrade repair, so a project "+
				"that upgrades gets a different organization handler from one that "+
				"installs the plugin today (grit#130)", decl)
		}
	}

	// The guards specifically, not just the function names: a copy with the
	// declarations and without the checks is the worst of the two outcomes.
	for _, guard := range []string{
		"h.canManageRoles(c, orgID)",
		"OwnedBy(c.Request.Context(), orgID,",
		"authz.CannotGrant(c, grants)",
		"services.CheckOrgRoleGrants(&probe)",
		"h.DB.Migrator().HasColumn(feature, tenant.ColumnName)",
	} {
		if !strings.Contains(repair, guard) {
			t.Errorf("the upgrade repair's copy is missing the guard %q (grit#130)", guard)
		}
	}
}

// The flow the two ceilings have to permit BETWEEN them.
//
// There were two: one deciding what the owner may write into a role, one
// deciding what they may hand out. Each was defensible alone, and together
// they made the feature impossible: the owner could author the role that
// delegates member management and then could not give it to anybody, because
// handing it out was measured against grants they personally hold and an owner
// holds none. Nothing compared the pair, and only running the whole flow found
// it, which is the third time this shape appeared in one day.
//
// So the pair is what this checks: both call sites must reach the same
// function. A test per ceiling passes while the pair is broken.
func TestBothCeilingsAreTheSameCeiling(t *testing.T) {
	src := memberHandlerSource(t)

	write := handlerBody(t, src, "grantsAcceptable")
	if !strings.Contains(write, "h.grantsBeyondAuthor(c, orgID, grants)") {
		t.Error("authoring a role does not go through grantsBeyondAuthor")
	}

	hand := handlerBody(t, src, "roleBeyondCaller")
	if !strings.Contains(hand, "h.grantsBeyondAuthor(c, orgID, role.GrantsList())") {
		t.Error("handing out a role is measured against a DIFFERENT ceiling from " +
			"authoring one. The owner can then write a role they cannot give to " +
			"anybody, which is the feature refusing its only user (grit#130)")
	}
	if strings.Contains(hand, "authz.CannotGrant(c, role.GrantsList())") {
		t.Error("roleBeyondCaller still holds the caller to their own grants " +
			"directly, which is the ceiling an organization's owner cannot clear")
	}

	// And a role belonging to somebody else is nobody's to hand out: the id was
	// the only thing checked, so one customer could give their member a role
	// another customer owns, and that customer's next edit changed what the
	// first one's people could do.
	if !strings.Contains(hand, `role.OwnerOrgID != "" && role.OwnerOrgID != orgID`) {
		t.Error("roleBeyondCaller does not check who owns the role, so one " +
			"organization can hand out another's (grit#130)")
	}
}
