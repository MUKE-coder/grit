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
		filepath.Join(apiRoot, "internal", "models", "erasable.go"):      apiErasableModelsGo(),
		filepath.Join(apiRoot, "internal", "models", "erasable_test.go"): apiErasableModelsTestGo(),
		filepath.Join(apiRoot, "internal", "services", "gdpr.go"):        apiGDPRServiceGo(),
		filepath.Join(apiRoot, "internal", "services", "gdpr_test.go"):   apiGDPRTestGo(),
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

import (
	"context"

	"gorm.io/gorm"

	"{{MODULE}}/internal/erasure"
	"{{MODULE}}/internal/files"
	"{{MODULE}}/internal/storage"
)

// The framework's own tables whose rows exist only to serve a user and carry no
// independent compliance value: hard-deleted when that user is erased.
//
// Resources generated with --owned-by register themselves from their own model
// files, so erasing a user deletes the records they own too. The activity log
// is deliberately absent: its rows hold only a UUID, so anonymizing the user
// anonymizes them, and deleting them would break the audit hash chain.
func init() {
	for _, m := range []interface{}{
		&Session{}, &PasswordResetToken{}, &UserRole{},
		&TwoFactorConfig{}, &TrustedDevice{}, &TOTPPendingToken{},
		&DashboardLayout{}, &Notification{},
	} {
		erasure.Register(m, "user_id")
	}

	// Uploads own files, and deleting the row that names a file is not
	// deleting the file. An erasure that reports success while the documents
	// are still in the bucket is the one outcome the feature exists to prevent
	// (grit#202), and the objects are unreferenced afterwards, so nothing can
	// find them to clean up later either.
	//
	// This covers Grit's own uploads table. A resource generated with
	// --owned-by that has a file column registers its own cleanup from its
	// model file, through files.EraseOwnedFiles: erasing the owner takes the
	// files on those rows too, whoever uploaded them (grit#311). A file column
	// on a model you wrote by hand is yours to register the same way: see the
	// GDPR docs.
	erasure.RegisterWithCleanup(&Upload{}, "user_id", deleteUploadObjects)
}

// deleteUploadObjects removes the stored objects of the uploads about to be
// erased, inside the erasure transaction and before the rows go.
//
// Every object, not every row: an optimised image is the stored file plus its
// renditions plus the original kept for reprocessing, and until v3.544.0 this
// deleted the first of those and reported the erasure complete. The original
// is the untouched upload, so what survived a right-to-erasure request was
// the full-resolution copy of the photo somebody asked to have deleted
// (grit#221, and the same promise as grit#202).
//
// Rows are read in batches rather than all at once: an account with ten
// thousand files should not be a ten-thousand-row slice in memory during a
// transaction.
//
// A storage disk that is not configured is not an error here. The rows still
// name keys, but there is no bucket they could be in, and failing the erasure
// would leave personal data in the database over an object store that was
// never set up.
func deleteUploadObjects(tx *gorm.DB, ownerColumn, userID string) error {
	disk := storage.Disks.Default()
	if disk == nil {
		return nil
	}

	ctx := context.Background()

	// The derived column may not exist yet: grit upgrade rewrites this file
	// and does not run grit migrate. Naming a missing column here would fail
	// the erasure transaction, which is worse than the leak it closes, so the
	// schema is asked first.
	// id is selected and the destination declares it because FindInBatches
	// pages by reading the primary key off the LAST element of the
	// destination slice. A destination without that field panics there, which
	// a slice of plain paths was: this cleanup worked for an account with up
	// to 500 uploads and crashed the erasure for one with more, and 500 is
	// the batch size the batching exists for.
	columns := []string{"id", "path"}
	if files.UploadsRecordDerived(tx) {
		columns = append(columns, "derived")
	}

	// The type tag is load-bearing. GORM cannot infer a column type for a
	// slice that implements Valuer and Scanner, and without the tag it logs
	// "unsupported data type" and runs the query anyway, so every row comes
	// back with Derived unset and the cleanup reports success having deleted
	// only the paths. A wrong answer that looks like a right one, which is
	// the shape of the bug this whole change is about.
	type row struct {
		ID      string
		Path    string
		Derived files.KeyList ` + "`" + `gorm:"type:text"` + "`" + `
	}
	var batch []row
	rows := tx.Unscoped().Model(&Upload{}).
		Where(ownerColumn+" = ?", userID).
		Where("path <> ''")

	return rows.Select(columns).FindInBatches(&batch, 500, func(_ *gorm.DB, _ int) error {
		// batch is the destination slice: read it, not the tx, which carries
		// no query of its own inside the callback.
		if len(batch) == 0 {
			return nil
		}
		keys := make([]string, 0, len(batch)*2)
		for _, r := range batch {
			if r.Path != "" {
				keys = append(keys, r.Path)
			}
			for _, derived := range r.Derived {
				if derived != "" {
					keys = append(keys, derived)
				}
			}
		}
		if len(keys) == 0 {
			return nil
		}
		return disk.DeleteMany(ctx, keys)
	}).Error
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
	model   interface{}
	column  string
	cleanup Cleanup
}

// Cleanup disposes of whatever a row owns outside the database, and runs
// inside the erasure transaction BEFORE those rows are deleted, because what
// it needs is written on them.
//
// It is given the owner column and the user, not the rows, so it can decide
// for itself what to read: a table with a thousand files wants to stream them
// rather than load them.
//
// An error fails the whole erasure, and that is deliberate. A right-to-erasure
// request that reports success while the files are still in the bucket is the
// one outcome the feature exists to prevent (grit#202).
type Cleanup func(tx *gorm.DB, ownerColumn, userID string) error

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
	RegisterWithCleanup(model, ownerColumn, nil)
}

// RegisterWithCleanup is Register for a model whose rows own something outside
// the database: a stored object, a row in another system, a scheduled job.
//
//	erasure.RegisterWithCleanup(&Upload{}, "user_id", deleteStoredObjects)
//
// Without this, erasing a user deletes the row that names a file and leaves
// the file, which is both a failed erasure and an object nothing can find
// again to clean up later.
func RegisterWithCleanup(model interface{}, ownerColumn string, cleanup Cleanup) {
	mu.Lock()
	defer mu.Unlock()
	targets = append(targets, target{model: model, column: ownerColumn, cleanup: cleanup})
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
			// Before the delete: a cleanup reads the rows to learn what they
			// own, and after the delete there is nothing left to read.
			if t.cleanup != nil {
				if err := t.cleanup(tx, t.column, userID); err != nil {
					return nil, 0, fmt.Errorf("cleaning up %s: %w", table, err)
				}
			}
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

// apiErasableModelsTestGo is the behavioural test for the uploads cleanup: a
// real local disk and a real database, because the three string checks in the
// framework's own suite only prove the code says the right thing.
func apiErasableModelsTestGo() string {
	return `package models

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/files"
	"{{MODULE}}/internal/storage"
)

// Erasing a user deletes every object their uploads own, not one per row.
//
// grit#221: an optimised image is the stored file, one object per rendition a
// profile asked for, and the original kept under a private prefix so it can
// be reprocessed later. This cleanup deleted the path column and reported the
// erasure complete, so what survived a right-to-erasure request was the
// untouched full-resolution copy of the photo somebody asked to have deleted.
//
// Nothing would have caught it afterwards either: a rendition has no uploads
// row, so the orphan sweep has no candidate to find.
func TestErasingAUserDeletesEveryObjectTheirUploadsOwn(t *testing.T) {
	db, root := erasureDisk(t)

	objects := []string{
		"uploads/2026/11/a.webp",
		"uploads/2026/11/a-thumb.webp",
		"originals/2026/11/a.jpg",
	}
	for _, key := range objects {
		writeObject(t, root, key)
	}

	if err := db.Exec(
		"INSERT INTO uploads (id, path, derived, user_id) VALUES (?, ?, ?, ?)",
		"u1", objects[0],
		"[\"uploads/2026/11/a-thumb.webp\",\"originals/2026/11/a.jpg\"]",
		"victim").Error; err != nil {
		t.Fatal(err)
	}

	if err := deleteUploadObjects(db, "user_id", "victim"); err != nil {
		t.Fatalf("deleteUploadObjects: %v", err)
	}

	for _, key := range objects {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(key))); err == nil {
			t.Errorf("%s survived the erasure", key)
		}
	}
}

// Another user's files are untouched, which is the direction that costs
// somebody their data rather than some bytes.
func TestErasingAUserLeavesEverybodyElsesObjects(t *testing.T) {
	db, root := erasureDisk(t)

	writeObject(t, root, "uploads/mine.webp")
	writeObject(t, root, "uploads/mine-thumb.webp")
	writeObject(t, root, "uploads/theirs.webp")
	writeObject(t, root, "uploads/theirs-thumb.webp")

	for _, row := range []struct{ id, path, derived, user string }{
		{"u1", "uploads/mine.webp", "[\"uploads/mine-thumb.webp\"]", "victim"},
		{"u2", "uploads/theirs.webp", "[\"uploads/theirs-thumb.webp\"]", "bystander"},
	} {
		if err := db.Exec(
			"INSERT INTO uploads (id, path, derived, user_id) VALUES (?, ?, ?, ?)",
			row.id, row.path, row.derived, row.user).Error; err != nil {
			t.Fatal(err)
		}
	}

	if err := deleteUploadObjects(db, "user_id", "victim"); err != nil {
		t.Fatalf("deleteUploadObjects: %v", err)
	}

	for _, gone := range []string{"uploads/mine.webp", "uploads/mine-thumb.webp"} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(gone))); err == nil {
			t.Errorf("%s survived the erasure", gone)
		}
	}
	for _, kept := range []string{"uploads/theirs.webp", "uploads/theirs-thumb.webp"} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(kept))); err != nil {
			t.Errorf("%s was deleted, and it belongs to somebody else: %v", kept, err)
		}
	}
}

// And a schema without the derived column still erases.
//
// grit upgrade rewrites this file and does not run grit migrate, so this is
// the state an upgraded project is in until somebody migrates. Naming a
// column that is not there would fail the erasure transaction, which is a
// worse outcome than the leak it closes: the request errors and the personal
// data stays in the database.
func TestErasureStillRunsBeforeTheDerivedColumnExists(t *testing.T) {
	db, root := erasureDisk(t)
	if err := db.Exec("ALTER TABLE uploads DROP COLUMN derived").Error; err != nil {
		t.Fatal(err)
	}
	if files.UploadsRecordDerived(db) {
		t.Fatal("the column was dropped and the probe still says it is there")
	}

	writeObject(t, root, "uploads/a.webp")
	if err := db.Exec(
		"INSERT INTO uploads (id, path, user_id) VALUES (?, ?, ?)",
		"u1", "uploads/a.webp", "victim").Error; err != nil {
		t.Fatal(err)
	}

	if err := deleteUploadObjects(db, "user_id", "victim"); err != nil {
		t.Fatalf("the erasure failed on a schema without the column: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "uploads", "a.webp")); err == nil {
		t.Error("the stored file survived, so the fallback deletes nothing at all")
	}
}

// erasureDisk is a database with an uploads table and a local disk registered
// as the default, so deleteUploadObjects has somewhere to delete from.
func erasureDisk(t *testing.T) (*gorm.DB, string) {
	t.Helper()

	root := t.TempDir()
	disk, err := storage.NewLocal(storage.LocalConfig{
		Root:      root,
		PublicURL: "http://localhost:8080/files",
		Secret:    "test-secret",
	})
	if err != nil {
		t.Fatal(err)
	}
	previous := storage.Disks.Default()
	storage.Disks.SetDefault(disk)
	t.Cleanup(func() { storage.Disks.SetDefault(previous) })

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	// One connection: each :memory: connection is its own empty database, and
	// a second one answers "no such table" while the test still passes.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := db.Exec(` + "`" + `CREATE TABLE uploads (
		id TEXT PRIMARY KEY, path TEXT, derived TEXT, user_id TEXT, deleted_at DATETIME
	)` + "`" + `).Error; err != nil {
		t.Fatal(err)
	}
	return db, root
}

func writeObject(t *testing.T, root, key string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
}
`
}
