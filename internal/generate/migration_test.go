package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// migrationProject is a project root with just enough in it for the generator.
func migrationProject(t *testing.T, provider string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".env"),
		[]byte("APP_NAME=x\nDB_PROVIDER="+provider+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func readMigration(t *testing.T, root, version, slug, half string) string {
	t.Helper()
	path := filepath.Join(root, "apps", "api", "migrations", version+"_"+slug+"."+half+".sql")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(body)
}

func TestGenerateMigrationWritesBothHalves(t *testing.T) {
	root := migrationProject(t, "postgres")
	at := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	if err := generateMigrationAt(root, "triple", MigrationOptions{
		Name: "RenameCompanyToOrganisation",
		Now:  at,
	}); err != nil {
		t.Fatal(err)
	}

	up := readMigration(t, root, "20261105120000", "rename_company_to_organisation", "up")
	down := readMigration(t, root, "20261105120000", "rename_company_to_organisation", "down")

	// A placeholder pair has nothing to execute, which is what makes it safe to
	// commit before the SQL is written.
	for name, body := range map[string]string{"up": up, "down": down} {
		for _, line := range strings.Split(body, "\n") {
			if line = strings.TrimSpace(line); line != "" && !strings.HasPrefix(line, "--") {
				t.Errorf("the %s placeholder has something to run: %q", name, line)
			}
		}
	}
	if !strings.Contains(up, NoTransactionMarkerText) {
		t.Error("the up file does not mention the no-transaction marker, which is the one rule that cannot be guessed")
	}
}

// An index is the case the generator exists for, and the driver decides whether
// CONCURRENTLY may appear at all.
func TestGenerateMigrationIndexPerDriver(t *testing.T) {
	for _, tc := range []struct {
		provider     string
		wantUp       []string
		notWantUp    []string
		wantDown     []string
		wantUpMarker bool
	}{
		{
			provider:     "postgres",
			wantUp:       []string{"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_customer_id ON orders (customer_id);"},
			wantDown:     []string{"DROP INDEX CONCURRENTLY IF EXISTS idx_orders_customer_id;"},
			wantUpMarker: true,
		},
		{
			provider:  "mysql",
			wantUp:    []string{"CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders (customer_id);"},
			notWantUp: []string{"CONCURRENTLY IF"},
			// MySQL's DROP INDEX names the table and has no IF EXISTS.
			wantDown: []string{"DROP INDEX idx_orders_customer_id ON orders;"},
		},
		{
			provider:  "sqlite",
			wantUp:    []string{"CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders (customer_id);"},
			notWantUp: []string{"CONCURRENTLY IF"},
			wantDown:  []string{"DROP INDEX IF EXISTS idx_orders_customer_id;"},
		},
	} {
		t.Run(tc.provider, func(t *testing.T) {
			root := migrationProject(t, tc.provider)
			at := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
			if err := generateMigrationAt(root, "triple", MigrationOptions{
				Name:  "add_order_customer_index",
				Index: "orders:customer_id",
				Now:   at,
			}); err != nil {
				t.Fatal(err)
			}
			up := readMigration(t, root, "20261105120000", "add_order_customer_index", "up")
			down := readMigration(t, root, "20261105120000", "add_order_customer_index", "down")

			for _, want := range tc.wantUp {
				if !strings.Contains(up, want) {
					t.Errorf("up is missing %q:\n%s", want, up)
				}
			}
			for _, unwanted := range tc.notWantUp {
				if strings.Contains(up, unwanted) {
					t.Errorf("up has %q, which %s does not support", unwanted, tc.provider)
				}
			}
			for _, want := range tc.wantDown {
				if !strings.Contains(down, want) {
					t.Errorf("down is missing %q:\n%s", want, down)
				}
			}
			// The marker has to be the first line or the runner does not see it.
			hasMarker := strings.HasPrefix(up, NoTransactionMarkerText+"\n")
			if hasMarker != tc.wantUpMarker {
				t.Errorf("first line is the no-transaction marker = %v, want %v:\n%s",
					hasMarker, tc.wantUpMarker, strings.SplitN(up, "\n", 2)[0])
			}
			if strings.HasPrefix(down, NoTransactionMarkerText) != tc.wantUpMarker {
				t.Error("DROP INDEX CONCURRENTLY needs the marker as much as the build does")
			}
		})
	}
}

func TestGenerateMigrationUniqueAndCustomName(t *testing.T) {
	root := migrationProject(t, "postgres")
	at := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	if err := generateMigrationAt(root, "triple", MigrationOptions{
		Name:      "one_live_email",
		Index:     "users:email,active",
		Unique:    true,
		IndexName: "users_live_email",
		Now:       at,
	}); err != nil {
		t.Fatal(err)
	}
	up := readMigration(t, root, "20261105120000", "one_live_email", "up")
	if !strings.Contains(up, "CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS users_live_email ON users (email, active);") {
		t.Errorf("unique multi-column index is wrong:\n%s", up)
	}
}

// Two in one second used to produce one version, and the second migration then
// ran and failed to record itself.
func TestGenerateMigrationTakesTheNextFreeSecond(t *testing.T) {
	root := migrationProject(t, "sqlite")
	at := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	for _, name := range []string{"first", "second", "third"} {
		if err := generateMigrationAt(root, "triple", MigrationOptions{Name: name, Now: at}); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
	entries, err := os.ReadDir(filepath.Join(root, "apps", "api", "migrations"))
	if err != nil {
		t.Fatal(err)
	}
	versions := map[string]bool{}
	for _, e := range entries {
		version, _, _ := strings.Cut(e.Name(), "_")
		versions[version] = true
	}
	if len(versions) != 3 {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("three migrations generated in one second share %d version(s): %v", len(versions), names)
	}
}

func TestGenerateMigrationRefusesBadInput(t *testing.T) {
	root := migrationProject(t, "postgres")
	at := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name string
		opts MigrationOptions
		want string
	}{
		{"no name", MigrationOptions{Now: at}, "name is required"},
		{"unusable name", MigrationOptions{Name: "!!!", Now: at}, "usable file name"},
		{"index without a column", MigrationOptions{Name: "x", Index: "orders", Now: at}, "table:column"},
		{"index with an empty column list", MigrationOptions{Name: "x", Index: "orders:", Now: at}, "missing the column"},
		{"unique without index", MigrationOptions{Name: "x", Unique: true, Now: at}, "only means something with --index"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := generateMigrationAt(root, "triple", tc.opts)
			if err == nil {
				t.Fatal("no error")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error is %q, want it to mention %q", err, tc.want)
			}
		})
	}
}

// A single-architecture project has no apps/api.
func TestGenerateMigrationHonoursTheSingleArchitecture(t *testing.T) {
	root := migrationProject(t, "sqlite")
	at := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	if err := generateMigrationAt(root, "single", MigrationOptions{Name: "x", Now: at}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "migrations", "20261105120000_x.up.sql")); err != nil {
		t.Errorf("a single project's migration is not at the project root: %v", err)
	}
}

func TestProjectDBProviderReadsTheEnv(t *testing.T) {
	for env, want := range map[string]string{
		"DB_PROVIDER=mysql\n":            "mysql",
		"DB_PROVIDER=\"sqlite\"\n":       "sqlite",
		"DB_PROVIDER=memory\n":           "sqlite",
		"# DB_PROVIDER=mysql\n":          "postgres",
		"DB_PROVIDER=something-else\n":   "postgres",
		"APP_NAME=x\nDB_PROVIDER=mysql ": "mysql",
	} {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, ".env"), []byte(env), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := projectDBProvider(root); got != want {
			t.Errorf("projectDBProvider(%q) = %q, want %q", env, got, want)
		}
	}
}
