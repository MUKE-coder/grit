package scaffold

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/MUKE-coder/grit/v3/internal/ui"
)

// Resolving Go dependencies, with the terminal told what is happening.
//
// This used to be exec.Command("go", "mod", "tidy").CombinedOutput(): one
// blocking call that swallowed every byte the tool printed and returned them
// only on failure. On a warm module cache it finished in a second and nobody
// looked. On a cold one it downloads the whole dependency graph, which over a
// slow link is minutes, and the user saw a motionless spinner over a silent
// terminal for all of it. It was reported as the scaffolder hanging, which is
// the correct reading of the evidence they had.
//
// go mod tidy is not quiet during that wait. It prints a line per module to
// stderr:
//
//	go: downloading github.com/gin-gonic/gin v1.10.0
//
// So the wait is only opaque because we made it opaque. Streaming those lines
// onto the running step turns an unexplained pause into a progress report,
// without the command taking a millisecond longer.
//
// Nothing here tries to make the download faster. It is a one-time fill of the
// module cache over whatever connection the user has, and a scaffolder that
// quietly swapped their GOPROXY to be quicker would be a worse tool than a
// slow one.

// goModTidy runs `go mod tidy` in dir, reporting each module it downloads to
// the running step.
//
// steps may be nil, for callers that have no stepper to report to. On failure
// the whole output is returned in the error, as it was before: a tidy that
// fails needs its message read, not summarised.
func goModTidy(dir string, steps *ui.Stepper) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir

	// Both pipes, because tidy writes progress to stderr and the occasional
	// note to stdout, and an error needs whichever one carried it.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	var (
		mu        sync.Mutex
		collected strings.Builder
		wg        sync.WaitGroup
		count     int
	)

	drain := func(r io.Reader) {
		defer wg.Done()
		sc := bufio.NewScanner(r)
		// A single go.sum mismatch line can be long; the default 64KiB token
		// limit is enough, but a module path plus a version plus a hash is
		// worth the headroom.
		sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for sc.Scan() {
			line := sc.Text()

			mu.Lock()
			collected.WriteString(line)
			collected.WriteByte('\n')
			if mod, ok := downloadingModule(line); ok {
				count++
				if steps != nil {
					steps.Detail(fmt.Sprintf("downloading %s (%d)", mod, count))
				}
			}
			mu.Unlock()
		}
	}

	wg.Add(2)
	go drain(stdout)
	go drain(stderr)
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("running go mod tidy: %w\n%s", err, collected.String())
	}
	if steps != nil && count > 0 {
		steps.Detail(fmt.Sprintf("%d modules downloaded", count))
	}
	return nil
}

// downloadingModule pulls the module path out of a progress line.
//
//	go: downloading github.com/gin-gonic/gin v1.10.0  ->  github.com/gin-gonic/gin
//
// Anything else is not progress, and returns false rather than being shown:
// the running line has room for one short fact, and a warning rendered as
// progress is worse than no progress at all.
func downloadingModule(line string) (string, bool) {
	const prefix = "go: downloading "
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	rest := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	if rest == "" {
		return "", false
	}
	if i := strings.IndexByte(rest, ' '); i > 0 {
		rest = rest[:i]
	}
	return rest, true
}
