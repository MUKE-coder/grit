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
        background: `radial-gradient(60% 60% at 30% 20%, ${t.accent}22, transparent 60%), radial-gradient(50% 50% at 80% 80%, ${t.primary}1a, transparent 60%), ${t.heroBg}`,
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
