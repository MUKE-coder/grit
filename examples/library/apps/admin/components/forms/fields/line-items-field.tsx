"use client";

import { useMemo, useState, useRef, useEffect } from "react";
import type { FieldDefinition } from "@/lib/resource";
import { Plus, Trash2 } from "@/lib/icons";
import { RelationshipSelectField } from "./relationship-select-field";
import { formatNumberDisplay, parseFormattedNumber } from "./number-field";

interface LineItemsFieldProps {
  field: FieldDefinition;
  value: Record<string, unknown>[];
  onChange: (value: Record<string, unknown>[]) => void;
  error?: string;
}

// Heuristic: if the item has a quantity column and a rate/price column, show a
// derived per-row Total and a grand total. Display-only — only the declared
// columns are submitted; any stored total is the backend's business.
const QTY_RE = /(^|_)(qty|quantity)($|_)/i;
const RATE_RE = /(unit[_-]?rate|unit[_-]?price|(^|_)(rate|price|amount)($|_))/i;

export function LineItemsField({ field, value, onChange, error }: LineItemsFieldProps) {
  const cols = field.itemFields ?? [];
  const rows = Array.isArray(value) ? value : [];
  const noun = field.itemNoun ?? "item";

  const qtyKey = cols.find((c) => QTY_RE.test(c.key))?.key;
  const rateKey = cols.find((c) => RATE_RE.test(c.key))?.key;
  const showTotal = Boolean(qtyKey && rateKey);

  const rowTotal = (row: Record<string, unknown>) => {
    if (!qtyKey || !rateKey) return 0;
    return (Number(row[qtyKey]) || 0) * (Number(row[rateKey]) || 0);
  };
  // biome-ignore lint/correctness/useExhaustiveDependencies: rowTotal is a fresh closure over the column keys on every render.
  const grandTotal = useMemo(
    () => rows.reduce((sum, r) => sum + rowTotal(r), 0),
    [rows]
  );

  const blankRow = () => {
    const r: Record<string, unknown> = {};
    for (const c of cols) r[c.key] = c.type === "number" ? "" : c.defaultValue ?? "";
    return r;
  };
  const addRow = () => onChange([...rows, blankRow()]);
  const removeRow = (i: number) => onChange(rows.filter((_, idx) => idx !== i));
  const setCell = (i: number, key: string, v: unknown) =>
    onChange(rows.map((r, idx) => (idx === i ? { ...r, [key]: v } : r)));

  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between">
        <label className="block text-sm font-medium text-foreground">
          {field.label}
          {field.required && <span className="text-danger ml-1">*</span>}
        </label>
        <button
          type="button"
          onClick={addRow}
          className="inline-flex items-center gap-1 rounded-lg border border-border px-2.5 py-1.5 text-xs font-medium text-accent hover:bg-bg-hover transition-colors"
        >
          <Plus className="h-3.5 w-3.5" /> Add {noun}
        </button>
      </div>

      <div className="overflow-x-auto rounded-xl border border-border">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border bg-bg-tertiary/40 text-left">
              {cols.map((c) => (
                <th key={c.key} className="px-3 py-2 text-xs font-medium uppercase tracking-wide text-text-muted">
                  {c.label}
                </th>
              ))}
              {showTotal && (
                <th className="px-3 py-2 text-right text-xs font-medium uppercase tracking-wide text-text-muted">
                  Total
                </th>
              )}
              <th className="w-10 px-3 py-2" />
            </tr>
          </thead>
          <tbody>
            {rows.length === 0 ? (
              <tr>
                <td
                  colSpan={cols.length + (showTotal ? 2 : 1)}
                  className="px-3 py-6 text-center text-xs text-text-muted"
                >
                  No {noun}s yet — click &ldquo;Add {noun}&rdquo;.
                </td>
              </tr>
            ) : (
              rows.map((row, i) => (
                <tr key={i} className="border-b border-border/60 last:border-b-0">
                  {cols.map((c) => (
                    <td key={c.key} className="px-2 py-1.5 align-top">
                      <LineItemCell
                        col={c}
                        value={row[c.key]}
                        onChange={(v) => setCell(i, c.key, v)}
                      />
                    </td>
                  ))}
                  {showTotal && (
                    <td className="px-3 py-1.5 text-right font-medium text-foreground align-middle">
                      {rowTotal(row).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </td>
                  )}
                  <td className="px-2 py-1.5 text-right align-middle">
                    <button
                      type="button"
                      onClick={() => removeRow(i)}
                      className="rounded-md p-1 text-text-muted hover:bg-bg-hover hover:text-danger transition-colors"
                      aria-label={"Remove " + noun}
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
          {showTotal && rows.length > 0 && (
            <tfoot>
              <tr className="border-t border-border">
                <td colSpan={cols.length} className="px-3 py-2 text-right text-xs font-medium uppercase tracking-wide text-text-muted">
                  Total
                </td>
                <td className="px-3 py-2 text-right text-sm font-semibold text-foreground">
                  {grandTotal.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                </td>
                <td />
              </tr>
            </tfoot>
          )}
        </table>
      </div>
      {field.description && !error && <p className="text-xs text-text-muted">{field.description}</p>}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}

// A single editable cell, rendered by the column's field type. Keeps the inputs
// compact so a row reads like a spreadsheet line, not a stack of form fields.
function LineItemCell({
  col,
  value,
  onChange,
}: {
  col: FieldDefinition;
  value: unknown;
  onChange: (v: unknown) => void;
}) {
  const base =
    "w-full rounded-md border border-border bg-bg-tertiary px-2 py-1.5 text-sm text-foreground outline-none focus:border-accent";

  if (col.type === "number") {
    return <LineItemNumberCell col={col} value={value} onChange={onChange} className={base + " text-right"} />;
  }
  if (col.type === "select") {
    return (
      <select value={String(value ?? "")} onChange={(e) => onChange(e.target.value)} className={base}>
        <option value="">Select…</option>
        {col.options?.map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
    );
  }
  if (col.type === "relationship-select") {
    return (
      <RelationshipSelectField field={col} value={String(value ?? "")} onChange={onChange} />
    );
  }
  if (col.type === "date" || col.type === "datetime") {
    return (
      <input
        type={col.type === "datetime" ? "datetime-local" : "date"}
        value={String(value ?? "")}
        onChange={(e) => onChange(e.target.value)}
        className={base}
      />
    );
  }
  return (
    <input
      type="text"
      value={String(value ?? "")}
      placeholder={col.placeholder ?? col.label}
      onChange={(e) => onChange(e.target.value)}
      className={base}
    />
  );
}

// A comma-formatting number cell — mirrors NumberField's thousand-separator
// behaviour (1000 -> 1,000) inside the line-items table, storing the parsed
// numeric value in form state and keeping the caret put as commas shift.
function LineItemNumberCell({
  col,
  value,
  onChange,
  className,
}: {
  col: FieldDefinition;
  value: unknown;
  onChange: (v: unknown) => void;
  className: string;
}) {
  const kind = col.numberKind ?? "float";
  const opts = { allowDecimal: kind === "float", allowNegative: kind !== "uint" };
  const inputRef = useRef<HTMLInputElement>(null);
  const [display, setDisplay] = useState(() =>
    formatNumberDisplay(value === null || value === undefined ? "" : String(value), opts)
  );

  // biome-ignore lint/correctness/useExhaustiveDependencies: runs on an outside change of value only, never on typing.
  useEffect(() => {
    const parsed = parseFormattedNumber(display);
    if (parsed === value) return;
    setDisplay(value === "" || value == null ? "" : formatNumberDisplay(String(value), opts));
  }, [value]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const input = e.target;
    const raw = input.value;
    const cursorBefore = input.selectionStart ?? raw.length;
    let nonCommasBeforeCursor = 0;
    for (let i = 0; i < cursorBefore; i++) {
      if (raw[i] !== ",") nonCommasBeforeCursor++;
    }
    const formatted = formatNumberDisplay(raw, opts);
    setDisplay(formatted);
    onChange(parseFormattedNumber(formatted));
    requestAnimationFrame(() => {
      if (!inputRef.current) return;
      let pos = 0;
      let counted = 0;
      while (pos < formatted.length && counted < nonCommasBeforeCursor) {
        if (formatted[pos] !== ",") counted++;
        pos++;
      }
      inputRef.current.setSelectionRange(pos, pos);
    });
  };

  return (
    <input
      ref={inputRef}
      type="text"
      inputMode={opts.allowDecimal ? "decimal" : "numeric"}
      autoComplete="off"
      value={display}
      placeholder={col.placeholder}
      onChange={handleChange}
      className={className}
    />
  );
}
