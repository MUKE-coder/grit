package migrate

import (
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Two shapes of the same table: the second adds a column. No struct tags, so
// the index below is created explicitly and the test does not depend on how a
// tag is spelled.
type widget struct {
	ID   uint
	Name string
}

func (widget) TableName() string { return "widgets" }

type widgetWithColour struct {
	ID     uint
	Name   string
	Colour string
}

func (widgetWithColour) TableName() string { return "widgets" }

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := EnsureHistory(db); err != nil {
		t.Fatalf("ensure history: %v", err)
	}
	return db
}

// The whole promise in one test: a run that adds a column and an index is
// undone by dropping exactly those, and the table an earlier run created is
// still standing afterwards.
func TestRollbackDropsWhatTheRunAdded(t *testing.T) {
	db := testDB(t)

	before, err := Snapshot(db)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if len(before) != 0 {
		t.Fatalf("a database holding only the history should snapshot as empty, got %v", before)
	}

	if err := db.AutoMigrate(&widget{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	after, err := Snapshot(db)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	changes := Diff(before, after)
	if len(changes) != 1 || changes[0].Kind != KindTable || changes[0].OnTable != "widgets" {
		t.Fatalf("a new table should be one change, not its every column: %v", changes)
	}
	baseline, err := Record(db, changes, "v0.0.0-test", true)
	if err != nil {
		t.Fatalf("record: %v", err)
	}

	// Second run: a column, and an index on it.
	before = after
	if err := db.AutoMigrate(&widgetWithColour{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	if err := db.Exec("CREATE INDEX idx_widgets_colour ON widgets(colour)").Error; err != nil {
		t.Fatalf("create index: %v", err)
	}
	after, err = Snapshot(db)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	changes = Diff(before, after)

	var sawColumn, sawIndex bool
	for _, change := range changes {
		if change.Kind == KindColumn && change.Name == "colour" {
			sawColumn = true
		}
		if change.Kind == KindIndex && change.Name == "idx_widgets_colour" {
			sawIndex = true
		}
	}
	if !sawColumn || !sawIndex {
		t.Fatalf("expected the new column and its index in the diff, got %v", changes)
	}

	run, err := Record(db, changes, "v0.0.0-test", false)
	if err != nil {
		t.Fatalf("record: %v", err)
	}
	stored, err := ChangesOf(db, run.ID)
	if err != nil {
		t.Fatalf("changes of: %v", err)
	}
	if len(stored) != len(changes) {
		t.Fatalf("recorded %d changes, read back %d", len(changes), len(stored))
	}

	// A dry run reports the statements and touches nothing.
	statements, err := Rollback(db, run, stored, true)
	if err != nil {
		t.Fatalf("dry run: %v", err)
	}
	if len(statements) == 0 {
		t.Fatal("a dry run should still report the statements it would run")
	}
	if !db.Migrator().HasColumn("widgets", "colour") {
		t.Fatal("a dry run dropped a column")
	}

	applied, err := Rollback(db, run, stored, false)
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if len(applied) == 0 {
		t.Fatal("rollback reported no statements")
	}
	if db.Migrator().HasColumn("widgets", "colour") {
		t.Fatal("the column the run added survived the rollback")
	}
	if !db.Migrator().HasTable("widgets") {
		t.Fatal("rolling back the second run dropped the table the first run created")
	}

	pending, err := Undone(db, 0)
	if err != nil {
		t.Fatalf("undone: %v", err)
	}
	var stillOffered, baselineOffered bool
	for _, candidate := range pending {
		if candidate.ID == run.ID {
			stillOffered = true
		}
		if candidate.ID == baseline.ID {
			baselineOffered = true
		}
	}
	if stillOffered {
		t.Fatal("a rolled back run is still offered for rollback")
	}
	if !baselineOffered {
		t.Fatal("the baseline run should still be listed; refusing it is the command's job")
	}
}

// A rollback that failed halfway can be run again: anything already dropped is
// skipped rather than retried into an error.
func TestRollbackSkipsWhatIsAlreadyGone(t *testing.T) {
	db := testDB(t)
	if err := db.AutoMigrate(&widgetWithColour{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}

	changes := []Change{
		{Kind: KindColumn, OnTable: "widgets", Name: "colour"},
		{Kind: KindColumn, OnTable: "widgets", Name: "dropped_by_hand"},
		{Kind: KindIndex, OnTable: "widgets", Name: "idx_that_never_existed"},
	}
	run, err := Record(db, changes, "v0.0.0-test", false)
	if err != nil {
		t.Fatalf("record: %v", err)
	}

	applied, err := Rollback(db, run, changes, false)
	if err != nil {
		t.Fatalf("a rollback must skip a change that is already undone: %v", err)
	}
	if len(applied) != 1 {
		t.Fatalf("expected one statement for the one column that exists, got %v", applied)
	}
}

// Order matters: an index goes before the column it covers, and a column before
// the table it sits in.
func TestStatementsUndoInReverseOrder(t *testing.T) {
	db := testDB(t)
	statements := Statements(db, []Change{
		{Kind: KindTable, OnTable: "widgets"},
		{Kind: KindColumn, OnTable: "gadgets", Name: "colour"},
		{Kind: KindIndex, OnTable: "gadgets", Name: "idx_gadgets_colour"},
	})
	if len(statements) != 3 {
		t.Fatalf("expected three statements, got %v", statements)
	}
	if !strings.Contains(statements[0], "idx_gadgets_colour") {
		t.Fatalf("the index should be dropped first, got %q", statements[0])
	}
	if !strings.Contains(statements[1], "DROP COLUMN") {
		t.Fatalf("the column should be dropped second, got %q", statements[1])
	}
	if !strings.Contains(statements[2], "DROP TABLE") {
		t.Fatalf("the table should be dropped last, got %q", statements[2])
	}
}

// The history is not part of the schema it records, so a rollback can never
// drop the record of itself.
func TestSnapshotIgnoresTheHistoryItself(t *testing.T) {
	db := testDB(t)
	schema, err := Snapshot(db)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	for name := range schema {
		if historyTables[name] {
			t.Fatalf("%s belongs to the history, not to the schema it records", name)
		}
		// sqlite_sequence is here because the history's id is an AUTOINCREMENT. A
		// snapshot that counts it is never empty, so no run is ever a baseline.
		if internalTable(name) {
			t.Fatalf("%s belongs to the database, not to the schema it records", name)
		}
	}
}

// Runs recorded back to back, faster than the clock's millisecond, each get an
// ID of their own, in order. Two such runs once failed on the primary key, which
// is how this test's neighbour failed in CI.
func TestRunsRecordedBackToBackGetDistinctIDs(t *testing.T) {
	db := testDB(t)
	previous := ""
	for i := 0; i < 50; i++ {
		run, err := Record(db, nil, "v0.0.0-test", false)
		if err != nil {
			t.Fatalf("run %d: %v", i, err)
		}
		if run.ID <= previous {
			t.Fatalf("run %d has id %s, not after %s", i, run.ID, previous)
		}
		previous = run.ID
	}
}
