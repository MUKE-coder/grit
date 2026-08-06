'use client'

import { useState } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { ArrowRight } from 'lucide-react'
import { FrameworkLogo } from '@/components/framework-logo'
import {
  CAPABILITY_GROUPS,
  COMPARED,
  PROVISION_LABEL,
  PROVISION_SHORT,
  readyCount,
  TOTAL_ROWS,
  type Provision,
} from '@/config/capabilities'

/*
 * The capability matrix under the benchmark.
 *
 * Two layouts rather than one responsive table. A seven-column grid cannot be
 * made readable at 375px by scrolling alone: the row label is the thing you
 * need pinned, and a horizontally scrolled table hides it exactly when you
 * scroll to the column you care about. So phones get a per-row stack and wider
 * screens get the grid.
 *
 * Colour is never the only signal. Each state carries a distinct glyph and a
 * text label, so the table survives a monochrome print and a red-green
 * colour-blind reader.
 */

const STYLE: Record<Provision, { glyph: string; cls: string; dot: string }> = {
  generated: {
    glyph: '●',
    cls: 'text-primary',
    dot: 'bg-primary',
  },
  builtin: {
    glyph: '●',
    cls: 'text-success',
    dot: 'bg-success',
  },
  official: {
    glyph: '◐',
    cls: 'text-warning',
    dot: 'bg-warning',
  },
  community: {
    glyph: '○',
    cls: 'text-muted-foreground',
    dot: 'bg-muted-foreground/60',
  },
  none: {
    glyph: '·',
    cls: 'text-muted-foreground/35',
    dot: 'bg-muted-foreground/25',
  },
}

const ORDER: Provision[] = ['generated', 'builtin', 'official', 'community', 'none']

function Marker({ state, note }: { state: Provision; note?: string }) {
  const st = STYLE[state]
  const label = `${PROVISION_LABEL[state]}${note ? `. ${note}` : ''}`
  return (
    <span
      className={`inline-flex items-center justify-center text-base leading-none ${st.cls}`}
      title={label}
    >
      <span aria-hidden="true">{st.glyph}</span>
      <span className="sr-only">{label}</span>
    </span>
  )
}

export function CapabilityMatrix() {
  const [openRow, setOpenRow] = useState<string | null>(null)

  return (
    <div className="rounded-2xl border border-border/40 bg-[#0d0d12] p-5 sm:p-8">
      {/* ── Legend ─────────────────────────────────────────────── */}
      <div className="flex flex-wrap items-center justify-center gap-x-5 gap-y-2 mb-8 text-xs">
        {ORDER.map((s) => (
          <span key={s} className="inline-flex items-center gap-1.5 text-muted-foreground">
            <span className={`${STYLE[s].cls} text-base leading-none`} aria-hidden="true">
              {STYLE[s].glyph}
            </span>
            {PROVISION_SHORT[s]}
          </span>
        ))}
      </div>

      {/* ── Wide layout ────────────────────────────────────────── */}
      <div className="hidden md:block overflow-x-auto -mx-2 px-2">
        <table className="w-full border-collapse text-sm" style={{ minWidth: 720 }}>
          <caption className="sr-only">
            How each capability arrives in each framework: generated, built in, an official
            package, a third-party package, or not provided.
          </caption>
          <thead>
            <tr>
              <th scope="col" className="text-left font-medium text-muted-foreground pb-3 pr-4">
                Capability
              </th>
              {COMPARED.map((f) => (
                <th key={f.slug} scope="col" className="pb-3 px-1 align-bottom">
                  <span className="flex flex-col items-center gap-1.5">
                    {f.slug === 'grit' ? (
                      <Image
                        src={f.logo}
                        alt=""
                        width={26}
                        height={26}
                        className="rounded-[6px]"
                      />
                    ) : (
                      <FrameworkLogo
                        src={f.logo}
                        alt=""
                        onLight={'invertOnDark' in f ? f.invertOnDark : undefined}
                        className="h-[26px] max-w-[74px]"
                      />
                    )}
                    <span
                      className={
                        'text-[11px] font-medium ' +
                        (f.slug === 'grit' ? 'text-primary' : 'text-muted-foreground')
                      }
                    >
                      {f.name}
                    </span>
                  </span>
                </th>
              ))}
            </tr>
          </thead>

          {CAPABILITY_GROUPS.map((group) => (
            <tbody key={group.title}>
              <tr>
                <th
                  scope="colgroup"
                  colSpan={COMPARED.length + 1}
                  className="text-left pt-6 pb-2"
                >
                  <span className="tag-mono text-primary/80">{group.title}</span>
                </th>
              </tr>
              {group.rows.map((row) => (
                <tr key={row.id} className="border-t border-border/20">
                  <th scope="row" className="text-left font-normal py-2.5 pr-4 align-top">
                    <span className="block text-foreground">{row.label}</span>
                    <span className="block text-xs text-muted-foreground mt-0.5 max-w-md">
                      {row.detail}
                    </span>
                  </th>
                  {COMPARED.map((f) => (
                    <td
                      key={f.slug}
                      className={
                        'text-center py-2.5 px-1 align-top ' +
                        (f.slug === 'grit' ? 'bg-primary/[0.045]' : '')
                      }
                    >
                      <Marker
                        state={row.cells[f.slug]?.state ?? 'none'}
                        note={row.cells[f.slug]?.note}
                      />
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          ))}

          <tfoot>
            <tr className="border-t border-border/40">
              <th scope="row" className="text-left font-medium py-3 pr-4 text-foreground">
                Running without installing anything
              </th>
              {COMPARED.map((f) => (
                <td
                  key={f.slug}
                  className={
                    'text-center py-3 px-1 font-semibold tabular-nums ' +
                    (f.slug === 'grit'
                      ? 'bg-primary/[0.045] text-primary'
                      : 'text-muted-foreground')
                  }
                >
                  {readyCount(f.slug)}
                  <span className="text-muted-foreground/50 font-normal">/{TOTAL_ROWS}</span>
                </td>
              ))}
            </tr>
          </tfoot>
        </table>
      </div>

      {/* ── Phone layout ───────────────────────────────────────── */}
      <div className="md:hidden space-y-6">
        {CAPABILITY_GROUPS.map((group) => (
          <section key={group.title}>
            <h3 className="tag-mono text-primary/80 mb-2">{group.title}</h3>
            <ul className="space-y-1.5">
              {group.rows.map((row) => {
                const open = openRow === row.id
                return (
                  <li key={row.id} className="rounded-lg border border-border/30">
                    <button
                      type="button"
                      onClick={() => setOpenRow(open ? null : row.id)}
                      aria-expanded={open}
                      className="w-full min-h-[44px] flex items-center justify-between gap-3 px-3 py-2.5 text-left"
                    >
                      <span className="text-sm text-foreground">{row.label}</span>
                      <span className="flex items-center gap-2 shrink-0">
                        <Marker state={row.cells.grit?.state ?? 'none'} />
                        <span className="text-[10px] text-muted-foreground/70">
                          {open ? 'Hide' : 'Compare'}
                        </span>
                      </span>
                    </button>
                    {open && (
                      <div className="px-3 pb-3 pt-1 border-t border-border/20">
                        <p className="text-xs text-muted-foreground mb-2.5">{row.detail}</p>
                        <ul className="space-y-1.5">
                          {COMPARED.map((f) => {
                            const cell = row.cells[f.slug]
                            return (
                              <li key={f.slug} className="flex items-baseline gap-2 text-xs">
                                <Marker state={cell?.state ?? 'none'} />
                                <span
                                  className={
                                    'w-20 shrink-0 ' +
                                    (f.slug === 'grit'
                                      ? 'text-primary font-medium'
                                      : 'text-foreground')
                                  }
                                >
                                  {f.name}
                                </span>
                                <span className="text-muted-foreground">
                                  {cell?.note ?? PROVISION_SHORT[cell?.state ?? 'none']}
                                </span>
                              </li>
                            )
                          })}
                        </ul>
                      </div>
                    )}
                  </li>
                )
              })}
            </ul>
          </section>
        ))}

        <div className="rounded-lg border border-primary/25 bg-primary/[0.04] px-4 py-3">
          <p className="text-sm text-foreground">
            Running without installing anything:{' '}
            <strong className="text-primary">
              {readyCount('grit')}/{TOTAL_ROWS}
            </strong>{' '}
            in Grit. The next closest is {readyCount('django')}/{TOTAL_ROWS}.
          </p>
        </div>
      </div>

      {/* ── The caveat that makes the table worth trusting ─────── */}
      <div className="mt-7 pt-4 border-t border-border/40 flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3">
        <p className="text-xs text-muted-foreground leading-relaxed max-w-2xl">
          Every row here can be built in every framework listed. The table is not about what is
          possible, it is about how much is already running the first time you open the app.
          Django&apos;s admin and Laravel&apos;s queues and observability are first rate, and they
          are marked as such. Express is a router and Bun is a runtime, so Express scoring low is a
          statement about scope and not about quality, which is also why Bun is not a column here
          at all.
        </p>
        <Link
          href="/docs/batteries"
          className="shrink-0 inline-flex items-center gap-1 text-sm text-primary hover:underline"
        >
          What ships in the box
          <ArrowRight className="h-3.5 w-3.5" />
        </Link>
      </div>
    </div>
  )
}
