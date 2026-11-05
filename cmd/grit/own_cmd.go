package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// `grit own <path>`: whose file is this, and what will an upgrade do to it.
//
// The answer existed in two places and was askable from neither. The manifest
// holds it, the upgrade guard acts on it, and an agent can read it over MCP
// through grit_file_ownership. A person at a terminal read the manifest by hand
// or inferred it from the directory, and inferring it is how people get this
// wrong (grit#161).
//
// It is also the question behind the two most common support exchanges: somebody
// edits a file and an upgrade leaves it alone, or edits one and an upgrade
// replaces it, and either way they want to know why.
func ownCmd() *cobra.Command {
	var showAll bool

	cmd := &cobra.Command{
		Use:   "own [path...]",
		Short: "Say whether a file is yours or the framework's, and what an upgrade will do to it",
		Long: `Ownership is the idea the whole project rests on, and the one people get wrong
first. This answers it for a file, from the manifest Grit writes when it
generates, which is the same record ` + "`grit upgrade`" + ` consults.

  grit own apps/admin/components/tables/data-table.tsx
  grit own apps/api/internal/handlers/*.go
  grit own --all

Four answers:

  framework-owned, unchanged     an upgrade replaces it
  framework-owned, edited        an upgrade leaves it alone and names it
  yours                          Grit has no record of writing it
  deleted                        an upgrade writes it back, unless its
                                 directory went with it`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := findProjectRoot()
			if err != nil {
				return fmt.Errorf("not inside a Grit project: %w", err)
			}

			if showAll {
				return ownAll(root)
			}
			if len(args) == 0 {
				return fmt.Errorf("name a file, or pass --all for every file Grit has written")
			}

			for _, arg := range args {
				if err := ownOne(root, arg); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&showAll, "all", false,
		"Every file in the manifest, grouped by what an upgrade will do")
	return cmd
}

func ownOne(root, arg string) error {
	o, err := scaffold.OwnershipOf(root, arg)
	if err != nil {
		return err
	}

	headline := color.New(color.FgHiWhite, color.Bold)
	dim := color.New(color.FgHiBlack)

	// Colour carries the same thing the words do, never more: green for a file
	// an upgrade replaces, yellow for one it holds back, cyan for yours.
	var tone *color.Color
	switch {
	case o.Status == manifest.Unchanged:
		tone = color.New(color.FgHiGreen)
	case o.Status == manifest.Modified, o.Status == manifest.Missing:
		tone = color.New(color.FgHiYellow)
	default:
		tone = color.New(color.FgHiCyan)
	}

	fmt.Fprintln(color.Output)
	headline.Fprintf(color.Output, "  %s\n", o.Rel)
	tone.Fprintf(color.Output, "  %s\n", o.Headline)
	for _, line := range strings.Split(o.Detail, "\n") {
		dim.Fprintf(color.Output, "    %s\n", line)
	}
	fmt.Fprintln(color.Output)
	return nil
}

func ownAll(root string) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	if !manifest.Exists(root) {
		fmt.Fprintln(color.Output)
		color.New(color.FgHiYellow).Fprintln(color.Output,
			"  This project has no manifest, so there is no record of what Grit wrote.")
		color.New(color.FgHiBlack).Fprintln(color.Output,
			"    Generated before v3.147.0. An upgrade overwrites everything it knows how\n"+
				"    to write, so commit first and read the diff afterwards.")
		fmt.Fprintln(color.Output)
		return nil
	}

	groups := map[string][]string{}
	for _, rel := range m.Entries() {
		abs := filepath.Join(root, filepath.FromSlash(rel))
		o, err := scaffold.OwnershipOf(root, abs)
		if err != nil {
			continue
		}
		groups[o.Headline] = append(groups[o.Headline], rel)
	}

	headings := make([]string, 0, len(groups))
	for h := range groups {
		headings = append(headings, h)
	}
	sort.Strings(headings)

	fmt.Fprintln(color.Output)
	for _, h := range headings {
		paths := groups[h]
		sort.Strings(paths)
		color.New(color.FgHiWhite, color.Bold).Fprintf(color.Output,
			"  %s (%d)\n", h, len(paths))
		for _, p := range paths {
			color.New(color.FgHiBlack).Fprintf(color.Output, "    %s\n", p)
		}
		fmt.Fprintln(color.Output)
	}

	color.New(color.FgHiBlack).Fprintln(color.Output,
		"  Files you wrote yourself are not listed: the manifest records what Grit\n"+
			"  wrote, so anything absent from it is yours by definition. Ask about one\n"+
			"  directly to have that said out loud.")
	fmt.Fprintln(color.Output)
	return nil
}
