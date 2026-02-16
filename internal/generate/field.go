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
	FieldBelongsTo  FieldType = "belongs_to"
	FieldManyToMany FieldType = "many_to_many"
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

// GoType returns the Go type for this field.
func (f Field) GoType() string {
	switch FieldType(f.Type) {
	case FieldString, FieldText, FieldSlug:
		return "string"
	case FieldInt:
		return "int"
	case FieldUint, FieldBelongsTo:
		return "uint"
	case FieldFloat:
		return "float64"
	case FieldBool:
		return "bool"
	case FieldDatetime, FieldDate:
		return "*time.Time"
	case FieldManyToMany:
		return "[]uint"
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

	switch FieldType(f.Type) {
	case FieldString:
		parts = append(parts, "size:255")
	case FieldText:
		parts = append(parts, "type:text")
	case FieldDate:
		parts = append(parts, "type:date")
	case FieldSlug:
		parts = append(parts, "size:255", "uniqueIndex")
	case FieldBelongsTo:
		parts = append(parts, "index")
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
	case FieldString, FieldText, FieldSlug:
		return "string"
	case FieldInt, FieldUint, FieldFloat, FieldBelongsTo:
		return "number"
	case FieldBool:
		return "boolean"
	case FieldDatetime, FieldDate:
		return "string | null"
	case FieldManyToMany:
		return "number[]"
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
	case FieldText:
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
		base = `z.number().int().min(1, "Required")`
	case FieldManyToMany:
		base = "z.array(z.number().int()).optional()"
	default:
		base = "z.string()"
	}

	if !f.Required && FieldType(f.Type) != FieldDatetime && FieldType(f.Type) != FieldDate && FieldType(f.Type) != FieldSlug && FieldType(f.Type) != FieldBelongsTo && FieldType(f.Type) != FieldManyToMany {
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
	return FieldType(f.Type) == FieldString || FieldType(f.Type) == FieldText || FieldType(f.Type) == FieldSlug
}

// ValidFieldTypes returns all valid field type names.
func ValidFieldTypes() []string {
	return []string{"string", "text", "int", "uint", "float", "bool", "datetime", "date", "slug", "belongs_to", "many_to_many"}
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
