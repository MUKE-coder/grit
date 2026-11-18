package errorcodes

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// Every code in the translation table is a code that exists.
//
// The table is a side map rather than a seventh field on 124 positional
// literals, which is the readable choice and the one that lets a name rot: a
// code renamed in the catalogue leaves a key here pointing at nothing, and
// nothing about that fails to compile. The message would simply never be
// translated again, and the reader would see English forever.
func TestEveryTranslatedCodeExists(t *testing.T) {
	known := map[string]bool{}
	for _, entry := range All() {
		known[entry.Code] = true
	}

	for _, code := range Translated() {
		if !known[code] {
			t.Errorf("messageKeys has %q, which is not a code in the catalogue: "+
				"renamed, or a typo. Either way its message is never translated", code)
		}
	}
}

// And every key it points at is one the catalogues actually carry.
//
// The other direction of the same rot: a key renamed in en.json leaves this
// table pointing at a key no catalogue has, Lookup misses, and respond quietly
// keeps the English message. Nothing breaks, and the translation silently
// stops happening, which is the failure this whole change exists to end.
func TestEveryTranslatedKeyIsInTheCatalogues(t *testing.T) {
	english := readCatalogue(t, "en")
	if len(english) < 30 {
		t.Fatalf("read %d keys from the English catalogue, which is too few to "+
			"be the real one", len(english))
	}

	for _, code := range Translated() {
		key := MessageKeyOf(code)
		if !english[key] {
			t.Errorf("%s maps to %q, which the English catalogue does not have, "+
				"so it is never translated in any locale", code, key)
		}
	}
}

// The codes that carry the traffic must stay translated.
//
// Coverage here is a number somebody chose, so it needs a floor that fails
// when it drops rather than a comment saying what it was. The busiest codes
// are counted from the templates themselves: if a handler starts using a code
// heavily and nobody translates it, this says so.
//
// The threshold is deliberately low. A code used five times in the framework
// is one a user meets, and the point of keying on the code was that covering
// it is cheap.
func TestTheCodesHandlersActuallyUseAreTranslated(t *testing.T) {
	const minUses = 5

	uses := countCodeUses(t)
	translated := map[string]bool{}
	for _, code := range Translated() {
		translated[code] = true
	}

	// Codes used often and deliberately left in English, each with the reason.
	// A machine reads these, so a translation is work nobody would ever see.
	noHumanReader := map[string]string{
		"SENTINEL_OFF":         "the WAF's own admin endpoints, operator only",
		"ENDPOINT_NOT_ALLOWED": "answered to an API key, not a browser",
		"INVALID_API_KEY":      "answered to an API key, not a browser",
		"INVALID_SIGNATURE":    "answered to a webhook sender",
		"ORIGIN_NOT_ALLOWED":   "answered to a browser that was not asked to make the call",
		"UNKNOWN_SETTING":      "a developer calling the settings API wrongly",
	}

	var missing []string
	for code, n := range uses {
		if n < minUses || translated[code] {
			continue
		}
		if _, deliberate := noHumanReader[code]; deliberate {
			continue
		}
		missing = append(missing, code)
	}
	sort.Strings(missing)

	for _, code := range missing {
		t.Errorf("%s is used %d times and has no translated message. Add it to "+
			"messageKeys, or to noHumanReader here with the reason nobody reads it",
			code, uses[code])
	}

	// Guard the guard: a counter that found nothing would make the loop above
	// pass with any table at all.
	if uses["VALIDATION_ERROR"] < 50 {
		t.Errorf("counted %d uses of VALIDATION_ERROR, so the template walk is "+
			"not reading the handlers", uses["VALIDATION_ERROR"])
	}
	for _, code := range []string{"SENTINEL_OFF", "ENDPOINT_NOT_ALLOWED",
		"INVALID_API_KEY", "INVALID_SIGNATURE", "ORIGIN_NOT_ALLOWED",
		"UNKNOWN_SETTING"} {
		if uses[code] < minUses {
			continue
		}
		if _, listed := noHumanReader[code]; !listed {
			t.Errorf("%s is excused and is not in noHumanReader", code)
		}
	}
}

// countCodeUses is how many times each code is named in the templates that
// become a project, keyed by the code and not the Go constant.
func countCodeUses(t *testing.T) map[string]int {
	t.Helper()

	byConst := map[string]string{}
	for _, entry := range All() {
		byConst[GoConstName(entry.Code)] = entry.Code
	}

	uses := map[string]int{}
	ref := regexp.MustCompile(`respond\.(Code[A-Za-z0-9]+)`)
	root := filepath.Join("..", "..", "internal")

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		name := info.Name()
		if strings.HasSuffix(name, "_test.go") {
			return nil
		}
		if !strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, ".tmpl") {
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, match := range ref.FindAllStringSubmatch(string(body), -1) {
			if code, ok := byConst[match[1]]; ok {
				uses[code]++
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking the templates: %v", err)
	}
	return uses
}

// readCatalogue pulls the keys out of a locale's template, which lives in the
// scaffold package. Read as text rather than imported, because importing
// scaffold from here would be a cycle: scaffold already imports errorcodes.
func readCatalogue(t *testing.T, locale string) map[string]bool {
	t.Helper()

	path := filepath.Join("..", "scaffold", "api_i18n_locales.go")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the catalogues: %v", err)
	}

	// The locale's own function, up to the next one, so a key in es.json is
	// not read as one in en.json.
	marker := "func i18nLocale" + strings.ToUpper(locale) + "() string {"
	start := strings.Index(string(body), marker)
	if start == -1 {
		t.Fatalf("no catalogue function for %q", locale)
	}
	rest := string(body)[start+len(marker):]
	if end := strings.Index(rest, "\nfunc "); end != -1 {
		rest = rest[:end]
	}

	keys := map[string]bool{}
	for _, match := range regexp.MustCompile(`"([a-z_]+(?:\.[a-z_]+)+)":`).FindAllStringSubmatch(rest, -1) {
		keys[match[1]] = true
	}
	return keys
}
