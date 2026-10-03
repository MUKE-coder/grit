package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The /api/health handler, brought to the four states.
//
// Before v3.350.0 every component reported one boolean, and that boolean covered
// two opposite situations: Redis is down, and this deployment has no Redis. The
// overall status then had to exclude the second case by hand, in a condition
// that named Cache and nothing else:
//
//	if !dbStatus.OK || (svc.Cache != nil && !redisStatus.OK) {
//
// which is why a missing mailer and a missing object store were simply not in
// it, and why the admin page guessed at a state per component and guessed
// differently each time. Storage was not reported at all, which is how an upload
// answering 503 came with a System Health page that had nothing to say.
//
// This replaces the handler and eventBusStatus with the ones the generator
// writes today, and adds the import they need.
//
// Only the exact text Grit wrote is replaced. A handler somebody has edited is
// reported and left alone, because the right way to add a component has been
// health.Register since this release, and overwriting an edit to teach that is
// not a trade worth making.

// oldHealthCheckBlock is the handler as Grit wrote it from v3.253.0 to
// v3.349.0: queue counts from the snapshot, the mail driver, and the realtime
// stats beside the event bus.
const oldHealthCheckBlock = `	healthCheck := func(c *gin.Context) {
` +
	"\t\ttype compStatus struct {\n" +
	"\t\t\tOK         bool   `json:\"ok\"`\n" +
	"\t\t\tLatencyMS  int64  `json:\"latency_ms,omitempty\"`\n" +
	"\t\t\tTables     int    `json:\"tables,omitempty\"`\n" +
	"\t\t\tQueued     *int   `json:\"queued,omitempty\"`\n" +
	"\t\t\tActive     *int   `json:\"active,omitempty\"`\n" +
	"\t\t\tConfigured bool   `json:\"configured,omitempty\"`\n" +
	"\t\t\tDriver     string `json:\"driver,omitempty\"`\n" +
	"\t\t\tError      string `json:\"error,omitempty\"`\n" +
	"\t\t}\n" +
	`
		// Database ping + table count. We probe with a 500ms deadline so a
		// blocked write loop can't hang the health check.
		dbStatus := compStatus{OK: true}
		dbStart := time.Now()
		if sqlDB, err := db.DB(); err == nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
			defer cancel()
			if err := sqlDB.PingContext(ctx); err != nil {
				dbStatus.OK = false
				log.Printf("health: database ping failed: %v", err)
			}
		}
		dbStatus.LatencyMS = time.Since(dbStart).Milliseconds()
		if dbStatus.OK {
			// Best-effort table count. Dialect-aware, and 0 rather than an
			// error when the database cannot be asked: a missing tooltip
			// figure is not a health problem.
			dbStatus.Tables = database.TableCount(db)
		}

		// Redis ping. Reuse the same cache client the rest of the app uses
		// rather than opening a new connection — that way "Redis healthy"
		// on the dashboard means the same Redis the cache + jobs use.
		redisStatus := compStatus{}
		if svc.Cache != nil {
			redisStart := time.Now()
			ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
			defer cancel()
			if err := svc.Cache.Client().Ping(ctx).Err(); err != nil {
				redisStatus.OK = false
				log.Printf("health: redis ping failed: %v", err)
			} else {
				redisStatus.OK = true
			}
			redisStatus.LatencyMS = time.Since(redisStart).Milliseconds()
		}

		// Background jobs: up when Redis is, with queue counts from a snapshot
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

		// Email is reported by the driver main.go picked: smtp, resend, log and
		// so on. No mailer is "not configured", which the dashboard shows as a
		// dash rather than as down.
		mailStatus := compStatus{Configured: svc.Mailer != nil, OK: svc.Mailer != nil}
		if svc.Mailer != nil {
			mailStatus.Driver = svc.Mailer.Driver()
		}

		// Overall status — ok if every wired-up component is up. Components
		// that aren't configured (e.g. Redis off in a single-binary dev
		// run) don't drag the overall status down.
		overall := "ok"
		if !dbStatus.OK || (svc.Cache != nil && !redisStatus.OK) {
			overall = "degraded"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   overall,
			"version":  "0.1.0",
			"database": dbStatus,
			"redis":    redisStatus,
			"api":      compStatus{OK: true},
			"jobs":     jobsStatus,
			"email":    mailStatus,
			// The event bus reports itself. Dropped rising is the only signal
			// from outside that a subscriber is too slow or the queue too
			// small, and "did my webhook fire" deserves a better answer than
			// reading logs.
			"events": eventBusStatus(),

			// This replica's sockets, and what its hub delivered, dropped and
			// failed to publish to the others.
			"realtime": realtimeHub.Stats(),
		})
	}
`

// oldEventBusStatus is the event bus report that called a lifetime drop counter
// a failure.
const oldEventBusStatus = `func eventBusStatus() interface{} {
	bus := events.Default()
	if bus == nil {
		return map[string]interface{}{"ok": false, "configured": false}
	}
	s := bus.Stats()
	return map[string]interface{}{
		"ok":          s.Dropped == 0,
		"configured":  true,
		"subscribers": s.Subscribers,
		"queued":      s.Queued,
		"capacity":    s.Capacity,
		"dropped":     s.Dropped,
	}
}
`

// repairHealthStates rewrites /api/health to answer in four states.
func repairHealthStates(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	// internal/health arrives with the framework-owned files, early in the
	// upgrade. Without it this repair would write a handler that does not
	// compile, so it does nothing instead.
	if !fileExists(routes) || !fileContains(filepath.Join(apiRoot, "internal", "health", "health.go"), "func Overall(") {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	if err := repairSourceFile(root, m, routes, repairHealthStatesSource(opts.Module())); err != nil {
		return err
	}
	return nil
}

func repairHealthStatesSource(module string) func(string) (string, []string, []string) {
	return func(src string) (string, []string, []string) {
		if strings.Contains(src, "health.Overall(comps)") {
			return src, nil, nil
		}
		if !strings.Contains(src, "type compStatus struct {") {
			// Not a routes.go with a health check in it at all. Nothing to say:
			// a warning here would fire on every single and every plugin.
			return src, nil, nil
		}
		if strings.Count(src, oldHealthCheckBlock) != 1 || strings.Count(src, oldEventBusStatus) != 1 {
			return src, nil, []string{"/api/health is not the handler Grit wrote, so it still answers one boolean per component: " +
				"\"off\" and \"degraded\" are the same word there, and an operator cannot tell a missing Redis from a dead one. " +
				"Add health.Register(name, probe) for your own components and take the edit back out, or port the handler by hand"}
		}
		// The realtime stats are in the handler this writes, so a project whose
		// hub predates them would get a reference to a variable it has not got.
		// gofmt would catch it after the fact; saying so is better.
		if !strings.Contains(src, "realtimeHub.Stats()") {
			return src, nil, []string{"/api/health has no realtime stats in it yet, so the four-state rewrite is held back one upgrade"}
		}

		out := strings.Replace(src, oldHealthCheckBlock, healthCheckBlock, 1)
		out = strings.Replace(out, oldEventBusStatus, eventBusStatusBlock, 1)
		var ok bool
		if out, ok = addImportGroup(out, module+"/internal/health"); !ok {
			return src, nil, []string{"could not add the internal/health import to routes.go, so /api/health still answers one boolean per component"}
		}
		return out, []string{"/api/health answers ok, degraded, off or unknown per component, reports storage, and merges anything health.Register was given"}, nil
	}
}
