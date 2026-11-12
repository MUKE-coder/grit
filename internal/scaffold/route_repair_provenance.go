package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Which *_routes.go files an authorization repair may rewrite.
//
// repairReviewCriticals walks every file in internal/routes whose name ends in
// _routes.go and moves `m.Protected.<VERB>("/<plural>...", h.<Method>)` onto
// the staff group behind a permission. The reasoning was that the edit anchors
// on the exact text the generator wrote, so it cannot reach a developer's own
// code. That reasoning is wrong, and grit#213 is what it cost.
//
// A project wrote internal/routes/reference_interchange_routes.go by hand,
// mounting two student-facing endpoints scoped by owner inside the handler. It
// holds `h := &handlers.ReferenceInterchangeHandler{`, and its two routes are
// POST .../import and GET .../export, which are names in routePermission. So
// every anchor matched, the upgrade rewrote the file, and two endpoints a
// student needs went behind perm:projects.create and perm:projects.view. The
// first symptom would have been a 403 in production.
//
// The anchors cannot tell generated code from hand-written code, because the
// generator's shape is the convention the documentation teaches. What can tell
// them apart is the manifest: it records every file Grit wrote.
//
// # Why not simply require the manifest
//
// Because a project generated before v3.147.0 has no manifest, and every file
// in it reads as Untracked (manifest.Load says so in as many words). Gating on
// the manifest alone would quietly stop repairing the projects that most need
// repairing, which is the failure shape of a gate added inside an existing
// repair.
//
// So the question asked is not "is this file tracked" but "does this project's
// manifest know about its route files at all". If it tracks even one
// *_routes.go, the manifest is authoritative here and an untracked route file
// is the developer's. If it tracks none, there is nothing to reason with and
// the old behaviour stands, because a project with open resource routes is a
// real exposure and refusing to fix it is not the cautious option.
type routeProvenance struct {
	// known reports whether the manifest tracks any route file, which is what
	// makes its silence about one file meaningful.
	known bool
	// tracked is the set of absolute paths the manifest records.
	tracked map[string]bool
}

// routeFileProvenance reads the manifest's record of a project's route files.
func routeFileProvenance(root string, m *manifest.Manifest, files []string) routeProvenance {
	p := routeProvenance{tracked: map[string]bool{}}
	for _, path := range files {
		key, inside := manifest.Rel(root, path)
		if !inside {
			continue
		}
		if _, ok := m.Files[manifest.Key(key)]; ok {
			p.known = true
			p.tracked[path] = true
		}
	}
	return p
}

// mayRewrite reports whether an authorization repair may edit this file, and
// if not, why not in a sentence for the developer.
func (p routeProvenance) mayRewrite(path string) (bool, string) {
	if !p.known {
		// No manifest coverage of routes anywhere: pre-v3.147.0, or a manifest
		// that lost its entries. Behave as before.
		return true, ""
	}
	if p.tracked[path] {
		return true, ""
	}
	return false, "this file is not one Grit wrote, so its routes were left exactly as they are"
}

// describeRouteMove renders one route's before and after for the upgrade
// summary. An authorization change is worth naming even in a file that
// genuinely is generated: a project may have granted the old shape on purpose,
// and a rewritten file otherwise reads as success (grit#213).
func describeRouteMove(method, path, from, to string) string {
	return fmt.Sprintf("%s %s: %s → %s", method, path, from, to)
}

// untouchedRouteRewrites names the routes a repair would have moved in a file
// it is not allowed to touch, so the refusal says what was at stake rather
// than only that nothing happened.
func untouchedRouteRewrites(src string) []string {
	var out []string
	for _, line := range protectedRouteRe.FindAllStringSubmatch(src, -1) {
		method, plural, rest, name := line[2], line[3], line[4], line[5]
		perm, generated := routePermission[name]
		if !generated {
			continue
		}
		out = append(out, fmt.Sprintf("%s /%s%s (would have asked for perm:%s.%s)",
			method, plural, rest, plural, perm))
	}
	sort.Strings(out)
	return out
}

// handWrittenRouteFiles lists the route files in a project that Grit did not
// write, for `grit doctor` and for the upgrade's own reporting.
func handWrittenRouteFiles(root string, apiRoot string) ([]string, error) {
	m, err := manifest.Load(root)
	if err != nil {
		return nil, err
	}
	files, err := goSources(filepath.Join(apiRoot, "internal", "routes"))
	if err != nil {
		return nil, err
	}
	var routes []string
	for _, path := range files {
		if strings.HasSuffix(path, "_routes.go") {
			routes = append(routes, path)
		}
	}
	p := routeFileProvenance(root, m, routes)
	if !p.known {
		return nil, nil
	}
	var out []string
	for _, path := range routes {
		if !p.tracked[path] {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				rel = path
			}
			out = append(out, filepath.ToSlash(rel))
		}
	}
	sort.Strings(out)
	return out, nil
}

// fileHasText reports whether a file contains a string, for the small checks
// that decide whether a repair is wanted at all.
func fileHasText(path, want string) bool {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(string(raw), want)
}
