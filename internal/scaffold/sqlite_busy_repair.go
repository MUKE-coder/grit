package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// "database is locked" on a machine doing nothing.
//
// A SQLite file opened with GORM's defaults has no busy timeout and no WAL, so
// the second writer fails instantly instead of waiting, and the pool hands out
// several connections to compete over one file lock. Two writes a few
// milliseconds apart is enough: building the blueprints site, a subscription
// webhook lost its write to a background worker's tick, and what the customer
// saw was a purchase that granted nothing. Nothing in the logs said "SQLite" —
// only SQLITE_BUSY from a handler that had done nothing wrong.
//
// Three lines fix it: busy_timeout makes a writer wait its turn, WAL lets
// readers carry on during a write, and a pool of one turns the competition
// into a queue inside the process.
const (
	sqliteOpenOld = `	case strings.HasPrefix(dsn, "sqlite://"):
		db, err = gorm.Open(sqlite.Open(strings.TrimPrefix(dsn, "sqlite://")), gormCfg)
	case strings.HasPrefix(dsn, "sqlite:"):
		db, err = gorm.Open(sqlite.Open(strings.TrimPrefix(dsn, "sqlite:")), gormCfg)`
	sqliteOpenNew = `	case strings.HasPrefix(dsn, "sqlite://"):
		db, err = gorm.Open(sqlite.Open(sqliteDSN(strings.TrimPrefix(dsn, "sqlite://"))), gormCfg)
	case strings.HasPrefix(dsn, "sqlite:"):
		db, err = gorm.Open(sqlite.Open(sqliteDSN(strings.TrimPrefix(dsn, "sqlite:"))), gormCfg)`

	// The helper goes next to the other small helper in database.go.
	sqliteDSNAnchor = `// getEnvInt reads a whole-number env var. A malformed value falls back rather`

	sqliteDSNHelper = `// sqliteDSN adds the two pragmas a SQLite file needs to survive more than one
// writer.
//
// Without them a second writer fails at once with "database is locked", which
// is not a load problem: a background worker ticking while a webhook writes is
// enough, and the error surfaces as a failed request nobody can reproduce by
// hand. busy_timeout makes a writer wait its turn instead of giving up, and
// WAL lets reads carry on while one write is in flight.
//
// Anything the caller set is left alone, and :memory: gets neither: there is
// no file to journal and a memory database is one connection's own.
func sqliteDSN(path string) string {
	if path == "" || strings.Contains(path, ":memory:") || strings.Contains(path, "mode=memory") {
		return path
	}
	if !strings.Contains(path, "busy_timeout") {
		path += pragmaJoin(path) + "_pragma=busy_timeout(5000)"
	}
	if !strings.Contains(path, "journal_mode") {
		path += pragmaJoin(path) + "_pragma=journal_mode(WAL)"
	}
	return path
}

func pragmaJoin(path string) string {
	if strings.Contains(path, "?") {
		return "&"
	}
	return "?"
}

`

	sqlitePoolAnchor = `	maxIdle := getEnvInt("DB_MAX_IDLE_CONNS", maxOpen)
	if maxIdle < 1 || maxIdle > maxOpen {
		maxIdle = maxOpen
	}
`
	// Word for word what the template writes, so a repaired project and a fresh
	// one are the same file.
	sqlitePoolClamp = `
	// SQLite takes one connection, whatever the numbers above say. It was long
	// written here that SQLite "ignores most of these", and it does not: a pool
	// of several connections is several writers competing for one file lock, so
	// a transaction holding the write lock on one connection and a statement
	// wanting it on another wait for each other until busy_timeout gives up.
	// What comes back is "database is locked" on a machine doing almost
	// nothing: a background worker ticking while a webhook writes was enough to
	// lose the write, and it reads as a load problem it is not.
	//
	// One connection makes that a queue inside the process, which is what
	// SQLite wants. It is also the only correct setting for :memory:, where a
	// second connection is a second, empty database.
	if strings.HasPrefix(dsn, "sqlite:") {
		maxOpen, maxIdle = 1, 1
	}

`
)

func repairSQLiteBusySource(src string) (string, []string, []string) {
	var changes []string
	var warnings []string
	out := src

	if !strings.Contains(out, "func sqliteDSN(") {
		switch {
		case strings.Count(out, sqliteOpenOld) == 1 && strings.Count(out, sqliteDSNAnchor) == 1:
			out = strings.Replace(out, sqliteOpenOld, sqliteOpenNew, 1)
			out = strings.Replace(out, sqliteDSNAnchor, sqliteDSNHelper+sqliteDSNAnchor, 1)
			changes = append(changes, "SQLite waits its turn instead of failing with \"database is locked\" (busy_timeout + WAL)")
		case strings.Contains(out, "sqlite.Open("):
			warnings = append(warnings, "the SQLite connection in internal/database/database.go is not the one Grit wrote: append ?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL) to the file path, or a second writer fails instantly")
		}
	}

	if !strings.Contains(out, `if strings.HasPrefix(dsn, "sqlite:") {
		maxOpen, maxIdle = 1, 1`) {
		if strings.Count(out, sqlitePoolAnchor) == 1 {
			out = strings.Replace(out, sqlitePoolAnchor, sqlitePoolAnchor+sqlitePoolClamp, 1)
			changes = append(changes, "the connection pool is one connection on SQLite, so two writers queue instead of colliding")
		} else if strings.Contains(out, "SetMaxOpenConns") {
			warnings = append(warnings, "the connection pool block in internal/database/database.go is not the one Grit wrote: set MaxOpenConns and MaxIdleConns to 1 when the DSN is SQLite")
		}
	}

	return out, changes, warnings
}

// repairSQLiteBusy applies both to an existing project.
func repairSQLiteBusy(root string) error {
	path := filepath.Join(root, "apps", "api", "internal", "database", "database.go")
	if !fileExists(path) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, path, repairSQLiteBusySource)
}
