package doctor

import (
	"fmt"

	"github.com/MUKE-coder/grit/v3/internal/accessgen"
)

// checkAccessTableStale reports an access table that no longer matches the
// route files it was built from.
//
// The table says what each route demands of a caller, and the routes screen
// reads it. The running process cannot work that out for itself: gin's
// RouteInfo carries a route's last handler and nothing about the middleware in
// front of it, so a screen with no table has no honest answer and shows the
// access column as unknown.
//
// `grit new`, `grit generate resource`, `grit remove` and `grit upgrade` all
// rewrite it. What they cannot catch is a route added by hand, which is the
// normal way a project grows past its generated resources. A stale table is
// then the one failure mode that matters: the new route is missing from it
// entirely, so the screen says unknown, which is fine, or an existing route's
// guard was changed and the table still reports the old one, which is not.
//
// A warning rather than an error. Nothing is broken, the build passes, and the
// cost is one column on one screen. But it is exactly the kind of quietly-wrong
// this project treats as worth a line of output.
func checkAccessTableStale(p *project) []Finding {
	stale, err := accessgen.Stale(p.root)
	if err != nil {
		return []Finding{{
			Level:    "warning",
			Resource: "internal/access/registry.go",
			Message:  fmt.Sprintf("the route files could not be read to check the access table: %v", err),
			Fix:      "grit routes   # shows the same parse, with the error in context",
		}}
	}
	if !stale {
		return nil
	}

	return []Finding{{
		Level:    "warning",
		Resource: "internal/access/registry.go",
		Message: "the access table no longer matches the route files, so the routes screen " +
			"and the MCP route tool may report the wrong permission for a route whose guard changed",
		Fix: "grit upgrade   # regenerates it, or `grit generate resource` for a new one",
	}}
}
