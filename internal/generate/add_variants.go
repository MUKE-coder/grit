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

// AddVariants installs the normalised variant system and attaches it to one
// resource.
//
// Its own command rather than a flag on `generate resource` because two of the
// five tables are shared. Options and their values belong to the shop, not to a
// product: Colour is Colour whether it is on a shirt or a phone case, and a
// per-resource copy gives you fourteen spellings of it and a filter that can
// only match one.
//
// Running it a second time for another resource adds only that resource's
// tables and leaves the shared ones alone.
func AddVariants(resource string) error {
	root, err := scaffold.FindProjectRoot()
	if err != nil {
		return err
	}

	names := variantNames(resource)
	apiRoot := filepath.Join(root, "apps", "api")
	if !dirExists(apiRoot) {
		apiRoot = filepath.Join(root, "api")
	}
	if !dirExists(apiRoot) {
		return fmt.Errorf("no API found: run this from a Grit project")
	}

	modelPath := filepath.Join(apiRoot, "internal", "models", names.Snake+".go")
	if !fileExists(modelPath) {
		return fmt.Errorf(
			"no %s model found at internal/models/%s.go.\n"+
				"Generate the resource first:\n"+
				"  grit generate resource %s --fields \"name:string,price:float\"",
			names.Pascal, names.Snake, names.Pascal)
	}

	module, err := readModulePath(root, "")
	if err != nil {
		return err
	}

	green := color.New(color.FgHiGreen)
	purple := color.New(color.FgHiMagenta, color.Bold)
	purple.Printf("\n  Adding variants to %s\n\n", names.Pascal)

	release, err := manifest.Start(root, "", "variants")
	if err != nil {
		return err
	}
	defer func() { _ = release() }()

	hasSlug, hasArchivedAt := modelShape(modelPath)

	// The shared half, written once however many resources have variants.
	//
	// Each of these is shop-wide rather than per-resource: the option models
	// because Colour is Colour on any product, the published view types because
	// a second copy of them is a second declaration in package handlers, and the
	// option seeder because no one resource owns the library.
	shared := map[string]string{
		filepath.Join(apiRoot, "internal", "models", "option.go"):           scaffold.APIOptionModelGo(module),
		filepath.Join(apiRoot, "internal", "handlers", "option_public.go"):  scaffold.APIOptionPublicGo(module),
		filepath.Join(apiRoot, "internal", "database", "options_seeder.go"): scaffold.APIOptionSeederGo(module),
	}
	for path, body := range shared {
		rel := relFromRoot(root, path)
		if fileExists(path) {
			fmt.Printf("  • %s exists, left alone (shared across resources)\n", rel)
			continue
		}
		if err := writeFileWithDirs(path, body); err != nil {
			return fmt.Errorf("writing %s: %w", filepath.Base(path), err)
		}
		green.Printf("  ✓ %s\n", rel)
	}

	// The per-resource half.
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "handlers", names.Snake+"_variant_public.go"):  scaffold.APIVariantPublicGo(module, names.Pascal, names.Snake, names.Plural, hasSlug, hasArchivedAt),
		filepath.Join(apiRoot, "internal", "database", names.Snake+"_variants_seeder.go"): scaffold.APIVariantSeederGo(module, names.Pascal, names.Snake, names.Plural, hasSlug),
		filepath.Join(apiRoot, "internal", "models", names.Snake+"_variant.go"):           scaffold.APIVariantModelGo(module, names.Pascal, names.Snake),
		filepath.Join(apiRoot, "internal", "services", names.Snake+"_variants.go"):        scaffold.APIVariantServiceGo(module, names.Pascal, names.Snake, names.Plural),
		filepath.Join(apiRoot, "internal", "services", names.Snake+"_variants_test.go"):   scaffold.APIVariantServiceTestGo(module, names.Pascal, names.Snake),
		filepath.Join(apiRoot, "internal", "handlers", names.Snake+"_variant.go"):         scaffold.APIVariantHandlerGo(module, names.Pascal, names.Snake, names.Plural),
	}
	for path, body := range files {
		if err := writeFileWithDirs(path, body); err != nil {
			return fmt.Errorf("writing %s: %w", filepath.Base(path), err)
		}
		green.Printf("  ✓ %s\n", relFromRoot(root, path))
	}

	if err := registerVariantModels(apiRoot, names); err != nil {
		return err
	}
	if err := mountVariantRoutes(apiRoot, names); err != nil {
		return err
	}
	if err := mountVariantPublicRoute(apiRoot, names); err != nil {
		return err
	}
	if err := registerVariantSeeder(apiRoot, names); err != nil {
		return err
	}
	if err := registerOptionStats(apiRoot); err != nil {
		return err
	}
	if err := writeVariantAdmin(root, names); err != nil {
		return err
	}

	fmt.Println()
	green.Println("  Variants installed.")
	fmt.Println()
	fmt.Println("  Next:")
	fmt.Println("    grit migrate                       # create the five tables")
	fmt.Println("    grit seed                          # a Colour x Size matrix to look at")
	fmt.Printf("    The matrix editor is on any %s's detail page in the admin,\n", names.Pascal)
	fmt.Println("    and the shared option library is under Options in the sidebar.")
	fmt.Println()
	return nil
}

// registerOptionStats teaches the dashboard how to count options.
//
// The admin's dashboard builds a stat card and a latest-rows table for every
// resource in its registry, and the option library is now one of those. The
// server answers those from a whitelist rather than from the resource name,
// deliberately — the dispatch is what stops a compromised admin token dumping
// arbitrary tables by guessing — so a resource that is not in it gets a 400 on
// every dashboard load.
//
// Shared, and written once: options are shop-wide, so a second resource with
// variants must not add the case twice.
func registerOptionStats(apiRoot string) error {
	const marker = "// grit:resource-stats:dispatch"
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "services", "resource_stats_dispatch.go"): "\tcase \"options\":\n\t\treturn reflectiveResourceStats(db, resourceName, &models.Option{}, filter)",
		filepath.Join(apiRoot, "internal", "services", "chart_dispatch.go"):          "\tcase \"options\":\n\t\treturn reflectiveChart(db, &models.Option{}, params)",
	}

	added := 0
	for path, block := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if strings.Contains(string(data), `case "options":`) {
			continue
		}
		if err := injectBefore(path, marker, block); err != nil {
			fmt.Printf("  Could not register option stats in %s: %v\n", filepath.Base(path), err)
			continue
		}
		manifest.Refresh(path)
		added++
	}
	if added > 0 {
		fmt.Println("  ✓ Options registered with the dashboard stats")
	}
	return nil
}

// relFromRoot renders a path the way the rest of the CLI's output does.
func relFromRoot(root, path string) string {
	return strings.TrimPrefix(filepath.ToSlash(strings.TrimPrefix(path, root)), "/")
}

// registerVariantSeeder adds the seeder to Seed().
//
// After the resource's own seeder, which is what the marker gives for free:
// injectBefore appends, and a matrix cannot be attached to rows that do not
// exist yet.
func registerVariantSeeder(apiRoot string, names variantNameSet) error {
	path := filepath.Join(apiRoot, "internal", "database", "seed.go")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	call := "Seed" + names.Pascal + "Variants(db)"
	if strings.Contains(string(data), call) {
		fmt.Println("  • Seeder already registered")
		return nil
	}

	block := fmt.Sprintf("\tif err := %s; err != nil {\n"+
		"\t\treturn fmt.Errorf(\"seeding %s variants: %%w\", err)\n"+
		"\t}\n", call, names.Snake)

	if err := injectBefore(path, "// grit:seeders", block); err != nil {
		fmt.Printf("  Could not register the seeder automatically: %v\n", err)
		fmt.Printf("    Add this to Seed() in internal/database/seed.go:\n      %s\n", call)
		return nil
	}
	manifest.Refresh(path)
	fmt.Printf("  ✓ Registered Seed%sVariants with grit seed\n", names.Pascal)
	return nil
}

type variantNameSet struct {
	Pascal string
	Snake  string
	Kebab  string
	Plural string
}

func variantNames(resource string) variantNameSet {
	pascal := toPascalCase(resource)
	snake := toSnakeCase(pascal)
	return variantNameSet{
		Pascal: pascal,
		Snake:  snake,
		Kebab:  strings.ReplaceAll(snake, "_", "-"),
		Plural: Pluralize(snake),
	}
}

// modelShape reports the two columns the public endpoint and the seeder have to
// know about before they can be written.
//
// Read off the model rather than assumed. A resource generated without a slug
// has no such column, and a public lookup on it would be a SQL error rather
// than the 404 the caller deserves.
func modelShape(modelPath string) (hasSlug, hasArchivedAt bool) {
	data, err := os.ReadFile(modelPath)
	if err != nil {
		return false, false
	}
	content := string(data)
	return strings.Contains(content, "Slug "), strings.Contains(content, "ArchivedAt ")
}

// registerVariantModels adds the new models to AutoMigrate and GORM Studio.
//
// Without this the tables are never created and the first request fails on a
// relation that does not exist, which reads as a bug in Grit rather than a
// migration nobody ran.
func registerVariantModels(apiRoot string, names variantNameSet) error {
	// internal/models/user.go, which is where the registry and its marker live.
	// Pointing this at routes.go, as it did until v3.167.0, meant the marker was
	// never found and the four tables were never migrated: every variant request
	// then failed on a relation that does not exist, which reads as a bug in
	// Grit rather than as a migration nobody ran.
	path := filepath.Join(apiRoot, "internal", "models", "user.go")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)

	toAdd := []string{"Option", "OptionValue", names.Pascal + "Option", names.Pascal + "Variant"}
	var missing []string
	for _, model := range toAdd {
		// The registry lives inside package models, so entries are written
		// &Option{} and not &models.Option{}. The leading ampersand is what
		// keeps "Option" from matching "ProductOption" and deciding it is
		// already registered.
		if !strings.Contains(content, "&"+model+"{},") {
			missing = append(missing, model)
		}
	}
	if len(missing) == 0 {
		fmt.Println("  • Models already registered with AutoMigrate")
		return nil
	}

	block := ""
	for _, model := range missing {
		block += "\t\t&" + model + "{},\n"
	}
	block = strings.TrimSuffix(block, "\n")

	if err := injectBefore(path, "// grit:models", block); err != nil {
		fmt.Printf("  Could not register the models automatically: %v\n", err)
		fmt.Println("    Add these to the list in internal/models/user.go:")
		for _, model := range missing {
			fmt.Printf("      &%s{},\n", model)
		}
		return nil
	}
	manifest.Refresh(path)
	fmt.Printf("  ✓ Registered %d model(s) with AutoMigrate\n", len(missing))
	return nil
}

// mountVariantRoutes adds the option and variant endpoints.
//
// Split into a shared half and a per-resource half, because the option library
// is shop-wide and gin panics at boot on two handlers for one method and path.
// Mounting /options again for a second resource would take the whole API down,
// and the symptom is a server that will not start after a command that said it
// succeeded.
//
// The per-resource update route is /<kebab>-variants/:id rather than the
// obvious /<plural>/:id/variants/:variant, because gin routes a parameter and a
// static segment at the same position by panicking, and /variants/:id on its
// own collides the moment a second resource wants one.
func mountVariantRoutes(apiRoot string, names variantNameSet) error {
	path := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)

	if strings.Contains(content, names.Snake+"VariantHandler :=") {
		return upgradeVariantRoutes(path, content, names)
	}

	handler := fmt.Sprintf("\t\t%sVariantHandler := handlers.New%sVariantHandler(db)",
		names.Snake, names.Pascal)

	// The shared library, mounted by whichever resource got variants first.
	shared := ""
	if !strings.Contains(content, "VariantHandler.ListOptions") {
		shared = fmt.Sprintf(`
		protected.GET("/options", %[1]sVariantHandler.ListOptions)
		protected.POST("/options", %[1]sVariantHandler.CreateOption)
		protected.DELETE("/options/:id", %[1]sVariantHandler.DeleteOption)
		protected.POST("/options/:id/values", %[1]sVariantHandler.CreateOptionValue)
		protected.DELETE("/option-values/:id", %[1]sVariantHandler.DeleteOptionValue)`,
			names.Snake)
	}

	perResource := fmt.Sprintf(`
		protected.GET("/%[3]s/:id/variants", %[1]sVariantHandler.List)
		protected.POST("/%[3]s/:id/variants/generate", %[1]sVariantHandler.GenerateMatrix)
		protected.PUT("/%[3]s/:id/options", %[1]sVariantHandler.SetOptions)
		protected.PATCH("/%[2]s-variants/:id", %[1]sVariantHandler.Update)`,
		names.Snake, names.Kebab, names.Plural)

	if err := injectBefore(path, "// grit:routes:custom", handler+shared+perResource); err != nil {
		fmt.Printf("  Could not mount the routes automatically: %v\n", err)
		return nil
	}
	manifest.Refresh(path)
	if shared != "" {
		fmt.Println("  ✓ Mounted the shared /options library")
	}
	fmt.Printf("  ✓ Mounted /%s/:id/variants and /%s-variants/:id\n", names.Plural, names.Kebab)
	return nil
}

// upgradeVariantRoutes brings a project installed before v3.167.0 up to the
// current route layout.
//
// Two things changed, and both of them are the kind that only bites the second
// time. The per-variant update lived at /variants/:id, which collides the moment
// another resource wants one; and there was no way to delete an option at all,
// so the library only ever grew.
//
// Rewriting rather than appending, because appending a second PATCH for the
// same path is a gin panic at boot and the old path is not worth keeping: the
// admin is the only thing that ever called it.
func upgradeVariantRoutes(path, content string, names variantNameSet) error {
	updated := content
	changes := []string{}

	oldUpdate := fmt.Sprintf("protected.PATCH(\"/variants/:id\", %sVariantHandler.Update)", names.Snake)
	newUpdate := fmt.Sprintf("protected.PATCH(\"/%s-variants/:id\", %sVariantHandler.Update)",
		names.Kebab, names.Snake)
	if strings.Contains(updated, oldUpdate) {
		updated = strings.Replace(updated, oldUpdate, newUpdate, 1)
		changes = append(changes, "moved the update route to /"+names.Kebab+"-variants/:id")
	}

	// Only for whichever resource mounted the shared library, which is the one
	// whose handler carries CreateOption.
	createOption := fmt.Sprintf("protected.POST(\"/options\", %sVariantHandler.CreateOption)", names.Snake)
	deleteOption := fmt.Sprintf("protected.DELETE(\"/options/:id\", %sVariantHandler.DeleteOption)", names.Snake)
	if strings.Contains(updated, createOption) && !strings.Contains(updated, deleteOption) {
		updated = strings.Replace(updated, createOption, createOption+"\n\t\t"+deleteOption, 1)
		changes = append(changes, "added DELETE /options/:id")
	}

	if len(changes) == 0 {
		fmt.Println("  • Routes already mounted")
		return nil
	}
	if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
		return fmt.Errorf("upgrading the variant routes: %w", err)
	}
	manifest.Refresh(path)
	for _, change := range changes {
		fmt.Printf("  ✓ Routes: %s\n", change)
	}
	return nil
}

// mountVariantPublicRoute adds the storefront's read endpoint to the public
// group, where one exists.
//
// A project generated without --public on this resource has no such group, and
// that is not an error: variants are perfectly useful to an admin-only app. The
// command says what to do about it and carries on.
func mountVariantPublicRoute(apiRoot string, names variantNameSet) error {
	path := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)

	if strings.Contains(content, names.Snake+"PublicVariants") {
		fmt.Println("  • Public variant route already mounted")
		return nil
	}
	if !strings.Contains(content, "// grit:routes:public") {
		fmt.Printf("  • No public group found, so /public/%s/:key/variants was not mounted.\n", names.Plural)
		fmt.Printf("    Regenerate the resource with --public to get one, then run this again.\n")
		return nil
	}

	// Its own handler instance rather than the variable above, which is
	// declared in the protected block further down the file. A public route
	// referring to it would be reading a variable that does not exist yet.
	route := fmt.Sprintf(
		"\t\t%[1]sPublicVariants := handlers.New%[2]sVariantHandler(db)\n"+
			"\t\tpublicAPI.GET(\"/%[3]s/:key/variants\", %[1]sPublicVariants.ListPublic)",
		names.Snake, names.Pascal, names.Plural)

	if err := injectBefore(path, "// grit:routes:public", route); err != nil {
		fmt.Printf("  Could not mount the public variant route: %v\n", err)
		return nil
	}
	manifest.Refresh(path)
	fmt.Printf("  ✓ GET /api/v1/public/%s/:key/variants (API key required)\n", names.Plural)
	return nil
}
