package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// ─── M38: /api/health answers under the version prefix too ──────────────────
//
// The frontends' axios client rewrites every /api/... path to /api/v1/..., and
// the health probe was the one route mounted outside the version group. So the
// admin's System Health page asked for /api/v1/health, no route matched, and
// the unversioned fallback refuses to rewrite a path that already names a
// version: the page rendered "degraded" with a 404 behind it.
//
// The same handler is registered at both paths. /api/health stays because it is
// configured outside this repo: container health checks, load balancers and the
// desktop client's heartbeat all name it.

const healthRouteOld = "\tr.GET(\"/api/health\", func(c *gin.Context) {\n"

const healthRouteNew = `	// Registered at two paths, not one. The frontends' axios client rewrites
	// /api/... to /api/v1/..., so the admin's System Health page asked for
	// /api/v1/health, no route matched, and the unversioned fallback refused to
	// rewrite a path that already names the version: the page read "degraded"
	// with a 404 behind it. /api/health stays for probes, load balancers and the
	// desktop client's heartbeat, which are configured outside this repo.
	healthCheck := func(c *gin.Context) {
`

// repairHealthVersionedRoute applies the second half of M38 to an existing
// project's routes.go.
func repairHealthVersionedRoute(root string, opts Options) error {
	routes := filepath.Join(opts.APIRoot(root), "internal", "routes", "routes.go")
	if !fileExists(routes) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, routes, repairHealthVersionedRouteSource)
}

func repairHealthVersionedRouteSource(src string) (string, []string, []string) {
	if strings.Contains(src, `r.GET("/api/"+APIVersion+"/health", healthCheck)`) ||
		strings.Contains(src, `r.GET("/api/"+APIVersion+"/health"`) {
		return src, nil, nil
	}
	if strings.Count(src, healthRouteOld) != 1 {
		return src, nil, nil
	}
	at := strings.Index(src, healthRouteOld)
	// The handler literal's closing "\t})\n". goBraceEnd walks from the opening
	// brace of the function literal, so the close is exactly where it ends.
	open := strings.Index(src[at:], "{") + at
	end := goBraceEnd(src, open)
	if end < 0 {
		return src, nil, []string{"routes.go: could not find the end of the /api/health handler, so it still answers only at the unversioned path"}
	}
	// What follows the literal is ")\n" from the r.GET call.
	rest := src[end:]
	if !strings.HasPrefix(rest, ")\n") {
		return src, nil, []string{"routes.go: the /api/health route is not the one Grit wrote, so it still answers only at the unversioned path"}
	}

	out := src[:at] + healthRouteNew + src[at+len(healthRouteOld):end] +
		"\n\tr.GET(\"/api/health\", healthCheck)\n\tr.GET(\"/api/\"+APIVersion+\"/health\", healthCheck)\n" +
		rest[len(")\n"):]
	return out, []string{"/api/health answers under the version prefix too, so the admin's health page stops reading 404 as degraded"}, nil
}
