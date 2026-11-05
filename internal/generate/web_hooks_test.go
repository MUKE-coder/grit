package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// grit#146: the web hook was ninety-two lines of React Query per resource, the
// same file with the names changed, and every mutation input was typed
// Record<string, unknown> although CreateProductInput was being generated two
// directories away in packages/shared.
//
// So the compiler could not tell you that you sent `title` where the API wanted
// `name`, and a fix to the hook shape was one edit per resource per project.
func TestTheWebHookIsThinAndTyped(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	// The generator only writes web hooks when the directory is there.
	hooks := filepath.Join(root, "apps", "web", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}

	def, err := ParseInlineFields("Product", "name:string,price:float,active:bool")
	if err != nil {
		t.Fatal(err)
	}
	if err := newTestGenerator(root, module, def).Run(); err != nil {
		t.Fatalf("Generator.Run(): %v", err)
	}

	path := filepath.Join(hooks, "use-products.ts")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the generated hook: %v", err)
	}
	src := string(data)

	// Typed, from the shared package, which is generated from the Go model.
	for _, want := range []string{
		`import type { Product } from "@repo/shared/types";`,
		`import type { CreateProductInput, UpdateProductInput } from "@repo/shared/schemas";`,
		`import { resourceHooks } from "./use-resource";`,
		`resourceHooks<Product, CreateProductInput, UpdateProductInput>(`,
		`"/api/products",`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the generated hook is missing %q", want)
		}
	}

	// The thing that made the hook useless as a type: an input the compiler
	// cannot check.
	if strings.Contains(src, "Record<string, unknown>") {
		t.Error("the mutation input is still Record<string, unknown>, so the compiler " +
			"cannot tell you that you sent `title` where the API wanted `name`")
	}

	// The exported names are unchanged, because somebody's pages import them
	// and a rename would be a breaking change to buy a tidier shape.
	for _, want := range []string{
		"export const useProducts =",
		"export const useGetProduct =",
		"export const useCreateProduct =",
		"export const useUpdateProduct =",
		"export const useDeleteProduct =",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("%q is gone: every page importing it stops compiling", want)
		}
	}

	// Thin. Not an exact count, which would fail on a comment edit, but the
	// difference between three lines over a generic and ninety-two of React
	// Query is not subtle.
	if lines := strings.Count(src, "\n"); lines > 40 {
		t.Errorf("the generated hook is %d lines; it should be a thin wrapper over "+
			"use-resource.ts, not a copy of the hook bodies", lines)
	}

	// And none of the bodies may come back into the per-resource file.
	for _, gone := range []string{"useQueryClient(", "new URLSearchParams(", "invalidateQueries("} {
		if strings.Contains(src, gone) {
			t.Errorf("the generated hook contains %q, so a hook body has been copied "+
				"back into the per-resource file", gone)
		}
	}
}
