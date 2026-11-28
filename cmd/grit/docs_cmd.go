package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/docs"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// docsCmd is `grit docs`: the documentation, searchable offline and pinned to
// the version of the CLI that is installed.
//
// The reason it exists is version drift. Grit ships several releases a week, so
// a developer or an agent reading the docs site is reading a different version
// from the one their project is pinned to, and the symptom is code that calls a
// helper which did not exist yet or a flag that has since been renamed. This
// answers from an index embedded at build time, so what it says is true of the
// binary saying it, and it warns when the project disagrees.
func docsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Search the documentation for this exact version, offline",
		Long: "The documentation, embedded in the binary and pinned to its version.\n\n" +
			"Grit ships several releases a week, so the docs site is usually ahead of the\n" +
			"version a project is pinned to. This answers from an index built with this\n" +
			"binary, and says so when the project was scaffolded with a different one.",
	}

	cmd.AddCommand(docsSearchCmd(), docsShowCmd(), docsListCmd())
	return cmd
}

// wantsHelp reports whether the arguments are asking for the usage text.
//
// DisableFlagParsing turns off cobra's own handling of -h, so the command has
// to notice it. Without this, `grit docs search --help` searches for the word
// help.
func wantsHelp(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

// parseSearchArgs splits the arguments into this command's own options and
// the query, and takes the dashes off anything flag-shaped.
//
// The dashes come off because the index cannot hold them: its keywords are
// split on every non-alphanumeric when they are built, so --owned-by is
// stored as owned-by with nothing in front. Stripping here is what makes the
// two spellings meet.
//
// Only the leading dashes, and only when a word is left: a bare "-" stays as
// it is rather than becoming a search for every page containing a hyphen.
func parseSearchArgs(args []string, limit int, asJSON bool) (string, int, bool, error) {
	terms := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--json":
			asJSON = true
		case arg == "--limit":
			if i+1 >= len(args) {
				return "", 0, false, fmt.Errorf("--limit needs a number, e.g. --limit 3")
			}
			i++
			n, err := strconv.Atoi(args[i])
			if err != nil {
				return "", 0, false, fmt.Errorf("--limit takes a number, not %q", args[i])
			}
			limit = n
		case strings.HasPrefix(arg, "--limit="):
			n, err := strconv.Atoi(strings.TrimPrefix(arg, "--limit="))
			if err != nil {
				return "", 0, false, fmt.Errorf("--limit takes a number, not %q", strings.TrimPrefix(arg, "--limit="))
			}
			limit = n
		default:
			if bare := strings.TrimLeft(arg, "-"); bare != "" {
				terms = append(terms, bare)
			} else {
				terms = append(terms, arg)
			}
		}
	}
	return strings.Join(terms, " "), limit, asJSON, nil
}

func docsSearchCmd() *cobra.Command {
	var limit int
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Find the pages that answer a question",
		Example: "  grit docs search field types\n" +
			"  grit docs search --owned-by\n" +
			"  grit docs search owned by user --limit 3",
		// The query may itself be a flag, and that is the whole point.
		//
		// "What does --owned-by do" is asked by typing --owned-by, and cobra
		// answered "unknown flag" and exited 1: the one search most likely to
		// be a flag was the one search that could not be run. So this command
		// parses its own two options and treats everything else, dashes and
		// all, as the query.
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if wantsHelp(args) {
				return cmd.Help()
			}
			query, limit, asJSON, err := parseSearchArgs(args, limit, asJSON)
			if err != nil {
				return err
			}
			if query == "" {
				return fmt.Errorf("what should I search for? e.g. grit docs search field types")
			}
			hits, err := docs.Search(query, limit)
			if err != nil {
				return err
			}

			if asJSON {
				out, err := json.MarshalIndent(map[string]interface{}{
					"version": docs.Version(),
					"query":   query,
					"hits":    hits,
				}, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(out))
				return nil
			}

			if len(hits) == 0 {
				fmt.Printf("Nothing in the v%s documentation matches %q.\n", docs.Version(), query)
				fmt.Println("Try fewer words, or `grit docs list` to see what there is.")
				return nil
			}

			warnAboutVersionDrift(cmd)
			for _, hit := range hits {
				color.New(color.FgCyan, color.Bold).Printf("%s\n", hit.Title)
				fmt.Printf("  %s\n", hit.URL)
				if hit.Snippet != "" {
					fmt.Printf("  %s\n", hit.Snippet)
				}
				fmt.Println()
			}
			fmt.Printf("Read one with: grit docs show %s\n", hits[0].URL)
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 5, "How many results")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Machine-readable output")
	return cmd
}

func docsShowCmd() *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:     "show <url>",
		Short:   "Print one page",
		Args:    cobra.ExactArgs(1),
		Example: "  grit docs show /docs/concepts/field-types",
		RunE: func(cmd *cobra.Command, args []string) error {
			page, err := docs.Read(args[0])
			if err != nil {
				return err
			}

			if asJSON {
				out, err := json.MarshalIndent(page, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(out))
				return nil
			}

			warnAboutVersionDrift(cmd)
			color.New(color.FgCyan, color.Bold).Printf("%s\n", page.Title)
			fmt.Printf("%s  (v%s)\n\n", page.URL, docs.Version())
			if page.Description != "" {
				fmt.Printf("%s\n\n", page.Description)
			}
			if len(page.Headings) > 0 {
				heads := page.Headings
				if len(heads) > 25 {
					heads = heads[:25]
				}
				color.New(color.Faint).Println("Sections")
				for _, h := range heads {
					fmt.Printf("  - %s\n", h)
				}
				fmt.Println()
			}
			if page.Text != "" {
				fmt.Println(page.Text)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&asJSON, "json", false, "Machine-readable output")
	return cmd
}

func docsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List every page in the index",
		RunE: func(cmd *cobra.Command, args []string) error {
			urls := docs.URLs()
			fmt.Printf("%d pages, documentation for v%s\n\n", len(urls), docs.Version())
			for _, url := range urls {
				fmt.Println(url)
			}
			return nil
		},
	}
}

// warnAboutVersionDrift says so when the project was scaffolded with a
// different version from the one answering.
//
// Printed rather than refused: older documentation is still mostly right, and
// a developer who cannot read it at all is worse off than one who knows to
// check. It goes to stderr so `--json` consumers are unaffected and a human
// still sees it.
func warnAboutVersionDrift(cmd *cobra.Command) {
	root, err := scaffold.FindProjectRoot()
	if err != nil {
		// Not in a project, so there is nothing to disagree with.
		return
	}
	projectVersion := scaffold.ProjectVersion(root)
	if projectVersion == "" || projectVersion == docs.Version() {
		return
	}
	color.New(color.FgYellow).Fprintf(cmd.ErrOrStderr(),
		"These are the v%s docs and this project was scaffolded with v%s.\n"+
			"Anything new since then will not apply. Run `grit update` to match.\n\n",
		docs.Version(), projectVersion)
}
