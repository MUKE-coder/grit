package scaffold

// v3.29 form + table primitives.
//
//   components/ui/ResponsiveSheet.tsx — Dialog on >=md, Sheet on mobile.
//                                       Both share the same API so pages
//                                       wrap their form once and it adapts.
//   components/ui/CurrencyInput.tsx   — text input with thousands separators
//                                       that emits a number to the form.
//   components/ui/ResponsiveTable.tsx — table-on-desktop, card-list-on-mobile
//                                       primitive driven by a column config.
//   components/ui/IconButton.tsx      — auto-collapse button: text + icon on
//                                       desktop, icon-only on mobile.
//   lib/export.ts                     — exportToExcel + exportToPDF helpers
//                                       backed by xlsx + jspdf. Loaded
//                                       lazily so they don't bloat the bundle.

// adminResponsiveSheetComponent — renders as a centred Dialog on >=md and
// a bottom-anchored Sheet on mobile. Both share the same open/onClose
// API so pages don't branch on viewport.
func adminResponsiveSheetComponent() string {
	return `"use client";

import { useEffect, useRef, useState } from "react";
import type { ReactNode } from "react";
import { X } from "@/lib/icons";

interface ResponsiveSheetProps {
  open: boolean;
  onClose: () => void;
  title: string;
  description?: string;
  children: ReactNode;
  /** Footer rendered at the bottom (typically Cancel + Submit). */
  footer?: ReactNode;
  /** Max width on desktop. Defaults to 'lg' (~36rem). */
  size?: "sm" | "md" | "lg" | "xl";
}

const sizeClass: Record<NonNullable<ResponsiveSheetProps["size"]>, string> = {
  sm: "md:max-w-sm",
  md: "md:max-w-md",
  lg: "md:max-w-lg",
  xl: "md:max-w-2xl",
};

/**
 * Adapts modal style to viewport. Desktop (>=md): a centred dialog with
 * a slight backdrop blur. Mobile: a bottom-anchored sheet that slides up
 * and stops at 90vh so the user keeps context behind it. Both lock body
 * scroll when open and close on backdrop click + Escape.
 */
export function ResponsiveSheet({
  open,
  onClose,
  title,
  description,
  children,
  footer,
  size = "lg",
}: ResponsiveSheetProps) {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => { if (e.key === "Escape") onClose(); };
    document.addEventListener("keydown", onKey);
    document.body.style.overflow = "hidden";
    return () => {
      document.removeEventListener("keydown", onKey);
      document.body.style.overflow = "";
    };
  }, [open, onClose]);

  if (!open) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-end justify-center md:items-center"
      onClick={(e) => { if (e.target === e.currentTarget) onClose(); }}
    >
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" />

      {/* Panel */}
      <div
        ref={ref}
        role="dialog"
        aria-modal="true"
        aria-labelledby="responsive-sheet-title"
        className={
          "relative z-10 flex w-full flex-col bg-bg-elevated text-foreground shadow-2xl " +
          "rounded-t-2xl md:rounded-2xl " +
          "max-h-[90vh] md:max-h-[85vh] " +
          ("md:w-auto " + sizeClass[size])
        }
      >
        <header className="flex items-start justify-between border-b border-border px-5 py-4">
          <div className="min-w-0">
            <h2 id="responsive-sheet-title" className="text-lg font-semibold truncate">{title}</h2>
            {description && <p className="mt-0.5 text-sm text-text-secondary">{description}</p>}
          </div>
          <button
            type="button"
            onClick={onClose}
            aria-label="Close"
            className="ml-3 inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-text-muted hover:bg-bg-hover hover:text-foreground"
          >
            <X className="h-4 w-4" />
          </button>
        </header>

        <div className="flex-1 overflow-y-auto px-5 py-4">{children}</div>

        {footer && (
          <footer className="flex items-center justify-end gap-2 border-t border-border px-5 py-3">
            {footer}
          </footer>
        )}
      </div>
    </div>
  );
}
`
}

// adminCurrencyInputComponent — text input with auto thousands separators.
// Forwards a normalised numeric value to the form. Designed to plug into
// react-hook-form via Controller but works as a controlled input on its own.
func adminCurrencyInputComponent() string {
	return `"use client";

import { forwardRef, useEffect, useState } from "react";
import type { ChangeEvent, InputHTMLAttributes } from "react";

interface CurrencyInputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, "value" | "onChange" | "type"> {
  /** Numeric value. Pass undefined for empty. */
  value?: number | null;
  /** Fires with the parsed number (or null) on every change. */
  onChange?: (value: number | null) => void;
  /** Currency symbol prefix shown inside the input. Defaults to "$". */
  prefix?: string;
  /** Locale used for thousands separator. Defaults to "en-US". */
  locale?: string;
  /** Allow decimal portion. Defaults to true. */
  allowDecimal?: boolean;
}

/**
 * Formats numeric values with locale thousands separators while storing
 * the canonical number internally. Typing "3000" displays "3,000"; the
 * onChange callback receives 3000. Decimals are preserved when the
 * trailing "." is typed (we hold the raw string so the caret doesn't
 * jump while the user is still typing).
 */
export const CurrencyInput = forwardRef<HTMLInputElement, CurrencyInputProps>(function CurrencyInput(
  { value, onChange, prefix = "$", locale = "en-US", allowDecimal = true, className = "", ...rest },
  ref
) {
  const [display, setDisplay] = useState<string>("");

  // Sync display when value changes externally (form reset, parent edit).
  useEffect(() => {
    if (value === null || value === undefined || Number.isNaN(value)) {
      setDisplay("");
      return;
    }
    setDisplay(formatNumber(value, locale, allowDecimal));
  }, [value, locale, allowDecimal]);

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const raw = e.target.value;
    // Strip everything but digits + one decimal separator.
    const pattern = allowDecimal ? /[^0-9.]/g : /[^0-9]/g;
    const cleaned = raw.replace(pattern, "");

    if (cleaned === "") {
      setDisplay("");
      onChange?.(null);
      return;
    }

    // Preserve trailing "." so the user can keep typing decimals.
    const trailingDot = allowDecimal && cleaned.endsWith(".") && cleaned.indexOf(".") === cleaned.length - 1;

    const numeric = Number(cleaned);
    if (Number.isNaN(numeric)) {
      setDisplay(cleaned);
      return;
    }

    const formatted = trailingDot
      ? formatNumber(Math.floor(numeric), locale, false) + "."
      : formatNumber(numeric, locale, allowDecimal);

    setDisplay(formatted);
    onChange?.(numeric);
  };

  return (
    <div className="relative">
      <span className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-text-muted">
        {prefix}
      </span>
      <input
        {...rest}
        ref={ref}
        type="text"
        inputMode="decimal"
        value={display}
        onChange={handleChange}
        className={
          "w-full rounded-lg border border-border bg-bg-elevated pl-7 pr-3 py-2.5 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent " +
          className
        }
      />
    </div>
  );
});

function formatNumber(n: number, locale: string, allowDecimal: boolean): string {
  const opts: Intl.NumberFormatOptions = allowDecimal
    ? { minimumFractionDigits: 0, maximumFractionDigits: 2 }
    : { maximumFractionDigits: 0 };
  return new Intl.NumberFormat(locale, opts).format(n);
}
`
}

// adminIconButtonComponent — desktop shows label + icon, mobile shows
// icon only. Wrapper around a regular button that adapts its layout via
// Tailwind responsive utilities so pages don't branch on viewport.
func adminIconButtonComponent() string {
	return `"use client";

import { forwardRef } from "react";
import type { ButtonHTMLAttributes, ReactNode } from "react";

interface IconButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  /** Icon node (e.g. <Plus className="h-4 w-4" />). Required. */
  icon: ReactNode;
  /** Text label shown on >=sm screens. The label always serves as the
   *  aria-label on mobile, so screen readers know what the button does. */
  label: string;
  /** Visual variant. Defaults to "primary". */
  variant?: "primary" | "secondary" | "ghost" | "danger";
}

const variantClass: Record<NonNullable<IconButtonProps["variant"]>, string> = {
  primary: "bg-accent text-white hover:bg-accent-hover",
  secondary: "border border-border bg-bg-elevated text-foreground hover:bg-bg-hover",
  ghost: "text-text-secondary hover:bg-bg-hover hover:text-foreground",
  danger: "bg-danger text-white hover:opacity-90",
};

/**
 * Auto-collapsing CTA. Stays text + icon on >=sm; collapses to icon-only
 * on mobile so table rows + page headers don't blow out of the viewport.
 * Pass label as the readable name; it doubles as the aria-label when the
 * text hides.
 */
export const IconButton = forwardRef<HTMLButtonElement, IconButtonProps>(function IconButton(
  { icon, label, variant = "primary", className = "", ...rest },
  ref
) {
  return (
    <button
      {...rest}
      ref={ref}
      aria-label={label}
      className={
        "inline-flex items-center justify-center gap-2 rounded-lg px-3 py-2 text-sm font-medium transition-colors disabled:opacity-50 " +
        // Mobile is icon-only at 36x36; >=sm reveals the label.
        "h-9 w-9 sm:h-9 sm:w-auto sm:px-3.5 " +
        variantClass[variant] + " " +
        className
      }
    >
      {icon}
      <span className="hidden sm:inline">{label}</span>
    </button>
  );
});
`
}

// adminResponsiveTableComponent — turns a column config + rows into a
// table on desktop and a stacked card list on mobile. Pages pass typed
// rows + column accessors; the primitive handles the responsive switch.
func adminResponsiveTableComponent() string {
	return `"use client";

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
    return (
      <div className="rounded-xl border border-border bg-bg-elevated p-12 text-center">
        <div className="inline-block h-6 w-6 animate-spin rounded-full border-2 border-accent border-t-transparent" />
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
      {/* Desktop table */}
      <div className="hidden md:block overflow-x-auto rounded-xl border border-border bg-bg-elevated">
        <table className="min-w-full divide-y divide-border">
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
                      (c.align === "right" ? "text-right" : "text-left")
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
`
}

// adminExportLib — exportToExcel + exportToPDF helpers. xlsx + the
// React-PDF renderer load lazily so they don't bloat the initial bundle;
// pages call the helper, the import fires, and the file downloads. Devs
// supply row data + a column map.
//
// PDF rendering uses @react-pdf/renderer (component-based, JSX) rather
// than jsPDF (imperative). The trade-off is bundle size (~600KB vs ~200KB)
// for the ability to design PDFs as React components — much more flexible
// for the inevitable "add the company letterhead" follow-up.
func adminExportLib() string {
	return `// Export utilities for table data. xlsx + @react-pdf/renderer are heavy
// bundles (~300KB + ~600KB gzipped), so we lazy-import them at call
// time. Pages trigger an export from a button handler; the bundle only
// loads when the user actually exports.
//
// Usage:
//   import { exportToExcel, exportToPDF } from "@/lib/export";
//   const rows = users.map(u => ({ Email: u.email, Name: u.first_name }));
//   await exportToExcel(rows, "users");
//   await exportToPDF(rows, "users", "All Users");
//
// For PDFs with branded headers or non-table layouts, import
// @react-pdf/renderer directly, design your <Document> as JSX, and call
// pdf(doc).toBlob() yourself. exportToPDF here covers the common case.

export type ExportRow = Record<string, string | number | boolean | null>;

/**
 * Download an .xlsx file with the given rows. Each row's keys become
 * column headers. Lazy-loads the xlsx package.
 */
export async function exportToExcel(rows: ExportRow[], filename: string) {
  const XLSX = await import("xlsx");
  const ws = XLSX.utils.json_to_sheet(rows);
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, "Sheet1");
  XLSX.writeFile(wb, sanitize(filename) + ".xlsx");
}

/**
 * Download a .pdf file with the given rows. Renders a simple table.
 * Lazy-loads @react-pdf/renderer + React (the renderer needs createElement
 * at runtime). For richer layouts, design your own <Document> and call
 * pdf(<MyDoc/>).toBlob() directly.
 */
export async function exportToPDF(rows: ExportRow[], filename: string, title?: string) {
  const { Document, Page, View, Text, StyleSheet, pdf } = await import("@react-pdf/renderer");
  const React = await import("react");

  // Inline stylesheet — keeps the helper standalone. Override by writing
  // your own Document component when you need custom typography or
  // letterheads.
  const styles = StyleSheet.create({
    page: { padding: 36, fontFamily: "Helvetica", fontSize: 9, color: "#0f172a" },
    title: { fontSize: 14, fontWeight: 700, marginBottom: 12 },
    table: { width: "auto", borderStyle: "solid", borderColor: "#e2e8f0", borderWidth: 1 },
    row: { flexDirection: "row", borderBottomWidth: 1, borderBottomColor: "#e2e8f0" },
    headRow: { backgroundColor: "#f1f5f9", borderBottomWidth: 1, borderBottomColor: "#e2e8f0", flexDirection: "row" },
    cell: { padding: 6, flex: 1 },
    headCell: { padding: 6, flex: 1, fontWeight: 700 },
    empty: { padding: 12, textAlign: "center", color: "#94a3b8" },
  });

  const headers = rows.length > 0 ? Object.keys(rows[0]) : [];

  // Doc body is built imperatively (createElement) so this helper stays a
  // pure .ts file — no TSX compile step required for the export module.
  const doc = React.createElement(
    Document,
    null,
    React.createElement(
      Page,
      { size: "A4", style: styles.page },
      title ? React.createElement(Text, { style: styles.title }, title) : null,
      rows.length === 0
        ? React.createElement(Text, { style: styles.empty }, "No data")
        : React.createElement(
            View,
            { style: styles.table },
            React.createElement(
              View,
              { style: styles.headRow },
              headers.map((h) =>
                React.createElement(Text, { key: h, style: styles.headCell }, h)
              )
            ),
            rows.map((r, i) =>
              React.createElement(
                View,
                { key: i, style: styles.row },
                headers.map((h) =>
                  React.createElement(
                    Text,
                    { key: h, style: styles.cell },
                    r[h] == null ? "" : String(r[h])
                  )
                )
              )
            )
          )
    )
  );

  const blob = await pdf(doc).toBlob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = sanitize(filename) + ".pdf";
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

/**
 * Parse a user-selected .xlsx or .csv file into rows. Returns the first
 * sheet as an array of objects keyed by column header.
 */
export async function importFromExcel(file: File): Promise<ExportRow[]> {
  const XLSX = await import("xlsx");
  const buf = await file.arrayBuffer();
  const wb = XLSX.read(buf, { type: "array" });
  const firstSheet = wb.SheetNames[0];
  if (!firstSheet) return [];
  return XLSX.utils.sheet_to_json<ExportRow>(wb.Sheets[firstSheet]);
}

function sanitize(name: string): string {
  return name.replace(/[^a-z0-9_-]+/gi, "-").replace(/^-+|-+$/g, "") || "export";
}
`
}
