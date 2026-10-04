package models

import (
	"time"

	"gorm.io/gorm"

	"library/apps/api/internal/ids"
)

// Author represents a author in the system.
type Author struct {
	ID        string         `gorm:"primarykey;size:36" json:"id"`
	Name      string         `gorm:"size:255" json:"name" binding:"required" search:"trigram"`
	Bio       string         `gorm:"type:text" json:"bio" search:"trigram"`
	Website   string         `gorm:"size:2048" json:"website" format:"url" docs:"format:uri,example:https://example.com/pricing" search:"trigram"`
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

// BeforeCreate generates a UUID before inserting.
func (m *Author) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	return nil
}

// BeforeUpdate increments Version so offline clients can detect server-side updates.
func (m *Author) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}
