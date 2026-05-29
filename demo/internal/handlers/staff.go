package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/models"
	"gritdemo/internal/services"
)

type StaffHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

// validRoles is the canonical set of roles a staff member can hold in a
// business. Mirrors the constants in user_business_role.go.
var validRoles = map[string]bool{
	models.RoleAdmin:       true,
	models.RoleManager:     true,
	models.RoleCashier:     true,
	models.RoleStockClerk:  true,
	models.RoleLoanOfficer: true,
	models.RoleAccountant:  true,
}

var validWorkspaces = map[string]bool{
	models.WorkspaceLoans:  true,
	models.WorkspaceSpares: true,
	models.WorkspaceBoth:   true,
}

// List returns all staff members in the active business with their roles.
func (h *StaffHandler) List(c *gin.Context) {
	businessID, _ := c.Get("business_id")

	var roles []models.UserBusinessRole
	h.DB.Preload("User").Preload("Branch").
		Where("business_id = ?", businessID).
		Find(&roles)

	type StaffMember struct {
		UserID          uint   `json:"user_id"`
		Name            string `json:"name"`
		Email           string `json:"email"`
		Role            string `json:"role"`
		BranchID        *uint  `json:"branch_id"`
		BranchName      string `json:"branch_name,omitempty"`
		WorkspaceAccess string `json:"workspace_access"`
	}

	staff := make([]StaffMember, 0, len(roles))
	for _, r := range roles {
		ws := r.WorkspaceAccess
		if ws == "" {
			ws = models.WorkspaceBoth
		}
		s := StaffMember{
			UserID:          r.UserID,
			Name:            r.User.Name,
			Email:           r.User.Email,
			Role:            r.Role,
			BranchID:        r.BranchID,
			WorkspaceAccess: ws,
		}
		if r.Branch != nil {
			s.BranchName = r.Branch.Name
		}
		staff = append(staff, s)
	}

	c.JSON(http.StatusOK, gin.H{"data": staff})
}

// Create directly creates a staff user in the active business (admin only).
// This is the primary "add a staff member" path; invitations are kept as
// a fallback for self-service onboarding by email.
//
// If the email already belongs to a User, we link that user to the business
// instead of creating a duplicate account. Email collisions inside the same
// business return 409.
func (h *StaffHandler) Create(c *gin.Context) {
	businessID := c.GetUint("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can create staff"})
		return
	}

	var req struct {
		Name            string `json:"name" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=6"`
		Role            string `json:"role" binding:"required"`
		BranchID        *uint  `json:"branch_id"`
		WorkspaceAccess string `json:"workspace_access"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !validRoles[req.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}
	if req.WorkspaceAccess == "" {
		req.WorkspaceAccess = models.WorkspaceBoth
	}
	if !validWorkspaces[req.WorkspaceAccess] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace_access"})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Find or create user. We allow re-using an existing user account
	// (e.g. a contractor who works for two businesses) — but only if they
	// don't already have a role in THIS business.
	var user models.User
	err := h.DB.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		user = models.User{
			Name:  strings.TrimSpace(req.Name),
			Email: email,
		}
		if err := user.SetPassword(req.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		if err := h.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to look up user"})
		return
	} else {
		// User exists — make sure they're not already on this business.
		var existing int64
		h.DB.Model(&models.UserBusinessRole{}).
			Where("user_id = ? AND business_id = ?", user.ID, businessID).
			Count(&existing)
		if existing > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "This email already has access to this business"})
			return
		}
		// Optionally update name on the existing user record so it shows
		// the friendly name the admin just typed.
		if user.Name == "" && req.Name != "" {
			h.DB.Model(&user).Update("name", strings.TrimSpace(req.Name))
		}
	}

	ubr := models.UserBusinessRole{
		UserID:          user.ID,
		BusinessID:      businessID,
		BranchID:        req.BranchID,
		Role:            req.Role,
		WorkspaceAccess: req.WorkspaceAccess,
	}
	if err := h.DB.Create(&ubr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role"})
		return
	}

	h.Audit.Log(c, models.ActionCreate, "staff", &user.ID,
		fmt.Sprintf("Created staff %s (%s) as %s", user.Name, user.Email, req.Role),
		gin.H{"role": req.Role, "branch_id": req.BranchID, "workspace_access": req.WorkspaceAccess},
	)

	c.JSON(http.StatusCreated, gin.H{"data": gin.H{
		"user_id":          user.ID,
		"name":             user.Name,
		"email":            user.Email,
		"role":             ubr.Role,
		"branch_id":        ubr.BranchID,
		"workspace_access": ubr.WorkspaceAccess,
	}})
}

// UpdateRole updates a staff member's role, branch, or workspace access.
func (h *StaffHandler) UpdateRole(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can update staff roles"})
		return
	}

	userIDParam, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Role            string `json:"role" binding:"required"`
		BranchID        *uint  `json:"branch_id"`
		WorkspaceAccess string `json:"workspace_access"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role is required"})
		return
	}

	if !validRoles[req.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}
	if req.WorkspaceAccess != "" && !validWorkspaces[req.WorkspaceAccess] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace_access"})
		return
	}

	var ubr models.UserBusinessRole
	if err := h.DB.Where("user_id = ? AND business_id = ?", userIDParam, businessID).First(&ubr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff member not found in this business"})
		return
	}

	ubr.Role = req.Role
	ubr.BranchID = req.BranchID
	if req.WorkspaceAccess != "" {
		ubr.WorkspaceAccess = req.WorkspaceAccess
	}
	h.DB.Save(&ubr)

	h.Audit.Log(c, models.ActionUpdate, "staff", &ubr.UserID,
		fmt.Sprintf("Updated staff role to %s", ubr.Role),
		gin.H{"role": ubr.Role, "branch_id": ubr.BranchID, "workspace_access": ubr.WorkspaceAccess},
	)

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"user_id":          ubr.UserID,
		"role":             ubr.Role,
		"branch_id":        ubr.BranchID,
		"workspace_access": ubr.WorkspaceAccess,
	}})
}

// Remove removes a staff member from the business (admin only).
func (h *StaffHandler) Remove(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	currentUserID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can remove staff"})
		return
	}

	userIDParam, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Can't remove yourself
	if uint(userIDParam) == currentUserID.(uint) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot remove yourself from the business"})
		return
	}

	result := h.DB.Where("user_id = ? AND business_id = ?", userIDParam, businessID).
		Delete(&models.UserBusinessRole{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff member not found in this business"})
		return
	}

	uid := uint(userIDParam)
	h.Audit.Log(c, models.ActionDelete, "staff", &uid, "Removed staff member from business", nil)

	c.JSON(http.StatusOK, gin.H{"message": "Staff member removed"})
}
