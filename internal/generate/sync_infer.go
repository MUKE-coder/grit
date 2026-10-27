package generate

import (
	"strconv"
	"strings"
)

// Two writers produced the Zod schema for the same model, and the second one
// to run won.
//
// `grit generate resource` writes it from the definition, which knows every
// rule: money, the enum's values, the uuid, the min(1) on a required string,
// the refinement for a field that is required only while another field holds a
// value. `grit sync` then rewrote the same file from the Go struct, and emitted
// z.unknown() for money, z.string() for the enum, no uuid, no min, no
// refinement, and `version` as a required field on create although the API
// drops it. Sync is the documented step after generating, so the weaker file
// is the one projects ended up with: a create form built from it could not be
// submitted at all (grit#108).
//
// The fix is not a second set of rules in sync. It is to stop sync deciding.
// The Go model carries what it needs -- binding:"required", enum, requiredif,
// format, money.Money, the FileRef types -- so sync reconstructs the field the
// model was generated from and hands it to Field.ZodType(), the emitter the
// generator uses. One writer, so the next field type is right in both places
// the day it is added to one.
//
// What a tag cannot say, a tag now says: a slug column and the owner column of
// an --owned-by resource look exactly like ordinary columns, and both are
// filled by the server, so both carry a marker tag for this to read.

// inferField reconstructs the resource field a Go struct field was generated
// from. ok is false when the Go type says nothing useful, and the caller falls
// back to the plain type mapping.
func inferField(f GoField, all []GoField) (Field, bool) {
	out := Field{
		Name: f.Name,
		// Two tags for one rule, because gin's cannot be used on every type:
		// binding:"required" rejects 0 and false, so a required number carries
		// required:"true" instead, which is a declaration and not a check.
		Required: strings.Contains(f.Binding, "required") || f.Required,
	}

	// Required only while a sibling holds a value. The model never carries
	// binding:"required" for one of these, because that would refuse every row
	// the condition does not hold for, so the requiredif tag is the only thing
	// that says the rule exists.
	// A column with a default may be left out of a create: the database fills
	// it in. The generator knew this from the definition and sync did not, so
	// the two disagreed on every defaulted field (grit#114).
	//
	// "null" is not a default in this sense. It is how an optional foreign key
	// is written, and that case is decided by the relationship rather than by
	// the column.
	if d := gormOption(f.GORMTag, "default"); d != "" && !strings.EqualFold(d, "null") {
		out.Default = d
	}

	if on, value, found := strings.Cut(f.RequiredIf, "="); found {
		out.Required = true
		out.ShowWhen = &FieldCondition{Field: on, Value: value}
	}

	// A formatted column names its own type: format:"email", format:"rating:5".
	if kind, param, _ := strings.Cut(f.Format, ":"); kind != "" {
		out.Type = kind
		if FieldType(kind) == FieldRating {
			if n, err := strconv.Atoi(param); err == nil && n > 0 {
				out.Max = n
			}
		}
		return out, true
	}

	switch f.GoType {
	case "money.Money", "*money.Money":
		out.Type = string(FieldMoney)
		return out, true
	case "*files.FileRef", "files.FileRef":
		out.Type = string(FieldFile)
		return out, true
	case "files.FileRefs", "*files.FileRefs":
		out.Type = string(FieldFiles)
		return out, true
	case "jsontime.DateTime", "*jsontime.DateTime", "time.Time", "*time.Time":
		out.Type = string(FieldDatetime)
		return out, true
	case "jsontime.Date", "*jsontime.Date", "datatypes.Date", "*datatypes.Date":
		out.Type = string(FieldDate)
		return out, true
	case "bool", "*bool":
		out.Type = string(FieldBool)
		return out, true
	case "int", "int8", "int16", "int32", "int64", "*int", "*int64":
		out.Type = string(FieldInt)
		return out, true
	case "uint", "uint8", "uint16", "uint32", "uint64", "*uint":
		out.Type = string(FieldUint)
		return out, true
	case "float32", "float64", "*float64":
		out.Type = string(FieldFloat)
		return out, true
	}

	// A list of strings is a multi-choice when the column names its values and
	// a free list of tags when it does not.
	if strings.Contains(f.GoType, "JSONSlice[string]") || f.GoType == "[]string" {
		if f.Enum != "" {
			out.Type = string(FieldCheck)
			out.Options = enumOptions(f.Enum)
		} else {
			out.Type = string(FieldStringArray)
		}
		return out, true
	}

	if strings.HasPrefix(f.GoType, "string") || f.GoType == "*string" {
		switch {
		case f.Enum != "":
			out.Type = string(FieldSelect)
			out.Options = enumOptions(f.Enum)
		case isForeignKey(f, all):
			out.Type = string(FieldBelongsTo)
			// A pointer is the column saying it may be NULL, which is how an
			// optional relationship is written.
			if strings.HasPrefix(f.GoType, "*") {
				out.Required = false
			}
		case f.Sanitize == "html":
			out.Type = string(FieldRichtext)
		case strings.Contains(f.GORMTag, "type:text"):
			out.Type = string(FieldText)
		default:
			out.Type = string(FieldString)
		}
		return out, true
	}

	return out, false
}

// isForeignKey reports whether this column holds another record's id, by
// finding the association that names it.
//
// By the association rather than by the column name ending in "_id", because a
// hand-written model may have an external_id or a stripe_customer_id that is
// not a relationship at all, and declaring those uuid would reject every value
// they legitimately hold.
func isForeignKey(f GoField, all []GoField) bool {
	for _, other := range all {
		if other.Name == f.Name {
			continue
		}
		if fk := gormOption(other.GORMTag, "foreignKey"); fk == f.Name {
			return true
		}
	}
	return false
}

// selfReferencing reports whether this foreign key points back at its own
// model, which is how a tree's parent is written.
func selfReferencing(f GoField, s GoStruct) bool {
	for _, other := range s.Fields {
		if gormOption(other.GORMTag, "foreignKey") != f.Name {
			continue
		}
		if strings.TrimPrefix(other.GoType, "*") == s.Name {
			return true
		}
	}
	return false
}

// gormOption reads one setting out of a gorm tag: gormOption("size:36;index",
// "size") is "36".
func gormOption(tag, key string) string {
	for _, part := range strings.Split(tag, ";") {
		name, value, found := strings.Cut(strings.TrimSpace(part), ":")
		if found && strings.EqualFold(name, key) {
			return value
		}
	}
	return ""
}

// enumOptions turns the model's enum:"open|won|lost" into the choices the
// emitter expects. The label is the value: the model does not carry labels,
// and nothing downstream of the schema reads them.
func enumOptions(enum string) []FieldOption {
	var out []FieldOption
	for _, v := range strings.Split(enum, "|") {
		if v = strings.TrimSpace(v); v != "" {
			out = append(out, FieldOption{Value: v, Label: v})
		}
	}
	return out
}
