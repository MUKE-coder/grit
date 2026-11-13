package scaffold

import (
	"fmt"
	"strconv"
	"strings"
)

// themePaletteCSS is the nine themes, as CSS custom properties.
//
// Every surface a project ships reads the same variables -- the web app, the
// admin panel and the docs site -- so a single THEME=<name> in .env paints all
// of them. This was two identical copies, in the admin's stylesheet and the web
// app's; the docs site would have made a third, and a third copy is how a theme
// ends up applied to two frontends out of three.
//
// Values, not indirections: a surface maps these onto its own token names
// (--color-*, --color-fd-*) in its own stylesheet.
func themePaletteCSS() string {
	return `/* atlas (default) */
:root,
[data-theme="atlas"] {
  --bg-primary: #ffffff;
  --bg-secondary: #f8fafc;
  --bg-tertiary: #f1f5f9;
  --bg-elevated: #ffffff;
  --bg-hover: #f1f5f9;
  --border: #e2e8f0;
  --text-primary: #0f172a;
  --text-secondary: #475569;
  --text-muted: #64748b;
  --accent: #2563eb;
  --accent-hover: #1d4ed8;
  --accent-fg: #ffffff;
  --success: #10b981;
  --danger: #ef4444;
  --warning: #f59e0b;
  --info: #0ea5e9;
}

/* aurora — friendly, pastel, consumer SaaS */
/* aurora — Apple-inspired. Monochrome: near-black text and CTAs on white and
 * Apple's warm greys. Blue is reserved for links/info only, so the accent that
 * drives buttons stays black like iCloud's sign-in pill. */
[data-theme="aurora"] {
  --bg-primary: #fbfbfd;
  --bg-secondary: #ffffff;
  --bg-tertiary: #f5f5f7;
  --bg-elevated: #ffffff;
  --bg-hover: #f5f5f7;
  --border: #d2d2d7;
  --text-primary: #1d1d1f;
  --text-secondary: #424245;
  --text-muted: #86868b;
  --accent: #1d1d1f;
  --accent-hover: #000000;
  --accent-fg: #ffffff;
  --success: #10b981;
  --danger: #ef4444;
  --warning: #f59e0b;
  --info: #0071e3;
}

/* pulse — Cloudflare-inspired. Premium blue CTAs on a cool grey-blue canvas
 * with white elevated cards; Cloudflare orange is the single warm accent. */
[data-theme="pulse"] {
  --bg-primary: #f6f7f9;
  --bg-secondary: #ffffff;
  --bg-tertiary: #eef1f5;
  --bg-elevated: #ffffff;
  --bg-hover: #eef1f5;
  --border: #e0e4e9;
  --text-primary: #1d1f26;
  --text-secondary: #4b5563;
  --text-muted: #8a94a6;
  --accent: #0051c3;
  --accent-hover: #003d99;
  --accent-fg: #ffffff;
  --success: #16a34a;
  --danger: #dc2626;
  --warning: #f6821f;
  --info: #0051c3;
}

/* coral - rose, warm neutrals. The marketplace palette: the accent is a
 * statement colour, so the greys around it stay very plain. */
[data-theme="coral"] {
  --bg-primary: #ffffff;
  --bg-secondary: #f7f7f7;
  --bg-tertiary: #f0f0f0;
  --bg-elevated: #ffffff;
  --bg-hover: #f0f0f0;
  --border: #dddddd;
  --text-primary: #222222;
  --text-secondary: #494949;
  --text-muted: #717171;
  --accent: #e11d48;
  --accent-hover: #be123c;
  --accent-fg: #ffffff;
  --success: #059669;
  --danger: #dc2626;
  --warning: #d97706;
  --info: #2563eb;
}

/* amber - the storefront palette. Dark text on the accent rather than white:
 * amber is too light to carry white text at AA. */
[data-theme="amber"] {
  --bg-primary: #ffffff;
  --bg-secondary: #f7f8f8;
  --bg-tertiary: #eff1f1;
  --bg-elevated: #ffffff;
  --bg-hover: #eff1f1;
  --border: #d5d9d9;
  --text-primary: #0f1111;
  --text-secondary: #3f4545;
  --text-muted: #565959;
  --accent: #f59e0b;
  --accent-hover: #d97706;
  --accent-fg: #0f1111;
  --success: #047857;
  --danger: #b91c1c;
  --warning: #b45309;
  --info: #0369a1;
}

/* sky - crisp blue on cool greys. */
[data-theme="sky"] {
  --bg-primary: #ffffff;
  --bg-secondary: #f8fafc;
  --bg-tertiary: #eef4f9;
  --bg-elevated: #ffffff;
  --bg-hover: #eef4f9;
  --border: #dbe3ec;
  --text-primary: #0b1521;
  --text-secondary: #3a4a5e;
  --text-muted: #5b6b7f;
  --accent: #0284c7;
  --accent-hover: #0369a1;
  --accent-fg: #ffffff;
  --success: #059669;
  --danger: #dc2626;
  --warning: #d97706;
  --info: #0ea5e9;
}

/* mono - black and white. The accent is the text colour, which is the whole
 * idea: nothing on the screen competes for attention with the content. */
[data-theme="mono"] {
  --bg-primary: #ffffff;
  --bg-secondary: #fafafa;
  --bg-tertiary: #f5f5f5;
  --bg-elevated: #ffffff;
  --bg-hover: #f5f5f5;
  --border: #e5e5e5;
  --text-primary: #0a0a0a;
  --text-secondary: #525252;
  --text-muted: #737373;
  --accent: #0a0a0a;
  --accent-hover: #262626;
  --accent-fg: #ffffff;
  --success: #15803d;
  --danger: #b91c1c;
  --warning: #a16207;
  --info: #1d4ed8;
}

/* emerald - green on neutral greys. */
[data-theme="emerald"] {
  --bg-primary: #ffffff;
  --bg-secondary: #f9fafb;
  --bg-tertiary: #f3f4f6;
  --bg-elevated: #ffffff;
  --bg-hover: #f3f4f6;
  --border: #e5e7eb;
  --text-primary: #111827;
  --text-secondary: #4b5563;
  --text-muted: #6b7280;
  --accent: #059669;
  --accent-hover: #047857;
  --accent-fg: #ffffff;
  --success: #059669;
  --danger: #dc2626;
  --warning: #d97706;
  --info: #2563eb;
}

/* midnight — legacy v3.27 dark look. Opt in by setting THEME=midnight or
 * adding data-theme="midnight" on a specific surface. */
[data-theme="midnight"] {
  --bg-primary: #0a0a0f;
  --bg-secondary: #111118;
  --bg-tertiary: #1a1a24;
  --bg-elevated: #22222e;
  --bg-hover: #2a2a38;
  --border: #2a2a3a;
  --text-primary: #e8e8f0;
  --text-secondary: #9090a8;
  --text-muted: #7c7c96;
  --accent: #6c5ce7;
  --accent-hover: #7c6cf7;
  --accent-fg: #ffffff;
  --success: #00b894;
  --danger: #ff6b6b;
  --warning: #fdcb6e;
  --info: #74b9ff;
}

` + derivedTokenCSS()
}

// derivedTokenCSS is the second half of the palette: the tokens a theme does
// not state, worked out from the ones it does.
//
// The redesign asks for a recessed input fill, a stronger border, a pressed
// surface, and a tinted fill and hairline for the accent and each status
// colour. Written out per theme that is nine more tables of seven values, and
// the tenth theme somebody adds gets none of them: a pill with no fill is
// invisible, and nothing reports it. Derived, a new theme gets the whole set
// from the four colours it already declares.
//
// The selector is the load-bearing part. ":root" alone would resolve these
// once, against the root theme, and a section with its own data-theme would
// inherit that resolved value while its --accent changed underneath: the bug
// behind grit#98. "[data-theme]" matches the themed element itself, so each
// one resolves its own.
//
// color-mix into transparent rather than into a surface, because these sit on
// a card on one screen and on the page on the next, and a fill mixed into one
// of them is wrong on the other. The toaster's --success-bg and friends below
// stay opaque on purpose: a toast floats over content and has to hide it.
func derivedTokenCSS() string {
	return `/* Derived tokens. Every theme gets these from the colours it declares. */
:root,
[data-theme] {
  /* The fill of an input, a textarea, a select. Recessed from the card it
     sits on: in a dark theme the page is darker than the card, in a light
     theme it is the same white or a shade under it, and either way the
     input's border carries the edge. */
  --surface-input: var(--bg-primary);

  /* A row that is pressed or selected, one step past --bg-hover. */
  --bg-pressed: color-mix(in srgb, var(--bg-hover) 92%, var(--text-primary));

  /* The border of something you are pointing at, and of the dashed wells that
     say a thing is not configured. --border is too quiet to carry either. */
  --border-strong: color-mix(in srgb, var(--border) 70%, var(--text-muted));

  /* Tinted fill and hairline, per colour. 12% reads as a tint rather than a
     block at every lightness we ship; 32% is a visible edge that is still
     quieter than the text it surrounds. */
  --accent-soft: color-mix(in srgb, var(--accent) 12%, transparent);
  --accent-line: color-mix(in srgb, var(--accent) 32%, transparent);
  --success-soft: color-mix(in srgb, var(--success) 12%, transparent);
  --success-line: color-mix(in srgb, var(--success) 32%, transparent);
  --warning-soft: color-mix(in srgb, var(--warning) 12%, transparent);
  --warning-line: color-mix(in srgb, var(--warning) 32%, transparent);
  --danger-soft: color-mix(in srgb, var(--danger) 12%, transparent);
  --danger-line: color-mix(in srgb, var(--danger) 32%, transparent);
  --info-soft: color-mix(in srgb, var(--info) 12%, transparent);
  --info-line: color-mix(in srgb, var(--info) 32%, transparent);
}`
}

// themeVars pulls one theme's custom properties out of the shared palette.
//
// The nine palettes are CSS because that is what three of the four frontends
// consume. The mobile app is the fourth: NativeWind has no CSS variables, so
// its colours have to be literal values written at generation time. Reading
// them back out of the same block keeps one source of truth instead of a
// second table that drifts.
func themeVars(name string) map[string]string {
	want := `[data-theme="` + name + `"] {`
	css := themePaletteCSS()
	start := strings.Index(css, want)
	if start < 0 {
		if name == "atlas" {
			return map[string]string{}
		}
		return themeVars("atlas")
	}
	block := css[start+len(want):]
	if end := strings.Index(block, "}"); end >= 0 {
		block = block[:end]
	}

	vars := map[string]string{}
	for _, line := range strings.Split(block, "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(line), ";"))
		key, value, ok := strings.Cut(line, ":")
		if !ok || !strings.HasPrefix(key, "--") {
			continue
		}
		vars[strings.TrimPrefix(key, "--")] = strings.TrimSpace(value)
	}
	return vars
}

// themeColor is one value from a theme, with a fallback for when a palette ever
// loses a variable: a colour missing from a generated config is a build error in
// Tailwind and an invisible element in React Native, and neither says why.
func themeColor(theme, name, fallback string) string {
	if v, ok := themeVars(theme)[name]; ok && v != "" {
		return v
	}
	return fallback
}

// rgba turns a palette hex into a React Native colour string.
//
// NativeWind reaches most of the mobile app, but gradients, the status bar and
// the tab bar are set in JavaScript, and those need a value rather than a class.
func rgba(hex string, alpha float64) string {
	hex = strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	if len(hex) != 6 {
		return "rgba(0, 0, 0, " + strconv.FormatFloat(alpha, 'g', -1, 64) + ")"
	}
	n, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return "rgba(0, 0, 0, " + strconv.FormatFloat(alpha, 'g', -1, 64) + ")"
	}
	return fmt.Sprintf("rgba(%d, %d, %d, %s)", n>>16&0xff, n>>8&0xff, n&0xff,
		strconv.FormatFloat(alpha, 'g', -1, 64))
}

// themeTokenMapCSS maps the palette's tokens onto Tailwind's own token names,
// for every surface that paints with Tailwind utilities.
//
// Regular @theme, NOT @theme inline, and that is the load-bearing detail. The
// utilities compile to var(--color-*), whose values are var(--bg-*)
// indirections, so overriding --bg-* under [data-theme="..."] repaints the
// whole surface at runtime. @theme inline would substitute the values at build
// time and theme switching would silently stop working.
//
// One writer rather than two identical blocks in the admin's stylesheet and
// the web app's. They were already the same nineteen lines, and the scoped
// block below is the second thing in a row that had to be written twice to be
// right once.
func themeTokenMapCSS() string {
	return "@theme {\n" + themeTokenMap + `}

/* The same mappings again, for a theme scoped to a surface instead of to
 * <html>.
 *
 * Tailwind declares the @theme mappings above on :root, and a custom property
 * that references another one is resolved where it is declared. So a section
 * with its own data-theme redefines --bg-primary, but inherits a
 * --color-background that was already resolved at the root: the tokens change
 * and bg-background does not (grit#98). Re-declaring them here makes them
 * resolve against that element's own tokens, and its children inherit the
 * resolved value.
 *
 * Nested on purpose. A single [data-theme] is the page's own theme on <html>,
 * where the block above already does this. */
[data-theme] [data-theme],
[data-theme] .dark,
[data-theme] [data-theme-mode="dark"] {
` + themeTokenMap + `}
`
}

// themeTokenMap is the list, once, so the two blocks above cannot drift.
const themeTokenMap = `  --color-background: var(--bg-primary);
  --color-bg-secondary: var(--bg-secondary);
  --color-bg-tertiary: var(--bg-tertiary);
  --color-bg-elevated: var(--bg-elevated);
  --color-bg-hover: var(--bg-hover);
  --color-border: var(--border);
  --color-foreground: var(--text-primary);
  --color-text-secondary: var(--text-secondary);
  --color-text-muted: var(--text-muted);
  --color-accent: var(--accent);
  --color-accent-hover: var(--accent-hover);
  --color-accent-fg: var(--accent-fg);
  --color-success: var(--success);
  --color-danger: var(--danger);
  --color-warning: var(--warning);
  --color-info: var(--info);

  --color-surface-input: var(--surface-input);
  --color-bg-pressed: var(--bg-pressed);
  --color-border-strong: var(--border-strong);
  --color-accent-soft: var(--accent-soft);
  --color-accent-line: var(--accent-line);
  --color-success-soft: var(--success-soft);
  --color-success-line: var(--success-line);
  --color-warning-soft: var(--warning-soft);
  --color-warning-line: var(--warning-line);
  --color-danger-soft: var(--danger-soft);
  --color-danger-line: var(--danger-line);
  --color-info-soft: var(--info-soft);
  --color-info-line: var(--info-line);

  --font-sans: var(--font-display), system-ui, sans-serif;
  --font-mono: var(--font-mono), ui-monospace, monospace;
  --font-serif: var(--font-serif), Georgia, serif;
`
