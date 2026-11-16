package scaffold

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// Every catalogue must carry the same keys as the English one, and the same
// markers inside them.
//
// A key present in en.json and absent everywhere else does not fail a build: t()
// falls back to the key's English text, so a French page reports one English
// label in the middle of a translated table and nobody notices until somebody
// reads it. That is how fr and sw shipped without table.emptyAll and without
// five nav entries for months.
//
// A missing marker is worse than a missing key, because it reaches the user as a
// defect rather than as English. The API translates through fmt.Sprintf, so a
// Spanish message that lost a %s renders "%!s(MISSING)", and next-intl renders a
// web message that lost a {count} with the wrong number or none at all. Keys do
// not catch either: the key is there, the sentence is broken.
//
// The API has had the key check since it was written, but only inside the
// generated project (i18nTestGo emits internal/i18n/i18n_test.go), so the
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

	keys := map[string]bool{}
	for key := range catalogueMessages(t, code, body) {
		keys[key] = true
	}
	return keys
}

// catalogueMessages flattens a catalogue into dotted path and translated text.
func catalogueMessages(t *testing.T, code, body string) map[string]string {
	t.Helper()

	var doc map[string]any
	if err := json.Unmarshal([]byte(body), &doc); err != nil {
		t.Fatalf("%s catalogue is not valid JSON: %v\n%s", code, err, body)
	}

	messages := map[string]string{}
	var walk func(prefix string, node map[string]any)
	walk = func(prefix string, node map[string]any) {
		for key, value := range node {
			path := prefix + key
			if child, ok := value.(map[string]any); ok {
				walk(path+".", child)
				continue
			}
			text, ok := value.(string)
			if !ok {
				t.Fatalf("%s catalogue: %s is %T, not a string", code, path, value)
			}
			messages[path] = text
		}
	}
	walk("", doc)
	return messages
}

// TestAPICataloguesKeepTheirFormatVerbs holds every locale to the printf verbs
// its English message has. T() renders through fmt.Sprintf, so a translation
// that drops a verb shows the reader "%!s(MISSING)" in the middle of a sentence.
func TestAPICataloguesKeepTheirFormatVerbs(t *testing.T) {
	catalogues := map[string]func() string{
		"en": i18nLocaleEN,
		"es": i18nLocaleES,
		"fr": i18nLocaleFR,
		"sw": i18nLocaleSW,
	}

	want := catalogueMessages(t, "en", i18nLocaleEN())
	for _, code := range []string{"es", "fr", "sw"} {
		got := catalogueMessages(t, code, catalogues[code]())
		for _, key := range sortedMessageKeys(want) {
			w := formatVerbs(want[key])
			g := formatVerbs(got[key])
			if strings.Join(w, " ") != strings.Join(g, " ") {
				t.Errorf("internal/i18n/locales/%s.json: %q has %v, en has %v", code, key, g, w)
			}
		}
	}
}

// TestWebCataloguesKeepTheirPlaceholders is the same check for next-intl, whose
// placeholders are braces rather than printf verbs. A message that lost one
// renders with a blank where a number belongs.
func TestWebCataloguesKeepTheirPlaceholders(t *testing.T) {
	catalogues := map[string]func() string{
		"en": i18nMessagesEN,
		"es": i18nMessagesES,
		"fr": i18nMessagesFR,
		"sw": i18nMessagesSW,
	}

	want := catalogueMessages(t, "en", i18nMessagesEN())
	for _, code := range []string{"es", "fr", "sw"} {
		got := catalogueMessages(t, code, catalogues[code]())
		for _, key := range sortedMessageKeys(want) {
			w := placeholders(want[key])
			g := placeholders(got[key])
			if strings.Join(w, " ") != strings.Join(g, " ") {
				t.Errorf("messages/%s.json: %q has %v, en has %v", code, key, g, w)
			}
		}
	}
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

func sortedMessageKeys(messages map[string]string) []string {
	keys := make([]string, 0, len(messages))
	for key := range messages {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// formatVerb matches one printf verb. A literal percent is skipped by
// formatVerbs, because "%%" consumes no argument.
var formatVerb = regexp.MustCompile(`%[-+# 0]*[0-9]*(\.[0-9]+)?[A-Za-z%]`)

// formatVerbs lists the printf verbs in a message, sorted, so two translations
// of one key compare equal however they order the sentence.
func formatVerbs(message string) []string {
	verbs := []string{}
	for _, verb := range formatVerb.FindAllString(message, -1) {
		if verb != "%%" {
			verbs = append(verbs, verb)
		}
	}
	sort.Strings(verbs)
	return verbs
}

// messagePlaceholder matches a next-intl placeholder: {count}, {total}, {name}.
var messagePlaceholder = regexp.MustCompile(`\{([A-Za-z0-9_]+)\}`)

func placeholders(message string) []string {
	names := []string{}
	for _, match := range messagePlaceholder.FindAllStringSubmatch(message, -1) {
		names = append(names, match[1])
	}
	sort.Strings(names)
	return names
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

func TestCatalogueMarkersAreRead(t *testing.T) {
	if got := formatVerbs("at least %d of %s"); strings.Join(got, " ") != "%d %s" {
		t.Errorf("formatVerbs read %v, want [%%d %%s]", got)
	}
	if got := formatVerbs("100%% done"); len(got) != 0 {
		t.Errorf("formatVerbs counted a literal percent as a verb: %v", got)
	}
	if got := placeholders("{count} of {total}"); strings.Join(got, " ") != "count total" {
		t.Errorf("placeholders read %v, want [count total]", got)
	}

	// Both extractors have to find something in the real catalogues, or the two
	// tests above pass by comparing empty lists.
	api := catalogueMessages(t, "en", i18nLocaleEN())
	if got := formatVerbs(api["validation.min_length"]); len(got) != 2 {
		t.Errorf("the English API catalogue no longer exercises two verbs: %v", got)
	}
	web := catalogueMessages(t, "en", i18nMessagesEN())
	if got := placeholders(web["table.showing"]); len(got) != 3 {
		t.Errorf("the English web catalogue no longer exercises three placeholders: %v", got)
	}
}
