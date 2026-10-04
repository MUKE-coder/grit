package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"library/apps/api/internal/config"
	"library/apps/api/internal/database"
	"library/apps/api/internal/models"
)

// With no flags it runs every seeder, as it always has. With -resource and
// -count it tops one resource's table up to that many rows ("grit seed Contact
// --count 1000000" passes them), inserting in batches and resuming where an
// earlier run stopped.
func main() {
	resource := flag.String("resource", "", "seed only this resource, for example Contact")
	count := flag.Int64("count", 0, "with -resource: the number of rows the table should hold")
	flag.Parse()
	if *resource != "" && *count <= 0 {
		log.Fatal("-count must be a positive number of rows when -resource is given")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Ensure tables exist before seeding
	fmt.Println("Running migrations...")
	if err := models.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	if *resource != "" {
		fmt.Printf("Seeding %s up to %d rows...\n", *resource, *count)
		if err := database.SeedOne(db, *resource, *count); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
		fmt.Println("Done.")
		os.Exit(0)
	}

	fmt.Println("Seeding database...")
	if err := database.Seed(db); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	fmt.Println("Database seeded successfully.")
	os.Exit(0)
}
