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
