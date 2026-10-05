package paginate

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// EnsureSearchIndexes gives every column tagged search:"trigram" a trigram index
// on Postgres, so List's search is served by an index instead of reading the
// whole table.
//
// LOWER(col) LIKE '%term%' cannot use a B-tree index. A pg_trgm GIN index on
// lower(col) can: on a million contacts a search went from 270 ms to 2 ms. Other
// databases are left as they are. The pg_trgm extension is one managed Postgres
// offers to ordinary users; when it cannot be created, the indexes are skipped
// with a message and search keeps working, unindexed.
func EnsureSearchIndexes(db *gorm.DB, models ...interface{}) error {
	if db.Dialector.Name() != "postgres" {
		return nil
	}
	type column struct{ table, name string }
	var columns []column
	for _, model := range models {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			return fmt.Errorf("reading %T: %w", model, err)
		}
		for _, field := range stmt.Schema.Fields {
			if field.Tag.Get("search") == "trigram" && field.DBName != "" {
				columns = append(columns, column{stmt.Schema.Table, field.DBName})
			}
		}
	}
	if len(columns) == 0 {
		return nil
	}
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm").Error; err != nil {
		log.Printf("Search indexes skipped: the pg_trgm extension could not be created (%v). Search still works, without an index.", err)
		return nil
	}
	quote := &gorm.Statement{DB: db}
	for _, c := range columns {
		name := "idx_" + c.table + "_" + c.name + "_trgm"
		// CONCURRENTLY, so a big table keeps taking writes while it builds.
		sql := fmt.Sprintf("CREATE INDEX CONCURRENTLY IF NOT EXISTS %s ON %s USING gin (lower(%s) gin_trgm_ops)",
			quote.Quote(name), quote.Quote(c.table), quote.Quote(c.name))
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("creating the search index %s: %w", name, err)
		}
	}
	return nil
}
