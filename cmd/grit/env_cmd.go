package main

import (
	"fmt"
	"slices"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// envCmd is `grit env`: a working .env for a project you cloned.
func envCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "env",
		Short: "Create .env from .env.example with fresh secrets (for a project you cloned)",
		Long: "A cloned project has .env.example, with CHANGE_ME where each secret goes, and no .env,\n" +
			"and the API refuses to start on a placeholder. grit env copies .env.example to .env and\n" +
			"generates every CHANGE_ME the way grit new does. Run it again later and it fills only the\n" +
			"placeholders still left in .env; it never changes a value you have set, and never prints one.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := scaffold.FindProjectRoot()
			if err != nil {
				return err
			}
			created, filled, err := scaffold.FillEnv(root)
			if err != nil {
				return err
			}
			green := color.New(color.FgGreen)
			if created {
				green.Println("  ✓ .env created from .env.example")
			}
			if len(filled) > 0 {
				green.Printf("  ✓ %s\n", scaffold.EnvFillSummary(filled))
			} else if !created {
				fmt.Println("  .env has no placeholders left: nothing to do.")
			}
			if slices.Contains(filled, "FIELD_ENCRYPTION_KEY") {
				color.New(color.FgYellow).Println("  ⚠ FIELD_ENCRYPTION_KEY is new. That is right for a database of your own. To share a database")
				color.New(color.FgYellow).Println("    that already holds encrypted data, put the team's key in .env instead: this one cannot read it.")
			}
			return nil
		},
	}
}
