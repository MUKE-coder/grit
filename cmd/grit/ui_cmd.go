package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
	"github.com/MUKE-coder/grit/v3/internal/uiregistry"
)

func uiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ui",
		Short: "Browse and install Grit UI components",
		Long: "Grit UI is a registry of ready-made React components — marketing sections,\n" +
			"SaaS dashboards, ecommerce, auth flows and app layout.\n\n" +
			"Components are written into your repo as ordinary .tsx files you own. They\n" +
			"are shadcn registry items, so the same components work in any React project:\n\n" +
			"  npx shadcn@latest add " + uiregistry.DefaultBaseURL + "/r/hero-split-01.json\n\n" +
			"Browse them at " + uiregistry.DefaultBaseURL,
	}
	cmd.AddCommand(uiListCmd())
	cmd.AddCommand(uiAddCmd())
	return cmd
}

func uiListCmd() *cobra.Command {
	var category string

	cmd := &cobra.Command{
		Use:          "list",
		Short:        "List every component in the registry",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := uiregistry.List(cmd.Context())
			if err != nil {
				return err
			}

			want := strings.ToLower(strings.TrimSpace(category))
			filtered := items[:0:0]
			for _, it := range items {
				if want != "" && strings.ToLower(it.Category()) != want {
					continue
				}
				filtered = append(filtered, it)
			}

			if len(filtered) == 0 {
				if want != "" {
					return fmt.Errorf("no components in category %q\n\nAvailable: %s",
						category, strings.Join(categoryNames(items), ", "))
				}
				return fmt.Errorf("the registry returned no components")
			}

			nameWidth := len("NAME")
			for _, it := range filtered {
				if len(it.Name) > nameWidth {
					nameWidth = len(it.Name)
				}
			}

			bold := color.New(color.Bold)
			dim := color.New(color.FgHiBlack)

			fmt.Println()
			bold.Printf("  %-*s  %-11s %s\n", nameWidth, "NAME", "CATEGORY", "DESCRIPTION")
			fmt.Printf("  %s  %s %s\n",
				strings.Repeat("─", nameWidth), strings.Repeat("─", 11), strings.Repeat("─", 44))

			for _, it := range filtered {
				desc := it.Description
				if len(desc) > 62 {
					desc = desc[:59] + "..."
				}
				fmt.Printf("  %-*s  ", nameWidth, it.Name)
				dim.Printf("%-11s ", it.Category())
				fmt.Printf("%s\n", desc)
			}

			fmt.Println()
			dim.Printf("  %d component", len(filtered))
			if len(filtered) != 1 {
				dim.Print("s")
			}
			if want == "" {
				dim.Printf(" · categories: %s", strings.Join(categoryNames(items), ", "))
			}
			fmt.Println()
			dim.Printf("  Preview them at %s\n\n", uiregistry.BaseURL())
			return nil
		},
	}

	cmd.Flags().StringVar(&category, "category", "", "Filter by category (marketing, saas, ecommerce, auth, layout)")
	return cmd
}

func uiAddCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add <component>...",
		Short: "Install one or more components into your project",
		Long: "Downloads the component source and writes it to components/grit-ui/ in your\n" +
			"frontend app. The file is yours — edit it, rename it, delete it.\n\n" +
			"Run `grit ui list` to see what is available.",
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := scaffold.FindProjectRoot()
			if err != nil {
				return err
			}
			appDir, label, err := uiregistry.TargetDir(root)
			if err != nil {
				return err
			}

			green := color.New(color.FgHiGreen, color.Bold)
			dim := color.New(color.FgHiBlack)
			purple := color.New(color.FgHiMagenta, color.Bold)

			purple.Printf("\n  Installing into %s\n\n", label)

			deps := map[string]bool{}
			installed := 0

			for _, name := range args {
				component, err := uiregistry.Fetch(cmd.Context(), name)
				if err != nil {
					return err
				}
				path, err := uiregistry.Write(appDir, component, force)
				if err != nil {
					return err
				}
				for _, d := range component.Dependencies {
					deps[d] = true
				}

				green.Print("  ✓ ")
				fmt.Printf("%s", component.Title)
				dim.Printf("  components/grit-ui/%s.tsx\n", component.Name)
				_ = path
				installed++
			}

			fmt.Println()
			if len(deps) > 0 {
				names := make([]string, 0, len(deps))
				for d := range deps {
					names = append(names, d)
				}
				sort.Strings(names)
				// Named rather than installed: the package manager, workspace
				// layout and lockfile are the user's call, and a CLI that runs
				// an install in the wrong workspace is worse than one that tells
				// you what to run.
				dim.Printf("  Requires: %s\n", strings.Join(names, ", "))
				dim.Printf("  Install with: pnpm add %s\n\n", strings.Join(names, " "))
			}

			dim.Printf("  The components use Grit's design tokens. If this app was not\n")
			dim.Printf("  scaffolded by Grit, see %s/install\n\n", uiregistry.BaseURL())
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Overwrite a component file that already exists")
	return cmd
}

func categoryNames(items []uiregistry.Item) []string {
	seen := map[string]bool{}
	for _, it := range items {
		seen[it.Category()] = true
	}
	out := make([]string, 0, len(seen))
	for c := range seen {
		out = append(out, c)
	}
	sort.Strings(out)
	return out
}
