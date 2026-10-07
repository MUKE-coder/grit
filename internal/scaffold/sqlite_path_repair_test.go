package scaffold

import (
	"strings"
	"testing"
)

// The fix has to reach a project that already exists, because config.go is not
// in the upgrade's file map.
func TestRepairSQLitePathRewritesAnOldConfig(t *testing.T) {
	old := `package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func Load() (*Config, error) {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	cfg := &Config{
		PulseStorageDSN: getEnv("PULSE_STORAGE_DSN", "pulse.db"),
	}
	return cfg, nil
}

// resolveDatabaseURL builds the DSN from DB_PROVIDER and that provider's parts.
func resolveDatabaseURL() string {
	switch provider {
	case "sqlite":
		return "sqlite:" + getEnv("SQLITE_PATH", "./app.db")
	}
	return ""
}
`
	out, changes, warnings := repairSQLitePathSource(old)
	if len(warnings) != 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if len(changes) == 0 {
		t.Fatal("the repair reported no changes")
	}
	for _, want := range []string{
		"func ResolveSQLitePath(",
		"loadDotEnv()",
		`resolveSQLitePath(getEnv("SQLITE_PATH", "./app.db"))`,
		`resolveSQLitePath(getEnv("PULSE_STORAGE_DSN", "pulse.db"))`,
		`"path/filepath"`,
		`filepath.Join(envDir, "apps", "api", path)`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repaired config has no %s", want)
		}
	}
	if strings.Contains(out, `_ = godotenv.Load("../../.env")`) {
		t.Error("the old three-line load survived, so envDir is never set")
	}

	// Running it twice changes nothing more.
	again, changes2, _ := repairSQLitePathSource(out)
	if again != out || len(changes2) != 0 {
		t.Error("the repair is not idempotent")
	}

	// A config.go that has been rewritten past recognition is said out loud
	// rather than edited blind.
	_, _, warn := repairSQLitePathSource("package config\n\nfunc Load() {}\n")
	if len(warn) == 0 {
		t.Error("an unrecognisable config.go was edited silently")
	}
}
