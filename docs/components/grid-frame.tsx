// Blueprint frame (Livewire-style). A single fixed, full-viewport backdrop that
// gives marketing pages their signature graph-paper framing:
//   • two vertical rails that frame the centred content column and — because the
//     layer is `fixed` — run unbroken from the top of the header to the footer;
//   • a full-bleed grid texture behind every section;
//   • a soft primary gradient glow at the top for depth;
//   • crosshair "+" markers sitting on the rails.
// Mounted once per marketing page (landing, pitch, courses). It sits at -z-10:
// above the page's solid background, behind all content. Sections with their own
// translucent tints let the rails read through, creating the "wrapped cells" look.
//
// `width` must match the page's content container so the rails hug it (default
// matches Tailwind's max-w-6xl).
export function GridFrame({ width = 'max-w-6xl' }: { width?: string }) {
  return (
    <div aria-hidden className="pointer-events-none fixed inset-0 -z-10 overflow-hidden">
      {/* Full-bleed grid texture */}
      <div
        className="absolute inset-0"
        style={{
          backgroundImage:
            'linear-gradient(to right, hsl(var(--foreground) / 0.07) 1px, transparent 1px), linear-gradient(to bottom, hsl(var(--foreground) / 0.07) 1px, transparent 1px)',
          backgroundSize: '56px 56px',
        }}
      />
      {/* Soft primary glow at the top — the only "colour" moment, keeps the base consistent */}
      <div
        className="absolute inset-x-0 top-0 h-[640px]"
        style={{
          background:
            'radial-gradient(ellipse 60% 100% at 50% -10%, hsl(var(--primary) / 0.12), transparent 65%)',
        }}
      />
      {/* Centred vertical rails + corner crosshairs */}
      <div className={`relative mx-auto h-full w-full ${width}`}>
        <div className="h-full border-x border-border/70" />
        <span className="crosshair absolute left-0 top-[68px] -translate-x-1/2 text-foreground/30" style={{ width: 15, height: 15 }} />
        <span className="crosshair absolute right-0 top-[68px] translate-x-1/2 text-foreground/30" style={{ width: 15, height: 15 }} />
      </div>
    </div>
  )
}
