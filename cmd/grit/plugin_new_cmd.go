package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/plugin"
)

// grit#143: a plugin author started from an empty directory.
//
// The issue offered a template repository or this command and argued for this
// one, because a repository drifts from the manifest schema and a skeleton
// generated from templates inside the CLI cannot. The test in internal/plugin
// loads the output with LoadDir, which is the function `grit plugin add ./dir`
// calls, so the two cannot disagree without failing a build here.
func pluginNewCmd() *cobra.Command {
	var dir string
	var force bool

	cmd := &cobra.Command{
		Use:   "new <name>",
		Short: "Start a plugin of your own",
		Long: "Writes a working plugin directory: a manifest, a Go endpoint with a test, an\n" +
			"admin component behind a `when` clause, two injections, a README and a CI\n" +
			"workflow that installs the plugin into a fresh project and removes it again.\n\n" +
			"It is meant to be installed immediately, before you change anything, so the\n" +
			"loop is working before the plugin does anything interesting:\n\n" +
			"  grit plugin new audit-trail\n" +
			"  grit new testbed --triple\n" +
			"  cd testbed && grit plugin add ../audit-trail\n\n" +
			"Run it outside a project. A plugin is its own repository.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if err := plugin.ValidatePluginName(name); err != nil {
				return err
			}

			target := filepath.Join(dir, name)
			if dir == "" {
				target = name
			}

			// Checked before anything is written, because a half-written
			// plugin directory is worse than none: the manifest validates,
			// the files it names are missing, and the error names a file the
			// author never created.
			if existing, err := os.Stat(target); err == nil {
				if !force {
					kind := "a file"
					if existing.IsDir() {
						kind = "a directory"
					}
					return fmt.Errorf("%s is already %s\n\n"+
						"Pick another name, pass --dir to write it elsewhere, or --force to overwrite", target, kind)
				}
			}

			files := plugin.SkeletonFiles(name)
			paths := make([]string, 0, len(files))
			for p := range files {
				paths = append(paths, p)
			}
			sort.Strings(paths)

			for _, rel := range paths {
				full := filepath.Join(target, filepath.FromSlash(rel))
				if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
					return fmt.Errorf("creating %s: %w", filepath.Dir(full), err)
				}
				if err := os.WriteFile(full, []byte(files[rel]), 0o644); err != nil {
					return fmt.Errorf("writing %s: %w", full, err)
				}
			}

			printLogo()
			green := color.New(color.FgHiGreen)
			faint := color.New(color.Faint)

			fmt.Printf("  Plugin %s\n\n", color.CyanString(name))
			for _, rel := range paths {
				green.Printf("  ✓ %s\n", filepath.ToSlash(filepath.Join(name, rel)))
			}

			fmt.Println()
			fmt.Println("  It installs as it stands. Check that before you change it:")
			fmt.Println()
			fmt.Println("    " + color.CyanString("grit new testbed --triple"))
			fmt.Println("    " + color.CyanString("cd testbed"))
			fmt.Println("    " + color.CyanString("grit plugin add ../"+name))
			fmt.Println("    " + color.CyanString("cd apps/api && go build ./... && go test ./..."))
			fmt.Println("    " + color.CyanString("cd ../.. && grit plugin remove "+name+" --force"))
			fmt.Println("    " + color.CyanString("cd apps/api && go build ./..."))
			fmt.Println()
			fmt.Println("  That last build is the one worth running. Removal replays the")
			fmt.Println("  lockfile backwards, so it cannot take back an edit the lockfile")
			fmt.Println("  has no record of, and the build is where that shows up.")
			fmt.Println()
			fmt.Println("  " + faint.Sprint("The README in there covers the manifest, the markers you can inject"))
			fmt.Println("  " + faint.Sprint("at, and what not to do. The same workflow runs in CI on every push."))
			fmt.Println()
			fmt.Println("  Writing one: " + faint.Sprint("https://gritframework.dev/docs/plugins/authoring"))
			fmt.Println()
			return nil
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "", "Parent directory to create the plugin in (default: here)")
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite an existing directory of that name")
	return cmd
}
