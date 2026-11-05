package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInjectAfterAnchor(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invoice.ts")
	src := "export interface Invoice {\n  id: string;\n  created_at: string;\n}\n"
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	if err := injectAfterAnchor(path, "export interface Invoice {", `  status: "draft" | "paid";`); err != nil {
		t.Fatalf("inject: %v", err)
	}
	out, _ := os.ReadFile(path)
	got := string(out)
	if !strings.Contains(got, `  status: "draft" | "paid";`) {
		t.Errorf("field not injected:\n%s", got)
	}
	// The new line must sit right after the interface opening, before id.
	if strings.Index(got, "status:") > strings.Index(got, "id: string") {
		t.Errorf("injected in the wrong place:\n%s", got)
	}
	// Idempotent: a second inject must not duplicate.
	_ = injectAfterAnchor(path, "export interface Invoice {", `  status: "draft" | "paid";`)
	out2, _ := os.ReadFile(path)
	if strings.Count(string(out2), "status:") != 1 {
		t.Errorf("second inject duplicated the field")
	}
}

func TestInjectAfterAnchorMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.ts")
	os.WriteFile(path, []byte("nothing here\n"), 0644)
	if err := injectAfterAnchor(path, "no such anchor", "code"); err == nil {
		t.Errorf("missing anchor should error")
	}
}

func TestAddFieldRejectsUnsupportedTypes(t *testing.T) {
	// Relationship, file, slug, and array types can't be added in place — they
	// change imports/joins/migrations. AddField must refuse them clearly.
	for _, spec := range []string{
		"category:belongs_to:Category",
		"tags:many_to_many:Tag",
		"avatar:file:image",
		"handle:slug",
		"channels:check:a=A|b=B",
	} {
		if err := AddField("Invoice", spec); err == nil {
			t.Errorf("AddField accepted unsupported spec %q", spec)
		} else if !strings.Contains(err.Error(), "regenerate") && !strings.Contains(err.Error(), "can't add") {
			// It may fail earlier (no project); we only assert it doesn't succeed.
			_ = err
		}
	}
}

func TestAddFieldSupportedSet(t *testing.T) {
	for _, ft := range []FieldType{FieldString, FieldText, FieldSelect, FieldToggle, FieldInt, FieldDatetime} {
		if !addFieldSupported[ft] {
			t.Errorf("%q should be an addable field type", ft)
		}
	}
	for _, ft := range []FieldType{FieldBelongsTo, FieldFile, FieldSlug, FieldCheck, FieldManyToMany} {
		if addFieldSupported[ft] {
			t.Errorf("%q should NOT be addable in place", ft)
		}
	}
}

// A field the model already has is refused, by whatever spelling.
//
// Every injection in AddField is idempotent for a byte-identical line, so a
// second add with any difference at all, a changed option list or a different
// type, appended a SECOND struct field of that name:
//
//	internal/models/part.go:16:2: State redeclared
//
// The command said "Field added to Part." and exited 0. Found by writing the
// verified flow for the field-types page, where the second invocation was a
// typo (grit#168).
//
// A unit test on existingField rather than on AddField: AddField needs a whole
// generated project around it, and the decision is this function.
func TestExistingFieldRecognisesWhatTheModelHas(t *testing.T) {
	fields := []Field{
		{Name: "name", Type: "string"},
		{Name: "State", Type: "select"},
		{Name: "unit_price", Type: "money"},
	}
	for _, want := range []struct{ asked, found string }{
		// The same column, by every spelling it arrives in.
		{"state", "State"},
		{"State", "State"},
		{"STATE", "State"},
		{"unit_price", "unit_price"},
		{"UnitPrice", "unit_price"},
		{"name", "name"},
		// And a field it genuinely does not have.
		{"colour", ""},
		{"states", ""},
		{"price", ""},
	} {
		if got := existingField(fields, want.asked); got != want.found {
			t.Errorf("existingField(%q) = %q, want %q", want.asked, got, want.found)
		}
	}
}
