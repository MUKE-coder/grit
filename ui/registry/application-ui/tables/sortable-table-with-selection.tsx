'use client'

import { useEffect, useId, useMemo, useRef, useState } from 'react'
import { ArrowDown, ArrowUp, ChevronsUpDown, SearchX } from 'lucide-react'

/*
 * A data table with sortable columns, row selection and pagination.
 *
 * Self-contained on purpose. The obvious build is TanStack Table plus a set of
 * Table primitives, which is a good library and the wrong dependency for a
 * block whose entire job is to be a starting point you own. Nothing here
 * imports anything.
 *
 * ── Sorting ─────────────────────────────────────────────────────────────────
 * The sortable header is a <th> carrying aria-sort, with a <button> inside it.
 * Both halves matter and both are usually missing. aria-sort is what makes a
 * screen reader say "sorted ascending" when it reaches the column; the button
 * is what makes the header reachable and operable by keyboard at all. A <th>
 * with an onClick is a cell nobody can press.
 *
 * aria-sort goes on exactly one header. The attribute means "this table is
 * sorted by this column", so putting "none" on every other column is fine but
 * putting "ascending" on two is a contradiction.
 *
 * ── Selection ───────────────────────────────────────────────────────────────
 * The select-all checkbox uses the indeterminate property when the selection
 * is partial. indeterminate cannot be set from markup: it is a property, not
 * an attribute, so it needs a ref and an effect. Skipping it leaves a checkbox
 * that reads as unchecked while three rows are selected.
 *
 * Every row checkbox names its row. Twelve checkboxes all called "Select row"
 * are twelve identical entries in a screen reader's control list.
 *
 * ── Announcements ───────────────────────────────────────────────────────────
 * Sorting and selection both change the table without moving focus, so both go
 * through one status region. Sorting silently reorders rows under a keyboard
 * user's cursor otherwise, which is disorienting rather than merely quiet.
 *
 * Money is integer cents, formatted once with Intl.
 */

export interface Row {
  id: string
  name: string
  email: string
  status: 'Paid' | 'Pending' | 'Failed'
  /** Integer cents. */
  amount: number
  /** ISO date. */
  date: string
}

type ColumnKey = 'name' | 'status' | 'amount' | 'date'
type Direction = 'asc' | 'desc'

const COLUMNS: { key: ColumnKey; label: string; numeric?: boolean }[] = [
  { key: 'name', label: 'Customer' },
  { key: 'status', label: 'Status' },
  { key: 'amount', label: 'Amount', numeric: true },
  { key: 'date', label: 'Date' },
]

const ROWS: Row[] = [
  { id: '1', name: 'Adam Wathan', email: 'adam@example.com', status: 'Paid', amount: 24999, date: '2026-03-14' },
  { id: '2', name: 'Bessie Cooper', email: 'bessie@example.com', status: 'Pending', amount: 8400, date: '2026-03-11' },
  { id: '3', name: 'Devon Lane', email: 'devon@example.com', status: 'Paid', amount: 129900, date: '2026-02-28' },
  { id: '4', name: 'Arlene McCoy', email: 'arlene@example.com', status: 'Failed', amount: 4200, date: '2026-02-19' },
  { id: '5', name: 'Guy Hawkins', email: 'guy@example.com', status: 'Paid', amount: 37900, date: '2026-02-02' },
  { id: '6', name: 'Kristin Watson', email: 'kristin@example.com', status: 'Pending', amount: 15900, date: '2026-01-22' },
  { id: '7', name: 'Jenny Wilson', email: 'jenny@example.com', status: 'Paid', amount: 89900, date: '2026-01-11' },
  { id: '8', name: 'Priya Nair', email: 'priya@example.com', status: 'Paid', amount: 6200, date: '2025-12-30' },
  { id: '9', name: 'Jonas Keller', email: 'jonas@example.com', status: 'Failed', amount: 19900, date: '2025-12-14' },
  { id: '10', name: 'Riley Quinn', email: 'riley@example.com', status: 'Pending', amount: 44900, date: '2025-11-29' },
]

const money = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' })
const dates = new Intl.DateTimeFormat('en-US', { year: 'numeric', month: 'short', day: 'numeric' })

const STATUS_TONE: Record<Row['status'], string> = {
  Paid: 'bg-emerald-100 text-emerald-800 dark:bg-emerald-500/15 dark:text-emerald-300',
  Pending: 'bg-amber-100 text-amber-800 dark:bg-amber-500/15 dark:text-amber-300',
  Failed: 'bg-rose-100 text-rose-800 dark:bg-rose-500/15 dark:text-rose-300',
}

export default function SortableTableWithSelection({
  rows = ROWS,
  pageSize = 5,
  caption = 'Recent payments',
}: {
  rows?: Row[]
  pageSize?: number
  caption?: string
}) {
  const [sortKey, setSortKey] = useState<ColumnKey>('date')
  const [direction, setDirection] = useState<Direction>('desc')
  const [selected, setSelected] = useState<string[]>([])
  const [page, setPage] = useState(0)
  const [announcement, setAnnouncement] = useState('')

  const selectAllRef = useRef<HTMLInputElement>(null)
  const baseId = useId()

  const sorted = useMemo(() => {
    /* A copy: Array.prototype.sort mutates, and sorting the prop in place
       would reorder the caller's array. */
    return [...rows].sort((a, b) => {
      const factor = direction === 'asc' ? 1 : -1
      if (sortKey === 'amount') return (a.amount - b.amount) * factor
      return String(a[sortKey]).localeCompare(String(b[sortKey])) * factor
    })
  }, [rows, sortKey, direction])

  const pageCount = Math.max(1, Math.ceil(sorted.length / pageSize))
  const safePage = Math.min(page, pageCount - 1)
  const visible = sorted.slice(safePage * pageSize, safePage * pageSize + pageSize)

  const visibleIds = visible.map((row) => row.id)
  const selectedOnPage = visibleIds.filter((id) => selected.includes(id))
  const allOnPageSelected = visibleIds.length > 0 && selectedOnPage.length === visibleIds.length
  const someOnPageSelected = selectedOnPage.length > 0 && !allOnPageSelected

  /* indeterminate is a property, not an attribute, so it cannot be expressed
     in JSX and has to be written to the node. Without this the box reads as
     unchecked while three rows are selected. */
  useEffect(() => {
    if (selectAllRef.current) selectAllRef.current.indeterminate = someOnPageSelected
  }, [someOnPageSelected])

  function announce(message: string) {
    setAnnouncement('')
    requestAnimationFrame(() => setAnnouncement(message))
  }

  function toggleSort(key: ColumnKey) {
    const nextDirection: Direction = key === sortKey && direction === 'asc' ? 'desc' : 'asc'
    setSortKey(key)
    setDirection(nextDirection)
    setPage(0)
    const label = COLUMNS.find((column) => column.key === key)?.label ?? key
    announce(`Sorted by ${label}, ${nextDirection === 'asc' ? 'ascending' : 'descending'}.`)
  }

  function toggleRow(row: Row) {
    const next = selected.includes(row.id)
      ? selected.filter((id) => id !== row.id)
      : [...selected, row.id]
    setSelected(next)
    announce(`${row.name} ${next.includes(row.id) ? 'selected' : 'deselected'}. ${next.length} of ${rows.length} selected.`)
  }

  function toggleAllOnPage() {
    const next = allOnPageSelected
      ? selected.filter((id) => !visibleIds.includes(id))
      : [...new Set([...selected, ...visibleIds])]
    setSelected(next)
    announce(
      allOnPageSelected
        ? `Cleared this page. ${next.length} of ${rows.length} selected.`
        : `Selected this page. ${next.length} of ${rows.length} selected.`,
    )
  }

  function goToPage(next: number) {
    const clamped = Math.max(0, Math.min(next, pageCount - 1))
    setPage(clamped)
    announce(`Page ${clamped + 1} of ${pageCount}.`)
  }

  return (
    <div className="bg-gray-50 py-10 dark:bg-gray-950">
      <div className="mx-auto max-w-5xl px-4">
        {/* One region for sorting, selection and paging: none of them moves
            focus, so none of them is otherwise noticed. */}
        <p role="status" aria-live="polite" className="sr-only">
          {announcement}
        </p>

        <div className="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900">
          <div className="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 px-4 py-3 dark:border-white/10">
            <h2 className="text-sm font-semibold text-gray-900 dark:text-white">{caption}</h2>
            {selected.length > 0 && (
              <p className="text-sm text-gray-600 dark:text-gray-300">
                {selected.length} selected
                <button
                  type="button"
                  onClick={() => {
                    setSelected([])
                    announce('Selection cleared.')
                  }}
                  className="ml-3 font-medium text-indigo-700 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-indigo-400"
                >
                  Clear
                </button>
              </p>
            )}
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              {/* Names the table for anyone listing tables on the page. */}
              <caption className="sr-only">
                {caption}, sortable, {rows.length} rows
              </caption>

              <thead className="border-b border-gray-200 dark:border-white/10">
                <tr>
                  <th scope="col" className="w-12 px-4 py-3">
                    <input
                      ref={selectAllRef}
                      type="checkbox"
                      checked={allOnPageSelected}
                      onChange={toggleAllOnPage}
                      aria-label={`Select all ${visible.length} rows on this page`}
                      className="size-4 accent-indigo-600"
                    />
                  </th>

                  {COLUMNS.map((column) => {
                    const active = column.key === sortKey
                    return (
                      <th
                        key={column.key}
                        scope="col"
                        /* On the <th>, and on exactly one column. This is what
                           makes a reader say "sorted ascending" on arrival. */
                        aria-sort={
                          active ? (direction === 'asc' ? 'ascending' : 'descending') : 'none'
                        }
                        className={`px-4 py-3 font-medium text-gray-600 dark:text-gray-300 ${
                          column.numeric ? 'text-right' : ''
                        }`}
                      >
                        {/* A button, not an onClick on the th: a cell with a
                            handler is not focusable and not operable. */}
                        <button
                          type="button"
                          onClick={() => toggleSort(column.key)}
                          className={`inline-flex min-h-11 items-center gap-1.5 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 ${
                            column.numeric ? 'flex-row-reverse' : ''
                          } ${active ? 'text-gray-900 dark:text-white' : ''}`}
                        >
                          {column.label}
                          {active ? (
                            direction === 'asc' ? (
                              <ArrowUp aria-hidden="true" className="size-3.5" />
                            ) : (
                              <ArrowDown aria-hidden="true" className="size-3.5" />
                            )
                          ) : (
                            <ChevronsUpDown
                              aria-hidden="true"
                              className="size-3.5 text-gray-400 dark:text-gray-500"
                            />
                          )}
                          {/* The glyph is decoration. This says what pressing
                              it will do, which the arrow cannot. */}
                          <span className="sr-only">
                            {active && direction === 'asc'
                              ? ', sorted ascending, activate to sort descending'
                              : active
                                ? ', sorted descending, activate to sort ascending'
                                : ', activate to sort ascending'}
                          </span>
                        </button>
                      </th>
                    )
                  })}
                </tr>
              </thead>

              <tbody className="divide-y divide-gray-100 dark:divide-white/5">
                {visible.length === 0 ? (
                  <tr>
                    <td colSpan={COLUMNS.length + 1} className="px-4 py-16 text-center">
                      <SearchX aria-hidden="true" className="mx-auto size-10 text-gray-400" />
                      <p className="mt-3 font-medium text-gray-700 dark:text-gray-200">
                        Nothing to show
                      </p>
                    </td>
                  </tr>
                ) : (
                  visible.map((row) => {
                    const checked = selected.includes(row.id)
                    return (
                      <tr
                        key={row.id}
                        className={checked ? 'bg-indigo-50/60 dark:bg-indigo-500/10' : ''}
                      >
                        <td className="px-4 py-3">
                          <input
                            type="checkbox"
                            checked={checked}
                            onChange={() => toggleRow(row)}
                            /* Named per row. Twelve boxes called "Select row"
                               are twelve identical controls in a list. */
                            aria-label={`Select ${row.name}`}
                            className="size-4 accent-indigo-600"
                          />
                        </td>

                        {/* The row header: this is the cell that identifies the
                            row, so a reader can announce it as context for the
                            cells that follow. */}
                        <th scope="row" className="px-4 py-3 font-normal">
                          <span className="block font-medium text-gray-900 dark:text-white">
                            {row.name}
                          </span>
                          <span className="block text-xs text-gray-500 dark:text-gray-400">
                            {row.email}
                          </span>
                        </th>

                        <td className="px-4 py-3">
                          <span
                            className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${STATUS_TONE[row.status]}`}
                          >
                            {row.status}
                          </span>
                        </td>

                        <td className="px-4 py-3 text-right tabular-nums text-gray-900 dark:text-white">
                          {money.format(row.amount / 100)}
                        </td>

                        <td className="px-4 py-3 text-gray-600 dark:text-gray-300">
                          <time dateTime={row.date}>{dates.format(new Date(row.date))}</time>
                        </td>
                      </tr>
                    )
                  })
                )}
              </tbody>
            </table>
          </div>

          <nav
            aria-label="Table pages"
            className="flex items-center justify-between gap-3 border-t border-gray-200 px-4 py-3 dark:border-white/10"
          >
            <p className="text-sm text-gray-600 dark:text-gray-300">
              Page {safePage + 1} of {pageCount}
            </p>
            <div className="flex gap-2">
              <button
                type="button"
                onClick={() => goToPage(safePage - 1)}
                disabled={safePage === 0}
                className="inline-flex min-h-11 items-center rounded-lg border border-gray-300 px-3 text-sm text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:opacity-40 dark:border-white/15 dark:text-gray-200 dark:hover:bg-white/5"
              >
                Previous
              </button>
              <button
                type="button"
                onClick={() => goToPage(safePage + 1)}
                disabled={safePage >= pageCount - 1}
                className="inline-flex min-h-11 items-center rounded-lg border border-gray-300 px-3 text-sm text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:opacity-40 dark:border-white/15 dark:text-gray-200 dark:hover:bg-white/5"
              >
                Next
              </button>
            </div>
          </nav>
        </div>
      </div>
    </div>
  )
}
