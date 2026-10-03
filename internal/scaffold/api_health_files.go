package scaffold

// internal/health, and the /api/health handler that answers in its vocabulary.
//
// Every component used to report one boolean, and that boolean covered two
// opposite situations: Redis is down, and this deployment has no Redis. An
// operator cannot tell those apart, and the dashboard could not either, so it
// guessed, differently per component: Redis read `ok === false ? down : ok ? ok
// : unknown`, email read `configured ? ok : unknown`, and storage was not
// reported at all. Four states put the distinction in the response, where it
// belongs, and leave the page nothing to infer.
//
// The same change is what lets a probe come from somewhere else. health.Register
// takes a name and a function, Snapshot runs them, and the handler merges them
// into the response, so a plugin or anything wired in main.go reports itself
// without editing this handler.

// apiHealthGo emits internal/health/health.go.
func apiHealthGo() string { return tmpl("api/health/health.go") }

// apiHealthTestGo emits internal/health/health_test.go.
func apiHealthTestGo() string { return tmpl("api/health/health_test.go") }

// healthCheckBlock is the /api/health handler, from `healthCheck := func` to its
// closing brace.
//
// One source for the template and for the repair that brings an older project
// to it, because a repair with its own transcription of a handler is a copy that
// drifts, which is what the repair markers in templates/ exist to prevent.
const healthCheckBlock = `	healthCheck := func(c *gin.Context) {
		// Components answer in four states rather than a boolean: ok, degraded,
		// off and unknown. See internal/health for which is which. "off" is a
		// dependency this deployment never asked for and "unknown" is a probe
		// that could not find out, so neither lowers the overall status; only
		// "degraded" does.
		//
		// Snapshot runs whatever health.Register was given, from a plugin or
		// from your own main.go, and those components are reported beside the
		// framework's own.
		comps := health.Snapshot()

		// The API answered this request, which is the whole of what it can
		// honestly say about itself.
		comps["api"] = health.OK("answering requests", nil)

		// Database: a ping bounded at 500ms, so a blocked write loop cannot hang
		// the health check, and a table count for the dashboard's tooltip.
		dbStart := time.Now()
		if sqlDB, err := db.DB(); err != nil {
			comps["database"] = health.Degraded("the database handle is not available", nil)
		} else {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
			defer cancel()
			pingErr := sqlDB.PingContext(ctx)
			fields := map[string]any{"latency_ms": time.Since(dbStart).Milliseconds()}
			if pingErr != nil {
				// The error goes to the log and not into the response. A driver
				// error names the host, the port and often the user, and this
				// endpoint is public so a load balancer can reach it.
				log.Printf("health: database ping failed: %v", pingErr)
				comps["database"] = health.Degraded("the database did not answer a ping", fields)
			} else {
				// Best-effort table count. Dialect-aware, and 0 rather than an
				// error when the database cannot be asked: a missing tooltip
				// figure is not a health problem.
				fields["tables"] = database.TableCount(db)
				comps["database"] = health.OK("", fields)
			}
		}

		// Redis: the same cache client the rest of the app uses, rather than a
		// new connection, so "Redis healthy" here means the Redis the cache and
		// the queue are actually on.
		//
		// A nil cache has two causes and they are not the same news. REDIS_URL=
		// turned it off on purpose. Anything else means main.go dialled Redis at
		// boot, was refused, and carried on with caching, jobs and cron off,
		// which in production is an incident and on a laptop is Tuesday. That is
		// the one thing in this report that asks which environment it is in.
		switch {
		case svc.Cache != nil:
			redisStart := time.Now()
			ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
			defer cancel()
			pingErr := svc.Cache.Client().Ping(ctx).Err()
			fields := map[string]any{"latency_ms": time.Since(redisStart).Milliseconds()}
			if pingErr != nil {
				log.Printf("health: redis ping failed: %v", pingErr)
				comps["redis"] = health.Degraded("Redis did not answer a ping", fields)
			} else {
				comps["redis"] = health.OK("", fields)
			}
		case cfg.RedisURL == "":
			comps["redis"] = health.Off("REDIS_URL is empty, so caching, jobs and cron are off on purpose")
		case cfg.IsDevelopment():
			comps["redis"] = health.Off("Redis was not reachable when the API started, so caching, jobs and cron are off; set REDIS_URL= in .env to mean it")
		default:
			// No address in the detail: this endpoint is public.
			comps["redis"] = health.Degraded("Redis is configured and was not reachable when the API started, so caching, jobs and cron are off; restart the API once it answers", nil)
		}

		// Background jobs: up when Redis is, with queue counts from a snapshot
		// refreshed in the background. Counting asynq's keys here ran KEYS inside
		// Redis on every probe, stalling every Redis client while it ran, and
		// anyone can call /api/health.
		switch {
		case svc.Jobs != nil && svc.Cache != nil:
			if stats, ok := queueStats.Snapshot(); ok {
				comps["jobs"] = health.OK("", map[string]any{"queued": stats.Queued, "active": stats.Active})
			} else {
				// The counts come from a background refresh, so the first probe
				// after a restart has none yet. That is not the queue being
				// unwell, which is the distinction "unknown" exists for.
				comps["jobs"] = health.Unknown("the queue counts have not been sampled yet")
			}
		case comps["redis"].State == health.StateDegraded:
			// The queue is missing because Redis is, and the Redis report has
			// already said why in the words that fit this deployment.
			comps["jobs"] = health.Degraded("the Redis the queue runs on is not there", nil)
		default:
			comps["jobs"] = health.Off("no queue: asynq needs Redis")
		}

		// Email names the driver main.go picked: smtp, resend, log and so on.
		if svc.Mailer == nil {
			comps["email"] = health.Off("no mailer is configured; set MAIL_MAILER")
		} else {
			comps["email"] = health.OK("", map[string]any{"configured": true, "driver": svc.Mailer.Driver()})
		}

		// Storage, which nothing reported until v3.350.0: an upload answered 503,
		// the System Health page had nothing to say about why, and the answer was
		// that no object store was configured in that environment at all.
		//
		// The driver's name only. Describe() also names the bucket and the
		// endpoint, and this endpoint is public.
		if svc.Storage == nil {
			comps["storage"] = health.Off("no object store is configured; set STORAGE_DRIVER")
		} else {
			comps["storage"] = health.OK("", map[string]any{"configured": true, "driver": svc.Storage.Driver()})
		}

		// The event bus reports itself, because "did my webhook fire" deserves a
		// better answer than reading logs.
		comps["events"] = eventBusStatus()

		body := gin.H{
			"status":  health.Overall(comps),
			"version": "0.1.0",
			// This replica's sockets, and what its hub delivered, dropped and
			// failed to publish to the others. Per process, not per deployment.
			"realtime": realtimeHub.Stats(),
		}
		for name, report := range comps {
			// A registered probe cannot take a key the handler already wrote. A
			// plugin that calls itself "status" must not be able to tell a load
			// balancer the deployment is fine.
			if _, taken := body[name]; !taken {
				body[name] = report
			}
		}
		c.JSON(http.StatusOK, body)
	}
`

// eventBusStatusBlock is the event bus's own health report.
//
// It used to report ok as Dropped == 0. Dropped is a counter for the life of the
// process and nothing resets it, so one drop at eight in the morning made the
// deployment unwell until somebody restarted it, which is how an alarm gets
// ignored. The counter is still reported, and it is still the number to watch:
// rising is the signal, a lifetime total is not.
const eventBusStatusBlock = `func eventBusStatus() health.Report {
	bus := events.Default()
	if bus == nil {
		return health.Off("no event bus is running")
	}
	s := bus.Stats()
	return health.OK("", map[string]any{
		"configured":  true,
		"subscribers": s.Subscribers,
		"queued":      s.Queued,
		"capacity":    s.Capacity,
		"dropped":     s.Dropped,
	})
}
`
