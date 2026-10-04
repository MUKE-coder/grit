// Package migrate records what each migration run changed, so a bad run can be
// undone.
//
// Grit migrates with GORM AutoMigrate, which has no down scripts. It only ever
// adds: a table, a column, an index. Nothing it does drops a column. That makes
// the reverse of a run computable rather than hand-written: snapshot the schema
// before and after, store the difference, and a rollback is dropping what the
// run added, newest first.
//
// What no rollback can do is bring data back. Dropping a column that a run
// added loses everything written into it since, which is why "grit migrate down"
// prints every statement and will not run without --yes.
package migrate

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// The three kinds of change AutoMigrate makes. Each has an exact inverse.
const (
	KindTable  = "table"
	KindColumn = "column"
	KindIndex  = "index"
)

// Change is one thing a run added.
type Change struct {
	Kind    string // table, column or index
	OnTable string
	Name    string // the column or index; empty for a table
}

// String is how a change is shown to somebody about to undo it.
func (c Change) String() string {
	switch c.Kind {
	case KindTable:
		return "table " + c.OnTable
	case KindColumn:
		return "column " + c.OnTable + "." + c.Name
	default:
		return "index " + c.Name + " on " + c.OnTable
	}
}

// Run is one migration, as recorded.
type Run struct {
	ID           string
	AppliedAt    time.Time
	GritVersion  string
	Baseline     bool
	RolledBackAt *time.Time
}

// TableName pins the table, so renaming the type cannot orphan the history.
func (Run) TableName() string { return "grit_migrations" }

// RunChange is one change belonging to a run.
//
// Rows rather than a JSON blob: the history stays readable from psql, and there
// is no long-text column to declare.
type RunChange struct {
	ID      uint
	RunID   string
	Kind    string
	OnTable string
	Name    string
}

// TableName pins the table.
func (RunChange) TableName() string { return "grit_migration_changes" }

// historyTables are not part of the schema the history describes, so they never
// appear in a snapshot and a rollback can never drop them.
var historyTables = map[string]bool{
	"grit_migrations":        true,
	"grit_migration_changes": true,
}

// internalTable reports tables the database keeps for itself.
//
// SQLite reports sqlite_sequence, which exists because the history's own id is
// an AUTOINCREMENT. Counting it means a snapshot of an empty database is not
// empty, the first run is not recorded as a baseline, and a rollback offers to
// drop the whole schema. That is the bug this function exists for.
func internalTable(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasPrefix(lower, "sqlite_") ||
		strings.HasPrefix(lower, "pg_") ||
		strings.HasPrefix(lower, "sql_")
}

// Schema is what the database looks like: the columns and indexes of each table.
type Schema map[string]Table

// Table is one table's shape.
type Table struct {
	Columns map[string]bool
	Indexes map[string]bool
}

// EnsureHistory creates the history tables. Safe to call on every run.
func EnsureHistory(db *gorm.DB) error {
	if err := db.AutoMigrate(&Run{}, &RunChange{}); err != nil {
		return fmt.Errorf("creating the migration history: %w", err)
	}
	return nil
}

// Snapshot reads the shape of every table except the history's own.
func Snapshot(db *gorm.DB) (Schema, error) {
	tables, err := db.Migrator().GetTables()
	if err != nil {
		return nil, fmt.Errorf("listing tables: %w", err)
	}
	out := Schema{}
	for _, name := range tables {
		if historyTables[name] || internalTable(name) {
			continue
		}
		shape := Table{Columns: map[string]bool{}, Indexes: map[string]bool{}}
		columns, err := db.Migrator().ColumnTypes(name)
		if err != nil {
			return nil, fmt.Errorf("reading the columns of %s: %w", name, err)
		}
		for _, column := range columns {
			shape.Columns[column.Name()] = true
		}
		// Indexes are not worth failing a migration over: the dialects report
		// them differently, and dropping a column drops its indexes with it.
		if indexes, err := db.Migrator().GetIndexes(name); err == nil {
			for _, index := range indexes {
				shape.Indexes[index.Name()] = true
			}
		}
		out[name] = shape
	}
	return out, nil
}

// Diff is what was added between two snapshots: tables first, then columns, then
// indexes. A new table's own columns and indexes are not listed separately,
// because dropping the table takes them with it.
func Diff(before, after Schema) []Change {
	var tables, columns, indexes []Change
	for name, shape := range after {
		previous, existed := before[name]
		if !existed {
			tables = append(tables, Change{Kind: KindTable, OnTable: name})
			continue
		}
		for column := range shape.Columns {
			if !previous.Columns[column] {
				columns = append(columns, Change{Kind: KindColumn, OnTable: name, Name: column})
			}
		}
		for index := range shape.Indexes {
			if !previous.Indexes[index] {
				indexes = append(indexes, Change{Kind: KindIndex, OnTable: name, Name: index})
			}
		}
	}
	sortChanges(tables)
	sortChanges(columns)
	sortChanges(indexes)
	return append(append(tables, columns...), indexes...)
}

func sortChanges(changes []Change) {
	sort.Slice(changes, func(i, j int) bool {
		if changes[i].OnTable != changes[j].OnTable {
			return changes[i].OnTable < changes[j].OnTable
		}
		return changes[i].Name < changes[j].Name
	})
}

// Record stores a run and its changes.
//
// baseline marks the run that built the schema in the first place. Undoing that
// is not a rollback, it is dropping the database, so the command refuses it and
// points at --fresh instead.
func Record(db *gorm.DB, changes []Change, gritVersion string, baseline bool) (Run, error) {
	now := time.Now().UTC()
	run := Run{
		AppliedAt:   now,
		GritVersion: gritVersion,
		Baseline:    baseline,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		id, err := nextRunID(tx, now)
		if err != nil {
			return err
		}
		run.ID = id
		if err := tx.Create(&run).Error; err != nil {
			return err
		}
		for _, change := range changes {
			row := RunChange{RunID: run.ID, Kind: change.Kind, OnTable: change.OnTable, Name: change.Name}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return run, fmt.Errorf("recording the migration: %w", err)
	}
	return run, nil
}

// runIDFormat is a run's ID: when it was recorded, to the millisecond, which
// sorts in the order the runs happened.
const runIDFormat = "20060102-150405.000"

// nextRunID is now as a run ID, moved past the latest recorded run when it is not
// later. Two runs in the same millisecond, which a fast test makes and a clock
// stepping back can too, used to collide on the primary key.
func nextRunID(tx *gorm.DB, now time.Time) (string, error) {
	id := now.Format(runIDFormat)
	var latest Run
	res := tx.Order("id desc").Limit(1).Find(&latest)
	if res.Error != nil {
		return "", fmt.Errorf("reading the latest run: %w", res.Error)
	}
	if res.RowsAffected == 0 || id > latest.ID {
		return id, nil
	}
	last, err := time.Parse(runIDFormat, latest.ID)
	if err != nil {
		return "", fmt.Errorf("reading the id of the latest run, %q: %w", latest.ID, err)
	}
	return last.Add(time.Millisecond).Format(runIDFormat), nil
}

// ChangesOf reads back the changes of one run, in the order they were recorded.
func ChangesOf(db *gorm.DB, runID string) ([]Change, error) {
	var rows []RunChange
	if err := db.Where("run_id = ?", runID).Order("id asc").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("reading the changes of %s: %w", runID, err)
	}
	changes := make([]Change, 0, len(rows))
	for _, row := range rows {
		changes = append(changes, Change{Kind: row.Kind, OnTable: row.OnTable, Name: row.Name})
	}
	return changes, nil
}

// History returns the most recent runs, newest first.
func History(db *gorm.DB, limit int) ([]Run, error) {
	var runs []Run
	query := db.Order("id desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&runs).Error; err != nil {
		return nil, fmt.Errorf("reading the migration history: %w", err)
	}
	return runs, nil
}

// Undone returns the runs that can still be rolled back, newest first.
func Undone(db *gorm.DB, steps int) ([]Run, error) {
	var runs []Run
	query := db.Where("rolled_back_at IS NULL").Order("id desc")
	if steps > 0 {
		query = query.Limit(steps)
	}
	if err := query.Find(&runs).Error; err != nil {
		return nil, fmt.Errorf("reading the migration history: %w", err)
	}
	return runs, nil
}

// Statements is the SQL a rollback would run, in the order it runs it: indexes
// first, then columns, then tables, so nothing is dropped out from under
// something else.
func Statements(db *gorm.DB, changes []Change) []string {
	var statements []string
	for _, kind := range []string{KindIndex, KindColumn, KindTable} {
		for i := len(changes) - 1; i >= 0; i-- {
			change := changes[i]
			if change.Kind != kind {
				continue
			}
			statements = append(statements, statementFor(db, change))
		}
	}
	return statements
}

func statementFor(db *gorm.DB, change Change) string {
	switch change.Kind {
	case KindIndex:
		if db.Dialector.Name() == "mysql" {
			return "ALTER TABLE " + quote(db, change.OnTable) + " DROP INDEX " + quote(db, change.Name)
		}
		return "DROP INDEX IF EXISTS " + quote(db, change.Name)
	case KindColumn:
		// No IF EXISTS: MySQL and SQLite do not accept it for a column, which is
		// why every statement is guarded by a HasColumn check before it runs.
		return "ALTER TABLE " + quote(db, change.OnTable) + " DROP COLUMN " + quote(db, change.Name)
	default:
		if db.Dialector.Name() == "postgres" {
			return "DROP TABLE IF EXISTS " + quote(db, change.OnTable) + " CASCADE"
		}
		return "DROP TABLE IF EXISTS " + quote(db, change.OnTable)
	}
}

// quote wraps an identifier the way the dialect wants it. "\x60" is a backtick,
// which is how MySQL quotes.
func quote(db *gorm.DB, identifier string) string {
	if db.Dialector.Name() == "mysql" {
		return "\x60" + strings.ReplaceAll(identifier, "\x60", "\x60\x60") + "\x60"
	}
	return "\"" + strings.ReplaceAll(identifier, "\"", "\"\"") + "\""
}

// Rollback drops what a run added and marks the run rolled back. With dryRun it
// returns the statements and changes nothing.
//
// Every statement is skipped if the thing it drops is already gone, so a
// rollback that failed halfway can be run again. On Postgres the whole thing is
// one transaction; MySQL commits each DDL statement as it goes, so a failure
// there leaves the earlier drops applied and the run stays open for a re-run.
func Rollback(db *gorm.DB, run Run, changes []Change, dryRun bool) ([]string, error) {
	if dryRun {
		return Statements(db, changes), nil
	}

	applied := make([]string, 0, len(changes))
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, kind := range []string{KindIndex, KindColumn, KindTable} {
			for i := len(changes) - 1; i >= 0; i-- {
				change := changes[i]
				if change.Kind != kind || alreadyGone(tx, change) {
					continue
				}
				statement := statementFor(tx, change)
				if err := tx.Exec(statement).Error; err != nil {
					return fmt.Errorf("%s: %w", statement, err)
				}
				applied = append(applied, statement)
			}
		}
		now := time.Now().UTC()
		return tx.Model(&Run{}).Where("id = ?", run.ID).Update("rolled_back_at", now).Error
	})
	if err != nil {
		return applied, err
	}
	return applied, nil
}

// alreadyGone reports whether a change has already been undone, so re-running a
// half-applied rollback is harmless.
func alreadyGone(db *gorm.DB, change Change) bool {
	switch change.Kind {
	case KindTable:
		return !db.Migrator().HasTable(change.OnTable)
	case KindColumn:
		return !db.Migrator().HasTable(change.OnTable) || !db.Migrator().HasColumn(change.OnTable, change.Name)
	default:
		return !db.Migrator().HasTable(change.OnTable) || !db.Migrator().HasIndex(change.OnTable, change.Name)
	}
}

// Reset forgets the history, for --fresh: the schema it described is gone.
func Reset(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&RunChange{}, &Run{}); err != nil {
		return fmt.Errorf("clearing the migration history: %w", err)
	}
	return EnsureHistory(db)
}
