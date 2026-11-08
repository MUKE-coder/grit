package scaffold

import (
	"strings"
	"testing"
)

// The verification token is single-use, and React 18 StrictMode mounts an
// effect twice in development. The page used to POST twice: the first request
// verified the address, the second answered "invalid or expired", and that is
// the screen the reader saw for a verification that had worked.
func TestVerifyEmailPageSpendsTheTokenOnce(t *testing.T) {
	src := adminVerifyEmailPage()

	if !strings.Contains(src, "useRef") {
		t.Fatal("the verify page has no ref, so nothing stops the effect spending " +
			"the token twice under StrictMode")
	}
	if !strings.Contains(src, "if (sent.current === token) return;") {
		t.Error("the verify page does not skip a token it already sent")
	}
	if strings.Contains(src, "if (cancelled)") {
		t.Error("the verify page gates its state update on a cancellation flag; " +
			"StrictMode cancels the first run and blocks the second, leaving it " +
			"on \"Confirming your email\" after the address was verified")
	}
	if !strings.Contains(src, `apiClient.get<{ data?: { email_verified_at?: string | null } }>("/api/auth/me")`) {
		t.Error("a spent token is reported as a failure even when the signed-in " +
			"account is already verified; the page should confirm success there")
	}
}
