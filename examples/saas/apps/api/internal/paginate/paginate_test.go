package paginate

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type widget struct {
	ID        string `gorm:"primarykey;size:36"`
	Name      string
	Rank      int
	CreatedAt time.Time
}

func newDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	// One connection. Every connection to ":memory:" gets its own empty
	// database, so a second one would silently see no rows at all.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db handle: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&widget{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// seed inserts n rows with distinct, increasing created_at values.
//
// Distinct on purpose: rows written in the same instant are the case the id
// tiebreaker exists for, and TestCursorWalksTiedTimestamps covers that
// separately.
func seed(t *testing.T, db *gorm.DB, n int) {
	t.Helper()
	base := time.Now().Add(-time.Duration(n) * time.Minute)
	for i := 0; i < n; i++ {
		w := widget{
			ID:        pad(i),
			Name:      "w" + pad(i),
			Rank:      (i * 37) % 100,
			CreatedAt: base.Add(time.Duration(i) * time.Minute),
		}
		if err := db.Create(&w).Error; err != nil {
			t.Fatalf("seed %d: %v", i, err)
		}
	}
}

func pad(i int) string {
	s := "000" + itoa(i)
	return s[len(s)-4:]
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	out := ""
	for i > 0 {
		out = string(rune('0'+i%10)) + out
		i /= 10
	}
	return out
}

// walk pages all the way through with a cursor and returns what it saw.
func walk(t *testing.T, db *gorm.DB, cfg Config, p Params) []string {
	t.Helper()
	var seen []string
	cursor := ""
	for pages := 0; ; pages++ {
		if pages > 50 {
			t.Fatal("cursor never reported has_more=false; it is not advancing")
		}
		q := p
		q.Cursor = cursor
		q.CursorMode = true
		res, err := List[widget](db.Model(&widget{}), q, cfg)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		for _, w := range res.Data {
			seen = append(seen, w.Name)
		}
		if !res.Meta.HasMore {
			return seen
		}
		if res.Meta.NextCursor == "" {
			t.Fatal("has_more is true with no next_cursor: the walk cannot continue")
		}
		if res.Meta.NextCursor == cursor {
			t.Fatal("next_cursor did not change between pages")
		}
		cursor = res.Meta.NextCursor
	}
}

func assertUnique(t *testing.T, seen []string, want int) {
	t.Helper()
	if len(seen) != want {
		t.Errorf("saw %d rows, want %d", len(seen), want)
	}
	uniq := map[string]bool{}
	for _, s := range seen {
		if uniq[s] {
			t.Errorf("row %q came back on more than one page", s)
		}
		uniq[s] = true
	}
	if len(uniq) != want {
		t.Errorf("saw %d distinct rows, want %d", len(uniq), want)
	}
}

// The default sort is a timestamp, which is where this broke: the cursor
// encodes an RFC3339 string, and comparing a datetime column against text
// matches every row in SQLite, so page two was page one forever.
func TestCursorWalksEveryRowOnce(t *testing.T) {
	db := newDB(t)
	seed(t, db, 47)

	cfg := Config{Sortable: map[string]bool{"created_at": true, "rank": true}}
	seen := walk(t, db, cfg, Params{PageSize: 10, SortBy: "created_at", SortOrder: "desc"})
	assertUnique(t, seen, 47)
}

// A non-unique sort column: without the id tiebreaker in the WHERE and the
// ORDER BY, rows sharing a rank get skipped or repeated at page boundaries.
func TestCursorHandlesDuplicateSortValues(t *testing.T) {
	db := newDB(t)
	base := time.Now()
	for i := 0; i < 40; i++ {
		w := widget{ID: pad(i), Name: "w" + pad(i), Rank: i / 10, CreatedAt: base}
		if err := db.Create(&w).Error; err != nil {
			t.Fatalf("seed: %v", err)
		}
	}

	cfg := Config{Sortable: map[string]bool{"rank": true}}
	seen := walk(t, db, cfg, Params{PageSize: 7, SortBy: "rank", SortOrder: "asc"})
	assertUnique(t, seen, 40)
}

// Every row written in the same instant, which is what a seeder or an import
// produces. The timestamp alone cannot order these at all.
func TestCursorWalksTiedTimestamps(t *testing.T) {
	db := newDB(t)
	at := time.Now()
	for i := 0; i < 25; i++ {
		w := widget{ID: pad(i), Name: "w" + pad(i), Rank: i, CreatedAt: at}
		if err := db.Create(&w).Error; err != nil {
			t.Fatalf("seed: %v", err)
		}
	}

	cfg := Config{Sortable: map[string]bool{"created_at": true}}
	seen := walk(t, db, cfg, Params{PageSize: 6, SortBy: "created_at", SortOrder: "desc"})
	assertUnique(t, seen, 25)
}

// The reason to prefer a cursor: rows arriving mid-walk do not shift the
// window. With an offset, inserting five rows at the top makes page two repeat
// five rows page one already showed.
func TestCursorIsStableWhileRowsAreInserted(t *testing.T) {
	db := newDB(t)
	seed(t, db, 30)

	cfg := Config{Sortable: map[string]bool{"created_at": true}}
	p := Params{PageSize: 10, SortBy: "created_at", SortOrder: "desc", CursorMode: true}

	first, err := List[widget](db.Model(&widget{}), p, cfg)
	if err != nil {
		t.Fatalf("first page: %v", err)
	}

	later := time.Now().Add(time.Hour)
	for i := 0; i < 5; i++ {
		w := widget{ID: "new" + itoa(i), Name: "new" + itoa(i), CreatedAt: later}
		if err := db.Create(&w).Error; err != nil {
			t.Fatalf("insert: %v", err)
		}
	}

	p.Cursor = first.Meta.NextCursor
	second, err := List[widget](db.Model(&widget{}), p, cfg)
	if err != nil {
		t.Fatalf("second page: %v", err)
	}

	onFirst := map[string]bool{}
	for _, w := range first.Data {
		onFirst[w.Name] = true
	}
	for _, w := range second.Data {
		if onFirst[w.Name] {
			t.Errorf("row %q appeared on both pages", w.Name)
		}
		if len(w.Name) > 3 && w.Name[:3] == "new" {
			t.Errorf("row %q was inserted after the cursor was issued and still appeared", w.Name)
		}
	}
}

// Offset stays the default and keeps answering "how many pages", which the
// admin table needs and a cursor cannot give.
func TestOffsetRemainsTheDefault(t *testing.T) {
	db := newDB(t)
	seed(t, db, 25)

	cfg := Config{Sortable: map[string]bool{"created_at": true}}
	res, err := List[widget](db.Model(&widget{}), Params{Page: 1, PageSize: 10}, cfg)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if res.Meta.Total != 25 {
		t.Errorf("total = %d, want 25", res.Meta.Total)
	}
	if res.Meta.Pages != 3 {
		t.Errorf("pages = %d, want 3", res.Meta.Pages)
	}
	if res.Meta.Mode != "" {
		t.Errorf("mode = %q on an offset response, want empty", res.Meta.Mode)
	}
	if res.Meta.NextCursor != "" {
		t.Error("an offset response handed back a cursor")
	}
}

// A cursor response says so, because its total and pages are both zero and
// that is otherwise indistinguishable from an empty table.
func TestCursorResponseIsLabelled(t *testing.T) {
	db := newDB(t)
	seed(t, db, 5)

	cfg := Config{Sortable: map[string]bool{"created_at": true}}
	res, err := List[widget](db.Model(&widget{}), Params{PageSize: 2, CursorMode: true}, cfg)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if res.Meta.Mode != "cursor" {
		t.Errorf("mode = %q, want cursor", res.Meta.Mode)
	}
}

func TestTypedCursorValue(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)
	if got := typedCursorValue(now.Format(time.RFC3339Nano)); !got.(time.Time).Equal(now) {
		t.Errorf("a timestamp did not survive the round trip: %v", got)
	}
	if got := typedCursorValue("1999"); got != int64(1999) {
		t.Errorf("integer: got %v (%T), want int64", got, got)
	}
	if got := typedCursorValue("19.99"); got != 19.99 {
		t.Errorf("float: got %v (%T), want float64", got, got)
	}
	if got := typedCursorValue("draft"); got != "draft" {
		t.Errorf("string: got %v, want draft", got)
	}
}
