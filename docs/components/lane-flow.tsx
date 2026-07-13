// LaneFlow — a Next.js-docs-quality architecture diagram.
//
// Swim-lane columns, numbered/tinted boxes, optional grouped containers, and
// orthogonal elbow-arrow connectors that route deterministically. Rendered as a
// single responsive SVG (viewBox) so connectors stay pixel-locked to boxes and
// the whole thing scales to its container. Theme-aware via hsl(var(--…)).
//
//   <LaneFlow
//     id="auth"
//     lanes={['Client', 'Server', 'Auth · Database']}
//     nodes={[
//       { id: 'ui',   lane: 0, row: 0, title: 'React UI',    sub: 'login form', tone: 'blue' },
//       { id: 'gin',  lane: 1, row: 0, title: 'Gin Router',  tone: 'primary' },
//       { id: 'db',   lane: 2, row: 1, title: 'PostgreSQL',  sub: 'users table', tone: 'cyan' },
//     ]}
//     edges={[
//       { from: 'ui',  to: 'gin', label: 'POST /login' },
//       { from: 'gin', to: 'db',  label: 'lookup' },
//     ]}
//     groups={[{ lane: 1, rows: [0, 1], label: 'Session Management', badge: 1, tone: 'primary' }]}
//     legend={[{ tone: 'primary', label: 'Grit code' }, { tone: 'cyan', label: 'Managed' }]}
//   />

export type FlowTone =
  | 'primary'
  | 'blue'
  | 'green'
  | 'amber'
  | 'rose'
  | 'cyan'
  | 'violet'
  | 'default'

export interface FlowNode {
  id: string
  lane: number
  row: number
  title: string
  sub?: string
  tone?: FlowTone
  /** Small number/letter badge on the box's top-left corner. */
  badge?: number | string
}

export interface FlowEdge {
  from: string
  to: string
  label?: string
  dashed?: boolean
  tone?: FlowTone
}

export interface FlowGroup {
  lane: number
  /** Inclusive [startRow, endRow]. */
  rows: [number, number]
  label: string
  badge?: number | string
  tone?: FlowTone
}

// tone → HSL parts (hue sat% light%). Fills use low alpha; strokes are solid.
const TONES: Record<FlowTone, string> = {
  primary: 'var(--primary)',
  blue: '217 91% 60%',
  green: '160 84% 39%',
  amber: '38 92% 50%',
  rose: '350 89% 60%',
  cyan: '189 94% 43%',
  violet: '258 90% 66%',
  default: '215 16% 55%',
}

const hsl = (tone: FlowTone, alpha = 1) =>
  `hsl(${TONES[tone] ?? TONES.default}${alpha === 1 ? '' : ` / ${alpha}`})`

// Geometry (viewBox units).
const LANE_W = 232
const HEADER_H = 40
const ROW_H = 104
const PAD = 20
const BOX_W = 178
const BOX_H = 56
const R = 12 // elbow corner radius

export function LaneFlow({
  id = 'flow',
  lanes,
  nodes,
  edges = [],
  groups = [],
  legend,
  caption,
}: {
  id?: string
  lanes: string[]
  nodes: FlowNode[]
  edges?: FlowEdge[]
  groups?: FlowGroup[]
  legend?: { tone: FlowTone; label: string }[]
  caption?: string
}) {
  const numRows = Math.max(0, ...nodes.map((n) => n.row)) + 1
  const W = lanes.length * LANE_W
  const gridTop = HEADER_H + PAD
  const H = gridTop + numRows * ROW_H + PAD

  const laneCenterX = (lane: number) => lane * LANE_W + LANE_W / 2
  const rowCenterY = (row: number) => gridTop + row * ROW_H + ROW_H / 2

  const byId = new Map(nodes.map((n) => [n.id, n]))
  const cx = (n: FlowNode) => laneCenterX(n.lane)
  const cy = (n: FlowNode) => rowCenterY(n.row)

  const arrowId = `${id}-arrow`

  // Orthogonal connector path between two boxes, with rounded corners.
  function edgePath(a: FlowNode, b: FlowNode) {
    const ax = cx(a)
    const ay = cy(a)
    const bx = cx(b)
    const by = cy(b)

    // same lane → vertical
    if (a.lane === b.lane) {
      const down = by > ay
      const y1 = ay + (down ? BOX_H / 2 : -BOX_H / 2)
      const y2 = by + (down ? -BOX_H / 2 : BOX_H / 2)
      return { d: `M ${ax} ${y1} L ${ax} ${y2}`, mid: [ax, (y1 + y2) / 2] as const }
    }
    // same row → horizontal
    if (a.row === b.row) {
      const right = bx > ax
      const x1 = ax + (right ? BOX_W / 2 : -BOX_W / 2)
      const x2 = bx + (right ? -BOX_W / 2 : BOX_W / 2)
      return { d: `M ${x1} ${ay} L ${x2} ${ay}`, mid: [(x1 + x2) / 2, ay] as const }
    }
    // elbow across lanes: exit horizontally toward the target, drop/rise on a
    // shared "bus" column between the lanes, then enter the target horizontally.
    // Fan-outs from one box to a column of targets share the bus and look clean.
    const hDir = bx > ax ? 1 : -1
    const x1 = ax + (hDir * BOX_W) / 2
    const x2 = bx - (hDir * BOX_W) / 2
    const busX = (x1 + x2) / 2
    const vDir = by > ay ? 1 : -1
    const d = [
      `M ${x1} ${ay}`,
      `L ${busX - R * hDir} ${ay}`,
      `Q ${busX} ${ay} ${busX} ${ay + R * vDir}`,
      `L ${busX} ${by - R * vDir}`,
      `Q ${busX} ${by} ${busX + R * hDir} ${by}`,
      `L ${x2} ${by}`,
    ].join(' ')
    return { d, mid: [busX, (ay + by) / 2] as const }
  }

  return (
    <figure className="my-8 not-prose">
      <div className="overflow-x-auto rounded-xl border border-border/60 bg-card/40 p-2">
        <svg
          viewBox={`0 0 ${W} ${H}`}
          width="100%"
          style={{ minWidth: Math.min(W, 560), maxWidth: W, display: 'block', margin: '0 auto' }}
          role="img"
        >
          <defs>
            <marker
              id={arrowId}
              markerWidth="9"
              markerHeight="9"
              refX="7"
              refY="4.5"
              orient="auto"
              markerUnits="userSpaceOnUse"
            >
              <path d="M1 1 L8 4.5 L1 8 Z" fill="context-stroke" />
            </marker>
          </defs>

          {/* swim-lane backdrops + headers */}
          {lanes.map((lane, i) => (
            <g key={`lane-${i}`}>
              {i > 0 && (
                <line
                  x1={i * LANE_W}
                  y1={HEADER_H - 8}
                  x2={i * LANE_W}
                  y2={H - 8}
                  stroke="hsl(var(--border))"
                  strokeWidth="1"
                  strokeDasharray="4 5"
                />
              )}
              <text
                x={laneCenterX(i)}
                y={22}
                textAnchor="middle"
                fontSize="11"
                fontWeight="700"
                letterSpacing="0.08em"
                fill="hsl(var(--muted-foreground))"
                style={{ textTransform: 'uppercase' } as React.CSSProperties}
              >
                {lane}
              </text>
            </g>
          ))}

          {/* group containers (behind boxes) */}
          {groups.map((g, i) => {
            const x = laneCenterX(g.lane) - BOX_W / 2 - 14
            const yTop = rowCenterY(g.rows[0]) - BOX_H / 2 - 22
            const yBot = rowCenterY(g.rows[1]) + BOX_H / 2 + 12
            const tone = g.tone ?? 'default'
            return (
              <g key={`grp-${i}`}>
                <rect
                  x={x}
                  y={yTop}
                  width={BOX_W + 28}
                  height={yBot - yTop}
                  rx="14"
                  fill={hsl(tone, 0.04)}
                  stroke={hsl(tone, 0.5)}
                  strokeWidth="1"
                  strokeDasharray="5 5"
                />
                {g.badge != null && (
                  <>
                    <circle cx={x + 16} cy={yTop + 2} r="10" fill={hsl(tone, 1)} />
                    <text
                      x={x + 16}
                      y={yTop + 2}
                      textAnchor="middle"
                      dominantBaseline="central"
                      fontSize="11"
                      fontWeight="700"
                      fill="hsl(var(--background))"
                    >
                      {g.badge}
                    </text>
                  </>
                )}
                <text
                  x={g.badge != null ? x + 32 : x + 14}
                  y={yTop + 3}
                  dominantBaseline="central"
                  fontSize="11"
                  fontWeight="700"
                  fill={hsl(tone, 1)}
                >
                  {g.label}
                </text>
              </g>
            )
          })}

          {/* connectors */}
          {edges.map((e, i) => {
            const a = byId.get(e.from)
            const b = byId.get(e.to)
            if (!a || !b) return null
            const { d, mid } = edgePath(a, b)
            const tone = e.tone ?? 'default'
            return (
              <g key={`edge-${i}`}>
                <path
                  d={d}
                  fill="none"
                  stroke={hsl(tone, 0.85)}
                  strokeWidth="1.6"
                  strokeDasharray={e.dashed ? '5 4' : undefined}
                  markerEnd={`url(#${arrowId})`}
                />
                {e.label && (
                  <g>
                    <rect
                      x={mid[0] - e.label.length * 3.15 - 5}
                      y={mid[1] - 9}
                      width={e.label.length * 6.3 + 10}
                      height="18"
                      rx="5"
                      fill="hsl(var(--card))"
                      stroke="hsl(var(--border))"
                      strokeWidth="0.75"
                    />
                    <text
                      x={mid[0]}
                      y={mid[1] + 1}
                      textAnchor="middle"
                      dominantBaseline="central"
                      fontSize="10"
                      fontFamily="var(--font-mono, ui-monospace, monospace)"
                      fill="hsl(var(--muted-foreground))"
                    >
                      {e.label}
                    </text>
                  </g>
                )}
              </g>
            )
          })}

          {/* boxes */}
          {nodes.map((n) => {
            const x = cx(n) - BOX_W / 2
            const y = cy(n) - BOX_H / 2
            const tone = n.tone ?? 'default'
            return (
              <g key={n.id}>
                <rect
                  x={x}
                  y={y}
                  width={BOX_W}
                  height={BOX_H}
                  rx="10"
                  fill={hsl(tone, 0.1)}
                  stroke={hsl(tone, 0.6)}
                  strokeWidth="1.25"
                />
                <rect x={x} y={y} width="3.5" height={BOX_H} rx="1.75" fill={hsl(tone, 1)} />
                <text
                  x={cx(n)}
                  y={n.sub ? cy(n) - 7 : cy(n)}
                  textAnchor="middle"
                  dominantBaseline="central"
                  fontSize="12.5"
                  fontWeight="600"
                  fill="hsl(var(--foreground))"
                >
                  {n.title}
                </text>
                {n.sub && (
                  <text
                    x={cx(n)}
                    y={cy(n) + 10}
                    textAnchor="middle"
                    dominantBaseline="central"
                    fontSize="10.5"
                    fontFamily="var(--font-mono, ui-monospace, monospace)"
                    fill="hsl(var(--muted-foreground))"
                  >
                    {n.sub}
                  </text>
                )}
                {n.badge != null && (
                  <>
                    <circle cx={x} cy={y} r="9.5" fill={hsl(tone, 1)} />
                    <text
                      x={x}
                      y={y}
                      textAnchor="middle"
                      dominantBaseline="central"
                      fontSize="11"
                      fontWeight="700"
                      fill="hsl(var(--background))"
                    >
                      {n.badge}
                    </text>
                  </>
                )}
              </g>
            )
          })}
        </svg>
      </div>

      {(legend || caption) && (
        <figcaption className="mt-3 flex flex-wrap items-center justify-between gap-3 px-1 text-xs text-muted-foreground">
          {legend && (
            <div className="flex flex-wrap items-center gap-x-4 gap-y-1.5">
              {legend.map((l, i) => (
                <span key={i} className="inline-flex items-center gap-1.5">
                  <span
                    className="inline-block h-2.5 w-2.5 rounded-[3px]"
                    style={{ background: hsl(l.tone, 0.9) }}
                  />
                  {l.label}
                </span>
              ))}
            </div>
          )}
          {caption && <span className="italic opacity-80">{caption}</span>}
        </figcaption>
      )}
    </figure>
  )
}
