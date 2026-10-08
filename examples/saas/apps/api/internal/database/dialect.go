package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// SupportsReturning reports whether the connected dialect can hand back the
// written row from an INSERT or UPDATE.
//
// Postgres and SQLite can. MySQL cannot, and this is the important part: it
// does not error when asked. The write succeeds, the RETURNING clause is
// dropped, and the struct comes back with every database-assigned default
// still at its zero value. A handler that skipped its reload on the strength
// of RETURNING would then answer 201 with a half-empty record.
func SupportsReturning(db *gorm.DB) bool {
	switch db.Dialector.Name() {
	case "postgres", "sqlite":
		return true
	default:
		return false
	}
}

// Write returns a session for a single-statement write: no wrapping
// transaction, and RETURNING where the dialect has it.
//
// Skipping the transaction is safe only because the caller has already
// established there is exactly one statement. The generator decides that from
// the resource definition, where it can be known rather than assumed.
func Write(db *gorm.DB) *gorm.DB {
	tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})
	if SupportsReturning(db) {
		tx = tx.Clauses(clause.Returning{})
	}
	return tx
}

// TableCount returns how many tables the connected database holds, or 0 when
// the dialect cannot be asked. Purely informational: it feeds the "tables: N"
// figure on the health card.
//
// Three dialects need three different questions. information_schema exists on
// Postgres and MySQL but not on SQLite, and the two that have it spell "this
// database" differently. Asking the Postgres question everywhere logged a red
// SQL error on every health poll of a SQLite project, and quietly returned 0
// on MySQL, where current_schema() does not exist either.
//
// Errors are swallowed with the logger silenced, because a missing tooltip
// figure is not worth a stack of scary log lines on a healthy server.
func TableCount(db *gorm.DB) int {
	var query string
	switch db.Dialector.Name() {
	case "postgres":
		query = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = current_schema()"
	case "mysql":
		query = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE()"
	case "sqlite":
		query = "SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'"
	default:
		return 0
	}

	var count int
	quiet := db.Session(&gorm.Session{Logger: db.Logger.LogMode(logger.Silent)})
	if err := quiet.Raw(query).Scan(&count).Error; err != nil {
		return 0
	}
	return count
}
