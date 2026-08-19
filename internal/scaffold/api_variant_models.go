package scaffold

import "strings"

// Product variants, normalised.
//
// The shape here is the one every catalogue converges on once it has more than
// one axis of choice, and the reason to write it down rather than let each
// project improvise is that the improvised version is nearly always the flat
// one: a variants table with a colour column and a size column. That works
// until somebody sells a laptop, where the axes are memory and storage, and
// then the schema needs a migration to sell a product.
//
// Five tables:
//
//	options                a choice a product can offer: Colour, Size, Memory
//	option_values          one choice on it: Red, XXL, 32GB
//	<resource>_options     which options THIS product uses, and in what order
//	<resource>_variants    one buyable combination: sku, stock, price, images
//	variant_option_values  the values that define a variant
//
// Three decisions in there are worth stating, because each has an obvious
// alternative that is worse.
//
// **Options are global, not per product.** Colour is Colour whether it is on a
// shirt or a phone case, and a shop with per-product options ends up with
// fourteen spellings of Colour and a filter that can only match one of them.
//
// **affects_price lives on the option, not the value.** "Does memory change the
// price" is a fact about memory. Putting it on each value means 32GB can be
// price-affecting while 16GB is not, which is not a thing anyone means, and it
// is a data state you then have to defend against when resolving a price.
//
// **Stock and images live on the variant, not the value.** Red/XXL selling out
// while Blue/XXL is in stock is the normal case, and the photo of the red one
// is a photo of that combination. A value-level image is the swatch, which is
// a different picture doing a different job, and both exist here.
//
// Price resolution is deliberately not stored. A variant carries an optional
// override, and where there is none the price is the product's plus the deltas
// of the values whose option declares affects_price. Storing the resolved
// figure means a product price change silently leaves every variant on the old
// one, which is the bug that turns into a support ticket about a wrong receipt.

// apiOptionModelGo emits internal/models/option.go and option_value.go.
func apiOptionModelGo(module string) string {
	return `package models

import (
	"time"

	"gorm.io/gorm"

	"` + module + `/internal/files"
	"` + module + `/internal/ids"
)

// Option is one axis of choice: Colour, Size, Memory.
//
// Shared across every resource that offers variants. A shop with per-product
// options accumulates fourteen spellings of Colour, and a filter can only ever
// match one of them.
type Option struct {
	ID   string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	Name string ` + "`" + `gorm:"size:255;not null" json:"name" binding:"required"` + "`" + `
	Slug string ` + "`" + `gorm:"size:255;uniqueIndex" json:"slug"` + "`" + `

	// Kind tells a storefront how to draw it: "swatch" for colours, "size" for
	// a row of boxes, "select" for a dropdown. It is presentation, and it lives
	// here so every shop renders the same option the same way.
	Kind string ` + "`" + `gorm:"size:32;default:select" json:"kind"` + "`" + `

	// AffectsPrice is a fact about the axis, not about one value on it.
	// Memory changes a laptop's price; colour does not. Per-value would allow
	// "32GB is price-affecting but 16GB is not", which nobody means.
	AffectsPrice bool ` + "`" + `json:"affects_price"` + "`" + `

	Position int ` + "`" + `gorm:"index;default:0" json:"position"` + "`" + `

	Values []OptionValue ` + "`" + `gorm:"foreignKey:OptionID" json:"values,omitempty"` + "`" + `

	Version   int            ` + "`" + `gorm:"not null;default:1" json:"version"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
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
	ID       string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	OptionID string ` + "`" + `gorm:"size:36;index;not null" json:"option_id" binding:"required"` + "`" + `
	Option   Option ` + "`" + `gorm:"foreignKey:OptionID" json:"option,omitempty"` + "`" + `

	Label string ` + "`" + `gorm:"size:255;not null" json:"label" binding:"required"` + "`" + `
	Slug  string ` + "`" + `gorm:"size:255;index" json:"slug"` + "`" + `

	// Swatch is a CSS colour for a colour option. Empty for every other kind.
	Swatch string ` + "`" + `gorm:"size:32" json:"swatch"` + "`" + `

	// Image is the swatch picture: this value, on its own, as a thumbnail.
	// It is NOT the photo of the product in that colour, which belongs to the
	// variant, because that photo is of a combination rather than of a value.
	Image *files.FileRef ` + "`" + `gorm:"serializer:json" json:"image,omitempty"` + "`" + `

	// PriceDelta is added to the product's price when this value is chosen,
	// and only when its option declares AffectsPrice. Signed, so a smaller
	// size can cost less.
	PriceDelta float64 ` + "`" + `gorm:"default:0" json:"price_delta"` + "`" + `

	Position int ` + "`" + `gorm:"index;default:0" json:"position"` + "`" + `

	Version   int            ` + "`" + `gorm:"not null;default:1" json:"version"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
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
`
}

// apiVariantModelGo emits internal/models/<resource>_variant.go: the variant
// itself, the product-to-option join, and the variant-to-value join.
func apiVariantModelGo(module, pascal, snake string) string {
	lower := strings.ToLower(pascal)
	return `package models

import (
	"time"

	"gorm.io/gorm"

	"` + module + `/internal/files"
	"` + module + `/internal/ids"
)

// ` + pascal + `Option records that a ` + lower + ` offers a given option, and in
// what order the storefront should show it.
//
// Without this table every ` + lower + ` would offer every option in the shop, and
// a t-shirt would ask the customer to choose a memory size.
type ` + pascal + `Option struct {
	ID        string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	` + pascal + `ID string ` + "`" + `gorm:"size:36;index:idx_` + snake + `_option,unique;not null" json:"` + snake + `_id" binding:"required"` + "`" + `
	OptionID  string ` + "`" + `gorm:"size:36;index:idx_` + snake + `_option,unique;not null" json:"option_id" binding:"required"` + "`" + `
	Option    Option ` + "`" + `gorm:"foreignKey:OptionID" json:"option,omitempty"` + "`" + `
	Position  int    ` + "`" + `gorm:"index;default:0" json:"position"` + "`" + `

	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

func (m *` + pascal + `Option) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	return nil
}

// ` + pascal + `Variant is one buyable combination.
//
// Stock and images live here rather than on an option value, because Red/XXL
// selling out while Blue/XXL is in stock is the normal case, and the photo of
// the red one is a photo of the combination.
type ` + pascal + `Variant struct {
	ID        string ` + "`" + `gorm:"primarykey;size:36" json:"id"` + "`" + `
	` + pascal + `ID string ` + "`" + `gorm:"size:36;index;not null" json:"` + snake + `_id" binding:"required"` + "`" + `

	SKU string ` + "`" + `gorm:"size:255;index" json:"sku"` + "`" + `

	// PriceOverride wins outright when set. Nil means the price is resolved
	// from the ` + lower + ` plus the deltas of the values that define this
	// variant, which is what keeps a price change from leaving stale copies
	// behind on every row.
	PriceOverride *float64 ` + "`" + `json:"price_override,omitempty"` + "`" + `

	Stock  int  ` + "`" + `gorm:"default:0" json:"stock"` + "`" + `
	Active bool ` + "`" + `gorm:"default:true" json:"active"` + "`" + `

	Images files.FileRefs ` + "`" + `gorm:"serializer:json" json:"images,omitempty"` + "`" + `

	Position int ` + "`" + `gorm:"index;default:0" json:"position"` + "`" + `

	// OptionValues are the values that define this combination, one per option
	// the ` + lower + ` offers.
	OptionValues []OptionValue ` + "`" + `gorm:"many2many:variant_option_values;foreignKey:ID;joinForeignKey:VariantID;References:ID;joinReferences:OptionValueID" json:"option_values,omitempty"` + "`" + `

	Version   int            ` + "`" + `gorm:"not null;default:1" json:"version"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

func (m *` + pascal + `Variant) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ids.New()
	}
	return nil
}

func (m *` + pascal + `Variant) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	return nil
}

// InStock is the question a storefront actually asks. Published instead of the
// raw count, which is a business fact competitors enjoy.
func (m *` + pascal + `Variant) InStock() bool {
	return m.Active && m.Stock > 0
}
`
}

// Exported so grit add variants, which lives in internal/generate, can write
// these into a project. Same reason APIDialectGo and APIPaginateGo are
// exported: generate imports scaffold, never the other way round.
func APIOptionModelGo(module string) string { return apiOptionModelGo(module) }

func APIVariantModelGo(module, pascal, snake string) string {
	return apiVariantModelGo(module, pascal, snake)
}
