package scaffold

import (
	"os"
	"strings"
	"testing"
)

// Eight themes, eight sign-in layouts, and four places that have to agree
// about them: the theme registry, the shell files, the dispatcher and the CLI
// flag. A theme whose layout has no case falls through to the default, which
// means somebody picks `--theme emerald` and gets the atlas screen with no
// error anywhere.

var authThemes = map[string]string{
	"atlas":   "split-static",
	"aurora":  "centered",
	"pulse":   "split-carousel",
	"coral":   "modal",
	"amber":   "boxed",
	"sky":     "banner",
	"mono":    "showcase",
	"emerald": "quote",
}

func TestEveryThemeNamesALayoutTheDispatcherHandles(t *testing.T) {
	themes := sharedThemes()
	dispatcher := adminAuthShellDispatcher()

	for name, layout := range authThemes {
		if !strings.Contains(themes, `name: "`+name+`"`) {
			t.Errorf("theme %q is not in the registry", name)
			continue
		}
		if !strings.Contains(themes, `authLayout: "`+layout+`"`) {
			t.Errorf("theme %q names layout %q, which no theme declares", name, layout)
		}
		// split-static is the default arm rather than a case of its own.
		if layout != "split-static" && !strings.Contains(dispatcher, `case "`+layout+`":`) {
			t.Errorf("layout %q has no case in AuthShell, so that theme silently renders atlas", layout)
		}
		if !strings.Contains(themes, `"`+name+`"`) {
			t.Errorf("%q is missing from the ThemeName union", name)
		}
	}
}

func TestTheNewShellsAllPublishTheThemeTokens(t *testing.T) {
	// Every form control in the scaffold is styled against these custom
	// properties. A shell that omits one renders an input with no border and
	// no error anywhere to explain it.
	required := []string{
		"--auth-bg", "--auth-fg", "--auth-card", "--auth-border",
		"--auth-muted", "--auth-primary", "--auth-primary-fg",
		"--auth-accent", "--auth-radius",
	}

	for file, source := range authShellFiles() {
		name := shellName(file)
		for _, token := range required {
			if !strings.Contains(source, token) {
				t.Errorf("%s never sets %s", name, token)
			}
		}
		if !strings.Contains(source, "export function "+name+"(props: Props)") {
			t.Errorf("%s does not export a component of its own name", name)
		}
		// The form is the same component in every layout: a shell that renders
		// its own inputs is a shell that will drift from the validation.
		if !strings.Contains(source, "{children}") {
			t.Errorf("%s never renders the form", name)
		}
		if !strings.Contains(source, "<ErrorBanner message={errorMessage} />") {
			t.Errorf("%s swallows the server's error message", name)
		}
		if !strings.Contains(source, "switchLinks[mode]") {
			t.Errorf("%s has no way to reach sign-up or password reset", name)
		}
	}
}

func TestTheNewShellsAreWrittenIntoBothAdmins(t *testing.T) {
	next := readScaffoldSource(t, "admin_files.go")
	tanstack := readScaffoldSource(t, "admin_tanstack_files.go")

	for file := range authShellFiles() {
		if !strings.Contains(next, file) {
			t.Errorf("the Next admin never writes %s", file)
		}
		if !strings.Contains(tanstack, file) {
			t.Errorf("the TanStack admin never writes %s", file)
		}
	}
}

// The quote layout ships a placeholder testimonial. It must read as one.
//
// A scaffold that ships a fabricated quote under a plausible name is a
// scaffold that puts a lie in production the first time somebody forgets to
// edit it.
func TestThePlaceholderTestimonialAdmitsItIsOne(t *testing.T) {
	brandTS := sharedBrandConfig(Options{ProjectName: "demo"})
	if !strings.Contains(brandTS, "quoteAuthor:") {
		t.Fatal("brand.proof has no quote author")
	}
	author := brandTS[strings.Index(brandTS, "quoteAuthor:"):]
	author = author[:strings.Index(author, "\n")]
	if !strings.Contains(strings.ToLower(author), "replace") {
		t.Errorf("the placeholder quote is attributed to %s, which reads as a real endorsement", strings.TrimSpace(author))
	}
}

// readScaffoldSource reads one of this package's own files.
//
// The file maps need a root and an Options to build, and what this is checking
// is that the entry exists at all, which the source answers directly.
func readScaffoldSource(t *testing.T, name string) string {
	t.Helper()
	body, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("reading %s: %v", name, err)
	}
	return string(body)
}

// The desktop app has three shells, not eight, and maps the rest onto them.
//
// That mapping has to be explicit. A layout that reaches the default arm by
// accident renders Atlas for somebody who asked for Emerald, and nothing
// anywhere says so.
func TestTheDesktopDispatcherNamesEveryLayout(t *testing.T) {
	dispatcher := desktopClientAuthShell()
	for _, layout := range authThemes {
		if layout == "split-static" {
			continue // the default arm, named in a comment beside it
		}
		if !strings.Contains(dispatcher, `case "`+layout+`":`) {
			t.Errorf("the desktop AuthShell never names %q, so it falls through silently", layout)
		}
	}
}

// A theme has to repaint the dashboard too, not only the sign-in page.
//
// [data-theme="<name>"] in globals.css is what repaints the admin. Without a
// block, a theme styles its auth screens and leaves the dashboard on the
// default palette: coral sign-in, atlas dashboard, and nothing anywhere says
// so. The picker's own description promises a theme "drives dashboard tokens".
func TestEveryThemeRepaintsTheDashboard(t *testing.T) {
	css := adminScreenCSS()
	// The token every surface reads. A block that omits it is worse than no
	// block, because the theme half-applies.
	for _, theme := range ValidThemes {
		block := `[data-theme="` + theme + `"]`
		if !strings.Contains(css, block) {
			t.Errorf("theme %q has no %s block, so its dashboard stays on the default palette", theme, block)
			continue
		}
		body := css[strings.Index(css, block):]
		if end := strings.Index(body, "}"); end > 0 {
			body = body[:end]
		}
		for _, token := range []string{"--bg-primary", "--text-primary", "--accent", "--accent-fg"} {
			if !strings.Contains(body, token) {
				t.Errorf("theme %q sets no %s, so that surface falls back to another theme's colour", theme, token)
			}
		}
	}
}
