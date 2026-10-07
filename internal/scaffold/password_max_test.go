package scaffold

import (
	"strings"
	"testing"
)

// A password bcrypt cannot hash is a validation error, not a 500.
//
// golang.org/x/crypto/bcrypt refuses anything over 72 bytes outright. The
// rules had a minimum and no maximum, so this passphrase
//
//	correct horse battery staple correct horse battery staple correct horse 9!
//
// which is 74 bytes and passes every other rule, reached the model's
// BeforeCreate, failed to hash, and came back as
//
//	HTTP 500 {"code":"INTERNAL_ERROR","message":"Failed to create user"}
//
// Reproduced against a running API before this was written. A long passphrase
// is the strongest kind, so the people turned away were the ones doing the
// right thing.
func TestPasswordRulesHaveABcryptCeiling(t *testing.T) {
	src := apiPasswordRulesGo()

	if !strings.Contains(src, "MaxBytes = 72") {
		t.Error("no ceiling, so a long passphrase reaches bcrypt and the API answers 500")
	}
	// Bytes, not runes: an accented letter is two and an emoji is four, so
	// counting characters puts the 500 back for anybody not writing ASCII.
	if !strings.Contains(src, "len(candidate) > MaxBytes") {
		t.Error("the ceiling counts characters, which lets a short multi-byte password through to bcrypt")
	}
	// The checklist the admin draws is built from Rules(), so a rule that is
	// enforced and not listed is a save that fails with nothing marked red.
	if !strings.Contains(src, `ID: "max-length"`) {
		t.Error("the rule is enforced but not in Rules(), so the checklist cannot show it")
	}
	// And the error says what to do about it.
	if !strings.Contains(src, "cannot be longer than 72 bytes") {
		t.Error("the failure has no clause, so the API's sentence does not mention it")
	}

	// The client says it too, so the message arrives before the round trip.
	shared := sharedUserSchema()
	if !strings.Contains(shared, "TextEncoder().encode(v).length <= 72") {
		t.Error("the Zod schema has no ceiling, so the browser lets it through to a server error")
	}
}
