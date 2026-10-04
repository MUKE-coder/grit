package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"library/apps/api/internal/jobs"
	"library/apps/api/internal/mail"
	"library/apps/api/internal/models"
	"library/apps/api/internal/paginate"
	"library/apps/api/internal/respond"
	"library/apps/api/internal/services"
)

type TicketHandler struct {
	DB   *gorm.DB
	Mail *mail.Mailer // can be nil — the new-ticket email is skipped
	// Jobs is optional: with it the new-ticket email goes on the background
	// queue, so a provider that is briefly down does not lose it.
	Jobs *jobs.Client
}

type CreateTicketRequest struct {
	Subject     string `json:"subject" binding:"required,min=3,max=200"`
	Description string `json:"description" binding:"required,min=10"`
	Priority    string `json:"priority"`
	Labels      string `json:"labels"`
}

type TicketReplyRequest struct {
	Body string `json:"body" binding:"required,min=1"`
}

type AssignTicketRequest struct {
	AssigneeID string `json:"assignee_id" binding:"required"`
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
