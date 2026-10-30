package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Review 2.8: the sign-in routes were rate limited only when Sentinel was on.
//
// Every number lived inside mountSentinel, behind a check on SENTINEL_ENABLED.
// That is the default off in development and one environment variable in
// production, and with it off /auth/login, /auth/forgot-password,
// /auth/totp/verify and the passkey routes had no limiter at all. The
// per-account lockout still barred a single account after LOGIN_MAX_ATTEMPTS
// failures, so guessing one password was bounded; nothing bounded credential
// stuffing across many accounts from one address, or a forgot-password loop
// pointed at one mailbox.
//
// Three deliveries, because the fix is a new file, a table and a mount:
//
//   - middleware/authlimit.go arrives whole, framework code
//   - authRouteLimits gains the routes it was missing, and builtInAuthLimits
//     reads the same table for the limiter that needs no Sentinel
//   - routes.go mounts it on the versioned group
//
// The table and the limiter come from routesAuthLimitsFunc, the one constant
// the new-project template also uses, so the two cannot be written differently.

// repairAuthRateLimit brings an existing project up to a limiter that holds
// whether or not Sentinel is mounted.
func repairAuthRateLimit(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	if !fileExists(routes) {
		return nil
	}

	// The middleware is framework code and arrives whole, like limits.go.
	for path, content := range map[string]string{
		filepath.Join(apiRoot, "internal", "middleware", "authlimit.go"):      middlewareAuthLimitGo(),
		filepath.Join(apiRoot, "internal", "middleware", "authlimit_test.go"): middlewareAuthLimitTestGo(),
	} {
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", opts.Module())); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, routes, repairAuthRateLimitSource)
}

// repairAuthRateLimitSource replaces the limits table and mounts the limiter.
func repairAuthRateLimitSource(src string) (string, []string, []string) {
	if strings.Contains(src, "func builtInAuthLimits(") {
		return src, nil, nil
	}
	// A project that never had the table at all is one repairRoutesProduction
	// has not reached yet. It inserts the whole constant, this version of it
	// included, so there is nothing to do here.
	if !strings.Contains(src, "func authRouteLimits(dev bool) map[string]sentinel.Limit {") {
		return src, nil, nil
	}

	const warn = "routes.go is not the file Grit wrote: mount middleware.AuthRateLimit on the /api/<version> group with the limits from authRouteLimits, or turning Sentinel off leaves the sign-in routes with no rate limit"

	// The whole function, replaced rather than edited: the new one adds seven
	// routes and changes what development does, and patching a map literal in
	// place is how a repair ends up with half of each.
	old, ok := methodBody(src, "func authRouteLimits(dev bool) map[string]sentinel.Limit {")
	if !ok {
		return src, nil, []string{warn}
	}
	// routesAuthLimitsFunc carries its own doc comment, and the file already
	// has one above the function being replaced. Taking the text from the end
	// of that comment keeps exactly one.
	replacement := routesAuthLimitsFunc
	if i := strings.Index(replacement, "func authRouteLimits(dev bool)"); i > 0 {
		replacement = replacement[i:]
	}
	out := strings.Replace(src, old, replacement, 1)

	// And the mount. Anchored on the group's creation, which is the one line
	// in this file that establishes the prefix the table is keyed by.
	const anchor = "\tv1 := r.Group(\"/api/\" + APIVersion)\n"
	if strings.Count(out, anchor) != 1 {
		return src, nil, []string{warn}
	}
	out = strings.Replace(out, anchor, anchor+`
	// The limit on the sign-in routes, whether or not Sentinel is mounted.
	//
	// Review 2.8: every one of these numbers was read inside mountSentinel, so
	// SENTINEL_ENABLED=false left /auth/login, /auth/forgot-password and the
	// code routes with no limiter at all. This reads the same table and is the
	// floor under it. Mounted here rather than on the auth group because the
	// passkey and TOTP routes are mounted elsewhere and the table is keyed by
	// the whole path: a path it does not list costs one map lookup and is let
	// through.
	v1.Use(middleware.AuthRateLimit(builtInAuthLimits(cfg.AppEnv == "development")))
`, 1)

	return out, []string{"the sign-in routes are rate limited whether or not Sentinel is mounted, and the table covers magic links, backup codes, passkey logins and the SSO callbacks"}, nil
}
