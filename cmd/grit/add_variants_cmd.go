package main

import (
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/generate"
)

// addVariantsCmd is `grit add variants --resource Product`.
//
// A command rather than a flag on `generate resource` because two of the five
// tables it installs are shared. Options and their values belong to the shop:
// Colour is Colour whether it is on a shirt or a phone case, and a per-resource
// copy of them gives you fourteen spellings of it and a filter that can only
// match one.
func addVariantsCmd() *cobra.Command {
	var resource string

	cmd := &cobra.Command{
		Use:   "variants",
		Short: "Add product variants: options, values, and a combination matrix",
		Long: "Installs the normalised variant system and attaches it to one resource.\n\n" +
			"Five tables: options and their values (shared across the shop), which\n" +
			"options a given resource offers, the buyable combinations, and the join\n" +
			"between them.\n\n" +
			"A variant's price is resolved rather than stored: an override wins, and\n" +
			"otherwise it is the resource's price plus the deltas of the values whose\n" +
			"option declares affects_price. Storing it instead means a price change\n" +
			"leaves every variant on the old figure.",
		Example: "  grit add variants --resource Product",
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()
			cmd.SilenceUsage = true
			return generate.AddVariants(resource)
		},
	}

	cmd.Flags().StringVar(&resource, "resource", "Product",
		"The resource that offers variants")

	return cmd
}
