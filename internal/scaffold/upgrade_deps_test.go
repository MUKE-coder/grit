package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestVersionLess(t *testing.T) {
	for _, tc := range []struct {
		a, b string
		less bool
	}{
		{"v2.2.1", "v2.5.0", true},
		{"v2.5.0", "v2.5.0", false},
		{"v2.6.0", "v2.5.0", false},
		{"v1.0.1", "v1.1.0", true},
		// Numbers, not strings: v1.10.0 is the newer one.
		{"v1.10.0", "v1.9.0", false},
		{"v1.9.0", "v1.10.0", true},
		// A pseudo-version is a pre-release of v0.0.0.
		{"v0.0.0-20260529025319-478cdfa8ce5f", "v1.0.0", true},
		{"v1.0.0-rc.1", "v1.0.0", true},
		{"v1.0.0", "v1.0.0-rc.1", false},
		{"v2.0.0+incompatible", "v2.0.0", false},
	} {
		if got := versionLess(tc.a, tc.b); got != tc.less {
			t.Errorf("versionLess(%s, %s) = %v, want %v", tc.a, tc.b, got, tc.less)
		}
	}
}

// A project scaffolded on an older Sentinel stayed on it: grit upgrade rewrote
// framework files and never go.mod, and v2.2.2 is a security release. It is
// raised now, and a library already past its floor is left where it is.
func TestUpgradeRaisesFrameworkDepsAndNeverLowersThem(t *testing.T) {
	dir := t.TempDir()
	write := func(body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var fetched []string
	calls := 0
	orig := goGet
	goGet = func(_ string, specs ...string) error {
		calls++
		fetched = specs
		return nil
	}
	defer func() { goGet = orig }()

	write(`module shop/apps/api

go 1.24

require (
	github.com/MUKE-coder/gorm-studio v1.0.1
	github.com/MUKE-coder/pulse v0.0.0-20260529025319-478cdfa8ce5f
	github.com/MUKE-coder/sentinel/v2 v2.9.0
	github.com/gin-gonic/gin v1.11.0
)
`)
	raised, err := raiseFrameworkDeps(dir)
	if err != nil {
		t.Fatalf("raise: %v", err)
	}
	// Derived, not written out: a floor bump must not break a test about
	// whether raising works at all.
	want := []string{"github.com/MUKE-coder/gorm-studio@" + floorOf(t, "github.com/MUKE-coder/gorm-studio"),
		"github.com/MUKE-coder/pulse@" + floorOf(t, "github.com/MUKE-coder/pulse")}
	if !reflect.DeepEqual(fetched, want) {
		t.Errorf("fetched %v, want %v (and never a lower Sentinel than the project has)", fetched, want)
	}
	if len(raised) != 2 {
		t.Errorf("reported %d raises, want 2: %v", len(raised), raised)
	}

	// Everything at or past its floor: go get is not run at all.
	calls = 0
	write("module shop/apps/api\n\ngo 1.24\n\nrequire (\n" +
		requireAtFloors("github.com/MUKE-coder/gorm-studio",
			"github.com/MUKE-coder/pulse", "github.com/MUKE-coder/sentinel/v2") + ")\n")
	if raised, err := raiseFrameworkDeps(dir); err != nil || raised != nil || calls != 0 {
		t.Errorf("an up-to-date project ran go get %d time(s): %v %v", calls, raised, err)
	}

	// No go.mod: nothing to do, and no error.
	if raised, err := raiseFrameworkDeps(t.TempDir()); err != nil || raised != nil {
		t.Errorf("a directory with no go.mod: %v %v", raised, err)
	}
}

// floorOf returns the current floor for path, so a test can say "the floor"
// instead of naming a version that goes stale on the next bump.
func floorOf(t *testing.T, path string) string {
	t.Helper()
	for _, d := range frameworkDeps {
		if d.Path == path {
			return d.Floor
		}
	}
	t.Fatalf("%s is not a framework dependency", path)
	return ""
}

// requireAtFloors renders go.mod require lines holding each path exactly at its
// floor, for the case where nothing should be raised.
func requireAtFloors(paths ...string) string {
	var b strings.Builder
	for _, path := range paths {
		for _, d := range frameworkDeps {
			if d.Path == path {
				fmt.Fprintf(&b, "\t%s %s\n", d.Path, d.Floor)
			}
		}
	}
	return b.String()
}

// The four advisories published on 2026-11-24 reach an EXISTING project, and
// the scaffold moving does nothing for one. otel was in frameworkRequires,
// which only adds a module a project lacks, so a project that already had
// internal/tracing kept v1.45.0 through every upgrade; saml and websocket were
// in neither list. This runs the real floor list against a go.mod carrying all
// three vulnerable versions.
func TestTheNovemberAdvisoriesAreRaisedOnAnExistingProject(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(`module app

go 1.26.11

require (
	github.com/crewjam/saml v0.5.1
	github.com/gorilla/websocket v1.5.3
	go.opentelemetry.io/otel v1.45.0
	go.opentelemetry.io/otel/sdk v1.45.0
	go.opentelemetry.io/otel/trace v1.45.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.45.0
)
`), 0o644); err != nil {
		t.Fatal(err)
	}

	var fetched []string
	orig := goGet
	goGet = func(_ string, specs ...string) error { fetched = specs; return nil }
	defer func() { goGet = orig }()

	if _, err := raiseFrameworkDeps(dir); err != nil {
		t.Fatalf("raising: %v", err)
	}

	for _, want := range []string{
		"github.com/crewjam/saml@v0.5.2",
		"github.com/gorilla/websocket@v1.5.4-0.20240701034025-d67f41855da4",
		"go.opentelemetry.io/otel@v1.47.0",
		"go.opentelemetry.io/otel/sdk@v1.47.0",
		"go.opentelemetry.io/otel/trace@v1.47.0",
		"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp@v1.47.0",
	} {
		var found bool
		for _, got := range fetched {
			if got == want {
				found = true
			}
		}
		if !found {
			t.Errorf("grit upgrade does not raise %s, so an existing project keeps the "+
				"advisory however many times it upgrades. Raised: %v", want, fetched)
		}
	}
}

// websocket's floor is a PRE-RELEASE, because upstream fixed GO-2026-6278 on
// master and never tagged it. A pre-release sorts below the release it is a
// pre-release of, so the day v1.5.4 is cut, a project that moved to the real
// tag must not be dragged back to a commit from July 2024.
func TestThePreReleaseWebsocketFloorNeverDowngrades(t *testing.T) {
	const floor = "v1.5.4-0.20240701034025-d67f41855da4"

	// It is the floor the templates ship, or this test is about nothing.
	var declared string
	for _, d := range frameworkDeps {
		if d.Path == "github.com/gorilla/websocket" {
			declared = d.Floor
		}
	}
	if declared != floor {
		t.Fatalf("the websocket floor is %q, not the pre-release this test is about", declared)
	}

	for _, tc := range []struct {
		current string
		raise   bool
		why     string
	}{
		{"v1.5.3", true, "the vulnerable release the advisory names"},
		{"v1.5.2", true, "clean, but two releases of other fixes behind"},
		{floor, false, "already there: raising it again would run go get on every upgrade"},
		{"v1.5.4", false, "the real tag, once it exists, is ABOVE a pre-release of itself"},
		{"v1.6.0", false, "past it"},
	} {
		if got := versionLess(tc.current, floor); got != tc.raise {
			t.Errorf("websocket %s: raise = %v, want %v (%s)", tc.current, got, tc.raise, tc.why)
		}
	}
}
