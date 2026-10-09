package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The owner of an --owned-by row is whoever is signed in.
//
// The service says so in a comment and enforces it in code: it sets the owner
// from the request context on create and never reads it from the body. The
// generated admin form did not know that. It put a relationship picker on the
// create and edit screens, marked it required, and refused to submit until an
// operator chose a user, whose choice was then discarded.
//
// Verified against a running project: picking Jane Cooper stored a row owned
// by admin@example.com, the account that was signed in. Every unit test in
// this package passed while that was true, because none of them read the form.

// ownedResource is a Contact with an owner, as --owned-by user generates it.
func ownedResource(t *testing.T) (*Generator, Names, string) {
	t.Helper()
	root := t.TempDir()

	def, err := ParseInlineFields("Contact", "first_name:string,company:string,group:belongs_to:Group,user:belongs_to:User")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	def.OwnedBy = "user"

	g := newTestGenerator(root, "testapp/apps/api", def)
	if g.Definition.OwnerField() == nil {
		t.Fatal("the definition has no owner field, so this test proves nothing")
	}
	return g, g.Names(), root
}

func TestTheOwnerIsNotAFormField(t *testing.T) {
	g, names, root := ownedResource(t)

	if err := g.writeResourceDefinition(names); err != nil {
		t.Fatalf("writeResourceDefinition: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "admin", "resources", "contacts", "contacts.ts"))

	form := betweenMarkers(t, src, "grit:fields:auto-start", "grit:fields:auto-end")
	if strings.Contains(form, "user_id") {
		t.Errorf("the create and edit form asks for user_id, and the service overwrites it "+
			"with the signed-in user. The field is required, so the form cannot be "+
			"submitted without an answer that is then thrown away.\n\nform fields:\n%s", form)
	}
	// The other relationship is a real choice and stays.
	if !strings.Contains(form, "group_id") {
		t.Error("group_id went missing from the form; only the owner should be dropped")
	}
}

// The column stays. Who owns a row is worth seeing, and a table is a read.
func TestTheOwnerIsStillATableColumn(t *testing.T) {
	g, names, root := ownedResource(t)

	if err := g.writeResourceDefinition(names); err != nil {
		t.Fatalf("writeResourceDefinition: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "admin", "resources", "contacts", "contacts.ts"))

	cols := betweenMarkers(t, src, "grit:cols:auto-start", "grit:cols:auto-end")
	if !strings.Contains(cols, `label: "User"`) {
		t.Errorf("the owner is not a column any more:\n%s", cols)
	}
}

func TestTheOwnerIsNotInTheZodSchemas(t *testing.T) {
	g, names, root := ownedResource(t)

	if err := g.writeZodSchema(names); err != nil {
		t.Fatalf("writeZodSchema: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "packages", "shared", "schemas", "contact.ts"))

	if strings.Contains(src, "user_id") {
		t.Errorf("the create schema requires user_id on a resource whose owner is "+
			"server-assigned, so anything validating against it has to invent a "+
			"value that is discarded.\n\n%s", src)
	}
	if !strings.Contains(src, "group_id") {
		t.Error("group_id went missing from the schema; only the owner should be dropped")
	}
}

// A resource with no owner keeps every relationship it declared.
func TestASharedResourceKeepsAllItsRelationships(t *testing.T) {
	root := t.TempDir()
	def, err := ParseInlineFields("Note", "body:text,user:belongs_to:User")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, "testapp/apps/api", def)
	names := g.Names()

	if err := g.writeResourceDefinition(names); err != nil {
		t.Fatalf("writeResourceDefinition: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "admin", "resources", "notes", "notes.ts"))

	form := betweenMarkers(t, src, "grit:fields:auto-start", "grit:fields:auto-end")
	if !strings.Contains(form, "user_id") {
		t.Errorf("a shared resource lost its user relationship; without --owned-by "+
			"the user is an ordinary choice the form should make.\n%s", form)
	}
}

// betweenMarkers returns the text between two anchors.
func betweenMarkers(t *testing.T, src, start, end string) string {
	t.Helper()
	i := strings.Index(src, start)
	j := strings.Index(src, end)
	if i < 0 || j < 0 || j < i {
		t.Fatalf("markers %q and %q not found in order", start, end)
	}
	return src[i+len(start) : j]
}

// A belongs_to column names the field the related model is actually labelled
// by. It used to hardcode ".name", so a column pointing at any model without
// one rendered a dash: User has first_name, an Invoice has a number, a
// Collection has a title.
func TestABelongsToColumnUsesTheRelatedModelsLabel(t *testing.T) {
	root := t.TempDir()
	models := filepath.Join(root, "apps", "api", "internal", "models")
	if err := os.MkdirAll(models, 0o755); err != nil {
		t.Fatal(err)
	}
	// A related model whose label is a title, not a name.
	collection := "package models\n\ntype Collection struct {\n\tID    string `json:\"id\"`\n\tTitle string `gorm:\"size:255\" json:\"title\"`\n}\n"
	if err := os.WriteFile(filepath.Join(models, "collection.go"), []byte(collection), 0o644); err != nil {
		t.Fatal(err)
	}

	def, err := ParseInlineFields("Item", "name:string,collection:belongs_to:Collection")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, "testapp/apps/api", def)
	names := g.Names()

	if err := g.writeResourceDefinition(names); err != nil {
		t.Fatalf("writeResourceDefinition: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "admin", "resources", "items", "items.ts"))

	if !strings.Contains(src, `key: "collection.title"`) {
		cols := betweenMarkers(t, src, "grit:cols:auto-start", "grit:cols:auto-end")
		t.Errorf("the column does not read collection.title, so it renders empty for "+
			"every row:\n%s", cols)
	}
}
