package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// scaffoldGosecBudget is the longest the package's file list may be, in
// characters, as CI hands it to gosec.
//
// go/packages splits a command line longer than 16383 characters into chunks
// and loads each chunk as its own package (safeArgMax in
// x/tools/go/packages/packages.go, derived from the Windows limit, so Linux
// behaves the same). gosec hands it this package's files as a command line
// rather than a pattern, so once the list crosses that line the package is
// type-checked in halves, and the half that uses Options no longer sees
// scaffold.go. What gosec then reports is `undefined: Options` against code
// that `go build` and `go vet` both accept, which reads as a compile error in
// the repository rather than a limit in the scanner. It cost an afternoon.
//
// CI runs gosec from inside each package directory so the list is basenames,
// which is the cheapest way to stay well under. The budget is set below the
// real limit so this test fails while there is still room to decide what to
// do, rather than after the build has gone red.
//
// If this fails: split internal/scaffold into subpackages. Raising the number
// only works up to 16383, and the failure above it is not self-explanatory.
const scaffoldGosecBudget = 13000

func TestScaffoldFileListFitsGosecsCommandLine(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}

	total, count := 0, 0
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		// One argument, plus the separator between it and the next.
		total += len(filepath.Join(".", name)) + 1
		count++
	}

	if total > scaffoldGosecBudget {
		t.Errorf("the %d files in internal/scaffold are %d characters of gosec command line, "+
			"over the %d budget.\n\n"+
			"go/packages splits a list over 16383 characters into separately type-checked "+
			"chunks, and gosec then reports undefined symbols for code that compiles. "+
			"Split this package into subpackages; raising the budget runs out at 16383.",
			count, total, scaffoldGosecBudget)
	}
	t.Logf("%d files, %d characters, budget %d", count, total, scaffoldGosecBudget)
}
