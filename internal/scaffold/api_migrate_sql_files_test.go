package scaffold

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// The marker has one spelling, and the emitted runner is the thing that has to
// read what the generator writes.
//
// Two tests, one per side, would both pass while the pair was broken: the
// generator could write -- grit:notransaction and the runner would go on
// looking for -- grit:no-transaction, and the only symptom would be Postgres
// refusing a concurrent index build at the moment somebody ran it in
// production. So this reads the emitted source.
func TestEmittedRunnerReadsTheMarkerTheGeneratorWrites(t *testing.T) {
	src := apiMigrateSQLGo()
	want := `const NoTransactionMarker = "` + NoTransactionMarker + `"`
	if !strings.Contains(src, want) {
		t.Errorf("internal/migrate/sql.go does not declare %s", want)
	}
	if NoTransactionMarker != "-- grit:no-transaction" {
		t.Errorf("the marker is %q; it is in generated files in the wild, so changing it "+
			"orphans every migration that already has it", NoTransactionMarker)
	}
}

// The three files are a set: the runner imports nothing of the package, but
// cmd/migrate imports both, so delivering one without the other is a project
// that does not compile.
func TestMigrationFilesParse(t *testing.T) {
	for name, src := range map[string]string{
		"internal/migrate/sql.go":      apiMigrateSQLGo(),
		"internal/migrate/sql_test.go": apiMigrateSQLTestGo(),
		"migrations/migrations.go":     apiMigrationsPackageGo(),
	} {
		src = strings.ReplaceAll(src, "{{MODULE}}", "example/apps/api")
		if _, err := parser.ParseFile(token.NewFileSet(), name, src, parser.AllErrors); err != nil {
			t.Errorf("%s does not parse: %v", name, err)
		}
	}
}

// The embed pattern is the one line in migrations.go that can fail a build in a
// project that has done nothing wrong.
//
// *.sql would be the obvious pattern and is a build error until the first
// migration exists, because an embed pattern that matches nothing does not
// compile. * always matches: README.md and migrations.go are in there.
func TestMigrationsPackageEmbedsEverything(t *testing.T) {
	src := apiMigrationsPackageGo()
	if !strings.Contains(src, "//go:embed *\n") {
		t.Error("migrations/migrations.go does not embed *; *.sql fails to build until the first migration exists")
	}
	if strings.Contains(src, "//go:embed *.sql") {
		t.Error("migrations/migrations.go embeds *.sql, which is a build error in a project with no migrations yet")
	}
}

// The runner is reached through the hub of framework-owned files, so it travels
// on upgrade. cmd/migrate imports migrations, and cmd/migrate is rewritten by
// every upgrade: delivering the rewritten tool without the package it imports
// is a project that stops compiling.
func TestTheMigrateToolAndItsPackagesTravelTogether(t *testing.T) {
	tool := apiMigrateMainGo()
	if !strings.Contains(tool, `"{{MODULE}}/migrations"`) {
		t.Fatal("cmd/migrate/main.go does not import the migrations package")
	}
	for _, want := range []string{
		"migrate.EnsureSQLHistory(db)",
		"migrate.LoadSQLMigrations(migrations.FS)",
		"migrate.ApplySQL(db, all)",
		"migrate.RollbackSQL(db, all, steps)",
	} {
		if !strings.Contains(tool, want) {
			t.Errorf("cmd/migrate/main.go is missing %q", want)
		}
	}

	root := t.TempDir()
	opts := Options{ProjectName: "shop", Architecture: ArchTriple, Frontend: FrontendNext}
	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}
	if err := writeAPIFiles(root, opts); err != nil {
		t.Fatalf("writeAPIFiles: %v", err)
	}
	for _, path := range []string{
		"apps/api/migrations/migrations.go",
		"apps/api/migrations/README.md",
		"apps/api/internal/migrate/sql.go",
		"apps/api/internal/migrate/sql_test.go",
	} {
		if !fileExists(root + "/" + path) {
			t.Errorf("%s was not written", path)
		}
	}
}
