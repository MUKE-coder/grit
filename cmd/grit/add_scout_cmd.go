package main

import (
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// addScoutCmd implements `grit add scout`.
//
// Opt-in, unlike Sentinel and Pulse, which every project gets. Two reasons.
// Scout is pre-1.0, and more importantly a security scanner is something you
// decide to point at an application; one that appears because you ran
// `grit new` is one nobody chose to run, and an unchosen scanner is either
// ignored or, worse, trusted without anybody having read what it does.
func addScoutCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "scout",
		Short: "Add Scout security scanning to the API",
		Long: "Installs Scout and mounts its dashboard at /scout.\n\n" +
			"Scout scans the running application: headers, TLS, DNS, exposed\n" +
			"files, Dockerfiles, secrets in the repository, reachable CVEs in\n" +
			"Go dependencies, and the route table.\n\n" +
			"The route table is the part Grit can give it that a scanner\n" +
			"usually cannot have. Gin does not expose a route's middleware\n" +
			"chain at runtime, so a scanner looking from outside has to guess\n" +
			"which endpoints are protected. Grit already parses the route\n" +
			"files into internal/access for its own Routes screen, and this\n" +
			"hands that table over, so \"this route has no auth gate\" becomes\n" +
			"something Scout reads rather than infers.\n\n" +
			"It mounts disabled. Set SCOUT_ENABLED=true and SCOUT_PASSWORD in\n" +
			"apps/api/.env to start it. The default mode is safe: it sends no\n" +
			"attack payloads. Raising the mode sends real ones, so raise it\n" +
			"only against a target you own.",
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			root, err := scaffold.FindProjectRoot()
			if err != nil {
				return err
			}
			cHeading.Println("\n  Adding Scout to apps/api")

			return scaffold.AddScout(root, force)
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Rewrite internal/routes/scout.go if it already exists")
	return cmd
}
