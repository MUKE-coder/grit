package generate

import (
	"strings"
	"testing"
	"time"
)

func TestParseRetentionAcceptsDaysWeeksAndHours(t *testing.T) {
	day := 24 * time.Hour
	for in, want := range map[string]time.Duration{
		"90d":    90 * day,
		"2w":     14 * day,
		"36h":    36 * time.Hour,
		"1d":     day,
		" 30D ":  30 * day, // typed in a shell or a YAML file, case and spaces vary
		"36500d": 36500 * day,
	} {
		got, err := ParseRetention(in)
		if err != nil {
			t.Errorf("ParseRetention(%q): %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("ParseRetention(%q) = %v, want %v", in, got, want)
		}
	}
}

// Each of these would otherwise become a window nobody meant: zero or negative
// deletes every row on the first sweep, and a month guessed at 30 days deletes
// a day early in every 31-day month.
func TestParseRetentionRefusesWhatItCannotMeanExactly(t *testing.T) {
	for in, mention := range map[string]string{
		"":                   "needs a window",
		"0d":                 "whole number",
		"-5d":                "whole number",
		"d":                  "whole number",
		"90":                 "unknown unit",
		"1.5d":               "whole number",
		"1d12h":              "whole number",
		"6m":                 "days",
		"1y":                 "days",
		"90s":                "unknown unit",
		"36501d":             "100 years",
		"9999999999999999w":  "100 years", // would overflow time.Duration if multiplied first
		"99999999999999999h": "100 years",
	} {
		_, err := ParseRetention(in)
		if err == nil {
			t.Errorf("ParseRetention(%q) was accepted", in)
			continue
		}
		if !strings.Contains(err.Error(), mention) {
			t.Errorf("ParseRetention(%q): the error does not say what to do (want %q): %v", in, mention, err)
		}
	}
}

func TestRetentionWindowIsZeroWhenUnset(t *testing.T) {
	d := &ResourceDefinition{Name: "Message"}
	if w, err := d.RetentionWindow(); err != nil || w != 0 {
		t.Errorf("a resource without --retain got window %v, err %v", w, err)
	}
}

// Refused before anything is written, the same as every other combination the
// generated code cannot honour.
func TestRetainRefusesWhatTheSweepCouldNotKeep(t *testing.T) {
	bad := &Generator{Definition: &ResourceDefinition{Name: "Message", Retain: "soon"}, Module: "app"}
	if err := bad.checkFlagCombinations(); err == nil {
		t.Error("an unreadable --retain value was accepted")
	}

	// Append-only rows cannot be deleted at all, so the sweep would fail every
	// night while the policy said the rows were going.
	ledger := &Generator{Definition: &ResourceDefinition{Name: "LedgerEntry", Retain: "90d", AppendOnly: true}, Module: "app"}
	err := ledger.checkFlagCombinations()
	if err == nil {
		t.Error("--retain and --append-only were accepted together")
	} else if !strings.Contains(err.Error(), "--append-only") {
		t.Errorf("the error does not name the conflicting flag: %v", err)
	}

	ok := &Generator{Definition: &ResourceDefinition{Name: "Message", Retain: "90d"}, Module: "app"}
	if err := ok.checkFlagCombinations(); err != nil {
		t.Errorf("a plain --retain 90d was refused: %v", err)
	}
}
