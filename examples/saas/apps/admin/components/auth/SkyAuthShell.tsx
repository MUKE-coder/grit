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

export function SkyAuthShell(props: Props) {
  const { theme, mode, title, subtitle, children, errorMessage, showSocial = true } = props;
  const t = theme.colors;
  const f = theme.fonts;
  const sw = switchLinks[mode];
  return (
    <div
      className="min-h-screen"
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
      <header className="flex items-center justify-between px-6 py-5">
        <Link href="/" className="inline-flex items-center gap-2.5">
          <BrandMark color={t.primary} fg={t.primaryFg} size={32} />
          <span className="text-lg font-semibold tracking-tight">{brand.name}</span>
        </Link>
        <Link
          href="/"
          className="rounded-[var(--auth-radius)] border px-4 py-2 text-sm font-medium"
          style={{ borderColor: t.border, color: t.fg }}
        >
          Back to site
        </Link>
      </header>

      <main className="mx-auto flex w-full max-w-[560px] flex-col gap-8 px-6 pt-8 pb-16">
        <div>
          <h1 className="text-4xl font-bold leading-tight tracking-tight" style={{ fontFamily: f.display }}>{title}</h1>
          {subtitle && <p className="mt-2 text-base" style={{ color: t.muted }}>{subtitle}</p>}
        </div>

        <div className="space-y-5">
          <ErrorBanner message={errorMessage} />
          {/* Social first: on this layout most people arrive holding an
              identity they already have, and burying it under a password
              field asks them to make a new one. */}
          {showSocial && (
            <>
              <SocialAuthButtons />
              <SocialAuthDivider />
            </>
          )}
          {children}
          <p className="text-sm" style={{ color: t.muted }}>
            {sw.hint}{" "}
            <Link href={sw.href} className="font-medium underline underline-offset-2" style={{ color: t.primary }}>
              {sw.label}
            </Link>
          </p>
        </div>
      </main>
    </div>
  );
}
