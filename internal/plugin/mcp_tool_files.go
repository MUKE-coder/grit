package plugin

// mcpToolGo is internal/mcp/tool.go: what a tool is, and the registry.
func mcpToolGo(ctx Context) string {
	return `package mcp

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// Tool is one thing an agent can do.
//
// A struct and not an interface, because every field of it is read by two
// consumers: MCP's tools/list, and the in-app assistant. One catalogue means a
// tool written once works in both, and it means the permission a tool needs is
// declared in the same place as the routes it calls, so a connection's key can
// be narrowed to exactly those routes without anybody maintaining a second
// list.
type Tool struct {
	// Name is what the client calls. Snake case, verb first: list_invoices,
	// create_invoice. Stable: renaming one breaks every saved prompt.
	Name string
	// Title is for a person choosing tools in the app.
	Title string
	// Description is for the model. Say what it does and when to reach for it;
	// a model picks tools by reading these.
	Description string
	// Group buckets the tool in the connection screen: "Invoices", "Clients".
	Group string

	// Permission is the key a role must grant to use it, the same string roles
	// grant everywhere else. Empty means any signed-in caller.
	Permission string

	// Write marks a tool that changes something. It drives MCP's
	// readOnlyHint, and it is what lets a connection be made read-only.
	Write bool
	// Destructive marks one that removes or voids. It drives
	// destructiveHint, so a client can confirm before calling.
	Destructive bool

	// Routes are the method and path each call reaches, as "GET /api/v1/invoices".
	// A connection's key is restricted to the union of its tools' routes, so
	// this is a security boundary and not documentation: a tool that omits a
	// route it calls will be refused at runtime by the key's own allowlist.
	Routes []string

	// Schema is the JSON Schema for the arguments, as the object the client
	// sends. Built from a resource's fields by the generator, or written by
	// hand for a custom tool.
	Schema map[string]any

	// Exec runs it. The Call it is given carries the connection's identity and
	// replays requests through the app's own router: see call.go.
	Exec func(ctx context.Context, call *Call, args map[string]any) (any, error)
}

// ReadOnly is the inverse of Write, for MCP's readOnlyHint.
func (t Tool) ReadOnly() bool { return !t.Write }

// registry is the catalogue. One per process, filled from init()s, so a tools
// file registers itself and nothing has to maintain a list of them.
var registry = struct {
	mu    sync.RWMutex
	tools map[string]Tool
}{tools: map[string]Tool{}}

// Register adds a tool. Called from an init() in a tools file.
//
// A duplicate name panics, on purpose and at startup. Two tools with one name
// means one of them is unreachable, and which one would depend on the order
// the files happened to initialise in: a silent coin toss over what an agent
// can do to your data.
func Register(tools ...Tool) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	for _, tool := range tools {
		if tool.Name == "" {
			panic("mcp: a tool with no name")
		}
		if _, clash := registry.tools[tool.Name]; clash {
			panic(fmt.Sprintf("mcp: two tools named %q", tool.Name))
		}
		registry.tools[tool.Name] = tool
	}
}

// Lookup is one tool by name.
func Lookup(name string) (Tool, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	tool, ok := registry.tools[name]
	return tool, ok
}

// All is every registered tool, sorted by name.
//
// Sorted because tools/list goes to a model: a list that reshuffles between
// two identical sessions is a prompt that cannot be cached and a diff that
// cannot be read.
func All() []Tool {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	names := make([]string, 0, len(registry.tools))
	for name := range registry.tools {
		names = append(names, name)
	}
	sort.Strings(names)

	out := make([]Tool, 0, len(names))
	for _, name := range names {
		out = append(out, registry.tools[name])
	}
	return out
}

// Named is the tools a connection chose, in registry order, skipping any name
// that no longer exists.
//
// Skipping rather than failing: a connection made before a resource was
// removed names a tool that is gone, and the useful behaviour is for its other
// forty tools to go on working. ToolsMissing reports the gap for the screen.
func Named(names []string) []Tool {
	want := make(map[string]bool, len(names))
	for _, name := range names {
		want[name] = true
	}
	var out []Tool
	for _, tool := range All() {
		if want[tool.Name] {
			out = append(out, tool)
		}
	}
	return out
}

// ToolsMissing is the names a connection has that the registry does not, so
// the screen can say "three of its tools no longer exist" rather than quietly
// offering fewer.
func ToolsMissing(names []string) []string {
	var missing []string
	for _, name := range names {
		if _, ok := Lookup(name); !ok {
			missing = append(missing, name)
		}
	}
	sort.Strings(missing)
	return missing
}

// RoutesFor is the union of the routes a set of tools calls, sorted and
// deduplicated.
//
// This is what a connection's API key is restricted to. Derived rather than
// chosen: a person picking tools in the app is not also asked to pick routes,
// and cannot widen the key by accident.
func RoutesFor(tools []Tool) []string {
	seen := map[string]bool{}
	var routes []string
	for _, tool := range tools {
		for _, route := range tool.Routes {
			if route == "" || seen[route] {
				continue
			}
			seen[route] = true
			routes = append(routes, route)
		}
	}
	sort.Strings(routes)
	return routes
}

// PermissionsFor is the permissions a set of tools needs, so the screen can
// refuse to offer a tool the person's own role does not grant. An agent acting
// as somebody must not be able to do more than they can.
func PermissionsFor(tools []Tool) []string {
	seen := map[string]bool{}
	var perms []string
	for _, tool := range tools {
		if tool.Permission == "" || seen[tool.Permission] {
			continue
		}
		seen[tool.Permission] = true
		perms = append(perms, tool.Permission)
	}
	sort.Strings(perms)
	return perms
}
`
}
