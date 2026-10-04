package handlers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"library/apps/api/internal/authz"
	"library/apps/api/internal/models"
	"library/apps/api/internal/respond"
	"library/apps/api/internal/services"
)

// RoleHandler serves the roles + permissions API.
type RoleHandler struct {
	DB *gorm.DB
}

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{DB: db}
}

// roles is the service every method below reads and writes through. Built per
// request because it is one field and a struct literal, and because building it
// here means the wiring in routes.go did not have to change.
//
// What a role is, how its holders are counted and which writes share a
// transaction are decisions about data: see internal/services/role.go. What
// follows decides who may ask.
func (h *RoleHandler) roles() *services.RoleService {
	return &services.RoleService{DB: h.DB}
}

// users is the same thing for the one account this handler reads: the user whose
// roles are being replaced.
func (h *RoleHandler) users() *services.UserService {
	return &services.UserService{DB: h.DB}
}

// roleDTO is a Role plus the derived fields the UI needs.
type roleDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Grants      []string `json:"grants"`
	// Expanded is grants with wildcards resolved against the catalog. The UI
	// renders checkboxes from this and never reimplements wildcard matching —
	// a duplicated matcher is how the reference implementation's Go and
	// TypeScript rules drifted apart.
	Expanded  []string `json:"expanded"`
	IsSystem  bool     `json:"is_system"`
	UserCount int64    `json:"user_count"`
}

func (h *RoleHandler) toDTO(r models.Role, userCount int64) roleDTO {
	grants := r.GrantsList()
	return roleDTO{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Grants:      grants,
		Expanded:    authz.Expand(grants),
		IsSystem:    r.IsSystem,
		UserCount:   userCount,
	}
}

// Catalog returns the permission tree so the UI can render the editor.
func (h *RoleHandler) Catalog(c *gin.Context) {
	respond.OK(c, gin.H{
		"modules": authz.Catalog(),
		"keys":    authz.Keys(),
	})
}

// List returns every role with how many users hold it.
func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.roles().All(c.Request.Context())
	if err != nil {
		respond.Internal(c, err)
		return
	}
	// The count the list shows beside each role. Its error used to go nowhere,
	// so a failure drew every role as held by nobody.
	counts, err := h.roles().HolderCounts(c.Request.Context())
	if err != nil {
		respond.Internal(c, err)
		return
	}

	out := make([]roleDTO, 0, len(roles))
	for _, r := range roles {
		out = append(out, h.toDTO(r, counts[r.ID]))
	}
	respond.OK(c, out)
}

func (h *RoleHandler) Get(c *gin.Context) {
	role, err := h.roles().ByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respond.NotFound(c, "Role not found")
		return
	}
	n, err := h.roles().Holders(c.Request.Context(), role.ID)
	if err != nil {
		respond.Internal(c, err)
		return
	}
	respond.OK(c, h.toDTO(*role, n))
}

type RoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Grants      []string `json:"grants"`
}

// validateGrants rejects keys the catalog doesn't know.
//
// Without this a typo is stored happily and simply never authorises anything,
// which presents as "I ticked the box and it doesn't work".
func validateGrants(grants []string) (bad []string) {
	known := map[string]bool{}
	for _, k := range authz.Keys() {
		known[k] = true
	}
	for _, g := range grants {
		if g == "*" || known[g] {
			continue
		}
		// Wildcards are fine as long as they match something real.
		if strings.HasSuffix(g, ".*") && len(authz.Expand([]string{g})) > 0 {
			continue
		}
		bad = append(bad, g)
	}
	return bad
}

func (h *RoleHandler) Create(c *gin.Context) {
	var in RoleRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, "Name is required")
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		respond.BadRequest(c, "Name is required")
		return
	}
	if bad := validateGrants(in.Grants); len(bad) > 0 {
		respond.BadRequest(c, "Unknown permissions: "+strings.Join(bad, ", "))
		return
	}
	// A role hands out at most what its author holds, or a roles.create or
	// roles.edit holder could write "*" into a role and take it.
	if beyond := authz.CannotGrant(c, in.Grants); len(beyond) > 0 {
		respond.Forbidden(c, "You cannot grant permissions you do not hold: "+strings.Join(beyond, ", "))
		return
	}

	taken, err := h.roles().NameTaken(c.Request.Context(), in.Name)
	if err != nil {
		respond.Internal(c, err)
		return
	}
	if taken {
		respond.Conflict(c, "A role with that name already exists")
		return
	}

	role := models.Role{Name: in.Name, Description: in.Description}
	if err := role.SetGrants(in.Grants); err != nil {
		respond.Internal(c, err)
		return
	}
	// The service drops the grant cache when the write lands.
	if err := h.roles().Create(c.Request.Context(), &role); err != nil {
		respond.Internal(c, err)
		return
	}
	respond.Created(c, h.toDTO(role, 0), "Role created")
}

func (h *RoleHandler) Update(c *gin.Context) {
	role, err := h.roles().ByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respond.NotFound(c, "Role not found")
		return
	}

	var in RoleRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, "Name is required")
		return
	}
	if bad := validateGrants(in.Grants); len(bad) > 0 {
		respond.BadRequest(c, "Unknown permissions: "+strings.Join(bad, ", "))
		return
	}
	// A role hands out at most what its author holds, or a roles.create or
	// roles.edit holder could write "*" into a role and take it.
	if beyond := authz.CannotGrant(c, in.Grants); len(beyond) > 0 {
		respond.Forbidden(c, "You cannot grant permissions you do not hold: "+strings.Join(beyond, ", "))
		return
	}

	// Only an ADMIN changes a built-in role: taking "*" out of ADMIN locks out
	// every administrator.
	if role.IsSystem && !authz.IsAdmin(c) {
		respond.Forbidden(c, "Only an ADMIN can change a built-in role")
		return
	}

	// System roles keep their name. Routes and the legacy role fallback resolve
	// by name, so renaming ADMIN silently strips access. Enforced here, not only
	// in the UI, because the UI is not a security boundary.
	newName := strings.TrimSpace(in.Name)
	if role.IsSystem && !strings.EqualFold(newName, role.Name) {
		respond.Forbidden(c, "Built-in roles cannot be renamed")
		return
	}
	if !role.IsSystem && newName != "" {
		role.Name = newName
	}

	role.Description = in.Description
	if err := role.SetGrants(in.Grants); err != nil {
		respond.Internal(c, err)
		return
	}
	if err := h.roles().Save(c.Request.Context(), role); err != nil {
		respond.Internal(c, err)
		return
	}

	n, err := h.roles().Holders(c.Request.Context(), role.ID)
	if err != nil {
		respond.Internal(c, err)
		return
	}
	respond.OK(c, h.toDTO(*role, n), "Role updated")
}

func (h *RoleHandler) Delete(c *gin.Context) {
	role, err := h.roles().ByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respond.NotFound(c, "Role not found")
		return
	}
	if role.IsSystem {
		respond.Forbidden(c, "Built-in roles cannot be deleted")
		return
	}

	// A role still assigned to users cannot be deleted: removing it would
	// silently strip those users of the permissions it grants. Reassign or
	// remove them first, then the role becomes deletable.
	//
	// Both ways a role reaches a user are counted, which is the service's
	// Assignments: see why there.
	assigned, err := h.roles().Assignments(c.Request.Context(), role)
	if err != nil {
		respond.Internal(c, err)
		return
	}
	if assigned > 0 {
		respond.BadRequest(c, fmt.Sprintf("Role is assigned to %d user(s); reassign them before deleting it", assigned))
		return
	}

	if err := h.roles().Delete(c.Request.Context(), role); err != nil {
		respond.Internal(c, err)
		return
	}
	respond.OK(c, gin.H{"id": role.ID}, "Role deleted")
}

type AssignRolesRequest struct {
	RoleIDs []string `json:"role_ids"`
}

// AssignUserRoles replaces a user's role assignments.
func (h *RoleHandler) AssignUserRoles(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.users().GetByID(c.Request.Context(), userID)
	if err != nil {
		respond.NotFound(c, "User not found")
		return
	}

	var in AssignRolesRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, "role_ids is required")
		return
	}

	// Reject unknown role ids rather than silently dropping them. The rows are
	// read once here and used again by the ceiling check below, which used to
	// count them and then fetch them.
	roles, err := h.roles().ByIDs(c.Request.Context(), in.RoleIDs)
	if err != nil {
		respond.Internal(c, err)
		return
	}
	if len(roles) != len(in.RoleIDs) {
		respond.BadRequest(c, "One or more roles do not exist")
		return
	}

	// The ceiling applies to assignment too: only an ADMIN changes an ADMIN's
	// roles or hands out ADMIN, and nobody else hands out a role that grants more
	// than they hold. Before, a users.edit holder could assign themselves ADMIN.
	if !authz.IsAdmin(c) {
		if authz.IsAdminAccount(h.DB, user) {
			respond.Forbidden(c, "Only an ADMIN can change an ADMIN's roles")
			return
		}
		for i := range roles {
			if strings.EqualFold(roles[i].Name, models.RoleAdmin) {
				respond.Forbidden(c, "Only an ADMIN can assign the ADMIN role")
				return
			}
			if beyond := authz.CannotGrant(c, roles[i].GrantsList()); len(beyond) > 0 {
				respond.Forbidden(c, fmt.Sprintf("The %s role grants permissions you do not hold: %s", roles[i].Name, strings.Join(beyond, ", ")))
				return
			}
		}
	}

	// The assignments and the legacy users.role column are one change: see
	// ReplaceUserRoles.
	if err := h.roles().ReplaceUserRoles(c.Request.Context(), userID, in.RoleIDs); err != nil {
		respond.Internal(c, err)
		return
	}

	respond.OK(c, gin.H{"user_id": userID, "role_ids": in.RoleIDs}, "Roles updated")
}

// MyPermissions returns the caller's expanded permissions, for the frontend's
// can() helper and nav gating.
func (h *RoleHandler) MyPermissions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := userID.(string)

	grants, err := authz.GrantsFor(h.DB, id)
	if err != nil {
		respond.Internal(c, err)
		return
	}
	respond.OK(c, gin.H{
		"grants":      grants,
		"permissions": authz.Expand(grants),
		"is_super":    authz.HasAll(grants),
	})
}
