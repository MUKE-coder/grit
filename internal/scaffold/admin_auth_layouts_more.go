package scaffold

import "strings"

// Five more sign-in layouts, one per theme.
//
//	coral   → modal     a card over a blurred glimpse of the app
//	amber   → boxed     a plain bordered box under a wordmark
//	sky     → banner    a top bar and one bold heading, social first
//	mono    → showcase  a fine grid, with a panel of proof beside the form
//	emerald → quote     a narrow form column, a customer quote filling the rest
//
// They exist because "which auth screen" is the first thing anybody changes in
// a generated app and the last thing they want to build. The form inside is
// the same component in all eight layouts, so none of this touches validation,
// the second factor, or the sign-in call: a layout decides where things sit.
//
// Each file is generated from one preamble and one style block, below, rather
// than five copies of them. Five copies of a token list is five chances for one
// theme to quietly stop publishing --auth-radius.

// authShellPreamble is the imports, props and switch-links every themed shell
// shares. component is the exported function name.
func authShellPreamble(component string) string {
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

export function ` + component + `(props: Props) {`
}

// authShellVars is the inline style object that publishes the active theme's
// tokens as --auth-* custom properties.
//
// Every form control in the scaffold is styled against these, so a layout that
// forgets one renders a control with no border rather than an obvious error.
const authShellVars = `        ["--auth-bg" as string]: t.bg,
        ["--auth-fg" as string]: t.fg,
        ["--auth-card" as string]: t.card,
        ["--auth-border" as string]: t.border,
        ["--auth-muted" as string]: t.muted,
        ["--auth-primary" as string]: t.primary,
        ["--auth-primary-fg" as string]: t.primaryFg,
        ["--auth-accent" as string]: t.accent,
        ["--auth-radius" as string]: theme.radius,`

// authShellCommon is the destructuring every layout body starts with.
const authShellCommon = `  const { theme, mode, title, subtitle, children, errorMessage, showSocial = true } = props;
  const t = theme.colors;
  const f = theme.fonts;
  const sw = switchLinks[mode];
`

// authShell assembles one shell file: the shared preamble, the destructuring
// every body starts with, then the layout's own JSX.
func authShell(component, body string) string {
	return authShellPreamble(component) + "\n" + authShellCommon + body
}

// adminCoralAuthShell: the modal layout.
//
// The blurred cards behind it are deliberately not a photograph. A stock hero
// image in a scaffold is the first thing anybody deletes, and until they do it
// tells every visitor the product is a template.
func adminCoralAuthShell() string {
	body := `  return (
    <div
      className="relative flex min-h-screen items-center justify-center overflow-hidden px-4 py-10"
      style={{
        fontFamily: f.ui,
        background: t.bg,
        color: t.fg,
` + authShellVars + `
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
`
	return authShell("CoralAuthShell", body)
}

// adminAmberAuthShell: the boxed layout.
func adminAmberAuthShell() string {
	body := `  return (
    <div
      className="flex min-h-screen flex-col items-center"
      style={{
        fontFamily: f.ui,
        background: t.bg,
        color: t.fg,
` + authShellVars + `
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
`
	return authShell("AmberAuthShell", body)
}

// adminSkyAuthShell: the banner layout, social sign-in first.
func adminSkyAuthShell() string {
	body := `  return (
    <div
      className="min-h-screen"
      style={{
        fontFamily: f.ui,
        background: t.bg,
        color: t.fg,
` + authShellVars + `
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
`
	return authShell("SkyAuthShell", body)
}

// adminMonoAuthShell: the showcase layout.
func adminMonoAuthShell() string {
	body := `  const points = brand.proof.points;

  return (
    <div
      className="grid min-h-screen lg:grid-cols-[1.25fr_1fr]"
      style={{
        fontFamily: f.ui,
        background: t.bg,
        color: t.fg,
` + authShellVars + `
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
`
	return authShell("MonoAuthShell", body)
}

// adminEmeraldAuthShell: the quote layout.
func adminEmeraldAuthShell() string {
	body := `  const quote = brand.proof.quote;
  const author = brand.proof.quoteAuthor;
  const initials = author.split(" ").map((part: string) => part.charAt(0)).join("").slice(0, 2);

  return (
    <div
      className="grid min-h-screen lg:grid-cols-[minmax(0,560px)_1fr]"
      style={{
        fontFamily: f.ui,
        background: t.bg,
        color: t.fg,
` + authShellVars + `
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
`
	return authShell("EmeraldAuthShell", body)
}

// authShellFiles maps the five new shells to their file names, so the writers
// and the tests agree on one list.
func authShellFiles() map[string]string {
	return map[string]string{
		"CoralAuthShell.tsx":   adminCoralAuthShell(),
		"AmberAuthShell.tsx":   adminAmberAuthShell(),
		"SkyAuthShell.tsx":     adminSkyAuthShell(),
		"MonoAuthShell.tsx":    adminMonoAuthShell(),
		"EmeraldAuthShell.tsx": adminEmeraldAuthShell(),
	}
}

// authLayoutForShell says which layout each shell serves, for the dispatcher
// test.
func authLayoutForShell() map[string]string {
	return map[string]string{
		"CoralAuthShell":   "modal",
		"AmberAuthShell":   "boxed",
		"SkyAuthShell":     "banner",
		"MonoAuthShell":    "showcase",
		"EmeraldAuthShell": "quote",
	}
}

// shellName strips the .tsx from a file name in authShellFiles.
func shellName(file string) string {
	return strings.TrimSuffix(file, ".tsx")
}
