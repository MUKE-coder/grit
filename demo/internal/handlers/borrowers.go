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

type BorrowerHandler struct {
	DB      *gorm.DB
	Storage *storage.Storage
	Audit   *services.AuditService
}

func (h *BorrowerHandler) response(b models.Borrower) gin.H {
	resp := gin.H{
		"id":                  b.ID,
		"business_id":         b.BusinessID,
		"branch_id":           b.BranchID,
		"first_name":          b.FirstName,
		"last_name":           b.LastName,
		"full_name":           b.FullName(),
		"phone":               b.Phone,
		"alt_phone":           b.AltPhone,
		"email":               b.Email,
		"date_of_birth":       b.DateOfBirth,
		"gender":              b.Gender,
		"national_id":         b.NationalID,
		"address":             b.Address,
		"employment_status":   b.EmploymentStatus,
		"occupation":          b.Occupation,
		"employer":            b.Employer,
		"monthly_income":      b.MonthlyIncome,
		"next_of_kin_name":    b.NextOfKinName,
		"next_of_kin_phone":   b.NextOfKinPhone,
		"next_of_kin_relation": b.NextOfKinRelation,
		"credit_score":        b.CreditScore,
		"risk_level":          b.RiskLevel,
		"created_at":          b.CreatedAt,
		"updated_at":          b.UpdatedAt,
	}
	if b.Branch.ID != 0 {
		resp["branch"] = gin.H{"id": b.Branch.ID, "name": b.Branch.Name}
	}
	if h.Storage != nil {
		if b.PhotoKey != "" {
			if url, err := h.Storage.GetSignedURL(context.Background(), b.PhotoKey, time.Hour); err == nil {
				resp["photo_url"] = url
			}
		}
		if b.IDDocumentKey != "" {
			if url, err := h.Storage.GetSignedURL(context.Background(), b.IDDocumentKey, time.Hour); err == nil {
				resp["id_document_url"] = url
			}
		}
	}
	return resp
}

func (h *BorrowerHandler) List(c *gin.Context) {
	businessID := c.GetUint("business_id")
	q := h.DB.Preload("Branch").Where("business_id = ?", businessID)

	if branchID := c.Query("branch_id"); branchID != "" {
		q = q.Where("branch_id = ?", branchID)
	}
	if risk := c.Query("risk_level"); risk != "" {
		q = q.Where("risk_level = ?", risk)
	}
	if search := c.Query("search"); search != "" {
		like := "%" + strings.ToLower(search) + "%"
		q = q.Where("LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR phone LIKE ? OR national_id LIKE ?", like, like, "%"+search+"%", "%"+search+"%")
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
	q.Model(&models.Borrower{}).Count(&total)

	var rows []models.Borrower
	q.Order("created_at DESC").Offset((page - 1) * perPage).Limit(perPage).Find(&rows)

	out := make([]gin.H, 0, len(rows))
	for _, b := range rows {
		out = append(out, h.response(b))
	}
	c.JSON(http.StatusOK, gin.H{"data": out, "total": total, "page": page, "per_page": perPage})
}

func (h *BorrowerHandler) Get(c *gin.Context) {
	businessID := c.GetUint("business_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var b models.Borrower
	if err := h.DB.Preload("Branch").Where("id = ? AND business_id = ?", id, businessID).First(&b).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Borrower not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": h.response(b)})
}

func (h *BorrowerHandler) Create(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	if !canManageBorrowers(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req struct {
		BranchID          uint    `json:"branch_id" binding:"required"`
		FirstName         string  `json:"first_name" binding:"required"`
		LastName          string  `json:"last_name" binding:"required"`
		Phone             string  `json:"phone" binding:"required"`
		AltPhone          string  `json:"alt_phone"`
		Email             string  `json:"email"`
		DateOfBirth       *time.Time `json:"date_of_birth"`
		Gender            string  `json:"gender"`
		NationalID        string  `json:"national_id"`
		Address           string  `json:"address"`
		EmploymentStatus  string  `json:"employment_status"`
		Occupation        string  `json:"occupation"`
		Employer          string  `json:"employer"`
		MonthlyIncome     float64 `json:"monthly_income"`
		NextOfKinName     string  `json:"next_of_kin_name"`
		NextOfKinPhone    string  `json:"next_of_kin_phone"`
		NextOfKinRelation string  `json:"next_of_kin_relation"`
		RiskLevel         string  `json:"risk_level"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch, first_name, last_name and phone are required"})
		return
	}
	if !validateBranchInBusiness(h.DB, req.BranchID, businessID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch does not belong to this business"})
		return
	}
	if req.RiskLevel == "" {
		req.RiskLevel = models.RiskLevelMedium
	}

	b := models.Borrower{
		BusinessID:        businessID,
		BranchID:          req.BranchID,
		FirstName:         strings.TrimSpace(req.FirstName),
		LastName:          strings.TrimSpace(req.LastName),
		Phone:             strings.TrimSpace(req.Phone),
		AltPhone:          req.AltPhone,
		Email:             req.Email,
		DateOfBirth:       req.DateOfBirth,
		Gender:            req.Gender,
		NationalID:        strings.TrimSpace(req.NationalID),
		Address:           req.Address,
		EmploymentStatus:  req.EmploymentStatus,
		Occupation:        req.Occupation,
		Employer:          req.Employer,
		MonthlyIncome:     req.MonthlyIncome,
		NextOfKinName:     req.NextOfKinName,
		NextOfKinPhone:    req.NextOfKinPhone,
		NextOfKinRelation: req.NextOfKinRelation,
		RiskLevel:         req.RiskLevel,
		CreatedBy:         userID,
	}
	if err := h.DB.Create(&b).Error; err != nil {
		if isUniqueViolation(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "A borrower with that phone already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create borrower"})
		return
	}
	h.DB.Preload("Branch").First(&b, b.ID)
	h.Audit.Log(c, models.ActionCreate, "borrower", &b.ID, fmt.Sprintf("Added borrower %s (%s)", b.FullName(), b.Phone), gin.H{"branch_id": b.BranchID, "risk_level": b.RiskLevel})
	c.JSON(http.StatusCreated, gin.H{"data": h.response(b)})
}

func (h *BorrowerHandler) Update(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if !canManageBorrowers(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var b models.Borrower
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&b).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Borrower not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Whitelist of allowed update fields. Phone change is allowed but uniqueness is enforced by DB.
	allowed := map[string]bool{
		"branch_id": true, "first_name": true, "last_name": true, "phone": true, "alt_phone": true,
		"email": true, "date_of_birth": true, "gender": true, "national_id": true, "address": true,
		"employment_status": true, "occupation": true, "employer": true, "monthly_income": true,
		"next_of_kin_name": true, "next_of_kin_phone": true, "next_of_kin_relation": true,
		"credit_score": true, "risk_level": true,
	}
	updates := map[string]interface{}{}
	for k, v := range req {
		if allowed[k] {
			updates[k] = v
		}
	}
	if branchVal, ok := updates["branch_id"]; ok {
		if bid, ok := branchVal.(float64); ok {
			if !validateBranchInBusiness(h.DB, uint(bid), businessID) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Branch does not belong to this business"})
				return
			}
		}
	}
	if err := h.DB.Model(&b).Updates(updates).Error; err != nil {
		if isUniqueViolation(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "A borrower with that phone already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update borrower"})
		return
	}
	h.DB.Preload("Branch").First(&b, b.ID)
	h.Audit.Log(c, models.ActionUpdate, "borrower", &b.ID, fmt.Sprintf("Updated borrower %s", b.FullName()), nil)
	c.JSON(http.StatusOK, gin.H{"data": h.response(b)})
}

func (h *BorrowerHandler) Delete(c *gin.Context) {
	businessID := c.GetUint("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can delete borrowers"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var count int64
	h.DB.Model(&models.Loan{}).Where("borrower_id = ? AND status NOT IN ?", id, []string{models.LoanStatusCompleted, models.LoanStatusRejected, models.LoanStatusCancelled}).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete borrower with active loans"})
		return
	}
	var b models.Borrower
	_ = h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&b).Error
	h.DB.Where("id = ? AND business_id = ?", id, businessID).Delete(&models.Borrower{})
	bid := uint(id)
	h.Audit.Log(c, models.ActionDelete, "borrower", &bid, fmt.Sprintf("Deleted borrower %s", b.FullName()), nil)
	c.JSON(http.StatusOK, gin.H{"message": "Borrower deleted"})
}

// UploadDocument handles photo + ID document upload. Field "kind" decides which slot.
func (h *BorrowerHandler) UploadDocument(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if !canManageBorrowers(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	if h.Storage == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Storage not configured"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	kind := c.PostForm("kind") // "photo" or "id_document"
	if kind != "photo" && kind != "id_document" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kind must be 'photo' or 'id_document'"})
		return
	}

	var b models.Borrower
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&b).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Borrower not found"})
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	key := fmt.Sprintf("borrowers/%d/%d-%s-%s", businessID, b.ID, kind, header.Filename)
	if err := h.Storage.Upload(context.Background(), key, file, contentType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload"})
		return
	}

	if kind == "photo" {
		if b.PhotoKey != "" {
			_ = h.Storage.Delete(context.Background(), b.PhotoKey)
		}
		b.PhotoKey = key
	} else {
		if b.IDDocumentKey != "" {
			_ = h.Storage.Delete(context.Background(), b.IDDocumentKey)
		}
		b.IDDocumentKey = key
	}
	h.DB.Save(&b)
	url, _ := h.Storage.GetSignedURL(context.Background(), key, time.Hour)
	h.Audit.Log(c, models.ActionUpdate, "borrower", &b.ID, fmt.Sprintf("Uploaded %s for borrower %s", kind, b.FullName()), gin.H{"kind": kind})
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"url": url, "kind": kind}})
}

func canManageBorrowers(c *gin.Context) bool {
	role, _ := c.Get("user_role")
	return hasAnyRole(role, models.RoleAdmin, models.RoleManager, models.RoleLoanOfficer)
}
