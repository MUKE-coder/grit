package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/explain"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// explainCmd is `grit explain <file>`: what this file is, who wrote it, and
// whether editing it is a good idea.
//
// The question it answers is the one a generated project provokes and does not
// answer: a few hundred files appear, and nothing says which are yours, which
// are the framework's, and which will be rewritten next time you upgrade.
// Reading a file tells you what it does; it does not tell you whether your
// change to it will survive.
func explainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "explain <file>",
		Short: "What a file is for, who wrote it, and whether to edit it",
		Long: "Explains one file in a generated project: what it is for, which command\n" +
			"wrote it, which Grit version, what it was generated from, and whether an\n" +
			"upgrade will rewrite it.\n\n" +
			"The editing advice comes from the manifest rather than from the path, so it\n" +
			"is about this project: a file you have already edited is reported as yours,\n" +
			"because that is what grit upgrade will treat it as.",
		Example: "  grit explain apps/api/internal/models/product.go\n" +
			"  grit explain packages/shared/types/product.ts\n" +
			"  grit explain .env",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := scaffold.FindProjectRoot()
			if err != nil {
				return err
			}

			e, err := explain.Explain(root, args[0])
			if err != nil {
				return err
			}

			bold := color.New(color.Bold)
			gray := color.New(color.FgHiBlack)
			cyan := color.New(color.FgHiCyan)
			yellow := color.New(color.FgHiYellow)
			green := color.New(color.FgHiGreen)

			fmt.Println()
			bold.Printf("  %s\n", e.Rel)

			if _, statErr := os.Stat(filepath.Join(root, filepath.FromSlash(e.Rel))); statErr != nil {
				yellow.Printf("  (not on disk)\n")
			}
			fmt.Println()

			if e.Purpose != "" {
				fmt.Printf("  %s\n\n", wrapAt(e.Purpose, 74, "  "))
			} else {
				gray.Printf("  No description for this path shape. The provenance below is still true.\n\n")
			}

			// Provenance, from the manifest.
			switch e.Kind {
			case explain.Generated:
				cyan.Printf("  %-12s", "Written by")
			case explain.Edited:
				yellow.Printf("  %-12s", "Written by")
			case explain.Gone:
				yellow.Printf("  %-12s", "Written by")
			default:
				green.Printf("  %-12s", "Written by")
			}
			if e.Generator != "" {
				fmt.Printf("%s", e.Generator)
				if e.Version != "" {
					gray.Printf("  (Grit %s)", e.Version)
				}
				fmt.Println()
			} else {
				fmt.Println("not Grit")
			}

			if e.Source != "" {
				gray.Printf("  %-12s%s\n", "Derived from", e.Source)
			}

			fmt.Println()
			fmt.Printf("  %s\n", wrapAt(e.Advice(), 74, "  "))

			if e.Edit != "" {
				fmt.Println()
				bold.Printf("  Edit it?  ")
				fmt.Printf("%s\n", wrapAt(e.Edit, 64, "            "))
			}

			if e.Read != "" {
				fmt.Println()
				gray.Printf("  Read  https://gritframework.dev%s\n", e.Read)
			}
			fmt.Println()
			return nil
		},
	}
	return cmd
}

// wrapAt breaks text at width, indenting every line after the first.
//
// Written here rather than pulled in, because the one thing it must not do is
// break a path or a command in half, and a general wrapper would.
func wrapAt(s string, width int, indent string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return s
	}
	var out string
	line := words[0]
	for _, w := range words[1:] {
		if len(line)+1+len(w) > width {
			out += line + "\n" + indent
			line = w
			continue
		}
		line += " " + w
	}
	return out + line
}

func splitWords(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r == ' ' || r == '\n' || r == '\t' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}
