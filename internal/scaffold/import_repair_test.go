package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

const oldImportService = `package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"gorm.io/gorm/clause"

	"demo/internal/models"
)

func (s *ContactService) ImportCSV(ctx context.Context, jobID, path string) {
	db := s.db(ctx)
	record := func(fields map[string]interface{}) {}
	defer os.Remove(path)

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	reader := csv.NewReader(f)
	get := func(rec []string, key string) (string, bool) { return "", false }

	created, skipped, failed := 0, 0, 0
	rowErrors := []map[string]interface{}{}
	// The progress map has an "errors" key, which once passed for the import.
	checkpoint := map[string]interface{}{
		"errors": rowErrors,
	}
	_ = checkpoint
	rowNum := 1
	for {
		rec, err := reader.Read()
		if err != nil {
			break
		}
		item := models.Category{}
		if v, ok := get(rec, "group"); ok && v != "" {
			var rel models.Group
			if err := db.Where("name = ?", v).First(&rel).Error; err != nil {
				rel = models.Group{Name: v}
				db.Create(&rel)
			}
			item.GroupID = rel.ID
		}
		if v, ok := get(rec, "parent"); ok && v != "" {
			var rel models.Category
			if err := db.Where("title = ?", v).First(&rel).Error; err != nil {
				rel = models.Category{Title: v}
				db.Create(&rel)
			}
			item.ParentID = &rel.ID
		}
		_, _, _, _, _, _ = created, skipped, failed, rowErrors, rowNum, clause.OnConflict{}
		fmt.Println(item)
	}
	record(nil)
}
`

func TestRepairImportService(t *testing.T) {
	out, fixed, warn := repairImportServiceSource(oldImportService, "demo")
	if len(warn) > 0 || len(fixed) != 2 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	for _, want := range []string{
		ImportNameResolver("group", "Group", "name", "Name"),
		ImportNameResolver("parent", "Category", "title", "Title"),
		ImportNameAssign("group", "GroupID", false),
		ImportNameAssign("parent", "ParentID", true),
		ImportLimitGate,
		`"errors"`, `"gorm.io/gorm"`, `"demo/internal/imports"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing:\n%s", want)
		}
	}
	if strings.Contains(out, "\t\t\t\tdb.Create(&rel)\n") {
		t.Error("a lookup still creates without checking")
	}
	if strings.Index(out, "resolveGroup := func") > strings.Index(out, "created, skipped, failed := 0, 0, 0") {
		t.Error("the resolvers are declared after the rows start")
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("not valid Go: %v\n%s", err, out)
	}
	if again, fixed, _ := repairImportServiceSource(out, "demo"); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed the importer again")
	}
}

// An importer from before v3.251.0 has no record helper; it still takes turns.
func TestRepairImportServiceWithoutRecordHelper(t *testing.T) {
	src := strings.Replace(oldImportService, "\trecord := func(fields map[string]interface{}) {}\n", "", 1)
	src = strings.Replace(src, "\trecord(nil)\n", "", 1)
	out, fixed, warn := repairImportServiceSource(src, "demo")
	if len(warn) > 0 || len(fixed) != 2 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if !strings.Contains(out, importLimitGateDirect) || !strings.Contains(out, "\t\"log\"\n") {
		t.Errorf("the direct gate or its log import is missing:\n%s", out)
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("not valid Go: %v\n%s", err, out)
	}
	if again, fixed, _ := repairImportServiceSource(out, "demo"); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed the importer again")
	}
}

func TestRepairImportServiceLeavesAnEditedLookup(t *testing.T) {
	src := strings.Replace(oldImportService, "\t\t\t\tdb.Create(&rel)\n\t\t\t}\n\t\t\titem.GroupID", "\t\t\t\tdb.Save(&rel)\n\t\t\t}\n\t\t\titem.GroupID", 1)
	out, fixed, _ := repairImportServiceSource(src, "demo")
	if strings.Contains(out, "resolveGroup") {
		t.Error("an edited lookup was rewritten")
	}
	if !strings.Contains(out, "resolveParent") || len(fixed) != 2 {
		t.Errorf("the untouched lookup and the gate were not fixed: %v", fixed)
	}
}

func TestRepairWAFImportExclusion(t *testing.T) {
	fresh := apiRoutesGo()
	if !strings.Contains(fresh, `"/*/import"`) {
		t.Fatal("the routes template does not exclude imports from the WAF")
	}
	old := strings.Replace(fresh, wafImportExclusion, "", 1)
	out, fixed, warn := repairWAFImportExclusionSource(old)
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if out != fresh {
		t.Error("the repaired routes.go differs from a fresh one")
	}
	if again, fixed, _ := repairWAFImportExclusionSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed routes.go again")
	}
	if out, fixed, _ := repairWAFImportExclusionSource(fresh); out != fresh || len(fixed) > 0 {
		t.Error("a fresh routes.go still needs the repair")
	}
}

func TestImportCodegenNames(t *testing.T) {
	if got := ImportNameAssign("parent_category", "ParentCategoryID", false); !strings.Contains(got, "resolveParentCategory(v)") {
		t.Errorf("snake case column: %s", got)
	}
	if got := ImportNameResolver("parent_category", "Category", "name", "Name"); !strings.Contains(got, "parentCategoryIDs := map[string]string{}") {
		t.Errorf("snake case resolver: %s", got)
	}
	for name, src := range map[string]string{"limit.go": importsLimitGo(), "limit_test.go": importsLimitTestGo()} {
		if _, err := format.Source([]byte(src)); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}
