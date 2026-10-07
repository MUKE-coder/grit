package scaffold

import (
	"strings"
	"testing"
)

// One SQLite database per project, wherever the process runs from.
//
// SQLITE_PATH is relative and was resolved against the working directory. The
// CLI runs migrate and seed from apps/api; a built binary runs from wherever it
// is started. The two resolved ./app.db to different files, and the second was
// created empty: the server found no tables and returned 500 for every read
// while the data sat one directory away, with nothing saying so.
func TestSQLitePathIsAnchoredToTheProject(t *testing.T) {
	cfg := apiConfigGo()

	if !strings.Contains(cfg, `resolveSQLitePath(getEnv("SQLITE_PATH", "./app.db"))`) {
		t.Error("the database path is still whatever the working directory makes of it")
	}
	if !strings.Contains(cfg, "var envDir string") || !strings.Contains(cfg, "func loadDotEnv()") {
		t.Error("nothing records which directory .env came from")
	}
	// The three candidates have to stay, and the first one found is the anchor.
	for _, candidate := range []string{`".env", "../.env", "../../.env"`} {
		if !strings.Contains(cfg, candidate) {
			t.Errorf("the .env search lost %s", candidate)
		}
	}
	// An absolute path and a URI already name their file.
	if !strings.Contains(cfg, "filepath.IsAbs(path)") || !strings.Contains(cfg, `strings.HasPrefix(path, "file:")`) {
		t.Error("an absolute path or a file: URI would be rewritten")
	}
	// A project from before this keeps the database it has, and that file is
	// beside the API module, which is where the CLI ran migrate from.
	if !strings.Contains(cfg, "the database this project already has") {
		t.Error("an existing project's database would be abandoned for an empty one")
	}
	for _, legacy := range []string{`filepath.Join(envDir, "apps", "api", path)`, `filepath.Join(envDir, "api", path)`} {
		if !strings.Contains(cfg, legacy) {
			t.Errorf("an upgraded project's database at %s is not looked for, so a server started from the project root would begin an empty one", legacy)
		}
	}
	// Sentinel's threat log and Pulse's traces had the same split.
	if !strings.Contains(apiRoutesGo(), `config.ResolveSQLitePath("sentinel.db")`) {
		t.Error("Sentinel still keeps its threat log wherever the process started")
	}
	if !strings.Contains(cfg, `resolveSQLitePath(getEnv("PULSE_STORAGE_DSN", "pulse.db"))`) {
		t.Error("Pulse still keeps its traces wherever the process started")
	}
}
