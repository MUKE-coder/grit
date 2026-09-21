package generate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// ensureFieldTypeSupport makes sure the project can carry the formatted fields
// it is about to be given: internal/fieldtypes and its call at connect time,
// internal/phone and its module for a tel field, and the shared Zod schemas.
//
// A project scaffolded before these types existed has none of them, and a
// model whose format tag names a check the project cannot run would refuse
// every write. So they are added here, whether or not grit upgrade has run.
func (g *Generator) ensureFieldTypeSupport(fields []Field) error {
	formatted, tel := false, false
	for _, f := range fields {
		if f.IsFormatted() {
			formatted = true
		}
		if f.IsTel() {
			tel = true
		}
	}
	if !formatted {
		return nil
	}

	apiRoot := g.APIRoot()
	if wrote, err := scaffold.WriteFieldTypesPackage(apiRoot, g.Module); err != nil {
		return err
	} else if wrote {
		fmt.Println("  ✓ Added internal/fieldtypes (checks email, url, domain, country, color, percent, rating, time and json on every write)")
	}
	if fileExists(filepath.Join(apiRoot, "internal", "database", "database.go")) {
		if err := scaffold.EnsureFieldTypesWiring(apiRoot, g.Module); err != nil {
			return err
		}
	} else {
		fmt.Println("  ⚠ No internal/database/database.go: call fieldtypes.Install(db) after connecting,")
		fmt.Println("    or the formatted fields are stored unchecked.")
	}

	if tel {
		if wrote, err := scaffold.WritePhonePackage(apiRoot, g.Module); err != nil {
			return err
		} else if wrote {
			fmt.Println("  ✓ Added internal/phone (validates tel fields with libphonenumber)")
		}
		if !fileContainsText(filepath.Join(apiRoot, "go.mod"), scaffold.PhoneNumbersModule+" ") {
			spec := scaffold.PhoneNumbersModule + "@" + scaffold.PhoneNumbersVersion
			if err := goGetModule(apiRoot, spec); err != nil {
				return err
			}
			fmt.Printf("  ✓ Added %s %s\n", scaffold.PhoneNumbersModule, scaffold.PhoneNumbersVersion)
		}
	}

	if shared := g.SharedRoot(); dirExists(shared) {
		added, err := scaffold.WriteSharedFieldFormats(shared, tel)
		if err != nil {
			return err
		}
		for _, a := range added {
			fmt.Printf("  ✓ %s\n", relativeToRoot(g.Root, a))
		}
	}

	// The admin's inputs for these types arrive with grit upgrade. Without them
	// the form builder has no case for the type and the field is not drawn.
	builder := filepath.Join(g.AdminCodeRoot(), "components", "forms", "form-builder.tsx")
	if fileExists(builder) && !fileContainsText(builder, `case "tel":`) {
		fmt.Println("  ⚠ The admin panel predates the email, url, domain, tel, country, color, percent,")
		fmt.Println("    rating, time and json inputs, so these fields will not show in its forms.")
		fmt.Println("    Run grit upgrade, then pnpm install.")
	}
	return nil
}

// goGetModule adds a module to the API's go.mod. A variable so a test can see
// what would be fetched without the network.
var goGetModule = func(dir, spec string) error {
	cmd := exec.Command("go", "get", spec)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("adding %s: %w\n%s", spec, err, string(out))
	}
	return nil
}

// tsNamedImportRe finds `import { A, B } from "<module>";`.
func tsNamedImportRe(module string) *regexp.Regexp {
	return regexp.MustCompile(`(?m)^import \{([^}]*)\} from "` + regexp.QuoteMeta(module) + `";$`)
}

// ensureTSNamedImport makes a TypeScript file import name from module, adding
// it to an existing import of that module or as a new line under the zod one.
func ensureTSNamedImport(path, module, name string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	content := string(raw)
	re := tsNamedImportRe(module)
	if m := re.FindStringSubmatchIndex(content); m != nil {
		names := strings.Split(content[m[2]:m[3]], ",")
		for _, n := range names {
			if strings.TrimSpace(n) == name {
				return nil
			}
		}
		list := []string{}
		for _, n := range names {
			if n = strings.TrimSpace(n); n != "" {
				list = append(list, n)
			}
		}
		list = append(list, name)
		content = content[:m[2]] + " " + strings.Join(list, ", ") + " " + content[m[3]:]
		return os.WriteFile(path, []byte(content), 0o644)
	}
	line := fmt.Sprintf(`import { %s } from "%s";`, name, module)
	const zodImport = `import { z } from "zod";`
	if strings.Contains(content, zodImport) {
		content = strings.Replace(content, zodImport, zodImport+"\n"+line, 1)
	} else {
		content = line + "\n" + content
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
