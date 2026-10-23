package scaffold

import "strings"

// scoutRoutesGo emits internal/routes/scout.go.
//
// A file of its own rather than more lines in routes.go, following the
// per-resource route files: routes.go wires the application together and is
// hard enough to read without a scanner's configuration in the middle of it.
func scoutRoutesGo(modulePath string) string {
	return strings.ReplaceAll(scoutRoutesTemplate, "{{MODULE}}", modulePath)
}

const scoutRoutesTemplate = `package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"log"

	"github.com/MUKE-coder/scout"
	// Every check registers itself from here. Importing a single collector
	// instead gives a scan that quietly covers less than it says it does.
	_ "github.com/MUKE-coder/scout/collectors/all"
	"github.com/MUKE-coder/scout/core"
	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/access"
	"{{MODULE}}/internal/config"
)

// mountScout attaches the security scanner at /scout.
//
// Off unless SCOUT_ENABLED=true, and it refuses to mount without a password:
// the dashboard shows findings, which is a map of where the application is
// weakest, so an unauthenticated one is worse than no scanner at all.
func mountScout(r *gin.Engine, cfg *config.Config) {
	if !cfg.ScoutEnabled {
		return
	}
	if cfg.ScoutPassword == "" {
		log.Println("Scout is enabled but SCOUT_PASSWORD is empty, so it was not mounted")
		return
	}

	env := core.EnvDev
	switch cfg.AppEnv {
	case "production":
		env = core.EnvProd
	case "staging":
		env = core.EnvStaging
	}

	scout.Mount(r, scout.Config{
		AppName:     cfg.AppName,
		Environment: env,
		Username:    cfg.ScoutUsername,
		Password:    cfg.ScoutPassword,
		DBPath:      cfg.ScoutDBPath,
		JWTSecret:   scoutDashboardKey(cfg.JWTSecret),
	},
		scout.WithTarget(cfg.AppURL),
		scout.WithRoutes(scoutRoutes()...),
	)

	log.Printf("Scout mounted at /scout (%d routes, safe mode)", access.Count())
}

// scoutDashboardKey signs Scout's dashboard sessions.
//
// Derived from the application's JWT secret rather than being it. Sharing one
// signing key between two systems means a token minted by either is valid to
// the other, and the two have different ideas about who may hold one. Derived
// rather than random because a random key changes on restart, which signs
// everybody out of the dashboard every deploy, which is how a warning becomes
// something people route around.
func scoutDashboardKey(appSecret string) string {
	sum := sha256.Sum256([]byte(appSecret + "|scout-dashboard-v1"))
	return hex.EncodeToString(sum[:])
}

// scoutRoutes turns the generated access table into the route table Scout
// scans.
//
// This is the part worth having. Scout's strongest check reads a route's
// middleware chain and reports a sensitive route with no auth gate, and it
// cannot learn that chain from gin: a RouteInfo carries the last handler and
// nothing in front of it. The access table is built by parsing the route
// files, so it knows, and ` + "`grit doctor`" + ` fails when it has gone stale.
//
// The chain is reconstructed from what each route demands rather than copied
// from the source, because what Scout does with it is look for an auth gate.
// A route that requires nothing gets an empty chain, which is the truth about
// it: /api/v1/auth/login cannot require a signed-in user.
//
// Known rough edge: Scout's auth.missing check treats any path containing
// "/api" as sensitive, so in active mode it reports every public route here
// as missing auth. It is a safe-mode no-op, so a default scan is unaffected.
// See https://github.com/MUKE-coder/scout/issues/2.
func scoutRoutes() []core.Route {
	table := access.All()
	out := make([]core.Route, 0, len(table))
	for key, req := range table {
		method, path, ok := splitRouteKey(key)
		if !ok {
			continue
		}
		out = append(out, core.Route{
			Method:     method,
			Path:       path,
			Middleware: scoutMiddleware(req),
		})
	}
	return out
}

// scoutMiddleware names the gates in front of a route, in the order a request
// meets them.
func scoutMiddleware(req access.Requirement) []string {
	if req.Public() {
		return nil
	}
	mw := []string{"AuthRequired"}
	for _, role := range req.Roles {
		mw = append(mw, "RequireRole("+role+")")
	}
	for _, perm := range req.Permissions {
		mw = append(mw, "RequirePermission("+perm+")")
	}
	return mw
}

// splitRouteKey splits "GET /api/v1/widgets" into its two halves.
func splitRouteKey(key string) (method, path string, ok bool) {
	for i := 0; i < len(key); i++ {
		if key[i] == ' ' {
			return key[:i], key[i+1:], i > 0 && i+1 < len(key)
		}
	}
	return "", "", false
}
`
