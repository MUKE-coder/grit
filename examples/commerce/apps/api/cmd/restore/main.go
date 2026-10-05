package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"commerce/apps/api/internal/backup"
	"commerce/apps/api/internal/config"
	"commerce/apps/api/internal/database"
	"commerce/apps/api/internal/models"
)

// Replays a backup archive into the configured database. Runs migrations first
// (the archive carries data, not schema), then executes dump.sql in one
// transaction: every row lands, or none does.
func main() {
	migrate := flag.Bool("migrate", true, "Run migrations before restoring")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("usage: restore [--migrate=false] <backup.zip>")
	}
	path := flag.Arg(0)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if *migrate {
		fmt.Println("Running migrations...")
		if err := models.Migrate(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	}

	fmt.Printf("Restoring %s ...\n", path)
	man, err := backup.Restore(db, path)
	if err != nil {
		log.Fatalf("Restore failed: %v", err)
	}

	fmt.Printf("Restored %d tables, %d rows (archive generated %s)\n",
		len(man.Tables), man.TotalRows, man.GeneratedAt.Format(time.RFC3339))
}
