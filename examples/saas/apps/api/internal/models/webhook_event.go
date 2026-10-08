package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"saas/apps/api/internal/ids"
)

// WebhookEvent persists every webhook the API receives. ExternalID is
// the provider's own event ID — we use it as the idempotency key, so
// duplicate deliveries (Stripe retries, partner pings) become no-ops.
//
// Status lifecycle:
//
//	pending   — received + verified, handler not yet run
//	processed — handler returned nil
//	failed    — handler returned an error; HandlerError holds the message
//	skipped   — duplicate ExternalID — handler was bypassed
type WebhookEvent struct {
	ID        string `gorm:"primarykey;size:36" json:"id"`
	Provider  string `gorm:"size:50;not null;uniqueIndex:idx_webhook_provider_external,priority:1" json:"provider"`
	EventType string `gorm:"size:100;index" json:"event_type"`

	// ExternalID is the provider's own event id, and the second half of the
	// idempotency key. A pointer because it has to be NULL when a provider
	// does not supply one: every database allows repeated NULLs in a unique
	// index and none allow a repeated empty string, so storing "" would make
	// two unrelated anonymous events collide and silently drop the second.
	ExternalID *string `gorm:"size:255;uniqueIndex:idx_webhook_provider_external,priority:2" json:"external_id,omitempty"`
	// No explicit type: datatypes.JSON maps to jsonb on Postgres and json on
	// MySQL by itself. Naming jsonb here fails AutoMigrate on MySQL, which has
	// no such type.
	Payload      datatypes.JSON `json:"payload"`
	Status       string         `gorm:"size:20;index;not null;default:pending" json:"status"`
	HandlerError string         `gorm:"type:text" json:"handler_error,omitempty"`
	RetryCount   int            `gorm:"not null;default:0" json:"retry_count"`
	ProcessedAt  *time.Time     `json:"processed_at,omitempty"`
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
}

func (w *WebhookEvent) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = ids.New()
	}
	return nil
}

// The unique index is declared on the fields above rather than built here.
//
// It used to live in a method returning DDL as a string, which nothing
// called. The column had a plain index, the INSERT never failed, and the
// handler's "duplicate means already processed" branch was unreachable: every
// retried delivery ran the handler again. Declaring it as a tag means the
// migration creates it, and the constraint is where it can be seen.
