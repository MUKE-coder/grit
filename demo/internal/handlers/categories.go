package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/models"
)

type CategoryHandler struct {
	DB *gorm.DB
}

// List returns all categories for the active business.
func (h *CategoryHandler) List(c *gin.Context) {
	businessID, _ := c.Get("business_id")

	var categories []models.Category
	h.DB.Where("business_id = ?", businessID).Order("name ASC").Find(&categories)

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// Create creates a new category (admin/manager).
func (h *CategoryHandler) Create(c *gin.Context) {
	businessID := c.GetUint("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin && role != models.RoleManager {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins and managers can create categories"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category name is required"})
		return
	}

	name := strings.TrimSpace(req.Name)

	// Check for duplicate
	var existing models.Category
	if h.DB.Where("business_id = ? AND LOWER(name) = LOWER(?)", businessID, name).First(&existing).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Category already exists"})
		return
	}

	category := models.Category{
		BusinessID: businessID,
		Name:       name,
	}
	if err := h.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": category})
}

// Update renames a category.
func (h *CategoryHandler) Update(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin && role != models.RoleManager {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins and managers can update categories"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category models.Category
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category name is required"})
		return
	}

	category.Name = strings.TrimSpace(req.Name)
	h.DB.Save(&category)

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// Delete removes a category if no products use it.
func (h *CategoryHandler) Delete(c *gin.Context) {
	businessID, _ := c.Get("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin && role != models.RoleManager {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins and managers can delete categories"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Check if products use this category
	var count int64
	h.DB.Model(&models.Product{}).Where("category_id = ? AND business_id = ?", id, businessID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete category with existing products"})
		return
	}

	result := h.DB.Where("id = ? AND business_id = ?", id, businessID).Delete(&models.Category{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
