package routeparser

import (
	"bufio"
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
	f, err := os.Open(routesFile)
	if err != nil {
		return nil, fmt.Errorf("opening routes file: %w", err)
	}
	defer f.Close()

	var routes []Route
	var currentGroup string // semantic group name (public, protected, admin)

	// Patterns for route registration
	routeRe := regexp.MustCompile(`\b(\w+)\.(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\("([^"]+)",\s*(\w+[\w.]*)\)`)
	groupRe := regexp.MustCompile(`(\w+)\s*:?=\s*(?:\w+\.)?Group\("([^"]+)"\)`)
	useAuthRe := regexp.MustCompile(`\.Use\(middleware\.Auth`)
	useRoleRe := regexp.MustCompile(`\.Use\(middleware\.RequireRole\("(\w+)"\)`)

	// Map variable names to their group prefixes
	groupPrefixes := map[string]string{
		"r": "", // root router
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// Detect group definitions: auth := r.Group("/api/auth")
		if matches := groupRe.FindStringSubmatch(line); len(matches) >= 3 {
			varName := matches[1]
			prefix := matches[2]

			// Find parent prefix by checking what variable the Group is called on
			parentPrefix := ""
			for pName, pPrefix := range groupPrefixes {
				if strings.Contains(line, pName+".Group") || strings.Contains(line, pName+" .Group") {
					parentPrefix = pPrefix
					break
				}
			}
			groupPrefixes[varName] = parentPrefix + prefix
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

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading routes file: %w", err)
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
