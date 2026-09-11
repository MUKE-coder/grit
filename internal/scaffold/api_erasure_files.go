package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteErasurePackage writes internal/erasure under apiRoot. overwrite=false
// leaves an existing copy alone, for the generator on an older project.
func WriteErasurePackage(apiRoot, module string, overwrite bool) error {
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "erasure", "erasure.go"):      apiErasureGo(),
		filepath.Join(apiRoot, "internal", "erasure", "erasure_test.go"): apiErasureTestGo(),
	}
	for path, content := range files {
		if !overwrite {
			if _, err := os.Stat(path); err == nil {
				continue
			}
		}
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", module)); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

// writeErasureFiles brings an upgraded project's right-to-erasure up to date:
// the erasure package, the framework tables' registration, and the GDPR
// service that uses them. Only called where the GDPR service already exists,
// and manifest-guarded like every upgrade write, so an edited copy is reported
// rather than replaced.
func writeErasureFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	if err := WriteErasurePackage(apiRoot, opts.Module(), true); err != nil {
		return err
	}
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "models", "erasable.go"):    apiErasableModelsGo(),
		filepath.Join(apiRoot, "internal", "services", "gdpr.go"):      apiGDPRServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "gdpr_test.go"): apiGDPRTestGo(),
	}
	for path, content := range files {
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", opts.Module())); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

// apiErasableModelsGo registers the framework's own tables for erasure.
func apiErasableModelsGo() string {
	return `package models

import "{{MODULE}}/internal/erasure"

// The framework's own tables whose rows exist only to serve a user and carry no
// independent compliance value: hard-deleted when that user is erased.
//
// Resources generated with --owned-by register themselves from their own model
// files, so erasing a user deletes the records they own too. The activity log
// is deliberately absent: its rows hold only a UUID, so anonymizing the user
// anonymizes them, and deleting them would break the audit hash chain.
func init() {
	for _, m := range []interface{}{
		&Upload{}, &Session{}, &PasswordResetToken{}, &UserRole{},
		&TwoFactorConfig{}, &TrustedDevice{}, &TOTPPendingToken{},
		&DashboardLayout{}, &Notification{},
	} {
		erasure.Register(m, "user_id")
	}
}
`
}

func apiErasureGo() string {
	return `// Package erasure is what a right-to-erasure request deletes, and the deletion.
//
// A registry rather than a list, so a table joins by declaring itself: the
// framework's own tables from models/erasable.go, and every resource generated
// with --owned-by from its model file. It used to be nine framework tables
// named in the GDPR service, and erasing a user left every app record they
// owned in place while the deletion journal certified the erasure complete.
//
// It imports nothing from the app, because the models import it.
package erasure

import (
	"fmt"
	"regexp"
	"sync"

	"gorm.io/gorm"
)

type target struct {
	model  interface{}
	column string
}

var (
	mu      sync.RWMutex
	targets []target

	// identifier is what an owner column has to look like: it goes into the
	// WHERE clause as text, not as a bound parameter.
	identifier = regexp.MustCompile("^[a-z_][a-z0-9_]*$")
)

// Register marks rows of model as belonging to the user named in ownerColumn,
// so erasing that user deletes them. Call it from the model file's init().
func Register(model interface{}, ownerColumn string) {
	mu.Lock()
	defer mu.Unlock()
	targets = append(targets, target{model: model, column: ownerColumn})
}

// Scrub hard-deletes every registered row belonging to userID and anonymizes
// the user row in place, inside tx. It returns what was deleted, per table.
//
// Unscoped, because several models carry gorm.DeletedAt and a normal Delete
// only soft-deletes: the row, and the personal data in it, would remain.
// Tables that do not exist are skipped, so a database migrated with a subset of
// the models, a test's for one, can still run an erasure.
func Scrub(tx *gorm.DB, userID string) (map[string]int64, int, error) {
	mu.RLock()
	list := append([]target(nil), targets...)
	mu.RUnlock()

	counts := map[string]int64{}
	total := 0
	for _, t := range list {
		if !identifier.MatchString(t.column) {
			return nil, 0, fmt.Errorf("refusing the owner column %q", t.column)
		}
		stmt := &gorm.Statement{DB: tx}
		if err := stmt.Parse(t.model); err != nil {
			return nil, 0, fmt.Errorf("resolving %T: %w", t.model, err)
		}
		table := stmt.Schema.Table
		if !tx.Migrator().HasTable(table) {
			continue
		}
		var n int64
		if err := tx.Unscoped().Model(t.model).Where(t.column+" = ?", userID).Count(&n).Error; err != nil {
			return nil, 0, fmt.Errorf("counting %s: %w", table, err)
		}
		if n > 0 {
			if err := tx.Unscoped().Where(t.column+" = ?", userID).Delete(t.model).Error; err != nil {
				return nil, 0, fmt.Errorf("deleting %s: %w", table, err)
			}
		}
		counts[table] += n
		total += int(n)
	}

	// Anonymize the user row in place: scrub every personal column, keep the id
	// so references resolve to a tombstone, and keep a unique, non-routable
	// email so the unique index stays satisfiable if the row is read again.
	// By table name, so no model hook runs over data that is being destroyed.
	if err := tx.Table("users").Where("id = ?", userID).Updates(map[string]interface{}{
		"first_name":  "Erased",
		"last_name":   "User",
		"email":       "erased-" + userID + "@deleted.invalid",
		"password":    "",
		"avatar":      "",
		"job_title":   "",
		"bio":         "",
		"ip_address":  "",
		"mac_address": "",
		"google_id":   "",
		"github_id":   "",
		"active":      false,
		"role":        "USER",
	}).Error; err != nil {
		return nil, 0, fmt.Errorf("anonymizing user: %w", err)
	}
	return counts, total, nil
}
`
}

func apiErasureTestGo() string {
	return `package erasure

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
`
}
