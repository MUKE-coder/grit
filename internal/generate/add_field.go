package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// addFieldSupported lists the field types `grit g field` can add to an existing
// resource without rewriting imports or relationship wiring. Relationship, file,
// slug, and JSON-array types are intentionally excluded — those change imports,
// join tables, or migrations, so regenerating the resource is the safe path.
var addFieldSupported = map[FieldType]bool{
	FieldString: true, FieldText: true, FieldRichtext: true,
	FieldInt: true, FieldUint: true, FieldFloat: true,
	FieldBool: true, FieldToggle: true, FieldSelect: true,
	FieldDate: true, FieldDatetime: true,
}

// AddField adds a single field to an already-generated resource, in place: it
// injects the column into the Go model, the create/update Zod schemas, the
// TypeScript interface, and the admin resource's form fields and table columns.
// The database column itself is added by GORM on the next `grit migrate` (the
// model is the source of truth), so no migration file is written.
func AddField(resourceName, spec string) error {
	def, err := ParseInlineFields(resourceName, spec)
	if err != nil {
		return err
	}
	if len(def.Fields) == 0 {
		return fmt.Errorf("no field parsed from %q", spec)
	}
	f := def.Fields[len(def.Fields)-1]

	if !addFieldSupported[FieldType(f.Type)] {
		return fmt.Errorf("grit g field can't add a %q field yet: regenerate the resource "+
			"for relationship, file, slug, or array fields", f.Type)
	}

	// DefinitionFromModel confirms the resource exists and gives us the project
	// layout (single-binary vs monorepo, Next vs TanStack admin).
	existing, err := DefinitionFromModel(resourceName)
	if err != nil {
		return fmt.Errorf("resource %q not found (generate it first): %w", resourceName, err)
	}
	g, err := NewGenerator(existing)
	if err != nil {
		return err
	}
	names := BuildNames(existing)

	if err := g.injectModelField(names, f); err != nil {
		return err
	}
	if err := g.injectZodField(names, f); err != nil {
		return err
	}
	if err := g.injectTSField(names, f); err != nil {
		return err
	}
	if err := g.injectAdminField(names, f); err != nil {
		return err
	}
	return nil
}

// injectAfterAnchor inserts code on the line after the first line equal (trimmed)
// to anchor. Unlike a marker, the anchor is an existing structural line — the
// opening of a struct, schema, or interface — so this works on projects generated
// before the command existed.
func injectAfterAnchor(filePath, anchor, code string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filePath, err)
	}
	content := string(data)
	if alreadyInjected(content, strings.TrimRight(code, "\n")) {
		return nil // idempotent — re-running doesn't duplicate the field
	}
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == anchor {
			out := append([]string{}, lines[:i+1]...)
			out = append(out, strings.TrimRight(code, "\n"))
			out = append(out, lines[i+1:]...)
			return os.WriteFile(filePath, []byte(strings.Join(out, "\n")), 0644)
		}
	}
	return fmt.Errorf("anchor %q not found in %s", anchor, filePath)
}

func (g *Generator) injectModelField(names Names, f Field) error {
	path := filepath.Join(g.APIRoot(), "internal", "models", names.Snake+".go")
	goName := toPascalCase(f.Name)
	tags := fmt.Sprintf(`json:"%s"`, toSnakeCase(f.Name))
	if gt := f.GORMTag(); gt != "" {
		tags = fmt.Sprintf(`gorm:"%s" %s`, gt, tags)
	}
	if f.Required && f.GoType() == "string" {
		tags += ` binding:"required"`
	}
	line := fmt.Sprintf("\t%s %s `%s`", goName, f.GoType(), tags)
	return injectAfterAnchor(path, fmt.Sprintf("type %s struct {", names.Pascal), line)
}

func (g *Generator) injectZodField(names Names, f Field) error {
	path := filepath.Join(g.Root, "packages", "shared", "schemas", names.Kebab+".ts")
	snake := toSnakeCase(f.Name)
	if err := injectAfterAnchor(path,
		fmt.Sprintf("export const Create%sSchema = z.object({", names.Pascal),
		fmt.Sprintf("  %s: %s,", snake, f.ZodType())); err != nil {
		return err
	}
	updateZod := f.ZodType()
	if !strings.Contains(updateZod, ".optional()") && !strings.Contains(updateZod, ".nullable()") {
		updateZod += ".optional()"
	}
	return injectAfterAnchor(path,
		fmt.Sprintf("export const Update%sSchema = z.object({", names.Pascal),
		fmt.Sprintf("  %s: %s,", snake, updateZod))
}

func (g *Generator) injectTSField(names Names, f Field) error {
	path := filepath.Join(g.Root, "packages", "shared", "types", names.Kebab+".ts")
	return injectAfterAnchor(path,
		fmt.Sprintf("export interface %s {", names.Pascal),
		fmt.Sprintf("  %s: %s;", toSnakeCase(f.Name), f.TSType()))
}

// injectAdminField adds the form field and table column at the resource
// definition's auto markers. The resource file lives at apps/admin/resources for
// the Next admin and apps/admin/src/resources for the TanStack admin — inject
// into whichever exists.
func (g *Generator) injectAdminField(names Names, f Field) error {
	candidates := []string{
		filepath.Join(g.AdminRoot(), "resources", names.PluralKebab+".ts"),
		filepath.Join(g.AdminRoot(), "src", "resources", names.PluralKebab+".ts"),
	}
	label := strings.Join(splitPascal(toPascalCase(f.Name)), " ")
	snake := toSnakeCase(f.Name)

	formParts := []string{
		fmt.Sprintf(`key: "%s"`, snake),
		fmt.Sprintf(`label: "%s"`, label),
		fmt.Sprintf(`type: "%s"`, f.FormFieldType()),
	}
	if f.Required {
		formParts = append(formParts, "required: true")
	}
	if f.HasOptions() {
		formParts = append(formParts, "options: "+f.OptionsLiteral())
	}
	formLine := "      { " + strings.Join(formParts, ", ") + " },"
	colLine := fmt.Sprintf(`      { key: "%s", label: "%s" },`, snake, label)

	injected := false
	for _, path := range candidates {
		if _, err := os.Stat(path); err != nil {
			continue
		}
		if err := injectBefore(path, "// grit:fields:auto-end", formLine); err != nil {
			return err
		}
		if err := injectBefore(path, "// grit:cols:auto-end", colLine); err != nil {
			return err
		}
		injected = true
	}
	if !injected {
		return fmt.Errorf("admin resource definition for %s not found", names.PluralKebab)
	}
	return nil
}
