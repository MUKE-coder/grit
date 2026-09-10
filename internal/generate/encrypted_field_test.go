package generate

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The :encrypted modifier makes crypto.EncryptedString the field's own type.
//
// Following the docs by hand, changing a generated model's field to
// crypto.EncryptedString, broke the build in the handler and the CSV importer,
// and once patched, PUT and PATCH stored plaintext. Found on a patient-records
// app: POST wrote enc:v1: ciphertext, the next edit wrote "PUT: diabetic".

func TestEncryptedModifierParses(t *testing.T) {
	def, err := ParseInlineFields("Patient", "name:string,notes:text:encrypted")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if !def.Fields[1].Encrypted {
		t.Fatal("notes:text:encrypted did not mark the field encrypted")
	}
	for spec, why := range map[string]string{
		"age:int:encrypted":           "an int is not text",
		"ssn:string:unique:encrypted": "ciphertext differs on every write, so unique can never hold",
		"ssn:string:encrypted:unique": "the same, whatever the order",
		"notes:text:encrpyted":        "a typo must not be ignored",
	} {
		if _, err := ParseInlineFields("Patient", spec); err == nil {
			t.Errorf("%s was accepted (%s)", spec, why)
		}
	}
}

func TestEncryptedFieldIsTypedThroughout(t *testing.T) {
	const module = "clinic/apps/api"
	root := setupMinimalProject(t, module)
	g := newTestGenerator(root, module, mustFields(t, "Patient", "name:string,notes:text:encrypted"))
	names := g.Names()

	if err := g.writeGoModel(names); err != nil {
		t.Fatalf("model: %v", err)
	}
	if err := g.writeGoHandler(names); err != nil {
		t.Fatalf("handler: %v", err)
	}
	if err := g.writeGoImportHandler(names); err != nil {
		t.Fatalf("importer: %v", err)
	}
	api := filepath.Join(root, "apps", "api", "internal")
	model := readTestFile(t, filepath.Join(api, "models", "patient.go"))
	handler := readTestFile(t, filepath.Join(api, "handlers", "patient.go"))
	importer := readTestFile(t, filepath.Join(api, "handlers", "patient_import.go"))

	for name, src := range map[string]string{"model": model, "handler": handler, "importer": importer} {
		if !strings.Contains(src, `"clinic/apps/api/internal/crypto"`) {
			t.Errorf("the %s uses crypto.EncryptedString without importing it", name)
		}
		if _, err := parser.ParseFile(token.NewFileSet(), name+".go", src, parser.AllErrors); err != nil {
			t.Errorf("the %s does not parse: %v", name, err)
		}
	}

	if !regexp.MustCompile(`Notes\s+crypto\.EncryptedString\s+` + "`" + `gorm:"type:text"`).MatchString(model) {
		t.Error("the model column is not an EncryptedString stored as text")
	}
	if !regexp.MustCompile(`(?s)type CreatePatientRequest struct \{[^}]*Notes\s+crypto\.EncryptedString`).MatchString(handler) {
		t.Error("the create request still takes a plain string, which does not compile into the model")
	}
	// The update goes through a typed pointer, so the value reaching the map
	// is an EncryptedString and Value() encrypts it.
	if !strings.Contains(handler, `updates["notes"] = *req.Notes`) {
		t.Error("the update does not write the typed value")
	}
	if !strings.Contains(importer, "item.Notes = crypto.EncryptedString(v)") {
		t.Error("the CSV importer assigns a plain string to an encrypted field")
	}

	// Ciphertext differs on every write, so it cannot be searched, sorted or
	// filtered, and offering it would return nothing or noise.
	for _, col := range []string{`"notes"`} {
		if strings.Contains(g.buildHandlerSearchCols(), col) {
			t.Error("an encrypted column is searchable")
		}
	}
	if strings.Contains(g.buildSortableSet(), `"notes"`) {
		t.Error("an encrypted column is sortable")
	}
	if regexp.MustCompile(`Filterable: map\[string\]bool\{[^}]*"notes"`).MatchString(handler) {
		t.Error("an encrypted column is filterable")
	}
}
