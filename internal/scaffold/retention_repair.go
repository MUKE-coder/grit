package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The retention sweep (grit#260), wired into a project that already exists.
//
// The registry and the sweep are new files and arrive through
// writeFrameworkOwnedFiles. What needs a repair is what schedules them: the
// task type, the worker's handler and the cron entry, in three files a project
// may have made its own. The same terms as the ban sweep's repair: anchored on
// what Grit wrote, a no-op the second time, and a file that cannot be wired
// completely is left as it was and named, with the edit to make by hand.

func repairRetentionSweep(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	jobs := filepath.Join(apiRoot, "internal", "jobs")
	client := filepath.Join(jobs, "client.go")
	workers := filepath.Join(jobs, "workers.go")
	cron := filepath.Join(apiRoot, "internal", "cron", "cron.go")

	// In order, each step only where the one before it landed. The handler
	// names the task type, so wiring it without the constant is "undefined:
	// TypeRetentionSweep" and a project that does not compile. Scheduling a
	// task nothing handles compiles, and then fails every night with asynq's
	// "handler not found" while the policy claims rows are going.
	steps := []struct {
		path  string
		fn    func(string) (string, []string, []string)
		needs func() bool
	}{
		{client, repairRetentionTaskTypeSource, func() bool { return true }},
		{workers, repairRetentionWorkerSource, func() bool { return fileContains(client, "TypeRetentionSweep") }},
		{cron, repairRetentionCronSource, func() bool { return fileContains(workers, "func handleRetentionSweep") }},
	}
	for _, s := range steps {
		if !fileExists(s.path) || !s.needs() {
			continue
		}
		if err := repairSourceFile(root, m, s.path, s.fn); err != nil {
			return fmt.Errorf("%s: %w", s.path, err)
		}
	}
	return nil
}

// repairRetentionTaskTypeSource adds the constant beside the Trash sweep's, in
// the file that declares the task types.
func repairRetentionTaskTypeSource(src string) (string, []string, []string) {
	if strings.Contains(src, "TypeRetentionSweep") {
		return src, nil, nil
	}
	_, after, ok := lineWith(src, "TypeTrashPurgeExpired ")
	if !ok {
		return src, nil, []string{"its task-type list is not the one Grit wrote, so retention:sweep " +
			"has no constant: add TypeRetentionSweep = \"retention:sweep\""}
	}
	src = src[:after] + "\tTypeRetentionSweep = \"retention:sweep\"\n" + src[after:]
	return src, []string{"retention:sweep has a task type"}, nil
}

// repairRetentionWorkerSource registers and adds handleRetentionSweep, or
// changes nothing.
func repairRetentionWorkerSource(src string) (string, []string, []string) {
	// Either half means a run already happened. Guarding on only one would let
	// a file that has the other gain a second copy on every upgrade.
	if strings.Contains(src, "mux.HandleFunc(TypeRetentionSweep") || strings.Contains(src, "func handleRetentionSweep") {
		return src, nil, nil
	}
	// The sweep is the Trash's purge, and the handler reads the registry the
	// Trash's handler reads. A worker without that handler predates both.
	if !strings.Contains(src, "func handleTrashPurgeExpired") {
		return src, nil, []string{"its worker has no Trash sweep, which the retention sweep is built on, " +
			"so retention:sweep was not wired: add handleRetentionSweep from a new project's internal/jobs/workers.go"}
	}

	out := src
	var warn []string
	if _, after, ok := lineWith(out, "mux.HandleFunc(TypeTrashPurgeExpired"); ok {
		out = out[:after] + "\tmux.HandleFunc(TypeRetentionSweep, handleRetentionSweep(deps))\n" + out[after:]
	} else {
		warn = append(warn, "its worker mux is not the one Grit wrote, so retention:sweep has no handler: "+
			"add mux.HandleFunc(TypeRetentionSweep, handleRetentionSweep(deps))")
	}
	// Both are there already wherever the Trash handler is, which is checked
	// above. Asked anyway, because the cost of being wrong is "undefined:
	// services" in a file the user did not touch.
	for _, pkg := range []string{"services", "sync"} {
		added, ok := withInternalImport(out, pkg)
		if !ok {
			warn = append(warn, "its import block is not the one Grit wrote, so internal/"+pkg+
				" could not be imported for the retention sweep")
		}
		out = added
	}
	if len(warn) > 0 {
		return src, nil, warn
	}
	// At the end of the file: a top-level function can go anywhere, and the end
	// is the one place no project has moved.
	out = strings.TrimRight(out, "\n") + "\n\n" + retentionJobHandler + "\n"
	return out, []string{"rows past their retention window are removed nightly"}, nil
}

// repairRetentionCronSource schedules the sweep at the marker every generated
// cron.go has.
func repairRetentionCronSource(src string) (string, []string, []string) {
	if strings.Contains(src, `"retention:sweep"`) {
		return src, nil, nil
	}
	const marker = "\t// grit:cron-tasks\n"
	at := strings.Index(src, marker)
	if at < 0 {
		return src, nil, []string{"its cron.go has no grit:cron-tasks marker, so the retention sweep is not " +
			"scheduled: add a scheduler.Register for retention:sweep"}
	}
	src = src[:at] + retentionCronRegistration + src[at:]
	return src, []string{"the retention sweep runs daily at 03:00"}, nil
}

// retentionCronRegistration is the same entry a new project's cron.go has.
const retentionCronRegistration = `	// Per-resource retention (grit#260), daily at 03:00: rows of every resource
	// generated with --retain are removed once older than its window. Before
	// the orphan sweep at 03:15 on purpose: the files those rows pointed at are
	// released to it, so they go the same night rather than a day later.
	_, err = scheduler.Register("0 3 * * *", asynq.NewTask("retention:sweep", nil), asynq.MaxRetry(3), asynq.Timeout(2*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("registering retention sweep: %w", err)
	}
	RegisteredTasks = append(RegisteredTasks, Task{
		Name:     "Remove rows past their retention window",
		Schedule: "0 3 * * *",
		Type:     "retention:sweep",
	})

`

// retentionJobHandler is handleRetentionSweep in jobs/workers.go. One copy,
// spliced into the template too, so a new project and an upgraded one run the
// same code.
const retentionJobHandler = `// handleRetentionSweep removes the rows of every resource generated with
// --retain once they are older than its window (grit#260). The removing is
// services.RetentionService, which goes through the Trash's purge so the files
// are released to the orphan sweep rather than deleted twice over.
func handleRetentionSweep(deps WorkerDeps) func(ctx context.Context, task *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		if deps.DB == nil {
			return fmt.Errorf("database not configured")
		}
		if sync.Shared == nil {
			// The same as the Trash sweep: nothing has built a registry, so no
			// table can be read. Tomorrow's run gets another chance.
			log.Println("Retention sweep: no resource registry yet, nothing to do")
			return nil
		}
		svc := &services.RetentionService{DB: deps.DB, Registry: sync.Shared}
		removed, err := svc.Sweep(ctx)
		// Logged before the error is returned: one failing resource does not
		// stop the others, and what did go is worth a line either way.
		for table, n := range removed {
			log.Printf("Retention sweep: removed %d %s row(s) past their window", n, table)
		}
		if err != nil {
			return fmt.Errorf("retention sweep: %w", err)
		}
		if len(removed) == 0 {
			log.Println("Retention sweep: nothing past its window")
		}
		return nil
	}
}`
