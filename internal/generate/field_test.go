package generate

import (
	"strings"
	"testing"
)

// ── GoType ───────────────────────────────────────────────────────────────────

func TestField_GoType(t *testing.T) {
	tests := []struct {
		fieldType string
		want      string
	}{
		{"string", "string"},
		{"text", "string"},
		{"richtext", "string"},
		{"slug", "string"},
		{"int", "int"},
		{"uint", "uint"},
		{"belongs_to", "uint"},
		{"float", "float64"},
		{"bool", "bool"},
		{"datetime", "*time.Time"},
		{"date", "*time.Time"},
		{"many_to_many", "[]uint"},
		{"string_array", "datatypes.JSONSlice[string]"},
		{"unknown", "string"}, // fallback
	}

	for _, tt := range tests {
		f := Field{Type: tt.fieldType}
		if got := f.GoType(); got != tt.want {
			t.Errorf("Field{Type:%q}.GoType() = %q, want %q", tt.fieldType, got, tt.want)
		}
	}
}

// ── TSType ───────────────────────────────────────────────────────────────────

func TestField_TSType(t *testing.T) {
	tests := []struct {
		fieldType string
		want      string
	}{
		{"string", "string"},
		{"text", "string"},
		{"slug", "string"},
		{"richtext", "string"},
		{"int", "number"},
		{"uint", "number"},
		{"float", "number"},
		{"belongs_to", "number"},
		{"bool", "boolean"},
		{"datetime", "string | null"},
		{"date", "string | null"},
		{"many_to_many", "number[]"},
		{"string_array", "string[]"},
		{"unknown", "string"}, // fallback
	}

	for _, tt := range tests {
		f := Field{Type: tt.fieldType}
		if got := f.TSType(); got != tt.want {
			t.Errorf("Field{Type:%q}.TSType() = %q, want %q", tt.fieldType, got, tt.want)
		}
	}
}

// ── ZodType ──────────────────────────────────────────────────────────────────

func TestField_ZodType(t *testing.T) {
	tests := []struct {
		field       Field
		wantContain string // substring the ZodType must contain
	}{
		{Field{Type: "string", Required: true}, `z.string()`},
		{Field{Type: "string", Required: true}, `min(1`},
		{Field{Type: "string", Required: false}, `.optional()`},
		{Field{Type: "int"}, "z.number().int()"},
		{Field{Type: "float"}, "z.number()"},
		{Field{Type: "bool"}, "z.boolean()"},
		{Field{Type: "datetime"}, "z.string().nullable()"},
		{Field{Type: "belongs_to"}, "z.number().int()"},
		{Field{Type: "many_to_many"}, "z.array(z.number().int())"},
		{Field{Type: "string_array"}, "z.array(z.string())"},
	}

	for _, tt := range tests {
		got := tt.field.ZodType()
		if !strings.Contains(got, tt.wantContain) {
			t.Errorf("Field{Type:%q,Required:%v}.ZodType() = %q, does not contain %q",
				tt.field.Type, tt.field.Required, got, tt.wantContain)
		}
	}
}

// ── GORMTag ──────────────────────────────────────────────────────────────────

func TestField_GORMTag(t *testing.T) {
	tests := []struct {
		field       Field
		wantContain string // "" means we don't care about content, just no panic
	}{
		{Field{Type: "string"}, "size:255"},
		{Field{Type: "text"}, "type:text"},
		{Field{Type: "richtext"}, "type:text"},
		{Field{Type: "slug"}, "uniqueIndex"},
		{Field{Type: "belongs_to"}, "index"},
		{Field{Type: "string_array"}, "type:json"},
		{Field{Type: "bool"}, ""},
		{Field{Type: "string", Unique: true}, "uniqueIndex"},
		{Field{Type: "string", Default: "draft"}, "default:draft"},
	}

	for _, tt := range tests {
		got := tt.field.GORMTag()
		if tt.wantContain != "" && !strings.Contains(got, tt.wantContain) {
			t.Errorf("Field{Type:%q,Unique:%v,Default:%q}.GORMTag() = %q, does not contain %q",
				tt.field.Type, tt.field.Unique, tt.field.Default, got, tt.wantContain)
		}
	}
}

// ── Boolean helpers ───────────────────────────────────────────────────────────

func TestField_TypeHelpers(t *testing.T) {
	cases := []struct {
		typ            string
		isSlug         bool
		isBelongsTo    bool
		isManyToMany   bool
		isRelationship bool
		isStringArray  bool
		needsTimeImport bool
	}{
		{"slug", true, false, false, false, false, false},
		{"belongs_to", false, true, false, true, false, false},
		{"many_to_many", false, false, true, true, false, false},
		{"string_array", false, false, false, false, true, false},
		{"datetime", false, false, false, false, false, true},
		{"date", false, false, false, false, false, true},
		{"string", false, false, false, false, false, false},
	}

	for _, c := range cases {
		f := Field{Type: c.typ}
		if got := f.IsSlug(); got != c.isSlug {
			t.Errorf("Field{%q}.IsSlug() = %v, want %v", c.typ, got, c.isSlug)
		}
		if got := f.IsBelongsTo(); got != c.isBelongsTo {
			t.Errorf("Field{%q}.IsBelongsTo() = %v, want %v", c.typ, got, c.isBelongsTo)
		}
		if got := f.IsManyToMany(); got != c.isManyToMany {
			t.Errorf("Field{%q}.IsManyToMany() = %v, want %v", c.typ, got, c.isManyToMany)
		}
		if got := f.IsRelationship(); got != c.isRelationship {
			t.Errorf("Field{%q}.IsRelationship() = %v, want %v", c.typ, got, c.isRelationship)
		}
		if got := f.IsStringArray(); got != c.isStringArray {
			t.Errorf("Field{%q}.IsStringArray() = %v, want %v", c.typ, got, c.isStringArray)
		}
		if got := f.NeedsTimeImport(); got != c.needsTimeImport {
			t.Errorf("Field{%q}.NeedsTimeImport() = %v, want %v", c.typ, got, c.needsTimeImport)
		}
	}
}

// ── FKColumnName ──────────────────────────────────────────────────────────────

func TestField_FKColumnName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"category", "category_id"},
		{"author", "author_id"},
		{"user", "user_id"},
		{"blog_post", "blog_post_id"},
	}

	for _, tt := range tests {
		f := Field{Name: tt.name, Type: "belongs_to"}
		if got := f.FKColumnName(); got != tt.want {
			t.Errorf("Field{Name:%q}.FKColumnName() = %q, want %q", tt.name, got, tt.want)
		}
	}
}

// ── RelatedModelName ─────────────────────────────────────────────────────────

func TestField_RelatedModelName(t *testing.T) {
	tests := []struct {
		name         string
		relatedModel string
		want         string
	}{
		{"category", "", "Category"},
		{"author", "", "Author"},
		{"author", "User", "User"},
		{"blog_post", "", "BlogPost"},
	}

	for _, tt := range tests {
		f := Field{Name: tt.name, Type: "belongs_to", RelatedModel: tt.relatedModel}
		if got := f.RelatedModelName(); got != tt.want {
			t.Errorf("Field{Name:%q,RelatedModel:%q}.RelatedModelName() = %q, want %q",
				tt.name, tt.relatedModel, got, tt.want)
		}
	}
}

// ── ColumnFormat / FormFieldType / IsSortable / IsSearchable ─────────────────

func TestField_UIHelpers(t *testing.T) {
	tests := []struct {
		fieldType      string
		columnFormat   string
		formFieldType  string
		isSortable     bool
		isSearchable   bool
	}{
		{"string", "text", "text", true, true},
		{"text", "text", "textarea", false, true},
		{"richtext", "richtext", "richtext", false, true},
		{"int", "text", "number", true, false},
		{"bool", "boolean", "toggle", false, false},
		{"datetime", "relative", "datetime", true, false},
		{"date", "relative", "date", true, false},
		{"slug", "text", "", true, true},
		{"belongs_to", "text", "relationship-select", false, false},
		{"many_to_many", "text", "multi-relationship-select", false, false},
		{"string_array", "text", "images", false, false},
	}

	for _, tt := range tests {
		f := Field{Type: tt.fieldType}
		if got := f.ColumnFormat(); got != tt.columnFormat {
			t.Errorf("Field{%q}.ColumnFormat() = %q, want %q", tt.fieldType, got, tt.columnFormat)
		}
		if got := f.FormFieldType(); got != tt.formFieldType {
			t.Errorf("Field{%q}.FormFieldType() = %q, want %q", tt.fieldType, got, tt.formFieldType)
		}
		if got := f.IsSortable(); got != tt.isSortable {
			t.Errorf("Field{%q}.IsSortable() = %v, want %v", tt.fieldType, got, tt.isSortable)
		}
		if got := f.IsSearchable(); got != tt.isSearchable {
			t.Errorf("Field{%q}.IsSearchable() = %v, want %v", tt.fieldType, got, tt.isSearchable)
		}
	}
}
