package cron

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

// Only one replica runs the scheduler.
//
// asynq's schedulers do not coordinate, and every API replica built and ran
// one, so each scheduled job ran once per replica: two replicas, two token
// sweeps an hour, two orphan-upload cleanups a night. The replica holding a
// Redis lock runs it. The lock is a lease: a replica that dies stops renewing
// it, and another takes over when it lapses.
const (
	leaderKey   = "grit:cron:leader"
	leaderLease = 30 * time.Second
	leaderRenew = 10 * time.Second
)

// renewLease extends the lock only while this replica still holds it.
var renewLease = redis.NewScript(`if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("PEXPIRE", KEYS[1], ARGV[2]) else return 0 end`)

var stopOnce sync.Once

// lead waits for the cron lock, then runs the scheduler for as long as it
// holds it.
func (s *Scheduler) lead(redisURL string) {
	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		log.Printf("[cron] %v", err)
		return
	}
	client, ok := opt.MakeRedisClient().(redis.UniversalClient)
	if !ok {
		// A connection asynq understands that go-redis does not expose here:
		// better to run unlocked than not at all.
		log.Println("[cron] could not take the cron lock; running the scheduler on this replica")
		if err := s.scheduler.Start(); err != nil {
			log.Printf("[cron] %v", err)
		}
		return
	}

	id := replicaID()
	ctx := context.Background()
	for {
		held, err := client.SetNX(ctx, leaderKey, id, leaderLease).Result()
		if err == nil && held {
			break
		}
		time.Sleep(leaderRenew)
	}
	if err := s.scheduler.Start(); err != nil {
		log.Printf("[cron] %v", err)
		return
	}
	log.Println("[cron] this replica holds the cron lock and runs the scheduled jobs")

	t := time.NewTicker(leaderRenew)
	defer t.Stop()
	for range t.C {
		n, err := renewLease.Run(ctx, client, []string{leaderKey}, id, leaderLease.Milliseconds()).Int()
		if err != nil {
			// Redis blinked. The lease has time left in it; try again next tick.
			log.Printf("[cron] renewing the cron lock: %v", err)
			continue
		}
		if n == 0 {
			// This replica stalled past its lease and another took the lock.
			// Stop, so each job still runs once.
			log.Println("[cron] another replica holds the cron lock now; this one stops scheduling")
			s.shutdown()
			return
		}
	}
}

// shutdown stops the scheduler once. asynq's Shutdown closes a channel, and a
// second call would panic.
func (s *Scheduler) shutdown() {
	stopOnce.Do(s.scheduler.Shutdown)
}

func replicaID() string {
	host, _ := os.Hostname()
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return host + ":" + strconv.Itoa(os.Getpid()) + ":" + hex.EncodeToString(b)
}
