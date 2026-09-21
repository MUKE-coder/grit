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

	"{{MODULE}}/internal/ids"
	"gorm.io/gorm"
)
` + ticketConstantsBlock + `
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
	Labels      string         ` + "`" + `gorm:"size:255" json:"labels"` + "`" + `                           // comma-separated, capped by the ticket service
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
		t.ID = ids.New()
	}
` + ticketDefaultsNew + `	return nil
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
		r.ID = ids.New()
	}
	return nil
}
`
}

// ticketHandlerGo emits internal/handlers/ticket.go.
//
// Contact-app review M29: this file used to hold every query and every rule the
// ticket system has, and there was no ticket service at all. None of it could
// be called from a job, a seeder or a test without a *gin.Context. The handler
// binds, calls services.TicketService, and responds. It runs no query of its
// own, which is the pattern the generated resource handlers already follow.
//
// M32 and M33 live on in the service and here: a ticket the caller may not see
// is answered 404 by a scoped query, the last-activity order applies only when
// no sort was asked for, search is case-insensitive through paginate, and the
// roles are models.RoleAdmin and models.RoleEditor rather than literals.
func ticketHandlerGo() string {
	return `package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/jobs"
	"{{MODULE}}/internal/mail"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/paginate"
	"{{MODULE}}/internal/respond"
	"{{MODULE}}/internal/services"
)

type TicketHandler struct {
	DB   *gorm.DB
	Mail *mail.Mailer // can be nil — the new-ticket email is skipped
	// Jobs is optional: with it the new-ticket email goes on the background
	// queue, so a provider that is briefly down does not lose it.
	Jobs *jobs.Client
}

type CreateTicketRequest struct {
	Subject     string ` + "`" + `json:"subject" binding:"required,min=3,max=200"` + "`" + `
	Description string ` + "`" + `json:"description" binding:"required,min=10"` + "`" + `
	Priority    string ` + "`" + `json:"priority"` + "`" + `
	Labels      string ` + "`" + `json:"labels"` + "`" + `
}

type TicketReplyRequest struct {
	Body string ` + "`" + `json:"body" binding:"required,min=1"` + "`" + `
}

type AssignTicketRequest struct {
	AssigneeID string ` + "`" + `json:"assignee_id" binding:"required"` + "`" + `
}

// ticketListConfig is what the ticket list may be searched and sorted by.
var ticketListConfig = paginate.Config{
	Searchable:   []string{"subject", "description"},
	Sortable:     map[string]bool{"created_at": true, "priority": true, "status": true},
	DefaultSort:  "created_at",
	DefaultOrder: "desc",
}

// tickets is the service every method below calls. Built per request because it
// is three fields and a struct literal, and because building it here means the
// wiring in routes.go did not have to change.
func (h *TicketHandler) tickets() *services.TicketService {
	svc := &services.TicketService{DB: h.DB, Mail: h.Mail}
	// A typed nil pointer in an interface field is not a nil interface, so the
	// queue is only set when there really is one.
	if h.Jobs != nil {
		svc.Queue = h.Jobs
	}
	return svc
}

// ticketStaff reports whether the caller handles everyone's tickets.
//
// One definition, because there were five, each spelling ADMIN and EDITOR out
// as string literals. models.RoleAdmin and models.RoleEditor exist precisely so
// a project that renames a role renames it once.
func ticketStaff(c *gin.Context) bool {
	role, _ := c.Get("user_role")
	return role == models.RoleAdmin || role == models.RoleEditor
}

// actorOf reads who is acting from the request. With ticketStaff, the only
// place in the ticket code that knows what a gin context is.
func (h *TicketHandler) actorOf(c *gin.Context) services.TicketActor {
	return services.TicketActor{
		UserID: c.GetString("user_id"),
		Staff:  ticketStaff(c),
	}
}

// Create opens a ticket for the authenticated user. The service emails
// SUPPORT_EMAIL through the job queue and lights up every admin's bell.
//
//	POST /api/tickets
func (h *TicketHandler) Create(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}
	actor := h.actorOf(c)
	if actor.UserID == "" {
		respond.Fail(c, respond.CodeUnauthorized, "Sign in to open a ticket")
		return
	}

	ticket, err := h.tickets().Open(services.ContextOf(c), actor, services.NewTicket{
		Subject:     req.Subject,
		Description: req.Description,
		Priority:    req.Priority,
		Labels:      req.Labels,
	})
	if err != nil {
		respond.WriteError(c, err, "Could not open the ticket")
		return
	}

	services.LogActivityCtx(services.ContextOf(c), h.DB, services.ActivityArgs{
		Action:       "ticket.create",
		Severity:     "info",
		Summary:      fmt.Sprintf("Opened ticket %q (priority %s)", ticket.Subject, ticket.Priority),
		ResourceType: "ticket",
		ResourceID:   ticket.ID,
	})

	respond.Created(c, ticket, "Ticket opened")
}

// List returns tickets the caller can see. Regular users see their own;
// ADMIN/EDITOR see everything. Supports status (open|closed) + q filters.
//
//	GET /api/tickets?status=open&q=billing
func (h *TicketHandler) List(c *gin.Context) {
	q := h.tickets().Query(services.ContextOf(c), h.actorOf(c))

	params := paginate.Bind(c).
		With("status", c.Query("status")).
		With("priority", c.Query("priority")).
		With("assignee_id", c.Query("assignee_id"))

	// ?q= is this endpoint's spelling of ?search=, so hand it to paginate
	// rather than building the clause here. paginate compares with
	// LOWER(col) LIKE LOWER(?); a bare LIKE on Postgres means a search for
	// "Billing" finds nothing filed as "billing".
	if needle := c.Query("q"); needle != "" {
		params.Search = needle
	}

	// Newest activity first, but only when the caller asked for nothing.
	// Applied always, paginate's own ordering was appended after it, so
	// ?sort_by=priority never did more than break ties. The column the support
	// queue wants to sort by is not a column at all, which is why it cannot go
	// in Sortable.
	if !ticketListConfig.Sortable[params.SortBy] {
		q = q.Order("COALESCE(last_reply_at, created_at) DESC")
	}

	res, err := paginate.List[models.Ticket](q, params, ticketListConfig)
	if err != nil {
		respond.ServerError(c, "INTERNAL_ERROR", err, "Internal server error")
		return
	}
	c.JSON(http.StatusOK, res)
}

// Get returns one ticket with its replies. Same visibility rule as List.
//
//	GET /api/tickets/:id
func (h *TicketHandler) Get(c *gin.Context) {
	ticket, err := h.tickets().Visible(services.ContextOf(c), h.actorOf(c), c.Param("id"), true)
	if err != nil {
		respond.WriteError(c, err, "Could not load the ticket")
		return
	}
	respond.OK(c, ticket)
}

// Reply adds a message to the thread. Sets is_admin_reply when the
// caller is ADMIN/EDITOR so the UI can style staff replies.
//
//	POST /api/tickets/:id/reply
func (h *TicketHandler) Reply(c *gin.Context) {
	var req TicketReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}
	reply, ticket, err := h.tickets().Reply(services.ContextOf(c), h.actorOf(c), c.Param("id"), req.Body)
	if err != nil {
		respond.WriteError(c, err, "Could not add the reply")
		return
	}

	services.LogActivityCtx(services.ContextOf(c), h.DB, services.ActivityArgs{
		Action:       "ticket.reply",
		Severity:     "info",
		Summary:      fmt.Sprintf("Replied on ticket %q", ticket.Subject),
		ResourceType: "ticket",
		ResourceID:   ticket.ID,
	})

	respond.Created(c, reply, "Reply added")
}

// Close stamps ClosedAt + status. Only the owner or an admin can close.
//
//	PATCH /api/tickets/:id/close
func (h *TicketHandler) Close(c *gin.Context) {
	h.transitionStatus(c, models.TicketStatusClosed)
}

// Reopen flips status back to open + clears ClosedAt.
//
//	PATCH /api/tickets/:id/reopen
func (h *TicketHandler) Reopen(c *gin.Context) {
	h.transitionStatus(c, models.TicketStatusOpen)
}

// Assign points the ticket at an admin. Admins only.
//
//	PATCH /api/tickets/:id/assign
func (h *TicketHandler) Assign(c *gin.Context) {
	actor := h.actorOf(c)
	// Refused before the body is read, as it always was: a caller who may not
	// assign learns nothing from a validation message.
	if !actor.Staff {
		respond.WriteError(c, services.ErrTicketStaffOnly, "Admins only")
		return
	}
	var req AssignTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Fail(c, respond.CodeValidationError, err.Error())
		return
	}
	ticket, err := h.tickets().Assign(services.ContextOf(c), actor, c.Param("id"), req.AssigneeID)
	if err != nil {
		respond.WriteError(c, err, "Could not assign the ticket")
		return
	}

	services.LogActivityCtx(services.ContextOf(c), h.DB, services.ActivityArgs{
		Action:       "ticket.assign",
		Severity:     "info",
		Summary:      fmt.Sprintf("Assigned ticket %q", ticket.Subject),
		ResourceType: "ticket",
		ResourceID:   ticket.ID,
	})
	respond.OK(c, ticket, "Assignee updated")
}

func (h *TicketHandler) transitionStatus(c *gin.Context, status string) {
	ticket, err := h.tickets().SetStatus(services.ContextOf(c), h.actorOf(c), c.Param("id"), status)
	if err != nil {
		respond.WriteError(c, err, "Could not update the ticket")
		return
	}

	services.LogActivityCtx(services.ContextOf(c), h.DB, services.ActivityArgs{
		Action:       "ticket." + status,
		Severity:     "info",
		Summary:      fmt.Sprintf("Marked ticket %q as %s", ticket.Subject, status),
		ResourceType: "ticket",
		ResourceID:   ticket.ID,
	})
	respond.OK(c, ticket, "Status updated")
}
`
}

// ticketServiceGo emits internal/services/ticket.go.
func ticketServiceGo() string {
	return `package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/mail"
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/respond"
)

// Tickets, as a service rather than as a handler.
//
// Contact-app review M29: every query and every rule the ticket system has used
// to live in handlers/ticket.go, so none of it could be called from a job, a
// seeder or a test without a *gin.Context.
//
// Nothing here knows what an HTTP request is. The handler turns the request
// into a TicketActor, calls one method, and responds.

// TicketActor is who is acting, as much as the ticket rules need to know.
type TicketActor struct {
	UserID string
	// Staff see the whole queue and can assign; everyone else sees their own.
	Staff bool
}

// ticketError is a ticket rule the caller broke, carrying the code it comes
// back as. respond.WriteError answers it without the handler mapping anything.
type ticketError struct {
	message string
	code    respond.Code
}

func (e ticketError) Error() string           { return e.message }
func (e ticketError) ErrorCode() respond.Code { return e.code }

var (
	// ErrTicketNotFound is a ticket that does not exist, or one the actor may
	// not see. The two are deliberately the same answer: a 403 for somebody
	// else's ticket tells whoever is walking the id space which ids are real.
	ErrTicketNotFound error = ticketError{"Ticket not found", respond.CodeNotFound}
	// ErrTicketStaffOnly is an action only the support queue's owners can take.
	ErrTicketStaffOnly error = ticketError{"Admins only", respond.CodeForbidden}
)

// TicketService is everything opening, reading and answering a ticket does.
type TicketService struct {
	DB *gorm.DB
	// Mail is optional: without it the new-ticket email is skipped.
	Mail *mail.Mailer
	// Queue is optional. With it the new-ticket email goes to the background
	// worker, which retries a provider that is briefly down and survives a
	// restart; without it the send is inline. Leave it nil rather than
	// assigning a nil *jobs.Client: a typed nil in an interface is not nil.
	Queue mail.Enqueuer
}

// NewTicket is what a caller asks for when opening one.
type NewTicket struct {
	Subject     string
	Description string
	Priority    string
	Labels      string
}

// Open creates the ticket, notifies the admins and sends the support email.
func (s *TicketService) Open(ctx context.Context, actor TicketActor, in NewTicket) (*models.Ticket, error) {
	ticket := models.Ticket{
		UserID:      actor.UserID,
		Subject:     in.Subject,
		Description: in.Description,
		Priority:    in.Priority,
		Labels: NormalizeTicketLabels(in.Labels, MaxTicketLabels),
	}
	if err := s.DB.WithContext(ctx).Create(&ticket).Error; err != nil {
		return nil, fmt.Errorf("creating the ticket: %w", err)
	}

	// Hydrate the creator for the email and the notification. Best-effort: the
	// ticket is saved either way.
	var creator models.User
	if err := s.DB.WithContext(ctx).First(&creator, "id = ?", actor.UserID).Error; err != nil {
		log.Printf("tickets: loading the creator of %s: %v", ticket.ID, err)
	}

	s.announce(ctx, &ticket, &creator)
	return &ticket, nil
}

// announce emails the support inbox and lights up every admin's bell.
//
// It used to run in a goroutine the request started, which meant the email was
// never retried and was lost outright on a restart. The mail is queued now, and
// the notification rows are written before the response goes out: there are as
// many as there are admins, and that is a number a support queue can hold.
func (s *TicketService) announce(ctx context.Context, t *models.Ticket, creator *models.User) {
	switch {
	case s.Queue != nil:
		if err := QueueTicketCreatedEmail(ctx, s.Queue, t, creator); err != nil {
			log.Printf("tickets: queueing the email for %s: %v", t.ID, err)
		}
	case s.Mail != nil:
		// Not the request's context: the email carries its own timeout, so a
		// client that hangs up does not cancel the mail support is owed.
		if err := SendTicketCreatedEmail(s.Mail, t, creator); err != nil { //nolint:contextcheck // see above
			log.Printf("tickets: emailing support about %s: %v", t.ID, err)
		}
	}

	var admins []models.User
	if err := s.DB.WithContext(ctx).Where("role = ? AND active = ?", models.RoleAdmin, true).Find(&admins).Error; err != nil {
		log.Printf("tickets: listing admins to notify about %s: %v", t.ID, err)
		return
	}
	for _, a := range admins {
		n := models.Notification{
			UserID:   a.ID,
			Source:   "system",
			Severity: TicketSeverity(t.Priority),
			Title:    "New ticket: " + t.Subject,
			Body:     "Opened by " + creator.Email + ".",
			Link:     "/system/support/" + t.ID,
			Dedup:    "ticket-created:" + t.ID + ":" + a.ID,
		}
		// FirstOrCreate on the dedup key, so a duplicate fire is a no-op.
		if err := s.DB.WithContext(ctx).FirstOrCreate(&n, models.Notification{Dedup: n.Dedup}).Error; err != nil {
			log.Printf("tickets: notifying %s of ticket %s: %v", a.ID, t.ID, err)
		}
	}
}

// Query is the list query, scoped to what the actor may see. The caller orders,
// pages and filters it.
func (s *TicketService) Query(ctx context.Context, actor TicketActor) *gorm.DB {
	q := s.DB.WithContext(ctx).Model(&models.Ticket{}).Preload("User").Preload("Assignee")
	return s.scope(q, actor)
}

// scope is the visibility rule, in one place. Staff see the whole queue;
// everyone else sees the tickets they opened.
func (s *TicketService) scope(q *gorm.DB, actor TicketActor) *gorm.DB {
	if actor.Staff {
		return q
	}
	return q.Where("user_id = ?", actor.UserID)
}

// Visible loads a ticket the actor is allowed to see. withThread also loads the
// replies, oldest first, with their authors.
//
// The scope is part of the query, not a check after it, so a ticket the actor
// may not see comes back as ErrTicketNotFound exactly like one that never
// existed.
func (s *TicketService) Visible(ctx context.Context, actor TicketActor, id string, withThread bool) (*models.Ticket, error) {
	q := s.DB.WithContext(ctx)
	if withThread {
		q = q.Preload("User").Preload("Assignee").
			Preload("Replies", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
			Preload("Replies.User")
	}
	var t models.Ticket
	if err := s.scope(q, actor).First(&t, "id = ?", id).Error; err != nil {
		return nil, ErrTicketNotFound
	}
	return &t, nil
}

// Reply adds a message to the thread and touches the ticket's last reply.
func (s *TicketService) Reply(ctx context.Context, actor TicketActor, id, body string) (*models.TicketReply, *models.Ticket, error) {
	t, err := s.Visible(ctx, actor, id, false)
	if err != nil {
		return nil, nil, err
	}

	reply := models.TicketReply{
		TicketID:     t.ID,
		UserID:       actor.UserID,
		Body:         body,
		IsAdminReply: actor.Staff,
	}
	now := time.Now()
	if err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&reply).Error; err != nil {
			return err
		}
		return tx.Model(t).Update("last_reply_at", now).Error
	}); err != nil {
		return nil, nil, fmt.Errorf("saving the reply: %w", err)
	}
	t.LastReplyAt = &now
	return &reply, t, nil
}

// SetStatus closes or reopens a ticket. "closed" stamps ClosedAt.
func (s *TicketService) SetStatus(ctx context.Context, actor TicketActor, id, status string) (*models.Ticket, error) {
	t, err := s.Visible(ctx, actor, id, false)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{"status": status}
	if status == models.TicketStatusClosed {
		now := time.Now()
		updates["closed_at"] = &now
	} else {
		updates["closed_at"] = nil
	}
	if err := s.DB.WithContext(ctx).Model(t).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("updating the ticket status: %w", err)
	}
	return t, nil
}

// Assign points the ticket at somebody. Staff only.
func (s *TicketService) Assign(ctx context.Context, actor TicketActor, id, assigneeID string) (*models.Ticket, error) {
	if !actor.Staff {
		return nil, ErrTicketStaffOnly
	}
	t, err := s.Visible(ctx, actor, id, false)
	if err != nil {
		return nil, err
	}
	if err := s.DB.WithContext(ctx).Model(t).Update("assignee_id", assigneeID).Error; err != nil {
		return nil, fmt.Errorf("assigning the ticket: %w", err)
	}
	t.AssigneeID = assigneeID
	return t, nil
}

// TicketSeverity maps a ticket priority onto a notification severity.
func TicketSeverity(priority string) string {
	switch priority {
	case models.TicketPriorityCritical:
		return "critical"
	case models.TicketPriorityHigh:
		return "high"
	case models.TicketPriorityLow:
		return "low"
	default:
		return "medium"
	}
}

// MaxTicketLabels is how many labels a ticket keeps. It stops a paste into the
// field becoming 400 labels.
const MaxTicketLabels = 8

// NormalizeTicketLabels trims each label and keeps at most limit of them.
func NormalizeTicketLabels(raw string, limit int) string {
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
		if len(out) >= limit {
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
	"{{MODULE}}/internal/settings"
)

// SendTicketCreatedEmail forwards a freshly-opened ticket to the support
// inbox configured via SUPPORT_EMAIL in .env. Silently no-ops (with a
// log line) when SUPPORT_EMAIL or Resend keys are missing — keeps the
// dev experience flowing without forcing email setup.
//
// The body intentionally stays plain-text + minimal HTML so any inbox
// renders it. The "Reply in dashboard" link points at the admin panel.
` + ticketMailFuncOpen + ticketMailSettingCheck + `
	msg := TicketCreatedMessage(t, creator)
	if msg == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.SendMessage(ctx, msg)
}

// QueueTicketCreatedEmail puts the new-ticket email on the background queue.
//
// Contact-app review M30: this email used to be sent from a goroutine the
// request started, so a provider that was down for a minute lost it and a
// deploy dropped whatever was in flight. Queued, the worker retries it.
func QueueTicketCreatedEmail(ctx context.Context, q mail.Enqueuer, t *models.Ticket, creator *models.User) error {
	if !settings.Bool(ctx, "notifications.email_enabled") {
		log.Printf("ticket-mail: notification emails are turned off in settings, skipping ticket %s", t.ID)
		return nil
	}
	msg := TicketCreatedMessage(t, creator)
	if msg == nil {
		return nil
	}
	return mail.Queue(ctx, q, msg)
}

// TicketCreatedMessage builds the new-ticket email, or nil when SUPPORT_EMAIL
// is not set, which is the normal state of a development machine.
//
// Split out of SendTicketCreatedEmail so the same message can be queued or sent
// directly without the body existing in two places.
func TicketCreatedMessage(t *models.Ticket, creator *models.User) *mail.Message {
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

	return &mail.Message{To: []string{to}, Subject: subject, HTML: html}
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
