package scaffold

// This file generates the v3.28 theme-aware auth surface for apps/admin:
//
//   components/auth/AuthShell.tsx        — dispatcher, picks shell by theme
//   components/auth/AtlasAuthShell.tsx   — split-static, blue/white (default)
//   components/auth/AuroraAuthShell.tsx  — centered, Clerk/WorkOS style
//   components/auth/PulseAuthShell.tsx   — split-screen with hero carousel
//   components/auth/SocialAuthButtons.tsx — Google + GitHub, hides when
//                                          NEXT_PUBLIC_SOCIAL_AUTH_ENABLED=false
//
// Each shell takes a `mode` (login | sign-up | forgot | reset) so a single
// shell renders the right hero copy and right-rail link variants. The form
// content itself is passed as children — pages stay thin.

// adminAuthShellDispatcher returns the entry point components/auth/AuthShell.tsx.
// It reads the active theme from packages/shared/themes and renders the
// matching shell. This is the only file the page-level auth routes import.
func adminAuthShellDispatcher() string {
	return `"use client";

import type { ReactNode } from "react";
import { getTheme } from "@repo/shared/themes";
import { AtlasAuthShell } from "./AtlasAuthShell";
import { AuroraAuthShell } from "./AuroraAuthShell";
import { PulseAuthShell } from "./PulseAuthShell";

export type AuthMode = "login" | "sign-up" | "forgot" | "reset";

export interface AuthShellProps {
  mode: AuthMode;
  title: string;
  subtitle?: string;
  children: ReactNode;
  /** Server-side error from useLogin / useRegister etc. — shown in a banner. */
  errorMessage?: string;
  /** Show social auth buttons + divider. Defaults to true. */
  showSocial?: boolean;
}

// AuthShell is the single entrypoint pages use. It reads
// NEXT_PUBLIC_THEME, resolves the theme tokens, then dispatches to the
// shell that matches the theme's authLayout. Centralising the lookup
// means a page never knows which theme is active — easier to add a new
// theme later.
export function AuthShell(props: AuthShellProps) {
  const theme = getTheme(process.env.NEXT_PUBLIC_THEME);
  switch (theme.authLayout) {
    case "centered":
      return <AuroraAuthShell {...props} theme={theme} />;
    case "split-carousel":
      return <PulseAuthShell {...props} theme={theme} />;
    case "split-static":
    default:
      return <AtlasAuthShell {...props} theme={theme} />;
  }
}
`
}

// adminAuthSocialButtons returns components/auth/SocialAuthButtons.tsx.
// One source of truth for the Google + GitHub buttons — all three shells
// embed the same component so SOCIAL_AUTH_ENABLED=false hides them in
// every theme without per-shell edits. Returns null when disabled, which
// also lets the shells skip the "or continue with" divider via a sibling
// isSocialAuthEnabled() check.
func adminAuthSocialButtons() string {
	return `"use client";

import { isSocialAuthEnabled } from "@repo/shared/themes";

export function SocialAuthButtons() {
  if (!isSocialAuthEnabled()) return null;

  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  return (
    <div className="flex gap-3">
      <a
        href={apiUrl + "/api/auth/oauth/google"}
        className="flex flex-1 items-center justify-center gap-2 rounded-[var(--auth-radius)] border border-[var(--auth-border)] bg-white px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
      >
        <svg className="h-5 w-5" viewBox="0 0 24 24">
          <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4" />
          <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
          <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
          <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
        </svg>
        Google
      </a>
      <a
        href={apiUrl + "/api/auth/oauth/github"}
        className="flex flex-1 items-center justify-center gap-2 rounded-[var(--auth-radius)] border border-[var(--auth-border)] bg-[#24292f] px-4 py-2.5 text-sm font-medium text-white hover:bg-[#2f363d] transition-colors"
      >
        <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
          <path fillRule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0 1 12 6.844a9.59 9.59 0 0 1 2.504.337c1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.02 10.02 0 0 0 22 12.017C22 6.484 17.522 2 12 2z" clipRule="evenodd" />
        </svg>
        GitHub
      </a>
    </div>
  );
}

export function SocialAuthDivider() {
  if (!isSocialAuthEnabled()) return null;
  return (
    <div className="relative my-6">
      <div className="absolute inset-0 flex items-center">
        <div className="w-full border-t border-[var(--auth-border)]" />
      </div>
      <div className="relative flex justify-center text-xs">
        <span className="bg-[var(--auth-card)] px-3 text-[var(--auth-muted)]">
          or continue with
        </span>
      </div>
    </div>
  );
}
`
}

// adminAtlasAuthShell returns the Atlas split-static layout — the default.
// Left panel is a calm blue gradient with brand identity + tagline; right
// panel hosts the form. Inter throughout. Inspired by the Acme / dashboard
// reference (split-screen with dashboard preview on left).
func adminAtlasAuthShell() string {
	return `"use client";

import Link from "next/link";
import type { ReactNode } from "react";
import type { ThemeTokens } from "@repo/shared/themes";
import { brand } from "@repo/shared/brand";
import { SocialAuthButtons, SocialAuthDivider } from "./SocialAuthButtons";
import type { AuthMode } from "./AuthShell";

interface Props {
  theme: ThemeTokens;
  mode: AuthMode;
  title: string;
  subtitle?: string;
  children: ReactNode;
  errorMessage?: string;
  showSocial?: boolean;
}

const switchLinks: Record<AuthMode, { hint: string; href: string; label: string }> = {
  "login":    { hint: "Don't have an account?", href: "/sign-up",         label: "Create one" },
  "sign-up":  { hint: "Already have an account?", href: "/login",         label: "Sign in" },
  "forgot":   { hint: "Remembered your password?", href: "/login",        label: "Back to sign in" },
  "reset":    { hint: "Remembered your password?", href: "/login",        label: "Back to sign in" },
};

export function AtlasAuthShell({ theme, mode, title, subtitle, children, errorMessage, showSocial = true }: Props) {
  const t = theme.colors;
  const f = theme.fonts;
  const sw = switchLinks[mode];

  return (
    <div
      className="flex min-h-screen"
      style={{
        fontFamily: f.ui,
        background: t.bg,
        color: t.fg,
        // CSS variables consumed by SocialAuthButtons + form inputs so
        // the same component fits every theme without per-shell styling.
        ["--auth-bg" as string]: t.bg,
        ["--auth-fg" as string]: t.fg,
        ["--auth-card" as string]: t.card,
        ["--auth-border" as string]: t.border,
        ["--auth-muted" as string]: t.muted,
        ["--auth-primary" as string]: t.primary,
        ["--auth-primary-fg" as string]: t.primaryFg,
        ["--auth-accent" as string]: t.accent,
        ["--auth-radius" as string]: theme.radius,
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
      <div
        className="flex flex-1 items-center justify-center px-6 py-12"
        style={{ background: t.bg }}
      >
        <div className="w-full max-w-md space-y-8">
          {/* Brand for sub-lg breakpoints */}
          <div className="lg:hidden flex items-center justify-center gap-2 text-2xl font-bold" style={{ color: t.primary }}>
            <BrandMark />
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

          {showSocial && (
            <>
              <SocialAuthDivider />
              <SocialAuthButtons />
            </>
          )}

          <p className="text-center text-sm" style={{ color: t.muted }}>
            {sw.hint}{" "}
            <Link href={sw.href} className="font-medium" style={{ color: t.primary }}>
              {sw.label}
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}

function BrandMark() {
  if (brand.logo.image) {
    return <img src={brand.logo.image} alt={brand.name} className="h-8 w-8" />;
  }
  return (
    <span
      className="inline-flex h-8 w-8 items-center justify-center rounded-md text-white font-bold"
      style={{ background: "rgba(255,255,255,0.15)" }}
    >
      {brand.logo.text}
    </span>
  );
}
`
}

// adminAuroraAuthShell returns the Aurora centered layout — single card on
// a pastel wallpaper, no hero panel. WorkOS / Clerk inspired. The "hero
// background" token doubles as the page wallpaper here.
func adminAuroraAuthShell() string {
	return `"use client";

import Link from "next/link";
import type { ReactNode } from "react";
import type { ThemeTokens } from "@repo/shared/themes";
import { brand } from "@repo/shared/brand";
import { SocialAuthButtons, SocialAuthDivider } from "./SocialAuthButtons";
import type { AuthMode } from "./AuthShell";

interface Props {
  theme: ThemeTokens;
  mode: AuthMode;
  title: string;
  subtitle?: string;
  children: ReactNode;
  errorMessage?: string;
  showSocial?: boolean;
}

const switchLinks: Record<AuthMode, { hint: string; href: string; label: string }> = {
  "login":    { hint: "Don't have an account?", href: "/sign-up", label: "Sign up" },
  "sign-up":  { hint: "Already have an account?", href: "/login", label: "Log in" },
  "forgot":   { hint: "Remembered your password?", href: "/login", label: "Back to login" },
  "reset":    { hint: "Remembered your password?", href: "/login", label: "Back to login" },
};

export function AuroraAuthShell({ theme, mode, title, subtitle, children, errorMessage, showSocial = true }: Props) {
  const t = theme.colors;
  const f = theme.fonts;
  const sw = switchLinks[mode];

  return (
    <div
      className="flex min-h-screen items-center justify-center px-4 py-12"
      style={{
        fontFamily: f.ui,
        // Soft pastel wallpaper made from the heroBg token + a radial
        // gradient pulse so the card lifts off the background.
        background: ` + "`" + `radial-gradient(60% 60% at 30% 20%, ${t.accent}22, transparent 60%), radial-gradient(50% 50% at 80% 80%, ${t.primary}1a, transparent 60%), ${t.heroBg}` + "`" + `,
        color: t.fg,
        ["--auth-bg" as string]: t.bg,
        ["--auth-fg" as string]: t.fg,
        ["--auth-card" as string]: t.card,
        ["--auth-border" as string]: t.border,
        ["--auth-muted" as string]: t.muted,
        ["--auth-primary" as string]: t.primary,
        ["--auth-primary-fg" as string]: t.primaryFg,
        ["--auth-accent" as string]: t.accent,
        ["--auth-radius" as string]: theme.radius,
      } as React.CSSProperties}
    >
      <div
        className="w-full max-w-md rounded-2xl border shadow-xl p-8 space-y-6"
        style={{ background: t.card, borderColor: t.border }}
      >
        <div className="flex flex-col items-center text-center space-y-3">
          <BrandMark color={t.primary} />
          <div>
            <h2 className="text-2xl font-semibold" style={{ fontFamily: f.display }}>{title}</h2>
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

        {showSocial && (
          <>
            <SocialAuthDivider />
            <SocialAuthButtons />
          </>
        )}

        <p className="text-center text-sm" style={{ color: t.muted }}>
          {sw.hint}{" "}
          <Link href={sw.href} className="font-medium" style={{ color: t.primary }}>
            {sw.label}
          </Link>
        </p>
      </div>
    </div>
  );
}

function BrandMark({ color }: { color: string }) {
  if (brand.logo.image) {
    return <img src={brand.logo.image} alt={brand.name} className="h-10 w-10" />;
  }
  return (
    <span
      className="inline-flex h-10 w-10 items-center justify-center rounded-xl text-white font-bold text-lg"
      style={{ background: color }}
    >
      {brand.logo.text}
    </span>
  );
}
`
}

// adminPulseAuthShell returns the Pulse split-carousel layout — auth form
// on the left, rotating hero photography on the right. Falls back to a
// static hero card when brand.hero.images is empty.
func adminPulseAuthShell() string {
	return `"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import type { ReactNode } from "react";
import type { ThemeTokens } from "@repo/shared/themes";
import { brand } from "@repo/shared/brand";
import { SocialAuthButtons, SocialAuthDivider } from "./SocialAuthButtons";
import type { AuthMode } from "./AuthShell";

interface Props {
  theme: ThemeTokens;
  mode: AuthMode;
  title: string;
  subtitle?: string;
  children: ReactNode;
  errorMessage?: string;
  showSocial?: boolean;
}

const switchLinks: Record<AuthMode, { hint: string; href: string; label: string }> = {
  "login":    { hint: "New here?", href: "/sign-up",          label: "Create an account" },
  "sign-up":  { hint: "Already signed up?", href: "/login",   label: "Log in" },
  "forgot":   { hint: "Got it back?", href: "/login",         label: "Back to sign in" },
  "reset":    { hint: "Got it back?", href: "/login",         label: "Back to sign in" },
};

export function PulseAuthShell({ theme, mode, title, subtitle, children, errorMessage, showSocial = true }: Props) {
  const t = theme.colors;
  const f = theme.fonts;
  const sw = switchLinks[mode];

  return (
    <div
      className="flex min-h-screen"
      style={{
        fontFamily: f.ui,
        background: t.bg,
        color: t.fg,
        ["--auth-bg" as string]: t.bg,
        ["--auth-fg" as string]: t.fg,
        ["--auth-card" as string]: t.card,
        ["--auth-border" as string]: t.border,
        ["--auth-muted" as string]: t.muted,
        ["--auth-primary" as string]: t.primary,
        ["--auth-primary-fg" as string]: t.primaryFg,
        ["--auth-accent" as string]: t.accent,
        ["--auth-radius" as string]: theme.radius,
      } as React.CSSProperties}
    >
      {/* Left form panel */}
      <div className="flex flex-1 items-center justify-center px-6 py-12 lg:px-12">
        <div className="w-full max-w-md space-y-8">
          <div className="flex items-center gap-2 text-xl font-bold" style={{ color: t.fg }}>
            <BrandMark accent={t.accent} />
            <span style={{ fontFamily: f.display }}>{brand.name}</span>
          </div>

          <div>
            <h2 className="text-3xl font-bold" style={{ fontFamily: f.display }}>{title}</h2>
            {subtitle && <p className="mt-2 text-sm" style={{ color: t.muted }}>{subtitle}</p>}
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

          {showSocial && (
            <>
              <SocialAuthDivider />
              <SocialAuthButtons />
            </>
          )}

          <p className="text-center text-sm" style={{ color: t.muted }}>
            {sw.hint}{" "}
            <Link href={sw.href} className="font-medium" style={{ color: t.fg, textDecoration: "underline" }}>
              {sw.label}
            </Link>
          </p>
        </div>
      </div>

      {/* Right hero carousel */}
      <PulseHeroCarousel accent={t.accent} fg={t.heroFg} bg={t.heroBg} fontDisplay={f.display} />
    </div>
  );
}

function PulseHeroCarousel({ accent, fg, bg, fontDisplay }: { accent: string; fg: string; bg: string; fontDisplay: string }) {
  const images = brand.hero.images.filter(Boolean);
  const [idx, setIdx] = useState(0);

  // Crossfade between images on a timer. Stops scheduling new ticks
  // when there's only one image so the timer doesn't fire pointlessly.
  useEffect(() => {
    if (images.length < 2) return;
    const id = window.setInterval(
      () => setIdx((i) => (i + 1) % images.length),
      brand.hero.intervalMs
    );
    return () => window.clearInterval(id);
  }, [images.length]);

  return (
    <div className="hidden lg:flex lg:w-1/2 relative overflow-hidden" style={{ background: bg }}>
      {images.length === 0 && (
        <div className="flex flex-1 items-center justify-center p-12">
          <div className="max-w-md space-y-4 text-center" style={{ color: fg }}>
            <div
              className="inline-flex h-16 w-16 items-center justify-center rounded-full"
              style={{ background: accent }}
            >
              <span className="text-3xl font-bold" style={{ color: "#0f0f0f" }}>{brand.logo.text}</span>
            </div>
            <h3 className="text-2xl font-bold whitespace-pre-line" style={{ fontFamily: fontDisplay }}>
              {brand.tagline}
            </h3>
            <p className="text-sm opacity-70">{brand.description}</p>
          </div>
        </div>
      )}

      {images.map((src, i) => (
        <img
          key={src + i}
          src={src}
          alt=""
          className="absolute inset-0 h-full w-full object-cover transition-opacity duration-1000"
          style={{ opacity: i === idx ? 1 : 0 }}
        />
      ))}

      {/* Caption overlay */}
      {images.length > 0 && (
        <div className="relative z-10 flex flex-1 items-end p-12">
          <div className="space-y-3" style={{ color: "#ffffff" }}>
            <div
              className="inline-block rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-wider"
              style={{ background: accent, color: "#0f0f0f" }}
            >
              {brand.name}
            </div>
            <h3 className="text-3xl font-bold max-w-md whitespace-pre-line" style={{ fontFamily: fontDisplay, textShadow: "0 2px 12px rgba(0,0,0,0.5)" }}>
              {brand.tagline}
            </h3>
          </div>
        </div>
      )}

      {/* Dots */}
      {images.length > 1 && (
        <div className="absolute bottom-6 right-6 z-10 flex gap-2">
          {images.map((_, i) => (
            <button
              key={i}
              type="button"
              aria-label={"Go to slide " + (i + 1)}
              onClick={() => setIdx(i)}
              className="h-1.5 rounded-full transition-all"
              style={{
                width: i === idx ? 24 : 8,
                background: i === idx ? accent : "rgba(255,255,255,0.5)",
              }}
            />
          ))}
        </div>
      )}
    </div>
  );
}

function BrandMark({ accent }: { accent: string }) {
  if (brand.logo.image) {
    return <img src={brand.logo.image} alt={brand.name} className="h-7 w-7" />;
  }
  return (
    <span
      className="inline-flex h-7 w-7 items-center justify-center rounded text-black font-bold"
      style={{ background: accent }}
    >
      {brand.logo.text}
    </span>
  );
}
`
}
