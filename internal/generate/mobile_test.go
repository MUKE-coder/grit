package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mobileGen builds a generator rooted at a temp project that has an apps/expo
// directory, so hasMobileApp() is true and the mobile writers run.
func mobileGen(t *testing.T, name, fields string) (*Generator, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "apps", "expo"), 0o755); err != nil {
		t.Fatalf("mkdir expo: %v", err)
	}
	def, err := ParseInlineFields(name, fields)
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	return newTestGenerator(root, "store/apps/api", def), root
}

func TestMobile_HasMobileApp(t *testing.T) {
	g, _ := mobileGen(t, "Product", "name:string")
	if !g.hasMobileApp() {
		t.Fatal("hasMobileApp() = false, want true when apps/expo exists")
	}

	// A project with no apps/expo must not trigger mobile generation.
	bare := &Generator{Root: t.TempDir(), Definition: &ResourceDefinition{Name: "Product"}}
	if bare.hasMobileApp() {
		t.Fatal("hasMobileApp() = true for project without apps/expo")
	}
}

func TestMobile_GeneratesFilesAndRoutes(t *testing.T) {
	g, root := mobileGen(t, "Product", "name:string,slug:slug,price:int,category:belongs_to:Category")
	if err := g.writeMobileFiles(g.Names()); err != nil {
		t.Fatalf("writeMobileFiles: %v", err)
	}

	expo := filepath.Join(root, "apps", "expo")
	for _, rel := range []string{
		"hooks/use-products.ts",
		"app/products/index.tsx",
		"app/products/[id].tsx",
		"app/products/new.tsx",
		"app/products/edit/[id].tsx",
		"components/resource-forms/products-form.tsx",
		"components/ui/screen-header.tsx",
	} {
		if _, err := os.Stat(filepath.Join(expo, rel)); err != nil {
			t.Errorf("expected generated file %s: %v", rel, err)
		}
	}
}

func TestMobile_HookContent(t *testing.T) {
	g, root := mobileGen(t, "Product", "name:string,category:belongs_to:Category,thumbnail:file:image")
	if err := g.writeMobileFiles(g.Names()); err != nil {
		t.Fatalf("writeMobileFiles: %v", err)
	}
	hook := read(t, filepath.Join(root, "apps", "expo", "hooks", "use-products.ts"))

	wants := []string{
		"export function useProducts(",  // infinite list hook
		"useInfiniteQuery",              // pagination
		"export function useProduct(id", // single-item hook
		"useCreateProduct",              // mutation
		`api.get("/products?"`,          // list endpoint (mobile api has no /api prefix)
		"category_id: string;",          // belongs_to FK
		"category?: any;",               // preloaded relation
		"type FileRef =",                // inline FileRef for the file field
		"thumbnail: FileRef | null;",    // file field typed
	}
	for _, w := range wants {
		if !strings.Contains(hook, w) {
			t.Errorf("hook missing %q\n---\n%s", w, hook)
		}
	}
}

func TestMobile_ScreensContent(t *testing.T) {
	g, root := mobileGen(t, "Product", "name:string,description:text,thumbnail:file:image,category:belongs_to:Category")
	if err := g.writeMobileFiles(g.Names()); err != nil {
		t.Fatalf("writeMobileFiles: %v", err)
	}
	expo := filepath.Join(root, "apps", "expo")

	list := read(t, filepath.Join(expo, "app", "products", "index.tsx"))
	for _, w := range []string{
		`import { ScreenHeader } from "@/components/ui/screen-header";`,
		"useProducts(search)",                 // wired to the hook
		`router.push("/products/" + item.id)`, // navigates to detail
		"item.thumbnail?.url",                 // image thumbnail from file field
		"item.description",                    // subtitle from text field
		"showBack",                            // back button
	} {
		if !strings.Contains(list, w) {
			t.Errorf("list screen missing %q", w)
		}
	}

	detail := read(t, filepath.Join(expo, "app", "products", "[id].tsx"))
	for _, w := range []string{
		"useProduct(id)",
		`<Row label="Category" value={(item.category`, // belongs_to renders related name w/ FK fallback
		"item.thumbnail?.url",                         // hero image
		`title="Product"`,                             // singular header
	} {
		if !strings.Contains(detail, w) {
			t.Errorf("detail screen missing %q", w)
		}
	}
}

// TestMobile_TitleAndImageHeuristics covers the field-picking helpers that
// decide the row title and image source.
func TestMobile_TitleAndImageHeuristics(t *testing.T) {
	// "title" beats a leading string field; a file field wins the image slot.
	g, _ := mobileGen(t, "Article", "headline:string,title:string,cover:file:image")
	if got := g.mobileTitleField(); got != "title" {
		t.Errorf("mobileTitleField() = %q, want title", got)
	}
	if got := g.mobileImageExpr("item"); got != "item.cover?.url" {
		t.Errorf("mobileImageExpr() = %q, want item.cover?.url", got)
	}

	// No name/title/string -> falls back to id; URL-shaped string is the image.
	g2, _ := mobileGen(t, "Banner", "image:string,active:bool")
	if got := g2.mobileTitleField(); got != "image" {
		t.Errorf("mobileTitleField() = %q, want image (first string)", got)
	}
	if got := g2.mobileImageExpr("item"); got != "item.image" {
		t.Errorf("mobileImageExpr() = %q, want item.image", got)
	}
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}
