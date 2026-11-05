package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairRealtimeSecurity brings an existing project's realtime socket up to A1.
//
// The upgrader accepted every origin while the handshake authenticates with the
// grit_access cookie, so a page on another origin could open a socket as the
// signed-in user. MODULE_REALTIME=false left /api/ws mounted and the backplane
// subscribed, and nothing capped how many sockets an account could hold.
//
// guard.go is new framework code and arrives whole. The handler and routes.go
// are the developer's, and are edited only where they still read as Grit wrote
// them. The handler is repaired only once routes.go hands it the CORS origins:
// its origin check without that list would refuse the admin panel's own socket.
func repairRealtimeSecurity(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	realtimeDir := filepath.Join(apiRoot, "internal", "realtime")
	hub := filepath.Join(realtimeDir, "hub.go")
	handler := filepath.Join(apiRoot, "internal", "handlers", "realtime.go")
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	if !fileExists(hub) || !fileExists(handler) {
		return nil
	}
	if !fileContains(hub, "clients map[string]map[*Client]struct{}") {
		fmt.Println("  ⚠ internal/realtime/hub.go is not the hub Grit wrote, so the WebSocket still accepts a cookie handshake from any origin.\n" +
			"    Compare it with grit upgrade --diff, check the Origin before accepting the grit_access cookie, and run grit upgrade again.")
		return nil
	}

	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	module := opts.Module()
	// First delivery only. Refreshing these two is writeFrameworkOwnedFiles's
	// job, because handlers/realtime_sse.go calls into guard.go and travels in
	// that set: a guard.go left behind is a project that does not compile.
	for _, f := range []struct{ path, content string }{
		{filepath.Join(realtimeDir, "guard.go"), apiRealtimeGuardGo()},
		{filepath.Join(realtimeDir, "guard_test.go"), apiRealtimeGuardTestGo()},
	} {
		if fileExists(f.path) {
			continue
		}
		if err := writeFile(f.path, strings.ReplaceAll(f.content, "{{MODULE}}", module)); err != nil {
			return fmt.Errorf("writing %s: %w", f.path, err)
		}
		fmt.Printf("  ✓ %s: added\n", shownPath(root, f.path))
	}

	if fileExists(routes) {
		modules := fileContains(filepath.Join(apiRoot, "internal", "config", "config.go"), `"MODULE_REALTIME"`)
		if err := repairSourceFile(root, m, routes, func(src string) (string, []string, []string) {
			return repairRealtimeRoutesSource(src, modules)
		}); err != nil {
			return err
		}
	}

	// #92 item 8: /api/events, the SSE fallback. Run after the route repair
	// above, so a project being brought up to the gated block gets the stream
	// inside it rather than beside it.
	if fileExists(routes) {
		if err := repairSourceFile(root, m, routes, repairRealtimeSSERoute); err != nil {
			return err
		}
	}
	routesSrc, err := os.ReadFile(routes)
	if err != nil || !realtimeOriginsWired(string(routesSrc)) {
		return nil
	}
	if err := repairSourceFile(root, m, handler, repairRealtimeHandlerSource); err != nil {
		return err
	}

	handlerTest := filepath.Join(apiRoot, "internal", "handlers", "realtime_test.go")
	authTest := filepath.Join(apiRoot, "internal", "handlers", "auth_test.go")
	if fileContains(handler, "h.Hub.Admit(client)") && fileContains(authTest, "func newTestAuthSvc(") && !fileExists(handlerTest) {
		if err := writeFile(handlerTest, strings.ReplaceAll(apiRealtimeHandlerTestGo(), "{{MODULE}}", module)); err != nil {
			return fmt.Errorf("writing %s: %w", handlerTest, err)
		}
		fmt.Printf("  ✓ %s: added\n", shownPath(root, handlerTest))
	}
	return nil
}

func shownPath(root, path string) string {
	if rel, err := filepath.Rel(root, path); err == nil {
		return filepath.ToSlash(rel)
	}
	return path
}

// realtimeOriginsWired reports whether routes.go already hands the realtime
// socket its origin list, in either of the two shapes it has had: the package
// variable up to v3.478.0, and the registration against the hub after it.
//
// A project on the older shape is left on it. Its handler reads the package
// variable, so rewriting only routes.go would leave the socket with no origins
// at all and refuse the project's own admin panel.
func realtimeOriginsWired(src string) bool {
	return strings.Contains(src, "realtime.AllowedOrigins = ") ||
		strings.Contains(src, "SetAllowedOrigins(corsOrigins)")
}

// repairRealtimeRoutesSource hands the handler the CORS origins and gates the
// hub's backplane and /api/ws on MODULE_REALTIME. modules says whether the
// project's config has the flag at all.
func repairRealtimeRoutesSource(src string, modules bool) (string, []string, []string) {
	if !strings.Contains(src, "realtimeHandler.Connect") {
		return src, nil, nil
	}
	out := src
	var fixed, warnings []string

	if !realtimeOriginsWired(out) {
		if strings.Count(out, routesCORSOld) == 1 && strings.Count(out, routesRealtimeHandlerOld) == 1 {
			out = strings.Replace(out, routesCORSOld, routesCORSNew, 1)
			out = strings.Replace(out, routesRealtimeHandlerOld, routesRealtimeHandlerNew, 1)
			fixed = append(fixed, "the WebSocket takes a cookie handshake only from the origins CORS allows")
		} else {
			warnings = append(warnings, "does not build CORS or the realtime handler the way Grit wrote it: call realtimeHub.SetAllowedOrigins with your CORS origin list and run grit upgrade again, or the socket keeps accepting a cookie handshake from any origin")
		}
	}

	if modules && !strings.Contains(out, "if cfg.Modules.Realtime {") {
		if strings.Count(out, routesRealtimeHubOld) == 1 && strings.Count(out, routesRealtimeRouteOld) == 1 {
			out = strings.Replace(out, routesRealtimeHubOld, routesRealtimeHubNew, 1)
			out = strings.Replace(out, routesRealtimeRouteOld, routesRealtimeRouteNew, 1)
			fixed = append(fixed, "MODULE_REALTIME=false unmounts /api/ws and starts no Redis backplane")
		} else {
			warnings = append(warnings, "does not build the realtime hub or mount /api/ws the way Grit wrote it: wrap the /api/ws route and realtime.WithRedis in if cfg.Modules.Realtime, or MODULE_REALTIME=false leaves the socket open")
		}
	}
	return out, fixed, warnings
}

// repairRealtimeHandlerSource checks the Origin before accepting the cookie and
// registers through the capped Admit. All or nothing: the pieces share imports.
//
// The gate is the call, in either of the two shapes it has. A handler repaired
// before v3.479.0 calls realtime.CheckOrigin and one written after it calls
// h.Hub.CheckOrigin; both mean "the Origin is already checked here", and gating
// on the older name alone would make this repair run again on a project that
// has the newer shape, fail its first hunk count and warn about a handler that
// is correct. A v1 handler is left as it is: it enforces the right origins, it
// just reads them from a package variable.
// realtimeHandlerHunks is every edit this repair makes to handlers/realtime.go,
// old text to new.
//
// Declared here rather than inside the function because the test reverses it to
// build the file Grit wrote before the fix, and then checks that repairing that
// file reproduces the template byte for byte. Two lists would each pass their
// own test while disagreeing with each other, which is how a hunk added to the
// repair and not to the test went unnoticed until an upgrade warned about a
// handler it had itself just written.
var realtimeHandlerHunks = [][2]string{
	{realtimeHandlerImportOld, realtimeHandlerImportNew},
	{realtimeUpgraderOld, realtimeUpgraderNew},
	{realtimeHandlerStructOld, realtimeHandlerStructNew},
	{realtimeHandlerCtorOld, realtimeHandlerCtorNew},
	{realtimeUpgradeCallOld, realtimeUpgradeCallNew},
	{realtimeTokenOld, realtimeTokenNew},
	{realtimeRegisterOld, realtimeRegisterNew},
}

// realtimeRoutesHunks is the same for routes.go. The route and hub edits are
// applied conditionally, so this is the list in its complete form.
var realtimeRoutesHunks = [][2]string{
	{routesCORSOld, routesCORSNew},
	{routesRealtimeHubOld, routesRealtimeHubNew},
	{routesRealtimeHandlerOld, routesRealtimeHandlerNew},
	{routesRealtimeRouteOld, routesRealtimeRouteNew},
}

func repairRealtimeHandlerSource(src string) (string, []string, []string) {
	// Either shape of the check reads as "already repaired".
	checked := strings.Contains(src, "realtime.CheckOrigin(c.Request)") ||
		strings.Contains(src, "h.Hub.CheckOrigin(c.Request)")
	if !strings.Contains(src, "RealtimeHandler") || checked {
		return src, nil, nil
	}
	hunks := realtimeHandlerHunks
	for _, h := range hunks {
		if strings.Count(src, h[0]) != 1 {
			return src, nil, []string{"is not the handler Grit wrote: check the Origin with realtime.CheckOrigin before accepting the grit_access cookie, and register with h.Hub.Admit, or any site can open a socket as the signed-in user"}
		}
	}
	out := src
	for _, h := range hunks {
		out = strings.Replace(out, h[0], h[1], 1)
	}
	return out, []string{"the socket refuses a cookie handshake from an unlisted origin, and caps connections per user and per process"}, nil
}
