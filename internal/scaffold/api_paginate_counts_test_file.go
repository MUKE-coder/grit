package scaffold

// apiPaginateCountsTestGo is the generated project's test for ?counts=, the
// extra totals the admin's stat cards read from the list response.
func apiPaginateCountsTestGo() string {
	return `package paginate

import (
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type countedNote struct {
	ID        uint
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func newCountsDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	// One connection: each connection to ":memory:" is its own empty database.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db handle: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&countedNote{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	day := 24 * time.Hour
	now := time.Now()
	for _, n := range []countedNote{
		{Title: "fresh", CreatedAt: now.Add(-2 * day), UpdatedAt: now.Add(-2 * day)},
		{Title: "other", CreatedAt: now.Add(-3 * day), UpdatedAt: now.Add(-3 * day)},
		{Title: "edited", CreatedAt: now.Add(-10 * day), UpdatedAt: now.Add(-1 * day)},
		{Title: "old", CreatedAt: now.Add(-40 * day), UpdatedAt: now.Add(-40 * day)},
	} {
		n := n
		if err := db.Create(&n).Error; err != nil {
			t.Fatalf("seed %s: %v", n.Title, err)
		}
	}
	return db
}

func TestListCountsAnswerTheWindowsAskedFor(t *testing.T) {
	db := newCountsDB(t)
	res, err := List[countedNote](db.Model(&countedNote{}), Params{
		Page: 1, PageSize: 1,
		Counts: []string{"created_7d", "created_30d", "updated_7d"},
	}, Config{})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if res.Meta.Total != 4 || len(res.Data) != 1 {
		t.Fatalf("total %d with %d rows, want 4 with 1", res.Meta.Total, len(res.Data))
	}
	want := map[string]int64{"created_7d": 2, "created_30d": 3, "updated_7d": 3}
	if !reflect.DeepEqual(res.Meta.Counts, want) {
		t.Errorf("counts %v, want %v", res.Meta.Counts, want)
	}
}

func TestListCountsFollowTheSearch(t *testing.T) {
	db := newCountsDB(t)
	res, err := List[countedNote](db.Model(&countedNote{}), Params{
		Page: 1, PageSize: 20, Search: "other",
		Counts: []string{"created_7d", "created_30d"},
	}, Config{Searchable: []string{"title"}})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	// A card reading 3 above a table of one match is the mismatch the counts
	// exist to avoid.
	want := map[string]int64{"created_7d": 1, "created_30d": 1}
	if res.Meta.Total != 1 || !reflect.DeepEqual(res.Meta.Counts, want) {
		t.Errorf("total %d counts %v, want 1 and %v", res.Meta.Total, res.Meta.Counts, want)
	}
}

func TestListWithoutCountsLeavesThemOut(t *testing.T) {
	db := newCountsDB(t)
	res, err := List[countedNote](db.Model(&countedNote{}), Params{Page: 1, PageSize: 20}, Config{})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if res.Meta.Counts != nil {
		t.Errorf("counts %v on a request that asked for none", res.Meta.Counts)
	}
}

func TestParseCountsKeepsOnlyWellFormedNames(t *testing.T) {
	got := parseCounts("created_7d, updated_30d,deleted_7d,created_0d,created_xd,created_7d,updated_99999d,created_7")
	want := []string{"created_7d", "updated_30d"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseCounts = %v, want %v", got, want)
	}
}

func TestCountsForAHandWrittenList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newCountsDB(t)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/api/v1/notes?counts=created_7d,updated_7d", nil)
	counts, err := Counts(c, db.Model(&countedNote{}))
	if err != nil {
		t.Fatalf("counts: %v", err)
	}
	want := map[string]int64{"created_7d": 2, "updated_7d": 3}
	if !reflect.DeepEqual(counts, want) {
		t.Errorf("counts %v, want %v", counts, want)
	}

	// A new context: gin caches the query string it parsed for the first request.
	c, _ = gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/api/v1/notes", nil)
	if counts, err := Counts(c, db.Model(&countedNote{})); err != nil || counts != nil {
		t.Errorf("no ?counts= gave %v, %v; want nil, nil", counts, err)
	}
}
`
}
