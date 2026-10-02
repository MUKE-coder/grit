package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

func TestRepairSyncWiring(t *testing.T) {
	fresh := strings.ReplaceAll(apiDatabaseGo(), "{{MODULE}}", "demo")
	old := strings.Replace(fresh, syncConnectHook, "", 1)
	old = strings.Replace(old, "\t\"demo/internal/sync\"\n", "", 1)
	if old == fresh {
		t.Fatal("could not rebuild database.go from before the repair")
	}
	out, fixed, warn := repairSyncWiringSource(old, "demo")
	if len(warn) > 0 || len(fixed) != 1 || !strings.Contains(out, "sync.Install(db)") || !strings.Contains(out, "\"demo/internal/sync\"") {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("not valid Go: %v", err)
	}
	if again, fixed, _ := repairSyncWiringSource(out, "demo"); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed database.go again")
	}
	if out, fixed, _ := repairSyncWiringSource(fresh, "demo"); out != fresh || len(fixed) > 0 {
		t.Error("a fresh database.go still needs the repair")
	}
}

func TestRepairUpdatedAtIndex(t *testing.T) {
	generated := strings.Replace(oldGeneratedModel, "*time.Time  `", "*time.Time `", 1)
	generated = strings.Replace(generated, "\tUpdatedAt time.Time      `json:\"updated_at\"`\n", "\tUpdatedAt time.Time      `json:\"updated_at\"`\n", 1)
	out, fixed, _ := repairUpdatedAtIndexSource(generated, "contact.go")
	if len(fixed) != 1 || !strings.Contains(out, "\tUpdatedAt time.Time      `gorm:\"index\" json:\"updated_at\"`\n") {
		t.Fatalf("fixed %v:\n%s", fixed, out)
	}
	if again, fixed, _ := repairUpdatedAtIndexSource(out, "contact.go"); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed the model again")
	}

	blog := blogModelGo()
	oldBlog := strings.Replace(blog, "`gorm:\"index\" json:\"updated_at\"`", "`json:\"updated_at\"`", 1)
	if oldBlog == blog {
		t.Fatal("could not rebuild blog.go from before the repair")
	}
	if out, fixed, _ := repairUpdatedAtIndexSource(oldBlog, "blog.go"); out != blog || len(fixed) != 1 {
		t.Errorf("fixed %v; the repaired blog model differs from a fresh one", fixed)
	}
	if out, fixed, _ := repairUpdatedAtIndexSource(blog, "blog.go"); out != blog || len(fixed) > 0 {
		t.Error("a fresh blog model still needs the repair")
	}

	hand := "package models\n\ntype Note struct {\n\tUpdatedAt time.Time `json:\"updated_at\"`\n}\n"
	if out, fixed, _ := repairUpdatedAtIndexSource(hand, "note.go"); out != hand || len(fixed) > 0 {
		t.Error("a hand-written model was changed")
	}
}

func TestFreshSyncPullTemplates(t *testing.T) {
	handler := apiSyncHandlerGo()
	// The cursor is the request's, so it is parsed in the handler; the query it
	// turns into is the service's since v3.346.0.
	for _, want := range []string{"func parseSyncCursor(", "h.rows().Page(c.Request.Context()"} {
		if !strings.Contains(handler, want) {
			t.Errorf("the sync handler is missing %s", want)
		}
	}
	service := apiSyncServiceGo()
	for _, want := range []string{`Order("updated_at asc, id asc")`, "Unscoped()", "updated_at = ? AND id > ?"} {
		if !strings.Contains(service, want) {
			t.Errorf("the sync service is missing %s", want)
		}
	}
	for _, gone := range []string{"GREATEST(", "MAX(updated_at", "effectiveSyncTime"} {
		if strings.Contains(handler+service, gone) {
			t.Errorf("the sync pull still has %s", gone)
		}
	}
	if !strings.Contains(apiMigrateMainGo(), "sync.BackfillUpdatedAt(db, models.Models()...)") {
		t.Error("cmd/migrate does not backfill updated_at for older deletes")
	}
	for name, src := range map[string]string{"softdelete.go": syncSoftDeleteGo(), "softdelete_test.go": syncSoftDeleteTestGo()} {
		if _, err := format.Source([]byte(src)); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}
