package scaffold

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// Bugs reported against a real project, each asserted on the template text
// that carries it. The project code is a string here rather than something
// this binary can run, so these read the emitted source: narrower than a
// behaviour test and still enough to stop the exact thing coming back.

// grit#128: every generated /export answered 404 through the unversioned
// /api path, with the whole file in the body.
//
// The alias is a NoRoute handler, and gin sets the writer's status to 404 on
// the way in. HandleContext does not reset it, so a handler that streams
// bytes without calling WriteHeader flushes the leftover status. Handlers
// answering with c.JSON set their own and were fine, which is why only the
// downloads were broken.
func TestUnversionedAPIAliasClearsTheNoRouteStatus(t *testing.T) {
	src := apiRoutesGo()

	i := strings.Index(src, "func mountLegacyAPIAlias")
	if i < 0 {
		t.Fatal("mountLegacyAPIAlias is gone; this test no longer checks anything")
	}
	body := src[i:]
	if j := strings.Index(body, "\nfunc "); j > 0 {
		body = body[:j]
	}

	status := strings.Index(body, "c.Status(http.StatusOK)")
	handle := strings.Index(body, "r.HandleContext(c)")
	if status < 0 {
		t.Error("the alias does not reset the pending 404 before re-dispatching, " +
			"so a handler that writes bytes without WriteHeader answers 404 " +
			"with its body attached (grit#128)")
	}
	if status > 0 && handle > 0 && status > handle {
		t.Error("the status is reset after HandleContext, which is too late: " +
			"the handler has already written (grit#128)")
	}
}

// grit#124: no browser on another origin could use optimistic locking.
//
// If-Match was missing from the allow-list, so the preflight stripped it and
// the PATCH was never sent: no status, no log line, and a bare TypeError in
// the page. And ETag was not exposed, so even with the header allowed a
// client had no way to read the version it must send back.
func TestCORSCarriesTheOptimisticLockingHeaders(t *testing.T) {
	src := apiRoutesGo() + apiCorsMiddlewareGo()

	for _, want := range []string{`"If-Match"`, `"ETag"`} {
		if !strings.Contains(src, want) {
			t.Errorf("%s is missing from the CORS lists, so optimistic locking "+
				"cannot be used from a browser on another origin (grit#124)", want)
		}
	}
	if !strings.Contains(src, "Access-Control-Expose-Headers") {
		t.Error("no Access-Control-Expose-Headers, so cross-origin JavaScript " +
			"cannot read ETag off a response and has no version to send back " +
			"in If-Match (grit#124)")
	}
}

// grit#122: a Message-ID set by the caller was dropped by the SMTP driver and
// kept by the API drivers, so the same message got a different id depending
// on which one sent it. An app storing the id to match replies stored one
// that was never sent.
func TestBuildMIMEKeepsACallerMessageID(t *testing.T) {
	src := mailMIMEGo()

	if !strings.Contains(src, "resolveMessageID(m.Headers") {
		t.Error("buildMIME writes a generated Message-ID without looking at " +
			"the one the caller set (grit#122)")
	}
	if !strings.Contains(src, "func validMessageID(") {
		t.Error("a caller Message-ID is used without being checked; a value " +
			"carrying CR or LF is header injection (grit#122)")
	}
}

// grit#126: on the desktop app every accent tint rendered transparent and
// every primary button had dark text on a dark fill.
//
// Tailwind 3 emits no rule at all for bg-accent/10 when the colour is a bare
// CSS variable, and accent-fg was used in three places and never defined.
var desktopTinted = regexp.MustCompile(
	`(?:bg|text|border|ring|from|to|via|fill|stroke|divide|outline)-` +
		`(accent|accent-hover|surface-hover|foreground|foreground-secondary|foreground-muted)/\d`)

func TestDesktopTintedColoursCanTakeAnOpacity(t *testing.T) {
	src := desktopClientTailwindConfig()

	tinted := map[string]bool{}
	for _, m := range desktopTinted.FindAllStringSubmatch(desktopClientFileSource(t), -1) {
		tinted[m[1]] = true
	}
	if len(tinted) == 0 {
		t.Fatal("no tinted colour classes found in the desktop client; the " +
			"reader stopped matching and this test checks nothing")
	}

	for name := range tinted {
		// The colour has to be declared in a form Tailwind can put an alpha
		// into: color-mix with <alpha-value>, or a plain hex.
		decl := declarationFor(src, name)
		if decl == "" {
			t.Errorf("%s is tinted in the desktop client and not declared in "+
				"tailwind.config.ts", name)
			continue
		}
		if strings.Contains(decl, "var(--") && !strings.Contains(decl, "<alpha-value>") {
			t.Errorf("%s is declared as a bare CSS variable, so Tailwind emits "+
				"no rule for its tints and every one of them renders as "+
				"nothing (grit#126): %s", name, decl)
		}
	}

	if !strings.Contains(src, `"accent-fg"`) {
		t.Error("accent-fg is used on accent fills and not declared, so every " +
			"primary button keeps the inherited dark text (grit#126)")
	}
}

// declarationFor returns the config line declaring a colour, or "".
func declarationFor(config, name string) string {
	for _, line := range strings.Split(config, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, `"`+name+`":`) ||
			strings.HasPrefix(trimmed, name+":") {
			return trimmed
		}
	}
	return ""
}

// desktopClientFileSource is the desktop client's template text: the Go file
// that holds it, read from disk. Reading the source rather than calling each
// template function keeps this honest when a new screen is added.
func desktopClientFileSource(t *testing.T) string {
	t.Helper()
	for _, name := range []string{
		"desktop_client_files.go",
		"desktop_frontend_config.go",
	} {
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("reading %s: %v", name, err)
		}
		desktopSource += string(b)
	}
	return desktopSource
}

var desktopSource string

// grit#127: enterprise SSO had no login hook, so a project had no supported
// place to act on an SSO sign-in. In a multi-tenant app the obvious thing to
// do there is put the person in the organisation the connection belongs to,
// and applyGroupRoles only writes platform roles.
//
// A registry nothing calls is the failure mode worth guarding: it compiles,
// it is documented, and no hook ever runs.
func TestBothSSOCallbacksRunTheLoginHooks(t *testing.T) {
	for name, src := range map[string]string{
		"sso.go (OIDC callback)": tmpl("api/handlers/sso.go"),
		"saml.go (SAML ACS)":     tmpl("api/handlers/saml.go"),
	} {
		if !strings.Contains(src, "services.RunSSOLoginHooks(") {
			t.Errorf("%s does not run the SSO login hooks, so the registry in "+
				"services/sso_hooks.go exists and nothing calls it (grit#127)", name)
			continue
		}
		// After the group roles: a hook that reads the person's roles to
		// decide an organisation role needs them applied first.
		roles := strings.Index(src, "h.applyGroupRoles(")
		hooks := strings.Index(src, "services.RunSSOLoginHooks(")
		if roles >= 0 && hooks < roles {
			t.Errorf("%s runs the login hooks before applyGroupRoles, so a hook "+
				"reading the person's roles sees them unapplied (grit#127)", name)
		}
	}

	hooks := apiSSOHooksGo()
	for _, want := range []string{
		"func OnSSOLogin(",
		"func RunSSOLoginHooks(",
		"func ResetSSOLoginHooks(",
		"conn *models.SSOConnection",
		"groups []string",
	} {
		if !strings.Contains(hooks, want) {
			t.Errorf("services/sso_hooks.go is missing %s", want)
		}
	}
}

// grit#137: the SQLite test database did not enforce foreign keys, so a
// generated test that inserted a child with an empty parent id passed there
// and failed on Postgres. The test suite was green on the database nobody
// ships and red on the one everybody does.
func TestSQLiteTestDatabaseEnforcesForeignKeys(t *testing.T) {
	src := apiTestKitGo()

	if !strings.Contains(src, "PRAGMA foreign_keys = ON") {
		t.Fatal("the SQLite test database does not enforce foreign keys, so a " +
			"broken one passes locally and fails on Postgres (grit#137)")
	}

	// Before AutoMigrate, so the tables are built with it already on.
	pragma := strings.Index(src, "PRAGMA foreign_keys = ON")
	migrate := strings.Index(src, "db.AutoMigrate(models.Models()...)")
	if migrate > 0 && pragma > migrate {
		t.Error("foreign keys are enabled after the tables are migrated")
	}
}
