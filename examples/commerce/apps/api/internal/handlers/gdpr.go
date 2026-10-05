package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"commerce/apps/api/internal/authz"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/respond"
	"commerce/apps/api/internal/services"
)

// GDPRHandler serves the right-to-access and right-to-erasure endpoints.
type GDPRHandler struct {
	DB *gorm.DB
}

func NewGDPRHandler(db *gorm.DB) *GDPRHandler {
	return &GDPRHandler{DB: db}
}

// Export returns the full data bundle for a user as a downloadable JSON file.
// A user may export their own data; an admin may export anyone's.
func (h *GDPRHandler) Export(c *gin.Context) {
	targetID := c.Param("id")
	callerRole, _ := c.Get("user_role")
	if authz.CurrentUserID(c) != targetID && fmt.Sprint(callerRole) != models.RoleAdmin {
		respond.Fail(c, respond.CodeForbidden, "you may only export your own data")
		return
	}

	bundle, err := services.ExportUserData(h.DB.WithContext(c.Request.Context()), targetID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			respond.Fail(c, respond.CodeNotFound, "user not found")
			return
		}
		respond.ServerError(c, "EXPORT_FAILED", err, "Internal server error")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"user-%s-export.json\"", targetID))
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Status(http.StatusOK)
	// The activity log is written a page at a time, so the status is on the wire
	// before the last page is read. A failure after that is logged, and leaves a
	// file that does not parse rather than one that looks complete.
	if err := bundle.WriteJSON(c.Writer); err != nil {
		log.Printf("gdpr export: writing the bundle: %v", err)
	}
}

type EraseRequest struct {
	// Required, not optional. An erasure is irreversible and its journal row is
	// what an auditor reads a year later; an empty reason tells them nothing
	// about why a person's data was destroyed. The previous handler ignored the
	// bind error, so a request with no body at all, or with the field
	// misspelled, erased an account and recorded no reason for it.
	Reason string `json:"reason" binding:"required,min=3,max=500"`
}

// Erase fulfils a right-to-erasure request (admin only). It refuses to let an
// admin erase themselves — that would revoke their own access mid-request and
// orphan the operation.
func (h *GDPRHandler) Erase(c *gin.Context) {
	targetID := c.Param("id")
	callerID := authz.CurrentUserID(c)
	callerEmail, _ := c.Get("user_email")

	if callerID == targetID {
		respond.Fail(c, respond.CodeSelfErase, "you cannot erase your own account here")
		return
	}

	var req EraseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Validation(c, "A reason is required for an erasure", map[string]string{
			"reason": "Give the compliance reason for this erasure (3 to 500 characters). It is written to the deletion journal.",
		})
		return
	}

	journal, err := services.EraseUser(h.DB, targetID, callerID, fmt.Sprint(callerEmail), req.Reason)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			respond.Fail(c, respond.CodeNotFound, "user not found")
			return
		}
		respond.ServerError(c, "ERASE_FAILED", err, "Internal server error")
		return
	}

	// Record the erasure in the semantic activity log too, so it shows up in the
	// dashboard and flows out through the OCSF/SIEM export.
	//
	// LogActivityErr rather than a bare Create: the erasure itself has already
	// happened, so a lost audit row must not fail the request, but it is the
	// record an auditor asks for. The response says so rather than reporting a
	// clean erasure that left no trace.
	activityLogged := true
	if err := services.LogActivityErr(h.DB, c, services.ActivityArgs{
		UserID:       callerID,
		Action:       "user.gdpr_erase",
		Severity:     "warn",
		Summary:      fmt.Sprintf("Erased all personal data for user %s (%d records)", targetID, journal.RecordsAffected),
		ResourceType: "user",
		ResourceID:   targetID,
	}); err != nil {
		activityLogged = false
	}

	if !activityLogged {
		c.JSON(http.StatusOK, gin.H{
			"data":    journal,
			"message": "User data erased, but the activity row could not be written",
			"meta":    gin.H{"activity_logged": false},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": journal, "message": "User data erased"})
}

// Journal lists the deletion journal and replays its hash chain so the caller
// can see at a glance whether the record of erasures is intact (admin only).
func (h *GDPRHandler) Journal(c *gin.Context) {
	var rows []models.DeletionJournal
	if err := h.DB.WithContext(c.Request.Context()).Order("created_at desc, id desc").Find(&rows).Error; err != nil {
		respond.ServerError(c, "QUERY_FAILED", err, "Internal server error")
		return
	}
	verification, err := services.VerifyJournalChain(h.DB)
	if err != nil {
		respond.ServerError(c, "VERIFY_FAILED", err, "Internal server error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows, "meta": verification})
}
