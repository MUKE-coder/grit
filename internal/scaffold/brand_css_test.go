package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The whole value of brand.css is in what Grit does NOT do to it, so that is
// what these test: created when absent, left exactly alone when present, and
// never recorded in the manifest.

func mkFrontend(t *testing.T, root, app, cssDir string) string {
	t.Helper()
	dir := filepath.Join(root, "apps", app, cssDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("creating %s: %v", dir, err)
	}
	// brandCSSPath probes for globals.css, so the frontend needs one to be
	// recognised as a frontend at all.
	if err := os.WriteFile(filepath.Join(dir, "globals.css"), []byte(":root{}"), 0o644); err != nil {
		t.Fatalf("writing globals.css: %v", err)
	}
	return filepath.Join(dir, "brand.css")
}

func TestBrandCSSIsCreatedForEveryFrontend(t *testing.T) {
	root := t.TempDir()
	web := mkFrontend(t, root, "web", "app")     // App Router
	admin := mkFrontend(t, root, "admin", "src") // Vite

	made := writeBrandCSS(root, Options{ProjectName: "kaya"})
	if made != 2 {
		t.Fatalf("created %d brand.css files, want 2", made)
	}

	for _, path := range []string{web, admin} {
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading %s: %v", path, err)
		}
		text := string(body)
		if !strings.Contains(text, "kaya") {
			t.Errorf("%s does not name the project", path)
		}
		// Every rule ships commented out, so adding the file to an existing
		// project cannot change a single pixel of it. Checked by removing the
		// comments and seeing whether any rule is left, because the template
		// shows the selectors inside the comments and a plain substring search
		// cannot tell an example from a declaration.
		if live := stripCSSComments(text); strings.Contains(live, "{") {
			t.Errorf("%s declares a live rule; the seed must set nothing.\nleft after stripping comments: %q",
				path, strings.TrimSpace(live))
		}
		if !strings.Contains(text, "html[data-theme]") {
			t.Errorf("%s does not show the selector to use", path)
		}
	}
}

func TestBrandCSSIsNeverOverwritten(t *testing.T) {
	root := t.TempDir()
	path := mkFrontend(t, root, "web", "app")

	mine := "html[data-theme] { --accent: #25d366; }\n"
	if err := os.WriteFile(path, []byte(mine), 0o644); err != nil {
		t.Fatalf("writing mine: %v", err)
	}

	// Twice, because an upgrade runs on a project that has already had one.
	writeBrandCSS(root, Options{ProjectName: "kaya"})
	made := writeBrandCSS(root, Options{ProjectName: "kaya"})
	if made != 0 {
		t.Fatalf("reported %d created, want 0: an existing file was touched", made)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading back: %v", err)
	}
	if string(body) != mine {
		t.Fatalf("brand.css was rewritten.\n got: %q\nwant: %q", body, mine)
	}
}

// The reason brand.css bypasses writeFile. A manifest entry is what turns an
// edited file into "left alone, because you have edited these files" on every
// later upgrade, which is the warning this file exists to stop producing.
func TestBrandCSSIsNotInTheManifest(t *testing.T) {
	root := t.TempDir()
	mkFrontend(t, root, "web", "app")

	writeBrandCSS(root, Options{ProjectName: "kaya"})

	manifestPath := filepath.Join(root, ".grit", "manifest.json")
	body, err := os.ReadFile(manifestPath)
	if os.IsNotExist(err) {
		return // nothing recorded at all, which is the point
	}
	if err != nil {
		t.Fatalf("reading manifest: %v", err)
	}
	if strings.Contains(string(body), "brand.css") {
		t.Fatal("brand.css is recorded in the manifest; it must not be")
	}
}

func TestBrandCSSSkipsADirectoryWithNoStylesheet(t *testing.T) {
	root := t.TempDir()
	// An app directory with no globals.css is not a frontend with a palette.
	if err := os.MkdirAll(filepath.Join(root, "apps", "web", "app"), 0o755); err != nil {
		t.Fatal(err)
	}
	if made := writeBrandCSS(root, Options{ProjectName: "kaya"}); made != 0 {
		t.Fatalf("created %d, want 0", made)
	}
}

// brand.css only wins if it is imported after globals.css. A template that
// imports it first would compile, render, and silently do nothing.
func TestEveryLayoutImportsBrandCSSAfterGlobals(t *testing.T) {
	sources := map[string]string{
		"web_files.go":            readScaffoldSource(t, "web_files.go"),
		"admin_files.go":          readScaffoldSource(t, "admin_files.go"),
		"desktop_client_files.go": readScaffoldSource(t, "desktop_client_files.go"),
	}

	for name, src := range sources {
		globals := strings.Index(src, `import "./globals.css";`)
		brand := strings.Index(src, `import "./brand.css";`)
		if globals < 0 {
			t.Errorf("%s: no globals.css import found", name)
			continue
		}
		if brand < 0 {
			t.Errorf("%s: does not import brand.css", name)
			continue
		}
		if brand < globals {
			t.Errorf("%s: imports brand.css before globals.css, so it is overridden rather than overriding", name)
		}
	}
}

// stripCSSComments removes every /* ... */ block, so what is left is the CSS
// that would actually apply.
func stripCSSComments(css string) string {
	var out strings.Builder
	for {
		start := strings.Index(css, "/*")
		if start < 0 {
			out.WriteString(css)
			return out.String()
		}
		out.WriteString(css[:start])
		end := strings.Index(css[start:], "*/")
		if end < 0 {
			return out.String() // unterminated: the rest is comment
		}
		css = css[start+end+2:]
	}
}

// Writing brand.css is half the job: it only wins because it is imported
// after globals.css, and that import lives in the layout, which is
// framework-owned. On an upgrade a project that has edited its layout keeps
// its own copy, so v3.421.0 delivered the file and not the line, and the
// WhatsApp blueprint was painted green and stayed blue with nothing saying
// why.
func TestTheLayoutGetsTheImportEvenWhenGritDoesNotRewriteIt(t *testing.T) {
	root := t.TempDir()
	app := filepath.Join(root, "apps", "web", "app")
	if err := os.MkdirAll(app, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{
		"globals.css": ":root{}",
		"brand.css":   "/* yours */",
		// A layout the developer has edited: Grit will not overwrite it.
		"layout.tsx": "import type { Metadata } from \"next\";\n" +
			"import \"./globals.css\";\n" +
			"export default function RootLayout() { return null }\n",
	} {
		if err := os.WriteFile(filepath.Join(app, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if fixed := ensureBrandCSSImport(root); fixed != 1 {
		t.Fatalf("repaired %d layouts, want 1", fixed)
	}

	body, err := os.ReadFile(filepath.Join(app, "layout.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)

	globals := strings.Index(text, `import "./globals.css";`)
	brand := strings.Index(text, `import "./brand.css";`)
	if brand < 0 {
		t.Fatalf("brand.css is still not imported:\n%s", text)
	}
	if brand < globals {
		t.Fatalf("brand.css is imported BEFORE globals.css, so it loses:\n%s", text)
	}

	// Twice is the normal case: every upgrade runs this.
	if fixed := ensureBrandCSSImport(root); fixed != 0 {
		t.Fatalf("a second run repaired %d layouts, want 0", fixed)
	}
	again, err := os.ReadFile(filepath.Join(app, "layout.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	if string(again) != text {
		t.Fatalf("a second run changed the file:\n%s", again)
	}
}

// An app with no brand.css has nothing to import, and a layout with no
// globals.css import has no anchor to put it after. Neither is an error.
func TestTheImportRepairLeavesAnUnrelatedLayoutAlone(t *testing.T) {
	root := t.TempDir()
	app := filepath.Join(root, "apps", "web", "app")
	if err := os.MkdirAll(app, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(app, "globals.css"), []byte(":root{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	original := "export default function RootLayout() { return null }\n"
	if err := os.WriteFile(filepath.Join(app, "layout.tsx"), []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	if fixed := ensureBrandCSSImport(root); fixed != 0 {
		t.Fatalf("repaired %d layouts, want 0", fixed)
	}
	body, err := os.ReadFile(filepath.Join(app, "layout.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != original {
		t.Fatalf("the layout was changed:\n%s", body)
	}
}
