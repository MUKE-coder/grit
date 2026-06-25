package scaffold

// api_chart_files.go scaffolds the v3.31.47 Preset Chart builder.
//
// What gets generated:
//   - services/chart_dispatch.go -- reflective helper with 4 presets
//     (count_over_time, group_by, sum_over_time, avg_over_time) and
//     a per-resource dispatch switch sharing the
//     // grit:resource-stats:dispatch marker with v3.31.44.
//   - handlers/chart.go -- the admin endpoint at
//     GET /api/admin/dashboard/chart/:resource.
//
// Reuses the same marker as resource-stats (v3.31.44) so a single
// `grit generate` injection wires both the v3.31.44 sparkline + the
// v3.31.47 chart presets for a new resource.

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeChartFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "services", "chart_dispatch.go"): chartDispatchGo(),
		filepath.Join(apiRoot, "internal", "handlers", "chart.go"):          chartHandlerGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// chartDispatchGo emits the chart builder's service file: dispatch
// + reflective helper + four preset queries + field whitelist.
func chartDispatchGo() string {
	return `package services

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// v3.31.47 -- the Preset Chart builder. Operators pick a resource,
// a preset (one of four), and the preset's required parameters; the
// frontend renders the result with Recharts. The presets cover the
// bulk of admin-dashboard needs without introducing a query plane:
//
//	count_over_time   -- no field needed; daily counts for 30 days
//	group_by          -- requires string/bool field; top-N counts
//	sum_over_time     -- requires numeric field; daily sums
//	avg_over_time     -- requires numeric field; daily averages

// ChartRow is one point in the chart. X is the dimension (date for
// time series, label for group-by); Y is the aggregated metric.
type ChartRow struct {
	X interface{} ` + "`json:\"x\"`" + `
	Y float64     ` + "`json:\"y\"`" + `
}

// ChartResult is the response body.
type ChartResult struct {
	Preset string                 ` + "`json:\"preset\"`" + `
	Rows   []ChartRow             ` + "`json:\"rows\"`" + `
	Meta   map[string]interface{} ` + "`json:\"meta\"`" + `
}

// ChartParams carries the resolved query inputs. Validated at the
// dispatcher boundary; the helper assumes everything here is safe
// against injection.
type ChartParams struct {
	Preset string
	Field  string
	Limit  int
	Grain  string
}

// ComputeChart dispatches to the right model based on resourceName.
// Same pattern as ComputeResourceStats -- adding a resource means
// adding a case below the marker. The chart service shares the
// v3.31.44 marker so one generator injection covers both.
func ComputeChart(db *gorm.DB, resourceName string, params ChartParams) (*ChartResult, error) {
	switch resourceName {
	case "users":
		return reflectiveChart(db, &models.User{}, params)
	case "blogs":
		return reflectiveChart(db, &models.Blog{}, params)
	// grit:resource-stats:dispatch
	default:
		return nil, fmt.Errorf("custom charts not registered for %q", resourceName)
	}
}

func reflectiveChart(db *gorm.DB, model interface{}, params ChartParams) (*ChartResult, error) {
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 10
	}
	if params.Grain == "" {
		params.Grain = "day"
	}

	categoricalFields, numericFields := reflectFieldsByKind(model)

	switch params.Preset {
	case "count_over_time":
		return countOverTime(db, model, params)
	case "group_by":
		if params.Field == "" {
			return nil, fmt.Errorf("group_by requires a field")
		}
		if !categoricalFields[params.Field] {
			return nil, fmt.Errorf("field %q is not a valid group-by column (must be a string or bool)", params.Field)
		}
		return groupBy(db, model, params)
	case "sum_over_time":
		if params.Field == "" {
			return nil, fmt.Errorf("sum_over_time requires a field")
		}
		if !numericFields[params.Field] {
			return nil, fmt.Errorf("field %q is not numeric (sum requires int / float)", params.Field)
		}
		return aggOverTime(db, model, params, "SUM")
	case "avg_over_time":
		if params.Field == "" {
			return nil, fmt.Errorf("avg_over_time requires a field")
		}
		if !numericFields[params.Field] {
			return nil, fmt.Errorf("field %q is not numeric (avg requires int / float)", params.Field)
		}
		return aggOverTime(db, model, params, "AVG")
	default:
		return nil, fmt.Errorf("unknown chart preset %q", params.Preset)
	}
}

// countOverTime returns daily counts for the last 30 days.
func countOverTime(db *gorm.DB, model interface{}, params ChartParams) (*ChartResult, error) {
	cutoff := time.Now().AddDate(0, 0, -29)
	cutoff = time.Date(cutoff.Year(), cutoff.Month(), cutoff.Day(), 0, 0, 0, 0, cutoff.Location())

	type row struct {
		CreatedAt time.Time
	}
	var rows []row
	if err := db.Model(model).
		Select("created_at").
		Where("created_at >= ?", cutoff).
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("count_over_time: %w", err)
	}

	bucket := map[string]int64{}
	for _, r := range rows {
		bucket[r.CreatedAt.Format("2006-01-02")]++
	}
	out := make([]ChartRow, 30)
	for i := 0; i < 30; i++ {
		d := cutoff.AddDate(0, 0, i)
		key := d.Format("2006-01-02")
		out[i] = ChartRow{X: key, Y: float64(bucket[key])}
	}
	return &ChartResult{
		Preset: params.Preset,
		Rows:   out,
		Meta:   map[string]interface{}{"grain": params.Grain},
	}, nil
}

// groupBy returns top-N counts grouped by a string/bool column. The
// field name is validated before reaching this helper, so safe to
// splice into the SQL.
func groupBy(db *gorm.DB, model interface{}, params ChartParams) (*ChartResult, error) {
	type row struct {
		X interface{}
		Y int64
	}
	var rows []row
	if err := db.Model(model).
		Select(params.Field + " as x, COUNT(*) as y").
		Group(params.Field).
		Order("y DESC").
		Limit(params.Limit).
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("group_by: %w", err)
	}
	out := make([]ChartRow, 0, len(rows))
	for _, r := range rows {
		label := r.X
		if label == nil || label == "" {
			label = "(none)"
		}
		out = append(out, ChartRow{X: label, Y: float64(r.Y)})
	}
	return &ChartResult{
		Preset: params.Preset,
		Rows:   out,
		Meta: map[string]interface{}{
			"field": params.Field,
			"limit": params.Limit,
		},
	}, nil
}

// aggOverTime runs SUM or AVG of a numeric field over time.
func aggOverTime(db *gorm.DB, model interface{}, params ChartParams, agg string) (*ChartResult, error) {
	cutoff := time.Now().AddDate(0, 0, -29)
	cutoff = time.Date(cutoff.Year(), cutoff.Month(), cutoff.Day(), 0, 0, 0, 0, cutoff.Location())

	type row struct {
		CreatedAt time.Time
		Value     float64
	}
	var rows []row
	if err := db.Model(model).
		Select("created_at, " + params.Field + " as value").
		Where("created_at >= ?", cutoff).
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("%s_over_time: %w", strings.ToLower(agg), err)
	}

	sumByDay := map[string]float64{}
	cntByDay := map[string]int{}
	for _, r := range rows {
		key := r.CreatedAt.Format("2006-01-02")
		sumByDay[key] += r.Value
		cntByDay[key]++
	}

	out := make([]ChartRow, 30)
	for i := 0; i < 30; i++ {
		d := cutoff.AddDate(0, 0, i)
		key := d.Format("2006-01-02")
		var y float64
		if agg == "SUM" {
			y = sumByDay[key]
		} else if cntByDay[key] > 0 {
			y = sumByDay[key] / float64(cntByDay[key])
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

// reflectFieldsByKind walks the model's struct fields and returns
// two sets: categorical columns (string/bool) and numeric columns.
// Framework columns and FileRef columns are filtered out.
func reflectFieldsByKind(model interface{}) (categorical map[string]bool, numerics map[string]bool) {
	categorical = map[string]bool{}
	numerics = map[string]bool{}

	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	skip := map[string]bool{
		"id":         true,
		"created_at": true,
		"updated_at": true,
		"deleted_at": true,
		"version":    true,
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		jsonTag := f.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		name := strings.Split(jsonTag, ",")[0]
		if skip[name] {
			continue
		}

		ft := f.Type
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		typeName := ft.String()
		if strings.Contains(typeName, "FileRef") {
			continue
		}

		switch ft.Kind() {
		case reflect.String, reflect.Bool:
			categorical[name] = true
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			numerics[name] = true
		}
	}
	return
}
`
}

// chartHandlerGo emits the HTTP handler. Single GET endpoint reads
// preset + field + limit + grain query params and delegates to
// services.ComputeChart.
func chartHandlerGo() string {
	return `package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/services"
)

// v3.31.47 -- ChartHandler exposes the Preset Chart builder
// endpoint. One endpoint per resource; the preset is a query param
// and the service-side dispatcher decides which aggregation to run.
//
// Mounted at GET /api/admin/dashboard/chart/:resource.
type ChartHandler struct {
	DB *gorm.DB
}

func (h *ChartHandler) Get(c *gin.Context) {
	resource := c.Param("resource")
	if resource == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": "resource is required"},
		})
		return
	}

	preset := c.Query("preset")
	if preset == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "VALIDATION_ERROR", "message": "preset is required"},
		})
		return
	}

	limit := 10
	if raw := c.Query("limit"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	params := services.ChartParams{
		Preset: preset,
		Field:  c.Query("field"),
		Limit:  limit,
		Grain:  c.Query("grain"),
	}

	result, err := services.ComputeChart(h.DB, resource, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "CHART_FAILED", "message": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
`
}
