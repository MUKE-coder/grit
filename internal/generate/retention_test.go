package generate

import (
	"go/format"
	"os"
	"path/filepath"
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

func TestRetentionGoExprReadsAsThePolicy(t *testing.T) {
	for window, want := range map[time.Duration]string{
		90 * 24 * time.Hour: "90*24*time.Hour",
		14 * 24 * time.Hour: "14*24*time.Hour", // --retain 2w
		36 * time.Hour:      "36*time.Hour",
	} {
		if got := retentionGoExpr(window); got != want {
			t.Errorf("retentionGoExpr(%v) = %q, want %q", window, got, want)
		}
	}
}

// Rows other rows point at by foreign key would be refused by the database
// every night, so the combinations are refused up front instead.
func TestRetainRefusesRowsOtherRowsPointAt(t *testing.T) {
	tree := &Generator{Definition: &ResourceDefinition{Name: "Category", Retain: "90d", Tree: true}, Module: "app"}
	if err := tree.checkFlagCombinations(); err == nil || !strings.Contains(err.Error(), "--tree") {
		t.Errorf("--retain with --tree was not refused by name: %v", err)
	}
	items := &Generator{Definition: &ResourceDefinition{Name: "Invoice", Retain: "90d",
		Items: &ResourceDefinition{Name: "InvoiceItem"}}, Module: "app"}
	if err := items.checkFlagCombinations(); err == nil || !strings.Contains(err.Error(), "--items") {
		t.Errorf("--retain with --items was not refused by name: %v", err)
	}
}

// The policy lives in the model file, so the model has it, the package it
// imports exists, and removing the resource takes the policy with it. Two
// words, because a one-word fixture cannot tell name variants apart.
func TestRetainIsWrittenToTheModelAndRemovedWithIt(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)
	def, err := ParseInlineFields("ChatMessage", "body:text,attachment:file")
	if err != nil {
		t.Fatal(err)
	}
	def.Retain = "90d"

	if err := newTestGenerator(root, module, def).Run(); err != nil {
		t.Fatalf("Run(): %v", err)
	}

	api := filepath.Join(root, "apps", "api", "internal")
	model, err := os.ReadFile(filepath.Join(api, "models", "chat_message.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"myapp/apps/api/internal/retention"`,
		"func init() { retention.Register(&ChatMessage{}, 90*24*time.Hour) }",
	} {
		if !strings.Contains(string(model), want) {
			t.Errorf("the model is missing %s:\n%s", want, model)
		}
	}
	if _, err := format.Source(model); err != nil {
		t.Errorf("the model is not valid Go: %v", err)
	}
	// An older project has no package for the model to import.
	if _, err := os.Stat(filepath.Join(api, "retention", "retention.go")); err != nil {
		t.Errorf("the retention package was not written for the model to import: %v", err)
	}

	names := MakeNames("ChatMessage")
	t.Chdir(root)
	if err := RemoveResource("ChatMessage"); err != nil {
		t.Fatalf("RemoveResource: %v", err)
	}
	for _, line := range survivors(t, root, names) {
		t.Errorf("remove left this behind: %s", line)
	}
}

// Without --retain the model is exactly what it was before the flag existed.
func TestNoRetainWritesNoPolicy(t *testing.T) {
	g := &Generator{Definition: &ResourceDefinition{Name: "Message"}, Module: "app"}
	if got := g.retentionModelInit(MakeNames("Message")); got != "" {
		t.Errorf("a resource without --retain got a policy:\n%s", got)
	}
}
