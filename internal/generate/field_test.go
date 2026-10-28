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
		{"belongs_to", "string"}, // FK matches UUID string PK
		{"float", "float64"},
		{"bool", "bool"},
		// Not *time.Time: that unmarshals RFC3339 only, and the admin picker
		// sends "2001-08-06". See internal/jsontime.
		{"datetime", "*jsontime.DateTime"},
		{"date", "*jsontime.Date"},
		{"many_to_many", "[]string"},
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
		{"belongs_to", "string"}, // FK matches UUID string PK
		{"bool", "boolean"},
		{"datetime", "string | null"},
		{"date", "string | null"},
		{"many_to_many", "string[]"},
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
		{Field{Type: "belongs_to"}, "z.string().uuid"},
		{Field{Type: "many_to_many"}, "z.array(z.string().uuid())"},
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
		typ             string
		isSlug          bool
		isBelongsTo     bool
		isManyToMany    bool
		isRelationship  bool
		isStringArray   bool
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
		fieldType     string
		columnFormat  string
		formFieldType string
		isSortable    bool
		isSearchable  bool
	}{
		{"string", "text", "text", true, true},
		{"text", "text", "textarea", false, true},
		{"richtext", "richtext", "richtext", false, true},
		{"int", "text", "number", true, false},
		{"bool", "boolean", "toggle", false, false},
		// A bare "datetime" with no name is not an _at timestamp, so it
		// formats as a date. The named cases are covered by
		// TestColumnFormatTellsDatesFromTimestamps.
		{"datetime", "date", "datetime", true, false},
		{"date", "date", "date", true, false},
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

// TestColumnFormatTellsDatesFromTimestamps is grit#119.
//
// A calendar date rendered relative to the current moment reads wrong for
// every value: today's date showed "just now" and yesterday's "18h ago", on a
// column whose own form field said "date".
func TestColumnFormatTellsDatesFromTimestamps(t *testing.T) {
	cases := []struct {
		fieldType string
		name      string
		want      string
		why       string
	}{
		{"date", "due", "date", "a calendar date has no time of day to be relative to"},
		{"date", "closed_at", "date", "still a calendar date, whatever it is called"},
		{"datetime", "created_at", "relative", "a record of when something happened"},
		{"datetime", "CreatedAt", "relative", "the camelCase spelling of the same thing"},
		{"datetime", "published_at", "relative", ""},
		{"datetime", "starts", "date", "a moment somebody picked, not a timestamp"},
		{"datetime", "format", "date", "ends in 'at' only if you delete the underscore"},
		{"datetime", "seat", "date", ""},
	}

	for _, tc := range cases {
		f := Field{Type: tc.fieldType, Name: tc.name}
		if got := f.ColumnFormat(); got != tc.want {
			t.Errorf("Field{Type:%q, Name:%q}.ColumnFormat() = %q, want %q%s",
				tc.fieldType, tc.name, got, tc.want,
				map[bool]string{true: " (" + tc.why + ")", false: ""}[tc.why != ""])
		}
	}
}

// TestColumnFormatAgreesWithGuessFormat: the generator and `grit sync` pick
// the format for the same column by two different code paths, and they used
// to disagree, so a generated resource and a synced one formatted the same
// field differently.
func TestColumnFormatAgreesWithGuessFormat(t *testing.T) {
	cases := []struct{ name, goType, fieldType string }{
		{"created_at", "time.Time", "datetime"},
		{"updated_at", "time.Time", "datetime"},
		{"due", "time.Time", "date"},
		{"starts", "time.Time", "datetime"},
	}

	for _, tc := range cases {
		generated := Field{Type: tc.fieldType, Name: tc.name}.ColumnFormat()
		synced := guessFormat(GoField{Name: tc.name, GoType: tc.goType})
		if generated != synced {
			t.Errorf("%s: generate says %q, sync says %q; they must agree",
				tc.name, generated, synced)
		}
	}
}
