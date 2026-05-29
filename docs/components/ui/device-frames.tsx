'use client'

// Device frame primitives (Browser / Desktop / Mobile) used by the
// "Three Platforms, One Framework" section. Adapted from a 21st.dev
// inspiration pass — fully restyled to match Grit's dark/light brand
// (border-border, bg-card, primary accents) and given consistent chrome
// so the three frames read as a related family. Each frame ships with
// its own traffic-light controls + URL bar / titlebar / status bar.

import React from 'react'
import { cn } from '@/lib/utils'

// ─── BrowserFrame ────────────────────────────────────────────────────────
// Web target. Mac-style window controls + spoofed URL bar with Grit URL.
export function BrowserFrame({
  url = 'localhost:8080',
  children,
  className,
  glow,
}: {
  url?: string
  children?: React.ReactNode
  className?: string
  glow?: boolean
}) {
  return (
    <div
      className={cn(
        'relative rounded-xl border-2 border-border/60 bg-card overflow-hidden',
        'shadow-[0_2px_0_rgba(0,0,0,0.04),0_24px_64px_-16px_rgba(0,0,0,0.4)] dark:shadow-[0_0_0_1px_rgba(255,255,255,0.04),0_24px_64px_-16px_rgba(0,0,0,0.6)]',
        glow && 'ring-1 ring-primary/15',
        className,
      )}
    >
      {/* Chrome */}
      <div className="flex items-center gap-2 px-3 py-2 bg-card border-b border-border/60">
        <div className="flex items-center gap-1.5">
          <span className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
          <span className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
          <span className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
        </div>
        <div className="flex-1 flex justify-center">
          <div className="inline-flex items-center gap-1.5 rounded-md bg-background border border-border/60 px-2.5 py-1 max-w-[260px] truncate">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              className="h-3 w-3 text-emerald-500 shrink-0"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
              <path d="M7 11V7a5 5 0 0 1 10 0v4" />
            </svg>
            <span className="text-[11px] font-mono text-muted-foreground truncate">{url}</span>
          </div>
        </div>
        <div className="w-12" />
      </div>
      {/* Body */}
      <div className="relative bg-background">{children}</div>
    </div>
  )
}

// ─── DesktopFrame ────────────────────────────────────────────────────────
// Wails / Tauri desktop target. Slim title bar, no URL — just a window
// title centred + frameless dot decoration.
export function DesktopFrame({
  title = 'Grit Desktop',
  children,
  className,
  glow,
}: {
  title?: string
  children?: React.ReactNode
  className?: string
  glow?: boolean
}) {
  return (
    <div
      className={cn(
        'relative rounded-xl border-2 border-border/60 bg-card overflow-hidden',
        'shadow-[0_2px_0_rgba(0,0,0,0.04),0_24px_64px_-16px_rgba(0,0,0,0.4)] dark:shadow-[0_0_0_1px_rgba(255,255,255,0.04),0_24px_64px_-16px_rgba(0,0,0,0.6)]',
        glow && 'ring-1 ring-primary/15',
        className,
      )}
    >
      {/* Title bar */}
      <div className="flex items-center gap-2 px-3 py-2 bg-card border-b border-border/60">
        <div className="flex items-center gap-1.5">
          <span className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
          <span className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
          <span className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
        </div>
        <div className="flex-1 text-center">
          <span className="text-[11px] font-medium text-muted-foreground">{title}</span>
        </div>
        <div className="flex items-center gap-1.5 text-muted-foreground/50">
          {/* Minimize / maximize / close hint glyphs */}
          <span className="text-[10px]">—</span>
          <span className="text-[10px]">▢</span>
          <span className="text-[10px]">✕</span>
        </div>
      </div>
      {/* Body */}
      <div className="relative bg-background">{children}</div>
    </div>
  )
}

// ─── MobileFrame ─────────────────────────────────────────────────────────
// Mobile (Expo) target. Phone shell with rounded corners + status bar +
// dynamic-island cutout + home indicator. Inspired by 21st.dev's
// IPhoneMockup but trimmed to the chrome we actually need.
export function MobileFrame({
  children,
  className,
  glow,
  time = '9:41',
}: {
  children?: React.ReactNode
  className?: string
  glow?: boolean
  time?: string
}) {
  return (
    <div
      className={cn(
        'relative rounded-[28px] bg-[#0b0b0d] p-[6px]',
        'shadow-[0_24px_64px_-16px_rgba(0,0,0,0.6),0_2px_4px_rgba(0,0,0,0.2)] dark:shadow-[0_24px_64px_-16px_rgba(0,0,0,0.8)]',
        glow && 'ring-1 ring-primary/30',
        className,
      )}
    >
      {/* Inner screen */}
      <div className="relative overflow-hidden rounded-[22px] bg-background">
        {/* Status bar */}
        <div className="relative flex items-center justify-between px-5 pt-2 pb-1 text-[10px] font-semibold text-foreground">
          <span>{time}</span>
          {/* Dynamic island cutout — absolutely positioned over the status bar */}
          <span className="absolute left-1/2 top-1.5 -translate-x-1/2 h-4 w-20 rounded-full bg-black" />
          <div className="flex items-center gap-1">
            <span className="text-[9px]">●●●●</span>
            <span className="text-[9px]">5G</span>
            <span className="inline-block h-2.5 w-4 rounded-sm border border-foreground/60 relative">
              <span className="absolute inset-[1.5px] right-[3px] bg-foreground/70 rounded-[1px]" />
              <span className="absolute -right-[2.5px] top-[3px] h-1 w-0.5 bg-foreground/60 rounded-r-sm" />
            </span>
          </div>
        </div>
        {/* App body */}
        <div className="relative pt-3 pb-6">{children}</div>
        {/* Home indicator */}
        <div className="absolute bottom-1.5 left-1/2 -translate-x-1/2 h-1 w-24 rounded-full bg-foreground/40" />
      </div>
    </div>
  )
}
