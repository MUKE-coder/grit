package generate

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

// SeederOptions controls how a resource seeder is generated.
type SeederOptions struct {
	Faker bool // use gofakeit to generate many records instead of one static example
	Count int  // number of faker records (default 10)
}

// WriteSeeder generates internal/database/<plural>_seeder.go for the resource
// and registers its Seed<Plural> call in seed.go. Called by the generator
// (via --seed/--faker) and by the standalone "grit generate seeder" command.
func (g *Generator) WriteSeeder(opts SeederOptions) error {
	names := MakeNames(g.Definition.Name)
	content := g.seederContent(names, opts)
	path := filepath.Join(g.APIRoot(), "internal", "database", names.PluralKebab+"_seeder.go")
	if err := writeFileWithDirs(path, content); err != nil {
		return fmt.Errorf("writing seeder: %w", err)
	}
	if err := g.injectSeederCall(names); err != nil {
		return fmt.Errorf("registering seeder: %w", err)
	}
	fmt.Printf("  ✓ apps/api/internal/database/%s_seeder.go\n", names.PluralKebab)
	return nil
}

// injectSeederCall adds a Seed<Plural> call to seed.go's Seed() at the
// grit:seeders marker. Idempotent (injectBefore skips duplicates).
func (g *Generator) injectSeederCall(names Names) error {
	seedPath := filepath.Join(g.APIRoot(), "internal", "database", "seed.go")
	if !fileExists(seedPath) {
		return nil
	}
	call := "\tif err := Seed" + names.PluralPascal + "(db); err != nil {\n" +
		"\t\treturn fmt.Errorf(\"seeding " + names.Plural + ": %w\", err)\n" +
		"\t}\n"
	return injectBefore(seedPath, "// grit:seeders", call)
}

func (g *Generator) seederContent(names Names, opts SeederOptions) string {
	count := opts.Count
	if count <= 0 {
		count = 10
	}
	mode := "static"
	if opts.Faker {
		mode = "faker"
	}
	fieldLines, relPreamble, needsTime, needsFiles, needsStrings := g.seederFieldLines(mode)
	if relPreamble != "" {
		relPreamble = "\t// Link each row to an existing parent (loaded once).\n" + relPreamble + "\n"
	}

	// Assemble imports — only what the record actually references.
	var imports strings.Builder
	// fmt is needed by the "no parent rows exist yet" guard, which only appears
	// when the resource has a belongs_to. Keyed off the preamble rather than a
	// separate flag, because that guard is the only thing in a seeder that uses
	// it and an unused import will not compile.
	if relPreamble != "" {
		imports.WriteString("\t\"fmt\"\n")
	}
	imports.WriteString("\t\"log\"\n")
	if needsStrings {
		imports.WriteString("\t\"strings\"\n")
	}
	if needsTime {
		imports.WriteString("\t\"time\"\n")
	}
	imports.WriteString("\n")
	if opts.Faker {
		imports.WriteString("\t\"github.com/brianvoe/gofakeit/v7\"\n")
	}
	imports.WriteString("\t\"" + g.Module + "/internal/models\"\n")
	if needsFiles {
		imports.WriteString("\t\"" + g.Module + "/internal/files\"\n")
	}
	imports.WriteString("\t\"gorm.io/gorm\"\n")

	lower := strings.ToLower(names.Lower)
	header := "package database\n\nimport (\n" + imports.String() + ")\n\n"

	if opts.Faker {
		return header +
			"// Seed" + names.PluralPascal + " inserts fake " + names.Plural + " using gofakeit.\n" +
			"// Change the count (n) or swap the gofakeit calls for your own values.\n" +
			"func Seed" + names.PluralPascal + "(db *gorm.DB) error {\n" +
			"\tvar count int64\n" +
			"\tdb.Model(&models." + names.Pascal + "{}).Count(&count)\n" +
			"\tif count > 0 {\n" +
			"\t\tlog.Println(\"" + names.PluralPascal + " already seeded, skipping...\")\n" +
			"\t\treturn nil\n" +
			"\t}\n\n" +
			relPreamble +
			"\tconst n = " + fmt.Sprintf("%d", count) + "\n" +
			"\tfor i := 0; i < n; i++ {\n" +
			"\t\tr := models." + names.Pascal + "{\n" +
			fieldLines +
			"\t\t}\n" +
			"\t\tif err := db.Create(&r).Error; err != nil {\n" +
			"\t\t\tlog.Printf(\"Warning: failed to seed " + lower + ": %v\", err)\n" +
			"\t\t}\n" +
			"\t}\n" +
			"\tlog.Printf(\"Seeded %d " + lower + "\", n)\n" +
			"\treturn nil\n" +
			"}\n"
	}

	return header +
		"// Seed" + names.PluralPascal + " inserts sample " + names.Plural + ".\n" +
		"// Edit the values below or add more entries to the slice. Run with\n" +
		"// \"grit seed\" (or on migrate). Pass --faker to grit to generate many\n" +
		"// rows with gofakeit instead.\n" +
		"func Seed" + names.PluralPascal + "(db *gorm.DB) error {\n" +
		"\tvar count int64\n" +
		"\tdb.Model(&models." + names.Pascal + "{}).Count(&count)\n" +
		"\tif count > 0 {\n" +
		"\t\tlog.Println(\"" + names.PluralPascal + " already seeded, skipping...\")\n" +
		"\t\treturn nil\n" +
		"\t}\n\n" +
		relPreamble +
		"\trecords := []models." + names.Pascal + "{\n" +
		"\t\t{\n" +
		fieldLines +
		"\t\t},\n" +
		"\t}\n\n" +
		"\tfor _, r := range records {\n" +
		"\t\tif err := db.Create(&r).Error; err != nil {\n" +
		"\t\t\tlog.Printf(\"Warning: failed to seed " + lower + ": %v\", err)\n" +
		"\t\t}\n" +
		"\t}\n" +
		"\tlog.Printf(\"Seeded %d " + lower + "(s)\", len(records))\n" +
		"\treturn nil\n" +
		"}\n"
}

// DefinitionFromModel reconstructs a minimal ResourceDefinition by parsing an
// already-generated Go model. Used by the standalone "grit generate seeder"
// command (which doesn't get --fields). It locates the project the same way
// the generator does, so it works across every architecture.
func DefinitionFromModel(name string) (*ResourceDefinition, error) {
	root, err := findProjectRoot()
	if err != nil {
		return nil, err
	}
	names := MakeNames(name)
	for _, apiRoot := range []string{filepath.Join(root, "apps", "api"), root} {
		if modelFileFor(apiRoot, names) != "" {
			return definitionFromModelFile(apiRoot, name)
		}
	}
	return nil, fmt.Errorf("no model found for %q: generate the resource first (looked in apps/api/internal/models and internal/models)", names.Pascal)
}

// modelFileFor returns the path of a resource's model file, or "".
//
// Snake first, because that is what every generator in this package writes:
// filepath.Join(..., "models", names.Snake+".go"). Looking it up as Lower
// instead, which this file did until v3.168.0, is identical for a one-word
// resource and wrong for every other one: OrderItem is written to
// order_item.go and was looked for in orderitem.go, so `grit generate field`
// and `grit generate seeder` reported a model that was plainly on disk as
// missing.
//
// Lower is kept as a fallback for a model somebody wrote by hand under the
// flat name, which costs one stat call and saves a confusing failure.
func modelFileFor(apiRoot string, names Names) string {
	for _, candidate := range []string{names.Snake, names.Lower} {
		path := filepath.Join(apiRoot, "internal", "models", candidate+".go")
		if fileExists(path) {
			return path
		}
	}
	return ""
}

// definitionFromModelFile parses a specific model file into a definition.
func definitionFromModelFile(apiRoot, name string) (*ResourceDefinition, error) {
	names := MakeNames(name)
	modelPath := modelFileFor(apiRoot, names)
	if modelPath == "" {
		modelPath = filepath.Join(apiRoot, "internal", "models", names.Snake+".go")
	}
	structs, err := parseGoStructs(modelPath)
	if err != nil {
		return nil, fmt.Errorf("reading model for %s: %w (has the resource been generated?)", names.Pascal, err)
	}
	var target *GoStruct
	for i := range structs {
		if structs[i].Name == names.Pascal {
			target = &structs[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("model %s not found in %s", names.Pascal, modelPath)
	}

	// Columns the framework owns. None of them were declared by whoever
	// generated the resource, so none of them should turn up in a seeder for a
	// person to wonder about and then leave alone.
	//
	// archived_at is the one that bites: seeding it puts every row in the
	// archive, where the default list does not show them, and the seeder looks
	// like it worked. It is also *time.Time, so the generated file did not
	// compile, which is the only reason this was noticed.
	//
	// path, depth and position are the --tree columns. BeforeCreate computes
	// all three from the parent; a literal here describes a hierarchy that does
	// not exist.
	skip := map[string]bool{
		"id": true, "created_at": true, "updated_at": true,
		"deleted_at": true, "version": true, "slug": true,
		"archived_at": true,
		"path":        true,
		"depth":       true,
		"position":    true,
	}
	def := &ResourceDefinition{Name: names.Pascal}
	for _, gf := range target.Fields {
		jn := gf.JSONName
		if jn == "" || skip[jn] || strings.HasSuffix(jn, "_id") {
			continue // base column, slug (auto), or a FK we can't safely seed
		}
		ft := goTypeToFieldType(gf.GoType)
		if ft == "" {
			continue // relation / unknown type — user wires it up
		}
		// A nullable time is left NULL rather than seeded.
		//
		// Narrow on purpose. The value branch emits time.Now(), which does not
		// assign to a *time.Time, so the generated file would not compile. The
		// other pointer this hits is *files.FileRef, which seeds correctly as
		// &files.FileRef{...} and is worth keeping: a category with no image is
		// a duller demo than one with.
		if strings.HasPrefix(gf.GoType, "*") && ft == string(FieldDatetime) {
			continue
		}
		def.Fields = append(def.Fields, Field{Name: jn, Type: ft})
	}
	return def, nil
}

func goTypeToFieldType(t string) string {
	switch strings.TrimPrefix(t, "*") {
	case "string":
		return "string"
	case "int", "int64", "int32":
		return "int"
	case "uint", "uint64", "uint32":
		return "uint"
	case "float64", "float32":
		return "float"
	case "bool":
		return "bool"
	case "time.Time":
		return "datetime"
	case "files.FileRef":
		return "file"
	case "files.FileRefs":
		return "files"
	default:
		return ""
	}
}

// seederFieldLines returns the "GoField: value," lines for one record, a
// preamble that resolves belongs_to relations (loads the related ids so each
// row links to a real parent), and whether the time/files packages are needed.
// Slug (auto), m2m and string-array fields are skipped — the user wires those
// up by hand.
func (g *Generator) seederFieldLines(mode string) (lines, preamble string, needsTime, needsFiles, needsStrings bool) {
	var b, pre strings.Builder
	seenRel := map[string]bool{}
	for _, f := range g.Definition.Fields {
		if f.IsSlug() || f.IsManyToMany() || f.IsStringArray() {
			continue
		}
		// belongs_to: link to a real existing parent. We load the parent ids
		// once (Pluck) and pick one per row — random for faker, the first for
		// the static example. If none exist yet the FK is left empty.
		if f.IsBelongsTo() {
			rel := MakeNames(f.RelatedModelName())
			fkGo := toPascalCase(f.FKColumnName())
			idsVar := lowerCamel(rel.Lower) + "IDs"

			// A self-reference is left unset, so every seeded row is a root.
			//
			// Picking a random parent out of the table currently being seeded
			// builds a hierarchy nobody asked for, and can pick a row that ends
			// up being its own ancestor. Roots are the honest default: arrange
			// them by dragging in the admin, which is what the tree view is for.
			if rel.Pascal == toPascalCase(g.Definition.Name) {
				continue
			}

			if !seenRel[idsVar] {
				pre.WriteString("\tvar " + idsVar + " []string\n")
				pre.WriteString("\tdb.Model(&models." + rel.Pascal + "{}).Pluck(\"id\", &" + idsVar + ")\n")

				// A required foreign key with no parent rows to point at cannot
				// be seeded at all. Saying so beats writing "" and letting the
				// database answer with SQLSTATE 23503, which is what used to
				// happen and reads as a bug in Grit rather than a seed order
				// the project has not set up yet.
				pre.WriteString("\tif len(" + idsVar + ") == 0 {\n")
				pre.WriteString("\t\treturn fmt.Errorf(\"cannot seed " + g.Names().Lower +
					": no " + rel.Lower + " rows exist yet. Seed " + rel.Pascal +
					" first (generate it with --faker, or add rows in the admin), then run grit seed again\")\n")
				pre.WriteString("\t}\n")

				seenRel[idsVar] = true
			}
			if mode == "faker" {
				b.WriteString("\t\t\t" + fkGo + ": pickID(" + idsVar + "),\n")
			} else {
				b.WriteString("\t\t\t" + fkGo + ": firstID(" + idsVar + "),\n")
			}
			continue
		}
		goField := toPascalCase(f.Name)
		label := strings.Join(splitPascal(toPascalCase(f.Name)), " ")
		lower := strings.ToLower(f.Name)
		ft := FieldType(f.Type)
		faker := mode == "faker"
		var val string

		switch {
		case f.IsFile():
			needsFiles = true
			if faker {
				val = `&files.FileRef{URL: "https://picsum.photos/seed/" + gofakeit.UUID() + "/600/400", Name: "sample.jpg", MIME: "image/jpeg"}`
			} else {
				val = `&files.FileRef{URL: "https://picsum.photos/seed/` + toSnakeCase(f.Name) + `/600/400", Name: "sample.jpg", MIME: "image/jpeg"}`
			}
		case f.IsFiles():
			needsFiles = true
			if faker {
				val = `files.FileRefs{{URL: "https://picsum.photos/seed/" + gofakeit.UUID() + "/600/400", Name: "sample.jpg", MIME: "image/jpeg"}}`
			} else {
				val = `files.FileRefs{{URL: "https://picsum.photos/seed/` + toSnakeCase(f.Name) + `/600/400", Name: "sample.jpg", MIME: "image/jpeg"}}`
			}
		case ft == FieldInt:
			if faker {
				val = "gofakeit.Number(1, 100)"
			} else {
				val = "10"
			}
		case ft == FieldUint:
			if faker {
				val = "uint(gofakeit.Number(1, 100))"
			} else {
				val = "10"
			}
		case ft == FieldFloat:
			if faker {
				val = "gofakeit.Price(1, 1000)"
			} else {
				val = "9.99"
			}
		case ft == FieldBool:
			if faker {
				val = "gofakeit.Bool()"
			} else {
				val = "true"
			}
		case ft == FieldDate || ft == FieldDatetime:
			needsTime = true
			if faker {
				val = "gofakeit.Date()"
			} else {
				val = "time.Now()"
			}
		// A choice field has to be seeded from its own choices. Falling through
		// to the string branch below put gofakeit.Word() in a status column, so
		// a freshly seeded app opened on rows reading "moreover" and "ouch" —
		// values the form's own dropdown cannot offer and the API's validation
		// would reject. Anything rendering a status badge then has to cope with
		// a value that is not in the union its own generated type declares.
		case (ft == FieldSelect || ft == FieldRadio) && len(f.Options) > 0:
			values := make([]string, 0, len(f.Options))
			for _, o := range f.Options {
				values = append(values, strconv.Quote(o.Value))
			}
			if faker {
				val = "gofakeit.RandomString([]string{" + strings.Join(values, ", ") + "})"
			} else {
				val = values[0]
			}
		case ft == FieldText || ft == FieldRichtext:
			if faker {
				val = "gofakeit.Sentence(12)"
			} else {
				val = `"A sample ` + strings.ToLower(label) + `."`
			}
		default: // string
			if faker {
				switch {
				// A unique column comes first, before any of the nicer
				// generators below.
				//
				// gofakeit.Word() returns a real word from a finite list, so
				// forty rows on a unique column collide and the seeder logs a
				// constraint failure for the ones that lose. A SKU or a
				// reference number is exactly the field people mark unique, and
				// it is the one where a readable prefix plus entropy is more
				// useful than a plausible word anyway.
				case f.Unique:
					val = "strings.ToUpper(" + strconv.Quote(skuPrefix(f.Name)+"-") +
						" + gofakeit.LetterN(4) + gofakeit.DigitN(4))"
					needsStrings = true
				case strings.Contains(lower, "email"):
					val = "gofakeit.Email()"
				case strings.Contains(lower, "name"):
					val = "gofakeit.Name()"
				case strings.Contains(lower, "phone"):
					val = "gofakeit.Phone()"
				case strings.Contains(lower, "url"), strings.Contains(lower, "website"):
					val = "gofakeit.URL()"
				case strings.Contains(lower, "address"):
					val = "gofakeit.Street()"
				case strings.Contains(lower, "city"):
					val = "gofakeit.City()"
				case strings.Contains(lower, "country"):
					val = "gofakeit.Country()"
				case strings.Contains(lower, "company"):
					val = "gofakeit.Company()"
				case strings.Contains(lower, "color"), strings.Contains(lower, "colour"):
					val = "gofakeit.Color()"
				default:
					val = "gofakeit.Word()"
				}
			} else {
				if strings.Contains(lower, "email") {
					val = `"sample@example.com"`
				} else {
					val = `"Sample ` + label + `"`
				}
			}
		}
		b.WriteString("\t\t\t" + goField + ": " + val + ",\n")
	}
	return b.String(), pre.String(), needsTime, needsFiles, needsStrings
}

// skuPrefix builds a short readable prefix from a field name, so a faked
// unique value reads as SKU-AB1234 rather than as pure noise.
func skuPrefix(fieldName string) string {
	name := strings.ToUpper(toSnakeCase(fieldName))
	name = strings.ReplaceAll(name, "_", "")
	if len(name) > 3 {
		name = name[:3]
	}
	if name == "" {
		name = "REF"
	}
	return name
}
