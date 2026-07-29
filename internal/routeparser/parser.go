package routeparser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Route represents a parsed API route.
type Route struct {
	Method  string
	Path    string
	Handler string
	Group   string // middleware group (public, protected, admin)
}

// Parse reads a routes.go file and extracts all registered routes.
func Parse(routesFile string) ([]Route, error) {
	data, err := os.ReadFile(routesFile)
	if err != nil {
		return nil, fmt.Errorf("opening routes file: %w", err)
	}
	lines := strings.Split(string(data), "\n")

	// String constants are collected up front because a group prefix may be
	// built from one: routes.go declares `const APIVersion = "v1"` and mounts
	// everything under r.Group("/api/" + APIVersion). Without resolving that,
	// the prefix silently evaluates to nothing and every route is reported one
	// path segment short — /users instead of /api/v1/users, which is worse than
	// no output at all because it looks plausible.
	consts := parseStringConsts(lines)

	var routes []Route
	var currentGroup string // semantic group name (public, protected, admin)

	// Patterns for route registration
	routeRe := regexp.MustCompile(`\b(\w+)\.(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\("([^"]+)",\s*(\w+[\w.]*)\)`)
	// The receiver is captured rather than searched for: the previous version
	// looked for any known variable name appearing in the line, iterating a map
	// in random order, so a line mentioning two known variables could inherit
	// the wrong parent prefix on some runs and not others.
	groupRe := regexp.MustCompile(`(\w+)\s*:?=\s*(?:(\w+)\.)?Group\(([^()]*)\)`)
	useAuthRe := regexp.MustCompile(`\.Use\(middleware\.Auth`)
	useRoleRe := regexp.MustCompile(`\.Use\(middleware\.RequireRole\("(\w+)"\)`)

	// Map variable names to their group prefixes
	groupPrefixes := map[string]string{
		"r": "", // root router
	}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// Detect group definitions: auth := r.Group("/api/auth")
		if matches := groupRe.FindStringSubmatch(line); len(matches) >= 4 {
			varName := matches[1]
			receiver := matches[2]
			prefix := resolveStringExpr(matches[3], consts)

			groupPrefixes[varName] = groupPrefixes[receiver] + prefix
		}

		// Detect middleware usage for semantic grouping
		if useAuthRe.MatchString(line) {
			currentGroup = "protected"
		}
		if matches := useRoleRe.FindStringSubmatch(line); len(matches) >= 2 {
			currentGroup = strings.ToLower(matches[1])
		}

		// Detect route registrations
		if matches := routeRe.FindStringSubmatch(line); len(matches) >= 5 {
			varName := matches[1]
			method := matches[2]
			path := matches[3]
			handler := matches[4]

			// Resolve full path from group prefix
			prefix := groupPrefixes[varName]
			fullPath := prefix + path

			// Clean up double slashes
			fullPath = strings.ReplaceAll(fullPath, "//", "/")

			group := currentGroup
			if group == "" {
				group = "public"
			}

			routes = append(routes, Route{
				Method:  method,
				Path:    fullPath,
				Handler: handler,
				Group:   group,
			})
		}
	}

	// Sort by path then method
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path == routes[j].Path {
			return routes[i].Method < routes[j].Method
		}
		return routes[i].Path < routes[j].Path
	})

	return routes, nil
}

var (
	singleConstRe = regexp.MustCompile(`^const\s+(\w+)\s*(?:=|\w+\s*=)\s*"([^"]*)"`)
	blockConstRe  = regexp.MustCompile(`^(\w+)\s*(?:=|\w+\s*=)\s*"([^"]*)"`)
)

// parseStringConsts collects package-level string constants, both the
// single-line form and entries inside a const (...) block, so a group prefix
// assembled from one can be resolved.
func parseStringConsts(lines []string) map[string]string {
	consts := map[string]string{}
	inBlock := false

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		if strings.HasPrefix(line, "const (") {
			inBlock = true
			continue
		}
		if inBlock {
			if line == ")" {
				inBlock = false
				continue
			}
			if m := blockConstRe.FindStringSubmatch(line); len(m) == 3 {
				consts[m[1]] = m[2]
			}
			continue
		}
		if m := singleConstRe.FindStringSubmatch(line); len(m) == 3 {
			consts[m[1]] = m[2]
		}
	}
	return consts
}

// resolveStringExpr evaluates a Group() argument such as `"/api/auth"` or
// `"/api/" + APIVersion` into the prefix it produces.
//
// An identifier that cannot be resolved is rendered as {Name} rather than
// dropped. A path that is visibly incomplete tells you the parser hit something
// it does not understand; a path that silently lost a segment reads as correct
// and sends you at the wrong URL.
func resolveStringExpr(expr string, consts map[string]string) string {
	var out strings.Builder
	for _, part := range strings.Split(expr, "+") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if len(part) >= 2 && strings.HasPrefix(part, `"`) && strings.HasSuffix(part, `"`) {
			out.WriteString(part[1 : len(part)-1])
			continue
		}
		if v, ok := consts[part]; ok {
			out.WriteString(v)
			continue
		}
		out.WriteString("{" + part + "}")
	}
	return out.String()
}

// FindRoutesFile locates the routes.go file in a Grit project.
func FindRoutesFile(projectRoot string) (string, error) {
	// Single app: internal/routes/routes.go
	single := filepath.Join(projectRoot, "internal", "routes", "routes.go")
	if _, err := os.Stat(single); err == nil {
		return single, nil
	}

	// Monorepo: apps/api/internal/routes/routes.go
	mono := filepath.Join(projectRoot, "apps", "api", "internal", "routes", "routes.go")
	if _, err := os.Stat(mono); err == nil {
		return mono, nil
	}

	return "", fmt.Errorf("routes.go not found (checked internal/routes/ and apps/api/internal/routes/)")
}

// FormatTable formats routes as a printable table.
func FormatTable(routes []Route) string {
	if len(routes) == 0 {
		return "  No routes found."
	}

	// Calculate column widths
	methodW, pathW, handlerW, groupW := 6, 4, 7, 5
	for _, r := range routes {
		if len(r.Method) > methodW {
			methodW = len(r.Method)
		}
		if len(r.Path) > pathW {
			pathW = len(r.Path)
		}
		if len(r.Handler) > handlerW {
			handlerW = len(r.Handler)
		}
		if len(r.Group) > groupW {
			groupW = len(r.Group)
		}
	}

	var b strings.Builder
	divider := fmt.Sprintf("  %-*s  %-*s  %-*s  %-*s", methodW, strings.Repeat("─", methodW), pathW, strings.Repeat("─", pathW), handlerW, strings.Repeat("─", handlerW), groupW, strings.Repeat("─", groupW))

	fmt.Fprintf(&b, "  %-*s  %-*s  %-*s  %-*s\n", methodW, "METHOD", pathW, "PATH", handlerW, "HANDLER", groupW, "GROUP")
	b.WriteString(divider + "\n")

	for _, r := range routes {
		fmt.Fprintf(&b, "  %-*s  %-*s  %-*s  %-*s\n", methodW, r.Method, pathW, r.Path, handlerW, r.Handler, groupW, r.Group)
	}

	fmt.Fprintf(&b, "\n  %d routes total\n", len(routes))
	return b.String()
}

// GroupStack helpers
func pushGroup(stack []string, prefix string) []string {
	return append(stack, prefix)
}

func popGroup(stack []string) []string {
	if len(stack) == 0 {
		return stack
	}
	return stack[:len(stack)-1]
}

func currentPrefix(stack []string) string {
	return strings.Join(stack, "")
}

// Unused but keep for reference
var _ = pushGroup
var _ = popGroup
var _ = currentPrefix
