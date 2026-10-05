package database

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/services"
)

// SeedProductVariants gives the first few products a real matrix.
//
// A few rather than all of them, on purpose. A catalogue where every row has
// sixteen combinations is a catalogue where nothing tests the plain case, and
// the plain case is most of a real shop: a product with no options at all
// still has to render, price and sell.
func SeedProductVariants(db *gorm.DB) error {
	// How many products get options attached.
	const withVariants = 6

	var existing int64
	db.Model(&models.ProductVariant{}).Count(&existing)
	if existing > 0 {
		log.Println("Product variants already seeded, skipping...")
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
		log.Println("No options to attach, skipping product variants...")
		return nil
	}

	var rows []models.Product
	if err := db.Order("created_at asc").Limit(withVariants).Find(&rows).Error; err != nil {
		return fmt.Errorf("loading products to attach variants to: %w", err)
	}
	if len(rows) == 0 {
		log.Println("No products exist yet, skipping variants. Seed Product first, then run grit seed again.")
		return nil
	}

	variants := services.NewProductVariantService(db)
	seeded := 0

	for _, row := range rows {
		for position, option := range options {
			link := models.ProductOption{
				ProductID: row.ID,
				OptionID:  option.ID,
				Position:  position,
			}
			err := db.Where("product_id = ? AND option_id = ?", row.ID, option.ID).
				FirstOrCreate(&link).Error
			if err != nil {
				return fmt.Errorf("attaching %s to a product: %w", option.Name, err)
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
				"sku":   productSkuPrefix(row.Handle, row.ID) + "-" + productValueSuffix(variant.OptionValues),
				"stock": 4 + (i*7)%37,
			}
			// One combination in seven is out of stock. The disabled swatch is
			// most of the work on a product page and the hardest state to
			// remember to build, so the seed puts it on screen unasked.
			if i%7 == 3 {
				updates["stock"] = 0
			}
			if err := db.Model(&models.ProductVariant{}).
				Where("id = ?", variant.ID).Updates(updates).Error; err != nil {
				return fmt.Errorf("filling in a variant: %w", err)
			}
			seeded++
		}
	}

	log.Printf("Seeded %d product variants across %d products", seeded, len(rows))
	return nil
}

// productValueSuffix builds the readable half of a SKU from the values that define the
// combination: BLACK-XL rather than a serial number nobody can check against a
// shelf.
//
// From the label and not the slug. Slugs carry a uniqueness suffix, so the slug
// half of a SKU reads BLACK-ADF9C36E-XL-93661012, which is unique, unreadable,
// and no use at all to the person holding the box.
func productValueSuffix(values []models.OptionValue) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, strings.ToUpper(strings.ReplaceAll(value.Label, " ", "-")))
	}
	return strings.Join(parts, "-")
}

// productSkuPrefix is the leading half of a generated SKU.
//
// Falls back to the id where a row has no slug yet, which happens for anything
// written before the slug hook existed.
func productSkuPrefix(slug, id string) string {
	if slug == "" {
		return productShortID(id)
	}
	return strings.ToUpper(slug)
}

// productShortID is a readable stub of an id, for a prefix that has nothing better.
func productShortID(id string) string {
	if len(id) > 6 {
		id = id[:6]
	}
	return strings.ToUpper(id)
}
