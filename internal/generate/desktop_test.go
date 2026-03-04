package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ── BuildNames ──────────────────────────────────────────────────────────────

func TestBuildNames(t *testing.T) {
	def := &ResourceDefinition{Name: "blog_post"}
	names := BuildNames(def)

	checks := map[string]string{
		"Pascal":       "BlogPost",
		"Camel":        "blogPost",
		"Snake":        "blog_post",
		"Kebab":        "blog-post",
		"Lower":        "blogpost",
		"PluralPascal": "BlogPosts",
		"PluralSnake":  "blog_posts",
		"PluralKebab":  "blog-posts",
	}

	for label, want := range checks {
		var got string
		switch label {
		case "Pascal":
			got = names.Pascal
		case "Camel":
			got = names.Camel
		case "Snake":
			got = names.Snake
		case "Kebab":
			got = names.Kebab
		case "Lower":
			got = names.Lower
		case "PluralPascal":
			got = names.PluralPascal
		case "PluralSnake":
			got = names.PluralSnake
		case "PluralKebab":
			got = names.PluralKebab
		}
		if got != want {
			t.Errorf("BuildNames(%q).%s = %q, want %q", "blog_post", label, got, want)
		}
	}
}

// ── Desktop Model ───────────────────────────────────────────────────────────

func TestWriteDesktopModel(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "internal", "models"), 0755)

	g := &DesktopGenerator{
		Root:   root,
		Module: "test-app",
		Definition: &ResourceDefinition{
			Name: "product",
			Fields: []Field{
				{Name: "name", Type: "string", Required: true},
				{Name: "price", Type: "float"},
				{Name: "published", Type: "bool"},
			},
		},
	}
	names := BuildNames(g.Definition)

	if err := g.writeDesktopModel(names); err != nil {
		t.Fatalf("writeDesktopModel: %v", err)
	}

	path := filepath.Join(root, "internal", "models", "product.go")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("model file not created: %v", err)
	}
	content := string(data)

	for _, want := range []string{
		"package models",
		"type Product struct",
		"Name",
		"Price",
		"Published",
		"gorm.DeletedAt",
		"time.Time",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("model missing %q", want)
		}
	}

	if strings.Contains(content, "<MODULE>") {
		t.Error("model still contains <MODULE> placeholder")
	}
}

func TestWriteDesktopModelWithSlug(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "internal", "models"), 0755)

	g := &DesktopGenerator{
		Root:   root,
		Module: "test-app",
		Definition: &ResourceDefinition{
			Name: "article",
			Fields: []Field{
				{Name: "title", Type: "string", Required: true},
				{Name: "slug", Type: "slug", SlugSource: "title"},
				{Name: "content", Type: "text"},
			},
		},
	}
	names := BuildNames(g.Definition)

	if err := g.writeDesktopModel(names); err != nil {
		t.Fatalf("writeDesktopModel: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(root, "internal", "models", "article.go"))
	if err != nil {
		t.Fatalf("model file not created: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "BeforeCreate") {
		t.Error("slug model missing BeforeCreate hook")
	}
	if !strings.Contains(content, "slugify") {
		t.Error("slug model missing slugify call")
	}
}

// ── Desktop Service ─────────────────────────────────────────────────────────

func TestWriteDesktopService(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "internal", "service"), 0755)

	g := &DesktopGenerator{
		Root:   root,
		Module: "test-app",
		Definition: &ResourceDefinition{
			Name: "product",
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "price", Type: "float"},
			},
		},
	}
	names := BuildNames(g.Definition)

	if err := g.writeDesktopService(names); err != nil {
		t.Fatalf("writeDesktopService: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(root, "internal", "service", "product.go"))
	if err != nil {
		t.Fatalf("service file not created: %v", err)
	}
	content := string(data)

	for _, want := range []string{
		"package service",
		"type ProductService struct",
		"NewProductService",
		"func (s *ProductService) List(",
		"func (s *ProductService) ListAll(",
		"func (s *ProductService) GetByID(",
		"func (s *ProductService) Create(",
		"func (s *ProductService) Update(",
		"func (s *ProductService) Delete(",
		"PaginatedResult",
		"test-app/internal/models",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("service missing %q", want)
		}
	}

	if strings.Contains(content, "<MODULE>") {
		t.Error("service still contains <MODULE> placeholder")
	}
}

// ── Desktop List Page ───────────────────────────────────────────────────────

func TestWriteDesktopListPage(t *testing.T) {
	root := t.TempDir()

	g := &DesktopGenerator{
		Root:   root,
		Module: "test-app",
		Definition: &ResourceDefinition{
			Name: "product",
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "price", Type: "float"},
				{Name: "published", Type: "bool"},
			},
		},
	}
	names := BuildNames(g.Definition)

	if err := g.writeDesktopListPage(names); err != nil {
		t.Fatalf("writeDesktopListPage: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(root, "frontend", "src", "pages", "products", "index.tsx"))
	if err != nil {
		t.Fatalf("list page not created: %v", err)
	}
	content := string(data)

	for _, want := range []string{
		"ProductsListPage",
		"GetProducts",
		"DeleteProduct",
		"ExportProductsPDF",
		"ExportProductsExcel",
		"useQuery",
		"useMutation",
		"wailsjs/go/main/App",
		"lucide-react",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("list page missing %q", want)
		}
	}
}

// ── Desktop Form Page ───────────────────────────────────────────────────────

func TestWriteDesktopFormPage(t *testing.T) {
	root := t.TempDir()

	g := &DesktopGenerator{
		Root:   root,
		Module: "test-app",
		Definition: &ResourceDefinition{
			Name: "product",
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "price", Type: "float"},
				{Name: "published", Type: "bool"},
			},
		},
	}
	names := BuildNames(g.Definition)

	if err := g.writeDesktopFormPage(names); err != nil {
		t.Fatalf("writeDesktopFormPage: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(root, "frontend", "src", "pages", "products", "form.tsx"))
	if err != nil {
		t.Fatalf("form page not created: %v", err)
	}
	content := string(data)

	for _, want := range []string{
		"ProductFormPage",
		"GetProduct",
		"CreateProduct",
		"UpdateProduct",
		"useParams",
		"isEditing",
		"handleSubmit",
		"type=\"number\"",
		"wailsjs/go/main/App",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("form page missing %q", want)
		}
	}
}

// ── Desktop Injection ───────────────────────────────────────────────────────

func TestDesktopInjectAll(t *testing.T) {
	root := setupDesktopProject(t)

	g := &DesktopGenerator{
		Root:   root,
		Module: "test-app",
		Definition: &ResourceDefinition{
			Name: "product",
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "price", Type: "float"},
			},
		},
	}
	names := BuildNames(g.Definition)

	if err := g.injectDesktopAll(names); err != nil {
		t.Fatalf("injectDesktopAll: %v", err)
	}

	// Verify db.go
	assertFileContains(t, filepath.Join(root, "internal", "db", "db.go"), "&models.Product{},")

	// Verify main.go
	mainContent := readFileStr(t, filepath.Join(root, "main.go"))
	if !strings.Contains(mainContent, "productSvc := service.NewProductService(database)") {
		t.Error("main.go missing service init")
	}
	if !strings.Contains(mainContent, "productSvc, ") {
		t.Error("main.go missing app arg")
	}

	// Verify app.go
	appContent := readFileStr(t, filepath.Join(root, "app.go"))
	if !strings.Contains(appContent, "product *service.ProductService") {
		t.Error("app.go missing field")
	}
	if !strings.Contains(appContent, "func (a *App) GetProducts(") {
		t.Error("app.go missing GetProducts method")
	}
	if !strings.Contains(appContent, "func (a *App) CreateProduct(") {
		t.Error("app.go missing CreateProduct method")
	}
	if !strings.Contains(appContent, "func (a *App) ExportProductsPDF(") {
		t.Error("app.go missing ExportProductsPDF method")
	}

	// Verify types.go
	assertFileContains(t, filepath.Join(root, "internal", "models", "types.go"), "type ProductInput struct")

	// Verify studio
	assertFileContains(t, filepath.Join(root, "cmd", "studio", "main.go"), "&models.Product{},")

	// Verify App.tsx
	appTsxContent := readFileStr(t, filepath.Join(root, "frontend", "src", "App.tsx"))
	if !strings.Contains(appTsxContent, "import ProductListPage") {
		t.Error("App.tsx missing import")
	}
	if !strings.Contains(appTsxContent, "path=\"/products\"") {
		t.Error("App.tsx missing route")
	}

	// Verify sidebar.tsx
	assertFileContains(t, filepath.Join(root, "frontend", "src", "components", "layout", "sidebar.tsx"),
		"path: \"/products\"")
}

// ── Desktop Inject + Remove Round-trip ──────────────────────────────────────

func TestDesktopInjectAndRemoveRoundtrip(t *testing.T) {
	root := setupDesktopProject(t)

	// Save original file contents
	originals := map[string]string{}
	files := []string{
		filepath.Join(root, "internal", "db", "db.go"),
		filepath.Join(root, "main.go"),
		filepath.Join(root, "app.go"),
		filepath.Join(root, "internal", "models", "types.go"),
		filepath.Join(root, "frontend", "src", "App.tsx"),
		filepath.Join(root, "frontend", "src", "components", "layout", "sidebar.tsx"),
	}
	for _, f := range files {
		data, _ := os.ReadFile(f)
		originals[f] = string(data)
	}

	// Inject
	g := &DesktopGenerator{
		Root:   root,
		Module: "test-app",
		Definition: &ResourceDefinition{
			Name: "product",
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "price", Type: "float"},
			},
		},
	}
	names := BuildNames(g.Definition)

	if err := g.writeDesktopModel(names); err != nil {
		t.Fatalf("writeDesktopModel: %v", err)
	}
	if err := g.writeDesktopService(names); err != nil {
		t.Fatalf("writeDesktopService: %v", err)
	}
	if err := g.injectDesktopAll(names); err != nil {
		t.Fatalf("injectDesktopAll: %v", err)
	}

	// Verify injection happened
	assertFileContains(t, filepath.Join(root, "internal", "db", "db.go"), "Product")

	// Now remove (call the inner logic, not RemoveDesktopResource which needs project detection)
	modelFile := filepath.Join(root, "internal", "models", "product.go")
	serviceFile := filepath.Join(root, "internal", "service", "product.go")

	os.Remove(modelFile)
	os.Remove(serviceFile)

	// Reverse injections
	dbFile := filepath.Join(root, "internal", "db", "db.go")
	removeLinesContaining(dbFile, "&models.Product{}")

	mainFile := filepath.Join(root, "main.go")
	removeLinesContaining(mainFile, "productSvc := service.NewProductService")
	removeInlineText(mainFile, "productSvc, ")

	appFile := filepath.Join(root, "app.go")
	removeLinesContaining(appFile, "product *service.ProductService")
	removeInlineText(appFile, "product *service.ProductService, ")
	removeInlineText(appFile, "product: product, ")
	// Remove each method block
	methodSigs := []string{
		"func (a *App) GetProducts(",
		"func (a *App) GetProduct(",
		"func (a *App) CreateProduct(",
		"func (a *App) UpdateProduct(",
		"func (a *App) DeleteProduct(",
		"func (a *App) ExportProductsPDF(",
		"func (a *App) ExportProductsExcel(",
	}
	for _, sig := range methodSigs {
		removeGoFuncBlock(appFile, sig)
	}
	removeLinesContaining(appFile, "// ── Products")

	typesFile := filepath.Join(root, "internal", "models", "types.go")
	removeLineBlock(typesFile, "type ProductInput struct {", "}")

	appTsx := filepath.Join(root, "frontend", "src", "App.tsx")
	removeLinesContaining(appTsx, "pages/products/")
	removeLinesContaining(appTsx, "path=\"/products\"")
	removeLinesContaining(appTsx, "path=\"/products/new\"")
	removeLinesContaining(appTsx, "path=\"/products/:id/edit\"")

	sidebarFile := filepath.Join(root, "frontend", "src", "components", "layout", "sidebar.tsx")
	removeLinesContaining(sidebarFile, "path: \"/products\"")

	// Verify model/service files are gone
	if _, err := os.Stat(modelFile); err == nil {
		t.Error("model file should be deleted")
	}
	if _, err := os.Stat(serviceFile); err == nil {
		t.Error("service file should be deleted")
	}

	// Verify injections are cleaned
	for _, f := range files {
		content := readFileStr(t, f)
		if strings.Contains(content, "Product") && !strings.Contains(originals[f], "Product") {
			// Find the lines containing Product for debug
			for i, line := range strings.Split(content, "\n") {
				if strings.Contains(line, "Product") {
					t.Errorf("file %s line %d still contains Product: %q", filepath.Base(f), i+1, line)
				}
			}
		}
	}
}

// ── Helpers ─────────────────────────────────────────────────────────────────

// setupDesktopProject creates a minimal desktop project with marker files for injection tests.
func setupDesktopProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	dirs := []string{
		filepath.Join(root, "internal", "db"),
		filepath.Join(root, "internal", "models"),
		filepath.Join(root, "internal", "service"),
		filepath.Join(root, "cmd", "studio"),
		filepath.Join(root, "frontend", "src", "components", "layout"),
		filepath.Join(root, "frontend", "src", "pages"),
	}
	for _, d := range dirs {
		os.MkdirAll(d, 0755)
	}

	// db.go with marker
	writeDesktopTestFile(t, filepath.Join(root, "internal", "db", "db.go"), `package db

func migrate() {
	db.AutoMigrate(
		&models.User{},
		// grit:models
	)
}
`)

	// main.go with markers
	writeDesktopTestFile(t, filepath.Join(root, "main.go"), `package main

func main() {
	authSvc := service.NewAuthService(database)
	exportSvc := service.NewExportService()
	// grit:service-init

	app := NewApp(authSvc, exportSvc, /* grit:app-args */)
}
`)

	// app.go with markers
	writeDesktopTestFile(t, filepath.Join(root, "app.go"), `package main

type App struct {
	auth   *service.AuthService
	export *service.ExportService
	// grit:fields
}

func NewApp(auth *service.AuthService, export *service.ExportService, /* grit:constructor-params */) *App {
	return &App{auth: auth, export: export, /* grit:constructor-assign */}
}

// grit:methods
`)

	// types.go with marker
	writeDesktopTestFile(t, filepath.Join(root, "internal", "models", "types.go"), `package models

type AuthResponse struct {
	User User
}

// grit:input-types
`)

	// studio main.go with marker
	writeDesktopTestFile(t, filepath.Join(root, "cmd", "studio", "main.go"), `package main

func main() {
	db.AutoMigrate(
		&models.User{},
		// grit:studio-models
	)
}
`)

	// App.tsx with markers
	writeDesktopTestFile(t, filepath.Join(root, "frontend", "src", "App.tsx"), `import DashboardPage from "./pages/dashboard";
// grit:page-imports

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<DashboardPage />} />
      {/* grit:routes */}
    </Routes>
  );
}
`)

	// sidebar.tsx with markers
	writeDesktopTestFile(t, filepath.Join(root, "frontend", "src", "components", "layout", "sidebar.tsx"), `import {
  Home,
  // grit:nav-icons
} from "lucide-react";

const navItems = [
  { label: "Dashboard", path: "/", icon: Home },
  // grit:nav
];
`)

	return root
}

func writeDesktopTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}

func readFileStr(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(data)
}

func assertFileContains(t *testing.T, path, want string) {
	t.Helper()
	content := readFileStr(t, path)
	if !strings.Contains(content, want) {
		t.Errorf("file %s missing %q", filepath.Base(path), want)
	}
}
