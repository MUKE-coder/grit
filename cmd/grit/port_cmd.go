package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/port"
)

func portCmd() *cobra.Command {
	var to string
	var force bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "port <bundle|path> --to <app>",
		Short: "Copy a component out of the admin panel into another app",
		Long: `Copies an admin component, and everything it imports, into another app.

The admin has the interface work: a data table with sorting, filters, bulk
actions, export and import, a form builder, the dashboard widgets. A customer
portal usually wants some of that, and rebuilding it by hand is how a project
ends up with two tables that behave differently.

The files are copied, not shared. They land at the same path in the target app,
which is why they compile unchanged: every Grit frontend uses the same "@/*"
alias and the same design tokens. From then on they belong to that app and can
be edited there.

What it will not do is add packages to the target's package.json or invent a
file the admin does not have. Both are reported instead.

  grit port table --to apps/web
  grit port table-full --to apps/web
  grit port components/forms/form-builder.tsx --to apps/web
  grit port ui --to apps/desktop/frontend --dry-run`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := findProjectRoot()
			if err != nil {
				return err
			}
			adminRoot, err := adminAppRoot(root)
			if err != nil {
				return err
			}
			if to == "" {
				return fmt.Errorf("--to is required: the app to copy into, such as apps/web")
			}
			targetRoot := filepath.Join(root, filepath.FromSlash(to))
			if info, err := os.Stat(targetRoot); err != nil || !info.IsDir() {
				return fmt.Errorf("%s is not a directory in this project", to)
			}

			entries, summary := entriesFor(args[0])

			res, err := port.Run(port.Options{
				AdminRoot:  adminRoot,
				TargetRoot: targetRoot,
				Entries:    entries,
				Force:      force,
				DryRun:     dryRun,
			})
			if err != nil {
				return err
			}

			rel, _ := filepath.Rel(root, adminRoot)
			cHeading.Printf("\n  %s\n", summary)
			cMuted.Printf("  from %s to %s\n\n", filepath.ToSlash(rel), to)

			for _, f := range res.Copied {
				if dryRun {
					cMuted.Printf("  would write %s\n", f)
				} else {
					cOK.Printf("  ✓ %s\n", f)
				}
			}
			for _, f := range res.Skipped {
				cMuted.Printf("  · %s (already there, --force to replace)\n", f)
			}

			if len(res.Missing) > 0 {
				fmt.Println()
				noun := "packages"
				if len(res.Missing) == 1 {
					noun = "package"
				}
				cWarn.Printf("  %d %s not in %s/package.json\n", len(res.Missing), noun, to)
				cMuted.Println("  A workspace hoists these, so the app type-checks today and fails")
				cMuted.Println("  when it is built on its own. Add them yourself: this command does")
				cMuted.Println("  not write to your manifest.")
				for _, p := range res.Missing {
					fmt.Printf("    pnpm --filter %s add %s\n", filepath.Base(to), p)
				}
			} else if len(res.External) > 0 {
				fmt.Println()
				cMuted.Printf("  %d imported packages, all already in %s/package.json.\n",
					len(res.External), to)
			}

			if len(res.Unresolved) > 0 {
				fmt.Println()
				cWarn.Println("  Imports that could not be followed")
				cMuted.Println("  Each one is a file the admin does not have at that path. The copy")
				cMuted.Println("  will not compile until the target app provides it.")
				for _, u := range res.Unresolved {
					fmt.Printf("    %s\n", u)
				}
			}

			fmt.Println()
			if dryRun {
				cMuted.Printf("  %d files, nothing written.\n\n", len(res.Copied))
				return nil
			}
			cOK.Printf("  %d files copied.\n\n", len(res.Copied))
			cMuted.Println("  Next: type-check the target app, which is the only honest way to")
			cMuted.Println("  know this landed.")
			cMuted.Printf("    cd %s && pnpm type-check\n\n", to)
			return nil
		},
	}

	cmd.Flags().StringVar(&to, "to", "", "The app to copy into, such as apps/web (required)")
	cmd.Flags().BoolVar(&force, "force", false, "Replace files that are already in the target")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "List what would be copied, write nothing")

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "The bundles you can copy by name",
		RunE: func(cmd *cobra.Command, args []string) error {
			cHeading.Println("\n  Bundles")
			for _, b := range port.Bundles {
				fmt.Printf("    %-12s %s\n", b.Name, b.Summary)
			}
			fmt.Println()
			cMuted.Println("  Or give a path inside the admin, such as")
			cMuted.Println("  components/tables/table-filters.tsx")
			fmt.Println()
			return nil
		},
	})

	return cmd
}

// entriesFor turns the argument into entry files: a bundle name, or a path.
func entriesFor(arg string) ([]string, string) {
	for _, b := range port.Bundles {
		if strings.EqualFold(b.Name, arg) {
			return b.Entries, b.Summary
		}
	}
	return []string{arg}, "Copying " + arg
}

// adminAppRoot is the directory holding the admin's own code.
//
// Three shapes, because the admin is not always an app. A triple gives it
// apps/admin. A double puts it inside the web app as a route group, with its
// components under admin-panel/. A single keeps it as a section of the one
// frontend. Porting out of the second two is usually a no-op, and Run says so
// rather than copying a file onto itself.
func adminAppRoot(root string) (string, error) {
	candidates := []string{
		filepath.Join(root, "apps", "admin"),
		filepath.Join(root, "apps", "web", "admin-panel"),
		filepath.Join(root, "src", "admin-panel"),
		filepath.Join(root, "admin-panel"),
		filepath.Join(root, "frontend", "src", "admin-panel"),
	}
	for _, c := range candidates {
		if info, err := os.Stat(filepath.Join(c, "components")); err == nil && info.IsDir() {
			return c, nil
		}
	}
	return "", fmt.Errorf("no admin panel found in this project\n\n" +
		"Looked for apps/admin, apps/web/admin-panel, src/admin-panel and a few\n" +
		"others. There is nothing to copy from until a project has one.")
}
