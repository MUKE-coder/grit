package paginate

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// What the list page's insights panel asks for, over the list's own query.
//
// The four stat cards answer "how many", which is one number per question and
// is what ?counts= already does. A chart is a different question: how many per
// month, or how many of each status, which is one GROUP BY and not eight
// COUNTs. Asking for it over the same query is the point: the chart describes
// the rows the table is showing, filters, search and date window included, so
// filtering to one category redraws the chart for that category rather than
// leaving a whole-table chart above a filtered list, which would be a lie told
// in a prominent place.
//
// Both are opt-in per request and the panel is collapsed by default, so an
// ordinary list request pays for none of this.

// Bucket is one bar: a period, and how many rows fall in it.
type Bucket struct {
	Bucket string `json:"bucket"`
	Count  int64  `json:"count"`
}

// Slice is one value of a column and how many rows hold it.
type Slice struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

const (
	// maxBuckets caps ?series=: five years of months, or two months of days.
	maxBuckets = 60
	// maxBreakdownColumns caps ?breakdown=, which is one GROUP BY each.
	maxBreakdownColumns = 3
	// maxSlices caps the values returned per column. A breakdown of a column
	// with ten thousand distinct values is not a chart, so the largest few come
	// back and the rest are the caller's to ask about another way.
	maxSlices = 12
)

// seriesSpec is a parsed ?series=created_at:month:12.
type seriesSpec struct {
	column  string
	unit    string
	buckets int
}

// parseSeries reads ?series=<column>:<unit>:<buckets>, with the unit and the
// bucket count optional.
//
// The column is one of two literals rather than anything from the request,
// because it is interpolated into the SQL below: a time bucket cannot be a
// bind parameter, so there is nowhere safe to put a column name that came from
// a caller. created_at and updated_at are the two every model has.
func parseSeries(raw string) (seriesSpec, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return seriesSpec{}, false
	}
	parts := strings.Split(raw, ":")
	spec := seriesSpec{column: parts[0], unit: "month", buckets: 12}

	switch spec.column {
	case "created_at", "updated_at":
	default:
		return seriesSpec{}, false
	}
	if len(parts) > 1 && parts[1] != "" {
		switch parts[1] {
		case "day", "week", "month":
			spec.unit = parts[1]
		default:
			return seriesSpec{}, false
		}
	}
	if len(parts) > 2 && parts[2] != "" {
		n, err := strconv.Atoi(parts[2])
		if err != nil || n < 1 {
			return seriesSpec{}, false
		}
		if n > maxBuckets {
			n = maxBuckets
		}
		spec.buckets = n
	}
	return spec, true
}

// parseBreakdown keeps the columns from ?breakdown= that the caller is already
// allowed to filter by.
//
// Filterable is the right list: a column a client may say ?status=pending
// about can hold no secret that counting its values reveals. A column outside
// it is dropped rather than refused, the way an unknown filter is, so a panel
// asking for a column this resource does not have shows nothing instead of
// failing the list behind it.
func parseBreakdown(raw string, allowed map[string]bool) []string {
	if raw == "" {
		return nil
	}
	var out []string
	seen := map[string]bool{}
	for _, name := range strings.Split(raw, ",") {
		name = strings.TrimSpace(name)
		if name == "" || seen[name] || !allowed[name] {
			continue
		}
		seen[name] = true
		out = append(out, name)
		if len(out) == maxBreakdownColumns {
			break
		}
	}
	return out
}

// seriesBuckets groups the list's rows by period.
//
// Only periods that have rows come back. A month with nothing in it is absent
// rather than zero, because the gap depends on the window the caller is
// drawing and the client knows that window; filling it in here would mean
// inventing a start date the request never gave.
func seriesBuckets(query *gorm.DB, spec seriesSpec) ([]Bucket, error) {
	expr := bucketExpr(query, spec.column, spec.unit)

	var out []Bucket
	// Newest first with a limit, so a table with ten years in it answers with
	// the twelve months asked for rather than all hundred and twenty.
	err := query.Session(&gorm.Session{}).
		Select(expr + " AS bucket, COUNT(*) AS count").
		Group("bucket").
		Order("bucket DESC").
		Limit(spec.buckets).
		Scan(&out).Error
	if err != nil {
		return nil, fmt.Errorf("grouping by %s: %w", spec.unit, err)
	}

	// Back into reading order. The query wanted the newest; a chart wants time
	// running left to right.
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out, nil
}

// bucketExpr is the SQL that turns a timestamp into the label for its period,
// in each database Grit supports.
//
// A week is the date of its Monday rather than a week number, deliberately.
// Postgres and MySQL both count ISO weeks and SQLite does not, so a week number
// means three different things on three databases, while "the Monday of that
// week" means one thing everywhere and sorts correctly as text besides.
func bucketExpr(db *gorm.DB, column, unit string) string {
	switch db.Dialector.Name() {
	case "postgres":
		switch unit {
		case "day":
			return "to_char(" + column + ", 'YYYY-MM-DD')"
		case "week":
			return "to_char(date_trunc('week', " + column + "), 'YYYY-MM-DD')"
		default:
			return "to_char(" + column + ", 'YYYY-MM')"
		}
	case "mysql":
		switch unit {
		case "day":
			return "DATE_FORMAT(" + column + ", '%Y-%m-%d')"
		case "week":
			return "DATE_FORMAT(DATE_SUB(" + column + ", INTERVAL WEEKDAY(" + column + ") DAY), '%Y-%m-%d')"
		default:
			return "DATE_FORMAT(" + column + ", '%Y-%m')"
		}
	default: // sqlite
		switch unit {
		case "day":
			return "strftime('%Y-%m-%d', " + column + ")"
		case "week":
			// Forward to Sunday, then back six days: the Monday of the ISO
			// week, which is what the other two produce.
			return "strftime('%Y-%m-%d', " + column + ", 'weekday 0', '-6 days')"
		default:
			return "strftime('%Y-%m', " + column + ")"
		}
	}
}

// breakdownSlices counts the rows per value, for each column asked for.
//
// NULL and the empty string both come back as "", which the client shows as
// "None". They are the same thing to somebody reading a chart: nobody filled
// this in.
func breakdownSlices(query *gorm.DB, columns []string) (map[string][]Slice, error) {
	textCast := "TEXT"
	if query.Dialector.Name() == "mysql" {
		// MySQL has no CAST(x AS TEXT): CHAR is its string cast, and TEXT is a
		// syntax error rather than a slow path.
		textCast = "CHAR"
	}

	out := make(map[string][]Slice, len(columns))
	for _, column := range columns {
		var slices []Slice
		err := query.Session(&gorm.Session{}).
			Select("COALESCE(CAST(" + column + " AS " + textCast + "), '') AS value, COUNT(*) AS count").
			Group("value").
			Order("count DESC").
			Limit(maxSlices).
			Scan(&slices).Error
		if err != nil {
			return nil, fmt.Errorf("counting %s: %w", column, err)
		}
		out[column] = slices
	}
	return out, nil
}
