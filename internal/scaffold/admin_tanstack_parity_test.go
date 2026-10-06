package scaffold

import (
	"path/filepath"
	"strings"
	"testing"
)

// The two admin front-ends keep separate file lists, and the pages are shared.
//
// adminFileMap writes the Next.js admin; adminTanStackFileMap writes the Vite
// one, re-emitting most of the same source through nextToTanStack. A component
// added to the first map and used in a shared page therefore exists for Next.js
// only, and the Vite build fails on an unresolved import: not a type error
// somewhere, the whole frontend refusing to build.
//
// This has now shipped twice. v3.119.0 put TwoFactorCard into a profile page
// both variants render, and the sign-in-link release did the same with
// security-nudges.tsx, which every dashboard variant imports, so every --single
// project built from it died on `vite build` over one missing file. Neither the
// Go suite nor CI's admin build can see it: the TSX lives inside Go string
// literals, and CI builds the Next.js variant.
//
// Only components, lib and hooks are compared. Pages and routes legitimately
// differ (app/<x>/page.tsx against src/pages and src/routes), and the
// per-feature writers that use the adminLib and adminComponent helpers serve
// both front-ends by construction, so they are in neither map and belong in
// neither side of this comparison.
func TestBothAdminFrontendsGetTheSameSharedSource(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "app", Architecture: ArchTriple, Theme: "atlas"}

	next := adminRelativePaths(t, adminFileMap(root, opts), filepath.Join(root, "apps", "admin"))
	tanstack := adminRelativePaths(t, adminTanStackFileMap(root, opts), filepath.Join(root, "apps", "admin", "src"))

	shared := []string{"components/", "lib/", "hooks/"}
	compared := 0
	for file := range next {
		if !hasAnyAdminPrefix(file, shared) {
			continue
		}
		compared++
		if !tanstack[file] {
			t.Errorf("%s is written for the Next.js admin and not the Vite one: if any shared page imports it, `vite build` fails on an unresolved import", file)
		}
	}

	// A comparison of nothing passes. The first draft of this test matched the
	// wrong import alias and reported success with the broken list in place.
	if compared < 80 {
		t.Fatalf("only %d shared files were compared, so the paths have moved and this test proves nothing", compared)
	}
}

func hasAnyAdminPrefix(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

// adminRelativePaths reduces a file map to slash-separated paths below base, so
// the Next.js layout and the TanStack one (which adds a src/ level) compare.
func adminRelativePaths(t *testing.T, files map[string]string, base string) map[string]bool {
	t.Helper()
	out := map[string]bool{}
	for path := range files {
		rel, err := filepath.Rel(base, path)
		if err != nil || strings.HasPrefix(rel, "..") {
			continue
		}
		out[filepath.ToSlash(rel)] = true
	}
	if len(out) < 50 {
		t.Fatalf("only %d files below %s, so this is comparing the wrong directories", len(out), base)
	}
	return out
}

// grit.json records the frontend that was built.
//
// A single project had one answer for years, so Frontend was left empty and
// --single --next accepted the flag and produced a Vite app anyway. There are
// two answers now, and the file has to say which one is on disk: every tool that
// reads it believes it.
func TestGritJSONRecordsTheFrontendThatWasBuilt(t *testing.T) {
	// --single on its own is still the Vite SPA, and says so rather than leaving
	// the field empty for each reader to interpret.
	viteSingle := gritJSON(Options{ProjectName: "app", Architecture: ArchSingle, Version: "9.9.9"})
	if !strings.Contains(viteSingle, `"frontend": "tanstack"`) {
		t.Errorf("--single did not record the frontend it built:\n%s", viteSingle)
	}

	// And --single --next is a Next app, not a flag that gets ignored.
	nextSingle := gritJSON(Options{ProjectName: "app", Architecture: ArchSingle, Frontend: FrontendNext, Version: "9.9.9"})
	if !strings.Contains(nextSingle, `"frontend": "next"`) {
		t.Errorf("--single --next recorded something other than next, so the file disagrees with the directory beside it:\n%s", nextSingle)
	}

	// And it still reports the real choice where there is one to report.
	for _, tc := range []struct {
		frontend Frontend
		want     string
	}{
		{FrontendNext, `"frontend": "next"`},
		{FrontendTanStack, `"frontend": "tanstack"`},
	} {
		got := gritJSON(Options{ProjectName: "app", Architecture: ArchTriple, Frontend: tc.frontend, Version: "9.9.9"})
		if !strings.Contains(got, tc.want) {
			t.Errorf("a triple project with %s did not record %s:\n%s", tc.frontend, tc.want, got)
		}
	}
}
