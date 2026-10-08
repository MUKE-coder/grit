// Package cluster is what a set of API replicas has to agree on: a generation
// number per piece of per-process state, bumped by whichever replica changes
// the data behind it and watched by the others.
//
// A permission cache, an SSO provider registry: each is built in one process,
// and each was right only in the process that changed it. Behind a load
// balancer, a role revoked through one replica kept working on the rest until
// they restarted. A tiny table in the database the replicas already share is
// the whole mechanism: no Redis, no message bus, the same on Postgres, MySQL
// and SQLite.
//
// It imports nothing from the app, so any package can use it.
package cluster

import (
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Generation is one counter, a row per watched thing.
type Generation struct {
	Name      string `gorm:"primaryKey;size:64"`
	Value     int64  `gorm:"not null;default:0"`
	UpdatedAt time.Time
}

func (Generation) TableName() string { return "cluster_generations" }

// Migrate creates the table. NewWatch and Bump call it, so no migration has to
// know about it. Bumps are rare, a role or a connection changing, so checking
// the schema on each is cheaper than getting a once-per-process guard wrong.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Generation{})
}

// Bump moves name's generation on, telling every other replica its copy of
// whatever name stands for is stale.
func Bump(db *gorm.DB, name string) error {
	if err := Migrate(db); err != nil {
		return err
	}
	res := db.Model(&Generation{}).Where("name = ?", name).Update("value", gorm.Expr("value + 1"))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		return nil
	}
	// The first bump creates the row. Two replicas doing that at once is fine:
	// one insert wins, and either way the value has moved off zero.
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&Generation{Name: name, Value: 1}).Error
}

// Current is name's generation, or 0 before its first bump.
func Current(db *gorm.DB, name string) (int64, error) {
	var g Generation
	err := db.Where("name = ?", name).Limit(1).Find(&g).Error
	return g.Value, err
}

// Watch reports when name's generation has moved. It reads the database at
// most once per interval, and only one caller at a time does the reading, so
// it can sit on a request's hot path.
type Watch struct {
	db       *gorm.DB
	name     string
	interval time.Duration

	mu   sync.Mutex
	next time.Time
	seen int64
}

// NewWatch starts watching name from its current generation.
func NewWatch(db *gorm.DB, name string, interval time.Duration) *Watch {
	w := &Watch{db: db, name: name, interval: interval, next: time.Now().Add(interval)}
	if err := Migrate(db); err == nil {
		w.seen, _ = Current(db, name)
	}
	return w
}

// Changed reports whether the generation has moved since the last call that
// reported a change. A failed read reports none: the local copy is exactly as
// fresh as it was.
func (w *Watch) Changed() bool {
	w.mu.Lock()
	if time.Now().Before(w.next) {
		w.mu.Unlock()
		return false
	}
	w.next = time.Now().Add(w.interval)
	w.mu.Unlock()

	value, err := Current(w.db, w.name)

	w.mu.Lock()
	defer w.mu.Unlock()
	if err != nil || value == w.seen {
		return false
	}
	w.seen = value
	return true
}
