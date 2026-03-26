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
