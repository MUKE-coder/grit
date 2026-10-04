// Package appendonly makes a table refuse UPDATE and DELETE.
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

	"library/apps/api/internal/respond"
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
				// MySQL with binary logging on, which is the default and the norm on
				// RDS, refuses CREATE TRIGGER from a user without SUPER. Failing here
				// failed the whole migration, roles included, so --append-only made
				// grit migrate unusable on ordinary production MySQL. The GORM guard
				// still holds for everything that goes through the API; say what is
				// missing and how to add it, the way the no-trigger dialect branch
				// above does.
				if dialect == "mysql" && triggersNeedPrivilege(err) {
					log.Printf("append-only: MySQL refused the trigger on %s: binary logging is on and this user lacks SUPER. "+
						"The API still cannot change it (GORM refuses), but raw SQL can. To add the database guard, run "+
						"SET GLOBAL log_bin_trust_function_creators = 1 (on RDS, set it in the parameter group) and migrate again", table)
					break
				}
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

// triggersNeedPrivilege recognises MySQL error 1419: "You do not have the SUPER
// privilege and binary logging is enabled". Matched on the text rather than the
// driver's error type, so this package needs no MySQL import in a project that
// runs Postgres or SQLite.
func triggersNeedPrivilege(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "1419") || strings.Contains(msg, "binary logging is enabled")
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
