package paginate

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
