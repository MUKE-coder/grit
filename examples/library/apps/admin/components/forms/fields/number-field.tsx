"use client";

// v3.31.38 — NumberField with thousand-separator formatting. As the
// user types "3000" the input visually shows "3,000". The form's
// stored value is still a JS number (or "" for empty); the comma is
// purely a display affordance. Cursor position is preserved across
// reformat so editing in the middle of a number feels natural.
//
// Why type="text" instead of type="number": browsers reject every
// non-digit char (including ",") on number inputs, so the comma
// would never reach the DOM. We use type="text" + inputMode so the
// numeric keyboard still pops up on mobile.

import { useId, useEffect, useRef, useState } from "react";
import type { FieldDefinition } from "@/lib/resource";
import { GenerateButton } from "./generate-button";
import { Input } from "@/components/ui/input";

interface NumberFieldProps {
  field: FieldDefinition;
  value: string | number;
  onChange: (value: number | string) => void;
  error?: string;
  /** When set, renders a Generate button in the label row (see FieldDefinition.generate). */
  onGenerate?: () => void | Promise<void>;
}

export interface FormatOpts {
  allowDecimal: boolean;
  allowNegative: boolean;
}

// formatNumberDisplay strips invalid characters from `raw` (anything
// that's not a digit, optional leading minus, optional single decimal
// point) and inserts thousand-separators into the integer part. The
// decimal portion is kept verbatim so a mid-typed "3000." doesn't
// lose the dot on the way to "3,000.".
export function formatNumberDisplay(raw: string, opts: FormatOpts): string {
  if (raw === "" || raw == null) return "";
  let s = String(raw).replace(opts.allowDecimal ? /[^0-9.\-]/g : /[^0-9\-]/g, "");
  const negative = opts.allowNegative && s.startsWith("-");
  s = s.replace(/-/g, "");
  if (opts.allowDecimal) {
    const parts = s.split(".");
    if (parts.length > 2) {
      s = parts[0] + "." + parts.slice(1).join("");
    }
  } else {
    s = s.replace(/\./g, "");
  }
  if (s === "") return negative ? "-" : "";
  let intPart: string;
  let decPart: string | undefined;
  if (s.includes(".")) {
    const [a, b] = s.split(".");
    intPart = a;
    decPart = b;
  } else {
    intPart = s;
  }
  // Strip multiple leading zeros so "0123" → "123". Keep a single "0"
  // so the user can type "0", "0.", or "0.5" without it vanishing.
  if (intPart.length > 1) {
    intPart = intPart.replace(/^0+/, "") || "0";
  }
  const intFormatted = intPart.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  let out = (negative ? "-" : "") + intFormatted;
  if (decPart !== undefined) out += "." + decPart;
  return out;
}

// parseFormattedNumber strips commas and parses to a JS number. Returns
// "" for the "no value yet" states (empty, lone minus, lone dot) so
// react-hook-form's required validation can still distinguish empty
// from zero.
export function parseFormattedNumber(formatted: string): number | "" {
  if (formatted === "" || formatted === "-" || formatted === "." || formatted === "-.") {
    return "";
  }
  const cleaned = formatted.replace(/,/g, "");
  const n = Number(cleaned);
  return Number.isNaN(n) ? "" : n;
}

export function NumberField({ field, value, onChange, error, onGenerate }: NumberFieldProps) {
  const fieldId = useId();
  const kind = field.numberKind ?? "float";
  const allowDecimal = kind === "float";
  const allowNegative = kind !== "uint";

  const inputRef = useRef<HTMLInputElement>(null);
  const [display, setDisplay] = useState(() =>
    formatNumberDisplay(String(value ?? ""), { allowDecimal, allowNegative })
  );

  // Sync display when external value changes (form reset, edit-mode
  // hydration). Skip when the parsed display already matches -- avoids
  // stomping on mid-edit state like "3000." that parses to 3000.
  // biome-ignore lint/correctness/useExhaustiveDependencies: runs on an outside change of value only, never on typing.
  useEffect(() => {
    const parsed = parseFormattedNumber(display);
    if (parsed === value) return;
    setDisplay(
      value === "" || value == null
        ? ""
        : formatNumberDisplay(String(value), { allowDecimal, allowNegative })
    );
  }, [value]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const input = e.target;
    const raw = input.value;
    const cursorBefore = input.selectionStart ?? raw.length;
    // Count non-comma characters BEFORE the cursor in the user-typed
    // value, so we can place the cursor after the same number of
    // non-comma characters in the reformatted output. That keeps the
    // caret in the same logical position even as commas shift.
    let nonCommasBeforeCursor = 0;
    for (let i = 0; i < cursorBefore; i++) {
      if (raw[i] !== ",") nonCommasBeforeCursor++;
    }
    const formatted = formatNumberDisplay(raw, { allowDecimal, allowNegative });
    setDisplay(formatted);
    onChange(parseFormattedNumber(formatted));
    // Restore cursor after React paints the new value.
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
    <div className="space-y-1.5">
      <div className="flex items-center justify-between gap-2">
        <label htmlFor={fieldId} className="block text-sm font-medium text-foreground">
          {field.label}
          {field.required && <span className="text-danger ml-1">*</span>}
        </label>
        {onGenerate && <GenerateButton onGenerate={onGenerate} />}
      </div>

      <div className="flex">
        {field.prefix && (
          <span className="inline-flex items-center rounded-l-lg border border-r-0 border-border bg-bg-tertiary px-3 text-sm text-text-muted">
            {field.prefix}
          </span>
        )}
        <Input
          id={fieldId}
          ref={inputRef}
          type="text"
          inputMode={allowDecimal ? "decimal" : "numeric"}
          autoComplete="off"
          value={display}
          onChange={handleChange}
          placeholder={field.placeholder}
          invalid={!!error}
          className={field.prefix ? "rounded-l-none" : field.suffix ? "rounded-r-none" : ""}
        />
        {field.suffix && (
          <span className="inline-flex items-center rounded-r-lg border border-l-0 border-border bg-bg-tertiary px-3 text-sm text-text-muted">
            {field.suffix}
          </span>
        )}
      </div>

      {field.description && !error && (
        <p className="text-xs text-text-muted">{field.description}</p>
      )}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
