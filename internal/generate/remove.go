package generate

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RemoveResource removes a previously generated resource — deleting files and
// reversing all marker-based injections.
func RemoveResource(name string) error {
	root, err := findProjectRoot()
	if err != nil {
		return err
	}

	def := &ResourceDefinition{Name: name}
	gen := &Generator{Root: root, Definition: def}
	names := gen.Names()

	apiRoot := filepath.Join(root, "apps", "api")
	sharedRoot := filepath.Join(root, "packages", "shared")
	adminRoot := filepath.Join(root, "apps", "admin")
	webRoot := filepath.Join(root, "apps", "web")

	fmt.Printf("\n  Removing resource: %s\n\n", names.Pascal)

	// --- Delete generated files ---
	fmt.Println("  Deleting files...")

	filesToDelete := []string{
		filepath.Join(apiRoot, "internal", "models", names.Snake+".go"),
		filepath.Join(apiRoot, "internal", "services", names.Snake+".go"),
		filepath.Join(apiRoot, "internal", "handlers", names.Snake+".go"),
		filepath.Join(sharedRoot, "schemas", names.Kebab+".ts"),
		filepath.Join(sharedRoot, "types", names.Kebab+".ts"),
		filepath.Join(webRoot, "hooks", "use-"+names.PluralKebab+".ts"),
		filepath.Join(adminRoot, "hooks", "use-"+names.PluralKebab+".ts"),
		filepath.Join(adminRoot, "resources", names.PluralKebab+".ts"),
		filepath.Join(adminRoot, "app", "(dashboard)", "resources", names.PluralKebab, "page.tsx"),
	}

	for _, f := range filesToDelete {
		if _, err := os.Stat(f); err == nil {
			if err := os.Remove(f); err != nil {
				return fmt.Errorf("deleting %s: %w", f, err)
			}
			rel, _ := filepath.Rel(root, f)
			fmt.Printf("  ✗ %s\n", rel)
		}
	}

	// Remove empty page directory (e.g., apps/admin/app/(dashboard)/resources/posts/)
	pageDir := filepath.Join(adminRoot, "app", "(dashboard)", "resources", names.PluralKebab)
	if entries, err := os.ReadDir(pageDir); err == nil && len(entries) == 0 {
		os.Remove(pageDir)
	}

	fmt.Println()
	fmt.Println("  Cleaning injections...")

	// --- Reverse injections ---

	routesFile := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	modelFile := filepath.Join(apiRoot, "internal", "models", "user.go")
	schemaIndex := filepath.Join(sharedRoot, "schemas", "index.ts")
	typesIndex := filepath.Join(sharedRoot, "types", "index.ts")
	constantsIndex := filepath.Join(sharedRoot, "constants", "index.ts")
	registryFile := filepath.Join(adminRoot, "resources", "index.ts")

	// 1. Remove model from AutoMigrate
	if fileExists(modelFile) {
		if removeLinesContaining(modelFile, fmt.Sprintf("&%s{}", names.Pascal)) == nil {
			fmt.Println("  ✗ Removed model from AutoMigrate")
		}
	}

	// 2. Remove model from GORM Studio (inline)
	if fileExists(routesFile) {
		if removeInlineText(routesFile, fmt.Sprintf("&models.%s{}, ", names.Pascal)) == nil {
			fmt.Println("  ✗ Removed model from GORM Studio")
		}
	}

	// 3. Remove handler initialization block
	if fileExists(routesFile) {
		if removeLineBlock(routesFile,
			fmt.Sprintf("%sHandler := &handlers.%sHandler{", names.Camel, names.Pascal),
			"}") == nil {
			fmt.Println("  ✗ Removed handler initialization")
		}
	}

	// 4+5. Remove routes (protected + admin)
	if fileExists(routesFile) {
		removed := removeLinesContaining(routesFile, names.Camel+"Handler.")
		if removed == nil {
			fmt.Println("  ✗ Removed API routes")
		}
	}

	// 6. Remove role-restricted route group (if present)
	if fileExists(routesFile) {
		removeLineBlock(routesFile,
			fmt.Sprintf("// %s routes (restricted to", names.PluralPascal),
			"}")
		removeLineBlock(routesFile,
			fmt.Sprintf("%sGroup := protected.Group", names.Camel),
			"}")
	}

	// 7. Remove schema export
	if fileExists(schemaIndex) {
		if removeSchemaExportBlock(schemaIndex, names.Pascal, names.Kebab) == nil {
			fmt.Println("  ✗ Removed schema export")
		}
	}

	// 8. Remove type export
	if fileExists(typesIndex) {
		if removeLinesContaining(typesIndex, fmt.Sprintf(`from "./%s"`, names.Kebab)) == nil {
			fmt.Println("  ✗ Removed type export")
		}
	}

	// 9. Remove API route constants block
	if fileExists(constantsIndex) {
		upper := strings.ToUpper(names.Plural)
		if removeLineBlock(constantsIndex, fmt.Sprintf("  %s: {", upper), "},") == nil {
			fmt.Println("  ✗ Removed API route constants")
		}
	}

	// 10. Remove resource import from registry
	if fileExists(registryFile) {
		if removeLinesContaining(registryFile, fmt.Sprintf(`from "./%s"`, names.PluralKebab)) == nil {
			fmt.Println("  ✗ Removed resource import")
		}
	}

	// 11. Remove resource from registry list
	if fileExists(registryFile) {
		if removeLinesContaining(registryFile, fmt.Sprintf("%sResource,", names.Camel)) == nil {
			fmt.Println("  ✗ Removed resource from registry")
		}
	}

	fmt.Println()
	fmt.Printf("  ✅ Resource %s removed successfully!\n\n", names.Pascal)

	return nil
}

// removeLinesContaining removes all lines from a file that contain the given pattern.
func removeLinesContaining(filePath, pattern string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var result []string
	removed := false

	for _, line := range lines {
		if strings.Contains(line, pattern) {
			removed = true
			continue
		}
		result = append(result, line)
	}

	if !removed {
		return fmt.Errorf("pattern %q not found", pattern)
	}

	return os.WriteFile(filePath, []byte(strings.Join(result, "\n")), 0644)
}

// removeInlineText removes a specific text from any line in a file.
func removeInlineText(filePath, text string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	content := string(data)
	if !strings.Contains(content, text) {
		return fmt.Errorf("text %q not found", text)
	}

	newContent := strings.Replace(content, text, "", 1)
	return os.WriteFile(filePath, []byte(newContent), 0644)
}

// removeLineBlock removes a block of lines starting with startPattern and ending
// with the first line containing endPattern (inclusive).
func removeLineBlock(filePath, startPattern, endPattern string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var result []string
	inBlock := false
	removed := false

	for _, line := range lines {
		if !inBlock && strings.Contains(line, startPattern) {
			inBlock = true
			removed = true
			continue
		}
		if inBlock {
			if strings.Contains(strings.TrimSpace(line), endPattern) {
				inBlock = false
				continue
			}
			continue
		}
		result = append(result, line)
	}

	if !removed {
		return fmt.Errorf("block starting with %q not found", startPattern)
	}

	return os.WriteFile(filePath, []byte(strings.Join(result, "\n")), 0644)
}

// removeSchemaExportBlock removes a multi-line export block for a given schema
// from the schemas/index.ts file. Handles both single-line and multi-line exports.
func removeSchemaExportBlock(filePath, pascal, kebab string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var result []string
	inBlock := false
	removed := false
	fromPattern := fmt.Sprintf(`from "./%s"`, kebab)

	for i, line := range lines {
		// Single-line export: export { ... } from "./kebab";
		if !inBlock && strings.Contains(line, fromPattern) {
			removed = true
			continue
		}

		// Multi-line export: starts with "export {"
		if !inBlock && strings.Contains(line, "export {") {
			// Look ahead to see if this block references our schema
			isOurs := false
			for j := i; j < len(lines); j++ {
				if strings.Contains(lines[j], fromPattern) {
					isOurs = true
					break
				}
				if strings.Contains(lines[j], "from \"") && !strings.Contains(lines[j], fromPattern) {
					break
				}
			}
			if isOurs {
				inBlock = true
				removed = true
				continue
			}
		}

		if inBlock {
			if strings.Contains(line, fromPattern) {
				inBlock = false
				continue
			}
			continue
		}

		result = append(result, line)
	}

	if !removed {
		return fmt.Errorf("schema export for %q not found", pascal)
	}

	return os.WriteFile(filePath, []byte(strings.Join(result, "\n")), 0644)
}

// ConfirmRemoval prompts the user for confirmation.
func ConfirmRemoval() bool {
	fmt.Print("  Are you sure? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}
