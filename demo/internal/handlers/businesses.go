package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/models"
)

type BusinessHandler struct {
	DB *gorm.DB
}

// List returns all businesses the current user belongs to.
func (h *BusinessHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var roles []models.UserBusinessRole
	h.DB.Preload("Business").Where("user_id = ?", userID).Find(&roles)

	type BusinessResponse struct {
		ID      uint   `json:"id"`
		Name    string `json:"name"`
		Role    string `json:"role"`
		OwnerID uint   `json:"owner_id"`
	}

	businesses := make([]BusinessResponse, 0, len(roles))
	for _, r := range roles {
		businesses = append(businesses, BusinessResponse{
			ID:      r.Business.ID,
			Name:    r.Business.Name,
			Role:    r.Role,
			OwnerID: r.Business.OwnerID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": businesses})
}

// Create creates a new business with a default Main Branch and admin role.
func (h *BusinessHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Business name is required"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)

	var business models.Business
	var branch models.Branch

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		business = models.Business{
			Name:    req.Name,
			OwnerID: userID,
		}
		if err := tx.Create(&business).Error; err != nil {
			return err
		}

		branch = models.Branch{
			BusinessID: business.ID,
			Name:       "Main Branch",
			IsDefault:  true,
		}
		if err := tx.Create(&branch).Error; err != nil {
			return err
		}

		role := models.UserBusinessRole{
			UserID:     userID,
			BusinessID: business.ID,
			BranchID:   nil,
			Role:       models.RoleAdmin,
		}
		return tx.Create(&role).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create business"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"business": business,
			"branch":   branch,
		},
	})
}

// Update updates a business name (admin only).
func (h *BusinessHandler) Update(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can update business settings"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Business name is required"})
		return
	}

	var business models.Business
	if err := h.DB.First(&business, businessID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Business not found"})
		return
	}

	business.Name = strings.TrimSpace(req.Name)
	h.DB.Save(&business)

	c.JSON(http.StatusOK, gin.H{"data": business})
}
