package scaffold

import (
	"strings"
	"testing"
)

// The template and the upgrade repair say the same thing, so a project that
// upgrades ends up with what a new one gets.
func TestPresenceBeatRepairMatchesTheTemplate(t *testing.T) {
	tmpl := apiRealtimePresenceGo()
	if !strings.Contains(tmpl, presenceBeatNote) {
		t.Fatal("presence.go's template no longer contains presenceBeatNote: update the repair with it")
	}

	old := strings.Replace(tmpl, presenceBeatNote, presenceBeatCall, 1)
	fixed, changes, warnings := repairPresenceBeatContextSource(old)
	if fixed != tmpl || len(changes) != 1 || len(warnings) != 0 {
		t.Fatalf("the repair did not restore the template: %d changes, warnings %v", len(changes), warnings)
	}
	if again, changes, _ := repairPresenceBeatContextSource(fixed); again != fixed || len(changes) != 0 {
		t.Error("a second run changed the file")
	}
}
