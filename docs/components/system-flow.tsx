import type { ReactNode } from 'react'

/* ─────────────────────────────────────────────────────────────
   Request-flow diagrams for the system design pages.

   A boxes-and-arrows picture of how one request moves through a
   system, with each hop numbered so the prose underneath can
   refer to it. Declared as nodes on a grid and edges between
   them, and drawn as one inline SVG: a grid of HTML boxes
   cannot draw an arrow that turns a corner, and an image cannot
   change colour with the theme or be read by a screen reader.

   Everything here is presentational and server-rendered.
   ───────────────────────────────────────────────────────────── */

/** Where a node sits. Columns run left to right, rows top to bottom. */
export interface FlowNode {
  id: string
  label: string
  /** Second line, for a qualifier like "(Web/Mobile)". */
  sub?: string
  col: number
  row: number
  /** Drawn in the accent colour: the component the page is about. */
  accent?: boolean
}

export interface FlowEdge {
  from: string
  to: string
  /** The number in the badge. Omit for an unnumbered link. */
  step?: number
  /**
   * Dashed, for a hop that is conditional or asynchronous: an MFA
   * challenge that only some sign-ins need, a write to an audit log that
   * the response does not wait for.
   */
  dashed?: boolean
  /**
   * Which way the line leaves and enters. "h" goes across then down,
   * "v" goes down then across. The default picks whichever suits the
   * relative positions.
   */
  bend?: 'h' | 'v'
  /** Nudges the badge along the line, 0 to 1. Defaults to the middle. */
  at?: number
}

const COL_W = 220
const COL_GAP = 56
const ROW_H = 72
const ROW_GAP = 60
const PAD = 16

function nodeBox(n: FlowNode) {
  const x = PAD + n.col * (COL_W + COL_GAP)
  const y = PAD + n.row * (ROW_H + ROW_GAP)
  return { x, y, w: COL_W, h: ROW_H, cx: x + COL_W / 2, cy: y + ROW_H / 2 }
}

interface Pt {
  x: number
  y: number
}

/** Clearance kept between a line and a box it is not attached to. */
const CLEAR = 6

/** Does this axis-aligned segment run through a box it has no business in? */
function segmentHits(p: Pt, q: Pt, n: FlowNode) {
  const B = nodeBox(n)
  const left = B.x - CLEAR
  const right = B.x + B.w + CLEAR
  const top = B.y - CLEAR
  const bottom = B.y + B.h + CLEAR

  if (p.y === q.y) {
    if (p.y <= top || p.y >= bottom) return false
    return Math.max(Math.min(p.x, q.x), left) < Math.min(Math.max(p.x, q.x), right)
  }
  if (p.x === q.x) {
    if (p.x <= left || p.x >= right) return false
    return Math.max(Math.min(p.y, q.y), top) < Math.min(Math.max(p.y, q.y), bottom)
  }
  return false
}

function routeIsClear(pts: Pt[], a: FlowNode, b: FlowNode, nodes: FlowNode[]) {
  return nodes.every((n) => {
    if (n.id === a.id || n.id === b.id) return true
    for (let i = 1; i < pts.length; i++) {
      if (segmentHits(pts[i - 1], pts[i], n)) return false
    }
    return true
  })
}

/**
 * The polyline from one box to another.
 *
 * A request flow is nearly always a straight run along a row or down a column,
 * and those are drawn as themselves. Everything else turns a corner, and where
 * it turns matters: a line that takes the direct corner can run straight
 * through a box that happens to sit between the two it joins, which reads as a
 * connection that is not there.
 *
 * So the diagonal cases are a list of candidate routes, tried in order, and the
 * first one that touches nothing else is the one drawn. The caller's `bend` is
 * the first candidate, so a layout whose direct route is already clear is drawn
 * exactly as it was before; the alternatives only appear where the direct route
 * would have lied. The last two run in the gutter between columns or rows,
 * which is empty by construction, and are what makes a clear route almost
 * always available.
 */
function routePoints(a: FlowNode, b: FlowNode, bend: 'h' | 'v' | undefined, nodes: FlowNode[]): Pt[] {
  const A = nodeBox(a)
  const B = nodeBox(b)

  if (a.row === b.row) {
    const left = A.x < B.x
    return [
      { x: left ? A.x + A.w : A.x, y: A.cy },
      { x: left ? B.x : B.x + B.w, y: B.cy },
    ]
  }
  if (a.col === b.col) {
    const down = A.y < B.y
    return [
      { x: A.cx, y: down ? A.y + A.h : A.y },
      { x: B.cx, y: down ? B.y : B.y + B.h },
    ]
  }

  const right = A.x < B.x
  const down = A.y < B.y
  const exitX = right ? A.x + A.w : A.x
  const exitY = down ? A.y + A.h : A.y
  const enterX = right ? B.x : B.x + B.w
  const enterY = down ? B.y : B.y + B.h
  // The empty lane beside A, and the one beside B. Columns and rows are laid
  // out with a gap, and nothing is ever drawn in it.
  const laneX = right ? A.x + A.w + COL_GAP / 2 : A.x - COL_GAP / 2
  const laneY = down ? A.y + A.h + ROW_GAP / 2 : A.y - ROW_GAP / 2
  const nearX = right ? B.x - COL_GAP / 2 : B.x + B.w + COL_GAP / 2
  const nearY = down ? B.y - ROW_GAP / 2 : B.y + B.h + ROW_GAP / 2

  // Leaves the side and turns, and leaves the bottom and turns: the two direct
  // corners. Then the same two taken through a gutter instead of through the
  // target's centre line.
  const h: Pt[] = [
    { x: exitX, y: A.cy },
    { x: B.cx, y: A.cy },
    { x: B.cx, y: enterY },
  ]
  const v: Pt[] = [
    { x: A.cx, y: exitY },
    { x: A.cx, y: B.cy },
    { x: enterX, y: B.cy },
  ]
  const hLane: Pt[] = [
    { x: exitX, y: A.cy },
    { x: laneX, y: A.cy },
    { x: laneX, y: B.cy },
    { x: enterX, y: B.cy },
  ]
  const vLane: Pt[] = [
    { x: A.cx, y: exitY },
    { x: A.cx, y: laneY },
    { x: B.cx, y: laneY },
    { x: B.cx, y: enterY },
  ]
  const hNear: Pt[] = [
    { x: exitX, y: A.cy },
    { x: nearX, y: A.cy },
    { x: nearX, y: B.cy },
    { x: enterX, y: B.cy },
  ]
  const vNear: Pt[] = [
    { x: A.cx, y: exitY },
    { x: A.cx, y: nearY },
    { x: B.cx, y: nearY },
    { x: B.cx, y: enterY },
  ]

  const preferV = bend === 'v'
  const candidates = preferV
    ? [v, vLane, vNear, h, hLane, hNear]
    : [h, hLane, hNear, v, vLane, vNear]

  return candidates.find((pts) => routeIsClear(pts, a, b, nodes)) ?? candidates[0]
}

function pathD(pts: Pt[]) {
  return pts.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`).join(' ')
}

/**
 * A point some fraction along a polyline, measured by length rather than by
 * segment count, so a badge sits where it looks like the middle whichever
 * route was chosen.
 */
function pointAlong(pts: Pt[], t: number): Pt {
  const lengths = pts.slice(1).map((p, i) => Math.hypot(p.x - pts[i].x, p.y - pts[i].y))
  const total = lengths.reduce((sum, l) => sum + l, 0)
  if (total === 0) return pts[0]

  let travelled = total * Math.min(Math.max(t, 0), 1)
  for (let i = 0; i < lengths.length; i++) {
    if (travelled <= lengths[i] || i === lengths.length - 1) {
      const f = lengths[i] === 0 ? 0 : travelled / lengths[i]
      return {
        x: pts[i].x + (pts[i + 1].x - pts[i].x) * f,
        y: pts[i].y + (pts[i + 1].y - pts[i].y) * f,
      }
    }
    travelled -= lengths[i]
  }
  return pts[pts.length - 1]
}

export function SystemFlow({
  title,
  nodes,
  edges,
  caption,
}: {
  title: string
  nodes: FlowNode[]
  edges: FlowEdge[]
  caption?: ReactNode
}) {
  const byId = new Map(nodes.map((n) => [n.id, n]))
  const cols = Math.max(...nodes.map((n) => n.col)) + 1
  const rows = Math.max(...nodes.map((n) => n.row)) + 1
  const width = PAD * 2 + cols * COL_W + (cols - 1) * COL_GAP
  const height = PAD * 2 + rows * ROW_H + (rows - 1) * ROW_GAP

  // A sentence of the same thing, for anybody who cannot see the picture.
  // Without this the diagram is a decoration to a screen reader, and the
  // numbered steps in the prose refer to nothing.
  const described = edges
    .filter((e) => e.step !== undefined)
    .sort((a, b) => (a.step ?? 0) - (b.step ?? 0))
    .map((e) => `${e.step}. ${byId.get(e.from)?.label} to ${byId.get(e.to)?.label}`)
    .join('. ')

  return (
    <figure className="my-8">
      <div className="rounded-xl border border-border/50 bg-gradient-to-b from-card/40 to-background/40 p-4 sm:p-6">
        <p className="mb-4 text-center text-sm font-semibold tracking-tight">{title}</p>
        <div className="overflow-x-auto">
          <svg
            viewBox={`0 0 ${width} ${height}`}
            width="100%"
            style={{ minWidth: Math.min(width, 640) }}
            role="img"
            aria-label={`${title}. ${described}`}
            className="h-auto"
          >
            <defs>
              <marker
                id="sysflow-arrow"
                viewBox="0 0 10 10"
                refX="9"
                refY="5"
                markerWidth="6"
                markerHeight="6"
                orient="auto-start-reverse"
              >
                <path d="M 0 0 L 10 5 L 0 10 z" className="fill-foreground/55" />
              </marker>
            </defs>

            {edges.map((e, i) => {
              const a = byId.get(e.from)
              const b = byId.get(e.to)
              if (!a || !b) return null
              const pts = routePoints(a, b, e.bend, nodes)
              const badge = pointAlong(pts, e.at ?? 0.5)
              return (
                <g key={`${e.from}-${e.to}-${i}`}>
                  <path
                    d={pathD(pts)}
                    fill="none"
                    strokeWidth={1.75}
                    strokeDasharray={e.dashed ? '6 5' : undefined}
                    markerEnd="url(#sysflow-arrow)"
                    className="stroke-foreground/55"
                  />
                  {e.step !== undefined && (
                    <>
                      <circle cx={badge.x} cy={badge.y} r={13} className="fill-primary" />
                      <text
                        x={badge.x}
                        y={badge.y}
                        textAnchor="middle"
                        dominantBaseline="central"
                        className="fill-primary-foreground"
                        style={{ fontSize: 13, fontWeight: 700 }}
                      >
                        {e.step}
                      </text>
                    </>
                  )}
                </g>
              )
            })}

            {nodes.map((n) => {
              const { x, y, w, h, cx } = nodeBox(n)
              return (
                <g key={n.id}>
                  <rect
                    x={x}
                    y={y}
                    width={w}
                    height={h}
                    rx={10}
                    strokeWidth={1.5}
                    className={
                      n.accent
                        ? 'fill-primary/15 stroke-primary/60'
                        : 'fill-emerald-500/10 stroke-emerald-500/40'
                    }
                  />
                  <text
                    x={cx}
                    y={n.sub ? y + h / 2 - 9 : y + h / 2}
                    textAnchor="middle"
                    dominantBaseline="central"
                    className="fill-foreground"
                    style={{ fontSize: 15, fontWeight: 600 }}
                  >
                    {n.label}
                  </text>
                  {n.sub && (
                    <text
                      x={cx}
                      y={y + h / 2 + 11}
                      textAnchor="middle"
                      dominantBaseline="central"
                      className="fill-muted-foreground"
                      style={{ fontSize: 12.5 }}
                    >
                      {n.sub}
                    </text>
                  )}
                </g>
              )
            })}
          </svg>
        </div>
      </div>
      {caption && (
        <figcaption className="mt-2 text-center text-xs text-muted-foreground">{caption}</figcaption>
      )}
    </figure>
  )
}
