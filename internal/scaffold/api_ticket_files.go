package scaffold

// v3.30 — internal user-ticket system.
//
//   internal/models/ticket.go        — Ticket + TicketReply schemas
//   internal/handlers/ticket.go      — CRUD + reply + close/reopen + assign
//   internal/services/ticket_mail.go — Resend email-out on create
//
// Auth model:
//   - Any authenticated user can create tickets and reply on their own.
//   - Regular users see only their own; ADMIN sees the full queue and
//     can assign, close, reopen, reply with is_admin_reply=true.

// ticketModelGo emits internal/models/ticket.go.
func ticketModelGo() string {
	return `package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Ticket — single support request. Anyone authenticated can open one;
// ADMIN/EDITOR roles see the full queue, regular USER role sees their
// own. Status is open by default; closing stamps ClosedAt.
type Ticket struct {
	ID          string         ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	UserID      string         ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `
	Subject     string         ` + "`" + `gorm:"size:200;not null" json:"subject"` + "`" + `
	Description string         ` + "`" + `gorm:"type:text;not null" json:"description"` + "`" + `
	Status      string         ` + "`" + `gorm:"size:16;index;default:'open'" json:"status"` + "`" + `   // open | closed
	Priority    string         ` + "`" + `gorm:"size:16;index;default:'medium'" json:"priority"` + "`" + ` // low | medium | high | critical
	Labels      string         ` + "`" + `gorm:"size:255" json:"labels"` + "`" + `                           // comma-separated up to 8
	AssigneeID  string         ` + "`" + `gorm:"size:36;index" json:"assignee_id"` + "`" + `                // optional, must be ADMIN role
	LastReplyAt *time.Time     ` + "`" + `json:"last_reply_at"` + "`" + `                                    // touched by every reply
	ClosedAt    *time.Time     ` + "`" + `json:"closed_at"` + "`" + `
	CreatedAt   time.Time      ` + "`" + `gorm:"index" json:"created_at"` + "`" + `
	UpdatedAt   time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt   gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `

	// Eager-loaded relations for the detail page.
	User     *User           ` + "`" + `gorm:"foreignKey:UserID" json:"user,omitempty"` + "`" + `
	Assignee *User           ` + "`" + `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"` + "`" + `
	Replies  []TicketReply   ` + "`" + `gorm:"foreignKey:TicketID" json:"replies,omitempty"` + "`" + `
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	if t.Status == "" {
		t.Status = "open"
	}
	if t.Priority == "" {
		t.Priority = "medium"
	}
	return nil
}

// TicketReply — chronological thread under a ticket. IsAdminReply lets
// the dashboard style admin replies differently (badged, opposite side
// of the thread, etc.).
type TicketReply struct {
	ID            string         ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	TicketID      string         ` + "`" + `gorm:"size:36;index;not null" json:"ticket_id"` + "`" + `
	UserID        string         ` + "`" + `gorm:"size:36;index;not null" json:"user_id"` + "`" + `
	Body          string         ` + "`" + `gorm:"type:text;not null" json:"body"` + "`" + `
	IsAdminReply  bool           ` + "`" + `json:"is_admin_reply"` + "`" + `
	CreatedAt     time.Time      ` + "`" + `gorm:"index" json:"created_at"` + "`" + `
	UpdatedAt     time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt     gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `

	User *User ` + "`" + `gorm:"foreignKey:UserID" json:"user,omitempty"` + "`" + `
}

func (r *TicketReply) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}
`
}

// ticketHandlerGo emits internal/handlers/ticket_handler.go.
func ticketHandlerGo() string {
	return `package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/mail"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/paginate"
	"{{MODULE}}/internal/services"
)

type TicketHandler struct {
	DB   *gorm.DB
	Mail *mail.Mailer // can be nil — handler logs instead of emailing
}

type createTicketRequest struct {
	Subject     string ` + "`" + `json:"subject" binding:"required,min=3,max=200"` + "`" + `
	Description string ` + "`" + `json:"description" binding:"required,min=10"` + "`" + `
	Priority    string ` + "`" + `json:"priority"` + "`" + `
	Labels      string ` + "`" + `json:"labels"` + "`" + `
}

type replyRequest struct {
	Body string ` + "`" + `json:"body" binding:"required,min=1"` + "`" + `
}

type assignRequest struct {
	AssigneeID string ` + "`" + `json:"assignee_id" binding:"required"` + "`" + `
}

// Create opens a ticket for the authenticated user. Fires an email to
// SUPPORT_EMAIL when Resend is configured + a Notification for every
// ADMIN so the bell lights up.
//
//	POST /api/tickets
func (h *TicketHandler) Create(c *gin.Context) {
	var req createTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	userIDv, _ := c.Get("user_id")
	userID, _ := userIDv.(string)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "UNAUTHORIZED", "message": "sign in to open a ticket"},
		})
		return
	}

	// Cap labels at 8 — protects against accidental spam pasted into the
	// field. Trim each one so "bug , billing" doesn't store the space.
	labels := normalizeLabels(req.Labels, 8)

	ticket := models.Ticket{
		UserID:      userID,
		Subject:     req.Subject,
		Description: req.Description,
		Priority:    req.Priority,
		Labels:      labels,
	}
	if err := h.DB.Create(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "DB_ERROR", "message": err.Error()},
		})
		return
	}

	// Hydrate the user for the email + audit. Best-effort.
	var creator models.User
	h.DB.First(&creator, "id = ?", userID)

	// Fire-and-forget email + admin notifications. We don't fail the
	// request if either side trips — the ticket itself is persisted.
	go h.emitTicketCreated(&ticket, &creator)

	services.LogActivity(h.DB, c, services.ActivityArgs{
		Action:       "ticket.create",
		Severity:     "info",
		Summary:      fmt.Sprintf("Opened ticket %q (priority %s)", ticket.Subject, ticket.Priority),
		ResourceType: "ticket",
		ResourceID:   ticket.ID,
	})

	c.JSON(http.StatusCreated, gin.H{
		"data":    ticket,
		"message": "Ticket opened",
	})
}

// List returns tickets the caller can see. Regular users see their own;
// ADMIN/EDITOR see everything. Supports status (open|closed) + q filters.
//
//	GET /api/tickets?status=open&q=billing
func (h *TicketHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	isAdmin := role == "ADMIN" || role == "EDITOR"

	q := h.DB.Model(&models.Ticket{}).
		Preload("User").Preload("Assignee").
		Order("COALESCE(last_reply_at, created_at) DESC")

	if !isAdmin {
		q = q.Where("user_id = ?", userID)
	}
	params := paginate.Bind(c).
		With("status", c.Query("status")).
		With("priority", c.Query("priority")).
		With("assignee_id", c.Query("assignee_id"))
	if needle := c.Query("q"); needle != "" {
		q = q.Where("subject LIKE ? OR description LIKE ?", "%"+needle+"%", "%"+needle+"%")
	}

	res, err := paginate.List[models.Ticket](q, params, paginate.Config{
		Sortable:     map[string]bool{"created_at": true, "priority": true, "status": true},
		DefaultSort:  "created_at",
		DefaultOrder: "desc",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

// Get returns one ticket with replies. Same auth rule as List.
//
//	GET /api/tickets/:id
func (h *TicketHandler) Get(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	isAdmin := role == "ADMIN" || role == "EDITOR"

	var t models.Ticket
	q := h.DB.Preload("User").Preload("Assignee").
		Preload("Replies", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		Preload("Replies.User")
	if err := q.First(&t, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "ticket not found"},
		})
		return
	}
	if !isAdmin && t.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{"code": "FORBIDDEN", "message": "not your ticket"},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": t})
}

// Reply adds a message to the thread. Sets is_admin_reply when the
// caller is ADMIN/EDITOR so the UI can style staff replies.
//
//	POST /api/tickets/:id/reply
func (h *TicketHandler) Reply(c *gin.Context) {
	id := c.Param("id")
	var req replyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	userIDv, _ := c.Get("user_id")
	userID, _ := userIDv.(string)
	role, _ := c.Get("user_role")
	isAdmin := role == "ADMIN" || role == "EDITOR"

	var t models.Ticket
	if err := h.DB.First(&t, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "ticket not found"},
		})
		return
	}
	if !isAdmin && t.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{"code": "FORBIDDEN", "message": "not your ticket"},
		})
		return
	}

	reply := models.TicketReply{
		TicketID:     t.ID,
		UserID:       userID,
		Body:         req.Body,
		IsAdminReply: isAdmin,
	}
	now := time.Now()
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&reply).Error; err != nil {
			return err
		}
		return tx.Model(&t).Update("last_reply_at", now).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "DB_ERROR", "message": err.Error()},
		})
		return
	}

	services.LogActivity(h.DB, c, services.ActivityArgs{
		Action:       "ticket.reply",
		Severity:     "info",
		Summary:      fmt.Sprintf("Replied on ticket %q", t.Subject),
		ResourceType: "ticket",
		ResourceID:   t.ID,
	})

	c.JSON(http.StatusCreated, gin.H{"data": reply, "message": "Reply added"})
}

// Close stamps ClosedAt + status. Only the owner or an admin can close.
//
//	PATCH /api/tickets/:id/close
func (h *TicketHandler) Close(c *gin.Context) {
	h.transitionStatus(c, "closed")
}

// Reopen flips status back to open + clears ClosedAt.
//
//	PATCH /api/tickets/:id/reopen
func (h *TicketHandler) Reopen(c *gin.Context) {
	h.transitionStatus(c, "open")
}

// Assign points the ticket at an admin. Admins only.
//
//	PATCH /api/tickets/:id/assign
func (h *TicketHandler) Assign(c *gin.Context) {
	role, _ := c.Get("user_role")
	if role != "ADMIN" && role != "EDITOR" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{"code": "FORBIDDEN", "message": "admins only"},
		})
		return
	}
	var req assignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
		return
	}

	id := c.Param("id")
	var t models.Ticket
	if err := h.DB.First(&t, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "ticket not found"},
		})
		return
	}
	if err := h.DB.Model(&t).Update("assignee_id", req.AssigneeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "DB_ERROR", "message": err.Error()},
		})
		return
	}

	services.LogActivity(h.DB, c, services.ActivityArgs{
		Action:       "ticket.assign",
		Severity:     "info",
		Summary:      fmt.Sprintf("Assigned ticket %q", t.Subject),
		ResourceType: "ticket",
		ResourceID:   t.ID,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Assignee updated"})
}

func (h *TicketHandler) transitionStatus(c *gin.Context, status string) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	isAdmin := role == "ADMIN" || role == "EDITOR"

	var t models.Ticket
	if err := h.DB.First(&t, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "ticket not found"},
		})
		return
	}
	if !isAdmin && t.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{"code": "FORBIDDEN", "message": "not your ticket"},
		})
		return
	}

	updates := map[string]interface{}{"status": status}
	if status == "closed" {
		now := time.Now()
		updates["closed_at"] = &now
	} else {
		updates["closed_at"] = nil
	}
	if err := h.DB.Model(&t).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "DB_ERROR", "message": err.Error()},
		})
		return
	}

	services.LogActivity(h.DB, c, services.ActivityArgs{
		Action:       "ticket." + status,
		Severity:     "info",
		Summary:      fmt.Sprintf("Marked ticket %q as %s", t.Subject, status),
		ResourceType: "ticket",
		ResourceID:   t.ID,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

// emitTicketCreated is the fire-and-forget side-effect of opening a
// ticket: email SUPPORT_EMAIL via Resend (if configured), and create
// an in-app Notification for every ADMIN.
func (h *TicketHandler) emitTicketCreated(t *models.Ticket, creator *models.User) {
	// Email to SUPPORT_EMAIL — best-effort, never blocks the request.
	if h.Mail != nil {
		_ = services.SendTicketCreatedEmail(h.Mail, t, creator)
	}

	// Fan-out an admin notification per ADMIN — bell lights up.
	var admins []models.User
	h.DB.Where("role = ? AND active = ?", "ADMIN", true).Find(&admins)
	for _, a := range admins {
		n := models.Notification{
			UserID:   a.ID,
			Source:   "system",
			Severity: ticketSeverity(t.Priority),
			Title:    "New ticket: " + t.Subject,
			Body:     "Opened by " + creator.Email + ".",
			Link:     "/system/support/" + t.ID,
			Dedup:    "ticket-created:" + t.ID + ":" + a.ID,
		}
		// Ignore unique-index collisions — duplicate fires are no-ops.
		h.DB.FirstOrCreate(&n, models.Notification{Dedup: n.Dedup})
	}
}

func ticketSeverity(priority string) string {
	switch priority {
	case "critical":
		return "critical"
	case "high":
		return "high"
	case "low":
		return "low"
	default:
		return "medium"
	}
}

func normalizeLabels(raw string, cap int) string {
	if raw == "" {
		return ""
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
		if len(out) >= cap {
			break
		}
	}
	return strings.Join(out, ",")
}
`
}

// ticketMailGo emits internal/services/ticket_mail.go.
func ticketMailGo() string {
	return `package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"{{MODULE}}/internal/mail"
	"{{MODULE}}/internal/models"
)

// SendTicketCreatedEmail forwards a freshly-opened ticket to the support
// inbox configured via SUPPORT_EMAIL in .env. Silently no-ops (with a
// log line) when SUPPORT_EMAIL or Resend keys are missing — keeps the
// dev experience flowing without forcing email setup.
//
// The body intentionally stays plain-text + minimal HTML so any inbox
// renders it. The "Reply in dashboard" link points at the admin panel.
func SendTicketCreatedEmail(m *mail.Mailer, t *models.Ticket, creator *models.User) error {
	to := os.Getenv("SUPPORT_EMAIL")
	if to == "" {
		log.Printf("ticket-mail: SUPPORT_EMAIL not set, skipping email for ticket %s", t.ID)
		return nil
	}

	creatorLine := "unknown"
	if creator != nil {
		creatorLine = fmt.Sprintf("%s %s <%s>", creator.FirstName, creator.LastName, creator.Email)
	}

	subject := fmt.Sprintf("[Ticket #%s] %s", short(t.ID), t.Subject)
	dashURL := os.Getenv("ADMIN_URL")
	if dashURL == "" {
		dashURL = "http://localhost:3001"
	}

	html := fmt.Sprintf(` + "`" + `<!doctype html>
<html><body style="font-family: -apple-system, sans-serif; line-height: 1.55; color: #111;">
  <h2 style="margin: 0 0 12px 0;">New support ticket</h2>
  <p style="margin: 0 0 12px 0;"><strong>Subject:</strong> %s</p>
  <p style="margin: 0 0 12px 0;"><strong>Priority:</strong> %s &nbsp;|&nbsp; <strong>Labels:</strong> %s</p>
  <p style="margin: 0 0 12px 0;"><strong>From:</strong> %s</p>
  <hr style="border:none; border-top:1px solid #eee; margin: 16px 0;" />
  <pre style="white-space: pre-wrap; font-family: inherit; margin: 0 0 16px 0;">%s</pre>
  <p style="margin: 0;">
    <a href="%s/system/support/%s" style="display:inline-block; padding:10px 16px; background:#2563eb; color:white; text-decoration:none; border-radius:8px;">
      Reply in dashboard
    </a>
  </p>
</body></html>` + "`" + `,
		t.Subject, t.Priority, defaultIfEmpty(t.Labels, "—"),
		creatorLine, t.Description, dashURL, t.ID,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.SendRaw(ctx, to, subject, html)
}

func short(id string) string {
	if len(id) >= 8 {
		return id[:8]
	}
	return id
}

func defaultIfEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
`
}
