package scaffold

import (
	"os"
	"strings"
	"testing"
)

// The repair turns the old listener into the one the push plugin now writes,
// and a second run changes nothing.
func TestPushWebRepairMatchesThePlugin(t *testing.T) {
	tmpl, err := os.ReadFile("../plugin/push/expo_push.ts.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	current := strings.ReplaceAll(string(tmpl), "\r\n", "\n")
	if !strings.Contains(current, pushTapNew) {
		t.Fatal("the push plugin's template no longer has pushTapNew: update the repair with it")
	}
	old := strings.Replace(current, pushTapNew, pushTapOld, 1)
	fixed, changes, _ := repairPushWebSource(old)
	if fixed != current || len(changes) != 1 {
		t.Fatalf("the repair did not produce the template (%d changes)", len(changes))
	}
	if again, changes, _ := repairPushWebSource(fixed); again != fixed || len(changes) != 0 {
		t.Error("a second run changed the file")
	}
}
