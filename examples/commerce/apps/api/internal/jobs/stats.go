package jobs

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hibiken/asynq"
)

// QueueStats counts the tasks across every queue at one moment.
type QueueStats struct {
	// Queued is waiting to run: pending, scheduled and retrying.
	Queued int
	// Active is running now.
	Active int
	// Failed ran out of retries and was archived.
	Failed int
	// Taken is when the counts were read.
	Taken time.Time
}

// StatsCache serves queue counts without making the caller wait on Redis.
//
// A health check that counted asynq's keys ran KEYS, a scan of the whole
// keyspace, inside Redis on every probe, and Redis serves every other client
// (the cache, rate limits, the job queue itself) from that same thread. The
// inspector reads each queue's list and set sizes instead, and at most once per
// interval, in the background.
type StatsCache struct {
	inspector *asynq.Inspector
	interval  time.Duration

	mu         sync.Mutex
	snapshot   QueueStats
	ready      bool
	refreshing bool
	lastTry    time.Time
}

// NewStatsCache reads the queues at redisURL no more than once per interval.
func NewStatsCache(redisURL string, interval time.Duration) (*StatsCache, error) {
	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parsing redis URL for queue stats: %w", err)
	}
	return &StatsCache{inspector: asynq.NewInspector(opt), interval: interval}, nil
}

// Snapshot returns the latest counts, and whether there are any yet. It never
// waits: when the counts are older than the interval it starts a refresh and
// returns what it has.
func (s *StatsCache) Snapshot() (QueueStats, bool) {
	if s == nil {
		return QueueStats{}, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.refreshing && time.Since(s.lastTry) >= s.interval {
		s.refreshing = true
		s.lastTry = time.Now()
		go s.refresh()
	}
	return s.snapshot, s.ready
}

func (s *StatsCache) refresh() {
	next, err := s.read()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refreshing = false
	if err != nil {
		log.Printf("jobs: reading queue counts: %v", err)
		return
	}
	s.snapshot, s.ready = next, true
}

func (s *StatsCache) read() (QueueStats, error) {
	queues, err := s.inspector.Queues()
	if err != nil {
		return QueueStats{}, err
	}
	stats := QueueStats{Taken: time.Now()}
	for _, name := range queues {
		info, err := s.inspector.GetQueueInfo(name)
		if err != nil {
			return QueueStats{}, fmt.Errorf("queue %s: %w", name, err)
		}
		stats.Queued += info.Pending + info.Scheduled + info.Retry
		stats.Active += info.Active
		stats.Failed += info.Archived
	}
	return stats, nil
}

// Close releases the inspector's Redis connections.
func (s *StatsCache) Close() error {
	if s == nil {
		return nil
	}
	return s.inspector.Close()
}
