package doctor

import (
	"strings"
	"testing"
)

const invoiceModel = "apps/api/internal/models/invoice.go"

// cleanProject's Invoice without its retention window.
func invoiceWithoutRetention(t *testing.T) string {
	t.Helper()
	src := cleanProject()[invoiceModel]
	at := strings.Index(src, "\n\n// An owned")
	if at < 0 {
		at = strings.Index(src, "\n\nfunc init() { retention.Register(")
	}
	if at < 0 {
		t.Fatal("the clean Invoice has no retention window to take out, so this test proves nothing")
	}
	return src[:at] + "\n"
}

// The case the warning is for: a user's rows with nothing that removes them.
func TestAnOwnedResourceWithNoRetentionIsAsked(t *testing.T) {
	report := run(t, map[string]string{invoiceModel: invoiceWithoutRetention(t)})
	found := fired(report, "owned-resource-no-retention")
	if len(found) != 1 {
		t.Fatalf("got %d findings, want 1: %v", len(found), found)
	}
	f := found[0]
	if f.Level != "warning" || f.Resource != "Invoice" {
		t.Errorf("got a %s on %s, want a warning on Invoice", f.Level, f.Resource)
	}
	// It says what the model shows. Whether the data is personal is the
	// project's call, and the maintainer asked for exactly this wording on
	// grit#260.
	if !strings.Contains(f.Message, "owned by a user") || strings.Contains(strings.ToLower(f.Message), "personal") {
		t.Errorf("the message should give its own reasoning, not a verdict on the data: %q", f.Message)
	}
	// Both answers are in the fix, because "keep them" is often right.
	if !strings.Contains(f.Fix, "retention.Register(&Invoice{}") || !strings.Contains(f.Fix, "nothing to do") {
		t.Errorf("the fix should name the edit and say when none is needed: %q", f.Fix)
	}
}

// --confidential is the project saying the rows are sensitive, and the message
// says that it was.
func TestAConfidentialResourceWithNoRetentionSaysSo(t *testing.T) {
	service := strings.ReplaceAll(cleanProject()["apps/api/internal/services/invoice.go"],
		`authz.ScopeOwned(ctx, query, "user_id")`,
		`authz.ScopeConfidential(ctx, query, "user_id", "invoices.confidential")`)
	report := run(t, map[string]string{
		invoiceModel:                            invoiceWithoutRetention(t),
		"apps/api/internal/services/invoice.go": service,
	})
	found := fired(report, "owned-resource-no-retention")
	if len(found) != 1 || !strings.Contains(found[0].Message, "confidential") {
		t.Fatalf("want one finding that mentions confidential, got %v", found)
	}
}

// --retain refuses an append-only resource, so suggesting it is advice the
// generator turns down.
func TestAnAppendOnlyResourceIsNotToldToRetain(t *testing.T) {
	model := invoiceWithoutRetention(t) + "\nfunc init() { appendonly.Register(&Invoice{}) }\n"
	if found := fired(run(t, map[string]string{invoiceModel: model}), "owned-resource-no-retention"); len(found) > 0 {
		t.Errorf("an append-only resource was told to take a window it cannot have: %q", found[0].Message)
	}
}

// Nobody owns the rows, so this is not the question: resource-could-be-owned
// asks first whether somebody should.
func TestAnUnownedResourceIsNotAsked(t *testing.T) {
	model := strings.Replace(invoiceWithoutRetention(t), "func (m *Invoice) GetOwnerID() string { return m.UserID }\n", "", 1)
	if found := fired(run(t, map[string]string{invoiceModel: model}), "owned-resource-no-retention"); len(found) > 0 {
		t.Errorf("a resource nobody owns was asked about retention: %q", found[0].Message)
	}
}

// The window is read from the line --retain writes, in both shapes its
// duration takes: whole days, and hours.
func TestTheRetentionLineTheGeneratorWritesIsRead(t *testing.T) {
	for _, line := range []string{
		"func init() { retention.Register(&Invoice{}, 90*24*time.Hour) }",
		"func init() { retention.Register(&Invoice{}, 36*time.Hour) }",
	} {
		models := parseModels(invoiceWithoutRetention(t) + "\n" + line + "\n")
		if len(models) != 1 || !models[0].Retained {
			t.Errorf("not read as retained: %s", line)
		}
	}
}
