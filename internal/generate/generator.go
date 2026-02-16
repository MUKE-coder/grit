package generate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

// Generator holds context for generating a resource.
type Generator struct {
	Root       string // project root (where docker-compose.yml lives)
	Module     string // Go module path (e.g., myapp/apps/api)
	Definition *ResourceDefinition
	Roles      []string // optional: restrict routes to these roles (e.g., ["ADMIN", "EDITOR"])
}

// Names holds all the naming variants for a resource.
type Names struct {
	Pascal      string // Post
	Camel       string // post
	Snake       string // post
	Kebab       string // post
	Lower       string // post
	Plural      string // posts
	PluralPascal string // Posts
	PluralSnake string // posts
	PluralKebab string // posts
}

// NewGenerator creates a generator after detecting the project root and module path.
func NewGenerator(def *ResourceDefinition) (*Generator, error) {
	root, err := findProjectRoot()
	if err != nil {
		return nil, err
	}

	module, err := readModulePath(root)
	if err != nil {
		return nil, err
	}

	return &Generator{
		Root:       root,
		Module:     module,
		Definition: def,
	}, nil
}

// Run generates all files for the resource and injects into existing files.
func (g *Generator) Run() error {
	names := g.Names()

	fmt.Printf("\n  Generating resource: %s\n\n", names.Pascal)

	// 1. Create new files
	if err := g.writeGoModel(names); err != nil {
		return fmt.Errorf("writing Go model: %w", err)
	}
	fmt.Printf("  ✓ apps/api/internal/models/%s.go\n", names.Snake)

	if err := g.writeGoService(names); err != nil {
		return fmt.Errorf("writing Go service: %w", err)
	}
	fmt.Printf("  ✓ apps/api/internal/services/%s.go\n", names.Snake)

	if err := g.writeGoHandler(names); err != nil {
		return fmt.Errorf("writing Go handler: %w", err)
	}
	fmt.Printf("  ✓ apps/api/internal/handlers/%s.go\n", names.Snake)

	if err := g.writeZodSchema(names); err != nil {
		return fmt.Errorf("writing Zod schema: %w", err)
	}
	fmt.Printf("  ✓ packages/shared/schemas/%s.ts\n", names.Kebab)

	if err := g.writeTSTypes(names); err != nil {
		return fmt.Errorf("writing TS types: %w", err)
	}
	fmt.Printf("  ✓ packages/shared/types/%s.ts\n", names.Kebab)

	// Write hooks for web (if it exists)
	webHooksDir := filepath.Join(g.Root, "apps", "web", "hooks")
	if dirExists(webHooksDir) {
		if err := g.writeReactQueryHooks(names, "web"); err != nil {
			return fmt.Errorf("writing web hooks: %w", err)
		}
		fmt.Printf("  ✓ apps/web/hooks/use-%s.ts\n", names.PluralKebab)
	}

	// Write admin resource definition + page (if admin app exists)
	adminResourcesDir := filepath.Join(g.Root, "apps", "admin", "resources")
	if dirExists(adminResourcesDir) {
		if err := g.writeResourceDefinition(names); err != nil {
			return fmt.Errorf("writing resource definition: %w", err)
		}
		fmt.Printf("  ✓ apps/admin/resources/%s.ts\n", names.PluralKebab)

		if err := g.writeResourcePage(names); err != nil {
			return fmt.Errorf("writing resource page: %w", err)
		}
		fmt.Printf("  ✓ apps/admin/app/(dashboard)/resources/%s/page.tsx\n", names.PluralKebab)
	}

	fmt.Println()

	// 2. Inject into existing files
	fmt.Println("  Injecting into existing files...")

	if err := g.injectAll(names); err != nil {
		return fmt.Errorf("injecting code: %w", err)
	}

	// Resolve new Go dependencies if needed (e.g., gorm.io/datatypes for string_array)
	needsDatatypes := false
	for _, f := range g.Definition.Fields {
		if f.NeedsDatatypesImport() {
			needsDatatypes = true
			break
		}
	}
	if needsDatatypes {
		apiDir := filepath.Join(g.Root, "apps", "api")
		cmd := exec.Command("go", "get", "gorm.io/datatypes")
		cmd.Dir = apiDir
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("adding gorm.io/datatypes dependency: %w\n%s", err, string(out))
		}
		fmt.Println("  ✓ Added gorm.io/datatypes dependency")
	}

	fmt.Println()
	fmt.Printf("  ✅ Resource %s generated successfully!\n\n", names.Pascal)
	fmt.Printf("  Next steps:\n")
	fmt.Printf("    1. cd apps/api && go build ./...\n")
	fmt.Printf("    2. Restart the API server\n")
	fmt.Printf("    3. The admin panel will show %s in the sidebar\n\n", names.PluralPascal)

	return nil
}

// Names builds all naming variants from the resource name.
func (g *Generator) Names() Names {
	raw := g.Definition.Name

	pascal := toPascalCase(raw)
	snake := toSnakeCase(pascal)
	kebab := strings.ReplaceAll(snake, "_", "-")
	camel := strings.ToLower(pascal[:1]) + pascal[1:]
	lower := strings.ToLower(pascal)

	plural := Pluralize(snake)
	pluralPascal := toPascalCase(plural)
	pluralKebab := strings.ReplaceAll(plural, "_", "-")

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

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting working directory: %w", err)
	}

	// Walk up looking for docker-compose.yml or turbo.json (monorepo markers)
	for {
		if fileExists(filepath.Join(dir, "docker-compose.yml")) ||
			fileExists(filepath.Join(dir, "turbo.json")) {
			return dir, nil
		}

		// Also detect API-only projects (go.mod in apps/api or just go.mod)
		if fileExists(filepath.Join(dir, "apps", "api", "go.mod")) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root (no docker-compose.yml or turbo.json found)")
		}
		dir = parent
	}
}

func readModulePath(root string) (string, error) {
	goModPath := filepath.Join(root, "apps", "api", "go.mod")
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("reading go.mod: %w (expected at %s)", err, goModPath)
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module "), nil
		}
	}

	return "", fmt.Errorf("no module directive found in %s", goModPath)
}

func toPascalCase(s string) string {
	// Handle snake_case, kebab-case, and already PascalCase
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	if len(parts) == 0 {
		return s
	}

	result := ""
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		result += strings.ToUpper(part[:1]) + part[1:]
	}
	return result
}

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func writeFileWithDirs(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}
	return os.WriteFile(path, []byte(content), 0644)
}
