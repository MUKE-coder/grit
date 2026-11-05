package scaffold

import (
	"strings"
	"testing"
)

// grit#174: every generated Expo app had a Google button that could not work.
//
// It asked the API to send it back to a deep link and read the tokens out of
// the query string. The API ignored redirect_uri, finished the handshake, set
// HttpOnly cookies and redirected to the web frontend, so the in-app browser
// landed on a web page, the app never got its tokens, and cookies are no use to
// a bearer client. Three faults compounding, and the feature had never worked.

// The allowlist check has to come BEFORE the handshake starts. Refusing after
// gothic.BeginAuthHandler has already redirected to the provider is not
// refusing: the response is gone.
func TestTheRedirectIsCheckedBeforeTheHandshakeStarts(t *testing.T) {
	src := tmpl("api/handlers/auth_oauth.go")

	begin := strings.Index(src, "func (h *AuthHandler) OAuthBegin")
	if begin < 0 {
		t.Fatal("OAuthBegin is gone, so this test is reading the wrong file")
	}
	body := src[begin:]
	if end := strings.Index(body, "\nfunc "); end > 0 {
		body = body[:end]
	}

	check := strings.Index(body, "NativeRedirectAllowed")
	handshake := strings.Index(body, "gothic.BeginAuthHandler")
	if check < 0 {
		t.Fatal("OAuthBegin does not check redirect_uri against the allowlist, so it is " +
			"an open redirect that carries a sign-in code")
	}
	if handshake < 0 {
		t.Fatal("OAuthBegin no longer starts the handshake")
	}
	if check > handshake {
		t.Error("the allowlist check comes after gothic.BeginAuthHandler. By then the " +
			"response has already redirected to the provider, so refusing does nothing.")
	}
}

// A native client gets a code. Not the tokens, which is the fix that suggests
// itself and the one that puts a refresh token in browser history.
func TestTheNativeCallbackSendsACodeAndNotTokens(t *testing.T) {
	src := tmpl("api/handlers/auth_oauth.go")

	if !strings.Contains(src, "services.IssueOAuthCode") {
		t.Error("the callback does not issue a one-time code, so a native client has " +
			"nothing to exchange")
	}
	if !strings.Contains(src, "services.NativeRedirectWithCode") {
		t.Error("the callback does not build the native redirect")
	}

	// The thing that must never appear: a token in a redirect URL.
	for _, forbidden := range []string{
		"access_token=",
		"refresh_token=",
		"q.Set(\"access_token\"",
		"q.Set(\"refresh_token\"",
	} {
		if strings.Contains(src, forbidden) {
			t.Errorf("the OAuth callback contains %q. A URL lands in browser history, in "+
				"the Referer of the next request and in any proxy log on the way, so a "+
				"refresh token in one is a standing key to the account.", forbidden)
		}
	}
}

// The Expo client has to use the scheme app.json registers, not a hard-coded
// one. "myapp://callback" was a URL the app was not listening on at all, so
// even a correct API would have redirected into nothing.
func TestTheExpoClientUsesItsOwnScheme(t *testing.T) {
	src := expoAuthProvider()

	if strings.Contains(src, "myapp://") {
		t.Error("the Expo client still hard-codes myapp://, which is not the scheme " +
			"app.json registers: that is the project name")
	}
	if !strings.Contains(src, "Linking.createURL(") {
		t.Error("the callback URL is not built from the app's own scheme, so it can " +
			"drift from the scheme the app actually listens on")
	}
	if !strings.Contains(src, `api.post("/auth/oauth/exchange"`) {
		t.Error("the Expo client does not exchange the code for tokens")
	}
	// And it must not go back to reading tokens out of the URL.
	for _, forbidden := range []string{
		`searchParams.get("access_token")`,
		`searchParams.get("refresh_token")`,
	} {
		if strings.Contains(src, forbidden) {
			t.Errorf("the Expo client reads %s out of the callback URL again", forbidden)
		}
	}
	// A failure has to come back through the deep link too: this browser
	// cannot show the app an HTML error page.
	if !strings.Contains(src, `searchParams.get("error")`) {
		t.Error("the Expo client ignores an error the API sends back through the deep link")
	}
}

// Nothing configured means no native sign-in. A project with no mobile app must
// not have an open redirect sitting on its login route by default.
func TestTheAllowlistDefaultsToTheProjectsOwnScheme(t *testing.T) {
	env := envExampleFile(Options{ProjectName: "shop", Architecture: ArchAPI})

	if !strings.Contains(env, "OAUTH_NATIVE_REDIRECTS=shop://callback") {
		t.Error("the allowlist in .env.example does not carry the project's own scheme, " +
			"which is the deep link apps/expo/app.json registers")
	}
	if !strings.Contains(env, "allowlist") {
		t.Error("nothing in .env.example says this is an allowlist, which is the whole " +
			"reason it is not a pattern")
	}
}
