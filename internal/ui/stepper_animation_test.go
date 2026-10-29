package ui

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"
)

/*
 * "The step that is running animates" was a comment, in two files, for as
 * long as the stepper existed. It was not true: Start drew Frames()[0] once
 * and nothing drew another frame, so every stage showed a motionless glyph.
 *
 * It cost nothing while stages took 600ms. On a cold Go module cache the
 * dependency stage takes minutes, and a frozen spinner over a silent terminal
 * was reported, correctly, as the scaffolder having hung.
 *
 * So the claim is a test now.
 */

// newTestStepper is a Stepper that writes somewhere readable and believes it
// is interactive, which is the only mode that draws a running line at all.
func newTestStepper() (*Stepper, *syncBuffer) {
	buf := &syncBuffer{}
	return &Stepper{w: buf, interactive: true, column: 46}, buf
}

// syncBuffer is a bytes.Buffer that survives being written from the ticker
// goroutine while the test reads it.
type syncBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}

func TestTheRunningStepActuallyAnimates(t *testing.T) {
	s, buf := newTestStepper()

	s.Start("Go dependencies")
	// Long enough for several ticks. The stage is deliberately not finished
	// until after the wait, because that is the situation being tested: a
	// stage that takes a long time.
	time.Sleep(5 * tick)
	s.Finish()

	out := buf.String()

	frames := Frames()
	seen := map[string]bool{}
	for _, f := range frames {
		if strings.Contains(out, f) {
			seen[f] = true
		}
	}
	if len(seen) < 2 {
		t.Fatalf("only %d distinct spinner frame(s) were ever drawn, so the "+
			"spinner is static. Frames seen: %v\noutput:\n%q",
			len(seen), seen, out)
	}
}

func TestAFinishedStepLeavesNoSpinnerBehind(t *testing.T) {
	s, buf := newTestStepper()

	s.Start("Go dependencies")
	time.Sleep(3 * tick)
	s.Finish()

	// Give any goroutine that outlived Finish a chance to write after it.
	time.Sleep(3 * tick)

	out := buf.String()
	last := lastLine(out)
	for _, f := range Frames() {
		if strings.Contains(last, f) {
			t.Fatalf("a spinner frame was drawn after the step finished; the "+
				"last line is %q", last)
		}
	}
	if !strings.Contains(last, "Go dependencies") {
		t.Fatalf("the finished line is missing; last line is %q", last)
	}
}

// Detail is what turns an unexplained wait into a progress report, so it has
// to reach the screen.
func TestDetailReachesTheRunningLine(t *testing.T) {
	s, buf := newTestStepper()

	s.Start("Go dependencies")
	s.Detail("downloading github.com/gin-gonic/gin (3)")
	time.Sleep(2 * tick)
	s.Finish()

	if !strings.Contains(buf.String(), "github.com/gin-gonic/gin") {
		t.Fatalf("the detail never reached the terminal:\n%q", buf.String())
	}
}

// A detail set on a closed stage must not draw: it would appear below a
// finished step and erase the line above it on the next redraw.
func TestDetailAfterFinishDrawsNothing(t *testing.T) {
	s, buf := newTestStepper()

	s.Start("Go dependencies")
	s.Finish()
	before := buf.String()

	s.Detail("downloading something")
	if buf.String() != before {
		t.Fatalf("Detail wrote after the stage closed:\n%q", buf.String())
	}
}

// The non-interactive path draws no running line at all, and must still be
// safe to call Detail on: a pipe and CI both take it.
func TestNonInteractiveNeverDrawsARunningLine(t *testing.T) {
	buf := &syncBuffer{}
	s := &Stepper{w: buf, interactive: false, column: 46}

	s.Start("Go dependencies")
	s.Detail("downloading github.com/gin-gonic/gin")
	time.Sleep(2 * tick)
	s.Finish()

	out := buf.String()
	if strings.Contains(out, "\033[") {
		t.Fatalf("cursor control reached a non-interactive writer:\n%q", out)
	}
	for _, f := range Frames() {
		if strings.Contains(out, f) {
			t.Fatalf("a spinner frame reached a non-interactive writer:\n%q", out)
		}
	}
	if !strings.Contains(out, "Go dependencies") {
		t.Fatalf("the finished step is missing:\n%q", out)
	}
}

// The running line has to stay on one row: erase moves the cursor up exactly
// one line, so a wrapped line would leave a row behind for ever and then eat
// the line above it on the next redraw.
func TestALongDetailCannotWrapTheRunningLine(t *testing.T) {
	s, _ := newTestStepper()
	s.Start("Go dependencies")
	s.mu.Lock()
	s.detail = strings.Repeat("github.com/some/very/long/module/path", 10)
	line := s.runningLine()
	s.mu.Unlock()
	s.Finish()

	if width := len([]rune(stripANSI(line))); width > maxLineWidth+Gutters {
		t.Fatalf("the running line is %d columns wide, which will wrap: %q",
			width, stripANSI(line))
	}
}

// Gutters is the slack the width check allows for the leading gutter and the
// spinner glyph, which are not part of the budget runningLine works to.
const Gutters = 8

func stripANSI(s string) string {
	var out strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == 0x1b {
			for i < len(s) && s[i] != 'm' {
				i++
			}
			continue
		}
		out.WriteByte(s[i])
	}
	return out.String()
}

func lastLine(s string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	return lines[len(lines)-1]
}
