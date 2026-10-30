package scaffold

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// authPageURLs reads the panel's own file map for the pages in the (auth)
// route group and returns the URLs Next serves them at. A route group is not
// part of the URL, so app/(auth)/sign-up/page.tsx is /sign-up.
func authPageURLs(t *testing.T, root string, files map[string]string) []string {
	t.Helper()
	adminRoot := filepath.Join(root, "apps", "admin")
	seen := map[string]bool{}
	for path := range files {
		rel, err := filepath.Rel(adminRoot, path)
		if err != nil {
			continue
		}
		rel = filepath.ToSlash(rel)
		if !strings.HasPrefix(rel, "app/(auth)/") || !strings.HasSuffix(rel, "/page.tsx") {
			continue
		}
		route := strings.TrimSuffix(strings.TrimPrefix(rel, "app/(auth)"), "/page.tsx")
		if route == "" {
			continue
		}
		seen[route] = true
	}
	out := make([]string, 0, len(seen))
	for url := range seen {
		out = append(out, url)
	}
	sort.Strings(out)
	return out
}

// Review 2.15. The gate denies by default, which is the right default and
// means the one way to get it wrong is to leave an auth page out of
// PUBLIC_PATHS: the page that lets somebody sign in would redirect to the page
// that lets somebody sign in.
//
// So the list is not trusted, it is compared against the pages the panel
// actually has. Both sides are read here, because a test per side passes while
// the pair is broken.
func TestEveryAuthPageIsPublic(t *testing.T) {
	root := t.TempDir()
	files := adminFileMap(root, tripleOptions())

	urls := authPageURLs(t, root, files)
	if len(urls) < 5 {
		t.Fatalf("found only %d pages in the (auth) group (%v): the test is reading the wrong thing", len(urls), urls)
	}

	gate := files[filepath.Join(root, "apps", "admin", "proxy.ts")]
	if gate == "" {
		t.Fatal("apps/admin has no proxy.ts, so every panel page is served to a signed-out visitor")
	}
	for _, url := range urls {
		if !strings.Contains(gate, `"`+url+`"`) {
			t.Errorf("%s is an auth page and is not in PUBLIC_PATHS: a visitor sent there would be sent there again", url)
		}
	}
}

// matcherCoverage reads the matcher out of proxy.ts and answers, for a
// path, whether Next would run the middleware on it.
//
// The matcher is "/((?!a|b).*)": everything except the alternatives inside the
// negative lookahead. Go's regexp has no lookahead, so the alternatives are
// read out and each is tried on its own, anchored where the lookahead sits.
// What matters is the list of exclusions, and this reads the real one rather
// than a copy of it.
func matcherCoverage(t *testing.T, gate string) func(string) bool {
	t.Helper()
	const head = `matcher: ["/((?!`
	const tail = `).*)"]`
	start := strings.Index(gate, head)
	end := strings.Index(gate, tail)
	if start < 0 || end < start {
		t.Fatalf("proxy.ts has no matcher of the expected shape; its matcher line reads:\n%s",
			lineContaining(gate, "matcher:"))
	}

	// The matcher is a JavaScript string literal, so JavaScript unescapes it
	// before the regex engine ever sees it. Reading it as Go source would make
	// this test agree with a pattern the browser disagrees with: a lone
	// backslash in the template survives Go untouched and is dropped by
	// JavaScript, which turns an escaped dot into "any character" and takes
	// every panel page out of the matcher.
	inner := jsUnescape(gate[start+len(head) : end])
	alternatives := strings.Split(inner, "|")
	if len(alternatives) < 2 {
		t.Errorf("the matcher excludes only %v, so either Next's own output or the public/ files are gated", alternatives)
	}
	excluded := make([]*regexp.Regexp, 0, len(alternatives))
	for _, alt := range alternatives {
		re, err := regexp.Compile("^" + alt)
		if err != nil {
			t.Fatalf("the matcher's exclusion %q is not a valid regex: %v", alt, err)
		}
		excluded = append(excluded, re)
	}
	return func(path string) bool {
		rest := strings.TrimPrefix(path, "/")
		for _, re := range excluded {
			if re.MatchString(rest) {
				return false
			}
		}
		return true
	}
}

// jsUnescape reads a double-quoted JavaScript string literal's contents the
// way an engine does: a backslash pair is one backslash, and a backslash
// before anything else is dropped. Enough for a regex written in a string;
// \n and \uXXXX have no business in one.
func jsUnescape(literal string) string {
	var out strings.Builder
	for i := 0; i < len(literal); i++ {
		if literal[i] == '\\' && i+1 < len(literal) {
			i++
		}
		out.WriteByte(literal[i])
	}
	return out.String()
}

func lineContaining(body, needle string) string {
	for _, line := range strings.Split(body, "\n") {
		if strings.Contains(line, needle) {
			return line
		}
	}
	return "(nothing)"
}

// The gate has to let Next's own output and the public/ files through, or the
// login page a stranger is redirected to renders without its logo.
func TestTheGateLetsAssetsThrough(t *testing.T) {
	root := t.TempDir()
	gate := adminFileMap(root, tripleOptions())[filepath.Join(root, "apps", "admin", "proxy.ts")]
	covers := matcherCoverage(t, gate)

	for _, path := range []string{"/dashboard", "/system/api-keys", "/resources/posts", "/"} {
		if !covers(path) {
			t.Errorf("%s is a panel page and the matcher does not cover it, so it has no gate", path)
		}
	}
	for _, path := range []string{"/_next/static/chunk.js", "/_next/image", "/grit_logo.png", "/favicon.ico"} {
		if covers(path) {
			t.Errorf("%s is an asset and the matcher covers it, so a signed-out visitor's login page loads without it", path)
		}
	}
}

// The name of the cookie is declared once and imported, rather than written in
// both files. Two copies of a cookie name is the shape of bug this project
// keeps finding: the gate reads one spelling and the thing that opens it
// writes another, and each file is correct read on its own.
func TestTheGateAndTheMarkerShareOneCookieName(t *testing.T) {
	root := t.TempDir()
	files := adminFileMap(root, tripleOptions())
	admin := filepath.Join(root, "apps", "admin")

	lib := files[filepath.Join(admin, "lib", "admin-session.ts")]
	if lib == "" {
		t.Fatal("apps/admin/lib/admin-session.ts is not written")
	}
	if !strings.Contains(lib, `ADMIN_SESSION_COOKIE = "grit_admin_session"`) {
		t.Error("the marker's name is not declared in lib/admin-session.ts")
	}

	gate := files[filepath.Join(admin, "proxy.ts")]
	if !strings.Contains(gate, `from "@/lib/admin-session"`) {
		t.Error("proxy.ts does not import the cookie names, so it holds a second copy of them")
	}
	if strings.Contains(gate, `"grit_admin_session"`) {
		t.Error("proxy.ts spells the cookie name out as well as importing it")
	}

	client := files[filepath.Join(admin, "lib", "api-client.ts")]
	if !strings.Contains(client, `from "@/lib/admin-session"`) {
		t.Error("api-client.ts does not import the marker, so nothing ever writes the cookie the gate reads")
	}
	for _, call := range []string{"setAdminSession()", "clearAdminSession()"} {
		if !strings.Contains(client, call) {
			t.Errorf("api-client.ts never calls %s", call)
		}
	}
}

// Every sign-in route the panel has answers through the api client, and the
// marker is stamped from one list there rather than from each caller. A route
// missing from the list is a page that signs you in and then bounces you to
// the login, so the list is checked against the endpoints the panel posts to.
func TestEverySignInEndpointStampsTheMarker(t *testing.T) {
	root := t.TempDir()
	files := adminFileMap(root, tripleOptions())
	client := files[filepath.Join(root, "apps", "admin", "lib", "api-client.ts")]

	for _, endpoint := range []string{
		"/auth/login",
		"/auth/register",
		"/auth/refresh",
		"/auth/me",
		"/auth/totp/verify",
		"/auth/totp/backup-codes/verify",
		"/auth/magic-link/consume",
		"/auth/passkeys/login/finish",
	} {
		if !strings.Contains(client, `"`+endpoint+`"`) {
			t.Errorf("%s hands back a session and is not in SESSION_PATHS, so signing in that way leaves no marker and the next page redirects to the login", endpoint)
		}
	}
	// A challenge is not a session.
	if !strings.Contains(client, "totp_required") {
		t.Error("the marker is stamped on a 2FA challenge, which opens the panel to a browser that stopped at the second factor")
	}
}

// The panel inside a web app has a gate already, in that app's own
// proxy.ts, and a Next app reads one only at its root. A second
// file under admin-panel/ would be dead, and one written to the web app's root
// would replace the gate it already has.
func TestTheEmbeddedPanelDoesNotBringAMiddleware(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "app", Architecture: ArchDouble, Frontend: FrontendNext}
	for path := range embeddedAdminFileMap(root, opts) {
		if filepath.Base(path) == "proxy.ts" {
			t.Errorf("the embedded panel writes %s", path)
		}
	}
}

// The Vite admin's api-client is the Next.js one with the environment
// rewritten, so it imports the marker too. Without the file it does not build.
func TestTheViteAdminHasTheMarkerItImports(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendTanStack}
	files := adminTanStackFileMap(root, opts)
	lib := filepath.Join(root, "apps", "admin", "src", "lib", "admin-session.ts")
	if _, ok := files[lib]; !ok {
		t.Error("the Vite admin's api-client imports @/lib/admin-session and the file is not written")
	}
}
