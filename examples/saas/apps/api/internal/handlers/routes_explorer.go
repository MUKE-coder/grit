package handlers

import (
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"

	"saas/apps/api/internal/access"
	"saas/apps/api/internal/llms"
)

// RouteExplorerHandler serves the route table to the admin's Routes screen.
//
// # Where each column comes from
//
// Method, path and the final handler come from the router, which is the only
// thing that knows what was actually registered.
//
// Access comes from the generated access table, not from the router. Gin's
// RouteInfo carries a route's last handler and nothing about the middleware in
// front of it, so at runtime this process genuinely cannot tell that
// /api/v1/widgets sits behind RequireRole("ADMIN", "perm:widgets.view"). The
// table is built by `grit` from routes.go and the per-resource route files,
// where the truth is written, and `grit doctor` reports when it has gone stale.
//
// # The three states, and why there are three
//
// A route is either covered by the table, or it is not, and those are
// different answers. The reference at /docs, the profiler and the database
// browser mount their own routes and appear in the router's table without
// appearing in routes.go, so nothing here can speak for them: they come back
// with access_known false, and the screen renders them as unknown.
//
// It would be easier to call an unknown route public and show two states
// instead of three. That is the one thing this must not do. Guessed from the
// path, /api/v1/auth/me and all six passkey endpoints read as public, because
// they sit under /auth/, and every one of them needs a signed-in user.
// Somebody would read that column as a security statement.
type RouteExplorerHandler struct {
	// Engine is the router itself. The one handler in the project that holds it,
	// because the table is the thing being served.
	Engine *gin.Engine
}

// routeRow is one route, as the screen renders it.
type routeRow struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	// Handler is the function the route ends at, as the router names it. It is
	// the honest answer to "what runs" and it is not documentation.
	Handler string `json:"handler,omitempty"`
	// Group is what the route is about, read from the path. It is a heading, not
	// a permission.
	Group string `json:"group"`

	// AccessKnown is false when the generated table does not cover this route,
	// which is every route a mounted dashboard brought with it. The screen must
	// render that as unknown rather than as public.
	AccessKnown bool `json:"access_known"`
	// Access is what the route demands of a caller, when AccessKnown is true.
	Access access.Requirement `json:"access,omitzero"`
}

// List returns every registered route.
//
// Unpaginated on purpose. It is a few hundred rows of three short strings, the
// screen filters in the browser so typing is instant, and a page of routes that
// needed a second request to find the next one would be slower to use than the
// terminal command it exists to replace.
func (h *RouteExplorerHandler) List(c *gin.Context) {
	if h.Engine == nil {
		c.JSON(http.StatusOK, gin.H{"data": []routeRow{}, "count": 0})
		return
	}

	method := strings.ToUpper(strings.TrimSpace(c.Query("method")))
	contains := strings.ToLower(strings.TrimSpace(c.Query("contains")))

	infos := h.Engine.Routes()
	out := make([]routeRow, 0, len(infos))
	for _, info := range infos {
		if method != "" && info.Method != method {
			continue
		}
		if contains != "" && !strings.Contains(strings.ToLower(info.Path), contains) {
			continue
		}
		req, known := access.Of(info.Method, info.Path)
		out = append(out, routeRow{
			Method:      info.Method,
			Path:        info.Path,
			Handler:     llms.ShortHandlerName(info.Handler),
			Group:       llms.GroupOf(info.Path),
			AccessKnown: known,
			Access:      req,
		})
	}

	// Alphabetical within the API, and the mounted dashboards last. They are
	// the biggest group by a distance, 130 of the routes in a fresh project, and
	// they are somebody else's: the reference, the profiler, the database
	// browser. Sorting them first put the least useful thing at the top of the
	// page.
	sort.Slice(out, func(i, j int) bool {
		if out[i].Group != out[j].Group {
			iMounted, jMounted := out[i].Group == llms.GroupMounted, out[j].Group == llms.GroupMounted
			if iMounted != jMounted {
				return jMounted
			}
			return out[i].Group < out[j].Group
		}
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		return out[i].Method < out[j].Method
	})

	c.JSON(http.StatusOK, gin.H{
		"data":  out,
		"count": len(out),
		"covered": gin.H{
			"routes": access.Count(),
			"note": "Access comes from a table grit builds out of routes.go, because the " +
				"router does not carry the middleware in front of a route. Rows with " +
				"access_known false are routes a mounted dashboard registered, which this " +
				"project's route files never saw: unknown, not public.",
		},
	})
}
