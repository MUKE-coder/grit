package generate

// A versioned SQL pair, for the changes AutoMigrate cannot make.
//
// grit#162. Migrations are GORM AutoMigrate plus a computed history, which is
// additive: it adds tables, columns, indexes and constraints and takes nothing
// away. That is the right default and it covers almost everything. It does not
// cover a rename, a type change, a backfill, a drop, or adding an index to a
// large busy table, and until now a project that needed one of those left the
// framework's conventions entirely for golang-migrate or goose.
//
// The index case is the one that bites in production. A plain CREATE INDEX
// takes a write lock, so adding one to a big table is an outage rather than a
// migration. CREATE INDEX CONCURRENTLY does not, and Postgres refuses it inside
// a transaction, which is the detail people get wrong when they write it by
// hand. --index emits both halves correctly.

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// MigrationOptions is what `grit generate migration` was asked for.
type MigrationOptions struct {
	// Name becomes the file name, snake_cased: rename_company_to_organisation.
	Name string
	// Index, when set, pre-fills both files with an index build and its drop.
	// Written as table:column or table:col_a,col_b.
	Index string
	// Unique makes that index unique.
	Unique bool
	// IndexName overrides the derived idx_<table>_<columns>.
	IndexName string
	// Now is the version stamp's source. Zero means time.Now.
	Now time.Time
}

// migrationNameRe is what a migration may be called: the part of the file name
// after the version. Kept to the characters that are safe in a file name on
// every platform and readable in `grit migrate status`.
var migrationNameRe = regexp.MustCompile(`^[a-z0-9]+(_[a-z0-9]+)*$`)

// GenerateMigration writes the pair into the project in the current directory.
func GenerateMigration(opts MigrationOptions) error {
	root, err := findProjectRoot()
	if err != nil {
		return err
	}
	arch, _ := readGritJSON(root)
	return generateMigrationAt(root, arch, opts)
}

func generateMigrationAt(root, arch string, opts MigrationOptions) error {
	if strings.TrimSpace(opts.Name) == "" {
		return fmt.Errorf("a migration name is required (e.g. grit generate migration rename_company_to_organisation)")
	}
	slug := migrationSlug(opts.Name)
	if !migrationNameRe.MatchString(slug) {
		return fmt.Errorf("%q does not make a usable file name; use words, digits and underscores (e.g. backfill_order_totals)", opts.Name)
	}

	index, err := parseIndexSpec(opts.Index)
	if err != nil {
		return err
	}
	if opts.Unique && index == nil {
		return fmt.Errorf("--unique only means something with --index (e.g. --index users:email --unique)")
	}
	if index != nil {
		index.Unique = opts.Unique
		index.Name = opts.IndexName
		if index.Name == "" {
			index.Name = derivedIndexName(index)
		}
		index.Driver = projectDBProvider(root)
	}

	apiRoot, apiPrefix := root, ""
	if arch != "single" {
		apiRoot = filepath.Join(root, "apps", "api")
		apiPrefix = "apps/api/"
	}

	at := opts.Now
	if at.IsZero() {
		at = time.Now()
	}
	dir := filepath.Join(apiRoot, "migrations")
	version, err := freeVersion(dir, at)
	if err != nil {
		return err
	}
	if fileExists(filepath.Join(dir, version+"_"+slug+".up.sql")) {
		return fmt.Errorf("migrations/%s_%s.up.sql already exists", version, slug)
	}

	upPath := filepath.Join(dir, version+"_"+slug+".up.sql")
	downPath := filepath.Join(dir, version+"_"+slug+".down.sql")

	fmt.Printf("\n  Generating migration: %s\n\n", slug)

	if err := writeFileWithDirs(upPath, migrationUpSQL(version, slug, index)); err != nil {
		return fmt.Errorf("writing the up migration: %w", err)
	}
	fmt.Printf("  ✓ %smigrations/%s_%s.up.sql\n", apiPrefix, version, slug)

	if err := writeFileWithDirs(downPath, migrationDownSQL(version, slug, index)); err != nil {
		return fmt.Errorf("writing the down migration: %w", err)
	}
	fmt.Printf("  ✓ %smigrations/%s_%s.down.sql\n", apiPrefix, version, slug)

	fmt.Printf("\n  ✅ Migration generated.\n\n")
	if index != nil {
		fmt.Printf("  The index is written for you. Read it before running it:\n")
		if index.Driver == "postgres" {
			fmt.Printf("  CONCURRENTLY keeps the table writable and cannot run in a transaction,\n")
			fmt.Printf("  which is why the file opens with -- grit:no-transaction.\n\n")
		} else {
			fmt.Printf("  %s has no CONCURRENTLY, so this takes whatever lock it takes.\n\n", index.Driver)
		}
	} else {
		fmt.Printf("  Write the change in the up file and its reverse in the down file.\n\n")
	}
	fmt.Printf("  Then:\n\n")
	fmt.Printf("      grit migrate                  # AutoMigrate, then any pending pairs in order\n")
	fmt.Printf("      grit migrate status           # which pairs have run\n")
	fmt.Printf("      grit migrate down --sql       # undo the most recent pair\n\n")

	return nil
}

// freeVersion is the first timestamp from at onwards that no file in dir
// already uses.
//
// Two migrations generated in the same second share a version, and the first
// version of this took the clock's word for it: the second migration ran and
// then could not record itself, because the history's key was the version. The
// history's key is the whole id now, so this is belt as well as braces, but it
// also keeps the file names distinct, which is what makes `ls migrations` read
// as a sequence.
func freeVersion(dir string, at time.Time) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return at.UTC().Format("20060102150405"), nil
		}
		return "", fmt.Errorf("reading %s: %w", filepath.ToSlash(dir), err)
	}
	used := map[string]bool{}
	for _, e := range entries {
		if version, _, ok := strings.Cut(e.Name(), "_"); ok {
			used[version] = true
		}
	}
	// A minute of seconds is far more than any one command needs, and a loop
	// without a bound here would hang on a directory somebody filled by hand.
	for i := 0; i < 60; i++ {
		version := at.Add(time.Duration(i) * time.Second).UTC().Format("20060102150405")
		if !used[version] {
			return version, nil
		}
	}
	return "", fmt.Errorf("every second from %s onwards is already taken in migrations/; the versions there do not look like timestamps",
		at.UTC().Format("20060102150405"))
}

// migrationSlug turns whatever was typed into the file name's middle.
func migrationSlug(name string) string {
	s := strings.TrimSpace(name)
	// A name is as likely to arrive as AddOrderIndex or "add order index" as it
	// is already snake_cased.
	s = MakeNames(s).Snake
	s = strings.ToLower(s)
	s = strings.NewReplacer(" ", "_", "-", "_", ".", "_").Replace(s)
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}
	return strings.Trim(s, "_")
}

// IndexSpec is a --index argument, parsed.
type IndexSpec struct {
	Table   string
	Columns []string
	Unique  bool
	Name    string
	// Driver decides whether CONCURRENTLY is written.
	Driver string
}

// parseIndexSpec reads table:col or table:col_a,col_b.
func parseIndexSpec(spec string) (*IndexSpec, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil, nil
	}
	table, cols, ok := strings.Cut(spec, ":")
	if !ok {
		return nil, fmt.Errorf("--index wants table:column, for example --index orders:customer_id (got %q)", spec)
	}
	table = strings.TrimSpace(table)
	if table == "" {
		return nil, fmt.Errorf("--index is missing the table (e.g. --index orders:customer_id)")
	}
	out := &IndexSpec{Table: table}
	for _, c := range strings.Split(cols, ",") {
		if c = strings.TrimSpace(c); c != "" {
			out.Columns = append(out.Columns, c)
		}
	}
	if len(out.Columns) == 0 {
		return nil, fmt.Errorf("--index is missing the column(s) (e.g. --index orders:customer_id,created_at)")
	}
	return out, nil
}

// derivedIndexName is the convention GORM's own index names follow closely
// enough to read beside them.
func derivedIndexName(spec *IndexSpec) string {
	parts := append([]string{"idx", spec.Table}, spec.Columns...)
	name := strings.Join(parts, "_")
	// Postgres truncates an identifier at 63 bytes without complaining, so two
	// long names can collide after truncation. Better to be short on purpose.
	if len(name) > 63 {
		name = name[:63]
	}
	return name
}

// projectDBProvider reads DB_PROVIDER from the project's .env.
//
// It decides one thing: whether CONCURRENTLY is written. Guessing postgres
// would put a keyword MySQL and SQLite reject into the file, and guessing the
// other way would quietly drop the one line the index case exists for.
func projectDBProvider(root string) string {
	for _, name := range []string{".env", ".env.example"} {
		body, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(body), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "#") {
				continue
			}
			key, value, ok := strings.Cut(line, "=")
			if !ok || strings.TrimSpace(key) != "DB_PROVIDER" {
				continue
			}
			switch v := strings.ToLower(strings.Trim(strings.TrimSpace(value), `"'`)); v {
			case "postgres", "mysql", "sqlite":
				return v
			case "memory":
				// An in-memory database is SQLite.
				return "sqlite"
			}
		}
	}
	return "postgres"
}
