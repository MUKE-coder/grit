package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// writeAppendOnlyFiles writes internal/appendonly into a new or upgraded
// project.
//
// Every project gets it, whether or not anything is append-only yet. The two
// calls that install it, one in Connect and one in Migrate, are no-ops with
// nothing registered, and having them in place from the start is what lets
// `grit generate resource X --append-only` work on the day someone needs it
// rather than after a round of wiring.
func writeAppendOnlyFiles(root string, opts Options) error {
	return WriteAppendOnlyPackage(opts.APIRoot(root), opts.Module(), true)
}

// WriteAppendOnlyPackage writes the package under apiRoot.
//
// overwrite=false leaves an existing copy alone. The generator uses that on a
// project scaffolded before the package existed: it needs the package to be
// there, and has no business replacing one that is.
func WriteAppendOnlyPackage(apiRoot, module string, overwrite bool) error {
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "appendonly", "appendonly.go"):      appendOnlyGo(),
		filepath.Join(apiRoot, "internal", "appendonly", "appendonly_test.go"): appendOnlyTestGo(),
	}
	for path, content := range files {
		if !overwrite {
			if _, err := os.Stat(path); err == nil {
				continue
			}
		}
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

// The two hooks, as they appear in a freshly scaffolded project. Shared with
// EnsureAppendOnlyWiring so an older project ends up with exactly the same
// code as a new one.
const appendOnlyConnectHook = `	// Append-only tables refuse UPDATE and DELETE through this handle, so GORM
	// Studio's row editor, the CSV importer and every handler are bound by it,
	// not just the one service that remembered. See internal/appendonly.
	if err := appendonly.Install(db); err != nil {
		return nil, fmt.Errorf("installing append-only guard: %w", err)
	}

`

const appendOnlyMigrateHook = `	// And at the database, for whatever does not go through GORM: the SQL
	// editor in Studio, psql, a script somebody writes next year.
	if err := appendonly.InstallTriggers(db); err != nil {
		return fmt.Errorf("installing append-only triggers: %w", err)
	}

`

// EnsureAppendOnlyWiring makes a project call the append-only guard when it
// connects and when it migrates.
//
// A project scaffolded before the package existed calls neither. A model
// generated with --append-only would then register itself with a guard nobody
// installs: every route and admin button gone, so it looks finished, and any
// UPDATE still goes through. That is worse than no flag at all, so this either
// wires both calls or refuses and says where they go.
//
// Idempotent: a project that already makes the call is left as it is.
func EnsureAppendOnlyWiring(apiRoot, module string) error {
	steps := []struct {
		path, has, anchor, hook string
	}{
		{
			path:   filepath.Join(apiRoot, "internal", "database", "database.go"),
			has:    "appendonly.Install(db)",
			anchor: "\tsqlDB, err := db.DB()",
			hook:   appendOnlyConnectHook,
		},
		{
			path:   filepath.Join(apiRoot, "internal", "models", "user.go"),
			has:    "appendonly.InstallTriggers(db)",
			anchor: "\tif err := SeedRoles(db); err != nil {",
			hook:   appendOnlyMigrateHook,
		},
	}

	importPath := module + "/internal/appendonly"
	for _, s := range steps {
		raw, err := os.ReadFile(s.path)
		if err != nil {
			return fmt.Errorf("append-only needs %s: %w", s.path, err)
		}
		crlf := strings.Contains(string(raw), "\r\n")
		content := strings.ReplaceAll(string(raw), "\r\n", "\n")

		if strings.Contains(content, s.has) {
			continue
		}
		if !strings.Contains(content, s.anchor) {
			return fmt.Errorf("could not find where to install the append-only guard in %s.\n\n"+
				"Add this call yourself, then generate again:\n\n  %s", s.path, s.has)
		}
		content = strings.Replace(content, s.anchor, s.hook+s.anchor, 1)

		var ok bool
		if content, ok = addImportGroup(content, importPath); !ok {
			return fmt.Errorf("could not add the %q import to %s: add it by hand", importPath, s.path)
		}
		if crlf {
			content = strings.ReplaceAll(content, "\n", "\r\n")
		}
		if err := os.WriteFile(s.path, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", s.path, err)
		}
		// Grit made this edit, so the file should not start reading as one the
		// developer changed and upgrade should keep its hands off.
		manifest.Refresh(s.path)
		fmt.Printf("  ✓ Wired the append-only guard into %s\n", filepath.Base(s.path))
	}
	return nil
}

// addImportGroup adds path to the end of the file's import block, as a group
// of its own so gofmt has nothing to reorder. Reports false when there is no
// parenthesised import block to add it to.
func addImportGroup(content, path string) (string, bool) {
	line := "\t\"" + path + "\""
	if strings.Contains(content, line) {
		return content, true
	}
	start := strings.Index(content, "\nimport (\n")
	if start < 0 {
		return content, false
	}
	end := strings.Index(content[start:], "\n)\n")
	if end < 0 {
		return content, false
	}
	at := start + end
	return content[:at] + "\n\n" + line + content[at:], true
}

func appendOnlyGo() string {
	return `// Package appendonly makes a table refuse UPDATE and DELETE.
//
// For records that must never change once written: journal entries, audit
// events, consent receipts. A correction is a new row that reverses the old
// one, which is how an accountant expects it and how an auditor can follow it.
//
// Two layers, because each covers what the other cannot.
//
// Install registers GORM callbacks on the connection. Every write through GORM
// to a registered table is refused with a respond.Rule, so the caller gets 422
// and the reason. That covers the handlers, the CSV importer, the offline sync
// endpoint and GORM Studio's row editor, which all share the handle.
//
// InstallTriggers puts a trigger on each table, and grit migrate runs it. That
// covers what never touches GORM: the Studio SQL editor, psql, a script. It
// exists because the first layer alone was tested against a real ledger, and
// one UPDATE typed into the Studio SQL box unbalanced the books.
//
// Tables join by registering their model, which a resource generated with
// --append-only does from its model file.
package appendonly

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	"gorm.io/gorm"

	"{{MODULE}}/internal/respond"
)

var (
	mu       sync.RWMutex
	registry []interface{}
)

// Register marks a model's table append-only. Call it from the model file's
// init(), so the registration exists before the server connects or the
// migrate command runs, and goes away with the file.
func Register(model interface{}) {
	mu.Lock()
	defer mu.Unlock()
	registry = append(registry, model)
}

// Tables resolves every registered model to the table GORM writes it to.
//
// Through GORM's own parser rather than a guess from the type name, so a
// custom TableName() or an irregular plural cannot leave a table unguarded.
// GORM caches the parse, so this is cheap enough to call per statement.
func Tables(db *gorm.DB) ([]string, error) {
	mu.RLock()
	defer mu.RUnlock()
	tables := make([]string, 0, len(registry))
	for _, model := range registry {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			return nil, fmt.Errorf("resolving the table for %T: %w", model, err)
		}
		tables = append(tables, stmt.Schema.Table)
	}
	return tables, nil
}

// Install registers the GORM guard on db. Call it once, straight after
// connecting.
func Install(db *gorm.DB) error {
	if err := db.Callback().Update().Before("gorm:update").
		Register("appendonly:update", refuse); err != nil {
		return fmt.Errorf("registering the update guard: %w", err)
	}
	if err := db.Callback().Delete().Before("gorm:delete").
		Register("appendonly:delete", refuse); err != nil {
		return fmt.Errorf("registering the delete guard: %w", err)
	}
	return nil
}

func refuse(tx *gorm.DB) {
	if tx.Statement == nil || tx.Error != nil {
		return
	}
	// By table name, not by Go type. Studio's row editor writes with
	// db.Table("x") and no model at all, and it is the writer that most needs
	// stopping.
	table := tx.Statement.Table
	if table == "" && tx.Statement.Schema != nil {
		table = tx.Statement.Schema.Table
	}
	if table == "" {
		return
	}
	tables, err := Tables(tx)
	if err != nil {
		_ = tx.AddError(err)
		return
	}
	for _, t := range tables {
		if t == table {
			_ = tx.AddError(respond.Rule(
				"%s is append-only: records are never changed or deleted, add a correcting entry instead", table))
			return
		}
	}
}

// identifier is what a table name has to look like before it goes into DDL.
// Trigger statements cannot take the table as a bound parameter.
var identifier = regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*$")

// InstallTriggers puts an UPDATE and DELETE trigger on every registered table.
// Repeatable: grit migrate runs it every time.
func InstallTriggers(db *gorm.DB) error {
	tables, err := Tables(db)
	if err != nil {
		return err
	}
	if len(tables) == 0 {
		return nil
	}

	dialect := db.Dialector.Name()
	if dialect == "postgres" {
		if err := db.Exec(postgresFunction).Error; err != nil {
			return fmt.Errorf("creating the append-only trigger function: %w", err)
		}
	}

	for _, table := range tables {
		if !identifier.MatchString(table) {
			return fmt.Errorf("refusing to build a trigger for the table name %q", table)
		}
		statements := triggerSQL(dialect, table)
		if statements == nil {
			log.Printf("append-only: no trigger for the %s dialect, so %s is guarded by GORM only and raw SQL can still change it", dialect, table)
			continue
		}
		for _, statement := range statements {
			if err := db.Exec(statement).Error; err != nil {
				return fmt.Errorf("append-only trigger on %s: %w", table, err)
			}
		}
	}
	return nil
}

// Suspend switches this project's append-only triggers off for the rest of the
// transaction tx, and returns the function that switches them back on.
//
// For restoring a backup, which has to rewrite these tables wholesale and is
// the one writer the triggers must let through. Other sessions never see the
// gap: ALTER TABLE inside a transaction is invisible until it commits, and the
// triggers are back on before this one does.
//
// ALTER TABLE rather than session_replication_role, which needs a superuser;
// the table owner is all a restore should require. Postgres only, like the
// restore that calls it.
func Suspend(tx *gorm.DB) (func() error, error) {
	noop := func() error { return nil }
	if tx.Dialector.Name() != "postgres" {
		return noop, nil
	}
	tables, err := Tables(tx)
	if err != nil {
		return nil, err
	}
	type trigger struct{ table, name string }
	var found []trigger
	for _, t := range tables {
		if !identifier.MatchString(t) {
			return nil, fmt.Errorf("refusing to alter the table name %q", t)
		}
		var names []string
		if err := tx.Raw("SELECT tgname FROM pg_trigger WHERE tgrelid = to_regclass(?) AND tgname LIKE 'grit_append_only%'", t).
			Scan(&names).Error; err != nil {
			return nil, fmt.Errorf("finding the append-only triggers on %s: %w", t, err)
		}
		for _, n := range names {
			found = append(found, trigger{t, n})
		}
	}
	toggle := func(action string) error {
		for _, tr := range found {
			if err := tx.Exec("ALTER TABLE " + tr.table + " " + action + " TRIGGER " + tr.name).Error; err != nil {
				return fmt.Errorf("%s trigger %s on %s: %w", strings.ToLower(action), tr.name, tr.table, err)
			}
		}
		return nil
	}
	if err := toggle("DISABLE"); err != nil {
		return nil, err
	}
	return func() error { return toggle("ENABLE") }, nil
}

const postgresFunction = "CREATE OR REPLACE FUNCTION grit_append_only() RETURNS trigger AS $$\n" +
	"BEGIN\n" +
	"  RAISE EXCEPTION '% is append-only: % refused', TG_TABLE_NAME, TG_OP\n" +
	"    USING ERRCODE = 'restrict_violation';\n" +
	"END;\n" +
	"$$ LANGUAGE plpgsql"

func triggerSQL(dialect, table string) []string {
	switch dialect {
	case "postgres":
		return []string{
			"DROP TRIGGER IF EXISTS grit_append_only ON " + table,
			"CREATE TRIGGER grit_append_only BEFORE UPDATE OR DELETE ON " + table +
				" FOR EACH ROW EXECUTE PROCEDURE grit_append_only()",
			"DROP TRIGGER IF EXISTS grit_append_only_truncate ON " + table,
			"CREATE TRIGGER grit_append_only_truncate BEFORE TRUNCATE ON " + table +
				" FOR EACH STATEMENT EXECUTE PROCEDURE grit_append_only()",
		}
	case "sqlite":
		return []string{
			"CREATE TRIGGER IF NOT EXISTS grit_append_only_update_" + table + " BEFORE UPDATE ON " + table +
				" BEGIN SELECT RAISE(ABORT, '" + table + " is append-only: UPDATE refused'); END",
			"CREATE TRIGGER IF NOT EXISTS grit_append_only_delete_" + table + " BEFORE DELETE ON " + table +
				" BEGIN SELECT RAISE(ABORT, '" + table + " is append-only: DELETE refused'); END",
		}
	case "mysql":
		return []string{
			"DROP TRIGGER IF EXISTS grit_append_only_update_" + table,
			"CREATE TRIGGER grit_append_only_update_" + table + " BEFORE UPDATE ON " + table +
				" FOR EACH ROW SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '" + table + " is append-only: UPDATE refused'",
			"DROP TRIGGER IF EXISTS grit_append_only_delete_" + table,
			"CREATE TRIGGER grit_append_only_delete_" + table + " BEFORE DELETE ON " + table +
				" FOR EACH ROW SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '" + table + " is append-only: DELETE refused'",
		}
	}
	return nil
}
`
}

func appendOnlyTestGo() string {
	return `package appendonly

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/respond"
)

type ledgerRow struct {
	ID     uint
	Amount int
}

// open is an in-memory SQLite on a single connection. :memory: is per
// connection, so a pool of more than one would see empty databases.
func open(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&ledgerRow{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// only makes model the sole registration for the length of a test.
func only(t *testing.T, model interface{}) {
	t.Helper()
	mu.Lock()
	saved := registry
	registry = []interface{}{model}
	mu.Unlock()
	t.Cleanup(func() {
		mu.Lock()
		registry = saved
		mu.Unlock()
	})
}

func TestGORMWritesAreRefusedWithAReason(t *testing.T) {
	db := open(t)
	only(t, &ledgerRow{})
	if err := Install(db); err != nil {
		t.Fatalf("install: %v", err)
	}

	row := ledgerRow{Amount: 100}
	if err := db.Create(&row).Error; err != nil {
		t.Fatalf("creating must still work: %v", err)
	}

	if _, ok := respond.IsRule(db.Model(&row).Update("amount", 1).Error); !ok {
		t.Error("an update through the model was not refused as a rule")
	}
	if _, ok := respond.IsRule(db.Delete(&row).Error); !ok {
		t.Error("a delete through the model was not refused as a rule")
	}
	// GORM Studio's row editor: a table name and a map, no model.
	err := db.Table("ledger_rows").Where("id = ?", row.ID).Updates(map[string]interface{}{"amount": 2}).Error
	if _, ok := respond.IsRule(err); !ok {
		t.Error("an update by table name was not refused")
	}
}

func TestRawSQLIsRefusedByTheTrigger(t *testing.T) {
	db := open(t)
	only(t, &ledgerRow{})
	if err := InstallTriggers(db); err != nil {
		t.Fatalf("triggers: %v", err)
	}

	row := ledgerRow{Amount: 100}
	if err := db.Create(&row).Error; err != nil {
		t.Fatalf("creating must still work: %v", err)
	}
	if err := db.Exec("UPDATE ledger_rows SET amount = 1").Error; err == nil {
		t.Error("a raw UPDATE went through")
	}
	if err := db.Exec("DELETE FROM ledger_rows").Error; err == nil {
		t.Error("a raw DELETE went through")
	}

	var got ledgerRow
	if err := db.First(&got, row.ID).Error; err != nil || got.Amount != 100 {
		t.Errorf("the row changed: %+v, %v", got, err)
	}
}

// grit migrate runs on every deploy, so installing twice must be harmless.
func TestInstallTriggersIsRepeatable(t *testing.T) {
	db := open(t)
	only(t, &ledgerRow{})
	for i := 0; i < 2; i++ {
		if err := InstallTriggers(db); err != nil {
			t.Fatalf("run %d: %v", i+1, err)
		}
	}
}

// Nothing registered means nothing guarded, and a project that has not used
// the flag must not notice the package is there.
func TestUnregisteredTablesAreUntouched(t *testing.T) {
	db := open(t)
	only(t, nil)
	mu.Lock()
	registry = nil
	mu.Unlock()
	if err := Install(db); err != nil {
		t.Fatalf("install: %v", err)
	}
	if err := InstallTriggers(db); err != nil {
		t.Fatalf("triggers: %v", err)
	}
	row := ledgerRow{Amount: 100}
	db.Create(&row)
	if err := db.Model(&row).Update("amount", 1).Error; err != nil {
		t.Errorf("an ordinary table was guarded: %v", err)
	}
}
`
}
