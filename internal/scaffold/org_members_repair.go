package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairOrgMemberGuards closes two privilege escalations in the organization
// membership API that the multitenant plugin wrote into every project that
// installed it.
//
//  1. Create assigned the PLATFORM ADMIN role as the creator's membership
//     role. addOrgGrants merges a membership role's grants into every request
//     made with that organization active, and ADMIN holds "*", so anybody who
//     created an organization could read every user on the platform. On a
//     fresh project GET /api/users was 403 before creating one and 200 with
//     the whole table after.
//
//  2. AddMember and RemoveMember guarded on membership alone, and AddMember
//     wrote whatever role_id it was sent. An ordinary member could POST their
//     own user id with the ADMIN role's id: 403 before, 200 after.
//
// A project that has edited this handler is NOT silently skipped. This is a
// security fix, so a file the repair cannot apply to says so loudly and names
// what to do, rather than leaving somebody to believe an upgrade fixed it.
func repairOrgMemberGuards(root string, opts Options) error {
	path := filepath.Join(opts.APIRoot(root), "internal", "handlers", "organization.go")
	if !fileExists(path) {
		return nil // no multitenant plugin here
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, path, repairOrgMemberGuardsSource)
}

func repairOrgMemberGuardsSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func (h *OrganizationHandler) AddMember") {
		return src, nil, nil
	}
	if strings.Contains(src, "canManageMembers") {
		return src, nil, nil // already repaired
	}

	var changes, warnings []string
	out := src

	// 1. The creator gets no platform role.
	const createRole = `		var admin models.Role
		tx.Where("name = ?", models.RoleAdmin).First(&admin) // grit:role-lookup-ok searched-for, not run
		return tenant.Unscoped(tx).Create(&models.OrganizationMember{
			UserID:         uid,
			OrganizationID: org.ID,
			RoleID:         admin.ID,
		}).Error`
	if strings.Contains(out, createRole) {
		out = strings.Replace(out, createRole, `		// No role. The creator is the organization's owner, which is what
		// canManageMembers reads. This used to assign the PLATFORM ADMIN
		// role, whose grants addOrgGrants merged into every request made with
		// the organization active: anybody who created one could read every
		// user on the platform. See grit#130.
		return tenant.Unscoped(tx).Create(&models.OrganizationMember{
			UserID:         uid,
			OrganizationID: org.ID,
		}).Error`, 1)
		changes = append(changes, "creating an organization no longer grants the platform ADMIN role")
	} else if strings.Contains(out, "models.RoleAdmin") {
		warnings = append(warnings, "Create still references the platform ADMIN role and "+
			"this repair could not rewrite it. Anybody who creates an organization may hold "+
			"\"*\" inside it: check Create by hand (grit#130)")
	}

	// 2. Both member writes go through canManageMembers.
	guarded := 0
	for _, fn := range []string{"AddMember", "RemoveMember"} {
		at := strings.Index(out, "func (h *OrganizationHandler) "+fn)
		if at < 0 {
			continue
		}
		head, tail := out[:at], out[at:]
		const old = "if !h.isMember(c, orgID) {"
		if i := strings.Index(tail, old); i >= 0 && i < 400 {
			tail = tail[:i] + "if !h.canManageMembers(c, orgID) {" + tail[i+len(old):]
			out = head + tail
			guarded++
		}
	}
	if guarded > 0 {
		changes = append(changes, "member writes require permission to manage members, not just membership")
	}

	// 3. The ceiling on what role may be handed out.
	const memberWrite = `	member := models.OrganizationMember{UserID: in.UserID, OrganizationID: orgID, RoleID: in.RoleID}`
	if strings.Contains(out, memberWrite) && !strings.Contains(out, "roleBeyondCaller") {
		out = strings.Replace(out, memberWrite, `	// Nobody hands out more than they hold. See grit#130.
	if why := h.roleBeyondCaller(c, in.RoleID); why != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": why}})
		return
	}

`+memberWrite, 1)
		changes = append(changes, "a role beyond the caller's own grants cannot be handed out")
	}

	// 4. The two guards themselves, and the imports they need.
	const anchor = "// isMember verifies the caller belongs to the organization"
	if !strings.Contains(out, anchor) {
		return src, nil, []string{
			"isMember is not declared here, so the member guards were not added. " +
				"This project's organization handler has been rewritten: add " +
				"canManageMembers and roleBeyondCaller by hand (grit#130)"}
	}
	out = strings.Replace(out, anchor, orgMemberGuardsSource+anchor, 1)

	for _, imp := range []string{"errors", "strings"} {
		if added, ok := addImportGroup(out, imp); ok {
			out = added
		}
	}
	// The module path, taken from an import the file already has rather than
	// from opts: this runs on a generated file, where the path is real.
	if !strings.Contains(out, "/internal/authz\"") {
		module := ""
		if i := strings.Index(out, "/internal/models\""); i > 0 {
			if j := strings.LastIndex(out[:i], "\""); j >= 0 {
				module = out[j+1 : i]
			}
		}
		if module == "" {
			return src, nil, []string{
				"the guards need internal/authz and the module path could not be " +
					"read from this file's imports, so it was left alone: add the " +
					"import by hand (grit#130)"}
		}
		added, ok := addImportGroup(out, module+"/internal/authz")
		if !ok {
			return src, nil, []string{
				"the guards need internal/authz and this file has no import block " +
					"to add it to, so it was left alone (grit#130)"}
		}
		out = added
	}

	if len(changes) == 0 {
		return src, nil, warnings
	}
	return out, changes, warnings
}

// The same text the plugin writes, so an upgraded project and a freshly
// installed one read identically.
const orgMemberGuardsSource = `// PermissionMembersManage is what a per-organization role carries to let a
// customer's own administrator manage their members. Nothing grants it by
// default: until a project defines such a role, only the organization's owner
// and a platform administrator can.
const PermissionMembersManage = "organizations.members.manage"

// canManageMembers reports whether the caller may change who belongs to this
// organization, writing the refusal when they may not.
//
// Membership alone used to be the whole check, and it is not a check: an
// ordinary member could POST their own user id with the ADMIN role's id, and
// addOrgGrants then merged that role's grants into every request they made
// inside the organization. See grit#130.
func (h *OrganizationHandler) canManageMembers(c *gin.Context, orgID string) bool {
	if !h.isMember(c, orgID) {
		return false
	}
	if authz.IsAdmin(c) || len(authz.Beyond(authz.CallerGrants(c), []string{PermissionMembersManage})) == 0 {
		return true
	}

	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)

	var org models.Organization
	if err := tenant.Unscoped(h.DB).Where("id = ?", orgID).First(&org).Error; err == nil {
		if org.OwnerID != "" && org.OwnerID == uid {
			return true
		}
	}

	c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN",
		"message": "Only the organization's owner can change who belongs to it"}})
	return false
}

// roleBeyondCaller says why the caller may not hand out this role, or "" when
// they may. The ceiling is the caller's own grants, so member management
// cannot be used to climb.
func (h *OrganizationHandler) roleBeyondCaller(c *gin.Context, roleID string) string {
	if roleID == "" || authz.IsAdmin(c) {
		return ""
	}

	var role models.Role
	if err := h.DB.WithContext(c.Request.Context()).Where("id = ?", roleID).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ""
		}
		return "The role could not be checked"
	}
	if strings.EqualFold(role.Name, models.RoleAdmin) {
		return "Only an administrator can grant the " + role.Name + " role"
	}
	if beyond := authz.CannotGrant(c, role.GrantsList()); len(beyond) > 0 {
		return "The " + role.Name + " role grants permissions you do not hold: " +
			strings.Join(beyond, ", ")
	}
	return ""
}

`
