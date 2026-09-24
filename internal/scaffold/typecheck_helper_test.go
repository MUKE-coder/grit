package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

// Emit generated TSX into a real project so its own type checker can read it.
//
// The tests in this package match strings: they can tell you a component
// exists and cannot tell you it compiles. Both of these skip unless you point
// them at a project, and they exist because a stray backtick in a template
// once produced five files that every Go test accepted and tsc rejected:
//
//	SHELL_OUT=<project>/apps/admin/components/auth go test ./internal/scaffold/ -run WriteShells
//	SHARED_OUT=<project>/packages/shared          go test ./internal/scaffold/ -run WriteShared
//	cd <project>/apps/admin && npx tsc --noEmit

// Writes the generated shells into a real project so its tsc can check them.
func TestWriteShellsForTypeCheck(t *testing.T) {
	dir := os.Getenv("SHELL_OUT")
	if dir == "" {
		t.Skip("SHELL_OUT not set")
	}
	for file, source := range authShellFiles() {
		if err := os.WriteFile(filepath.Join(dir, file), []byte(source), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Logf("wrote %s", file)
	}
	if err := os.WriteFile(filepath.Join(dir, "AuthShell.tsx"), []byte(adminAuthShellDispatcher()), 0o644); err != nil {
		t.Fatal(err)
	}
}

// Writes the shared brand + themes into a real project for the same reason.
func TestWriteSharedForTypeCheck(t *testing.T) {
	dir := os.Getenv("SHARED_OUT")
	if dir == "" {
		t.Skip("SHARED_OUT not set")
	}
	files := map[string]string{
		"brand.config.ts": sharedBrandConfig(Options{ProjectName: "s321-fresh"}),
		"themes.ts":       sharedThemes(),
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Logf("wrote %s", name)
	}
}
