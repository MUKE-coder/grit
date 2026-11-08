package scaffold

import (
	"strings"
	"testing"
)

// grit#187: a brand-new project failed its own govulncheck gate on the first
// push, on a called vulnerability in crewjam/saml reached from
// POST /api/v1/auth/saml/:slug/acs, with Fixed in: N/A.
//
// Two halves, because neither is sufficient. govulncheck is static, so the
// flag does not make the gate green; and the allowlist entry does not remove a
// public endpoint from a project that will never use it.

// TestSSOIsAModuleAndDefaultsOff. The only module that defaults off, which is
// the whole point and the thing a careless edit would undo.
func TestSSOIsAModuleAndDefaultsOff(t *testing.T) {
	cfg := apiConfigGo()

	if !strings.Contains(cfg, "\tSSO bool\n") {
		t.Error("ModuleFlags has no SSO field, so there is no flag to gate the SAML routes with")
	}
	if !strings.Contains(cfg, `SSO: getEnv("MODULE_SSO", "false") == "true"`) {
		t.Error(`SSO does not default to false. It is the only module whose cost is a public ` +
			`endpoint: POST /auth/saml/:slug/acs parses an assertion through crewjam/saml, ` +
			`which carries a called advisory with no fixed version`)
	}
	// The other flags stay as they were. A sweep that turned them all off would
	// pass the line above and break every project.
	for _, on := range []string{"MODULE_AI", "MODULE_JOBS", "MODULE_REALTIME", "MODULE_TWOFACTOR"} {
		if !strings.Contains(cfg, on+`", "true"`) {
			t.Errorf("%s no longer defaults to true", on)
		}
	}
	// And it answers through both of the ways a module is asked about, or the
	// admin nav shows a screen whose routes are not mounted.
	if !strings.Contains(cfg, "case \"sso\":") {
		t.Error("ModuleFlags.Enabled does not answer for \"sso\"")
	}
	if !strings.Contains(cfg, `"sso":       m.SSO,`) {
		t.Error("ModuleFlags.Map omits sso, so /system/modules does not report it and the " +
			"admin nav cannot hide the screen")
	}
}

// TestSSORoutesAreGated: both groups, because /sso is the OIDC half and /saml
// is the half with the advisory, and leaving either mounted leaves a public
// auth endpoint on a project that does not use one.
func TestSSORoutesAreGated(t *testing.T) {
	routes := apiRoutesGo()

	for _, group := range []string{
		"\tif ssoEnabled {\n\t\tsso := auth.Group(\"/sso\")",
		"\tif ssoEnabled {\n\t\tsamlGroup := auth.Group(\"/saml\")",
	} {
		if !strings.Contains(routes, group) {
			t.Errorf("a route group is mounted unconditionally; wanted:\n%s", group)
		}
	}
	if !strings.Contains(routes, "mountAuthRoutes(v1, authHandler, ssoHandler, totpHandler, passkeyHandler, cfg.Modules.SSO)") {
		t.Error("mountAuthRoutes is not given the flag, so the guard reads a variable nobody set")
	}
	if !strings.Contains(routes, "passkeyHandler *handlers.PasskeyHandler, ssoEnabled bool) {") {
		t.Error("mountAuthRoutes does not take the flag")
	}
}

// TestRepairedSSOConfigIsTheGeneratedConfig.
//
// config.go is the developer's file, so an existing project gets the field from
// a repair rather than a rewrite. If the repair's text drifts from the
// template's, a repaired project and a new one disagree about the same field
// and the next reader cannot tell which is right.
func TestRepairedSSOConfigIsTheGeneratedConfig(t *testing.T) {
	want := apiConfigGo()

	// The same file as it stood before this change.
	old := strings.Replace(want, ssoField, "", 1)
	old = strings.Replace(old, ssoEnabledCase, "", 1)
	old = strings.Replace(old, ssoMapEntry, "", 1)
	old = strings.Replace(old, ssoRead, "", 1)
	if strings.Contains(old, "MODULE_SSO") {
		t.Fatal("could not build the pre-fix config.go, so this test asserts nothing")
	}

	got, changes, warnings := repairSSOModuleConfigSource(old)
	if len(warnings) != 0 {
		t.Errorf("unexpected warnings: %v", warnings)
	}
	if len(changes) == 0 {
		t.Error("the repair changed the file without saying what it did")
	}
	if got != want {
		t.Errorf("a repaired config.go is not the generated one.\n--- repaired ---\n%s\n--- generated ---\n%s",
			moduleFlagRegion(got), moduleFlagRegion(want))
	}

	// And again: a second upgrade finds nothing to do.
	twice, _, _ := repairSSOModuleConfigSource(got)
	if twice != got {
		t.Error("a second upgrade changed config.go again")
	}
}

// TestRepairedSSORoutesAreTheGeneratedRoutes, for the same reason.
func TestRepairedSSORoutesAreTheGeneratedRoutes(t *testing.T) {
	want := apiRoutesGo()

	old := strings.Replace(want, ssoRoutesNew, ssoRoutesOld, 1)
	old = strings.Replace(old, samlRoutesNew, samlRoutesOld, 1)
	old = strings.Replace(old,
		"passkeyHandler *handlers.PasskeyHandler, ssoEnabled bool) {",
		"passkeyHandler *handlers.PasskeyHandler) {", 1)
	old = strings.Replace(old,
		"mountAuthRoutes(v1, authHandler, ssoHandler, totpHandler, passkeyHandler, cfg.Modules.SSO)",
		"mountAuthRoutes(v1, authHandler, ssoHandler, totpHandler, passkeyHandler)", 1)
	if strings.Contains(old, "ssoEnabled") {
		t.Fatal("could not build the pre-fix routes.go, so this test asserts nothing")
	}

	got, changes, warnings := repairSSORoutesSource(old)
	if len(warnings) != 0 {
		t.Errorf("unexpected warnings: %v", warnings)
	}
	if len(changes) == 0 {
		t.Error("the repair changed the file without saying what it did")
	}
	if got != want {
		t.Error("a repaired routes.go is not the generated one")
	}
	twice, _, _ := repairSSORoutesSource(got)
	if twice != got {
		t.Error("a second upgrade changed routes.go again")
	}
}

// TestSSORepairWarnsOnAnEditedFile: a project that has rewritten either file
// gets told what is wrong, not a guess. The routes repair in particular must
// not half-apply: the signature and the call site are a pair, and changing one
// without the other is a project that does not compile.
func TestSSORepairWarnsOnAnEditedFile(t *testing.T) {
	want := apiRoutesGo()
	old := strings.Replace(want, ssoRoutesNew, ssoRoutesOld, 1)
	old = strings.Replace(old, samlRoutesNew, samlRoutesOld, 1)
	old = strings.Replace(old,
		"passkeyHandler *handlers.PasskeyHandler, ssoEnabled bool) {",
		"passkeyHandler *handlers.PasskeyHandler) {", 1)
	old = strings.Replace(old,
		"mountAuthRoutes(v1, authHandler, ssoHandler, totpHandler, passkeyHandler, cfg.Modules.SSO)",
		"mountAuthRoutes(v1, authHandler, ssoHandler, totpHandler, passkeyHandler)", 1)

	// Somebody has added a route to the SAML group.
	edited := strings.Replace(old, samlRoutesOld,
		"\tsamlGroup := auth.Group(\"/saml\")\n\t{\n\t\tsamlGroup.GET(\"/:slug/metadata\", ssoHandler.SAMLMetadata)\n\t}\n", 1)

	got, _, warnings := repairSSORoutesSource(edited)
	if got != edited {
		t.Error("the repair rewrote a routes.go it does not recognise")
	}
	var warned bool
	for _, w := range warnings {
		if strings.Contains(w, "MODULE_SSO") {
			warned = true
		}
	}
	if !warned {
		t.Errorf("the repair said nothing about a routes.go it could not gate: %v", warnings)
	}
}

// TestAllowlistShipsTheSAMLEntry. Every project inherits the same called
// vulnerability from the same framework code, so every project was going to
// write the same line, and until it did its CI was red with nothing it could do.
func TestAllowlistShipsTheSAMLEntry(t *testing.T) {
	allow := govulncheckAllowTXT()

	// The gate matches ^GO-dddd-dddd, so the entry has to start the line.
	var found bool
	for _, line := range strings.Split(allow, "\n") {
		if strings.HasPrefix(line, "GO-2026-6909 ") {
			found = true
			if !strings.Contains(line, "MODULE_SSO") {
				t.Error("the entry does not say what makes it unreachable, which is the " +
					"condition somebody has to re-read the day they enable SAML")
			}
			if !strings.Contains(line, "issues/187") {
				t.Error("the entry names no issue, which the file's own format requires")
			}
		}
	}
	if !found {
		t.Error("the allowlist does not ship the crewjam/saml entry, so a new project's " +
			"security.yml fails on its first push with no version to move to")
	}
}

// moduleFlagRegion trims config.go to the part these tests are about.
func moduleFlagRegion(src string) string {
	i := strings.Index(src, "type ModuleFlags struct")
	if i == -1 {
		return src
	}
	j := strings.Index(src[i:], "\n}\n")
	if j == -1 {
		return src[i:]
	}
	return src[i : i+j]
}
