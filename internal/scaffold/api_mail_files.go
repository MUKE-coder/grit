package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeMailFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	// The package is framework-owned: mailer.go, the drivers, FromConfig, their
	// tests and the mailtest fake. templates.go is written with them, and so is
	// the admin's Mail Preview handler, which reads the package's registry.
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "mail", "templates.go"):        mailTemplatesGo(),
		filepath.Join(apiRoot, "internal", "mail", "theme.go"):            mailThemeGo(),
		filepath.Join(apiRoot, "internal", "handlers", "mail_preview.go"): mailPreviewHandlerGo(),
	}
	for rel, content := range mailPackageFiles() {
		files[filepath.Join(apiRoot, "internal", "mail", filepath.FromSlash(rel))] = content
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// mailThemeGo emits internal/mail/theme.go.
func mailThemeGo() string {
	return `package mail

import (
	"html/template"
	"os"
	"strings"
)

// Palette is the handful of colours an email needs.
//
// Deliberately not the full design system: an email has a canvas, a card, two
// weights of text, a rule and a brand colour, and anything more is a colour a
// mail client will get wrong somewhere.
type Palette struct {
	Canvas   string // behind the card
	Card     string // the card itself
	Border   string
	Ink      string // headings
	Muted    string // body copy
	Accent   string // the brand, and the button
	AccentFG string // the label on that button
}

// palettes mirrors the light values in the admin's globals.css. Kept as literal
// hex rather than read from anywhere, because an email is rendered on somebody
// else's machine and there are no CSS variables there.
var palettes = map[string]Palette{
	"atlas":    {Canvas: "#f6f7f9", Card: "#ffffff", Border: "#e2e8f0", Ink: "#0f172a", Muted: "#475569", Accent: "#2563eb", AccentFG: "#ffffff"},
	"aurora":   {Canvas: "#f5f5f7", Card: "#ffffff", Border: "#d2d2d7", Ink: "#1d1d1f", Muted: "#424245", Accent: "#1d1d1f", AccentFG: "#ffffff"},
	"pulse":    {Canvas: "#f6f7f9", Card: "#ffffff", Border: "#e0e4e9", Ink: "#1d1f26", Muted: "#4b5563", Accent: "#0051c3", AccentFG: "#ffffff"},
	"coral":    {Canvas: "#fdf7f5", Card: "#ffffff", Border: "#efe0da", Ink: "#26140f", Muted: "#6b4a40", Accent: "#e2503f", AccentFG: "#ffffff"},
	"amber":    {Canvas: "#fdfaf3", Card: "#ffffff", Border: "#ece0c9", Ink: "#241c0c", Muted: "#6b5a38", Accent: "#b45309", AccentFG: "#ffffff"},
	"sky":      {Canvas: "#f4f9fd", Card: "#ffffff", Border: "#d9e6f2", Ink: "#0c1f2e", Muted: "#3f5c72", Accent: "#0284c7", AccentFG: "#ffffff"},
	"mono":     {Canvas: "#fafafa", Card: "#ffffff", Border: "#e5e5e5", Ink: "#171717", Muted: "#525252", Accent: "#171717", AccentFG: "#ffffff"},
	"emerald":  {Canvas: "#f9fafb", Card: "#ffffff", Border: "#e5e7eb", Ink: "#111827", Muted: "#4b5563", Accent: "#059669", AccentFG: "#ffffff"},
	"midnight": {Canvas: "#f6f7f9", Card: "#ffffff", Border: "#e2e8f0", Ink: "#0f172a", Muted: "#475569", Accent: "#6c5ce7", AccentFG: "#ffffff"},
}

// Theme is the theme name this deployment is using.
//
// The same THEME the two frontends read, so changing it in .env repaints the
// email with them rather than leaving one surface behind.
func Theme() string {
	name := strings.ToLower(strings.TrimSpace(os.Getenv("THEME")))
	if _, ok := palettes[name]; ok {
		return name
	}
	return "atlas"
}

// PaletteFor returns a theme's colours, falling back to the default rather than
// to zero values: an email with an empty colour in a style attribute is worse
// than an email in the wrong brand.
func PaletteFor(name string) Palette {
	if p, ok := palettes[strings.ToLower(name)]; ok {
		return p
	}
	return palettes["atlas"]
}

// Style is everything a template interpolates: the palette plus the type.
//
// The style strings live here rather than in each template because there are
// six templates and one idea of what a heading looks like. A template that
// wants something different writes it inline, which is still one place.
//
// They are template.CSS, not string. html/template treats the inside of a
// style attribute as a CSS context and filters anything interpolated into it
// down to ZgotmplZ, which is exactly what it should do for a value that came
// from a user and exactly wrong for one the package wrote itself. Marking them
// says these are ours.
func Style(theme string) map[string]interface{} {
	p := PaletteFor(theme)
	font := "-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif"

	heading := "font-family:" + font + ";font-size:20px;line-height:1.35;font-weight:700;letter-spacing:-0.01em;color:" + p.Ink + ";margin:0 0 12px;"
	body := "font-family:" + font + ";font-size:15px;line-height:1.65;color:" + p.Muted + ";margin:0 0 16px;"

	return map[string]interface{}{
		"Canvas":   p.Canvas,
		"Card":     p.Card,
		"Border":   p.Border,
		"Ink":      p.Ink,
		"Muted":    p.Muted,
		"Accent":   p.Accent,
		"AccentFG": p.AccentFG,
		"Font":     template.CSS(font),
		"H1":       template.CSS(heading),
		"P":        template.CSS(body),
		// The last paragraph in a card, with its bottom margin taken off so the
		// card's own padding is the only space below it.
		"PLast": template.CSS(body + "margin-bottom:0;"),
		// A button is a table cell with a background and a link inside it.
		// Padding on an inline <a> is ignored by Outlook, which turns the
		// button into a bare underlined link.
		"BtnCell": template.CSS("background-color:" + p.Accent + ";border-radius:10px;"),
		"BtnLink": template.CSS("display:inline-block;padding:13px 26px;font-family:" + font +
			";font-size:15px;font-weight:600;line-height:1;color:" + p.AccentFG + ";text-decoration:none;"),
	}
}
`
}

func mailTemplatesGo() string {
	return `package mail

// The six emails, as fragments. Layout is the shell they are rendered into.
//
// Fragments rather than whole documents, because six copies of one shell is how
// they drifted apart last time: three different footers and two different greys
// between them.
var EmailTemplates = map[string]string{
	"welcome":            welcomeTemplate,
	"password-reset":     passwordResetTemplate,
	"email-verification": emailVerificationTemplate,
	"notification":       notificationTemplate,
	"two-factor-code":    twoFactorCodeTemplate,
	"magic-link":         magicLinkTemplate,
}

// Layout wraps every fragment.
//
// Tables and style attributes, which is not a stylistic choice. Outlook on
// Windows renders mail through Word, which ignores most of a <style> block and
// all of flexbox; Gmail drops <style> when it clips a long message or when
// somebody forwards it. Anything that has to survive goes on the element.
//
// The style strings the fragments use ({{.H1}}, {{.P}} and the rest) are
// supplied by the renderer, so there is one definition of what a heading looks
// like rather than one per template.
const Layout = ` + "`" + `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<meta name="x-apple-disable-message-reformatting" />
<meta name="color-scheme" content="light" />
<meta name="supported-color-schemes" content="light" />
<title>{{.Subject}}</title>
</head>
<body style="margin:0;padding:0;background-color:{{.Canvas}};">
<div style="display:none;max-height:0;overflow:hidden;mso-hide:all;font-size:1px;line-height:1px;color:{{.Canvas}};">{{.Preheader}}&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;</div>
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color:{{.Canvas}};">
<tr>
<td align="center" style="padding:32px 12px;">
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="max-width:560px;width:100%;">
<tr>
<td align="left" style="padding:0 4px 18px;font-family:{{.Font}};font-size:17px;font-weight:700;letter-spacing:-0.01em;color:{{.Accent}};">{{.AppName}}</td>
</tr>
<tr>
<td style="background-color:{{.Card}};border:1px solid {{.Border}};border-radius:14px;padding:32px;">
{{.Content}}
</td>
</tr>
<tr>
<td align="left" style="padding:18px 4px 0;font-family:{{.Font}};font-size:12px;line-height:1.6;color:{{.Muted}};">
&copy; {{.Year}} {{.AppName}}. Sent to you because this address was used on {{.AppName}}.
</td>
</tr>
</table>
</td>
</tr>
</table>
</body>
</html>` + "`" + `

// Fragments: the inside of the card, and nothing else.
//
// Divs with explicit margins rather than h1/p, because Outlook adds margins of
// its own to block elements that no stylesheet can take back.

const welcomeTemplate = ` + "`" + `<div style="{{.H1}}">Welcome to {{.AppName}}</div>
<div style="{{.P}}">Hi {{.Name}}, your account is ready. There is nothing else to set up.</div>
{{if .ActionURL}}<table role="presentation" cellpadding="0" cellspacing="0" border="0" style="margin:0 0 20px;"><tr><td align="center" style="{{.BtnCell}}"><a href="{{.ActionURL}}" style="{{.BtnLink}}">{{.ActionText}}</a></td></tr></table>{{end}}
<div style="{{.PLast}}">If you did not create this account you can ignore this message, and nothing further will happen.</div>` + "`" + `

const passwordResetTemplate = ` + "`" + `<div style="{{.H1}}">Reset your password</div>
<div style="{{.P}}">Somebody asked to reset the password for this address. Use the button below to choose a new one.</div>
<table role="presentation" cellpadding="0" cellspacing="0" border="0" style="margin:0 0 20px;"><tr><td align="center" style="{{.BtnCell}}"><a href="{{.ResetURL}}" style="{{.BtnLink}}">Choose a new password</a></td></tr></table>
<div style="{{.PLast}}">The link works once and expires in an hour. If this was not you, nothing has changed and you can ignore this message.</div>` + "`" + `

const emailVerificationTemplate = ` + "`" + `<div style="{{.H1}}">{{.Title}}</div>
<div style="{{.P}}">{{.Message}}</div>
<table role="presentation" cellpadding="0" cellspacing="0" border="0" style="margin:0 0 20px;"><tr><td align="center" style="{{.BtnCell}}"><a href="{{.VerifyURL}}" style="{{.BtnLink}}">Confirm this address</a></td></tr></table>
<div style="{{.PLast}}">If the button does not work, paste this into your browser:<br /><span style="word-break:break-all;color:{{.Accent}};">{{.VerifyURL}}</span></div>` + "`" + `

const notificationTemplate = ` + "`" + `<div style="{{.H1}}">{{.Title}}</div>
<div style="{{.P}}">{{.Message}}</div>
{{if .ActionURL}}<table role="presentation" cellpadding="0" cellspacing="0" border="0" style="margin:0;"><tr><td align="center" style="{{.BtnCell}}"><a href="{{.ActionURL}}" style="{{.BtnLink}}">{{.ActionText}}</a></td></tr></table>{{end}}` + "`" + `

const magicLinkTemplate = ` + "`" + `<div style="{{.H1}}">Your sign-in link</div>
<div style="{{.P}}">Use the button below to sign in. There is no password to type.</div>
<table role="presentation" cellpadding="0" cellspacing="0" border="0" style="margin:0 0 20px;"><tr><td align="center" style="{{.BtnCell}}"><a href="{{.MagicURL}}" style="{{.BtnLink}}">Sign in</a></td></tr></table>
<div style="{{.PLast}}">The link lasts {{.Minutes}} minutes and works once. If you did not ask for it, somebody knows your email address, which is not a secret: they cannot sign in without this mailbox.</div>` + "`" + `

const twoFactorCodeTemplate = ` + "`" + `<div style="{{.H1}}">Your sign-in code</div>
<div style="{{.P}}">Type this into the page you already have open.</div>
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="margin:0 0 18px;">
<tr><td align="center" style="background-color:{{.Canvas}};border:1px solid {{.Border}};border-radius:12px;padding:20px 12px;font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;font-size:32px;font-weight:700;letter-spacing:10px;color:{{.Ink}};">{{.Code}}</td></tr>
</table>
<div style="{{.PLast}}">It expires in {{.Minutes}} minutes and works once. If you were not signing in, somebody has your password: change it. This code alone lets nobody in.</div>` + "`" + `
`
}
