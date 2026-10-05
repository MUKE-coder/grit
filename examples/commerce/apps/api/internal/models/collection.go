package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"commerce/apps/api/internal/files"
	"commerce/apps/api/internal/ids"
)

// Collection represents a collection in the system.
type Collection struct {
	ID          string         `gorm:"primarykey;size:36" json:"id"`
	Title       string         `gorm:"size:255" json:"title" binding:"required" search:"trigram"`
	Handle      string         `gorm:"size:255;uniqueIndex" json:"handle" search:"trigram"`
	Description string         `gorm:"type:text" json:"description" search:"trigram"`
	Image       *files.FileRef `gorm:"type:json" json:"image"`
	Version     int            `gorm:"not null;default:1" json:"version"`
	CreatedAt   time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"index" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	// ArchivedAt is the "put this away without destroying it" state, and it is
	// deliberately not DeletedAt. A soft delete is invisible to every query and
	// means the row is gone as far as the app is concerned; an archived row is
	// still listable, still exportable and still restorable in one click. The
	// list endpoint hides archived rows unless ?archived=true asks for them.
	ArchivedAt *time.Time `gorm:"index" json:"archived_at,omitempty"`
}

// BeforeCreate generates a UUID and auto-generates the slug before inserting.
func (m *Collection) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	if m.Handle == "" {
		m.Handle = slugify(fmt.Sprintf("%v", m.Title))
	}
	return nil
}

// BeforeUpdate increments Version so offline clients can detect server-side updates.
func (m *Collection) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}
