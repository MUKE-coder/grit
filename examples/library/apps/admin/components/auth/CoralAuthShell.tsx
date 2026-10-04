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

function BrandMark({ color, fg, size = 40 }: { color: string; fg: string; size?: number }) {
  if (brand.logo.image) {
    return <img src={brand.logo.image} alt={brand.name} style={{ height: size, width: size }} />;
  }
  return (
    <span
      className="inline-flex items-center justify-center rounded-xl font-bold"
      style={{ background: color, color: fg, height: size, width: size, fontSize: size * 0.45 }}
    >
      {brand.logo.text}
    </span>
  );
}

function ErrorBanner({ message }: { message?: string }) {
  if (!message) return null;
  return (
    <div
      role="alert"
      className="rounded-[var(--auth-radius)] border px-4 py-3 text-sm"
      style={{ borderColor: "#fecaca", background: "#fef2f2", color: "#b91c1c" }}
    >
      {message}
    </div>
  );
}

export function CoralAuthShell(props: Props) {
  const { theme, mode, title, subtitle, children, errorMessage, showSocial = true } = props;
  const t = theme.colors;
  const f = theme.fonts;
  const sw = switchLinks[mode];
  return (
    <div
      className="relative flex min-h-screen items-center justify-center overflow-hidden px-4 py-10"
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
      {/* A suggestion of the app behind the dialog, built from the theme's own
          tokens so it changes with the brand instead of shipping a photo. */}
      <div aria-hidden="true" className="pointer-events-none absolute inset-0 grid grid-cols-2 gap-5 p-8 opacity-70 blur-[3px] sm:grid-cols-4">
        {Array.from({ length: 8 }, (_, i) => (
          <div key={i} className="flex flex-col gap-3">
            <div
              className="aspect-[4/3] rounded-2xl"
              style={{ background: "linear-gradient(135deg, " + t.accent + "33, " + t.heroBg + ")" }}
            />
            <div className="h-3 w-3/4 rounded" style={{ background: t.border }} />
            <div className="h-3 w-1/2 rounded" style={{ background: t.border }} />
          </div>
        ))}
      </div>
      <div aria-hidden="true" className="absolute inset-0" style={{ background: "rgba(17,17,17,0.35)" }} />

      <div
        className="relative w-full max-w-[520px] rounded-[28px] px-6 py-10 shadow-2xl sm:px-12"
        style={{ background: t.card, color: t.fg }}
      >
        <div className="mb-7 flex flex-col items-center gap-4">
          <BrandMark color={t.primary} fg={t.primaryFg} />
          <div className="text-center">
            <h1 className="text-[26px] font-semibold leading-tight" style={{ fontFamily: f.display }}>{title}</h1>
            {subtitle && <p className="mt-1.5 text-sm" style={{ color: t.muted }}>{subtitle}</p>}
          </div>
        </div>

        <div className="space-y-5">
          <ErrorBanner message={errorMessage} />
          {children}
          {showSocial && (
            <>
              <SocialAuthDivider />
              <SocialAuthButtons />
            </>
          )}
          <p className="text-center text-sm" style={{ color: t.muted }}>
            {sw.hint}{" "}
            <Link href={sw.href} className="font-medium underline underline-offset-2" style={{ color: t.primary }}>
              {sw.label}
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
