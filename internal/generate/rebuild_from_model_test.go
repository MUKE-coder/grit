package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The fixture is testdata/invoice_model.go.txt: the literal bytes of a model
// written by a real `grit generate resource` run,
//
//	grit generate resource Invoice --fields "reference:string:unique,slug:slug,
//	  notes:text,total:money,paid:bool,status:select[draft, sent, paid],
//	  contact:email,client:belongs_to[Client]"
//
// not a model written by hand here. The first version of this test was
// hand-written and had money as a decimal float, which is a shape no generated
// project has: a money column is money.Money, embedded, two columns behind one
// field, because one column holding "19.99 USD" is not something you can SUM.
// A fixture that invents the format proves only that the reader handles the
// invention.
//
// It is a .txt rather than a .go so it is not compiled as part of this
// package: it names a module that does not exist here, and it is data.
func writeRebuildProject(t *testing.T) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "invoice_model.go.txt"))
	if err != nil {
		t.Fatalf("reading the generated-model fixture: %v", err)
	}
	root := t.TempDir()
	models := filepath.Join(root, "apps", "api", "internal", "models")
	if err := os.MkdirAll(models, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(models, "invoice.go"), raw, 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// The point of grit#87's second ask: --force should not require retyping the
// field list, because a retyped list is how a column silently changes type
// under a table full of data.
func TestARebuildReadsEveryDeclaredColumn(t *testing.T) {
	root := writeRebuildProject(t)
	def, err := definitionFromModelFileMode(filepath.Join(root, "apps", "api"), "Invoice", modelForRebuild)
	if err != nil {
		t.Fatal(err)
	}

	got := map[string]Field{}
	for _, f := range def.Fields {
		got[f.Name] = f
	}

	// Each of these is a different way the read goes wrong.
	cases := []struct {
		name, wantType string
		check          func(Field) string
	}{
		// text and string are the same Go type and differ only in the column.
		// Rebuilding text as string narrows it, which on Postgres is a
		// migration that can fail on rows that already exist.
		{"notes", "text", nil},
		// money.Money is embedded: two columns, one field, and nothing about
		// the Go type says "money" to a reader that only knows scalars.
		{"total", "money", nil},
		// A slug says in its own tag what it is derived from. Without that it
		// rebuilds as a plain string and stops being a slug at all.
		{"slug", "slug", func(f Field) string {
			if f.SlugSource != "reference" {
				return "slug source is " + f.SlugSource + ", want reference"
			}
			return ""
		}},
		// A choice column's values live in its enum tag.
		{"status", "select", func(f Field) string {
			if len(f.Options) != 3 {
				return "status came back with no options"
			}
			return ""
		}},
		// A formatted column says what it is in its format tag.
		{"contact", "email", nil},
		// The foreign key, spelled from the column, with the model it points
		// at. A seeder skips this and a rebuild must not.
		{"client", "belongs_to", func(f Field) string {
			if f.RelatedModel != "Client" {
				return "related model is " + f.RelatedModel + ", want Client"
			}
			return ""
		}},
		{"reference", "string", func(f Field) string {
			if !f.Unique {
				return "reference lost its unique"
			}
			return ""
		}},
		{"paid", "bool", nil},
	}

	for _, tc := range cases {
		f, ok := got[tc.name]
		if !ok {
			t.Errorf("%s is missing from the rebuilt definition", tc.name)
			continue
		}
		if f.Type != tc.wantType {
			t.Errorf("%s is type %q, want %q", tc.name, f.Type, tc.wantType)
		}
		if tc.check != nil {
			if msg := tc.check(f); msg != "" {
				t.Errorf("%s: %s", tc.name, msg)
			}
		}
	}

	// And nothing the framework owns, which nobody declared and which a
	// rebuild would turn into an ordinary column.
	for _, unwanted := range []string{
		"id", "created_at", "updated_at", "deleted_at", "version",
		"archived_at", "client_id",
	} {
		if _, ok := got[unwanted]; ok {
			t.Errorf("%s is a framework column and should not be in a rebuild", unwanted)
		}
	}
}

// The seeder reads the same file through the same parser and has to keep its
// own answer: no slug, which BeforeCreate computes, and no foreign key, which
// would point at a row that may not exist.
func TestTheSeederKeepsItsOwnSkips(t *testing.T) {
	root := writeRebuildProject(t)
	def, err := definitionFromModelFileMode(filepath.Join(root, "apps", "api"), "Invoice", modelForSeeder)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range def.Fields {
		if f.Name == "slug" {
			t.Error("a seeder would invent a slug that BeforeCreate computes")
		}
		if f.Name == "client" || strings.HasSuffix(f.Name, "_id") {
			t.Errorf("a seeder would invent the foreign key %q", f.Name)
		}
	}
	// It still reads the ordinary columns, so the mode narrowed what it skips
	// rather than breaking it.
	if len(def.Fields) == 0 {
		t.Error("the seeder definition came back empty")
	}
}

// The printed list has to be pasteable: somebody reads it, edits it and hands
// it back as --fields. A list that does not parse is worse than none.
func TestThePrintedFieldListParsesBack(t *testing.T) {
	root := writeRebuildProject(t)
	def, err := definitionFromModelFileMode(filepath.Join(root, "apps", "api"), "Invoice", modelForRebuild)
	if err != nil {
		t.Fatal(err)
	}
	spelled := fieldsSpelling(def)

	reparsed, err := ParseInlineFields("Invoice", spelled)
	if err != nil {
		t.Fatalf("the list printed to the user does not parse back:\n  %s\n  %v", spelled, err)
	}
	if len(reparsed.Fields) != len(def.Fields) {
		t.Errorf("the printed list has %d fields and the definition has %d:\n  %s",
			len(reparsed.Fields), len(def.Fields), spelled)
	}
	for i, f := range reparsed.Fields {
		if f.Name != def.Fields[i].Name || f.Type != def.Fields[i].Type {
			t.Errorf("field %d round-tripped as %s:%s, was %s:%s",
				i, f.Name, f.Type, def.Fields[i].Name, def.Fields[i].Type)
		}
	}
	// The parts most easily lost in a round trip through text.
	for _, want := range []string{
		"notes:text", "total:money", "client:belongs_to[Client]", "reference:string:unique",
	} {
		if !strings.Contains(spelled, want) {
			t.Errorf("the printed list is missing %q:\n  %s", want, spelled)
		}
	}
}

// A model this reader cannot make sense of has to say so rather than rebuild a
// resource with fields missing. Silence here rewrites a model.
func TestARebuildRefusesAModelItCannotRead(t *testing.T) {
	root := t.TempDir()
	models := filepath.Join(root, "apps", "api", "internal", "models")
	if err := os.MkdirAll(models, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(models, "ghost.go"),
		[]byte("package models\n\ntype Ghost struct {\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	def, err := definitionFromModelFileMode(filepath.Join(root, "apps", "api"), "Ghost", modelForRebuild)
	if err == nil && len(def.Fields) > 0 {
		t.Errorf("a model with no columns rebuilt as %d fields", len(def.Fields))
	}
}

// Requiredness is the property a rebuild is most likely to lose quietly, and
// losing it is a migration that drops a NOT NULL nobody asked to drop.
//
// It was lost, on every column, until a byte diff of the rebuilt model against
// the one it had just read showed three fields missing `binding:"required"`
// and a required foreign key rebuilt as *string with default:null. The cause
// is that GoField.Required reads a tag named required:"true", which a
// generated model does not have: it says binding:"required".
func TestARebuildKeepsRequired(t *testing.T) {
	root := writeRebuildProject(t)
	def, err := definitionFromModelFileMode(filepath.Join(root, "apps", "api"), "Invoice", modelForRebuild)
	if err != nil {
		t.Fatal(err)
	}

	required := map[string]bool{}
	for _, f := range def.Fields {
		required[f.Name] = f.Required
	}
	// The fixture declares these three with binding:"required".
	for _, name := range []string{"reference", "status", "client"} {
		if !required[name] {
			t.Errorf("%s lost its required, so a rebuild would relax the column", name)
		}
	}
	// And one that is genuinely optional, so this is not just "everything is
	// required" passing by accident.
	if required["notes"] {
		t.Error("notes is not declared required and came back required")
	}
}

// columnIsRequired has to read the tag a generated model actually writes.
func TestRequiredIsReadFromTheBindingTag(t *testing.T) {
	cases := []struct {
		field GoField
		want  bool
	}{
		{GoField{Binding: "required"}, true},
		{GoField{Binding: "required,email"}, true},
		{GoField{Binding: "omitempty"}, false},
		{GoField{Binding: ""}, false},
		// The other spelling, kept so a model that does carry it still works.
		{GoField{Required: true}, true},
	}
	for _, tc := range cases {
		if got := columnIsRequired(tc.field); got != tc.want {
			t.Errorf("binding %q required=%v: got %v, want %v",
				tc.field.Binding, tc.field.Required, got, tc.want)
		}
	}
}
