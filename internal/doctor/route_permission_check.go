package doctor

import (
	"fmt"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/routeparser"
)

// A route you wrote that stands behind authentication and demands no
// permission (grit#179, proposal 1).
//
// `grit routes --check` is the one for CI, with an exit code and the full
// list. This is the same question asked by the command people already run, so
// the answer arrives without having to know the flag exists.
//
// Why it is an error rather than a warning: forgetting a guard makes a route
// MORE permissive. The mistake nobody notices is the one that costs something,
// and with open registration "any signed-in account" is anybody who filled in
// a form.
func checkRoutePermissions(p *project) []Finding {
	res, err := routeparser.Check(p.root, false)
	if err != nil {
		// No routes file, or a project this cannot read. Not a finding: the
		// other checks will say more useful things about such a project.
		return nil
	}
	if res.OK() || !res.ManifestKnown {
		return nil
	}

	// Grouped by file, because that is how somebody fixes them: one file,
	// one sitting. A flat list of twenty routes across four files reads as
	// twenty problems rather than four.
	byFile := map[string][]string{}
	var order []string
	for _, f := range res.Undeclared {
		if _, seen := byFile[f.File]; !seen {
			order = append(order, f.File)
		}
		byFile[f.File] = append(byFile[f.File], f.Route.Method+" "+f.Route.Path)
	}

	var findings []Finding
	for _, file := range order {
		routes := byFile[file]
		shown := routes
		if len(shown) > 4 {
			shown = append(append([]string{}, routes[:4]...),
				fmt.Sprintf("and %d more", len(routes)-4))
		}
		findings = append(findings, Finding{
			Level:    "error",
			Resource: file,
			Message: fmt.Sprintf("%s behind authentication and demand no permission, so any signed-in account can reach them: %s",
				plural(len(routes), "1 route stands", fmt.Sprintf("%d routes stand", len(routes))),
				strings.Join(shown, ", ")),
			Fix: `Declare what each one needs: perm("payroll.edit") on the route, which is ` +
				`RequireRole("ADMIN", "perm:payroll.edit") written short. For a route that ` +
				`genuinely needs none, because its handler scopes every row to the caller, say ` +
				`openToAnySignedIn("the reason"). Then ` + "`grit routes --check`" + ` passes and ` +
				`keeps passing.`,
		})
	}
	return findings
}
