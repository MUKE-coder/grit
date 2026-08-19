package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// The admin side of `grit add variants`.
//
// Every file it writes is shared. The matrix editor takes the resource off the
// detail page's own props, so Products and Courses drive the same component and
// there is no generated second copy to go stale when one of them is improved.
// The only per-resource thing here is a single line in that resource's custom
// file, saying which slot the editor hangs off.

// writeVariantAdmin installs the matrix editor and the option library.
func writeVariantAdmin(root string, names variantNameSet) error {
	adminRoot := filepath.Join(root, "apps", "admin")
	if !dirExists(adminRoot) {
		fmt.Println("  • No apps/admin found, so the matrix UI was skipped (the API is complete).")
		return nil
	}

	// The TanStack admin keeps its pages under src/ and routes them through
	// TanStack Router rather than the app directory. The slot names and the
	// hooks are the same, but the route file is not, so rather than write a
	// page.tsx it will never load, this says so and stops.
	if dirExists(filepath.Join(adminRoot, "src", "routes")) &&
		!dirExists(filepath.Join(adminRoot, "app")) {
		fmt.Println("  • This is the TanStack admin. The API and the hooks are installed;")
		fmt.Println("    the matrix editor's route file is written for the Next admin only.")
		return nil
	}

	green := color.New(color.FgHiGreen)

	files := map[string]string{
		filepath.Join(adminRoot, "hooks", "use-variants.ts"):                               scaffold.AdminVariantHooksTS(),
		filepath.Join(adminRoot, "components", "variants", "variant-matrix.tsx"):           scaffold.AdminVariantMatrixTSX(),
		filepath.Join(adminRoot, "components", "variants", "option-library.tsx"):           scaffold.AdminOptionLibraryTSX(),
		filepath.Join(adminRoot, "resources", "options", "options.ts"):                     scaffold.AdminOptionsResourceTS(),
		filepath.Join(adminRoot, "app", "(dashboard)", "resources", "options", "page.tsx"): scaffold.AdminOptionsPageTSX(),
	}

	for path, body := range files {
		rel := relFromRoot(root, path)
		// Overwritten on a second run, unlike the resource's custom file. These
		// are the framework's, and a project that reruns the command after an
		// upgrade should get the improved editor rather than the one it had.
		if err := writeFileWithDirs(path, body); err != nil {
			return fmt.Errorf("writing %s: %w", filepath.Base(path), err)
		}
		green.Printf("  ✓ %s\n", rel)
	}

	if err := registerOptionResource(adminRoot); err != nil {
		return err
	}
	return attachMatrixToResource(adminRoot, names)
}

// registerOptionResource adds the option library to the admin registry, which
// is what puts it in the sidebar.
func registerOptionResource(adminRoot string) error {
	path := filepath.Join(adminRoot, "resources", "index.ts")
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("  • No resources/index.ts found, so Options was not added to the sidebar.")
		return nil
	}
	content := string(data)

	if strings.Contains(content, "optionResource") {
		fmt.Println("  • Options already in the sidebar")
		return nil
	}

	if err := injectBefore(path, "// grit:resources",
		`import { optionResource } from "./options/options";`); err != nil {
		fmt.Printf("  Could not register the option library: %v\n", err)
		return nil
	}
	if err := injectBefore(path, "// grit:resource-list", "  optionResource,"); err != nil {
		fmt.Printf("  Could not add the option library to the list: %v\n", err)
		return nil
	}
	manifest.Refresh(path)
	fmt.Println("  ✓ Options added to the admin sidebar")
	return nil
}

// attachMatrixToResource points the resource's DetailAside slot at the matrix.
//
// This is the one file here that belongs to the developer: the generator writes
// it once and never again, which is the promise that makes it safe to put your
// own cell renderers in. So the edit is attempted only on a customisation that
// is still empty, and anything else gets the two lines to paste rather than a
// rewrite of work somebody did.
func attachMatrixToResource(adminRoot string, names variantNameSet) error {
	path := findCustomFile(adminRoot, names)
	if path == "" {
		printMatrixByHand(names, "")
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		printMatrixByHand(names, path)
		return nil
	}
	content := string(data)

	if strings.Contains(content, "VariantMatrix") {
		fmt.Println("  • Matrix editor already attached to the detail page")
		return nil
	}

	// The pristine shape the generator writes. Matched exactly, because a
	// customisation with anything at all in it is somebody's work and a
	// regex-driven merge into it is how that work gets lost.
	empty := "= {};"
	if !strings.Contains(content, empty) {
		printMatrixByHand(names, path)
		return nil
	}

	updated := strings.Replace(content, empty,
		"= {\n  // The variant matrix, on this resource's detail page.\n"+
			"  components: { DetailAside: VariantMatrix },\n};", 1)

	// Import added after the body edit, so a failure to find the anchor leaves
	// the file exactly as it was rather than importing something unused.
	anchor := "import type { ResourceCustomisation }"
	if !strings.Contains(updated, anchor) {
		printMatrixByHand(names, path)
		return nil
	}
	updated = strings.Replace(updated, anchor,
		"import { VariantMatrix } from \"@/components/variants/variant-matrix\";\n"+anchor, 1)

	if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
		return fmt.Errorf("attaching the matrix editor: %w", err)
	}
	manifest.Refresh(path)
	color.New(color.FgHiGreen).Printf("  ✓ Matrix editor attached to the %s detail page\n", names.Pascal)
	return nil
}

// findCustomFile locates the resource's customisation file, whichever layout
// the admin uses.
func findCustomFile(adminRoot string, names variantNameSet) string {
	pluralKebab := strings.ReplaceAll(names.Plural, "_", "-")
	for _, candidate := range []string{
		filepath.Join(adminRoot, "resources", pluralKebab, pluralKebab+".custom.tsx"),
		filepath.Join(adminRoot, "resources", pluralKebab+".custom.tsx"),
	} {
		if fileExists(candidate) {
			return candidate
		}
	}
	return ""
}

// printMatrixByHand tells the developer where the two lines go.
func printMatrixByHand(names variantNameSet, path string) {
	where := "your " + names.Pascal + " customisation file"
	if path != "" {
		where = filepath.Base(path)
	}
	fmt.Printf("  • %s already has customisations, so it was left alone.\n", where)
	fmt.Println("    Add these two lines to put the matrix on the detail page:")
	fmt.Println(`      import { VariantMatrix } from "@/components/variants/variant-matrix";`)
	fmt.Println("      components: { DetailAside: VariantMatrix },")
}
