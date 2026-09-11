package generate

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// method returns one generated method, from its signature to its closing brace.
func method(t *testing.T, src, signature string) string {
	t.Helper()
	start := strings.Index(src, signature)
	if start < 0 {
		t.Fatalf("no %s in the generated file", signature)
	}
	end := strings.Index(src[start:], "\n}\n")
	if end < 0 {
		t.Fatalf("%s never closes", signature)
	}
	return src[start : start+end]
}

func mustParse(t *testing.T, name, src string) {
	t.Helper()
	if _, err := parser.ParseFile(token.NewFileSet(), name, src, parser.AllErrors); err != nil {
		t.Fatalf("%s does not parse: %v", name, err)
	}
}

// An owned resource checked ownership on list, read, update and delete, and on
// nothing else. Reproduced on a live project: a second ordinary account printed
// the first account's record as a PDF and rewrote it with PATCH.
func TestOwnedResourceGuardsTheSideDoors(t *testing.T) {
	_, _, h, s := ownedFiles(t)

	// PDF reads through the guarded GetByID; Patch loads through the guarded load.
	if !strings.Contains(method(t, h, "func (h *InvoiceHandler) PDF("), "h.service().GetByID(") {
		t.Error("PDF loads a row by id without checking who owns it")
	}
	if !strings.Contains(method(t, s, "func (s *InvoiceService) Patch("), "s.load(ctx, id)") {
		t.Error("Patch loads a row by id without checking who owns it")
	}
	if !strings.Contains(method(t, s, "func (s *InvoiceService) load("), "authz.Owns(ctx, &item)") {
		t.Error("the write path's load does not check who owns the row")
	}
	if !strings.Contains(method(t, s, "func (s *InvoiceService) Export("), `authz.ScopeOwned(ctx, query, "user_id")`) {
		t.Error("Export is not scoped, so it returns every row in the table")
	}
	if !strings.Contains(method(t, s, "func (s *InvoiceService) Bulk("), `authz.ScopeOwned(ctx, scope, "user_id")`) {
		t.Error("Bulk is not scoped, so a role allowed to reach it acts on anyone's rows")
	}
}

// Every CSV export was an empty 200. FindInBatches hands its callback a fresh
// session with no query on it, and the rows were re-read through tx.Scan.
func TestExportReadsTheBatchItWasGiven(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	def, err := ParseInlineFields("Product", "name:string,price:float")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.writeGoHandler(g.Names()); err != nil {
		t.Fatalf("handler: %v", err)
	}
	if err := g.writeGoService(g.Names()); err != nil {
		t.Fatalf("service: %v", err)
	}
	h := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "product.go"))
	s := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "services", "product.go"))
	mustParse(t, "handlers/product.go", h)
	mustParse(t, "services/product.go", s)
	export := method(t, s, "func (s *ProductService) Export(")

	if strings.Contains(export, "tx.Scan(") {
		t.Error("Export re-reads each batch through tx.Scan, which finds nothing")
	}
	if !strings.Contains(export, "return each(rows)") || strings.Count(export, "FindInBatches(&rows,") != 1 {
		t.Error("Export does not hand on the slice FindInBatches fills")
	}
	// A sort in front of the key FindInBatches pages by repeats rows.
	if strings.Contains(export, ".Order(") {
		t.Error("Export sorts ahead of the primary key it pages by")
	}
	handlerExport := method(t, h, "func (h *ProductHandler) Export(")
	if strings.Count(handlerExport, "h.service().Export(") != 2 {
		t.Error("CSV and XLSX should both read through the service's batches")
	}
	for _, src := range []string{export, handlerExport} {
		if strings.Contains(src, "_ = err") {
			t.Error("Export still discards its error")
		}
	}
}

// The importer took the owner from a CSV column, so an ordinary account could
// file records under anyone's name, and it created a user for every email it
// did not recognise.
func TestOwnedImportBelongsToTheImporter(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	g, names := ownedGenerator(t, root, module)
	if err := g.writeGoImportHandler(names); err != nil {
		t.Fatalf("import handler: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "invoice_import.go"))
	mustParse(t, "invoice_import.go", src)

	for _, want := range []string{
		"go h.runImportInvoice(job.ID, tmpPath, authz.CurrentUserID(c), authz.IsAdmin(c))",
		"item.UserID = ownerID",
		`ok && v != "" && canAssignOwner {`,
		`"shop/apps/api/internal/authz"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the owned importer is missing %s", want)
		}
	}
	if strings.Contains(src, "rel = models.User{") {
		t.Error("the importer still creates users")
	}
}

func TestImportNeverCreatesUsers(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	def, err := ParseInlineFields("Post", "title:string,author:belongs_to:User")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.writeGoImportHandler(g.Names()); err != nil {
		t.Fatalf("import handler: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "post_import.go"))
	mustParse(t, "post_import.go", src)
	if strings.Contains(src, "rel = models.User{") || strings.Contains(src, "canAssignOwner") {
		t.Error("a shared resource's importer creates users or grew owner logic")
	}
	if !strings.Contains(src, "no user with") {
		t.Error("an unknown user should fail the row with a message")
	}
}

// Transition loaded the row by id and moved it, whoever asked.
func TestOwnedWorkflowChecksTheOwner(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	path := filepath.Join(t.TempDir(), "ticket.yaml")
	body := `name: Ticket
owned_by: user
fields:
  - name: subject
    type: string
  - name: status
    type: select
    options:
      - value: open
      - value: closed
    workflow:
      initial: open
      terminal: [closed]
      transitions:
        - action: close
          from: [open]
          to: closed
  - name: user
    type: belongs_to
    related_model: User
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	def, err := LoadFromYAML(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.writeWorkflowService(g.Names(), def.WorkflowField()); err != nil {
		t.Fatalf("workflow service: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "services", "ticket_workflow.go"))
	mustParse(t, "ticket_workflow.go", src)
	if !strings.Contains(src, "item.GetOwnerID() != authz.CurrentUserID(c)") {
		t.Error("an owned resource's transitions do not check the owner")
	}
	if !strings.Contains(src, `"shop/apps/api/internal/authz"`) {
		t.Error("the owner check has no authz import")
	}
}

// The tree endpoints know nothing about owners, so an owned tree is refused
// rather than generated with a hole in it.
func TestOwnedTreeIsRefused(t *testing.T) {
	const module = "shop/apps/api"
	g, _ := ownedGenerator(t, setupMinimalProject(t, module), module)
	g.Definition.Tree = true
	err := g.checkFlagCombinations()
	if err == nil || !strings.Contains(err.Error(), "--tree and --owned-by") {
		t.Fatalf("an owned tree was accepted: %v", err)
	}
	g.Definition.Tree = false
	if err := g.checkFlagCombinations(); err != nil {
		t.Fatalf("an owned resource without --tree was refused: %v", err)
	}
}
