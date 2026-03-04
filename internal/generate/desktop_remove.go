package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v2/internal/project"
)

// removeGoFuncBlock removes a Go function block starting with the given signature
// and ending at the first closing brace at column 0 (no indentation).
func removeGoFuncBlock(filePath, funcSig string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var result []string
	inBlock := false
	removed := false

	for _, line := range lines {
		if !inBlock && strings.Contains(line, funcSig) {
			inBlock = true
			removed = true
			continue
		}
		if inBlock {
			// Match closing brace at column 0 (no leading whitespace)
			if strings.TrimRight(line, " \t\r") == "}" {
				inBlock = false
				continue
			}
			continue
		}
		result = append(result, line)
	}

	if !removed {
		return fmt.Errorf("function %q not found", funcSig)
	}

	return os.WriteFile(filePath, []byte(strings.Join(result, "\n")), 0644)
}

// RemoveDesktopResource removes a previously generated desktop resource.
func RemoveDesktopResource(name string) error {
	info, err := project.DetectProject()
	if err != nil {
		return fmt.Errorf("detecting project: %w", err)
	}
	if info.Type != project.ProjectDesktop {
		return fmt.Errorf("not inside a desktop project")
	}

	def := &ResourceDefinition{Name: name}
	names := BuildNames(def)
	root := info.Root

	fmt.Printf("\n  Removing desktop resource: %s\n\n", names.Pascal)

	// --- Delete generated files ---
	fmt.Println("  Deleting files...")

	filesToDelete := []string{
		filepath.Join(root, "internal", "models", names.Snake+".go"),
		filepath.Join(root, "internal", "service", names.Snake+".go"),
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

	// Delete frontend page directory
	pageDir := filepath.Join(root, "frontend", "src", "pages", names.Plural)
	if _, err := os.Stat(pageDir); err == nil {
		if err := os.RemoveAll(pageDir); err != nil {
			return fmt.Errorf("deleting %s: %w", pageDir, err)
		}
		fmt.Printf("  ✗ frontend/src/pages/%s/\n", names.Plural)
	}

	fmt.Println()
	fmt.Println("  Cleaning injections...")

	// --- Reverse injections ---

	// 1. Remove model from AutoMigrate — db.go
	dbFile := filepath.Join(root, "internal", "db", "db.go")
	if fileExists(dbFile) {
		if removeLinesContaining(dbFile, fmt.Sprintf("&models.%s{}", names.Pascal)) == nil {
			fmt.Println("  ✗ Removed from db.go")
		}
	}

	// 2. Remove service init — main.go
	mainFile := filepath.Join(root, "main.go")
	if fileExists(mainFile) {
		if removeLinesContaining(mainFile, fmt.Sprintf("%sSvc := service.New%sService", names.Camel, names.Pascal)) == nil {
			fmt.Println("  ✗ Removed service init from main.go")
		}
	}

	// 3. Remove app arg — main.go (inline)
	if fileExists(mainFile) {
		removeInlineText(mainFile, fmt.Sprintf("%sSvc, ", names.Camel))
	}

	// 4. Remove app field — app.go
	appFile := filepath.Join(root, "app.go")
	if fileExists(appFile) {
		removeLinesContaining(appFile, fmt.Sprintf("%s *service.%sService", names.Camel, names.Pascal))
	}

	// 5. Remove constructor param — app.go (inline)
	if fileExists(appFile) {
		removeInlineText(appFile, fmt.Sprintf("%s *service.%sService, ", names.Camel, names.Pascal))
	}

	// 6. Remove constructor assign — app.go (inline)
	if fileExists(appFile) {
		removeInlineText(appFile, fmt.Sprintf("%s: %s, ", names.Camel, names.Camel))
	}

	// 7. Remove bound methods — app.go (remove each function block)
	if fileExists(appFile) {
		methodSigs := []string{
			fmt.Sprintf("func (a *App) Get%s(", names.PluralPascal),
			fmt.Sprintf("func (a *App) Get%s(", names.Pascal),
			fmt.Sprintf("func (a *App) Create%s(", names.Pascal),
			fmt.Sprintf("func (a *App) Update%s(", names.Pascal),
			fmt.Sprintf("func (a *App) Delete%s(", names.Pascal),
			fmt.Sprintf("func (a *App) Export%sPDF(", names.PluralPascal),
			fmt.Sprintf("func (a *App) Export%sExcel(", names.PluralPascal),
		}
		for _, sig := range methodSigs {
			removeGoFuncBlock(appFile, sig)
		}
		// Remove the section header
		removeLinesContaining(appFile, fmt.Sprintf("// ── %s ", names.PluralPascal))
	}

	// 8. Remove input type — types.go
	typesFile := filepath.Join(root, "internal", "models", "types.go")
	if fileExists(typesFile) {
		removeLineBlock(typesFile,
			fmt.Sprintf("type %sInput struct {", names.Pascal),
			"}")
	}

	// 9. Remove studio model — cmd/studio/main.go
	studioFile := filepath.Join(root, "cmd", "studio", "main.go")
	if fileExists(studioFile) {
		removeLinesContaining(studioFile, fmt.Sprintf("&models.%s{}", names.Pascal))
	}

	// 10. Remove page imports — App.tsx
	appTsx := filepath.Join(root, "frontend", "src", "App.tsx")
	if fileExists(appTsx) {
		removeLinesContaining(appTsx, fmt.Sprintf("pages/%s/", names.Plural))
	}

	// 11. Remove routes — App.tsx
	if fileExists(appTsx) {
		removeLinesContaining(appTsx, fmt.Sprintf("path=\"/%s\"", names.Plural))
		removeLinesContaining(appTsx, fmt.Sprintf("path=\"/%s/new\"", names.Plural))
		removeLinesContaining(appTsx, fmt.Sprintf("path=\"/%s/:id/edit\"", names.Plural))
	}

	// 12. Remove sidebar nav — sidebar.tsx
	sidebarFile := filepath.Join(root, "frontend", "src", "components", "layout", "sidebar.tsx")
	if fileExists(sidebarFile) {
		removeLinesContaining(sidebarFile, fmt.Sprintf("path: \"/%s\"", names.Plural))
	}

	fmt.Println()
	fmt.Printf("  ✅ Resource %s removed successfully!\n\n", names.Pascal)

	return nil
}
