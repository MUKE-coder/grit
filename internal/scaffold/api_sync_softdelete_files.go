package scaffold

// syncSoftDeleteGo emits internal/sync/softdelete.go: a soft delete also moves
// updated_at, so sync pull can page on one indexed column (H20 in the
// contact-app review).
func syncSoftDeleteGo() string {
	return `package sync

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
			Where(d + " IS NOT NULL AND " + u + " < " + d).
			UpdateColumn(updated.DBName, gorm.Expr(d)).Error; err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", stmt.Schema.Table, err))
		}
	}
	return errors.Join(errs...)
}
`
}

// syncSoftDeleteTestGo emits internal/sync/softdelete_test.go.
func syncSoftDeleteTestGo() string {
	return `package sync

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type syncNote struct {
	ID        string ` + "`" + `gorm:"primarykey"` + "`" + `
	Title     string
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type syncPlain struct {
	ID        string ` + "`" + `gorm:"primarykey"` + "`" + `
	UpdatedAt time.Time
}

func syncDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&syncNote{}, &syncPlain{}); err != nil {
		t.Fatal(err)
	}
	return db
}

// Every way a handler deletes a row moves updated_at with deleted_at.
func TestSoftDeleteMovesUpdatedAt(t *testing.T) {
	db := syncDB(t)
	if err := Install(db); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-time.Hour)
	for _, id := range []string{"a", "b", "c"} {
		if err := db.Create(&syncNote{ID: id, UpdatedAt: old}).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Delete(&syncNote{}, "id = ?", "a").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Where("id = ?", "b").Delete(&syncNote{}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Delete(&syncNote{ID: "c"}).Error; err != nil {
		t.Fatal(err)
	}
	var rows []syncNote
	if err := db.Unscoped().Order("id").Find(&rows).Error; err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		if !row.DeletedAt.Valid || !row.UpdatedAt.Equal(row.DeletedAt.Time) {
			t.Errorf("%s: updated_at %v, deleted_at %v; a soft delete must move both", row.ID, row.UpdatedAt, row.DeletedAt)
		}
	}
}

// A hard delete, and a model with no soft delete, are left as GORM does them.
func TestOtherDeletesAreUnchanged(t *testing.T) {
	db := syncDB(t)
	if err := Install(db); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&syncNote{ID: "a"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&syncPlain{ID: "p"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Unscoped().Delete(&syncNote{}, "id = ?", "a").Error; err != nil {
		t.Errorf("hard delete: %v", err)
	}
	if err := db.Delete(&syncPlain{}, "id = ?", "p").Error; err != nil {
		t.Errorf("delete without soft delete: %v", err)
	}
	var n int64
	if err := db.Unscoped().Model(&syncNote{}).Count(&n).Error; err != nil || n != 0 {
		t.Errorf("%d notes after a hard delete, err %v", n, err)
	}
}

// Rows deleted before Install get updated_at moved by the backfill, once.
func TestBackfillUpdatedAt(t *testing.T) {
	db := syncDB(t)
	old := time.Now().Add(-2 * time.Hour)
	if err := db.Create(&syncNote{ID: "gone", UpdatedAt: old}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&syncNote{ID: "kept", UpdatedAt: old}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Delete(&syncNote{}, "id = ?", "gone").Error; err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := BackfillUpdatedAt(db, &syncNote{}, &syncPlain{}); err != nil {
			t.Fatal(err)
		}
	}
	var gone, kept syncNote
	if err := db.Unscoped().First(&gone, "id = ?", "gone").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.First(&kept, "id = ?", "kept").Error; err != nil {
		t.Fatal(err)
	}
	if !gone.UpdatedAt.Equal(gone.DeletedAt.Time) {
		t.Errorf("a deleted row kept updated_at %v, deleted at %v", gone.UpdatedAt, gone.DeletedAt.Time)
	}
	if !kept.UpdatedAt.Equal(old) && kept.UpdatedAt.Sub(old).Abs() > time.Millisecond {
		t.Errorf("a live row's updated_at moved to %v", kept.UpdatedAt)
	}
}
`
}
