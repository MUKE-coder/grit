package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// grit#131: every organization drew auto numbers from one counter.
//
// sequence.Counter is keyed by (name, bucket) and the hook passed a fixed
// name, so organization A got 0001, B got 0002, and A's next was 0003. A
// sees a gap it cannot explain, and each organization can work out how many
// numbers the others issue. Gap-free numbering per business is a legal
// requirement in many countries, which is most of what auto is for.
func TestAutoFieldOnATenantOwnedResourceCountsPerOrganization(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Invoice", "title:string,number:string:auto:INV")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	def.TenantOwned = true

	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	model := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "invoice.go"))

	assertContains(t, "model", model,
		// The organization comes off the context, not off the struct: the
		// stamping hook's order relative to this one is not something to
		// depend on.
		`tenant.FromContext(tx.Statement.Context)`,
		// And it is part of the counter's name.
		`Name:   "invoice_number" + ":" + orgID,`,
		// The refusal needs fmt, and a missing import is a model that does
		// not compile.
		`"fmt"`,
	)

	// It fails closed. A number drawn from another organization's counter is
	// worse than a write that did not happen.
	if !strings.Contains(model, "if !ok {") {
		t.Error("the hook does not refuse the write when there is no organization " +
			"on the context, so it falls back to a shared counter (grit#131)")
	}
}

// A resource that is not tenant-owned keeps the plain counter: there is no
// organization to scope by, and adding one would be a name nothing matches.
func TestAutoFieldOnAPlainResourceIsUnchanged(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Invoice", "title:string,number:string:auto:INV")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	model := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "invoice.go"))

	if !strings.Contains(model, `Name:   "invoice_number",`) {
		t.Error("a resource that is not tenant-owned no longer uses the plain counter name")
	}
	if strings.Contains(model, "tenant.FromContext") {
		t.Error("a resource that is not tenant-owned reads an organization it does not have")
	}
}
