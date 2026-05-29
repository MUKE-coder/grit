package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gritdemo/internal/dgateway"
	"gritdemo/internal/models"
	"gritdemo/internal/services"
)

// DGatewayHandler hosts mobile-money collection endpoints.
// nil DG client = the gateway is not configured; handlers return 503.
type DGatewayHandler struct {
	DB *gorm.DB
	DG *dgateway.Client

	DefaultProvider string // "iotec" or "relworx"
	WebhookCallbackURL string // Public URL where DGateway can POST webhooks (e.g. https://demo.gritframework.dev/api/dgateway/webhook)
}

// Collect creates a pending Repayment and asks DGateway to initiate a
// request-to-pay against the borrower's phone. The repayment stays in PENDING
// status until either the webhook arrives or a poll observes a final state.
func (h *DGatewayHandler) Collect(c *gin.Context) {
	if h.DG == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Mobile money gateway not configured"})
		return
	}

	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	if !canCollectRepayments(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req struct {
		LoanID     uint    `json:"loan_id" binding:"required"`
		ScheduleID *uint   `json:"schedule_id"`
		Amount     float64 `json:"amount" binding:"required"`
		Phone      string  `json:"phone" binding:"required"`
		Provider   string  `json:"provider"` // optional override
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loan_id, amount, and phone are required"})
		return
	}
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
		return
	}

	// Loan must belong to this business and be active.
	var loan models.Loan
	if err := h.DB.Where("id = ? AND business_id = ?", req.LoanID, businessID).First(&loan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Loan does not belong to this business"})
		return
	}
	if loan.Status != models.LoanStatusActive {
		c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Cannot collect on a %s loan", loan.Status)})
		return
	}

	provider := req.Provider
	if provider == "" {
		provider = h.DefaultProvider
	}

	// Create the Repayment row first so we have an internal ID to correlate.
	repayment := models.Repayment{
		LoanID:           req.LoanID,
		ScheduleID:       req.ScheduleID,
		Amount:           req.Amount,
		CollectionDate:   time.Now(),
		PaymentMethod:    models.RepaymentMethodMobileMoney,
		Status:           models.RepaymentStatusPending,
		DGatewayProvider: provider,
		DGatewayPhone:    req.Phone,
		CollectedBy:      userID,
	}
	repayment.Receipt = fmt.Sprintf("KM-RC-%d-%d", time.Now().Unix(), userID)
	if err := h.DB.Create(&repayment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create repayment"})
		return
	}

	// Initiate the collection. If DGateway fails, mark the repayment as failed
	// and surface the error — better than leaving a dangling pending row.
	collectReq := dgateway.CollectRequest{
		Amount:      req.Amount,
		Currency:    "UGX",
		PhoneNumber: req.Phone,
		Provider:    provider,
		Reference:   repayment.Receipt,
		Description: fmt.Sprintf("Loan %s repayment", loan.LoanNumber),
		CallbackURL: h.WebhookCallbackURL,
	}
	res, err := h.DG.Collect(c.Request.Context(), collectReq)
	if err != nil {
		now := time.Now()
		h.DB.Model(&repayment).Updates(map[string]interface{}{
			"status":               models.RepaymentStatusFailed,
			"dgateway_fail_reason": err.Error(),
			"verified_at":          now,
		})
		c.JSON(http.StatusBadGateway, gin.H{"error": "Mobile money request failed: " + err.Error(), "repayment_id": repayment.ID})
		return
	}

	// Stamp the gateway reference back on the repayment so polls/webhooks can find it.
	h.DB.Model(&repayment).Update("dgateway_reference", res.Reference)
	repayment.DGatewayReference = res.Reference

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"repayment_id":        repayment.ID,
			"dgateway_reference":  res.Reference,
			"status":              res.Status,
			"provider":            res.Provider,
			"amount":              res.Amount,
			"message":             "Approve the prompt on your phone to complete the payment.",
		},
	})
}

// Verify polls DGateway for the latest status of a repayment's transaction
// and applies the loan-balance update if it just completed. Idempotent — safe
// to call repeatedly while waiting on the borrower to approve the prompt.
func (h *DGatewayHandler) Verify(c *gin.Context) {
	if h.DG == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Mobile money gateway not configured"})
		return
	}
	businessID := c.GetUint("business_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var r models.Repayment
	err := h.DB.Joins("JOIN loans ON loans.id = repayments.loan_id").
		Where("repayments.id = ? AND loans.business_id = ?", id, businessID).
		First(&r).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repayment not found"})
		return
	}
	if r.DGatewayReference == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This repayment has no DGateway reference"})
		return
	}

	res, err := h.DG.Verify(c.Request.Context(), r.DGatewayReference)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Verify failed: " + err.Error()})
		return
	}

	if r.Status == models.RepaymentStatusPending {
		switch res.Status {
		case dgateway.StatusCompleted:
			if err := h.applyCompletedRepayment(r, res); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			r.Status = models.RepaymentStatusApproved
		case dgateway.StatusFailed:
			h.markRepaymentFailed(&r, res.FailureReason)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"repayment_id":  r.ID,
			"status":        r.Status,
			"gateway_status": res.Status,
			"failure_reason": res.FailureReason,
		},
	})
}

// Webhook is the public endpoint DGateway calls when a transaction settles.
// Critically: this route is NOT behind the internal API key (DGateway can't
// know our key) — instead it validates by re-checking the transaction with
// DGateway itself before applying any state changes. The webhook is just a
// "wake up" signal; we trust DGateway only via our outbound verify call.
func (h *DGatewayHandler) Webhook(c *gin.Context) {
	if h.DG == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	var event dgateway.WebhookEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	if event.Data.Reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing reference"})
		return
	}

	// Find the repayment by reference. If we don't own it, ignore (acknowledge so
	// DGateway doesn't keep retrying).
	var r models.Repayment
	if err := h.DB.Where("dgateway_reference = ?", event.Data.Reference).First(&r).Error; err != nil {
		log.Printf("dgateway webhook: unknown reference %s, ignoring", event.Data.Reference)
		c.Status(http.StatusOK)
		return
	}

	// Re-verify with DGateway to confirm the webhook isn't spoofed.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := h.DG.Verify(ctx, event.Data.Reference)
	if err != nil {
		log.Printf("dgateway webhook: verify failed for %s: %v", event.Data.Reference, err)
		c.Status(http.StatusServiceUnavailable)
		return
	}

	// Apply the verified state. If we already processed this reference, the
	// applyCompletedRepayment guard makes this idempotent.
	if r.Status == models.RepaymentStatusPending {
		switch res.Status {
		case dgateway.StatusCompleted:
			if err := h.applyCompletedRepayment(r, res); err != nil {
				log.Printf("dgateway webhook: apply failed for %s: %v", event.Data.Reference, err)
				c.Status(http.StatusInternalServerError)
				return
			}
		case dgateway.StatusFailed:
			h.markRepaymentFailed(&r, res.FailureReason)
		}
	}

	c.Status(http.StatusOK)
}

// applyCompletedRepayment marks a pending mobile-money repayment as approved
// and applies it to the loan balance, all in one transaction.
func (h *DGatewayHandler) applyCompletedRepayment(r models.Repayment, v *dgateway.VerifyResponse) error {
	return h.DB.Transaction(func(tx *gorm.DB) error {
		// Re-read inside the tx to guard against race with manual approval.
		var fresh models.Repayment
		if err := tx.First(&fresh, r.ID).Error; err != nil {
			return err
		}
		if fresh.Status != models.RepaymentStatusPending {
			return nil // already settled
		}

		now := time.Now()
		updates := map[string]interface{}{
			"status":                models.RepaymentStatusApproved,
			"verified_at":           now,
			"dgateway_confirmed_at": now,
			"dgateway_fee":          v.Fee,
			"dgateway_net_amount":   v.NetAmount,
		}
		if err := tx.Model(&fresh).Updates(updates).Error; err != nil {
			return err
		}
		return services.ApplyRepayment(tx, fresh.LoanID, fresh.Amount, fresh.ScheduleID)
	})
}

// markRepaymentFailed flips a pending mobile-money repayment to failed.
// Does NOT touch the loan balance — the money never moved.
func (h *DGatewayHandler) markRepaymentFailed(r *models.Repayment, reason string) {
	now := time.Now()
	if reason == "" {
		reason = "unknown"
	}
	h.DB.Model(r).Updates(map[string]interface{}{
		"status":               models.RepaymentStatusFailed,
		"dgateway_fail_reason": reason,
		"verified_at":          now,
	})
}

// Compile-time check we wired everything.
var _ = errors.New
