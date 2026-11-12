package scaffold

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// Taking a role out of a project.
//
// `grit add role` existed and had no counterpart, so a role added by mistake,
// or one a project outgrew, had to be unpicked from seven places by hand. Miss
// one and the result is quiet: a constant nothing uses, or worse, a value in
// the Zod enum that the Go model will not accept, so the form offers a role
// the API refuses.
//
// # What it refuses to do
//
// ADMIN and USER stay. The framework's own code names them: the seeder makes
// an ADMIN, every generated route guard asks for ADMIN, and USER is what a
// registration assigns. Removing either leaves a project that compiles and
// cannot sign anybody in, which is not a thing to let somebody do by typing
// one command.
//
// # What it cannot check
//
// Whether any account still holds the role. That needs the database, and this
// command reads files, like every other one here: the database may be
// Postgres, it may not be running, and "could not connect" would have become a
// silent skip of the only safety check there was. So it says what to run
// instead, every time, rather than pretending to know.
//
// It matters because those rows keep a value nothing recognises afterwards:
// the API refuses a role it cannot resolve, so the accounts stop being able to
// sign in, and the admin renders a blank badge where the role used to be.

// protectedRoles are the two the framework itself depends on.
var protectedRoles = map[string]string{
	"ADMIN": "the seeder creates one, and every generated route guard asks for it",
	"USER":  "registration assigns it, so without it a new account has no role at all",
}

// RemoveRole takes a role out of every file `grit add role` puts it in.
func RemoveRole(roleName string, force bool) error {
	root, err := FindProjectRoot()
	if err != nil {
		return fmt.Errorf("not inside a Grit project: %w", err)
	}

	upper := strings.ToUpper(strings.TrimSpace(roleName))
	if upper == "" {
		return fmt.Errorf("name the role to remove")
	}
	if why, protected := protectedRoles[upper]; protected {
		return fmt.Errorf("%s cannot be removed: %s", upper, why)
	}

	pascal := rolePascal(upper)
	label := roleLabel(upper)
	paths := roleFilesFor(root)

	if paths.Model == "" {
		return fmt.Errorf("no internal/models/user.go found, so there is nothing to remove the role from")
	}

	// Is it even there? Removing a role a project does not have should say
	// so rather than report a successful no-op.
	model, err := os.ReadFile(paths.Model)
	if err != nil {
		return fmt.Errorf("reading the user model: %w", err)
	}
	if !strings.Contains(string(model), fmt.Sprintf("Role%s ", pascal)) &&
		!strings.Contains(string(model), fmt.Sprintf("Role%s=", pascal)) {
		return fmt.Errorf("this project has no %s role", upper)
	}

	// Code that still names the role. A guard asking for it keeps compiling
	// after the constant is gone only if it spells the role as a string, and
	// then it asks for something nothing can be.
	stillNamed := filesNaming(root, upper)
	if len(stillNamed) > 0 && !force {
		return fmt.Errorf(
			"%s is still named in %d file(s):\n    %s\n"+
				"Change those first, or pass --force to remove the declarations anyway",
			upper, len(stillNamed), strings.Join(stillNamed, "\n    "))
	}

	spinner := color.New(color.FgHiBlack)
	removed := 0

	// 1. The Go constant.
	if n := removeLinesMatching(paths.Model, regexp.MustCompile(
		`^\s*Role`+regexp.QuoteMeta(pascal)+`\s*=\s*"`+regexp.QuoteMeta(upper)+`"\s*$`)); n > 0 {
		spinner.Printf("  ✓ Removed Role%s from %s\n", pascal, shown(root, paths.Model))
		removed += n
	}

	// 2. The Zod enum, which is one line holding every role.
	if paths.Schema != "" {
		if removeFromList(paths.Schema, fmt.Sprintf(`"%s"`, upper)) {
			spinner.Printf("  ✓ Removed %q from the Zod enum\n", upper)
			removed++
		}
	}

	// 3. The TypeScript union, likewise.
	if paths.Types != "" {
		if removeFromUnion(paths.Types, upper) {
			spinner.Printf("  ✓ Removed %q from the TypeScript union\n", upper)
			removed++
		}
	}

	// 4. The shared constant.
	if paths.Constants != "" {
		if n := removeLinesMatching(paths.Constants, regexp.MustCompile(
			`^\s*`+regexp.QuoteMeta(upper)+`\s*:\s*"`+regexp.QuoteMeta(upper)+`"\s*,?\s*$`)); n > 0 {
			spinner.Printf("  ✓ Removed %s from ROLES\n", upper)
			removed += n
		}
	}

	// 5, 6, 7. The admin's badge, filter and form option.
	if paths.Admin != "" {
		patterns := []struct {
			re   *regexp.Regexp
			what string
		}{
			{regexp.MustCompile(`^\s*` + regexp.QuoteMeta(upper) + `:\s*\{[^}]*\},?\s*$`), "badge"},
			{regexp.MustCompile(`^\s*\{\s*label:\s*"` + regexp.QuoteMeta(label) +
				`",\s*value:\s*"` + regexp.QuoteMeta(upper) + `"\s*\},?\s*$`), "filter and form option"},
		}
		for _, p := range patterns {
			if n := removeLinesMatching(paths.Admin, p.re); n > 0 {
				spinner.Printf("  ✓ Removed the %s %s from %s\n", upper, p.what, shown(root, paths.Admin))
				removed += n
			}
		}
	}

	if removed == 0 {
		return fmt.Errorf("nothing was removed: the role's lines are not in the shape Grit wrote them")
	}

	for _, what := range paths.missing(root) {
		spinner.Printf("  - No %s in this project\n", what)
	}

	fmt.Println()
	green := color.New(color.FgHiGreen, color.Bold)
	green.Printf("  ✓ Role %q removed from %d place(s)\n", upper, removed)
	fmt.Println()

	gray := color.New(color.FgHiBlack)
	gray.Println("  Accounts that still hold it:")
	gray.Printf("    This reads files, not your database, so it cannot tell you.\n")
	gray.Printf("    Any account left on %s cannot sign in: the API refuses a role\n", upper)
	gray.Printf("    it cannot resolve. Move them before you deploy:\n\n")
	gray.Printf("      SELECT COUNT(*) FROM users WHERE role = '%s';\n", upper)
	gray.Printf("      UPDATE users SET role = 'USER' WHERE role = '%s';\n\n", upper)
	gray.Println("  Then:")
	gray.Printf("    Rebuild the API and run `pnpm type-check`.\n")
	fmt.Println()
	return nil
}

// filesNaming finds project source files that still mention the role, so
// removing its declaration does not leave a guard asking for something nothing
// can be.
//
// Only the places a person writes: the generated declarations are what this
// command is about to remove, and listing them would be noise.
func filesNaming(root, role string) []string {
	needle := []byte(`"` + role + `"`)
	skipDir := map[string]bool{
		"node_modules": true, ".next": true, ".git": true, "dist": true,
		"build": true, ".grit": true, "vendor": true,
	}
	declarations := map[string]bool{
		"user.go": true, "user.ts": true, "index.ts": true, "users.ts": true,
	}

	var out []string
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if skipDir[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		switch filepath.Ext(path) {
		case ".go", ".ts", ".tsx":
		default:
			return nil
		}
		if declarations[info.Name()] {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		if bytes.Contains(raw, needle) {
			out = append(out, shown(root, path))
		}
		return nil
	})
	sort.Strings(out)
	return out
}

// removeLinesMatching deletes every line matching re, and returns how many.
func removeLinesMatching(path string, re *regexp.Regexp) int {
	raw, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	src := string(raw)
	crlf := strings.Contains(src, "\r\n")
	lines := strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n")

	kept := make([]string, 0, len(lines))
	removed := 0
	for _, line := range lines {
		if re.MatchString(line) {
			removed++
			continue
		}
		kept = append(kept, line)
	}
	if removed == 0 {
		return 0
	}
	out := strings.Join(kept, "\n")
	if crlf {
		out = strings.ReplaceAll(out, "\n", "\r\n")
	}
	if err := writeRepaired(path, out); err != nil {
		return 0
	}
	return removed
}

// removeFromList takes one quoted entry out of a bracketed list on one line,
// tidying the comma either side so the result still parses.
func removeFromList(path, entry string) bool {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	src := string(raw)
	for _, form := range []string{entry + ", ", ", " + entry, entry + ",", "," + entry, entry} {
		if strings.Contains(src, form) {
			out := strings.Replace(src, form, "", 1)
			if err := writeRepaired(path, out); err != nil {
				return false
			}
			return true
		}
	}
	return false
}

// removeFromUnion takes one member out of a TypeScript string union, which is
// separated by | rather than by commas.
func removeFromUnion(path, member string) bool {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	src := string(raw)
	quoted := `"` + member + `"`
	for _, form := range []string{` | ` + quoted, quoted + ` | `, `|` + quoted, quoted + `|`, quoted} {
		if strings.Contains(src, form) {
			out := strings.Replace(src, form, "", 1)
			if err := writeRepaired(path, out); err != nil {
				return false
			}
			return true
		}
	}
	return false
}
