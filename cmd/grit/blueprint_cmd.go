package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/blueprint"
)

func blueprintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "blueprint",
		Aliases: []string{"blueprints"},
		Short:   "Complete applications built with Grit",
		Long: "A blueprint is a whole application, in its own repository: not something\n" +
			"this command installs, but something you clone and work from.\n\n" +
			"A plugin adds a feature to your project. A blueprint is a project.",
	}
	cmd.AddCommand(blueprintListCmd())
	cmd.AddCommand(blueprintInfoCmd())
	return cmd
}

func blueprintListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List the available blueprints",
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			faint := color.New(color.Faint)

			fmt.Println("  Blueprints:")
			fmt.Println("  " + faint.Sprint("complete applications, each in its own repository") + "\n")

			planned := false
			for _, b := range blueprint.All() {
				if b.Access == blueprint.Planned && !planned {
					planned = true
					fmt.Println("  " + faint.Sprint("Announced, not built yet:"))
					fmt.Println()
				}

				fmt.Printf("  %-16s %s\n", color.CyanString(b.Name), blueprintAccessLabel(b.Access))
				fmt.Printf("  %-16s %s\n", "", faint.Sprint(b.Summary))
				if len(b.Surfaces) > 0 {
					fmt.Printf("  %-16s %s\n", "",
						faint.Sprint(strings.Join(b.SurfaceList(), " · ")))
				}
				fmt.Println()
			}

			fmt.Println("  grit blueprint info <name>   what it teaches, and where it is")
			fmt.Printf("  %s\n", faint.Sprintf(
				"$%d buys one. $%d a month, $%d a year or $%d once buys every one of them.",
				blueprint.PriceSingleUSD, blueprint.PriceMembershipMonthlyUSD,
				blueprint.PriceMembershipYearlyUSD, blueprint.PriceLifetimeUSD))
			fmt.Println("  " + faint.Sprint(blueprint.Site))
			fmt.Println()
			return nil
		},
	}
}

func blueprintInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info <name>",
		Short: "Show what a blueprint teaches and where to get it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := strings.ToLower(strings.TrimSpace(args[0]))
			b, ok := blueprint.Get(name)
			if !ok {
				return fmt.Errorf("unknown blueprint %q\n\nAvailable: %s",
					args[0], strings.Join(blueprint.Names(), ", "))
			}

			printLogo()
			faint := color.New(color.Faint)

			fmt.Printf("  %s  %s\n", color.New(color.Bold).Sprint(b.Title), blueprintAccessLabel(b.Access))
			fmt.Printf("  %s\n\n", faint.Sprint(b.Summary))

			fmt.Printf("  %-14s %s\n", "Architecture", b.Architecture)
			if len(b.Surfaces) > 0 {
				fmt.Printf("  %-14s %s\n", "Surfaces", strings.Join(b.SurfaceList(), ", "))
			}
			fmt.Println()

			if len(b.Teaches) > 0 {
				fmt.Println("  What you open it for:")
				for _, t := range b.Teaches {
					fmt.Printf("    • %s\n", t)
				}
				fmt.Println()
			}

			switch b.Access {
			case blueprint.Open:
				fmt.Printf("  %s\n", b.Repo)
				fmt.Printf("  %s\n", faint.Sprintf("git clone %s", b.Repo))
				fmt.Printf("  %s\n", faint.Sprintf("%s has the steps to run it", b.Page()))
			case blueprint.Paid:
				fmt.Printf("  %s\n", faint.Sprintf(
					"$%d once, or $%d a month for every blueprint. Paying invites your",
					blueprint.PriceSingleUSD, blueprint.PriceMembershipMonthlyUSD))
				fmt.Printf("  %s\n", faint.Sprint("GitHub account to the private repository."))
				fmt.Printf("  %s\n", b.Page())
			case blueprint.Planned:
				fmt.Printf("  %s\n", faint.Sprint("Not built yet. It is on the list, which is why it is listed."))
				fmt.Printf("  %s\n", blueprint.Site)
			}
			fmt.Println()
			return nil
		},
	}
}

func blueprintAccessLabel(a blueprint.Access) string {
	switch a {
	case blueprint.Open:
		return color.GreenString("open source")
	case blueprint.Paid:
		return color.YellowString("membership")
	default:
		return color.New(color.Faint).Sprint("planned")
	}
}
