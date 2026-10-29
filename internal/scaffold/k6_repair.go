package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

// repairK6Suite fixes a load-test suite that never ran.
//
// The shared journey in tests/k6/lib/common.js asked for /health, which no
// Grit API serves, and /api/blogs, which only a project with a Blog resource
// serves. Run against a real project every request failed, every threshold
// breached, and all six load tests reported a broken system before measuring
// anything about it. One of its checks was "body is non-empty", which an
// error response passes, so the suite was not silent so much as misleading.
//
// The journey now signs in once, reads /api/health and /api/profile, and
// reads one of your own resources when RESOURCE names it. The six entry
// scripts re-export setup, because k6 only calls it when the entry script
// does.
//
// Only files Grit still owns are rewritten: the upgrade guard keeps an edited
// one, and a project that deleted a script keeps it deleted.
func repairK6Suite(root string, opts Options) error {
	dir := filepath.Join(root, "tests", "k6")
	if !dirExists(dir) {
		return nil
	}

	files := map[string]string{
		filepath.Join(dir, "README.md"):        k6ReadmeMD(opts),
		filepath.Join(dir, "lib", "common.js"): k6CommonJS(),
		filepath.Join(dir, "smoke.js"):         k6SmokeJS(),
		filepath.Join(dir, "average-load.js"):  k6AverageLoadJS(),
		filepath.Join(dir, "stress.js"):        k6StressJS(),
		filepath.Join(dir, "spike.js"):         k6SpikeJS(),
		filepath.Join(dir, "soak.js"):          k6SoakJS(),
		filepath.Join(dir, "breakpoint.js"):    k6BreakpointJS(),
	}

	for path, content := range files {
		if !fileExists(path) {
			// A script the project deleted stays deleted. lib/common.js is the
			// exception: every remaining script imports it, so a suite without
			// it does not run at all.
			if filepath.Base(path) != "common.js" {
				continue
			}
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return fmt.Errorf("creating %s: %w", filepath.Dir(path), err)
			}
		}
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}
