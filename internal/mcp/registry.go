package mcp

import (
	"encoding/json"
	"fmt"
	"sort"
)

// tool is one thing an agent can call.
//
// Description is written for a model deciding whether to call something, so each
// one says what the answer is good for and what it cannot tell you. "Parsed from
// source, not from a running server" is in most of them for a reason: an agent
// that took these for live readings would draw confident conclusions from a
// stale checkout.
type tool struct {
	name        string
	description string
	schema      map[string]interface{}
	// mutates marks a tool that writes to the project. A read-only server does
	// not register it. See mode.go.
	mutates bool
	run     func(*Server, json.RawMessage) (string, error)
}

// noArgs is the schema for a tool that takes nothing.
var noArgs = map[string]interface{}{
	"type":       "object",
	"properties": map[string]interface{}{},
}

// allTools is every tool this server knows how to answer, in both modes.
//
// One list. tools/list filters it by mode and tools/call looks a name up in the
// same filtered set, so the two cannot disagree about what exists.
func allTools() []tool {
	return []tool{
		{
			name: "grit_project_info",
			description: "Describe the Grit project: architecture (single/double/triple/api/mobile), " +
				"frontend framework, Go module path, CLI version it was scaffolded with, and which " +
				"apps exist. Call this first to learn the layout before assuming where files live. " +
				"Read from grit.json on disk.",
			schema: noArgs,
			run:    (*Server).projectInfoTool,
		},
		{
			name: "grit_list_routes",
			description: "List every registered API route with its HTTP method, full path " +
				"(including the /api/v1 prefix), handler, and middleware group (public, protected, " +
				"admin). Use this to find the correct URL and required auth level before writing a " +
				"client call. Parsed from routes.go, so it reflects the checkout rather than a " +
				"running server.",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"method": map[string]interface{}{
						"type":        "string",
						"description": "Optional HTTP method filter, e.g. GET or POST.",
					},
					"contains": map[string]interface{}{
						"type":        "string",
						"description": "Optional case-insensitive substring filter on the path, e.g. \"users\".",
					},
				},
			},
			run: (*Server).listRoutes,
		},
		{
			name: "grit_describe_models",
			description: "List the GORM models with their fields, Go types, JSON names, and GORM " +
				"tags. Use this to learn the exact shape of a request or response body, the column " +
				"constraints, and the relationships between tables. Parsed from internal/models. " +
				"Note: this is the model, which is what the authenticated endpoints return. The " +
				"public endpoints a resource generated with --public serves return an allowlist " +
				"held in internal/handlers/<name>_public.go instead, which is deliberately " +
				"narrower and is edited by hand; read that file for the public shape. A money " +
				"field is two columns, <name>_amount and <name>_currency, and the sortable and " +
				"filterable one is <name>_amount.",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"model": map[string]interface{}{
						"type":        "string",
						"description": "Optional model name, e.g. \"User\". Omit to list them all.",
					},
				},
			},
			run: (*Server).describeModels,
		},
		{
			name: "grit_list_resources",
			description: "List the resources `grit generate resource` has created in this project, " +
				"with the model, service, handler, admin page and hooks each one owns. Use this to " +
				"find out whether a thing already exists before generating it, and to learn which " +
				"files a change to one resource has to touch. Read from the files on disk.",
			schema: noArgs,
			run:    (*Server).listResources,
		},
		{
			name: "grit_file_ownership",
			description: "Report, for the files Grit generated, which are still exactly as it " +
				"wrote them, which the developer has edited, and which have been deleted. Nothing " +
				"else can answer this: it comes from the manifest Grit records when it writes. Use " +
				"it before proposing an edit, because changing a pristine file means `grit upgrade` " +
				"will stop updating it, and before `grit upgrade`, because a modified file is where " +
				"a conflict will be reported.",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"status": map[string]interface{}{
						"type":        "string",
						"description": "Optional filter: unchanged, modified, missing. Omit for all of them.",
					},
					"contains": map[string]interface{}{
						"type":        "string",
						"description": "Optional case-insensitive substring filter on the path.",
					},
				},
			},
			run: (*Server).fileOwnership,
		},
		{
			name: "grit_doctor",
			description: "Run Grit's project audit and return the findings as data: security and " +
				"configuration mistakes that compile fine and are wrong, each with the check that " +
				"found it, what it is about, and the fix. Use it before claiming a project is ready " +
				"to deploy, and after generating a resource that holds personal data. The same " +
				"checks as the `grit doctor` command.",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"level": map[string]interface{}{
						"type":        "string",
						"description": "Optional filter: error or warning. Omit for both.",
					},
				},
			},
			run: (*Server).doctorReport,
		},
		{
			name: "grit_list_permissions",
			description: "List every permission key this application understands, grouped as the " +
				"admin's permission tree groups them. Use this to find the exact string to pass to " +
				"a route guard or a role, instead of inventing one that will silently never match. " +
				"Parsed from internal/authz/permissions.go, which is the project's own catalogue " +
				"and is meant to be edited.",
			schema: noArgs,
			run:    (*Server).listPermissions,
		},
		{
			name: "grit_env_keys",
			description: "List the environment variables this project reads, with the comment that " +
				"documents each one and whether .env sets it. Values are never returned, only " +
				"whether something is set, because this file holds the database password and the " +
				"signing keys. Use it to find the right variable name and to see what a deployment " +
				"is missing.",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"contains": map[string]interface{}{
						"type":        "string",
						"description": "Optional case-insensitive substring filter on the key, e.g. \"mail\".",
					},
				},
			},
			run: (*Server).envKeys,
		},
		{
			name: "grit_cli_reference",
			description: "List the Grit CLI's commands and flags, from the binary that is running " +
				"this server. Use it to propose a command that exists with flags that exist, " +
				"instead of one from a different version. Changing this project is done through " +
				"these commands: they write the code and inject the wiring that makes it reachable.",
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "Optional command path, e.g. \"generate resource\". Omit to list them all.",
					},
				},
			},
			run: (*Server).cliReference,
		},
		{
			name: "grit_generate_resource",
			description: "Generate a full resource: the GORM model, service, handler and routes, " +
				"the Zod schema and TypeScript types, the React Query hooks and the admin page, " +
				"with the wiring injected into the router and the registries. This WRITES FILES. " +
				"Prefer it over hand-writing those files: hand-written ones compile and are missing " +
				"the injections that make them reachable. Run `grit_list_resources` first to check " +
				"the name is free.",
			mutates: true,
			schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The resource name, singular and PascalCase, e.g. \"Product\".",
					},
					"fields": map[string]interface{}{
						"type": "string",
						"description": "Field definitions, e.g. \"title:string,price:money,published:bool\". " +
							"Call grit_cli_reference for the field types this version supports.",
					},
					"owned_by": map[string]interface{}{
						"type": "string",
						"description": "Optional: scope every row to its owner through this belongs_to-User " +
							"field, e.g. \"user\". Rows are then readable and writable only by their owner.",
					},
					"seed": map[string]interface{}{
						"type":        "boolean",
						"description": "Also generate a seeder with one example record.",
					},
				},
				"required": []string{"name"},
			},
			run: (*Server).generateResource,
		},
	}
}

// tools is the set this server exposes, which is what its mode allows.
func (s *Server) tools() map[string]tool {
	out := map[string]tool{}
	for _, t := range allTools() {
		if !s.mode().allows(t.mutates) {
			continue
		}
		out[t.name] = t
	}
	return out
}

func (s *Server) mode() Mode {
	if s.Mode == "" {
		return ModeRead
	}
	return s.Mode
}

// toolDefinitions is the tools/list payload: the tools this server has, in a
// stable order so a client's own logging does not churn.
func (s *Server) toolDefinitions() []map[string]interface{} {
	set := s.tools()
	names := make([]string, 0, len(set))
	for name := range set {
		names = append(names, name)
	}
	sort.Strings(names)

	out := make([]map[string]interface{}, 0, len(names))
	for _, name := range names {
		t := set[name]
		out = append(out, map[string]interface{}{
			"name":        t.name,
			"description": t.description,
			"inputSchema": t.schema,
		})
	}
	return out
}

// callTool dispatches a tools/call.
//
// The lookup is the same map tools/list was built from, so a tool this server
// does not have is not callable even by a client that knows its name from
// somewhere else.
func (s *Server) callTool(name string, args json.RawMessage) (string, error) {
	t, ok := s.tools()[name]
	if !ok {
		// Say which it is. A model that asked for a generator on a read-only
		// server should be told to ask the human for --mode write rather than
		// conclude the tool does not exist and hand-write the files.
		for _, known := range allTools() {
			if known.name == name {
				return "", fmt.Errorf("%s writes files and this server is running read-only; "+
					"start it with --mode write to allow it, or run the equivalent CLI command yourself", name)
			}
		}
		return "", fmt.Errorf("unknown tool %q", name)
	}
	return t.run(s, args)
}
