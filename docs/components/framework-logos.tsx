// Inline SVG framework / runtime / service logos for the hero pill row.
// Keep these as pure JSX (no external icons) so the page renders without
// network round-trips and stays SSR-safe. Swap with custom files dropped
// into /public/images/icons/ later if branding requires it.

import React from 'react'

interface LogoProps { className?: string }
const SIZE = 'h-5 w-5'

export const GoLogo = ({ className = SIZE }: LogoProps) => (
  <svg viewBox="0 0 64 64" className={className} aria-label="Go">
    <g fill="#00ADD8">
      <path d="M4 26h12a1 1 0 010 2H4a1 1 0 010-2zM18 32h22a1 1 0 010 2H18a1 1 0 010-2zM7 38h12a1 1 0 010 2H7a1 1 0 010-2z" />
      <path d="M38 24c-7 0-12 4-12 12 0 6 4 12 12 12s12-6 12-12c0-7-5-12-12-12zm-3 16a2.5 2.5 0 110-5 2.5 2.5 0 010 5zm6 0a2.5 2.5 0 110-5 2.5 2.5 0 010 5z" />
    </g>
  </svg>
)

export const ReactLogo = ({ className = SIZE }: LogoProps) => (
  <svg viewBox="0 0 24 24" className={className} aria-label="React">
    <circle cx="12" cy="12" r="2.1" fill="#61DAFB" />
    <g fill="none" stroke="#61DAFB" strokeWidth="1">
      <ellipse cx="12" cy="12" rx="10" ry="4.5" />
      <ellipse cx="12" cy="12" rx="10" ry="4.5" transform="rotate(60 12 12)" />
      <ellipse cx="12" cy="12" rx="10" ry="4.5" transform="rotate(-60 12 12)" />
    </g>
  </svg>
)

export const VueLogo = ({ className = SIZE }: LogoProps) => (
  <svg viewBox="0 0 24 24" className={className} aria-label="Vue">
    <path d="M12 2 2 22h6l4-7 4 7h6L12 2Z" fill="#42B883" />
    <path d="M12 6l6 10h-3l-3-5-3 5H6L12 6Z" fill="#35495E" />
  </svg>
)

export const SvelteLogo = ({ className = SIZE }: LogoProps) => (
  <svg viewBox="0 0 24 24" className={className} aria-label="Svelte">
    <path
      fill="#FF3E00"
      d="M19.6 5.2C17.6 2.6 14.1 2.1 11.4 4l-4.5 2.9c-1.3.8-2.2 2.1-2.5 3.6-.3 1.2-.1 2.4.4 3.5-.4.5-.6 1.2-.8 1.8-.3 1.4-.1 2.8.6 4 1.9 2.6 5.4 3.2 8.1 1.3l4.5-2.9c1.3-.8 2.2-2.1 2.5-3.6.3-1.2.1-2.4-.4-3.5.4-.5.6-1.2.8-1.8.3-1.4.1-2.8-.6-4Z"
    />
  </svg>
)

export const NextLogo = ({ className = SIZE }: LogoProps) => (
  <svg viewBox="0 0 24 24" className={className} aria-label="Next.js">
    <circle cx="12" cy="12" r="11" fill="black" />
    <path d="M9.5 7.2v9.6M9.5 7.2l6.2 9.6M15.7 7.2v6" stroke="white" strokeWidth="1.4" fill="none" strokeLinecap="round" />
  </svg>
)

export const TanStackLogo = ({ className = SIZE }: LogoProps) => (
  <svg viewBox="0 0 24 24" className={className} aria-label="TanStack">
    <circle cx="12" cy="12" r="11" fill="#FF4154" />
    <path d="M8 8h8v2H8zm0 3h8v2H8zm0 3h5v2H8z" fill="white" />
  </svg>
)

export const TypeScriptLogo = ({ className = SIZE }: LogoProps) => (
  <svg viewBox="0 0 24 24" className={className} aria-label="TypeScript">
    <rect x="1" y="1" width="22" height="22" rx="2" fill="#3178C6" />
    <path d="M13 11h6v1.6h-2.2v6.1h-1.7v-6.1H13V11zm-7 0v1.6h2.5V19h1.7v-6.4h2.5V11H6z" fill="white" />
  </svg>
)

export const TailwindLogo = ({ className = SIZE }: LogoProps) => (
  <svg viewBox="0 0 24 24" className={className} aria-label="Tailwind CSS">
    <path
      fill="#38BDF8"
      d="M12 6c-2.5 0-4 1.2-4.5 3.5 1-1.2 2-1.7 3.5-1.4.9.2 1.5.9 2.2 1.7 1.2 1.2 2.5 2.7 5.3 2.7 2.5 0 4-1.2 4.5-3.5-1 1.2-2 1.7-3.5 1.4-.9-.2-1.5-.9-2.2-1.7C16.1 7.5 14.8 6 12 6zm-4.5 6c-2.5 0-4 1.2-4.5 3.5 1-1.2 2-1.7 3.5-1.4.9.2 1.5.9 2.2 1.7 1.2 1.2 2.5 2.7 5.3 2.7 2.5 0 4-1.2 4.5-3.5-1 1.2-2 1.7-3.5 1.4-.9-.2-1.5-.9-2.2-1.7-1.2-1.2-2.5-2.7-5.3-2.7z"
    />
  </svg>
)

export const PostgresLogo = ({ className = SIZE }: LogoProps) => (
  <svg viewBox="0 0 24 24" className={className} aria-label="PostgreSQL">
    <ellipse cx="12" cy="5" rx="9" ry="3" fill="#336791" />
    <path d="M3 5v14c0 1.7 4 3 9 3s9-1.3 9-3V5" fill="#336791" />
    <ellipse cx="12" cy="5" rx="9" ry="3" fill="none" stroke="white" strokeWidth="0.6" />
    <ellipse cx="12" cy="12" rx="9" ry="3" fill="none" stroke="white" strokeWidth="0.6" opacity="0.5" />
  </svg>
)

export const RedisLogo = ({ className = SIZE }: LogoProps) => (
  <svg viewBox="0 0 24 24" className={className} aria-label="Redis">
    <rect x="2" y="6" width="20" height="12" rx="2" fill="#DC382D" />
    <circle cx="7" cy="12" r="1.5" fill="white" />
    <circle cx="12" cy="12" r="1.5" fill="white" />
    <circle cx="17" cy="12" r="1.5" fill="white" />
  </svg>
)

export const DockerLogo = ({ className = SIZE }: LogoProps) => (
  <svg viewBox="0 0 24 24" className={className} aria-label="Docker">
    <path
      fill="#2496ED"
      d="M21.8 11.4c-.2-.2-1-.7-2.5-.7-.4 0-.8.1-1.2.2-.3-1.6-1.3-2.4-1.4-2.4l-.5-.3-.3.5c-.4.7-.7 1.5-.8 2.4 0 .6.1 1.1.3 1.5-.7.4-1.7.5-2 .5H2.4c-.2 0-.4.2-.4.4-.1 1.5.2 4.4 2.2 6.6 1.5 1.7 3.7 2.6 6.7 2.6 6.4 0 11.2-3 13.4-8.4.9 0 2.8 0 3.7-1.8.1-.1.2-.4-.2-.6z M3 10h2v2H3zM6 10h2v2H6zM9 10h2v2H9zM12 10h2v2h-2zM6 7h2v2H6zM9 7h2v2H9zM12 7h2v2h-2zM12 4h2v2h-2z"
    />
  </svg>
)
