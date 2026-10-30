package scaffold

import (
	"strings"
	"testing"
)

// The web sign-in could not complete for an account with 2FA on. The API
// answers with a pending token and no user, and the page stamped the session
// marker and redirected anyway, so the password was accepted, nothing was
// signed in, and the person landed on the home page as a guest (grit#105).
//
// The admin had handled this for some time. The fix never reached the web
// template, which is the recurring shape of these: one of the two apps learns
// something and the other does not.
func TestWebLoginHandlesTheTwoFactorChallenge(t *testing.T) {
	page := webThemedLoginPage()

	for _, want := range []string{
		"totp_required",
		"pending_token",
		"Two-factor authentication",
		"/api/auth/totp/verify",
		"/api/auth/totp/backup-codes/verify",
		"Use a backup code",
		"Trust this device",
	} {
		if !strings.Contains(page, want) {
			t.Errorf("the web login page is missing %q", want)
		}
	}

	// And the marker is stamped only once there really is a session.
	// Stamping it beside a pending token is what told the middleware
	// somebody was signed in when they were not.
	if strings.Count(page, "setWebSessionMarker()") != 1 {
		t.Error("setWebSessionMarker should be called from one place, after the " +
			"login has actually completed")
	}
}

// The hook has the same hole: data.data.user is undefined on a challenge.
func TestWebUseLoginGuardsTheChallenge(t *testing.T) {
	src := webUseAuth()

	for _, want := range []string{
		"export interface TOTPChallenge",
		"function isChallenge(",
		"if (isChallenge(data.data)) return;",
		"export function useVerifyTOTP()",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the web auth hooks are missing %q", want)
		}
	}
}

// The API's OAuth callback sets the cookies itself and redirects with no
// query parameters. The web page read tokens from the query string, so a
// successful sign-in fell through to "Authentication failed", and the branch
// that would have handled them stored tokens where JavaScript could read them
// (grit#106).
func TestWebOAuthCallbackDoesNotTouchTokens(t *testing.T) {
	src := webAuthCallback()

	// Comments may discuss the old behaviour. Code may not do it.
	var code []string
	for _, line := range strings.Split(src, "\n") {
		if t := strings.TrimSpace(line); !strings.HasPrefix(t, "//") {
			code = append(code, line)
		}
	}
	body := strings.Join(code, "\n")

	for _, forbidden := range []string{
		`searchParams.get("access_token")`,
		`Cookies.set`,
		"js-cookie",
	} {
		if strings.Contains(body, forbidden) {
			t.Errorf("the callback still handles tokens: %q", forbidden)
		}
	}
	if !strings.Contains(body, "setWebSessionMarker()") {
		t.Error("the callback must stamp the web session marker, or proxy.ts " +
			"bounces the next protected page straight back to /login")
	}
	if !strings.Contains(body, `searchParams.get("error")`) {
		t.Error("the callback should still handle ?error=")
	}
}

// "(auth)" is a route group, so it does not appear in the URL. The callback
// lived at app/(auth)/callback and therefore served /callback, while the API
// redirects to /auth/callback: every web OAuth sign-in landed on a 404.
func TestWebOAuthCallbackIsAtTheURLTheAPIRedirectsTo(t *testing.T) {
	var found bool
	for _, f := range webAuthFiles("web", Options{}) {
		p := strings.ReplaceAll(f.path, "\\", "/")
		if strings.HasSuffix(p, "auth/callback/page.tsx") {
			found = true
			if strings.Contains(p, "(auth)/callback/") {
				t.Errorf("the callback is at %s, which serves /callback", p)
			}
		}
	}
	if !found {
		t.Error("no callback page is written at a path serving /auth/callback")
	}
}
