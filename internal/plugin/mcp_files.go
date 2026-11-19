package plugin

// mcpFiles is everything the plugin writes.
//
// The API half only. The connection screens are a later release; until then a
// connection is made over the API, which is also how a headless project would
// do it.
func mcpFiles(ctx Context) map[string]string {
	p := func(rel string) string { return apiPath(ctx, rel) }

	return map[string]string{
		p("internal/mcp/tool.go"):                     mcpToolGo(ctx),
		p("internal/mcp/transport.go"):                mcpTransportGo(ctx),
		p("internal/mcp/call.go"):                     mcpCallGo(ctx),
		p("internal/mcp/args.go"):                     mcpArgsGo(ctx),
		p("internal/mcp/mount.go"):                    mcpMountGo(ctx),
		p("internal/routes/mcp_wire.go"):              mcpWireGo(ctx),
		p("internal/mcp/tools_builtin.go"):            mcpBuiltinToolsGo(ctx),
		p("internal/mcp/custom_tools.go"):             mcpCustomToolsGo(ctx),
		p("internal/mcp/transport_test.go"):           mcpTransportTestGo(ctx),
		p("internal/mcp/args_test.go"):                mcpArgsTestGo(ctx),
		p("internal/models/mcp.go"):                   mcpModelsGo(ctx),
		p("internal/services/mcp_connection.go"):      mcpConnectionServiceGo(ctx),
		p("internal/handlers/mcp_connection.go"):      mcpConnectionHandlerGo(ctx),
		p("internal/services/mcp_connection_test.go"): mcpConnectionServiceTestGo(ctx),
	}
}

// mcpInjections wires it into the app.
func mcpInjections(ctx Context) []Injection {
	return []Injection{
		{
			// The registration list lives in models/user.go, which is where
			// Models() is declared, not in a models.go that does not exist.
			File:   apiPath(ctx, "internal/models/user.go"),
			Marker: "// grit:models",
			Code:   "\t\t&MCPConnection{},\n\t\t&MCPToolCall{},",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:handlers",
			Code:   "\tmcpConnectionHandler := handlers.NewMCPConnectionHandler(db)",
		},
		{
			// Mounted at the root and not under /api, because it carries its
			// own credential and answers JSON-RPC rather than REST. The engine
			// is handed over so a tool can replay a request through it: see
			// internal/mcp/call.go for why that is the whole design.
			//
			// One line, with the wiring in routes/mcp_wire.go. internal/mcp
			// cannot import internal/services, because the service that
			// creates a connection imports internal/mcp to read its tools, so
			// the adapter has to live in a package above both.
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:routes:custom",
			Code: "\t// The MCP server. Outside the /api groups: it authenticates itself.\n" +
				"\tmountMCP(r, db)",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:routes:protected",
			Code: "\t\tconnections := protected.Group(\"/mcp-connections\")\n" +
				"\t\t{\n" +
				"\t\t\tconnections.GET(\"\", mcpConnectionHandler.List)\n" +
				"\t\t\tconnections.POST(\"\", mcpConnectionHandler.Create)\n" +
				"\t\t\tconnections.DELETE(\"/:id\", mcpConnectionHandler.Revoke)\n" +
				"\t\t\tconnections.GET(\"/:id/calls\", mcpConnectionHandler.Calls)\n" +
				"\t\t\tconnections.GET(\"/tools\", mcpConnectionHandler.Tools)\n" +
				"\t\t}",
		},
	}
}
