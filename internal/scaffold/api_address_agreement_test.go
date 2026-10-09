package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// Where the frontend thinks the API is.
//
// This went wrong five times in one sweep, once per tier, and each time the
// symptom pointed somewhere else:
//
//	Next      the CSP read NEXT_PUBLIC_API_URL and the client read API_URL,
//	          so moving APP_PORT blocked every request as a console violation
//	Vite      the CSP was right and the client fell back to 8080, so sign-in
//	          was blocked on any project that had moved the port
//	Expo      lib/api.ts held the literal 8080 and the seeder wrote names an
//	          Expo app cannot read, so the phone called a dead port
//	desktop   api-client.ts held the literal 8080 and the dev proxy repeated
//	          it, so sign-in answered 500 from a different project's server
//	single    the server-side blog client fell back to 8080 while everything
//	          else in the same app followed API_URL
//
// The shape is always the same: the address is written down in several places,
// some of them literals, and only some of them move. These tests are the
// standing answer, so the sixth one fails here rather than in somebody's
// browser.

// envReadsIn pulls the environment variable names out of an expression, in
// order.
var envReadsIn = regexp.MustCompile(`(?:process\.env|import\.meta\.env)\.([A-Z_][A-Z0-9_]*)`)

// apiBaseExpression returns the text of a declaration that resolves the API's
// address, found by the name it is assigned to.
func apiBaseExpression(t *testing.T, name, src, decl string) string {
	t.Helper()
	i := strings.Index(src, decl)
	if i < 0 {
		t.Fatalf("%s declares no %s", name, decl)
	}
	rest := src[i:]
	// To the end of the statement, which is the first ";" at depth zero. Close
	// enough: none of these expressions contains a string with a semicolon.
	if j := strings.Index(rest, ";"); j >= 0 {
		return rest[:j]
	}
	return rest
}

// Every browser-facing base reads API_URL somewhere in its chain, because that
// is the name the generated root .env sets.
func TestEveryNextAPIBaseReadsAPIURL(t *testing.T) {
	opts := Options{ProjectName: "app", Frontend: FrontendNext}

	for name, src := range map[string]string{
		"apps/admin/next.config.ts": adminNextConfig(opts),
		"apps/web/next.config.ts":   webNextConfig(opts),
	} {
		expr := apiBaseExpression(t, name, src, "const API_ORIGIN = toOrigin(")
		if !strings.Contains(expr, "process.env.API_URL") {
			t.Errorf("%s: the Content-Security-Policy's origin does not read API_URL, so it "+
				"authorises a different host from the one the client calls", name)
		}
	}

	// The server-rendered blog client, which is a separate file with its own
	// chain and was the one that still said 8080.
	blog := webBlogAPILib()
	expr := apiBaseExpression(t, "lib/blog-api.ts", blog, "const API_URL = (")
	if !strings.Contains(expr, "process.env.API_URL") {
		names := envReadsIn.FindAllStringSubmatch(expr, -1)
		var got []string
		for _, m := range names {
			got = append(got, m[1])
		}
		t.Errorf("lib/blog-api.ts reads %v and not API_URL, so a project that moves its port "+
			"has one file calling 8080 while every other file follows the setting", got)
	}
}

// A Vite app's config resolves the address once and hands it to both the proxy
// and the bundle, so the two cannot disagree.
func TestEveryViteConfigResolvesTheAddressOnce(t *testing.T) {
	viteOpts := Options{ProjectName: "app", Frontend: FrontendTanStack}

	for name, src := range map[string]string{
		"apps/admin/vite.config.ts":            adminTanStackViteConfig(),
		"apps/web/vite.config.ts":              webTanStackViteConfig(viteOpts),
		"apps/desktop/frontend/vite.config.ts": desktopClientViteConfig(),
	} {
		if !strings.Contains(src, "APP_PORT") {
			t.Errorf("%s does not read APP_PORT, so nothing follows the API when it moves", name)
		}
		if !strings.Contains(src, "define: clientEnv") {
			t.Errorf("%s does not hand the resolved address to the browser bundle", name)
		}
		if !strings.Contains(src, "target: apiTarget") {
			t.Errorf("%s's dev proxy does not use the resolved address", name)
		}
	}
}

// And no client module carries a port of its own. A literal here is invisible
// until somebody moves the API, and then it answers from whatever else happens
// to be listening.
func TestNoClientModuleCarriesItsOwnPort(t *testing.T) {
	opts := Options{ProjectName: "app", Frontend: FrontendNext}

	// Each of these resolves the API's address for one app. A fallback of
	// "http://localhost:8080" is allowed only where the expression has already
	// read the configured names first.
	for name, src := range map[string]string{
		"apps/expo/lib/api.ts":                    expoAPIClient(),
		"apps/desktop/frontend/lib/api-client.ts": desktopClientApiClientTS(),
		"lib/api-core.ts":                         apiCoreTS(opts),
	} {
		for _, line := range strings.Split(src, "\n") {
			if !strings.Contains(line, "localhost:8080") {
				continue
			}
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "*") {
				continue
			}
			// A fallback at the end of a chain that reads the environment
			// first is fine. A bare assignment is not.
			if envReadsIn.MatchString(line) || strings.Contains(line, "API_PORT") {
				continue
			}
			t.Errorf("%s names a port without reading the environment first:\n  %s", name, trimmed)
		}
	}
}

// A single Next project's frontend could not reach its API in development.
//
// In this shape the Go binary embeds the built frontend and serves it, so the
// browser is right to call its own origin, and api-core.ts does. In
// development they are two processes, Next on 3000 and the API on APP_PORT,
// and nothing joined them: every page that loaded data got a 404 from Next.
// The Vite sibling of the same architecture has had a dev proxy since it was
// written.
func TestASingleNextProjectBridgesToItsAPI(t *testing.T) {
	single := webNextConfig(Options{
		ProjectName: "app", Architecture: ArchSingle, Frontend: FrontendNext,
	})

	if !strings.Contains(single, "async rewrites()") {
		t.Fatal("a single Next project has no rewrite to its API, so the browser's " +
			"same-origin calls 404 against Next in development")
	}
	for _, path := range []string{"/api/:path*", "/studio/:path*", "/files/:path*"} {
		if !strings.Contains(single, path) {
			t.Errorf("the rewrite does not cover %s", path)
		}
	}
	if !strings.Contains(single, "process.env.APP_PORT") {
		t.Error("the rewrite does not follow APP_PORT, so it breaks the moment the API moves")
	}

	// A monorepo frontend calls the API cross-origin by its configured URL,
	// which is what the CSP authorises. A rewrite there would be a second
	// answer to a question that already has one.
	mono := webNextConfig(Options{
		ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext,
	})
	if strings.Contains(mono, "async rewrites()") {
		t.Error("a monorepo frontend has a rewrite it does not need")
	}
}

// And it reads its own .env.
//
// The hoist path was "../../", which is the project root from apps/web and two
// directories above the project in a single one. THEME, SOCIAL_AUTH_ENABLED
// and API_URL were all invisible to that shape.
func TestTheEnvHoistPointsAtTheProjectRoot(t *testing.T) {
	single := webNextConfig(Options{
		ProjectName: "app", Architecture: ArchSingle, Frontend: FrontendNext,
	})
	if strings.Contains(single, `resolve(process.cwd(), "..", "..", ".env")`) {
		t.Error("a single project looks for .env two directories above itself")
	}
	if !strings.Contains(single, `resolve(process.cwd(), ".", ".env")`) {
		t.Error("a single project does not read the .env beside it")
	}

	mono := webNextConfig(Options{
		ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext,
	})
	if !strings.Contains(mono, `resolve(process.cwd(), "..", "..", ".env")`) {
		t.Error("a monorepo app no longer reads the project root's .env")
	}
}
