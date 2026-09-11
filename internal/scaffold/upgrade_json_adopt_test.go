package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Next 16 rewrote tsconfig.json on the first build ("jsx": "react-jsx", the
// .next/dev/types include, and its own formatting), and from then on every
// upgrade reported the file as edited by you.
func TestTSConfigTemplatesAreWhatNextWants(t *testing.T) {
	for name, src := range map[string]string{"admin": adminTSConfig(), "web": webTSConfig()} {
		for _, want := range []string{`"jsx": "react-jsx"`, `".next/dev/types/**/*.ts"`} {
			if !strings.Contains(src, want) {
				t.Errorf("the %s tsconfig lacks %s, so next build rewrites it", name, want)
			}
		}
	}
}

func TestAdoptNextTSConfig(t *testing.T) {
	root := t.TempDir()
	app := filepath.Join(root, "apps", "admin")
	path := filepath.Join(app, "tsconfig.json")
	template := adminTSConfig()

	// Recorded as the old template, then rewritten by Next: same meaning as
	// the new template, different formatting.
	writeTestFile(t, path, strings.Replace(template, `"jsx": "react-jsx"`, `"jsx": "preserve"`, 1))
	release, err := manifest.Start(root, "3.221.0", "scaffold")
	if err != nil {
		t.Fatal(err)
	}
	manifest.Note(path, strings.Replace(template, `"jsx": "react-jsx"`, `"jsx": "preserve"`, 1))
	if err := release(); err != nil {
		t.Fatal(err)
	}
	nextRewrite := "{\n  \"compilerOptions\": {\n" + strings.Join(strings.Fields(template), " ")[len("{ \"compilerOptions\": {"):]
	writeTestFile(t, path, nextRewrite)

	templates := map[string]string{path: template}
	release, _ = manifest.Start(root, "3.222.0", "scaffold")
	if n := adoptNextTSConfig(app, templates); n != 1 {
		t.Fatalf("a Next rewrite meaning exactly the template was not adopted (%d)", n)
	}
	_ = release()

	// A real difference stays the user's.
	writeTestFile(t, path, strings.Replace(template, `"strict": true`, `"strict": false`, 1))
	release, _ = manifest.Start(root, "3.222.0", "scaffold")
	defer func() { _ = release() }()
	if n := adoptNextTSConfig(app, templates); n != 0 {
		t.Error("a tsconfig with a real edit was adopted")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}
