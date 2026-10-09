package devlog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func project(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func write(t *testing.T, root, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, Dir), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(Path(root), []byte(body), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

// The log exists so that something other than a person can read what the
// application printed. The tests are about whether it would have found this
// week's bugs.

func TestWritingAndTailing(t *testing.T) {
	root := project(t)

	w, err := Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	for i := 1; i <= 10; i++ {
		fmt.Fprintf(w, "line %d\n", i)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	if !Exists(root) {
		t.Fatal("Exists says there is no log after writing one")
	}

	lines, err := Tail(root, 3)
	if err != nil {
		t.Fatalf("Tail: %v", err)
	}
	if len(lines) != 3 || lines[0] != "line 8" || lines[2] != "line 10" {
		t.Fatalf("Tail(3) = %v", lines)
	}
}

// A nil writer is usable and does nothing.
//
// Because logging must never stop the dev servers: if the file cannot be
// opened, `grit start` carries on without it, and every caller would otherwise
// need to check.
func TestANilWriterIsSafe(t *testing.T) {
	var w *Writer
	if _, err := w.Write([]byte("anything\n")); err != nil {
		t.Fatalf("writing to a nil log: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("closing a nil log: %v", err)
	}
}

// The previous run is kept, because "it worked a minute ago" is the question
// the log is most often opened to answer.
func TestThePreviousRunIsKept(t *testing.T) {
	root := project(t)

	first, err := Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	fmt.Fprintln(first, "the first run")
	first.Close()

	second, err := Open(root)
	if err != nil {
		t.Fatalf("Open again: %v", err)
	}
	fmt.Fprintln(second, "the second run")
	second.Close()

	current, err := Tail(root, 10)
	if err != nil {
		t.Fatalf("Tail: %v", err)
	}
	if strings.Join(current, "\n") != "the second run" {
		t.Fatalf("the current log holds %v", current)
	}

	previous, err := os.ReadFile(Path(root) + ".1")
	if err != nil {
		t.Fatalf("the previous run was not kept: %v", err)
	}
	if !strings.Contains(string(previous), "the first run") {
		t.Fatalf("the previous log holds %q", previous)
	}
}

// The failures this is for, taken from the ones that actually happened.
func TestErrorsFindsWhatMatters(t *testing.T) {
	root := project(t)
	write(t, root, strings.Join([]string{
		"[api] 2026/10/09 03:10:01 starting on :8080",
		"[web] ready in 812ms",
		"[api] internal/handlers/contact.go:42:19: undefined: services.NewContactService",
		"[api] 2026/10/09 03:10:14 [GIN] | 500 |   1.2ms | POST /api/v1/contacts",
		"[api] panic: runtime error: invalid memory address or nil pointer dereference",
		"[api] goroutine 42 [running]:",
		"[api] main.(*ContactHandler).Create(0x0)",
		"[web] compiled successfully",
	}, "\n")+"\n")

	problems, err := Errors(root, 20)
	if err != nil {
		t.Fatalf("Errors: %v", err)
	}
	if len(problems) < 3 {
		t.Fatalf("found %d problems in a log with a compile error, a 500 and a panic", len(problems))
	}

	all := ""
	for _, p := range problems {
		all += p.Text + "\n" + strings.Join(p.More, "\n") + "\n"
	}
	for _, want := range []string{
		"undefined: services.NewContactService",
		"500",
		"panic: runtime error",
	} {
		if !strings.Contains(all, want) {
			t.Errorf("the failures do not include %q", want)
		}
	}

	// A panic's stack is the useful part, so it has to come along with it.
	stackFound := false
	for _, p := range problems {
		if strings.Contains(p.Text, "panic:") {
			for _, more := range p.More {
				if strings.Contains(more, "ContactHandler") {
					stackFound = true
				}
			}
		}
	}
	if !stackFound {
		t.Error("the panic came back without its stack, which is the half that says where")
	}
}

// A quiet log has no problems, and says so rather than inventing them.
func TestAQuietLogIsQuiet(t *testing.T) {
	root := project(t)
	write(t, root, strings.Join([]string{
		"[api] 2026/10/09 03:10:01 starting on :8080",
		"[api] 2026/10/09 03:10:02 [GIN] | 200 |  1.1ms | GET /api/v1/contacts",
		"[web] ready in 812ms",
		"[web] compiled successfully",
	}, "\n")+"\n")

	problems, err := Errors(root, 20)
	if err != nil {
		t.Fatalf("Errors: %v", err)
	}
	if len(problems) != 0 {
		t.Fatalf("found %d problems in a healthy log: %+v", len(problems), problems)
	}
}

// A compile failure is read from the hot reloader's own file, because when the
// build fails the binary never starts and prints nothing at all.
func TestABuildFailureIsFound(t *testing.T) {
	root := project(t)
	write(t, root, "[api] starting\n")

	tmp := filepath.Join(root, "apps", "api", "tmp")
	if err := os.MkdirAll(tmp, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	body := "internal/services/contact.go:18:2: imported and not used: \"fmt\"\nexit status 1\n"
	if err := os.WriteFile(filepath.Join(tmp, "build-errors.log"), []byte(body), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	problems, err := Errors(root, 20)
	if err != nil {
		t.Fatalf("Errors: %v", err)
	}
	if len(problems) == 0 {
		t.Fatal("a build failure was not reported")
	}
	if !strings.Contains(problems[0].Text, "build failed") {
		t.Fatalf("the first problem is %q, want the build failure first", problems[0].Text)
	}
	if !strings.Contains(strings.Join(problems[0].More, "\n"), "imported and not used") {
		t.Errorf("the build failure came back without the compiler's message: %v", problems[0].More)
	}
}

// An empty build-errors.log means the last build succeeded, and air leaves the
// file behind. Reporting it would mean every healthy project looks broken.
func TestAnEmptyBuildLogIsNotAFailure(t *testing.T) {
	root := project(t)
	write(t, root, "[api] starting\n")

	tmp := filepath.Join(root, "apps", "api", "tmp")
	if err := os.MkdirAll(tmp, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "build-errors.log"), nil, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	problems, err := Errors(root, 20)
	if err != nil {
		t.Fatalf("Errors: %v", err)
	}
	for _, p := range problems {
		if strings.Contains(p.Text, "build failed") {
			t.Fatal("an empty build-errors.log was reported as a build failure")
		}
	}
}

// No log at all is a thing to say, not an error to raise: the answer is "run
// grit start", and a tool that errors here makes an agent conclude the project
// is broken.
func TestNoLogIsNotAnError(t *testing.T) {
	root := project(t)

	if Exists(root) {
		t.Fatal("Exists says there is a log in an empty project")
	}
	if _, err := Tail(root, 10); err == nil {
		t.Error("Tail on a missing log should say it is missing")
	}
}

// The cap stops a reload loop filling a disk, and says that it did.
func TestTheLogStopsAtItsCap(t *testing.T) {
	root := project(t)

	w, err := Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	// A line a little over 1 KB, written past the cap.
	line := strings.Repeat("x", 1024) + "\n"
	for written := 0; written < maxBytes+(2<<20); written += len(line) {
		if _, err := w.Write([]byte(line)); err != nil {
			t.Fatalf("Write: %v", err)
		}
	}
	w.Close()

	info, err := os.Stat(Path(root))
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Size() > maxBytes+4096 {
		t.Fatalf("the log grew to %d bytes past a %d cap", info.Size(), maxBytes)
	}

	lines, err := Tail(root, 5)
	if err != nil {
		t.Fatalf("Tail: %v", err)
	}
	if !strings.Contains(strings.Join(lines, "\n"), "stopped recording") {
		t.Error("the log stopped recording without saying so, so a reader would think " +
			"the application went quiet")
	}
}

// air's build log can say nothing but that it failed, and saying so beats
// repeating it.
//
// It appends without newlines, so three retries leave one line reading
// "exit status 1exit status 1exit status 1". That was reported as the context
// for the failure, which is worse than no context: it looks like output and
// carries nothing. The compiler's real message goes to the dev log.
func TestABuildLogWithOnlyAnExitStatus(t *testing.T) {
	root := project(t)
	write(t, root, "[api] stat cmd/server: directory not found\n")

	tmp := filepath.Join(root, "apps", "api", "tmp")
	if err := os.MkdirAll(tmp, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Exactly what air leaves behind, with no trailing newline.
	if err := os.WriteFile(filepath.Join(tmp, "build-errors.log"),
		[]byte("exit status 1exit status 1exit status 1"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	problems, err := Errors(root, 20)
	if err != nil {
		t.Fatalf("Errors: %v", err)
	}
	if len(problems) == 0 {
		t.Fatal("a failed build was not reported at all")
	}

	build := problems[0]
	if !strings.Contains(build.Text, "build failed") {
		t.Fatalf("the first problem is %q", build.Text)
	}
	context := strings.Join(build.More, " ")
	if strings.Contains(context, "exit status 1exit status 1") {
		t.Fatal("the repeated exit status was handed back as context")
	}
	if !strings.Contains(context, "development log") {
		t.Errorf("the note does not say where the real message is: %q", context)
	}
}

// And a build log with a real diagnostic still shows it.
func TestABuildLogWithARealMessage(t *testing.T) {
	root := project(t)
	write(t, root, "[api] starting\n")

	tmp := filepath.Join(root, "apps", "api", "tmp")
	if err := os.MkdirAll(tmp, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	body := "internal/services/contact.go:18:2: imported and not used: \"fmt\"\nexit status 1\n"
	if err := os.WriteFile(filepath.Join(tmp, "build-errors.log"), []byte(body), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	problems, err := Errors(root, 20)
	if err != nil {
		t.Fatalf("Errors: %v", err)
	}
	if len(problems) == 0 {
		t.Fatal("not reported")
	}
	if !strings.Contains(strings.Join(problems[0].More, "\n"), "imported and not used") {
		t.Fatalf("the compiler's message was dropped: %v", problems[0].More)
	}
}

// The failure a real run found that the first pattern list missed.
//
// A triple-tier project started before `pnpm install` is a total outage, and
// every line it produces avoids the word error: "'next' is not recognized",
// "apps/web dev: Failed", ERR_PNPM_RECURSIVE_RUN_FIRST_FAIL, "Exit status 1".
// grit logs --errors reported that nothing had failed, which is the worst
// answer available: it is confident and wrong, and it sends the reader looking
// somewhere else.
func TestAMissingToolchainIsAFailure(t *testing.T) {
	root := project(t)
	write(t, root, strings.Join([]string{
		"[web] Scope: 2 of 5 workspace projects",
		"[web] apps/web dev$ rm -rf .next && next dev --webpack --port 3000",
		"[web] apps/web dev: 'next' is not recognized as an internal or external command,",
		"[web] apps/web dev: operable program or batch file.",
		"[web] apps/admin dev: Failed",
		"[web]  ERR_PNPM_RECURSIVE_RUN_FIRST_FAIL  @app/admin@0.1.0 dev: `next dev`",
		"[web] Exit status 1",
		"[web]  ELIFECYCLE  Command failed with exit code 1.",
	}, "\n")+"\n")

	problems, err := Errors(root, 20)
	if err != nil {
		t.Fatalf("Errors: %v", err)
	}
	if len(problems) == 0 {
		t.Fatal("a project whose frontend toolchain is missing was reported as healthy")
	}

	all := ""
	for _, p := range problems {
		all += p.Text + "\n" + strings.Join(p.More, "\n") + "\n"
	}
	for _, want := range []string{"not recognized", "ERR_PNPM"} {
		if !strings.Contains(all, want) {
			t.Errorf("the failures do not include %q", want)
		}
	}
}

// And the widened patterns must not fire on an ordinary successful run, which
// is the half that keeps the tool worth calling.
func TestTheWiderPatternsStayQuietOnASuccess(t *testing.T) {
	root := project(t)
	write(t, root, strings.Join([]string{
		"[web] apps/web dev$ next dev --port 3000",
		"[web] ready in 812ms",
		"[web] compiled successfully",
		"[api] 2026/10/09 05:00:01 starting on :8099",
		"[api] 2026/10/09 05:00:02 [GIN] | 200 |  1.1ms | GET /api/v1/contacts",
		"[api] 2026/10/09 05:00:03 [GIN] | 201 |  4.3ms | POST /api/v1/contacts",
	}, "\n")+"\n")

	problems, err := Errors(root, 20)
	if err != nil {
		t.Fatalf("Errors: %v", err)
	}
	if len(problems) != 0 {
		for _, p := range problems {
			t.Errorf("fired on: %s", p.Text)
		}
		t.Fatal("the widened patterns report failures in a healthy run")
	}
}

// A single-process start writes the same shape as the combined one.
//
// The log existed only when `grit start` ran everything at once, so the
// projects least likely to do that had nothing for an agent to read: an --api
// project has no other way to start, and a --mobile project runs its API and
// its Expo app as separate commands by design. Found by starting a mobile
// project's API with `grit start server` and watching `grit logs` say there
// was no log.
func TestPrefixedWritesTheSameShapeAsACombinedStart(t *testing.T) {
	root := t.TempDir()
	log, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}

	w := Prefixed(log, "api")
	// Written in the chunks a pipe produces, not in whole lines: a tag belongs
	// at the start of a line, and a Writer is handed whatever arrived.
	for _, chunk := range []string{"Server star", "ting on port 8096\npanic: ", "nope\n"} {
		if _, err := w.Write([]byte(chunk)); err != nil {
			t.Fatal(err)
		}
	}
	if err := log.Close(); err != nil {
		t.Fatal(err)
	}

	lines, err := Tail(root, 10)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"[api] Server starting on port 8096", "[api] panic: nope"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d: %q", len(lines), len(want), lines)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Errorf("line %d: got %q, want %q", i, lines[i], want[i])
		}
	}
}

// And the colour a dev server writes does not reach the file.
func TestPrefixedStripsColour(t *testing.T) {
	root := t.TempDir()
	log, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Prefixed(log, "expo").Write([]byte("\x1b[32mready\x1b[39m in 900ms\n")); err != nil {
		t.Fatal(err)
	}
	if err := log.Close(); err != nil {
		t.Fatal(err)
	}

	lines, err := Tail(root, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != "[expo] ready in 900ms" {
		t.Fatalf("got %q, want [\"[expo] ready in 900ms\"]", lines)
	}
}

// The tag is what Errors uses to keep a failure's context to the process that
// failed, so a single-process log has to carry it.
func TestASingleProcessLogIsSearchableByProcess(t *testing.T) {
	root := t.TempDir()
	log, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Prefixed(log, "api").Write([]byte("listening\nError: bind: address already in use\n")); err != nil {
		t.Fatal(err)
	}
	if err := log.Close(); err != nil {
		t.Fatal(err)
	}

	problems, err := Errors(root, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) == 0 {
		t.Fatal("a bind failure in a single-process log was not reported as an error")
	}
}
