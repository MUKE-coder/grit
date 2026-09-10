package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// A --items child whose fields need imports gets them.
//
// The inline Items request struct is built from the child's fields and written
// into the parent's handler, so a money field on the child puts money.Money in
// that file. The import scan walked only the parent's fields, so nothing added
// internal/money and the generated project did not build:
//
//	internal/handlers/journal_entry.go:281:10: undefined: money
//
// Found generating a double-entry journal, where every line is money by
// definition.
func TestItemsChildFieldsBringTheirImports(t *testing.T) {
	const module = "ledger/apps/api"
	root := setupMinimalProject(t, module)

	child, err := ParseInlineFields("JournalLine", "account:belongs_to:Account,debit:money,credit:money")
	if err != nil {
		t.Fatalf("child fields: %v", err)
	}
	parent, err := ParseInlineFields("JournalEntry", "reference:string,memo:text")
	if err != nil {
		t.Fatalf("parent fields: %v", err)
	}
	parent.Items = child

	g := newTestGenerator(root, module, parent)
	if err := g.writeGoHandler(g.Names()); err != nil {
		t.Fatalf("handler: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "journal_entry.go"))

	if strings.Contains(src, "money.Money") && !strings.Contains(src, module+"/internal/money") {
		t.Error("the handler uses money.Money and does not import internal/money")
	}
}

// A belongs_to on a --items child reaches the request and the insert.
//
// The builder skipped every belongs_to, on the grounds that the FK back to the
// parent is set by the association and clients never send it per row. True of
// the parent's FK. Not true of a belongs_to pointing anywhere else, and those
// are the ones the client has to send: the account a journal line posts to,
// the product an invoice line bills.
//
// So the request struct had no way to say which account, the create loop set
// nothing, and every line was inserted with an empty account_id while the
// model declared the column required. No error, just wrong rows.
func TestItemsChildKeepsItsOwnBelongsTo(t *testing.T) {
	const module = "ledger/apps/api"
	root := setupMinimalProject(t, module)

	child, err := ParseInlineFields("JournalLine", "account:belongs_to:Account,debit:money,credit:money")
	if err != nil {
		t.Fatalf("child fields: %v", err)
	}
	parent, err := ParseInlineFields("JournalEntry", "reference:string")
	if err != nil {
		t.Fatalf("parent fields: %v", err)
	}
	parent.Items = child

	g := newTestGenerator(root, module, parent)
	if err := g.writeGoHandler(g.Names()); err != nil {
		t.Fatalf("handler: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "journal_entry.go"))

	// The FK, named the way the parent's own belongs_to fields are named. The
	// association itself is a *Account and would not compile here.
	if !strings.Contains(src, `AccountID string`) {
		t.Error("the items request struct cannot say which account the line posts to")
	}
	if !strings.Contains(src, "AccountID: it.AccountID") {
		t.Error("the create loop never sets AccountID, so lines insert with it empty")
	}
	if strings.Contains(src, "Account *models.Account") || strings.Contains(src, "Account Account") {
		t.Error("the association leaked into the request struct instead of the FK")
	}

	// And the parent's own FK stays out: the association sets it on create, so
	// asking a client to send it per row would be wrong.
	if strings.Contains(src, "JournalEntryID string `json:\"journal_entry_id\"") {
		t.Error("the parent FK leaked into the per-row request struct")
	}
}
