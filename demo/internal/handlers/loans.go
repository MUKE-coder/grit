package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/models"
	"gritdemo/internal/services"
)

type LoanHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

// List returns loans for the active business with optional status / borrower / branch filters.
func (h *LoanHandler) List(c *gin.Context) {
	businessID := c.GetUint("business_id")
	q := h.DB.Preload("Borrower").Preload("Motorcycle").Preload("LoanProduct").Preload("Branch").
		Where("business_id = ?", businessID)

	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if borrowerID := c.Query("borrower_id"); borrowerID != "" {
		q = q.Where("borrower_id = ?", borrowerID)
	}
	if branchID := c.Query("branch_id"); branchID != "" {
		q = q.Where("branch_id = ?", branchID)
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
	q.Model(&models.Loan{}).Count(&total)

	var loans []models.Loan
	q.Order("created_at DESC").Offset((page - 1) * perPage).Limit(perPage).Find(&loans)

	c.JSON(http.StatusOK, gin.H{"data": loans, "total": total, "page": page, "per_page": perPage})
}

// Get returns a single loan with its full schedule and repayment history.
func (h *LoanHandler) Get(c *gin.Context) {
	businessID := c.GetUint("business_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var loan models.Loan
	if err := h.DB.Preload("Borrower").Preload("Motorcycle").Preload("LoanProduct").Preload("Branch").
		Where("id = ? AND business_id = ?", id, businessID).First(&loan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	var schedule []models.RepaymentSchedule
	h.DB.Where("loan_id = ?", loan.ID).Order("installment_number ASC").Find(&schedule)

	var repayments []models.Repayment
	h.DB.Preload("Collector").Where("loan_id = ?", loan.ID).Order("collection_date DESC").Find(&repayments)

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"loan":       loan,
		"schedule":   schedule,
		"repayments": repayments,
	}})
}

// Schedule returns just the repayment schedule for a loan.
func (h *LoanHandler) Schedule(c *gin.Context) {
	businessID := c.GetUint("business_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var loan models.Loan
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&loan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}
	var schedule []models.RepaymentSchedule
	h.DB.Where("loan_id = ?", loan.ID).Order("installment_number ASC").Find(&schedule)
	c.JSON(http.StatusOK, gin.H{"data": schedule})
}

// Create drafts a new loan (status=pending) and generates its repayment schedule.
func (h *LoanHandler) Create(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	if !canManageLoans(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req struct {
		BranchID         uint      `json:"branch_id" binding:"required"`
		BorrowerID       uint      `json:"borrower_id" binding:"required"`
		MotorcycleID     *uint     `json:"motorcycle_id"`
		LoanProductID    uint      `json:"loan_product_id" binding:"required"`
		PrincipalAmount  float64   `json:"principal_amount" binding:"required"`
		InitialDeposit   float64   `json:"initial_deposit"`
		Duration         int       `json:"duration"`
		InterestRate     float64   `json:"interest_rate"`
		InterestMethod   string    `json:"interest_method"`
		RepaymentCycle   string    `json:"repayment_cycle"`
		FirstPaymentDate time.Time `json:"first_payment_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required fields missing or invalid"})
		return
	}
	if !validateBranchInBusiness(h.DB, req.BranchID, businessID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Branch does not belong to this business"})
		return
	}
	// Borrower must belong to the same business.
	var bcount int64
	h.DB.Model(&models.Borrower{}).Where("id = ? AND business_id = ?", req.BorrowerID, businessID).Count(&bcount)
	if bcount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Borrower does not belong to this business"})
		return
	}
	// If a motorcycle was specified, verify ownership and availability.
	if req.MotorcycleID != nil {
		var moto models.Motorcycle
		if err := h.DB.Where("id = ? AND business_id = ?", *req.MotorcycleID, businessID).First(&moto).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Motorcycle does not belong to this business"})
			return
		}
		if moto.Status != models.MotorcycleStatusAvailable {
			c.JSON(http.StatusConflict, gin.H{"error": "Motorcycle is not available"})
			return
		}
	}

	loan, err := services.CreateLoan(h.DB, services.CreateLoanInput{
		BusinessID:       businessID,
		BranchID:         req.BranchID,
		BorrowerID:       req.BorrowerID,
		MotorcycleID:     req.MotorcycleID,
		LoanProductID:    req.LoanProductID,
		PrincipalAmount:  req.PrincipalAmount,
		InitialDeposit:   req.InitialDeposit,
		Duration:         req.Duration,
		InterestRate:     req.InterestRate,
		InterestMethod:   req.InterestMethod,
		RepaymentCycle:   req.RepaymentCycle,
		FirstPaymentDate: req.FirstPaymentDate,
		CreatedBy:        userID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.Audit.Log(c, models.ActionCreate, "loan", &loan.ID,
		fmt.Sprintf("Created loan %s for UGX %s", loan.LoanNumber, formatMoney(loan.PrincipalAmount)),
		gin.H{"borrower_id": loan.BorrowerID, "principal": loan.PrincipalAmount, "loan_product_id": loan.LoanProductID},
	)

	c.JSON(http.StatusCreated, gin.H{"data": loan})
}

// Approve flips a pending loan to approved.
func (h *LoanHandler) Approve(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	if !hasAnyRole(role, models.RoleAdmin, models.RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins and managers can approve loans"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var loan models.Loan
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&loan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}
	if err := services.ApproveLoan(h.DB, loan.ID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.DB.First(&loan, loan.ID)
	h.Audit.Log(c, models.ActionApprove, "loan", &loan.ID, fmt.Sprintf("Approved loan %s", loan.LoanNumber), nil)
	c.JSON(http.StatusOK, gin.H{"data": loan})
}

// Reject flips a pending loan to rejected.
func (h *LoanHandler) Reject(c *gin.Context) {
	businessID := c.GetUint("business_id")
	role, _ := c.Get("user_role")
	if !hasAnyRole(role, models.RoleAdmin, models.RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins and managers can reject loans"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var loan models.Loan
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&loan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}
	if loan.Status != models.LoanStatusPending {
		c.JSON(http.StatusConflict, gin.H{"error": "Only pending loans can be rejected"})
		return
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&loan).Update("status", models.LoanStatusRejected).Error; err != nil {
			return err
		}
		// Free the motorcycle so it's sellable again.
		if loan.MotorcycleID != nil {
			return tx.Model(&models.Motorcycle{}).Where("id = ?", *loan.MotorcycleID).
				Update("status", models.MotorcycleStatusAvailable).Error
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.Audit.Log(c, models.ActionReject, "loan", &loan.ID, fmt.Sprintf("Rejected loan %s", loan.LoanNumber), nil)
	c.JSON(http.StatusOK, gin.H{"message": "Loan rejected"})
}

// Disburse flips an approved loan to active.
func (h *LoanHandler) Disburse(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	if !hasAnyRole(role, models.RoleAdmin, models.RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins and managers can disburse loans"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var loan models.Loan
	if err := h.DB.Where("id = ? AND business_id = ?", id, businessID).First(&loan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}
	if err := services.DisburseLoan(h.DB, loan.ID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.DB.First(&loan, loan.ID)
	h.Audit.Log(c, models.ActionDisburse, "loan", &loan.ID,
		fmt.Sprintf("Disbursed loan %s — UGX %s", loan.LoanNumber, formatMoney(loan.PrincipalAmount)), nil)
	c.JSON(http.StatusOK, gin.H{"data": loan})
}

func canManageLoans(c *gin.Context) bool {
	role, _ := c.Get("user_role")
	return hasAnyRole(role, models.RoleAdmin, models.RoleManager, models.RoleLoanOfficer)
}
