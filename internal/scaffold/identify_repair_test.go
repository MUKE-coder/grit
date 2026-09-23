package scaffold

import (
	"strings"
	"testing"
)

func TestIdentifyRepairAddsItToAnOlderProject(t *testing.T) {
	fresh := apiAuthMiddlewareGo()
	mustFormatGo(t, "auth.go", fresh)

	block := identifyFunc()
	if block == "" || !strings.Contains(block, "func Identify(db *gorm.DB") {
		t.Fatal("the Identify function cannot be lifted out of the template")
	}

	old := strings.Replace(fresh, block, "", 1)
	if old == fresh || strings.Contains(old, "func Identify(") {
		t.Fatal("the fixture still has Identify in it")
	}

	out, changes, warnings := repairIdentifySource(old)
	if len(warnings) != 0 {
		t.Fatalf("repairing an untouched auth.go warned: %v", warnings)
	}
	if out != fresh {
		t.Error("the repaired auth.go is not the one the template writes")
	}
	if len(changes) != 1 {
		t.Errorf("changes = %v, want one", changes)
	}
	mustFormatGo(t, "auth.go", out)

	if again, changes, _ := repairIdentifySource(out); again != out || len(changes) != 0 {
		t.Error("the Identify repair is not idempotent")
	}
}

// auth.go is a file people edit. One that no longer looks like Grit's gets a
// note rather than having its edits overwritten.
func TestIdentifyRepairLeavesAnEditedFileAlone(t *testing.T) {
	edited := "package middleware\n\nfunc Auth() {}\n"
	out, changes, warnings := repairIdentifySource(edited)
	if out != edited || len(changes) != 0 || len(warnings) != 1 {
		t.Errorf("an edited auth.go should be left alone with a note (changes %v, warnings %v)", changes, warnings)
	}
}
