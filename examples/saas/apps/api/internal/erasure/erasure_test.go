package erasure

import (
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type testUser struct {
	ID         string
	FirstName  string
	LastName   string
	Email      string
	Password   string
	Avatar     string
	JobTitle   string
	Bio        string
	IPAddress  string
	MACAddress string
	GoogleID   string
	GithubID   string
	Active     bool
	Role       string
}

func (testUser) TableName() string { return "users" }

type ownedNote struct {
	ID     uint
	UserID string
	Body   string
}

// never is registered and never migrated, as in a test database built from a
// subset of the models.
type never struct {
	ID     uint
	UserID string
}

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
	if err := db.AutoMigrate(&testUser{}, &ownedNote{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	mu.Lock()
	saved := targets
	targets = nil
	mu.Unlock()
	t.Cleanup(func() {
		mu.Lock()
		targets = saved
		mu.Unlock()
	})
	return db
}

// Erasing a user deletes what they own and leaves everyone else's alone.
// Before the registry, rows in an owned resource survived the erasure.
func TestScrubDeletesOwnedRowsAndAnonymizes(t *testing.T) {
	db := open(t)
	Register(&ownedNote{}, "user_id")
	Register(&never{}, "user_id")

	db.Create(&testUser{ID: "u1", FirstName: "Real", Email: "real@example.com", Active: true})
	db.Create(&testUser{ID: "u2", FirstName: "Other", Email: "other@example.com", Active: true})
	db.Create(&ownedNote{UserID: "u1", Body: "diagnosis"})
	db.Create(&ownedNote{UserID: "u1", Body: "prescription"})
	db.Create(&ownedNote{UserID: "u2", Body: "someone else's"})

	counts, total, err := Scrub(db, "u1")
	if err != nil {
		t.Fatalf("scrub: %v", err)
	}
	if total != 2 || counts["owned_notes"] != 2 {
		t.Errorf("deleted %d (%v), want the 2 rows u1 owned", total, counts)
	}

	var left []ownedNote
	db.Find(&left)
	if len(left) != 1 || left[0].UserID != "u2" {
		t.Errorf("rows left: %+v, want only u2's", left)
	}

	var u1, u2 testUser
	db.First(&u1, "id = ?", "u1")
	db.First(&u2, "id = ?", "u2")
	if !strings.HasPrefix(u1.Email, "erased-") || u1.FirstName != "Erased" || u1.Active {
		t.Errorf("u1 was not anonymized: %+v", u1)
	}
	if u2.Email != "other@example.com" {
		t.Errorf("u2 was touched: %+v", u2)
	}
}

func TestScrubRefusesAnUnsafeColumn(t *testing.T) {
	db := open(t)
	Register(&ownedNote{}, "user_id; DROP TABLE users")
	if _, _, err := Scrub(db, "u1"); err == nil {
		t.Fatal("an owner column that is not an identifier was accepted")
	}
}
