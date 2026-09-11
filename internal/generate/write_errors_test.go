package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// Generated handlers report the error a write came back with.
//
// Five places bound err and never used it, returning a constant 500:
//
//	if err := h.scoped(c).Create(&item).Error; err != nil {
//	    c.JSON(http.StatusInternalServerError, gin.H{... "Failed to create x"})
//
// So a rule enforced in a GORM hook, which is where a rule has to live if
// Studio and the CSV importer are to respect it too, reached neither the
// client nor the log. Found building a double-entry ledger, where "USD does
// not balance: debits 100.00 USD, credits 99.99 USD" arrived as "Failed to
// create journalentry" with status 500.
func TestGeneratedHandlerDoesNotSwallowWriteErrors(t *testing.T) {
	const module = "ledger/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Invoice", "reference:string,total:money")
	if err != nil {
		t.Fatalf("fields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.writeGoHandler(g.Names()); err != nil {
		t.Fatalf("handler: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "invoice.go"))

	// Every write path hands the service's error to fail.
	for _, want := range []string{
		`h.fail(c, err, "Failed to create invoice")`,
		`h.fail(c, err, "Failed to update invoice")`,
		`h.fail(c, err, "Failed to patch invoice")`,
		`h.fail(c, err, "Failed to delete invoice")`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("missing: %s", want)
		}
	}
	// And fail passes on whatever it does not answer itself: a broken rule as
	// 422 with its message, anything else as a logged 500.
	if !strings.Contains(method(t, src, "func (h *InvoiceHandler) fail("), "respond.WriteError(c, err, fallback)") {
		t.Error("fail does not hand the error to respond.WriteError")
	}
	if !strings.Contains(src, module+"/internal/respond") {
		t.Error("the handler calls respond and does not import it")
	}
	if strings.Contains(src, `"message": "Failed to create invoice",`) {
		t.Error("create still returns a constant 500 and discards err")
	}
}
