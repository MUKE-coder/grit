'use client'

import { useState } from 'react'
import { Check, Terminal } from 'lucide-react'

/**
 * The themes, as actual screenshots of an actual admin.
 *
 * These are not mockups. Each one is the same generated dashboard — same
 * resources, same seeded data — rendered under a different theme and captured at
 * 1440px. That matters more than it sounds: a mockup proves a designer can draw,
 * a screenshot proves the framework produces it.
 *
 * Retake them the same way if the admin changes: run a generated project, flip
 * data-theme on <html>, screenshot at 1440×860. Anything else and the four stop
 * being comparable, which is the whole point of showing them together.
 */

interface Theme {
  key: string
  name: string
  tagline: string
  accent: string
  image: string
  /** False for midnight: `grit new --theme midnight` fails ValidateTheme. */
  scaffoldFlag: boolean
  dark?: boolean
}

const THEMES: Theme[] = [
  {
    key: 'atlas',
    name: 'Atlas',
    tagline: 'The default. Professional blue on white — teams, dashboards, internal tools.',
    accent: '#2563eb',
    image: '/images/themes/atlas.png',
    scaffoldFlag: true,
  },
  {
    key: 'aurora',
    name: 'Aurora',
    tagline: 'Apple-inspired monochrome. Near-black type and CTAs, colour reserved for meaning.',
    accent: '#1d1d1f',
    image: '/images/themes/aurora.png',
    scaffoldFlag: true,
  },
  {
    key: 'pulse',
    name: 'Pulse',
    tagline: 'Deeper blue with white elevated cards. Reads as infrastructure rather than SaaS.',
    accent: '#0051c3',
    image: '/images/themes/pulse.png',
    scaffoldFlag: true,
  },
  {
    key: 'midnight',
    name: 'Midnight',
    tagline:
      'The legacy dark palette, violet accent. A runtime variant rather than a full theme — it repaints the dashboard, but reuses the Atlas fonts and auth pages.',
    accent: '#6c5ce7',
    image: '/images/themes/midnight.png',
    scaffoldFlag: false,
    dark: true,
  },
]

export function ThemeShowcase() {
  const [active, setActive] = useState(THEMES[0].key)
  const theme = THEMES.find((t) => t.key === active) ?? THEMES[0]

  return (
    <div>
      <div className="flex flex-wrap items-center gap-2 mb-6">
        {THEMES.map((t) => {
          const selected = t.key === active
          return (
            <button
              key={t.key}
              type="button"
              onClick={() => setActive(t.key)}
              aria-pressed={selected}
              className={`inline-flex items-center gap-2 rounded-full border px-4 py-2 text-sm font-medium transition-colors ${
                selected
                  ? 'border-primary/40 bg-primary/10 text-foreground'
                  : 'border-border/60 text-muted-foreground hover:text-foreground hover:border-border'
              }`}
            >
              <span
                aria-hidden
                className="h-3 w-3 rounded-full ring-1 ring-black/10"
                style={{ backgroundColor: t.accent }}
              />
              {t.name}
              {selected && <Check className="h-3.5 w-3.5 text-primary" />}
            </button>
          )
        })}
      </div>

      <div className="grid lg:grid-cols-[1fr_18rem] gap-6 items-start">
        {/* Browser chrome around the shot, so it reads as a running app rather
            than a marketing image. */}
        <div className="rounded-xl overflow-hidden border border-border bg-card/40 shadow-[0_24px_64px_-16px_rgba(2,6,23,0.5)]">
          <div className="flex items-center gap-2 px-3.5 py-2.5 bg-card/70 border-b border-border/60">
            <div className="flex items-center gap-1.5">
              <span className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
              <span className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
              <span className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
            </div>
            <span className="mx-auto text-[11px] font-mono text-muted-foreground">
              localhost:3001/dashboard
            </span>
          </div>
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            key={theme.key}
            src={theme.image}
            alt={`The Grit admin panel in the ${theme.name} theme`}
            width={1440}
            height={860}
            className="w-full h-auto block"
            loading="lazy"
          />
        </div>

        <div>
          <div className="flex items-center gap-2 mb-2">
            <h3 className="text-lg font-semibold">{theme.name}</h3>
            {!theme.scaffoldFlag && (
              <span className="rounded-full border border-border/60 px-2 py-0.5 text-[10px] font-mono uppercase tracking-wider text-muted-foreground/70">
                runtime only
              </span>
            )}
          </div>
          <p className="text-sm text-muted-foreground leading-relaxed mb-5">{theme.tagline}</p>

          {theme.scaffoldFlag && (
            <div className="rounded-xl border border-border/60 bg-card/40 p-4 mb-4">
              <div className="flex items-center gap-1.5 mb-2 text-[11px] font-mono uppercase tracking-wider text-muted-foreground/70">
                <Terminal className="h-3 w-3" />
                At scaffold time
              </div>
              <code className="block text-[12px] font-mono text-foreground/90 break-all">
                grit new myapp --theme {theme.key}
              </code>
            </div>
          )}

          <div className="rounded-xl border border-border/60 bg-card/40 p-4">
            <div className="text-[11px] font-mono uppercase tracking-wider text-muted-foreground/70 mb-2">
              {theme.scaffoldFlag ? 'Or any time after' : 'In .env, any time'}
            </div>
            <code className="block text-[12px] font-mono text-foreground/90">
              THEME={theme.key}
            </code>
            <p className="text-[12px] text-muted-foreground mt-2 leading-relaxed">
              {theme.scaffoldFlag ? (
                <>
                  One line in <code className="text-foreground/70">.env</code>. Themes are CSS
                  variables, so switching one costs a restart, not a rewrite.
                </>
              ) : (
                <>
                  Midnight has no <code className="text-foreground/70">--theme</code> flag &mdash;
                  it is a runtime variant, so <code className="text-foreground/70">.env</code> is
                  the only way in.
                </>
              )}
            </p>
          </div>
        </div>
      </div>

      <p className="text-xs text-muted-foreground/70 mt-5">
        Real screenshots of a generated admin — same resources, same seeded data, same
        viewport. Only the theme changes.
      </p>
    </div>
  )
}
