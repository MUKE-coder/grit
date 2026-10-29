package ui

import (
	"fmt"
	"io"
	"os"
	"sync"
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
//
// The running line really does animate. It did not until v3.421.0: Start drew
// Frames()[0] once and nothing ever drew the second frame, so every stage
// showed a motionless glyph. Nobody noticed while the stages took 600ms. On a
// cold module cache "Go dependencies" takes minutes, and a frozen spinner over
// a silent terminal is indistinguishable from a hang, which is what it was
// reported as.
type Stepper struct {
	w           io.Writer
	interactive bool
	column      int

	// mu guards everything the ticker goroutine touches, which is the writer
	// and the fields the running line is drawn from.
	mu      sync.Mutex
	label   string
	detail  string
	open    bool
	started time.Time
	frame   int

	stop chan struct{} // closed to ask the ticker to stop
	gone chan struct{} // closed by the ticker once it has
}

// tick is how often the running line is redrawn. Fast enough to read as
// motion, slow enough that a long stage is not thousands of writes.
const tick = 120 * time.Millisecond

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

	s.mu.Lock()
	s.label = label
	s.detail = ""
	s.started = time.Now()
	s.frame = 0
	s.open = true
	interactive := s.interactive
	s.mu.Unlock()

	if !interactive {
		return
	}

	s.mu.Lock()
	fmt.Fprintln(s.w, s.runningLine())
	s.mu.Unlock()

	s.stop = make(chan struct{})
	s.gone = make(chan struct{})
	go s.animate(s.stop, s.gone)
}

// Detail sets the note shown after the label on the running line: what the
// stage is doing right now, when it knows.
//
// "Go dependencies" alone says nothing about whether a four minute wait is
// progress or a hang. "Go dependencies  downloading github.com/gin-gonic/gin"
// answers that without the reader having to trust us.
//
// Safe to call from another goroutine, and a no-op once the stage is closed.
func (s *Stepper) Detail(detail string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.open || !s.interactive {
		return
	}
	s.detail = detail
	s.redrawLocked()
}

// animate rewrites the running line until the stage is closed.
func (s *Stepper) animate(stop <-chan struct{}, gone chan<- struct{}) {
	defer close(gone)
	t := time.NewTicker(tick)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
			s.mu.Lock()
			if s.open {
				s.frame++
				s.redrawLocked()
			}
			s.mu.Unlock()
		}
	}
}

// redrawLocked erases the running line and draws it again. Caller holds mu.
func (s *Stepper) redrawLocked() {
	s.eraseLocked()
	fmt.Fprintln(s.w, s.runningLine())
}

// runningLine is the whole line, already truncated to one terminal row.
//
// One row is not cosmetic: erase moves the cursor up exactly one line, so a
// line that wrapped would leave its first row on screen for ever and every
// later redraw would eat the line above it instead.
func (s *Stepper) runningLine() string {
	frames := Frames()
	line := Running(frames[s.frame%len(frames)], Muted.Render(s.label))

	// Elapsed, once there is enough of it to be worth reading. This is the
	// part that says "still working" without anybody trusting the spinner.
	visible := len(s.label)
	if elapsed := time.Since(s.started); elapsed >= 2*time.Second {
		took := fmt.Sprintf("%.0fs", elapsed.Seconds())
		line += pad(s.label, s.column) + Muted.Render(took)
		visible = s.column + len(took)
	}

	if s.detail != "" {
		room := maxLineWidth - visible - 4
		if room > 8 {
			line += "  " + Faint.Render(truncate(s.detail, room))
		}
	}
	return line
}

// maxLineWidth is the width the running line is kept inside. Narrower than
// the 80 every terminal has, because the gutter and the glyph are in it too.
const maxLineWidth = 72

func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	if max <= 1 {
		return ""
	}
	return string(r[:max-1]) + "…"
}

// Finish closes the last stage. Calling it twice does nothing the second time.
func (s *Stepper) Finish() { s.done() }

// Fail closes the running stage as the one that stopped everything, so the
// error that follows has something to attach to.
func (s *Stepper) Fail() {
	if !s.stopAnimation() {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.eraseLocked()
	fmt.Fprintln(s.w, Gutter+Error.Render(GlyphFail)+" "+s.label)
	s.open = false
}

// done writes the finished line for whatever was running.
func (s *Stepper) done() {
	if !s.stopAnimation() {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.eraseLocked()
	fmt.Fprintln(s.w, StepTimed(s.label, time.Since(s.started), s.column))
	s.open = false
}

// stopAnimation halts the ticker and waits for it, so that nothing can draw a
// running line after the finished one. Reports whether a stage was open.
func (s *Stepper) stopAnimation() bool {
	s.mu.Lock()
	open := s.open
	s.mu.Unlock()
	if !open {
		return false
	}
	if s.stop != nil {
		close(s.stop)
		<-s.gone
		s.stop, s.gone = nil, nil
	}
	return true
}

// eraseLocked removes the running line, when there was one to draw.
func (s *Stepper) eraseLocked() {
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
