package generate

import (
	"strings"
	"testing"
)

// The generator and grit sync both write packages/shared/schemas/<name>.ts,
// and they wrote it in two incompatible shapes: the generator keyed by the
// JSON tag, sync keyed in camelCase. Whichever ran last decided whether the
// app worked, because the API decodes with ShouldBindJSON against
// `json:"contact_id"` and ignores a payload keyed contactId (grit#102).
//
// This is the agreement test. It does not care which code path produced a
// key; it cares that both produce the one the API reads.
func TestSyncKeysSchemasByTheJSONTag(t *testing.T) {
	s := GoStruct{
		Name: "Deal",
		Fields: []GoField{
			{Name: "ID", GoType: "string", JSONName: "id"},
			{Name: "Title", GoType: "string", JSONName: "title"},
			{Name: "ContactID", GoType: "*string", JSONName: "contact_id"},
			{Name: "Contact", GoType: "*Contact", JSONName: "contact", GORMTag: "foreignKey:ContactID"},
			{Name: "ClosedAt", GoType: "*jsontime.DateTime", JSONName: "closed_at"},
			{Name: "CreatedAt", GoType: "time.Time", JSONName: "created_at"},
		},
	}
	out := buildZodSchema(s)

	for _, want := range []string{"contact_id:", "closed_at:", "created_at:"} {
		if !strings.Contains(out, want) {
			t.Errorf("the schema does not use the JSON tag %q", want)
		}
	}
	for _, unwanted := range []string{"contactId:", "closedAt:", "createdAt:"} {
		if strings.Contains(out, unwanted) {
			t.Errorf("the schema still writes %q, which the API ignores", unwanted)
		}
	}

	// An association is not a create or update field.
	create := out[strings.Index(out, "CreateDealSchema"):]
	if strings.Contains(create[:strings.Index(create, "})")], "  contact:") {
		t.Error("the create schema includes the association struct")
	}

	// And a date is a string, not an unknown that validates nothing.
	if strings.Contains(out, "closed_at: z.unknown()") {
		t.Error("a jsontime column fell through to z.unknown()")
	}
}

// An optional relationship is optional in the schema too. After the #100 fix
// the Go side was right and the schema still demanded a uuid, so a form built
// from it refused to save a deal with no contact.
func TestGeneratedSchemaMakesAnOptionalRelationshipOptional(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Deal", "title:string")
	if err != nil {
		t.Fatal(err)
	}
	def.Fields = append(def.Fields,
		Field{Name: "contact", Type: "belongs_to", RelatedModel: "Contact"},
		Field{Name: "company", Type: "belongs_to", RelatedModel: "Company", Required: true})

	g := newTestGenerator(root, module, def)
	// The model first: the schema is written from it rather than from the
	// definition, which is the order Generator.Run uses.
	if err := g.writeGoModel(g.Names()); err != nil {
		t.Fatalf("writeGoModel: %v", err)
	}
	if err := g.writeZodSchema(g.Names()); err != nil {
		t.Fatal(err)
	}
	out := readTestFile(t, root+"/packages/shared/schemas/deal.ts")

	if !strings.Contains(out, "contact_id: z.string().uuid(\"Invalid ID\").nullable().optional()") {
		t.Errorf("an optional relationship should be nullable and optional:\n%s", out)
	}
	if strings.Contains(out, "company_id: z.string().uuid(\"Invalid ID\").nullable().optional()") {
		t.Error("a required relationship was made optional")
	}
}

// "Required while it is shown" is a predicate, not a flag. Marked required in
// the object it rejects every row the condition does not hold for: a
// lost_reason required when status is lost made every open deal invalid.
func TestConditionalRequiredBecomesARefinement(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Deal", "title:string")
	if err != nil {
		t.Fatal(err)
	}
	// The field the condition reads has to be in the definition. It was not,
	// and this test asserted a refinement reading a key the schema does not
	// have, which is a TypeScript error: the output it demanded never
	// compiled. That is the bug grit#115 turned out to be, in miniature.
	def.Fields = append(def.Fields,
		Field{
			Name:    "status",
			Type:    "select",
			Options: []FieldOption{{Value: "open", Label: "Open"}, {Value: "lost", Label: "Lost"}},
		},
		Field{
			Name:     "lost_reason",
			Type:     "string",
			Required: true,
			ShowWhen: &FieldCondition{Field: "status", Value: "lost"},
		})

	g := newTestGenerator(root, module, def)
	// The model first: the schema is written from it rather than from the
	// definition, which is the order Generator.Run uses.
	if err := g.writeGoModel(g.Names()); err != nil {
		t.Fatalf("writeGoModel: %v", err)
	}
	if err := g.writeZodSchema(g.Names()); err != nil {
		t.Fatal(err)
	}
	out := readTestFile(t, root+"/packages/shared/schemas/deal.ts")

	if !strings.Contains(out, "lost_reason: z.string().min(1, \"Required\").optional()") {
		t.Errorf("a conditional field must be optional in the object:\n%s", out)
	}
	for _, want := range []string{
		".superRefine((v, ctx) => {",
		`if (v.status === "lost" && !v.lost_reason)`,
		`path: ["lost_reason"]`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the refinement is missing %q", want)
		}
	}
}

// And a resource with no conditional fields gets no refinement, or every
// schema carries an empty closure.
func TestNoRefinementWithoutAConditionalField(t *testing.T) {
	if got := zodRefinements(nil); got != "" {
		t.Errorf("expected no refinement, got %q", got)
	}
}
