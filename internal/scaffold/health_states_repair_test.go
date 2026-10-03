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
	for name, page := range map[string]string{
		"admin":   adminSystemHealthPage(),
		"desktop": desktopClientSystemHealthPage(),
	} {
		if !strings.Contains(page, `"ok" | "degraded" | "off" | "unknown"`) {
			t.Errorf("%s: the page does not know the four states", name)
		}
		if !strings.Contains(page, "storage") {
			t.Errorf("%s: the page still does not report storage", name)
		}
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
