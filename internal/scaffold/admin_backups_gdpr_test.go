package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// Every icon any admin file imports is a named export of lib/icons.
//
// lib/icons has three blocks: the import from lucide-react, the iconMap a
// resource definition looks a name up in, and the export block a page imports
// from. Being in the first two does not put a name in the third, and the
// failure is a build error in a generated project, which is to say in somebody
// else's repository rather than in CI here.
//
// It has now happened at least three times (Workflow, Info, Package), each
// time found by standing a project up rather than by a test, so here is the
// test. It reads the names every admin file actually asks for instead of a
// list somebody has to remember to extend.
func TestEveryIconAnAdminFileImportsIsExported(t *testing.T) {
	icons := adminIconMap()
	exportAt := strings.Index(icons, "export {")
	if exportAt < 0 {
		t.Fatal("lib/icons has no export block")
	}
	exports := icons[exportAt:]

	// A name counts as exported either by that block or by its own
	// declaration: getIcon is a function this file exports directly, and it is
	// imported from here like any icon.
	exported := func(name string) bool {
		return strings.Contains(exports, "\n  "+name+",") ||
			strings.Contains(icons, "export function "+name+"(") ||
			strings.Contains(icons, "export const "+name+" ") ||
			strings.Contains(icons, "export type "+name+" ")
	}

	opts := Options{ProjectName: "shop", Version: "9.9.9", Architecture: ArchTriple, Frontend: FrontendNext}
	root := t.TempDir()
	files := adminFileMap(root, opts)

	// The map is not the whole admin: four writers of their own put files
	// beside it, and one of those may import an icon too.
	for _, write := range []func(string, Options) error{
		writeAdminResourceDashboardWidgets,
		writeAdminCustomChartFiles,
		writeAdminPasskeyFiles,
		writePasswordStrengthFiles,
	} {
		if err := write(root, opts); err != nil {
			t.Fatalf("running an admin writer: %v", err)
		}
	}
	_ = filepath.Walk(filepath.Join(root, "apps", "admin"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if _, already := files[path]; !already {
			if b, readErr := os.ReadFile(path); readErr == nil {
				files[path] = string(b)
			}
		}
		return nil
	})

	importBlock := regexp.MustCompile(`import\s*\{([^}]*)\}\s*from\s*"@/lib/icons"`)
	missing := map[string][]string{}

	for path, body := range files {
		for _, m := range importBlock.FindAllStringSubmatch(body, -1) {
			for _, raw := range strings.Split(m[1], ",") {
				name := strings.TrimSpace(raw)
				// "X as Y" imports X; a type-only import is not a value.
				if i := strings.Index(name, " as "); i >= 0 {
					name = strings.TrimSpace(name[:i])
				}
				name = strings.TrimPrefix(name, "type ")
				name = strings.TrimSpace(name)
				if name == "" || strings.HasPrefix(name, "//") {
					continue
				}
				if !exported(name) {
					missing[name] = append(missing[name], path)
				}
			}
		}
	}

	if len(missing) == 0 {
		return
	}
	names := make([]string, 0, len(missing))
	for name := range missing {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		t.Errorf("%q is imported from @/lib/icons but is not in its export block, so the "+
			"generated project will not compile. First seen in %s", name, missing[name][0])
	}
}

// The retention number has one source.
//
// The backups page says "3 of 4 kept", and that 4 is backup.KeepBackups in Go.
// Typing it into the page as well makes two numbers that have to agree with no
// way to notice when they stop: raise retention to eight and the screen goes
// on saying four, confidently and wrongly. It travels with the schedule
// instead.
func TestTheBackupRetentionCountComesFromTheServer(t *testing.T) {
	api := backupHandlerGo()
	admin := adminUseBackups()
	page := adminBackupsPage()

	if !strings.Contains(api, `gin.H{"keep": backup.KeepBackups}`) {
		t.Error("GetSettings no longer sends the retention count, so the admin has nothing " +
			"to read and will have to hardcode it")
	}
	if !strings.Contains(admin, `.data.meta?.keep ?? 0`) {
		t.Error("the admin no longer reads the retention count from the response")
	}
	if !strings.Contains(page, `keep ? held + " of " + keep + " kept"`) {
		t.Error("the snapshots card does not show the retention window from the server")
	}
	// The number itself must not appear as a literal on the page.
	for _, bad := range []string{`"of 4"`, `of 4 kept`, `KeepBackups = 4`} {
		if strings.Contains(page, bad) {
			t.Errorf("the backups page hardcodes the retention count (%q)", bad)
		}
	}
}

// Erasing a person needs their address typed out.
//
// A reason field is a record, not a confirmation: it can be filled in without
// ever reading whose account is about to go. This page has exactly one
// catastrophic mistake available to it, which is erasing the wrong person, and
// it is irreversible by design. So the address is typed, and the reason is
// kept as well because the API requires it and the journal records it.
func TestErasingAPersonRequiresTypingTheirEmail(t *testing.T) {
	page := adminGDPRPage()

	if !strings.Contains(page, "const emailMatches =") {
		t.Fatal("the erase confirmation no longer compares a typed address")
	}
	if !strings.Contains(page, "typed.trim().toLowerCase() === picked.email.trim().toLowerCase()") {
		t.Error("the typed address is not compared to the picked user's address")
	}
	if !strings.Contains(page, "disabled={!emailMatches || reason.trim().length < 3}") {
		t.Error("the erase button is not gated on both the typed address and a reason; " +
			"the API refuses an erasure with no reason, and the address is the confirmation")
	}
	if !strings.Contains(page, "The email has to match exactly.") {
		t.Error("nothing says why the erase button is disabled")
	}
}

// Both pages are drawn in the project's theme, not in Tailwind's palette.
//
// The GDPR page was written against red-500 and emerald-500 directly. Those
// are fixed hues: they do not move when a project picks a theme, so on seven
// of the eight themes this page's destructive button was a different red from
// every other destructive button in the panel, and nothing in the build said
// so. A hardcoded palette colour is the one styling mistake that is invisible
// in the theme it was written in.
func TestTheBackupAndGDPRPagesUseTheProjectsOwnColours(t *testing.T) {
	pages := map[string]string{
		"Backups": adminBackupsPage(),
		"GDPR":    adminGDPRPage(),
	}
	// Tailwind's own scales. The project's colours are tokens with no number.
	palette := regexp.MustCompile(
		`\b(?:bg|text|border|ring|from|to|via)-(?:red|green|emerald|blue|amber|yellow|orange|` +
			`purple|violet|indigo|pink|rose|teal|cyan|sky|lime|slate|gray|zinc|neutral|stone)-\d{2,3}\b`)

	for name, page := range pages {
		if hits := palette.FindAllString(page, -1); len(hits) > 0 {
			t.Errorf("%s uses Tailwind palette colours %v instead of the project's tokens, "+
				"so it stays that hue on every theme", name, dedupe(hits))
		}
		// And the softer version of the same mistake.
		for _, bad := range []string{"bg-success/", "bg-danger/", "bg-warning/", "bg-info/", "bg-accent/"} {
			if strings.Contains(page, bad) {
				t.Errorf("%s mixes its own tint with %q; use the soft and line tokens", name, bad)
			}
		}
	}
}

func dedupe(in []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}
