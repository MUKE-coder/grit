package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/plugin"
	"github.com/MUKE-coder/grit/v3/internal/project"
)

func pluginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Install and manage Grit plugins",
		Long: "Plugins extend a Grit project with generated code — models, routes, pages and\n" +
			"migrations written into your repo, not hidden behind a runtime dependency.\n\n" +
			"Everything an install does is recorded in .grit/plugins.lock.json, and\n" +
			"`grit plugin remove` replays that record backwards. Commit the lockfile.",
	}

	cmd.AddCommand(pluginListCmd())
	cmd.AddCommand(pluginInfoCmd())
	cmd.AddCommand(pluginAddCmd())
	cmd.AddCommand(pluginUpdateCmd())
	cmd.AddCommand(pluginRemoveCmd())

	return cmd
}

// findProjectRoot walks up looking for grit.json.
//
// grit.json is the canonical marker every `grit new` project gets, whatever its
// architecture. project.DetectProject keys off turbo.json/wails.json instead,
// which an --api project has neither of — so it can't be used here.
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "grit.json")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no grit.json found in this directory or any parent")
		}
		dir = parent
	}
}

// pluginContext builds the Context a plugin needs from the detected project.
func pluginContext() (plugin.Context, error) {
	root, err := findProjectRoot()
	if err != nil {
		return plugin.Context{}, fmt.Errorf("not inside a Grit project: %w", err)
	}

	ctx := plugin.Context{
		Root:         root,
		Architecture: "triple",
		Frontend:     "next",
	}

	// Module path comes from the API's go.mod — plugins emit imports with it.
	if info, err := project.DetectProjectFrom(root); err == nil && info.Module != "" {
		ctx.Module = info.Module
	}

	// grit.json is the source of truth for shape; fall back to sane defaults so
	// a hand-made project still works.
	if data, err := os.ReadFile(filepath.Join(root, "grit.json")); err == nil {
		var meta struct {
			Architecture string `json:"architecture"`
			Frontend     string `json:"frontend"`
		}
		if json.Unmarshal(data, &meta) == nil {
			if meta.Architecture != "" {
				ctx.Architecture = meta.Architecture
			}
			if meta.Frontend != "" {
				ctx.Frontend = meta.Frontend
			}
		}
	}

	// Single-binary apps keep Go code at the root; monorepos use apps/api.
	if ctx.Architecture == "single" {
		ctx.APIRoot = root
	} else {
		ctx.APIRoot = filepath.Join(root, "apps", "api")
	}

	// Fall back to reading the module from go.mod directly, since
	// DetectProjectFrom only recognises monorepo/desktop layouts.
	if ctx.Module == "" {
		ctx.Module = readGoModule(filepath.Join(ctx.APIRoot, "go.mod"))
	}
	return ctx, nil
}

// readGoModule returns the module path from a go.mod, or "" if unreadable.
func readGoModule(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}

func pluginListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			// Installed state is only knowable inside a project; listing what's
			// available should still work anywhere.
			installed := map[string]string{}
			if ctx, err := pluginContext(); err == nil {
				if lock, err := plugin.LoadLock(ctx.Root); err == nil {
					for _, p := range lock.Plugins {
						installed[p.Name] = p.Version
					}
				}
			}

			all := plugin.All()
			if len(all) == 0 {
				fmt.Println("  No plugins are available yet.")
				return nil
			}

			fmt.Println("  Available plugins:")
			fmt.Println()
			for _, p := range all {
				status := color.New(color.Faint).Sprint("not installed")
				if v, ok := installed[p.Name]; ok {
					status = color.GreenString("installed (v%s)", v)
				}
				fmt.Printf("  %-16s %s\n", color.CyanString(p.Name), status)
				fmt.Printf("  %-16s %s\n", "", color.New(color.Faint).Sprint(p.Summary))
				fmt.Println()
			}
			printModulePlugins()

			fmt.Println("  grit plugin info <name>   details")
			fmt.Println("  grit plugin add <name>    install")
			fmt.Println()
			return nil
		},
	}
}

// modulePlugins are the packages at github.com/MUKE-coder/grit-plugins.
//
// They are installed with go get rather than grit plugin add, so they cannot
// be folded into the list above without implying a command that does not work.
// Listing them anyway matters more than the tidiness of one list: without it
// the CLI reads as a complete catalogue of five, and the nine it does not
// mention are the ones people go looking for first.
var modulePlugins = []struct{ name, summary string }{
	{"websockets", "Realtime hub with rooms, presence and typing indicators"},
	{"notifications", "In-app and push notifications with a delivery log"},
	{"search", "Full-text search across resources"},
	{"stripe", "Subscriptions and billing (one-off checkout is grit plugin add stripe)"},
	{"oauth", "Sign in with Google, GitHub and friends"},
	{"i18n", "Translated API messages and locale negotiation"},
	{"video", "Video upload, transcode and playback"},
	{"export", "Scheduled exports to CSV and XLSX"},
	{"conference", "Audio and video rooms"},
}

func printModulePlugins() {
	fmt.Println("  Go module plugins:")
	fmt.Println("  " + color.New(color.Faint).Sprint(
		"installed with go get, from github.com/MUKE-coder/grit-plugins") + "\n")
	for _, p := range modulePlugins {
		fmt.Printf("  %-16s %s\n", color.CyanString(p.name),
			color.New(color.Faint).Sprint(p.summary))
	}
	fmt.Println()
	fmt.Println("  " + color.New(color.Faint).Sprint(
		"go get github.com/MUKE-coder/grit-plugins/grit-<name>"))
	fmt.Println()
}

func pluginInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info <name>",
		Short: "Show what a plugin does",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()
			p, err := plugin.Get(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("  %s  %s\n\n", color.CyanString(p.Name), color.New(color.Faint).Sprintf("v%s", p.Version))
			fmt.Printf("  %s\n\n", p.Summary)
			if p.Description != "" {
				for _, line := range strings.Split(strings.TrimSpace(p.Description), "\n") {
					fmt.Printf("  %s\n", line)
				}
				fmt.Println()
			}
			if len(p.Requires) > 0 {
				fmt.Printf("  Requires: %s\n\n", strings.Join(p.Requires, ", "))
			}
			if len(p.GoDeps) > 0 {
				fmt.Println("  Go dependencies:")
				for _, d := range p.GoDeps {
					fmt.Printf("    %s %s\n", d.Name, d.Version)
				}
				fmt.Println()
			}
			return nil
		},
	}
}

func pluginAddCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add <name|./path>",
		Short: "Install a plugin into this project",
		Long: `Install a built-in plugin by name, or one of your own from a directory:

  grit plugin add impersonate
  grit plugin add ./plugins/product-reviews

A plugin directory holds a plugin.json manifest and the files it installs.
Everything else is identical: the same lockfile records what was written, and
` + "`" + `grit plugin remove` + "`" + ` replays it backwards either way.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			// An explicit path means a plugin of your own. Requiring ./ or an
			// absolute path rather than sniffing for a directory keeps a local
			// folder from ever shadowing a built-in that shares its name.
			var p plugin.Plugin
			var err error
			if plugin.IsDirRef(args[0]) {
				p, err = plugin.LoadDir(args[0])
			} else {
				p, err = plugin.Get(args[0])
			}
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}
			ctx, err := pluginContext()
			if err != nil {
				return err
			}

			fmt.Printf("  Installing %s v%s\n\n", color.CyanString(p.Name), p.Version)

			// A plugin writes code and edits files. Uncommitted work makes that
			// hard to review or undo, so say so before touching anything.
			if !force && isDirtyGitTree(ctx.Root) {
				fmt.Println(color.YellowString("  ⚠ You have uncommitted changes."))
				fmt.Println("    Installing writes files and patches existing ones; committing first")
				fmt.Println("    makes the change reviewable with `git diff`.")
				fmt.Println("    Re-run with --force to install anyway.")
				fmt.Println()
				return fmt.Errorf("working tree is dirty")
			}

			record, err := plugin.Install(ctx, p)
			if err != nil {
				return err
			}

			fmt.Println()
			fmt.Printf("  %s %s installed (%d files, %d patches)\n",
				color.GreenString("✅"), p.Name, len(record.Files), len(record.Injections))

			if len(p.GoDeps) > 0 {
				fmt.Println()
				fmt.Println("  Add the Go dependencies:")
				for _, d := range p.GoDeps {
					fmt.Printf("    go get %s@%s\n", d.Name, d.Version)
				}
			}
			if len(p.NextSteps) > 0 {
				fmt.Println()
				fmt.Println("  Next steps:")
				for i, s := range p.NextSteps {
					fmt.Printf("    %d. %s\n", i+1, s)
				}
			}
			fmt.Println()
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Install even with uncommitted changes")
	return cmd
}

func pluginRemoveCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:     "remove <name>",
		Aliases: []string{"rm"},
		Short:   "Remove an installed plugin",
		Long: "Deletes the files the plugin wrote and reverts the snippets it injected,\n" +
			"using the record in .grit/plugins.lock.json.\n\n" +
			"Anything you edited by hand is reported rather than overwritten — removal\n" +
			"never guesses at a block that no longer matches what was installed.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			ctx, err := pluginContext()
			if err != nil {
				return err
			}

			if !force {
				fmt.Printf("\n  ⚠ This removes every file and patch from %q.\n", args[0])
				if !confirm() {
					fmt.Println("\n  Cancelled.")
					return nil
				}
			}

			fmt.Println()
			warnings, err := plugin.Remove(ctx.Root, args[0])
			if err != nil {
				return err
			}

			fmt.Println()
			if len(warnings) > 0 {
				fmt.Println(color.YellowString("  Needs your attention:"))
				for _, w := range warnings {
					fmt.Printf("    • %s\n", w)
				}
				fmt.Println()
			}
			fmt.Printf("  %s %s removed\n\n", color.GreenString("✅"), args[0])
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Skip the confirmation prompt")
	return cmd
}

// confirm asks for a y/N before a destructive step.
func confirm() bool {
	fmt.Print("  Continue? [y/N]: ")
	var answer string
	// A read error means no usable stdin (CI, a pipe). Declining is the safe
	// default — silently proceeding would delete files nobody agreed to.
	if _, err := fmt.Scanln(&answer); err != nil {
		return false
	}
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes"
}

// isDirtyGitTree reports whether the project has uncommitted changes.
//
// Best-effort: if git is missing or this isn't a repo, report clean rather than
// blocking the install. The check is a courtesy so the user can review a
// plugin's edits with `git diff`, not a safety mechanism.
func isDirtyGitTree(root string) bool {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(out))) > 0
}

// pluginUpdateCmd exists because a plugin fix used to reach nobody: add refuses
// once a plugin is installed, and remove-then-add loses local edits and moves
// the injections, which has broken a build.
func pluginUpdateCmd() *cobra.Command {
	var force bool
	var all bool
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "update [name]",
		Short: "Bring installed plugins up to this CLI's version of them",
		Long: `Rewrite the files an installed plugin owns, and add the ones it has gained.

  grit plugin update stripe
  grit plugin update --all

Only files nobody has touched are rewritten: each one is fingerprinted when the
plugin writes it, so an edited file is reported and left exactly as it is.
Patches already in place stay where they sit, and only a new one is applied.
Nothing is ever deleted.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()
			cmd.SilenceUsage = true

			if len(args) == 0 && !all {
				return fmt.Errorf("name a plugin, or pass --all")
			}

			ctx, err := pluginContext()
			if err != nil {
				return err
			}

			lock, err := plugin.LoadLock(ctx.Root)
			if err != nil {
				return err
			}
			names := args
			if all {
				names = lock.Installed()
			}
			if len(names) == 0 {
				fmt.Println("  No plugins are installed.")
				fmt.Println()
				return nil
			}

			// Rewriting files is the kind of change you want to read in a
			// diff, exactly like an install.
			if !force && isDirtyGitTree(ctx.Root) {
				fmt.Println(color.YellowString("  ⚠ You have uncommitted changes."))
				fmt.Println("    Updating rewrites files the plugin owns; committing first")
				fmt.Println("    makes the change reviewable with `git diff`.")
				fmt.Println("    Re-run with --force to update anyway.")
				fmt.Println()
				return fmt.Errorf("working tree is dirty")
			}

			for _, name := range names {
				p, err := plugin.Get(name)
				if err != nil {
					// Almost always a plugin of your own, installed from a
					// directory: Grit carries no copy of it, so there is
					// nothing here to update it from. Not a reason to fail,
					// and not worth phrasing as "you typed it wrong".
					fmt.Printf("  %s %s is not one of Grit's own plugins, so there is nothing to update it from\n",
						color.YellowString("⚠"), name)
					continue
				}
				res, err := plugin.UpdateWith(ctx, p, plugin.UpdateOptions{Overwrite: overwrite})
				if err != nil {
					return err
				}
				printUpdate(res)
			}
			fmt.Println()
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Update even with uncommitted changes")
	cmd.Flags().BoolVar(&all, "all", false, "Update every installed plugin")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Take the plugin's version of every file, including ones you have edited")
	return cmd
}

func printUpdate(res *plugin.UpdateResult) {
	version := res.To
	if res.From != res.To {
		version = res.From + " → " + res.To
	}
	if res.Blocked {
		fmt.Printf("  %s %s (%s): nothing was changed\n", color.YellowString("⚠"), res.Name, version)
		fmt.Println("    This project installed the plugin before Grit recorded what it wrote,")
		fmt.Println("    so an edit of yours and a change in the plugin look the same here. Doing")
		fmt.Println("    half of an update would leave a project that does not build, so none of")
		fmt.Println("    it was done. These are the files it cannot vouch for:")
		for _, f := range res.Unverified {
			fmt.Printf("      %s\n", f)
		}
		color.New(color.FgHiCyan).Printf("    grit plugin update %s --overwrite   # take the plugin's version, then read git diff\n", res.Name)
		return
	}
	if !res.Changed() && len(res.Edited) == 0 && len(res.Unverified) == 0 {
		fmt.Printf("  %s %s (%s) is already up to date\n", color.GreenString("✓"), res.Name, version)
		return
	}

	fmt.Printf("  %s %s (%s)\n", color.CyanString("→"), res.Name, version)
	for _, f := range res.Added {
		fmt.Printf("    %s %s (new)\n", color.GreenString("+"), f)
	}
	for _, f := range res.Replaced {
		fmt.Printf("    %s %s\n", color.GreenString("✓"), f)
	}
	for _, f := range res.Injected {
		fmt.Printf("    %s patched %s\n", color.GreenString("✓"), f)
	}
	for _, d := range res.Deps {
		fmt.Printf("    %s %s (run pnpm install)\n", color.GreenString("✓"), d)
	}
	if len(res.Edited) > 0 {
		fmt.Println(color.YellowString("    Left alone, because you have edited them since the plugin wrote them:"))
		for _, f := range res.Edited {
			fmt.Printf("      %s\n", f)
		}
	}
	if len(res.Unverified) > 0 {
		fmt.Println(color.YellowString("    Left alone, because this project was set up before plugin files were fingerprinted,"))
		fmt.Println(color.YellowString("    so an edit and an improvement look the same here. Compare with a fresh project:"))
		for _, f := range res.Unverified {
			fmt.Printf("      %s\n", f)
		}
	}
	if len(res.Edited)+len(res.Unverified) > 0 {
		cyan := color.New(color.FgHiCyan)
		cyan.Printf("    grit plugin update %s --overwrite   # take the plugin's version, then read git diff\n", res.Name)
		if len(res.Added) > 0 {
			fmt.Println("    A file left behind is often the one the new file needs, so if the build")
			fmt.Println("    breaks, that is what to look at first.")
		}
	}
	if len(res.Orphaned) > 0 {
		fmt.Println("    No longer part of the plugin, and kept in case you use them:")
		for _, f := range res.Orphaned {
			fmt.Printf("      %s\n", f)
		}
	}
}
