// Package paginate provides a generic list/sort/search/paginate helper
// used by every resource's List endpoint. The goal: one source of truth
// for page clamping, sort whitelisting, and search construction so that
// new resources don't drift on the boilerplate. Works with any GORM model.
//
// Usage (handler side):
//
//	func (h *ShopHandler) List(c *gin.Context) {
//	    res, err := paginate.List[models.Shop](
//	        h.DB.WithContext(c.Request.Context()).Model(&models.Shop{}).Preload("Building"),
//	        paginate.Bind(c),
//	        paginate.Config{
//	            Searchable:   []string{"shop_number", "description"},
//	            Sortable:     map[string]bool{"created_at": true, "monthly_rent": true},
//	            DefaultSort:  "created_at",
//	            DefaultOrder: "desc",
//	        },
//	    )
//	    if err != nil {
//	        respond.WriteError(c, err, "Failed to fetch shops")
//	        return
//	    }
//	    c.JSON(http.StatusOK, res)
//	}
package paginate

import (
	"encoding/base64"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Defaults applied when the request query is empty or out of range.
const (
	DefaultPage       = 1
	DefaultPageSize   = 20
	MaxPageSize       = 100
	DefaultSortColumn = "created_at"
	DefaultSortOrder  = "desc"
)

// Params is the normalized query state for a List request.
// Produced by Bind(c). Filters is free-form extra WHERE col = val clauses.
// Cursor (when present) drives cursor-mode pagination — see Config.CursorMode.
type Params struct {
	Page      int
	PageSize  int
	Search    string
	SortBy    string
	SortOrder string
	Cursor    string // opaque base64 from a previous Result.Meta.NextCursor

	// CursorMode asks for keyset pagination on this request.
	//
	// Set by ?mode=cursor, and implied by sending a ?cursor= at all, so a
	// client that follows meta.next_cursor never has to say so twice. Offset
	// stays the default because the admin table needs page numbers and a
	// total, which keyset pagination cannot give you.
	CursorMode bool
	Filters    map[string]any
	// QueryFilters are raw query-string params that are NOT reserved words.
	// Untrusted: kept apart from Filters (which handlers set in Go) so the
	// whitelist can be applied to these and only these.
	QueryFilters map[string]string

	// v3.31.34 — date filter. DateField is the column (default
	// "created_at"); DateFrom/DateTo are inclusive bounds. When both
	// are zero values, no date filter is applied. Set via the
	// ?created_from=/?created_to= query params (or the legacy
	// ?created_since=Nd shortcut used by the stat cards).
	DateField string
	DateFrom  time.Time
	DateTo    time.Time

	// Counts names extra totals to return beside the page, from ?counts=:
	// created_7d is rows created in the last 7 days, updated_30d rows updated
	// in the last 30. Each is one COUNT over the same filters and search as
	// Total. The admin's stat cards ask for theirs on the list request rather
	// than sending a request per card.
	Counts []string
}

// With returns a copy of Params with an additional filter applied.
// Empty string values are ignored so handlers can pipe c.Query() directly.
//
//	paginate.Bind(c).With("building_id", c.Query("building_id"))
func (p Params) With(key string, value any) Params {
	if s, ok := value.(string); ok && s == "" {
		return p
	}
	if value == nil {
		return p
	}
	if p.Filters == nil {
		p.Filters = map[string]any{key: value}
		return p
	}
	// Copy the map so we don't mutate the caller's Params.
	copied := make(map[string]any, len(p.Filters)+1)
	for k, v := range p.Filters {
		copied[k] = v
	}
	copied[key] = value
	p.Filters = copied
	return p
}

// Config describes which columns the caller has declared searchable / sortable
// for a particular resource. Anything not in Sortable falls back to DefaultSort.
type Config struct {
	Searchable []string        // columns included in case-insensitive search
	Sortable   map[string]bool // whitelist for sort_by values

	// Filterable whitelists columns that may be filtered from the query
	// string, so ?status=pending becomes WHERE status = 'pending'.
	//
	// A whitelist and not a free-for-all, because the column name is
	// interpolated into the SQL: without it, ?"1=1 OR x"= would be a query
	// the caller wrote. Anything not listed here is ignored rather than
	// rejected, so an unknown param is never an error.
	Filterable map[string]bool

	// RangeFilterable whitelists numeric columns that accept a window from
	// the query string: ?price_min=50&price_max=200 becomes
	// WHERE price >= 50 AND price <= 200.
	//
	// Separate from Filterable because the two answer different questions and
	// a column often wants one and not the other. Equality on a price is
	// almost never what a caller means; a range on a status is meaningless.
	// Same whitelist reasoning as above: the column name reaches the WHERE
	// clause, and the bound is parameterised.
	//
	// Either bound may be omitted. A bound that does not parse as a number is
	// ignored rather than rejected, so a hand-typed URL degrades to a wider
	// result set instead of an error page.
	RangeFilterable map[string]bool

	// InFilterable whitelists columns that accept a comma-separated list:
	// ?category_id=a,b,c becomes WHERE category_id IN (a,b,c).
	//
	// Only for id columns, and that is a deliberate restriction rather than an
	// accident of naming. Splitting on commas is wrong for anything a person
	// types, because "Smith, John" is one value; it is safe for ids, where a
	// comma never appears. So this is opt-in per column and never inferred from
	// the value.
	//
	// The case it exists for: a category tree. "Products in Electronics" means
	// Electronics and every category under it, which is a list of ids the client
	// already has from the category it fetched.
	InFilterable map[string]bool

	DefaultSort  string // fallback sort column (defaults to "created_at")
	DefaultOrder string // fallback sort order (defaults to "desc")

	// CursorMode opts into cursor-based pagination (default is offset/page).
	// When true, the response carries Meta.NextCursor + Meta.HasMore instead
	// of Page/Pages/Total. Cursor is opaque base64 of (sort_value, id) so
	// pages stay stable when rows insert mid-pagination.
	CursorMode bool

	// IncludeTotal asks cursor mode to also run COUNT(*). Slow on big
	// tables — leave off unless your UI shows a "X of Y" indicator.
	IncludeTotal bool
}

// Meta is the pagination envelope, matching Grit's existing response shape.
// Cursor mode populates NextCursor + HasMore; offset mode populates
// Page + Pages. Total is shared (always set in offset mode; opt-in in
// cursor mode via Config.IncludeTotal).
type Meta struct {
	// No omitempty on the counts. Zero is an answer: an empty result set that
	// reports {"page":1,"page_size":20} and no total leaves every client doing
	// meta.total with undefined, which renders as a blank stat card rather
	// than a nought and turns arithmetic into NaN.
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	Pages      int    `json:"pages"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more,omitempty"`

	// Mode is "cursor" on a keyset response and absent on an offset one.
	//
	// It is here because total, page and pages are all zero in cursor mode,
	// and a client cannot otherwise tell that apart from an empty table. The
	// counts are zero because counting the whole set on every page is the
	// cost keyset pagination exists to avoid; ask for it with
	// Config.IncludeTotal when you genuinely need it.
	Mode string `json:"mode,omitempty"`

	// Counts answers ?counts=, keyed by the names asked for. See Params.Counts.
	Counts map[string]int64 `json:"counts,omitempty"`
}

// Result wraps the paginated data in the canonical { data, meta } envelope.
type Result[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

// coerceFilterValue turns a query-string value into something the column can
// actually be compared against.
//
// Only booleans need this, and only because the two databases disagree about
// what to do with a string. Postgres reads WHERE active = 'true' as a boolean
// and answers correctly; MySQL stores the column as tinyint(1), coerces the
// non-numeric string to 0, and quietly returns nothing. Neither errors, so the
// bug is a filter that silently matches no rows.
//
// Numbers stay strings on purpose: both databases coerce those identically,
// and a varchar column holding "12" would break if this second-guessed it.
func coerceFilterValue(val string, dataType schema.DataType) any {
	if dataType != schema.Bool {
		return val
	}
	switch strings.ToLower(val) {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	}
	return val
}

// Bind reads page / page_size / search / sort_by / sort_order from the Gin
// context, clamps them, and returns a normalized Params.
//
// v3.31.34 — also parses the date-filter query params:
//
//	?created_from=2026-01-01&created_to=2026-12-31
//	?created_since=7d   (legacy shortcut: last 7 days)
//	?date_field=published_at   (override the default "created_at" column)
//
// Both _from and _to are inclusive. Dates without time components are
// snapped to the start (00:00) for _from and end (23:59:59) for _to.
func Bind(c *gin.Context) Params {
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(DefaultPage)))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(DefaultPageSize)))

	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}

	dateField := c.Query("date_field")
	if dateField == "" {
		dateField = "created_at"
	}
	dateFrom, dateTo := parseDateRange(c)

	return Params{
		Page:         page,
		PageSize:     pageSize,
		Search:       c.Query("search"),
		SortBy:       c.Query("sort_by"),
		SortOrder:    c.Query("sort_order"),
		Cursor:       c.Query("cursor"),
		CursorMode:   c.Query("mode") == "cursor" || c.Query("cursor") != "",
		DateField:    dateField,
		DateFrom:     dateFrom,
		DateTo:       dateTo,
		QueryFilters: collectQueryFilters(c),
		Counts:       parseCounts(c.Query("counts")),
	}
}

// reservedParams are the query keys pagination owns. Everything else is a
// candidate column filter, to be checked against Config.Filterable before it
// is used.
var reservedParams = map[string]bool{
	"page": true, "page_size": true, "search": true,
	"sort_by": true, "sort_order": true, "cursor": true,
	"date_field": true, "date_from": true, "date_to": true,
	"mode":          true,
	"created_since": true, "created_from": true, "created_to": true,
	"updated_since": true, "archived": true, "format": true,
	"counts": true,
}

func collectQueryFilters(c *gin.Context) map[string]string {
	out := map[string]string{}
	for key, values := range c.Request.URL.Query() {
		if reservedParams[key] || len(values) == 0 || values[0] == "" {
			continue
		}
		out[key] = values[0]
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// parseDateRange reads the three supported date-window query params and
// returns the resolved (from, to) bounds. Zero values mean "no bound".
//
// Precedence: explicit created_from/created_to wins over created_since.
// This lets a UI date picker override a stat card's "last 7 days" link
// without surprising clobber.
func parseDateRange(c *gin.Context) (time.Time, time.Time) {
	var from, to time.Time
	if since := c.Query("created_since"); since != "" {
		if d, ok := parseRelativeDuration(since); ok {
			from = time.Now().Add(-d)
		}
	}
	if s := c.Query("created_from"); s != "" {
		if t, err := parseDateInput(s, false); err == nil {
			from = t
		}
	}
	if s := c.Query("created_to"); s != "" {
		if t, err := parseDateInput(s, true); err == nil {
			to = t
		}
	}
	return from, to
}

// parseRelativeDuration parses "7d", "30d", "12h", "1w" into a
// time.Duration. Used by the stat-card shortcut ?created_since=7d so
// hand-written URLs stay short. Returns ok=false on unrecognised input
// (caller falls back to no bound rather than failing the request).
func parseRelativeDuration(s string) (time.Duration, bool) {
	if len(s) < 2 {
		return 0, false
	}
	unit := s[len(s)-1]
	nStr := s[:len(s)-1]
	n, err := strconv.Atoi(nStr)
	if err != nil || n < 0 {
		return 0, false
	}
	switch unit {
	case 'h':
		return time.Duration(n) * time.Hour, true
	case 'd':
		return time.Duration(n) * 24 * time.Hour, true
	case 'w':
		return time.Duration(n) * 7 * 24 * time.Hour, true
	case 'm':
		// month = 30 days. Good enough for stats; calendar-accurate
		// month math isn't worth the dep.
		return time.Duration(n) * 30 * 24 * time.Hour, true
	}
	return 0, false
}

// parseDateInput parses an ISO date or datetime string. If endOfDay is
// true and the input is a bare date, it snaps to 23:59:59.999 so the
// _to bound is inclusive of the whole day the user picked.
func parseDateInput(s string, endOfDay bool) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		if endOfDay {
			return t.Add(24*time.Hour - time.Nanosecond), nil
		}
		return t, nil
	}
	return time.Time{}, fmt.Errorf("invalid date %q", s)
}

// isSafeDateColumn reports whether the client-supplied date_field column is
// safe to interpolate into a WHERE fragment. Only bare identifiers that are
// either a known timestamp column or an explicitly-declared sortable column
// are allowed; everything else is rejected so the caller falls back to
// "created_at". This is the guard behind the date-filter injection fix.
func isSafeDateColumn(col string, cfg Config) bool {
	if col == "" {
		return false
	}
	// Reject anything that isn't a plain snake_case identifier up front —
	// no spaces, parens, quotes, or SQL operators can survive this.
	for _, r := range col {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_') {
			return false
		}
	}
	if col == "created_at" || col == "updated_at" {
		return true
	}
	return cfg.Sortable[col]
}

// List runs the query with search / sort / filters / pagination applied and
// returns a typed Result[T]. The caller is expected to have already set the
// model and any relevant Preload() chains on the passed-in *gorm.DB.
//
// Invariants enforced:
//   - page >= 1, 1 <= page_size <= MaxPageSize
//   - sort_by must be in cfg.Sortable, else cfg.DefaultSort (or DefaultSortColumn)
//   - sort_order must be "asc" or "desc", else cfg.DefaultOrder (or DefaultSortOrder)
//   - search is applied case-insensitively across cfg.Searchable columns (nothing if empty)
func List[T any](query *gorm.DB, p Params, cfg Config) (Result[T], error) {
	// Normalize sort_by against the whitelist.
	sortBy := p.SortBy
	if sortBy == "" || !cfg.Sortable[sortBy] {
		sortBy = cfg.DefaultSort
		if sortBy == "" {
			sortBy = DefaultSortColumn
		}
	}

	// Normalize sort_order.
	sortOrder := p.SortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = cfg.DefaultOrder
		if sortOrder == "" {
			sortOrder = DefaultSortOrder
		}
	}

	// Equality filters set by the handler in Go. Trusted: the column names
	// are literals in our own source.
	for col, val := range p.Filters {
		query = query.Where(col+" = ?", val)
	}

	// Equality filters from the query string (?status=pending). Untrusted, so
	// every column is checked against the whitelist before it reaches the SQL.
	// This is what makes the admin's filter dropdowns and tab strips work:
	// before it existed, Bind never collected them and the comment above
	// promised a feature the code did not have.
	if len(p.QueryFilters) > 0 {
		// The model's schema is parsed once so each value can be bound as the
		// column's real type. See coerceFilterValue for why that matters.
		var zero T
		stmt := &gorm.Statement{DB: query}
		parsed := stmt.Parse(&zero) == nil

		for col, val := range p.QueryFilters {
			if !cfg.Filterable[col] {
				continue
			}
			var bound any = val
			if parsed && stmt.Schema != nil {
				if f := stmt.Schema.LookUpField(col); f != nil {
					bound = coerceFilterValue(val, f.DataType)
				}
			}
			query = query.Where(col+" = ?", bound)
		}

		// Id lists. Same untrusted map, same whitelist rule.
		for col := range cfg.InFilterable {
			raw, ok := p.QueryFilters[col]
			if !ok || raw == "" {
				continue
			}
			parts := strings.Split(raw, ",")
			values := make([]string, 0, len(parts))
			for _, part := range parts {
				if trimmed := strings.TrimSpace(part); trimmed != "" {
					values = append(values, trimmed)
				}
			}
			if len(values) == 0 {
				continue
			}
			query = query.Where(col+" IN ?", values)
		}

		// Range windows. Read from the same untrusted QueryFilters map, so a
		// column has to be declared RangeFilterable before "price_min" can
		// become a WHERE on price.
		for col := range cfg.RangeFilterable {
			for suffix, op := range map[string]string{"_min": ">=", "_max": "<="} {
				raw, ok := p.QueryFilters[col+suffix]
				if !ok || raw == "" {
					continue
				}
				bound, err := strconv.ParseFloat(raw, 64)
				if err != nil {
					// A bound nobody can parse widens the window rather than
					// failing the request. ?price_min=cheap is a typo, not an
					// attack, and an error page is a worse answer than results.
					continue
				}
				query = query.Where(col+" "+op+" ?", bound)
			}
		}
	}

	// v3.31.34 — date-window filter. DateField defaults to "created_at"
	// in Bind() so we only need to apply when at least one bound is set.
	//
	// SECURITY (v3.31.84): DateField arrives straight from the ?date_field=
	// query param, so it MUST be whitelisted before it touches a WHERE
	// fragment — GORM treats the condition string as raw SQL. We allow the
	// two always-present timestamp columns plus anything the resource
	// already declared sortable (a strict, developer-controlled set), and
	// fall back to "created_at" on anything else. This closes the
	// date_field SQL-injection vector reachable by any authenticated user.
	dateField := p.DateField
	if !isSafeDateColumn(dateField, cfg) {
		dateField = "created_at"
	}
	if !p.DateFrom.IsZero() {
		query = query.Where(dateField+" >= ?", p.DateFrom)
	}
	if !p.DateTo.IsZero() {
		query = query.Where(dateField+" <= ?", p.DateTo)
	}

	// Apply search across configured columns.
	if p.Search != "" && len(cfg.Searchable) > 0 {
		clause, args := buildSearchClause(cfg.Searchable, p.Search)
		query = query.Where(clause, args...)
	}

	// Either the handler insists, or the request asked.
	if cfg.CursorMode || p.CursorMode {
		return listCursor[T](query, p, cfg, sortBy, sortOrder)
	}

	var result Result[T]

	// Count first (before Order/Offset/Limit so Count reflects the whole match).
	if err := countTotal(query, &result.Meta.Total); err != nil {
		return result, err
	}
	if len(p.Counts) > 0 {
		counts, err := countWindows(query, p.Counts)
		if err != nil {
			return result, err
		}
		result.Meta.Counts = counts
	}

	// Then fetch the page.
	offset := (p.Page - 1) * p.PageSize
	if err := query.
		Order(sortBy + " " + sortOrder).
		Offset(offset).
		Limit(p.PageSize).
		Find(&result.Data).Error; err != nil {
		return result, err
	}

	result.Meta.Page = p.Page
	result.Meta.PageSize = p.PageSize
	result.Meta.Pages = 0
	if result.Meta.Total > 0 && p.PageSize > 0 {
		result.Meta.Pages = int(math.Ceil(float64(result.Meta.Total) / float64(p.PageSize)))
	}

	return result, nil
}

// maxCounts caps ?counts=, so one request cannot ask for any number of COUNTs.
const maxCounts = 8

// parseCounts keeps the well-formed names from ?counts=: created_<N>d or
// updated_<N>d, with N from 1 to 3660. Anything else is dropped rather than
// failing the list; a stat card asking for a count it cannot have shows a dash.
func parseCounts(raw string) []string {
	if raw == "" {
		return nil
	}
	var out []string
	seen := map[string]bool{}
	for _, name := range strings.Split(raw, ",") {
		name = strings.TrimSpace(name)
		if _, _, ok := countWindow(name); !ok || seen[name] {
			continue
		}
		seen[name] = true
		out = append(out, name)
		if len(out) == maxCounts {
			break
		}
	}
	return out
}

// countWindow turns created_7d into the column it filters and how far back.
// The column is one of two literals, never text from the request.
func countWindow(name string) (string, time.Duration, bool) {
	var column, rest string
	switch {
	case strings.HasPrefix(name, "created_"):
		column, rest = "created_at", strings.TrimPrefix(name, "created_")
	case strings.HasPrefix(name, "updated_"):
		column, rest = "updated_at", strings.TrimPrefix(name, "updated_")
	default:
		return "", 0, false
	}
	if !strings.HasSuffix(rest, "d") {
		return "", 0, false
	}
	days, err := strconv.Atoi(strings.TrimSuffix(rest, "d"))
	if err != nil || days < 1 || days > 3660 {
		return "", 0, false
	}
	return column, time.Duration(days) * 24 * time.Hour, true
}

// countWindows runs one COUNT per name over query, which carries the list's
// filters, search and date window but no order or page.
func countWindows(query *gorm.DB, names []string) (map[string]int64, error) {
	// To the minute, so the same request moments later builds the same SQL and
	// the count cache can answer it. Compared in local time, as GORM writes
	// created_at: SQLite compares times as text, and a UTC bound against local
	// timestamps put the window off by the machine's UTC offset.
	now := time.Now().UTC().Truncate(time.Minute)
	out := make(map[string]int64, len(names))
	for _, name := range names {
		column, window, ok := countWindow(name)
		if !ok {
			continue
		}
		var n int64
		if err := countTotal(query.Session(&gorm.Session{}).Where(column+" >= ?", now.Add(-window).In(time.Local)), &n); err != nil {
			return nil, fmt.Errorf("counting %s: %w", name, err)
		}
		out[name] = n
	}
	return out, nil
}

// Counts answers ?counts= for a handler that builds its own list query instead
// of calling List. query is that list's query with its filters and search, before
// order and paging. It returns nil when the request asked for nothing.
func Counts(c *gin.Context, query *gorm.DB) (map[string]int64, error) {
	names := parseCounts(c.Query("counts"))
	if len(names) == 0 {
		return nil, nil
	}
	return countWindows(query, names)
}

// listCursor implements cursor-based pagination. The cursor is an
// opaque base64 of (sort_value, id) so pages stay stable when rows
// insert mid-pagination. We fetch PageSize+1 rows to detect HasMore
// without a separate count query.
func listCursor[T any](query *gorm.DB, p Params, cfg Config, sortBy, sortOrder string) (Result[T], error) {
	var result Result[T]

	if cfg.IncludeTotal {
		countQuery := query.Session(&gorm.Session{})
		if err := countQuery.Count(&result.Meta.Total).Error; err != nil {
			return result, err
		}
	}

	if p.Cursor != "" {
		sortVal, lastID, err := decodeCursor(p.Cursor)
		if err == nil {
			op := "<"
			if sortOrder == "asc" {
				op = ">"
			}
			// Postgres tuple comparison: (sort_col, id) < (val, id).
			// Works on SQLite too. The id tiebreaker keeps the cursor
			// stable when sort_value collides on multiple rows.
			//
			// typedCursorValue matters more than it looks: handing the raw
			// string over compares a timestamp column against text, which
			// SQLite answers true for every row, and page two comes back
			// identical to page one.
			query = query.Where(fmt.Sprintf("(%s, id) %s (?, ?)", sortBy, op),
				typedCursorValue(sortVal), lastID)
		}
	}

	limit := p.PageSize + 1
	if err := query.
		Order(sortBy + " " + sortOrder).
		Order("id " + sortOrder).
		Limit(limit).
		Find(&result.Data).Error; err != nil {
		return result, err
	}

	if len(result.Data) > p.PageSize {
		result.Data = result.Data[:p.PageSize]
		result.Meta.HasMore = true
	}

	if len(result.Data) > 0 {
		last := result.Data[len(result.Data)-1]
		sortVal, id := extractCursor(last, sortBy)
		if id != "" {
			result.Meta.NextCursor = encodeCursor(sortVal, id)
		}
	}

	result.Meta.PageSize = p.PageSize
	result.Meta.Mode = "cursor"
	return result, nil
}

// typedCursorValue turns the cursor's text back into something the column can
// be compared against.
//
// The cursor is text because it travels in a URL, but the column is not: a
// timestamp compared against a string, or an integer compared against '500',
// is a different comparison in every database and the wrong one in SQLite.
// Binding the value at its own type lets the driver do what it does for every
// other parameter.
//
// The order is deliberate. RFC3339 first, because that is what extractCursor
// writes for a time and it is specific enough not to swallow anything else.
// Integers before floats, so an id-like value does not become 1.0. A string
// that is none of those was a string in the column too.
func typedCursorValue(s string) any {
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t
	}
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}

// EncodeCursor / DecodeCursor are exported for handlers that build
// custom cursors (e.g. nested resource links).
func EncodeCursor(sortValue, id string) string      { return encodeCursor(sortValue, id) }
func DecodeCursor(s string) (string, string, error) { return decodeCursor(s) }

func encodeCursor(sortVal, id string) string {
	return base64.URLEncoding.EncodeToString([]byte(sortVal + "|" + id))
}

func decodeCursor(s string) (string, string, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", "", fmt.Errorf("invalid cursor: %w", err)
	}
	parts := strings.SplitN(string(b), "|", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid cursor format")
	}
	return parts[0], parts[1], nil
}

// extractCursor reflects on the last row to pull out the sort field
// + ID. The sort field is stored as snake_case (matching the column),
// so we convert to PascalCase for the Go struct field lookup.
func extractCursor(item interface{}, sortBy string) (string, string) {
	rv := reflect.ValueOf(item)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return "", ""
	}

	idVal := rv.FieldByName("ID")
	if !idVal.IsValid() || idVal.Kind() != reflect.String {
		return "", ""
	}
	id := idVal.String()

	sortField := lookupSortField(rv, sortBy)
	if !sortField.IsValid() {
		// No field to read means no usable cursor. Returning the id alone
		// would encode an empty sort value, and the next page would compare
		// against "" and return the whole table from the top.
		return "", ""
	}

	if t, ok := sortField.Interface().(time.Time); ok {
		return t.Format(time.RFC3339Nano), id
	}
	return fmt.Sprintf("%v", sortField.Interface()), id
}

// lookupSortField finds the struct field behind a column name.
//
// Usually that is one field: "created_at" is CreatedAt. An embedded struct is
// two, because GORM flattens it with a prefix: a money field declared as
// Price money.Money becomes the columns price_amount and price_currency, and
// there is no PriceAmount field to find. So a name that does not resolve
// whole is split at each underscore and tried as a path, longest prefix first,
// which finds Price then Amount.
//
// Getting this wrong is not a crash. FieldByName returns an invalid Value, the
// cursor encodes an empty sort value, and the next page compares against ""
// and starts again from the top: an infinite scroll that repeats its first
// page forever.
func lookupSortField(rv reflect.Value, column string) reflect.Value {
	if f := rv.FieldByName(snakeToPascal(column)); f.IsValid() {
		return f
	}

	parts := strings.Split(column, "_")
	for split := len(parts) - 1; split >= 1; split-- {
		outer := rv.FieldByName(snakeToPascal(strings.Join(parts[:split], "_")))
		if !outer.IsValid() {
			continue
		}
		if outer.Kind() == reflect.Ptr {
			if outer.IsNil() {
				return reflect.Value{}
			}
			outer = outer.Elem()
		}
		if outer.Kind() != reflect.Struct {
			continue
		}
		if inner := outer.FieldByName(snakeToPascal(strings.Join(parts[split:], "_"))); inner.IsValid() {
			return inner
		}
	}
	return reflect.Value{}
}

// snakeToPascal turns "created_at" into "CreatedAt".
func snakeToPascal(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

// buildSearchClause builds "LOWER(col1) LIKE LOWER(?) OR ..." with the
// same wildcard-wrapped search term repeated as each arg.
func buildSearchClause(cols []string, term string) (string, []any) {
	clause := ""
	args := make([]any, 0, len(cols))
	wild := "%" + term + "%"
	for i, col := range cols {
		if i > 0 {
			clause += " OR "
		}
		clause += "LOWER(" + col + ") LIKE LOWER(?)"
		args = append(args, wild)
	}
	return clause, args
}
