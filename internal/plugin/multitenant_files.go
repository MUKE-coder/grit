package plugin

import "strings"

// mtMiddleware emits the middleware that resolves the active organization.
func mtMiddleware(ctx Context) string {
	src := `package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/authz"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/tenant"
)

// The browser is told it may send the header this middleware reads.
//
// Without it the preflight succeeds and the real request is never sent: no
// status, no log line, an axios network error with nothing behind it. It only
// shows up once a user has a second organization, because until then the API
// picks the single one without being told (grit#107).
//
// Registered from here rather than edited into cors.go, so it arrives with this
// file and goes away when grit plugin remove multitenant deletes it.
func init() { AllowRequestHeader("X-Organization-ID") }

// Tenant resolves the caller's active organization and puts it on the request
// context, where the GORM scoping callbacks pick it up.
//
// Resolution order:
//  1. the X-Organization-ID header (explicit; what SPA clients send)
//  2. the user's only organization, if they belong to exactly one
//
// No subdomain resolution: it forces DNS and TLS decisions on every deployment
// and makes local development awkward.
//
// Membership is ALWAYS verified against the database. Trusting the header alone
// would let any authenticated user read another tenant's data by editing one
// request header — the whole isolation model rests on this check.
func Tenant(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.Next() // unauthenticated; nothing to scope
			return
		}
		uid, _ := userID.(string)

		requested := c.GetHeader("X-Organization-ID")

		var memberships []models.OrganizationMember
		if err := db.Where("user_id = ?", uid).Find(&memberships).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not resolve organization"},
			})
			c.Abort()
			return
		}

		active := ""
		switch {
		case requested != "":
			for _, m := range memberships {
				if m.OrganizationID == requested {
					active = requested
					break
				}
			}
			if active == "" {
				// Asked for an organization they do not belong to. 404 rather
				// than 403: a 403 confirms it exists and lets an attacker
				// enumerate tenant ids.
				c.JSON(http.StatusNotFound, gin.H{
					"error": gin.H{"code": "NOT_FOUND", "message": "Organization not found"},
				})
				c.Abort()
				return
			}
		case len(memberships) == 1:
			active = memberships[0].OrganizationID
		}

		if active != "" {
			c.Set("org_id", active)
			roleID := roleFor(memberships, active)
			c.Set("org_role_id", roleID)
			addOrgGrants(c, db, roleID)
			c.Request = c.Request.WithContext(tenant.WithOrg(c.Request.Context(), active))
		}
		c.Next()
	}
}

// addOrgGrants gives the caller the permissions their role in THIS organization
// carries, for this request only.
//
// A per-organization role adds to the platform ones rather than replacing them:
// an organization cannot take away a permission the platform granted, and an
// existing project that has no membership roles behaves exactly as before. Until
// this existed the membership role was stored, put on the context, and read by
// nothing, so being an administrator of one organization meant being one
// everywhere.
func addOrgGrants(c *gin.Context, db *gorm.DB, roleID string) {
	if roleID == "" {
		return
	}
	extra, err := authz.GrantsForRole(db, roleID)
	if err != nil || len(extra) == 0 {
		// Not fatal: the request continues with the platform grants, which fails
		// closed rather than 500ing every route because one role row is bad.
		return
	}
	existing, _ := c.Get("user_grants")
	list, _ := existing.([]string)

	seen := make(map[string]bool, len(list)+len(extra))
	merged := make([]string, 0, len(list)+len(extra))
	for _, grant := range append(append([]string{}, list...), extra...) {
		if seen[grant] {
			continue
		}
		seen[grant] = true
		merged = append(merged, grant)
	}
	c.Set("user_grants", merged)
}

func roleFor(memberships []models.OrganizationMember, orgID string) string {
	for _, m := range memberships {
		if m.OrganizationID == orgID {
			return m.RoleID
		}
	}
	return ""
}

// RequireOrganization rejects a request with no active organization. Put it on
// routes that operate on tenant data.
func RequireOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := c.Get("org_id"); !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "NO_ORGANIZATION",
					"message": "Select an organization first (send the X-Organization-ID header)",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
`
	return strings.ReplaceAll(strings.ReplaceAll(src, "~", "`"), "{{MODULE}}", ctx.Module)
}

// mtHandler emits the organization + membership API.
func mtHandler(ctx Context) string {
	src := `package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/authz"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
	"{{MODULE}}/internal/tenant"
)

// OrganizationHandler serves organization + membership endpoints.
type OrganizationHandler struct {
	DB *gorm.DB
}

func NewOrganizationHandler(db *gorm.DB) *OrganizationHandler {
	return &OrganizationHandler{DB: db}
}

func (h *OrganizationHandler) roles() *services.RoleService {
	return &services.RoleService{DB: h.DB}
}

// MyOrganizations lists the organizations the caller belongs to.
//
// Reads memberships first and looks organizations up by id, rather than listing
// organizations and filtering — so a bug here can only ever return too FEW
// rows, never another tenant's.
func (h *OrganizationHandler) MyOrganizations(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)

	var memberships []models.OrganizationMember
	if err := h.DB.Where("user_id = ?", uid).Find(&memberships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not load organizations"}})
		return
	}
	if len(memberships) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []models.Organization{}})
		return
	}

	ids := make([]string, 0, len(memberships))
	for _, m := range memberships {
		ids = append(ids, m.OrganizationID)
	}

	// Organizations are not themselves tenant-scoped rows, so this lookup is
	// exempt from the scoping callbacks.
	var orgs []models.Organization
	if err := tenant.Unscoped(h.DB).Where("id IN ?", ids).Find(&orgs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not load organizations"}})
		return
	}

	active, _ := c.Get("org_id")
	c.JSON(http.StatusOK, gin.H{"data": orgs, "active": active})
}

type createOrgInput struct {
	Name        string ~json:"name" binding:"required"~
	Description string ~json:"description"~
}

// Create makes an organization and adds the caller as its owner.
func (h *OrganizationHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)

	var in createOrgInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": "Name is required"}})
		return
	}

	org := models.Organization{
		Name:        strings.TrimSpace(in.Name),
		Slug:        slugify(in.Name),
		Description: in.Description,
		OwnerID:     uid,
		Active:      true,
	}

	// Both writes in one transaction: an organization with no members is
	// unreachable, because every read goes through membership.
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tenant.Unscoped(tx).Create(&org).Error; err != nil {
			return err
		}
		// No role. The creator is the organization's owner, which is what
		// canManageMembers reads, and that is the whole of what they need.
		//
		// This used to assign the PLATFORM ADMIN role, whose grants
		// addOrgGrants then merged into every request made with the
		// organization active. ADMIN holds "*", so anybody who created an
		// organization could read every user on the platform rather than
		// their own members: 403 before, 200 and the whole table after. That
		// is what GHSA-4j45-98gm-mv9v describes about a membership role
		// holding a platform grant, and the plugin was doing it by default.
		// See grit#130.
		return tenant.Unscoped(tx).Create(&models.OrganizationMember{
			UserID:         uid,
			OrganizationID: org.ID,
		}).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not create organization"}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": org, "message": "Organization created"})
}

// Members lists an organization's members. Callers must belong to it.
func (h *OrganizationHandler) Members(c *gin.Context) {
	if !h.isMember(c, c.Param("id")) {
		return
	}
	var members []models.OrganizationMember
	if err := tenant.Unscoped(h.DB).Where("organization_id = ?", c.Param("id")).Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not load members"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": members})
}

type addMemberInput struct {
	UserID string ~json:"user_id" binding:"required"~
	RoleID string ~json:"role_id"~
}

// AddMember adds a user to an organization with a role.
func (h *OrganizationHandler) AddMember(c *gin.Context) {
	orgID := c.Param("id")
	if !h.canManageMembers(c, orgID) {
		return
	}

	var in addMemberInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION_ERROR", "message": "user_id is required"}})
		return
	}

	// Nobody hands out more than they hold. This is the guard that closes the
	// escalation even where the one above is relaxed, so it is checked on the
	// role itself rather than on who is asking.
	if why := h.roleBeyondCaller(c, orgID, in.RoleID); why != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": why}})
		return
	}

	member := models.OrganizationMember{UserID: in.UserID, OrganizationID: orgID, RoleID: in.RoleID}
	// Upsert: re-adding an existing member updates their role rather than
	// failing on the composite primary key.
	if err := tenant.Unscoped(h.DB).Save(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not add member"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": member, "message": "Member added"})
}

// RemoveMember removes a user from an organization.
func (h *OrganizationHandler) RemoveMember(c *gin.Context) {
	orgID := c.Param("id")
	if !h.canManageMembers(c, orgID) {
		return
	}

	// The owner keeps access, or an organization can be orphaned with nobody
	// able to administer it.
	var org models.Organization
	if err := tenant.Unscoped(h.DB).Where("id = ?", orgID).First(&org).Error; err == nil {
		if org.OwnerID == c.Param("userId") {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "The owner cannot be removed"}})
			return
		}
	}

	if err := tenant.Unscoped(h.DB).
		Where("organization_id = ? AND user_id = ?", orgID, c.Param("userId")).
		Delete(&models.OrganizationMember{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "Could not remove member"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}

// The permissions a per-organization role carries so a customer's own
// administrator can run their organization. Nothing grants either by default:
// until somebody authors such a role, only the organization's owner and a
// platform administrator can do these things.
//
// Two segments, because that is the scheme: a key is "<feature>.<action>" and
// authz.Keys() builds every valid one from the catalogue. Before v3.442.0 the
// first of these had three segments, under the organizations prefix, and so
// could not be authored through the roles UI (an unknown permission), could
// not be saved into an organization role (that prefix is forbidden there), and
// was never in the catalogue at all. Nothing can hold the old name.
const (
	PermissionMembersManage = "orgmembers.manage"
	PermissionRolesManage   = "orgroles.manage"
)

// orgAdminKeys are the keys an organization's OWNER may hand out without
// holding them, so that delegation works at all. See grantsBeyondAuthor.
var orgAdminKeys = []string{PermissionMembersManage, PermissionRolesManage}

// canManageMembers reports whether the caller may change who belongs to this
// organization, writing the refusal when they may not.
//
// Membership alone used to be the whole check, and it is not a check: an
// ordinary member could POST their own user id with the ADMIN role's id, and
// addOrgGrants then merged that role's grants into every request they made
// inside the organization. Reproduced on a fresh project: the same GET
// /api/users answered 403 before the call and 200 after it. See grit#130.
//
// Who may: a platform administrator, the organization's owner, or a member
// whose grants carry orgmembers.manage. That last is what a
// per-organization role gives a customer's own administrator, and until a
// project defines such a role only the owner can do it, which is the safe
// default rather than the convenient one.
func (h *OrganizationHandler) canManageMembers(c *gin.Context, orgID string) bool {
	if !h.isMember(c, orgID) {
		return false
	}
	if authz.IsAdmin(c) || len(authz.Beyond(authz.CallerGrants(c), []string{PermissionMembersManage})) == 0 {
		return true
	}

	if h.isOwner(c, orgID) {
		return true
	}

	c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN",
		"message": "Only the organization's owner can change who belongs to it"}})
	return false
}

// canManageRoles reports whether the caller may author this organization's own
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

// roleBeyondCaller says why the caller may not hand out this role, or "" when
// they may.
//
// The ceiling is grantsBeyondAuthor, the SAME one that decides what may be
// written into a role. It used to be authz.CannotGrant on its own, and the two
// disagreed: the organization's owner could author the role that delegates
// member management and then could not hand it to anybody, because handing it
// out was measured against grants they personally hold and an owner holds
// none. Two defensible rules, in two places, that between them made the
// feature impossible. See grit#130.
//
// A role id that matches nothing grants nothing and is allowed through, which
// is the same answer as sending no role at all.
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

	// Another organization's role is not this organization's to hand out. The
	// id was the only thing checked before, so one customer could give their
	// member a role another customer owns, and that customer's next edit to it
	// changed what the first one's people could do. That is the complaint
	// grit#130 opens with, arriving by a different door.
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

type orgRoleInput struct {
	Name        string   ~json:"name" binding:"required"~
	Description string   ~json:"description"~
	Grants      []string ~json:"grants"~
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

// isMember verifies the caller belongs to the organization, writing a 404 and
// returning false when they do not.
//
// 404 rather than 403: a 403 confirms the organization exists, which lets an
// attacker enumerate tenants.
func (h *OrganizationHandler) isMember(c *gin.Context, orgID string) bool {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)

	var count int64
	tenant.Unscoped(h.DB).Model(&models.OrganizationMember{}).
		Where("user_id = ? AND organization_id = ?", uid, orgID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "Organization not found"}})
		return false
	}
	return true
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '_':
			b.WriteRune('-')
		}
	}
	return strings.Trim(b.String(), "-")
}
`
	return strings.ReplaceAll(strings.ReplaceAll(src, "~", "`"), "{{MODULE}}", ctx.Module)
}

// mtTenantTests emits isolation tests into the generated project.
//
// These matter more than most generated tests: they assert that a query for one
// organization cannot see another's rows, and that an unscoped context fails
// CLOSED. A regression here is a data breach, not a bug.
func mtTenantTests(ctx Context) string {
	src := `package tenant_test

import (
	"context"
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/tenant"
)

// Widget is a tenant-owned model for these tests.
type Widget struct {
	ID   string ~gorm:"primarykey;size:36"~
	Name string
	tenant.Owned
}

func newDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&Widget{}); err != nil {
		t.Fatal(err)
	}
	if err := tenant.RegisterScoping(db); err != nil {
		t.Fatal(err)
	}
	return db
}

func seed(t *testing.T, db *gorm.DB, org, id, name string) {
	t.Helper()
	ctx := tenant.WithOrg(context.Background(), org)
	if err := db.WithContext(ctx).Create(&Widget{ID: id, Name: name}).Error; err != nil {
		t.Fatal(err)
	}
}

// A batch insert is the ordinary GORM way to add several rows, and it reaches
// the stamp callback with a slice rather than a struct. Before the callback
// handled both, this panicked inside reflect and took the process down.
func TestStamp_BatchCreate(t *testing.T) {
	db := newDB(t)
	ctx := tenant.WithOrg(context.Background(), "org-a")

	rows := []Widget{{ID: "w1", Name: "one"}, {ID: "w2", Name: "two"}}
	if err := db.WithContext(ctx).Create(&rows).Error; err != nil {
		t.Fatalf("batch create: %v", err)
	}
	for _, r := range rows {
		if r.OrgID != "org-a" {
			t.Errorf("row %s has org %q, want org-a", r.ID, r.OrgID)
		}
	}

	// And the rows that landed are scoped, not just the structs in memory.
	var found []Widget
	if err := db.WithContext(ctx).Find(&found).Error; err != nil {
		t.Fatal(err)
	}
	if len(found) != 2 {
		t.Fatalf("read back %d rows, want 2", len(found))
	}
}

// Slices of pointers are just as common, and a pointer has no fields to set
// either: the element has to be dereferenced first.
func TestStamp_BatchCreatePointers(t *testing.T) {
	db := newDB(t)
	ctx := tenant.WithOrg(context.Background(), "org-b")

	rows := []*Widget{{ID: "p1", Name: "one"}, {ID: "p2", Name: "two"}}
	if err := db.WithContext(ctx).Create(&rows).Error; err != nil {
		t.Fatalf("batch create: %v", err)
	}
	for _, r := range rows {
		if r.OrgID != "org-b" {
			t.Errorf("row %s has org %q, want org-b", r.ID, r.OrgID)
		}
	}
}

// An explicit organization still wins over the one on the context, the same
// as it does for a single row.
func TestStamp_BatchCreateKeepsAnExplicitOrg(t *testing.T) {
	db := newDB(t)
	ctx := tenant.WithOrg(context.Background(), "org-a")

	rows := []Widget{
		{ID: "m1", Name: "default"},
		{ID: "m2", Name: "explicit", Owned: tenant.Owned{OrgID: "org-z"}},
	}
	if err := db.WithContext(ctx).Create(&rows).Error; err != nil {
		t.Fatalf("batch create: %v", err)
	}
	if rows[0].OrgID != "org-a" {
		t.Errorf("m1 has org %q, want org-a", rows[0].OrgID)
	}
	if rows[1].OrgID != "org-z" {
		t.Errorf("m2 has org %q; an explicit value must not be overwritten", rows[1].OrgID)
	}
}

// The core guarantee: one organization cannot read another's rows.
func TestScoping_IsolatesOrganizations(t *testing.T) {
	db := newDB(t)
	seed(t, db, "org-a", "w1", "alpha")
	seed(t, db, "org-b", "w2", "beta")

	var got []Widget
	ctx := tenant.WithOrg(context.Background(), "org-a")
	if err := db.WithContext(ctx).Find(&got).Error; err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 {
		t.Fatalf("org-a saw %d widgets, want 1 — CROSS-TENANT LEAK", len(got))
	}
	if got[0].Name != "alpha" {
		t.Errorf("org-a saw %q", got[0].Name)
	}
}

// Create stamps the active organization, so a row cannot land in the wrong one.
func TestScoping_StampsOrgOnCreate(t *testing.T) {
	db := newDB(t)
	seed(t, db, "org-a", "w1", "alpha")

	var w Widget
	if err := tenant.Unscoped(db).Where("id = ?", "w1").First(&w).Error; err != nil {
		t.Fatal(err)
	}
	if w.OrgID != "org-a" {
		t.Errorf("OrgID = %q, want org-a", w.OrgID)
	}
}

// No active organization must FAIL, not silently return everything. This is the
// property that makes the whole design safe.
func TestScoping_FailsClosedWithoutOrg(t *testing.T) {
	db := newDB(t)
	seed(t, db, "org-a", "w1", "alpha")

	var got []Widget
	err := db.WithContext(context.Background()).Find(&got).Error
	if err == nil {
		t.Fatalf("query without an organization returned %d rows; it must error", len(got))
	}
}

// Updates and deletes are scoped too — otherwise one tenant could modify
// another's data even without reading it.
func TestScoping_UpdateAndDeleteAreScoped(t *testing.T) {
	db := newDB(t)
	seed(t, db, "org-a", "w1", "alpha")
	seed(t, db, "org-b", "w2", "beta")

	ctxA := tenant.WithOrg(context.Background(), "org-a")
	if err := db.WithContext(ctxA).Model(&Widget{}).Where("1 = 1").Update("name", "hacked").Error; err != nil {
		t.Fatal(err)
	}
	var b Widget
	tenant.Unscoped(db).Where("id = ?", "w2").First(&b)
	if b.Name != "beta" {
		t.Errorf("org-a modified org-b's row: %q", b.Name)
	}

	if err := db.WithContext(ctxA).Where("1 = 1").Delete(&Widget{}).Error; err != nil {
		t.Fatal(err)
	}
	var count int64
	tenant.Unscoped(db).Model(&Widget{}).Where("id = ?", "w2").Count(&count)
	if count != 1 {
		t.Error("org-a deleted org-b's row")
	}
}

// Unscoped is the deliberate escape hatch for cross-tenant work.
func TestUnscoped_SeesEverything(t *testing.T) {
	db := newDB(t)
	seed(t, db, "org-a", "w1", "alpha")
	seed(t, db, "org-b", "w2", "beta")

	var all []Widget
	if err := tenant.Unscoped(db).Find(&all).Error; err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Errorf("Unscoped saw %d widgets, want 2", len(all))
	}
}

// A handle kept in a variable has to behave like a fresh one every time.
//
// Unscoped used to return GORM's chainable instance, which shares one
// Statement across every call made on it, so a kept handle carried the
// previous call's WHERE clause and its error into the next (grit#103):
//
//	tx = tenant.Unscoped(tx)
//	tx.Where("slug = ?", slug).Take(&existing)  // ErrRecordNotFound
//	return tx.Create(&org).Error                // also "record not found"
//
// The one-shot form in the doc comment was always fine, and keeping the
// handle is the natural next step, so the failure found people who had done
// nothing wrong and pointed nowhere near scoping.
func TestUnscoped_HandleDoesNotLeakBetweenStatements(t *testing.T) {
	db := newDB(t)

	un := tenant.Unscoped(db)

	// A miss, which is expected here and must not be contagious.
	var found Widget
	err := un.Where("name = ?", "nothing-by-this-name").Take(&found).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected the lookup to miss, got %v", err)
	}

	if err := un.Create(&Widget{ID: "w1", Name: "after a miss"}).Error; err != nil {
		t.Fatalf("the create inherited the previous call's error: %v", err)
	}
	if err := un.Create(&Widget{ID: "w2", Name: "second"}).Error; err != nil {
		t.Fatal(err)
	}

	// And two reads must not combine their conditions into one that matches
	// nothing.
	var first, second []Widget
	if err := un.Where("name = ?", "after a miss").Find(&first).Error; err != nil {
		t.Fatal(err)
	}
	if err := un.Where("name = ?", "second").Find(&second).Error; err != nil {
		t.Fatal(err)
	}
	if len(first) != 1 || len(second) != 1 {
		t.Errorf("two reads through one handle combined their conditions: got %d and %d, want 1 and 1",
			len(first), len(second))
	}
}

// Models without an OrgID are untouched — users, roles and settings are global.
func TestScoping_IgnoresNonTenantModels(t *testing.T) {
	type Global struct {
		ID   string ~gorm:"primarykey;size:36"~
		Name string
	}
	db := newDB(t)
	if err := db.AutoMigrate(&Global{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&Global{ID: "g1", Name: "x"}).Error; err != nil {
		t.Fatalf("a non-tenant model must not require an organization: %v", err)
	}
	var got []Global
	if err := db.Find(&got).Error; err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Errorf("non-tenant query returned %d rows, want 1", len(got))
	}
}

// ── Scan, Rows and Row ───────────────────────────────────────────────────
//
// GORM's Scan does NOT go through the Query callback. It calls Rows(), which
// executes the ROW processor:
//
//	func (db *DB) Scan(dest any) *DB      { ... tx.Rows() ... }
//	func (db *DB) Rows() (*sql.Rows, err) { tx.callbacks.Row().Execute(tx) }
//
// Scoping was once registered on Query, Update, Delete and Create but not
// Row, so Find and First behaved while Scan read every organization's rows.
// The ordinary reads behaving is what hid it: the leak was only on the
// aggregate paths, and the generated insights are built entirely from those.
//
// Reported by a Grit user who found it in their own multi-tenant project.

func seedBothOrgs(t *testing.T, db *gorm.DB) {
	t.Helper()
	seed(t, db, "org-a", "a1", "alpha")
	seed(t, db, "org-b", "b1", "bravo")
	seed(t, db, "org-b", "b2", "charlie")
}

func TestScan_IsScoped(t *testing.T) {
	db := newDB(t)
	seedBothOrgs(t, db)

	ctx := tenant.WithOrg(context.Background(), "org-a")
	var got []Widget
	if err := db.WithContext(ctx).Model(&Widget{}).Scan(&got).Error; err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("Scan returned %d rows, want 1.\n\n"+
			"This organization can read another's data. Check that "+
			"RegisterScoping still registers the Row callback: Scan goes "+
			"through Row, not Query.", len(got))
	}
}

func TestRows_IsScoped(t *testing.T) {
	db := newDB(t)
	seedBothOrgs(t, db)

	ctx := tenant.WithOrg(context.Background(), "org-a")
	rows, err := db.WithContext(ctx).Model(&Widget{}).Rows()
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	seen := 0
	for rows.Next() {
		seen++
	}
	if seen != 1 {
		t.Fatalf("Rows() yielded %d rows, want 1", seen)
	}
}

// The shape the generated insights use: a grouped COUNT(*) through Scan.
// paginate/insights.go does this for the time series and for every breakdown
// column, so an unscoped Row callback means an organization's own admin
// screen counts everybody's rows.
func TestAggregateScan_IsScoped(t *testing.T) {
	db := newDB(t)
	seedBothOrgs(t, db)

	type slice struct {
		Value string
		Count int64
	}

	ctx := tenant.WithOrg(context.Background(), "org-a")
	var out []slice
	err := db.WithContext(ctx).Model(&Widget{}).
		Select("name AS value, COUNT(*) AS count").
		Group("value").
		Scan(&out).Error
	if err != nil {
		t.Fatal(err)
	}

	var total int64
	for _, s := range out {
		total += s.Count
	}
	if total != 1 {
		t.Fatalf("the breakdown counted %d rows, want 1. This organization's "+
			"insights are counting every organization's rows.", total)
	}
}

// Raw SQL has no parsed schema, so it is not scoped and is not meant to be.
// Pinned so that scoping Row is never mistaken for making raw queries safe.
func TestRaw_IsNotScoped(t *testing.T) {
	db := newDB(t)
	seedBothOrgs(t, db)

	ctx := tenant.WithOrg(context.Background(), "org-a")
	var n int64
	if err := db.WithContext(ctx).Raw("SELECT COUNT(*) FROM widgets").Scan(&n).Error; err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("raw count is %d, want 3. Raw SQL is deliberately unscoped: "+
			"write the org filter yourself, or go through the model.", n)
	}
}
`
	return strings.ReplaceAll(strings.ReplaceAll(src, "~", "`"), "{{MODULE}}", ctx.Module)
}
