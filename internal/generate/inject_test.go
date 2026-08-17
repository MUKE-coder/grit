package generate

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ── ensureHandlerInit ────────────────────────────────────────────────────────

// The regression this function exists for: regenerating a resource after it
// gained a file field used to append a second "xHandler :=" block, because the
// guard compared the whole block and the new one carried Storage. The project
// then failed to compile with "no new variables on left side of :=".
func TestEnsureHandlerInit(t *testing.T) {
	const routes = `package routes

func Setup(db *gorm.DB, svc *Services) {
	productHandler := &handlers.ProductHandler{
		DB: db,
	}
	// grit:handlers
}
`

	t.Run("adds the handler when it is absent", func(t *testing.T) {
		f := writeTempFile(t, "routes.go", `package routes

func Setup(db *gorm.DB, svc *Services) {
	// grit:handlers
}
`)
		what, err := ensureHandlerInit(f, "product", "Product", false)
		if err != nil {
			t.Fatalf("ensureHandlerInit: %v", err)
		}
		if !strings.Contains(what, "Injected") {
			t.Errorf("expected an injection, reported %q", what)
		}
		got := readFile(t, f)
		if !strings.Contains(got, "productHandler := &handlers.ProductHandler{") {
			t.Errorf("handler not written:\n%s", got)
		}
		if strings.Contains(got, "Storage:") {
			t.Error("a resource with no file fields must not ask for Storage")
		}
	})

	t.Run("does not declare the same handler twice", func(t *testing.T) {
		f := writeTempFile(t, "routes.go", routes)
		// The file field is the case that used to slip past the guard.
		if _, err := ensureHandlerInit(f, "product", "Product", true); err != nil {
			t.Fatalf("ensureHandlerInit: %v", err)
		}
		got := readFile(t, f)
		if n := strings.Count(got, "productHandler := &handlers.ProductHandler{"); n != 1 {
			t.Fatalf("handler declared %d times, want 1:\n%s", n, got)
		}
		// And the reason it is only once: Storage was added to the block that
		// was already there, not written as a second one.
		if !strings.Contains(got, "Storage: svc.Storage,") {
			t.Errorf("a resource that gained a file field must get Storage wired:\n%s", got)
		}
		if _, err := format.Source([]byte(got)); err != nil {
			t.Errorf("result does not parse as Go: %v\n%s", err, got)
		}
	})

	t.Run("is idempotent", func(t *testing.T) {
		f := writeTempFile(t, "routes.go", routes)
		for i := 0; i < 3; i++ {
			if _, err := ensureHandlerInit(f, "product", "Product", true); err != nil {
				t.Fatalf("run %d: %v", i, err)
			}
		}
		got := readFile(t, f)
		if n := strings.Count(got, "Storage: svc.Storage,"); n != 1 {
			t.Errorf("Storage wired %d times, want 1:\n%s", n, got)
		}
	})

	t.Run("leaves a hand-edited block alone", func(t *testing.T) {
		f := writeTempFile(t, "routes.go", `package routes

func Setup(db *gorm.DB, svc *Services) {
	productHandler := &handlers.ProductHandler{
		DB:      db,
		Storage: svc.Storage,
		Pricing: pricingService,
	}
	// grit:handlers
}
`)
		what, err := ensureHandlerInit(f, "product", "Product", true)
		if err != nil {
			t.Fatalf("ensureHandlerInit: %v", err)
		}
		if !strings.Contains(what, "leaving it alone") {
			t.Errorf("expected to leave the block alone, reported %q", what)
		}
		if !strings.Contains(readFile(t, f), "Pricing: pricingService,") {
			t.Error("a hand-added field was removed")
		}
	})
}

// ── injectBefore ─────────────────────────────────────────────────────────────

func TestInjectBefore(t *testing.T) {
	t.Run("inserts code on line before marker", func(t *testing.T) {
		f := writeTempFile(t, "file.go", `package models

func Models() []interface{} {
	return []interface{}{
		&User{},
		// grit:models
	}
}
`)
		if err := injectBefore(f, "// grit:models", "\t\t&Post{},"); err != nil {
			t.Fatalf("injectBefore error: %v", err)
		}
		got := readFile(t, f)
		if !strings.Contains(got, "\t\t&Post{},\n\t\t// grit:models") {
			t.Errorf("injected code not found in expected position:\n%s", got)
		}
		// Original marker must still be present
		if !strings.Contains(got, "// grit:models") {
			t.Error("marker was removed after injection")
		}
		// &User{} must still be there
		if !strings.Contains(got, "&User{},") {
			t.Error("&User{} was removed after injection")
		}
	})

	t.Run("marker not found returns error", func(t *testing.T) {
		f := writeTempFile(t, "file.go", "package x\n")
		err := injectBefore(f, "// grit:missing", "code")
		if err == nil {
			t.Error("expected error for missing marker, got nil")
		}
	})

	t.Run("idempotent marker still present after second inject", func(t *testing.T) {
		f := writeTempFile(t, "file.go", `// grit:models
`)
		if err := injectBefore(f, "// grit:models", "&Post{},"); err != nil {
			t.Fatalf("first inject: %v", err)
		}
		if err := injectBefore(f, "// grit:models", "&Tag{},"); err != nil {
			t.Fatalf("second inject: %v", err)
		}
		got := readFile(t, f)
		if !strings.Contains(got, "&Post{},") {
			t.Error("&Post{} missing after second inject")
		}
		if !strings.Contains(got, "&Tag{},") {
			t.Error("&Tag{} missing after second inject")
		}
		if !strings.Contains(got, "// grit:models") {
			t.Error("marker missing after second inject")
		}
	})

	t.Run("file not found returns error", func(t *testing.T) {
		err := injectBefore("/tmp/does-not-exist-grit-inject.go", "// marker", "code")
		if err == nil {
			t.Error("expected error for missing file, got nil")
		}
	})
}

// ── injectInline ─────────────────────────────────────────────────────────────

func TestInjectInline(t *testing.T) {
	t.Run("inserts code immediately before marker", func(t *testing.T) {
		f := writeTempFile(t, "routes.go", `studio.Mount(r, db, []interface{}{&models.User{}, /* grit:studio */}, cfg)
`)
		if err := injectInline(f, "/* grit:studio */", "&models.Post{}, "); err != nil {
			t.Fatalf("injectInline error: %v", err)
		}
		got := readFile(t, f)
		if !strings.Contains(got, "&models.Post{}, /* grit:studio */") {
			t.Errorf("inline code not found in expected position:\n%s", got)
		}
	})

	// Generated Go is gofmt'd on write, and gofmt removes the optional trailing
	// comma from a single-line composite literal. The injector must not depend
	// on that comma being there — without this, the second resource generated
	// into a project produced "&models.User{} &models.Post{}", a bitwise AND
	// that failed to compile with a mismatched-types error far from the cause.
	t.Run("supplies the separator when gofmt dropped the trailing comma", func(t *testing.T) {
		f := writeTempFile(t, "routes.go", `studio.Mount(r, db, []interface{}{&models.User{} /* grit:studio */}, cfg)
`)
		if err := injectInline(f, "/* grit:studio */", "&models.Post{}, "); err != nil {
			t.Fatalf("injectInline error: %v", err)
		}
		got := readFile(t, f)
		if !strings.Contains(got, "&models.User{}, &models.Post{}, /* grit:studio */") {
			t.Errorf("missing separator between elements:\n%s", got)
		}
		// The result has to be valid Go, not merely look right.
		if _, err := format.Source([]byte("package x\nfunc f() {\n" + got + "\n}\n")); err != nil {
			t.Errorf("injected result does not parse: %v\n%s", err, got)
		}
	})

	// The empty-list case must NOT gain a leading comma.
	t.Run("no separator after an opening bracket", func(t *testing.T) {
		f := writeTempFile(t, "app.go", `func NewApp(/* grit:constructor-params */) *App {}
`)
		if err := injectInline(f, "/* grit:constructor-params */", "postSvc *service.PostService, "); err != nil {
			t.Fatalf("injectInline error: %v", err)
		}
		got := readFile(t, f)
		if strings.Contains(got, "(, ") {
			t.Errorf("leading comma inserted into an empty list:\n%s", got)
		}
		if !strings.Contains(got, "NewApp(postSvc *service.PostService, /* grit:constructor-params */)") {
			t.Errorf("unexpected result:\n%s", got)
		}
	})

	t.Run("marker not found returns error", func(t *testing.T) {
		f := writeTempFile(t, "routes.go", "package x\n")
		err := injectInline(f, "/* grit:missing */", "code")
		if err == nil {
			t.Error("expected error for missing marker, got nil")
		}
	})

	t.Run("file not found returns error", func(t *testing.T) {
		err := injectInline("/tmp/does-not-exist-grit-inline.go", "/* marker */", "code")
		if err == nil {
			t.Error("expected error for missing file, got nil")
		}
	})
}

// ── guessLucideIcon ───────────────────────────────────────────────────────────

func TestGuessLucideIcon(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Post", "FileText"},
		{"Article", "Newspaper"},
		{"Comment", "MessageSquare"},
		{"Category", "FolderTree"},
		{"Product", "Package"},
		{"Order", "ShoppingCart"},
		{"User", "Users"},
		{"Task", "CheckSquare"},
		{"Event", "Calendar"},
		{"Invoice", "Receipt"},
		// Unknown name falls back to "Database"
		{"Widget", "Database"},
		{"Foo", "Database"},
	}

	for _, tt := range tests {
		got := guessLucideIcon(tt.name)
		if got != tt.want {
			t.Errorf("guessLucideIcon(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return path
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readFile: %v", err)
	}
	return string(data)
}
