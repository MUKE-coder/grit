package scaffold

import "strings"

// Seed data for the variant system.
//
// A shop with variants installed and no variants in it is a storefront you
// cannot look at: the picker renders nothing, the price range collapses to one
// figure, and the disabled-swatch state that is most of the work on a product
// page never appears. So the seeder writes a small library and a real matrix,
// and deliberately leaves one combination out of stock so that state is on
// screen the first time anybody opens the page.
//
// Deterministic rather than faked. Stock and overrides come from the row's
// index, so two people running `grit seed` see the same catalogue and a
// screenshot in the docs keeps matching the thing it documents.

// APIOptionSeederGo emits internal/database/options_seeder.go: the shared
// option library.
//
// Written once however many resources offer variants, like every other shared
// piece of this system. Colour and Size are the two axes almost every shop
// starts with, and one of them affects price so the resolution logic is
// exercised by the seed rather than only by the tests.
func APIOptionSeederGo(module string) string {
	return `package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"` + module + `/internal/models"
)

// SeedOptions writes the shared option library: Colour and Size.
//
// Idempotent, and called by every variant seeder rather than registered in
// Seed() on its own. Options are shop-wide, so "which seeder owns them" has no
// good answer, and having each one ensure they exist removes the ordering
// question entirely.
func SeedOptions(db *gorm.DB) error {
	// Colour does not affect price. That is the honest default, and it is worth
	// seeding the flag as false so the first thing anybody reads in the admin is
	// an option where the price column is empty for a reason.
	colour := []models.OptionValue{
		{Label: "Black", Swatch: "#111118", Position: 0},
		{Label: "White", Swatch: "#f4f4f6", Position: 1},
		{Label: "Navy", Swatch: "#1f2a44", Position: 2},
		{Label: "Sand", Swatch: "#d8cbb4", Position: 3},
	}
	if _, err := ensureOption(db,
		models.Option{Name: "Colour", Kind: "swatch", AffectsPrice: false, Position: 0},
		colour); err != nil {
		return err
	}

	// Size does, and only at the top end, which is what the price_delta column
	// is for: an XL costs more to make and the delta says so without needing a
	// second product.
	size := []models.OptionValue{
		{Label: "S", Position: 0},
		{Label: "M", Position: 1},
		{Label: "L", Position: 2},
		{Label: "XL", PriceDelta: 2.50, Position: 3},
	}
	if _, err := ensureOption(db,
		models.Option{Name: "Size", Kind: "size", AffectsPrice: true, Position: 1},
		size); err != nil {
		return err
	}

	return nil
}

// ensureOption creates an option and its values if it is not already there,
// and returns whichever one now is.
//
// Matched on name rather than on the slug the model derives. The slug would be
// the better key if this file could compute it, but that means reimplementing
// the model's hook here, and the copy goes stale the first time somebody
// improves the real one.
func ensureOption(db *gorm.DB, option models.Option, values []models.OptionValue) (models.Option, error) {
	// Find with a limit rather than First. Absence is the expected case on a
	// fresh database, and First reports it as ErrRecordNotFound, which GORM logs
	// in red. A seeder whose happy path prints two errors is a seeder people
	// stop reading the output of.
	var existing models.Option
	if err := db.Preload("Values").Where("name = ?", option.Name).
		Limit(1).Find(&existing).Error; err != nil {
		return models.Option{}, fmt.Errorf("looking up the %s option: %w", option.Name, err)
	}
	if existing.ID != "" {
		return existing, nil
	}

	if err := db.Create(&option).Error; err != nil {
		return models.Option{}, fmt.Errorf("creating the %s option: %w", option.Name, err)
	}
	for i := range values {
		values[i].OptionID = option.ID
		if err := db.Create(&values[i]).Error; err != nil {
			return models.Option{}, fmt.Errorf("adding %s to %s: %w", values[i].Label, option.Name, err)
		}
	}
	option.Values = values
	log.Printf("Seeded the %s option with %d values", option.Name, len(values))
	return option, nil
}
`
}

// APIVariantSeederGo emits internal/database/<snake>_variants_seeder.go.
//
// hasSlug decides where the generated SKUs get their prefix. A resource with a
// slug gives readable ones (AURA-TEE-BLACK-XL); without one the prefix falls
// back to a slice of the id, which is ugly and still unique, and unique is the
// part that matters to a warehouse.
func APIVariantSeederGo(module, pascal, snake, plural string, hasSlug bool) string {
	lower := strings.ToLower(pascal)

	prefix := "shortID(row.ID)"
	if hasSlug {
		prefix = "skuPrefix(row.Slug, row.ID)"
	}

	return `package database

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"

	"` + module + `/internal/models"
	"` + module + `/internal/services"
)

// Seed` + pascal + `Variants gives the first few ` + plural + ` a real matrix.
//
// A few rather than all of them, on purpose. A catalogue where every row has
// sixteen combinations is a catalogue where nothing tests the plain case, and
// the plain case is most of a real shop: a ` + lower + ` with no options at all
// still has to render, price and sell.
func Seed` + pascal + `Variants(db *gorm.DB) error {
	// How many ` + plural + ` get options attached.
	const withVariants = 6

	var existing int64
	db.Model(&models.` + pascal + `Variant{}).Count(&existing)
	if existing > 0 {
		log.Println("` + pascal + ` variants already seeded, skipping...")
		return nil
	}

	// Idempotent, and called here rather than registered separately: options are
	// shop-wide, so no one seeder owns them.
	if err := SeedOptions(db); err != nil {
		return err
	}

	var options []models.Option
	if err := db.Where("name IN ?", []string{"Colour", "Size"}).
		Order("position asc").Find(&options).Error; err != nil {
		return fmt.Errorf("loading the seeded options: %w", err)
	}
	if len(options) == 0 {
		log.Println("No options to attach, skipping ` + lower + ` variants...")
		return nil
	}

	var rows []models.` + pascal + `
	if err := db.Order("created_at asc").Limit(withVariants).Find(&rows).Error; err != nil {
		return fmt.Errorf("loading ` + plural + ` to attach variants to: %w", err)
	}
	if len(rows) == 0 {
		log.Println("No ` + plural + ` exist yet, skipping variants. Seed ` + pascal + ` first, then run grit seed again.")
		return nil
	}

	variants := services.New` + pascal + `VariantService(db)
	seeded := 0

	for _, row := range rows {
		for position, option := range options {
			link := models.` + pascal + `Option{
				` + pascal + `ID: row.ID,
				OptionID: option.ID,
				Position: position,
			}
			err := db.Where("` + snake + `_id = ? AND option_id = ?", row.ID, option.ID).
				FirstOrCreate(&link).Error
			if err != nil {
				return fmt.Errorf("attaching %s to a ` + lower + `: %w", option.Name, err)
			}
		}

		// The same generator the admin's button calls, so the seed cannot drift
		// from what the UI produces.
		if _, err := variants.Generate(row.ID, 0); err != nil {
			return fmt.Errorf("generating the matrix: %w", err)
		}

		created, err := variants.VariantsFor(row.ID)
		if err != nil {
			return err
		}
		for i, variant := range created {
			updates := map[string]any{
				"sku":   ` + prefix + ` + "-" + valueSuffix(variant.OptionValues),
				"stock": 4 + (i*7)%37,
			}
			// One combination in seven is out of stock. The disabled swatch is
			// most of the work on a product page and the hardest state to
			// remember to build, so the seed puts it on screen unasked.
			if i%7 == 3 {
				updates["stock"] = 0
			}
			if err := db.Model(&models.` + pascal + `Variant{}).
				Where("id = ?", variant.ID).Updates(updates).Error; err != nil {
				return fmt.Errorf("filling in a variant: %w", err)
			}
			seeded++
		}
	}

	log.Printf("Seeded %d ` + lower + ` variants across %d ` + plural + `", seeded, len(rows))
	return nil
}

// valueSuffix builds the readable half of a SKU from the values that define the
// combination: BLACK-XL rather than a serial number nobody can check against a
// shelf.
//
// From the label and not the slug. Slugs carry a uniqueness suffix, so the slug
// half of a SKU reads BLACK-ADF9C36E-XL-93661012, which is unique, unreadable,
// and no use at all to the person holding the box.
func valueSuffix(values []models.OptionValue) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, strings.ToUpper(strings.ReplaceAll(value.Label, " ", "-")))
	}
	return strings.Join(parts, "-")
}
` + skuPrefixHelpers(hasSlug)
}

// skuPrefixHelpers emits whichever prefix helper the seeder above referenced.
//
// One or the other, never both: an unused function is not a build failure in Go
// but it is a question the next reader has to answer, and the answer is "a
// generator emitted it for a case your project is not in".
func skuPrefixHelpers(hasSlug bool) string {
	if hasSlug {
		return `
// skuPrefix is the leading half of a generated SKU.
//
// Falls back to the id where a row has no slug yet, which happens for anything
// written before the slug hook existed.
func skuPrefix(slug, id string) string {
	if slug == "" {
		return shortID(id)
	}
	return strings.ToUpper(slug)
}

// shortID is a readable stub of an id, for a prefix that has nothing better.
func shortID(id string) string {
	if len(id) > 6 {
		id = id[:6]
	}
	return strings.ToUpper(id)
}
`
	}
	return `
// shortID is a readable stub of an id, used as a SKU prefix because this
// resource has no slug to build one from.
func shortID(id string) string {
	if len(id) > 6 {
		id = id[:6]
	}
	return strings.ToUpper(id)
}
`
}
