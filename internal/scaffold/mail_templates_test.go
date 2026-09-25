package scaffold

import (
	"strings"
	"testing"
)

// Email is not a web page, and the rules are not stylistic.
//
// A <style> block is the one part of an HTML email that is not reliably
// delivered: Outlook on Windows renders through Word and drops most of it, and
// Gmail strips it when a message is clipped or forwarded. The six templates
// that shipped were each a full document with its own <style> block, so for a
// large share of recipients they arrived as unstyled serif text.
func TestEmailsAreStyledWhereAClientWillSeeIt(t *testing.T) {
	// Comments in the emitter mention <style>; the emitted templates must not
	// contain one, so the check is on the rendered markup only.
	src := stripGoComments(mailTemplatesGo())

	if strings.Contains(src, "<style>") {
		t.Error("a template still carries a <style> block, which Outlook and a clipped Gmail drop")
	}
	if strings.Contains(src, "class=") {
		t.Error("a template still uses a class, which has nothing to match without a <style> block")
	}
	if !strings.Contains(src, `role="presentation"`) {
		t.Error("the layout is not table-based, so Outlook will not lay it out")
	}
	// A button is a table cell with a background, not a padded <a>: Outlook
	// ignores padding on an inline element and renders a bare link.
	if !strings.Contains(src, `style="{{.BtnCell}}"`) || !strings.Contains(src, `style="{{.BtnLink}}"`) {
		t.Error("the buttons do not use the two-part cell and link styles")
	}
	// One shell, six fragments. Six documents is how they drifted last time.
	if strings.Count(src, "<!DOCTYPE") != 1 {
		t.Errorf("%d documents in the mail package; there should be one layout", strings.Count(src, "<!DOCTYPE"))
	}
	if !strings.Contains(src, "{{.Preheader}}") {
		t.Error("no preheader, so the inbox preview repeats the heading")
	}
}

// The colour comes from the theme, so a project built with --theme emerald
// does not send purple email.
func TestEmailsUseTheProjectsTheme(t *testing.T) {
	// The palettes in theme.go are hex by necessity: an email is rendered on
	// somebody else's machine and there are no CSS variables there. What must
	// not carry a colour is a template.
	for _, hex := range []string{"#6c5ce7", "#0a0a0f", "#111118", "#e8e8f0"} {
		if strings.Contains(mailTemplatesGo(), hex) {
			t.Errorf("a template hardcodes %s instead of taking the theme's colour", hex)
		}
	}

	src := mailThemeGo()
	if !strings.Contains(src, "func Theme() string {") {
		t.Fatal("nothing reads THEME, so email cannot follow the two frontends")
	}
	// Every theme the CLI offers has a palette here, or a project built with it
	// sends mail in somebody else's colours.
	for _, name := range ValidThemes {
		if !strings.Contains(src, `"`+name+`":`) {
			t.Errorf("no email palette for theme %q", name)
		}
	}

	// html/template filters anything interpolated into a style attribute down
	// to ZgotmplZ unless it is marked as CSS the package wrote itself. Without
	// this every style attribute in every email is the literal text ZgotmplZ,
	// which is how the first version of this shipped in testing.
	if !strings.Contains(src, "template.CSS(heading)") {
		t.Error("the style strings are not template.CSS, so html/template will blank them")
	}
}

// All six are in Mail Preview. Two were not, which is how the sign-in code
// email kept a letter-spacing nobody had looked at.
func TestEveryEmailIsInThePreview(t *testing.T) {
	preview := mailPreviewGo()
	for _, name := range []string{
		"welcome", "password-reset", "email-verification",
		"notification", "two-factor-code", "magic-link",
	} {
		if !strings.Contains(preview, `builtIn("`+name+`"`) {
			t.Errorf("%s is sent by the app and missing from Mail Preview", name)
		}
	}
	// And the preview renders through the same layout the mailer uses, rather
	// than a copy that can drift.
	if !strings.Contains(preview, `template.New("layout").Parse(Layout)`) {
		t.Error("the preview parses its own layout")
	}
	if !strings.Contains(preview, "data := Style(Theme())") {
		t.Error("the preview does not give the layout a palette, so it renders with empty colours")
	}
}

// stripGoComments drops // lines, so a comment explaining a rule does not read
// as a breach of it.
func stripGoComments(src string) string {
	var kept []string
	for _, line := range strings.Split(src, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}
