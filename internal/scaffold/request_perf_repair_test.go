package scaffold

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func mustFormatGo(t *testing.T, name, src string) {
	t.Helper()
	if _, err := format.Source([]byte(strings.ReplaceAll(src, "{{MODULE}}", "example.com/app"))); err != nil {
		t.Errorf("%s is not valid Go: %v", name, err)
	}
}

// Every template holds its handler database calls to the request's context.
func TestTemplatesBindRequestContext(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") || name == "request_perf_repair.go" {
			continue
		}
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, n := bindRequestContext(string(raw)); n != 0 {
			t.Errorf("%s has %d handler database calls without the request's context", name, n)
		}
	}
}

func TestBindRequestContext(t *testing.T) {
	src := `package handlers

func (h *NoteHandler) List(c *gin.Context) {
	var notes []Note
	h.DB.Where("a = ?", 1).Find(&notes)
	h.DB.WithContext(c.Request.Context()).Count(&n)
	stats, _ := services.ComputeResourceStats(h.DB, "notes", filter)
	go func() {
		h.DB.Create(&Audit{})
	}()
	// h.DB.Delete in a comment stays
	_ = "h.DB.Find in a string stays"
}

func helper(db *gorm.DB) {
	h.DB.Find(nil)
}

func (h *NoteHandler) load(c *gin.Context, id string) (*Note, error) {
	var note Note
	return &note, h.DB.Where("id = ?", id).First(&note).Error
}

func (h *NoteHandler) record(note *Note) error {
	return h.DB.Save(note).Error
}
`
	out, n := bindRequestContext(src)
	if n != 3 {
		t.Fatalf("bound %d calls, want 3:\n%s", n, out)
	}
	if !strings.Contains(out, "return h.DB.Save(note).Error") {
		t.Error("a method with no request was bound to one")
	}
	if cleaned, ok := dropUnusedImport("package x\n\nimport (\n\t\"fmt\"\n\t\"time\"\n)\n\nfunc f() { fmt.Println() }\n", "time"); !ok || strings.Contains(cleaned, "\"time\"") {
		t.Errorf("dropUnusedImport kept an unused import:\n%s", cleaned)
	}
	if _, ok := dropUnusedImport("package x\n\nimport (\n\t\"time\"\n)\n\nvar d = time.Second\n", "time"); ok {
		t.Error("dropUnusedImport removed an import that is used")
	}
	for _, want := range []string{
		`h.DB.WithContext(c.Request.Context()).Where("a = ?", 1)`,
		`h.DB.WithContext(c.Request.Context()).Where("id = ?", id)`,
		`services.ComputeResourceStats(h.DB.WithContext(c.Request.Context()), "notes", filter)`,
		"\t\th.DB.Create(&Audit{})",
		"\th.DB.Find(nil)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if again, n := bindRequestContext(out); n != 0 || again != out {
		t.Error("binding is not idempotent")
	}
	mustFormatGo(t, "the bound fixture", out)
}

func TestWorkersAndAuthLookupBindContext(t *testing.T) {
	raw, _ := os.ReadFile("api_jobs_files.go")
	for old, bound := range workerDBCalls {
		// The audit prune's transaction moved into audit.Prune, which takes ctx.
		if old == "return deps.DB.Transaction(" && strings.Contains(jobsWorkersGo(), "audit.Prune(ctx, deps.DB,") {
			continue
		}
		if strings.Contains(string(raw), old) || !strings.Contains(string(raw), bound) {
			t.Errorf("the job workers still call %q without the job's context", old)
		}
	}
	raw, _ = os.ReadFile("api_files.go")
	if strings.Contains(string(raw), authUserLookupOld) {
		t.Error("the auth middleware's user lookup ignores the request's context")
	}
}

func TestFlagExposuresAreBatched(t *testing.T) {
	src := apiFlagsGo()
	if !strings.Contains(src, "func (e *Engine) writeExposures()") || !strings.Contains(src, "CreateInBatches(fresh, 256)") {
		t.Error("flag exposures are not written by one batching writer")
	}
	if strings.Contains(src, "err := e.db.WithContext(ctx).Create(&models.FlagExposure{") {
		t.Error("a flag check still inserts its exposure on its own")
	}
	mustFormatGo(t, "flags.go", src)
}

func TestAPIKeyVerificationIsCached(t *testing.T) {
	src := apiAPIKeyServiceGo()
	for _, want := range []string{"verifiedKeys.Load(cacheKey)", "lastTouched.Load(id)", "verifiedKeys.Range("} {
		if !strings.Contains(src, want) {
			t.Errorf("the API key service is missing %q", want)
		}
	}
	mustFormatGo(t, "services/api_key.go", src)
}

func TestStatsCountPerDayInSQL(t *testing.T) {
	for name, src := range map[string]string{"resource_stats_dispatch.go": resourceStatsDispatchGo(), "chart_dispatch.go": chartDispatchGo()} {
		if strings.Contains(src, `Select("created_at")`) || !strings.Contains(src, "dayExpr(db") {
			t.Errorf("%s still loads rows to count them in Go", name)
		}
		mustFormatGo(t, name, src)
	}
	mustFormatGo(t, "day_bucket.go", servicesDayBucketGo())

	old := "package services\n\n// buildDailySeries old doc\n// second line\nfunc buildDailySeries(db *gorm.DB, model interface{}) ([]ResourceStatsBucket, error) {\n" +
		"\tvar rows []row\n\terr := db.Model(model).Select(\"created_at\").Scan(&rows).Error\n\treturn nil, err\n}\n\nfunc other() {}\n"
	out, ok := replaceGoFunc(old, "func buildDailySeries(db *gorm.DB, model interface{}) ([]ResourceStatsBucket, error) {", dailySeriesFunc)
	if !ok || strings.Contains(out, "old doc") || !strings.Contains(out, "dayExpr(db") || !strings.Contains(out, "func other() {}") {
		t.Errorf("replaceGoFunc did not swap the function and its doc:\n%s", out)
	}
	mustFormatGo(t, "the replaced stats file", out)
}

func TestImageWorkIsBounded(t *testing.T) {
	src := mediaTransformGo()
	for _, want := range []string{"transformSlots <- struct{}{}", "defer func() { <-transformSlots }()"} {
		if !strings.Contains(src, want) {
			t.Errorf("the media transform is missing %q", want)
		}
	}
	mustFormatGo(t, "media/transform.go", src)
	raw, _ := os.ReadFile("api_storage_files.go")
	if !strings.Contains(string(raw), "renditionsWG.Wait()") {
		t.Error("renditions still upload one after another")
	}
}
