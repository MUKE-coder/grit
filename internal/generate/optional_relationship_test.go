package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// A belongs_to ignored its own required flag and always came out
// binding:"required" (grit#100). So a deal with no contact yet could not be
// saved, and the workaround, storing "", is refused by Postgres with
// SQLSTATE 23503 the moment the foreign key constraint exists.
//
// Six surfaces decide this and they used to decide it separately. These
// assertions are about agreement between them as much as about any one.
func TestOptionalBelongsToIsOptionalEverywhere(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Deal", "title:string")
	if err != nil {
		t.Fatal(err)
	}
	// One required relationship and one that is not, in the same resource, so
	// a fix that simply made everything optional fails here too.
	def.Fields = append(def.Fields,
		Field{Name: "contact", Type: "belongs_to", RelatedModel: "Contact"},
		Field{Name: "company", Type: "belongs_to", RelatedModel: "Company", Required: true},
	)

	g := newTestGenerator(root, module, def)
	names := g.Names()
	if err := g.writeGoModel(names); err != nil {
		t.Fatalf("writeGoModel: %v", err)
	}
	model := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "deal.go"))
	if !strings.Contains(model, `ContactID *string`) {
		t.Error("an optional belongs_to should be a *string, so the column can hold NULL")
	}
	if !strings.Contains(model, "default:null") {
		t.Error("an optional foreign key needs default:null, or the column defaults to ''")
	}
	if strings.Contains(model, `json:"contact_id" binding:"required"`) {
		t.Error("an optional belongs_to still carries binding:\"required\"")
	}
	if !strings.Contains(model, `json:"company_id" binding:"required"`) {
		t.Error("a required belongs_to lost its binding:\"required\"")
	}
	if !strings.Contains(model, `CompanyID string`) {
		t.Error("a required belongs_to should stay a plain string")
	}
}

// relOptional is the one predicate the model, the handler, the TypeScript
// type, the Zod schema, the admin form and the importer all read. They each
// decided it for themselves before, and three said required while the fourth
// said nullable.
func TestRelOptional(t *testing.T) {
	def, err := ParseInlineFields("Category", "name:string")
	if err != nil {
		t.Fatal(err)
	}
	g := &Generator{Definition: def}

	cases := []struct {
		name string
		f    Field
		want bool
	}{
		{"a self-reference is always optional, a tree root has no parent",
			Field{Name: "parent", Type: "belongs_to", RelatedModel: "Category"}, true},
		{"a self-reference stays optional even when marked required",
			Field{Name: "parent", Type: "belongs_to", RelatedModel: "Category", Required: true}, true},
		{"an unmarked relationship is optional, like every other field type",
			Field{Name: "owner", Type: "belongs_to", RelatedModel: "User"}, true},
		{"a required relationship is not",
			Field{Name: "owner", Type: "belongs_to", RelatedModel: "User", Required: true}, false},
	}
	for _, c := range cases {
		if got := g.relOptional(c.f); got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

// The handler calls optional<X>ID for a nullable key, and the helper was only
// emitted for a self-reference. The first optional relationship produced a
// handler calling a function that was never written.
func TestOptionalIDHelperIsEmittedForAnyOptionalRelationship(t *testing.T) {
	def, err := ParseInlineFields("Deal", "title:string")
	if err != nil {
		t.Fatal(err)
	}
	def.Fields = append(def.Fields,
		Field{Name: "contact", Type: "belongs_to", RelatedModel: "Contact"})

	g := &Generator{Definition: def, Module: "myapp/apps/api"}
	parts := g.crud(MakeNames("Deal"))
	if !strings.Contains(parts.optionalID, "func optionalDealID(") {
		t.Error("a resource with an optional relationship needs the optionalID helper")
	}

	// And it stays conditional: an unused function does not compile.
	required, err := ParseInlineFields("Invoice", "number:string")
	if err != nil {
		t.Fatal(err)
	}
	required.Fields = append(required.Fields,
		Field{Name: "client", Type: "belongs_to", RelatedModel: "Client", Required: true})
	g2 := &Generator{Definition: required, Module: "myapp/apps/api"}
	if parts := g2.crud(MakeNames("Invoice")); parts.optionalID != "" {
		t.Error("a resource whose relationships are all required should not get the helper")
	}
}

// An optional foreign key is a *string, and the Zod generated from it must say
// so. It fell through to z.unknown(), so the one surface that should catch a
// bad relationship id validated nothing at all.
func TestPointerColumnsBecomeNullableZod(t *testing.T) {
	if got := goTypeToZod("*string", "size:36;index;default:null"); got != "z.string().nullable().optional()" {
		t.Errorf("a *string column should be nullable and optional in Zod, got %s", got)
	}
	if got := goTypeToZod("string", "size:36;index"); got != "z.string()" {
		t.Errorf("a plain string should stay required, got %s", got)
	}
}
