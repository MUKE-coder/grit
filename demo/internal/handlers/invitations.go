package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"gritdemo/internal/config"
	"gritdemo/internal/mail"
	"gritdemo/internal/models"
	"gritdemo/internal/services"
)

type InvitationHandler struct {
	DB          *gorm.DB
	Mailer      *mail.Mailer
	AuthService *services.AuthService
	Config      *config.Config
	Audit       *services.AuditService
}

// Create creates an invitation and sends an email (admin only).
func (h *InvitationHandler) Create(c *gin.Context) {
	businessID := c.GetUint("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can invite staff"})
		return
	}

	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Role     string `json:"role" binding:"required"`
		BranchID *uint  `json:"branch_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and role are required"})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Validate role
	validRoles := map[string]bool{
		models.RoleAdmin: true, models.RoleManager: true,
		models.RoleCashier: true, models.RoleStockClerk: true,
	}
	if !validRoles[req.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	// Check if already a member
	var existingRole models.UserBusinessRole
	var existingUser models.User
	if h.DB.Where("email = ?", req.Email).First(&existingUser).Error == nil {
		if h.DB.Where("user_id = ? AND business_id = ?", existingUser.ID, businessID).First(&existingRole).Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "This user is already a member of this business"})
			return
		}
	}

	// Check for pending invitation
	var existingInvite models.Invitation
	if h.DB.Where("email = ? AND business_id = ? AND accepted_at IS NULL AND expires_at > ?",
		req.Email, businessID, time.Now()).First(&existingInvite).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "An invitation has already been sent to this email"})
		return
	}

	// Create invitation
	token := uuid.New().String()
	invitation := models.Invitation{
		BusinessID: businessID,
		BranchID:   req.BranchID,
		Email:      req.Email,
		Role:       req.Role,
		Token:      token,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := h.DB.Create(&invitation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invitation"})
		return
	}

	// Get business name for email
	var business models.Business
	h.DB.First(&business, businessID)

	// Send invitation email
	inviteURL := h.Config.FrontendURL + "/invite/" + token
	if h.Mailer != nil {
		go func() {
			_ = h.Mailer.SendInvitation(req.Email, business.Name, req.Role, inviteURL)
		}()
	}

	h.Audit.Log(c, models.ActionCreate, "invitation", &invitation.ID, fmt.Sprintf("Invited %s as %s", invitation.Email, invitation.Role), gin.H{"branch_id": invitation.BranchID})
	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":         invitation.ID,
			"email":      invitation.Email,
			"role":       invitation.Role,
			"expires_at": invitation.ExpiresAt,
			"invite_url": inviteURL,
		},
	})
}

// List returns pending invitations for the active business (admin only).
func (h *InvitationHandler) List(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can view invitations"})
		return
	}

	var invitations []models.Invitation
	h.DB.Preload("Branch").
		Where("business_id = ? AND accepted_at IS NULL", businessID).
		Order("created_at DESC").
		Find(&invitations)

	c.JSON(http.StatusOK, gin.H{"data": invitations})
}

// Revoke deletes a pending invitation (admin only).
func (h *InvitationHandler) Revoke(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can revoke invitations"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invitation ID"})
		return
	}

	result := h.DB.Where("id = ? AND business_id = ? AND accepted_at IS NULL", id, businessID).
		Delete(&models.Invitation{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invitation not found"})
		return
	}

	iid := uint(id)
	h.Audit.Log(c, models.ActionDelete, "invitation", &iid, "Revoked invitation", nil)
	c.JSON(http.StatusOK, gin.H{"message": "Invitation revoked"})
}

// GetInvite validates an invitation token (public, no auth required).
func (h *InvitationHandler) GetInvite(c *gin.Context) {
	token := c.Param("token")

	var invitation models.Invitation
	if err := h.DB.Preload("Business").Where("token = ?", token).First(&invitation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invitation not found"})
		return
	}

	if invitation.AcceptedAt != nil {
		c.JSON(http.StatusGone, gin.H{"error": "Invitation has already been accepted"})
		return
	}

	if time.Now().After(invitation.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "Invitation has expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"email":         invitation.Email,
			"role":          invitation.Role,
			"business_name": invitation.Business.Name,
		},
	})
}

// AcceptInvite accepts an invitation (public, no auth required).
func (h *InvitationHandler) AcceptInvite(c *gin.Context) {
	token := c.Param("token")

	var invitation models.Invitation
	if err := h.DB.Where("token = ?", token).First(&invitation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invitation not found"})
		return
	}

	if invitation.AcceptedAt != nil {
		c.JSON(http.StatusGone, gin.H{"error": "Invitation has already been accepted"})
		return
	}

	if time.Now().After(invitation.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "Invitation has expired"})
		return
	}

	var req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	c.ShouldBindJSON(&req)

	var user models.User
	isNewUser := false

	// Check if user already exists
	if err := h.DB.Where("email = ?", invitation.Email).First(&user).Error; err != nil {
		// New user — name and password required
		if req.Name == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Name and password are required for new users"})
			return
		}
		isNewUser = true
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if isNewUser {
			user = models.User{
				Name:  strings.TrimSpace(req.Name),
				Email: invitation.Email,
			}
			if err := user.SetPassword(req.Password); err != nil {
				return err
			}
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
		}

		// Create role
		role := models.UserBusinessRole{
			UserID:     user.ID,
			BusinessID: invitation.BusinessID,
			BranchID:   invitation.BranchID,
			Role:       invitation.Role,
		}
		if err := tx.Create(&role).Error; err != nil {
			return err
		}

		// Mark invitation as accepted
		now := time.Now()
		invitation.AcceptedAt = &now
		return tx.Save(&invitation).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept invitation"})
		return
	}

	// Generate JWT
	tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
			"tokens": tokens,
		},
	})
}
