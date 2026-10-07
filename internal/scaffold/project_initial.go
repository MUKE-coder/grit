package scaffold

import (
	"strings"
	"unicode"
)

// projectInitial is the letter in a generated app's logo mark.
//
// It was the letter G, in the web navbar, the single's navbar and the desktop
// title bar, beside the project's own name: every app shipped wearing the
// framework's initial. A project called contacts gets a C.
//
// Non-letters are skipped, so "2fa-portal" is an F rather than a 2, and a name
// with no letter at all falls back to the one it used to have.
func projectInitial(opts Options) string {
	for _, r := range opts.ProjectName {
		if unicode.IsLetter(r) {
			return strings.ToUpper(string(r))
		}
	}
	return "G"
}
