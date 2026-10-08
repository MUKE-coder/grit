package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
)

// SeedPlan describes how to top one table up to a number of rows.
//
// Make builds the row with the given index. The index runs from the number of
// rows already in the table up to Target, so a seeder that mixes it into a
// unique column (an email, a slug, a phone number) never collides with rows an
// earlier run inserted, and a run that stopped partway resumes where it
// stopped.
type SeedPlan[T any] struct {
	Target int64
	Make   func(i int64) T
	// Batch is the rows per INSERT. Zero picks the largest batch the database
	// accepts for this table's column count, up to 1,000.
	Batch int
	// Writers is how many batches are written at once. Zero picks 1 on SQLite,
	// which has a single writer lock, and 3 elsewhere.
	Writers int
}

// SeedTopUp inserts rows until the table holds plan.Target of them.
//
// It counts what is there first and inserts only the difference, so running it
// twice does nothing the second time. Rows are built in chunks by a small pool
// of goroutines and written in batches, one transaction per batch, so memory
// stays flat however many rows are asked for. GORM hooks still run for every
// row. The first batch that fails stops the run and is returned as the error:
// a seeder that logs and carries on reports success over a half-empty table.
func SeedTopUp[T any](db *gorm.DB, name string, plan SeedPlan[T]) error {
	if plan.Make == nil {
		return fmt.Errorf("seeding %s: no Make function", name)
	}
	var have int64
	if err := db.Model(new(T)).Count(&have).Error; err != nil {
		return fmt.Errorf("seeding %s: counting existing rows: %w", name, err)
	}
	if have >= plan.Target {
		log.Printf("%s: %d rows, target %d, nothing to seed", name, have, plan.Target)
		return nil
	}
	need := plan.Target - have

	batch := plan.Batch
	if batch <= 0 {
		batch = seedBatchSize(db, new(T))
	}
	writers := plan.Writers
	if writers <= 0 {
		writers = 3
		if db.Dialector.Name() == "sqlite" {
			writers = 1
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chunks := make(chan []T, writers*2)
	var next atomic.Int64 // index of the next chunk to build
	chunkCount := (need + int64(batch) - 1) / int64(batch)

	makers := runtime.NumCPU()
	if makers > 4 {
		makers = 4
	}
	var makersDone sync.WaitGroup
	for m := 0; m < makers; m++ {
		makersDone.Add(1)
		go func() {
			defer makersDone.Done()
			for {
				k := next.Add(1) - 1
				if k >= chunkCount {
					return
				}
				start := have + k*int64(batch)
				end := start + int64(batch)
				if end > plan.Target {
					end = plan.Target
				}
				rows := make([]T, 0, end-start)
				for i := start; i < end; i++ {
					rows = append(rows, plan.Make(i))
				}
				select {
				case chunks <- rows:
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() { makersDone.Wait(); close(chunks) }()

	var written atomic.Int64
	var firstErr error
	var errOnce sync.Once
	fail := func(err error) {
		errOnce.Do(func() { firstErr = err; cancel() })
	}

	started := time.Now()
	stopProgress := make(chan struct{})
	go func() {
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				logSeedProgress(name, written.Load(), need, started)
			case <-stopProgress:
				return
			}
		}
	}()

	var writersDone sync.WaitGroup
	for w := 0; w < writers; w++ {
		writersDone.Add(1)
		go func() {
			defer writersDone.Done()
			for rows := range chunks {
				if ctx.Err() != nil {
					continue // drain so the makers can finish
				}
				err := db.Transaction(func(tx *gorm.DB) error {
					return tx.CreateInBatches(rows, len(rows)).Error
				})
				if err != nil {
					fail(fmt.Errorf("seeding %s: inserting rows %d to %d: %w", name, have+written.Load(), have+written.Load()+int64(len(rows)), err))
					continue
				}
				written.Add(int64(len(rows)))
			}
		}()
	}
	writersDone.Wait()
	close(stopProgress)

	if firstErr != nil {
		return firstErr
	}
	elapsed := time.Since(started)
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	log.Printf("%s: seeded %d rows in %s (%.0f rows/s), now %d; heap in use %d MB",
		name, written.Load(), elapsed.Round(time.Millisecond), float64(written.Load())/elapsed.Seconds(),
		have+written.Load(), mem.HeapInuse>>20)
	return nil
}

func logSeedProgress(name string, done, total int64, started time.Time) {
	elapsed := time.Since(started).Seconds()
	if done == 0 || elapsed == 0 {
		log.Printf("%s: 0 of %d rows", name, total)
		return
	}
	rate := float64(done) / elapsed
	left := time.Duration(float64(total-done)/rate) * time.Second
	log.Printf("%s: %d of %d rows (%.0f rows/s, about %s left)", name, done, total, rate, left.Round(time.Second))
}

// seedBatchSize is the largest batch the database accepts for model, up to
// 1,000 rows. A multi-row INSERT sends one parameter per column per row, and
// each database caps the parameters in one statement.
func seedBatchSize(db *gorm.DB, model any) int {
	limit := 65535 // Postgres and MySQL
	if db.Dialector.Name() == "sqlite" {
		limit = 32766
	}
	stmt := &gorm.Statement{DB: db}
	cols := 32
	if err := stmt.Parse(model); err == nil && len(stmt.Schema.DBNames) > 0 {
		cols = len(stmt.Schema.DBNames)
	}
	size := limit / cols
	if size > 1000 {
		size = 1000
	}
	if size < 1 {
		size = 1
	}
	return size
}

// seedTarget is one seeder that can be run on its own with a row count.
type seedTarget struct {
	name  string
	names []string
	run   func(db *gorm.DB, target int64) error
}

var (
	seedTargetsMu sync.Mutex
	seedTargets   []seedTarget
)

// RegisterSeeder makes a seeder reachable as "grit seed <Resource> --count N".
// Generated seeders call it from init; names are matched without regard to
// case, so "Contact", "contact" and "contacts" all find the same seeder.
func RegisterSeeder(name string, run func(db *gorm.DB, target int64) error, aliases ...string) {
	seedTargetsMu.Lock()
	defer seedTargetsMu.Unlock()
	names := []string{strings.ToLower(name)}
	for _, a := range aliases {
		names = append(names, strings.ToLower(a))
	}
	seedTargets = append(seedTargets, seedTarget{name: name, names: names, run: run})
}

// ErrUnknownSeeder is returned by SeedOne for a name no seeder registered.
var ErrUnknownSeeder = errors.New("no seeder with that name supports --count")

// SeedOne tops one resource's table up to target rows.
func SeedOne(db *gorm.DB, name string, target int64) error {
	want := strings.ToLower(name)
	seedTargetsMu.Lock()
	var found *seedTarget
	var known []string
	for i := range seedTargets {
		known = append(known, seedTargets[i].name)
		for _, n := range seedTargets[i].names {
			if n == want {
				found = &seedTargets[i]
			}
		}
	}
	seedTargetsMu.Unlock()
	if found == nil {
		hint := "no seeder in this project supports it yet"
		if len(known) > 0 {
			sort.Strings(known)
			hint = "seeders that do: " + strings.Join(known, ", ")
		}
		return fmt.Errorf("%w: %q (%s). A seeder generated before --count existed can be regenerated with grit generate seeder <Resource> --faker",
			ErrUnknownSeeder, name, hint)
	}
	return found.run(db, target)
}
