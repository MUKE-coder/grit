package scaffold

import (
	"strings"
	"testing"
)

// An upgrade repair is a pair of anchors: the shape it expects to find, and the
// shape it writes. The second one has to still be in the templates, because the
// repair uses it to decide it has nothing to do. When it drifts, the repair
// stops recognising a project Grit itself wrote, prints "not the file Grit
// wrote" and returns, and everything after that point in the same function is
// skipped too. That is a failure mode with no test behind it and no symptom
// except a warning people learn to scroll past.
//
// It happened: deliverVerificationEmail grew a string return, the anchors ended
// at the opening brace, and repairQueuedAuthMail warned on every single
// upgrade. This is the test that would have caught it the same day.
func TestThinServicesAnchorsMatchTemplates(t *testing.T) {
	// The three anchors the repair reads as "this project already has it".
	// These are the ones that must still be in the templates. The rewrite
	// targets are not checked against the templates on purpose: an old project
	// passed a models.User value where the template now passes *user, so the
	// text the repair writes is correct for the file it is fixing and is not
	// meant to be identical to what a fresh scaffold writes.
	for _, c := range []struct {
		what, anchor, in string
	}{
		{"the AuthHandler Jobs field", "Jobs *jobs.Client", apiAuthHandlerGo()},
		{"the verification function", verifyFuncNew, apiAuthEmailVerificationGo()},
		{"the password reset function", resetFuncNew, apiAuthPasswordResetGo()},
	} {
		if !strings.Contains(c.in, c.anchor) {
			t.Errorf("repairQueuedAuthMail will not recognise its own output: %s is no longer in the template.\nanchor: %q", c.what, c.anchor)
		}
	}

	// And the shapes it repairs from must NOT be in the templates, or the
	// repair would rewrite a current project back and forth.
	for _, c := range []struct {
		what, anchor, in string
	}{
		{"the old verification signature", verifyFuncOld, apiAuthEmailVerificationGo()},
		{"the old reset signature", resetFuncOld, apiAuthPasswordResetGo()},
		{"the verification goroutine", verifyGoroutineCall, apiAuthHandlerGo()},
		{"the reset goroutine", resetGoroutineCall, apiAuthPasswordResetGo()},
	} {
		if strings.Contains(c.in, c.anchor) {
			t.Errorf("%s is still in the template, so the scaffold ships the thing the repair exists to remove", c.what)
		}
	}
}
