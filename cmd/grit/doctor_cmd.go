package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/doctor"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// doctorCmd is `grit doctor`: an audit of the project for the mistakes that do
// not announce themselves.
//
// Separate from `grit sync doctor`, which checks one thing (the offline sync
// configuration) and keeps its own command.
func doctorCmd() *cobra.Command {
	var asJSON bool
	var queryLog string

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Audit the project for silent security and configuration mistakes",
		Long: "Checks the project for the mistakes that do not announce themselves: an\n" +
			"encrypted field with no key, a resource nothing scopes to its owner, a table\n" +
			"shared across organizations, a database browser with no login, dashboard\n" +
			"credentials still at their defaults.\n\n" +
			"Exits non-zero when anything is an error, so CI can run it.\n\n" +
			"Everything it reads is a file, so it needs no database and no running\n" +
			"server. The one exception reads a file a run produced: record one with\n" +
			"GRIT_QUERY_LOG=.grit/queries.jsonl and it reports the same query issued\n" +
			"many times inside one request, which is what an N+1 looks like.",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := scaffold.FindProjectRoot()
			if err != nil {
				return err
			}
			report, err := doctor.RunWith(root, doctor.Options{QueryLog: queryLog})
			if err != nil {
				return err
			}

			if asJSON {
				out, err := json.MarshalIndent(report, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(out))
				if report.Errors() > 0 {
					os.Exit(1)
				}
				return nil
			}

			red := color.New(color.FgHiRed)
			yellow := color.New(color.FgHiYellow)
			cyan := color.New(color.FgHiCyan)
			green := color.New(color.FgHiGreen)
			gray := color.New(color.FgHiBlack)

			fmt.Println()
			cyan.Printf("  %d checks over %d resource(s)\n", report.Checks, len(report.Resources))
			fmt.Println()

			for _, f := range report.Findings {
				printer, mark := yellow, "⚠"
				if f.Level == "error" {
					printer, mark = red, "✗"
				}
				if f.Resource != "" {
					printer.Printf("  %s %s %s\n", mark, f.Resource, f.Message)
				} else {
					printer.Printf("  %s %s\n", mark, f.Message)
				}
				if f.Fix != "" {
					gray.Printf("      %s\n", f.Fix)
				}
				gray.Printf("      (%s)\n", f.Check)
			}

			if len(report.Findings) == 0 {
				green.Println("  ✓ Nothing to report.")
			} else {
				fmt.Println()
				gray.Printf("  %d error(s), %d warning(s)\n", report.Errors(), report.Warnings())
			}
			fmt.Println()

			if report.Errors() > 0 {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&asJSON, "json", false, "Print the report as JSON, for CI")
	cmd.Flags().StringVar(&queryLog, "queries", "",
		"A recorded query log to read (default .grit/queries.jsonl when it exists). "+
			"Record one with GRIT_QUERY_LOG=<path> while the app or its tests run.")
	return cmd
}
