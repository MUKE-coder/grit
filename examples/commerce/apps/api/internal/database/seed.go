package database

import (
	"fmt"

	"gorm.io/gorm"
)

// Seed runs every seeder. Seeders live in their own <name>_seeder.go files in
// this package — edit those to change the seed data, or run
// "grit generate seeder <Resource>" to add a new one.
func Seed(db *gorm.DB) error {
	if err := SeedUsers(db); err != nil {
		return fmt.Errorf("seeding users: %w", err)
	}

	if err := SeedAPIKeys(db); err != nil {
		return fmt.Errorf("seeding api keys: %w", err)
	}

	if err := SeedBlogs(db); err != nil {
		return fmt.Errorf("seeding blogs: %w", err)
	}

	if err := SeedCollections(db); err != nil {
		return fmt.Errorf("seeding collections: %w", err)
	}

	if err := SeedProducts(db); err != nil {
		return fmt.Errorf("seeding products: %w", err)
	}

	if err := SeedProductVariants(db); err != nil {
		return fmt.Errorf("seeding product variants: %w", err)
	}

	if err := SeedPages(db); err != nil {
		return fmt.Errorf("seeding pages: %w", err)
	}

	// grit:seeders

	return nil
}
