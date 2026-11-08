package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// Every bookkeeping table the migrate package owns is excluded from a schema
// snapshot.
//
// A run is recorded as the baseline, the one that built the schema, when the
// snapshot taken before it was empty:
//
//	baseline := len(before) == 0
//
// and Snapshot skips the package's own tables so that its bookkeeping does not
// count as the application's schema. v3.480.0 added a third such table,
// grit_sql_migrations, created by EnsureSQLHistory before that snapshot is
// taken, and did not add it to historyTables. So a fresh database no longer
// looked empty, no run was ever the baseline, and `grit migrate down` offered
// to roll back the run that created all 42 tables instead of refusing. On a
// fresh project the recorded run had baseline = 0 and --down --yes dropped the
// schema.
//
// The list is built from the types now, which is what makes it hard to forget.
// This test is the backstop for that: it reads every TableName in the package
// and insists the exclusion list knows about it.
var tableNameDecl = regexp.MustCompile(`func \((\w+)\) TableName\(\) string \{ return "(grit_\w+)" \}`)

func TestSnapshotExcludesEveryBookkeepingTable(t *testing.T) {
	// Both files that declare a table in package migrate.
	sources := map[string]string{
		"internal/migrate/history.go": apiMigrateHistoryGo(),
		"internal/migrate/sql.go":     apiMigrateSQLGo(),
	}

	declared := map[string]string{} // table name -> the type that declares it
	for name, src := range sources {
		for _, m := range tableNameDecl.FindAllStringSubmatch(src, -1) {
			declared[m[2]] = m[1] + " in " + name
		}
	}
	if len(declared) < 3 {
		t.Fatalf("found %d bookkeeping tables, expected at least 3 (the two run tables and "+
			"the SQL history); the regex has stopped matching", len(declared))
	}

	history := apiMigrateHistoryGo()
	start := strings.Index(history, "var historyTables = map[string]bool{")
	if start == -1 {
		t.Fatal("history.go has no historyTables")
	}
	// To the closing brace of the map, which is the one at the start of a line.
	// The first "}" after start is inside Run{}.TableName().
	end := strings.Index(history[start:], "\n}")
	if end == -1 {
		t.Fatal("historyTables is not closed")
	}
	block := history[start : start+end]

	for table, where := range declared {
		// The list names the type, not the string, so look for either: a
		// future change back to literals should still pass.
		typeName := strings.Fields(where)[0]
		if !strings.Contains(block, typeName+"{}.TableName()") && !strings.Contains(block, `"`+table+`"`) {
			t.Errorf("%s is not excluded from a schema snapshot (%s).\n\n"+
				"    Snapshot counts it as part of the application's schema, so the "+
				"snapshot taken\n"+
				"    before the first migration is not empty, no run is recorded as the "+
				"baseline, and\n"+
				"    grit migrate down will roll back the run that built the database "+
				"instead of\n"+
				"    refusing. Add it to historyTables in history.go.", table, where)
		}
	}
}

// And the decision itself, which is the line the bug was in.
func TestBaselineIsAnEmptySnapshot(t *testing.T) {
	main := apiMigrateMainGo()
	if !strings.Contains(main, "baseline := len(before) == 0") {
		t.Error("cmd/migrate no longer decides the baseline from an empty snapshot; " +
			"if that has changed, TestSnapshotExcludesEveryBookkeepingTable is testing " +
			"the wrong thing")
	}
	// The snapshot has to be taken before the migration, or it is a snapshot of
	// the schema the migration just created and nothing is ever a baseline.
	snapshot := strings.Index(main, "before, err := migrate.Snapshot(db)")
	migrateCall := strings.Index(main, "if err := models.Migrate(db); err != nil {")
	if snapshot == -1 || migrateCall == -1 {
		t.Fatal("cmd/migrate does not take a snapshot and then migrate")
	}
	if snapshot > migrateCall {
		t.Error("the before-snapshot is taken after models.Migrate, so it is not a before-snapshot")
	}
}
