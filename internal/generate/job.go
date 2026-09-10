package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// JobOptions describes one `grit generate job` request.
type JobOptions struct {
	Name string
	// Cron, when set, also runs the job on a schedule: a five-field spec such
	// as "30 23 * * *", or a descriptor such as @daily or "@every 15m".
	Cron string
}

// GenerateJob writes a background job and registers it with the worker, and
// with the cron scheduler when a schedule is given.
//
// Found building end-of-day reconciliation for a ledger. The cron docs said new
// scheduled tasks are injected at the grit:cron-tasks marker "when you use
// grit add cron", and there was no such command. A custom job was five hand
// edits across two framework files, one of which, the worker's handler list,
// had no marker to anchor on. A job enqueued without its handler registered is
// accepted by Redis and then fails on the worker, far from the line that
// caused it.
func GenerateJob(opts JobOptions) error {
	root, err := findProjectRoot()
	if err != nil {
		return err
	}
	arch, _ := readGritJSON(root)
	apiRoot := root
	if arch != "single" {
		apiRoot = filepath.Join(root, "apps", "api")
	}
	return generateJobAt(apiRoot, opts)
}

func generateJobAt(apiRoot string, opts JobOptions) error {
	if opts.Name == "" {
		return fmt.Errorf("a job name is required (e.g. grit generate job ReconcileLedger)")
	}
	if opts.Cron != "" {
		if err := validateCronSpec(opts.Cron); err != nil {
			return err
		}
	}

	names := MakeNames(opts.Name)
	taskType := names.Snake

	jobPath := filepath.Join(apiRoot, "internal", "jobs", names.Snake+".go")
	workersPath := filepath.Join(apiRoot, "internal", "jobs", "workers.go")
	cronPath := filepath.Join(apiRoot, "internal", "cron", "cron.go")

	if _, err := os.Stat(jobPath); err == nil {
		return fmt.Errorf("internal/jobs/%s.go already exists. It holds your handler, so it is not overwritten", names.Snake)
	}
	if _, err := os.Stat(workersPath); err != nil {
		return fmt.Errorf("internal/jobs/workers.go not found: this project has no background worker to register the job with")
	}
	if opts.Cron != "" {
		if _, err := os.Stat(cronPath); err != nil {
			return fmt.Errorf("--cron needs internal/cron/cron.go, which this project does not have")
		}
	}

	fmt.Printf("\n  Generating job: %s (task type %q)\n\n", names.Pascal, taskType)

	if err := writeFileWithDirs(jobPath, jobFileGo(names, taskType)); err != nil {
		return fmt.Errorf("writing the job: %w", err)
	}
	fmt.Printf("  ✓ internal/jobs/%s.go\n", names.Snake)

	if err := registerJobHandler(workersPath, names); err != nil {
		return err
	}
	fmt.Println("  ✓ Registered the handler with the worker")

	if opts.Cron != "" {
		if err := registerCronTask(cronPath, names, taskType, opts.Cron); err != nil {
			return err
		}
		fmt.Printf("  ✓ Scheduled: %s (listed in the admin under Cron)\n", opts.Cron)
	}

	fmt.Printf("\n  ✅ Job generated. Write the work in handle%s, then queue it:\n\n", names.Pascal)
	fmt.Printf("      jobsClient.Enqueue%s(ctx, jobs.%sPayload{}, jobs.EnqueueOption{\n", names.Pascal, names.Pascal)
	fmt.Printf("          IdempotencyKey: \"%s:\" + day, // the same key twice is one run, not two\n", taskType)
	fmt.Printf("      })\n\n")
	return nil
}

// validateCronSpec catches the mistakes worth catching before they reach a
// scheduler that would only log them at startup.
func validateCronSpec(spec string) error {
	s := strings.TrimSpace(spec)
	if strings.HasPrefix(s, "@") {
		switch s {
		case "@yearly", "@annually", "@monthly", "@weekly", "@daily", "@midnight", "@hourly":
			return nil
		}
		if every := strings.TrimPrefix(s, "@every "); every != s {
			if _, err := time.ParseDuration(strings.TrimSpace(every)); err == nil {
				return nil
			}
		}
		return fmt.Errorf("--cron %q is not a schedule asynq understands: use @daily, @hourly, \"@every 15m\" or five fields", spec)
	}
	if len(strings.Fields(s)) != 5 {
		return fmt.Errorf("--cron wants five fields (minute hour day-of-month month day-of-week), e.g. \"30 23 * * *\" for 23:30 every day; got %q", spec)
	}
	return nil
}

// registerJobHandler adds the handler to the worker's ServeMux.
//
// At the grit:jobs marker when the project has one. A project scaffolded before
// the marker gets the line just before the worker starts, plus the marker, so
// the next job has somewhere to go.
func registerJobHandler(path string, names Names) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	crlf := strings.Contains(string(raw), "\r\n")
	s := strings.ReplaceAll(string(raw), "\r\n", "\n")

	line := fmt.Sprintf("\tmux.HandleFunc(Type%s, handle%s(deps))", names.Pascal, names.Pascal)
	if strings.Contains(s, line) {
		return nil
	}
	const marker = "\t// grit:jobs\n"
	const start = "\tgo func() {\n\t\tif err := srv.Run(mux); err != nil {"
	switch {
	case strings.Contains(s, marker):
		s = strings.Replace(s, marker, line+"\n"+marker, 1)
	case strings.Contains(s, start):
		s = strings.Replace(s, start, line+"\n"+marker+"\n"+start, 1)
	default:
		return fmt.Errorf("could not find where handlers are registered in internal/jobs/workers.go.\n\n"+
			"Add this after the other mux.HandleFunc lines, or the job is queued and never run:\n\n  %s", strings.TrimSpace(line))
	}
	if crlf {
		s = strings.ReplaceAll(s, "\n", "\r\n")
	}
	return os.WriteFile(path, []byte(s), 0644)
}

// registerCronTask schedules the job at the grit:cron-tasks marker, and lists
// it in RegisteredTasks so the admin's cron page shows it.
func registerCronTask(path string, names Names, taskType, spec string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	crlf := strings.Contains(string(raw), "\r\n")
	s := strings.ReplaceAll(string(raw), "\r\n", "\n")

	if strings.Contains(s, fmt.Sprintf("asynq.NewTask(%q, nil)", taskType)) {
		return nil
	}
	const marker = "\t// grit:cron-tasks"
	if !strings.Contains(s, marker) {
		return fmt.Errorf("internal/cron/cron.go has no // grit:cron-tasks marker, so the schedule was not added.\n\n"+
			"Register it by hand in cron.New:\n\n  scheduler.Register(%q, asynq.NewTask(%q, nil))", spec, taskType)
	}
	human := strings.Join(splitPascal(names.Pascal), " ")
	block := fmt.Sprintf("\t_, err = scheduler.Register(%q, asynq.NewTask(%q, nil))\n"+
		"\tif err != nil {\n"+
		"\t\treturn nil, fmt.Errorf(\"registering %s: %%w\", err)\n"+
		"\t}\n"+
		"\tRegisteredTasks = append(RegisteredTasks, Task{\n"+
		"\t\tName:     %q,\n"+
		"\t\tSchedule: %q,\n"+
		"\t\tType:     %q,\n"+
		"\t})\n\n", spec, taskType, taskType, human, spec, taskType)
	s = strings.Replace(s, marker, block+marker, 1)
	if crlf {
		s = strings.ReplaceAll(s, "\n", "\r\n")
	}
	return os.WriteFile(path, []byte(s), 0644)
}

// jobFileGo is the job's own file. It is the user's from the moment it is
// written, which is why a second generate refuses rather than overwrites.
func jobFileGo(names Names, taskType string) string {
	return strings.NewReplacer(
		"{{P}}", names.Pascal,
		"{{T}}", taskType,
		"{{BT}}", "`",
	).Replace(`package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

// Type{{P}} is the task type the worker routes on.
const Type{{P}} = "{{T}}"

// {{P}}Payload is what one run needs. Keep it to IDs and small values: it is
// serialised into Redis, and a retry sees exactly what the first attempt saw.
// A scheduled run arrives with an empty payload.
type {{P}}Payload struct {
	// Day string {{BT}}json:"day,omitempty"{{BT}}
}

// Enqueue{{P}} queues one run. Pass an IdempotencyKey whenever the same piece of
// work might be queued twice, a nightly run for one date say: the second
// enqueue is refused rather than run again.
func (c *Client) Enqueue{{P}}(ctx context.Context, p {{P}}Payload, opts ...EnqueueOption) error {
	return c.Enqueue(ctx, Type{{P}}, p, opts...)
}

// handle{{P}} runs one task. Return an error to be retried with backoff. Return
// nil once the work is done, or when it can never succeed, since retrying only
// repeats the failure.
func handle{{P}}(deps WorkerDeps) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		var p {{P}}Payload
		if len(task.Payload()) > 0 {
			if err := json.Unmarshal(task.Payload(), &p); err != nil {
				// A payload that does not decode now never will.
				return fmt.Errorf("decoding a %s payload: %v: %w", Type{{P}}, err, asynq.SkipRetry)
			}
		}

		// The work goes here. deps.DB, deps.Mailer, deps.Storage and
		// deps.Cache are available.
		log.Printf("%s: ran", Type{{P}})
		_ = deps
		return nil
	}
}
`)
}
