package scaffold

import (
	"strings"
	"testing"
)

// Stage 6: read replicas, and the two things a hand-rolled split gets wrong.
//
// dbresolver routes per statement rather than per call site, which is the whole
// reason this is one environment variable rather than an audit of every query
// in the project. The rule that matters most is the one nobody writes by hand:
// a read inside a transaction goes to the primary, so a balance check inside
// the transaction that debits the balance cannot read a replica.
func TestReadReplicasRouteWithoutTouchingHandlers(t *testing.T) {
	src := apiReplicasGo()

	if !strings.Contains(src, "dbresolver.Register(cfg)") {
		t.Fatal("replicas are not registered as a GORM plugin, so every call site would have to choose")
	}
	if !strings.Contains(src, `os.Getenv("DATABASE_REPLICA_URLS")`) {
		t.Error("replicas are not configured from the environment")
	}
	// Empty means none, and none is the right answer for almost every project.
	if !strings.Contains(src, "if raw == \"\" {\n\t\treturn 0, nil\n\t}") {
		t.Error("an unset DATABASE_REPLICA_URLS is not treated as no replicas")
	}
	// The replica pool gets the same ceiling as the primary's, or it becomes
	// the thing that runs out of connections.
	for _, want := range []string{"SetMaxOpenConns", "SetMaxIdleConns", "SetConnMaxLifetime"} {
		if !strings.Contains(src, want) {
			t.Errorf("the replica pool has no %s, so it is unbounded", want)
		}
	}
	// Forcing either side, for the reads whose answer decides a write.
	if !strings.Contains(src, "func Primary(db *gorm.DB) *gorm.DB") ||
		!strings.Contains(src, "dbresolver.Write") {
		t.Error("nothing can pin a read to the primary")
	}
	// A password must not reach an error message.
	if !strings.Contains(src, "func redactDSN(") {
		t.Error("a bad replica URL would put its password in the error")
	}
	// Postgres only, said out loud rather than failing oddly later.
	if !strings.Contains(src, "read replicas are Postgres only") {
		t.Error("a non-Postgres replica URL is not refused with a reason")
	}
}

// The bug this is really about: you post, the feed loads from a replica that
// has not caught up, and your own post is missing.
func TestReadYourOwnWritesNeedsNoSharedStore(t *testing.T) {
	mw := apiReadYourWritesGo()

	if !strings.Contains(mw, "readYourWritesCookie") {
		t.Fatal("the recent-write marker is not a cookie")
	}
	// A cookie rather than Redis, deliberately: it travels with the person who
	// wrote, costs no lookup, and cannot itself be stale. The comment in the
	// file says so, which is why this reads code rather than prose.
	for _, line := range strings.Split(mw, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}
		if strings.Contains(strings.ToLower(line), "redis") {
			t.Errorf("the marker reaches a shared store, which puts Redis on the read path: %q",
				strings.TrimSpace(line))
		}
	}
	if !strings.Contains(mw, "c.Writer.Status() < 300") {
		t.Error("a failed write still pins later reads to the primary")
	}
	if !strings.Contains(mw, "http.SameSiteLaxMode") || !strings.Contains(mw, "secure") {
		t.Error("the cookie is not set with SameSite and Secure")
	}

	helper := apiForRequestGo()
	if !strings.Contains(helper, "func ForRequest(c *gin.Context, db *gorm.DB) *gorm.DB") {
		t.Fatal("handlers have no helper to pick the right handle")
	}
	if !strings.Contains(helper, "return Primary(db)") {
		t.Error("the sticky window does not actually route to the primary")
	}
}

// Stage 7: cache-aside, with the three things a hand-written version misses.
func TestCacheRememberIsTheWholePattern(t *testing.T) {
	src := apiCacheRememberGo()

	if !strings.Contains(src, "func Remember[T any](") {
		t.Fatal("there is no cache-aside helper, so every call site reimplements it")
	}
	// 1. A cache outage makes the app slower, not broken.
	if !strings.Contains(src, "falling through to the source") {
		t.Error("a cache read failure does not fall through to the loader")
	}
	// 2. Jitter, or a class of keys written together expires together.
	if !strings.Contains(src, "func jitter(") || !strings.Contains(src, "jitter(ttl)") {
		t.Error("TTLs are not jittered, so a deploy's keys all expire in the same second")
	}
	// 3. One rebuilder per key when a hot key expires.
	if !strings.Contains(src, "SetNX") || !strings.Contains(src, "waitForRebuild") {
		t.Error("no stampede protection: every concurrent miss would run the same query")
	}
	// The value outlives the request that computed it.
	if !strings.Contains(src, "context.WithoutCancel(ctx), key, value") {
		t.Error("a client hanging up throws away the value it just paid for")
	}
	// And a number, so "is the cache working" is answerable.
	if !strings.Contains(src, "func Stats()") {
		t.Error("there is no hit rate, so the cache cannot be tuned")
	}
}

// Stage 0: measure first. Percentiles rather than an average, because an
// average hides exactly the users worth hearing about.
func TestLatencyIsMeasuredInPercentiles(t *testing.T) {
	src := apiLatencyMiddlewareGo()

	for _, want := range []string{"P50", "P95", "P99", "func Latencies()"} {
		if !strings.Contains(src, want) {
			t.Errorf("no %s", want)
		}
	}
	// A ring, so measuring costs nothing and never grows.
	if !strings.Contains(src, "latencyWindow") || !strings.Contains(src, "% latencyWindow") {
		t.Error("samples are not kept in a fixed ring, so the measurement leaks memory")
	}
	if strings.Contains(src, "append(latency.ms") {
		t.Error("samples are appended to an unbounded slice")
	}
}

// The advisor names one thing, and most of the time that thing is nothing.
func TestScaleAdvisorRecommendsOneThingAtATime(t *testing.T) {
	src := apiScaleHandlerGo()

	if !strings.Contains(src, "func verdict(r ScaleReport)") {
		t.Fatal("there is no verdict, only measurements")
	}
	// Connection exhaustion is checked before anything else: it is the only
	// failure that arrives as errors while every dashboard looks healthy.
	iConn := strings.Index(src, "Stage 5: connection exhaustion")
	iIndex := strings.Index(src, "Stage 6: missing indexes")
	if iConn < 0 || iIndex < 0 || iConn > iIndex {
		t.Error("the decision tree does not check connection exhaustion first")
	}
	// An index before a machine, always.
	if !strings.Contains(src, "Do this before any replica") {
		t.Error("missing indexes do not take precedence over adding a replica")
	}
	// The most common answer has to be "nothing", or the tool becomes noise.
	if !strings.Contains(src, "Adding infrastructure now buys complexity and no speed") {
		t.Error("a healthy system is not told to do nothing")
	}
	// Admin only: it reports connection counts and query shapes.
	if !strings.Contains(apiScaleHandlerGo(), "ScaleHandler") {
		t.Error("no handler")
	}
}

// Every stage this now covers has to be reachable from the environment file,
// or nobody finds it.
func TestScalingIsDocumentedWhereItIsConfigured(t *testing.T) {
	env := envCloudExampleFile(Options{ProjectName: "shop", Architecture: ArchTriple, Frontend: FrontendNext})

	if !strings.Contains(env, "DATABASE_REPLICA_URLS=") {
		t.Fatal("the replica setting is not in .env.example, so nobody discovers it")
	}
	// And the warning that matters more than the setting.
	if !strings.Contains(env, "it should stay empty until") {
		t.Error(".env does not say that replicas are for later")
	}
}
