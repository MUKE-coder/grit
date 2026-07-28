package generate

import (
	"strings"
	"testing"
)

// TestColumnClick_FirstColumnLinks — the first plain column in a generated
// resource is click-to-open (onClick: "link"), and only that one is.
func TestColumnClick_FirstColumnLinks(t *testing.T) {
	def, err := ParseInlineFields("Invoice", "number:string,status:string,total:float")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	g := &Generator{Definition: def}
	content := g.resourceDefinitionFileContent(g.Names())

	// Isolate the columns block (before the form "fields:" key).
	formIdx := strings.Index(content, "fields:")
	if formIdx < 0 {
		t.Fatalf("no form fields block:\n%s", content)
	}
	columns := content[:formIdx]

	if n := strings.Count(columns, `onClick: "link"`); n != 1 {
		t.Fatalf("expected exactly 1 onClick:\"link\" column, got %d:\n%s", n, columns)
	}
	// It must be attached to the first column (number), not a later one.
	numberIdx := strings.Index(columns, `key: "number"`)
	statusIdx := strings.Index(columns, `key: "status"`)
	linkIdx := strings.Index(columns, `onClick: "link"`)
	if numberIdx < 0 || linkIdx < numberIdx || (statusIdx > 0 && linkIdx > statusIdx) {
		t.Errorf("onClick:\"link\" not on the first (number) column:\n%s", columns)
	}
}

// TestColumnClick_BelongsToNotLinked — a leading relationship column is not the
// click target; the first plain column is (linking a related-entity name to
// THIS resource's detail would be misleading).
func TestColumnClick_BelongsToNotLinked(t *testing.T) {
	def, err := ParseInlineFields("Invoice", "customer:belongs_to:Customer,number:string")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	g := &Generator{Definition: def}
	content := g.resourceDefinitionFileContent(g.Names())
	formIdx := strings.Index(content, "fields:")
	columns := content[:formIdx]

	// The customer.name column must not carry the link.
	custLine := ""
	for _, line := range strings.Split(columns, "\n") {
		if strings.Contains(line, `key: "customer.name"`) {
			custLine = line
		}
	}
	if custLine == "" {
		t.Fatalf("expected a customer.name column:\n%s", columns)
	}
	if strings.Contains(custLine, "onClick") {
		t.Errorf("relationship column should not be click-to-open: %s", custLine)
	}
	if !strings.Contains(columns, `onClick: "link"`) {
		t.Errorf("expected the number column to be click-to-open:\n%s", columns)
	}
}
