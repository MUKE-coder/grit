package scaffold

import (
	"go/format"
	"os"
	"path/filepath"
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

// The settings have to land in the file the project actually reads.
//
// They did not. addScoutConfig wrote to apps/api/.env.example, which a Grit
// project does not have: .env.example sits at the project root, beside
// grit.json, and that is also the only .env `grit doctor` reads. The write was
// skipped silently while the command printed a tick next to ".env.example",
// so `grit add scout` reported configuring something it had not touched.
func TestScoutConfigWritesToTheProjectRootEnvExample(t *testing.T) {
	root := t.TempDir()
	apiDir := filepath.Join(root, "apps", "api")
	cfgDir := filepath.Join(apiDir, "internal", "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}

	const config = `package config

type Config struct {
	// Observability (Pulse v1.0)
	PulseEnabled bool
}

func Load() *Config {
	return &Config{
		PulseEnabled: getEnv("PULSE_ENABLED", "true") == "true",
	}
}
`
	if err := os.WriteFile(filepath.Join(cfgDir, "config.go"), []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	rootEnv := filepath.Join(root, ".env.example")
	if err := os.WriteFile(rootEnv, []byte("APP_ENV=development\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	label, err := addScoutConfig(root, apiDir)
	if err != nil {
		t.Fatalf("addScoutConfig: %v", err)
	}
	if label != ".env.example" {
		t.Errorf("label = %q", label)
	}

	written, err := os.ReadFile(rootEnv)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"SCOUT_ENABLED=false", "SCOUT_USERNAME=admin", "SCOUT_PASSWORD=", "SCOUT_DB=scout.db"} {
		if !strings.Contains(string(written), want) {
			t.Errorf("the project's .env.example is missing %q", want)
		}
	}
	if _, err := os.Stat(filepath.Join(apiDir, ".env.example")); err == nil {
		t.Error("an apps/api/.env.example was created; the project's one is at the root")
	}

	// Twice is a no-op.
	before := string(written)
	if _, err := addScoutConfig(root, apiDir); err != nil {
		t.Fatalf("second run: %v", err)
	}
	after, _ := os.ReadFile(rootEnv)
	if string(after) != before {
		t.Error("running twice changed .env.example again")
	}
}

// No .env.example means the settings have nowhere to go, and saying so beats
// printing a tick beside a file that was never written.
func TestScoutConfigRefusesWithNoEnvExample(t *testing.T) {
	root := t.TempDir()
	cfgDir := filepath.Join(root, "apps", "api", "internal", "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := "package config\n\ntype Config struct {\n\t// Observability (Pulse v1.0)\n\tPulseEnabled bool\n}\n\nfunc Load() *Config {\n\treturn &Config{\n\t\tPulseEnabled: true,\n\t}\n}\n"
	if err := os.WriteFile(filepath.Join(cfgDir, "config.go"), []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := addScoutConfig(root, filepath.Join(root, "apps", "api")); err == nil {
		t.Fatal("addScoutConfig succeeded with no .env.example to write to")
	}
}
