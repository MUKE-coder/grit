package scaffold

import (
	"strings"
	"testing"
)

// v3349Routes is routes.go as v3.349.0 wrote it: the health handler that
// reported one boolean per component, before the four states replaced it.
//
// Built from oldHealthCheckBlock rather than by undoing today's template, which
// is what these two tests did until v3.350.0. The input was rebuilt from the
// current template, so when the template moved the input moved with it and the
// tests described a transition nobody was making. oldHealthCheckBlock is fixed
// text, taken out of a project a release generated.
func v3349Routes(t *testing.T) string {
	t.Helper()
	fresh := apiRoutesGo()
	old := strings.Replace(fresh, healthCheckBlock, oldHealthCheckBlock, 1)
	old = strings.Replace(old, eventBusStatusBlock, oldEventBusStatus, 1)
	old = strings.Replace(old, "\t\"{{MODULE}}/internal/health\"\n", "", 1)
	if old == fresh || !strings.Contains(old, "type compStatus struct {") ||
		strings.Contains(old, "/internal/health\"") || strings.Contains(old, "health.Overall(") {
		t.Fatal("could not rebuild routes.go as v3.349.0 wrote it")
	}
	return old
}

func oldHealthRoutes(t *testing.T) string {
	t.Helper()
	fresh := v3349Routes(t)
	old := strings.Replace(fresh, healthJobsBlock+"\n", oldHealthJobsBlock, 1)
	old = strings.Replace(old, healthQueueFields, oldHealthQueueField, 1)
	old = strings.Replace(old, healthQueueStatsSetup, "", 1)
	if old == fresh || !strings.Contains(old, "redis.call('keys', 'asynq:*')") {
		t.Fatal("could not rebuild routes.go from before the repair")
	}
	return old
}

func TestRepairHealthQueueProbe(t *testing.T) {
	out, fixed, warn := repairHealthQueueProbeSource(oldHealthRoutes(t))
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if out != v3349Routes(t) {
		t.Error("the repaired routes.go differs from the one that release wrote")
	}
	if again, fixed, _ := repairHealthQueueProbeSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed routes.go again")
	}
}

func TestRepairHealthQueueProbeWarnsOnEditedCode(t *testing.T) {
	src := strings.Replace(oldHealthRoutes(t), "jobsStatus.QueueKeys = n", "jobsStatus.QueueKeys = n + 0", 1)
	if out, fixed, warn := repairHealthQueueProbeSource(src); out != src || len(fixed) > 0 || len(warn) != 1 {
		t.Errorf("fixed %v, warned %v", fixed, warn)
	}
}

func TestFreshTemplatesNeedNoHealthQueueRepair(t *testing.T) {
	src := v3349Routes(t)
	if strings.Contains(src, "redis.call('keys'") || strings.Contains(src, "queue_keys") {
		t.Error("/api/health still counts asynq keys inside Redis")
	}
	if out, fixed, warn := repairHealthQueueProbeSource(src); out != src || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("routes.go still needs the repair (fixed %v, warn %v)", fixed, warn)
	}
	if page := adminSystemHealthPage(); strings.Contains(page, "queue_keys") || !strings.Contains(page, "data.jobs.queued") {
		t.Error("the admin health page still reads queue_keys")
	}
}
