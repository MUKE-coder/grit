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

export function EmeraldAuthShell(props: Props) {
  const { theme, mode, title, subtitle, children, errorMessage, showSocial = true } = props;
  const t = theme.colors;
  const f = theme.fonts;
  const sw = switchLinks[mode];
  const quote = brand.proof.quote;
  const author = brand.proof.quoteAuthor;
  const initials = author.split(" ").map((part: string) => part.charAt(0)).join("").slice(0, 2);

  return (
    <div
      className="grid min-h-screen lg:grid-cols-[minmax(0,560px)_1fr]"
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
      <section className="flex flex-col px-8 py-8 sm:px-14">
        <Link href="/" className="inline-flex items-center gap-2.5 text-xl font-semibold">
          <BrandMark color={t.primary} fg={t.primaryFg} size={28} />
          {brand.name}
        </Link>

        <div className="my-auto flex w-full max-w-[400px] flex-col gap-7 py-12">
          <div>
            <h1 className="text-[30px] font-semibold leading-tight" style={{ fontFamily: f.display }}>{title}</h1>
            {subtitle && <p className="mt-1.5 text-sm" style={{ color: t.muted }}>{subtitle}</p>}
          </div>
          <div className="space-y-5">
            <ErrorBanner message={errorMessage} />
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
        </div>
      </section>

      <aside
        className="hidden items-center justify-center border-l p-16 lg:flex"
        style={{ borderColor: t.border, background: t.heroBg, color: t.heroFg }}
      >
        <figure className="flex max-w-lg flex-col gap-6">
          <blockquote className="text-3xl leading-snug">
            <span aria-hidden="true" className="mb-2 block text-6xl leading-none opacity-30">&ldquo;</span>
            {quote}
          </blockquote>
          <figcaption className="flex items-center gap-3">
            <span
              className="grid h-10 w-10 place-items-center rounded-full text-sm font-semibold"
              style={{ background: t.primary, color: t.primaryFg }}
            >
              {initials}
            </span>
            <span className="text-sm font-medium">{author}</span>
          </figcaption>
        </figure>
      </aside>
    </div>
  );
}
