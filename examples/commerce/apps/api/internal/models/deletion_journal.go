package models

import (
	"time"

	"commerce/apps/api/internal/ids"
	"gorm.io/gorm"
)

// DeletionJournal is one tamper-evident record of a GDPR erasure. It proves a
// deletion happened without retaining the personal data that was deleted: it
// holds the erased user's opaque id (anonymized once their row is scrubbed), the
// admin who performed it, per-table counts, and a hash-chain link. Read-only —
// never updated or deleted, or the chain breaks.
type DeletionJournal struct {
	ID              string `gorm:"primarykey;size:36" json:"id"`
	DeletedUserID   string `gorm:"size:36;index" json:"deleted_user_id"` // UUID only — no PII
	ActorID         string `gorm:"size:36" json:"actor_id"`              // who ran the erasure
	ActorEmail      string `gorm:"size:255" json:"actor_email"`          // the admin, not the erased user
	Reason          string `gorm:"size:500" json:"reason"`
	RecordsAffected int    `json:"records_affected"`
	Counts          string `gorm:"type:text" json:"counts"` // JSON: {"uploads":3,"sessions":2}

	PrevHash string `gorm:"size:64" json:"prev_hash"`        // hex sha256, "" for the genesis row
	Hash     string `gorm:"size:64;uniqueIndex" json:"hash"` // hex sha256(prev_hash || canonical)

	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

func (d *DeletionJournal) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = ids.New()
	}
	return nil
}
