package scaffold

import (
	"fmt"

	"github.com/MUKE-coder/grit/v3/internal/accessgen"
	"github.com/fatih/color"
)

// RegenerateAccessTable rewrites internal/access/registry.go from the route
// files.
//
// Called from the upgrade command after everything that writes a route, which
// is not the same place as the end of Upgrade: the money repair, the
// write-checks repair and the plugin update all run after it, and all three can
// add routes. Written inside Upgrade it recorded a table six routes short, and
// nothing noticed until an idempotency harness upgraded a project twice and
// diffed the results (review 2.5).
//
// An upgrade is also the one moment a project that has never had this table
// gets one, so a missing table is written rather than skipped.
func RegenerateAccessTable(root string) {
	res, err := accessgen.Generate(root)
	if err != nil {
		fmt.Printf("  ! the access table could not be regenerated: %v\n", err)
		fmt.Printf("    the routes screen will show access as unknown until it is\n")
		return
	}
	if res.Changed {
		color.New(color.FgHiGreen).Printf(
			"  \u2713 %s covers %d routes, so the routes screen can show what each one needs\n",
			res.Path, res.Routes)
	}
}
