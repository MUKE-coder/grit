package generate

import (
	"fmt"
	"strings"
)

// FieldType represents a supported Grit field type.
type FieldType string

const (
	FieldString     FieldType = "string"
	FieldText       FieldType = "text"
	FieldInt        FieldType = "int"
	FieldUint       FieldType = "uint"
	FieldFloat      FieldType = "float"
	FieldBool       FieldType = "bool"
	FieldDatetime   FieldType = "datetime"
	FieldDate       FieldType = "date"
	FieldSlug       FieldType = "slug"
	FieldRichtext   FieldType = "richtext"
	FieldBelongsTo   FieldType = "belongs_to"
	FieldManyToMany  FieldType = "many_to_many"
	FieldStringArray FieldType = "string_array"
	FieldFile        FieldType = "file"  // single FileRef
	FieldFiles       FieldType = "files" // []FileRef
)

// Field describes a single field in a resource.
type Field struct {
	Name         string `yaml:"name"`
	Type         string `yaml:"type"`
	Required     bool   `yaml:"required"`
	Unique       bool   `yaml:"unique"`
	Default      string `yaml:"default"`
	SlugSource   string `yaml:"slug_source"`
	RelatedModel string `yaml:"related_model"`

	// FileAccepts is the resolved list of accept-aliases for a file/files
	// field. Source: the third position of name:file:<accept-list>. May be
	// a single alias ("image") or a bracketed list ("[pdf,doc,image]").
	// Values are aliases ("image", "video", "pdf", "all", ...) — see
	// resolveFileMIMEs() in the scaffolded API for the runtime mapping.
	FileAccepts []string `yaml:"file_accepts,omitempty"`
}

// IsSlug returns true if this field is an auto-generated slug.
func (f Field) IsSlug() bool {
	return FieldType(f.Type) == FieldSlug
}

// IsBelongsTo returns true if this field is a belongs_to relationship.
func (f Field) IsBelongsTo() bool {
	return FieldType(f.Type) == FieldBelongsTo
}

// IsManyToMany returns true if this field is a many_to_many relationship.
func (f Field) IsManyToMany() bool {
	return FieldType(f.Type) == FieldManyToMany
}

// IsRelationship returns true if this field is any relationship type.
func (f Field) IsRelationship() bool {
	return f.IsBelongsTo() || f.IsManyToMany()
}

// IsStringArray returns true if this field is a string array (JSON).
func (f Field) IsStringArray() bool {
	return FieldType(f.Type) == FieldStringArray
}

// IsFile returns true if this field is a single FileRef.
func (f Field) IsFile() bool {
	return FieldType(f.Type) == FieldFile
}

// IsFiles returns true if this field is a []FileRef.
func (f Field) IsFiles() bool {
	return FieldType(f.Type) == FieldFiles
}

// IsFileField returns true for either file or files.
func (f Field) IsFileField() bool {
	return f.IsFile() || f.IsFiles()
}

// NeedsDatatypesImport returns true if this field requires "gorm.io/datatypes" import.
func (f Field) NeedsDatatypesImport() bool {
	return FieldType(f.Type) == FieldStringArray
}

// NeedsFilesImport returns true if this field requires the models/files
// package import (for the FileRef type) in the generated Go model.
func (f Field) NeedsFilesImport() bool {
	return f.IsFileField()
}

// GoType returns the Go type for this field.
func (f Field) GoType() string {
	switch FieldType(f.Type) {
	case FieldString, FieldText, FieldSlug, FieldRichtext:
		return "string"
	case FieldInt:
		return "int"
	case FieldUint:
		return "uint"
	case FieldBelongsTo:
		// FK columns match the referenced model's UUID string PK.
		return "string"
	case FieldFloat:
		return "float64"
	case FieldBool:
		return "bool"
	case FieldDatetime, FieldDate:
		return "*time.Time"
	case FieldManyToMany:
		return "[]string"
	case FieldStringArray:
		return "datatypes.JSONSlice[string]"
	case FieldFile:
		return "*files.FileRef"
	case FieldFiles:
		return "files.FileRefs"
	default:
		return "string"
	}
}

// GORMTag returns the GORM struct tag for this field.
func (f Field) GORMTag() string {
	// Relationship fields handle their own GORM tags in the template
	if f.IsManyToMany() {
		return ""
	}

	parts := []string{}

	name := strings.ToLower(f.Name)

	switch FieldType(f.Type) {
	case FieldString:
		// Heuristic: URL-shaped fields blow past 255 in the wild
		// (UTM-tagged tracking links, signed S3 URLs, etc.). Bump to 500.
		// Long-form "description"-style fields really want type:text so
		// PG doesn't truncate. Default everything else to size:255.
		switch {
		case isURLField(name):
			parts = append(parts, "size:500")
		case isLongTextField(name):
			parts = append(parts, "type:text")
		default:
			parts = append(parts, "size:255")
		}
	case FieldText, FieldRichtext:
		parts = append(parts, "type:text")
	case FieldDate:
		parts = append(parts, "type:date")
	case FieldSlug:
		parts = append(parts, "size:255", "uniqueIndex")
	case FieldBelongsTo:
		// FK matches the referenced model's UUID string PK.
		parts = append(parts, "size:36", "index")
	case FieldStringArray:
		parts = append(parts, "type:json")
	case FieldFile, FieldFiles:
		// FileRef / FileRefs implement Value / Scan via the files package,
		// so GORM stores them as JSON. type:json signals jsonb on Postgres
		// (otherwise we'd get text and lose efficient querying).
		parts = append(parts, "type:json")
	case FieldFloat:
		// Heuristic: money-shaped fields need fixed-precision storage
		// to avoid float rounding (1.99 + 0.01 = 1.9999999...).
		// Use decimal(12,2) — 10 digits before the decimal, 2 after —
		// which is plenty for any individual transaction.
		if isMoneyField(name) {
			parts = append(parts, "type:decimal(12,2)")
		}
	}

	if f.Unique && FieldType(f.Type) != FieldSlug {
		parts = append(parts, "uniqueIndex")
	}
	if !f.Required {
		switch FieldType(f.Type) {
		case FieldString, FieldText:
			// strings default to "" which is fine
		default:
			// no extra tag needed
		}
	}
	if f.Default != "" {
		parts = append(parts, fmt.Sprintf("default:%s", f.Default))
	}

	if len(parts) == 0 {
		return ""
	}

	tag := ""
	for i, p := range parts {
		if i > 0 {
			tag += ";"
		}
		tag += p
	}
	return tag
}

// TSType returns the TypeScript type for this field.
func (f Field) TSType() string {
	switch FieldType(f.Type) {
	case FieldString, FieldText, FieldSlug, FieldRichtext:
		return "string"
	case FieldInt, FieldUint, FieldFloat:
		return "number"
	case FieldBelongsTo:
		// FK columns match the referenced model's UUID string PK.
		return "string"
	case FieldBool:
		return "boolean"
	case FieldDatetime, FieldDate:
		return "string | null"
	case FieldManyToMany:
		return "string[]"
	case FieldStringArray:
		return "string[]"
	case FieldFile:
		return "FileRef | null"
	case FieldFiles:
		return "FileRef[]"
	default:
		return "string"
	}
}

// ZodType returns the Zod validator for this field.
func (f Field) ZodType() string {
	base := ""
	switch FieldType(f.Type) {
	case FieldString:
		base = "z.string()"
		if f.Required {
			base += `.min(1, "Required")`
		}
	case FieldText, FieldRichtext:
		base = "z.string()"
	case FieldSlug:
		base = "z.string()"
	case FieldInt:
		base = "z.number().int()"
	case FieldUint:
		base = "z.number().int().nonnegative()"
	case FieldFloat:
		base = "z.number()"
	case FieldBool:
		base = "z.boolean()"
	case FieldDatetime, FieldDate:
		base = "z.string().nullable()"
	case FieldBelongsTo:
		// FK columns are UUID strings matching the referenced model's PK.
		base = `z.string().uuid("Invalid ID")`
	case FieldManyToMany:
		base = "z.array(z.string().uuid()).optional()"
	case FieldStringArray:
		base = "z.array(z.string()).optional()"
	case FieldFile:
		base = "FileRefSchema.nullable()"
		if f.Required {
			base = "FileRefSchema"
		}
	case FieldFiles:
		base = "z.array(FileRefSchema).default([])"
	default:
		base = "z.string()"
	}

	if !f.Required && FieldType(f.Type) != FieldDatetime && FieldType(f.Type) != FieldDate && FieldType(f.Type) != FieldSlug && FieldType(f.Type) != FieldRichtext && FieldType(f.Type) != FieldBelongsTo && FieldType(f.Type) != FieldManyToMany && FieldType(f.Type) != FieldStringArray && FieldType(f.Type) != FieldFile && FieldType(f.Type) != FieldFiles {
		base += ".optional()"
	}

	return base
}

// NeedsTimeImport returns true if this field requires "time" import in Go.
func (f Field) NeedsTimeImport() bool {
	return FieldType(f.Type) == FieldDatetime || FieldType(f.Type) == FieldDate
}

// ColumnFormat returns the DataTable column format for this field type.
func (f Field) ColumnFormat() string {
	switch FieldType(f.Type) {
	case FieldBool:
		return "boolean"
	case FieldDatetime, FieldDate:
		return "relative"
	case FieldRichtext:
		return "richtext"
	case FieldFile:
		return "file"
	case FieldFiles:
		return "files"
	default:
		return "text"
	}
}

// FormFieldType returns the form builder field type for this field type.
// Returns "" for auto-generated fields like slug (excluded from forms).
func (f Field) FormFieldType() string {
	switch FieldType(f.Type) {
	case FieldString:
		return "text"
	case FieldText:
		return "textarea"
	case FieldRichtext:
		return "richtext"
	case FieldInt, FieldUint, FieldFloat:
		return "number"
	case FieldBool:
		return "toggle"
	case FieldDatetime:
		return "datetime"
	case FieldDate:
		return "date"
	case FieldSlug:
		return ""
	case FieldBelongsTo:
		return "relationship-select"
	case FieldManyToMany:
		return "multi-relationship-select"
	case FieldStringArray:
		return "images"
	case FieldFile:
		return "file"
	case FieldFiles:
		return "files"
	default:
		return "text"
	}
}

// IsSortable returns true if this field type should be sortable by default.
func (f Field) IsSortable() bool {
	switch FieldType(f.Type) {
	case FieldString, FieldInt, FieldUint, FieldFloat, FieldDatetime, FieldDate, FieldSlug:
		return true
	default:
		return false
	}
}

// IsSearchable returns true if this field type should be searchable by default.
func (f Field) IsSearchable() bool {
	return FieldType(f.Type) == FieldString || FieldType(f.Type) == FieldText || FieldType(f.Type) == FieldSlug || FieldType(f.Type) == FieldRichtext
}

// ValidFieldTypes returns all valid field type names.
func ValidFieldTypes() []string {
	return []string{"string", "text", "richtext", "int", "uint", "float", "bool", "datetime", "date", "slug", "belongs_to", "many_to_many", "string_array", "file", "files"}
}

// FKColumnName returns the foreign key column name for a belongs_to field.
// e.g., "category" → "category_id", "author" → "author_id"
func (f Field) FKColumnName() string {
	name := toSnakeCase(toPascalCase(f.Name))
	if !strings.HasSuffix(name, "_id") {
		name += "_id"
	}
	return name
}

// RelatedModelName returns the PascalCase related model name.
// Uses the explicit RelatedModel if set, otherwise infers from field name.
func (f Field) RelatedModelName() string {
	if f.RelatedModel != "" {
		return toPascalCase(f.RelatedModel)
	}
	// Infer from field name: "category" → "Category", "author" → "Author"
	name := f.Name
	// Strip _id suffix if present
	name = strings.TrimSuffix(name, "_id")
	name = strings.TrimSuffix(name, "Id")
	return toPascalCase(name)
}

// isURLField returns true for field names that are very likely to hold
// a URL — those blow past size:255 in the wild (UTM-tagged links,
// signed S3 URLs, profile picture URLs from external IDPs).
func isURLField(name string) bool {
	if strings.HasSuffix(name, "_url") {
		return true
	}
	switch name {
	case "url", "image", "avatar", "thumbnail", "logo", "cover", "icon", "banner", "photo":
		return true
	}
	return false
}

// isLongTextField returns true for field names that conventionally
// hold long-form text — these really want PG TEXT instead of VARCHAR.
func isLongTextField(name string) bool {
	switch name {
	case "description", "notes", "content", "body", "summary", "bio", "details", "comment", "comments", "message":
		return true
	}
	return false
}

// isMoneyField returns true for float field names that conventionally
// hold money — those need fixed-precision storage to avoid float
// rounding artifacts (decimal(12,2) gives 10 whole digits + 2 cents).
func isMoneyField(name string) bool {
	suffixes := []string{"_amount", "_price", "_total", "_cost", "_fee", "_balance", "_rent", "_salary", "_wage", "_value", "_revenue", "_deposit"}
	for _, suf := range suffixes {
		if strings.HasSuffix(name, suf) {
			return true
		}
	}
	switch name {
	case "amount", "price", "total", "cost", "fee", "balance", "subtotal":
		return true
	}
	return false
}
