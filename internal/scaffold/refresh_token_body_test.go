package scaffold

import (
	"strings"
	"testing"
)

// Review 2.10: the refresh token went to every client in the JSON body, beside
// the HttpOnly cookie carrying the same value, while the docs said scripts
// never see it.
//
// Three halves to the fix, and this is the test that holds all three in one
// place, because they only work together: the API omits it for browsers, the
// native clients ask for it, and the single frontend stops keeping it where a
// script can read it.
func TestTheRefreshTokenIsNotHandedToBrowsers(t *testing.T) {
	auth := strings.ReplaceAll(apiAuthServiceGo(), "{{MODULE}}", "example.com/app")

	for _, want := range []string{
		// Omitted rather than empty, so a browser's response has no such key.
		"RefreshToken string `json:\"refresh_token,omitempty\"`",
		"func (t TokenPair) ForClient(c *gin.Context) TokenPair {",
		"func WantsTokenInBody(c *gin.Context) bool {",
		"func RequestIsBrowser(c *gin.Context) bool {",
		// A native client says so, because the desktop app is a webview and
		// looks exactly like a browser.
		`const NativeClientHeader = "X-Grit-Client"`,
		`"Sec-Fetch-Site", "Sec-Fetch-Mode", "Sec-Fetch-Dest", "Origin"`,
	} {
		if !strings.Contains(auth, want) {
			t.Errorf("the auth service is missing %q", want)
		}
	}

	// Every handler that returns a pair goes through ForClient. One that does
	// not is a route that still hands the token to a browser, and there is no
	// other way to tell which.
	handlers := map[string]string{
		"auth.go":       apiAuthHandlerGo(),
		"magic_link.go": magicLinkHandlerGo("example.com/app"),
		"passkey.go":    passkeyHandlerGo(),
		"totp.go":       totpHandlerGo(),
	}
	for name, src := range handlers {
		src = strings.ReplaceAll(src, "{{MODULE}}", "example.com/app")
		if n := strings.Count(src, `"tokens": tokens,`); n != 0 {
			t.Errorf("%s returns the raw token pair at %d place(s): use tokens.ForClient(c)", name, n)
		}
		if !strings.Contains(src, "tokens.ForClient(c)") {
			t.Errorf("%s never calls ForClient, so it either returns no tokens or returns them unfiltered", name)
		}
	}
}

// The clients that keep the token somewhere better than a cookie say so.
func TestNativeClientsDeclareThemselves(t *testing.T) {
	for name, src := range map[string]string{
		"the Expo client":    expoAPIClient(),
		"the desktop client": desktopClientApiClientTS(),
	} {
		if !strings.Contains(src, `"X-Grit-Client": "native"`) {
			t.Errorf("%s does not ask for the refresh token, so it will sign a user out when its access token expires", name)
		}
	}
}

// The single architecture's frontend kept the refresh token in localStorage,
// which is the worse half of the same finding: the credential the HttpOnly
// cookie exists to hide, put somewhere every script on the page can read.
func TestTheSingleFrontendDoesNotStoreTheRefreshToken(t *testing.T) {
	src := singleAuthLib() + viteAPIClientWithAuth()

	if strings.Contains(src, "localStorage.setItem(REFRESH_KEY") {
		t.Error("the single frontend still writes the refresh token to localStorage")
	}
	if strings.Contains(src, "localStorage.getItem(REFRESH_KEY)") {
		t.Error("the single frontend still reads the refresh token from localStorage")
	}
	// It is still removed, so an app upgraded from an older version clears the
	// copy it made before.
	if !strings.Contains(src, "localStorage.removeItem(REFRESH_KEY)") {
		t.Error("nothing clears a refresh token an older version of this app stored")
	}
	// And the refresh call relies on the cookie.
	if !strings.Contains(src, `axios.post(API_URL + "/api/auth/refresh", {}, { withCredentials: true })`) {
		t.Error("the refresh interceptor does not send the cookie")
	}
	if strings.Contains(src, "refresh_token: rt") || strings.Contains(src, "refresh_token: refreshToken") {
		t.Error("the single frontend still sends the refresh token in a request body it read from storage")
	}
}
