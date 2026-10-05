package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
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
