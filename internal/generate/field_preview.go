package generate

import (
	"sort"
	"strings"
)

// What one field type produces across the stack, for the field builder on the
// website.
//
// ── Why this is computed rather than written down ────────────────────────
//
// The field builder used to show only the command it was building. The
// question it could not answer is the one people actually have, which is
// "what do I get if I write this": the Go column, the TypeScript, the Zod
// rule, the admin control.
//
// Writing that out as a table in a TSX file would be a second description of
// what the generator does, and the day somebody changes GoType or
// FormFieldType the website would start describing a framework that no longer
// exists, with nothing failing. So the table is built here, by asking the same
// methods the generator asks, and TestFieldPreviewDataIsCurrent fails when the
// committed copy drifts from it.
//
// It is still only a preview of one field. The whole file a resource produces
// is a different question (grit#256) and this does not pretend to answer it.

// FieldShape is what one field becomes, across the stack.
type FieldShape struct {
	// Go is the struct field's type, as the model declares it.
	Go string `json:"go"`
	// GORM is the struct tag that decides the column, empty when the field
	// needs none.
	GORM string `json:"gorm"`
	// TS is the TypeScript type in packages/shared.
	TS string `json:"ts"`
	// Zod is the validation rule the create and update schemas carry.
	Zod string `json:"zod"`
	// Control is the admin form control, empty for a column the form does not
	// show because the server fills it.
	Control string `json:"control"`
}

// FieldPreview is one field type and what the generator makes of it, in each
// combination of the two modifiers that change the answer.
//
// required and unique are not decoration: Required adds .min(1, "Required") to
// a string's Zod rule and drops the .optional(), and Unique adds uniqueIndex
// to the column. A preview that ignored them would show .optional() beside a
// ticked "required" box, which is a confident interface being wrong, and the
// first version of this did exactly that.
//
// The three variants are omitted when they are identical to Base, which most
// types are for Unique, so the generated file stays readable.
type FieldPreview struct {
	// Type is the name written in --fields, such as "money".
	Type string `json:"type"`
	// Base is the field with neither modifier.
	Base FieldShape `json:"base"`
	// Required, Unique and RequiredUnique are nil when they match Base.
	Required       *FieldShape `json:"required,omitempty"`
	Unique         *FieldShape `json:"unique,omitempty"`
	RequiredUnique *FieldShape `json:"requiredUnique,omitempty"`
	// Note is the one thing worth knowing that the shape does not say. Empty
	// for the types that hold no surprise.
	Note string `json:"note,omitempty"`
}

// previewNotes are the surprises. Each one is a thing somebody has to find out
// some other way, and most of them were found the hard way by somebody.
var previewNotes = map[string]string{
	"money":        "Two columns: the amount in minor units and the currency. Never a float.",
	"slug":         "Filled by the server from another column, so it is not on the form.",
	"vector":       "The dimension is part of the column type, so it cannot be changed later without a migration.",
	"file":         "Stores a reference, not the bytes. Deleting the row releases it to the orphan sweep.",
	"files":        "A list of references. The same release rules as file.",
	"richtext":     "Sanitised on write, so what is stored is already safe to render.",
	"belongs_to":   "Writes the foreign key column and the association, and the admin gets a picker.",
	"many_to_many": "Writes a join table. The admin gets a multi-select.",
	"select":       "The choices are a database CHECK as well as a dropdown. The union and the enum above show an example; yours take their place.",
	"radio":        "The same column as select, with a different control. The union and the enum above show an example; your choices take their place.",
	"check":        "Several choices in one JSON column, not a join table. The union and the enum above show an example; yours take their place.",
	"json":         "Stored as JSON and typed as unknown, so nothing downstream assumes a shape.",
	"string_array": "One column holding a list, not a join table.",
	"date":         "A day with no time, so it does not move when a reader is in another timezone.",
}

// previewSample is a field of the given type, filled in enough to be rendered.
//
// Several types need more than a name to mean anything: a select with no
// choices, a belongs_to with no model and a vector with no dimension are not
// fields anybody could write, and asking the generator about them would
// describe a shape it never emits.
func previewSample(fieldType string) Field {
	f := Field{Name: "example", Type: fieldType}
	switch FieldType(fieldType) {
	case FieldSelect, FieldRadio, FieldCheck:
		f.Options = []FieldOption{{Value: "draft", Label: "Draft"}, {Value: "sent", Label: "Sent"}}
	case FieldBelongsTo, FieldOneToOne, FieldManyToMany:
		f.RelatedModel = "User"
	case FieldVector:
		// The dimension lives in Max, which is where VectorDim reads it.
		f.Max = 1536
	case FieldSlug:
		f.SlugSource = "title"
	}
	return f
}

// shapeOf asks the generator what this field becomes.
func shapeOf(f Field) FieldShape {
	return FieldShape{
		Go:      f.GoType(),
		GORM:    f.GORMTag(),
		TS:      f.TSType(),
		Zod:     f.ZodType(),
		Control: f.FormFieldType(),
	}
}

// variant returns the shape with those modifiers, or nil when it is the same
// as base and the generated file is better off not repeating it.
func variant(t string, base FieldShape, required, unique bool) *FieldShape {
	f := previewSample(t)
	f.Required = required
	f.Unique = unique
	got := shapeOf(f)
	if got == base {
		return nil
	}
	return &got
}

// FieldPreviews is every field type and what it produces, sorted by name so
// the generated file does not reshuffle between runs.
func FieldPreviews() []FieldPreview {
	types := ValidFieldTypes()
	sort.Strings(types)

	out := make([]FieldPreview, 0, len(types))
	for _, t := range types {
		base := shapeOf(previewSample(t))
		out = append(out, FieldPreview{
			Type:           t,
			Base:           base,
			Required:       variant(t, base, true, false),
			Unique:         variant(t, base, false, true),
			RequiredUnique: variant(t, base, true, true),
			Note:           previewNotes[t],
		})
	}
	return out
}

// FieldPreviewTS renders the table as the TypeScript module the website reads.
//
// Written rather than served, because the field builder is a static page and a
// request for this would be a request that can fail on a page whose whole job
// is to answer instantly.
func FieldPreviewTS() string {
	var b strings.Builder
	b.WriteString(`// Code generated from internal/generate. DO NOT EDIT.
//
// What each field type produces across the stack, asked of the same methods
// the generator uses: Field.GoType, GORMTag, TSType, ZodType and
// FormFieldType. Regenerate with:
//
//	go test ./internal/generate -run TestFieldPreviewDataIsCurrent -update
//
// TestFieldPreviewDataIsCurrent fails when this file drifts from them, so the
// field builder cannot describe a framework that no longer exists.
//
// required and unique each get their own shape, because they change the
// answer: Required adds .min(1, "Required") to a string and drops the
// .optional(), and Unique adds uniqueIndex to the column. They are left out
// when they are the same as base.

export type FieldShape = {
  /** The struct field's type, as the model declares it. */
  go: string
  /** The struct tag that decides the column, empty when it needs none. */
  gorm: string
  /** The TypeScript type in packages/shared. */
  ts: string
  /** The rule the create and update schemas carry. */
  zod: string
  /** The admin form control, empty when the server fills the column. */
  control: string
}

export type FieldPreview = {
  /** The name written in --fields, such as "money". */
  type: string
  /** With neither modifier. */
  base: FieldShape
  /** Absent when the modifier changes nothing. */
  required?: FieldShape
  unique?: FieldShape
  requiredUnique?: FieldShape
  /** The one thing the shape does not say. */
  note?: string
}

/** The shape for a field with these modifiers, falling back to base. */
export function shapeFor(
  preview: FieldPreview,
  modifiers: { required?: boolean; unique?: boolean },
): FieldShape {
  const { required, unique } = modifiers
  if (required && unique) return preview.requiredUnique ?? preview.base
  if (required) return preview.required ?? preview.base
  if (unique) return preview.unique ?? preview.base
  return preview.base
}

export const FIELD_PREVIEWS: Record<string, FieldPreview> = {
`)
	for _, p := range FieldPreviews() {
		b.WriteString("  " + tsKey(p.Type) + ": {\n")
		b.WriteString("    type: " + tsString(p.Type) + ",\n")
		writeShape(&b, "base", &p.Base)
		writeShape(&b, "required", p.Required)
		writeShape(&b, "unique", p.Unique)
		writeShape(&b, "requiredUnique", p.RequiredUnique)
		if p.Note != "" {
			b.WriteString("    note: " + tsString(p.Note) + ",\n")
		}
		b.WriteString("  },\n")
	}
	b.WriteString("}\n")
	return b.String()
}

// writeShape emits one shape, or nothing when it is nil.
func writeShape(b *strings.Builder, key string, shape *FieldShape) {
	if shape == nil {
		return
	}
	b.WriteString("    " + key + ": {\n")
	b.WriteString("      go: " + tsString(shape.Go) + ",\n")
	b.WriteString("      gorm: " + tsString(shape.GORM) + ",\n")
	b.WriteString("      ts: " + tsString(shape.TS) + ",\n")
	b.WriteString("      zod: " + tsString(shape.Zod) + ",\n")
	b.WriteString("      control: " + tsString(shape.Control) + ",\n")
	b.WriteString("    },\n")
}

// tsString is a single-quoted TypeScript literal, which is what the rest of
// docs/config is written in.
func tsString(s string) string {
	return "'" + strings.NewReplacer(`\`, `\\`, `'`, `\'`, "\n", `\n`).Replace(s) + "'"
}
