package services

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
)

// ResourceStats is the shape returned by the dashboard's per-resource
// stats endpoint. Total is the row count within the active date range
// (or all-time when no range is given); Series is one bucket per day
// for the sparkline (always the last 30 days regardless of the date
// filter, so the sparkline never collapses to a single bar); Latest is
// up to `LatestLimit` rows, newest first.
type ResourceStats struct {
	Resource string                   `json:"resource"`
	Total    int64                    `json:"total"`
	Series   []ResourceStatsBucket    `json:"series"`
	Latest   []map[string]interface{} `json:"latest"`
}

// ResourceStatsBucket is one day in the 30-day sparkline. Date is
// always YYYY-MM-DD; Count is the number of rows whose `created_at`
// fell within that calendar day (local time of the server).
type ResourceStatsBucket struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
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
// resourceName. Generator-driven: each `grit generate resource Foo`
// appends a `case "foos": ...` entry below the marker.
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
	case "collections":
		return reflectiveResourceStats(db, resourceName, &models.Collection{}, filter)

	case "products":
		return reflectiveResourceStats(db, resourceName, &models.Product{}, filter)

	case "options":
		return reflectiveResourceStats(db, resourceName, &models.Option{}, filter)
	case "pages":
		return reflectiveResourceStats(db, resourceName, &models.Page{}, filter)

	case "carts":
		return reflectiveResourceStats(db, resourceName, &models.Cart{}, filter)

	case "cart_items":
		return reflectiveResourceStats(db, resourceName, &models.CartItem{}, filter)

	// grit:resource-stats:dispatch
	default:
		return nil, fmt.Errorf("dashboard stats not registered for %q", resourceName)
	}
}

// reflectiveResourceStats computes Total + 30-day sparkline + Latest
// using one GORM scope per piece. The model parameter must be a
// pointer to a struct (e.g. `&models.Product{}`); the function
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
			// Include the entire `To` day -- the DateFilter passes
			// YYYY-MM-DD which would otherwise truncate at 00:00:00.
			to = t.Add(24 * time.Hour).Add(-time.Second)
		}
	}
	return from, to
}

// buildDailySeries returns 30 buckets, one per UTC calendar day for the last 30
// days. The database counts the rows per day; loading every timestamp to count
// them in Go cost a row per record created that month, per widget, per load.
func buildDailySeries(db *gorm.DB, model interface{}) ([]ResourceStatsBucket, error) {
	cutoff := lastThirtyDays()
	day := dayExpr(db, "created_at")

	type row struct {
		Day string
		N   int64
	}
	var rows []row
	if err := db.Model(model).
		Select(day+" AS day, COUNT(*) AS n").
		Where("created_at >= ?", asStored(cutoff)).
		Group(day).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	counts := make(map[string]int64, len(rows))
	for _, r := range rows {
		counts[r.Day] = r.N
	}
	series := make([]ResourceStatsBucket, 30)
	for i := 0; i < 30; i++ {
		key := cutoff.AddDate(0, 0, i).Format("2006-01-02")
		series[i] = ResourceStatsBucket{Date: key, Count: counts[key]}
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
