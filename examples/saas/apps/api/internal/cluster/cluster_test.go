package cluster

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func openCluster(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	return db
}

// A bump on one replica is what tells the others their copy is stale.
func TestABumpIsSeenOnce(t *testing.T) {
	db := openCluster(t)
	w := NewWatch(db, "authz", 0)
	if w.Changed() {
		t.Fatal("nothing was bumped, and the watch reported a change")
	}
	if err := Bump(db, "authz"); err != nil {
		t.Fatalf("bump: %v", err)
	}
	if !w.Changed() {
		t.Fatal("a bump was not seen")
	}
	if w.Changed() {
		t.Fatal("the same bump was reported twice")
	}
	if err := Bump(db, "sso"); err != nil {
		t.Fatalf("bump: %v", err)
	}
	if w.Changed() {
		t.Fatal("another name's bump was reported")
	}
}

// Changed sits on the request path, so it must not read on every call.
func TestAWatchReadsAtMostOncePerInterval(t *testing.T) {
	db := openCluster(t)
	w := NewWatch(db, "authz", time.Hour)
	if err := Bump(db, "authz"); err != nil {
		t.Fatalf("bump: %v", err)
	}
	if w.Changed() {
		t.Fatal("the watch read the database before its interval was up")
	}
}
