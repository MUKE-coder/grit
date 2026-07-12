package generate

import (
	"fmt"
	"path/filepath"
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
	fieldLines, needsTime, needsFiles := g.seederFieldLines(mode)

	// Assemble imports — only what the record actually references.
	var imports strings.Builder
	imports.WriteString("\t\"log\"\n")
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
		if fileExists(filepath.Join(apiRoot, "internal", "models", names.Lower+".go")) {
			return definitionFromModelFile(apiRoot, name)
		}
	}
	return nil, fmt.Errorf("no model found for %q — generate the resource first (looked in apps/api/internal/models and internal/models)", names.Pascal)
}

// definitionFromModelFile parses a specific model file into a definition.
func definitionFromModelFile(apiRoot, name string) (*ResourceDefinition, error) {
	names := MakeNames(name)
	modelPath := filepath.Join(apiRoot, "internal", "models", names.Lower+".go")
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

	skip := map[string]bool{
		"id": true, "created_at": true, "updated_at": true,
		"deleted_at": true, "version": true, "slug": true,
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

// seederFieldLines returns the "GoField: value," lines for one record, plus
// whether the time and files packages are needed. Slug (auto), belongs_to
// (needs a real FK), m2m and string-array fields are skipped — the user wires
// those up by hand.
func (g *Generator) seederFieldLines(mode string) (lines string, needsTime, needsFiles bool) {
	var b strings.Builder
	for _, f := range g.Definition.Fields {
		if f.IsSlug() || f.IsBelongsTo() || f.IsManyToMany() || f.IsStringArray() {
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
		case ft == FieldText || ft == FieldRichtext:
			if faker {
				val = "gofakeit.Sentence(12)"
			} else {
				val = `"A sample ` + strings.ToLower(label) + `."`
			}
		default: // string
			if faker {
				switch {
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
	return b.String(), needsTime, needsFiles
}
