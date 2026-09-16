package scaffold

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review M16 to M20, in the files upgrade does not deliver whole.
// The audit writer and prune (audit.go), the sync push (handlers/sync.go), the
// outbox relay (outbox/relay.go, events/durable.go) and the orphan cleanup
// (files/lifecycle.go, storage/storage.go) arrive whole; the activity service,
// the cron schedule, the prune job and main.go are repaired here.

// ─── M18: activity rows are written by a batching writer ─────────────────────

// activityRowFuncs is LogActivity and the row builder the batching writer
// shares with it, in services/activity.go.
const activityRowFuncs = `// LogActivity writes a UserActivity row before it returns. Picks actor + IP +
// user-agent from the request context automatically, falling back to
// args.UserID when the caller is in an unauthenticated handler (auth flows).
func LogActivity(db *gorm.DB, c *gin.Context, args ActivityArgs) {
	if err := LogActivityErr(db, c, args); err != nil {
		// Audit failures are non-fatal but worth knowing about: log and keep
		// moving. In production, wire a metric to alert on a sudden surge in
		// these (suggests DB write pressure).
		log.Printf("activity: failed to write %s: %v", args.Action, err)
	}
}

// LogActivityErr is LogActivity for the handful of events where losing the row
// matters more than the noise: a GDPR erasure, a key revocation, anything an
// auditor will later ask you to produce. It writes the same row and hands back
// the error instead of logging it, so the caller can tell the client the audit
// trail is incomplete.
func LogActivityErr(db *gorm.DB, c *gin.Context, args ActivityArgs) error {
	row := activityRow(c, args)
	return db.Create(&row).Error
}

// activityRow builds the row for args from the request: the actor, the client
// IP and the user agent. It reads the request, so it runs before the handler
// returns even when the row is written later.
//
// The building itself is activityRowCtx, in request_meta.go, which takes those
// three facts from a context.Context instead. This is the gin-shaped front door
// to it, and the only reason it still exists is that handlers hold a gin
// context and jobs do not.
func activityRow(c *gin.Context, args ActivityArgs) models.UserActivity {
	return activityRowCtx(ContextOf(c), args)
}`

// The create, update and delete helpers queue their row instead of writing it.
const (
	cudActivityCallOld = "\tLogActivity(db, c, ActivityArgs{\n\t\tAction:       strings.ToLower(entityType) + \"."
	cudActivityCallNew = "\tqueueActivity(db, c, ActivityArgs{\n\t\tAction:       strings.ToLower(entityType) + \"."
)

// servicesActivityWriterGo is internal/services/activity_writer.go.
func servicesActivityWriterGo() string {
	return `package services

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// The activity feed's create, update and delete rows are written here, in
// batches, rather than inside the request that caused them. Written inline,
// each was a transaction of its own in every create, update and delete, which
// on Postgres is three round trips before the response goes out.
//
// The row is built in the request, while its IP address and user agent can
// still be read, and handed to a writer. A burst larger than the queue writes
// inline rather than dropping rows. FlushActivity, called at shutdown, waits
// for what is still queued.
//
// Sign-ins, security events and other LogActivity calls are still written
// before the request returns.

const (
	activityQueueSize = 4096
	activityBatchSize = 256
	// activityLinger is how long the writer waits for more rows after the first
	// before it writes. Without it, steady traffic that never arrives at the
	// same instant is written one row at a time, which saves the request its
	// wait but not the database its commits.
	activityLinger = 50 * time.Millisecond
)

type activityWriter struct {
	db   *gorm.DB
	rows chan models.UserActivity
}

var (
	activityWriters sync.Map // *gorm.Config -> *activityWriter
	activityPending atomic.Int64
)

// queueActivity builds the row now and writes it soon.
func queueActivity(db *gorm.DB, c *gin.Context, args ActivityArgs) {
	row := activityRow(c, args)
	w := activityWriterFor(db)
	activityPending.Add(1)
	select {
	case w.rows <- row:
	default:
		w.write([]models.UserActivity{row})
	}
}

// activityWriterFor returns the writer for db's connection, starting it on
// first use. Keyed by the connection's config rather than the handle: a handle
// bound to a request's context is a new *gorm.DB every time, and its context
// ends with the request, before the row is written.
func activityWriterFor(db *gorm.DB) *activityWriter {
	if w, ok := activityWriters.Load(db.Config); ok {
		return w.(*activityWriter)
	}
	w := &activityWriter{
		db:   db.Session(&gorm.Session{NewDB: true, Context: context.Background()}),
		rows: make(chan models.UserActivity, activityQueueSize),
	}
	actual, loaded := activityWriters.LoadOrStore(db.Config, w)
	if !loaded {
		go w.run()
	}
	return actual.(*activityWriter)
}

func (w *activityWriter) run() {
	for first := range w.rows {
		batch := []models.UserActivity{first}
		linger := time.NewTimer(activityLinger)
	collect:
		for len(batch) < activityBatchSize {
			select {
			case row := <-w.rows:
				batch = append(batch, row)
			case <-linger.C:
				break collect
			}
		}
		linger.Stop()
		w.write(batch)
	}
}

func (w *activityWriter) write(rows []models.UserActivity) {
	defer activityPending.Add(-int64(len(rows)))
	err := w.db.CreateInBatches(rows, activityBatchSize).Error
	if err == nil {
		return
	}
	if len(rows) == 1 {
		log.Printf("activity: failed to write %s: %v", rows[0].Action, err)
		return
	}
	// One bad row must not take the rest of the batch with it.
	for i := range rows {
		if err := w.db.Create(&rows[i]).Error; err != nil {
			log.Printf("activity: failed to write %s: %v", rows[i].Action, err)
		}
	}
}

// FlushActivity waits until every queued activity row is written, or ctx ends.
// cmd/server/main.go calls it once the server has stopped taking requests.
func FlushActivity(ctx context.Context) {
	for activityPending.Load() > 0 {
		select {
		case <-ctx.Done():
			log.Printf("activity: %d rows were not written before shutdown", activityPending.Load())
			return
		case <-time.After(10 * time.Millisecond):
		}
	}
}
`
}

const (
	subscribersAuditCommentOld = `// Sync, and deliberately so. The activity row should exist before the caller
// is told the write succeeded, and this is the one subscriber that legitimately
// needs the request context: the feed records IP and user agent.`
	subscribersAuditCommentNew = `// Sync, because this is the one subscriber that needs the request context: the
// feed records IP and user agent. The row is built here and written by a
// batching writer (services/activity_writer.go), so the request does not wait
// on the insert.`
)

// ─── M19: the relay stops on shutdown ────────────────────────────────────────

// shutdownDrain follows srv.Shutdown in main.go.
const shutdownDrain = `
	// After the last request: stop claiming outbox messages this process will
	// not live to deliver, and write the activity rows still queued.
	events.StopRelay()
	services.FlushActivity(ctx)
`

var serverShutdownBlock = regexp.MustCompile(`(?m)^\tif err := srv\.Shutdown\(ctx\); err != nil \{\n(?:\t\t[^\n]*\n)+\t\}\n`)

func repairShutdownDrainSource(src, module string) (string, []string, []string) {
	if strings.Contains(src, "events.StopRelay()") || !strings.Contains(src, "srv.Shutdown(ctx)") {
		return src, nil, nil
	}
	warn := []string{"main.go is not the file Grit wrote: after srv.Shutdown, call events.StopRelay() and services.FlushActivity(ctx), or a replica that stops keeps claiming outbox messages and loses the activity rows still queued"}
	loc := serverShutdownBlock.FindStringIndex(src)
	if loc == nil {
		return src, nil, warn
	}
	out := src[:loc[1]] + shutdownDrain + src[loc[1]:]
	for _, pkg := range []string{"/internal/events", "/internal/services"} {
		var ok bool
		if out, ok = addImportGroup(out, module+pkg); !ok {
			return src, nil, warn
		}
	}
	return out, []string{"on shutdown the server stops the outbox relay and writes the activity rows still queued"}, nil
}

// ─── M20: scheduled jobs retry less, and the prune goes in chunks ────────────

// cronRetryOptions gives each built-in scheduled task a retry limit and a
// deadline in place of asynq's 25 retries.
var cronRetryOptions = [][2]string{
	{`asynq.NewTask("tokens:cleanup", nil))`, `asynq.NewTask("tokens:cleanup", nil), asynq.MaxRetry(3), asynq.Timeout(10*time.Minute))`},
	{`asynq.NewTask("uploads:cleanup_orphans", nil))`, `asynq.NewTask("uploads:cleanup_orphans", nil), asynq.MaxRetry(3), asynq.Timeout(time.Hour))`},
	{`asynq.NewTask("audit:prune", nil))`, `asynq.NewTask("audit:prune", nil), asynq.MaxRetry(3), asynq.Timeout(2*time.Hour))`},
}

func repairCronRetriesSource(src string) (string, []string, []string) {
	out := src
	for _, pair := range cronRetryOptions {
		out = strings.Replace(out, pair[0], pair[1], 1)
	}
	if out == src {
		return src, nil, nil
	}
	if !strings.Contains(out, "\t\"time\"\n") {
		var ok bool
		if out, ok = addImportGroup(out, "time"); !ok {
			return src, nil, []string{"cron/cron.go is not the file Grit wrote: give the built-in tasks asynq.MaxRetry(3) and a Timeout, or a failing job is retried 25 times"}
		}
	}
	return out, []string{"scheduled jobs retry 3 times and have a deadline, instead of 25 retries"}, nil
}

// auditPruneHandler is handleAuditPrune in jobs/workers.go.
const auditPruneHandler = `// handleAuditPrune trims the tamper-evident activity log through audit.Prune,
// which deletes in chunks of audit.PruneChunk, one transaction each, and records
// every chunk in the chain so what remains still verifies.
//
// Retention is AUDIT_RETENTION_DAYS, default 365. Zero or negative disables
// pruning entirely, which is the right default for anyone who has not thought
// about it: silently discarding audit history is worse than a large table.
func handleAuditPrune(deps WorkerDeps) func(ctx context.Context, task *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		if deps.DB == nil {
			return fmt.Errorf("database not configured")
		}

		days := 365
		if v := os.Getenv("AUDIT_RETENTION_DAYS"); v != "" {
			parsed, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("AUDIT_RETENTION_DAYS must be a whole number of days, got %q", v)
			}
			days = parsed
		}
		if days <= 0 {
			log.Println("Audit prune skipped: AUDIT_RETENTION_DAYS <= 0 (retain forever)")
			return nil
		}

		removed, err := audit.Prune(ctx, deps.DB, time.Now().AddDate(0, 0, -days))
		if err != nil {
			return fmt.Errorf("pruning the activity log after %d entries: %w", removed, err)
		}
		if removed > 0 {
			log.Printf("Audit prune complete: removed %d entries older than %d days", removed, days)
		}
		return nil
	}
}`

const auditPruneSignature = "func handleAuditPrune(deps WorkerDeps) func(ctx context.Context, task *asynq.Task) error {"

func repairAuditPruneSource(src, module string) (string, []string, []string) {
	if strings.Contains(src, "audit.Prune(ctx, deps.DB,") || !strings.Contains(src, auditPruneSignature) {
		return src, nil, nil
	}
	out, ok := replaceGoFunc(src, auditPruneSignature, auditPruneHandler+"\n")
	if !ok {
		return src, nil, []string{"jobs/workers.go is not the file Grit wrote: prune the activity log with audit.Prune, or pruning breaks the chain's verification"}
	}
	for _, path := range []string{"errors", module + "/internal/models", "gorm.io/gorm"} {
		out, _ = dropUnusedImport(out, path)
	}
	return out, []string{"the activity-log prune deletes in chunks and leaves a chain that still verifies"}, nil
}

func repairActivitySource(src string) (string, []string, []string) {
	if strings.Contains(src, "func activityRow(") || !strings.Contains(src, "func LogActivity(db *gorm.DB, c *gin.Context, args ActivityArgs) {") {
		return src, nil, nil
	}
	warn := []string{"services/activity.go is not the file Grit wrote: build the row in the request and write it from a batching writer, or every create, update and delete waits on an activity insert"}
	if !strings.Contains(src, "args.Metadata") || !strings.Contains(src, "ResolveClientIP(c)") || !strings.Contains(src, cudActivityCallOld) {
		return src, nil, warn
	}
	out, ok := replaceGoFunc(src, "func LogActivity(db *gorm.DB, c *gin.Context, args ActivityArgs) {", activityRowFuncs+"\n")
	if !ok {
		return src, nil, warn
	}
	out = strings.ReplaceAll(out, cudActivityCallOld, cudActivityCallNew)
	return out, []string{"create, update and delete activity rows are written in batches after the request"}, nil
}

// repairBackgroundWork applies M18 to M20 to the files an existing project keeps.
func repairBackgroundWork(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	services := filepath.Join(apiRoot, "internal", "services")
	writer := filepath.Join(services, "activity_writer.go")

	if activity := filepath.Join(services, "activity.go"); fileExists(activity) {
		if err := repairSourceFile(root, m, activity, repairActivitySource); err != nil {
			return err
		}
		if fileContains(activity, "queueActivity(db, c,") && !fileExists(writer) {
			if err := writeFile(writer, strings.ReplaceAll(servicesActivityWriterGo(), "{{MODULE}}", opts.Module())); err != nil {
				return fmt.Errorf("writing %s: %w", writer, err)
			}
		}
	}
	if subs := filepath.Join(services, "event_subscribers.go"); fileExists(subs) {
		if err := repairSourceFile(root, m, subs, func(src string) (string, []string, []string) {
			if !strings.Contains(src, subscribersAuditCommentOld) || !fileExists(writer) {
				return src, nil, nil
			}
			return strings.Replace(src, subscribersAuditCommentOld, subscribersAuditCommentNew, 1), []string{"the audit subscriber's comment describes the batching writer"}, nil
		}); err != nil {
			return err
		}
	}

	if cron := filepath.Join(apiRoot, "internal", "cron", "cron.go"); fileExists(cron) {
		if err := repairSourceFile(root, m, cron, repairCronRetriesSource); err != nil {
			return err
		}
	}
	workers := filepath.Join(apiRoot, "internal", "jobs", "workers.go")
	if fileExists(workers) && fileContains(filepath.Join(apiRoot, "internal", "audit", "audit.go"), "func Prune(ctx context.Context") {
		if err := repairSourceFile(root, m, workers, func(src string) (string, []string, []string) {
			return repairAuditPruneSource(src, opts.Module())
		}); err != nil {
			return err
		}
	}

	// The clients Grit ships push in chunks the server accepts.
	if engine := filepath.Join(root, "packages", "sync", "src", "engine.ts"); fileExists(engine) {
		if err := repairTextFile(root, m, engine, repairSyncEnginePushSource); err != nil {
			return err
		}
	}
	if engine := filepath.Join(root, "apps", "desktop", "sync", "engine.go"); fileExists(engine) {
		if err := repairSourceFile(root, m, engine, repairDesktopPushSource); err != nil {
			return err
		}
	}

	if !fileContains(filepath.Join(apiRoot, "internal", "events", "durable.go"), "func StopRelay()") || !fileExists(writer) {
		return nil
	}
	for _, server := range []string{filepath.Join(apiRoot, "cmd", "server", "main.go"), filepath.Join(root, "main.go")} {
		if !fileExists(server) {
			continue
		}
		if err := repairSourceFile(root, m, server, func(src string) (string, []string, []string) {
			return repairShutdownDrainSource(src, opts.Module())
		}); err != nil {
			return err
		}
	}
	return nil
}

// ─── M17: the sync clients push in chunks ────────────────────────────────────

const (
	tsLastSyncedKey = "const LAST_SYNCED_KEY = \"last_synced_at\";\n"
	tsMaxPushConst  = "/** The most changes the server takes in one push. */\nconst MAX_PUSH_CHANGES = 500;\n"
	tsPushHeadOld   = `    const all = await this.adapter.listOutbox();
    const entries = all
      .filter((e) => !e.hasConflict)
      .sort((a, b) => a.createdAt - b.createdAt);
    if (entries.length === 0) return { pushed: 0, conflicts: 0, overridden: 0 };

`
	tsPushHeadNew = `    const all = await this.adapter.listOutbox();
    const queued = all
      .filter((e) => !e.hasConflict)
      .sort((a, b) => a.createdAt - b.createdAt);
    if (queued.length === 0) return { pushed: 0, conflicts: 0, overridden: 0 };

    let pushed = 0;
    let conflicts = 0;
    let overridden = 0;
    // The server takes at most MAX_PUSH_CHANGES a push, so a long offline
    // stretch goes up in several.
    for (let start = 0; start < queued.length; start += MAX_PUSH_CHANGES) {
      const entries = queued.slice(start, start + MAX_PUSH_CHANGES);
`
	tsPushDecl = "    let pushed = 0;\n    let conflicts = 0;\n    let overridden = 0;\n"
	tsPushTail = "    await this.refreshCounts();\n    return { pushed, conflicts, overridden };"
)

// repairSyncEnginePushSource makes packages/sync push its outbox in chunks of
// MAX_PUSH_CHANGES: the loop that sent it whole now runs once per chunk.
func repairSyncEnginePushSource(src string) (string, []string, []string) {
	if strings.Contains(src, "MAX_PUSH_CHANGES") || !strings.Contains(src, "/sync/push") {
		return src, nil, nil
	}
	warn := []string{"packages/sync/src/engine.ts is not the file Grit wrote: push the outbox in chunks of 500, or a device with more queued changes is refused by the server"}
	head := strings.Index(src, tsPushHeadOld)
	if head < 0 || strings.Count(src, tsPushHeadOld) != 1 || !strings.Contains(src, tsLastSyncedKey) {
		return src, nil, warn
	}
	bodyStart := head + len(tsPushHeadOld)
	bodyLen := strings.Index(src[bodyStart:], tsPushTail)
	if bodyLen < 0 {
		return src, nil, warn
	}
	body := src[bodyStart : bodyStart+bodyLen]
	if strings.Count(body, tsPushDecl) != 1 {
		return src, nil, warn
	}
	body = strings.TrimRight(strings.Replace(body, tsPushDecl, "", 1), "\n") + "\n"
	lines := strings.SplitAfter(body, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			lines[i] = "  " + line
		}
	}
	out := src[:head] + tsPushHeadNew + strings.Join(lines, "") + "    }\n\n" + src[bodyStart+bodyLen:]
	out = strings.Replace(out, tsLastSyncedKey, tsLastSyncedKey+tsMaxPushConst, 1)
	return out, []string{"the web sync client pushes its outbox 500 changes at a time"}, nil
}

const (
	desktopPushBatchDoc = "// PushBatch is the JSON shape /api/sync/push expects.\n"
	desktopMaxPushConst = "// MaxPushChanges is the most changes the server takes in one push.\nconst MaxPushChanges = 500\n\n"
	desktopPushOld      = "\treturn e.pushEntries(entries)\n}\n"
	desktopPushNew      = `	// The server takes at most MaxPushChanges a push, so a long offline stretch
	// goes up in several.
	pushed, conflicts := 0, 0
	for start := 0; start < len(entries); start += MaxPushChanges {
		end := start + MaxPushChanges
		if end > len(entries) {
			end = len(entries)
		}
		p, c, err := e.pushEntries(entries[start:end])
		pushed += p
		conflicts += c
		if err != nil {
			return pushed, conflicts, err
		}
	}
	return pushed, conflicts, nil
}
`
)

func repairDesktopPushSource(src string) (string, []string, []string) {
	if strings.Contains(src, "MaxPushChanges") || !strings.Contains(src, "/sync/push") {
		return src, nil, nil
	}
	if strings.Count(src, desktopPushOld) != 1 || strings.Count(src, desktopPushBatchDoc) != 1 {
		return src, nil, []string{"apps/desktop/sync/engine.go is not the file Grit wrote: push the outbox in chunks of 500, or a device with more queued changes is refused by the server"}
	}
	out := strings.Replace(src, desktopPushOld, desktopPushNew, 1)
	out = strings.Replace(out, desktopPushBatchDoc, desktopMaxPushConst+desktopPushBatchDoc, 1)
	return out, []string{"the desktop sync client pushes its outbox 500 changes at a time"}, nil
}
