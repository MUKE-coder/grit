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

type BranchHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

// List returns all branches for the active business.
func (h *BranchHandler) List(c *gin.Context) {
	businessID, _ := c.Get("business_id")

	var branches []models.Branch
	h.DB.Where("business_id = ?", businessID).Order("is_default DESC, name ASC").Find(&branches)

	c.JSON(http.StatusOK, gin.H{"data": branches})
}

// Create creates a new branch (admin only).
func (h *BranchHandler) Create(c *gin.Context) {
	businessID := c.GetUint("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can create branches"})
		return
	}

	var req struct {
		Name    string `json:"name" binding:"required"`
		Address string `json:"address"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch name is required"})
		return
	}

	branch := models.Branch{
		BusinessID: businessID,
		Name:       strings.TrimSpace(req.Name),
		Address:    strings.TrimSpace(req.Address),
		IsDefault:  false,
	}

	if err := h.DB.Create(&branch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create branch"})
		return
	}

	h.Audit.Log(c, models.ActionCreate, "branch", &branch.ID, fmt.Sprintf("Created branch %q", branch.Name), nil)

	c.JSON(http.StatusCreated, gin.H{"data": branch})
}

// Update updates a branch (admin only).
func (h *BranchHandler) Update(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can update branches"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	var branch models.Branch
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&branch).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}

	var req struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// The default Main Branch is permanent — it backstops every business so
	// users always have somewhere to attach inventory + sales. Block edits
	// to its name/address; admins who really want to "rename" the default
	// can create a second branch and mark that as their main one in code,
	// or remove the is_default flag manually via the DB.
	if branch.IsDefault {
		c.JSON(http.StatusConflict, gin.H{"error": "The default Main Branch can't be edited"})
		return
	}

	if req.Name != "" {
		branch.Name = strings.TrimSpace(req.Name)
	}
	if req.Address != "" {
		branch.Address = strings.TrimSpace(req.Address)
	}

	h.DB.Save(&branch)
	h.Audit.Log(c, models.ActionUpdate, "branch", &branch.ID, fmt.Sprintf("Updated branch %q", branch.Name), nil)
	c.JSON(http.StatusOK, gin.H{"data": branch})
}

// Delete removes a non-default branch. The is_default branch is permanent —
// requesting its deletion returns 409 Conflict.
func (h *BranchHandler) Delete(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can delete branches"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	var branch models.Branch
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&branch).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}
	if branch.IsDefault {
		c.JSON(http.StatusConflict, gin.H{"error": "The default Main Branch can't be deleted"})
		return
	}

	// Refuse if any inventory / motorcycles / loans reference this branch.
	var attached int64
	h.DB.Model(&models.Motorcycle{}).Where("branch_id = ?", branch.ID).Count(&attached)
	if attached > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Branch has motorcycles attached. Transfer them first."})
		return
	}
	h.DB.Model(&models.Borrower{}).Where("branch_id = ?", branch.ID).Count(&attached)
	if attached > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Branch has borrowers attached."})
		return
	}
	h.DB.Model(&models.Stock{}).Where("branch_id = ? AND quantity > 0", branch.ID).Count(&attached)
	if attached > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Branch still has stock. Transfer it out first."})
		return
	}

	bid := branch.ID
	name := branch.Name
	h.DB.Delete(&branch)
	h.Audit.Log(c, models.ActionDelete, "branch", &bid, fmt.Sprintf("Deleted branch %q", name), nil)
	c.JSON(http.StatusOK, gin.H{"message": "Branch deleted"})
}
