package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// The bridge reads the context key the middleware writes.
//
// Two spellings of one fact, and they cannot be made one. middleware imports
// internal/i18n for Default, so the bridge, which lives in package i18n,
// cannot import middleware back: it repeats the key as a literal instead.
//
// Nothing fails if they drift. The bridge reads a key nothing ever set, every
// lookup falls back, and every API error is English again with no error
// anywhere. Exactly what this release was fixing, reintroduced silently. So
// the comparison lives here, where both templates are in reach.
//
// This guard is on the Grit side and not in the generated project because a
// test inside package i18n cannot import middleware either: the cycle is real
// at test time too, which is how the first version of this test was found.
func TestTheBridgeReadsTheKeyTheMiddlewareWrites(t *testing.T) {
	bridgeKey := constValue(t, i18nRespondBridgeGo(), "localeKey")
	middlewareKey := constValue(t, i18nMiddlewareGo(), "LocaleKey")

	if bridgeKey == "" || middlewareKey == "" {
		t.Fatalf("could not read both keys: bridge %q, middleware %q",
			bridgeKey, middlewareKey)
	}
	if bridgeKey != middlewareKey {
		t.Errorf("the bridge reads %q and the middleware writes %q, so no "+
			"translation is ever found and every API error stays English",
			bridgeKey, middlewareKey)
	}
}

// constValue reads a string constant's value out of a template.
func constValue(t *testing.T, body, name string) string {
	t.Helper()
	re := regexp.MustCompile(`const ` + regexp.QuoteMeta(name) + ` = "([^"]*)"`)
	if match := re.FindStringSubmatch(body); match != nil {
		return match[1]
	}
	return ""
}

// The bridge is written, and written into a package something imports.
//
// An init only runs in a package that is linked in. internal/response is the
// obvious home for this and the wrong one: nothing imports it, so its init
// would never run and the registration would silently not happen, in a way no
// test of the catalogues could see. middleware imports i18n, routes import
// middleware, main imports routes.
func TestTheBridgeIsWrittenIntoAPackageThatIsLinked(t *testing.T) {
	bridge := i18nRespondBridgeGo()

	if !strings.Contains(bridge, "package i18n") {
		t.Error("the bridge is not in package i18n, so its init may never run")
	}
	if !strings.Contains(bridge, "respond.Translator = func") {
		t.Error("the bridge does not register a translator, so respond has none")
	}
	// It must not translate the default locale: a handler's own message often
	// says more than the code's general one, and overriding English would
	// discard the better sentence at every call site at once.
	if !strings.Contains(bridge, "if locale == Default {") {
		t.Error("the bridge translates the default locale too, which replaces " +
			"every handler's own English message with the code's general one")
	}
	// Lookup and not T: T answers with the key when nothing has it, which
	// respond would then show the reader in place of the English message.
	if !strings.Contains(bridge, "Lookup(locale, key)") {
		t.Error("the bridge uses something other than Lookup, so an untranslated " +
			"code reaches the reader as a catalogue key")
	}
	// Comments stripped first. The bridge's own doc comment says the words
	// "middleware.LocaleOf" while explaining why it must not call it, so a
	// Contains check over the whole file matches the prose that explains it
	// and fails on a correct bridge.
	if strings.Contains(stripComments(bridge), "middleware.") {
		t.Error("the bridge imports middleware, which imports i18n: a cycle")
	}
}

// respond keeps working with no translator at all.
//
// respond is in every project and the catalogues are in none by default, so
// the hook being nil is the normal case and has to be the cheap one.
func TestRespondNeedsNoTranslator(t *testing.T) {
	codes := apiRespondCodesGo(nil)

	if !strings.Contains(codes, "var Translator func(c *gin.Context, key string) string") {
		t.Error("respond has no translator hook, so nothing can translate its errors")
	}
	if !strings.Contains(codes, "if Translator == nil || c == nil {") {
		t.Error("respond does not guard a nil translator, so a project without " +
			"i18n panics on its first error")
	}
	// Fail and FailT share one writer, or the two disagree about the status and
	// the envelope, and FailT's specific message is overwritten by Fail's
	// lookup of the code's general one.
	if !strings.Contains(codes, "func write(c *gin.Context, code Code, message string") {
		t.Error("Fail and FailT do not share a writer")
	}
	if strings.Contains(codes, "func FailT") && !strings.Contains(codes, "write(c, code, message, details...)") {
		t.Error("FailT writes through Fail, which translates again and discards " +
			"the message FailT exists to deliver")
	}
}
