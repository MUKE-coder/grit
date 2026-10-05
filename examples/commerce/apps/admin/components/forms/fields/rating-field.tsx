"use client";

import { useId, useState } from "react";
import type { FieldDefinition } from "@/lib/resource";

interface RatingFieldProps {
  field: FieldDefinition;
  value: number | null;
  onChange: (value: number | null) => void;
  error?: string;
}

export function Star({ filled, className }: { filled: boolean; className?: string }) {
  return (
    <svg aria-hidden="true" viewBox="0 0 24 24" className={className} fill={filled ? "currentColor" : "none"} stroke="currentColor" strokeWidth="1.5" strokeLinejoin="round">
      <path d="M12 2.5l2.9 6.1 6.6.8-4.9 4.6 1.3 6.6L12 17.3l-5.9 3.3 1.3-6.6-4.9-4.6 6.6-.8z" />
    </svg>
  );
}

/**
 * Stars from 1 to field.max (5 unless the resource says otherwise).
 *
 * The stars are radio buttons, visually hidden behind the icons, so the
 * keyboard works the way it does for any radio group: Tab reaches the group,
 * the arrow keys move the rating, and a screen reader reads "3 stars, 3 of 5".
 */
export function RatingField({ field, value, onChange, error }: RatingFieldProps) {
  const name = useId();
  const legendId = useId();
  const max = field.max ?? 5;
  const [hover, setHover] = useState<number | null>(null);
  const shown = hover ?? value ?? 0;
  return (
    <fieldset className="space-y-1.5" aria-describedby={error ? name + "-error" : undefined}>
      <legend id={legendId} className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </legend>
      <div className="flex items-center gap-3">
        <div className="flex items-center gap-0.5" onMouseLeave={() => setHover(null)}>
          {Array.from({ length: max }, (_, i) => i + 1).map((n) => (
            <label key={n} className="cursor-pointer" onMouseEnter={() => setHover(n)}>
              <input
                type="radio"
                name={name}
                value={n}
                checked={value === n}
                onChange={() => onChange(n)}
                aria-label={n === 1 ? "1 star" : n + " stars"}
                className="peer sr-only"
              />
              <Star
                filled={n <= shown}
                className={
                  "h-6 w-6 rounded-sm transition-colors peer-focus-visible:ring-2 peer-focus-visible:ring-accent " +
                  (n <= shown ? "text-warning" : "text-text-muted")
                }
              />
            </label>
          ))}
        </div>
        <span className="text-xs tabular-nums text-text-muted">{value ? value + " / " + max : "Not rated"}</span>
        {!field.required && value ? (
          <button type="button" onClick={() => onChange(null)} className="text-xs text-text-secondary underline-offset-2 hover:underline">
            Clear
          </button>
        ) : null}
      </div>
      {field.description && !error && <p className="text-xs text-text-muted">{field.description}</p>}
      {error && <p id={name + "-error"} className="text-xs text-danger">{error}</p>}
    </fieldset>
  );
}
