package scaffold

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/fatih/color"
)

// Review 2.12: the scaffold aborted mid-write when `go mod tidy` failed.
//
// The reviewer hit it offline on a cold module cache: 441 files of 731 written,
// no Docker, no shared package, no web, no admin, nothing cleaned up, and the
// full `grit new` usage block printed under the error. A half-written project
// is worse than either finishing or not starting, because it looks like a
// project.
//
// The desktop path in the same file already did the right thing and said why.
// This is the test that keeps the other two from going back, because the
// tempting shape when reading that code is `return err`.
func TestResolvingDependenciesIsNotFatal(t *testing.T) {
	src, err := os.ReadFile("scaffold.go")
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(src), "\n")

	call := regexp.MustCompile(`if err := goModTidy\(`)
	var sites int
	for i, line := range lines {
		if !call.MatchString(line) {
			continue
		}
		sites++
		// The three lines after the call are the body. A bare `return err`
		// there is the bug: it leaves the project half written.
		body := strings.Join(lines[i+1:min(i+4, len(lines))], "\n")
		if strings.Contains(body, "return err") {
			t.Errorf("scaffold.go:%d returns when `go mod tidy` fails, which leaves the "+
				"project half written: no Docker, no shared package, no frontends, and "+
				"nothing cleaned up. Report it and carry on, as the desktop path does "+
				"(review 2.12).\n%s", i+1, body)
		}
	}
	if sites < 3 {
		t.Errorf("found %d goModTidy call sites, want at least 3: this test is not "+
			"checking what it thinks it is", sites)
	}
}

// And the failure is carried to the end, because by then the warning has
// scrolled past a hundred lines and the person is about to run a command that
// needs the go.sum.
func TestAPendingTidyIsReportedAtTheEnd(t *testing.T) {
	t.Cleanup(ResetDependenciesPending)

	ResetDependenciesPending()
	if got := DependenciesPending(); got != "" {
		t.Errorf("a fresh scaffold reports %q pending", got)
	}

	// What the call site does when tidy fails. Output goes to stdout, which
	// the test does not care about; the recorded directory is the contract.
	captureStdout(t, func() { tidyFailed(filepath.Join("apps", "api"), errors.New("dial tcp: no route to host")) })
	if got := DependenciesPending(); got != filepath.Join("apps", "api") {
		t.Errorf("DependenciesPending() = %q, want the API directory", got)
	}

	// A second scaffold in the same process does not inherit it.
	ResetDependenciesPending()
	if got := DependenciesPending(); got != "" {
		t.Errorf("after a reset it still reports %q", got)
	}
}

// The message has to say what to run, because "could not resolve Go
// dependencies" on its own leaves somebody with a project and no idea what is
// missing from it.
func TestTheTidyWarningSaysWhatToRun(t *testing.T) {
	t.Cleanup(ResetDependenciesPending)
	out := captureStdout(t, func() {
		tidyFailed(filepath.Join("apps", "api"), errors.New("no network"))
	})
	for _, want := range []string{"go mod tidy", "apps", "api", "go.sum"} {
		if !strings.Contains(out, want) {
			t.Errorf("the warning does not mention %q:\n%s", want, out)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// captureStdout returns what fn printed.
//
// color.Output, not os.Stdout. The color package resolves its writer once at
// init and writes there, so swapping os.Stdout captures nothing: the first
// version of this test reported that the warning said none of the four things
// it says, which was the test looking in the wrong place.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	var buf bytes.Buffer
	previous := color.Output
	color.Output = &buf
	defer func() { color.Output = previous }()
	fn()
	return buf.String()
}
