package doctor

import (
	"os"
	"path/filepath"
	"testing"
)

func migrationsProject(t *testing.T, files map[string]string) *project {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "migrations")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return &project{root: root, apiRoot: root}
}

func TestSQLMigrationCheckFindsAHalfPair(t *testing.T) {
	p := migrationsProject(t, map[string]string{
		"20261105120000_rename.up.sql": "ALTER TABLE a RENAME TO b;\n",
	})
	findings := checkSQLMigrations(p)
	if len(findings) != 1 {
		t.Fatalf("got %d finding(s), want 1: %+v", len(findings), findings)
	}
	if findings[0].Level != "error" {
		t.Errorf("level is %q; a directory that cannot be loaded stops every pending migration", findings[0].Level)
	}
	if got := findings[0].Resource; got != "migrations/20261105120000_rename" {
		t.Errorf("resource is %q", got)
	}
}

// An up file that changes the schema with an empty down file rolls back
// silently, which reads as a rollback that worked.
func TestSQLMigrationCheckFindsAnEmptyDown(t *testing.T) {
	p := migrationsProject(t, map[string]string{
		"20261105120000_rename.up.sql":   "ALTER TABLE a RENAME TO b;\n",
		"20261105120000_rename.down.sql": "-- nothing here yet\n",
	})
	findings := checkSQLMigrations(p)
	if len(findings) != 1 {
		t.Fatalf("got %d finding(s), want 1: %+v", len(findings), findings)
	}
	if findings[0].Level != "warning" {
		t.Errorf("level is %q, want warning: there are changes that genuinely cannot be undone", findings[0].Level)
	}
}

// A pair that is empty on both sides is a placeholder nobody has written yet.
// Reporting it would mean `grit generate migration` leaves the project with a
// warning, which trains people to ignore the warnings.
func TestSQLMigrationCheckIgnoresAFreshPlaceholder(t *testing.T) {
	header := "-- 20261105120000_rename.up.sql\n--\n-- Write the change here.\n"
	p := migrationsProject(t, map[string]string{
		"20261105120000_rename.up.sql":   header,
		"20261105120000_rename.down.sql": "-- Write the reverse here.\n",
	})
	if findings := checkSQLMigrations(p); len(findings) != 0 {
		t.Errorf("a freshly generated pair was reported: %+v", findings)
	}
}

func TestSQLMigrationCheckPassesACompletePair(t *testing.T) {
	p := migrationsProject(t, map[string]string{
		"README.md":                      "# migrations\n",
		"migrations.go":                  "package migrations\n",
		"20261105120000_rename.up.sql":   "ALTER TABLE a RENAME TO b;\n",
		"20261105120000_rename.down.sql": "ALTER TABLE b RENAME TO a;\n",
	})
	if findings := checkSQLMigrations(p); len(findings) != 0 {
		t.Errorf("a complete pair was reported: %+v", findings)
	}
}

// A project from before SQL migrations existed has no directory, and that is
// not a finding.
func TestSQLMigrationCheckIsSilentWithNoDirectory(t *testing.T) {
	p := &project{root: t.TempDir(), apiRoot: t.TempDir()}
	if findings := checkSQLMigrations(p); len(findings) != 0 {
		t.Errorf("a project with no migrations directory was reported: %+v", findings)
	}
}
