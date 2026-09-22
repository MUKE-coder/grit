package scaffold

import (
	"fmt"
	"go/format"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// grit add variants mounted the variant routes on the protected group, which
// every signed-in user reaches, while the resource's own routes ask for ADMIN
// or a permission. So any customer could PATCH a variant's price override to a
// cent and check out at that price, or create and delete the shop's options.
// They go on the staff group now, with the resource's own view and edit
// permissions. Found building the storefront blueprint, whose checkout prices
// every line from the variant.

var (
	protectedVariantRoute = regexp.MustCompile(`(?m)^[ \t]*protected\.(GET|POST|PUT|PATCH|DELETE)\("([^"]+)", (\w+)VariantHandler\.(\w+)\)$`)
	variantListRoute      = regexp.MustCompile(`^/(.+)/:id/variants$`)
)

// VariantStaffRoute is one variant route on the staff group, asking for the
// resource's view permission to read and its edit permission to change.
func VariantStaffRoute(method, route, plural, handler, action string) string {
	perm := "edit"
	if method == "GET" {
		perm = "view"
	}
	return fmt.Sprintf("\tstaff.%s(%q, middleware.RequireRole(\"ADMIN\", \"perm:%s.%s\"), %sVariantHandler.%s)",
		method, route, plural, perm, handler, action)
}

// VariantRoutesOnStaff moves every variant route in routes.go from the
// protected group to the staff group. It reports whether anything moved, and
// warns when it cannot: a routes.go with no staff group declared before them.
func VariantRoutesOnStaff(src string) (string, []string, []string) {
	matches := protectedVariantRoute.FindAllStringSubmatchIndex(src, -1)
	if len(matches) == 0 {
		return src, nil, nil
	}
	staffAt := strings.Index(src, "staff := v1.Group(")
	if staffAt < 0 || staffAt > matches[0][0] {
		return src, nil, []string{"the variant routes are on the protected group, which every signed-in user reaches: move them behind middleware.RequireRole(\"ADMIN\", \"perm:<resource>.edit\"), or a customer can change a variant's price"}
	}

	// A handler's permissions are its resource's, read off its list route.
	plurals := map[string]string{}
	for _, m := range protectedVariantRoute.FindAllStringSubmatch(src, -1) {
		if list := variantListRoute.FindStringSubmatch(m[2]); list != nil {
			plurals[m[3]] = list[1]
		}
	}
	var warnings []string
	out := protectedVariantRoute.ReplaceAllStringFunc(src, func(line string) string {
		m := protectedVariantRoute.FindStringSubmatch(line)
		plural, ok := plurals[m[3]]
		if !ok {
			warnings = append(warnings, "could not tell which resource "+strings.TrimSpace(line)+" belongs to: put it behind the resource's permissions by hand")
			return line
		}
		return VariantStaffRoute(m[1], m[2], plural, m[3], m[4])
	})
	if formatted, err := format.Source([]byte(out)); err == nil {
		out = string(formatted)
	}
	return out, []string{"the variant routes ask for the resource's permissions, instead of any signed-in user"}, warnings
}

// repairVariantRoutes applies it to an existing project's routes.go.
func repairVariantRoutes(root string) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, apiRoot := range []string{filepath.Join(root, "apps", "api"), filepath.Join(root, "api"), root} {
		if path := filepath.Join(apiRoot, "internal", "routes", "routes.go"); fileExists(path) {
			return repairTextFile(root, m, path, VariantRoutesOnStaff)
		}
	}
	return nil
}
