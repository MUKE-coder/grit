package scaffold

import (
	"encoding/json"
	"regexp"
	"sort"
	"testing"
)

// Every catalogue must carry the same keys as the English one.
//
// A key present in en.json and absent everywhere else does not fail a build: t()
// falls back to the key's English text, so a French page reports one English
// label in the middle of a translated table and nobody notices until somebody
// reads it. That is how fr and sw shipped without table.emptyAll and without
// five nav entries for months.
//
// The API has had this check since it was written, but only inside the generated
// project (i18nLocaleES's template emits internal/i18n/i18n_test.go), so the
// templates themselves were never compared to each other. This is that
// comparison, and it runs against Grit's own build.
func TestWebCataloguesShareTheSameKeys(t *testing.T) {
	catalogues := map[string]func() string{
		"en": i18nMessagesEN,
		"es": i18nMessagesES,
		"fr": i18nMessagesFR,
		"sw": i18nMessagesSW,
	}

	want := catalogueKeys(t, "en", i18nMessagesEN())
	for _, code := range []string{"es", "fr", "sw"} {
		got := catalogueKeys(t, code, catalogues[code]())
		for _, key := range missingKeys(want, got) {
			t.Errorf("messages/%s.json is missing %q (it will show English there)", code, key)
		}
		for _, key := range missingKeys(got, want) {
			t.Errorf("messages/%s.json has %q, which en.json does not", code, key)
		}
	}
}

func TestAPICataloguesShareTheSameKeys(t *testing.T) {
	catalogues := map[string]func() string{
		"en": i18nLocaleEN,
		"es": i18nLocaleES,
		"fr": i18nLocaleFR,
		"sw": i18nLocaleSW,
	}

	want := catalogueKeys(t, "en", i18nLocaleEN())
	for _, code := range []string{"es", "fr", "sw"} {
		got := catalogueKeys(t, code, catalogues[code]())
		for _, key := range missingKeys(want, got) {
			t.Errorf("internal/i18n/locales/%s.json is missing %q", code, key)
		}
		for _, key := range missingKeys(got, want) {
			t.Errorf("internal/i18n/locales/%s.json has %q, which en.json does not", code, key)
		}
	}
}

// The switcher offers what lib/locale.ts lists, not what the app ships. A
// catalogue with no LOCALES entry is a language nobody can choose; a LOCALES
// entry with no catalogue is a switcher option that renders its raw keys.
func TestLocaleSwitcherOffersEveryCatalogue(t *testing.T) {
	lib := i18nLocaleLibTS()
	codes := localeLibCodes(t, lib)

	for _, code := range []string{"en", "es", "fr", "sw"} {
		if !codes[code] {
			t.Errorf("lib/locale.ts does not offer %q, so its catalogue is unreachable", code)
		}
	}
	for code := range codes {
		switch code {
		case "en", "es", "fr", "sw":
		default:
			t.Errorf("lib/locale.ts offers %q, which has no catalogue", code)
		}
	}
}

// catalogueKeys flattens a catalogue into dotted paths, so the nested Next.js
// catalogues and the API's flat dotted keys compare the same way. A nested
// object is a section, never a translated value.
func catalogueKeys(t *testing.T, code, body string) map[string]bool {
	t.Helper()

	var doc map[string]any
	if err := json.Unmarshal([]byte(body), &doc); err != nil {
		t.Fatalf("%s catalogue is not valid JSON: %v\n%s", code, err, body)
	}

	keys := map[string]bool{}
	var walk func(prefix string, node map[string]any)
	walk = func(prefix string, node map[string]any) {
		for key, value := range node {
			path := prefix + key
			if child, ok := value.(map[string]any); ok {
				walk(path+".", child)
				continue
			}
			keys[path] = true
		}
	}
	walk("", doc)
	return keys
}

var localeLibEntry = regexp.MustCompile(`code:\s*'([a-z]{2}(?:-[A-Za-z]{2,4})?)'`)

// localeLibCodes reads the codes out of LOCALES in the emitted lib/locale.ts.
func localeLibCodes(t *testing.T, lib string) map[string]bool {
	t.Helper()

	codes := map[string]bool{}
	for _, match := range localeLibEntry.FindAllStringSubmatch(lib, -1) {
		codes[match[1]] = true
	}
	if len(codes) == 0 {
		t.Fatal("LOCALES in lib/locale.ts has no readable entries; the regexp and the template disagree")
	}
	return codes
}

// missingKeys returns the keys in want that got lacks, sorted so a failing test
// reads the same on every run.
func missingKeys(want, got map[string]bool) []string {
	var missing []string
	for key := range want {
		if !got[key] {
			missing = append(missing, key)
		}
	}
	sort.Strings(missing)
	return missing
}

// Guards the guard: a parser that finds nothing would make every comparison
// above pass, which is the shape of a test that only pretends to run.
func TestCatalogueKeysReadsBothCatalogueShapes(t *testing.T) {
	if got := catalogueKeys(t, "en", i18nMessagesEN()); !got["table.emptyAll"] {
		t.Error("a nested key was not read")
	}
	if got := catalogueKeys(t, "en", i18nLocaleEN()); !got["errors.not_found"] {
		t.Error("a flat dotted key was not read")
	}
}
