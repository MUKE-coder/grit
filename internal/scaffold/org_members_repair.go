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
	if err := repairSourceFile(root, m, path, repairOrgMemberGuardsSource); err != nil {
		return err
	}
	// And the routes, in a different file: handler methods nothing calls are
	// not a feature. A project whose routes could not be reached is told so by
	// the repair itself rather than discovering four missing endpoints.
	// The role endpoints read tenant.ColumnName, which a project that
	// installed the plugin earlier does not have. The constant goes in before
	// the handler that references it, or the project stops compiling.
	tenantGo := filepath.Join(opts.APIRoot(root), "internal", "tenant", "tenant.go")
	if fileExists(tenantGo) {
		if err := repairSourceFile(root, m, tenantGo, repairTenantColumnName); err != nil {
			return err
		}
	}
	routes := filepath.Join(opts.APIRoot(root), "internal", "routes", "routes.go")
	if !fileExists(routes) {
		return nil
	}
	return repairSourceFile(root, m, routes, repairOrgRoleRoutes)
}

// repairTenantColumnName names the organization column once, in the package
// that owns it.
//
// Three things had to agree about the string "org_id": the struct tag on
// tenant.Owned, the WHERE clause the scoping callbacks add, and whatever asks
// whether a table is tenant-owned at all. They were three literals.
func repairTenantColumnName(src string) (string, []string, []string) {
	const anchor = "// Owned marks a model as belonging to an organization."
	if !strings.Contains(src, anchor) {
		return src, nil, nil // not the tenant package this plugin wrote
	}
	if strings.Contains(src, "ColumnName = \"org_id\"") {
		return src, nil, nil
	}
	out := strings.Replace(src, anchor, `// ColumnName is the column that marks a row as belonging to an organization.
//
// One name, in one place, because three things have to agree about it: the
// struct tag on Owned, the WHERE clause the callbacks add, and anything that
// asks whether a table is tenant-owned at all.
const ColumnName = "org_id"

`+anchor, 1)
	out = strings.Replace(out, `Column: clause.Column{Table: table, Name: "org_id"},`,
		`Column: clause.Column{Table: table, Name: ColumnName},`, 1)
	return out, []string{"the organization column is named once, as tenant.ColumnName"}, nil
}

// repairOrgRoleRoutes mounts the per-organization role endpoints in a project
// that already has the membership ones.
func repairOrgRoleRoutes(src string) (string, []string, []string) {
	const anchor = `		protected.DELETE("/organizations/:id/members/:userId", orgHandler.RemoveMember)`
	if !strings.Contains(src, anchor) {
		return src, nil, nil // not a multitenant project, or hand-rewritten routes
	}
	if strings.Contains(src, `orgHandler.CreateRole`) {
		return src, nil, nil
	}
	out := strings.Replace(src, anchor, anchor+`

		// An organization's own roles. Any member may read them; the owner, a
		// platform administrator, or a member holding orgroles.manage may
		// write them. See grit#130.
		protected.GET("/organizations/:id/roles", orgHandler.Roles)
		protected.POST("/organizations/:id/roles", orgHandler.CreateRole)
		protected.PUT("/organizations/:id/roles/:roleId", orgHandler.UpdateRole)
		protected.DELETE("/organizations/:id/roles/:roleId", orgHandler.DeleteRole)`, 1)
	return out, []string{"the per-organization role endpoints are mounted"}, nil
}

func repairOrgMemberGuardsSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func (h *OrganizationHandler) AddMember") {
		return src, nil, nil
	}

	var changes, warnings []string
	out := src

	// The key v3.438.0 invented had three segments, and a permission key has
	// two: authz.Keys() could not produce it, so no role could be authored
	// holding it. Renamed wherever it was written, including in a project the
	// earlier repair already reached.
	// Bare as well as quoted: the guard's own doc comment names the key, and a
	// comment naming a constant that no longer exists is how the next person
	// writes the wrong one.
	if strings.Contains(out, "organizations.members.manage") {
		out = strings.ReplaceAll(out, "organizations.members.manage", "orgmembers.manage")
		changes = append(changes, "the member-management permission is orgmembers.manage, a key the catalogue can produce")
	}

	if strings.Contains(out, "canManageMembers") {
		// Already guarded. The role API may still be missing.
		roled, rc, rw := addOrgRoleAPI(out)
		changes = append(changes, rc...)
		warnings = append(warnings, rw...)
		if len(changes) == 0 {
			return src, nil, warnings
		}
		return roled, changes, warnings
	}

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
	if why := h.roleBeyondCaller(c, orgID, in.RoleID); why != "" {
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
	if !strings.Contains(out, "/internal/authz\"") {
		module := moduleFromImports(out)
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

	roled, rc, rw := addOrgRoleAPI(out)
	out = roled
	changes = append(changes, rc...)
	warnings = append(warnings, rw...)

	if len(changes) == 0 {
		return src, nil, warnings
	}
	return out, changes, warnings
}

// addOrgRoleAPI adds the per-organization role endpoints to a handler that has
// the membership guards but not the roles (grit#130).
func addOrgRoleAPI(src string) (string, []string, []string) {
	if strings.Contains(src, "func (h *OrganizationHandler) CreateRole") {
		return src, nil, nil
	}
	const anchor = "// isMember verifies the caller belongs to the organization"
	if !strings.Contains(src, anchor) {
		return src, nil, []string{
			"isMember is not declared here, so the per-organization role endpoints " +
				"were not added. This project's organization handler has been " +
				"rewritten: add them by hand or reinstall the plugin (grit#130)"}
	}

	out := src
	block := orgRoleSource
	// A project the earlier repair reached declares PermissionMembersManage on
	// its own, so the roles key and the owner-delegation list come with the
	// block instead of with the guards.
	if !strings.Contains(out, "PermissionRolesManage") {
		block = `const PermissionRolesManage = "orgroles.manage"

// orgAdminKeys are the keys an organization's OWNER may hand out without
// holding them, so that delegation works at all.
var orgAdminKeys = []string{PermissionMembersManage, PermissionRolesManage}

` + block
	}
	out = strings.Replace(out, anchor, block+anchor, 1)

	// The accessor the endpoints read the role service through.
	const ctor = `func NewOrganizationHandler(db *gorm.DB) *OrganizationHandler {
	return &OrganizationHandler{DB: db}
}`
	if strings.Contains(out, ctor) && !strings.Contains(out, "func (h *OrganizationHandler) roles()") {
		out = strings.Replace(out, ctor, ctor+`

func (h *OrganizationHandler) roles() *services.RoleService {
	return &services.RoleService{DB: h.DB}
}`, 1)
	}

	if !strings.Contains(out, "/internal/services\"") {
		module := moduleFromImports(out)
		if module == "" {
			return src, nil, []string{
				"the role endpoints need internal/services and the module path could " +
					"not be read from this file's imports, so they were left out " +
					"(grit#130)"}
		}
		added, ok := addImportGroup(out, module+"/internal/services")
		if !ok {
			return src, nil, []string{
				"the role endpoints need internal/services and this file has no " +
					"import block to add it to, so they were left out (grit#130)"}
		}
		out = added
	}

	return out, []string{"an organization can manage its own roles through /organizations/:id/roles"}, nil
}

// moduleFromImports reads the project's module path out of an import this file
// already has, rather than from Options: this runs on a generated file, where
// the path is real and right there.
func moduleFromImports(src string) string {
	if i := strings.Index(src, "/internal/models\""); i > 0 {
		if j := strings.LastIndex(src[:i], "\""); j >= 0 {
			return src[j+1 : i]
		}
	}
	return ""
}

// The same text the plugin writes, so an upgraded project and a freshly
// installed one read identically.
const orgMemberGuardsSource = `// The permissions a per-organization role carries so a customer's own
// administrator can run their organization. Nothing grants either by default:
// until somebody authors such a role, only the organization's owner and a
// platform administrator can do these things.
//
// Two segments, because that is the scheme authz.Keys() builds keys in.
// Before v3.442.0 the first of these had three, under the organizations
// prefix, and so could not be authored through any UI or saved into an
// organization role. The old key is deliberately not written out here: this
// text is written INTO a project, and the rename above is a ReplaceAll, so a
// comment quoting the old name would be rewritten on every later upgrade and
// report a change that is not one.
const (
	PermissionMembersManage = "orgmembers.manage"
	PermissionRolesManage   = "orgroles.manage"
)

// orgAdminKeys are the keys an organization's OWNER may hand out without
// holding them, so that delegation works at all.
var orgAdminKeys = []string{PermissionMembersManage, PermissionRolesManage}

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
// they may. The ceiling is grantsBeyondAuthor, the same one that decides what
// may be written into a role: two separate ceilings disagreed, and between
// them made delegation impossible. See grit#130.
func (h *OrganizationHandler) roleBeyondCaller(c *gin.Context, orgID, roleID string) string {
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
	// Another organization's role is not this organization's to hand out.
	if role.OwnerOrgID != "" && role.OwnerOrgID != orgID {
		return "That role belongs to another organization"
	}
	if strings.EqualFold(role.Name, models.RoleAdmin) {
		return "Only an administrator can grant the " + role.Name + " role"
	}
	if why := h.grantsBeyondAuthor(c, orgID, role.GrantsList()); len(why) > 0 {
		return "The " + role.Name + " role grants more than you may hand out. " +
			strings.Join(why, " ")
	}
	return ""
}

`

// orgRoleSource is the per-organization role API: the endpoints an
// organization manages its own roles through (grit#130).
//
// Sliced out of internal/plugin/multitenant_files.go rather than retyped, and
// TestOrgRoleRepairMatchesPlugin fails if the two drift. An upgraded project
// keeps the canManageMembers that the earlier repair wrote, which does its own
// owner lookup rather than calling isOwner: identical behaviour, and a repair
// that rewrites less is a repair that applies more often.
const orgRoleSource = `// canManageRoles reports whether the caller may author this organization's own
// roles, writing the refusal when they may not.
//
// The same three as member management: a platform administrator, the
// organization's owner, or a member holding orgroles.manage. Authoring a role
// is how permissions enter an organization, so it is guarded at least as
// tightly as handing one out.
func (h *OrganizationHandler) canManageRoles(c *gin.Context, orgID string) bool {
	if !h.isMember(c, orgID) {
		return false
	}
	if authz.IsAdmin(c) || len(authz.Beyond(authz.CallerGrants(c), []string{PermissionRolesManage})) == 0 {
		return true
	}
	if h.isOwner(c, orgID) {
		return true
	}

	c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN",
		"message": "Only the organization's owner can change its roles"}})
	return false
}

// isOwner reports whether the caller owns the organization. It answers nothing
// to the caller: a refusal is the calling check's to write, with its own
// wording.
func (h *OrganizationHandler) isOwner(c *gin.Context, orgID string) bool {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)
	if uid == "" {
		return false
	}

	var org models.Organization
	if err := tenant.Unscoped(h.DB).WithContext(c.Request.Context()).
		Where("id = ?", orgID).First(&org).Error; err != nil {
		return false
	}
	return org.OwnerID != "" && org.OwnerID == uid
}

type orgRoleInput struct {
	Name        string   ` + "`" + `json:"name" binding:"required"` + "`" + `
	Description string   ` + "`" + `json:"description"` + "`" + `
	Grants      []string ` + "`" + `json:"grants"` + "`" + `
}

// Roles lists the roles this organization defines for itself.
//
// Any member may read them: picking a role for a new colleague means seeing
// the names, and the names are not a secret from the people they apply to.
// The platform's own roles are NOT in this list, and an organization cannot
// read or assign them.
func (h *OrganizationHandler) Roles(c *gin.Context) {
	orgID := c.Param("id")
	if !h.isMember(c, orgID) {
		return
	}
	roles, err := h.roles().ForOrg(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not load roles"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// CreateRole authors a role owned by this organization.
//
// What may go in it: the organization's own administration, and permissions on
// models that are tenant-owned. Anything else has to come from a grant the
// author holds. See grantsBeyondAuthor.
func (h *OrganizationHandler) CreateRole(c *gin.Context) {
	orgID := c.Param("id")
	if !h.canManageRoles(c, orgID) {
		return
	}

	var in orgRoleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": "Name is required"}})
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": "Name is required"}})
		return
	}
	if !h.grantsAcceptable(c, orgID, in.Grants) {
		return
	}

	taken, err := h.roles().NameTakenFor(c.Request.Context(), orgID, in.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not check the name"}})
		return
	}
	if taken {
		c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "CONFLICT", "message": "This organization already has a role with that name"}})
		return
	}

	role := models.Role{Name: in.Name, Description: in.Description, OwnerOrgID: orgID}
	if err := role.SetGrants(in.Grants); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not store the permissions"}})
		return
	}
	// The service checks the organization ceiling again on the way in, which
	// is deliberate duplication: this handler is not the only writer, and the
	// one that matters is the one closest to the database.
	if err := h.roles().Create(c.Request.Context(), &role); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": role, "message": "Role created"})
}

// UpdateRole changes one of this organization's own roles.
func (h *OrganizationHandler) UpdateRole(c *gin.Context) {
	orgID := c.Param("id")
	if !h.canManageRoles(c, orgID) {
		return
	}

	// By owner as well as by id, so another organization's role cannot be
	// edited by passing its id. A role that exists and belongs to somebody
	// else answers 404, the same as one that does not exist.
	role, err := h.roles().OwnedBy(c.Request.Context(), orgID, c.Param("roleId"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "Role not found"}})
		return
	}

	var in orgRoleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": "Name is required"}})
		return
	}
	if !h.grantsAcceptable(c, orgID, in.Grants) {
		return
	}

	if name := strings.TrimSpace(in.Name); name != "" && name != role.Name {
		taken, err := h.roles().NameTakenFor(c.Request.Context(), orgID, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not check the name"}})
			return
		}
		if taken {
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "CONFLICT", "message": "This organization already has a role with that name"}})
			return
		}
		role.Name = name
	}
	role.Description = in.Description
	if err := role.SetGrants(in.Grants); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not store the permissions"}})
		return
	}
	if err := h.roles().Save(c.Request.Context(), role); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": role, "message": "Role updated"})
}

// DeleteRole removes one of this organization's own roles.
//
// A role somebody still holds is refused rather than deleted, because deleting
// it changes what those people can do without anyone deciding that.
func (h *OrganizationHandler) DeleteRole(c *gin.Context) {
	orgID := c.Param("id")
	if !h.canManageRoles(c, orgID) {
		return
	}

	role, err := h.roles().OwnedBy(c.Request.Context(), orgID, c.Param("roleId"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "Role not found"}})
		return
	}

	var held int64
	if err := tenant.Unscoped(h.DB).WithContext(c.Request.Context()).
		Model(&models.OrganizationMember{}).
		Where("organization_id = ? AND role_id = ?", orgID, role.ID).
		Count(&held).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not check who holds the role"}})
		return
	}
	if held > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "CONFLICT",
			"message": "This role is held by members. Move them to another role first."}})
		return
	}

	if err := h.roles().Delete(c.Request.Context(), role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not delete the role"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role deleted"})
}

// grantsAcceptable checks the permissions a caller is putting in a role and
// writes the refusal itself. Three questions, in order of how much they cost.
func (h *OrganizationHandler) grantsAcceptable(c *gin.Context, orgID string, grants []string) bool {
	if bad := unknownOrgGrants(grants); len(bad) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR",
			"message": "Unknown permissions: " + strings.Join(bad, ", ")}})
		return false
	}
	// What an organization role may hold at all: not the platform's own
	// administration, and never a wildcard. The same check the service runs.
	probe := models.Role{OwnerOrgID: orgID}
	if err := probe.SetGrants(grants); err == nil {
		if err := services.CheckOrgRoleGrants(&probe); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": err.Error()}})
			return false
		}
	}
	if why := h.grantsBeyondAuthor(c, orgID, grants); len(why) > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN",
			"message": strings.Join(why, " ")}})
		return false
	}
	return true
}

// unknownOrgGrants returns the keys that are not real permissions.
//
// The catalogue, plus the two keys that administer an organization. Those two
// are deliberately not in the catalogue: they mean nothing in a platform role,
// and offering them on the operator's roles page would invite somebody to
// grant "manage every organization's members", which is not what the handlers
// that read them do.
func unknownOrgGrants(grants []string) (bad []string) {
	known := map[string]bool{}
	for _, k := range authz.Keys() {
		known[k] = true
	}
	for _, k := range orgAdminKeys {
		known[k] = true
	}
	for _, g := range grants {
		g = strings.TrimSpace(g)
		// "*" is not unknown, it is all of them. CheckOrgRoleGrants refuses it
		// and says why; calling it an unknown permission here sent somebody
		// looking for a typo.
		if g == "*" || known[g] {
			continue
		}
		bad = append(bad, g)
	}
	return bad
}

// grantsBeyondAuthor returns, for each permission this caller may not hand
// out, a sentence saying why. Empty means they may hand out all of them.
//
// The ceiling is the author's own grants, as on the platform: a role carrying
// a permission its author does not hold is how somebody climbs. A member
// holding orgroles.manage is held to exactly that.
//
// The organization's OWNER is held to a wider one, because they hold no role
// by default: a ceiling of "only what you hold" would have refused them every
// role they tried to write, including the one that delegates their own job,
// and the feature would exist and refuse its only user. What they may hand out
// without holding it is in whyNotDelegable.
//
// Wildcards do not reach here: CheckOrgRoleGrants refuses "*" first, with its
// own wording.
func (h *OrganizationHandler) grantsBeyondAuthor(c *gin.Context, orgID string, grants []string) []string {
	if authz.IsAdmin(c) {
		return nil
	}
	owner := h.isOwner(c, orgID)

	var why []string
	var plainly []string
	for _, g := range authz.CannotGrant(c, grants) {
		if !owner {
			plainly = append(plainly, g)
			continue
		}
		if reason := h.whyNotDelegable(g); reason != "" {
			why = append(why, reason)
		}
	}
	if len(plainly) > 0 {
		why = append(why, "You cannot grant permissions you do not hold: "+
			strings.Join(plainly, ", ")+".")
	}
	return why
}

// whyNotDelegable says why an organization's owner may not put this permission
// in one of their roles, or "" when they may.
//
// Two things qualify. The keys that administer the organization itself, which
// have no model behind them. And a permission whose model is TENANT-OWNED,
// which is the real test and the one the forbidden-prefix list cannot make:
//
// "products.view" in an organization role is safe exactly when products
// carries org_id, because the tenant callback then scopes every read to the
// organization in the request, and a member holding it sees their own rows.
// Without org_id nothing scopes the query and the same grant hands them every
// organization's products. Same permission, same role, opposite blast radius,
// and only the schema knows which.
//
// So the schema is asked. The convention a permission key rests on is that a
// resource's feature key is its table name: "products.view" on the products
// table, the same name /api/products serves.
//
// The reason is spelled out because the first version of this said "you cannot
// grant permissions you do not hold", which is true and is not the reason: the
// fix it implies is to grant the owner that permission, which opens the very
// hole the refusal exists to close.
func (h *OrganizationHandler) whyNotDelegable(key string) string {
	if isOrgAdminKey(key) {
		return ""
	}
	feature := key
	if i := strings.Index(key, "."); i > 0 {
		feature = key[:i]
	}
	if feature == "" {
		return "\"" + key + "\" is not a permission on anything."
	}
	if !h.DB.Migrator().HasTable(feature) {
		return "\"" + key + "\" has no table behind it (" + feature +
			"), so this organization cannot own what it grants."
	}
	if !h.DB.Migrator().HasColumn(feature, tenant.ColumnName) {
		return "\"" + key + "\" cannot go in an organization's role because " +
			feature + " is not tenant-owned: it has no " + tenant.ColumnName +
			" column, so nothing scopes it and a member holding this would see " +
			"every organization's " + feature + ". Embed tenant.Owned in the " +
			feature + " model and migrate, or grant it through a platform role."
	}
	return ""
}

func isOrgAdminKey(key string) bool {
	for _, k := range orgAdminKeys {
		if key == k {
			return true
		}
	}
	return false
}

`
