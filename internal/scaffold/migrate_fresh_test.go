package scaffold

import (
	"strings"
	"testing"
)

// grit#189: `grit migrate --fresh` dropped grit_sql_migrations and then read it.
//
// The ordering is the whole bug, and it is invisible in either file alone.
// cmd/migrate/main.go calls EnsureSQLHistory near the top, correctly, because
// everything below needs the table. Then the --fresh branch drops every table,
// including that one, and the SQL half is read a few lines further down:
//
//	ERROR: relation "grit_sql_migrations" does not exist (SQLSTATE 42P01)
//
// The run then skips every checked-in pair, on a database that otherwise looks
// correctly migrated. Those pairs are the ones AutoMigrate cannot express, so
// what goes missing is triggers, backfills and extensions: including the
// append-only battery's trigger, which leaves a table unprotected while the
// app still believes it is not.

// TestFreshResetsBothHistories reads the one thing that matters: that the
// --fresh branch puts back both of the tables it just dropped.
func TestFreshResetsBothHistories(t *testing.T) {
	src := apiMigrateMainGo()

	start := strings.Index(src, "if *fresh {")
	if start == -1 {
		t.Fatal("cmd/migrate/main.go has no --fresh branch")
	}
	end := strings.Index(src[start:], "\n\tbefore, err := migrate.Snapshot(db)")
	if end == -1 {
		t.Fatal("could not find the end of the --fresh branch")
	}
	branch := src[start : start+end]

	for _, want := range []string{
		"database.DropAll(db)",
		"migrate.Reset(db)",
		"migrate.ResetSQL(db)",
	} {
		if !strings.Contains(branch, want) {
			t.Errorf("the --fresh branch does not call %s, so the history it dropped "+
				"is read back missing and every checked-in SQL pair is skipped:\n%s", want, branch)
		}
	}

	// And the order: both resets after the drop, or they are undone by it.
	drop := strings.Index(branch, "database.DropAll(db)")
	for _, after := range []string{"migrate.Reset(db)", "migrate.ResetSQL(db)"} {
		if at := strings.Index(branch, after); at < drop {
			t.Errorf("%s runs before DropAll, which drops the table it just made", after)
		}
	}
}

// TestResetSQLPutsTheTableBack: dropping without ensuring would leave --fresh
// failing in a different place.
func TestResetSQLPutsTheTableBack(t *testing.T) {
	src := apiMigrateSQLGo()

	start := strings.Index(src, "func ResetSQL(db *gorm.DB) error {")
	if start == -1 {
		t.Fatal("internal/migrate/sql.go has no ResetSQL")
	}
	body := src[start:]
	if end := strings.Index(body, "\n}\n"); end != -1 {
		body = body[:end]
	}
	if !strings.Contains(body, "DropTable(&AppliedSQLMigration{})") {
		t.Error("ResetSQL does not drop the history, so --fresh keeps a history of migrations whose schema is gone")
	}
	if !strings.Contains(body, "EnsureSQLHistory(db)") {
		t.Error("ResetSQL does not recreate the history, so the next read of it fails")
	}
}

// TestFreshAsksBeforeDropping: it is the most destructive thing the command
// does and was the only one that did it without asking, while `migrate down`
// has a confirmation.
func TestFreshAsksBeforeDropping(t *testing.T) {
	src := apiMigrateMainGo()

	start := strings.Index(src, "if *fresh {")
	if start == -1 {
		t.Fatal("no --fresh branch")
	}
	branch := src[start:]
	if end := strings.Index(branch, "database.DropAll(db)"); end != -1 {
		branch = branch[:end]
	}
	if !strings.Contains(branch, "confirmed()") {
		t.Errorf("--fresh drops every table without asking:\n%s", branch)
	}
	if !strings.Contains(branch, "!*yes") {
		t.Error("--fresh has no way to skip the confirmation, so it cannot run in CI")
	}
}
