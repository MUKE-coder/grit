package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateDesktopDirectories(t *testing.T) {
	root := t.TempDir()

	if err := createDesktopDirectories(root); err != nil {
		t.Fatalf("createDesktopDirectories error: %v", err)
	}

	required := []string{
		filepath.Join(root, "internal", "config"),
		filepath.Join(root, "internal", "db"),
		filepath.Join(root, "internal", "models"),
		filepath.Join(root, "internal", "service"),
		filepath.Join(root, "cmd", "studio"),
		filepath.Join(root, "build"),
		filepath.Join(root, "frontend", "src", "components", "ui"),
		filepath.Join(root, "frontend", "src", "components", "layout"),
		filepath.Join(root, "frontend", "src", "routes", "_layout"),
		filepath.Join(root, "frontend", "src", "hooks"),
		filepath.Join(root, "frontend", "src", "lib"),
	}

	for _, dir := range required {
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("expected directory %s was not created", dir)
		}
	}
}

func TestDesktopRootFiles(t *testing.T) {
	root := t.TempDir()
	opts := DesktopOptions{ProjectName: "test-desktop"}

	if err := createDesktopDirectories(root); err != nil {
		t.Fatalf("createDesktopDirectories: %v", err)
	}
	if err := writeDesktopRootFiles(root, opts); err != nil {
		t.Fatalf("writeDesktopRootFiles: %v", err)
	}

	expected := []string{
		"wails.json",
		"go.mod",
		".gitignore",
		".env",
		".env.example",
		"README.md",
	}

	for _, f := range expected {
		path := filepath.Join(root, f)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s was not created", f)
		}
	}

	// Verify go.mod has correct module name
	goMod, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}
	if !strings.Contains(string(goMod), "module test-desktop") {
		t.Error("go.mod does not contain correct module name")
	}
	if !strings.Contains(string(goMod), "wailsapp/wails") {
		t.Error("go.mod does not contain wails dependency")
	}

	// Verify wails.json has project name
	wailsJSON, err := os.ReadFile(filepath.Join(root, "wails.json"))
	if err != nil {
		t.Fatalf("reading wails.json: %v", err)
	}
	if !strings.Contains(string(wailsJSON), "test-desktop") {
		t.Error("wails.json does not contain project name")
	}
}

func TestDesktopGoFiles(t *testing.T) {
	root := t.TempDir()
	opts := DesktopOptions{ProjectName: "test-desktop"}

	if err := createDesktopDirectories(root); err != nil {
		t.Fatalf("createDesktopDirectories: %v", err)
	}
	if err := writeDesktopGoFiles(root, opts); err != nil {
		t.Fatalf("writeDesktopGoFiles: %v", err)
	}

	// Check main.go and app.go exist
	for _, f := range []string{"main.go", "app.go"} {
		path := filepath.Join(root, f)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s was not created", f)
		}
	}

	// Verify module substitution in main.go
	mainGo, err := os.ReadFile(filepath.Join(root, "main.go"))
	if err != nil {
		t.Fatalf("reading main.go: %v", err)
	}
	if strings.Contains(string(mainGo), "<MODULE>") {
		t.Error("main.go still contains <MODULE> placeholder")
	}
	if !strings.Contains(string(mainGo), "test-desktop/internal/config") {
		t.Error("main.go does not contain substituted module path")
	}

	// Verify app.go has bound methods
	appGo, err := os.ReadFile(filepath.Join(root, "app.go"))
	if err != nil {
		t.Fatalf("reading app.go: %v", err)
	}
	for _, method := range []string{"Login", "Register", "GetBlogs", "CreateBlog", "GetContacts", "ExportBlogsPDF", "MinimiseWindow", "CloseApp"} {
		if !strings.Contains(string(appGo), method) {
			t.Errorf("app.go missing method: %s", method)
		}
	}
}

func TestDesktopGoInternalFiles(t *testing.T) {
	root := t.TempDir()
	opts := DesktopOptions{ProjectName: "test-desktop"}

	if err := createDesktopDirectories(root); err != nil {
		t.Fatalf("createDesktopDirectories: %v", err)
	}
	if err := writeDesktopGoInternalFiles(root, opts); err != nil {
		t.Fatalf("writeDesktopGoInternalFiles: %v", err)
	}

	expected := []string{
		filepath.Join("internal", "config", "config.go"),
		filepath.Join("internal", "db", "db.go"),
		filepath.Join("internal", "models", "user.go"),
		filepath.Join("internal", "models", "blog.go"),
		filepath.Join("internal", "models", "contact.go"),
		filepath.Join("internal", "models", "types.go"),
	}

	for _, f := range expected {
		path := filepath.Join(root, f)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s was not created", f)
		}
	}

	// Verify no <MODULE> placeholders remain
	for _, f := range expected {
		content, err := os.ReadFile(filepath.Join(root, f))
		if err != nil {
			continue
		}
		if strings.Contains(string(content), "<MODULE>") {
			t.Errorf("%s still contains <MODULE> placeholder", f)
		}
	}
}

func TestDesktopGoServiceFiles(t *testing.T) {
	root := t.TempDir()
	opts := DesktopOptions{ProjectName: "test-desktop"}

	if err := createDesktopDirectories(root); err != nil {
		t.Fatalf("createDesktopDirectories: %v", err)
	}
	if err := writeDesktopGoServiceFiles(root, opts); err != nil {
		t.Fatalf("writeDesktopGoServiceFiles: %v", err)
	}

	expected := []string{
		filepath.Join("internal", "service", "auth.go"),
		filepath.Join("internal", "service", "blog.go"),
		filepath.Join("internal", "service", "contact.go"),
		filepath.Join("internal", "service", "export.go"),
	}

	for _, f := range expected {
		path := filepath.Join(root, f)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s was not created", f)
		}
	}

	// Verify auth service has key functions
	authGo, err := os.ReadFile(filepath.Join(root, "internal", "service", "auth.go"))
	if err != nil {
		t.Fatalf("reading auth.go: %v", err)
	}
	for _, fn := range []string{"NewAuthService", "Register", "Login", "GetUser"} {
		if !strings.Contains(string(authGo), fn) {
			t.Errorf("auth.go missing function: %s", fn)
		}
	}

	// Verify export service has PDF and Excel
	exportGo, err := os.ReadFile(filepath.Join(root, "internal", "service", "export.go"))
	if err != nil {
		t.Fatalf("reading export.go: %v", err)
	}
	for _, fn := range []string{"ExportPDF", "ExportExcel", "gofpdf", "excelize"} {
		if !strings.Contains(string(exportGo), fn) {
			t.Errorf("export.go missing: %s", fn)
		}
	}
}

func TestDesktopFrontendFiles(t *testing.T) {
	root := t.TempDir()
	opts := DesktopOptions{ProjectName: "test-desktop"}

	if err := createDesktopDirectories(root); err != nil {
		t.Fatalf("createDesktopDirectories: %v", err)
	}

	writers := []struct {
		name string
		fn   func(string, DesktopOptions) error
	}{
		{"config", writeDesktopFrontendConfig},
		{"app", writeDesktopFrontendAppFiles},
		{"layout", writeDesktopFrontendLayoutFiles},
		{"ui", writeDesktopFrontendUIFiles},
		{"pages", writeDesktopFrontendPageFiles},
		{"table-form", writeDesktopFrontendTableFormFiles},
	}

	for _, w := range writers {
		if err := w.fn(root, opts); err != nil {
			t.Fatalf("write %s: %v", w.name, err)
		}
	}

	// Verify key frontend files exist
	expected := []string{
		filepath.Join("frontend", "package.json"),
		filepath.Join("frontend", "vite.config.ts"),
		filepath.Join("frontend", "tailwind.config.ts"),
		filepath.Join("frontend", "tsconfig.json"),
		filepath.Join("frontend", "index.html"),
		filepath.Join("frontend", "src", "main.tsx"),
		filepath.Join("frontend", "src", "index.css"),
		filepath.Join("frontend", "src", "lib", "utils.ts"),
		filepath.Join("frontend", "src", "lib", "query-client.ts"),
		filepath.Join("frontend", "src", "hooks", "use-auth.ts"),
		filepath.Join("frontend", "src", "hooks", "use-theme.ts"),
		filepath.Join("frontend", "src", "components", "layout", "title-bar.tsx"),
		filepath.Join("frontend", "src", "components", "layout", "sidebar.tsx"),
		filepath.Join("frontend", "src", "components", "layout", "draggable-window.tsx"),
		filepath.Join("frontend", "src", "components", "ui", "button.tsx"),
		filepath.Join("frontend", "src", "components", "ui", "input.tsx"),
		filepath.Join("frontend", "src", "components", "ui", "card.tsx"),
		filepath.Join("frontend", "src", "components", "data-table.tsx"),
		filepath.Join("frontend", "src", "components", "form-builder.tsx"),
		filepath.Join("frontend", "src", "routes", "__root.tsx"),
		filepath.Join("frontend", "src", "routes", "login.tsx"),
		filepath.Join("frontend", "src", "routes", "register.tsx"),
		filepath.Join("frontend", "src", "routes", "_layout.tsx"),
		filepath.Join("frontend", "src", "routes", "_layout", "index.tsx"),
		filepath.Join("frontend", "src", "routes", "_layout", "blogs.index.tsx"),
		filepath.Join("frontend", "src", "routes", "_layout", "blogs.new.tsx"),
		filepath.Join("frontend", "src", "routes", "_layout", "blogs.$id.edit.tsx"),
		filepath.Join("frontend", "src", "routes", "_layout", "contacts.index.tsx"),
		filepath.Join("frontend", "src", "routes", "_layout", "contacts.new.tsx"),
		filepath.Join("frontend", "src", "routes", "_layout", "contacts.$id.edit.tsx"),
	}

	for _, f := range expected {
		path := filepath.Join(root, f)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s was not created", f)
		}
	}

	// Verify index.css has Grit theme
	css, err := os.ReadFile(filepath.Join(root, "frontend", "src", "index.css"))
	if err != nil {
		t.Fatalf("reading index.css: %v", err)
	}
	for _, v := range []string{"--background", "--accent", "--foreground", "--bg-elevated", "#6c5ce7"} {
		if !strings.Contains(string(css), v) {
			t.Errorf("index.css missing CSS variable: %s", v)
		}
	}

	// Verify title-bar has wails drag property
	titleBar, err := os.ReadFile(filepath.Join(root, "frontend", "src", "components", "layout", "title-bar.tsx"))
	if err != nil {
		t.Fatalf("reading title-bar.tsx: %v", err)
	}
	if !strings.Contains(string(titleBar), "wails-draggable") {
		t.Error("title-bar.tsx missing --wails-draggable CSS property")
	}

	// Verify main.tsx uses TanStack Router with hash history
	mainTsx, err := os.ReadFile(filepath.Join(root, "frontend", "src", "main.tsx"))
	if err != nil {
		t.Fatalf("reading main.tsx: %v", err)
	}
	if !strings.Contains(string(mainTsx), "createRouter") {
		t.Error("main.tsx does not use createRouter")
	}
	if !strings.Contains(string(mainTsx), "createHashHistory") {
		t.Error("main.tsx does not use createHashHistory")
	}

	// Verify __root.tsx uses createRootRoute
	rootRoute, err := os.ReadFile(filepath.Join(root, "frontend", "src", "routes", "__root.tsx"))
	if err != nil {
		t.Fatalf("reading __root.tsx: %v", err)
	}
	if !strings.Contains(string(rootRoute), "createRootRoute") {
		t.Error("__root.tsx does not use createRootRoute")
	}
}

func TestDesktopMarkers(t *testing.T) {
	root := t.TempDir()
	opts := DesktopOptions{ProjectName: "test-desktop"}

	if err := createDesktopDirectories(root); err != nil {
		t.Fatalf("createDesktopDirectories: %v", err)
	}
	if err := writeDesktopGoFiles(root, opts); err != nil {
		t.Fatalf("writeDesktopGoFiles: %v", err)
	}
	if err := writeDesktopGoInternalFiles(root, opts); err != nil {
		t.Fatalf("writeDesktopGoInternalFiles: %v", err)
	}
	if err := writeDesktopFrontendAppFiles(root, opts); err != nil {
		t.Fatalf("writeDesktopFrontendAppFiles: %v", err)
	}
	if err := writeDesktopFrontendLayoutFiles(root, opts); err != nil {
		t.Fatalf("writeDesktopFrontendLayoutFiles: %v", err)
	}
	if err := writeDesktopStudioFiles(root, opts); err != nil {
		t.Fatalf("writeDesktopStudioFiles: %v", err)
	}

	markers := map[string][]string{
		filepath.Join(root, "app.go"): {
			"// grit:imports",
			"// grit:fields",
			"/* grit:constructor-params */",
			"/* grit:constructor-assign */",
			"// grit:methods",
		},
		filepath.Join(root, "main.go"): {
			"// grit:service-init",
			"/* grit:app-args */",
		},
		filepath.Join(root, "internal", "db", "db.go"): {
			"// grit:models",
		},
		filepath.Join(root, "internal", "models", "types.go"): {
			"// grit:input-types",
		},
		filepath.Join(root, "frontend", "src", "components", "layout", "sidebar.tsx"): {
			"// grit:nav-icons",
			"// grit:nav",
		},
		filepath.Join(root, "cmd", "studio", "main.go"): {
			"// grit:studio-models",
		},
	}

	for file, expectedMarkers := range markers {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("reading %s: %v", file, err)
			continue
		}
		for _, marker := range expectedMarkers {
			if !strings.Contains(string(content), marker) {
				t.Errorf("%s missing marker: %s", filepath.Base(file), marker)
			}
		}
	}
}

func TestDesktopStudioFiles(t *testing.T) {
	root := t.TempDir()
	opts := DesktopOptions{ProjectName: "test-desktop"}

	if err := createDesktopDirectories(root); err != nil {
		t.Fatalf("createDesktopDirectories: %v", err)
	}
	if err := writeDesktopStudioFiles(root, opts); err != nil {
		t.Fatalf("writeDesktopStudioFiles: %v", err)
	}

	studioMain := filepath.Join(root, "cmd", "studio", "main.go")
	if _, err := os.Stat(studioMain); err != nil {
		t.Fatal("cmd/studio/main.go was not created")
	}

	content, err := os.ReadFile(studioMain)
	if err != nil {
		t.Fatalf("reading studio main.go: %v", err)
	}
	if strings.Contains(string(content), "<MODULE>") {
		t.Error("studio main.go still contains <MODULE> placeholder")
	}
	if !strings.Contains(string(content), "test-desktop/internal/config") {
		t.Error("studio main.go does not contain substituted module path")
	}
}
