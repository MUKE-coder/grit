"use client";

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
          <BrandMark color={t.primary} fg={t.primaryFg} />
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

        <p className="text-sm opacity-60">Built with Grit: Go + React framework</p>
      </div>

      {/* Right form panel. main, not div: this is the page's content, and
          landmarks are how somebody skips the hero panel to reach it. */}
      <main
        className="flex flex-1 items-center justify-center px-6 py-12"
        style={{ background: t.bg }}
      >
        <div className="w-full max-w-md space-y-8">
          {/* Brand for sub-lg breakpoints */}
          <div className="lg:hidden flex items-center justify-center gap-2 text-2xl font-bold" style={{ color: t.primary }}>
            <BrandMark color={t.primary} fg={t.primaryFg} />
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
      </main>
    </div>
  );
}

// The brand colour, not a wash of the hero's own background: six of the eight
// themes have a light hero panel, and white on fifteen percent white is nothing
// at all. The letter was invisible on every one of them.
function BrandMark({ color, fg }: { color: string; fg: string }) {
  if (brand.logo.image) {
    return <img src={brand.logo.image} alt={brand.name} className="h-8 w-8" />;
  }
  return (
    <span
      className="inline-flex h-8 w-8 items-center justify-center rounded-md font-bold"
      style={{ background: color, color: fg }}
    >
      {brand.logo.text}
    </span>
  );
}
