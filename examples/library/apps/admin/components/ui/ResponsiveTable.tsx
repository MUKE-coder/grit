"use client";

import type { ReactNode } from "react";

export interface TableColumn<T> {
  key: string;
  /** Header label. */
  header: string;
  /** Cell renderer. Receives the row. */
  cell: (row: T) => ReactNode;
  /** Hide this column on mobile cards. */
  hideOnMobile?: boolean;
  /** Right-align (numbers, money). */
  align?: "left" | "right";
  /** Fixed pixel width — useful for action / status columns. */
  width?: number;
  /** Make this column shrink-to-fit instead of share remaining space.
   *  Pair with the width prop for tight controls (e.g. row actions). */
  fixed?: boolean;
  /** Override the truncation behaviour. Default: text-ellipsis on overflow.
   *  Set "wrap" to allow line-wrap (e.g. summary / description columns). */
  overflow?: "truncate" | "wrap";
}

interface ResponsiveTableProps<T> {
  columns: TableColumn<T>[];
  rows: T[];
  /** Unique key per row. */
  rowKey: (row: T) => string;
  /** Optional row click handler. */
  onRowClick?: (row: T) => void;
  /** Empty state when rows.length === 0. */
  emptyMessage?: string;
  /** Loading state. */
  loading?: boolean;
}

/**
 * Renders <table> on >=md and a card list on <md. The card view stacks
 * label + value pairs vertically using the column header as the label,
 * which means it stays in sync as columns change without a separate
 * mobile config. Columns flagged hideOnMobile are dropped from cards.
 */
export function ResponsiveTable<T>({
  columns,
  rows,
  rowKey,
  onRowClick,
  emptyMessage = "No records found",
  loading,
}: ResponsiveTableProps<T>) {
  if (loading) {
    // Render a skeleton that mirrors the live table geometry so layout
    // doesn't jump when data arrives. Column count drives the placeholder
    // shape so wide and narrow tables both look right.
    const skeletonCols = columns.length || 4;
    return (
      <div className="animate-pulse">
        <div className="hidden md:block overflow-hidden rounded-xl border border-border bg-bg-elevated">
          <div className="flex gap-4 border-b border-border px-4 py-3">
            {Array.from({ length: skeletonCols }).map((_, i) => (
              <div key={i} className="h-3.5 flex-1 max-w-[120px] rounded bg-bg-hover" />
            ))}
          </div>
          <div className="divide-y divide-border">
            {Array.from({ length: 6 }).map((_, i) => (
              <div key={i} className="flex gap-4 px-4 py-3.5">
                {Array.from({ length: skeletonCols }).map((_, j) => (
                  <div key={j} className="h-3.5 flex-1 rounded bg-bg-hover" />
                ))}
              </div>
            ))}
          </div>
        </div>
        <ul className="md:hidden space-y-3">
          {Array.from({ length: 4 }).map((_, i) => (
            <li key={i} className="rounded-xl border border-border bg-bg-elevated p-4 space-y-2">
              <div className="h-3.5 w-1/2 rounded bg-bg-hover" />
              <div className="h-3.5 w-3/4 rounded bg-bg-hover" />
              <div className="h-3.5 w-1/3 rounded bg-bg-hover" />
            </li>
          ))}
        </ul>
      </div>
    );
  }

  if (rows.length === 0) {
    return (
      <div className="rounded-xl border border-border bg-bg-elevated p-12 text-center text-sm text-text-muted">
        {emptyMessage}
      </div>
    );
  }

  return (
    <>
      {/* Desktop table — table-fixed so column widths follow the config
          and long cells truncate cleanly instead of forcing horizontal
          scroll. Columns without an explicit width share remaining space. */}
      <div className="hidden md:block overflow-x-auto rounded-xl border border-border bg-bg-elevated">
        <table className="w-full table-fixed divide-y divide-border">
          <colgroup>
            {columns.map((c) => (
              <col
                key={c.key}
                style={c.width ? { width: c.width + "px" } : c.fixed ? { width: "1%" } : undefined}
              />
            ))}
          </colgroup>
          <thead>
            <tr>
              {columns.map((c) => (
                <th
                  key={c.key}
                  className={
                    "px-4 py-3 text-xs font-semibold uppercase tracking-wider text-text-muted " +
                    (c.align === "right" ? "text-right" : "text-left")
                  }
                >
                  {c.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {rows.map((row) => (
              <tr
                key={rowKey(row)}
                onClick={onRowClick ? () => onRowClick(row) : undefined}
                className={onRowClick ? "cursor-pointer hover:bg-bg-hover" : ""}
              >
                {columns.map((c) => (
                  <td
                    key={c.key}
                    className={
                      "px-4 py-3 text-sm text-foreground " +
                      (c.align === "right" ? "text-right" : "text-left") + " " +
                      (c.overflow === "wrap" ? "whitespace-normal break-words" : "truncate")
                    }
                  >
                    {c.cell(row)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Mobile cards */}
      <ul className="md:hidden space-y-3">
        {rows.map((row) => (
          <li
            key={rowKey(row)}
            onClick={onRowClick ? () => onRowClick(row) : undefined}
            className={
              "rounded-xl border border-border bg-bg-elevated p-4 " +
              (onRowClick ? "cursor-pointer active:bg-bg-hover" : "")
            }
          >
            <dl className="divide-y divide-border">
              {columns
                .filter((c) => !c.hideOnMobile)
                .map((c) => (
                  <div key={c.key} className="grid grid-cols-3 gap-3 py-2 first:pt-0 last:pb-0">
                    <dt className="col-span-1 text-xs font-medium uppercase tracking-wide text-text-muted">
                      {c.header}
                    </dt>
                    <dd className="col-span-2 text-sm text-foreground">{c.cell(row)}</dd>
                  </div>
                ))}
            </dl>
          </li>
        ))}
      </ul>
    </>
  );
}
