package appendonly

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"saas/apps/api/internal/respond"
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
type triggerErr string

func (e triggerErr) Error() string { return string(e) }

// The MySQL refusal is recognised, and nothing else is: a syntax error must still
// fail the migration rather than be logged and skipped.
func TestMySQLTriggerPrivilegeIsRecognised(t *testing.T) {
	refused := triggerErr("Error 1419 (HY000): You do not have the SUPER privilege and binary logging is enabled (you *might* want to use the less safe log_bin_trust_function_creators variable)")
	if !triggersNeedPrivilege(refused) {
		t.Error("error 1419 was not recognised, so grit migrate still fails on binlog MySQL")
	}
	if triggersNeedPrivilege(triggerErr("Error 1064 (42000): You have an error in your SQL syntax")) {
		t.Error("an unrelated error was treated as the privilege refusal and would be skipped")
	}
}

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
