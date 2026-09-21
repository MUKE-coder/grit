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
	FieldBelongsTo  FieldType = "belongs_to"
	FieldManyToMany FieldType = "many_to_many"
	// One-to-one: a belongs_to whose foreign key is unique.
	//
	// Declared on the side that holds the key, so a Profile declares its
	// User rather than the other way round. That is the side the constraint
	// can actually live on, and a unique index is the entire difference
	// between "many profiles per user" and "one".
	FieldOneToOne FieldType = "one_to_one"
	// Money: an integer count of minor units plus its currency.
	//
	// Not float. 0.1 + 0.2 is not 0.3 in binary floating point, and after
	// tax, a discount and a refund the ledger drifts. float stays available
	// and is still right for weights, ratings and percentages.
	FieldMoney       FieldType = "money"
	FieldStringArray FieldType = "string_array"
	FieldFile        FieldType = "file"  // single FileRef
	FieldFiles       FieldType = "files" // []FileRef

	// Option-backed types. select = one value from a fixed list (dropdown);
	// check = zero-or-more values (checkbox group, stored as a JSON array);
	// toggle = an on/off boolean (a friendlier alias for bool). Options are
	// declared as value=Label pairs: "status:select:draft=Draft|paid=Paid".
	FieldSelect FieldType = "select"
	FieldRadio  FieldType = "radio"
	FieldCheck  FieldType = "check"
	FieldToggle FieldType = "toggle"

	// Formatted types: a value with a shape the API checks and normalises on
	// every write (internal/fieldtypes), a matching Zod schema, and an admin
	// input built for it. Declared as name:type, with an option in the third
	// position for the two that take one: phone:tel:UG sets the default
	// country, score:rating:10 the number of stars.
	FieldEmail   FieldType = "email"
	FieldURL     FieldType = "url"
	FieldDomain  FieldType = "domain"
	FieldTel     FieldType = "tel"
	FieldCountry FieldType = "country"
	FieldColor   FieldType = "color"
	FieldPercent FieldType = "percent"
	FieldRating  FieldType = "rating"
	FieldTime    FieldType = "time"
	FieldJSON    FieldType = "json"
)

// DefaultRatingMax is the number of stars a rating field has when the
// declaration does not say.
const DefaultRatingMax = 5

// Field describes a single field in a resource.
type Field struct {
	Name         string `yaml:"name"`
	Type         string `yaml:"type"`
	Required     bool   `yaml:"required"`
	Unique       bool   `yaml:"unique"`
	Default      string `yaml:"default"`
	SlugSource   string `yaml:"slug_source"`
	RelatedModel string `yaml:"related_model"`

	// Encrypted stores the column as crypto.EncryptedString: AES-256-GCM at
	// rest, plaintext in code and over the wire. Set by the :encrypted
	// modifier, on string, text and richtext fields only.
	Encrypted bool `yaml:"encrypted,omitempty"`

	// FileAccepts is the resolved list of accept-aliases for a file/files
	// field. Source: the third position of name:file:<accept-list>. May be
	// a single alias ("image") or a bracketed list ("[pdf,doc,image]").
	// Values are aliases ("image", "video", "pdf", "all", ...) — see
	// resolveFileMIMEs() in the scaffolded API for the runtime mapping.
	FileAccepts []string `yaml:"file_accepts,omitempty"`

	// Options holds the value=Label choices for select/check fields, parsed
	// from the third position of name:select:v1=L1|v2=L2.
	Options []FieldOption `yaml:"options,omitempty"`

	// Workflow, when set on a select field, turns it from a column that
	// accepts any of its options into a state machine with guarded
	// transitions. Nil for an ordinary select.
	Workflow *WorkflowSpec `yaml:"workflow,omitempty"`

	// Auto marks a string field that is auto-generated from a sequence in the
	// model's BeforeCreate hook (declared name:string:auto or name:string:auto:PREFIX).
	// The generator wires the sequence infra + BeforeCreate call, makes the field
	// optional, and drops it from the create/edit form. AutoPrefix is the sequence
	// prefix (defaults to the first three letters of the resource, uppercased).
	Auto       bool   `yaml:"auto,omitempty"`
	AutoPrefix string `yaml:"auto_prefix,omitempty"`

	// DefaultCountry is the ISO 3166-1 alpha-2 code a tel field parses local
	// numbers against and its picker starts on, or the value a country field
	// starts on. From phone:tel:UG. Empty means no default.
	DefaultCountry string `yaml:"default_country,omitempty"`

	// Max is the number of stars on a rating field (score:rating:10). Zero
	// means DefaultRatingMax.
	Max int `yaml:"max,omitempty"`
}

// IsFormatted reports one of the formatted types: email, url, domain, tel,
// country, color, percent, rating, time and json.
func (f Field) IsFormatted() bool {
	switch FieldType(f.Type) {
	case FieldEmail, FieldURL, FieldDomain, FieldTel, FieldCountry,
		FieldColor, FieldPercent, FieldRating, FieldTime, FieldJSON:
		return true
	}
	return false
}

// IsTel reports a phone number field.
func (f Field) IsTel() bool { return FieldType(f.Type) == FieldTel }

// RatingMax is the number of stars a rating field offers.
func (f Field) RatingMax() int {
	if f.Max > 0 {
		return f.Max
	}
	return DefaultRatingMax
}

// FormatTag is the value of the model's format:"..." struct tag, which is what
// internal/fieldtypes reads to check and normalise the column on every write.
// Empty for a field that is not a formatted type.
func (f Field) FormatTag() string {
	if !f.IsFormatted() {
		return ""
	}
	switch FieldType(f.Type) {
	case FieldTel:
		if f.DefaultCountry != "" {
			return "tel:" + f.DefaultCountry
		}
	case FieldRating:
		return fmt.Sprintf("rating:%d", f.RatingMax())
	}
	return f.Type
}

// DocsTag is the gin-docs struct tag for a formatted field, so the OpenAPI
// reference names the format and shows a value that would be accepted.
func (f Field) DocsTag() string {
	switch FieldType(f.Type) {
	case FieldEmail:
		return "format:email,example:ada@example.com"
	case FieldURL:
		return "format:uri,example:https://example.com/pricing"
	case FieldDomain:
		return "format:hostname,example:example.co.ug"
	case FieldTel:
		return "format:e164,example:+256772123456"
	case FieldCountry:
		return "format:iso-3166-alpha-2,example:UG"
	case FieldColor:
		return "format:hex-color,example:#6c5ce7"
	case FieldPercent:
		return "description:A percentage from 0 to 100,example:12.5"
	case FieldRating:
		return fmt.Sprintf("description:Stars from 1 to %d (0 is unrated),example:4", f.RatingMax())
	case FieldTime:
		return "format:time,example:14:30"
	case FieldJSON:
		return "description:Any JSON value"
	}
	return ""
}

// SharedSchema names the shared Zod schema a formatted field validates with,
// and the module in packages/shared/schemas it comes from. Empty for the rest.
func (f Field) SharedSchema() (module, name string) {
	switch FieldType(f.Type) {
	case FieldEmail:
		return "./field-formats", "EmailSchema"
	case FieldURL:
		return "./field-formats", "UrlSchema"
	case FieldDomain:
		return "./field-formats", "DomainSchema"
	case FieldTel:
		return "./phone", "PhoneSchema"
	case FieldCountry:
		return "./field-formats", "CountrySchema"
	case FieldColor:
		return "./field-formats", "ColorSchema"
	case FieldPercent:
		return "./field-formats", "PercentSchema"
	case FieldRating:
		return "./field-formats", "ratingSchema"
	case FieldTime:
		return "./field-formats", "TimeSchema"
	case FieldJSON:
		return "./field-formats", "JsonValueSchema"
	}
	return "", ""
}

// IsAuto reports a string field that is auto-generated from a sequence.
func (f Field) IsAuto() bool { return f.Auto }

// FieldOption is one choice in a select or check field: Value is stored,
// Label is shown.
type FieldOption struct {
	Value string `yaml:"value"`
	Label string `yaml:"label"`
}

// IsSelect reports a single-choice dropdown field.
func (f Field) IsSelect() bool { return FieldType(f.Type) == FieldSelect }

// IsRadio reports a single-choice radio-button field.
func (f Field) IsRadio() bool { return FieldType(f.Type) == FieldRadio }

// IsCheck reports a multi-choice checkbox-group field (stored as a JSON array).
func (f Field) IsCheck() bool { return FieldType(f.Type) == FieldCheck }

// IsToggle reports an on/off boolean field.
func (f Field) IsToggle() bool { return FieldType(f.Type) == FieldToggle }

// HasOptions reports whether this field carries a value=Label option list.
func (f Field) HasOptions() bool { return len(f.Options) > 0 }

// OptionValues returns just the stored values, in declared order — used to
// build Zod enums and TS unions.
func (f Field) OptionValues() []string {
	vals := make([]string, len(f.Options))
	for i, o := range f.Options {
		vals[i] = o.Value
	}
	return vals
}

// tsUnion renders the option values as a TypeScript string-literal union,
// e.g. `"draft" | "sent" | "paid"`.
func (f Field) tsUnion() string {
	parts := make([]string, len(f.Options))
	for i, o := range f.Options {
		parts[i] = fmt.Sprintf("%q", o.Value)
	}
	return strings.Join(parts, " | ")
}

// zodEnumList renders the option values as a comma-separated quoted list for
// z.enum([...]), e.g. `"draft", "sent", "paid"`.
func (f Field) zodEnumList() string {
	parts := make([]string, len(f.Options))
	for i, o := range f.Options {
		parts[i] = fmt.Sprintf("%q", o.Value)
	}
	return strings.Join(parts, ", ")
}

// OptionsLiteral renders the options as a TS array-of-objects for the admin
// form definition: `[{ value: "draft", label: "Draft" }, ...]`.
func (f Field) OptionsLiteral() string {
	parts := make([]string, len(f.Options))
	for i, o := range f.Options {
		parts[i] = fmt.Sprintf("{ value: %q, label: %q }", o.Value, o.Label)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// IsSlug returns true if this field is an auto-generated slug.
func (f Field) IsSlug() bool {
	return FieldType(f.Type) == FieldSlug
}

// IsBelongsTo returns true if this field is a belongs_to relationship.
//
// one_to_one counts. It is a belongs_to with a unique foreign key, and every
// consumer of this, the preload, the admin picker, the CSV import, the mobile
// and desktop clients, wants exactly the same treatment. Only the GORM tag and
// the duplicate-key error message differ.
func (f Field) IsBelongsTo() bool {
	t := FieldType(f.Type)
	return t == FieldBelongsTo || t == FieldOneToOne
}

// IsOneToOne returns true only for the unique variant.
func (f Field) IsOneToOne() bool {
	return FieldType(f.Type) == FieldOneToOne
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
	return FieldType(f.Type) == FieldStringArray || FieldType(f.Type) == FieldCheck || FieldType(f.Type) == FieldJSON
}

// NeedsFilesImport returns true if this field requires the models/files
// package import (for the FileRef type) in the generated Go model.
func (f Field) NeedsFilesImport() bool {
	return f.IsFileField()
}

// GoType returns the Go type for this field.
func (f Field) GoType() string {
	// Still text to everyone but the database. The type carries the
	// encryption, so the model, both request structs and the create literal
	// agree about it without a conversion anywhere; a plain string in any one
	// of them was a compile error, and in an update map, plaintext at rest.
	if f.Encrypted {
		return "crypto.EncryptedString"
	}
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
	case FieldMoney:
		return "money.Money"
	case FieldBool, FieldToggle:
		return "bool"
	case FieldSelect, FieldRadio:
		return "string"
	case FieldCheck:
		return "datatypes.JSONSlice[string]"
	case FieldDate:
		// Not *time.Time. time.Time.UnmarshalJSON accepts RFC3339 and nothing
		// else, while the admin's picker sends "2001-08-06" on purpose, so
		// every date field failed to save. See internal/jsontime.
		return "*jsontime.Date"
	case FieldDatetime:
		// Same reason: <input type="datetime-local"> sends "2001-08-06T14:30",
		// which has neither seconds nor a zone and is not RFC3339 either.
		return "*jsontime.DateTime"
	case FieldManyToMany:
		return "[]string"
	case FieldStringArray:
		return "datatypes.JSONSlice[string]"
	case FieldFile:
		return "*files.FileRef"
	case FieldFiles:
		return "files.FileRefs"
	case FieldEmail, FieldURL, FieldDomain, FieldTel, FieldCountry, FieldColor, FieldTime:
		return "string"
	case FieldPercent:
		return "float64"
	case FieldRating:
		return "int"
	case FieldJSON:
		// jsonb on Postgres, JSON on MySQL and SQLite, which stores it as text.
		return "datatypes.JSON"
	default:
		return "string"
	}
}

// GORMTag returns the GORM struct tag for this field.
func (f Field) GORMTag() string {
	// Ciphertext is longer than the text it holds (a nonce, a tag, base64 and a
	// prefix), so a size:255 column would reject a 200-character value.
	if f.Encrypted {
		return "type:text"
	}
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
	case FieldMoney:
		// Two columns behind one field: price_amount BIGINT and price_currency.
		// One column holding "19.99 USD" is not something you can SUM.
		return "" // the embedded tag is added by the model template
	case FieldOneToOne:
		// uniqueIndex, not index. Without it this is a belongs_to wearing a
		// different name: the database would happily accept a second row
		// pointing at the same parent, and "one to one" would be a comment
		// rather than a constraint.
		parts = append(parts, "size:36", "uniqueIndex")
	case FieldStringArray, FieldCheck:
		parts = append(parts, "type:json")
	case FieldSelect, FieldRadio:
		parts = append(parts, "size:255")
	case FieldFile, FieldFiles:
		// FileRef / FileRefs implement Value / Scan via the files package,
		// so GORM stores them as JSON. type:json signals jsonb on Postgres
		// (otherwise we'd get text and lose efficient querying).
		parts = append(parts, "type:json")
	case FieldEmail:
		// The longest address SMTP allows.
		parts = append(parts, "size:254")
	case FieldURL:
		parts = append(parts, "size:2048")
	case FieldDomain:
		// The longest name DNS allows.
		parts = append(parts, "size:253")
	case FieldTel:
		// E.164 is at most 16 characters with the plus.
		parts = append(parts, "size:32")
	case FieldCountry:
		parts = append(parts, "size:2")
	case FieldColor:
		parts = append(parts, "size:7")
	case FieldTime:
		parts = append(parts, "size:5")
	case FieldPercent:
		// 0.00 to 100.00 exactly, which a float column would not promise.
		parts = append(parts, "type:decimal(5,2)")
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
	case FieldBool, FieldToggle:
		return "boolean"
	case FieldSelect, FieldRadio:
		if f.HasOptions() {
			return f.tsUnion()
		}
		return "string"
	case FieldCheck:
		if f.HasOptions() {
			return "(" + f.tsUnion() + ")[]"
		}
		return "string[]"
	case FieldDatetime, FieldDate:
		return "string | null"
	case FieldManyToMany:
		return "string[]"
	case FieldStringArray:
		return "string[]"
	case FieldMoney:
		return "Money"
	case FieldFile:
		return "FileRef | null"
	case FieldFiles:
		return "FileRef[]"
	case FieldPercent, FieldRating:
		return "number"
	case FieldJSON:
		return "unknown"
	default:
		// email, url, domain, tel, country, color and time are strings.
		return "string"
	}
}

// ZodType returns the Zod validator for this field.
func (f Field) ZodType() string {
	if f.IsFormatted() {
		return f.formattedZodType()
	}
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
	case FieldBool, FieldToggle:
		base = "z.boolean()"
	case FieldSelect, FieldRadio:
		if f.HasOptions() {
			base = "z.enum([" + f.zodEnumList() + "])"
		} else {
			base = "z.string()"
		}
	case FieldCheck:
		if f.HasOptions() {
			base = "z.array(z.enum([" + f.zodEnumList() + "])).optional()"
		} else {
			base = "z.array(z.string()).optional()"
		}
	case FieldDatetime, FieldDate:
		base = "z.string().nullable()"
	case FieldBelongsTo:
		// FK columns are UUID strings matching the referenced model's PK.
		base = `z.string().uuid("Invalid ID")`
	case FieldManyToMany:
		base = "z.array(z.string().uuid()).optional()"
	case FieldStringArray:
		base = "z.array(z.string()).optional()"
	case FieldMoney:
		// Not z.number(). The currency is half the value, and a schema that
		// validates only the amount lets a client post a price with no
		// currency at all, which the API would then default to USD.
		base = "MoneySchema"
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

	if !f.Required && FieldType(f.Type) != FieldDatetime && FieldType(f.Type) != FieldDate && FieldType(f.Type) != FieldSlug && FieldType(f.Type) != FieldRichtext && FieldType(f.Type) != FieldBelongsTo && FieldType(f.Type) != FieldManyToMany && FieldType(f.Type) != FieldStringArray && FieldType(f.Type) != FieldCheck && FieldType(f.Type) != FieldFile && FieldType(f.Type) != FieldFiles {
		base += ".optional()"
	}

	return base
}

// formattedZodType is the Zod validator for a formatted field: the shared
// schema, and for an optional one the empty value a form sends when the field
// is left blank.
func (f Field) formattedZodType() string {
	_, name := f.SharedSchema()
	switch FieldType(f.Type) {
	case FieldRating, FieldPercent:
		if FieldType(f.Type) == FieldRating {
			name = fmt.Sprintf("ratingSchema(%d)", f.RatingMax())
		}
		if !f.Required {
			return name + ".nullable().optional()"
		}
		return name
	case FieldJSON:
		if !f.Required {
			return name + ".optional()"
		}
		return name
	}
	if !f.Required {
		return name + `.or(z.literal("")).optional()`
	}
	return name
}

// NeedsTimeImport returns true if this field requires "time" import in Go.
func (f Field) NeedsTimeImport() bool {
	return FieldType(f.Type) == FieldDatetime || FieldType(f.Type) == FieldDate
}

// NeedsJSONTimeImport reports whether this field's Go type comes from
// internal/jsontime.
func (f Field) NeedsJSONTimeImport() bool {
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
	case FieldMoney:
		return "money"
	case FieldFile:
		return "file"
	case FieldFiles:
		return "files"
	case FieldEmail:
		return "email"
	case FieldURL:
		return "link"
	case FieldDomain, FieldTel, FieldCountry, FieldColor, FieldPercent, FieldRating, FieldTime, FieldJSON:
		return f.Type
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
	case FieldBool, FieldToggle:
		return "toggle"
	case FieldSelect:
		return "select"
	case FieldRadio:
		return "radio"
	case FieldCheck:
		return "checkbox-group"
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
	case FieldMoney:
		return "money"
	case FieldFile:
		return "file"
	case FieldFiles:
		return "files"
	case FieldEmail, FieldURL, FieldDomain, FieldTel, FieldCountry, FieldColor,
		FieldPercent, FieldRating, FieldTime, FieldJSON:
		// Each has an input of its own, named after the type.
		return f.Type
	default:
		return "text"
	}
}

// IsSortable returns true if this field type should be sortable by default.
func (f Field) IsSortable() bool {
	switch FieldType(f.Type) {
	case FieldString, FieldInt, FieldUint, FieldFloat, FieldDatetime, FieldDate, FieldSlug:
		return true
	case FieldEmail, FieldURL, FieldDomain, FieldTel, FieldCountry, FieldColor,
		FieldPercent, FieldRating, FieldTime:
		return true
	case FieldMoney:
		// Sorts on <field>_amount, which the handler whitelists under that
		// name. The ordering is exact because the amount is an integer.
		return true
	default:
		return false
	}
}

// IsSearchable returns true if this field type should be searchable by default.
func (f Field) IsSearchable() bool {
	// Ciphertext is different on every write, so LIKE can never match it.
	if f.Encrypted {
		return false
	}
	switch FieldType(f.Type) {
	case FieldString, FieldText, FieldSlug, FieldRichtext,
		FieldEmail, FieldURL, FieldDomain, FieldTel:
		return true
	}
	return false
}

// ValidFieldTypes returns all valid field type names.
func ValidFieldTypes() []string {
	return []string{"string", "text", "richtext", "int", "uint", "float", "bool", "toggle", "select", "radio", "check", "datetime", "date", "money", "slug", "belongs_to", "one_to_one", "many_to_many", "string_array", "file", "files",
		"email", "url", "domain", "tel", "country", "color", "percent", "rating", "time", "json"}
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

// HasWorkflow reports whether this field is a state machine.
func (f Field) HasWorkflow() bool { return f.Workflow != nil }
