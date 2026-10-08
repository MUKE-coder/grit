package jsontime

import (
	"encoding/json"
	"testing"
	"time"
)

// The bug these types exist for: the admin date picker sends "2001-08-06",
// and a *time.Time refuses it.
func TestDateAcceptsWhatTheBrowserSends(t *testing.T) {
	for _, in := range []string{
		`"2001-08-06"`,
		`"2001-08-06T00:00:00Z"`,
		`"2001-08-06T14:30"`,
		`"2001-08-06T14:30:00"`,
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
	if err := json.Unmarshal([]byte(`"2001-08-06T14:30:00"`), &d); err != nil {
		t.Fatal(err)
	}
	if d.Hour() != 0 || d.Minute() != 0 {
		t.Errorf("expected midnight, got %s", d.Format(time.RFC3339))
	}
}

// What goes out is what the picker sends back in, unchanged.
func TestDateRoundTrips(t *testing.T) {
	var d Date
	if err := json.Unmarshal([]byte(`"2001-08-06"`), &d); err != nil {
		t.Fatal(err)
	}
	out, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != `"2001-08-06"` {
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
	if err := json.Unmarshal([]byte(`"2001-08-06"`), &d); err != nil {
		t.Fatal(err)
	}
	if got := d.Format("2006-01-02"); got != "2001-08-06" {
		t.Errorf("the date moved to %s under a western server zone", got)
	}
}

func TestEmptyAndNullAreZero(t *testing.T) {
	for _, in := range []string{`""`, "null"} {
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
	if err := json.Unmarshal([]byte(`"the sixth of August"`), &d); err == nil {
		t.Error("expected an error, and one that says what the format should be")
	}
}

// datetime-local sends no zone and often no seconds. Neither is RFC3339.
func TestDateTimeAcceptsDatetimeLocal(t *testing.T) {
	var d DateTime
	if err := json.Unmarshal([]byte(`"2001-08-06T14:30"`), &d); err != nil {
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
