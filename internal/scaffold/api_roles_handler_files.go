package scaffold

import "strings"

// roleHandlerGo emits internal/handlers/role.go — the roles + permissions API.
//
// Endpoints:
//
//	GET    /api/admin/permissions        the catalog tree, for the UI
//	GET    /api/admin/roles              roles + user counts
//	POST   /api/admin/roles              create
//	GET    /api/admin/roles/:id          one role
//	PUT    /api/admin/roles/:id          update name/description/grants
//	DELETE /api/admin/roles/:id          delete (blocked for system roles)
//	PUT    /api/admin/users/:id/roles    assign roles to a user
//
// Guard rails learned from the reference implementation:
//
//   - System roles are protected SERVER-SIDE, not just greyed out in the UI.
//     Shoppleet only disables the name input in React, so a direct PUT can
//     rename ADMIN and quietly break every route that still names it.
//   - Every write calls authz.Invalidate(), so a revoked permission takes effect
//     on the next request rather than lingering in a cache.
//   - Grants are validated against the catalog. A typo like "product.create"
//     (singular) would otherwise be stored silently and simply never match.
func roleHandlerGo() string {
	src := `package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/authz"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/respond"
)

// RoleHandler serves the roles + permissions API.
type RoleHandler struct {
	DB *gorm.DB
}

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{DB: db}
}

// roleDTO is a Role plus the derived fields the UI needs.
type roleDTO struct {
	ID          string   ~json:"id"~
	Name        string   ~json:"name"~
	Description string   ~json:"description"~
	Grants      []string ~json:"grants"~
	// Expanded is grants with wildcards resolved against the catalog. The UI
	// renders checkboxes from this and never reimplements wildcard matching —
	// a duplicated matcher is how the reference implementation's Go and
	// TypeScript rules drifted apart.
	Expanded  []string ~json:"expanded"~
	IsSystem  bool     ~json:"is_system"~
	UserCount int64    ~json:"user_count"~
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
	var roles []models.Role
	if err := h.DB.Order("is_system DESC, name ASC").Find(&roles).Error; err != nil {
		respond.Internal(c, err)
		return
	}

	// One grouped count instead of a query per role.
	type countRow struct {
		RoleID string
		N      int64
	}
	var rows []countRow
	h.DB.Model(&models.UserRole{}).
		Select("role_id, count(*) as n").
		Group("role_id").
		Scan(&rows)
	counts := map[string]int64{}
	for _, r := range rows {
		counts[r.RoleID] = r.N
	}

	out := make([]roleDTO, 0, len(roles))
	for _, r := range roles {
		out = append(out, h.toDTO(r, counts[r.ID]))
	}
	respond.OK(c, out)
}

func (h *RoleHandler) Get(c *gin.Context) {
	var role models.Role
	if err := h.DB.Where("id = ?", c.Param("id")).First(&role).Error; err != nil {
		respond.NotFound(c, "Role not found")
		return
	}
	var n int64
	h.DB.Model(&models.UserRole{}).Where("role_id = ?", role.ID).Count(&n)
	respond.OK(c, h.toDTO(role, n))
}

type roleInput struct {
	Name        string   ~json:"name" binding:"required"~
	Description string   ~json:"description"~
	Grants      []string ~json:"grants"~
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
	var in roleInput
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

	var existing models.Role
	if err := h.DB.Where("name = ?", in.Name).First(&existing).Error; err == nil {
		respond.Conflict(c, "A role with that name already exists")
		return
	}

	role := models.Role{Name: in.Name, Description: in.Description}
	if err := role.SetGrants(in.Grants); err != nil {
		respond.Internal(c, err)
		return
	}
	if err := h.DB.Create(&role).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	authz.Invalidate()
	respond.Created(c, h.toDTO(role, 0), "Role created")
}

func (h *RoleHandler) Update(c *gin.Context) {
	var role models.Role
	if err := h.DB.Where("id = ?", c.Param("id")).First(&role).Error; err != nil {
		respond.NotFound(c, "Role not found")
		return
	}

	var in roleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, "Name is required")
		return
	}
	if bad := validateGrants(in.Grants); len(bad) > 0 {
		respond.BadRequest(c, "Unknown permissions: "+strings.Join(bad, ", "))
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
	if err := h.DB.Save(&role).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	authz.Invalidate()

	var n int64
	h.DB.Model(&models.UserRole{}).Where("role_id = ?", role.ID).Count(&n)
	respond.OK(c, h.toDTO(role, n), "Role updated")
}

func (h *RoleHandler) Delete(c *gin.Context) {
	var role models.Role
	if err := h.DB.Where("id = ?", c.Param("id")).First(&role).Error; err != nil {
		respond.NotFound(c, "Role not found")
		return
	}
	if role.IsSystem {
		respond.Forbidden(c, "Built-in roles cannot be deleted")
		return
	}

	// Drop the assignments too, or users keep a dangling role_id that resolves
	// to nothing and silently loses them permissions.
	if err := h.DB.Where("role_id = ?", role.ID).Delete(&models.UserRole{}).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	if err := h.DB.Delete(&role).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	authz.Invalidate()
	respond.OK(c, gin.H{"id": role.ID}, "Role deleted")
}

type assignInput struct {
	RoleIDs []string ~json:"role_ids"~
}

// AssignUserRoles replaces a user's role assignments.
func (h *RoleHandler) AssignUserRoles(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := h.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		respond.NotFound(c, "User not found")
		return
	}

	var in assignInput
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, "role_ids is required")
		return
	}

	// Reject unknown role ids rather than silently dropping them.
	if len(in.RoleIDs) > 0 {
		var n int64
		h.DB.Model(&models.Role{}).Where("id IN ?", in.RoleIDs).Count(&n)
		if int(n) != len(in.RoleIDs) {
			respond.BadRequest(c, "One or more roles do not exist")
			return
		}
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
			return err
		}
		for _, rid := range in.RoleIDs {
			if err := tx.Create(&models.UserRole{UserID: userID, RoleID: rid}).Error; err != nil {
				return err
			}
		}
		// Keep the legacy users.role string in step with the primary role.
		// Grant resolution prefers user_roles, but routes still guarded by a
		// role NAME read this column — leaving it stale is how the reference
		// implementation produced spurious 403s.
		legacy := ""
		if len(in.RoleIDs) > 0 {
			var first models.Role
			if err := tx.Where("id = ?", in.RoleIDs[0]).First(&first).Error; err == nil {
				legacy = strings.ToUpper(first.Name)
			}
		}
		return tx.Model(&models.User{}).Where("id = ?", userID).Update("role", legacy).Error
	})
	if err != nil {
		respond.Internal(c, err)
		return
	}
	authz.Invalidate()

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
`
	return strings.ReplaceAll(src, "~", "`")
}
