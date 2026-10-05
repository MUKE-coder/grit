package scaffold

import (
	"strings"
	"testing"
)

// grit migrate --fresh has to work on the three databases a project can use.
//
// DropAll asked pg_tables and dropped with CASCADE, both Postgres only, so
// `grit migrate --fresh` failed outright on a SQLite project ("no such table:
// pg_tables") and on MySQL. The dialect switch it needed was two functions away
// in the same generated package, in TableCount, with a comment about this exact
// mistake having been made once already.
func TestDropAllAsksEachDialectItsOwnQuestion(t *testing.T) {
	src := apiMigrateGo()

	// The three listings, each in its own branch.
	for dialect, query := range map[string]string{
		"postgres": "SELECT tablename FROM pg_tables WHERE schemaname = current_schema()",
		"mysql":    "SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()",
		"sqlite":   "SELECT name FROM sqlite_master WHERE type = 'table'",
	} {
		if !strings.Contains(src, `case "`+dialect+`":`) {
			t.Errorf("DropAll has no branch for %s", dialect)
		}
		if !strings.Contains(src, query) {
			t.Errorf("DropAll does not ask %s its own question (%q)", dialect, query)
		}
	}

	// 'public' hardcoded would drop nothing in a project with its own
	// search_path, and report success while doing it.
	if strings.Contains(src, "schemaname = 'public'") {
		t.Error("DropAll hardcodes the public schema instead of current_schema()")
	}

	// CASCADE is Postgres-only: MySQL rejects the keyword outright.
	if strings.Count(src, "CASCADE") > 0 && !strings.Contains(src, `if db.Dialector.Name() == "postgres" {`) {
		t.Error("DropAll uses CASCADE without confining it to Postgres")
	}

	// The other two need foreign keys off first, because dropping in an
	// arbitrary order hits a table another one still references.
	for _, off := range []string{"SET FOREIGN_KEY_CHECKS = 0", "PRAGMA foreign_keys = OFF"} {
		if !strings.Contains(src, off) {
			t.Errorf("DropAll does not disable foreign keys with %q", off)
		}
	}

	// MySQL reads a double-quoted name as a string literal unless ANSI_QUOTES
	// is set, so DROP TABLE "users" is a syntax error there.
	if !strings.Contains(src, "`%s`") {
		t.Error("DropAll does not backtick-quote table names for MySQL")
	}

	// An unknown dialect says so rather than silently dropping nothing.
	if !strings.Contains(src, "does not know how to drop tables on") {
		t.Error("DropAll does not report a dialect it cannot handle")
	}

	mustFormatGo(t, "database/migrate.go", src)
}
