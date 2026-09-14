package scaffold

// paginateCountGo emits internal/paginate/count.go: list totals reused across
// the pages of a list until a write to the table (H18 in the contact-app
// review).
func paginateCountGo() string {
	return `package paginate

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
)

// CountTTL is the longest a list total is reused. A write to the table through
// the database handle ends it at once; the TTL covers the writes that handle
// cannot see, from raw SQL on another connection or another instance of the API.
var CountTTL = 15 * time.Second

// countCacheLimit bounds how many totals are kept. Past it they are all dropped
// and counted again, which costs one count each.
const countCacheLimit = 5000

type countEntry struct {
	total   int64
	gen     int64
	counted time.Time
}

// countScope tracks writes through one database handle.
type countScope struct {
	tables sync.Map     // table name -> *atomic.Int64
	raw    atomic.Int64 // raw SQL names no table, so it moves every total on
}

func (s *countScope) table(name string) *atomic.Int64 {
	if v, ok := s.tables.Load(name); ok {
		return v.(*atomic.Int64)
	}
	v, _ := s.tables.LoadOrStore(name, new(atomic.Int64))
	return v.(*atomic.Int64)
}

var (
	// Keyed by the handle's callbacks, which every session and transaction made
	// from it shares. Not by *gorm.Config: WithContext copies that, so a cache
	// keyed on it missed on every request.
	countScopes sync.Map // callbacks -> *countScope
	countMu     sync.Mutex
	countTotals = map[string]countEntry{}
)

// Install lets List reuse a list's total across its pages, until a write to the
// table through db, or CountTTL.
//
// Counting the whole match on every page was most of what a list page cost: on
// a million Postgres rows the count took 75 ms and the indexed page under 1 ms.
// Without Install, every page counts, as before.
func Install(db *gorm.DB) error {
	if _, ok := countScopes.Load(db.Callback()); ok {
		return nil
	}
	scope := &countScope{}
	written := func(tx *gorm.DB) {
		if tx.Statement == nil || tx.Statement.Table == "" {
			scope.raw.Add(1)
			return
		}
		scope.table(tx.Statement.Table).Add(1)
	}
	callbacks := db.Callback()
	if err := callbacks.Create().After("gorm:create").Register("paginate:count_create", written); err != nil {
		return err
	}
	if err := callbacks.Update().After("gorm:update").Register("paginate:count_update", written); err != nil {
		return err
	}
	if err := callbacks.Delete().After("gorm:delete").Register("paginate:count_delete", written); err != nil {
		return err
	}
	if err := callbacks.Raw().After("gorm:raw").Register("paginate:count_raw", func(*gorm.DB) { scope.raw.Add(1) }); err != nil {
		return err
	}
	countScopes.Store(db.Callback(), scope)
	return nil
}

// countTotal counts what query matches, or reuses the total from an earlier
// page when nothing has been written to the table since.
func countTotal(query *gorm.DB, total *int64) error {
	v, ok := countScopes.Load(query.Callback())
	if !ok {
		return query.Count(total).Error
	}
	scope := v.(*countScope)

	// The SQL the count would run, with its arguments, is the key: the same
	// filters, search and tenant scope give the same key, and nothing else does.
	var dryTotal int64
	dry := query.Session(&gorm.Session{DryRun: true}).Count(&dryTotal)
	if dry.Error != nil || dry.Statement.Table == "" {
		return query.Count(total).Error
	}
	key := fmt.Sprintf("%p|%s|%v", query.Callback(), dry.Statement.SQL.String(), dry.Statement.Vars)
	// Read before counting: a write that lands during the count moves the
	// generation past the one stored with it, so the total is not reused.
	gen := scope.raw.Load() + scope.table(dry.Statement.Table).Load()

	now := time.Now()
	countMu.Lock()
	entry, hit := countTotals[key]
	countMu.Unlock()
	if hit && entry.gen == gen && now.Sub(entry.counted) < CountTTL {
		*total = entry.total
		return nil
	}

	if err := query.Count(total).Error; err != nil {
		return err
	}
	countMu.Lock()
	if len(countTotals) >= countCacheLimit {
		countTotals = map[string]countEntry{}
	}
	countTotals[key] = countEntry{total: *total, gen: gen, counted: now}
	countMu.Unlock()
	return nil
}
`
}

// paginateSearchIndexGo emits internal/paginate/search_index.go.
func paginateSearchIndexGo() string {
	return `package paginate

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// EnsureSearchIndexes gives every column tagged search:"trigram" a trigram index
// on Postgres, so List's search is served by an index instead of reading the
// whole table.
//
// LOWER(col) LIKE '%term%' cannot use a B-tree index. A pg_trgm GIN index on
// lower(col) can: on a million contacts a search went from 270 ms to 2 ms. Other
// databases are left as they are. The pg_trgm extension is one managed Postgres
// offers to ordinary users; when it cannot be created, the indexes are skipped
// with a message and search keeps working, unindexed.
func EnsureSearchIndexes(db *gorm.DB, models ...interface{}) error {
	if db.Dialector.Name() != "postgres" {
		return nil
	}
	type column struct{ table, name string }
	var columns []column
	for _, model := range models {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			return fmt.Errorf("reading %T: %w", model, err)
		}
		for _, field := range stmt.Schema.Fields {
			if field.Tag.Get("search") == "trigram" && field.DBName != "" {
				columns = append(columns, column{stmt.Schema.Table, field.DBName})
			}
		}
	}
	if len(columns) == 0 {
		return nil
	}
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm").Error; err != nil {
		log.Printf("Search indexes skipped: the pg_trgm extension could not be created (%v). Search still works, without an index.", err)
		return nil
	}
	quote := &gorm.Statement{DB: db}
	for _, c := range columns {
		name := "idx_" + c.table + "_" + c.name + "_trgm"
		// CONCURRENTLY, so a big table keeps taking writes while it builds.
		sql := fmt.Sprintf("CREATE INDEX CONCURRENTLY IF NOT EXISTS %s ON %s USING gin (lower(%s) gin_trgm_ops)",
			quote.Quote(name), quote.Quote(c.table), quote.Quote(c.name))
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("creating the search index %s: %w", name, err)
		}
	}
	return nil
}
`
}

// paginateCountTestGo emits internal/paginate/count_test.go.
func paginateCountTestGo() string {
	return `package paginate

import (
	"context"
	"testing"
	"time"
)

// A total is reused across pages until a write ends it: through GORM at once,
// behind its back once CountTTL passes.
func TestListReusesATotalUntilAWrite(t *testing.T) {
	db := newDB(t)
	if err := Install(db); err != nil {
		t.Fatal(err)
	}
	seed(t, db, 30)
	list := func(page int) int64 {
		t.Helper()
		// Through WithContext, as a service lists: a new session every call.
		res, err := List[widget](db.WithContext(context.Background()).Model(&widget{}), Params{Page: page, PageSize: 10}, Config{})
		if err != nil {
			t.Fatal(err)
		}
		return res.Meta.Total
	}
	if got := list(1); got != 30 {
		t.Fatalf("total %d, want 30", got)
	}

	// A row GORM did not write: the handle cannot see it, so the total stands.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sqlDB.Exec("INSERT INTO widgets (id, name, rank, created_at) VALUES ('raw-1', 'raw', 0, ?)", time.Now()); err != nil {
		t.Fatal(err)
	}
	if got := list(2); got != 30 {
		t.Errorf("page 2 total %d; the total was counted again instead of reused", got)
	}

	// A write through GORM ends it at once.
	if err := db.Create(&widget{ID: "gorm-1", Name: "gorm", CreatedAt: time.Now()}).Error; err != nil {
		t.Fatal(err)
	}
	if got := list(3); got != 32 {
		t.Errorf("total %d after a write, want 32", got)
	}

	// And CountTTL ends one it could not see.
	defer func(ttl time.Duration) { CountTTL = ttl }(CountTTL)
	CountTTL = time.Millisecond
	list(1)
	if _, err := sqlDB.Exec("INSERT INTO widgets (id, name, rank, created_at) VALUES ('raw-2', 'raw', 0, ?)", time.Now()); err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Millisecond)
	if got := list(1); got != 33 {
		t.Errorf("total %d after CountTTL, want 33", got)
	}
}

// Different filters are different totals.
func TestCountCacheKeysOnTheQuery(t *testing.T) {
	db := newDB(t)
	if err := Install(db); err != nil {
		t.Fatal(err)
	}
	seed(t, db, 12)
	all, err := List[widget](db.Model(&widget{}), Params{Page: 1, PageSize: 5}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	one, err := List[widget](db.Model(&widget{}).Where("id = ?", pad(0)), Params{Page: 1, PageSize: 5}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	if all.Meta.Total != 12 || one.Meta.Total != 1 {
		t.Errorf("totals %d and %d, want 12 and 1", all.Meta.Total, one.Meta.Total)
	}
}

// Without Install every page counts, as before.
func TestListCountsEveryPageWithoutInstall(t *testing.T) {
	db := newDB(t)
	seed(t, db, 5)
	if _, err := List[widget](db.Model(&widget{}), Params{Page: 1, PageSize: 2}, Config{}); err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sqlDB.Exec("INSERT INTO widgets (id, name, rank, created_at) VALUES ('raw-1', 'raw', 0, ?)", time.Now()); err != nil {
		t.Fatal(err)
	}
	res, err := List[widget](db.Model(&widget{}), Params{Page: 1, PageSize: 2}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Meta.Total != 6 {
		t.Errorf("total %d, want 6", res.Meta.Total)
	}
}

// Trigram indexes are a Postgres feature; anywhere else there is nothing to do.
func TestSearchIndexesAreForPostgresOnly(t *testing.T) {
	type tagged struct {
		ID   string ` + "`" + `gorm:"primarykey"` + "`" + `
		Name string ` + "`" + `search:"trigram"` + "`" + `
	}
	if err := EnsureSearchIndexes(newDB(t), &tagged{}); err != nil {
		t.Errorf("SQLite: %v", err)
	}
}
`
}
