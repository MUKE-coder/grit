package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review M11 and M14, in the files upgrade does not deliver whole.
// M12, M13 and M15 change files upgrade already delivers, so they live in the
// templates alone.

// ─── M11: database calls carry the request's context ────────────────────────

// dropUnusedImport removes path from the import block when nothing after the
// block refers to it, reporting whether it did. A mention in a comment keeps it.
func dropUnusedImport(src, path string) (string, bool) {
	start := strings.Index(src, "\nimport (\n")
	if start < 0 {
		return src, false
	}
	end := strings.Index(src[start:], "\n)\n")
	if end < 0 {
		return src, false
	}
	end += start
	line := "\t\"" + path + "\"\n"
	block := src[start:end]
	if !strings.Contains(block+"\n", "\n"+line) {
		return src, false
	}
	name := path[strings.LastIndex(path, "/")+1:]
	if strings.Contains(src[end:], name+".") {
		return src, false
	}
	return src[:start] + strings.Replace(block+"\n", "\n"+line, "\n", 1)[:len(block)+1-len(line)] + src[end+1:], true
}

var (
	// A handler, or a method that takes the request first: resolveUser(c, conn).
	ginHandlerSig  = regexp.MustCompile(`func \(h \*\w+\) \w+\(c \*gin\.Context(?:, [^)]*)?\)[^{\n]*\{`)
	handlerDBCall  = regexp.MustCompile(`\bh\.DB\.`)
	computeCallDB  = regexp.MustCompile(`\bservices\.(ComputeChart|ComputeResourceStats)\(h\.DB,`)
	goFuncLiteral  = regexp.MustCompile(`\bgo func\(`)
	requestCtxBind = "WithContext(c.Request.Context())."
)

// bindRequestContext binds the database calls inside each gin handler method to
// the request's context, so a request the client abandoned stops its queries
// instead of holding a pool connection until they finish. Calls inside a go func
// are left alone: that work is meant to outlive the request.
func bindRequestContext(src string) (string, int) {
	var b strings.Builder
	last, n := 0, 0
	for _, loc := range ginHandlerSig.FindAllStringIndex(src, -1) {
		if loc[0] < last {
			continue
		}
		end := goBraceEnd(src, loc[1]-1)
		if end < 0 {
			continue
		}
		body := src[loc[1]:end]
		skip := append(goFuncBodies(body), commentAndStringSpans(body)...)

		var nb strings.Builder
		at := 0
		for _, m := range handlerDBCall.FindAllStringIndex(body, -1) {
			if inSpans(skip, m[0]) || strings.HasPrefix(body[m[1]:], "WithContext(") {
				continue
			}
			nb.WriteString(body[at:m[1]])
			nb.WriteString(requestCtxBind)
			at = m[1]
			n++
		}
		nb.WriteString(body[at:])
		rewritten := nb.String()

		rewritten = computeCallDB.ReplaceAllStringFunc(rewritten, func(call string) string {
			n++
			return strings.TrimSuffix(call, "h.DB,") + "h.DB.WithContext(c.Request.Context()),"
		})

		b.WriteString(src[last:loc[1]])
		b.WriteString(rewritten)
		last = end
	}
	b.WriteString(src[last:])
	return b.String(), n
}

// goBraceEnd returns the offset just past the brace that closes the one at
// open, skipping string and rune literals and line comments, or -1.
func goBraceEnd(src string, open int) int {
	depth := 0
	for i := open; i < len(src); i++ {
		switch ch := src[i]; ch {
		case '"', '\'':
			for j := i + 1; j < len(src); j++ {
				if src[j] == '\\' {
					j++
					continue
				}
				if src[j] == ch || src[j] == '\n' {
					i = j
					break
				}
			}
		case '/':
			if i+1 < len(src) && src[i+1] == '/' {
				if nl := strings.IndexByte(src[i:], '\n'); nl >= 0 {
					i += nl
				} else {
					return -1
				}
			}
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}
	return -1
}

func goFuncBodies(body string) [][2]int {
	var spans [][2]int
	for _, m := range goFuncLiteral.FindAllStringIndex(body, -1) {
		open := strings.IndexByte(body[m[1]:], '{')
		if open < 0 {
			continue
		}
		start := m[1] + open
		if end := goBraceEnd(body, start); end > 0 {
			spans = append(spans, [2]int{start, end})
		}
	}
	return spans
}

// commentAndStringSpans marks line comments and double-quoted strings, where
// h.DB is text rather than a call.
func commentAndStringSpans(src string) [][2]int {
	var spans [][2]int
	for i := 0; i < len(src); i++ {
		switch {
		case src[i] == '"':
			j := i + 1
			for j < len(src) && src[j] != '"' && src[j] != '\n' {
				if src[j] == '\\' {
					j++
				}
				j++
			}
			spans = append(spans, [2]int{i, j + 1})
			i = j
		case src[i] == '/' && i+1 < len(src) && src[i+1] == '/':
			j := strings.IndexByte(src[i:], '\n')
			if j < 0 {
				j = len(src) - i
			}
			spans = append(spans, [2]int{i, i + j})
			i += j
		}
	}
	return spans
}

func inSpans(spans [][2]int, at int) bool {
	for _, s := range spans {
		if at >= s[0] && at < s[1] {
			return true
		}
	}
	return false
}

// The auth middleware's user lookup, and the three jobs that write.
const (
	authUserLookupOld = "\t\tif err := db.Where(\"id = ?\", claims.UserID).First(&user).Error; err != nil {\n"
	authUserLookupNew = "\t\tif err := db.WithContext(c.Request.Context()).Where(\"id = ?\", claims.UserID).First(&user).Error; err != nil {\n"
)

var workerDBCalls = map[string]string{
	"deps.DB.Model(&models.Upload{})": "deps.DB.WithContext(ctx).Model(&models.Upload{})",
	"result := deps.DB.Exec(":         "result := deps.DB.WithContext(ctx).Exec(",
	"return deps.DB.Transaction(":     "return deps.DB.WithContext(ctx).Transaction(",
}

// ─── M14: stats and charts count per day in the database ────────────────────

const servicesDayBucket = `package services

import "gorm.io/gorm"

// dayExpr is a SQL expression for column's calendar day in UTC, as YYYY-MM-DD,
// in the dialect db speaks. Grouping by it counts rows per day in the database;
// the stats and chart services used to load every row of the last 30 days into
// the API and count them there.
func dayExpr(db *gorm.DB, column string) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "to_char(" + column + " AT TIME ZONE 'UTC', 'YYYY-MM-DD')"
	case "mysql":
		return "DATE_FORMAT(" + column + ", '%Y-%m-%d')"
	default: // sqlite
		return "strftime('%Y-%m-%d', " + column + ")"
	}
}

// lastThirtyDays is the first moment of the 30 UTC days that end today.
func lastThirtyDays() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -29)
}

// asStored is t written the way GORM writes created_at: the same instant, in
// local time. SQLite compares times as text, so a UTC cutoff against local
// timestamps was off by the machine's UTC offset. Postgres compares instants
// and is unaffected.
func asStored(t time.Time) time.Time {
	return t.In(time.Local)
}
`

const dailySeriesFunc = `// buildDailySeries returns 30 buckets, one per UTC calendar day for the last 30
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
`

const countOverTimeFunc = `// countOverTime returns daily counts for the last 30 UTC days, counted by the
// database rather than by loading every row.
func countOverTime(db *gorm.DB, model interface{}, params ChartParams) (*ChartResult, error) {
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
		return nil, fmt.Errorf("count_over_time: %w", err)
	}

	counts := make(map[string]int64, len(rows))
	for _, r := range rows {
		counts[r.Day] = r.N
	}
	out := make([]ChartRow, 30)
	for i := 0; i < 30; i++ {
		key := cutoff.AddDate(0, 0, i).Format("2006-01-02")
		out[i] = ChartRow{X: key, Y: float64(counts[key])}
	}
	return &ChartResult{
		Preset: params.Preset,
		Rows:   out,
		Meta:   map[string]interface{}{"grain": params.Grain},
	}, nil
}
`

const aggOverTimeFunc = `// aggOverTime runs SUM or AVG of a numeric field per UTC day over the last 30
// days, in the database. The field is validated as numeric before this runs.
func aggOverTime(db *gorm.DB, model interface{}, params ChartParams, agg string) (*ChartResult, error) {
	cutoff := lastThirtyDays()
	day := dayExpr(db, "created_at")

	type row struct {
		Day   string
		Total float64
		N     int64
	}
	var rows []row
	if err := db.Model(model).
		Select(day+" AS day, COALESCE(SUM("+params.Field+"), 0) AS total, COUNT(*) AS n").
		Where("created_at >= ?", asStored(cutoff)).
		Group(day).
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("%s_over_time: %w", strings.ToLower(agg), err)
	}

	byDay := make(map[string]row, len(rows))
	for _, r := range rows {
		byDay[r.Day] = r
	}
	out := make([]ChartRow, 30)
	for i := 0; i < 30; i++ {
		key := cutoff.AddDate(0, 0, i).Format("2006-01-02")
		r := byDay[key]
		y := r.Total
		if agg != "SUM" {
			y = 0
			if r.N > 0 {
				y = r.Total / float64(r.N)
			}
		}
		out[i] = ChartRow{X: key, Y: y}
	}
	return &ChartResult{
		Preset: params.Preset,
		Rows:   out,
		Meta: map[string]interface{}{
			"field": params.Field,
			"agg":   strings.ToLower(agg),
			"grain": params.Grain,
		},
	}, nil
}
`

// servicesDayBucketGo is internal/services/day_bucket.go. dayExpr and
// lastThirtyDays live apart from both services that use them, so either file
// can be removed or regenerated without taking the helper with it.
func servicesDayBucketGo() string {
	return strings.Replace(servicesDayBucket, `import "gorm.io/gorm"`, "import (\n\t\"time\"\n\n\t\"gorm.io/gorm\"\n)", 1)
}

// replaceGoFunc replaces the function that starts with signature, and the doc
// comment directly above it, with replacement.
func replaceGoFunc(src, signature, replacement string) (string, bool) {
	at := strings.Index(src, signature)
	if at < 0 {
		return src, false
	}
	end := goBraceEnd(src, at+len(signature)-1)
	if end < 0 {
		return src, false
	}
	if end < len(src) && src[end] == '\n' {
		end++
	}
	start := at
	for start > 0 {
		prev := strings.LastIndex(src[:start-1], "\n") + 1
		if !strings.HasPrefix(strings.TrimLeft(src[prev:start-1], "\t "), "//") {
			break
		}
		start = prev
	}
	return src[:start] + replacement + src[end:], true
}

// repairRequestPerformance applies M11 and M14 to an existing project.
func repairRequestPerformance(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	handlers, _ := filepath.Glob(filepath.Join(apiRoot, "internal", "handlers", "*.go"))
	for _, path := range handlers {
		if strings.HasSuffix(path, "_test.go") || !gritWroteFile(root, path) {
			continue
		}
		name := filepath.Base(path)
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			out, n := bindRequestContext(src)
			if n == 0 {
				return src, nil, nil
			}
			return out, []string{fmt.Sprintf("%s: %d database calls follow the request's context", name, n)}, nil
		}); err != nil {
			return err
		}
	}

	if path := filepath.Join(apiRoot, "internal", "middleware", "auth.go"); fileExists(path) {
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			if !strings.Contains(src, authUserLookupOld) {
				return src, nil, nil
			}
			return strings.Replace(src, authUserLookupOld, authUserLookupNew, 1), []string{"the auth middleware's user lookup follows the request's context"}, nil
		}); err != nil {
			return err
		}
	}

	if path := filepath.Join(apiRoot, "internal", "jobs", "workers.go"); fileExists(path) {
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			out := src
			for old, bound := range workerDBCalls {
				out = strings.Replace(out, old, bound, 1)
			}
			if out == src {
				return src, nil, nil
			}
			return out, []string{"background jobs' database calls stop when the job's context ends"}, nil
		}); err != nil {
			return err
		}
	}

	services := filepath.Join(apiRoot, "internal", "services")
	stats := filepath.Join(services, "resource_stats_dispatch.go")
	chart := filepath.Join(services, "chart_dispatch.go")
	if !fileExists(stats) && !fileExists(chart) {
		return nil
	}
	if bucket := filepath.Join(services, "day_bucket.go"); !fileExists(bucket) {
		if err := os.WriteFile(bucket, []byte(servicesDayBucketGo()), 0o644); err != nil {
			return err
		}
	}
	for path, funcs := range map[string][][2]string{
		stats: {{"func buildDailySeries(db *gorm.DB, model interface{}) ([]ResourceStatsBucket, error) {", dailySeriesFunc}},
		chart: {
			{"func countOverTime(db *gorm.DB, model interface{}, params ChartParams) (*ChartResult, error) {", countOverTimeFunc},
			{"func aggOverTime(db *gorm.DB, model interface{}, params ChartParams, agg string) (*ChartResult, error) {", aggOverTimeFunc},
		},
	} {
		if !fileExists(path) {
			continue
		}
		name := filepath.Base(path)
		fs := funcs
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			out := src
			var changes []string
			if !strings.Contains(src, "dayExpr(db") {
				for _, f := range fs {
					next, ok := replaceGoFunc(out, f[0], f[1])
					if !ok {
						return src, nil, []string{name + " is not the file Grit wrote: count rows per day with GROUP BY in SQL instead of loading every row since the cutoff"}
					}
					out = next
				}
				changes = append(changes, name+" counts rows per day in the database")
			}
			// The day counts moved the cutoff arithmetic into day_bucket.go, so a
			// file can be left importing time without using it, which does not build.
			if trimmed, ok := dropUnusedImport(out, "time"); ok {
				out = trimmed
				if len(changes) == 0 {
					changes = append(changes, name+" no longer imports time, which it stopped using")
				}
			}
			if out == src {
				return src, nil, nil
			}
			return out, changes, nil
		}); err != nil {
			return err
		}
	}
	return nil
}
