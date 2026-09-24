package scaffold

import (
	"strings"
	"testing"
)

// Three things a person hits on their first afternoon with a generated
// project, each of which used to end with them stuck and no way forward.

// Verifying your own email in development.
//
// The banner says "confirm your email" and the only way to do it is to open
// the server log or dig a .html file out of storage/mail. Nobody does that, so
// the banner sits there for the life of every dev project and stops meaning
// anything. In development the endpoint hands the link back with the response.
func TestDevelopmentHandsBackTheVerificationLink(t *testing.T) {
	handler := apiAuthEmailVerificationGo()

	if !strings.Contains(handler, "func (h *AuthHandler) deliverVerificationEmail(ctx context.Context, user models.User) string {") {
		t.Fatal("deliverVerificationEmail does not return the link, so nothing can offer it")
	}
	if !strings.Contains(handler, `body["verify_url"] = link`) {
		t.Error("the send endpoint never returns the link, so development still has to read the log")
	}

	// And only in development. Anywhere else this hands a working verification
	// link to whoever asked for one, which is the entire secret.
	guard := `if h.Config.AppEnv == "development" && link != ""`
	if !strings.Contains(handler, guard) {
		t.Error("the link is not guarded on APP_ENV, so production would hand it out too")
	}
	idx := strings.Index(handler, guard)
	if used := strings.Index(handler, `body["verify_url"] = link`); used < idx {
		t.Error("the link is set before the environment is checked")
	}

	// Every exit returns something, or the caller gets a link for a token that
	// was never stored.
	for _, failure := range []string{
		"generating token for %s: %v\", user.Email, err)\n\t\treturn \"\"",
		"storing token for %s: %v\", user.Email, err)\n\t\treturn \"\"",
	} {
		if !strings.Contains(handler, failure) {
			t.Errorf("a failure path does not return an empty link: %s", failure)
		}
	}
}

// The register path does not want the link, and says so rather than letting a
// silently dropped return value look like an oversight.
func TestRegisterIgnoresTheVerificationLinkOnPurpose(t *testing.T) {
	if !strings.Contains(apiAuthHandlerGo(), "_ = h.deliverVerificationEmail(c.Request.Context(), user)") {
		t.Error("register does not acknowledge the returned link")
	}
}

// "Authenticator is not working."
//
// It almost always is. The phone's clock has drifted past the 30 second window
// and the server says "Invalid verification code. Make sure your authenticator
// app is synced.", which sends people to delete the entry and re-scan a QR
// that was never the problem. The server can measure the drift, so it says how
// far out the clock is and which way.
func TestARefusedCodeSaysWhenTheClockIsTheProblem(t *testing.T) {
	lib := totpServiceGo()

	if !strings.Contains(lib, "func SkewSeconds(secret, code string) (int, bool) {") {
		t.Fatal("nothing measures how far out a refused code is")
	}
	// The measurement is not an acceptance: the steps the validator already
	// tried are skipped, and the caller still refuses the code.
	if !strings.Contains(lib, "continue // already tried and refused by ValidateCodeStep") {
		t.Error("SkewSeconds re-checks the accepted window, so it could be mistaken for a wider one")
	}
	if !strings.Contains(lib, "SkewSearch = 10") {
		t.Error("no bound on the search, so a refusal could cost an unbounded number of HMACs")
	}
	if !strings.Contains(lib, "Window = 1") {
		t.Error("the accepted window changed; the skew search assumes it is still one step")
	}

	handler := totpHandlerGo()
	if !strings.Contains(handler, "func totpRefusal(secret, code string) string {") {
		t.Fatal("the handler has no refusal message builder")
	}
	if !strings.Contains(handler, "respond.Fail(c, respond.CodeInvalidTOTPCode, totpRefusal(req.Secret, req.Code))") {
		t.Error("enable still returns the flat 'invalid code' message")
	}
	if strings.Contains(handler, `"Invalid verification code. Make sure your authenticator app is synced."`) {
		t.Error("the old message that blames the app rather than the clock is still shipped")
	}
	if !strings.Contains(handler, "clock is about %d seconds %s this server's") {
		t.Error("the refusal never says how far out the clock is")
	}
	// A code that is simply wrong must not be told its clock is wrong.
	if !strings.Contains(handler, `return "Invalid verification code. Check you are reading the code for this account`) {
		t.Error("a genuinely wrong code has no message of its own")
	}
	// fmt.Sprintf is used, so the file has to import it.
	if !strings.Contains(handler, "\t\"fmt\"\n") {
		t.Error("the handler formats a message without importing fmt, so the project will not compile")
	}
}

// A refused password has to say why in a sentence somebody can read.
//
// The checklist labels are noun phrases. Dropped after "Your password needs "
// they produced "Your password needs not a password everyone tries first and
// nothing from your name or email.", which is what the API returned for
// months. Each rule now carries its own clause.
func TestARefusedPasswordExplainsItselfInEnglish(t *testing.T) {
	rules := apiPasswordRulesGo()

	if !strings.Contains(rules, "Clause string ") {
		t.Fatal("rules have no sentence form, so the message is built from checklist labels")
	}
	for _, clause := range []string{
		`Clause: "needs at least 8 characters"`,
		`Clause: "cannot be one of the passwords everyone tries first"`,
		`Clause: "cannot contain your name or your email address"`,
	} {
		if !strings.Contains(rules, clause) {
			t.Errorf("missing clause: %s", clause)
		}
	}
	if strings.Contains(rules, `return "Your password needs " + parts[0] + "."`) {
		t.Error("the sentence is still assembled with a fixed 'needs', which only fits two of the four rules")
	}
	if !strings.Contains(rules, `return "Your password " + parts[0] + "."`) {
		t.Error("the single-failure sentence does not come from the clause")
	}
	// Ordered by the rule list, not by the order they failed, so two failures
	// always read the same way.
	if !strings.Contains(rules, "for _, rule := range Rules() {\n\t\tif broken[rule.ID]") {
		t.Error("the clauses are not emitted in rule order, so the message varies between requests")
	}
	// The clause is for the API error only; the checklist has no use for it.
	if !strings.Contains(rules, "Clause string ") || !strings.Contains(rules, `json:"-"`) {
		t.Error("Clause is serialised to the checklist, which does not want it")
	}
}

// The banner can be put away.
//
// It is the first thing on every page of a project whose mail is not wired up
// yet, which is every project on day one. Dismissing it for the session is not
// dismissing the requirement: it comes back at the next sign-in.
func TestTheVerifyEmailBannerCanBeDismissed(t *testing.T) {
	banner := adminEmailVerifiedBanner()

	if !strings.Contains(banner, "sessionStorage") {
		t.Fatal("the dismissal is not stored, so it cannot survive a route change")
	}
	if strings.Contains(banner, "localStorage") {
		t.Error("dismissed in localStorage would outlive the sign-in, and the banner would never return")
	}
	if !strings.Contains(banner, `aria-label="Dismiss until next sign-in"`) {
		t.Error("the dismiss control has no accessible name, so it is an unlabelled button to a screen reader")
	}
	// Storage throws in a private window and in an iframe with cookies
	// blocked. A banner is not worth a blank page.
	if !strings.Contains(banner, "try {") || !strings.Contains(banner, "catch") {
		t.Error("storage is read without a guard, so a private window would crash the layout")
	}
	// And in development it offers the link rather than telling you to check
	// an inbox that will never receive anything.
	if !strings.Contains(banner, "verify_url") {
		t.Error("the banner ignores the link development hands back")
	}
}
