package errorcodes

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The templates are scanned for three shapes, which are the three ways a
// generated project writes an error today:
//
//	c.JSON(http.StatusX, gin.H{"error": gin.H{"code": "CODE", ...}})
//	fail(c, http.StatusX, "CODE", ...)
//	code = "CODE"   (inside a switch that sets the status separately)
var (
	codeInEnvelope = regexp.MustCompile(`"code":\s*"([A-Z][A-Z0-9_]+)"`)
	codeAssigned   = regexp.MustCompile(`\bcode\s*=\s*"([A-Z][A-Z0-9_]+)"`)
	codeInFail     = regexp.MustCompile(`\bfail\(c,\s*[^,]+,\s*"([A-Z][A-Z0-9_]+)"`)
	statusNamed    = regexp.MustCompile(`http\.Status([A-Za-z]+)`)
)

// statusValues are the names gin handlers use, mapped to their numbers. Only the
// ones the templates actually use: an unrecognised name is reported rather than
// guessed, because guessing is how a status mismatch would slip through.
var statusValues = map[string]int{
	"OK": 200, "Created": 201, "NoContent": 204,
	"BadRequest": 400, "Unauthorized": 401, "Forbidden": 403, "NotFound": 404,
	"MethodNotAllowed": 405, "Conflict": 409, "Gone": 410,
	"RequestEntityTooLarge": 413, "UnsupportedMediaType": 415,
	"UnprocessableEntity": 422, "TooManyRequests": 429,
	"InternalServerError": 500, "NotImplemented": 501, "BadGateway": 502,
	"ServiceUnavailable": 503,
}

type emission struct {
	code   string
	status int // 0 when the status is set elsewhere, which is the switch shape
	file   string
	line   int
}

// scanTemplates walks the packages that hold the Go templates and returns every
// error code they emit.
func scanTemplates(t *testing.T) []emission {
	t.Helper()
	roots := []string{
		filepath.Join("..", "scaffold"),
		filepath.Join("..", "generate"),
	}

	var out []emission
	for _, root := range roots {
		entries, err := os.ReadDir(root)
		if err != nil {
			t.Fatalf("reading %s: %v", root, err)
		}
		for _, file := range entries {
			name := file.Name()
			if file.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
				continue
			}
			path := filepath.Join(root, name)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("reading %s: %v", path, err)
			}
			lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
			for i, line := range lines {
				var found []string
				for _, pattern := range []*regexp.Regexp{codeInEnvelope, codeAssigned, codeInFail} {
					for _, match := range pattern.FindAllStringSubmatch(line, -1) {
						found = append(found, match[1])
					}
				}
				if len(found) == 0 {
					continue
				}
				// The status belongs to the statement, which can start a few lines
				// above the code: c.JSON opens, gin.H follows, the code is inside it.
				window := strings.Join(lines[max(0, i-5):i+1], "\n")
				status := 0
				if names := statusNamed.FindAllStringSubmatch(window, -1); len(names) > 0 {
					last := names[len(names)-1][1]
					value, known := statusValues[last]
					if !known {
						t.Errorf("%s:%d: http.Status%s is not in statusValues; add it", path, i+1, last)
					}
					status = value
				}
				for _, code := range found {
					out = append(out, emission{code: code, status: status, file: path, line: i + 1})
				}
			}
		}
	}
	if len(out) < 200 {
		t.Fatalf("only %d error emissions found in the templates; the scanner has stopped matching", len(out))
	}
	return out
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Every code a generated project can return has to be in the catalogue.
//
// Without this, the catalogue is a document that drifts. With it, inventing a code
// inline fails the build, and the only way to add one is to say what it means and
// what a caller should do about it.
func TestEveryEmittedCodeIsCatalogued(t *testing.T) {
	missing := map[string][]string{}
	for _, e := range scanTemplates(t) {
		if _, ok := Lookup(e.code); !ok {
			where := fmt.Sprintf("%s:%d", filepath.Base(e.file), e.line)
			missing[e.code] = append(missing[e.code], where)
		}
	}
	codes := make([]string, 0, len(missing))
	for code := range missing {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	for _, code := range codes {
		t.Errorf("%s is returned by %s but is not in the catalogue: add it to catalog.go with its status, meaning and what a client should do",
			code, strings.Join(missing[code], ", "))
	}
}

// And it always arrives with the same status.
//
// This is the check that matters to whoever is writing the frontend. One code
// with two statuses cannot be branched on: VALIDATION_ERROR was 422 in
// thirty-eight places and 400 in twenty-five, and INVALID_TOKEN was both 401 and
// 400, so a client had to handle every pairing it happened to see.
func TestEmittedStatusMatchesTheCatalogue(t *testing.T) {
	for _, e := range scanTemplates(t) {
		entry, ok := Lookup(e.code)
		if !ok || e.status == 0 {
			continue // reported by the test above, or set elsewhere in a switch
		}
		if e.status != entry.Status {
			t.Errorf("%s:%d: %s is returned with %d here and catalogued as %d; one code, one status",
				e.file, e.line, e.code, e.status, entry.Status)
		}
	}
}

// The catalogue itself has to be usable: no duplicates, and every entry says what
// happened and what to do about it.
func TestCatalogueIsWellFormed(t *testing.T) {
	seen := map[string]bool{}
	for _, entry := range catalog {
		if seen[entry.Code] {
			t.Errorf("%s appears twice in the catalogue", entry.Code)
		}
		seen[entry.Code] = true

		if entry.Status < 400 || entry.Status > 599 {
			t.Errorf("%s has status %d, which is not an error", entry.Code, entry.Status)
		}
		if entry.Category == "" || entry.Area == "" {
			t.Errorf("%s is missing its category or area", entry.Code)
		}
		for field, value := range map[string]string{"meaning": entry.Meaning, "client": entry.Client} {
			if len(value) < 20 {
				t.Errorf("%s has no useful %s: %q", entry.Code, field, value)
			}
			if !strings.HasSuffix(value, ".") {
				t.Errorf("%s's %s should read as a sentence: %q", entry.Code, field, value)
			}
		}
	}
}
