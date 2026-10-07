package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// SQLITE_PATH was resolved against the working directory, so a project had two
// databases and no way to tell.
//
// The CLI runs migrate and seed from the API module, so the project's data went
// to apps/api/app.db. A server started anywhere else resolved the same ./app.db
// to a different file, created it empty, and answered every read with a 500
// while the real rows sat one directory away. The home page said "No blog posts
// yet" and the log said "no such table", neither of which points at the cause.
//
// internal/config/config.go is not in the upgrade's file map, so without this
// the fix would have reached new projects only.
func repairSQLitePath(root string) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, dir := range []string{
		filepath.Join(root, "apps", "api"),
		filepath.Join(root, "api"),
		root,
	} {
		path := filepath.Join(dir, "internal", "config", "config.go")
		if !fileExists(path) {
			continue
		}
		if err := repairTextFile(root, m, path, repairSQLitePathSource); err != nil {
			return err
		}
	}
	return nil
}

const (
	sqliteOldEnvLoad = `	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")`

	sqliteOldDSN = `return "sqlite:" + getEnv("SQLITE_PATH", "./app.db")`
	sqliteNewDSN = `return "sqlite:" + resolveSQLitePath(getEnv("SQLITE_PATH", "./app.db"))`

	sqliteOldPulse = `PulseStorageDSN: getEnv("PULSE_STORAGE_DSN", "pulse.db"),`
	sqliteNewPulse = `PulseStorageDSN: resolveSQLitePath(getEnv("PULSE_STORAGE_DSN", "pulse.db")),`
)

// sqlitePathHelpers is the block inserted ahead of resolveDatabaseURL. Kept
// here rather than regenerated from the template, because the template is one
// long string and this has to apply to a file somebody may have edited.
const sqlitePathHelpers = `// envDir is the directory holding the .env this process loaded, or "" when there
// is none. A relative path in that file is relative to it, not to wherever the
// process happens to have been started.
var envDir string

// loadDotEnv reads the project's .env from wherever the binary is.
//
// godotenv does not overwrite a variable that is already set, so the nearest
// file wins and the environment wins over all of them.
func loadDotEnv() {
	for _, candidate := range []string{".env", "../.env", "../../.env"} {
		if _, err := os.Stat(candidate); err != nil {
			continue
		}
		_ = godotenv.Load(candidate)
		if envDir == "" {
			envDir = filepath.Dir(candidate)
		}
	}
}

// resolveSQLitePath anchors a relative database file to the project.
//
// The CLI runs migrate and seed from the API module and the server may run from
// anywhere, so a path resolved against the working directory is two different
// files. The second one is created empty, which is the part that hurts: the
// server starts, finds no tables, and returns a 500 for every read while the
// data sits in the first.
//
// A URI form (file:...) and an absolute path are left alone: both already say
// exactly which file they mean.
func resolveSQLitePath(path string) string {
	return ResolveSQLitePath(path)
}

// ResolveSQLitePath is resolveSQLitePath, for the other things in this app that
// keep a SQLite file of their own: Sentinel's threat log and Pulse's traces.
func ResolveSQLitePath(path string) string {
	if path == "" || envDir == "" || filepath.IsAbs(path) || strings.HasPrefix(path, "file:") {
		return path
	}
	anchored := filepath.Join(envDir, path)
	if _, err := os.Stat(anchored); err == nil {
		return anchored
	}
	// A project created before this rule has its database where the old one put
	// it: beside whichever directory the CLI ran migrate from, which is the API
	// module. An existing file wins over starting an empty one beside .env,
	// because an empty database that looks exactly like the real one is the
	// whole failure this is about.
	for _, legacy := range []string{
		path,
		filepath.Join(envDir, "apps", "api", path),
		filepath.Join(envDir, "api", path),
	} {
		if _, err := os.Stat(legacy); err != nil {
			continue
		}
		if abs, err := filepath.Abs(legacy); err == nil {
			log.Printf("sqlite: using %s, the database this project already has. Move it beside .env to have every command agree on it wherever they run from.", abs)
		}
		return legacy
	}
	return anchored
}

`

// repairSQLitePathSource rewrites one config.go.
func repairSQLitePathSource(src string) (string, []string, []string) {
	var changes []string

	if strings.Contains(src, "func ResolveSQLitePath(") {
		return src, nil, nil
	}

	anchor := "// resolveDatabaseURL builds the DSN"
	if !strings.Contains(src, anchor) {
		// A config.go this far from the template is the developer's. Saying so
		// beats editing it blind.
		return src, nil, []string{"could not find resolveDatabaseURL, so the SQLite path was left as it is: a server started outside the API directory will use a different database file"}
	}

	if strings.Contains(src, sqliteOldEnvLoad) {
		src = strings.Replace(src, sqliteOldEnvLoad, "\tloadDotEnv()", 1)
		changes = append(changes, "reads .env once and remembers where it was")
	}
	src = strings.Replace(src, anchor, sqlitePathHelpers+anchor, 1)
	changes = append(changes, "anchors the SQLite database to the project instead of the working directory")

	if strings.Contains(src, sqliteOldDSN) {
		src = strings.Replace(src, sqliteOldDSN, sqliteNewDSN, 1)
	}
	if strings.Contains(src, sqliteOldPulse) {
		src = strings.Replace(src, sqliteOldPulse, sqliteNewPulse, 1)
	}

	// The helpers need these three. Every scaffolded config.go imports os, log
	// and strings already; filepath is the one a very old file can be missing.
	src = ensureConfigImport(src, `"path/filepath"`, `"os"`)

	return src, changes, nil
}

// ensureConfigImport adds one import to config.go's block, after an import that
// is certain to be there, so the result stays gofmt-stable.
func ensureConfigImport(src, want, after string) string {
	if strings.Contains(src, want) {
		return src
	}
	idx := strings.Index(src, "\n\t"+after+"\n")
	if idx < 0 {
		return src
	}
	insertAt := idx + len("\n\t"+after)
	return src[:insertAt] + "\n\t" + want + src[insertAt:]
}
