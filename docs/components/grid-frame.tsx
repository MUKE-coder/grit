// Blueprint frame (Livewire-style). A full-page backdrop that gives marketing
// pages their signature graph-paper framing:
//   • two vertical rails that frame the centred content column and run the full
//     height of the page (top of the header to the footer);
//   • a full-bleed grid texture behind every section;
//   • a soft primary gradient glow at the top for depth;
//   • crosshair "+" markers sitting on the rails.
// Mounted once per page (landing, pitch, courses, …). It is `absolute inset-0`
// inside the page's `relative isolate` wrapper, so it spans the whole document
// AND — crucially — centres its rails in the EXACT same context as flow content
// (wrapper width = viewport minus scrollbar). A `fixed` layer instead centres on
// the full viewport, leaving the rails ~scrollbar/2 off the section edges on big
// screens. Absolute keeps them perfectly aligned at every screen size.
// It sits at -z-10: above the wrapper's background, behind all content.
//
// `width` must match the page's content container so the rails hug it (default
// matches Tailwind's max-w-6xl).
export function GridFrame({ width = 'max-w-6xl' }: { width?: string }) {
  return (
    <div aria-hidden className="pointer-events-none absolute inset-0 -z-10 overflow-hidden">
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
      {/* Two structural rails framing the content column. On a fixed layer so
          they run unbroken from the header to the footer. Drawn in the FOREGROUND
          colour (the `border` token is nearly invisible on the dark base) so the
          rails actually read in both themes. */}
      <div className={`relative mx-auto h-full w-full ${width}`}>
        <div className="absolute inset-y-0 left-0 w-px bg-foreground/20" />
        <div className="absolute inset-y-0 right-0 w-px bg-foreground/20" />
        <span className="crosshair absolute left-0 top-[68px] -translate-x-1/2 text-foreground/40" style={{ width: 16, height: 16 }} />
        <span className="crosshair absolute right-0 top-[68px] translate-x-1/2 text-foreground/40" style={{ width: 16, height: 16 }} />
        <span className="crosshair absolute left-0 bottom-24 -translate-x-1/2 text-foreground/30" style={{ width: 14, height: 14 }} />
        <span className="crosshair absolute right-0 bottom-24 translate-x-1/2 text-foreground/30" style={{ width: 14, height: 14 }} />
      </div>
    </div>
  )
}
