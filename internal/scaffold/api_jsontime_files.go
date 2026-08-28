package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeJSONTimeFiles writes internal/jsontime: the wire types for date and
// datetime fields.
func writeJSONTimeFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "jsontime", "jsontime.go"):      jsonTimeGo(),
		filepath.Join(apiRoot, "internal", "jsontime", "jsontime_test.go"): jsonTimeTestGo(),
	}
	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func jsonTimeGo() string {
	return `// Package jsontime carries the date and datetime types that cross the wire.
//
// The problem it solves: a date field used to be a *time.Time, and
// time.Time.UnmarshalJSON accepts RFC3339 and nothing else. The admin's date
// picker sends "2001-08-06", deliberately, because building a Date object from
// a string is how a date of birth ends up a day earlier for anyone west of
// Greenwich. So every date field failed to save with
//
//	parsing time "2001-08-06" as "2006-01-02T15:04:05Z07:00": cannot parse "" as "T"
//
// Rather than making the browser send a timezone it does not have, these types
// accept what a browser actually sends and store what was actually meant.
package jsontime

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Layouts accepted on the way in, widest first.
//
// The middle two are what <input type="datetime-local"> produces: no zone, and
// often no seconds. Neither is RFC3339, and both are what you get from a real
// browser, so both are accepted rather than rejected on principle.
var inputLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02T15:04",
	"2006-01-02",
}

func parse(raw string) (time.Time, error) {
	raw = strings.TrimSpace(strings.Trim(raw, ` + "`" + `"` + "`" + `))
	if raw == "" || raw == "null" {
		return time.Time{}, nil
	}
	for _, layout := range inputLayouts {
		// Parsed in UTC, not Local. A date has no zone, and resolving one
		// against the server's zone is what shifts a birthday by a day when
		// the server moves.
		if t, err := time.ParseInLocation(layout, raw, time.UTC); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot read %q as a date: expected YYYY-MM-DD or RFC3339", raw)
}

// Date is a calendar date: a year, a month and a day. No time, no zone.
//
// Marshals back as "2006-01-02", so what the API returns is exactly what the
// picker sends, and a value survives a round trip unchanged.
type Date struct{ time.Time }

func (d *Date) UnmarshalJSON(b []byte) error {
	t, err := parse(string(b))
	if err != nil {
		return err
	}
	// Truncated to the day on the way in. Anything finer is not part of what a
	// date means, and keeping it would let a stray time leak into a comparison.
	if !t.IsZero() {
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return []byte(` + "`" + `"` + "`" + ` + d.Format("2006-01-02") + ` + "`" + `"` + "`" + `), nil
}

// Value and Scan let GORM store this like any other timestamp column.
func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

func (d *Date) Scan(v any) error {
	if v == nil {
		d.Time = time.Time{}
		return nil
	}
	switch t := v.(type) {
	case time.Time:
		// Read back as the calendar day that was stored, with the zone the
		// driver attached discarded rather than applied.
		d.Time = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		return nil
	case string:
		parsed, err := parse(t)
		if err != nil {
			return err
		}
		d.Time = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
		return nil
	case []byte:
		return d.Scan(string(t))
	}
	return fmt.Errorf("cannot scan %T into a Date", v)
}

// DateTime is a timestamp that accepts the shapes a browser actually sends.
//
// Marshals as RFC3339, because unlike a date this genuinely is an instant and
// the zone is part of the value.
type DateTime struct{ time.Time }

func (d *DateTime) UnmarshalJSON(b []byte) error {
	t, err := parse(string(b))
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return []byte(` + "`" + `"` + "`" + ` + d.Format(time.RFC3339) + ` + "`" + `"` + "`" + `), nil
}

func (d DateTime) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

func (d *DateTime) Scan(v any) error {
	if v == nil {
		d.Time = time.Time{}
		return nil
	}
	switch t := v.(type) {
	case time.Time:
		d.Time = t
		return nil
	case string:
		parsed, err := parse(t)
		if err != nil {
			return err
		}
		d.Time = parsed
		return nil
	case []byte:
		return d.Scan(string(t))
	}
	return fmt.Errorf("cannot scan %T into a DateTime", v)
}
`
}

func jsonTimeTestGo() string {
	return `package jsontime

import (
	"encoding/json"
	"testing"
	"time"
)

// The bug these types exist for: the admin date picker sends "2001-08-06",
// and a *time.Time refuses it.
func TestDateAcceptsWhatTheBrowserSends(t *testing.T) {
	for _, in := range []string{
		` + "`" + `"2001-08-06"` + "`" + `,
		` + "`" + `"2001-08-06T00:00:00Z"` + "`" + `,
		` + "`" + `"2001-08-06T14:30"` + "`" + `,
		` + "`" + `"2001-08-06T14:30:00"` + "`" + `,
	} {
		var d Date
		if err := json.Unmarshal([]byte(in), &d); err != nil {
			t.Fatalf("%s was rejected: %v", in, err)
		}
		if d.Year() != 2001 || d.Month() != time.August || d.Day() != 6 {
			t.Errorf("%s parsed as %s", in, d.Format(time.RFC3339))
		}
	}
}

// A date has no time. Anything finer is dropped so it cannot leak into a
// comparison later.
func TestDateTruncatesToTheDay(t *testing.T) {
	var d Date
	if err := json.Unmarshal([]byte(` + "`" + `"2001-08-06T14:30:00"` + "`" + `), &d); err != nil {
		t.Fatal(err)
	}
	if d.Hour() != 0 || d.Minute() != 0 {
		t.Errorf("expected midnight, got %s", d.Format(time.RFC3339))
	}
}

// What goes out is what the picker sends back in, unchanged.
func TestDateRoundTrips(t *testing.T) {
	var d Date
	if err := json.Unmarshal([]byte(` + "`" + `"2001-08-06"` + "`" + `), &d); err != nil {
		t.Fatal(err)
	}
	out, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != ` + "`" + `"2001-08-06"` + "`" + ` {
		t.Errorf("round trip changed the value: %s", out)
	}
}

// The whole point of parsing in UTC: a birthday must not move because the
// server did.
func TestDateDoesNotShiftWithTheServerZone(t *testing.T) {
	west, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("tzdata not available on this platform")
	}
	original := time.Local
	time.Local = west
	defer func() { time.Local = original }()

	var d Date
	if err := json.Unmarshal([]byte(` + "`" + `"2001-08-06"` + "`" + `), &d); err != nil {
		t.Fatal(err)
	}
	if got := d.Format("2006-01-02"); got != "2001-08-06" {
		t.Errorf("the date moved to %s under a western server zone", got)
	}
}

func TestEmptyAndNullAreZero(t *testing.T) {
	for _, in := range []string{` + "`" + `""` + "`" + `, "null"} {
		var d Date
		if err := json.Unmarshal([]byte(in), &d); err != nil {
			t.Fatalf("%s should be accepted as empty: %v", in, err)
		}
		if !d.IsZero() {
			t.Errorf("%s should be the zero date", in)
		}
		out, _ := json.Marshal(d)
		if string(out) != "null" {
			t.Errorf("an empty date should marshal as null, got %s", out)
		}
	}
}

func TestGarbageIsRejected(t *testing.T) {
	var d Date
	if err := json.Unmarshal([]byte(` + "`" + `"the sixth of August"` + "`" + `), &d); err == nil {
		t.Error("expected an error, and one that says what the format should be")
	}
}

// datetime-local sends no zone and often no seconds. Neither is RFC3339.
func TestDateTimeAcceptsDatetimeLocal(t *testing.T) {
	var d DateTime
	if err := json.Unmarshal([]byte(` + "`" + `"2001-08-06T14:30"` + "`" + `), &d); err != nil {
		t.Fatalf("datetime-local was rejected: %v", err)
	}
	if d.Hour() != 14 || d.Minute() != 30 {
		t.Errorf("time lost: %s", d.Format(time.RFC3339))
	}
}

// A driver hands back a time in whatever zone it likes. The stored calendar
// day is what was meant.
func TestDateScanKeepsTheStoredDay(t *testing.T) {
	east := time.FixedZone("east", 3*60*60)
	var d Date
	if err := d.Scan(time.Date(2001, 8, 6, 0, 0, 0, 0, east)); err != nil {
		t.Fatal(err)
	}
	if got := d.Format("2006-01-02"); got != "2001-08-06" {
		t.Errorf("scanned back as %s", got)
	}
}

func TestValueIsNilWhenEmpty(t *testing.T) {
	var d Date
	v, err := d.Value()
	if err != nil {
		t.Fatal(err)
	}
	if v != nil {
		t.Errorf("an empty date should store as NULL, got %v", v)
	}
}
`
}
