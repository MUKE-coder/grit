package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// A repair's historical anchor must keep the spelling it had, not the current
// one.
//
// A repair finds the code it is fixing by matching what an older project
// actually contains. Those snapshots are frozen by definition: the constants
// named loginBlockOld286, and the inline blocks in token_repair.go, are what
// a project scaffolded at that version has on disk.
//
// v3.546.0 converted 53 handlers from
// respond.Fail(c, respond.CodeValidationError, err.Error()) to
// respond.Invalid(c, err), with a walk over internal/. The walk rewrote two
// anchors along with the real call sites, and an anchor updated to the new
// spelling matches nothing in any project that predates the change: the repair
// reports nothing to do, forever, silently. One was caught by the test that
// feeds the repair an older handler. The other had no such test and was found
// by reading the diff.
//
// So this is the guard for the class rather than for those two. A repair that
// matches on an old spelling has to keep it.
func TestHistoricalAnchorsKeepTheOldSpelling(t *testing.T) {
	// The old spelling, and the repair that has to go on recognising it.
	const oldBind = `respond.Fail(c, respond.CodeValidationError, err.Error())`

	anchors := map[string]string{
		"auth_disclosure_repair.go": "loginBlockOld286",
		"token_repair.go":           "bindOlds",
	}

	for file, holder := range anchors {
		body, err := os.ReadFile(filepath.Join(file))
		if err != nil {
			t.Fatalf("reading %s: %v", file, err)
		}
		src := string(body)

		if !strings.Contains(src, holder) {
			t.Errorf("%s no longer declares %s, so this guard is watching nothing",
				file, holder)
			continue
		}
		if !strings.Contains(src, oldBind) {
			t.Errorf("%s no longer contains the pre-v3.546.0 bind spelling, so its "+
				"repair matches nothing in a project scaffolded before the change "+
				"and reports nothing to do. Anchors are snapshots: leave them as "+
				"they were and add the new spelling beside them", file)
		}
	}
}

// And no template still answers a bind failure with the validator's own text.
//
// The point of grit#225: err.Error() on a validator failure names the Go type
// and the tag, with newlines, and leaves error.details empty, which is the
// field the documentation tells clients to read. respond.Invalid fills it.
//
// Repairs are excluded, because their job is to recognise exactly that string.
func TestNoTemplateAnswersABindFailureWithTheValidatorsOwnText(t *testing.T) {
	const bad = `respond.Fail(c, respond.CodeValidationError, err.Error())`

	var offenders []string
	err := filepath.Walk("..", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		name := info.Name()
		if !strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, ".tmpl") {
			return nil
		}
		// A repair matches the old text on purpose, and a test asserts on it.
		if strings.Contains(name, "repair") || strings.HasSuffix(name, "_test.go") {
			return nil
		}
		body, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if strings.Contains(string(body), bad) {
			offenders = append(offenders, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking the templates: %v", err)
	}

	for _, path := range offenders {
		t.Errorf("%s answers a bind failure with err.Error(), which sends the "+
			"validator's own sentence to the client and leaves error.details "+
			"empty. Use respond.Invalid(c, err)", path)
	}

	// Guard the guard: the walk has to be reading the handlers at all.
	found := 0
	_ = filepath.Walk("..", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if !strings.HasSuffix(info.Name(), ".tmpl") && !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}
		body, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		found += strings.Count(string(body), "respond.Invalid(c, err)")
		return nil
	})
	if found < 40 {
		t.Errorf("found only %d respond.Invalid call sites, so either the "+
			"conversion was reverted or this walk reads nothing", found)
	}
}

// The helper itself has to fill details, and answer 400 rather than 422 for a
// body that never parsed.
func TestInvalidFillsTheDetailsMap(t *testing.T) {
	body := stripComments(apiRespondGo())

	if !regexp.MustCompile(`func Invalid\(c \*gin\.Context, err error\)`).MatchString(body) {
		t.Fatal("respond has no Invalid helper")
	}
	if !strings.Contains(body, "ValidationFields(err)") {
		t.Error("Invalid does not use ValidationFields, so details stays empty " +
			"and the client has nothing to put against its inputs")
	}
	if !strings.Contains(body, "CodeInvalidBody") {
		t.Error("Invalid answers 422 for a body that never parsed; a shape whose " +
			"fields could not be checked is a 400")
	}
}

// And no template hands an internal error to the client.
//
// INTERNAL_ERROR's own catalogue entry says the message is deliberately vague
// and the detail is in the server log. respond.Fail(c, CodeInternalError,
// err.Error()) sends the Go error instead, which for a write is a GORM error
// carrying schema names and sometimes SQL. respond.Internal logs it and
// answers opaquely, which is what the entry describes.
//
// One site did this, in the generated tree handler's rebuild endpoint. The
// guard is here for the next one.
func TestNoTemplateHandsAnInternalErrorToTheClient(t *testing.T) {
	const bad = "respond.Fail(c, respond.CodeInternalError, err.Error())"

	var offenders []string
	err := filepath.Walk("..", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		name := info.Name()
		if !strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, ".tmpl") {
			return nil
		}
		if strings.Contains(name, "repair") || strings.HasSuffix(name, "_test.go") {
			return nil
		}
		body, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if strings.Contains(string(body), bad) {
			offenders = append(offenders, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking the templates: %v", err)
	}
	for _, path := range offenders {
		t.Errorf("%s answers INTERNAL_ERROR with the Go error's own text, which "+
			"can carry schema names and SQL. Use respond.Internal(c, err)", path)
	}
}
