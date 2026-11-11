package scaffold

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/errorcodes"
)

// grit#208: a project code has to reach BOTH generated files, because the
// property that makes the catalogue worth having is that the Go constant and
// the TypeScript union are one list. A Go-only extension would give the API
// its code and leave the frontend unable to type a case for the one error the
// product most needs to handle well.

func planLimitEntry() errorcodes.Entry {
	return errorcodes.Entry{
		Code:     "THESIS_CHECKS_EXHAUSTED",
		Status:   402,
		Category: errorcodes.CategoryLimit,
		Area:     "billing",
		Meaning:  "The account has used all three of its similarity checks.",
		Client:   "Show the upgrade prompt. Retrying will not help.",
	}
}

func TestAProjectCodeReachesTheGoFile(t *testing.T) {
	src := apiRespondCodesGo([]errorcodes.Entry{planLimitEntry()})

	// The constant, which is what a handler writes.
	if !strings.Contains(src, "CodeThesisChecksExhausted") {
		t.Error("no typed constant for the project's code, so a handler still has to write " +
			"respond.Code(\"...\") and get a 500")
	}
	// And its status, which is the half that was missing: StatusOf not knowing
	// the code is what answered 500.
	if !strings.Contains(src, "http.StatusPaymentRequired") {
		t.Error("the project's code is in the catalogue without its status")
	}
	if !strings.Contains(src, "THESIS_CHECKS_EXHAUSTED") {
		t.Error("the wire code is missing")
	}
}

func TestAProjectCodeReachesTheTypeScript(t *testing.T) {
	src := sharedErrorsTS([]errorcodes.Entry{planLimitEntry()})

	if !strings.Contains(src, "'THESIS_CHECKS_EXHAUSTED'") {
		t.Error("the project's code is not in the ApiErrorCode union, so an exhaustive switch " +
			"cannot have a case for it: the frontend half of grit#208")
	}
	if !strings.Contains(src, "402") {
		t.Error("the code is in the union without its status")
	}
}

// Grit's own codes are still all there. An extension point that drops the
// built-ins would be a worse bug than the one it fixes.
func TestProjectCodesDoNotDisplaceTheCatalogue(t *testing.T) {
	with := apiRespondCodesGo([]errorcodes.Entry{planLimitEntry()})
	for _, entry := range errorcodes.All() {
		if !strings.Contains(with, `"`+entry.Code+`"`) {
			t.Fatalf("%s is missing once a project declares a code of its own", entry.Code)
		}
	}
}

// A project cannot change what one of Grit's codes answers with. LoadProject
// refuses the file, so this never arrives, and if it did the built-in wins:
// the framework's own handlers must not have their status moved under them.
func TestAProjectCannotChangeABuiltInStatus(t *testing.T) {
	shadow := errorcodes.Entry{
		Code:     "VALIDATION_ERROR",
		Status:   418,
		Category: errorcodes.CategoryRequest,
		Area:     "project",
		Meaning:  "m",
		Client:   "c",
	}
	src := apiRespondCodesGo([]errorcodes.Entry{shadow})
	if strings.Contains(src, "StatusTeapot") {
		t.Error("a project entry overrode a built-in code's status")
	}
	// Exactly one entry for it, not two: two is how a code acquires two
	// statuses depending on which is read.
	if n := strings.Count(src, `CodeValidationError: {`); n != 1 {
		t.Errorf("VALIDATION_ERROR appears %d times in the catalogue, want 1", n)
	}
}

// No project codes is the common case and must change nothing at all, because
// every existing project is in it.
func TestNoProjectCodesGeneratesExactlyWhatItDidBefore(t *testing.T) {
	if apiRespondCodesGo(nil) != apiRespondCodesGo([]errorcodes.Entry{}) {
		t.Error("nil and empty produce different files")
	}
	if sharedErrorsTS(nil) != sharedErrorsTS([]errorcodes.Entry{}) {
		t.Error("nil and empty produce different TypeScript")
	}
}

// The generated file says where to put a code of your own. The reporter's
// first problem was not that it could not be done but that codes.go said DO
// NOT EDIT and nothing said what to do instead.
func TestTheGeneratedFilesSayWhereProjectCodesGo(t *testing.T) {
	for name, src := range map[string]string{
		"codes.go":  apiRespondCodesGo(nil),
		"errors.ts": sharedErrorsTS(nil),
	} {
		if !strings.Contains(src, errorcodes.ProjectFile) {
			t.Errorf("%s does not name %s, so a reader told DO NOT EDIT has nowhere to go",
				name, errorcodes.ProjectFile)
		}
	}
}

// And Fail reports an unknown code instead of answering 500 in silence.
func TestFailReportsAnUndeclaredCode(t *testing.T) {
	src := apiRespondCodesGo(nil)
	if !strings.Contains(src, "reportUndeclaredCode(code)") {
		t.Error("an undeclared code is still answered 500 with nothing logged, so the on-call " +
			"gets a page and no way to tell which code did it (grit#208)")
	}
	if !strings.Contains(src, errorcodes.ProjectFile) {
		t.Error("the report does not say how to declare the code")
	}
	// Once per code: a handler on a hot path would otherwise bury the line in
	// its own repetitions.
	if !strings.Contains(src, "reported.LoadOrStore") {
		t.Error("the report is not deduplicated")
	}
}

// Every status in Grit's own catalogue renders as a named constant, not a
// number. The map was hand-kept and 402 was missing, so PLAN_LIMIT_REACHED
// came out as a bare 402 in a file whose stated point is to read like
// something a person wrote.
func TestEveryCatalogueStatusHasANamedConstant(t *testing.T) {
	for _, entry := range errorcodes.All() {
		name := goStatusName(entry.Status)
		if !strings.HasPrefix(name, "http.Status") {
			t.Errorf("%s is %d, which renders as %q rather than an http.Status constant",
				entry.Code, entry.Status, name)
		}
	}
}

// And every status a project may legally declare renders as valid Go, named or
// numeric. A name that does not exist in net/http would stop the project
// compiling, which is why the map is read out of net/http rather than guessed.
func TestEveryDeclarableStatusRendersAsValidGo(t *testing.T) {
	for status := 400; status <= 599; status++ {
		if http.StatusText(status) == "" {
			continue // not declarable: LoadProject refuses it
		}
		name := goStatusName(status)
		switch {
		case strings.HasPrefix(name, "http.Status"):
			// Named. The map came from net/http's source, so it exists.
		case name == strconv.Itoa(status):
			// Numeric, which compiles.
		default:
			t.Errorf("status %d renders as %q, which is neither", status, name)
		}
	}
	// The two Go names nobody would derive correctly from http.StatusText.
	if got := goStatusName(418); got != "http.StatusTeapot" {
		t.Errorf("418 renders as %q; net/http calls it StatusTeapot", got)
	}
	if got := goStatusName(414); got != "http.StatusRequestURITooLong" {
		t.Errorf("414 renders as %q; net/http calls it StatusRequestURITooLong", got)
	}
}
