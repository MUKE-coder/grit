package doctor

import (
	"strings"
	"testing"
)

// The hint has to name a command that refreshes the table and nothing else.
//
// It suggested `grit upgrade`, which moves the whole project onto the CLI's
// templates, and `grit generate resource`, which rewrites the route file and
// drops the hand-written guard that made the table stale in the first place
// (grit#104). Somebody following either one loses work or makes a far larger
// decision than the warning called for.
func TestStaleAccessTableHintNamesTheNarrowCommand(t *testing.T) {
	// The check only fires on a real project, so assert the constant the
	// finding carries rather than standing one up for a single string.
	fix := accessTableStaleFix

	if !strings.Contains(fix, "grit routes --write") {
		t.Errorf("the hint should name grit routes --write, got %q", fix)
	}
	if strings.Contains(fix, "grit upgrade") {
		t.Errorf("the hint still suggests grit upgrade for a stale table: %q", fix)
	}
}
