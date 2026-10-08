package database

import (
	"fmt"

	"gorm.io/gorm"
)

// DropAll drops every table in the database. Used by grit migrate --fresh.
//
// Three dialects need three different questions and three different DROPs.
// This asked pg_tables and dropped with CASCADE, both of which are Postgres
// only, so grit migrate --fresh failed outright on a SQLite project ("no such
// table: pg_tables") and on MySQL. The same mistake had already been made and
// fixed in TableCount, in dialect.go, two functions away.
func DropAll(db *gorm.DB) error {
	var list, quote string
	switch name := db.Dialector.Name(); name {
	case "postgres":
		// current_schema() rather than 'public', so a project with its own
		// search_path drops its own tables and not nothing.
		list = "SELECT tablename FROM pg_tables WHERE schemaname = current_schema()"
		quote = `"%s"`
	case "mysql":
		list = "SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()"
		// Backticks: MySQL reads a double-quoted name as a string literal
		// unless ANSI_QUOTES is set, and DROP TABLE "users" is a syntax error.
		quote = "`%s`"
	case "sqlite":
		list = "SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'"
		quote = `"%s"`
	default:
		return fmt.Errorf("grit migrate --fresh does not know how to drop tables on %s", name)
	}

	var tables []string
	if err := db.Raw(list).Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}
	if len(tables) == 0 {
		return nil
	}

	// Foreign keys have to come off first. Dropping in an arbitrary order hits
	// a table another one still references, and on SQLite that is an error
	// rather than something the database works around.
	//
	// Postgres does not need this because CASCADE is per-statement, which is
	// why it is still there and only there: MySQL rejects the keyword and
	// SQLite ignores it, so neither can be given it.
	switch db.Dialector.Name() {
	case "mysql":
		if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
			return fmt.Errorf("failed to disable foreign key checks: %w", err)
		}
		defer db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	case "sqlite":
		if err := db.Exec("PRAGMA foreign_keys = OFF").Error; err != nil {
			return fmt.Errorf("failed to disable foreign key enforcement: %w", err)
		}
		defer db.Exec("PRAGMA foreign_keys = ON")
	}

	cascade := ""
	if db.Dialector.Name() == "postgres" {
		cascade = " CASCADE"
	}

	for _, table := range tables {
		stmt := "DROP TABLE IF EXISTS " + fmt.Sprintf(quote, table) + cascade
		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	return nil
}
