package main

import (
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// addOfflineCmd is `grit add offline`.
//
// The API has served /api/sync/pull and /api/sync/push since v3.60, and the
// resource generator has registered every model with the sync registry since
// then too. What was missing was a client outside apps/desktop, so this
// installs one that apps/web, apps/admin and apps/expo can all import.
func addOfflineCmd() *cobra.Command {
	var models string

	cmd := &cobra.Command{
		Use:   "offline",
		Short: "Add offline-first sync to the web, admin and mobile apps",
		Long: "Installs packages/sync: a local mirror, an outbox for changes made while\n" +
			"offline, and optimistic-lock conflict handling, over a storage interface\n" +
			"with IndexedDB, expo-sqlite and in-memory implementations.\n\n" +
			"By default it mirrors every model the API has registered for sync. Pass\n" +
			"--models to narrow it.",
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			purple := color.New(color.FgHiMagenta, color.Bold)
			purple.Printf("\n  Adding offline-first sync\n\n")

			var wanted []string
			for _, m := range strings.Split(models, ",") {
				if trimmed := strings.TrimSpace(m); trimmed != "" {
					wanted = append(wanted, trimmed)
				}
			}

			root, err := scaffold.FindProjectRoot()
			if err != nil {
				return err
			}
			release, err := manifest.Start(root, version, "offline")
			if err != nil {
				return err
			}
			defer func() { _ = release() }()

			return scaffold.AddOffline(scaffold.AddOfflineOptions{
				Models:  wanted,
				Version: version,
			})
		},
	}

	cmd.Flags().StringVar(&models, "models", "",
		"Comma-separated models to mirror, in plural snake_case (default: all registered)")

	return cmd
}
