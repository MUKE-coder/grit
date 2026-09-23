package scaffold

import (
	"strings"
	"testing"
)

// A SQLite file with no busy timeout fails the second writer instantly, and a
// pool of several connections makes sure there is a second writer. Both are in
// the template, and the repair puts both into a project that predates them.
func TestSQLiteOpensWithPragmasAndOneConnection(t *testing.T) {
	db := apiDatabaseGo()
	mustFormatGo(t, "database.go", db)

	if !strings.Contains(db, "func sqliteDSN(") || !strings.Contains(db, "busy_timeout(5000)") || !strings.Contains(db, "journal_mode(WAL)") {
		t.Error("database.go opens SQLite without busy_timeout and WAL")
	}
	if strings.Contains(db, `sqlite.Open(strings.TrimPrefix(dsn, "sqlite://"))`) {
		t.Error("a SQLite DSN still reaches GORM without the pragmas")
	}
	if !strings.Contains(db, `if strings.HasPrefix(dsn, "sqlite:") {
		maxOpen, maxIdle = 1, 1`) {
		t.Error("the pool hands SQLite more than one connection, which is how two writers collide")
	}
}

// :memory: takes neither pragma: there is no file to journal, and the database
// belongs to the one connection that opened it.
func TestSQLiteDSNLeavesMemoryAndExistingPragmasAlone(t *testing.T) {
	src := apiDatabaseGo()
	for _, want := range []string{
		`if path == "" || strings.Contains(path, ":memory:")`,
		`if !strings.Contains(path, "busy_timeout") {`,
		`if !strings.Contains(path, "journal_mode") {`,
		`if strings.Contains(path, "?") {`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("sqliteDSN does not guard %s", want)
		}
	}
}

func TestSQLiteBusyRepair(t *testing.T) {
	fresh := apiDatabaseGo()

	old := strings.Replace(fresh, sqliteOpenNew, sqliteOpenOld, 1)
	old = strings.Replace(old, sqliteDSNHelper, "", 1)
	old = strings.Replace(old, sqlitePoolClamp, "", 1)
	if old == fresh {
		t.Fatal("the fixture is not older than the template")
	}

	out, changes, warnings := repairSQLiteBusySource(old)
	if len(warnings) != 0 {
		t.Fatalf("repairing an untouched database.go warned: %v", warnings)
	}
	if len(changes) != 2 {
		t.Errorf("expected the pragmas and the pool, got %v", changes)
	}
	if out != fresh {
		t.Error("the repaired database.go is not the one the template writes")
	}

	again, changes, _ := repairSQLiteBusySource(out)
	if again != out || len(changes) != 0 {
		t.Error("the SQLite repair is not idempotent")
	}

	// A hand-edited file is told what to do rather than rewritten.
	edited := `package database

func Connect(dsn string) {
	db, err = gorm.Open(sqlite.Open(dsn), gormCfg)
	sqlDB.SetMaxOpenConns(maxOpen)
}
`
	_, changes, warnings = repairSQLiteBusySource(edited)
	if len(changes) != 0 || len(warnings) != 2 {
		t.Errorf("an edited database.go should get two notes and no edits (changes %v, warnings %v)", changes, warnings)
	}

	// A project with no SQLite anywhere is left alone.
	postgresOnly := "package database\n\nfunc Connect(dsn string) {}\n"
	out, changes, warnings = repairSQLiteBusySource(postgresOnly)
	if out != postgresOnly || len(changes) != 0 || len(warnings) != 0 {
		t.Error("a project that never opens SQLite should not be touched")
	}
}
