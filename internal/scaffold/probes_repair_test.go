package scaffold

import (
	"strings"
	"testing"
)

// The repair's output is the template, exactly.
//
// Two things that have to agree: the routes.go a new project gets, and the
// routes.go an upgraded project is repaired into. A test over either alone
// passes while the pair is broken, and the way it breaks is quiet: a new
// project serves /readyz, an upgraded one answers 404, and nothing says so
// because a file with a function nobody calls still compiles.
func TestTheProbesRepairProducesTheTemplate(t *testing.T) {
	fresh := apiRoutesGo()
	if !strings.Contains(fresh, "mountProbes(r, db)") {
		t.Fatal("the routes template does not mount the probes, so new projects " +
			"do not serve /healthz or /readyz")
	}

	// routes.go as it was before either edit: no mount, and the probes still
	// in Pulse's request log. Both, because the repair makes both and a test
	// that undoes one covers one.
	before := strings.Replace(fresh, probesMountNew, probesMountOld, 1)
	before = strings.Replace(before,
		`pulse.WithExcludePaths("/studio/*", "/sentinel/*", "/docs/*", "/pulse/*",
				"/healthz", "/readyz"),`,
		`pulse.WithExcludePaths("/studio/*", "/sentinel/*", "/docs/*", "/pulse/*"),`, 1)
	if before == fresh || strings.Contains(before, "mountProbes(r, db)") ||
		strings.Contains(before, `"/healthz", "/readyz"),`) {
		t.Fatal("could not rebuild routes.go from before the two edits")
	}

	out, fixed, warn := repairProbesMountSource(before)
	if len(warn) > 0 || len(fixed) != 2 {
		t.Fatalf("fixed %v, warned %v: expected the mount and the Pulse exclusion", fixed, warn)
	}
	if out != fresh {
		t.Error("the repaired file is not the template. The two have to be the " +
			"same bytes, or a new project and an upgraded one serve different " +
			"routes from the same version.")
	}
}

// Running it twice changes nothing.
//
// Every repair runs on every upgrade, so the second pass over an already
// repaired project is the normal case, not an edge one.
func TestTheProbesRepairIsIdempotent(t *testing.T) {
	fresh := apiRoutesGo()
	out, fixed, warn := repairProbesMountSource(fresh)
	if out != fresh || len(fixed) != 0 || len(warn) != 0 {
		t.Errorf("a second pass changed the file: fixed %v, warned %v", fixed, warn)
	}
}

// The guard is on the call and not on the prose beside it.
//
// A Contains check that matches the comment explaining the mount is a repair
// that no-ops for ever and reports success, which has happened here before.
func TestTheProbesGuardIsNotItsOwnComment(t *testing.T) {
	// The comment alone, with the call taken out: the repair must still act.
	fresh := apiRoutesGo()
	commentOnly := strings.Replace(fresh, "\tmountProbes(r, db)\n", "", 1)
	if commentOnly == fresh {
		t.Fatal("could not remove the call")
	}
	if !strings.Contains(commentOnly, "internal/routes/probes_routes.go") {
		t.Fatal("the comment went with the call, so this case cannot arise")
	}

	_, fixed, _ := repairProbesMountSource(commentOnly)
	if len(fixed) != 1 {
		t.Error("the repair skipped a file that has the comment and not the call")
	}
}

// A project with no recognisable anchor is told, rather than left silent.
func TestTheProbesRepairWarnsWhenItCannotFindTheAnchor(t *testing.T) {
	out, fixed, warn := repairProbesMountSource("package routes\n\nfunc SetupRoutes() {}\n")
	if len(fixed) != 0 {
		t.Errorf("claimed to fix a file it did not understand: %v", fixed)
	}
	if len(warn) != 1 {
		t.Fatalf("said nothing about a file it could not repair: %v", warn)
	}
	if !strings.Contains(warn[0], "mountProbes(r, db)") {
		t.Errorf("the warning does not say what to add by hand: %q", warn[0])
	}
	if strings.Contains(out, "mountProbes") {
		t.Error("it edited the file anyway")
	}
}

// Both probe files travel on upgrade.
//
// The file is framework-owned on purpose: a readiness rule from before anybody
// thought about which dependencies should gate rotation is exactly the kind of
// thing that should not stay behind in an old project.
func TestTheProbeFilesAreFrameworkOwned(t *testing.T) {
	body := read(t, "api_framework_owned_files.go")
	for _, want := range []string{
		`"internal", "routes", "probes_routes.go"`,
		`"internal", "routes", "probes_routes_test.go"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("%s is not in writeFrameworkOwnedFiles, so an upgraded "+
				"project never gets it", want)
		}
	}
}

// Liveness probes nothing.
//
// The response to a failed liveness probe is to kill the container, so a
// liveness check that pings Redis turns a Redis blip into a rolling restart of
// every replica: a degraded service becomes an outage. This is the one rule in
// that file worth a guard, because the tempting edit is to make /healthz
// "more useful".
func TestLivenessTouchesNoDependency(t *testing.T) {
	body := tmpl("api/routes/probes_routes.go")

	at := strings.Index(body, `r.GET("/healthz"`)
	if at < 0 {
		t.Fatal("no /healthz route")
	}
	end := strings.Index(body[at:], `r.GET("/readyz"`)
	if end < 0 {
		t.Fatal("no /readyz route")
	}
	liveness := body[at : at+end]

	for _, dependency := range []string{"Ping", "db.", "sqlDB", "Snapshot()", "Cache", "Storage"} {
		if strings.Contains(liveness, dependency) {
			t.Errorf("the liveness handler touches %q. A failed liveness probe "+
				"kills the container, so probing a dependency there restarts "+
				"every replica when that dependency blips.", dependency)
		}
	}
}

// A new project does not log its own probes.
//
// Not a nicety. A kubelet calls these every few seconds, so one replica
// contributes thousands of records a day: they dominate the request log, and a
// flood of 2ms requests drags the RED charts down until the latency of the
// thing people actually use is invisible inside it.
func TestTheProbesAreNotInPulsesRequestLog(t *testing.T) {
	fresh := apiRoutesGo()
	at := strings.Index(fresh, "pulse.WithExcludePaths(")
	if at < 0 {
		t.Fatal("no Pulse exclude list")
	}
	end := strings.Index(fresh[at:], ")")
	list := fresh[at : at+end]
	for _, path := range []string{"/healthz", "/readyz"} {
		if !strings.Contains(list, path) {
			t.Errorf("%s is not excluded from Pulse: %s", path, list)
		}
	}
}
