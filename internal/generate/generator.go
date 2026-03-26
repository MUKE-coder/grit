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
	Root         string // project root
	Module       string // Go module path (e.g., myapp/apps/api or myapp)
	Architecture string // "single", "double", "triple", "api", "mobile"
	Frontend     string // "next", "tanstack"
	Definition   *ResourceDefinition
	Roles        []string // optional: restrict routes to these roles
}

// APIRoot returns the base directory for Go files.
func (g *Generator) APIRoot() string {
	if g.Architecture == "single" {
		return g.Root
	}
	return filepath.Join(g.Root, "apps", "api")
}

// AdminRoot returns the base directory for admin frontend files.
func (g *Generator) AdminRoot() string {
	return filepath.Join(g.Root, "apps", "admin")
}

// UseTanStack returns true if the project uses TanStack Router.
func (g *Generator) UseTanStack() bool {
	return g.Frontend == "tanstack"
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

	// Read grit.json for architecture and frontend
	arch, frontend := readGritJSON(root)

	module, err := readModulePath(root, arch)
	if err != nil {
		return nil, err
	}

	return &Generator{
		Root:         root,
		Module:       module,
		Architecture: arch,
		Frontend:     frontend,
		Definition:   def,
	}, nil
}

// Run generates all files for the resource and injects into existing files.
func (g *Generator) Run() error {
	names := g.Names()
	apiRoot := g.APIRoot()
	apiPrefix := "apps/api/" // for display
	if g.Architecture == "single" {
		apiPrefix = ""
	}

	fmt.Printf("\n  Generating resource: %s\n\n", names.Pascal)

	// 1. Create new Go files
	if err := g.writeGoModel(names); err != nil {
		return fmt.Errorf("writing Go model: %w", err)
	}
	fmt.Printf("  ✓ %sinternal/models/%s.go\n", apiPrefix, names.Snake)

	if err := g.writeGoService(names); err != nil {
		return fmt.Errorf("writing Go service: %w", err)
	}
	fmt.Printf("  ✓ %sinternal/services/%s.go\n", apiPrefix, names.Snake)

	if err := g.writeGoHandler(names); err != nil {
		return fmt.Errorf("writing Go handler: %w", err)
	}
	fmt.Printf("  ✓ %sinternal/handlers/%s.go\n", apiPrefix, names.Snake)

	// Shared types (monorepo only)
	sharedDir := filepath.Join(g.Root, "packages", "shared")
	if dirExists(sharedDir) {
		if err := g.writeZodSchema(names); err != nil {
			return fmt.Errorf("writing Zod schema: %w", err)
		}
		fmt.Printf("  ✓ packages/shared/schemas/%s.ts\n", names.Kebab)

		if err := g.writeTSTypes(names); err != nil {
			return fmt.Errorf("writing TS types: %w", err)
		}
		fmt.Printf("  ✓ packages/shared/types/%s.ts\n", names.Kebab)
	}

	// Write hooks for web app (monorepo: apps/web, check both hooks/ and src/hooks/)
	webHooksDir := filepath.Join(g.Root, "apps", "web", "hooks")
	webTanStackHooksDir := filepath.Join(g.Root, "apps", "web", "src", "hooks")
	if dirExists(webHooksDir) {
		if err := g.writeReactQueryHooks(names, "web"); err != nil {
			return fmt.Errorf("writing web hooks: %w", err)
		}
		fmt.Printf("  ✓ apps/web/hooks/use-%s.ts\n", names.PluralKebab)
	} else if dirExists(webTanStackHooksDir) {
		if err := g.writeReactQueryHooksTanStack(names, filepath.Join(g.Root, "apps", "web", "src", "hooks")); err != nil {
			return fmt.Errorf("writing web hooks: %w", err)
		}
		fmt.Printf("  ✓ apps/web/src/hooks/use-%s.ts\n", names.PluralKebab)
	}

	// Write hooks for single app frontend
	singleHooksDir := filepath.Join(g.Root, "frontend", "src", "hooks")
	if g.Architecture == "single" && dirExists(singleHooksDir) {
		if err := g.writeReactQueryHooksTanStack(names, singleHooksDir); err != nil {
			return fmt.Errorf("writing frontend hooks: %w", err)
		}
		fmt.Printf("  ✓ frontend/src/hooks/use-%s.ts\n", names.PluralKebab)
	}

	// Write admin resource definition + page (if admin app exists)
	adminResourcesDir := filepath.Join(g.Root, "apps", "admin", "resources")
	adminTanStackResourcesDir := filepath.Join(g.Root, "apps", "admin", "src", "resources")
	if dirExists(adminResourcesDir) {
		// Next.js admin
		if err := g.writeResourceDefinition(names); err != nil {
			return fmt.Errorf("writing resource definition: %w", err)
		}
		fmt.Printf("  ✓ apps/admin/resources/%s.ts\n", names.PluralKebab)

		if err := g.writeResourcePage(names); err != nil {
			return fmt.Errorf("writing resource page: %w", err)
		}
		fmt.Printf("  ✓ apps/admin/app/(dashboard)/resources/%s/page.tsx\n", names.PluralKebab)
	} else if dirExists(adminTanStackResourcesDir) {
		// TanStack admin
		if err := g.writeResourceDefinitionTanStack(names); err != nil {
			return fmt.Errorf("writing resource definition: %w", err)
		}
		fmt.Printf("  ✓ apps/admin/src/resources/%s.ts\n", names.PluralKebab)

		if err := g.writeResourcePageTanStack(names); err != nil {
			return fmt.Errorf("writing resource page: %w", err)
		}
		fmt.Printf("  ✓ apps/admin/src/routes/_dashboard/resources/%s.tsx\n", names.PluralKebab)
	}

	fmt.Println()

	// 2. Inject into existing files
	fmt.Println("  Injecting into existing files...")

	if err := g.injectAll(names); err != nil {
		return fmt.Errorf("injecting code: %w", err)
	}

	// Resolve new Go dependencies if needed
	needsDatatypes := false
	for _, f := range g.Definition.Fields {
		if f.NeedsDatatypesImport() {
			needsDatatypes = true
			break
		}
	}
	if needsDatatypes {
		cmd := exec.Command("go", "get", "gorm.io/datatypes")
		cmd.Dir = apiRoot
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("adding gorm.io/datatypes dependency: %w\n%s", err, string(out))
		}
		fmt.Println("  ✓ Added gorm.io/datatypes dependency")
	}

	fmt.Println()
	fmt.Printf("  ✅ Resource %s generated successfully!\n\n", names.Pascal)
	fmt.Printf("  Next steps:\n")
	if g.Architecture == "single" {
		fmt.Printf("    1. go build ./...\n")
		fmt.Printf("    2. Restart the server\n")
	} else {
		fmt.Printf("    1. cd apps/api && go build ./...\n")
		fmt.Printf("    2. Restart the API server\n")
	}
	if g.Architecture == "triple" {
		fmt.Printf("    3. The admin panel will show %s in the sidebar\n", names.PluralPascal)
	}
	fmt.Println()

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

	for {
		// grit.json is the canonical project marker
		if fileExists(filepath.Join(dir, "grit.json")) {
			return dir, nil
		}
		// Fallback: monorepo markers
		if fileExists(filepath.Join(dir, "docker-compose.yml")) ||
			fileExists(filepath.Join(dir, "turbo.json")) {
			return dir, nil
		}
		// Fallback: single app (go.mod + internal/ at root)
		if fileExists(filepath.Join(dir, "go.mod")) && dirExists(filepath.Join(dir, "internal")) {
			return dir, nil
		}
		// Fallback: monorepo with apps/api
		if fileExists(filepath.Join(dir, "apps", "api", "go.mod")) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root (no grit.json, docker-compose.yml, or turbo.json found)")
		}
		dir = parent
	}
}

// readGritJSON reads the grit.json manifest for architecture and frontend.
// Returns defaults ("triple", "next") if file doesn't exist.
func readGritJSON(root string) (string, string) {
	data, err := os.ReadFile(filepath.Join(root, "grit.json"))
	if err != nil {
		// Fallback: detect architecture from directory structure
		if dirExists(filepath.Join(root, "frontend")) && !dirExists(filepath.Join(root, "apps")) {
			return "single", "tanstack"
		}
		return "triple", "next"
	}

	// Simple JSON parsing without encoding/json to avoid import bloat
	content := string(data)
	arch := extractJSONValue(content, "architecture")
	frontend := extractJSONValue(content, "frontend")

	if arch == "" {
		arch = "triple"
	}
	if frontend == "" {
		frontend = "next"
	}

	return arch, frontend
}

// extractJSONValue does simple string extraction from JSON (avoids encoding/json import).
func extractJSONValue(json, key string) string {
	search := fmt.Sprintf(`"%s"`, key)
	idx := strings.Index(json, search)
	if idx < 0 {
		return ""
	}

	// Find the value after the colon
	rest := json[idx+len(search):]
	colonIdx := strings.Index(rest, ":")
	if colonIdx < 0 {
		return ""
	}
	rest = rest[colonIdx+1:]

	// Find the quoted value
	startQuote := strings.Index(rest, `"`)
	if startQuote < 0 {
		return ""
	}
	rest = rest[startQuote+1:]
	endQuote := strings.Index(rest, `"`)
	if endQuote < 0 {
		return ""
	}
	return rest[:endQuote]
}

func readModulePath(root, arch string) (string, error) {
	// Single app: go.mod at project root
	// Monorepo: go.mod at apps/api/
	var goModPath string
	if arch == "single" {
		goModPath = filepath.Join(root, "go.mod")
	} else {
		goModPath = filepath.Join(root, "apps", "api", "go.mod")
	}

	data, err := os.ReadFile(goModPath)
	if err != nil {
		// Try the other location as fallback
		alt := filepath.Join(root, "go.mod")
		if arch != "single" {
			alt = filepath.Join(root, "apps", "api", "go.mod")
		}
		data, err = os.ReadFile(alt)
		if err != nil {
			return "", fmt.Errorf("reading go.mod: %w (checked %s)", err, goModPath)
		}
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module "), nil
		}
	}

	return "", fmt.Errorf("no module directive found in go.mod")
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
