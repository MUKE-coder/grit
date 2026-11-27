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

// A record's lines go with it, and a tree goes leaf first, so each is accepted
// and says so in the model, alone or together.
func TestRetainShapesThePolicyForRowsOtherRowsPointAt(t *testing.T) {
	items := &Generator{Definition: &ResourceDefinition{Name: "PurchaseOrder", Retain: "90d",
		Items: &ResourceDefinition{Name: "PurchaseOrderLine"}}, Module: "app"}
	if err := items.checkFlagCombinations(); err != nil {
		t.Errorf("--retain with --items was refused: %v", err)
	}
	// Two words, so the foreign key is the snake name GORM gives the
	// has-many's PurchaseOrderID, not a lowercased "purchaseorder_id".
	want := `retention.Register(&PurchaseOrder{}, 90*24*time.Hour, retention.WithLines(&PurchaseOrderLine{}, "purchase_order_id")) }`
	got := items.retentionModelInit(MakeNames("PurchaseOrder"))
	if !strings.Contains(got, want) {
		t.Errorf("the lines are not registered with their record:\n%s", got)
	}
	if !strings.Contains(got, "// Its PurchaseOrderLine rows go with it, whatever their own age.") {
		t.Errorf("the model does not say its lines go with it:\n%s", got)
	}

	tree := &Generator{Definition: &ResourceDefinition{Name: "Category", Retain: "90d", Tree: true}, Module: "app"}
	if err := tree.checkFlagCombinations(); err != nil {
		t.Errorf("--retain with --tree was refused: %v", err)
	}
	if got := tree.retentionModelInit(MakeNames("Category")); !strings.Contains(got, "retention.AsTree()) }") ||
		!strings.Contains(got, "// A row stays while any row below it is still inside the window.") {
		t.Errorf("the tree is not registered as one, or the model does not say what that means:\n%s", got)
	}

	// Both: the lines of a node go when the node goes, as on #310, so a node
	// kept for its younger children keeps its lines too.
	both := &Generator{Definition: &ResourceDefinition{Name: "Project", Retain: "90d", Tree: true,
		Items: &ResourceDefinition{Name: "ProjectTask"}}, Module: "app"}
	if err := both.checkFlagCombinations(); err != nil {
		t.Errorf("--retain with --tree and --items together was refused: %v", err)
	}
	if got := both.retentionModelInit(MakeNames("Project")); !strings.Contains(got,
		`retention.WithLines(&ProjectTask{}, "project_id"), retention.AsTree()) }`) {
		t.Errorf("a tree with lines is not registered as both:\n%s", got)
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

// The window is shown the way it was asked for, tidied, so neither the model
// comment nor the admin repeats a typo back.
func TestRetentionSpecIsTheWindowTidied(t *testing.T) {
	for in, want := range map[string]string{
		"90d": "90d", " 90D ": "90d", "090d": "90d", "2w": "2w", "36h": "36h",
		"": "", "soon": "", "6m": "",
	} {
		d := &ResourceDefinition{Name: "Message", Retain: in}
		if got := d.RetentionSpec(); got != want {
			t.Errorf("RetentionSpec(%q) = %q, want %q", in, got, want)
		}
	}
}

// --retain reaches the admin as the window its list page reads, and a resource
// without it says nothing there.
func TestRetainIsShownInTheAdminDefinition(t *testing.T) {
	const module = "chat/apps/api"
	root := setupMinimalProject(t, module)

	def := mustFields(t, "ChatMessage", "body:text")
	def.Retain = " 90D"
	g := newTestGenerator(root, module, def)
	if src := g.adminDefinitionContent(g.Names()); !strings.Contains(src, "\n  retention: \"90d\",") {
		t.Errorf("the admin definition does not carry the window:\n%s", src)
	}

	plain := newTestGenerator(root, module, mustFields(t, "ChatMessage", "body:text"))
	if src := plain.adminDefinitionContent(plain.Names()); strings.Contains(src, "retention:") {
		t.Error("a resource without --retain claims a window in the admin")
	}
}

// A project that took --retain before lines and trees could be retained has
// the package without the options. Generating into it would write a model that
// does not compile, so it is refused before anything is written.
func TestRetainRefusesAnOlderRegistryForLinesAndTrees(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	older := "package retention\n\nimport \"time\"\n\n" +
		"func Register(model interface{}, window time.Duration) {}\n"
	path := filepath.Join(root, "apps", "api", "internal", "retention", "retention.go")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(older), 0o644); err != nil {
		t.Fatal(err)
	}

	def := mustFields(t, "Invoice", "number:string")
	def.Items = mustFields(t, "InvoiceLine", "description:string")
	def.Retain = "90d"
	err := newTestGenerator(root, module, def).Run()
	if err == nil || !strings.Contains(err.Error(), "grit upgrade") {
		t.Fatalf("an older registry was not refused with the way out: %v", err)
	}
	for _, model := range []string{"invoice.go", "invoice_line.go"} {
		if _, err := os.Stat(filepath.Join(root, "apps", "api", "internal", "models", model)); err == nil {
			t.Errorf("%s was written although the generate was refused", model)
		}
	}
}
