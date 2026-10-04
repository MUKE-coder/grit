// Package jsontime carries the date and datetime types that cross the wire.
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
	raw = strings.TrimSpace(strings.Trim(raw, `"`))
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
	return []byte(`"` + d.Format("2006-01-02") + `"`), nil
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
//
// Whole seconds, on purpose: it is for datetimes people enter, where more
// digits are noise. For a time the system records and compares, such as when
// a message was read, use time.Time, which keeps the fraction: here, a message
// and a read receipt a few milliseconds apart marshal as the same second.
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
	return []byte(`"` + d.Format(time.RFC3339) + `"`), nil
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
