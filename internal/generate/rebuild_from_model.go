package generate

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Rebuilding a resource definition from the model Grit already wrote, for
// `grit generate resource <Name> --force` with no --fields.
//
// grit#87: the mobile and desktop screens are written once, at generation, and
// no command updates them. The only route to a regenerated screen is
//
//	grit generate resource <Name> --fields "<every field, including the new one>" --force
//
// which means retyping the complete field list from the model. That is tedious
// and, more to the point, easy to get subtly wrong: a type remembered as
// "string" that was "text", a missing unique, a belongs_to spelled without its
// target. The retyped list silently becomes the new truth, and the first
// symptom is a column that changed type under a table full of data.
//
// So --force with no --fields reads the model. The list a rebuild uses is then
// the one the database actually has, which is the only list that can be right.
//
// # Why this is not a second parser
//
// definitionFromModelFile already reads a model into a definition, through the
// shared parseGoStructs, for `grit generate seeder`. Writing another reader
// beside it would be two parsers over one file format, disagreeing quietly
// from the first time somebody fixed one of them.
//
// They differ in what they skip, not in how they read, so the difference is a
// parameter. A seeder must not invent a slug or a foreign key: the slug is
// computed and the FK would point at a row that may not exist. A rebuild must
// keep both, because dropping them from the definition drops the columns from
// the regenerated model.

// modelReadMode says what a reader of a generated model is for.
type modelReadMode int

const (
	// modelForSeeder omits columns a seeder must not invent: the slug, which
	// is computed, and foreign keys, which would point at nothing.
	modelForSeeder modelReadMode = iota
	// modelForRebuild keeps every column somebody declared, because the
	// definition it produces is about to rewrite the model itself.
	modelForRebuild
)

// RebuildDefinitionFromModel reads a generated model and returns the
// definition that would produce it, plus the --fields spelling it corresponds
// to.
//
// The spelling is returned so the caller can print it. That is not decoration:
// --force is destructive, and somebody has to be able to read the list before
// it is used, and to paste it back with an edit if it is wrong.
func RebuildDefinitionFromModel(name string) (*ResourceDefinition, string, error) {
	root, err := findProjectRoot()
	if err != nil {
		return nil, "", err
	}
	names := MakeNames(name)
	for _, apiRoot := range []string{filepath.Join(root, "apps", "api"), root} {
		if modelFileFor(apiRoot, names) == "" {
			continue
		}
		def, err := definitionFromModelFileMode(apiRoot, name, modelForRebuild)
		if err != nil {
			return nil, "", err
		}
		if len(def.Fields) == 0 {
			return nil, "", fmt.Errorf(
				"no declarable fields found on %s, so there is nothing to rebuild from: pass --fields",
				names.Pascal)
		}
		return def, fieldsSpelling(def), nil
	}
	return nil, "", fmt.Errorf(
		"no model found for %q: generate the resource first (looked in apps/api/internal/models and internal/models)",
		names.Pascal)
}

// fieldsSpelling renders a definition back into the --fields string a person
// would have typed. Shown to the user, never parsed back: the definition is
// already in hand, and a round trip through text is a second thing to keep in
// step with this one.
func fieldsSpelling(def *ResourceDefinition) string {
	parts := make([]string, 0, len(def.Fields))
	for _, f := range def.Fields {
		spec := f.Name + ":" + f.Type
		switch {
		case f.IsBelongsTo() && f.RelatedModel != "":
			spec = f.Name + ":belongs_to[" + f.RelatedModel + "]"
		case len(f.Options) > 0:
			values := make([]string, 0, len(f.Options))
			for _, o := range f.Options {
				values = append(values, o.Value)
			}
			spec = f.Name + ":" + f.Type + "[" + strings.Join(values, ", ") + "]"
		}
		if f.Unique {
			spec += ":unique"
		}
		parts = append(parts, spec)
	}
	return strings.Join(parts, ",")
}

// belongsToTargetFor names the model a <target>_id column points at.
//
// Read from the GORM tag when it says so, and from the column name otherwise.
// The tag is preferred because a column may be named differently from its
// target, and a guess from the name would then generate a relation to a model
// that does not exist.
func belongsToTargetFor(gf GoField, owner *GoStruct) string {
	if ref := foreignKeyTargetFromTag(gf.GORMTag); ref != "" {
		return ref
	}
	base := strings.TrimSuffix(gf.JSONName, "_id")
	if base == "" {
		return ""
	}
	pascal := toPascalCase(base)
	// The association field beside it, which is the model's own statement of
	// what the FK points at: `User *models.User` or `User User`.
	for _, other := range owner.Fields {
		if other.Name != pascal {
			continue
		}
		t := strings.TrimPrefix(strings.TrimPrefix(other.GoType, "*"), "models.")
		if t != "" && t[0] >= 'A' && t[0] <= 'Z' {
			return t
		}
	}
	return pascal
}

// foreignKeyTargetFromTag reads a references clause out of a GORM tag.
func foreignKeyTargetFromTag(tag string) string {
	for _, part := range strings.Split(tag, ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "references:") {
			// references:ID names a column, not a model, so it says nothing
			// about the target and the column name is the better source.
			return ""
		}
	}
	return ""
}

// frameworkGoTypes are the Go types the generator emits that are not scalars,
// each with the --fields type that produces it.
//
// goTypeToFieldType was written for seeders and answers "" for all of these,
// which for a seeder means "skip it, the user wires it up". For a rebuild
// there is no such thing as skipping: the definition is about to rewrite the
// model, and a column missing from it is a column dropped from the table.
//
// money.Money is the one that found this. It is embedded, two columns behind
// one field, and nothing about the Go type says money to a reader that only
// knows scalars.
var frameworkGoTypes = map[string]FieldType{
	"money.Money":     FieldMoney,
	"files.FileRef":   FieldFile,
	"*files.FileRef":  FieldFile,
	"files.FileRefs":  FieldFiles,
	"vectors.Vector":  FieldVector,
	"pq.StringArray":  FieldStringArray,
	"datatypes.JSON":  FieldJSON,
	"json.RawMessage": FieldJSON,
}

// rebuildFieldFor turns one parsed column into the field a rebuild declares,
// or says why it cannot.
//
// It reuses the two tag readers the seeder already has, because a column's
// format and enum tags are the model's own statement of what it is and there
// should not be a second opinion about them.
//
// The error matters more than the mapping. This feeds a destructive command,
// so a column this cannot read has to stop the rebuild and name itself. The
// alternative is a regenerated model with a column quietly missing, and the
// first symptom of that is the migration that drops it.
func rebuildFieldFor(gf GoField) (Field, error) {
	if f, ok := fieldFromFormatTag(gf.JSONName, gf.Format); ok {
		f.Unique = gormTagSaysUnique(gf.GORMTag)
		f.Required = columnIsRequired(gf)
		return f, nil
	}
	if f, ok := fieldFromEnumTag(gf.JSONName, gf.Enum); ok {
		f.Required = columnIsRequired(gf)
		return f, nil
	}
	// A slug says, in its own tag, which column it is derived from. Without
	// that the rebuild would emit a plain string and the column would stop
	// being a slug.
	if gf.Slug != "" {
		return Field{Name: gf.JSONName, Type: string(FieldSlug), SlugSource: gf.Slug}, nil
	}
	if ft, ok := frameworkGoTypes[gf.GoType]; ok {
		return Field{
			Name:     gf.JSONName,
			Type:     string(ft),
			Required: columnIsRequired(gf),
		}, nil
	}

	ft := goTypeToFieldType(gf.GoType)
	if ft == "" {
		return Field{}, fmt.Errorf(
			"cannot tell what kind of field %q is from its Go type %s, so a rebuild would drop the column.\n"+
				"Pass --fields with the list you want, including this one",
			gf.JSONName, gf.GoType)
	}
	// text and string are the same Go type and differ only in the column.
	// Rebuilding a text column as a string narrows it, which on Postgres is a
	// migration that can fail on rows that already exist.
	if ft == string(FieldString) && strings.Contains(gf.GORMTag, "type:text") {
		ft = string(FieldText)
	}
	return Field{
		Name:     gf.JSONName,
		Type:     ft,
		Unique:   gormTagSaysUnique(gf.GORMTag),
		Required: columnIsRequired(gf),
	}, nil
}

// gormTagSaysUnique reports a column declared unique, either way of spelling it.
func gormTagSaysUnique(tag string) bool {
	for _, part := range strings.Split(tag, ";") {
		switch strings.TrimSpace(part) {
		case "unique", "uniqueIndex":
			return true
		}
		if strings.HasPrefix(strings.TrimSpace(part), "uniqueIndex:") {
			return true
		}
	}
	return false
}

// isAssociationField reports a struct field GORM fills by Preload rather than
// a column of this table: `Client *Client`, `Items []Item`.
//
// Its column is the <name>_id beside it, which is read separately. Telling
// them apart matters for a rebuild in a way it never did for a seeder: a
// seeder quietly skips anything it cannot type, and a rebuild has to refuse,
// so an association reaching the refusal would make every resource with a
// relationship unrebuildable.
func isAssociationField(gf GoField) bool {
	if strings.Contains(gf.GORMTag, "foreignKey:") || strings.Contains(gf.GORMTag, "many2many:") {
		return true
	}
	t := strings.TrimPrefix(strings.TrimPrefix(gf.GoType, "[]"), "*")
	t = strings.TrimPrefix(t, "models.")
	if t == "" || t == gf.GoType && !strings.HasPrefix(gf.GoType, "[]") && !strings.HasPrefix(gf.GoType, "*") {
		return false
	}
	// A bare Pascal-case name with no package: another model in this package.
	// time.Time and gorm.DeletedAt carry a dot and are not caught here.
	return !strings.Contains(t, ".") && t[0] >= 'A' && t[0] <= 'Z'
}

// columnIsRequired reports a column the model declares as required.
//
// Read from the binding tag, not from GoField.Required. That field is filled
// from a tag literally named required:"true", which a generated model does not
// have: requiredness lives in binding:"required", which parseGoStructs puts in
// Binding. So GoField.Required is false on every generated column, and a
// rebuild trusting it dropped `required` from three fields and turned a
// required foreign key into *string with default:null. Found by diffing the
// rebuilt model against the one it had just read.
func columnIsRequired(gf GoField) bool {
	if gf.Required {
		return true
	}
	for _, part := range strings.Split(gf.Binding, ",") {
		if strings.TrimSpace(part) == "required" {
			return true
		}
	}
	return false
}
