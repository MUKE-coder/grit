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

	// The shared half, written once however many resources have variants.
	optionsPath := filepath.Join(apiRoot, "internal", "models", "option.go")
	if fileExists(optionsPath) {
		fmt.Println("  • internal/models/option.go exists, left alone (options are shared)")
	} else {
		if err := writeFileWithDirs(optionsPath, scaffold.APIOptionModelGo(module)); err != nil {
			return fmt.Errorf("writing option models: %w", err)
		}
		green.Println("  ✓ internal/models/option.go (Option, OptionValue)")
	}

	// The per-resource half.
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "models", names.Snake+"_variant.go"):
			scaffold.APIVariantModelGo(module, names.Pascal, names.Snake),
		filepath.Join(apiRoot, "internal", "services", names.Snake+"_variants.go"):
			scaffold.APIVariantServiceGo(module, names.Pascal, names.Snake, names.Plural),
		filepath.Join(apiRoot, "internal", "services", names.Snake+"_variants_test.go"):
			scaffold.APIVariantServiceTestGo(module, names.Pascal, names.Snake),
		filepath.Join(apiRoot, "internal", "handlers", names.Snake+"_variant.go"):
			scaffold.APIVariantHandlerGo(module, names.Pascal, names.Snake, names.Plural),
	}
	for path, body := range files {
		if err := writeFileWithDirs(path, body); err != nil {
			return fmt.Errorf("writing %s: %w", filepath.Base(path), err)
		}
		green.Printf("  ✓ %s\n", strings.TrimPrefix(filepath.ToSlash(strings.TrimPrefix(path, root)), "/"))
	}

	if err := registerVariantModels(apiRoot, names); err != nil {
		return err
	}
	if err := mountVariantRoutes(apiRoot, names); err != nil {
		return err
	}

	fmt.Println()
	green.Println("  Variants installed.")
	fmt.Println()
	fmt.Println("  Next:")
	fmt.Println("    grit migrate                       # create the five tables")
	fmt.Printf("    Options live at /api/v1/options, and the matrix at\n")
	fmt.Printf("    /api/v1/%s/:id/variants.\n", names.Plural)
	fmt.Println()
	return nil
}

type variantNameSet struct {
	Pascal string
	Snake  string
	Plural string
}

func variantNames(resource string) variantNameSet {
	pascal := toPascalCase(resource)
	return variantNameSet{
		Pascal: pascal,
		Snake:  toSnakeCase(pascal),
		Plural: Pluralize(toSnakeCase(pascal)),
	}
}

// registerVariantModels adds the new models to AutoMigrate and GORM Studio.
//
// Without this the tables are never created and the first request fails on a
// relation that does not exist, which reads as a bug in Grit rather than a
// migration nobody ran.
func registerVariantModels(apiRoot string, names variantNameSet) error {
	path := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(data)

	toAdd := []string{"Option", "OptionValue", names.Pascal + "Option", names.Pascal + "Variant"}
	var missing []string
	for _, model := range toAdd {
		// Anchored on the marker's own line so "Option" does not match
		// "ProductOption" and decide it is already registered.
		if !strings.Contains(content, "&models."+model+"{},") {
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
func mountVariantRoutes(apiRoot string, names variantNameSet) error {
	path := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	if strings.Contains(string(data), names.Snake+"VariantHandler") {
		fmt.Println("  • Routes already mounted")
		return nil
	}

	block := fmt.Sprintf(`		%[1]sVariantHandler := handlers.New%[2]sVariantHandler(db)
		protected.GET("/options", %[1]sVariantHandler.ListOptions)
		protected.POST("/options", %[1]sVariantHandler.CreateOption)
		protected.POST("/options/:id/values", %[1]sVariantHandler.CreateOptionValue)
		protected.DELETE("/option-values/:id", %[1]sVariantHandler.DeleteOptionValue)
		protected.GET("/%[3]s/:id/variants", %[1]sVariantHandler.List)
		protected.POST("/%[3]s/:id/variants/generate", %[1]sVariantHandler.GenerateMatrix)
		protected.PATCH("/variants/:id", %[1]sVariantHandler.Update)
		protected.PUT("/%[3]s/:id/options", %[1]sVariantHandler.SetOptions)`,
		names.Snake, names.Pascal, names.Plural)

	if err := injectBefore(path, "// grit:routes:custom", block); err != nil {
		fmt.Printf("  Could not mount the routes automatically: %v\n", err)
		return nil
	}
	manifest.Refresh(path)
	fmt.Printf("  ✓ Mounted /options and /%s/:id/variants\n", names.Plural)
	return nil
}
