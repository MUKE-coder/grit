package scaffold

// The shipped tests for internal/migrate/sql.go.
//
// The statement splitter is the part worth testing hard. It exists because pgx
// refuses more than one command per Exec and the MySQL driver needs
// multiStatements on the DSN, so a .sql file has to be split before it runs.
// Splitting on every semicolon turns one PL/pgSQL function into three syntax
// errors, and that is exactly the SQL people reach for a migration file to
// write.

// apiMigrateSQLTestGo emits internal/migrate/sql_test.go.
func apiMigrateSQLTestGo() string {
	return `package migrate

import (
	"testing"
	"testing/fstest"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestSplitSQLKeepsSemicolonsThatAreInsideSomething(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
		want []string
	}{
		{
			name: "two statements",
			body: "ALTER TABLE a RENAME TO b;\nALTER TABLE c RENAME TO d;\n",
			want: []string{"ALTER TABLE a RENAME TO b", "ALTER TABLE c RENAME TO d"},
		},
		{
			name: "a semicolon inside a string is not a statement boundary",
			body: "UPDATE t SET note = 'a; b' WHERE id = 1;",
			want: []string{"UPDATE t SET note = 'a; b' WHERE id = 1"},
		},
		{
			name: "an escaped quote does not end the string",
			body: "UPDATE t SET note = 'it''s; fine' WHERE id = 1;",
			want: []string{"UPDATE t SET note = 'it''s; fine' WHERE id = 1"},
		},
		{
			name: "a quoted identifier may hold one too",
			body: ` + "`" + `ALTER TABLE "odd;name" RENAME TO ok;` + "`" + `,
			want: []string{` + "`" + `ALTER TABLE "odd;name" RENAME TO ok` + "`" + `},
		},
		{
			// Every PL/pgSQL function and trigger is written this way, and a
			// naive split makes three broken fragments of it.
			name: "a dollar-quoted body is one statement",
			body: "CREATE FUNCTION f() RETURNS int AS $$\nBEGIN\n  PERFORM 1;\n  RETURN 2;\nEND;\n$$ LANGUAGE plpgsql;\nSELECT 1;",
			want: []string{
				"CREATE FUNCTION f() RETURNS int AS $$\nBEGIN\n  PERFORM 1;\n  RETURN 2;\nEND;\n$$ LANGUAGE plpgsql",
				"SELECT 1",
			},
		},
		{
			name: "a tagged dollar quote too",
			body: "CREATE FUNCTION f() RETURNS int AS $body$ SELECT 1; $body$ LANGUAGE sql;",
			want: []string{"CREATE FUNCTION f() RETURNS int AS $body$ SELECT 1; $body$ LANGUAGE sql"},
		},
		{
			name: "a line comment may mention one",
			body: "-- drop it; then rebuild\nDROP INDEX idx_a;",
			want: []string{"-- drop it; then rebuild\nDROP INDEX idx_a"},
		},
		{
			name: "a file that is only comments runs nothing",
			body: "-- grit:no-transaction\n--\n-- Write the change here.\n",
			want: nil,
		},
		{
			name: "a missing final semicolon is still a statement",
			body: "DROP INDEX idx_a",
			want: []string{"DROP INDEX idx_a"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, SplitSQL(tc.body))
		})
	}
}

func TestNoTransactionMarkerIsTheFirstLineOnly(t *testing.T) {
	assert.True(t, hasNoTransactionMarker("-- grit:no-transaction\nCREATE INDEX CONCURRENTLY x ON t (c);"))
	assert.True(t, hasNoTransactionMarker("\n\n-- grit:no-transaction\nSELECT 1;"))
	// A comment about transactions further down is a comment, not a request.
	assert.False(t, hasNoTransactionMarker("-- the index build\n-- grit:no-transaction\nSELECT 1;"))
	assert.False(t, hasNoTransactionMarker("SELECT 1;"))
}

func TestLoadSQLMigrationsPairsAndOrdersThem(t *testing.T) {
	fsys := fstest.MapFS{
		"README.md": {Data: []byte("# migrations")},
		// Out of order on purpose: the version decides, not the directory.
		"20261105120000_second.up.sql":   {Data: []byte("SELECT 2;")},
		"20261105120000_second.down.sql": {Data: []byte("SELECT -2;")},
		"20261101090000_first.up.sql":    {Data: []byte("-- grit:no-transaction\nSELECT 1;")},
		"20261101090000_first.down.sql":  {Data: []byte("SELECT -1;")},
	}
	all, err := LoadSQLMigrations(fsys)
	require.NoError(t, err)
	require.Len(t, all, 2)

	assert.Equal(t, "20261101090000_first", all[0].ID())
	assert.Equal(t, "20261105120000_second", all[1].ID())
	assert.True(t, all[0].UpNoTransaction, "the marker on the first line was not read")
	assert.False(t, all[1].UpNoTransaction)
	assert.Len(t, all[0].Checksum, 64)
}

// A pair with one half is a migration that cannot be rolled back, and the
// moment to find that out is not the moment you need to roll it back.
func TestLoadSQLMigrationsRefusesAHalfPair(t *testing.T) {
	_, err := LoadSQLMigrations(fstest.MapFS{
		"20261101090000_first.up.sql": {Data: []byte("SELECT 1;")},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no .down.sql")

	_, err = LoadSQLMigrations(fstest.MapFS{
		"20261101090000_first.down.sql": {Data: []byte("SELECT 1;")},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no .up.sql")
}

func newSQLMigrateDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, EnsureSQLHistory(db))
	return db
}

func TestApplySQLRunsEachPairOnceAndRecordsIt(t *testing.T) {
	db := newSQLMigrateDB(t)
	all, err := LoadSQLMigrations(fstest.MapFS{
		"20261101090000_make_notes.up.sql":   {Data: []byte("CREATE TABLE notes (id integer primary key, body text);")},
		"20261101090000_make_notes.down.sql": {Data: []byte("DROP TABLE notes;")},
		"20261102090000_seed_notes.up.sql":   {Data: []byte("INSERT INTO notes (body) VALUES ('a; b');")},
		"20261102090000_seed_notes.down.sql": {Data: []byte("DELETE FROM notes;")},
	})
	require.NoError(t, err)

	applied, err := ApplySQL(db, all)
	require.NoError(t, err)
	assert.Equal(t, []string{"20261101090000_make_notes", "20261102090000_seed_notes"}, applied)

	var body string
	require.NoError(t, db.Raw("SELECT body FROM notes").Scan(&body).Error)
	assert.Equal(t, "a; b", body, "the semicolon inside the string was treated as a statement boundary")

	// A second run applies nothing, which is the whole point of recording them.
	again, err := ApplySQL(db, all)
	require.NoError(t, err)
	assert.Empty(t, again)

	pending, err := PendingSQL(db, all)
	require.NoError(t, err)
	assert.Empty(t, pending)
}

func TestRollbackSQLUndoesTheMostRecentPair(t *testing.T) {
	db := newSQLMigrateDB(t)
	all, err := LoadSQLMigrations(fstest.MapFS{
		"20261101090000_make_notes.up.sql":   {Data: []byte("CREATE TABLE notes (id integer primary key, body text);")},
		"20261101090000_make_notes.down.sql": {Data: []byte("DROP TABLE notes;")},
		"20261102090000_seed_notes.up.sql":   {Data: []byte("INSERT INTO notes (body) VALUES ('x');")},
		"20261102090000_seed_notes.down.sql": {Data: []byte("DELETE FROM notes;")},
	})
	require.NoError(t, err)
	_, err = ApplySQL(db, all)
	require.NoError(t, err)

	undone, err := RollbackSQL(db, all, 1)
	require.NoError(t, err)
	assert.Equal(t, []string{"20261102090000_seed_notes"}, undone)

	var count int64
	require.NoError(t, db.Raw("SELECT count(*) FROM notes").Scan(&count).Error)
	assert.Zero(t, count)

	// Rolled back means pending again.
	pending, err := PendingSQL(db, all)
	require.NoError(t, err)
	require.Len(t, pending, 1)
	assert.Equal(t, "20261102090000_seed_notes", pending[0].ID())
}

// Editing a pair that has already run is the classic mistake, and it is silent:
// the database holds one version and the file holds another.
func TestChangedSQLNoticesAnEditedPair(t *testing.T) {
	db := newSQLMigrateDB(t)
	before := fstest.MapFS{
		"20261101090000_make_notes.up.sql":   {Data: []byte("CREATE TABLE notes (id integer primary key);")},
		"20261101090000_make_notes.down.sql": {Data: []byte("DROP TABLE notes;")},
	}
	all, err := LoadSQLMigrations(before)
	require.NoError(t, err)
	_, err = ApplySQL(db, all)
	require.NoError(t, err)

	changed, err := ChangedSQL(db, all)
	require.NoError(t, err)
	assert.Empty(t, changed)

	edited, err := LoadSQLMigrations(fstest.MapFS{
		"20261101090000_make_notes.up.sql":   {Data: []byte("CREATE TABLE notes (id integer primary key, body text);")},
		"20261101090000_make_notes.down.sql": {Data: []byte("DROP TABLE notes;")},
	})
	require.NoError(t, err)
	changed, err = ChangedSQL(db, edited)
	require.NoError(t, err)
	assert.Equal(t, []string{"20261101090000_make_notes"}, changed)
}

// A failing file leaves nothing behind, and the error names the statement.
func TestApplySQLStopsAtTheFailingStatement(t *testing.T) {
	db := newSQLMigrateDB(t)
	all, err := LoadSQLMigrations(fstest.MapFS{
		"20261101090000_bad.up.sql":   {Data: []byte("CREATE TABLE ok (id integer primary key);\nTHIS IS NOT SQL;")},
		"20261101090000_bad.down.sql": {Data: []byte("DROP TABLE ok;")},
	})
	require.NoError(t, err)

	applied, err := ApplySQL(db, all)
	require.Error(t, err)
	assert.Empty(t, applied)
	assert.Contains(t, err.Error(), "20261101090000_bad.up.sql")
	assert.Contains(t, err.Error(), "THIS IS NOT SQL")
	assert.False(t, db.Migrator().HasTable("ok"),
		"the first statement stayed applied, so the file did not run in a transaction")
}
`
}
