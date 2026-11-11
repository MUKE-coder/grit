package scaffold

import (
	"fmt"
	"sort"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/errorcodes"
)

// The error catalogue, rendered into the two languages that have to agree about
// it.
//
// Both files are generated from internal/errorcodes rather than written, because
// the failure being fixed is exactly that kind of drift: the API returned a code,
// the frontend switched on a string somebody remembered, and nothing anywhere
// listed the set. A generated pair cannot disagree, and a test in the CLI refuses
// any code the templates emit that the catalogue does not carry.

// goConstName turns VALIDATION_ERROR into CodeValidationError.
//
// The rule itself lives in the catalogue package, because the test that checks
// every code a template emits is catalogued has to read the same constants this
// writes, and two copies of a naming rule is one copy too many.
func goConstName(code string) string {
	return errorcodes.GoConstName(code)
}

// apiRespondCodesGo returns internal/respond/codes.go: the catalogue as Go.
// projectCodes is .grit/error-codes.yaml, already validated, or nil.
//
// Merged into the same sorted list rather than emitted as a second map: a
// handler asks StatusOf one question and must get one answer, and two maps is
// how a code acquires two statuses depending on which is consulted. The
// comment on each entry says which list it came from, so it is still obvious
// what is yours.
func apiRespondCodesGo(projectCodes []errorcodes.Entry) string {
	entries := mergeCodes(errorcodes.All(), projectCodes)

	var b strings.Builder
	b.WriteString(`// Code generated from Grit's error catalogue. DO NOT EDIT.
//
// Every error this API returns has a code, and a code always arrives with the
// same status. That second half is the part a frontend needs and the part that
// used to be untrue: VALIDATION_ERROR came back as 422 in some handlers and 400
// in others, so a client had to handle whichever pairing it had happened to see.
//
// The same catalogue generates packages/shared/types/errors.ts, so a TypeScript
// switch over these codes can be exhaustive, and the table at /docs/backend/errors.
//
// Your own codes go in .grit/error-codes.yaml, NOT here: this file is
// regenerated and an entry added to it is lost on the next upgrade. Codes
// declared there are folded into the catalogue below and into the TypeScript,
// so the API and the frontend stay one list.
//
// Regenerate with: grit sync (or grit upgrade, which also does it)
package respond

import (
	"log"
	"net/http"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
)

// Code is one of the documented error codes. Use the constants below rather than
// a string literal: a typo in a literal is a code no client has ever heard of.
type Code string

// Category is what kind of problem a code reports, which is what decides how a
// caller should treat it.
type Category string

const (
	CategoryRequest    Category = "request"
	CategoryAuth       Category = "auth"
	CategoryPermission Category = "permission"
	CategoryNotFound   Category = "notfound"
	CategoryConflict   Category = "conflict"
	CategoryState      Category = "state"
	CategoryLimit      Category = "limit"
	CategoryServer     Category = "server"
	CategoryUpstream   Category = "upstream"
	CategoryDisabled   Category = "disabled"
)

const (
`)

	// Constants, with the meaning as the doc comment so it shows on hover.
	width := 0
	for _, entry := range entries {
		if n := len(goConstName(entry.Code)); n > width {
			width = n
		}
	}
	for _, entry := range entries {
		b.WriteString(fmt.Sprintf("\t// %s — %s\n", entry.Code, entry.Meaning))
		b.WriteString(fmt.Sprintf("\t%-*s Code = %q\n", width, goConstName(entry.Code), entry.Code))
	}

	b.WriteString(`)

// Meaning is everything the catalogue knows about one code.
type Meaning struct {
	Status   int
	Category Category
	// Area is the part of the API that raises it.
	Area string
	// What happened, from the caller's side.
	Meaning string
	// What the caller should do about it.
	Client string
}

// catalogue is the whole set. Generated, so it matches the documentation exactly.
var catalogue = map[Code]Meaning{
`)
	for _, entry := range entries {
		b.WriteString(fmt.Sprintf("\t%s: {Status: %s, Category: %s, Area: %q,\n\t\tMeaning: %q,\n\t\tClient:  %q},\n",
			goConstName(entry.Code), goStatusName(entry.Status), goCategoryName(string(entry.Category)),
			entry.Area, entry.Meaning, entry.Client))
	}

	b.WriteString(`}

// StatusOf is the status a code is always returned with, or 0 if the code is not
// one of the documented ones.
func StatusOf(code Code) int {
	return catalogue[code].Status
}

// Lookup returns what the catalogue says about a code.
func Lookup(code Code) (Meaning, bool) {
	meaning, ok := catalogue[code]
	return meaning, ok
}

// Codes lists every documented code, sorted, which is useful in a test that
// asserts your own handlers only return codes a client has been told about.
func Codes() []Code {
	out := make([]Code, 0, len(catalogue))
	for code := range catalogue {
		out = append(out, code)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// Fail writes a documented error, at the status the catalogue gives it.
//
// Prefer this to writing c.JSON with a status and a code side by side: that is
// how one code ends up with two statuses.
//
//	respond.Fail(c, respond.CodeForbidden, "Only an owner can archive an invoice")
//
// Pass per-field messages for a validation failure:
//
//	respond.Fail(c, respond.CodeValidationError, "Check the highlighted fields",
//	    map[string]string{"email": "That address is already in use"})
//
// # An unknown code
//
// Code is a string type, so respond.Code("WHATEVER") compiles. It is still a
// bug: nothing documents the code, the frontend cannot branch on it, and the
// status has to be guessed. So it is answered 500 and reported once, by name,
// with what to do about it.
//
// The 500 is deliberate and the report is the new part. Before it, a product
// whose handler refused a request with an undeclared code answered every one
// of those callers a 500 and the on-call was paged for what was a client
// error (grit#208). Nothing in the response or the log said which code, and
// it reads as correct in review.
func Fail(c *gin.Context, code Code, message string, details ...map[string]string) {
	status := StatusOf(code)
	if status == 0 {
		reportUndeclaredCode(code)
		status = http.StatusInternalServerError
	}
	body := Error{Code: string(code), Message: message}
	if len(details) > 0 && len(details[0]) > 0 {
		body.Details = details[0]
	}
	c.AbortWithStatusJSON(status, gin.H{"error": body})
}

// reported is the codes already complained about. Once each: a handler on a
// hot path would otherwise write the same line on every request, which buries
// it in the noise it is competing with.
var reported sync.Map

func reportUndeclaredCode(code Code) {
	if _, already := reported.LoadOrStore(code, true); already {
		return
	}
	log.Printf("respond: %q is not a declared error code, so this answered 500. "+
		"Declare it in .grit/error-codes.yaml and run grit sync, or use one of the "+
		"codes in this file. respond.Code(%q) compiles and is still a bug.", code, code)
}
`)
	return b.String()
}

// goStatusName turns 422 into http.StatusUnprocessableEntity, so the generated
// file reads like something a person wrote.
func goStatusName(status int) string {
	// Every 4xx and 5xx constant net/http declares, read out of net/http's own
	// source rather than typed: 418 is StatusTeapot and not StatusImATeapot,
	// and 414 is StatusRequestURITooLong, so a name derived from
	// http.StatusText is wrong for both and the generated file would not
	// compile.
	//
	// A hand-kept list of fifteen until v3.505.4, and 402 was not among them,
	// so PLAN_LIMIT_REACHED and any project code with an unusual status came
	// out as a bare number: valid Go, and the odd one out in a file whose
	// stated point is to read like something a person wrote.
	names := map[int]string{
		400: "http.StatusBadRequest",
		401: "http.StatusUnauthorized",
		402: "http.StatusPaymentRequired",
		403: "http.StatusForbidden",
		404: "http.StatusNotFound",
		405: "http.StatusMethodNotAllowed",
		406: "http.StatusNotAcceptable",
		407: "http.StatusProxyAuthRequired",
		408: "http.StatusRequestTimeout",
		409: "http.StatusConflict",
		410: "http.StatusGone",
		411: "http.StatusLengthRequired",
		412: "http.StatusPreconditionFailed",
		413: "http.StatusRequestEntityTooLarge",
		414: "http.StatusRequestURITooLong",
		415: "http.StatusUnsupportedMediaType",
		416: "http.StatusRequestedRangeNotSatisfiable",
		417: "http.StatusExpectationFailed",
		418: "http.StatusTeapot",
		421: "http.StatusMisdirectedRequest",
		422: "http.StatusUnprocessableEntity",
		423: "http.StatusLocked",
		424: "http.StatusFailedDependency",
		425: "http.StatusTooEarly",
		426: "http.StatusUpgradeRequired",
		428: "http.StatusPreconditionRequired",
		429: "http.StatusTooManyRequests",
		431: "http.StatusRequestHeaderFieldsTooLarge",
		451: "http.StatusUnavailableForLegalReasons",
		500: "http.StatusInternalServerError",
		501: "http.StatusNotImplemented",
		502: "http.StatusBadGateway",
		503: "http.StatusServiceUnavailable",
		504: "http.StatusGatewayTimeout",
		505: "http.StatusHTTPVersionNotSupported",
		506: "http.StatusVariantAlsoNegotiates",
		507: "http.StatusInsufficientStorage",
		508: "http.StatusLoopDetected",
		510: "http.StatusNotExtended",
		511: "http.StatusNetworkAuthenticationRequired",
	}
	if name, ok := names[status]; ok {
		return name
	}
	// A status outside the standard set. Numeric, which compiles.
	return fmt.Sprintf("%d", status)
}

// goCategoryName is an explicit map, not a capitalisation: "notfound" becomes
// CategoryNotFound, and a rule that only upper-cases the first letter produced
// CategoryNotfound, which the generated file then referred to and did not compile.
func goCategoryName(category string) string {
	names := map[string]string{
		"request":    "CategoryRequest",
		"auth":       "CategoryAuth",
		"permission": "CategoryPermission",
		"notfound":   "CategoryNotFound",
		"conflict":   "CategoryConflict",
		"state":      "CategoryState",
		"limit":      "CategoryLimit",
		"server":     "CategoryServer",
		"upstream":   "CategoryUpstream",
		"disabled":   "CategoryDisabled",
	}
	if name, ok := names[category]; ok {
		return name
	}
	// A category with no constant would be a compile error in the generated file,
	// so say so here instead, where a test will see it.
	return "Category_UNKNOWN_" + category
}

// sharedErrorsTS returns packages/shared/types/errors.ts.
//
// The point of generating this is the exhaustive switch: with a union type, a
// frontend that handles five codes and forgets the sixth fails to compile rather
// than falling through to "Something went wrong".
func sharedErrorsTS(projectCodes []errorcodes.Entry) string {
	entries := mergeCodes(errorcodes.All(), projectCodes)

	var b strings.Builder
	b.WriteString(`// Code generated from Grit's error catalogue. DO NOT EDIT.
//
// Every error the API returns has a code, and each code always arrives with the
// same HTTP status. Switch on the code, not on the status: the status groups
// errors, the code says which one it is.
//
//   import { API_ERRORS, type ApiErrorCode } from '@/shared/types/errors'
//
//   function explain(code: ApiErrorCode) {
//     return API_ERRORS[code].client   // what the person should do about it
//   }
//
// This project's own codes are declared in .grit/error-codes.yaml and appear
// here with the rest, which is what keeps a switch over ApiErrorCode
// exhaustive for the errors the product actually raises.
//
// Regenerate with: grit sync (or grit upgrade, which also does it)

export type ApiErrorCategory =
  | 'request'
  | 'auth'
  | 'permission'
  | 'notfound'
  | 'conflict'
  | 'state'
  | 'limit'
  | 'server'
  | 'upstream'
  | 'disabled'

export type ApiErrorCode =
`)
	for _, entry := range entries {
		b.WriteString(fmt.Sprintf("  | '%s'\n", entry.Code))
	}

	b.WriteString(`
export interface ApiErrorInfo {
  /** The status this code always arrives with. */
  status: number
  category: ApiErrorCategory
  /** Which part of the API raises it. */
  area: string
  /** What happened, from the caller's side. */
  meaning: string
  /** What to do about it. */
  client: string
}

export const API_ERRORS: Record<ApiErrorCode, ApiErrorInfo> = {
`)
	for _, entry := range entries {
		b.WriteString(fmt.Sprintf("  %s: {\n", entry.Code))
		b.WriteString(fmt.Sprintf("    status: %d,\n", entry.Status))
		b.WriteString(fmt.Sprintf("    category: '%s',\n", entry.Category))
		b.WriteString(fmt.Sprintf("    area: '%s',\n", entry.Area))
		b.WriteString(fmt.Sprintf("    meaning: %s,\n", tsString(entry.Meaning)))
		b.WriteString(fmt.Sprintf("    client: %s,\n", tsString(entry.Client)))
		b.WriteString("  },\n")
	}

	b.WriteString(`}

/** The shape of the error envelope every endpoint uses. */
export interface ApiErrorEnvelope {
  error: {
    code: string
    message: string
    /** Per-field messages, on a validation failure. */
    details?: Record<string, string>
  }
}

/** Narrows an unknown string to a documented code. */
export function isApiErrorCode(value: unknown): value is ApiErrorCode {
  return typeof value === 'string' && value in API_ERRORS
}

/**
 * The documented code behind a failure, narrowed to the union.
 *
 * Takes either the response body or the error a client library threw, because
 * both shapes turn up in practice. Returns null when the code is not one of the
 * documented ones, which covers two cases worth keeping apart from a known code:
 * a proxy or gateway answered instead of the API, or one of your own handlers
 * returned a code of its own. For the raw string, including your own codes, use
 * apiErrorCode from ./api.
 */
export function documentedErrorCode(value: unknown): ApiErrorCode | null {
  const asError = value as { response?: { data?: ApiErrorEnvelope } } | undefined
  const code = asError?.response?.data?.error?.code ?? (value as ApiErrorEnvelope | undefined)?.error?.code
  return isApiErrorCode(code) ? code : null
}

/** What the person in front of the screen should do about it. */
export function apiErrorAdvice(code: ApiErrorCode): string {
  return API_ERRORS[code].client
}

/** Codes that are worth retrying unchanged. Everything else needs a decision. */
export function isRetryable(code: ApiErrorCode): boolean {
  const category = API_ERRORS[code].category
  return category === 'server' || category === 'upstream' || category === 'limit' || category === 'disabled'
}

/** Every documented code, sorted, for a table or a test. */
export const API_ERROR_CODES = Object.keys(API_ERRORS).sort() as ApiErrorCode[]
`)
	return b.String()
}

// tsString quotes a string for TypeScript, preferring single quotes the way the
// rest of the shared package is written.
func tsString(value string) string {
	escaped := strings.ReplaceAll(value, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, "'", `\'`)
	return "'" + escaped + "'"
}

// mergeCodes is Grit's catalogue plus the project's own, sorted by code.
//
// One list, because StatusOf asks one question. A second map consulted after
// the first is how a code acquires two statuses depending on which is read,
// which is the inconsistency the catalogue was built to end.
//
// A project code that shadows a built-in never reaches here: LoadProject
// refuses the file and says which code and what Grit answers it with. If one
// somehow did, the built-in wins, so the framework's own handlers cannot have
// their status changed under them.
func mergeCodes(built, project []errorcodes.Entry) []errorcodes.Entry {
	if len(project) == 0 {
		return built
	}
	known := make(map[string]bool, len(built))
	for _, e := range built {
		known[e.Code] = true
	}
	out := make([]errorcodes.Entry, 0, len(built)+len(project))
	out = append(out, built...)
	for _, e := range project {
		if known[e.Code] {
			continue
		}
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out
}

// projectErrorCodes is .grit/error-codes.yaml for the project at root.
//
// A malformed file is reported and skipped rather than aborting. Aborting
// would block an upgrade over a typo in a file the upgrade does not own, and
// skipping is not silent in practice: the constant a handler uses is then
// absent from the generated file, so the build fails right after this warning
// says why. `grit sync` refuses outright, which is the command whose whole job
// this is.
func projectErrorCodes(root string) []errorcodes.Entry {
	codes, err := errorcodes.LoadProject(root)
	if err != nil {
		fmt.Printf("  ⚠ %v\n", err)
		fmt.Printf("    The codes in that file are not in this build. Fix it and run grit sync.\n")
		return nil
	}
	return codes
}
