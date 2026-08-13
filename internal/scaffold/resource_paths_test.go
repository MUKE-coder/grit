package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFileT(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

func read(t *testing.T, path string) string {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(body)
}

// The migration moves files someone wrote by hand, so the thing worth proving
// is that everything arrives and the imports still point at it.
func TestMigrateResourceLayoutMovesFilesAndFixesImports(t *testing.T) {
	admin := t.TempDir()
	resources := filepath.Join(admin, "resources")

	writeFileT(t, filepath.Join(resources, "products.ts"), "export const productResource = {}\n")
	writeFileT(t, filepath.Join(resources, "products.custom.tsx"), "// hand written\n")
	writeFileT(t, filepath.Join(resources, "users.ts"), "export const usersResource = {}\n")
	// No overlay for users: a resource that never had one must still move.
	writeFileT(t, filepath.Join(resources, "index.ts"), strings.Join([]string{
		`import { productResource } from "./products";`,
		`import { usersResource } from "./users";`,
		`export const resources = [productResource, usersResource];`,
	}, "\n"))
	// A page importing by alias rather than by relative path.
	writeFileT(t, filepath.Join(admin, "app", "resources", "products", "page.tsx"),
		`import { productResource } from "@/resources/products";`)

	moved, err := MigrateResourceLayout(resources)
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if moved != 2 {
		t.Errorf("moved %d resources, want 2", moved)
	}

	for _, want := range []string{
		"products/products.ts",
		"products/products.custom.tsx",
		"users/users.ts",
	} {
		if !isFile(filepath.Join(resources, filepath.FromSlash(want))) {
			t.Errorf("%s was not created", want)
		}
	}
	for _, gone := range []string{"products.ts", "products.custom.tsx", "users.ts"} {
		if isFile(filepath.Join(resources, gone)) {
			t.Errorf("%s is still in the root", gone)
		}
	}
	if read(t, filepath.Join(resources, "products", "products.custom.tsx")) != "// hand written\n" {
		t.Error("the overlay's contents changed in the move")
	}

	index := read(t, filepath.Join(resources, "index.ts"))
	if !strings.Contains(index, `from "./products/products"`) {
		t.Errorf("registry import not rewritten:\n%s", index)
	}
	if strings.Contains(index, `from "./products"`+";") {
		t.Errorf("old registry import still present:\n%s", index)
	}

	page := read(t, filepath.Join(admin, "app", "resources", "products", "page.tsx"))
	if !strings.Contains(page, `"@/resources/products/products"`) {
		t.Errorf("alias import not rewritten: %s", page)
	}
}

// Running it twice must not produce products/products/products.
func TestMigrateResourceLayoutIsIdempotent(t *testing.T) {
	admin := t.TempDir()
	resources := filepath.Join(admin, "resources")
	writeFileT(t, filepath.Join(resources, "products.ts"), "x\n")
	writeFileT(t, filepath.Join(resources, "index.ts"), `import { p } from "./products";`)
	writeFileT(t, filepath.Join(admin, "page.tsx"), `import { p } from "@/resources/products";`)

	for range 3 {
		if _, err := MigrateResourceLayout(resources); err != nil {
			t.Fatalf("migrate: %v", err)
		}
	}

	if isFile(filepath.Join(resources, "products", "products", "products.ts")) {
		t.Error("the migration nested a folder inside itself")
	}
	index := read(t, filepath.Join(resources, "index.ts"))
	if strings.Count(index, "products/products") != 1 {
		t.Errorf("import rewritten more than once: %s", index)
	}
	page := read(t, filepath.Join(admin, "page.tsx"))
	if strings.Count(page, "products/products") != 1 {
		t.Errorf("alias rewritten more than once: %s", page)
	}
}

// index.ts is the registry, not a resource. Moving it would break every import
// in the app at once.
func TestMigrateResourceLayoutLeavesTheRegistryAlone(t *testing.T) {
	resources := t.TempDir()
	writeFileT(t, filepath.Join(resources, "index.ts"), "export const resources = [];")
	writeFileT(t, filepath.Join(resources, "registry.ts"), "// legacy name for the same thing")

	moved, err := MigrateResourceLayout(resources)
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if moved != 0 {
		t.Errorf("moved %d, want 0", moved)
	}
	if !isFile(filepath.Join(resources, "index.ts")) {
		t.Error("index.ts was moved")
	}
	if !isFile(filepath.Join(resources, "registry.ts")) {
		t.Error("registry.ts was moved")
	}
}

func TestFindResourceDefPrefersTheFolderLayout(t *testing.T) {
	resources := t.TempDir()
	flat := filepath.Join(resources, "products.ts")
	nested := filepath.Join(resources, "products", "products.ts")

	writeFileT(t, flat, "flat")
	if got := FindResourceDef(resources, "products"); got != flat {
		t.Errorf("flat only: got %s", got)
	}

	writeFileT(t, nested, "nested")
	if got := FindResourceDef(resources, "products"); got != nested {
		t.Errorf("both present: got %s, want the nested one", got)
	}

	if got := FindResourceDef(resources, "nothing"); got != "" {
		t.Errorf("missing resource: got %s, want empty", got)
	}
}
