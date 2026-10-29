package scaffold

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// UpgradeDryRun says what an upgrade would do to this project, and changes
// nothing in it.
//
// --diff was documented as the way to look before upgrading ("Not a dry run:
// use --diff for that") and it applies the upgrade: on a real project it
// rewrote 46 tracked files and added about 40. A flag that promises a look
// and writes is the worst shape a flag can have. See grit#133.
//
// This runs the REAL upgrade, on a copy, and diffs the two trees. The
// alternative was a dry-run switch threaded through every writer, and there
// are 54 places that write a file without passing through guardedWrite: a
// dry run that still writes some of them would be the same bug wearing a
// different name. Running the real thing somewhere else cannot drift from
// what an upgrade does, because it IS what an upgrade does.
func UpgradeDryRun(uOpts UpgradeOptions) error {
	root, err := FindProjectRoot()
	if err != nil {
		return err
	}

	cyan := color.New(color.FgHiCyan)
	muted := color.New(color.FgHiBlack)

	cyan.Printf("\n  Dry run: upgrading a copy of this project, changing nothing here.\n")
	muted.Printf("  Everything the upgrade prints below happens to the copy.\n")
	muted.Printf("  Copying %s\n", filepath.Base(root))

	temp, err := os.MkdirTemp("", "grit-dryrun-*")
	if err != nil {
		return fmt.Errorf("making a working copy: %w", err)
	}
	defer os.RemoveAll(temp)

	copyRoot := filepath.Join(temp, filepath.Base(root))
	before, err := copyProjectForDryRun(root, copyRoot)
	if err != nil {
		return fmt.Errorf("making a working copy: %w", err)
	}
	muted.Printf("  %d files copied\n\n", len(before))

	// The upgrade finds its project by walking up from the working directory,
	// so the copy has to be the working directory while it runs.
	previous, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(copyRoot); err != nil {
		return fmt.Errorf("entering the working copy: %w", err)
	}
	runErr := Upgrade(uOpts)
	if err := os.Chdir(previous); err != nil {
		return fmt.Errorf("leaving the working copy: %w", err)
	}
	if runErr != nil {
		return runErr
	}

	// The upgrade has just printed a full run, ending in "Upgrade
	// complete". Said plainly, because that line is alarming when you
	// asked for a dry run and it is describing a directory that is about
	// to be deleted.
	muted.Printf("\n  ---- end of the copy's output. Nothing above happened to this project. ----\n")

	after, err := hashTree(copyRoot)
	if err != nil {
		return fmt.Errorf("reading the working copy: %w", err)
	}

	reportDryRun(before, after)
	return nil
}

// dryRunSkip are directories an upgrade neither reads nor writes, and that
// make a copy slow enough to look broken.
var dryRunSkip = map[string]bool{
	"node_modules": true,
	".next":        true,
	".git":         true,
	".turbo":       true,
	"dist":         true,
	"build":        true,
	".expo":        true,
	"vendor":       true,
}

// copyProjectForDryRun copies the project and returns a hash per relative
// path, which is the "before" to diff against.
func copyProjectForDryRun(src, dst string) (map[string]string, error) {
	hashes := map[string]string{}

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}
		if info.IsDir() {
			if dryRunSkip[info.Name()] && rel != "." {
				return filepath.SkipDir
			}
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		sum, copyErr := copyAndHash(path, filepath.Join(dst, rel))
		if copyErr != nil {
			// A file being written by something else is not worth failing a
			// dry run over; it shows up as unreadable rather than as a change.
			return nil
		}
		hashes[filepath.ToSlash(rel)] = sum
		return nil
	})

	return hashes, err
}

func copyAndHash(src, dst string) (string, error) {
	in, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return "", err
	}
	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(out, h), in); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func hashTree(root string) (map[string]string, error) {
	hashes := map[string]string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if dryRunSkip[info.Name()] && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		sum := sha256.Sum256(raw)
		hashes[filepath.ToSlash(rel)] = hex.EncodeToString(sum[:])
		return nil
	})
	return hashes, err
}

// reportDryRun prints what the upgrade did to the copy, which is what it
// would do here.
func reportDryRun(before, after map[string]string) {
	var added, changed, removed []string

	for path, sum := range after {
		switch old, had := before[path]; {
		case !had:
			added = append(added, path)
		case old != sum:
			changed = append(changed, path)
		}
	}
	for path := range before {
		if _, still := after[path]; !still {
			removed = append(removed, path)
		}
	}
	sort.Strings(added)
	sort.Strings(changed)
	sort.Strings(removed)

	green := color.New(color.FgHiGreen)
	yellow := color.New(color.FgHiYellow)
	red := color.New(color.FgHiRed)
	cyan := color.New(color.FgHiCyan)
	muted := color.New(color.FgHiBlack)

	fmt.Println()
	cyan.Printf("  What the upgrade would do here\n\n")

	section := func(c *color.Color, label string, paths []string) {
		if len(paths) == 0 {
			return
		}
		c.Printf("  %s (%d)\n", label, len(paths))
		for _, p := range paths {
			fmt.Printf("      %s\n", p)
		}
		fmt.Println()
	}

	section(green, "added", added)
	section(yellow, "changed", changed)
	section(red, "removed", removed)

	if len(added)+len(changed)+len(removed) == 0 {
		muted.Printf("  Nothing. This project is already at this version.\n\n")
		return
	}

	muted.Printf("  Nothing in this project was touched. Run grit upgrade to apply it.\n")
	// .grit/manifest.json changes on every upgrade by construction, so saying
	// so here saves somebody working out why a "no change" upgrade has one.
	if _, ok := changedContains(changed, ".grit/manifest.json"); ok {
		muted.Printf("  .grit/manifest.json always changes: it records what Grit wrote.\n")
	}
	fmt.Println()
}

func changedContains(paths []string, want string) (int, bool) {
	for i, p := range paths {
		if strings.EqualFold(p, want) {
			return i, true
		}
	}
	return 0, false
}
