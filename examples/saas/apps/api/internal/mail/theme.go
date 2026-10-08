package mail

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

	// template.CSS on every rule below, because html/template will not put an
	// unmarked string into a style attribute: it writes ZgotmplZ instead, and
	// the mail arrives unstyled with no error anywhere.
	//
	// Safe because none of it comes from outside. Every value is either a
	// constant in this file or a field of the palette, which is chosen by name
	// from a fixed set in this package. Nothing a user or a request can reach
	// is interpolated here, so there is no input to escape.
	//nolint:gosec // G203: our own stylesheet, no external input
	//#nosec G203
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
