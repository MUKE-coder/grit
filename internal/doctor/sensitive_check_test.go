package doctor

import (
	"strings"
	"testing"
)

// The rules engine refuses to filter on a column whose name is on a list it
// carries. The list is a guess, and this check is the prompt to stop relying
// on it: a column named like a secret, that the API returns and an expression
// can read back a character at a time, should say so on the model instead.
func TestAColumnNamedLikeASecretIsFlagged(t *testing.T) {
	p := &project{
		resources: []string{"Wallet"},
		models: map[string]*model{"Wallet": {Name: "Wallet", Fields: []field{
			{Go: "Label", Type: "string", JSON: "label"},
			{Go: "RecoveryPhrase", Type: "string", JSON: "recovery_phrase"},
			{Go: "WebhookSecret", Type: "string", JSON: "webhook_secret"},
			{Go: "RefreshToken", Type: "string", JSON: "refresh_token"},
		}}},
	}

	findings := checkSensitiveColumns(p)
	if len(findings) != 3 {
		t.Fatalf("got %d findings, want 3: %v", len(findings), findings)
	}
	for _, f := range findings {
		if f.Level != "warning" {
			t.Errorf("%s is a %s; the name is a strong signal and not a certainty", f.Resource, f.Level)
		}
		if !strings.Contains(f.Fix, ":sensitive") {
			t.Errorf("the fix for %s does not name the modifier that settles it: %q", f.Resource, f.Fix)
		}
	}
}

// A column that has declared itself is the point of the exercise, and must go
// quiet. It was flagged once, which would have taught people that declaring it
// changes nothing.
func TestADeclaredColumnIsNotFlagged(t *testing.T) {
	p := &project{
		resources: []string{"Wallet"},
		models: map[string]*model{"Wallet": {Name: "Wallet", Fields: []field{
			{Go: "RecoveryPhrase", Type: "string", JSON: "recovery_phrase", Sensitive: true},
		}}},
	}
	if findings := checkSensitiveColumns(p); len(findings) != 0 {
		t.Errorf("a declared column was still flagged: %q", findings[0].Message)
	}
}

// The names that merely contain one of the words. Two false warnings is all it
// takes for somebody to stop reading the output, which is why the match is on
// the end of the name rather than anywhere in it.
func TestTheLookalikesAreNotFlagged(t *testing.T) {
	for name, column := range map[string]string{
		"a count of tokens":   "token_count",
		"an office party":     "secret_santa",
		"when it was issued":  "token_issued_at",
		"a public identifier": "public_key_id",
		"an ordinary word":    "description",
		"a password policy":   "password_policy_version",
	} {
		p := &project{
			resources: []string{"Thing"},
			models: map[string]*model{"Thing": {Name: "Thing", Fields: []field{
				{Go: "X", Type: "string", JSON: column},
			}}},
		}
		if findings := checkSensitiveColumns(p); len(findings) != 0 {
			t.Errorf("%s (%s) was flagged: %q", name, column, findings[0].Message)
		}
	}
}

// An encrypted column is the fix pii-column-not-encrypted recommends. Warning
// about it from a second check teaches people that doing the right thing
// changes nothing.
func TestAnEncryptedColumnIsNotFlagged(t *testing.T) {
	p := &project{
		resources: []string{"Invoice"},
		models: map[string]*model{"Invoice": {Name: "Invoice", Fields: []field{
			{Go: "Secret", Type: "crypto.EncryptedString", JSON: "secret"},
		}}},
	}
	if findings := checkSensitiveColumns(p); len(findings) != 0 {
		t.Errorf("an encrypted column was flagged: %q", findings[0].Message)
	}
}

// And the model parser has to see the tag at all. It reads the line, because
// fieldRe stops at the json tag and never reaches the rest of it: a test of
// the check alone would pass while every real project was flagged.
func TestTheModelParserReadsTheSensitiveTag(t *testing.T) {
	src := "type Wallet struct {\n" +
		"\tLabel string `gorm:\"size:255\" json:\"label\"`\n" +
		"\tRecoveryPhrase string `gorm:\"type:text\" json:\"recovery_phrase\" sensitive:\"true\"`\n" +
		"}\n"
	models := parseModels(src)
	if len(models) != 1 {
		t.Fatalf("got %d models, want 1", len(models))
	}
	byJSON := map[string]field{}
	for _, f := range models[0].Fields {
		byJSON[f.JSON] = f
	}
	if byJSON["label"].Sensitive {
		t.Error("an ordinary column was read as sensitive")
	}
	if !byJSON["recovery_phrase"].Sensitive {
		t.Error("the tag is on the line and the parser did not see it")
	}
}
