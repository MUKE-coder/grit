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
	return fmt.Sprintf(`// brand.config.ts: single source of truth for your app's identity.
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

  /** Logo — single-character fallback when no image is set. The full
   *  variant is used in the expanded sidebar + auth pages; the mark
   *  variant is used in the collapsed sidebar and favicons. Both are
   *  optional and resolve from /public. */
  logo: {
    text: %q,
    /** Wide logo (text + icon) shown in the expanded sidebar. */
    image: "" as string | "",
    /** Square mark shown in the collapsed sidebar, ~32x32. */
    mark: "" as string | "",
  },

  /** Hero imagery — used by the Pulse theme's auth carousel and as a
   *  fallback wallpaper by Atlas/Aurora when set. Paths resolve from
   *  /public.
   *
   *  Empty by default: the scaffold ships no photography, so listing paths
   *  here would 404 on the login page of every new project. Drop your own
   *  images into public/hero/ and add them below — provide at least 3 for
   *  the carousel to feel alive. While this is empty the auth screens fall
   *  back to a themed gradient. */
  hero: {
    images: [] as string[],
    /** Carousel rotation interval in ms. */
    intervalMs: 5000,
  },

  /** Copy for the auth layouts that put proof beside the form: the mono
   *  theme's showcase panel and the emerald theme's quote. Every value is
   *  yours to rewrite, and the layouts fall back to these rather than
   *  rendering an empty panel.
   *
   *  The quote is deliberately not attributed to a real person or company.
   *  A scaffold that ships a fabricated testimonial with a plausible name is
   *  a scaffold that ships a lie to production the first time somebody
   *  forgets to edit it. */
  proof: {
    headline: "Built on Grit.",
    subheadline: "Go on the back, React on the front, one command to ship.",
    points: [
      "Auth, RBAC and an admin panel on the first run",
      "One API for web, mobile and desktop",
      "Deploy with a single command",
    ] as string[],
    quote: "We replaced four services and a month of wiring with one command.",
    quoteAuthor: "Replace this with a real one",
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

export type ThemeName =
  | "atlas"
  | "aurora"
  | "pulse"
  | "coral"
  | "amber"
  | "sky"
  | "mono"
  | "emerald";

// The shape of the sign-in screen. A theme picks one; the form inside is the
// same in every layout, so adding one of these never touches a form.
export type AuthLayout =
  | "split-static"
  | "split-carousel"
  | "centered"
  | "modal"
  | "boxed"
  | "banner"
  | "showcase"
  | "quote";

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
      bg: "#fbfbfd",
      fg: "#1d1d1f",
      card: "#ffffff",
      border: "#d2d2d7",
      muted: "#86868b",
      primary: "#1d1d1f",
      primaryFg: "#ffffff",
      accent: "#1d1d1f",
      // Aurora's centered layout has no hero panel, so heroBg is the page
      // wallpaper — Apple's light section grey.
      heroBg: "#f5f5f7",
      heroFg: "#1d1d1f",
    },
    radius: "0.75rem",
    authLayout: "centered",
  },

  pulse: {
    name: "pulse",
    fonts: {
      ui: '"Onest", "Inter", system-ui, sans-serif',
      display: '"Onest", "Inter", system-ui, sans-serif',
    },
    colors: {
      bg: "#f6f7f9",
      fg: "#1d1f26",
      card: "#ffffff",
      border: "#e0e4e9",
      muted: "#8a94a6",
      primary: "#0051c3",
      primaryFg: "#ffffff",
      accent: "#f6821f",
      // Deep Cloudflare blue hero panel behind the split-carousel — premium,
      // white text.
      heroBg: "#003682",
      heroFg: "#ffffff",
    },
    radius: "0.5rem",
    authLayout: "split-carousel",
  },

  // Coral: a card floating over a blurred glimpse of the app behind it, the
  // shape a marketplace uses when signing in is an interruption rather than a
  // destination.
  coral: {
    name: "coral",
    fonts: {
      ui: '"Inter", system-ui, -apple-system, sans-serif',
      display: '"Inter", system-ui, -apple-system, sans-serif',
    },
    colors: {
      bg: "#f7f7f7",
      fg: "#222222",
      card: "#ffffff",
      border: "#dddddd",
      muted: "#717171",
      primary: "#e11d48",
      primaryFg: "#ffffff",
      accent: "#f43f5e",
      heroBg: "#fff1f2",
      heroFg: "#222222",
    },
    radius: "0.75rem",
    authLayout: "modal",
  },

  // Amber: a plain bordered box under a wordmark, with the legal text and the
  // footer links a storefront is obliged to carry. Deliberately unfashionable.
  amber: {
    name: "amber",
    fonts: {
      ui: '"Inter", Arial, system-ui, sans-serif',
      display: '"Inter", Arial, system-ui, sans-serif',
    },
    colors: {
      bg: "#ffffff",
      fg: "#0f1111",
      card: "#ffffff",
      border: "#d5d9d9",
      muted: "#565959",
      primary: "#f59e0b",
      primaryFg: "#0f1111",
      accent: "#b45309",
      heroBg: "#fffbeb",
      heroFg: "#0f1111",
    },
    radius: "0.375rem",
    authLayout: "boxed",
  },

  // Sky: a top bar and one bold left-aligned heading, social sign-in first.
  // For products where most people arrive with an existing identity.
  sky: {
    name: "sky",
    fonts: {
      ui: '"Inter", system-ui, -apple-system, sans-serif',
      display: '"Inter Display", "Inter", system-ui, sans-serif',
    },
    colors: {
      bg: "#ffffff",
      fg: "#0b1521",
      card: "#f8fafc",
      border: "#dbe3ec",
      muted: "#5b6b7f",
      primary: "#0284c7",
      primaryFg: "#ffffff",
      accent: "#0ea5e9",
      heroBg: "#e0f2fe",
      heroFg: "#0b1521",
    },
    radius: "0.5rem",
    authLayout: "banner",
  },

  // Mono: black and white on a fine grid, with a showcase panel beside the
  // form. The look developer tools reach for when the product is the proof.
  mono: {
    name: "mono",
    fonts: {
      ui: '"Inter", system-ui, -apple-system, sans-serif',
      display: '"Inter", system-ui, -apple-system, sans-serif',
    },
    colors: {
      bg: "#ffffff",
      fg: "#0a0a0a",
      card: "#ffffff",
      border: "#e5e5e5",
      muted: "#737373",
      primary: "#0a0a0a",
      primaryFg: "#ffffff",
      accent: "#404040",
      heroBg: "#fafafa",
      heroFg: "#0a0a0a",
    },
    radius: "0.5rem",
    authLayout: "showcase",
  },

  // Emerald: a narrow form column with a customer quote filling the rest.
  // The quote is the argument, so it gets the larger half.
  emerald: {
    name: "emerald",
    fonts: {
      ui: '"Inter", system-ui, -apple-system, sans-serif',
      display: '"Inter", system-ui, -apple-system, sans-serif',
    },
    colors: {
      bg: "#ffffff",
      fg: "#111827",
      card: "#ffffff",
      border: "#e5e7eb",
      muted: "#6b7280",
      primary: "#059669",
      primaryFg: "#ffffff",
      accent: "#10b981",
      heroBg: "#ecfdf5",
      heroFg: "#064e3b",
    },
    radius: "0.5rem",
    authLayout: "quote",
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
