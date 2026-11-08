package scaffold

// Giving an existing project the SSO module flag.
//
// grit#187: SSO mounts four public endpoints, and POST /auth/saml/:slug/acs
// parses a signed assertion through github.com/crewjam/saml, which carries a
// called advisory with no fixed version. Most projects never configure an IdP,
// so MODULE_SSO now decides whether any of it is mounted, defaulting off for a
// new project.
//
// Two things make this an upgrade problem rather than a template change.
//
// routes.go is framework-owned and rewritten, and it now reads
// cfg.Modules.SSO. config.go is the developer's and is repaired rather than
// rewritten, because a project may have added configuration of its own. So
// without this, every existing project's upgrade produces a routes.go that
// names a field its config.go does not have, and the project stops compiling.
// That is the shape of five bugs already.
//
// And the default. A new project gets SSO off, which is the point. An existing
// project that has SAML configured would lose it silently, so the upgrade
// writes MODULE_SSO=true into its .env: the flag arrives, nothing changes
// today, and turning it off is then a decision somebody makes by reading the
// line rather than by upgrading.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

const (
	// The field, after the last of the flags that were there before it.
	ssoFieldAnchor = "\tTwoFactor bool // TOTP / 2FA\n"
	// Byte for byte the template's text, so a repaired config.go is the
	// generated one. TestRepairedSSOConfigIsTheGeneratedConfig holds that.
	ssoField = "\t// Enterprise SSO: OIDC and SAML. The one module that is OFF by default.\n" +
		"\t//\n" +
		"\t// Every other battery costs nothing to leave on. This one mounts four\n" +
		"\t// public endpoints, and POST /auth/saml/:slug/acs parses a signed\n" +
		"\t// assertion through github.com/crewjam/saml, which carries a called\n" +
		"\t// advisory with no fixed version: GO-2026-6909, a nil-pointer panic and a\n" +
		"\t// namespace-blind Subject.NameID. Most projects never configure an IdP, so\n" +
		"\t// the route was surface nobody asked for (grit#187).\n" +
		"\t//\n" +
		"\t// Turning it on is a deliberate act, and the docs say what it costs.\n" +
		"\tSSO bool\n"

	// The switch arm, so Enabled(\"sso\") answers and the admin nav can hide it.
	ssoEnabledAnchor = "\tcase \"twofactor\":\n\t\treturn m.TwoFactor\n"
	ssoEnabledCase   = "\tcase \"sso\":\n\t\treturn m.SSO\n"

	// The map the /system/modules endpoint returns.
	ssoMapAnchor = "\t\t\"flags\":     m.Flags,\n"
	ssoMapEntry  = "\t\t\"sso\":       m.SSO,\n"

	// The read. Defaults false here, like the template: an existing project
	// gets its true from .env, written below.
	ssoReadAnchor = "\t\t\tTwoFactor: getEnv(\"MODULE_TWOFACTOR\", \"true\") == \"true\",\n"
	ssoRead       = "\t\t\t// The one that defaults off: see ModuleFlags.SSO.\n" +
		"\t\t\tSSO: getEnv(\"MODULE_SSO\", \"false\") == \"true\",\n"

	ssoConfigWarn = "config.go is not the file Grit wrote, so it has no SSO module flag while " +
		"routes.go now reads one: add `SSO bool` to ModuleFlags, a \"sso\" arm to Enabled and Map, " +
		"and `SSO: getEnv(\"MODULE_SSO\", \"false\") == \"true\"` where the other modules are read, " +
		"or the project will not compile"
)

// addSSOModuleFlag repairs config.go and preserves the project's behaviour.
func addSSOModuleFlag(root string, opts Options) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	apiRoot := opts.APIRoot(root)
	cfg := filepath.Join(apiRoot, "internal", "config", "config.go")
	if !fileExists(cfg) {
		return nil
	}
	// Only a project that has the module flags at all. One from before them is
	// not one this applies to.
	if !fileContains(cfg, "ModuleFlags") {
		return nil
	}
	if err := repairSourceFile(root, m, cfg, repairSSOModuleConfigSource); err != nil {
		return err
	}

	// routes.go only after config.go has the field, and only if it does: a
	// routes.go naming cfg.Modules.SSO against a config that has no such field
	// is a project that does not compile, and an upgrade must not be the reason
	// for that.
	if !fileContains(cfg, "MODULE_SSO") {
		return nil
	}
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	if fileExists(routes) {
		if err := repairSourceFile(root, m, routes, repairSSORoutesSource); err != nil {
			return err
		}
	}

	// And only once the field is really there. An .env line for a flag the
	// config cannot read would say the feature is on while nothing reads it.
	if !fileContains(cfg, "MODULE_SSO") {
		return nil
	}
	return preserveSSOInEnv(root)
}

func repairSSOModuleConfigSource(src string) (string, []string, []string) {
	if strings.Contains(src, "MODULE_SSO") {
		return src, nil, nil
	}
	for _, anchor := range []string{ssoFieldAnchor, ssoEnabledAnchor, ssoMapAnchor, ssoReadAnchor} {
		if strings.Count(src, anchor) != 1 {
			return src, nil, []string{ssoConfigWarn}
		}
	}
	out := strings.Replace(src, ssoFieldAnchor, ssoFieldAnchor+ssoField, 1)
	out = strings.Replace(out, ssoEnabledAnchor, ssoEnabledAnchor+ssoEnabledCase, 1)
	out = strings.Replace(out, ssoMapAnchor, ssoMapAnchor+ssoMapEntry, 1)
	out = strings.Replace(out, ssoReadAnchor, ssoReadAnchor+ssoRead, 1)
	return out, []string{"SSO is a module now, off by default in a new project: MODULE_SSO decides " +
		"whether the OIDC and SAML endpoints are mounted at all"}, nil
}

// preserveSSOInEnv writes MODULE_SSO=true, so an upgrade is not the reason a
// project's SAML stops working.
//
// The compiled default is false, which is right for a new project and wrong as
// a silent change to one that exists. Writing the value makes the upgrade a
// no-op at runtime and leaves a line somebody can read and change.
func preserveSSOInEnv(root string) error {
	path := filepath.Join(root, ".env")
	raw, err := os.ReadFile(path)
	if err != nil {
		// No .env is a project that configures the environment elsewhere, where
		// the compiled default applies and Grit has nothing to write.
		return nil
	}
	src := string(raw)
	if strings.Contains(src, "MODULE_SSO") {
		return nil
	}

	nl := "\n"
	if strings.Contains(src, "\r\n") {
		nl = "\r\n"
	}
	if !strings.HasSuffix(src, nl) {
		src += nl
	}
	src += nl +
		"# Enterprise SSO (OIDC + SAML). Written by grit upgrade as true because it is" + nl +
		"# what this project did before the flag existed, so upgrading changed nothing." + nl +
		"#" + nl +
		"# A new project gets false. It mounts four public endpoints, and" + nl +
		"# POST /auth/saml/:slug/acs parses a signed assertion through crewjam/saml," + nl +
		"# which has a called advisory with no fixed version (GO-2026-6909). If you do" + nl +
		"# not configure an IdP, set this to false and the routes are never mounted." + nl +
		"MODULE_SSO=true" + nl

	if err := writeRepaired(path, src); err != nil {
		return err
	}
	fmt.Println("  ✓ .env: MODULE_SSO=true, which is what this project did before the flag " +
		"existed. Set it to false if you do not use SAML or OIDC: it mounts four public " +
		"endpoints, one of which carries an advisory with no fix (grit#187)")
	return nil
}

// The two route groups, as Grit wrote them, and wrapped in the flag.
//
// routes.go is the developer's file: it is repaired, not rewritten, which is
// why an existing project keeps its unguarded groups until this runs. Anchored
// on the exact text, so a project that has edited either group is warned rather
// than rewritten.
const (
	ssoRoutesOld = `	sso := auth.Group("/sso")
	{
		sso.POST("/discover", ssoHandler.Discover)
		sso.GET("/:slug", ssoHandler.Begin)
		sso.GET("/:slug/callback", ssoHandler.Callback)
	}
`
	ssoRoutesNew = `	//
	// Only when MODULE_SSO=true. Nothing is mounted otherwise, so a project
	// with no IdP has no SSO surface at all rather than endpoints that answer
	// for a connection that does not exist.
	if ssoEnabled {
		sso := auth.Group("/sso")
		{
			sso.POST("/discover", ssoHandler.Discover)
			sso.GET("/:slug", ssoHandler.Begin)
			sso.GET("/:slug/callback", ssoHandler.Callback)
		}
	}
`
	samlRoutesOld = `	samlGroup := auth.Group("/saml")
	{
		samlGroup.GET("/:slug/metadata", ssoHandler.SAMLMetadata)
		samlGroup.GET("/:slug", ssoHandler.SAMLBegin)
		samlGroup.POST("/:slug/acs", ssoHandler.SAMLACS)
	}
`
	samlRoutesNew = `	//
	// /acs is the one that parses an assertion through crewjam/saml, which is
	// where GO-2026-6909 lives with no fixed version, so MODULE_SSO=false is
	// what keeps it off the router (grit#187).
	if ssoEnabled {
		samlGroup := auth.Group("/saml")
		{
			samlGroup.GET("/:slug/metadata", ssoHandler.SAMLMetadata)
			samlGroup.GET("/:slug", ssoHandler.SAMLBegin)
			samlGroup.POST("/:slug/acs", ssoHandler.SAMLACS)
		}
	}
`
	ssoRoutesWarn = "routes.go does not mount the SSO and SAML groups the way Grit wrote them, " +
		"so MODULE_SSO cannot gate them here: wrap both auth.Group(\"/sso\") and " +
		"auth.Group(\"/saml\") in `if cfg.Modules.SSO`, or POST /auth/saml/:slug/acs stays " +
		"mounted on every project whether or not it uses SAML (grit#187)"
)

// repairSSORoutesSource wraps the groups, and passes the flag to the function
// that mounts them.
func repairSSORoutesSource(src string) (string, []string, []string) {
	if strings.Contains(src, "ssoEnabled") {
		return src, nil, nil
	}
	// Nothing to do in a project with no SSO routes at all.
	if !strings.Contains(src, ssoRoutesOld) && !strings.Contains(src, samlRoutesOld) {
		return src, nil, nil
	}
	if strings.Count(src, ssoRoutesOld) != 1 || strings.Count(src, samlRoutesOld) != 1 {
		return src, nil, []string{ssoRoutesWarn}
	}

	// mountAuthRoutes has to see the flag. Both the signature and the call,
	// or the file does not compile: the pair is the whole repair.
	const sigOld = "func mountAuthRoutes(v1 *gin.RouterGroup, authHandler *handlers.AuthHandler, ssoHandler *handlers.SSOHandler, totpHandler *handlers.TOTPHandler, passkeyHandler *handlers.PasskeyHandler) {"
	const callOld = "mountAuthRoutes(v1, authHandler, ssoHandler, totpHandler, passkeyHandler)"
	if strings.Count(src, sigOld) != 1 || strings.Count(src, callOld) != 1 {
		return src, nil, []string{ssoRoutesWarn}
	}

	out := strings.Replace(src, ssoRoutesOld, ssoRoutesNew, 1)
	out = strings.Replace(out, samlRoutesOld, samlRoutesNew, 1)
	// A bool rather than the whole config, which is what the template does:
	// the function needs one answer, and taking *config.Config would need an
	// import this repair cannot be sure is already there.
	out = strings.Replace(out, sigOld,
		"func mountAuthRoutes(v1 *gin.RouterGroup, authHandler *handlers.AuthHandler, ssoHandler *handlers.SSOHandler, totpHandler *handlers.TOTPHandler, passkeyHandler *handlers.PasskeyHandler, ssoEnabled bool) {", 1)
	out = strings.Replace(out, callOld,
		"mountAuthRoutes(v1, authHandler, ssoHandler, totpHandler, passkeyHandler, cfg.Modules.SSO)", 1)

	return out, []string{"the OIDC and SAML endpoints are mounted only when MODULE_SSO=true, " +
		"so a project with no IdP has no SSO surface"}, nil
}
