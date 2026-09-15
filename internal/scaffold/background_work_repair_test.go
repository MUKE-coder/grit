package scaffold

import (
	"strings"
	"testing"
)

// The LogActivity services/activity.go had before M18.
const oldLogActivityFunc = `// LogActivity writes a UserActivity row. Picks actor + IP + user-agent
// from the request context automatically, falling back to args.UserID
// when the caller is in an unauthenticated handler (auth flows).
func LogActivity(db *gorm.DB, c *gin.Context, args ActivityArgs) {
	userID := args.UserID
	if userID == "" {
		if v, ok := c.Get("user_id"); ok {
			if s, ok := v.(string); ok {
				userID = s
			}
		}
	}

	var metaJSON string
	if args.Metadata != nil {
		if b, err := json.Marshal(args.Metadata); err == nil {
			metaJSON = string(b)
		}
	}

	row := models.UserActivity{
		UserID:       userID,
		Action:       args.Action,
		Severity:     args.Severity,
		Summary:      args.Summary,
		ResourceType: args.ResourceType,
		ResourceID:   args.ResourceID,
		IPAddress:    ResolveClientIP(c),
		UserAgent:    c.GetHeader("User-Agent"),
		Metadata:     metaJSON,
	}

	if err := db.Create(&row).Error; err != nil {
		log.Printf("activity: failed to write %s: %v", args.Action, err)
	}
}`

func TestAuditBatchesWritesAndPrunesKeepTheChain(t *testing.T) {
	src := apiAuditGo()
	for _, want := range []string{
		"return tx.CreateInBatches(entries, 256).Error",
		"func Prune(ctx context.Context, db *gorm.DB, cutoff time.Time) (int64, error) {",
		"pruneRecorded(ctx, db, e.PrevHash)",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("audit.go is missing %q", want)
		}
	}
	if strings.Contains(src, "if err := tx.Create(e).Error; err != nil {") {
		t.Error("the audit writer still inserts entries one at a time")
	}
	mustFormatGo(t, "audit.go", src)

	workers := jobsWorkersGo()
	if !strings.Contains(workers, "audit.Prune(ctx, deps.DB,") || strings.Contains(workers, "DELETE FROM activity_logs") {
		t.Error("the prune job still deletes the whole backlog in one statement")
	}
	mustFormatGo(t, "jobs/workers.go", workers)

	old := "package jobs\n\nimport (\n\t\"context\"\n\t\"errors\"\n\t\"fmt\"\n\t\"log\"\n\t\"os\"\n\t\"strconv\"\n\t\"time\"\n\n\t\"github.com/hibiken/asynq\"\n\t\"gorm.io/gorm\"\n\n\t\"example.com/app/internal/audit\"\n\t\"example.com/app/internal/models\"\n)\n\n" +
		"type WorkerDeps struct{ DB *gorm.DB }\n\n" +
		"// handleAuditPrune trims the log.\n" + auditPruneSignature + "\n\treturn func(ctx context.Context, task *asynq.Task) error {\n\t\tif errors.Is(nil, nil) {\n\t\t\tvar e models.ActivityLog\n\t\t\t_, _ = audit.Canonical(&e)\n\t\t}\n\t\treturn nil\n\t}\n}\n"
	out, fixed, warn := repairAuditPruneSource(old, "example.com/app")
	if len(fixed) != 1 || len(warn) != 0 || !strings.Contains(out, "audit.Prune(ctx, deps.DB,") {
		t.Fatalf("the prune job was not repaired (%v %v):\n%s", fixed, warn, out)
	}
	if strings.Contains(out, "\"errors\"") || strings.Contains(out, "internal/models\"") || !strings.Contains(out, "\"gorm.io/gorm\"") {
		t.Errorf("the repair left the imports wrong:\n%s", out)
	}
	if again, _, _ := repairAuditPruneSource(out, "example.com/app"); again != out {
		t.Error("the prune repair is not idempotent")
	}
	mustFormatGo(t, "the repaired workers.go", out)
}

func TestActivityRowsAreQueued(t *testing.T) {
	current := userActivityServiceGo()
	if strings.Count(current, cudActivityCallNew) != 3 || strings.Contains(current, cudActivityCallOld) {
		t.Error("the create, update and delete helpers do not all queue their rows")
	}
	old := strings.Replace(current, activityRowFuncs, oldLogActivityFunc, 1)
	old = strings.ReplaceAll(old, cudActivityCallNew, cudActivityCallOld)
	out, fixed, warn := repairActivitySource(old)
	if len(fixed) != 1 || len(warn) != 0 || out != current {
		t.Errorf("the repaired activity.go is not the template (%v %v)", fixed, warn)
	}
	if again, _, _ := repairActivitySource(out); again != out {
		t.Error("the activity repair is not idempotent")
	}
	mustFormatGo(t, "services/activity.go", current)
	mustFormatGo(t, "services/activity_writer.go", servicesActivityWriterGo())
	if !strings.Contains(apiEventsSubscribersGo(), subscribersAuditCommentNew) {
		t.Error("the audit subscriber still says the row is written before the request returns")
	}
}

func TestShutdownStopsTheRelayAndFlushesActivity(t *testing.T) {
	for name, src := range map[string]string{"cmd/server/main.go": apiMainGo(Options{}), "main.go": singleMainGo(Options{})} {
		if !strings.Contains(src, shutdownDrain) || !strings.Contains(src, "\t\"{{MODULE}}/internal/events\"\n") {
			t.Errorf("%s does not stop the relay and flush activity on shutdown", name)
		}
		old := strings.Replace(src, shutdownDrain, "", 1)
		old = strings.Replace(old, "\t\"{{MODULE}}/internal/events\"\n", "", 1)
		out, fixed, warn := repairShutdownDrainSource(old, "{{MODULE}}")
		if len(fixed) != 1 || len(warn) != 0 || !strings.Contains(out, shutdownDrain) || !strings.Contains(out, "\t\"{{MODULE}}/internal/events\"\n") {
			t.Errorf("%s was not repaired (%v %v)", name, fixed, warn)
		}
		if again, _, _ := repairShutdownDrainSource(out, "{{MODULE}}"); again != out {
			t.Errorf("the %s repair is not idempotent", name)
		}
		mustFormatGo(t, "the repaired "+name, out)
	}
}

func TestSyncPushIsBounded(t *testing.T) {
	h := apiSyncHandlerGo()
	for _, want := range []string{"const MaxPushChanges = 500", "rows := h.loadCurrent(c, req.Changes)", "current, err := h.current(c, rows, ch, proto)"} {
		if !strings.Contains(h, want) {
			t.Errorf("handlers/sync.go is missing %q", want)
		}
	}
	if strings.Contains(h, `First(current, "id = ?", ch.ID)`) {
		t.Error("a push still reads each row it changes with a query of its own")
	}
	mustFormatGo(t, "handlers/sync.go", h)
	if !strings.Contains(syncEngineTS(), "start += MAX_PUSH_CHANGES") {
		t.Error("the web sync client still sends its whole outbox in one push")
	}
	desktop := desktopSyncEngineGo()
	if !strings.Contains(desktop, "start += MaxPushChanges") {
		t.Error("the desktop sync client still sends its whole outbox in one push")
	}
	mustFormatGo(t, "desktop sync/engine.go", desktop)
}

func TestRelaySkipsLockedRowsPrunesAndStops(t *testing.T) {
	relay := outboxRelayGo()
	for _, want := range []string{`Options: "SKIP LOCKED"`, "Prune(r.DB.WithContext(ctx), r.KeepDelivered)", "r.release(batch[i:])"} {
		if !strings.Contains(relay, want) {
			t.Errorf("outbox/relay.go is missing %q", want)
		}
	}
	mustFormatGo(t, "outbox/relay.go", relay)
	durable := apiEventsDurableGo()
	if !strings.Contains(durable, "func StopRelay() {") || strings.Contains(durable, "go relay.Start(context.Background())") {
		t.Error("the event relay still cannot be stopped")
	}
	mustFormatGo(t, "events/durable.go", durable)
}

func TestCleanupJobsAreBounded(t *testing.T) {
	cron := cronSchedulerGo()
	old := cron
	for _, pair := range cronRetryOptions {
		if !strings.Contains(cron, pair[1]) {
			t.Errorf("cron.go does not register %s", pair[1])
		}
		old = strings.Replace(old, pair[1], pair[0], 1)
	}
	out, fixed, _ := repairCronRetriesSource(old)
	for _, pair := range cronRetryOptions {
		if !strings.Contains(out, pair[1]) {
			t.Errorf("the cron repair did not give %s its options (%v)", pair[0], fixed)
		}
	}
	if again, _, _ := repairCronRetriesSource(out); again != out {
		t.Error("the cron repair is not idempotent")
	}
	mustFormatGo(t, "cron/cron.go", cron)

	lifecycle := filesLifecycleGo()
	if !strings.Contains(lifecycle, "many.DeleteMany(ctx, keys)") || !strings.Contains(lifecycle, "Limit(orphanPage)") {
		t.Error("the orphan cleanup still loads every orphan and deletes each one alone")
	}
	mustFormatGo(t, "files/lifecycle.go", lifecycle)
	storage := storageServiceGo()
	if !strings.Contains(storage, "func (s *Storage) DeleteMany(ctx context.Context, keys []string) error {") {
		t.Error("storage has no batch delete")
	}
	mustFormatGo(t, "storage/storage.go", storage)
}

func TestRepairChunksClientPushes(t *testing.T) {
	current := syncEngineTS()
	head := strings.Index(current, tsPushHeadNew)
	end := strings.Index(current, "    }\n\n"+tsPushTail)
	if head < 0 || end < 0 || !strings.Contains(current, tsLastSyncedKey+tsMaxPushConst) {
		t.Fatal("the web sync client does not push in chunks")
	}
	lines := strings.SplitAfter(current[head+len(tsPushHeadNew):end], "\n")
	for i, line := range lines {
		lines[i] = strings.TrimPrefix(line, "  ")
	}
	body := strings.Replace(strings.Join(lines, "")+"\n", "    for (let i = 0; i < results.length; i += 1) {\n", tsPushDecl+"    for (let i = 0; i < results.length; i += 1) {\n", 1)
	old := strings.Replace(current[:head]+tsPushHeadOld+body+current[end+len("    }\n\n"):], tsMaxPushConst, "", 1)
	out, fixed, warn := repairSyncEnginePushSource(old)
	if out != current || len(fixed) != 1 || len(warn) != 0 {
		t.Errorf("the repaired engine.ts is not the template (%v %v)", fixed, warn)
	}
	if again, _, _ := repairSyncEnginePushSource(out); again != out {
		t.Error("the engine.ts repair is not idempotent")
	}

	desktop := desktopSyncEngineGo()
	old = strings.Replace(strings.Replace(desktop, desktopPushNew, desktopPushOld, 1), desktopMaxPushConst, "", 1)
	out, fixed, warn = repairDesktopPushSource(old)
	if out != desktop || len(fixed) != 1 || len(warn) != 0 {
		t.Errorf("the repaired desktop engine is not the template (%v %v)", fixed, warn)
	}
}
