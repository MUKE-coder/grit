package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeCSS(t *testing.T, appRoot, body string) {
	t.Helper()
	dir := filepath.Join(appRoot, "app")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "globals.css"), []byte(body), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

// An app whose stylesheet is still on the v3 directives keeps the v3 plugin.
//
// This is the case that broke a real project: upgrade wrote the v4 postcss
// config while the manifest guard held back the edited package.json and
// globals.css, so a v4 plugin parsed a v3 stylesheet and every utility built
// from tailwind.config.ts became "Cannot apply unknown utility".
func TestV3StylesheetKeepsTheV3PostCSSPlugin(t *testing.T) {
	root := t.TempDir()
	writeCSS(t, root, "@tailwind base;\n@tailwind components;\n@tailwind utilities;\n")

	got := postCSSConfigFor(root)
	if !strings.Contains(got, "tailwindcss: {}") {
		t.Error("a v3 stylesheet did not get the v3 plugin")
	}
	// Checked against the code rather than the whole file: the comment names
	// the v4 plugin on purpose, as part of the migration instructions.
	if strings.Contains(withoutComments(got), "@tailwindcss/postcss") {
		t.Error("a v3 stylesheet got the v4 plugin; the app will not build")
	}
	// And it says how to migrate, because the alternative is a file that looks
	// out of date for no stated reason.
	if !strings.Contains(got, "@config") {
		t.Error("the v3 config does not explain the migration")
	}
}

func TestV4StylesheetGetsTheV4PostCSSPlugin(t *testing.T) {
	root := t.TempDir()
	writeCSS(t, root, "@import \"tailwindcss\";\n@theme {\n  --color-accent: #6c5ce7;\n}\n")

	got := postCSSConfigFor(root)
	if !strings.Contains(got, "@tailwindcss/postcss") {
		t.Error("a v4 stylesheet did not get the v4 plugin")
	}
	if strings.Contains(got, "autoprefixer") {
		t.Error("v4 prefixes for itself; autoprefixer can mangle its @property rules")
	}
}

// A fresh scaffold has no stylesheet yet when the config is written, and the
// one about to be written is v4.
func TestNoStylesheetYetGetsTheV4PostCSSPlugin(t *testing.T) {
	got := postCSSConfigFor(t.TempDir())
	if !strings.Contains(got, "@tailwindcss/postcss") {
		t.Error("a fresh scaffold did not get the v4 plugin")
	}
}

// A stylesheet that imports tailwindcss is v4 even if the word "@tailwind"
// appears later in a comment, so the v4 check runs first.
func TestV4DetectionBeatsAStrayMention(t *testing.T) {
	root := t.TempDir()
	writeCSS(t, root, "@import \"tailwindcss\";\n/* migrated from @tailwind base; */\n")

	if stylesheetIsV3(root) {
		t.Error("a v4 stylesheet was read as v3 because of a comment")
	}
}

func TestStylesheetFoundInSrcLayouts(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "src", "styles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "globals.css"),
		[]byte("@tailwind base;\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	if !stylesheetIsV3(root) {
		t.Error("a v3 stylesheet under src/styles was not found; TanStack apps live there")
	}
}

// withoutComments drops // lines so an assertion about the config can ignore
// the prose above it.
func withoutComments(s string) string {
	var kept []string
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}
