package models

import (
	"time"

	"commerce/apps/api/internal/ids"
	"gorm.io/gorm"
)

// FormSubmission is an audit-log row for one successful public form
// submission. Created by FormShareHandler.PublicSubmit after the
// dispatcher returns a record ID. Never modified; soft-deletable so
// admins can prune old rows without losing the share's history.
type FormSubmission struct {
	ID           string `gorm:"primarykey;size:36" json:"id"`
	ShareID      string `gorm:"size:36;not null;index" json:"share_id"`
	ResourceName string `gorm:"size:64;not null;index" json:"resource_name"`
	RecordID     string `gorm:"size:36;not null;index" json:"record_id"`
	// IP and UserAgent are best-effort — set from gin.Context at
	// submission time. Truncated to fit the column; long UAs are
	// trimmed to 500 chars.
	IP        string         `gorm:"size:64" json:"ip"`
	UserAgent string         `gorm:"size:500" json:"user_agent"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *FormSubmission) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = ids.New()
	}
	return nil
}
