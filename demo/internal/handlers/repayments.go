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

type RepaymentHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

// List returns repayments scoped to the business via the parent loan.
func (h *RepaymentHandler) List(c *gin.Context) {
	businessID := c.GetUint("business_id")
	q := h.DB.Preload("Loan").Preload("Collector").
		Joins("JOIN loans ON loans.id = repayments.loan_id").
		Where("loans.business_id = ?", businessID)

	if status := c.Query("status"); status != "" {
		q = q.Where("repayments.status = ?", status)
	}
	if loanID := c.Query("loan_id"); loanID != "" {
		q = q.Where("repayments.loan_id = ?", loanID)
	}
	if method := c.Query("payment_method"); method != "" {
		q = q.Where("repayments.payment_method = ?", method)
	}
	if from := c.Query("from_date"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			q = q.Where("repayments.collection_date >= ?", t)
		}
	}
	if to := c.Query("to_date"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			q = q.Where("repayments.collection_date < ?", t.AddDate(0, 0, 1))
		}
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
	q.Model(&models.Repayment{}).Count(&total)

	var rows []models.Repayment
	q.Order("repayments.collection_date DESC").Offset((page - 1) * perPage).Limit(perPage).Find(&rows)

	c.JSON(http.StatusOK, gin.H{"data": rows, "total": total, "page": page, "per_page": perPage})
}

func (h *RepaymentHandler) Get(c *gin.Context) {
	businessID := c.GetUint("business_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var r models.Repayment
	err := h.DB.Preload("Loan").Preload("Schedule").Preload("Collector").
		Joins("JOIN loans ON loans.id = repayments.loan_id").
		Where("repayments.id = ? AND loans.business_id = ?", id, businessID).
		First(&r).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repayment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

// Create records a new repayment. Cash-method repayments start as PENDING and
// must be approved by a manager. Mobile-money repayments are initiated via
// the dedicated DGateway endpoint, not here.
func (h *RepaymentHandler) Create(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	if !canCollectRepayments(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req struct {
		LoanID         uint    `json:"loan_id" binding:"required"`
		ScheduleID     *uint   `json:"schedule_id"`
		Amount         float64 `json:"amount" binding:"required"`
		PaymentMethod  string  `json:"payment_method" binding:"required"`
		TransactionRef string  `json:"transaction_ref"`
		CollectionDate *time.Time `json:"collection_date"`
		Notes          string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Loan, amount and payment_method are required"})
		return
	}
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
		return
	}
	switch req.PaymentMethod {
	case models.RepaymentMethodCash, models.RepaymentMethodBank, models.RepaymentMethodCheque:
	case models.RepaymentMethodMobileMoney:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Use the DGateway endpoint to collect mobile-money repayments"})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported payment_method"})
		return
	}

	// Validate the loan belongs to this business and is collectable.
	var loan models.Loan
	if err := h.DB.Where("id = ? AND business_id = ?", req.LoanID, businessID).First(&loan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Loan does not belong to this business"})
		return
	}
	if loan.Status != models.LoanStatusActive {
		c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Cannot collect on a %s loan", loan.Status)})
		return
	}

	collectionDate := time.Now()
	if req.CollectionDate != nil {
		collectionDate = *req.CollectionDate
	}

	r := models.Repayment{
		LoanID:         req.LoanID,
		ScheduleID:     req.ScheduleID,
		Amount:         req.Amount,
		CollectionDate: collectionDate,
		PaymentMethod:  req.PaymentMethod,
		TransactionRef: req.TransactionRef,
		Status:         models.RepaymentStatusPending,
		Notes:          req.Notes,
		CollectedBy:    userID,
	}
	r.Receipt = fmt.Sprintf("KM-RC-%d-%d", time.Now().Unix(), userID)

	if err := h.DB.Create(&r).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record repayment"})
		return
	}

	h.Audit.Log(c, models.ActionCreate, "repayment", &r.ID,
		fmt.Sprintf("Recorded UGX %s repayment on %s", formatMoney(r.Amount), loan.LoanNumber),
		gin.H{"loan_id": r.LoanID, "amount": r.Amount, "method": r.PaymentMethod},
	)

	c.JSON(http.StatusCreated, gin.H{"data": r})
}

// Approve verifies a pending repayment and applies it against the loan balance.
func (h *RepaymentHandler) Approve(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	if !hasAnyRole(role, models.RoleAdmin, models.RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins and managers can approve repayments"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var r models.Repayment
		err := tx.Joins("JOIN loans ON loans.id = repayments.loan_id").
			Where("repayments.id = ? AND loans.business_id = ?", id, businessID).
			First(&r).Error
		if err != nil {
			return fmt.Errorf("not_found")
		}
		if r.Status != models.RepaymentStatusPending {
			return fmt.Errorf("only pending repayments can be approved (current: %s)", r.Status)
		}
		now := time.Now()
		if err := tx.Model(&r).Updates(map[string]interface{}{
			"status":      models.RepaymentStatusApproved,
			"verified_by": userID,
			"verified_at": now,
		}).Error; err != nil {
			return err
		}
		return services.ApplyRepayment(tx, r.LoanID, r.Amount, r.ScheduleID)
	})
	if err != nil {
		if err.Error() == "not_found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Repayment not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rid := uint(id)
	h.Audit.Log(c, models.ActionApprove, "repayment", &rid, fmt.Sprintf("Approved repayment #%d", rid), nil)
	c.JSON(http.StatusOK, gin.H{"message": "Repayment approved"})
}

// Reject flips a pending repayment to rejected. Does NOT affect loan balance
// (the money never landed in our books in the first place).
func (h *RepaymentHandler) Reject(c *gin.Context) {
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	if !hasAnyRole(role, models.RoleAdmin, models.RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins and managers can reject repayments"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)

	var r models.Repayment
	err := h.DB.Joins("JOIN loans ON loans.id = repayments.loan_id").
		Where("repayments.id = ? AND loans.business_id = ?", id, businessID).
		First(&r).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repayment not found"})
		return
	}
	if r.Status != models.RepaymentStatusPending {
		c.JSON(http.StatusConflict, gin.H{"error": "Only pending repayments can be rejected"})
		return
	}
	now := time.Now()
	notes := r.Notes
	if req.Reason != "" {
		notes = notes + "\n[Rejected] " + req.Reason
	}
	h.DB.Model(&r).Updates(map[string]interface{}{
		"status":      models.RepaymentStatusRejected,
		"verified_by": userID,
		"verified_at": now,
		"notes":       notes,
	})
	h.Audit.Log(c, models.ActionReject, "repayment", &r.ID, fmt.Sprintf("Rejected repayment #%d", r.ID), gin.H{"reason": req.Reason})
	c.JSON(http.StatusOK, gin.H{"message": "Repayment rejected"})
}

func canCollectRepayments(c *gin.Context) bool {
	role, _ := c.Get("user_role")
	return hasAnyRole(role, models.RoleAdmin, models.RoleManager, models.RoleLoanOfficer, models.RoleCashier)
}
