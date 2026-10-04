package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"library/apps/api/internal/appendonly"
	"library/apps/api/internal/crypto"
	"library/apps/api/internal/fieldtypes"
	"library/apps/api/internal/paginate"
	"library/apps/api/internal/sanitize"
	"library/apps/api/internal/sync"
)

// Connect establishes a database connection using the provided DSN.
//
// Driver is chosen by DSN shape:
//   - "sqlite://path" or "sqlite:path"  → SQLite (file or :memory:)
//   - anything else                     → Postgres
//
// Examples:
//
//	DATABASE_URL=sqlite:./bench.db
//	DATABASE_URL=sqlite::memory:
//	DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=disable
func Connect(dsn string) (*gorm.DB, error) {
	logLevel := logger.Warn
	if os.Getenv("DB_LOG_LEVEL") == "info" {
		logLevel = logger.Info
	} else if os.Getenv("DB_LOG_LEVEL") == "silent" {
		logLevel = logger.Silent
	}
	// Two GORM knobs that get recommended a lot. Both are off by default here,
	// and both defaults were measured rather than assumed — k6, 50 concurrent
	// writers, 4 CPUs, three 20-second runs each, median req/s on inserts:
	//
	//	off / off        740     what ships
	//	PrepareStmt      738     no difference
	//	+ skip the tx  1,294     +75%
	//
	// PrepareStmt caches a prepared statement per connection so a query is
	// planned once instead of per request. On this workload it measured as
	// nothing — the cache is mutex-guarded and under concurrency the contention
	// cancels out the saved planning. It also breaks against a connection pooler
	// in transaction mode (pgbouncer, RDS Proxy), because server-side prepared
	// statements do not survive a pooler that hands each transaction a different
	// backend. No measured gain, real downsides, so: opt in with
	// DB_PREPARED_STATEMENTS=true if your own numbers disagree.
	gormCfg := &gorm.Config{
		Logger:      logger.Default.LogMode(logLevel),
		PrepareStmt: os.Getenv("DB_PREPARED_STATEMENTS") == "true",
	}

	// Skipping the default transaction is the one that actually pays — GORM
	// wraps every Create, Update and Delete in an implicit transaction, so a
	// single-row insert costs BEGIN + INSERT + COMMIT where one round trip would
	// do. Turning it off was worth 75% here.
	//
	// It is still off by default, and that is a correctness decision rather than
	// a cautious one. The resource generator emits models with relations, and
	// saving an invoice with its line items is several INSERTs that GORM's
	// implicit transaction is currently what makes atomic — the generated
	// handler does not open its own. Without it, a failure halfway through
	// leaves an invoice holding some of its lines, with nothing logged and
	// nobody the wiser until the totals stop adding up.
	//
	// So: if your resources are flat, DB_SKIP_DEFAULT_TRANSACTION=true is close
	// to free throughput. If you generate anything with line items, leave it
	// alone until the generated handlers wrap their own writes.
	if os.Getenv("DB_SKIP_DEFAULT_TRANSACTION") == "true" {
		gormCfg.SkipDefaultTransaction = true
	}

	var (
		db  *gorm.DB
		err error
	)

	switch {
	case strings.HasPrefix(dsn, "sqlite://"):
		db, err = gorm.Open(sqlite.Open(sqliteDSN(strings.TrimPrefix(dsn, "sqlite://"))), gormCfg)
	case strings.HasPrefix(dsn, "sqlite:"):
		db, err = gorm.Open(sqlite.Open(sqliteDSN(strings.TrimPrefix(dsn, "sqlite:"))), gormCfg)
	case strings.HasPrefix(dsn, "mysql://"), strings.HasPrefix(dsn, "mysql:"):
		// go-sql-driver wants "user:pass@tcp(host:port)/db", not a URL, so the
		// scheme is stripped rather than parsed. parseTime is not optional:
		// without it DATETIME columns arrive as []byte and every time.Time
		// field on every model fails to scan.
		my := strings.TrimPrefix(strings.TrimPrefix(dsn, "mysql://"), "mysql:")
		if !strings.Contains(my, "parseTime=") {
			sep := "?"
			if strings.Contains(my, "?") {
				sep = "&"
			}
			my += sep + "parseTime=true&loc=UTC"
		}
		db, err = gorm.Open(mysql.Open(my), gormCfg)
	default:
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}), gormCfg)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Map-based updates to an EncryptedString column are encrypted too. GORM only
	// runs a column type's Value() when the value already has that type, and the
	// generated update and PATCH handlers write maps of plain strings, so without
	// this the first edit to an encrypted column stored plaintext.
	if err := crypto.Install(db); err != nil {
		return nil, fmt.Errorf("installing field encryption for map updates: %w", err)
	}

	// Rich text is sanitised on its way into the database: every field tagged
	// sanitize:"html", through every write that has a model (create, update,
	// PATCH, bulk edit, CSV import, sync push, GORM Studio's row editor).
	if err := sanitize.Install(db); err != nil {
		return nil, fmt.Errorf("installing the HTML sanitiser: %w", err)
	}

	// Append-only tables refuse UPDATE and DELETE through this handle, so GORM
	// Studio's row editor, the CSV importer and every handler are bound by it,
	// not just the one service that remembered. See internal/appendonly.
	if err := appendonly.Install(db); err != nil {
		return nil, fmt.Errorf("installing append-only guard: %w", err)
	}

	// A list's total is reused across its pages until a write to the table, so
	// paging does not count the whole match every time. See internal/paginate.
	if err := paginate.Install(db); err != nil {
		return nil, fmt.Errorf("installing the list total cache: %w", err)
	}

	// A soft delete also moves updated_at, so sync pull pages on one indexed
	// column and still carries deletes. See internal/sync.
	if err := sync.Install(db); err != nil {
		return nil, fmt.Errorf("installing the sync soft-delete hook: %w", err)
	}

	// Formatted columns (email, url, domain, tel, country, color, percent,
	// rating, time, json) are checked and normalised on every write that has a
	// model: create, update, PATCH, bulk edit, CSV import and sync push. A
	// value that is not what its column says is refused with a 422 naming the
	// field. See internal/fieldtypes.
	if err := fieldtypes.Install(db); err != nil {
		return nil, fmt.Errorf("installing the field type checks: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Connection pool settings. Postgres and MySQL use every knob; SQLite is
	// handled at the bottom, and not the way this comment used to claim.
	//
	// The number that matters is the total across processes, not this one.
	// Every replica opens its own pool, and Sentinel opens a second one in each
	// (SENTINEL_DB_MAX_OPEN_CONNS, see SentinelPoolSize), all against the
	// server's max_connections: 100 on a default Postgres, 151 on MySQL, a few
	// of them reserved for superusers. The default here was 100, so one replica
	// could take the whole server and three asked for 330 connections. Size it
	// as
	//
	//	DB_MAX_OPEN_CONNS = (max_connections - reserve) / replicas - SENTINEL_DB_MAX_OPEN_CONNS
	//
	// The defaults, 25 here and 5 for Sentinel, are that sum for three replicas
	// on a default Postgres with 10 left for migrations and psql. Behind
	// pgbouncer, as docker-compose.prod.yml runs it, Postgres only ever sees
	// pgbouncer's DEFAULT_POOL_SIZE, and these bound how much one process queues
	// there instead.
	//
	// Idle defaults to Open, and that default matters more than it looks. When
	// idle is lower, a request past the idle limit returns its connection to a
	// full pool, so the connection is closed, and the next request opens a new
	// one, which makes Postgres fork a backend process. Under concurrency that
	// is a connection storm, and it surfaces as database CPU rather than as
	// anything you would think to look for in the application.
	//
	// If your queries are heavy enough to saturate the database (an unindexed
	// COUNT over a large table on every request, say) a smaller pool acts as
	// admission control and can measure faster, because queueing in the app is
	// cheaper than thrashing in Postgres. Start here, then measure.
	maxOpen := getEnvInt("DB_MAX_OPEN_CONNS", 25)
	if maxOpen < 1 {
		maxOpen = 1
	}
	maxIdle := getEnvInt("DB_MAX_IDLE_CONNS", maxOpen)
	if maxIdle < 1 || maxIdle > maxOpen {
		maxIdle = maxOpen
	}

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

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	if maxOpen == 1 && strings.HasPrefix(dsn, "sqlite:") {
		log.Println("Database pool: one connection, because SQLite writes one at a time")
	} else {
		log.Printf("Database pool: up to %d connections from this process (DB_MAX_OPEN_CONNS)", maxOpen)
	}

	log.Println("Database connected successfully")
	return db, nil
}

// sqliteDSN adds the two pragmas a SQLite file needs to survive more than one
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

// getEnvInt reads a whole-number env var. A malformed value falls back rather
// than failing the boot — a typo in DB_MAX_OPEN_CONNS should not stop the app
// from starting.
func getEnvInt(key string, fallback int) int {
	if raw := os.Getenv(key); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			return n
		}
		log.Printf("warning: %s=%q is not a number, using %d", key, raw, fallback)
	}
	return fallback
}

// SentinelPoolSize is the connection pool Sentinel opens for its own storage.
// On Postgres that is a second pool on the app's server, in every replica, so
// it counts against max_connections the same way DB_MAX_OPEN_CONNS does.
// Sentinel writes from a batching background pipeline and needs few.
func SentinelPoolSize() (maxOpen, maxIdle int) {
	maxOpen = getEnvInt("SENTINEL_DB_MAX_OPEN_CONNS", 5)
	if maxOpen < 1 {
		maxOpen = 1
	}
	maxIdle = getEnvInt("SENTINEL_DB_MAX_IDLE_CONNS", maxOpen)
	if maxIdle < 1 || maxIdle > maxOpen {
		maxIdle = maxOpen
	}
	return maxOpen, maxIdle
}
