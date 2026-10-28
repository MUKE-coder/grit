package scaffold

// apiPaginateTestGo emits internal/paginate/paginate_test.go.
//
// Cursor pagination shipped once before without working. The machinery was
// there, but nothing could reach it: cursor mode was a compile-time Config
// field no handler set, so the code was never run and the fact that page two
// came back identical to page one went unnoticed. These tests exist so that
// cannot happen quietly again.
func apiPaginateTestGo() string {
	return `package paginate

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type widget struct {
	ID        string ` + "`" + `gorm:"primarykey;size:36"` + "`" + `
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

// A page-numbered list repeated rows and skipped others once the table was
// big enough.
//
// The page query sorted by the requested column and stopped there, so rows
// sharing a value had no defined order and the database was free to return
// them differently for each query. Postgres does exactly that as soon as it
// plans a parallel external sort: the order of tied rows depends on which
// worker finished first. Measured on 500,000 rows with repeating names, the
// same row came back on two pages while another never appeared, with no error
// anywhere (grit#112).
//
// Asserted on the SQL rather than on the rows, deliberately. SQLite returns
// ties in insertion order, so a test that walks the pages passes whether or
// not the tiebreaker is there, and would have gone on passing through the
// whole of this bug.
func TestPagesBreakTiesOnID(t *testing.T) {
	db := newDB(t)
	seed(t, db, 3)

	var orders []string
	if err := db.Callback().Query().After("gorm:query").Register("test:order", func(tx *gorm.DB) {
		sql := tx.Statement.SQL.String()
		if strings.Contains(sql, "LIMIT") && strings.Contains(sql, "ORDER BY") {
			orders = append(orders, sql[strings.Index(sql, "ORDER BY"):])
		}
	}); err != nil {
		t.Fatalf("register callback: %v", err)
	}

	cfg := Config{Sortable: map[string]bool{"rank": true}}
	if _, err := List[widget](db.Model(&widget{}), Params{
		Page: 2, PageSize: 1, SortBy: "rank", SortOrder: "asc",
	}, cfg); err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(orders) != 1 {
		t.Fatalf("expected one page query, saw %d: %q", len(orders), orders)
	}
	if !strings.HasPrefix(orders[0], "ORDER BY rank asc,id asc") {
		t.Errorf("page query ordered %q, want rank then id: without the second\n"+
			"column, tied rows have no defined order and pages repeat and skip", orders[0])
	}
}

// And sorting by id does not ask for it twice.
func TestSortingByIDDoesNotRepeatTheTiebreaker(t *testing.T) {
	db := newDB(t)
	seed(t, db, 3)

	var orders []string
	if err := db.Callback().Query().After("gorm:query").Register("test:order:id", func(tx *gorm.DB) {
		sql := tx.Statement.SQL.String()
		if strings.Contains(sql, "LIMIT") && strings.Contains(sql, "ORDER BY") {
			orders = append(orders, sql[strings.Index(sql, "ORDER BY"):])
		}
	}); err != nil {
		t.Fatalf("register callback: %v", err)
	}

	cfg := Config{Sortable: map[string]bool{"id": true}}
	if _, err := List[widget](db.Model(&widget{}), Params{
		Page: 1, PageSize: 1, SortBy: "id", SortOrder: "asc",
	}, cfg); err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(orders) == 1 && strings.Contains(orders[0], "id asc,id asc") {
		t.Errorf("ordered %q: id is already unique", orders[0])
	}
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

// ── The deferred join ────────────────────────────────────────────────────
//
// Past deepPageOffset, List fetches a page in two queries: the primary keys
// from the index, then the rows by key. It is an optimisation, so the only
// thing that matters is that nobody can tell. These drive List across the
// boundary on a table big enough to cross it, with ties on the sort column
// (where an unstable order shows up as a repeated or missing row) and with a
// preload (which the key query has to skip and the row query has to run).
//
// The preload test is not hypothetical: the first implementation shared the
// caller's Statement, so clearing preloads for the key query cleared them for
// the row query too and every association came back empty. This caught it.

type djAuthor struct {
	ID   uint   ` + "`" + `gorm:"primarykey"` + "`" + `
	Name string
}

type djPost struct {
	ID       uint   ` + "`" + `gorm:"primarykey"` + "`" + `
	Title    string ` + "`" + `gorm:"index"` + "`" + `
	AuthorID uint
	Author   djAuthor
}

const (
	djRows     = 3000 // comfortably past deepPageOffset
	djPageSize = 25
)

func newDeferredJoinDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db handle: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&djAuthor{}, &djPost{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	authors := []djAuthor{{Name: "Ada"}, {Name: "Grace"}, {Name: "Alan"}}
	if err := db.Create(&authors).Error; err != nil {
		t.Fatalf("authors: %v", err)
	}

	// Titles repeat every 10 rows on purpose: 300 groups of ties.
	posts := make([]djPost, 0, djRows)
	for i := 0; i < djRows; i++ {
		posts = append(posts, djPost{
			Title:    fmt.Sprintf("title-%04d", i/10),
			AuthorID: authors[i%len(authors)].ID,
		})
	}
	if err := db.CreateInBatches(&posts, 500).Error; err != nil {
		t.Fatalf("posts: %v", err)
	}
	return db
}

func djConfig() Config {
	return Config{
		Sortable:    map[string]bool{"title": true, "id": true},
		DefaultSort: "id",
	}
}

func djIDs(rows []djPost) []uint {
	out := make([]uint, len(rows))
	for i, r := range rows {
		out[i] = r.ID
	}
	return out
}

// TestDeepPageMatchesPlainOffset: the optimisation returns the same page.
func TestDeepPageMatchesPlainOffset(t *testing.T) {
	db := newDeferredJoinDB(t)

	for _, sort := range []struct{ by, order string }{
		{"title", "asc"},
		{"title", "desc"},
		{"id", "asc"},
	} {
		for _, page := range []int{1, 40, 41, 80, 120} {
			offset := (page - 1) * djPageSize

			got, err := List[djPost](db.Model(&djPost{}), Params{
				Page: page, PageSize: djPageSize, SortBy: sort.by, SortOrder: sort.order,
			}, djConfig())
			if err != nil {
				t.Fatalf("List page %d: %v", page, err)
			}

			var want []djPost
			if err := orderStably(db.Model(&djPost{}), sort.by, sort.order).
				Offset(offset).Limit(djPageSize).Find(&want).Error; err != nil {
				t.Fatalf("plain offset: %v", err)
			}

			gotIDs, wantIDs := djIDs(got.Data), djIDs(want)
			if len(gotIDs) != len(wantIDs) {
				t.Fatalf("sort=%s %s page %d: got %d rows, plain OFFSET gives %d",
					sort.by, sort.order, page, len(gotIDs), len(wantIDs))
			}
			for i := range gotIDs {
				if gotIDs[i] != wantIDs[i] {
					t.Fatalf("sort=%s %s page %d: row %d is id %d, plain OFFSET gives id %d",
						sort.by, sort.order, page, i, gotIDs[i], wantIDs[i])
				}
			}
		}
	}
}

// TestDeepPagesPartitionTheTable walks every page across the shallow/deep
// boundary: each row exactly once, none missing, none twice.
func TestDeepPagesPartitionTheTable(t *testing.T) {
	db := newDeferredJoinDB(t)

	seen := map[uint]int{}
	for page := 1; page <= djRows/djPageSize; page++ {
		res, err := List[djPost](db.Model(&djPost{}), Params{
			Page: page, PageSize: djPageSize, SortBy: "title", SortOrder: "asc",
		}, djConfig())
		if err != nil {
			t.Fatalf("page %d: %v", page, err)
		}
		for _, r := range res.Data {
			seen[r.ID]++
		}
	}

	if len(seen) != djRows {
		t.Errorf("walked every page and saw %d distinct rows, want %d", len(seen), djRows)
	}
	for id, n := range seen {
		if n != 1 {
			t.Errorf("row %d appeared %d times across the pages", id, n)
		}
	}
}

// TestPreloadsRunOnADeepPage: the key query drops preloads, the row query
// has to still run them.
func TestPreloadsRunOnADeepPage(t *testing.T) {
	db := newDeferredJoinDB(t)

	res, err := List[djPost](db.Model(&djPost{}).Preload("Author"), Params{
		Page: 80, PageSize: djPageSize, SortBy: "title", SortOrder: "asc",
	}, djConfig())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(res.Data) == 0 {
		t.Fatal("deep page is empty")
	}
	for _, r := range res.Data {
		if r.Author.ID == 0 || r.Author.Name == "" {
			t.Fatalf("row %d came back with no Author: the deferred join dropped "+
				"the preload instead of running it on the row query", r.ID)
		}
	}
}

// TestDeepPageFiltersApplyToBothQueries: the caller's WHERE has to reach
// both, or a deep page returns rows the filter excluded.
func TestDeepPageFiltersApplyToBothQueries(t *testing.T) {
	db := newDeferredJoinDB(t)

	res, err := List[djPost](db.Model(&djPost{}).Where("id % 2 = 0"), Params{
		Page: 50, PageSize: djPageSize, SortBy: "id", SortOrder: "asc",
	}, djConfig())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(res.Data) == 0 {
		t.Fatal("deep filtered page is empty")
	}
	for _, r := range res.Data {
		if r.ID%2 != 0 {
			t.Fatalf("row %d is not in the filtered set: the deferred join lost "+
				"the caller's WHERE on one of its two queries", r.ID)
		}
	}
}

// TestDeepPagePastTheEndIsEmpty: no error, and no WHERE id IN ().
func TestDeepPagePastTheEndIsEmpty(t *testing.T) {
	db := newDeferredJoinDB(t)

	res, err := List[djPost](db.Model(&djPost{}), Params{
		Page: 9999, PageSize: djPageSize, SortBy: "title", SortOrder: "asc",
	}, djConfig())
	if err != nil {
		t.Fatalf("a page past the end errored: %v", err)
	}
	if len(res.Data) != 0 {
		t.Errorf("a page past the end returned %d rows", len(res.Data))
	}
	if res.Meta.Total != djRows {
		t.Errorf("total is %d, want %d", res.Meta.Total, djRows)
	}
}
`
}
