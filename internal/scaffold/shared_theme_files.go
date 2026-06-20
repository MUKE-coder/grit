package scaffold

import "fmt"

// sharedBrandConfig generates packages/shared/brand.config.ts — the single
// source of truth for the scaffolded app's identity: name, tagline, logo
// path, hero copy, hero image set (used by the Pulse theme's carousel), and
// social links. Auth pages and dashboards across apps/admin and apps/web
// import this so a rebrand is one edit, not a grep + replace.
//
// The shape is intentionally a plain object (not a function or factory) so
// IDE autocomplete works and Next.js can tree-shake unused fields.
func sharedBrandConfig(opts Options) string {
	return fmt.Sprintf(`// brand.config.ts — single source of truth for your app's identity.
//
// Imported across apps/admin and apps/web by auth pages, dashboards,
// emails, and metadata. Edit this file once to rebrand the entire app:
// name, tagline, logo, hero photography, social links.
//
// Want runtime-driven branding (set from .env or a Settings UI)? Mirror
// the fields below into env variables and read them with process.env at
// build time. The file shape stays the same.

export const brand = {
  /** Display name shown on auth pages, the sidebar header, and emails. */
  name: %q,

  /** One-line hero headline on the auth split panel. Keep it under 50 chars. */
  tagline: "Manage everything\nin one place.",

  /** Secondary line under the tagline. 1–2 short sentences. */
  description:
    "The admin dashboard for your application. Monitor, manage, and control your platform from a single screen.",

  /** Logo — single-character fallback when no image is set, plus an
   *  optional image path resolved from /public. */
  logo: {
    text: %q,
    image: "" as string | "",
  },

  /** Hero imagery — used by the Pulse theme's auth carousel and as a
   *  fallback wallpaper by Atlas/Aurora when set. Paths resolve from
   *  /public. Provide at least 3 for the carousel to feel alive. */
  hero: {
    images: [
      "/hero/01.jpg",
      "/hero/02.jpg",
      "/hero/03.jpg",
    ] as string[],
    /** Carousel rotation interval in ms. */
    intervalMs: 5000,
  },

  /** Optional brand color overrides. Leave empty strings to inherit the
   *  active theme's palette. Useful for keeping the theme structure but
   *  swapping just the primary accent. */
  colors: {
    primary: "" as string | "",
    accent: "" as string | "",
  },

  /** Social links — surfaced in the auth page footer and the dashboard
   *  user menu. Leave empty to hide. */
  social: {
    twitter: "",
    linkedin: "",
    github: "",
    youtube: "",
  },

  /** Legal links — shown in the auth footer for compliance. */
  legal: {
    termsUrl: "/terms",
    privacyUrl: "/privacy",
  },
} as const;

export type Brand = typeof brand;
`, displayName(opts.ProjectName), firstUpper(opts.ProjectName)[:1])
}

// firstUpper returns name with its first character uppercased.
// Used for the single-char logo fallback.
func firstUpper(name string) string {
	if name == "" {
		return "A"
	}
	if name[0] >= 'a' && name[0] <= 'z' {
		return string(name[0]-32) + name[1:]
	}
	return name
}

// displayName turns "my-app" into "My App" for the brand.name default.
func displayName(name string) string {
	if name == "" {
		return "Acme"
	}
	out := []byte(name)
	upper := true
	for i, b := range out {
		if b == '-' || b == '_' {
			out[i] = ' '
			upper = true
			continue
		}
		if upper && b >= 'a' && b <= 'z' {
			out[i] = b - 32
			upper = false
		} else {
			upper = false
		}
	}
	return string(out)
}

// sharedThemes generates packages/shared/themes.ts — the three theme
// palettes shipped with Grit v3.28+. Each theme is a flat token bag that
// auth pages and dashboards consume directly. No CSS-in-JS, no runtime
// theme provider; Tailwind classes pick up the tokens via inline style or
// CSS variables set at the layout root.
//
// Themes:
//   - atlas:  team/organisation. Light, sharp, Inter. Split-screen auth.
//   - aurora: consumer SaaS. Pastel, friendly, Geist. Centered auth.
//   - pulse:  ecommerce/brand. Warm, bold, Onest + DM Serif. Split with carousel.
func sharedThemes() string {
	return `// themes.ts — Grit v3.28 theme token palettes.
//
// Three theme variants ship out of the box. The active theme is picked at
// scaffold time via 'grit new --theme=<name>' and can be overridden at
// runtime via THEME=<name> in .env. Apps consume tokens through getTheme().
//
// Adding a new theme: extend ThemeName, add an entry to themes, restart
// the dev servers so the env -> token lookup picks it up.

export type ThemeName = "atlas" | "aurora" | "pulse";

export type AuthLayout = "split-static" | "split-carousel" | "centered";

export interface ThemeFonts {
  /** Font family used for body text and form inputs. */
  ui: string;
  /** Font family used for headings and display copy. */
  display: string;
  /** Optional monospace family for code and numerics. */
  mono?: string;
}

export interface ThemeColors {
  /** Page background. */
  bg: string;
  /** Default text color on bg. */
  fg: string;
  /** Card/surface background, slightly raised from bg. */
  card: string;
  /** Card border / divider line. */
  border: string;
  /** Muted/secondary text. */
  muted: string;
  /** Brand primary action color (CTAs, links, focus rings). */
  primary: string;
  /** Foreground color that pairs with primary backgrounds. */
  primaryFg: string;
  /** Brand accent for highlights, chips, badges. */
  accent: string;
  /** Hero panel background on split-screen auth layouts. */
  heroBg: string;
  /** Hero panel foreground. */
  heroFg: string;
}

export interface ThemeTokens {
  name: ThemeName;
  fonts: ThemeFonts;
  colors: ThemeColors;
  /** Border radius for cards, buttons, inputs. */
  radius: string;
  /** Auth page layout this theme is designed around. */
  authLayout: AuthLayout;
}

export const themes: Record<ThemeName, ThemeTokens> = {
  atlas: {
    name: "atlas",
    fonts: {
      ui: '"Inter", system-ui, -apple-system, sans-serif',
      display: '"Inter Display", "Inter", system-ui, sans-serif',
    },
    colors: {
      bg: "#ffffff",
      fg: "#0f172a",
      card: "#f8fafc",
      border: "#e2e8f0",
      muted: "#64748b",
      primary: "#2563eb",
      primaryFg: "#ffffff",
      accent: "#4f46e5",
      heroBg: "#4f46e5",
      heroFg: "#ffffff",
    },
    radius: "0.625rem",
    authLayout: "split-static",
  },

  aurora: {
    name: "aurora",
    fonts: {
      ui: '"Geist", "Inter", system-ui, sans-serif',
      display: '"Geist", "Inter", system-ui, sans-serif',
    },
    colors: {
      bg: "#fafaf9",
      fg: "#1c1917",
      card: "#ffffff",
      border: "#e7e5e4",
      muted: "#78716c",
      primary: "#7c3aed",
      primaryFg: "#ffffff",
      accent: "#a855f7",
      // Aurora's centered layout doesn't use a hero panel, so heroBg
      // is reused as the page wallpaper for the pastel background.
      heroBg: "#ede9fe",
      heroFg: "#1c1917",
    },
    radius: "0.75rem",
    authLayout: "centered",
  },

  pulse: {
    name: "pulse",
    fonts: {
      ui: '"Onest", "Inter", system-ui, sans-serif',
      display: '"DM Serif Display", "Onest", serif',
    },
    colors: {
      bg: "#fafaf9",
      fg: "#0f0f0f",
      card: "#ffffff",
      border: "#e7e5e4",
      muted: "#737373",
      primary: "#0f0f0f",
      primaryFg: "#ffffff",
      accent: "#fbbf24",
      heroBg: "#f5f5f4",
      heroFg: "#0f0f0f",
    },
    radius: "0.5rem",
    authLayout: "split-carousel",
  },
};

/**
 * getTheme resolves a name to a token bag, falling back to atlas on
 * anything unknown so a typo in .env can't crash render. Pass the value of
 * process.env.NEXT_PUBLIC_THEME or read from a server component env helper.
 */
export function getTheme(name?: string): ThemeTokens {
  const key = (name || "").toLowerCase() as ThemeName;
  return themes[key] || themes.atlas;
}

/**
 * The default theme baked at scaffold time. Auth and dashboard code can
 * import this directly when there's no runtime env to read (server
 * components, build-time metadata).
 */
export const defaultTheme = getTheme(
  typeof process !== "undefined" ? process.env.NEXT_PUBLIC_THEME : undefined
);

/**
 * isSocialAuthEnabled reads NEXT_PUBLIC_SOCIAL_AUTH_ENABLED at the call
 * site so auth pages can conditionally render Google/GitHub buttons. The
 * env is set in .env as SOCIAL_AUTH_ENABLED=true|false; Next.js exposes
 * it on the client through the NEXT_PUBLIC_ prefix wired in next.config.
 */
export function isSocialAuthEnabled(): boolean {
  const v = typeof process !== "undefined"
    ? process.env.NEXT_PUBLIC_SOCIAL_AUTH_ENABLED
    : undefined;
  // Default to enabled when the env is unset — matches v3.27 behavior.
  if (v === undefined || v === "") return true;
  return v === "true" || v === "1";
}
`
}
