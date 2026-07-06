package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAlreadyInjected(t *testing.T) {
	content := "func f() {\n\t\t&Post{},\n\t\t// grit:models\n}\n"
	cases := []struct {
		name string
		code string
		want bool
	}{
		{"exact", "\t\t&Post{},", true},
		{"reindented", "  &Post{},", true},       // whitespace-insensitive
		{"collapsed", "&Post{},", true},          // no leading indent
		{"absent", "\t\t&Comment{},", false},     // different resource
		{"prefix not a match", "&Pos{},", false}, // must be the whole snippet
		{"empty", "", false},                     // empty never counts as present
	}
	for _, tt := range cases {
		if got := alreadyInjected(content, tt.code); got != tt.want {
			t.Errorf("%s: alreadyInjected(%q) = %v, want %v", tt.name, tt.code, got, tt.want)
		}
	}
}

func TestInjectBefore_Idempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "routes.go")
	base := "package routes\n\nfunc R() {\n\t// grit:models\n}\n"
	if err := os.WriteFile(path, []byte(base), 0o644); err != nil {
		t.Fatal(err)
	}
	code := "\t&Post{},"

	for i := 0; i < 3; i++ {
		if err := injectBefore(path, "// grit:models", code); err != nil {
			t.Fatalf("injectBefore #%d: %v", i, err)
		}
	}

	got, _ := os.ReadFile(path)
	if n := strings.Count(string(got), "&Post{}"); n != 1 {
		t.Fatalf("expected &Post{} injected exactly once, got %d\n%s", n, got)
	}
}

func TestInjectInline_Idempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "studio.go")
	base := "studio.Register(/* grit:studio */)\n"
	if err := os.WriteFile(path, []byte(base), 0o644); err != nil {
		t.Fatal(err)
	}
	code := "&models.Post{}, "

	for i := 0; i < 3; i++ {
		if err := injectInline(path, "/* grit:studio */", code); err != nil {
			t.Fatalf("injectInline #%d: %v", i, err)
		}
	}

	got, _ := os.ReadFile(path)
	if n := strings.Count(string(got), "&models.Post{}"); n != 1 {
		t.Fatalf("expected model registered exactly once, got %d\n%s", n, got)
	}
}
