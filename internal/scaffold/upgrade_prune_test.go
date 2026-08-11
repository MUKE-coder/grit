package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte("// test\n"), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestExistsExactIsCaseSensitiveEverywhere(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "Providers.tsx"))

	if !existsExact(filepath.Join(dir, "Providers.tsx")) {
		t.Error("existsExact missed a file that is there")
	}
	// The whole point: os.Stat says yes to this on Windows and macOS.
	if existsExact(filepath.Join(dir, "providers.tsx")) {
		t.Error("existsExact matched a different capitalisation — the prune would " +
			"then delete the only copy believing a second one exists")
	}
}

func TestPruneAdminStraysRemovesDuplicatesOnly(t *testing.T) {
	admin := t.TempDir()

	// A stray beside its real counterpart: goes.
	write(t, filepath.Join(admin, "components", "tables", "DataTable.tsx"))
	write(t, filepath.Join(admin, "components", "tables", "data-table.tsx"))

	// A stray with no counterpart: stays, because deleting it would take the
	// last copy of something the app may still import.
	write(t, filepath.Join(admin, "components", "widgets", "StatsCard.tsx"))

	// Not in the map at all: untouched.
	write(t, filepath.Join(admin, "components", "tables", "date-filter.tsx"))

	removed := pruneAdminStrays(admin)
	if removed != 1 {
		t.Errorf("removed %d files, want 1", removed)
	}

	gone := filepath.Join(admin, "components", "tables", "DataTable.tsx")
	if existsExact(gone) {
		t.Error("the duplicate DataTable.tsx should have been removed")
	}
	for _, keep := range []string{
		filepath.Join(admin, "components", "tables", "data-table.tsx"),
		filepath.Join(admin, "components", "widgets", "StatsCard.tsx"),
		filepath.Join(admin, "components", "tables", "date-filter.tsx"),
	} {
		if !existsExact(keep) {
			t.Errorf("%s should have been left alone", filepath.Base(keep))
		}
	}
}

// The regression this whole helper exists for: on Windows the stray and the
// real file are the same file, and a Stat-based prune deletes it.
func TestPruneKeepsTheOnlyCopyOfACaseOnlyPair(t *testing.T) {
	admin := t.TempDir()
	write(t, filepath.Join(admin, "components", "shared", "Providers.tsx"))

	if n := pruneAdminStrays(admin); n != 0 {
		t.Errorf("removed %d files, want 0", n)
	}
	entries, err := os.ReadDir(filepath.Join(admin, "components", "shared"))
	if err != nil || len(entries) != 1 {
		t.Fatalf("expected the single Providers file to survive, got %v (err %v)", entries, err)
	}
}

func TestPruneUnwiredI18nRespectsAnInstalledNextIntl(t *testing.T) {
	for _, tc := range []struct {
		name    string
		pkg     string
		want    int
		survive bool
	}{
		{"no next-intl", `{"dependencies":{"next":"^16.0.0"}}`, 6, false},
		{"next-intl installed", `{"dependencies":{"next-intl":"^3.26.0"}}`, 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			admin := t.TempDir()
			if err := os.WriteFile(filepath.Join(admin, "package.json"), []byte(tc.pkg), 0644); err != nil {
				t.Fatal(err)
			}
			for _, f := range []string{
				"i18n/request.ts", "lib/locale.ts", "components/language-switcher.tsx",
				"messages/en.json", "messages/fr.json", "messages/sw.json",
			} {
				write(t, filepath.Join(admin, filepath.FromSlash(f)))
			}

			if n := pruneUnwiredI18n(admin); n != tc.want {
				t.Errorf("removed %d, want %d", n, tc.want)
			}
			survived := existsExact(filepath.Join(admin, "i18n", "request.ts"))
			if survived != tc.survive {
				t.Errorf("request.ts survived = %v, want %v", survived, tc.survive)
			}
		})
	}
}

// Every stray must name a real file that the scaffold actually writes,
// otherwise the "is the real one present?" guard can never pass and the entry
// is dead weight that quietly does nothing.
func TestStrayMapPointsAtRealScaffoldPaths(t *testing.T) {
	for stray, real := range adminStrayFiles {
		if stray == real {
			t.Errorf("%s maps to itself", stray)
		}
		if strings.Contains(real, "\\") {
			t.Errorf("%s: real path must use forward slashes, got %q", stray, real)
		}
	}
}
