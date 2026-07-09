package generate

import (
	"fmt"

	"github.com/MUKE-coder/grit/v3/internal/project"
)

// DesktopGenerator holds context for generating a resource in a desktop project.
type DesktopGenerator struct {
	Root       string
	Module     string
	Definition *ResourceDefinition
}

// RunDesktop is the entry point for desktop resource generation.
func RunDesktop(def *ResourceDefinition) error {
	info, err := project.DetectProject()
	if err != nil {
		return fmt.Errorf("detecting project: %w", err)
	}
	if info.Type != project.ProjectDesktop {
		return fmt.Errorf("not inside a desktop project")
	}

	g := &DesktopGenerator{
		Root:       info.Root,
		Module:     info.Module,
		Definition: def,
	}

	return g.Run()
}

// BuildNames creates naming variants from a resource definition (shared between web and desktop).
func BuildNames(def *ResourceDefinition) Names {
	raw := def.Name
	pascal := toPascalCase(raw)
	snake := toSnakeCase(pascal)
	kebab := replaceAll(snake, "_", "-")
	camel := toLowerFirst(pascal)
	lower := toLower(pascal)
	plural := Pluralize(snake)
	pluralPascal := toPascalCase(plural)
	pluralKebab := replaceAll(plural, "_", "-")

	return Names{
		Pascal:       pascal,
		Camel:        camel,
		Snake:        snake,
		Kebab:        kebab,
		Lower:        lower,
		Plural:       plural,
		PluralPascal: pluralPascal,
		PluralSnake:  plural,
		PluralKebab:  pluralKebab,
	}
}

// Run generates all files for a desktop resource and injects into existing files.
func (g *DesktopGenerator) Run() error {
	names := BuildNames(g.Definition)

	// many_to_many and string_array aren't wired for desktop resources yet
	// (m2m needs join-table handling + a multi-select UI; string_array needs
	// the datatypes column + an array editor). Fail early with an actionable
	// message instead of emitting a project that won't compile.
	for _, f := range g.Definition.Fields {
		if f.IsManyToMany() {
			return fmt.Errorf("field %q: many_to_many isn't supported for desktop apps yet — model the link as a separate resource with two belongs_to fields instead", f.Name)
		}
		if f.IsStringArray() {
			return fmt.Errorf("field %q: string_array isn't supported for desktop apps yet — use a text field, or a related resource", f.Name)
		}
	}

	fmt.Printf("\n  Generating desktop resource: %s\n\n", names.Pascal)

	// 1. Write Go model
	if err := g.writeDesktopModel(names); err != nil {
		return fmt.Errorf("writing model: %w", err)
	}
	fmt.Printf("  ✓ internal/models/%s.go\n", names.Snake)

	// 2. Write Go service
	if err := g.writeDesktopService(names); err != nil {
		return fmt.Errorf("writing service: %w", err)
	}
	fmt.Printf("  ✓ internal/service/%s.go\n", names.Snake)

	// 3. Write frontend list route
	if err := g.writeDesktopListRoute(names); err != nil {
		return fmt.Errorf("writing list route: %w", err)
	}
	fmt.Printf("  ✓ frontend/src/routes/_layout/%s.index.tsx\n", names.Plural)

	// 4. Write frontend new route
	if err := g.writeDesktopNewRoute(names); err != nil {
		return fmt.Errorf("writing new route: %w", err)
	}
	fmt.Printf("  ✓ frontend/src/routes/_layout/%s.new.tsx\n", names.Plural)

	// 5. Write frontend edit route
	if err := g.writeDesktopEditRoute(names); err != nil {
		return fmt.Errorf("writing edit route: %w", err)
	}
	fmt.Printf("  ✓ frontend/src/routes/_layout/%s.$id.edit.tsx\n", names.Plural)

	fmt.Println()
	fmt.Println("  Injecting into existing files...")

	// 5. Inject into existing files
	if err := g.injectDesktopAll(names); err != nil {
		return fmt.Errorf("injecting code: %w", err)
	}

	fmt.Println()
	fmt.Printf("  ✅ Resource %s generated successfully!\n\n", names.Pascal)
	fmt.Printf("  Next steps:\n")
	fmt.Printf("    1. Restart wails dev\n")
	fmt.Printf("    2. %s will appear in the sidebar\n\n", names.PluralPascal)

	return nil
}

// helpers to avoid importing strings in this file
func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}

func toLowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]+32) + s[1:]
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		result[i] = c
	}
	return string(result)
}
