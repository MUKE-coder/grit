import { useEffect, useId, useState } from "react";
import { Input } from "@/components/ui/input";
import type { FieldDefinition } from "@/lib/resource";
import { currencyExponent, fromMajor, toMajor, type Money } from "@repo/shared/types";

interface MoneyFieldProps {
  field: FieldDefinition;
  value: Money | null;
  onChange: (value: Money) => void;
  error?: string;
}

// The currencies offered when a field does not name its own. Deliberately
// short: a shop that trades in one currency should set field.currencies to
// that one currency, and a list of 180 codes is a worse experience than a
// list of eight.
const DEFAULT_CURRENCIES = ["USD", "EUR", "GBP", "UGX", "KES", "NGN", "ZAR", "JPY"];

/** Amount in major units, currency alongside it.
 *
 * The input is a text input rather than type="number" on purpose. A number
 * input silently discards what it considers invalid while the user is still
 * typing, so "19." disappears mid-keystroke and the caret jumps; and its
 * spinner steps by 1, which is wrong for every currency at once.
 */
export function MoneyField({ field, value, onChange, error }: MoneyFieldProps) {
  const fieldId = useId();
  const currencies = field.currencies ?? DEFAULT_CURRENCIES;
  const currency = value?.currency || field.defaultCurrency || currencies[0] || "USD";
  const digits = currencyExponent(currency);

  // What is in the box, held separately from the value that leaves the field.
  //
  // It has to be its own state. Deriving the text from the stored minor units
  // rewrites it on every keystroke: type the "." of 19.99 and the box becomes
  // "19.00" with the caret behind it. Round-tripping the draft only when it
  // disagrees with the value keeps "19." on screen while it is still being
  // typed, and still picks up an edit-mode load or a form reset.
  const [draft, setDraft] = useState(() =>
    value ? toMajor(value).toFixed(digits) : ""
  );

  // biome-ignore lint/correctness/useExhaustiveDependencies: runs on an outside change of value only, never on typing.
  useEffect(() => {
    const shown = Number.parseFloat(draft);
    const held = value ? toMajor(value) : NaN;
    if (!Number.isNaN(shown) && shown === held) return;
    setDraft(value ? toMajor(value).toFixed(digits) : "");
  }, [value]);

  const handleAmount = (raw: string) => {
    // One leading minus, one decimal point, digits. Anything else never
    // reaches the box, so there is no invalid state to report.
    const cleaned = raw
      .replace(/[^0-9.-]/g, "")
      .replace(/(?!^)-/g, "")
      .replace(/^([^.]*\.)|\./g, "$1");
    setDraft(cleaned);
    if (cleaned === "" || cleaned === "-" || cleaned === ".") {
      onChange({ amount: 0, currency });
      return;
    }
    const major = Number.parseFloat(cleaned);
    if (Number.isNaN(major)) return;
    onChange(fromMajor(major, currency));
  };

  // Tidy the draft to the currency's own decimals once the user leaves the
  // box. Doing it any earlier is what makes the field fight back.
  const handleBlur = () => {
    if (draft === "") return;
    setDraft(toMajor(value ?? { amount: 0, currency }).toFixed(digits));
  };

  // Switching currency keeps the figure the user is looking at, not the minor
  // units behind it. Someone typing 19.99 and then picking UGX means 19.99
  // shillings, not the 0.1999 that reinterpreting 1999 minor units would give.
  const handleCurrency = (next: string) => {
    const major = toMajor(value ?? { amount: 0, currency });
    onChange(fromMajor(major, next));
    setDraft(major.toFixed(currencyExponent(next)));
  };

  return (
    <div className="space-y-1.5">
      <label htmlFor={fieldId} className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>

      <div className="flex">
        <Input
          id={fieldId}
          type="text"
          inputMode="decimal"
          autoComplete="off"
          className="rounded-r-none font-mono tabular-nums"
          value={draft}
          onBlur={handleBlur}
          onChange={(e) => handleAmount(e.target.value)}
          placeholder={field.placeholder ?? (0).toFixed(digits)}
          invalid={!!error}
        />
        <select
          aria-label={field.label + " currency"}
          className="rounded-r-lg border border-l-0 border-border bg-bg-tertiary px-3 text-sm text-text-secondary focus:outline-none focus:ring-2 focus:ring-accent"
          value={currency}
          onChange={(e) => handleCurrency(e.target.value)}
        >
          {currencies.map((c) => (
            <option key={c} value={c}>
              {c}
            </option>
          ))}
        </select>
      </div>

      {field.description && !error && (
        <p className="text-xs text-text-muted">{field.description}</p>
      )}
      {!field.description && !error && digits === 0 && (
        <p className="text-xs text-text-muted">
          {currency} has no minor unit, so this is a whole number.
        </p>
      )}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
