package scaffold

import (
	"strings"
	"testing"
)

// grit#194: the test database has to be the database the app ships against.
//
// Both findings in that issue are the same shape, and the shape is worth
// naming: testkit builds its schema with AutoMigrate, and AutoMigrate is not
// what `grit migrate` does. Anything the migrate command adds afterwards is
// missing from every test, silently, and the things it adds afterwards are
// exactly the ones AutoMigrate cannot express.

// The append-only triggers are installed, or an append-only table accepts
// UPDATE and DELETE in every test.
//
// The guarantee --append-only sells is that the rule survives somebody with a
// SQL client. A project that writes the obvious test for that gets a red build
// caused by the harness, and the natural next move is to weaken the test until
// it passes, which is the opposite of what anybody wants. A project that does
// not write it has a table it believes is immutable and a suite that would
// never notice if it stopped being.
func TestTestkitInstallsTheAppendOnlyTriggers(t *testing.T) {
	src := apiTestKitGo()

	if !strings.Contains(src, "appendonly.InstallTriggers(db)") {
		t.Error("testkit does not install the append-only triggers, so an append-only table " +
			"accepts UPDATE and DELETE in every test and the guarantee is untested")
	}
	if !strings.Contains(src, "{{MODULE}}/internal/appendonly") {
		t.Error("testkit calls appendonly.InstallTriggers and does not import the package")
	}

	// Both dialects. A fix applied to the SQLite path alone would leave the
	// Postgres path, which is where the TRUNCATE guard lives, untested.
	if n := strings.Count(src, "installTriggers(tb, db)"); n != 2 {
		t.Errorf("installTriggers is called %d times; it belongs after both AutoMigrate "+
			"calls, the SQLite one and the Postgres one", n)
	}

	// And after the migration, not before: there is no table to put a trigger
	// on until AutoMigrate has run.
	for _, after := range []string{
		"migrating the test database",
		"migrating schema %s",
	} {
		at := strings.Index(src, after)
		if at == -1 {
			t.Fatalf("could not find the AutoMigrate error for %q", after)
		}
		call := strings.Index(src[at:], "installTriggers(tb, db)")
		if call == -1 {
			t.Errorf("installTriggers is not called after the AutoMigrate that reports %q", after)
		}
	}
}

// Postgres search_path keeps public on it, or a model with a vector column
// cannot migrate and takes every test in the package with it.
//
// pgvector's type lives in whichever schema the extension was created in,
// which is public. A path holding only the test's own schema makes AutoMigrate
// fail with `type "vector" does not exist` on the first model that has an
// embedding, including for tests that have nothing to do with embeddings. New
// in practice because vector shipped in v3.491.0.
func TestTestkitSearchPathKeepsPublic(t *testing.T) {
	src := apiTestKitGo()

	if !strings.Contains(src, `"SET search_path TO " + schema + ", public"`) {
		t.Error(`postgresDB does not keep public on the search_path, so a model with a ` +
			`vector column fails to migrate with "type \"vector\" does not exist" and ` +
			`every Postgres test in the package fails with it`)
	}
	// The isolation the file promises is unchanged: the test's own schema is
	// still first, which is where CREATE TABLE lands.
	path := strings.Index(src, `"SET search_path TO " + schema`)
	if path == -1 {
		t.Fatal("no search_path is set at all")
	}
	if strings.Contains(src, `"SET search_path TO public`) {
		t.Error("public is first on the path, so the test's tables land in the shared schema")
	}
}

// grit#198: the checked-in SQL pairs reach the test schema, and a test can
// insist on them.
//
// Third instance of the shape above, and the worst of the three. A missing
// trigger fails loudly, in the refusal direction: a test asserting an UPDATE
// is refused gets an accepted UPDATE and goes red. A missing CHECK fails
// silently in the acceptance direction: a test asserting that a valid row is
// accepted goes green while the constraint is not in the database at all.
func TestTestkitAppliesTheCheckedInSQLMigrations(t *testing.T) {
	src := apiTestKitGo()

	for _, want := range []string{
		"migrate.LoadSQLMigrations(migrations.FS)",
		"migrate.EnsureSQLHistory(db)",
		"migrate.ApplySQL(db, all)",
		`"{{MODULE}}/migrations"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("testkit does not %s, so anything a project expresses in a pair "+
				"(a CHECK, a trigger, a backfill, a partial index) is absent from every test", want)
		}
	}

	// After AutoMigrate: a pair alters tables AutoMigrate has to have made.
	mig := strings.Index(src, "migrating schema %s")
	apply := strings.Index(src, "applySQLMigrations(tb, db)")
	if mig == -1 || apply == -1 {
		t.Fatal("could not find the Postgres AutoMigrate and the apply call")
	}
	if apply < mig {
		t.Error("the pairs are applied before AutoMigrate, so they alter tables that do not exist yet")
	}
}

// The guard has to be answerable before a database exists, because that is
// where a guard goes, and it must not be a package-level flag: DB mutating one
// is not safe under t.Parallel, and a guard placed first would always fail.
func TestRequireSQLMigrationsIsOrderIndependent(t *testing.T) {
	src := apiTestKitGo()

	if strings.Contains(src, "var sqlMigrationsApplied") {
		t.Error("the guard reads a package-level flag that DB sets: a test calling it before " +
			"DB always fails, and a global mutated by DB is not safe under t.Parallel")
	}
	if !strings.Contains(src, "func sqlMigrationsEnabled() bool") {
		t.Error("there is no environment-only predicate for whether the pairs will be applied")
	}
	// One predicate, read by both the guard and the behaviour, or they drift.
	if n := strings.Count(src, "sqlMigrationsEnabled()"); n < 2 {
		t.Errorf("sqlMigrationsEnabled is read %d times; the guard and applySQLMigrations "+
			"both have to ask it or they can disagree", n)
	}

	// It fails rather than skips. A skipped guarantee reads as one that held.
	at := strings.Index(src, "func RequireSQLMigrations(")
	if at == -1 {
		t.Fatal("there is no RequireSQLMigrations")
	}
	body := src[at:]
	if end := strings.Index(body, "\n}"); end != -1 {
		body = body[:end]
	}
	// The call, not the word: the message names SkipSQLMigrationsEnv, which an
	// unanchored search for "Skip" matched. Checking for the substring rather
	// than the call is how this test first failed against correct code.
	for _, call := range []string{"tb.Skip(", "tb.Skipf(", "tb.SkipNow("} {
		if strings.Contains(body, call) {
			t.Errorf("RequireSQLMigrations calls %s. A guarantee that skipped is not a "+
				"guarantee that held, and the default fall back to SQLite makes that easy "+
				"to miss:\n%s", call, body)
		}
	}
	if !strings.Contains(body, "Fatalf") {
		t.Errorf("RequireSQLMigrations does not fail:\n%s", body)
	}
}
