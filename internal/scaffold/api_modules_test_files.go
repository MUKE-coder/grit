package scaffold

// configModulesTestGo emits internal/config/modules_test.go.
//
// The load-bearing assertion is TestModules_DefaultOn: every module defaults to
// TRUE, so an app upgrading to a version that adds this never silently loses a
// feature. The others pin that disabling one module doesn't disturb the rest,
// that an unknown name fails CLOSED (a typo hides a feature rather than
// exposing one), and that Map() stays in step with the struct — the admin hides
// nav from Map(), so a missing key leaves a dead link.
func configModulesTestGo() string {
	return `package config

import (
	"os"
	"testing"
)

// Load validates unrelated required settings; satisfy them so these tests
// exercise the module flags and nothing else.
func init() {
	os.Setenv("JWT_SECRET", "test-secret-key-for-module-flag-tests-only")
}

// Default ON: an existing app upgrading must behave exactly as before.
func TestModules_DefaultOn(t *testing.T) {
	for _, k := range []string{"MODULE_AI", "MODULE_JOBS", "MODULE_BACKUP"} {
		os.Unsetenv(k)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for name, on := range cfg.Modules.Map() {
		if !on {
			t.Errorf("module %q defaulted to OFF; upgrading an app must not disable features", name)
		}
	}
}

// Setting a flag false turns exactly that module off.
func TestModules_DisableOne(t *testing.T) {
	os.Setenv("MODULE_AI", "false")
	defer os.Unsetenv("MODULE_AI")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Modules.AI {
		t.Error("MODULE_AI=false did not disable the AI module")
	}
	if !cfg.Modules.Jobs {
		t.Error("disabling AI must not affect other modules")
	}
	if cfg.Modules.Enabled("ai") {
		t.Error("Enabled(\"ai\") disagrees with the struct field")
	}
	if !cfg.Modules.Enabled("jobs") {
		t.Error("Enabled(\"jobs\") disagrees with the struct field")
	}
}

// An unknown name must report false — a typo should hide a feature, never
// silently expose one.
func TestModules_UnknownNameIsOff(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Modules.Enabled("not-a-module") {
		t.Error("unknown module name reported enabled")
	}
}

// Map must cover every field, or the admin can't hide a nav entry it should.
func TestModules_MapCoversAllFields(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	m := cfg.Modules.Map()
	if len(m) != 11 {
		t.Errorf("Map has %d entries; keep it in step with the ModuleFlags fields", len(m))
	}
	for name := range m {
		if !cfg.Modules.Enabled(name) {
			t.Errorf("Map key %q is not handled by Enabled()", name)
		}
	}
}
`
}
