package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"library/apps/api/internal/ids"
)

// Message is one thing to be delivered, and the record that it was.
//
// Kept after delivery rather than deleted. The row is the evidence that an
// event was published, which is what you want when a downstream system says it
// never received one, and it is what a replay reads from.
type OutboxMessage struct {
	ID    string `gorm:"type:varchar(36);primaryKey" json:"id"`
	Topic string `gorm:"size:255;not null;index:idx_outbox_topic" json:"topic"`

	// Key is the caller's idempotency key for this message, unique across the
	// table. Enqueueing the same key twice is a no-op rather than an error, so
	// a retried request cannot publish the same event twice.
	//
	// Empty means "no key", and an empty string cannot be unique across many
	// rows, so it is stored as NULL. That is why this is a pointer.
	Key *string `gorm:"size:255;uniqueIndex" json:"key,omitempty"`

	Payload datatypes.JSON `gorm:"type:json" json:"payload"`

	Status    string `gorm:"size:20;not null;default:'pending';index:idx_outbox_claim,priority:1" json:"status"`
	Attempts  int    `gorm:"not null;default:0" json:"attempts"`
	LastError string `gorm:"type:text" json:"last_error,omitempty"`

	// AvailableAt is when the relay may next try. It moves forward on each
	// failure, which is the backoff.
	AvailableAt time.Time `gorm:"not null;index:idx_outbox_claim,priority:2" json:"available_at"`

	ClaimedBy   string     `gorm:"size:64" json:"claimed_by,omitempty"`
	ClaimedAt   *time.Time `json:"claimed_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`

	CreatedAt time.Time `gorm:"index" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (OutboxMessage) TableName() string { return "outbox_messages" }

// BeforeCreate fills the id and the first availability.
func (m *OutboxMessage) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	if m.AvailableAt.IsZero() {
		m.AvailableAt = time.Now()
	}
	if m.Status == "" {
		m.Status = "pending"
	}
	return nil
}
