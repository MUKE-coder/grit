package database

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
