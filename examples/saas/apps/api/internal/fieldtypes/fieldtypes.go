// Package fieldtypes checks and normalises the formatted columns a generated
// model declares with a format tag:
//
//	Email string `gorm:"size:254" json:"email" format:"email"`
//	Phone string `gorm:"size:32" json:"phone" format:"tel:UG"`
//	Score int    `json:"score" format:"rating:10"`
//
// Install registers GORM callbacks that run before every create and update, so
// the rule holds on every way in rather than in whichever handler remembered:
// the generated create, update and PATCH, bulk edit, the CSV importer, sync
// push and GORM Studio's row editor. A value is either stored in its canonical
// form (an email lowercased, a domain without its scheme, a phone number in
// E.164) or the write is refused with an *Error, which the respond package
// answers with 422 and the field's message.
//
// Raw SQL and db.Table(...).Updates(...) have no model, so no field to read a
// tag from, and are not checked.
package fieldtypes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/idna"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"saas/apps/api/internal/respond"
)

// Normalizer checks one value and returns it in its stored form. param is
// what follows the colon in the tag: the default country of tel:UG, the stars
// of rating:10. A value that fails comes back with an error phrased for the
// person who typed it ("is not a valid email address").
type Normalizer func(value any, param string) (any, error)

var (
	registryMu sync.RWMutex
	registry   = map[string]Normalizer{}
)

// Register adds a format. internal/phone registers tel this way, so the
// libphonenumber metadata is only compiled into a project that has a phone
// field.
func Register(format string, fn Normalizer) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[format] = fn
}

func lookup(format string) (Normalizer, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	fn, ok := registry[format]
	return fn, ok
}

func init() {
	Register("email", stringRule(Email))
	Register("url", stringRule(URL))
	Register("domain", stringRule(Domain))
	Register("country", stringRule(Country))
	Register("color", stringRule(Color))
	Register("time", stringRule(TimeOfDay))
	Register("percent", normalizePercent)
	Register("rating", normalizeRating)
	Register("json", normalizeJSON)
}

// Error is a value refused by its column's format. It names the field, so the
// response can say which input is wrong.
type Error struct {
	Field   string
	Message string
}

func (e *Error) Error() string { return label(e.Field) + " " + e.Message }

// FieldErrors is the per-field detail the validation response carries.
func (e *Error) FieldErrors() map[string]string {
	return map[string]string{e.Field: label(e.Field) + " " + e.Message}
}

// ErrorCode answers 422 through respond.WriteError, including from a respond
// package older than the FieldErrors detail.
func (e *Error) ErrorCode() respond.Code { return respond.CodeValidationError }

// label turns phone_number into "Phone number".
func label(field string) string {
	s := strings.ReplaceAll(field, "_", " ")
	if s == "" {
		return "Value"
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// ── the rules ───────────────────────────────────────────────────────────────

// Email lowercases and trims an address and refuses anything that is not a
// bare address with a dotted domain. "Ada <ada@example.com>" is refused rather
// than quietly reduced: a display name in an email column is a form bug.
func Email(s string) (string, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	addr, err := mail.ParseAddress(s)
	if err != nil || addr.Address != s || len(s) > 254 {
		return "", errors.New("is not a valid email address")
	}
	at := strings.LastIndex(s, "@")
	if at < 1 || !strings.Contains(s[at+1:], ".") {
		return "", errors.New("is not a valid email address")
	}
	return s, nil
}

// URL accepts http and https addresses with a host, and nothing else. A
// javascript: or data: URL stored here would be a link in the admin.
func URL(s string) (string, error) {
	s = strings.TrimSpace(s)
	u, err := url.Parse(s)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || len(s) > 2048 {
		return "", errors.New("must be an http or https address, such as https://example.com")
	}
	if _, err := Domain(u.Hostname()); err != nil && !isIP(u.Hostname()) && u.Hostname() != "localhost" {
		return "", errors.New("must be an http or https address, such as https://example.com")
	}
	return s, nil
}

var ipHost = regexp.MustCompile("^[0-9.]+$|:")

func isIP(host string) bool { return ipHost.MatchString(host) }

var (
	domainLabel = regexp.MustCompile("^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$")
	topLevel    = regexp.MustCompile("^([a-z]{2,63}|xn--[a-z0-9-]+)$")
)

// Domain stores a bare host name in lowercase ASCII. A pasted scheme, path,
// port or trailing dot is removed, and an internationalised name is stored in
// its punycode form (münchen.de becomes xn--mnchen-3ya.de), which is what DNS
// resolves and what keeps a unique index honest: one domain, one spelling.
func Domain(s string) (string, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if i := strings.Index(s, "://"); i >= 0 {
		s = s[i+3:]
	}
	if i := strings.IndexAny(s, "/?#"); i >= 0 {
		s = s[:i]
	}
	if i := strings.LastIndex(s, "@"); i >= 0 {
		s = s[i+1:]
	}
	if i := strings.Index(s, ":"); i >= 0 {
		s = s[:i]
	}
	s = strings.TrimSuffix(s, ".")
	ascii, err := idna.Lookup.ToASCII(s)
	if err != nil || len(ascii) > 253 || !strings.Contains(ascii, ".") {
		return "", errors.New("is not a valid domain, such as example.com")
	}
	labels := strings.Split(ascii, ".")
	for _, l := range labels {
		if !domainLabel.MatchString(l) {
			return "", errors.New("is not a valid domain, such as example.com")
		}
	}
	if tld := labels[len(labels)-1]; !topLevel.MatchString(tld) {
		return "", errors.New("is not a valid domain, such as example.com")
	}
	return ascii, nil
}

// Country stores an ISO 3166-1 alpha-2 code in upper case.
func Country(s string) (string, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if !countrySet[s] {
		return "", errors.New("is not an ISO 3166-1 country code, such as UG")
	}
	return s, nil
}

var hexColor = regexp.MustCompile("^#?([0-9a-f]{3}|[0-9a-f]{6})$")

// Color stores #rrggbb in lower case. #abc is expanded, and a missing # is
// added.
func Color(s string) (string, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	m := hexColor.FindStringSubmatch(s)
	if m == nil {
		return "", errors.New("must be a hex colour, such as #6c5ce7")
	}
	hex := m[1]
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	return "#" + hex, nil
}

var clock = regexp.MustCompile("^([01]?[0-9]|2[0-3]):([0-5][0-9])(:[0-5][0-9](\\.[0-9]+)?)?$")

// TimeOfDay stores HH:MM on a 24-hour clock. Seconds, which an
// <input type="time"> can send, are dropped.
func TimeOfDay(s string) (string, error) {
	m := clock.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return "", errors.New("must be a time of day, such as 14:30")
	}
	hour, _ := strconv.Atoi(m[1])
	return fmt.Sprintf("%02d:%s", hour, m[2]), nil
}

// Percent checks a value from 0 to 100 and rounds it to the two decimals the
// column holds.
func Percent(v float64) (float64, error) {
	if math.IsNaN(v) || v < 0 || v > 100 {
		return 0, errors.New("must be between 0 and 100")
	}
	return math.Round(v*100) / 100, nil
}

// Rating checks a whole number of stars from 1 to max. Zero is "not rated".
func Rating(v, max int) (int, error) {
	if v < 0 || v > max {
		return 0, fmt.Errorf("must be a whole number of stars from 1 to %d", max)
	}
	return v, nil
}

// ── adapting the rules to what a write carries ──────────────────────────────

// stringRule adapts a string rule. The empty string is left alone: whether a
// field may be blank is the binding tag's business, not the format's.
func stringRule(rule func(string) (string, error)) Normalizer {
	return func(value any, _ string) (any, error) {
		switch v := value.(type) {
		case string:
			if v == "" {
				return v, nil
			}
			return rule(v)
		case *string:
			if v == nil || *v == "" {
				return v, nil
			}
			out, err := rule(*v)
			if err != nil {
				return nil, err
			}
			return &out, nil
		case nil:
			return nil, nil
		}
		return nil, errors.New("must be text")
	}
}

func toFloat(value any) (float64, bool, error) {
	switch v := value.(type) {
	case nil:
		return 0, false, nil
	case float64:
		return v, true, nil
	case float32:
		return float64(v), true, nil
	case int:
		return float64(v), true, nil
	case int64:
		return float64(v), true, nil
	case json.Number:
		f, err := v.Float64()
		return f, err == nil, err
	case string:
		if strings.TrimSpace(v) == "" {
			return 0, false, nil
		}
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		return f, err == nil, err
	case *float64:
		if v == nil {
			return 0, false, nil
		}
		return *v, true, nil
	}
	return 0, false, errors.New("not a number")
}

func normalizePercent(value any, _ string) (any, error) {
	f, ok, err := toFloat(value)
	if err != nil {
		return nil, errors.New("must be a number between 0 and 100")
	}
	if !ok {
		return value, nil
	}
	return Percent(f)
}

func normalizeRating(value any, param string) (any, error) {
	max := 5
	if n, err := strconv.Atoi(param); err == nil && n > 0 {
		max = n
	}
	f, ok, err := toFloat(value)
	if err != nil || (ok && f != math.Trunc(f)) {
		return nil, fmt.Errorf("must be a whole number of stars from 1 to %d", max)
	}
	if !ok {
		return value, nil
	}
	return Rating(int(f), max)
}

// normalizeJSON checks raw JSON and turns a decoded value (what a PATCH body
// holds) back into JSON text, so the column is always given valid JSON.
func normalizeJSON(value any, _ string) (any, error) {
	switch v := value.(type) {
	case nil:
		return nil, nil
	case datatypes.JSON:
		if len(v) == 0 {
			return v, nil
		}
		if !json.Valid(v) {
			return nil, errors.New("is not valid JSON")
		}
		return v, nil
	case []byte:
		if len(v) == 0 {
			return datatypes.JSON(nil), nil
		}
		if !json.Valid(v) {
			return nil, errors.New("is not valid JSON")
		}
		return datatypes.JSON(v), nil
	case json.RawMessage:
		if !json.Valid(v) {
			return nil, errors.New("is not valid JSON")
		}
		return datatypes.JSON(v), nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, errors.New("is not valid JSON")
	}
	return datatypes.JSON(raw), nil
}

// Normalize runs format (for example "tel:UG") over value, as a write would.
func Normalize(format string, value any) (any, error) {
	name, param, _ := strings.Cut(format, ":")
	fn, ok := lookup(name)
	if !ok {
		return nil, fmt.Errorf("no check registered for format %q", name)
	}
	return fn(value, param)
}

// ── the GORM callbacks ──────────────────────────────────────────────────────

// Install checks every formatted field on its way into the database. Call it
// once, straight after connecting.
func Install(db *gorm.DB) error {
	if err := db.Callback().Create().Before("gorm:create").Register("fieldtypes:create", check); err != nil {
		return err
	}
	return db.Callback().Update().Before("gorm:update").Register("fieldtypes:update", check)
}

type formatted struct {
	field  *schema.Field
	format string
}

var fieldCache sync.Map // *schema.Schema -> []formatted

func formattedFields(s *schema.Schema) []formatted {
	if cached, ok := fieldCache.Load(s); ok {
		return cached.([]formatted)
	}
	var out []formatted
	for _, f := range s.Fields {
		if tag := f.Tag.Get("format"); tag != "" {
			out = append(out, formatted{field: f, format: tag})
		}
	}
	fieldCache.Store(s, out)
	return out
}

func jsonName(f *schema.Field) string {
	if tag := f.Tag.Get("json"); tag != "" {
		if name, _, _ := strings.Cut(tag, ","); name != "" && name != "-" {
			return name
		}
	}
	return f.DBName
}

func check(tx *gorm.DB) {
	stmt := tx.Statement
	if stmt == nil || stmt.Schema == nil {
		return
	}
	fields := formattedFields(stmt.Schema)
	if len(fields) == 0 {
		return
	}
	// An update that names its columns in a map is checked for those columns
	// only. The row it is applied to was checked when it was written, and a
	// value stored before a rule existed should not block an unrelated edit.
	switch dest := stmt.Dest.(type) {
	case map[string]interface{}:
		checkMap(tx, fields, dest)
		return
	case *map[string]interface{}:
		if dest != nil {
			checkMap(tx, fields, *dest)
		}
		return
	}
	if stmt.Dest != nil {
		checkValue(tx, fields, reflect.ValueOf(stmt.Dest))
	}
}

func checkMap(tx *gorm.DB, fields []formatted, values map[string]interface{}) {
	for key, value := range values {
		for _, f := range fields {
			if key != f.field.DBName && key != f.field.Name && key != jsonName(f.field) {
				continue
			}
			out, err := Normalize(f.format, value)
			if err != nil {
				_ = tx.AddError(&Error{Field: jsonName(f.field), Message: err.Error()})
				return
			}
			values[key] = out
		}
	}
}

func checkValue(tx *gorm.DB, fields []formatted, v reflect.Value) {
	if !v.IsValid() {
		return
	}
	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			checkValue(tx, fields, v.Index(i))
		}
	case reflect.Struct:
		if !v.CanAddr() || v.Type() != fields[0].field.Schema.ModelType {
			return
		}
		ctx := tx.Statement.Context
		if ctx == nil {
			ctx = context.Background()
		}
		for _, f := range fields {
			value, zero := f.field.ValueOf(ctx, v)
			if zero {
				continue
			}
			out, err := Normalize(f.format, value)
			if err != nil {
				_ = tx.AddError(&Error{Field: jsonName(f.field), Message: err.Error()})
				return
			}
			if !reflect.DeepEqual(out, value) {
				if err := f.field.Set(ctx, v, out); err != nil {
					_ = tx.AddError(err)
					return
				}
			}
		}
	}
}
