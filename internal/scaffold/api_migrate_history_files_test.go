package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Grit has no migration files, so a rollback can only come from a record of what
// each run actually changed. These tests hold the parts of that record that are
// easy to break from the outside: the tables it writes to, what it refuses, and
// the fact that it reaches existing projects at all.
func TestMigrateHistoryPackageRecordsAndReverses(t *testing.T) {
	src := apiMigrateHistoryGo()
	durableMustFormat(t, "internal/migrate/history.go", src)
	durableMustContain(t, "internal/migrate/history.go", src,
		`func (Run) TableName() string { return "grit_migrations" }`,
		`func (RunChange) TableName() string { return "grit_migration_changes" }`,
		"func Snapshot(db *gorm.DB) (Schema, error)",
		"func Diff(before, after Schema) []Change",
		"func Rollback(db *gorm.DB, run Run, changes []Change, dryRun bool) ([]string, error)",
		"func Undone(db *gorm.DB, steps int) ([]Run, error)",
		// Indexes, then columns, then tables: a column cannot be dropped while an
		// index covers it, and a table must go last or the drops inside it fail.
		`for _, kind := range []string{KindIndex, KindColumn, KindTable} {`,
		// Every drop is guarded, so a rollback that failed halfway can be re-run.
		"func alreadyGone(db *gorm.DB, change Change) bool",
	)

	// The history describes the schema; it is not part of it. If a snapshot could
	// see these tables, the first rollback would drop the record of itself.
	durableMustContain(t, "internal/migrate/history.go", src,
		"var historyTables = map[string]bool{",
		"if historyTables[name] || internalTable(name) {",
		// SQLite reports sqlite_sequence, created by the history's own
		// AUTOINCREMENT id. Counted, it makes every database look non-empty, so
		// no first run is a baseline and a rollback offers to drop everything.
		"func internalTable(name string) bool {",
	)

	// Postgres, MySQL and SQLite do not spell these the same way.
	durableMustContain(t, "internal/migrate/history.go", src,
		`if db.Dialector.Name() == "mysql" {`,
		`return "DROP TABLE IF EXISTS " + quote(db, change.OnTable) + " CASCADE"`,
	)
}

func TestMigrateHistoryTestsShipIntoTheProject(t *testing.T) {
	src := apiMigrateHistoryTestGo()
	durableMustFormat(t, "internal/migrate/history_test.go", src)
	durableMustContain(t, "internal/migrate/history_test.go", src,
		"func TestRollbackDropsWhatTheRunAdded(t *testing.T)",
		"func TestRollbackSkipsWhatIsAlreadyGone(t *testing.T)",
		"func TestStatementsUndoInReverseOrder(t *testing.T)",
		"func TestSnapshotIgnoresTheHistoryItself(t *testing.T)",
	)
}

// cmd/migrate is the command around the package: it records a run, and it is
// where the two refusals live.
func TestMigrateCommandHasStatusAndDown(t *testing.T) {
	src := apiMigrateMainGo()
	durableMustFormat(t, "cmd/migrate/main.go", src)
	durableMustContain(t, "cmd/migrate/main.go", src,
		`status := flag.Bool("status"`,
		`down := flag.Bool("down"`,
		`steps := flag.Int("steps"`,
		`dryRun := flag.Bool("dry-run"`,
		`yes := flag.Bool("yes"`,
		"if err := migrate.EnsureHistory(db); err != nil {",
		"changes := migrate.Diff(before, after)",
		"baseline := len(before) == 0",
		// The first run on an empty database built the schema. Undoing that is
		// dropping the database, which is what --fresh is for.
		"is the baseline, the run that built this schema",
		// Dropping a column loses its data, so nothing runs unasked.
		"if !yes && !confirmed() {",
	)

	// --fresh threw away the schema the history described.
	durableMustContain(t, "cmd/migrate/main.go", src,
		"if err := migrate.Reset(db); err != nil {",
	)

	// A migration that worked must not report failure because the bookkeeping
	// after it did not.
	durableMustContain(t, "cmd/migrate/main.go", src,
		"Migrations completed, but recording the run failed",
	)
}

// The package has to arrive in existing projects, and in the same upgrade as the
// cmd/migrate that imports it. Split across the two sets, grit upgrade would
// write a main.go importing a package that never came, and the project would
// stop compiling.
func TestMigrateHistoryTravelsOnUpgrade(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "app", Architecture: ArchAPI}

	if err := writeFrameworkOwnedFiles(root, opts); err != nil {
		t.Fatal(err)
	}
	if err := writeMigrateSeedFiles(root, opts); err != nil {
		t.Fatal(err)
	}

	apiRoot := opts.APIRoot(root)
	for _, rel := range []string{
		"internal/migrate/history.go",
		"internal/migrate/history_test.go",
		"cmd/migrate/main.go",
	} {
		if _, err := os.Stat(filepath.Join(apiRoot, filepath.FromSlash(rel))); err != nil {
			t.Errorf("%s is not written by the upgrade path: %v", rel, err)
		}
	}

	main, err := os.ReadFile(filepath.Join(apiRoot, filepath.FromSlash("cmd/migrate/main.go")))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(main), "/internal/migrate\"") {
		t.Error("cmd/migrate/main.go does not import the migrate package it depends on")
	}
	if strings.Contains(string(main), "{{MODULE}}") {
		t.Error("cmd/migrate/main.go still has an unreplaced module placeholder")
	}
}
