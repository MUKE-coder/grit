package scaffold

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The files an existing project has: today's templates with the retention
// sweep taken out. Built from the templates rather than written by hand, so the
// anchors are tested against what Grit really writes.
func preRetentionClient() string {
	src := strings.ReplaceAll(jobsClientGo(), "{{MODULE}}", "app/apps/api")
	start, after, ok := lineWith(src, "TypeRetentionSweep ")
	if !ok {
		panic("the client template has no TypeRetentionSweep to take out")
	}
	return src[:start] + src[after:]
}

func preRetentionWorkers() string {
	src := strings.ReplaceAll(jobsWorkersGo(), "{{MODULE}}", "app/apps/api")
	src = strings.Replace(src, "\tmux.HandleFunc(TypeRetentionSweep, handleRetentionSweep(deps))\n", "", 1)
	return strings.Replace(src, retentionJobHandler+"\n\n", "", 1)
}

func preRetentionCron() string {
	return strings.Replace(cronSchedulerGo(), retentionCronRegistration, "", 1)
}

func TestThePreRetentionFixturesHaveNoRetention(t *testing.T) {
	for name, src := range map[string]string{
		"client": preRetentionClient(), "workers": preRetentionWorkers(), "cron": preRetentionCron(),
	} {
		if strings.Contains(src, "retention:sweep") || strings.Contains(src, "RetentionSweep") {
			t.Errorf("the %s fixture still has the retention sweep in it, so the repair tests prove nothing", name)
		}
	}
}

// Each repair wires a project scaffolded before the sweep, produces valid Go,
// and changes nothing the second time, which is what makes running upgrade
// twice safe.
func TestTheRetentionRepairsWireAnOlderProject(t *testing.T) {
	cases := []struct {
		name string
		fn   func(string) (string, []string, []string)
		src  string
		want []string
	}{
		{"task type", repairRetentionTaskTypeSource, preRetentionClient(),
			[]string{`TypeRetentionSweep = "retention:sweep"`}},
		{"worker", repairRetentionWorkerSource, preRetentionWorkers(),
			[]string{"mux.HandleFunc(TypeRetentionSweep, handleRetentionSweep(deps))", "func handleRetentionSweep("}},
		{"cron", repairRetentionCronSource, preRetentionCron(),
			[]string{`asynq.NewTask("retention:sweep", nil)`, `Type:     "retention:sweep"`}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			once, fixed, warn := c.fn(c.src)
			if len(warn) > 0 {
				t.Fatalf("warned on the shape Grit writes: %v", warn)
			}
			if len(fixed) == 0 {
				t.Fatal("the repair did nothing to a project that needs it")
			}
			for _, w := range c.want {
				if !strings.Contains(once, w) {
					t.Errorf("the repaired file is missing %q", w)
				}
			}
			if _, err := format.Source([]byte(once)); err != nil {
				t.Errorf("the repaired file is not valid Go: %v", err)
			}
			twice, fixed2, warn2 := c.fn(once)
			if twice != once || len(fixed2) > 0 || len(warn2) > 0 {
				t.Error("running upgrade twice changes the file twice")
			}
		})
	}
}

// A project already on the sweep is left exactly as it is.
func TestTheRetentionRepairsLeaveANewProjectAlone(t *testing.T) {
	for name, c := range map[string]struct {
		fn  func(string) (string, []string, []string)
		src string
	}{
		"task type": {repairRetentionTaskTypeSource, jobsClientGo()},
		"worker":    {repairRetentionWorkerSource, jobsWorkersGo()},
		"cron":      {repairRetentionCronSource, cronSchedulerGo()},
	} {
		if out, fixed, warn := c.fn(c.src); out != c.src || len(fixed) > 0 || len(warn) > 0 {
			t.Errorf("%s: a file Grit wrote with the sweep in it was changed (%v, %v)", name, fixed, warn)
		}
	}
}

// A project older than the Trash sweep has no Trash job to anchor on, and
// still has everything retention needs. examples/library is one, and the first
// version of this repair declined it.
func TestTheRetentionRepairsWireAProjectOlderThanTheTrash(t *testing.T) {
	client := strings.Replace(preRetentionClient(), "TypeTrashPurgeExpired ", "TypeUnrelatedThing ", 1)
	workers := strings.Replace(preRetentionWorkers(),
		"\tmux.HandleFunc(TypeTrashPurgeExpired, handleTrashPurgeExpired(deps))\n", "", 1)

	for name, c := range map[string]struct {
		fn   func(string) (string, []string, []string)
		src  string
		want string
	}{
		"task type": {repairRetentionTaskTypeSource, client, `TypeRetentionSweep = "retention:sweep"`},
		"worker":    {repairRetentionWorkerSource, workers, "mux.HandleFunc(TypeRetentionSweep, handleRetentionSweep(deps))\n\t// grit:jobs"},
	} {
		out, fixed, warn := c.fn(c.src)
		if len(warn) > 0 || len(fixed) == 0 {
			t.Errorf("%s: declined a project older than the Trash (%v)", name, warn)
			continue
		}
		if !strings.Contains(out, c.want) {
			t.Errorf("%s: not wired where expected, missing %q", name, c.want)
		}
		if _, err := format.Source([]byte(out)); err != nil {
			t.Errorf("%s: the repaired file is not valid Go: %v", name, err)
		}
	}
}

// A worker the repair cannot wire completely is left as it was: a handler with
// no registration is dead code, and a registration with no handler does not
// compile.
func TestTheRetentionWorkerRepairIsAllOrNothing(t *testing.T) {
	src := strings.Replace(preRetentionWorkers(),
		"\tmux.HandleFunc(TypeTrashPurgeExpired, handleTrashPurgeExpired(deps))\n", "", 1)
	src = strings.Replace(src, "\t// grit:jobs\n", "", 1)

	out, fixed, warn := repairRetentionWorkerSource(src)
	if len(warn) == 0 {
		t.Error("a mux with neither anchor was repaired without a word")
	}
	if len(fixed) > 0 || out != src {
		t.Error("the file was changed although it could not be wired completely")
	}
}

// Each file is only touched once the one before it is wired, and the worker
// only where what the handler calls is there: the handler names the task type
// and the sweep, and a schedule with no handler fails every night.
func TestTheRetentionSweepIsOnlyScheduledOnceItCanRun(t *testing.T) {
	setup := func(t *testing.T, client string, withSweep bool) (root string, read func(string) string) {
		root = t.TempDir()
		api := filepath.Join(root, "apps", "api", "internal")
		files := map[string]string{
			filepath.Join(api, "jobs", "client.go"):   client,
			filepath.Join(api, "jobs", "workers.go"):  preRetentionWorkers(),
			filepath.Join(api, "cron", "cron.go"):     preRetentionCron(),
			filepath.Join(api, "sync", "registry.go"): "package sync\n\ntype Registry struct{}\n\nvar Shared *Registry\n",
		}
		if withSweep {
			files[filepath.Join(api, "services", "retention.go")] = "package services\n"
		}
		for path, src := range files {
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return root, func(rel string) string {
			b, err := os.ReadFile(filepath.Join(api, rel))
			if err != nil {
				t.Fatal(err)
			}
			return string(b)
		}
	}
	unwired := func(t *testing.T, read func(string) string) {
		t.Helper()
		if read("jobs/workers.go") != preRetentionWorkers() {
			t.Error("the worker was wired to something the project does not have, so it does not compile")
		}
		if read("cron/cron.go") != preRetentionCron() {
			t.Error("the sweep was scheduled with nothing to handle it")
		}
	}

	t.Run("the task type cannot be added", func(t *testing.T) {
		own := strings.Replace(preRetentionClient(), "TypeTrashPurgeExpired ", "TypeBinPurge ", 1)
		own = strings.Replace(own, "TypeAuditPrune ", "TypeLogPrune ", 1)
		root, read := setup(t, own, true)
		if err := repairRetentionSweep(root, Options{ProjectName: "app"}); err != nil {
			t.Fatal(err)
		}
		unwired(t, read)
	})

	t.Run("the sweep is not there", func(t *testing.T) {
		root, read := setup(t, preRetentionClient(), false)
		if err := repairRetentionSweep(root, Options{ProjectName: "app"}); err != nil {
			t.Fatal(err)
		}
		unwired(t, read)
	})

	t.Run("an older project", func(t *testing.T) {
		root, read := setup(t, preRetentionClient(), true)
		if err := repairRetentionSweep(root, Options{ProjectName: "app"}); err != nil {
			t.Fatal(err)
		}
		for rel, want := range map[string]string{
			"jobs/client.go":  "TypeRetentionSweep",
			"jobs/workers.go": "func handleRetentionSweep(",
			"cron/cron.go":    `"retention:sweep"`,
		} {
			if !strings.Contains(read(rel), want) {
				t.Errorf("%s was not wired: missing %s", rel, want)
			}
		}
	})
}
