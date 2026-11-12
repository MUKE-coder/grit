package routeparser

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Refusing a route that declares no permission (grit#179, proposal 1).
//
// A generated resource's routes carry their permission because the generator
// writes it. A route somebody writes carries whatever they remembered, and a
// route with nothing in front of it is reachable by any signed-in account.
// Forgetting makes a route MORE permissive, which is the wrong way round: the
// mistake nobody notices is the one that costs something.
//
// # Why this does not check every route
//
// A freshly generated project has 67 routes behind authentication that demand
// no permission, and almost all of them are right: /auth/me, /api-keys,
// /auth/passkeys/:id. They are the caller's own things, scoped in the handler,
// and asking each to declare a permission would be 67 false positives on the
// first run. A check that is mostly wrong is one people learn to pass with
// --no-verify.
//
// So the default is the routes the PROJECT wrote: files the manifest has no
// record of Grit writing. That is exactly the reporter's situation, where 105
// hand-written routes needed a 364-line guard and a 429-line catalog to get
// this property back.
//
// Strict checks everything, for a project willing to annotate the framework's
// routes too. It is a flag rather than the default because the first run of a
// new check has to be believable.

// Finding is one route that declares nothing.
type Finding struct {
	Route Route
	// File is where the route is written, relative to the project root.
	File string
	// Framework is true when Grit wrote the file, so the finding is only
	// reported in strict mode.
	Framework bool
}

// CheckResult is what a check run found.
type CheckResult struct {
	// Undeclared are the routes behind authentication that demand no
	// permission and give no reason.
	Undeclared []Finding
	// Declared counts the routes that do say, either with a permission or
	// with openToAnySignedIn.
	Declared int
	// Checked is how many routes the run looked at.
	Checked int
	// Strict records which mode produced this.
	Strict bool
	// ManifestKnown reports whether the manifest could tell framework files
	// from hand-written ones. When it cannot, nothing is treated as
	// hand-written and the non-strict run says so rather than passing in
	// silence.
	ManifestKnown bool
}

// OK reports a run with nothing to fix.
func (r CheckResult) OK() bool { return len(r.Undeclared) == 0 }

// Check reads a project's routes and reports the ones that declare no
// permission.
//
// root is the project root. strict includes the routes in files Grit wrote.
func Check(root string, strict bool) (*CheckResult, error) {
	routesFile, err := FindRoutesFile(root)
	if err != nil {
		return nil, err
	}
	routes, err := ParseProject(routesFile)
	if err != nil {
		return nil, err
	}

	owned, known, err := frameworkRouteFiles(root, filepath.Dir(routesFile))
	if err != nil {
		return nil, err
	}

	res := &CheckResult{Strict: strict, ManifestKnown: known}
	for _, rt := range routes {
		if !rt.Access.Authenticated {
			continue // public: a different decision, made elsewhere and visible
		}
		res.Checked++
		if !rt.Access.NeedsPermission() {
			res.Declared++
			continue
		}

		file := relSlash(root, rt.File)
		isFramework := known && owned[file]
		if isFramework && !strict {
			continue
		}
		res.Undeclared = append(res.Undeclared, Finding{
			Route:     rt,
			File:      file,
			Framework: isFramework,
		})
	}

	sort.Slice(res.Undeclared, func(i, j int) bool {
		if res.Undeclared[i].File != res.Undeclared[j].File {
			return res.Undeclared[i].File < res.Undeclared[j].File
		}
		return res.Undeclared[i].Route.Path < res.Undeclared[j].Route.Path
	})
	return res, nil
}

// frameworkRouteFiles is the set of route files the manifest says Grit wrote,
// and whether the manifest knows about route files at all.
//
// The second return matters. A project from before the manifest has no entries
// for anything, so every file reads as hand-written and a non-strict run would
// report all 67 framework routes. Knowing that it cannot tell is what lets the
// caller say so instead.
func frameworkRouteFiles(root, routesDir string) (map[string]bool, bool, error) {
	m, err := manifest.Load(root)
	if err != nil {
		return nil, false, err
	}
	entries, err := os.ReadDir(routesDir)
	if err != nil {
		return nil, false, fmt.Errorf("reading %s: %w", routesDir, err)
	}

	owned := map[string]bool{}
	known := false
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		abs := filepath.Join(routesDir, e.Name())
		key, inside := manifest.Rel(root, abs)
		if !inside {
			continue
		}
		if _, tracked := m.Files[manifest.Key(key)]; tracked {
			owned[relSlash(root, abs)] = true
			known = true
		}
	}
	return owned, known, nil
}

// relSlash renders a path relative to the root with forward slashes.
func relSlash(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(rel)
}
