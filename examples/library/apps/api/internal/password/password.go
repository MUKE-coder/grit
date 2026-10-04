package password

import (
	"strings"
	"unicode"
)

// MinLength is the floor. Length beats composition for real-world strength, so
// this is the rule that carries the most weight here.
const MinLength = 8

// Rule is one thing a password has to satisfy, named the way it is shown.
type Rule struct {
	// ID is what the admin's checklist matches on, so the wording can change
	// on either side without the two drifting apart.
	ID string `json:"id"`
	// Label is the sentence a person reads.
	Label string `json:"label"`
	// Required marks the rules that block a save. All four do today; the field
	// exists so an advisory rule can be added without a second mechanism.
	Required bool `json:"required"`
	// Clause is the same rule as a piece of a sentence, for the error the API
	// returns. Label is a checklist item and reads as one; dropped into
	// "Your password ..." it produces "Your password needs not a password
	// everyone tries first", which is the sort of thing nobody notices until
	// a user quotes it back. Kept out of the JSON because the checklist has
	// no use for it.
	Clause string `json:"-"`
}

// Rules is the list, in the order the checklist shows them.
func Rules() []Rule {
	return []Rule{
		{ID: "length", Label: "At least 8 characters", Required: true,
			Clause: "needs at least 8 characters"},
		{ID: "variety", Label: "Letters and something else: a number, a symbol or a space", Required: true,
			Clause: "needs letters and something else: a number, a symbol or a space"},
		{ID: "not-common", Label: "Not a password everyone tries first", Required: true,
			Clause: "cannot be one of the passwords everyone tries first"},
		{ID: "not-personal", Label: "Nothing from your name or email", Required: true,
			Clause: "cannot contain your name or your email address"},
	}
}

// Check returns the ids of the rules this password fails.
//
// about is what the account already tells us: an email, a name. Passing it is
// what makes "nothing from your name" possible, and passing nothing is fine at
// registration time when there is nothing to compare against yet.
func Check(candidate string, about ...string) []string {
	var failed []string

	if len([]rune(candidate)) < MinLength {
		failed = append(failed, "length")
	}

	var letters, others bool
	for _, r := range candidate {
		switch {
		case unicode.IsLetter(r):
			letters = true
		default:
			// A digit, a symbol, a space. Anything that is not a letter.
			others = true
		}
	}
	if !letters || !others {
		failed = append(failed, "variety")
	}

	if IsCommon(candidate) {
		failed = append(failed, "not-common")
	}

	if usesPersonal(candidate, about) {
		failed = append(failed, "not-personal")
	}

	return failed
}

// Valid is Check with the ids thrown away.
func Valid(candidate string, about ...string) bool {
	return len(Check(candidate, about...)) == 0
}

// Message turns the failures into one sentence for an API error.
//
// Built from each rule's Clause, in the order Rules lists them rather than the
// order they happened to fail, so the same two failures always read the same
// way.
func Message(failed []string) string {
	broken := map[string]bool{}
	for _, id := range failed {
		broken[id] = true
	}
	parts := make([]string, 0, len(failed))
	for _, rule := range Rules() {
		if broken[rule.ID] && rule.Clause != "" {
			parts = append(parts, rule.Clause)
		}
	}
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return "Your password " + parts[0] + "."
	default:
		// A serial comma before the last clause, because two of the clauses
		// contain an "and" of their own and the sentence stops parsing
		// without it.
		return "Your password " + strings.Join(parts[:len(parts)-1], ", ") + ", and " + parts[len(parts)-1] + "."
	}
}

// usesPersonal reports whether the password contains something the account
// already publishes about the person.
//
// Compared case-insensitively and only for pieces of four characters or more:
// a surname like "Ng" inside an otherwise good password is a coincidence, and
// rejecting it would be a rule nobody could satisfy.
func usesPersonal(candidate string, about []string) bool {
	lower := strings.ToLower(candidate)
	for _, raw := range about {
		for _, piece := range personalPieces(raw) {
			if len(piece) >= 4 && strings.Contains(lower, piece) {
				return true
			}
		}
	}
	return false
}

// personalPieces breaks an email or a name into the parts worth comparing.
func personalPieces(raw string) []string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return nil
	}
	if at := strings.Index(raw, "@"); at > 0 {
		// The local part, and the domain without its suffix: somebody at
		// acme.com should not be using "acme" either.
		raw = raw[:at] + " " + strings.SplitN(raw[at+1:], ".", 2)[0]
	}
	pieces := strings.FieldsFunc(raw, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	return pieces
}

// IsCommon reports whether this is one of the passwords attackers try first.
//
// Deliberately a short list rather than a ten-million-line corpus: the long
// lists belong in a service, and the value of the first few hundred is most of
// the value there is. Everything here is normalised, so Password1 and
// PASSWORD1 are the same entry.
func IsCommon(candidate string) bool {
	_, found := commonPasswords[strings.ToLower(strings.TrimSpace(candidate))]
	return found
}

var commonPasswords = func() map[string]struct{} {
	list := []string{
		"123456", "password", "12345678", "qwerty", "123456789", "12345", "1234", "111111",
		"1234567", "dragon", "123123", "baseball", "abc123", "football", "monkey", "letmein",
		"shadow", "master", "666666", "qwertyuiop", "123321", "mustang", "1234567890",
		"michael", "654321", "superman", "1qaz2wsx", "7777777", "121212", "000000",
		"qazwsx", "123qwe", "killer", "trustno1", "jordan", "jennifer", "zxcvbnm", "asdfgh",
		"hunter", "buster", "soccer", "harley", "batman", "andrew", "tigger", "sunshine",
		"iloveyou", "2000", "charlie", "robert", "thomas", "hockey", "ranger", "daniel",
		"starwars", "klaster", "112233", "george", "computer", "michelle", "jessica",
		"pepper", "1111", "zxcvbn", "555555", "11111111", "131313", "freedom", "777777",
		"pass", "maggie", "159753", "aaaaaa", "ginger", "princess", "joshua", "cheese",
		"amanda", "summer", "love", "ashley", "nicole", "chelsea", "biteme", "matthew",
		"access", "yankees", "987654321", "dallas", "austin", "thunder", "taylor", "matrix",
		"william", "corvette", "hello", "martin", "heather", "secret", "merlin", "diamond",
		"1234qwer", "gfhjkm", "hammer", "silver", "222222", "88888888", "anthony", "justin",
		"test", "bailey", "q1w2e3r4t5", "patrick", "internet", "scooter", "orange", "11111",
		"golfer", "cookie", "richard", "samantha", "bigdog", "guitar", "jackson", "whatever",
		"mickey", "chicken", "sparky", "snoopy", "maverick", "phoenix", "camaro", "sexy",
		"peanut", "morgan", "welcome", "falcon", "cowboy", "ferrari", "samsung", "andrea",
		"smokey", "steelers", "joseph", "mercedes", "dakota", "arsenal", "eagles", "melissa",
		"boomer", "booboo", "spider", "nascar", "monster", "tigers", "yellow", "xxxxxx",
		"123123123", "gateway", "marina", "diablo", "bulldog", "qwer1234", "compaq",
		"purple", "hardcore", "banana", "junior", "hannah", "123654", "porsche", "lakers",
		"iceman", "money", "cowboys", "987654", "london", "tennis", "999999", "ncc1701",
		"coffee", "scooby", "0000", "miller", "boston", "q1w2e3r4", "fuckoff", "brandon",
		"yamaha", "chester", "mother", "forever", "johnny", "edward", "333333", "oliver",
		"redsox", "player", "nikita", "knight", "fender", "barney", "midnight", "please",
		"brandy", "chicago", "badboy", "iwantu", "slayer", "rangers", "charles", "angel",
		"flower", "bigdaddy", "rabbit", "wizard", "bigdick", "jasper", "enter", "rachel",
		"chris", "steven", "winner", "adidas", "victoria", "natasha", "1q2w3e4r", "jasmine",
		"winter", "prince", "panties", "marine", "ghbdtn", "fishing", "cocacola", "casper",
		"james", "232323", "raiders", "888888", "marlboro", "gandalf", "asdfasdf", "crystal",
		"87654321", "12344321", "sexsex", "golden", "blowme", "bigtits", "8675309", "panther",
		"lauren", "angela", "bitch", "spanky", "thx1138", "angels", "madison", "winston",
		"shannon", "mike", "toyota", "blowjob", "jordan23", "canada", "sophie", "Password",
		"apples", "dick", "tiger", "razz", "123abc", "pokemon", "qazxsw", "55555", "qwaszx",
		"muffin", "johnson", "murphy", "cooper", "jonathan", "liverpoo", "david", "danielle",
		"golf", "green", "123456a", "honda", "network", "xxxxxxxx", "admin", "letmein1",
		"password1", "password123", "passw0rd", "p@ssw0rd", "qwerty123", "iloveyou1",
		"welcome1", "admin123", "root", "toor", "changeme", "secret123", "trustno1!",
	}
	set := make(map[string]struct{}, len(list))
	for _, entry := range list {
		set[strings.ToLower(entry)] = struct{}{}
	}
	return set
}()
