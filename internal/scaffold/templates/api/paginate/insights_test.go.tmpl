package paginate

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type insightRow struct {
	ID        uint
	Title     string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func newInsightsDB(t *testing.T, rows []insightRow) *gorm.DB {
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
	if err := db.AutoMigrate(&insightRow{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	for i := range rows {
		if err := db.Create(&rows[i]).Error; err != nil {
			t.Fatalf("seeding: %v", err)
		}
	}
	return db
}

func at(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 12, 0, 0, 0, time.Local)
}

func TestParseSeriesDefaultsToTwelveMonths(t *testing.T) {
	spec, ok := parseSeries("created_at")
	if !ok {
		t.Fatal("a bare column was refused")
	}
	if spec.unit != "month" || spec.buckets != 12 {
		t.Errorf("got %+v", spec)
	}
}

func TestParseSeriesReadsUnitAndCount(t *testing.T) {
	spec, ok := parseSeries("updated_at:day:30")
	if !ok {
		t.Fatal("refused")
	}
	if spec.column != "updated_at" || spec.unit != "day" || spec.buckets != 30 {
		t.Errorf("got %+v", spec)
	}
}

// The column reaches the SQL, so anything but the two literals is refused
// rather than quoted and hoped for.
func TestParseSeriesRefusesAnyOtherColumn(t *testing.T) {
	for _, raw := range []string{
		"",
		"name",
		"created_at); DROP TABLE users;--",
		"created_at:century",
		"created_at:month:0",
		"created_at:month:zero",
	} {
		if _, ok := parseSeries(raw); ok {
			t.Errorf("accepted %q", raw)
		}
	}
}

func TestParseSeriesCapsTheBucketCount(t *testing.T) {
	spec, _ := parseSeries("created_at:day:5000")
	if spec.buckets != maxBuckets {
		t.Errorf("buckets %d, want the cap %d", spec.buckets, maxBuckets)
	}
}

func TestParseBreakdownKeepsOnlyFilterableColumns(t *testing.T) {
	allowed := map[string]bool{"status": true, "role": true}
	got := parseBreakdown("status, role ,secret,status", allowed)
	if len(got) != 2 || got[0] != "status" || got[1] != "role" {
		t.Errorf("got %v", got)
	}
}

func TestParseBreakdownCapsTheColumnCount(t *testing.T) {
	allowed := map[string]bool{"a": true, "b": true, "c": true, "d": true}
	if got := parseBreakdown("a,b,c,d", allowed); len(got) != maxBreakdownColumns {
		t.Errorf("got %d columns, want the cap %d", len(got), maxBreakdownColumns)
	}
}

func TestSeriesBucketsCountsByMonthOldestFirst(t *testing.T) {
	db := newInsightsDB(t, []insightRow{
		{Title: "a", CreatedAt: at(2026, time.January, 3)},
		{Title: "b", CreatedAt: at(2026, time.January, 20)},
		{Title: "c", CreatedAt: at(2026, time.March, 1)},
	})

	spec, _ := parseSeries("created_at:month:12")
	got, err := seriesBuckets(db.Model(&insightRow{}), spec)
	if err != nil {
		t.Fatalf("seriesBuckets: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %+v, want a bucket for January and one for March", got)
	}
	// A chart reads left to right, so the oldest comes first.
	if got[0].Bucket != "2026-01" || got[0].Count != 2 {
		t.Errorf("first bucket %+v", got[0])
	}
	if got[1].Bucket != "2026-03" || got[1].Count != 1 {
		t.Errorf("second bucket %+v", got[1])
	}
}

// A month with nothing in it is absent rather than zero: the window belongs to
// whoever is drawing the chart, and this cannot invent one.
func TestSeriesBucketsLeavesEmptyPeriodsOut(t *testing.T) {
	db := newInsightsDB(t, []insightRow{{Title: "only", CreatedAt: at(2026, time.June, 9)}})
	spec, _ := parseSeries("created_at:month:12")
	got, _ := seriesBuckets(db.Model(&insightRow{}), spec)
	if len(got) != 1 || got[0].Bucket != "2026-06" {
		t.Errorf("got %+v", got)
	}
}

// Every database labels a week by the date of its Monday, so the three agree.
// 2026-06-09 is a Tuesday; its Monday is the 8th.
func TestSeriesBucketsLabelsAWeekByItsMonday(t *testing.T) {
	db := newInsightsDB(t, []insightRow{
		{Title: "tue", CreatedAt: at(2026, time.June, 9)},
		{Title: "sun", CreatedAt: at(2026, time.June, 14)},
	})
	spec, _ := parseSeries("created_at:week:8")
	got, err := seriesBuckets(db.Model(&insightRow{}), spec)
	if err != nil {
		t.Fatalf("seriesBuckets: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %+v, want both days in one week", got)
	}
	if got[0].Bucket != "2026-06-08" || got[0].Count != 2 {
		t.Errorf("got %+v, want the Monday 2026-06-08 with 2", got[0])
	}
}

// The point of running over the list's query: the chart describes the rows the
// table is showing.
func TestSeriesBucketsRespectsTheQueryItIsGiven(t *testing.T) {
	db := newInsightsDB(t, []insightRow{
		{Title: "keep", Status: "live", CreatedAt: at(2026, time.January, 3)},
		{Title: "drop", Status: "draft", CreatedAt: at(2026, time.January, 4)},
	})
	spec, _ := parseSeries("created_at:month:12")
	got, _ := seriesBuckets(db.Model(&insightRow{}).Where("status = ?", "live"), spec)
	if len(got) != 1 || got[0].Count != 1 {
		t.Errorf("got %+v, want only the live row counted", got)
	}
}

func TestBreakdownSlicesCountsValuesLargestFirst(t *testing.T) {
	db := newInsightsDB(t, []insightRow{
		{Title: "a", Status: "live"},
		{Title: "b", Status: "live"},
		{Title: "c", Status: "draft"},
	})
	got, err := breakdownSlices(db.Model(&insightRow{}), []string{"status"})
	if err != nil {
		t.Fatalf("breakdownSlices: %v", err)
	}
	slices := got["status"]
	if len(slices) != 2 {
		t.Fatalf("got %+v", slices)
	}
	if slices[0].Value != "live" || slices[0].Count != 2 {
		t.Errorf("first slice %+v", slices[0])
	}
	if slices[1].Value != "draft" || slices[1].Count != 1 {
		t.Errorf("second slice %+v", slices[1])
	}
}

// Nothing filled in is one answer, not two.
func TestBreakdownSlicesReportsEmptyAsOneValue(t *testing.T) {
	db := newInsightsDB(t, []insightRow{
		{Title: "a", Status: ""},
		{Title: "b", Status: ""},
	})
	got, _ := breakdownSlices(db.Model(&insightRow{}), []string{"status"})
	if len(got["status"]) != 1 || got["status"][0].Value != "" || got["status"][0].Count != 2 {
		t.Errorf("got %+v", got["status"])
	}
}

func TestBucketExprCoversEveryDatabaseGritSupports(t *testing.T) {
	db := newInsightsDB(t, nil)
	for _, unit := range []string{"day", "week", "month"} {
		if expr := bucketExpr(db, "created_at", unit); expr == "" {
			t.Errorf("no expression for %s", unit)
		}
	}
	// An unknown unit cannot reach here through parseSeries, and if it ever
	// does it must still be SQL rather than an empty select list.
	if expr := bucketExpr(db, "created_at", "fortnight"); expr == "" {
		t.Error("an unknown unit produced no expression at all")
	}
}
