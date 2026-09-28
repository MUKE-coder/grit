package scaffold

// Stage 0 of the scaling guide: measure before you scale.
//
// The guide's thesis is that the skill is knowing which box you need next, and
// that you cannot know without data. Everything below exists so that question
// has an answer a person can read rather than a dashboard they have to
// interpret.
//
// The measuring lives in the API, not the CLI, for one reason: these numbers
// are only true in production. A laptop's p99 and a laptop's connection count
// describe a laptop. `grit scale` reads this endpoint over the network, which
// means it can be pointed at the thing that is actually struggling.

// apiLatencyMiddlewareGo emits internal/middleware/latency.go.
func apiLatencyMiddlewareGo() string {
	return `package middleware

import (
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Request latency, kept in a ring so the numbers cost nothing.
//
// Percentiles rather than an average, because an average hides exactly the
// users you need to hear about: at a thousand requests a second, a p99 of two
// seconds is ten people every second having a bad time, and a mean of 40ms
// says everything is fine.
//
// A fixed window of the most recent requests rather than a time series: this
// answers "how is it right now", which is the question that decides whether to
// add a box. History is Pulse's job.
const latencyWindow = 4096

var latency = struct {
	sync.Mutex
	ms    [latencyWindow]float64
	next  int
	count int
	since time.Time
	total int64
}{since: time.Now()}

// Latency records how long each request took.
func Latency() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		ms := float64(time.Since(start).Microseconds()) / 1000

		latency.Lock()
		latency.ms[latency.next] = ms
		latency.next = (latency.next + 1) % latencyWindow
		if latency.count < latencyWindow {
			latency.count++
		}
		latency.total++
		latency.Unlock()
	}
}

// LatencySnapshot is what the window says right now.
type LatencySnapshot struct {
	Samples int     ` + "`" + `json:"samples"` + "`" + `
	P50     float64 ` + "`" + `json:"p50_ms"` + "`" + `
	P95     float64 ` + "`" + `json:"p95_ms"` + "`" + `
	P99     float64 ` + "`" + `json:"p99_ms"` + "`" + `
	Max     float64 ` + "`" + `json:"max_ms"` + "`" + `
	// RPS since the process started, which is an average over the whole
	// lifetime and therefore a floor rather than a peak. Peak belongs to a
	// load test, and the guide says to run one.
	RPS float64 ` + "`" + `json:"requests_per_second"` + "`" + `
}

// Latencies reports the percentiles of the recent window.
func Latencies() LatencySnapshot {
	latency.Lock()
	n := latency.count
	sample := make([]float64, n)
	copy(sample, latency.ms[:n])
	total := latency.total
	since := latency.since
	latency.Unlock()

	out := LatencySnapshot{Samples: n}
	if elapsed := time.Since(since).Seconds(); elapsed > 0 {
		out.RPS = float64(total) / elapsed
	}
	if n == 0 {
		return out
	}
	sort.Float64s(sample)
	out.P50 = percentile(sample, 0.50)
	out.P95 = percentile(sample, 0.95)
	out.P99 = percentile(sample, 0.99)
	out.Max = sample[n-1]
	return out
}

// percentile on a sorted slice, nearest-rank. Exact enough for a number whose
// job is to decide whether to buy a machine.
func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	i := int(float64(len(sorted))*p) - 1
	if i < 0 {
		i = 0
	}
	if i >= len(sorted) {
		i = len(sorted) - 1
	}
	return sorted[i]
}
`
}

// apiScaleHandlerGo emits internal/handlers/scale.go.
func apiScaleHandlerGo() string {
	return `package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"{{MODULE}}/internal/cache"
	"{{MODULE}}/internal/database"
	"{{MODULE}}/internal/middleware"
)

// ScaleHandler answers "which scaling stage am I at, and what is the one thing
// to do next".
//
// The rule the whole thing follows: never recommend the next box until the
// current one is actually hurting, and never recommend two at once. Most
// answers should be "nothing, you have headroom", and that answer is the
// feature. An admin panel that nags you to add read replicas at forty requests
// a minute has taught you to ignore it.
type ScaleHandler struct {
	DB    *gorm.DB
	Cache *cache.Cache
}

// ScaleReport is the measurement plus the verdict.
type ScaleReport struct {
	Latency middleware.LatencySnapshot ` + "`" + `json:"latency"` + "`" + `
	DB      dbFacts                    ` + "`" + `json:"database"` + "`" + `
	Cache   cacheFacts                 ` + "`" + `json:"cache"` + "`" + `
	Stage   string                     ` + "`" + `json:"stage"` + "`" + `
	Next    string                     ` + "`" + `json:"next"` + "`" + `
	Why     string                     ` + "`" + `json:"why"` + "`" + `
	Notes   []string                   ` + "`" + `json:"notes"` + "`" + `
}

type dbFacts struct {
	Driver          string      ` + "`" + `json:"driver"` + "`" + `
	OpenConnections int         ` + "`" + `json:"open_connections"` + "`" + `
	MaxConnections  int         ` + "`" + `json:"max_connections"` + "`" + `
	PoolMax         int         ` + "`" + `json:"pool_max_per_instance"` + "`" + `
	Replicas        int         ` + "`" + `json:"replicas"` + "`" + `
	ReplicationLag  float64     ` + "`" + `json:"replication_lag_seconds"` + "`" + `
	SlowQueries     []slowQuery ` + "`" + `json:"slowest_queries"` + "`" + `
	SeqScanTables   []seqScan   ` + "`" + `json:"tables_scanned_end_to_end"` + "`" + `
	// QueryStats is false when pg_stat_statements is not installed, which is
	// the default on self-hosted Postgres. Without it the slowest-query list is
	// empty and this report is half blind.
	QueryStats bool ` + "`" + `json:"query_stats_available"` + "`" + `
}

type slowQuery struct {
	Query  string  ` + "`" + `json:"query"` + "`" + `
	Calls  int64   ` + "`" + `json:"calls"` + "`" + `
	MeanMS float64 ` + "`" + `json:"mean_ms"` + "`" + `
	TotalMS float64 ` + "`" + `json:"total_ms"` + "`" + `
}

type seqScan struct {
	Table    string ` + "`" + `json:"table"` + "`" + `
	SeqScans int64  ` + "`" + `json:"seq_scans"` + "`" + `
	LiveRows int64  ` + "`" + `json:"live_rows"` + "`" + `
}

type cacheFacts struct {
	Configured bool    ` + "`" + `json:"configured"` + "`" + `
	Hits       int64   ` + "`" + `json:"hits"` + "`" + `
	Misses     int64   ` + "`" + `json:"misses"` + "`" + `
	HitRate    float64 ` + "`" + `json:"hit_rate"` + "`" + `
}

// Report is GET /api/v1/admin/scale.
func (h *ScaleHandler) Report(c *gin.Context) {
	r := ScaleReport{Latency: middleware.Latencies()}
	r.DB = h.dbFacts()
	r.Cache = h.cacheFacts()
	r.Stage, r.Next, r.Why, r.Notes = verdict(r)
	c.JSON(http.StatusOK, gin.H{"data": r})
}

func (h *ScaleHandler) dbFacts() dbFacts {
	f := dbFacts{Driver: h.DB.Dialector.Name()}
	f.PoolMax = envInt("DB_MAX_OPEN_CONNS", 25)

	if raw := strings.TrimSpace(os.Getenv("DATABASE_REPLICA_URLS")); raw != "" {
		for _, u := range strings.Split(raw, ",") {
			if strings.TrimSpace(u) != "" {
				f.Replicas++
			}
		}
		if lag, err := database.ReplicationLag(h.DB); err == nil {
			f.ReplicationLag = lag.Seconds()
		}
	}

	if f.Driver != "postgres" {
		return f
	}

	// Connections against the ceiling. This is the Stage 5 question and the
	// only one whose answer is an error rather than slowness.
	_ = h.DB.Raw("SELECT count(*) FROM pg_stat_activity").Scan(&f.OpenConnections).Error
	var maxConn string
	if err := h.DB.Raw("SELECT current_setting('max_connections')").Scan(&maxConn).Error; err == nil {
		f.MaxConnections, _ = strconv.Atoi(maxConn)
	}

	// Slowest by total time spent, which is the ranking that matters: a query
	// taking 4ms a million times costs more than one taking 2s once.
	// pg_stat_statements is an extension and may not be installed; a missing
	// one is not an error, it is a note.
	// Silenced: the extension is off by default on self-hosted Postgres, and
	// a red ERROR line on every poll for an ordinary condition is how people
	// learn to stop reading their error log. The report says so instead.
	quiet := h.DB.Session(&gorm.Session{Logger: h.DB.Logger.LogMode(gormlogger.Silent)})
	f.QueryStats = quiet.Raw(` + "`" + `
		SELECT query, calls, mean_exec_time AS mean_ms, total_exec_time AS total_ms
		FROM pg_stat_statements
		WHERE query NOT LIKE '%pg_stat_%'
		ORDER BY total_exec_time DESC
		LIMIT 5` + "`" + `).Scan(&f.SlowQueries).Error == nil

	// Tables being read end to end. A sequential scan on a big table is the
	// Stage 6 free fix: an index, not a replica.
	_ = quiet.Raw(` + "`" + `
		SELECT relname AS table, seq_scan AS seq_scans, n_live_tup AS live_rows
		FROM pg_stat_user_tables
		WHERE seq_scan > 0 AND n_live_tup > 10000
		ORDER BY seq_scan * n_live_tup DESC
		LIMIT 5` + "`" + `).Scan(&f.SeqScanTables).Error

	return f
}

func (h *ScaleHandler) cacheFacts() cacheFacts {
	f := cacheFacts{Configured: h.Cache != nil}
	f.Hits, f.Misses, f.HitRate = cache.Stats()
	return f
}

// verdict is the guide's decision tree, in order. The first thing that is
// actually breaking wins, and nothing after it is mentioned as a "next step",
// because two changes at once means not knowing which one helped.
func verdict(r ScaleReport) (stage, next, why string, notes []string) {
	db := r.DB

	if db.Driver == "sqlite" {
		return "Stage 1: one server, one database",
			"Nothing. Move to Postgres before you think about scaling.",
			"SQLite serialises writes, so every scaling question after this one has the same answer.",
			notes
	}

	// Stage 5 first, because it is the only failure that is an error rather
	// than slowness: the app returns 500s while every dashboard looks healthy.
	if db.MaxConnections > 0 {
		used := float64(db.OpenConnections) / float64(db.MaxConnections)
		if used > 0.8 {
			return "Stage 5: connection exhaustion",
				fmt.Sprintf("Put pgbouncer in front, or lower DB_MAX_OPEN_CONNS (currently %d per instance).", db.PoolMax),
				fmt.Sprintf("%d of %d connections are in use. At 100%% the API returns errors while the database looks idle.",
					db.OpenConnections, db.MaxConnections),
				notes
		}
	}

	// Stage 6 step 1: the free fix. An index before a machine, always.
	if len(db.SeqScanTables) > 0 {
		t := db.SeqScanTables[0]
		return "Stage 6: missing indexes",
			fmt.Sprintf("EXPLAIN ANALYZE the queries hitting %s, then add the index. Do this before any replica.", t.Table),
			fmt.Sprintf("%s has %d rows and has been read end to end %d times. A missing index looks exactly like a capacity problem and costs nothing to fix.",
				t.Table, t.LiveRows, t.SeqScans),
			notes
	}

	if db.Replicas > 0 && db.ReplicationLag > 5 {
		return "Stage 6: replicas are behind",
			"Find what is loading the primary, or add replica capacity. Consider routing reads back to the primary until lag recovers.",
			fmt.Sprintf("Replication lag is %.0fs. The read-your-writes window is 5s, so users are seeing stale data outside it.", db.ReplicationLag),
			notes
	}

	// Stage 7: the same expensive answer computed over and over.
	if r.Cache.Configured && r.Cache.Hits+r.Cache.Misses > 1000 && r.Cache.HitRate < 0.5 {
		return "Stage 7: the cache is not earning its keep",
			"Check what you are caching and for how long. A hit rate under 50% usually means the TTL is shorter than the gap between reads.",
			fmt.Sprintf("Hit rate is %.0f%% over %d lookups.", r.Cache.HitRate*100, r.Cache.Hits+r.Cache.Misses),
			notes
	}

	if r.Latency.Samples > 100 && r.Latency.P99 > 1000 {
		return "Stage 2 or 3: the app is the bottleneck",
			"Check CPU on the API instances. Below 70% it is not capacity, it is a slow dependency: look at the slowest queries below.",
			fmt.Sprintf("p99 is %.0fms over the last %d requests, while p50 is %.0fms. That gap is a tail, not a trend.",
				r.Latency.P99, r.Latency.Samples, r.Latency.P50),
			notes
	}

	// Nothing is breaking. Say so, plainly, and do not invent work.
	if db.Replicas == 0 {
		notes = append(notes, "No read replicas, and you should not add any until the database's CPU is pegged by reads with the indexes already right.")
	}
	if !r.Cache.Configured {
		notes = append(notes, "No cache configured. Add one when the same expensive query shows up near the top of the slowest list.")
	}
	if db.Driver == "postgres" && !db.QueryStats {
		notes = append(notes, "pg_stat_statements is not installed, so the slowest-query list is empty. CREATE EXTENSION pg_stat_statements; and add it to shared_preload_libraries.")
	}
	if r.Latency.Samples < 100 {
		notes = append(notes, "Fewer than 100 requests measured, so the percentiles are not worth much yet. Run a load test against staging.")
	}

	return "Healthy",
		"Nothing. Adding infrastructure now buys complexity and no speed.",
		fmt.Sprintf("p99 %.0fms over %d requests, %d of %d database connections in use.",
			r.Latency.P99, r.Latency.Samples, db.OpenConnections, db.MaxConnections),
		notes
}

func envInt(key string, fallback int) int {
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil && v > 0 {
		return v
	}
	return fallback
}
`
}
