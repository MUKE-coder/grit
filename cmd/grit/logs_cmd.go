package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/devlog"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// logsCmd is `grit logs`: what the development servers printed.
//
// For a person as well as an agent. The reason it exists for both is the same:
// a Grit application logs to stderr, so when `grit start` is in another
// terminal or has scrolled, what it said is gone. The failures that take
// longest to find are the ones that never reach an exit code, and they are all
// in here.
func logsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Show what the development servers printed",
		Long: "`grit start` records everything the API and the frontends print into\n" +
			".grit/logs/dev.log, and keeps the previous run beside it.\n\n" +
			"  grit logs              the tail of this run\n" +
			"  grit logs --errors     only the failures, with their stack traces\n" +
			"  grit logs --grep 500   lines matching something",
	}

	var lines int
	var errorsOnly bool
	var grep string
	var asJSON bool

	cmd.Flags().IntVarP(&lines, "lines", "n", 80, "How many lines from the end")
	cmd.Flags().BoolVar(&errorsOnly, "errors", false, "Only the failures, with their context")
	cmd.Flags().StringVar(&grep, "grep", "", "Only lines containing this (case-insensitive)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Machine-readable output")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		root, err := scaffold.FindProjectRoot()
		if err != nil {
			return err
		}

		if errorsOnly {
			problems, err := devlog.Errors(root, 50)
			if err != nil || !devlog.Exists(root) {
				return noLogYet()
			}
			if asJSON {
				out, err := json.MarshalIndent(problems, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(out))
				return nil
			}
			if len(problems) == 0 {
				color.New(color.FgGreen).Println("  Nothing failed in this run.")
				color.New(color.FgHiBlack).Println("  Not every failure says the word error; " +
					"run `grit logs` to read the tail.")
				return nil
			}
			for _, problem := range problems {
				color.New(color.FgRed, color.Bold).Printf("%s\n", problem.Text)
				for _, more := range problem.More {
					color.New(color.FgHiBlack).Printf("    %s\n", more)
				}
				fmt.Println()
			}
			return nil
		}

		if !devlog.Exists(root) {
			return noLogYet()
		}
		tail, err := devlog.Tail(root, lines)
		if err != nil {
			return err
		}
		if grep != "" {
			needle := strings.ToLower(grep)
			kept := tail[:0]
			for _, line := range tail {
				if strings.Contains(strings.ToLower(line), needle) {
					kept = append(kept, line)
				}
			}
			tail = kept
		}

		if asJSON {
			out, err := json.MarshalIndent(tail, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(out))
			return nil
		}
		for _, line := range tail {
			fmt.Println(line)
		}
		return nil
	}

	return cmd
}

// noLogYet explains rather than failing with a missing file, because the fix
// is a command and not a bug.
func noLogYet() error {
	color.New(color.FgYellow).Println("  No development log yet.")
	fmt.Printf("  It is written by `grit start`, into %s/%s.\n", devlog.Dir, devlog.Name)
	fmt.Println("  Start the servers, reproduce what you are chasing, then run this again.")
	return nil
}
