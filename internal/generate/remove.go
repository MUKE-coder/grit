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

	// Every artefact `generate resource` can write, plus the differently-named
	// files the SCAFFOLD ships for its demo resources.
	//
	// Two naming conventions exist and both must be handled:
	//   - generated:  handlers/blog.go, services/blog.go
	//   - scaffolded: handlers/blog_handler.go, services/blog_service.go,
	//                 database/blogs_seeder.go
	// Missing the scaffold set is why `grit remove resource Blog` used to delete
	// only the model and leave the app unable to compile.
	//
	// Paths that don't exist are skipped, so listing every variant is safe.
	filesToDelete := []string{
		// --- Go API: generated naming ---
		filepath.Join(apiRoot, "internal", "models", names.Snake+".go"),
		filepath.Join(apiRoot, "internal", "services", names.Snake+".go"),
		filepath.Join(apiRoot, "internal", "handlers", names.Snake+".go"),
		filepath.Join(apiRoot, "internal", "handlers", names.Snake+"_import.go"),
		// --- Go API: scaffolded demo-resource naming ---
		filepath.Join(apiRoot, "internal", "handlers", names.Snake+"_handler.go"),
		filepath.Join(apiRoot, "internal", "services", names.Snake+"_service.go"),
		filepath.Join(apiRoot, "internal", "database", names.PluralSnake+"_seeder.go"),
		// --- shared package ---
		filepath.Join(sharedRoot, "schemas", names.Kebab+".ts"),
		filepath.Join(sharedRoot, "types", names.Kebab+".ts"),
		// --- web (Next.js + Vite layouts) ---
		filepath.Join(webRoot, "hooks", "use-"+names.PluralKebab+".ts"),
		filepath.Join(webRoot, "src", "hooks", "use-"+names.PluralKebab+".ts"),
		// --- admin: Next.js ---
		filepath.Join(adminRoot, "hooks", "use-"+names.PluralKebab+".ts"),
		filepath.Join(adminRoot, "resources", names.PluralKebab, names.PluralKebab+".ts"),
		filepath.Join(adminRoot, "resources", names.PluralKebab+".ts"),
		filepath.Join(adminRoot, "app", "(dashboard)", "resources", names.PluralKebab, "page.tsx"),
		// --- admin: TanStack/Vite ---
		filepath.Join(adminRoot, "src", "hooks", "use-"+names.PluralKebab+".ts"),
		filepath.Join(adminRoot, "src", "resources", names.PluralKebab, names.PluralKebab+".ts"),
		filepath.Join(adminRoot, "src", "resources", names.PluralKebab+".ts"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "resources", names.PluralKebab+".tsx"),
		filepath.Join(adminRoot, "src", "pages", "resources", names.PluralKebab+".tsx"),
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

	// The customisation overlay, which is hand-written and so gets different
	// treatment from everything above. Left in place it does not merely go
	// stale — it imports a type the shared package no longer exports, and the
	// admin stops type-checking over a resource that is supposed to be gone.
	// An untouched stub is deleted; anything with real work in it is set aside
	// as .bak, which keeps it off the TypeScript build without throwing it away.
	for _, overlay := range []string{
		filepath.Join(adminRoot, "resources", names.PluralKebab, names.PluralKebab+".custom.tsx"),
		filepath.Join(adminRoot, "resources", names.PluralKebab+".custom.tsx"),
		filepath.Join(adminRoot, "src", "resources", names.PluralKebab, names.PluralKebab+".custom.tsx"),
		filepath.Join(adminRoot, "src", "resources", names.PluralKebab+".custom.tsx"),
	} {
		body, err := os.ReadFile(overlay)
		if err != nil {
			continue
		}
		rel, _ := filepath.Rel(root, overlay)
		// The stub's only statement is an empty object literal. Anything else
		// means someone wrote something here.
		if strings.Contains(string(body), "= {};") {
			if err := os.Remove(overlay); err != nil {
				return fmt.Errorf("deleting %s: %w", overlay, err)
			}
			fmt.Printf("  ✗ %s\n", rel)
			continue
		}
		if err := os.Rename(overlay, overlay+".bak"); err != nil {
			return fmt.Errorf("archiving %s: %w", overlay, err)
		}
		fmt.Printf("  → %s.bak (your customisations, kept out of the build)\n", rel)
	}

	// Whole directories owned by the resource. RemoveAll (not Remove) because
	// these contain nested dynamic routes — admin's [id]/page.tsx and web's
	// [slug]/page.tsx — so an empty-dir check would never fire.
	dirsToDelete := []string{
		// Only if empty: a .bak the operator wanted kept lives in here.
		filepath.Join(adminRoot, "resources", names.PluralKebab),
		filepath.Join(adminRoot, "src", "resources", names.PluralKebab),
		filepath.Join(adminRoot, "app", "(dashboard)", "resources", names.PluralKebab),
		filepath.Join(webRoot, "app", names.Kebab),
	}
	for _, d := range dirsToDelete {
		if _, err := os.Stat(d); err != nil {
			continue
		}
		if err := os.RemoveAll(d); err == nil {
			rel, _ := filepath.Rel(root, d)
			fmt.Printf("  ✗ %s%c\n", rel, filepath.Separator)
		}
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

	// 1b. Remove the offline-sync registration.
	//
	// MUST run before the inline model-list surgery below. That step strips
	// ", &models.X{}" wherever it appears, which also matches
	//     syncRegistry.Register("xs", &models.X{})
	// leaving `Register("xs")` — a call with too few arguments that the
	// whole-line removal could no longer recognise.
	if fileExists(routesFile) {
		removeLinesContaining(routesFile,
			fmt.Sprintf("syncRegistry.Register(%q, &models.%s{})", names.Plural, names.Pascal))
	}

	// 2. Remove the model from every inline []interface{} model list (GORM Studio
	// mount, Pulse config, ...). Both separator positions are tried: the model
	// may sit mid-slice ("&M{}, ") or last ( ", &M{}" ), and only handling the
	// former left the trailing entry behind.
	if fileExists(routesFile) {
		removedAny := removeInlineText(routesFile, fmt.Sprintf("&models.%s{}, ", names.Pascal)) == nil
		if removeInlineText(routesFile, fmt.Sprintf(", &models.%s{}", names.Pascal)) == nil {
			removedAny = true
		}
		if removedAny {
			fmt.Println("  ✗ Removed model from studio/model lists")
		}
	}

	// 3. Remove handler initialization — two shapes exist:
	//   generated:  xHandler := &handlers.XHandler{ ... }   (multi-line literal)
	//   scaffolded: xHandler := handlers.NewXHandler(db)    (single-line ctor)
	if fileExists(routesFile) {
		removedInit := removeLineBlock(routesFile,
			fmt.Sprintf("%sHandler := &handlers.%sHandler{", names.Camel, names.Pascal),
			"}") == nil
		if removeLinesContaining(routesFile,
			fmt.Sprintf("%sHandler := handlers.New%sHandler(", names.Camel, names.Pascal)) == nil {
			removedInit = true
		}
		if removedInit {
			fmt.Println("  ✗ Removed handler initialization")
		}
	}

	// 3b. Remove a public route group keyed on the plural name, e.g.
	//     blogs := r.Group("/api/blogs") { ... }
	if fileExists(routesFile) {
		removeLineBlock(routesFile,
			fmt.Sprintf("%s := r.Group(", names.Plural),
			"}")
	}

	// 4+5. Remove routes (protected + admin)
	if fileExists(routesFile) {
		removed := removeLinesContaining(routesFile, names.Camel+"Handler.")
		if removed == nil {
			fmt.Println("  ✗ Removed API routes")
		}
	}

	// 5c. Remove the API-reference entries.
	if fileExists(routesFile) {
		if removeDocsRoutes(routesFile, "/api/"+apiVersion+"/"+names.Plural) == nil {
			fmt.Println("  ✗ Removed the /docs entries")
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

	// 10. Remove the resource import from the registry (Next.js and Vite admins).
	for _, reg := range []string{registryFile, filepath.Join(adminRoot, "src", "resources", "index.ts")} {
		if !fileExists(reg) {
			continue
		}
		// No closing quote in the needle: the import is "./carriers" in a flat
		// project and "./carriers/carriers" in a foldered one, and matching the
		// quote would only catch the first.
		if removeLinesContaining(reg, fmt.Sprintf(`from "./%s`, names.PluralKebab)) == nil {
			fmt.Println("  ✗ Removed resource import")
		}
	}

	// 11. Remove the resource from the registry list.
	//
	// Both spellings must be tried: `generate resource` exports the singular
	// (categoryResource) while the scaffold's demo resources use the plural
	// (blogsResource, usersResource). Handling only the singular left
	// "blogsResource," in the array with its import gone, which compiled but
	// failed type-check with "Cannot find name 'blogsResource'".
	for _, reg := range []string{registryFile, filepath.Join(adminRoot, "src", "resources", "index.ts")} {
		if !fileExists(reg) {
			continue
		}
		removed := removeLinesContaining(reg, fmt.Sprintf("%sResource,", names.Camel)) == nil
		if removeLinesContaining(reg, fmt.Sprintf("%sResource,", names.Plural)) == nil {
			removed = true
		}
		if removed {
			fmt.Println("  ✗ Removed resource from registry")
		}
	}

	// 13+14. Reverse the two switch-dispatch injections. These were the reason a
	// removed resource still failed to compile: both files keep a `case` arm
	// referencing models.<Name>, which no longer exists.
	//
	// generate also ADDS the encoding/json import to form_share_dispatch, so the
	// helper prunes imports that the removal just orphaned — otherwise we trade
	// "undefined: models.X" for "imported and not used".
	formShareFile := filepath.Join(apiRoot, "internal", "services", "form_share_dispatch.go")
	if fileExists(formShareFile) {
		// Two arms are injected: the create-dispatch and the public-fields lookup.
		n := removeCaseBlocks(formShareFile, fmt.Sprintf("case %q:", names.Pascal))
		pruneUnusedImports(formShareFile)
		if n > 0 {
			fmt.Println("  ✗ Removed form-share dispatch")
		}
	}

	// Every switch-dispatch file keyed by the plural resource name. Adding a new
	// dispatch to the scaffold? Add it here too, or removal silently leaves a
	// `case` arm referencing a deleted model.
	for _, d := range []struct{ file, label string }{
		{"resource_stats_dispatch.go", "resource-stats dispatch"},
		{"chart_dispatch.go", "chart dispatch"},
	} {
		path := filepath.Join(apiRoot, "internal", "services", d.file)
		if !fileExists(path) {
			continue
		}
		if removeCaseBlocks(path, fmt.Sprintf("case %q:", names.Plural)) > 0 {
			pruneUnusedImports(path)
			fmt.Printf("  ✗ Removed %s\n", d.label)
		}
	}

	// 15. Remove the seeder call. Scaffolded demo resources register a
	// Seed<Plural>(db) block in database/seed.go; deleting the seeder file
	// without this leaves an undefined reference.
	seedFile := filepath.Join(apiRoot, "internal", "database", "seed.go")
	if fileExists(seedFile) {
		if removeLineBlock(seedFile,
			fmt.Sprintf("if err := Seed%s(db); err != nil {", names.PluralPascal),
			"}") == nil {
			fmt.Println("  ✗ Removed seeder call")
		}
	}

	// 15b. Remove the resource's permissions from the authz catalog.
	//
	// Uses the module Key line as the anchor and cuts to the closing brace of
	// that literal. Leaving it behind would keep a deleted resource showing up
	// in the roles UI as a grantable permission — the same generate/remove drift
	// that made removal leave broken code in the first place.
	permsFile := filepath.Join(apiRoot, "internal", "authz", "permissions.go")
	if fileExists(permsFile) {
		if removeStructBlock(permsFile, fmt.Sprintf("Key:  %q,", names.Plural)) {
			fmt.Println("  ✗ Removed permissions from the authz catalog")
		}
	}

	// 16. Blog only: strip the "Recent Posts" section from the web home page.
	//
	// This is the one piece of demo content that lives inside a hand-designed
	// page rather than in a file of its own, so it can't be deleted wholesale.
	// The scaffold wraps it in grit:home:blog-* markers precisely so removal can
	// cut it out; without this the home page keeps importing the deleted
	// use-blogs hook and the web app fails to build.
	if strings.EqualFold(names.Pascal, "Blog") {
		homePage := filepath.Join(webRoot, "app", "page.tsx")
		if fileExists(homePage) {
			cut := removeMarkedRegion(homePage, "grit:home:blog-start", "grit:home:blog-end")
			if removeMarkedRegion(homePage, "grit:home:blog-hook-start", "grit:home:blog-hook-end") {
				cut = true
			}
			if removeLinesContaining(homePage, "// grit:home:blog-import") == nil {
				cut = true
			}
			removeLinesContaining(homePage, `from "@/hooks/use-blogs"`)
			if cut {
				fmt.Println("  ✗ Removed blog section from web home page")
			}
		}
	}

	fmt.Println()
	fmt.Printf("  ✅ Resource %s removed successfully!\n\n", names.Pascal)

	return nil
}

// removeStructBlock deletes a composite-literal element identified by a line it
// contains, e.g. the permission-catalog entry anchored on `Key:  "products",`.
//
// Brace-counting rather than indentation-matching: the catalog literal nests
// three levels (Module > Group > Feature), so scanning for "a line at the same
// indent" would stop at the first inner closing brace and leave a fragment
// behind. Starts from the anchor's opening "{" line and consumes until the brace
// depth returns to zero.
//
// Braces inside string literals would confuse the counter, but catalog entries
// hold identifiers and display names, so this stays honest for its one caller.
func removeStructBlock(filePath, anchor string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	lines := strings.Split(string(data), "\n")

	anchorIdx := -1
	for i, l := range lines {
		if strings.Contains(l, anchor) {
			anchorIdx = i
			break
		}
	}
	if anchorIdx == -1 {
		return false
	}

	// Walk back to the "{" that opens this element.
	start := -1
	for i := anchorIdx; i >= 0; i-- {
		if strings.HasSuffix(strings.TrimSpace(lines[i]), "{") {
			start = i
			break
		}
	}
	if start == -1 {
		return false
	}

	depth := 0
	end := -1
	for i := start; i < len(lines); i++ {
		depth += strings.Count(lines[i], "{") - strings.Count(lines[i], "}")
		if depth <= 0 {
			end = i
			break
		}
	}
	if end == -1 {
		return false
	}

	out := append([]string{}, lines[:start]...)
	out = append(out, lines[end+1:]...)
	return os.WriteFile(filePath, []byte(strings.Join(out, "\n")), 0644) == nil
}

// removeMarkedRegion deletes the lines between two grit markers, inclusive,
// reporting whether anything was cut.
//
// Used for demo content embedded in hand-designed files, where there is no
// standalone file to delete — the scaffold brackets the region so removal can
// take it out without parsing JSX.
func removeMarkedRegion(filePath, startMarker, endMarker string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	lines := strings.Split(string(data), "\n")

	start, end := -1, -1
	for i, l := range lines {
		if start == -1 && strings.Contains(l, startMarker) {
			start = i
			continue
		}
		if start != -1 && strings.Contains(l, endMarker) {
			end = i
			break
		}
	}
	// An unmatched start marker means the file was edited by hand; leave it
	// alone rather than truncating to the end of the file.
	if start == -1 || end == -1 {
		return false
	}

	out := append([]string{}, lines[:start]...)
	out = append(out, lines[end+1:]...)
	return os.WriteFile(filePath, []byte(strings.Join(out, "\n")), 0644) == nil
}

// removeCaseBlocks deletes every `case <label>:` arm from a Go switch, returning
// how many it removed.
//
// An arm runs from its `case` line until the next line at the same indentation
// that starts a new arm (`case`/`default:`) or is a grit injection marker —
// mirroring the shape `generate resource` writes. Trailing blank lines inside the
// arm are dropped with it so the switch doesn't accumulate gaps.
//
// A resource can appear in more than one switch in the same file (form_share
// dispatch injects both a create arm and a public-fields arm), hence "Blocks".
func removeCaseBlocks(filePath, caseLabel string) int {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0
	}

	lines := strings.Split(string(data), "\n")
	var out []string
	removed := 0

	for i := 0; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) != caseLabel {
			out = append(out, lines[i])
			continue
		}

		// Consume the arm.
		removed++
		indent := len(lines[i]) - len(strings.TrimLeft(lines[i], " \t"))
		i++
		for i < len(lines) {
			trimmed := strings.TrimSpace(lines[i])
			lineIndent := len(lines[i]) - len(strings.TrimLeft(lines[i], " \t"))
			isBoundary := lineIndent <= indent && trimmed != "" &&
				(strings.HasPrefix(trimmed, "case ") ||
					strings.HasPrefix(trimmed, "default:") ||
					strings.HasPrefix(trimmed, "// grit:") ||
					trimmed == "}")
			if isBoundary {
				break
			}
			i++
		}
		i-- // the outer loop's i++ will land us on the boundary line

		// Drop blank lines the arm left behind.
		for len(out) > 0 && strings.TrimSpace(out[len(out)-1]) == "" {
			out = out[:len(out)-1]
		}
	}

	if removed == 0 {
		return 0
	}
	if err := os.WriteFile(filePath, []byte(strings.Join(out, "\n")), 0644); err != nil {
		return 0
	}
	return removed
}

// pruneUnusedImports drops imports the file no longer references, so removing
// the last consumer of an import doesn't just trade "undefined: models.X" for
// "imported and not used".
//
// Only intended for the dispatch files we edit above — it is not a general
// goimports. Two details matter:
//   - Usage is tested with comments stripped. These files document the codegen
//     contract in prose ("re-marshals fields into the typed model via
//     json.Marshal(fields)"), and matching that comment made the import look
//     used when it wasn't.
//   - Blank (_) and dot (.) imports are never touched; their whole purpose is a
//     side effect the selector test can't see.
func pruneUnusedImports(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")

	// Locate the import block. Single-line `import "x"` files have nothing to
	// prune for our purposes, so only the parenthesised form is handled.
	start, end := -1, -1
	for i, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "import (") {
			start = i
			continue
		}
		if start != -1 && strings.TrimSpace(l) == ")" {
			end = i
			break
		}
	}
	if start == -1 || end == -1 {
		return
	}

	// Body with comments stripped — what counts as real usage.
	var body strings.Builder
	for i, l := range lines {
		if i >= start && i <= end {
			continue // the import block itself is not usage
		}
		if idx := strings.Index(l, "//"); idx != -1 {
			l = l[:idx]
		}
		body.WriteString(l)
		body.WriteString("\n")
	}
	code := body.String()

	var kept []string
	changed := false
	for i := start + 1; i < end; i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			kept = append(kept, line)
			continue
		}
		// Side-effect imports stay, always.
		if strings.HasPrefix(trimmed, "_ ") || strings.HasPrefix(trimmed, ". ") {
			kept = append(kept, line)
			continue
		}

		q := strings.Index(trimmed, `"`)
		if q == -1 {
			kept = append(kept, line)
			continue
		}
		path := strings.Trim(trimmed[q:], `"`)

		pkg := path
		if idx := strings.LastIndex(path, "/"); idx != -1 {
			pkg = path[idx+1:]
		}
		if alias := strings.TrimSpace(trimmed[:q]); alias != "" {
			pkg = alias
		}

		if strings.Contains(code, pkg+".") {
			kept = append(kept, line)
			continue
		}
		changed = true
	}
	if !changed {
		return
	}

	// Collapse blank runs the pruning left inside the block.
	var block []string
	for _, l := range kept {
		if strings.TrimSpace(l) == "" && (len(block) == 0 || strings.TrimSpace(block[len(block)-1]) == "") {
			continue
		}
		block = append(block, l)
	}
	for len(block) > 0 && strings.TrimSpace(block[len(block)-1]) == "" {
		block = block[:len(block)-1]
	}

	out := append([]string{}, lines[:start+1]...)
	out = append(out, block...)
	out = append(out, lines[end:]...)
	os.WriteFile(filePath, []byte(strings.Join(out, "\n")), 0644)
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

// removeDocsRoutes deletes the docs.Route(...) chains this resource added to
// routes.go.
//
// The chains are method-call fluent blocks, not brace-delimited, so neither
// removeLinesContaining nor removeLineBlock can express them: a chain starts at
// a docs.Route("<METHOD> <base>...") line and runs until the first line that
// does not end in a dot. Matching includes the closing quote so a resource
// whose plural prefixes another ("order" vs "orders") cannot take its
// neighbour's routes with it.
func removeDocsRoutes(filePath, base string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	result := make([]string, 0, len(lines))
	removed := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		isStart := strings.Contains(line, "docs.Route(\"") &&
			(strings.Contains(line, base+"\"") || strings.Contains(line, base+"/:id\""))
		if !isStart {
			result = append(result, line)
			continue
		}

		removed = true
		// Consume the chain: every line but the last ends with a dot.
		for i < len(lines) && strings.HasSuffix(strings.TrimRight(lines[i], " 	"), ".") {
			i++
		}
		// i now sits on the chain's final line, which the loop's i++ skips.
	}

	if !removed {
		return fmt.Errorf("no docs routes for %q", base)
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
