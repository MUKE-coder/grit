package models

import (
	"time"

	"gorm.io/gorm"
	"library/apps/api/internal/ids"
)

// A ticket's status and priority. The column is a string, so these are the
// values it holds; the admin panel offers the same four priorities.
const (
	TicketStatusOpen   = "open"
	TicketStatusClosed = "closed"

	TicketPriorityLow      = "low"
	TicketPriorityMedium   = "medium"
	TicketPriorityHigh     = "high"
	TicketPriorityCritical = "critical"
)

// Ticket — single support request. Anyone authenticated can open one;
// ADMIN/EDITOR roles see the full queue, regular USER role sees their
// own. Status is open by default; closing stamps ClosedAt.
type Ticket struct {
	ID          string         `gorm:"primarykey;size:36" json:"id"`
	UserID      string         `gorm:"size:36;index;not null" json:"user_id"`
	Subject     string         `gorm:"size:200;not null" json:"subject"`
	Description string         `gorm:"type:text;not null" json:"description"`
	Status      string         `gorm:"size:16;index;default:'open'" json:"status"`     // open | closed
	Priority    string         `gorm:"size:16;index;default:'medium'" json:"priority"` // low | medium | high | critical
	Labels      string         `gorm:"size:255" json:"labels"`                         // comma-separated, capped by the ticket service
	AssigneeID  string         `gorm:"size:36;index" json:"assignee_id"`               // optional, must be ADMIN role
	LastReplyAt *time.Time     `json:"last_reply_at"`                                  // touched by every reply
	ClosedAt    *time.Time     `json:"closed_at"`
	CreatedAt   time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Eager-loaded relations for the detail page.
	User     *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Assignee *User         `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	Replies  []TicketReply `gorm:"foreignKey:TicketID" json:"replies,omitempty"`
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = ids.New()
	}
	if t.Status == "" {
		t.Status = TicketStatusOpen
	}
	if t.Priority == "" {
		t.Priority = TicketPriorityMedium
	}
	return nil
}

// TicketReply — chronological thread under a ticket. IsAdminReply lets
// the dashboard style admin replies differently (badged, opposite side
// of the thread, etc.).
type TicketReply struct {
	ID           string         `gorm:"primarykey;size:36" json:"id"`
	TicketID     string         `gorm:"size:36;index;not null" json:"ticket_id"`
	UserID       string         `gorm:"size:36;index;not null" json:"user_id"`
	Body         string         `gorm:"type:text;not null" json:"body"`
	IsAdminReply bool           `json:"is_admin_reply"`
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (r *TicketReply) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = ids.New()
	}
	return nil
}
