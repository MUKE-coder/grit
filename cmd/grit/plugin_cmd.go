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
			fmt.Println("  grit plugin info <name>   details")
			fmt.Println("  grit plugin add <name>    install")
			fmt.Println()
			return nil
		},
	}
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
		Use:   "add <name>",
		Short: "Install a plugin into this project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			printLogo()

			p, err := plugin.Get(args[0])
			if err != nil {
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
