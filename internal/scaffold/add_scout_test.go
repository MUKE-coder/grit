package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

// The generated file has to compile in somebody's project, and the usual way
// a template stops compiling is an edit to the string it lives in.
func TestScoutRoutesGoParses(t *testing.T) {
	src := scoutRoutesGo("example.com/app")
	if _, err := format.Source([]byte(src)); err != nil {
		t.Fatalf("generated internal/routes/scout.go does not parse: %v", err)
	}
	if strings.Contains(src, "{{MODULE}}") {
		t.Error("the module placeholder survived into the output")
	}
	for _, want := range []string{
		`"example.com/app/internal/access"`,
		`"example.com/app/internal/config"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated file does not import %s", want)
		}
	}
}

// The integration is the route table. Without WithRoutes, Scout falls back to
// gin introspection, which carries no middleware, which is exactly the thing
// this command exists to supply.
func TestScoutMountHandsOverTheRouteTable(t *testing.T) {
	src := scoutRoutesGo("example.com/app")
	for _, want := range []string{
		"scout.WithRoutes(scoutRoutes()...)",
		"access.All()",
		"collectors/all",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated file is missing %q", want)
		}
	}
}

// A public route reports an empty chain rather than a made-up one. Scout reads
// the chain to decide whether a route is gated, so inventing an entry here
// would mean reporting an unprotected route as protected.
func TestScoutMiddlewareLeavesPublicRoutesEmpty(t *testing.T) {
	src := scoutRoutesGo("example.com/app")
	if !strings.Contains(src, "if req.Public() {\n\t\treturn nil\n\t}") {
		t.Error("scoutMiddleware does not return an empty chain for a public route")
	}
	for _, want := range []string{
		`mw := []string{"AuthRequired"}`,
		`"RequireRole("+role+")"`,
		`"RequirePermission("+perm+")"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated file is missing %q", want)
		}
	}
}

// The dashboard key is derived from the app secret, not equal to it: one
// signing key shared between two auth systems means a token minted by either
// is valid to the other.
func TestScoutDashboardKeyIsDerivedNotShared(t *testing.T) {
	src := scoutRoutesGo("example.com/app")
	if strings.Contains(src, "JWTSecret:   cfg.JWTSecret") {
		t.Error("the dashboard is signed with the application's own JWT secret")
	}
	if !strings.Contains(src, "JWTSecret:   scoutDashboardKey(cfg.JWTSecret)") {
		t.Error("the dashboard key is not derived from the app secret")
	}
	if !strings.Contains(src, "sha256.Sum256") {
		t.Error("scoutDashboardKey does not hash")
	}
}

// It mounts disabled, and refuses to mount without a password. The dashboard
// lists where the application is weakest.
func TestScoutRefusesToMountUnconfigured(t *testing.T) {
	src := scoutRoutesGo("example.com/app")
	if !strings.Contains(src, "if !cfg.ScoutEnabled {") {
		t.Error("generated file does not check ScoutEnabled")
	}
	if !strings.Contains(src, `if cfg.ScoutPassword == "" {`) {
		t.Error("generated file mounts without requiring a password")
	}
}

// The pin is exact. A range is how a generated project once adopted a
// one-day-old release that did not work.
func TestScoutVersionIsPinnedExactly(t *testing.T) {
	if strings.ContainsAny(scoutVersion, "^~><= ") {
		t.Errorf("scoutVersion %q is a range, not a pin", scoutVersion)
	}
	if !strings.HasPrefix(scoutVersion, "v") {
		t.Errorf("scoutVersion %q is not a Go version", scoutVersion)
	}
	for _, r := range scoutRequires {
		if !strings.HasPrefix(r.version, "v") {
			t.Errorf("%s is pinned to %q, which is not a Go version", r.module, r.version)
		}
	}
}
