'use client'

import { useState } from 'react'
import { Check } from 'lucide-react'

/** The Grit UI mark, inlined so the block stays self-contained. */
function GritMark({ className = 'size-8' }: { className?: string }) {
  return (
    <svg viewBox="0 0 32 32" fill="none" aria-hidden="true" className={className}>
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M8 0H24A8 8 0 0 1 32 8V24A8 8 0 0 1 24 32H8A8 8 0 0 1 0 24V8A8 8 0 0 1 8 0ZM16 9.4A6.6 6.6 0 1 0 21.4 19.9V17.4H17.2A1.7 1.7 0 0 1 17.2 14H23.1A1.7 1.7 0 0 1 24.8 15.7V20.6A1.7 1.7 0 0 1 24.4 21.7A10 10 0 1 1 22.6 8.2A1.7 1.7 0 0 1 20.4 10.8A6.6 6.6 0 0 0 16 9.4Z"
        fill="currentColor"
      />
    </svg>
  )
}

const NAV = ['Platform', 'Documentation', 'Pricing', 'Community', 'Enterprise']

/** Neutral, invented company names. Real third-party marks in a demo logo
 *  cloud read as customer endorsements that were never given. */
const LOGOS = ['NORTHWIND', 'ACME', 'GLOBEX', 'INITECH', 'UMBRELLA', 'SOYLENT']

export default function DarkWithEmailCapture({
  title = 'Radically better full-stack tooling',
  subtitle = 'Ship higher-quality software faster. Be the hero of your engineering team.',
  ctaLabel = 'Start for free',
  placeholder = 'Your work e-mail',
  footnote = 'Start building for free or',
  footnoteLink = 'book a demo',
}: {
  title?: string
  subtitle?: string
  ctaLabel?: string
  placeholder?: string
  footnote?: string
  footnoteLink?: string
}) {
  const [email, setEmail] = useState('')
  const [submitted, setSubmitted] = useState(false)

  return (
    // This block is dark in both themes by design — it is a dark hero, not a
    // light hero that inverts. Fixed neutrals rather than dark: variants.
    <div className="relative isolate overflow-hidden bg-gray-950">
      {/* Two opposing light shafts, the signature of this layout */}
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 -z-10 opacity-70"
        style={{
          background:
            'radial-gradient(60rem 40rem at -10% 40%, rgba(129,140,248,0.18), transparent 60%), radial-gradient(60rem 40rem at 110% 40%, rgba(129,140,248,0.18), transparent 60%)',
        }}
      />
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 -z-10"
        style={{
          backgroundImage:
            'radial-gradient(rgba(255,255,255,0.07) 1px, transparent 1px)',
          backgroundSize: '4px 4px',
          maskImage: 'radial-gradient(50rem 40rem at 50% 30%, black, transparent)',
          WebkitMaskImage: 'radial-gradient(50rem 40rem at 50% 30%, black, transparent)',
        }}
      />

      {/* Nav */}
      <header className="border-b border-white/5">
        <nav
          aria-label="Global"
          className="mx-auto flex max-w-7xl items-center gap-8 px-6 py-5"
        >
          <a href="#" className="flex items-center gap-2 text-white">
            <GritMark className="size-7" />
            <span className="text-base font-semibold tracking-tight">Grit</span>
          </a>
          <div className="hidden gap-7 lg:flex">
            {NAV.map((item) => (
              <a
                key={item}
                href="#"
                className="text-sm text-gray-400 transition-colors hover:text-white"
              >
                {item}
              </a>
            ))}
          </div>
          <div className="ml-auto flex items-center gap-4">
            <a
              href="#"
              className="hidden text-sm text-gray-400 transition-colors hover:text-white sm:block"
            >
              Sign in
            </a>
            <a
              href="#"
              className="rounded-lg bg-indigo-500 px-4 py-1.5 text-sm font-semibold text-white transition-colors hover:bg-indigo-400"
            >
              Sign up
            </a>
          </div>
        </nav>
      </header>

      {/* Hero */}
      <div className="mx-auto max-w-4xl px-6 py-28 text-center sm:py-36">
        <div className="mb-10 flex items-center justify-center gap-2.5 text-white">
          <GritMark className="size-7" />
          <span className="text-lg font-semibold tracking-tight">Grit Framework</span>
        </div>

        <h1 className="bg-gradient-to-b from-white to-gray-400 bg-clip-text text-5xl font-bold tracking-tight text-balance text-transparent sm:text-7xl">
          {title}
        </h1>
        <p className="mx-auto mt-7 max-w-xl text-lg text-gray-400">{subtitle}</p>

        <form
          onSubmit={(e) => {
            e.preventDefault()
            setSubmitted(true)
          }}
          className="mx-auto mt-10 flex max-w-xl flex-col gap-3 sm:flex-row"
        >
          <label htmlFor="hero-email" className="sr-only">
            Work e-mail
          </label>
          <input
            id="hero-email"
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder={placeholder}
            className="min-w-0 flex-1 rounded-xl border border-white/10 bg-white/5 px-5 py-3.5 text-base text-white placeholder:text-gray-500 focus:border-indigo-400/60 focus:outline-none"
          />
          <button
            type="submit"
            className="inline-flex items-center justify-center gap-2 rounded-xl bg-indigo-500 px-6 py-3.5 text-base font-semibold text-white transition-colors hover:bg-indigo-400"
          >
            {submitted ? (
              <>
                <Check aria-hidden="true" className="size-4" />
                Check your inbox
              </>
            ) : (
              ctaLabel
            )}
          </button>
        </form>

        <p className="mt-5 text-sm text-gray-500">
          {footnote}{' '}
          <a href="#" className="text-gray-300 underline underline-offset-4">
            {footnoteLink}
          </a>
        </p>
      </div>

      {/* Logo cloud */}
      <div className="border-t border-white/5">
        <div className="mx-auto flex max-w-7xl flex-wrap items-center justify-center gap-x-12 gap-y-6 px-6 py-10">
          {LOGOS.map((logo) => (
            <span
              key={logo}
              className="font-mono text-sm tracking-[0.15em] text-gray-600"
            >
              {logo}
            </span>
          ))}
        </div>
      </div>
    </div>
  )
}
