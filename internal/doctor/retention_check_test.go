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

// An invoice's lines go with it (--items), so they are covered by its window
// even though their own model file registers nothing. The record names them,
// in a file the doctor reads before or after theirs.
func TestLinesAreCoveredByTheirRecordsWindow(t *testing.T) {
	line := "package models\n\n" +
		"type InvoiceLine struct {\n" +
		"\tID        string `gorm:\"primaryKey\" json:\"id\"`\n" +
		"\tInvoiceID string `gorm:\"size:36;index\" json:\"invoice_id\"`\n" +
		"\tUserID    string `gorm:\"size:36;index\" json:\"user_id\"`\n" +
		"}\n\n" +
		"func (m *InvoiceLine) GetOwnerID() string { return m.UserID }\n"
	record := invoiceWithoutRetention(t) +
		"\nfunc init() { retention.Register(&Invoice{}, 90*24*time.Hour, retention.WithLines(&InvoiceLine{}, \"invoice_id\")) }\n"
	handler := "package handlers\n\ntype BulkInvoiceLineRequest struct{}\n"

	report := run(t, map[string]string{
		invoiceModel: record,
		"apps/api/internal/models/invoice_line.go":   line,
		"apps/api/internal/handlers/invoice_line.go": handler,
	})
	for _, f := range fired(report, "owned-resource-no-retention") {
		t.Errorf("%s was asked about retention: %s", f.Resource, f.Message)
	}
	if len(report.Resources) != 2 {
		t.Fatalf("the lines were not read as a resource, so this test proves nothing: %v", report.Resources)
	}
}

// The suggested line carries the options the resource's shape needs, the same
// ones --retain writes. Without them it compiles and the database refuses the
// sweep every night.
func TestTheSuggestedLineFitsTheResourcesShape(t *testing.T) {
	plain := invoiceWithoutRetention(t)
	withLines := strings.Replace(plain, "\tUser   User",
		"\tItems  []InvoiceLine `gorm:\"foreignKey:InvoiceID\" json:\"items\"`\n\tUser   User", 1)
	asTree := strings.Replace(plain, "\tUser   User",
		"\t// --- tree (grit --tree) ---------------------------------------------------\n\tUser   User", 1)
	if withLines == plain || asTree == plain {
		t.Fatal("the fixture's Invoice changed shape, so this test proves nothing")
	}

	for name, c := range map[string]struct{ model, want string }{
		"lines": {withLines, `retention.Register(&Invoice{}, 90*24*time.Hour, retention.WithLines(&InvoiceLine{}, "invoice_id")) }`},
		"tree":  {asTree, "retention.Register(&Invoice{}, 90*24*time.Hour, retention.AsTree()) }"},
		"plain": {plain, "retention.Register(&Invoice{}, 90*24*time.Hour) }"},
	} {
		found := fired(run(t, map[string]string{invoiceModel: c.model}), "owned-resource-no-retention")
		if len(found) != 1 || !strings.Contains(found[0].Fix, c.want) {
			t.Errorf("%s: the fix does not suggest %s: %v", name, c.want, found)
		}
	}
}
