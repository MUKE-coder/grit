package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeMigrateSeedFiles(root string, opts Options) error {
	apiRoot := filepath.Join(root, "apps", "api")
	module := opts.ProjectName + "/apps/api"

	files := map[string]string{
		filepath.Join(apiRoot, "cmd", "migrate", "main.go"):          apiMigrateMainGo(),
		filepath.Join(apiRoot, "cmd", "seed", "main.go"):             apiSeedMainGo(),
		filepath.Join(apiRoot, "internal", "database", "seed.go"):    apiSeedGo(),
		filepath.Join(apiRoot, "internal", "database", "migrate.go"): apiMigrateGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// apiMigrateMainGo returns the cmd/migrate/main.go entrypoint.
func apiMigrateMainGo() string {
	return `package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"{{MODULE}}/internal/config"
	"{{MODULE}}/internal/database"
	"{{MODULE}}/internal/models"
)

func main() {
	fresh := flag.Bool("fresh", false, "Drop all tables before migrating")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if *fresh {
		fmt.Println("Dropping all tables...")
		if err := database.DropAll(db); err != nil {
			log.Fatalf("Failed to drop tables: %v", err)
		}
		fmt.Println("All tables dropped.")
	}

	fmt.Println("Running migrations...")
	if err := models.AutoMigrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migrations completed successfully.")
	os.Exit(0)
}
`
}

// apiSeedMainGo returns the cmd/seed/main.go entrypoint.
func apiSeedMainGo() string {
	return `package main

import (
	"fmt"
	"log"
	"os"

	"{{MODULE}}/internal/config"
	"{{MODULE}}/internal/database"
	"{{MODULE}}/internal/models"
)

func main() {
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
	if err := models.AutoMigrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Seeding database...")
	if err := database.Seed(db); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	fmt.Println("Database seeded successfully.")
	os.Exit(0)
}
`
}

// apiMigrateGo returns the database/migrate.go helper with DropAll.
func apiMigrateGo() string {
	return `package database

import (
	"fmt"

	"gorm.io/gorm"
)

// DropAll drops all tables in the database.
// Used by the migrate --fresh command.
func DropAll(db *gorm.DB) error {
	// Get all table names
	var tables []string
	if err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}

	if len(tables) == 0 {
		return nil
	}

	// Disable foreign key checks and drop all tables
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %q CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	return nil
}
`
}

// apiSeedGo returns the database/seed.go with the default seeder.
func apiSeedGo() string {
	return `package database

import (
	"fmt"
	"log"

	"{{MODULE}}/internal/models"
	"gorm.io/gorm"
)

// Seed populates the database with initial data.
// Add your seeders to this function.
func Seed(db *gorm.DB) error {
	if err := seedAdminUser(db); err != nil {
		return fmt.Errorf("seeding admin user: %w", err)
	}

	if err := seedDemoUsers(db); err != nil {
		return fmt.Errorf("seeding demo users: %w", err)
	}

	// grit:seeders

	return nil
}

// seedAdminUser creates the default admin account.
func seedAdminUser(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Where("email = ?", "admin@example.com").Count(&count)
	if count > 0 {
		log.Println("Admin user already exists, skipping...")
		return nil
	}

	admin := models.User{
		Name:   "Admin",
		Email:  "admin@example.com",
		Password: "password",
		Role:   "admin",
		Active: true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("creating admin user: %w", err)
	}

	log.Println("Created admin user: admin@example.com / password")
	return nil
}

// seedDemoUsers creates sample user accounts for development.
func seedDemoUsers(db *gorm.DB) error {
	users := []models.User{
		{Name: "Jane Cooper", Email: "jane@example.com", Password: "password", Role: "editor", Active: true},
		{Name: "Robert Fox", Email: "robert@example.com", Password: "password", Role: "user", Active: true},
		{Name: "Emily Davis", Email: "emily@example.com", Password: "password", Role: "user", Active: true},
		{Name: "Michael Chen", Email: "michael@example.com", Password: "password", Role: "user", Active: false},
	}

	for _, u := range users {
		var count int64
		db.Model(&models.User{}).Where("email = ?", u.Email).Count(&count)
		if count > 0 {
			continue
		}

		if err := db.Create(&u).Error; err != nil {
			log.Printf("Warning: failed to create user %s: %v", u.Email, err)
			continue
		}
		log.Printf("Created user: %s / password", u.Email)
	}

	return nil
}
`
}
