package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/MUKE-coder/grit/v3/internal/mcp"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

func mcpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Expose the project to AI coding agents over MCP",
		Long: "The Model Context Protocol lets an AI coding agent ask Grit about your project\n" +
			"instead of guessing from a grep. Point your agent at `grit mcp serve` and it can\n" +
			"read the real route table, the real model definitions, the real permission\n" +
			"catalogue, which files you have edited, and what `grit doctor` finds.\n\n" +
			"The default server is read-only, and not by checking a flag: the tools that\n" +
			"write files are not in it. `--mode write` builds the server that has them.",
	}
	cmd.AddCommand(mcpServeCmd())
	return cmd
}

func mcpServeCmd() *cobra.Command {
	var projectDir string
	var mode string

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
			"Read mode (the default) answers questions: grit_project_info, grit_list_routes,\n" +
			"grit_describe_models, grit_list_resources, grit_list_permissions,\n" +
			"grit_file_ownership, grit_env_keys, grit_doctor, grit_cli_reference.\n\n" +
			"Write mode adds the generators, which write files: grit_generate_resource.\n" +
			"They are absent from a read-only server rather than refused by it, so nothing\n" +
			"a client sends can reach one.",
		// This command is launched by an MCP client, which surfaces stderr in a
		// log. Dumping the full usage block after a runtime error buries the one
		// line that says what went wrong. Flag errors still print usage.
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			serverMode, err := mcp.ParseMode(mode)
			if err != nil {
				return err
			}

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
			// goes to stderr, where clients log it and users can see it. It
			// names the mode, because a server that can write files should say
			// so somewhere a person will read.
			fmt.Fprintf(os.Stderr, "grit mcp serve — project %s, mode %s\n", root, serverMode)

			srv := &mcp.Server{
				Root:     root,
				Version:  version,
				Mode:     serverMode,
				Commands: commandTree(cmd.Root()),
			}
			return srv.Serve(os.Stdin, os.Stdout)
		},
	}

	cmd.Flags().StringVar(&projectDir, "project", "",
		"Project root (defaults to searching upward from the working directory)")
	cmd.Flags().StringVar(&mode, "mode", "read",
		"read exposes the tools that answer questions; write also exposes the generators, which write files")

	return cmd
}

// commandTree flattens cobra's commands into what grit_cli_reference reports.
//
// Read from the running binary rather than from a list, so an agent is told
// what this version accepts. A list would be a second copy of the command tree,
// and the first thing it would do is go out of date.
func commandTree(root *cobra.Command) []mcp.CommandInfo {
	var out []mcp.CommandInfo

	var walk func(cmd *cobra.Command, prefix string)
	walk = func(cmd *cobra.Command, prefix string) {
		for _, child := range cmd.Commands() {
			if child.Hidden || child.Name() == "help" || child.Name() == "completion" {
				continue
			}
			path := strings.TrimSpace(prefix + " " + child.Name())
			// A command that only groups others, such as `grit generate`, is
			// still worth listing: it names the family an agent is looking for.
			out = append(out, mcp.CommandInfo{
				Path:  path,
				Short: child.Short,
				Long:  child.Long,
				Flags: flagList(child),
			})
			walk(child, path)
		}
	}
	walk(root, "")
	return out
}

func flagList(cmd *cobra.Command) []mcp.FlagInfo {
	var out []mcp.FlagInfo
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" {
			return
		}
		out = append(out, mcp.FlagInfo{
			Name:    "--" + f.Name,
			Type:    f.Value.Type(),
			Usage:   f.Usage,
			Default: f.DefValue,
		})
	})
	return out
}
