// Package services > audit.go
//
// Audit logging — fire-and-forget. Every mutating handler calls
// audit.Log(c, action, resource, &id, description) AFTER its main work
// succeeds. The write happens on a background goroutine so a slow audit
// table never blocks the user's response.
//
// We deliberately swallow audit errors and log them server-side instead
// of bubbling them up: a failed audit row should not break a successful
// sale. If you need stronger guarantees, switch to a synchronous Save
// inside the handler's transaction.
package services

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"gritdemo/internal/models"
)

// AuditService writes activity rows. One instance per app, injected into
// route registration so handlers can pull it from gin.Context.
type AuditService struct {
	DB *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{DB: db}
}

// Log records an activity. Reads business_id, user_id, ip, ua from the
// gin.Context. Safe to call on any handler — if business_id or user_id
// aren't set (public route, system task), it skips silently.
//
// metadata may be nil. Pass any JSON-serializable value (map, struct).
func (s *AuditService) Log(
	c *gin.Context,
	action, resource string,
	resourceID *uint,
	description string,
	metadata any,
) {
	if s == nil || s.DB == nil {
		return
	}
	businessID := c.GetUint("business_id")
	userID := c.GetUint("user_id")
	if businessID == 0 || userID == 0 {
		return
	}

	var meta datatypes.JSON
	if metadata != nil {
		if b, err := json.Marshal(metadata); err == nil {
			meta = b
		}
	}

	row := models.Activity{
		BusinessID:  businessID,
		UserID:      userID,
		Action:      action,
		Resource:    resource,
		ResourceID:  resourceID,
		Description: description,
		Metadata:    meta,
		IPAddress:   c.ClientIP(),
		UserAgent:   trim(c.Request.UserAgent(), 511),
	}

	// Async write — handler returns immediately, audit lands eventually.
	// Capture db reference so closure doesn't outlive context-bound vars.
	db := s.DB
	go func() {
		if err := db.Create(&row).Error; err != nil {
			log.Printf("audit log failed (action=%s resource=%s): %v", action, resource, err)
		}
	}()
}

// LogFor lets system code (login handler, webhooks) record activity for a
// known business+user without a populated gin.Context. ipAddress / userAgent
// are still read from c when available.
func (s *AuditService) LogFor(
	c *gin.Context,
	businessID, userID uint,
	action, resource string,
	resourceID *uint,
	description string,
	metadata any,
) {
	if s == nil || s.DB == nil || businessID == 0 || userID == 0 {
		return
	}

	var meta datatypes.JSON
	if metadata != nil {
		if b, err := json.Marshal(metadata); err == nil {
			meta = b
		}
	}

	ip, ua := "", ""
	if c != nil && c.Request != nil {
		ip = c.ClientIP()
		ua = trim(c.Request.UserAgent(), 511)
	}

	row := models.Activity{
		BusinessID:  businessID,
		UserID:      userID,
		Action:      action,
		Resource:    resource,
		ResourceID:  resourceID,
		Description: description,
		Metadata:    meta,
		IPAddress:   ip,
		UserAgent:   ua,
	}

	db := s.DB
	go func() {
		if err := db.Create(&row).Error; err != nil {
			log.Printf("audit log failed (action=%s resource=%s): %v", action, resource, err)
		}
	}()
}

// UintPtr is a tiny helper so callsites can write `audit.UintPtr(s.ID)` instead of
// pulling ResourceID into a temp variable.
func UintPtr(v uint) *uint { return &v }

func trim(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}
