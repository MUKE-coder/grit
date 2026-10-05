package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"commerce/apps/api/internal/ids"
)

// Page represents a page in the system.
type Page struct {
	ID        string         `gorm:"primarykey;size:36" json:"id"`
	Title     string         `gorm:"size:255" json:"title" binding:"required" search:"trigram"`
	Handle    string         `gorm:"size:255;uniqueIndex" json:"handle" search:"trigram"`
	Body      string         `gorm:"type:text" json:"body" sanitize:"html" search:"trigram"`
	Version   int            `gorm:"not null;default:1" json:"version"`
	CreatedAt time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"index" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	// ArchivedAt is the "put this away without destroying it" state, and it is
	// deliberately not DeletedAt. A soft delete is invisible to every query and
	// means the row is gone as far as the app is concerned; an archived row is
	// still listable, still exportable and still restorable in one click. The
	// list endpoint hides archived rows unless ?archived=true asks for them.
	ArchivedAt *time.Time `gorm:"index" json:"archived_at,omitempty"`
}

// BeforeCreate generates a UUID and auto-generates the slug before inserting.
func (m *Page) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	if m.Handle == "" {
		m.Handle = slugify(fmt.Sprintf("%v", m.Title))
	}
	return nil
}

// BeforeUpdate increments Version so offline clients can detect server-side updates.
func (m *Page) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}
