package scaffold

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Behind a load balancer, a role revoked through one replica kept working on
// the others until they restarted: the grants cache was per process.
// Reproduced live with two instances over one database.
func TestPermissionCacheIsSharedAcrossReplicas(t *testing.T) {
	src := authzGrantsGo()
	for _, want := range []string{
		"func Share(db *gorm.DB)",
		`cluster.Bump(sharedDB, "authz")`,
		"sharedWatch.Changed()",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("authz grants.go is missing %s", want)
		}
	}
	routes := apiRoutesGo()
	if !strings.Contains(routes, "authz.Share(db)") {
		t.Error("the routes template does not share the permission cache")
	}
}

// An SSO connection created on one replica was unknown to the others.
func TestSSORegistriesFollowOtherReplicas(t *testing.T) {
	if !strings.Contains(apiSSOServiceGo(), "func WatchSSO(") {
		t.Error("nothing rebuilds the SSO registries when another replica changes them")
	}
	if !strings.Contains(apiSSOHandlerGo(), "services.AnnounceSSOChange(h.DB)") {
		t.Error("changing a connection does not tell the other replicas")
	}
	if !strings.Contains(apiRoutesGo(), "services.WatchSSO(db, ssoRegistry, samlRegistry)") {
		t.Error("the routes template does not start the SSO watch")
	}
}

// Every replica ran every scheduled job.
func TestCronRunsOnOneReplica(t *testing.T) {
	src := cronSchedulerGo()
	if !strings.Contains(src, "go s.lead(s.redisURL)") || strings.Contains(src, "s.scheduler.Run()") {
		t.Error("the scheduler still starts on every replica")
	}
	leader := cronLeaderGo()
	for _, want := range []string{"SetNX(ctx, leaderKey, id, leaderLease)", "renewLease.Run(", "stopOnce.Do(s.scheduler.Shutdown)"} {
		if !strings.Contains(leader, want) {
			t.Errorf("leader.go is missing %s", want)
		}
	}
}

// An existing cron.go, which grit add job writes into, is patched rather than
// replaced.
func TestRepairCronLeaderPatchesTheGeneratedLines(t *testing.T) {
	old := "package cron\n\nimport (\n\t\"fmt\"\n\t\"log\"\n\n\t\"github.com/hibiken/asynq\"\n)\n\n"
	for _, p := range cronPatches {
		old += p.old + "\n\n"
	}
	old += "func New(redisURL string) (*Scheduler, error) {\n\tscheduler := asynq.NewScheduler(nil, nil)\n\t_ = fmt.Sprint()\n" + cronPatches[1].old + "\n}\n"
	// The return line is patched where New returns; keep only that one copy.
	old = strings.Replace(old, cronPatches[1].old+"\n\n", "", 1)

	out, fixed, warn := repairCronSource(old)
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("the repaired cron.go is not valid Go: %v\n%s", err, out)
	}
	if !strings.Contains(out, "go s.lead(s.redisURL)") || strings.Contains(out, "\t\"log\"\n") {
		t.Errorf("not repaired:\n%s", out)
	}
	if again, fixed, _ := repairCronSource(out); again != out || len(fixed) != 0 {
		t.Error("a second upgrade changed cron.go again")
	}
}

func TestClusterWiringIsGatedAndIdempotent(t *testing.T) {
	api := t.TempDir()
	put := func(rel, body string) {
		t.Helper()
		p := filepath.Join(api, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	routes := "package routes\n\nimport (\n\t\"x/internal/handlers\"\n)\n\nfunc setup() {\n" +
		"\troleHandler := handlers.NewRoleHandler(db)\n" +
		"\tssoHandler := handlers.NewSSOHandler(db, authService, cfg, ssoRegistry, samlRegistry)\n}\n"
	put("internal/routes/routes.go", routes)

	// Neither function exists yet: nothing may be injected.
	if err := ensureClusterWiring(api, "x"); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(filepath.Join(api, "internal", "routes", "routes.go"))
	if string(got) != routes {
		t.Fatal("routes.go was wired to functions the project does not have")
	}

	put("internal/authz/grants.go", "func Share(db *gorm.DB) {}")
	put("internal/services/sso.go", "func WatchSSO(db *gorm.DB) {}")
	for i := 0; i < 2; i++ {
		if err := ensureClusterWiring(api, "x"); err != nil {
			t.Fatal(err)
		}
	}
	got, _ = os.ReadFile(filepath.Join(api, "internal", "routes", "routes.go"))
	s := string(got)
	if strings.Count(s, "authz.Share(db)") != 1 || strings.Count(s, "services.WatchSSO(") != 1 {
		t.Errorf("wired %d and %d times, want once each:\n%s",
			strings.Count(s, "authz.Share(db)"), strings.Count(s, "services.WatchSSO("), s)
	}
	if !strings.Contains(s, `"x/internal/authz"`) {
		t.Error("authz is called without being imported")
	}
}
