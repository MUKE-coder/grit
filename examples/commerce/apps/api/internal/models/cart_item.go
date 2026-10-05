package models

import (
	"time"

	"gorm.io/gorm"

	"commerce/apps/api/internal/ids"
	"commerce/apps/api/internal/money"
)

// CartItem represents a cartitem in the system.
type CartItem struct {
	ID           string         `gorm:"primarykey;size:36" json:"id"`
	CartID       string         `gorm:"size:36;index" json:"cart_id" binding:"required"`
	Cart         *Cart          `gorm:"foreignKey:CartID" json:"cart,omitempty"`
	ProductID    string         `gorm:"size:36;index" json:"product_id" binding:"required"`
	Product      *Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	VariantID    string         `gorm:"size:255" json:"variant_id" binding:"required" search:"trigram"`
	Quantity     int            `json:"quantity"`
	UnitPrice    money.Money    `gorm:"embedded;embeddedPrefix:unit_price_" json:"unit_price"`
	Title        string         `gorm:"size:255" json:"title" binding:"required" search:"trigram"`
	VariantLabel string         `gorm:"size:255" json:"variant_label" binding:"required" search:"trigram"`
	ImageURL     string         `gorm:"size:500" json:"image_url" binding:"required" search:"trigram"`
	Version      int            `gorm:"not null;default:1" json:"version"`
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"index" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	// ArchivedAt is the "put this away without destroying it" state, and it is
	// deliberately not DeletedAt. A soft delete is invisible to every query and
	// means the row is gone as far as the app is concerned; an archived row is
	// still listable, still exportable and still restorable in one click. The
	// list endpoint hides archived rows unless ?archived=true asks for them.
	ArchivedAt *time.Time `gorm:"index" json:"archived_at,omitempty"`
}

// BeforeCreate generates a UUID before inserting.
func (m *CartItem) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	return nil
}

// BeforeUpdate increments Version so offline clients can detect server-side updates.
func (m *CartItem) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}
