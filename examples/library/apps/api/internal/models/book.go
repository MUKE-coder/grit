package models

import (
	"time"

	"gorm.io/gorm"

	"library/apps/api/internal/files"
	"library/apps/api/internal/ids"
	"library/apps/api/internal/jsontime"
	"library/apps/api/internal/money"
)

// Book represents a book in the system.
type Book struct {
	ID        string         `gorm:"primarykey;size:36" json:"id"`
	Title     string         `gorm:"size:255" json:"title" binding:"required" search:"trigram"`
	Isbn      string         `gorm:"size:255" json:"isbn" binding:"required" search:"trigram"`
	Summary   string         `gorm:"type:text" json:"summary" sanitize:"html" search:"trigram"`
	Published *jsontime.Date `gorm:"type:date" json:"published"`
	Price     money.Money    `gorm:"embedded;embeddedPrefix:price_" json:"price"`
	Cover     *files.FileRef `gorm:"type:json" json:"cover"`
	Genre     string         `gorm:"size:255" json:"genre" binding:"required"`
	AuthorID  string         `gorm:"size:36;index" json:"author_id" binding:"required"`
	Author    *Author        `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
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
func (m *Book) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	return nil
}

// BeforeUpdate increments Version so offline clients can detect server-side updates.
func (m *Book) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}
