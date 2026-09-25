package scaffold

import (
	"strings"
	"testing"
)

// A passkey you register has to be able to sign you in.
//
// The API has had /auth/passkeys/login/begin and /finish since passkeys
// shipped, and lib/webauthn.ts has had toRequestOptions and encodeAssertion,
// written for exactly this. No frontend ever called any of them, so the
// account page invited you to "add one and you can sign in on this device
// without a password" and then offered no way to do it.
func TestPasskeyLoginIsWiredUp(t *testing.T) {
	hooks := adminUseAuth()

	if !strings.Contains(hooks, "export function usePasskeyLogin()") {
		t.Fatal("nothing calls the passkey login endpoints")
	}
	for _, part := range []string{
		`"/api/auth/passkeys/login/begin"`,
		`"/api/auth/passkeys/login/finish?session=" + encodeURIComponent(begun.data.session_id)`,
		"navigator.credentials.get({",
		"toRequestOptions(begun.data.options)",
		"encodeAssertion(assertion)",
	} {
		if !strings.Contains(hooks, part) {
			t.Errorf("the passkey login hook is missing %s", part)
		}
	}
	// The helpers it needs are imported, or the file does not compile.
	if !strings.Contains(hooks, `from "@/lib/webauthn"`) {
		t.Error("the webauthn helpers are not imported")
	}
	// Same landing as a password sign-in: the server issues the same tokens
	// and records the same session row, so nothing downstream should branch.
	if !strings.Contains(hooks, `router.push(data.data.user.role === "USER" ? "/profile" : "/dashboard")`) {
		t.Error("a passkey sign-in does not land where a password sign-in lands")
	}
}

// And the sign-in page offers it, with a fingerprint on it.
func TestSignInPageOffersThePasskey(t *testing.T) {
	page := adminThemedLoginPage()

	if !strings.Contains(page, "Sign in with a passkey") {
		t.Fatal("the sign-in page does not offer a passkey")
	}
	if !strings.Contains(page, "<Fingerprint className=") {
		t.Error("the passkey button has no fingerprint icon")
	}
	// An icon in the admin's icon map is not automatically a named export.
	if !strings.Contains(adminIconMap(), "  Fingerprint,") {
		t.Error("Fingerprint is not exported from @/lib/icons, so the import fails at build time")
	}

	// Shown only where the browser can produce one. The server is never asked,
	// because before anybody identifies themselves it does not know whether
	// this person has a passkey, and asking would tell an attacker which
	// addresses have accounts.
	if !strings.Contains(page, "passkeysSupported().then((ok)") {
		t.Error("the button is not gated on the browser's own capability check")
	}
	if !strings.Contains(page, "{passkeyReady && (") {
		t.Error("the button renders even where no authenticator exists")
	}

	// Cancelling is not an error. Telling somebody "passkey sign-in failed"
	// because they closed the sheet is how a button stops being trusted.
	if !strings.Contains(page, `?.name === "NotAllowedError") return;`) {
		t.Error("cancelling the passkey prompt is reported as a failure")
	}
}
