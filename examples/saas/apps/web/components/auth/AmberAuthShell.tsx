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

export function AmberAuthShell(props: Props) {
  const { theme, mode, title, subtitle, children, errorMessage, showSocial = true } = props;
  const t = theme.colors;
  const f = theme.fonts;
  const sw = switchLinks[mode];
  return (
    <div
      className="flex min-h-screen flex-col items-center"
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
      <div className="flex w-full flex-col items-center px-4 pt-8">
        <Link href="/" className="mb-5 inline-flex items-center gap-2.5">
          <BrandMark color={t.primary} fg={t.primaryFg} size={32} />
          <span className="text-xl font-bold tracking-tight" style={{ fontFamily: f.display }}>{brand.name}</span>
        </Link>

        <div
          className="w-full max-w-[400px] rounded-[var(--auth-radius)] border px-6 py-7 sm:px-8"
          style={{ borderColor: t.border, background: t.card }}
        >
          <h1 className="mb-5 text-[26px] font-medium leading-tight" style={{ fontFamily: f.display }}>{title}</h1>
          {subtitle && <p className="-mt-3 mb-5 text-sm" style={{ color: t.muted }}>{subtitle}</p>}

          <div className="space-y-5">
            <ErrorBanner message={errorMessage} />
            {children}
            {showSocial && (
              <>
                <SocialAuthDivider />
                <SocialAuthButtons />
              </>
            )}
          </div>

          <p className="mt-5 text-xs leading-5" style={{ color: t.muted }}>
            By continuing you agree to the{" "}
            <Link href="/terms" className="underline underline-offset-2">Terms of Service</Link> and{" "}
            <Link href="/privacy" className="underline underline-offset-2">Privacy Policy</Link>.
          </p>
        </div>

        <p className="mt-5 text-sm" style={{ color: t.muted }}>
          {sw.hint}{" "}
          <Link href={sw.href} className="font-medium underline underline-offset-2" style={{ color: t.accent }}>
            {sw.label}
          </Link>
        </p>
      </div>

      <footer
        className="mt-auto flex w-full flex-col items-center gap-2 border-t py-6 text-xs"
        style={{ borderColor: t.border, color: t.muted }}
      >
        <nav className="flex gap-4">
          <Link href="/terms" className="hover:underline">Conditions of use</Link>
          <Link href="/privacy" className="hover:underline">Privacy notice</Link>
        </nav>
        <p>&copy; {new Date().getFullYear()} {brand.name}</p>
      </footer>
    </div>
  );
}
