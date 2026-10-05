"use client";

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
