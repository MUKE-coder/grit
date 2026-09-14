package scaffold

import (
	"go/format"
	"os"
	"strings"
	"testing"
)

const oldGeneratedModel = "package models\n\nimport (\n\t\"time\"\n\n\t\"gorm.io/gorm\"\n)\n\n" +
	"type Contact struct {\n" +
	"\tID        string         `gorm:\"primarykey;size:36\" json:\"id\"`\n" +
	"\tName      string         `gorm:\"size:255\" json:\"name\" binding:\"required\"`\n" +
	"\tEmail     string         `gorm:\"size:255\" json:\"email\" binding:\"required\"`\n" +
	"\tPhone     string         `gorm:\"size:255\" json:\"phone\"`\n" +
	"\tVersion   int            `gorm:\"not null;default:1\" json:\"version\"`\n" +
	"\tCreatedAt time.Time      `json:\"created_at\"`\n" +
	"\tUpdatedAt time.Time      `json:\"updated_at\"`\n" +
	"\tDeletedAt gorm.DeletedAt `gorm:\"index\" json:\"-\"`\n" +
	generatedArchivedMarker +
	"}\n"

func TestRepairModelIndexes(t *testing.T) {
	// As the template writes it, and as gofmt leaves it in a generated project.
	for _, src := range []string{oldGeneratedModel, strings.Replace(oldGeneratedModel, "*time.Time  `", "*time.Time `", 1)} {
		repairModelIndexesCase(t, src)
	}
}

func repairModelIndexesCase(t *testing.T, src string) {
	t.Helper()
	searchable := map[string][]string{"Contact": {"name", "email"}}
	out, fixed, warn := repairModelIndexesSource(src, "contact.go", searchable)
	if len(warn) > 0 || len(fixed) != 2 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	for _, want := range []string{
		"`gorm:\"index\" json:\"created_at\"`",
		"`gorm:\"size:255\" json:\"name\" binding:\"required\" search:\"trigram\"`",
		"`gorm:\"size:255\" json:\"email\" binding:\"required\" search:\"trigram\"`",
		"`gorm:\"size:255\" json:\"phone\"`\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in:\n%s", want, out)
		}
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("not valid Go: %v", err)
	}
	if again, fixed, _ := repairModelIndexesSource(out, "contact.go", searchable); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed the model again")
	}
}

func TestRepairModelIndexesLeavesOtherModels(t *testing.T) {
	src := "package models\n\ntype Note struct {\n\tCreatedAt time.Time `json:\"created_at\"`\n}\n"
	if out, fixed, _ := repairModelIndexesSource(src, "note.go", nil); out != src || len(fixed) > 0 {
		t.Errorf("a hand-written model was changed: %v", fixed)
	}
}

func TestRepairFrameworkModelsIndexCreatedAt(t *testing.T) {
	for base, fresh := range map[string]string{"upload.go": apiUploadModelGo(), "user.go": apiUserModelGo()} {
		old := strings.Replace(fresh, "`gorm:\"index\" json:\"created_at\"`", "`json:\"created_at\"`", 1)
		if old == fresh {
			t.Fatalf("%s: could not rebuild the model from before the repair", base)
		}
		out, fixed, _ := repairModelIndexesSource(old, base, nil)
		if out != fresh || len(fixed) != 1 {
			t.Errorf("%s: fixed %v; the repaired model differs from a fresh one", base, fixed)
		}
	}
}

func TestSearchableColumnsReadsServiceConfigs(t *testing.T) {
	dir := t.TempDir()
	service := "package services\n\nvar contactListConfig = paginate.Config{\n\tSearchable: []string{\"name\", \"email\"},\n\tSortable:   map[string]bool{\"id\": true},\n}\n"
	if err := os.WriteFile(dir+"/contact.go", []byte(service), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := searchableColumns(dir)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got["Contact"], ",") != "name,email" {
		t.Errorf("searchable %v", got)
	}
}

func TestRepairPaginateWiring(t *testing.T) {
	fresh := strings.ReplaceAll(apiDatabaseGo(), "{{MODULE}}", "demo")
	old := strings.Replace(fresh, paginateConnectHook, "", 1)
	old = strings.Replace(old, "\t\"demo/internal/paginate\"\n", "", 1)
	if old == fresh {
		t.Fatal("could not rebuild database.go from before the repair")
	}
	out, fixed, warn := repairPaginateWiringSource(old, "demo")
	if len(warn) > 0 || len(fixed) != 1 || !strings.Contains(out, "paginate.Install(db)") || !strings.Contains(out, "\"demo/internal/paginate\"") {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("not valid Go: %v", err)
	}
	if again, fixed, _ := repairPaginateWiringSource(out, "demo"); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed database.go again")
	}
}

func TestFreshTemplatesNeedNoListPerformanceRepair(t *testing.T) {
	db := strings.ReplaceAll(apiDatabaseGo(), "{{MODULE}}", "demo")
	if out, fixed, _ := repairPaginateWiringSource(db, "demo"); out != db || len(fixed) > 0 {
		t.Error("database.go still needs paginate.Install")
	}
	for base, src := range map[string]string{"upload.go": apiUploadModelGo(), "user.go": apiUserModelGo()} {
		if out, fixed, _ := repairModelIndexesSource(src, base, nil); out != src || len(fixed) > 0 {
			t.Errorf("%s still needs the created_at index", base)
		}
	}
	if !strings.Contains(apiMigrateMainGo(), "paginate.EnsureSearchIndexes(db, models.Models()...)") {
		t.Error("cmd/migrate does not build the search indexes")
	}
	if !strings.Contains(apiPaginateGo(), "countTotal(query, &result.Meta.Total)") {
		t.Error("List still counts every page")
	}
	for name, src := range map[string]string{"count.go": paginateCountGo(), "search_index.go": paginateSearchIndexGo(), "count_test.go": paginateCountTestGo()} {
		if _, err := format.Source([]byte(src)); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}
