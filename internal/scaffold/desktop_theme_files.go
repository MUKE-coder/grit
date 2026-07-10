package scaffold

import "fmt"

// ═══════════════════════════════════════════════════════════════════
// Desktop theming — shares packages/shared/themes.ts with the admin
// panel so `grit new --theme=atlas|aurora|pulse` styles the desktop
// client identically.
//
// The desktop's Tailwind config already maps its colour utilities onto
// CSS variables (--bg-primary, --text-foreground, --accent, …). Driving
// those variables from the shared token bag means the dashboard,
// settings, sidebar, topbar AND every generated resource screen adopt
// the active theme with no per-page changes. The auth pages get the
// same three shells the admin uses (split-static / centered /
// split-carousel).
// ═══════════════════════════════════════════════════════════════════

// desktopThemeFontURL returns the Google Fonts stylesheet the active theme
// needs. Themes pick different families (atlas: Inter, aurora: Geist,
// pulse: Onest + DM Serif Display); JetBrains Mono is always loaded for
// numerics and code.
func desktopThemeFontURL(theme string) string {
	families := "family=Inter:wght@400;500;600;700"
	switch theme {
	case "aurora":
		families = "family=Geist:wght@400;500;600;700"
	case "pulse":
		families = "family=Onest:wght@400;500;600;700&family=DM+Serif+Display"
	}
	return "https://fonts.googleapis.com/css2?" + families +
		"&family=JetBrains+Mono:wght@400;500;600&display=swap"
}

// desktopClientThemeTokens generates lib/theme-tokens.ts — the bridge from
// the shared theme palette to the CSS variables the desktop's Tailwind
// config consumes.
func desktopClientThemeTokens(opts Options) string {
	theme := opts.Theme
	if theme == "" {
		theme = "atlas"
	}

	return fmt.Sprintf(`import { getTheme, type ThemeTokens } from "@repo/shared/themes";

// The theme chosen at scaffold time (grit new --theme=%s). Overridable at
// runtime with VITE_THEME in the frontend env.
//
// NOTE: the shared defaultTheme export reads process.env.NEXT_PUBLIC_THEME,
// which doesn't exist under Vite — so resolve explicitly here.
const SCAFFOLD_THEME = %q;

export const activeTheme: ThemeTokens = getTheme(
  (import.meta.env.VITE_THEME as string | undefined) || SCAFFOLD_THEME,
);

export type ColorMode = "light" | "dark";

// themeCssVars maps the shared token bag onto the CSS variables declared in
// globals.css and consumed by tailwind.config.ts. Because every desktop
// surface (dashboard, settings, sidebar, topbar, generated resource
// screens) is already written against these variables, setting them here is
// all it takes for the whole app to adopt the theme.
//
// Light mode uses the theme's own palette verbatim, so it matches the admin
// panel. Dark mode keeps neutral dark surfaces — a desktop app gets used at
// night and a "dark" atlas would just be white — but adopts the theme's
// brand colours so the two modes read as one product.
export function themeCssVars(mode: ColorMode): Record<string, string> {
  const c = activeTheme.colors;

  const base: Record<string, string> = {
    "--accent": c.primary,
    "--accent-hover": c.accent,
    "--primary-fg": c.primaryFg,
    "--hero-bg": c.heroBg,
    "--hero-fg": c.heroFg,
    "--font-ui": activeTheme.fonts.ui,
    "--font-display": activeTheme.fonts.display,
    "--radius": activeTheme.radius,
  };

  if (mode === "light") {
    return {
      ...base,
      "--bg-primary": c.bg,
      "--bg-secondary": c.card,
      "--bg-tertiary": c.card,
      "--bg-elevated": c.bg,
      "--bg-hover": c.border,
      "--border": c.border,
      "--border-subtle": c.border,
      "--text-foreground": c.fg,
      "--text-secondary": c.muted,
      "--text-muted": c.muted,
    };
  }

  return {
    ...base,
    "--bg-primary": "#0a0a0f",
    "--bg-secondary": "#111118",
    "--bg-tertiary": "#1a1a24",
    "--bg-elevated": "#22222e",
    "--bg-hover": "#2a2a38",
    "--border": "#2a2a3a",
    "--border-subtle": "#1f1f2b",
    "--text-foreground": "#e8e8f0",
    "--text-secondary": "#9090a8",
    "--text-muted": "#606078",
  };
}

// applyThemeVars writes the variables onto <html>. Called by ThemeProvider
// whenever the light/dark mode changes.
export function applyThemeVars(mode: ColorMode): void {
  const root = document.documentElement;
  const vars = themeCssVars(mode);
  for (const key of Object.keys(vars)) {
    root.style.setProperty(key, vars[key]);
  }
}
`, theme, theme)
}

// desktopClientAuthShell generates components/auth/AuthShell.tsx — the
// dispatcher that picks the shell matching the active theme's authLayout.
// Mirrors the admin panel's AuthShell so both apps present the same login.
func desktopClientAuthShell() string {
	return `import type { ReactNode } from "react";
import { activeTheme } from "@/lib/theme-tokens";
import { AtlasAuthShell } from "./AtlasAuthShell";
import { AuroraAuthShell } from "./AuroraAuthShell";
import { PulseAuthShell } from "./PulseAuthShell";

export type AuthMode = "login" | "sign-up";

export interface AuthShellProps {
  mode: AuthMode;
  title: string;
  subtitle?: string;
  children: ReactNode;
  errorMessage?: string;
}

// AuthShell renders the layout the active theme was designed around:
//   atlas  -> split-static    (hero panel left, form right)
//   aurora -> centered        (single card on a pastel wallpaper)
//   pulse  -> split-carousel  (form left, editorial hero right)
export function AuthShell(props: AuthShellProps) {
  const theme = activeTheme;
  switch (theme.authLayout) {
    case "centered":
      return <AuroraAuthShell theme={theme} {...props} />;
    case "split-carousel":
      return <PulseAuthShell theme={theme} {...props} />;
    case "split-static":
    default:
      return <AtlasAuthShell theme={theme} {...props} />;
  }
}

// Shared prop contract for the three shells.
export interface ShellProps extends AuthShellProps {
  theme: typeof activeTheme;
}

// switchLinks maps the auth mode to its "other" page. Desktop routes are
// /auth/login and /auth/register (the admin uses /login and /sign-up).
export const switchLinks: Record<AuthMode, { hint: string; to: string; label: string }> = {
  login: { hint: "Don't have an account?", to: "/auth/register", label: "Create one" },
  "sign-up": { hint: "Already have an account?", to: "/auth/login", label: "Sign in" },
};

// authVars are the CSS variables the form inputs read, so one set of inputs
// fits every theme without per-shell styling.
export function authVars(theme: typeof activeTheme): Record<string, string> {
  const t = theme.colors;
  return {
    "--auth-bg": t.bg,
    "--auth-fg": t.fg,
    "--auth-card": t.card,
    "--auth-border": t.border,
    "--auth-muted": t.muted,
    "--auth-primary": t.primary,
    "--auth-primary-fg": t.primaryFg,
    "--auth-accent": t.accent,
    "--auth-radius": theme.radius,
  };
}
`
}

// desktopClientBrandMark generates components/auth/BrandMark.tsx — the small
// logo lockup used by all three shells.
func desktopClientBrandMark() string {
	return `import { brand } from "@repo/shared/brand.config";

export function BrandMark({ tint }: { tint?: string }) {
  if (brand.logo.image) {
    return <img src={brand.logo.image} alt={brand.name} className="h-8 w-8" />;
  }
  return (
    <span
      className="inline-flex h-8 w-8 items-center justify-center rounded-md font-bold text-white"
      style={{ background: tint || "rgba(255,255,255,0.15)" }}
    >
      {brand.logo.text}
    </span>
  );
}
`
}

// desktopClientAuthField generates components/auth/AuthField.tsx — the input
// and submit-button styling shared by login/register. Everything reads the
// --auth-* variables the shell sets, so one set of styles fits all three
// themes without per-theme forms.
func desktopClientAuthField() string {
	return `import type { CSSProperties, ReactNode } from "react";

export const authInputCls =
  "w-full h-11 rounded-[var(--auth-radius)] border px-3.5 text-[14px] outline-none transition-colors focus:ring-2";

export const authInputStyle: CSSProperties = {
  borderColor: "var(--auth-border)",
  background: "var(--auth-card)",
  color: "var(--auth-fg)",
};

export function AuthSubmit({ disabled, children }: { disabled?: boolean; children: ReactNode }) {
  return (
    <button
      type="submit"
      disabled={disabled}
      className="flex h-11 w-full items-center justify-center gap-2 rounded-[var(--auth-radius)] text-[15px] font-semibold transition-opacity hover:opacity-90 disabled:opacity-50"
      style={{ background: "var(--auth-primary)", color: "var(--auth-primary-fg)" }}
    >
      {children}
    </button>
  );
}
`
}

// desktopAtlasAuthShell — split-static: hero panel on the left, form right.
// The default, and a direct port of the admin's AtlasAuthShell.
func desktopAtlasAuthShell() string {
	return `import { Link } from "@tanstack/react-router";
import { brand } from "@repo/shared/brand.config";
import { BrandMark } from "./BrandMark";
import { authVars, switchLinks, type ShellProps } from "./AuthShell";

export function AtlasAuthShell({ theme, mode, title, subtitle, children, errorMessage }: ShellProps) {
  const t = theme.colors;
  const f = theme.fonts;
  const sw = switchLinks[mode];

  return (
    <div
      className="flex min-h-full"
      style={{
        fontFamily: f.ui,
        background: t.bg,
        color: t.fg,
        ...authVars(theme),
      } as React.CSSProperties}
    >
      {/* Left hero panel */}
      <div
        className="hidden lg:flex lg:w-1/2 flex-col justify-between p-12"
        style={{ background: t.heroBg, color: t.heroFg }}
      >
        <div className="flex items-center gap-2 text-2xl font-bold">
          <BrandMark />
          <span style={{ fontFamily: f.display }}>{brand.name}</span>
        </div>

        <div className="space-y-4 max-w-md">
          <h1
            className="text-4xl font-bold leading-tight whitespace-pre-line"
            style={{ fontFamily: f.display }}
          >
            {brand.tagline}
          </h1>
          <p className="text-lg opacity-80">{brand.description}</p>
        </div>

        <p className="text-sm opacity-60">Built with Grit — Go + React framework</p>
      </div>

      {/* Right form panel */}
      <div className="flex flex-1 items-center justify-center px-6 py-12" style={{ background: t.bg }}>
        <div className="w-full max-w-md space-y-8">
          <div className="lg:hidden flex items-center justify-center gap-2 text-2xl font-bold" style={{ color: t.primary }}>
            <BrandMark tint={t.primary} />
            <span style={{ fontFamily: f.display }}>{brand.name}</span>
          </div>

          <div>
            <h2 className="text-2xl font-bold" style={{ fontFamily: f.display }}>{title}</h2>
            {subtitle && <p className="mt-2" style={{ color: t.muted }}>{subtitle}</p>}
          </div>

          {errorMessage && (
            <div
              className="rounded-[var(--auth-radius)] border px-4 py-3 text-sm"
              style={{ borderColor: "#fecaca", background: "#fef2f2", color: "#b91c1c" }}
            >
              {errorMessage}
            </div>
          )}

          {children}

          <p className="text-center text-sm" style={{ color: t.muted }}>
            {sw.hint}{" "}
            <Link to={sw.to} className="font-medium" style={{ color: t.primary }}>
              {sw.label}
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
`
}

// desktopAuroraAuthShell — centered: one card on a pastel wallpaper.
func desktopAuroraAuthShell() string {
	return `import { Link } from "@tanstack/react-router";
import { BrandMark } from "./BrandMark";
import { authVars, switchLinks, type ShellProps } from "./AuthShell";

export function AuroraAuthShell({ theme, mode, title, subtitle, children, errorMessage }: ShellProps) {
  const t = theme.colors;
  const f = theme.fonts;
  const sw = switchLinks[mode];

  // Soft pastel wallpaper: two radial pulses over the heroBg token so the
  // card lifts off the background. Built with concatenation rather than a
  // template literal to stay inside the scaffold's raw strings.
  const wallpaper =
    "radial-gradient(60% 60% at 30% 20%, " + t.accent + "22, transparent 60%), " +
    "radial-gradient(50% 50% at 80% 80%, " + t.primary + "1a, transparent 60%), " +
    t.heroBg;

  return (
    <div
      className="flex min-h-full items-center justify-center px-4 py-12"
      style={{
        fontFamily: f.ui,
        background: wallpaper,
        color: t.fg,
        ...authVars(theme),
      } as React.CSSProperties}
    >
      <div
        className="w-full max-w-md rounded-2xl border shadow-xl p-8 space-y-6"
        style={{ background: t.card, borderColor: t.border }}
      >
        <div className="flex flex-col items-center text-center space-y-3">
          <BrandMark tint={t.primary} />
          <div>
            <h2 className="text-2xl font-bold" style={{ fontFamily: f.display }}>{title}</h2>
            {subtitle && <p className="mt-1 text-sm" style={{ color: t.muted }}>{subtitle}</p>}
          </div>
        </div>

        {errorMessage && (
          <div
            className="rounded-[var(--auth-radius)] border px-4 py-3 text-sm"
            style={{ borderColor: "#fecaca", background: "#fef2f2", color: "#b91c1c" }}
          >
            {errorMessage}
          </div>
        )}

        {children}

        <p className="text-center text-sm" style={{ color: t.muted }}>
          {sw.hint}{" "}
          <Link to={sw.to} className="font-medium" style={{ color: t.primary }}>
            {sw.label}
          </Link>
        </p>
      </div>
    </div>
  );
}
`
}

// desktopPulseAuthShell — split-carousel: form left, editorial hero right.
func desktopPulseAuthShell() string {
	return `import { Link } from "@tanstack/react-router";
import { brand } from "@repo/shared/brand.config";
import { BrandMark } from "./BrandMark";
import { authVars, switchLinks, type ShellProps } from "./AuthShell";

export function PulseAuthShell({ theme, mode, title, subtitle, children, errorMessage }: ShellProps) {
  const t = theme.colors;
  const f = theme.fonts;
  const sw = switchLinks[mode];

  return (
    <div
      className="flex min-h-full"
      style={{
        fontFamily: f.ui,
        background: t.bg,
        color: t.fg,
        ...authVars(theme),
      } as React.CSSProperties}
    >
      {/* Left form panel */}
      <div className="flex flex-1 items-center justify-center px-6 py-12 lg:px-12">
        <div className="w-full max-w-md space-y-8">
          <div className="flex items-center gap-2 text-xl font-bold" style={{ color: t.fg }}>
            <BrandMark tint={t.primary} />
            <span style={{ fontFamily: f.display }}>{brand.name}</span>
          </div>

          <div>
            <h2 className="text-3xl font-bold" style={{ fontFamily: f.display }}>{title}</h2>
            {subtitle && <p className="mt-2" style={{ color: t.muted }}>{subtitle}</p>}
          </div>

          {errorMessage && (
            <div
              className="rounded-[var(--auth-radius)] border px-4 py-3 text-sm"
              style={{ borderColor: "#fecaca", background: "#fef2f2", color: "#b91c1c" }}
            >
              {errorMessage}
            </div>
          )}

          {children}

          <p className="text-sm" style={{ color: t.muted }}>
            {sw.hint}{" "}
            <Link to={sw.to} className="font-medium underline" style={{ color: t.fg }}>
              {sw.label}
            </Link>
          </p>
        </div>
      </div>

      {/* Right editorial hero */}
      <div
        className="hidden lg:flex lg:w-1/2 flex-col justify-between p-12"
        style={{ background: t.heroBg, color: t.heroFg }}
      >
        <span
          className="self-end rounded-full px-3 py-1 text-xs font-semibold"
          style={{ background: t.accent, color: t.fg }}
        >
          {brand.name}
        </span>

        <div className="space-y-6 max-w-md">
          <h1
            className="text-5xl leading-[1.05] whitespace-pre-line"
            style={{ fontFamily: f.display }}
          >
            {brand.tagline}
          </h1>
          <p className="text-lg opacity-70">{brand.description}</p>
          <div className="h-1 w-24 rounded-full" style={{ background: t.accent }} />
        </div>

        <p className="text-sm opacity-50">Built with Grit — Go + React framework</p>
      </div>
    </div>
  );
}
`
}
