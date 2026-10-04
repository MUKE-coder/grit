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

export function MonoAuthShell(props: Props) {
  const { theme, mode, title, subtitle, children, errorMessage, showSocial = true } = props;
  const t = theme.colors;
  const f = theme.fonts;
  const sw = switchLinks[mode];
  const points = brand.proof.points;

  return (
    <div
      className="grid min-h-screen lg:grid-cols-[1.25fr_1fr]"
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
      <section className="relative flex flex-col items-center overflow-hidden px-6 py-8">
        <div
          aria-hidden="true"
          className="pointer-events-none absolute inset-0"
          style={{
            backgroundImage:
              "linear-gradient(" + t.border + " 1px, transparent 1px), linear-gradient(90deg, " + t.border + " 1px, transparent 1px)",
            backgroundSize: "72px 72px",
            maskImage: "radial-gradient(ellipse at top, black 20%, transparent 70%)",
            WebkitMaskImage: "radial-gradient(ellipse at top, black 20%, transparent 70%)",
          }}
        />
        <Link href="/" className="relative text-2xl font-bold tracking-tight" style={{ fontFamily: f.display }}>
          {brand.name}
        </Link>

        <div className="relative my-auto flex w-full max-w-[420px] flex-col gap-7 py-12">
          <div className="text-center">
            <h1 className="text-2xl font-semibold leading-tight" style={{ fontFamily: f.display }}>{title}</h1>
            {subtitle && <p className="mt-1.5 text-sm" style={{ color: t.muted }}>{subtitle}</p>}
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
      </section>

      <aside
        className="hidden flex-col justify-center gap-10 border-l p-12 lg:flex"
        style={{ borderColor: t.border, background: t.heroBg }}
      >
        <div
          className="flex flex-col gap-6 rounded-2xl p-8 shadow-xl"
          style={{ background: t.fg, color: t.bg }}
        >
          <p className="text-2xl font-semibold leading-snug">{brand.proof.headline}</p>
          <p className="text-sm opacity-80">{brand.proof.subheadline}</p>
        </div>
        <ul className="grid gap-3">
          {points.map((point: string) => (
            <li key={point} className="flex items-center gap-3 text-sm" style={{ color: t.muted }}>
              <span aria-hidden="true" style={{ color: t.fg }}>&#10003;</span>
              {point}
            </li>
          ))}
        </ul>
      </aside>
    </div>
  );
}
