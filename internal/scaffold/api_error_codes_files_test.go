package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/errorcodes"
)

// The Go and TypeScript halves are generated from one catalogue so they cannot
// disagree. These tests hold the part that generation alone does not guarantee:
// that both files carry every code, that the Go one compiles, and that they reach
// an existing project.
func TestGeneratedCodesCarryTheWholeCatalogue(t *testing.T) {
	goSrc := apiRespondCodesGo()
	tsSrc := sharedErrorsTS()
	durableMustFormat(t, "internal/respond/codes.go", goSrc)

	for _, entry := range errorcodes.All() {
		// The constants are aligned, so the name and the value are checked apart:
		// asserting one padded line is asserting gofmt's column widths.
		if !strings.Contains(goSrc, goConstName(entry.Code)) ||
			!strings.Contains(goSrc, fmt.Sprintf("Code = %q", entry.Code)) {
			t.Errorf("codes.go has no constant for %s", entry.Code)
		}
		if !strings.Contains(tsSrc, fmt.Sprintf("  | '%s'\n", entry.Code)) {
			t.Errorf("errors.ts does not have %s in the union", entry.Code)
		}
		if !strings.Contains(tsSrc, fmt.Sprintf("  %s: {\n", entry.Code)) {
			t.Errorf("errors.ts has no API_ERRORS entry for %s", entry.Code)
		}
	}

	// The catalogue's own count, so a generator that silently emits half the list
	// is caught rather than passing on the codes it happens to include. Counted
	// from the code union only: the category union above it is also a list of
	// quoted members.
	at := strings.Index(tsSrc, "export type ApiErrorCode =")
	if at < 0 {
		t.Fatal("errors.ts has no ApiErrorCode union")
	}
	union := tsSrc[at : strings.Index(tsSrc[at:], "export interface ApiErrorInfo")+at]
	if got, want := strings.Count(union, "  | '"), len(errorcodes.All()); got != want {
		t.Errorf("errors.ts declares %d codes in the union, the catalogue has %d", got, want)
	}
}

// A handler gets the status from the catalogue rather than pairing one with a code
// by hand, which is how the same code ended up with two statuses.
func TestGeneratedCodesTakeTheStatusFromTheCatalogue(t *testing.T) {
	goSrc := apiRespondCodesGo()
	durableMustContain(t, "internal/respond/codes.go", goSrc,
		"func Fail(c *gin.Context, code Code, message string, details ...map[string]string)",
		"status := StatusOf(code)",
		// An undocumented code is a bug in the handler, and a 500 says so rather
		// than inventing a status for it.
		"status = http.StatusInternalServerError",
		"CodeValidationError: {Status: http.StatusUnprocessableEntity",
		"CodeNotFound: {Status: http.StatusNotFound",
		"CodeVersionConflict: {Status: http.StatusConflict",
	)

	ts := sharedErrorsTS()
	durableMustContain(t, "packages/shared/types/errors.ts", ts,
		"export type ApiErrorCode =",
		"export const API_ERRORS: Record<ApiErrorCode, ApiErrorInfo> = {",
		"export function isApiErrorCode(",
		"export function documentedErrorCode(",
		"export function isRetryable(",
		"  VALIDATION_ERROR: {\n    status: 422,",
		"  NOT_FOUND: {\n    status: 404,",
	)

	// Every identifier the generated file refers to has to be declared in it.
	// gofmt parses, it does not resolve names, so a file that was valid Go and
	// referred to CategoryNotfound (for the notfound category) passed this test
	// and failed to compile in a real project.
	for _, reference := range regexp.MustCompile(`Category[A-Za-z_]+`).FindAllString(goSrc, -1) {
		if reference == "Category" {
			continue
		}
		// The declarations are aligned, so the gap before Category is whitespace of
		// whatever width gofmt chose.
		declared := regexp.MustCompile(regexp.QuoteMeta(reference) + `\s+Category = `)
		if !declared.MatchString(goSrc) {
			t.Errorf("codes.go refers to %s, which it does not declare", reference)
		}
	}

	// The TS must not carry Go escaping: a catalogue sentence with an apostrophe
	// in it would otherwise end the string early and break every build.
	if strings.Contains(ts, `\"`) {
		t.Error("errors.ts has Go-style escaped quotes in it")
	}
	for _, line := range strings.Split(ts, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "meaning: '") && !strings.HasPrefix(trimmed, "client: '") {
			continue
		}
		if !strings.HasSuffix(trimmed, "',") {
			t.Errorf("a catalogue sentence is not a closed TS string: %s", trimmed)
		}
	}
}

// Both files have to arrive in an existing project: generated handlers call
// respond.Fail, so a project with the older respond package would not compile,
// and a frontend switching on a stale union is the drift this is meant to end.
func TestCatalogueFilesTravelOnUpgrade(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "app", Architecture: ArchTriple}

	if err := writeRespondFiles(root, opts); err != nil {
		t.Fatal(err)
	}
	if err := writeSharedFiles(root, opts); err != nil {
		t.Fatal(err)
	}
	// The scaffold path too, and separately: the first version of this test called
	// only writeRespondFiles, which is the upgrade path, and passed while a
	// freshly scaffolded project had no codes.go and did not compile.
	newRoot := t.TempDir()
	if err := writeAPIFiles(newRoot, opts); err != nil {
		t.Fatal(err)
	}
	fresh := filepath.Join(opts.APIRoot(newRoot), filepath.FromSlash("internal/respond/codes.go"))
	if _, err := os.Stat(fresh); err != nil {
		t.Errorf("a new project gets no internal/respond/codes.go, so generated handlers calling respond.Fail will not compile: %v", err)
	}

	codes := filepath.Join(opts.APIRoot(root), filepath.FromSlash("internal/respond/codes.go"))
	data, err := os.ReadFile(codes)
	if err != nil {
		t.Fatalf("internal/respond/codes.go was not written: %v", err)
	}
	if strings.Contains(string(data), "{{MODULE}}") {
		t.Error("codes.go still has an unreplaced module placeholder")
	}

	errorsTS := filepath.Join(root, filepath.FromSlash("packages/shared/types/errors.ts"))
	if _, err := os.Stat(errorsTS); err != nil {
		t.Errorf("packages/shared/types/errors.ts was not written: %v", err)
	}

	// And the barrel exports it, or nothing in the frontend can import it by the
	// name the docs give.
	index, err := os.ReadFile(filepath.Join(root, filepath.FromSlash("packages/shared/types/index.ts")))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"API_ERRORS", "type ApiErrorCode", `from "./errors"`} {
		if !strings.Contains(string(index), want) {
			t.Errorf("types/index.ts does not export %s", want)
		}
	}
}

// Constant names are part of the API of the generated package, so the mapping
// from code to name is pinned: initialisms stay whole, which is the part a
// naive converter gets wrong.
func TestGoConstNames(t *testing.T) {
	for code, want := range map[string]string{
		"VALIDATION_ERROR":    "CodeValidationError",
		"NOT_FOUND":           "CodeNotFound",
		"AI_RATE_LIMITED":     "CodeAIRateLimited",
		"INVALID_CSV":         "CodeInvalidCSV",
		"SMS_FAILED":          "CodeSMSFailed",
		"CSRF_INVALID":        "CodeCSRFInvalid",
		"TOTP_ERROR":          "CodeTOTPError",
		"DB_ERROR":            "CodeDBError",
		"PDF_ERROR":           "CodePDFError",
		"API_KEY_REQUIRED":    "CodeAPIKeyRequired",
		"INVALID_TOTP_CODE":   "CodeInvalidTOTPCode",
		"PASSKEY_REJECTED":    "CodePasskeyRejected",
		"VERSION_CONFLICT":    "CodeVersionConflict",
		"STORAGE_UNAVAILABLE": "CodeStorageUnavailable",
	} {
		if got := goConstName(code); got != want {
			t.Errorf("goConstName(%q) = %q, want %q", code, got, want)
		}
	}
}
