package handlers

import (
	"net/http"
	goruntime "runtime"

	"github.com/gin-gonic/gin"

	"saas/apps/api/internal/respond"
	"saas/apps/api/internal/services"
)

type ObservabilityHandler struct {
	Bridge *services.SecObsBridge
}

// pulseOverview is /pulse/api/overview: pulse.Overview plus pulse.RouteStats.
// Latencies are nanoseconds, Go's time.Duration JSON form.
type pulseOverview struct {
	TotalRequests int64   `json:"total_requests"`
	TotalErrors   int64   `json:"total_errors"`
	ErrorRate     float64 `json:"error_rate"`
	AvgLatency    int64   `json:"avg_latency"`
	P95Latency    int64   `json:"p95_latency"`
	RPM           float64 `json:"rpm"`
	TopRoutes     []struct {
		Method       string  `json:"method"`
		Path         string  `json:"path"`
		RequestCount int64   `json:"request_count"`
		ErrorRate    float64 `json:"error_rate"`
		AvgLatency   int64   `json:"avg_latency"`
		P50Latency   int64   `json:"p50_latency"`
		P95Latency   int64   `json:"p95_latency"`
		P99Latency   int64   `json:"p99_latency"`
	} `json:"top_routes"`
}

// pulseRuntime is /pulse/api/runtime/current: a RuntimeMetric snapshot.
type pulseRuntime struct {
	HeapAlloc    uint64 `json:"heap_alloc"`
	NumGoroutine int    `json:"num_goroutine"`
	NumGC        uint32 `json:"num_gc"`
}

// pulseN1 is /pulse/api/database/n1/ranked: N1Ranking[] under "data".
type pulseN1 struct {
	Data []struct {
		Route            string  `json:"route"`
		AvgQueriesPerHit float64 `json:"avg_queries_per_hit"`
		FirstSeen        string  `json:"first_seen"`
	} `json:"data"`
}

// pulseErrors is /pulse/api/errors: ErrorRecord[] under "data".
type pulseErrors struct {
	Data []struct {
		ID           string `json:"id"`
		Route        string `json:"route"`
		ErrorMessage string `json:"error_message"`
		LastSeen     string `json:"last_seen"`
	} `json:"data"`
}

// Pulse durations come back as nanoseconds (Go's time.Duration JSON form).
func nsToMs(ns int64) float64 { return float64(ns) / 1_000_000.0 }

// Summary returns a flat performance envelope built from Pulse's overview,
// runtime, ranked N+1 and errors endpoints. Missing or zero data is fine
// — the React page renders "—" for unset numeric fields and shows the
// "No X yet" empty state for empty arrays.
func (h *ObservabilityHandler) Summary(c *gin.Context) {
	if h.Bridge == nil {
		respond.Fail(c, respond.CodePulseOff, "Pulse is not enabled")
		return
	}

	ctx := c.Request.Context()
	var up upstream

	get := func(name, path string, out interface{}) {
		up.get(name, func() error { return h.Bridge.PulseGet(ctx, path, out) })
	}

	var overview pulseOverview
	get("overview", "/pulse/api/overview?range=1h", &overview)

	var rt pulseRuntime
	get("runtime", "/pulse/api/runtime/current", &rt)

	var n1Resp pulseN1
	get("n1_detections", "/pulse/api/database/n1/ranked?range=1h&limit=10", &n1Resp)

	var errResp pulseErrors
	get("recent_errors", "/pulse/api/errors?limit=10&resolved=false", &errResp)

	// Nothing answered. A page of noughts would say this app served no
	// requests, had no errors and ran at nought milliseconds.
	if up.allFailed(4) {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": gin.H{
				"code":    "PULSE_UNAVAILABLE",
				"message": "Pulse is not answering",
				"details": gin.H{"degraded": up.names()},
			},
		})
		return
	}

	// Build slowest_routes from the overview's TopRoutes (already ranked
	// by Pulse — top requests rather than top latency, but a useful proxy).
	slowest := make([]gin.H, 0, len(overview.TopRoutes))
	for _, r := range overview.TopRoutes {
		slowest = append(slowest, gin.H{
			"route":      r.Path,
			"method":     r.Method,
			"requests":   r.RequestCount,
			"avg":        nsToMs(r.AvgLatency),
			"p95":        nsToMs(r.P95Latency),
			"p99":        nsToMs(r.P99Latency),
			"error_rate": r.ErrorRate,
		})
	}

	n1List := make([]gin.H, 0, len(n1Resp.Data))
	for _, n := range n1Resp.Data {
		n1List = append(n1List, gin.H{
			"route":       n.Route,
			"query_count": int(n.AvgQueriesPerHit),
			"first_seen":  n.FirstSeen,
		})
	}

	recentErrs := make([]gin.H, 0, len(errResp.Data))
	for _, e := range errResp.Data {
		recentErrs = append(recentErrs, gin.H{
			"id":         e.ID,
			"route":      e.Route,
			"message":    e.ErrorMessage,
			"created_at": e.LastSeen,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"latency": gin.H{
			"p50": nsToMs(overview.AvgLatency), // Pulse overview lacks p50 — use avg as proxy
			"p95": nsToMs(overview.P95Latency),
			"p99": nsToMs(overview.P95Latency), // p99 not in overview either
			"avg": nsToMs(overview.AvgLatency),
		},
		"traffic": gin.H{
			"throughput": overview.RPM / 60.0, // convert RPM → req/s
			"total":      overview.TotalRequests,
		},
		"errors": gin.H{
			"rate":        overview.ErrorRate,
			"active_open": overview.TotalErrors,
		},
		"saturation": gin.H{
			"goroutines": rt.NumGoroutine,
			"heap_mb":    float64(rt.HeapAlloc) / (1024.0 * 1024.0),
			"gc_cycles":  rt.NumGC,
			"cpu_cores":  goruntime.NumCPU(),
		},
		"slowest_routes": slowest,
		"n1_detections":  n1List,
		"recent_errors":  recentErrs,
		// The panels whose numbers are missing rather than zero.
		"degraded": up.names(),
	})
}
