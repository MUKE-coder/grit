package scaffold

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strconv"
	"strings"
	"testing"
)

// assertImportsWhatItUses fails when src uses pkg.Name for a package it does
// not import. The parser leaves such a pkg unresolved, the way it leaves an
// imported package's name, so every selector on an unresolved identifier has
// to name an import.
func assertImportsWhatItUses(t *testing.T, name, src string) {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), name, src, parser.SkipObjectResolution)
	if err != nil {
		t.Errorf("%s does not parse: %v", name, err)
		return
	}
	imported := map[string]bool{}
	for _, spec := range file.Imports {
		p, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			t.Fatal(err)
		}
		local := path.Base(p)
		if spec.Name != nil {
			local = spec.Name.Name
		}
		imported[local] = true
	}
	// Local names that shadow nothing: receivers, parameters and variables
	// are lowercase and never collide with these package names.
	for _, pkg := range []string{"respond", "files", "models", "paginate", "services", "http", "gin", "gorm", "strings", "errors", "fmt", "time", "sort"} {
		used := false
		ast.Inspect(file, func(n ast.Node) bool {
			if sel, ok := n.(*ast.SelectorExpr); ok {
				if id, ok := sel.X.(*ast.Ident); ok && id.Name == pkg {
					used = true
				}
			}
			return !used
		})
		if used && !imported[pkg] {
			t.Errorf("%s uses %s without importing it", name, pkg)
		}
	}
}

func TestVariantTemplatesImportWhatTheyUse(t *testing.T) {
	const module = "shop/apps/api"
	for name, src := range map[string]string{
		"handler":         APIVariantHandlerGo(module, "Product", "product", "products"),
		"public":          APIVariantPublicGo(module, "Product", "product", "products", true, true),
		"public, no slug": APIVariantPublicGo(module, "Product", "product", "products", false, false),
		"option public":   APIOptionPublicGo(module),
		"service":         APIVariantServiceGo(module, "Product", "product", "products"),
		"models":          APIVariantModelGo(module, "Product", "product"),
		"option models":   APIOptionModelGo(module),
		"seeder":          APIVariantSeederGo(module, "Product", "product", "products", true),
	} {
		assertImportsWhatItUses(t, name, src)
	}
}

func TestVariantPublicImportRepair(t *testing.T) {
	current := APIVariantPublicGo("shop/apps/api", "Product", "product", "products", true, true)
	if out, changes, _ := repairVariantPublicImportSource(current); out != current || len(changes) != 0 {
		t.Error("a current handler was changed")
	}
	old := strings.Replace(current, "\t\"shop/apps/api/internal/respond\"\n", "", 1)
	fixed, changes, warnings := repairVariantPublicImportSource(old)
	if fixed != current || len(changes) != 1 || len(warnings) != 0 {
		t.Errorf("the repair did not restore the template (%v, %v)", changes, warnings)
	}
}

// topLevelNames is every name a file declares at package level.
func topLevelNames(t *testing.T, name, src string) []string {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), name, src, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("%s does not parse: %v", name, err)
	}
	var names []string
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil {
				names = append(names, d.Name.Name)
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					names = append(names, s.Name.Name)
				case *ast.ValueSpec:
					for _, n := range s.Names {
						names = append(names, n.Name)
					}
				}
			}
		}
	}
	return names
}

// Two resources with variants share the services, database, handlers and
// models packages, so nothing either writes may be declared by both.
func TestTwoResourcesWithVariantsDeclareNothingTwice(t *testing.T) {
	const module = "shop/apps/api"
	perPackage := func(pascal, snake, plural string) map[string][]string {
		return map[string][]string{
			"services": {APIVariantServiceGo(module, pascal, snake, plural), APIVariantServiceTestGo(module, pascal, snake)},
			"database": {APIVariantSeederGo(module, pascal, snake, plural, true)},
			"handlers": {APIVariantHandlerGo(module, pascal, snake, plural), APIVariantPublicGo(module, pascal, snake, plural, true, true)},
			"models":   {APIVariantModelGo(module, pascal, snake)},
		}
	}
	first, second := perPackage("Product", "product", "products"), perPackage("Gadget", "gadget", "gadgets")
	for pkg, files := range first {
		declared := map[string]bool{}
		for i, src := range files {
			for _, n := range topLevelNames(t, pkg+strconv.Itoa(i), src) {
				declared[n] = true
			}
		}
		for i, src := range second[pkg] {
			for _, n := range topLevelNames(t, pkg+strconv.Itoa(i), src) {
				if declared[n] {
					t.Errorf("%s: both resources declare %s", pkg, n)
				}
			}
		}
	}
}
