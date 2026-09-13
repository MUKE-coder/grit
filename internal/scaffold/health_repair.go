package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// healthQueueFields replaced the queue_keys figure on /api/health.
const healthQueueFields = "\t\t\tQueued     *int   `json:\"queued,omitempty\"`\n" +
	"\t\t\tActive     *int   `json:\"active,omitempty\"`\n"

const oldHealthQueueField = "\t\t\tQueueKeys  int    `json:\"queue_keys,omitempty\"`\n"

// healthQueueStatsSetup goes before the health route: the snapshot it reads.
const healthQueueStatsSetup = `	// Queue counts for /api/health, read by a background refresh at most every
	// 30 seconds, so a probe never waits on Redis for them.
	var queueStats *jobs.StatsCache
	if svc.Jobs != nil {
		qs, err := jobs.NewStatsCache(cfg.RedisURL, 30*time.Second)
		if err != nil {
			log.Printf("Queue counts on /api/health are off: %v", err)
		}
		queueStats = qs
	}

`

// healthJobsBlock is the jobs probe inside /api/health.
const healthJobsBlock = `		// Background jobs: up when Redis is, with queue counts from a snapshot
		// refreshed in the background. Counting asynq's keys here ran KEYS inside
		// Redis on every probe, stalling every Redis client while it ran. If asynq
		// isn't wired (Jobs == nil), report unconfigured rather than "down" so the
		// dashboard distinguishes the cases.
		jobsStatus := compStatus{}
		if svc.Jobs != nil && svc.Cache != nil {
			jobsStatus.OK = redisStatus.OK
			if stats, ok := queueStats.Snapshot(); ok {
				jobsStatus.Queued, jobsStatus.Active = &stats.Queued, &stats.Active
			}
		}
`

const healthCheckAnchor = "\t// Health check\n\t// /api/health probes every infrastructure dependency"

// repairHealthQueueProbe brings a project up to the fix for H15 in the
// contact-app review: /api/health counted asynq's keys with KEYS inside Redis
// on every probe. With a million keys that stalled every Redis client for over
// a second, and anyone can call /api/health.
//
// jobs/stats.go is framework code and arrives whole. routes.go is the
// developer's, so the change there anchors on what Grit wrote.
func repairHealthQueueProbe(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	if !fileExists(routes) || !fileExists(filepath.Join(apiRoot, "internal", "jobs", "client.go")) {
		return nil
	}
	stats := filepath.Join(apiRoot, "internal", "jobs", "stats.go")
	for path, content := range map[string]string{
		stats: jobsStatsGo(),
		filepath.Join(apiRoot, "internal", "jobs", "stats_test.go"): jobsStatsTestGo(),
	} {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	if !fileContains(stats, "func NewStatsCache(") {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, routes, repairHealthQueueProbeSource)
}

// oldHealthJobsBlock is the jobs probe as Grit wrote it before v3.253.0, with
// the blank line after it. Only this exact text is replaced.
const oldHealthJobsBlock = `		// Background-jobs queue — count active asynq keys as a liveness
		// signal. If asynq isn't wired (Jobs == nil), report unconfigured
		// rather than "down" so the dashboard distinguishes the cases.
		jobsStatus := compStatus{}
		if svc.Jobs != nil && svc.Cache != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
			defer cancel()
			n, err := svc.Cache.Client().Eval(ctx,
				"local total = 0\nfor _, k in ipairs(redis.call('keys', 'asynq:*')) do total = total + 1 end\nreturn total",
				[]string{}).Int()
			if err == nil {
				jobsStatus.OK = true
				jobsStatus.QueueKeys = n
			} else {
				// Fall back to a simple ping so a "no keys yet" install still
				// reports OK rather than down.
				if perr := svc.Cache.Client().Ping(ctx).Err(); perr == nil {
					jobsStatus.OK = true
				}
			}
		}

`

func repairHealthQueueProbeSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "redis.call('keys', 'asynq:*')") {
		return src, nil, nil
	}
	start := strings.Index(src, "\t\t// Background-jobs queue")
	end := strings.Index(src, "\t\t// Email is")
	if start < 0 || end < start || src[start:end] != oldHealthJobsBlock ||
		strings.Count(src, oldHealthQueueField) != 1 || strings.Count(src, healthCheckAnchor) != 1 {
		return src, nil, []string{"/api/health in routes.go counts asynq keys with KEYS inside Redis and is not the code Grit wrote: read the counts from jobs.StatsCache, or every health probe stalls Redis"}
	}
	out := src[:start] + healthJobsBlock + "\n" + src[end:]
	out = strings.Replace(out, oldHealthQueueField, healthQueueFields, 1)
	out = strings.Replace(out, healthCheckAnchor, healthQueueStatsSetup+healthCheckAnchor, 1)
	return out, []string{"/api/health reads queue counts from a background snapshot instead of running KEYS inside Redis"}, nil
}
