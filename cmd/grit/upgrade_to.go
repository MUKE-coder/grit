package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/MUKE-coder/grit/v3/internal/selfupdate"
)

// `grit upgrade --to <version>`: apply one named version's framework files.
//
// Until this, grit upgrade had one speed. It applied the CLI's own templates,
// and the only way to pin a project's framework version was to pin the binary
// on the machine and in CI. With several releases a week and six hundred tags
// since v3.0.0, "latest" is a long way from wherever an older project is, and
// moving one step at a time was not possible (grit#158).
//
// It runs the other version's binary rather than trying to apply its templates
// from this one, because the templates ARE the binary: there is no archive of
// past templates to read, and pretending otherwise would mean applying this
// version's files under another version's name. The downloaded binary is
// checksum-verified against the release's signed manifest, cached under
// ~/.grit/versions, and never put on PATH.
func upgradeToVersion(to string, passthrough []string) error {
	to = strings.TrimPrefix(strings.TrimSpace(to), "v")
	if to == "" {
		return fmt.Errorf("--to needs a version, for example --to 3.455.0")
	}
	if to == version {
		cMuted.Printf("  --to %s is the version this CLI is, so this is an ordinary upgrade.\n\n", to)
		return nil // the caller carries on with the normal path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("finding your home directory, which is where the version cache lives: %w", err)
	}
	cache := filepath.Join(home, ".grit", "versions")

	cMuted.Printf("  This CLI is v%s. Fetching v%s to do the upgrade.\n", version, to)
	cMuted.Println("  The binary is verified against the release's signed checksums and")
	cMuted.Printf("  cached in %s. It is not installed on your PATH.\n\n", cache)

	bin, err := selfupdate.FetchVersion(to, cache)
	if err != nil {
		return fmt.Errorf("getting Grit v%s: %w", to, err)
	}

	// Everything the person typed, minus --to itself, so --diff, --plan,
	// --dry-run and --force all work exactly as they do without it.
	args := append([]string{"upgrade"}, passthrough...)

	fmt.Println()
	color.New(color.FgHiBlack).Printf("  → %s %s\n\n", bin, strings.Join(args, " "))

	run := exec.Command(bin, args...)
	run.Stdout, run.Stderr, run.Stdin = os.Stdout, os.Stderr, os.Stdin
	if err := run.Run(); err != nil {
		return fmt.Errorf("Grit v%s could not finish the upgrade: %w", to, err)
	}
	return nil
}

// upgradePassthrough is the flags to hand the other binary: the ones that
// change what the upgrade does, and not --to, which has already been consumed.
//
// Rebuilt from the parsed flags rather than taken from os.Args, so an alias or
// a different spelling (-f, --force=true) arrives as the other binary expects
// it. A flag added to `grit upgrade` and forgotten here is silently dropped, so
// TestUpgradePassthroughCoversEveryFlag reads the command's own flag set.
func upgradePassthrough(force, showDiff, writePlan, dryRun bool) []string {
	var args []string
	if force {
		args = append(args, "--force")
	}
	if showDiff {
		args = append(args, "--diff")
	}
	if writePlan {
		args = append(args, "--plan")
	}
	if dryRun {
		args = append(args, "--dry-run")
	}
	return args
}
