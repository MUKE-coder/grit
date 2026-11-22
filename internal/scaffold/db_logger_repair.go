package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The GORM logger a project connects with, for projects that already exist.
//
// A new project gets the right one from apiDatabaseGo. database.go is not
// rewritten by an upgrade, so without this an existing project keeps printing
// a red block with a file and a line number every time a seeder looks for a
// row it expects to be absent, which on a first run reads as a failure
// (grit#238).
//
// Its own pass and its own anchor: the pool repair edits a different part of
// the same file, and two repairs that converge on one marker cannot both run.

// dbLoggerLineOld is the one-line logger database.go carried before this.
const dbLoggerLineOld = `		Logger:      logger.Default.LogMode(logLevel),`

// dbLoggerBlockNew keeps everything the default logger did, and stops
// reporting the one error every caller already handles.
const dbLoggerBlockNew = `		// ErrRecordNotFound is never worth a log line: every caller that can
		// get one handles it, and GORM prints it in red with a file and a line
		// number, which on a first run looks exactly like a failure. Seeding a
		// new project printed two such blocks for two lookups whose whole
		// purpose was to miss (grit#238). Everything else about the default
		// logger is kept, so a slow query still reports itself.
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),`

// repairDBLogger gives an existing project the quieter logger.
func repairDBLogger(root string, opts Options) error {
	path := filepath.Join(opts.APIRoot(root), "internal", "database", "database.go")
	if !fileExists(path) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, path, repairDBLoggerSource)
}

// repairDBLoggerSource is the edit itself.
//
// Guarded on `IgnoreRecordNotFoundError:` with its colon rather than on the
// bare word, because the word also appears in the comment that explains the
// change: a Contains check on the prose matches the explanation and makes the
// repair a silent no-op on exactly the projects that still need it.
func repairDBLoggerSource(src string) (string, []string, []string) {
	if strings.Contains(src, "IgnoreRecordNotFoundError:") {
		return src, nil, nil
	}
	if strings.Count(src, dbLoggerLineOld) != 1 {
		return src, nil, []string{
			"the GORM logger is not the line Grit wrote, so it was left alone: " +
				"to stop seeders printing a red block for a lookup that is meant to " +
				"miss, build the logger with logger.Config{IgnoreRecordNotFoundError: true}",
		}
	}
	// log, os and time are already imported by every database.go Grit wrote,
	// for the pool settings and the DSN. Checked rather than assumed, because
	// adding an import here would mean parsing the block.
	for _, pkg := range []string{`"log"`, `"os"`, `"time"`} {
		if !strings.Contains(src, pkg) {
			return src, nil, []string{
				"database.go does not import " + pkg + ", so the quieter GORM " +
					"logger was not applied",
			}
		}
	}
	out := strings.Replace(src, dbLoggerLineOld, dbLoggerBlockNew, 1)
	return out, []string{"GORM no longer logs the record-not-found a seeder expects"}, nil
}
