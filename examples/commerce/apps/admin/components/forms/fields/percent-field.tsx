"use client";

import { useEffect, useId, useState } from "react";
import type { FieldDefinition } from "@/lib/resource";
import { Input } from "@/components/ui/input";

interface PercentFieldProps {
  field: FieldDefinition;
  value: number | null;
  onChange: (value: number | null) => void;
  error?: string;
}

/** A percentage from 0 to 100, typed as a number with the % beside it. */
export function PercentField({ field, value, onChange, error }: PercentFieldProps) {
  const inputId = useId();
  const [draft, setDraft] = useState(value === null || value === undefined ? "" : String(value));

  // biome-ignore lint/correctness/useExhaustiveDependencies: only value; the draft is this field's own
  useEffect(() => {
    if (Number.parseFloat(draft) === value) return;
    setDraft(value === null || value === undefined ? "" : String(value));
  }, [value]);

  const handle = (raw: string) => {
    // Digits and one decimal point; anything else never reaches the box.
    const cleaned = raw.replace(/[^0-9.]/g, "").replace(/^([^.]*\.)|\./g, "$1");
    setDraft(cleaned);
    if (cleaned === "" || cleaned === ".") {
      onChange(null);
      return;
    }
    const n = Number.parseFloat(cleaned);
    if (!Number.isNaN(n)) onChange(n);
  };

  return (
    <div className="space-y-1.5">
      <label htmlFor={inputId} className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      <div className="flex">
        <Input
          id={inputId}
          type="text"
          inputMode="decimal"
          autoComplete="off"
          value={draft}
          onChange={(e) => handle(e.target.value)}
          placeholder={field.placeholder ?? "0"}
          invalid={!!error}
          className="rounded-r-none font-mono tabular-nums"
        />
        <span aria-hidden="true" className="inline-flex items-center rounded-r-lg border border-l-0 border-border bg-bg-tertiary px-3 text-sm text-text-muted">
          %
        </span>
      </div>
      {field.description && !error && <p className="text-xs text-text-muted">{field.description}</p>}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
