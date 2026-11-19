package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func init() { Register(mcpPlugin()) }

// mcpPlugin gives a project an MCP server, so an agent can use the app through
// the same rules a person does.
//
// Design, and the one decision that matters most:
//
//   - **A tool call goes through the app's own router.** Exec builds an HTTP
//     request and runs it in process against the gin engine, with the
//     connection's key and organization header. So authentication, tenancy,
//     role and policy checks, module switches, record scopes, validation and
//     the audit log all apply exactly as they do for a browser, and they are
//     checked again on every call. A tool that talked to GORM directly would
//     be a second implementation of every rule in the app, and it would drift.
//
//   - **One connection is one key, never a shared secret.** A person makes a
//     connection in the app and it mints a Grit secret API key that acts as
//     them, limited by the key's endpoint allowlist to exactly the routes the
//     chosen tools call. Revoking the connection revokes the key. The
//     alternative, one MCP_API_KEY for everything, cannot say who did what and
//     cannot be narrowed.
//
//   - **Stateless transport.** Streamable HTTP at POST /mcp, mounted outside
//     the /api groups because it authenticates itself. Every request carries
//     its key, so there is no session to keep, lose or pin to one process.
//
//   - **One registry, two consumers.** A Tool is a struct, and the in-app
//     assistant and MCP read the same catalogue, so a tool written once works
//     in both.
//
// Written from two hand-built implementations, Atlas CRM and Outpost, the
// second of which had about sixty tools. grit#224 has the full reference.
func mcpPlugin() Plugin {
	return Plugin{
		Name:    "mcp",
		Version: "1.0.0",
		Summary: "An MCP server for your app, so an agent works through the same rules a person does",
		Description: `Exposes your app to Claude, Cursor and any other MCP client, at POST /mcp.

  • Streamable HTTP, stateless: protocol 2025-11-25 down to 2024-11-05,
    with initialize, ping, tools/list and tools/call
  • One connection is one API key, acting as the person who made it and
    narrowed by the key's endpoint allowlist to the routes its tools call
  • Every tool call runs in process through your own gin engine, so tenancy,
    roles, policies, validation and the audit log apply unchanged
  • Every call is recorded: who, which tool, the arguments, the status and
    how long it took
  • A tool registry your own code can add to, for the actions that are not
    CRUD

Generated CRUD tools per resource come next; this release is the server, the
connections, and the registry hand-written tools register into.

Requires API keys (Grit v3.300.0+) and the activity log.`,

		NextSteps: []string{
			"Rebuild the API:  cd apps/api && go build ./...",
			"Run the migration: grit migrate",
			"Make a connection: POST /api/v1/mcp-connections with a name and its tools",
			"Point a client at it: see /docs/plugins/mcp for the mcpServers entry",
		},

		Requires:   nil,
		Preflight:  mcpPreflight,
		Files:      mcpFiles,
		Injections: mcpInjections,
	}
}

// mcpPreflight refuses a project that cannot mint the keys a connection needs.
//
// The plugin's whole auth model is "one connection is one API key", so a
// project without the API key service would install, compile, and then be
// unable to make a single connection. Said before anything is written, with
// the fix, because an older project can be brought up to this.
func mcpPreflight(ctx Context) error {
	for _, need := range []struct{ rel, decl, why string }{
		{"internal/services/api_key.go", "func GenerateAPIKey(",
			"which is what mints a connection's key"},
		{"internal/models/api_key.go", "func (k *APIKey) AllowsEndpoint(",
			"so a connection could not be narrowed to the routes its own tools call"},
	} {
		if !declares(ctx, need.rel, need.decl) {
			return fmt.Errorf("%s has no %s, %s\n\nRun `grit upgrade` first.",
				need.rel, strings.TrimSuffix(strings.TrimPrefix(need.decl, "func "), "("), need.why)
		}
	}
	return nil
}

// declares reports whether a file in the project's API declares something.
//
// On the declaration and not the bare name, because a doc comment that
// mentions GenerateAPIKey would satisfy a looser check and the preflight would
// pass on a project that cannot mint a key.
func declares(ctx Context, rel, decl string) bool {
	body, err := os.ReadFile(filepath.Join(ctx.Root, apiPath(ctx, rel)))
	if err != nil {
		return false
	}
	return strings.Contains(string(body), decl)
}
