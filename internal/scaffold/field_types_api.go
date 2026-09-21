package scaffold

// The API half of the formatted field types: internal/fieldtypes, which every
// API ships with, and internal/phone, which a project gets the first time it
// generates a tel field.

// PhoneNumbersModule is the Go port of Google's libphonenumber that
// internal/phone validates and formats with (MIT). Added to a project only
// when it has a tel field: the metadata it embeds is a megabyte and a half,
// which a project without phone numbers has no reason to download.
const (
	PhoneNumbersModule  = "github.com/nyaruka/phonenumbers"
	PhoneNumbersVersion = "v1.8.1"
)

// fieldTypesConnectHook installs the checks straight after connecting, beside
// the sanitiser, so every write with a model goes through them.
const fieldTypesConnectHook = `	// Formatted columns (email, url, domain, tel, country, color, percent,
	// rating, time, json) are checked and normalised on every write that has a
	// model: create, update, PATCH, bulk edit, CSV import and sync push. A
	// value that is not what its column says is refused with a 422 naming the
	// field. See internal/fieldtypes.
	if err := fieldtypes.Install(db); err != nil {
		return nil, fmt.Errorf("installing the field type checks: %w", err)
	}

`

func apiFieldTypesGo() string {
	return `// Package fieldtypes checks and normalises the formatted columns a generated
// model declares with a format tag:
//
//	Email string ` + "`" + `gorm:"size:254" json:"email" format:"email"` + "`" + `
//	Phone string ` + "`" + `gorm:"size:32" json:"phone" format:"tel:UG"` + "`" + `
//	Score int    ` + "`" + `json:"score" format:"rating:10"` + "`" + `
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

	"{{MODULE}}/internal/respond"
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
`
}

func apiFieldTypesCountriesGo() string {
	return `package fieldtypes

// CountryCodes is ISO 3166-1 alpha-2, the 249 officially assigned codes. The
// admin's country picker and the shared Zod schema carry the same list.
var CountryCodes = []string{
` + countryCodesLiteral("\t") + `
}

var countrySet = func() map[string]bool {
	m := make(map[string]bool, len(CountryCodes))
	for _, c := range CountryCodes {
		m[c] = true
	}
	return m
}()
`
}

func apiFieldTypesSamplesGo() string {
	return `package fieldtypes

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/datatypes"
)

// Sample values for seeders. Each takes the row number and mixes it in, so a
// thousand seeded rows on a unique column do not collide, and each returns a
// value its own rule accepts. The fieldtypes tests check both.

var nonLetters = regexp.MustCompile("[^a-z0-9]+")

func slug(s string) string {
	s = strings.Trim(nonLetters.ReplaceAllString(strings.ToLower(s), ""), "-")
	if s == "" {
		return "user"
	}
	return s
}

// SampleEmail builds first.last.N@domain.
func SampleEmail(first, last, domain string, n int) string {
	d, err := Domain(domain)
	if err != nil {
		d = "example.com"
	}
	return fmt.Sprintf("%s.%s.%d@%s", slug(first), slug(last), n, d)
}

var samplePaths = []string{"pricing", "about", "blog", "docs", "careers", "contact", "products", "support"}

// SampleURL builds https://domain/page/N.
func SampleURL(domain string, n int) string {
	d, err := Domain(domain)
	if err != nil {
		d = "example.com"
	}
	return fmt.Sprintf("https://%s/%s/%d", d, samplePaths[n%len(samplePaths)], n)
}

var sampleTLDs = []string{"com", "co.ug", "io", "org", "net", "dev", "co.uk", "co.ke", "app", "africa"}

// SampleDomain builds wordN.tld from a spread of top-level domains.
func SampleDomain(word string, n int) string {
	return fmt.Sprintf("%s%d.%s", slug(word), n, sampleTLDs[n%len(sampleTLDs)])
}

// SampleCountry walks the list with a stride that is coprime with its length,
// so consecutive rows land on countries all over it rather than AD, AE, AF.
func SampleCountry(n int) string {
	if n < 0 {
		n = -n
	}
	return CountryCodes[(n*97)%len(CountryCodes)]
}

// SampleColor spreads rows over the colour space. The multiplier is odd, so
// the first 16.7 million rows get different colours.
func SampleColor(n int) string {
	m := int64(n % 0x1000000)
	if m < 0 {
		m = -m
	}
	return fmt.Sprintf("#%06x", (m*2654435761)%0x1000000)
}

// SamplePercent is a value with at most one decimal, clustered below 40 the
// way discounts, tax rates and completion figures tend to be.
func SamplePercent(n int) float64 {
	v := (n*7919 + 13) % 1001 // 0..1000
	if n%3 != 0 {
		v = v * 4 / 10
	}
	return float64(v) / 10
}

// sampleStars leans towards four and five out of five, which is what real
// review distributions look like.
var sampleStars = []int{5, 4, 5, 3, 4, 5, 4, 2, 5, 4, 1, 4, 5, 3, 4, 5}

// SampleRating is a whole number of stars from 1 to max.
func SampleRating(n, max int) int {
	if max <= 0 {
		max = 5
	}
	s := sampleStars[n%len(sampleStars)]
	return (s-1)*(max-1)/4 + 1
}

// SampleTime is a time in business hours, 08:00 to 17:45 in quarter hours.
func SampleTime(n int) string {
	slot := (n * 7) % 40
	return fmt.Sprintf("%02d:%02d", 8+slot/4, (slot%4)*15)
}

var samplePlans = []string{"free", "starter", "pro", "enterprise"}

// SampleJSON is a small object, the kind of settings blob a json column holds.
func SampleJSON(n int) datatypes.JSON {
	return datatypes.JSON(fmt.Sprintf("{\"plan\":%q,\"seats\":%d,\"trial\":%t,\"tags\":[\"row-%d\"]}",
		samplePlans[n%len(samplePlans)], 1+n%50, n%4 == 0, n))
}
`
}

func apiFieldTypesTestGo() string {
	return `package fieldtypes

import (
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestRules(t *testing.T) {
	good := map[string][2]string{
		"email":   {"  Ada@Example.COM ", "ada@example.com"},
		"url":     {"https://example.com/pricing?x=1", "https://example.com/pricing?x=1"},
		"domain":  {"https://www.Example.co.ug/about", "www.example.co.ug"},
		"country": {"ug", "UG"},
		"color":   {"#ABC", "#aabbcc"},
		"time":    {"9:05:30", "09:05"},
	}
	for format, c := range good {
		got, err := Normalize(format, c[0])
		if err != nil || got != c[1] {
			t.Errorf("%s(%q) = %v, %v; want %q", format, c[0], got, err, c[1])
		}
	}
	if got, err := Normalize("domain", "münchen.de"); err != nil || got != "xn--mnchen-3ya.de" {
		t.Errorf("an IDN is stored as punycode: got %v, %v", got, err)
	}

	bad := map[string]string{
		"email":   "Ada <ada@example.com>",
		"url":     "javascript:alert(1)",
		"domain":  "not a domain",
		"country": "XX",
		"color":   "blue",
		"time":    "25:00",
		"percent": "101",
		"rating":  "4.5",
		"json":    "",
	}
	for format, v := range bad {
		var value any = v
		if format == "json" {
			value = datatypes.JSON("{nope")
		}
		if _, err := Normalize(format, value); err == nil {
			t.Errorf("%s accepted %q", format, v)
		}
	}
	if _, err := Normalize("rating:10", 7.0); err != nil {
		t.Errorf("rating:10 refused 7: %v", err)
	}
	if _, err := Normalize("rating:5", 7.0); err == nil {
		t.Error("rating:5 accepted 7")
	}
	if got, _ := Normalize("percent", 12.345); got != 12.35 {
		t.Errorf("percent rounds to two decimals: got %v", got)
	}
}

func TestSamplesPassTheirRules(t *testing.T) {
	seen := map[string]bool{}
	for n := 0; n < 5000; n++ {
		values := map[string]any{
			"email":     SampleEmail("Ada", "Lovelace", "example.com", n),
			"url":       SampleURL("example.com", n),
			"domain":    SampleDomain("acme", n),
			"country":   SampleCountry(n),
			"color":     SampleColor(n),
			"time":      SampleTime(n),
			"percent":   SamplePercent(n),
			"rating:5":  SampleRating(n, 5),
			"rating:10": SampleRating(n, 10),
			"json":      SampleJSON(n),
		}
		for format, v := range values {
			if _, err := Normalize(format, v); err != nil {
				t.Fatalf("sample %d for %s (%v) fails its own rule: %v", n, format, v, err)
			}
		}
		for _, unique := range []string{values["email"].(string), values["url"].(string), values["domain"].(string), values["color"].(string)} {
			if seen[unique] {
				t.Fatalf("sample %d repeats %q", n, unique)
			}
			seen[unique] = true
		}
	}
}

type contact struct {
	ID      uint
	Email   string         ` + "`" + `json:"email" format:"email"` + "`" + `
	Color   string         ` + "`" + `json:"color" format:"color"` + "`" + `
	Score   int            ` + "`" + `json:"score" format:"rating:5"` + "`" + `
	Meta    datatypes.JSON ` + "`" + `json:"meta" format:"json"` + "`" + `
	Comment string
}

// Every way a generated API writes a row.
func TestEveryWriteIsChecked(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := Install(db); err != nil {
		t.Fatalf("install: %v", err)
	}
	if err := db.AutoMigrate(&contact{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	row := contact{Email: "Ada@Example.com", Color: "#ABC", Score: 4}
	if err := db.Create(&row).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	var stored contact
	db.First(&stored, row.ID)
	if stored.Email != "ada@example.com" || stored.Color != "#aabbcc" {
		t.Errorf("create stored %q and %q", stored.Email, stored.Color)
	}

	var fe *Error
	err = db.Create(&contact{Email: "nope"}).Error
	if !errors.As(err, &fe) || fe.Field != "email" {
		t.Fatalf("an invalid create was not refused with the field: %v", err)
	}

	if err := db.Model(&row).Updates(map[string]interface{}{"color": "FFF"}).Error; err != nil {
		t.Fatalf("Updates(map): %v", err)
	}
	db.First(&stored, row.ID)
	if stored.Color != "#ffffff" {
		t.Errorf("Updates(map) stored %q", stored.Color)
	}
	if err := db.Model(&row).Updates(map[string]interface{}{"score": 9}).Error; !errors.As(err, &fe) {
		t.Errorf("Updates(map) accepted 9 stars out of 5: %v", err)
	}
	if err := db.Model(&row).Update("email", "bad").Error; !errors.As(err, &fe) {
		t.Errorf("Update(column) accepted a bad email: %v", err)
	}
	// A PATCH body holds decoded JSON, which is stored as JSON text.
	if err := db.Model(&row).Updates(map[string]interface{}{"meta": map[string]interface{}{"plan": "pro"}}).Error; err != nil {
		t.Fatalf("Updates(map) with a decoded object: %v", err)
	}
	db.First(&stored, row.ID)
	if string(stored.Meta) != "{\"plan\":\"pro\"}" {
		t.Errorf("json stored %s", stored.Meta)
	}
	// An untagged column is not touched.
	if err := db.Model(&row).Update("comment", "anything").Error; err != nil {
		t.Errorf("an untagged column was checked: %v", err)
	}
}
`
}

func apiPhoneGo() string {
	return `// Package phone validates, normalises and formats phone numbers with
// libphonenumber's metadata (github.com/nyaruka/phonenumbers, MIT).
//
// Importing it registers the tel format with internal/fieldtypes, which is
// why a model with a tel field imports it for effect:
//
//	import _ "{{MODULE}}/internal/phone"
package phone

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"

	"{{MODULE}}/internal/fieldtypes"
)

func init() {
	fieldtypes.Register("tel", func(value any, region string) (any, error) {
		switch v := value.(type) {
		case string:
			if v == "" {
				return v, nil
			}
			return Normalize(v, region)
		case *string:
			if v == nil || *v == "" {
				return v, nil
			}
			out, err := Normalize(*v, region)
			if err != nil {
				return nil, err
			}
			return &out, nil
		case nil:
			return nil, nil
		}
		return nil, errors.New("must be text")
	})
}

// Normalize returns raw in E.164 (+256772123456), or an error when it is not
// a number that can be dialled in its country. region is the ISO 3166-1 code
// a number written without its country code is read in; with none, the number
// has to start with + and its country code.
func Normalize(raw, region string) (string, error) {
	region = strings.ToUpper(strings.TrimSpace(region))
	if region == "" {
		region = "ZZ"
	}
	num, err := phonenumbers.Parse(strings.TrimSpace(raw), region)
	if err != nil {
		if region == "ZZ" {
			return "", errors.New("must be an international number starting with + and the country code, such as +256772123456")
		}
		return "", errors.New("is not a phone number")
	}
	if !phonenumbers.IsValidNumber(num) {
		country := phonenumbers.GetRegionCodeForNumber(num)
		if country == "" || country == "ZZ" {
			return "", errors.New("is not a valid phone number")
		}
		return "", fmt.Errorf("is not a valid phone number for %s", country)
	}
	return phonenumbers.Format(num, phonenumbers.E164), nil
}

// Country is the ISO 3166-1 code a stored E.164 number belongs to, or "".
func Country(e164 string) string {
	num, err := phonenumbers.Parse(e164, "ZZ")
	if err != nil {
		return ""
	}
	return phonenumbers.GetRegionCodeForNumber(num)
}

// International formats a stored number for people: +256 772 123456.
func International(e164 string) string {
	num, err := phonenumbers.Parse(e164, "ZZ")
	if err != nil {
		return e164
	}
	return phonenumbers.Format(num, phonenumbers.INTERNATIONAL)
}

// sampleRegions is the spread seeded numbers come from.
var sampleRegions = []string{
	"UG", "KE", "TZ", "RW", "NG", "GH", "ZA", "EG", "GB", "US",
	"CA", "DE", "FR", "IN", "BR", "AU", "JP", "AE", "MX", "NL",
}

// Sample returns a valid mobile number in E.164 for seeders, a different one
// for every n. It starts from libphonenumber's example mobile number for a
// country and replaces its last five digits with digits derived from n, then
// checks the result against the country's metadata, so the numbers are real
// shapes rather than random digits.
func Sample(n int) string {
	if n < 0 {
		n = -n
	}
	region := sampleRegions[n%len(sampleRegions)]
	k := n / len(sampleRegions)
	example := phonenumbers.GetExampleNumberForType(region, phonenumbers.MOBILE)
	if example == nil {
		return ""
	}
	national := phonenumbers.GetNationalSignificantNumber(example)
	if len(national) <= 5 {
		return phonenumbers.Format(example, phonenumbers.E164)
	}
	prefix := fmt.Sprintf("+%d%s", example.GetCountryCode(), national[:len(national)-5])
	for attempt := 0; attempt < 10; attempt++ {
		// 7919 is prime and coprime with 10^5, so k maps to distinct suffixes
		// for the first hundred thousand rows per country.
		suffix := ((k+attempt*10000)*7919 + 12345) % 100000
		candidate, err := phonenumbers.Parse(fmt.Sprintf("%s%05d", prefix, suffix), "ZZ")
		if err == nil && phonenumbers.IsValidNumber(candidate) {
			return phonenumbers.Format(candidate, phonenumbers.E164)
		}
	}
	return phonenumbers.Format(example, phonenumbers.E164)
}
`
}

func apiPhoneTestGo() string {
	return `package phone

import (
	"testing"

	"{{MODULE}}/internal/fieldtypes"
)

func TestNormalize(t *testing.T) {
	good := map[[2]string]string{
		{"+256 772 123456", ""}:   "+256772123456",
		{"0772 123456", "UG"}:     "+256772123456",
		{"+44 7400 123456", ""}:   "+447400123456",
		{"(201) 555-0123", "US"}:  "+12015550123",
	}
	for in, want := range good {
		got, err := Normalize(in[0], in[1])
		if err != nil || got != want {
			t.Errorf("Normalize(%q, %q) = %q, %v; want %q", in[0], in[1], got, err, want)
		}
	}
	for _, bad := range []string{"+256 12", "12345", "not a number", "+1 555 0000"} {
		if _, err := Normalize(bad, ""); err == nil {
			t.Errorf("Normalize accepted %q", bad)
		}
	}
}

func TestTheTelFormatIsRegistered(t *testing.T) {
	got, err := fieldtypes.Normalize("tel:UG", "0772 123456")
	if err != nil || got != "+256772123456" {
		t.Fatalf("tel:UG = %v, %v", got, err)
	}
}

func TestSamplesAreValidAndDistinct(t *testing.T) {
	seen := map[string]bool{}
	for n := 0; n < 5000; n++ {
		s := Sample(n)
		if _, err := Normalize(s, ""); err != nil {
			t.Fatalf("Sample(%d) = %q is not valid: %v", n, s, err)
		}
		if seen[s] {
			t.Fatalf("Sample(%d) repeats %q", n, s)
		}
		seen[s] = true
	}
}
`
}
