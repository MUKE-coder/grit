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

type LoanProductHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

// List returns all loan products for the active business.
func (h *LoanProductHandler) List(c *gin.Context) {
	businessID := c.GetUint("business_id")
	q := h.DB.Where("business_id = ?", businessID)
	if active := c.Query("active"); active == "true" {
		q = q.Where("is_active = ?", true)
	}
	var products []models.LoanProduct
	q.Order("name ASC").Find(&products)
	c.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *LoanProductHandler) Get(c *gin.Context) {
	businessID := c.GetUint("business_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var p models.LoanProduct
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *LoanProductHandler) Create(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if !canManageLoanProducts(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req struct {
		Name               string  `json:"name" binding:"required"`
		Description        string  `json:"description"`
		MinAmount          float64 `json:"min_amount" binding:"required"`
		MaxAmount          float64 `json:"max_amount" binding:"required"`
		MinDuration        int     `json:"min_duration" binding:"required"`
		MaxDuration        int     `json:"max_duration" binding:"required"`
		InterestMethod     string  `json:"interest_method" binding:"required"`
		InterestRate       float64 `json:"interest_rate"`
		RepaymentCycle     string  `json:"repayment_cycle" binding:"required"`
		RequiresCollateral bool    `json:"requires_collateral"`
		GracePeriodDays    int     `json:"grace_period_days"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required fields missing"})
		return
	}
	if req.MinAmount > req.MaxAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_amount cannot exceed max_amount"})
		return
	}
	if req.MinDuration > req.MaxDuration {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_duration cannot exceed max_duration"})
		return
	}
	if req.InterestMethod != models.InterestMethodFlat && req.InterestMethod != models.InterestMethodReducingBalance {
		c.JSON(http.StatusBadRequest, gin.H{"error": "interest_method must be 'flat' or 'reducing_balance'"})
		return
	}
	switch req.RepaymentCycle {
	case models.RepaymentCycleWeekly, models.RepaymentCycleBiweekly, models.RepaymentCycleMonthly:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "repayment_cycle must be 'weekly', 'biweekly', or 'monthly'"})
		return
	}

	p := models.LoanProduct{
		BusinessID:         businessID,
		Name:               strings.TrimSpace(req.Name),
		Description:        req.Description,
		MinAmount:          req.MinAmount,
		MaxAmount:          req.MaxAmount,
		MinDuration:        req.MinDuration,
		MaxDuration:        req.MaxDuration,
		InterestMethod:     req.InterestMethod,
		InterestRate:       req.InterestRate,
		RepaymentCycle:     req.RepaymentCycle,
		RequiresCollateral: req.RequiresCollateral,
		GracePeriodDays:    req.GracePeriodDays,
		IsActive:           true,
	}
	if err := h.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create loan product"})
		return
	}
	h.Audit.Log(c, models.ActionCreate, "loan_product", &p.ID, fmt.Sprintf("Created loan product %q", p.Name), gin.H{"interest_method": p.InterestMethod, "interest_rate": p.InterestRate, "min_amount": p.MinAmount, "max_amount": p.MaxAmount})
	c.JSON(http.StatusCreated, gin.H{"data": p})
}

func (h *LoanProductHandler) Update(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if !canManageLoanProducts(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var p models.LoanProduct
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan product not found"})
		return
	}

	var req struct {
		Name               *string  `json:"name"`
		Description        *string  `json:"description"`
		MinAmount          *float64 `json:"min_amount"`
		MaxAmount          *float64 `json:"max_amount"`
		MinDuration        *int     `json:"min_duration"`
		MaxDuration        *int     `json:"max_duration"`
		InterestMethod     *string  `json:"interest_method"`
		InterestRate       *float64 `json:"interest_rate"`
		RepaymentCycle     *string  `json:"repayment_cycle"`
		RequiresCollateral *bool    `json:"requires_collateral"`
		GracePeriodDays    *int     `json:"grace_period_days"`
		IsActive           *bool    `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if req.Name != nil {
		p.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.MinAmount != nil {
		p.MinAmount = *req.MinAmount
	}
	if req.MaxAmount != nil {
		p.MaxAmount = *req.MaxAmount
	}
	if req.MinDuration != nil {
		p.MinDuration = *req.MinDuration
	}
	if req.MaxDuration != nil {
		p.MaxDuration = *req.MaxDuration
	}
	if req.InterestMethod != nil {
		p.InterestMethod = *req.InterestMethod
	}
	if req.InterestRate != nil {
		p.InterestRate = *req.InterestRate
	}
	if req.RepaymentCycle != nil {
		p.RepaymentCycle = *req.RepaymentCycle
	}
	if req.RequiresCollateral != nil {
		p.RequiresCollateral = *req.RequiresCollateral
	}
	if req.GracePeriodDays != nil {
		p.GracePeriodDays = *req.GracePeriodDays
	}
	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}

	h.DB.Save(&p)
	h.Audit.Log(c, models.ActionUpdate, "loan_product", &p.ID, fmt.Sprintf("Updated loan product %q", p.Name), nil)
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *LoanProductHandler) Delete(c *gin.Context) {
	businessID := c.GetUint("business_id")
	role, _ := c.Get("user_role")
	if role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can delete loan products"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	// Refuse if loans exist on this product.
	var count int64
	h.DB.Model(&models.Loan{}).Where("loan_product_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete a loan product with existing loans — deactivate it instead"})
		return
	}
	h.DB.Where("id = ? AND business_id = ?", id, businessID).Delete(&models.LoanProduct{})
	lpid := uint(id)
	h.Audit.Log(c, models.ActionDelete, "loan_product", &lpid, fmt.Sprintf("Deleted loan product #%d", lpid), nil)
	c.JSON(http.StatusOK, gin.H{"message": "Loan product deleted"})
}

func canManageLoanProducts(c *gin.Context) bool {
	role, _ := c.Get("user_role")
	return hasAnyRole(role, models.RoleAdmin, models.RoleManager)
}
