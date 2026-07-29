package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/mcp"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

func mcpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Expose the project to AI coding agents over MCP",
		Long: "The Model Context Protocol lets an AI coding agent ask Grit about your project\n" +
			"instead of guessing from a grep. Point your agent at `grit mcp serve` and it can\n" +
			"read the real route table, the real model definitions, and the real layout.\n\n" +
			"Every tool is read-only and answers from your source files — no running server,\n" +
			"no database, no credentials. The agent still has to call the CLI to change\n" +
			"anything, so every change shows up in your diff.",
	}
	cmd.AddCommand(mcpServeCmd())
	return cmd
}

func mcpServeCmd() *cobra.Command {
	var projectDir string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the MCP server on stdio",
		Long: "Speaks the Model Context Protocol over stdin/stdout. Run by an MCP client, not\n" +
			"usually by hand.\n\n" +
			"Register it with Claude Code:\n\n" +
			"  claude mcp add grit -- grit mcp serve --project /path/to/project\n\n" +
			"Or add it to an MCP client config:\n\n" +
			"  {\n" +
			"    \"mcpServers\": {\n" +
			"      \"grit\": {\n" +
			"        \"command\": \"grit\",\n" +
			"        \"args\": [\"mcp\", \"serve\", \"--project\", \"/path/to/project\"]\n" +
			"      }\n" +
			"    }\n" +
			"  }\n\n" +
			"Tools: grit_project_info, grit_list_routes, grit_describe_models.",
		// This command is launched by an MCP client, which surfaces stderr in a
		// log. Dumping the full usage block after a runtime error buries the one
		// line that says what went wrong. Flag errors still print usage.
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			root := projectDir
			if root == "" {
				detected, err := scaffold.FindProjectRoot()
				if err != nil {
					return err
				}
				root = detected
			}
			if _, err := os.Stat(root); err != nil {
				return fmt.Errorf("project directory %s: %w", root, err)
			}

			// stdout carries the protocol and nothing else — a banner here
			// would desynchronise the client's JSON stream. The startup line
			// goes to stderr, where clients log it and users can see it.
			fmt.Fprintf(os.Stderr, "grit mcp serve — project %s\n", root)

			srv := &mcp.Server{Root: root, Version: version}
			return srv.Serve(os.Stdin, os.Stdout)
		},
	}

	cmd.Flags().StringVar(&projectDir, "project", "",
		"Project root (defaults to searching upward from the working directory)")

	return cmd
}
