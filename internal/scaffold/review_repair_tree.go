package scaffold

import (
	"fmt"
	"regexp"
	"strings"
)

// treeRouteRe matches a tree route as the generator wrote it into routes.go
// before v3.241.0: on the protected group, open to any signed-in account.
var treeRouteRe = regexp.MustCompile(`(?m)^([ \t]*)protected\.(GET|PATCH|POST)\("/(\w+)(/tree|/:id/breadcrumbs|/:id/move|/reorder|/rebuild-tree)", (\w+)TreeHandler\.(GetTree|GetBreadcrumbs|Move|Reorder|RebuildPaths)\)[ \t]*$`)

var treeRoutePermission = map[string]string{
	"GetTree":        "view",
	"GetBreadcrumbs": "view",
	"Move":           "edit",
	"Reorder":        "edit",
	"RebuildPaths":   "edit",
}

// repairTreeRoutesSource moves a shared tree's routes in routes.go behind the
// resource's permissions, or onto the admin group where the project has no staff
// group. A move or a rebuild rewrites every row beneath a node. Scoped
// (tenant-owned) trees keep their protected routes.
func repairTreeRoutesSource(src string, scoped map[string]bool) (string, []string, []string) {
	staff := strings.Contains(src, "middleware.RequireStaff()")
	moved := 0
	out := treeRouteRe.ReplaceAllStringFunc(src, func(line string) string {
		p := treeRouteRe.FindStringSubmatch(line)
		indent, method, resource, rest, camel, call := p[1], p[2], p[3], p[4], p[5], p[6]
		if scoped[strings.ToUpper(camel[:1])+camel[1:]] {
			return line
		}
		moved++
		if staff {
			return fmt.Sprintf(`%sstaff.%s("/%s%s", middleware.RequireRole("ADMIN", "perm:%s.%s"), %sTreeHandler.%s)`,
				indent, method, resource, rest, resource, treeRoutePermission[call], camel, call)
		}
		return fmt.Sprintf(`%sadmin.%s("/%s%s", %sTreeHandler.%s)`, indent, method, resource, rest, camel, call)
	})
	if moved == 0 {
		return src, nil, nil
	}
	where := "the admin group"
	if staff {
		where = "a permission"
	}
	return out, []string{fmt.Sprintf("%d tree %s now ask for %s instead of any signed-in account",
		moved, plural(moved, "route", "routes"), where)}, nil
}
