"use client";

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
