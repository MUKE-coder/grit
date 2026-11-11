package scaffold

import (
	"strings"
	"testing"
)

// The error page can fix the error most likely to put somebody on it.
//
// Next splits an app into hashed chunks, and a deploy replaces them. A tab
// that was open across one asks for a filename that is gone, gets a 404 and
// throws, and the browser shows "Application error: a client-side exception
// has occurred". Nothing is wrong with the app.
//
// The page already existed and looked fine. What it offered was reset(), and
// reset re-renders the same tree, which asks for the same missing file. So the
// only button on the page could not fix the commonest reason for seeing it.
// Observed on Grit's own site after three releases in one day.
func TestTheErrorPageRecoversFromAStaleChunk(t *testing.T) {
	for name, src := range map[string]string{
		"apps/web":   webErrorPage(),
		"apps/admin": adminErrorPage(),
	} {
		if !strings.Contains(src, "function isStaleChunk(error: Error): boolean {") {
			t.Errorf("%s: the error page cannot tell a stale chunk from any other failure, "+
				"so it offers a reset that cannot help", name)
			continue
		}

		// Matched on the message, because the class name does not survive
		// minification, and because three runtimes word it three ways.
		for _, phrase := range []string{
			"chunkloaderror",
			"loading chunk",
			"failed to fetch dynamically imported module",
			"importing a module script failed",
		} {
			if !strings.Contains(src, phrase) {
				t.Errorf("%s: the stale-chunk check does not match %q, so that runtime's "+
					"wording falls through to the generic branch", name, phrase)
			}
		}

		// A reload, which fetches the new filenames. reset() cannot.
		if !strings.Contains(src, "window.location.reload()") {
			t.Errorf("%s: nothing reloads, so a stale chunk has no way out of the page", name)
		}

		// And exactly one automatic reload per tab. Without the guard, an
		// error this misidentifies becomes an infinite refresh, which is
		// worse than the error it was trying to recover from.
		if !strings.Contains(src, `const RELOAD_GUARD = "grit:reloaded-after-stale-chunk"`) {
			t.Errorf("%s: the automatic reload has no guard and can loop", name)
		}
		if !strings.Contains(src, "sessionStorage.getItem(RELOAD_GUARD)") {
			t.Errorf("%s: the guard is declared and never read", name)
		}
		// Storage can throw (private mode, storage disabled). Reloading
		// anyway, with no guard, is the loop this is avoiding.
		if !strings.Contains(src, "} catch {") {
			t.Errorf("%s: a sessionStorage failure is not handled, and reloading without "+
				"the guard is the loop it exists to prevent", name)
		}

		// The copy changes, because "Something went wrong" is wrong: nothing
		// did, the app was updated.
		if !strings.Contains(src, "A new version was deployed") {
			t.Errorf("%s: a stale chunk is still reported as an unexpected error", name)
		}
	}
}

// One source for the detection, not two.
//
// The two pages are separate templates with different copy and different
// button helpers, and the question "is this a stale chunk?" has to have one
// answer. Two copies of a list of substrings drift, and the drift is silent:
// each page keeps working for the wordings its own copy still knows.
func TestBothErrorPagesShareTheStaleChunkCheck(t *testing.T) {
	web, admin := webErrorPage(), adminErrorPage()

	const start = "function isStaleChunk(error: Error): boolean {"
	const end = "  return true;\n}"

	extract := func(src string) string {
		at := strings.Index(src, start)
		to := strings.Index(src[at:], end)
		if at == -1 || to == -1 {
			return ""
		}
		return src[at : at+to+len(end)]
	}

	a, b := extract(web), extract(admin)
	if a == "" || b == "" {
		t.Fatal("could not read the helper out of both pages")
	}
	if a != b {
		t.Error("the web and admin error pages carry different stale-chunk checks. " +
			"They will disagree about what counts as one, and the page that falls " +
			"behind keeps working for the wordings it still knows, silently")
	}
}
