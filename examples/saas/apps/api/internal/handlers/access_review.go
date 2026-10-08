package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"saas/apps/api/internal/respond"
	"saas/apps/api/internal/services"
)

// AccessReviewHandler exposes the recertification workflow. Admin-only: the
// point of an access review is that a privileged reviewer certifies access, and
// its evidence must not be editable by the people whose access it covers.
type AccessReviewHandler struct {
	DB *gorm.DB
}

func NewAccessReviewHandler(db *gorm.DB) *AccessReviewHandler {
	return &AccessReviewHandler{DB: db}
}

// List returns campaigns newest first, each with its decision counts.
func (h *AccessReviewHandler) List(c *gin.Context) {
	// Two queries rather than four per campaign, and the counts' errors are no
	// longer dropped: see AccessReviewSummaries.
	out, err := services.AccessReviewSummaries(c.Request.Context(), h.DB)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "failed to load reviews")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// Get returns one campaign with all its items.
func (h *AccessReviewHandler) Get(c *gin.Context) {
	review, err := services.AccessReviewWithItems(c.Request.Context(), h.DB, c.Param("id"))
	if err != nil {
		respond.Fail(c, respond.CodeNotFound, "access review not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": review})
}

type OpenReviewRequest struct {
	Name string `json:"name" binding:"required"`
	Note string `json:"note"`
}

// Open starts a campaign, snapshotting all current role assignments.
func (h *AccessReviewHandler) Open(c *gin.Context) {
	var req OpenReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}
	review, err := services.OpenAccessReview(h.DB, req.Name, req.Note, c.GetString("user_id"), c.GetString("user_email"))
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "failed to open review")
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data":    review,
		"message": "Access review opened",
	})
}

type ReviewDecisionRequest struct {
	Decision string `json:"decision" binding:"required"`
	Note     string `json:"note"`
}

// Decide records an approve/revoke on one item. A revoke removes the grant.
func (h *AccessReviewHandler) Decide(c *gin.Context) {
	var req ReviewDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}
	item, err := services.DecideAccessReviewItem(
		h.DB, c, c.Param("id"), c.Param("itemId"),
		req.Decision, req.Note,
		c.GetString("user_id"), c.GetString("user_email"),
	)
	if err != nil {
		status, code := http.StatusBadRequest, "INVALID_DECISION"
		// errors.Is rather than a switch on err: a switch compares with ==,
		// so the day one of these is wrapped with %w every case stops matching
		// and the caller gets a 400 that says nothing.
		switch {
		case errors.Is(err, services.ErrReviewClosed):
			code = "REVIEW_CLOSED"
		case errors.Is(err, services.ErrItemDecided):
			code = "ITEM_LOCKED"
		case errors.Is(err, gorm.ErrRecordNotFound):
			status, code = http.StatusNotFound, "NOT_FOUND"
		}
		c.JSON(status, gin.H{"error": gin.H{"code": code, "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item, "message": "Decision recorded"})
}

// Complete signs off a campaign once every item is decided.
func (h *AccessReviewHandler) Complete(c *gin.Context) {
	review, err := services.CompleteAccessReview(h.DB, c.Param("id"), c.GetString("user_id"), c.GetString("user_email"))
	if err != nil {
		status, code := http.StatusBadRequest, "CANNOT_COMPLETE"
		switch {
		case errors.Is(err, services.ErrReviewIncomplete):
			code = "REVIEW_INCOMPLETE"
		case errors.Is(err, services.ErrReviewClosed):
			code = "REVIEW_CLOSED"
		case errors.Is(err, gorm.ErrRecordNotFound):
			status, code = http.StatusNotFound, "NOT_FOUND"
		}
		c.JSON(status, gin.H{"error": gin.H{"code": code, "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": review, "message": "Access review completed"})
}
