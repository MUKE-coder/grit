package middleware

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
	Samples int     `json:"samples"`
	P50     float64 `json:"p50_ms"`
	P95     float64 `json:"p95_ms"`
	P99     float64 `json:"p99_ms"`
	Max     float64 `json:"max_ms"`
	// RPS since the process started, which is an average over the whole
	// lifetime and therefore a floor rather than a peak. Peak belongs to a
	// load test, and the guide says to run one.
	RPS float64 `json:"requests_per_second"`
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
