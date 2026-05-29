'use client'

// Hubfly-inspired hub-and-spoke section. Four input source cards on the
// left feed into a central "Your Grit App" hub, which fans out to seven
// runtime capability cards on the right.
//
// Layout discipline: card vertical centers and SVG path endpoints share
// the SAME percentage anchor system, so paths land precisely on the
// cards regardless of viewport width. Implementation:
//   - The container is aspect-[5/3] (matches the 1000×600 SVG viewBox).
//   - Each card is absolutely positioned by `top: X%` with translateY(-50%)
//     to center vertically; cards have an x position that lines up with
//     the SVG endpoint exactly.
//   - SVG endpoints (SRC_X / OUT_X) sit right where the rails end, so
//     paths emerge cleanly from the card edge — no floating gap, no
//     anchor dot stranded in space.
//
// Animation: paths render statically (no scroll reveal); only the pulse
// dot rides each path forever via <animateMotion>. The whole composition
// is "live" the moment it paints.

import {
  GitBranch,
  Container as DockerIcon,
  FileCode2,
  Layers,
  Shield,
  HardDrive,
  Globe,
  Network,
  Cpu,
  Terminal as TerminalIcon,
  Activity,
} from 'lucide-react'

const SOURCES = [
  { icon: GitBranch,   label: 'Git push deploys',     hue: 'text-emerald-500' },
  { icon: DockerIcon,  label: 'Docker images',        hue: 'text-sky-500' },
  { icon: FileCode2,   label: 'Resource generator',   hue: 'text-violet-500' },
  { icon: Layers,      label: 'Compose stacks',       hue: 'text-amber-500' },
]

const OUTPUTS = [
  { icon: Shield,       label: 'Auth + RBAC',           hue: 'text-emerald-500' },
  { icon: HardDrive,    label: 'File storage',          hue: 'text-amber-500' },
  { icon: Globe,        label: 'Domains + SSL',         hue: 'text-sky-500' },
  { icon: Network,      label: 'Realtime hub',          hue: 'text-rose-500' },
  { icon: Cpu,          label: 'AI Gateway',            hue: 'text-violet-500' },
  { icon: Activity,     label: 'Pulse observability',   hue: 'text-cyan-500' },
  { icon: TerminalIcon, label: 'CLI deploy',            hue: 'text-orange-500' },
]

// Layout anchors — both the SVG and the card columns use these numbers.
// 1000×600 viewBox; 5:3 aspect ratio of the container.
const VB_W = 1000
const VB_H = 600

// Rail widths (in viewBox units)
const RAIL_W = 260
const GAP = 24   // gap between rail edge and path start

const SRC_X = RAIL_W + GAP            // 284 — left path emerges here
const OUT_X = VB_W - RAIL_W - GAP     // 716 — right path terminates here

const HUB = { x: VB_W / 2, y: VB_H / 2 } // 500, 300

// Distribute N anchors evenly down the full height with a small inset
// at top/bottom so the topmost and bottommost cards don't kiss the edges.
function evenYs(n: number, vbH = VB_H, inset = 60): number[] {
  const span = vbH - inset * 2
  return Array.from({ length: n }, (_, i) =>
    n === 1 ? vbH / 2 : inset + (span * i) / (n - 1)
  )
}

const SRC_YS = evenYs(SOURCES.length)
const OUT_YS = evenYs(OUTPUTS.length)

// Cubic Bezier from one anchor to another. `bend` controls how far the
// control points reach horizontally so the curve hooks gently rather
// than slashes across.
function bezier(x1: number, y1: number, x2: number, y2: number, bend: number) {
  return `M ${x1} ${y1} C ${x1 + bend} ${y1}, ${x2 - bend} ${y2}, ${x2} ${y2}`
}

// Convert a viewBox coord to a percentage for absolute-positioning cards.
const pct = (v: number, total: number) => `${(v / total) * 100}%`

export function HubAndSpoke() {
  return (
    <section className="relative py-24 px-6 overflow-hidden border-t border-border/40">
      {/* Layered backdrop */}
      <div className="absolute inset-0 -z-30 bg-gradient-to-b from-background via-card/30 to-background" />
      <div
        className="absolute inset-0 -z-20 opacity-[0.5]"
        style={{
          backgroundImage:
            'radial-gradient(circle at 1px 1px, hsl(var(--foreground) / 0.06) 1px, transparent 0)',
          backgroundSize: '32px 32px',
        }}
      />
      <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 -z-10 h-[400px] w-[400px] rounded-full bg-primary/[0.06] blur-[100px]" />

      <div className="max-w-6xl mx-auto">

        <div className="text-center mb-16">
          <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-5">
            <span className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse" /> One Project · Every Primitive
          </span>
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight text-foreground mb-5 leading-[1.1]">
            Ship anything. Wire everything.
          </h2>
          <p className="text-base md:text-lg text-muted-foreground max-w-2xl mx-auto leading-relaxed">
            From <code className="text-foreground/90 font-mono">git push</code> to a running
            app with auth, storage, jobs, realtime, and AI all on the same wire — no glue
            code, no dashboard tabs to hunt through.
          </p>
        </div>

        {/* Topology — relative container. SVG fills behind, cards float above. */}
        <div className="relative w-full aspect-[5/3] max-w-[1100px] mx-auto">

          {/* Animated connection lines */}
          <svg
            viewBox={`0 0 ${VB_W} ${VB_H}`}
            className="absolute inset-0 w-full h-full pointer-events-none"
            preserveAspectRatio="none"
          >
            <defs>
              <linearGradient id="hub-line-grad" x1="0" y1="0" x2="1" y2="0">
                <stop offset="0%" stopColor="hsl(var(--primary))" stopOpacity="0.25" />
                <stop offset="50%" stopColor="hsl(var(--primary))" stopOpacity="0.9" />
                <stop offset="100%" stopColor="hsl(var(--primary))" stopOpacity="0.25" />
              </linearGradient>
              <radialGradient id="hub-pulse-grad">
                <stop offset="0%" stopColor="hsl(var(--primary))" stopOpacity="1" />
                <stop offset="100%" stopColor="hsl(var(--primary))" stopOpacity="0" />
              </radialGradient>
            </defs>

            {/* Source paths (left rail → hub) */}
            {SRC_YS.map((y, i) => {
              const d = bezier(SRC_X, y, HUB.x, HUB.y, 90)
              return (
                <g key={`src-${i}`}>
                  <path
                    d={d}
                    fill="none"
                    stroke="url(#hub-line-grad)"
                    strokeWidth="1.6"
                    strokeLinecap="round"
                  />
                  <circle r="4" fill="url(#hub-pulse-grad)">
                    <animateMotion dur={`${3 + i * 0.4}s`} repeatCount="indefinite" path={d} />
                  </circle>
                </g>
              )
            })}

            {/* Output paths (hub → right rail) */}
            {OUT_YS.map((y, i) => {
              const d = bezier(HUB.x, HUB.y, OUT_X, y, 90)
              return (
                <g key={`out-${i}`}>
                  <path
                    d={d}
                    fill="none"
                    stroke="url(#hub-line-grad)"
                    strokeWidth="1.6"
                    strokeLinecap="round"
                  />
                  <circle r="4" fill="url(#hub-pulse-grad)">
                    <animateMotion dur={`${2.6 + i * 0.3}s`} repeatCount="indefinite" path={d} />
                  </circle>
                </g>
              )
            })}
          </svg>

          {/* LEFT rail — source cards. right edge sits at RAIL_W. */}
          <div
            className="absolute inset-y-0 left-0"
            style={{ width: pct(RAIL_W, VB_W) }}
          >
            {SOURCES.map((s, i) => {
              const Icon = s.icon
              return (
                <div
                  key={s.label}
                  className="absolute right-0 flex items-center gap-3 rounded-xl border-2 border-border/60 bg-card/95 backdrop-blur px-4 py-3 shadow-[0_2px_0_rgba(0,0,0,0.04),0_8px_24px_-8px_rgba(0,0,0,0.15)] dark:shadow-[0_0_0_1px_rgba(255,255,255,0.03),0_8px_24px_-8px_rgba(0,0,0,0.5)] hover:border-primary/50 hover:-translate-x-0.5 transition-all"
                  style={{
                    top: pct(SRC_YS[i], VB_H),
                    transform: 'translateY(-50%)',
                    width: 'calc(100% - 8px)',
                  }}
                >
                  <div className={`h-8 w-8 rounded-lg bg-card border border-border/50 flex items-center justify-center ${s.hue}`}>
                    <Icon className="h-4 w-4" />
                  </div>
                  <span className="text-sm font-medium text-foreground">{s.label}</span>
                </div>
              )
            })}
          </div>

          {/* CENTER hub */}
          <div
            className="absolute"
            style={{
              left: pct(HUB.x, VB_W),
              top: pct(HUB.y, VB_H),
              transform: 'translate(-50%, -50%)',
            }}
          >
            <div className="relative">
              {/* Pulsing primary halo */}
              <div className="absolute inset-0 rounded-2xl bg-primary/30 blur-2xl scale-110 animate-pulse" />

              <div className="relative rounded-2xl border-2 border-primary/40 bg-card px-6 py-5 shadow-[0_24px_64px_-16px_rgba(0,0,0,0.35),0_0_0_4px_hsl(var(--primary)/0.08)]">
                <div className="flex items-center gap-3 mb-3">
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img
                    src="/grit_logo.png"
                    alt="Grit"
                    className="h-10 w-10 rounded-xl shadow-[0_4px_12px_-2px_hsl(var(--primary)/0.4)]"
                  />
                  <div>
                    <div className="font-bold text-foreground leading-tight">Grit Framework</div>
                    <div className="text-[10px] font-mono text-muted-foreground tracking-wider uppercase">Your Project</div>
                  </div>
                </div>
                <div className="flex items-center gap-1.5">
                  <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 border border-emerald-500/30 px-2 py-0.5 text-[9px] font-mono text-emerald-600 dark:text-emerald-400">
                    <span className="h-1 w-1 rounded-full bg-emerald-500 animate-pulse" /> LIVE
                  </span>
                  <span className="inline-flex items-center gap-1 rounded-full bg-primary/10 border border-primary/30 px-2 py-0.5 text-[9px] font-mono text-primary">
                    v3.23
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* RIGHT rail — output cards. left edge sits at (VB_W - RAIL_W). */}
          <div
            className="absolute inset-y-0 right-0"
            style={{ width: pct(RAIL_W, VB_W) }}
          >
            {OUTPUTS.map((o, i) => {
              const Icon = o.icon
              return (
                <div
                  key={o.label}
                  className="absolute left-0 flex items-center gap-3 rounded-xl border-2 border-border/60 bg-card/95 backdrop-blur px-4 py-2.5 shadow-[0_2px_0_rgba(0,0,0,0.04),0_8px_24px_-8px_rgba(0,0,0,0.15)] dark:shadow-[0_0_0_1px_rgba(255,255,255,0.03),0_8px_24px_-8px_rgba(0,0,0,0.5)] hover:border-primary/50 hover:translate-x-0.5 transition-all"
                  style={{
                    top: pct(OUT_YS[i], VB_H),
                    transform: 'translateY(-50%)',
                    width: 'calc(100% - 8px)',
                  }}
                >
                  <div className={`h-7 w-7 rounded-lg bg-card border border-border/50 flex items-center justify-center ${o.hue}`}>
                    <Icon className="h-3.5 w-3.5" />
                  </div>
                  <span className="text-[13px] font-medium text-foreground">{o.label}</span>
                </div>
              )
            })}
          </div>

        </div>

      </div>
    </section>
  )
}
