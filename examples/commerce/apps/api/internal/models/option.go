package models

import (
	"time"

	"gorm.io/gorm"

	"commerce/apps/api/internal/files"
	"commerce/apps/api/internal/ids"
)

// Option is one axis of choice: Colour, Size, Memory.
//
// Shared across every resource that offers variants. A shop with per-product
// options accumulates fourteen spellings of Colour, and a filter can only ever
// match one of them.
type Option struct {
	ID   string `gorm:"primarykey;size:36" json:"id"`
	Name string `gorm:"size:255;not null" json:"name" binding:"required"`
	Slug string `gorm:"size:255;uniqueIndex" json:"slug"`

	// Kind tells a storefront how to draw it: "swatch" for colours, "size" for
	// a row of boxes, "select" for a dropdown. It is presentation, and it lives
	// here so every shop renders the same option the same way.
	Kind string `gorm:"size:32;default:select" json:"kind"`

	// AffectsPrice is a fact about the axis, not about one value on it.
	// Memory changes a laptop's price; colour does not. Per-value would allow
	// "32GB is price-affecting but 16GB is not", which nobody means.
	AffectsPrice bool `json:"affects_price"`

	Position int `gorm:"index;default:0" json:"position"`

	Values []OptionValue `gorm:"foreignKey:OptionID" json:"values,omitempty"`

	Version   int            `gorm:"not null;default:1" json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *Option) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	if m.Slug == "" {
		m.Slug = slugify(m.Name)
	}
	return nil
}

func (m *Option) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}

// OptionValue is one choice on an option: Red, XXL, 32GB.
type OptionValue struct {
	ID       string `gorm:"primarykey;size:36" json:"id"`
	OptionID string `gorm:"size:36;index;not null" json:"option_id" binding:"required"`
	Option   Option `gorm:"foreignKey:OptionID" json:"option,omitempty"`

	Label string `gorm:"size:255;not null" json:"label" binding:"required"`
	Slug  string `gorm:"size:255;index" json:"slug"`

	// Swatch is a CSS colour for a colour option. Empty for every other kind.
	Swatch string `gorm:"size:32" json:"swatch"`

	// Image is the swatch picture: this value, on its own, as a thumbnail.
	// It is NOT the photo of the product in that colour, which belongs to the
	// variant, because that photo is of a combination rather than of a value.
	Image *files.FileRef `gorm:"serializer:json" json:"image,omitempty"`

	// PriceDelta is added to the product's price when this value is chosen,
	// and only when its option declares AffectsPrice. Signed, so a smaller
	// size can cost less.
	PriceDelta float64 `gorm:"default:0" json:"price_delta"`

	Position int `gorm:"index;default:0" json:"position"`

	Version   int            `gorm:"not null;default:1" json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *OptionValue) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	if m.Slug == "" {
		m.Slug = slugify(m.Label)
	}
	return nil
}

func (m *OptionValue) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}
