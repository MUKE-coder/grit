package models

import (
	"time"

	"gorm.io/gorm"
	"saas/apps/api/internal/ids"
)

// ActivityLog records every successful authenticated mutation, with a
// tamper-evident hash chain — each row's Hash is SHA-256 of (PrevHash
// || canonical(this_row)). Mutating any row breaks the chain on the
// next VerifyChain pass.
//
// The payload digest is a SHA-256 of the request body so we have
// evidence of what was sent without storing PII verbatim. Read-only —
// no updates, no deletes (use a separate retention job to prune old
// rows; deletion still breaks the chain so it must rebuild from a
// safe checkpoint).
type ActivityLog struct {
	ID            string `gorm:"primarykey;size:36" json:"id"`
	UserID        string `gorm:"size:36;index" json:"user_id"`
	Method        string `gorm:"size:10" json:"method"`
	Path          string `gorm:"size:500;index" json:"path"`
	Status        int    `json:"status"`
	PayloadDigest string `gorm:"size:64" json:"payload_digest"` // sha256 hex
	IPAddress     string `gorm:"size:45" json:"ip_address"`
	UserAgent     string `gorm:"size:500" json:"user_agent"`
	DurationMS    int64  `json:"duration_ms"`
	// What an audited read served, or what a reseal covered. Empty on every
	// other entry, and left out of the hash when empty, so entries written
	// before these fields existed hash exactly as they did.
	Resource    string    `gorm:"size:100;index" json:"resource,omitempty"`
	ResourceIDs string    `gorm:"type:text" json:"resource_ids,omitempty"`
	RecordCount int       `json:"record_count,omitempty"`
	PrevHash    string    `gorm:"size:64" json:"prev_hash"`        // hex sha256, "" for the genesis row
	Hash        string    `gorm:"size:64;uniqueIndex" json:"hash"` // hex sha256(prev_hash || canonical)
	CreatedAt   time.Time `gorm:"index" json:"created_at"`
}

func (a *ActivityLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = ids.New()
	}
	return nil
}
