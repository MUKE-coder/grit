package audit

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"commerce/apps/api/internal/models"
)

func openChain(t *testing.T) *gorm.DB {
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
	if err := db.AutoMigrate(&models.ActivityLog{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func write(t *testing.T, db *gorm.DB, n int) {
	t.Helper()
	batch := make([]models.ActivityLog, n)
	for i := range batch {
		batch[i] = models.ActivityLog{UserID: "u1", Method: "POST", Path: "/api/v1/notes", Status: 201}
	}
	if err := appendBatch(db, batch); err != nil {
		t.Fatalf("append: %v", err)
	}
}

// Postgres keeps microseconds and MySQL milliseconds, and the hash covers
// created_at. Before v3.215.0 the stamp carried nanoseconds, the database
// rounded them away, and the chain failed verification on its first row.
func TestStampSurvivesTheDatabasesPrecision(t *testing.T) {
	prev := time.Time{}
	for i := 0; i < 50; i++ {
		s := nextStamp(prev)
		if !s.Equal(s.Round(time.Millisecond)) {
			t.Fatalf("stamp %v has precision a database would round away", s)
		}
		if !prev.IsZero() && !s.After(prev) {
			t.Fatalf("stamp %v is not after %v, so verify order would differ from write order", s, prev)
		}
		prev = s
	}
}

func TestChainVerifies(t *testing.T) {
	db := openChain(t)
	write(t, db, 40)
	if err := AppendChained(db, &models.ActivityLog{UserID: "u2", Method: "SECURITY", Path: "security.login.failed"}); err != nil {
		t.Fatalf("append chained: %v", err)
	}
	write(t, db, 10)
	status, err := VerifyChain(context.Background(), db)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !status.Valid || status.TotalEntries != 51 {
		t.Fatalf("chain did not verify: %+v", status)
	}
}

// A changed row is found, a reseal has to name it, and the reseal leaves its
// own record in the chain.
func TestResealNamesTheBreakAndRecordsItself(t *testing.T) {
	db := openChain(t)
	write(t, db, 5)
	var third models.ActivityLog
	db.Order("created_at asc, id asc").Offset(2).First(&third)
	db.Exec("UPDATE activity_logs SET status = 500 WHERE id = ?", third.ID)

	status, _ := VerifyChain(context.Background(), db)
	if status.Valid || status.BrokenAtID != third.ID {
		t.Fatalf("the changed row was not found: %+v", status)
	}
	if _, err := Reseal(context.Background(), db, "not-the-break", "admin", "127.0.0.1", "test"); err == nil {
		t.Fatal("a reseal that did not name the first bad entry was accepted")
	}
	n, err := Reseal(context.Background(), db, third.ID, "admin", "127.0.0.1", "test")
	if err != nil {
		t.Fatalf("reseal: %v", err)
	}
	if n != 3 {
		t.Errorf("resealed %d entries, want the 3 from the break onward", n)
	}
	status, _ = VerifyChain(context.Background(), db)
	if !status.Valid {
		t.Fatalf("still broken after the reseal: %+v", status)
	}
	var last models.ActivityLog
	db.Order("created_at desc, id desc").First(&last)
	if last.Method != "SECURITY" || !strings.HasPrefix(last.Path, "audit.chain.resealed") || last.UserID != "admin" {
		t.Errorf("the reseal did not record itself: %+v", last)
	}
	if _, err := Reseal(context.Background(), db, third.ID, "admin", "", ""); err == nil {
		t.Error("a reseal of a chain that verifies was accepted")
	}
}
