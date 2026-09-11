package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// writeAuditChainFiles brings an upgraded project's activity-log chain up to
// date.
//
// Upgrade does not regenerate API code in general. This is the exception for
// the same reason as the crypto and backup packages: before v3.215.0 the chain
// hashed a timestamp more precise than Postgres and MySQL store, so it failed
// verification on the first row of every such deployment, and a project
// carrying the old writer keeps producing entries that can never verify.
// Manifest-guarded, so an edited copy is reported rather than replaced.
func writeAuditChainFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	write := func(path, content string) error {
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", opts.Module())); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
		return nil
	}

	// The model first. The new audit package reads fields it adds, so if an
	// edited copy is held back, writing the rest would stop the API compiling.
	model := filepath.Join(apiRoot, "internal", "models", "activity_log.go")
	if err := write(model, apiActivityLogModelGo()); err != nil {
		return err
	}
	if !fileContains(model, "ResourceIDs") {
		fmt.Println("  ⚠ models/activity_log.go has been edited, so the audit chain fix was not applied.\n" +
			"    Run grit upgrade --diff to see what it needs.")
		return nil
	}

	for _, f := range []struct{ path, content string }{
		{filepath.Join(apiRoot, "internal", "audit", "audit.go"), apiAuditGo()},
		{filepath.Join(apiRoot, "internal", "audit", "chain_test.go"), apiAuditChainTestGo()},
		{filepath.Join(apiRoot, "internal", "middleware", "activity.go"), apiActivityMiddlewareGo()},
		{filepath.Join(apiRoot, "internal", "middleware", "security_log.go"), securityLogGo()},
		{filepath.Join(apiRoot, "internal", "handlers", "activity.go"), apiActivityHandlerGo()},
	} {
		if err := write(f.path, f.content); err != nil {
			return err
		}
	}

	// Mounting the route before the handler has the method would stop the
	// API compiling, so an edited handler that was held back gets no route.
	if !fileContains(filepath.Join(apiRoot, "internal", "handlers", "activity.go"), "func (h *ActivityHandler) Reseal(") {
		return nil
	}
	return ensureAuditResealRoute(apiRoot)
}

func fileContains(path, needle string) bool {
	data, err := os.ReadFile(path)
	return err == nil && strings.Contains(string(data), needle)
}

const auditResealRoute = "\t\tadmin.POST(\"/admin/activity/reseal\", activityHandler.Reseal)\n"

// ensureAuditResealRoute mounts the reseal endpoint beside the integrity check
// in an existing routes.go. Idempotent; says where the line goes when the file
// has been reshaped past recognition.
func ensureAuditResealRoute(apiRoot string) error {
	path := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil // no routes.go: nothing mounts the activity handler here
	}
	crlf := strings.Contains(string(raw), "\r\n")
	content := strings.ReplaceAll(string(raw), "\r\n", "\n")
	if strings.Contains(content, "activityHandler.Reseal") {
		return nil
	}
	const anchor = "\t\tadmin.GET(\"/admin/activity/integrity\", activityHandler.VerifyIntegrity)\n"
	if !strings.Contains(content, anchor) {
		fmt.Printf("  ⚠ Could not find the activity integrity route in routes.go. Add this beside it\n" +
			"    to be able to reseal a chain that cannot verify:\n\n" + auditResealRoute + "\n")
		return nil
	}
	content = strings.Replace(content, anchor, anchor+auditResealRoute, 1)
	if crlf {
		content = strings.ReplaceAll(content, "\n", "\r\n")
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	manifest.Refresh(path)
	fmt.Println("  ✓ Mounted POST /api/admin/activity/reseal")
	return nil
}

// apiAuditChainTestGo ships with every project, so the chain is checked where
// it runs and not only in Grit's own tests.
func apiAuditChainTestGo() string {
	return `package audit

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"{{MODULE}}/internal/models"
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
`
}
