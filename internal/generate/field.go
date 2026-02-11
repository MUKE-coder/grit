package generate

import "fmt"

// FieldType represents a supported Grit field type.
type FieldType string

const (
	FieldString   FieldType = "string"
	FieldText     FieldType = "text"
	FieldInt      FieldType = "int"
	FieldUint     FieldType = "uint"
	FieldFloat    FieldType = "float"
	FieldBool     FieldType = "bool"
	FieldDatetime FieldType = "datetime"
	FieldDate     FieldType = "date"
)

// Field describes a single field in a resource.
type Field struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
	Unique   bool   `yaml:"unique"`
	Default  string `yaml:"default"`
}

// GoType returns the Go type for this field.
func (f Field) GoType() string {
	switch FieldType(f.Type) {
	case FieldString, FieldText:
		return "string"
	case FieldInt:
		return "int"
	case FieldUint:
		return "uint"
	case FieldFloat:
		return "float64"
	case FieldBool:
		return "bool"
	case FieldDatetime, FieldDate:
		return "*time.Time"
	default:
		return "string"
	}
}

// GORMTag returns the GORM struct tag for this field.
func (f Field) GORMTag() string {
	parts := []string{}

	switch FieldType(f.Type) {
	case FieldString:
		parts = append(parts, "size:255")
	case FieldText:
		parts = append(parts, "type:text")
	case FieldDate:
		parts = append(parts, "type:date")
	}

	if f.Unique {
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
	case FieldString, FieldText:
		return "string"
	case FieldInt, FieldUint, FieldFloat:
		return "number"
	case FieldBool:
		return "boolean"
	case FieldDatetime, FieldDate:
		return "string | null"
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
	default:
		base = "z.string()"
	}

	if !f.Required && FieldType(f.Type) != FieldDatetime && FieldType(f.Type) != FieldDate {
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
	case FieldFloat:
		return "text"
	default:
		return "text"
	}
}

// FormFieldType returns the form builder field type for this field type.
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
	default:
		return "text"
	}
}

// IsSortable returns true if this field type should be sortable by default.
func (f Field) IsSortable() bool {
	switch FieldType(f.Type) {
	case FieldString, FieldInt, FieldUint, FieldFloat, FieldDatetime, FieldDate:
		return true
	default:
		return false
	}
}

// IsSearchable returns true if this field type should be searchable by default.
func (f Field) IsSearchable() bool {
	return FieldType(f.Type) == FieldString || FieldType(f.Type) == FieldText
}

// ValidFieldTypes returns all valid field type names.
func ValidFieldTypes() []string {
	return []string{"string", "text", "int", "uint", "float", "bool", "datetime", "date"}
}
