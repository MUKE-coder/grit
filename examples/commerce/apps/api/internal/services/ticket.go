package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"

	"commerce/apps/api/internal/mail"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/respond"
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
		Labels:      NormalizeTicketLabels(in.Labels, MaxTicketLabels),
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
