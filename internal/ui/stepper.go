package ui

import (
	"fmt"
	"io"
	"os"
	"time"
)

// A Stepper prints a list of stages the way the style guide asks: the one that
// is running shows as running, and a finished one is a check and a noun.
//
// "Go API", not "Scaffolding Go API...". The difference shows at the end, when
// somebody reads the whole block at once: a list of finished nouns is a summary
// of what they got, and a list of present participles is a transcript of a
// process that is over.
//
// One call per stage. Start closes the previous stage as finished and opens the
// next, because the code it wraps is a sequence of "do this, then do that" and
// a second call per stage is one more thing to forget. Finish closes the last.
//
// An interactive terminal gets the running line rewritten in place. Anything
// else (a pipe, a CI log, NO_COLOR) gets one plain line per finished stage and
// no cursor control at all, because a log full of escape codes is worse than a
// log with no spinner in it.
type Stepper struct {
	w           io.Writer
	interactive bool
	column      int

	label   string
	open    bool
	started time.Time
}

// NewStepper writes to stdout, deciding once whether it may move the cursor.
//
// column is where timings line up. It only affects alignment, never
// correctness: pass the longest label you expect plus a few.
func NewStepper(column int) *Stepper {
	return &Stepper{w: os.Stdout, interactive: Interactive(), column: column}
}

// Start finishes the stage that was running and begins a new one.
func (s *Stepper) Start(label string) {
	s.done()
	s.label = label
	s.started = time.Now()
	s.open = true
	if s.interactive {
		fmt.Fprintln(s.w, Running(Frames()[0], Muted.Render(label)))
	}
}

// Finish closes the last stage. Calling it twice does nothing the second time.
func (s *Stepper) Finish() { s.done() }

// Fail closes the running stage as the one that stopped everything, so the
// error that follows has something to attach to.
func (s *Stepper) Fail() {
	if !s.open {
		return
	}
	s.erase()
	fmt.Fprintln(s.w, Gutter+Error.Render(GlyphFail)+" "+s.label)
	s.open = false
}

// done writes the finished line for whatever was running.
func (s *Stepper) done() {
	if !s.open {
		return
	}
	s.erase()
	fmt.Fprintln(s.w, StepTimed(s.label, time.Since(s.started), s.column))
	s.open = false
}

// erase removes the running line, when there was one to draw.
func (s *Stepper) erase() {
	if s.interactive {
		fmt.Fprint(s.w, "\033[1A\033[2K\r")
	}
}

// Interactive reports whether stdout is a terminal that will act on cursor
// codes, and whether the user has asked for plain output.
//
// NO_COLOR is honoured here as well as by lipgloss: somebody who sets it wants
// a log they can read, and a line rewritten in place is as much of a problem
// as a colour.
func Interactive() bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("CI") != "" {
		return false
	}
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
