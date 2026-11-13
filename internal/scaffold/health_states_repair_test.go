package scaffold

import (
	"strings"
	"testing"
)

func TestRepairHealthStates(t *testing.T) {
	out, fixed, warn := repairHealthStatesSource("demo")(v3349Routes(t))
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	for _, want := range []string{
		`"demo/internal/health"`,
		"health.Overall(comps)",
		`comps["storage"] = health.Off(`,
		`comps["redis"] = health.Off(`,
		"func eventBusStatus() health.Report {",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repaired routes.go is missing %q", want)
		}
	}
	if strings.Contains(out, "type compStatus struct {") {
		t.Error("the boolean-per-component struct survived the repair")
	}
	if strings.Contains(out, `overall := "ok"`) {
		t.Error("the hand-written overall status survived the repair")
	}
}

func TestRepairHealthStatesIsIdempotent(t *testing.T) {
	first, _, _ := repairHealthStatesSource("demo")(v3349Routes(t))
	again, fixed, warn := repairHealthStatesSource("demo")(first)
	if again != first || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("a second upgrade changed routes.go again (fixed %v, warn %v)", fixed, warn)
	}
}

func TestRepairHealthStatesLeavesAnEditedHandlerAlone(t *testing.T) {
	// Somebody who added their own component to the handler gets told about
	// health.Register, not overwritten. The whole reason the repair matches the
	// exact text is that this handler is in a file people edit.
	src := strings.Replace(v3349Routes(t),
		`"email":    mailStatus,`,
		`"email":    mailStatus,`+"\n\t\t\t\"billing\": billingStatus(),", 1)
	out, fixed, warn := repairHealthStatesSource("demo")(src)
	if out != src || len(fixed) > 0 || len(warn) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if !strings.Contains(warn[0], "health.Register") {
		t.Errorf("the warning should name the supported way to add a component, got %q", warn[0])
	}
}

func TestRepairHealthStatesSaysNothingAboutAFileWithNoHealthCheck(t *testing.T) {
	// A single's routes, a plugin's routes, anything without the handler. A
	// warning here would fire on every upgrade of a project that was never
	// wrong.
	src := "package routes\n\nimport (\n\t\"github.com/gin-gonic/gin\"\n)\n\nfunc Setup(r *gin.Engine) {}\n"
	out, fixed, warn := repairHealthStatesSource("demo")(src)
	if out != src || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("fixed %v, warned %v", fixed, warn)
	}
}

func TestRepairHealthStatesWaitsForTheRealtimeStats(t *testing.T) {
	// The handler this writes reports realtimeHub.Stats(), so a project whose
	// hub predates them would get a reference to a variable it has not got. The
	// realtime repair runs first and this one catches up next time.
	src := strings.Replace(v3349Routes(t), `"realtime": realtimeHub.Stats(),`, "", 1)
	out, fixed, warn := repairHealthStatesSource("demo")(src)
	if out != src || len(fixed) > 0 || len(warn) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
}

func TestFreshRoutesAnswerInFourStates(t *testing.T) {
	src := apiRoutesGo()
	for _, want := range []string{
		`"{{MODULE}}/internal/health"`,
		"comps := health.Snapshot()",
		`comps["storage"] = health.OK(`,
		`"status":  health.Overall(comps),`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("a fresh routes.go is missing %q", want)
		}
	}
	// The two words the four states exist to separate must not be the same word
	// any more: nothing in the handler decides the overall status by hand.
	if strings.Contains(src, "svc.Cache != nil && !redisStatus.OK") {
		t.Error("the overall status still excludes an unconfigured Redis by hand")
	}
	if out, fixed, warn := repairHealthStatesSource("demo")(src); out != src || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("a fresh routes.go still needs the repair (fixed %v, warn %v)", fixed, warn)
	}
}

func TestHealthPackageIsFrameworkOwned(t *testing.T) {
	// It has to travel on upgrade, because the repair above writes health.OK
	// into routes.go and a project without the package would not compile.
	root := t.TempDir()
	opts := Options{ProjectName: "demo", Architecture: ArchTriple, Frontend: FrontendNext}
	if err := writeFrameworkOwnedFiles(root, opts); err != nil {
		t.Fatalf("writeFrameworkOwnedFiles: %v", err)
	}
	for _, rel := range []string{"internal/health/health.go", "internal/health/health_test.go"} {
		if !fileExists(opts.APIRoot(root) + "/" + rel) {
			t.Errorf("%s was not written", rel)
		}
	}
}

func TestHealthPageReadsTheStateItIsSent(t *testing.T) {
	// The four states, wherever each surface declares them.
	//
	// The admin's moved to hooks/use-system-health in v3.535.0, which is the
	// one query the top bar, the sidebar badge, the hub tiles and this page
	// all read: two copies of a type describing one endpoint's response is two
	// places for "off" to come to mean different things. The desktop client
	// has its own, because it is a separate application.
	for name, src := range map[string]string{
		"admin":   adminUseSystemHealth(),
		"desktop": desktopClientSystemHealthPage(),
	} {
		if !strings.Contains(src, `"ok" | "degraded" | "off" | "unknown"`) {
			t.Errorf("%s: does not know the four states", name)
		}
		if !strings.Contains(src, "storage") {
			t.Errorf("%s: still does not report storage", name)
		}
	}

	// And the page reads them rather than declaring its own.
	if strings.Contains(adminSystemHealthPage(), `type HealthState =`) {
		t.Error("the health page declares HealthState again. The shared query already " +
			"exports it, and two declarations of one endpoint's response is how the " +
			"pill in the top bar and this page come to disagree about what off means")
	}
	if !strings.Contains(adminSystemHealthPage(), `from "@/hooks/use-system-health"`) {
		t.Error("the health page does not read the shared health query, so it fetches " +
			"the same endpoint a second time and can report a different answer from " +
			"the top bar on the same screen")
	}
	admin := adminSystemHealthPage()
	// The guesses that each read the boolean differently.
	for _, gone := range []string{
		`data.redis?.ok === false ? "down"`,
		`data.email?.configured ? "ok" : "unknown"`,
		`return { status: "ok", database: { ok: true }, api: { ok: true } };`,
	} {
		if strings.Contains(admin, gone) {
			t.Errorf("the admin page still infers a state: %q", gone)
		}
	}
}
