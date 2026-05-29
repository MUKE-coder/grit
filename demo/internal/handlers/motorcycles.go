package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/models"
	"gritdemo/internal/services"
	"gritdemo/internal/storage"
)

type MotorcycleHandler struct {
	DB      *gorm.DB
	Storage *storage.Storage
	Audit   *services.AuditService
}

func (h *MotorcycleHandler) response(m models.Motorcycle) gin.H {
	resp := gin.H{
		"id":            m.ID,
		"business_id":   m.BusinessID,
		"branch_id":     m.BranchID,
		"name":          m.Name,
		"number_plate":  m.NumberPlate,
		"chassis_no":    m.ChassisNo,
		"engine_no":     m.EngineNo,
		"color":         m.Color,
		"year_of_make":  m.YearOfMake,
		"cost_price":    m.CostPrice,
		"selling_price": m.SellingPrice,
		"status":        m.Status,
		"notes":         m.Notes,
		"created_at":    m.CreatedAt,
		"updated_at":    m.UpdatedAt,
	}
	if m.Branch.ID != 0 {
		resp["branch"] = gin.H{"id": m.Branch.ID, "name": m.Branch.Name}
	}
	if m.ImageKey != "" && h.Storage != nil {
		if url, err := h.Storage.GetSignedURL(context.Background(), m.ImageKey, time.Hour); err == nil {
			resp["image_url"] = url
		}
	}
	return resp
}

// List returns motorcycles for the active business with optional filters.
func (h *MotorcycleHandler) List(c *gin.Context) {
	businessID := c.GetUint("business_id")
	q := h.DB.Preload("Branch").Where("business_id = ?", businessID)

	if branchID := c.Query("branch_id"); branchID != "" {
		q = q.Where("branch_id = ?", branchID)
	}
	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if search := c.Query("search"); search != "" {
		like := "%" + strings.ToLower(search) + "%"
		q = q.Where("LOWER(name) LIKE ? OR LOWER(number_plate) LIKE ?", like, like)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "50"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 200 {
		perPage = 50
	}

	var total int64
	q.Model(&models.Motorcycle{}).Count(&total)

	var rows []models.Motorcycle
	q.Order("created_at DESC").Offset((page - 1) * perPage).Limit(perPage).Find(&rows)

	out := make([]gin.H, 0, len(rows))
	for _, m := range rows {
		out = append(out, h.response(m))
	}
	c.JSON(http.StatusOK, gin.H{"data": out, "total": total, "page": page, "per_page": perPage})
}

// Get returns a single motorcycle.
func (h *MotorcycleHandler) Get(c *gin.Context) {
	businessID := c.GetUint("business_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var m models.Motorcycle
	if err := h.DB.Preload("Branch").Where("id = ? AND business_id = ?", id, businessID).First(&m).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Motorcycle not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": h.response(m)})
}

// Create adds a motorcycle to inventory.
func (h *MotorcycleHandler) Create(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	if !canManageMotorcycles(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req struct {
		BranchID     uint    `json:"branch_id" binding:"required"`
		Name         string  `json:"name" binding:"required"`
		NumberPlate  string  `json:"number_plate" binding:"required"`
		ChassisNo    string  `json:"chassis_no"`
		EngineNo     string  `json:"engine_no"`
		Color        string  `json:"color"`
		YearOfMake   int     `json:"year_of_make"`
		CostPrice    float64 `json:"cost_price"`
		SellingPrice float64 `json:"selling_price" binding:"required"`
		Notes        string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch, name, number_plate and selling_price are required"})
		return
	}

	if !validateBranchInBusiness(h.DB, req.BranchID, businessID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch does not belong to this business"})
		return
	}

	m := models.Motorcycle{
		BusinessID:   businessID,
		BranchID:     req.BranchID,
		Name:         strings.TrimSpace(req.Name),
		NumberPlate:  strings.ToUpper(strings.ReplaceAll(req.NumberPlate, " ", "")),
		ChassisNo:    strings.TrimSpace(req.ChassisNo),
		EngineNo:     strings.TrimSpace(req.EngineNo),
		Color:        req.Color,
		YearOfMake:   req.YearOfMake,
		CostPrice:    req.CostPrice,
		SellingPrice: req.SellingPrice,
		Notes:        req.Notes,
		Status:       models.MotorcycleStatusAvailable,
		CreatedBy:    userID,
	}
	if err := h.DB.Create(&m).Error; err != nil {
		if isUniqueViolation(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "A motorcycle with that number plate already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create motorcycle"})
		return
	}

	h.DB.Preload("Branch").First(&m, m.ID)
	h.Audit.Log(c, models.ActionCreate, "motorcycle", &m.ID, fmt.Sprintf("Created motorcycle %q (%s)", m.Name, m.NumberPlate), gin.H{"branch_id": m.BranchID, "selling_price": m.SellingPrice})
	c.JSON(http.StatusCreated, gin.H{"data": h.response(m)})
}

// Update edits motorcycle metadata. Status changes go through dedicated lifecycle endpoints.
func (h *MotorcycleHandler) Update(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if !canManageMotorcycles(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var m models.Motorcycle
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&m).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Motorcycle not found"})
		return
	}

	var req struct {
		BranchID     *uint    `json:"branch_id"`
		Name         *string  `json:"name"`
		ChassisNo    *string  `json:"chassis_no"`
		EngineNo     *string  `json:"engine_no"`
		Color        *string  `json:"color"`
		YearOfMake   *int     `json:"year_of_make"`
		CostPrice    *float64 `json:"cost_price"`
		SellingPrice *float64 `json:"selling_price"`
		Notes        *string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if req.BranchID != nil {
		if !validateBranchInBusiness(h.DB, *req.BranchID, businessID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Branch does not belong to this business"})
			return
		}
		m.BranchID = *req.BranchID
	}
	if req.Name != nil {
		m.Name = strings.TrimSpace(*req.Name)
	}
	if req.ChassisNo != nil {
		m.ChassisNo = *req.ChassisNo
	}
	if req.EngineNo != nil {
		m.EngineNo = *req.EngineNo
	}
	if req.Color != nil {
		m.Color = *req.Color
	}
	if req.YearOfMake != nil {
		m.YearOfMake = *req.YearOfMake
	}
	if req.CostPrice != nil {
		m.CostPrice = *req.CostPrice
	}
	if req.SellingPrice != nil {
		m.SellingPrice = *req.SellingPrice
	}
	if req.Notes != nil {
		m.Notes = *req.Notes
	}

	h.DB.Save(&m)
	h.DB.Preload("Branch").First(&m, m.ID)
	h.Audit.Log(c, models.ActionUpdate, "motorcycle", &m.ID, fmt.Sprintf("Updated motorcycle %q (%s)", m.Name, m.NumberPlate), nil)
	c.JSON(http.StatusOK, gin.H{"data": h.response(m)})
}

// Delete soft-deletes a motorcycle (only if it's never been sold or financed).
func (h *MotorcycleHandler) Delete(c *gin.Context) {
	businessID := c.GetUint("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can delete motorcycles"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var m models.Motorcycle
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&m).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Motorcycle not found"})
		return
	}
	if m.Status == models.MotorcycleStatusSold || m.Status == models.MotorcycleStatusOnLoan {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete a motorcycle that has been sold or financed"})
		return
	}

	mid := uint(id)
	name := m.Name
	plate := m.NumberPlate
	h.DB.Delete(&m)
	h.Audit.Log(c, models.ActionDelete, "motorcycle", &mid, fmt.Sprintf("Deleted motorcycle %q (%s)", name, plate), nil)
	c.JSON(http.StatusOK, gin.H{"message": "Motorcycle deleted"})
}

// Transfer moves a motorcycle to another branch. Only allowed when available.
func (h *MotorcycleHandler) Transfer(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if !canManageMotorcycles(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		BranchID uint   `json:"branch_id" binding:"required"`
		Note     string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch ID is required"})
		return
	}
	if !validateBranchInBusiness(h.DB, req.BranchID, businessID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch does not belong to this business"})
		return
	}

	var m models.Motorcycle
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&m).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Motorcycle not found"})
		return
	}
	if m.Status != models.MotorcycleStatusAvailable {
		c.JSON(http.StatusConflict, gin.H{"error": "Only available motorcycles can be transferred"})
		return
	}

	m.BranchID = req.BranchID
	if req.Note != "" {
		m.Notes = strings.TrimSpace(m.Notes + "\n[Transfer] " + req.Note)
	}
	h.DB.Save(&m)
	h.DB.Preload("Branch").First(&m, m.ID)
	h.Audit.Log(c, models.ActionTransfer, "motorcycle", &m.ID, fmt.Sprintf("Transferred motorcycle %q (%s) to branch %s", m.Name, m.NumberPlate, m.Branch.Name), gin.H{"branch_id": m.BranchID, "note": req.Note})
	c.JSON(http.StatusOK, gin.H{"data": h.response(m)})
}

// UploadImage stores a motorcycle photo on R2.
func (h *MotorcycleHandler) UploadImage(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if !canManageMotorcycles(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	if h.Storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Storage not configured"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var m models.Motorcycle
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&m).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Motorcycle not found"})
		return
	}
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image"})
		return
	}

	if m.ImageKey != "" {
		_ = h.Storage.Delete(context.Background(), m.ImageKey)
	}
	key := fmt.Sprintf("motorcycles/%d/%d-%s", businessID, m.ID, header.Filename)
	if err := h.Storage.Upload(context.Background(), key, file, contentType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}
	m.ImageKey = key
	h.DB.Save(&m)
	url, _ := h.Storage.GetSignedURL(context.Background(), key, time.Hour)
	h.Audit.Log(c, models.ActionUpdate, "motorcycle", &m.ID, fmt.Sprintf("Uploaded image for motorcycle %q (%s)", m.Name, m.NumberPlate), nil)
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"image_url": url}})
}

func canManageMotorcycles(c *gin.Context) bool {
	role, _ := c.Get("user_role")
	return role == models.RoleAdmin || role == models.RoleManager
}
