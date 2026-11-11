package errorcodes

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Codes a project adds to the catalogue, from .grit/error-codes.yaml.
//
// Grit's catalogue cannot hold every code every product needs. A SaaS refusing
// a request because the account has used everything its plan allows is a
// domain fact, not a framework one, and the nearest built-in codes are all
// wrong in a way a client can feel: RATE_LIMITED says "retry shortly", which
// is false, and FORBIDDEN is indistinguishable from a permission failure, so a
// UI can only tell them apart by matching on the message string (grit#208).
//
// Before this there was nowhere to put one. internal/respond/codes.go is
// generated and says DO NOT EDIT, so an entry added there is lost on the next
// upgrade. The workaround that suggests itself compiles and fails at runtime:
//
//	respond.Fail(c, respond.Code("PLAN_LIMIT_REACHED"), "...")
//
// Code is a string type, so that is legal Go, and StatusOf does not know the
// code, so every caller who hit their plan limit was answered 500.
//
// # Why a file in the project and not a Go file in internal/respond
//
// Because both Go and TypeScript are generated from one list, and that is the
// property worth having: an exhaustive TypeScript switch over error codes is
// only exhaustive if the frontend knows about the project's codes too. A Go
// extension file would give the API its code and leave the frontend unable to
// type a case for the one error the product most needs to handle well.

// ProjectFile is where a project declares its own codes, under the project
// root. Inside .grit/, which a generated project commits: only backups, slots
// and the token cache are ignored there.
const ProjectFile = ".grit/error-codes.yaml"

type projectDocument struct {
	Codes []projectEntry `yaml:"codes"`
}

type projectEntry struct {
	Code     string `yaml:"code"`
	Status   int    `yaml:"status"`
	Category string `yaml:"category"`
	Area     string `yaml:"area"`
	Meaning  string `yaml:"meaning"`
	Client   string `yaml:"client"`
}

// codeName is SHOUTING_SNAKE_CASE, which is what every code in the catalogue
// is and what the generated Go constant name is derived from.
var codeName = regexp.MustCompile(`^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$`)

// LoadProject reads the project's own codes, sorted by code.
//
// A missing file is not an error: most projects have none. A malformed one is,
// and loudly, because the alternative is generating a catalogue that silently
// lacks the code a handler is about to use, which is the failure this exists to
// remove.
func LoadProject(root string) ([]Entry, error) {
	path := filepath.Join(root, filepath.FromSlash(ProjectFile))
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading %s: %w", ProjectFile, err)
	}

	var doc projectDocument
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("%s is not valid YAML: %w", ProjectFile, err)
	}

	seen := map[string]bool{}
	out := make([]Entry, 0, len(doc.Codes))
	for i, e := range doc.Codes {
		entry, err := e.validate(i, seen)
		if err != nil {
			return nil, err
		}
		seen[entry.Code] = true
		out = append(out, entry)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out, nil
}

// validate turns one declaration into an Entry, or says what is wrong with it
// in terms of the file the person edited.
//
// Every rule here is a failure that would otherwise be silent or confusing
// later, which is the bar for refusing somebody's file.
func (e projectEntry) validate(index int, seen map[string]bool) (Entry, error) {
	where := fmt.Sprintf("%s: codes[%d]", ProjectFile, index)
	if e.Code != "" {
		where = fmt.Sprintf("%s: %s", ProjectFile, e.Code)
	}

	switch {
	case e.Code == "":
		return Entry{}, fmt.Errorf("%s has no code", where)
	case !codeName.MatchString(e.Code):
		// The generated Go constant name is derived from this, so a lower-case
		// or hyphenated code produces a constant nobody would guess, or one
		// that does not compile.
		return Entry{}, fmt.Errorf("%s is not a code name: use SHOUTING_SNAKE_CASE, as every "+
			"code in the catalogue does", where)
	case seen[e.Code]:
		return Entry{}, fmt.Errorf("%s is declared twice", where)
	}

	// Shadowing a built-in would give one code two statuses depending on which
	// entry won, which is the exact inconsistency the catalogue exists to end.
	if _, built := Lookup(e.Code); built {
		return Entry{}, fmt.Errorf("%s is already one of Grit's codes, at %d. Pick another name: "+
			"two entries for one code means its status depends on which of them wins",
			where, StatusOf(e.Code))
	}

	if e.Status < 400 || e.Status > 599 || http.StatusText(e.Status) == "" {
		return Entry{}, fmt.Errorf("%s has status %d, which is not an HTTP error status. An error "+
			"code is answered with a 4xx or a 5xx", where, e.Status)
	}
	category := Category(strings.ToLower(strings.TrimSpace(e.Category)))
	if !validCategory(category) {
		return Entry{}, fmt.Errorf("%s has category %q. The categories are: %s",
			where, e.Category, strings.Join(categoryNames(), ", "))
	}

	// Both of these are shown to people: Meaning in the error table on
	// /docs/backend/errors and as the TypeScript doc comment, Client in both.
	// An entry without them is a code a frontend developer finds and cannot
	// act on, which is the state the catalogue was built to replace.
	if strings.TrimSpace(e.Meaning) == "" {
		return Entry{}, fmt.Errorf("%s has no meaning: one sentence on what happened, from the "+
			"caller's side", where)
	}
	if strings.TrimSpace(e.Client) == "" {
		return Entry{}, fmt.Errorf("%s has no client: what the caller should do about it. This is "+
			"the field a frontend developer actually wants", where)
	}

	area := strings.ToLower(strings.TrimSpace(e.Area))
	if area == "" {
		area = "project"
	}

	return Entry{
		Code:     e.Code,
		Status:   e.Status,
		Category: category,
		Area:     area,
		Meaning:  strings.TrimSpace(e.Meaning),
		Client:   strings.TrimSpace(e.Client),
	}, nil
}

// Categories is every category a code may be in, in the order the docs list
// them.
func Categories() []Category {
	return []Category{
		CategoryRequest, CategoryAuth, CategoryPermission, CategoryNotFound,
		CategoryConflict, CategoryState, CategoryLimit, CategoryServer,
		CategoryUpstream, CategoryDisabled,
	}
}

func validCategory(c Category) bool {
	for _, known := range Categories() {
		if known == c {
			return true
		}
	}
	return false
}

func categoryNames() []string {
	out := make([]string, 0, len(Categories()))
	for _, c := range Categories() {
		out = append(out, string(c))
	}
	return out
}

// StarterProjectFile is the commented-out example a new project gets, so the
// extension point is found by reading the directory rather than the docs.
func StarterProjectFile() string {
	return `# Error codes this project adds to Grit's catalogue.
#
# Grit's own codes are in internal/respond/codes.go, which is generated and
# must not be edited: an entry added there is lost on the next grit upgrade.
# Add yours here instead. They are folded into the same two generated files, so
# the API and the frontend stay one list:
#
#   apps/api/internal/respond/codes.go     the typed constant and its status
#   packages/shared/types/errors.ts        so a TypeScript switch stays exhaustive
#
# Run ` + "`grit sync`" + ` after editing this file.
#
# Every field is required except area, which groups the code in the error table
# at /docs/backend/errors and defaults to "project".

codes: []

# codes:
#   - code: THESIS_ALREADY_SUBMITTED      # SHOUTING_SNAKE_CASE
#     status: 409                         # a 4xx or 5xx
#     category: state                     # request, auth, permission, notfound,
#                                         # conflict, state, limit, server,
#                                         # upstream, disabled
#     area: submissions
#     meaning: The thesis has already been submitted and cannot be changed.
#     client: Show the submitted version. Withdrawing it is a separate action.
`
}
