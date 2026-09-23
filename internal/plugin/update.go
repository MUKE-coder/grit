package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MUKE-coder/grit/v3/internal/codefmt"
)

// Updating an installed plugin in place.
//
// A plugin fix used to reach nobody. `grit plugin add` refuses once a plugin is
// installed, `grit upgrade` never touched plugin files, and remove-then-add was
// the only way through. That path loses every local edit and re-inserts the
// injections at their markers, which moves them relative to code the app wrote
// since: doing it to a shop built on Grit put the payment service below the
// code that constructs things from it, and the API stopped compiling. So the
// day the Stripe plugin learned to record what a renewal charged, no existing
// project could have it.
//
// An update therefore rewrites only what it can prove nobody has touched. Each
// file is fingerprinted when the plugin writes it; on an update a file whose
// fingerprint still matches is replaced, and one that differs is left exactly
// as it is and reported. Files the plugin has newly added simply arrive, which
// is what most plugin releases are. Injections already in place are left where
// they sit, and only a genuinely new one is applied.

// UpdateResult is what one update did, for the caller to print.
type UpdateResult struct {
	Name string
	From string
	To   string

	// Added are files the plugin did not write before.
	Added []string
	// Replaced are files it wrote, that nobody has edited since.
	Replaced []string
	// Edited are files left alone because they no longer match what the plugin
	// wrote. The person who edited them decides what to do.
	Edited []string
	// Unverified are files from a lockfile written before fingerprints existed
	// and which differ from the current template. They are left alone too: with
	// no record of what was written, an edit and an improvement look the same.
	Unverified []string
	// Injected are the new patches this update applied.
	Injected []string
	// Deps are dependencies added, as "apps/web: stripe@^9".
	Deps []string
	// Orphaned are files the lockfile has that the plugin no longer writes.
	// Never deleted: the app may import them.
	Orphaned []string

	// Blocked is set when nothing was done at all, because this project was
	// installed before plugin files were fingerprinted and the update cannot
	// tell an edit from an improvement. Half an update is worse than none: the
	// new file a release adds is usually the one another file has to change to
	// use, so stopping keeps the project consistent and one command fixes it.
	Blocked bool
}

// Changed says whether anything happened, so a caller can stay quiet.
func (r *UpdateResult) Changed() bool {
	return len(r.Added)+len(r.Replaced)+len(r.Injected)+len(r.Deps) > 0
}

// UpdateOptions changes how careful an update is.
type UpdateOptions struct {
	// Overwrite takes the plugin's version of every file, including ones you
	// have edited and ones this project is too old to vouch for.
	//
	// It exists because being careful is not always enough: a plugin that adds
	// a file usually changes another one to use it, and leaving that one behind
	// gives you a project that does not compile. The way out has to be one
	// command and a git diff, not a hand-merge nobody can review.
	Overwrite bool
}

// Update brings an installed plugin up to the version this CLI carries.
func Update(ctx Context, p Plugin) (*UpdateResult, error) {
	return UpdateWith(ctx, p, UpdateOptions{})
}

// UpdateWith is Update with the options spelled out.
func UpdateWith(ctx Context, p Plugin, opts UpdateOptions) (*UpdateResult, error) {
	lock, err := LoadLock(ctx.Root)
	if err != nil {
		return nil, err
	}
	installed, ok := lock.Find(p.Name)
	if !ok {
		return nil, fmt.Errorf("%s is not installed: grit plugin add %s", p.Name, p.Name)
	}

	res := &UpdateResult{Name: p.Name, From: installed.Version, To: p.Version}
	if installed.FileHashes == nil {
		installed.FileHashes = map[string]string{}
	}

	rendered := map[string]string{}
	if p.Files != nil {
		rendered = p.Files(ctx)
	}

	paths := make([]string, 0, len(rendered))
	for rel := range rendered {
		paths = append(paths, rel)
	}
	sortStrings(paths)

	known := map[string]bool{}
	for _, rel := range installed.Files {
		known[rel] = true
	}

	// A first pass decides whether this project can be updated at all. Doing
	// it before anything is written is what keeps a refusal from leaving half
	// an update behind.
	if !opts.Overwrite {
		for _, rel := range paths {
			abs := filepath.Join(ctx.Root, rel)
			current, err := os.ReadFile(abs)
			if err != nil {
				continue // missing: it will simply be written
			}
			want := codefmt.File(abs, rendered[rel])
			if sameText(string(current), want) {
				// Already the new content: fingerprint it now, so a project
				// that predates fingerprints gets quieter every release
				// instead of staying unverifiable forever.
				installed.FileHashes[rel] = hashOf(want)
				continue
			}
			if _, vouched := installed.FileHashes[rel]; !vouched {
				res.Unverified = append(res.Unverified, rel)
			}
		}
		if len(res.Unverified) > 0 {
			res.Blocked = true
			// The fingerprints learned above are worth keeping even though
			// nothing else was done.
			if err := lock.Save(ctx.Root); err != nil {
				return nil, err
			}
			return res, nil
		}
	}

	for _, rel := range paths {
		abs := filepath.Join(ctx.Root, rel)
		want := codefmt.File(abs, rendered[rel])

		current, err := os.ReadFile(abs)
		switch {
		case os.IsNotExist(err):
			// New file in this version of the plugin, or one somebody deleted.
			if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
				return nil, fmt.Errorf("creating dir for %s: %w", rel, err)
			}
			if err := os.WriteFile(abs, []byte(want), 0644); err != nil {
				return nil, fmt.Errorf("writing %s: %w", rel, err)
			}
			installed.FileHashes[rel] = hashOf(want)
			if !known[rel] {
				installed.Files = append(installed.Files, rel)
				known[rel] = true
			}
			res.Added = append(res.Added, rel)
			continue
		case err != nil:
			return nil, fmt.Errorf("reading %s: %w", rel, err)
		}

		if sameText(string(current), want) {
			// Already the new content. Record the fingerprint so the next
			// update can vouch for it even if this lockfile predates them.
			installed.FileHashes[rel] = hashOf(want)
			continue
		}

		recorded, vouched := installed.FileHashes[rel]
		switch {
		case !vouched && !opts.Overwrite:
			res.Unverified = append(res.Unverified, rel)
		case vouched && recorded != hashOf(string(current)) && !opts.Overwrite:
			res.Edited = append(res.Edited, rel)
		default:
			if err := os.WriteFile(abs, []byte(want), 0644); err != nil {
				return nil, fmt.Errorf("writing %s: %w", rel, err)
			}
			installed.FileHashes[rel] = hashOf(want)
			res.Replaced = append(res.Replaced, rel)
		}
	}

	for _, rel := range installed.Files {
		if _, still := rendered[rel]; !still {
			res.Orphaned = append(res.Orphaned, rel)
		}
	}

	// Injections. injectBefore is idempotent against the lines above its own
	// marker, so one already in place costs nothing and is not moved.
	if p.Injections != nil {
		for _, inj := range p.Injections(ctx) {
			abs := filepath.Join(ctx.Root, inj.File)
			if !fileExists(abs) {
				continue
			}
			if hasInjection(installed.Injections, inj.File, inj.Code) {
				continue
			}
			if err := injectBefore(abs, inj.Marker, inj.Code); err != nil {
				if inj.Optional {
					continue
				}
				return nil, fmt.Errorf("injecting into %s: %w", inj.File, err)
			}
			installed.Injections = append(installed.Injections, LockedInjection{
				File:   inj.File,
				Marker: inj.Marker,
				Code:   inj.Code,
			})
			res.Injected = append(res.Injected, inj.File)
		}
	}

	for _, d := range p.NodeDeps {
		if d.Workspace == "" {
			continue
		}
		pkg := filepath.Join(ctx.Root, d.Workspace, "package.json")
		if !fileExists(pkg) {
			continue
		}
		added, err := addNodeDependency(pkg, d.Name, d.Version)
		if err != nil {
			return nil, fmt.Errorf("adding %s to %s: %w", d.Name, d.Workspace, err)
		}
		if added {
			res.Deps = append(res.Deps, fmt.Sprintf("%s: %s %s", d.Workspace, d.Name, d.Version))
		}
	}
	installed.NodeDeps = p.NodeDeps

	now := time.Now().UTC()
	installed.UpdatedAt = &now
	installed.Version = p.Version
	sortStrings(installed.Files)
	if err := lock.Save(ctx.Root); err != nil {
		return nil, err
	}
	return res, nil
}

// hasInjection says whether this exact patch is already recorded.
func hasInjection(applied []LockedInjection, file, code string) bool {
	for _, inj := range applied {
		if inj.File == file && inj.Code == code {
			return true
		}
	}
	return false
}

// sameText compares two files as text, so a line ending is not an edit.
//
// It has to be: git rewrites line endings on checkout on Windows, an older CLI
// wrote some of these files with CRLF, and an editor may convert a file nobody
// meant to change. Counting any of those as "you edited this" would make the
// update refuse to touch a project that is in fact untouched.
func sameText(a, b string) bool {
	return normalizeEOL(a) == normalizeEOL(b)
}

func normalizeEOL(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}
