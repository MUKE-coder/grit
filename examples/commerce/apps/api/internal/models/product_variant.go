package models

import (
	"time"

	"gorm.io/gorm"

	"commerce/apps/api/internal/files"
	"commerce/apps/api/internal/ids"
)

// ProductOption records that a product offers a given option, and in
// what order the storefront should show it.
//
// Without this table every product would offer every option in the shop, and
// a t-shirt would ask the customer to choose a memory size.
type ProductOption struct {
	ID        string `gorm:"primarykey;size:36" json:"id"`
	ProductID string `gorm:"size:36;index:idx_product_option,unique;not null" json:"product_id" binding:"required"`
	OptionID  string `gorm:"size:36;index:idx_product_option,unique;not null" json:"option_id" binding:"required"`
	Option    Option `gorm:"foreignKey:OptionID" json:"option,omitempty"`
	Position  int    `gorm:"index;default:0" json:"position"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *ProductOption) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	return nil
}

// ProductOptionValue narrows an axis to the values THIS product
// offers.
//
// Options are shop-wide on purpose: Colour is Colour everywhere, so a filter
// can match across the catalogue. That leaves a gap any real shop hits on its
// second product, which is that a shirt in ecru and navy and a tee in black,
// sand and olive cannot share the Colour axis: offering the axis offered every
// colour in the shop. Found building a storefront on Grit, which worked around
// it with two colour options, one per product.
//
// No rows for an option means every value of it, so a product that has
// never been narrowed behaves exactly as before and nothing has to be
// backfilled.
type ProductOptionValue struct {
	ID            string      `gorm:"primarykey;size:36" json:"id"`
	ProductID     string      `gorm:"size:36;index:idx_product_option_value,unique;not null" json:"product_id" binding:"required"`
	OptionValueID string      `gorm:"size:36;index:idx_product_option_value,unique;not null" json:"option_value_id" binding:"required"`
	OptionValue   OptionValue `gorm:"foreignKey:OptionValueID" json:"option_value,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *ProductOptionValue) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	return nil
}

// ProductVariant is one buyable combination.
//
// Stock and images live here rather than on an option value, because Red/XXL
// selling out while Blue/XXL is in stock is the normal case, and the photo of
// the red one is a photo of the combination.
type ProductVariant struct {
	ID        string `gorm:"primarykey;size:36" json:"id"`
	ProductID string `gorm:"size:36;index;not null" json:"product_id" binding:"required"`

	SKU string `gorm:"size:255;index" json:"sku"`

	// PriceOverride wins outright when set. Nil means the price is resolved
	// from the product plus the deltas of the values that define this
	// variant, which is what keeps a price change from leaving stale copies
	// behind on every row.
	PriceOverride *float64 `json:"price_override,omitempty"`

	Stock  int  `gorm:"default:0" json:"stock"`
	Active bool `gorm:"default:true" json:"active"`

	Images files.FileRefs `gorm:"serializer:json" json:"images,omitempty"`

	Position int `gorm:"index;default:0" json:"position"`

	// OptionValues are the values that define this combination, one per option
	// the product offers.
	OptionValues []OptionValue `gorm:"many2many:variant_option_values;foreignKey:ID;joinForeignKey:VariantID;References:ID;joinReferences:OptionValueID" json:"option_values,omitempty"`

	Version   int            `gorm:"not null;default:1" json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *ProductVariant) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	return nil
}

func (m *ProductVariant) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}

// InStock is the question a storefront actually asks. Published instead of the
// raw count, which is a business fact competitors enjoy.
func (m *ProductVariant) InStock() bool {
	return m.Active && m.Stock > 0
}
