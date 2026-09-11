package scaffold

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// writeClusterFiles brings an upgraded project's multi-replica coordination up
// to date: the cluster package, the permission cache that uses it, and a cron
// scheduler that runs on one replica only.
//
// Run behind a load balancer, three pieces of per-process state disagreed:
// a permission revoked through one replica kept working on the others, an SSO
// connection created on one was unknown to the rest, and every replica ran
// every scheduled job. Upgrade does not regenerate API code in general; this is
// the exception because the first of those is a security hole.
//
// Runs before the SSO files are written (they import the cluster package), and
// ensureClusterWiring runs after them.
func writeClusterFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	write := func(path, content string) error {
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", opts.Module())); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
		return nil
	}

	pkg := filepath.Join(apiRoot, "internal", "cluster", "cluster.go")
	if err := write(pkg, apiClusterGo()); err != nil {
		return err
	}
	if err := write(filepath.Join(apiRoot, "internal", "cluster", "cluster_test.go"), apiClusterTestGo()); err != nil {
		return err
	}
	if !fileContains(pkg, "func Bump(") {
		fmt.Println("  ⚠ internal/cluster has been edited, so replicas still keep separate permission caches.")
		return nil
	}

	if fileExists(filepath.Join(apiRoot, "internal", "authz", "grants.go")) {
		if err := write(filepath.Join(apiRoot, "internal", "authz", "grants.go"), authzGrantsGo()); err != nil {
			return err
		}
	}
	if fileExists(filepath.Join(apiRoot, "internal", "cron", "cron.go")) {
		if err := repairCronLeader(root, apiRoot, opts.Module()); err != nil {
			fmt.Printf("  ⚠ %v\n", err)
		}
	}
	return nil
}

const (
	clusterAuthzShare = "\t// Permission caches are per process. Share makes a role change on one\n" +
		"\t// replica reach every other within a second.\n" +
		"\tauthz.Share(db)\n"
	clusterSSOWatch = "\t// Rebuild the SSO registries when another replica changes a connection.\n" +
		"\tservices.WatchSSO(db, ssoRegistry, samlRegistry)\n"
)

// ensureClusterWiring adds the two start-up calls to an existing routes.go,
// each only once the function it calls is present, so a held-back file never
// leaves routes.go calling something undefined.
func ensureClusterWiring(apiRoot, module string) error {
	path := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	crlf := strings.Contains(string(raw), "\r\n")
	content := strings.ReplaceAll(string(raw), "\r\n", "\n")
	changed := false

	if fileContains(filepath.Join(apiRoot, "internal", "authz", "grants.go"), "func Share(") &&
		!strings.Contains(content, "authz.Share(db)") {
		const anchor = "\troleHandler := handlers.NewRoleHandler(db)\n"
		if strings.Contains(content, anchor) {
			content = strings.Replace(content, anchor, anchor+clusterAuthzShare, 1)
			var ok bool
			if content, ok = addImportGroup(content, module+"/internal/authz"); ok {
				changed = true
			}
		} else {
			fmt.Printf("  ⚠ Could not find where to share the permission cache. Add this to routes.go\n"+
				"    once the database is connected:\n\n%s", clusterAuthzShare)
		}
	}
	if fileContains(filepath.Join(apiRoot, "internal", "services", "sso.go"), "func WatchSSO(") &&
		!strings.Contains(content, "services.WatchSSO(") {
		const anchor = "\tssoHandler := handlers.NewSSOHandler(db, authService, cfg, ssoRegistry, samlRegistry)\n"
		if strings.Contains(content, anchor) {
			content = strings.Replace(content, anchor, anchor+clusterSSOWatch, 1)
			changed = true
		} else {
			fmt.Printf("  ⚠ Could not find the SSO registries in routes.go. Add this after they are built:\n\n%s",
				clusterSSOWatch)
		}
	}
	if !changed {
		return nil
	}
	if _, err := format.Source([]byte(content)); err != nil {
		return fmt.Errorf("routes.go left alone: wiring the replicas did not produce valid Go (%v)", err)
	}
	if crlf {
		content = strings.ReplaceAll(content, "\n", "\r\n")
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	manifest.Refresh(path)
	fmt.Println("  ✓ Replicas now share permission and SSO changes")
	return nil
}

// repairCronLeader puts the scheduler behind the cron lock in an existing
// project. cron.go is not rewritten, because grit add job writes into it, so
// the four generated lines that change are patched where they are still
// recognisable, and leader.go is added beside it.
func repairCronLeader(root, apiRoot, module string) error {
	if err := writeFile(filepath.Join(apiRoot, "internal", "cron", "leader.go"),
		strings.ReplaceAll(cronLeaderGo(), "{{MODULE}}", module)); err != nil {
		return err
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, filepath.Join(apiRoot, "internal", "cron", "cron.go"), repairCronSource)
}

var cronPatches = []struct{ old, new string }{
	{"type Scheduler struct {\n\tscheduler *asynq.Scheduler\n}",
		"type Scheduler struct {\n\tscheduler *asynq.Scheduler\n\tredisURL  string\n}"},
	{"\treturn &Scheduler{scheduler: scheduler}, nil",
		"\treturn &Scheduler{scheduler: scheduler, redisURL: redisURL}, nil"},
	{"// Start begins executing scheduled tasks.\nfunc (s *Scheduler) Start() error {\n\tgo func() {\n\t\tif err := s.scheduler.Run(); err != nil {\n\t\t\tlog.Printf(\"Cron scheduler error: %v\", err)\n\t\t}\n\t}()\n\treturn nil\n}",
		"// Start runs the scheduler on whichever replica holds the cron lock; see\n// leader.go. Every replica used to run its own, so each job ran once per replica.\nfunc (s *Scheduler) Start() error {\n\tgo s.lead(s.redisURL)\n\treturn nil\n}"},
	{"func (s *Scheduler) Stop() {\n\ts.scheduler.Shutdown()\n}",
		"func (s *Scheduler) Stop() {\n\ts.shutdown()\n}"},
}

func repairCronSource(src string) (string, []string, []string) {
	if strings.Contains(src, "go s.lead(") {
		return src, nil, nil
	}
	out := src
	for _, p := range cronPatches {
		if !strings.Contains(out, p.old) {
			return src, nil, []string{"cron.go is too changed to patch, so every replica still runs the scheduled jobs. " +
				"Start the scheduler with go s.lead(redisURL) and stop it with s.shutdown(); see leader.go"}
		}
		out = strings.Replace(out, p.old, p.new, 1)
	}
	// Start was the only user of log.
	if !strings.Contains(out, "log.") {
		out = strings.Replace(out, "\t\"log\"\n", "", 1)
	}
	return out, []string{"scheduled jobs run on one replica"}, nil
}

func apiClusterGo() string {
	return `// Package cluster is what a set of API replicas has to agree on: a generation
// number per piece of per-process state, bumped by whichever replica changes
// the data behind it and watched by the others.
//
// A permission cache, an SSO provider registry: each is built in one process,
// and each was right only in the process that changed it. Behind a load
// balancer, a role revoked through one replica kept working on the rest until
// they restarted. A tiny table in the database the replicas already share is
// the whole mechanism: no Redis, no message bus, the same on Postgres, MySQL
// and SQLite.
//
// It imports nothing from the app, so any package can use it.
package cluster

import (
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Generation is one counter, a row per watched thing.
type Generation struct {
	Name      string ` + "`" + `gorm:"primaryKey;size:64"` + "`" + `
	Value     int64  ` + "`" + `gorm:"not null;default:0"` + "`" + `
	UpdatedAt time.Time
}

func (Generation) TableName() string { return "cluster_generations" }

// Migrate creates the table. NewWatch and Bump call it, so no migration has to
// know about it. Bumps are rare, a role or a connection changing, so checking
// the schema on each is cheaper than getting a once-per-process guard wrong.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Generation{})
}

// Bump moves name's generation on, telling every other replica its copy of
// whatever name stands for is stale.
func Bump(db *gorm.DB, name string) error {
	if err := Migrate(db); err != nil {
		return err
	}
	res := db.Model(&Generation{}).Where("name = ?", name).Update("value", gorm.Expr("value + 1"))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		return nil
	}
	// The first bump creates the row. Two replicas doing that at once is fine:
	// one insert wins, and either way the value has moved off zero.
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&Generation{Name: name, Value: 1}).Error
}

// Current is name's generation, or 0 before its first bump.
func Current(db *gorm.DB, name string) (int64, error) {
	var g Generation
	err := db.Where("name = ?", name).Limit(1).Find(&g).Error
	return g.Value, err
}

// Watch reports when name's generation has moved. It reads the database at
// most once per interval, and only one caller at a time does the reading, so
// it can sit on a request's hot path.
type Watch struct {
	db       *gorm.DB
	name     string
	interval time.Duration

	mu   sync.Mutex
	next time.Time
	seen int64
}

// NewWatch starts watching name from its current generation.
func NewWatch(db *gorm.DB, name string, interval time.Duration) *Watch {
	w := &Watch{db: db, name: name, interval: interval, next: time.Now().Add(interval)}
	if err := Migrate(db); err == nil {
		w.seen, _ = Current(db, name)
	}
	return w
}

// Changed reports whether the generation has moved since the last call that
// reported a change. A failed read reports none: the local copy is exactly as
// fresh as it was.
func (w *Watch) Changed() bool {
	w.mu.Lock()
	if time.Now().Before(w.next) {
		w.mu.Unlock()
		return false
	}
	w.next = time.Now().Add(w.interval)
	w.mu.Unlock()

	value, err := Current(w.db, w.name)

	w.mu.Lock()
	defer w.mu.Unlock()
	if err != nil || value == w.seen {
		return false
	}
	w.seen = value
	return true
}
`
}

func apiClusterTestGo() string {
	return `package cluster

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func openCluster(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	return db
}

// A bump on one replica is what tells the others their copy is stale.
func TestABumpIsSeenOnce(t *testing.T) {
	db := openCluster(t)
	w := NewWatch(db, "authz", 0)
	if w.Changed() {
		t.Fatal("nothing was bumped, and the watch reported a change")
	}
	if err := Bump(db, "authz"); err != nil {
		t.Fatalf("bump: %v", err)
	}
	if !w.Changed() {
		t.Fatal("a bump was not seen")
	}
	if w.Changed() {
		t.Fatal("the same bump was reported twice")
	}
	if err := Bump(db, "sso"); err != nil {
		t.Fatalf("bump: %v", err)
	}
	if w.Changed() {
		t.Fatal("another name's bump was reported")
	}
}

// Changed sits on the request path, so it must not read on every call.
func TestAWatchReadsAtMostOncePerInterval(t *testing.T) {
	db := openCluster(t)
	w := NewWatch(db, "authz", time.Hour)
	if err := Bump(db, "authz"); err != nil {
		t.Fatalf("bump: %v", err)
	}
	if w.Changed() {
		t.Fatal("the watch read the database before its interval was up")
	}
}
`
}

func cronLeaderGo() string {
	return `package cron

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

// Only one replica runs the scheduler.
//
// asynq's schedulers do not coordinate, and every API replica built and ran
// one, so each scheduled job ran once per replica: two replicas, two token
// sweeps an hour, two orphan-upload cleanups a night. The replica holding a
// Redis lock runs it. The lock is a lease: a replica that dies stops renewing
// it, and another takes over when it lapses.
const (
	leaderKey   = "grit:cron:leader"
	leaderLease = 30 * time.Second
	leaderRenew = 10 * time.Second
)

// renewLease extends the lock only while this replica still holds it.
var renewLease = redis.NewScript(` + "`" + `if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("PEXPIRE", KEYS[1], ARGV[2]) else return 0 end` + "`" + `)

var stopOnce sync.Once

// lead waits for the cron lock, then runs the scheduler for as long as it
// holds it.
func (s *Scheduler) lead(redisURL string) {
	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		log.Printf("[cron] %v", err)
		return
	}
	client, ok := opt.MakeRedisClient().(redis.UniversalClient)
	if !ok {
		// A connection asynq understands that go-redis does not expose here:
		// better to run unlocked than not at all.
		log.Println("[cron] could not take the cron lock; running the scheduler on this replica")
		if err := s.scheduler.Start(); err != nil {
			log.Printf("[cron] %v", err)
		}
		return
	}

	id := replicaID()
	ctx := context.Background()
	for {
		held, err := client.SetNX(ctx, leaderKey, id, leaderLease).Result()
		if err == nil && held {
			break
		}
		time.Sleep(leaderRenew)
	}
	if err := s.scheduler.Start(); err != nil {
		log.Printf("[cron] %v", err)
		return
	}
	log.Println("[cron] this replica holds the cron lock and runs the scheduled jobs")

	t := time.NewTicker(leaderRenew)
	defer t.Stop()
	for range t.C {
		n, err := renewLease.Run(ctx, client, []string{leaderKey}, id, leaderLease.Milliseconds()).Int()
		if err != nil {
			// Redis blinked. The lease has time left in it; try again next tick.
			log.Printf("[cron] renewing the cron lock: %v", err)
			continue
		}
		if n == 0 {
			// This replica stalled past its lease and another took the lock.
			// Stop, so each job still runs once.
			log.Println("[cron] another replica holds the cron lock now; this one stops scheduling")
			s.shutdown()
			return
		}
	}
}

// shutdown stops the scheduler once. asynq's Shutdown closes a channel, and a
// second call would panic.
func (s *Scheduler) shutdown() {
	stopOnce.Do(s.scheduler.Shutdown)
}

func replicaID() string {
	host, _ := os.Hostname()
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return host + ":" + strconv.Itoa(os.Getpid()) + ":" + hex.EncodeToString(b)
}
`
}
