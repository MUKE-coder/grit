package jobs

import (
	"testing"
	"time"
)

// A project without Redis has no StatsCache, and the health check still asks it.
func TestNilStatsCacheHasNoCounts(t *testing.T) {
	var s *StatsCache
	if _, ok := s.Snapshot(); ok {
		t.Error("a nil StatsCache reported counts")
	}
	if err := s.Close(); err != nil {
		t.Error(err)
	}
}

func TestStatsCacheRefusesABadRedisURL(t *testing.T) {
	if _, err := NewStatsCache("not a redis url", time.Second); err == nil {
		t.Error("a bad Redis URL was accepted")
	}
}

// Snapshot must answer at once even when Redis cannot be reached: the health
// check calls it on every probe.
func TestSnapshotDoesNotWaitOnRedis(t *testing.T) {
	s, err := NewStatsCache("redis://127.0.0.1:1", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	started := time.Now()
	if _, ok := s.Snapshot(); ok {
		t.Error("counts reported before any were read")
	}
	if took := time.Since(started); took > 50*time.Millisecond {
		t.Errorf("Snapshot took %v; it must not wait on Redis", took)
	}
}
