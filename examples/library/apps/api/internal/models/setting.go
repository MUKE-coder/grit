package models

import "time"

// Setting is one stored override.
//
// Keyed by (key, scope, scope_id) rather than by key alone, because the same
// setting can hold a global value, a per-tenant value and a per-user value at
// once, and resolution picks the most specific one that exists.
//
// ScopeID is empty for a global row. A partial unique index would express
// that better on Postgres, but a plain composite unique works on all three
// dialects Grit supports and the empty string is a perfectly good sentinel
// for "no scope".
type Setting struct {
	ID        string    `gorm:"primarykey;size:36" json:"id"`
	Key       string    `gorm:"size:100;not null;uniqueIndex:idx_setting_scope" json:"key"`
	Scope     string    `gorm:"size:16;not null;default:global;uniqueIndex:idx_setting_scope" json:"scope"`
	ScopeID   string    `gorm:"size:36;not null;default:'';uniqueIndex:idx_setting_scope" json:"scope_id"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedBy string    `gorm:"size:36" json:"updated_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Setting) TableName() string { return "settings" }
