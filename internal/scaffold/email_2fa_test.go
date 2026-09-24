package scaffold

import (
	"strings"
	"testing"
)

// A second factor by email, for the people who will not install an app.
//
// The rule that matters throughout: never switch on a factor somebody cannot
// receive. That is why setup is two steps, why a deployment with no mailer is
// refused rather than allowed to lock an account, and why the sign-in says so
// out loud if mail stops working later.
func TestEmailSecondFactorIsTwoSteps(t *testing.T) {
	totp := totpHandlerGo()

	if !strings.Contains(totp, "func (h *TOTPHandler) SendEmailSetupCode(") {
		t.Error("no way to ask for a setup code, so email can only be turned on blind")
	}
	if !strings.Contains(totp, "func (h *TOTPHandler) EnableEmail(") {
		t.Error("no way to confirm the setup code")
	}
	if !strings.Contains(totp, "This deployment cannot send email yet") {
		t.Error("a deployment with no mailer can turn on a factor it cannot deliver")
	}
	// Enabled is only committed once the code comes back.
	send := totp[strings.Index(totp, "SendEmailSetupCode("):]
	send = send[:strings.Index(send, "func (h *TOTPHandler) EnableEmail(")]
	if strings.Contains(send, "config.Enabled = true") {
		t.Error("asking for a code enables the factor before the code is confirmed")
	}
}

// The code is the whole factor, so it is generated and compared like one.
func TestEmailCodesAreGeneratedAndComparedProperly(t *testing.T) {
	lib := totpServiceGo()

	if !strings.Contains(lib, "rand.Int(rand.Reader") {
		t.Error("the sign-in code is not from crypto/rand, so it may be predictable")
	}
	if !strings.Contains(lib, "subtle.ConstantTimeCompare") {
		t.Error("the code comparison is not constant time")
	}
	if !strings.Contains(lib, `fmt.Sprintf("%06d"`) {
		t.Error("the code is not padded to six digits, so some codes would be shorter than others")
	}
}

// Stored hashed, spent on use, and counted against the same attempt limit as
// an authenticator code.
func TestEmailCodesAreStoredHashedAndSpent(t *testing.T) {
	models := twoFactorModelsGo()
	if !strings.Contains(models, "CodeHash") {
		t.Error("the pending token has nowhere to hold an emailed code")
	}
	if strings.Contains(models, "Code string") {
		t.Error("the emailed code is stored in the clear")
	}

	handler := totpHandlerGo()
	if !strings.Contains(handler, "totp.EmailCodeMatches(pending.CodeHash, req.Code)") {
		t.Error("verifying does not check the emailed code")
	}
	if !strings.Contains(handler, `Update("code_hash", "")`) {
		t.Error("an emailed code is not spent when it is used, so it works twice")
	}
	if !strings.Contains(handler, "h.failSecondFactor(pending, user)") {
		t.Error("a wrong emailed code does not count against the attempt limit")
	}
}

// The sign-in has to know which factor to ask for, and say so.
func TestSignInAsksForTheRightFactor(t *testing.T) {
	login := loginBlockNew
	if !strings.Contains(login, "models.TwoFactorMethodEmail") {
		t.Error("the sign-in challenge does not know about codes by email")
	}
	if !strings.Contains(login, `"method": totpConfig.Method`) {
		t.Error("the sign-in response does not say which factor it asked for")
	}
	if !strings.Contains(login, "no mail configured") {
		t.Error("a sign-in with mail switched off leaves somebody at a code box waiting for nothing")
	}

	page := adminThemedLoginPage()
	if !strings.Contains(page, `factorMethod === "email"`) {
		t.Error("the sign-in page tells everybody to open an authenticator, including the people whose code is in their inbox")
	}
}

// An existing project has its own routes.go and auth.go, so both halves need a
// repair or the feature arrives half-wired: a button that 500s, or a factor
// that can be turned on and never asked for.
func TestEmailSecondFactorReachesExistingProjects(t *testing.T) {
	routes := emailTwoFactorAnchor + "\n" + emailTwoFactorHandlerOld
	out, changes, warnings := repairEmailTwoFactorRoutesSource(routes)
	if len(warnings) != 0 {
		t.Fatalf("repairing routes warned: %v", warnings)
	}
	if len(changes) != 2 {
		t.Errorf("expected the routes and the handler wiring, got %v", changes)
	}
	if !strings.Contains(out, "totpHandler.SendEmailSetupCode") || !strings.Contains(out, "Mailer: svc.Mailer") {
		t.Error("the repair left the routes or the mailer out")
	}
	if again, changes, _ := repairEmailTwoFactorRoutesSource(out); again != out || len(changes) != 0 {
		t.Error("the routes repair is not idempotent")
	}

	// And the sign-in half.
	old := strings.Replace(loginBlockNew, emailChallengeNew, emailChallengeOld, 1)
	old = strings.Replace(old, emailChallengeTokenNew, emailChallengeTokenOld, 1)
	if old == loginBlockNew {
		t.Fatal("the fixture is not older than the template")
	}
	out, changes, warnings = repairEmailTwoFactorLoginSource(old)
	if len(warnings) != 0 || len(changes) != 1 {
		t.Fatalf("repairing the sign-in: changes %v, warnings %v", changes, warnings)
	}
	if !strings.Contains(out, "models.TwoFactorMethodEmail") {
		t.Error("the sign-in still cannot ask for a code by email")
	}
	if again, changes, _ := repairEmailTwoFactorLoginSource(out); again != out || len(changes) != 0 {
		t.Error("the sign-in repair is not idempotent")
	}
}

// The email itself: a code to type, never a link to click.
func TestTheCodeEmailHasNoLink(t *testing.T) {
	mail := mailTemplatesGo()
	start := strings.Index(mail, "const twoFactorCodeTemplate = ")
	if start < 0 {
		t.Fatal("there is no template for the code email")
	}
	body := mail[start:]
	if end := strings.Index(body[10:], "\nconst "); end > 0 {
		body = body[:end+10]
	}
	if strings.Contains(body, "<a ") || strings.Contains(body, "href=") {
		t.Error("the code email contains a link, which is the shape of every phishing message ever written")
	}
	if !strings.Contains(body, "{{.Code}}") {
		t.Error("the code email does not contain the code")
	}
	if !strings.Contains(mail, `"two-factor-code":    twoFactorCodeTemplate`) {
		t.Error("the template is not registered, so sending it renders nothing")
	}
}

// The desktop sign-in has to learn which factor an account uses, the same way
// the admin does.
//
// It shipped with the state and the two messages that read it, and nothing
// calling the setter: every account using codes by email was told to open an
// authenticator app it had never set up. Only a generated project's tsc caught
// it, as an unused variable, which is a long way from what was actually wrong.
func TestEverySignInSaysWhereTheCodeCameFrom(t *testing.T) {
	for name, page := range map[string]string{
		"the admin":   adminThemedLoginPage(),
		"the desktop": desktopClientLoginRoute(),
	} {
		if !strings.Contains(page, "factorMethod") {
			t.Errorf("%s sign-in does not know which factor the account uses", name)
			continue
		}
		if !strings.Contains(page, "setFactorMethod(") {
			t.Errorf("%s sign-in never reads the method off the challenge, so an emailed code is announced as an authenticator app", name)
		}
		if !strings.Contains(page, `factorMethod === "email"`) {
			t.Errorf("%s sign-in does not say the code is coming by email", name)
		}
	}
}
