package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/swap"
	"github.com/spf13/cobra"
)

func swapCmd() *cobra.Command {
	var (
		swapForce     bool
		swapSkipCheck bool
		swapList      bool
		swapRevert    bool
	)

	cmd := &cobra.Command{
		Use:   "swap [slot] [variant]",
		Short: "Replace an admin component everywhere at once",
		Long: `Replace a canonical admin component — a "slot" — with a variant from Grit UI.

This is not the same as ` + "`grit ui add`" + `. Adding gives you a new file to import
where you like. Swapping overwrites the one file every call site already imports,
so one command restyles every button in the admin without touching an import.

  grit swap --list                 what can be swapped, and to what
  grit swap button glow-ring       swap the button slot
  grit swap button --revert        put the previous one back

The previous file is always backed up to .grit/swaps/ first, and the admin is
type-checked afterwards — a variant that does not compile against your call sites
is rolled back rather than left in place.`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			if swapList {
				return runSwapList(ctx, args)
			}

			root, err := os.Getwd()
			if err != nil {
				return err
			}
			project, err := swap.FindProject(root)
			if err != nil {
				return err
			}

			if len(args) == 0 {
				return runSwapStatus(project)
			}
			slot := args[0]
			if !contains(swap.KnownSlots, slot) {
				return fmt.Errorf("unknown slot %q\n\nSwappable slots: %s",
					slot, strings.Join(swap.KnownSlots, ", "))
			}

			if swapRevert {
				name, err := swap.Revert(project, slot)
				if err != nil {
					return err
				}
				fmt.Printf("  Reverted %s from %s\n", slot, name)
				fmt.Printf("  %s\n", rel(project.Root, project.SlotPath(slot)))
				return nil
			}

			if len(args) < 2 {
				return fmt.Errorf("which variant?\n\n  grit swap %s <variant>\n\nRun `grit swap --list %s` to see them", slot, slot)
			}
			variant := args[1]

			fmt.Printf("  Swapping %s → %s\n", slot, variant)
			if !swapSkipCheck {
				fmt.Printf("  Type-checking after the write (--skip-check to skip)...\n")
			}

			res, err := swap.Apply(ctx, project, slot, variant, swap.Options{
				Force:     swapForce,
				SkipCheck: swapSkipCheck,
			})
			if err != nil {
				return err
			}

			fmt.Printf("\n  ✅ %s is now %s (%s)\n", res.Slot, res.Variant, res.Contract)
			fmt.Printf("     %s\n", rel(project.Root, res.Path))
			fmt.Printf("     backup: %s\n", rel(project.Root, res.BackupPath))
			if res.Checked {
				fmt.Printf("     type-check passed\n")
			}
			fmt.Printf("\n  Undo with: grit swap %s --revert\n", res.Slot)
			return nil
		},
	}

	cmd.Flags().BoolVar(&swapList, "list", false, "List swappable slots and variants")
	cmd.Flags().BoolVar(&swapRevert, "revert", false, "Restore the slot's most recent backup")
	cmd.Flags().BoolVar(&swapForce, "force", false, "Overwrite a slot file you have edited")
	cmd.Flags().BoolVar(&swapSkipCheck, "skip-check", false, "Skip the post-swap type-check")
	// Cobra otherwise prints the whole usage block under a real error, which
	// buries the one line that actually says what went wrong.
	cmd.SilenceUsage = true
	return cmd
}

func runSwapList(ctx context.Context, args []string) error {
	slot := ""
	if len(args) > 0 {
		slot = args[0]
	}
	variants, err := swap.ListVariants(ctx, slot)
	if err != nil {
		return err
	}
	if len(variants) == 0 {
		if slot != "" {
			return fmt.Errorf("no variants published for slot %q", slot)
		}
		return fmt.Errorf("no swappable variants found in the registry")
	}

	fmt.Println()
	current := ""
	for _, v := range variants {
		if v.Slot != current {
			current = v.Slot
			fmt.Printf("  %s\n", strings.ToUpper(v.Slot))
		}
		tag := ""
		if v.Pro {
			tag = "  [pro]"
		}
		fmt.Printf("    %-18s %s%s\n", v.Name, v.Title, tag)
		if v.Description != "" {
			fmt.Printf("    %-18s %s\n", "", truncate(v.Description, 78))
		}
	}
	fmt.Printf("\n  grit swap <slot> <variant>\n\n")
	return nil
}

func runSwapStatus(p *swap.Project) error {
	state := p.LoadState()
	fmt.Printf("\n  Slots in %s\n\n", p.Label)
	for _, slot := range swap.KnownSlots {
		path := p.SlotPath(slot)
		if _, err := os.Stat(path); err != nil {
			fmt.Printf("    %-10s (not present)\n", slot)
			continue
		}
		if rec, ok := state[slot]; ok {
			fmt.Printf("    %-10s %s  (%s, swapped %s)\n", slot, rec.Variant, rec.Contract, rec.SwappedAt)
			continue
		}
		// Never swapped, so it is whatever the scaffolder wrote. Read the
		// contract out of the file rather than assuming — the project may
		// predate the current CLI.
		b, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("    %-10s (unreadable)\n", slot)
			continue
		}
		fmt.Printf("    %-10s default  (%s)\n", slot, swap.ContractOf(string(b)))
	}
	fmt.Printf("\n  grit swap --list   to see what else is available\n\n")
	return nil
}

func rel(root, path string) string {
	if r, err := filepath.Rel(root, path); err == nil {
		return filepath.ToSlash(r)
	}
	return filepath.ToSlash(path)
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
