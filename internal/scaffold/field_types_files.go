package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// fieldTypesFiles is internal/fieldtypes, which ships with every API.
func fieldTypesFiles(apiRoot string) map[string]string {
	dir := filepath.Join(apiRoot, "internal", "fieldtypes")
	return map[string]string{
		filepath.Join(dir, "fieldtypes.go"):      apiFieldTypesGo(),
		filepath.Join(dir, "countries.go"):       apiFieldTypesCountriesGo(),
		filepath.Join(dir, "samples.go"):         apiFieldTypesSamplesGo(),
		filepath.Join(dir, "fieldtypes_test.go"): apiFieldTypesTestGo(),
	}
}

// phoneFiles is internal/phone, written when a project first has a tel field.
func phoneFiles(apiRoot string) map[string]string {
	dir := filepath.Join(apiRoot, "internal", "phone")
	return map[string]string{
		filepath.Join(dir, "phone.go"):      apiPhoneGo(),
		filepath.Join(dir, "phone_test.go"): apiPhoneTestGo(),
	}
}

// writeFieldTypesFiles writes internal/fieldtypes through the manifest guard,
// for the scaffold and for grit upgrade, and refreshes internal/phone where a
// project already has it. The phone package is never added here: a project
// without a tel field has no use for its module.
func writeFieldTypesFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()
	files := fieldTypesFiles(apiRoot)
	if fileExists(filepath.Join(apiRoot, "internal", "phone", "phone.go")) {
		for path, body := range phoneFiles(apiRoot) {
			files[path] = body
		}
	}
	for path, body := range files {
		if err := writeFile(path, strings.ReplaceAll(body, "{{MODULE}}", module)); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

// writeMissing writes each file that is not there yet and reports whether it
// wrote any. The generator uses it, so a project that has edited one of these
// files keeps its edit.
func writeMissing(files map[string]string, module string) (bool, error) {
	wrote := false
	for path, body := range files {
		if fileExists(path) {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return wrote, err
		}
		if err := os.WriteFile(path, []byte(strings.ReplaceAll(body, "{{MODULE}}", module)), 0o644); err != nil {
			return wrote, fmt.Errorf("writing %s: %w", path, err)
		}
		manifest.Refresh(path)
		wrote = true
	}
	return wrote, nil
}

// WriteFieldTypesPackage adds internal/fieldtypes to a project that does not
// have it. For the generator, ahead of the first formatted field.
func WriteFieldTypesPackage(apiRoot, module string) (bool, error) {
	return writeMissing(fieldTypesFiles(apiRoot), module)
}

// WritePhonePackage adds internal/phone to a project that does not have it.
// The caller adds PhoneNumbersModule to go.mod.
func WritePhonePackage(apiRoot, module string) (bool, error) {
	return writeMissing(phoneFiles(apiRoot), module)
}

// EnsureFieldTypesWiring makes the project call fieldtypes.Install when it
// connects. Idempotent, and it refuses rather than guesses when database.go
// has been reshaped past recognition.
func EnsureFieldTypesWiring(apiRoot, module string) error {
	path := filepath.Join(apiRoot, "internal", "database", "database.go")
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("the field type checks need %s: %w", path, err)
	}
	crlf := strings.Contains(string(raw), "\r\n")
	content := strings.ReplaceAll(string(raw), "\r\n", "\n")
	if strings.Contains(content, "fieldtypes.Install(db)") {
		return nil
	}
	const anchor = "\tsqlDB, err := db.DB()"
	if !strings.Contains(content, anchor) {
		return fmt.Errorf("could not find where to install the field type checks in %s.\n\n"+
			"Add this call after connecting, or email, tel and the other formatted fields are stored unchecked:\n\n  fieldtypes.Install(db)", path)
	}
	content = strings.Replace(content, anchor, fieldTypesConnectHook+anchor, 1)
	// Where the template has it, between crypto and paginate, when that group
	// is there; a group of its own otherwise.
	importLine := "\t\"" + module + "/internal/fieldtypes\"\n"
	if cryptoLine := "\t\"" + module + "/internal/crypto\"\n"; strings.Contains(content, cryptoLine) {
		content = strings.Replace(content, cryptoLine, cryptoLine+importLine, 1)
	} else {
		var ok bool
		if content, ok = addImportGroup(content, module+"/internal/fieldtypes"); !ok {
			return fmt.Errorf("could not add the fieldtypes import to %s: add it by hand", path)
		}
	}
	if crlf {
		content = strings.ReplaceAll(content, "\n", "\r\n")
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	manifest.Refresh(path)
	fmt.Println("  ✓ Wired the field type checks into database.go")
	return nil
}

// repairFieldTypes is the upgrade step: the package, and the call that
// installs it.
func repairFieldTypes(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	if !fileExists(filepath.Join(apiRoot, "internal", "database", "database.go")) {
		return nil
	}
	if err := writeFieldTypesFiles(root, opts); err != nil {
		return err
	}
	return EnsureFieldTypesWiring(apiRoot, opts.Module())
}

// ── packages/shared ─────────────────────────────────────────────────────────

// sharedFieldFormatFiles are the shared Zod schemas generated resources import
// for their formatted fields.
func sharedFieldFormatFiles(sharedRoot string) map[string]string {
	return map[string]string{
		filepath.Join(sharedRoot, "schemas", "field-formats.ts"): sharedFieldFormatsSchema(),
		filepath.Join(sharedRoot, "schemas", "phone.ts"):         sharedPhoneSchema(),
	}
}

// WriteSharedFieldFormats adds the shared schemas to a project that lacks
// them, and libphonenumber-js to the package that compiles them when tel is
// in use. Returns what it added.
func WriteSharedFieldFormats(sharedRoot string, tel bool) ([]string, error) {
	var added []string
	for path, body := range sharedFieldFormatFiles(sharedRoot) {
		if fileExists(path) {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return added, err
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return added, fmt.Errorf("writing %s: %w", path, err)
		}
		manifest.Refresh(path)
		added = append(added, path)
	}
	if tel {
		// packages/shared/package.json in a monorepo; the SPA's own
		// package.json when the shared files are mirrored into frontend/src.
		pkg := filepath.Join(sharedRoot, "package.json")
		if !fileExists(pkg) {
			pkg = filepath.Join(sharedRoot, "..", "..", "package.json")
		}
		if fileExists(pkg) && !fileContains(pkg, `"libphonenumber-js"`) {
			ok, err := insertMissingAfter(pkg, `"dependencies": {`, []string{
				fmt.Sprintf(`    "libphonenumber-js": %q,`, libphonenumberVersion),
			})
			if err != nil {
				return added, err
			}
			if ok {
				added = append(added, pkg+" (libphonenumber-js, run pnpm install)")
			}
		}
	}
	return added, nil
}

// ensureSharedFieldFormats is the upgrade step for packages/shared: the two
// schema modules, which are new files, so writing them costs nothing a
// project has.
func ensureSharedFieldFormats(sharedRoot string) error {
	if !dirExists(sharedRoot) {
		return nil
	}
	for path, body := range sharedFieldFormatFiles(sharedRoot) {
		if err := writeFile(path, body); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}
