package scaffold

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// grit#204: `grit upgrade` rewrote internal/authz/actor.go with an import of
// internal/confidential and never wrote the package, so an upgraded project
// did not compile. One file referenced it and nothing provided it.
//
// The two delivery paths diverged silently, which is the part worth guarding.
// writeFrameworkOwnedFiles travels on upgrade; the new-project map does not.
// A file in the travelling set that imports a package only the new-project map
// writes is correct for `grit new` and broken for every existing project, and
// nothing says so until somebody runs go build.
//
// Read statically rather than by upgrading a project, because the question is
// about the two maps and not about any one project's state: it needs no
// toolchain, no network and no fixture, and it answers in a few seconds.
//
// The module, not {{MODULE}}: writeFrameworkOwnedFiles substitutes it before
// writing, so a pattern matching the placeholder matches nothing on disk. The
// first version of this test did that and passed with the fix reverted, which
// is the only reason the vacuity was noticed.
func internalImportRe(module string) *regexp.Regexp {
	return regexp.MustCompile(`"` + regexp.QuoteMeta(module) + `/internal/([a-z0-9_]+)"`)
}

// deliveredElsewhere are the packages an owned file may import without this
// set writing them, because an upgraded project is certain to have them.
//
// Two reasons qualify, and each entry says which:
//
//   - another writer that `grit upgrade` also calls provides it, so the pair
//     still travels together, just not in one map; or
//   - every project has had it since long before any of this mattered, so
//     there is no project to deliver it to.
//
// Nothing else qualifies. A package introduced in a release is not "obviously
// always there", however it looks in hindsight: internal/confidential looked
// exactly like that, and the only projects that had it were the ones generated
// after it shipped (grit#204). Adding a name here to make this test pass is
// how that bug ships again, so add the package to the map instead.
var deliveredElsewhere = map[string]bool{
	// Written by another writer upgrade.go calls.
	"appendonly": true, // writeAppendOnlyFiles
	"audit":      true, // writeAuditChainFiles
	"crypto":     true, // writeCryptoFiles
	"events":     true, // writeCodegenRuntimeFiles
	"ids":        true, // writeCodegenRuntimeFiles
	"paginate":   true, // writeCodegenRuntimeFiles
	"respond":    true, // writeRespondFiles
	"sms":        true, // writeRecoveryFiles

	// Older than the problem. No project that can be upgraded lacks these.
	"config": true,
	"mail":   true,
	"sync":   true,
}

func TestEveryOwnedImportIsDeliveredToo(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "probe", Architecture: ArchTriple}
	apiRoot := opts.APIRoot(root)

	// The map as writeFrameworkOwnedFiles builds it, by writing it and reading
	// the result back: the map is a local, and a copy of it here would be the
	// second thing that has to agree with the first.
	if err := writeFrameworkOwnedFiles(root, opts); err != nil {
		t.Fatalf("writing the owned set: %v", err)
	}

	written := map[string]bool{} // internal/<pkg> this set provides
	imported := map[string][]string{}
	importRe := internalImportRe(opts.Module())

	for _, path := range goFilesUnder(t, apiRoot) {
		rel, err := filepath.Rel(apiRoot, path)
		if err != nil {
			t.Fatal(err)
		}
		parts := strings.Split(filepath.ToSlash(rel), "/")
		if len(parts) >= 2 && parts[0] == "internal" {
			written[parts[1]] = true
		}
		for _, m := range importRe.FindAllStringSubmatch(readFileOrEmpty(path), -1) {
			imported[m[1]] = append(imported[m[1]], filepath.ToSlash(rel))
		}
	}

	var missing []string
	for pkg, by := range imported {
		if written[pkg] || deliveredElsewhere[pkg] {
			continue
		}
		sort.Strings(by)
		missing = append(missing, pkg+" (imported by "+strings.Join(by, ", ")+")")
	}
	sort.Strings(missing)

	if len(missing) > 0 {
		t.Errorf("writeFrameworkOwnedFiles travels on upgrade and imports %d internal package(s) "+
			"it does not write:\n  %s\n\nAn upgraded project gets the import and not the package, "+
			"so it does not compile (grit#204). Either add the package to this same map, or, if "+
			"every project has had it since before upgrades existed, add it to deliveredElsewhere with "+
			"a reason.", len(missing), strings.Join(missing, "\n  "))
	}
}

// The specific regression: the file and the package it imports are in one set.
func TestActorAndTheConfidentialRegistryTravelTogether(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "probe", Architecture: ArchTriple}
	apiRoot := opts.APIRoot(root)
	if err := writeFrameworkOwnedFiles(root, opts); err != nil {
		t.Fatalf("writing the owned set: %v", err)
	}

	actor := filepath.Join(apiRoot, "internal", "authz", "actor.go")
	if !fileContains(actor, "/internal/confidential\"") {
		t.Skip("actor.go no longer imports the confidential registry")
	}
	registry := filepath.Join(apiRoot, "internal", "confidential", "confidential.go")
	if !fileExists(registry) {
		t.Error("the upgrade set writes an actor.go importing internal/confidential and does not " +
			"write the package, so an upgraded project does not compile (grit#204)")
	}
}

// goFilesUnder is every .go file below dir.
func goFilesUnder(t *testing.T, dir string) []string {
	t.Helper()
	var out []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".go") {
			out = append(out, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking %s: %v", dir, err)
	}
	sort.Strings(out)
	return out
}
