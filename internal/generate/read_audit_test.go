package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The activity log recorded every write and no read, so "who opened this
// patient's chart" had no answer. A resource generated with --audit-reads
// marks every read path.
func TestAuditReadsMarksEveryRead(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	g, names := ownedGenerator(t, root, module)
	g.Definition.AuditReads = true
	if err := g.writeGoHandler(names); err != nil {
		t.Fatalf("handler: %v", err)
	}
	h := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "invoice.go"))
	mustParse(t, "invoice.go", h)

	checks := map[string]string{
		"func (h *InvoiceHandler) List(":    `audit.Read(c, "invoices", ids...)`,
		"func (h *InvoiceHandler) GetByID(": `audit.Read(c, "invoices", item.ID)`,
		"func (h *InvoiceHandler) PDF(":     `audit.Read(c, "invoices", item.ID)`,
		"func (h *InvoiceHandler) Export(":  `audit.ReadCount(c, "invoices", exported)`,
	}
	for sig, want := range checks {
		if !strings.Contains(method(t, h, sig), want) {
			t.Errorf("%s does not record its read (%s)", sig, want)
		}
	}
	if !strings.Contains(method(t, h, "func (h *InvoiceHandler) Export("), `audit.ReadCount(c, "invoices", len(all))`) {
		t.Error("the XLSX export does not record its read")
	}
	if !strings.Contains(h, `"shop/apps/api/internal/audit"`) {
		t.Error("the audit import is missing")
	}
}

func TestResourcesWithoutAuditReadsStayQuiet(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	def, err := ParseInlineFields("Product", "name:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.writeGoHandler(g.Names()); err != nil {
		t.Fatalf("handler: %v", err)
	}
	h := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "product.go"))
	mustParse(t, "product.go", h)
	if strings.Contains(h, "audit.") || strings.Contains(h, "exported") {
		t.Error("a resource generated without --audit-reads records its reads")
	}
}

// Marks the middleware drops would look finished and record nothing.
func TestAuditReadsNeedsALogThatRecordsReads(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	g, _ := ownedGenerator(t, root, module)
	g.Definition.AuditReads = true
	if err := g.prepareReadAudit(); err == nil || !strings.Contains(err.Error(), "grit upgrade") {
		t.Fatalf("a project whose log cannot record reads was accepted: %v", err)
	}

	api := g.APIRoot()
	for path, body := range map[string]string{
		filepath.Join(api, "internal", "audit", "audit.go"):           "func Read(c *gin.Context, resource string, ids ...string) {}",
		filepath.Join(api, "internal", "middleware", "activity.go"): "mark, ok := audit.ReadMarkOf(c)",
	} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := g.prepareReadAudit(); err != nil {
		t.Fatalf("an upgraded project was refused: %v", err)
	}
}

func TestAuditReadsRefusesPublicAndTree(t *testing.T) {
	const module = "shop/apps/api"
	g, _ := ownedGenerator(t, setupMinimalProject(t, module), module)
	g.Definition.OwnedBy = ""
	g.Definition.AuditReads = true

	g.Definition.Public = true
	if err := g.checkFlagCombinations(); err == nil || !strings.Contains(err.Error(), "--public") {
		t.Errorf("--audit-reads with --public was accepted: %v", err)
	}
	g.Definition.Public = false
	g.Definition.Tree = true
	if err := g.checkFlagCombinations(); err == nil || !strings.Contains(err.Error(), "--tree") {
		t.Errorf("--audit-reads with --tree was accepted: %v", err)
	}
	g.Definition.Tree = false
	if err := g.checkFlagCombinations(); err != nil {
		t.Errorf("a plain audited resource was refused: %v", err)
	}
}
