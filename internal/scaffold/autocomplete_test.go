package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Every email field a generated project renders carries an autocomplete
// attribute, and the right one.
//
// grit#180: the password field on a generated sign-in said current-password and
// the email field beside it said nothing, so there was no username field for a
// password manager to anchor to. Chrome says so on a fresh scaffold:
//
//	[DOM] Input elements should have autocomplete attributes (suggested: "username")
//
// Saving is unreliable without the pair and filling often does not happen,
// which is what pushes somebody back to typing a password they reuse.
//
// This is a sweep rather than a list of the pages known today, because the
// cost in that issue was the count: the four admin styles and the themed pages
// each carry their own copy of the same three forms, and the web, desktop and
// single apps carry theirs again, so one attribute is a 28-file change. A test
// naming the pages would pass for a new auth style that forgot it.

// emailInput finds a type="email" and the function that writes it.
var (
	emailInputRe   = regexp.MustCompile(`type="email"`)
	goFuncRe       = regexp.MustCompile(`func (\w+)\(`)
	// Either spelling. React wants autoComplete; plain HTML, which the
	// first-run setup page is because it is served by the Go binary before
	// any frontend has been built, wants autocomplete. A sweep that knows
	// only the JSX casing reports a correct HTML form as a missing attribute,
	// which teaches people to ignore it.
	autoCompleteRe = regexp.MustCompile(`(?i)autocomplete="([a-z-]+)"`)
)

// credentialWords name a form whose email field is the credential's identifier.
var credentialWords = []string{"login", "signin", "signup", "register", "forgot", "reset"}

// profileWords name a form where it is the reader's own address, with a session
// already. `username` there would offer to overwrite the saved credential for
// the site, so the two cases are genuinely different.
var profileWords = []string{"profile", "account", "settings"}

func TestEveryEmailFieldHasAnAutocomplete(t *testing.T) {
	for _, path := range scaffoldSourceFiles(t) {
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading %s: %v", path, err)
		}
		lines := strings.Split(string(body), "\n")
		current := ""

		for i, line := range lines {
			if m := goFuncRe.FindStringSubmatch(line); m != nil {
				current = m[1]
			}
			if !emailInputRe.MatchString(line) {
				continue
			}
			// A Playwright selector is not an input.
			if strings.Contains(line, "input[type=") {
				continue
			}

			// The element, as far as the attribute could reasonably sit.
			lo, hi := i, i+1
			for lo > 0 && !strings.Contains(lines[lo], "<input") {
				lo--
				if i-lo > 8 {
					break
				}
			}
			for hi < len(lines) && !strings.Contains(lines[hi-1], "/>") {
				hi++
				if hi-i > 12 {
					break
				}
			}
			element := strings.Join(lines[lo:hi], "\n")

			name := filepath.Base(path)
			found := autoCompleteRe.FindStringSubmatch(element)
			if found == nil {
				t.Errorf("%s:%d (%s): a type=\"email\" input with no autoComplete, so a "+
					"password manager has no username field to pair with the password:\n%s",
					name, i+1, current, element)
				continue
			}

			want := wantedToken(current)
			if want != "" && found[1] != want {
				t.Errorf("%s:%d (%s): autoComplete=%q, want %q. A credential field is "+
					"\"username\", which is what pairs it with the password; the reader's own "+
					"address on an account page is \"email\".",
					name, i+1, current, found[1], want)
			}
		}
	}
}

// wantedToken classifies a page by the function that writes it. An empty answer
// means this test will not guess: any attribute is accepted, and the one above
// has already insisted there is one.
func wantedToken(fn string) string {
	lower := strings.ToLower(fn)
	// Profile first: singleAccountProfileRoute contains neither credential word,
	// but a name could carry both and the session decides.
	for _, w := range profileWords {
		if strings.Contains(lower, w) {
			return "email"
		}
	}
	for _, w := range credentialWords {
		if strings.Contains(lower, w) {
			return "username"
		}
	}
	return ""
}

// scaffoldSourceFiles is every file in this package that can write a page, plus
// the .tmpl files beside it. Read as text: the attribute is in a string
// literal, and the only question is whether it is there.
func scaffoldSourceFiles(t *testing.T) []string {
	t.Helper()
	var out []string
	for _, root := range []string{".", "templates"} {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			name := info.Name()
			switch {
			case strings.HasSuffix(name, "_test.go"):
				return nil
			case strings.HasSuffix(name, ".go"), strings.HasSuffix(name, ".tmpl"):
				out = append(out, path)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walking %s: %v", root, err)
		}
	}
	if len(out) < 50 {
		t.Fatalf("only %d source files found, so this sweep is not sweeping", len(out))
	}
	return out
}
