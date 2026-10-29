package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"unicode"
)

// A generated project's tests and benchmarks register users, and registration
// applies the project's own password rules. A fixture that breaks one of them
// does not fail in an obvious way: the request answers 422, and the test
// either fails on a require with no explanation of why, or does not check and
// quietly measures the validation error path instead of the thing it names.
//
// BenchmarkAuthLogin shipped with the password "benchpassword99" and the
// surname "Bench". usesPersonal rejects a password containing any piece of the
// name or address four characters or longer, so registration refused it and
// the benchmark failed on its own require, in every scaffolded project, while
// the register benchmark beside it used a password that passes.
//
// This mirrors the rules in internal/password so a fixture cannot be written
// that the project it is written into will refuse. The rules are duplicated
// rather than imported because that package is a template, not code this
// binary links; if they change there, they change here.
const (
	fixtureMinLength = 8
	fixturePieceMin  = 4
)

// Every map literal in a template that posts a registration.
var fixtureLiteral = regexp.MustCompile(
	`(?s)map\[string\]string\{(.*?)\}`)

var fixtureField = regexp.MustCompile(`"(first_name|last_name|email|password)":\s*"([^"]*)"`)

func TestTemplateTestFixturePasswordsPassTheProjectsOwnRules(t *testing.T) {
	dir := filepath.Join("templates", "api", "handlers")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}

	checked := 0
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), "_test.go.tmpl") {
			continue
		}
		src, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("reading %s: %v", e.Name(), err)
		}

		for _, lit := range fixtureLiteral.FindAllStringSubmatch(string(src), -1) {
			fields := map[string]string{}
			for _, f := range fixtureField.FindAllStringSubmatch(lit[1], -1) {
				fields[f[1]] = f[2]
			}
			password, ok := fields["password"]
			// Only a literal that carries a name is a registration; a login
			// posts an email and a password, and its password is whatever was
			// registered.
			if !ok || fields["first_name"] == "" {
				continue
			}
			checked++

			about := []string{fields["email"], fields["first_name"], fields["last_name"]}
			for _, reason := range fixtureFails(password, about) {
				t.Errorf("%s: the fixture password %q fails the project's own "+
					"rule %q, so registering it answers 422 and the test or "+
					"benchmark measures that instead. about=%v",
					e.Name(), password, reason, about)
			}
		}
	}

	if checked < 3 {
		t.Errorf("found only %d registration fixtures in %s; the reader "+
			"stopped matching and this test is now checking nothing",
			checked, dir)
	}
}

// fixtureFails mirrors password.Check for the rules a fixture can break.
// IsCommon is not mirrored: it reads a word list that lives in the template.
func fixtureFails(candidate string, about []string) []string {
	var failed []string

	if len([]rune(candidate)) < fixtureMinLength {
		failed = append(failed, "length")
	}

	var letters, others bool
	for _, r := range candidate {
		if unicode.IsLetter(r) {
			letters = true
		} else {
			others = true
		}
	}
	if !letters || !others {
		failed = append(failed, "variety")
	}

	lower := strings.ToLower(candidate)
	for _, raw := range about {
		for _, piece := range fixturePieces(raw) {
			if len(piece) >= fixturePieceMin && strings.Contains(lower, piece) {
				failed = append(failed, "not-personal")
				return failed
			}
		}
	}
	return failed
}

// fixturePieces mirrors password.personalPieces: the parts of a name or an
// address worth comparing, which for an email is the local part and the
// domain without its suffix.
func fixturePieces(raw string) []string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return nil
	}
	if at := strings.Index(raw, "@"); at > 0 {
		raw = raw[:at] + " " + strings.SplitN(raw[at+1:], ".", 2)[0]
	}
	return strings.FieldsFunc(raw, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}
