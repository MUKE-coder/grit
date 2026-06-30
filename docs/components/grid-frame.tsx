// Blueprint frame (Tailwind-Plus / portfolio style). A full-page backdrop that
// frames a centred content column and treats the space around it as part of the
// grid — so every page reads as one symmetric system:
//   • a centred content column (max-w) carrying a faint grid texture;
//   • two rails on the column edges (the borders between column and gutters);
//   • diagonal-striped GUTTERS outside the rails so big screens are framed, not
//     empty;
//   • crosshair "+" markers on the rails that gently pulse;
//   • a slow light that drifts down each rail — subtle, not exaggerated.
//
// It is `absolute inset-0` inside the page's `relative isolate` wrapper, so it
// spans the whole document AND centres in the SAME context as flow content —
// the rails hug the section edges at every screen size. Sits at -z-10: above the
// wrapper background, behind all content.
//
// `width` must match the page's content container so the rails hug it (default
// matches Tailwind's max-w-6xl = 72rem).
export function GridFrame({ width = '72rem' }: { width?: string }) {
  return (
    <div aria-hidden className="pointer-events-none absolute inset-0 -z-10 overflow-hidden">
      {/* Soft primary glow at the top — the only colour moment. */}
      <div
        className="absolute inset-x-0 top-0 h-[640px]"
        style={{
          background:
            'radial-gradient(ellipse 60% 100% at 50% -10%, hsl(var(--primary) / 0.12), transparent 65%)',
        }}
      />

      {/* gutter │ content column │ gutter — center column = min(100%, width). */}
      <div
        className="absolute inset-0 grid"
        style={{ gridTemplateColumns: `1fr min(100%, ${width}) 1fr` }}
      >
        {/* Left gutter (striped) — its right border is the left rail. */}
        <div className="bg-grit-stripes border-r border-foreground/20" />

        {/* Content column */}
        <div className="relative">
          {/* Faint grid texture inside the column only. */}
          <div className="absolute inset-0 bg-grit-grid-sm opacity-90" />
          {/* Slow light drifting down each rail. */}
          <div className="rail-scan absolute left-0 h-28 w-px bg-gradient-to-b from-transparent via-primary/60 to-transparent" />
          <div className="rail-scan absolute right-0 h-28 w-px bg-gradient-to-b from-transparent via-primary/60 to-transparent" style={{ animationDelay: '4s' }} />
          {/* Crosshair markers on the rails (gentle pulse). */}
          <span className="crosshair animate-pulse-glow absolute left-0 top-[68px] -translate-x-1/2 text-foreground/45" style={{ width: 16, height: 16 }} />
          <span className="crosshair animate-pulse-glow absolute right-0 top-[68px] translate-x-1/2 text-foreground/45" style={{ width: 16, height: 16 }} />
          <span className="crosshair animate-pulse-glow absolute left-0 bottom-24 -translate-x-1/2 text-foreground/30" style={{ width: 14, height: 14 }} />
          <span className="crosshair animate-pulse-glow absolute right-0 bottom-24 translate-x-1/2 text-foreground/30" style={{ width: 14, height: 14 }} />
        </div>

        {/* Right gutter (striped) — its left border is the right rail. */}
        <div className="bg-grit-stripes border-l border-foreground/20" />
      </div>
    </div>
  )
}
