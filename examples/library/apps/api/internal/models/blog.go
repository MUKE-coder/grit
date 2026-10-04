package models

import (
	"time"

	"gorm.io/gorm"
	"library/apps/api/internal/ids"
)

// Blog represents a blog post in the system.
type Blog struct {
	ID          string         `gorm:"primarykey;size:36" json:"id"`
	Title       string         `gorm:"size:255;not null" json:"title" binding:"required"`
	Slug        string         `gorm:"size:255;uniqueIndex" json:"slug"`
	Content     string         `gorm:"type:text" json:"content" sanitize:"html"`
	Image       string         `gorm:"size:500" json:"image"`
	Excerpt     string         `gorm:"size:500" json:"excerpt"`
	Published   bool           `gorm:"default:false" json:"published"`
	PublishedAt *time.Time     `json:"published_at"`
	Version     int            `gorm:"not null;default:1" json:"version"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `gorm:"index" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate auto-generates a UUID and the slug before inserting.
func (b *Blog) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = ids.New()
	}
	if b.Slug == "" {
		b.Slug = slugify(b.Title)
	}
	return nil
}

// BeforeUpdate increments Version so offline clients can detect that
// a record they edited has moved on.
func (b *Blog) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}
