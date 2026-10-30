package scaffold

import (
	"strings"
	"testing"
)

// The repaired handler is the handler the template writes, byte for byte.
//
// Two things that must agree with nothing comparing them: the template holds
// the wording, the repair holds a copy of it, and a repair that drifts leaves
// an upgraded project checked less than a new one. The template is stripped of
// the pieces this repair adds, repaired, and compared.
func TestOAuthVerifiedRepairMatchesTheTemplate(t *testing.T) {
	fresh := strings.ReplaceAll(apiAuthOAuthGo(), "{{MODULE}}", "example.com/app")

	old := fresh
	for _, piece := range []string{
		oauthVerifiedPredicate,
		oauthVerifiedCompute,
		oauthVerifiedRefusal,
		oauthVerifiedNotice,
	} {
		if strings.Count(old, piece) != 1 {
			t.Fatalf("the template does not hold this piece exactly once, so this test is comparing nothing:\n%s", piece)
		}
		old = strings.Replace(old, piece, "", 1)
	}
	if strings.Count(old, oauthVerifiedNewAccount) != 1 {
		t.Fatal("the template does not hold the new-account block")
	}
	old = strings.Replace(old, oauthVerifiedNewAccount, oauthVerifiedNewAccountOld, 1)

	// The stripped version is the handler as it shipped: it reads the address
	// and never asks whether the provider verified it.
	if strings.Contains(old, "providerVerifiedTheEmail") {
		t.Fatal("stripping left the predicate behind")
	}
	if !strings.Contains(old, "EmailVerifiedAt: &now,") {
		t.Fatal("stripping did not restore the unconditional verified timestamp")
	}

	out, changes, warnings := repairOAuthVerifiedEmailSource(old)
	if len(warnings) != 0 {
		t.Fatalf("the repair refused the handler it is written for: %v", warnings)
	}
	if len(changes) != 1 {
		t.Fatalf("changes %v", changes)
	}
	if out != fresh {
		t.Error("the repaired handler differs from the one the template writes")
	}
	mustBeValidGo(t, out)

	// A fresh handler needs nothing, and running it twice changes nothing.
	if again, changes, warnings := repairOAuthVerifiedEmailSource(fresh); again != fresh || len(changes)+len(warnings) != 0 {
		t.Errorf("a fresh handler still needs the repair (changes %v warnings %v)", changes, warnings)
	}
	if again, _, _ := repairOAuthVerifiedEmailSource(out); again != out {
		t.Error("repairing twice changed it again")
	}
}

// The repair waits for the one before it. This one decides when the password
// may be cleared; adding that condition to a handler with no such block would
// be a condition on nothing, and would leave the predicate in a file that
// never consults it.
func TestOAuthVerifiedRepairWaitsForTheUnverifiedLinkRepair(t *testing.T) {
	src := "package handlers\n\n// OAuthCallback completes the OAuth flow\nfunc x() {}\n"
	out, changes, warnings := repairOAuthVerifiedEmailSource(src)
	if out != src || len(changes)+len(warnings) != 0 {
		t.Errorf("a handler without the unverified-link block was changed (changes %v warnings %v)", changes, warnings)
	}
}

// What the predicate is for, stated as a test so the policy is written down
// somewhere a reader will find it: an address a provider has not verified does
// not reach an account that already exists.
func TestTheVerifiedPolicyIsInTheHandler(t *testing.T) {
	src := strings.ReplaceAll(apiAuthOAuthGo(), "{{MODULE}}", "example.com/app")

	for _, want := range []string{
		// GitHub asserts it by construction, Google says so in the payload,
		// and anything else is not trusted by default.
		`if provider == "github" {`,
		`for _, key := range []string{"email_verified", "verified_email"} {`,
		// A new account is only verified when the provider verified it.
		"if emailVerified {\n\t\t\t\tuser.EmailVerifiedAt = &now\n\t\t\t}",
		// An existing account is not matched on an unverified claim, unless
		// this provider account is already linked to it.
		"if !emailVerified && !alreadyLinked {",
		// And a removed password is notified.
		`"Your password was removed", body)`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the OAuth handler is missing %q", want)
		}
	}

	// The unconditional version is gone.
	if strings.Contains(src, "EmailVerifiedAt: &now,") {
		t.Error("a new account is still marked verified whatever the provider said")
	}
}
