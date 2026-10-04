package cron

import (
	"log"

	"library/apps/api/internal/cache"
	"library/apps/api/internal/config"
)

// Start initializes the cron scheduler with the project's default tasks and
// begins executing scheduled jobs in the background.
//
// Returns the running Scheduler (so callers can Stop() it on shutdown) and
// any startup error. The cache argument is reserved for tasks that need
// shared cache access — it is currently unused but kept stable so adding a
// cache-backed task later is a drop-in change.
func Start(cfg *config.Config, _ *cache.Cache) (*Scheduler, error) {
	if cfg.RedisURL == "" {
		log.Println("Cron scheduler disabled: REDIS_URL not configured")
		return nil, nil
	}

	s, err := New(cfg.RedisURL)
	if err != nil {
		return nil, err
	}

	if err := s.Start(); err != nil {
		return nil, err
	}

	log.Println("Cron scheduler ready: it runs on the replica holding the cron lock")
	return s, nil
}
