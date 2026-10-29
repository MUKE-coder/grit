package plugin

import (
	"strings"
	"testing"
)

// Two privilege escalations in the organization membership API, both found
// while working on grit#130 and both reproduced on a scaffolded project.
//
// These read the emitted source rather than running it, because the handler
// is a template this binary does not link. Narrow, and enough to stop the
// exact shapes coming back.

func memberHandlerSource(t *testing.T) string {
	t.Helper()
	src := mtHandler(Context{Module: "example/apps/api"})
	if !strings.Contains(src, "func (h *OrganizationHandler) AddMember") {
		t.Fatal("the organization handler no longer declares AddMember; this test checks nothing")
	}
	return src
}

// Anybody who created an organization became a platform administrator inside
// it. Create assigned the PLATFORM ADMIN role as the creator's membership
// role, and addOrgGrants merges a membership role's grants into the request:
// ADMIN holds "*", so the creator could read every user on the platform.
// Measured on a fresh project: GET /api/users was 403 before creating an
// organization and 200 with the whole table after.
func TestCreatingAnOrganizationGrantsNoPlatformRole(t *testing.T) {
	src := memberHandlerSource(t)

	start := strings.Index(src, "func (h *OrganizationHandler) Create")
	if start < 0 {
		t.Fatal("Create is gone; this test checks nothing")
	}
	body := src[start:]
	if end := strings.Index(body, "\n// Members lists"); end > 0 {
		body = body[:end]
	}

	if strings.Contains(body, "models.RoleAdmin") {
		t.Error("Create still looks up the platform ADMIN role. A membership " +
			"role's grants are merged into every request made with the " +
			"organization active, so this hands the creator \"*\" across the " +
			"whole platform (grit#130, and the shape GHSA-4j45-98gm-mv9v is about)")
	}
	if strings.Contains(body, "RoleID:") {
		t.Error("Create assigns a membership role to the creator. The creator " +
			"is the organization's owner, which is what canManageMembers " +
			"reads; a role on top of that is a grant nobody asked for (grit#130)")
	}
}

// An ordinary member could POST their own user id with the ADMIN role's id
// and be a platform administrator from the next request: 403 before, 200
// after. AddMember checked isMember, which is not a check.
func TestMemberWritesAreGuardedTwice(t *testing.T) {
	src := memberHandlerSource(t)

	for _, fn := range []string{"AddMember", "RemoveMember"} {
		start := strings.Index(src, "func (h *OrganizationHandler) "+fn)
		if start < 0 {
			t.Errorf("%s is gone; this test checks less than it should", fn)
			continue
		}
		body := src[start:]
		if end := strings.Index(body[1:], "\nfunc "); end > 0 {
			body = body[:end+1]
		}
		if !strings.Contains(body, "h.canManageMembers(c,") {
			t.Errorf("%s does not check who may manage members: membership "+
				"alone lets any member change the membership (grit#130)", fn)
		}
		if strings.Contains(body, "if !h.isMember(c,") {
			t.Errorf("%s guards on membership alone, which is the bug (grit#130)", fn)
		}
	}

	// The ceiling is the half that still holds if a project relaxes the
	// first check, so it is checked separately.
	add := src[strings.Index(src, "func (h *OrganizationHandler) AddMember"):]
	if !strings.Contains(add, "h.roleBeyondCaller(c,") {
		t.Error("AddMember writes whatever role_id it is sent, so member " +
			"management can be used to hand out more than the caller holds (grit#130)")
	}

	for _, want := range []string{
		"func (h *OrganizationHandler) canManageMembers(",
		"func (h *OrganizationHandler) roleBeyondCaller(",
		"authz.CannotGrant(c, role.GrantsList())",
		`const PermissionMembersManage = "organizations.members.manage"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the organization handler is missing %s", want)
		}
	}

	// Only a platform administrator hands out ADMIN, whatever else they hold.
	if !strings.Contains(src, "strings.EqualFold(role.Name, models.RoleAdmin)") {
		t.Error("nothing singles out the ADMIN role, so an owner holding a " +
			"wide enough set could hand out the platform administrator role (grit#130)")
	}
}
