package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

const oldServerTimeouts = "\t\tHandler:      router,\n\t\tReadTimeout:  15 * time.Second,\n\t\tWriteTimeout: 15 * time.Second,\n\t\tIdleTimeout:  60 * time.Second,\n"

func TestRepairRequestLimitsRoutes(t *testing.T) {
	fresh := apiRoutesGo()
	old := strings.Replace(fresh, routesRequestLimitsBlock, oldBodyCapLine, 1)
	if old == fresh {
		t.Fatal("could not rebuild routes.go from before the repair")
	}
	out, fixed, warn := repairRequestLimitsRoutesSource(old)
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if out != fresh {
		t.Error("the repaired routes.go differs from a fresh one")
	}
	if again, fixed, _ := repairRequestLimitsRoutesSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed routes.go again")
	}
}

func TestRepairRequestLimitsWarnsOnAnEditedCap(t *testing.T) {
	src := "package routes\n\nfunc f() {\n\tr.Use(middleware.MaxBodySize(50 << 20))\n}\n"
	if out, fixed, warn := repairRequestLimitsRoutesSource(src); out != src || len(fixed) > 0 || len(warn) != 1 {
		t.Errorf("fixed %v, warned %v", fixed, warn)
	}
}

func TestRepairServerTimeouts(t *testing.T) {
	for name, fresh := range map[string]string{
		"double": apiMainGo(Options{ProjectName: "x"}),
		"single": singleMainGo(Options{ProjectName: "x"}),
	} {
		old := strings.Replace(fresh, "\t\tHandler: router,\n"+serverTimeoutFields, oldServerTimeouts, 1)
		old = strings.Replace(old, "\t\tAddr:    ", "\t\tAddr:         ", 1)
		if old == fresh {
			t.Fatalf("%s: could not rebuild main.go from before the repair", name)
		}
		out, fixed, warn := repairServerTimeoutsSource(old)
		if len(warn) > 0 || len(fixed) != 1 {
			t.Fatalf("%s: fixed %v, warned %v", name, fixed, warn)
		}
		if out != fresh {
			t.Errorf("%s: the repaired main.go differs from a fresh one", name)
		}
		if again, fixed, _ := repairServerTimeoutsSource(out); again != out || len(fixed) > 0 {
			t.Errorf("%s: a second upgrade changed main.go again", name)
		}
	}
}

func TestRepairServerTimeoutsWarnsOnAnEditedServer(t *testing.T) {
	src := "package main\n\nvar srv = &http.Server{\n\tAddr:         addr,\n\tWriteTimeout: 15 * time.Second,\n}\n"
	if out, fixed, warn := repairServerTimeoutsSource(src); out != src || len(fixed) > 0 || len(warn) != 1 {
		t.Errorf("fixed %v, warned %v", fixed, warn)
	}
}

func TestFreshTemplatesNeedNoRequestLimitsRepair(t *testing.T) {
	if src := apiRoutesGo(); !strings.Contains(src, "middleware.RequestLimits(") || strings.Contains(src, "middleware.MaxBodySize(") {
		t.Error("routes.go still caps every body with MaxBodySize")
	}
	for name, src := range map[string]string{
		"double": apiMainGo(Options{ProjectName: "x"}),
		"single": singleMainGo(Options{ProjectName: "x"}),
	} {
		if out, fixed, warn := repairServerTimeoutsSource(src); out != src || len(fixed) > 0 || len(warn) > 0 {
			t.Errorf("%s: main.go still needs the repair (fixed %v, warn %v)", name, fixed, warn)
		}
	}
	for name, src := range map[string]string{"limits.go": middlewareLimitsGo(), "limits_test.go": middlewareLimitsTestGo()} {
		if _, err := format.Source([]byte(src)); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}
