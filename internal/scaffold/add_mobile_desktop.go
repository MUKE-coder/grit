package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// There was no way to add the Expo app or the desktop app to a project that
// already existed: `grit add` had neither, and --expo and --desktop live only
// on `grit new`. So the way in was to scaffold a second project and copy the
// directory across, which is what the reporter of grit#125 did, and what
// Outpost did.
//
// A copied-in app has no manifest records, so until v3.444.0 the next upgrade
// treated every file in it as Grit's and unedited: it replaced the tab bar,
// the API client, the auth provider and the query client, and wrote back six
// screens the project had deleted. The guard covers that now, but the first
// upgrade after a copy still restores deleted files once, because a file with
// no record and no bytes on disk is indistinguishable from a file this
// version has newly added.
//
// These commands remove the copying. They write the app through the same
// function `grit new` uses, with the manifest recording as it goes, so the
// provenance is right from the first upgrade rather than the second.

// AddExpo writes the Expo app into a project that does not have one.
func AddExpo(root string, force bool) error {
	return addApp(root, force, appAddition{
		Key:   "expo",
		Name:  "Expo app",
		Dir:   filepath.Join("apps", "expo"),
		Write: writeExpoFiles,
		Next: []string{
			installStep(root, filepath.Join("apps", "expo")),
			"grit start expo",
		},
	})
}

// AddDesktop writes the Wails desktop app into a project that does not have
// one.
func AddDesktop(root string, force bool) error {
	return addApp(root, force, appAddition{
		Key:   "desktop",
		Name:  "desktop app",
		Dir:   filepath.Join("apps", "desktop"),
		Write: writeDesktopClientFiles,
		Next: []string{
			installStep(root, filepath.Join("apps", "desktop", "frontend")),
			"cd apps/desktop && wails dev",
		},
	})
}

// appAddition is one app that can be added after the fact: the grit.json flag
// that turns it on, the directory it lives in, and the writer that fills it.
type appAddition struct {
	Key   string
	Name  string
	Dir   string
	Write func(root string, opts Options) error
	Next  []string
}

func addApp(root string, force bool, app appAddition) error {
	opts := projectOptions(root)

	dir := filepath.Join(root, app.Dir)
	if dirExists(dir) && !force {
		// Not an error: somebody running this twice wants to know where they
		// stand, not to be told off.
		fmt.Printf("  %s is already there, so the %s was left alone.\n",
			filepath.ToSlash(app.Dir), app.Name)
		fmt.Printf("  grit add %s --force   # write the framework's version over it\n", app.Key)
		return nil
	}

	// grit.json first, because the writers read it: the root package.json's
	// dev script, pnpm-workspace.yaml's packages, and grit.config.ts all gate
	// on it, and an app that is on disk and not in grit.json is one the
	// upgrade will not maintain.
	if err := setAppIncluded(root, app.Key, true); err != nil {
		return err
	}
	switch app.Key {
	case "expo":
		opts.IncludeExpo = true
	case "desktop":
		opts.IncludeDesktop = true
	}

	// Recorded as it is written, which is the whole point of the command.
	// Without this the files go down with no provenance and the next upgrade
	// cannot tell the project's edits from the template (grit#125).
	release, err := manifest.Start(root, DefaultVersion, "add-"+app.Key)
	if err != nil {
		return err
	}
	writeErr := app.Write(root, opts)
	if err := release(); err != nil && writeErr == nil {
		writeErr = err
	}
	if writeErr != nil {
		return fmt.Errorf("writing the %s: %w", app.Name, writeErr)
	}

	// And the root files that mention it. These are framework-owned and may
	// have been edited, so they go through the same writer the upgrade uses:
	// an edited one is left alone and reported rather than replaced.
	if err := rewireRootForApp(root, opts); err != nil {
		return err
	}

	fmt.Printf("\n  ✓ %s written to %s\n", app.Name, filepath.ToSlash(app.Dir))
	fmt.Printf("  ✓ grit.json now says %s: true, so grit upgrade maintains it\n", app.Key)
	fmt.Println()
	fmt.Println("  Next:")
	for _, line := range app.Next {
		fmt.Printf("    %s\n", line)
	}
	return nil
}

// installStep is how to install this app's dependencies, which is not the same
// answer in every project.
//
// Both apps are self-contained: neither depends on packages/shared, so both
// work in an --api project, which has no root package.json and no pnpm
// workspace at all. Telling somebody there to run `pnpm install` at the root
// sends them to a file that does not exist.
func installStep(root, dir string) string {
	if fileExists(filepath.Join(root, "package.json")) {
		return "pnpm install"
	}
	return "cd " + filepath.ToSlash(dir) + " && pnpm install"
}

// rewireRootForApp rewrites the root files whose content depends on which apps
// a project has: the workspace globs, the dev scripts, and the config.
func rewireRootForApp(root string, opts Options) error {
	files := map[string]string{
		filepath.Join(root, "pnpm-workspace.yaml"): pnpmWorkspaceFor(
			!dirExists(filepath.Join(root, "apps")) && dirExists(filepath.Join(root, "frontend")),
			opts.ShouldIncludeDesktop(), opts.ShouldIncludeExpo()),
		filepath.Join(root, "package.json"):   rootPackageJSON(opts),
		filepath.Join(root, "grit.config.ts"): gritConfig(opts),
	}
	for path, body := range files {
		if _, err := os.Stat(path); err != nil {
			continue // a project that does not have this file does not want one
		}
		if err := writeFile(path, body); err != nil {
			return fmt.Errorf("rewiring %s: %w", filepath.Base(path), err)
		}
	}
	return nil
}

// setAppIncluded flips one app's flag in grit.json.
//
// Edited as text rather than decoded and re-encoded, because grit.json is a
// file people read and comment on, and a round trip through a struct drops
// every key this version of Grit does not know about.
func setAppIncluded(root, key string, on bool) error {
	path := filepath.Join(root, "grit.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading grit.json: %w", err)
	}
	src := string(data)

	want := fmt.Sprintf("%q: %t", key, on)
	if strings.Contains(src, want) {
		return nil
	}
	from := fmt.Sprintf("%q: %t", key, !on)
	if !strings.Contains(src, from) {
		return fmt.Errorf("grit.json has no %q flag to set, so the upgrade would not "+
			"maintain the app: add %s to it by hand", key, want)
	}
	return os.WriteFile(path, []byte(strings.Replace(src, from, want, 1)), 0o644)
}

// projectOptions describes the project as it stands, the way the upgrade reads
// it: from grit.json and from what is on disk, not from flags nobody passed.
func projectOptions(root string) Options {
	opts := Options{
		ProjectName:  filepath.Base(root),
		Style:        readProjectStyle(root),
		Theme:        readProjectTheme(root),
		Version:      DefaultVersion,
		Frontend:     readProjectFrontend(root),
		Architecture: detectArchitecture(root),
	}
	// Which apps the project has, probed from disk. The architecture carries
	// web and admin; these two are the optional extras.
	opts.IncludeExpo = dirExists(filepath.Join(root, "apps", "expo"))
	opts.IncludeDesktop = dirExists(filepath.Join(root, "apps", "desktop", "frontend"))
	opts.Normalize()
	return opts
}
