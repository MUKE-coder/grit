package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/config"
	"gritdemo/internal/models"
	"gritdemo/internal/services"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	DB          *gorm.DB
	AuthService *services.AuthService
	Config      *config.Config
	Audit       *services.AuditService
}

// RegisterRequest is the request body for registration.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest is the request body for login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register creates a new user with a default business, branch, and admin role.
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	// Check if email already exists
	var existing models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Use transaction to create user + business + branch + role
	var user models.User
	var business models.Business
	var branch models.Branch

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create User
		user = models.User{
			Name:  req.Name,
			Email: req.Email,
		}
		if err := user.SetPassword(req.Password); err != nil {
			return err
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 2. Create default Business
		business = models.Business{
			Name:    req.Name + "'s Business",
			OwnerID: user.ID,
		}
		if err := tx.Create(&business).Error; err != nil {
			return err
		}

		// 3. Create default Branch
		branch = models.Branch{
			BusinessID: business.ID,
			Name:       "Main Branch",
			IsDefault:  true,
		}
		if err := tx.Create(&branch).Error; err != nil {
			return err
		}

		// 4. Create UserBusinessRole (admin, all branches)
		role := models.UserBusinessRole{
			UserID:     user.ID,
			BusinessID: business.ID,
			BranchID:   nil, // nil = all branches
			Role:       models.RoleAdmin,
		}
		if err := tx.Create(&role).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	// Generate JWT
	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
			"business": gin.H{
				"id":   business.ID,
				"name": business.Name,
			},
			"branch": gin.H{
				"id":   branch.ID,
				"name": branch.Name,
			},
			"tokens": tokens,
		},
	})
}

// Login authenticates a user and returns JWT tokens.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Find user
	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Get all businesses the user belongs to
	var roles []models.UserBusinessRole
	h.DB.Preload("Business").Preload("Branch").Where("user_id = ?", user.ID).Find(&roles)

	type BusinessInfo struct {
		ID              uint   `json:"id"`
		Name            string `json:"name"`
		Role            string `json:"role"`
		BranchID        *uint  `json:"branch_id"`
		WorkspaceAccess string `json:"workspace_access"`
	}

	businesses := make([]BusinessInfo, 0, len(roles))
	for _, r := range roles {
		ws := r.WorkspaceAccess
		if ws == "" {
			ws = models.WorkspaceBoth
		}
		businesses = append(businesses, BusinessInfo{
			ID:              r.BusinessID,
			Name:            r.Business.Name,
			Role:            r.Role,
			BranchID:        r.BranchID,
			WorkspaceAccess: ws,
		})
	}

	// Generate JWT
	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Audit: log a login event into every business this user belongs to so
	// the admin of any of them can see it. We use LogFor since gin.Context
	// doesn't yet have business_id set during login.
	for _, r := range roles {
		h.Audit.LogFor(c, r.BusinessID, user.ID,
			models.ActionLogin, "auth", nil,
			fmt.Sprintf("%s signed in", user.Name), nil,
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
			"businesses": businesses,
			"tokens":     tokens,
		},
	})
}

// Me returns the current user's profile and their businesses/roles.
func (h *AuthHandler) Me(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	user := userVal.(models.User)

	// Get all businesses and roles
	var roles []models.UserBusinessRole
	h.DB.Preload("Business").Preload("Branch").Where("user_id = ?", user.ID).Find(&roles)

	type RoleInfo struct {
		BusinessID      uint   `json:"business_id"`
		BusinessName    string `json:"business_name"`
		Role            string `json:"role"`
		BranchID        *uint  `json:"branch_id"`
		BranchName      string `json:"branch_name,omitempty"`
		WorkspaceAccess string `json:"workspace_access"`
	}

	roleInfos := make([]RoleInfo, 0, len(roles))
	for _, r := range roles {
		ws := r.WorkspaceAccess
		if ws == "" {
			ws = models.WorkspaceBoth
		}
		ri := RoleInfo{
			BusinessID:      r.BusinessID,
			BusinessName:    r.Business.Name,
			Role:            r.Role,
			BranchID:        r.BranchID,
			WorkspaceAccess: ws,
		}
		if r.Branch != nil {
			ri.BranchName = r.Branch.Name
		}
		roleInfos = append(roleInfos, ri)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
			"roles": roleInfos,
		},
	})
}

// Logout acknowledges a logout (stateless JWT — client deletes token).
func (h *AuthHandler) Logout(c *gin.Context) {
	// business_id may not be set if the client doesn't send X-Business-ID;
	// the audit helper will skip silently in that case.
	if userVal, ok := c.Get("user"); ok {
		if u, ok := userVal.(models.User); ok {
			h.Audit.Log(c, models.ActionLogout, "auth", nil, fmt.Sprintf("%s signed out", u.Name), nil)
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Refresh validates a refresh token and returns a new token pair.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	claims, err := h.AuthService.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Verify user still exists
	var user models.User
	if err := h.DB.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"tokens": tokens,
		},
	})
}
