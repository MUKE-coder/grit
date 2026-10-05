package services

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
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
	X interface{} `json:"x"`
	Y float64     `json:"y"`
}

// ChartResult is the response body.
type ChartResult struct {
	Preset string                 `json:"preset"`
	Rows   []ChartRow             `json:"rows"`
	Meta   map[string]interface{} `json:"meta"`
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
	case "options":
		return reflectiveChart(db, &models.Option{}, params)
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

// countOverTime returns daily counts for the last 30 UTC days, counted by the
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

// aggOverTime runs SUM or AVG of a numeric field per UTC day over the last 30
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
