package generate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// Sync parses Go model files and regenerates TypeScript types and Zod schemas.
func Sync() error {
	root, err := findProjectRoot()
	if err != nil {
		return err
	}

	modelsDir := filepath.Join(root, "apps", "api", "internal", "models")
	if !dirExists(modelsDir) {
		return fmt.Errorf("models directory not found at %s", modelsDir)
	}

	sharedRoot := sharedRootFrom(root)
	if !dirExists(sharedRoot) {
		return fmt.Errorf("shared package not found at %s: sync requires a monorepo project", sharedRoot)
	}

	fmt.Print("\n  Syncing Go types → TypeScript...\n\n")

	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		return fmt.Errorf("reading models directory: %w", err)
	}

	var synced int
	var adminFieldsAdded int
	var adminWarnings []string
	var uiDrift []clientUIDrift
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		filePath := filepath.Join(modelsDir, entry.Name())
		structs, err := parseGoStructs(filePath)
		if err != nil {
			fmt.Printf("  ⚠ Skipping %s: %v\n", entry.Name(), err)
			continue
		}

		for _, s := range structs {
			// Skip User — it has custom schemas already
			if s.Name == "User" {
				continue
			}

			kebab := strings.ReplaceAll(toSnakeCase(s.Name), "_", "-")

			// Write TypeScript type
			tsContent := buildTSType(s)
			tsPath := filepath.Join(sharedRoot, "types", kebab+".ts")
			if err := writeFileWithDirs(tsPath, tsContent); err != nil {
				return fmt.Errorf("writing %s: %w", tsPath, err)
			}
			fmt.Printf("  ✓ packages/shared/types/%s.ts\n", kebab)

			// Write Zod schema
			zodContent := buildZodSchema(s)
			zodPath := filepath.Join(sharedRoot, "schemas", kebab+".ts")
			if err := writeFileWithDirs(zodPath, zodContent); err != nil {
				return fmt.Errorf("writing %s: %w", zodPath, err)
			}
			fmt.Printf("  ✓ packages/shared/schemas/%s.ts\n", kebab)

			// v3.31.16: marker-based partial regen of the admin resource
			// file. Adds any new model fields to the auto-section of
			// columns + form; leaves customised entries alone.
			added, warnings, syncErr := SyncAdminResource(root, s)
			if syncErr != nil {
				fmt.Printf("  ⚠ admin auto-add for %s: %v\n", s.Name, syncErr)
			} else if added > 0 {
				pluralKebab := strings.ReplaceAll(Pluralize(toSnakeCase(s.Name)), "_", "-")
				// The definition sits in a folder of its own on a current
				// project and flat on an older one. SyncAdminResource just
				// resolved that to edit the file; resolve it again to name it,
				// rather than printing a path built from the slug that is
				// wrong in whichever layout the guess did not pick.
				def := scaffold.FindResourceDef(filepath.Join(adminRootFrom(root), "resources"), pluralKebab)
				fmt.Printf("  ✓ %s  (added %d field%s to columns + form)\n",
					relFromRoot(root, def), added, plural(added))
				adminFieldsAdded += added
			}
			adminWarnings = append(adminWarnings, warnings...)

			// The admin was just brought up to date in place. The mobile and
			// desktop screens cannot be, so say which fields they never saw.
			uiDrift = append(uiDrift, findClientUIDrift(root, s)...)

			synced++
		}
	}

	if synced == 0 {
		fmt.Println("  No models found to sync (User is skipped).")
	} else {
		fmt.Printf("\n  ✅ Synced %d model(s) to TypeScript + Zod\n", synced)
		if adminFieldsAdded > 0 {
			fmt.Printf("  ✅ Auto-added %d field%s to admin resource files\n",
				adminFieldsAdded, plural(adminFieldsAdded))
		}
		for _, w := range adminWarnings {
			fmt.Println(w)
		}
		fmt.Println()
		printClientUIDrift(uiDrift)
	}

	return nil
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// GoStruct represents a parsed Go struct.
type GoStruct struct {
	Name   string
	Fields []GoField
}

// GoField represents a parsed Go struct field.
type GoField struct {
	Name     string
	GoType   string
	JSONName string
	GORMTag  string
	// Format is the formatted field type the column was generated as (email,
	// tel:UG, rating:10, ...), read from its format tag. Empty for the rest.
	Format string
	// Enum is the pipe-separated values a choice column may hold, read from
	// its enum tag. Empty for every other column.
	Enum string

	// The rest of what the model already says about the column, which sync
	// used to ignore and then write a schema that contradicted (grit#108).
	//
	// Binding is the gin binding tag: "required" is the only part that
	// matters here. RequiredIf is the condition under which the column is
	// required, as "status=lost". Sanitize tells rich text from plain text,
	// which are the same Go type. Slug names the column a slug is derived
	// from, and Owner marks the owner column of an --owned-by resource:
	// both are filled by the server, and neither is otherwise distinguishable
	// from an ordinary string.
	Binding    string
	Required   bool
	RequiredIf string
	Sanitize   string
	Slug       string
	Owner      bool
}

// parseGoStructs parses all struct declarations from a Go file.
// ParseGoStructs is the public entry point — needed by the
// internal/expose package which scaffolds web pages from a Go model.
func ParseGoStructs(filePath string) ([]GoStruct, error) {
	return parseGoStructs(filePath)
}

func parseGoStructs(filePath string) ([]GoStruct, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", filePath, err)
	}

	var structs []GoStruct

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			s := GoStruct{Name: typeSpec.Name.Name}

			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue // embedded field
				}

				goField := GoField{
					Name:   field.Names[0].Name,
					GoType: typeToString(field.Type),
				}

				// Parse struct tags
				if field.Tag != nil {
					tag := field.Tag.Value
					tag = strings.Trim(tag, "`")

					goField.JSONName = extractTag(tag, "json")
					goField.GORMTag = extractTag(tag, "gorm")
					goField.Format = extractTag(tag, "format")
					goField.Enum = extractTag(tag, "enum")
					goField.Binding = extractTag(tag, "binding")
					goField.Required = extractTag(tag, "required") == "true"
					goField.RequiredIf = extractTag(tag, "requiredif")
					goField.Sanitize = extractTag(tag, "sanitize")
					goField.Slug = extractTag(tag, "slug")
					goField.Owner = extractTag(tag, "owner") == "true"
				}

				// Skip fields with json:"-"
				if goField.JSONName == "-" {
					continue
				}

				// Default JSON name to snake_case of Go name
				if goField.JSONName == "" {
					goField.JSONName = toSnakeCase(goField.Name)
				}

				s.Fields = append(s.Fields, goField)
			}

			structs = append(structs, s)
		}
	}

	return structs, nil
}

func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt)
	case *ast.IndexExpr:
		return typeToString(t.X) + "[" + typeToString(t.Index) + "]"
	default:
		return "interface{}"
	}
}

func extractTag(tag, key string) string {
	search := key + `:"`
	idx := strings.Index(tag, search)
	if idx == -1 {
		return ""
	}
	rest := tag[idx+len(search):]
	end := strings.Index(rest, `"`)
	if end == -1 {
		return ""
	}
	value := rest[:end]
	// json:"name,omitempty" → just "name"
	if comma := strings.Index(value, ","); comma != -1 {
		value = value[:comma]
	}
	return value
}

func buildTSType(s GoStruct) string {
	var fields strings.Builder
	for _, f := range s.Fields {
		tsType := goTypeToTS(f.GoType)
		fields.WriteString(fmt.Sprintf("  %s: %s;\n", f.JSONName, tsType))
	}

	return fmt.Sprintf("export interface %s {\n%s}\n", s.Name, fields.String())
}

// buildZodSchema writes a model's create, update and entity schemas.
//
// It decides nothing about a field. inferField reconstructs the resource field
// the column was generated from and Field.ZodType() renders it, which is the
// same call the generator makes, so the file sync writes says what the file the
// generator wrote said. See sync_infer.go for why that matters.
func buildZodSchema(s GoStruct) string {
	var createFields strings.Builder
	var updateFields strings.Builder
	var refinements []conditionalRequired
	var shaped []Field

	for _, f := range s.Fields {
		// Skip auto-generated fields
		if isAutoField(f.JSONName) {
			continue
		}
		// An association is not a field. Including it put `contact:
		// z.unknown()` beside `contact_id` in the schema, which validates
		// nothing and invites a form to send the whole related record back.
		if isAssociation(f, s.Fields) {
			continue
		}
		// Filled by the server. A slug is derived from another column by a
		// hook, and the owner of an --owned-by row is the signed-in user,
		// whatever the body said. Asking a create form for either is asking
		// for a value that will be discarded, and the form cannot be
		// submitted without inventing one.
		if f.Slug != "" || f.Owner {
			continue
		}

		field, zodType, ok := syncZod(f, s)
		if ok {
			shaped = append(shaped, field)

			// "Required while it is shown" is a predicate, not a flag.
			// Marking it required in the object rejects every row the
			// condition does not hold for.
			if field.ShowWhen != nil && field.Required {
				refinements = append(refinements, conditionalRequired{
					Field: f.JSONName,
					On:    toSnakeCase(field.ShowWhen.Field),
					Value: field.ShowWhen.Value,
				})
				if !strings.Contains(zodType, ".optional()") {
					zodType += ".optional()"
				}
			}
		}

		// The JSON tag's own name, not camelCase.
		//
		// The API decodes with ShouldBindJSON against `json:"contact_id"`, so
		// a payload keyed contactId is a payload the API ignores: the field
		// arrives empty and the error names a field the form did fill in.
		// The generator always wrote snake_case here and sync rewrote the
		// same file in camelCase, so whichever ran last decided whether the
		// app worked (grit#102).
		key := f.JSONName
		if key == "" {
			key = toSnakeCase(f.Name)
		}

		createFields.WriteString(fmt.Sprintf("  %s: %s,\n", key, zodType))

		updateZod := zodType
		if !strings.Contains(updateZod, ".optional()") && !strings.Contains(updateZod, ".nullable()") {
			updateZod += ".optional()"
		}
		updateFields.WriteString(fmt.Sprintf("  %s: %s,\n", key, updateZod))
	}

	// The entity schema is emitted alongside the create and update pair.
	//
	// Not decoration: the scaffold's own packages/shared/schemas/index.ts
	// re-exports <Name>Schema for every resource it ships. Writing this file
	// without one left the barrel importing a member that no longer existed,
	// which is a TS2305 inside @repo/shared, which fails the type-check of the
	// admin AND the web app. The reported symptom was that sync broke the
	// login page, and nothing in the project imported the missing schema at
	// all: the stale re-export was the entire fault.
	//
	// It is also the useful one. Create and update describe what goes in;
	// this describes what comes back, id and timestamps included, which is
	// what you validate an API response against. A slug and an owner column
	// belong here for the same reason: the server fills them, and then sends
	// them.
	var entityFields strings.Builder
	for _, f := range s.Fields {
		// Keyed by the JSON tag, like the create and update pair above and
		// like everything the API actually sends. This one was camelCase too,
		// and it is the schema you validate a response against, so it
		// rejected every response the API had just produced.
		key := f.JSONName
		if key == "" {
			key = toSnakeCase(f.Name)
		}
		// An association appears in a response, so unlike create and update
		// it belongs here. It is another resource's shape, though, and this
		// builder has no way to reach for it, so it stays permissive rather
		// than pretending to validate.
		zodType := fieldZod(f)
		switch {
		case isAssociation(f, s.Fields):
			zodType = "z.unknown().optional()"
		case isAutoField(f.JSONName):
			// id, the timestamps and the version, from their Go types alone.
			// These are the framework's columns and the API always sends
			// them: inferring a resource field from one would make id
			// optional and created_at nullable, which is a response this API
			// does not produce.
		default:
			if field, zod, ok := syncZod(f, s); ok {
				shaped = append(shaped, field)
				zodType = zod
			}
		}
		entityFields.WriteString(fmt.Sprintf("  %s: %s,\n", key, zodType))
	}

	return fmt.Sprintf(`%s

export const %sSchema = z.object({
%s});

export const Create%sSchema = z.object({
%s})%s;

export const Update%sSchema = z.object({
%s});

export type %s = z.infer<typeof %sSchema>;
export type Create%sInput = z.infer<typeof Create%sSchema>;
export type Update%sInput = z.infer<typeof Update%sSchema>;
`, zodSchemaImports(shaped),
		s.Name, entityFields.String(),
		s.Name, createFields.String(), zodRefinements(refinements),
		s.Name, updateFields.String(),
		s.Name, s.Name, s.Name, s.Name, s.Name, s.Name)
}

// syncZod renders one column, by reconstructing the resource field it was
// generated from and asking the generator's own emitter.
//
// Two rules live here rather than in Field.ZodType, because the generator
// applies them at its call site too: a self-reference has to accept the empty
// string, since that is how "this is a root" is spelled, and a belongs_to never
// takes .optional() from ZodType, so an optional relationship says so itself.
//
// ok is false when the Go type says nothing a resource field can be inferred
// from, and the caller keeps the plain type mapping.
func syncZod(f GoField, s GoStruct) (Field, string, bool) {
	field, ok := inferField(f, s.Fields)
	if !ok {
		return field, fieldZod(f), false
	}
	zod := field.ZodType()
	if FieldType(field.Type) == FieldBelongsTo {
		if selfReferencing(f, s) {
			zod += `.or(z.literal(""))`
		}
		if !field.Required {
			zod += ".nullable().optional()"
		}
	}
	return field, zod, true
}

func goTypeToTS(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "*time.Time":
		return "string | null"
	case "time.Time":
		return "string"
	case "gorm.DeletedAt":
		return "string | null"
	// The API's date wrappers marshal to a string. Without these they fell
	// through to unknown, so every date column in a generated type was
	// something the frontend had to cast before it could be read.
	case "jsontime.DateTime", "jsontime.Date":
		return "string"
	case "*jsontime.DateTime", "*jsontime.Date":
		return "string | null"
	default:
		if strings.Contains(goType, "JSONSlice[string]") {
			return "string[]"
		}
		if strings.HasPrefix(goType, "[]") {
			inner := goTypeToTS(goType[2:])
			return inner + "[]"
		}
		if strings.HasPrefix(goType, "*") {
			inner := goTypeToTS(goType[1:])
			return inner + " | null"
		}
		return "unknown"
	}
}

// fieldZod is a parsed field's Zod validator: its format's rule when it has a
// format tag, its Go type's otherwise. Self-contained, because sync writes
// schema files that import nothing but zod.
func fieldZod(f GoField) string {
	kind, param, _ := strings.Cut(f.Format, ":")
	switch FieldType(kind) {
	case FieldEmail:
		return `z.string().email().or(z.literal(""))`
	case FieldURL:
		return `z.string().url().regex(/^https?:\/\//i).or(z.literal(""))`
	case FieldDomain:
		return `z.string().regex(/^([a-z0-9-]+\.)+[a-z0-9-]{2,63}$/i).or(z.literal(""))`
	case FieldTel:
		return `z.string().regex(/^\+[1-9][0-9]{6,14}$/).or(z.literal(""))`
	case FieldCountry:
		return `z.string().regex(/^[A-Z]{2}$/).or(z.literal(""))`
	case FieldColor:
		return `z.string().regex(/^#[0-9a-f]{6}$/).or(z.literal(""))`
	case FieldTime:
		return `z.string().regex(/^([01][0-9]|2[0-3]):[0-5][0-9]$/).or(z.literal(""))`
	case FieldPercent:
		return "z.number().min(0).max(100)"
	case FieldRating:
		max := DefaultRatingMax
		if n, err := strconv.Atoi(param); err == nil && n > 0 {
			max = n
		}
		return fmt.Sprintf("z.number().int().min(0).max(%d)", max)
	case FieldJSON:
		return "z.unknown()"
	}
	return goTypeToZod(f.GoType, f.GORMTag)
}

func goTypeToZod(goType, gormTag string) string {
	switch goType {
	case "string":
		if strings.Contains(gormTag, "type:text") {
			return "z.string()"
		}
		return "z.string()"
	case "int", "int8", "int16", "int32", "int64":
		return "z.number().int()"
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return "z.number().int().nonnegative()"
	case "float32", "float64":
		return "z.number()"
	case "bool":
		return "z.boolean()"
	case "*time.Time":
		return "z.string().nullable()"
	case "time.Time":
		return "z.string()"
	case "datatypes.Date", "*datatypes.Date":
		return "z.string()"
	// The API's own date and datetime wrappers. They marshal to a string, and
	// falling through to z.unknown() meant the schema validated nothing for
	// every date column a resource had.
	case "jsontime.DateTime", "jsontime.Date":
		return "z.string()"
	case "*jsontime.DateTime", "*jsontime.Date":
		return "z.string().nullable()"
	default:
		if strings.Contains(goType, "JSONSlice[string]") {
			return "z.array(z.string()).optional()"
		}
		// A pointer is the column saying it may be NULL, which is how an
		// optional foreign key is written. Without this an optional
		// belongs_to came out as z.unknown(): the schema validated nothing,
		// so the one surface that should have caught a bad relationship id
		// accepted anything at all (grit#100).
		if strings.HasPrefix(goType, "*") {
			inner := goTypeToZod(strings.TrimPrefix(goType, "*"), gormTag)
			if inner == "z.unknown()" {
				return inner
			}
			return inner + ".nullable().optional()"
		}
		return "z.unknown()"
	}
}

func isAutoField(jsonName string) bool {
	auto := map[string]bool{
		"id":         true,
		"created_at": true,
		"updated_at": true,
		"deleted_at": true,
		// The framework's own archive bookkeeping, which sync was offering as
		// an editable date on every form (grit#101).
		"archived_at": true,
		// Optimistic locking. The row's version travels as an ETag and comes
		// back as If-Match, never in the body: the service drops it from
		// every write. Required on create, as sync had it, it asked a form
		// for a number nobody has and refused to submit without one
		// (grit#108).
		"version": true,
	}
	return auto[jsonName]
}
