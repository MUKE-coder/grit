package models

import (
	"time"

	"gorm.io/gorm"
)

// Upload represents a file uploaded to storage.
type Upload struct {
	ID           string `gorm:"primarykey;size:36" json:"id"`
	Filename     string `gorm:"size:255;not null" json:"filename"`
	OriginalName string `gorm:"size:255;not null" json:"original_name"`
	MimeType     string `gorm:"size:100;not null" json:"mime_type"`
	Size         int64  `gorm:"not null" json:"size"`
	Path         string `gorm:"size:500;not null;index" json:"path"`
	URL          string `gorm:"size:500" json:"url"`
	ThumbnailURL string `gorm:"size:500" json:"thumbnail_url"`
	UserID       string `gorm:"size:36;index;not null" json:"user_id"`
	User         User   `gorm:"foreignKey:UserID" json:"-"`
	Version      int    `gorm:"not null;default:1" json:"version"`
	// v3.31.33 -- claimed_at is set when a parent record references this
	// upload's path/key via a FileRef column. NULL means abandoned, and
	// the daily orphan cleanup cron deletes the S3 object + DB row when
	// the upload is older than 24h and still unclaimed.
	ClaimedAt *time.Time     `gorm:"index" json:"claimed_at,omitempty"`
	CreatedAt time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeUpdate increments Version on every server-side write so offline
// clients can detect that a record they edited has moved on.
func (u *Upload) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}
