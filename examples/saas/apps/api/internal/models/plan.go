package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"saas/apps/api/internal/ids"
	"saas/apps/api/internal/money"
)

// Plan represents a plan in the system.
type Plan struct {
	ID        string         `gorm:"primarykey;size:36" json:"id"`
	Name      string         `gorm:"size:255" json:"name" binding:"required" search:"trigram"`
	Slug      string         `gorm:"size:255;uniqueIndex" json:"slug" search:"trigram"`
	Price     money.Money    `gorm:"embedded;embeddedPrefix:price_" json:"price"`
	Interval  string         `gorm:"size:255" json:"interval" binding:"required"`
	Seats     int            `json:"seats"`
	Blurb     string         `gorm:"type:text" json:"blurb" search:"trigram"`
	Active    bool           `json:"active"`
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
func (m *Plan) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	if m.Slug == "" {
		m.Slug = slugify(fmt.Sprintf("%v", m.Name))
	}
	return nil
}

// BeforeUpdate increments Version so offline clients can detect server-side updates.
func (m *Plan) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}
