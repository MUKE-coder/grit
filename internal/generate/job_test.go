package generate

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
)

// The cron docs promised a "grit add cron" that did not exist, and a custom job
// was five hand edits across two framework files, one of which had no marker.

const workersWithMarker = `package jobs

func StartWorker() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailSend, handleEmailSend(deps))
	// grit:jobs

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Printf("Worker error: %v", err)
		}
	}()
}
`

// A project scaffolded before the marker existed.
const workersWithoutMarker = `package jobs

func StartWorker() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailSend, handleEmailSend(deps))

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Printf("Worker error: %v", err)
		}
	}()
}
`

const cronWithMarker = `package cron

func New(redisURL string) (*Scheduler, error) {
	redisOpt, err := asynq.ParseRedisURI(redisURL)
	_ = redisOpt

	// grit:cron-tasks

	return &Scheduler{scheduler: scheduler}, nil
}
`

func jobProject(t *testing.T, workers string) string {
	t.Helper()
	api := filepath.Join(t.TempDir(), "apps", "api")
	writeTestFile(t, filepath.Join(api, "internal", "jobs", "workers.go"), workers)
	writeTestFile(t, filepath.Join(api, "internal", "cron", "cron.go"), cronWithMarker)
	return api
}

func TestGenerateJobWritesAndRegisters(t *testing.T) {
	api := jobProject(t, workersWithMarker)
	if err := generateJobAt(api, JobOptions{Name: "ReconcileLedger"}); err != nil {
		t.Fatalf("generate: %v", err)
	}

	job := readTestFile(t, filepath.Join(api, "internal", "jobs", "reconcile_ledger.go"))
	for _, want := range []string{
		`const TypeReconcileLedger = "reconcile_ledger"`,
		"func (c *Client) EnqueueReconcileLedger(",
		"func handleReconcileLedger(deps WorkerDeps)",
		"asynq.SkipRetry",
	} {
		if !strings.Contains(job, want) {
			t.Errorf("job file is missing %s", want)
		}
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "job.go", job, parser.AllErrors); err != nil {
		t.Fatalf("the generated job does not parse: %v", err)
	}

	workers := readTestFile(t, filepath.Join(api, "internal", "jobs", "workers.go"))
	if !strings.Contains(workers, "mux.HandleFunc(TypeReconcileLedger, handleReconcileLedger(deps))") {
		t.Error("the handler was never registered, so the job would be queued and never run")
	}
	if strings.Index(workers, "handleReconcileLedger") > strings.Index(workers, "// grit:jobs") {
		t.Error("the registration went after the marker, so the next job would land above it")
	}
	// No --cron, no schedule.
	if strings.Contains(readTestFile(t, filepath.Join(api, "internal", "cron", "cron.go")), "reconcile_ledger") {
		t.Error("a job without --cron was scheduled anyway")
	}
}

func TestGenerateJobSchedulesWithCron(t *testing.T) {
	api := jobProject(t, workersWithMarker)
	if err := generateJobAt(api, JobOptions{Name: "ReconcileLedger", Cron: "30 23 * * *"}); err != nil {
		t.Fatalf("generate: %v", err)
	}
	cron := readTestFile(t, filepath.Join(api, "internal", "cron", "cron.go"))
	for _, want := range []string{
		`scheduler.Register("30 23 * * *", asynq.NewTask("reconcile_ledger", nil))`,
		`Name:     "Reconcile Ledger"`,
		`Type:     "reconcile_ledger"`,
	} {
		if !strings.Contains(cron, want) {
			t.Errorf("cron.go is missing %s", want)
		}
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "cron.go", cron, parser.AllErrors); err != nil {
		t.Fatalf("cron.go no longer parses: %v", err)
	}
}

// A project from before the marker still gets the handler registered, and the
// marker, so the next job has somewhere to go.
func TestGenerateJobOnAnOlderWorker(t *testing.T) {
	api := jobProject(t, workersWithoutMarker)
	if err := generateJobAt(api, JobOptions{Name: "SyncRates"}); err != nil {
		t.Fatalf("generate: %v", err)
	}
	workers := readTestFile(t, filepath.Join(api, "internal", "jobs", "workers.go"))
	if !strings.Contains(workers, "mux.HandleFunc(TypeSyncRates, handleSyncRates(deps))") {
		t.Error("the handler was not registered on a worker without the marker")
	}
	if !strings.Contains(workers, "// grit:jobs") {
		t.Error("the marker was not left behind for the next job")
	}
}

// The job file holds the user's handler, so a second generate refuses.
func TestGenerateJobDoesNotOverwrite(t *testing.T) {
	api := jobProject(t, workersWithMarker)
	if err := generateJobAt(api, JobOptions{Name: "SyncRates"}); err != nil {
		t.Fatalf("first: %v", err)
	}
	path := filepath.Join(api, "internal", "jobs", "sync_rates.go")
	writeTestFile(t, path, readTestFile(t, path)+"\n// my work\n")

	if err := generateJobAt(api, JobOptions{Name: "SyncRates"}); err == nil {
		t.Fatal("a second generate overwrote the handler")
	}
	if !strings.Contains(readTestFile(t, path), "// my work") {
		t.Error("the user's code is gone")
	}
}

func TestGenerateJobRejectsBadSchedules(t *testing.T) {
	for _, spec := range []string{"every day", "30 23 * *", "@fortnightly", "@every soon"} {
		api := jobProject(t, workersWithMarker)
		if err := generateJobAt(api, JobOptions{Name: "Nightly", Cron: spec}); err == nil {
			t.Errorf("--cron %q was accepted", spec)
		}
	}
	for _, spec := range []string{"30 23 * * *", "@daily", "@every 15m", "0 */6 * * 1-5"} {
		if err := validateCronSpec(spec); err != nil {
			t.Errorf("--cron %q was refused: %v", spec, err)
		}
	}
}
