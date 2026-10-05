// Package testkit is the shared setup and assertions for this project's tests.
//
// # Why a package and not a helper per test file
//
// Every handler test needs the same four things: a database with the schema on
// it, a config, a request, and a way to say what the response should have been.
// Written out per file that is about fifteen lines of preamble before the first
// assertion, and the fifteenth copy is the one that gets a detail wrong.
//
// # The database your tests run against
//
// By default this is SQLite in memory, which is fast and needs nothing
// installed. SQLite is not what you deploy to, and the difference is not
// academic: Postgres stores timestamps to microseconds and SQLite to
// nanoseconds, so a hash chain over a timestamp verifies locally and fails in
// production. JSONB, ILIKE, numeric precision and ordering of nulls all differ
// too.
//
// So set GRIT_TEST_DATABASE_URL to a Postgres connection string and the whole
// suite runs against the real engine instead:
//
//	docker compose up -d postgres
//	GRIT_TEST_DATABASE_URL=postgres://grit:grit@localhost:5432/grit go test ./...
//
// Each call to DB gets its own schema, created and dropped around the test, so
// tests stay isolated and can run in parallel. Worth doing before a release even
// if you do not do it every day.
package testkit

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"commerce/apps/api/internal/config"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/services"
)

// DatabaseURLEnv names the variable that points the suite at a real database.
const DatabaseURLEnv = "GRIT_TEST_DATABASE_URL"

// schemaCounter makes each schema name unique within a process. The test name
// alone is not enough: subtests repeat names, and a table-driven test that runs
// the same case twice would collide with itself.
var schemaCounter atomic.Int64

// DB returns a database with this project's schema on it, and registers its
// cleanup.
//
// Every model in models.Models() is migrated, so a resource generated tomorrow
// is in the test database without anybody editing this file. That is the whole
// reason the list lives in one place.
func DB(tb testing.TB) *gorm.DB {
	tb.Helper()

	if url := strings.TrimSpace(os.Getenv(DatabaseURLEnv)); url != "" {
		return postgresDB(tb, url)
	}
	return sqliteDB(tb)
}

func sqliteDB(tb testing.TB) *gorm.DB {
	tb.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		tb.Fatalf("opening the test database: %v", err)
	}

	// One connection, and only one.
	//
	// Every connection to ":memory:" gets its OWN empty database. The pool opens
	// a second one as soon as two queries overlap, which they do whenever a
	// handler writes something on a goroutine: the next read lands on a fresh,
	// empty database and finds nothing. It passed or failed depending on timing,
	// which reads as a broken feature rather than a test-harness problem.
	sqlDB, err := db.DB()
	if err != nil {
		tb.Fatalf("reaching the connection pool: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := db.AutoMigrate(models.Models()...); err != nil {
		tb.Fatalf("migrating the test database: %v", err)
	}
	return db
}

// postgresDB gives the test its own schema on a real Postgres, and drops it
// afterwards.
//
// A schema rather than a database because creating a database needs a separate
// connection to another one and takes about a second; a schema is immediate and
// isolates just as well once search_path is set.
func postgresDB(tb testing.TB, url string) *gorm.DB {
	tb.Helper()

	schema := fmt.Sprintf("test_%d_%s", schemaCounter.Add(1), sanitize(tb.Name()))
	if len(schema) > 60 {
		schema = schema[:60] // Postgres identifiers stop at 63 bytes.
	}

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		tb.Fatalf("connecting to %s: %v\n\nUnset %s to fall back to SQLite.", DatabaseURLEnv, err, DatabaseURLEnv)
	}

	// The identifier is built here, not taken from input, and sanitize leaves
	// only letters, digits and underscores, so there is nothing to quote around.
	if err := db.Exec("CREATE SCHEMA " + schema).Error; err != nil {
		tb.Fatalf("creating schema %s: %v", schema, err)
	}
	tb.Cleanup(func() {
		_ = db.Exec("DROP SCHEMA " + schema + " CASCADE").Error
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	})

	if err := db.Exec("SET search_path TO " + schema).Error; err != nil {
		tb.Fatalf("setting search_path: %v", err)
	}
	// One connection, so every statement sees the search_path set above. A
	// second connection from the pool would default to public and migrate the
	// schema the whole team shares.
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxOpenConns(1)
	}

	if err := db.AutoMigrate(models.Models()...); err != nil {
		tb.Fatalf("migrating schema %s: %v", schema, err)
	}
	return db
}

// sanitize reduces a test name to something usable as an SQL identifier.
func sanitize(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	return b.String()
}

// Config returns a config suitable for tests: real expiry windows, a secret
// that is obviously not a secret, and development mode so the parts that check
// it behave the way a developer running the suite expects.
func Config() *config.Config {
	return &config.Config{
		AppEnv:           "development",
		JWTSecret:        "test-secret-key-for-testing-only",
		JWTAccessExpiry:  15 * time.Minute,
		JWTRefreshExpiry: 7 * 24 * time.Hour,
	}
}

// SignIn issues a usable access token for the user, the way signing in does.
//
// Not just a token pair. Every access token carries the id of the session that
// minted it, and the auth middleware refuses one whose session row is missing,
// which is what makes signing out of one device actually sign that device out.
// A test that calls GenerateTokenPair alone gets a token the middleware rejects,
// with "Invalid or expired token", which reads as a broken token rather than a
// missing fixture and costs an afternoon.
//
// The user must already exist: the middleware loads the account on every
// request, so a token for a user who was never written is refused too.
func SignIn(tb testing.TB, db *gorm.DB, auth *services.AuthService, user *models.User) string {
	tb.Helper()

	tokens, err := auth.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		tb.Fatalf("issuing tokens for %s: %v", user.Email, err)
	}
	if _, err := services.CreateSessionCtx(context.Background(), db, services.RequestMeta{
		UserAgent: "grit-testkit",
		IP:        "127.0.0.1",
	}, user.ID, tokens.RefreshToken); err != nil {
		tb.Fatalf("recording the session for %s: %v", user.Email, err)
	}
	return tokens.AccessToken
}

// Auth builds the auth service a test's middleware needs, wired to db.
//
// db is passed rather than left nil on purpose. A service with no database
// skips the session check, so a test built that way passes against a middleware
// the project does not run.
func Auth(cfg *config.Config, db *gorm.DB) *services.AuthService {
	return &services.AuthService{
		Secret:        cfg.JWTSecret,
		AccessExpiry:  cfg.JWTAccessExpiry,
		RefreshExpiry: cfg.JWTRefreshExpiry,
		DB:            db,
	}
}

// User writes a user with the role, for a test that needs somebody to be.
//
// The password is stored as given; the model hashes it on create. Tests that
// sign in over HTTP need the plaintext, which is why it is a constant here
// rather than random.
func User(tb testing.TB, db *gorm.DB, email, role string) *models.User {
	tb.Helper()

	user := &models.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     email,
		Password:  Password,
		Role:      role,
		Active:    true,
	}
	if err := db.Create(user).Error; err != nil {
		tb.Fatalf("creating %s: %v", email, err)
	}
	return user
}

// Password is what User sets, for a test that signs in over HTTP rather than
// with SignIn.
const Password = "password-for-testing"
