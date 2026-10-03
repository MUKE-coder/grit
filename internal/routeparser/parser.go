// Package routeparser reads a generated project's routes.go and reports every
// route it registers, together with what stands in front of that route.
//
// # Why this is a parser and not a runtime question
//
// Gin's router can be asked for its table, and the answer carries a route's
// method, path and last handler. It carries nothing about the middleware in
// front of it: RouteInfo has no chain, and the chain is where authorization
// lives. So the only place the full truth exists is the source file that built
// the router, and reading it is the only way to get at it.
//
// # What the previous version got wrong, and why it mattered
//
// Three defects, all of the same kind: each one produced a plausible answer
// rather than no answer.
//
//  1. A route registered with middleware in front of it was dropped entirely.
//     The pattern required the handler to follow the path immediately, so
//     `admin.GET("/users", middleware.RequireRole("ADMIN", "perm:users.view"),
//     h.List)` matched nothing. In a fresh project that silently lost 63 of
//     171 routes, and specifically the 55 that name a permission.
//
//  2. Authentication was detected with a pattern for `.Use(middleware.Auth`,
//     and the generated file writes `.Use(middleware.APIKeyOrAuth(db,
//     middleware.Auth(...)))`. So no group was ever recognised as protected
//     and every authenticated route was reported as public, including
//     /api/v1/auth/me. A route table that calls an authenticated endpoint
//     public is worse than one that says nothing: somebody reads it as a
//     security statement.
//
//  3. Group prefixes were resolved as the file was read, so a group declared
//     in one function from a parameter belonging to another resolved against
//     an empty parent and lost a path segment.
//
// All three are fixed by doing less guessing: the group graph is built first
// and resolved afterwards, the middleware chain is accumulated per group
// rather than tracked as one running variable, and route arguments are split
// with a paren-aware scanner instead of a regexp that assumes their shape.
package routeparser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Route represents a parsed API route and what it demands of a caller.
type Route struct {
	Method  string
	Path    string
	Handler string
	Group   string // middleware group (public, protected, staff, admin)
	Access  Access
}

// Access is what stands in front of a route: the middleware its group carries,
// inherited down the group chain, plus anything passed on the route itself.
//
// The fields are deliberately separate rather than one string, because the
// question has two halves that fail differently. Authenticated answers "does
// this need a caller at all", and getting it wrong marks a protected endpoint
// public. Roles and Permissions answer "which caller", and getting that wrong
// is a smaller mistake. Summary joins them for display.
type Access struct {
	// Guards are the authorization-relevant middleware in front of the route,
	// outermost first, named as they appear in the source. Middleware that does
	// not decide access (logging, gzip, request IDs) is not listed: a guard list
	// that includes everything is a list nobody reads.
	Guards []string

	// Roles and Permissions are what the route demands. Either satisfies it:
	// RequireRole("ADMIN", "perm:users.view") admits an ADMIN or anybody granted
	// users.view, which is how the generated code uses it.
	Roles       []string
	Permissions []string

	// Staff records that RequireStaff stands in front of the route. It is a
	// gate rather than a demand: it admits the ADMIN role or anybody holding
	// any permission at all, and the route behind it then names the permission
	// it actually needs. So it is kept out of Roles, because reporting
	// "staff or ADMIN or settings.view" states a requirement that is not one.
	Staff bool

	// Authenticated is true when something in front of the route establishes who
	// is calling: a session, a bearer token or an API key.
	Authenticated bool

	// APIKeyOnly is true when an API key alone is accepted, with no user session.
	APIKeyOnly bool

	// Optional is true when the route identifies a caller if one is present and
	// serves anonymous requests otherwise, which is what Identify does.
	Optional bool
}

// Summary renders the access requirement as one short phrase, for a table
// column or a UI badge.
//
// "public" is only ever returned when nothing in front of the route
// establishes a caller. It is the one answer that must not be produced by
// default, so it is produced last and only when every other field is empty.
func (a Access) Summary() string {
	var parts []string

	switch {
	case len(a.Roles) > 0 || len(a.Permissions) > 0:
		var demands []string
		demands = append(demands, a.Roles...)
		demands = append(demands, a.Permissions...)
		parts = append(parts, strings.Join(demands, " or "))
	case a.APIKeyOnly:
		parts = append(parts, "api key")
	case a.Staff:
		parts = append(parts, "any permission")
	case a.Authenticated:
		parts = append(parts, "signed in")
	case a.Optional:
		parts = append(parts, "optional")
	}

	if len(parts) == 0 {
		return "public"
	}
	return strings.Join(parts, " + ")
}

// guardKind is how a middleware affects access, for the ones that do.
type guardKind int

const (
	// guardNone is every middleware that does not decide access.
	guardNone guardKind = iota
	// guardSession establishes a caller from a session, bearer token or key.
	guardSession
	// guardAPIKey accepts an API key and no user session.
	guardAPIKey
	// guardOptional identifies a caller when there is one, and serves anonymous
	// requests otherwise.
	guardOptional
	// guardRole reads its demands from the call's string arguments.
	guardRole
	// guardStaff admits the ADMIN role or any permission at all.
	guardStaff
	// guardPermissionFor names a permission the route cannot know in advance,
	// because the resource is a path parameter.
	guardPermissionFor
)

// accessMiddleware is every middleware in the generated project that changes
// what a caller must present. Anything absent from this map is treated as not
// affecting access, which is the safe direction: an unknown middleware never
// turns a protected route public, it only fails to explain why a route is
// protected.
//
// Keyed by the function name as routes.go writes it, without the package
// qualifier.
var accessMiddleware = map[string]guardKind{
	"Auth":                 guardSession,
	"APIKeyOrAuth":         guardSession,
	"RequireAPIKey":        guardAPIKey,
	"Identify":             guardOptional,
	"RequireRole":          guardRole,
	"RequireStaff":         guardStaff,
	"RequirePermissionFor": guardPermissionFor,
}

// group is one router group as declared in the source.
type group struct {
	parent string
	prefix string
	// uses are the middleware calls applied to the whole group, in source order.
	uses []string
}

var (
	// groupRe matches `x := y.Group(expr)` and `x = y.Group(expr)`, capturing
	// the new variable, the receiver and the argument.
	groupRe = regexp.MustCompile(`(\w+)\s*:?=\s*(?:(\w+)\.)?Group\(([^()]*)\)`)
	// useRe matches the start of `x.Use(middleware.Name(` or `x.Use(gin.Name(`,
	// capturing the receiver and the middleware name. The arguments are read
	// separately, because a Use call can span lines.
	useRe = regexp.MustCompile(`\b(\w+)\.Use\(\s*(?:\w+\.)?(\w+)\(`)
	// routeStartRe matches `x.METHOD(` and nothing about what follows, because
	// what follows is read by a paren-aware scanner.
	routeStartRe = regexp.MustCompile(`\b(\w+)\.(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\(`)
	// callRe matches a middleware call such as `middleware.RequireRole(` inside
	// a route's argument list.
	callRe = regexp.MustCompile(`^\s*(?:\w+\.)?(\w+)\((.*)\)\s*$`)
)

// Parse reads a routes.go file and extracts all registered routes.
func Parse(routesFile string) ([]Route, error) {
	data, err := os.ReadFile(routesFile)
	if err != nil {
		return nil, fmt.Errorf("opening routes file: %w", err)
	}
	return ParseSource(string(data)), nil
}

// ParseSource is Parse against source already in memory, which is what the
// tests and the registry generator use.
func ParseSource(src string) []Route {
	routes, _ := parseFile(src, nil)
	sortRoutes(routes)
	return routes
}

// ParseProject reads routes.go and every per-resource route file beside it,
// which together are the whole table.
//
// Resources stopped writing into routes.go: each one owns
// internal/routes/<name>_routes.go and registers through the Mount struct, so
// its routes say m.Staff.GET(...) rather than staff.GET(...). Those files carry
// no group declarations of their own, so they are parsed with routes.go's group
// table seeded under the Mount field names, read from the Mount literal rather
// than assumed: the local variables are whatever routes.go chose to call them.
func ParseProject(routesFile string) ([]Route, error) {
	data, err := os.ReadFile(routesFile)
	if err != nil {
		return nil, fmt.Errorf("opening routes file: %w", err)
	}
	main := string(data)

	routes, groups := parseFile(main, nil)

	seed := map[string]*group{}
	for name, g := range groups {
		seed[name] = g
	}
	for field, local := range parseMountAliases(strings.Split(main, "\n")) {
		if _, ok := groups[local]; ok {
			seed[field] = &group{parent: local}
		}
	}

	files, err := FindResourceRouteFiles(routesFile)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		src, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("opening %s: %w", filepath.Base(file), err)
		}
		more, _ := parseFile(string(src), seed)
		routes = append(routes, more...)
	}

	sortRoutes(routes)
	return routes, nil
}

// mountFieldRe matches one line of the Mount literal, `Staff: staff,`.
var mountFieldRe = regexp.MustCompile(`^(V1|Public|Protected|Admin|Staff):\s*(\w+),?$`)

// parseMountAliases reads the Mount literal in routes.go and reports which
// local group variable each router field was handed.
func parseMountAliases(lines []string) map[string]string {
	out := map[string]string{}
	inLiteral := false
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if strings.Contains(line, "mountResources(&Mount{") {
			inLiteral = true
			continue
		}
		if !inLiteral {
			continue
		}
		if strings.HasPrefix(line, "})") {
			break
		}
		if m := mountFieldRe.FindStringSubmatch(line); len(m) == 3 {
			out[m[1]] = m[2]
		}
	}
	return out
}

func sortRoutes(routes []Route) {
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path == routes[j].Path {
			return routes[i].Method < routes[j].Method
		}
		return routes[i].Path < routes[j].Path
	})
}

// parseFile returns the routes declared in src plus the group table it built,
// starting from seed so a resource route file can resolve groups that routes.go
// declared.
func parseFile(src string, seed map[string]*group) ([]Route, map[string]*group) {
	lines := strings.Split(src, "\n")

	// String constants are collected up front because a group prefix may be
	// built from one: routes.go declares `const APIVersion = "v1"` and mounts
	// everything under r.Group("/api/" + APIVersion). Without resolving that,
	// the prefix silently evaluates to nothing and every route is reported one
	// path segment short, /users instead of /api/v1/users, which is worse than
	// no output at all because it looks plausible.
	consts := parseStringConsts(lines)

	// Pass one builds the group graph and each group's middleware. Nothing is
	// resolved yet: a group may be declared from a function parameter whose
	// value is assigned hundreds of lines later, so resolving as we read is how
	// the previous version lost path segments.
	groups := map[string]*group{
		"r": {}, // the engine itself, the root of every chain
	}
	for name, g := range seed {
		groups[name] = g
	}
	ensure := func(name string) *group {
		if g, ok := groups[name]; ok {
			return g
		}
		g := &group{}
		groups[name] = g
		return g
	}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		if m := groupRe.FindStringSubmatch(line); len(m) >= 4 {
			g := ensure(m[1])
			g.parent = m[2]
			g.prefix = resolveStringExpr(m[3], consts)
		}

		// A group can take middleware on a line of its own, and the generated
		// file does exactly that for the protected, staff and admin groups.
		if m := useRe.FindStringSubmatch(line); len(m) >= 3 {
			start := strings.Index(line, m[0]) + len(m[0]) - 1
			call := m[2] + balanced(line[start:])
			ensure(m[1]).uses = append(ensure(m[1]).uses, call)
		}
	}

	// Pass two reads the routes themselves.
	var routes []Route
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		m := routeStartRe.FindStringSubmatch(line)
		if len(m) < 3 {
			continue
		}
		recv, method := m[1], m[2]

		// The call is completed across lines if it does not close on this
		// one. Three routes in a generated project pass their handler as a
		// literal spanning four lines, and skipping them left three gaps in
		// the middle of the project's own API that the routes screen then
		// had to show as unknown.
		//
		// Only this call is joined, not every unbalanced line: a whole
		// RegisterRoutes(func(m *Mount) { ... }) body collapsed into one
		// line leaves only its first route findable.
		open := strings.Index(line, m[0]) + len(m[0]) - 1
		call, consumed := completeCall(lines, i, line[open:])
		if call == "" {
			continue
		}
		i += consumed

		args := splitArgs(trimParens(call))
		if len(args) < 2 {
			continue
		}

		// The path is a string expression, not necessarily a literal: the
		// health check registers "/api/" + APIVersion + "/health". Unquoting
		// the whole argument read that as the path /api/"+APIVersion+"/health.
		if !isStringExpr(args[0], consts) {
			continue
		}
		path := resolveStringExpr(args[0], consts)

		fullPath := joinPath(prefixOf(groups, recv, nil), path)

		// Group is the bucket the route is mounted in, read from the group
		// chain alone. Access is the whole requirement, chain plus the
		// middleware on the route itself. They are separate because a route on
		// the staff group naming perm:widgets.view belongs to staff and demands
		// that permission, and one column cannot say both.
		access := accessFor(groups, recv, args[1:len(args)-1])

		routes = append(routes, Route{
			Method:  method,
			Path:    fullPath,
			Handler: strings.TrimSpace(args[len(args)-1]),
			Group:   groupName(accessFor(groups, recv, nil)),
			Access:  access,
		})
	}

	return routes, groups
}

// prefixOf walks a group up to the root and returns the path it sits under.
//
// seen guards against a cycle. A group graph read out of source can contain one
// if a name is reused, and a parser that loops forever on a malformed file is a
// worse failure than one that stops early.
func prefixOf(groups map[string]*group, name string, seen map[string]bool) string {
	g, ok := groups[name]
	if !ok || name == "" {
		return ""
	}
	if seen == nil {
		seen = map[string]bool{}
	}
	if seen[name] {
		return ""
	}
	seen[name] = true
	return prefixOf(groups, g.parent, seen) + g.prefix
}

// chainOf returns every middleware in front of a group, outermost first.
func chainOf(groups map[string]*group, name string, seen map[string]bool) []string {
	g, ok := groups[name]
	if !ok || name == "" {
		return nil
	}
	if seen == nil {
		seen = map[string]bool{}
	}
	if seen[name] {
		return nil
	}
	seen[name] = true
	return append(chainOf(groups, g.parent, seen), g.uses...)
}

// accessFor reads a route's access requirement from its group chain and its own
// inline middleware.
func accessFor(groups map[string]*group, recv string, inline []string) Access {
	var out Access
	for _, call := range append(chainOf(groups, recv, nil), inline...) {
		applyGuard(&out, call)
	}
	return out
}

// applyGuard folds one middleware call into an access requirement.
func applyGuard(a *Access, call string) {
	m := callRe.FindStringSubmatch(strings.TrimSpace(call))
	if len(m) < 3 {
		return
	}
	name, rawArgs := m[1], m[2]

	kind, ok := accessMiddleware[name]
	if !ok || kind == guardNone {
		return
	}
	a.Guards = append(a.Guards, name)

	switch kind {
	case guardSession:
		a.Authenticated = true
		// A session guard supersedes an optional one: APIKeyOrAuth wrapping Auth
		// means a caller is required, whatever an Identify further out said.
		a.Optional = false
	case guardAPIKey:
		a.Authenticated = true
		a.APIKeyOnly = true
	case guardOptional:
		if !a.Authenticated {
			a.Optional = true
		}
	case guardStaff:
		a.Authenticated = true
		a.Staff = true
	case guardRole:
		a.Authenticated = true
		a.APIKeyOnly = false
		for _, arg := range splitArgs(rawArgs) {
			lit, ok := stringLit(arg)
			if !ok {
				continue
			}
			if perm, isPerm := strings.CutPrefix(lit, "perm:"); isPerm {
				a.Permissions = appendUnique(a.Permissions, perm)
				continue
			}
			a.Roles = appendUnique(a.Roles, lit)
		}
	case guardPermissionFor:
		a.Authenticated = true
		// The resource is a path parameter, so the permission is only knowable
		// per request. Naming it ":resource.view" says which half is dynamic
		// rather than pretending the route has no requirement.
		args := splitArgs(rawArgs)
		if len(args) == 2 {
			param, _ := stringLit(args[0])
			action, _ := stringLit(args[1])
			a.Permissions = appendUnique(a.Permissions, ":"+param+"."+action)
		}
	}
}

// groupName is the coarse bucket a route falls in, kept because the table and
// the generated TypeScript have always had a GROUP column.
//
// Role names are matched without regard to case here and only here. The
// generated code writes "ADMIN" and a hand-written project may write "admin";
// both mean the same bucket. The middleware's own comparison stays exact,
// because that one decides access rather than a column heading.
func groupName(a Access) string {
	switch {
	case a.Staff:
		return "staff"
	case containsAnyFold(a.Roles, "admin"):
		return "admin"
	case len(a.Permissions) > 0:
		return "staff"
	case a.APIKeyOnly:
		return "api-key"
	case a.Authenticated:
		return "protected"
	default:
		return "public"
	}
}

// maxJoin caps how many source lines one call may absorb.
//
// A file with an unclosed paren would otherwise swallow the rest of itself,
// and the point of a cap is that a malformed file costs one route rather than
// every route after it.
const maxJoin = 60

// completeCall returns the parenthesised call starting at head, continuing
// onto the following lines when it does not close on this one, and reports how
// many extra lines it consumed.
//
// An empty result means the call never closed, which a caller should treat as
// "skip this": guessing at half an argument list is how a wrong answer gets
// produced, and a wrong answer in this package becomes a wrong access column.
func completeCall(lines []string, i int, head string) (string, int) {
	if out := balanced(head); parenDepth(out) == 0 && strings.HasSuffix(out, ")") {
		return out, 0
	}

	joined := stripLineComment(head)
	for n := 1; n <= maxJoin && i+n < len(lines); n++ {
		joined += " " + strings.TrimSpace(stripLineComment(lines[i+n]))
		if out := balanced(joined); parenDepth(out) == 0 && strings.HasSuffix(out, ")") {
			return out, n
		}
	}
	return "", 0
}

// parenDepth is the net change in parenthesis nesting across s, ignoring
// parens inside string and rune literals.
func parenDepth(s string) int {
	depth := 0
	inStr := byte(0)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inStr != 0 {
			switch c {
			case '\\':
				i++
			case inStr:
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '\'', '`':
			inStr = c
		case '(':
			depth++
		case ')':
			depth--
		}
	}
	return depth
}

// stripLineComment removes a trailing // comment, leaving one inside a string
// literal alone. Joining lines without this would comment out the rest of the
// call, including the handler.
func stripLineComment(s string) string {
	inStr := byte(0)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inStr != 0 {
			switch c {
			case '\\':
				i++
			case inStr:
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '\'', '`':
			inStr = c
		case '/':
			if i+1 < len(s) && s[i+1] == '/' {
				return s[:i]
			}
		}
	}
	return s
}

// balanced returns s from its first opening paren through the matching close,
// inclusive, respecting string literals and nested brackets. If the parens
// never close on this line it returns what there is, so a caller can decide to
// skip rather than mis-split.
func balanced(s string) string {
	depth := 0
	inStr := byte(0)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inStr != 0 {
			switch c {
			case '\\':
				i++
			case inStr:
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '\'', '`':
			inStr = c
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			depth--
			if depth == 0 {
				return s[:i+1]
			}
		}
	}
	return s
}

// trimParens drops one matched pair of outer parens.
func trimParens(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '(' && s[len(s)-1] == ')' {
		return s[1 : len(s)-1]
	}
	return s
}

// splitArgs splits an argument list on top-level commas, respecting string
// literals and nested brackets. This is the whole reason route registrations
// with middleware in front of them parse now: their arguments contain commas
// inside parens, which is what the old pattern could not express.
func splitArgs(s string) []string {
	var out []string
	depth := 0
	inStr := byte(0)
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inStr != 0 {
			switch c {
			case '\\':
				i++
			case inStr:
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '\'', '`':
			inStr = c
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			depth--
		case ',':
			if depth == 0 {
				out = append(out, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	if rest := strings.TrimSpace(s[start:]); rest != "" {
		out = append(out, rest)
	}
	return out
}

// isStringExpr reports whether expr is something resolveStringExpr can
// evaluate: a quoted literal, a known constant, or a concatenation of those.
//
// It exists so a route whose path this parser cannot work out is skipped rather
// than reported at a path that is subtly wrong, which is the failure mode this
// package keeps having.
func isStringExpr(expr string, consts map[string]string) bool {
	parts := strings.Split(expr, "+")
	if len(parts) == 0 {
		return false
	}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return false
		}
		if _, ok := stringLit(part); ok {
			continue
		}
		if _, ok := consts[part]; ok {
			continue
		}
		return false
	}
	return true
}

// stringLit unquotes a double-quoted Go string literal, reporting whether the
// argument was one at all.
func stringLit(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1], true
	}
	return "", false
}

// joinPath concatenates a group prefix and a route path without producing a
// double slash or losing a root path.
func joinPath(prefix, path string) string {
	full := prefix + path
	for strings.Contains(full, "//") {
		full = strings.ReplaceAll(full, "//", "/")
	}
	if full == "" {
		return "/"
	}
	return full
}

func appendUnique(list []string, v string) []string {
	for _, existing := range list {
		if existing == v {
			return list
		}
	}
	return append(list, v)
}

func containsAnyFold(list []string, want ...string) bool {
	for _, v := range list {
		for _, w := range want {
			if strings.EqualFold(v, w) {
				return true
			}
		}
	}
	return false
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

// FindResourceRouteFiles returns every per-resource route file beside routes.go.
//
// Resources stopped writing into routes.go in v3.2xx: each one owns
// internal/routes/<name>_routes.go. A route table built from routes.go alone
// therefore misses every generated resource, which is most of what a real
// project serves.
func FindResourceRouteFiles(routesFile string) ([]string, error) {
	dir := filepath.Dir(routesFile)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}

	var out []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, "_routes.go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		out = append(out, filepath.Join(dir, name))
	}
	sort.Strings(out)
	return out, nil
}

// FormatTable formats routes as a printable table.
func FormatTable(routes []Route) string {
	if len(routes) == 0 {
		return "  No routes found."
	}

	headers := []string{"METHOD", "PATH", "HANDLER", "GROUP", "ACCESS"}
	rows := make([][]string, 0, len(routes))
	for _, r := range routes {
		rows = append(rows, []string{r.Method, r.Path, r.Handler, r.Group, r.Access.Summary()})
	}

	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	line := func(cells []string) string {
		parts := make([]string, len(cells))
		for i, cell := range cells {
			parts[i] = fmt.Sprintf("%-*s", widths[i], cell)
		}
		return "  " + strings.TrimRight(strings.Join(parts, "  "), " ")
	}

	var b strings.Builder
	b.WriteString(line(headers) + "\n")

	dividers := make([]string, len(headers))
	for i := range headers {
		dividers[i] = strings.Repeat("─", widths[i])
	}
	b.WriteString(line(dividers) + "\n")

	for _, row := range rows {
		b.WriteString(line(row) + "\n")
	}

	fmt.Fprintf(&b, "\n  %d routes total\n", len(routes))
	return b.String()
}
