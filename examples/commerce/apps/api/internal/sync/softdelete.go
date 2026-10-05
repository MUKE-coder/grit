package sync

import (
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var deletedAtType = reflect.TypeOf(gorm.DeletedAt{})

// Install makes a soft delete also set updated_at, to the same instant as
// deleted_at, on every model that has both.
//
// Sync pull pages on updated_at. A soft delete used to set only deleted_at, so
// pull ordered by the later of the two: an expression no index can serve, which
// sorted the whole table on every pull, and one MySQL has no function for, so
// pull failed there outright.
func Install(db *gorm.DB) error {
	return db.Callback().Delete().Before("gorm:delete").Register("sync:soft_delete_touches_updated_at", touchUpdatedAt)
}

func touchUpdatedAt(tx *gorm.DB) {
	stmt := tx.Statement
	if tx.Error != nil || stmt.Schema == nil || stmt.Unscoped {
		return
	}
	deleted, updated := stmt.Schema.LookUpField("DeletedAt"), stmt.Schema.LookUpField("UpdatedAt")
	if deleted == nil || updated == nil || deleted.FieldType != deletedAtType {
		return
	}
	// GORM's soft delete puts SET deleted_at into this clause as it builds the
	// statement, and keeps a Builder already on the clause. The Builder writes
	// updated_at beside it, with the same value.
	set := stmt.Clauses["SET"]
	set.Builder = func(c clause.Clause, b clause.Builder) {
		assignments, _ := c.Expression.(clause.Set)
		if len(assignments) == 1 && assignments[0].Column.Name == deleted.DBName {
			assignments = append(assignments, clause.Assignment{Column: clause.Column{Name: updated.DBName}, Value: assignments[0].Value})
		}
		b.WriteString("SET ")
		assignments.Build(b)
	}
	stmt.Clauses["SET"] = set
}

// BackfillUpdatedAt moves updated_at up to deleted_at on rows soft-deleted
// before Install existed, so pull carries those deletes too. It changes nothing
// on a second run; grit migrate runs it.
func BackfillUpdatedAt(db *gorm.DB, models ...interface{}) error {
	var errs []error
	for _, model := range models {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			errs = append(errs, fmt.Errorf("reading %T: %w", model, err))
			continue
		}
		deleted, updated := stmt.Schema.LookUpField("DeletedAt"), stmt.Schema.LookUpField("UpdatedAt")
		if deleted == nil || updated == nil || deleted.FieldType != deletedAtType || !db.Migrator().HasTable(model) {
			continue
		}
		d, u := stmt.Quote(deleted.DBName), stmt.Quote(updated.DBName)
		if err := db.Unscoped().Model(model).
			Where(d+" IS NOT NULL AND "+u+" < "+d).
			UpdateColumn(updated.DBName, gorm.Expr(d)).Error; err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", stmt.Schema.Table, err))
		}
	}
	return errors.Join(errs...)
}
