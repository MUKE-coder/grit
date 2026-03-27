package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ── ValidateProjectName ───────────────────────────────────────────────────────

func TestValidateProjectName(t *testing.T) {
	valid := []string{
		"my-app",
		"myapp",
		"my-project-123",
		"a",
		"abc123",
		"hello-world",
		"x1y2z3",
	}
	for _, name := range valid {
		if err := ValidateProjectName(name); err != nil {
			t.Errorf("ValidateProjectName(%q) returned unexpected error: %v", name, err)
		}
	}

	invalid := []struct {
		name string
		desc string
	}{
		{"My-App", "uppercase letters"},
		{"my app", "space"},
		{"-my-app", "starts with hyphen"},
		{"my-app-", "ends with hyphen"},
		{"my_app", "underscore"},
		{"", "empty string"},
		{"123app", "starts with digit"},
		{"My App!", "special characters"},
	}
	for _, tc := range invalid {
		if err := ValidateProjectName(tc.name); err == nil {
			t.Errorf("ValidateProjectName(%q) expected error (%s), got nil", tc.name, tc.desc)
		}
	}
}

func TestResolveScaffoldRoot(t *testing.T) {
	t.Run("explicit in-place", func(t *testing.T) {
		root, inPlace, err := resolveScaffoldRoot(Options{ProjectName: "my-app", InPlace: true})
		if err != nil {
			t.Fatalf("resolveScaffoldRoot returned error: %v", err)
		}
		if root != "." {
			t.Fatalf("root = %q, want .", root)
		}
		if !inPlace {
			t.Fatal("inPlace = false, want true")
		}
	})

	t.Run("auto in-place when cwd matches project name", func(t *testing.T) {
		prev, err := os.Getwd()
		if err != nil {
			t.Fatalf("getwd: %v", err)
		}

		target := filepath.Join(t.TempDir(), "my-app")
		if err := os.MkdirAll(target, 0755); err != nil {
			t.Fatalf("mkdirall: %v", err)
		}
		if err := os.Chdir(target); err != nil {
			t.Fatalf("chdir target: %v", err)
		}
		t.Cleanup(func() {
			_ = os.Chdir(prev)
		})

		root, inPlace, err := resolveScaffoldRoot(Options{ProjectName: "my-app"})
		if err != nil {
			t.Fatalf("resolveScaffoldRoot returned error: %v", err)
		}
		if root != "." {
			t.Fatalf("root = %q, want .", root)
		}
		if !inPlace {
			t.Fatal("inPlace = false, want true")
		}
	})

	t.Run("default creates named directory", func(t *testing.T) {
		root, inPlace, err := resolveScaffoldRoot(Options{ProjectName: "my-app"})
		if err != nil {
			t.Fatalf("resolveScaffoldRoot returned error: %v", err)
		}
		if root != "my-app" {
			t.Fatalf("root = %q, want my-app", root)
		}
		if inPlace {
			t.Fatal("inPlace = true, want false")
		}
	})
}

func TestEnsureTargetDirectory(t *testing.T) {
	t.Run("non-in-place rejects existing directory", func(t *testing.T) {
		root := filepath.Join(t.TempDir(), "existing")
		if err := os.MkdirAll(root, 0755); err != nil {
			t.Fatalf("mkdirall: %v", err)
		}

		err := ensureTargetDirectory(root, false, false)
		if err == nil {
			t.Fatal("expected error for existing directory, got nil")
		}
	})

	t.Run("in-place allows empty directory", func(t *testing.T) {
		root := t.TempDir()
		if err := ensureTargetDirectory(root, true, false); err != nil {
			t.Fatalf("ensureTargetDirectory returned error: %v", err)
		}
	})

	t.Run("in-place rejects non-empty directory by default", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "placeholder.txt"), []byte("x"), 0644); err != nil {
			t.Fatalf("writefile: %v", err)
		}

		err := ensureTargetDirectory(root, true, false)
		if err == nil {
			t.Fatal("expected error for non-empty directory, got nil")
		}
		if !strings.Contains(err.Error(), "--force") {
			t.Fatalf("expected --force hint in error, got: %v", err)
		}
	})

	t.Run("in-place force allows non-empty directory", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "placeholder.txt"), []byte("x"), 0644); err != nil {
			t.Fatalf("writefile: %v", err)
		}

		if err := ensureTargetDirectory(root, true, true); err != nil {
			t.Fatalf("ensureTargetDirectory returned error with force: %v", err)
		}
	})
}

// ── ValidateStyle ─────────────────────────────────────────────────────────────

func TestValidateStyle(t *testing.T) {
	// All valid styles should pass
	for _, style := range ValidStyles {
		o := &Options{Style: style}
		if err := o.ValidateStyle(); err != nil {
			t.Errorf("ValidateStyle(%q) returned unexpected error: %v", style, err)
		}
		if o.Style != style {
			t.Errorf("ValidateStyle(%q) changed Style to %q", style, o.Style)
		}
	}

	// Empty defaults to "default"
	o := &Options{}
	if err := o.ValidateStyle(); err != nil {
		t.Errorf("ValidateStyle(\"\") returned unexpected error: %v", err)
	}
	if o.Style != "default" {
		t.Errorf("ValidateStyle(\"\") did not set Style to %q, got %q", "default", o.Style)
	}

	// Invalid styles
	for _, bad := range []string{"neon", "flat", "bootstrap", "material"} {
		o2 := &Options{Style: bad}
		if err := o2.ValidateStyle(); err == nil {
			t.Errorf("ValidateStyle(%q) expected error, got nil", bad)
		}
	}
}

// ── ValidStyles coverage ──────────────────────────────────────────────────────

func TestValidStyles_NotEmpty(t *testing.T) {
	if len(ValidStyles) == 0 {
		t.Fatal("ValidStyles is empty")
	}
	// "default" must always be present so zero-value Options work
	found := false
	for _, s := range ValidStyles {
		if s == "default" {
			found = true
			break
		}
	}
	if !found {
		t.Error("ValidStyles does not include \"default\"")
	}
}

// ── ShouldInclude* ────────────────────────────────────────────────────────────

func TestOptions_ShouldInclude(t *testing.T) {
	tests := []struct {
		name          string
		opts          Options
		includeWeb    bool
		includeAdmin  bool
		includeShared bool
		includeExpo   bool
		includeDocs   bool
	}{
		{
			name:          "default (all false)",
			opts:          Options{},
			includeWeb:    true,
			includeAdmin:  true,
			includeShared: true,
			includeExpo:   false,
			includeDocs:   false,
		},
		{
			name:          "api-only",
			opts:          Options{APIOnly: true},
			includeWeb:    false,
			includeAdmin:  false,
			includeShared: false,
			includeExpo:   false,
			includeDocs:   false,
		},
		{
			name:          "mobile-only",
			opts:          Options{MobileOnly: true},
			includeWeb:    false,
			includeAdmin:  false,
			includeShared: true,
			includeExpo:   true,
			includeDocs:   false,
		},
		{
			name:          "include-expo",
			opts:          Options{IncludeExpo: true},
			includeWeb:    true,
			includeAdmin:  true,
			includeShared: true,
			includeExpo:   true,
			includeDocs:   false,
		},
		{
			name:          "full",
			opts:          Options{Full: true},
			includeWeb:    true,
			includeAdmin:  true,
			includeShared: true,
			includeExpo:   true,
			includeDocs:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.opts.Normalize()
			if got := tt.opts.ShouldIncludeWeb(); got != tt.includeWeb {
				t.Errorf("ShouldIncludeWeb() = %v, want %v", got, tt.includeWeb)
			}
			if got := tt.opts.ShouldIncludeAdmin(); got != tt.includeAdmin {
				t.Errorf("ShouldIncludeAdmin() = %v, want %v", got, tt.includeAdmin)
			}
			if got := tt.opts.ShouldIncludeShared(); got != tt.includeShared {
				t.Errorf("ShouldIncludeShared() = %v, want %v", got, tt.includeShared)
			}
			if got := tt.opts.ShouldIncludeExpo(); got != tt.includeExpo {
				t.Errorf("ShouldIncludeExpo() = %v, want %v", got, tt.includeExpo)
			}
			if got := tt.opts.ShouldIncludeDocs(); got != tt.includeDocs {
				t.Errorf("ShouldIncludeDocs() = %v, want %v", got, tt.includeDocs)
			}
		})
	}
}

// ── createDirectories ─────────────────────────────────────────────────────────

func TestCreateDirectories_APIOnly(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "test-app", APIOnly: true}
	opts.Normalize()

	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories error: %v", err)
	}

	required := []string{
		filepath.Join(root, "apps", "api", "cmd", "server"),
		filepath.Join(root, "apps", "api", "cmd", "migrate"),
		filepath.Join(root, "apps", "api", "cmd", "seed"),
		filepath.Join(root, "apps", "api", "internal", "config"),
		filepath.Join(root, "apps", "api", "internal", "database"),
		filepath.Join(root, "apps", "api", "internal", "models"),
		filepath.Join(root, "apps", "api", "internal", "handlers"),
		filepath.Join(root, "apps", "api", "internal", "middleware"),
		filepath.Join(root, "apps", "api", "internal", "services"),
		filepath.Join(root, "apps", "api", "internal", "routes"),
		filepath.Join(root, "apps", "api", "internal", "cache"),
		filepath.Join(root, "apps", "api", "internal", "mail"),
		filepath.Join(root, "apps", "api", "internal", "jobs"),
		filepath.Join(root, "apps", "api", "internal", "cron"),
		filepath.Join(root, "apps", "api", "internal", "ai"),
	}
	for _, dir := range required {
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("expected directory %s was not created: %v", dir, err)
		}
	}

	// Web/admin/shared/expo/docs dirs must NOT exist in api-only mode
	absent := []string{
		filepath.Join(root, "apps", "web"),
		filepath.Join(root, "apps", "admin"),
		filepath.Join(root, "packages", "shared"),
		filepath.Join(root, "apps", "expo"),
		filepath.Join(root, "apps", "docs"),
	}
	for _, dir := range absent {
		if _, err := os.Stat(dir); err == nil {
			t.Errorf("directory %s should not exist for api-only mode", dir)
		}
	}
}

func TestCreateDirectories_Default(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "test-app"}
	opts.Normalize()

	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories error: %v", err)
	}

	// Web and admin should be created (including test dirs)
	for _, dir := range []string{
		filepath.Join(root, "apps", "web", "app"),
		filepath.Join(root, "apps", "web", "__tests__"),
		filepath.Join(root, "apps", "admin", "app"),
		filepath.Join(root, "apps", "admin", "__tests__"),
		filepath.Join(root, "packages", "shared"),
		filepath.Join(root, "e2e"),
	} {
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("expected directory %s was not created: %v", dir, err)
		}
	}

	// Expo and docs should NOT exist
	for _, dir := range []string{
		filepath.Join(root, "apps", "expo"),
		filepath.Join(root, "apps", "docs"),
	} {
		if _, err := os.Stat(dir); err == nil {
			t.Errorf("directory %s should not exist in default mode", dir)
		}
	}
}

func TestCreateDirectories_TanStackWeb(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "test-app", Architecture: ArchTriple, Frontend: FrontendTanStack}

	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories error: %v", err)
	}

	// TanStack web uses src/ instead of app/
	for _, dir := range []string{
		filepath.Join(root, "apps", "web", "src", "routes"),
		filepath.Join(root, "apps", "web", "src", "components"),
		filepath.Join(root, "apps", "web", "src", "hooks"),
		filepath.Join(root, "apps", "web", "src", "lib"),
	} {
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("expected TanStack directory %s was not created: %v", dir, err)
		}
	}

	// Next.js app/ directory should NOT exist for TanStack
	nextAppDir := filepath.Join(root, "apps", "web", "app")
	if _, err := os.Stat(nextAppDir); err == nil {
		t.Errorf("Next.js app/ directory should not exist when using TanStack Router")
	}
}

func TestCreateSingleDirectories(t *testing.T) {
	root := t.TempDir()

	if err := createSingleDirectories(root); err != nil {
		t.Fatalf("createSingleDirectories error: %v", err)
	}

	// Go internal dirs at root level (no apps/api/)
	for _, dir := range []string{
		filepath.Join(root, "cmd", "server"),
		filepath.Join(root, "internal", "config"),
		filepath.Join(root, "internal", "models"),
		filepath.Join(root, "internal", "handlers"),
		filepath.Join(root, "internal", "routes"),
		filepath.Join(root, "internal", "cache"),
		filepath.Join(root, "frontend", "src", "routes"),
		filepath.Join(root, "frontend", "src", "components"),
	} {
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("expected directory %s was not created: %v", dir, err)
		}
	}

	// apps/ directory should NOT exist
	if _, err := os.Stat(filepath.Join(root, "apps")); err == nil {
		t.Error("apps/ directory should not exist in single app architecture")
	}
}

func TestSingleAppAPIFiles(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "my-single-app", Architecture: ArchSingle, Frontend: FrontendTanStack}

	if err := createSingleDirectories(root); err != nil {
		t.Fatalf("createSingleDirectories: %v", err)
	}
	if err := writeAPIFiles(root, opts); err != nil {
		t.Fatalf("writeAPIFiles: %v", err)
	}

	// Go files should be at root/internal/, not root/apps/api/internal/
	for _, f := range []string{
		filepath.Join(root, "go.mod"),
		filepath.Join(root, "internal", "config", "config.go"),
		filepath.Join(root, "internal", "models", "user.go"),
		filepath.Join(root, "internal", "handlers", "auth.go"),
		filepath.Join(root, "internal", "routes", "routes.go"),
	} {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("expected file %s was not created: %v", f, err)
		}
	}

	// apps/api/ should NOT exist
	if _, err := os.Stat(filepath.Join(root, "apps")); err == nil {
		t.Error("apps/ directory should not exist for single app")
	}

	// go.mod should use root module path (not project-name/apps/api)
	data, _ := os.ReadFile(filepath.Join(root, "go.mod"))
	content := string(data)
	if !strings.Contains(content, "module my-single-app") {
		t.Error("go.mod should contain 'module my-single-app', not 'module my-single-app/apps/api'")
	}
	if strings.Contains(content, "apps/api") {
		t.Error("go.mod should NOT contain 'apps/api' for single app")
	}
}

func TestSingleAppFrontendFiles(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "my-single-app", Architecture: ArchSingle, Frontend: FrontendTanStack}

	if err := createSingleDirectories(root); err != nil {
		t.Fatalf("createSingleDirectories: %v", err)
	}
	if err := writeSingleFrontendFiles(root, opts); err != nil {
		t.Fatalf("writeSingleFrontendFiles: %v", err)
	}

	for _, f := range []string{
		filepath.Join(root, "frontend", "package.json"),
		filepath.Join(root, "frontend", "vite.config.ts"),
		filepath.Join(root, "frontend", "index.html"),
		filepath.Join(root, "frontend", "src", "main.tsx"),
		filepath.Join(root, "frontend", "src", "routes", "__root.tsx"),
		filepath.Join(root, "frontend", "src", "routes", "index.tsx"),
	} {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("expected file %s was not created: %v", f, err)
		}
	}
}

func TestSingleMainGoEmbed(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "my-single-app", Architecture: ArchSingle, Frontend: FrontendTanStack}

	if err := createSingleDirectories(root); err != nil {
		t.Fatalf("createSingleDirectories: %v", err)
	}
	if err := writeSingleMainGo(root, opts); err != nil {
		t.Fatalf("writeSingleMainGo: %v", err)
	}

	mainPath := filepath.Join(root, "cmd", "server", "main.go")
	data, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatalf("reading main.go: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "go:embed frontend/dist") {
		t.Error("single app main.go should contain go:embed directive")
	}
	if !strings.Contains(content, "my-single-app/internal") {
		t.Error("single app main.go should use root module path imports")
	}
	if strings.Contains(content, "apps/api") {
		t.Error("single app main.go should NOT reference apps/api")
	}
}

func TestOptions_APIRoot(t *testing.T) {
	single := Options{ProjectName: "app", Architecture: ArchSingle}
	if got := single.APIRoot("/tmp/app"); got != "/tmp/app" {
		t.Errorf("Single APIRoot = %s, want /tmp/app", got)
	}

	triple := Options{ProjectName: "app", Architecture: ArchTriple}
	want := filepath.Join("/tmp/app", "apps", "api")
	if got := triple.APIRoot("/tmp/app"); got != want {
		t.Errorf("Triple APIRoot = %s, want %s", got, want)
	}
}

func TestOptions_Module(t *testing.T) {
	single := Options{ProjectName: "my-app", Architecture: ArchSingle}
	if got := single.Module(); got != "my-app" {
		t.Errorf("Single Module = %s, want my-app", got)
	}

	triple := Options{ProjectName: "my-app", Architecture: ArchTriple}
	if got := triple.Module(); got != "my-app/apps/api" {
		t.Errorf("Triple Module = %s, want my-app/apps/api", got)
	}
}

func TestWriteWebTanStackFiles(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "test-app", Architecture: ArchTriple, Frontend: FrontendTanStack}

	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}
	if err := writeWebTanStackFiles(root, opts); err != nil {
		t.Fatalf("writeWebTanStackFiles: %v", err)
	}

	files := []string{
		filepath.Join(root, "apps", "web", "package.json"),
		filepath.Join(root, "apps", "web", "vite.config.ts"),
		filepath.Join(root, "apps", "web", "index.html"),
		filepath.Join(root, "apps", "web", "tailwind.config.ts"),
		filepath.Join(root, "apps", "web", "tsconfig.json"),
		filepath.Join(root, "apps", "web", "src", "main.tsx"),
		filepath.Join(root, "apps", "web", "src", "globals.css"),
		filepath.Join(root, "apps", "web", "src", "routes", "__root.tsx"),
		filepath.Join(root, "apps", "web", "src", "routes", "index.tsx"),
		filepath.Join(root, "apps", "web", "src", "routes", "blog", "index.tsx"),
		filepath.Join(root, "apps", "web", "src", "routes", "blog", "$slug.tsx"),
		filepath.Join(root, "apps", "web", "src", "components", "navbar.tsx"),
		filepath.Join(root, "apps", "web", "src", "components", "footer.tsx"),
		filepath.Join(root, "apps", "web", "src", "lib", "api.ts"),
		filepath.Join(root, "apps", "web", "src", "hooks", "use-blogs.ts"),
	}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("expected file %s was not created: %v", f, err)
		}
	}

	// Verify package.json contains vite, not next
	data, _ := os.ReadFile(filepath.Join(root, "apps", "web", "package.json"))
	content := string(data)
	if !strings.Contains(content, "vite") {
		t.Error("TanStack package.json should contain vite dependency")
	}
	if strings.Contains(content, "next") {
		t.Error("TanStack package.json should NOT contain next dependency")
	}
	if !strings.Contains(content, "@tanstack/react-router") {
		t.Error("TanStack package.json should contain @tanstack/react-router")
	}
}

func TestWriteAdminTanStackFiles(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "test-app", Architecture: ArchTriple, Frontend: FrontendTanStack}

	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}
	if err := writeAdminTanStackFiles(root, opts); err != nil {
		t.Fatalf("writeAdminTanStackFiles: %v", err)
	}

	files := []string{
		filepath.Join(root, "apps", "admin", "package.json"),
		filepath.Join(root, "apps", "admin", "vite.config.ts"),
		filepath.Join(root, "apps", "admin", "index.html"),
		filepath.Join(root, "apps", "admin", "src", "main.tsx"),
		filepath.Join(root, "apps", "admin", "src", "routes", "__root.tsx"),
		filepath.Join(root, "apps", "admin", "src", "routes", "_auth", "login.tsx"),
		filepath.Join(root, "apps", "admin", "src", "routes", "_auth", "sign-up.tsx"),
		filepath.Join(root, "apps", "admin", "src", "routes", "_dashboard.tsx"),
		filepath.Join(root, "apps", "admin", "src", "routes", "_dashboard", "dashboard.tsx"),
		filepath.Join(root, "apps", "admin", "src", "routes", "_dashboard", "resources", "users.tsx"),
		filepath.Join(root, "apps", "admin", "src", "components", "layout", "sidebar.tsx"),
		filepath.Join(root, "apps", "admin", "src", "components", "tables", "data-table.tsx"),
		filepath.Join(root, "apps", "admin", "src", "components", "forms", "form-builder.tsx"),
		filepath.Join(root, "apps", "admin", "src", "components", "widgets", "stats-card.tsx"),
		filepath.Join(root, "apps", "admin", "src", "lib", "api-client.ts"),
		filepath.Join(root, "apps", "admin", "src", "hooks", "use-auth.ts"),
		filepath.Join(root, "apps", "admin", "src", "resources", "users.ts"),
	}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("expected file %s was not created: %v", f, err)
		}
	}

	// Verify package.json uses Vite, not Next.js
	data, _ := os.ReadFile(filepath.Join(root, "apps", "admin", "package.json"))
	content := string(data)
	if !strings.Contains(content, "vite") {
		t.Error("Admin TanStack package.json should contain vite")
	}
	if strings.Contains(content, "\"next\"") {
		t.Error("Admin TanStack package.json should NOT contain next")
	}

	// Verify no "use client" directives in TanStack components
	dtData, _ := os.ReadFile(filepath.Join(root, "apps", "admin", "src", "components", "tables", "data-table.tsx"))
	if strings.Contains(string(dtData), "use client") {
		t.Error("TanStack admin data-table.tsx should not contain 'use client' directive")
	}
}

func TestCreateDirectories_Double(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "test-app", Architecture: ArchDouble, Frontend: FrontendNext}

	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories error: %v", err)
	}

	// Web and API should exist
	for _, dir := range []string{
		filepath.Join(root, "apps", "api", "cmd", "server"),
		filepath.Join(root, "apps", "web", "app"),
		filepath.Join(root, "packages", "shared"),
	} {
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("expected directory %s was not created: %v", dir, err)
		}
	}

	// Admin must NOT exist
	adminDir := filepath.Join(root, "apps", "admin")
	if _, err := os.Stat(adminDir); err == nil {
		t.Errorf("admin directory should not exist in double architecture")
	}
}

func TestOptions_DoubleArchitecture(t *testing.T) {
	opts := Options{Architecture: ArchDouble, Frontend: FrontendNext}
	if !opts.ShouldIncludeWeb() {
		t.Error("Double should include web")
	}
	if opts.ShouldIncludeAdmin() {
		t.Error("Double should NOT include admin")
	}
	if !opts.ShouldUseTurborepo() {
		t.Error("Double should use Turborepo")
	}
	if !opts.ShouldIncludeShared() {
		t.Error("Double should include shared")
	}
	if opts.ShouldIncludeSingleSPA() {
		t.Error("Double should NOT be a single SPA")
	}
}

func TestOptions_SingleArchitecture(t *testing.T) {
	opts := Options{Architecture: ArchSingle, Frontend: FrontendTanStack}
	if opts.ShouldIncludeWeb() {
		t.Error("Single should NOT include web (Turborepo app)")
	}
	if opts.ShouldIncludeAdmin() {
		t.Error("Single should NOT include admin")
	}
	if !opts.ShouldIncludeSingleSPA() {
		t.Error("Single should include SPA")
	}
	if opts.ShouldUseTurborepo() {
		t.Error("Single should NOT use Turborepo")
	}
	if !opts.ShouldIncludeFrontend() {
		t.Error("Single should include frontend")
	}
	if !opts.UseTanStack() {
		t.Error("Single with TanStack should return UseTanStack()=true")
	}
}

func TestCreateDirectories_Full(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "test-app", Full: true}
	opts.Normalize()

	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories error: %v", err)
	}

	required := []string{
		filepath.Join(root, "apps", "api", "cmd", "server"),
		filepath.Join(root, "apps", "web", "app"),
		filepath.Join(root, "apps", "admin", "app"),
		filepath.Join(root, "packages", "shared"),
		filepath.Join(root, "apps", "expo", "app"),
		filepath.Join(root, "apps", "docs", "app"),
	}
	for _, dir := range required {
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("expected directory %s was not created: %v", dir, err)
		}
	}
}

// ── writeAPIFiles template substitution ──────────────────────────────────────

func TestWriteAPIFiles_ModuleSubstitution(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "my-project"}
	opts.Normalize()

	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}
	if err := writeAPIFiles(root, opts); err != nil {
		t.Fatalf("writeAPIFiles: %v", err)
	}

	// go.mod must contain the correct module declaration
	goModPath := filepath.Join(root, "apps", "api", "go.mod")
	goModData, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}
	goMod := string(goModData)
	if !strings.Contains(goMod, "module my-project/apps/api") {
		t.Errorf("go.mod missing expected module declaration; got:\n%s", firstN(goMod, 300))
	}
	if strings.Contains(goMod, "{{MODULE}}") {
		t.Error("go.mod still contains unreplaced {{MODULE}} placeholder")
	}

	// main.go must not have any unreplaced placeholders
	mainGoPath := filepath.Join(root, "apps", "api", "cmd", "server", "main.go")
	mainData, err := os.ReadFile(mainGoPath)
	if err != nil {
		t.Fatalf("reading main.go: %v", err)
	}
	if strings.Contains(string(mainData), "{{MODULE}}") {
		t.Error("main.go still contains unreplaced {{MODULE}} placeholder")
	}

	// routes.go must not have any unreplaced placeholders
	routesPath := filepath.Join(root, "apps", "api", "internal", "routes", "routes.go")
	routesData, err := os.ReadFile(routesPath)
	if err != nil {
		t.Fatalf("reading routes.go: %v", err)
	}
	if strings.Contains(string(routesData), "{{MODULE}}") {
		t.Error("routes.go still contains unreplaced {{MODULE}} placeholder")
	}

	// .air.toml must be present
	airPath := filepath.Join(root, "apps", "api", ".air.toml")
	if _, err := os.Stat(airPath); err != nil {
		t.Errorf(".air.toml was not created: %v", err)
	}
}

func TestWriteAPIFiles_KeyFilesExist(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "acme"}
	opts.Normalize()

	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}
	if err := writeAPIFiles(root, opts); err != nil {
		t.Fatalf("writeAPIFiles: %v", err)
	}

	files := []string{
		filepath.Join(root, "apps", "api", "go.mod"),
		filepath.Join(root, "apps", "api", ".gitignore"),
		filepath.Join(root, "apps", "api", ".air.toml"),
		filepath.Join(root, "apps", "api", "cmd", "server", "main.go"),
		filepath.Join(root, "apps", "api", "internal", "config", "config.go"),
		filepath.Join(root, "apps", "api", "internal", "database", "database.go"),
		filepath.Join(root, "apps", "api", "internal", "models", "user.go"),
		filepath.Join(root, "apps", "api", "internal", "models", "upload.go"),
		filepath.Join(root, "apps", "api", "internal", "services", "auth.go"),
		filepath.Join(root, "apps", "api", "internal", "handlers", "auth.go"),
		filepath.Join(root, "apps", "api", "internal", "handlers", "user.go"),
		filepath.Join(root, "apps", "api", "internal", "middleware", "auth.go"),
		filepath.Join(root, "apps", "api", "internal", "middleware", "cors.go"),
		filepath.Join(root, "apps", "api", "internal", "middleware", "logger.go"),
		filepath.Join(root, "apps", "api", "internal", "routes", "routes.go"),
		// Test + bench files generated for the API
		filepath.Join(root, "apps", "api", "internal", "handlers", "auth_test.go"),
		filepath.Join(root, "apps", "api", "internal", "handlers", "user_test.go"),
		filepath.Join(root, "apps", "api", "internal", "handlers", "bench_test.go"),
	}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("expected file %s was not created: %v", f, err)
		}
	}
}

// ── writeFrontendTestFiles ────────────────────────────────────────────────────

func TestWriteFrontendTestFiles_KeyFilesExist(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "acme"} // default mode: web + admin
	opts.Normalize()

	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}
	if err := writeFrontendTestFiles(root, opts); err != nil {
		t.Fatalf("writeFrontendTestFiles: %v", err)
	}

	files := []string{
		// Web vitest
		filepath.Join(root, "apps", "web", "vitest.config.ts"),
		filepath.Join(root, "apps", "web", "vitest.setup.ts"),
		filepath.Join(root, "apps", "web", "__tests__", "navbar.test.tsx"),
		filepath.Join(root, "apps", "web", "__tests__", "footer.test.tsx"),
		// Admin vitest
		filepath.Join(root, "apps", "admin", "vitest.config.ts"),
		filepath.Join(root, "apps", "admin", "vitest.setup.ts"),
		filepath.Join(root, "apps", "admin", "__tests__", "login.test.tsx"),
		filepath.Join(root, "apps", "admin", "__tests__", "utils.test.ts"),
		// Playwright E2E
		filepath.Join(root, "playwright.config.ts"),
		filepath.Join(root, "e2e", "auth.spec.ts"),
		filepath.Join(root, "e2e", "admin.spec.ts"),
	}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("expected file %s was not created: %v", f, err)
		}
	}
}

func TestWriteFrontendTestFiles_APIOnly_NoFiles(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "acme", APIOnly: true}
	opts.Normalize()

	if err := createDirectories(root, opts); err != nil {
		t.Fatalf("createDirectories: %v", err)
	}
	if err := writeFrontendTestFiles(root, opts); err != nil {
		t.Fatalf("writeFrontendTestFiles: %v", err)
	}

	// API-only mode: no frontend test files should be created
	absent := []string{
		filepath.Join(root, "playwright.config.ts"),
		filepath.Join(root, "apps", "web", "vitest.config.ts"),
		filepath.Join(root, "apps", "admin", "vitest.config.ts"),
	}
	for _, f := range absent {
		if _, err := os.Stat(f); err == nil {
			t.Errorf("file %s should not exist in API-only mode", f)
		}
	}
}

// ── writeFile helper ──────────────────────────────────────────────────────────

func TestWriteFile_CreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "file.txt")
	if err := writeFile(path, "hello world"); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("file content = %q, want %q", string(data), "hello world")
	}
}

func TestWriteFile_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")
	if err := writeFile(path, "first"); err != nil {
		t.Fatalf("writeFile (first): %v", err)
	}
	if err := writeFile(path, "second"); err != nil {
		t.Fatalf("writeFile (second): %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}
	if string(data) != "second" {
		t.Errorf("file content = %q, want %q", string(data), "second")
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

// firstN returns the first n bytes of s as a string.
func firstN(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
