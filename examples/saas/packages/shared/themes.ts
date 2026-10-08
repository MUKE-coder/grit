// themes.ts — Grit v3.28 theme token palettes.
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
