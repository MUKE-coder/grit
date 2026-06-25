package scaffold

// api_resource_stats_files.go scaffolds the v3.31.44 per-resource
// dashboard stats endpoint.
//
// What gets generated:
//   - services/resource_stats_dispatch.go — marker-driven switch over
//     resourceName -> reflective helper that computes Total + 30-day
//     sparkline + Latest N rows for any GORM model.
//   - handlers/resource_stats.go          — admin handler at
//     GET /api/admin/dashboard/resource-stats/:resource. Accepts the
//     same date params the DateFilter component sends
//     (created_since=7d | created_from + created_to).
//
// Generator wiring: each `grit generate resource Foo` injects a new
// `case "foos": ...` at the `// grit:resource-stats:dispatch` marker.
// Pre-v3.31.44 projects don't have the marker; the generator prints a
// one-line warning instead of failing.

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeResourceStatsFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "services", "resource_stats_dispatch.go"): resourceStatsDispatchGo(),
		filepath.Join(apiRoot, "internal", "handlers", "resource_stats.go"):          resourceStatsHandlerGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// resourceStatsDispatchGo emits the marker-driven dispatcher used by
// the dashboard. Every model registered here gets a "Total + sparkline"
// stat card and a "Latest N" table on the admin dashboard.
//
// The reflective helper assumes the model has a `created_at` column
// (every Grit model does, via the framework's BaseColumns convention).
// Models without `created_at` will produce empty series.
func resourceStatsDispatchGo() string {
	return `package services

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// ResourceStats is the shape returned by the dashboard's per-resource
// stats endpoint. Total is the row count within the active date range
// (or all-time when no range is given); Series is one bucket per day
// for the sparkline (always the last 30 days regardless of the date
// filter, so the sparkline never collapses to a single bar); Latest is
// up to ` + "`LatestLimit`" + ` rows, newest first.
type ResourceStats struct {
	Resource string                   ` + "`json:\"resource\"`" + `
	Total    int64                    ` + "`json:\"total\"`" + `
	Series   []ResourceStatsBucket    ` + "`json:\"series\"`" + `
	Latest   []map[string]interface{} ` + "`json:\"latest\"`" + `
}

// ResourceStatsBucket is one day in the 30-day sparkline. Date is
// always YYYY-MM-DD; Count is the number of rows whose ` + "`created_at`" + `
// fell within that calendar day (local time of the server).
type ResourceStatsBucket struct {
	Date  string ` + "`json:\"date\"`" + `
	Count int64  ` + "`json:\"count\"`" + `
}

// ResourceStatsFilter mirrors the query params the DateFilter
// component sends. Either Since (relative duration like "7d") or
// From/To (explicit YYYY-MM-DD bounds) -- never both. LatestLimit
// caps the Latest slice; defaults to 10 when zero.
type ResourceStatsFilter struct {
	Since       string // "1d" | "7d" | "30d"; empty = all-time
	From        string // "YYYY-MM-DD"; empty = no lower bound
	To          string // "YYYY-MM-DD"; empty = no upper bound
	LatestLimit int    // default 10
}

// ComputeResourceStats dispatches to the right model based on
// resourceName. Generator-driven: each ` + "`grit generate resource Foo`" + `
// appends a ` + "`case \"foos\": ...`" + ` entry below the marker.
//
// The dispatch is also the security boundary -- only the whitelisted
// resources are reachable through the dashboard endpoint, so a
// compromised admin token can't dump arbitrary tables by guessing
// resource names.
func ComputeResourceStats(db *gorm.DB, resourceName string, filter ResourceStatsFilter) (*ResourceStats, error) {
	if filter.LatestLimit <= 0 {
		filter.LatestLimit = 10
	}
	switch resourceName {
	case "users":
		return reflectiveResourceStats(db, resourceName, &models.User{}, filter)
	case "blogs":
		return reflectiveResourceStats(db, resourceName, &models.Blog{}, filter)
	// grit:resource-stats:dispatch
	default:
		return nil, fmt.Errorf("dashboard stats not registered for %q", resourceName)
	}
}

// reflectiveResourceStats computes Total + 30-day sparkline + Latest
// using one GORM scope per piece. The model parameter must be a
// pointer to a struct (e.g. ` + "`&models.Product{}`" + `); the function
// reflects on it to find the slice type for the Latest query.
func reflectiveResourceStats(db *gorm.DB, resourceName string, model interface{}, filter ResourceStatsFilter) (*ResourceStats, error) {
	fromT, toT := parseStatsRange(filter)

	// ── Total ────────────────────────────────────────────────────
	// Counts rows within the active date range; no range = count all
	// rows (the "All time" preset on the DateFilter).
	countQ := db.Model(model)
	if !fromT.IsZero() {
		countQ = countQ.Where("created_at >= ?", fromT)
	}
	if !toT.IsZero() {
		countQ = countQ.Where("created_at <= ?", toT)
	}
	var total int64
	if err := countQ.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("counting %s: %w", resourceName, err)
	}

	// ── Series ───────────────────────────────────────────────────
	// Always the last 30 days, regardless of the active date filter.
	// The sparkline is meant to show "what's the trend like?", not
	// "what's inside the filter?" -- which would collapse to a flat
	// line under the "Today" preset.
	series, err := buildDailySeries(db, model)
	if err != nil {
		return nil, fmt.Errorf("building series for %s: %w", resourceName, err)
	}

	// ── Latest ───────────────────────────────────────────────────
	// Newest-first, scoped to the active date range. Reflects on the
	// model type to build a typed slice -- this lets GORM scan
	// associations correctly, and lets us JSON-marshal each row to
	// honor the model's json tags (notably json:"-" on PasswordHash).
	latestQ := db.Model(model).Order("created_at DESC").Limit(filter.LatestLimit)
	if !fromT.IsZero() {
		latestQ = latestQ.Where("created_at >= ?", fromT)
	}
	if !toT.IsZero() {
		latestQ = latestQ.Where("created_at <= ?", toT)
	}

	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	sliceType := reflect.SliceOf(modelType)
	slicePtr := reflect.New(sliceType)
	if err := latestQ.Find(slicePtr.Interface()).Error; err != nil {
		return nil, fmt.Errorf("listing latest %s: %w", resourceName, err)
	}
	latest := sanitiseLatest(slicePtr.Elem())

	return &ResourceStats{
		Resource: resourceName,
		Total:    total,
		Series:   series,
		Latest:   latest,
	}, nil
}

// parseStatsRange turns the filter's Since/From/To into concrete time
// bounds. Returns (zero, zero) when no filter is set -- callers
// should treat zero bounds as "no constraint" and skip the WHERE
// clause entirely.
func parseStatsRange(filter ResourceStatsFilter) (time.Time, time.Time) {
	now := time.Now()
	if filter.Since != "" {
		switch filter.Since {
		case "1d":
			return now.AddDate(0, 0, -1), time.Time{}
		case "7d":
			return now.AddDate(0, 0, -7), time.Time{}
		case "30d":
			return now.AddDate(0, 0, -30), time.Time{}
		}
		return time.Time{}, time.Time{}
	}
	var from, to time.Time
	if filter.From != "" {
		if t, err := time.Parse("2006-01-02", filter.From); err == nil {
			from = t
		}
	}
	if filter.To != "" {
		if t, err := time.Parse("2006-01-02", filter.To); err == nil {
			// Include the entire ` + "`To`" + ` day -- the DateFilter passes
			// YYYY-MM-DD which would otherwise truncate at 00:00:00.
			to = t.Add(24 * time.Hour).Add(-time.Second)
		}
	}
	return from, to
}

// buildDailySeries returns 30 buckets, one per calendar day for the
// last 30 days. Uses a pull-then-bucket approach (one SELECT for the
// raw timestamps, then in-memory grouping) rather than a SQL
// GROUP BY DATE() because the DATE() function's return type differs
// between SQLite and Postgres -- this way we get identical behavior
// across both supported drivers.
//
// For very high-volume tables this would be slow; v3.31.44 ships the
// simple version. Optimization (server-side aggregation with driver-
// specific syntax) is a follow-up if the page gets heavy.
func buildDailySeries(db *gorm.DB, model interface{}) ([]ResourceStatsBucket, error) {
	cutoff := time.Now().AddDate(0, 0, -29) // 30 days inclusive of today
	cutoff = time.Date(cutoff.Year(), cutoff.Month(), cutoff.Day(), 0, 0, 0, 0, cutoff.Location())

	type row struct {
		CreatedAt time.Time
	}
	var rows []row
	err := db.Model(model).
		Select("created_at").
		Where("created_at >= ?", cutoff).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	bucketMap := make(map[string]int64, 30)
	for _, r := range rows {
		key := r.CreatedAt.Format("2006-01-02")
		bucketMap[key]++
	}

	series := make([]ResourceStatsBucket, 30)
	for i := 0; i < 30; i++ {
		d := cutoff.AddDate(0, 0, i)
		key := d.Format("2006-01-02")
		series[i] = ResourceStatsBucket{Date: key, Count: bucketMap[key]}
	}
	return series, nil
}

// sanitiseLatest walks a reflect.Value of []ModelType and turns each
// element into a map[string]interface{} via JSON round-trip. The JSON
// step honors json:"-" tags, so columns like PasswordHash never reach
// the response. Returns an empty slice (never nil) so the JSON
// response always carries an array, not null.
func sanitiseLatest(slice reflect.Value) []map[string]interface{} {
	n := slice.Len()
	out := make([]map[string]interface{}, 0, n)
	for i := 0; i < n; i++ {
		item := slice.Index(i).Interface()
		b, err := json.Marshal(item)
		if err != nil {
			continue
		}
		var m map[string]interface{}
		if err := json.Unmarshal(b, &m); err != nil {
			continue
		}
		out = append(out, m)
	}
	return out
}
`
}

// resourceStatsHandlerGo emits the admin HTTP handler. Single GET
// endpoint reads ?created_since|created_from|created_to and delegates
// to services.ComputeResourceStats.
func resourceStatsHandlerGo() string {
	return `package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/services"
)

// ResourceStatsHandler exposes per-resource dashboard stats. One
// endpoint per resource, mounted at
// GET /api/admin/dashboard/resource-stats/:resource.
//
// The handler is intentionally thin -- date-range parsing and the
// per-resource dispatch live in services.ComputeResourceStats so the
// generator only has to inject one line per new resource (a switch
// case) without touching the HTTP layer.
type ResourceStatsHandler struct {
	DB *gorm.DB
}

// Get returns the dashboard stat bundle for one resource. Accepts the
// same query params the DateFilter component already sends to the
// resource list endpoint, so the dashboard widgets and the per-page
// stats agree on the wire shape.
//
// Query params:
//
//	?created_since=7d                 — relative window
//	?created_from=2026-01-01&created_to=2026-01-31
//	                                  — explicit range
//	?limit=10                         — latest-rows cap (default 10)
func (h *ResourceStatsHandler) Get(c *gin.Context) {
	resource := c.Param("resource")
	if resource == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": "resource is required"},
		})
		return
	}

	limit := 10
	if raw := c.Query("limit"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	filter := services.ResourceStatsFilter{
		Since:       c.Query("created_since"),
		From:        c.Query("created_from"),
		To:          c.Query("created_to"),
		LatestLimit: limit,
	}

	stats, err := services.ComputeResourceStats(h.DB, resource, filter)
	if err != nil {
		// 404 (unknown resource) and 500 (DB error) look similar at
		// this layer -- the dispatcher returns a typed sentinel error
		// for the former in a future release. For now both map to 400
		// so the dashboard widget can render a friendly error state
		// without crashing.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "STATS_FAILED", "message": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
`
}
