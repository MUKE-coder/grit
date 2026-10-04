package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// UseReplicas sends reads to the replicas named in DATABASE_REPLICA_URLS.
//
// Comma-separated, and empty means no replicas, which is what almost every
// project should have for a long time. A replica is the answer to a database
// whose CPU is pegged by reads after the indexes are right; adding one before
// that is money spent on a problem you do not have yet.
//
// dbresolver routes per statement rather than per call site: SELECT goes to a
// replica, INSERT/UPDATE/DELETE go to the primary, and anything inside a
// transaction goes to the primary including its reads. That last rule is the
// one that matters. A balance check inside the transaction that debits the
// balance must not read a replica, and with this it cannot, whatever the
// handler was written to do.
//
// Returns the number of replicas in use, so the boot log can say so.
func UseReplicas(db *gorm.DB) (int, error) {
	raw := strings.TrimSpace(os.Getenv("DATABASE_REPLICA_URLS"))
	if raw == "" {
		return 0, nil
	}

	var dialectors []gorm.Dialector
	for _, url := range strings.Split(raw, ",") {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}
		// Postgres only, and deliberately: a replica is streaming replication,
		// which SQLite does not have and MySQL spells differently enough that
		// claiming support without testing it would be a lie in a config file.
		if !strings.HasPrefix(url, "postgres://") && !strings.HasPrefix(url, "postgresql://") {
			return 0, fmt.Errorf("replica %q is not a postgres:// URL; read replicas are Postgres only", redactDSN(url))
		}
		dialectors = append(dialectors, postgres.Open(url))
	}
	if len(dialectors) == 0 {
		return 0, nil
	}

	cfg := dbresolver.Config{
		Replicas: dialectors,
		// Random is the only policy the plugin ships, and it is the right one
		// here anyway: every replica is a full copy, so there is nothing to be
		// clever about, and random spreads a burst better than round robin
		// does across a fleet of API replicas that each keep their own counter.
		Policy: dbresolver.RandomPolicy{},
	}

	// The same per-connection limits the primary gets. A replica pool that is
	// unbounded while the primary's is capped is how a replica becomes the
	// thing that runs out of connections.
	plugin := dbresolver.Register(cfg).
		SetConnMaxIdleTime(10 * time.Minute).
		SetConnMaxLifetime(30 * time.Minute).
		SetMaxIdleConns(getEnvInt("DB_MAX_IDLE_CONNS", 5)).
		SetMaxOpenConns(getEnvInt("DB_MAX_OPEN_CONNS", 25))

	if err := db.Use(plugin); err != nil {
		return 0, fmt.Errorf("registering read replicas: %w", err)
	}

	log.Printf("Read replicas: %d, reads routed to them (writes and transactions stay on the primary)", len(dialectors))
	return len(dialectors), nil
}

// Primary pins a query to the primary, whatever dbresolver would have chosen.
//
// For the reads that must not be stale: a balance before a debit, a seat before
// a sale, anything whose answer decides a write. Inside a transaction this is
// already the behaviour and the call is redundant but harmless.
//
//	var acct Account
//	database.Primary(tx).First(&acct, "id = ?", id)
func Primary(db *gorm.DB) *gorm.DB {
	return db.Clauses(dbresolver.Write)
}

// Replica forces a read onto a replica even during the read-your-writes window.
//
// For the reads where staleness is the point rather than a risk: a dashboard, a
// report, a public feed. Using it on anything a person is about to act on is
// how a stale number becomes a wrong decision.
func Replica(db *gorm.DB) *gorm.DB {
	return db.Clauses(dbresolver.Read)
}

// ReplicationLag is how far the replicas are behind, in seconds.
//
// Postgres only, and it returns 0 with no error when there are no replicas or
// the database cannot answer: a monitoring helper that fails the request it is
// measuring is worse than no monitoring.
//
// Worth alerting on. A replica that is minutes behind is not a replica, it is a
// backup that answers queries, and every read-your-writes window in the app is
// too short for it.
func ReplicationLag(db *gorm.DB) (time.Duration, error) {
	if strings.TrimSpace(os.Getenv("DATABASE_REPLICA_URLS")) == "" {
		return 0, nil
	}
	var seconds float64
	err := Replica(db).Raw(
		"SELECT COALESCE(EXTRACT(EPOCH FROM (now() - pg_last_xact_replay_timestamp())), 0)",
	).Scan(&seconds).Error
	if err != nil {
		return 0, err
	}
	if seconds < 0 {
		seconds = 0
	}
	return time.Duration(seconds * float64(time.Second)), nil
}

// redactDSN keeps a password out of an error message.
func redactDSN(dsn string) string {
	at := strings.LastIndex(dsn, "@")
	scheme := strings.Index(dsn, "://")
	if at < 0 || scheme < 0 || at < scheme {
		return dsn
	}
	return dsn[:scheme+3] + "***" + dsn[at:]
}
