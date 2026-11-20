package scaffold

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Three files decide what "@" means in an admin, and they have to agree.
//
// grit#229: tsconfig.json mapped @/* to ./src/*, vite.config.ts resolved @ to
// ./src, and vitest.config.ts resolved it to the app root. So the app ran, tsc
// passed, and a test that imported app code failed before it started with
// "Failed to resolve import @/lib/utils". Reported with a one-line
// reproduction, which is the only reason it was ever found: nothing in a green
// suite points at it.
//
// This reads all three in one test, because a test per file passes while the
// set disagrees, which is the whole shape of the bug.
func TestTheAdminsAliasesAgree(t *testing.T) {
	// The spelling each file uses, and what it means.
	//
	//	tsconfig:  "@/*": ["./src/*"]   or ["./*"]
	//	vite:      '@': path.resolve(__dirname, './src')
	//	vitest:    "@": resolve(__dirname, "./src")
	tsconfigAlias := regexp.MustCompile(`"@/\*":\s*\[\s*"\./(?:(src)/)?\*"`)
	viteAlias := regexp.MustCompile(`['"]@['"]:\s*path\.resolve\(__dirname,\s*['"]\./(src)?['"]?\)?`)
	vitestAlias := regexp.MustCompile(`"@":\s*resolve\(__dirname,\s*"\.(?:/(src))?"\)`)

	for _, frontend := range []Frontend{FrontendNext, FrontendTanStack} {
		opts := Options{Architecture: ArchTriple, Frontend: frontend}

		vitest := adminVitestConfig(opts)
		match := vitestAlias.FindStringSubmatch(vitest)
		if match == nil {
			t.Fatalf("%s: no @ alias found in the admin's vitest config, so this "+
				"guard is watching nothing", frontend)
		}
		got := match[1] // "src" or ""

		want := ""
		if frontend == FrontendTanStack {
			want = "src"
		}
		if got != want {
			t.Errorf("%s: the vitest alias points @ at %q, and the app points it at %q. "+
				"A test importing @/lib/utils then resolves a different file from the "+
				"one the component gets, or none at all.",
				frontend, "./"+got, "./"+want)
		}
	}

	// And the claim above about what the app says: read from the templates
	// rather than asserted from memory, because the day one of them moves is
	// the day this test has to fail.
	for name, body := range map[string]string{
		"the Vite admin's tsconfig":    adminTanStackTSConfig(),
		"the Vite admin's vite.config": adminTanStackViteConfig(),
	} {
		var found string
		if m := tsconfigAlias.FindStringSubmatch(body); m != nil {
			found = m[1]
		} else if m := viteAlias.FindStringSubmatch(body); m != nil {
			found = m[1]
		} else {
			t.Errorf("%s names no @ alias, so the vitest config above is matched "+
				"against an assumption rather than against the app", name)
			continue
		}
		if found != "src" {
			t.Errorf("%s points @ at %q, not at ./src, so the vitest config is now "+
				"the odd one out in the other direction", name, "./"+found)
		}
	}
}

// Every admin test imports lib the way every other file in that admin does.
//
// The TanStack tests used to say "@/src/lib", which resolved only because the
// vitest alias pointed at the app root. That is a workaround for grit#229
// rather than a convention, and leaving it in place after the fix would mean
// one admin whose tests import differently from its own components.
func TestNoAdminTestUsesTheSrcWorkaround(t *testing.T) {
	// Comments stripped, because the fix's own explanation names the spelling
	// it removed: a Contains check over the whole file matches the prose about
	// the bug and fails for ever after. Which this guard did on its first run.
	writer := stripGoComments(read(t, "frontend_test_files.go"))
	if strings.Contains(writer, `"@/src/lib"`) {
		t.Error(`the test writer still spells an import "@/src/lib". With the ` +
			`vitest alias fixed, "@/lib" is what the app uses and what resolves.`)
	}
}

// The component tests reach both admins.
//
// Eight of them were written for the Next admin alone, because of grit#229,
// and the gap was invisible: both suites were green, one just had less in it.
func TestBothAdminsGetTheComponentTests(t *testing.T) {
	for _, frontend := range []Frontend{FrontendNext, FrontendTanStack} {
		root := t.TempDir()
		opts := Options{Architecture: ArchTriple, Frontend: frontend}
		if err := writeFrontendTestFiles(root, opts); err != nil {
			t.Fatal(err)
		}

		for _, name := range []string{
			"data-table.test.tsx",
			"table-pagination.test.tsx",
			"form-schema.test.tsx",
			"if-match.test.tsx",
			"dropzone.test.tsx",
			"field-inputs.test.tsx",
			"conditional-fields.test.tsx",
			"password-strength.test.tsx",
		} {
			path := filepath.Join(root, "apps", "admin", "__tests__", name)
			if !fileExists(path) {
				t.Errorf("%s: %s was not written", frontend, name)
			}
		}

		// The one that genuinely belongs to one admin: proxy.ts is the Next
		// admin's edge gate and the Vite admin guards routes in the router.
		proxy := filepath.Join(root, "apps", "admin", "__tests__", "proxy.test.ts")
		if frontend == FrontendTanStack && fileExists(proxy) {
			t.Error("the Vite admin got a test for a proxy.ts it does not have")
		}
		if frontend == FrontendNext && !fileExists(proxy) {
			t.Error("the Next admin lost its proxy test")
		}
	}
}
