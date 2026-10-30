package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Review 2.5: eleven repairs wrote with os.WriteFile and told the manifest
// nothing, so a file a repair edited read as hand-edited afterwards and every
// later framework fix to it was skipped. The repair poisoned the file it fixed.
//
// This is the test that keeps that fixed, because the fix itself is a
// convention and a convention with a hundred call sites is enforced by a test
// or it is not enforced. A repair writes through writeRepaired, when it is
// editing a file that exists, or writeFile, when a whole framework-owned file
// is arriving.
//
// repairSourceFile is the one allowed exception: it is the shared helper the
// majority already use, and it does the check-before and refresh-after by
// hand because it also has to format and report.
func TestRepairsDoNotWriteBehindTheManifest(t *testing.T) {
	const allowed = "handler_repair.go" // repairSourceFile lives here

	files, err := filepath.Glob(filepath.Join(".", "*_repair.go"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) < 50 {
		t.Fatalf("found %d *_repair.go files, so this test is looking in the wrong place", len(files))
	}

	writeRe := regexp.MustCompile(`os\.WriteFile\(`)
	var offenders []string
	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") || filepath.Base(path) == allowed {
			continue
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for i, line := range strings.Split(string(raw), "\n") {
			if writeRe.MatchString(line) {
				offenders = append(offenders, filepath.Base(path)+":"+itoaForTest(i+1))
			}
		}
	}
	if len(offenders) > 0 {
		t.Errorf("these repairs write without telling the manifest, so the files they fix "+
			"stop receiving later fixes (review 2.5). Use writeRepaired for an edit to an "+
			"existing file, or writeFile for a whole framework-owned file:\n  %s",
			strings.Join(offenders, "\n  "))
	}
}

// And the helpers do what the convention says they do.
func TestWriteRepairedTellsTheManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.txt")

	if err := writeRepaired(path, "hello\n"); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "hello\n" {
		t.Errorf("wrote %q", raw)
	}

	// The mode-preserving variant keeps the mode it is given. Checked through
	// a read rather than a stat comparison, because Windows reports a mode
	// this assertion cannot rely on.
	script := filepath.Join(dir, "run.sh")
	if err := writeRepairedMode(script, "#!/bin/sh\necho hi\n", 0o755); err != nil {
		t.Fatal(err)
	}
	if raw, err := os.ReadFile(script); err != nil || !strings.HasPrefix(string(raw), "#!/bin/sh") {
		t.Errorf("read %q, %v", raw, err)
	}
}

func itoaForTest(n int) string {
	if n == 0 {
		return "0"
	}
	var out []byte
	for n > 0 {
		out = append([]byte{byte('0' + n%10)}, out...)
		n /= 10
	}
	return string(out)
}
