package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

// The profile repair writes a call to h.users(), so it has to add that method to
// a handler that has not got one, and only to such a handler: an accessor nobody
// calls is a method the project's own linter reports as unused.
func TestProfileRepairAddsTheUsersAccessor(t *testing.T) {
	fresh := apiUserHandlerGo()
	if out, ok := withUsersAccessor(fresh); !ok || out != fresh {
		t.Error("the current template already has the accessor and should be left as it is")
	}

	// A handler from before v3.286.0: no CurrentPassword, so the profile repair
	// fires, and no accessor, so it has to come with it.
	old := strings.ReplaceAll(userHandlerFixture(t, "v3.339.0.go.txt"), "{{MODULE}}", "example.com/app")
	old = strings.Replace(old, profileRequestNew, profileRequestAnchor, 1)
	old = strings.Replace(old, profileEmailNew, profileEmailOld, 1)
	old = strings.Replace(old, profileRevokeNew, profileRevokeOld, 1)
	old = strings.Replace(old, "\t\"strings\"\n", "", 1)
	if strings.Contains(old, "h.users()") || strings.Contains(old, "CurrentPassword string") {
		t.Fatal("could not rebuild a handler from before the profile checks")
	}

	out, fixed, warn := repairProfileChangeSource(old)
	if len(warn) != 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v warnings %v", fixed, warn)
	}
	if !strings.Contains(out, "func (h *UserHandler) users() *services.UserService {") {
		t.Fatal("the repaired handler calls h.users() without declaring it")
	}
	if strings.Contains(out, "// Create creates a new user (admin only).\n\nfunc (h *UserHandler) users()") ||
		strings.Contains(out, "// Create creates a new user (admin only).\nfunc (h *UserHandler) users()") {
		t.Error("the accessor was spliced between Create's comment and Create")
	}
	declared := strings.Index(out, "type UserHandler struct {")
	accessor := strings.Index(out, "func (h *UserHandler) users()")
	called := strings.Index(out, "h.users().EmailTakenByAnother(")
	if !(declared < accessor && accessor < called) {
		t.Errorf("the accessor is in the wrong place (type %d, accessor %d, call %d)", declared, accessor, called)
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("the repaired handler is not valid Go: %v", err)
	}
	// Twice is once.
	again, _, _ := repairProfileChangeSource(out)
	if again != out {
		t.Error("a second pass changed it again")
	}
}

// A handler Grit did not write is left alone rather than handed a call to a
// method it has no way to resolve.
func TestProfileRepairLeavesAForeignHandlerAlone(t *testing.T) {
	const theirs = "package handlers\n\nfunc (h *UserHandler) UpdateProfile(c *gin.Context) {\n\tpanic(\"ours\")\n}\n"
	out, fixed, warn := repairProfileChangeSource(theirs)
	if out != theirs || len(fixed) != 0 || len(warn) != 1 {
		t.Errorf("changed a handler Grit did not write: fixed %v warnings %v", fixed, warn)
	}
	if _, ok := withUsersAccessor(theirs); ok {
		t.Error("claimed it could add the accessor to a file with no UserHandler type")
	}
}
