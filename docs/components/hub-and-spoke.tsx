'use client'

// Hubfly-inspired hub-and-spoke section. Four input source cards on the
// left feed into a central "Your Grit App" hub, which fans out to six
// runtime capability cards on the right. Connection paths are animated
// SVG (draw-in on viewport entry, then a pulse dot loops along each
// path so the whole thing feels alive). Pure SVG + framer-motion, no
// canvas, no GSAP — keeps it snappy on mobile.

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
import { motion } from 'framer-motion'

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

// Build a cubic Bezier path between two anchor points. The control
// points are offset horizontally so the curve always hooks into the
// hub from the side rather than crossing it.
function bezier(x1: number, y1: number, x2: number, y2: number, bend = 80) {
  return `M ${x1} ${y1} C ${x1 + bend} ${y1}, ${x2 - bend} ${y2}, ${x2} ${y2}`
}

export function HubAndSpoke() {
  // Card layout in SVG coordinates. The SVG is 1000×600 viewBox so the
  // values are unitless and scale with the parent width.
  const HUB = { x: 500, y: 300 }
  const SRC_X = 220
  const OUT_X = 780

  const sourceY = (i: number, total: number) => 120 + (i * 360) / Math.max(total - 1, 1)
  const outY = (i: number, total: number) => 80 + (i * 440) / Math.max(total - 1, 1)

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
      {/* Center halo */}
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

        {/* Topology — relative container; SVG fills behind, cards float above */}
        <div className="relative w-full aspect-[5/3] max-w-[1100px] mx-auto">

          {/* Animated connection lines */}
          <svg
            viewBox="0 0 1000 600"
            className="absolute inset-0 w-full h-full pointer-events-none"
            preserveAspectRatio="none"
          >
            <defs>
              <linearGradient id="hub-line-grad" x1="0" y1="0" x2="1" y2="0">
                <stop offset="0%" stopColor="hsl(var(--primary))" stopOpacity="0.2" />
                <stop offset="50%" stopColor="hsl(var(--primary))" stopOpacity="0.8" />
                <stop offset="100%" stopColor="hsl(var(--primary))" stopOpacity="0.2" />
              </linearGradient>
              <radialGradient id="hub-pulse-grad">
                <stop offset="0%" stopColor="hsl(var(--primary))" stopOpacity="1" />
                <stop offset="100%" stopColor="hsl(var(--primary))" stopOpacity="0" />
              </radialGradient>
            </defs>

            {/* Source paths (left → hub) */}
            {SOURCES.map((_, i) => {
              const y = sourceY(i, SOURCES.length)
              const d = bezier(SRC_X, y, HUB.x, HUB.y, 90)
              return (
                <g key={`src-${i}`}>
                  <motion.path
                    d={d}
                    fill="none"
                    stroke="url(#hub-line-grad)"
                    strokeWidth="1.6"
                    strokeLinecap="round"
                    initial={{ pathLength: 0, opacity: 0 }}
                    whileInView={{ pathLength: 1, opacity: 1 }}
                    viewport={{ once: true, margin: '-100px' }}
                    transition={{ duration: 1.2, delay: 0.15 * i, ease: 'easeInOut' }}
                  />
                  {/* Pulse dot riding the path */}
                  <circle r="4" fill="url(#hub-pulse-grad)">
                    <animateMotion dur={`${3 + i * 0.4}s`} repeatCount="indefinite" path={d} />
                  </circle>
                </g>
              )
            })}

            {/* Output paths (hub → right) */}
            {OUTPUTS.map((_, i) => {
              const y = outY(i, OUTPUTS.length)
              const d = bezier(HUB.x, HUB.y, OUT_X, y, 90)
              return (
                <g key={`out-${i}`}>
                  <motion.path
                    d={d}
                    fill="none"
                    stroke="url(#hub-line-grad)"
                    strokeWidth="1.6"
                    strokeLinecap="round"
                    initial={{ pathLength: 0, opacity: 0 }}
                    whileInView={{ pathLength: 1, opacity: 1 }}
                    viewport={{ once: true, margin: '-100px' }}
                    transition={{ duration: 1.2, delay: 0.6 + 0.1 * i, ease: 'easeInOut' }}
                  />
                  <circle r="4" fill="url(#hub-pulse-grad)">
                    <animateMotion dur={`${2.6 + i * 0.3}s`} repeatCount="indefinite" path={d} />
                  </circle>
                </g>
              )
            })}

            {/* Anchor dots where lines meet cards */}
            {SOURCES.map((_, i) => (
              <circle key={`sa-${i}`} cx={SRC_X} cy={sourceY(i, SOURCES.length)} r="3" fill="hsl(var(--primary))" />
            ))}
            {OUTPUTS.map((_, i) => (
              <circle key={`oa-${i}`} cx={OUT_X} cy={outY(i, OUTPUTS.length)} r="3" fill="hsl(var(--primary))" />
            ))}
          </svg>

          {/* LEFT column — source cards */}
          <div className="absolute inset-y-0 left-0 w-[28%] flex flex-col justify-around">
            {SOURCES.map((s, i) => {
              const Icon = s.icon
              return (
                <motion.div
                  key={s.label}
                  initial={{ opacity: 0, x: -20 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true, margin: '-100px' }}
                  transition={{ duration: 0.5, delay: 0.1 + 0.1 * i }}
                  className="group flex items-center gap-3 rounded-xl border-2 border-border/60 bg-card/90 backdrop-blur px-4 py-3 shadow-[0_2px_0_rgba(0,0,0,0.04),0_8px_24px_-8px_rgba(0,0,0,0.15)] dark:shadow-[0_0_0_1px_rgba(255,255,255,0.03),0_8px_24px_-8px_rgba(0,0,0,0.5)] hover:border-primary/50 hover:-translate-x-0.5 transition-all"
                >
                  <div className={`h-8 w-8 rounded-lg bg-card border border-border/50 flex items-center justify-center ${s.hue}`}>
                    <Icon className="h-4 w-4" />
                  </div>
                  <span className="text-sm font-medium text-foreground">{s.label}</span>
                </motion.div>
              )
            })}
          </div>

          {/* CENTER hub */}
          <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2">
            <motion.div
              initial={{ opacity: 0, scale: 0.85 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true, margin: '-100px' }}
              transition={{ duration: 0.6, delay: 0.4, ease: 'backOut' }}
              className="relative"
            >
              {/* Pulsing primary halo */}
              <div className="absolute inset-0 rounded-2xl bg-primary/30 blur-2xl scale-110 animate-pulse" />

              <div className="relative rounded-2xl border-2 border-primary/40 bg-card px-6 py-5 shadow-[0_24px_64px_-16px_rgba(0,0,0,0.35),0_0_0_4px_hsl(var(--primary)/0.08)]">
                <div className="flex items-center gap-3 mb-3">
                  <div className="h-10 w-10 rounded-xl bg-gradient-to-br from-primary to-sky-400 flex items-center justify-center font-bold text-white text-lg shadow-[inset_0_-2px_0_rgba(0,0,0,0.15),0_4px_12px_-2px_hsl(var(--primary)/0.4)]">
                    G
                  </div>
                  <div>
                    <div className="font-bold text-foreground leading-tight">Grit Framework</div>
                    <div className="text-[10px] font-mono text-muted-foreground tracking-wider uppercase">Your Project</div>
                  </div>
                </div>
                {/* Status pills */}
                <div className="flex items-center gap-1.5">
                  <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 border border-emerald-500/30 px-2 py-0.5 text-[9px] font-mono text-emerald-600 dark:text-emerald-400">
                    <span className="h-1 w-1 rounded-full bg-emerald-500 animate-pulse" /> LIVE
                  </span>
                  <span className="inline-flex items-center gap-1 rounded-full bg-primary/10 border border-primary/30 px-2 py-0.5 text-[9px] font-mono text-primary">
                    v3.23
                  </span>
                </div>
              </div>
            </motion.div>
          </div>

          {/* RIGHT column — output cards */}
          <div className="absolute inset-y-0 right-0 w-[28%] flex flex-col justify-around">
            {OUTPUTS.map((o, i) => {
              const Icon = o.icon
              return (
                <motion.div
                  key={o.label}
                  initial={{ opacity: 0, x: 20 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true, margin: '-100px' }}
                  transition={{ duration: 0.5, delay: 0.6 + 0.08 * i }}
                  className="group flex items-center gap-3 rounded-xl border-2 border-border/60 bg-card/90 backdrop-blur px-4 py-2.5 shadow-[0_2px_0_rgba(0,0,0,0.04),0_8px_24px_-8px_rgba(0,0,0,0.15)] dark:shadow-[0_0_0_1px_rgba(255,255,255,0.03),0_8px_24px_-8px_rgba(0,0,0,0.5)] hover:border-primary/50 hover:translate-x-0.5 transition-all"
                >
                  <div className={`h-7 w-7 rounded-lg bg-card border border-border/50 flex items-center justify-center ${o.hue}`}>
                    <Icon className="h-3.5 w-3.5" />
                  </div>
                  <span className="text-[13px] font-medium text-foreground">{o.label}</span>
                </motion.div>
              )
            })}
          </div>

        </div>

      </div>
    </section>
  )
}
