package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
	"github.com/MUKE-coder/grit/v3/internal/testrunner"
)

func testCmd() *cobra.Command {
	var (
		onlyGo   bool
		onlyNode bool
		e2e      bool
		race     bool
		cover    bool
	)

	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run every test suite in the project",
		Long: "Runs the Go tests, the frontend tests, and — with --e2e — the Playwright suite,\n" +
			"then prints one report.\n\n" +
			"Which suites exist depends on the architecture you scaffolded, and you should\n" +
			"not have to remember which. Suites that do not apply are reported as skipped\n" +
			"with the reason, never quietly dropped.\n\n" +
			"Examples:\n" +
			"  grit test                 # Go + frontend\n" +
			"  grit test --go --race     # Go only, with the race detector\n" +
			"  grit test --e2e           # include Playwright (needs the app running)\n" +
			"  grit test --cover         # Go coverage",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := scaffold.FindProjectRoot()
			if err != nil {
				return err
			}

			suites := testrunner.Discover(root, testrunner.Options{
				Race:  race,
				Cover: cover,
				E2E:   e2e,
			})

			// Selection flags narrow the set. --e2e both opts the suite in and,
			// on its own, does NOT restrict to it: `grit test --e2e` should mean
			// "everything, including e2e", which is what people expect.
			var keys []string
			if onlyGo {
				keys = append(keys, "go")
			}
			if onlyNode {
				keys = append(keys, "node")
			}
			if len(keys) > 0 && e2e {
				keys = append(keys, "e2e")
			}
			suites = testrunner.Filter(suites, keys)

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Printf("\n  Running tests in %s\n", root)

			results := testrunner.Run(suites, os.Stdout, os.Stderr)
			printReport(results)

			if testrunner.Failed(results) {
				// The report above already named what failed; returning a bare
				// error here would print a redundant second line.
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&onlyGo, "go", false, "Run only the Go tests")
	cmd.Flags().BoolVar(&onlyNode, "node", false, "Run only the frontend tests")
	cmd.Flags().BoolVar(&e2e, "e2e", false, "Include the Playwright end-to-end suite")
	cmd.Flags().BoolVar(&race, "race", false, "Enable the Go race detector")
	cmd.Flags().BoolVar(&cover, "cover", false, "Report Go coverage")

	return cmd
}

func printReport(results []testrunner.Result) {
	var (
		pass, fail, skip int
		labelWidth       = len("RUNNER")
	)
	for _, r := range results {
		if len(r.Label) > labelWidth {
			labelWidth = len(r.Label)
		}
	}

	green := color.New(color.FgHiGreen, color.Bold)
	red := color.New(color.FgHiRed, color.Bold)
	dim := color.New(color.FgHiBlack)
	bold := color.New(color.Bold)

	fmt.Println()
	bold.Printf("  %-*s  %-6s  %s\n", labelWidth, "RUNNER", "STATUS", "TIME")
	fmt.Printf("  %s  %s  %s\n",
		strings.Repeat("─", labelWidth), strings.Repeat("─", 6), strings.Repeat("─", 8))

	for _, r := range results {
		switch r.Status {
		case testrunner.StatusPass:
			pass++
			fmt.Printf("  %-*s  ", labelWidth, r.Label)
			green.Printf("%-6s", "PASS")
			fmt.Printf("  %s\n", formatDuration(r.Duration))
		case testrunner.StatusFail:
			fail++
			fmt.Printf("  %-*s  ", labelWidth, r.Label)
			red.Printf("%-6s", "FAIL")
			fmt.Printf("  %s\n", formatDuration(r.Duration))
		case testrunner.StatusSkip:
			skip++
			fmt.Printf("  %-*s  ", labelWidth, r.Label)
			dim.Printf("%-6s", "SKIP")
			dim.Printf("  %s\n", r.SkipReason)
		}
	}

	fmt.Println()
	summary := fmt.Sprintf("  %d passed", pass)
	if fail > 0 {
		summary += fmt.Sprintf(" · %d failed", fail)
	}
	if skip > 0 {
		summary += fmt.Sprintf(" · %d skipped", skip)
	}

	if fail > 0 {
		red.Println(summary)
	} else {
		green.Println(summary)
	}
	fmt.Println()
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}
