package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestAutoFieldParsing verifies the `name:string:auto[:PREFIX]` grammar produces
// an auto-number field: optional (never required), string-typed, prefix captured.
func TestAutoFieldParsing(t *testing.T) {
	def, err := ParseInlineFields("Invoice", "number:string:auto:INV,status:string")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	f := fieldByName(def, "number")
	if !f.IsAuto() {
		t.Fatalf("number is not auto: %+v", f)
	}
	if f.AutoPrefix != "INV" {
		t.Errorf("AutoPrefix = %q, want INV", f.AutoPrefix)
	}
	if f.Required {
		t.Error("auto field must not be required (server fills it in BeforeCreate)")
	}
	if f.GoType() != "string" {
		t.Errorf("auto GoType = %q, want string", f.GoType())
	}
}

// TestAutoFieldParsing_NoPrefix — the prefix segment is optional.
func TestAutoFieldParsing_NoPrefix(t *testing.T) {
	def, err := ParseInlineFields("Invoice", "number:string:auto")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	f := fieldByName(def, "number")
	if !f.IsAuto() {
		t.Fatalf("number is not auto: %+v", f)
	}
	if f.AutoPrefix != "" {
		t.Errorf("AutoPrefix = %q, want empty", f.AutoPrefix)
	}
}

// TestAutoFieldParsing_NonStringRejected — auto is only valid on string fields.
func TestAutoFieldParsing_NonStringRejected(t *testing.T) {
	if _, err := ParseInlineFields("Invoice", "seq:int:auto"); err == nil {
		t.Fatal("expected error for auto on a non-string field, got nil")
	}
}

// TestAutoField_ModelHook — a generated model auto-numbers in BeforeCreate by
// calling sequence.Next directly (not the services wrapper, which would be an
// import cycle), imports internal/sequence, and leaves the column non-required.
func TestAutoField_ModelHook(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Invoice", "number:string:auto:INV,status:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	model := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "invoice.go"))

	// BeforeCreate must set the number only when blank, via sequence.Next.
	assertContains(t, "model", model,
		`if m.Number == ""`,
		`sequence.Next(tx, sequence.Config{`,
		`Prefix: "INV"`,
		module+"/internal/sequence",
	)
	// It must NOT call the services wrapper from a model (import cycle).
	if strings.Contains(model, "services.Next") {
		t.Errorf("model calls services.Next* — that is a models→services import cycle:\n%s", model)
	}
	// The auto column is server-filled, so it must not be binding:"required".
	for _, line := range strings.Split(model, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "Number ") && strings.Contains(line, `binding:"required"`) {
			t.Errorf("auto Number field must not be required: %s", line)
		}
	}

	// The shared sequence package must have been stood up.
	seqPkg := filepath.Join(root, "apps", "api", "internal", "sequence", "sequence.go")
	if got := readTestFile(t, seqPkg); !strings.Contains(got, "func Next(") {
		t.Errorf("internal/sequence/sequence.go missing Next(): %s", seqPkg)
	}
}

// TestAutoField_HiddenFromForm — the auto field is dropped from the admin
// create/edit form (server-filled) but still appears as a table column.
func TestAutoField_HiddenFromForm(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Invoice", "number:string:auto:INV,status:string")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)

	content := g.resourceDefinitionFileContent(g.Names())

	// Split into columns vs form.fields so we can assert per-section.
	formIdx := strings.Index(content, "fields:")
	if formIdx < 0 {
		t.Fatalf("resource definition has no form fields block:\n%s", content)
	}
	columnsPart := content[:formIdx]
	formPart := content[formIdx:]

	if !strings.Contains(columnsPart, `"number"`) {
		t.Errorf("auto field should still appear as a table column:\n%s", columnsPart)
	}
	if strings.Contains(formPart, `key: "number"`) {
		t.Errorf("auto field must NOT appear in the create/edit form:\n%s", formPart)
	}
	// status is a normal field — it must still be in the form.
	if !strings.Contains(formPart, `key: "status"`) {
		t.Errorf("non-auto field status missing from form:\n%s", formPart)
	}
}
