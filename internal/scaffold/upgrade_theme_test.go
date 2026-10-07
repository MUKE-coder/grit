package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

// An upgrade keeps the project's theme.
//
// Upgrade reads the style, the frontend and the architecture off the project,
// and did not read the theme. Files that bake colour values in, the mobile
// app's Tailwind config among them, were therefore regenerated as atlas: an
// emerald project came out blue on its next upgrade and nothing said so.
func TestUpgradeReadsTheProjectTheme(t *testing.T) {
	root := t.TempDir()
	env := "APP_PORT=8080\n# a comment\nTHEME=emerald\nVITE_THEME=emerald\n"
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte(env), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := readProjectTheme(root); got != "emerald" {
		t.Errorf("got %q, want emerald", got)
	}

	// A theme nobody has heard of is not a theme.
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte("THEME=chartreuse\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := readProjectTheme(root); got != "" {
		t.Errorf("an unknown theme came back as %q", got)
	}

	// No .env at all is not an error: Normalize fills in the default.
	if got := readProjectTheme(t.TempDir()); got != "" {
		t.Errorf("a project with no .env came back as %q", got)
	}
}
