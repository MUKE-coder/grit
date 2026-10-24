package scaffold

import (
	"strings"
	"testing"
)

// The seeder mints a live server-to-server key and is the only moment its
// plaintext exists. Printing it is right in development and wrong everywhere
// else: a log line outlives the terminal, gets shipped to an aggregator, and
// is readable by everyone who can read logs.
//
// The same codebase already makes this call twice. The password reset handler
// refuses to print a reset link outside development, with a comment saying a
// live link in a log file is a credential, and the users seeder refuses the
// development password outside development. This file was the one that still
// printed, which a Scout scan found on 2026-10-24.
func TestAPIKeySeederDoesNotPrintTheSecretOutsideDevelopment(t *testing.T) {
	src := apiAPIKeySeederGo()

	if !strings.Contains(src, "if seedDevDefaults() {") {
		t.Fatal("the seeder prints the secret without checking the environment")
	}

	// The guard has to be around the secret, not the publishable key: one is
	// safe in a browser by design and the other is not.
	secretAt := strings.Index(src, `log.Printf("  Secret       %s", secret)`)
	if secretAt < 0 {
		t.Fatal("could not find the line that prints the secret")
	}
	guardAt := strings.Index(src, "if seedDevDefaults() {")
	if guardAt > secretAt {
		t.Error("the environment check comes after the secret is printed")
	}

	// And the operator has to be told why nothing appeared, or they will
	// assume the seeder failed.
	if !strings.Contains(src, "not printed") {
		t.Error("production says nothing about why the secret is absent")
	}
}
