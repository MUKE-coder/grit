'use client'

import { useState } from 'react'
import type React from 'react'
import { LogIn, MonitorSmartphone, ShieldCheck, KeySquare, Building2, Check } from 'lucide-react'
import { CodeBlock } from '@/components/code-block'

/**
 * Authentication, as shipped.
 *
 * Written against the scaffold source rather than from memory, which is why
 * two things are stated as absent: there is no email-verification flow (the
 * User model has EmailVerifiedAt, but only an IdP ever sets it), and there is
 * no API-key auth. Claiming either would be the kind of thing a reader checks
 * in thirty seconds and never trusts you again after.
 *
 * The provider list is exactly the two goth providers the scaffold imports.
 * Adding a logo here without adding the provider package is a lie with a
 * design budget.
 */

interface Feature {
  key: string
  label: string
  icon: React.ComponentType<{ className?: string }>
  headline: string
  blurb: string
  points: string[]
  /** A tab carries screenshots or a code sample. */
  image?: string
  imageLabel?: string
  /** Several shots, switched by sub-tabs, for features that are a flow. */
  shots?: { src: string; label: string; chrome: string }[]
  code?: string
  language?: string
  codeFile?: string
}

const FEATURES: Feature[] = [
  {
    key: 'signin',
    label: 'Sign in',
    icon: LogIn,
    headline: 'The pages you would have built on day one',
    blurb:
      'Login, register, forgot-password and reset-password, themed with the rest of the admin and wired to a Go API that already hashes with bcrypt and issues a JWT pair.',
    points: [
      'Access token (15m) + refresh token (168h), both configurable',
      'HttpOnly cookies for the browser, Bearer headers for mobile and desktop',
      'CSRF enforced on cookie-authenticated mutations, transparent to Bearer clients',
      'Login rate-limited to 5 attempts / 15 min in production; register to 3',
    ],
    image: '/images/auth/login.png',
    imageLabel: 'localhost:3001/login',
  },
  {
    key: 'sessions',
    label: 'Sessions & revocation',
    icon: MonitorSmartphone,
    headline: 'Every refresh token is a row you can kill',
    blurb:
      'A JWT you cannot revoke is a liability. Every refresh token is backed by a session record storing only a SHA-256 hash, so signing a device out actually signs it out.',
    points: [
      'Refresh rotation with replay detection — a reused token kills the session',
      'Idle timeout of 7 days, absolute timeout of 30 days',
      'Per-device sign-out, plus revoke-all that spares the device you are on',
      'Changing your password signs out every other device automatically',
    ],
    image: '/images/auth/sessions.png',
    imageLabel: 'localhost:3001/profile',
  },
  {
    key: 'social',
    label: 'Social login',
    icon: KeySquare,
    headline: 'Google and GitHub, off by default',
    blurb:
      'OAuth via goth, finding-or-creating by email and linking the provider to an existing local account. The buttons stay hidden until you set the credentials — a provider button that no provider backs is worse than none.',
    points: [
      'Google and GitHub are the two providers the scaffold ships',
      'Registers a provider only when its client ID is present',
      'Marks the email verified, since the IdP already did that work',
      'SOCIAL_AUTH_ENABLED controls whether the buttons render at all',
    ],
    code: `# .env — the buttons appear once these are set
SOCIAL_AUTH_ENABLED=true

GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
GITHUB_CLIENT_ID=...
GITHUB_CLIENT_SECRET=...

# Where the API sends the browser back to
OAUTH_FRONTEND_URL=http://localhost:3001

# Callback URLs to register with each provider:
#   <APP_URL>/api/auth/oauth/google/callback
#   <APP_URL>/api/auth/oauth/github/callback`,
    language: 'bash',
    codeFile: '.env',
  },
  {
    key: 'totp',
    label: 'Two-factor',
    icon: ShieldCheck,
    headline: 'TOTP with backup codes and trusted devices',
    blurb:
      'Scan a QR, confirm with a live code, and save ten backup codes. From then on a sign-in from an untrusted device asks for a code — or a backup code, if the phone is gone. All of it in the generated admin, not a wiring exercise.',
    points: [
      'The QR is rendered by the API, so no client ships an encoder or touches the raw secret',
      'Ten backup codes, hashed at rest and consumed on use — shown exactly once',
      'Trust a device for 30 days; list them, revoke one, or revoke all',
      'Login returns a short-lived pending token instead of a session when 2FA is due',
      'TOTP_ISSUER controls the name shown in the authenticator app',
    ],
    shots: [
      { src: '/images/auth/2fa-setup.png', label: 'Enrol', chrome: 'localhost:3001/profile' },
      { src: '/images/auth/2fa-challenge.png', label: 'Challenge', chrome: 'localhost:3001/login' },
      { src: '/images/auth/2fa-on.png', label: 'Manage', chrome: 'localhost:3001/profile' },
    ],
  },
  {
    key: 'rbac',
    label: 'Roles & SSO',
    icon: Building2,
    headline: 'Permissions in the database, SSO per customer',
    blurb:
      'Roles are rows with a grants array, so they change without a deploy. Every generated resource registers its own permissions. For enterprise buyers, SSO connections are rows too — OIDC or SAML, one per customer.',
    points: [
      'ADMIN / EDITOR / USER seeded, then editable — ADMIN holds the "*" wildcard',
      'Permission keys are <feature>.<action>: products.create, users.view',
      'Middleware takes either a role or perm:<key>, so both styles work',
      'SSO connections live in the database: OIDC or SAML 2.0, added per customer from the admin',
    ],
    code: `// A resource generated today registers its own permissions.
// Nothing to add to a constants file.

middleware.RequireRole("ADMIN")              // by role
middleware.RequireRole("perm:products.edit") // by permission

// Seeded roles — non-destructive, so your edits survive a re-seed
ADMIN  → ["*"]
EDITOR → ["uploads.create", "uploads.view",
          "uploads.delete", "users.view"]
USER   → []`,
    language: 'go',
    codeFile: 'internal/middleware/auth.go',
  },
]

export function AuthShowcase() {
  const [active, setActive] = useState(FEATURES[0].key)
  const [shot, setShot] = useState(0)
  const feature = FEATURES.find((f) => f.key === active) ?? FEATURES[0]
  // Guard the index: switching tabs leaves it pointing past a shorter list.
  const shots = feature.shots ?? []
  const current = shots[Math.min(shot, shots.length - 1)]

  const pick = (key: string) => {
    setActive(key)
    setShot(0)
  }

  return (
    <div>
      <div role="tablist" aria-label="Authentication feature" className="flex flex-wrap gap-2 mb-8">
        {FEATURES.map((f) => {
          const selected = f.key === active
          return (
            <button
              key={f.key}
              type="button"
              role="tab"
              aria-selected={selected}
              onClick={() => pick(f.key)}
              className={`inline-flex items-center gap-2 rounded-full border px-4 py-2 text-sm font-medium transition-colors ${
                selected
                  ? 'border-primary/40 bg-primary/10 text-foreground'
                  : 'border-border/60 text-muted-foreground hover:text-foreground hover:border-border'
              }`}
            >
              <f.icon className="h-4 w-4" />
              {f.label}
            </button>
          )
        })}
      </div>

      <div className="grid lg:grid-cols-[1fr_20rem] gap-8 lg:gap-10 items-start">
        <div>
          {shots.length > 0 ? (
            <>
              <div role="tablist" aria-label="Two-factor step" className="flex gap-1.5 mb-3">
                {shots.map((sh, i) => {
                  const on = i === Math.min(shot, shots.length - 1)
                  return (
                    <button
                      key={sh.src}
                      type="button"
                      role="tab"
                      aria-selected={on}
                      onClick={() => setShot(i)}
                      className={`rounded-lg border px-3 py-1.5 text-[12px] font-medium transition-colors ${
                        on
                          ? 'border-primary/40 bg-primary/10 text-foreground'
                          : 'border-border/50 text-muted-foreground hover:text-foreground'
                      }`}
                    >
                      <span className="font-mono text-[10px] text-muted-foreground/60 mr-1.5">
                        {i + 1}
                      </span>
                      {sh.label}
                    </button>
                  )
                })}
              </div>

              <div className="rounded-xl overflow-hidden border border-border bg-card/40 shadow-[0_24px_64px_-16px_rgba(2,6,23,0.5)]">
                <div className="flex items-center gap-2 px-3.5 py-2.5 bg-card/70 border-b border-border/60">
                  <div className="flex items-center gap-1.5">
                    <span className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                    <span className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                    <span className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                  </div>
                  <span className="mx-auto text-[11px] font-mono text-muted-foreground">
                    {current.chrome}
                  </span>
                </div>
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  key={current.src}
                  src={current.src}
                  alt={`${feature.label}: ${current.label}`}
                  className="w-full h-auto block"
                  loading="lazy"
                />
              </div>
            </>
          ) : feature.image ? (
            <div className="rounded-xl overflow-hidden border border-border bg-card/40 shadow-[0_24px_64px_-16px_rgba(2,6,23,0.5)]">
              <div className="flex items-center gap-2 px-3.5 py-2.5 bg-card/70 border-b border-border/60">
                <div className="flex items-center gap-1.5">
                  <span className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                  <span className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                  <span className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                </div>
                <span className="mx-auto text-[11px] font-mono text-muted-foreground">
                  {feature.imageLabel}
                </span>
              </div>
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                key={feature.key}
                src={feature.image}
                alt={`${feature.label}: ${feature.headline}`}
                className="w-full h-auto block"
                loading="lazy"
              />
            </div>
          ) : (
            <div className="rounded-xl overflow-hidden border border-border bg-card/40">
              <div className="flex items-center gap-2 px-4 py-2.5 bg-card/60 border-b border-border/60">
                <span className="text-[11.5px] font-mono text-muted-foreground">{feature.codeFile}</span>
              </div>
              <CodeBlock
                key={feature.key}
                code={feature.code ?? ''}
                language={feature.language ?? 'go'}
                className="!border-0 !rounded-none !shadow-none !bg-transparent dark:!bg-transparent !m-0"
              />
            </div>
          )}
        </div>

        <div>
          <h3 className="text-lg font-semibold mb-3 leading-snug">{feature.headline}</h3>
          <p className="text-sm text-muted-foreground leading-relaxed mb-5">{feature.blurb}</p>
          <ul className="space-y-2.5">
            {feature.points.map((p) => (
              <li key={p} className="flex gap-2 text-[12.5px] text-foreground/80 leading-relaxed">
                <Check className="h-3.5 w-3.5 shrink-0 mt-0.5 text-primary" strokeWidth={2.5} />
                {p}
              </li>
            ))}
          </ul>
        </div>
      </div>

      {/* Stating the gaps costs one line and buys the rest of the page. */}
      <p className="text-[11.5px] text-muted-foreground/70 mt-8 pt-5 border-t border-border/40 leading-relaxed max-w-3xl">
        Not included, so you are not surprised later: there is no email-verification flow (the
        field exists and social sign-in sets it, but nothing sends a verification mail), and no
        API-key authentication &mdash; machine callers use the same JWT flow.
      </p>
    </div>
  )
}
