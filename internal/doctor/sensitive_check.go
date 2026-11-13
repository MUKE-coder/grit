package doctor

import "strings"

// secretSuffixes are column endings that almost always name a value the API
// should not send back.
//
// Suffixes rather than substrings, which is the whole reason this check can be
// on by default. "token" as a substring matches token_count and "secret"
// matches secret_santa, and two false warnings is all it takes for somebody to
// stop reading the output. A column that ENDS in _token is a token.
var secretSuffixes = []string{
	"password", "passphrase", "secret", "token", "api_key", "apikey",
	"access_key", "private_key", "signing_key", "mnemonic", "seed_phrase",
	"recovery_phrase", "recovery_code", "recovery_codes", "security_answer",
	"otp", "pin_code", "credential", "credentials",
}

// looksLikeASecret reports whether the column name is one of those words, or
// ends in one after an underscore.
func looksLikeASecret(column string) bool {
	column = strings.ToLower(column)
	for _, word := range secretSuffixes {
		if column == word || strings.HasSuffix(column, "_"+word) {
			return true
		}
	}
	return false
}

// A column whose name says it holds a secret, which the API sends back in
// every response.
//
// The rules engine has a list of names it refuses to filter on, and this check
// exists because that list is a guess. It can only know the spellings somebody
// thought of, so a project with a column called wallet_recovery_phrase gets no
// protection from it at all: the name is filterable, and `?filter=
// wallet_recovery_phrase ~ "aband"` reads it back one character at a time
// whether or not the response contains it.
//
// Declaring the column sensitive settles all of it at once. It stops
// resolving in an expression, it leaves the admin table, it leaves the export,
// and the model's MarshalJSON drops the key, so no response carries it.
//
// A warning rather than an error: "ends in _token" is a strong signal and not
// a certainty, and the project may have a column called display_token that is
// genuinely public. Saying so is one modifier away, and saying nothing is
// also allowed.
func checkSensitiveColumns(p *project) []Finding {
	var out []Finding
	for _, name := range p.resources {
		for _, f := range p.models[name].Fields {
			if f.Sensitive || !looksLikeASecret(f.JSON) {
				continue
			}
			// Already encrypted. Not the same promise (encryption is about a
			// database dump, this is about a response), but somebody has
			// clearly thought about this column, and pii-column-not-encrypted
			// is the check that told them to. Warning about the fix we
			// recommended is how a linter gets switched off.
			if strings.Contains(f.Type, "crypto.") {
				continue
			}
			out = append(out, Finding{
				Level:    "warning",
				Resource: name + "." + f.JSON,
				Message:  "is named like a secret and is sent in every response, and an expression can filter on it",
				Fix: "regenerate the field as " + f.JSON + ":string:sensitive, or add `sensitive:\"true\"` to its tag. " +
					"Leave it if the value really is public.",
			})
		}
	}
	return out
}
