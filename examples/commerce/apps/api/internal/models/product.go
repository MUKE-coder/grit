package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"commerce/apps/api/internal/files"
	"commerce/apps/api/internal/ids"
	"commerce/apps/api/internal/money"
)

// Product represents a product in the system.
type Product struct {
	ID            string         `gorm:"primarykey;size:36" json:"id"`
	Title         string         `gorm:"size:255" json:"title" binding:"required" search:"trigram"`
	Handle        string         `gorm:"size:255;uniqueIndex" json:"handle" search:"trigram"`
	Description   string         `gorm:"type:text" json:"description" sanitize:"html" search:"trigram"`
	Price         money.Money    `gorm:"embedded;embeddedPrefix:price_" json:"price"`
	FeaturedImage *files.FileRef `gorm:"type:json" json:"featured_image"`
	Images        files.FileRefs `gorm:"type:json" json:"images"`
	Available     bool           `json:"available"`
	CollectionID  string         `gorm:"size:36;index" json:"collection_id" binding:"required"`
	Collection    *Collection    `gorm:"foreignKey:CollectionID" json:"collection,omitempty"`
	Version       int            `gorm:"not null;default:1" json:"version"`
	CreatedAt     time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"index" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	// ArchivedAt is the "put this away without destroying it" state, and it is
	// deliberately not DeletedAt. A soft delete is invisible to every query and
	// means the row is gone as far as the app is concerned; an archived row is
	// still listable, still exportable and still restorable in one click. The
	// list endpoint hides archived rows unless ?archived=true asks for them.
	ArchivedAt *time.Time `gorm:"index" json:"archived_at,omitempty"`
}

// BeforeCreate generates a UUID and auto-generates the slug before inserting.
func (m *Product) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	if m.Handle == "" {
		m.Handle = slugify(fmt.Sprintf("%v", m.Title))
	}
	return nil
}

// BeforeUpdate increments Version so offline clients can detect server-side updates.
func (m *Product) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}
