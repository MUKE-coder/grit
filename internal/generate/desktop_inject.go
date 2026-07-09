package generate

import (
	"fmt"
	"path/filepath"
	"strings"
)

// injectDesktopAll injects code into all existing desktop project files that have markers.
func (g *DesktopGenerator) injectDesktopAll(names Names) error {

	// 1. Inject model into AutoMigrate — db.go
	dbFile := filepath.Join(g.Root, "internal", "db", "db.go")
	if fileExists(dbFile) {
		if err := injectBefore(dbFile, "// grit:models",
			fmt.Sprintf("\t\t&models.%s{},", names.Pascal)); err != nil {
			return fmt.Errorf("injecting model into db.go: %w", err)
		}
		fmt.Println("  ✓ Injected into db.go")
	}

	// 2. Inject service init — main.go
	mainFile := filepath.Join(g.Root, "main.go")
	if fileExists(mainFile) {
		if err := injectBefore(mainFile, "// grit:service-init",
			fmt.Sprintf("\t%sSvc := service.New%sService(database)", names.Camel, names.Pascal)); err != nil {
			return fmt.Errorf("injecting service init into main.go: %w", err)
		}
		fmt.Println("  ✓ Injected into main.go (service init)")
	}

	// 3. Inject app args — main.go (inline)
	if fileExists(mainFile) {
		if err := injectInline(mainFile, "/* grit:app-args */",
			fmt.Sprintf("%sSvc, ", names.Camel)); err != nil {
			return fmt.Errorf("injecting app args into main.go: %w", err)
		}
		fmt.Println("  ✓ Injected into main.go (app args)")
	}

	// 4. Inject app field — app.go
	appFile := filepath.Join(g.Root, "app.go")
	if fileExists(appFile) {
		if err := injectBefore(appFile, "// grit:fields",
			fmt.Sprintf("\t%s *service.%sService", names.Camel, names.Pascal)); err != nil {
			return fmt.Errorf("injecting field into app.go: %w", err)
		}
		fmt.Println("  ✓ Injected into app.go (field)")
	}

	// 5. Inject constructor param — app.go (inline)
	if fileExists(appFile) {
		if err := injectInline(appFile, "/* grit:constructor-params */",
			fmt.Sprintf("%s *service.%sService, ", names.Camel, names.Pascal)); err != nil {
			return fmt.Errorf("injecting constructor param into app.go: %w", err)
		}
		fmt.Println("  ✓ Injected into app.go (constructor param)")
	}

	// 6. Inject constructor assign — app.go (inline)
	if fileExists(appFile) {
		if err := injectInline(appFile, "/* grit:constructor-assign */",
			fmt.Sprintf("%s: %s, ", names.Camel, names.Camel)); err != nil {
			return fmt.Errorf("injecting constructor assign into app.go: %w", err)
		}
		fmt.Println("  ✓ Injected into app.go (constructor assign)")
	}

	// 7. Inject bound methods — app.go
	if fileExists(appFile) {
		methods := g.buildBoundMethods(names)
		if err := injectBefore(appFile, "// grit:methods", methods); err != nil {
			return fmt.Errorf("injecting methods into app.go: %w", err)
		}
		fmt.Println("  ✓ Injected into app.go (methods)")
	}

	// 8. Inject input type — types.go
	typesFile := filepath.Join(g.Root, "internal", "models", "types.go")
	if fileExists(typesFile) {
		// file/files inputs are files.FileRef(s) — make sure the package is
		// imported. Idempotent: injectBefore skips if it's already there.
		if g.definitionHasFileField() {
			imp := "\t\"" + g.Module + "/internal/files\""
			if err := injectBefore(typesFile, "// grit:input-imports", imp); err != nil {
				return fmt.Errorf("injecting files import into types.go: %w", err)
			}
		}
		inputStruct := g.buildInputStruct(names)
		if err := injectBefore(typesFile, "// grit:input-types", inputStruct); err != nil {
			return fmt.Errorf("injecting input type into types.go: %w", err)
		}
		fmt.Println("  ✓ Injected into types.go")
	}

	// 9. Inject studio model — cmd/studio/main.go
	studioFile := filepath.Join(g.Root, "cmd", "studio", "main.go")
	if fileExists(studioFile) {
		if err := injectBefore(studioFile, "// grit:studio-models",
			fmt.Sprintf("\t\t&models.%s{},", names.Pascal)); err != nil {
			return fmt.Errorf("injecting studio model: %w", err)
		}
		fmt.Println("  ✓ Injected into cmd/studio/main.go")
	}

	// 10. Inject sidebar — sidebar.tsx
	sidebarFile := filepath.Join(g.Root, "frontend", "src", "components", "layout", "sidebar.tsx")
	if fileExists(sidebarFile) {
		iconName := guessLucideIcon(names.Pascal)

		// 12a. Inject icon import
		if err := injectBefore(sidebarFile, "// grit:nav-icons",
			fmt.Sprintf("  %s,", iconName)); err != nil {
			return fmt.Errorf("injecting nav icon into sidebar.tsx: %w", err)
		}

		// 12b. Inject nav item
		navItem := fmt.Sprintf("    { label: \"%s\", path: \"/%s\", icon: %s },",
			names.PluralPascal, names.Plural, iconName)
		if err := injectBefore(sidebarFile, "// grit:nav", navItem); err != nil {
			return fmt.Errorf("injecting nav item into sidebar.tsx: %w", err)
		}

		fmt.Println("  ✓ Injected into sidebar.tsx")
	}

	return nil
}

// definitionHasFileField reports whether any field is a file or files type, in
// which case models and the input type reference the files package.
func (g *DesktopGenerator) definitionHasFileField() bool {
	for _, f := range g.Definition.Fields {
		if f.IsFile() || f.IsFiles() {
			return true
		}
	}
	return false
}

// buildBoundMethods generates the Wails-bound CRUD + export methods for app.go.
func (g *DesktopGenerator) buildBoundMethods(names Names) string {
	var b strings.Builder

	// Section header
	b.WriteString(fmt.Sprintf("// ── %s ──────────────────────────────\n\n", names.PluralPascal))

	// List
	b.WriteString(fmt.Sprintf("func (a *App) Get%s(page, pageSize int, search string) (*models.PaginatedResult, error) {\n", names.PluralPascal))
	b.WriteString(fmt.Sprintf("\treturn a.%s.List(page, pageSize, search)\n", names.Camel))
	b.WriteString("}\n\n")

	// GetByID
	b.WriteString(fmt.Sprintf("func (a *App) Get%s(id uint) (*models.%s, error) {\n", names.Pascal, names.Pascal))
	b.WriteString(fmt.Sprintf("\treturn a.%s.GetByID(id)\n", names.Camel))
	b.WriteString("}\n\n")

	// Create
	b.WriteString(fmt.Sprintf("func (a *App) Create%s(input models.%sInput) (*models.%s, error) {\n", names.Pascal, names.Pascal, names.Pascal))
	b.WriteString(fmt.Sprintf("\treturn a.%s.Create(input)\n", names.Camel))
	b.WriteString("}\n\n")

	// Update
	b.WriteString(fmt.Sprintf("func (a *App) Update%s(id uint, input models.%sInput) (*models.%s, error) {\n", names.Pascal, names.Pascal, names.Pascal))
	b.WriteString(fmt.Sprintf("\treturn a.%s.Update(id, input)\n", names.Camel))
	b.WriteString("}\n\n")

	// Delete
	b.WriteString(fmt.Sprintf("func (a *App) Delete%s(id uint) error {\n", names.Pascal))
	b.WriteString(fmt.Sprintf("\treturn a.%s.Delete(id)\n", names.Camel))
	b.WriteString("}\n\n")

	// ExportPDF
	headers, accessors := g.buildExportFieldsInfo()
	b.WriteString(fmt.Sprintf("func (a *App) Export%sPDF() (string, error) {\n", names.PluralPascal))
	b.WriteString(fmt.Sprintf("\titems, err := a.%s.ListAll()\n", names.Camel))
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\treturn \"\", err\n")
	b.WriteString("\t}\n")
	b.WriteString(fmt.Sprintf("\theaders := []string{%s}\n", headers))
	b.WriteString("\tvar rows [][]string\n")
	b.WriteString("\tfor _, item := range items {\n")
	b.WriteString(fmt.Sprintf("\t\trows = append(rows, []string{%s})\n", accessors))
	b.WriteString("\t}\n")
	b.WriteString(fmt.Sprintf("\treturn a.export.ExportPDF(\"%s\", headers, rows)\n", names.PluralPascal))
	b.WriteString("}\n\n")

	// ExportExcel
	b.WriteString(fmt.Sprintf("func (a *App) Export%sExcel() (string, error) {\n", names.PluralPascal))
	b.WriteString(fmt.Sprintf("\titems, err := a.%s.ListAll()\n", names.Camel))
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\treturn \"\", err\n")
	b.WriteString("\t}\n")
	b.WriteString(fmt.Sprintf("\theaders := []string{%s}\n", headers))
	b.WriteString("\tvar rows [][]string\n")
	b.WriteString("\tfor _, item := range items {\n")
	b.WriteString(fmt.Sprintf("\t\trows = append(rows, []string{%s})\n", accessors))
	b.WriteString("\t}\n")
	b.WriteString(fmt.Sprintf("\treturn a.export.ExportExcel(\"%s\", headers, rows)\n", names.PluralPascal))
	b.WriteString("}")

	return b.String()
}

// buildExportFieldsInfo returns the headers string and accessors string for export methods.
func (g *DesktopGenerator) buildExportFieldsInfo() (string, string) {
	var headers []string
	var accessors []string

	for _, f := range g.Definition.Fields {
		// Skip fields that don't render as a single CSV/PDF cell: slug (internal),
		// m2m / string arrays (multi-valued), and file/files (a FileRef object,
		// not a string — exporting an image URL to a spreadsheet isn't useful).
		if f.IsSlug() || f.IsManyToMany() || f.IsStringArray() || f.IsFile() || f.IsFiles() {
			continue
		}

		pascalField := toPascalCase(f.Name)
		label := toTitleCase(f.Name)
		headers = append(headers, fmt.Sprintf("%q", label))

		switch {
		case FieldType(f.Type) == FieldBool:
			accessors = append(accessors, fmt.Sprintf("fmt.Sprintf(\"%%v\", item.%s)", pascalField))
		case FieldType(f.Type) == FieldDatetime || FieldType(f.Type) == FieldDate:
			accessors = append(accessors, fmt.Sprintf("item.%s.Format(\"2006-01-02\")", pascalField))
		case FieldType(f.Type) == FieldInt || FieldType(f.Type) == FieldUint || FieldType(f.Type) == FieldFloat || FieldType(f.Type) == FieldBelongsTo:
			accessors = append(accessors, fmt.Sprintf("fmt.Sprintf(\"%%v\", item.%s)", pascalField))
		default:
			// string, text, richtext
			accessors = append(accessors, fmt.Sprintf("item.%s", pascalField))
		}
	}

	return strings.Join(headers, ", "), strings.Join(accessors, ", ")
}

// buildInputStruct generates the Input struct for types.go.
func (g *DesktopGenerator) buildInputStruct(names Names) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("type %sInput struct {\n", names.Pascal))

	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}

		pascalField := toPascalCase(f.Name)
		snakeField := toSnakeCase(toPascalCase(f.Name))
		goType := f.GoType()

		b.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n", pascalField, goType, snakeField))
	}

	b.WriteString("}\n")

	return b.String()
}

// toTitleCase converts a snake_case or camelCase field name to Title Case.
func toTitleCase(s string) string {
	pascal := toPascalCase(s)
	var words []string
	start := 0
	for i := 1; i < len(pascal); i++ {
		if pascal[i] >= 'A' && pascal[i] <= 'Z' {
			words = append(words, pascal[start:i])
			start = i
		}
	}
	words = append(words, pascal[start:])
	return strings.Join(words, " ")
}
