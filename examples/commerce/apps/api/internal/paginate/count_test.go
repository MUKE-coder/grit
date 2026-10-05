package paginate

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
		ID   string `gorm:"primarykey"`
		Name string `search:"trigram"`
	}
	if err := EnsureSearchIndexes(newDB(t), &tagged{}); err != nil {
		t.Errorf("SQLite: %v", err)
	}
}
