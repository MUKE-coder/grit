// Package devlog records what the development servers printed, so that
// something other than a person can read it.
//
// A Grit application logs to stderr, which is right for a container and leaves
// an agent with nothing: Laravel's Boost can read `storage/logs/laravel.log`
// and there was no equivalent here. The consequence is not theoretical. The
// bugs that mattered most in a week of building a Grit application were all
// found by reading output by hand, and none of them showed in an exit code: a
// protocol message sent in the wrong mode, an authentication exchange failing
// on a parse and reporting itself as a bad password, a server answering "no" to
// every request for encryption. A person tailing a file found each one in
// minutes. An agent could not look.
//
// So `grit start` tees everything it prints into .grit/logs/dev.log, and the
// MCP server can read it back. The file is local state and is not committed:
// .grit/ is already gitignored apart from the manifest.
package devlog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Dir is where the logs live, under the project's existing local-state
// directory.
const Dir = ".grit/logs"

// Name is the current run's log.
const Name = "dev.log"

// maxBytes caps one run's log.
//
// A dev server left running for a day with a reload loop writes a great deal,
// and a log that fills a disk is a worse bug than the one it was helping with.
// The cap is per run and the previous run is kept, so there is always the
// current session and the one before it, which is what "it worked a minute
// ago" needs.
const maxBytes = 16 << 20

// Path is the current log file for a project.
func Path(root string) string {
	return filepath.Join(root, Dir, Name)
}

// previousPath is last run's.
func previousPath(root string) string {
	return filepath.Join(root, Dir, Name+".1")
}

// Writer is a size-capped file that stops writing rather than growing.
type Writer struct {
	file    *os.File
	written int64
	full    bool
}

// Open starts a new log for this run, keeping the previous one.
//
// Returns a nil Writer and no error when the log cannot be opened. Failing to
// record output must never stop the dev servers from starting: the log is a
// convenience and the servers are the point.
func Open(root string) (*Writer, error) {
	dir := filepath.Join(root, Dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	// Keep one generation. Rename rather than copy, so a long log costs
	// nothing to rotate.
	_ = os.Remove(previousPath(root))
	_ = os.Rename(Path(root), previousPath(root))

	file, err := os.OpenFile(Path(root), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return &Writer{file: file}, nil
}

// Write records a line, and stops at the cap.
func (w *Writer) Write(p []byte) (int, error) {
	// A nil Writer is usable and does nothing, so every caller does not need
	// to check whether logging was available.
	if w == nil || w.file == nil {
		return len(p), nil
	}
	if w.full {
		return len(p), nil
	}
	if w.written+int64(len(p)) > maxBytes {
		w.full = true
		_, _ = fmt.Fprintf(w.file, "\n[grit] this log reached %d MB and stopped recording. "+
			"Restart `grit start` for a fresh one.\n", maxBytes>>20)
		return len(p), nil
	}
	n, err := w.file.Write(p)
	w.written += int64(n)
	return len(p), err
}

// Close flushes and closes.
func (w *Writer) Close() error {
	if w == nil || w.file == nil {
		return nil
	}
	return w.file.Close()
}

// Tail returns the last n lines of the current log, oldest first.
func Tail(root string, n int) ([]string, error) {
	if n <= 0 {
		n = 100
	}
	lines, err := readLines(Path(root))
	if err != nil {
		return nil, err
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return lines, nil
}

// errorLine is what counts as something going wrong.
//
// Deliberately broad. A missed error is an agent concluding the application is
// healthy; a false positive is one extra line it can read and discard. The
// patterns cover Go panics and build failures, the API's own error logging,
// Next.js and Vite failures, and HTTP statuses in the 5xx range, which is how
// Gin's request log reports a handler that blew up.
var errorLine = regexp.MustCompile(`(?i)` + strings.Join([]string{
	`\bpanic:`,
	`\bfatal\b`,
	`\[error\]`,
	`\berror\b.*:`,
	`\bERRO\b`,
	`\bWARN(ING)?\b.*fail`,
	`cannot find|undefined:|no such file`,
	`\b5\d\d\b\s*\|`, // Gin's request log: "| 500 |"
	`unhandled|uncaught`,
	`failed to|failure`,
	// The ones a real run found that this list missed.
	//
	// A triple-tier project started before `pnpm install` printed "'next' is
	// not recognized as an internal or external command", "apps/web dev:
	// Failed", ERR_PNPM_RECURSIVE_RUN_FIRST_FAIL and "Exit status 1", and
	// grit logs --errors reported that nothing had failed. Every one of those
	// is a total outage and none of them contains the word error.
	`not recognized as an internal`,
	`command not found|not found:`,
	`ERR_[A-Z_]+`,
	`FAILED?`,
	`ELIFECYCLE`,
	`[Ee]xit status [1-9]`,
	`\.go:\d+:\d+:`, // a compiler diagnostic
	`Module not found|Failed to compile`,
}, "|"))

// A panic's stack follows the panic line and is the useful part, so a few
// lines after a match come along with it.
const trailingContext = 12

// Problem is one thing that went wrong, with the lines that explain it.
type Problem struct {
	// Line is the line number in the log, so an agent can ask for more around
	// it rather than the whole file.
	Line int      `json:"line"`
	Text string   `json:"text"`
	More []string `json:"more,omitempty"`
}

// Errors reads the recent problems out of the development log and the build
// error log, newest last.
//
// Both sources, because they fail differently and an agent needs whichever
// happened: a compile error means the binary never started, and a panic means
// it started and died.
func Errors(root string, limit int) ([]Problem, error) {
	if limit <= 0 {
		limit = 20
	}

	problems := make([]Problem, 0, limit)

	// air writes compile failures here, and they are the most common thing to
	// go wrong straight after generating code.
	for _, build := range buildErrorLogs(root) {
		lines, err := readLines(build)
		if err != nil || len(lines) == 0 {
			continue
		}
		// The newest failure only, with the noise recognised for what it is.
		//
		// air appends to this file without newlines, so three failed builds leave
		// one 39-byte line reading "exit status 1exit status 1exit status 1". It
		// carries no diagnostic at all: the compiler's actual message went to the
		// dev log. Reporting the repetition as context made the answer worse than
		// saying nothing, so a file that is only exit statuses is reported as a
		// failed build and points at where the reason is.
		lines = withoutRepeats(lastLines(lines, 25))
		body := strings.TrimSpace(strings.Join(lines, "\n"))
		if body == "" {
			continue
		}
		problem := Problem{
			Line: 0,
			Text: fmt.Sprintf("build failed (%s)", mustRel(root, build)),
		}
		if onlyExitStatus(body) {
			problem.More = []string{
				"The hot reloader recorded only an exit status, not a reason. " +
					"The compiler's message is in the development log: read the [api] lines.",
			}
		} else {
			problem.More = lines
		}
		problems = append(problems, problem)
		continue

	}

	lines, err := readLines(Path(root))
	if err != nil && len(problems) == 0 {
		return nil, err
	}

	for i, line := range lines {
		if !errorLine.MatchString(line) {
			continue
		}
		problem := Problem{Line: i + 1, Text: strings.TrimSpace(line)}
		end := i + 1 + trailingContext
		if end > len(lines) {
			end = len(lines)
		}
		// Context comes only from the same process.
		//
		// `grit start` interleaves the API, the frontends and the hot
		// reloader, so the lines after an [api] failure are usually [web]
		// telling you Vite is ready. Taking them made every stack trace
		// arrive wrapped in unrelated noise, which was obvious the first time
		// a real log was read and invisible in a synthetic one.
		owner := processPrefix(line)
		for _, more := range lines[i+1 : end] {
			if strings.TrimSpace(more) == "" {
				continue
			}
			if processPrefix(more) != owner {
				continue
			}
			// Stop at the next unrelated error: it gets its own entry.
			if errorLine.MatchString(more) {
				break
			}
			problem.More = append(problem.More, strings.TrimSpace(more))
		}
		problems = append(problems, problem)
	}

	if len(problems) > limit {
		problems = problems[len(problems)-limit:]
	}
	return problems, nil
}

// buildErrorLogs finds air's compile-failure logs, wherever the API lives.
func buildErrorLogs(root string) []string {
	candidates := []string{
		filepath.Join(root, "apps", "api", "tmp", "build-errors.log"),
		filepath.Join(root, "api", "tmp", "build-errors.log"),
		filepath.Join(root, "tmp", "build-errors.log"),
	}
	found := make([]string, 0, 1)
	for _, path := range candidates {
		info, err := os.Stat(path)
		if err != nil || info.Size() == 0 {
			continue
		}
		found = append(found, path)
	}
	sort.Strings(found)
	return found
}

// Exists reports whether there is a log to read, so a tool can say "run grit
// start first" instead of "no such file".
func Exists(root string) bool {
	info, err := os.Stat(Path(root))
	return err == nil && info.Size() > 0
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	// Webpack and turbo print very long lines, and a truncated error message
	// is the one you need in full.
	scanner.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		return lines, err
	}
	return lines, nil
}

// onlyExitStatus reports whether a build log says nothing but that it failed.
//
// air writes "exit status 1" and appends on each retry without a newline, so
// the whole file can be that phrase repeated. There is no diagnostic to show.
func onlyExitStatus(body string) bool {
	stripped := strings.ReplaceAll(body, "exit status", "")
	for _, r := range stripped {
		if r != ' ' && r != '\n' && r != '\t' && r != '\r' && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

// lastLines is the tail of a slice.
func lastLines(lines []string, n int) []string {
	if len(lines) <= n {
		return lines
	}
	return lines[len(lines)-n:]
}

// withoutRepeats drops lines already seen, keeping order.
//
// A compiler prints the same diagnostic for every file that fails the same way,
// and a reader needs it once.
func withoutRepeats(lines []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		seen[trimmed] = true
		out = append(out, trimmed)
	}
	return out
}

// processPrefix is the "[api] " style tag grit start writes, or "" when a line
// has none.
func processPrefix(line string) string {
	if !strings.HasPrefix(line, "[") {
		return ""
	}
	if end := strings.Index(line, "]"); end > 0 {
		return line[:end+1]
	}
	return ""
}

func mustRel(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(rel)
}

// ansiEscape matches the colour codes a dev server writes. They belong on a
// terminal and are noise in a file written for something that cannot see
// colour.
var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")

// StripANSI removes those codes.
func StripANSI(s string) string { return ansiEscape.ReplaceAllString(s, "") }

// Prefixed wraps a log so each line written to it is tagged with the process
// it came from, exactly as the combined `grit start` writes them.
//
// The tag is what tells an API panic from a bundler's output when both are in
// one file, and Errors uses it to keep the context around a failure to the
// process that failed. A single-process start writes the same shape so that
// `grit logs` cannot tell the two cases apart.
func Prefixed(w io.Writer, label string) io.Writer {
	return &prefixWriter{w: w, prefix: "[" + label + "] "}
}

// prefixWriter buffers until a newline, because a Writer is handed whatever
// chunk the pipe produced and a tag belongs at the start of a line rather than
// at the start of a read.
type prefixWriter struct {
	w      io.Writer
	prefix string
	buf    []byte
}

func (p *prefixWriter) Write(b []byte) (int, error) {
	n := len(b)
	p.buf = append(p.buf, b...)
	for {
		i := bytes.IndexByte(p.buf, '\n')
		if i < 0 {
			break
		}
		line := strings.TrimRight(string(p.buf[:i]), "\r")
		p.buf = p.buf[i+1:]
		if _, err := fmt.Fprintln(p.w, p.prefix+StripANSI(line)); err != nil {
			return n, err
		}
	}
	return n, nil
}
