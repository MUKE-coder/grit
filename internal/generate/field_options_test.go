package generate

import "testing"

func fieldByName(def *ResourceDefinition, name string) Field {
	for _, f := range def.Fields {
		if f.Name == name {
			return f
		}
	}
	return Field{}
}

func TestSelectFieldParsing(t *testing.T) {
	def, err := ParseInlineFields("Invoice", "status:select:draft=Draft|sent=Sent|paid=Paid")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	f := fieldByName(def, "status")
	if !f.IsSelect() {
		t.Fatalf("status not a select: %q", f.Type)
	}
	if len(f.Options) != 3 {
		t.Fatalf("options = %d, want 3", len(f.Options))
	}
	if f.Options[0].Value != "draft" || f.Options[0].Label != "Draft" {
		t.Errorf("option 0 = %+v", f.Options[0])
	}
	if f.GoType() != "string" {
		t.Errorf("select GoType = %q, want string", f.GoType())
	}
	if got := f.ZodType(); got != `z.enum(["draft", "sent", "paid"])` {
		t.Errorf("select ZodType = %q", got)
	}
	if got := f.TSType(); got != `"draft" | "sent" | "paid"` {
		t.Errorf("select TSType = %q", got)
	}
	if f.FormFieldType() != "select" {
		t.Errorf("select FormFieldType = %q", f.FormFieldType())
	}
}

func TestCheckFieldParsing(t *testing.T) {
	def, err := ParseInlineFields("Post", "tags:check:news=News|ops=Ops")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	f := fieldByName(def, "tags")
	if !f.IsCheck() {
		t.Fatalf("tags not a check: %q", f.Type)
	}
	if f.GoType() != "datatypes.JSONSlice[string]" {
		t.Errorf("check GoType = %q", f.GoType())
	}
	if !f.NeedsDatatypesImport() {
		t.Errorf("check must pull in the datatypes import")
	}
	if got := f.ZodType(); got != `z.array(z.enum(["news", "ops"])).optional()` {
		t.Errorf("check ZodType = %q", got)
	}
	if got := f.TSType(); got != `("news" | "ops")[]` {
		t.Errorf("check TSType = %q", got)
	}
	if f.FormFieldType() != "checkbox-group" {
		t.Errorf("check FormFieldType = %q", f.FormFieldType())
	}
}

func TestToggleField(t *testing.T) {
	def, err := ParseInlineFields("Task", "active:toggle")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	f := fieldByName(def, "active")
	if !f.IsToggle() || f.GoType() != "bool" || f.FormFieldType() != "toggle" {
		t.Errorf("toggle mapping wrong: go=%q form=%q", f.GoType(), f.FormFieldType())
	}
	if f.ZodType() != "z.boolean().optional()" {
		t.Errorf("toggle ZodType = %q", f.ZodType())
	}
}

func TestBareOptionValueHumanized(t *testing.T) {
	def, err := ParseInlineFields("Order", "state:select:in_progress|done")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	f := fieldByName(def, "state")
	if f.Options[0].Value != "in_progress" || f.Options[0].Label != "In Progress" {
		t.Errorf("bare value not humanized: %+v", f.Options[0])
	}
	if f.Options[1].Label != "Done" {
		t.Errorf("option 1 label = %q", f.Options[1].Label)
	}
}

func TestSelectRequiresOptions(t *testing.T) {
	if _, err := ParseInlineFields("X", "status:select"); err == nil {
		t.Errorf("select with no options should error")
	}
}

func TestDuplicateOptionRejected(t *testing.T) {
	if _, err := ParseInlineFields("X", "s:select:a=A|a=B"); err == nil {
		t.Errorf("duplicate option value should error")
	}
}

func TestOptionsLiteral(t *testing.T) {
	def, _ := ParseInlineFields("Invoice", "status:select:draft=Draft|paid=Paid")
	f := fieldByName(def, "status")
	if got := f.OptionsLiteral(); got != `[{ value: "draft", label: "Draft" }, { value: "paid", label: "Paid" }]` {
		t.Errorf("OptionsLiteral = %q", got)
	}
}

func TestRadioField(t *testing.T) {
	def, err := ParseInlineFields("Survey", "rating:radio:low=Low|high=High")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	f := fieldByName(def, "rating")
	if !f.IsRadio() {
		t.Fatalf("rating not a radio: %q", f.Type)
	}
	// radio is single-choice like select — string, enum, union — but renders as
	// radio buttons, not a dropdown.
	if f.GoType() != "string" {
		t.Errorf("radio GoType = %q, want string", f.GoType())
	}
	if got := f.ZodType(); got != `z.enum(["low", "high"])` {
		t.Errorf("radio ZodType = %q", got)
	}
	if f.FormFieldType() != "radio" {
		t.Errorf("radio FormFieldType = %q, want radio", f.FormFieldType())
	}
	if len(f.Options) != 2 || f.Options[0].Label != "Low" {
		t.Errorf("radio options = %+v", f.Options)
	}
}
