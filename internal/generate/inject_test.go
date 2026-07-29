package generate

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
