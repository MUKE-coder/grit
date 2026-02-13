package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// AddRole adds a new role to all relevant files in a scaffolded Grit project.
func AddRole(roleName string) error {
	root, err := FindProjectRoot()
	if err != nil {
		return fmt.Errorf("not inside a Grit project: %w", err)
	}

	upper := strings.ToUpper(roleName)
	pascal := rolePascal(upper)
	label := roleLabel(upper)

	spinner := color.New(color.FgHiBlack)
	count := 0

	apiRoot := filepath.Join(root, "apps", "api")
	sharedRoot := filepath.Join(root, "packages", "shared")
	adminRoot := filepath.Join(root, "apps", "admin")

	// 1. Go model — add role constant
	modelFile := filepath.Join(apiRoot, "internal", "models", "user.go")
	if fileExists(modelFile) {
		data, _ := os.ReadFile(modelFile)
		if strings.Contains(string(data), fmt.Sprintf("Role%s", pascal)) {
			return fmt.Errorf("role %q already exists in %s", upper, "models/user.go")
		}
		constLine := fmt.Sprintf("\tRole%s = %q", pascal, upper)
		if err := injectBefore(modelFile, "// grit:roles", constLine); err != nil {
			spinner.Printf("  - Skipped models/user.go (marker not found)\n")
		} else {
			spinner.Printf("  ✓ Added Role%s constant to models/user.go\n", pascal)
			count++
		}
	}

	// 2. Zod schema — add to role enum
	schemaFile := filepath.Join(sharedRoot, "schemas", "user.ts")
	if fileExists(schemaFile) {
		if err := replaceOnMarkerLine(schemaFile, "// grit:role-enum",
			`])`, fmt.Sprintf(`, "%s"])`, upper)); err == nil {
			spinner.Printf("  ✓ Added %q to Zod schema enum\n", upper)
			count++
		} else {
			spinner.Printf("  - Skipped schemas/user.ts (%v)\n", err)
		}
	}

	// 3. TypeScript type — add to union
	typesFile := filepath.Join(sharedRoot, "types", "user.ts")
	if fileExists(typesFile) {
		if err := replaceOnMarkerLine(typesFile, "// grit:role-union",
			`";`, fmt.Sprintf(`" | "%s";`, upper)); err == nil {
			spinner.Printf("  ✓ Added %q to TypeScript union type\n", upper)
			count++
		} else {
			spinner.Printf("  - Skipped types/user.ts (%v)\n", err)
		}
	}

	// 4. Constants — add to ROLES object
	constantsFile := filepath.Join(sharedRoot, "constants", "index.ts")
	if fileExists(constantsFile) {
		entry := fmt.Sprintf("  %s: %q,", upper, upper)
		if err := injectBefore(constantsFile, "// grit:role-constants", entry); err != nil {
			spinner.Printf("  - Skipped constants/index.ts (marker not found)\n")
		} else {
			spinner.Printf("  ✓ Added %s to ROLES constants\n", upper)
			count++
		}
	}

	// 5. Admin users resource — badge config
	usersResource := filepath.Join(adminRoot, "resources", "users.ts")
	if fileExists(usersResource) {
		badgeLine := fmt.Sprintf(`          %s: { color: "success", label: "%s" },`, upper, label)
		if err := injectBefore(usersResource, "// grit:role-badges", badgeLine); err != nil {
			spinner.Printf("  - Skipped users resource badges (marker not found)\n")
		} else {
			spinner.Printf("  ✓ Added %s badge config\n", upper)
			count++
		}

		// 6. Filter options
		filterLine := fmt.Sprintf(`          { label: "%s", value: "%s" },`, label, upper)
		if err := injectBefore(usersResource, "// grit:role-filters", filterLine); err != nil {
			spinner.Printf("  - Skipped users resource filters (marker not found)\n")
		} else {
			spinner.Printf("  ✓ Added %s filter option\n", label)
			count++
		}

		// 7. Form select options
		optionLine := fmt.Sprintf(`          { label: "%s", value: "%s" },`, label, upper)
		if err := injectBefore(usersResource, "// grit:role-options", optionLine); err != nil {
			spinner.Printf("  - Skipped users resource form options (marker not found)\n")
		} else {
			spinner.Printf("  ✓ Added %s form option\n", label)
			count++
		}
	}

	if count == 0 {
		return fmt.Errorf("no files were updated — make sure role markers exist (run `grit upgrade` to add them)")
	}

	fmt.Println()
	green := color.New(color.FgHiGreen, color.Bold)
	green.Printf("  ✓ Role %q added to %d location(s)\n", upper, count)
	fmt.Println()

	gray := color.New(color.FgHiBlack)
	gray.Println("  Next steps:")
	gray.Printf("    1. Run `grit migrate` to update the database\n")
	gray.Printf("    2. Update sidebar visibility in components/layout/sidebar.tsx if needed\n")
	fmt.Println()

	return nil
}

// injectBefore finds a marker in a file and inserts code on the line before it.
func injectBefore(filePath, marker, code string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filePath, err)
	}

	content := string(data)
	idx := strings.Index(content, marker)
	if idx == -1 {
		return fmt.Errorf("marker %q not found in %s", marker, filePath)
	}

	lineStart := idx
	for lineStart > 0 && content[lineStart-1] != '\n' {
		lineStart--
	}

	newContent := content[:lineStart] + code + "\n" + content[lineStart:]
	return os.WriteFile(filePath, []byte(newContent), 0644)
}

// replaceOnMarkerLine finds the line containing the marker,
// then performs a string replacement on that line only.
func replaceOnMarkerLine(filePath, marker, old, replacement string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filePath, err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	found := false
	for i, line := range lines {
		if strings.Contains(line, marker) {
			if !strings.Contains(line, old) {
				return fmt.Errorf("pattern %q not found on marker line", old)
			}
			lines[i] = strings.Replace(line, old, replacement, 1)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("marker %q not found in %s", marker, filePath)
	}

	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
}

// rolePascal converts "CONTENT_MANAGER" to "ContentManager".
func rolePascal(upper string) string {
	words := strings.Split(strings.ToLower(upper), "_")
	var result string
	for _, w := range words {
		if len(w) > 0 {
			result += strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return result
}

// roleLabel converts "CONTENT_MANAGER" to "Content Manager".
func roleLabel(upper string) string {
	words := strings.Split(strings.ToLower(upper), "_")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
